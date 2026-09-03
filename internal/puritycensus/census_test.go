package puritycensus

import (
	"strings"
	"testing"
)

// The census reads this file too, so every synthetic subject line below is assembled
// from parts. A literal would report this file instead of the fixture.

// TestCensusDiagnosesForbiddenImportAmbientEffectAndParallel drives the three checks
// over one in-memory source and proves each diagnostic carries the file and the line.
// (Coverage row LQ13.)
func TestCensusDiagnosesForbiddenImportAmbientEffectAndParallel(t *testing.T) {
	source := strings.Join([]string{
		"\t\"github.com/gibbonmi/bench/" + "internal/git\"",
		"\tvalue := os." + "Getenv(\"HOME\")",
		"\t" + "t." + "Parallel()",
	}, "\n")
	diagnostics := diagnose("fixture.go", source, PolicyPackage())
	if len(diagnostics) != 3 {
		t.Fatalf("the census yields %d diagnostics, want 3: %v", len(diagnostics), diagnostics)
	}
	wants := []string{
		"fixture.go:1: forbidden import",
		"fixture.go:2: ambient effect",
		"fixture.go:3: t.Parallel",
	}
	for i, want := range wants {
		if !strings.HasPrefix(diagnostics[i], want) {
			t.Errorf("diagnostic %d is %q, want the prefix %q", i, diagnostics[i], want)
		}
	}
}

// TestCensusRefusesProcessBackedFixture proves exec.Command stays in the ambient set,
// so a process-backed git fixture inside a pure package reds. (Coverage row LQ15.)
func TestCensusRefusesProcessBackedFixture(t *testing.T) {
	line := "\tcmd := " + "exec" + "." + "Command(\"git\", \"status\")"
	diagnostics := diagnose("fixture_test.go", line, PolicyPackage())
	if len(diagnostics) != 1 || !strings.Contains(diagnostics[0], "exec"+".Command") {
		t.Fatalf("the census accepts a process-backed fixture: %v", diagnostics)
	}
}

// TestCensusExemptsOnlyTheWrapperEffects proves the self-exempt file name suppresses
// the ambient and parallel checks for the wrapper that drives the census, and that the
// forbidden-import check still applies there.
func TestCensusExemptsOnlyTheWrapperEffects(t *testing.T) {
	source := strings.Join([]string{
		"\tvalue := os." + "Getenv(\"HOME\")",
		"\t" + "t." + "Parallel()",
		"\t\"github.com/gibbonmi/bench/" + "internal/git\"",
	}, "\n")
	diagnostics := diagnose(censusFile, source, PolicyPackage())
	if len(diagnostics) != 1 || !strings.Contains(diagnostics[0], "forbidden import") {
		t.Fatalf("the wrapper exemption covers the wrong checks: %v", diagnostics)
	}
}

// TestCensusExemptsItsOwnImportPath proves a wrapper may import this package under the
// leaf policy, which otherwise forbids every path under internal/.
func TestCensusExemptsItsOwnImportPath(t *testing.T) {
	if diagnostics := diagnose("wrapper_test.go", "\t"+helperImport, LeafPackage()); len(diagnostics) != 0 {
		t.Fatalf("the leaf policy refuses the census helper itself: %v", diagnostics)
	}
}

// TestCensusReadsCodeBeforeTheComment proves a mention inside a comment cannot trip a
// check, which keeps the policy documentation in this repository writable.
func TestCensusReadsCodeBeforeTheComment(t *testing.T) {
	line := "// a comment naming os." + "Getenv( and " + "t." + "Parallel("
	if diagnostics := diagnose("doc.go", line, PolicyPackage()); len(diagnostics) != 0 {
		t.Fatalf("a comment trips the census: %v", diagnostics)
	}
}
