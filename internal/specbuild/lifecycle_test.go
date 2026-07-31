package specbuild

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/worktree"
)

func TestAssignBindsTicketAtCandidateInOwnedWorktree(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\nAssumptions: Go is installed\n\n- [ ] [R06-R09] create a real assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "ticket details")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatalf("Start: %v", err)
	}

	assignment, status, err := service.Assign(context.Background(), "build demo", "one.md", "request 1")
	if err != nil {
		t.Fatalf("Assign: %v", err)
	}
	if assignment.ID == "" || assignment.Path == "" || assignment.Base != status.Subject {
		t.Fatalf("assignment = %#v, status = %#v", assignment, status)
	}
	if got, want := assignment.Rows, []string{"R06", "R07", "R08", "R09"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("rows = %v, want %v", got, want)
	}
	if got, want := assignment.Fence, []string{"internal/specbuild"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("fence = %v, want %v", got, want)
	}
	if got, want := assignment.Assumptions, []string{"Go is installed"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("assumptions = %v, want %v", got, want)
	}
	if got := git(t, assignment.Path, "rev-parse", "HEAD"); got != status.Subject {
		t.Fatalf("assignment base = %s, want %s", got, status.Subject)
	}
}

func TestStatusRefusesMalformedAndUnknownDurableState(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func([]byte) []byte
	}{
		{"malformed", func([]byte) []byte { return []byte("not JSON\n") }},
		{"unknown field", func(data []byte) []byte {
			return []byte(strings.TrimSuffix(string(data), "}\n") + `,"unknown":true}` + "\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			service := New(root, greenGate{}, nil)
			if _, err := service.Start(context.Background(), "build demo"); err != nil {
				t.Fatal(err)
			}
			path, err := service.statePath("build demo")
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, tc.mutate(data), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := service.Status("build demo"); err == nil {
				t.Fatal("Status accepted invalid durable state")
			}
		})
	}
}

func TestAssignAllowsSiblingsAndScopesRequestIdentity(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "two.md"), "# Two\n\nOwnership fence: internal/specbuild\n\n- [ ] [R07] sibling assignment\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "second ticket")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatal(err)
	}
	first, _, err := service.Assign(context.Background(), "build demo", "one.md", "same")
	if err != nil {
		t.Fatal(err)
	}
	if replay, _, err := service.Assign(context.Background(), "build demo", "one.md", "same"); err != nil || replay.ID != first.ID || replay.Path != first.Path || replay.Base != first.Base {
		t.Fatalf("replay = %#v, %v", replay, err)
	}
	if _, _, err := service.Assign(context.Background(), "build demo", "two.md", "same"); err == nil {
		t.Fatal("request reuse for another ticket was accepted")
	}
	second, _, err := service.Assign(context.Background(), "build demo", "two.md", "other")
	if err != nil || first.Base != second.Base || first.Path == second.Path {
		t.Fatalf("siblings = %#v, %#v, %v", first, second, err)
	}
}

func TestAssignRejectsHostileTicketsBeforeWorktreeCreation(t *testing.T) {
	root := repo(t)
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "malformed.md"), "# Malformed\n")
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "no-fence.md"), "# No fence\n\n- [ ] [R09] missing authority\n")
	if err := os.Symlink("one.md", filepath.Join(root, "specs", "build demo", "tickets", "one-link")); err != nil {
		t.Fatal(err)
	}
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "hostile ticket fixtures")
	owner := &countingOwner{}
	service := New(root, greenGate{}, owner)
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "specs", "build demo", "tickets", "one-fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "info", "exclude"), []byte("specs/build demo/tickets/one-fifo\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, ticket := range []string{"missing.md", "../spec.md", "malformed.md", "no-fence.md", "one-dir", "one-fifo", "one-link"} {
		if ticket == "one-dir" {
			if err := os.Mkdir(filepath.Join(root, "specs", "build demo", "tickets", ticket), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		if _, _, err := service.Assign(context.Background(), "build demo", ticket, ticket); err == nil {
			t.Fatalf("Assign(%q) unexpectedly succeeded", ticket)
		}
	}
	if owner.calls != 0 {
		t.Fatalf("worktree owner calls = %d", owner.calls)
	}
}

func TestStatusHasDefinitiveEmptyAndActiveProjections(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	empty, err := service.Status("build demo")
	if err != nil || empty != (Status{Slug: "build demo", State: "empty", Next: "bench spec build start build demo"}) {
		t.Fatalf("empty status = %#v, %v", empty, err)
	}
	active, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	got, err := service.Status("build demo")
	if err != nil || got != active {
		t.Fatalf("active status = %#v, %v", got, err)
	}
}

func TestStatusProjectsDurableTerminalState(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	if _, err := service.Start(context.Background(), "build demo"); err != nil {
		t.Fatal(err)
	}
	run, found, err := service.load("build demo")
	if err != nil || !found {
		t.Fatalf("load: found:%v err:%v", found, err)
	}
	run.Terminal = true
	if err := service.save(run); err != nil {
		t.Fatal(err)
	}
	status, err := service.Status("build demo")
	if err != nil || status != (Status{Slug: "build demo", State: "terminal", Subject: run.CandidateTip}) {
		t.Fatalf("terminal status = %#v, %v", status, err)
	}
}

func TestStatusUsesTheGitCommonDirectory(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)
	active, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(t.TempDir(), "linked")
	git(t, root, "worktree", "add", "--detach", "-q", linked, "HEAD")
	fromLinked, err := New(linked, greenGate{}, nil).Status("build demo")
	if err != nil || fromLinked != active {
		t.Fatalf("linked status = %#v, %v", fromLinked, err)
	}
}

func TestLiteralSlugAndTicketUseOpaqueIdentities(t *testing.T) {
	root := repo(t)
	slug := "build [special]*"
	write(t, filepath.Join(root, "specs", slug, "spec.md"), "# Special\n\nStatus: staged\n")
	write(t, filepath.Join(root, "specs", slug, "tickets", "ticket [*].md"), "# Special ticket\n\nOwnership fence: internal/specbuild\n\n- [ ] [R52] literal input\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "literal spec")
	service := New(root, greenGate{}, realOwner{})
	if _, err := service.Start(context.Background(), slug); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if _, _, err := service.Assign(context.Background(), slug, "ticket [*].md", "literal request"); err != nil {
		t.Fatalf("Assign: %v", err)
	}
	refs := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/candidate/")
	if strings.Contains(refs, " ") || strings.Contains(refs, "[") || strings.Contains(refs, "*") {
		t.Fatalf("candidate ref is not opaque: %q", refs)
	}
	if status, err := service.Status("build demo"); err != nil || status.State != "empty" {
		t.Fatalf("neighbor spec status = %#v, %v", status, err)
	}
}

type countingGate struct{ calls int }

func (g *countingGate) Bootstrap(ctx context.Context, root, branch, tip string) error {
	g.calls++
	return greenGate{}.Bootstrap(ctx, root, branch, tip)
}

type rejectGate struct{}

func (rejectGate) Bootstrap(context.Context, string, string, string) error {
	return fmt.Errorf("missing evidence")
}

type countingOwner struct{ calls int }

func (o *countingOwner) Create(context.Context, string, string, string, string) (OwnedWorktree, error) {
	o.calls++
	return OwnedWorktree{}, nil
}

func stringsSplitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

type realOwner struct{}

func (realOwner) Create(_ context.Context, root, request, label, start string) (OwnedWorktree, error) {
	created, err := worktree.Create(root, request, label, nil, start)
	if err != nil {
		return OwnedWorktree{}, err
	}
	return OwnedWorktree{ID: created.Assignment.ID, Path: created.Path}, nil
}

type greenGate struct{}

func (greenGate) Bootstrap(_ context.Context, root, branch, tip string) error {
	cmd := exec.Command("git", "-C", root, "update-ref", "refs/bench/green/"+branch, tip, "")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("establish green: %s", out)
	}
	return nil
}

func repo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "test@example.invalid"}, {"config", "user.name", "Test"}} {
		git(t, root, args...)
	}
	write(t, filepath.Join(root, "specs", "build demo", "spec.md"), "# Build demo\n\nStatus: staged\n")
	write(t, filepath.Join(root, "specs", "build demo", "tickets", "one.md"), "# One\n\nOwnership fence: internal/specbuild\n\n- [ ] [R01] exact start\n")
	write(t, filepath.Join(root, "tracked.txt"), "base\n")
	git(t, root, "add", ".")
	git(t, root, "commit", "-qm", "initial")
	return root
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(bytesTrimSpace(out))
}

func bytesTrimSpace(v []byte) []byte {
	for len(v) > 0 && (v[0] == ' ' || v[0] == '\n' || v[0] == '\r' || v[0] == '\t') {
		v = v[1:]
	}
	for len(v) > 0 && (v[len(v)-1] == ' ' || v[len(v)-1] == '\n' || v[len(v)-1] == '\r' || v[len(v)-1] == '\t') {
		v = v[:len(v)-1]
	}
	return v
}
