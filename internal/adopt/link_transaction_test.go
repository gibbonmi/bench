package adopt

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransactionalLinkAdoptsUnownedAdapterThroughSymlinkParent(t *testing.T) {
	for _, tc := range []struct {
		name       string
		content    string
		target     string
		wantCode   int
		wantAbsent string
	}{
		{name: "converged", content: "same\n", wantCode: 0},
		{name: "divergent", content: "different\n", target: "../adapter-mirror", wantCode: 1, wantAbsent: "accepted.txt"},
		{name: "foreign-identical", content: "same\n", target: "../adapter-mirror", wantCode: 1, wantAbsent: "accepted.txt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			runAdoptGit(t, root, "init", "-q")
			if err := os.MkdirAll(filepath.Join(root, ".agents", "commands"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md"), []byte("same\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if tc.target != "" {
				if err := os.MkdirAll(filepath.Join(root, "adapter-mirror"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "adapter-mirror", "bench-implement-spec.md"), []byte(tc.content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.MkdirAll(filepath.Join(root, ".claude"), 0o755); err != nil {
				t.Fatal(err)
			}
			target := "../.agents/commands"
			if tc.target != "" {
				target = tc.target
			}
			if err := os.Symlink(target, filepath.Join(root, ".claude", "commands")); err != nil {
				t.Fatal(err)
			}
			kitAsset := filepath.Join(t.TempDir(), "bench-implement-spec.md")
			if err := os.WriteFile(kitAsset, []byte("same\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			plan := []planEntry{{src: kitAsset, rel: ".claude/commands/bench-implement-spec.md", kind: "adapter"}}
			if tc.wantAbsent != "" {
				plan = append([]planEntry{{rel: tc.wantAbsent, kind: "inline", content: "accepted\n"}}, plan...)
			}
			var stdout, stderr bytes.Buffer
			code, _ := transactionalLink(root, t.TempDir(), "copy", "test", plan, &stdout, &stderr)
			if code != tc.wantCode {
				t.Fatalf("transactionalLink exit = %d, want %d; stderr=%q", code, tc.wantCode, stderr.String())
			}
			if tc.wantCode == 0 {
				if got, err := os.ReadFile(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")); err != nil || string(got) != tc.content {
					t.Fatalf("resolved adapter content = %q, %v", got, err)
				}
			} else if _, err := os.Lstat(filepath.Join(root, tc.wantAbsent)); !os.IsNotExist(err) {
				t.Fatalf("accepted write was promoted: %v", err)
			}
			if tc.wantCode != 0 && !strings.Contains(stderr.String(), "has a symlink parent directory") {
				t.Fatalf("stderr = %q, want symlink-parent conflict", stderr.String())
			}
		})
	}
}
