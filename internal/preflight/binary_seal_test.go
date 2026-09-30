package preflight

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestBinarySealRedsAMismatchedDestination covers BF26. A build whose
// published dist/bench fails its seal reds binary-seal, carries the
// verifier's refusal whole, and makes the whole verdict red, so a --full run
// names the rebuild before it opens the transaction.
func TestBinarySealRedsAMismatchedDestination(t *testing.T) {
	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent = true
	refusal := &freshness.Refusal{Executable: "/root/dist/bench", RepairRoot: "/root", Cause: errors.New("source digest mismatch")}
	f.BinarySealRefusal = refusal

	v := Decide(f)
	c, ok := checkRow(v, "binary-seal")
	if !ok {
		t.Fatalf("Checks = %#v, want a binary-seal row", v.Checks)
	}
	if c.Verdict != verdictRed || c.Detail != refusal.Error() {
		t.Errorf("binary-seal row = %+v, want the red detailed %q", c, refusal)
	}
	if !v.Red {
		t.Errorf("Verdict.Red = false, want true when the destination seal mismatches")
	}
	if at := rowIndex(v, "binary-seal"); at != rowIndex(v, "kit-pin")+1 {
		t.Errorf("binary-seal sits at index %d, want directly after kit-pin at %d", at, rowIndex(v, "kit-pin"))
	}
}

// TestBinarySealNotApplicableWithoutADestinationBinary covers BF27. A root
// that publishes no dist/bench reports the row not applicable, keeps the
// verdict green, and leaves every other row exactly as it was.
func TestBinarySealNotApplicableWithoutADestinationBinary(t *testing.T) {
	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent = false

	v := Decide(f)
	c, ok := checkRow(v, "binary-seal")
	if !ok {
		t.Fatalf("Checks = %#v, want a binary-seal row", v.Checks)
	}
	if c.Verdict != verdictNA || c.Detail != "" || c.Next != "" {
		t.Errorf("binary-seal row = %+v, want not-applicable with no detail and no remedy", c)
	}
	if v.Red {
		t.Errorf("Verdict.Red = true, want false when no destination binary is published")
	}
	for _, other := range v.Checks {
		if other.Check == "binary-seal" || other.Check == "diff-nonempty" {
			continue
		}
		if other.Verdict != verdictGreen {
			t.Errorf("%s row = %+v, want green and unchanged by an absent binary", other.Check, other)
		}
	}
}

// TestBinarySealIsBuildModeOnly pins the mode gate. Story 26 names
// `bench preflight build`, so a review preflight renders no binary-seal row
// at all, even over a root whose published binary fails its seal.
func TestBinarySealIsBuildModeOnly(t *testing.T) {
	f := baseFacts()
	f.BinarySealPresent = true
	f.BinarySealRefusal = &freshness.Refusal{Executable: "/root/dist/bench", RepairRoot: "/root", Cause: errors.New("source digest mismatch")}

	v := Decide(f)
	if _, ok := checkRow(v, "binary-seal"); ok {
		t.Errorf("review Checks = %#v, want no binary-seal row", v.Checks)
	}
	if v.Red {
		t.Errorf("Verdict.Red = true, want false: a review preflight does not grade the binary")
	}
}

// sealFixtureRoot is the smallest rebuildable root the seal primitives accept: a module
// whose ./cmd/bench resolves, plus the build-inputs manifest the source digest reads.
func sealFixtureRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range map[string]string{
		"go.mod":                      "module example.com/sealfixture\n\ngo 1.25\n",
		"cmd/bench/main.go":           "package main\n\nfunc main() {}\n",
		"scripts/go-build.sh":         "#!/usr/bin/env bash\n",
		freshness.BuildInputsManifest: freshness.BuildInputLine("build_script", "scripts/go-build.sh"),
	} {
		preflighttest.MustWriteFile(t, filepath.Join(root, filepath.FromSlash(name)), body)
	}
	return root
}

// publishSealedBinary seals a dist/bench pair for root through the same
// publisher a build runs, so the fixture's green case is a real seal rather
// than a hand-written sidecar.
func publishSealedBinary(t *testing.T, root string) string {
	t.Helper()
	staged := filepath.Join(root, "staged-bench")
	preflighttest.MustWriteFile(t, staged, "Bench executable")
	if err := os.Chmod(staged, 0o755); err != nil {
		t.Fatal(err)
	}
	executable := freshness.PublishedExecutable(root)
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "1.2.3"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	return executable
}

// TestBinarySealFactsGradeARootOnDisk covers BF26 and BF27 at the gatherer
// seam. Both row tests set the facts by hand, so the path spelling and the
// verifier call answer to nothing. This drives binarySealFacts against real
// trees instead. A dangling dist/bench is the state the two derivations
// disagreed on: doctor lstats and reds it, so the build preflight grades it
// too rather than reading a broken link as no binary at all.
func TestBinarySealFactsGradeARootOnDisk(t *testing.T) {
	for _, tc := range []struct {
		name    string
		place   func(*testing.T, string)
		present bool
		refused bool
	}{
		{
			name:  "no binary published",
			place: func(*testing.T, string) {},
		},
		{
			name:    "sealed pair",
			place:   func(t *testing.T, root string) { publishSealedBinary(t, root) },
			present: true,
		},
		{
			name: "sources changed after the seal",
			place: func(t *testing.T, root string) {
				publishSealedBinary(t, root)
				preflighttest.MustWriteFile(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() { _ = 1 }\n")
			},
			present: true,
			refused: true,
		},
		{
			name: "dangling symbolic link",
			place: func(t *testing.T, root string) {
				executable := freshness.PublishedExecutable(root)
				if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(root, "absent-bench"), executable); err != nil {
					t.Fatal(err)
				}
			},
			present: true,
			refused: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := sealFixtureRoot(t)
			tc.place(t, root)

			present, refusal := binarySealFacts(root)
			if present != tc.present {
				t.Errorf("binarySealFacts present = %v, want %v (refusal %v)", present, tc.present, refusal)
			}
			if refused := refusal != nil; refused != tc.refused {
				t.Errorf("binarySealFacts refusal = %v, want refused = %v", refusal, tc.refused)
			}
		})
	}
}

// TestDanglingDestinationBinaryRedsBinarySealInBuildMode carries the
// gathered facts of a broken dist/bench link through to the verdict. The
// row an operator reads is the whole point: a --full build that trusts a
// link to nothing runs an executable nobody can name.
func TestDanglingDestinationBinaryRedsBinarySealInBuildMode(t *testing.T) {
	root := sealFixtureRoot(t)
	executable := freshness.PublishedExecutable(root)
	if err := os.MkdirAll(filepath.Dir(executable), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "absent-bench"), executable); err != nil {
		t.Fatal(err)
	}

	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent, f.BinarySealRefusal = binarySealFacts(root)

	v := Decide(f)
	c, ok := checkRow(v, "binary-seal")
	if !ok {
		t.Fatalf("Checks = %#v, want a binary-seal row", v.Checks)
	}
	if c.Verdict != verdictRed {
		t.Errorf("binary-seal row = %+v, want red over a dangling dist/bench", c)
	}
	if !v.Red {
		t.Errorf("Verdict.Red = false, want true when dist/bench links to nothing")
	}
}

// TestBinarySealRemedyNamesTheWorktreeBuildVerb grades the remedy a stale binary
// carries. An assignment worktree sits under the pool path, and the hooks refuse a cd
// into the pool path, so the owned row names `bench worktree build` and no cell names
// the verifier's cd route. A root no assignment owns keeps the verifier's refusal whole.
func TestBinarySealRemedyNamesTheWorktreeBuildVerb(t *testing.T) {
	root := sealFixtureRoot(t)
	publishSealedBinary(t, root)
	preflighttest.MustWriteFile(t, filepath.Join(root, "cmd", "bench", "main.go"), "package main\n\nfunc main() { _ = 1 }\n")

	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	f.BinarySealPresent, f.BinarySealRefusal = binarySealFacts(root)

	f.AssignmentTarget = "a1b2c3"
	owned, _ := checkRow(Decide(f), "binary-seal")
	if want := usage.WorktreeBuildFor(f.AssignmentTarget); owned.Verdict != verdictRed || owned.Next != want {
		t.Errorf("owned binary-seal row = %+v, want red with next %q", owned, want)
	}
	// The expectation above reads the fill that the row uses, so this check is the one that
	// sees a fill that leaves the slot open.
	if !strings.HasSuffix(owned.Next, " "+f.AssignmentTarget) {
		t.Errorf("owned binary-seal next %q does not name the target %q", owned.Next, f.AssignmentTarget)
	}
	for _, cell := range []string{owned.Detail, owned.Next} {
		if strings.Contains(cell, "cd ") || strings.Contains(cell, "scripts/go-build.sh") {
			t.Errorf("owned binary-seal cell %q names the cd route the hooks refuse in the pool path", cell)
		}
	}

	f.AssignmentTarget = ""
	unowned, _ := checkRow(Decide(f), "binary-seal")
	if unowned.Verdict != verdictRed || !strings.Contains(unowned.Detail, freshness.RebuildAction(root)) || unowned.Next != "" {
		t.Errorf("unowned binary-seal row = %+v, want red carrying the verifier's rebuild action %q", unowned, freshness.RebuildAction(root))
	}
}

// rowIndex is the position of one named row, or -1. Order is part of the
// verdict's contract, so a test that grades placement reads it here.
func rowIndex(v Verdict, name string) int {
	for i, c := range v.Checks {
		if c.Check == name {
			return i
		}
	}
	return -1
}
