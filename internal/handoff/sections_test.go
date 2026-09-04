package handoff

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/status"
)

// HS6, stories 8 and 11. Two live assignments each write their own section. The second
// writer must leave the first one's bytes exactly as they were: a run that rendered the
// whole document from its own facts would drop the sibling's pins and its State with them.
func TestCommandRewritesOneSectionAndLeavesTheSiblingIdentical(t *testing.T) {
	root := benchRepo(t)
	first := activeAssignment(t, root, "hs-first", "first tree")
	second := activeAssignment(t, root, "hs-second", "second tree")
	document := filepath.Join(root, status.HandoffFile)

	runIn(t, first.Worktree, nil)
	planted := sectionBytes(t, document, first.Request)
	if len(planted) == 0 {
		t.Fatalf("the first run wrote no section for %q\n%s", first.Request, read(t, document))
	}

	runIn(t, second.Worktree, nil)
	kept := sectionBytes(t, document, first.Request)
	if planted != kept {
		t.Fatalf("the sibling section changed\nbefore:\n%s\nafter:\n%s", planted, kept)
	}
	if len(sectionBytes(t, document, second.Request)) == 0 {
		t.Fatalf("the second run wrote no section for %q\n%s", second.Request, read(t, document))
	}
}

// HS7, story 9. The primary checkout owns main, so a phase close with nothing live still
// writes. A verb that demanded an assignment would refuse that close.
func TestCommandFromThePrimaryCheckoutWritesMain(t *testing.T) {
	root := benchRepo(t)
	runIn(t, root, nil)
	if !strings.Contains(read(t, filepath.Join(root, status.HandoffFile)), handoffdoc.MainHeading+"\n") {
		t.Fatalf("primary-checkout handoff carries no %q section\n%s", handoffdoc.MainHeading, read(t, filepath.Join(root, status.HandoffFile)))
	}
}

// HS8, story 10. A checkout that is neither primary nor an active assignment's tree owns
// no section. It refuses and writes nothing: a match by label or by path string would let
// it adopt a section some other request owns.
func TestCommandRefusesACheckoutThatOwnsNoSection(t *testing.T) {
	root := benchRepo(t)
	runIn(t, root, nil)
	document := filepath.Join(root, status.HandoffFile)
	before := read(t, document)

	stray := linkedWorktree(t, root, "bench/stray")
	out, code := runAt(t, stray, nil)
	if code != 1 {
		t.Fatalf("handoff from an unowned checkout = (%q, %d), want exit 1", out, code)
	}
	if after := read(t, document); after != before {
		t.Fatalf("a refused run rewrote the document\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// benchRepo builds a primary checkout that ignores its capture directory, which is the
// decided shape: the handoff is a local pin the primary checkout holds, so a worktree run
// writes there rather than into its own tree.
func benchRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q", "-b", "main")
	gitRun(t, root, "config", "user.email", "t@example.com")
	gitRun(t, root, "config", "user.name", "t")
	write(t, filepath.Join(root, ".gitignore"), "capture/\n")
	write(t, filepath.Join(root, ".bench", "keep"), "\n")
	write(t, filepath.Join(root, "ROADMAP.md"), "# Roadmap\n")
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-qm", "base")
	return root
}

// nextIdentity hands out the 32-hex owner and assignment identities the ledger validates.
// Each call takes the next pair, so two assignments in one repo never collide.
var nextIdentity = 0

// activeAssignment registers one active assignment owning a fresh worktree, and returns
// the record. It goes in through intent.PutAssignment, so the command reads exactly the
// shape every worktree command writes.
func activeAssignment(t *testing.T, root, token, label string) intent.Assignment {
	t.Helper()
	nextIdentity++
	owner := fmt.Sprintf("%032x", 2*nextIdentity)
	id := fmt.Sprintf("%032x", 2*nextIdentity+1)
	branch := intent.AssignmentBranchRef(owner, id)
	a := intent.Assignment{
		Schema:       intent.AssignmentRecordSchema,
		ID:           id,
		OwnerID:      owner,
		Request:      intent.RequestDigest(token),
		RequestToken: token,
		Label:        label,
		Start:        gitOut(t, root, "rev-parse", "HEAD"),
		Branch:       branch,
		Worktree:     linkedWorktree(t, root, strings.TrimPrefix(branch, "refs/heads/")),
		State:        intent.StateActive,
	}
	if err := intent.PutAssignment(root, a); err != nil {
		t.Fatalf("PutAssignment: %v", err)
	}
	return a
}

// linkedWorktree adds one linked checkout on a new branch and returns its cleaned path.
// The ledger stores a cleaned absolute path, and so does the record built over this.
func linkedWorktree(t *testing.T, root, branch string) string {
	t.Helper()
	path := filepath.Clean(filepath.Join(t.TempDir(), "tree"))
	gitRun(t, root, "worktree", "add", "-q", "-b", branch, path, "HEAD")
	return path
}

// runIn runs the command from dir and fails the test on any non-zero exit.
func runIn(t *testing.T, dir string, args []string) string {
	t.Helper()
	out, code := runAt(t, dir, args)
	if code != 0 {
		t.Fatalf("handoff from %s = (%q, %d), want exit 0", dir, out, code)
	}
	return out
}

// runAt runs the command with dir as the working directory and returns its output and
// exit code. The command resolves its own root from the process directory, so the chdir
// is what puts it in one checkout rather than another.
func runAt(t *testing.T, dir string, args []string) (string, int) {
	t.Helper()
	t.Chdir(dir)
	return Command(args)
}

// sectionBytes returns one section's bytes, from its heading to the next level-two
// heading. A comparison over these bytes is what catches a rewrite that regenerated the
// sibling's fields instead of re-emitting them.
func sectionBytes(t *testing.T, path, key string) string {
	t.Helper()
	heading := handoffdoc.MainHeading
	if key != handoffdoc.MainKey {
		heading = handoffdoc.RequestHeadingPrefix + key
	}
	_, after, found := strings.Cut(read(t, path), heading+"\n")
	if !found {
		return ""
	}
	if end := strings.Index(after, "\n## "); end >= 0 {
		return heading + "\n" + after[:end]
	}
	return heading + "\n" + after
}

func read(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func gitRun(t *testing.T, root string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func gitOut(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}
