package prose

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// write puts a fixture at the repository-relative path rel under root and makes every
// parent directory it needs.
func write(t *testing.T, root, rel, body string) string {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("make directory for %s: %v", rel, err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

// longSentence is a document with one sentence over the word bound. Every fixture that
// asks whether a subject is graded plants this text.
func longSentence() string { return words(40) + ".\n" }

// atLimit returns a clean document of exactly bounds.ControlRecordLimit bytes, so the
// walk grades a subject at the bound and refuses only the subject over it.
func atLimit() string {
	unit := "ok.\n\n"
	body := strings.Repeat(unit, int(bounds.ControlRecordLimit)/len(unit))
	return body + strings.Repeat("x", int(bounds.ControlRecordLimit)-len(body))
}

// TestGrade grades one temporary root per row: the walk, the exclusion grammar, and the
// subject classification each own their rows.
func TestGrade(t *testing.T) {
	for _, tt := range []struct {
		name    string
		build   func(*testing.T, string)
		count   int
		wantSub string
	}{
		{
			name: "PD22 skipped directories are not graded",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
				write(t, root, ".bench/prose-exclusions", "tests/canary/ fixtures carry planted content\n")
				for _, rel := range []string{
					"tests/canary/planted.md",
					"node_modules/pkg/readme.md",
					"dist/report.md",
					".git/note.md",
					"internal/testdata/sample.md",
				} {
					write(t, root, rel, longSentence())
				}
			},
		},
		{
			name: "PD43 Go and shell comments are not graded",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
				write(t, root, ".bench/prose-exclusions", "")
				write(t, root, "main.go", "package main\n\n// "+words(40)+".\nfunc main() {}\n")
				write(t, root, "run.sh", "#!/bin/sh\n# "+words(40)+".\n")
			},
		},
		{
			name: "PD20 a subject at the byte bound is graded",
			build: func(t *testing.T, root string) {
				write(t, root, "big.md", atLimit())
				write(t, root, ".bench/prose-exclusions", "")
			},
		},
		{
			name: "PD20 a subject over the byte bound is refused",
			build: func(t *testing.T, root string) {
				write(t, root, "big.md", atLimit()+"y")
				write(t, root, ".bench/prose-exclusions", "")
			},
			count:   1,
			wantSub: "refused",
		},
		{
			name: "PD20 a non-UTF-8 subject is refused",
			build: func(t *testing.T, root string) {
				write(t, root, "binary.md", "head\xff\xfetail")
				write(t, root, ".bench/prose-exclusions", "")
			},
			count:   1,
			wantSub: "refused",
		},
		{
			name: "PD20 a symbolic link to a file is refused",
			build: func(t *testing.T, root string) {
				target := write(t, t.TempDir(), "outside.md", longSentence())
				requireSymlink(t, target, filepath.Join(root, "link.md"))
				write(t, root, ".bench/prose-exclusions", "")
			},
			count:   1,
			wantSub: "refused",
		},
		{
			name: "PD20 a symbolic link to a directory is not reported",
			build: func(t *testing.T, root string) {
				outside := t.TempDir()
				write(t, outside, "inner.md", longSentence())
				write(t, root, "README.md", "Short prose.\n")
				requireSymlink(t, outside, filepath.Join(root, "linked.md"))
				write(t, root, ".bench/prose-exclusions", "")
			},
		},
		{
			name: "PD20 a FIFO is refused",
			build: func(t *testing.T, root string) {
				requireFifo(t, filepath.Join(root, "pipe.md"))
				write(t, root, ".bench/prose-exclusions", "")
			},
			count:   1,
			wantSub: "refused",
		},
		{
			name: "a graded subject reports its path and counts",
			build: func(t *testing.T, root string) {
				write(t, root, "docs/guide.md", longSentence())
				write(t, root, ".bench/prose-exclusions", "")
			},
			count:   1,
			wantSub: `"docs/guide.md" line 1: sentence of 40 words`,
		},
		{
			name: "an excluded file is not graded",
			build: func(t *testing.T, root string) {
				write(t, root, "docs/guide.md", longSentence())
				write(t, root, ".bench/prose-exclusions", "docs/guide.md the record keeps its text\n")
			},
		},
		{
			name: "a root with no subject yields no diagnostic",
			build: func(t *testing.T, root string) {
				write(t, root, "main.go", "package main\n")
			},
		},
		{
			name: "an absent exclusion file reds",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
			},
			count:   1,
			wantSub: "the exclusion file is absent",
		},
		{
			name: "a malformed row reds",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
				write(t, root, ".bench/prose-exclusions", "# a comment\ndocs/\n")
			},
			count:   1,
			wantSub: "malformed exclusion row",
		},
		{
			name: "a duplicate row reds",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
				write(t, root, ".bench/prose-exclusions", "README.md first reason\nREADME.md second reason\n")
			},
			count:   1,
			wantSub: "duplicate exclusion row",
		},
		{
			name: "a glob row reds",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
				write(t, root, ".bench/prose-exclusions", "docs/*.md a glob is not a path\n")
			},
			count:   1,
			wantSub: "uses a glob character",
		},
		{
			name: "a row that names an absent path reds",
			build: func(t *testing.T, root string) {
				write(t, root, "README.md", "Short prose.\n")
				write(t, root, ".bench/prose-exclusions", "docs/gone.md the file left\n")
			},
			count:   1,
			wantSub: "names an absent path",
		},
		{
			name: "C2 a directory row with no trailing slash reds",
			build: func(t *testing.T, root string) {
				write(t, root, "docs/guide.md", longSentence())
				write(t, root, ".bench/prose-exclusions", "docs the record keeps its text\n")
			},
			count:   1,
			wantSub: "directory row needs a trailing slash",
		},
		{
			name: "a row with no trailing newline parses",
			build: func(t *testing.T, root string) {
				write(t, root, "docs/guide.md", longSentence())
				write(t, root, ".bench/prose-exclusions", "docs/ the record keeps its text")
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			tt.build(t, root)
			got := Grade(root)
			if len(got) != tt.count {
				t.Fatalf("Grade() = %q, want %d diagnostics", got, tt.count)
			}
			if tt.wantSub != "" && !strings.Contains(got[0], tt.wantSub) {
				t.Errorf("Grade() = %q, want a diagnostic that holds %q", got[0], tt.wantSub)
			}
		})
	}
}

// requireSymlink makes one link or reports the host capability the assertion needs.
func requireSymlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symbolic links unavailable on this filesystem: %v", err))
	}
}

// requireFifo makes one FIFO or reports the host capability the assertion needs. The
// FIFO has no writer, so an implementation that opens before it checks the type blocks
// in open(2) rather than returning a refusal.
func requireFifo(t *testing.T, path string) {
	t.Helper()
	if err := syscall.Mkfifo(path, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
	}
}
