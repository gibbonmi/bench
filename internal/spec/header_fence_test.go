package spec

import "testing"

// TestMetadataSkipsAFencedStatus pins the one fence classification of the header readers: a
// tilde fence and an indented backtick fence each hide a spec path from LiveSpecSlugs and a
// column-zero `Status:` line from Facts. A first-line frontmatter block hides them too.
func TestMetadataSkipsAFencedStatus(t *testing.T) {
	for name, content := range map[string]string{
		"tilde fence":       "Status: staged\n\n~~~\nspecs/hidden/spec.md\nStatus: implemented\n~~~\n",
		"indented fence":    "Status: staged\n\n  ```\nspecs/hidden/spec.md\nStatus: implemented\n  ```\n",
		"frontmatter block": "---\nspecs/hidden/spec.md\nStatus: implemented\n---\nStatus: staged\n",
	} {
		t.Run(name, func(t *testing.T) {
			if got := LiveSpecSlugs([]byte(content)); len(got) != 0 {
				t.Errorf("LiveSpecSlugs = %q, want none", got)
			}
			root := t.TempDir()
			writeFolderSpec(t, root, "s", content)
			facts, err := Facts(root)
			if err != nil || len(facts) != 1 || facts[0].Status != "staged" {
				t.Errorf("Facts = %#v, %v, want one fact with the status staged", facts, err)
			}
		})
	}
}

// TestImplementedSkipsAFencedStatus pins the flip over the reader's lines: a fenced copy of the
// status line is no second match, and a CRLF spec flips with each carriage return kept.
func TestImplementedSkipsAFencedStatus(t *testing.T) {
	for name, tc := range map[string]struct{ in, want string }{
		"fenced copy": {"Status: staged\n````\n```\nStatus: staged\n````\n", "Status: implemented\n````\n```\nStatus: staged\n````\n"},
		"crlf":        {"# s\r\nStatus: staged\r\nbody\r\n", "# s\r\nStatus: implemented\r\nbody\r\n"},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := Implemented([]byte(tc.in))
			if err != nil || string(got) != tc.want {
				t.Errorf("Implemented = %q, %v, want %q", got, err, tc.want)
			}
		})
	}
}
