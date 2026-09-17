//go:build system

package systemtest

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// evidenceJourney is one disposable repository with a staged build spec and the Bench
// home its assignments live in.
type evidenceJourney struct {
	root, home, base string
}

// newEvidenceJourney seeds a build-ready spec: one ticket, its completion plan, and the
// canonical guidance sources a build charge requires.
func newEvidenceJourney(t *testing.T) evidenceJourney {
	t.Helper()
	root, err := os.MkdirTemp(owner.root, "charge-evidence [journey]-")
	if err != nil {
		t.Fatal(err)
	}
	home, err := os.MkdirTemp(owner.root, "charge-evidence [home]-")
	if err != nil {
		t.Fatal(err)
	}
	systemGit(t, root, "init", "-q", "-b", "main")
	systemGit(t, root, "config", "user.email", "bench@local")
	systemGit(t, root, "config", "user.name", "bench")
	for path, body := range map[string]string{
		".agents/skills/bench-craft-delegate/SKILL.md":                            "# Delegation skill\n",
		".agents/skills/bench-craft-delegate/references/delegation-discipline.md": "# Delegation procedure\n",
		".agents/commands/bench-implement-spec.md":                                "# Build phase\n",
	} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	spec := "# x\n\nStatus: staged\n\n## User stories\n1. Prepare evidence.\n\n### Acceptance coverage map\n| row | story | behavior | seam | why it catches the failure |\n|---|---|---|---|---|\n| E1 | 1 | prepares | command | catches failure |\n\n## Ownership fences\n\n- `source.txt`\n- `reviews/x.md`\n"
	recordtest.Prepare(t, root, 1, "specs/x/spec.md", spec)
	return evidenceJourney{root: root, home: home, base: systemGitOutput(t, root, "rev-parse", "HEAD")}
}

func (j evidenceJourney) env(extra ...string) []string {
	return append([]string{"BENCH_HOME=" + j.home, "BENCH_SYSTEM_ROOT=" + j.root, "BENCH_COMMAND_OBSERVE=1"}, extra...)
}

// assignment creates one Bench assignment and commits body to source.txt in it, unless
// from names the label of a sibling whose tip the assignment starts at.
func (j evidenceJourney) assignment(t *testing.T, label, body string, from ...string) systemLandingWorktree {
	t.Helper()
	extra := []string{}
	if len(from) > 0 {
		extra = []string{"--from", from[0]}
	}
	worktree := systemCreateLandingWorktree(t, j.root, j.home, label, label+" request", extra...)
	if body != "" {
		systemCommit(t, worktree.path, "source.txt", body, "change "+label)
	}
	worktree.tip = systemGitOutput(t, worktree.path, "rev-parse", "HEAD")
	return worktree
}

func (j evidenceJourney) prepareArgs(worktree systemLandingWorktree, extra ...string) []string {
	return append([]string{"preflight", "build", "x", "--charge", "--ticket", "1.md", "--base", j.base, "--source-tip", worktree.tip}, extra...)
}

var preparedIdentity = regexp.MustCompile(`(?m)^  "?(sha256:[0-9a-f]{64})"?,build,`)

func (j evidenceJourney) prepare(t *testing.T, worktree systemLandingWorktree, extra ...string) processResult {
	t.Helper()
	return systemSelected(t, worktree.path, j.env(), j.prepareArgs(worktree, extra...)...)
}

func identityOf(t *testing.T, result processResult) string {
	t.Helper()
	match := preparedIdentity.FindStringSubmatch(result.stdout)
	if result.code != 0 || match == nil || !strings.Contains(result.stderr, "command-registry:preflight") {
		t.Fatalf("preparation = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	return match[1]
}

func (j evidenceJourney) packs(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(j.root, ".git", "bench-charge-evidence"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".lock") {
			names = append(names, entry.Name())
		}
	}
	return names
}

// requiredBytes reads the exact pack size a capacity refusal names for an empty store.
func (j evidenceJourney) requiredBytes(t *testing.T, worktree systemLandingWorktree) uint64 {
	t.Helper()
	result := j.prepare(t, worktree, "--max-store-bytes", "1")
	match := regexp.MustCompile(`retry with --max-store-bytes ([0-9]+)`).FindStringSubmatch(result.stdout)
	if result.code != 1 || match == nil {
		t.Fatalf("sizing refusal = (%d, %q)", result.code, result.stdout)
	}
	size, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return size
}

// startPaused starts one preparation that pauses at stage and waits for its marker.
func (j evidenceJourney) startPaused(t *testing.T, worktree systemLandingWorktree, stage string, extra ...string) (*exec.Cmd, string, func() string) {
	t.Helper()
	marker := filepath.Join(j.home, worktree.request+" "+stage+".marker")
	cmd, stdout, _ := systemStartSelected(t, worktree.path, j.env("BENCH_EVIDENCE_PAUSE="+stage+":"+marker), j.prepareArgs(worktree, extra...)...)
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	waitForFile(t, marker, "the preparation to reach stage "+stage)
	return cmd, marker, stdout.String
}

// waitForFile waits a bounded window until path exists.
func waitForFile(t *testing.T, path, wait string) {
	t.Helper()
	window := bounds.TestDeadline(0)
	for deadline := time.Now().Add(window); ; time.Sleep(10 * time.Millisecond) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal(bounds.TestTimeoutVerdict(wait, window))
		}
	}
}

func waitExit(t *testing.T, cmd *exec.Cmd) int {
	t.Helper()
	return systemExitCode(cmd.Wait())
}

// TestEvidenceLocationIndependence is CE17: sibling assignments at one source tip prepare
// one identity, and the assignment stays outside it.
func TestEvidenceLocationIndependence(t *testing.T) {
	j := newEvidenceJourney(t)
	first := j.assignment(t, "location-first", "changed\n")
	second := j.assignment(t, "location-second", "", "location-first")
	if first.tip != second.tip || first.path == second.path {
		t.Fatalf("siblings = %+v, %+v", first, second)
	}
	one, two := j.prepare(t, first), j.prepare(t, second)
	if identityOf(t, one) != identityOf(t, two) || !strings.Contains(one.stdout, first.assignment) || !strings.Contains(two.stdout, second.assignment) {
		t.Fatalf("sibling preparations = %q and %q", one.stdout, two.stdout)
	}
	if packs := j.packs(t); len(packs) != 1 {
		t.Fatalf("sibling preparations published %v", packs)
	}
}

// TestEvidenceQuotaConcurrency is CE77. The first writer pauses between its capacity
// calculation and its temporary growth. The second writer pauses just before it requests the
// writer lock and is released first; only after it reports its resumption does the first
// writer resume. A writer lock that did not serialize the two lets the second writer count
// an empty store and publish past the quota.
func TestEvidenceQuotaConcurrency(t *testing.T) {
	j := newEvidenceJourney(t)
	first := j.assignment(t, "quota-first", "first\n")
	second := j.assignment(t, "quota-second", "second body\n")
	sizeFirst, sizeSecond := j.requiredBytes(t, first), j.requiredBytes(t, second)
	quota := strconv.FormatUint(max(sizeFirst, sizeSecond)+min(sizeFirst, sizeSecond)-1, 10)
	paused, marker, _ := j.startPaused(t, first, "capacity", "--max-store-bytes", quota)
	waiting, secondMarker, secondOut := j.startPaused(t, second, "writer-lock", "--max-store-bytes", quota)
	if err := os.Remove(secondMarker); err != nil {
		t.Fatal(err)
	}
	waitForFile(t, secondMarker+".resumed", "the second writer to leave its pause")
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	if code := waitExit(t, paused); code != 0 {
		t.Fatalf("paused writer exit = %d", code)
	}
	if code := waitExit(t, waiting); code != 1 || !strings.Contains(secondOut(), "evidence capacity") {
		t.Fatalf("second writer = (%d, %q)", code, secondOut())
	}
	if packs := j.packs(t); len(packs) != 1 {
		t.Fatalf("concurrent writers left %v under the quota", packs)
	}
}

// TestEvidenceInterruptedPublication is CE86 and CE87.
func TestEvidenceInterruptedPublication(t *testing.T) {
	for _, stage := range []string{"staged", "verifying"} {
		t.Run(stage, func(t *testing.T) {
			j := newEvidenceJourney(t)
			worktree := j.assignment(t, "interrupt-"+stage, "interrupted\n")
			cmd, _, _ := j.startPaused(t, worktree, stage)
			if err := cmd.Process.Kill(); err != nil {
				t.Fatal(err)
			}
			_ = cmd.Wait()
			entries := j.packs(t)
			if len(entries) != 1 || !strings.HasPrefix(entries[0], "tmp-") {
				t.Fatalf("interrupted %s left %v, want one temporary pack and no artifact", stage, entries)
			}
			identity := identityOf(t, j.prepare(t, worktree))
			if read := systemSelected(t, worktree.path, j.env(), "preflight", "evidence", identity); read.code != 0 || !strings.HasPrefix(read.stdout, "page[1]") {
				t.Fatalf("read after interruption = (%d, %q)", read.code, read.stdout)
			}
		})
	}
}

// TestEvidenceConcurrentPublication is CE89 and CE90.
func TestEvidenceConcurrentPublication(t *testing.T) {
	t.Run("CE89 identical writers", func(t *testing.T) {
		j := newEvidenceJourney(t)
		first := j.assignment(t, "identical-first", "identical\n")
		second := j.assignment(t, "identical-second", "", "identical-first")
		paused, marker, stdout := j.startPaused(t, first, "staged")
		blocked, blockedOut, _ := systemStartSelected(t, second.path, j.env(), j.prepareArgs(second)...)
		if err := os.Remove(marker); err != nil {
			t.Fatal(err)
		}
		if code := waitExit(t, paused); code != 0 {
			t.Fatalf("first writer exit = %d", code)
		}
		published := j.packs(t)
		before, err := os.Stat(filepath.Join(j.root, ".git", "bench-charge-evidence", published[0]))
		if err != nil {
			t.Fatal(err)
		}
		if code := waitExit(t, blocked); code != 0 {
			t.Fatalf("second writer = (%d, %q)", code, blockedOut.String())
		}
		after, err := os.Stat(filepath.Join(j.root, ".git", "bench-charge-evidence", published[0]))
		firstID := preparedIdentity.FindStringSubmatch(stdout())
		secondID := preparedIdentity.FindStringSubmatch(blockedOut.String())
		if err != nil || !os.SameFile(before, after) || !after.ModTime().Equal(before.ModTime()) || firstID == nil || secondID == nil || firstID[1] != secondID[1] || len(j.packs(t)) != 1 {
			t.Fatalf("identical writers replaced or duplicated the artifact: %v", j.packs(t))
		}
	})
	t.Run("CE90 distinct writers", func(t *testing.T) {
		j := newEvidenceJourney(t)
		first := j.assignment(t, "distinct-first", "first\n")
		second := j.assignment(t, "distinct-second", "second\n")
		one, _, _ := systemStartSelected(t, first.path, j.env(), j.prepareArgs(first)...)
		two, _, _ := systemStartSelected(t, second.path, j.env(), j.prepareArgs(second)...)
		if waitExit(t, one) != 0 || waitExit(t, two) != 0 {
			t.Fatal("a distinct writer failed")
		}
		if packs := j.packs(t); len(packs) != 2 {
			t.Fatalf("distinct writers left %v", packs)
		}
	})
}
