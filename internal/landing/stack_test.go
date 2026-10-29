package landing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/gittest"
)

// stackLanding is one destination and the reviewed member sources of a stack. Every member
// starts from one review base, so the members compose in any order unless a row makes two
// of them change one path.
type stackLanding struct {
	root, base string
	members    []stackSource
}

type stackSource struct{ name, worktree, branch, tip, folder string }

// newStackLanding commits, on the branch of member i, the member's own file, its tickets
// folder, and the files of extra[i]. Member indexes start at 1.
func newStackLanding(t *testing.T, n int, extra map[int]map[string]string) stackLanding {
	t.Helper()
	root := fixture(t)
	f := stackLanding{root: root, base: gittest.Output(t, root, "rev-parse", "HEAD")}
	for i := 1; i <= n; i++ {
		name := fmt.Sprintf("member-%d", i)
		worktree := filepath.Join(t.TempDir(), name)
		gittest.MustRun(t, root, "worktree", "add", "-qb", name, worktree, f.base)
		folder := "specs/" + name
		files := map[string]string{name: name + " bytes\n", folder + "/tickets/01.md": "# " + name + "\n"}
		maps.Copy(files, extra[i])
		for path, content := range files {
			write(t, worktree, path, content)
		}
		gittest.MustRun(t, worktree, "add", "-A")
		gittest.MustRun(t, worktree, "commit", "-qm", "review "+name)
		f.members = append(f.members, stackSource{
			name: name, worktree: worktree, branch: "refs/heads/" + name,
			tip: gittest.Output(t, worktree, "rev-parse", "HEAD"), folder: folder,
		})
	}
	return f
}

// request is the stack request that lands every member in fixture order, each member
// closing its own tickets folder.
func (f stackLanding) request(t *testing.T) StackRequest {
	t.Helper()
	fingerprint := func(root string) string {
		got, err := CheckoutFingerprint(root)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	members := make([]ReviewedRequest, len(f.members))
	for i, m := range f.members {
		members[i] = ReviewedRequest{
			Source: m.branch, SourceTip: m.tip, ReviewBase: f.base,
			SourceWorktree: m.worktree, SourceFingerprint: fingerprint(m.worktree),
			ClosePath: m.folder, Message: "land " + m.name,
		}
	}
	return StackRequest{
		Root: f.root, Destination: "refs/heads/main", DestinationBase: f.base,
		DestinationFingerprint: fingerprint(f.root), Members: members,
		Stdout: io.Discard, Stderr: io.Discard,
	}
}

// stackRecord is what a stack owner observed: each graded tree, the stack members of the
// last grade, and each ref update.
type stackRecord struct {
	graded, updates []string
	members         []gate.StackMember
	// duringGrade runs inside the grade, so a row can move a source while the gate runs.
	duringGrade func()
}

// stackOwner grades every tree with kind and records each grade and each ref update. The
// ref update still runs, so a row reads the published ref.
func stackOwner(kind authorization.Kind) (Owner, *stackRecord) {
	record := &stackRecord{}
	o := greenOwner()
	o.authorize = func(ctx context.Context, _, tree string, _, _ io.Writer) authorization.Result {
		record.graded = append(record.graded, tree)
		record.members = gate.StackMembers(ctx)
		if record.duringGrade != nil {
			record.duringGrade()
		}
		return authorization.Result{Kind: kind}
	}
	o.updateRef = func(root, ref, next, old string) error {
		record.updates = append(record.updates, ref+" "+old+" "+next)
		return updateRef(root, ref, next, old)
	}
	return o, record
}

// errAdmissionRefused is the refusal of a refusing stackAdmissionRecord.
var errAdmissionRefused = errors.New("admission refused")

// stackAdmissionRecord records what each decision received. It admits every member unless
// refuseMember names the 1-based member that CheckMember refuses or refusePublish is set.
type stackAdmissionRecord struct {
	checked       []string
	published     int
	trees         []string
	refuseMember  int
	refusePublish bool
}

func (a *stackAdmissionRecord) CheckMember(i int, tree string) error {
	if i != len(a.checked) {
		return fmt.Errorf("member %d checked out of order", i)
	}
	a.checked = append(a.checked, tree)
	if i+1 == a.refuseMember {
		return errAdmissionRefused
	}
	return nil
}

func (a *stackAdmissionRecord) PublishStack(k int, trees []string, publish func() error) error {
	if a.refusePublish {
		return errAdmissionRefused
	}
	a.published, a.trees = k, slices.Clone(trees)
	return publish()
}

// landGreenStack lands an n-member stack under a green grade and fails the row on a refusal.
func landGreenStack(t *testing.T, n int) (stackLanding, StackResult, *stackRecord, *stackAdmissionRecord) {
	t.Helper()
	f := newStackLanding(t, n, nil)
	o, record := stackOwner(authorization.Green)
	admission := &stackAdmissionRecord{}
	got, err := o.LandStack(context.Background(), f.request(t), admission)
	if err != nil {
		t.Fatal(err)
	}
	return f, got, record, admission
}

func treeNames(t *testing.T, root, tree string) []string {
	t.Helper()
	return strings.Split(gittest.Output(t, root, "ls-tree", "-r", "--name-only", tree), "\n")
}

func assertStackUnpublished(t *testing.T, f stackLanding, record *stackRecord) {
	t.Helper()
	if len(record.updates) != 0 || gittest.Output(t, f.root, "rev-parse", "refs/heads/main") != f.base {
		t.Fatalf("a refused stack moved the destination: updates %q", record.updates)
	}
}

func TestLandStackGradesOnce(t *testing.T) {
	f, got, record, _ := landGreenStack(t, 3)
	top := gittest.Output(t, f.root, "rev-parse", got.Merges[len(got.Merges)-1]+"^{tree}")
	if !slices.Equal(record.graded, []string{top}) {
		t.Fatalf("graded trees = %q, want the top tree %s once", record.graded, top)
	}
	t.Run("the grade carries each member in file order", func(t *testing.T) {
		f := newStackLanding(t, 3, nil)
		o, record := stackOwner(authorization.Green)
		request := f.request(t)
		request.Members[2].ClosePath = ""
		if _, err := o.LandStack(context.Background(), request, &stackAdmissionRecord{}); err != nil {
			t.Fatal(err)
		}
		want := []gate.StackMember{{Source: f.members[0].tip, Folder: f.members[0].folder}, {Source: f.members[1].tip, Folder: f.members[1].folder}, {Source: f.members[2].tip}}
		if !slices.Equal(record.members, want) {
			t.Fatalf("graded stack members = %+v, want %+v", record.members, want)
		}
	})
}

func TestLandStackChainsMemberMerges(t *testing.T) {
	f, got, _, _ := landGreenStack(t, 3)
	if got.Published() != len(f.members) || got.Kind != authorization.Green || got.DestinationBase != f.base {
		t.Fatalf("stack result = %+v", got)
	}
	previous := f.base
	for i, m := range f.members {
		merge := got.Merges[i]
		if parents := gittest.Output(t, f.root, "rev-list", "--parents", "-n", "1", merge); parents != merge+" "+previous+" "+m.tip {
			t.Fatalf("member %d merge parents = %q, want %s then %s", i+1, parents, previous, m.tip)
		}
		if message := gittest.Output(t, f.root, "log", "-1", "--format=%B", merge); message != "land "+m.name {
			t.Fatalf("member %d merge message = %q", i+1, message)
		}
		previous = merge
	}
}

func TestLandStackPublishesWithOneUpdate(t *testing.T) {
	f, got, record, _ := landGreenStack(t, 3)
	top := got.Merges[len(got.Merges)-1]
	if want := []string{"refs/heads/main " + f.base + " " + top}; !slices.Equal(record.updates, want) {
		t.Fatalf("ref updates = %q, want %q", record.updates, want)
	}
	if gittest.Output(t, f.root, "rev-parse", "refs/heads/main") != top {
		t.Fatal("the destination does not name the top merge")
	}
}

func TestLandStackClosesEachFolderInItsMerge(t *testing.T) {
	f, got, _, _ := landGreenStack(t, 2)
	for _, m := range f.members {
		if gittest.Output(t, f.root, "ls-tree", "-r", "--name-only", m.tip, "--", m.folder) == "" {
			t.Fatalf("fixture tip of %s carries no tickets folder; the close is unobservable", m.name)
		}
	}
	for i, merge := range got.Merges {
		for _, m := range f.members[:i+1] {
			if listed := gittest.Output(t, f.root, "ls-tree", "-r", "--name-only", merge, "--", m.folder); listed != "" {
				t.Fatalf("merge %d still carries the folder of %s: %q", i+1, m.name, listed)
			}
		}
		if gittest.Output(t, f.root, "show", merge+":"+f.members[i].name) != f.members[i].name+" bytes" {
			t.Fatalf("merge %d lost the reviewed bytes of its member", i+1)
		}
	}
}

func TestLandStackAdmitsEachMemberTree(t *testing.T) {
	f, got, _, admission := landGreenStack(t, 3)
	if len(admission.checked) != len(f.members) {
		t.Fatalf("admission checked %d members, want %d", len(admission.checked), len(f.members))
	}
	for i, tree := range admission.checked {
		// The own landing tree is the destination with this member alone, its folder closed.
		if want := []string{"foreign", f.members[i].name, "named"}; !slices.Equal(treeNames(t, f.root, tree), want) {
			t.Fatalf("member %d admission tree holds %q, want %q", i+1, treeNames(t, f.root, tree), want)
		}
	}
	if got.Published() != admission.published || !slices.Equal(admission.trees, admission.checked) {
		t.Fatalf("publication decided %d members over %q, want %d over %q", admission.published, admission.trees, got.Published(), admission.checked)
	}
	t.Run("a refused member stops the stack before the grade", func(t *testing.T) {
		f := newStackLanding(t, 3, nil)
		o, record := stackOwner(authorization.Green)
		admission := &stackAdmissionRecord{refuseMember: 2}
		got, err := o.LandStack(context.Background(), f.request(t), admission)
		var member StackMemberError
		if !errors.As(err, &member) || member.Member != 2 || !errors.Is(err, errAdmissionRefused) {
			t.Fatalf("refused member = %v, want the admission refusal of member 2", err)
		}
		if len(record.graded) != 0 || got.Published() != 0 || admission.published != 0 {
			t.Fatalf("a refused member graded %d trees and published %d", len(record.graded), got.Published())
		}
		assertStackUnpublished(t, f, record)
	})
}

func TestLandStackStopsOnInfrastructure(t *testing.T) {
	f := newStackLanding(t, 3, nil)
	o, record := stackOwner(authorization.Infrastructure)
	admission := &stackAdmissionRecord{}
	got, err := o.LandStack(context.Background(), f.request(t), admission)
	var refusal AuthorizationRefusal
	if !errors.As(err, &refusal) || refusal.Result.Kind != authorization.Infrastructure || got.Kind != authorization.Infrastructure {
		t.Fatalf("infrastructure stack = %+v, %v", got, err)
	}
	if len(record.graded) != 1 || got.Published() != 0 || admission.published != 0 {
		t.Fatalf("infrastructure stack graded %d trees and published %d", len(record.graded), got.Published())
	}
	assertStackUnpublished(t, f, record)
}

func TestLandStackRefusesAMemberConflict(t *testing.T) {
	f := newStackLanding(t, 3, map[int]map[string]string{1: {"named": "member one"}, 2: {"named": "member two"}})
	o, record := stackOwner(authorization.Green)
	_, err := o.LandStack(context.Background(), f.request(t), &stackAdmissionRecord{})
	var member StackMemberError
	var conflict ConflictError
	if !errors.As(err, &member) || member.Member != 2 || !errors.As(err, &conflict) {
		t.Fatalf("member conflict = %v, want a conflict error of member 2", err)
	}
	if len(record.graded) != 0 {
		t.Fatalf("a conflicted stack graded %q", record.graded)
	}
	assertStackUnpublished(t, f, record)
}

func TestLandStackRechecksEachMember(t *testing.T) {
	f := newStackLanding(t, 3, nil)
	o, record := stackOwner(authorization.Green)
	record.duringGrade = func() {
		gittest.MustRun(t, f.members[1].worktree, "commit", "--allow-empty", "-qm", "moved during the grade")
	}
	_, err := o.LandStack(context.Background(), f.request(t), &stackAdmissionRecord{})
	var member StackMemberError
	if !errors.As(err, &member) || member.Member != 2 {
		t.Fatalf("moved member = %v, want a refusal of member 2", err)
	}
	assertStackUnpublished(t, f, record)
	t.Run("a changed destination checkout refuses the publication", func(t *testing.T) {
		f := newStackLanding(t, 3, nil)
		o, record := stackOwner(authorization.Green)
		record.duringGrade = func() { write(t, f.root, "changed-during-the-grade", "bytes\n") }
		got, err := o.LandStack(context.Background(), f.request(t), &stackAdmissionRecord{})
		var member StackMemberError
		if err == nil || errors.As(err, &member) || got.Published() != 0 {
			t.Fatalf("changed destination = %+v, %v, want a destination refusal", got, err)
		}
		assertStackUnpublished(t, f, record)
	})
	t.Run("a refused publication publishes nothing", func(t *testing.T) {
		f := newStackLanding(t, 3, nil)
		o, record := stackOwner(authorization.Green)
		got, err := o.LandStack(context.Background(), f.request(t), &stackAdmissionRecord{refusePublish: true})
		if !errors.Is(err, errAdmissionRefused) || got.Published() != 0 || len(got.Merges) != 0 {
			t.Fatalf("refused publication = %+v, %v, want the admission refusal and no merges", got, err)
		}
		assertStackUnpublished(t, f, record)
	})
}

func TestLandStackRefusesAStagedSpecRequest(t *testing.T) {
	f := newStackLanding(t, 2, nil)
	o, record := stackOwner(authorization.Green)
	request := f.request(t)
	request.Members[1].ClosePath, request.Members[1].SpecPath = "", "specs/member-2/spec.md"
	_, err := o.LandStack(context.Background(), request, &stackAdmissionRecord{})
	var member StackMemberError
	if !errors.As(err, &member) || member.Member != 2 || len(record.graded) != 0 {
		t.Fatalf("staged-spec member = %v after %d grades, want a refusal of member 2 before any grade", err, len(record.graded))
	}
	assertStackUnpublished(t, f, record)
}
