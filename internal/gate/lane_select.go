package gate

// The lane's subject: what the composed tree changes against its base, and the lane the
// change list is graded under.

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/packagesurface"
	"github.com/gibbonmi/bench/internal/prose"
)

// composedChangeFields is the field count of one raw-diff entry's metadata: the two
// modes, the two object IDs, and the status letter.
const composedChangeFields = 5

// ComposedChange is one entry of the raw diff between the base tree and the composed
// tree. Status is Git's own status letter. SrcMode and DstMode are the six-digit modes
// of the two sides, and one of them is `000000` when that side holds nothing. Path
// carries the file's own bytes, so a name with a space or a byte above ASCII survives.
type ComposedChange struct {
	Status  string
	SrcMode string
	DstMode string
	Path    string
}

// ComposedChanges lists what the tree changes against the base commit's tree. It is the
// one derivation of a lane's subject: a caller that named a directory reaches the files
// under it, and a caller that named nothing at all reaches an empty list.
//
// Rename detection is off, so a rename arrives as one deletion and one addition and each
// side classifies by its own path. The NUL framing is load-bearing: under the default
// `core.quotepath` a newline-framed name with a byte above ASCII arrives C-quoted, so a
// reader would carry a path no file has.
func ComposedChanges(root, base, tree string) ([]ComposedChange, error) {
	raw, err := benchgit.Raw("-C", root, "diff", "--raw", "--no-renames", "-z", base+"^{tree}", tree)
	if err != nil {
		return nil, fmt.Errorf("gate: composed change list unavailable: %w", err)
	}
	frames := strings.Split(string(raw), "\x00")
	var changes []ComposedChange
	for i := 0; i+1 < len(frames); i += 2 {
		change, err := parseComposedChange(frames[i], frames[i+1])
		if err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	return changes, nil
}

// parseComposedChange reads one `--raw -z` entry, whose metadata frame is
// `:<srcmode> <dstmode> <srcsha> <dstsha> <status>` and whose path is the frame after it.
func parseComposedChange(meta, path string) (ComposedChange, error) {
	fields := strings.Fields(strings.TrimPrefix(meta, ":"))
	if !strings.HasPrefix(meta, ":") || len(fields) != composedChangeFields || path == "" {
		return ComposedChange{}, fmt.Errorf("gate: unreadable composed change entry %q", meta)
	}
	return ComposedChange{Status: fields[4], SrcMode: fields[0], DstMode: fields[1], Path: path}, nil
}

// proseSubject answers the prose placeholder's paths: the changed Markdown the composed
// tree holds as a regular file. A deletion leaves the tree no file to grade, so it
// contributes no subject.
func proseSubject(changes []ComposedChange) []string {
	var paths []string
	for _, change := range changes {
		if strings.HasSuffix(change.Path, ".md") && regularFile(change.DstMode) {
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

// LaneForCommit resolves the lane a worktree commit at root runs. It answers a nil lane
// for a root that declares none, so one read tells a caller both whether a lane exists
// and what it is. It applies the gate's own kit-root selection, so a caller outside this
// package asks the lane question once.
func LaneForCommit(root string) (*Lane, error) {
	kit := kitRoot(root)
	checks, err := LaneFor(root, kit)
	if err != nil || checks == nil {
		return nil, err
	}
	// The one built-in-versus-manifest decision. LaneFor answers the built-in lane under
	// this same predicate, so the selection switch and the lane come from one call.
	if sameDirectory(root, kit) {
		return &Lane{Checks: checks, Selective: true}, nil
	}
	return &Lane{Checks: checks, Kit: kit}, nil
}

// LaneRequest is one lane run. Root is the repository whose Git dir receives the record
// and whose object store holds Tree. Tree is the composed snapshot the lane grades. Lane
// names the lane in its record. Checks is the declared check list, resolved through
// LaneFor. Changes is the composed change list the prose placeholder and the selection
// resolve from. Selective runs the classes the changes carry rather than the whole
// declared list. Kit is the source root the run binary is built from; empty selects the
// private checkout, which is the composed tree itself.
type LaneRequest struct {
	Root      string
	Kit       string
	Tree      string
	Lane      string
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

// The two Git modes that make a change unclassifiable. A symbolic link's bytes are a
// path and not the content its name suggests, and a gitlink names a commit in another
// repository. Neither side's name tells the lane what a check would grade.
const (
	symlinkMode = "120000"
	gitlinkMode = "160000"
)

// UnknownClass is the class of a path no row claims. It selects every declared check,
// because a path the table does not know is a path the lane cannot narrow safely.
const UnknownClass = "unknown"

// PathClass is one row of the lane's path-class table. Name is the class a lane line
// reports. Match answers whether a path belongs to the class, and it reads the composed
// checkout's embed targets because one row's membership is derived rather than spelled.
// Checks names the lane checks the class selects, by check name, so a row can name a
// check the declared lane does not carry.
type PathClass struct {
	Name   string
	Match  func(path string, embedTargets []string) bool
	Checks []string
}

// laneClasses is the path-class table in table order: the four content classes, then the
// document families the registry binds. It is the one source for what a composed change
// selects, and the profile's `selected by` column renders from it.
var laneClasses = append([]PathClass{
	{
		Name:   "go-source",
		Match:  func(path string, _ []string) bool { return strings.HasSuffix(path, ".go") },
		Checks: []string{"gofmt", "vet", "build"},
	},
	{
		Name: "go-build-input",
		Match: func(path string, embedTargets []string) bool {
			return path == "go.mod" || path == "go.sum" || slices.Contains(embedTargets, path)
		},
		Checks: []string{"vet", "build"},
	},
	{
		Name:   "markdown",
		Match:  func(path string, _ []string) bool { return strings.HasSuffix(path, ".md") },
		Checks: []string{"prose"},
	},
	{
		Name:   "prose-policy",
		Match:  func(path string, _ []string) bool { return path == prose.ExclusionFile },
		Checks: []string{"prose"},
	},
}, documentClasses()...)

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
			Match:  func(path string, _ []string) bool { return match(path) },
			Checks: documentRegistryChecks(family.source),
		})
	}
	return classes
}

// documentRegistryChecks names the dev-tier registry checks bound to any of the given
// input sources, in registry order. It is the one derivation of the family-to-check fact:
// a class row reads it for its own source, and the kit lane declares its document rows
// from the whole set. A check the registry adds therefore joins both with no second list.
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

// documentLaneChecks is the kit lane's document half: one check per dev-tier registry
// check a document family binds, in registry order. Each runs through the lane's own run
// binary, so the lane builds no second executable and the check grades the composed
// checkout that binary was built from.
func documentLaneChecks() []Phase {
	sources := make([]registry.InputSource, 0, len(documentFamilies))
	for _, family := range documentFamilies {
		sources = append(sources, family.source)
	}
	names := documentRegistryChecks(sources...)
	checks := make([]Phase, 0, len(names))
	for _, name := range names {
		checks = append(checks, Phase{Name: name, Argv: []string{runBinaryArgvToken, "test", "--check", name}})
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
func SelectLane(checks []Phase, changes []ComposedChange, embedTargets []string) ([]Phase, []string) {
	matched := map[string]bool{}
	for _, change := range changes {
		for _, name := range changeClasses(change, embedTargets) {
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

// changeClasses answers the classes one composed change carries. Either side's mode is
// read, so a typechange to or from a link, and a deleted submodule pointer, are unknown
// rather than classified by a name the tree no longer backs.
func changeClasses(change ComposedChange, embedTargets []string) []string {
	for _, mode := range []string{change.SrcMode, change.DstMode} {
		if mode == symlinkMode || mode == gitlinkMode {
			return []string{UnknownClass}
		}
	}
	var names []string
	for _, class := range laneClasses {
		if class.Match(change.Path, embedTargets) {
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
// reads the embed targets from the composed checkout, so the membership rule grades the
// tree under commit rather than the working tree beside it.
func selectLaneChecks(req LaneRequest, checkout string) (checks []Phase, names, classes []string, err error) {
	if !req.Selective {
		return req.Checks, laneCheckNames(req.Checks), nil, nil
	}
	embedTargets, err := packagesurface.EmbedTargets(checkout)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gate: lane embed targets unavailable: %w", err)
	}
	selected, classes := SelectLane(req.Checks, req.Changes, embedTargets)
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
