package worktree

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every path-taking verb accepts the label, the id, or an unambiguous 8-12 character
// prefix of either. The resolver is shared; path proves each address form resolves to
// the one worktree, clean proves a verb consumes it, and release closes end to end.
func TestVerbsResolveIdentifierOperands(t *testing.T) {
	root, creation := newOwnedAssignment(t, "operand-forms")
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
		if code := PathCommand(root, []string{target}, &stdout, &stderr); code != 0 {
			t.Fatalf("path %q exited %d: %s", target, code, stderr.String())
		}
		if strings.TrimSpace(stdout.String()) != creation.Path {
			t.Fatalf("path %q printed %q, want %q", target, stdout.String(), creation.Path)
		}
	}
	var planned, stderr bytes.Buffer
	if code := CleanCommand([]string{creation.Assignment.ID[:10]}, &planned, &stderr); code != 0 {
		t.Fatalf("clean by id prefix exited %d: %s", code, planned.String())
	}
	if !strings.Contains(planned.String(), creation.Path) {
		t.Fatalf("clean by id prefix planned another target: %s", planned.String())
	}
	var released bytes.Buffer
	if code := ReleaseCommand(root, []string{"--request", "landed-operand-forms", creation.Assignment.Label}, &released, &stderr); code != 0 {
		t.Fatalf("release by label exited %d: %s", code, stderr.String())
	}
}

// An ambiguous prefix and a prefix under 8 characters each stay unresolved, so a short
// word can never grab a worktree another assignment also answers to.
func TestPrefixOperandRefusals(t *testing.T) {
	root := newWorktreeRepo(t)
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "bench-home"))
	mustCreate(t, root, "req-prefix-a", "prefix-shared-a")
	mustCreate(t, root, "req-prefix-b", "prefix-shared-b")
	chdir(t, root)
	for name, target := range map[string]string{
		"ambiguous": "prefix-share",
		"too short": "prefix-",
	} {
		var stdout, stderr bytes.Buffer
		if code := PathCommand(root, []string{target}, &stdout, &stderr); code == 0 {
			t.Fatalf("%s prefix %q resolved: %s", name, target, stdout.String())
		}
		if !strings.Contains(stderr.String(), "not one active Bench-owned worktree") {
			t.Fatalf("%s prefix %q lost the named refusal: %s", name, target, stderr.String())
		}
	}
}

// `clean --apply` accepts a fingerprint prefix of at least 8 characters: one plan
// carries one digest, so the prefix is unambiguous and applies the same plan.
func TestCleanApplyAcceptsAFingerprintPrefix(t *testing.T) {
	root, creation := newOwnedAssignment(t, "fp-prefix")
	chdir(t, root)
	var planned, stderr bytes.Buffer
	if code := CleanCommand([]string{creation.Path}, &planned, &stderr); code != 0 {
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
		if code := CleanCommand([]string{creation.Path, "--apply", bad}, &refused, &stderr); code == 0 || strings.Contains(refused.String(), ",removed,") {
			t.Fatalf("%s %q was not refused: %s", name, bad, refused.String())
		}
	}
	var applied bytes.Buffer
	if code := CleanCommand([]string{creation.Path, "--apply", fingerprint[:12]}, &applied, &stderr); code != 0 {
		t.Fatalf("apply with a prefix exited %d: %s", code, applied.String())
	}
	if !strings.Contains(applied.String(), ",removed,") {
		t.Fatalf("prefix apply did not remove: %s", applied.String())
	}
}
