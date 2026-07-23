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

// WaitForTwoLegMarkers waits for the fast marker before allowing the slow marker.
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
