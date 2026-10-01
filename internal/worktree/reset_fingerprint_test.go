package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestResetFingerprintTracksTheCheckout(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-fingerprint")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before != after, "tracked edit did not change fingerprint: %s", after)
}

func TestResetFingerprintTracksTheContent(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-content")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("first\n"), 0o644)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("second\n"), 0o644)
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before != after, "content edit did not change fingerprint: %s", after)
}

func TestResetFingerprintIgnoresIgnoredFiles(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-ignored")
	commitInWorktree(t, creation.Path, ".gitignore", "output\n", "ignore output")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(creation.Path, "output"), []byte("build\n"), 0o644)
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before == after, "ignored output changed fingerprint: %s != %s", before, after)
}

func TestResetFingerprintTracksTheTip(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-tip")
	commitInWorktree(t, creation.Path, "ahead", "one\n", "ahead")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	commitInWorktree(t, creation.Path, "ahead", "two\n", "ahead again")
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before != after, "new commit did not change fingerprint")
}

func TestResetFingerprintTracksTheRef(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-ref")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	gitRun(t, creation.Path, "switch", "--detach", "HEAD")
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before != after, "detach did not change fingerprint")
}

func TestResetFingerprintTracksTheCheckpoint(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-checkpoint")
	commitInWorktree(t, creation.Path, "ahead", "ahead\n", "ahead")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("dirty\n"), 0o644)
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", gitOutput(t, creation.Path, "rev-parse", "HEAD"), creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before != after, "checkpoint did not change fingerprint")
}

func TestResetFingerprintTracksTheIndex(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "reset-index")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("staged one\n"), 0o644)
	gitRun(t, creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("working\n"), 0o644)
	status := gitOutput(t, creation.Path, "status", "--porcelain=v1")
	plan := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	before := plan.mustFingerprint(t)
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("staged two\n"), 0o644)
	gitRun(t, creation.Path, "add", "README.md")
	mustWrite(t, filepath.Join(creation.Path, "README.md"), []byte("working\n"), 0o644)
	requireTest(t, gitOutput(t, creation.Path, "status", "--porcelain=v1") == status, "fixture changed the status")
	plan = runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID}})
	requireTest(t, plan.exit == 0, "fingerprint plan = %d %s %s", plan.exit, plan.stdout, plan.stderr)
	after := plan.mustFingerprint(t)
	requireTest(t, before != after, "staged-only edit did not change fingerprint: %s", after)
	result := runVerb(t, verbReset, verbCall{root: root, home: home, args: []string{"--to", creation.Assignment.Start, creation.Assignment.ID, "--apply", before}})
	requireTest(t, result.exit == 1 && strings.Contains(result.stdout, "reset plan is stale") && gitOutput(t, creation.Path, "show", ":README.md") == "staged two",
		"stale staged apply = %d %s %s", result.exit, result.stdout, result.stderr)
}
