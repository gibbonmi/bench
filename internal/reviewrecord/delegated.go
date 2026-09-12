package reviewrecord

import (
	"errors"
	"fmt"
	"sort"
)

// Assignment is one delegated author occupancy of a ticket. A history's last
// entry supplies the ticket's effective author. Earlier entries stay, so a
// historical occurrence keeps the author and source that produced it.
type Assignment struct {
	Session      string `json:"session"`
	Assignment   string `json:"assignment"`
	Model        string `json:"model"`
	Effort       string `json:"effort"`
	Source       string `json:"source"`
	NativeRef    string `json:"native_ref"`
	Predecessor  string `json:"predecessor,omitempty"`
	Trigger      string `json:"trigger,omitempty"`
	Stopped      string `json:"stopped,omitempty"`
	Preserved    string `json:"preserved,omitempty"`
	Reassessment string `json:"reassessment,omitempty"`
}

// Execution declares the delegated run. A version 2 plan owns exactly one.
type Execution struct {
	Mode                string                  `json:"mode"`
	RunID               string                  `json:"run_id"`
	OrchestratorSession string                  `json:"orchestrator_session"`
	AuthorLimit         int                     `json:"author_limit"`
	Assignments         map[string][]Assignment `json:"assignments"`
}

// verificationRoles is the closed set of verification evidence roles. A version
// 1 record uses the author role for every obligation. A version 2 record adds
// the integration role, which only the orchestrator's final results carry.
func verificationRoles() []string {
	return []string{"author-verification", "integration-verification"}
}

// Triggers is the closed set of author-transfer reasons.
func Triggers() []string {
	return []string{"no-progress", "terminal-failure", "cap-exhausted", "session-lost"}
}

// Delegated reports whether the plan carries the version 2 delegated form.
func (p Plan) Delegated() bool { return p.Version == 2 && p.Execution != nil }

// Author is a ticket's effective author. An undispatched ticket has none.
func (p Plan) Author(ticket string) (string, bool) {
	if !p.Delegated() {
		return "", false
	}
	history := p.Execution.Assignments[ticket]
	if len(history) == 0 {
		return "", false
	}
	return history[len(history)-1].Session, true
}

// Participants lists the orchestrator and every session that authored in this
// run. An independent review session must appear in neither role.
func (p Plan) Participants() []string {
	if !p.Delegated() {
		return nil
	}
	result := []string{p.Execution.OrchestratorSession}
	for _, name := range p.ticketNames() {
		for _, item := range p.Execution.Assignments[name] {
			if !contains(result, item.Session) {
				result = append(result, item.Session)
			}
		}
	}
	return result
}

// ticketNames is every planned ticket basename in a stable order, so a refusal
// names the same ticket on every run.
func (p Plan) ticketNames() []string {
	names := []string{}
	for _, chunk := range p.Chunks {
		names = append(names, chunk.Tickets...)
	}
	sort.Strings(names)
	return names
}

// validateExecution grades the delegated declaration against the planned
// tickets. The plan is the sole identity authority, so every refusal here
// denies acceptance rather than warns.
func validateExecution(plan Plan) error {
	if plan.Version == 1 {
		if plan.Execution != nil {
			return errors.New("invalid execution declaration in a version 1 plan")
		}
		return nil
	}
	execution := plan.Execution
	if execution == nil {
		return errors.New("missing execution declaration; a version 2 plan declares one delegated run")
	}
	if execution.Mode != "delegate" {
		return fmt.Errorf("invalid execution mode %q; use delegate", execution.Mode)
	}
	if execution.RunID == "" || execution.OrchestratorSession == "" {
		return errors.New("missing run id or orchestrator session")
	}
	if execution.AuthorLimit < 1 {
		return errors.New("invalid author limit; declare at least one concurrent writer")
	}
	planned := map[string]bool{}
	for _, name := range plan.ticketNames() {
		planned[name] = true
	}
	for name := range execution.Assignments {
		if !planned[name] {
			return fmt.Errorf("invalid assignment ticket %s; name a planned ticket", name)
		}
	}
	owners := map[string]string{}
	for _, name := range plan.ticketNames() {
		history, present := execution.Assignments[name]
		if !present {
			return fmt.Errorf("missing assignment history for ticket %s; an undispatched ticket declares an empty history", name)
		}
		for i, item := range history {
			if err := validateAssignment(item, i, history, execution.OrchestratorSession); err != nil {
				return fmt.Errorf("ticket %s assignment %d: %w", name, i, err)
			}
			if owner, taken := owners[item.Session]; taken && owner != name {
				return fmt.Errorf("ticket %s: session already authors ticket %s; one session owns one ticket", name, owner)
			}
			owners[item.Session] = name
		}
	}
	return nil
}

func validateAssignment(item Assignment, index int, history []Assignment, orchestrator string) error {
	if item.Session == "" || item.Assignment == "" || item.Model == "" || item.Effort == "" || item.Source == "" || item.NativeRef == "" {
		return errors.New("missing session, assignment, model, effort, source, or native dispatch reference")
	}
	if item.Session == orchestrator {
		return errors.New("author and orchestrator identities must differ")
	}
	if index == 0 {
		if item.Predecessor != "" || item.Trigger != "" || item.Stopped != "" || item.Preserved != "" || item.Reassessment != "" {
			return errors.New("invalid replacement fields on a first assignment")
		}
		return nil
	}
	previous := history[index-1]
	if item.Predecessor != previous.Session || item.Session == previous.Session {
		return errors.New("invalid predecessor; name the distinct preceding author")
	}
	if !contains(Triggers(), item.Trigger) {
		return fmt.Errorf("invalid replacement trigger %q", item.Trigger)
	}
	if item.Stopped == "" {
		return errors.New("missing stopped-writer evidence; confirm the old writer stopped")
	}
	if item.Preserved == "" {
		return errors.New("missing preserved source")
	}
	if (item.Reassessment != "") != (item.Trigger == "no-progress") {
		return errors.New("a no-progress trigger requires a reassessment, and no other trigger carries one")
	}
	return nil
}

// verifier resolves who owes a verification requirement, and under which role.
// Version 1 keeps one implementation session for every obligation. Version 2
// reads each chunk obligation's owning ticket from this same frozen plan, so a
// historical chunk keeps the author its own plan named.
func verifier(plan Plan, record Record, requirement Requirement, final bool) (string, string, error) {
	if !plan.Delegated() {
		return record.ImplementationSession, "author-verification", nil
	}
	if final {
		return plan.Execution.OrchestratorSession, "integration-verification", nil
	}
	session, dispatched := plan.Author(requirement.Ticket)
	if !dispatched {
		return "", "", fmt.Errorf("verification %s names undispatched ticket %s; dispatch an author before its obligation", requirement.ID, requirement.Ticket)
	}
	return session, "author-verification", nil
}

// reviewExclusions is the set of identities that cannot supply independent
// review. Version 1 excludes its one implementation session. Version 2 excludes
// the orchestrator and every current and former author in the run.
func reviewExclusions(plan Plan, record Record) []string {
	if plan.Delegated() {
		return plan.Participants()
	}
	return []string{record.ImplementationSession}
}
