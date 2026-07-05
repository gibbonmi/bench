package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func checkPackageCoreAndGuards(root string) []string {
	var diags []string
	diags = append(diags, checkPackageFiles(root)...)
	diags = append(diags, checkGoCore(root)...)
	diags = append(diags, checkReleaseWorkflow(root)...)
	diags = append(diags, checkGuardDescribeManifests(root)...)
	return diags
}

func checkPackageFiles(root string) []string {
	path := filepath.Join(root, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkg struct {
		Files []string `json:"files"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	var diags []string
	for _, file := range pkg.Files {
		if !exists(filepath.Join(root, filepath.FromSlash(file))) {
			diags = append(diags, "package.json files[] missing "+file)
		}
	}
	if len(pkg.Files) == 0 {
		return diags
	}

	probe := runAtCleanEnv(root, "npm", "pack", "--dry-run", "--json")
	if probe != nil && probe.ExitCode != 0 {
		diags = append(diags, "npm pack --dry-run failed")
	} else if probe != nil {
		diags = append(diags, checkNpmPackAssets(probe.Stdout)...)
	}
	return append(diags, checkRepoOnlyPackageClaims(root)...)
}

func checkNpmPackAssets(packJSON string) []string {
	var packs []struct {
		Files []struct {
			Path string `json:"path"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(packJSON), &packs); err != nil {
		return []string{"npm pack --dry-run JSON unreadable: " + err.Error()}
	}
	files := map[string]bool{}
	if len(packs) > 0 {
		for _, file := range packs[0].Files {
			files[file.Path] = true
		}
	}
	var diags []string
	for _, required := range []string{
		"bin/bench.sh",
		"bin/bench-postinstall.sh",
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
	} {
		if !files[required] {
			diags = append(diags, "npm package missing "+required)
		}
	}
	for _, forbidden := range []string{".claude/settings.local.json"} {
		if files[forbidden] {
			diags = append(diags, "npm package includes local-only file "+forbidden)
		}
	}
	return diags
}

func checkRepoOnlyPackageClaims(root string) []string {
	// Mirrors the package fragment's lightweight prose sweep over shipped markdown.
	var diags []string
	files := packageMarkdownFiles(root)
	claimRe := regexp.MustCompile(`(?i)\b(ship|ships|shipped|shipping|package|packaged|tarball|installable|included|includes)\b`)
	headingRe := regexp.MustCompile(`(?i)^#{1,6}\s+.*\b(ship|ships|shipped|shipping|package|packaged|tarball|installable|included|includes|surfaces?)\b`)
	repoOnlyRe := regexp.MustCompile(`(?i)\b(repo-only|development context|local development|not shipped|not in the npm package|not part of the npm package)\b`)
	for _, file := range files {
		lines := strings.Split(readIfExists(file), "\n")
		inClaimSection := false
		for i, line := range lines {
			if strings.HasPrefix(line, "#") {
				inClaimSection = headingRe.MatchString(line)
			}
			if !(inClaimSection || claimRe.MatchString(line)) || repoOnlyRe.MatchString(line) {
				continue
			}
			for _, repoOnlyPath := range []string{"projects/", "specs/", "decisions/", "tests/"} {
				if strings.Contains(line, repoOnlyPath) {
					diags = append(diags, fmt.Sprintf("%s:%d claims repo-only path '%s' is shipped/package content; label it repo-only development context", slashRel(root, file), i+1, repoOnlyPath))
				}
			}
		}
	}
	return diags
}

func checkGoCore(root string) []string {
	if !exists(filepath.Join(root, "go.mod")) {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return []string{"go.mod present but no Go toolchain on PATH — the compiled core is load-bearing; install Go"}
	}
	var diags []string
	if probe := runAtCleanEnv(root, "gofmt", "-l", "."); probe != nil && strings.TrimSpace(probe.Stdout) != "" {
		diags = append(diags, "gofmt: unformatted Go files: "+strings.Join(strings.Fields(probe.Stdout), " "))
	}
	buildHelper := filepath.Join(root, "scripts", "go-build.sh")
	if exists(buildHelper) {
		if probe := runAtCleanEnv(root, "bash", buildHelper, root, filepath.Join(root, "dist", "bench")); probe == nil || probe.ExitCode != 0 {
			diags = append(diags, "go build failed")
		}
	} else if probe := runAtCleanEnv(root, "go", "build", "./..."); probe == nil || probe.ExitCode != 0 {
		diags = append(diags, "go build failed")
	}
	if probe := runAtCleanEnv(root, "go", "vet", "./..."); probe == nil || probe.ExitCode != 0 {
		diags = append(diags, "go vet failed")
	}
	testPackages, ok := goCoreTestPackages(root)
	if !ok {
		diags = append(diags, "go list failed")
	} else if len(testPackages) > 0 {
		args := append([]string{"go", "test"}, testPackages...)
		if probe := runAtCleanEnv(root, args...); probe == nil || probe.ExitCode != 0 {
			diags = append(diags, "go test failed")
		}
	}
	if exists(filepath.Join(root, "scripts", "platforms.json")) && exists(buildHelper) {
		matrix, err := platformMatrix(filepath.Join(root, "scripts", "platforms.json"))
		if err != nil {
			diags = append(diags, "platform matrix unreadable: "+err.Error())
		}
		tmp, err := os.MkdirTemp("", "bench-cross-*")
		if err != nil {
			diags = append(diags, "cross-compile setup failed: "+err.Error())
		} else {
			defer os.RemoveAll(tmp)
			for _, target := range matrix {
				env := append(conformanceSubprocessEnv(), "GOOS="+target.Goos, "GOARCH="+target.Goarch)
				probe := runAtEnv(root, env, "bash", buildHelper, root, filepath.Join(tmp, "bench-"+target.Goos+"-"+target.Goarch))
				if probe == nil || probe.ExitCode != 0 {
					diags = append(diags, fmt.Sprintf("cross-compile failed: %s/%s", target.Goos, target.Goarch))
				}
			}
		}
	}
	return diags
}

func goCoreTestPackages(root string) ([]string, bool) {
	probe := runAtCleanEnv(root, "go", "list", "./...")
	if probe == nil || probe.ExitCode != 0 {
		return nil, false
	}
	var packages []string
	for _, pkg := range strings.Fields(probe.Stdout) {
		// The gate runs internal/contract separately with BENCH_CONTRACT_ROOT set.
		if pkg == "internal/contract" || strings.HasSuffix(pkg, "/internal/contract") {
			continue
		}
		packages = append(packages, pkg)
	}
	return packages, true
}

func checkReleaseWorkflow(root string) []string {
	if !exists(filepath.Join(root, "scripts", "platforms.json")) {
		return nil
	}
	wf := filepath.Join(root, ".github", "workflows", "release.yml")
	if !exists(wf) {
		return []string{"release workflow missing (.github/workflows/release.yml)"}
	}
	text := readIfExists(wf)
	var diags []string
	if !regexp.MustCompile(`(?m)^\s*tags:`).MatchString(text) {
		diags = append(diags, "release workflow does not trigger on tags")
	}
	if !strings.Contains(text, "scripts/platforms.json") {
		diags = append(diags, "release workflow does not derive targets from the matrix (scripts/platforms.json)")
	}
	if !strings.Contains(text, "scripts/gen-platform-packages.sh") {
		diags = append(diags, "release workflow does not run the platform-package generator")
	}
	if !strings.Contains(text, "npm publish") {
		diags = append(diags, "release workflow does not publish to npm")
	}
	if !strings.Contains(text, "provenance") {
		diags = append(diags, "release workflow does not publish with provenance")
	}
	return diags
}

func checkGuardDescribeManifests(root string) []string {
	var diags []string
	for _, guard := range []string{"block-dangerous-git", "check-agent-line", "stop", "session-start"} {
		path := filepath.Join(root, ".bench", "hooks", guard+".sh")
		if !exists(path) {
			continue
		}
		probe := runAtCleanEnv(root, "bash", path, "--describe")
		if probe == nil || probe.ExitCode != 0 {
			exit := 1
			if probe != nil {
				exit = probe.ExitCode
			}
			diags = append(diags, fmt.Sprintf("guard %s --describe did not exit 0 (exit %d)", guard, exit))
			continue
		}
		for _, key := range []string{"name", "boundary", "denies", "why"} {
			if !regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `: `).MatchString(probe.Stdout) {
				diags = append(diags, fmt.Sprintf("guard %s --describe manifest missing %s key", guard, key))
			}
		}
		if guard == "session-start" && !regexp.MustCompile(`(?m)^denies: nothing \(informational\)$`).MatchString(probe.Stdout) {
			diags = append(diags, "session-start --describe is not classified informational (denies: nothing)")
		}
	}
	return diags
}

func packageMarkdownFiles(root string) []string {
	path := filepath.Join(root, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var pkg struct {
		Files []string `json:"files"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	var out []string
	for _, entry := range pkg.Files {
		full := filepath.Join(root, filepath.FromSlash(entry))
		info, err := os.Stat(full)
		if err != nil {
			continue
		}
		if info.IsDir() {
			_ = filepath.WalkDir(full, func(path string, d os.DirEntry, err error) error {
				if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
					out = append(out, path)
				}
				return nil
			})
			continue
		}
		if strings.HasSuffix(full, ".md") {
			out = append(out, full)
		}
	}
	return uniqueSorted(out)
}

type platformTarget struct {
	Goos   string `json:"goos"`
	Goarch string `json:"goarch"`
}

func platformMatrix(path string) ([]platformTarget, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var targets []platformTarget
	if err := json.Unmarshal(data, &targets); err != nil {
		return nil, err
	}
	return targets, nil
}
