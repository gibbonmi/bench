package subprocess

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

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
// no longer exists. Each expiry of poll observes pid again, and a poll that never expires
// observes it only once. A ctx that ends first ends the wait with the context error.
func Stop(ctx context.Context, pid int, poll Window) error {
	if err := onePID(pid); err != nil {
		return err
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}
	for {
		alive, err := Alive(pid)
		if err != nil || !alive {
			return err
		}
		expired, stop := poll()
		select {
		case <-ctx.Done():
			stop()
			return ctx.Err()
		case <-expired:
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
