package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestGoBuildRefusesWhenFlagRegistryReadFails grades the read of the Go build flags. A
// failed read that the builder does not see leaves the flag list empty, and go build then
// runs without the registry's reproducibility flags and version stamp. That build exits
// zero, so only the refusal keeps an unreproducible binary out of dist/. The shim fails only the
// registry read, so the version read and the toolchain stay real.
func TestGoBuildRefusesWhenFlagRegistryReadFails(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	root = gittest.KitCopy(t, root)
	node, err := exec.LookPath("node")
	if err != nil {
		t.Fatalf("node is required to run the builder: %v", err)
	}
	shimDir := t.TempDir()
	shim := fmt.Sprintf("#!/bin/sh\nfor arg do\n  case \"$arg\" in\n    */'%s') echo 'shim: registry read refused' >&2; exit 1 ;;\n  esac\ndone\nexec '%s' \"$@\"\n",
		buildInput(t, root, "go_requirements"), node)
	if err := os.WriteFile(filepath.Join(shimDir, "node"), []byte(shim), 0o755); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "bench")

	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, out)
	cmd.Dir = root
	env := capability.WithoutEnvironment(capability.WithoutEnvironment(os.Environ(), runbinary.Env), "PATH")
	cmd.Env = append(env, "PATH="+shimDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("build with a failed flag registry read exited zero\n%s", output)
	}
	for _, path := range []string{out, out + ".seal"} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("build with a failed flag registry read left %q: %v\n%s", path, err, output)
		}
	}
}
