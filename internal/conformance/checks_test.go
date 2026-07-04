package conformance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/lines"
)

func RunConformance(root, kitRoot string) []string {
	var diags []string
	diags = append(diags, checkLoadValidityMetadata(root)...)
	diags = append(diags, checkSkillsIndexAndCommandAdapters(root)...)
	diags = append(diags, checkDocsCurrencyAndWorkflow(root, kitRoot)...)
	diags = append(diags, checkLineRouting(root)...)
	diags = append(diags, checkPackageCoreAndGuards(root)...)
	return diags
}

func containsDiagnostic(diags []string, want string) bool {
	for _, diag := range diags {
		if strings.Contains(diag, want) {
			return true
		}
	}
	return false
}

func checkLoadValidityMetadata(root string) []string {
	var diags []string
	diags = append(diags, checkShellSyntax(root)...)
	diags = append(diags, checkExecutableGitModes(root)...)
	diags = append(diags, checkExtensionlessGateRefs(root)...)
	diags = append(diags, checkJSONValidity(root)...)
	diags = append(diags, checkCodexHooks(root)...)
	diags = append(diags, checkSkillFrontmatter(root)...)
	diags = append(diags, checkCraftSkillNames(root)...)
	diags = append(diags, checkClaudeSkillMirror(root)...)
	diags = append(diags, checkSharedRuleSingleSource(root)...)
	return diags
}

func checkShellSyntax(root string) []string {
	var diags []string
	for _, pattern := range []string{
		"bin/*.sh",
		".bench/gate-*.sh",
		".bench/skills-index.sh",
		".bench/hooks/*.sh",
		".bench/lib/*.sh",
		"scripts/*.sh",
	} {
		matches, _ := filepath.Glob(filepath.Join(root, filepath.FromSlash(pattern)))
		sort.Strings(matches)
		for _, path := range matches {
			cmd := exec.Command("bash", "-n", path)
			if err := cmd.Run(); err != nil {
				diags = append(diags, fmt.Sprintf("bash syntax error in %s", slashRel(root, path)))
			}
		}
	}
	return diags
}

func checkExecutableGitModes(root string) []string {
	var diags []string
	for _, pattern := range []string{"bin/bench.sh", ".bench/hooks/*.sh", ".bench/adapters/*"} {
		matches, _ := filepath.Glob(filepath.Join(root, filepath.FromSlash(pattern)))
		sort.Strings(matches)
		for _, path := range matches {
			rel := slashRel(root, path)
			cmd := exec.Command("git", "ls-files", "-s", rel)
			cmd.Dir = root
			out, err := cmd.Output()
			if err != nil || len(out) == 0 {
				continue
			}
			fields := strings.Fields(string(out))
			if len(fields) > 0 && fields[0] != "100755" {
				diags = append(diags, fmt.Sprintf("%s is not executable in git (mode %s); the harness runs it as a command path", rel, fields[0]))
			}
		}
	}
	return diags
}

func checkExtensionlessGateRefs(root string) []string {
	path := filepath.Join(root, "bin", "bench.sh")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	refRE := regexp.MustCompile(`\.bench/[A-Za-z0-9_.-]+`)
	var bad []string
	for _, ref := range refRE.FindAllString(string(data), -1) {
		if ref == ".bench/gate" || ref == ".bench/done" {
			bad = append(bad, ref)
		}
	}
	sort.Strings(bad)
	bad = uniqueStrings(bad)
	if len(bad) == 0 {
		return nil
	}
	return []string{fmt.Sprintf("bin/bench.sh has extensionless gate/done refs (%s); the contract is .sh", strings.Join(bad, " ")+" ")}
}

func checkJSONValidity(root string) []string {
	var diags []string
	for _, rel := range []string{"package.json", ".claude/settings.json", ".codex/hooks.json"} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			if os.IsNotExist(err) {
				diags = append(diags, "JSON file missing: "+rel)
			}
			continue
		}
		var v any
		if err := json.Unmarshal(data, &v); err != nil {
			diags = append(diags, "invalid JSON in "+rel)
		}
	}
	return diags
}

func checkCodexHooks(root string) []string {
	path := filepath.Join(root, ".codex", "hooks.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg struct {
		Hooks map[string][]struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	var diags []string
	var stopCommands, preCommands []string
	for _, group := range cfg.Hooks["Stop"] {
		for _, hook := range group.Hooks {
			stopCommands = append(stopCommands, hook.Command)
		}
	}
	if !anyContains(stopCommands, ".bench/hooks/stop.sh") {
		diags = append(diags, "codex hooks.json Stop event does not run .bench/hooks/stop.sh")
	}
	hasBashMatcher := false
	for _, group := range cfg.Hooks["PreToolUse"] {
		if group.Matcher == "Bash" {
			hasBashMatcher = true
		}
		for _, hook := range group.Hooks {
			preCommands = append(preCommands, hook.Command)
		}
	}
	if !hasBashMatcher {
		diags = append(diags, "codex hooks.json PreToolUse Bash matcher missing")
	}
	if !anyContains(preCommands, ".bench/hooks/block-dangerous-git.sh") {
		diags = append(diags, "codex hooks.json PreToolUse does not run .bench/hooks/block-dangerous-git.sh")
	}
	return diags
}

func checkSkillFrontmatter(root string) []string {
	var diags []string
	files, _ := filepath.Glob(filepath.Join(root, ".agents", "skills", "*", "SKILL.md"))
	sort.Strings(files)
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		first := string(data)
		if i := strings.IndexByte(first, '\n'); i >= 0 {
			first = first[:i]
		}
		if first != "---" {
			diags = append(diags, fmt.Sprintf("%s missing frontmatter", slashRel(root, path)))
		}
	}
	return diags
}

func checkCraftSkillNames(root string) []string {
	var diags []string
	dirs, _ := filepath.Glob(filepath.Join(root, ".agents", "skills", "bench-craft-*"))
	sort.Strings(dirs)
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		base := filepath.Base(dir)
		expected := "craft-" + strings.TrimPrefix(base, "bench-craft-")
		actual := frontmatterField(filepath.Join(dir, "SKILL.md"), "name")
		if actual != expected {
			diags = append(diags, fmt.Sprintf("craft skill '%s' visible name is '%s'; expected '%s'", base, actual, expected))
		}
	}

	files, _ := filepath.Glob(filepath.Join(root, ".agents", "skills", "*", "SKILL.md"))
	sort.Strings(files)
	for _, path := range files {
		dir := filepath.Base(filepath.Dir(path))
		if exists(filepath.Join(root, ".agents", "commands", dir+".md")) {
			continue
		}
		actual := frontmatterField(path, "name")
		if strings.HasPrefix(actual, "bench-") {
			diags = append(diags, fmt.Sprintf("non-command skill '%s' uses bench-* visible name '%s'", dir, actual))
		}
	}
	return diags
}

func checkClaudeSkillMirror(root string) []string {
	if !exists(filepath.Join(root, ".claude", "skills")) && !exists(filepath.Join(root, ".agents", "skills")) {
		return nil
	}
	var diags []string
	dirs, _ := filepath.Glob(filepath.Join(root, ".agents", "skills", "bench-craft-*"))
	sort.Strings(dirs)
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		base := filepath.Base(dir)
		if !exists(filepath.Join(root, ".claude", "skills", base)) {
			diags = append(diags, fmt.Sprintf("craft skill '%s' missing from .claude/skills (Claude loses the guidance surface)", base))
		}
	}
	claudeDirs, _ := filepath.Glob(filepath.Join(root, ".claude", "skills", "*"))
	sort.Strings(claudeDirs)
	for _, dir := range claudeDirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		base := filepath.Base(dir)
		if !strings.HasPrefix(base, "bench-craft-") {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s is not a craft skill (phase adapters are Codex-only; it duplicates the slash menu)", base))
			continue
		}
		if !exists(filepath.Join(dir, "SKILL.md")) {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s does not resolve to a SKILL.md (broken adapter link)", base))
		}
	}
	return diags
}

func checkSharedRuleSingleSource(root string) []string {
	bench := readIfExists(filepath.Join(root, ".bench", "BENCH.md"))
	agents := readIfExists(filepath.Join(root, "AGENTS.md"))
	readme := readIfExists(filepath.Join(root, "README.md"))
	if bench == "" && agents == "" && readme == "" {
		return nil
	}
	var diags []string
	for _, marker := range []string{
		"you never grade your own work",
		"Declare the line before a long run",
		"Document for the teammate who just walked in",
		"One small change at a time, repo stays green",
		"Clear beats dense",
		"You are the worker; I am the reviewer",
		"Right-size the process",
		"never silently rewrite your own rules",
	} {
		if !strings.Contains(bench, marker) {
			diags = append(diags, fmt.Sprintf("shared rule missing from canonical .bench/BENCH.md: %q", marker))
		}
		if strings.Contains(agents, marker) {
			diags = append(diags, fmt.Sprintf("shared rule duplicated in AGENTS.md (it must live only in .bench/BENCH.md): %q", marker))
		}
		if strings.Contains(readme, marker) {
			diags = append(diags, fmt.Sprintf("shared rule duplicated in README.md (README must point to .bench/BENCH.md instead of restating shared rules): %q", marker))
		}
	}
	if !strings.Contains(agents, "canonical in `.bench/BENCH.md`") {
		diags = append(diags, "AGENTS.md lost its pointer to the canonical .bench/BENCH.md shared rules")
	}
	for _, heading := range []string{"## The four invariants", "## The three layers"} {
		if strings.Contains(readme, heading) {
			diags = append(diags, "README.md duplicates shared operating-guide section "+heading+"; point to .bench/BENCH.md instead")
		}
	}
	return diags
}

func checkSkillsIndexAndCommandAdapters(root string) []string {
	var diags []string
	diags = append(diags, checkSkillsIndex(root)...)
	diags = append(diags, checkCommandGuideReferences(root)...)
	diags = append(diags, checkCodexCommandAdapters(root)...)
	diags = append(diags, checkRoadmapPromotionAnchors(root)...)
	return diags
}

func checkDocsCurrencyAndWorkflow(root, kitRoot string) []string {
	var diags []string
	diags = append(diags, checkStaleCommandReferences(root)...)
	diags = append(diags, checkColdPickupCLILists(root)...)
	diags = append(diags, checkAXIProfileAnchors(root)...)
	diags = append(diags, checkBenchReferenceTokenDiet(root)...)
	diags = append(diags, checkCommandFirstAnchors(root)...)
	diags = append(diags, checkWorkflowAnchors(root)...)
	diags = append(diags, checkSkillsIndexGenerateVerify(root, kitRoot)...)
	diags = append(diags, checkCoverageMaps(root)...)
	return diags
}

func checkStaleCommandReferences(root string) []string {
	commandsDir := filepath.Join(root, ".agents", "commands")
	validSlash := map[string]bool{"/model": true}
	commandFiles, _ := filepath.Glob(filepath.Join(commandsDir, "*.md"))
	for _, file := range commandFiles {
		validSlash["/"+strings.TrimSuffix(filepath.Base(file), ".md")] = true
	}
	validCodex := map[string]bool{}
	for slash := range validSlash {
		if strings.HasPrefix(slash, "/bench-") {
			validCodex["$"+strings.TrimPrefix(slash, "/")] = true
		}
	}

	var files []string
	for _, rel := range []string{
		"README.md",
		"AGENTS.md",
		".bench/BENCH.md",
		".bench/BENCH-reference.md",
		".bench/learnings.md",
		"CONTEXT.md",
		"HANDOFF.md",
		"CHANGELOG.md",
	} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if exists(path) {
			files = append(files, path)
		}
	}
	for _, rel := range []string{"specs", "decisions", ".agents"} {
		files = append(files, walkConformanceDocs(filepath.Join(root, filepath.FromSlash(rel)))...)
	}
	files = uniqueSorted(files)

	knownStale := map[string]bool{
		"/resynthesize":   true,
		"/spec":           true,
		"/grill":          true,
		"/start-ideation": true,
		"/setup":          true,
		"/build":          true,
		"/prep-shift":     true,
		"/fix-bug":        true,
		"/verify-gate":    true,
		"/map":            true,
		"/diagnose":       true,
		"/review":         true,
		"/verify":         true,
		"/shift":          true,
	}
	historicalMarker := regexp.MustCompile(`(?m)^<!-- command-currency: historical -->$`)
	slashRef := regexp.MustCompile(`(^|[\s([` + "`" + `"'])/([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])`)
	codexRef := regexp.MustCompile(`(^|[\s([` + "`" + `"'])\$([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])`)

	var diags []string
	for _, file := range files {
		text := readIfExists(file)
		rel := slashRel(root, file)
		if historicalMarker.MatchString(text) {
			continue
		}
		switch rel {
		case ".bench/learnings.md":
			text = strings.Split(text, "<!-- entries below -->")[0]
		case "CHANGELOG.md":
			if idx := strings.Index(text, "\n## "); idx >= 0 {
				text = text[:idx]
			}
		}
		for i, line := range strings.Split(text, "\n") {
			for _, match := range slashRef.FindAllStringSubmatch(line, -1) {
				token := "/" + match[2]
				if !validSlash[token] && (strings.HasPrefix(token, "/bench-") || knownStale[token]) {
					diags = append(diags, fmt.Sprintf("stale command reference %s in %s:%d", token, rel, i+1))
				}
			}
			for _, match := range codexRef.FindAllStringSubmatch(line, -1) {
				token := "$" + match[2]
				if strings.HasPrefix(token, "$bench-") && !validCodex[token] {
					diags = append(diags, fmt.Sprintf("stale Codex adapter reference %s in %s:%d", token, rel, i+1))
				}
			}
		}
	}
	return diags
}

func checkColdPickupCLILists(root string) []string {
	bench := readIfExists(filepath.Join(root, "bin", "bench.sh"))
	if bench == "" {
		return nil
	}
	cmdRE := regexp.MustCompile(`(?m)^  ([a-z][a-z-]*)\)\s`)
	var commands []string
	for _, match := range cmdRE.FindAllStringSubmatch(bench, -1) {
		commands = append(commands, match[1])
	}
	sort.Strings(commands)
	known := map[string]bool{}
	for _, command := range commands {
		known[command] = true
	}
	docRef := regexp.MustCompile("`bench ([a-z][a-z-]*)\\b")
	var diags []string
	for _, rel := range []string{".bench/BENCH.md"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		text := readIfExists(path)
		if text == "" {
			continue
		}
		for _, command := range commands {
			if !strings.Contains(text, "bench "+command) {
				diags = append(diags, fmt.Sprintf("%s does not list CLI command 'bench %s'", rel, command))
			}
		}
	}
	for _, rel := range []string{"HANDOFF.md", ".bench/BENCH.md", ".bench/BENCH-reference.md"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		text := readIfExists(path)
		if text == "" {
			continue
		}
		for _, match := range docRef.FindAllStringSubmatch(text, -1) {
			if !known[match[1]] {
				diags = append(diags, fmt.Sprintf("%s documents unknown CLI command 'bench %s' (removed or renamed in bin/bench.sh?)", rel, match[1]))
			}
		}
	}
	return diags
}

func checkAXIProfileAnchors(root string) []string {
	text := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	var diags []string
	if text == "" {
		return nil
	}
	if !strings.Contains(text, "bench diff") {
		diags = append(diags, "projects/benchkit.md does not name bench diff on the AXI seam")
	}
	if !strings.Contains(text, "bench coverage") {
		diags = append(diags, "projects/benchkit.md does not name bench coverage on the AXI seam")
	}
	return diags
}

func checkBenchReferenceTokenDiet(root string) []string {
	benchPath := filepath.Join(root, ".bench", "BENCH.md")
	if !exists(benchPath) {
		return nil
	}
	refRel := ".bench/BENCH-reference.md"
	refPath := filepath.Join(root, filepath.FromSlash(refRel))
	if !exists(refPath) {
		return []string{refRel + " missing: the token-diet reference file the operating guide points to"}
	}
	var diags []string
	bench := readIfExists(benchPath)
	ref := readIfExists(refPath)
	if !strings.Contains(bench, "BENCH-reference.md") {
		diags = append(diags, ".bench/BENCH.md does not point to .bench/BENCH-reference.md (agents can't find the moved lookup sections)")
	}
	claude := readIfExists(filepath.Join(root, "CLAUDE.md"))
	for _, line := range strings.Split(claude, "\n") {
		if strings.TrimSpace(line) == "@.bench/BENCH-reference.md" {
			diags = append(diags, ".bench/BENCH-reference.md is @-imported by CLAUDE.md; it must stay on-demand (referenced by path, never imported, or the token diet regresses)")
			break
		}
	}
	benchHeads := h2Headings(bench)
	refHeads := h2Headings(ref)
	var dup []string
	for head := range benchHeads {
		if refHeads[head] {
			dup = append(dup, head)
		}
	}
	sort.Strings(dup)
	if len(dup) > 0 {
		diags = append(diags, "section heading present in both .bench/BENCH.md and .bench/BENCH-reference.md (single-source violation): "+strings.Join(dup, "|")+"|")
	}
	return diags
}

func checkCommandFirstAnchors(root string) []string {
	var diags []string
	readme := readIfExists(filepath.Join(root, "README.md"))
	if readme == "" {
		diags = append(diags, "README.md missing")
	} else {
		firstH2 := ""
		for _, line := range strings.Split(readme, "\n") {
			if strings.HasPrefix(line, "## ") {
				firstH2 = line
				break
			}
		}
		if firstH2 != "## Reviewer quick start" {
			if firstH2 == "" {
				firstH2 = "(none)"
			}
			diags = append(diags, fmt.Sprintf("README first H2 is '%s'; expected '## Reviewer quick start'", firstH2))
		}
	}

	commandFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "commands", "*.md"))
	sort.Strings(commandFiles)
	for _, file := range commandFiles {
		rel := slashRel(root, file)
		text := readIfExists(file)
		if !regexp.MustCompile(`(?m)^## Entry orientation$`).MatchString(text) {
			diags = append(diags, rel+" missing Entry orientation")
		}
		if !regexp.MustCompile(`(?m)^## Exit handoff$`).MatchString(text) {
			diags = append(diags, rel+" missing Exit handoff")
		}
	}
	return diags
}

func checkWorkflowAnchors(root string) []string {
	var diags []string
	require := func(rel, needle string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !exists(path) {
			diags = append(diags, "acceptance coverage anchor file missing: "+rel)
			return
		}
		if !strings.Contains(readIfExists(path), needle) {
			diags = append(diags, fmt.Sprintf("%s missing acceptance coverage anchor: %s", rel, needle))
		}
	}

	require(".agents/commands/bench-write-spec.md", "acceptance coverage map")
	require(".agents/commands/bench-write-spec.md", "why it catches the failure")
	require(".agents/commands/bench-write-spec.md", "red signal")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "acceptance row")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "not TDD-able")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "call count")
	require(".agents/commands/bench-implement-spec.md", "coverage table")
	require(".agents/commands/bench-implement-spec.md", "already covered")
	require(".agents/commands/bench-implement-spec.md", "turning red-to-green")
	require(".agents/commands/bench-implement-spec.md", "bench coverage <spec>")
	require(".agents/commands/bench-review-implementation.md", "acceptance coverage map")
	require(".agents/commands/bench-review-implementation.md", "mapped behavior")
	require(".agents/commands/bench-review-implementation.md", "bench diff")
	require(".agents/commands/bench-final-check.md", ".bench/gate.sh")
	require(".agents/commands/bench-final-check.md", "BENCH_GATE")
	require(".agents/commands/bench-write-spec.md", "seam diagram")
	require(".agents/commands/bench-write-spec.md", "tests attach here")
	require(".agents/commands/bench-write-spec.md", "edge inventory")
	require(".agents/commands/bench-write-spec.md", "Won't handle")
	require(".agents/commands/bench-write-spec.md", "hostile-input checklist")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "floor, not the ceiling")
	require(".agents/skills/bench-craft-seams/SKILL.md", "failure modes")
	require(".agents/skills/bench-craft-seams/SKILL.md", "structure.budgets")
	require(".agents/commands/bench-review-implementation.md", "## Coverage")
	require(".agents/commands/bench-review-implementation.md", "Coverage axis")
	require(".agents/commands/bench-setup-repo.md", "hostile-input checklist")
	require("projects/benchkit.md", "hostile-input checklist")
	require(".agents/commands/bench-setup-repo.md", "craft-gate")
	require(".agents/commands/bench-final-check.md", "craft-gate")
	require(".agents/commands/bench-review-implementation.md", "craft-review")
	require(".agents/skills/bench-craft-review/SKILL.md", "an edge nobody decided")
	require(".agents/commands/bench-review-implementation.md", "craft-delegate")
	require(".agents/skills/bench-craft-delegate/SKILL.md", "a claim, not a result")
	require(".agents/commands/bench-implement-spec.md", "When the build stops short")
	require(".agents/commands/bench-write-spec.md", "Superseded by")
	require(".agents/commands/bench-debug.md", "before launching the shift")
	require(".agents/commands/bench-shape-idea.md", "## Handoff")
	require(".agents/commands/bench-shape-idea.md", "Hostile-input owner")
	require(".agents/commands/bench-shape-idea.md", "Dependency order")
	require(".agents/commands/bench-shape-idea.md", "n/a \u2014")
	require(".agents/commands/bench-write-spec.md", "map's Handoff")
	require(".agents/commands/bench-write-spec.md", "spec-retire:")
	require(".agents/commands/bench-write-spec.md", "Status: staged")
	require(".agents/commands/bench-write-spec.md", "new session on the mid tier")
	require(".agents/commands/bench-implement-spec.md", "Status: implemented")
	require(".agents/commands/bench-debug.md", "diff-filter=D")

	shapeIdeaPath := filepath.Join(root, ".agents", "commands", "bench-shape-idea.md")
	if exists(shapeIdeaPath) && strings.Contains(collapseSpace(readIfExists(shapeIdeaPath)), "straight to `/bench-write-spec`") {
		diags = append(diags, ".agents/commands/bench-shape-idea.md reintroduces the skip-to-spec bypass fragment; every idea must yield a map with a Handoff before a spec")
	}
	writeSpecPath := filepath.Join(root, ".agents", "commands", "bench-write-spec.md")
	if exists(writeSpecPath) && !strings.Contains(collapseSpace(readIfExists(writeSpecPath)), "refuses to run without") {
		diags = append(diags, ".agents/commands/bench-write-spec.md dropped the map-required entry contract (refuses to run without a complete map)")
	}

	readme := readIfExists(filepath.Join(root, "README.md"))
	if readme != "" {
		if !strings.Contains(readme, "session-start.sh") {
			diags = append(diags, "README layout omits .bench/hooks/session-start.sh")
		}
		if !strings.Contains(readme, "bench.sh") {
			diags = append(diags, "README layout omits the real bin/bench.sh filename")
		}
		if !strings.Contains(readme, "benchkit.md") {
			diags = append(diags, "README layout omits projects/benchkit.md")
		}
		if strings.Contains(readme, "\u2502   \u2514\u2500\u2500 bench                 #") {
			diags = append(diags, "README layout still names bin/bench instead of bin/bench.sh")
		}
	}

	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")); text != "" && !strings.Contains(text, "craft-line") {
		diags = append(diags, "bench-implement-spec.md does not reference craft-line")
	}
	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-write-spec.md")); text != "" {
		if !strings.Contains(text, "craft-line") {
			diags = append(diags, "bench-write-spec.md does not reference craft-line")
		}
		if !strings.Contains(text, "model and effort") {
			diags = append(diags, "bench-write-spec.md does not mandate per-story model and effort")
		}
	}
	if text := readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md")); text != "" && !strings.Contains(text, "BENCH_MODEL") {
		diags = append(diags, "BENCH-reference.md adapter contract does not document BENCH_MODEL")
	}
	return diags
}

func checkSkillsIndexGenerateVerify(root, kitRoot string) []string {
	if !exists(filepath.Join(kitRoot, ".bench", "skills-index.sh")) {
		return nil
	}
	tmp, err := os.MkdirTemp("", "bench-skills-index-*")
	if err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	defer os.RemoveAll(tmp)
	if err := os.MkdirAll(filepath.Join(tmp, ".bench"), 0o755); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	if err := os.MkdirAll(filepath.Join(tmp, ".agents", "skills", "zeta-skill"), 0o755); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	if err := os.WriteFile(filepath.Join(tmp, ".agents", "skills", "zeta-skill", "SKILL.md"), []byte("---\nname: zeta-skill\ndescription: d\nindex: doing zeta things\n---\n"), 0o644); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	if err := os.WriteFile(filepath.Join(tmp, ".bench", "BENCH-reference.md"), []byte("# Reference\n\n<!-- bench:skills-index:start -->\n<!-- bench:skills-index:end -->\n"), 0o644); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	script := filepath.Join(kitRoot, ".bench", "skills-index.sh")
	if probe := runAt(tmp, "bash", script, "--check"); probe == nil || probe.ExitCode == 0 {
		return []string{"skills-index generate/verify contract failed: check passed on an empty index block"}
	}
	if probe := runAt(tmp, "bash", script, "--write"); probe == nil || probe.ExitCode != 0 {
		return []string{"skills-index generate/verify contract failed: --write failed"}
	}
	generated := readIfExists(filepath.Join(tmp, ".bench", "BENCH-reference.md"))
	if !strings.Contains(generated, "- doing zeta things \u2192 `.agents/skills/zeta-skill/SKILL.md`") {
		return []string{"skills-index generate/verify contract failed: --write did not generate the entry from frontmatter"}
	}
	if probe := runAt(tmp, "bash", script, "--check"); probe == nil || probe.ExitCode != 0 {
		return []string{"skills-index generate/verify contract failed: check red right after --write"}
	}
	before := readIfExists(filepath.Join(tmp, ".bench", "BENCH-reference.md"))
	if probe := runAt(tmp, "bash", script, "--write"); probe == nil || probe.ExitCode != 0 {
		return []string{"skills-index generate/verify contract failed: second --write failed"}
	}
	if before != readIfExists(filepath.Join(tmp, ".bench", "BENCH-reference.md")) {
		return []string{"skills-index generate/verify contract failed: --write is not idempotent"}
	}
	return nil
}

func checkCoverageMaps(root string) []string {
	specsDir := filepath.Join(root, "specs")
	if !exists(specsDir) {
		return nil
	}
	matches, _ := filepath.Glob(filepath.Join(specsDir, "*.md"))
	sort.Strings(matches)
	var diags []string
	for _, path := range matches {
		out, code := coverage.Command([]string{"--check", path})
		if code == 0 {
			continue
		}
		if strings.TrimSpace(out) == "" {
			diags = append(diags, fmt.Sprintf("%s coverage --check failed (exit %d) with no message", slashRel(root, path), code))
			continue
		}
		for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			diags = append(diags, strings.TrimPrefix(line, "error: "))
		}
	}
	return diags
}

func checkLineRouting(root string) []string {
	var diags []string
	diags = append(diags, checkLineBinding(root)...)
	diags = append(diags, checkClaudeAgentHookWiring(root)...)
	diags = append(diags, checkAgentHookBehavior(root)...)
	diags = append(diags, checkAdapterLineGuards(root)...)
	return diags
}

func checkLineBinding(root string) []string {
	path := filepath.Join(root, ".bench", "lines.env")
	if !exists(path) {
		return []string{"lines.env missing: .bench/lines.env is the tier binding enforcement reads"}
	}
	content := []byte(readIfExists(path))
	binding := lines.ParseBinding(content)
	var diags []string
	tiers := []struct {
		label string
		key   string
		value string
	}{
		{"top", "BENCH_TIER_TOP", binding.Top},
		{"mid", "BENCH_TIER_MID", binding.Mid},
		{"cheap", "BENCH_TIER_CHEAP", binding.Cheap},
	}
	modelID := regexp.MustCompile(`^claude-[a-z0-9][a-z0-9.-]*$`)
	for _, tier := range tiers {
		if tier.value == "" {
			diags = append(diags, fmt.Sprintf("lines.env tier unset: %s has no value in .bench/lines.env", tier.key))
		} else if !modelID.MatchString(tier.value) {
			diags = append(diags, fmt.Sprintf("lines.env tier malformed: %s='%s' is not a model id", tier.label, tier.value))
		}
	}
	aliases := []struct {
		key   string
		value string
	}{
		{"BENCH_ALIAS_TOP", binding.AliasTop},
		{"BENCH_ALIAS_MID", binding.AliasMid},
		{"BENCH_ALIAS_CHEAP", binding.AliasCheap},
	}
	aliasRe := regexp.MustCompile(`^[a-z0-9-]+$`)
	for _, alias := range aliases {
		if !regexp.MustCompile(`(?m)^[ \t]*` + regexp.QuoteMeta(alias.key) + `=`).Match(content) {
			continue
		}
		if !aliasRe.MatchString(alias.value) {
			diags = append(diags, fmt.Sprintf("lines.env alias malformed: %s='%s' is not a bare alias", alias.key, alias.value))
		}
	}

	profile := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	if profile != "" {
		for _, tier := range tiers {
			if tier.value == "" {
				continue
			}
			if !strings.Contains(profile, tier.value) {
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not name bound model id '%s' (%s in lines.env)", tier.value, tier.key))
			}
		}
		for _, alias := range aliases {
			if alias.value == "" {
				continue
			}
			want := alias.key + "=" + alias.value
			if !strings.Contains(profile, want) {
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not carry alias declaration %s", want))
			}
		}
	}
	return diags
}

func checkClaudeAgentHookWiring(root string) []string {
	path := filepath.Join(root, ".claude", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg struct {
		Hooks map[string][]struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	for _, group := range cfg.Hooks["PreToolUse"] {
		if group.Matcher != "Agent" {
			continue
		}
		for _, hook := range group.Hooks {
			if strings.Contains(hook.Command, ".bench/hooks/check-agent-line.sh") {
				return nil
			}
		}
	}
	return []string{"claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh"}
}

func checkAgentHookBehavior(root string) []string {
	hook := filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")
	realBench := filepath.Join(root, "bin", "bench.sh")
	if !exists(hook) {
		return nil
	}
	if !exists(realBench) {
		probe := runWithInput(root, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "bash", hook)
		if probe != nil && probe.ExitCode == 0 {
			return []string{"check-agent-line.sh does not deny an unbound model"}
		}
		return nil
	}

	bindir, cleanup, err := wrapperStubDir(realBench)
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanup()
	env := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
	var diags []string

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\nBENCH_ALIAS_MID=opus\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupRouted()

	hookCase := func(label, cwd, input, wantErr string, wantExit int) {
		probe := runWithInputEnv(cwd, env, input, "bash", hook)
		if probe == nil || probe.ExitCode != wantExit {
			got := -1
			if probe != nil {
				got = probe.ExitCode
			}
			diags = append(diags, fmt.Sprintf("check-agent-line.sh %s exit %d (want %d)", label, got, wantExit))
			return
		}
		if wantErr != "" && !strings.Contains(probe.Stderr, wantErr) && !strings.Contains(probe.Stdout, wantErr) {
			diags = append(diags, fmt.Sprintf("check-agent-line.sh %s did not warn with %q", label, wantErr))
		}
	}
	hookCase("denies a bound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-opus-4-8"}}`, "", 0)
	hookCase("denies a declared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus"}}`, "", 0)
	hookCase("does not deny an undeclared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"sonnet"}}`, "", 2)
	hookCase("does not deny an unbound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "", 2)
	hookCase("does not fail open on malformed stdin", routed, `not json at all`, "not parseable as JSON", 0)
	hookCase("does not fail open on a missing model field", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "no resolvedModel/model field", 0)

	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	hookCase("does not fail open without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "no .bench/lines.env", 0)

	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupPartial()
	hookCase("does not fail open on an incomplete binding", partial, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "unset or empty", 0)
	return diags
}

func checkAdapterLineGuards(root string) []string {
	if !exists(filepath.Join(root, ".bench", "adapters")) {
		return nil
	}
	var diags []string
	realBench := filepath.Join(root, "bin", "bench.sh")
	bindir := ""
	cleanup := func() {}
	if exists(realBench) {
		var err error
		bindir, cleanup, err = adapterStubDir(realBench)
		if err != nil {
			return []string{"adapter line guard setup failed: " + err.Error()}
		}
		defer cleanup()
	}
	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			diags = append(diags, "adapter missing from .bench/adapters: "+name)
			continue
		}
		text := readIfExists(path)
		hasResolveCall := regexp.MustCompile(`(?m)^[[:space:]]*model="\$\([^)]*resolve-model`).MatchString(text)
		if !hasResolveCall {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
	}
	if bindir == "" {
		return diags
	}

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupRouted()
	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupPartial()

	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			continue
		}
		envBase := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
		bound := runAtEnv(routed, append(envBase, "BENCH_MODEL=claude-opus-4-8"), "bash", path, "--line probe prompt")
		if bound == nil || bound.ExitCode != 0 || !strings.Contains(bound.Stdout, "claude-opus-4-8") || !strings.Contains(bound.Stdout, "--line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass BENCH_MODEL and a dash-leading prompt to the harness in a routed repo", name))
		}
		unset := runAtEnv(routed, envBase, "bash", path, "line probe prompt")
		if unset != nil && unset.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
		}
		unbound := runAtEnv(routed, append(envBase, "BENCH_MODEL=claude-nonexistent-9"), "bash", path, "line probe prompt")
		if unbound != nil && unbound.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
		pass := runAtEnv(unrouted, envBase, "bash", path, "line probe prompt")
		if pass == nil || pass.ExitCode != 0 || !strings.Contains(pass.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass through in an unrouted repo", name))
		}
		explicit := runAtEnv(unrouted, append(envBase, "BENCH_MODEL=claude-anything-7"), "bash", path, "line probe prompt")
		if explicit == nil || explicit.ExitCode != 0 || !strings.Contains(explicit.Stdout, "claude-anything-7") || !strings.Contains(explicit.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass an explicit BENCH_MODEL through in an unrouted repo", name))
		}
		partialProbe := runAtEnv(partial, append(envBase, "BENCH_MODEL=claude-anything-7"), "bash", path, "line probe prompt")
		if partialProbe == nil || partialProbe.ExitCode != 0 || !strings.Contains(partialProbe.Stdout, "claude-anything-7") || !strings.Contains(partialProbe.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not fall back to passthrough on an incomplete binding", name))
		}
	}
	return diags
}

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
	if probe := runAtCleanEnv(root, "go", "test", "./..."); probe == nil || probe.ExitCode != 0 {
		diags = append(diags, "go test failed")
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

func checkSkillsIndex(root string) []string {
	var diags []string
	var expected []string
	skillFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "skills", "*", "SKILL.md"))
	sort.Strings(skillFiles)
	for _, file := range skillFiles {
		name := filepath.Base(filepath.Dir(file))
		if exists(filepath.Join(root, ".agents", "commands", name+".md")) {
			continue
		}
		index := frontmatterField(file, "index")
		if index == "" {
			diags = append(diags, fmt.Sprintf("skill '%s' missing index: frontmatter (the skills index is generated)", name))
			continue
		}
		line := fmt.Sprintf("- %s \u2192 `.agents/skills/%s/SKILL.md`", index, name)
		if note := frontmatterField(file, "index-note"); note != "" {
			line += " + " + note
		}
		expected = append(expected, line)
	}

	const start = "<!-- bench:skills-index:start -->"
	const end = "<!-- bench:skills-index:end -->"
	refPath := filepath.Join(root, ".bench", "BENCH-reference.md")
	ref := readIfExists(refPath)
	if ref == "" {
		return append(diags, ".bench/BENCH-reference.md missing (skills index unverifiable)")
	}
	actual, ok := markerBlock(ref, start, end)
	if !ok {
		return append(diags, ".bench/BENCH-reference.md skills-index markers missing (bench:skills-index)")
	}
	if strings.Join(expected, "\n") == strings.Join(actual, "\n") {
		return diags
	}

	expectedByName := map[string]string{}
	for _, line := range expected {
		if name := skillNameFromIndexLine(line); name != "" {
			expectedByName[name] = line
		}
	}
	actualByName := map[string]string{}
	for _, line := range actual {
		if name := skillNameFromIndexLine(line); name != "" {
			actualByName[name] = line
		}
	}
	attributed := false
	for name, line := range expectedByName {
		if actualByName[name] == line {
			continue
		}
		if _, ok := actualByName[name]; ok {
			diags = append(diags, fmt.Sprintf("skills index entry for '%s' drifted from its frontmatter (regenerate: .bench/skills-index.sh --write)", name))
		} else {
			diags = append(diags, fmt.Sprintf("skills index missing entry for skill '%s' (regenerate: .bench/skills-index.sh --write)", name))
		}
		attributed = true
	}
	for name := range actualByName {
		if _, ok := expectedByName[name]; !ok {
			diags = append(diags, fmt.Sprintf("skills index entry '%s' has no indexed .agents/skills/%s on disk (regenerate: .bench/skills-index.sh --write)", name, name))
			attributed = true
		}
	}
	if !attributed {
		diags = append(diags, "skills index block drifted from generated form (regenerate: .bench/skills-index.sh --write)")
	}
	return diags
}

func checkCommandGuideReferences(root string) []string {
	var diags []string
	guide := readIfExists(filepath.Join(root, ".bench", "BENCH.md")) + "\n" + readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md"))
	commandFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "commands", "*.md"))
	sort.Strings(commandFiles)
	for _, file := range commandFiles {
		name := strings.TrimSuffix(filepath.Base(file), ".md")
		if !strings.Contains(guide, "/"+name) {
			diags = append(diags, fmt.Sprintf("command '/%s' on disk but not referenced in the operating guide (.bench/BENCH.md or .bench/BENCH-reference.md)", name))
		}
	}
	return diags
}

func checkCodexCommandAdapters(root string) []string {
	var diags []string
	guide := readIfExists(filepath.Join(root, ".bench", "BENCH.md")) + "\n" + readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md"))
	commandFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "commands", "*.md"))
	sort.Strings(commandFiles)
	for _, file := range commandFiles {
		name := strings.TrimSuffix(filepath.Base(file), ".md")
		adapter := filepath.Join(root, ".agents", "skills", name, "SKILL.md")
		metadata := filepath.Join(root, ".agents", "skills", name, "agents", "openai.yaml")
		if !exists(adapter) {
			diags = append(diags, fmt.Sprintf("command '%s' has no Codex adapter skill at .agents/skills/%s/SKILL.md", name, name))
			continue
		}
		adapterText := readIfExists(adapter)
		if frontmatterField(adapter, "name") != name {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' frontmatter name does not match command", name))
		}
		if !strings.Contains(adapterText, ".agents/commands/"+name+".md") {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' does not reference .agents/commands/%s.md", name, name))
		}
		if !exists(metadata) {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' missing agents/openai.yaml explicit-invocation metadata", name))
			continue
		}
		if !strings.Contains(readIfExists(metadata), "allow_implicit_invocation: false") {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' does not disable implicit invocation", name))
		}
		if !strings.Contains(guide, "$"+name) {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' is not documented in the operating guide (.bench/BENCH.md or .bench/BENCH-reference.md)", name))
		}
	}
	return diags
}

func checkRoadmapPromotionAnchors(root string) []string {
	path := filepath.Join(root, ".agents", "commands", "bench-shape-idea.md")
	if !exists(path) {
		return nil
	}
	text := readIfExists(path)
	var diags []string
	if !strings.Contains(text, "ROADMAP.md") {
		diags = append(diags, "/bench-shape-idea does not reference ROADMAP.md (roadmap promotion seam)")
	}
	if !regexp.MustCompile(`(?i)remove|delete`).MatchString(text) {
		diags = append(diags, "/bench-shape-idea does not describe removing a promoted roadmap entry")
	}
	return diags
}

func markerBlock(text, start, end string) ([]string, bool) {
	lines := strings.Split(text, "\n")
	inBlock := false
	var out []string
	for _, line := range lines {
		switch {
		case line == end:
			return out, inBlock
		case inBlock:
			out = append(out, line)
		case line == start:
			inBlock = true
		}
	}
	return nil, false
}

func skillNameFromIndexLine(line string) string {
	re := regexp.MustCompile(`\.agents/skills/([a-z0-9-]+)/SKILL\.md`)
	m := re.FindStringSubmatch(line)
	if len(m) != 2 {
		return ""
	}
	return m[1]
}

func frontmatterField(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	fence := 0
	prefix := key + ":"
	for _, line := range strings.Split(string(data), "\n") {
		if line == "---" {
			fence++
			continue
		}
		if fence == 1 && strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
		if fence > 1 {
			return ""
		}
	}
	return ""
}

func readIfExists(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func anyContains(values []string, needle string) bool {
	for _, value := range values {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

func slashRel(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

func walkConformanceDocs(dir string) []string {
	var files []string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == "SKILL.md" ||
			strings.HasSuffix(name, ".md") ||
			strings.HasSuffix(name, ".yaml") ||
			strings.HasSuffix(name, ".yml") ||
			strings.HasSuffix(name, ".json") ||
			strings.HasSuffix(name, ".sh") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func uniqueSorted(values []string) []string {
	sort.Strings(values)
	if len(values) == 0 {
		return values
	}
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

func h2Headings(text string) map[string]bool {
	out := map[string]bool{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			out[line] = true
		}
	}
	return out
}

func collapseSpace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// runProbe captures stdout and stderr separately: probes like the npm pack
// JSON parse read stdout alone, and subprocess stderr chatter (npm's update
// notifier, warnings) must not corrupt it.
func runProbe(cmd *exec.Cmd, args []string) *Probe {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		if cmd.ProcessState != nil {
			exitCode = cmd.ProcessState.ExitCode()
		}
	}
	return &Probe{Args: append([]string(nil), args...), ExitCode: exitCode, Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
}

func runAt(dir string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	return runProbe(cmd, args)
}

func runAtCleanEnv(dir string, args ...string) *Probe {
	return runAtEnv(dir, conformanceSubprocessEnv(), args...)
}

func runAtEnv(dir string, env []string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = env
	return runProbe(cmd, args)
}

func runWithInput(dir, input string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(input)
	return runProbe(cmd, args)
}

func runWithInputEnv(dir string, env []string, input string, args ...string) *Probe {
	if len(args) == 0 {
		return nil
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = strings.NewReader(input)
	return runProbe(cmd, args)
}

func conformanceSubprocessEnv() []string {
	env := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_CONFORMANCE_ROOT=") {
			continue
		}
		env = append(env, kv)
	}
	return env
}

func wrapperStubDir(realBench string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "bench-wrapper-stub-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	content := "#!/usr/bin/env bash\nexec " + shellQuote(realBench) + " \"$@\"\n"
	if err := os.WriteFile(filepath.Join(dir, "bench"), []byte(content), 0o755); err != nil {
		cleanup()
		return "", func() {}, err
	}
	return dir, cleanup, nil
}

func adapterStubDir(realBench string) (string, func(), error) {
	dir, cleanup, err := wrapperStubDir(realBench)
	if err != nil {
		return "", cleanup, err
	}
	for _, name := range []string{"claude", "codex", "opencode"} {
		content := "#!/usr/bin/env bash\nprintf '%s\\n' \"$@\"\n"
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o755); err != nil {
			cleanup()
			return "", func() {}, err
		}
	}
	return dir, cleanup, nil
}

func tempGitRepoWithLines(linesEnv string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "bench-line-repo-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	if probe := runAt(dir, "git", "init", "-q"); probe == nil || probe.ExitCode != 0 {
		cleanup()
		return "", func() {}, fmt.Errorf("git init failed")
	}
	if err := os.MkdirAll(filepath.Join(dir, ".bench"), 0o755); err != nil {
		cleanup()
		return "", func() {}, err
	}
	if err := os.WriteFile(filepath.Join(dir, ".bench", "lines.env"), []byte(linesEnv), 0o644); err != nil {
		cleanup()
		return "", func() {}, err
	}
	return dir, cleanup, nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
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
