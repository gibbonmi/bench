package commitment

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/sanitize"
)

var identityPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9._-]*$`)

// Parse validates one complete policy document.
func Parse(data []byte) (Policy, error) {
	var policy Policy
	if err := jsonfile.DecodeExactDocument(data, &policy); err != nil {
		return Policy{}, fmt.Errorf("commitment policy: %w", err)
	}
	if err := Validate(policy); err != nil {
		return Policy{}, err
	}
	return policy, nil
}

// Bytes returns the canonical tracked representation of policy.
func Bytes(policy Policy) ([]byte, error) {
	if err := Validate(policy); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode commitment policy: %w", err)
	}
	return append(data, '\n'), nil
}

// Validate checks policy identities and references without reading repository state.
func Validate(policy Policy) error {
	if policy.Version != 1 {
		return fmt.Errorf("commitment policy: unsupported version %d", policy.Version)
	}
	if len(policy.Milestones) == 0 {
		return errors.New("commitment policy: milestones are empty")
	}
	milestones := map[string]bool{}
	outcomes := map[string]bool{}
	criteria := map[string]bool{}
	sources := map[string]string{}
	sourcePaths := map[string]string{}
	dependencies := map[string][]string{}
	for _, milestone := range policy.Milestones {
		if !validIdentity(milestone.ID) || milestones[milestone.ID] {
			return fmt.Errorf("commitment policy: invalid or duplicate milestone %q", milestone.ID)
		}
		milestones[milestone.ID] = true
		criterionCount := 0
		if len(milestone.Outcomes) == 0 {
			return fmt.Errorf("commitment policy: milestone %q has no outcomes", milestone.ID)
		}
		for _, outcome := range milestone.Outcomes {
			if !validIdentity(outcome.ID) || outcomes[outcome.ID] {
				return fmt.Errorf("commitment policy: invalid or duplicate outcome %q", outcome.ID)
			}
			outcomes[outcome.ID] = true
			dependencies[outcome.ID] = append([]string(nil), outcome.Dependencies...)
			criterionCount += len(outcome.Criteria)
			for _, criterion := range outcome.Criteria {
				if !validIdentity(criterion.ID) || criteria[criterion.ID] || strings.TrimSpace(criterion.Text) == "" || !sanitize.LineSafe(criterion.Text) {
					return fmt.Errorf("commitment policy: invalid or duplicate criterion %q", criterion.ID)
				}
				criteria[criterion.ID] = true
			}
			if err := validateDeliverables(outcome); err != nil {
				return err
			}
			for _, source := range outcome.Sources {
				if err := validateSource(source); err != nil {
					return err
				}
				if _, exists := sources[source.ID]; exists || sourcePaths[source.Path] != "" {
					return fmt.Errorf("commitment policy: source %q has more than one owner", source.ID)
				} else {
					sources[source.ID] = outcome.ID
					sourcePaths[source.Path] = outcome.ID
				}
			}
		}
		if milestone.ID == policy.ActiveMilestone && criterionCount == 0 {
			return fmt.Errorf("commitment policy: active milestone %q has no criteria", milestone.ID)
		}
	}
	if policy.ActiveMilestone != "" && !milestones[policy.ActiveMilestone] {
		return fmt.Errorf("commitment policy: active milestone %q is unknown", policy.ActiveMilestone)
	}
	for outcome, deps := range dependencies {
		seen := map[string]bool{}
		for _, dependency := range deps {
			if !outcomes[dependency] || dependency == outcome || seen[dependency] {
				return fmt.Errorf("commitment policy: invalid dependency %q for %q", dependency, outcome)
			}
			seen[dependency] = true
		}
	}
	if cycle := dependencyCycle(dependencies); cycle != "" {
		return fmt.Errorf("commitment policy: dependency cycle at %q", cycle)
	}
	for _, grant := range policy.ParallelGrants {
		if len(grant.Outcomes) < 2 {
			return errors.New("commitment policy: parallel grant needs at least two outcomes")
		}
		seen := map[string]bool{}
		for _, outcome := range grant.Outcomes {
			if !outcomes[outcome] || seen[outcome] {
				return fmt.Errorf("commitment policy: invalid parallel grant outcome %q", outcome)
			}
			seen[outcome] = true
		}
	}
	recorded := map[[2]string]bool{}
	for _, delivery := range policy.Deliveries {
		key := [2]string{delivery.Outcome, delivery.Binding}
		binding, bound := deliveredBinding(policy, delivery)
		if !bound || recorded[key] || delivery.Identity != binding.Source.Identity || !lineValue(delivery.Source) || !lineValue(delivery.Evidence) {
			return fmt.Errorf("commitment policy: invalid delivery for outcome %q", delivery.Outcome)
		}
		recorded[key] = true
	}
	delivered := deliveredOutcomes(policy)
	completed := map[string]bool{}
	for _, completion := range policy.Completions {
		milestone, known := findMilestone(policy, completion.Milestone)
		if !known || completed[completion.Milestone] || completion.Milestone == policy.ActiveMilestone || !lineValue(completion.Verification) ||
			slices.ContainsFunc(milestone.Outcomes, func(outcome Outcome) bool { return !delivered[outcome.ID] }) {
			return fmt.Errorf("commitment policy: invalid completion of milestone %q", completion.Milestone)
		}
		completed[completion.Milestone] = true
	}
	return nil
}

func validIdentity(value string) bool { return identityPattern.MatchString(value) }

func lineValue(value string) bool { return value != "" && sanitize.LineSafe(value) }

func dependencyCycle(graph map[string][]string) string {
	const (
		unseen = iota
		visiting
		visited
	)
	state := map[string]int{}
	var visit func(string) string
	visit = func(node string) string {
		if state[node] == visiting {
			return node
		}
		if state[node] == visited {
			return ""
		}
		state[node] = visiting
		for _, dependency := range graph[node] {
			if cycle := visit(dependency); cycle != "" {
				return cycle
			}
		}
		state[node] = visited
		return ""
	}
	for node := range graph {
		if cycle := visit(node); cycle != "" {
			return cycle
		}
	}
	return ""
}

func validateSource(source SourceBinding) error {
	if !validIdentity(source.ID) || !repositoryPath(source.Path) || !sanitize.LineSafe(source.Identity) || source.Identity == "" {
		return fmt.Errorf("commitment policy: invalid source binding %q", source.ID)
	}
	return nil
}

// repositoryPath reports whether value is one clean relative path inside the repository.
func repositoryPath(value string) bool {
	path := filepath.ToSlash(filepath.Clean(value))
	return value != "" && sanitize.LineSafe(value) && path == value && !filepath.IsAbs(value) && path != ".." && !strings.HasPrefix(path, "../")
}

func validateDeliverables(outcome Outcome) error {
	ids, paths := map[string]bool{}, map[string]bool{}
	for _, binding := range outcome.Deliverables {
		if err := validateSource(binding.Source); err != nil {
			return err
		}
		if ids[binding.Source.ID] || paths[binding.Source.Path] {
			return fmt.Errorf("commitment policy: duplicate deliverable %q", binding.Source.ID)
		}
		ids[binding.Source.ID], paths[binding.Source.Path] = true, true
		seen := map[string]bool{}
		for _, id := range binding.Obligations {
			found := false
			for _, source := range outcome.Sources {
				if source.ID == id {
					found = true
				}
			}
			if !found || seen[id] {
				return fmt.Errorf("commitment policy: invalid deliverable obligation %q", id)
			}
			seen[id] = true
		}
	}
	return nil
}
