// Package reviewrecord reads source-bound verification and review evidence.
package reviewrecord

import (
	"fmt"
	"slices"
)

type NativeRef struct {
	Ref     string `json:"ref"`
	Digest  string `json:"digest"`
	Excerpt string `json:"excerpt"`
}
type Evidence struct {
	ID           string    `json:"id"`
	Performer    string    `json:"performer"`
	Role         string    `json:"role"`
	Model        string    `json:"model"`
	Effort       string    `json:"effort"`
	SourceDigest string    `json:"source_digest"`
	State        string    `json:"state"`
	Outcome      string    `json:"outcome"`
	NativeRef    NativeRef `json:"native_ref"`
}
type Probe struct {
	Mutation  string    `json:"mutation"`
	Outcome   string    `json:"outcome"`
	ExitCode  int       `json:"exit_code"`
	Restore   string    `json:"restore"`
	NativeRef NativeRef `json:"native_ref"`
}
type Verification struct {
	Evidence
	Requirement string `json:"requirement"`
	Command     string `json:"command"`
	ExitCode    *int   `json:"exit_code"`
	Probe       *Probe `json:"probe,omitempty"`
}
type Review struct {
	Evidence
	Axis       string   `json:"axis"`
	Base       string   `json:"base"`
	Tip        string   `json:"tip"`
	FindingIDs []string `json:"finding_ids"`
	Supersedes []string `json:"supersedes"`
}
type Chunk struct {
	ID             string         `json:"id"`
	Base           string         `json:"base"`
	Tip            string         `json:"tip"`
	PlanDigest     string         `json:"plan_digest"`
	SourceDigest   string         `json:"source_digest"`
	AcceptanceRows []string       `json:"acceptance_rows"`
	Verification   []Verification `json:"verification"`
	Reviews        []Review       `json:"reviews"`
	Confirmations  []Confirmation `json:"confirmations,omitempty"`
}
type Completion struct {
	State          string            `json:"state"`
	SourceDigest   string            `json:"source_digest"`
	Performer      string            `json:"performer"`
	Reconciliation map[string]string `json:"reconciliation"`
	Verification   []Verification    `json:"verification"`
}
type Amendment struct {
	From     string              `json:"from"`
	To       string              `json:"to"`
	ChunkIDs map[string][]string `json:"chunk_ids"`
}

// Fold is one recorded merge of a default-branch commit into the build source. After
// names the chunk that the fold follows. From is the first parent, Commit the merge, and
// Incoming the default-branch commit. Resolution lists, in byte order, each path where
// Commit differs from the automatic merge of From and Incoming. TestPackages lists the
// packages of the test files that differ between From and Commit.
type Fold struct {
	ID               string   `json:"id"`
	After            string   `json:"after"`
	From             string   `json:"from"`
	Commit           string   `json:"commit"`
	Incoming         string   `json:"incoming"`
	SourceDigest     string   `json:"source_digest"`
	Resolution       []string `json:"resolution"`
	ResolutionDigest string   `json:"resolution_digest"`
	TestPackages     []string `json:"test_packages"`
	Reviews          []Review `json:"reviews"`
}
type Record struct {
	Version               int         `json:"version"`
	Spec                  string      `json:"spec"`
	PlanDigest            string      `json:"plan_digest"`
	ImplementationSession string      `json:"implementation_session"`
	Chunks                []Chunk     `json:"chunks"`
	Completion            Completion  `json:"completion"`
	Amendments            []Amendment `json:"amendments,omitempty"`
	Folds                 []Fold      `json:"folds,omitempty"`
}

// CheckReviews grades a chunk's three axes. The excluded set holds every
// identity that cannot supply independent review. Under distinct, one session
// supplies at most one axis for the chunk.
func CheckReviews(chunk Chunk, excluded []string, distinct bool) error {
	reviewerCount := 0
	if distinct {
		reviewerCount = 3
	}
	return checkReviews(chunk.reviewed(), excluded, reviewerCount)
}

// reviewed is one frozen pair that takes independent review: a chunk entry or a fold entry.
// scope names the pair in a refusal, at holds its repair target, and reviews is its list.
// carried marks the reviewed pair of a valid confirmation: a result there may hold open
// findings, because the confirmation closes them.
type reviewed struct {
	scope             string
	at                Repair
	base, tip, digest string
	reviews           *[]Review
	carried           bool
}

// reviewed is the frozen pair of a chunk entry.
func (c *Chunk) reviewed() reviewed {
	return reviewed{scope: "chunk " + c.ID, at: Repair{Chunk: c.ID}, base: c.Base, tip: c.Tip, digest: c.SourceDigest, reviews: &c.Reviews}
}

// reviewed is the pair of a fold: its first parent and the fold commit.
func (f *Fold) reviewed() reviewed {
	return reviewed{scope: "fold " + f.ID, at: Repair{Fold: f.ID}, base: f.From, tip: f.Commit, digest: f.SourceDigest, reviews: &f.Reviews}
}

func checkReviews(pair reviewed, excluded []string, reviewerCount int) error {
	// Each refusal below repairs through the review form: the pair lacks a current
	// independent result of the named axis, or of the axis set.
	refused := func(axis string, err error) error {
		repair := pair.at
		repair.Form, repair.Axis, repair.Err = FormReview, axis, err
		return repair
	}
	axisOwners := map[string]string{}
	for _, axis := range Axes() {
		current := currentReview(*pair.reviews, axis)
		if current == nil {
			return refused(axis, fmt.Errorf("%s: missing %s; record the native review result", pair.scope, axis))
		}
		if current.Role != "independent-review" || slices.Contains(excluded, current.Performer) || current.Performer == "" {
			return refused(axis, fmt.Errorf("%s: invalid %s performer or role; obtain independent review", pair.scope, axis))
		}
		axisOwners[current.Performer] = axis
		if current.State != "completed" || (!pair.carried && (current.Outcome != "pass" || len(current.FindingIDs) != 0)) {
			return refused(axis, fmt.Errorf("%s: %s %s; resolve findings and obtain a completed result", pair.scope, current.State, axis))
		}
		if current.SourceDigest != pair.digest || current.Base != pair.base || current.Tip != pair.tip {
			return refused(axis, fmt.Errorf("%s: stale %s; review the current frozen pair", pair.scope, axis))
		}
	}
	if reviewerCount == 3 && len(axisOwners) != 3 {
		return refused("", fmt.Errorf("%s: use three distinct review sessions", pair.scope))
	}
	if reviewerCount == 1 && len(axisOwners) != 1 {
		return refused("", fmt.Errorf("%s: unified review requires one reviewer for all three axes", pair.scope))
	}
	return nil
}

// currentReview returns the current result of axis in reviews: the last result of that
// axis, or nil when reviews holds none.
func currentReview(reviews []Review, axis string) *Review {
	var current *Review
	for i := range reviews {
		if reviews[i].Axis == axis {
			current = &reviews[i]
		}
	}
	return current
}

func Axes() []string { return []string{"Standards", "Spec", "Coverage"} }
