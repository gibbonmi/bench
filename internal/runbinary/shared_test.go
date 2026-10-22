package runbinary

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// countingFactory answers a factory whose builder counts its calls and whose verifier
// accepts every executable, so only the selection rule decides the result.
func countingFactory(t *testing.T, builds *int) Factory {
	t.Helper()
	return Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			*builds++
			return os.WriteFile(output, []byte("private"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	}
}

// inheritedExecutable answers a regular executable that the counting verifier accepts.
func inheritedExecutable(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(path, []byte("inherited"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSharedInheritedSelectionWinsWithZeroBuilds(t *testing.T) {
	inherited := inheritedExecutable(t)
	builds := 0
	shared := &Shared{Factory: countingFactory(t, &builds), Raw: inherited, Present: true}
	t.Cleanup(func() { _ = shared.Close() })
	for _, caller := range []string{"first caller", "second caller"} {
		path, err := shared.Path(context.Background(), t.TempDir())
		if err != nil || path != inherited {
			t.Fatalf("%s selected %q, %v; want the inherited %q unchanged", caller, path, err, inherited)
		}
	}
	if builds != 0 {
		t.Fatalf("an inherited selection ran the builder %d times, want 0", builds)
	}
}

func TestSharedInvalidInheritedSelectionRefusesBeforeAnyBuild(t *testing.T) {
	for _, tc := range []struct {
		name, raw, source string
	}{
		{"missing", filepath.Join(t.TempDir(), "absent"), t.TempDir()},
		{"relative", "dist/bench", t.TempDir()},
		{"empty", "", t.TempDir()},
		{"no source root", inheritedExecutable(t), ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builds := 0
			shared := &Shared{Factory: countingFactory(t, &builds), Raw: tc.raw, Present: true}
			t.Cleanup(func() { _ = shared.Close() })
			for range 2 {
				if path, err := shared.Path(context.Background(), tc.source); err == nil {
					t.Fatalf("the present %s value %q selected %q, want a refusal", tc.name, tc.raw, path)
				}
			}
			if builds != 0 {
				t.Fatalf("the present %s value %q ran the builder %d times, want 0", tc.name, tc.raw, builds)
			}
		})
	}
}
