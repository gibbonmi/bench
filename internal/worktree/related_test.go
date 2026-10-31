package worktree

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// ER32, ER33: a spec filter relates an assignment only through a scope that the commitment
// state gives it: a binding's deliverable, or a continuation's declared scope that holds the
// spec folder or lies within it. A label, a branch, or a path that names the spec relates
// nothing. An assignment with no scope is excluded and counted as unbound, and the matched
// rows follow assignment-ID order. A slug names its live folder spec, and a path operand
// names its spec from the caller's directory.
func TestRelatedWorktreeFilter(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	id := func(c byte) string { return strings.Repeat(string(c), 32) }
	assignment := func(c byte, label string) intent.Assignment {
		return intent.Assignment{
			ID: id(c), Label: label, RequestToken: "request-" + string(c), State: intent.StateComplete,
			Branch: "refs/heads/specs/x/spec.md", Worktree: filepath.Join(root, "specs", "x", string(c)),
		}
	}
	ledger := intent.Ledger{
		Assignments: []intent.Assignment{
			assignment('e', "specs/x/spec.md"),
			assignment('d', "continuation inside the folder"),
			assignment('a', "bound deliverable"),
			assignment('c', "continuation over the folder"),
			assignment('b', "neighbor spec"),
			assignment('f', "other continuation"),
		},
		Commitment: &intent.CommitmentState{
			Bindings: []intent.DeliveryBinding{
				{Assignment: id('a'), Deliverable: "specs/x/spec.md"},
				{Assignment: id('b'), Deliverable: "specs/xy/spec.md"},
				{Assignment: id('9'), Deliverable: "specs/x/spec.md"},
			},
			Continuations: []intent.LegacyContinuation{
				{Assignment: id('c'), Scope: []string{"internal/other", "specs"}},
				{Assignment: id('d'), Scope: []string{"specs/x/tickets"}},
				{Assignment: id('f'), Scope: []string{"specs/xy", "specs/x.md"}},
			},
		},
	}
	for _, tc := range []struct {
		name, spec, summary string
		ids                 []string
	}{
		{name: "bindings and continuation scopes", spec: "specs/x/spec.md", summary: "specs/x/spec.md,3,1", ids: []string{id('a'), id('c'), id('d')}},
		{name: "only the continuation over every spec folder", spec: "specs/none/spec.md", summary: "specs/none/spec.md,1,1", ids: []string{id('c')}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			content, err := relatedProjection(defaultJoins(), &git.Run{}, root, ledger, tc.spec)
			if err != nil {
				t.Fatal(err)
			}
			document, err := axitest.DecodeDocument(string(content))
			if err != nil {
				t.Fatalf("decode %q: %v", content, err)
			}
			summary, err := document.Rows(relatedTable)
			if err != nil || len(summary) != 1 {
				t.Fatalf("related rows = %#v, %v", summary, err)
			}
			cells := summary[0].(map[string]any)
			if got := fmt.Sprint(cells["spec"], ",", cells["matched"], ",", cells["unbound"]); got != tc.summary {
				t.Fatalf("summary = %s, want %s in %q", got, tc.summary, content)
			}
			rows, err := document.Rows(listTable)
			if err != nil {
				t.Fatal(err)
			}
			var ids []string
			for _, row := range rows {
				ids = append(ids, row.(map[string]any)["id"].(string))
			}
			if strings.Join(ids, " ") != strings.Join(tc.ids, " ") {
				t.Fatalf("ids = %q, want %q", ids, tc.ids)
			}
		})
	}
	sub := filepath.Join(root, "sub")
	for _, tc := range []struct{ cwd, operand, want string }{
		{root, "x", "specs/x/spec.md"},
		{root, "x.md", "specs/x/spec.md"},
		{sub, "x", "specs/x/spec.md"},
		{root, "./specs/x/spec.md", "specs/x/spec.md"},
		{sub, "../specs/x/spec.md", "specs/x/spec.md"},
		{sub, filepath.Join(root, "specs", "x", "spec.md"), "specs/x/spec.md"},
		{root, "specs/x", ""},
		{root, "../specs/x/spec.md", ""},
		{sub, "specs/x", ""},
		{root, filepath.Join(filepath.Dir(root), "spec.md"), ""},
	} {
		got, err := relatedSpec(root, tc.cwd, tc.operand)
		if (err != nil) != (tc.want == "") || err == nil && got != tc.want {
			t.Fatalf("relatedSpec(%q) from %s = %q, %v, want %q", tc.operand, tc.cwd, got, err, tc.want)
		}
	}
}
