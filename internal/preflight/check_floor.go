package preflight

import (
	"context"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/tickets"
)

const checkFloorName = "check-floor"

// testCommand is the command that each red names.
const testCommand = reviewrecord.TestCommand

// commandPackageElement is the floor element of the command tests. A ticket owes it when
// it writes into a bound command package or touches an embed.
const commandPackageElement = "./cmd/bench"

// CheckFloorFacts is the evidence check-floor grades: the floor of each ticket, and the
// selection of each plan command. The gatherer composes both; whether an element is
// satisfied is Decide's.
type CheckFloorFacts struct {
	// Elements holds the floor elements of each ticket, keyed by ticket basename.
	Elements map[string][]string
	// Selections holds the prepared selection of each plan command that is one simple
	// `bench test` call that the grammar accepts, keyed by the command text. Any other
	// command is absent, so it satisfies nothing.
	Selections map[string]testreport.Request
	// LoaderError is the package loader message, empty when every floor loaded.
	LoaderError string
}

// gatherCheckFloor composes the floor of each ticket and prepares each plan command. A
// loader failure stops the composition, because a partial floor would grade as a false
// green.
func gatherCheckFloor(ctx context.Context, root string, parsed []tickets.Ticket, plan reviewrecord.Plan) CheckFloorFacts {
	facts := CheckFloorFacts{Elements: map[string][]string{}, Selections: map[string]testreport.Request{}}
	conformance := gate.RootConformanceRuns(root, gate.KitRoot(root))
	for _, ticket := range parsed {
		elements, err := floorElements(ctx, root, ticket.Writes, conformance)
		if err != nil {
			facts.LoaderError = err.Error()
			return facts
		}
		facts.Elements[ticket.Name] = elements
	}
	for _, chunk := range plan.Chunks {
		for _, requirement := range chunk.Verification {
			if selection, ok := testreport.PlanSelection(root, requirement.Command); ok {
				facts.Selections[requirement.Command] = selection
			}
		}
	}
	return facts
}

// floorElements is one ticket's floor: root conformance where it applies, each check and
// each package that the ticket writes, and the command tests for a command package or an
// embed. An entry can name a command package that holds no package yet, and a directory
// entry can hold a command package below it, so both the entries and the floor packages
// are asked.
func floorElements(ctx context.Context, root string, writes []string, conformance bool) ([]string, error) {
	floor, err := testreport.WritesFloor(ctx, root, writes)
	if err != nil {
		return nil, err
	}
	var elements []string
	if conformance {
		elements = append(elements, testreport.RootConformanceCheckName)
	}
	elements = append(append(elements, floor.Checks...), floor.Packages...)
	commandWrite := slices.ContainsFunc(writes, func(entry string) bool {
		path, _ := tickets.WritesPath(entry)
		return tickets.CommandPackage(path)
	}) || slices.ContainsFunc(floor.Packages, func(pattern string) bool {
		return tickets.CommandPackage(strings.TrimPrefix(pattern, "./"))
	})
	if (floor.Embed || commandWrite) && !slices.Contains(elements, commandPackageElement) {
		elements = append(elements, commandPackageElement)
	}
	return elements, nil
}

// checkFloorCheck grades each chunk once: the union of its tickets' floors against the
// union of its requirements, whatever ticket a requirement names. A ticket that no chunk
// names meets no element. An invalid plan is completion-plan's red, so this row does not
// grade it a second time.
func checkFloorCheck(f Facts) CheckResult {
	if f.CompletionPlanError != "" {
		return notApplicable(checkFloorName)
	}
	if f.CheckFloor.LoaderError != "" {
		return red(checkFloorName, "package index not readable: "+f.CheckFloor.LoaderError)
	}
	var omitted []string
	planned := map[string]bool{}
	for _, chunk := range f.CompletionPlan.Chunks {
		var elements []string
		for _, name := range chunk.Tickets {
			planned[name] = true
			for _, element := range f.CheckFloor.Elements[name] {
				if !slices.Contains(elements, element) {
					elements = append(elements, element)
				}
			}
		}
		omitted = append(omitted, f.CheckFloor.omitted(chunkOwner(chunk.ID), elements, chunk.Verification)...)
	}
	for _, ticket := range f.Tickets {
		if !planned[ticket.Name] {
			omitted = append(omitted, f.CheckFloor.omitted(ticket.Name, f.CheckFloor.Elements[ticket.Name], nil)...)
		}
	}
	if len(omitted) > 0 {
		return red(checkFloorName, floorRedDetail(omitted))
	}
	return green(checkFloorName)
}

// omitted is the red item of each element of owner that no requirement satisfies.
func (facts CheckFloorFacts) omitted(owner string, elements []string, requirements []reviewrecord.Requirement) []string {
	var omitted []string
	for _, element := range elements {
		satisfied := slices.ContainsFunc(requirements, func(requirement reviewrecord.Requirement) bool {
			selection, ok := facts.Selections[requirement.Command]
			return ok && satisfies(selection, element)
		})
		if !satisfied {
			omitted = append(omitted, omittedElement(owner, element))
		}
	}
	return omitted
}

// chunkOwner is the owner name of a chunk's red items.
func chunkOwner(id string) string {
	return "chunk " + id
}

// omittedElement is the red item of one unsatisfied element of owner, a chunk or a ticket
// that no chunk names.
func omittedElement(owner, element string) string {
	return owner + ": " + element + " (add " + floorCommand(element) + ")"
}

// floorRedDetail is the red detail over the omitted items, in plan chunk order and then in
// the order of the tickets that no chunk names.
func floorRedDetail(omitted []string) string {
	return "completion plan omits a floor check: " + strings.Join(omitted, "; ")
}

// satisfies applies the satisfaction rules to one selection. A form that runs no test
// satisfies nothing, and a run with a run pattern is filtered, `--run all` included. Only an
// unfiltered run of the named check satisfies a check element. A package element needs an
// unfiltered package run that is not `--changed`, whose target equals the element or covers
// it at a segment boundary. A run of several package expressions satisfies each element
// that one of its expressions satisfies.
func satisfies(selection testreport.Request, element string) bool {
	if !selection.Runs() {
		return false
	}
	if testreport.IsNamedCheck(element) {
		return selection.Form() == testreport.CheckForm && selection.Target() == element && !selection.Filtered()
	}
	if selection.Form() != testreport.PackageForm || selection.Filtered() || selection.Changed() {
		return false
	}
	return slices.ContainsFunc(selection.Packages(), func(target string) bool {
		return testreport.PackageCovers(target, element)
	})
}

// floorCommand is the exact command that satisfies one element.
func floorCommand(element string) string {
	if testreport.IsNamedCheck(element) {
		return testCommand + " --check " + element
	}
	return reviewrecord.PackageTestCommand(element)
}
