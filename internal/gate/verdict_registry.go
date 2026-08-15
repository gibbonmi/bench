package gate

import (
	"errors"
	"time"
)

type verdictRecordClass struct {
	name     string
	fields   []string
	validate func([]byte, verdictRecord, time.Time) error
}

var verdictRecordClasses = []verdictRecordClass{
	{
		name:     "full verdict",
		fields:   []string{"oracle", "recorded_at", "schema", "state", "status", "tree"},
		validate: validateFullRecord,
	},
	{
		name:     "partial verdict",
		fields:   []string{"executed", "oracle", "recorded_at", "schema", "skip_evidence", "skipped", "state", "status", "tree"},
		validate: validatePartialRecord,
	},
	{
		name:     "check-partial verdict",
		fields:   []string{"check_evidence", "check_executed", "check_inherited", "oracle", "recorded_at", "schema", "state", "status", "tree"},
		validate: validateCheckPartialRecord,
	},
	{
		name:     "combined-partial verdict",
		fields:   []string{"check_evidence", "check_executed", "check_inherited", "executed", "oracle", "recorded_at", "schema", "skip_evidence", "skipped", "state", "status", "tree"},
		validate: validateCombinedPartialRecord,
	},
	{
		name:     "pending",
		fields:   []string{"oracle", "owner_pid", "schema", "started_at", "state", "tree"},
		validate: validatePendingRecord,
	},
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
