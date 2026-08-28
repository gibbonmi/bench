package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// The builder owns CGO_ENABLED, so the caller's environment cannot change the bytes.
// Go derives that variable from the host: a cross-compiling runner gets 0, and a runner
// whose host equals its target gets 1. While the callers owned the setting, the release
// workflow built each Darwin package on a Linux runner and rebuilt it on a macOS runner.
// The two builds disagreed, and scripts/native-proof.sh reported "rebuilt binary differs
// from package". This test drives the same disagreement through the one accused script.
func TestGoBuildIgnoresAmbientCGO(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	base := capability.WithoutEnvironment(os.Environ(), runbinary.Env)

	build := func(cgo, name string) []byte {
		t.Helper()
		path := filepath.Join(out, name)
		cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), "--mode", "artifact", root, path)
		cmd.Dir = root
		cmd.Env = append(capability.WithoutEnvironment(base, "CGO_ENABLED"), "CGO_ENABLED="+cgo)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("build with CGO_ENABLED=%s: %v\n%s", cgo, err, output)
		}
		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	if disabled, enabled := build("0", "cgo-off"), build("1", "cgo-on"); !bytes.Equal(disabled, enabled) {
		t.Fatalf("go-build.sh output depends on ambient CGO_ENABLED: %d bytes with cgo off, %d bytes with cgo on", len(disabled), len(enabled))
	}
}
