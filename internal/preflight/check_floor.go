package preflight

import (
	"context"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/shellcommand"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/tickets"
)

const checkFloorName = "check-floor"

// testCommand is the command whose plan requirements can satisfy a floor element, and the
// command that each red names.
const testCommand = "bench test"

// commandPackageElement is the floor element of the command tests. A ticket owes it when
// it writes into a bound command package or touches an embed.
const commandPackageElement = "./cmd/bench"

// CheckFloorFacts is the evidence check-floor grades: the floor of each planned ticket,
// and the selection of each plan command. The gatherer composes both; whether an element
// is satisfied is Decide's.
type CheckFloorFacts struct {
	// Elements holds the floor elements of each planned ticket, keyed by ticket basename.
	Elements map[string][]string
	// Selections holds the prepared selection of each plan command that is one simple
	// `bench test` call that the grammar accepts, keyed by the command text. Any other
	// command is absent, so it satisfies nothing.
	Selections map[string]testreport.Request
	// LoaderError is the package loader message, empty when every floor loaded.
	LoaderError string
}

// gatherCheckFloor composes the floor of each ticket that the plan names and prepares
// each plan command. A loader failure stops the composition, because a partial floor
// would grade as a false green.
func gatherCheckFloor(ctx context.Context, root string, parsed []tickets.Ticket, plan reviewrecord.Plan) CheckFloorFacts {
	facts := CheckFloorFacts{Elements: map[string][]string{}, Selections: map[string]testreport.Request{}}
	conformance := gate.RootConformanceRuns(root, gate.KitRoot(root))
	for _, ticket := range parsed {
		if _, planned := ticketRequirements(plan, ticket.Name); !planned {
			continue
		}
		elements, err := floorElements(ctx, root, ticket.Writes, conformance)
		if err != nil {
			facts.LoaderError = err.Error()
			return facts
		}
		facts.Elements[ticket.Name] = elements
	}
	for _, chunk := range plan.Chunks {
		for _, requirement := range chunk.Verification {
			if selection, ok := planSelection(root, requirement.Command); ok {
				facts.Selections[requirement.Command] = selection
			}
		}
	}
	return facts
}

// floorElements is one ticket's floor: root conformance where it applies, each package
// that the ticket writes, and the command tests for a command package or an embed.
func floorElements(ctx context.Context, root string, writes []string, conformance bool) ([]string, error) {
	floor, err := testreport.WritesFloor(ctx, root, writes)
	if err != nil {
		return nil, err
	}
	var elements []string
	if conformance {
		elements = append(elements, testreport.RootConformanceCheckName)
	}
	elements = append(elements, floor.Packages...)
	commandEntry := slices.ContainsFunc(writes, func(entry string) bool {
		path, _ := tickets.WritesPath(entry)
		return tickets.CommandPackage(path)
	})
	if (floor.Embed || commandEntry) && !slices.Contains(elements, commandPackageElement) {
		elements = append(elements, commandPackageElement)
	}
	return elements, nil
}

// planSelection prepares one plan command. Only a lexed stream of one simple command
// whose words start with `bench test` reaches the grammar, so a compound line, a foreign
// executable, or an assignment prefix satisfies nothing.
func planSelection(root, command string) (testreport.Request, bool) {
	stream := shellcommand.Parse(command)
	if stream.Unlexed || len(stream.Commands) != 1 {
		return testreport.Request{}, false
	}
	span := stream.Commands[0]
	words := shellcommand.ProjectCommandWords(stream.Tokens[span.Start:span.End])
	if len(words) < 2 || strings.Join(words[:2], " ") != testCommand {
		return testreport.Request{}, false
	}
	request, line, _ := testreport.Prepare(root, words[2:])
	return request, line == ""
}

// ticketRequirements answers the requirements that grade one ticket, and whether the plan
// names the ticket. A version 1 chunk grades every ticket against all of its requirements.
// A delegated chunk grades each ticket against the requirements that name it.
func ticketRequirements(plan reviewrecord.Plan, name string) ([]reviewrecord.Requirement, bool) {
	for _, chunk := range plan.Chunks {
		if !slices.Contains(chunk.Tickets, name) {
			continue
		}
		var owned []reviewrecord.Requirement
		for _, requirement := range chunk.Verification {
			if !plan.Delegated() || requirement.Ticket == name {
				owned = append(owned, requirement)
			}
		}
		return owned, true
	}
	return nil, false
}

// checkFloorCheck grades each planned ticket's floor against its requirements. An
// invalid plan is completion-plan's red, so this row does not grade it a second time.
func checkFloorCheck(f Facts) CheckResult {
	if f.CompletionPlanError != "" {
		return notApplicable(checkFloorName)
	}
	if f.CheckFloor.LoaderError != "" {
		return red(checkFloorName, "package index not readable: "+f.CheckFloor.LoaderError)
	}
	var omitted []string
	for _, ticket := range f.Tickets {
		requirements, _ := ticketRequirements(f.CompletionPlan, ticket.Name)
		for _, element := range f.CheckFloor.Elements[ticket.Name] {
			satisfied := slices.ContainsFunc(requirements, func(requirement reviewrecord.Requirement) bool {
				selection, ok := f.CheckFloor.Selections[requirement.Command]
				return ok && satisfies(selection, element)
			})
			if !satisfied {
				omitted = append(omitted, ticket.Name+": "+element+" (add "+floorCommand(element)+")")
			}
		}
	}
	if len(omitted) > 0 {
		return red(checkFloorName, "completion plan omits a floor check: "+strings.Join(omitted, "; "))
	}
	return green(checkFloorName)
}

// satisfies applies the satisfaction rules to one selection. A form that runs no test
// satisfies nothing. Only the root-conformance
// check satisfies its element. A package element needs an unfiltered package run that is
// not `--changed`, whose target equals the element or covers it at a segment boundary.
func satisfies(selection testreport.Request, element string) bool {
	if !selection.Runs() {
		return false
	}
	if element == testreport.RootConformanceCheckName {
		return selection.Form() == testreport.CheckForm && selection.Target() == element
	}
	if selection.Form() != testreport.PackageForm || selection.Run() != testreport.AllTests || selection.Changed() {
		return false
	}
	// The parent `.` of `./...` covers every `./` element, the root element `./.` included.
	target := selection.Target()
	parent, recursive := strings.CutSuffix(target, "/...")
	switch {
	case target == element:
		return true
	case !recursive:
		return false
	}
	return tickets.Covers(parent, element)
}

// floorCommand is the exact command that satisfies one element.
func floorCommand(element string) string {
	if element == testreport.RootConformanceCheckName {
		return testCommand + " --check " + element
	}
	return testCommand + " --package " + element
}
