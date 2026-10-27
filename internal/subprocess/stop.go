package subprocess

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// stopPoll is the interval at which Stop observes whether its process has exited.
const stopPoll = 10 * time.Millisecond

// Alive reports whether process pid exists. Only ESRCH proves that it is absent, so every
// other probe failure is an error.
func Alive(pid int) (bool, error) {
	if err := onePID(pid); err != nil {
		return false, err
	}
	switch observation := protection.Classify(syscall.Kill(pid, 0)); observation.State {
	case protection.GroupPresent:
		return true, nil
	case protection.GroupAbsent:
		return false, nil
	default:
		return false, observation.Err
	}
}

// Stop asks process pid to end with SIGTERM, one of CancelSignals, and returns only after pid
// no longer exists. A ctx that ends first ends the wait with the context error.
func Stop(ctx context.Context, pid int) error {
	if err := onePID(pid); err != nil {
		return err
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}
	ticker := time.NewTicker(stopPoll)
	defer ticker.Stop()
	for {
		alive, err := Alive(pid)
		if err != nil || !alive {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// onePID refuses a pid that names no one process: zero and one select the caller's own group
// and every process, and a negative pid selects a group.
func onePID(pid int) error {
	if pid <= 1 {
		return fmt.Errorf("pid %d names no one process", pid)
	}
	return nil
}
