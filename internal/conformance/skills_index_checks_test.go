package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	kitpayload "github.com/gibbonmi/bench"
)

// TestClaudeSkillMirrorClassifiesUserInvokedSkillSymlinks pins the three
// non-craft postures of checkClaudeSkillMirror: a symlink resolving to its own
// .agents/skills/<name>/SKILL.md is a genuine user-invoked skill and passes; a
// name shared with an .agents/commands file is a phase adapter and stays red;
// anything else — a dangling link, a plain directory — stays red as well.
func TestClaudeSkillMirrorClassifiesUserInvokedSkillSymlinks(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	link := func(target, rel string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
	}

	write(".agents/skills/prototype/SKILL.md", "---\nname: prototype\n---\n")
	link("../../.agents/skills/prototype", ".claude/skills/prototype")

	write(".agents/commands/bench-write-spec.md", "command\n")
	write(".agents/skills/bench-write-spec/SKILL.md", "---\nname: bench-write-spec\n---\n")
	link("../../.agents/skills/bench-write-spec", ".claude/skills/bench-write-spec")

	link("../../.agents/skills/ghost", ".claude/skills/ghost")

	write(".claude/skills/plaindir/SKILL.md", "---\nname: plaindir\n---\n")
	write(".agents/skills/plaindir/SKILL.md", "---\nname: plaindir\n---\n")

	got := checkClaudeSkillMirror(root)
	sort.Strings(got)
	want := []string{
		".claude/skills/bench-write-spec is not a craft skill (phase adapters are Codex-only; it duplicates the slash menu)",
		".claude/skills/ghost does not resolve to .agents/skills/ghost/SKILL.md (broken adapter link)",
		".claude/skills/plaindir is neither a craft skill nor a user-invoked skill symlink into .agents/skills (phase adapters are Codex-only; it duplicates the slash menu)",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("diagnostics = \n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func checkSkillsIndexAndCommandAdapters(root string) []string {
	var diags []string
	diags = append(diags, checkSkillsIndex(root)...)
	diags = append(diags, checkCommandGuideReferences(root)...)
	diags = append(diags, checkCodexCommandAdapters(root)...)
	diags = append(diags, checkRoadmapPromotionAnchors(root)...)
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
		base := filepath.Base(dir)
		if strings.HasPrefix(base, "bench-craft-") {
			info, err := os.Stat(dir)
			if err != nil || !info.IsDir() {
				continue
			}
			if !exists(filepath.Join(dir, "SKILL.md")) {
				diags = append(diags, fmt.Sprintf(".claude/skills/%s does not resolve to a SKILL.md (broken adapter link)", base))
			}
			continue
		}
		// A non-craft entry is admissible only as a user-invoked skill: a symlink
		// resolving to its own .agents/skills/<name>/SKILL.md. A name shared with a
		// command file is a phase adapter regardless of where its link points.
		if exists(filepath.Join(root, ".agents", "commands", base+".md")) {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s is not a craft skill (phase adapters are Codex-only; it duplicates the slash menu)", base))
			continue
		}
		linkInfo, err := os.Lstat(dir)
		if err != nil || linkInfo.Mode()&os.ModeSymlink == 0 {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s is neither a craft skill nor a user-invoked skill symlink into .agents/skills (phase adapters are Codex-only; it duplicates the slash menu)", base))
			continue
		}
		resolved, resolveErr := filepath.EvalSymlinks(dir)
		canonical, canonicalErr := filepath.EvalSymlinks(filepath.Join(root, ".agents", "skills", base))
		if resolveErr != nil || canonicalErr != nil || resolved != canonical || !exists(filepath.Join(canonical, "SKILL.md")) {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s does not resolve to .agents/skills/%s/SKILL.md (broken adapter link)", base, base))
		}
	}
	return diags
}

// kitOnlySkillSources reads root's consumer-payload allowlist and returns the skill
// sources it withholds. The allowlist is the one source of who receives an asset, so
// the index expectation reads it here exactly as the generator does instead of naming
// the kit-only skills a second time. A tree with no allowlist (a stripped fixture) has
// nothing withheld, which is also what the generator concludes.
func kitOnlySkillSources(root string) map[string]bool {
	var rows []kitpayload.PayloadRow
	if err := json.Unmarshal([]byte(readIfExists(filepath.Join(root, ".bench", "consumer-payload.json"))), &rows); err != nil {
		return nil
	}
	out := map[string]bool{}
	for _, source := range kitpayload.PayloadKitOnlyPrefixes(rows) {
		out[source] = true
	}
	return out
}

func checkSkillsIndex(root string) []string {
	var diags []string
	var expected []string
	kitOnly := kitOnlySkillSources(root)
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
		if kitOnly[".agents/skills/"+name] {
			line += " (kit-only)"
		}
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
	if !strings.Contains(collapseSpace(text), "working prioritization document") {
		diags = append(diags, "/bench-shape-idea does not describe ROADMAP.md as the working prioritization document")
	}
	if !strings.Contains(collapseSpace(text), "roadmap row in place") {
		diags = append(diags, "/bench-shape-idea does not state that pulling an item leaves its roadmap row in place (row presence is status)")
	}
	if strings.Contains(collapseSpace(text), "capture-and-forget") {
		diags = append(diags, "/bench-shape-idea reintroduces capture-sink roadmap wording; capture/IDEAS.md is the capture inbox and ROADMAP.md the working plan")
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
	// The allowlist is JSON, so its key order carries no meaning: this fixture writes
	// "audience" before "source" precisely because a reader that only matches one fixed
	// order would drop the marker and generate a consumer-visible skill for a withheld
	// one. Marker text is not re-stated here — the assertion below reads the generated
	// line as a whole.
	if err := os.WriteFile(filepath.Join(tmp, ".bench", "consumer-payload.json"),
		[]byte(`[{ "audience": "kit-only", "mode": "0644", "tree": true, "source": ".agents/skills/zeta-skill" }]`), 0o644); err != nil {
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
	if !strings.Contains(generated, "- doing zeta things \u2192 `.agents/skills/zeta-skill/SKILL.md` (kit-only)") {
		return []string{"skills-index generate/verify contract failed: --write did not generate the entry from frontmatter and the allowlist's audience"}
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
