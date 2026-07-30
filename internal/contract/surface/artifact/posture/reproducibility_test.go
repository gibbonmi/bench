package posture

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

func TestGoBuildIgnoresCheckoutTopology(t *testing.T) {
	root := contract.SubjectRoot(t)
	clone := filepath.Join(t.TempDir(), "isolated-source")
	command := exec.Command("git", "clone", "--quiet", "--no-hardlinks", root, clone)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("clone isolated source: %v\n%s", err, output)
	}
	// Mirror the working tree onto the clone so uncommitted DELETIONS (a file removed in
	// the source but still tracked at HEAD) are reflected too — a tar overlay can only add
	// or modify, leaving a deleted source file resurrected in the clone and drifting the
	// binary. --delete makes the clone tree a true mirror; .git and dist stay untouched.
	if _, err := exec.LookPath("rsync"); err != nil {
		capability.Capability(t, capability.Tool, fmt.Sprintf("reproducibility probe needs rsync on PATH: %v", err))
	}
	overlay := exec.Command("rsync", "-a", "--delete", "--exclude=/.git", "--exclude=/dist", root+"/", clone)
	if output, err := overlay.CombinedOutput(); err != nil {
		t.Fatalf("overlay source snapshot: %v\n%s", err, output)
	}
	outputs := []string{filepath.Join(t.TempDir(), "worktree-binary"), filepath.Join(t.TempDir(), "clone-binary")}
	for index, source := range []string{root, clone} {
		build := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), source, outputs[index])
		build.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=darwin", "GOARCH=arm64")
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build checkout %d: %v\n%s", index, err, output)
		}
	}
	first, err := os.ReadFile(outputs[0])
	if err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(outputs[1])
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(first) != sha256.Sum256(second) {
		t.Fatal("checkout topology changed release binary bytes")
	}
}
