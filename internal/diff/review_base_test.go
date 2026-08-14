package diff

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// runGit runs a git command in the current working directory (set by the caller via
// t.Chdir) and fails the test loudly on a non-zero exit, carrying the combined output
// so a fixture-seeding mistake is diagnosable rather than a bare exit code.
func runGit(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// mustWriteFile writes path relative to the current working directory, creating
// parent directories as needed.
func mustWriteFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// commitFile writes content to name and commits it with msg, returning the new
// commit's sha.
func commitFile(t *testing.T, name, content, msg string) string {
	t.Helper()
	mustWriteFile(t, name, content)
	runGit(t, "add", name)
	runGit(t, "commit", "-q", "-m", msg)
	return runGit(t, "rev-parse", "HEAD")
}

// seedDivergedRepo builds a repo with a "main" branch and a "feature" branch that
// diverge after a shared ancestor:
//
//	C0 -- C1 (main) -- C3 (main only)
//	        \
//	         C2 (feature, HEAD)
//
// merge-base(main, feature) == C1. It returns the root dir and the three shas the
// scenarios key off: c0 (a reachable ancestor distinct from the merge-base), c1
// (the merge-base), and c3 (reachable but not an ancestor of feature's HEAD).
func seedDivergedRepo(t *testing.T) (root, c0, c1, c3 string) {
	t.Helper()
	root = t.TempDir()
	t.Chdir(root)
	runGit(t, "init", "-q", "-b", "main")
	runGit(t, "config", "user.email", "t@example.com")
	runGit(t, "config", "user.name", "t")

	c0 = commitFile(t, "f.txt", "0", "c0")
	c1 = commitFile(t, "f.txt", "1", "c1")
	runGit(t, "checkout", "-q", "-b", "feature")
	commitFile(t, "f.txt", "2", "c2 feature")
	runGit(t, "checkout", "-q", "main")
	c3 = commitFile(t, "f.txt", "3", "c3 main only")
	runGit(t, "checkout", "-q", "feature")
	return root, c0, c1, c3
}

// TestResolveReviewBase_RecordedPrecedence is DX1: a recorded benchBase key naming a
// reachable ancestor wins over merge-base even when the two would disagree.
func TestResolveReviewBase_RecordedPrecedence(t *testing.T) {
	root, c0, c1, _ := seedDivergedRepo(t)
	runGit(t, "config", "branch.feature.benchBase", c0)

	base, method, errKind, errHint := ResolveReviewBase(root)
	if errKind != "" || errHint != "" {
		t.Fatalf("ResolveReviewBase returned error (%q, %q), want none", errKind, errHint)
	}
	if base != c0 || method != "recorded" {
		t.Errorf("ResolveReviewBase = (%q, %q), want (%q, %q)", base, method, c0, "recorded")
	}
	if base == c1 {
		t.Fatalf("fixture invalid: recorded base equals merge-base %q; precedence is untestable", c1)
	}
}

// TestResolveReviewBase_FallbackMergeBase is DX2: with no recorded key, the default-
// branch merge-base answers, labeled `merge-base`.
func TestResolveReviewBase_FallbackMergeBase(t *testing.T) {
	root, _, c1, _ := seedDivergedRepo(t)

	base, method, errKind, errHint := ResolveReviewBase(root)
	if errKind != "" || errHint != "" {
		t.Fatalf("ResolveReviewBase returned error (%q, %q), want none", errKind, errHint)
	}
	if base != c1 || method != "merge-base" {
		t.Errorf("ResolveReviewBase = (%q, %q), want (%q, %q)", base, method, c1, "merge-base")
	}
}

// TestResolveReviewBase_UnreachableRecordedSha is DX3 (unreachable half): a recorded
// sha absent from the object database falls back to merge-base, carrying the exact
// loud label.
func TestResolveReviewBase_UnreachableRecordedSha(t *testing.T) {
	root, _, c1, _ := seedDivergedRepo(t)
	runGit(t, "config", "branch.feature.benchBase", "1234567890123456789012345678901234567890")

	base, method, errKind, errHint := ResolveReviewBase(root)
	if errKind != "" || errHint != "" {
		t.Fatalf("ResolveReviewBase returned error (%q, %q), want none", errKind, errHint)
	}
	if base != c1 || method != "merge-base (recorded sha unreachable)" {
		t.Errorf("ResolveReviewBase = (%q, %q), want (%q, %q)", base, method, c1, "merge-base (recorded sha unreachable)")
	}
}

// TestResolveReviewBase_NonAncestorRecordedSha is DX3 (non-ancestor half): a recorded
// sha that is reachable but not an ancestor of HEAD falls back to merge-base, carrying
// the exact loud label.
func TestResolveReviewBase_NonAncestorRecordedSha(t *testing.T) {
	root, _, c1, c3 := seedDivergedRepo(t)
	runGit(t, "config", "branch.feature.benchBase", c3)

	base, method, errKind, errHint := ResolveReviewBase(root)
	if errKind != "" || errHint != "" {
		t.Fatalf("ResolveReviewBase returned error (%q, %q), want none", errKind, errHint)
	}
	if base != c1 || method != "merge-base (recorded sha not an ancestor)" {
		t.Errorf("ResolveReviewBase = (%q, %q), want (%q, %q)", base, method, c1, "merge-base (recorded sha not an ancestor)")
	}
}

// TestResolveReviewBase_Unresolvable is DX4: with no resolvable default branch, the
// exported resolution returns the existing cannot-resolve error kind and hint, never
// an empty base paired with a clean return.
func TestResolveReviewBase_Unresolvable(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	runGit(t, "init", "-q", "-b", "trunk")
	runGit(t, "config", "user.email", "t@example.com")
	runGit(t, "config", "user.name", "t")
	commitFile(t, "f.txt", "0", "c0")
	runGit(t, "checkout", "-q", "-b", "other")
	commitFile(t, "f.txt", "1", "c1")

	base, method, errKind, errHint := ResolveReviewBase(root)
	if base != "" || method != "" {
		t.Errorf("ResolveReviewBase base/method = (%q, %q), want (\"\", \"\") on error", base, method)
	}
	wantKind := "cannot resolve a review base"
	wantHint := "this repository has no resolvable default branch; record one with: git config branch.<name>.benchBase <sha>"
	if errKind != wantKind || errHint != wantHint {
		t.Errorf("ResolveReviewBase error = (%q, %q), want (%q, %q)", errKind, errHint, wantKind, wantHint)
	}
}

// TestResolveBranchRangeConsumesExport is DX5: resolveBranchRange (and therefore bare
// `bench diff`) agrees byte-for-byte with the exported resolution — no second,
// divergence-capable derivation survives inline.
func TestResolveBranchRangeConsumesExport(t *testing.T) {
	_, _, c1, _ := seedDivergedRepo(t)

	out, code := Command(nil)
	if code != 0 {
		t.Fatalf("Command(nil) exit = %d, want 0; output:\n%s", code, out)
	}
	// The base cell renders quoted when the sha starts with a zero (the toon
	// layer treats a leading-zero string as numeric-looking), so the assertion
	// accepts both renderings rather than flaking on the fixture repo's sha.
	if !strings.HasPrefix(out, "revision[1]{branch,default,ahead,behind,base,method,head}:\n  feature,main,") ||
		(!strings.Contains(out, c1+",merge-base") && !strings.Contains(out, `"`+c1+`",merge-base`)) {
		t.Errorf("Command(nil) output = %q, want revision carrying merge-base %q", out, c1)
	}
}

// TestResolveBranchRangeHasNoInlineDerivation asserts the extraction actually moved
// the fallback/error logic out of resolveBranchRange rather than merely adding a
// second copy: its body must no longer reference the default-branch resolver or
// derive a merge-base directly.
func TestResolveBranchRangeHasNoInlineDerivation(t *testing.T) {
	src, err := os.ReadFile("diff.go")
	if err != nil {
		t.Fatalf("read diff.go: %v", err)
	}
	body := funcBody(t, string(src), "resolveBranchRange")
	for _, banned := range []string{"git.ResolvedDefault(", `"merge-base"`, `git.Output("merge-base"`} {
		if strings.Contains(body, banned) {
			t.Errorf("resolveBranchRange body still contains %q; inline derivation was not removed:\n%s", banned, body)
		}
	}
}

// funcBody extracts the source text of the named top-level function, from its
// `func <name>` line up to (not including) the next top-level `func ` line — good
// enough for a single-file grep assertion without pulling in go/parser.
func funcBody(t *testing.T, src, name string) string {
	t.Helper()
	start := strings.Index(src, "func "+name+"(")
	if start == -1 {
		t.Fatalf("function %s not found", name)
	}
	rest := src[start+len("func "+name+"("):]
	next := strings.Index(rest, "\nfunc ")
	if next == -1 {
		return rest
	}
	return rest[:next]
}
