package handoff

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/status"
)

// HS13, story 16. A State that pins a real commit off the tip's ancestry is a stale resume
// target: the session it sends there builds on a history this section no longer holds. The
// run refuses, prints the offending line, and leaves the document byte for byte.
func TestCommandRefusesAStatePinningAnOffAncestryCommit(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	off := offAncestryCommit(t, root)
	before := seedState(t, document, "The build resumes from `"+off+"`.")

	out, code := runAt(t, root, nil)
	if code != 1 {
		t.Fatalf("handoff over an off-ancestry pin = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, "The build resumes from `"+off+"`.") {
		t.Errorf("the refusal does not print the offending line\n%s", out)
	}
	// The reason, not the exit code alone. A commit the tip has lost and an object that
	// was never a commit need different repairs, and the printed reason is what separates
	// them for the reader.
	if !strings.Contains(out, faultOffAncestry) {
		t.Errorf("the refusal does not give the ancestry reason %q\n%s", faultOffAncestry, out)
	}
	if strings.Contains(out, faultNotACommit) {
		t.Errorf("a real commit drew the not-a-commit reason\n%s", out)
	}
	if after := read(t, document); after != before {
		t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// HS14, story 17. State is prose. Two English words that happen to spell hex are not pins,
// and a scan that judged shape alone would refuse a sentence.
func TestCommandAcceptsHexShapedEnglishInState(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	seedState(t, document, "The `facade` holds while the `decade` of debt is paid.")

	runIn(t, root, nil)
}

// HS15, story 18. A tree hash exists in the object store and is still not a commit, so it
// names no resume target. The `^{commit}` peel is what tells the two apart.
//
// The reason is asserted, not the exit code alone. This tree is also off the tip's
// ancestry, so a scan that dropped the peel would still exit 1 — by the wrong route, under
// the wrong reason. The reason text is what makes the peel load-bearing here.
func TestCommandRefusesAStateNamingATreeHash(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	tree := gitOut(t, root, "rev-parse", "HEAD^{tree}")
	seedState(t, document, "The tree is `"+tree+"`.")

	out, code := runAt(t, root, nil)
	if code != 1 {
		t.Fatalf("handoff over a tree hash = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, faultNotACommit) {
		t.Errorf("the refusal does not give the not-a-commit reason %q\n%s", faultNotACommit, out)
	}
	if strings.Contains(out, faultOffAncestry) {
		t.Errorf("a tree hash drew the ancestry reason, so the commit peel never ran\n%s", out)
	}
}

// The scan's length boundary, both sides. Seven is Git's shortest unambiguous
// abbreviation, so a seven-character token is a pin and a six-character one is a word.
func TestStateScanLengthBoundary(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	off := offAncestryCommit(t, root)
	for _, tc := range []struct {
		name  string
		token string
		want  int
	}{
		{"seven characters are scanned", off[:7], 1},
		{"six characters are not", off[:6], 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seedState(t, document, "Resume at `"+tc.token+"`.")
			if out, code := runAt(t, root, nil); code != tc.want {
				t.Fatalf("handoff over %q = (%q, %d), want exit %d", tc.token, out, code, tc.want)
			}
		})
	}
}

// offAncestryCommit lands one commit in a linked worktree on its own branch and returns
// its full sha. The object store is shared, so the primary checkout can resolve the commit
// while main's tip does not contain it.
func offAncestryCommit(t *testing.T, root string) string {
	t.Helper()
	a := activeAssignment(t, root, "hs-off-ancestry", "hs-off-ancestry")
	commitIn(t, a.Worktree, "off the ancestry")
	return gitOut(t, a.Worktree, "rev-parse", "HEAD")
}

// seedSection plants main's reviewer-owned fields and returns the document's bytes, so a
// refusal test can compare against exactly what the refused run read. It goes in through
// the leaf package's own writer, so the planted bytes are the ones a real run would parse.
func seedSection(t *testing.T, document, next, state string) string {
	t.Helper()
	if err := handoffdoc.WriteSection(document, handoffdoc.Section{Key: handoffdoc.MainKey, Next: next, State: state}); err != nil {
		t.Fatalf("seed section: %v", err)
	}
	return read(t, document)
}

func seedState(t *testing.T, document, state string) string {
	t.Helper()
	return seedSection(t, document, "", state)
}
