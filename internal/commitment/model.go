// Package commitment owns delivery-policy validation and transition decisions.
package commitment

import (
	"path"
	"strings"
)

// PolicyPath is the project-owned tracked commitment.
const PolicyPath = ".bench/commitment.json"

// Policy is the tracked delivery commitment.
type Policy struct {
	Version         int             `json:"version"`
	Milestones      []Milestone     `json:"milestones"`
	ActiveMilestone string          `json:"active_milestone,omitempty"`
	ParallelGrants  []ParallelGrant `json:"parallel_grants,omitempty"`
	Deliveries      []DeliveryFact  `json:"deliveries,omitempty"`
}

// Milestone is one immutable milestone identity and its ordered outcomes.
type Milestone struct {
	ID       string    `json:"id"`
	Outcomes []Outcome `json:"outcomes"`
}

// Outcome is one immutable delivery outcome.
type Outcome struct {
	ID           string            `json:"id"`
	Criteria     []Criterion       `json:"criteria"`
	Sources      []SourceBinding   `json:"sources"`
	Deliverables []DeliveryBinding `json:"deliverables,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
}

// DeliveryBinding identifies an approved deliverable and its fully satisfied sources.
type DeliveryBinding struct {
	Source      SourceBinding `json:"source"`
	Obligations []string      `json:"obligations,omitempty"`
}

// Criterion is one immutable outcome criterion.
type Criterion struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// SourceBinding ties an outcome to exact repository content.
type SourceBinding struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Identity string `json:"identity"`
}

// ParallelGrant authorizes the named outcomes to run together.
type ParallelGrant struct {
	Outcomes []string `json:"outcomes"`
}

// DeliveryFact is broker-authored evidence for one delivered outcome.
type DeliveryFact struct {
	Milestone string `json:"milestone"`
	Outcome   string `json:"outcome"`
	Binding   string `json:"binding"`
	Identity  string `json:"identity"`
}

// Projection is the canonical active-policy view consumed by commands.
type Projection struct {
	Milestone string
	Outcomes  []string
}

// Selection projects the active milestone and ordered outcomes.
func Selection(policy Policy) Projection {
	for _, milestone := range policy.Milestones {
		if milestone.ID != policy.ActiveMilestone {
			continue
		}
		projection := Projection{Milestone: milestone.ID, Outcomes: make([]string, 0, len(milestone.Outcomes))}
		for _, outcome := range milestone.Outcomes {
			projection.Outcomes = append(projection.Outcomes, outcome.ID)
		}
		return projection
	}
	return Projection{}
}

// PlanningPath classifies documentation paths. Promotion paths must be named by a planning artifact.
func PlanningPath(name, mode string, promotions []string) bool {
	if mode != "000000" && mode != "100644" {
		return false
	}
	if name == PolicyPath {
		return true
	}
	if path.Base(name) == "AGENTS.md" || path.Base(name) == "CLAUDE.md" {
		return false
	}
	if path.Ext(name) != ".md" {
		return false
	}
	if name == "ROADMAP.md" {
		return true
	}
	for _, prefix := range []string{"specs/", "roadmap/", "capture/", "research/", "decisions/"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	for _, promoted := range promotions {
		if name == promoted && (name == "CONTEXT.md" || strings.HasPrefix(name, "docs/adr/")) {
			return true
		}
	}
	return false
}
