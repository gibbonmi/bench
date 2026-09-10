// Promotion-broker tests for the landing command: the notice a broker-changing diff
// prints before publication, and the refresh effect that republishes the destination's
// own broker after publication.
package worktree

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

// brokerChangingLanding is the landing fixture whose reviewed diff changes the promotion
// broker's own build inputs. Both install-step rows read the same notice, so they compose
// the same destination rather than each building one.
func brokerChangingLanding(t *testing.T, request string) (root string, creation Creation, base, tip, home string) {
	t.Helper()
	root, creation, _, _, _, home = publicLandingFixture(t, request, "", "")
	writeGoMainFixture(t, root)
	mustWrite(t, filepath.Join(root, "scripts", "go-build.inputs"), []byte("build_script=scripts/go-build.sh\n"), 0o644)
	spec := filepath.Join(root, "specs", "x", "spec.md")
	body, err := os.ReadFile(spec)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, spec, append(body, []byte("- `scripts/go-build.sh`\n")...), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "broker build inputs")
	gitRun(t, creation.Path, "rebase", "main")
	base = gitOutput(t, root, "rev-parse", "HEAD")
	commitInWorktree(t, creation.Path, "scripts/go-build.sh", "#!/bin/sh\n# next broker\nexit 0\n", "change broker source")
	return root, creation, base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), home
}

// writeGoMainFixture writes the resolvable Go main package and build script a broker
// fixture's destination stands on. The seal's source digest is the build-input closure of
// that package, so a destination without it reports an unresolvable tree rather than the
// broker state the row is about.
func writeGoMainFixture(t *testing.T, root string) {
	t.Helper()
	mustWrite(t, filepath.Join(root, "go.mod"), []byte("module benchfixture\n\ngo 1.22\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "cmd", "bench"), 0o755)
	mustWrite(t, filepath.Join(root, "cmd", "bench", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644)
	mustMkdirAll(t, filepath.Join(root, "scripts"), 0o755)
	mustWrite(t, filepath.Join(root, "scripts", "go-build.sh"), []byte("#!/bin/sh\nexit 0\n"), 0o755)
}

// kitCheckoutJoins is the landing seam set whose checkout predicate answers a fixed
// verdict. The real predicate reads the kit root from the process environment, so a
// fixture that bound that environment would leave the package's parallel set.
func kitCheckoutJoins(kit bool) joins {
	j := defaultJoins()
	j.kitSourceCheckout = func(string) bool { return kit }
	return j
}

// TestLandCommandReportsInstallStepForABrokerChangingDiff is SOL16 and BF19. A reviewed
// diff that changes the promotion broker's own build inputs lands as source, but the
// installed broker keeps authority. The landing must name the install step so the
// operator does not expect source publication to replace it. In the kit source checkout
// that step is the stamped rebuild and bench doctor --fix, because bench repair reads a
// pin manifest the source tree does not carry.
func TestLandCommandReportsInstallStepForABrokerChangingDiff(t *testing.T) {
	t.Parallel()
	request := "land-owner-broker-change"
	root, creation, base, tip, home := brokerChangingLanding(t, request)

	var stdout, stderr bytes.Buffer
	code := landWith(kitCheckoutJoins(true), root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:refresh") {
		t.Fatalf("broker-changing landing = (%d, %q, %q), want a failed refresh", code, stdout.String(), stderr.String())
	}
	notice := brokerNoticeLine(t, stderr.String())
	if !strings.Contains(notice, freshness.RebuildAction(root)) {
		t.Fatalf("kit-checkout landing named no rebuild: %q", notice)
	}
	if !strings.Contains(notice, "bench doctor --fix") {
		t.Fatalf("kit-checkout landing named no publication step: %q", notice)
	}
}

// TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout is BF20. An installed kit
// carries the pin manifest bench repair reads, so the notice keeps that route there and
// names no source-tree rebuild.
func TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout(t *testing.T) {
	t.Parallel()
	request := "land-owner-broker-change-installed"
	root, creation, base, tip, home := brokerChangingLanding(t, request)

	var stdout, stderr bytes.Buffer
	code := landWith(kitCheckoutJoins(false), root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:refresh") {
		t.Fatalf("broker-changing landing = (%d, %q, %q), want a failed refresh", code, stdout.String(), stderr.String())
	}
	notice := brokerNoticeLine(t, stderr.String())
	if !strings.Contains(notice, "bench repair") {
		t.Fatalf("installed-kit landing named no install step: %q", notice)
	}
	if strings.Contains(notice, freshness.RebuildAction(root)) {
		t.Fatalf("installed-kit landing named the source-tree rebuild: %q", notice)
	}
}

// brokerNoticeLine is the install-step notice inside a landing's diagnostics. The refresh
// effect's own failure names the rebuild command too, so a row about which route the
// notice names reads that one line rather than the whole stream.
func brokerNoticeLine(t *testing.T, diagnostics string) string {
	t.Helper()
	for _, line := range strings.Split(diagnostics, "\n") {
		if strings.Contains(line, "the installed broker keeps authority") {
			return line
		}
	}
	t.Fatalf("landing printed no broker install notice: %q", diagnostics)
	return ""
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
func brokerDestinationFixture(t *testing.T, request string) (root string, creation Creation, base, tip, home string) {
	t.Helper()
	root, creation, _, _, _, home = publicLandingFixture(t, request, "dist/bench", "dist/")
	if err := os.Remove(filepath.Join(creation.Path, "dist", "bench")); err != nil {
		t.Fatal(err)
	}
	writeGoMainFixture(t, root)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "broker sources")
	commitLandingBuildInputs(t, root, "build_script=scripts/go-build.sh\n")
	base = gitOutput(t, root, "rev-parse", "HEAD")
	gitRun(t, creation.Path, "rebase", "main")
	return root, creation, base, gitOutput(t, creation.Path, "rev-parse", "HEAD"), home
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
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	j, calls := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("not-applicable")) {
		t.Fatalf("landing without build inputs = (%d, %q, %q), want a not-applicable refresh", code, stdout.String(), stderr.String())
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
	root, creation, base, tip, home := brokerDestinationFixture(t, request)
	if err := publishVerifyingBroker(t, root, freshness.PublishedExecutable(root)); err != nil {
		t.Fatal(err)
	}
	j, calls := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("complete")) {
		t.Fatalf("landing with a fresh broker = (%d, %q, %q), want a complete refresh", code, stdout.String(), stderr.String())
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
	root, creation, base, tip, home := brokerDestinationFixture(t, request)
	j, calls := refreshJoins(func(root, executable string) error {
		return publishVerifyingBroker(t, root, executable)
	})

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), wantEffects("complete")) {
		t.Fatalf("landing with a stale broker = (%d, %q, %q), want a complete refresh at exit 0", code, stdout.String(), stderr.String())
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
	root, creation, base, _, home := brokerDestinationFixture(t, request)
	sibling, tip := foldLandingSibling(t, root, home, request+"-sibling", creation)
	j, calls := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), wantEffects("failed", "pending")) {
		t.Fatalf("failed refresh = (%d, %q, %q), want exit 3 with a failed refresh", code, stdout.String(), stderr.String())
	}
	requirePresent(t, sibling.Path, "sibling worktree")
	if strings.Contains(stderr.String(), "landing cleanup{") {
		t.Fatalf("failed refresh stderr = %q, want no cleanup plan row", stderr.String())
	}
	if !strings.Contains(stdout.String(), "worktree=incomplete:refresh,next=bench worktree land --resume ") {
		t.Fatalf("failed refresh record = %q, want an incomplete:refresh cell with a resume", stdout.String())
	}
	if *calls != 1 {
		t.Fatalf("refresh build calls = %d, want exactly one", *calls)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	if published == base {
		t.Fatal("the failed refresh unpublished the landing")
	}
	if got := projectGreenMarker(t, root); got != published {
		t.Fatalf("project-green marker = %q, want the published commit %q", got, published)
	}
	if got := gitOutput(t, root, "rev-parse", "HEAD"); got != published {
		t.Fatalf("destination checkout = %q, want the published commit %q", got, published)
	}
}

// TestLandRefreshReadsTheSeal is LC8. A build that installs the executable and stops
// before its seal leaves a set that cannot authenticate. The refresh reads the seal, so
// the half-installed set reports failed rather than complete.
func TestLandRefreshReadsTheSeal(t *testing.T) {
	t.Parallel()
	request := "land-refresh-unsealed"
	root, creation, base, tip, home := brokerDestinationFixture(t, request)
	j, calls := refreshJoins(func(_, executable string) error {
		mustMkdirAll(t, filepath.Dir(executable), 0o755)
		mustWrite(t, executable, []byte("#!/bin/sh\nexit 0\n"), 0o755)
		return nil
	})

	var stdout, stderr bytes.Buffer
	code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), wantEffects("failed")) {
		t.Fatalf("unsealed refresh = (%d, %q, %q), want exit 3 with a failed refresh", code, stdout.String(), stderr.String())
	}
	if *calls != 1 {
		t.Fatalf("refresh build calls = %d, want exactly one", *calls)
	}
}

// TestLandEffectsRowPrecedesTheLandedRecord is LC20. Every existing reader takes the
// landed record from the last stdout line, so the effects table sits immediately before
// it rather than after it.
func TestLandEffectsRowPrecedesTheLandedRecord(t *testing.T) {
	t.Parallel()
	request := "land-effects-row-order"
	root, creation, base, tip, _, home := publicLandingFixture(t, request, "", "")
	j, _ := refreshJoins(nil)

	var stdout, stderr bytes.Buffer
	if code := landWith(j, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 0 {
		t.Fatalf("landing = (%d, %q, %q), want a released landing", code, stdout.String(), stderr.String())
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) < 4 {
		t.Fatalf("landing stdout = %q, want the effects table before the landed record", stdout.String())
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
	root, creation, base, tip, home := brokerDestinationFixture(t, request)
	working, calls := refreshJoins(func(root, executable string) error {
		return publishVerifyingBroker(t, root, executable)
	})
	interrupted := working
	interrupted.releaseLandingAssignment = func(joins, string, string, []string, io.Writer, io.Writer) int { return 1 }

	var stdout, stderr bytes.Buffer
	if code := landWith(interrupted, root, home, "", landArgs(request, base, tip, creation.Path), &stdout, &stderr); code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:release") {
		t.Fatalf("interrupted landing = (%d, %q, %q)", code, stdout.String(), stderr.String())
	}
	if *calls != 0 {
		t.Fatalf("interrupted landing refresh calls = %d, want none", *calls)
	}
	published := gitOutput(t, root, "rev-parse", "main")
	args := []string{"--resume", published, "--request", request, "--base", base, "--source-tip", tip, "--spec", "x", creation.Path}

	stdout.Reset()
	stderr.Reset()
	if code := resumeLandWith(working, root, home, args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), wantEffects("complete")) {
		t.Fatalf("first resume = (%d, %q, %q), want a completed refresh", code, stdout.String(), stderr.String())
	}
	if *calls != 1 {
		t.Fatalf("first resume refresh calls = %d, want exactly one", *calls)
	}

	stdout.Reset()
	stderr.Reset()
	if code := resumeLandWith(working, root, home, args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), wantEffects("complete")) {
		t.Fatalf("second resume = (%d, %q, %q), want the refresh left alone", code, stdout.String(), stderr.String())
	}
	if *calls != 1 {
		t.Fatalf("second resume refresh calls = %d, want no second build", *calls)
	}
	if got := gitOutput(t, root, "rev-parse", "main"); got != published {
		t.Fatalf("resume republished: main = %s, want %s", got, published)
	}
}
