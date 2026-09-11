package gate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

func TestCompletionFile(t *testing.T) {
	for _, kind := range []string{"regular", "executable", "missing", "symlink", "oversized"} {
		t.Run(kind, func(t *testing.T) {
			f := recordtest.New(t, 1)
			path := "completion.txt"
			switch kind {
			case "regular", "executable":
				f.Write(path, "source\n")
				if kind == "executable" {
					if err := os.Chmod(filepath.Join(f.Root, path), 0755); err != nil {
						t.Fatal(err)
					}
				}
			case "symlink":
				if err := os.Symlink("source.txt", filepath.Join(f.Root, path)); err != nil {
					t.Fatal(err)
				}
			case "oversized":
				f.Write(path, string(make([]byte, int(bounds.ControlRecordLimit)+1)))
			}
			if kind != "missing" {
				f.Commit("completion file case")
			}
			generation, err := captureProspectiveTree(f.Root, f.Tree())
			if err != nil {
				t.Fatal(err)
			}
			result, err := newGateEvaluation(f.Root).completionFile(generation, path)
			if kind != "regular" && kind != "executable" {
				if err == nil {
					t.Fatalf("%s completion file admitted: %d bytes", kind, len(result.data))
				}
				return
			}
			mode := "100644"
			if kind == "executable" {
				mode = "100755"
			}
			if err != nil || string(result.data) != "source\n" || result.mode != mode {
				t.Fatalf("completion file = %+v, %v", result, err)
			}
		})
	}
}
