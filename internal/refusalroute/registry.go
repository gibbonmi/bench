// Package refusalroute owns the refusal faces of the write verbs: the one ordered face
// inventory, the route each face names, and the one constructor a refusal travels
// through. It imports no write-verb package, so each write verb imports it with no cycle.
package refusalroute

import "fmt"

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

// inventory is the authoritative, ordered inventory of the write verbs' refusal faces. A
// verb adds its face here rather than composing a route at the site that refuses.
var inventory = []Face{}

// New is the one constructor a registered face travels through.
func New(name string, facts Facts) Refusal { return newIn(inventory, name, facts) }

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
