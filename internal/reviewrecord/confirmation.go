package reviewrecord

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/testgap"
)

// ConfirmationRole is the evidence role of a coordinator confirmation.
const ConfirmationRole = "coordinator-confirmation"

// Confirmation is the coordinator's record that a test-only repair of a chunk closes the
// findings of the axis results at an earlier pair. Base and ReviewedTip are that reviewed
// pair, and RepairedTip is the chunk tip after the repair; the evidence source digest is
// the repaired digest. FindingIDs lists the findings that the repair closes, Supersedes
// names the prior confirmation of the chunk, and Probe is the coordinator's probe of the
// repair.
type Confirmation struct {
	Evidence
	Base        string   `json:"base"`
	ReviewedTip string   `json:"reviewed_tip"`
	RepairedTip string   `json:"repaired_tip"`
	FindingIDs  []string `json:"finding_ids"`
	Supersedes  []string `json:"supersedes"`
	Probe       Probe    `json:"probe"`
}

// validateConfirmations grades the confirmation list of one chunk entry. Each confirmation
// after the first supersedes the one before it.
func validateConfirmations(items []Confirmation) error {
	prior := ""
	for _, item := range items {
		if err := validateEvidence(item.Evidence); err != nil {
			return err
		}
		if item.Role != ConfirmationRole {
			return fmt.Errorf("%s: invalid confirmation role %q; use %s", item.ID, item.Role, ConfirmationRole)
		}
		if !objectID.MatchString(item.Base) || !objectID.MatchString(item.ReviewedTip) || !objectID.MatchString(item.RepairedTip) {
			return fmt.Errorf("%s: invalid confirmation pair", item.ID)
		}
		if !supersedes(prior, item.Supersedes) {
			return fmt.Errorf("%s: invalid supersession; retain and name the preceding confirmation", item.ID)
		}
		if err := validateProbe(item.ID, item.Probe); err != nil {
			return err
		}
		prior = item.ID
	}
	return nil
}

// ReviewedPair is the earlier frozen pair of a chunk that holds its three current axis
// results: the chunk base and Tip. Digest is the source digest of the Tip tree, and Findings
// is the union of the finding IDs of the three results, sorted.
type ReviewedPair struct {
	Tip, Digest string
	Findings    []string
}

// errAtChunkPair marks three current axis results that sit at the chunk pair, so no earlier
// pair holds them.
var errAtChunkPair = errors.New("the axis results sit at the chunk pair")

// FindReviewedPair derives the reviewed pair of chunk from the current result of each axis.
// spec names the record that the source digest excludes. It refuses a missing axis, results
// at different pairs or at the chunk pair, a reviewed tip equal to the chunk base, a
// reviewed tip outside the line from the chunk base to the chunk tip, and a delta from the
// reviewed digest to the chunk digest that testgap.Prove refuses. Each refusal names the
// chunk.
func FindReviewedPair(root, spec string, chunk Chunk) (ReviewedPair, error) {
	scope := chunk.reviewed().scope
	tip, mixed, pairs, findings := "", false, []string{}, []string{}
	for _, axis := range Axes() {
		current := currentReview(chunk.Reviews, axis)
		if current == nil {
			return ReviewedPair{}, fmt.Errorf("%s: missing %s; a reviewed pair needs a current result of each axis", scope, axis)
		}
		if tip == "" {
			tip = current.Tip
		}
		mixed = mixed || current.Base != chunk.Base || current.Tip != tip
		pairs = append(pairs, axis+" at "+current.Base+".."+current.Tip)
		findings = append(findings, current.FindingIDs...)
	}
	if mixed {
		return ReviewedPair{}, fmt.Errorf("%s: mixed review pairs: %s; record the three axes at one pair from the chunk base %s", scope, strings.Join(pairs, ", "), chunk.Base)
	}
	if tip == chunk.Tip {
		return ReviewedPair{}, fmt.Errorf("%s: %w", scope, errAtChunkPair)
	}
	// The ancestry query is reflexive, so it accepts a reviewed tip equal to the chunk base.
	if tip == chunk.Base {
		return ReviewedPair{}, fmt.Errorf("%s: empty reviewed pair %s..%s; review the three axes at a tip after the chunk base", scope, chunk.Base, tip)
	}
	descends, err := isAncestor(root, chunk.Base, tip)
	if err != nil {
		return ReviewedPair{}, fmt.Errorf("%s: %w", scope, err)
	}
	if !descends {
		return ReviewedPair{}, fmt.Errorf("%s: reviewed tip %s does not descend from the chunk base %s", scope, tip, chunk.Base)
	}
	precedes, err := isAncestor(root, tip, chunk.Tip)
	if err != nil {
		return ReviewedPair{}, fmt.Errorf("%s: %w", scope, err)
	}
	if !precedes {
		return ReviewedPair{}, fmt.Errorf("%s: reviewed tip %s is not an ancestor of the chunk tip %s", scope, tip, chunk.Tip)
	}
	tree, err := sourceTree(root, tip)
	if err != nil {
		return ReviewedPair{}, fmt.Errorf("%s: %w", scope, err)
	}
	digest, err := SourceDigest(root, tree, spec)
	if err != nil {
		return ReviewedPair{}, fmt.Errorf("%s: %w", scope, err)
	}
	if err := testgap.Prove(root, digest, chunk.SourceDigest); err != nil {
		return ReviewedPair{}, fmt.Errorf("%s: stale axis results at the reviewed tip %s: the delta to the chunk tip is not test-only: %w; review the three axes at the chunk pair", scope, tip, err)
	}
	slices.Sort(findings)
	return ReviewedPair{Tip: tip, Digest: digest, Findings: slices.Compact(findings)}, nil
}

// checkChunkReviews grades the reviews of chunk. With the three current axis results at the
// chunk pair, or with no reviewed pair and no confirmation, the strict rule applies. A
// reviewed pair with no confirmation refuses through the confirmation form, before the strict
// rule grades an outcome. Otherwise the latest confirmation must hold, and the results carry
// from the reviewed pair. A refusal of the reviewed pair repairs through the review form, and a
// refusal of the confirmation through the confirmation form. plan is the frozen plan of chunk.
func checkChunkReviews(root string, record Record, plan Plan, chunk *Chunk, excluded []string, reviewerCount int) error {
	pair := chunk.reviewed()
	reviewed, err := FindReviewedPair(root, record.Spec, *chunk)
	if errors.Is(err, errAtChunkPair) || err != nil && len(chunk.Confirmations) == 0 {
		return checkReviews(pair, excluded, reviewerCount)
	}
	refused := func(form string, err error) error {
		repair := pair.at
		repair.Form, repair.Err = form, err
		return repair
	}
	if err != nil {
		return refused(FormReview, err)
	}
	if len(chunk.Confirmations) == 0 {
		return refused(FormConfirmation, fmt.Errorf("%s: stale %s at the reviewed tip %s; the delta to the chunk tip is test-only, so confirm the repair or review the three axes at the chunk pair", pair.scope, strings.Join(Axes(), ", "), reviewed.Tip))
	}
	if err := confirms(*chunk, reviewed, completionReconciler(record, plan)); err != nil {
		return refused(FormConfirmation, err)
	}
	pair.tip, pair.digest, pair.carried = reviewed.Tip, reviewed.Digest, true
	return checkReviews(pair, excluded, reviewerCount)
}

// confirms grades the latest confirmation of chunk against reviewed, the reviewed pair of its
// axis results. reconciler is the one performer that confirms. Each refusal names the chunk
// and the failed condition.
func confirms(chunk Chunk, reviewed ReviewedPair, reconciler string) error {
	latest := chunk.Confirmations[len(chunk.Confirmations)-1]
	scope := chunk.confirmationScope(latest.ID)
	switch {
	case latest.Base != chunk.Base || latest.RepairedTip != chunk.Tip || latest.SourceDigest != chunk.SourceDigest:
		return fmt.Errorf("%s covers the stale pair %s..%s; confirm the chunk pair %s..%s", scope, latest.Base, latest.RepairedTip, chunk.Base, chunk.Tip)
	case latest.Performer != reconciler:
		return fmt.Errorf("%s performer %s is not the completion reconciler %s", scope, latest.Performer, reconciler)
	case latest.State != "completed":
		return fmt.Errorf("%s state is %s; record a completed confirmation", scope, latest.State)
	}
	if err := probeHolds(scope, latest.Probe); err != nil {
		return err
	}
	if latest.ReviewedTip != reviewed.Tip {
		return fmt.Errorf("%s names the reviewed tip %s, and the axis results sit at %s", scope, latest.ReviewedTip, reviewed.Tip)
	}
	return closesOpen(scope, latest.FindingIDs, reviewed)
}

// confirmationScope names the confirmation id of c in a refusal.
func (c *Chunk) confirmationScope(id string) string {
	return c.reviewed().scope + ": confirmation " + id
}

// probeHolds refuses a confirmation probe that did not bite or did not restore. scope names the
// confirmation that holds the probe.
func probeHolds(scope string, probe Probe) error {
	bit, restored := probeBites(probe)
	if !bit {
		return fmt.Errorf("%s probe did not bite: outcome %s, focused exit %d; record a probe that bites", scope, probe.Outcome, probe.ExitCode)
	}
	if !restored {
		return fmt.Errorf("%s probe restore %s; restore and verify the source", scope, probe.Restore)
	}
	return nil
}

// closesOpen refuses findings that differ from the open findings at reviewed. scope names the
// confirmation that closes them.
func closesOpen(scope string, findings []string, reviewed ReviewedPair) error {
	if sameSet(findings, reviewed.Findings) {
		return nil
	}
	return fmt.Errorf("%s closes the findings [%s], and the open findings at the reviewed pair are [%s]; list each open finding", scope, strings.Join(findings, " "), strings.Join(reviewed.Findings, " "))
}
