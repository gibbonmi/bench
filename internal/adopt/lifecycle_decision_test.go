package adopt

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func decision(path, kind, reason, action, fingerprint string, order int) LifecycleDecision {
	return LifecycleDecision{Path: path, Kind: kind, Reason: reason, ManifestAction: action, Fingerprint: fingerprint, Order: order}
}

// requirePlan compares a whole plan and names each differing decision.
func requirePlan(t *testing.T, got, want LifecyclePlan) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return
	}
	for i := 0; i < len(got.Decisions) || i < len(want.Decisions); i++ {
		var g, w LifecycleDecision
		if i < len(got.Decisions) {
			g = got.Decisions[i]
		}
		if i < len(want.Decisions) {
			w = want.Decisions[i]
		}
		if g != w {
			t.Errorf("decision %d = %+v, want %+v", i, g, w)
		}
	}
	t.Fatalf("operations = %v, want %v; refusal = %q", got.Operations, want.Operations, got.Refusal)
}

// linkObservations holds one link observation for each execution partition. The slice
// order is the caller order, and it differs from the path order on purpose.
func linkObservations() []LifecycleObservation {
	regular := func(o LifecycleObservation) LifecycleObservation {
		o.Present, o.FileKind, o.Perm, o.StagedKind = true, "regular", 0o644, "regular"
		if o.StagedPerm == 0 {
			o.StagedPerm = 0o644
		}
		return o
	}
	observations := []LifecycleObservation{
		{Path: "seed-present", PlanKind: "seed", Present: true, FileKind: "regular"},
		{Path: "absent", PlanKind: "file", Staged: "staged"},
		regular(LifecycleObservation{Path: "clean-old", PlanKind: "file", Recorded: "old", Observed: "old", Staged: "new"}),
		regular(LifecycleObservation{Path: "converged", PlanKind: "file", Recorded: "old", Observed: "new", Staged: "new"}),
		regular(LifecycleObservation{Path: "mode", PlanKind: "file", Recorded: "same", Observed: "same", Staged: "same", StagedPerm: 0o755}),
		regular(LifecycleObservation{Path: "strict-mode", PlanKind: "file", Recorded: "same", Observed: "same", Staged: "same", StagedPerm: 0o755, Strict: true}),
		{Path: "adapter-canonical", PlanKind: "adapter", Present: true, FileKind: "regular", Observed: "bytes", Staged: "link", StagedKind: "symlink",
			ResolvedEquivalent: true, CanonicalAdapter: true, SymlinkParent: true},
		regular(LifecycleObservation{Path: "modified", PlanKind: "file", Recorded: "old", Observed: "edit", Staged: "new"}),
		regular(LifecycleObservation{Path: "collision", PlanKind: "file", Observed: "new", Staged: "new"}),
		{Path: "seed-absent", PlanKind: "seed"},
		{Path: "seed-owned", PlanKind: "seed", Present: true, FileKind: "symlink", Recorded: "old"},
		{Path: "dropped-clean", Present: true, FileKind: "regular", Recorded: "h", Observed: "h"},
		{Path: "dropped-modified", Present: true, FileKind: "regular", Recorded: "h", Observed: "edit"},
		{Path: "dropped-absent", Recorded: "h"},
		{Path: "claude-modified", PlanKind: "claude", Present: true, FileKind: "regular", Recorded: "old", Observed: "edit", Staged: "current"},
		{Path: "claude-absent", PlanKind: "claude", Staged: "current"},
		{Path: "claude-known", PlanKind: "claude", Present: true, FileKind: "regular", KnownForm: true, Recorded: "old", Observed: "legacy", Staged: "current"},
		{Path: "claude-empty", PlanKind: "claude", Present: true, FileKind: "regular", Observed: "empty", Staged: "current"},
		{Path: "claude-special", PlanKind: "claude", Present: true, FileKind: "other", Staged: "current"},
	}
	for i := range observations {
		observations[i].Order = i
	}
	return observations
}

func TestCompleteLifecycleDecisions(t *testing.T) {
	t.Run("legacy", func(t *testing.T) {
		got := PlanLifecycle(LifecycleInput{
			Current:  []Asset{{"changed", "old"}, {"removed", "old"}, {"CLAUDE.md", "keep"}, {"same", "s"}, {"pinned", "a"}},
			Desired:  []Asset{{"added", "new"}, {"changed", "new"}, {"same", "s"}, {"pinned", "b"}},
			Preserve: []string{"CLAUDE.md", "pinned", "preserve-only"},
		})
		want := LifecyclePlan{
			Decisions: []LifecycleDecision{
				decision("CLAUDE.md", "preserve", "requested-preserve", "retain", "keep", 0),
				decision("added", "add", "", "record", "new", 0),
				decision("changed", "change", "", "record", "new", 0),
				decision("pinned", "change", "", "record", "b", 0),
				decision("removed", "remove", "", "omit", "", 0),
				decision("same", "preserve", "unchanged", "retain", "s", 0),
			},
			Operations: []Operation{{"add", "added"}, {"change", "changed"}, {"change", "pinned"}, {"remove", "removed"}},
		}
		requirePlan(t, got, want)
	})
	t.Run("link", func(t *testing.T) {
		got := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: linkObservations()}})
		want := LifecyclePlan{
			Decisions: []LifecycleDecision{
				decision("absent", "add", "", "record", "staged", 1),
				decision("adapter-canonical", "preserve", "converged", "record", "bytes", 6),
				decision("claude-absent", "add", "", "record", "current", 15),
				decision("claude-empty", "preserve", "project-owned", "omit", "", 17),
				decision("claude-known", "change", "", "record", "current", 16),
				decision("claude-modified", "preserve", "modified-managed", "retain", "old", 14),
				decision("claude-special", "preserve", "project-owned", "omit", "", 18),
				decision("clean-old", "change", "", "record", "new", 2),
				decision("collision", "preserve", "project-owned", "omit", "", 8),
				decision("converged", "preserve", "converged", "record", "new", 3),
				decision("dropped-absent", "preserve", "absent", "omit", "", 13),
				decision("dropped-clean", "remove", "", "omit", "", 11),
				decision("dropped-modified", "preserve", "kept-modified-removed", "retain", "h", 12),
				decision("mode", "change", "", "record", "same", 4),
				decision("modified", "preserve", "modified-managed", "retain", "old", 7),
				decision("seed-absent", "add", "", "omit", "", 9),
				decision("seed-owned", "preserve", "seed-owned", "retain", "old", 10),
				decision("seed-present", "preserve", "seed-owned", "omit", "", 0),
				decision("strict-mode", "preserve", "modified-managed", "retain", "same", 5),
			},
			Operations: []Operation{{"add", "absent"}, {"add", "claude-absent"}, {"change", "claude-known"}, {"change", "clean-old"},
				{"remove", "dropped-clean"}, {"change", "mode"}, {"add", "seed-absent"}},
		}
		requirePlan(t, got, want)
	})
	t.Run("unlink", func(t *testing.T) {
		got := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "unlink", Observations: []LifecycleObservation{
			{Path: "clean", Order: 0, Present: true, FileKind: "regular", Recorded: "h", Observed: "h"},
			{Path: "edited", Order: 1, Present: true, FileKind: "regular", Recorded: "h", Observed: "x"},
			{Path: "../escape", Order: 2, Escapes: true, Present: true, FileKind: "regular", Recorded: "h", Observed: "h"},
			{Path: "linked/asset", Order: 3, SymlinkParent: true, Present: true, FileKind: "regular", Recorded: "h", Observed: "h"},
			{Path: "dir", Order: 4, Present: true, FileKind: "other", Recorded: "h"},
			{Path: "gone", Order: 5, Recorded: "h"},
		}}})
		want := LifecyclePlan{
			Decisions: []LifecycleDecision{
				decision("../escape", "preserve", "refused", "retain", "h", 2),
				decision("clean", "remove", "", "omit", "", 0),
				decision("dir", "preserve", "refused", "retain", "h", 4),
				decision("edited", "preserve", "modified", "retain", "h", 1),
				decision("gone", "preserve", "absent", "omit", "", 5),
				decision("linked/asset", "preserve", "refused", "retain", "h", 3),
			},
			Operations: []Operation{{"remove", "clean"}},
		}
		requirePlan(t, got, want)
	})
	t.Run("present empty execution", func(t *testing.T) {
		current := []Asset{{"kept", "a"}}
		if got := PlanLifecycle(LifecycleInput{Current: current, Execution: &LifecycleInventory{}}); !reflect.DeepEqual(got, LifecyclePlan{}) {
			t.Fatalf("present empty execution = %#v, want an empty plan", got)
		}
		if got := PlanLifecycle(LifecycleInput{Current: current}); len(got.Decisions) != 1 {
			t.Fatalf("nil execution = %#v, want the legacy decision", got)
		}
	})
	for _, tc := range []struct {
		name, action, path string
		observations       []LifecycleObservation
	}{
		{name: "symlink parent write", action: "link", path: "linked/asset", observations: []LifecycleObservation{
			{Path: "linked/asset", PlanKind: "inline", SymlinkParent: true, Staged: "s"}}},
		{name: "foreign adapter", action: "link", path: ".claude/commands/x.md", observations: []LifecycleObservation{
			{Path: ".claude/commands/x.md", PlanKind: "adapter", Present: true, FileKind: "regular", Observed: "bytes", Staged: "link",
				StagedKind: "symlink", ResolvedEquivalent: true, SymlinkParent: true}}},
		{name: "seed symlink parent", action: "link", path: "linked/seed", observations: []LifecycleObservation{
			{Path: "linked/seed", PlanKind: "seed", Present: true, FileKind: "regular", SymlinkParent: true}}},
		{name: "dropped escape", action: "link", path: "../escape", observations: []LifecycleObservation{
			{Path: "../escape", Escapes: true, Present: true, Recorded: "h", Observed: "h"}}},
		{name: "dropped special", action: "link", path: "dir", observations: []LifecycleObservation{
			{Path: "dir", Present: true, FileKind: "other", Recorded: "h"}}},
		{name: "duplicate", action: "unlink", path: "twice", observations: []LifecycleObservation{
			{Path: "twice", Recorded: "h"}, {Path: "twice", Recorded: "h", Order: 1}}},
		{name: "unknown action", action: "relink", path: "relink", observations: []LifecycleObservation{{Path: "x"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: tc.action, Observations: tc.observations}})
			if !strings.Contains(got.Refusal, tc.path) || got.Decisions != nil || got.Operations != nil {
				t.Fatalf("PlanLifecycle() = %#v, want a refusal naming %q", got, tc.path)
			}
		})
	}
}

func TestLifecycleDecisionOrder(t *testing.T) {
	reversed := func(observations []LifecycleObservation) []LifecycleObservation {
		out := make([]LifecycleObservation, len(observations))
		for i, o := range observations {
			out[len(out)-1-i] = o
		}
		return out
	}
	t.Run("execution", func(t *testing.T) {
		observations := linkObservations()
		forward := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: observations}})
		backward := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: reversed(observations)}})
		if len(forward.Decisions) != len(observations) || !reflect.DeepEqual(forward, backward) {
			t.Fatalf("reordered observations changed the plan:\n%#v\n%#v", forward, backward)
		}
		if !sort.SliceIsSorted(forward.Decisions, func(i, j int) bool { return forward.Decisions[i].Path < forward.Decisions[j].Path }) {
			t.Fatalf("decisions are not in path order: %#v", forward.Decisions)
		}
	})
	t.Run("legacy", func(t *testing.T) {
		current := []Asset{{"b", "old"}, {"a", "old"}, {"c", "same"}}
		desired := []Asset{{"c", "same"}, {"d", "new"}, {"b", "new"}}
		forward := PlanLifecycle(LifecycleInput{Current: current, Desired: desired})
		backward := PlanLifecycle(LifecycleInput{Current: []Asset{current[2], current[1], current[0]}, Desired: []Asset{desired[2], desired[1], desired[0]}})
		paths := []string{}
		for _, d := range forward.Decisions {
			paths = append(paths, d.Path)
		}
		if !reflect.DeepEqual(paths, []string{"a", "b", "c", "d"}) || !reflect.DeepEqual(forward, backward) {
			t.Fatalf("reordered inventories changed the plan:\n%#v\n%#v", forward, backward)
		}
	})
	t.Run("first caller refusal", func(t *testing.T) {
		observations := []LifecycleObservation{
			{Path: "a-late", Order: 7, PlanKind: "inline", SymlinkParent: true},
			{Path: "z-early", Order: 2, PlanKind: "inline", SymlinkParent: true},
		}
		forward := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: observations}})
		backward := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: reversed(observations)}})
		if !strings.Contains(forward.Refusal, "z-early") || !reflect.DeepEqual(forward, backward) {
			t.Fatalf("refusals = %q then %q, want the earliest caller refusal for z-early", forward.Refusal, backward.Refusal)
		}
	})
}
