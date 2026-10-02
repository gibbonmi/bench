package repairtest

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/gittest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepairRefusesPublicNamespace(t *testing.T) {
	s := linkedSession(t)
	hook := filepath.Join(s.root, ".codex", "hooks.json")
	original, err := os.ReadFile(hook)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	namespace := filepath.Join(s.home, "bench-state", "compatibility-repairs")
	if err := os.MkdirAll(namespace, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(namespace, 0o755); err != nil {
		t.Fatal(err)
	}
	_, stdout, stderr := s.doctor("--fix")
	if data, err := os.ReadFile(hook); err == nil && bytes.Equal(data, original) {
		t.Fatalf("repair published through a non-private recovery namespace: %s %s", stdout, stderr)
	}
	if _, err := os.Lstat(hook); !os.IsNotExist(err) {
		t.Fatalf("refused namespace changed target: %v", err)
	}
	info, err := os.Stat(namespace)
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("repair changed the existing namespace mode: %v", err)
	}
}

func TestRepairPreservesRepositoryRefs(t *testing.T) {
	s := linkedSession(t)
	s.track(t)
	branch := strings.TrimSpace(gittest.Output(t, s.root, "symbolic-ref", "--short", "HEAD"))
	gittest.Output(t, s.root, "remote", "add", "origin", s.root)
	gittest.Output(t, s.root, "update-ref", "refs/remotes/origin/"+branch, "HEAD")
	before := gittest.Output(t, s.root, "for-each-ref")
	_, stdout, stderr := s.doctor("--fix")
	after := gittest.Output(t, s.root, "for-each-ref")
	if after != before {
		t.Fatalf("repair changed repository refs outside its recovery record: %s %s", stdout, stderr)
	}
}
