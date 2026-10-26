package gate

// The reader-table derivation over canned test logs. Each fixture holds a root-conformance
// log that reads every covered family, so the witness holds unless a test removes a read.

import (
	"slices"
	"strings"
	"testing"
)

const (
	// tableRoot is the graded root of the canned logs.
	tableRoot = "/kit"
	// tableReader is the test package whose log each test varies.
	tableReader = "internal/reader"
)

// witnessLog is a root-conformance log that reads one path below root of each family except
// skipped.
func witnessLog(root, skipped string) string {
	log := testLogHeader + "\n"
	for _, family := range coveredFamilies {
		if family == skipped {
			continue
		}
		read := family
		if strings.HasSuffix(family, "/") {
			read += "witness.md"
		}
		log += "open " + root + "/" + read + "\n"
	}
	return log
}

// tableSources are the sources of the witness log and of tableReader with log.
func tableSources(log string) readerSources {
	return readerSources{
		Root:     tableRoot,
		Packages: []string{conformancePackagePath, tableReader},
		Logs:     map[string][]byte{conformancePackagePath: []byte(witnessLog(tableRoot, "")), tableReader: []byte(log)},
		Families: coveredFamilies,
	}
}

// deriveTable derives the table of in or fails the test.
func deriveTable(t *testing.T, in readerSources) readerTable {
	t.Helper()
	table, err := deriveReaderTable(in)
	if err != nil {
		t.Fatalf("derivation = %v; want a table", err)
	}
	return table
}

// readFamilies answers the families of table that pkg reads.
func readFamilies(table readerTable, pkg string) []string {
	var families []string
	for family, readers := range table {
		if slices.Contains(readers, pkg) {
			families = append(families, family)
		}
	}
	slices.Sort(families)
	return families
}

// TestReaderTableMapsAReadToItsFamily requires that a read of roadmap/FT1.md makes its
// package a reader of roadmap/.
func TestReaderTableMapsAReadToItsFamily(t *testing.T) {
	table := deriveTable(t, tableSources(testLogHeader+"\nopen "+tableRoot+"/roadmap/FT1.md\n"))
	if got := readFamilies(table, tableReader); !slices.Equal(got, []string{"roadmap/"}) {
		t.Fatalf("families of %s = %q; want roadmap/", tableReader, got)
	}
}

// TestReaderTableCoversANewPathInAKnownFamily requires that the table derived from a read of
// roadmap/FT1.md names that package for a change of roadmap/FT2.md.
func TestReaderTableCoversANewPathInAKnownFamily(t *testing.T) {
	table := deriveTable(t, tableSources(testLogHeader+"\nopen "+tableRoot+"/roadmap/FT1.md\n"))
	family, covered := familyOf(coveredFamilies, "roadmap/FT2.md")
	if !covered || !slices.Contains(table[family], tableReader) {
		t.Fatalf("readers of roadmap/FT2.md = %q in family %q; want %s", table[family], family, tableReader)
	}
}

// TestReaderTableResolvesATrackedSymlink requires that a read below a tracked link counts as
// a read of the link target's family.
func TestReaderTableResolvesATrackedSymlink(t *testing.T) {
	in := tableSources(testLogHeader + "\nopen " + tableRoot + "/.claude/commands/x.md\n")
	in.Links = map[string]string{".claude/commands": "../.agents/commands"}
	if got := readFamilies(deriveTable(t, in), tableReader); !slices.Equal(got, []string{".agents/"}) {
		t.Fatalf("families of %s = %q; want .agents/", tableReader, got)
	}
}

// TestReaderTableSkipsTheOwnDirectory requires that a read inside the package's own
// directory adds no family, even where that directory lies inside a family.
func TestReaderTableSkipsTheOwnDirectory(t *testing.T) {
	const own = "docs/example"
	in := tableSources("")
	in.Packages = []string{conformancePackagePath, own}
	in.Logs[own] = []byte(testLogHeader + "\nopen " + tableRoot + "/" + own + "/testdata/a.md\nstat testdata\n")
	if got := readFamilies(deriveTable(t, in), own); len(got) != 0 {
		t.Fatalf("families of %s = %q; want none", own, got)
	}
}

// TestReaderTableFollowsChdir requires that a relative path after a chdir line resolves
// against the new directory.
func TestReaderTableFollowsChdir(t *testing.T) {
	table := deriveTable(t, tableSources(testLogHeader+"\nchdir ../../specs\nopen a/spec.md\n"))
	if got := readFamilies(table, tableReader); !slices.Equal(got, []string{"specs/"}) {
		t.Fatalf("families of %s = %q; want specs/", tableReader, got)
	}
}

// TestReaderTableIgnoresPathsOutsideTheRoot requires that a path outside the graded root, a
// chdir line, and a getenv line each add no family.
func TestReaderTableIgnoresPathsOutsideTheRoot(t *testing.T) {
	log := testLogHeader + "\nopen /elsewhere" + tableRoot + "/docs/a.md\nstat " + tableRoot + "-other/docs/b.md\nchdir " + tableRoot + "/docs\ngetenv docs/c.md\n"
	if got := readFamilies(deriveTable(t, tableSources(log)), tableReader); len(got) != 0 {
		t.Fatalf("families of %s = %q; want none", tableReader, got)
	}
}

// requireReaderRefusal requires that the derivation of in refuses and names pkg.
func requireReaderRefusal(t *testing.T, in readerSources, pkg string) {
	t.Helper()
	table, err := deriveReaderTable(in)
	if err == nil || table != nil || !strings.Contains(err.Error(), pkg) {
		t.Fatalf("derivation = %v, %v; want a refusal that names %s", table, err, pkg)
	}
}

// TestReaderTableRefusesAMalformedLog requires that a log whose first line is not the header
// refuses the table and names its package.
func TestReaderTableRefusesAMalformedLog(t *testing.T) {
	requireReaderRefusal(t, tableSources("# not a test log\nopen "+tableRoot+"/docs/a.md\n"), tableReader)
}

// TestReaderTableRefusesAnUnknownOperation requires that a line of an unknown operation
// refuses the table.
func TestReaderTableRefusesAnUnknownOperation(t *testing.T) {
	requireReaderRefusal(t, tableSources(testLogHeader+"\nwrite "+tableRoot+"/docs/a.md\n"), tableReader)
}

// TestReaderTableRefusesAnEmptyLog requires that an empty log refuses the table.
func TestReaderTableRefusesAnEmptyLog(t *testing.T) {
	requireReaderRefusal(t, tableSources(""), tableReader)
}

// TestReaderTableRefusesAMissingLog requires that three listed test packages with two logs
// refuse the table and name the package without a log.
func TestReaderTableRefusesAMissingLog(t *testing.T) {
	const unlogged = "internal/unlogged"
	in := tableSources(testLogHeader + "\n")
	in.Packages = append(in.Packages, unlogged)
	requireReaderRefusal(t, in, unlogged)
}

// TestReaderTableRequiresTheConformanceWitness requires that a table in which the
// root-conformance package reads no docs/ path is refused.
func TestReaderTableRequiresTheConformanceWitness(t *testing.T) {
	in := tableSources(testLogHeader + "\nopen " + tableRoot + "/docs/a.md\n")
	in.Logs[conformancePackagePath] = []byte(witnessLog(tableRoot, "docs/"))
	requireReaderRefusal(t, in, "docs/")
}

// TestReaderTableJoinsTheChildReaders requires that the declared child row makes
// internal/releaseevidence a reader of CHANGELOG.md with no logged read.
func TestReaderTableJoinsTheChildReaders(t *testing.T) {
	in := tableSources(testLogHeader + "\n")
	in.Children = childReaders
	if got := readFamilies(deriveTable(t, in), "internal/releaseevidence"); !slices.Equal(got, []string{"CHANGELOG.md"}) {
		t.Fatalf("families of internal/releaseevidence = %q; want CHANGELOG.md", got)
	}
}
