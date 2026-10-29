//go:build system

package systemtest

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// landableSpec is the smallest staged spec that a landing accepts: one story, one coverage
// row, and one chunk. header holds any header lines after the status line.
func landableSpec(title, header string) string {
	return "# " + title + "\n\nStatus: staged\n" + header + "\n## User stories\n1. Land source.\n\n### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n| E1 | 1 | lands | command | catches failure |\n" + preflighttest.ChunkTable("E1")
}

func systemLandingRaceFixture(t *testing.T) (root, home, tally, trees, ready, release string) {
	t.Helper()
	var err error
	root, err = os.MkdirTemp(owner.root, "landing-race [journey]-")
	if err != nil {
		t.Fatal(err)
	}
	if result := owner.runAt(root, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init landing race = (%d, %q)", result.code, result.stderr)
	}
	// The landing creates its own commit in this repository, so the identity belongs in
	// the config. A per-command -c leaves the product's commit without an author.
	systemIdentity(t, root)
	home, err = os.MkdirTemp(owner.root, "landing-race [home]-")
	if err != nil {
		t.Fatal(err)
	}
	tally = filepath.Join(home, "gate-tally")
	trees = filepath.Join(home, "gate-trees")
	ready, release = landingBarrier(t, home)
	f := testrepo.NewGateFixture(t.TempDir())
	f.Environment = []string{"LAND_GATE_TALLY", "LAND_GATE_TREES", "LAND_RACE_READY", "LAND_RACE_RELEASE"}
	gate := "set -eu\nruntime=$1\n" + f.Command("grep") + " -q '^Status: implemented$' \"$runtime/specs/x/spec.md\"\ntree=$(" + f.Command("git") + " -C \"$runtime\" write-tree)\nprintf '%s\\n' \"$tree\" >> \"$LAND_GATE_TREES\"\nif [ -f \"$runtime/loser.txt\" ]; then\n  printf l >> \"$LAND_GATE_TALLY\"\n  if [ ! -f \"$runtime/winner.txt\" ]; then\n    " + landingRaceBarrier(f, landingBarrierWindow) + "  fi\nelse\n  printf w >> \"$LAND_GATE_TALLY\"\nfi\n"
	if err := f.Write(root, gate, gate); err != nil {
		t.Fatal(err)
	}
	specBody := landableSpec("x", "")
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("retained-output\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The race writes these files, so the ticket `Writes:` line names them and derives the fence.
	recordtest.Prepare(t, root, 1, "specs/x/spec.md", specBody).SetTicketWrites("loser.txt (new), winner.txt (new)")
	// Each landing worktree publishes production files, so the race delivers an approved spec.
	commitmenttest.SeedAdmission(t, root, "specs/x/spec.md")
	commitmenttest.Commit(t, root, "approve fixture delivery")
	base := systemGitOutput(t, root, "rev-parse", "HEAD")
	systemGit(t, root, "update-ref", "refs/bench/green/main", base)
	return root, home, tally, trees, ready, release
}

// landingBarrier makes the ready and release FIFO pair of one landing barrier in dir. A
// landing waits on the pair that its LAND_RACE_READY and LAND_RACE_RELEASE values name.
func landingBarrier(t *testing.T, dir string) (ready, release string) {
	t.Helper()
	ready, release = filepath.Join(dir, "loser-ready"), filepath.Join(dir, "loser-release")
	for _, fifo := range []string{ready, release} {
		if err := syscall.Mkfifo(fifo, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return ready, release
}

// openChunkGaps is the shared open-gap seed under the name that the journey fixtures call.
var openChunkGaps = commitmenttest.OpenChunkGaps

// linkedDeliverable is the staged spec that a linked project's adoption approves.
const linkedDeliverable = "specs/x/spec.md"

// linkedProject is one project that `bench setup` linked before any commitment policy
// exists, and the installed distribution whose promotion broker publishes its landings.
type linkedProject struct {
	root, home, wrapper string
	install             landRouteInstall
}

// newLinkedProject links a disposable project that holds one staged spec and no board, as
// `bench setup` leaves it. Its gate is a plain green script, so each landing grades
// commitment authority and not project checks. The installed broker is the selected
// executable, bound by the manifest that the land route authenticates.
func newLinkedProject(t *testing.T) linkedProject {
	t.Helper()
	root := journeyRepo(t, "linked-adoption [project]-")
	systemIdentity(t, root)
	home, err := os.MkdirTemp(owner.root, "linked-adoption [home]-")
	if err != nil {
		t.Fatal(err)
	}
	commitmenttest.Write(t, root, "go.mod", "module example.com/linked\n\ngo 1.22\n")
	commitmenttest.Write(t, root, linkedDeliverable, commitmenttest.StagedBody)
	project := linkedProject{root: root, home: home, wrapper: filepath.Join(root, ".bench", "bin", "bench.sh"), install: newLandRouteInstall(t)}
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	setup := owner.runAt(root, project.env(), owner.selected.path, "setup", "--yes")
	if setup.code != 0 {
		t.Fatalf("linked project setup = (%d, %q, %q)", setup.code, setup.stdout, setup.stderr)
	}
	testrepo.NewGateFixture(t.TempDir()).MustWrite(t, root, "exit 0\n", "")
	commitmenttest.Commit(t, root, "link the project")
	broker, err := os.ReadFile(owner.selected.path)
	if err != nil {
		t.Fatal(err)
	}
	project.install.bindBroker(t, broker)
	return project
}

func (p linkedProject) env() []string {
	return []string{"BENCH_HOME=" + p.home, "BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
}

// run launches the project's installed wrapper in dir.
func (p linkedProject) run(t *testing.T, dir string, args ...string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runAt(dir, p.env(), "bash", append([]string{p.wrapper}, args...)...)
}

// land publishes source through the installed broker from the operator's own shell, so
// no inherited routing override reaches the land route.
func (p linkedProject) land(t *testing.T, source systemLandingWorktree, base string) processResult {
	t.Helper()
	return owner.runAt(p.root, []string{"BENCH_HOME=" + p.home, "BENCH_KIT", "BENCH_RUN_BINARY", "BENCH_WRAPPER"}, "bash", p.install.wrapper,
		"worktree", "land", "--request", source.request, "--base", base, "--source-tip", source.tip, "-m", "publish "+source.request, source.path)
}

// refuseLanding requires the installed broker to refuse source with want before it
// publishes anything.
func (p linkedProject) refuseLanding(t *testing.T, source systemLandingWorktree, base, want string) {
	t.Helper()
	before := systemGitOutput(t, p.root, "show-ref", "--head")
	refused := p.land(t, source, base)
	if refused.code != 1 || !strings.Contains(refused.stdout, want) {
		t.Fatalf("landing %s = (%d, %q, %q), want a refusal naming %q", source.request, refused.code, refused.stdout, refused.stderr, want)
	}
	if after := systemGitOutput(t, p.root, "show-ref", "--head"); after != before {
		t.Fatalf("refused landing %s moved refs: got %q want %q", source.request, after, before)
	}
}

// publish requires the installed broker to land source.
func (p linkedProject) publish(t *testing.T, source systemLandingWorktree, base string) {
	t.Helper()
	landed := p.land(t, source, base)
	if landed.code != 0 || !strings.Contains(landed.stdout, "landed{") {
		t.Fatalf("landing %s = (%d, %q, %q)", source.request, landed.code, landed.stdout, landed.stderr)
	}
}

// commitSource records the planning or production change in source with plain Git, which no
// Bench guard grades, so only the installed broker decides its publication.
func commitSource(t *testing.T, source *systemLandingWorktree, message string) {
	t.Helper()
	commitmenttest.Commit(t, source.path, message)
	source.tip = systemGitOutput(t, source.path, "rev-parse", "HEAD")
}

// propose writes the adoption policy that approves linkedDeliverable into planning and
// commits it with no approval.
func (p linkedProject) propose(t *testing.T, planning *systemLandingWorktree) {
	t.Helper()
	commitmenttest.SeedAdmission(t, planning.path, linkedDeliverable)
	commitSource(t, planning, "propose the adoption policy")
}

// approve plans the proposed policy and approves that exact plan through the installed
// wrapper. The project has no board, so approval stages the committed proposal alone and
// leaves the planning checkout clean.
func (p linkedProject) approve(t *testing.T, planning *systemLandingWorktree) {
	t.Helper()
	planned := p.run(t, planning.path, "commitment", "plan", "--input", filepath.FromSlash(commitment.PolicyPath))
	match := regexp.MustCompile(`(?m)^  "(sha256:[0-9a-f]{64})",`).FindStringSubmatch(planned.stdout)
	if planned.code != 0 || match == nil {
		t.Fatalf("adoption plan = (%d, %q, %q)", planned.code, planned.stdout, planned.stderr)
	}
	approved := p.run(t, planning.path, "commitment", "approve", "--plan", match[1], "--decision", "reviewer adoption", "--delayed", "none", "--removed", "none")
	if approved.code != 0 {
		t.Fatalf("adoption approval = (%d, %q, %q)", approved.code, approved.stdout, approved.stderr)
	}
	if status := systemGitOutput(t, planning.path, "status", "--porcelain", "--untracked-files=all"); status != "" {
		t.Fatalf("approval without a board changed the planning checkout: %q", status)
	}
}

// adopt proposes, approves, and publishes the adoption policy from a fresh planning run.
func (p linkedProject) adopt(t *testing.T) {
	t.Helper()
	base := systemGitOutput(t, p.root, "rev-parse", "HEAD")
	planning := systemCreateLandingWorktree(t, p.root, p.home, "planning", "planning")
	p.propose(t, &planning)
	p.approve(t, &planning)
	p.publish(t, planning, base)
}

// boundedWait prefixes a fixture's shell wait with the ambient timeout, so a fixture that
// its test abandons ends after window. The wait stays in the caller's process group,
// because that group is the registered user that the wait holds.
func boundedWait(f *testrepo.GateFixture, window time.Duration) string {
	return f.Command("timeout") + " --foreground " + strconv.FormatFloat(window.Seconds(), 'f', -1, 64) + " "
}

// landingBarrierWindow bounds a fixture gate's landing barrier. The barrier waits on its own
// test, and the test binary's timeout bounds that test.
var landingBarrierWindow = bounds.TestDeadline(bounds.GoTestTimeout)

// landingRaceBarrier is the shell line where a fixture gate signals LAND_RACE_READY and then
// waits on LAND_RACE_RELEASE. Each FIFO open blocks until its partner opens the other end,
// so the ambient timeout bounds both steps. A gate that loses its partner ends after window.
func landingRaceBarrier(f *testrepo.GateFixture, window time.Duration) string {
	return boundedWait(f, window) + f.Command("sh") + " -c 'printf r > \"$1\" && IFS= read -r _ < \"$2\"' barrier \"$LAND_RACE_READY\" \"$LAND_RACE_RELEASE\"\n"
}

// artifactOracleFile names the file under the fixture home where the artifact gate publishes
// its process ID. The gate leads the oracle group, so that ID is the group's ID.
const artifactOracleFile = "artifact-oracle.pid"

func configureArtifactLandingFixture(t *testing.T, root string) {
	t.Helper()
	f := testrepo.NewGateFixture(t.TempDir())
	git := f.Command("git")
	// The selected executable resolves Go build inputs before the phase starts.
	f.Command("go")
	gate := "set -eu\nroot=${1:-$(" + git + " rev-parse --show-toplevel)}\nkit=${BENCH_KIT:?}\nbench=${BENCH_RUN_BINARY:?}\nprintf '%s\\n' \"$$\" > \"$HOME/" + artifactOracleFile + "\"\nexec " + f.Command("env") + " BENCH_KIT=\"$kit\" BENCH_RUN_BINARY=\"$bench\" \"$bench\" gate-phases \"$root\"\n"
	barrier := landingRaceBarrier(f, landingBarrierWindow)
	phase := "#!/bin/sh\nset -eu\nruntime=$(" + git + " rev-parse --show-toplevel)\n" + f.Command("grep") + " -q '^Status: implemented$' \"$runtime/specs/x/spec.md\"\ntree=$(" + git + " -C \"$runtime\" write-tree)\nprintf '%s\\n' \"$tree\" >> \"$LAND_GATE_TREES\"\nif [ -f \"$runtime/loser.txt\" ]; then\n  printf l >> \"$LAND_GATE_TALLY\"\nelse\n  printf w >> \"$LAND_GATE_TALLY\"\nfi\n" + barrier
	if err := os.WriteFile(filepath.Join(root, ".bench", "landing-race-phase.sh"), []byte(f.Script(phase)), 0o755); err != nil {
		t.Fatal(err)
	}
	phases := "{\"phases\":[{\"name\":\"landing-race\",\"argv\":[\".bench/landing-race-phase.sh\"]}]}\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "phases.json"), []byte(phases), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "cmd", "bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module landingrace\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// The stub reads the two positionals off the end and creates the output directory, as
	// the real go-build.sh does; that script also takes options ahead of the positionals.
	dirname := f.Command("dirname")
	build := "set -eu\nshift $(($# - 2))\nroot=$1\nout=$2\nstaged=$out.staged\n" + f.Command("mkdir") + " -p \"$(" + dirname + " \"$out\")\"\n" + f.Command("cp") + " \"$LAND_BASELINE_BENCH\" \"$staged\"\n" + f.Command("chmod") + " 0700 \"$staged\"\n\"$staged\" freshness-publish \"$root\" \"$out\" \"$(" + dirname + " \"$out\")\" 1.2.3\n"
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.sh"), []byte(f.Script(build)), 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(root, filepath.FromSlash(freshness.BuildInputsManifest))
	if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, []byte(freshness.BuildInputLine("build_script", "scripts/go-build.sh")), 0o644); err != nil {
		t.Fatal(err)
	}
	f.Environment = []string{"HOME", "LAND_BASELINE_BENCH", "LAND_GATE_TALLY", "LAND_GATE_TREES", "LAND_RACE_READY", "LAND_RACE_RELEASE"}
	if err := f.Write(root, gate, gate); err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "add", ".")
	systemGit(t, root, systemCommitArgs("-qm", "artifact recovery fixture")...)
	base := systemGitOutput(t, root, "rev-parse", "HEAD")
	systemGit(t, root, "update-ref", "refs/bench/green/main", base)
}

// releaseArtifactBarrier releases the phase that waits on release. The phase opens release
// directly after its barrier, so the release waits a bounded window for that reader.
func releaseArtifactBarrier(t *testing.T, release string) {
	t.Helper()
	window := bounds.TestDeadline(0)
	for deadline := time.Now().Add(window); ; time.Sleep(10 * time.Millisecond) {
		// A nonblocking open of a FIFO that no reader holds fails at once with ENXIO.
		releaseWriter, err := os.OpenFile(release, os.O_WRONLY|syscall.O_NONBLOCK, 0)
		if err == nil {
			if _, err := releaseWriter.Write([]byte("go\n")); err != nil {
				t.Fatal(err)
			}
			if err := releaseWriter.Close(); err != nil {
				t.Fatal(err)
			}
			return
		}
		if !errors.Is(err, syscall.ENXIO) {
			t.Fatal(err)
		}
		if time.Now().After(deadline) {
			t.Fatal(bounds.TestTimeoutVerdict("the phase to open its release barrier", window))
		}
	}
}

// artifactOracle reads the oracle group that the artifact gate published under home.
func artifactOracle(t *testing.T, home string) int {
	t.Helper()
	published, err := os.ReadFile(filepath.Join(home, artifactOracleFile))
	if err != nil {
		t.Fatal(err)
	}
	oracle, err := strconv.Atoi(strings.TrimSpace(string(published)))
	if err != nil {
		t.Fatal(err)
	}
	return oracle
}

// releaseArtifactOracle releases the phase of a landing that died, and waits until the
// landing's surviving oracle group ends. That group is a registered bundle user, and it
// writes under the fixture home until it ends.
func releaseArtifactOracle(t *testing.T, release string, oracle int) {
	t.Helper()
	releaseArtifactBarrier(t, release)
	processtest.WaitAbsent(t, oracle)
}

// artifactBarrierWindow bounds a landing's way to its prospective gate barrier. The landing
// reads git worktree state on the way there, so that read is the longest bound it contains.
var artifactBarrierWindow = bounds.TestDeadline(bounds.WorktreeListTimeout)

// awaitArtifactBarrier waits for a landing subprocess to open its prospective gate barrier.
func awaitArtifactBarrier(t *testing.T, ready string, done <-chan error, stdout, stderr *bytes.Buffer) {
	t.Helper()
	window := artifactBarrierWindow
	observed := make(chan error, 1)
	go func() {
		file, err := os.Open(ready)
		if err != nil {
			observed <- err
			return
		}
		defer file.Close()
		byteRead := make([]byte, 1)
		_, err = file.Read(byteRead)
		if err == nil && string(byteRead) != "r" {
			err = errors.New("unexpected prospective gate barrier")
		}
		observed <- err
	}()
	select {
	case err := <-observed:
		if err != nil {
			t.Fatal(err)
		}
	case err := <-done:
		t.Fatalf("prospective owner ended before the artifact barrier: %v, stdout=%q, stderr=%q", err, stdout.String(), stderr.String())
	case <-time.After(window):
		t.Fatal(bounds.TestTimeoutVerdict("the prospective owner to reach the artifact barrier", window))
	}
}
