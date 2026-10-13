package probe

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/probe/probetest"
)

// lockUnavailableLine opens every refusal for a subject lock that cannot be set up.
const lockUnavailableLine = "error: probe subject lock unavailable — "

// bitRow is the verdict row of the accepted probe over a mutated run that bites.
func bitRow(t *testing.T) string {
	t.Helper()
	return probeRow(t, "bit", "clamp.go", "swap", "failed", 1, "TestClampNegative", "yes")
}

// requireLockRefusal asserts a refusal that answers before the baseline: the exit, the
// subject unchanged, and no run child started.
func requireLockRefusal(t *testing.T, f *fixture, out string, code int, ok bool) {
	t.Helper()
	if !ok || code != 1 {
		t.Fatalf("refusal = (%q, %d), want the lock refusal and 1", out, code)
	}
	requireSubjectBytes(t, f, clampSource)
	requireNoRunChild(t, f)
}

// PS58, PS60: a second probe on a held subject refuses before any write or run child, so two
// probes never mutate one file.
func TestSubjectLockRefusesASecondProbe(t *testing.T) {
	f := refusalFixture(t)
	release := startSubjectHolder(t, f.subject)
	defer release()
	before := listRecordSurfaces(t, f)
	out, code := runProbe(t, probeArgs()...)
	requireLockRefusal(t, f, out, code, out == "error: probe subject in use — another probe holds clamp.go\n")
	if after := listRecordSurfaces(t, f); after != before {
		t.Fatalf("record surfaces changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// PS61: the gate lock is read before the subject lock, so a live gate names itself even when
// another probe also holds the subject.
func TestSubjectLockFollowsTheGateLock(t *testing.T) {
	f := refusalFixture(t)
	releaseGate := startLockHolder(t, gateLockPath(f.root))
	defer releaseGate()
	releaseSubject := startSubjectHolder(t, f.subject)
	defer releaseSubject()
	out, code := runProbe(t, probeArgs()...)
	requireLockRefusal(t, f, out, code, out == gateRefusal)
}

// PS59, PS75: the subject is read again under the lock, so a change inside the start-read
// window refuses rather than preserving bytes that another probe wrote.
func TestSubjectLockRereadsTheSubject(t *testing.T) {
	changed := "error: probe subject changed — clamp.go changed before the subject lock; run the probe again\n"
	t.Run("bytes", func(t *testing.T) {
		f := refusalFixture(t)
		writeFixtureFile(t, f.subject, strings.Replace(clampSource, "keeps n at", "holds n at", 1), 0o644)
		setHook(t, &startReadHook, func() { writeFixtureFile(t, f.subject, clampSource, 0o644) })
		out, code := runProbe(t, probeArgs()...)
		requireRefused(t, f, out, code, changed, 1)
	})
	t.Run("mode", func(t *testing.T) {
		f := refusalFixture(t)
		setHook(t, &startReadHook, func() {
			if err := os.Chmod(f.subject, 0o755); err != nil {
				t.Fatal(err)
			}
		})
		out, code := runProbe(t, probeArgs()...)
		requireRefused(t, f, out, code, changed, 1)
	})
	t.Run("deleted", func(t *testing.T) {
		f := refusalFixture(t)
		setHook(t, &startReadHook, func() {
			if err := os.Remove(f.subject); err != nil {
				t.Fatal(err)
			}
		})
		out, code := runProbe(t, probeArgs()...)
		if want := "error: probe subject unavailable — clamp.go is absent\n"; out != want || code != 1 {
			t.Fatalf("refusal = (%q, %d), want (%q, 1)", out, code, want)
		}
		requireHomeEmpty(t, f)
		requireNoRunChild(t, f)
	})
}

// PS62: the lock names one subject, so a probe on another file runs beside a held lock.
func TestSubjectLockAllowsAnotherSubject(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	writeFixtureFile(t, filepath.Join(f.root, "other.txt"), "alpha\n", 0o644)
	release := startSubjectHolder(t, f.subject)
	defer release()
	out, code := runProbe(t, "other.txt", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampNegative$")
	requireRow(t, out, code, 0, probeRow(t, "bit", "other.txt", "swap", "failed", 1, "TestClampNegative", "yes"))
}

// PS63: an interrupted probe releases its lock, so the next probe on the subject reaches its
// verdict and the interrupted one leaves no lock entry behind.
func TestSubjectLockReleasesOnInterrupt(t *testing.T) {
	f := newFixture(t)
	installStubStarts(t, f, probetest.Start{Events: cannedPass}, probetest.Start{Hook: "sleep 60\n", Events: cannedPass})
	if out, want := interruptedProbe(t, f, 2), probeRow(t, "invalid", "clamp.go", "swap", "interrupted", 0, nil, "yes"); !strings.Contains(out, want) {
		t.Fatalf("helper stdout = %q, want the interrupted row %q", out, want)
	}
	requireHomeEmpty(t, f)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 0, bitRow(t))
}

// PS64: a lock file whose holder died grants the lock, because a record lock dies with its
// process, and the probe removes that file at its end.
func TestSubjectLockIgnoresAFileWithoutAHolder(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	writeFixtureFile(t, subjectLockPath(f.subject), "left by a killed probe\n", 0o600)
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 0, bitRow(t))
	requireHomeEmpty(t, f)
}

// PS65, PS66: the held descriptor must name the file at the lock path, so a replaced lock file
// gives a retry, and a replacement on every attempt refuses rather than retrying forever.
func TestSubjectLockRechecksTheIdentity(t *testing.T) {
	t.Run("replaced once", func(t *testing.T) {
		f := newFixture(t)
		installStubGo(t, f, cannedFailure("caught"), 1, "")
		calls := 0
		setHook(t, &lockOwnerHook, func(path string) {
			if calls++; calls == 1 {
				replaceLockFile(t, path)
			}
		})
		out, code := runProbe(t, probeArgs()...)
		requireRow(t, out, code, 0, bitRow(t))
		if calls != 2 {
			t.Fatalf("lock-owner hook calls = %d, want 2: one replaced file and one retry", calls)
		}
		requireHomeEmpty(t, f)
	})
	t.Run("replaced on every attempt", func(t *testing.T) {
		f := refusalFixture(t)
		calls := 0
		setHook(t, &lockOwnerHook, func(path string) {
			calls++
			replaceLockFile(t, path)
		})
		out, code := runProbe(t, probeArgs()...)
		requireLockRefusal(t, f, out, code, strings.HasPrefix(out, lockUnavailableLine) && strings.Count(out, "\n") == 1)
		if calls != lockAttempts {
			t.Fatalf("lock-owner hook calls = %d, want %d", calls, lockAttempts)
		}
	})
}

// PS67: a lock directory path that is not a directory refuses before the baseline, so a
// setup fault never starts an unlocked run.
func TestSubjectLockRefusesAnUnusableDirectory(t *testing.T) {
	f := refusalFixture(t)
	writeFixtureFile(t, filepath.Dir(subjectLockPath(f.subject)), "not a directory\n", 0o644)
	out, code := runProbe(t, probeArgs()...)
	requireLockRefusal(t, f, out, code, strings.HasPrefix(out, lockUnavailableLine) && strings.Count(out, "\n") == 1)
}

// Only the contention rule names another probe, so a lock error that no holder causes refuses
// as a setup fault before the baseline.
func TestSubjectLockReadsOnlyContentionAsInUse(t *testing.T) {
	f := refusalFixture(t)
	setHook(t, &takeRecordLock, func(*os.File) error { return syscall.ENOLCK })
	out, code := runProbe(t, probeArgs()...)
	requireLockRefusal(t, f, out, code, strings.HasPrefix(out, lockUnavailableLine) && strings.HasSuffix(out, syscall.ENOLCK.Error()+"\n"))
}

// An open that answers ENOENT starts the next attempt, because a holder's release removes the
// lock directory between this attempt's create and its open. A third ENOENT refuses.
func TestSubjectLockRetriesAnAbsentLockDirectory(t *testing.T) {
	t.Run("removed once", func(t *testing.T) {
		f := newFixture(t)
		installStubGo(t, f, cannedFailure("caught"), 1, "")
		calls := 0
		setHook(t, &lockOpenHook, func(path string) {
			if calls++; calls == 1 {
				removeLockDir(t, path)
			}
		})
		out, code := runProbe(t, probeArgs()...)
		requireRow(t, out, code, 0, bitRow(t))
		if calls != 2 {
			t.Fatalf("lock-open hook calls = %d, want 2: one absent directory and one retry", calls)
		}
		requireHomeEmpty(t, f)
	})
	t.Run("removed on every attempt", func(t *testing.T) {
		f := refusalFixture(t)
		calls := 0
		setHook(t, &lockOpenHook, func(path string) {
			calls++
			removeLockDir(t, path)
		})
		out, code := runProbe(t, probeArgs()...)
		requireLockRefusal(t, f, out, code, strings.HasPrefix(out, lockUnavailableLine) && strings.HasSuffix(out, syscall.ENOENT.Error()+"\n"))
		if calls != lockAttempts {
			t.Fatalf("lock-open hook calls = %d, want %d", calls, lockAttempts)
		}
	})
}

// The open never follows a link at the lock path, so a planted link refuses before the
// baseline and creates nothing at the link's target.
func TestSubjectLockRefusesALinkAtTheLockPath(t *testing.T) {
	f := refusalFixture(t)
	path := subjectLockPath(f.subject)
	target := filepath.Join(scratchDir(t), "target")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
	out, code := runProbe(t, probeArgs()...)
	requireLockRefusal(t, f, out, code, strings.HasPrefix(out, lockUnavailableLine) && strings.HasSuffix(out, syscall.ELOOP.Error()+"\n"))
	if _, err := os.Lstat(target); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("the link target exists after the refusal: %v", err)
	}
}

// The lock directory and the lock file belong to their owner alone, so no other user can open
// a subject's lock file and hold its lock.
func TestSubjectLockKeepsItsFilesPrivate(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	var dirMode, fileMode fs.FileMode
	setHook(t, &lockOwnerHook, func(path string) {
		dirMode, fileMode = permissions(t, filepath.Dir(path)), permissions(t, path)
	})
	out, code := runProbe(t, probeArgs()...)
	requireRow(t, out, code, 0, bitRow(t))
	if dirMode != 0o700 || fileMode != 0o600 {
		t.Fatalf("lock modes = (%o, %o), want (700, 600)", dirMode, fileMode)
	}
}
