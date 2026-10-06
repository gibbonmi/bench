// Broker-notice tests for the landing command: the install step a broker-changing diff
// names after the refresh effect reports its result.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
)

// brokerChangingLanding is the landing fixture whose reviewed diff changes the promotion
// broker's own build inputs. Both install-step rows read the same notice, so they compose
// the same destination rather than each building one.
func brokerChangingLanding(t *testing.T, request string) landingFixture {
	t.Helper()
	f := publicLandingFixture(t, request, "", "")
	writeGoMainFixture(t, f.root)
	mustWrite(t, filepath.Join(f.root, filepath.FromSlash(freshness.BuildInputsManifest)), []byte(freshness.BuildInputLine("build_script", "scripts/go-build.sh")), 0o644)
	spec := filepath.Join(f.root, "specs", "x", "spec.md")
	body, err := os.ReadFile(spec)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, spec, withFenceEntry(body, "scripts/go-build.sh"), 0o644)
	// The fence edit changes the approved deliverable, so the same commit re-approves it.
	commitmenttest.SeedAdmission(t, f.root, "specs/x/spec.md")
	gitRun(t, f.root, "add", ".")
	gitRun(t, f.root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "broker build inputs")
	gitRun(t, f.creation.Path, "rebase", "main")
	commitmenttest.Rebind(t, f.creation.Path, request, "specs/x/spec.md")
	refreshLandingEvidence(t, f.creation.Path, gitOutput(t, f.root, "rev-parse", "HEAD"))
	f.base = gitOutput(t, f.root, "rev-parse", "HEAD")
	commitInWorktree(t, f.creation.Path, "scripts/go-build.sh", "#!/bin/sh\n# next broker\nexit 0\n", "change broker source")
	refreshLandingEvidence(t, f.creation.Path, f.base)
	f.tip = gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	return f
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
	f := brokerChangingLanding(t, request)
	call := f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)
	call.kit = f.root

	r := runVerb(t, verbLand, call)
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:refresh") {
		t.Fatalf("broker-changing landing = (%d, %q, %q), want a failed refresh", r.exit, r.stdout, r.stderr)
	}
	notice := brokerNoticeLine(t, r.stderr)
	if !strings.Contains(notice, freshness.RebuildAction(f.root)) {
		t.Fatalf("kit-checkout landing named no rebuild: %q", notice)
	}
	if !strings.Contains(notice, "bench doctor --fix") {
		t.Fatalf("kit-checkout landing named no publication step: %q", notice)
	}
	if failed := strings.Index(r.stderr, "landing refresh failed: "); failed < 0 || failed > strings.Index(r.stderr, notice) {
		t.Fatalf("broker notice printed before the refresh effect ran: %q", r.stderr)
	}
}

// TestLandBrokerNoticeAfterACompleteRefresh reads the broker notice after a complete
// refresh. In the kit checkout the refresh has already republished the broker, its seal,
// and its manifest, so a notice that named the rebuild would report a step the landing no
// longer owes. Off the kit checkout the refresh republishes the destination's own
// executable, not the installed broker, so the notice still names the install route.
func TestLandBrokerNoticeAfterACompleteRefresh(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name       string
		kitAtRoot  bool
		wantNotice bool
	}{
		{name: "kit-checkout", kitAtRoot: true, wantNotice: false},
		{name: "installed", kitAtRoot: false, wantNotice: true},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			request := "land-owner-broker-change-refreshed-" + row.name
			f := brokerChangingLanding(t, request)
			j, _ := refreshJoins(func(root, executable string) error {
				return publishVerifyingBroker(t, root, executable)
			})
			call := f.callWith(j, landArgs(request, f.base, f.tip, f.creation.Path)...)
			call.kit = t.TempDir()
			if row.kitAtRoot {
				call.kit = f.root
			}

			r := runVerb(t, verbLand, call)
			if r.exit != 0 || !strings.Contains(r.stdout, wantEffects("complete")) {
				t.Fatalf("refreshed broker-changing landing = (%d, %q, %q), want a complete refresh at exit 0", r.exit, r.stdout, r.stderr)
			}
			if !row.wantNotice {
				if strings.Contains(r.stderr, brokerNoticeAnchor) {
					t.Fatalf("refreshed kit-checkout landing named a manual rebuild: %q", r.stderr)
				}
				return
			}
			if notice := brokerNoticeLine(t, r.stderr); !strings.Contains(notice, "bench repair") {
				t.Fatalf("installed-kit landing named no install step: %q", notice)
			}
		})
	}
}

// TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout is BF20. An installed kit
// carries the pin manifest bench repair reads, so the notice keeps that route there and
// names no source-tree rebuild.
func TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout(t *testing.T) {
	t.Parallel()
	request := "land-owner-broker-change-installed"
	f := brokerChangingLanding(t, request)
	call := f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)
	call.kit = t.TempDir()

	r := runVerb(t, verbLand, call)
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:refresh") {
		t.Fatalf("broker-changing landing = (%d, %q, %q), want a failed refresh", r.exit, r.stdout, r.stderr)
	}
	notice := brokerNoticeLine(t, r.stderr)
	if !strings.Contains(notice, "bench repair") {
		t.Fatalf("installed-kit landing named no install step: %q", notice)
	}
	if strings.Contains(notice, freshness.RebuildAction(f.root)) {
		t.Fatalf("installed-kit landing named the source-tree rebuild: %q", notice)
	}
}

// brokerNoticeAnchor is the phrase that marks the broker install notice in a landing's
// diagnostics.
const brokerNoticeAnchor = "the installed broker keeps authority"

// brokerNoticeLine is the install-step notice inside a landing's diagnostics. The refresh
// effect's own failure names the rebuild command too, so a row about which route the
// notice names reads that one line rather than the whole stream.
func brokerNoticeLine(t *testing.T, diagnostics string) string {
	t.Helper()
	for _, line := range strings.Split(diagnostics, "\n") {
		if strings.Contains(line, brokerNoticeAnchor) {
			return line
		}
	}
	t.Fatalf("landing printed no broker install notice: %q", diagnostics)
	return ""
}
