// Resume refusal tests: destructive destination state, a non-ancestor review base, a stale marker, and an evicted receipt.
package worktree

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

func TestResumeLandCommandPublicRefusesDestructiveDestinationState(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	for _, journey := range []struct {
		name  string
		later bool
	}{
		{name: "PL25"},
		{name: "PL30", later: true},
	} {
		for _, tc := range []struct {
			name    string
			allowed bool
			detail  string
			setup   func(*testing.T, string)
		}{
			{name: "clean"},
			{name: "declared ignored output", allowed: true, setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("dist/\n"), 0o644)
				mustMkdirAll(t, filepath.Join(root, "dist"), 0o755)
				mustWrite(t, filepath.Join(root, "dist", "out"), []byte("build output\n"), 0o600)
			}},
			{name: "staged changes", detail: "landing destination has staged changes", setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, "staged.txt"), []byte("staged\n"), 0o600)
				gitRun(t, root, "add", "staged.txt")
			}},
			{name: "tracked-worktree changes", detail: "landing destination has tracked-worktree changes", setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, "owned.txt"), []byte("caller bytes\n"), 0o600)
			}},
			{name: "untracked collision", detail: "landing destination has untracked collisions", setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, "untracked-collision.txt"), []byte("caller bytes\n"), 0o600)
			}},
			{name: "ignored residue", detail: "landing destination has ignored residue", setup: func(t *testing.T, root string) {
				mustWrite(t, filepath.Join(root, ".git", "info", "exclude"), []byte("ignored-residue\n"), 0o644)
				mustWrite(t, filepath.Join(root, "ignored-residue"), []byte("caller bytes\n"), 0o600)
			}},
			{name: "nested repository", detail: "landing destination has nested repositories", setup: func(t *testing.T, root string) {
				nested := filepath.Join(root, "nested")
				mustMkdirAll(t, nested, 0o755)
				gitRun(t, nested, "init", "-q", "-b", "main")
				mustWrite(t, filepath.Join(nested, "nested.txt"), []byte("base\n"), 0o644)
				gitRun(t, nested, "add", "nested.txt")
				gitRun(t, nested, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "nested")
			}},
		} {
			t.Run(journey.name+"/"+tc.name, func(t *testing.T) {
				request := "resume-destination-state-" + strings.ReplaceAll(journey.name+"-"+tc.name, " ", "-")
				root, creation, base, tip, tally, _ := publicLandingFixture(t, request, "private/output", "dist/")
				land := func(args ...string) (int, string, string) {
					var stdout, stderr bytes.Buffer
					cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
					cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
					return exitCode(cmd.Run()), stdout.String(), stderr.String()
				}
				if code, stdout, stderr := land("--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path); code != 3 || !strings.Contains(stdout, "worktree=incomplete:release") || !strings.Contains(stderr, "worktree retained (ignored)") {
					t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout, stderr)
				}
				published := gitOutput(t, root, "rev-parse", "main")
				mustRemove(t, filepath.Join(creation.Path, "private", "output"))
				if journey.later {
					commitInWorktree(t, root, "destination-after-publication", "forward\n", "destination movement")
				}
				if tc.setup != nil {
					tc.setup(t, root)
				}
				before := resumeDestinationState(t, root)
				code, stdout, stderr := land("--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path)
				if tc.setup == nil || tc.allowed {
					if code != 0 || !strings.Contains(stdout, "worktree=released}") || stderr != "" {
						t.Fatalf("resume = (%d, %q, %q)", code, stdout, stderr)
					}
					return
				}
				if code != 1 || stdout != "refused{detail="+tc.detail+"}\n" || stderr != "" {
					t.Fatalf("destructive-state refusal = (%d, %q, %q), want detail %q", code, stdout, stderr, tc.detail)
				}
				if after := resumeDestinationState(t, root); after != before {
					t.Fatalf("refusal changed destination state:\nbefore:\n%safter:\n%s", before, after)
				}
				if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
					t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
				}
			})
		}
	}
}

func TestResumeLandCommandRefusesNonAncestorReviewBaseWithoutMutation(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	request := "resume-nonancestor-base"
	root, creation, base, tip, tally, _ := publicLandingFixture(t, request, "", "")
	run := func(args ...string) (int, string, string) {
		var stdout, stderr bytes.Buffer
		cmd := descendant(t, binary, append([]string{"worktree", "land"}, args...)...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
		return exitCode(cmd.Run()), stdout.String(), stderr.String()
	}
	if code, stdout, stderr := run("--request", request, "--base", base, "--source-tip", tip, "--spec", "x", "-m", "land reviewed source", creation.Path); code != 0 || !strings.Contains(stdout, "worktree=released}") || stderr == "" {
		t.Fatalf("landing = (%d, %q, %q)", code, stdout, stderr)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	commitInWorktree(t, root, "destination-after-publication", "forward\n", "destination movement")
	destination := gitOutput(t, root, "rev-parse", "main")
	code, stdout, stderr := run("--resume", published, "--request", request, "--base", destination, "--source-tip", tip, "--spec", "x", creation.Path)
	if code != 1 || !strings.Contains(stdout, "review base does not authenticate the published source") || stderr != "" {
		t.Fatalf("nonancestor resume = (%d, %q, %q)", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != destination {
		t.Fatalf("nonancestor resume moved destination: got %s want %s", got, destination)
	}
	if got := gitOutput(t, root, "rev-parse", "refs/bench/green/main"); got != published {
		t.Fatalf("nonancestor resume moved project-green: got %s want %s", got, published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("nonancestor resume reran gate: tally=%q error=%v", got, err)
	}
}

func TestResumeLandCommandRefusesAbsentOrBehindMarkerAfterDestinationMoves(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		marker func(*testing.T, string, string, string)
	}{
		{name: "absent", marker: func(t *testing.T, root, _, _ string) { gitRun(t, root, "update-ref", "-d", "refs/bench/green/main") }},
		{name: "behind", marker: func(t *testing.T, root, base, _ string) { gitRun(t, root, "update-ref", "refs/bench/green/main", base) }},
		{name: "divergent", marker: func(t *testing.T, root, _, tip string) { gitRun(t, root, "update-ref", "refs/bench/green/main", tip) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "resume-marker-" + tc.name
			root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
			working := defaultJoins()
			broken := working
			broken.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 1 }
			var stdout, stderr bytes.Buffer
			if code := landWith(broken, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
				t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			commitInWorktree(t, root, "destination-after-publication", "forward\n", "destination movement")
			destination := gitOutput(t, root, "rev-parse", "main")
			tc.marker(t, root, base, tip)
			stdout.Reset()
			stderr.Reset()
			args := []string{"--resume", gitOutput(t, root, "rev-parse", "main~1"), "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
			if code := landWith(working, root, home, "", args, &stdout, &stderr); code != 1 || !strings.HasPrefix(stdout.String(), "refused{detail=project-green marker") || stderr.Len() != 0 {
				t.Fatalf("marker refusal = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			if got := gitOutput(t, root, "rev-parse", "main"); got != destination {
				t.Fatalf("marker refusal moved destination: got %s want %s", got, destination)
			}
			if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
				t.Fatalf("resume reran gate: tally=%q error=%v", got, err)
			}
		})
	}
}

func TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted(t *testing.T) {
	t.Parallel()
	request := "resume-evicted-receipt"
	root, creation, base, tip, tally, home := publicLandingFixture(t, request, "", "")
	var stdout, stderr bytes.Buffer
	if code := LandCommand(root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "worktree=released}") {
		t.Fatalf("landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	published := gitOutput(t, root, "rev-parse", "main")
	for i := 0; i < intent.MaxCleanupReceipts; i++ {
		receipt := intent.CleanupReceipt{
			Schema: intent.CleanupReceiptSchema, Repo: root, Operation: "eviction-test", Target: filepath.Join(root, fmt.Sprintf("target-%d", i)),
			Fingerprint: intent.RequestDigest(fmt.Sprintf("evict-%d", i)), State: intent.ReceiptComplete, Phase: intent.ReceiptPhaseTerminal, Action: "removed",
		}
		if err := intent.PutCleanupReceipt(root, receipt); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "update-ref", "-d", "refs/bench/green/main")
	gitRun(t, root, "read-tree", "main^")
	gitRun(t, root, "checkout-index", "-a", "-f")
	staged := gitOutput(t, root, "diff", "--cached", "--name-only")
	stdout.Reset()
	stderr.Reset()
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}
	if code := LandCommand(root, home, "", args, &stdout, &stderr); code != 1 || !strings.Contains(stdout.String(), "missing-terminal-receipt") || stderr.Len() != 0 {
		t.Fatalf("evicted resume = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != published {
		t.Fatalf("evicted resume moved destination: got %s want %s", got, published)
	}
	if descendant(t, "git", "-C", root, "show-ref", "--verify", "--quiet", "refs/bench/green/main").Run() == nil {
		t.Fatal("evicted resume recreated project-green marker")
	}
	if got := gitOutput(t, root, "diff", "--cached", "--name-only"); got != staged {
		t.Fatalf("evicted resume reconciled checkout: got %q want %q", got, staged)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
		t.Fatalf("evicted resume reran gate: tally=%q error=%v", got, err)
	}
}

func resumeDestinationState(t *testing.T, root string) string {
	t.Helper()
	index, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	if err != nil {
		t.Fatal(err)
	}
	assignment, err := os.ReadFile(filepath.Join(root, ".git", intent.Filename))
	if err != nil {
		t.Fatal(err)
	}
	marker := gitOutput(t, root, "rev-parse", "refs/bench/green/main")
	return fmt.Sprintf("refs=%smarker=%sindex=%x\nassignment=%x\nworktree=%s", gitOutput(t, root, "show-ref", "--head"), marker, index, assignment, resumeWorktreeBytes(t, root))
}

func resumeWorktreeBytes(t *testing.T, root string) string {
	t.Helper()
	var snapshot strings.Builder
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == filepath.Join(root, ".git") && entry.IsDir() {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(&snapshot, "%s %o ", filepath.ToSlash(rel), info.Mode())
		if info.Mode().IsRegular() {
			body, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Fprintf(&snapshot, "%x", body)
		} else if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			snapshot.WriteString(target)
		}
		snapshot.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot.String()
}
