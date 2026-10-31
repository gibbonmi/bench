package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/census"
)

// TestRawReadAdvice is ER41 and ER43 through the follow-on guard entry. An eligible raw read
// exits 0 with the native read form on stderr: a whole-file read for cat and tail, and a
// line range for head and a simple sed line print. Stdin reads, examples, transformations,
// other sed programs, unresolved operands, and unknown parses exit 0 with no advice. A pool
// denial keeps its exit and its exact diagnostic. An owned read records under its verified
// assignment, and a primary read records nothing.
func TestRawReadAdvice(t *testing.T) {
	target := poolTarget(t)
	envelope := func(command string) string {
		quoted, _ := json.Marshal(command)
		return `{"tool_input":{"command":` + string(quoted) + `}}`
	}
	file := func(path string) string { return rawReadAdvice(census.Read{Path: path}) }
	lines := func(path string, first, last int) string {
		return rawReadAdvice(census.Read{Path: path, First: first, Last: last})
	}
	if got, want := lines("notes.md", 3, 7), "bench read notes.md --lines 3:7"; !strings.Contains(got, want) {
		t.Fatalf("range advice = %q, want it to name %q", got, want)
	}
	if got, want := file("a b.md"), "bench read 'a b.md' "; !strings.Contains(got, want) {
		t.Fatalf("file advice = %q, want it to name %q", got, want)
	}
	cases := []struct {
		command string
		want    int
		message string
	}{
		{command: "cat notes.md", message: file("notes.md")},
		{command: "cat 'a b.md'", message: file("a b.md")},
		{command: "timeout 5 cat notes.md", message: file("notes.md")},
		{command: "head notes.md", message: lines("notes.md", 1, 10)},
		{command: "head -n 5 notes.md", message: lines("notes.md", 1, 5)},
		{command: "head -n5 notes.md", message: lines("notes.md", 1, 5)},
		{command: "tail notes.md", message: file("notes.md")},
		{command: "tail -n 5 notes.md", message: file("notes.md")},
		{command: "sed -n '3,7p' notes.md", message: lines("notes.md", 3, 7)},
		{command: "sed -n '4p' notes.md", message: lines("notes.md", 4, 4)},
		{command: "cat"},
		{command: "cat -"},
		{command: "cat <<EOF\nx\nEOF"},
		{command: "cat <<< x"},
		{command: "cat < notes.md"},
		{command: "cat notes.md > out.md"},
		{command: "cat a.md b.md"},
		{command: "cat -n notes.md"},
		{command: "cat -A notes.md"},
		{command: "head -c 5 notes.md"},
		{command: "head -n 0 notes.md"},
		{command: "head -n -2 notes.md"},
		{command: "tail -n +5 notes.md"},
		{command: "sed 's/a/b/' notes.md"},
		{command: "sed -n '/re/p' notes.md"},
		{command: "sed -i -n '1p' notes.md"},
		{command: "sed -n '7,3p' notes.md"},
		{command: "sed -n '1,$p' notes.md"},
		{command: "cat $VAR"},
		{command: `cat "$VAR"`},
		{command: "cat *.md"},
		{command: "cat ~/notes.md"},
		{command: "cat $(echo notes.md)"},
		{command: "cat `echo notes.md`"},
		{command: "cat notes.md | grep x"},
		{command: "cat notes.md; ls"},
		{command: "xargs cat"},
		{command: `echo "cat notes.md"`},
		{command: "cat 'notes.md"},
		{command: "X=" + target + " cat notes.md", want: 2, message: benchguard.PoolReferenceMessage(target)},
		{command: "cd " + target + " && cat notes.md", want: 2, message: benchguard.PoolReferenceMessage(target)},
	}
	rows := make([]guardEnvelopeCase, 0, len(cases))
	for _, tc := range cases {
		rows = append(rows, guardEnvelopeCase{name: tc.command, in: envelope(tc.command), want: tc.want, message: tc.message})
	}
	checkGuardEnvelope(t, guardBenchFollowOn, rows)

	home, repo, creation := censusAssignment(t, "raw-read-advice")
	read := []guardEnvelopeCase{{name: "eligible read", in: envelope("cat README.md"), message: file("README.md")}}
	t.Chdir(repo)
	checkGuardEnvelope(t, guardBenchFollowOn, read)
	if counts, _ := census.Counts(home, repo); len(counts) != 0 {
		t.Fatalf("primary read counts = %v, want none", counts)
	}
	t.Chdir(creation.Path)
	checkGuardEnvelope(t, guardBenchFollowOn, read)
	if got := census.HeadBreakdown(home, repo, creation.Assignment.ID); got != "cat=1" {
		t.Fatalf("owned read breakdown = %q, want cat=1", got)
	}
}
