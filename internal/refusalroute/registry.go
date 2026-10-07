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

// The named facts the registered routes read. A raising site keys Facts.Values with these
// names, so a slot and the value that fills it share one spelling.
const (
	// FactRerun is the caller's own command, composed from the flags the caller passed.
	FactRerun = "rerun"
	// FactLabel is the label of the assignment that owns the refusing tree.
	FactLabel = "label"
	// FactCheckoutGit is the Git command addressed at the conflicted checkout.
	FactCheckoutGit = "checkout-git"
	// FactConflictCommit is the commit whose composition conflicted.
	FactConflictCommit = "full-destination-commit"
)

// rerun ends a route with the caller's own command, so one paste finishes the recovery.
var rerun = Command(Composed(FactRerun))

// destinationClean is the repair a destination that carries uncommitted work demands. The
// first run and the resume refuse on the same destination state, so both faces name it.
var destinationClean = Instruction(Text("commit the destination's uncommitted work, or discard it"))

// review sends a repaired source back through review before the landing re-runs.
var review = Command(Text("/bench-review-implementation"))

// handMerge is the hand repair a composition conflict demands, up to the commit that
// records the resolution. The landing and the merge refuse the same conflict, so both name
// this one repair. Neither verb composes the repair itself, so the merge step is raw Git
// and the step says so.
var handMerge = []Step{
	Command(Composed(FactCheckoutGit), Text("merge"), Fact(FactConflictCommit), Text("(bench worktree merge refuses this conflict; resolve it by hand)")),
	Command(Text("bench commit")),
}

// inventory is the authoritative, ordered inventory of the write verbs' refusal faces. A
// verb adds its face here rather than composing a route at the site that refuses.
var inventory = []Face{
	{
		Verb:      Land,
		Name:      "destination-not-clean",
		Sentence:  "landing destination is not clean",
		Authority: Reviewer,
		Route:     []Step{destinationClean, rerun},
	},
	{
		// The paths are the operator's own files where the landing writes. Git refuses to
		// overwrite an untracked one and overwrites an ignored one without a word, so the
		// operator moves them before the landing runs.
		Verb:      Land,
		Name:      "destination-collision",
		Sentence:  "landing destination has untracked or ignored files where the landing writes",
		Authority: Reviewer,
		Route:     []Step{Instruction(Text("move the refusal_paths entries out of the landing checkout")), rerun},
	},
	{
		// The raising site re-points the caller's re-run at the source tip the tree holds,
		// so a moved tip leaves exactly one command to run.
		Verb:      Land,
		Name:      "source-tip-mismatch",
		Sentence:  "worktree source tip mismatch",
		Authority: Agent,
		Route:     []Step{rerun},
	},
	{
		Verb:      Land,
		Name:      "source-not-clean",
		Sentence:  "reviewed source is not clean",
		Authority: Agent,
		Route: []Step{
			TreeCommand("bench commit", Fact(FactLabel), Text("-m"), Operator("msg"), Text("--"), Operators("path")),
			rerun,
		},
	},
	{
		// The fence already authorizes a path that a fold of the default branch brings in
		// unchanged, so every refused path is a write of the build, and the route keeps the
		// caller's --base.
		Verb:      Land,
		Name:      "source-not-fenced",
		Sentence:  "reviewed source range or ownership fence is invalid",
		Authority: Agent,
		Route: []Step{
			Instruction(Text("take the refusal_paths entries out of the reviewed range, or declare them under the spec's ## Ownership fences")),
			rerun,
		},
	},
	{
		// The composition names the conflicted paths in its own sentence, so the face
		// declares none.
		Verb:      Land,
		Name:      "composition-conflict",
		Authority: Reviewer,
		Route:     append(append([]Step{}, handMerge...), review, rerun),
	},
	{
		// The source already holds a merge in progress, so a second merge is not the
		// repair: finishing the merge records the resolution and needs no separate commit.
		Verb:      Land,
		Name:      "composition-conflict-pending",
		Authority: Reviewer,
		Route: []Step{
			Command(Composed(FactCheckoutGit), Text("merge --continue (resolve the conflicted paths of the merge in progress first)")),
			review,
			rerun,
		},
	},
	{
		// The residue policy owns the sentence this face prints.
		Verb:      Land,
		Name:      "resume-destination-residue",
		Authority: Reviewer,
		Route:     []Step{destinationClean, rerun},
	},
	{
		// The landing policy owns the sentence this face prints. The marker belongs to the
		// published landing on the primary checkout, so the gate runs there.
		Verb:      Land,
		Name:      "resume-marker",
		Authority: Agent,
		Route:     []Step{TreeCommand("bench gate", Text("primary")), rerun},
	},
	{
		// The landing published, and a later step did not finish. The re-run is the
		// caller's resume of that published landing.
		Verb:      Land,
		Name:      "land-incomplete",
		Authority: Agent,
		Route:     []Step{rerun},
	},
	{
		// A landing cause that only the reviewer clears and that has no reviewer step of
		// its own. The refusal's own sentence names the cause.
		Verb:      Land,
		Name:      "land-handback",
		Authority: Reviewer,
		Route:     []Step{Instruction(Text("clear the cause that the refusal names")), rerun},
	},
}

// New is the one constructor a registered face travels through.
func New(name string, facts Facts) Refusal { return newIn(inventory, name, facts) }

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

// HandMerge renders the hand merge of a composition conflict over the raising site's facts.
// The merge verb prints it as the repair of its own conflict refusal.
func HandMerge(facts Facts) string { return Face{Route: handMerge}.Render(facts) }

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
