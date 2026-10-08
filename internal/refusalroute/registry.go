// Package refusalroute owns the refusal faces of the write verbs: the one ordered face
// inventory, the route each face names, and the one constructor a refusal travels
// through. It imports no write-verb package, so each write verb imports it with no cycle.
package refusalroute

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

// NextField is the one record label a write verb prints a refusal's route under.
const NextField = "next"

// Verb names the write verb whose rule raises a face. Another write verb can print the
// face, so the verb states the owner of the rule and not the printer.
type Verb string

const (
	Commit     Verb = "commit"
	Merge      Verb = "merge"
	Reset      Verb = "reset"
	Land       Verb = "land"
	Gate       Verb = "gate"
	Commitment Verb = "commitment"
)

// Authority names who may run a face's route. An agent route clears its cause with Bench
// verbs and file edits inside the agent's own worktree. Any other clear is the
// reviewer's, and its route renders behind the reviewer marker.
type Authority string

const (
	Agent    Authority = "agent"
	Reviewer Authority = "reviewer"
)

// Face is one registered refusal. Sentence is empty when a policy owns the sentence at
// the refusal, and the raising site then supplies the observed one.
type Face struct {
	Verb      Verb
	Name      string
	Sentence  string
	Authority Authority
	Route     []Step
}

// Facts are the observed values a raising site supplies to the constructor. An observed
// Sentence replaces the face's declared one, because it carries the cause the declared
// one drops. Preface is a sentence the route states ahead of its first step. Values holds
// the named facts that the route's slots read.
type Facts struct {
	Sentence string
	Paths    []string
	Preface  string
	Values   map[string]string
}

// Refusal is the typed refusal the constructor returns. Route is the face's route,
// rendered over the raising site's facts.
type Refusal struct {
	Face     Face
	Sentence string
	Paths    []string
	Route    string
}

func (r Refusal) Error() string { return r.Sentence }

// Raised is a refusal that a policy raises for whichever write verb runs the policy: the
// name of its registered face, the cause, and the values the policy observed. The policy
// does not know the verb that prints it, so it renders no route; Printed renders it.
type Raised struct {
	Name   string
	Err    error
	Values map[string]string
}

func (r Raised) Error() string { return r.Err.Error() }
func (r Raised) Unwrap() error { return r.Err }

// Printed is the refusal that a verb prints for cause: the face that cause raised, or else
// face. values are the printing verb's own facts, and the raising policy's values join
// them. The sentence is the cause's whole message, so the context that wraps a raised
// refusal stays; a nil cause keeps the face's declared sentence.
func Printed(face string, cause error, values map[string]string) Refusal {
	facts := Facts{Values: map[string]string{}}
	maps.Copy(facts.Values, values)
	if cause != nil {
		facts.Sentence = cause.Error()
	}
	var raised Raised
	if errors.As(cause, &raised) {
		face = raised.Name
		maps.Copy(facts.Values, raised.Values)
	}
	return New(face, facts)
}

// The named facts the registered routes read. A raising site keys Facts.Values with these
// names, so a slot and the value that fills it share one spelling.
const (
	// FactRerun is the caller's own command, composed from the flags the caller passed.
	FactRerun = "rerun"
	// FactLabel is the label of the assignment that owns the refusing tree.
	FactLabel = "label"
	// FactSiblingLabel is the label of the sibling assignment whose tree a merge reads.
	FactSiblingLabel = "sibling-label"
	// FactFrom is the incoming commit of a merge, as the caller spelled it.
	FactFrom = "from"
	// FactResetPlan is a reset verb's plan command: the plan at the commit a verb published,
	// or the plan of the checkout as it stands.
	FactResetPlan = "reset-plan"
	// FactRecovery is the route that clears an assignment record whose tree is missing,
	// which the record's own landedness picks.
	FactRecovery = "recovery"
	// FactCheckoutGit is the Git command addressed at the conflicted checkout.
	FactCheckoutGit = "checkout-git"
	// FactConflictCommit is the commit whose composition conflicted.
	FactConflictCommit = "full-destination-commit"
	// FactAssignmentID is the id of the assignment that owns the checkout a step addresses.
	FactAssignmentID = "id"
	// FactCheckoutCommand is a command the raising site composed that takes a checkout
	// path as its last word.
	FactCheckoutCommand = "command"
	// FactCheckout is the checkout path. A path that is not line-safe prints the
	// placeholder, and a route that can look the path up names the lookup first.
	FactCheckout = "checkout"
	// FactPublishedCommit is the commit a verb published before its checkout failed to
	// follow it.
	FactPublishedCommit = "published-commit"
	// FactArguments is the caller's own arguments after a tree target, each rendered by Arg.
	FactArguments = "arguments"
	// FactSlug is the slug of the spec whose evidence a route reads.
	FactSlug = "slug"
	// FactTicket is the path of the light-path ticket whose Writes line a route edits.
	FactTicket = "ticket"
	// FactMilestone is the id of the milestone whose verification a route records. Its
	// placeholder is the commitment grammar's <id>, and no route reads it beside
	// FactAssignmentID.
	FactMilestone = "id"
)

// Arg renders one value of a command that a raising site composes: shell-quoted, or the
// placeholder of slot when the value is absent or not line-safe. A composed command
// renders its values here, so it and a route's own slots follow one rule.
func Arg(slot, value string) string { return quoted(slot, value) }

// rerun ends a route with the caller's own command, so one paste finishes the recovery.
var rerun = Command(Composed(FactRerun))

// commitAt commits the operator's repair in the worktree that the named label fact
// addresses. Every face whose repair is the agent's own commit names this one step.
func commitAt(label string) Step {
	return TreeCommand("bench commit", Fact(label), Text("-m"), Operator("msg"), Text("--"), Operators("path"))
}

// doctor diagnoses an infrastructure outcome ahead of the caller's re-run.
var doctor = Command(Text("bench doctor"))

// clearCause is the reviewer's step of each verb's handback face: a cause that only the
// reviewer clears and that has no reviewer step of its own. The refusal's own sentence
// names the cause.
var clearCause = Instruction(Text("clear the cause that the refusal names"))

// handback is the route of each handback face whose verb re-runs as the caller wrote it.
var handback = []Step{clearCause, rerun}

// handMerge is the hand repair a composition conflict demands, up to the commit that
// records the resolution. The landing and the merge refuse the same conflict, so both name
// this one repair. Neither verb composes the repair itself, so the merge step is raw Git
// and the step says so.
var handMerge = []Step{
	Command(Composed(FactCheckoutGit), Text("merge"), Fact(FactConflictCommit), Text("(bench worktree merge refuses this conflict; resolve it by hand)")),
	Command(Text("bench commit")),
}

// pathLookup prints the checkout path of an assignment, which the operator puts in place
// of the checkout placeholder.
var pathLookup = Command(Text("bench worktree path"), Fact(FactAssignmentID))

// atCheckout is a composed command addressed at a checkout through the checkout slot.
var atCheckout = Command(Composed(FactCheckoutCommand), Fact(FactCheckout))

// inventory is the authoritative, ordered inventory of the write verbs' refusal faces. Each
// verb declares its faces in its own file of this package, and a verb adds its face there
// rather than composing a route at the site that refuses. The verbs join in this order.
var inventory = slices.Concat(landFaces, mergeFaces, resetFaces, commitFaces, gateFaces, commitmentFaces)

// New is the one constructor a registered face travels through.
func New(name string, facts Facts) Refusal { return newIn(inventory, name, facts) }

// Inventory returns every registered face, in registry order. A check over all the faces
// reads them here, so the check and the registry cannot drift.
func Inventory() []Face { return slices.Clone(inventory) }

// Faces returns the registered faces of one verb, in registry order. A walk over a verb's
// faces reads them here, so the walk and the registry cannot drift.
func Faces(verb Verb) []Face {
	var faces []Face
	for _, face := range inventory {
		if face.Verb == verb {
			faces = append(faces, face)
		}
	}
	return faces
}

// Sentence is the sentence a registered face declares, and the empty sentence for a face
// whose sentence a policy owns or for a name outside the registry.
func Sentence(name string) string {
	face, _ := faceNamed(inventory, name)
	return face.Sentence
}

// AtCheckout renders a command whose last word is a checkout path that is not line-safe.
// No quoting makes a control byte pasteable, and `bench worktree exec` refuses a Bench
// child, so the route looks the path up by the assignment id and runs the command with
// the checkout placeholder. With no id there is nothing to look up, and the
// command stands alone.
func AtCheckout(command, id string) string {
	route := []Step{atCheckout}
	if id != "" {
		route = []Step{pathLookup, atCheckout}
	}
	return Face{Route: route}.Render(Facts{Values: map[string]string{FactCheckoutCommand: command, FactAssignmentID: id}})
}

// newIn is the constructor over an injectable inventory, so a test can supply its own faces.
func newIn(faces []Face, name string, facts Facts) Refusal {
	face, ok := faceNamed(faces, name)
	if !ok {
		// A name outside the registry is a programming fault, not an operator condition.
		// It refuses rather than panics, because a write verb must not abort the
		// operator's session over its own bookkeeping. The fault's own sentence replaces
		// any observed one, so the fault cannot hide behind the cause it was to name.
		face = Face{
			Name:      name,
			Sentence:  "refusal face " + name + " is unregistered",
			Authority: Reviewer,
			Route:     []Step{Instruction(Text("declare this refusal face in the refusal-route registry"))},
		}
		facts.Sentence = ""
	}
	sentence := facts.Sentence
	if sentence == "" {
		sentence = face.Sentence
	}
	return Refusal{Face: face, Sentence: sentence, Paths: facts.Paths, Route: face.Render(facts)}
}

// faceNamed finds the face that declares name. The uniqueness walk keeps the answer to
// one face.
func faceNamed(faces []Face, name string) (Face, bool) {
	for _, face := range faces {
		if face.Name == name {
			return face, true
		}
	}
	return Face{}, false
}

// uniqueNames refuses a face list in which two faces declare one name, because the
// lookup would then return the first and hide the second.
func uniqueNames(faces []Face) error {
	seen := make(map[string]bool, len(faces))
	for _, face := range faces {
		if seen[face.Name] {
			return fmt.Errorf("refusal face %q is declared twice", face.Name)
		}
		seen[face.Name] = true
	}
	return nil
}
