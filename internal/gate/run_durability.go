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
	record := verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: plan.Tree, Oracle: plan.Oracle, RecordedAt: recordedAt.UTC().Truncate(time.Second).Format(time.RFC3339)}
	return durableReplaceAt(replacer, dir, evidenceName(plan), record)
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
