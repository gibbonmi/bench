package env

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestGitPolicyHasOneSource(t *testing.T) {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatal(err)
	}
	root := strings.TrimSpace(string(output))
	findings, err := gitPolicyFindings(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range findings {
		t.Error(finding)
	}
}

func gitPolicyFindings(root string) ([]string, error) {
	files, err := exec.Command("git", "-C", root, "ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", "*.go").Output()
	if err != nil {
		return nil, err
	}
	var findings []string
	for _, name := range strings.Split(string(files), "\x00") {
		if name == "" || name == "internal/env/kit_run.go" || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(name)))
		if file.State == bounds.StateAbsent {
			continue
		}
		if file.State.Failed() {
			return nil, fmt.Errorf("%s: %s", name, file.Reason)
		}
		if bytes.Contains(file.Data, []byte("maintenance.auto")) {
			findings = append(findings, "git maintenance policy outside its owner: "+name)
		}
	}
	return findings, nil
}

func TestGitPolicyScanRejectsSpecialSources(t *testing.T) {
	for _, kind := range []string{"symlink", "fifo"} {
		t.Run(kind, func(t *testing.T) {
			root := gittest.Repo(t)
			path := filepath.Join(root, "source.go")
			// Git omits an untracked FIFO, so retain the source path in its index first.
			if err := os.WriteFile(path, nil, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := exec.Command("git", "-C", root, "add", "source.go").Run(); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if kind == "symlink" {
				target := filepath.Join(t.TempDir(), "target")
				if err := os.WriteFile(target, nil, 0o600); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, err.Error())
				}
			} else if err := syscall.Mkfifo(path, 0o600); err != nil {
				capability.Capability(t, capability.Fifo, err.Error())
			}
			if _, err := gitPolicyFindings(root); err == nil || !strings.Contains(err.Error(), "source.go") {
				t.Fatalf("special source refusal = %v, want an error naming source.go", err)
			}
		})
	}
}
