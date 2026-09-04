package handoff

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os/exec"
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

// HS27, story 28. A State that opens a fence and never closes it swallows every section
// below it, so a later run appends a duplicate of the section it could not see. The run
// refuses on the line that opened the fence and leaves the document alone.
func TestCommandRefusesAStateWithAnUnterminatedFence(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	before := seedState(t, document, "Resume here:\n\n```console\n$ bench gate")

	out, code := runAt(t, root, nil)
	if code != 1 {
		t.Fatalf("handoff over an unterminated fence = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, "```console") {
		t.Errorf("the refusal does not print the line that opened the fence\n%s", out)
	}
	assertNamesLine(t, out, document, before, "```console")
	if after := read(t, document); after != before {
		t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// HS29, story 29. A level-two heading in State opens a section this grammar has no key
// for, and every later verb then refuses the file. The run refuses while the writer still
// has the text, and prints the heading they wrote.
func TestCommandRefusesAStateHeadingOutsideAFence(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	before := seedState(t, document, "The build is live.\n\n## Open questions\n\nNone.")

	out, code := runAt(t, root, nil)
	if code != 1 {
		t.Fatalf("handoff over an unfenced heading = (%q, %d), want exit 1", out, code)
	}
	if !strings.Contains(out, "## Open questions") {
		t.Errorf("the refusal does not print the offending line\n%s", out)
	}
	assertNamesLine(t, out, document, before, "## Open questions")
	if after := read(t, document); after != before {
		t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// HS31, story 31. A session quotes an earlier handoff line inside a fence. The sha in that
// block reports a run that already happened; it pins nothing here, so a real off-ancestry
// commit inside a fence is not a stale resume target.
//
// The sha is backticked inside the fence, which is what makes the row bite: the token rule
// is a backticked hex run, so an unquoted sha in the block would pass a scan that skipped no
// fence at all.
func TestCommandAcceptsAnOffAncestryCommitInsideAFence(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	off := offAncestryCommit(t, root)
	seedState(t, document, "The earlier close wrote:\n\n```\nThe build resumes from `"+off+"`.\n```\n\nThe tree is clean.")

	runIn(t, root, nil)
}

// HS32, story 32. An abbreviation that expands to two objects resolves to neither, so both
// `cat-file` probes fail and the token reads as prose under the exit code alone. It is a
// pin that sends a cold session nowhere, and it gets its own reason.
func TestCommandRefusesAnAmbiguousAbbreviationInState(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	prefix := ambiguousPrefix(t, root)
	seedState(t, document, "Resume from `"+prefix+"`.")

	out, code := runAt(t, root, nil)
	if code != 1 {
		t.Fatalf("handoff over %q = (%q, %d), want exit 1", prefix, out, code)
	}
	if !strings.Contains(out, faultAmbiguous) {
		t.Errorf("the refusal does not give the ambiguity reason %q\n%s", faultAmbiguous, out)
	}
}

// ambiguousPrefix writes two blobs whose object ids share their first seven hex digits and
// returns that prefix. Git then expands the prefix to two objects, which is the state the
// scan has to tell apart from prose.
//
// The search runs over Git's own blob hash in this process, so only the two colliding
// blobs are ever written. Hashing every candidate through `git hash-object` would take
// tens of thousands of subprocesses to reach the same pair. The candidate texts are a
// counted sequence, so every run finds the same collision.
func ambiguousPrefix(t *testing.T, root string) string {
	t.Helper()
	seen := map[string]string{}
	for i := range 5_000_000 {
		content := fmt.Sprintf("hs32 collision probe %d\n", i)
		sum := sha1.Sum([]byte(fmt.Sprintf("blob %d\x00%s", len(content), content)))
		prefix := hex.EncodeToString(sum[:])[:7]
		if other, found := seen[prefix]; found {
			writeBlob(t, root, other)
			writeBlob(t, root, content)
			return prefix
		}
		seen[prefix] = content
	}
	t.Fatal("no two probe blobs shared a seven-digit object id prefix")
	return ""
}

// writeBlob stores one blob in the repository's object store.
func writeBlob(t *testing.T, root, content string) {
	t.Helper()
	cmd := exec.Command("git", "-C", root, "hash-object", "-w", "--stdin")
	cmd.Stdin = strings.NewReader(content)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("hash-object: %v\n%s", err, out)
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

// assertNamesLine checks that a parser refusal points the reader at the document and the
// line the offending text sits on. The expected line is read out of the seeded bytes rather
// than counted by hand, so a change to the rendered shape moves the expectation with it.
func assertNamesLine(t *testing.T, out, document, content, offending string) {
	t.Helper()
	line := 0
	for i, text := range strings.Split(content, "\n") {
		if strings.TrimRight(text, " \t") == offending {
			line = i + 1
			break
		}
	}
	if line == 0 {
		t.Fatalf("the seeded document holds no line %q\n%s", offending, content)
	}
	if want := fmt.Sprintf("%s:%d", document, line); !strings.Contains(out, want) {
		t.Errorf("the refusal does not name %s\n%s", want, out)
	}
}

func seedState(t *testing.T, document, state string) string {
	t.Helper()
	return seedSection(t, document, "", state)
}
