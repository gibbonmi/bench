package repairpilot

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/assessment"
	"github.com/gibbonmi/bench/internal/toon"
)

type reportTable struct {
	Name    string
	Columns []string
	Rows    [][]string
}

func renderFullReport(document Document, summary reportSummary) (string, error) {
	observations := sortedObservations(document.Observations)
	tables := []reportTable{
		sequenceReportTable(observations),
		observationReportTable(document, observations),
		failureReportTable(observations),
		referenceReportTable("observation_refs", observations),
		repairReportTable(observations),
		rerunReportTable(observations),
		endpointReportTable(observations),
		endpointReferenceReportTable(observations),
		auditReportTable(document.Audits),
		auditReferenceReportTable(document.Audits),
		intervalReportTable(observations),
		comparisonReportTable(observations),
		gapReportTable(summary.Gaps),
	}
	var out strings.Builder
	for _, table := range tables {
		rendered, err := toon.Table(table.Name, table.Columns, table.Rows)
		if err != nil {
			return "", err
		}
		out.WriteString(rendered)
	}
	return out.String(), nil
}

func sequenceReportTable(items []observation) reportTable {
	states := map[sequence]string{}
	for _, item := range items {
		if _, exists := states[item.Sequence]; !exists {
			states[item.Sequence] = "incomplete"
		}
		if item.Endpoint != nil {
			states[item.Sequence] = "complete"
		}
	}
	keys := make([]sequence, 0, len(states))
	for key := range states {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return compareSequences(keys[i], keys[j]) < 0
	})
	rows := make([][]string, 0, len(keys))
	for _, key := range keys {
		rows = append(rows, []string{key.Source, key.Spec, key.Chunk, states[key]})
	}
	return reportTable{"sequences", []string{"source", "spec", "chunk", "state"}, rows}
}

func observationReportTable(document Document, items []observation) reportTable {
	columns := []string{"source", "spec", "chunk", "observed_at", "id", "assignment_id", "session_id", "source_revision", "stage", "kind", "effective_label", "failure_completeness", "started_at", "ended_at", "proposal_label", "proposal_reason"}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		proposalLabel, proposalReason := "", ""
		if item.ProgressProposal != nil {
			proposalLabel, proposalReason = item.ProgressProposal.Label, item.ProgressProposal.Reason
		}
		rows = append(rows, []string{item.Sequence.Source, item.Sequence.Spec, item.Sequence.Chunk, formatReportTime(&item.ObservedAt), item.ID, item.AssignmentID, item.SessionID, item.SourceRevision, item.Stage, item.Kind, effectiveProgressLabel(document, item.ID), item.FailureCompleteness, formatReportTime(item.StartedAt), formatReportTime(item.EndedAt), proposalLabel, proposalReason})
	}
	return reportTable{"observations", columns, rows}
}

func failureReportTable(items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		for index, failure := range item.Failures {
			rows = append(rows, []string{item.ID, strconv.Itoa(index), failure.Check, failure.Identity, failure.Diagnostic, failure.Ownership, strconv.FormatBool(failure.Blocking), failure.Reference.Producer, failure.Reference.Native})
		}
	}
	return reportTable{"failures", []string{"observation_id", "index", "check", "identity", "diagnostic", "ownership", "blocking", "producer", "native"}, rows}
}

func referenceReportTable(name string, items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		for index, reference := range item.References {
			rows = append(rows, []string{item.ID, strconv.Itoa(index), reference.Producer, reference.Native})
		}
	}
	return reportTable{name, []string{"observation_id", "index", "producer", "native"}, rows}
}

func repairReportTable(items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		if item.Repair != nil {
			rows = append(rows, []string{item.ID, item.Repair.Hypothesis, item.Repair.IntendedChange, item.Repair.VerificationReference.Producer, item.Repair.VerificationReference.Native})
		}
	}
	return reportTable{"repairs", []string{"observation_id", "hypothesis", "intended_change", "verification_producer", "verification_native"}, rows}
}

func rerunReportTable(items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		if item.Rerun != nil {
			rows = append(rows, []string{item.ID, strconv.FormatBool(item.Rerun.UnchangedContent), item.Rerun.VerificationReference.Producer, item.Rerun.VerificationReference.Native})
		}
	}
	return reportTable{"reruns", []string{"observation_id", "unchanged_content", "verification_producer", "verification_native"}, rows}
}

func endpointReportTable(items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		if item.Endpoint != nil {
			rows = append(rows, []string{item.ID, item.Endpoint.Kind, item.Endpoint.Reference.Producer, item.Endpoint.Reference.Native})
		}
	}
	return reportTable{"endpoints", []string{"observation_id", "kind", "producer", "native"}, rows}
}

func endpointReferenceReportTable(items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		if item.Endpoint == nil {
			continue
		}
		for index, reference := range item.Endpoint.CoveredBlockerRefs {
			rows = append(rows, []string{item.ID, strconv.Itoa(index), reference.Producer, reference.Native})
		}
	}
	return reportTable{"endpoint_blocker_refs", []string{"observation_id", "index", "producer", "native"}, rows}
}

func sortedAudits(items []audit) []audit {
	ordered := append([]audit{}, items...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].ObservationID != ordered[j].ObservationID {
			return ordered[i].ObservationID < ordered[j].ObservationID
		}
		return ordered[i].ID < ordered[j].ID
	})
	return ordered
}

func auditReportTable(items []audit) reportTable {
	rows := [][]string{}
	for _, item := range sortedAudits(items) {
		label, target, producer, native := "", "", "", ""
		if item.Conclusion != nil {
			label, target = item.Conclusion.Label, item.Conclusion.RepairedTarget
			producer, native = item.Conclusion.VerificationReference.Producer, item.Conclusion.VerificationReference.Native
		}
		resolves := append([]string{}, item.ResolvesAuditIDs...)
		sort.Strings(resolves)
		rows = append(rows, []string{item.ObservationID, item.ID, item.Auditor, item.Reason, label, target, producer, native, strings.Join(resolves, " ")})
	}
	return reportTable{"audits", []string{"observation_id", "audit_id", "auditor", "reason", "conclusion", "repaired_target", "verification_producer", "verification_native", "resolves_audit_ids"}, rows}
}

func auditReferenceReportTable(items []audit) reportTable {
	rows := [][]string{}
	for _, item := range sortedAudits(items) {
		for index, reference := range item.EvidenceReferences {
			rows = append(rows, []string{item.ObservationID, item.ID, strconv.Itoa(index), reference.Producer, reference.Native})
		}
	}
	return reportTable{"audit_refs", []string{"observation_id", "audit_id", "index", "producer", "native"}, rows}
}

func intervalReportTable(items []observation) reportTable {
	rows := [][]string{}
	for _, item := range items {
		if item.StartedAt == nil && item.EndedAt == nil && item.IntervalReference == nil {
			continue
		}
		producer, native := referenceFields(item.IntervalReference)
		status := "complete"
		if item.StartedAt == nil && item.EndedAt == nil {
			status = "unbounded"
		} else if item.StartedAt == nil || item.EndedAt == nil {
			status = "partial"
		} else if item.IntervalReference == nil || !assessment.ValidReference(*item.IntervalReference) {
			status = "unproven"
		}
		rows = append(rows, []string{item.ID, item.AssignmentID, formatReportTime(item.StartedAt), formatReportTime(item.EndedAt), status, producer, native})
	}
	return reportTable{"intervals", []string{"observation_id", "assignment_id", "started_at", "ended_at", "status", "producer", "native"}, rows}
}

func comparisonReportTable(items []observation) reportTable {
	rows := [][]string{}
	for i, left := range items {
		for _, right := range items[i+1:] {
			if status, comparable := intervalComparisonStatus(left, right); comparable {
				rows = append(rows, []string{left.ID, right.ID, status})
			}
		}
	}
	return reportTable{"interval_comparisons", []string{"left_observation_id", "right_observation_id", "status"}, rows}
}

func gapReportTable(gaps []evidenceGap) reportTable {
	rows := make([][]string, 0, len(gaps))
	for _, gap := range gaps {
		rows = append(rows, []string{gap.ObservationID, gap.Kind, gap.Detail})
	}
	return reportTable{"evidence_gaps", []string{"observation_id", "kind", "detail"}, rows}
}

func referenceFields(reference *assessment.Reference) (string, string) {
	if reference == nil {
		return "", ""
	}
	return reference.Producer, reference.Native
}
