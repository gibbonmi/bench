package packagesurface

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

var RequiredPackAssets = []string{
	"bin/bench.sh",
	"bin/bench-repair-binary.mjs",
	"bin/bench-postinstall.sh",
	"projects/gl-axi.md",
	"projects/regroup.md",
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
	".bench/hooks/worktree-lifecycle.sh",
	".bench/lib/resolve-bench.sh",
	".claude/README.md",
	".claude/settings.json",
	".codex/hooks.json",
}

// ForbiddenPackAssets are paths the npm tarball must never carry: build/dev-only
// files and the kit's own working docs (adoption generates a consumer's AGENTS.md,
// CLAUDE.md, and HANDOFF equivalents from constants, so shipping the kit's copies
// would export a stale handoff and a duplicate agreement). Both the surface contract
// test and the conformance package check iterate this one list.
var ForbiddenPackAssets = []string{
	"projects/benchkit.md",
	".bench/gate.sh",
	"scripts/gen-platform-packages.sh",
	".claude/settings.local.json",
	"HANDOFF.md",
	"CLAUDE.md",
	"AGENTS.md",
}

// RequiredBuildPackAssets derives the Go build inputs that an npm git install's
// prepare script needs. Source directories are packaged recursively, so a split
// package automatically joins this expectation without a second file registry.
func RequiredBuildPackAssets(root string) ([]string, error) {
	assets := []string{
		"go.mod",
		"go.sum",
		"internal/releaseevidence/requirements.json",
		"scripts/go-build.sh",
	}
	for _, dir := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			assets = append(assets, filepath.ToSlash(rel))
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Strings(assets)
	return assets, nil
}
