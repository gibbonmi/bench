// Package injectedports owns the injected-port audit. It derives the ports that the
// required packages of a graded root inject, and it grades that inventory against the
// admission rows its caller supplies.
package injectedports

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// These prefixes name the five failure modes, so a red names its cause without
// archaeology. They are exported because the canary fixture's EXPECT and the conformance
// bite proof both quote them. A message typed a second time drifts from the one that
// fires.
const (
	UnregisteredPortMessage   = "injected port has no registry row"
	MissingPortTestMessage    = "injected port registry row names a test the tree does not declare"
	EmptyPortExemptionMessage = "injected port registry row is exempt with an empty reason"
	ZeroPortInventoryMessage  = "injected port derivation found no ports"
	OrphanPortRowMessage      = "injected port registry row names a port the derivation no longer reports"
)

// Row is one derived port's disposition. The disposition is either a real-producer test
// or an exemption whose reason someone had to write down. A real-producer test drives the
// actual producer through the consuming surface. A row carrying neither is the
// empty-exemption red.
type Row struct {
	Package, Port      string
	TestFile, TestName string
	Exempt             string
}

// Policy is the advertisement half of the audit. Rows hold one disposition per derived
// port. Required names the packages whose source the audit derives, each of which must
// yield ports when the graded root holds it.
type Policy struct {
	Rows     []Row
	Required []string
}

// Check grades policy against a fresh derivation over root and returns the sorted,
// deduplicated diagnostics.
//
// A required package that the graded tree does not hold is not graded at all. An adopting
// repo has none of these directories, and a check that reds there is a check nobody
// keeps. A package the tree does hold must yield ports, the fail-closed posture; see
// ZeroPortInventoryMessage.
//
// The registry-to-derived direction is graded here too, not only derived-to-registry.
// Every row of a package whose derivation is non-zero must name a port that derivation
// still reports, or it fires OrphanPortRowMessage. Without this half, deleting an
// injection arm together with its junction test leaves a green row behind. That row stays
// green for as long as the package's other ports keep the inventory non-zero.
// ZeroPortInventoryMessage closes that same silence for a package that derives nothing at
// all.
func Check(root string, policy Policy) []string {
	registered := make(map[string]Row, len(policy.Rows))
	for _, row := range policy.Rows {
		registered[row.Package+"."+row.Port] = row
	}
	derived := make(map[string]bool, len(policy.Rows))
	gradedNonZero := make(map[string]bool, len(policy.Required))
	var diags []string
	for _, pkg := range policy.Required {
		dir := filepath.Join(root, filepath.FromSlash(pkg))
		if absentPackage(dir) {
			continue
		}
		ports, parseDiags := derivedInjectedPorts(dir, pkg)
		diags = append(diags, parseDiags...)
		if len(ports) == 0 {
			diags = append(diags, fmt.Sprintf(
				"%s: %s is present in the graded tree and the audit records it as declaring them, so an empty inventory would grade every registry row as satisfied — the derivation fails closed instead of reporting silence",
				ZeroPortInventoryMessage, pkg))
			continue
		}
		gradedNonZero[pkg] = true
		for _, port := range ports {
			key := port.pkg + "." + port.name
			derived[key] = true
			row, found := registered[key]
			if !found {
				diags = append(diags, fmt.Sprintf(
					"%s: %s.%s is injected at %s — add a row to injectedPortRegistry naming the test that drives its real producer, or the written reason it has none; an unregistered port is the audit's gap recurring unseen",
					UnregisteredPortMessage, port.pkg, port.name, port.site))
				continue
			}
			diags = append(diags, rowDiags(root, row)...)
		}
	}
	for _, row := range policy.Rows {
		if !gradedNonZero[row.Package] || derived[row.Package+"."+row.Port] {
			continue
		}
		diags = append(diags, fmt.Sprintf(
			"%s: %s.%s is registered but the derivation no longer reports it — a stale advertisement whose port or injection arm left the tree must not stay green",
			OrphanPortRowMessage, row.Package, row.Port))
	}
	slices.Sort(diags)
	return slices.Compact(diags)
}

// rowDiags grades one row's disposition.
func rowDiags(root string, row Row) []string {
	if row.TestName == "" {
		if strings.TrimSpace(row.Exempt) == "" {
			return []string{fmt.Sprintf(
				"%s: %s.%s names no real-producer test and gives no reason — an exemption nobody had to write down is an unregistered port with better manners",
				EmptyPortExemptionMessage, row.Package, row.Port)}
		}
		return nil
	}
	if !fileDeclaresTest(filepath.Join(root, filepath.FromSlash(row.TestFile)), row.TestName) {
		return []string{fmt.Sprintf(
			"%s: %s.%s advertises %s:%s, which that file does not declare — a deleted or renamed junction test must not leave a green advertisement behind (tripwire: this catches removal, not a test gutted back to a fake)",
			MissingPortTestMessage, row.Package, row.Port, row.TestFile, row.TestName)}
	}
	return nil
}

// absentPackage reports whether the graded tree lacks dir. Only a missing path or a
// non-directory is absent. Any other stat failure leaves the package present, so the
// derivation reports that it cannot read it instead of grading silence.
func absentPackage(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return errors.Is(err, fs.ErrNotExist)
	}
	return !info.IsDir()
}

// fileDeclaresTest reports whether path declares a top-level func of the given name. It
// parses rather than searches, so a name written in a comment or a string is not a match.
// The tripwire must answer for the test the binary would run.
func fileDeclaresTest(path, name string) bool {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return false
	}
	for _, decl := range file.Decls {
		fn, isFunc := decl.(*ast.FuncDecl)
		if isFunc && fn.Recv == nil && fn.Name.Name == name {
			return true
		}
	}
	return false
}
