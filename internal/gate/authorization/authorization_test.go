package authorization

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

func TestAuthorizeClassifiesAndValidatesOwnerEvidence(t *testing.T) {
	root := authorizationRepo(t)
	base := gitOutput(t, root, "write-tree")
	if got := Authorize(t.Context(), root, base); got.Kind != Inherited || Validate(root, base, got.Evidence) {
		t.Fatalf("unproved red = %+v, validate=%v; want inherited non-reusable", got, Validate(root, base, got.Evidence))
	}

	os.Remove(filepath.Join(root, "fail"))
	gitRun(t, root, "add", "-u")
	gitRun(t, root, "commit", "-q", "-m", "green base")
	branch := gitOutput(t, root, "branch", "--show-current")
	tip := gitOutput(t, root, "rev-parse", "HEAD")
	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap accepted a subject with no green evidence")
	}
	if out, err := exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", "refs/bench/green/"+branch).CombinedOutput(); err == nil {
		t.Fatalf("missing-evidence bootstrap created a marker: %s", out)
	}
	if got := gate.Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed inherited green = %+v", got)
	}
	if err := Bootstrap(root, branch, tip, ""); err != nil {
		t.Fatalf("bootstrap exact green: %v", err)
	}
	if err := Bootstrap(root, branch, tip, ""); err != nil {
		t.Fatalf("idempotent bootstrap: %v", err)
	}
	if got := gate.ValidateProjectGreen(root, branch); !got.ReusableGreen {
		t.Fatalf("bootstrapped project green = %+v", got)
	}
	gitRun(t, root, "checkout", "-q", "--detach")
	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap accepted detached HEAD")
	}
	gitRun(t, root, "checkout", "-q", branch)
	previous := gitOutput(t, root, "rev-parse", "HEAD^")
	gitRun(t, root, "update-ref", "refs/bench/green/"+branch, previous, tip)
	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap overwrote a conflicting marker")
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/"+branch); got != previous {
		t.Fatalf("conflicting marker = %s, want %s", got, previous)
	}
	if err := Bootstrap(root, branch, tip, previous); err != nil {
		t.Fatalf("advance expected descendant marker: %v", err)
	}
	if err := Bootstrap(root, branch, tip, previous); err != nil {
		t.Fatalf("replay advanced marker: %v", err)
	}
	os.WriteFile(filepath.Join(root, "later"), []byte("stale\n"), 0o644)
	gitRun(t, root, "add", "later")
	gitRun(t, root, "commit", "-q", "-m", "stale descendant")
	stale := gitOutput(t, root, "rev-parse", "HEAD")
	if err := Bootstrap(root, branch, stale, tip); err == nil {
		t.Fatal("bootstrap accepted a descendant without exact green evidence")
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/"+branch); got != tip {
		t.Fatalf("stale-evidence marker = %s, want %s", got, tip)
	}
	gitRun(t, root, "reset", "--hard", tip)
	greenTree := gitOutput(t, root, "write-tree")
	green := Authorize(t.Context(), root, greenTree)
	if green.Kind != Green || !Validate(root, greenTree, green.Evidence) {
		t.Fatalf("green authorization = %+v, validate=%v", green, Validate(root, greenTree, green.Evidence))
	}

	os.WriteFile(filepath.Join(root, "fail"), []byte("red\n"), 0o644)
	gitRun(t, root, "add", "fail")
	redTree := gitOutput(t, root, "write-tree")
	gitRun(t, root, "reset", "--hard", "HEAD")
	red := Authorize(t.Context(), root, redTree)
	if red.Kind != Candidate || Validate(root, redTree, red.Evidence) || Validate(root, greenTree, red.Evidence) {
		t.Fatalf("candidate red = %+v, same=%v other=%v", red, Validate(root, redTree, red.Evidence), Validate(root, greenTree, red.Evidence))
	}
}

func TestAuthorizeClassifiesMissingGateAsInfrastructure(t *testing.T) {
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	os.WriteFile(filepath.Join(root, "tracked"), []byte("x\n"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "base")
	tree := gitOutput(t, root, "write-tree")
	got := Authorize(t.Context(), root, tree)
	if got.Kind != Infrastructure || Validate(root, tree, got.Evidence) {
		t.Fatalf("missing gate = %+v, validate=%v", got, Validate(root, tree, got.Evidence))
	}
}

func TestBootstrapRefusesDivergentExpectedMarker(t *testing.T) {
	root, branch, previous, tip := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	divergent := siblingCommit(t, root, previous)
	gitRun(t, root, "update-ref", marker, divergent)

	// The marker matches the expectation, so the update-ref compare-and-swap alone
	// would happily fast-forward the marker off its own branch of history.
	if err := Bootstrap(root, branch, tip, divergent); err == nil {
		t.Fatal("bootstrap accepted an expected marker that is not an ancestor of the tip")
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != divergent {
		t.Fatalf("divergent marker = %s, want it left at %s", got, divergent)
	}
}

func TestBootstrapFastForwardsAncestorMarker(t *testing.T) {
	root, branch, previous, tip := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	gitRun(t, root, "update-ref", marker, previous)

	if err := Bootstrap(root, branch, tip, previous); err != nil {
		t.Fatalf("bootstrap over a strict ancestor marker: %v", err)
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != tip {
		t.Fatalf("advanced marker = %s, want the tip %s", got, tip)
	}
}

func TestAdvanceMarkerRecognizesLaggingMarker(t *testing.T) {
	root, branch, previous, lineage := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	gitRun(t, root, "update-ref", marker, previous)
	gitRun(t, root, "commit", "--allow-empty", "-qm", "destination")
	destination := gitOutput(t, root, "rev-parse", "HEAD")

	if err := AdvanceMarker(t.Context(), root, branch, destination, lineage); err != nil {
		t.Fatalf("advance from lagging marker: %v", err)
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != destination {
		t.Fatalf("marker = %s, want destination %s", got, destination)
	}
}

func TestAdvanceMarkerRefusesDivergentMarker(t *testing.T) {
	root, branch, parent, lineage := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	gitRun(t, root, "commit", "--allow-empty", "-qm", "destination")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	divergent := siblingCommit(t, root, parent)
	gitRun(t, root, "update-ref", marker, divergent)

	if err := AdvanceMarker(t.Context(), root, branch, destination, lineage); err == nil {
		t.Fatal("advance accepted a divergent marker")
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != divergent {
		t.Fatalf("marker = %s, want divergent marker %s", got, divergent)
	}
}

func TestAdvanceMarkerRefusesMarkerBetweenLineageAndDestination(t *testing.T) {
	root, branch, _, lineage := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	gitRun(t, root, "commit", "--allow-empty", "-qm", "between")
	between := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, root, "commit", "--allow-empty", "-qm", "destination")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, root, "update-ref", marker, between)

	if err := AdvanceMarker(t.Context(), root, branch, destination, lineage); err == nil {
		t.Fatal("advance accepted a marker between lineage and destination")
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != between {
		t.Fatalf("marker = %s, want between marker %s", got, between)
	}
}

func TestAdvanceMarkerRefusesUndecidableMarkerAncestry(t *testing.T) {
	root, branch, actual, _ := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	tree := gitOutput(t, root, "rev-parse", actual+"^{tree}")
	intermediate := gitOutput(t, root, "commit-tree", tree, "-p", actual, "-m", "missing intermediate")
	lineage := gitOutput(t, root, "commit-tree", tree, "-p", intermediate, "-m", "lineage")
	destination := gitOutput(t, root, "commit-tree", tree, "-p", lineage, "-m", "destination")
	gitRun(t, root, "update-ref", marker, actual)
	if err := os.Remove(filepath.Join(root, ".git", "objects", intermediate[:2], intermediate[2:])); err != nil {
		t.Fatalf("unlink intermediate object: %v", err)
	}

	if err := AdvanceMarker(t.Context(), root, branch, destination, lineage); err == nil {
		t.Fatal("advance accepted an unreadable marker ancestry")
	}
	if got := markerFile(t, root, branch); got != actual {
		t.Fatalf("marker = %s, want marker %s", got, actual)
	}
}

func TestBootstrapIsNoOpWhenMarkerEqualsTip(t *testing.T) {
	root, branch, previous, tip := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	gitRun(t, root, "update-ref", marker, tip)

	// A marker already at the tip is the answer whatever the caller expected, so the
	// ancestor rule never gets to reject an expectation there is nothing to write over.
	if err := Bootstrap(root, branch, tip, siblingCommit(t, root, previous)); err != nil {
		t.Fatalf("bootstrap with the marker already at the tip: %v", err)
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != tip {
		t.Fatalf("marker = %s, want the tip %s untouched", got, tip)
	}
}

func TestBootstrapFailsClosedOnUnreadableMarker(t *testing.T) {
	root, branch, previous, tip := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	unreadable := siblingCommit(t, root, previous)
	gitRun(t, root, "update-ref", marker, unreadable)
	if err := os.Remove(filepath.Join(root, ".git", "objects", unreadable[:2], unreadable[2:])); err != nil {
		t.Fatalf("unlink marker object: %v", err)
	}

	if err := Bootstrap(root, branch, tip, previous); err == nil {
		t.Fatal("bootstrap read an unpeelable marker as an absent one")
	}
	if got := markerFile(t, root, branch); got != unreadable {
		t.Fatalf("marker = %s, want the unreadable %s rather than a fresh one", got, unreadable)
	}
}

func TestBootstrapRefusesDanglingSymbolicMarker(t *testing.T) {
	root, branch, _, tip := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	missing := "refs/bench/missing-marker"
	gitRun(t, root, "symbolic-ref", marker, missing)
	before := gitOutput(t, root, "symbolic-ref", marker)

	if err := Bootstrap(root, branch, tip, ""); err == nil {
		t.Fatal("bootstrap treated a dangling symbolic marker as absent")
	}
	if got := gitOutput(t, root, "symbolic-ref", marker); got != before {
		t.Fatalf("marker = %q, want dangling symbolic marker %q left untouched", got, before)
	}
	if out, err := exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", missing).CombinedOutput(); err == nil {
		t.Fatalf("bootstrap created the dangling marker target: %s", out)
	}
}

func TestBootstrapExpectationIsAFullObjectID(t *testing.T) {
	root, branch, previous, tip := greenBootstrapRepo(t)
	marker := "refs/bench/green/" + branch
	plant := func() { gitRun(t, root, "update-ref", marker, previous) }

	plant()
	if err := Bootstrap(root, branch, tip, previous); err != nil {
		t.Fatalf("bootstrap with the marker's own object id: %v", err)
	}
	plant()
	if err := Bootstrap(root, branch, tip, marker); err == nil {
		t.Fatal("bootstrap resolved a ref name as the expectation")
	}
	plant()
	if err := Bootstrap(root, branch, tip, previous[:7]); err == nil {
		t.Fatal("bootstrap resolved an abbreviated object id as the expectation")
	}
	if got := gitOutput(t, root, "rev-parse", marker); got != previous {
		t.Fatalf("marker = %s, want it left at %s by the refused expectations", got, previous)
	}
}

// greenBootstrapRepo returns a repository whose tip carries reusable green evidence,
// with previous naming the tip's parent.
func greenBootstrapRepo(t *testing.T) (root, branch, previous, tip string) {
	t.Helper()
	root = authorizationRepo(t)
	os.Remove(filepath.Join(root, "fail"))
	gitRun(t, root, "add", "-u")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "green base")
	branch = gitOutput(t, root, "branch", "--show-current")
	tip = gitOutput(t, root, "rev-parse", "HEAD")
	previous = gitOutput(t, root, "rev-parse", "HEAD^")
	if got := gate.Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("seed green evidence = %+v", got)
	}
	return root, branch, previous, tip
}

// siblingCommit writes a commit beside parent that no branch reaches, leaving the
// working tree and HEAD where the caller's green evidence needs them.
func siblingCommit(t *testing.T, root, parent string) string {
	t.Helper()
	return gitOutput(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench",
		"commit-tree", parent+"^{tree}", "-p", parent, "-m", "divergent")
}

func markerFile(t *testing.T, root, branch string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, ".git", "refs", "bench", "green", branch))
	if err != nil {
		t.Fatalf("read marker ref: %v", err)
	}
	return string(bytesTrimSpace(raw))
}

func authorizationRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	os.MkdirAll(filepath.Join(root, ".bench"), 0o755)
	os.WriteFile(filepath.Join(root, ".bench", "gate.sh"), []byte("#!/bin/sh\ntest ! -f fail\n"), 0o755)
	os.WriteFile(filepath.Join(root, ".bench", "gate-inputs.json"), []byte(`{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`), 0o644)
	os.WriteFile(filepath.Join(root, "fail"), []byte("red\n"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "red")
	return root
}

func gitRun(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

func gitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return string(bytesTrimSpace(out))
}

func bytesTrimSpace(value []byte) []byte {
	start, end := 0, len(value)
	for start < end && (value[start] == ' ' || value[start] == '\n' || value[start] == '\r' || value[start] == '\t') {
		start++
	}
	for end > start && (value[end-1] == ' ' || value[end-1] == '\n' || value[end-1] == '\r' || value[end-1] == '\t') {
		end--
	}
	return value[start:end]
}
