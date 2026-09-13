package repairpilot

import (
	"sort"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/assessment"
	"github.com/gibbonmi/bench/internal/toon"
)

var requiredReportClasses = []string{"stalled", "productive", "unchanged", "overlap"}

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
	deadline := collectionDeadline(document)
	summary := reportSummary{
		State: "active", Sample: "provisional", ActivatedAt: formatReportTime(&document.ActivatedAt), EndsAt: formatReportTime(&deadline),
		Classes: map[string]bool{}, Gaps: reportEvidenceGaps(document), Completed: completedSequenceCount(document), Legacy: summarizeDocument(document),
	}
	if document.CutoffAt != nil {
		summary.State = "stopped"
		summary.CutoffAt = formatReportTime(document.CutoffAt)
	} else if !now.Before(deadline) {
		summary.State = "stopped"
		summary.CutoffAt = formatReportTime(&deadline)
	}
	sequences := map[sequence]bool{}
	for _, item := range document.Observations {
		sequences[item.Sequence] = true
		label := effectiveProgressLabel(document, item.ID)
		if label == "unknown" {
			summary.Unknown++
		}
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
		hasStart := item.StartedAt != nil
		hasEnd := item.EndedAt != nil
		switch {
		case hasStart != hasEnd:
			gaps = append(gaps, evidenceGap{item.ID, "partial-interval", "both interval bounds are required for comparison"})
		case (hasStart || hasEnd) && (item.IntervalReference == nil || !assessment.ValidReference(*item.IntervalReference)):
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
			if intervalsOverlap(left, right) {
				return true
			}
		}
	}
	return false
}

func intervalsOverlap(left, right observation) bool {
	if left.AssignmentID == right.AssignmentID || left.StartedAt == nil || left.EndedAt == nil || right.StartedAt == nil || right.EndedAt == nil {
		return false
	}
	if left.IntervalReference == nil || right.IntervalReference == nil || !assessment.ValidReference(*left.IntervalReference) || !assessment.ValidReference(*right.IntervalReference) {
		return false
	}
	if !left.StartedAt.Before(*left.EndedAt) || !right.StartedAt.Before(*right.EndedAt) {
		return false
	}
	return left.StartedAt.Before(*right.EndedAt) && right.StartedAt.Before(*left.EndedAt)
}

func sortedObservations(items []observation) []observation {
	ordered := append([]observation{}, items...)
	sort.SliceStable(ordered, func(i, j int) bool {
		left, right := ordered[i], ordered[j]
		if left.Sequence.Source != right.Sequence.Source {
			return left.Sequence.Source < right.Sequence.Source
		}
		if left.Sequence.Spec != right.Sequence.Spec {
			return left.Sequence.Spec < right.Sequence.Spec
		}
		if left.Sequence.Chunk != right.Sequence.Chunk {
			return left.Sequence.Chunk < right.Sequence.Chunk
		}
		if !left.ObservedAt.Equal(right.ObservedAt) {
			return left.ObservedAt.Before(right.ObservedAt)
		}
		return left.ID < right.ID
	})
	return ordered
}

func formatReportTime(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
