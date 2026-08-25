package gate

import (
	"errors"
	"slices"
	"time"
)

type verdictRecordClass struct {
	name   string
	fields []string
	// reuseRefusal is why a green of this class may never be reused, and "" when it
	// may. The class declares it, so a class added to the store states its own reuse
	// posture instead of inheriting whatever a name match happens to give it.
	reuseRefusal string
	validate     func([]byte, verdictRecord, time.Time) error
}

// isGateVerdict reports whether this class grades the oracle. The oracle identity is
// the field that says so, so the one field list already in the class answers it. A
// second marker beside that list could disagree with it.
func (c verdictRecordClass) isGateVerdict() bool {
	return slices.Contains(c.fields, "oracle")
}

var verdictRecordClasses = []verdictRecordClass{
	{
		name:     "full verdict",
		fields:   []string{"oracle", "recorded_at", "schema", "state", "status", "tree"},
		validate: validateFullRecord,
	},
	{
		name:         "partial verdict",
		reuseRefusal: "partial verdict",
		fields:       []string{"executed", "oracle", "recorded_at", "schema", "skip_evidence", "skipped", "state", "status", "tree"},
		validate:     validatePartialRecord,
	},
	{
		name:         "check-partial verdict",
		reuseRefusal: "partial verdict",
		fields:       []string{"check_evidence", "check_executed", "check_inherited", "oracle", "recorded_at", "schema", "state", "status", "tree"},
		validate:     validateCheckPartialRecord,
	},
	{
		name:         "combined-partial verdict",
		reuseRefusal: "partial verdict",
		fields:       []string{"check_evidence", "check_executed", "check_inherited", "executed", "oracle", "recorded_at", "schema", "skip_evidence", "skipped", "state", "status", "tree"},
		validate:     validateCombinedPartialRecord,
	},
	{
		name:     "pending",
		fields:   []string{"oracle", "owner_pid", "schema", "started_at", "state", "tree"},
		validate: validatePendingRecord,
	},
	{
		name:         "lane record",
		fields:       []string{"lane", "outcome", "recorded_at", "run_binary", "schema", "tree"},
		reuseRefusal: "lane record",
		validate:     validateLaneRecord,
	},
}

// validateLaneRecord grades one fast-lane record. A lane grades a declared check list
// and not the oracle, so the record carries an outcome rather than a verdict status.
// The run binary is recorded as its content address, which is what ties the outcome to
// the code that produced it.
func validateLaneRecord(_ []byte, r verdictRecord, now time.Time) error {
	if r.Lane == "" || (r.Outcome != lanePass && r.Outcome != laneFail) || !isContentAddress(r.RunBinary) {
		return errors.New("invalid lane record")
	}
	if tm, err := strictRecordTime(r.RecordedAt); err != nil || tm.After(now) {
		return errors.New("invalid lane record time")
	}
	return nil
}

func selectVerdictRecordClass(data []byte) (verdictRecordClass, error) {
	for _, class := range verdictRecordClasses {
		if err := requireObjectFields(data, class.fields); err == nil {
			return class, nil
		}
	}
	return verdictRecordClass{}, errors.New("invalid record fields")
}

func validateFullRecord(_ []byte, r verdictRecord, now time.Time) error {
	return validateReadyRecord(r, now)
}

func validatePartialRecord(data []byte, r verdictRecord, now time.Time) error {
	if err := validateReadyRecord(r, now); err != nil {
		return err
	}
	tm, _ := time.Parse(time.RFC3339, r.RecordedAt)
	return validatePartition(data, r, tm)
}

func validateCheckPartialRecord(data []byte, r verdictRecord, now time.Time) error {
	if err := validateReadyRecord(r, now); err != nil {
		return err
	}
	tm, _ := time.Parse(time.RFC3339, r.RecordedAt)
	return validateCheckPartition(data, r, tm)
}

func validateCombinedPartialRecord(data []byte, r verdictRecord, now time.Time) error {
	if err := validateReadyRecord(r, now); err != nil {
		return err
	}
	tm, _ := time.Parse(time.RFC3339, r.RecordedAt)
	if err := validatePartition(data, r, tm); err != nil {
		return err
	}
	return validateCheckPartition(data, r, tm)
}

func validatePendingRecord(_ []byte, r verdictRecord, now time.Time) error {
	if r.State != Pending || r.Status != "" || r.RecordedAt != "" || r.StartedAt == "" || r.OwnerPID <= 0 {
		return errors.New("invalid pending")
	}
	if tm, err := strictRecordTime(r.StartedAt); err != nil || tm.After(now) {
		return errors.New("invalid pending time")
	}
	return nil
}

func validateReadyRecord(r verdictRecord, now time.Time) error {
	if r.State != Ready || (r.Status != "green" && r.Status != "red" && r.Status != "timeout") || r.RecordedAt == "" || r.StartedAt != "" || r.OwnerPID != 0 {
		return errors.New("invalid ready")
	}
	if tm, err := strictRecordTime(r.RecordedAt); err != nil || tm.After(now) {
		return errors.New("invalid ready time")
	}
	return nil
}
