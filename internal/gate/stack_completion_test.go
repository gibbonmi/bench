package gate

// The stack proof. Each test commits a destination and member commits on it in an outcome
// fixture, crafts a graded tree, and grades that tree under a stack context on a landing of
// the destination.

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

// stackBase is the destination of every stack fixture. Member 1 changes docs/one.md, member 2
// changes docs/two.md, and no member changes docs/kept.md.
var stackBase = map[string]string{"docs/one.md": "one\n", "docs/two.md": "two\n", "docs/kept.md": "kept\n"}

// stackFixture is an outcome fixture whose destination commit holds stackBase.
type stackFixture struct {
	t                 *testing.T
	root, destination string
}

func newStackFixture(t *testing.T) *stackFixture {
	t.Helper()
	root := outcomeFixture(t)
	for path, body := range stackBase {
		outcomeWrite(t, root, path, body, 0o644)
	}
	outcomeCommit(t, root, "destination")
	return &stackFixture{t: t, root: root, destination: outcomeGit(t, root, "rev-parse", "HEAD")}
}

// tree answers the destination tree with each file of files written and each path of
// deleted removed.
func (f *stackFixture) tree(files map[string]string, deleted ...string) string {
	f.t.Helper()
	outcomeGit(f.t, f.root, "read-tree", "--reset", "-u", f.destination)
	for path, body := range files {
		outcomeWrite(f.t, f.root, path, body, 0o644)
	}
	for _, path := range deleted {
		outcomeGit(f.t, f.root, "rm", "-q", "--", path)
	}
	outcomeGit(f.t, f.root, "add", "-A")
	return outcomeGit(f.t, f.root, "write-tree")
}

// member commits the tree of files and deleted on the destination, and answers the commit.
func (f *stackFixture) member(files map[string]string, deleted ...string) string {
	f.t.Helper()
	return outcomeGit(f.t, f.root, "commit-tree", "-p", f.destination, "-m", "member", f.tree(files, deleted...))
}

// pair answers the spec-less members 1 and 2 of stackBase.
func (f *stackFixture) pair() []StackMember {
	f.t.Helper()
	return []StackMember{
		{Source: f.member(map[string]string{"docs/one.md": "ONE\n"})},
		{Source: f.member(map[string]string{"docs/two.md": "TWO\n"})},
	}
}

// stack answers ctx with a landing on the destination and the stack context of members.
func (f *stackFixture) stack(members []StackMember) context.Context {
	return WithStack(WithLanding(context.Background(), f.destination, ""), members)
}

// grade answers the checkpoint fault of tree under the stack context of members.
func (f *stackFixture) grade(members []StackMember, tree string) error {
	f.t.Helper()
	graded, err := captureProspectiveTree(f.root, tree)
	if err != nil {
		f.t.Fatal(err)
	}
	_, err = checkpointEvaluation(f.stack(members), newGateEvaluation(f.root)).applyCheckpoint(graded, subject{})
	return err
}

// requireFault requires that err is a completion proof fault with the message want.
func requireFault(t *testing.T, err error, want string) {
	t.Helper()
	if !errors.As(err, new(completionProofError)) || err.Error() != want {
		t.Fatalf("stack proof = %v; want the proof fault %q", err, want)
	}
}

// exactPair is the exact tree of the members of pair.
var exactPair = map[string]string{"docs/one.md": "ONE\n", "docs/two.md": "TWO\n"}

// The exact composition of two members passes the stack proof.
func TestStackProofAcceptsTheExactStack(t *testing.T) {
	f := newStackFixture(t)
	if err := f.grade(f.pair(), f.tree(exactPair)); err != nil {
		t.Fatalf("exact stack tree = %v; want a pass", err)
	}
}

// A tree whose bytes at a member's changed path differ from that member's refuses, so the
// proof reads the entries and not the path set alone.
func TestStackProofRefusesAChangedMemberPath(t *testing.T) {
	f := newStackFixture(t)
	err := f.grade(f.pair(), f.tree(map[string]string{"docs/one.md": "rewritten\n", "docs/two.md": "TWO\n"}))
	if !errors.As(err, new(completionProofError)) || !strings.Contains(err.Error(), "docs/one.md") {
		t.Fatalf("rewritten member path = %v; want a proof fault that names docs/one.md", err)
	}
}

// A tree that changes a path that no member changes refuses, also for a stack of spec-less
// members.
func TestStackProofRefusesAnUnreviewedDelta(t *testing.T) {
	f := newStackFixture(t)
	err := f.grade(f.pair(), f.tree(map[string]string{"docs/one.md": "ONE\n", "docs/two.md": "TWO\n", "docs/kept.md": "changed\n"}))
	requireFault(t, err, "stack composition changes docs/kept.md; review the destination delta")
}

// A tree that keeps an entry under a member's tickets folder refuses, and the tree that
// closes that folder passes.
func TestStackProofRefusesAKeptFolder(t *testing.T) {
	f := newStackFixture(t)
	const folder, ticket = "specs/example/tickets", "specs/example/tickets/01-ticket.md"
	members := []StackMember{
		{Source: f.member(map[string]string{"docs/one.md": "ONE\n", ticket: "ticket\n"}), Folder: folder},
		{Source: f.member(map[string]string{"docs/two.md": "TWO\n"})},
	}
	if err := f.grade(members, f.tree(exactPair)); err != nil {
		t.Fatalf("closed folder = %v; want a pass", err)
	}
	err := f.grade(members, f.tree(map[string]string{"docs/one.md": "ONE\n", "docs/two.md": "TWO\n", ticket: "ticket\n"}))
	requireFault(t, err, "completion keeps closed "+ticket+"; review the delivery closure")
}

// A capture/ path that two members change passes with bytes of neither member, because the
// settle policy writes the union there.
func TestStackProofLeavesASharedCapturePathToTheSettlePolicy(t *testing.T) {
	f := newStackFixture(t)
	members := []StackMember{
		{Source: f.member(map[string]string{"capture/learnings.md": "one\n"})},
		{Source: f.member(map[string]string{"capture/learnings.md": "two\n"})},
	}
	if err := f.grade(members, f.tree(map[string]string{"capture/learnings.md": "one\ntwo\n"})); err != nil {
		t.Fatalf("settled capture path = %v; want a pass", err)
	}
}

// A path outside capture/ that two members change refuses and names both members, so the
// bytes of the last owner never overwrite the first in silence.
func TestStackProofRefusesAnOverlap(t *testing.T) {
	f := newStackFixture(t)
	members := []StackMember{
		{Source: f.member(map[string]string{"docs/one.md": "first\n"})},
		{Source: f.member(map[string]string{"docs/one.md": "second\n"})},
	}
	err := f.grade(members, f.tree(map[string]string{"docs/one.md": "second\n"}))
	requireFault(t, err, "stack members 1 and 2 change docs/one.md; land them apart")
}

// A destination path that one member deletes passes only when the tree lacks it.
func TestStackProofHonorsAMemberDeletion(t *testing.T) {
	f := newStackFixture(t)
	members := []StackMember{{Source: f.member(nil, "docs/kept.md")}, f.pair()[1]}
	if err := f.grade(members, f.tree(map[string]string{"docs/two.md": "TWO\n"}, "docs/kept.md")); err != nil {
		t.Fatalf("deleted path absent = %v; want a pass", err)
	}
	if err := f.grade(members, f.tree(map[string]string{"docs/two.md": "TWO\n"})); err == nil || !strings.Contains(err.Error(), "docs/kept.md") {
		t.Fatalf("deleted path kept = %v; want a refusal that names docs/kept.md", err)
	}
}

// A one-member stack context accepts its exact tree, so each bisect prefix can grade.
func TestStackProofAcceptsAOneMemberPrefix(t *testing.T) {
	f := newStackFixture(t)
	members := f.pair()[:1]
	if err := f.grade(members, f.tree(map[string]string{"docs/one.md": "ONE\n"})); err != nil {
		t.Fatalf("one-member prefix = %v; want a pass", err)
	}
}

// A green of the top tree under no stack context does not answer a stack grade whose proof
// fails, because the proof runs before the reuse lookup.
func TestStackProofRunsOnAReusedGreen(t *testing.T) {
	f := newStackFixture(t)
	members, top := f.pair(), f.tree(exactPair)
	if result := executeTreeWithOwner(context.Background(), f.root, top, io.Discard, io.Discard, nil, reuseFreshGreen, nil); result.ActionExit != 0 {
		t.Fatalf("green of the top tree = %+v", result)
	}
	var stderr bytes.Buffer
	result := executeTreeWithOwner(f.stack(members[:1]), f.root, top, io.Discard, &stderr, nil, reuseFreshGreen, nil)
	if result.ActionExit == 0 || !strings.Contains(stderr.String(), "stack composition changes docs/two.md") {
		t.Fatalf("stack grade after a green = %+v %q; want the proof refusal", result, stderr.String())
	}
	requireRuns(t, f.root, 1, "after the refused stack grade")
}

// A landing with a stack context names the landing stack case and runs the whole gate.
func TestReaderPlanRunsWholeForALandingStack(t *testing.T) {
	f := newLandingFixture(t, nil)
	if _, eligible, out := f.decide(WithStack(f.landing(), []StackMember{{Source: f.destination}}), docsChange, nil, f.root); eligible || out != "gate: whole gate (reader set refused: landing stack)\n" {
		t.Fatalf("plan of a landing stack = %t %q; want the whole gate with the landing stack line", eligible, out)
	}
}

// A stack landing over an uncovered path names the landing stack case, which comes before
// the change cases.
func TestReaderPlanNamesTheStackCaseFirst(t *testing.T) {
	f := newLandingFixture(t, nil)
	_, eligible, out := f.decide(WithStack(f.landing(), []StackMember{{Source: f.destination}}), func(f *landingFixture) { f.put("AGENTS.md", "change\n") }, nil, f.root)
	if eligible || out != "gate: whole gate (reader set refused: landing stack)\n" {
		t.Fatalf("plan of a stack over an uncovered path = %t %q; want the landing stack line", eligible, out)
	}
}
