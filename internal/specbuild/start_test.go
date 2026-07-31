package specbuild

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/spec"
)

func TestStartCreatesRunFromExactGreenEvidence(t *testing.T) {
	root := repo(t)
	service := New(root, greenGate{}, nil)

	status, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if status.State != "active" || status.Subject != git(t, root, "rev-parse", "HEAD") || status.Next != "bench spec build assign build demo" {
		t.Fatalf("Start status = %#v", status)
	}
	branch := git(t, root, "symbolic-ref", "--short", "HEAD")
	if git(t, root, "rev-parse", "refs/bench/green/"+branch) != git(t, root, "rev-parse", "HEAD") {
		t.Fatal("green marker does not name the initial tip")
	}
	refs := git(t, root, "for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/specbuild/candidate/")
	if refs == "" || refs[len(refs)-40:] != git(t, root, "rev-parse", "HEAD") {
		t.Fatalf("candidate ref = %q", refs)
	}
}

func TestStartResumeAndConflictsDoNotDuplicateRun(t *testing.T) {
	tests := []struct {
		name   string
		change func(*testing.T, string)
	}{
		{"branch", func(t *testing.T, root string) { git(t, root, "checkout", "-qb", "other") }},
		{"tip", func(t *testing.T, root string) {
			write(t, filepath.Join(root, "later"), "x")
			git(t, root, "add", ".")
			git(t, root, "commit", "-qm", "later")
		}},
		{"dirt", func(t *testing.T, root string) { write(t, filepath.Join(root, "dirty"), "x") }},
		{"detached", func(t *testing.T, root string) { git(t, root, "checkout", "--detach", "-q") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			service := New(root, greenGate{}, nil)
			if _, err := service.Start(context.Background(), "build demo"); err != nil {
				t.Fatal(err)
			}
			tc.change(t, root)
			if _, err := service.Start(context.Background(), "build demo"); err == nil {
				t.Fatal("Start unexpectedly accepted a changed working subject")
			}
			if got := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/specbuild/candidate/"); len(stringsSplitLines(got)) != 1 {
				t.Fatalf("candidate identities = %q", got)
			}
		})
	}
	root := repo(t)
	gate := &countingGate{}
	service := New(root, gate, nil)
	first, err := service.Start(context.Background(), "build demo")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Start(context.Background(), "build demo")
	if err != nil || second != first || gate.calls != 1 {
		t.Fatalf("resume = %#v, %v; gate calls = %d", second, err, gate.calls)
	}
}

func TestStartWithoutEvidenceNamesOneRecoveryAndDoesNotMutate(t *testing.T) {
	root := repo(t)
	service := New(root, rejectGate{}, nil)
	if _, err := service.Start(context.Background(), "build demo"); err == nil || !strings.Contains(err.Error(), "bench gate, then retry start") {
		t.Fatalf("Start error = %v", err)
	}
	if got := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/"); got != "" {
		t.Fatalf("failed start created refs: %q", got)
	}
	status, err := service.Status("build demo")
	if err != nil || status.State != "empty" {
		t.Fatalf("status after failed start = %#v, %v", status, err)
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

func TestStartRefusesConflictingCandidateAndInvalidPriorState(t *testing.T) {
	t.Run("candidate compare-and-swap", func(t *testing.T) {
		root := repo(t)
		gate := &countingGate{}
		service := New(root, gate, nil)
		_, resolved, _, ok, err := spec.Resolve(root, "build demo")
		if err != nil || !ok {
			t.Fatalf("Resolve: %v, %v", ok, err)
		}
		candidate := "refs/bench/specbuild/candidate/" + digest(resolved)
		before := git(t, root, "rev-parse", "HEAD")
		write(t, filepath.Join(root, "later"), "x")
		git(t, root, "add", ".")
		git(t, root, "commit", "-qm", "later")
		git(t, root, "update-ref", candidate, before)
		if _, err := service.Start(context.Background(), "build demo"); err == nil {
			t.Fatal("Start overwrote a pre-existing candidate")
		}
		if got := git(t, root, "rev-parse", candidate); got != before {
			t.Fatalf("candidate = %s, want %s", got, before)
		}
		if gate.calls != 0 {
			t.Fatalf("gate calls = %d", gate.calls)
		}
		if got := git(t, root, "for-each-ref", "--format=%(refname)", "refs/bench/green/"); got != "" {
			t.Fatalf("candidate conflict created a green marker: %q", got)
		}
		if _, found, err := service.load("build demo"); err != nil || found {
			t.Fatalf("run state after candidate conflict = found:%v err:%v", found, err)
		}
	})
	for _, tc := range []struct {
		name   string
		mutate func(*record)
	}{
		{"conflicting run identity", func(run *record) { run.Run = "another-run" }},
		{"incomplete assignments", func(run *record) { run.Assignments = nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := repo(t)
			service := New(root, greenGate{}, nil)
			if _, err := service.Start(context.Background(), "build demo"); err != nil {
				t.Fatal(err)
			}
			run, found, err := service.load("build demo")
			if err != nil || !found {
				t.Fatalf("load initial run: found:%v err:%v", found, err)
			}
			candidate := run.Candidate
			before := git(t, root, "rev-parse", candidate)
			tc.mutate(&run)
			if err := service.save(run); err != nil {
				t.Fatal(err)
			}
			if _, err := service.Start(context.Background(), "build demo"); err == nil {
				t.Fatal("Start accepted invalid durable state")
			}
			if got := git(t, root, "rev-parse", candidate); got != before {
				t.Fatalf("candidate changed from %s to %s", before, got)
			}
		})
	}
}
