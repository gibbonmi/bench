package contract

import (
	"os"
	"testing"
	"time"
)

type markerWaitClock struct{ now time.Time }

func (c *markerWaitClock) sleep(d time.Duration) { c.now = c.now.Add(d) }

func TestWaitForTwoLegMarkers(t *testing.T) {
	t.Run("fast marker miss uses the fast deadline", func(t *testing.T) {
		clock := &markerWaitClock{}
		miss := WaitForTwoLegMarkers("fast", "slow", 5*time.Second, 60*time.Second,
			func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }, nil,
			func() time.Time { return clock.now }, clock.sleep)
		if miss != MarkerWaitFast || clock.now != (time.Time{}).Add(5*time.Second) {
			t.Fatalf("miss=%q at %s, want fast at 5s", miss, clock.now)
		}
	})

	t.Run("slow marker may use its own deadline", func(t *testing.T) {
		clock := &markerWaitClock{}
		miss := WaitForTwoLegMarkers("fast", "slow", 5*time.Second, 60*time.Second,
			func(path string) (os.FileInfo, error) {
				if path == "fast" || !clock.now.Before((time.Time{}).Add(30*time.Second)) {
					return nil, nil
				}
				return nil, os.ErrNotExist
			}, nil, func() time.Time { return clock.now }, clock.sleep)
		if miss != "" || clock.now != (time.Time{}).Add(30*time.Second) {
			t.Fatalf("miss=%q at %s, want success at 30s", miss, clock.now)
		}
	})

	t.Run("child exit returns immediately", func(t *testing.T) {
		clock := &markerWaitClock{}
		exit := make(chan struct{})
		close(exit)
		miss := WaitForTwoLegMarkers("fast", "slow", 5*time.Second, 60*time.Second,
			func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }, exit,
			func() time.Time { return clock.now }, clock.sleep)
		if miss != MarkerWaitExited || !clock.now.IsZero() {
			t.Fatalf("miss=%q at %s, want immediate exit", miss, clock.now)
		}
	})

	t.Run("observed markers precede exit", func(t *testing.T) {
		for _, tc := range []struct {
			name        string
			slowPresent bool
			want        MarkerWaitMiss
		}{
			{name: "fast marker advances to slow leg", want: MarkerWaitExited},
			{name: "both markers succeed", slowPresent: true},
		} {
			t.Run(tc.name, func(t *testing.T) {
				exit := make(chan struct{})
				close(exit)
				slowObserved := false
				miss := WaitForTwoLegMarkers("fast", "slow", 5*time.Second, 60*time.Second,
					func(path string) (os.FileInfo, error) {
						if path == "fast" {
							return nil, nil
						}
						slowObserved = true
						if tc.slowPresent {
							return nil, nil
						}
						return nil, os.ErrNotExist
					}, exit, time.Now, time.Sleep)
				if !slowObserved {
					t.Fatal("fast marker lost to exit instead of advancing to slow leg")
				}
				if miss != tc.want {
					t.Fatalf("miss=%q, want %q", miss, tc.want)
				}
			})
		}
	})
}
