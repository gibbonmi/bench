package probe

import (
	"errors"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// The probe replaces only its subject: the mutation first, then the restore.
var (
	mutationReplacement = durablefiletest.Fault{Destination: "subject", Ordinal: 1}
	restoreReplacement  = durablefiletest.Fault{Destination: "subject", Ordinal: 2}
)

// subjectFaults routes every replacement of one probe to the subject destination.
func subjectFaults() *durablefiletest.Router {
	return &durablefiletest.Router{
		Destination: func(string) string { return mutationReplacement.Destination },
		Fail:        map[durablefiletest.Fault]map[string]error{},
	}
}

// probeWith runs the accepted probe invocation over f with each subject replacement through
// replacer. The mutated run catches the mutation, so a probe that reached it would bite.
func probeWith(t *testing.T, f *fixture, replacer durablefile.Replacer) (string, int) {
	t.Helper()
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	parsed, line, _ := usage.Parse(grammar, probeArgs())
	if line != "" {
		t.Fatalf("parse refused the accepted invocation: %q", line)
	}
	return runWith(replacer, f.root, parsed)
}

// D27: a directory-sync failure after the mutation's rename reaches the caller as the
// mutation failure at exit 1, after the baseline and before any focused run.
func TestProbeReplacement(t *testing.T) {
	f := newFixture(t)
	router := subjectFaults()
	cause := errors.New("injected mutation directory-sync failure")
	router.Fail[mutationReplacement] = map[string]error{"directory-sync": cause}

	out, code := probeWith(t, f, router.Replacer())

	want := toon.Errorf("probe mutation failed", router.Failure(mutationReplacement.Destination, "directory-sync", true, cause)) + "\n"
	if out != want || code != 1 {
		t.Fatalf("probe = (%q, %d), want (%q, 1)", out, code, want)
	}
	requireBaselineOnly(t, f)
}

// D42: each failure after the mutation's rename may leave the mutated bytes visible, so the
// probe restores the start bytes before it releases the copy.
func TestProbeMutationRecoverySuccess(t *testing.T) {
	for _, stage := range durablefiletest.PostRenameStages {
		t.Run(stage, func(t *testing.T) {
			f := newFixture(t)
			router := subjectFaults()
			router.Fail[mutationReplacement] = map[string]error{stage: errors.New("injected mutation " + stage + " failure")}

			probeWith(t, f, router.Replacer())

			durablefiletest.RequireBytes(t, f.subject, clampSource)
			requireHomeEmpty(t, f)
		})
	}
}

// D41: a failed restore after a failed mutation keeps the copy, names it, and exits 2 with
// both causes.
func TestProbeMutationRecoveryFailure(t *testing.T) {
	f := newFixture(t)
	router := subjectFaults()
	mutationCause := errors.New("injected mutation directory-sync failure")
	restoreCause := errors.New("injected restore rename failure")
	router.Fail[mutationReplacement] = map[string]error{"directory-sync": mutationCause}
	router.Fail[restoreReplacement] = map[string]error{"rename": restoreCause}

	out, code := probeWith(t, f, router.Replacer())

	if code != 2 {
		t.Fatalf("exit = %d, want 2\n%s", code, out)
	}
	kept := preservedCopy(t, f)
	for _, want := range []string{kept, mutationCause.Error(), restoreCause.Error()} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout = %q, want it to name %q", out, want)
		}
	}
	durablefiletest.RequireBytes(t, kept, clampSource)
}

// D34: a restore whose directory sync fails is not proven, even though the subject already
// shows the start bytes, so the probe keeps the copy and exits 2.
func TestProbeFailedRestore(t *testing.T) {
	f := newFixture(t)
	router := subjectFaults()
	router.Fail[restoreReplacement] = map[string]error{"directory-sync": errors.New("injected restore directory-sync failure")}

	out, code := probeWith(t, f, router.Replacer())

	requireRow(t, out, code, 2, probeRow(t, "restore-failed", "clamp.go", "swap", "failed", 1, "no"))
	durablefiletest.RequireBytes(t, preservedCopy(t, f), clampSource)
}

// D43: a mutation that fails before its rename leaves the start bytes, so the probe
// releases the copy without a restore.
func TestProbeMutationPreparationFailure(t *testing.T) {
	f := newFixture(t)
	rec := &durablefiletest.Recorder{Fail: map[string]error{"rename": errors.New("injected mutation rename failure")}}

	probeWith(t, f, rec.Replacer())

	durablefiletest.RequireBytes(t, f.subject, clampSource)
	requireHomeEmpty(t, f)
}
