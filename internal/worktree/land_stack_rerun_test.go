// Stack re-run tests: each stack refusal ends with the stack's own re-run, and an override
// pins the tip that its member lands.
package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

// stackMemberDetail is the opening of the refusal of stack member member, whose label the
// refusal names.
func stackMemberDetail(member string, m stackFixtureMember) string {
	return "stack member " + member + " (" + m.creation.Assignment.Label + "): "
}

// stackRefusalNext reads the next= value of the refusal of stack member member, whose
// label the refusal names.
func stackRefusalNext(t *testing.T, stdout, member string, m stackFixtureMember) string {
	t.Helper()
	next, printed := landingFaceNext(stdout, stackMemberDetail(member, m))
	if !printed {
		t.Fatalf("stack stdout = %q, want the refusal of member %s", stdout, member)
	}
	return next
}

// advanceMember commits one more change to the member's own file and answers the new tip.
func advanceMember(t *testing.T, m stackFixtureMember) string {
	t.Helper()
	commitInWorktree(t, m.creation.Path, m.file, "moved bytes\n", "move the member")
	return gitOutput(t, m.creation.Path, "rev-parse", "HEAD")
}

// TestLandStackRoutesEndWithTheStackRerun is LS86. In a three-member stack, the route of a
// source-not-clean member 2 ends with the stack re-run that adds member 2's repaired slot.
func TestLandStackRoutesEndWithTheStackRerun(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-route", 3, nil)
	m := f.members[1]
	mustWrite(t, filepath.Join(m.creation.Path, m.file), []byte("dirty\n"), 0o644)
	file, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	rerun := "bench worktree land --stack '" + file + "' --source-tip '" + m.request + "'=<repaired-tip>"
	if want, next := landingRoute(faceSourceNotClean, rerun, labelOf(m.creation)), stackRefusalNext(t, r.stdout, "2", m); r.exit != 1 || next != want {
		t.Fatalf("stack = (%d, next %q), want exit 1 and next %q", r.exit, next, want)
	}
}

// TestLandStackRefusesAStaleMemberTip is LS87. A member whose branch moved past its file
// source_tip refuses with source-tip-mismatch, which names the tip that the branch holds,
// and the re-run pins that tip.
func TestLandStackRefusesAStaleMemberTip(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-stale", 2, nil)
	m := f.members[1]
	moved := advanceMember(t, m)
	file, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 1 || !strings.Contains(r.stdout, "wanted="+moved+",") {
		t.Fatalf("stale member = (%d, %q, %q), want the mismatch naming %s", r.exit, r.stdout, r.stderr, moved)
	}
	rerun := "bench worktree land --stack '" + file + "' --source-tip '" + m.request + "'='" + moved + "'"
	if want, next := landingRoute(faceSourceTipMismatch, rerun, nil), stackRefusalNext(t, r.stdout, "2", m); next != want {
		t.Fatalf("stale member next = %q, want %q", next, want)
	}
}

// TestLandStackRerunKeepsEachOverride is LS94. A call that pins member 1 with an override
// and refuses member 2 re-runs with the member 1 override and member 2's repaired slot.
func TestLandStackRerunKeepsEachOverride(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-keeps", 2, nil)
	first, second := f.members[0], f.members[1]
	moved := advanceMember(t, first)
	mustWrite(t, filepath.Join(second.creation.Path, second.file), []byte("dirty\n"), 0o644)
	file, args := f.stack(t, first.request+"="+moved)
	r := runVerb(t, verbLand, f.call(args...))
	rerun := "bench worktree land --stack '" + file + "' --source-tip '" + first.request + "'='" + moved + "' --source-tip '" + second.request + "'=<repaired-tip>"
	if want, next := landingRoute(faceSourceNotClean, rerun, labelOf(second.creation)), stackRefusalNext(t, r.stdout, "2", second); r.exit != 1 || next != want {
		t.Fatalf("stack = (%d, next %q), want exit 1 and next %q", r.exit, next, want)
	}
	if strings.Count(r.stdout, "refused{") != 1 {
		t.Fatalf("stack stdout = %q, want member 2's refusal alone", r.stdout)
	}
}

// TestStackRerunRules composes the re-run of one three-member stack whose call pinned
// members 2 and 3. A refused member that holds an override gets its new tip in that
// override's place. A removed or a reset member loses its override, and a route that leaves
// one member re-runs that member's single landing.
func TestStackRerunRules(t *testing.T) {
	t.Parallel()
	route := stackRoute{
		file: "/stack.json",
		input: StackInput{
			Members: []StackMember{
				{Request: "r1", Base: "b1", SourceTip: "t1", Spec: "s1"},
				{Request: "r2", Base: "b2", SourceTip: "t2", Spec: "s2"},
				{Request: "r3", Base: "b3", SourceTip: "t3", Spec: "s3"},
			},
			Overrides: []StackOverride{{Request: "r2", Commit: "c2"}, {Request: "r3", Commit: "c3"}},
		},
		paths:       []string{"/w1", "/w2", "/w3"},
		assignments: []string{"a1", "a2", "a3"},
	}
	stack := "bench worktree land --stack '/stack.json'"
	for _, row := range []struct {
		name string
		fact stackRerun
		want string
	}{
		{"repair-replaces", stackRerun{refused: 1, repairs: true}, stack + " --source-tip 'r2'=<repaired-tip> --source-tip 'r3'='c3'"},
		{"read-replaces", stackRerun{refused: 1, read: "t9"}, stack + " --source-tip 'r2'='t9' --source-tip 'r3'='c3'"},
		{"removed", stackRerun{refused: -1, removed: []int{2}}, stack + " --source-tip 'r2'='c2'"},
		{"reset", stackRerun{refused: -1, reset: []int{1}}, stack + " --source-tip 'r3'='c3'"},
		{"one-repaired", stackRerun{refused: 0, repairs: true, removed: []int{1, 2}}, landingRerunAt("r1", "b1", repairedSourceTipFlag, "s1", "/w1", "a1")},
		{"one-pinned", stackRerun{refused: -1, removed: []int{0, 1}}, landingRerunAt("r3", "b3", landingSourceTipFlag("c3"), "s3", "/w3", "a3")},
	} {
		if got := route.rerun(row.fact); got != row.want {
			t.Errorf("%s: rerun = %q, want %q", row.name, got, row.want)
		}
	}
}

// TestLandStackOverrideReplacesTheFileTip is LS95. An override replaces its member's file
// source_tip, and that member's merge takes the override tip as its second parent.
func TestLandStackOverrideReplacesTheFileTip(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-override", 2, nil)
	m := f.members[1]
	moved := advanceMember(t, m)
	_, args := f.stack(t, m.request+"="+moved)
	r := runVerb(t, verbLand, f.call(args...))
	if r.exit != 0 {
		t.Fatalf("stack with an override = (%d, %q, %q), want exit 0", r.exit, r.stdout, r.stderr)
	}
	if parent := gitOutput(t, f.root, "rev-parse", "main^2"); parent != moved {
		t.Fatalf("second parent of member 2's merge = %s, want the override tip %s", parent, moved)
	}
}
