package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestGoBuildSubjectModePublishesTheStampedVersion grades the wiring between the builder's
// version stamp and the manifest the land route authenticates. The publication library
// binds whatever version its caller hands it, so only a real subject-mode build shows
// whether the script hands it the stamped one. A manifest that carries "dev" fails the
// land route's version comparison, and no library test reaches the binding.
func TestGoBuildSubjectModePublishesTheStampedVersion(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "bench")

	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, out)
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("subject build: %v\n%s", err, output)
	}

	fields, err := brokermanifest.Read(filepath.Join(filepath.Dir(out), brokermanifest.Name))
	if err != nil {
		t.Fatal(err)
	}
	want := stampedPackageVersion(t, root)
	if fields["version"] != want {
		t.Fatalf("published manifest version = %q, want the stamped version %q", fields["version"], want)
	}
}

// stampedPackageVersion reads the version through the same two operands the builder walks:
// the package-version input the build owner manifest names, and that file's version field.
// Reading the same source is what keeps the expectation from becoming a second copy of the
// released version. It refuses "dev", because an unstamped build is the failure this row
// exists to catch, not an expectation to mirror.
func stampedPackageVersion(t *testing.T, root string) string {
	t.Helper()
	inputs, err := os.ReadFile(filepath.Join(root, "scripts", "go-build.inputs"))
	if err != nil {
		t.Fatal(err)
	}
	source := ""
	for _, line := range strings.Split(string(inputs), "\n") {
		if key, value, ok := strings.Cut(line, "="); ok && key == "package_version" {
			source = value
		}
	}
	if source == "" {
		t.Fatal("scripts/go-build.inputs names no package_version input")
	}
	body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(source)))
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(body, &pkg); err != nil {
		t.Fatal(err)
	}
	if pkg.Version == "" || pkg.Version == "dev" {
		t.Fatalf("%s carries no stamped version: %q", source, pkg.Version)
	}
	return pkg.Version
}
