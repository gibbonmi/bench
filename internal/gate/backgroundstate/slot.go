package backgroundstate

import "github.com/gibbonmi/bench/internal/subprocess/protection"

// Slot is the running grade of the machine queue: its entry, and the pid of its process with
// the boot in which that process runs. The empty slot holds no grade.
type Slot struct {
	Entry
	Owner int    `json:"owner_pid"`
	Boot  string `json:"boot_id"`
}

// SlotFor answers the slot of process owner, which runs the grade of entry, in the current boot.
func SlotFor(entry Entry, owner int) (Slot, error) {
	boot, err := protection.BootID()
	if err != nil {
		return Slot{}, err
	}
	return Slot{Entry: entry, Owner: owner, Boot: boot}, nil
}

// ReadSlot reads the running grade of the machine queue under the Bench home home. A queue with
// no running grade answers the empty slot.
func ReadSlot(home string) (Slot, error) {
	queue, err := ReadQueue(home)
	if err != nil || queue.Running == nil {
		return Slot{}, err
	}
	return *queue.Running, nil
}

// Live reports whether the process of slot still runs. The empty slot is not live.
func (s Slot) Live() (bool, error) {
	if s == (Slot{}) {
		return false, nil
	}
	return liveOwner(s.Owner, s.Boot)
}

// Claim answers the claim that the grade of slot holds on its key, so Stop can end that grade
// before it writes its claim.
func (s Slot) Claim() Claim {
	return Claim{Key: s.Key, Worktree: s.Worktree, Owner: s.Owner, Boot: s.Boot}
}

func (s Slot) valid() error {
	if err := s.Entry.valid(); err != nil {
		return err
	}
	if err := required("boot_id", s.Boot); err != nil {
		return err
	}
	return validOwner(s.Owner)
}
