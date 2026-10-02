package repairtest

import (
	"bytes"
	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/brokermanifest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func kitSession(t *testing.T) (session, string, string) {
	t.Helper()
	s := newSession(t)
	t.Setenv("BENCH_KIT", s.root)
	bin := filepath.Join(s.home, ".local", "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	wrapper := filepath.Join(s.root, "bench-wrapper")
	source := []byte("#!/bin/sh\n# kit source sentinel\n")
	if err := os.WriteFile(wrapper, source, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_WRAPPER", wrapper)
	return s, wrapper, bin
}

func TestCompatibilityKitRepair(t *testing.T) {
	s, wrapper, bin := kitSession(t)
	source, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := s.doctor("--fix")
	shim := filepath.Join(bin, "bench")
	data, err := os.ReadFile(shim)
	if code != 1 || err != nil || string(data) != adopt.ShimContent(wrapper)+"\n" {
		t.Fatalf("kit repair did not install canonical shim: %d, %v, %s %s", code, err, stdout, stderr)
	}
	id := repairID(t, stdout)
	targets := []string{shim, filepath.Join(s.root, brokermanifest.Name), filepath.Join(s.root, ".git", "hooks", "pre-push")}
	for _, path := range targets {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("kit repair omitted %s: %v", path, err)
		}
	}
	code, stdout, stderr = s.doctor("--undo", id)
	for _, path := range targets {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Errorf("kit undo did not restore absence for %s: %v; %d %s %s", path, err, code, stdout, stderr)
		}
	}
	data, err = os.ReadFile(wrapper)
	if err != nil || !bytes.Equal(data, source) {
		t.Fatalf("kit source changed: %q, %v", data, err)
	}
	if _, err := os.Lstat(filepath.Join(s.root, ".bench", "link-manifest.tsv")); !os.IsNotExist(err) {
		t.Fatalf("kit repair installed consumer payload: %v", err)
	}
}

func TestCompatibilityKitModifiedConflict(t *testing.T) {
	for _, kind := range []string{"shim mode", "shim content", "hook mode", "hook content", "broker content"} {
		t.Run(kind, func(t *testing.T) {
			s, _, bin := kitSession(t)
			_, stdout, stderr := s.doctor("--fix")
			if stderr != "" {
				t.Fatalf("kit fixture repair: %s %s", stdout, stderr)
			}
			path := filepath.Join(bin, "bench")
			if strings.HasPrefix(kind, "hook") {
				path = filepath.Join(s.root, ".git", "hooks", "pre-push")
			}
			if strings.HasPrefix(kind, "broker") {
				path = filepath.Join(s.root, brokermanifest.Name)
			}
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasSuffix(kind, "content") {
				before = append(before, []byte("\nuser change\n")...)
			}
			if err := os.WriteFile(path, before, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, 0o600); err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := s.doctor("--fix")
			after, err := os.ReadFile(path)
			info, statErr := os.Stat(path)
			if err != nil || statErr != nil || !bytes.Equal(after, before) || info.Mode().Perm() != 0o600 {
				t.Fatalf("kit repair changed modified asset: %s %v %v", path, err, statErr)
			}
			if code != 1 || !strings.Contains(stdout+stderr, "modified-managed") {
				t.Fatalf("kit repair omitted modified conflict for %s: %d %s %s", kind, code, stdout, stderr)
			}
		})
	}
}

func TestCompatibilityKitUpdatesCanonicalBroker(t *testing.T) {
	s, _, _ := kitSession(t)
	_, stdout, stderr := s.doctor("--fix")
	if stderr != "" {
		t.Fatalf("kit fixture repair: %s %s", stdout, stderr)
	}
	broker, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	path, _, err := brokermanifest.Write(s.root, broker, "0.9.0")
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, stdout, stderr = s.doctor("--fix")
	if stderr != "" {
		t.Fatalf("repair refused an unchanged canonical broker manifest: %s %s", stdout, stderr)
	}
	id := repairID(t, stdout)
	after, err := os.ReadFile(path)
	if err != nil || bytes.Equal(before, after) {
		t.Fatalf("repair did not update canonical broker binding: %v", err)
	}
	_, stdout, stderr = s.doctor("--undo", id)
	restored, err := os.ReadFile(path)
	if stderr != "" || err != nil || !bytes.Equal(restored, before) {
		t.Fatalf("undo lost canonical broker preimage: %v, %s %s", err, stdout, stderr)
	}
}

func TestCompatibilityKitRepairsAliasedShim(t *testing.T) {
	s, wrapper, bin := kitSession(t)
	alias := filepath.Join(s.home, "stable-bin")
	if err := os.Symlink(bin, alias); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", alias+string(os.PathListSeparator)+os.Getenv("PATH"))
	path := filepath.Join(bin, "bench")
	before := []byte(adopt.ShimContent(filepath.Join(s.root, "previous-wrapper")) + "\n")
	if err := os.WriteFile(path, before, 0o755); err != nil {
		t.Fatal(err)
	}
	_, stdout, stderr := s.doctor("--fix")
	after, err := os.ReadFile(path)
	if stderr != "" || err != nil || string(after) != adopt.ShimContent(wrapper)+"\n" {
		t.Fatalf("repair refused canonical shim through directory alias: %v %s %s", err, stdout, stderr)
	}
	id := repairID(t, stdout)
	_, stdout, stderr = s.doctor("--undo", id)
	restored, err := os.ReadFile(path)
	if stderr != "" || err != nil || !bytes.Equal(restored, before) {
		t.Fatalf("undo did not restore aliased shim: %v %s %s", err, stdout, stderr)
	}
}
