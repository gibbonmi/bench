// Tests for the frontmatter read of a CRLF skill file.
package skillsindex

import (
	"path/filepath"
	"testing"
)

// TestFrontmatterFieldReadsCRLF pins both sides of the frontmatter marker predicate on
// a CRLF skill file, and the unclosed opener of such a file.
func TestFrontmatterFieldReadsCRLF(t *testing.T) {
	for _, row := range []struct {
		name, content, want string
	}{
		{"CRLF block", "---\r\nindex: x\r\n---\r\n", "x"},
		{"CRLF opener with a trailing space", "--- \r\nindex: x\r\n---\r\n", ""},
		{"CRLF opener with no closer", "---\r\nindex: x\r\n", ""},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "SKILL.md", row.content)
			if got := FrontmatterField(filepath.Join(root, "SKILL.md"), "index"); got != row.want {
				t.Errorf("%s: index = %q, want %q", row.name, got, row.want)
			}
		})
	}
}
