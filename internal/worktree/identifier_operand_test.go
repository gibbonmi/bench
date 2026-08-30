package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/intent"
)

// Every path-taking verb accepts the label, the id, or an unambiguous 8-12 character
// prefix of either. The resolver is shared; path proves each address form resolves to
// the one worktree, clean proves a verb consumes it, and release closes end to end.
func TestVerbsResolveIdentifierOperands(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "operand-forms")
	chdir(t, root)
	targets := []string{
		creation.Assignment.ID,
		creation.Assignment.Label,
		creation.Assignment.ID[:10],
		creation.Assignment.ID[:12],
		creation.Assignment.Label[:8],
	}
	for _, target := range targets {
		var stdout, stderr bytes.Buffer
		if code := PathCommand(root, home, []string{target}, &stdout, &stderr); code != 0 {
			t.Fatalf("path %q exited %d: %s", target, code, stderr.String())
		}
		if strings.TrimSpace(stdout.String()) != creation.Path {
			t.Fatalf("path %q printed %q, want %q", target, stdout.String(), creation.Path)
		}
	}
	var planned, stderr bytes.Buffer
	if code := CleanCommand(root, home, []string{creation.Assignment.ID[:10]}, &planned, &stderr); code != 0 {
		t.Fatalf("clean by id prefix exited %d: %s", code, planned.String())
	}
	if !strings.Contains(planned.String(), creation.Path) {
		t.Fatalf("clean by id prefix planned another target: %s", planned.String())
	}
	var released bytes.Buffer
	if code := ReleaseCommand(root, home, []string{"--request", "landed-operand-forms", creation.Assignment.Label}, &released, &stderr); code != 0 {
		t.Fatalf("release by label exited %d: %s", code, stderr.String())
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
		var stdout, stderr bytes.Buffer
		if code := PathCommand(root, home, []string{refusal.target}, &stdout, &stderr); code == 0 {
			t.Fatalf("%s prefix %q resolved: %s", name, refusal.target, stdout.String())
		}
		if want := "bench worktree path: " + refusal.reason + "\nnext=" + nextList + "\n"; stderr.String() != want {
			t.Fatalf("%s prefix %q printed %q, want %q", name, refusal.target, stderr.String(), want)
		}
	}
}

// `clean --apply` accepts a fingerprint prefix of at least 8 characters: one plan
// carries one digest, so the prefix is unambiguous and applies the same plan.
func TestCleanApplyAcceptsAFingerprintPrefix(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "fp-prefix")
	chdir(t, root)
	var planned, stderr bytes.Buffer
	if code := CleanCommand(root, home, []string{creation.Path}, &planned, &stderr); code != 0 {
		t.Fatalf("plan exited %d: %s", code, planned.String())
	}
	fingerprint := regexp.MustCompile(`[0-9a-f]{64}`).FindString(planned.String())
	if fingerprint == "" {
		t.Fatalf("plan carried no fingerprint: %s", planned.String())
	}
	for name, bad := range map[string]string{
		"seven-character prefix": fingerprint[:7],
		"uppercase prefix":       "ABCDEF01",
	} {
		var refused bytes.Buffer
		if code := CleanCommand(root, home, []string{creation.Path, "--apply", bad}, &refused, &stderr); code == 0 || strings.Contains(refused.String(), ",removed,") {
			t.Fatalf("%s %q was not refused: %s", name, bad, refused.String())
		}
	}
	var applied bytes.Buffer
	if code := CleanCommand(root, home, []string{creation.Path, "--apply", fingerprint[:12]}, &applied, &stderr); code != 0 {
		t.Fatalf("apply with a prefix exited %d: %s", code, applied.String())
	}
	if !strings.Contains(applied.String(), ",removed,") {
		t.Fatalf("prefix apply did not remove: %s", applied.String())
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
			root, creation, _ := newOwnedAssignment(t, "resolver-inactive")
			a := creation.Assignment
			a.State = intent.StateComplete
			mustNoError(t, intent.PutAssignment(root, a))
			return root, a.Label, "assignment " + a.ID + " is not active", nextList
		}},
		{name: "owner marker", setup: func(t *testing.T) (string, string, string, string) {
			root, creation, _ := newOwnedAssignment(t, "resolver-marker")
			rewriteMarkerOwner(t, creation.Path, strings.Repeat("a", 32))
			return root, creation.Assignment.Label, "owner marker does not match assignment " + creation.Assignment.ID, nextList
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
			root, creation, _ := newOwnedAssignment(t, "resolver-missing-landed")
			landAssignment(t, root, creation, "landed.txt")
			mustNoError(t, os.RemoveAll(creation.Path))
			return root, creation.Assignment.Label, "worktree tree is missing", "bench worktree clean --landed"
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
			var stdout, stderr bytes.Buffer
			if code := PathCommand(root, Home(), []string{target}, &stdout, &stderr); code != 1 {
				t.Fatalf("path %q exited %d, want 1: %s", target, code, stderr.String())
			}
			if want := "bench worktree path: " + reason + "\nnext=" + next + "\n"; stderr.String() != want {
				t.Errorf("path %q printed %q, want %q", target, stderr.String(), want)
			}
			stdout.Reset()
			stderr.Reset()
			if code := ExecCommand(root, Home(), []string{target, "--", "true"}, nil, &stdout, &stderr); code != 1 {
				t.Fatalf("exec %q exited %d, want 1: %s", target, code, stderr.String())
			}
			if want := "bench worktree exec: " + reason + "\nnext=" + next + "\n"; stderr.String() != want {
				t.Errorf("exec %q printed %q, want %q", target, stderr.String(), want)
			}
			stdout.Reset()
			stderr.Reset()
			if code := ShowCommand(root, Home(), []string{target, "HEAD:x"}, &stdout, &stderr); code != 1 {
				t.Fatalf("show %q exited %d, want 1: %s", target, code, stderr.String())
			}
			if want := "bench worktree show: " + reason + "\nnext=" + next + "\n"; stderr.String() != want {
				t.Errorf("show %q printed %q, want %q", target, stderr.String(), want)
			}
			stdout.Reset()
			stderr.Reset()
			// WF8: build shares the resolver and the printer, so a broken target reads the
			// same way through it as through path, exec, and show.
			if code := BuildCommand(root, Home(), []string{target}, &stdout, &stderr); code != 1 {
				t.Fatalf("build %q exited %d, want 1: %s", target, code, stderr.String())
			}
			if want := "bench worktree build: " + reason + "\nnext=" + next + "\n"; stderr.String() != want {
				t.Errorf("build %q printed %q, want %q", target, stderr.String(), want)
			}
		})
	}
}

// TestTargetVerbsShareOneRefusalPrinter is LR15, F6, and S6: one broken target yields
// byte-identical stderr from every target-taking verb once the verb prefix is stripped,
// and each ends with the same route line, so the three cannot drift.
func TestTargetVerbsShareOneRefusalPrinter(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "shared-printer")
	rewriteMarkerOwner(t, creation.Path, strings.Repeat("b", 32))
	chdir(t, root)
	target := creation.Assignment.Label
	var stdout, pathErr, execErr, showErr, buildErr bytes.Buffer
	if code := PathCommand(root, home, []string{target}, &stdout, &pathErr); code != 1 {
		t.Fatalf("path exited %d: %s", code, pathErr.String())
	}
	stdout.Reset()
	if code := ExecCommand(root, home, []string{target, "--", "true"}, nil, &stdout, &execErr); code != 1 {
		t.Fatalf("exec exited %d: %s", code, execErr.String())
	}
	stdout.Reset()
	if code := ShowCommand(root, home, []string{target, "HEAD:x"}, &stdout, &showErr); code != 1 {
		t.Fatalf("show exited %d: %s", code, showErr.String())
	}
	stdout.Reset()
	// WF8: build is the fourth verb through the one printer.
	if code := BuildCommand(root, home, []string{target}, &stdout, &buildErr); code != 1 {
		t.Fatalf("build exited %d: %s", code, buildErr.String())
	}
	pathTail, pathFound := strings.CutPrefix(pathErr.String(), "bench worktree path: ")
	execTail, execFound := strings.CutPrefix(execErr.String(), "bench worktree exec: ")
	showTail, showFound := strings.CutPrefix(showErr.String(), "bench worktree show: ")
	buildTail, buildFound := strings.CutPrefix(buildErr.String(), "bench worktree build: ")
	if !pathFound || !execFound || !showFound || !buildFound {
		t.Fatalf("verb prefixes missing: path=%q exec=%q show=%q build=%q", pathErr.String(), execErr.String(), showErr.String(), buildErr.String())
	}
	if pathTail != execTail || pathTail != showTail || pathTail != buildTail {
		t.Errorf("path tail %q, exec tail %q, show tail %q, and build tail %q differ", pathTail, execTail, showTail, buildTail)
	}
	if want := "owner marker does not match assignment " + creation.Assignment.ID + "\nnext=" + nextList + "\n"; pathTail != want {
		t.Errorf("refusal tail = %q, want %q", pathTail, want)
	}
}
