package settlepolicy

import "testing"

// TestSettlePolicyAnswersTheCaptureRule drives the settle decision with literal stage
// records and observes the verdict of each path. It builds no repository.
// (Coverage rows LS20, LS21.)
func TestSettlePolicyAnswersTheCaptureRule(t *testing.T) {
	blob := func(path string, stage int) StageRecord {
		return StageRecord{Mode: "100644", OID: "0123456789abcdef0123456789abcdef01234567", Stage: stage, Path: path}
	}
	for _, tc := range []struct {
		name     string
		records  []StageRecord
		want     []Verdict
		reason   string
		refusing []string
	}{
		{
			name:    "the handoff takes the source stage",
			records: []StageRecord{blob("capture/session-handoff.md", 2), blob("capture/session-handoff.md", 3)},
			want: []Verdict{{
				Path: "capture/session-handoff.md", Kind: VerdictSide, Side: SideSource,
				Record: blob("capture/session-handoff.md", 3),
			}},
		},
		{
			name:    "the two journals union",
			records: []StageRecord{blob("capture/learnings.md", 2), blob("capture/learnings.md", 3), blob("capture/IDEAS.md", 3)},
			want: []Verdict{
				{Path: "capture/learnings.md", Kind: VerdictUnion, Side: SideUnion, Stages: map[int]StageRecord{
					2: blob("capture/learnings.md", 2), 3: blob("capture/learnings.md", 3),
				}},
				{Path: "capture/IDEAS.md", Kind: VerdictUnion, Side: SideUnion, Stages: map[int]StageRecord{
					3: blob("capture/IDEAS.md", 3),
				}},
			},
		},
		{
			name:    "another capture path takes the destination stage",
			records: []StageRecord{blob("capture/agent-performance/claude-models.md", 2), blob("capture/agent-performance/claude-models.md", 3)},
			want: []Verdict{{
				Path: "capture/agent-performance/claude-models.md", Kind: VerdictSide, Side: SideDestination,
				Record: blob("capture/agent-performance/claude-models.md", 2),
			}},
		},
		{
			name:    "a capture path holding a space settles by the table",
			records: []StageRecord{blob("capture/path with a space.md", 2), blob("capture/path with a space.md", 3)},
			want: []Verdict{{
				Path: "capture/path with a space.md", Kind: VerdictSide, Side: SideDestination,
				Record: blob("capture/path with a space.md", 2),
			}},
		},
		{
			name:    "an absent winning stage answers a removal",
			records: []StageRecord{blob("capture/agent-performance/claude-models.md", 1), blob("capture/agent-performance/claude-models.md", 3)},
			want: []Verdict{{
				Path: "capture/agent-performance/claude-models.md", Kind: VerdictRemove, Side: SideDestination,
			}},
		},
		{
			name:    "a union whose two sides are both absent answers a removal",
			records: []StageRecord{blob("capture/learnings.md", 1)},
			want: []Verdict{{
				Path: "capture/learnings.md", Kind: VerdictRemove, Side: SideUnion,
			}},
		},
		{
			// The two rows below engage the table through a trailing capture record, so
			// each observes the reason the earlier out-of-table record wins.
			name:     "the prefix boundary sits outside the table",
			records:  []StageRecord{blob("capture.md", 2), blob("capture.md", 3), blob("capture/notes.md", 2)},
			reason:   ReasonOutsideTable,
			refusing: []string{"capture.md"},
		},
		{
			name:     "the table check beats the mode check inside one record",
			records:  []StageRecord{{Mode: "120000", OID: "a", Stage: 3, Path: "capture.md"}, blob("capture/notes.md", 2)},
			reason:   ReasonOutsideTable,
			refusing: []string{"capture.md"},
		},
		{
			name:    "a conflict the table names nowhere engages no rule",
			records: []StageRecord{blob("named", 2), blob("named", 3)},
		},
		{
			name: "the earlier record wins the reason",
			records: []StageRecord{
				{Mode: "120000", OID: "a", Stage: 3, Path: "capture/learnings.md"},
				blob("capture.md", 2),
			},
			reason:   ReasonNonRegularMode,
			refusing: []string{"capture/learnings.md"},
		},
		{
			name:     "a non-regular mode under the table refuses",
			records:  []StageRecord{{Mode: "120000", OID: "a", Stage: 3, Path: "capture/learnings.md"}},
			reason:   ReasonNonRegularMode,
			refusing: []string{"capture/learnings.md"},
		},
		{
			name: "unequal regular modes on the two sides disagree",
			records: []StageRecord{
				blob("capture/session-handoff.md", 2),
				{Mode: "100755", OID: "b", Stage: 3, Path: "capture/session-handoff.md"},
			},
			reason:   ReasonModeDisagreement,
			refusing: []string{"capture/session-handoff.md"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Settle(tc.records)
			if got.Refusal.Reason != tc.reason {
				t.Fatalf("reason = %q, want %q", got.Refusal.Reason, tc.reason)
			}
			if tc.reason != "" {
				if len(got.Refusal.Paths) != len(tc.refusing) {
					t.Fatalf("refusing paths = %v, want %v", got.Refusal.Paths, tc.refusing)
				}
				for i, path := range tc.refusing {
					if got.Refusal.Paths[i] != path {
						t.Fatalf("refusing paths = %v, want %v", got.Refusal.Paths, tc.refusing)
					}
				}
				if len(got.Verdicts) != 0 {
					t.Fatalf("refusal carries verdicts %v", got.Verdicts)
				}
				return
			}
			if len(got.Verdicts) != len(tc.want) {
				t.Fatalf("verdicts = %+v, want %+v", got.Verdicts, tc.want)
			}
			for i, want := range tc.want {
				have := got.Verdicts[i]
				if have.Path != want.Path || have.Kind != want.Kind || have.Side != want.Side || have.Record != want.Record {
					t.Fatalf("verdict %d = %+v, want %+v", i, have, want)
				}
				if len(have.Stages) != len(want.Stages) {
					t.Fatalf("verdict %d stages = %+v, want %+v", i, have.Stages, want.Stages)
				}
				for stage, record := range want.Stages {
					if have.Stages[stage] != record {
						t.Fatalf("verdict %d stage %d = %+v, want %+v", i, stage, have.Stages[stage], record)
					}
				}
			}
		})
	}
}

// TestConflictKindClassifiesModeLists drives the mode-list classifier with literal
// mode lists. It builds no repository. (Coverage row LS23.)
func TestConflictKindClassifiesModeLists(t *testing.T) {
	for _, tc := range []struct {
		name  string
		modes []string
		want  string
	}{
		{name: "a gitlink beats the ordinary disagreement", modes: []string{"100644", "160000"}, want: "gitlink"},
		{name: "a symlink beats the ordinary disagreement", modes: []string{"100644", "120000"}, want: "symlink"},
		{name: "two unequal ordinary modes disagree", modes: []string{"100644", "100755"}, want: "mode"},
		{name: "one repeated mode is textual", modes: []string{"100644", "100644"}, want: "textual"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ConflictKind(tc.modes); got != tc.want {
				t.Fatalf("ConflictKind(%v) = %q, want %q", tc.modes, got, tc.want)
			}
		})
	}
}

// TestRefuseBuildsTheAdapterReason observes the refusal constructor the adapter uses
// after its own union text merge fails. (Coverage row LS20.)
func TestRefuseBuildsTheAdapterReason(t *testing.T) {
	got := Refuse(ReasonUnionContentNotText, []string{"capture/learnings.md"})
	if got.Reason != "union content not text" || len(got.Paths) != 1 || got.Paths[0] != "capture/learnings.md" {
		t.Fatalf("Refuse = %+v", got)
	}
}
