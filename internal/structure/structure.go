// Package structure ports `bench structure`. It is the FILE TOO LONG / DIR
// CROWDED structural-debt check over the tracked source tree, plus the
// violation count `bench status` reads. One engine (Check) drives the human
// report, the status count, and the `--since <base>` touched-scope run. The
// length/crowding rules, the per-path budget overrides, and the two caps have
// a single definition — the two-derivations bug class this slice ends.
//
// The shell sent its debt summary and malformed-budget warnings to stderr. This
// command layer folds every line (violations, warnings, the debt summary, the
// ok / no-source messages) into one stdout report. This is behavior-preserving
// because every CLI contract captures the command with 2>&1.
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
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:   "bench structure",
	Help:  "usage: bench structure [--since <base>]",
	Flags: []usage.Flag{{Name: "--since", HasValue: true}},
}

// sourceRe matches a tracked source file by its trailing extension — the exact
// `py|ts|…|sh` set the shell grep filtered on. RE2 tries every alternation branch,
// so a longer extension (tsx, cpp) still anchors to `$` where a prefix would not.
var sourceRe = regexp.MustCompile(`\.(py|ts|tsx|js|jsx|go|rs|java|rb|kt|scala|cs|cpp|cc|c|h|hpp|sh)$`)

// Check runs the structural-debt length/crowding rules over the given file list (the raw
// split of a git query, filtered to source extensions here). Querying git and propagating
// its error is the caller's job — checkAll for the whole tracked tree, Touched for the
// touched scope. So the rules have one home, and each git query owns its own error. It
// returns the full human report the CLI prints and the violation count. Exit-code mapping
// is the caller's (Command's) job.
//
// allMode is set for the whole-tree scan and cleared for the touched (--since) scope. Both
// honor the .bench/structure-accept exclusion (an accepted over-budget file drops out of
// the count). Only the whole-tree scan can judge the accept list's completeness. So the
// accepted: section, the malformed-row warnings, and the stale-row warnings are all-mode only.
func Check(root string, rawFiles []string, allMode bool) (report string, violations int) {
	report, violations, _ = evaluate(root, rawFiles, allMode)
	return report, violations
}

func evaluate(root string, rawFiles []string, allMode bool) (report string, violations int, facts []Fact) {
	maxLines := envInt("BENCH_MAX_LINES", 400)
	maxFiles := envInt("BENCH_MAX_DIR_FILES", 12)
	budgets, budgetWarnings := loadBudgets(filepath.Join(root, ".bench", "structure.budgets"))
	accepts, acceptWarnings, acceptErr := loadAccepts(filepath.Join(root, ".bench", "structure-accept"))

	lines := append([]string(nil), budgetWarnings...)

	// Fail-closed read error (FT29): a present-but-unreadable accept file is loud. A named
	// line and a forced non-zero count return through the same path both the report and
	// ViolationCount read. No surface can silently observe an empty accept list. This is
	// the one intentional exception to "non-zero only on real violations"; the ordinary
	// accept states (absent, malformed, stale) never change the exit code.
	if acceptErr != nil {
		lines = append(lines, "structure-accept: present but unreadable: "+acceptErr.Error())
		return strings.Join(lines, "\n") + "\n", 1, nil
	}

	files := filterSources(rawFiles)

	if len(files) == 0 {
		lines = append(lines, "structure: no tracked source files to check")
		return strings.Join(lines, "\n") + "\n", 0, nil
	}

	if allMode {
		lines = append(lines, acceptWarnings...)
		lines = append(lines, staleAcceptWarnings(accepts, files)...)
	}

	var acceptedLines []string

	// FILE TOO LONG: a file whose newline count exceeds its cap (its exact-path budget
	// override, else the global BENCH_MAX_LINES). Missing and non-regular paths are
	// skipped. An accepted subject is excluded from the count and recorded for the
	// accepted: section.
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
			fact := Fact{Kind: "file", Path: f, Actual: n, Limit: limit, State: "violation"}
			if reason, ok := accepts[f]; ok {
				acceptedLines = append(acceptedLines, fmt.Sprintf("accepted: %s — %s", f, reason))
				fact.State, fact.Detail = "accepted", reason
				facts = append(facts, fact)
				continue
			}
			fact.Detail = fmt.Sprintf("FILE TOO LONG   %d lines (max %d)   %s", n, limit, f)
			facts = append(facts, fact)
			lines = append(lines, fact.Detail)
			violations++
		}
	}

	// DIR CROWDED: a directory whose source-file count exceeds its cap (its `<dir>/`
	// budget override, else the global BENCH_MAX_DIR_FILES). Directories are the path up
	// to the last '/', or '.' when none, counted over the same source list in sorted order.
	// An accepted `<dir>/` key is likewise excluded and recorded, reusing budgetFor's keying.
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
		key := dirs[i] + "/"
		limit := budgetFor(budgets, key, maxFiles)
		if count > limit {
			fact := Fact{Kind: "directory", Path: key, Actual: count, Limit: limit, State: "violation"}
			if reason, ok := accepts[key]; ok {
				acceptedLines = append(acceptedLines, fmt.Sprintf("accepted: %s — %s", key, reason))
				fact.State, fact.Detail = "accepted", reason
				facts = append(facts, fact)
			} else {
				fact.Detail = fmt.Sprintf("DIR CROWDED     %d source files (max %d), group into modules   %s", count, limit, key)
				facts = append(facts, fact)
				lines = append(lines, fact.Detail)
				violations++
			}
		}
		i = j
	}

	if violations > 0 {
		lines = append(lines, fmt.Sprintf("structural debt: %d issue(s). Split along responsibility (see the craft-seams skill); don't fragment to beat the number.", violations))
	} else {
		lines = append(lines, fmt.Sprintf("structure ok (≤%d lines/file, ≤%d source files/dir)", maxLines, maxFiles))
	}

	// The accepted: section (all-mode only) prints one reasoned line per suppressed
	// violation plus a count, in a section separate from the live violations above.
	if allMode && len(acceptedLines) > 0 {
		lines = append(lines, acceptedLines...)
		lines = append(lines, fmt.Sprintf("accepted: %d file(s) over budget by reviewer grant (see .bench/structure-accept)", len(acceptedLines)))
	}

	sort.Slice(facts, func(i, j int) bool { return facts[i].Path < facts[j].Path })
	return strings.Join(lines, "\n") + "\n", violations, facts
}

// ViolationCount is the count `bench status` reads: the violation count of an all-files
// check through the same engine. So the figure cannot drift from the report.
func ViolationCount(root string) int {
	_, violations, err := checkAll(root)
	if err != nil {
		// EXPLICIT tolerate (audit #1, read side): `bench status` is an ambient advisory
		// board the SessionStart hook consumes. A git-query failure degrades this count
		// to zero rather than crashing the hook. `bench structure` is the loud-error path
		// for the same query.
		return 0
	}
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
// budget). A non-digit budget yields a warning and is dropped, so the global cap
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

// loadAccepts reads <root>/.bench/structure-accept into a path→reason map plus warnings for
// malformed rows. It mirrors loadBudgets' comment/whitespace discipline: strip from the first
// '#'; a blank remainder is skipped; the last line may lack a trailing newline. A path
// listed twice keeps its first reason. Three deliberate differences apply.
//
//   - The reason is the whole remainder after the first whitespace-delimited path token, not
//     a Fields split, so a one-clause reason keeps its internal spaces.
//   - A row with a path but no reason is malformed: warned and not honored. A reason is the
//     price of acceptance.
//   - The read-error posture is fail-closed: a missing file is an empty list with no error.
//     Any other read error (present but unreadable) is returned so Check is loud (FT29),
//     never a silent empty list at exit 0.
func loadAccepts(path string) (map[string]string, []string, error) {
	accepts := map[string]string{}
	var warnings []string
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return accepts, nil, nil
		}
		return nil, nil, err
	}
	for _, line := range strings.Split(string(content), "\n") {
		stripped := line
		if i := strings.IndexByte(line, '#'); i >= 0 {
			stripped = line[:i]
		}
		trimmed := strings.TrimLeft(stripped, " \t")
		if trimmed == "" {
			continue
		}
		subj, reason := trimmed, ""
		if i := strings.IndexAny(trimmed, " \t"); i >= 0 {
			subj, reason = trimmed[:i], strings.TrimSpace(trimmed[i+1:])
		}
		if reason == "" {
			warnings = append(warnings, "structure-accept: ignoring malformed line (no reason): "+stripped)
			continue
		}
		if _, seen := accepts[subj]; seen {
			continue
		}
		accepts[subj] = reason
	}
	return accepts, warnings, nil
}

// staleAcceptWarnings reports each accept row whose subject is not a known scanned subject:
// neither a scanned source file nor a scanned `<dir>/` key. So the accept list cannot
// quietly accumulate dead entries. A subject present but under budget suppresses nothing,
// yet is not stale. Warning on that inert case is a separate honesty check left out of
// scope. The result is sorted for a deterministic report, since the map's own iteration
// order is random.
func staleAcceptWarnings(accepts map[string]string, files []string) []string {
	known := make(map[string]bool, len(files)*2)
	for _, f := range files {
		known[f] = true
		dir := "."
		if i := strings.LastIndex(f, "/"); i >= 0 {
			dir = f[:i]
		}
		known[dir+"/"] = true
	}
	var warnings []string
	for subj := range accepts {
		if !known[subj] {
			warnings = append(warnings, "structure-accept: stale accept row (not a scanned source file): "+subj)
		}
	}
	sort.Strings(warnings)
	return warnings
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
	parsed, line, code := usage.Parse(grammar, args)
	if line != "" {
		return line + "\n", code
	}
	base, touched := parsed.Flags["--since"]
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if touched {
		report, violations, terr := Touched(root, base)
		if terr != nil {
			fmt.Fprintln(os.Stderr, gitOpError("diff", terr))
			return "", 1
		}
		return report, exitOf(violations)
	}
	report, violations, cerr := checkAll(root)
	if cerr != nil {
		fmt.Fprintln(os.Stderr, gitOpError("ls-files", cerr))
		return "", 1
	}
	return report, exitOf(violations)
}

// gitOpError builds the stderr line for a failed git query: `git <op> failed: <err>`. It
// is a pure helper. A package unit test pins the exact message shape, while the
// os.Stderr write in Command stays the thin process-boundary rim. Map-dispatched
// commands return only stdout, so a command reaching stderr writes it directly (the
// resolveModel pattern).
func gitOpError(op string, err error) string {
	return fmt.Sprintf("git %s failed: %v", op, err)
}

// Touched runs the length/crowding rules over only the files changed between base and HEAD
// (`git diff --diff-filter=ACMR base..HEAD`). It returns the report, the violation count,
// and the git error if the diff query itself failed. It is the one diff call site. It is
// also the one source of the touched-scope query. The `--since` subcommand and the shift
// loop's refactor gate both read it. Command propagates the error (loud stderr + exit 1);
// the shift gate tolerates it because its own gate run is that worktree's loud oracle.
func Touched(root, base string) (report string, violations int, err error) {
	out, err := git.Output("-C", root, "diff", "--name-only", "--diff-filter=ACMR", base+"..HEAD")
	if err != nil {
		return "", 0, err
	}
	report, violations = Check(root, strings.Split(out, "\n"), false)
	return report, violations, nil
}

func exitOf(violations int) int {
	if violations > 0 {
		return 1
	}
	return 0
}
