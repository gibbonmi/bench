package reclaimpolicy

import (
	"reflect"
	"strings"
	"testing"
)

// deadChild is a child whose pointer provably names an absent repository —
// the only shape the predicate may count toward a reclaim.
func deadChild(name, target string) ChildFacts {
	return ChildFacts{Name: name, Shape: ShapeDir, Pointer: PointerFacts{
		Shape: ShapeRegular, Body: "gitdir: " + target + "\n", TargetExistence: ExistenceAbsent}}
}

// TestClassifyKeyTable is the RP1 typed reclaim table. Each case is one named
// partition of the reclaim edge inventory: protected (current repository,
// live target), uncertain (unreadable key, child, pointer, or target; the
// unfilled zero-value existence), hostile (symlink at every level, FIFO
// pointer, repository directory, malformed and relative pointers, mixed
// live-and-dead key), and provably dead (empty key, all-dead children). A
// filesystem-only shortcut misclassifies a live, uncertain, drifted, or
// hostile key in one of these partitions by name.
func TestClassifyKeyTable(t *testing.T) {
	cases := []struct {
		partition string
		facts     KeyFacts
		verdict   string
		reason    string
		targets   []string
	}{
		{"protected: current repository key", KeyFacts{Name: "cur", Current: true},
			VerdictRetain, "key belongs to the current repository", nil},
		{"protected: live gitdir target", KeyFacts{Name: "live", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, Body: "gitdir: /repo/.git/worktrees/wt\n", TargetExistence: ExistencePresent}}}},
			VerdictRetain, "child wt gitdir: target exists", nil},
		{"uncertain: unreadable key", KeyFacts{Name: "k", Shape: ShapeUnreadable, ShapeErr: "permission denied"},
			VerdictRetain, "key cannot be read: permission denied", nil},
		{"uncertain: missing key", KeyFacts{Name: "k", Shape: ShapeMissing, ShapeErr: "no such file"},
			VerdictRetain, "key cannot be read: no such file", nil},
		{"uncertain: unlistable key", KeyFacts{Name: "k", Shape: ShapeDir, ListErr: "permission denied"},
			VerdictRetain, "key contents cannot be listed: permission denied", nil},
		{"uncertain: unreadable child entry", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeUnreadable, ShapeErr: "permission denied"}}},
			VerdictRetain, "entry wt cannot be read: permission denied", nil},
		{"uncertain: unreadable pointer", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeUnreadable, ShapeErr: "permission denied"}}}},
			VerdictRetain, "child wt .git cannot be read: permission denied", nil},
		{"uncertain: pointer read failure", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, ReadErr: "input/output error"}}}},
			VerdictRetain, "child wt .git cannot be read: input/output error", nil},
		{"uncertain: unprobeable target", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, Body: "gitdir: /gone/.git\n", TargetErr: "permission denied"}}}},
			VerdictRetain, "child wt gitdir: target cannot be read: permission denied", nil},
		{"uncertain: zero-value existence fails closed", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, Body: "gitdir: /gone/.git\n"}}}},
			VerdictRetain, "child wt gitdir: target cannot be read: ", nil},
		{"hostile: symlink key", KeyFacts{Name: "k", Shape: ShapeSymlink},
			VerdictRetain, "key is a symlink", nil},
		{"hostile: non-directory key", KeyFacts{Name: "k", Shape: ShapeRegular},
			VerdictRetain, "key is not a directory", nil},
		{"hostile: symlink child", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{Name: "wt", Shape: ShapeSymlink}}},
			VerdictRetain, "entry wt is a symlink", nil},
		{"hostile: non-directory child", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{Name: "notes.txt", Shape: ShapeRegular}}},
			VerdictRetain, "entry notes.txt is not a directory", nil},
		{"hostile: missing pointer", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeMissing, ShapeErr: "no such file"}}}},
			VerdictRetain, "child wt holds no .git entry", nil},
		{"hostile: symlink pointer", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeSymlink}}}},
			VerdictRetain, "child wt .git is a symlink", nil},
		{"hostile: repository-directory pointer", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeDir}}}},
			VerdictRetain, "child wt .git is a repository directory", nil},
		{"hostile: FIFO pointer", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeOther}}}},
			VerdictRetain, "child wt .git is not a regular file", nil},
		{"hostile: pointer without gitdir line", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, Body: "ref: refs/heads/main\n", TargetExistence: ExistenceAbsent}}}},
			VerdictRetain, "child wt .git carries no gitdir: target", nil},
		{"hostile: relative gitdir target", KeyFacts{Name: "k", Shape: ShapeDir, Children: []ChildFacts{{
			Name: "wt", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, Body: "gitdir: ../repo/.git\n", TargetExistence: ExistenceAbsent}}}},
			VerdictRetain, `child wt gitdir: target "../repo/.git" is not absolute`, nil},
		{"hostile: mixed live and dead pointers retained whole", KeyFacts{Name: "mixed", Shape: ShapeDir, Children: []ChildFacts{
			deadChild("a-dead", "/gone/.git/worktrees/a"),
			{Name: "b-live", Shape: ShapeDir, Pointer: PointerFacts{Shape: ShapeRegular, Body: "gitdir: /repo/.git/worktrees/b\n", TargetExistence: ExistencePresent}}}},
			VerdictRetain, "child b-live gitdir: target exists", nil},
		{"dead: empty key", KeyFacts{Name: "empty", Shape: ShapeDir},
			VerdictReclaim, "key holds nothing", nil},
		{"dead: every child points at an absent repository", KeyFacts{Name: "dead", Shape: ShapeDir, Children: []ChildFacts{
			deadChild("a", "/gone/.git/worktrees/a"), deadChild("b", "/gone/.git/worktrees/b")}},
			VerdictReclaim, "every child points at an absent repository", []string{"/gone/.git/worktrees/a", "/gone/.git/worktrees/b"}},
	}
	for _, c := range cases {
		t.Run(c.partition, func(t *testing.T) {
			got := ClassifyKey(c.facts)
			want := KeyVerdict{Key: c.facts.Name, Verdict: c.verdict, Reason: c.reason, Targets: c.targets}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("ClassifyKey = %#v, want %#v", got, want)
			}
		})
	}
}

// TestGitdirTargetTable pins the pointer parse: only a gitdir: line yields a
// target, the trailing space a repository path may carry survives, and a
// blank value parses as no target at all.
func TestGitdirTargetTable(t *testing.T) {
	cases := []struct {
		name, body, target string
		ok                 bool
	}{
		{"plain", "gitdir: /repo/.git/worktrees/wt\n", "/repo/.git/worktrees/wt", true},
		{"crlf and indent", " \tgitdir: /repo/.git\r\n", "/repo/.git", true},
		{"trailing space kept", "gitdir: /repo /.git \n", "/repo /.git ", true},
		{"no gitdir line", "ref: refs/heads/main\n", "", false},
		{"blank value", "gitdir: \n", "", false},
		{"empty body", "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			target, ok := GitdirTarget(c.body)
			if target != c.target || ok != c.ok {
				t.Fatalf("GitdirTarget(%q) = (%q,%v), want (%q,%v)", c.body, target, ok, c.target, c.ok)
			}
		})
	}
}

// TestFingerprintMaterialCoversExactlyTheReclaimableSubset pins drift
// sensitivity: retained keys leave the material alone, while a reclaimable
// key's name and each of its targets move it.
func TestFingerprintMaterialCoversExactlyTheReclaimableSubset(t *testing.T) {
	base := []KeyVerdict{
		{Key: "dead", Verdict: VerdictReclaim, Targets: []string{"/gone/a"}},
		{Key: "live", Verdict: VerdictRetain, Reason: "child wt gitdir: target exists"},
	}
	flatten := func(verdicts []KeyVerdict) string {
		var joined []string
		for _, part := range FingerprintMaterial(verdicts) {
			joined = append(joined, string(part))
		}
		return strings.Join(joined, "\x00")
	}
	got := flatten(base)
	if want := FingerprintVersion + "\x00key\x00dead\x001\x00target\x00/gone/a"; got != want {
		t.Fatalf("material = %q, want %q", got, want)
	}
	retainedMoved := []KeyVerdict{base[0], {Key: "other-live", Verdict: VerdictRetain}}
	if flatten(retainedMoved) != got {
		t.Fatalf("a retained key moved the material")
	}
	targetMoved := []KeyVerdict{{Key: "dead", Verdict: VerdictReclaim, Targets: []string{"/gone/b"}}, base[1]}
	if flatten(targetMoved) == got {
		t.Fatalf("a changed target left the material unchanged")
	}
}

// TestPlanDrift pins the drift verdict: only an exact match lets an apply act.
func TestPlanDrift(t *testing.T) {
	if PlanDrift("abc", "abc") {
		t.Fatalf("a matching fingerprint read as drift")
	}
	if !PlanDrift("abc", "abd") {
		t.Fatalf("a moved pool read as current")
	}
}

// TestApplyDecisionsTable pins the apply-side decisions: retained keys carry
// their protecting reason into the past tense, removal targets stay bounded
// to direct pool children, a requalification failure protects its key, and a
// removal error refuses to report success.
func TestApplyDecisionsTable(t *testing.T) {
	retained := RetainedOnApply(KeyVerdict{Key: "live", Verdict: VerdictRetain, Reason: "child wt gitdir: target exists", Targets: []string{"x"}})
	if want := (KeyVerdict{Key: "live", Verdict: VerdictRetained, Reason: "child wt gitdir: target exists"}); !reflect.DeepEqual(retained, want) {
		t.Fatalf("RetainedOnApply = %#v, want %#v", retained, want)
	}

	for _, key := range []string{"", ".", "..", "../decoy", "nested/child", "/abs/elsewhere"} {
		_, refusal, ok := RemovalBounds("/home/pool", key)
		if ok || refusal.Verdict != VerdictRetained || !strings.Contains(refusal.Reason, "not a direct child of /home/pool") {
			t.Fatalf("RemovalBounds(%q) = %#v ok=%v, want a bounded refusal", key, refusal, ok)
		}
	}
	target, _, ok := RemovalBounds("/home/pool", "dead-key")
	if !ok || target != "/home/pool/dead-key" {
		t.Fatalf("RemovalBounds(dead-key) = %q ok=%v, want the direct child", target, ok)
	}

	if _, ok := RemovalRequalified("k", KeyVerdict{Key: "k", Verdict: VerdictReclaim}); !ok {
		t.Fatalf("a still-qualifying key was refused")
	}
	refusal, ok := RemovalRequalified("k", KeyVerdict{Key: "k", Verdict: VerdictRetain, Reason: "child wt gitdir: target exists"})
	if ok || refusal.Reason != "key stopped qualifying before removal: child wt gitdir: target exists" || refusal.Verdict != VerdictRetained {
		t.Fatalf("RemovalRequalified = %#v ok=%v, want a protecting refusal", refusal, ok)
	}

	if got := RemovalOutcome("k", ""); got.Verdict != VerdictRemoved || got.Reason != "key removed" {
		t.Fatalf("RemovalOutcome success = %#v", got)
	}
	if got := RemovalOutcome("k", "permission denied"); got.Verdict != VerdictRetained || got.Reason != "removal failed: permission denied" {
		t.Fatalf("RemovalOutcome failure = %#v", got)
	}
}

// TestApplyIncompleteTable pins the exit verdict: a planned key that was not
// removed makes the run incomplete, while a plan-retained key never does.
func TestApplyIncompleteTable(t *testing.T) {
	planned := []KeyVerdict{
		{Key: "dead", Verdict: VerdictReclaim},
		{Key: "live", Verdict: VerdictRetain},
	}
	cases := []struct {
		name    string
		applied []KeyVerdict
		want    bool
	}{
		{"all planned keys removed", []KeyVerdict{{Key: "dead", Verdict: VerdictRemoved}, {Key: "live", Verdict: VerdictRetained}}, false},
		{"planned key survived", []KeyVerdict{{Key: "dead", Verdict: VerdictRetained, Reason: "removal failed: x"}, {Key: "live", Verdict: VerdictRetained}}, true},
		{"only unplanned retentions", []KeyVerdict{{Key: "live", Verdict: VerdictRetained}}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ApplyIncomplete(planned, c.applied); got != c.want {
				t.Fatalf("ApplyIncomplete = %v, want %v", got, c.want)
			}
		})
	}
}

// TestReclaimableCount pins the plan aggregate's derivation.
func TestReclaimableCount(t *testing.T) {
	verdicts := []KeyVerdict{
		{Verdict: VerdictReclaim}, {Verdict: VerdictRetain}, {Verdict: VerdictReclaim},
	}
	if got := ReclaimableCount(verdicts); got != 2 {
		t.Fatalf("ReclaimableCount = %d, want 2", got)
	}
}
