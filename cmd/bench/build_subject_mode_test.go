package main

import (
	"bytes"
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
	preserveWrapperManifest(t, root)

	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, out)
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("subject build: %v\n%s", err, output)
	}

	fields, err := brokermanifest.Read(subjectManifestPath(root))
	if err != nil {
		t.Fatal(err)
	}
	want := stampedPackageVersion(t, root)
	if fields["version"] != want {
		t.Fatalf("published manifest version = %q, want the stamped version %q", fields["version"], want)
	}
	// The wrapper's own directory is the only place bin/bench.sh and the doctor row look,
	// so a manifest beside the executable binds nothing either reader can find. The digest
	// comparison keeps the found manifest bound to the executable this build published,
	// rather than to whatever a previous run left in the wrapper's directory.
	published, err := brokermanifest.Digest(out)
	if err != nil {
		t.Fatal(err)
	}
	if fields["sha256"] != published {
		t.Fatalf("manifest at %s binds digest %q, want the published executable digest %q", subjectManifestPath(root), fields["sha256"], published)
	}
	// The digest alone cannot tell this build's manifest from a leftover one, because the
	// build is reproducible and any prior manifest carries the same digest. The bound path
	// is what changes: it names the executable this run published.
	bound := out
	if resolved, err := filepath.EvalSymlinks(out); err == nil {
		bound = resolved
	}
	if fields["path"] != bound {
		t.Fatalf("manifest at %s binds path %q, want this build's executable %q", subjectManifestPath(root), fields["path"], bound)
	}
}

// TestReleasePreflightBuildDoesNotRebindThePromotionBroker grades the auxiliary build
// scripts/release-preflight.sh runs. That script builds dist/bench-preflight, which is a
// build output a clean deletes, and the wrapper's manifest names whatever the last
// subject-mode build published. A manifest bound to the preflight artifact therefore
// survives the artifact, and bin/bench.sh execs the missing path as the promotion broker:
// the landing refuses at exit 127 with nothing wrong in the tree. Only the build that
// publishes the checkout's own dist/bench may rebind that file.
//
// The row runs the real script, because the defect is which operands the script hands the
// builder. A test that restated those operands would grade its own copy of them.
func TestReleasePreflightBuildDoesNotRebindThePromotionBroker(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	preserveWrapperManifest(t, root)
	before, beforeErr := os.ReadFile(subjectManifestPath(root))
	if beforeErr != nil && !os.IsNotExist(beforeErr) {
		t.Fatal(beforeErr)
	}

	// The script skips its build when the artifact is already present, so the row removes
	// it to reach the build at all. It is a build output, which a clean deletes anyway.
	artifact := filepath.Join(root, "dist", "bench-preflight")
	if err := os.Remove(artifact); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	// The script execs the built binary once the build finishes. This row grades the build,
	// not the verb, so the exit status of that exec says nothing and is ignored.
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "release-preflight.sh"), "--help")
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	_, _ = cmd.CombinedOutput()
	if _, err := os.Stat(artifact); err != nil {
		t.Fatalf("release-preflight build produced no auxiliary executable: %v", err)
	}

	after, afterErr := os.ReadFile(subjectManifestPath(root))
	if os.IsNotExist(beforeErr) {
		if !os.IsNotExist(afterErr) {
			t.Fatalf("release-preflight build published %s", subjectManifestPath(root))
		}
		return
	}
	if afterErr != nil || !bytes.Equal(after, before) {
		t.Fatalf("release-preflight build rebound %s: %v", subjectManifestPath(root), afterErr)
	}
}

// subjectManifestPath names where a subject-mode build publishes the broker manifest: the
// wrapper's own directory, which is the one place the land route and the doctor row read.
func subjectManifestPath(root string) string {
	return filepath.Join(root, "bin", brokermanifest.Name)
}

// preserveWrapperManifest restores whatever the checkout's wrapper directory held before a
// build test republished over it. The row has to run against the real module root, because
// only that root carries the build inputs, so it writes where a developer's own dev-install
// manifest lives.
func preserveWrapperManifest(t *testing.T, root string) {
	t.Helper()
	path := subjectManifestPath(root)
	before, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	existed := err == nil
	t.Cleanup(func() {
		if !existed {
			_ = os.Remove(path)
			return
		}
		if err := os.WriteFile(path, before, 0o644); err != nil {
			t.Error(err)
		}
	})
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
