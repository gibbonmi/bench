package artifact

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactBuilderRejectsDirtyAndUntrackedSourceState(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	for _, test := range []struct {
		name   string
		expect string
		mutate func(*testing.T, string)
	}{
		{name: "tracked", expect: "bench artifacts: source state must be clean and tracked at HEAD", mutate: func(t *testing.T, source string) {
			file := filepath.Join(source, "package.json")
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file, append(data, '\n'), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "untracked", expect: "bench artifacts: source state must be clean and tracked at HEAD", mutate: func(t *testing.T, source string) {
			if err := os.WriteFile(filepath.Join(source, "untracked-release-input"), []byte("uncommitted\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "status error", expect: "bench artifacts: could not verify source state at HEAD", mutate: func(t *testing.T, source string) {
			if err := os.WriteFile(filepath.Join(source, ".git", "index"), []byte("not a git index\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "dirty initialized submodule", expect: "bench artifacts: source state must be clean and tracked at HEAD", mutate: dirtyArtifactSubmodule},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := committedHostileArtifactSource(t, root)
			test.mutate(t, source)
			prepared := t.TempDir()
			if err := os.WriteFile(filepath.Join(prepared, "prepared-artifact"), []byte("prepared\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			ready := filepath.Join(t.TempDir(), "promotion-ready")
			command := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, filepath.Join(t.TempDir(), "artifacts"))
			command.Env = promotionTestEnv(prepared, ready)
			output, err := runArtifactBuildThroughPromotionSeam(t, command, ready)
			if err == nil {
				t.Fatalf("artifact builder accepted %s source state:\n%s", test.name, output)
			}
			if !strings.Contains(string(output), test.expect) {
				t.Fatalf("artifact builder did not attribute %s source state:\n%s", test.name, output)
			}
		})
	}
}

func TestArtifactBuilderRefusesMissingBinaryPinManifest(t *testing.T) {
	root := contract.SubjectRoot(t)
	source := committedHostileArtifactSource(t, root)
	command := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, filepath.Join(t.TempDir(), "artifacts"))
	command.Env = append(os.Environ(), "BENCH_TEST_SKIP_PIN_MANIFEST=1", "BENCH_REPRO_BUILD=1")
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("artifact builder packed wrapper without binary pins:\n%s", output)
	}
	if !strings.Contains(string(output), "binary pin manifest is missing or empty") {
		t.Fatalf("missing pin refusal was not attributed:\n%s", output)
	}
}

func dirtyArtifactSubmodule(t *testing.T, source string) {
	t.Helper()
	origin := t.TempDir()
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "submodule@example.invalid"}, {"config", "user.name", "submodule"}} {
		if output, err := exec.Command("git", append([]string{"-C", origin}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("prepare submodule %v: %v\n%s", args, err, output)
		}
	}
	if err := os.WriteFile(filepath.Join(origin, "tracked"), []byte("clean\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "tracked"}, {"commit", "-qm", "submodule source"}} {
		if output, err := exec.Command("git", append([]string{"-C", origin}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("commit submodule %v: %v\n%s", args, err, output)
		}
	}
	add := exec.Command("git", "-c", "protocol.file.allow=always", "-C", source, "submodule", "add", "-q", origin, "nested-release-input")
	if output, err := add.CombinedOutput(); err != nil {
		t.Fatalf("add submodule: %v\n%s", err, output)
	}
	commit := exec.Command("git", "-C", source, "commit", "-qam", "add initialized submodule")
	if output, err := commit.CombinedOutput(); err != nil {
		t.Fatalf("commit submodule: %v\n%s", err, output)
	}
	if err := os.WriteFile(filepath.Join(source, "nested-release-input", "tracked"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
