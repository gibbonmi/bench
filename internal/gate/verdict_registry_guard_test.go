package gate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

func TestVerdictRecordClassRegistryMatchesExpectation(t *testing.T) {
	want := []struct {
		name   string
		fields []string
	}{
		{
			name:   "full verdict",
			fields: []string{"oracle", "recorded_at", "schema", "state", "status", "tree"},
		},
		{
			name:   "partial verdict",
			fields: []string{"executed", "oracle", "recorded_at", "schema", "skip_evidence", "skipped", "state", "status", "tree"},
		},
		{
			name:   "check-partial verdict",
			fields: []string{"check_evidence", "check_executed", "check_inherited", "oracle", "recorded_at", "schema", "state", "status", "tree"},
		},
		{
			name:   "combined-partial verdict",
			fields: []string{"check_evidence", "check_executed", "check_inherited", "executed", "oracle", "recorded_at", "schema", "skip_evidence", "skipped", "state", "status", "tree"},
		},
		{
			name:   "pending",
			fields: []string{"oracle", "owner_pid", "schema", "started_at", "state", "tree"},
		},
		{
			name:   "lane record",
			fields: []string{"lane", "outcome", "recorded_at", "run_binary", "schema", "tree"},
		},
	}

	if len(verdictRecordClasses) != len(want) {
		t.Fatalf("record-class count = %d, want %d", len(verdictRecordClasses), len(want))
	}
	for i, row := range want {
		got := verdictRecordClasses[i]
		if got.name != row.name || !slices.Equal(got.fields, row.fields) {
			t.Errorf("record class %d = (%q, %v), want (%q, %v)", i, got.name, got.fields, row.name, row.fields)
		}
		if got.validate == nil {
			t.Errorf("record class %d has no validator", i)
		}
	}
}

func TestVerdictRecordClassesAtInspect(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	for _, tc := range []struct {
		name  string
		state State
		make  func(string, string, time.Time) []byte
	}{
		{name: "full_ready", state: Ready, make: inspectFullRecord},
		{name: "partial_ready", state: Ready, make: inspectPartialRecord},
		{name: "check_partial_ready", state: Ready, make: inspectCheckPartialRecord},
		{name: "combined_partial_ready", state: Ready, make: inspectCombinedPartialRecord},
		{name: "mixed_class_invalid", state: Invalid, make: inspectMixedClassRecord},
		{name: "full_invalid_status", state: Invalid, make: inspectFullInvalidStatusRecord},
		{name: "partial_executed_skipped_overlap_invalid", state: Invalid, make: inspectPartialOverlapRecord},
		{name: "check_partial_missing_evidence_invalid", state: Invalid, make: inspectCheckMissingEvidenceRecord},
		{name: "combined_partial_executed_skipped_overlap_invalid", state: Invalid, make: inspectCombinedOverlapRecord},
		{name: "pending_owner_pid_zero_invalid", state: Invalid, make: inspectPendingOwnerZeroRecord},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := outcomeFixture(t)
			subject, err := buildSubject(root)
			if err != nil {
				t.Fatal(err)
			}
			writeInspectRecord(t, root, tc.make(subject.Tree, subject.Oracle, now))
			if got := Inspect(root); got.State != tc.state {
				t.Fatalf("Inspect state = %s, want %s", got.State, tc.state)
			}
		})
	}
}

func writeInspectRecord(t *testing.T, root string, data []byte) {
	t.Helper()
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	if err := os.WriteFile(filepath.Join(gitdir, benchgit.GateCacheFile), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

func inspectReadyRecord(tree, oracle string, now time.Time) verdictRecord {
	return verdictRecord{
		Schema: verdictSchema, State: Ready, Status: "green", Tree: tree, Oracle: oracle,
		RecordedAt: now.Add(-time.Minute).Format(time.RFC3339),
	}
}

func inspectFullRecord(tree, oracle string, now time.Time) []byte {
	return inspectJSON(inspectReadyRecord(tree, oracle, now))
}

func inspectPartialRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	record.Executed = []string{"conformance", "conformance-suite"}
	record.Skipped = []string{"build", "vet"}
	record.SkipEvidence = map[string]skipEvidence{
		"build": {Seal: strings.Repeat("b", 64)},
		"vet":   {Identity: strings.Repeat("c", 64), AuthoredAt: now.Add(-time.Minute).Format(time.RFC3339)},
	}
	return inspectJSON(record)
}

func inspectCheckPartialRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	addCheckPartition(&record)
	record.CheckEvidence = map[string]skipEvidence{
		record.CheckInherited[0]: {Identity: strings.Repeat("d", 64), AuthoredAt: now.Add(-time.Minute).Format(time.RFC3339)},
	}
	return inspectJSON(record)
}

func inspectCombinedPartialRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	record.Executed = []string{"conformance", "conformance-suite"}
	record.Skipped = []string{"build", "vet"}
	record.SkipEvidence = map[string]skipEvidence{
		"build": {Seal: strings.Repeat("b", 64)},
		"vet":   {Identity: strings.Repeat("c", 64), AuthoredAt: now.Add(-time.Minute).Format(time.RFC3339)},
	}
	addCheckPartition(&record)
	record.CheckEvidence = map[string]skipEvidence{
		record.CheckInherited[0]: {Identity: strings.Repeat("d", 64), AuthoredAt: now.Add(-time.Minute).Format(time.RFC3339)},
	}
	return inspectJSON(record)
}

func addCheckPartition(record *verdictRecord) {
	for _, check := range registry.Checks {
		if !check.RunsAt(registry.Dev) {
			continue
		}
		if check.Meta || len(record.CheckInherited) > 0 {
			record.CheckExecuted = append(record.CheckExecuted, check.Name)
		} else {
			record.CheckInherited = []string{check.Name}
		}
	}
}

func inspectMixedClassRecord(tree, oracle string, now time.Time) []byte {
	return inspectJSON(map[string]any{
		"schema": 1, "state": "ready", "status": "green", "tree": tree, "oracle": oracle,
		"recorded_at": now.Add(-time.Minute).Format(time.RFC3339), "executed": []string{"build"},
	})
}

func inspectFullInvalidStatusRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	record.Status = "bogus"
	return inspectJSON(record)
}

func inspectPartialOverlapRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	record.Executed = []string{"build"}
	record.Skipped = []string{"build"}
	record.SkipEvidence = map[string]skipEvidence{"build": {Seal: strings.Repeat("b", 64)}}
	return inspectJSON(record)
}

func inspectCheckMissingEvidenceRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	addCheckPartition(&record)
	record.CheckEvidence = map[string]skipEvidence{
		"not-inherited": {Identity: strings.Repeat("d", 64), AuthoredAt: now.Add(-time.Minute).Format(time.RFC3339)},
	}
	return inspectJSON(record)
}

func inspectCombinedOverlapRecord(tree, oracle string, now time.Time) []byte {
	record := inspectReadyRecord(tree, oracle, now)
	record.Executed = []string{"build"}
	record.Skipped = []string{"build"}
	record.SkipEvidence = map[string]skipEvidence{"build": {Seal: strings.Repeat("b", 64)}}
	addCheckPartition(&record)
	record.CheckEvidence = map[string]skipEvidence{
		record.CheckInherited[0]: {Identity: strings.Repeat("d", 64), AuthoredAt: now.Add(-time.Minute).Format(time.RFC3339)},
	}
	return inspectJSON(record)
}

func inspectPendingOwnerZeroRecord(tree, oracle string, now time.Time) []byte {
	return inspectJSON(map[string]any{
		"schema": 1, "state": "pending", "tree": tree, "oracle": oracle,
		"started_at": now.Add(-time.Minute).Format(time.RFC3339), "owner_pid": 0,
	})
}

func inspectJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}
