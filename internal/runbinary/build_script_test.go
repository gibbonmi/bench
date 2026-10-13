package runbinary

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestBuildEnvironmentRefusesWithoutAnAbsoluteHome(t *testing.T) {
	t.Parallel()
	built, err := buildEnvironment([]string{"PATH=/usr/bin"})
	if err == nil {
		t.Fatalf("buildEnvironment = %#v, want an error", built)
	}
	if !strings.Contains(err.Error(), "HOME") {
		t.Fatalf("error = %q, want it to name HOME", err)
	}
}

func TestBuildLeavesTheWrapperManifestUntouched(t *testing.T) {
	kit, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	for _, state := range []string{"absent", "present"} {
		t.Run(state, func(t *testing.T) {
			root := gittest.KitCopy(t, kit)
			manifest := filepath.Join(root, "bin", brokermanifest.Name)
			before := []byte("original wrapper manifest\n")
			if state == "present" {
				if err := os.WriteFile(manifest, before, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			output := filepath.Join(t.TempDir(), "bench")
			if err := Build(context.Background(), root, output); err != nil {
				t.Fatalf("private build: %v", err)
			}
			after, afterErr := os.ReadFile(manifest)
			if state == "absent" {
				if !os.IsNotExist(afterErr) {
					t.Fatalf("private build published %s", manifest)
				}
				return
			}
			if afterErr != nil || !bytes.Equal(after, before) {
				t.Fatalf("private build changed %s: %v", manifest, afterErr)
			}
		})
	}
}

// TestBuildSubjectLeavesTheManifestDirectoryToTheScript grades the one operand that parts
// the two builders. The private build names a manifest directory beside its throwaway
// executable; the subject build names none, so the script publishes the manifest beside
// the checkout's own wrapper.
func TestBuildSubjectLeavesTheManifestDirectoryToTheScript(t *testing.T) {
	source, recorded := t.TempDir(), filepath.Join(t.TempDir(), "argv")
	if err := os.MkdirAll(filepath.Join(source, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	script := "#!/usr/bin/env bash\nprintf '%s ' \"$@\" > " + recorded + "\n"
	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "bench")
	argv := func(build func(context.Context, string, string) error) string {
		if err := build(context.Background(), source, output); err != nil {
			t.Fatal(err)
		}
		recordedArgv, err := os.ReadFile(recorded)
		if err != nil {
			t.Fatal(err)
		}
		return string(recordedArgv)
	}
	if got := argv(Build); !strings.Contains(got, "--manifest-dir "+filepath.Dir(output)) {
		t.Fatalf("private build argv = %q, want the throwaway's own manifest directory", got)
	}
	if got := argv(BuildSubject); strings.Contains(got, "--manifest-dir") {
		t.Fatalf("subject build argv = %q, want no manifest directory of its own", got)
	}
}
