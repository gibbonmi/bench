package gate

// The launch of a background grade: after a lane pass publishes a tip in a kit worktree, the
// launch starts `bench gate-background <worktree> <tip>` through the wrapper of that worktree as
// a detached session, and it returns without a wait for the grade.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/subprocess"
)

// StartBackground starts one background grade of tip, which a lane pass just published in the
// worktree root. kit is the kit that the verb read, and home is the Bench home of the machine
// slot. A grade starts only in the worktree of an active assignment for which KitForRoot answers
// root itself, and only when root holds its executable wrapper, because the grade runs as that
// wrapper. While a live grade holds the slot, no grade starts. The publication stands whatever
// the launch does, so a launch that fails prints one line to stderr and the verb keeps its exit.
func StartBackground(home, kit, root, tip string, stderr io.Writer) {
	if KitForRoot(kit, root) != root || !RealFS().Executable(backgroundWrapper(root)) {
		return
	}
	if _, assigned := intent.AssignmentForWorktree(root); !assigned {
		return
	}
	if err := takeSlot(home, root, tip); err != nil {
		fmt.Fprintf(stderr, "background grade: not started: %v\n", err)
	}
}

// takeSlot starts the grade of tip in root when no live grade holds the slot under home, and
// it records that grade in the slot. It holds the slot lock throughout, so two launches never
// both find the slot empty. A launch that finds the lock held starts none, because the holder
// takes the slot. A slot that does not read retains its grade, so it starts none.
func takeSlot(home, root, tip string) (err error) {
	release, held, err := backgroundstate.LockSlot(home)
	if err != nil || !held {
		return err
	}
	defer func() { err = errors.Join(err, release()) }()
	slot, err := backgroundstate.ReadSlot(home)
	if err != nil {
		return err
	}
	if live, err := slot.Live(); err != nil || live {
		return err
	}
	pid, err := startGrade(root, tip)
	if err != nil {
		return err
	}
	if slot, err = backgroundstate.SlotFor(root, tip, pid); err != nil {
		return err
	}
	return backgroundstate.WriteSlot(durablefile.Native(), home, slot)
}

// startGrade starts the background grade of tip in root through the lifetime owner, as a
// detached session, and answers its pid. The owned group is only the detach helper, which exits
// once the grade starts, so the teardown window of a process bounds it, and a cancellation kills
// it at once.
func startGrade(root, tip string) (int, error) {
	ctx, cancel := bounds.Context(context.Background(), bounds.FixedWindow(bounds.ProcessFinalWindow))
	defer cancel()
	cmd := exec.Command(backgroundWrapper(root), backgroundVerb, root, tip)
	cmd.Dir = root
	detached, err := subprocess.Detach(cmd)
	if err != nil {
		return 0, err
	}
	return detached.PID(subprocess.RunInherited(ctx, detached.Cmd, subprocess.Policy{}, detached.Streams, subprocess.Interrupt(syscall.SIGKILL)))
}

// backgroundWrapper answers the wrapper of the worktree root.
func backgroundWrapper(root string) string {
	return filepath.Join(root, "bin", "bench.sh")
}
