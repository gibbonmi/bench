// Package refusalroute owns the refusal faces of the write verbs: the one ordered face
// inventory, the route each face names, and the one constructor a refusal travels
// through. It imports no write-verb package, so each write verb imports it with no cycle.
package refusalroute

import (
	"fmt"
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
	// factCheckout is the checkout path. Only a path that is not line-safe takes the
	// lookup route, so its slot prints the placeholder.
	factCheckout = "checkout"
)

// rerun ends a route with the caller's own command, so one paste finishes the recovery.
var rerun = Command(Composed(FactRerun))

// commitAt commits the operator's repair in the worktree that the named label fact
// addresses. Every face whose repair is the agent's own commit names this one step.
func commitAt(label string) Step {
	return TreeCommand("bench commit", Fact(label), Text("-m"), Operator("msg"), Text("--"), Operators("path"))
}

// doctor diagnoses an infrastructure outcome ahead of the caller's re-run.
var doctor = Command(Text("bench doctor"))

// handback is the route of each verb's handback face: a cause that only the reviewer
// clears and that has no reviewer step of its own. The refusal's own sentence names the
// cause.
var handback = []Step{Instruction(Text("clear the cause that the refusal names")), rerun}

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
var atCheckout = Command(Composed(FactCheckoutCommand), Fact(factCheckout))

// inventory is the authoritative, ordered inventory of the write verbs' refusal faces. Each
// verb declares its faces in its own file of this package, and a verb adds its face there
// rather than composing a route at the site that refuses. The verbs join in this order.
var inventory = slices.Concat(landFaces, mergeFaces, resetFaces)

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
