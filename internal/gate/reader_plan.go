package gate

// The landing plan: a gate run that carries a landing value runs the reader set unless a case
// of the plan holds. The first case that holds sends the run to the whole gate, and a landing
// names that case before the phases.

import (
	"context"
	"fmt"
	"io"
	"slices"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/packagesurface"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// readerChainBound is the count of first-parent commits after the green marker at which a
// landing runs the whole gate.
const readerChainBound = 5

// readerRefusedFormat is the line that names the first case that refused the reader set.
const readerRefusedFormat = "gate: whole gate (reader set refused: %s)\n"

// The cases of the landing plan, in plan order.
const (
	refusedNotKit     = "not a kit run of the built-in table"
	refusedNoLanding  = "no landing value"
	refusedStack      = "landing stack"
	refusedNoMarker   = "no green marker of the destination"
	refusedUncovered  = "uncovered path"
	refusedBuildInput = "embed or shellcheck path"
	refusedLink       = "link or gitlink"
	refusedNoTable    = "no reader table"
	refusedIdentity   = "new oracle identity"
	refusedHold       = "reader hold"
	refusedChainGap   = "chain gap"
	refusedChainBound = "reader landing bound"
)

type landingKey struct{}

// landingValue is the landing value of a gate run: the destination commit and the green
// marker commit. An empty marker is an absent marker.
type landingValue struct {
	destination, marker string
}

// WithLanding puts the destination commit and the green marker commit of a reviewed landing
// in ctx, so the gate run reads the landing plan. An empty marker names no green marker.
func WithLanding(ctx context.Context, destination, marker string) context.Context {
	return context.WithValue(ctx, landingKey{}, landingValue{destination: destination, marker: marker})
}

// landingOf answers the landing value of ctx, and false when ctx carries none.
func landingOf(ctx context.Context) (landingValue, bool) {
	landing, landed := ctx.Value(landingKey{}).(landingValue)
	return landing, landed
}

// readerPlanner runs the landing plan. shellcheck answers the shellcheck inventory of a root.
type readerPlanner struct {
	shellcheck func(root string) []string
}

// plan answers the reader run of plan and whether the run is the reader set. When a case
// refuses the reader set of a kit landing, plan prints the case to stdout. A run that is not
// a kit run, and a run without a landing value, print no line, so a linked repository and a
// plain gate run keep their output.
func (p readerPlanner) plan(ctx context.Context, stdout io.Writer, runtimeRoot, storageRoot string, plan subject, selection *runbinary.Selection) (readerRun, bool) {
	landing, landed := landingOf(ctx)
	run, refusal := p.decide(landing, landed, stackGradeOf(ctx) != nil, runtimeRoot, storageRoot, plan, selection)
	if refusal == "" {
		return run, true
	}
	if landed && refusal != refusedNotKit {
		fmt.Fprintf(stdout, readerRefusedFormat, refusal)
	}
	return readerRun{}, false
}

// decide answers the reader run, or the first case of the landing plan that holds. A read
// that fails refuses under the case that reads it.
func (p readerPlanner) decide(landing landingValue, landed, stacked bool, runtimeRoot, storageRoot string, plan subject, selection *runbinary.Selection) (readerRun, string) {
	if !kitRun(runtimeRoot, storageRoot, plan, selection) {
		return readerRun{}, refusedNotKit
	}
	if !landed {
		return readerRun{}, refusedNoLanding
	}
	if stacked {
		return readerRun{}, refusedStack
	}
	if landing.marker == "" {
		return readerRun{}, refusedNoMarker
	}
	if ancestor, err := benchgit.IsAncestor(storageRoot, landing.marker, landing.destination); err != nil || !ancestor {
		return readerRun{}, refusedNoMarker
	}
	destinationTree, err := commitTree(storageRoot, landing.destination)
	if err != nil {
		return readerRun{}, refusedNoMarker
	}
	// A gitlink change refuses the reader set, so the list keeps it whatever the configuration.
	changes, err := benchgit.TreeChangesIncludingSubmodules(storageRoot, landing.destination+"^{tree}", plan.Tree)
	if err != nil {
		return readerRun{}, refusedUncovered
	}
	var families []string
	for _, change := range changes {
		family, covered := familyOf(coveredFamilies, change.Path)
		if !covered {
			return readerRun{}, refusedUncovered
		}
		if !slices.Contains(families, family) {
			families = append(families, family)
		}
	}
	directives, err := packagesurface.ScanEmbedDirectives(runtimeRoot)
	if err != nil {
		return readerRun{}, refusedBuildInput
	}
	inventory := p.shellcheck(runtimeRoot)
	for _, change := range changes {
		if directives.Matches(change.Path) || slices.Contains(inventory, change.Path) {
			return readerRun{}, refusedBuildInput
		}
	}
	if slices.ContainsFunc(changes, linkChange) {
		return readerRun{}, refusedLink
	}
	markerTree, err := commitTree(storageRoot, landing.marker)
	if err != nil {
		return readerRun{}, refusedNoTable
	}
	table, err := loadReaderTable(storageRoot, markerTree)
	if err != nil {
		return readerRun{}, refusedNoTable
	}
	// An open subject has no identity that a table can share.
	if !plan.Closed || table.Identity != plan.Identity {
		return readerRun{}, refusedIdentity
	}
	if held, err := readerHeld(storageRoot, landing.marker); err != nil || held {
		return readerRun{}, refusedHold
	}
	chain, landings, err := readerChain(storageRoot, landing)
	if err != nil || len(landings) != len(chain) {
		return readerRun{}, refusedChainGap
	}
	if len(chain) >= readerChainBound {
		return readerRun{}, refusedChainBound
	}
	return newReaderRun(plan, readerScheduleOf(table.Readers, families), destinationTree, landing.marker), ""
}

// readerChain answers the first-parent chain after the green marker of landing, up to its
// destination, and the commits of that chain whose tree holds a reader record.
func readerChain(root string, landing landingValue) (chain, landings []string, err error) {
	chain, err = benchgit.FirstParentsAfter(root, landing.marker, landing.destination)
	if err != nil {
		return nil, nil, err
	}
	for _, commit := range chain {
		tree, err := commitTree(root, commit)
		if err != nil {
			return nil, nil, err
		}
		recorded, err := ReaderRecorded(root, tree)
		if err != nil {
			return nil, nil, err
		}
		if recorded {
			landings = append(landings, commit)
		}
	}
	return chain, landings, nil
}

// commitTree answers the tree of commit in root.
func commitTree(root, commit string) (string, error) {
	return benchgit.Output("-C", root, "rev-parse", "--verify", commit+"^{tree}")
}

// readerScheduleOf answers the reader set of the changed families: the readers of each
// family in table, the declared child readers of each family, and the root-conformance
// package, each once and sorted. The race flag holds when a race-table package is among them.
func readerScheduleOf(table readerTable, families []string) readerSchedule {
	packages := []string{conformancePackagePath}
	for _, family := range families {
		packages = append(packages, table[family]...)
		for _, child := range childReaders {
			if child.Family == family {
				packages = append(packages, child.Package)
			}
		}
	}
	slices.Sort(packages)
	var schedule readerSchedule
	for _, pkg := range slices.Compact(packages) {
		operand := "./" + pkg
		schedule.Packages = append(schedule.Packages, operand)
		schedule.Race = schedule.Race || slices.ContainsFunc(raceTests, func(test raceTest) bool { return test.PackagePath == operand })
	}
	return schedule
}
