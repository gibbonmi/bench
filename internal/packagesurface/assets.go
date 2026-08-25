package packagesurface

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

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
	".bench/hooks/worktree-lifecycle.sh",
	".bench/hooks/block-bench-follow-on.sh",
	".bench/lib/resolve-bench.sh",
	".claude/README.md",
	".claude/settings.json",
	".codex/hooks.json",
}

// ForbiddenPackAssets are paths the npm tarball must never carry: build/dev-only
// files and the kit's own working docs. Adoption generates a consumer's AGENTS.md,
// CLAUDE.md, and HANDOFF equivalents from constants, so shipping the kit's copies
// would export a stale handoff and a duplicate agreement. Both the surface contract
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
// Root-level non-test sources (the module root can carry its own package alongside
// cmd/ and internal/) and every source's //go:embed targets are derived the same way,
// so a new root package or a new embed joins the expectation without a second
// registry.
func RequiredBuildPackAssets(root string) ([]string, error) {
	assets := []string{
		"go.mod",
		"go.sum",
		"scripts/go-build.sh",
		"scripts/go-build.inputs",
	}
	var goFiles []string
	rootEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range rootEntries {
		if entry.IsDir() || !isBuildGoSource(entry.Name()) {
			continue
		}
		assets = append(assets, entry.Name())
		goFiles = append(goFiles, filepath.Join(root, entry.Name()))
	}
	for _, dir := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || !isBuildGoSource(entry.Name()) {
				return err
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			assets = append(assets, filepath.ToSlash(rel))
			goFiles = append(goFiles, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	embeds, err := embeddedPackAssets(root, goFiles)
	if err != nil {
		return nil, err
	}
	assets = append(assets, embeds...)
	sort.Strings(assets)
	return assets, nil
}

func isBuildGoSource(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

var goEmbedDirective = regexp.MustCompile(`^//go:embed\s+(.+)$`)

// embeddedPackAssets reads each Go source in goFiles and returns the repo-relative
// path of every //go:embed target it names, resolved against that source's own
// directory (embed patterns are directory-relative). A source with no embed
// directive contributes nothing.
func embeddedPackAssets(root string, goFiles []string) ([]string, error) {
	var assets []string
	for _, file := range goFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		dir := filepath.Dir(file)
		for _, line := range strings.Split(string(data), "\n") {
			match := goEmbedDirective.FindStringSubmatch(strings.TrimSpace(line))
			if match == nil {
				continue
			}
			for _, pattern := range strings.Fields(match[1]) {
				pattern = strings.Trim(pattern, `"`)
				rel, err := filepath.Rel(root, filepath.Join(dir, filepath.FromSlash(pattern)))
				if err != nil {
					return nil, err
				}
				assets = append(assets, filepath.ToSlash(rel))
			}
		}
	}
	return assets, nil
}
