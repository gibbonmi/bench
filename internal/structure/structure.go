// Package structure ports `bench structure`: the FILE TOO LONG / DIR CROWDED
// structural-debt check over the tracked source tree, plus the violation count
// `bench status` reads. One engine (Check) drives the human report, the status
// count, and the `--since <base>` touched-scope run, so the length/crowding rules,
// the per-path budget overrides, and the two caps have a single definition — the
// two-derivations bug class this slice ends.
//
// The shell sent its debt summary and malformed-budget warnings to stderr; this
// command layer folds every line (violations, warnings, the debt summary, the ok /
// no-source messages) into one stdout report, which is behavior-preserving because
// every CLI contract captures the command with 2>&1.
package structure

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// sourceRe matches a tracked source file by its trailing extension — the exact
// `py|ts|…|sh` set the shell grep filtered on. RE2 tries every alternation branch,
// so a longer extension (tsx, cpp) still anchors to `$` where a prefix would not.
var sourceRe = regexp.MustCompile(`\.(py|ts|tsx|js|jsx|go|rs|java|rb|kt|scala|cs|cpp|cc|c|h|hpp|sh)$`)

// Check runs the structural-debt check under root. mode "all" scans the whole tree
// via `git ls-files`; any other mode uses the passed scopedFiles list (the touched
// scope). It returns the full human report the CLI prints and the violation count;
// exit-code mapping is the caller's (Command's) job.
func Check(root, mode string, scopedFiles []string) (report string, violations int) {
	maxLines := envInt("BENCH_MAX_LINES", 400)
	maxFiles := envInt("BENCH_MAX_DIR_FILES", 12)
	budgets, warnings := loadBudgets(filepath.Join(root, ".bench", "structure.budgets"))

	var files []string
	if mode == "all" {
		out, _ := git.Output("-C", root, "ls-files")
		files = filterSources(strings.Split(out, "\n"))
	} else {
		files = filterSources(scopedFiles)
	}

	lines := append([]string(nil), warnings...)
	if len(files) == 0 {
		lines = append(lines, "structure: no tracked source files to check")
		return strings.Join(lines, "\n") + "\n", 0
	}

	// FILE TOO LONG: a file whose newline count exceeds its cap (its exact-path budget
	// override, else the global BENCH_MAX_LINES). Missing/non-regular paths are skipped.
	for _, f := range files {
		info, err := os.Stat(filepath.Join(root, f))
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		n := 0
		if content, err := os.ReadFile(filepath.Join(root, f)); err == nil {
			n = bytes.Count(content, []byte{'\n'})
		}
		limit := budgetFor(budgets, f, maxLines)
		if n > limit {
			lines = append(lines, fmt.Sprintf("FILE TOO LONG   %d lines (max %d)   %s", n, limit, f))
			violations++
		}
	}

	// DIR CROWDED: a directory whose source-file count exceeds its cap (its `<dir>/`
	// budget override, else the global BENCH_MAX_DIR_FILES). Directories are the path up
	// to the last '/', or '.' when none, counted over the same source list in sorted order.
	dirs := make([]string, 0, len(files))
	for _, f := range files {
		dir := "."
		if i := strings.LastIndex(f, "/"); i >= 0 {
			dir = f[:i]
		}
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	for i := 0; i < len(dirs); {
		j := i
		for j < len(dirs) && dirs[j] == dirs[i] {
			j++
		}
		count := j - i
		limit := budgetFor(budgets, dirs[i]+"/", maxFiles)
		if count > limit {
			lines = append(lines, fmt.Sprintf("DIR CROWDED     %d source files (max %d), group into modules   %s/", count, limit, dirs[i]))
			violations++
		}
		i = j
	}

	if violations > 0 {
		lines = append(lines, fmt.Sprintf("structural debt: %d issue(s). Split along responsibility (see the craft-seams skill); don't fragment to beat the number.", violations))
		return strings.Join(lines, "\n") + "\n", violations
	}
	lines = append(lines, fmt.Sprintf("structure ok (≤%d lines/file, ≤%d source files/dir)", maxLines, maxFiles))
	return strings.Join(lines, "\n") + "\n", 0
}

// ViolationCount is the count `bench status` reads: the second return of an "all"
// Check, computed through the same engine so the figure cannot drift from the report.
func ViolationCount(root string) int {
	_, violations := Check(root, "all", nil)
	return violations
}

// filterSources keeps the entries whose trailing extension is a source extension,
// dropping the empty splits `git ls-files` / a diff list leave behind.
func filterSources(files []string) []string {
	var out []string
	for _, f := range files {
		if f != "" && sourceRe.MatchString(f) {
			out = append(out, f)
		}
	}
	return out
}

// loadBudgets reads <root>/.bench/structure.budgets into a path→budget map plus the
// warnings for malformed lines. Each line is stripped from its first '#'; a blank or
// whitespace-only remainder is skipped; the rest splits on whitespace into (path,
// budget); a non-digit budget yields a warning and is dropped, so the global cap
// stands. A missing file is not an error. The last line may lack a trailing newline.
// A path listed twice keeps its first budget, matching the shell's `awk … {exit}`.
func loadBudgets(path string) (map[string]int, []string) {
	budgets := map[string]int{}
	var warnings []string
	content, err := os.ReadFile(path)
	if err != nil {
		return budgets, warnings
	}
	for _, line := range strings.Split(string(content), "\n") {
		stripped := line
		if i := strings.IndexByte(line, '#'); i >= 0 {
			stripped = line[:i]
		}
		fields := strings.Fields(stripped)
		if len(fields) == 0 {
			continue
		}
		if len(fields) < 2 || !allDigits(fields[1]) {
			warnings = append(warnings, "structure.budgets: ignoring malformed line: "+stripped)
			continue
		}
		if _, seen := budgets[fields[0]]; seen {
			continue
		}
		n, _ := strconv.Atoi(fields[1])
		budgets[fields[0]] = n
	}
	return budgets, warnings
}

// budgetFor returns the exact-key override or the fallback cap.
func budgetFor(budgets map[string]int, key string, fallback int) int {
	if b, ok := budgets[key]; ok {
		return b
	}
	return fallback
}

// envInt reads an integer environment override, falling back to def when the value
// is unset, empty, or not a valid integer.
func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// allDigits reports whether s is a non-empty run of ASCII digits (the `^[0-9]+$` test).
func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// Command implements `bench structure`. No args → an "all" scan of the tree.
// `--since <base>` → scan only the files a diff of base..HEAD touched. Unknown
// argument → usage on stdout, exit 2; outside a repo → structured error, exit 1.
// The exit code is 1 when the report carries any violation, else 0.
func Command(args []string) (string, int) {
	touched := false
	switch {
	case len(args) == 0:
	case args[0] == "-h" || args[0] == "--help":
		return "usage: bench structure [--since <base>]\n", 0
	case args[0] == "--since":
		if len(args) < 2 {
			return toon.Usage("bench structure", "--since") + "\n", 2
		}
		touched = true
	default:
		return toon.Usage("bench structure", args[0]) + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if touched {
		out, _ := git.Output("-C", root, "diff", "--name-only", "--diff-filter=ACMR", args[1]+"..HEAD")
		report, violations := Check(root, "touched", strings.Split(out, "\n"))
		return report, exitOf(violations)
	}
	report, violations := Check(root, "all", nil)
	return report, exitOf(violations)
}

func exitOf(violations int) int {
	if violations > 0 {
		return 1
	}
	return 0
}
