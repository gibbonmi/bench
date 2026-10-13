package probe

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canonicalpath"
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

// probeWith runs the probe invocation that probeArgs answers for spelled over f, with each
// subject replacement through replacer. The mutated run catches the mutation, so a probe
// that reached it would bite.
func probeWith(t *testing.T, f *fixture, replacer durablefile.Replacer, spelled ...string) (string, int) {
	t.Helper()
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	parsed, line, _ := usage.Parse(grammar, probeArgs(spelled...))
	if line != "" {
		t.Fatalf("parse refused the accepted invocation: %q", line)
	}
	return runWith(replacer, f.root, parsed)
}

// D27: the mutation renames the mutated bytes in at the subject mode, and a directory-sync
// failure after that rename reaches the caller as the mutation failure at exit 1, after the
// baseline and before any focused run.
func TestProbeReplacement(t *testing.T) {
	f := newFixture(t)
	// The subject mode differs from the fixture default, so a mutation that drops it shows.
	if err := os.Chmod(f.subject, 0o640); err != nil {
		t.Fatal(err)
	}
	router := subjectFaults()
	cause := errors.New("injected mutation directory-sync failure")
	router.Fail[mutationReplacement] = map[string]error{"directory-sync": cause}
	replacer := router.Replacer()
	rename := replacer.Rename
	replacer.Rename = func(oldpath, newpath string) error {
		err := rename(oldpath, newpath)
		if router.Count(mutationReplacement.Destination) == mutationReplacement.Ordinal {
			durablefiletest.RequireFile(t, newpath, 0o640, strings.Replace(clampSource, "n < 0", "n > 0", 1))
		}
		return err
	}

	out, code := probeWith(t, f, replacer)

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

			out, code := probeWith(t, f, router.Replacer())

			if code != 1 {
				t.Fatalf("exit = %d, want 1\n%s", code, out)
			}
			requireBaselineOnly(t, f)
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
	// The first line is the mutation failure, so only the preserved block follows it.
	_, block, _ := strings.Cut(out, "\n")
	for _, want := range []string{kept, mutationCause.Error(), restoreCause.Error()} {
		if !strings.Contains(block, want) {
			t.Fatalf("preserved block = %q, want it to name %q", block, want)
		}
	}
	durablefiletest.RequireBytes(t, kept, clampSource)
}

// D41: a copy path the encoder refuses still reaches the caller with every cause, on the
// mutation-failure branch and on the focused-run branch alike.
func TestProbeNamesACopyTheEncoderRefuses(t *testing.T) {
	// The BEL byte sits in the base name, so the copy path carries it too.
	name := "data\aalert.txt"
	mutationCause := errors.New("injected mutation directory-sync failure")
	restoreCause := errors.New("injected restore rename failure")
	for _, tc := range []struct {
		label  string
		faults map[durablefiletest.Fault]map[string]error
	}{
		{"mutation failure", map[durablefiletest.Fault]map[string]error{
			mutationReplacement: {"directory-sync": mutationCause},
			restoreReplacement:  {"rename": restoreCause},
		}},
		{"focused run", map[durablefiletest.Fault]map[string]error{
			restoreReplacement: {"rename": restoreCause},
		}},
	} {
		t.Run(tc.label, func(t *testing.T) {
			f := newFixture(t)
			writeFixtureFile(t, filepath.Join(f.root, name), "alpha\n", 0o644)
			router := subjectFaults()
			router.Fail = tc.faults

			out, code := probeWith(t, f, router.Replacer(), name, "--swap", "alpha", "--with", "beta", "--package", "./")

			if code != 2 {
				t.Fatalf("exit = %d, want 2\n%s", code, out)
			}
			kept := preservedCopy(t, f)
			for _, cause := range tc.faults {
				for _, err := range cause {
					if !strings.Contains(out, err.Error()) {
						t.Fatalf("stdout = %q, want it to name %q", out, err.Error())
					}
				}
			}
			if !strings.Contains(out, strconv.Quote(kept)) {
				t.Fatalf("stdout = %q, want it to name the copy %s", out, strconv.Quote(kept))
			}
			durablefiletest.RequireBytes(t, kept, "alpha\n")
		})
	}
}

// D34: a restore whose directory sync fails is not proven, even though the subject already
// shows the start bytes, so the probe keeps the copy and exits 2.
func TestProbeFailedRestore(t *testing.T) {
	f := newFixture(t)
	router := subjectFaults()
	router.Fail[restoreReplacement] = map[string]error{"directory-sync": errors.New("injected restore directory-sync failure")}

	out, code := probeWith(t, f, router.Replacer())

	requireRow(t, out, code, 2, probeRow(t, "restore-failed", "clamp.go", "swap", "failed", 1, "no"))
	preservedCopy(t, f)
}

// D30: a subject outside the repository refuses before preserve or any replacement
// starts, and that subject keeps its bytes.
func TestProbeSubjectRefusalSkipsReplacement(t *testing.T) {
	f := newFixture(t)
	outside := filepath.Join(scratchDir(t), "clamp.go")
	writeFixtureFile(t, outside, clampSource, 0o644)
	_, display := canonicalpath.Operand(f.root, "", outside)
	rec := &durablefiletest.Recorder{}

	out, code := probeWith(t, f, rec.Replacer(), outside, "--swap", "n < 0", "--with", "n > 0", "--package", "./")

	requireRefused(t, f, out, code, unavailable(display, "outside the repository"), 1)
	durablefiletest.RequireNoReplacement(t, rec.Trace)
	durablefiletest.RequireBytes(t, outside, clampSource)
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
