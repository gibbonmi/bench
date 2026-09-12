package testreport

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/sanitize"
)

func TestFullFailureDiagnostics(t *testing.T) {
	first := "first " + strings.Repeat("x", bounds.PreviewRuneLimit)
	for _, tc := range []struct {
		name      string
		events    []event
		fragments []string
		escapes   []string
		first     string
		want      Outcome
	}{
		{
			name: "test diagnostics",
			events: []event{
				{Action: "run", Package: "z", Test: "TestLater"},
				{Action: "output", Package: "z", Test: "TestLater", Output: first + "\nmiddle\x1b[31m diagnostic\n"},
				{Action: "output", Package: "z", Test: "TestLater", Output: "last diagnostic\n--- FAIL: TestLater (0.00s)\n"},
				{Action: "fail", Package: "z", Test: "TestLater"},
				{Action: "fail", Package: "z"},
				{Action: "run", Package: "a", Test: "TestGroup"},
				{Action: "run", Package: "a", Test: "TestGroup/child"},
				{Action: "output", Package: "a", Test: "TestGroup/child", Output: "child first\nchild center\nchild last\n"},
				{Action: "fail", Package: "a", Test: "TestGroup/child"},
				{Action: "fail", Package: "a", Test: "TestGroup"},
				{Action: "run", Package: "a", Test: "TestEmpty"},
				{Action: "fail", Package: "a", Test: "TestEmpty"},
				{Action: "fail", Package: "a"},
			},
			fragments: []string{"a,TestEmpty,", "no diagnostic emitted", "a,TestGroup/child,", "child first", "child center", "child last", "z,TestLater,", first, "middle", "last diagnostic"},
			escapes:   []string{`child first\\nchild center\\nchild last`, `middle\\u001b[31m diagnostic\\nlast diagnostic`},
			first:     first,
			want:      Outcome{Kind: OutcomeFailed, FailedTests: 3, Ran: 4},
		},
		{
			name: "compiler diagnostics",
			events: []event{
				{Action: "build-output", ImportPath: "compiler", Output: "first compiler diagnostic\nmiddle compiler diagnostic\n"},
				{Action: "build-output", ImportPath: "compiler", Output: "last compiler diagnostic\nFAIL\n"},
				{Action: "build-fail", ImportPath: "compiler"},
			},
			fragments: []string{"first compiler diagnostic", "middle compiler diagnostic", "last compiler diagnostic"},
			escapes:   []string{`first compiler diagnostic\\nmiddle compiler diagnostic\\nlast compiler diagnostic`},
			first:     "first compiler diagnostic",
			want:      Outcome{Kind: OutcomeBuildFailed},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			set := cannedSet{exit: 1}
			for _, e := range tc.events {
				encoded, err := json.Marshal(e)
				if err != nil {
					t.Fatal(err)
				}
				set.events = append(set.events, string(encoded))
			}
			installCannedGo(t, set)
			installCannedSelection(t)
			for _, full := range []bool{false, true} {
				args := []string{"--package", "./..."}
				if full {
					args = append(args, "--full")
				}
				request, line, code := Prepare(t.TempDir(), args)
				if line != "" {
					t.Fatalf("Prepare = (%d,%q)", code, line)
				}
				outcome, output, code := Execute(t.TempDir(), request)
				if code != 1 || outcome != tc.want {
					t.Fatalf("full=%t outcome=%+v code=%d, want %+v and1\n%s", full, outcome, code, tc.want, output)
				}
				if strings.ContainsAny(output, "\x1b\t") || strings.Contains(output, "--- FAIL:") {
					t.Fatalf("unsafe or runner output:\n%s", output)
				}
				if !full {
					if !strings.Contains(output, sanitize.Preview(tc.first)) || strings.Contains(output, "middle") || strings.Contains(output, "last diagnostic") {
						t.Fatalf("default preview changed:\n%s", output)
					}
					continue
				}
				// These expectations pin escaped separators and controls at the command seam.
				for _, escaped := range tc.escapes {
					if !strings.Contains(output, escaped) {
						t.Errorf("full output lost escaped diagnostic bytes %q:\n%s", escaped, output)
					}
				}
				previous := -1
				for _, fragment := range tc.fragments {
					position := strings.Index(output, fragment)
					if position <= previous {
						t.Errorf("full output lost or reordered %q:\n%s", fragment, output)
					}
					previous = position
				}
			}
		})
	}
}
