package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/git"
)

// gitOpError pins the exact stderr message shape Command writes on a git-query failure:
// the git-op name plus the underlying error, `git <op> failed: <err>`.
func TestGitOpErrorShape(t *testing.T) {
	got := gitOpError("ls-files", errors.New("index file corrupt"))
	if want := "git ls-files failed: index file corrupt"; got != want {
		t.Errorf("gitOpError = %q, want %q", got, want)
	}
}

// git is subprocess-based, so the seam is a real tree: init a repo, write files,
// `git add`, and exercise Check/ViolationCount/Command against it. Identity is set so
// commits don't error.
func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	run(t, root, "init")
	run(t, root, "config", "user.email", "t@example.com")
	run(t, root, "config", "user.name", "t")
	return root
}

func run(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := git.Output(append([]string{"-C", root}, args...)...)
	if err != nil {
		t.Fatalf("git %v: %v (%s)", args, err, out)
	}
	return out
}

func write(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func add(t *testing.T, root string) { run(t, root, "add", "-A") }

// lines returns a body with exactly n newline bytes (wc -l semantics).
func lines(n int) string { return strings.Repeat("x\n", n) }

// A file over the line budget is flagged with the exact FILE TOO LONG substring; a
// granted higher per-path budget then suppresses it and the tree reports ok.
func TestFileTooLongAndBudget(t *testing.T) {
	root := initRepo(t)
	write(t, root, "long.go", lines(401))
	write(t, root, "ok.go", lines(10))
	add(t, root)

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "FILE TOO LONG   401 lines (max 400)   long.go") {
		t.Errorf("missing FILE TOO LONG line:\n%s", report)
	}
	if v != 1 {
		t.Errorf("violations = %d, want 1", v)
	}

	// A higher budget for the exact path suppresses the violation.
	write(t, root, ".bench/structure.budgets", "long.go 500\n")
	report, v, _ = checkAll(root)
	if strings.Contains(report, "FILE TOO LONG") {
		t.Errorf("granted budget did not suppress:\n%s", report)
	}
	if v != 0 || !strings.Contains(report, "structure ok") {
		t.Errorf("want clean ok report, got v=%d:\n%s", v, report)
	}
}

// A tightening (lower) budget on the UNTERMINATED last line of the budgets file must
// still parse and trigger the violation against the reduced cap.
func TestTighteningBudgetUnterminatedLastLine(t *testing.T) {
	root := initRepo(t)
	write(t, root, "mid.sh", lines(200))
	add(t, root)
	// No trailing newline on the final line; a leading comment line precedes it.
	write(t, root, ".bench/structure.budgets", "# caps\nmid.sh 100")

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "200 lines (max 100)   mid.sh") {
		t.Errorf("tightened cap not applied:\n%s", report)
	}
	if v != 1 {
		t.Errorf("violations = %d, want 1", v)
	}
}

// A malformed budgets line warns and is dropped, so the global cap still governs.
func TestMalformedBudgetWarns(t *testing.T) {
	root := initRepo(t)
	write(t, root, "big.go", lines(401))
	add(t, root)
	write(t, root, ".bench/structure.budgets", "big.go notanumber\n")

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "structure.budgets: ignoring malformed line: big.go notanumber") {
		t.Errorf("missing malformed warning:\n%s", report)
	}
	if !strings.Contains(report, "FILE TOO LONG   401 lines (max 400)   big.go") {
		t.Errorf("global cap did not govern after malformed line:\n%s", report)
	}
	if v != 1 {
		t.Errorf("violations = %d, want 1", v)
	}
}

// A path listed twice keeps its FIRST budget, matching the shell's first-wins awk
// (`$1==p {print $2; exit}`) — a later duplicate line does not override it.
func TestDuplicateBudgetKeyFirstWins(t *testing.T) {
	root := initRepo(t)
	write(t, root, "dup.go", lines(401))
	add(t, root)
	// First budget 100 (tightening — dup.go is over it); a second, looser 500 must
	// NOT win, or the violation would be suppressed.
	write(t, root, ".bench/structure.budgets", "dup.go 100\ndup.go 500\n")

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "401 lines (max 100)   dup.go") {
		t.Errorf("first budget (100) did not win over the later duplicate:\n%s", report)
	}
	if v != 1 {
		t.Errorf("violations = %d, want 1", v)
	}
}

// A non-numeric env cap falls back to the default (400/12); the shell fed garbage
// into arithmetic (yielding cap 0), which this Go port deliberately does not preserve.
func TestNonNumericEnvCapFallsBack(t *testing.T) {
	t.Setenv("BENCH_MAX_LINES", "notanumber")
	root := initRepo(t)
	write(t, root, "under.go", lines(399)) // under the 400 default, over a garbage cap-0
	add(t, root)

	report, v, _ := checkAll(root)
	if v != 0 || !strings.Contains(report, "structure ok") {
		t.Errorf("non-numeric BENCH_MAX_LINES should fall back to 400, got v=%d:\n%s", v, report)
	}
}

// DIR CROWDED fires strictly above the cap: 12 files ok, 13 over at the default; a
// granted directory budget (keyed with a trailing slash) then suppresses it.
func TestDirCrowded(t *testing.T) {
	root := initRepo(t)
	for i := 0; i < 12; i++ {
		write(t, root, fmt.Sprintf("pkg/f%02d.go", i), "x\n")
	}
	add(t, root)
	if _, v, _ := checkAll(root); v != 0 {
		t.Errorf("12 files should be ok, got %d violations", v)
	}

	write(t, root, "pkg/f12.go", "x\n")
	add(t, root)
	report, v, _ := checkAll(root)
	if !strings.Contains(report, "DIR CROWDED     13 source files (max 12), group into modules   pkg/") {
		t.Errorf("missing DIR CROWDED line:\n%s", report)
	}
	if v != 1 {
		t.Errorf("violations = %d, want 1", v)
	}

	write(t, root, ".bench/structure.budgets", "pkg/ 20\n")
	report, v, _ = checkAll(root)
	if strings.Contains(report, "DIR CROWDED") {
		t.Errorf("granted dir budget did not suppress:\n%s", report)
	}
	if v != 0 {
		t.Errorf("violations = %d, want 0", v)
	}
}

// A directory whose name contains a space is kept whole — never split on whitespace.
func TestSpaceDirPreserved(t *testing.T) {
	t.Setenv("BENCH_MAX_DIR_FILES", "1")
	root := initRepo(t)
	write(t, root, "space dir/file1.sh", "x\n")
	write(t, root, "space dir/file2.sh", "x\n")
	add(t, root)

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "group into modules   space dir/") {
		t.Errorf("space dir not preserved whole:\n%s", report)
	}
	if strings.Contains(report, "   ./") || strings.Contains(report, "   dir/") {
		t.Errorf("path was split on whitespace:\n%s", report)
	}
	if v == 0 {
		t.Errorf("expected a crowded violation for the space dir")
	}
}

// ViolationCount is the integer status reads, off the same engine.
func TestViolationCount(t *testing.T) {
	root := initRepo(t)
	write(t, root, "a.go", lines(401))
	write(t, root, "b.go", lines(401))
	write(t, root, "c.go", lines(10))
	add(t, root)
	if got := ViolationCount(root); got != 2 {
		t.Errorf("ViolationCount = %d, want 2", got)
	}
}

// No tracked source files → the definitive message, exit 0.
func TestNoSourceFiles(t *testing.T) {
	root := initRepo(t)
	write(t, root, "README.md", "prose\n")
	add(t, root)
	report, v, _ := checkAll(root)
	if !strings.Contains(report, "structure: no tracked source files to check") || v != 0 {
		t.Errorf("want no-source message, got v=%d:\n%s", v, report)
	}
}

// Command maps violations to exit codes and rejects unknown arguments.
func TestCommand(t *testing.T) {
	root := initRepo(t)
	t.Chdir(root)

	if r, c := Command([]string{"--bogus"}); c != 2 || !strings.Contains(r, "usage:") {
		t.Errorf("unknown arg: report %q exit %d", r, c)
	}

	write(t, root, "ok.go", "x\n")
	add(t, root)
	if r, c := Command(nil); c != 0 || !strings.Contains(r, "structure ok") {
		t.Errorf("clean tree: report %q exit %d", r, c)
	}

	write(t, root, "big.go", lines(401))
	add(t, root)
	if r, c := Command(nil); c != 1 || !strings.Contains(r, "FILE TOO LONG") {
		t.Errorf("violation: report %q exit %d", r, c)
	}
}

// Command --since scopes the check to the files a diff of base..HEAD touched.
func TestCommandSince(t *testing.T) {
	root := initRepo(t)
	write(t, root, "old.go", "x\n")
	add(t, root)
	run(t, root, "commit", "-m", "base")
	base := run(t, root, "rev-parse", "HEAD")

	write(t, root, "new.go", lines(401))
	add(t, root)
	run(t, root, "commit", "-m", "add big")

	t.Chdir(root)
	r, c := Command([]string{"--since", base})
	if c != 1 || !strings.Contains(r, "FILE TOO LONG   401 lines (max 400)   new.go") {
		t.Errorf("since scope: report %q exit %d", r, c)
	}
	if strings.Contains(r, "old.go") {
		t.Errorf("untouched file leaked into --since scope:\n%s", r)
	}
}

// A valid accept row excludes its over-budget file from the violation count and
// records it in the separate accepted: section with its reason and a count, so the
// suppression is per-file, reasoned on the page, and reversible — not a silent bump.
func TestAcceptExcludesAndPrints(t *testing.T) {
	root := initRepo(t)
	write(t, root, "long.go", lines(401))
	write(t, root, "ok.go", lines(10))
	write(t, root, ".bench/structure-accept", "long.go cohesive and barely over budget\n")
	add(t, root)

	report, v, _ := checkAll(root)
	if v != 0 {
		t.Errorf("accepted file left in the count: v=%d\n%s", v, report)
	}
	if strings.Contains(report, "FILE TOO LONG") {
		t.Errorf("accepted file still printed as a live violation:\n%s", report)
	}
	if !strings.Contains(report, "accepted: long.go — cohesive and barely over budget") {
		t.Errorf("missing per-file accepted line:\n%s", report)
	}
	if !strings.Contains(report, "1 file(s)") {
		t.Errorf("missing accepted count:\n%s", report)
	}
	if !strings.Contains(report, "structure ok") {
		t.Errorf("want the ok summary once the only violation is suppressed:\n%s", report)
	}
}

// Absence of the accept file is an empty list at exit 0 — never an error — the
// distinct-from-unreadable half of the FT29 posture. A clean tree stays clean and
// no accept line leaks: an impl that errored on a missing file would fail here.
func TestAcceptAbsentIsEmpty(t *testing.T) {
	root := initRepo(t)
	write(t, root, "ok.go", lines(10))
	add(t, root)

	report, v, _ := checkAll(root)
	if v != 0 || !strings.Contains(report, "structure ok") {
		t.Errorf("absent accept file should be a clean empty list, got v=%d:\n%s", v, report)
	}
	if strings.Contains(report, "structure-accept") || strings.Contains(report, "accepted:") {
		t.Errorf("absent accept file leaked an accept line:\n%s", report)
	}
}

// A row with a path but no reason is malformed: reported and NOT honored, so the
// file it names stays counted — a reason is the price of acceptance.
func TestAcceptMalformedRowNotHonored(t *testing.T) {
	root := initRepo(t)
	write(t, root, "long.go", lines(401))
	write(t, root, ".bench/structure-accept", "long.go\n")
	add(t, root)

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "structure-accept: ignoring malformed line (no reason): long.go") {
		t.Errorf("missing malformed warning:\n%s", report)
	}
	if v != 1 || !strings.Contains(report, "FILE TOO LONG   401 lines (max 400)   long.go") {
		t.Errorf("reasonless row was honored (v=%d) — acceptance without a reason:\n%s", v, report)
	}
}

// A valid accept row whose path is not a scanned source file is reported stale, so
// the list cannot quietly accumulate dead entries. It never touches the live count.
func TestAcceptStaleRowWarns(t *testing.T) {
	root := initRepo(t)
	write(t, root, "long.go", lines(401))
	write(t, root, ".bench/structure-accept", "gone.go it moved away\n")
	add(t, root)

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "structure-accept: stale accept row (not a scanned source file): gone.go") {
		t.Errorf("missing stale warning:\n%s", report)
	}
	if v != 1 {
		t.Errorf("stale row should not change the live count: v=%d\n%s", v, report)
	}
}

// An accept row whose path contains a space is owned deterministically: the first
// whitespace-delimited token is the path, the remainder the reason. A spaced path is
// not a scanned source file, so it is reported stale — never misattributed or panicked.
func TestAcceptWhitespacePathIsStale(t *testing.T) {
	root := initRepo(t)
	write(t, root, "long.go", lines(401))
	write(t, root, ".bench/structure-accept", "space file.go a reason with several words\n")
	add(t, root)

	report, v, _ := checkAll(root)
	if !strings.Contains(report, "structure-accept: stale accept row (not a scanned source file): space") {
		t.Errorf("spaced path not owned deterministically as first-token:\n%s", report)
	}
	if v != 1 {
		t.Errorf("spaced-path row wrongly suppressed a live violation: v=%d\n%s", v, report)
	}
}

// A present-but-unreadable accept file is LOUD — a non-zero result with a named line
// — never a silently empty list at exit 0 (the FT29 false-empty defect). The forced
// non-zero flows through the same count both the report and ViolationCount read.
func TestAcceptUnreadableIsLoud(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses file permissions; cannot simulate an unreadable file")
	}
	root := initRepo(t)
	write(t, root, "ok.go", lines(10)) // clean tree: no live violation to mask the loud one
	add(t, root)
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "structure-accept"), []byte("ok.go fine\n"), 0o000); err != nil {
		t.Fatal(err)
	}

	report, v, _ := checkAll(root)
	if v == 0 {
		t.Errorf("unreadable accept file was swallowed to exit 0 (false-empty):\n%s", report)
	}
	if !strings.Contains(report, "structure-accept: present but unreadable") {
		t.Errorf("missing loud named line:\n%s", report)
	}
	if ViolationCount(root) == 0 {
		t.Errorf("ViolationCount silently observed an empty accept list")
	}
}
