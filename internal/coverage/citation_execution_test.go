package coverage

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/gate"
)

// TestCitationPackageScope grades the package arm of the citation check: a cited file is
// evidence only where one complete test phase both selects its package and builds it.
func TestCitationPackageScope(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		rel      string
		source   string
		want     string
	}{
		{
			name: "package scope and tags from different phases do not combine",
			manifest: `{"phases":[
				{"name":"scope","argv":["go","test","-tags=scope","./match"]},
				{"name":"match","argv":["go","test","-tags=match","./other"]}
			]}`,
			rel:    "match/match_test.go",
			source: "//go:build match\n\npackage match\n\nfunc TestPresent() {}\n",
			want:   "which no executed tag set builds",
		},
		{
			name:     "recursive pattern excludes testdata",
			manifest: `{"phases":[{"name":"test","argv":["go","test","./..."]}]}`,
			rel:      "testdata/fixture/fixture_test.go",
			source:   "package fixture\n\nfunc TestPresent() {}\n",
			want:     "which no executed test phase selects",
		},
		{
			name:     "recursive pattern excludes underscore-prefixed packages",
			manifest: `{"phases":[{"name":"test","argv":["go","test","./..."]}]}`,
			rel:      "_private/fixture/fixture_test.go",
			source:   "package fixture\n\nfunc TestPresent() {}\n",
			want:     "which no executed test phase selects",
		},
		{
			name:     "explicit package selects evidence under testdata",
			manifest: `{"phases":[{"name":"test","argv":["go","test","./testdata/fixture"]}]}`,
			rel:      "testdata/fixture/fixture_test.go",
			source:   "package fixture\n\nfunc TestPresent() {}\n",
		},
		{
			name:     "absent package operands select the effective phase directory",
			manifest: `{"phases":[{"name":"test","argv":["go","test"]}]}`,
			rel:      "fixture_test.go",
			source:   "package fixture\n\nfunc TestPresent() {}\n",
		},
		{
			// The first phase selects the cited package but refuses its build constraint.
			// The second phase does both, which is what the citation needs. A check that
			// stopped at the first selecting phase, or that required every phase to
			// agree, would red here.
			name: "one phase that both selects and builds accepts the citation",
			manifest: `{"phases":[
				{"name":"plain","argv":["go","test","./match"]},
				{"name":"tagged","argv":["go","test","-tags=match","./match"]}
			]}`,
			rel:    "match/match_test.go",
			source: "//go:build match\n\npackage match\n\nfunc TestPresent() {}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			specPath := citationExecutionSpec(t, test.manifest, test.rel, test.source)
			v := checkFilesOf(t, specPath)
			if test.want == "" {
				if len(v) != 0 {
					t.Fatalf("CheckFiles = %#v, want no violation", v)
				}
				return
			}
			if len(v) != 1 || !strings.Contains(v[0], "coverage map row 1") ||
				!strings.Contains(v[0], test.rel) || !strings.Contains(v[0], test.want) {
				t.Fatalf("CheckFiles = %#v, want row, file, and %q", v, test.want)
			}
		})
	}
}

// TestCitationPackageScopeFailure grades the failure arm: a phase whose package operands
// the Go loader cannot expand rejects the citation and names the row and the file. A
// pattern the loader resolves to nothing is not that failure; it selects nothing.
func TestCitationPackageScopeFailure(t *testing.T) {
	const source = "package fixture\n\nfunc TestPresent() {}\n"

	t.Run("an unresolvable package operand rejects the citation", func(t *testing.T) {
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./missing"]}]}`,
			"fixture_test.go", source)

		assertExpansionRefusal(t, checkFilesOf(t, specPath), "fixture_test.go")
	})

	t.Run("an absent Go executable rejects the citation", func(t *testing.T) {
		// The phase manifest resolves the phase table without the toolchain, so the only
		// step this PATH deprives is the package loader itself.
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."]}]}`,
			"fixture_test.go", source)
		t.Setenv("PATH", t.TempDir())

		assertExpansionRefusal(t, checkFilesOf(t, specPath), "fixture_test.go")
	})

	t.Run("a pattern that matches no package is no expansion failure", func(t *testing.T) {
		// Go excludes testdata from a recursive pattern, so this phase resolves to no
		// package at all. The loader says so on its error stream and exits zero, which
		// leaves the citation unselected rather than unexpandable.
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./testdata/..."]}]}`,
			"testdata/fixture/fixture_test.go", source)

		v := checkFilesOf(t, specPath)
		if len(v) != 1 || !strings.Contains(v[0], "which no executed test phase selects") {
			t.Fatalf("CheckFiles = %#v, want the selection violation rather than an expansion failure", v)
		}
	})
}

// TestCitationPhaseEnv grades the phase-environment arm. A phase declares its own
// platform, and `go test` would select and compile different files under it. So the check
// reads the phase's environment rather than the host's on both sides: the package loader,
// and the file's build constraint.
func TestCitationPhaseEnv(t *testing.T) {
	arch := crossArch()
	// One fixture source serves every cgo case, so the cases differ only in the phase
	// environment under grade.
	const cgoSource = "//go:build cgo\n\npackage linked\n\nfunc TestPresent() {}\n"

	t.Run("the loader child runs under the phase environment", func(t *testing.T) {
		root := t.TempDir()
		writeUnder(t, filepath.Join(root, "go.mod"), "module example.test/phaseenv\n\ngo 1.24\n")
		// This package holds one file only the cross architecture builds. A recursive
		// pattern drops the directory on every other architecture, so the loader lists it
		// only when the loader itself runs under the phase's own GOARCH.
		writeUnder(t, filepath.Join(root, "arch", "arch_"+arch+".go"), "package arch\n")

		dirs, err := selectedPackageDirs(gate.TestExecution{
			Name: "test", Packages: []string{"./..."}, Dir: root, Env: []string{"GOARCH=" + arch},
		})
		if err != nil {
			t.Fatalf("selectedPackageDirs: %v", err)
		}
		if !containsDirectory(dirs, filepath.Join(root, "arch")) {
			t.Fatalf("selectedPackageDirs = %q, want the package the phase's own GOARCH selects", dirs)
		}
	})

	t.Run("a phase GOARCH override builds an arch-suffixed file", func(t *testing.T) {
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"GOARCH":"`+arch+`"}}]}`,
			"arch/fixture_"+arch+"_test.go",
			"package arch\n\nfunc TestPresent() {}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation under the phase's own GOARCH", v)
		}
	})

	// The two cgo cases are symmetric on purpose. One of them contradicts the host's own
	// CgoEnabled default whichever way that default falls, so the pair bites on any host.
	t.Run("a phase that disables cgo refuses a cgo-tagged file", func(t *testing.T) {
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"CGO_ENABLED":"0"}}]}`,
			"linked/linked_test.go", cgoSource)

		v := checkFilesOf(t, specPath)
		if len(v) != 1 || !strings.Contains(v[0], "which no executed tag set builds (//go:build cgo)") {
			t.Fatalf("CheckFiles = %#v, want the cgo constraint refused", v)
		}
	})

	t.Run("a phase that enables cgo builds a cgo-tagged file", func(t *testing.T) {
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"CGO_ENABLED":"1"}}]}`,
			"linked/linked_test.go", cgoSource)

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation when the phase enables cgo", v)
		}
	})

	t.Run("a cross-architecture phase refuses a cgo-tagged file", func(t *testing.T) {
		// Go turns cgo off for a cross-compiled phase unless that phase names CGO_ENABLED
		// itself, whatever the host's own default is. So an architecture override alone
		// must reach the file constraint.
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"GOARCH":"`+arch+`"}}]}`,
			"linked/linked_test.go", cgoSource)

		v := checkFilesOf(t, specPath)
		if len(v) != 1 || !strings.Contains(v[0], "which no executed tag set builds (//go:build cgo)") {
			t.Fatalf("CheckFiles = %#v, want the phase architecture to refuse cgo", v)
		}
	})

	t.Run("an empty cgo value reads as an unset one", func(t *testing.T) {
		// Go reads an empty value as an unset variable and keeps its own default. So this
		// phase must reach the same verdict as a phase that declares no environment.
		empty := checkFilesOf(t, citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"CGO_ENABLED":""}}]}`,
			"linked/linked_test.go", cgoSource))
		absent := checkFilesOf(t, citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."]}]}`,
			"linked/linked_test.go", cgoSource))
		if !reflect.DeepEqual(empty, absent) {
			t.Fatalf("empty CGO_ENABLED = %#v, want the absent verdict %#v", empty, absent)
		}
	})

	t.Run("an empty platform value keeps the host platform", func(t *testing.T) {
		// An empty GOOS declares no platform, so the cited file's own GOOS suffix still
		// has to match the platform the phase runs on.
		specPath := citationExecutionSpec(t,
			`{"phases":[{"name":"test","argv":["go","test","./..."],"env":{"GOOS":""}}]}`,
			"host/fixture_"+runtime.GOOS+"_test.go", "package host\n\nfunc TestPresent() {}\n")

		if v := checkFilesOf(t, specPath); len(v) != 0 {
			t.Fatalf("CheckFiles = %#v, want no violation under an empty GOOS", v)
		}
	})
}

// crossArch names an architecture the host is not, so a phase's own override is
// observable rather than a restatement of the host's platform.
func crossArch() string {
	if runtime.GOARCH == "386" {
		return "amd64"
	}
	return "386"
}

// TestCitationPackageScopeTimeout grades the loader's time bound. A cited FIFO never
// starts the loader, because the grammar refuses a non-regular file first, but a FIFO
// planted anywhere else under a phase's package tree still blocks the loader in open(2)
// on another row's honest citation. The bound must end that wait and report it as the
// expansion failure it is, rather than hold the check open with no deadline.
func TestCitationPackageScopeTimeout(t *testing.T) {
	specPath := citationExecutionSpec(t,
		`{"phases":[{"name":"test","argv":["go","test","./..."]}]}`,
		"fixture_test.go", "package fixture\n\nfunc TestPresent() {}\n")
	// The stub sleeps far past the shrunk bound, and it sleeps in a child of the shell
	// the kernel actually executes. Only a process-group kill reaches that child, so a
	// bound that killed the shell alone would still wait on the inherited output pipe.
	stub := t.TempDir()
	writeExecutable(t, filepath.Join(stub, "go"), "#!/bin/sh\nsleep 600 &\nwait\n")
	// The stub shadows the real toolchain rather than replacing PATH, because the script
	// itself needs the commands the rest of PATH provides.
	t.Setenv("PATH", stub+string(os.PathListSeparator)+os.Getenv("PATH"))
	shrinkPackageLoadTimeout(t, 50*time.Millisecond)

	started := time.Now()
	v := checkFilesOf(t, specPath)
	elapsed := time.Since(started)

	assertExpansionRefusal(t, v, "fixture_test.go")
	if !strings.Contains(v[0], `executed phase "test"`) || !strings.Contains(v[0], "timed out") {
		t.Fatalf("CheckFiles = %#v, want the phase name and a timeout", v)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("CheckFiles took %s, want the loader bound to end the wait promptly", elapsed)
	}
}

// shrinkPackageLoadTimeout installs a test-only loader bound and restores the production
// one, so a hang test costs milliseconds rather than the real bound.
func shrinkPackageLoadTimeout(t *testing.T, limit time.Duration) {
	t.Helper()
	previous := packageLoadTimeout
	packageLoadTimeout = limit
	t.Cleanup(func() { packageLoadTimeout = previous })
}

// writeExecutable writes an executable file at an absolute path, creating its parents.
func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	writeUnder(t, path, content)
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("Chmod(%q): %v", path, err)
	}
}

func assertExpansionRefusal(t *testing.T, v []string, rel string) {
	t.Helper()
	if len(v) != 1 || !strings.Contains(v[0], "coverage map row 1") ||
		!strings.Contains(v[0], rel) || !strings.Contains(v[0], "could not expand packages") {
		t.Fatalf("CheckFiles = %#v, want an expansion violation naming the row and file", v)
	}
}

// citationExecutionSpec writes one temporary Go module, the phase manifest that declares
// its test phases, and the folder spec whose single mapped row cites rel. The module
// root holds a space, because the loader prints a filesystem path and the check must
// keep such a path whole.
func citationExecutionSpec(t *testing.T, manifest, rel, source string) string {
	t.Helper()
	// A prospective landing gate points the phase schedule at its baseline, which would
	// answer this root with the kit's table instead of the manifest written here.
	t.Setenv(gate.BaselinePolicyEnv, "")
	root := filepath.Join(t.TempDir(), "with space")
	writeUnder(t, filepath.Join(root, "go.mod"), "module example.test/citation\n\ngo 1.24\n")
	writeUnder(t, filepath.Join(root, filepath.FromSlash(canary.PhaseManifestPath)), manifest+"\n")
	writeUnder(t, filepath.Join(root, "other", "other.go"), "package other\n")
	cited := filepath.Join(root, filepath.FromSlash(rel))
	writeUnder(t, cited, source)
	writeCompiledCompanion(t, cited, source)
	return citedSpecPath(t, root, rel)
}

// writeCompiledCompanion gives a cited fixture file one companion that no constraint
// excludes. Go refuses an explicit operand for a directory whose every file a constraint
// excludes, and drops such a directory out of a recursive pattern. So a package holding
// the cited file alone could never be selected at all, and the arm under grade would
// never be reached.
func writeCompiledCompanion(t *testing.T, cited, source string) {
	t.Helper()
	writeUnder(t, filepath.Join(filepath.Dir(cited), "compiled.go"), "package "+packageOf(t, source)+"\n")
}

// packageOf reads the package clause a fixture source declares, so the companion file
// joins that package rather than making the directory hold two.
func packageOf(t *testing.T, source string) string {
	t.Helper()
	for _, line := range strings.Split(source, "\n") {
		if name, found := strings.CutPrefix(line, "package "); found {
			return strings.TrimSpace(name)
		}
	}
	t.Fatalf("fixture source declares no package clause: %q", source)
	return ""
}

func citedSpecPath(t *testing.T, root, rel string) string {
	t.Helper()
	path := filepath.Join(root, "specs", "citation", "spec.md")
	writeUnder(t, path, "# citation fixture\n\n## User stories\n\n1. A story.\n\n"+
		"### Acceptance coverage map\n\n| story | behavior | seam | why it catches the failure |\n"+
		"|---|---|---|---|\n| 1 | a behavior | `"+rel+"` (`TestPresent`) | why |\n")
	return path
}
