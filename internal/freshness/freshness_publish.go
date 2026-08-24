// Atomic publication of an executable and its seal for package freshness.
package freshness

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess"
)

// publicationBackupPattern names the copies a publication keeps of the pair it replaces.
const publicationBackupPattern = ".bench-publish-backup-*"

// sealTemporaryPattern names the staging file a seal write promotes into path.
func sealTemporaryPattern(path string) string { return filepath.Base(path) + ".tmp-*" }

var replacePublicationFile = os.Rename

// Publish replaces executable with staged and writes its matching content seal.
func Publish(root, staged, executable string) error {
	if err := publicationTarget(executable); err != nil {
		return fmt.Errorf("publish executable %q: %w", executable, err)
	}
	if err := publicationTarget(sealPath(executable)); err != nil {
		return fmt.Errorf("publish seal %q: %w", sealPath(executable), err)
	}
	encoded, err := sealContents(root, staged)
	if err != nil {
		return err
	}
	transaction, err := beginPublication(executable)
	if err != nil {
		return err
	}
	defer transaction.close()
	if err := transaction.install(staged); err != nil {
		return fmt.Errorf("publish executable: %w", err)
	}
	if err := transaction.sealWith(encoded); err != nil {
		if restoreErr := transaction.rollback(); restoreErr != nil {
			return fmt.Errorf("publish seal: %w; restore prior pair: %v", err, restoreErr)
		}
		return fmt.Errorf("publish seal: %w", err)
	}
	return nil
}

// sealContents encodes the seal that describes staged as built from root's sources.
func sealContents(root, staged string) ([]byte, error) {
	sources, err := Digest(root)
	if err != nil {
		return nil, err
	}
	binary, err := ExecutableDigest(staged)
	if err != nil {
		return nil, fmt.Errorf("stage executable %q: %w", staged, err)
	}
	return json.Marshal(seal{
		Schema:     sealSchema,
		Sources:    sources,
		Executable: binary,
	})
}

// publicationStepGrace bounds how long an arriving termination waits for a publication
// step to finish before restoring anyway. A step is one rename or one small write. A step
// still running after the grace means the filesystem is wedged. Honoring the termination
// then matters more than waiting for a step that may never return.
const publicationStepGrace = 2 * time.Second

// publication owns replacing an executable and its seal as one outcome. The pair is only
// consistent before the executable moves and after the seal lands. The transaction both
// restores the prior pair when the seal fails and answers a termination arriving in
// between. The invoking shell holds no rollback state and needs none.
type publication struct {
	executable    string
	backup        string
	hadExecutable bool
	sealBackup    string
	hadSeal       bool
	// step serializes the two renames against the termination restore, so a signal lands
	// strictly between steps rather than inside one.
	step     sync.Mutex
	resolved atomic.Bool
	signals  chan os.Signal
	// sealTemporary holds the staging file of a seal write that has not yet landed. The
	// termination handler reads it while that write is still in flight. The name lives on
	// the transaction rather than only in the writer's own frame.
	sealTemporary atomic.Pointer[string]
}

func beginPublication(executable string) (*publication, error) {
	transaction := &publication{executable: executable}
	backup, hadExecutable, err := publicationBackup(executable, true)
	if err != nil {
		return nil, fmt.Errorf("backup executable: %w", err)
	}
	transaction.backup, transaction.hadExecutable = backup, hadExecutable
	sealBackup, hadSeal, err := publicationBackup(sealPath(executable), false)
	if err != nil {
		transaction.remove(backup)
		return nil, fmt.Errorf("backup seal: %w", err)
	}
	transaction.sealBackup, transaction.hadSeal = sealBackup, hadSeal
	transaction.watch()
	return transaction, nil
}

// watch arms the transaction against the terminations an operator or a supervisor sends.
// The restore has to run before the process ends. A process cannot honor a signal and
// keep running as if it had not arrived. The handler therefore exits under the shell
// convention for the signal it received.
func (p *publication) watch() {
	p.signals = make(chan os.Signal, 1)
	signal.Notify(p.signals, subprocess.CancelSignals...)
	go func() {
		received, delivered := <-p.signals
		if !delivered {
			return
		}
		held := p.awaitStep()
		if !p.resolved.Swap(true) {
			_ = p.restore()
		}
		p.removeTemporaries()
		if held {
			p.step.Unlock()
		}
		publicationExit(received)
	}()
}

func (p *publication) awaitStep() bool {
	deadline := time.Now().Add(publicationStepGrace)
	for {
		if p.step.TryLock() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (p *publication) install(staged string) error {
	p.step.Lock()
	defer p.step.Unlock()
	return replacePublicationFile(staged, p.executable)
}

func (p *publication) sealWith(encoded []byte) error {
	p.step.Lock()
	defer p.step.Unlock()
	if err := writeSeal(sealPath(p.executable), encoded, p.trackSealTemporary); err != nil {
		return err
	}
	// Resolving under the same lock the restore takes stops a termination arriving on
	// the heels of a landed seal from undoing a complete publication.
	p.resolved.Store(true)
	return nil
}

func (p *publication) rollback() error {
	p.step.Lock()
	defer p.step.Unlock()
	if p.resolved.Swap(true) {
		return nil
	}
	return p.restore()
}

func (p *publication) restore() error {
	if err := restorePriorFile(p.backup, p.executable, p.hadExecutable); err != nil {
		return err
	}
	return restorePriorFile(p.sealBackup, sealPath(p.executable), p.hadSeal)
}

func (p *publication) close() {
	signal.Stop(p.signals)
	close(p.signals)
	p.removeTemporaries()
}

// removeTemporaries deletes every temporary the transaction owns. Both the termination
// handler and close call it, and either may run after the other. A removal of what is
// already gone is the normal case, not a failure.
func (p *publication) removeTemporaries() {
	p.remove(p.backup)
	p.remove(p.sealBackup)
	if temporary := p.sealTemporary.Load(); temporary != nil {
		p.remove(*temporary)
	}
}

func (p *publication) trackSealTemporary(path string) {
	if path == "" {
		p.sealTemporary.Store(nil)
		return
	}
	p.sealTemporary.Store(&path)
}

func (p *publication) remove(backup string) {
	if backup != "" {
		_ = os.Remove(backup)
	}
}

// restorePriorFile puts back what publication found at path: the backed-up bytes, or
// nothing at all where there was nothing before.
func restorePriorFile(backup, path string, existed bool) error {
	if existed {
		return os.Rename(backup, path)
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func publicationExit(received os.Signal) {
	code := 128
	if signaled, ok := received.(syscall.Signal); ok {
		code += int(signaled)
	}
	os.Exit(code)
}

func publicationBackup(path string, executable bool) (string, bool, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	contents, err := secureContents(path, executable)
	if err != nil {
		return "", false, err
	}
	backup, err := os.CreateTemp(filepath.Dir(path), publicationBackupPattern)
	if err != nil {
		return "", false, err
	}
	name := backup.Name()
	if _, err := backup.Write(contents); err != nil {
		_ = backup.Close()
		_ = os.Remove(name)
		return "", false, err
	}
	if err := backup.Chmod(info.Mode().Perm()); err != nil {
		_ = backup.Close()
		_ = os.Remove(name)
		return "", false, err
	}
	if err := backup.Close(); err != nil {
		_ = os.Remove(name)
		return "", false, err
	}
	return name, true, nil
}

func publicationTarget(path string) error {
	if err := rejectSymlinkComponents(filepath.Dir(path)); err != nil {
		return err
	}
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("is a symbolic link")
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("is not a regular file")
	}
	return nil
}
