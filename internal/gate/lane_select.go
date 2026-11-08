package gate

// The lane's subject: what the composed tree changes against its base, and the lane the
// change list is graded under.

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/packagesurface"
	"github.com/gibbonmi/bench/internal/prose"
)

// ComposedChange is one raw change in the lane's subject.
type ComposedChange = benchgit.TreeChange

// ComposedChanges lists what the tree changes against the base commit's tree.
func ComposedChanges(root, base, tree string) ([]ComposedChange, error) {
	changes, err := benchgit.TreeChanges(root, base+"^{tree}", tree)
	if err != nil {
		return nil, fmt.Errorf("gate: composed change list unavailable: %w", err)
	}
	return changes, nil
}

// proseSubject answers the prose placeholder's paths: the changed Markdown the composed
// tree holds as a regular file. A deletion leaves the tree no file to grade, so it
// contributes no subject.
func proseSubject(changes []ComposedChange) []string {
	var paths []string
	for _, change := range changes {
		if prose.IsSubjectName(change.Path) && regularFile(change.DstMode) {
			paths = append(paths, change.Path)
		}
	}
	return paths
}

func regularFile(mode string) bool {
	return mode == "100644" || mode == "100755"
}

// Lane is the fast lane a worktree commit runs in place of the whole-project gate.
// Checks is the declared check list. Kit is the source root the Bench-owned checks are
// built from, and it is empty when the graded root is the kit root itself, which selects
// the private checkout of the composed tree: the kit grades with its own code. Selective
// runs the checks the composed changes select, which the kit's built-in lane does and a
// project's declared lane does not.
type Lane struct {
	Checks    []Phase
	Kit       string
	Selective bool
}

// LaneRequest is one lane run. Root is the repository whose Git dir receives the record
// and whose object store holds Tree. Tree is the composed snapshot the lane grades. Lane
// names the lane in its record. Checks is the declared check list, resolved through
// LaneFor. Changes is the composed change list the prose placeholder and the selection
// resolve from. Selective runs the classes the changes carry rather than the whole
// declared list. Kit is the source root the run binary is built from; empty selects the
// private checkout, which is the composed tree itself. Base is the commit the graded tree
// is measured against, and the base placeholder resolves to it.
type LaneRequest struct {
	Root      string
	Kit       string
	Tree      string
	Lane      string
	Base      string
	Checks    []Phase
	Changes   []ComposedChange
	Selective bool
	Stdout    io.Writer
	Stderr    io.Writer
}

// LaneResult is what one lane run decided. Outcome is "pass" or "fail". Check names the
// first check that failed, and Diagnostic is that check's first output line, so a caller
// can name the failure without re-reading the stream. Checks names the checks the run
// actually graded, and Classes names the path classes that selected them, so a caller
// states the lane line without repeating the selection. RunBinary is the content address
// of the executable the Bench-owned checks ran.
type LaneResult struct {
	Outcome    string
	Check      string
	Diagnostic string
	Checks     []string
	Classes    []string
	Tree       string
	Lane       string
	RunBinary  string
	RecordedAt time.Time
}

// Passed reports the one question a caller acts on.
func (r LaneResult) Passed() bool { return r.Outcome == lanePass }

// UnknownClass is the class of a path no row claims. It selects every declared check,
// because a path the table does not know is a path the lane cannot narrow safely.
const UnknownClass = "unknown"

// PathClass is one row of the lane's path-class table. Name is the class a lane line
// reports. Match answers whether a path belongs to the class, and it reads facts of the
// composed checkout because some rows derive their membership rather than spell it.
// Checks names the lane checks the class selects, by check name, so a row can name a
// check the declared lane does not carry.
type PathClass struct {
	Name   string
	Match  func(path string, tree TreeFacts) bool
	Checks []string
}

// TreeFacts are the path sets the lane reads from the composed checkout. EmbedTargets
// are the targets of the Go embed directives, and FixturePins are the paths the retained
// canary fixtures pin.
type TreeFacts struct {
	EmbedTargets []string
	FixturePins  []string
}

// laneClasses orders the content classes, the anchor class, and the document families.
// It is the one source for what a composed change selects, and the profile's
// `selected by` column renders from it.
var laneClasses = append([]PathClass{
	{
		Name:   "go-source",
		Match:  func(path string, _ TreeFacts) bool { return strings.HasSuffix(path, ".go") },
		Checks: []string{"gofmt", "vet", "build", "structure"},
	},
	{
		Name: "go-build-input",
		Match: func(path string, tree TreeFacts) bool {
			return path == "go.mod" || path == "go.sum" || slices.Contains(tree.EmbedTargets, path)
		},
		Checks: []string{"vet", "build"},
	},
	// The docs-currency and prose-budget checks read Markdown across the tree, so every
	// Markdown change selects them. The class claims no copy of either check's file list.
	{
		Name:   "markdown",
		Match:  func(path string, _ TreeFacts) bool { return prose.IsSubjectName(path) },
		Checks: []string{"prose", "docs-currency-workflow", "guidance-prose-budgets"},
	},
	{
		Name:   "prose-policy",
		Match:  func(path string, _ TreeFacts) bool { return path == prose.ExclusionFile },
		Checks: []string{"prose"},
	},
	anchorRegistryClass(),
	// A fixture copies or quotes the bytes of each path it pins, so a reflow of a pinned
	// line can stop the fixture from biting while every document check stays green.
	{
		Name:   "canary-pin",
		Match:  func(path string, tree TreeFacts) bool { return slices.Contains(tree.FixturePins, path) },
		Checks: []string{fixtureBiteCheck},
	},
	// The registry owns the files that the route bypass test scans, so the class claims
	// exactly those files.
	{
		Name:   "write-verb-source",
		Match:  func(path string, _ TreeFacts) bool { return registry.RouteBypassScanned(path) },
		Checks: []string{routeBypassCheck},
	},
}, documentClasses()...)

func anchorRegistryClass() PathClass {
	files := map[string]bool{}
	for _, anchor := range anchors.Entries() {
		files[anchor.File] = true
	}
	return PathClass{
		Name:   "anchor-registry",
		Match:  func(path string, _ TreeFacts) bool { return files[path] },
		Checks: []string{"docs-currency-workflow"},
	}
}

// documentFamilies binds each document class to the paths it claims. A row's name is the
// registry input source itself, so the lane and the registry cannot spell one family two
// ways, and a row states no check name: the registry's own binding answers that.
var documentFamilies = []struct {
	source registry.InputSource
	match  func(path string) bool
}{
	{registry.InputRoadmapBoard, func(path string) bool {
		return path == "ROADMAP.md" || strings.HasPrefix(path, "roadmap/")
	}},
	{registry.InputDecisionDocuments, decisionDocument},
	{registry.InputCaptureRetros, func(path string) bool {
		return strings.HasPrefix(path, "capture/retros/")
	}},
	{registry.InputBenchkitProfile, func(path string) bool {
		return path == "projects/benchkit.md"
	}},
}

// decisionDocument claims the two places a decision map lives: the repository's own
// `decisions/` tree, and the `decisions/` tree a spec folder carries. A rule that read
// only the top-level tree would miss every spec-local map.
func decisionDocument(path string) bool {
	if strings.HasPrefix(path, "decisions/") {
		return true
	}
	rest, found := strings.CutPrefix(path, "specs/")
	if !found {
		return false
	}
	slug, tail, split := strings.Cut(rest, "/")
	return split && slug != "" && strings.HasPrefix(tail, "decisions/")
}

// documentClasses renders the document families as class-table rows, in family order.
func documentClasses() []PathClass {
	classes := make([]PathClass, 0, len(documentFamilies))
	for _, family := range documentFamilies {
		match := family.match
		classes = append(classes, PathClass{
			Name:   string(family.source),
			Match:  func(path string, _ TreeFacts) bool { return match(path) },
			Checks: documentRegistryChecks(family.source),
		})
	}
	return classes
}

// documentRegistryChecks names the dev-tier checks bound to these input sources, in
// registry order. Document classes derive their checks from these bindings.
func documentRegistryChecks(sources ...registry.InputSource) []string {
	var names []string
	for _, check := range registry.Checks {
		if !slices.Contains(sources, check.Inputs) || !check.RunsAt(registry.Dev) {
			continue
		}
		names = append(names, check.Name)
	}
	return names
}

// documentLaneChecks declares the registry checks that path classes select, in registry
// order. Each check uses the lane's run binary to grade the composed checkout.
func documentLaneChecks() []Phase {
	var checks []Phase
	for _, check := range registry.Checks {
		for _, class := range laneClasses {
			if check.RunsAt(registry.Dev) && slices.Contains(class.Checks, check.Name) {
				checks = append(checks, Phase{Name: check.Name, Argv: []string{runBinaryArgvToken, "test", "--check", check.Name}})
				break
			}
		}
	}
	return checks
}

// LaneClasses answers the path-class table in table order. It is the read seam for a
// caller that advertises the table rather than applies it.
func LaneClasses() []PathClass {
	return slices.Clone(laneClasses)
}

// SelectLane answers the declared checks the composed changes select, in declared order
// and without a duplicate, and the classes that selected them in table order. A class
// that names a check the lane does not declare adds nothing, so the result is always a
// subsequence of the declared list. The unknown class selects every declared check, and
// it is reported last because the table does not carry it.
func SelectLane(checks []Phase, changes []ComposedChange, tree TreeFacts) ([]Phase, []string) {
	matched := map[string]bool{}
	for _, change := range changes {
		for _, name := range changeClasses(change, tree) {
			matched[name] = true
		}
	}
	var classes []string
	selected := map[string]bool{}
	for _, class := range laneClasses {
		if !matched[class.Name] {
			continue
		}
		classes = append(classes, class.Name)
		for _, name := range class.Checks {
			selected[name] = true
		}
	}
	if matched[UnknownClass] {
		return checks, append(classes, UnknownClass)
	}
	var kept []Phase
	for _, check := range checks {
		if selected[check.Name] {
			kept = append(kept, check)
		}
	}
	return kept, classes
}

// linkChange reports whether either side of change is a symbolic link or a gitlink. A
// symbolic link's bytes are a path and not the content its name suggests, and a gitlink
// names a commit in another repository, so neither mode's name tells what a check would
// grade. Either side's mode is read, so a typechange to or from a link, and a deleted
// submodule pointer, count.
func linkChange(change ComposedChange) bool {
	for _, mode := range []string{change.SrcMode, change.DstMode} {
		if mode == benchgit.SymlinkMode || mode == benchgit.GitlinkMode {
			return true
		}
	}
	return false
}

// changeClasses answers the classes one composed change carries. A link change is unknown
// rather than classified by a name the tree no longer backs.
func changeClasses(change ComposedChange, tree TreeFacts) []string {
	if linkChange(change) {
		return []string{UnknownClass}
	}
	var names []string
	for _, class := range laneClasses {
		if class.Match(change.Path, tree) {
			names = append(names, class.Name)
		}
	}
	if len(names) == 0 {
		return []string{UnknownClass}
	}
	return names
}

// selectLaneChecks answers the checks one run grades, their names, and the classes that
// selected them. A lane that is not selective runs its declared list, and it reports no
// class: a project's own lane keeps the meaning its manifest gives it. A selective lane
// reads the embed targets and the fixture pins from the composed checkout, so the
// membership rule grades the tree under commit rather than the working tree beside it.
func selectLaneChecks(req LaneRequest, checkout string) (checks []Phase, names, classes []string, err error) {
	if !req.Selective {
		return req.Checks, laneCheckNames(req.Checks), nil, nil
	}
	embedTargets, err := packagesurface.EmbedTargets(checkout)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gate: lane embed targets unavailable: %w", err)
	}
	pins, err := canary.FixturePins(checkout)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gate: lane fixture pins unavailable: %w", err)
	}
	tree := TreeFacts{EmbedTargets: embedTargets, FixturePins: slices.Sorted(maps.Keys(pins))}
	selected, classes := SelectLane(req.Checks, req.Changes, tree)
	if len(selected) == 0 {
		return nil, nil, nil, errors.New("gate: lane selects no check")
	}
	return selected, laneCheckNames(selected), classes, nil
}

func laneCheckNames(checks []Phase) []string {
	names := make([]string, len(checks))
	for i, check := range checks {
		names[i] = check.Name
	}
	return names
}
