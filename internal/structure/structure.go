// Package structure ports `bench structure`. It is the FILE TOO LONG / DIR
// CROWDED structural-debt check over the tracked source tree, plus the
// violation count `bench status` reads. One engine (Check) drives the human
// report, the status count, and the `--since <base>` run; Growth reds a file
// grown past its limit. The length/crowding rules, the per-path budget
// overrides, and the two caps have a single definition — the two-derivations
// bug class this slice ends.
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
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for this subcommand.
// Arity, flag recognition, `--`, and help all come from there rather than a local switch.
var grammar = usage.Grammar{
	Cmd:   "bench structure",
	Help:  "usage: bench structure [--since <base> | --growth <base>]",
	Flags: []usage.Flag{{Name: "--since", HasValue: true}, {Name: "--growth", HasValue: true}},
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
	growthBase, growth := parsed.Flags["--growth"]
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if growth && !touched {
		report, violations, gerr := Growth(root, growthBase)
		if gerr != nil {
			fmt.Fprintln(os.Stderr, gitOpError("diff", gerr))
			return "", 1
		}
		return report, exitOf(violations)
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

// growthChange is one changed file of the growth query: the path it carries at the tip, and
// the path the base knows it by, which differ only for an exact rename. addedAtTip marks a
// path the base holds nothing for, whose base count is therefore zero.
type growthChange struct {
	path, basePath string
	addedAtTip     bool
}

// Growth is the growth ratchet: it reds a source file only when the file both exceeds its
// limit and gained lines since base. Existing debt therefore stays soft and a split is never
// punished, while a file that grows past its cap is named on the spot. The limit and the
// grant list come from the engine the all scan reads, so a `structure.budgets` row and a
// `.bench/structure-accept` row mean here what they mean there, and the accept read keeps the
// fail-closed posture (FT29). It returns the report, the count of grown files, and the git
// error when the diff query failed; Command renders that error as `--since` does.
func Growth(root, base string) (report string, violations int, err error) {
	maxLines := envInt("BENCH_MAX_LINES", 400)
	budgets, budgetWarnings := loadBudgets(filepath.Join(root, ".bench", "structure.budgets"))
	accepts, _, acceptErr := loadAccepts(filepath.Join(root, ".bench", "structure-accept"))

	lines := append([]string(nil), budgetWarnings...)
	if acceptErr != nil {
		lines = append(lines, "structure-accept: present but unreadable: "+acceptErr.Error())
		return strings.Join(lines, "\n") + "\n", 1, nil
	}

	changes, err := growthChanges(root, base)
	if err != nil {
		return "", 0, err
	}

	paths := make([]string, 0, len(changes))
	for _, change := range changes {
		paths = append(paths, change.path)
	}
	sources := make(map[string]bool, len(paths))
	for _, p := range filterSources(paths) {
		sources[p] = true
	}

	// FILE GREW: a source file whose tip count exceeds both its limit and its base count.
	for _, change := range changes {
		if !sources[change.path] {
			continue
		}
		if _, granted := accepts[change.path]; granted {
			continue
		}
		limit := budgetFor(budgets, change.path, maxLines)
		tip := growthTipCount(root, change.path)
		if tip <= limit {
			continue
		}
		was := 0
		if !change.addedAtTip {
			was = growthBaseCount(root, base, change.basePath)
		}
		if tip <= was {
			continue
		}
		lines = append(lines, fmt.Sprintf("FILE GREW       %d lines, was %d (max %d)   %s", tip, was, limit, change.path))
		violations++
	}

	if violations > 0 {
		lines = append(lines, fmt.Sprintf("structure growth: %d file(s) grew past budget. Split along responsibility (see the craft-seams skill), or record a reviewer grant in .bench/structure-accept.", violations))
	} else {
		lines = append(lines, fmt.Sprintf("structure growth ok (no source file grew past its budget since %s)", base))
	}
	return strings.Join(lines, "\n") + "\n", violations, nil
}

// growthChanges lists what base..HEAD changed, with the two framing choices the mode needs.
// The NUL framing is load-bearing: under the default `core.quotepath` a newline-framed name
// with a byte above ASCII arrives C-quoted, so a reader would carry a path no file has and
// the extension filter would drop it. `-M100%` pairs exact renames, so a pure move reads
// its old path's count rather than reading as an addition.
func growthChanges(root, base string) ([]growthChange, error) {
	raw, err := git.Raw("-C", root, "diff", "--name-status", "-z", "-M100%", "--diff-filter=ACMR", base+"..HEAD")
	if err != nil {
		return nil, err
	}
	frames := strings.Split(string(raw), "\x00")
	var changes []growthChange
	for i := 0; i < len(frames); i++ {
		status := frames[i]
		if status == "" {
			continue
		}
		// A rename or copy spends three frames (status, old path, new path) and every other
		// status two. A truncated tail is unreadable as a change, so it ends the scan.
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			if i+2 >= len(frames) {
				break
			}
			changes = append(changes, growthChange{path: frames[i+2], basePath: frames[i+1], addedAtTip: status[0] == 'C'})
			i += 2
			continue
		}
		if i+1 >= len(frames) {
			break
		}
		changes = append(changes, growthChange{path: frames[i+1], basePath: frames[i+1], addedAtTip: status == "A"})
		i++
	}
	return changes, nil
}

// growthTipCount is the working tree's newline count for path, counted exactly as evaluate
// counts a scanned file. A missing or non-regular path counts zero.
func growthTipCount(root, path string) int {
	info, err := os.Stat(filepath.Join(root, path))
	if err != nil || !info.Mode().IsRegular() {
		return 0
	}
	content, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		return 0
	}
	return bytes.Count(content, []byte{'\n'})
}

// growthBaseCount is the newline count of path's blob at base. An unreadable blob counts
// zero, which is the count of a path the base holds nothing for.
func growthBaseCount(root, base, path string) int {
	blob, err := git.Raw("-C", root, "show", base+":"+path)
	if err != nil {
		return 0
	}
	return bytes.Count(blob, []byte{'\n'})
}

func exitOf(violations int) int {
	if violations > 0 {
		return 1
	}
	return 0
}
