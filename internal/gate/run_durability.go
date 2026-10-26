package gate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// persistRecovering installs rec as the verdict cache of gitdir. When that fails, it
// attempts once to install pending in its place. The error keeps the original cause, and
// a failed recovery adds its own cause beside it.
func persistRecovering(replacer durablefile.Replacer, gitdir string, rec, pending verdictRecord) error {
	err := durableReplace(replacer, gitdir, rec)
	if err == nil {
		return nil
	}
	if recovery := durableReplace(replacer, gitdir, pending); recovery != nil {
		return fmt.Errorf("%w; pending recovery failed: %w", err, recovery)
	}
	return err
}

// persistenceFailure reports a failed persistence of category with its cause as an
// operational result.
func persistenceFailure(root string, gateExit int, stderr io.Writer, category string, err error) Result {
	return operational(root, gateExit, stderr, category+": "+err.Error())
}

// unprotectedFailure keeps the pending record after a failed run whose protection is not
// safe, so no verdict claims the tree. It reports the process exit beside the protection cause.
func unprotectedFailure(root string, plan subject, rc int, stderr io.Writer, run processGroupResult) Result {
	return retainedFailure(root, plan, rc, stderr, fmt.Sprintf("gate terminal protection failed: process exit %d; %v", run.Outcome.Process.Code, run.StartErr))
}

// retainedFailure keeps the pending record, so no verdict claims the tree, and invalidates
// retained evidence. It reports message with any invalidation failure beside it.
func retainedFailure(root string, plan subject, rc int, stderr io.Writer, message string) Result {
	if err := invalidateEvidence(root, plan); err != nil {
		message += "; gate evidence invalidation failed: " + err.Error()
	}
	return operational(root, rc, stderr, message)
}

// persistInterruptedIfGreen demotes a reusable green verdict of root to a pending record
// so that a run which cannot take the execution lock leaves no reusable green behind.
func persistInterruptedIfGreen(replacer durablefile.Replacer, root, gitdir string, plan subject) error {
	if !inspectAt(root, time.Now().UTC()).ReusableGreen {
		return nil
	}
	pending := interruptedRecord(plan, time.Now().UTC())
	return persistRecovering(replacer, gitdir, pending, pending)
}

// besideDemotion returns message with a failed interrupted-green demotion beside it.
func besideDemotion(message string, err error) string {
	if err == nil {
		return message
	}
	return message + "; interrupted-green demotion failed: " + err.Error()
}

// reuseWritesNothing answers whether a reuse before the gate lock answers with no write. The
// gate cache can already describe this subject's fresh green. Another run can hold the gate
// lock, and a reuse does not wait for it.
func reuseWritesNothing(root string, plan subject, now time.Time) bool {
	if inspectSubjectAt(root, plan, now).ReusableGreen {
		return true
	}
	gitdir, err := benchgit.AdminDir(root)
	if err != nil {
		return false
	}
	held, err := lockHeld(gitdir)
	return err == nil && held
}

// recordReusedGreen installs reuse as the gate cache of gitdir, unless that cache already
// describes this subject's fresh green. The record keeps the time of the store record, so a
// reuse does not move the freshness window.
func recordReusedGreen(replacer durablefile.Replacer, root, gitdir string, plan subject, reuse Inspection, now time.Time) error {
	if inspectSubjectAt(root, plan, now).ReusableGreen {
		return nil
	}
	return durableReplace(replacer, gitdir, readyVerdict(plan, "green", reuse.RecordedAt))
}

func interruptedRecord(plan subject, now time.Time) verdictRecord {
	return verdictRecord{Schema: 1, State: Pending, Tree: plan.Tree, Oracle: plan.Oracle, StartedAt: now.UTC().Truncate(time.Second).Format(time.RFC3339), OwnerPID: os.Getpid()}
}

func retainGreen(replacer durablefile.Replacer, root string, plan subject, recordedAt time.Time) error {
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return err
	}
	dir := filepath.Join(gitdir, "bench-gate-evidence")
	if err := ensureEvidenceDir(gitdir, dir); err != nil {
		return err
	}
	record := verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Grade: evidenceName(plan), RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339)}
	return durableReplaceAt(replacer, dir, record.Grade, record)
}

func ensureEvidenceDir(parent, dir string) error {
	if err := os.Mkdir(dir, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
		return errors.New("invalid evidence directory")
	}
	parentFile, err := os.Open(parent)
	if err != nil {
		return err
	}
	defer parentFile.Close()
	return parentFile.Sync()
}
