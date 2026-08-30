package preflight

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/diff"
)

// TestCommandConformantTree is C1, the tracer: five green rows by name, exit 0, and a
// byte-identical second run.
func TestCommandConformantTree(t *testing.T) {
	_, slug := seedConformant(t)

	first, code := Command([]string{"review", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, first)
	}
	for _, name := range []string{"base-current", "paths-authorized", "rows-owned", "rows-membership", "diff-nonempty"} {
		if !strings.Contains(first, name+",green") {
			t.Errorf("output missing green %s row:\n%s", name, first)
		}
	}
	// WF32: the rendered header carries the next column on every run, green
	// included. A Next field the renderer never reads would leave the old header.
	if !strings.Contains(first, "checks[5]{check,verdict,detail,next}") {
		t.Errorf("output missing the four-column checks header:\n%s", first)
	}
	second, code2 := Command([]string{"review", slug})
	if code2 != 0 || second != first {
		t.Errorf("second run = (%q, %d), want byte-identical to first (%q, 0)", second, code2, first)
	}
}

// TestCommandStaleBase is C2: the default branch advanced past the branch point makes
// base-current the red row, exit 1.
func TestCommandStaleBase(t *testing.T) {
	_, slug := seedConformant(t)
	runGit(t, "checkout", "-q", "main")
	mustWriteFile(t, "unrelated.txt", "advance main\n")
	runGit(t, "add", "unrelated.txt")
	runGit(t, "commit", "-q", "-m", "advance main")
	runGit(t, "checkout", "-q", "feature")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "base-current,red") {
		t.Errorf("output missing red base-current row:\n%s", out)
	}
	// WF37: the remedy reaches the rendered row, so the gatherer's default-branch
	// and assignment facts are proved beside the pure decision. This tree is no
	// worktree of an assignment, so the id is the placeholder.
	if !strings.Contains(out, "base-current,red,default branch tip is not an ancestor of HEAD,bench worktree merge --from main <target>") {
		t.Errorf("output missing the base-current row's remedy:\n%s", out)
	}
}

// TestCommandOutOfFencePath is C3: a tracked change outside every fence entry makes
// paths-authorized red naming the path.
func TestCommandOutOfFencePath(t *testing.T) {
	_, slug := seedConformant(t)
	mustWriteFile(t, "unfenced/other.go", "package other\n")
	runGit(t, "add", "unfenced/other.go")
	runGit(t, "commit", "-q", "-m", "out of fence")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,red") || !strings.Contains(out, "unfenced/other.go") {
		t.Errorf("output missing red paths-authorized naming unfenced/other.go:\n%s", out)
	}
}

// TestCommandUnstagedOutOfFenceRed is RS1. A tracked file sits outside
// every fence entry, already identical between the review base and HEAD.
// It gets edited in the worktree with no `git add` or commit. It still
// makes paths-authorized red naming the path. This proves the
// changed-set resolution folds in tracked-worktree edits and not just
// the base..HEAD commit range.
//
// The file is committed at the shared base commit itself (never recommitted
// on feature), so a base..HEAD diff sees no change to it at all. Only the
// uncommitted worktree edit distinguishes it.
func TestCommandUnstagedOutOfFenceRed(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing PF1 and PF2.\n")
	mustWriteFile(t, "unfenced/other.go", "package other\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")
	mustWriteFile(t, "unfenced/other.go", "package other\n\n// edited, never staged\n")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,red") || !strings.Contains(out, "unfenced/other.go") {
		t.Errorf("output missing red paths-authorized naming unfenced/other.go:\n%s", out)
	}
}

// TestCommandFencePrefixBoundary is C3's prefix-boundary half: internal/git2 must not
// match a fence entry of internal/git.
func TestCommandFencePrefixBoundary(t *testing.T) {
	slug := "boundary"
	initRepo(t)
	body := "# boundary\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n" +
		"|---|---|---|---|---|\n| PF1 | 1 | does x | cli seam | catches z |\n\n" +
		"## Ownership fences\n\n- `internal/git`\n- `reviews/boundary.md`\n"
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing PF1.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/git2/thing.go", "package git2\n")
	runGit(t, "add", "internal/git2/thing.go")
	runGit(t, "commit", "-q", "-m", "c1")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,red") || !strings.Contains(out, "internal/git2/thing.go") {
		t.Errorf("fence internal/git must not authorize internal/git2/thing.go:\n%s", out)
	}
}

// TestCommandUncitedRow is C4: one declared row ID cited by no ticket file makes
// rows-owned red naming the uncited ID.
func TestCommandUncitedRow(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing only PF1.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "rows-owned,red") || !strings.Contains(out, "PF2") {
		t.Errorf("output missing red rows-owned naming PF2:\n%s", out)
	}
}

// TestCommandPhantomAndForeignTag is C5: a ticket token under the spec's own tag
// naming no declared row (PF99) makes rows-membership red; a foreign-tag token (FT93)
// is ignored.
func TestCommandPhantomAndForeignTag(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Cites PF1, PF2, PF99, and FT93.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "rows-membership,red") || !strings.Contains(out, "PF99") {
		t.Errorf("output missing red rows-membership naming PF99:\n%s", out)
	}
	if strings.Contains(out, "FT93") {
		t.Errorf("a foreign-tag token must be ignored, not named in any row:\n%s", out)
	}
}

// TestCommandTicketsSubdirRowOwned is RG2/subdir-row-owned: a declared row cited only
// by a ticket file nested under tickets/sub/ is owned — enumeration recurses rather
// than skipping the subdirectory.
func TestCommandTicketsSubdirRowOwned(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing only PF1.\n")
	mustWriteFile(t, "specs/"+slug+"/tickets/sub/x.md", "Nested ticket citing PF2.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")

	out, code := Command([]string{"review", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	if !strings.Contains(out, "rows-owned,green") {
		t.Errorf("a row cited only under tickets/sub/ must still be owned:\n%s", out)
	}
}

// TestCommandTicketsSubdirPhantomDetected is RG2/subdir-phantom-detected: a phantom
// own-tag token that appears only in a nested tickets/sub/ file is still detected.
func TestCommandTicketsSubdirPhantomDetected(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing PF1 and PF2.\n")
	mustWriteFile(t, "specs/"+slug+"/tickets/sub/x.md", "Nested ticket citing PF99.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "rows-membership,red") || !strings.Contains(out, "PF99") {
		t.Errorf("a phantom token nested under tickets/sub/ must still be caught:\n%s", out)
	}
}

// TestCommandEmptyDiff is C6: an empty changed set in review mode makes diff-nonempty
// red.
func TestCommandEmptyDiff(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing PF1 and PF2.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature") // no further commits: HEAD == merge-base

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "diff-nonempty,red") {
		t.Errorf("output missing red diff-nonempty:\n%s", out)
	}
}

// TestCommandControlBytePath is C7's refusal half. A changed path
// carrying an ESC control byte exits 1 with the unrepresentable-TOON-cell
// error rather than a mangled table or a silently sanitized path.
//
// The refusal is unconditional. It fires before verdict rendering over
// every changed path, not only one that would otherwise reach a red
// row's detail cell. See TestCommandFencedControlBytePathReds for the
// fenced, otherwise-all-green case that exercises that distinction
// directly.
func TestCommandControlBytePath(t *testing.T) {
	_, slug := seedConformant(t)
	hostile := "unfenced/a\x1bb.go"
	mustWriteFile(t, hostile, "package example\n")
	runGit(t, "add", "--", hostile)
	runGit(t, "commit", "-q", "-m", "hostile path")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "unrepresentable TOON cell") {
		t.Errorf("output = %q, want the unrepresentable-TOON-cell error", out)
	}
	if strings.Contains(out, "checks[") {
		t.Errorf("a control-byte path must never render a checks table:\n%s", out)
	}
}

// TestCommandFencedControlBytePathReds is RG3/fenced-esc-path-reds. A
// changed path carrying an ESC control byte still exits 1 with the
// unrepresentable-TOON-cell error. This holds even when the path is
// authorized by an ownership fence and every check would otherwise be
// green. PF7's refusal is unconditional, not gated on the path reaching
// a red row's rendered detail cell.
func TestCommandFencedControlBytePathReds(t *testing.T) {
	_, slug := seedConformant(t)
	hostile := "internal/" + slug + "/a\x1bb.go"
	mustWriteFile(t, hostile, "package example\n")
	runGit(t, "add", "--", hostile)
	runGit(t, "commit", "-q", "-m", "fenced hostile path")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.Contains(out, "unrepresentable TOON cell") {
		t.Errorf("output = %q, want the unrepresentable-TOON-cell error", out)
	}
	if strings.Contains(out, "checks[") {
		t.Errorf("a fenced control-byte path must never render a checks table:\n%s", out)
	}
}

// TestCommandSpaceAndGlobPath is C7's render half: a path with a space or glob
// character renders escaped and authorizes correctly rather than false-reddening.
func TestCommandSpaceAndGlobPath(t *testing.T) {
	_, slug := seedConformant(t)
	fancy := "internal/" + slug + "/a b*.go"
	mustWriteFile(t, fancy, "package example\n")
	runGit(t, "add", "--", fancy)
	runGit(t, "commit", "-q", "-m", "space and glob path")

	out, code := Command([]string{"review", slug})
	if code != 0 {
		t.Fatalf("Command exit = %d, want 0; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,green") {
		t.Errorf("a fenced space/glob path must stay authorized:\n%s", out)
	}
}

// TestCommandRecordedBaseKey is C8: a recorded branch.<name>.benchBase past
// an out-of-fence commit keeps paths-authorized green. Removing the key
// falls back to merge-base and turns the same tree red. The CLI observably
// consumes the recorded-key resolution bench diff itself uses. Bare bench
// diff output stays byte-identical.
func TestCommandRecordedBaseKey(t *testing.T) {
	_, slug := seedConformant(t) // HEAD = c1 on feature, base commit = c0
	mustWriteFile(t, "unfenced/other.go", "package other\n")
	runGit(t, "add", "unfenced/other.go")
	runGit(t, "commit", "-q", "-m", "out of fence")
	c2 := runGit(t, "rev-parse", "HEAD") // past (includes) the out-of-fence commit
	mustWriteFile(t, "internal/"+slug+"/bar.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/bar.go")
	runGit(t, "commit", "-q", "-m", "authorized change after the out-of-fence commit")

	beforeDiff, beforeCode := diff.Command(nil)

	runGit(t, "config", "branch.feature.benchBase", c2)
	out, code := Command([]string{"review", slug})
	if code != 0 {
		t.Fatalf("with benchBase recorded past the out-of-fence commit, Command exit = %d, want 0; output:\n%s", code, out)
	}
	if !strings.Contains(out, "paths-authorized,green") {
		t.Errorf("recorded-key base must exclude the out-of-fence commit from the diff:\n%s", out)
	}

	runGit(t, "config", "--unset", "branch.feature.benchBase")
	out2, code2 := Command([]string{"review", slug})
	if code2 != 1 {
		t.Fatalf("with the key removed, Command exit = %d, want 1; output:\n%s", code2, out2)
	}
	if !strings.Contains(out2, "paths-authorized,red") {
		t.Errorf("merge-base fallback must surface the out-of-fence commit:\n%s", out2)
	}

	afterDiff, afterCode := diff.Command(nil)
	if afterDiff != beforeDiff || afterCode != beforeCode {
		t.Errorf("bare bench diff output changed across preflight's use of the export: before=(%q,%d) after=(%q,%d)", beforeDiff, beforeCode, afterDiff, afterCode)
	}
}
