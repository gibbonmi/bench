// Tests for the consumer-payload allowlist states that gate check and write.
package skillsindex

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnparseableAllowlistRefusesBothWriteAndCheck(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/zeta/SKILL.md", "---\nname: zeta\nindex: doing zeta things\n---\n")
	writeFile(t, root, ".bench/consumer-payload.json", "{not json")
	refPath := filepath.Join(root, ".bench", "BENCH-reference.md")
	writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
	before, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}

	err = Write(context.Background(), root)
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
	// kit-only marker in the generated block would be missing. Grading the block
	// against it would attribute the reader's defect to the skills.
	wantCheck := ".bench/consumer-payload.json unreadable: kit-only marking unresolved"
	if got := Check(root); strings.Join(got, "\n") != wantCheck {
		t.Fatalf("check with an unparseable allowlist =\n%s\nwant\n%s", strings.Join(got, "\n"), wantCheck)
	}
}

// TestPresentInvalidAllowlistStatesRefuseCheckAndWrite walks the present-but-unusable
// partition through the canonical parser. Only absence is optional; a tree with no
// allowlist withholds nothing. Empty bytes, broken syntax, and every semantic row
// defect leave kit-only marking unresolved, and must stop both callers with the
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
			if err := Write(context.Background(), root); err == nil || err.Error() != ErrAllowlistUnreadable.Error() {
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
	if err := Write(context.Background(), root); err != nil {
		t.Fatalf("write with no allowlist = %v, want success", err)
	}
	if got := Check(root); len(got) != 0 {
		t.Fatalf("check after write with no allowlist = %s, want none", strings.Join(got, "\n"))
	}
}
