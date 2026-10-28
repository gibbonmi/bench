package backgroundstate

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// The slot record and its lock sit beside the machine queue under the Bench home.
const (
	slotFile = "slot"
	slotLock = "slot.lock"
)

// Slot is the one running grade of the machine: the worktree and the tip that it grades, and
// the pid of its process with the boot in which that process runs. The empty slot holds no
// grade.
type Slot struct {
	Worktree string `json:"worktree"`
	Tip      string `json:"tip"`
	Owner    int    `json:"owner_pid"`
	Boot     string `json:"boot_id"`
}

// SlotFor answers the slot of process owner, which grades tip in worktree, in the current boot.
func SlotFor(worktree, tip string, owner int) (Slot, error) {
	boot, err := protection.BootID()
	if err != nil {
		return Slot{}, err
	}
	return Slot{Worktree: worktree, Tip: tip, Owner: owner, Boot: boot}, nil
}

// WriteSlot installs slot as the slot under the Bench home home.
func WriteSlot(replacer durablefile.Replacer, home string, slot Slot) error {
	if err := slot.valid(); err != nil {
		return err
	}
	if err := ensureDir(filepath.Join(home, queueDir)); err != nil {
		return err
	}
	return writeRecord(replacer, SlotPath(home), slot)
}

// ReadSlot reads the slot under the Bench home home. An absent slot is the empty slot.
func ReadSlot(home string) (Slot, error) {
	var slot Slot
	if err := readRecord(SlotPath(home), &slot, func() error { return slot.valid() }); err != nil {
		return Slot{}, err
	}
	return slot, nil
}

// Live reports whether the process of slot still runs. The empty slot is not live.
func (s Slot) Live() (bool, error) {
	if s == (Slot{}) {
		return false, nil
	}
	return liveOwner(s.Owner, s.Boot)
}

func (s Slot) valid() error {
	if err := required("worktree", s.Worktree, "tip", s.Tip, "boot_id", s.Boot); err != nil {
		return err
	}
	return validOwner(s.Owner)
}

// LockSlot takes the lock of the slot under the Bench home home without a wait, and answers its
// release. held is false when another process holds the lock, and then no release is due.
func LockSlot(home string) (release func() error, held bool, err error) {
	if err := ensureDir(filepath.Join(home, queueDir)); err != nil {
		return nil, false, err
	}
	lock, err := os.OpenFile(filepath.Join(home, queueDir, slotLock), os.O_CREATE|os.O_RDWR, recordMode)
	if err != nil {
		return nil, false, err
	}
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		closeErr := lock.Close()
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, false, closeErr
		}
		return nil, false, errors.Join(err, closeErr)
	}
	return lock.Close, true, nil
}

// SlotPath answers the path of the slot under the Bench home home.
func SlotPath(home string) string {
	return filepath.Join(home, queueDir, slotFile)
}
