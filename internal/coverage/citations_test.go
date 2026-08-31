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
	"github.com/gibbonmi/bench/internal/gate"
)

// citedSpec writes a folder spec at <root>/specs/<slug>/spec.md whose single mapped
// row carries seam as its seam cell, and returns that spec's path. fences, when
// non-empty, is appended verbatim so a test can shape the section it grades.
func citedSpec(t *testing.T, slug, seam, fences string) (root, specPath string) {
	t.Helper()
	return citedCellSpec(t, slug, seam, "w", fences)
}

// citedCellSpec is citedSpec with the why cell under the test's control, so a test can
// prove which cell the grammar reads.
func citedCellSpec(t *testing.T, slug, seam, why, fences string) (root, specPath string) {
	t.Helper()
	// A synthetic fixture root is not the kit. Without this pin the fixture becomes its
	// own kit whenever the ambient one is unset, which materializes the system phase and
	// its ./internal/systemtest operand in a tree that holds no such package. The tree
	// itself then decides the census: a go.mod gets the toolchain test phase, and a root
	// without one gets no test phase at all.
	t.Setenv("BENCH_KIT", t.TempDir())
	root = t.TempDir()
	body := "# " + slug + "\n\n" + stories + "\n### Acceptance coverage map\n" + hdrReduced +
		"| 1 | b | " + seam + " | " + why + " |\n" + fences
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
		"package x\n\nfunc TestPresent(t *testing.T) {}\n\nfunc TestOther(t *testing.T) {\n\tt.Run(\"a case\", func(t *testing.T) {})\n}\n")

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
	cited := filepath.Join(root, filepath.FromSlash(rel))
	writeUnder(t, cited, header+"package x\n\nfunc TestPresent(t *testing.T) {}\n")
	// The cited package also holds one file no constraint excludes. Go drops a directory
	// whose every file is excluded out of a recursive pattern, so a package holding the
	// cited file alone is never selected, and the tag arm this helper grades is never
	// reached.
	writeUnder(t, filepath.Join(filepath.Dir(cited), "compiled.go"), "package x\n")
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
		// The system phase carries one package operand, ./internal/systemtest, so that
		// package is the only place a system-tagged file is executed evidence.
		root, specPath := citedFileSpec(t, "internal/systemtest/sys_test.go", "//go:build system\n\n")
		goModuleRoot(t, root)
		t.Setenv("BENCH_KIT", root) // the system phase materializes only where the root is its own kit

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

	t.Run("a symlinked path is a violation", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/link_test.go` (`TestPresent`)", "")
		goModuleRoot(t, root)
		target := filepath.Join(root, "internal", "x", "real_test.go")
		writeUnder(t, target, "package x\n\nfunc TestPresent(t *testing.T) {}\n")
		if err := os.Symlink(target, filepath.Join(root, "internal", "x", "link_test.go")); err != nil {
			t.Fatalf("Symlink: %v", err)
		}

		want := "coverage map row 1 cites 'internal/x/link_test.go', which is not a regular file"
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
		goModuleRoot(t, root)
		// A prospective landing gate points the phase schedule at its baseline, which
		// would answer this root with the kit's table instead of the manifest below.
		t.Setenv(gate.BaselinePolicyEnv, "")
		writeUnder(t, filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath)),
			`{"phases":[{"name":"test","argv":["go","test","-tags=customsuite","./..."]}]}`+"\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for a manifest-declared tag", v)
		}
	})
}

// TestCitationOutsideTheRepoIsRefused pins the containment rule: a cited path that
// escapes the tree the spec resolves against is a violation, not a resolved file.
func TestCitationOutsideTheRepoIsRefused(t *testing.T) {
	_, specPath := citedSpec(t, "s", "`../outside/foo_test.go` (`TestPresent`)", "")

	want := "coverage map row 1 cites '../outside/foo_test.go', which is not a path inside the repo"
	if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
		t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
	}
}

// mentionSpecFile writes the cited file every mention and subtest case grades. body is
// the file's declarations after its package clause.
func mentionSpecFile(t *testing.T, root, body string) {
	t.Helper()
	writeUnder(t, filepath.Join(root, "internal", "x", "foo_test.go"), "package x\n\n"+body)
}

// TestMentionIsNotACitation grades the mention rule: a seam-cell test path with no name
// list claims evidence it never names, and the why cell stays outside the grammar.
func TestMentionIsNotACitation(t *testing.T) {
	t.Run("a seam-cell path with no name list is a violation", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go`", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {}\n")

		want := "coverage map row 1 mentions 'internal/x/foo_test.go' without a cited name list"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("an empty name list stays a non-citation", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` ()", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for an empty list", v)
		}
	})

	t.Run("a why-cell path is not graded", func(t *testing.T) {
		root, specPath := citedCellSpec(t, "s", "review-owned: the Standards axis reads the type",
			"`internal/x/foo_test.go` names the shape", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for a why-cell mention", v)
		}
	})
}

// TestSubtestSegmentResolves grades the subtest rule: a cited segment must appear as a
// t.Run string literal, and a file with a computed t.Run name is exempt.
func TestSubtestSegmentResolves(t *testing.T) {
	t.Run("a declared segment passes", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent/a case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tt.Run(\"a case\", func(t *testing.T) {})\n}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for a declared segment", v)
		}
	})

	t.Run("an absent segment is a violation", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent/renamed case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tt.Run(\"a case\", func(t *testing.T) {})\n}\n")

		want := "coverage map row 1 cites 'TestPresent/renamed case', whose subtest 'renamed case' " +
			"is no t.Run name in 'internal/x/foo_test.go'"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("a non-literal t.Run name exempts the file", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent/renamed case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tfor _, c := range cases {\n\t\tt.Run(c.name, func(t *testing.T) {})\n\t}\n}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation where a t.Run name is computed", v)
		}
	})

	t.Run("a foreign Run call neither names nor exempts", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent/renamed case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tcmd.Run()\n\tt.Run(\"a case\", func(t *testing.T) {})\n}\n")

		want := "coverage map row 1 cites 'TestPresent/renamed case', whose subtest 'renamed case' " +
			"is no t.Run name in 'internal/x/foo_test.go'"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})

	t.Run("a raw-string t.Run name resolves", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestPresent/a case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tt.Run(`a case`, func(t *testing.T) {})\n}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation for a raw-string subtest name", v)
		}
	})

	t.Run("the parent function still resolves in an exempt file", func(t *testing.T) {
		root, specPath := citedSpec(t, "s", "`internal/x/foo_test.go` (`TestAbsent/any case`)", "")
		mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tt.Run(c.name, func(t *testing.T) {})\n}\n")

		want := "coverage map row 1 cites 'TestAbsent/any case', which 'internal/x/foo_test.go' does not declare"
		if v := checkFilesOf(t, specPath); len(v) != 1 || v[0] != want {
			t.Fatalf("CheckFiles = %#v, want exactly %q", v, want)
		}
	})
}

// citationClassesSpec writes the three-row map both class tests grade: a mention, a
// stale subtest segment, and a citation into a never-executed file. header is the spec
// preamble that decides whether the map is historical.
func citationClassesSpec(t *testing.T, header string) string {
	t.Helper()
	root := t.TempDir()
	body := "# s\n\n" + header + stories + "\n### Acceptance coverage map\n" + hdrReducedID +
		"| MT1 | 1 | b | `internal/x/foo_test.go` | w |\n" +
		"| XT2 | 1 | c | `internal/x/foo_test.go` (`TestPresent/renamed case`) | w |\n" +
		"| MT3 | 1 | d | `internal/x/stress_test.go` (`TestStress`) | w |\n"
	specPath := filepath.Join(root, "specs", "s", "spec.md")
	writeUnder(t, specPath, body)
	mentionSpecFile(t, root, "func TestPresent(t *testing.T) {\n\tt.Run(\"a case\", func(t *testing.T) {})\n}\n")
	writeUnder(t, filepath.Join(root, "internal", "x", "stress_test.go"),
		"//go:build stress\n\npackage x\n\nfunc TestStress(t *testing.T) {}\n")
	goModuleRoot(t, root)
	t.Setenv("BENCH_KIT", t.TempDir()) // the fixture is not its own kit; only its go.mod decides the census
	return specPath
}

// TestHistoricalSpecSilencesCitationChecks pins the opt-out over every citation check: a
// historical spec carries a mixed-tag row-ID set, a mention, a stale subtest, and a
// citation into a never-executed file, and stays silent. A partial opt-out would break
// the documented contract.
func TestHistoricalSpecSilencesCitationChecks(t *testing.T) {
	specPath := citationClassesSpec(t, historicalMarker+"\n\n")

	if v := checkFilesOf(t, specPath); len(v) != 0 {
		t.Fatalf("CheckFiles = %#v, want no violation for a historical spec", v)
	}
}

// TestParseSpecReturnsCitationViolationClasses pins the shared entry point the review
// preflight reads: the mixed-tag, mention, subtest, and unexecuted-constraint classes
// all reach a caller outside this package, so the preflight and the gate agree.
func TestParseSpecReturnsCitationViolationClasses(t *testing.T) {
	specPath := citationClassesSpec(t, "")

	optIn, ids, violations, err := ParseSpec(specPath)
	if err != nil {
		t.Fatalf("ParseSpec(%q): %v", specPath, err)
	}
	if !optIn || len(ids) != 3 {
		t.Fatalf("ParseSpec = (optIn %v, ids %#v), want an opted-in map of three rows", optIn, ids)
	}
	for _, want := range []string{
		"coverage map row ids carry more than one tag (MT, XT); a map declares one tag",
		"coverage map row 1 mentions 'internal/x/foo_test.go' without a cited name list",
		"coverage map row 2 cites 'TestPresent/renamed case', whose subtest 'renamed case' is no t.Run name in 'internal/x/foo_test.go'",
		"coverage map row 3 cites 'internal/x/stress_test.go', which no executed tag set builds (//go:build stress)",
	} {
		if !hasViolation(violations, want) {
			t.Fatalf("ParseSpec violations = %#v, want one holding %q", violations, want)
		}
	}
}
