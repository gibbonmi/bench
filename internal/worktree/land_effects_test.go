// Promotion-broker tests for the landing command: the refresh effect that republishes the
// destination's own broker after publication.
package worktree

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
)

// writeGoMainFixture writes a resolvable Go main package and its build script, for a
// fixture whose root must hold Go code. A broker fixture's destination stands on it,
// because the seal's source digest is the build-input closure of that package. A merge
// fixture's target stands on it too, because the built-in lane needs a Go package to grade.
func writeGoMainFixture(t *testing.T, root string) {
	t.Helper()
	mustWrite(t, filepath.Join(root, "go.mod"), []byte("module benchfixture\n\ngo 1.22\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "cmd", "bench"), 0o755)
	mustWrite(t, filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "scripts"), 0o755)
	mustWrite(t, filepath.Join(root, "scripts", "go-build.sh"), []byte("#!/bin/sh\nexit 0\n"), 0o755)
}

// projectGreenMarker reads the destination's project-green marker, answering the empty
// string when no marker is recorded. An absent marker is an ordinary state before the
// first landing, so it is a value here rather than a test failure.
func projectGreenMarker(t *testing.T, root string) string {
	t.Helper()
	output, err := descendant(t, "git", "-C", root, "rev-parse", "--verify", "--quiet", "refs/bench/green/main").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// wantEffects is the two-row effects table a landing prints for a refresh result and the
// cleanup result beside it. It composes the cleanup word from the owner's own constants,
// so a rename there moves the expectation with it. It spells the row header, the row
// layout, and the caller's refresh word here, so a drift in any of the three fails this
// comparison.
//
// A caller that states the refresh word alone takes the cleanup word its own fixture
// expects: a failed refresh starts no cleanup, and a fixture that carries no sibling of
// its own settles an empty set. A caller whose fixture expects another word states it.
func wantEffects(refresh string, cleanup ...string) string {
	result := effectComplete
	if refresh == effectFailed {
		result = effectPending
	}
	if len(cleanup) > 0 {
		result = cleanup[0]
	}
	return "effects[2]{effect,result}:\n  refresh," + refresh + "\n  cleanup," + result + "\n"
}

// brokerDestinationFixture is the public landing fixture whose destination declares Bench
// build inputs and ignores the directory its published executable lands in. The source
// worktree's own ignored residue is removed, so the release settles and the effects run.
func brokerDestinationFixture(t *testing.T, request string) landingFixture {
	t.Helper()
	f := publicLandingFixture(t, request, "dist/bench", "dist/")
	if err := os.Remove(filepath.Join(f.creation.Path, "dist", "bench")); err != nil {
		t.Fatal(err)
	}
	writeGoMainFixture(t, f.root)
	gitRun(t, f.root, "add", ".")
	gitRun(t, f.root, withIdentity("commit", "-qm", "broker sources")...)
	commitLandingBuildInputs(t, f.root)
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	gitRun(t, f.creation.Path, "rebase", "main")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	return f
}

// refreshJoins is the landing seam set whose refresh build is install, beside the count
// of its calls. A fixture answers the build itself, so no row waits on a real compile.
func refreshJoins(install func(root, executable string) error) (joins, *int) {
	calls := 0
	j := defaultJoins()
	j.buildSubject = func(_ context.Context, root, executable string) error {
		calls++
		if install == nil {
			return nil
		}
		return install(root, executable)
	}
	return j, &calls
}

// publishVerifyingBroker installs an executable and its seal at root's published path, as
// one transaction, so the destination verifies against its own sources.
func publishVerifyingBroker(t *testing.T, root, executable string) error {
	t.Helper()
	staged := filepath.Join(t.TempDir(), "staged-bench")
	mustWrite(t, staged, []byte("#!/bin/sh\nexit 0\n"), 0o755)
	mustMkdirAll(t, filepath.Dir(executable), 0o755)
	return freshness.Publish(root, staged, executable, filepath.Dir(executable), "0.0.0-fixture")
}

// TestLandSkipsTheRefreshWithoutBuildInputs is LC2. A destination that declares no Bench
// build inputs has no broker of its own to refresh, so the effect reports not-applicable
// and the build seam is never reached.
func TestLandSkipsTheRefreshWithoutBuildInputs(t *testing.T) {
	t.Parallel()
	request := "land-refresh-no-inputs"
	f := publicLandingFixture(t, request, "", "")
	j, calls := refreshJoins(nil)

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 0 || !strings.Contains(r.stdout, wantEffects("not-applicable")) {
		t.Fatalf("landing without build inputs = (%d, %q, %q), want a not-applicable refresh", r.exit, r.stdout, r.stderr)
	}
	if *calls != 0 {
		t.Fatalf("refresh build calls = %d, want none", *calls)
	}
}

// TestLandSkipsAFreshBroker is LC3. A destination whose published executable already
// verifies against its own sources is refreshed by the publication itself, so the effect
// completes with no compile.
func TestLandSkipsAFreshBroker(t *testing.T) {
	t.Parallel()
	request := "land-refresh-already-fresh"
	f := brokerDestinationFixture(t, request)
	if err := publishVerifyingBroker(t, f.root, freshness.PublishedExecutable(f.root)); err != nil {
		t.Fatal(err)
	}
	j, calls := refreshJoins(nil)

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 0 || !strings.Contains(r.stdout, wantEffects("complete")) {
		t.Fatalf("landing with a fresh broker = (%d, %q, %q), want a complete refresh", r.exit, r.stdout, r.stderr)
	}
	if *calls != 0 {
		t.Fatalf("refresh build calls = %d, want none", *calls)
	}
}

// TestLandRefreshesTheBrokerAfterPublication is LC1. A destination that declares Bench
// build inputs and carries a stale executable is rebuilt once after the release, and the
// landing then reports both effects and exits zero.
func TestLandRefreshesTheBrokerAfterPublication(t *testing.T) {
	t.Parallel()
	request := "land-refresh-after-publication"
	f := brokerDestinationFixture(t, request)
	j, calls := refreshJoins(func(root, executable string) error {
		return publishVerifyingBroker(t, root, executable)
	})

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, wantEffects("complete")) {
		t.Fatalf("landing with a stale broker = (%d, %q, %q), want a complete refresh at exit 0", r.exit, r.stdout, r.stderr)
	}
	if *calls != 1 {
		t.Fatalf("refresh build calls = %d, want exactly one", *calls)
	}
}

// TestLandReportsAFailedRefresh is LC4, LC5, and LC13. A build that returns without
// publishing a verified executable leaves the destination's broker stale. The landing
// reports the failed effect with the resume every incomplete step names, and it
// unpublishes nothing. A failed effect stops every later effect, so the cleanup reports
// pending and the folded sibling's checkout is untouched.
func TestLandReportsAFailedRefresh(t *testing.T) {
	t.Parallel()
	request := "land-refresh-failed"
	f := brokerDestinationFixture(t, request)
	folded := foldLandingSibling(t, f.root, f.home, request+"-sibling", f.creation)
	j, calls := refreshJoins(nil)

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, folded.tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, wantEffects("failed", "pending")) {
		t.Fatalf("failed refresh = (%d, %q, %q), want exit 3 with a failed refresh", r.exit, r.stdout, r.stderr)
	}
	requirePresent(t, folded.sibling.Path, "sibling worktree")
	if strings.Contains(r.stderr, "landing cleanup{") {
		t.Fatalf("failed refresh stderr = %q, want no cleanup plan row", r.stderr)
	}
	if !strings.Contains(r.stdout, "worktree=incomplete:refresh,next=bench worktree land --resume ") {
		t.Fatalf("failed refresh record = %q, want an incomplete:refresh cell with a resume", r.stdout)
	}
	if *calls != 1 {
		t.Fatalf("refresh build calls = %d, want exactly one", *calls)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	if published == f.base {
		t.Fatal("the failed refresh unpublished the landing")
	}
	if got := projectGreenMarker(t, f.root); got != published {
		t.Fatalf("project-green marker = %q, want the published commit %q", got, published)
	}
	if got := gitOutput(t, f.root, "rev-parse", "HEAD"); got != published {
		t.Fatalf("destination checkout = %q, want the published commit %q", got, published)
	}
}

// TestLandRefreshReadsTheSeal is LC8. A build that installs the executable and stops
// before its seal leaves a set that cannot authenticate. The refresh reads the seal, so
// the half-installed set reports failed rather than complete.
func TestLandRefreshReadsTheSeal(t *testing.T) {
	t.Parallel()
	request := "land-refresh-unsealed"
	f := brokerDestinationFixture(t, request)
	j, calls := refreshJoins(func(_, executable string) error {
		mustMkdirAll(t, filepath.Dir(executable), 0o755)
		mustWrite(t, executable, []byte("#!/bin/sh\nexit 0\n"), 0o755)
		return nil
	})

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, wantEffects("failed")) {
		t.Fatalf("unsealed refresh = (%d, %q, %q), want exit 3 with a failed refresh", r.exit, r.stdout, r.stderr)
	}
	if *calls != 1 {
		t.Fatalf("refresh build calls = %d, want exactly one", *calls)
	}
}

// openInstallBuild claims the outcome that the install fixture's binding names, and moves the
// binding to the default-branch identity of installSpec, so the bound assignment holds an
// open build.
func openInstallBuild(t *testing.T, root string) {
	t.Helper()
	identity, err := commitrepo.SourceIdentity(root, "main", installSpec)
	mustNoError(t, err)
	commitmenttest.EditCommitment(t, root, func(state *intent.CommitmentState) {
		state.Bindings[0].Identity = identity
		state.Claims = append(state.Claims, intent.OutcomeClaim{Milestone: state.Bindings[0].Milestone, Outcome: state.Bindings[0].Outcome})
	})
}

// openUndeclaredBuild claims a second outcome of the install fixture, and binds it to a new
// assignment whose default-branch spec declares no persisted field.
func openUndeclaredBuild(t *testing.T, f ownedAssignment) {
	t.Helper()
	const plain = "specs/plain/spec.md"
	mustMkdirAll(t, filepath.Join(f.root, filepath.Dir(plain)), 0o755)
	commitInWorktree(t, f.root, plain, "# plain\n\nStatus: staged\n", "an undeclared spec")
	identity, err := commitrepo.SourceIdentity(f.root, "main", plain)
	mustNoError(t, err)
	other := mustCreate(t, f.root, f.home, "undeclared-build", "undeclared-build")
	commitmenttest.EditCommitment(t, f.root, func(state *intent.CommitmentState) {
		state.Bindings = append(state.Bindings, intent.DeliveryBinding{Assignment: other.Assignment.ID, Request: other.Assignment.Request, Milestone: "fixture", Outcome: "plain", Deliverable: plain, Identity: identity})
		state.Claims = append(state.Claims, intent.OutcomeClaim{Milestone: "fixture", Outcome: "plain"})
	})
}

// TestRefreshHoldsAnInstalledCandidate is CB51 to CB53, CB62, and CB68. The primary
// executable has no seal, so, like an installed candidate, it never verifies against the
// primary sources. The refresh holds it only while a declared build is open and the primary
// manifest binds its digest, or while the ledger does not read. Other rows rebuild as before.
func TestRefreshHoldsAnInstalledCandidate(t *testing.T) {
	t.Parallel()
	undeclared := installShape{inputs: true, deliverable: installSpec}
	invalid := installShape{inputs: true, declaration: "Persisted fields: gate-evidence:x\n", deliverable: installSpec}
	for _, row := range []struct {
		name                             string
		shape                            installShape
		earlier, open, bound, unreadable bool
		want                             string
	}{
		{name: "a declared build with the manifest bound holds", shape: declaredInstall, open: true, bound: true, want: effectHeld},
		{name: "a declared build after an undeclared build holds", shape: declaredInstall, earlier: true, open: true, bound: true, want: effectHeld},
		{name: "a declared build with another manifest digest rebuilds", shape: declaredInstall, open: true, want: effectComplete},
		{name: "an undeclared open build rebuilds", shape: undeclared, open: true, bound: true, want: effectComplete},
		{name: "an invalid declaration counts as none", shape: invalid, open: true, bound: true, want: effectComplete},
		{name: "an unreadable ledger holds", shape: declaredInstall, open: true, unreadable: true, want: effectHeld},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := newInstallAssignment(t, "refresh-hold-"+strings.ReplaceAll(row.name, " ", "-"), row.shape)
			if row.earlier {
				openUndeclaredBuild(t, f)
			}
			if row.open {
				openInstallBuild(t, f.root)
			}
			if row.unreadable {
				plantCommitmentField(t, f.root)
			}
			primary, manifest := freshness.PublishedExecutable(f.root), manifestPath(f.root)
			mustMkdirAll(t, filepath.Dir(primary), 0o755)
			mustWrite(t, primary, []byte("installed-candidate"), 0o755)
			_, _, err := brokermanifest.Write(wrapperDir(f.root), primary, "1.2.3")
			mustNoError(t, err)
			if !row.bound {
				mustWrite(t, primary, []byte("rebuilt-candidate"), 0o755)
			}
			before := fixtureFileText(t, primary) + fixtureFileText(t, manifest)
			j, calls := refreshJoins(func(root, executable string) error {
				return publishVerifyingBroker(t, root, executable)
			})
			var stderr strings.Builder

			got := refreshBroker(j, nil, f.root, &stderr)
			requireTest(t, got == row.want, "refresh = %q (%q), want %q", got, stderr.String(), row.want)
			if row.want == effectHeld {
				after := fixtureFileText(t, primary) + fixtureFileText(t, manifest)
				requireTest(t, *calls == 0 && after == before, "held refresh: build calls %d, primary set changed %t; want no write", *calls, after != before)
				return
			}
			requireTest(t, *calls == 1 && freshness.Verify(f.root, primary) == nil, "rebuilt refresh: build calls %d, verify %v; want one verified rebuild", *calls, freshness.Verify(f.root, primary))
		})
	}
	// A held refresh is not a failure, so the landing exits 0 and the cleanup still runs. The
	// broker notice reads a held refresh as it reads a complete one.
	t.Run("a held landing exits 0", func(t *testing.T) {
		t.Parallel()
		request := "land-refresh-held"
		f := brokerChangingLanding(t, request)
		j, calls := refreshJoins(nil)
		j.heldInstall = func(commitrepo.Store) bool { return true }
		call := f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...)
		call.kit = f.root

		r := runVerb(t, verbLand, call)
		requireTest(t, r.exit == 0 && strings.Contains(r.stdout, wantEffects(effectHeld)), "held landing = (%d, %q, %q), want a held refresh at exit 0", r.exit, r.stdout, r.stderr)
		requireTest(t, *calls == 0 && !strings.Contains(r.stderr, brokerNoticeAnchor), "held landing: build calls %d, stderr %q; want no build and no broker notice", *calls, r.stderr)
	})
}

// TestLandEffectsRowPrecedesTheLandedRecord is LC20. Every existing reader takes the
// landed record from the last stdout line, so the effects table sits immediately before
// it rather than after it.
func TestLandEffectsRowPrecedesTheLandedRecord(t *testing.T) {
	t.Parallel()
	request := "land-effects-row-order"
	f := publicLandingFixture(t, request, "", "")
	j, _ := refreshJoins(nil)

	r := runVerb(t, verbLand, f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 0 {
		t.Fatalf("landing = (%d, %q, %q), want a released landing", r.exit, r.stdout, r.stderr)
	}
	lines := strings.Split(strings.TrimSuffix(r.stdout, "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("landing stdout = %q, want the effects table before the landed record", r.stdout)
	}
	tail := strings.Join(lines[len(lines)-4:], "\n") + "\n"
	if !strings.HasPrefix(tail, wantEffects("not-applicable")) {
		t.Fatalf("landing stdout tail = %q, want the effects table first", tail)
	}
	if !strings.HasPrefix(lines[len(lines)-1], "landed{") {
		t.Fatalf("last stdout line = %q, want the landed record", lines[len(lines)-1])
	}
}

// TestResumeReadsEffectStateFromTheTree is LC19. The resume keeps no journal of its own:
// it reads the destination's published executable and finishes the refresh an interrupted
// run left unfinished, then it leaves the finished refresh alone on the next call.
func TestResumeReadsEffectStateFromTheTree(t *testing.T) {
	t.Parallel()
	request := "land-refresh-resume-from-tree"
	f := brokerDestinationFixture(t, request)
	working, calls := refreshJoins(func(root, executable string) error {
		return publishVerifyingBroker(t, root, executable)
	})
	interrupted := working
	interrupted.releaseLandingAssignment = func(joins, ambient, string, []string, io.Writer, io.Writer) int { return 1 }

	r := runVerb(t, verbLand, f.callWith(interrupted, landArgs(request, f.base, f.tip, f.creation.Path)...))
	r.mustViaJoins(t)
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:release") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", r.exit, r.stdout, r.stderr)
	}
	if *calls != 0 {
		t.Fatalf("interrupted landing refresh calls = %d, want none", *calls)
	}
	published := gitOutput(t, f.root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", f.base, "--source-tip", f.tip, "--spec", "x", f.creation.Path}

	if r := runVerb(t, verbLandResume, f.callWith(working, args...)); r.exit != 0 || !strings.Contains(r.stdout, wantEffects("complete")) {
		t.Fatalf("first resume = (%d, %q, %q), want a completed refresh", r.exit, r.stdout, r.stderr)
	}
	if *calls != 1 {
		t.Fatalf("first resume refresh calls = %d, want exactly one", *calls)
	}

	if r := runVerb(t, verbLandResume, f.callWith(working, args...)); r.exit != 0 || !strings.Contains(r.stdout, wantEffects("complete")) {
		t.Fatalf("second resume = (%d, %q, %q), want the refresh left alone", r.exit, r.stdout, r.stderr)
	}
	if *calls != 1 {
		t.Fatalf("second resume refresh calls = %d, want no second build", *calls)
	}
	if got := gitOutput(t, f.root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main = %s, want %s", got, published)
	}
}
