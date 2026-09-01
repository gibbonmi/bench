package contract

import (
	"os"
	"time"
)

type MarkerWaitMiss string

const (
	MarkerWaitFast   MarkerWaitMiss = "fast"
	MarkerWaitSlow   MarkerWaitMiss = "slow"
	MarkerWaitExited MarkerWaitMiss = "exited"
)

// WaitForTwoLegMarkers polls fastPath until fastDeadline, then polls slowPath
// until slowDeadline. It accepts one stat function, an optional exit signal,
// and clock and sleep functions. A marker observed by stat takes precedence
// over an exit signal. It returns empty after both markers, MarkerWaitFast or
// MarkerWaitSlow after that leg's deadline, or MarkerWaitExited on exit.
func WaitForTwoLegMarkers(fastPath, slowPath string, fastDeadline, slowDeadline time.Duration, stat func(string) (os.FileInfo, error), exit <-chan struct{}, now func() time.Time, sleep func(time.Duration)) MarkerWaitMiss {
	if miss := waitForMarker(fastPath, fastDeadline, MarkerWaitFast, stat, exit, now, sleep); miss != "" {
		return miss
	}
	if miss := waitForMarker(slowPath, slowDeadline, MarkerWaitSlow, stat, exit, now, sleep); miss != "" {
		return miss
	}
	return ""
}

func waitForMarker(path string, deadline time.Duration, missed MarkerWaitMiss, stat func(string) (os.FileInfo, error), exit <-chan struct{}, now func() time.Time, sleep func(time.Duration)) MarkerWaitMiss {
	until := now().Add(deadline)
	for now().Before(until) {
		if _, err := stat(path); err == nil {
			return ""
		}
		select {
		case <-exit:
			return MarkerWaitExited
		default:
		}
		sleep(10 * time.Millisecond)
	}
	return missed
}
