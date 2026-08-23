package lifecyclepolicy

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// forbiddenImport names the dependencies a pure policy owner must never take:
// the parent effect adapter, the Git owner, process execution, and the bounds
// constants the parent boundary supplies as explicit threshold facts.
var forbiddenImport = regexp.MustCompile(`"(os/exec|syscall|github\.com/gibbonmi/bench/internal/git|github\.com/gibbonmi/bench/internal/bounds|github\.com/gibbonmi/bench/internal/worktree)"`)

// ambientEffect names the ambient process reads, mutations, and descendant
// starts a pure decision owner must never perform.
var ambientEffect = regexp.MustCompile(`\b(os\.Getenv|os\.LookupEnv|os\.Setenv|os\.Getwd|os\.Chdir|os\.Environ|time\.Now|exec\.Command|exec\.CommandContext)\(`)

// TestPurePackageSourceCensus enforces the pure-package boundary: no import of
// the parent package, internal/git, os/exec, or syscall, and no ambient
// environment, directory, clock, or descendant effect anywhere in the package.
// Test files obey the same rule, so a fixture cannot smuggle an effect in.
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
				t.Errorf("%s:%d: forbidden import %s in the pure policy package", name, i+1, match)
			}
			if match := ambientEffect.FindString(code); match != "" && name != "purity_census_test.go" {
				t.Errorf("%s:%d: ambient effect %s in the pure policy package", name, i+1, match)
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
