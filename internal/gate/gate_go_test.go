package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// TestGofmtStepGradesOnlyTheVisibleTree is the bite proof for the gofmt step's scope. An
// unformatted file under an ignored path is local state no commit carries, so the step
// stays green over it. The same file at a visible untracked path is about to be committed,
// so the step reds and names it.
func TestGofmtStepGradesOnlyTheVisibleTree(t *testing.T) {
	root := gittest.RepoIgnoring(t, ".logs/")
	unformatted := []byte("package overlay\nfunc  f( ) {}\n")
	write := func(rel string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, unformatted, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write(".logs/copy/overlay.go")
	var stdout, stderr bytes.Buffer
	if code := GateGoCommand([]string{"gofmt", root}, &stdout, &stderr); code != 0 {
		t.Fatalf("ignored unformatted file: exit %d, stderr %q, want green", code, stderr.String())
	}

	write("overlay/overlay.go")
	stderr.Reset()
	code := GateGoCommand([]string{"gofmt", root}, &stdout, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "gofmt: unformatted Go files: overlay/overlay.go") {
		t.Fatalf("visible unformatted file: exit %d, stderr %q, want red naming overlay/overlay.go", code, stderr.String())
	}
	if strings.Contains(stderr.String(), ".logs/") {
		t.Fatalf("visible unformatted file: stderr %q also names the ignored copy", stderr.String())
	}
}
