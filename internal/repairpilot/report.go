package repairpilot

import (
	"sort"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/assessment"
	"github.com/gibbonmi/bench/internal/toon"
)

var requiredReportClasses = []string{"stalled", "productive", "unchanged", "overlap"}

const (
	intervalAbsent    = "absent"
	intervalComplete  = "complete"
	intervalPartial   = "partial"
	intervalUnbounded = "unbounded"
	intervalUnproven  = "unproven"
)

type reportSummary struct {
	State       string
	Sample      string
	ActivatedAt string
	EndsAt      string
	CutoffAt    string
	Completed   int
	Incomplete  int
	Unknown     int
	Gaps        []evidenceGap
	Classes     map[string]bool
	Legacy      pilotSummary
}

type evidenceGap struct {
	ObservationID string
	Kind          string
	Detail        string
}

type collectionState struct {
	Name     string
	Deadline time.Time
	Cutoff   *time.Time
}

func renderPilotReport(document Document, now time.Time, full bool) (string, int) {
	summary := summarizeReport(document, now)
	columns := []string{"state", "activated_at", "collection_ends_at", "cutoff_at", "sample", "completed_sequences", "incomplete_sequences", "unknown_labels", "evidence_gaps", "comparable_repeats"}
	row := []string{summary.State, summary.ActivatedAt, summary.EndsAt, summary.CutoffAt, summary.Sample, strconv.Itoa(summary.Completed), strconv.Itoa(summary.Incomplete), strconv.Itoa(summary.Unknown), strconv.Itoa(len(summary.Gaps)), strconv.Itoa(summary.Legacy.ComparableRepeats)}
	for _, label := range append(append([]string{}, progressLabels...), "unknown") {
		columns = append(columns, "progress_"+label)
		row = append(row, strconv.Itoa(summary.Legacy.Progress[label]))
	}
	out, err := toon.Table("repair_pilot", columns, [][]string{row})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	classRows := make([][]string, 0, len(requiredReportClasses))
	for _, class := range requiredReportClasses {
		status := "missing"
		if summary.Classes[class] {
			status = "present"
		}
		classRows = append(classRows, []string{class, status})
	}
	classTable, err := toon.Table("required_classes", []string{"class", "status"}, classRows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out += classTable
	if !full {
		return out, 0
	}
	detail, err := renderFullReport(document, summary)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + detail, 0
}

func summarizeReport(document Document, now time.Time) reportSummary {
	if document.Version == 0 {
		return reportSummary{State: "inactive", Sample: "inconclusive", Classes: map[string]bool{}}
	}
	state := deriveCollectionState(document, now)
	legacy := summarizeDocument(document)
	summary := reportSummary{
		State: state.Name, Sample: "provisional", ActivatedAt: formatReportTime(&document.ActivatedAt), EndsAt: formatReportTime(&state.Deadline), CutoffAt: formatReportTime(state.Cutoff),
		Classes: map[string]bool{}, Gaps: reportEvidenceGaps(document), Completed: completedSequenceCount(document), Unknown: legacy.Progress["unknown"], Legacy: legacy,
	}
	sequences := map[sequence]bool{}
	for _, item := range document.Observations {
		sequences[item.Sequence] = true
		label := effectiveProgressLabel(document, item.ID)
		if item.Kind == "repair" && (label == "productive" || label == "stalled") {
			summary.Classes[label] = true
		}
		if item.Kind == "rerun" && item.Rerun != nil && item.Rerun.UnchangedContent && assessment.ValidReference(item.Rerun.VerificationReference) {
			summary.Classes["unchanged"] = true
		}
	}
	summary.Incomplete = len(sequences) - summary.Completed
	summary.Classes["overlap"] = hasAssignmentOverlap(document.Observations)
	if summary.State == "stopped" {
		summary.Sample = "coverage-complete"
		for _, class := range requiredReportClasses {
			if !summary.Classes[class] {
				summary.Sample = "inconclusive"
				break
			}
		}
	}
	return summary
}

func reportEvidenceGaps(document Document) []evidenceGap {
	gaps := []evidenceGap{}
	for _, item := range document.Observations {
		switch classifyIntervalEvidence(item) {
		case intervalUnbounded:
			gaps = append(gaps, evidenceGap{item.ID, "unbounded-interval", "interval provenance has no activity bounds"})
		case intervalPartial:
			gaps = append(gaps, evidenceGap{item.ID, "partial-interval", "both interval bounds are required for comparison"})
		case intervalUnproven:
			gaps = append(gaps, evidenceGap{item.ID, "unproven-interval", "interval provenance is required for comparison"})
		}
	}
	sort.Slice(gaps, func(i, j int) bool {
		if gaps[i].ObservationID != gaps[j].ObservationID {
			return gaps[i].ObservationID < gaps[j].ObservationID
		}
		return gaps[i].Kind < gaps[j].Kind
	})
	return gaps
}

func hasAssignmentOverlap(items []observation) bool {
	ordered := sortedObservations(items)
	for i, left := range ordered {
		for _, right := range ordered[i+1:] {
			if status, comparable := intervalComparisonStatus(left, right); comparable && status == "true" {
				return true
			}
		}
	}
	return false
}

func intervalComparisonStatus(left, right observation) (string, bool) {
	if left.AssignmentID == right.AssignmentID || left.Sequence != right.Sequence {
		return "", false
	}
	leftStatus := classifyIntervalEvidence(left)
	rightStatus := classifyIntervalEvidence(right)
	if leftStatus == intervalAbsent || rightStatus == intervalAbsent {
		return "", false
	}
	if leftStatus != intervalComplete || rightStatus != intervalComplete {
		return "unknown", true
	}
	if !left.StartedAt.Before(*left.EndedAt) || !right.StartedAt.Before(*right.EndedAt) {
		return "false", true
	}
	return strconv.FormatBool(left.StartedAt.Before(*right.EndedAt) && right.StartedAt.Before(*left.EndedAt)), true
}

func classifyIntervalEvidence(item observation) string {
	hasStart := item.StartedAt != nil
	hasEnd := item.EndedAt != nil
	switch {
	case !hasStart && !hasEnd && item.IntervalReference == nil:
		return intervalAbsent
	case !hasStart && !hasEnd:
		return intervalUnbounded
	case hasStart != hasEnd:
		return intervalPartial
	case item.IntervalReference == nil || !assessment.ValidReference(*item.IntervalReference):
		return intervalUnproven
	default:
		return intervalComplete
	}
}

func sortedObservations(items []observation) []observation {
	ordered := append([]observation{}, items...)
	sort.SliceStable(ordered, func(i, j int) bool {
		left, right := ordered[i], ordered[j]
		if comparison := compareSequences(left.Sequence, right.Sequence); comparison != 0 {
			return comparison < 0
		}
		if !left.ObservedAt.Equal(right.ObservedAt) {
			return left.ObservedAt.Before(right.ObservedAt)
		}
		return left.ID < right.ID
	})
	return ordered
}

func compareSequences(left, right sequence) int {
	for _, pair := range [][2]string{{left.Source, right.Source}, {left.Spec, right.Spec}, {left.Chunk, right.Chunk}} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	return 0
}

func deriveCollectionState(document Document, now time.Time) collectionState {
	deadline := collectionDeadline(document)
	state := collectionState{Name: "active", Deadline: deadline}
	if document.CutoffAt != nil {
		state.Name, state.Cutoff = "stopped", document.CutoffAt
	} else if !now.Before(deadline) {
		cutoff := deadline
		state.Name, state.Cutoff = "stopped", &cutoff
	}
	return state
}

func formatReportTime(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
