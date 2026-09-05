package structure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// The growth mode's fixtures are the since-mode fixtures plus a base commit to
// compare against: initRepo, write, add, and lines all come from structure_test.go.
// commit and headSha name the two steps every growth fixture repeats.
func commit(t *testing.T, root, msg string) {
	t.Helper()
	add(t, root)
	run(t, root, "commit", "-m", msg)
}

func headSha(t *testing.T, root string) string {
	t.Helper()
	return run(t, root, "rev-parse", "HEAD")
}

// growthRun runs `bench structure --growth <base>` against root and returns the
// report and the exit code, the way an operator invokes the mode.
func growthRun(t *testing.T, root, base string) (string, int) {
	t.Helper()
	t.Chdir(root)
	return Command([]string{"--growth", base})
}

// SR42: an over-budget file that gained lines prints one FILE GREW row carrying the
// tip count, the base count, the limit, and the path, then one summary line, at exit 1.
func TestGrowthOverBudgetFileGrew(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "big.go", lines(12))
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "big.go", lines(15))
	commit(t, root, "grow")

	report, code := growthRun(t, root, base)
	if want := "FILE GREW       15 lines, was 12 (max 10)   big.go"; !strings.Contains(report, want) {
		t.Errorf("missing growth row %q:\n%s", want, report)
	}
	if !strings.Contains(report, "structure growth: 1 file(s) grew past budget.") {
		t.Errorf("missing growth summary line:\n%s", report)
	}
	if code != 1 {
		t.Errorf("exit = %d, want 1:\n%s", code, report)
	}

	// SR57's revert half: the same command over the same base is silent once the
	// planted growth is undone, so the ratchet releases what it caught.
	write(t, root, "big.go", lines(12))
	commit(t, root, "revert the growth")
	report, code = growthRun(t, root, base)
	if code != 0 || !strings.Contains(report, "structure growth ok") {
		t.Errorf("the reverted growth stayed red at exit %d:\n%s", code, report)
	}
}

// SR43: an over-budget file that lost lines, and one that kept its count while its
// content changed, are both silent. Existing debt stays soft, and a split is never punished.
func TestGrowthShrankOrHeldIsSilent(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "shrink.go", lines(20))
	write(t, root, "held.go", lines(20))
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "shrink.go", lines(15))
	write(t, root, "held.go", strings.Repeat("x\n", 19)+"y\n")
	commit(t, root, "edit")

	report, code := growthRun(t, root, base)
	if strings.Contains(report, "FILE GREW") {
		t.Errorf("a file that shrank or held its count was flagged:\n%s", report)
	}
	if code != 0 {
		t.Errorf("exit = %d, want 0:\n%s", code, report)
	}
	if !strings.Contains(report, "structure growth ok") {
		t.Errorf("missing the ok line:\n%s", report)
	}
}

// SR44: a file at or under its limit is free to gain lines. The limit is half the rule.
func TestGrowthUnderLimitGrewIsSilent(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "small.go", lines(2))
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "small.go", lines(10))
	commit(t, root, "grow within the limit")

	report, code := growthRun(t, root, base)
	if strings.Contains(report, "FILE GREW") {
		t.Errorf("an edit that stayed at the limit was flagged:\n%s", report)
	}
	if code != 0 || !strings.Contains(report, "structure growth ok") {
		t.Errorf("want the ok line at exit 0, got exit %d:\n%s", code, report)
	}
}

// SR45: a file absent at the base counts zero, so a fresh oversized file reds with `was 0`.
func TestGrowthAddedOverLimitCountsZeroAtBase(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "old.go", lines(1))
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "fresh.go", lines(15))
	commit(t, root, "add a big file")

	report, code := growthRun(t, root, base)
	if want := "FILE GREW       15 lines, was 0 (max 10)   fresh.go"; !strings.Contains(report, want) {
		t.Errorf("missing added-file row %q:\n%s", want, report)
	}
	if code != 1 {
		t.Errorf("exit = %d, want 1:\n%s", code, report)
	}
}

// SR46: a structure-accept row exempts its file, so a reviewer grant holds against growth too.
func TestGrowthAcceptedFileGrewIsSilent(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "big.go", lines(12))
	write(t, root, ".bench/structure-accept", "big.go cohesive and barely over budget\n")
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "big.go", lines(15))
	commit(t, root, "grow")

	report, code := growthRun(t, root, base)
	if strings.Contains(report, "FILE GREW") {
		t.Errorf("a granted file was flagged for growth:\n%s", report)
	}
	if code != 0 || !strings.Contains(report, "structure growth ok") {
		t.Errorf("want the ok line at exit 0, got exit %d:\n%s", code, report)
	}
}

// SR47: the limit is the engine's limit for the path. A `structure.budgets` row governs
// both directions: silence within the grant, and a row naming that grant past it.
func TestGrowthHonorsPerPathBudget(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "wide.go", lines(12))
	write(t, root, ".bench/structure.budgets", "wide.go 20\n")
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "wide.go", lines(18))
	commit(t, root, "grow within the grant")
	report, code := growthRun(t, root, base)
	if strings.Contains(report, "FILE GREW") {
		t.Errorf("growth within the granted budget was flagged:\n%s", report)
	}
	if code != 0 {
		t.Errorf("exit = %d, want 0:\n%s", code, report)
	}

	write(t, root, "wide.go", lines(25))
	commit(t, root, "grow past the grant")
	report, code = growthRun(t, root, base)
	if want := "FILE GREW       25 lines, was 12 (max 20)   wide.go"; !strings.Contains(report, want) {
		t.Errorf("missing granted-budget row %q:\n%s", want, report)
	}
	if code != 1 {
		t.Errorf("exit = %d, want 1:\n%s", code, report)
	}
}

// SR48: an exact rename pairs with its old path, so a pure move of an over-budget file
// compares against the same count and passes. Read as an addition it would red at `was 0`.
func TestGrowthExactRenameIsSilent(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "old.go", lines(12))
	commit(t, root, "base")
	base := headSha(t, root)

	run(t, root, "mv", "old.go", "new.go")
	commit(t, root, "move")

	report, code := growthRun(t, root, base)
	if strings.Contains(report, "FILE GREW") {
		t.Errorf("a pure rename was read as an addition:\n%s", report)
	}
	if code != 0 || !strings.Contains(report, "structure growth ok") {
		t.Errorf("want the ok line at exit 0, got exit %d:\n%s", code, report)
	}
}

// SR49: the query is NUL-framed, so a name with a byte above ASCII arrives as its own
// bytes rather than the C-quoted spelling `core.quotepath` produces for a newline frame.
// A C-quoted name would also lose the `.go` anchor and drop out of the source filter.
func TestGrowthNonASCIIPathKeepsItsBytes(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "old.go", lines(1))
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "café.go", lines(15))
	commit(t, root, "add a big file with an accented name")

	report, code := growthRun(t, root, base)
	if want := "FILE GREW       15 lines, was 0 (max 10)   café.go"; !strings.Contains(report, want) {
		t.Errorf("missing accented-name row %q:\n%s", want, report)
	}
	if strings.Contains(report, `\303`) {
		t.Errorf("the path arrived C-quoted:\n%s", report)
	}
	if code != 1 {
		t.Errorf("exit = %d, want 1:\n%s", code, report)
	}
}

// SR50: the growth mode reads the accept list through the same loader, so a
// present-but-unreadable file is loud here too — never a silently empty grant list.
func TestGrowthAcceptUnreadableIsLoud(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses file permissions; cannot simulate an unreadable file")
	}
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "ok.go", lines(1))
	commit(t, root, "base")
	base := headSha(t, root)

	write(t, root, "ok.go", lines(2))
	commit(t, root, "edit")

	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "structure-accept"), []byte("ok.go fine\n"), 0o000); err != nil {
		t.Fatal(err)
	}

	report, code := growthRun(t, root, base)
	if !strings.Contains(report, "structure-accept: present but unreadable") {
		t.Errorf("missing the loud named line:\n%s", report)
	}
	if code != 1 {
		t.Errorf("exit = %d, want 1:\n%s", code, report)
	}
}

// SR51: a base that does not resolve is loud. Command writes the stderr line and exits 1;
// the returned error carries the git failure, and gitOpError renders the `git diff failed:`
// shape Command prints.
func TestGrowthUnresolvableBaseIsLoud(t *testing.T) {
	root := initRepo(t)
	write(t, root, "ok.go", lines(1))
	commit(t, root, "base")

	_, _, err := Growth(root, "no-such-base")
	if err == nil {
		t.Fatal("an unresolvable base returned no error")
	}
	if line := gitOpError("diff", err); !strings.HasPrefix(line, "git diff failed:") {
		t.Errorf("stderr line = %q, want the `git diff failed:` shape", line)
	}

	report, code := growthRun(t, root, "no-such-base")
	if report != "" || code != 1 {
		t.Errorf("unresolvable base: report %q exit %d, want an empty report at exit 1", report, code)
	}
}

// SR62: the mode reads the tip from the working tree, so it fires in the shape the fast
// lane runs it — a detached checkout whose HEAD is the base and whose working tree holds
// the composed tree. A query between two commits sees nothing there.
func TestGrowthReadsTheWorkingTreeNotHEAD(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "10")
	root := initRepo(t)
	write(t, root, "big.go", lines(12))
	commit(t, root, "base")
	base := headSha(t, root)

	// The composed tree is an index-only object: the file grows, `write-tree` names the
	// tree, and no commit carries it. That is what the lane hands its private checkout.
	write(t, root, "big.go", lines(15))
	add(t, root)
	tree := run(t, root, "write-tree")

	checkout := filepath.Join(t.TempDir(), "lane")
	run(t, root, "worktree", "add", "--detach", checkout, base)
	run(t, checkout, "read-tree", "--reset", "-u", tree)

	if head := run(t, checkout, "rev-parse", "HEAD"); head != base {
		t.Fatalf("checkout HEAD = %q, want the base %q", head, base)
	}

	report, code := growthRun(t, checkout, base)
	if want := "FILE GREW       15 lines, was 12 (max 10)   big.go"; !strings.Contains(report, want) {
		t.Errorf("missing growth row %q:\n%s", want, report)
	}
	if code != 1 {
		t.Errorf("exit = %d, want 1:\n%s", code, report)
	}
	if head := run(t, checkout, "rev-parse", "HEAD"); head != base {
		t.Errorf("HEAD moved to %q; the mode must read the working tree, not a commit", head)
	}
}

// SR62: an empty `--growth` value names no base. The grammar refuses it, so the mode
// never resolves it to a self-comparison that passes and prints a dangling base.
func TestGrowthEmptyBaseRefuses(t *testing.T) {
	root := initRepo(t)
	write(t, root, "ok.go", lines(1))
	commit(t, root, "base")
	t.Chdir(root)

	report, code := Command([]string{"--growth", ""})
	if want := `usage: bench structure (unknown argument: --growth "")` + "\n"; report != want {
		t.Errorf("report = %q, want %q", report, want)
	}
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

// SR63: `--since` and `--growth` name two different scopes, so a call carrying both is a
// mistyped invocation. It refuses with the usage line and one reason rather than running
// one flag and dropping the other in silence.
func TestSinceWithGrowthRefuses(t *testing.T) {
	root := initRepo(t)
	write(t, root, "ok.go", lines(1))
	commit(t, root, "base")
	base := headSha(t, root)
	t.Chdir(root)

	report, code := Command([]string{"--since", base, "--growth", base})
	if want := grammar.Help + " (--since and --growth exclude each other)" + "\n"; report != want {
		t.Errorf("report = %q, want %q", report, want)
	}
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}
