// Package chargesource names the fixed canonical sources every charge preparation reads.
// The preflight producer and its shared test fixtures both read these paths, and neither
// may import the other, so the paths live here.
package chargesource

const (
	// DelegateSkill is the delegation skill a charge returns.
	DelegateSkill = ".agents/skills/bench-craft-delegate/SKILL.md"
	// DelegateProcedure is the delegation procedure a charge returns.
	DelegateProcedure = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	// BuildPhase is the build phase command a build charge checks.
	BuildPhase = ".agents/commands/bench-implement-spec.md"
	// ReviewSkill is the review skill a review charge checks.
	ReviewSkill = ".agents/skills/bench-craft-review/SKILL.md"
	// ReviewPhase is the review phase command a review charge returns.
	ReviewPhase = ".agents/commands/bench-review-implementation.md"
)
