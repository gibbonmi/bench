package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
)

// wantProseTable renders the pass table this verb owes for the named paths. The block
// name, the field order, and the verdict cell are written here; the escaping derives
// through the encoder that renders the real output, so a path the encoder quotes cannot
// red the expectation for the wrong reason.
func wantProseTable(t *testing.T, paths ...string) string {
	t.Helper()
	rows := make([][]string, 0, len(paths))
	for _, path := range paths {
		rows = append(rows, []string{path, "pass"})
	}
	out, err := toon.Table("prose", []string{"path", "verdict"}, rows)
	if err != nil {
		t.Fatalf("render the expected prose table: %v", err)
	}
	return out
}

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
	if !strings.Contains(out, strings.TrimSpace(words(27))) {
		t.Fatalf("stdout = %q, want the offending sentence", out)
	}
	if strings.Contains(out, "prose[") {
		t.Fatalf("stdout = %q, want no pass table on a red list", out)
	}
}

// TestGateProseCommandCleanList exits 0 on a clean named file and states its verdict as a
// `prose[N]{path,verdict}` table, so a caller reads the pass rather than inferring it from
// the exit code.
func TestGateProseCommandCleanList(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")
	write(t, root, "docs/notes.md", "Short prose.\n")

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root, "--", "docs/notes.md"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if want := wantProseTable(t, "docs/notes.md"); stdout.String() != want {
		t.Fatalf("stdout = %q, want the pass table %q", stdout.String(), want)
	}
	if strings.Contains(stdout.String(), "green") {
		t.Fatalf("stdout = %q, want no green token: the lane states that a pass is not a graded green", stdout.String())
	}
}

// TestGateProseCommandEmptyPathList passes with no findings when the lane names no
// Markdown, and its pass table carries no row.
func TestGateProseCommandEmptyPathList(t *testing.T) {
	root := t.TempDir()
	write(t, root, ".bench/prose-exclusions", "")

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if want := wantProseTable(t); stdout.String() != want {
		t.Fatalf("stdout = %q, want the empty pass table %q", stdout.String(), want)
	}
	if strings.Contains(stdout.String(), "green") {
		t.Fatalf("stdout = %q, want no green token on an empty path list", stdout.String())
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
	if want := gateProseUsage + "\n"; stderr.String() != want {
		t.Fatalf("stderr = %q, want usage %q", stderr.String(), want)
	}
}

// TestGateProseCommandHelp prints the command contract without treating help as a refusal.
func TestGateProseCommandHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if want := gateProseUsage + "\n"; stdout.String() != want {
		t.Fatalf("stdout = %q, want usage %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
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

// TestGateProseCommandRefusesAMalformedExclusionRowWithNoSubject is PL14. The lane runs
// this verb on an empty path list when a commit changes the exclusion list alone, so the
// policy is graded before any subject is. A verb that read the policy only per subject
// would pass an empty list and let a malformed row land.
func TestGateProseCommandRefusesAMalformedExclusionRowWithNoSubject(t *testing.T) {
	root := t.TempDir()
	write(t, root, "docs/notes.md", "A short sentence.\n")
	write(t, root, ".bench/prose-exclusions", "docs/notes.md\n")

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "malformed exclusion row") {
		t.Fatalf("stdout = %q, want the malformed row named", stdout.String())
	}
}

// TestGateProseCommandFileRootIsUsageError covers the file-root refusal: the root operand
// must be a directory, so a file root is a malformed argument and never reaches the
// exclusion file.
func TestGateProseCommandFileRootIsUsageError(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "README.md", "Short prose.\n")
	root := filepath.Join(dir, "README.md")

	var stdout, stderr bytes.Buffer
	code := GateProseCommand([]string{root}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("exit = %d, want 2; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty on a usage error", stdout.String())
	}
	err := stderr.String()
	if !strings.Contains(err, root) {
		t.Fatalf("stderr = %q, want it to name the operand %q", err, root)
	}
	if !strings.Contains(err, "must be a directory") {
		t.Fatalf("stderr = %q, want it to say the root must be a directory", err)
	}
	if !strings.Contains(err, gateProseUsage) {
		t.Fatalf("stderr = %q, want the usage line %q", err, gateProseUsage)
	}
	if strings.Contains(err, "prose-exclusions") {
		t.Fatalf("stderr = %q, want no exclusion-file diagnostic for a file root", err)
	}
}
