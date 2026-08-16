package skillsindex

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// The verb's mode contract: --check (the default) reports drift and exits 1, the
// conflicting pair is refused with usage and touches nothing, --write clears the drift
// and exits 0, and a following --check on the freshly written tree reports nothing.
func TestCommand(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
	run(t, root, "add", ".")
	run(t, root, "commit", "-q", "-m", "fixture base")
	t.Chdir(root)

	wantCheck := "skills index missing entry for skill 'alpha' (regenerate: bench skills-index --write)\n"
	report, code := Command(nil)
	if report != wantCheck || code != 1 {
		t.Fatalf("no-args (default --check) on drifted tree = %q, %d, want %q, 1", report, code, wantCheck)
	}

	report, code = Command([]string{"--check"})
	if report != wantCheck || code != 1 {
		t.Fatalf("--check on drifted tree = %q, %d, want %q, 1", report, code, wantCheck)
	}

	before, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	wantConflict := grammar.Help + " (--check and --write are mutually exclusive)\n"
	report, code = Command([]string{"--check", "--write"})
	if report != wantConflict || code != 2 {
		t.Fatalf("--check --write = %q, %d, want %q, 2", report, code, wantConflict)
	}
	after, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("--check --write changed bytes:\n%q\nwant\n%q", after, before)
	}

	report, code = Command([]string{"--write"})
	if report != "" || code != 0 {
		t.Fatalf("--write = %q, %d, want \"\", 0", report, code)
	}

	report, code = Command([]string{"--check"})
	if report != "" || code != 0 {
		t.Fatalf("--check after --write = %q, %d, want \"\", 0", report, code)
	}

	report, code = Command([]string{"--help"})
	if code != 0 || report != grammar.Help+"\n" {
		t.Fatalf("--help = %q, %d, want %q, 0", report, code, grammar.Help+"\n")
	}

	report, code = Command([]string{"--bogus"})
	if code != 2 {
		t.Fatalf("--bogus code = %d, want 2 (report %q)", code, report)
	}
}

// Argument validation is a verdict on the arguments alone: outside a repository the
// conflicting pair still earns usage and exit 2, not the not-in-repo refusal.
func TestConflictingModesAreRefusedOutsideARepository(t *testing.T) {
	t.Chdir(t.TempDir())

	wantConflict := grammar.Help + " (--check and --write are mutually exclusive)\n"
	report, code := Command([]string{"--check", "--write"})
	if report != wantConflict || code != 2 {
		t.Fatalf("--check --write outside a repo = %q, %d, want %q, 2", report, code, wantConflict)
	}
}

// run is the fixture git driver shared by this file's Command test.
func run(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// barrierMarker is the one byte sequence the interrupted-write handshake speaks. The
// child publishes it on an inherited pipe once the real replacement has created and
// written its sibling temp, so the parent learns the process is inside the vulnerable
// interval instead of guessing at it with a sleep or a directory scan.
const barrierMarker = "skills-index:pre-replacement\n"

// barrierPipeFD is the descriptor the marker arrives on: ExtraFiles[0] is fd 3 in the
// child, the first free descriptor after the standard three.
const barrierPipeFD = 3

// TestWriteBarrierHelperProcess is the child half of the handshake, inert unless the
// parent selects it. It installs the barrier the production replacement path reaches
// after the temp exists, publishes the marker there, and then blocks on the same
// cancellation context `bench skills-index --write` derives from SIGINT — so the signal
// lands exactly inside the interval where an unguarded replacement would leave residue.
func TestWriteBarrierHelperProcess(t *testing.T) {
	root := os.Getenv("BENCH_SKILLS_INDEX_BARRIER_ROOT")
	if root == "" {
		return
	}
	pipe := os.NewFile(barrierPipeFD, "barrier")
	preReplacementBarrier = func(ctx context.Context) {
		if _, err := pipe.WriteString(barrierMarker); err != nil {
			os.Exit(3)
		}
		<-ctx.Done()
	}
	if err := os.Chdir(root); err != nil {
		os.Exit(3)
	}
	_, code := Command([]string{"--write"})
	os.Exit(code)
}

// TestSIGINTBeforeReplacementLeavesNoResidueAndKeepsReferenceBytes is HI10: a fresh
// process interrupted between the temp's creation and the rename exits nonzero, leaves
// the reference authoritative, and takes its temp with it. Timing is a handshake, not a
// guess: the parent blocks reading the child's marker and signals only after it arrives.
func TestSIGINTBeforeReplacementLeavesNoResidueAndKeepsReferenceBytes(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	original := reference("")
	writeFile(t, root, ".bench/BENCH-reference.md", original)

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	cmd := exec.Command(os.Args[0], "-test.run=^TestWriteBarrierHelperProcess$")
	cmd.Env = append(os.Environ(), "BENCH_SKILLS_INDEX_BARRIER_ROOT="+root)
	cmd.ExtraFiles = []*os.File{writer}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	// The parent's copy of the write end must go, or a child that dies before the
	// barrier would leave this read blocked on a descriptor the parent holds open.
	writer.Close()
	defer cmd.Process.Kill()

	marker := make([]byte, len(barrierMarker))
	await(t, "the pre-replacement marker", func() error {
		_, err := io.ReadFull(reader, marker)
		return err
	})
	if string(marker) != barrierMarker {
		t.Fatalf("marker = %q, want %q", marker, barrierMarker)
	}

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	waitErr := await(t, "the interrupted child's exit", cmd.Wait)
	var exit *exec.ExitError
	if !errors.As(waitErr, &exit) || exit.ExitCode() == 0 {
		t.Fatalf("interrupted child exit = %v, want a nonzero exit status", waitErr)
	}

	kept, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(kept) != original {
		t.Fatalf("reference after an interrupted write =\n%q\nwant\n%q", kept, original)
	}
	if residue := siblingTemps(t, root); len(residue) != 0 {
		t.Fatalf("interrupted write left %v, want no .bench/.skills-index-* residue", residue)
	}
}
