package conformance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/skillsindex"
)

// TestClaudeSkillMirrorClassifiesStandaloneSkillSymlinks pins the three
// non-craft postures of checkClaudeSkillMirror: a symlink resolving to its own
// .agents/skills/<name>/SKILL.md is a genuine standalone skill and passes; a
// name shared with an .agents/commands file is a phase adapter and stays red;
// anything else — a dangling link, a plain directory — stays red as well.
func TestClaudeSkillMirrorClassifiesStandaloneSkillSymlinks(t *testing.T) {
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
		".claude/skills/plaindir is neither a craft skill nor a standalone skill symlink into .agents/skills (phase adapters are Codex-only; it duplicates the slash menu)",
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
		// A non-craft entry is admissible only as a standalone skill: a symlink
		// resolving to its own .agents/skills/<name>/SKILL.md. A name shared with a
		// command file is a phase adapter regardless of where its link points.
		if exists(filepath.Join(root, ".agents", "commands", base+".md")) {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s is not a craft skill (phase adapters are Codex-only; it duplicates the slash menu)", base))
			continue
		}
		linkInfo, err := os.Lstat(dir)
		if err != nil || linkInfo.Mode()&os.ModeSymlink == 0 {
			diags = append(diags, fmt.Sprintf(".claude/skills/%s is neither a craft skill nor a standalone skill symlink into .agents/skills (phase adapters are Codex-only; it duplicates the slash menu)", base))
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

// checkSkillsIndex grades the committed skills index through the module that
// generates it, so the gate's oracle and `bench skills-index` cannot disagree.
func checkSkillsIndex(root string) []string {
	return skillsindex.Check(root)
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

// TestSkillsIndexConformanceCarriesNoSecondReader is the row that sees the cheapest
// wrong refactor: adding internal/skillsindex and leaving the old parsers in place.
// The three banned literals are the skills index's derivable facts — marker text,
// allowlist path, line format — and the module is their one source, so a copy in
// either conformance file is a second reader by definition. Every declaration except
// this function is in scope, so hoisting a literal to a package-level const or a
// helper trips the guard rather than evading it.
func TestSkillsIndexConformanceCarriesNoSecondReader(t *testing.T) {
	const guard = "TestSkillsIndexConformanceCarriesNoSecondReader"
	banned := []string{
		"<!-- bench:skills-index:start -->",
		"consumer-payload.json",
		"→ `.agents/skills/",
	}
	kitRoot := NewHarness(t).KitRoot
	for _, rel := range []string{"skills_index_checks_test.go", "checks_test.go"} {
		fset := token.NewFileSet()
		path := filepath.Join(kitRoot, "internal", "conformance", rel)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, declaration := range file.Decls {
			if function, ok := declaration.(*ast.FuncDecl); ok && function.Name.Name == guard {
				continue
			}
			ast.Inspect(declaration, func(node ast.Node) bool {
				literal, ok := node.(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING {
					return true
				}
				value, err := strconv.Unquote(literal.Value)
				if err != nil {
					return true
				}
				for _, banned := range banned {
					if strings.Contains(value, banned) {
						t.Errorf("%s carries skills-index literal %q (%s); internal/skillsindex is its one source", rel, banned, fset.Position(literal.Pos()))
					}
				}
				return true
			})
		}
	}
}
