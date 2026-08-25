package worktree

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// summaryFor renders the resume summary for a listing of orphans alone, which is the
// surface these tests are about.
// The counted header is fixed for all of them and is covered where the sweep produces it.
func summaryFor(orphans []OrphanCandidate) string {
	return renderResumeSummary(ResumeResult{Orphans: orphans})
}

func summaryLines(summary string) []string {
	return strings.Split(strings.TrimSuffix(summary, "\n"), "\n")
}

// oneShellArgument asks a real shell what the quoted text parses to.
// The quoting rule is graded by the thing that will run it, not by a second copy of the rule.
func oneShellArgument(t *testing.T, quoted string) (int, string) {
	t.Helper()
	out, err := descendant(t, "sh", "-c", "set -- "+quoted+`; printf '%d\n%s' "$#" "$1"`).Output()
	mustNoError(t, err)
	count, rest, _ := strings.Cut(string(out), "\n")
	var argc int
	if _, scanErr := fmt.Sscanf(count, "%d", &argc); scanErr != nil {
		t.Fatalf("shell reported no argument count for %s: %q", quoted, out)
	}
	return argc, rest
}

// TestResumeSummaryQuotesHostilePaths grades the emitted argument twice: against the
// POSIX single-quoting a reader would predict, and against a shell that parses it.
// A pool path is operator-influenced.
// An unquoted one splits on a space or expands on a glob, and the pasted command names a
// different tree, or none.
func TestResumeSummaryQuotesHostilePaths(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"space":        "/pool/a b/wt",
		"glob":         "/pool/wt*?[abc]",
		"quote":        "/pool/it's here",
		"substitution": "/pool/$(touch pwned)`id`",
		"semicolon":    "/pool/wt; rm -rf /",
	}
	quoting := map[string]string{
		"space":        `'/pool/a b/wt'`,
		"glob":         `'/pool/wt*?[abc]'`,
		"quote":        `'/pool/it'\''s here'`,
		"substitution": "'/pool/$(touch pwned)`id`'",
		"semicolon":    `'/pool/wt; rm -rf /'`,
	}
	for name, path := range cases {
		t.Run(name, func(t *testing.T) {
			want := "bench worktree clean " + quoting[name]
			summary := summaryFor([]OrphanCandidate{{ID: "a1", Path: path}})
			requireTest(t, strings.Contains(summary, want), "summary does not carry %q:\n%s", want, summary)
			argc, arg := oneShellArgument(t, quoting[name])
			requireTest(t, argc == 1 && arg == path, "shell parsed %s into %d arguments, first %q", quoting[name], argc, arg)
		})
	}
}

// TestResumeSummaryPreservesLineStructure holds the summary's line count to its record
// count. renderResumeSummary writes raw lines and escapes nothing, so quoting alone is
// not enough.
// Single-quoting makes a newline literal but still emits the byte.
// One embedded newline would forge a whole extra summary line for a reader.
// A raw ESC byte would drive the terminal that prints it.
func TestResumeSummaryPreservesLineStructure(t *testing.T) {
	t.Parallel()
	for name, path := range map[string]string{
		"newline": "/pool/wt\nreconciled 999; failed 0",
		"escape":  "/pool/wt\x1b[2J\x1b[H",
		"return":  "/pool/wt\rmasked",
	} {
		t.Run(name, func(t *testing.T) {
			summary := summaryFor([]OrphanCandidate{{ID: "a1", Path: path}, {ID: "a2", Path: "/pool/plain"}})
			lines := summaryLines(summary)
			requireTest(t, len(lines) == 3, "summary emitted %d lines for a header and 2 orphans:\n%q", len(lines), summary)
			body := strings.ReplaceAll(summary, "\n", "")
			requireTest(t, !strings.ContainsFunc(body, unicode.IsControl),
				"summary carries a raw control byte outside its line terminators:\n%q", summary)
			requireTest(t, strings.Contains(summary, "a1") && strings.Contains(lines[2], "/pool/plain"),
				"the unsafe path's record went unreported or displaced its neighbour:\n%q", summary)
		})
	}
}

// TestResumeSummaryFallbackNamesTheLedgerRow grades the control-byte fallback's wording,
// not only its line count. The line stands where a pasteable command cannot.
// It hands over the one handle that still resolves: the assignment ID, which is the key
// `bench worktree list` reports.
// The fallback must not send the reader after a path, because no route emits one.
func TestResumeSummaryFallbackNamesTheLedgerRow(t *testing.T) {
	t.Parallel()
	line := summaryLines(summaryFor([]OrphanCandidate{{ID: "a1", Path: "/pool/wt\nforged"}}))[1]
	requireTest(t, strings.Contains(line, "orphan a1:") && strings.Contains(line, "id row in bench worktree list"),
		"fallback line does not point at the ledger row `bench worktree list` reports: %q", line)
	requireTest(t, !strings.Contains(line, "by path"),
		"fallback line sends the reader after a path no route emits: %q", line)
}

// TestResumeSummaryCapsListings pins the bound and its honesty together.
// A cap that does not state what it withholds would read as "that is all of them".
// That is the one way this output could mislead rather than help.
func TestResumeSummaryCapsListings(t *testing.T) {
	t.Parallel()
	var orphans []OrphanCandidate
	for i := range 5 {
		orphans = append(orphans, OrphanCandidate{ID: fmt.Sprintf("o%d", i), Path: fmt.Sprintf("/pool/o%d", i)})
	}

	lines := summaryLines(summaryFor(orphans))
	requireTest(t, len(lines) == 5, "capped summary emitted %d lines, want a header and 3+1 orphan:\n%s", len(lines), strings.Join(lines, "\n"))
	requireTest(t, lines[4] == "and 2 more (5 total)", "orphan cap line = %q", lines[4])
	requireTest(t, strings.Contains(lines[3], "/pool/o2") && !strings.Contains(strings.Join(lines, "\n"), "/pool/o3"),
		"the cap did not bite at the third orphan:\n%s", strings.Join(lines, "\n"))

	atCap := summaryLines(summaryFor(orphans[:3]))
	requireTest(t, len(atCap) == 4, "summary exactly at the cap emitted %d lines, want 4:\n%s", len(atCap), strings.Join(atCap, "\n"))
	requireTest(t, !strings.Contains(strings.Join(atCap, "\n"), " more"),
		"a listing exactly at the cap claims withheld records:\n%s", strings.Join(atCap, "\n"))
}

// resumeReclaimableCount reads the count back out of the rendered resume summary, which
// is the only number an operator ever sees. Reading the struct field instead would grade
// the count against itself rather than against what was reported.
func resumeReclaimableCount(t *testing.T, summary string) int {
	t.Helper()
	match := regexp.MustCompile(`pool: (\d+) reclaimable keys`).FindStringSubmatch(summary)
	if match == nil {
		return 0
	}
	count, err := strconv.Atoi(match[1])
	mustNoError(t, err)
	return count
}

// planReclaimableCount reads the reclaimable column out of the plan's aggregate row.
// The comparison is against what `bench worktree reclaim` actually advertises as its
// target count, not against a second call of the predicate from the test.
func planReclaimableCount(t *testing.T, out string) int {
	t.Helper()
	match := regexp.MustCompile(`(?m)^\s+\d+,(\d+),\d+,"?[0-9a-f]{64}"?$`).FindStringSubmatch(out)
	requireTest(t, match != nil, "plan carries no aggregate row: %q", out)
	count, err := strconv.Atoi(match[1])
	mustNoError(t, err)
	return count
}

func mustResumeClean(t *testing.T, root, home string) (string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := ResumeCleanCommand(root, home, nil, &stdout, &stderr)
	requireTest(t, code == 0, "resume-clean code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	return stdout.String(), code
}

// plantHostilePool fills a pool with the shapes the predicate must separate, plus the
// current repository's own empty key, which is excluded before the predicate runs.
func plantHostilePool(t *testing.T, pool, root string) {
	t.Helper()
	plantDeadChild(t, pool, "dead-key", "wt")
	plantDeadChild(t, pool, "another-dead-key", "wt")
	mustMkdirAll(t, filepath.Join(pool, "empty-key"), 0o700)
	plantLiveChild(t, pool, "live-key", "wt")
	mustWrite(t, filepath.Join(pool, "stray-key-file"), []byte("mine\n"), 0o644)
	mustMkdirAll(t, filepath.Join(pool, filepath.Base(Pool(canonicalRoot(root)))), 0o700)
}

// [RS1] The debris is invisible without a line that both counts it and names the verb.
// A resume that reports the count but not the command leaves the operator hunting for a
// verb.
// One that prints the line over a clean pool trains the operator to ignore it.
func TestResumeSummaryNamesTheReclaimCommandOnlyWhenKeysAreReclaimable(t *testing.T) {
	t.Parallel()
	summary := renderResumeSummary(ResumeResult{ReclaimableKeys: 2})
	requireTest(t, strings.Contains(summary, "pool: 2 reclaimable keys") && strings.Contains(summary, "bench worktree reclaim"),
		"resume summary does not count the reclaimable keys and name the verb:\n%s", summary)

	clean := renderResumeSummary(ResumeResult{})
	requireTest(t, !strings.Contains(clean, "reclaim"),
		"resume summary advertises reclamation over a pool with nothing to reclaim:\n%s", clean)
}

// [RS3] The ambient number is the one an operator trusts without re-checking, so it has to
// be the verb's own target count. A second walk of the pool would drift from the predicate
// the command plans with.
// The drift would only show up as a plan that names a different set than resume promised.
func TestResumeReclaimableCountEqualsWhatTheVerbWouldTarget(t *testing.T) {
	t.Parallel()
	t.Run("hostile pool", func(t *testing.T) {
		pool, root, home := newReclaimPool(t)
		plantHostilePool(t, pool, root)

		summary, _ := mustResumeClean(t, root, home)
		plan, code := mustReclaim(t, root, home)
		requireTest(t, code == 0, "reclaim code=%d out=%q", code, plan)
		want := planReclaimableCount(t, plan)
		requireTest(t, want == 3, "the hostile pool plans %d reclaimable keys, want the two dead and the empty one: %q", want, plan)
		requireTest(t, resumeReclaimableCount(t, summary) == want,
			"resume reported %d reclaimable keys, the plan targets %d:\n%s\n%s", resumeReclaimableCount(t, summary), want, summary, plan)
	})
	t.Run("clean pool", func(t *testing.T) {
		pool, root, home := newReclaimPool(t)
		plantLiveChild(t, pool, "live-key", "wt")
		mustMkdirAll(t, filepath.Join(pool, filepath.Base(Pool(canonicalRoot(root)))), 0o700)

		summary, _ := mustResumeClean(t, root, home)
		plan, code := mustReclaim(t, root, home)
		requireTest(t, code == 0, "reclaim code=%d out=%q", code, plan)
		want := planReclaimableCount(t, plan)
		requireTest(t, want == 0, "the clean pool plans %d reclaimable keys, want none: %q", want, plan)
		requireTest(t, resumeReclaimableCount(t, summary) == want,
			"resume reported %d reclaimable keys over a clean pool:\n%s", resumeReclaimableCount(t, summary), summary)
	})
}

// [RS2] Resume runs unattended at session start, outside any tree the gate observes. It
// reports the pool and never touches it.
// The whole recursive listing must survive a resume byte-identical, including the keys it
// just counted as reclaimable.
func TestResumeCleanRemovesNoPoolKey(t *testing.T) {
	t.Parallel()
	pool, root, home := newReclaimPool(t)
	plantHostilePool(t, pool, root)
	before := poolListing(t, pool)

	summary, _ := mustResumeClean(t, root, home)
	requireTest(t, resumeReclaimableCount(t, summary) > 0, "the fixture pool reported nothing reclaimable:\n%s", summary)
	requireTest(t, poolListing(t, pool) == before,
		"the pool changed across a resume:\nbefore\n%s\nafter\n%s", before, poolListing(t, pool))
}

// [P4] A resume that cannot read the pool still succeeds at its own work.
// It must not report a zero it has no basis for.
// That is the one state where the ambient count could disagree with what the verb would say.
func TestResumeReportsAnUnreadablePoolRatherThanZero(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	pool := filepath.Join(home, "worktrees")
	mustMkdirAll(t, pool, 0o700)
	mustNoError(t, os.Chmod(pool, 0o000))
	t.Cleanup(func() { _ = os.Chmod(pool, 0o700) })

	result, err := conservativeCleanupAt(root, home, currentTime())
	requireTest(t, err == nil, "resume failed over an unreadable pool: %v", err)
	requireTest(t, result.PoolUnreadable != nil, "an unreadable pool reported no failure")
	requireTest(t, result.ReclaimableKeys == 0, "an unreadable pool reported %d keys", result.ReclaimableKeys)

	summary := renderResumeSummary(result)
	requireTest(t, strings.Contains(summary, "pool: not read"), "summary %q does not report the unreadable pool", summary)
	requireTest(t, !strings.Contains(summary, "0 reclaimable"), "summary %q reports a zero it cannot support", summary)
}
