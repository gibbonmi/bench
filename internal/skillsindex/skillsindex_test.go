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
	// alpha collides two diagnostics on one skill: it declares no trigger AND still
	// carries a committed entry. A per-skill map would drop one of them, so the pair
	// is asserted here rather than left to the canaries — none of which has both a
	// field-less skill and a non-empty committed block.
	writeFile(t, root, ".bench/BENCH-reference.md", reference(
		"- doing alpha things → `.agents/skills/alpha/SKILL.md`\n"+
			"- stale beta wording → `.agents/skills/beta/SKILL.md`\n"+
			"- doing gamma things → `.agents/skills/gamma/SKILL.md`\n"))

	want := []string{
		"skill 'alpha' missing index: frontmatter (the skills index is generated)",
		"skills index entry 'alpha' has no indexed .agents/skills/alpha on disk (regenerate: bench skills-index --write)",
		"skills index entry for 'beta' drifted from its frontmatter (regenerate: bench skills-index --write)",
		"skills index entry 'gamma' has no indexed .agents/skills/gamma on disk (regenerate: bench skills-index --write)",
	}
	if got := Check(root); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("attributed diagnostics =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
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
	// Check keeps today's posture: it marks nothing and adds no allowlist diagnostic,
	// so the gate gains no red the index itself did not earn.
	wantCheck := "skills index missing entry for skill 'zeta' (regenerate: bench skills-index --write)"
	if got := Check(root); strings.Join(got, "\n") != wantCheck {
		t.Fatalf("check with an unparseable allowlist =\n%s\nwant\n%s", strings.Join(got, "\n"), wantCheck)
	}
}
