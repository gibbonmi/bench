package reviewrecord

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// AssignmentCall is the caller's part of one ticket dispatch. Source is a full commit ID
// that must be HEAD, and Ticket is a planned ticket basename. Assignment is the new
// occurrence, and RecordAssignment sets its Source. RunID, OrchestratorSession,
// AuthorLimit, and ReviewMode assert the plan's execution declaration; they never
// change it, and an empty ReviewMode asserts the default mode. A non-nil BeforeWrite reads
// the graded plan before the plan write, and its refusal leaves the plan unchanged.
type AssignmentCall struct {
	Source, Ticket                         string
	Assignment                             Assignment
	RunID, OrchestratorSession, ReviewMode string
	AuthorLimit                            int
	BeforeWrite                            func(Plan) error
}

// RecordAssignment appends one assignment to the history of a ticket in the version 2
// completion plan of spec. The working spec and every planned ticket must hold their
// bytes at the source. The complete candidate plan passes the reader's grade before the
// write, and only the plan payload of the spec changes. The answer is the history as
// written.
func RecordAssignment(root, spec string, call AssignmentCall) ([]Assignment, error) {
	if _, err := RecordPath(spec); err != nil {
		return nil, err
	}
	if head, err := benchgit.ResolveCommit(root, "HEAD", "--quiet"); err != nil || head != call.Source {
		return nil, fmt.Errorf("source %s is not HEAD %s; a dispatch writes only the current plan", call.Source, head)
	}
	tree, err := sourceTree(root, call.Source)
	if err != nil {
		return nil, err
	}
	read := func(name string) ([]byte, error) { return unchanged(root, tree, name) }
	data, err := read(spec)
	if err != nil {
		return nil, err
	}
	plan, err := parsePlan(data)
	if err != nil {
		return nil, err
	}
	if !plan.Delegated() {
		return nil, fmt.Errorf("the plan at %s has version %d; a dispatch needs a version 2 plan with an execution declaration, and version 1 conversion is an authoring step", call.Source, plan.Version)
	}
	execution := plan.Execution
	for _, field := range [][3]string{{"run", call.RunID, execution.RunID}, {"orchestrator", call.OrchestratorSession, execution.OrchestratorSession},
		{"author limit", strconv.Itoa(call.AuthorLimit), strconv.Itoa(execution.AuthorLimit)}, {"review mode", call.ReviewMode, execution.ReviewMode}} {
		if field[1] != field[2] {
			return nil, fmt.Errorf("%s %q differs from the plan's %q; a dispatch asserts the execution declaration and never changes it", field[0], field[1], field[2])
		}
	}
	history, planned := execution.Assignments[call.Ticket]
	if !planned {
		return nil, fmt.Errorf("ticket %s has no assignment history in the plan; name a planned ticket", call.Ticket)
	}
	item := call.Assignment
	item.Source = call.Source
	history = append(history, item)
	execution.Assignments[call.Ticket] = history
	encoded, err := payload(plan, "")
	if err != nil {
		return nil, err
	}
	span, err := locate(data, planFence)
	if err != nil {
		return nil, err
	}
	candidate := span.splice(data, encoded)
	if _, err := checkPlan(spec, candidate, read); err != nil {
		return nil, err
	}
	for _, name := range append([]string{spec}, plan.ticketPaths(spec)...) {
		if _, err := read(name); err != nil {
			return nil, err
		}
	}
	if call.BeforeWrite != nil {
		if err := call.BeforeWrite(plan); err != nil {
			return nil, err
		}
	}
	return history, replace(filepath.Join(root, filepath.FromSlash(spec)), candidate)
}

// ticketPaths is the repository path of every planned ticket beside spec.
func (p Plan) ticketPaths(spec string) []string {
	paths := []string{}
	for _, name := range p.ticketNames() {
		paths = append(paths, ticketPath(spec, name))
	}
	return paths
}

// unchanged returns the bytes of the file at name in tree when the working file holds
// the same bytes, so a write never replaces an input that differs from its source.
func unchanged(root, tree, name string) ([]byte, error) {
	source, err := benchgit.ReadTreeFile(root, tree, name)
	if err != nil {
		return nil, err
	}
	working, err := readFile(root, name)
	if err != nil {
		return nil, fmt.Errorf("working %s: %w", name, err)
	}
	if !bytes.Equal(working, source) {
		return nil, fmt.Errorf("working %s differs from its source; commit or discard the change first", name)
	}
	return source, nil
}
