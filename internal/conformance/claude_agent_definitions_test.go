package conformance

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/skillsindex"
)

const (
	// claudeAgentsDir is the Claude adapter's agents surface. It is a real directory and
	// not a symlink into the portable tree, because Codex has no agents surface to share.
	claudeAgentsDir = ".claude/agents"
	// claudeAgentPrefix marks the files this check owns. A consumer's own agent files sit
	// beside the Bench ones and stay ungraded.
	claudeAgentPrefix = "bench-"
	// claudeAgentSkill states the routing rule. The check reads it in both directions, so
	// neither a shipped agent nobody is told to use nor a rule naming a missing type can
	// pass.
	claudeAgentSkill = ".agents/skills/bench-craft-delegate/SKILL.md"
)

// claudeAgentForbiddenTools are the tool schemas whose absence is the whole point of a
// Bench agent type. A delegate that can spawn, publish, or block the user is not the cheap
// read-only or fenced-write role the routing rule assigns it.
var claudeAgentForbiddenTools = []string{"Agent", "Artifact", "AskUserQuestion"}

// claudeAgentRequiredTools are the two tools every Bench delegate needs. Read is how it
// sees its charge's inputs, and Bash is how it runs a worktree exec.
var claudeAgentRequiredTools = []string{"Read", "Bash"}

// claudeAgentSkillName matches a backticked agent basename in the routing rule. A Bench
// agent basename carries no space, so the pattern separates a named type from the skill's
// backticked commands and paths. The prefix comes from the constant the file filter also
// reads, so the naming convention has one source.
var claudeAgentSkillName = regexp.MustCompile("`(" + regexp.QuoteMeta(claudeAgentPrefix) + "[a-z0-9]+(?:-[a-z0-9]+)*)`")

// checkClaudeAgentDefinitions grades every Bench agent file in the Claude adapter. It
// grades each file against the shape a cheap delegate needs. It then reconciles the
// shipped set with the routing rule the delegate skill states.
func checkClaudeAgentDefinitions(root string) []string {
	files, diags := claudeAgentFiles(root)
	named := claudeAgentNamesInSkill(root)
	for _, rel := range files {
		diags = append(diags, claudeAgentFileDiagnostics(root, rel, named)...)
	}
	shipped := map[string]bool{}
	for _, rel := range files {
		shipped[strings.TrimSuffix(path.Base(rel), ".md")] = true
	}
	for _, name := range named {
		if !shipped[name] {
			diags = append(diags, "claude-agent missing: "+claudeAgentSkill+" names '"+name+"', which has no file in "+claudeAgentsDir)
		}
	}
	return diags
}

// claudeAgentFileDiagnostics grades one agent file. Each fault reports and the file keeps
// going. One file can carry more than one fault, and a reader who repairs it wants the
// whole list.
func claudeAgentFileDiagnostics(root, rel string, named []string) []string {
	base := strings.TrimSuffix(path.Base(rel), ".md")
	full := filepath.Join(root, filepath.FromSlash(rel))
	// The shared producer classifier owns shape and bounded bytes. This check owns only
	// what each state costs the agent report.
	subject := bounds.ClassifyNoFollow(full)
	if subject.State != bounds.StateParsed {
		return []string{"claude-agent subject refused: " + rel + " is not a readable regular file (" + subject.Reason + ")"}
	}
	var diags []string
	if name := skillsindex.FrontmatterField(full, "name"); name != base {
		diags = append(diags, "claude-agent name mismatch: "+rel+" declares name '"+name+"', want '"+base+"'")
	}
	if model := skillsindex.FrontmatterField(full, "model"); model != "" {
		diags = append(diags, "claude-agent model declared: "+rel+" declares model '"+model+"'; a Bench agent type declares no model, so the charge carries the bound tier token")
	}
	tools := claudeAgentTools(skillsindex.FrontmatterField(full, "tools"))
	if len(tools) == 0 {
		diags = append(diags, "claude-agent tools absent: "+rel+" declares no tools, so it inherits the full tool set")
	} else {
		for _, tool := range claudeAgentForbiddenTools {
			if tools[tool] {
				diags = append(diags, "claude-agent tool refused: "+rel+" lists '"+tool+"', which a Bench delegate must not carry")
			}
		}
		for _, tool := range claudeAgentRequiredTools {
			if !tools[tool] {
				diags = append(diags, "claude-agent tool missing: "+rel+" lists no '"+tool+"' tool")
			}
		}
	}
	found := false
	for _, name := range named {
		if name == base {
			found = true
		}
	}
	if !found {
		diags = append(diags, "claude-agent unnamed in skill: "+base+" ships but "+claudeAgentSkill+" never names it")
	}
	return diags
}

// claudeAgentFiles lists the Bench agent files in stable order. An absent agents directory
// yields no files and no diagnostic of its own: the routing rule names both types, so the
// missing-agent reconciliation reports each one by name instead.
func claudeAgentFiles(root string) (files, diags []string) {
	dir := filepath.Join(root, filepath.FromSlash(claudeAgentsDir))
	info, err := os.Lstat(dir)
	switch {
	case err != nil:
		return nil, nil
	case info.Mode()&os.ModeSymlink != 0:
		return nil, []string{"claude-agent subject refused: " + claudeAgentsDir + " is a symbolic link, not a regular directory"}
	case !info.IsDir():
		return nil, []string{"claude-agent subject refused: " + claudeAgentsDir + " is not a directory"}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, []string{"claude-agent subject unreadable: " + claudeAgentsDir + ": " + err.Error()}
	}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, claudeAgentPrefix) || !strings.HasSuffix(name, ".md") {
			continue
		}
		files = append(files, path.Join(claudeAgentsDir, name))
	}
	sort.Strings(files)
	return files, nil
}

// claudeAgentNamesInSkill reads the routing rule's backticked agent basenames in sorted
// order. An absent or refused skill file yields no names, and the unnamed-in-skill
// diagnostic then reports every shipped agent.
func claudeAgentNamesInSkill(root string) []string {
	subject := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(claudeAgentSkill)))
	if subject.State != bounds.StateParsed {
		return nil
	}
	seen := map[string]bool{}
	var names []string
	for _, match := range claudeAgentSkillName.FindAllStringSubmatch(string(subject.Data), -1) {
		if !seen[match[1]] {
			seen[match[1]] = true
			names = append(names, match[1])
		}
	}
	sort.Strings(names)
	return names
}

// claudeAgentTools splits the frontmatter's comma-separated tool list. The frontmatter
// reader returns the first value line, so a list folded onto a second line contributes
// only its first entries and the required-tool diagnostics then report the rest.
func claudeAgentTools(value string) map[string]bool {
	tools := map[string]bool{}
	for _, field := range strings.Split(value, ",") {
		if field = strings.TrimSpace(field); field != "" {
			tools[field] = true
		}
	}
	return tools
}

// TestClaudeAgentDefinitionsIgnoresAConsumerPrefix pins the ownership boundary. A linked
// repo keeps its own agent files in the same directory, and grading them would export this
// kit's tool policy to a consumer's agents.
func TestClaudeAgentDefinitionsIgnoresAConsumerPrefix(t *testing.T) {
	root := t.TempDir()
	writeClaudeAgentFixture(t, root, "team-helper.md", "---\nname: something-else\nmodel: opus\n---\n")
	writeClaudeAgentFixture(t, root, "bench-reviewer.md", "---\nname: bench-reviewer\ntools: Read, Bash\n---\n")
	writeClaudeAgentSkill(t, root, "The axis runs as `bench-reviewer`.\n")

	for _, diagnostic := range checkClaudeAgentDefinitions(root) {
		if strings.Contains(diagnostic, "team-helper") {
			t.Fatalf("check graded a consumer agent file: %s", diagnostic)
		}
	}
}

// TestClaudeAgentDefinitionsReadTheSkillBothWays pins the two reconciliation directions in
// one tree. A check that reads only the files ships a dead type, and one that reads only
// the skill sends the charge back to the general-purpose type.
func TestClaudeAgentDefinitionsReadTheSkillBothWays(t *testing.T) {
	root := t.TempDir()
	writeClaudeAgentFixture(t, root, "bench-orphan.md", "---\nname: bench-orphan\ntools: Read, Bash\n---\n")
	writeClaudeAgentSkill(t, root, "The write delegation runs as `bench-writer`.\n")

	diagnostics := strings.Join(checkClaudeAgentDefinitions(root), "\n")
	for _, want := range []string{
		"claude-agent unnamed in skill: bench-orphan",
		"claude-agent missing: " + claudeAgentSkill + " names 'bench-writer'",
	} {
		if !strings.Contains(diagnostics, want) {
			t.Errorf("diagnostics missing %q:\n%s", want, diagnostics)
		}
	}
}

// TestClaudeAgentDefinitionsRefuseANonRegularSubject pins the refusal path. A symbolic link
// where the agents directory belongs would grade whatever it points at under the adapter's
// canonical path.
func TestClaudeAgentDefinitionsRefuseANonRegularSubject(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(root, ".claude", "agents")); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
	diagnostics := strings.Join(checkClaudeAgentDefinitions(root), "\n")
	if !strings.Contains(diagnostics, "claude-agent subject refused: "+claudeAgentsDir+" is a symbolic link") {
		t.Fatalf("diagnostics do not refuse a symlinked agents directory:\n%s", diagnostics)
	}
}

// TestClaudeAgentDefinitionsGradeEveryPolicedTool drives every policed tool by an
// independently authored name. The independence is load-bearing: a test that ranges over
// the production slice passes after the slice loses a member, so a trim to {"Agent"} runs
// silent and hands a cheap delegate the publish or block tool. The equality assertion
// below is what turns that trim red.
func TestClaudeAgentDefinitionsGradeEveryPolicedTool(t *testing.T) {
	forbidden := []string{"Agent", "Artifact", "AskUserQuestion"}
	required := []string{"Read", "Bash"}
	if !slices.Equal(claudeAgentForbiddenTools, forbidden) {
		t.Fatalf("claudeAgentForbiddenTools = %v, want %v", claudeAgentForbiddenTools, forbidden)
	}
	if !slices.Equal(claudeAgentRequiredTools, required) {
		t.Fatalf("claudeAgentRequiredTools = %v, want %v", claudeAgentRequiredTools, required)
	}
	for _, tool := range forbidden {
		t.Run("forbidden/"+tool, func(t *testing.T) {
			root := t.TempDir()
			writeClaudeAgentFixture(t, root, "bench-reviewer.md", "---\nname: bench-reviewer\ntools: Read, Bash, "+tool+"\n---\n")
			writeClaudeAgentSkill(t, root, "The axis runs as `bench-reviewer`.\n")

			want := "claude-agent tool refused: " + claudeAgentsDir + "/bench-reviewer.md lists '" + tool + "'"
			if diagnostics := strings.Join(checkClaudeAgentDefinitions(root), "\n"); !strings.Contains(diagnostics, want) {
				t.Fatalf("diagnostics missing %q:\n%s", want, diagnostics)
			}
		})
	}
	for _, tool := range required {
		t.Run("required/"+tool, func(t *testing.T) {
			root := t.TempDir()
			var kept []string
			for _, other := range required {
				if other != tool {
					kept = append(kept, other)
				}
			}
			writeClaudeAgentFixture(t, root, "bench-reviewer.md", "---\nname: bench-reviewer\ntools: Glob, "+strings.Join(kept, ", ")+"\n---\n")
			writeClaudeAgentSkill(t, root, "The axis runs as `bench-reviewer`.\n")

			want := "claude-agent tool missing: " + claudeAgentsDir + "/bench-reviewer.md lists no '" + tool + "' tool"
			if diagnostics := strings.Join(checkClaudeAgentDefinitions(root), "\n"); !strings.Contains(diagnostics, want) {
				t.Fatalf("diagnostics missing %q:\n%s", want, diagnostics)
			}
		})
	}
}

// TestClaudeAgentDefinitionsRefuseAnUnreadableAgentFile pins the per-file refusal. It is a
// different code path from the directory refusal: a symbolic link at an agent file would
// otherwise skip every other diagnostic for that file and report nothing.
func TestClaudeAgentDefinitionsRefuseAnUnreadableAgentFile(t *testing.T) {
	root := t.TempDir()
	writeClaudeAgentFixture(t, root, "bench-writer.md", "---\nname: bench-writer\ntools: Read, Bash\n---\n")
	link := filepath.Join(root, filepath.FromSlash(claudeAgentsDir), "bench-reviewer.md")
	if err := os.Symlink(filepath.Join(root, "elsewhere.md"), link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
	writeClaudeAgentSkill(t, root, "The axis runs as `bench-reviewer`, and the write runs as `bench-writer`.\n")

	want := "claude-agent subject refused: " + claudeAgentsDir + "/bench-reviewer.md is not a readable regular file"
	if diagnostics := strings.Join(checkClaudeAgentDefinitions(root), "\n"); !strings.Contains(diagnostics, want) {
		t.Fatalf("diagnostics missing %q:\n%s", want, diagnostics)
	}
}

// TestClaudeAgentDefinitionsRefuseAPlainFileWhereTheDirectoryBelongs pins the third subject
// state. A regular file at the adapter's agents path is neither absent nor a link, and an
// unrefused one would enumerate nothing and report a clean adapter.
func TestClaudeAgentDefinitionsRefuseAPlainFileWhereTheDirectoryBelongs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(claudeAgentsDir)), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	want := "claude-agent subject refused: " + claudeAgentsDir + " is not a directory"
	if diagnostics := strings.Join(checkClaudeAgentDefinitions(root), "\n"); !strings.Contains(diagnostics, want) {
		t.Fatalf("diagnostics missing %q:\n%s", want, diagnostics)
	}
}

func writeClaudeAgentFixture(t *testing.T, root, name, body string) {
	t.Helper()
	dir := filepath.Join(root, filepath.FromSlash(claudeAgentsDir))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeClaudeAgentSkill(t *testing.T, root, body string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(claudeAgentSkill))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
