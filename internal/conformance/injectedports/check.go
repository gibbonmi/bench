// Package injectedports owns the injected-port audit. It derives the test seams that the
// production packages of a graded root declare, and it grades that inventory against the
// admission rows its caller supplies.
package injectedports

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
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

// Row is one derived seam's disposition. Kind, Package, and Port are its source
// identity: Port names the port type, the package variable, or the setter function.
//
// A named-port row carries either a real-producer test or an exemption whose reason
// someone had to write down. A real-producer test drives the actual producer through the
// consuming surface. A row carrying neither is the empty-exemption red. A fault port has
// no production producer, so its row instead names Owner and the evidence fields of a
// variable row; see faultDiags.
//
// A variable row carries its policy facts: the production owner that reads the variable,
// the consumer tests ("path::TestName") that install a replacement, the failure purpose,
// and the isolation evidence. A setter row carries only Canonical, the variable seams it
// finally replaces, and borrows every other fact from their rows.
//
// A runtime row names its branch function as Port and carries the same evidence; its
// consumers reach the branch through their production calls. An environment row names
// the hook as Port, one reader body of Package as Owner, and consumers that set the hook.
// Its source identity is the hook name alone.
type Row struct {
	Kind               Kind
	Package, Port      string
	TestFile, TestName string
	Exempt             string
	Owner              string
	Consumers          []string
	Purpose, Isolation string
	Canonical          []string
}

// Kind distinguishes the seam forms that a row admits. The zero value is a named port.
type Kind string

const (
	KindPort        Kind = ""
	KindVariable    Kind = "variable"
	KindSetter      Kind = "setter"
	KindRuntime     Kind = "runtime"
	KindEnvironment Kind = "environment"
)

// Policy is the advertisement half of the audit. Rows hold one disposition per derived
// seam. Required names the packages that must yield named ports when the graded root holds
// them. It does not limit discovery: the audit derives every production package.
//
// Families supplies the members of a local name that source constants cannot name, keyed
// "pkg.owner.name", such as a phase name that a release registry of the graded root
// lists. A caller derives each family from graded-root data. A supplied empty family
// declares no members, so its reader reads no hook; a missing key leaves the name unresolved.
type Policy struct {
	Rows     []Row
	Required []string
	Families map[string][]string
}

// Check grades policy against one derivation over the production source of root and
// returns the sorted, deduplicated diagnostics. The audit loads that source once; every
// seam form derives from the same parse.
//
// A required package that the graded tree does not hold is not graded at all. An adopting
// repo has none of these directories, and a check that reds there is a check nobody
// keeps. A required package the tree does hold must yield ports, the fail-closed posture;
// see ZeroPortInventoryMessage.
//
// The registry-to-derived direction is graded here too, not only derived-to-registry.
// Every row of a package that the graded tree holds must name a seam that the derivation
// still reports, or it fires OrphanPortRowMessage. Without this half, deleting an
// injection arm together with its junction test leaves a green row behind. A required
// package that derives nothing reports ZeroPortInventoryMessage for its rows instead.
func Check(root string, policy Policy) []string {
	tree, diags := loadSourceTree(root)
	census := deriveSeams(tree)
	diags = append(diags, duplicateRows(policy.Rows)...)
	diags = append(diags, census.portDiags(root, tree, policy)...)
	diags = append(diags, census.seamDiags(tree, policy)...)
	slices.Sort(diags)
	return slices.Compact(diags)
}

// portDiags grades the named-port rows against the ports of every production package.
func (c *seamCensus) portDiags(root string, tree *sourceTree, policy Policy) []string {
	ports := c.namedPorts(tree)
	registered := map[string]Row{}
	for _, row := range policy.Rows {
		if row.Kind == KindPort {
			registered[row.Package+"."+row.Port] = row
		}
	}
	var diags []string
	empty := map[string]bool{}
	for _, pkg := range policy.Required {
		if len(ports[pkg]) > 0 || absentPackage(filepath.Join(root, filepath.FromSlash(pkg))) {
			continue
		}
		empty[pkg] = true
		diags = append(diags, fmt.Sprintf(
			"%s: %s is present in the graded tree and the audit records it as declaring them, so an empty inventory would grade every registry row as satisfied — the derivation fails closed instead of reporting silence",
			ZeroPortInventoryMessage, pkg))
	}
	derived := map[string]bool{}
	for _, pkg := range sortedKeys(ports) {
		for _, port := range ports[pkg] {
			key := port.pkg + "." + port.name
			derived[key] = true
			row, found := registered[key]
			switch {
			case !found:
				diags = append(diags, fmt.Sprintf(
					"%s: %s.%s is injected at %s — add a row to injectedPortRegistry naming the test that drives its real producer, or the written reason it has none; an unregistered port is the audit's gap recurring unseen",
					UnregisteredPortMessage, port.pkg, port.name, port.site))
			case row.Owner != "":
				diags = append(diags, c.faultDiags(row)...)
			default:
				diags = append(diags, c.rowDiags(row)...)
			}
		}
	}
	for _, row := range policy.Rows {
		if row.Kind != KindPort || derived[row.Package+"."+row.Port] || empty[row.Package] || absentPackage(filepath.Join(root, filepath.FromSlash(row.Package))) {
			continue
		}
		diags = append(diags, orphanDiag(row))
	}
	return diags
}

func orphanDiag(row Row) string {
	return fmt.Sprintf(
		"%s: %s.%s is registered but the derivation no longer reports it — a stale advertisement whose port or injection arm left the tree must not stay green",
		OrphanPortRowMessage, row.Package, row.Port)
}

// rowDiags grades one named-port row's disposition.
func (c *seamCensus) rowDiags(row Row) []string {
	if row.TestName == "" {
		if strings.TrimSpace(row.Exempt) == "" {
			return []string{fmt.Sprintf(
				"%s: %s.%s names no real-producer test and gives no reason — an exemption nobody had to write down is an unregistered port with better manners",
				EmptyPortExemptionMessage, row.Package, row.Port)}
		}
		return nil
	}
	if !c.declaresTest(row.TestFile, row.TestName) {
		return []string{fmt.Sprintf(
			"%s: %s.%s advertises %s:%s, which that file does not declare — a deleted or renamed junction test must not leave a green advertisement behind (tripwire: this catches removal, not a test gutted back to a fake)",
			MissingPortTestMessage, row.Package, row.Port, row.TestFile, row.TestName)}
	}
	return nil
}

// absentPackage reports whether the graded tree lacks dir. A missing path, a path below a
// non-directory, and a non-directory are absent. Any other stat failure leaves the package
// present, so the derivation reports that it cannot read it instead of grading silence.
func absentPackage(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR)
	}
	return !info.IsDir()
}

// declaresTest reports whether the parsed test file declares a top-level func of the given
// name. The census parsed it, so a name written in a comment or a string is not a match:
// the tripwire answers for the test the binary would run.
func (c *seamCensus) declaresTest(file, name string) bool {
	scope, known := c.testFiles[file]
	entry := c.tests[scope+"."+name]
	return known && entry != nil && entry.file.rel == file
}
