// Tests for the two file-resolving checks: seam-cell test citations and the
// review-pickup fence member.
package coverage

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

// citedSpec writes a folder spec at <root>/specs/<slug>/spec.md whose single mapped
// row carries seam as its seam cell, and returns that spec's path. fences, when
// non-empty, is appended verbatim so a test can shape the section it grades.
func citedSpec(t *testing.T, slug, seam, fences string) (root, specPath string) {
	t.Helper()
	root = t.TempDir()
	body := "# " + slug + "\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced +
		"| 1 | b | " + seam + " | w |\n" + fences
	specPath = filepath.Join(root, "specs", slug, "spec.md")
	writeUnder(t, specPath, body)
	return root, specPath
}

// writeUnder writes content at an absolute path, creating its parents.
func writeUnder(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", dir, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

// checkFilesOf parses specPath and grades it, the way ParseSpec and Command both do.
func checkFilesOf(t *testing.T, specPath string) []string {
	t.Helper()
	content, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", specPath, err)
	}
	return CheckFiles(parse(content), specPath)
}

// hasViolation reports whether any violation holds every fragment.
func hasViolation(v []string, fragments ...string) bool {
	for _, msg := range v {
		all := true
		for _, f := range fragments {
			if !strings.Contains(msg, f) {
				all = false
			}
		}
		if all {
			return true
		}
	}
	return false
}

// TestCheckFilesRejectAnUndeclaredCitedTestName is the citation check's red case: the
// cited file exists, but it declares no function of the cited name, so the row's
// evidence does not exist.
func TestCheckFilesRejectAnUndeclaredCitedTestName(t *testing.T) {
	root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestAbsent`)", "")
	writeUnder(t, filepath.Join(root, "internal", "x", "foo_test.go"), "package x\n\nfunc TestPresent(t *testing.T) {}\n")

	v := checkFilesOf(t, specPath)
	if !hasViolation(v, "coverage map row 1", "TestAbsent") {
		t.Fatalf("CheckFiles = %#v, want a violation naming row 1 and TestAbsent", v)
	}
}

// TestCheckFilesRejectAnAbsentCitedFile pins the other half of the citation check: a
// row may cite a file that was never written, or was renamed after the map was.
func TestCheckFilesRejectAnAbsentCitedFile(t *testing.T) {
	_, specPath := citedSpec(t, "s", "`internal/x/gone_test.go` (`TestPresent`)", "")

	v := checkFilesOf(t, specPath)
	if !hasViolation(v, "coverage map row 1", "internal/x/gone_test.go") {
		t.Fatalf("CheckFiles = %#v, want a violation naming row 1 and the absent file", v)
	}
}

// TestCheckFilesAcceptADeclaredCitation is the citation check's green case, including
// a subtest citation whose leading segment is the declared function.
func TestCheckFilesAcceptADeclaredCitation(t *testing.T) {
	root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent`, `TestOther/a case`)", "")
	writeUnder(t, filepath.Join(root, "internal", "x", "foo_test.go"),
		"package x\n\nfunc TestPresent(t *testing.T) {}\n\nfunc TestOther(t *testing.T) {}\n")

	if v := checkFilesOf(t, specPath); len(v) != 0 {
		t.Fatalf("CheckFiles = %#v, want no violation", v)
	}
}

// TestCheckFilesIgnoreARowThatCitesNoTestFile pins the no-citation case. A
// review-owned or prose seam cell names no file, so it can add no citation violation.
func TestCheckFilesIgnoreARowThatCitesNoTestFile(t *testing.T) {
	_, specPath := citedSpec(t, "s", "review-owned: the Standards axis reads the type", "")

	if v := checkFilesOf(t, specPath); len(v) != 0 {
		t.Fatalf("CheckFiles = %#v, want no violation", v)
	}
}

// TestCheckFilesRequireTheReviewPickupFenceMember pins the pickup rule: a folder spec
// that declares ownership fences must authorize its own review file, or the review
// cannot land with the build.
func TestCheckFilesRequireTheReviewPickupFenceMember(t *testing.T) {
	_, specPath := citedSpec(t, "s", "cli seam", "\n## Ownership fences\n\n- `internal/x/`\n")

	v := checkFilesOf(t, specPath)
	if !hasViolation(v, "reviews/s.md") {
		t.Fatalf("CheckFiles = %#v, want a violation naming reviews/s.md", v)
	}
}

// TestCheckFilesAcceptTheReviewPickupFenceMember is the pickup check's green case.
func TestCheckFilesAcceptTheReviewPickupFenceMember(t *testing.T) {
	_, specPath := citedSpec(t, "s", "cli seam", "\n## Ownership fences\n\n- `internal/x/`\n- `reviews/s.md`\n")

	if v := checkFilesOf(t, specPath); len(v) != 0 {
		t.Fatalf("CheckFiles = %#v, want no violation", v)
	}
}

// TestCheckFilesRejectAPickupBelowASubsection pins the shared parser's section bound
// at this check. The fence section ends at any level-2-or-deeper heading, so a pickup
// written under a subsection is outside the section and authorizes nothing.
func TestCheckFilesRejectAPickupBelowASubsection(t *testing.T) {
	_, specPath := citedSpec(t, "s", "cli seam", "\n## Ownership fences\n\n### Paths\n\n- `internal/x/`\n- `reviews/s.md`\n")

	v := checkFilesOf(t, specPath)
	if !hasViolation(v, "reviews/s.md") {
		t.Fatalf("CheckFiles = %#v, want a violation naming reviews/s.md", v)
	}
}

// TestCheckFilesSkipASpecWithNoFenceSection pins that the pickup rule grades only a
// declared section. A spec that declares no fences is not yet at that stage.
func TestCheckFilesSkipASpecWithNoFenceSection(t *testing.T) {
	_, specPath := citedSpec(t, "s", "cli seam", "")

	if v := checkFilesOf(t, specPath); len(v) != 0 {
		t.Fatalf("CheckFiles = %#v, want no violation", v)
	}
}

// TestCheckFilesExemptAHistoricalSpec pins the historical opt-out over both new
// checks, exactly as it covers every other one.
func TestCheckFilesExemptAHistoricalSpec(t *testing.T) {
	root := t.TempDir()
	body := "# s\n\n" + historicalMarker + "\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced +
		"| 1 | b | `internal/x/gone_test.go` (`TestAbsent`) | w |\n\n## Ownership fences\n\n- `internal/x/`\n"
	specPath := filepath.Join(root, "specs", "s", "spec.md")
	writeUnder(t, specPath, body)

	if v := checkFilesOf(t, specPath); len(v) != 0 {
		t.Fatalf("CheckFiles = %#v, want no violation for a historical spec", v)
	}
}

// TestCommandCheckExitsOneOnAnUnresolvedCitation pins the citation check at the
// surface a caller reads: `bench coverage --check` exits 1 and names the row and the
// name it could not resolve.
func TestCommandCheckExitsOneOnAnUnresolvedCitation(t *testing.T) {
	root, _ := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestAbsent`)", "")
	writeUnder(t, filepath.Join(root, "internal", "x", "foo_test.go"), "package x\n\nfunc TestPresent(t *testing.T) {}\n")
	t.Chdir(root)

	out, code := Command([]string{"--check", "specs/s/spec.md"})
	if code != 1 || !strings.Contains(out, "row 1") || !strings.Contains(out, "TestAbsent") {
		t.Fatalf("Command = (%d, %q), want exit 1 naming row 1 and TestAbsent", code, out)
	}
}

// TestCommandCheckExitsOneOnAMissingReviewPickup pins the pickup check at the same
// surface, naming the fence member the spec omits.
func TestCommandCheckExitsOneOnAMissingReviewPickup(t *testing.T) {
	root, _ := citedSpec(t, "s", "cli seam", "\n## Ownership fences\n\n- `internal/x/`\n")
	t.Chdir(root)

	out, code := Command([]string{"--check", "specs/s/spec.md"})
	if code != 1 || !strings.Contains(out, "reviews/s.md") {
		t.Fatalf("Command = (%d, %q), want exit 1 naming reviews/s.md", code, out)
	}
}

// citedFileSpec writes a folder spec whose one mapped row cites rel, and writes the
// cited file with header ahead of its package clause. The file always declares the
// cited name, so only the execution arm can red.
func citedFileSpec(t *testing.T, rel, header string) (root, specPath string) {
	t.Helper()
	root, specPath = citedSpec(t, "s", "`"+rel+"` (`TestPresent`)", "")
	writeUnder(t, filepath.Join(root, filepath.FromSlash(rel)), header+"package x\n\nfunc TestPresent(t *testing.T) {}\n")
	return root, specPath
}

// goModuleRoot makes root a Go module. That is what gives the root the built-in kit
// phase table, whose test phases are the untagged set and the system set, so the census
// under grade is the one the gate would run here.
func goModuleRoot(t *testing.T, root string) {
	t.Helper()
	writeUnder(t, filepath.Join(root, "go.mod"), "module example.test\n\ngo 1.21\n")
}

// TestCitationUnexecutedConstraint grades the execution arm of the citation check: a
// cited test is evidence only when some tag set the gate executes compiles its file.
func TestCitationUnexecutedConstraint(t *testing.T) {
	t.Run("an untagged file passes", func(t *testing.T) {
		root, specPath := citedFileSpec(t, "internal/x/plain_test.go", "")
		goModuleRoot(t, root)

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for an unconstrained file", v)
		}
	})

	t.Run("a system-tagged file passes", func(t *testing.T) {
		root, specPath := citedFileSpec(t, "internal/x/sys_test.go", "//go:build system\n\n")
		goModuleRoot(t, root)
		t.Setenv("BENCH_KIT", root) // the system set exists only where the root is its own kit; the gate's ambient kit must not leak in

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for the executed system set", v)
		}
	})

	t.Run("a stress-tagged file is a violation", func(t *testing.T) {
		root, specPath := citedFileSpec(t, "internal/x/stress_test.go", "//go:build stress\n\n")
		goModuleRoot(t, root)

		want := "coverage map row 1 cites 'internal/x/stress_test.go', which no executed tag set builds (//go:build stress)"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("a foreign GOOS filename suffix is a violation", func(t *testing.T) {
		root, specPath := citedFileSpec(t, "internal/x/foo_windows_test.go", "")
		goModuleRoot(t, root)

		want := "coverage map row 1 cites 'internal/x/foo_windows_test.go', which no executed tag set builds (the filename's GOOS or GOARCH suffix)"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("a malformed constraint is a violation naming the file", func(t *testing.T) {
		root, specPath := citedFileSpec(t, "internal/x/bad_test.go", "//go:build system &&\n\n")
		goModuleRoot(t, root)

		want := "coverage map row 1 cites 'internal/x/bad_test.go', whose //go:build expression does not parse: unexpected end of expression"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("a non-regular path is a violation", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/fifo_test.go` (`TestPresent`)", "")
		goModuleRoot(t, root)
		fifo := filepath.Join(root, "internal", "x", "fifo_test.go")
		if err := os.MkdirAll(filepath.Dir(fifo), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := syscall.Mkfifo(fifo, 0o644); err != nil {
			t.Fatalf("mkfifo: %v", err)
		}

		want := "coverage map row 1 cites 'internal/x/fifo_test.go', which is not a regular file"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("a root with no test phase is inapplicable", func(t *testing.T) {
		// A root that is neither a Go module nor the kit itself has no test phase at
		// all: the toolchain phases need the module, and the system phase grades the
		// kit alone.
		t.Setenv("BENCH_KIT", t.TempDir())
		_, specPath := citedFileSpec(t, "internal/x/stress_test.go", "//go:build stress\n\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation where the gate runs no test phase", v)
		}
	})

	t.Run("a manifest-declared custom tag passes", func(t *testing.T) {
		root, specPath := citedFileSpec(t, "internal/x/custom_test.go", "//go:build customsuite\n\n")
		writeUnder(t, filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath)),
			`{"phases":[{"name":"test","argv":["go","test","-tags=customsuite","./..."]}]}`+"\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for a manifest-declared tag", v)
		}
	})
}
