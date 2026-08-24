// Stable-owner tests for the landing command: the public landing runs entirely under
// the invoked owner process; repository executables, their seals, and candidate build
// code never join the promotion.
package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion is SOL01. The
// candidate tree carries its own build entry point, a go-build.sh that records a
// marker. A candidate-owned promotion rebuilds through that script and re-runs the
// landing under the result; the stable owner completes the landing without ever
// executing it.
func TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion(t *testing.T) {
	request := "land-owner-no-candidate-code"
	root, creation, _, _, tally := publicLandingFixture(t, request, "", "")
	marker := filepath.Join(t.TempDir(), "candidate-ran")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	mustWrite(t, filepath.Join(root, "scripts", "go-build.sh"), []byte("#!/bin/sh\nprintf ran > "+marker+"\nexit 1\n"), 0o755)
	gitRun(t, root, "add", "scripts/go-build.sh")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "candidate build entry")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, filepath.Join(root, "dist", "bench"), landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
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
	request := "land-owner-single-process"
	root, creation, _, _, tally := publicLandingFixture(t, request, "", "")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	ownerPid := os.Getpid()
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, filepath.Join(root, "dist", "bench"), landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
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
	binary := testRunBinary(t)
	request := "land-owner-forged-primary"
	root, creation, _, _, tally := publicLandingFixture(t, request, "dist/bench", "dist/")
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
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
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
	request := "land-owner-broker-change"
	root, creation, _, _, _ := publicLandingFixture(t, request, "", "")
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
	code := LandCommand(root, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
		t.Fatalf("broker-changing landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "bench repair") {
		t.Fatalf("broker-changing landing named no install step: %q", stderr.String())
	}
}
