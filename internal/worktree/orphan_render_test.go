package worktree

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"unicode"
)

// summaryFor renders the resume summary for a listing of orphans and preserved records
// alone, which is the surface these tests are about; the counted header is fixed for all
// of them and is covered where the sweep produces it.
func summaryFor(orphans []OrphanCandidate, preserved []PreservedOrphan) string {
	return renderResumeSummary(ResumeResult{Orphans: orphans, Preserved: preserved})
}

func summaryLines(summary string) []string {
	return strings.Split(strings.TrimSuffix(summary, "\n"), "\n")
}

// oneShellArgument asks a real shell what the quoted text parses to, so the quoting rule
// is graded by the thing that will run it rather than by a second copy of the rule.
func oneShellArgument(t *testing.T, quoted string) (int, string) {
	t.Helper()
	out, err := exec.Command("sh", "-c", "set -- "+quoted+`; printf '%d\n%s' "$#" "$1"`).Output()
	mustNoError(t, err)
	count, rest, _ := strings.Cut(string(out), "\n")
	var argc int
	if _, scanErr := fmt.Sscanf(count, "%d", &argc); scanErr != nil {
		t.Fatalf("shell reported no argument count for %s: %q", quoted, out)
	}
	return argc, rest
}

// TestResumeSummaryQuotesHostilePaths grades the emitted argument twice: against the
// POSIX single-quoting a reader would predict, and against a shell that parses it. A
// pool path is operator-influenced, so an unquoted one splits on a space or expands on a
// glob and the pasted command names a different tree, or none.
func TestResumeSummaryQuotesHostilePaths(t *testing.T) {
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
			summary := summaryFor([]OrphanCandidate{{ID: "a1", Path: path}}, nil)
			requireTest(t, strings.Contains(summary, want), "summary does not carry %q:\n%s", want, summary)
			argc, arg := oneShellArgument(t, quoting[name])
			requireTest(t, argc == 1 && arg == path, "shell parsed %s into %d arguments, first %q", quoting[name], argc, arg)
		})
	}
}

// TestResumeSummaryPreservesLineStructure holds the summary's line count to its record
// count. renderResumeSummary writes raw lines and escapes nothing, so quoting is not
// enough: single-quoting makes a newline literal but still emits the byte, and one
// embedded newline would forge a whole extra summary line for a reader — and a raw ESC
// would drive the terminal that prints it.
func TestResumeSummaryPreservesLineStructure(t *testing.T) {
	for name, path := range map[string]string{
		"newline": "/pool/wt\nreconciled 999; failed 0",
		"escape":  "/pool/wt\x1b[2J\x1b[H",
		"return":  "/pool/wt\rmasked",
	} {
		t.Run(name, func(t *testing.T) {
			summary := summaryFor([]OrphanCandidate{{ID: "a1", Path: path}, {ID: "a2", Path: "/pool/plain"}}, nil)
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

// TestResumeSummaryCapsListings pins the bound and its honesty together: a cap that did
// not state what it withheld would read as "that is all of them", which is the one way
// bounding this output could mislead rather than help.
func TestResumeSummaryCapsListings(t *testing.T) {
	var orphans []OrphanCandidate
	var preserved []PreservedOrphan
	for i := range 5 {
		orphans = append(orphans, OrphanCandidate{ID: fmt.Sprintf("o%d", i), Path: fmt.Sprintf("/pool/o%d", i)})
		preserved = append(preserved, PreservedOrphan{ID: fmt.Sprintf("p%d", i), Ref: fmt.Sprintf("refs/bench/recovery/o/p/%d", i)})
	}

	lines := summaryLines(summaryFor(orphans, preserved))
	requireTest(t, len(lines) == 9, "capped summary emitted %d lines, want a header, 3+1 orphan and 3+1 preserved:\n%s", len(lines), strings.Join(lines, "\n"))
	requireTest(t, lines[4] == "and 2 more (5 total)", "orphan cap line = %q", lines[4])
	requireTest(t, lines[8] == "and 2 more (5 total)", "preserved cap line = %q", lines[8])
	requireTest(t, strings.Contains(lines[3], "/pool/o2") && !strings.Contains(strings.Join(lines, "\n"), "/pool/o3"),
		"the cap did not bite at the third orphan:\n%s", strings.Join(lines, "\n"))

	atCap := summaryLines(summaryFor(orphans[:3], preserved[:3]))
	requireTest(t, len(atCap) == 7, "summary exactly at the cap emitted %d lines, want 7:\n%s", len(atCap), strings.Join(atCap, "\n"))
	requireTest(t, !strings.Contains(strings.Join(atCap, "\n"), " more"),
		"a listing exactly at the cap claims withheld records:\n%s", strings.Join(atCap, "\n"))
}
