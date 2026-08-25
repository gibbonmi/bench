package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// write puts a fixture at the repository-relative path rel under root and makes every
// parent directory it needs.
func write(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("make directory for %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

// words returns a document with one sentence of n words, over the 25-word bound at
// n=27.
func words(n int) string {
	fields := make([]string, n)
	for i := range fields {
		fields[i] = "word" + strconv.Itoa(i)
	}
	return strings.Join(fields, " ") + ".\n"
}

// TestGateProseCommandFindsAnOverLongSentence is OG33: bench gate-prose <root> -- <file>
// exits 1 and names the file and the line for a 27-word sentence.
func TestGateProseCommandFindsAnOverLongSentence(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")
	write(t, root, "docs/notes.md", words(27))

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root, "--", "docs/notes.md"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("exit = %d, want 1; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, `"docs/notes.md"`) {
		t.Fatalf("stdout = %q, want it to name docs/notes.md", out)
	}
	if !strings.Contains(out, "line 1") {
		t.Fatalf("stdout = %q, want it to name line 1", out)
	}
}

// TestGateProseCommandCleanList exits 0 on a clean named file.
func TestGateProseCommandCleanList(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")
	write(t, root, "docs/notes.md", "Short prose.\n")

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root, "--", "docs/notes.md"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

// TestGateProseCommandEmptyPathList passes with no findings when the lane names no
// Markdown.
func TestGateProseCommandEmptyPathList(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

// TestGateProseCommandUnknownFlagIsUsageError exits 2 and never reads as a clean list.
func TestGateProseCommandUnknownFlagIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{"--bogus"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("exit = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty on a usage error", stdout.String())
	}
}

// TestGateProseCommandMissingRootIsUsageError exits 2 on an empty argv.
func TestGateProseCommandMissingRootIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := GateProseCommand(nil, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("exit = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}
