package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/poolkey"
)

// TestRawReadAdvice is ER41 and ER43 through the follow-on guard entry. An eligible raw read
// exits 0 with the native read form on stderr: a whole-file read for cat and tail, and a
// line range for head and a simple sed line print. Stdin reads, examples, transformations,
// other sed programs, unresolved and controlled operands, and unknown parses exit 0 with no
// advice. A read whose file resolves outside the current repository, such as a pool path read
// from the primary root, gets no advice either. A pool denial keeps its exit and its exact
// diagnostic. An owned read records under its verified assignment, an owned read outside the
// repository and a primary read record nothing, and a pool-path read keeps its pool attribution.
func TestRawReadAdvice(t *testing.T) {
	home, repo, creation := censusAssignment(t, "raw-read-advice")
	target := filepath.Join(poolkey.Pools(home), "repository", "assignment")
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
	checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{{name: "owned read outside the repository", in: envelope("cat /etc/hosts")}})
	if got := census.HeadBreakdown(home, repo, creation.Assignment.ID); got != "cat=1" {
		t.Fatalf("owned outside read breakdown = %q, want no owned record, cat=1", got)
	}

	inside := filepath.Join(repo, "README.md")
	pooled := filepath.Join(creation.Path, "README.md")
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
		{command: "xargs cat notes.md"},
		{command: "sed -i 3p notes.md"},
		{command: "sed -e 3p notes.md"},
		{command: "cat notes\u009b.md"},
		{command: "cat notes\u0085.md"},
		{command: "cat " + inside, message: file(inside)},
		{command: "cat /etc/hosts"},
		{command: "cat ../outside"},
		{command: "cat " + pooled},
		{command: `echo "cat notes.md"`},
		{command: "cat 'notes.md"},
		{command: "X=" + target + " cat notes.md", want: 2, message: benchguard.PoolReferenceMessage(target)},
		{command: "cd " + target + " && cat notes.md", want: 2, message: benchguard.PoolReferenceMessage(target)},
	}
	rows := make([]guardEnvelopeCase, 0, len(cases))
	for _, tc := range cases {
		rows = append(rows, guardEnvelopeCase{name: tc.command, in: envelope(tc.command), want: tc.want, message: tc.message})
	}
	t.Chdir(repo)
	checkGuardEnvelope(t, guardBenchFollowOn, rows)
	if got := census.HeadBreakdown(home, repo, creation.Assignment.ID); got != "cat=2" {
		t.Fatalf("pool-path read breakdown = %q, want the owned read and the pool attribution, cat=2", got)
	}

	sub := filepath.Join(repo, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(sub)
	checkGuardEnvelope(t, guardBenchFollowOn, []guardEnvelopeCase{
		{name: "relative read from a subdirectory", in: envelope("cat notes.md"), message: file("notes.md")},
		{name: "relative climb inside the repository", in: envelope("cat ../README.md"), message: file("../README.md")},
	})
}
