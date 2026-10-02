package worktree

import (
	"path/filepath"
	"strings"
	"testing"
)

// The FT169 landing surface: one reviewed source must not pay six refusal
// round-trips. Each subtest drives the real land command through the public
// fixture and asserts the exact symptom the 2026-08-22 landing paid for.

func landSurface(t *testing.T, request string) landingFixture {
	t.Helper()
	return publicLandingFixture(t, request, "", "")
}

// processHomeCall builds a land call at the fixture's root and the process's Bench home,
// which is the home each landing surface test lands at. It runs no verb.
func (f landingFixture) processHomeCall(args ...string) verbCall {
	return repoHome{f.root, Home()}.call(args...)
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
	f := landSurface(t, request)
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base[:12], f.tip[:12], f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") || !strings.Contains(r.stdout, "source_base="+f.base+",source_tip="+f.tip+",") {
		t.Fatalf("abbreviated landing = (%d, %q, %q), want released with full identities", r.exit, r.stdout, r.stderr)
	}
	parents := strings.Fields(gitOutput(t, f.root, "rev-list", "--parents", "-n", "1", "main"))
	if len(parents) != 3 || parents[1] != f.base || parents[2] != f.tip {
		t.Fatalf("published parents = %q, want %s and %s", parents, f.base, f.tip)
	}
}

func TestLandCommandComposesCaptureOntoMovedDestination(t *testing.T) {
	t.Parallel()
	request := "land-surface-capture-conflict"
	f := landSurface(t, request)
	base := seedCaptureBase(t, f.root, f.creation.Path, map[string]string{
		"capture/session-handoff.md": "handoff base\n",
		"capture/learnings.md":       "learnings base\n",
	})
	commitInWorktree(t, f.creation.Path, "capture/session-handoff.md", "handoff source\n", "source handoff")
	commitInWorktree(t, f.creation.Path, "capture/learnings.md", "learnings source\n", "source learnings")
	refreshLandingEvidence(t, f.creation.Path, base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.root, "capture/session-handoff.md", "handoff destination\n", "destination handoff")
	commitInWorktree(t, f.root, "capture/learnings.md", "learnings destination\n", "destination learnings")
	destination := gitOutput(t, f.root, "rev-parse", "main")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "completion composition changes capture/learnings.md") || gitOutput(t, f.root, "rev-parse", "main") != destination {
		t.Fatalf("unreviewed capture composition = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	tip = foldCompletionComposition(t, f.root, f.creation.Path, base, tip, destination)
	r = runVerb(t, verbLand, f.processHomeCall(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("capture-conflict landing = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "show", "main:capture/session-handoff.md"); got != "handoff source" {
		t.Fatalf("published handoff = %q, want the source's", got)
	}
	// `git merge-file --union` publishes the base lines, then the destination side,
	// then the source side. Both sides replaced the one base line, so only the two
	// appended lines remain, in that order.
	published, err := descendant(t, "git", "-C", f.root, "show", "main:capture/learnings.md").Output()
	if want := "learnings destination\nlearnings source\n"; err != nil || string(published) != want {
		t.Fatalf("published learnings = %q (%v), want exactly %q", published, err, want)
	}
}

// WL19: the landing discloses each settled phase-owned path with its verb, so a
// union the merge did not decide is visible on stderr rather than silent.
func TestLandCommandDisclosesAUnionResolution(t *testing.T) {
	t.Parallel()
	request := "land-surface-union-disclosure"
	f := landSurface(t, request)
	base := seedCaptureBase(t, f.root, f.creation.Path, map[string]string{"capture/learnings.md": "learnings base\n"})
	commitInWorktree(t, f.creation.Path, "capture/learnings.md", "learnings base\nlearnings source\n", "source learnings")
	refreshLandingEvidence(t, f.creation.Path, base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.root, "capture/learnings.md", "learnings base\nlearnings destination\n", "destination learnings")
	destination := gitOutput(t, f.root, "rev-parse", "main")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "completion composition changes capture/learnings.md") || gitOutput(t, f.root, "rev-parse", "main") != destination {
		t.Fatalf("unreviewed union = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	lines := 0
	for _, line := range strings.Split(r.stderr, "\n") {
		if strings.HasPrefix(line, "landing composition{resolved=") {
			lines++
			if !strings.Contains(line, "capture/learnings.md:union") {
				t.Fatalf("disclosure = %q, want capture/learnings.md settled by union", line)
			}
		}
	}
	if lines != 1 {
		t.Fatalf("disclosure lines = %d in %q, want exactly one", lines, r.stderr)
	}
	tip = foldCompletionComposition(t, f.root, f.creation.Path, base, tip, destination)
	r = runVerb(t, verbLand, f.processHomeCall(landArgs(request, base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("reviewed union landing = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
}

func TestLandCommandAuthorizesCaptureOutsideTheFence(t *testing.T) {
	t.Parallel()
	request := "land-surface-capture-fence"
	f := landSurface(t, request)
	mustMkdirAll(t, filepath.Join(f.creation.Path, "capture"), 0o755)
	commitInWorktree(t, f.creation.Path, "capture/learnings.md", "learning\n", "phase-owned learning")
	refreshLandingEvidence(t, f.creation.Path, f.base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("capture landing = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "show", "main:capture/learnings.md"); got != "learning" {
		t.Fatalf("published learning = %q", got)
	}
}

func TestLandCommandReportsEveryRefusalInOnePreflight(t *testing.T) {
	t.Parallel()
	request := "land-surface-one-preflight"
	f := landSurface(t, request)
	mustWrite(t, filepath.Join(f.root, "tracked.txt"), []byte("dirty\n"), 0o644)
	mustWrite(t, filepath.Join(f.creation.Path, "scratch"), []byte("scratch\n"), 0o600)
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "landing destination is not clean") || !strings.Contains(r.stdout, "reviewed source is not clean") {
		t.Fatalf("two-refusal preflight = (%d, %q, %q), want both refusals named", r.exit, r.stdout, r.stderr)
	}
	// LRS3: every landing-preflight route ends with the caller's own re-run, so a repair
	// does not cost the operator the flags it passed. The registry is the face set, so the
	// walk covers each preflight face: the two this run raises are read from the printed
	// record, and the rest from the route the face composes over the same re-run.
	rerun := "bench worktree land --request '" + request + "' --base '" + f.base +
		"' --source-tip '" + f.tip + "' --spec 'x' -m <message> '" + f.creation.Path + "'"
	tail := "; then " + rerun
	for _, face := range landingRefusalFaces {
		if face.stage != stagePreflight {
			continue
		}
		if next, printed := landingFaceNext(r.stdout, face.detail); printed {
			if !strings.HasSuffix(next, tail) {
				t.Fatalf("%s next = %q in %q, want a repair ending %q", face.name, next, r.stdout, tail)
			}
			continue
		}
		if route := face.route(rerun); !strings.HasSuffix(route, rerun) {
			t.Fatalf("%s route = %q, want it to end with the caller's own re-run %q", face.name, route, rerun)
		}
	}
}

// TestLandCommandReportsIdentityAndDestinationInOnePreflight is LR10. The destination
// proof and the identity proof are independent, so one run has to name both; a
// first-refusal-exits rewrite would hide the second for a whole run.
func TestLandCommandReportsIdentityAndDestinationInOnePreflight(t *testing.T) {
	t.Parallel()
	request := "land-surface-identity-preflight"
	f := landSurface(t, request)
	mustWrite(t, filepath.Join(f.root, "tracked.txt"), []byte("dirty\n"), 0o644)
	r := runVerb(t, verbLand, f.processHomeCall(landArgs("unknown-request", f.base, f.tip, f.creation.Path)...))
	both := strings.Contains(r.stdout, "refused{detail=landing destination is not clean") &&
		strings.Contains(r.stdout, "refused{detail=request token matches no assignment")
	if r.exit != 1 || !both || strings.Count(r.stdout, "refused{") != 2 {
		t.Fatalf("identity-and-destination preflight = (%d, %q, %q), want exactly two refusals", r.exit, r.stdout, r.stderr)
	}
	// LRS20: the destination proof refuses before the assignment resolves, so its re-run
	// addresses the operator's own worktree path rather than an assignment id.
	next, printed := landingFaceNext(r.stdout, landingRefusalFaceByName(faceDestinationNotClean).detail)
	if !printed || strings.Contains(next, "bench worktree exec") || !strings.HasSuffix(next, " '"+f.creation.Path+"'") {
		t.Fatalf("destination next = %q (printed=%t) in %q, want a re-run ending with the operator's own path", next, printed, r.stdout)
	}
	// LRS9: the assignment fault stopped the source proofs of its own group, so its route
	// says so and the operator expects a second refusal after the repair. The destination
	// group ended at its own fault, so its route says no such thing.
	if strings.Contains(next, laterProofsSkipped) {
		t.Fatalf("destination next = %q, want no skipped-proof sentence from a group that ran to its end", next)
	}
	assignmentNext, printed := landingFaceNext(r.stdout, "request token matches no assignment")
	if !printed || !strings.HasPrefix(assignmentNext, laterProofsSkipped+"; ") {
		t.Fatalf("assignment next = %q (printed=%t) in %q, want the skipped-proof sentence ahead of the repair", assignmentNext, printed, r.stdout)
	}
}

func TestLandCommandFenceRefusalNamesThePath(t *testing.T) {
	t.Parallel()
	request := "land-surface-fence-path"
	f := landSurface(t, request)
	commitInWorktree(t, f.creation.Path, "stray.txt", "stray\n", "out of fence")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "paths_total=1") || !strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:\n  stray.txt\n") {
		t.Fatalf("fence refusal = (%d, %q, %q), want the unfenced path in a refusal_paths row", r.exit, r.stdout, r.stderr)
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
			f := specLessLandingFixture(t, request)
			commitInWorktree(t, f.creation.Path, "stray.txt", "stray\n", "out of fence")
			tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
			args := specLessLandArgs(request, f.base, tip, f.creation.Path)
			if tc.specArg {
				args = landArgs(request, f.base, tip, f.creation.Path)
			}
			r := runVerb(t, verbLand, f.processHomeCall(args...))
			if tc.specArg {
				if r.exit != 1 || !strings.Contains(r.stdout, "ownership fence is invalid") || !strings.Contains(r.stdout, "stray.txt") {
					t.Fatalf("spec-backed fence refusal = (%d, %q, %q), want the offending path named", r.exit, r.stdout, r.stderr)
				}
				return
			}
			if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
				t.Fatalf("spec-less out-of-fence landing = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
			}
			if got := gitOutput(t, f.root, "show", "main:stray.txt"); got != "stray" {
				t.Fatalf("published stray path = %q", got)
			}
		})
	}
}

func TestLandCommandConflictRefusalNamesThePath(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-path"
	f := landSurface(t, request)
	commitInWorktree(t, f.root, "owned.txt", "destination bytes\n", "destination conflict")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "composition conflict: textual") || !strings.Contains(r.stdout, "owned.txt") {
		t.Fatalf("conflict refusal = (%d, %q, %q), want the conflicted path named", r.exit, r.stdout, r.stderr)
	}
}

// WL16: a board file is outside the rule table, so a conflict on ROADMAP.md refuses
// and names the path the repair starts from.
func TestLandCommandConflictOnTheBoardNamesTheBoardPath(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-roadmap"
	f := specLessLandingFixture(t, request)
	commitInWorktree(t, f.creation.Path, "ROADMAP.md", "board source\n", "source board")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.root, "ROADMAP.md", "board destination\n", "destination board")
	r := runVerb(t, verbLand, f.processHomeCall(specLessLandArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "composition conflict: textual") || !strings.Contains(r.stdout, "ROADMAP.md") {
		t.Fatalf("board conflict refusal = (%d, %q, %q), want ROADMAP.md named", r.exit, r.stdout, r.stderr)
	}
}

// WL18: the conflict refusal names the source repair in order, and the re-run carries
// the destination as the new base.
func TestLandCommandConflictRefusalNamesTheSourceRepair(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-repair"
	f := landSurface(t, request)
	commitInWorktree(t, f.root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, f.root, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, f.tip, f.creation.Path)...))
	wantNext := "next=git -C '" + f.creation.Path + "' merge '" + destination +
		"' (bench worktree merge refuses this conflict; resolve it by hand); then bench commit; then /bench-review-implementation; then " +
		"bench worktree land --request <request> --base '" + destination +
		"' --source-tip <repaired-source-tip> --spec 'x' -m <message> '" + f.creation.Path + "'}"
	if r.exit != 1 || !strings.Contains(r.stdout, "composition conflict: textual") || !strings.Contains(r.stdout, wantNext) {
		t.Fatalf("conflict repair next = (%d, %q, %q), want %q", r.exit, r.stdout, r.stderr, wantNext)
	}
	if strings.Contains(r.stdout, "no Bench verb") {
		t.Fatalf("conflict next still denies the merge verb: %q", r.stdout)
	}
	if strings.Contains(r.stdout, request) {
		t.Fatalf("conflict next leaked the caller token: %q", r.stdout)
	}
}

// A spec-less landing re-runs spec-less, so its conflict next names no --spec.
func TestLandCommandSpecLessConflictNextNamesNoSpec(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-spec-less"
	f := specLessLandingFixture(t, request)
	commitInWorktree(t, f.creation.Path, "ROADMAP.md", "board source\n", "source board")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.root, "ROADMAP.md", "board destination\n", "destination board")
	r := runVerb(t, verbLand, f.processHomeCall(specLessLandArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, " --source-tip <repaired-source-tip> -m <message> '") || strings.Contains(r.stdout, "--spec") {
		t.Fatalf("spec-less conflict next = (%d, %q, %q), want no --spec", r.exit, r.stdout, r.stderr)
	}
}

// Edge under WL16: a conflicted path that carries a control byte renders through the
// sanitized paths table, so no raw control byte reaches the terminal.
func TestLandCommandConflictOnAControlBytePathRendersSanitized(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-control-byte"
	f := specLessLandingFixture(t, request)
	name := "board\x1bfile.md"
	commitInWorktree(t, f.creation.Path, name, "source bytes\n", "source control-byte path")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.root, name, "destination bytes\n", "destination control-byte path")
	r := runVerb(t, verbLand, f.processHomeCall(specLessLandArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "composition conflict: textual") {
		t.Fatalf("control-byte conflict = (%d, %q, %q), want a refusal", r.exit, r.stdout, r.stderr)
	}
	if strings.ContainsRune(r.stdout, '\x1b') || !strings.Contains(r.stdout, `"board\\u001bfile.md"`) {
		t.Fatalf("control-byte path render = %q, want the escaped path and no raw control byte", r.stdout)
	}
}

// Edge under WL18: a source worktree path that is not line-safe cannot be pasted, so
// both repair steps that address it take the assignment pointer form.
func TestLandCommandConflictNextPointsThroughUnsafePath(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-unsafe-path"
	home := filepath.Join(t.TempDir(), "bench\n\x1bhome")
	f := publicLandingFixtureAtHome(t, request, "", "", home)
	commitInWorktree(t, f.root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, f.root, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, f.tip, f.creation.Path)...))
	wantNext := "next=bench worktree exec " + f.creation.Assignment.ID + " -- git merge '" + destination +
		"' (bench worktree merge refuses this conflict; resolve it by hand); then bench commit; then /bench-review-implementation; then " +
		"bench worktree exec " + f.creation.Assignment.ID + " -- bench worktree land --request <request> --base '" + destination +
		"' --source-tip <repaired-source-tip> --spec 'x' -m <message> .}"
	if r.exit != 1 || strings.ContainsRune(r.stdout, '\x1b') || !strings.Contains(r.stdout, wantNext) {
		t.Fatalf("unsafe-path conflict next = (%d, %q, %q), want the pointer form %q", r.exit, r.stdout, r.stderr, wantNext)
	}
	if strings.Contains(r.stdout, "no Bench verb") {
		t.Fatalf("unsafe-path conflict next still denies the merge verb: %q", r.stdout)
	}
}

// Edge under WL18: a --spec slug that is not line-safe cannot be pasted, so the
// conflict next carries the `<spec>` placeholder and no raw control byte. A
// tickets-only folder is the one spec shape whose name reaches the landing verbatim.
func TestLandCommandConflictNextPlaceholdsAnUnsafeSpec(t *testing.T) {
	t.Parallel()
	request := "land-surface-conflict-unsafe-spec"
	f := landSurface(t, request)
	slug := "close\x1bme"
	mustMkdirAll(t, filepath.Join(f.creation.Path, "specs", slug, "tickets"), 0o755)
	commitInWorktree(t, f.creation.Path, filepath.Join("specs", slug, "tickets", "one.md"), "Ticket.\n", "tickets-only folder")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	commitInWorktree(t, f.root, "owned.txt", "destination bytes\n", "destination conflict")
	destination := gitOutput(t, f.root, "rev-parse", "HEAD")
	args := []string{"--request", request, "--base", f.base, "--source-tip", tip, "--spec", slug, "-m", "land", f.creation.Path}
	r := runVerb(t, verbLand, f.processHomeCall(args...))
	wantNext := "bench worktree land --request <request> --base '" + destination +
		"' --source-tip <repaired-source-tip> --spec <spec> -m <message> '" + f.creation.Path + "'}"
	if r.exit != 1 || !strings.Contains(r.stdout, "composition conflict: textual") || !strings.Contains(r.stdout, wantNext) {
		t.Fatalf("unsafe-spec conflict next = (%d, %q, %q), want the placeholder form %q", r.exit, r.stdout, r.stderr, wantNext)
	}
	if strings.Contains(r.stdout, "no Bench verb") {
		t.Fatalf("unsafe-spec conflict next still denies the merge verb: %q", r.stdout)
	}
	if strings.ContainsRune(r.stdout, '\x1b') {
		t.Fatalf("unsafe-spec conflict next leaked a raw control byte: %q", r.stdout)
	}
}
