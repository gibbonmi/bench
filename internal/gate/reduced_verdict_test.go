package gate

import (
	"bytes"
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

var (
	partialTestIdentity = strings.Repeat("b", 64)
	partialTestSeal     = strings.Repeat("c", 64)
)

// partialTestRecord carries both evidence forms at once — an ancestor slot for a scoped
// component and a reused seal for build — so a case that drops one form is refused by the
// same record every other case starts from.
func partialTestRecord(now time.Time) verdictRecord {
	return verdictRecord{
		Schema:     verdictSchema,
		State:      Ready,
		Status:     "green",
		Tree:       reducedTestTree,
		Oracle:     reducedTestOracle,
		RecordedAt: now.Format(time.RFC3339),
		Executed:   []string{"conformance", "conformance-suite"},
		Skipped:    []string{"build", "vet"},
		SkipEvidence: map[string]skipEvidence{
			"build": {Seal: partialTestSeal},
			"vet":   {Identity: partialTestIdentity, AuthoredAt: now.Add(-90 * time.Minute).Format(time.RFC3339)},
		},
	}
}

// partialTestEvidence reaches into a marshalled record's evidence map so a case can bend one
// entry without restating the whole record.
func partialTestEvidence(object map[string]any, component string) map[string]any {
	return object["skip_evidence"].(map[string]any)[component].(map[string]any)
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

// [PC2a] A partial verdict is a third record class, and the evidence is the whole of what
// makes it readable: a record naming which components were skipped without saying what
// covered each one claims a grading nobody can trace. The executed set, both evidence forms,
// and the projection consumers read all have to survive the write-read pair, and the record
// has to re-serialize to the same bytes — reuse compares records byte for byte, so a
// partition that round-trips into a second spelling would read as a second verdict.
func TestPartialVerdictRoundTrips(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	want := partialTestRecord(now)
	data, loaded := roundTripReducedVerdict(t, want, now)
	if loaded.state != Ready || loaded.reason != "" {
		t.Fatalf("loaded partial verdict = %s/%q, want ready with no reason (bytes %q)", loaded.state, loaded.reason, data)
	}
	if !reflect.DeepEqual(loaded.record, want) {
		t.Fatalf("partial round-trip = %+v, want %+v (bytes %q)", loaded.record, want, data)
	}
	again, _ := roundTripReducedVerdict(t, loaded.record, now)
	if !bytes.Equal(again, data) {
		t.Fatalf("re-serialized partial record = %q, want %q", again, data)
	}

	got := loaded.record.partition()
	if got == nil {
		t.Fatalf("partial record projected no partition (bytes %q)", data)
	}
	if !reflect.DeepEqual(got.Executed, want.Executed) {
		t.Fatalf("projected executed set = %q, want %q", got.Executed, want.Executed)
	}
	wantSkips := []ComponentSkip{
		{Component: "build", Seal: partialTestSeal},
		{Component: "vet", Identity: partialTestIdentity, AuthoredAt: now.Add(-90 * time.Minute)},
	}
	if !reflect.DeepEqual(got.Skipped, wantSkips) {
		t.Fatalf("projected skips = %+v, want %+v", got.Skipped, wantSkips)
	}
	if fullTestRecord(now).partition() != nil || reducedTestRecord(now).partition() != nil {
		t.Fatal("a full or reduced record projected a partition, so a consumer cannot read the nil as whole-tree grading")
	}
}

// [PC2b] The skipped set and the evidence map are cross-checked in both directions. A skip
// with no evidence credits a component nobody graded; evidence naming a component the record
// did not skip proves nothing the record claims. Each entry is then graded against exactly
// one evidence form: an identity that is not a content address addresses no slot, and an
// entry reaching for both forms describes no evidence at all.
func TestPartialVerdictRequiresEvidencePerSkip(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		mutate func(map[string]any)
		want   State
	}{
		{name: "intact partial record", want: Ready},
		{
			name:   "skip with no evidence entry",
			mutate: func(o map[string]any) { delete(o["skip_evidence"].(map[string]any), "vet") },
			want:   Invalid,
		},
		{
			name: "evidence for a component that ran",
			mutate: func(o map[string]any) {
				evidence := o["skip_evidence"].(map[string]any)
				evidence["conformance"] = evidence["vet"]
				delete(evidence, "vet")
			},
			want: Invalid,
		},
		{
			name: "evidence beyond the skipped set",
			mutate: func(o map[string]any) {
				o["skip_evidence"].(map[string]any)["conformance"] = map[string]any{"seal": partialTestSeal}
			},
			want: Invalid,
		},
		{
			name:   "empty evidence entry",
			mutate: func(o map[string]any) { o["skip_evidence"].(map[string]any)["vet"] = map[string]any{} },
			want:   Invalid,
		},
		{
			name:   "ancestor identity that is not a content address",
			mutate: func(o map[string]any) { partialTestEvidence(o, "vet")["identity"] = "the-vet-component" },
			want:   Invalid,
		},
		{
			name: "ancestor identity spelled in upper case",
			mutate: func(o map[string]any) {
				partialTestEvidence(o, "vet")["identity"] = strings.ToUpper(partialTestIdentity)
			},
			want: Invalid,
		},
		{
			name:   "ancestor evidence without its authored time",
			mutate: func(o map[string]any) { delete(partialTestEvidence(o, "vet"), "authored_at") },
			want:   Invalid,
		},
		{
			name: "slot authored after the run that read it",
			mutate: func(o map[string]any) {
				partialTestEvidence(o, "vet")["authored_at"] = now.Add(time.Minute).Format(time.RFC3339)
			},
			want: Invalid,
		},
		{
			name:   "ancestor evidence also carrying a seal",
			mutate: func(o map[string]any) { partialTestEvidence(o, "vet")["seal"] = partialTestSeal },
			want:   Invalid,
		},
		{
			name:   "seal evidence also carrying an identity",
			mutate: func(o map[string]any) { partialTestEvidence(o, "build")["identity"] = partialTestIdentity },
			want:   Invalid,
		},
		{
			name:   "seal that is not a content address",
			mutate: func(o map[string]any) { partialTestEvidence(o, "build")["seal"] = "built-recently" },
			want:   Invalid,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertPartialRecordLoads(t, now, tc.mutate, tc.want)
		})
	}
}

// [PC2b] The partition's two halves are a coherent split of the run: every name is a real
// component named once, a component cannot be both executed and skipped, and one partition
// has exactly one spelling so two records over one run cannot differ byte for byte.
func TestPartialVerdictRequiresACoherentComponentSet(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"empty executed set", func(o map[string]any) { o["executed"] = []any{} }},
		{"empty skipped set", func(o map[string]any) { o["skipped"] = []any{}; o["skip_evidence"] = map[string]any{} }},
		{"unnamed executed component", func(o map[string]any) { o["executed"] = []any{"", "conformance"} }},
		{"duplicated executed component", func(o map[string]any) { o["executed"] = []any{"conformance", "conformance"} }},
		{"unsorted executed set", func(o map[string]any) { o["executed"] = []any{"conformance-suite", "conformance"} }},
		{"unsorted skipped set", func(o map[string]any) { o["skipped"] = []any{"vet", "build"} }},
		{
			name: "component both executed and skipped",
			mutate: func(o map[string]any) {
				o["executed"] = []any{"conformance", "conformance-suite", "vet"}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertPartialRecordLoads(t, now, tc.mutate, Invalid)
		})
	}
}

// [PS23] The three ready classes are alternatives, never a spectrum. A record carrying
// fields of two of them names no class the loader can resolve, and resolving it by guess
// would credit one class's evidence for the other's claim — so it is refused for the shape
// it holds rather than read as whichever class was checked first.
func TestMixedClassRecordRefuses(t *testing.T) {
	now := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	partial := partialTestRecord(now)
	reduced := reducedTestRecord(now)
	cases := []struct {
		name   string
		record verdictRecord
		mutate func(map[string]any)
	}{
		{
			name:   "partial shape carrying an ancestor",
			record: partial,
			mutate: func(o map[string]any) { o["ancestor"] = reducedTestAncestor },
		},
		{
			name:   "partial shape carrying the reduced marker",
			record: partial,
			mutate: func(o map[string]any) { o["reduced"] = true },
		},
		{
			name:   "reduced shape carrying a partition",
			record: reduced,
			mutate: func(o map[string]any) {
				o["executed"] = partial.Executed
				o["skipped"] = partial.Skipped
				o["skip_evidence"] = reducedTestObject(t, partial)["skip_evidence"]
			},
		},
		{
			name:   "full shape carrying an executed set",
			record: fullTestRecord(now),
			mutate: func(o map[string]any) { o["executed"] = partial.Executed },
		},
		{
			name:   "every ready field of every class at once",
			record: reduced,
			mutate: func(o map[string]any) {
				for name, value := range reducedTestObject(t, partial) {
					o[name] = value
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := reducedTestObject(t, tc.record)
			tc.mutate(object)
			body, err := json.Marshal(object)
			if err != nil {
				t.Fatal(err)
			}
			got := loadVerdict(writeReducedTestCache(t, t.TempDir(), body), now)
			if got.state != Invalid || got.reason != "invalid cache record" {
				t.Fatalf("loaded %q = %s/%q, want an invalid cache record (bytes %q)", tc.name, got.state, got.reason, body)
			}
		})
	}
}

// assertPartialRecordLoads writes a bent partial record at the cache's required mode and
// grades the loader's verdict on it, so a case is refused for its contents rather than for
// the file it arrived in.
func assertPartialRecordLoads(t *testing.T, now time.Time, mutate func(map[string]any), want State) {
	t.Helper()
	object := reducedTestObject(t, partialTestRecord(now))
	if mutate != nil {
		mutate(object)
	}
	body, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	got := loadVerdict(writeReducedTestCache(t, t.TempDir(), body), now)
	if got.state != want {
		t.Fatalf("loaded partial record = %s/%q, want %s (bytes %q)", got.state, got.reason, want, body)
	}
	if want == Invalid && got.reason != "invalid cache record" {
		t.Fatalf("rejection reason = %q, want %q", got.reason, "invalid cache record")
	}
}
