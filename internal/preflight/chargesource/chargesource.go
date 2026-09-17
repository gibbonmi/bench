// Package chargesource names the fixed canonical sources every build charge reads. The
// preflight producer and its shared test fixtures both read these paths, and neither may
// import the other, so the paths live here.
package chargesource

const (
	// DelegateSkill is the delegation skill a build charge returns.
	DelegateSkill = ".agents/skills/bench-craft-delegate/SKILL.md"
	// DelegateProcedure is the delegation procedure a build charge returns.
	DelegateProcedure = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	// BuildPhase is the build phase command a build charge checks.
	BuildPhase = ".agents/commands/bench-implement-spec.md"
)
