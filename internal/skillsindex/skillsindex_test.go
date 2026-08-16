package skillsindex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFile is the fixture builder every case below shares: a temp root grown one
// tracked-path-shaped file at a time, so the tests drive the package through its
// public seam rather than through the fence scan's internals.
func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func reference(block string) string {
	return "# Reference\n\n<!-- bench:skills-index:start -->\n" + block + "<!-- bench:skills-index:end -->\n"
}

func renderedBlock(t *testing.T, root string) string {
	t.Helper()
	entries, err := Entries(root)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, entry := range entries {
		if entry.Indexed() {
			lines = append(lines, entry.Line())
		}
	}
	return strings.Join(lines, "\n")
}

func TestEntriesRenderAlphabeticallySkippingCommandAdapters(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\nindex-note: and the alpha note\n---\n")
	writeFile(t, root, ".agents/skills/gamma/SKILL.md", "---\nname: gamma\nindex: doing gamma things\n---\n")
	// A skill sharing a name with a command file is a phase adapter: the slash menu
	// already surfaces it, so it must contribute no line.
	writeFile(t, root, ".agents/commands/bench-write-spec.md", "command\n")
	writeFile(t, root, ".agents/skills/bench-write-spec/SKILL.md", "---\nname: bench-write-spec\nindex: writing specs\n---\n")
	writeFile(t, root, ".agents/skills/zeta/SKILL.md", "---\nname: zeta\nindex: doing zeta things\n---\n")
	// The allowlist is JSON, so its key order carries no meaning: this row orders
	// "audience" before "source" precisely because a reader matching one fixed order
	// would drop the marker and advertise a withheld skill to consumers.
	writeFile(t, root, ".bench/consumer-payload.json",
		`[{ "audience": "kit-only", "mode": "0644", "tree": true, "source": ".agents/skills/zeta" }]`)

	want := "- doing alpha things → `.agents/skills/alpha/SKILL.md` + and the alpha note\n" +
		"- doing gamma things → `.agents/skills/gamma/SKILL.md`\n" +
		"- doing zeta things → `.agents/skills/zeta/SKILL.md` (kit-only)"
	if got := renderedBlock(t, root); got != want {
		t.Fatalf("rendered block =\n%s\nwant\n%s", got, want)
	}
}

func TestFrontmatterFieldReadsOnlyTheLeadingFence(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "SKILL.md")
	writeFile(t, root, "SKILL.md", "---\nname: probe\nindex: first value\nindex: second value\n---\nindex-after: past the fence\n")

	if got, want := FrontmatterField(path, "index"), "first value"; got != want {
		t.Fatalf("first value: index = %q, want %q", got, want)
	}
	if got := FrontmatterField(path, "index-after"); got != "" {
		t.Fatalf("key past the closing fence read as %q, want empty", got)
	}
	if got := FrontmatterField(path, "index-note"); got != "" {
		t.Fatalf("absent index-note read as %q, want empty", got)
	}
}

func TestCheckAndWriteGenerateVerifyContract(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	writeFile(t, root, ".agents/skills/zeta/SKILL.md", "---\nname: zeta\nindex: doing zeta things\n---\n")
	writeFile(t, root, ".bench/consumer-payload.json",
		`[{ "audience": "kit-only", "mode": "0644", "tree": true, "source": ".agents/skills/zeta" }]`)
	refPath := filepath.Join(root, ".bench", "BENCH-reference.md")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))

	want := []string{
		"skills index missing entry for skill 'alpha' (regenerate: bench skills-index --write)",
		"skills index missing entry for skill 'zeta' (regenerate: bench skills-index --write)",
	}
	if got := Check(root); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("check on an empty block =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if err := Write(root); err != nil {
		t.Fatal(err)
	}
	if got := Check(root); len(got) != 0 {
		t.Fatalf("check right after write = %s, want none", strings.Join(got, "\n"))
	}
	generated, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	wantRef := reference("- doing alpha things → `.agents/skills/alpha/SKILL.md`\n" +
		"- doing zeta things → `.agents/skills/zeta/SKILL.md` (kit-only)\n")
	if string(generated) != wantRef {
		t.Fatalf("generated reference =\n%q\nwant\n%q", generated, wantRef)
	}
	if err := Write(root); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(generated) {
		t.Fatalf("second write changed bytes:\n%q\nwant\n%q", after, generated)
	}
	info, err := os.Stat(refPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o644 {
		t.Fatalf("reference mode = %o, want 644", got)
	}
}

func TestCheckAttributesEveryDriftInAlphabeticalOrder(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\n---\n")
	writeFile(t, root, ".agents/skills/beta/SKILL.md", "---\nname: beta\nindex: doing beta things\n---\n")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(
		"- stale beta wording \u2192 `.agents/skills/beta/SKILL.md`\n"+
			"- doing gamma things \u2192 `.agents/skills/gamma/SKILL.md`\n"))

	want := []string{
		"skill 'alpha' missing index: frontmatter (the skills index is generated)",
		"skills index entry for 'beta' drifted from its frontmatter (regenerate: bench skills-index --write)",
		"skills index entry 'gamma' has no indexed .agents/skills/gamma on disk (regenerate: bench skills-index --write)",
	}
	if got := Check(root); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("attributed diagnostics =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// A skill can earn two diagnostics at once: it declares no trigger and still carries a
// committed entry. Keying one string per skill would drop the first of the pair, and no
// canary fixture pairs a field-less skill with a non-empty committed block.
func TestCheckKeepsBothDiagnosticsForACollidingSkill(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\n---\n")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(
		"- doing alpha things \u2192 `.agents/skills/alpha/SKILL.md`\n"))

	want := []string{
		"skill 'alpha' missing index: frontmatter (the skills index is generated)",
		"skills index entry 'alpha' has no indexed .agents/skills/alpha on disk (regenerate: bench skills-index --write)",
	}
	if got := Check(root); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("colliding diagnostics =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestZeroSkillRootRendersAnEmptyBlockAndPasses(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
	if err := Write(root); err != nil {
		t.Fatal(err)
	}
	generated, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(generated) != reference("") {
		t.Fatalf("zero-skill reference =\n%q\nwant\n%q", generated, reference(""))
	}
	if got := Check(root); len(got) != 0 {
		t.Fatalf("check on a zero-skill root = %s, want none", strings.Join(got, "\n"))
	}
}

func TestUnparseableAllowlistRefusesWriteAndLeavesCheckAlone(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/zeta/SKILL.md", "---\nname: zeta\nindex: doing zeta things\n---\n")
	writeFile(t, root, ".bench/consumer-payload.json", "{not json")
	refPath := filepath.Join(root, ".bench", "BENCH-reference.md")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
	before, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}

	err = Write(root)
	want := ".bench/consumer-payload.json unreadable: kit-only marking unresolved (write refused)"
	if err == nil || err.Error() != want {
		t.Fatalf("write refusal = %v, want %q", err, want)
	}
	after, readErr := os.ReadFile(refPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(after) != string(before) {
		t.Fatalf("refused write changed bytes:\n%q\nwant\n%q", after, before)
	}
	// Check marks nothing and adds no allowlist diagnostic, so the gate gains no red
	// the index itself did not earn.
	wantCheck := "skills index missing entry for skill 'zeta' (regenerate: bench skills-index --write)"
	if got := Check(root); strings.Join(got, "\n") != wantCheck {
		t.Fatalf("check with an unparseable allowlist =\n%s\nwant\n%s", strings.Join(got, "\n"), wantCheck)
	}
}

// hostileRoot is a repository root whose own name carries every character a glob would
// consume. Discovery must read it as literal text: a pattern-based walk matches nothing
// here and silently reports a repository with no skills at all.
func hostileRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "hostile [ ] * ? root")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestEntriesEnumerateLiterallyUnderAGlobShapedRoot(t *testing.T) {
	root := hostileRoot(t)
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	line := "- doing alpha things → `.agents/skills/alpha/SKILL.md`"

	if got := renderedBlock(t, root); got != line {
		t.Fatalf("rendered block = %q, want %q", got, line)
	}

	writeFile(t, root, ".bench/BENCH-reference.md", reference(line+"\n"))
	if diags := Check(root); len(diags) != 0 {
		t.Fatalf("Check diagnostics = %v, want none", diags)
	}
	if err := Write(root); err != nil {
		t.Fatalf("Write: %v", err)
	}
	after, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := reference(line + "\n"); string(after) != want {
		t.Fatalf("Write erased the generated row:\n%q\nwant\n%q", string(after), want)
	}
}

// TestOrphanIsDiagnosedBeforeAdapterSuppression pairs the two command-named directories
// on purpose: if suppression ran before classification, the orphan would vanish behind
// its adapter and the index would report a clean tree with a missing producer in it.
func TestOrphanIsDiagnosedBeforeAdapterSuppression(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	writeFile(t, root, ".agents/commands/bench-orphan.md", "command\n")
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "bench-orphan"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, ".agents/commands/bench-valid.md", "command\n")
	writeFile(t, root, ".agents/skills/bench-valid/SKILL.md", "---\nname: bench-valid\nindex: valid adapter skill\n---\n")

	line := "- doing alpha things → `.agents/skills/alpha/SKILL.md`"
	if got := renderedBlock(t, root); got != line {
		t.Fatalf("rendered block = %q, want %q", got, line)
	}

	original := reference(line + "\n")
	writeFile(t, root, ".bench/BENCH-reference.md", original)
	wantDiag := ".agents/skills/bench-orphan/SKILL.md missing (orphan skill directory)"
	if got := strings.Join(Check(root), "\n"); got != wantDiag {
		t.Fatalf("Check =\n%s\nwant\n%s", got, wantDiag)
	}
	if err := Write(root); err == nil || err.Error() != wantDiag {
		t.Fatalf("Write error = %v, want %s", err, wantDiag)
	}
	after, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != original {
		t.Fatalf("Write changed the reference:\n%q", string(after))
	}
}

// TestEmptySkillBytesStayDistinctFromAbsence holds the two states apart: absence is an
// orphan directory, while zero bytes are a readable skill that simply declared no
// trigger, and each keeps its own diagnostic.
func TestEmptySkillBytesStayDistinctFromAbsence(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/empty/SKILL.md", "")
	if err := os.MkdirAll(filepath.Join(root, ".agents", "skills", "gone"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))

	want := "skill 'empty' missing index: frontmatter (the skills index is generated)\n" +
		".agents/skills/gone/SKILL.md missing (orphan skill directory)"
	if got := strings.Join(Check(root), "\n"); got != want {
		t.Fatalf("Check =\n%s\nwant\n%s", got, want)
	}
}

// TestRefusedCommandNamedSkillIsDiagnosedBeforeAdapterSuppression is the untrustworthy-bytes
// half of the same ordering: suppression removes a row from the rendered block, it does not
// grant amnesty from diagnosis, so a symlinked SKILL.md behind an adapter still refuses.
func TestRefusedCommandNamedSkillIsDiagnosedBeforeAdapterSuppression(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	writeFile(t, root, ".agents/commands/bench-refused.md", "command\n")
	writeFile(t, root, ".agents/skills/bench-refused/elsewhere.md", "---\nname: bench-refused\nindex: smuggled\n---\n")
	if err := os.Symlink("elsewhere.md", filepath.Join(root, ".agents", "skills", "bench-refused", "SKILL.md")); err != nil {
		t.Fatal(err)
	}

	line := "- doing alpha things → `.agents/skills/alpha/SKILL.md`"
	original := reference(line + "\n")
	writeFile(t, root, ".bench/BENCH-reference.md", original)

	wantPrefix := ".agents/skills/bench-refused/SKILL.md refused: "
	diags := Check(root)
	if len(diags) != 1 || !strings.HasPrefix(diags[0], wantPrefix) {
		t.Fatalf("Check = %v, want one diagnostic starting %q", diags, wantPrefix)
	}
	if err := Write(root); err == nil || !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Fatalf("Write error = %v, want one starting %q", err, wantPrefix)
	}
	after, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != original {
		t.Fatalf("Write changed the reference:\n%q", string(after))
	}
}
