package anchors

import "testing"

// TestBuildCostAnchors keeps the build-cost guidance rows red-capable. Each needle and
// diagnostic is written here independently of the registry, so a rewrite that drops a row
// cannot define itself green.
func TestBuildCostAnchors(t *testing.T) {
	const implement = ".agents/commands/bench-implement-spec.md"
	const plan = ".agents/skills/bench-craft-spec/references/completion-plan.md"
	const drain = ".agents/commands/bench-drain.md"
	const issues = "docs/known-issues.md"
	const reference = ".bench/BENCH-reference.md"
	const synthesis = ".agents/skills/bench-craft-synthesis/SKILL.md"
	const guide = ".bench/BENCH.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, needle: "then runs each verification requirement once at that commit", want: "build cost: the implement phase restored the run of each verification requirement", forbidden: true},
		{file: implement, needle: "check-floor", want: "build cost: the implement phase restored the retired check-floor row", forbidden: true},
		{file: plan, needle: "check-floor", want: "build cost: the completion-plan reference restored the retired check-floor row", forbidden: true},
		{file: drain, needle: "60 min", want: "build cost: the drain repeats the wall budget of the cost renderer", forbidden: true},
		{file: drain, needle: "7.1", want: "build cost: the drain repeats the bookkeeping budget of the cost renderer", forbidden: true},
		{file: issues, needle: "Run it as `bench preflight --in <target> evidence <id> --check-current`.", want: "build cost: the known issues dropped the --in form of the evidence check"},
		{file: issues, needle: "names the exec form", want: "build cost: the known-issues retire column restored the exec form", forbidden: true},
		{file: reference, needle: "For a worktree, run `bench probe --in <target> ...`.", want: "build cost: the command notes dropped the --in form of the worktree probe"},
		{file: reference, needle: "`bench worktree exec <target> -- bench probe ...`", want: "build cost: the command notes restored the exec probe form", forbidden: true},
		{file: synthesis, needle: "Take the green verdict from the whole-project gate. That gate is the landing's gate, or `bench gate --in <target>` for a batch that waits for approval.", want: "build cost: craft-synthesis dropped the whole-project gate or its --in batch form"},
		{file: guide, needle: "record every plan or gate expansion with `bench learning`", want: "build cost: the operating guide restored the learning for each Writes: widen", forbidden: true},
	}}.check(t)
}
