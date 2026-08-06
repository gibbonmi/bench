package posture

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactModeCommandTraceExcludesPublication(t *testing.T) {
	root := contract.SubjectRoot(t)
	trace, output := runArtifactPublicationTrace(t, root)
	if diags := artifactPublicationTraceDiagnostics(trace, output); len(diags) != 0 {
		t.Fatalf("artifact publication trace: %v", diags)
	}
}

func TestArtifactModeCommandTraceBitesOnPublishThenDelete(t *testing.T) {
	root := writeArtifactBuilderFixture(t)
	script := filepath.Join(root, "scripts", "go-build.sh")
	body, err := os.ReadFile(script)
	if err != nil {
		t.Fatal(err)
	}
	needle := "  go build \"${go_build_flags[@]}\" -o \"$staged\" ./cmd/bench\n"
	mutation := needle + "  \"$staged\" freshness-publish \"$modroot\" \"$out\"\n  rm -f -- \"$out.seal\"\n  exit 0\n"
	if !strings.Contains(string(body), needle) {
		t.Fatal("artifact build seam not found")
	}
	if err := os.WriteFile(script, []byte(strings.Replace(string(body), needle, mutation, 1)), 0o755); err != nil {
		t.Fatal(err)
	}

	trace, output := runArtifactPublicationTrace(t, root)
	diags := strings.Join(artifactPublicationTraceDiagnostics(trace, output), "\n")
	if !strings.Contains(diags, "executed staged Bench publication") {
		t.Fatalf("publish-then-delete diagnostics = %q, want staged-publication execution", diags)
	}
}

func runArtifactPublicationTrace(t *testing.T, root string) ([]string, string) {
	t.Helper()
	dir := t.TempDir()
	tracePath := filepath.Join(dir, "trace")
	bin := filepath.Join(dir, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFakeGoBuilder(t, bin)
	output := filepath.Join(dir, "artifact")
	command := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, output)
	command.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"), "BENCH_PUBLICATION_TRACE="+tracePath, "GOOS=js", "GOARCH=wasm")
	if out, err := command.CombinedOutput(); err != nil {
		t.Fatalf("artifact builder: %v\n%s", err, out)
	}
	data, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n"), output
}

func artifactPublicationTraceDiagnostics(trace []string, output string) []string {
	var diags []string
	builds := 0
	for _, line := range trace {
		switch {
		case strings.HasPrefix(line, "go:build "):
			builds++
		case strings.HasPrefix(line, "go:run "):
			diags = append(diags, "artifact mode invoked go run")
		case strings.HasPrefix(line, "bench:"):
			diags = append(diags, "artifact mode executed staged Bench publication: "+line)
		}
	}
	if builds != 1 {
		diags = append(diags, fmt.Sprintf("artifact mode go build count = %d, want 1", builds))
	}
	if info, err := os.Stat(output); err != nil || !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		diags = append(diags, fmt.Sprintf("artifact output is not a regular executable: %v, %v", info, err))
	}
	if _, err := os.Lstat(output + ".seal"); !os.IsNotExist(err) {
		diags = append(diags, fmt.Sprintf("artifact output retained a freshness seal: %v", err))
	}
	return diags
}

// writeArtifactBuilderFixture copies the builder and the inputs it reads into a scratch
// root, so a contract can mutate the builder script without touching the subject tree.
func writeArtifactBuilderFixture(t *testing.T) string {
	t.Helper()
	source := contract.SubjectRoot(t)
	root := t.TempDir()
	for _, rel := range []string{"scripts/go-build.sh", "scripts/go-build.inputs", "package.json", "internal/releaseevidence/requirements.json"} {
		data, err := os.ReadFile(filepath.Join(source, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") {
			mode = 0o755
		}
		if err := os.WriteFile(path, data, mode); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
