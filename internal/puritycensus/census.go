// Package puritycensus owns the source census a pure package runs over its own
// directory. The forbidden-import patterns, the ambient-effect set, the t.Parallel
// ban, and the self-exempt file name live here once, so the policy owners and the
// leaf owners grade against one policy instead of four copies. Only tests import
// this package.
package puritycensus

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// helperImport is this package's own import path. Every wrapper takes it, and this
// package is itself a pure leaf, so the census never reports it. The raw string is
// the exact text the leaf pattern matches, so this declaration exempts itself.
const helperImport = `"github.com/gibbonmi/bench/internal/puritycensus"`

// censusFile is the wrapper file name every scanned package uses. The wrapper drives
// the census, so the ambient and parallel checks exempt it.
const censusFile = "purity_census_test.go"

// parallelCall is the banned call token. It is assembled from parts because the
// census reads this file too, and a literal would report this line.
var parallelCall = "t." + "Parallel("

// policyImport names the dependencies a pure policy owner must never take: the parent
// effect adapter, the Git owner, process execution, the system-call surface, and the
// bounds constants the parent boundary supplies as explicit threshold facts.
var policyImport = regexp.MustCompile(`"(os/exec|syscall|github\.com/gibbonmi/bench/internal/git|github\.com/gibbonmi/bench/internal/bounds|github\.com/gibbonmi/bench/internal/worktree)"`)

// leafImport names the dependencies a leaf owner must never take: any Bench package
// under internal/, process execution, and the system-call surface. Every importer of a
// leaf sits under internal/, so one internal import there can close a cycle or drag an
// effect into the derivation.
var leafImport = regexp.MustCompile(`"(os/exec|syscall|github\.com/gibbonmi/bench/internal/[^"]*)"`)

// ambientEffect names the ambient process reads, mutations, and descendant starts a
// pure owner must never perform. exec.Command and exec.CommandContext stay in the set,
// so a process-backed fixture counts as an effect.
var ambientEffect = regexp.MustCompile(`\b(os\.Getenv|os\.LookupEnv|os\.Setenv|os\.Getwd|os\.Chdir|os\.Environ|time\.Now|exec\.Command|exec\.CommandContext)\(`)

// Policy carries the four facts a census grades against. The fields stay unexported,
// so a caller takes a policy from PolicyPackage or LeafPackage and never assembles one.
type Policy struct {
	forbiddenImport *regexp.Regexp
	ambientEffect   *regexp.Regexp
	subject         string
	selfExempt      string
}

// PolicyPackage returns the policy a pure worktree decision owner runs. Its parent
// package owns every effect, and the parent boundary hands the owner its threshold
// facts, so the owner takes neither.
func PolicyPackage() Policy {
	return Policy{
		forbiddenImport: policyImport,
		ambientEffect:   ambientEffect,
		subject:         "the pure policy package",
		selfExempt:      censusFile,
	}
}

// LeafPackage returns the policy a leaf owner runs. A leaf imports nothing under
// internal/, so no importer of it can close a cycle.
func LeafPackage() Policy {
	return Policy{
		forbiddenImport: leafImport,
		ambientEffect:   ambientEffect,
		subject:         "the leaf package",
		selfExempt:      censusFile,
	}
}

// Sources names the Go files one census run read, in directory order.
type Sources []string

// MustHold fails the test when the scanned set lacks name. A wrapper names its own
// package source, so a census pointed at another directory reds here.
func (s Sources) MustHold(t testing.TB, name string) {
	t.Helper()
	if !slices.Contains(s, name) {
		t.Fatalf("the census scanned %v, which lacks %s: it read the wrong directory", []string(s), name)
	}
}

// Scan grades every Go source in dir against policy and reports one diagnostic per
// violation. Test files obey the same rule, so a fixture cannot smuggle an effect in.
// It returns the scanned file names, so the caller proves the directory it read.
func Scan(t testing.TB, dir string, policy Policy) Sources {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var scanned Sources
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		scanned = append(scanned, name)
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		for _, diagnostic := range diagnose(name, string(data), policy) {
			t.Error(diagnostic)
		}
	}
	if len(scanned) == 0 {
		t.Fatal("census found no Go sources")
	}
	return scanned
}

// diagnose grades one source's content and returns one message per violation, each
// carrying the file and the line. A comment cannot trip a check, so the scan reads the
// code before the first comment marker only.
func diagnose(name, content string, policy Policy) []string {
	var diagnostics []string
	exempt := name == policy.selfExempt
	for i, line := range strings.Split(content, "\n") {
		code := line
		if idx := strings.Index(code, "//"); idx >= 0 {
			code = code[:idx]
		}
		if match := policy.forbiddenImport.FindString(code); match != "" && match != helperImport {
			diagnostics = append(diagnostics, fmt.Sprintf("%s:%d: forbidden import %s in %s", name, i+1, match, policy.subject))
		}
		if exempt {
			continue
		}
		if match := policy.ambientEffect.FindString(code); match != "" {
			diagnostics = append(diagnostics, fmt.Sprintf("%s:%d: ambient effect %s in %s", name, i+1, match, policy.subject))
		}
		if strings.Contains(code, parallelCall) {
			diagnostics = append(diagnostics, fmt.Sprintf("%s:%d: t.Parallel is out of scope for this spec", name, i+1))
		}
	}
	return diagnostics
}
