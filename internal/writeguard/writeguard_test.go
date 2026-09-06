package writeguard

import (
	"errors"
	"strings"
	"testing"
)

// checkerFor builds a checker whose four facts are stated rather than read. primary
// names the one root that is a primary checkout, and the two sets name the paths git
// tracks and ignores under it.
func checkerFor(primary string, tracked, ignored []string) Checker {
	return Checker{
		RootAt: func(dir string) (string, error) {
			for _, root := range []string{primary, "/pool/wt"} {
				if dir == root || strings.HasPrefix(dir, root+"/") {
					return root, nil
				}
			}
			return "", errors.New("not in a git repository")
		},
		IsPrimary: func(root string) (bool, error) { return root == primary, nil },
		IsTracked: func(_, path string) bool { return contains(tracked, path) },
		IsIgnored: func(_, path string) bool { return contains(ignored, path) },
	}
}

func contains(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func TestClassify(t *testing.T) {
	checker := checkerFor("/repo", []string{"/repo/README.md", "/pool/wt/README.md"}, []string{"/repo/capture/IDEAS.md"})
	tests := []struct {
		name, path string
		blocked    bool
	}{
		{name: "tracked under the primary checkout", path: "/repo/README.md", blocked: true},
		{name: "the same relative path under a worktree", path: "/pool/wt/README.md"},
		{name: "ignored under the primary checkout", path: "/repo/capture/IDEAS.md"},
		{name: "untracked under the primary checkout", path: "/repo/notes.txt"},
		{name: "outside every repository", path: "/tmp/nowhere.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verdict := Classify(tt.path, checker)
			if verdict.Blocked != tt.blocked {
				t.Fatalf("Classify(%q).Blocked = %v, want %v", tt.path, verdict.Blocked, tt.blocked)
			}
			if verdict.Path != tt.path {
				t.Fatalf("Classify(%q).Path = %q", tt.path, verdict.Path)
			}
		})
	}
}

// TestClassifyAllowsAnUnreadableFact proves the guard denies only on a fact it read. A
// root or a checkout answer that errors leaves the write open.
func TestClassifyAllowsAnUnreadableFact(t *testing.T) {
	always := func(_, _ string) bool { return true }
	tests := []struct {
		name    string
		checker Checker
	}{
		{
			name: "the root cannot be read",
			checker: Checker{
				RootAt:    func(string) (string, error) { return "", errors.New("no repository") },
				IsPrimary: func(string) (bool, error) { return true, nil },
				IsTracked: always,
				IsIgnored: func(_, _ string) bool { return false },
			},
		},
		{
			name: "the checkout kind cannot be read",
			checker: Checker{
				RootAt:    func(string) (string, error) { return "/repo", nil },
				IsPrimary: func(string) (bool, error) { return false, errors.New("no admin directory") },
				IsTracked: always,
				IsIgnored: func(_, _ string) bool { return false },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Classify("/repo/README.md", tt.checker).Blocked {
				t.Fatal("an unreadable fact refused the write")
			}
		})
	}
}

func TestPathFromEnvelope(t *testing.T) {
	tests := []struct {
		name, envelope, want string
	}{
		{
			name:     "an absolute path",
			envelope: `{"tool_name":"Edit","tool_input":{"file_path":"/repo/README.md"},"cwd":"/elsewhere"}`,
			want:     "/repo/README.md",
		},
		{
			name:     "a relative path against the envelope cwd",
			envelope: `{"tool_name":"Write","tool_input":{"file_path":"docs/x.md"},"cwd":"/repo"}`,
			want:     "/repo/docs/x.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := PathFromEnvelope([]byte(tt.envelope))
			if err != nil {
				t.Fatalf("PathFromEnvelope: %v", err)
			}
			if path != tt.want {
				t.Fatalf("PathFromEnvelope = %q, want %q", path, tt.want)
			}
		})
	}
}

func TestPathFromEnvelopeRefusesAnUnreadablePath(t *testing.T) {
	tests := []struct{ name, envelope string }{
		{"not JSON", `{`},
		{"no tool_input", `{"tool_name":"Edit"}`},
		{"no file_path", `{"tool_name":"Edit","tool_input":{"command":"ls"}}`},
		{"a non-string file_path", `{"tool_input":{"file_path":7}}`},
		{"an empty file_path", `{"tool_input":{"file_path":""}}`},
		{"a control byte", "{\"tool_input\":{\"file_path\":\"/repo/a\\u0000b\"}}"},
		{"a relative file_path with no cwd", `{"tool_input":{"file_path":"README.md"}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if path, err := PathFromEnvelope([]byte(tt.envelope)); err == nil {
				t.Fatalf("PathFromEnvelope returned %q for an unreadable envelope", path)
			}
		})
	}
}

// TestMessageNamesThePathAndTheRepair proves the refusal carries both halves a reader
// needs: which file, and where the write belongs instead.
func TestMessageNamesThePathAndTheRepair(t *testing.T) {
	message := Verdict{Blocked: true, Path: "/repo/README.md"}.Message()

	if !strings.HasPrefix(message, "BLOCKED: ") {
		t.Fatalf("the refusal does not open with BLOCKED: %q", message)
	}
	if !strings.Contains(message, "/repo/README.md") {
		t.Fatalf("the refusal does not name the path: %q", message)
	}
	if !strings.Contains(message, "bench worktree create") {
		t.Fatalf("the refusal does not name the worktree route: %q", message)
	}
	if strings.Contains(message, "\n") {
		t.Fatalf("the refusal spans more than one line: %q", message)
	}
}
