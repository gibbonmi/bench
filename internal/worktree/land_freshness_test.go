// Stable-owner tests for the landing command: the public landing runs entirely under
// the invoked owner process; repository executables, their seals, and candidate build
// code never join the promotion.
package worktree

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	"github.com/gibbonmi/bench/internal/landing"
)

// TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion is SOL01 and LC6. The
// candidate tree carries its own build entry point, a go-build.sh that records the
// destination branch it observed. The refresh effect runs that script once, and nothing
// before it does: the branch the script read is the published commit, so the candidate
// code reached no step that decided which bytes were published.
func TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion(t *testing.T) {
	t.Parallel()
	request := "land-owner-no-candidate-code"
	f := publicLandingFixture(t, request, "", "")
	marker := filepath.Join(t.TempDir(), "candidate-ran")
	commitLandingBuildInputs(t, f.root)
	mustWrite(t, filepath.Join(f.root, "scripts", "go-build.sh"), []byte("#!/bin/sh\ngit rev-parse main > "+marker+"\nexit 1\n"), 0o755)
	gitRun(t, f.root, "add", "scripts/go-build.sh")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "candidate build entry")
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, wantEffects("failed")) {
		t.Fatalf("stable-owner landing = (%d, %q, %q), want a failed refresh", r.exit, r.stdout, r.stderr)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	if got := strings.TrimSpace(fixtureFileText(t, marker)); got != published {
		t.Fatalf("candidate build entry observed branch %q, want the published commit %q", got, published)
	}
	if strings.Contains(r.stderr, "rebuilt") {
		t.Fatalf("landing rebuilt an executable: %q", r.stderr)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want one prospective run", got, err)
	}
}

// TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease is SOL04 and LC7. The
// invoked owner carries the complete landing — publication, marker, reconcile, release,
// and the effects — in one process. A rebuild-and-re-exec path would either replace the
// process or surface its rebuild disclosure; neither may appear. This destination
// declares build inputs but carries no build entry point, so the refresh reports failed
// under that same owner rather than handing the landing to another process.
func TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease(t *testing.T) {
	t.Parallel()
	request := "land-owner-single-process"
	f := publicLandingFixture(t, request, "", "")
	commitLandingBuildInputs(t, f.root)
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	ownerPid := os.Getpid()
	r := runVerb(t, verbLand, f.call(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, wantEffects("failed")) {
		t.Fatalf("single-owner landing = (%d, %q, %q), want a failed refresh", r.exit, r.stdout, r.stderr)
	}
	if os.Getpid() != ownerPid {
		t.Fatalf("owner process identity changed: %d -> %d", ownerPid, os.Getpid())
	}
	if strings.Contains(r.stderr, "rebuilt") {
		t.Fatalf("owner re-executed through a rebuild: %q", r.stderr)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want exactly one gate under one owner", got, err)
	}
	if _, err := os.Stat(f.creation.Path); !os.IsNotExist(err) {
		t.Fatalf("release did not complete under the owner process: %v", err)
	}
}

// TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal is SOL17 and LC4, at the real
// owner through the process seam. A forged dist/bench and adjacent seal sit at the
// primary repository path. The stable owner publishes without consulting them, and the
// refresh then reads that pair rather than running it: the forged seal authenticates
// nothing, and the destination carries no build entry point, so the effect reports
// failed at exit 3 over a published commit that stands.
func TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "land-owner-forged-primary"
	f := publicLandingFixture(t, request, "dist/bench", "dist/")
	commitLandingBuildInputs(t, f.root)
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	marker := filepath.Join(t.TempDir(), "forged-ran")
	mustMkdirAll(t, filepath.Join(f.root, "dist"), 0o755)
	mustWrite(t, filepath.Join(f.root, "dist", "bench"), []byte("#!/bin/sh\nprintf ran > "+marker+"\n"), 0o755)
	mustWrite(t, filepath.Join(f.root, "dist", "bench.seal"), []byte(`{"schema":1,"sources":"`+strings.Repeat("a", 64)+`","executable":"`+strings.Repeat("b", 64)+`"}`), 0o644)

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", f.creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = f.root, &stdout, &stderr
	code := exitCode(cmd.Run())
	if code != 3 || !strings.Contains(stdout.String(), wantEffects("failed")) {
		t.Fatalf("forged-primary landing = (%d, %q, %q), want a failed refresh", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("forged primary executable ran during promotion: %v", err)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got == base {
		t.Fatal("landing published nothing")
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want one prospective run", got, err)
	}
}

// redProspectiveGateLanding is the public landing fixture whose composed prospective
// tree carries a red gate. The refusal it produces is the gate's own verdict, so every
// landing proof ahead of publication has already passed.
func redProspectiveGateLanding(t *testing.T, request string) landingFixture {
	t.Helper()
	f := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(f.root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nprintf g >> '"+f.tally+"'\nexit 1\n"), 0o755)
	gitRun(t, f.root, "add", ".bench/gate-prospective.sh")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "red prospective gate")
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	return f
}

// temporaryProspectiveArtifacts lists the private prospective bundles and gate
// executables left under dir. They are the landing's own temporary storage, so after a
// landing settles either way none may remain.
func temporaryProspectiveArtifacts(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var residue []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prospectiveartifact.BundlePrefix) || strings.HasPrefix(entry.Name(), "bench-run-") {
			residue = append(residue, entry.Name())
		}
	}
	return residue
}

// TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate is SOL11. The
// prospective gate is the only authority that can release the destination update. Its
// red leaves the destination ref, the project-green marker, and the source assignment
// exactly as the landing found them.
func TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate(t *testing.T) {
	t.Parallel()
	request := "land-owner-red-gate"
	f := redProspectiveGateLanding(t, request)
	marker := projectGreenMarker(t, f.root)

	r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.HasPrefix(r.stdout, "refused{") {
		t.Fatalf("red prospective gate = (%d, %q, %q), want a refusal", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != f.base {
		t.Fatalf("red gate published: main = %s, want %s", got, f.base)
	}
	if got := projectGreenMarker(t, f.root); got != marker {
		t.Fatalf("red gate advanced the project-green marker: %q, want %q", got, marker)
	}
	if _, err := os.Stat(f.creation.Path); err != nil {
		t.Fatalf("red gate released the reviewed source: %v", err)
	}
	if got, err := os.ReadFile(f.tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want one prospective run", got, err)
	}
}

// TestLandCommandRemovesEveryTemporaryProspectiveArtifact is SOL15. The owner
// materializes the prospective tree and its gate executable in private temporary
// storage. Neither may outlive the landing, on the published path or on a refusal.
func TestLandCommandRemovesEveryTemporaryProspectiveArtifact(t *testing.T) {
	for _, tc := range []struct {
		name string
		red  bool
		want int
	}{
		{name: "published", want: 0},
		{name: "refused", red: true, want: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "land-owner-residue-" + tc.name
			var f landingFixture
			if tc.red {
				f = redProspectiveGateLanding(t, request)
			} else {
				f = publicLandingFixture(t, request, "", "")
			}
			private := t.TempDir()
			bindEnv(t, "TMPDIR", private)

			r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...))
			if r.exit != tc.want {
				t.Fatalf("landing exit = %d, want %d; stdout=%q stderr=%q", r.exit, tc.want, r.stdout, r.stderr)
			}
			if residue := temporaryProspectiveArtifacts(t, private); len(residue) != 0 {
				t.Fatalf("landing left temporary prospective artifacts %v", residue)
			}
		})
	}
}

// TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing is SOL14. The
// destination update is the commit point: marker, reconcile, and release all run after
// it. Each one's failure resumes to a released landing that composes and publishes
// nothing a second time.
func TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing(t *testing.T) {
	// Each stage breaks its own seam in a copy of the caller's value. The resume then
	// runs under the untouched value, so it faces a working stage exactly as a retry does.
	t.Parallel()
	for _, tc := range []struct {
		name   string
		break_ func(joins) joins
	}{
		{name: "marker", break_: func(j joins) joins {
			j.advanceLandingMarker = func(context.Context, string, string, string, string) error {
				return errors.New("injected marker interruption")
			}
			return j
		}},
		{name: "reconcile", break_: func(j joins) joins {
			j.reconcileLanding = func(joins, string, string, string, string) error {
				return errors.New("injected reconciliation interruption")
			}
			return j
		}},
		{name: "release", break_: func(j joins) joins {
			j.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 1 }
			return j
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := "land-owner-resume-" + tc.name
			f := publicLandingFixture(t, request, "", "")
			publications := 0
			working := defaultJoins()
			oldLand := working.landReviewed
			working.landReviewed = func(ctx context.Context, request landing.ReviewedRequest) (landing.ReviewedResult, error) {
				publications++
				return oldLand(ctx, request)
			}
			broken := tc.break_(working)

			r := runVerb(t, verbLand, f.callWith(broken, landArgs(request, f.base, f.tip, f.creation.Path)...))
			if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:"+tc.name) {
				t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			published := gitOutput(t, f.root, "rev-parse", "main")

			args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}
			r = runVerb(t, verbLand, f.callWith(working, args...))
			if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
				t.Fatalf("resume = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
			}
			if got := gitOutput(t, f.root, "rev-parse", "main"); got != published {
				t.Fatalf("resume republished: main = %s, want %s", got, published)
			}
			if publications != 1 {
				t.Fatalf("landing compositions = %d, want exactly one publication", publications)
			}
		})
	}
}

// TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate is SOL10 at the
// owner-to-gate transport. The schedule-resolution tests in internal/gate name the
// baseline themselves, so they grade the manifest lookup alone. Only a landing proves
// the owner actually hands the destination to the gate it started; an omission there
// silently returns phase selection to the candidate tree.
func TestLandCommandCarriesTheBaselineScheduleRootIntoTheProspectiveGate(t *testing.T) {
	t.Parallel()
	request := "land-owner-baseline-transport"
	f := publicLandingFixture(t, request, "", "")
	recorded := filepath.Join(t.TempDir(), "baseline")
	mustWrite(t, filepath.Join(f.root, ".bench", "gate-prospective.sh"),
		[]byte("#!/bin/sh\nset -eu\nprintf '%s' \"${BENCH_GATE_BASELINE:-}\" > "+recorded+"\nprintf g >> '"+f.tally+"'\n"), 0o755)
	gitRun(t, f.root, "add", ".bench/gate-prospective.sh")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "record the baseline schedule root")
	base := gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 0 {
		t.Fatalf("landing = (%d, %q, %q), want a released landing", r.exit, r.stdout, r.stderr)
	}
	got := strings.TrimSpace(fixtureFileText(t, recorded))
	if got == "" {
		t.Fatal("the landing owner handed the prospective gate no baseline schedule root")
	}
	if !sameDirectoryAs(t, got, f.root) {
		t.Fatalf("baseline schedule root = %q, want the landing destination %q", got, f.root)
	}
}

// fixtureFileText reads one recorded fixture value.
func fixtureFileText(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

// sameDirectoryAs compares two directory paths by file identity, so a differently
// spelled path to the same destination still matches.
func sameDirectoryAs(t *testing.T, a, b string) bool {
	t.Helper()
	ai, err := os.Stat(a)
	if err != nil {
		return false
	}
	bi, err := os.Stat(b)
	if err != nil {
		t.Fatal(err)
	}
	return os.SameFile(ai, bi)
}
