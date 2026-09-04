package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestGoBuildArtifactModeWritesNoManifest keeps the release path inert. The build targets
// an architecture this host cannot execute, so the seal and the broker manifest, which only
// the built executable can write, are the observable proof that the mode executed nothing.
// A manifest write added to this branch would run the new binary and break every
// cross-compiled release.
func TestGoBuildArtifactModeWritesNoManifest(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	before, beforeErr := os.ReadFile(subjectManifestPath(root))
	if beforeErr != nil && !os.IsNotExist(beforeErr) {
		t.Fatal(beforeErr)
	}
	out := filepath.Join(t.TempDir(), "bench")
	target := "arm64"
	if runtime.GOARCH == target {
		target = "amd64"
	}

	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, out)
	cmd.Dir = root
	cmd.Env = append(capability.WithoutEnvironment(os.Environ(), runbinary.Env), "GOOS="+runtime.GOOS, "GOARCH="+target)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("artifact build: %v\n%s", err, output)
	}

	if _, err := os.Stat(out); err != nil {
		t.Fatalf("artifact build produced no executable: %v", err)
	}
	for _, path := range []string{out + ".seal", filepath.Join(filepath.Dir(out), brokermanifest.Name)} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("artifact build wrote %q: %v", path, err)
		}
	}
	// The subject path now publishes into the wrapper's directory rather than beside the
	// executable, so the "wrote no manifest" claim has to be graded there as well.
	if after, err := os.ReadFile(subjectManifestPath(root)); !sameWrapperManifest(after, err, before, beforeErr) {
		t.Fatalf("artifact build changed %s", subjectManifestPath(root))
	}
}

// sameWrapperManifest compares two reads of the wrapper's manifest, counting "absent both
// times" as unchanged. An artifact build must leave that file exactly as it found it,
// present or not.
func sameWrapperManifest(after []byte, afterErr error, before []byte, beforeErr error) bool {
	if os.IsNotExist(afterErr) || os.IsNotExist(beforeErr) {
		return os.IsNotExist(afterErr) && os.IsNotExist(beforeErr)
	}
	if afterErr != nil || beforeErr != nil {
		return false
	}
	return bytes.Equal(after, before)
}
