// Package testrepo owns reusable repository materialization for behavioral tests.
package testrepo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CommitWorkingTree copies tracked and visible untracked files into destination. It
// commits that exact snapshot, so tests can authenticate the tree they are grading.
func CommitWorkingTree(source, destination string) error {
	listed := exec.Command("git", "-C", source, "ls-files", "-z", "--cached", "--others", "--exclude-standard")
	output, err := listed.Output()
	if err != nil {
		return err
	}
	for _, relative := range strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00") {
		if relative == "" {
			continue
		}
		from, to := filepath.Join(source, relative), filepath.Join(destination, relative)
		info, err := os.Lstat(from)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(from)
			if err == nil {
				err = os.Symlink(target, to)
			}
			if err != nil {
				return err
			}
			continue
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("source entry is not regular: %s", relative)
		}
		data, err := os.ReadFile(from)
		if err == nil {
			err = os.WriteFile(to, data, info.Mode().Perm())
		}
		if err != nil {
			return err
		}
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "testrepo@example.invalid"},
		{"config", "user.name", "testrepo"},
	} {
		command := exec.Command("git", append([]string{"-C", destination}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("git %v: %w: %s", args, err, strings.TrimSpace(string(output)))
		}
	}
	return CommitAll(destination, "authenticated test snapshot")
}

// CommitAll authenticates the current visible working tree in an initialized test repo.
func CommitAll(root, message string) error {
	for _, args := range [][]string{{"add", "-f", "."}, {"commit", "-qm", message}} {
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("git %v: %w: %s", args, err, strings.TrimSpace(string(output)))
		}
	}
	return nil
}
