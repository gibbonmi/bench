package gate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

const (
	reducedTestTree     = "0123456789abcdef0123456789abcdef01234567"
	reducedTestAncestor = "89abcdef0123456789abcdef0123456789abcdef"
)

var reducedTestOracle = strings.Repeat("a", 64)

func reducedTestRecord(now time.Time) verdictRecord {
	return verdictRecord{
		Schema:             verdictSchema,
		State:              Ready,
		Status:             "green",
		Tree:               reducedTestTree,
		Oracle:             reducedTestOracle,
		RecordedAt:         now.Format(time.RFC3339),
		Reduced:            true,
		Phases:             ReducedScope().IncludedPhases(),
		Ancestor:           reducedTestAncestor,
		AncestorRecordedAt: now.Add(-30 * time.Minute).Format(time.RFC3339),
	}
}

func fullTestRecord(now time.Time) verdictRecord {
	return verdictRecord{
		Schema:     verdictSchema,
		State:      Ready,
		Status:     "green",
		Tree:       reducedTestTree,
		Oracle:     reducedTestOracle,
		RecordedAt: now.Format(time.RFC3339),
	}
}

// roundTripReducedVerdict writes a record through the durable writer and reads it back
// through the loader, so a field the writer drops or the loader refuses is a failure of
// the pair rather than of either half read in isolation.
func roundTripReducedVerdict(t *testing.T, rec verdictRecord, now time.Time) ([]byte, loadedVerdict) {
	t.Helper()
	dir := t.TempDir()
	if err := durableReplace(dir, rec); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, benchgit.GateCacheFile)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data, loadVerdict(path, now)
}

// writeReducedTestCache installs raw bytes at the cache's required mode, so a case that
// hand-builds a record is refused for its contents rather than for its metadata.
func writeReducedTestCache(t *testing.T, dir string, body []byte) string {
	t.Helper()
	path := filepath.Join(dir, benchgit.GateCacheFile)
	if err := os.WriteFile(path, append(body, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func reducedTestObject(t *testing.T, rec verdictRecord) map[string]any {
	t.Helper()
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatal(err)
	}
	return object
}

// [R15] A reduced verdict is a second record class, so the marker, the phases that
// actually ran, and the ancestor whose evidence the rest of the run inherited all have to
// survive the write-read pair — an inherited field that silently drops leaves a record
// claiming a full grading it never had. The ancestor's recorded time travels with it
// unchanged: it attributes the inherited evidence to the run that produced it, and a
// re-stamp would dress an ever-older full green as recent.
func TestReducedVerdictRecordsAncestor(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)

	t.Run("reduced record", func(t *testing.T) {
		want := reducedTestRecord(now)
		data, loaded := roundTripReducedVerdict(t, want, now)
		if loaded.state != Ready || loaded.reason != "" {
			t.Fatalf("loaded reduced verdict = %s/%q, want ready with no reason (bytes %q)", loaded.state, loaded.reason, data)
		}
		if !reflect.DeepEqual(loaded.record, want) {
			t.Fatalf("reduced round-trip = %+v, want %+v (bytes %q)", loaded.record, want, data)
		}
		if len(loaded.record.Phases) == 0 {
			t.Fatalf("reduced round-trip kept no executed phases: %q", data)
		}
	})

	t.Run("full record unchanged", func(t *testing.T) {
		want := fullTestRecord(now)
		data, loaded := roundTripReducedVerdict(t, want, now)
		wantBytes := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":%q}`+"\n",
			want.Tree, want.Oracle, want.RecordedAt)
		if string(data) != wantBytes {
			t.Fatalf("full-shape bytes = %q, want %q — the reduced class leaked into the existing shape", data, wantBytes)
		}
		if loaded.state != Ready || loaded.reason != "" || !reflect.DeepEqual(loaded.record, want) {
			t.Fatalf("full round-trip = %s/%q %+v, want ready with %+v", loaded.state, loaded.reason, loaded.record, want)
		}
	})

	// A reduced verdict answers only for the phases that could observe its changeset, so
	// it is readable evidence but never the whole-tree green a reuse would credit.
	t.Run("not a reusable whole-tree green", func(t *testing.T) {
		root := gateTestRepo(t, "#!/usr/bin/env bash\nexit 0\n", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`)
		if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
			t.Fatalf("seed execution = %+v, want green", got)
		}
		plan := mustSubject(t, root)
		seeded := time.Now().UTC().Truncate(time.Second)
		reduced := reducedTestRecord(seeded)
		reduced.Tree, reduced.Oracle = plan.Tree, plan.Oracle
		if err := durableReplace(filepath.Dir(cachePath(t, root)), reduced); err != nil {
			t.Fatal(err)
		}
		got := inspectAt(root, seeded)
		if got.State != Ready || !got.Reduced || got.ReusableGreen || got.Reason != "reduced verdict" {
			t.Fatalf("reduced inspection = %+v, want a readable ready record that is not a reusable green", got)
		}
	})
}

// [R16] The loader validates an exact field set per state, and the second class doubles
// the ways a record can be half-written. Every case here carries a coherent subset of one
// shape plus a fragment of the other: the loader has no sound way to guess which class was
// meant, and guessing credits phases nobody ran, so each is refused.
func TestVerdictLoaderRejectsMixedShape(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	reduced := reducedTestRecord(now)
	full := fullTestRecord(now)

	cases := []struct {
		name   string
		record verdictRecord
		mutate func(map[string]any)
		want   State
	}{
		{name: "intact reduced record", record: reduced, want: Ready},
		{name: "intact full record", record: full, want: Ready},
		{
			name:   "ancestor without the reduced marker",
			record: reduced,
			mutate: func(o map[string]any) { delete(o, "reduced") },
			want:   Invalid,
		},
		{
			name:   "reduced marker denied by its own value",
			record: reduced,
			mutate: func(o map[string]any) { o["reduced"] = false },
			want:   Invalid,
		},
		{
			name:   "reduced marker without an ancestor",
			record: reduced,
			mutate: func(o map[string]any) { delete(o, "ancestor"); delete(o, "ancestor_recorded_at") },
			want:   Invalid,
		},
		{
			name:   "reduced marker without the executed phases",
			record: reduced,
			mutate: func(o map[string]any) { delete(o, "phases") },
			want:   Invalid,
		},
		{
			name:   "ancestor without its recorded time",
			record: reduced,
			mutate: func(o map[string]any) { delete(o, "ancestor_recorded_at") },
			want:   Invalid,
		},
		{
			name:   "empty executed phase list",
			record: reduced,
			mutate: func(o map[string]any) { o["phases"] = []any{} },
			want:   Invalid,
		},
		{
			name:   "inherited time re-stamped past the record's own",
			record: reduced,
			mutate: func(o map[string]any) { o["ancestor_recorded_at"] = now.Add(time.Minute).Format(time.RFC3339) },
			want:   Invalid,
		},
		{
			name:   "full shape carrying an ancestor",
			record: full,
			mutate: func(o map[string]any) { o["ancestor"] = reducedTestAncestor },
			want:   Invalid,
		},
		{
			name:   "full shape carrying the reduced marker",
			record: full,
			mutate: func(o map[string]any) { o["reduced"] = true },
			want:   Invalid,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := reducedTestObject(t, tc.record)
			if tc.mutate != nil {
				tc.mutate(object)
			}
			body, err := json.Marshal(object)
			if err != nil {
				t.Fatal(err)
			}
			got := loadVerdict(writeReducedTestCache(t, t.TempDir(), body), now)
			if got.state != tc.want {
				t.Fatalf("loaded %q = %s/%q, want %s (bytes %q)", tc.name, got.state, got.reason, tc.want, body)
			}
			if tc.want == Invalid && got.reason != "invalid cache record" {
				t.Fatalf("rejection reason for %q = %q, want %q", tc.name, got.reason, "invalid cache record")
			}
		})
	}

	// The reduced class belongs to a completed verdict; a pending record reaching for it
	// would claim inherited evidence for a run still in flight.
	t.Run("pending record carrying the reduced marker", func(t *testing.T) {
		object := reducedTestObject(t, verdictRecord{
			Schema: verdictSchema, State: Pending, Tree: reducedTestTree, Oracle: reducedTestOracle,
			StartedAt: now.Format(time.RFC3339), OwnerPID: 4321, Reduced: true,
		})
		body, err := json.Marshal(object)
		if err != nil {
			t.Fatal(err)
		}
		if got := loadVerdict(writeReducedTestCache(t, t.TempDir(), body), now); got.state != Invalid {
			t.Fatalf("pending record with the reduced marker = %s/%q, want invalid", got.state, got.reason)
		}
	})
}
