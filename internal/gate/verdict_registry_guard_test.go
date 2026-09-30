package gate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

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
		{name: "legacy_partial_invalid", state: Invalid, make: inspectLegacyPartialRecord},
		{name: "legacy_check_partial_invalid", state: Invalid, make: inspectLegacyCheckPartialRecord},
		{name: "legacy_combined_partial_invalid", state: Invalid, make: inspectLegacyCombinedPartialRecord},
		{name: "mixed_class_invalid", state: Invalid, make: inspectMixedClassRecord},
		{name: "full_invalid_status", state: Invalid, make: inspectFullInvalidStatusRecord},
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

// The partial classes are retired, and nothing writes them. A legacy on-disk record in any
// of their field sets must read as an invalid cache, so it can never pass as a full green.
func inspectLegacyPartialRecord(tree, oracle string, now time.Time) []byte {
	return inspectLegacyRecord(tree, oracle, now, legacyPartitionFields(now))
}

func inspectLegacyCheckPartialRecord(tree, oracle string, now time.Time) []byte {
	return inspectLegacyRecord(tree, oracle, now, legacyCheckPartitionFields(now))
}

func inspectLegacyCombinedPartialRecord(tree, oracle string, now time.Time) []byte {
	fields := legacyPartitionFields(now)
	for name, value := range legacyCheckPartitionFields(now) {
		fields[name] = value
	}
	return inspectLegacyRecord(tree, oracle, now, fields)
}

func legacyPartitionFields(now time.Time) map[string]any {
	return map[string]any{
		"executed": []string{"conformance"},
		"skipped":  []string{"build"},
		"skip_evidence": map[string]any{
			"build": map[string]string{"identity": strings.Repeat("c", 64), "authored_at": now.Add(-time.Minute).Format(time.RFC3339)},
		},
	}
}

func legacyCheckPartitionFields(now time.Time) map[string]any {
	return map[string]any{
		"check_executed":  []string{"conformance-meta"},
		"check_inherited": []string{"line-routing"},
		"check_evidence": map[string]any{
			"line-routing": map[string]string{"identity": strings.Repeat("d", 64), "authored_at": now.Add(-time.Minute).Format(time.RFC3339)},
		},
	}
}

func inspectLegacyRecord(tree, oracle string, now time.Time, extra map[string]any) []byte {
	record := map[string]any{
		"schema": verdictSchema, "state": "ready", "status": "green", "tree": tree, "oracle": oracle,
		"recorded_at": now.Add(-time.Minute).Format(time.RFC3339),
	}
	for name, value := range extra {
		record[name] = value
	}
	return inspectJSON(record)
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
