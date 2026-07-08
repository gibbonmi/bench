package packagesurface

var RequiredPackAssets = []string{
	"bin/bench.sh",
	"bin/bench-repair-binary.mjs",
	"bin/bench-postinstall.sh",
	"projects/gl-axi.md",
	".agents/commands/bench-implement-spec.md",
	".agents/skills/bench-craft-seams/SKILL.md",
	".agents/skills/bench-implement-spec/SKILL.md",
	".agents/skills/bench-implement-spec/agents/openai.yaml",
	".bench/BENCH.md",
	".bench/BENCH-reference.md",
	".bench/adapters/claude",
	".bench/adapters/codex",
	".bench/adapters/opencode",
	".bench/hooks/stop.sh",
	".bench/lib/resolve-bench.sh",
	".claude/README.md",
	".codex/hooks.json",
}

// ForbiddenPackAssets are paths the npm tarball must never carry: build/dev-only
// files and the kit's own working docs (adoption generates a consumer's AGENTS.md,
// CLAUDE.md, and HANDOFF equivalents from constants, so shipping the kit's copies
// would export a stale handoff and a duplicate agreement). Both the surface contract
// test and the conformance package check iterate this one list.
//
// dist/bench is deliberately in neither list: with dist/ in files[] the npx-from-git
// prepare build packs it (present on a built tree, absent in the CI publish checkout),
// so requiring or forbidding it would make the dry-run shape check tree-state-dependent.
var ForbiddenPackAssets = []string{
	".bench/gate.sh",
	"scripts/go-build.sh",
	"scripts/gen-platform-packages.sh",
	".claude/settings.local.json",
	"HANDOFF.md",
	"CLAUDE.md",
	"AGENTS.md",
}
