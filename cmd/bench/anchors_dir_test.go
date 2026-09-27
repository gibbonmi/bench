package main

import (
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/toon"
)

// anchorsDirFields is the directory form's table header.
var anchorsDirFields = []string{"file", "anchors", "diagnostics", "verdict"}

// expectedAnchoredFiles is the directory form's independent file expectation: each
// distinct registered file below dir, sorted by path. It reads anchors.Entries, not the
// production inventory, because a production inventory that drops a file must turn these
// tests red; the ticket's probe records that red.
func expectedAnchoredFiles(dir string) []string {
	var files []string
	for _, anchor := range anchors.Entries() {
		if dir != "." && !strings.HasPrefix(anchor.File, dir+"/") || slices.Contains(files, anchor.File) {
			continue
		}
		files = append(files, anchor.File)
	}
	slices.Sort(files)
	return files
}

// registeredAnchorCount is the anchor cell of a file's row: the number of registry
// entries that name the file.
func registeredAnchorCount(file string) int {
	count := 0
	for _, anchor := range anchors.Entries() {
		if anchor.File == file {
			count++
		}
	}
	return count
}

// expectedAnchorsDirTable renders the expected rows for dir under root. A file in verdicts
// takes its literal verdict; every other file reads pass with no diagnostic and fail
// otherwise.
func expectedAnchorsDirTable(t *testing.T, root, dir string, verdicts map[string]string) string {
	t.Helper()
	var rows [][]any
	for _, file := range expectedAnchoredFiles(dir) {
		count := registeredAnchorCount(file)
		diagnostics := len(anchors.EvaluatePath(root, file).Diagnostics)
		verdict, literal := verdicts[file]
		if !literal {
			verdict = "fail"
			if diagnostics == 0 {
				verdict = "pass"
			}
		}
		rows = append(rows, []any{file, count, diagnostics, verdict})
	}
	want, err := toon.TableTyped("files", anchorsDirFields, rows)
	if err != nil {
		t.Fatal(err)
	}
	return want
}

// firstNestedAnchorDir answers the directory of the first registered file that sits below
// the repository root, so the nested case follows the registry as it grows.
func firstNestedAnchorDir(t *testing.T) string {
	t.Helper()
	for _, anchor := range anchors.Entries() {
		if strings.Contains(anchor.File, "/") {
			return path.Dir(anchor.File)
		}
	}
	t.Fatal("registry has no nested anchor file")
	return ""
}

func TestAnchorsDirectoryGradesEachAnchoredFile(t *testing.T) {
	nested := firstNestedAnchorDir(t)
	for _, operand := range []string{".", nested, "./" + nested + "/"} {
		t.Run(operand, func(t *testing.T) {
			dir := path.Clean(operand)
			root := newAXIEnvelopeRepo(t)
			writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), readAnchorsFixture(t, "AGENTS.md"))
			writeAXIFixture(t, filepath.Join(root, filepath.FromSlash(nested), ".keep"), "fixture\n")
			if len(expectedAnchoredFiles(dir)) < 2 {
				t.Fatalf("directory %q holds fewer than two anchored files; the case cannot show one row per file", dir)
			}
			result := runAXICommandAt(t, root, []string{"anchors", operand})
			if result.code != 1 || result.stderr != "" {
				t.Fatalf("anchors %s = %#v, want exit 1 for the absent files and no stderr", operand, result)
			}
			verdicts := map[string]string{}
			if dir == "." {
				verdicts["AGENTS.md"] = "pass"
			}
			want := expectedAnchorsDirTable(t, root, dir, verdicts)
			if !strings.HasPrefix(result.stdout, want) {
				t.Fatalf("anchors %s stdout = %q, want the table %q", operand, result.stdout, want)
			}
			if rest := strings.TrimPrefix(result.stdout, want); !strings.HasPrefix(rest, "help[1]{cmd,why}:\n") || !strings.Contains(rest, "bench anchors <file>") {
				t.Fatalf("anchors %s help = %q, want one row that names the file form", operand, rest)
			}
			if dir != "." && strings.Contains(result.stdout, "\n  AGENTS.md,") {
				t.Fatalf("anchors %s stdout = %q, want no row for a file outside the directory", operand, result.stdout)
			}
		})
	}
}

func TestAnchorsDirectoryReportsARefusedFile(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	target := filepath.Join(root, "target.md")
	writeAXIFixture(t, target, readAnchorsFixture(t, "AGENTS.md"))
	requireAnchorSymlink(t, target, filepath.Join(root, "AGENTS.md"))
	result := runAXICommandAt(t, root, []string{"anchors", "."})
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("anchors . = %#v, want exit 1 and no stderr", result)
	}
	row, err := toon.TableTyped("files", anchorsDirFields, [][]any{{"AGENTS.md", registeredAnchorCount("AGENTS.md"), 1, "refused"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Split(row, "\n")[1] + "\n"; !strings.Contains(result.stdout, want) {
		t.Fatalf("anchors . stdout = %q, want the refused row %q", result.stdout, want)
	}
}

func TestAnchorsDirectoryWithNoAnchoredFileAnswersAnEmptySet(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "nested", "guide.md"), "# An ordinary guide\n")
	for _, operand := range []string{"nested", "nested/deep"} {
		result := runAXICommandAt(t, root, []string{"anchors", operand})
		if want := "files[0]{file,anchors,diagnostics,verdict}:\nhelp[0]{cmd,why}:\n"; result.code != 0 || result.stderr != "" || result.stdout != want {
			t.Fatalf("anchors %s = %#v, want exit 0 and the explicit empty set %q", operand, result, want)
		}
	}
}
