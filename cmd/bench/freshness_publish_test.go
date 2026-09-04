package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/freshness"
)

func TestFreshnessPublishRequiresRootDestinationManifestDirAndVersion(t *testing.T) {
	for _, args := range [][]string{nil, {"root"}, {"root", "output"}, {"root", "output", "manifest-dir"}, {"root", "output", "manifest-dir", "1.2.3", "other-stage"}} {
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
	manifestDir := t.TempDir()
	var stderr bytes.Buffer
	if code := freshnessPublish([]string{root, output, manifestDir, "1.2.3"}, staged, &stderr); code != 0 {
		t.Fatalf("freshnessPublish = %d, stderr=%q", code, stderr.String())
	}
	if err := freshness.Verify(root, output); err != nil {
		t.Fatalf("published executable does not match its seal: %v", err)
	}
	// The verb has to hand its third operand through to the publisher. A manifest that
	// landed beside the output instead would leave every reader of the wrapper's own
	// directory with nothing.
	if _, err := brokermanifest.Read(filepath.Join(manifestDir, brokermanifest.Name)); err != nil {
		t.Fatalf("freshnessPublish did not publish the manifest in its named directory: %v", err)
	}
}
