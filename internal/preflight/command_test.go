package preflight

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/toon"
)

// --- fixture harness: mirrors internal/coverage's TestCommand and internal/diff's
// review_base_test.go patterns — real git commands in a throwaway repo, exact
// output/exit assertions against the public Command entry point only. ---

func runGit(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func mustWriteFile(t *testing.T, path, body string) {
	t.Helper()
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q): %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// specBody renders a bootstrap-conformant spec for slug: staged status, a valid
// opted-in coverage map declaring row PF1 and PF2, and one backticked fence entry
// authorizing internal/<slug>/.
func specBody(slug string, extraFenceLines ...string) string {
	var b strings.Builder
	b.WriteString("# " + slug + "\n\nStatus: staged\n\n")
	b.WriteString("## User stories\n1. As a, I want b, so c.\n\n")
	b.WriteString("### Acceptance coverage map\n")
	b.WriteString("| row | story | behavior | seam | red signal | why it catches the failure |\n")
	b.WriteString("|---|---|---|---|---|---|\n")
	b.WriteString("| PF1 | 1 | does x | cli seam | cmd fails | catches z |\n")
	b.WriteString("| PF2 | 1 | does y | cli seam | cmd fails | catches w |\n")
	b.WriteString("\n## Ownership fences\n\n")
	b.WriteString("- `internal/" + slug + "/` (implementation)\n")
	for _, line := range extraFenceLines {
		b.WriteString(line + "\n")
	}
	return b.String()
}

// initRepo starts a fresh git repo at t.TempDir(), chdir'd into it, with an author
// identity configured — the common prefix of every seed function below.
func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	runGit(t, "init", "-q", "-b", "main")
	runGit(t, "config", "user.email", "t@example.com")
	runGit(t, "config", "user.name", "t")
	return root
}

// seedConformant builds the PF1 tracer fixture: a base commit on main carrying a
// bootstrap-conformant spec and a ticket citing both declared rows, then a feature
// branch with one authorized change — the tree every check should answer green over.
func seedConformant(t *testing.T) (root, slug string) {
	t.Helper()
	slug = "example"
	root = initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	mustWriteFile(t, "specs/"+slug+"/tickets/one.md", "Ticket citing PF1 and PF2.\n")
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")
	runGit(t, "checkout", "-q", "-b", "feature")
	mustWriteFile(t, "internal/"+slug+"/foo.go", "package example\n")
	runGit(t, "add", "internal/"+slug+"/foo.go")
	runGit(t, "commit", "-q", "-m", "c1")
	return root, slug
}

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

// TestCommandFencePrefixBoundary is C3's prefix-boundary half: internal/git2 must not
// match a fence entry of internal/git.
func TestCommandFencePrefixBoundary(t *testing.T) {
	slug := "boundary"
	initRepo(t)
	body := "# boundary\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n| row | story | behavior | seam | red signal | why it catches the failure |\n" +
		"|---|---|---|---|---|---|\n| PF1 | 1 | does x | cli seam | cmd fails | catches z |\n\n" +
		"## Ownership fences\n\n- `internal/git`\n"
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

// TestCommandControlBytePath is C7's refusal half: a changed path carrying an ESC
// control byte exits 1 with the unrepresentable-TOON-cell error rather than a mangled
// table or a silently sanitized path.
func TestCommandControlBytePath(t *testing.T) {
	// The hostile path is deliberately unfenced: an authorized path never reaches the
	// rendered detail, so only an unauthorized (and therefore named) control-byte
	// path exercises the TOON sink's refusal.
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

// TestCommandRecordedBaseKey is C8: a recorded branch.<name>.benchBase past an
// out-of-fence commit keeps paths-authorized green; removing the key falls back to
// merge-base and turns the same tree red — the CLI observably consumes the recorded-
// key resolution bench diff itself uses. Bare bench diff output stays byte-identical.
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

// TestCommandTicketsAbsent is C12's absent half: review mode with tickets/ absent is a
// structured red.
func TestCommandTicketsAbsent(t *testing.T) {
	slug := "example"
	initRepo(t)
	mustWriteFile(t, "specs/"+slug+"/spec.md", specBody(slug))
	runGit(t, "add", ".")
	runGit(t, "commit", "-q", "-m", "c0")

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error:") {
		t.Errorf("output = %q, want a structured error", out)
	}
}

// TestCommandSpecialFileInTickets is C12's refusal half: a FIFO inside tickets/ is
// refused before reading, named in the error, without blocking. The goroutine/timeout
// wrapper is the test-level bound the ticket requires: a hang mutation fails the test
// instead of the suite.
func TestCommandSpecialFileInTickets(t *testing.T) {
	_, slug := seedConformant(t)
	fifoPath := filepath.Join("specs", slug, "tickets", "fifo")
	if err := syscall.Mkfifo(fifoPath, 0o644); err != nil {
		t.Fatalf("Mkfifo: %v", err)
	}

	type result struct {
		out  string
		code int
	}
	done := make(chan result, 1)
	go func() {
		out, code := Command([]string{"review", slug})
		done <- result{out, code}
	}()
	select {
	case r := <-done:
		if r.code != 1 {
			t.Fatalf("Command exit = %d, want 1; output:\n%s", r.code, r.out)
		}
		if !strings.Contains(r.out, "fifo") {
			t.Errorf("output must name the refused fifo:\n%s", r.out)
		}
	case <-time.After(bounds.TestDeadline(5 * time.Second)):
		t.Fatal("Command blocked reading a FIFO instead of refusing it by mode")
	}
}

// TestCommandNotInRepo is C13's not-in-repo branch: outside a git repository the
// command answers the standard not-in-repo error.
func TestCommandNotInRepo(t *testing.T) {
	t.Chdir(t.TempDir())
	out, code := Command([]string{"review", "example"})
	if want := toon.NotInRepo() + "\n"; out != want || code != 1 {
		t.Errorf("Command = (%q, %d), want (%q, 1)", out, code, want)
	}
}

// TestCommandUsageBranches is C13's four usage branches: missing mode, unknown mode,
// missing slug, and unknown flag each exit 2.
func TestCommandUsageBranches(t *testing.T) {
	initRepo(t)
	runGit(t, "commit", "-q", "--allow-empty", "-m", "c0")

	cases := []struct {
		name string
		args []string
	}{
		{"missing mode", nil},
		{"unknown mode", []string{"frob", "example"}},
		{"build mode not yet accepted", []string{"build", "example"}},
		{"missing slug", []string{"review"}},
		{"unknown flag", []string{"--frob"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, code := Command(c.args)
			if code != 2 {
				t.Errorf("Command(%v) = (%q, %d), want exit 2", c.args, out, code)
			}
		})
	}
}

// TestCommandMissingSpecNamesPath is H1's missing half (mutation
// H1/missing-spec-names-path): no spec at all answers "spec not found" naming the
// tried path, exit 1.
func TestCommandMissingSpecNamesPath(t *testing.T) {
	slug := "example"
	initRepo(t)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: spec not found") {
		t.Errorf("output = %q, want the spec-not-found error", out)
	}
	if !strings.Contains(out, filepath.Join("specs", slug, "spec.md")) {
		t.Errorf("output must name the tried spec path:\n%s", out)
	}
}

// TestCommandDanglingSymlinkClassifiedBroken is H1's symlink half (mutation
// H1/dangling-symlink-classified): a dangling symlink where the spec should be is
// classified broken, not an authoritative empty state — a distinct error from the
// missing-spec case, naming the spec path.
func TestCommandDanglingSymlinkClassifiedBroken(t *testing.T) {
	slug := "example"
	initRepo(t)
	specPath := filepath.Join("specs", slug, "spec.md")
	if err := os.MkdirAll(filepath.Dir(specPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.Symlink("missing-target.md", specPath); err != nil {
		t.Fatalf("Symlink: %v", err)
	}

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if strings.HasPrefix(out, "error: spec not found") {
		t.Errorf("a dangling symlink must not be classified as authoritative absence:\n%s", out)
	}
	if !strings.HasPrefix(out, "error: spec not readable") {
		t.Errorf("output = %q, want the spec-not-readable error", out)
	}
	if !strings.Contains(out, specPath) {
		t.Errorf("output must name the spec path:\n%s", out)
	}
}

// TestCommandInvalidCoverageMapCarriesMessage is H2's invalid-map half (mutation
// H2/validator-message-carried): a coverage map that fails validation surfaces the
// validator's own message, exit 1.
func TestCommandInvalidCoverageMapCarriesMessage(t *testing.T) {
	slug := "example"
	initRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| row | story | behavior | seam | red signal | why it catches the failure |\n" +
		"|---|---|---|---|---|---|\n" +
		"| PF1 | 1 | does x | cli seam | cmd fails |\n" + // 5 cells where 6 are wanted
		"\n## Ownership fences\n\n- `internal/" + slug + "/`\n"
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: coverage map invalid") {
		t.Errorf("output = %q, want the coverage-map-invalid error", out)
	}
	if !strings.Contains(out, "row 1 has 5 cells (want 6)") {
		t.Errorf("output must carry the validator's own message:\n%s", out)
	}
}

// TestCommandLegacyMapNamesOptIn is H2's legacy-map half (mutation
// H2/optin-hint-named): a valid legacy 5-cell map (no `row` header column) is
// refused with an error naming the row-ID opt-in, exit 1.
func TestCommandLegacyMapNamesOptIn(t *testing.T) {
	slug := "example"
	initRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| story | behavior | seam | red signal | why it catches the failure |\n" +
		"|---|---|---|---|---|\n" +
		"| 1 | does x | cli seam | cmd fails | catches z |\n" +
		"\n## Ownership fences\n\n- `internal/" + slug + "/`\n"
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: coverage map not opted into row IDs") {
		t.Errorf("output = %q, want the opt-in error naming the row-ID convention", out)
	}
	if !strings.Contains(out, "`row`") {
		t.Errorf("output must name the leading `row` column opt-in:\n%s", out)
	}
}

// TestCommandFencesAbsentError is H3's absent half (mutation H3/fences-absent-error):
// a spec with no `## Ownership fences` section at all answers a structured error,
// exit 1.
func TestCommandFencesAbsentError(t *testing.T) {
	slug := "example"
	initRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| row | story | behavior | seam | red signal | why it catches the failure |\n" +
		"|---|---|---|---|---|---|\n" +
		"| PF1 | 1 | does x | cli seam | cmd fails | catches z |\n"
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: ownership fences empty") {
		t.Errorf("output = %q, want the ownership-fences error", out)
	}
}

// TestCommandFencesEmptyError is H3's empty half (mutation H3/fences-empty-error):
// an `## Ownership fences` section present but holding no entry answers the same
// structured error as the absent case, exit 1.
func TestCommandFencesEmptyError(t *testing.T) {
	slug := "example"
	initRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| row | story | behavior | seam | red signal | why it catches the failure |\n" +
		"|---|---|---|---|---|---|\n" +
		"| PF1 | 1 | does x | cli seam | cmd fails | catches z |\n" +
		"\n## Ownership fences\n\n## Next section\n"
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: ownership fences empty") {
		t.Errorf("output = %q, want the ownership-fences error", out)
	}
}

// TestCommandFencesParenTokenNeverAuthorizes is H3's parenthesized-token half
// (mutation H3/paren-token-never-authorizes): a fences section whose only
// backticked token sits inside parentheses answers the same structured error as an
// empty section, exit 1 — a parenthetical mention is never an authorization.
func TestCommandFencesParenTokenNeverAuthorizes(t *testing.T) {
	slug := "example"
	initRepo(t)
	body := "# " + slug + "\n\nStatus: staged\n\n## User stories\n1. As a, I want b, so c.\n\n" +
		"### Acceptance coverage map\n" +
		"| row | story | behavior | seam | red signal | why it catches the failure |\n" +
		"|---|---|---|---|---|---|\n" +
		"| PF1 | 1 | does x | cli seam | cmd fails | catches z |\n" +
		"\n## Ownership fences\n\n- see also (`internal/" + slug + "/`)\n"
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: ownership fences empty") {
		t.Errorf("a parenthesized backticked token must not authorize: output = %q", out)
	}
}

// TestCommandNonStagedNamesFoundStatus is H4 (mutation H4/found-status-named): a
// spec whose Status: is anything but staged answers a structured error naming the
// found status, exit 1.
func TestCommandNonStagedNamesFoundStatus(t *testing.T) {
	slug := "example"
	initRepo(t)
	body := strings.Replace(specBody(slug), "Status: staged", "Status: implemented", 1)
	mustWriteFile(t, "specs/"+slug+"/spec.md", body)

	out, code := Command([]string{"review", slug})
	if code != 1 {
		t.Fatalf("Command exit = %d, want 1; output:\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: spec not staged") {
		t.Errorf("output = %q, want the spec-not-staged error", out)
	}
	if !strings.Contains(out, "Status: implemented") {
		t.Errorf("output must name the found status:\n%s", out)
	}
}

// TestCommandTrailingNewlineParity is C14: a spec whose last line lacks a trailing
// newline, and a ticket file whose last line lacks one, each parse identically to
// their terminated forms.
func TestCommandTrailingNewlineParity(t *testing.T) {
	slug := "example"
	terminated := specBody(slug)
	unterminated := strings.TrimSuffix(terminated, "\n")
	if unterminated == terminated {
		t.Fatal("fixture invalid: specBody already lacks a trailing newline")
	}

	for _, tc := range []struct {
		name      string
		specBody  string
		ticketDoc string
	}{
		{"terminated", terminated, "Ticket citing PF1 and PF2.\n"},
		{"unterminated spec", unterminated, "Ticket citing PF1 and PF2.\n"},
		{"unterminated ticket", terminated, "Ticket citing PF1 and PF2."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			initRepo(t)
			mustWriteFile(t, "specs/"+slug+"/spec.md", tc.specBody)
			mustWriteFile(t, "specs/"+slug+"/tickets/one.md", tc.ticketDoc)
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
			for _, name := range []string{"rows-owned,green", "rows-membership,green"} {
				if !strings.Contains(out, name) {
					t.Errorf("%s: output missing %s (trailing-newline handling dropped the last citation):\n%s", tc.name, name, out)
				}
			}
		})
	}
}
