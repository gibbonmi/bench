package structure

// This file holds the limit engine every structure surface reads: the two on-disk
// lists (.bench/structure.budgets and .bench/structure-accept), their parse and
// warning discipline, the per-path cap lookup, and the environment override. The
// report, the status count, the touched scope, the growth mode, and the gate's
// grant validator all take their limits from here, so a cap has one definition.
// It pairs with accept_validation.go, which grades the same accept list.

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

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
