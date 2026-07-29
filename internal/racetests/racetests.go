package racetests

import (
	"path"
	"strings"
)

// Test is one regression test the race phase must execute.
type Test struct {
	PackagePath string
	Name        string
}

// Tests names the regression tests that need race instrumentation because the ordinary
// suite cannot observe their concurrent failure modes.
var Tests = []Test{
	{PackagePath: "./internal/worktree", Name: "TestConcurrentCleanupRecordsOneTransaction"},
	{PackagePath: "./internal/guards", Name: "TestScanTimeoutPreservesPartialRowsAndHonestCounts"},
	{PackagePath: "./internal/guards", Name: "TestScanEnumerationTimeoutUsesUnknownCounts"},
}

// SyntheticSources supplies grouped test declarations for a fixture that needs every
// registered race test without carrying a second hand-maintained list of their names.
func SyntheticSources() map[string]string {
	return SyntheticSourcesFor(Tests)
}

// SyntheticSourcesFor supplies grouped test declarations for the registered subset a
// fixture needs.
func SyntheticSourcesFor(tests []Test) map[string]string {
	byPackage := map[string][]string{}
	for _, test := range tests {
		byPackage[test.PackagePath] = append(byPackage[test.PackagePath], test.Name)
	}
	sources := make(map[string]string, len(byPackage))
	for packagePath, names := range byPackage {
		var source strings.Builder
		source.WriteString("package ")
		source.WriteString(path.Base(packagePath))
		source.WriteString("\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"testing\"\n)\n")
		for _, name := range names {
			source.WriteString("\nfunc ")
			source.WriteString(name)
			source.WriteString("(t *testing.T) { fmt.Fprintln(os.Stderr, \"race test noise\") }\n")
		}
		sources[packagePath+"/race_registry_test.go"] = source.String()
	}
	return sources
}
