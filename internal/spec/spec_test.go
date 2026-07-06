package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeSpec writes content to <dir>/specs/<slug>.md and returns the path.
func writeSpec(t *testing.T, dir, slug, content string) string {
	t.Helper()
	specsDir := filepath.Join(dir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(specsDir, slug+".md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFlipRewritesOnlyTheStatusLine(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "space separator, trailing newline",
			in:   "# spec\n\nStatus: staged\n\n## body\n",
			want: "# spec\n\nStatus: implemented\n\n## body\n",
		},
		{
			name: "tab separator preserved",
			in:   "Status:\tstaged\nbody\n",
			want: "Status:\timplemented\nbody\n",
		},
		{
			name: "no trailing newline preserved",
			in:   "# spec\nStatus: staged",
			want: "# spec\nStatus: implemented",
		},
		{
			name: "trailing whitespace on status line preserved",
			in:   "Status:  staged  \n",
			want: "Status:  implemented  \n",
		},
		{
			name: "other staged word untouched",
			in:   "Status: staged\nThe work is staged elsewhere.\n",
			want: "Status: implemented\nThe work is staged elsewhere.\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := writeSpec(t, dir, "s", tc.in)
			resolved, err := Flip(dir, "s")
			if err != nil {
				t.Fatalf("Flip: %v", err)
			}
			if resolved != path {
				t.Errorf("resolved = %q, want %q", resolved, path)
			}
			got, _ := os.ReadFile(path)
			if string(got) != tc.want {
				t.Errorf("content = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFlipErrors(t *testing.T) {
	cases := []struct {
		name    string
		content string // "" means don't create the file (not-found)
		wantSub string
	}{
		{"not found", "", "not found"},
		{"no staged line", "# spec\nStatus: draft\n", "no `Status: staged`"},
		{"already implemented", "# spec\nStatus: implemented\n", "no `Status: staged`"},
		{"more than one staged line", "Status: staged\nStatus: staged\n", "expected exactly one"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.content != "" {
				writeSpec(t, dir, "s", tc.content)
			}
			_, err := Flip(dir, "s")
			if err == nil {
				t.Fatal("Flip: expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantSub)
			}
		})
	}
}

func TestFlipReRunIsNonDestructive(t *testing.T) {
	dir := t.TempDir()
	path := writeSpec(t, dir, "s", "Status: staged\n")
	if _, err := Flip(dir, "s"); err != nil {
		t.Fatalf("first flip: %v", err)
	}
	after, _ := os.ReadFile(path)
	if _, err := Flip(dir, "s"); err == nil {
		t.Fatal("second flip: expected error on already-implemented spec")
	}
	again, _ := os.ReadFile(path)
	if string(again) != string(after) {
		t.Errorf("second flip mutated the file: %q -> %q", after, again)
	}
}

func TestResolveConvention(t *testing.T) {
	dir := t.TempDir()
	slugPath := writeSpec(t, dir, "mine", "Status: staged\n")

	// bare slug resolves to <base>/specs/<slug>.md
	_, resolved, _, ok, err := Resolve(dir, "mine")
	if err != nil || !ok {
		t.Fatalf("bare slug: ok=%v err=%v", ok, err)
	}
	if resolved != slugPath {
		t.Errorf("bare slug resolved = %q, want %q", resolved, slugPath)
	}

	// a path argument resolves as-given (cwd-relative), not via the fallback
	abs := slugPath
	_, resolved, _, ok, err = Resolve(dir, abs)
	if err != nil || !ok || resolved != abs {
		t.Errorf("path arg: resolved=%q ok=%v err=%v, want %q", resolved, ok, err, abs)
	}

	// a slug containing '/' gets no fallback
	_, _, tried, ok, _ := Resolve(dir, "sub/mine")
	if ok {
		t.Error("slug with '/': expected no resolution")
	}
	if len(tried) != 1 || tried[0] != "sub/mine" {
		t.Errorf("slug with '/': tried = %v, want only the as-given form", tried)
	}
}

func TestResolveBaseAnchorsFallbackFromAnyCwd(t *testing.T) {
	root := t.TempDir()
	writeSpec(t, root, "anchored", "Status: staged\n")
	// Resolve from a different cwd: the fallback must still hit <root>/specs.
	_, resolved, _, ok, err := Resolve(root, "anchored")
	if err != nil || !ok {
		t.Fatalf("anchored resolve: ok=%v err=%v", ok, err)
	}
	if resolved != filepath.Join(root, "specs", "anchored.md") {
		t.Errorf("resolved = %q, want root-anchored path", resolved)
	}
}
