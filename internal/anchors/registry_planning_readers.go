package anchors

// planningReaderAnchors pin the exact reader's guidance: the focused-read route to the
// native file reader, the profile's root and worktree read contracts, and the reference's
// account of the raw-read advice.
var planningReaderAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-cli/SKILL.md", Kind: Require, Class: Command, Needle: "`bench read <path>`, then each continuation it offers", Diagnostic: "exact planning readers: bench-craft-cli dropped the bench read focused-read route"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Require, Class: Route, Needle: "`bench read` returns exact file", Diagnostic: "exact planning readers: project profile dropped the bench read contract"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Require, Class: Route, Needle: "`bench worktree read --in <label|primary>` pages the same file", Diagnostic: "exact planning readers: project profile dropped the bench worktree read contract"},
	{Group: AfterImplementSpec, File: ".bench/BENCH-reference.md", Kind: Require, Class: Route, Needle: "the matching `bench read` form on stderr", Diagnostic: "exact planning readers: BENCH-reference dropped the raw-read advice account"},
}
