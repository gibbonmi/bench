package canonicalpath

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// forbiddenImport names the dependencies the canonical-path owner must never take: any
// Bench package under internal/, process execution, and the system-call surface. Every
// importer of this package sits under internal/, so one internal import here can close a
// cycle or drag an effect into the derivation.
var forbiddenImport = regexp.MustCompile(`"(os/exec|syscall|github\.com/gibbonmi/bench/internal/[^"]*)"`)

// ambientEffect names the ambient process reads, mutations, and descendant starts the
// owner must never perform.
var ambientEffect = regexp.MustCompile(`\b(os\.Getenv|os\.LookupEnv|os\.Setenv|os\.Getwd|os\.Chdir|os\.Environ|time\.Now|exec\.Command|exec\.CommandContext)\(`)

// TestPurePackageSourceCensus enforces the leaf-package boundary: no import of any Bench
// internal package, os/exec, or syscall, and no ambient environment, directory, clock, or
// descendant effect anywhere in the package. Test files obey the same rule, so a fixture
// cannot smuggle an effect in.
func TestPurePackageSourceCensus(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	sources := 0
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		sources++
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			code := line
			if idx := strings.Index(code, "//"); idx >= 0 {
				code = code[:idx]
			}
			if match := forbiddenImport.FindString(code); match != "" {
				t.Errorf("%s:%d: forbidden import %s in the canonical-path owner", name, i+1, match)
			}
			if match := ambientEffect.FindString(code); match != "" && name != "purity_census_test.go" {
				t.Errorf("%s:%d: ambient effect %s in the canonical-path owner", name, i+1, match)
			}
			if name != "purity_census_test.go" && strings.Contains(code, "t.Parallel(") {
				t.Errorf("%s:%d: t.Parallel is out of scope for this spec", name, i+1)
			}
		}
	}
	if sources == 0 {
		t.Fatal("census found no Go sources")
	}
}
