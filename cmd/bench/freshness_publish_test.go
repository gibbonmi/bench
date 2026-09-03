package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

func TestFreshnessPublishRequiresRootDestinationAndVersion(t *testing.T) {
	for _, args := range [][]string{nil, {"root"}, {"root", "output"}, {"root", "output", "1.2.3", "other-stage"}} {
		var stderr bytes.Buffer
		if code := freshnessPublish(args, "ignored", &stderr); code != 2 {
			t.Fatalf("freshnessPublish(%q) = %d, want usage refusal; stderr=%q", args, code, stderr.String())
		}
	}

}

func TestFreshnessCheckRefusesMissingOwnExecutable(t *testing.T) {
	root := t.TempDir()
	var stderr bytes.Buffer

	code := freshnessCheck([]string{root}, filepath.Join(root, "dist", "bench"), &stderr)
	if code != 1 {
		t.Fatalf("freshnessCheck missing executable exit = %d, want 1", code)
	}
	want := freshness.RebuildAction(root)
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("freshnessCheck stderr = %q, want rebuild action %q", stderr.String(), want)
	}
}

func TestFreshnessPublishBindsTheInvokedExecutable(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	staged := filepath.Join(t.TempDir(), "staged bench")
	binary, err := os.ReadFile(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(staged, binary, 0o755); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "published bench")
	var stderr bytes.Buffer
	if code := freshnessPublish([]string{root, output, "1.2.3"}, staged, &stderr); code != 0 {
		t.Fatalf("freshnessPublish = %d, stderr=%q", code, stderr.String())
	}
	if err := freshness.Verify(root, output); err != nil {
		t.Fatalf("published executable does not match its seal: %v", err)
	}
}
