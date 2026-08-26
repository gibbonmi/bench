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

	"github.com/gibbonmi/bench/internal/landing"
)

// TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion is SOL01. The
// candidate tree carries its own build entry point, a go-build.sh that records a
// marker. A candidate-owned promotion rebuilds through that script and re-runs the
// landing under the result; the stable owner completes the landing without ever
// executing it.
func TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion(t *testing.T) {
	t.Parallel()
	request := "land-owner-no-candidate-code"
	root, creation, _, _, tally, home := publicLandingFixture(t, request, "", "")
	marker := filepath.Join(t.TempDir(), "candidate-ran")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	mustWrite(t, filepath.Join(root, "scripts", "go-build.sh"), []byte("#!/bin/sh\nprintf ran > "+marker+"\nexit 1\n"), 0o755)
	gitRun(t, root, "add", "scripts/go-build.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "candidate build entry")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, filepath.Join(root, "dist", "bench"), landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("stable-owner landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("candidate landing code ran during its own promotion: %v", err)
	}
	if strings.Contains(stderr.String(), "rebuilt") {
		t.Fatalf("landing rebuilt an executable: %q", stderr.String())
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want one prospective run", got, err)
	}
}

// TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease is SOL04. The
// invoked owner carries the complete landing — publication, marker, reconcile, and
// release — in one process. A rebuild-and-re-exec path would either replace the
// process or surface its rebuild disclosure; neither may appear.
func TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease(t *testing.T) {
	t.Parallel()
	request := "land-owner-single-process"
	root, creation, _, _, tally, home := publicLandingFixture(t, request, "", "")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	ownerPid := os.Getpid()
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, filepath.Join(root, "dist", "bench"), landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("single-owner landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	if os.Getpid() != ownerPid {
		t.Fatalf("owner process identity changed: %d -> %d", ownerPid, os.Getpid())
	}
	if strings.Contains(stderr.String(), "rebuilt") {
		t.Fatalf("owner re-executed through a rebuild: %q", stderr.String())
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want exactly one gate under one owner", got, err)
	}
	if _, err := os.Stat(creation.Path); !os.IsNotExist(err) {
		t.Fatalf("release did not complete under the owner process: %v", err)
	}
}

// TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal is SOL17, at the real owner
// through the process seam. A forged dist/bench and adjacent seal sit at the primary
// repository path; the stable owner completes the landing without consulting or
// executing them.
func TestLandCommandIgnoresAForgedPrimaryExecutableAndSeal(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "land-owner-forged-primary"
	root, creation, _, _, tally, _ := publicLandingFixture(t, request, "dist/bench", "dist/")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	marker := filepath.Join(t.TempDir(), "forged-ran")
	mustMkdirAll(t, filepath.Join(root, "dist"), 0o755)
	mustWrite(t, filepath.Join(root, "dist", "bench"), []byte("#!/bin/sh\nprintf ran > "+marker+"\n"), 0o755)
	mustWrite(t, filepath.Join(root, "dist", "bench.seal"), []byte(`{"schema":1,"sources":"`+strings.Repeat("a", 64)+`","executable":"`+strings.Repeat("b", 64)+`"}`), 0o644)

	var stdout, stderr bytes.Buffer
	cmd := descendant(t, binary, "worktree", "land", "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path)
	cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
	code := exitCode(cmd.Run())
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("forged-primary landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("forged primary executable ran during promotion: %v", err)
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got == base {
		t.Fatal("landing published nothing")
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("gate tally = %q, %v, want one prospective run", got, err)
	}
}

// TestLandCommandReportsInstallStepForABrokerChangingDiff is SOL16. A reviewed diff
// that changes the promotion broker's own build inputs lands as source, but the
// installed broker keeps authority. The landing must name the install step so the
// operator does not expect source publication to replace it.
func TestLandCommandReportsInstallStepForABrokerChangingDiff(t *testing.T) {
	t.Parallel()
	request := "land-owner-broker-change"
	root, creation, _, _, _, home := publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(root, "go.mod"), []byte("module benchfixture\n\ngo 1.22\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "cmd", "bench"), 0o755)
	mustWrite(t, filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "scripts"), 0o755)
	mustWrite(t, filepath.Join(root, "scripts", "go-build.sh"), []byte("#!/bin/sh\nexit 0\n"), 0o755)
	mustWrite(t, filepath.Join(root, "scripts", "go-build.inputs"), []byte("build_script=scripts/go-build.sh\n"), 0o644)
	spec := filepath.Join(root, "specs", "x", "spec.md")
	body, err := os.ReadFile(spec)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, spec, append(body, []byte("- `scripts/go-build.sh`\n")...), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "broker build inputs")
	gitRun(t, creation.Path, "rebase", "main")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "scripts/go-build.sh", "#!/bin/sh\n# next broker\nexit 0\n", "change broker source")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
		t.Fatalf("broker-changing landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "bench repair") {
		t.Fatalf("broker-changing landing named no install step: %q", stderr.String())
	}
}

// redProspectiveGateLanding is the public landing fixture whose composed prospective
// tree carries a red gate. The refusal it produces is the gate's own verdict, so every
// landing proof ahead of publication has already passed.
func redProspectiveGateLanding(t *testing.T, request string) (root string, creation Creation, base, tip, tally, home string) {
	t.Helper()
	root, creation, _, _, tally, home = publicLandingFixture(t, request, "", "")
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"), []byte("#!/bin/sh\nprintf g >> '"+tally+"'\nexit 1\n"), 0o755)
	gitRun(t, root, "add", ".bench/gate-prospective.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "red prospective gate")
	base = gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	return root, creation, base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), tally, home
}

// temporaryProspectiveArtifacts lists the private prospective checkouts and gate
// executables left under dir. Both are the landing's own temporary storage, so after a
// landing settles either way none may remain.
func temporaryProspectiveArtifacts(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var residue []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "bench-gate-subject-") || strings.HasPrefix(entry.Name(), "bench-run-") {
			residue = append(residue, entry.Name())
		}
	}
	return residue
}

// projectGreenMarker reads the destination's project-green marker, answering the empty
// string when no marker is recorded. An absent marker is an ordinary state before the
// first landing, so it is a value here rather than a test failure.
func projectGreenMarker(t *testing.T, root string) string {
	t.Helper()
	output, err := descendant(t, "git", "-C", root, "rev-parse", "--verify", "--quiet", "refs/bench/green/main").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate is SOL11. The
// prospective gate is the only authority that can release the destination update. Its
// red leaves the destination ref, the project-green marker, and the source assignment
// exactly as the landing found them.
func TestLandCommandLeavesTheDestinationUnchangedAfterARedProspectiveGate(t *testing.T) {
	t.Parallel()
	request := "land-owner-red-gate"
	root, creation, base, tip, tally, home := redProspectiveGateLanding(t, request)
	marker := projectGreenMarker(t, root)

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 1 || !strings.HasPrefix(stdout.String(), "refused{") {
		t.Fatalf("red prospective gate = (%d, %q, %q), want a refusal", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != base {
		t.Fatalf("red gate published: main = %s, want %s", got, base)
	}
	if got := projectGreenMarker(t, root); got != marker {
		t.Fatalf("red gate advanced the project-green marker: %q, want %q", got, marker)
	}
	if _, err := os.Stat(creation.Path); err != nil {
		t.Fatalf("red gate released the reviewed source: %v", err)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
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
			var root, base, tip, home string
			var creation Creation
			if tc.red {
				root, creation, base, tip, _, home = redProspectiveGateLanding(t, request)
			} else {
				root, creation, base, tip, _, home = publicLandingFixture(t, request, "", "")
			}
			private := t.TempDir()
			bindEnv(t, "TMPDIR", private)

			var stdout, stderr bytes.Buffer
			if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != tc.want {
				t.Fatalf("landing exit = %d, want %d; stdout=%q stderr=%q", code, tc.want, stdout.String(), stderr.String())
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
			root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
			publications := 0
			working := defaultJoins()
			oldLand := working.landReviewed
			working.landReviewed = func(ctx context.Context, request landing.ReviewedRequest) (landing.ReviewedResult, error) {
				publications++
				return oldLand(ctx, request)
			}
			broken := tc.break_(working)

			var stdout, stderr bytes.Buffer
			if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:"+tc.name) {
				t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			published := gitOutput(t, root, "rev-parse", "main")

			stdout.Reset()
			stderr.Reset()
			args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
			if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released,census=0}") {
				t.Fatalf("resume = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			if got := gitOutput(t, root, "rev-parse", "main"); got != published {
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
	root, creation, _, _, tally, home := publicLandingFixture(t, request, "", "")
	recorded := filepath.Join(t.TempDir(), "baseline")
	mustWrite(t, filepath.Join(root, ".bench", "gate-prospective.sh"),
		[]byte("#!/bin/sh\nset -eu\nprintf '%s' \"${BENCH_GATE_BASELINE:-}\" > "+recorded+"\nprintf g >> '"+tally+"'\n"), 0o755)
	gitRun(t, root, "add", ".bench/gate-prospective.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "record the baseline schedule root")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 0 {
		t.Fatalf("landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	got := strings.TrimSpace(fixtureFileText(t, recorded))
	if got == "" {
		t.Fatal("the landing owner handed the prospective gate no baseline schedule root")
	}
	if !sameDirectoryAs(t, got, root) {
		t.Fatalf("baseline schedule root = %q, want the landing destination %q", got, root)
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
