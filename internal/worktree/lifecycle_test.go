package worktree

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestReclaimable pins the four-way lease decision the black-box lease-hardening
// contract exercises but cannot cheaply enumerate: a recorded numeric pid gates on
// liveness, and unreadable/empty content reclaims only once aged past the threshold —
// so a fresh-empty writer mid-claim is never stolen while a legacy/crashed lease is.
func TestReclaimable(t *testing.T) {
	now := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	dead := func(int) bool { return false }
	live := func(int) bool { return true }

	cases := []struct {
		name    string
		content string
		age     time.Duration
		alive   func(int) bool
		want    bool
	}{
		{"live pid respected", "4242 2026-07-04T11:59:00Z", 30 * time.Second, live, false},
		{"dead pid reclaimed", "4242 2026-07-04T11:59:00Z", 30 * time.Second, dead, true},
		{"non-numeric legacy aged out reclaimed", "garbage content", 2 * time.Minute, dead, true},
		{"fresh empty writer mid-claim respected", "", 5 * time.Second, dead, false},
		{"aged-out empty crash reclaimed", "", 2 * time.Minute, dead, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mtime := now.Add(-tc.age)
			if got := reclaimable([]byte(tc.content), mtime, now, tc.alive); got != tc.want {
				t.Errorf("reclaimable(%q, age %s) = %v, want %v", tc.content, tc.age, got, tc.want)
			}
		})
	}
}

// TestCandidateNameStaysInPool pins that a minted candidate never escapes the pool
// directory — a wrong name would mint outside the pool and silently break warm reuse.
func TestCandidateNameStaysInPool(t *testing.T) {
	pool := "/home/x/.bench/worktrees/bench-123"
	got := candidateName(pool, 1751630400, 4242, 2)
	if filepath.Dir(got) != pool {
		t.Errorf("candidateName parent = %q, want %q", filepath.Dir(got), pool)
	}
	if want := filepath.Join(pool, "1751630400-4242-2"); got != want {
		t.Errorf("candidateName = %q, want %q", got, want)
	}
	if strings.ContainsAny(filepath.Base(got), "/") {
		t.Errorf("candidate base %q escaped the pool", filepath.Base(got))
	}
}
