package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRefreshUsesBoundedNoninteractiveFetch(t *testing.T) {
	marker := installFakeRefreshGit(t, "printf '%s|%s' \"$GIT_TERMINAL_PROMPT\" \"$*\" > \"$BENCH_REFRESH_MARKER\"")
	got := Refresh("/repo with space")
	if got.Status != "refreshed" {
		t.Fatalf("Refresh = %#v", got)
	}
	data, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	want := "0|-C /repo with space fetch -q --no-recurse-submodules origin"
	if string(data) != want {
		t.Fatalf("fake git record = %q, want %q", data, want)
	}
}

func TestRefreshFailureAndTimeoutAreNonfatalAndDetailed(t *testing.T) {
	t.Run("exit", func(t *testing.T) {
		installFakeRefreshGit(t, "printf 'origin said no\\n' >&2; exit 23")
		got := Refresh("/repo")
		if got.Status != "failed" || !strings.Contains(got.Detail, "origin said no") {
			t.Fatalf("Refresh = %#v", got)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		installFakeRefreshGit(t, "sleep 5")
		old := refreshTimeout
		refreshTimeout = 20 * time.Millisecond
		t.Cleanup(func() { refreshTimeout = old })
		started := time.Now()
		got := Refresh("/repo")
		if got.Status != "failed" || !strings.Contains(got.Detail, "timeout") || time.Since(started) > time.Second {
			t.Fatalf("Refresh = %#v after %s", got, time.Since(started))
		}
	})
}

func TestRefreshOfflineStartsNoGitAndNamesFlag(t *testing.T) {
	marker := installFakeRefreshGit(t, ": > \"$BENCH_REFRESH_MARKER\"")
	t.Setenv("BENCH_OFFLINE", "1")
	got := Refresh("/repo")
	if got.Status != "offline" || got.Detail != "BENCH_OFFLINE=1" {
		t.Fatalf("Refresh = %#v", got)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("offline refresh started git: %v", err)
	}
}

func installFakeRefreshGit(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	marker := filepath.Join(dir, "marker")
	script := filepath.Join(dir, "git")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BENCH_REFRESH_MARKER", marker)
	return marker
}
