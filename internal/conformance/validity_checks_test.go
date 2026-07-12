package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

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
				Command string          `json:"command"`
				Timeout json.RawMessage `json:"timeout"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	var diags []string
	events := make([]string, 0, len(cfg.Hooks))
	for event := range cfg.Hooks {
		events = append(events, event)
	}
	sort.Strings(events)
	for _, event := range events {
		for _, group := range cfg.Hooks[event] {
			for _, hook := range group.Hooks {
				if len(hook.Timeout) > 0 {
					diags = append(diags, fmt.Sprintf("codex hooks.json %s hook sets a timeout key", event))
				}
			}
		}
	}
	var stopCommands, sessionCommands, preCommands []string
	for _, group := range cfg.Hooks["Stop"] {
		for _, hook := range group.Hooks {
			stopCommands = append(stopCommands, hook.Command)
		}
	}
	if !anyContains(stopCommands, ".bench/hooks/stop.sh") {
		diags = append(diags, "codex hooks.json Stop event does not run .bench/hooks/stop.sh")
	}
	for _, group := range cfg.Hooks["SessionStart"] {
		for _, hook := range group.Hooks {
			sessionCommands = append(sessionCommands, hook.Command)
		}
	}
	if !anyContains(sessionCommands, ".bench/hooks/session-start.sh") {
		diags = append(diags, "codex hooks.json SessionStart event does not run .bench/hooks/session-start.sh")
	}
	hasBashMatcher := false
	hasAgentMatcher := false
	for _, group := range cfg.Hooks["PreToolUse"] {
		if group.Matcher == "Bash" {
			hasBashMatcher = true
		}
		if group.Matcher == "Agent" {
			hasAgentMatcher = true
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
	if hasAgentMatcher || anyContains(preCommands, ".bench/hooks/check-agent-line.sh") {
		diags = append(diags, "codex hooks.json must not claim an Agent intent writer")
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

func TestRunConformanceDistinguishesAbsentAndEmptyInputs(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	assertStructuredPhaseDiags := func(label string, diags []string) {
		t.Helper()
		if !containsDiagnostic(diags, structuredPhaseUnavailable) {
			t.Errorf("%s shared rules did not fail closed with %q:\n%s", label, structuredPhaseUnavailable, strings.Join(diags, "\n"))
		}
	}

	absent := RunConformance(root, NewHarness(t).KitRoot)
	assertStructuredPhaseDiags("absent", absent)
	if !containsDiagnostic(absent, "JSON file missing: package.json") {
		t.Fatalf("absent package.json diagnostic missing:\n%s", strings.Join(absent, "\n"))
	}
	if !containsDiagnostic(absent, "lines.env missing: .bench/lines.env") {
		t.Fatalf("absent lines.env diagnostic missing:\n%s", strings.Join(absent, "\n"))
	}

	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".bench", "BENCH.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	empty := RunConformance(root, NewHarness(t).KitRoot)
	assertStructuredPhaseDiags("empty", empty)
	if !containsDiagnostic(empty, "lines.env tier unset: BENCH_TIER_TOP has no value") {
		t.Fatalf("empty lines.env diagnostic missing:\n%s", strings.Join(empty, "\n"))
	}
}

func TestStructuredPhaseContractIgnoresInactiveGuidance(t *testing.T) {
	const activeWithoutProgress = `# Bench Operating Guide

## How to talk to me

- **Structured Bench phase conversation:** Apply the named clauses
  ` + "`progress`, `exit`, `omission`, and `cohesion`" + ` proportionally.
  - **Exit:** Lead with the result.
  - **Omission:** Omit empty sections.
  - **Cohesion:** Keep related prose together.
`
	const progress = "A substantial update uses **Status:** and **Next:** labels."
	cases := map[string]string{
		"HTML comment":    "\n<!--   - **Progress:** " + progress + " -->\n\n## Workflow\n",
		"quotation":       "\n>   - **Progress:** " + progress + "\n\n## Workflow\n",
		"negated clause":  "\n  - **Progress:** Do not use **Status:** or **Next:** labels.\n\n## Workflow\n",
		"negating bullet": "\n- Do not follow this obsolete **Progress:** clause: " + progress + "\n\n## Workflow\n",
		"other section":   "\n## Workflow\n\n  - **Progress:** " + progress + "\n",
	}
	for name, decoy := range cases {
		t.Run(name, func(t *testing.T) {
			diags := checkStructuredPhaseContract(activeWithoutProgress + decoy)
			want := ".bench/BENCH.md dropped the structured Bench phase progress clause"
			if !containsDiagnostic(diags, want) {
				t.Fatalf("inactive guidance satisfied the progress clause; want %q in diagnostics:\n%s", want, strings.Join(diags, "\n"))
			}
		})
	}
}

func TestStructuredPhaseContractRejectsNegatedDeclaration(t *testing.T) {
	guide := `# Bench Operating Guide

## How to talk to me

- **Structured Bench phase conversation:** Do not apply the named clauses
  ` + "`progress`" + ` proportionally.
  - **Progress:** Use labels for substantial updates.
`
	diags := checkStructuredPhaseContract(guide)
	want := ".bench/BENCH.md negated or emptied the structured Bench phase contract declaration"
	if !containsDiagnostic(diags, want) {
		t.Fatalf("negated declaration stayed active; want %q in diagnostics:\n%s", want, strings.Join(diags, "\n"))
	}
}
