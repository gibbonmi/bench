package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/intent"
)

// Every path-taking verb accepts the label, the id, or an unambiguous 8-12 character
// prefix of either. The resolver is shared; path proves each address form resolves to
// the one worktree, clean proves a verb consumes it, and release closes end to end.
func TestVerbsResolveIdentifierOperands(t *testing.T) {
	f := newOwnedAssignment(t, "operand-forms")
	chdir(t, f.root)
	targets := []string{
		f.creation.Assignment.ID,
		f.creation.Assignment.Label,
		f.creation.Assignment.ID[:10],
		f.creation.Assignment.ID[:12],
		f.creation.Assignment.Label[:8],
	}
	for _, target := range targets {
		r := runVerb(t, verbPath, f.call(target))
		if r.exit != 0 {
			t.Fatalf("path %q exited %d: %s", target, r.exit, r.stderr)
		}
		if strings.TrimSpace(r.stdout) != f.creation.Path {
			t.Fatalf("path %q printed %q, want %q", target, r.stdout, f.creation.Path)
		}
	}
	planned := runVerb(t, verbClean, f.call(f.creation.Assignment.ID[:10]))
	if planned.exit != 0 {
		t.Fatalf("clean by id prefix exited %d: %s", planned.exit, planned.stdout)
	}
	if !strings.Contains(planned.stdout, f.creation.Path) {
		t.Fatalf("clean by id prefix planned another target: %s", planned.stdout)
	}
	if r := runVerb(t, verbRelease, f.call("--request", "landed-operand-forms", f.creation.Assignment.Label)); r.exit != 0 {
		t.Fatalf("release by label exited %d: %s", r.exit, r.stderr)
	}
}

// An ambiguous prefix and a prefix under 8 characters each stay unresolved, so a short
// word can never grab a worktree another assignment also answers to.
func TestPrefixOperandRefusals(t *testing.T) {
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	mustCreate(t, root, home, "req-prefix-a", "prefix-shared-a")
	mustCreate(t, root, home, "req-prefix-b", "prefix-shared-b")
	chdir(t, root)
	for name, refusal := range map[string]struct{ target, reason string }{
		"ambiguous": {"prefix-share", "target is ambiguous: " + strings.Join(ledgerOrderIDs(t, root), ", ")},
		"too short": {"prefix-", "target is unassigned"},
	} {
		r := runVerb(t, verbPath, repoHome{root, home}.call(refusal.target))
		if r.exit == 0 {
			t.Fatalf("%s prefix %q resolved: %s", name, refusal.target, r.stdout)
		}
		if want := "bench worktree path: " + refusal.reason + "\nnext=" + nextList + "\n"; r.stderr != want {
			t.Fatalf("%s prefix %q printed %q, want %q", name, refusal.target, r.stderr, want)
		}
	}
}

// `clean --apply` accepts a fingerprint prefix of at least 8 characters: one plan
// carries one digest, so the prefix is unambiguous and applies the same plan.
func TestCleanApplyAcceptsAFingerprintPrefix(t *testing.T) {
	f := newOwnedAssignment(t, "fp-prefix")
	chdir(t, f.root)
	planned := runVerb(t, verbClean, f.call(f.creation.Path))
	if planned.exit != 0 {
		t.Fatalf("plan exited %d: %s", planned.exit, planned.stdout)
	}
	fingerprint := planned.mustFingerprint(t)
	for name, bad := range map[string]string{
		"seven-character prefix": fingerprint[:7],
		"uppercase prefix":       "ABCDEF01",
	} {
		if refused := runVerb(t, verbClean, f.call(f.creation.Path, "--apply", bad)); refused.exit == 0 || strings.Contains(refused.stdout, ",removed,") {
			t.Fatalf("%s %q was not refused: %s", name, bad, refused.stdout)
		}
	}
	applied := runVerb(t, verbClean, f.call(f.creation.Path, "--apply", fingerprint[:12]))
	if applied.exit != 0 {
		t.Fatalf("apply with a prefix exited %d: %s", applied.exit, applied.stdout)
	}
	if !strings.Contains(applied.stdout, ",removed,") {
		t.Fatalf("prefix apply did not remove: %s", applied.stdout)
	}
}

// resolverRefusalCase breaks one dimension of a target and names the two lines both
// target-taking verbs must print for it. setup returns the repository root, the target
// the operator types, the resolver's reason, and the verb the refusal routes to.
type resolverRefusalCase struct {
	name  string
	setup func(t *testing.T) (root, target, reason, next string)
}

// nextList is the route every refusal before the target resolves names.
const nextList = "bench worktree list"

func resolverRefusalCases() []resolverRefusalCase {
	return []resolverRefusalCase{
		{name: "unassigned", setup: func(t *testing.T) (string, string, string, string) {
			root := newWorktreeRepo(t)
			return root, "no-such-target", "target is unassigned", nextList
		}},
		{name: "ambiguous", setup: func(t *testing.T) (string, string, string, string) {
			root := newWorktreeRepo(t)
			home := filepath.Join(t.TempDir(), "bench-home")
			mustCreate(t, root, home, "req-collide-a", "collide-shared-a")
			mustCreate(t, root, home, "req-collide-b", "collide-shared-b")
			return root, "collide-shar", "target is ambiguous: " + strings.Join(ledgerOrderIDs(t, root), ", "), nextList
		}},
		{name: "inactive", setup: func(t *testing.T) (string, string, string, string) {
			f := newOwnedAssignment(t, "resolver-inactive")
			a := f.creation.Assignment
			a.State = intent.StateComplete
			mustNoError(t, intent.PutAssignment(f.root, a))
			return f.root, a.Label, "assignment " + a.ID + " is not active", nextList
		}},
		{name: "owner marker", setup: func(t *testing.T) (string, string, string, string) {
			f := newOwnedAssignment(t, "resolver-marker")
			rewriteMarkerOwner(t, f.creation.Path, strings.Repeat("a", 32))
			return f.root, f.creation.Assignment.Label, "owner marker does not match assignment " + f.creation.Assignment.ID, nextList
		}},
		// F7 and F9: a removed tree is refused by name, and an assignment whose branch has
		// not landed leaves through its own release.
		{name: "missing tree unlanded", setup: func(t *testing.T) (string, string, string, string) {
			// The home sits under a directory that holds a `'`, so the recovery path is one
			// the operator cannot paste unless the producer quotes it as axi does.
			root := newWorktreeRepo(t)
			home := filepath.Join(t.TempDir(), "it's", "bench-home")
			creation := mustCreate(t, root, home, "landed-resolver-missing", "landedness")
			makeUnlandedAssignment(t, creation)
			mustNoError(t, os.RemoveAll(creation.Path))
			a := creation.Assignment
			return root, a.Label, "worktree tree is missing", "bench worktree release --request " + a.RequestToken + " " + axi.ShellQuote(a.Worktree)
		}},
		// F7 and F8: a landed assignment leaves with the batch clean instead.
		{name: "missing tree landed", setup: func(t *testing.T) (string, string, string, string) {
			f := newOwnedAssignment(t, "resolver-missing-landed")
			landAssignment(t, f.root, f.creation, "landed.txt")
			mustNoError(t, os.RemoveAll(f.creation.Path))
			return f.root, f.creation.Assignment.Label, "worktree tree is missing", "bench worktree clean --landed"
		}},
	}
}

// ledgerOrderIDs lists the assignment ids in ledger order, which is the order the ambiguity
// refusal names them in.
func ledgerOrderIDs(t *testing.T, root string) []string {
	t.Helper()
	assignments, err := intent.Assignments(root)
	mustNoError(t, err)
	ids := make([]string, 0, len(assignments))
	for _, a := range assignments {
		ids = append(ids, a.ID)
	}
	return ids
}

// TestTargetVerbsNameTheResolverReason is LR11 through LR14 and F5, F7, F8, and F9: an
// operator reads the check that failed rather than one blanket sentence, and then reads
// the one verb that answers it.
func TestTargetVerbsNameTheResolverReason(t *testing.T) {
	for _, refusalCase := range resolverRefusalCases() {
		t.Run(refusalCase.name, func(t *testing.T) {
			root, target, reason, next := refusalCase.setup(t)
			chdir(t, root)
			at := repoHome{root, Home()}
			path := runVerb(t, verbPath, at.call(target))
			if path.exit != 1 {
				t.Fatalf("path %q exited %d, want 1: %s", target, path.exit, path.stderr)
			}
			if want := "bench worktree path: " + reason + "\nnext=" + next + "\n"; path.stderr != want {
				t.Errorf("path %q printed %q, want %q", target, path.stderr, want)
			}
			exec := runVerb(t, verbExec, at.call(target, "--", "true"))
			if exec.exit != 1 {
				t.Fatalf("exec %q exited %d, want 1: %s", target, exec.exit, exec.stderr)
			}
			if want := "bench worktree exec: " + reason + "\nnext=" + next + "\n"; exec.stderr != want {
				t.Errorf("exec %q printed %q, want %q", target, exec.stderr, want)
			}
			show := runVerb(t, verbShow, at.call(target, "HEAD:x"))
			if show.exit != 1 {
				t.Fatalf("show %q exited %d, want 1: %s", target, show.exit, show.stderr)
			}
			if want := "bench worktree show: " + reason + "\nnext=" + next + "\n"; show.stderr != want {
				t.Errorf("show %q printed %q, want %q", target, show.stderr, want)
			}
			// WF8: build shares the resolver and the printer, so a broken target reads the
			// same way through it as through path, exec, and show.
			build := runVerb(t, verbBuild, at.call(target))
			if build.exit != 1 {
				t.Fatalf("build %q exited %d, want 1: %s", target, build.exit, build.stderr)
			}
			if want := "bench worktree build: " + reason + "\nnext=" + next + "\n"; build.stderr != want {
				t.Errorf("build %q printed %q, want %q", target, build.stderr, want)
			}
		})
	}
}

// TestTargetVerbsShareOneRefusalPrinter is LR15, F6, and S6: one broken target yields
// byte-identical stderr from every target-taking verb once the verb prefix is stripped,
// and each ends with the same route line, so the three cannot drift.
func TestTargetVerbsShareOneRefusalPrinter(t *testing.T) {
	f := newOwnedAssignment(t, "shared-printer")
	rewriteMarkerOwner(t, f.creation.Path, strings.Repeat("b", 32))
	chdir(t, f.root)
	target := f.creation.Assignment.Label
	path := runVerb(t, verbPath, f.call(target))
	if path.exit != 1 {
		t.Fatalf("path exited %d: %s", path.exit, path.stderr)
	}
	exec := runVerb(t, verbExec, f.call(target, "--", "true"))
	if exec.exit != 1 {
		t.Fatalf("exec exited %d: %s", exec.exit, exec.stderr)
	}
	show := runVerb(t, verbShow, f.call(target, "HEAD:x"))
	if show.exit != 1 {
		t.Fatalf("show exited %d: %s", show.exit, show.stderr)
	}
	// WF8: build is the fourth verb through the one printer.
	build := runVerb(t, verbBuild, f.call(target))
	if build.exit != 1 {
		t.Fatalf("build exited %d: %s", build.exit, build.stderr)
	}
	pathTail, pathFound := strings.CutPrefix(path.stderr, "bench worktree path: ")
	execTail, execFound := strings.CutPrefix(exec.stderr, "bench worktree exec: ")
	showTail, showFound := strings.CutPrefix(show.stderr, "bench worktree show: ")
	buildTail, buildFound := strings.CutPrefix(build.stderr, "bench worktree build: ")
	if !pathFound || !execFound || !showFound || !buildFound {
		t.Fatalf("verb prefixes missing: path=%q exec=%q show=%q build=%q", path.stderr, exec.stderr, show.stderr, build.stderr)
	}
	if pathTail != execTail || pathTail != showTail || pathTail != buildTail {
		t.Errorf("path tail %q, exec tail %q, show tail %q, and build tail %q differ", pathTail, execTail, showTail, buildTail)
	}
	if want := "owner marker does not match assignment " + f.creation.Assignment.ID + "\nnext=" + nextList + "\n"; pathTail != want {
		t.Errorf("refusal tail = %q, want %q", pathTail, want)
	}
}
