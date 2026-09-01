package worktree

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// The FT169 landing surface: one reviewed source must not pay six refusal
// round-trips. Each subtest drives the real land command through the public
// fixture and asserts the exact symptom the 2026-08-22 landing paid for.

func landSurface(t *testing.T, request string) (string, Creation, string, string) {
	t.Helper()
	root, creation, base, tip, _, _ := publicLandingFixture(t, request, "", "")
	return root, creation, base, tip
}

func landIn(t *testing.T, root string, args []string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := LandCommand(root, Home(), "", args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// seedCaptureBase commits capture files on the destination and rebases the
// source onto them, so both sides share the files in their merge base.
func seedCaptureBase(t *testing.T, root, source string, files map[string]string) string {
	t.Helper()
	for name, body := range files {
		mustMkdirAll(t, filepath.Dir(filepath.Join(root, filepath.FromSlash(name))), 0o755)
		mustWrite(t, filepath.Join(root, filepath.FromSlash(name)), []byte(body), 0o644)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "seed capture")
	gitRun(t, source, "rebase", "main")
	return gitOutput(t, root, "rev-parse", "HEAD")
}

func TestLandCommandAcceptsAbbreviatedIdentities(t *testing.T) {
	t.Parallel()
	request := "land-surface-abbreviated"
	root, creation, base, tip := landSurface(t, request)
	code, stdout, stderr := landIn(t, root, landArgs(request, base[:12], tip[:12], creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") || !strings.Contains(stdout, "source_base="+base+",source_tip="+tip+",") {
		t.Fatalf("abbreviated landing = (%d, %q, %q), want released with full identities", code, stdout, stderr)
	}
	parents := strings.Fields(gitOutput(t, root, "rev-list", "--parents", "-n", "1", "main"))
	if len(parents) != 3 || parents[1] != base || parents[2] != tip {
		t.Fatalf("published parents = %q, want %s and %s", parents, base, tip)
	}
}

func TestLandCommandComposesCaptureOntoMovedDestination(t *testing.T) {
	t.Parallel()
	request := "land-surface-capture-conflict"
	root, creation, _, _ := landSurface(t, request)
	base := seedCaptureBase(t, root, creation.Path, map[string]string{
		"capture/session-handoff.md": "handoff base\n",
		"capture/learnings.md":       "learnings base\n",
	})
	commitInWorktree(t, creation.Path, "capture/session-handoff.md", "handoff source\n", "source handoff")
	commitInWorktree(t, creation.Path, "capture/learnings.md", "learnings source\n", "source learnings")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, "capture/session-handoff.md", "handoff destination\n", "destination handoff")
	commitInWorktree(t, root, "capture/learnings.md", "learnings destination\n", "destination learnings")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") {
		t.Fatalf("capture-conflict landing = (%d, %q, %q), want released", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "show", "main:capture/session-handoff.md"); got != "handoff source" {
		t.Fatalf("published handoff = %q, want the source's", got)
	}
	// `git merge-file --union` publishes the base lines, then the destination side,
	// then the source side. Both sides replaced the one base line, so only the two
	// appended lines remain, in that order.
	published, err := descendant(t, "git", "-C", root, "show", "main:capture/learnings.md").Output()
	if want := "learnings destination\nlearnings source\n"; err != nil || string(published) != want {
		t.Fatalf("published learnings = %q (%v), want exactly %q", published, err, want)
	}
}

// WL19: the landing discloses each settled phase-owned path with its verb, so a
// union the merge did not decide is visible on stderr rather than silent.
func TestLandCommandDisclosesAUnionResolution(t *testing.T) {
	t.Parallel()
	request := "land-surface-union-disclosure"
	root, creation, _, _ := landSurface(t, request)
	base := seedCaptureBase(t, root, creation.Path, map[string]string{"capture/learnings.md": "learnings base\n"})
	commitInWorktree(t, creation.Path, "capture/learnings.md", "learnings base\nlearnings source\n", "source learnings")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, "capture/learnings.md", "learnings base\nlearnings destination\n", "destination learnings")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") {
		t.Fatalf("union landing = (%d, %q, %q), want released", code, stdout, stderr)
	}
	lines := 0
	for _, line := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(line, "landing composition{resolved=") {
			lines++
			if !strings.Contains(line, "capture/learnings.md:union") {
				t.Fatalf("disclosure = %q, want capture/learnings.md settled by union", line)
			}
		}
	}
	if lines != 1 {
		t.Fatalf("disclosure lines = %d in %q, want exactly one", lines, stderr)
	}
}

func TestLandCommandAuthorizesCaptureOutsideTheFence(t *testing.T) {
	t.Parallel()
	request := "land-surface-capture-fence"
	root, creation, base, _ := landSurface(t, request)
	mustMkdirAll(t, filepath.Join(creation.Path, "capture"), 0o755)
	commitInWorktree(t, creation.Path, "capture/learnings.md", "learning\n", "phase-owned learning")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") {
		t.Fatalf("capture landing = (%d, %q, %q), want released", code, stdout, stderr)
	}
	if got := gitOutput(t, root, "show", "main:capture/learnings.md"); got != "learning" {
		t.Fatalf("published learning = %q", got)
	}
}

func TestLandCommandReportsEveryRefusalInOnePreflight(t *testing.T) {
	t.Parallel()
	request := "land-surface-one-preflight"
	root, creation, base, tip := landSurface(t, request)
	mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
	mustWrite(t, filepath.Join(creation.Path, "scratch"), []byte("scratch\n"), 0o600)
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "landing destination is not clean") || !strings.Contains(stdout, "reviewed source is not clean") {
		t.Fatalf("two-refusal preflight = (%d, %q, %q), want both refusals named", code, stdout, stderr)
	}
	// LRS3: every landing-preflight route ends with the caller's own re-run, so a repair
	// does not cost the operator the flags it passed.
	tail := "; then bench worktree land --request '" + request + "' --base '" + base +
		"' --source-tip '" + tip + "' --spec 'x' -m <message> '" + creation.Path + "'"
	for _, face := range []string{faceDestinationNotClean, faceSourceNotClean} {
		next, printed := landingFaceNext(stdout, landingRefusalFaceByName(face).detail)
		if !printed || !strings.HasSuffix(next, tail) {
			t.Fatalf("%s next = %q (printed=%t) in %q, want a repair ending %q", face, next, printed, stdout, tail)
		}
	}
}

// TestLandCommandReportsIdentityAndDestinationInOnePreflight is LR10. The destination
// proof and the identity proof are independent, so one run has to name both; a
// first-refusal-exits rewrite would hide the second for a whole run.
func TestLandCommandReportsIdentityAndDestinationInOnePreflight(t *testing.T) {
	t.Parallel()
	request := "land-surface-identity-preflight"
	root, creation, base, tip := landSurface(t, request)
	mustWrite(t, filepath.Join(root, "dirty"), []byte("dirty\n"), 0o600)
	code, stdout, stderr := landIn(t, root, landArgs("unknown-request", base, tip, creation.Path))
	both := strings.Contains(stdout, "refused{detail=landing destination is not clean") &&
		strings.Contains(stdout, "refused{detail=request token matches no assignment")
	if code != 1 || !both || strings.Count(stdout, "refused{") != 2 {
		t.Fatalf("identity-and-destination preflight = (%d, %q, %q), want exactly two refusals", code, stdout, stderr)
	}
	// LRS20: the destination proof refuses before the assignment resolves, so its re-run
	// addresses the operator's own worktree path rather than an assignment id.
	next, printed := landingFaceNext(stdout, landingRefusalFaceByName(faceDestinationNotClean).detail)
	if !printed || strings.Contains(next, "bench worktree exec") || !strings.HasSuffix(next, " '"+creation.Path+"'") {
		t.Fatalf("destination next = %q (printed=%t) in %q, want a re-run ending with the operator's own path", next, printed, stdout)
	}
}

func TestLandCommandFenceRefusalNamesThePath(t *testing.T) {
	t.Parallel()
	request := "land-surface-fence-path"
	root, creation, base, _ := landSurface(t, request)
	commitInWorktree(t, creation.Path, "stray.txt", "stray\n", "out of fence")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "paths_total=1") || !strings.Contains(stdout, "refusal_paths[1]{path}:\n  stray.txt\n") {
		t.Fatalf("fence refusal = (%d, %q, %q), want the unfenced path in a refusal_paths row", code, stdout, stderr)
	}
}

// WL4 and WL21: the fence rides with the spec. A spec-backed landing still refuses a
// path no fence names; the same source lands when no spec names a fence.
func TestLandCommandFenceRidesWithTheSpec(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		specArg bool
	}{
		{name: "spec-backed", specArg: true},
		{name: "spec-less"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := "land-surface-fence-" + tc.name
			root, creation, base, _, _, _ := specLessLandingFixture(t, request)
			commitInWorktree(t, creation.Path, "stray.txt", "stray\n", "out of fence")
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
			args := specLessLandArgs(request, base, tip, creation.Path)
			if tc.specArg {
				args = landArgs(request, base, tip, creation.Path)
			}
			code, stdout, stderr := landIn(t, root, args)
			if tc.specArg {
				if code != 1 || !strings.Contains(stdout, "ownership fence is invalid") || !strings.Contains(stdout, "stray.txt") {
					t.Fatalf("spec-backed fence refusal = (%d, %q, %q), want the offending path named", code, stdout, stderr)
				}
				return
			}
			if code != 0 || !strings.Contains(stdout, "worktree=released,census=0}") {
				t.Fatalf("spec-less out-of-fence landing = (%d, %q, %q), want released", code, stdout, stderr)
			}
			if got := gitOutput(t, root, "show", "main:stray.txt"); got != "stray" {
				t.Fatalf("published stray path = %q", got)
			}
		})
	}
}

func TestLandCommandConflictRefusalNamesThePath(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-path"
	root, creation, base, tip := landSurface(t, request)
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "composition conflict: textual") || !strings.Contains(stdout, "owned.txt") {
		t.Fatalf("conflict refusal = (%d, %q, %q), want the conflicted path named", code, stdout, stderr)
	}
}

// WL16: a board file is outside the rule table, so a conflict on ROADMAP.md refuses
// and names the path the repair starts from.
func TestLandCommandConflictOnTheBoardNamesTheBoardPath(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-roadmap"
	root, creation, base, _, _, _ := specLessLandingFixture(t, request)
	commitInWorktree(t, creation.Path, "ROADMAP.md", "board source\n", "source board")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, "ROADMAP.md", "board destination\n", "destination board")
	code, stdout, stderr := landIn(t, root, specLessLandArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "composition conflict: textual") || !strings.Contains(stdout, "ROADMAP.md") {
		t.Fatalf("board conflict refusal = (%d, %q, %q), want ROADMAP.md named", code, stdout, stderr)
	}
}

// WL18: the conflict refusal names the source repair in order, and the re-run carries
// the destination as the new base.
func TestLandCommandConflictRefusalNamesTheSourceRepair(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-repair"
	root, creation, base, tip := landSurface(t, request)
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	wantNext := "next=git -C '" + creation.Path + "' merge '" + destination +
		"' (bench worktree merge refuses this conflict; resolve it by hand); then bench commit; then /bench-review-implementation; then " +
		"bench worktree land --request <request> --base '" + destination +
		"' --source-tip <repaired-source-tip> --spec 'x' -m <message> '" + creation.Path + "'}"
	if code != 1 || !strings.Contains(stdout, "composition conflict: textual") || !strings.Contains(stdout, wantNext) {
		t.Fatalf("conflict repair next = (%d, %q, %q), want %q", code, stdout, stderr, wantNext)
	}
	if strings.Contains(stdout, "no Bench verb") {
		t.Fatalf("conflict next still denies the merge verb: %q", stdout)
	}
	if strings.Contains(stdout, request) {
		t.Fatalf("conflict next leaked the caller token: %q", stdout)
	}
}

// A spec-less landing re-runs spec-less, so its conflict next names no --spec.
func TestLandCommandSpecLessConflictNextNamesNoSpec(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-spec-less"
	root, creation, base, _, _, _ := specLessLandingFixture(t, request)
	commitInWorktree(t, creation.Path, "ROADMAP.md", "board source\n", "source board")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, "ROADMAP.md", "board destination\n", "destination board")
	code, stdout, stderr := landIn(t, root, specLessLandArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, " --source-tip <repaired-source-tip> -m <message> '") || strings.Contains(stdout, "--spec") {
		t.Fatalf("spec-less conflict next = (%d, %q, %q), want no --spec", code, stdout, stderr)
	}
}

// Edge under WL16: a conflicted path that carries a control byte renders through the
// sanitized paths table, so no raw control byte reaches the terminal.
func TestLandCommandConflictOnAControlBytePathRendersSanitized(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-control-byte"
	root, creation, base, _, _, _ := specLessLandingFixture(t, request)
	name := "board\x1bfile.md"
	commitInWorktree(t, creation.Path, name, "source bytes\n", "source control-byte path")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, name, "destination bytes\n", "destination control-byte path")
	code, stdout, stderr := landIn(t, root, specLessLandArgs(request, base, tip, creation.Path))
	if code != 1 || !strings.Contains(stdout, "composition conflict: textual") {
		t.Fatalf("control-byte conflict = (%d, %q, %q), want a refusal", code, stdout, stderr)
	}
	if strings.ContainsRune(stdout, '\x1b') || !strings.Contains(stdout, `"board\\u001bfile.md"`) {
		t.Fatalf("control-byte path render = %q, want the escaped path and no raw control byte", stdout)
	}
}

// Edge under WL18: a source worktree path that is not line-safe cannot be pasted, so
// both repair steps that address it take the assignment pointer form.
func TestLandCommandConflictNextPointsThroughUnsafePath(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-unsafe-path"
	home := filepath.Join(t.TempDir(), "bench\n\x1bhome")
	root, creation, base, tip, _ := publicLandingFixtureAtHome(t, request, "", "", home)
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	code, stdout, stderr := landIn(t, root, landArgs(request, base, tip, creation.Path))
	wantNext := "next=bench worktree exec " + creation.Assignment.ID + " -- git merge '" + destination +
		"' (bench worktree merge refuses this conflict; resolve it by hand); then bench commit; then /bench-review-implementation; then " +
		"bench worktree exec " + creation.Assignment.ID + " -- bench worktree land --request <request> --base '" + destination +
		"' --source-tip <repaired-source-tip> --spec 'x' -m <message> .}"
	if code != 1 || strings.ContainsRune(stdout, '\x1b') || !strings.Contains(stdout, wantNext) {
		t.Fatalf("unsafe-path conflict next = (%d, %q, %q), want the pointer form %q", code, stdout, stderr, wantNext)
	}
	if strings.Contains(stdout, "no Bench verb") {
		t.Fatalf("unsafe-path conflict next still denies the merge verb: %q", stdout)
	}
}

// Edge under WL18: a --spec slug that is not line-safe cannot be pasted, so the
// conflict next carries the `<spec>` placeholder and no raw control byte. A
// tickets-only folder is the one spec shape whose name reaches the landing verbatim.
func TestLandCommandConflictNextPlaceholdsAnUnsafeSpec(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-unsafe-spec"
	root, creation, base, _ := landSurface(t, request)
	slug := "close\x1bme"
	mustMkdirAll(t, filepath.Join(creation.Path, "specs", slug, "tickets"), 0o755)
	commitInWorktree(t, creation.Path, filepath.Join("specs", slug, "tickets", "one.md"), "Ticket.\n", "tickets-only folder")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, root, "rev-parse", "HEAD")
	args := []string{"--request", request, "--base", base, "--source-tip", tip, "--spec", slug, "-m", "land", creation.Path}
	code, stdout, stderr := landIn(t, root, args)
	wantNext := "bench worktree land --request <request> --base '" + destination +
		"' --source-tip <repaired-source-tip> --spec <spec> -m <message> '" + creation.Path + "'}"
	if code != 1 || !strings.Contains(stdout, "composition conflict: textual") || !strings.Contains(stdout, wantNext) {
		t.Fatalf("unsafe-spec conflict next = (%d, %q, %q), want the placeholder form %q", code, stdout, stderr, wantNext)
	}
	if strings.Contains(stdout, "no Bench verb") {
		t.Fatalf("unsafe-spec conflict next still denies the merge verb: %q", stdout)
	}
	if strings.ContainsRune(stdout, '\x1b') {
		t.Fatalf("unsafe-spec conflict next leaked a raw control byte: %q", stdout)
	}
}
