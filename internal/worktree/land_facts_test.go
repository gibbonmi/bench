package worktree

import (
	"bytes"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/worktree/landingpolicy"
)

// The tests in this file are the FA1 focused fact-adapter proofs: one real-Git
// test per named landing fact group — source, destination, resume,
// publication, release, and residue — each asserting the exact typed facts
// the parent boundary hands the pure landing policy.

// TestLandingResidueFactAdapterTranslatesRealDestination is the residue fact
// group. A real checkout with one staged, one unstaged, one untracked, and one
// ignored path must translate into exactly those typed status entries, and the
// lazy staged-content fact must answer from the real index.
func TestLandingResidueFactAdapterTranslatesRealDestination(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	mustWrite(t, filepath.Join(root, ".gitignore"), []byte("junk/\n"), 0o644)
	gitRun(t, root, "add", ".gitignore")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare ignored dir")
	destination := gitOutput(t, root, "rev-parse", "HEAD")

	mustWrite(t, filepath.Join(root, "tracked.txt"), []byte("staged edit\n"), 0o644)
	gitRun(t, root, "add", "tracked.txt")
	mustWrite(t, filepath.Join(root, "README.md"), []byte("unstaged edit\n"), 0o644)
	mustWrite(t, filepath.Join(root, "stray.txt"), []byte("stray\n"), 0o600)
	mustMkdirAll(t, filepath.Join(root, "junk"), 0o755)
	mustWrite(t, filepath.Join(root, "junk", "out"), []byte("ignored\n"), 0o600)

	facts := destinationResidueFacts(root, destination, destination, destination+"^")
	if !facts.NestedClean || !facts.StatusReadable || !facts.StatusWellFormed {
		t.Fatalf("healthy destination facts = %+v", facts)
	}
	want := []landingpolicy.StatusEntry{
		{Status: " M", Path: "README.md"},
		{Status: "M ", Path: "tracked.txt"},
		{Status: "??", Path: "stray.txt"},
		{Status: "!!", Path: "junk/out"},
	}
	if !reflect.DeepEqual(facts.Entries, want) {
		t.Fatalf("residue entries = %+v, want %+v", facts.Entries, want)
	}
	if !facts.DestinationAtPublished {
		t.Fatalf("destination-at-published fact = false for identical commits")
	}
	if facts.StagedMatchesPublished() {
		t.Fatalf("staged-matches-published fact = true for a staged edit beyond the base")
	}
	if facts.IgnoredDeclared() {
		t.Fatalf("ignored-declared fact = true with no build-output declaration")
	}
	markProof(t, "landing/adapter/facts")
}

// TestLandingMarkerFactAdapterTranslatesGreenMarker is the resume fact group.
// The adapter must read the real ancestry between the published landing and
// the recorded green marker.
func TestLandingMarkerFactAdapterTranslatesGreenMarker(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	published := gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, root, "later.txt", "later\n", "destination moves on")
	destination := gitOutput(t, root, "rev-parse", "HEAD")

	same := resumeMarkerFacts(root, published, published, "")
	if !same.DestinationAtPublished {
		t.Fatalf("unmoved destination facts = %+v, want DestinationAtPublished", same)
	}
	moved := resumeMarkerFacts(root, destination, published, destination)
	if moved.DestinationAtPublished || !moved.MarkerPresent || !moved.MarkerReachesPublished() {
		t.Fatalf("covering marker facts = %+v, want present and reaching", moved)
	}
	absent := resumeMarkerFacts(root, destination, published, "")
	if absent.MarkerPresent {
		t.Fatalf("absent marker facts = %+v, want MarkerPresent=false", absent)
	}
	divergent := resumeMarkerFacts(root, destination, destination, published)
	if !divergent.MarkerPresent || divergent.MarkerReachesPublished() {
		t.Fatalf("behind marker facts = %+v, want present and not reaching", divergent)
	}
}

// TestLandingPublicationFactAdapterTranslatesMergeCommit is the publication
// fact group. A real two-parent landing merge must translate into the exact
// typed identity, reachability, parent, and range facts.
func TestLandingPublicationFactAdapterTranslatesMergeCommit(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	base := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, root, "checkout", "-q", "-b", "source")
	commitInWorktree(t, root, "owned.txt", "reviewed\n", "reviewed source")
	tip := gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, root, "checkout", "-q", "main")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "merge", "-q", "--no-ff", "-m", "landing", "source")
	published := gitOutput(t, root, "rev-parse", "HEAD")

	facts, sourceBase, destinationBase := publicationFacts(root, published, published, base, tip, "")
	want := landingpolicy.PublicationFacts{
		Requested: published, Resolved: true, Published: published,
		ReachableFromDestination: true, ParentsReadable: true,
		Parents: []string{published, base, tip}, RequestedSource: tip,
		RangeAuthenticates: true,
	}
	if !reflect.DeepEqual(facts, want) {
		t.Fatalf("publication facts = %+v, want %+v", facts, want)
	}
	if sourceBase != base || destinationBase != base {
		t.Fatalf("bases = (%q, %q), want both %q", sourceBase, destinationBase, base)
	}

	foreign, _, _ := publicationFacts(root, published, base, base, tip, "")
	if !foreign.Resolved || foreign.Published != base || foreign.ReachableFromDestination && len(foreign.Parents) == 3 {
		t.Fatalf("non-merge publication facts = %+v, want a first-parent record", foreign)
	}
	unresolved, _, _ := publicationFacts(root, published, strings.Repeat("f", 40), base, tip, "")
	if unresolved.Resolved {
		t.Fatalf("unresolved value facts = %+v, want Resolved=false", unresolved)
	}
}

// TestLandingDestinationFactAdapterTranslatesCleanCheckout is the destination
// fact group. A clean default-branch checkout must translate into the exact
// tip, branch, marker, and fingerprint facts.
func TestLandingDestinationFactAdapterTranslatesCleanCheckout(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	head := gitOutput(t, root, "rev-parse", "HEAD")
	tip, branch, marker, fingerprint, err := landingDestination(root)
	if err != nil || tip != head || branch != "main" || marker != "" {
		t.Fatalf("destination facts = (%q, %q, %q, %v), want (%q, main, empty marker, nil)", tip, branch, marker, err, head)
	}
	wantFingerprint, err := landing.CheckoutFingerprint(root)
	if err != nil || fingerprint != wantFingerprint {
		t.Fatalf("destination fingerprint = %q, %v, want %q", fingerprint, err, wantFingerprint)
	}
}

// TestLandingSourceFactAdapterTranslatesRealWorktree is the source fact group.
// A real reviewed assignment worktree must translate into the exact typed
// base, tip, and fingerprint facts, with no spec named.
func TestLandingSourceFactAdapterTranslatesRealWorktree(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	base := gitOutput(t, root, "rev-parse", "HEAD")
	creation := mustCreate(t, root, home, "source-fact-adapter", "source facts")
	commitInWorktree(t, creation.Path, "owned.txt", "reviewed bytes\n", "reviewed source")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	fact, err := landingSource(root, creation.Assignment, base, tip, "")
	if err != nil {
		t.Fatalf("landingSource: %v", err)
	}
	wantFingerprint, err := landing.CheckoutFingerprint(creation.Path)
	if err != nil {
		t.Fatal(err)
	}
	want := landingSourceFact{base: base, tip: tip, fingerprint: wantFingerprint}
	if !reflect.DeepEqual(fact, want) {
		t.Fatalf("source facts = %+v, want %+v", fact, want)
	}
}

// TestLandingReleaseFactAdapterTranslatesReleaseExit is the release fact
// group. A real published landing whose release step fails must translate the
// nonzero release exit into the terminal release-interruption outcome, and the
// policy's terminal verdict must own the rendered state and exit code.
func TestLandingReleaseFactAdapterTranslatesReleaseExit(t *testing.T) {
	request := "release-fact-adapter"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	original := releaseLandingAssignment
	releaseLandingAssignment = func(string, string, []string, io.Writer, io.Writer) int { return 1 }
	t.Cleanup(func() { releaseLandingAssignment = original })
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	want := landingpolicy.Terminal(landingpolicy.TerminalFacts{FailedStep: "release", Active: true})
	if code != want.ExitCode || !strings.Contains(stdout.String(), "worktree="+want.WorktreeState+",next=") {
		t.Fatalf("release-failure landing = (%d, %q), want exit %d and state %q", code, stdout.String(), want.ExitCode, want.WorktreeState)
	}
}
