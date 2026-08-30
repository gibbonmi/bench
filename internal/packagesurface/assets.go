package packagesurface

import (
	"errors"
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
	goFiles, err := buildGoSources(root)
	if err != nil {
		return nil, err
	}
	for _, file := range goFiles {
		rel, err := filepath.Rel(root, file)
		if err != nil {
			return nil, err
		}
		assets = append(assets, filepath.ToSlash(rel))
	}
	embeds, err := EmbedTargets(root)
	if err != nil {
		return nil, err
	}
	assets = append(assets, embeds...)
	sort.Strings(assets)
	return assets, nil
}

// EmbedTargets returns the repo-relative slash path of every //go:embed target the
// checkout's build sources name, resolved against the naming source's own directory
// because an embed pattern is directory-relative. The paths carry the form a
// composed change carries, so a caller can test a changed path for membership
// directly. A source with no embed directive contributes nothing.
func EmbedTargets(root string) ([]string, error) {
	goFiles, err := buildGoSources(root)
	if err != nil {
		return nil, err
	}
	var targets []string
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
				targets = append(targets, filepath.ToSlash(rel))
			}
		}
	}
	return targets, nil
}

// buildGoSources returns the absolute path of every non-test Go source the build
// packages: the module root's own package, plus everything under cmd/ and internal/.
// Both pack derivations walk this one enumeration, so a new source directory joins
// them together. An absent cmd/ or internal/ contributes nothing, because a module
// that carries one of the two is a legitimate tree.
func buildGoSources(root string) ([]string, error) {
	var goFiles []string
	rootEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range rootEntries {
		if entry.IsDir() || !isBuildGoSource(entry.Name()) {
			continue
		}
		goFiles = append(goFiles, filepath.Join(root, entry.Name()))
	}
	for _, dir := range []string{"cmd", "internal"} {
		base := filepath.Join(root, dir)
		err := filepath.WalkDir(base, func(path string, entry fs.DirEntry, err error) error {
			if err != nil && path == base && errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			if err != nil || entry.IsDir() || !isBuildGoSource(entry.Name()) {
				return err
			}
			goFiles = append(goFiles, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return goFiles, nil
}

func isBuildGoSource(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}

var goEmbedDirective = regexp.MustCompile(`^//go:embed\s+(.+)$`)
