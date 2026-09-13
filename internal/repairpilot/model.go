package repairpilot

import (
	"time"

	"github.com/gibbonmi/bench/internal/assessment"
)

var progressLabels = []string{"productive", "stalled", "unchanged"}

type sequence struct {
	Source string `json:"source"`
	Spec   string `json:"spec"`
	Chunk  string `json:"chunk"`
}

type failure struct {
	Check      string               `json:"check"`
	Identity   string               `json:"identity,omitempty"`
	Diagnostic string               `json:"diagnostic"`
	Ownership  string               `json:"ownership"`
	Blocking   bool                 `json:"blocking"`
	Reference  assessment.Reference `json:"reference"`
}

type repair struct {
	Hypothesis            string               `json:"hypothesis"`
	IntendedChange        string               `json:"intended_change"`
	VerificationReference assessment.Reference `json:"verification_reference"`
}

type rerun struct {
	VerificationReference assessment.Reference `json:"verification_reference"`
	UnchangedContent      bool                 `json:"unchanged_content"`
}

type endpoint struct {
	Kind               string                 `json:"kind"`
	Reference          assessment.Reference   `json:"reference"`
	CoveredBlockerRefs []assessment.Reference `json:"covered_blocker_refs,omitempty"`
}

type progressProposal struct {
	Label  string `json:"label"`
	Reason string `json:"reason"`
}

type observation struct {
	ID                  string                 `json:"id"`
	ObservedAt          time.Time              `json:"observed_at"`
	Sequence            sequence               `json:"sequence"`
	AssignmentID        string                 `json:"assignment_id"`
	SessionID           string                 `json:"session_id"`
	SourceRevision      string                 `json:"source_revision"`
	Stage               string                 `json:"stage"`
	Kind                string                 `json:"kind"`
	FailureCompleteness string                 `json:"failure_completeness,omitempty"`
	Failures            []failure              `json:"failures,omitempty"`
	References          []assessment.Reference `json:"references"`
	StartedAt           *time.Time             `json:"started_at,omitempty"`
	EndedAt             *time.Time             `json:"ended_at,omitempty"`
	IntervalReference   *assessment.Reference  `json:"interval_reference,omitempty"`
	Repair              *repair                `json:"repair,omitempty"`
	Rerun               *rerun                 `json:"rerun,omitempty"`
	Endpoint            *endpoint              `json:"endpoint,omitempty"`
	ProgressProposal    *progressProposal      `json:"progress_proposal,omitempty"`
}

type progressConclusion struct {
	Label                 string               `json:"label"`
	RepairedTarget        string               `json:"repaired_target,omitempty"`
	VerificationReference assessment.Reference `json:"verification_reference,omitempty"`
}

type audit struct {
	ID                 string                 `json:"id"`
	ObservationID      string                 `json:"observation_id"`
	Auditor            string                 `json:"auditor"`
	EvidenceReferences []assessment.Reference `json:"evidence_references"`
	Reason             string                 `json:"reason"`
	Conclusion         *progressConclusion    `json:"conclusion,omitempty"`
	ResolvesAuditIDs   []string               `json:"resolves_audit_ids,omitempty"`
}

type recordInput struct {
	Version     int          `json:"version"`
	Observation *observation `json:"observation,omitempty"`
	Audit       *audit       `json:"audit,omitempty"`
}
