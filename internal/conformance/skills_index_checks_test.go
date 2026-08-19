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

	"github.com/gibbonmi/bench/internal/bounds"
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
		// Classified before it is opened: this is the first check to touch a skill file,
		// so a FIFO here would block the whole gate rather than one diagnostic.
		classified := bounds.ClassifyNoFollow(path)
		if classified.State.Failed() {
			diags = append(diags, fmt.Sprintf("%s refused: %s", slashRel(root, path), classified.Reason))
			continue
		}
		first := string(classified.Data)
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

// explicitOnlyDescriptionMarker is the phrasing that tells Codex a phase fires only on
// a typed invocation. On a row the policy marks implicitly invocable it is the one
// sentence that would keep the flip from ever firing, so it is graded as an anti-trigger.
const explicitOnlyDescriptionMarker = "Use only when the reviewer invokes"

// phaseInvocationPolicy is the reviewed invocation posture of every Bench phase: whether
// the Claude model may reach for the command on its own, and whether Codex may invoke the
// adapter skill implicitly. It restates the frontmatter facts it grades on purpose — that
// independence is what turns a silent flip of either surface red, the named exception to
// the one-source rule, and the canary fixtures beside it demonstrate the reds.
var phaseInvocationPolicy = map[string]struct {
	claudeModelInvocable bool
	codexImplicit        bool
}{
	// The reviewer's 2026-08-19 settle: the bug path fires on symptoms, on both harnesses.
	"bench-debug": {claudeModelInvocable: true, codexImplicit: true},

	"bench":                       {claudeModelInvocable: true},
	"bench-final-check":           {claudeModelInvocable: true},
	"bench-implement-spec":        {claudeModelInvocable: true},
	"bench-review-implementation": {claudeModelInvocable: true},
	"bench-shape-idea":            {claudeModelInvocable: true},
	"bench-write-spec":            {claudeModelInvocable: true},

	"bench-assess":     {},
	"bench-deepen":     {},
	"bench-drain":      {},
	"bench-setup-repo": {},
	"bench-update-kit": {},
	"bench-what-next":  {},
}

// invocationYAMLValue spells the openai.yaml line a policy value demands. Both the match
// and the mismatch branch read it, so the file's grammar has one spelling in the check.
func invocationYAMLValue(implicit bool) string {
	return fmt.Sprintf("allow_implicit_invocation: %t", implicit)
}

// frontmatterHasKey reports whether a frontmatter block declares key at all, empty value
// included — the presence question FrontmatterField cannot answer, since it returns the
// same empty string for an absent key and a declared-but-blank one. A mention in the body
// is not frontmatter and stays inert.
func frontmatterHasKey(text, key string) bool {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return false
	}
	for _, line := range lines[1:] {
		if line == "---" {
			return false
		}
		if strings.HasPrefix(line, key+":") {
			return true
		}
	}
	return false
}

func checkCodexCommandAdapters(root string) []string {
	var diags []string
	guide := readIfExists(filepath.Join(root, ".bench", "BENCH.md")) + "\n" + readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md"))
	commandFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "commands", "*.md"))
	sort.Strings(commandFiles)
	for _, file := range commandFiles {
		name := strings.TrimSuffix(filepath.Base(file), ".md")
		// The policy lookup runs ahead of every adapter check: a phase nobody declared a
		// trigger for reds as undeclared even when its adapter is missing too, so a new
		// command cannot arrive triggerless behind a louder diagnostic.
		policy, declared := phaseInvocationPolicy[name]
		if !declared {
			diags = append(diags, fmt.Sprintf("command '%s' has no invocation-policy row; every phase declares its Claude and Codex trigger before it ships", name))
			continue
		}
		if hasKey := frontmatterHasKey(readIfExists(file), "disable-model-invocation"); hasKey == policy.claudeModelInvocable {
			if hasKey {
				diags = append(diags, fmt.Sprintf("command '%s' frontmatter carries disable-model-invocation though the invocation policy makes it model-invocable on Claude", name))
			} else {
				diags = append(diags, fmt.Sprintf("command '%s' frontmatter declares no disable-model-invocation though the invocation policy disables it on Claude", name))
			}
		}
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
		commandName := name
		if name == "bench-what-next" {
			commandName = "bench-drain"
		}
		if !strings.Contains(adapterText, ".agents/commands/"+commandName+".md") {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' does not reference .agents/commands/%s.md", name, commandName))
		}
		if frontmatterHasKey(adapterText, "disable-model-invocation") {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' SKILL.md frontmatter carries the inert disable-model-invocation key; agents/openai.yaml is the Codex invocation-policy surface", name))
		}
		if !exists(metadata) {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' missing agents/openai.yaml explicit-invocation metadata", name))
			continue
		}
		metadataText := readIfExists(metadata)
		switch {
		case strings.Contains(metadataText, invocationYAMLValue(policy.codexImplicit)):
		case strings.Contains(metadataText, invocationYAMLValue(!policy.codexImplicit)):
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' agents/openai.yaml spells allow_implicit_invocation: %t against the invocation policy's %t", name, !policy.codexImplicit, policy.codexImplicit))
		default:
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' agents/openai.yaml declares no allow_implicit_invocation value (undeclared invocation policy)", name))
		}
		if policy.codexImplicit && strings.Contains(collapseSpace(frontmatterField(adapter, "description")), explicitOnlyDescriptionMarker) {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' description carries explicit-only phrasing (%q) though the invocation policy makes it implicitly invocable", name, explicitOnlyDescriptionMarker))
		}
		if !strings.Contains(guide, "$"+name) {
			diags = append(diags, fmt.Sprintf("Codex adapter '%s' is not documented in the operating guide (.bench/BENCH.md or .bench/BENCH-reference.md)", name))
		}
	}
	// The other completeness direction: a row survives the command file it grades.
	// The pass is unconditional over the whole table, so a narrow root — every canary
	// fixture in this family materializes one phase, not thirteen — collects a stale
	// row for each phase it does not carry. Scoping the pass to roots that look
	// complete would make the check guess which absences are real, so the noise stays
	// and the fixtures isolate their own red by its text rather than by the diagnostic
	// count.
	phases := make([]string, 0, len(phaseInvocationPolicy))
	for name := range phaseInvocationPolicy {
		phases = append(phases, name)
	}
	sort.Strings(phases)
	for _, name := range phases {
		if !exists(filepath.Join(root, ".agents", "commands", name+".md")) {
			diags = append(diags, fmt.Sprintf("invocation policy declares phase '%s' but .agents/commands/%s.md is absent (stale policy row)", name, name))
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

// TestSkillsIndexConformanceCarriesNoSecondReader keeps the skills index's derivable
// facts — marker text, allowlist path, line format — out of the conformance files.
// internal/skillsindex is their one source, so a copy in either file is a second reader
// by definition. Every declaration except this function is in scope, so hoisting a
// literal to a package-level const or a helper trips the guard rather than evading it.
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

// TestCodexAdapterInvocationPolicyGradesParseEdges pins the two policy edges no canary
// fixture can reach from the real tree: an agents/openai.yaml that declares nothing is
// red as undeclared rather than passing as "not the wrong value", and the explicit-only
// description clause is red on a row the policy marks implicitly invocable — the
// anti-trigger that would leave a flipped yaml unable to fire.
func TestCodexAdapterInvocationPolicyGradesParseEdges(t *testing.T) {
	adapterBody := func(name, description string) string {
		return fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\nRead `.agents/commands/%s.md`.\n", name, description, name)
	}
	for _, tc := range []struct {
		name        string
		phase       string
		description string
		yaml        string
		want        string
	}{
		{
			name:        "empty yaml",
			phase:       "bench-write-spec",
			description: "Explicit Codex adapter. Use only when the reviewer invokes $bench-write-spec.",
			yaml:        "",
			want:        "Codex adapter 'bench-write-spec' agents/openai.yaml declares no allow_implicit_invocation value (undeclared invocation policy)",
		},
		{
			name:        "yaml spelling neither value",
			phase:       "bench-write-spec",
			description: "Explicit Codex adapter. Use only when the reviewer invokes $bench-write-spec.",
			yaml:        "policy:\n  allow_implicit_invocation: maybe\n",
			want:        "Codex adapter 'bench-write-spec' agents/openai.yaml declares no allow_implicit_invocation value (undeclared invocation policy)",
		},
		{
			name:        "explicit-only description on an implicit row",
			phase:       "bench-debug",
			description: "Codex adapter for the bug path. Use only when the reviewer invokes $bench-debug.",
			yaml:        "policy:\n  allow_implicit_invocation: true\n",
			want:        `Codex adapter 'bench-debug' description carries explicit-only phrasing ("Use only when the reviewer invokes") though the invocation policy makes it implicitly invocable`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
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
			write(".bench/BENCH.md", fmt.Sprintf("/%s and $%s\n", tc.phase, tc.phase))
			write(".bench/BENCH-reference.md", "reference\n")
			write(".agents/commands/"+tc.phase+".md", "---\ndescription: phase\n---\n")
			write(".agents/skills/"+tc.phase+"/SKILL.md", adapterBody(tc.phase, tc.description))
			write(".agents/skills/"+tc.phase+"/agents/openai.yaml", tc.yaml)

			if diags := checkCodexCommandAdapters(root); !containsDiagnostic(diags, tc.want) {
				t.Fatalf("diagnostics = %q, want one containing %q", diags, tc.want)
			}
		})
	}
}

// TestCommandInvocationPolicyGradesTableCompleteness pins the two directions no canary
// fixture can reach: a policy row whose command file left the tree reds as stale (the
// table is compiled in, so omitting the file is the only way to age a row), and a
// command body that merely quotes the frontmatter key keeps its declared policy green.
func TestCommandInvocationPolicyGradesTableCompleteness(t *testing.T) {
	write := func(t *testing.T, root, rel, content string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// materializePolicyRoot writes a BASE-style tree for every declared phase except the
	// ones named, each phase carrying the frontmatter its own policy row demands.
	materializePolicyRoot := func(t *testing.T, omit ...string) string {
		t.Helper()
		omitted := map[string]bool{}
		for _, name := range omit {
			omitted[name] = true
		}
		root := t.TempDir()
		var guide strings.Builder
		for name, policy := range phaseInvocationPolicy {
			fmt.Fprintf(&guide, "/%s and $%s\n", name, name)
			if omitted[name] {
				continue
			}
			disable := ""
			if !policy.claudeModelInvocable {
				disable = "disable-model-invocation: true\n"
			}
			write(t, root, ".agents/commands/"+name+".md", fmt.Sprintf("---\ndescription: the %s phase\n%s---\n\nBody.\n", name, disable))
			target := name
			if name == "bench-what-next" {
				target = "bench-drain"
			}
			write(t, root, ".agents/skills/"+name+"/SKILL.md", fmt.Sprintf("---\nname: %s\ndescription: Codex adapter for %s.\n---\n\nRead `.agents/commands/%s.md`.\n", name, name, target))
			write(t, root, ".agents/skills/"+name+"/agents/openai.yaml", fmt.Sprintf("policy:\n  %s\n", invocationYAMLValue(policy.codexImplicit)))
		}
		write(t, root, ".bench/BENCH.md", guide.String())
		write(t, root, ".bench/BENCH-reference.md", "reference\n")
		return root
	}

	t.Run("intact policy root is green", func(t *testing.T) {
		if diags := checkCodexCommandAdapters(materializePolicyRoot(t)); len(diags) != 0 {
			t.Fatalf("diagnostics over an intact policy root = %q, want none", diags)
		}
	})

	t.Run("stale table row", func(t *testing.T) {
		diags := checkCodexCommandAdapters(materializePolicyRoot(t, "bench-debug"))
		const want = "invocation policy declares phase 'bench-debug' but .agents/commands/bench-debug.md is absent (stale policy row)"
		if !containsDiagnostic(diags, want) {
			t.Fatalf("diagnostics = %q, want one containing %q", diags, want)
		}
		for _, diag := range diags {
			if strings.Contains(diag, "bench-write-spec") {
				t.Fatalf("stale-row grading also red on a present phase: %q", diag)
			}
		}
	})

	t.Run("body prose mention of the key is inert", func(t *testing.T) {
		root := materializePolicyRoot(t)
		// bench-write-spec is model-invocable, so a graded key would red the phase; the
		// mention lands in the body, below the closing frontmatter fence.
		write(t, root, ".agents/commands/bench-write-spec.md",
			"---\ndescription: the bench-write-spec phase\n---\n\nA command file disables Claude's own reach for a phase with the\n`disable-model-invocation: true` frontmatter key.\n")
		if diags := checkCodexCommandAdapters(root); len(diags) != 0 {
			t.Fatalf("body-prose mention of disable-model-invocation reddened the check: %q", diags)
		}
	})

	t.Run("undeclared phase reds ahead of its missing adapter", func(t *testing.T) {
		root := materializePolicyRoot(t)
		write(t, root, ".agents/commands/ghost.md", "---\ndescription: undeclared\n---\n")
		diags := checkCodexCommandAdapters(root)
		const want = "command 'ghost' has no invocation-policy row; every phase declares its Claude and Codex trigger before it ships"
		if !containsDiagnostic(diags, want) {
			t.Fatalf("diagnostics = %q, want one containing %q", diags, want)
		}
		for _, diag := range diags {
			if strings.Contains(diag, "command 'ghost' has no Codex adapter skill") {
				t.Fatalf("adapter-existence grading preempted the undeclared-policy red: %q", diag)
			}
		}
	})
}
