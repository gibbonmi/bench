package skillsindex

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
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
	// Check refuses on the same footing as Write: an unresolved allowlist means every
	// kit-only marker in the generated block would be missing, so grading the block
	// against it would attribute the reader's defect to the skills.
	wantCheck := ".bench/consumer-payload.json unreadable: kit-only marking unresolved"
	if got := Check(root); strings.Join(got, "\n") != wantCheck {
		t.Fatalf("check with an unparseable allowlist =\n%s\nwant\n%s", strings.Join(got, "\n"), wantCheck)
	}
}

// TestPresentInvalidAllowlistStatesRefuseCheckAndWrite walks the present-but-unusable
// partition through the canonical parser. Only absence is optional — a tree with no
// allowlist withholds nothing — while empty bytes, broken syntax, and every semantic
// row defect leave kit-only marking unresolved and must stop both callers with the
// reference bytes intact.
func TestPresentInvalidAllowlistStatesRefuseCheckAndWrite(t *testing.T) {
	for name, allowlist := range map[string]string{
		"empty":            "",
		"invalid JSON":     "{not json",
		"unknown audience": `[{"source":".agents/skills/zeta","audience":"everyone"}]`,
		"empty source":     `[{"source":"","audience":"kit-only"}]`,
		"unsafe source":    `[{"source":"../escape","audience":"kit-only"}]`,
		"duplicate source": `[{"source":"a","audience":"kit-only"},{"source":"a","audience":"consumer"}]`,
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".agents/skills/zeta/SKILL.md", "---\nname: zeta\nindex: doing zeta things\n---\n")
			writeFile(t, root, ".bench/consumer-payload.json", allowlist)
			writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
			refPath := filepath.Join(root, ".bench", "BENCH-reference.md")
			before, err := os.ReadFile(refPath)
			if err != nil {
				t.Fatal(err)
			}
			if err := Write(root); err == nil || err.Error() != ErrAllowlistUnreadable.Error() {
				t.Fatalf("write with a %s allowlist = %v, want %v", name, err, ErrAllowlistUnreadable)
			}
			after, err := os.ReadFile(refPath)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Fatalf("refused write changed bytes:\n%q\nwant\n%q", after, before)
			}
			want := ".bench/consumer-payload.json unreadable: kit-only marking unresolved"
			if got := Check(root); strings.Join(got, "\n") != want {
				t.Fatalf("check with a %s allowlist =\n%s\nwant\n%s", name, strings.Join(got, "\n"), want)
			}
		})
	}
}

// TestAbsentAllowlistStaysOptional is the control the partition above needs: absence
// is the one present-invalid state's opposite, and it must still index cleanly.
func TestAbsentAllowlistStaysOptional(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/zeta/SKILL.md", "---\nname: zeta\nindex: doing zeta things\n---\n")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
	if err := Write(root); err != nil {
		t.Fatalf("write with no allowlist = %v, want success", err)
	}
	if got := Check(root); len(got) != 0 {
		t.Fatalf("check after write with no allowlist = %s, want none", strings.Join(got, "\n"))
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

// A skill file is attacker-shaped input, so only a fence that opens at byte zero and
// closes may authorize an index field: a late fence lets prose smuggle one in, and an
// unclosed opener lets the whole body act as frontmatter.
func TestFrontmatterFieldRequiresCompleteLeadingFence(t *testing.T) {
	for _, row := range []struct {
		name    string
		content string
		want    string
	}{
		{"complete leading fence", "---\nname: probe\nindex: taken\n---\n", "taken"},
		{"unclosed opener without a trailing newline", "---\nname: probe\nindex: taken", ""},
		{"closing fence as a last line without a trailing newline", "---\nname: probe\nindex: taken\n---", "taken"},
		{"duplicate key keeps the first value", "---\nindex: first value\nindex: second value\n---\n", "first value"},
		{"duplicate key keeps an empty first value", "---\nindex:\nindex: second value\n---\n", ""},
		{"text before the opener", "prose\n---\nindex: smuggled\n---\n", ""},
		{"blank line before the opener", "\n---\nindex: smuggled\n---\n", ""},
		{"unclosed opener", "---\nname: probe\nindex: smuggled\n", ""},
		{"no fence at all", "index: smuggled\n", ""},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "SKILL.md", row.content)
			if got := FrontmatterField(filepath.Join(root, "SKILL.md"), "index"); got != row.want {
				t.Fatalf("index = %q, want %q", got, row.want)
			}
		})
	}
}

// referenceSnapshot records what the reference is without ever reading a path the
// classifier would refuse: a FIFO opened here would block the test itself, which is
// the failure the production reader exists to avoid.
func referenceSnapshot(t *testing.T, path string) string {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		return "absent"
	}
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		target, err := os.Readlink(path)
		if err != nil {
			t.Fatal(err)
		}
		return "symlink " + target
	case info.Mode().IsRegular():
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return "regular " + string(data)
	}
	return "other " + info.Mode().Type().String()
}

// await runs one reference-reading call under a deadline. A reader that reopens the
// path directly blocks in open(2) on the FIFO row forever, and a hung package test
// reports as a suite-wide timeout rather than as this row, so the bound is here.
func await[T any](t *testing.T, what string, call func() T) T {
	t.Helper()
	done := make(chan T, 1)
	go func() { done <- call() }()
	select {
	case got := <-done:
		return got
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatalf("%s blocked on the reference, so it opened the path before classifying it", what)
	}
	panic("unreachable")
}

// TestReferenceProducerStatesStayDistinctAndBlockWrite is HI4: every hostile or
// degenerate reference disposition keeps its own diagnostic — empty is not missing,
// and an untrustworthy object is neither — and none of them lets Write touch the
// bytes it could not read.
func TestReferenceProducerStatesStayDistinctAndBlockWrite(t *testing.T) {
	for _, tc := range []struct {
		name    string
		plant   func(*testing.T, string)
		want    string
		refusal bool
	}{
		{name: "missing", want: referenceRel + " missing (skills index unverifiable)"},
		{
			name:  "empty",
			plant: func(t *testing.T, path string) { writeAt(t, path, "") },
			want:  referenceRel + " empty (skills index unverifiable)",
		},
		{
			name: "oversized",
			plant: func(t *testing.T, path string) {
				writeAt(t, path, strings.Repeat("a", int(bounds.ControlRecordLimit)+1))
			},
			refusal: true,
		},
		{
			name:    "invalid UTF-8",
			plant:   func(t *testing.T, path string) { writeAt(t, path, reference("")+"\xff\xfe") },
			refusal: true,
		},
		{
			name: "directory",
			plant: func(t *testing.T, path string) {
				if err := os.MkdirAll(path, 0o755); err != nil {
					t.Fatal(err)
				}
			},
			refusal: true,
		},
		{
			name: "fifo",
			plant: func(t *testing.T, path string) {
				if err := syscall.Mkfifo(path, 0o644); err != nil {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
				}
			},
			refusal: true,
		},
		{
			name: "live symlink",
			plant: func(t *testing.T, path string) {
				target := filepath.Join(filepath.Dir(path), "elsewhere.md")
				writeAt(t, target, reference(""))
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
			refusal: true,
		},
		{
			name: "dangling symlink",
			plant: func(t *testing.T, path string) {
				if err := os.Symlink(filepath.Join(filepath.Dir(path), "gone.md"), path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
			refusal: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
			path := filepath.Join(root, filepath.FromSlash(referenceRel))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if tc.plant != nil {
				tc.plant(t, path)
			}
			before := referenceSnapshot(t, path)

			diags := await(t, "Check", func() []string { return Check(root) })
			err := await(t, "Write", func() error { return Write(root) })
			if err == nil {
				t.Fatal("Write accepted a reference it could not read")
			}
			if got := referenceSnapshot(t, path); got != before {
				t.Fatalf("Write changed the reference it refused: %.60q, want %.60q", got, before)
			}
			if tc.refusal {
				prefix := ReferenceRefusalPrefix()
				if !hasPrefixed(diags, prefix) {
					t.Fatalf("Check reported %q, want a diagnostic prefixed %q", diags, prefix)
				}
				if !strings.HasPrefix(err.Error(), prefix) || strings.TrimSpace(strings.TrimPrefix(err.Error(), prefix)) == "" {
					t.Fatalf("Write reported %q, want an attributed refusal prefixed %q", err, prefix)
				}
				return
			}
			if !hasExact(diags, tc.want) {
				t.Fatalf("Check reported %q, want %q", diags, tc.want)
			}
			if err.Error() != tc.want {
				t.Fatalf("Write reported %q, want %q", err, tc.want)
			}
		})
	}
}

// TestControlRunesNeverReachTheRenderedLine is the sink contract: the index line is
// line-structured markdown, so any control rune in a rendered field could split or
// forge an entry and is refused before rendering, while ordinary graphic Unicode —
// accented Latin, CJK, emoji, arrows — still renders as exactly one line. The
// permitted half is asserted alongside the refused one because an ASCII-only fix
// passes the refusals and quietly loses the rest of the control partition.
func TestControlRunesNeverReachTheRenderedLine(t *testing.T) {
	for _, row := range []struct {
		name    string
		key     string
		value   string
		refused bool
	}{
		{"tab in the trigger", "index", "doing\tthings", true},
		{"carriage return forging a second line", "index", "safe → `x`\r- forged → `y`", true},
		{"escape in the trigger", "index", "doing\x1b[31m things", true},
		{"bell in the trigger", "index", "doing\a things", true},
		{"nul in the trigger", "index", "doing\x00 things", true},
		{"delete in the trigger", "index", "doing\x7f things", true},
		{"c1 next-line in the trigger", "index", "doing\u0085things", true},
		{"tab in the note", "index-note", "a\tnote", true},
		{"carriage return in the note", "index-note", "a note\r- forged", true},
		{"escape in the note", "index-note", "a\x1bnote", true},
		{"graphic unicode in the trigger", "index", "café 文書 🚀 →", false},
		{"graphic unicode in the note", "index-note", "註記 ✨ →", false},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := t.TempDir()
			trigger, note := "doing probe things", "a probe note"
			if row.key == "index" {
				trigger = row.value
			} else {
				note = row.value
			}
			writeFile(t, root, ".agents/skills/probe/SKILL.md", "---\nindex: "+trigger+"\nindex-note: "+note+"\n---\n")
			original := reference("")
			writeFile(t, root, ".bench/BENCH-reference.md", original)

			block := renderedBlock(t, root)
			if !row.refused {
				want := "- " + trigger + " → `.agents/skills/probe/SKILL.md` + " + note
				if block != want {
					t.Fatalf("rendered block =\n%q\nwant\n%q", block, want)
				}
				return
			}
			if block != "" {
				t.Fatalf("refused field rendered a line: %q", block)
			}
			wantPrefix := ".agents/skills/probe/SKILL.md refused: "
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
		})
	}
}

func writeAt(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func hasExact(diags []string, want string) bool {
	for _, diag := range diags {
		if diag == want {
			return true
		}
	}
	return false
}

func hasPrefixed(diags []string, prefix string) bool {
	for _, diag := range diags {
		if strings.HasPrefix(diag, prefix) && strings.TrimSpace(strings.TrimPrefix(diag, prefix)) != "" {
			return true
		}
	}
	return false
}

// A newline cannot survive the fence scan's line split, so the sink table above cannot
// construct one from a file; the predicate is asserted directly so a later parser that
// does yield multi-line values still finds the sink closed.
func TestControlRefusalCoversTheNewlineTheParserCannotYield(t *testing.T) {
	if controlRefusal("index", "safe\n- forged") == "" {
		t.Fatal("newline accepted in an index field")
	}
	if got := controlRefusal("index", "café 文書 🚀 →"); got != "" {
		t.Fatalf("graphic unicode refused: %s", got)
	}
}
