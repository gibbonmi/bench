package commitment

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// The terminal criterion result states. Only a verified result can complete a milestone.
const (
	ResultVerified = "verified"
	ResultUnmet    = "unmet"
	ResultBlocked  = "blocked"
)

// MilestoneEvidence is the verification input for one milestone. It gives one result for
// each criterion at one published revision, and names that revision's retained green gate
// evidence. A green gate is necessary delivery evidence, but it proves no criterion.
type MilestoneEvidence struct {
	Version  int               `json:"version"`
	Revision string            `json:"revision"`
	Gate     string            `json:"gate"`
	Results  []CriterionResult `json:"results"`
}

// CriterionResult is one criterion's result. Identity is the criterion identity that the
// result examined. Evidence is a native reference: a repository path at the examined
// revision. Assessment is the reviewer's explicit assessment of the outcome that the
// evidence shows.
type CriterionResult struct {
	Criterion  string `json:"criterion"`
	Identity   string `json:"identity"`
	Result     string `json:"result"`
	Evidence   string `json:"evidence"`
	Assessment string `json:"assessment"`
}

// Verification is the receipt of one successful milestone verification. It binds the
// milestone, the policy and published revision that it examined, and the repository
// object of each evidence reference. A completion proposal consumes it.
type Verification struct {
	ID        string              `json:"id"`
	Milestone string              `json:"milestone"`
	Policy    string              `json:"policy"`
	Revision  string              `json:"revision"`
	Gate      string              `json:"gate"`
	Results   []VerifiedCriterion `json:"results"`
}

// VerifiedCriterion is one verified result and the object its evidence reference names.
type VerifiedCriterion struct {
	Result CriterionResult `json:"result"`
	Object string          `json:"object"`
}

// Resolver returns the repository object that a native evidence reference names at the
// examined revision.
type Resolver func(reference string) (string, error)

// ParseEvidence decodes one complete milestone evidence document.
func ParseEvidence(data []byte) (MilestoneEvidence, error) {
	var evidence MilestoneEvidence
	if err := jsonfile.DecodeExactDocument(data, &evidence); err != nil {
		return MilestoneEvidence{}, fmt.Errorf("milestone evidence: %w", err)
	}
	if evidence.Version != 1 {
		return MilestoneEvidence{}, fmt.Errorf("milestone evidence: unsupported version %d", evidence.Version)
	}
	return evidence, nil
}

// CriterionIdentity is the immutable identity of criterion: its identifier and statement.
func CriterionIdentity(criterion Criterion) string {
	return identity([]byte(criterion.ID + "\x00" + criterion.Text))
}

// VerifyMilestone checks evidence for the active milestone of policy, whose canonical
// identity is policyIdentity, at the published revision. Every outcome of the milestone
// must already be delivered. Every criterion needs exactly one verified result for its
// current identity, with a resolvable evidence reference and a reviewer assessment. The
// check validates evidence identity and coverage; it does not infer the truth of prose.
func VerifyMilestone(policy Policy, policyIdentity, revision, milestone string, evidence MilestoneEvidence, resolve Resolver) (Verification, error) {
	if milestone != policy.ActiveMilestone {
		return Verification{}, fmt.Errorf("milestone %q is not the active milestone", milestone)
	}
	if remaining := Remaining(policy); len(remaining) > 0 {
		return Verification{}, fmt.Errorf("milestone %q has undelivered outcomes %s", milestone, strings.Join(remaining, ","))
	}
	if evidence.Revision != revision {
		return Verification{}, fmt.Errorf("evidence revision %q is not the published revision %s", evidence.Revision, revision)
	}
	gate, err := resolveReference(evidence.Gate, resolve)
	if err != nil {
		return Verification{}, fmt.Errorf("gate evidence: %w", err)
	}
	criteria := milestoneCriteria(policy, milestone)
	results := map[string]VerifiedCriterion{}
	for _, result := range evidence.Results {
		index := slices.IndexFunc(criteria, func(criterion Criterion) bool { return criterion.ID == result.Criterion })
		if index < 0 {
			return Verification{}, fmt.Errorf("criterion %q is unknown to milestone %q", result.Criterion, milestone)
		}
		if _, seen := results[result.Criterion]; seen {
			return Verification{}, fmt.Errorf("criterion %q has a duplicate result", result.Criterion)
		}
		if result.Identity != CriterionIdentity(criteria[index]) {
			return Verification{}, fmt.Errorf("criterion %q result is stale: it examined identity %q, and the criterion is %s", result.Criterion, result.Identity, CriterionIdentity(criteria[index]))
		}
		switch result.Result {
		case ResultVerified:
		case ResultUnmet, ResultBlocked:
			return Verification{}, fmt.Errorf("criterion %q is %s", result.Criterion, result.Result)
		default:
			return Verification{}, fmt.Errorf("criterion %q has unknown result %q", result.Criterion, result.Result)
		}
		object, err := resolveReference(result.Evidence, resolve)
		if err != nil {
			return Verification{}, fmt.Errorf("criterion %q evidence: %w", result.Criterion, err)
		}
		if strings.TrimSpace(result.Assessment) == "" || !sanitize.LineSafe(result.Assessment) {
			return Verification{}, fmt.Errorf("criterion %q has no reviewer outcome assessment", result.Criterion)
		}
		results[result.Criterion] = VerifiedCriterion{Result: result, Object: object}
	}
	verification := Verification{Milestone: milestone, Policy: policyIdentity, Revision: revision, Gate: gate}
	for _, criterion := range criteria {
		result, found := results[criterion.ID]
		if !found {
			return Verification{}, fmt.Errorf("criterion %q has no result", criterion.ID)
		}
		verification.Results = append(verification.Results, result)
	}
	data, err := json.Marshal(verification)
	if err != nil {
		return Verification{}, fmt.Errorf("encode milestone verification: %w", err)
	}
	verification.ID = identity(data)
	return verification, nil
}

// Completes checks that v is the current receipt for completion. It must verify that
// milestone under the predecessor policy, and proposed must keep exactly the criteria that
// it verified, so a completion cannot add, remove, or weaken a criterion.
func (v Verification) Completes(completion Completion, predecessor string, proposed Policy) error {
	if v.ID != completion.Verification || v.Milestone != completion.Milestone {
		return fmt.Errorf("verification %q does not verify milestone %q", completion.Verification, completion.Milestone)
	}
	if v.Policy != predecessor {
		return fmt.Errorf("verification %q is stale: it examined policy %s, and the current policy is %s", v.ID, v.Policy, predecessor)
	}
	var verified, proposedCriteria []string
	for _, result := range v.Results {
		verified = append(verified, result.Result.Identity)
	}
	for _, criterion := range milestoneCriteria(proposed, completion.Milestone) {
		proposedCriteria = append(proposedCriteria, CriterionIdentity(criterion))
	}
	if !slices.Equal(verified, proposedCriteria) {
		return fmt.Errorf("completion of milestone %q changes the criteria that verification %q examined", completion.Milestone, v.ID)
	}
	return nil
}

// Completions returns the milestone completions that proposed adds to current. A proposal
// keeps every published delivery fact and recorded completion. A new completion must
// complete the active milestone, clear it, and activate no successor: activation remains
// its own explicit transition.
func Completions(current *Policy, proposed Policy) ([]Completion, error) {
	var predecessor Policy
	if current != nil {
		predecessor = *current
	}
	if !slices.Equal(predecessor.Deliveries, proposed.Deliveries) {
		return nil, errors.New("a proposal cannot change a published delivery fact")
	}
	for _, completion := range predecessor.Completions {
		if !slices.Contains(proposed.Completions, completion) {
			return nil, fmt.Errorf("a proposal cannot remove the completion of milestone %q", completion.Milestone)
		}
	}
	var added []Completion
	for _, completion := range proposed.Completions {
		if slices.Contains(predecessor.Completions, completion) {
			continue
		}
		if completion.Milestone != predecessor.ActiveMilestone || proposed.ActiveMilestone != "" {
			return nil, fmt.Errorf("completion of milestone %q must clear that active milestone and activate none", completion.Milestone)
		}
		added = append(added, completion)
	}
	return added, nil
}

func findMilestone(policy Policy, id string) (Milestone, bool) {
	index := slices.IndexFunc(policy.Milestones, func(milestone Milestone) bool { return milestone.ID == id })
	if index < 0 {
		return Milestone{}, false
	}
	return policy.Milestones[index], true
}

func milestoneCriteria(policy Policy, id string) []Criterion {
	milestone, _ := findMilestone(policy, id)
	var criteria []Criterion
	for _, outcome := range milestone.Outcomes {
		criteria = append(criteria, outcome.Criteria...)
	}
	return criteria
}

func resolveReference(reference string, resolve Resolver) (string, error) {
	if !repositoryPath(reference) {
		return "", fmt.Errorf("native evidence reference %q is incomplete", reference)
	}
	object, err := resolve(reference)
	if err != nil {
		return "", fmt.Errorf("native evidence reference %q does not resolve: %w", reference, err)
	}
	return object, nil
}
