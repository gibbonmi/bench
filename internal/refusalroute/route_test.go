package refusalroute

import "testing"

// TestRouteRendering pins the rendered route of injected faces, so each case controls its
// own steps and facts. Each expectation spells the route from the declared steps and the
// fixture facts. It covers RR05, RR06, and RR13.
func TestRouteRendering(t *testing.T) {
	commitAt := TreeCommand("bench commit", Fact("label"), Text("-m"), Operator("msg"), Text("--"), Operators("path"))
	cases := []struct {
		name  string
		face  Face
		facts Facts
		want  string
	}{
		{
			name: "a reviewer route starts with the reviewer marker",
			face: Face{Authority: Reviewer, Route: []Step{Instruction(Text("finish the merge in progress"))}},
			want: "reviewer: finish the merge in progress",
		},
		{
			name: "an agent route joins its steps in declared order",
			face: Face{Authority: Agent, Route: []Step{
				Instruction(Text("repair the red")),
				Command(Text("bench doctor")),
				Command(Composed("rerun")),
			}},
			facts: Facts{Values: map[string]string{"rerun": "bench worktree land --request r1"}},
			want:  "repair the red; then bench doctor; then bench worktree land --request r1",
		},
		{
			name:  "a tree-scoped step names its worktree first after the verb",
			face:  Face{Authority: Agent, Route: []Step{commitAt}},
			facts: Facts{Values: map[string]string{"label": "ft-build"}},
			want:  "bench commit --in 'ft-build' -m <msg> -- <path>...",
		},
		{
			name: "a tree-scoped step with no owning assignment prints the label placeholder",
			face: Face{Authority: Agent, Route: []Step{commitAt}},
			want: "bench commit --in <label> -m <msg> -- <path>...",
		},
		{
			name: "a tree-scoped step at the primary checkout names it",
			face: Face{Authority: Agent, Route: []Step{TreeCommand("bench gate", Text("primary"))}},
			want: "bench gate --in primary",
		},
		{
			name:  "a value that is not line-safe prints its placeholder",
			face:  Face{Authority: Agent, Route: []Step{Command(Text("bench worktree reset --to"), Fact("published-commit"), Fact("checkout"))}},
			facts: Facts{Values: map[string]string{"published-commit": "abc123", "checkout": "/tmp/a\nb"}},
			want:  "bench worktree reset --to 'abc123' <checkout>",
		},
		{
			name:  "a line-safe value is shell-quoted",
			face:  Face{Authority: Agent, Route: []Step{Command(Text("bench worktree reset --to"), Fact("published-commit"), Fact("checkout"))}},
			facts: Facts{Values: map[string]string{"published-commit": "abc123", "checkout": "/tmp/it's here"}},
			want:  `bench worktree reset --to 'abc123' '/tmp/it'\''s here'`,
		},
		{
			name:  "a composed command that is not line-safe prints its placeholder",
			face:  Face{Authority: Agent, Route: []Step{Command(Text("bench doctor")), Command(Composed("rerun"))}},
			facts: Facts{Values: map[string]string{"rerun": "bench worktree land\x1b"}},
			want:  "bench doctor; then <rerun>",
		},
		{
			name:  "a preface comes first after the reviewer marker",
			face:  Face{Authority: Reviewer, Route: []Step{Instruction(Text("finish the merge in progress")), Command(Composed("rerun"))}},
			facts: Facts{Preface: "later proofs in this group did not run", Values: map[string]string{"rerun": "bench worktree land"}},
			want:  "reviewer: later proofs in this group did not run; finish the merge in progress; then bench worktree land",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.face.Render(tc.facts); got != tc.want {
				t.Fatalf("rendered route:\n got %q\nwant %q", got, tc.want)
			}
		})
	}
}
