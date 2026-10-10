package roadmap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/usage"
)

func journalPath(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(root, filepath.FromSlash(learnings.JournalPath))
}

var learningArgs = []string{"the", "gate", "hid", "a", "rule", "--what", "it failed twice", "--right", "read the map first"}

// newProseGradedRepo is newRepo plus an empty exclusion policy, so the fixture's prose
// grade runs against real content instead of the absent-policy diagnostic.
func newProseGradedRepo(t *testing.T) string {
	t.Helper()
	root := newRepo(t)
	full := filepath.Join(root, filepath.FromSlash(prose.ExclusionFile))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// words returns an n-word sentence, plain tokens with no punctuation to count, closed
// with a period so the grader reads it as one complete sentence.
func words(n int) string {
	list := make([]string, n)
	for i := range list {
		list[i] = fmt.Sprintf("w%d", i+1)
	}
	return strings.Join(list, " ") + "."
}

// TestLearningRefusesAnOverBoundSentenceAndWritesNothing covers story 9's writer half: a
// --what sentence over the 25-word bound must not reach the journal the live-tree prose
// sweep will grade later. The 26th word is what turns this red; the 25-word twin proves
// the bound is honored at the boundary, not just past it. The composed bullet reads "-
// **What happened:** <words>.": the list marker is stripped before grading, but the
// two-word bold label joins the sentence the grader counts, so passing wantTotal-2 words
// closes the sentence at exactly wantTotal words.
func TestLearningRefusesAnOverBoundSentenceAndWritesNothing(t *testing.T) {
	const label = 2
	for _, tc := range []struct {
		name       string
		wantTotal  int
		wantCode   int
		wantWrites bool
	}{
		{"25 words writes", 25, 0, true},
		{"26 words refuses", 26, 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newProseGradedRepo(t)
			args := append([]string{"a", "title"}, "--what", words(tc.wantTotal-label), "--right", "read the map first")
			out, code := LearningCommand(args)
			if code != tc.wantCode {
				t.Fatalf("code = %d, want %d\n%s", code, tc.wantCode, out)
			}
			_, statErr := os.Stat(journalPath(t, root))
			if wrote := statErr == nil; wrote != tc.wantWrites {
				t.Fatalf("journal exists = %v, want %v", wrote, tc.wantWrites)
			}
			if tc.wantWrites {
				return
			}
			// A prose refusal is a content error: exit 1 above, and no usage line here.
			if strings.Contains(out, learningGrammar.Help) {
				t.Fatalf("out = %q, want no usage line on a prose refusal", out)
			}
			if !strings.HasPrefix(out, "refused --what: ") {
				t.Fatalf("out = %q, want it to lead with the refused --what flag", out)
			}
			if !strings.Contains(out, fmt.Sprintf("sentence of %d words is over the 25-word bound", tc.wantTotal)) {
				t.Fatalf("out = %q, want the grader's own sentence diagnostic", out)
			}
		})
	}
}

// TestLearningRefusalNamesTheOffendingEntryLine covers story 9's second half: the
// refusal names which composed line is over bound, the --right bullet here, so the
// author can shorten the right field without opening the file to find it.
func TestLearningRefusalNamesTheOffendingEntryLine(t *testing.T) {
	root := newProseGradedRepo(t)
	out, code := LearningCommand([]string{"a", "title", "--what", "it happened", "--right", words(26)})
	if code != 1 {
		t.Fatalf("code = %d, want 1\n%s", code, out)
	}
	if !strings.Contains(out, "line 3") || !strings.HasPrefix(out, "refused --right: ") {
		t.Fatalf("out = %q, want it to name line 3 and the --right flag", out)
	}
	if _, err := os.Stat(journalPath(t, root)); err == nil {
		t.Fatal("journal should not have been created on refusal")
	}
}

// TestLearningRoundTripsThroughParser is the verb's whole reason to exist: the entry it
// appends is one the journal parser reads back as open, with no malformed record, and
// the parser's own strict `[open]` rule is what grades it.
func TestLearningRoundTripsThroughParser(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(journalPath(t, root), []byte(learnings.JournalSchemaHeading+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := LearningCommand(learningArgs)
	if code != 0 || out != "captured: the gate hid a rule\n" {
		t.Fatalf("stdout/exit = %q/%d, want captured on exit 0", out, code)
	}
	data, err := os.ReadFile(journalPath(t, root))
	if err != nil {
		t.Fatal(err)
	}
	entries, malformed := learnings.Parse(data)
	if len(malformed) != 0 {
		t.Fatalf("appended entry is malformed: %+v\n%s", malformed, data)
	}
	if len(entries) != 1 || entries[0].Title != "the gate hid a rule" || entries[0].State != "open" {
		t.Fatalf("entries = %+v, want one open entry", entries)
	}
	want := "- **What happened:** it failed twice\n- **Right behavior:** read the map first\n- **Proposed rule change:** none"
	if entries[0].Body != want {
		t.Fatalf("body = %q, want %q", entries[0].Body, want)
	}
}

// TestLearningTwoEntriesStayApart covers a second append onto a file the first one
// ended: a blank line separates the headings, so both parse as their own entry.
func TestLearningTwoEntriesStayApart(t *testing.T) {
	root := newRepo(t)
	for i := 0; i < 2; i++ {
		if _, code := LearningCommand(append([]string{"--rule", "add a verb"}, learningArgs...)); code != 0 {
			t.Fatalf("append %d exited %d", i, code)
		}
	}
	data, err := os.ReadFile(journalPath(t, root))
	if err != nil {
		t.Fatal(err)
	}
	entries, malformed := learnings.Parse(data)
	if len(entries) != 2 || len(malformed) != 0 {
		t.Fatalf("entries=%d malformed=%d\n%s", len(entries), len(malformed), data)
	}
}

// headCommit returns the full ID of HEAD in the fixture checkout.
func headCommit(t *testing.T, root string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("rev-parse HEAD: %v", err)
	}
	return strings.TrimSpace(string(out))
}

// TestLearningResolvesLinksTheNamedOpenEntry covers the link verb: an abbreviated commit
// resolves to its full ID, only the entry with the exact title carries it, and the entry
// stays open and well-formed for the parser.
func TestLearningResolvesLinksTheNamedOpenEntry(t *testing.T) {
	root := newRepo(t)
	for _, title := range []string{"first lesson", "second lesson"} {
		if out, code := LearningCommand([]string{title, "--what", "it failed", "--right", "read the map"}); code != 0 {
			t.Fatalf("seed %q = %q/%d", title, out, code)
		}
	}
	head := headCommit(t, root)
	out, code := LearningCommand([]string{"--resolves", "second lesson", "--commit", head[:10]})
	if code != 0 || out != "resolved: second lesson by "+head+"\n" {
		t.Fatalf("resolve = %q/%d, want the full commit ID on exit 0", out, code)
	}
	data, err := os.ReadFile(journalPath(t, root))
	if err != nil {
		t.Fatal(err)
	}
	entries, malformed := learnings.Parse(data)
	if len(malformed) != 0 || len(entries) != 2 {
		t.Fatalf("entries=%d malformed=%+v\n%s", len(entries), malformed, data)
	}
	if entries[0].ResolvedBy != "" || entries[1].ResolvedBy != head || entries[1].State != "open" {
		t.Fatalf("entries = %+v, want only the second linked to %s and still open", entries, head)
	}
}

// TestLearningResolveRefusalsLeaveTheJournalUnchanged covers each content refusal: exit 1,
// a named cause, and the journal bytes as they were.
func TestLearningResolveRefusalsLeaveTheJournalUnchanged(t *testing.T) {
	root := newRepo(t)
	head := headCommit(t, root)
	for _, title := range []string{"first lesson", "first lesson", "linked lesson"} {
		if out, code := LearningCommand([]string{title, "--what", "it failed", "--right", "read the map"}); code != 0 {
			t.Fatalf("seed %q = %q/%d", title, out, code)
		}
	}
	if out, code := LearningCommand([]string{"--resolves", "linked lesson", "--commit", head}); code != 0 {
		t.Fatalf("link = %q/%d", out, code)
	}
	before, err := os.ReadFile(journalPath(t, root))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, title, commit, want string
	}{
		{"unknown title", "missing lesson", head, "no open entry has that title"},
		{"two entries share the title", "first lesson", head, "more than one open entry has that title"},
		{"entry already linked", "linked lesson", head, "already names a resolving commit"},
		{"commit does not resolve", "linked lesson", "0000000000000000000000000000000000000001", "does not resolve"},
		{"commit reads as an option", "first lesson", "--all", "does not resolve"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, code := LearningCommand([]string{"--resolves", tc.title, "--commit", tc.commit})
			if code != 1 || !strings.Contains(out, tc.want) {
				t.Fatalf("resolve = %q/%d, want exit 1 naming %q", out, code, tc.want)
			}
			after, err := os.ReadFile(journalPath(t, root))
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatalf("journal changed on refusal:\n%s", after)
			}
		})
	}
}

// TestLearningResolveUsageExitsTwo covers the argument shapes the link form refuses: each
// half without the other, and the resolve flag beside a new entry's title or body.
func TestLearningResolveUsageExitsTwo(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"resolves without commit", []string{"--resolves", "a lesson"}},
		{"commit without resolves", []string{"--commit", "HEAD"}},
		{"title joins resolves", []string{"--resolves", "a lesson", "--commit", "HEAD", "extra"}},
		{"what joins resolves", []string{"--resolves", "a lesson", "--commit", "HEAD", "--what", "x"}},
		{"right joins resolves", []string{"--resolves", "a lesson", "--commit", "HEAD", "--right", "x"}},
		{"rule joins resolves", []string{"--resolves", "a lesson", "--commit", "HEAD", "--rule", "x"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			out, code := LearningCommand(tc.args)
			if code != 2 || !strings.HasPrefix(out, "usage: bench learning") {
				t.Fatalf("args %q: got %q/%d, want the usage line on exit 2", tc.args, out, code)
			}
			if _, err := os.Stat(journalPath(t, root)); err == nil {
				t.Fatal("usage refusal created the journal")
			}
		})
	}
}

// TestLearningUsageExitsTwo covers the argument shapes the grammar refuses: no title,
// and a title without both body flags. None of them creates the journal.
func TestLearningUsageExitsTwo(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want string
	}{
		{nil, learningGrammar.Help + "\n"},
		{[]string{"--what", "x", "--right", "y"}, learningGrammar.Help + "\n"},
		{[]string{"title", "--what", "x"}, "usage: bench learning (missing argument: --what and --right)\n"},
		{[]string{"title", "--right", "y"}, "usage: bench learning (missing argument: --what and --right)\n"},
	} {
		root := newRepo(t)
		out, code := LearningCommand(tc.args)
		if code != 2 || out != tc.want {
			t.Fatalf("args %q: got %q/%d, want %q/2", tc.args, out, code, tc.want)
		}
		if _, err := os.Stat(journalPath(t, root)); err == nil {
			t.Fatalf("args %q: journal should not have been created", tc.args)
		}
	}
}

// TestLearningRefusesPrimaryCheckout covers the tracked-journal landing boundary.
func TestLearningRefusesPrimaryCheckout(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	t.Chdir(primary)
	out, code := LearningCommand(learningArgs)
	if code != 1 || out != usage.PrimaryCheckoutRefusal()+"\n" {
		t.Fatalf("primary checkout = %q/%d, want the shared refusal on exit 1", out, code)
	}
	if _, err := os.Stat(journalPath(t, primary)); err == nil {
		t.Fatal("refused learning still appended to the journal")
	}
}

// TestLearningIgnoredJournalRedirectsWorktreeWrite covers an ignored journal from a
// linked worktree: the append lands in the primary checkout's copy, the one that
// survives the worktree's release.
func TestLearningIgnoredJournalRedirectsWorktreeWrite(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	if err := os.WriteFile(filepath.Join(primary, ".gitignore"), []byte(learnings.JournalPath+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", ".gitignore"}, {"commit", "-q", "-m", "ignore journal"}} {
		if out, err := exec.Command("git", append([]string{"-C", primary}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	linked := newLinkedWorktree(t, primary)
	t.Chdir(linked)
	if out, code := LearningCommand(learningArgs); code != 0 {
		t.Fatalf("ignored journal in worktree = %q/%d, want exit 0", out, code)
	}
	if _, err := os.Stat(journalPath(t, linked)); err == nil {
		t.Fatal("append landed in the worktree copy, not the primary checkout")
	}
	data, err := os.ReadFile(journalPath(t, primary))
	if err != nil {
		t.Fatalf("journal not created on primary: %v", err)
	}
	if entries, malformed := learnings.Parse(data); len(entries) != 1 || len(malformed) != 0 {
		t.Fatalf("primary journal: entries=%d malformed=%d\n%s", len(entries), len(malformed), data)
	}
}
