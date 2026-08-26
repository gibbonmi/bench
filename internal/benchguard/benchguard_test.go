package benchguard

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestClassifyFollowOns(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(path string) (string, error) {
		if path == "/work/kit-command" {
			return "/work/bin/bench.sh", nil
		}
		return "", errors.New("not bench")
	}}
	for _, command := range []string{"bench help | touch marker", "bench gate 2>&1", "X=1 bench help && touch marker", "bash -lc 'bench help || touch marker'", "bench worktree exec a -- bench gate > marker"} {
		if !Classify(command, resolver) {
			t.Errorf("Classify(%q) = false, want true", command)
		}
	}
	for _, command := range []string{"bench gate --fresh", "rg bench AGENTS.md", "bench worktree exec a -- bash -lc 'go test && go vet'", "cat <<'EOF'\nbench gate | tail\nEOF"} {
		if Classify(command, resolver) {
			t.Errorf("Classify(%q) = true, want false", command)
		}
	}
}

func TestClassifyResolvesBareAliasFromProcessPath(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "bench.sh")
	if err := os.WriteFile(target, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(dir, "kit-command")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	if !Classify("kit-command help | touch marker", DefaultResolver()) {
		t.Fatal("bare PATH alias did not classify as Bench")
	}
}

func TestCommandFromEnvelope(t *testing.T) {
	if _, err := CommandFromEnvelope([]byte(`{"tool_input":{"command":"bench help\\u0026\\u0026 touch marker"}}`)); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{`{`, `{}`, `{"tool_input":{"command":1}}`, `{"tool_input":{"command":""}}`, `{"tool_input":{"command":"\u0000"}}`} {
		if _, err := CommandFromEnvelope([]byte(input)); err == nil {
			t.Errorf("CommandFromEnvelope(%q) succeeded", input)
		}
	}
}

// TestInvokesBenchWalksOneWrapperLevel proves the exported Bench test reads a
// wrapped string, so a `bash -c` call that runs Bench is a Bench call.
func TestInvokesBenchWalksOneWrapperLevel(t *testing.T) {
	resolver := Resolver{Getwd: func() (string, error) { return "/work", nil }, EvalSymlinks: func(string) (string, error) { return "", errors.New("not bench") }}
	for _, command := range []string{"bench status", "bash -c 'bench status'", "bench worktree exec a -- sed -i x /pool/id/b", "sh -c 'cd /tmp && bench gate'"} {
		if !InvokesBench(command, resolver) {
			t.Errorf("InvokesBench(%q) = false, want true", command)
		}
	}
	for _, command := range []string{"ls", "rg bench AGENTS.md", "bash -c 'ls /pool/id'"} {
		if InvokesBench(command, resolver) {
			t.Errorf("InvokesBench(%q) = true, want false", command)
		}
	}
}
