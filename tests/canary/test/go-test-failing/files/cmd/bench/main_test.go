package main

import "testing"

// Intentionally failing: proves the gate's `go test` check bites. If that check rots
// into an always-pass, this fixture stops going red and the canary sweep catches it.
func TestCanaryFails(t *testing.T) {
	t.Fatal("canary: intentional failure")
}
