package repairpilot

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gibbonmi/bench/internal/assessment"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

func record(options Options, path string) (string, int) {
	input, err := readRecordInput(path)
	if err != nil {
		return refusal(err.Error())
	}
	var document Document
	err = withPilotLock(options, false, func(files FileOps) error {
		loaded, state, loadErr := load(options)
		if loadErr != nil {
			return loadErr
		}
		if state == bounds.StateAbsent {
			return errors.New("repair pilot is inactive; activate it before recording evidence")
		}
		document = loaded
		changed, applyErr := applyRecordInput(&document, input, options.Now.UTC())
		if applyErr != nil {
			return applyErr
		}
		if !changed {
			return nil
		}
		return (Store{Path: documentPath(options), Files: files}).Replace(document)
	})
	if err != nil {
		return refusal(err.Error())
	}
	return renderDocumentStatus(document, options.Now)
}

func readRecordInput(path string) (recordInput, error) {
	if err := bounds.RefuseLinks(path); err != nil {
		return recordInput{}, fmt.Errorf("record input path is unsafe: %w", err)
	}
	classified := bounds.ClassifyNoFollow(path)
	if classified.State != bounds.StateParsed {
		return recordInput{}, fmt.Errorf("record input is %s", classified.State)
	}
	var input recordInput
	if err := jsonfile.DecodeDocument(classified.Data, &input); err != nil {
		return recordInput{}, fmt.Errorf("record input is malformed: %w", err)
	}
	if input.Version != 1 || (input.Observation == nil) == (input.Audit == nil) {
		return recordInput{}, errors.New("record input must contain exactly one version 1 observation or audit")
	}
	return input, nil
}

func applyRecordInput(document *Document, input recordInput, now time.Time) (bool, error) {
	if input.Observation != nil {
		return appendObservation(document, *input.Observation, now)
	}
	return appendAudit(document, *input.Audit)
}

func appendObservation(document *Document, candidate observation, now time.Time) (bool, error) {
	for _, stored := range document.Observations {
		if stored.ID != candidate.ID {
			continue
		}
		if reflect.DeepEqual(stored, candidate) {
			return false, nil
		}
		return false, fmt.Errorf("observation id %q conflicts with retained evidence", candidate.ID)
	}
	deadline := collectionDeadline(*document)
	if !now.Before(deadline) || document.CutoffAt != nil || completedSequenceCount(*document) >= 10 {
		return false, errors.New("repair pilot collection is stopped; observations cannot be appended")
	}
	if err := validateObservation(*document, candidate, now, deadline); err != nil {
		return false, err
	}
	document.Observations = append(document.Observations, candidate)
	if candidate.Endpoint != nil && completedSequenceCount(*document) == 10 {
		cutoff := now.UTC()
		document.CutoffAt = &cutoff
	}
	return true, nil
}

func validateObservation(document Document, candidate observation, now, deadline time.Time) error {
	if !assessment.ValidID(candidate.ID) || !assessment.ValidID(candidate.Sequence.Source) || !assessment.ValidID(candidate.Sequence.Spec) || !assessment.ValidID(candidate.Sequence.Chunk) || !assessment.ValidID(candidate.AssignmentID) || !assessment.ValidID(candidate.SessionID) || !assessment.ValidID(candidate.SourceRevision) {
		return errors.New("observation is missing identity or attribution evidence")
	}
	if candidate.Stage != "pre-review" && candidate.Stage != "post-review" {
		return errors.New("observation has unsupported work stage")
	}
	if !validReferences(candidate.References) {
		return errors.New("observation requires native evidence references")
	}
	for _, stamp := range []*time.Time{&candidate.ObservedAt, candidate.StartedAt, candidate.EndedAt} {
		if stamp == nil {
			continue
		}
		if stamp.IsZero() || stamp.Before(document.ActivatedAt) || stamp.After(now) || !stamp.Before(deadline) {
			return errors.New("observation timestamp is outside the active collection window")
		}
	}
	if candidate.StartedAt != nil && candidate.EndedAt != nil && candidate.EndedAt.Before(*candidate.StartedAt) {
		return errors.New("observation interval ends before it starts")
	}
	if (candidate.StartedAt != nil || candidate.EndedAt != nil) && (candidate.IntervalReference == nil || !assessment.ValidReference(*candidate.IntervalReference)) {
		return errors.New("observation interval requires native evidence")
	}
	if candidate.IntervalReference != nil && !assessment.ValidReference(*candidate.IntervalReference) {
		return errors.New("observation interval has invalid evidence")
	}
	if err := validateObservationKind(candidate); err != nil {
		return err
	}
	prior := sequenceObservations(document, candidate.Sequence)
	if len(prior) == 0 && (candidate.Kind != "failure" || !hasBlockingFailure(candidate)) {
		return errors.New("the first sequence observation must be a failure with a blocking result")
	}
	for _, item := range prior {
		if item.Endpoint != nil {
			return errors.New("a completed sequence cannot accept another observation")
		}
	}
	if candidate.Endpoint != nil {
		return validateEndpoint(prior, candidate)
	}
	return nil
}

func validateObservationKind(candidate observation) error {
	if len(candidate.Failures) != 0 {
		if candidate.FailureCompleteness != "complete" && candidate.FailureCompleteness != "first-only" && candidate.FailureCompleteness != "unknown" {
			return errors.New("failure set has unsupported completeness")
		}
		for _, item := range candidate.Failures {
			if item.Check == "" || item.Diagnostic == "" || item.Identity != "" && !assessment.ValidID(item.Identity) || !assessment.ValidReference(item.Reference) || !validOwnership(item.Ownership) {
				return errors.New("failure is missing its check, diagnostic, ownership, or reference")
			}
		}
	}
	switch candidate.Kind {
	case "failure":
		if len(candidate.Failures) == 0 || candidate.Repair != nil || candidate.Rerun != nil {
			return errors.New("failure observation has an invalid payload")
		}
	case "repair":
		if candidate.Repair == nil || candidate.Rerun != nil || candidate.Repair.Hypothesis == "" || candidate.Repair.IntendedChange == "" || !assessment.ValidReference(candidate.Repair.VerificationReference) {
			return errors.New("repair observation requires its hypothesis, intended change, and verification reference")
		}
	case "rerun":
		if candidate.Rerun == nil || candidate.Repair != nil || !assessment.ValidReference(candidate.Rerun.VerificationReference) {
			return errors.New("rerun observation requires a verification reference")
		}
	default:
		return errors.New("observation has unsupported kind")
	}
	if candidate.ProgressProposal != nil && (candidate.ProgressProposal.Reason == "" || !validProgressLabel(candidate.ProgressProposal.Label)) {
		return errors.New("progress proposal has unsupported label or no reason")
	}
	return nil
}

func validateEndpoint(prior []observation, candidate observation) error {
	value := candidate.Endpoint
	if !assessment.ValidReference(value.Reference) || len(value.CoveredBlockerRefs) != 0 && !validReferences(value.CoveredBlockerRefs) {
		return errors.New("sequence endpoint requires evidence")
	}
	switch value.Kind {
	case "reviewer-handoff":
		return nil
	case "verified-closure":
		covered := referenceSet(value.CoveredBlockerRefs)
		for _, item := range append(append([]observation{}, prior...), candidate) {
			for _, blocker := range item.Failures {
				if blocker.Blocking && !covered[blocker.Reference] {
					return fmt.Errorf("verified closure does not cover blocker %q", blocker.Reference)
				}
			}
		}
		return nil
	default:
		return errors.New("sequence endpoint has unsupported kind")
	}
}

func appendAudit(document *Document, candidate audit) (bool, error) {
	for _, stored := range document.Audits {
		if stored.ID != candidate.ID {
			continue
		}
		if reflect.DeepEqual(stored, candidate) {
			return false, nil
		}
		return false, fmt.Errorf("audit id %q conflicts with retained evidence", candidate.ID)
	}
	if !assessment.ValidID(candidate.ID) || !assessment.ValidID(candidate.ObservationID) || !assessment.ValidID(candidate.Auditor) || candidate.Reason == "" || !validReferences(candidate.EvidenceReferences) || !validIDs(candidate.ResolvesAuditIDs) {
		return false, errors.New("audit is missing identity, auditor, reason, or evidence references")
	}
	if observationByID(*document, candidate.ObservationID) == nil {
		return false, errors.New("audit names an observation that is not retained")
	}
	if candidate.Conclusion != nil && !validProgressLabel(candidate.Conclusion.Label) {
		return false, errors.New("audit conclusion has unsupported label")
	}
	if len(candidate.ResolvesAuditIDs) != 0 {
		if len(candidate.ResolvesAuditIDs) < 2 || candidate.Conclusion == nil || !supportedConclusion(*candidate.Conclusion) {
			return false, errors.New("resolution audit requires two earlier audits and a supported conclusion")
		}
		labels := map[string]bool{}
		for _, id := range candidate.ResolvesAuditIDs {
			prior := auditByID(*document, id)
			if prior == nil || prior.ObservationID != candidate.ObservationID || prior.Conclusion == nil || !supportedConclusion(*prior.Conclusion) {
				return false, errors.New("resolution audit does not cite retained audits for one observation")
			}
			labels[prior.Conclusion.Label] = true
		}
		if len(labels) < 2 {
			return false, errors.New("resolution audit does not cite conflicting supported conclusions")
		}
	}
	document.Audits = append(document.Audits, candidate)
	return true, nil
}

func effectiveProgressLabel(document Document, observationID string) string {
	labels := map[string]bool{}
	for _, item := range document.Audits {
		if item.ObservationID != observationID || item.Conclusion == nil || !supportedConclusion(*item.Conclusion) {
			continue
		}
		if len(item.ResolvesAuditIDs) >= 2 {
			labels = map[string]bool{item.Conclusion.Label: true}
			continue
		}
		labels[item.Conclusion.Label] = true
	}
	if len(labels) != 1 {
		return "unknown"
	}
	for label := range labels {
		return label
	}
	return "unknown"
}

func supportedConclusion(value progressConclusion) bool {
	return validProgressLabel(value.Label) && value.RepairedTarget != "" && assessment.ValidReference(value.VerificationReference)
}

func validProgressLabel(value string) bool {
	for _, label := range progressLabels {
		if value == label {
			return true
		}
	}
	return false
}

func validOwnership(value string) bool {
	return value == "diff-owned" || value == "inherited" || value == "spec-predicted" || value == "unknown"
}

func comparableFailure(left, right failure) bool {
	return left.Identity != "" && right.Identity != "" && left.Check == right.Check && left.Identity == right.Identity
}

func collectionDeadline(document Document) time.Time {
	return document.ActivatedAt.UTC().AddDate(0, 0, 14)
}

func completedSequenceCount(document Document) int {
	completed := map[sequence]bool{}
	for _, item := range document.Observations {
		if item.Endpoint != nil {
			completed[item.Sequence] = true
		}
	}
	return len(completed)
}

func sequenceObservations(document Document, key sequence) []observation {
	items := []observation{}
	for _, item := range document.Observations {
		if item.Sequence == key {
			items = append(items, item)
		}
	}
	return items
}

func observationByID(document Document, id string) *observation {
	for i := range document.Observations {
		if document.Observations[i].ID == id {
			return &document.Observations[i]
		}
	}
	return nil
}

func auditByID(document Document, id string) *audit {
	for i := range document.Audits {
		if document.Audits[i].ID == id {
			return &document.Audits[i]
		}
	}
	return nil
}

func hasBlockingFailure(value observation) bool {
	for _, item := range value.Failures {
		if item.Blocking {
			return true
		}
	}
	return false
}

func validReferences(values []assessment.Reference) bool {
	if len(values) == 0 {
		return false
	}
	for _, value := range values {
		if !assessment.ValidReference(value) {
			return false
		}
	}
	return true
}

func validIDs(values []string) bool {
	for _, value := range values {
		if !assessment.ValidID(value) {
			return false
		}
	}
	return true
}

func referenceSet(values []assessment.Reference) map[assessment.Reference]bool {
	result := map[assessment.Reference]bool{}
	for _, value := range values {
		result[value] = true
	}
	return result
}
