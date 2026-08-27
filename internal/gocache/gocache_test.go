package gocache

import (
	"slices"
	"strings"
	"testing"
)

// C01, C02, C05: the derivation answers from HOME alone, whatever else the slice carries.
func TestDirDerivesFromHomeAlone(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		env  []string
		want string
	}{
		{name: "absolute-home", env: []string{"HOME=/home/agent"}, want: "/home/agent/.cache/bench/go-build"},
		{name: "xdg-present", env: []string{"HOME=/home/agent", "XDG_CACHE_HOME=/var/xdg"}, want: "/home/agent/.cache/bench/go-build"},
		{name: "space-in-home", env: []string{"HOME=/home/a gent"}, want: "/home/a gent/.cache/bench/go-build"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir, err := Dir(tc.env)
			if err != nil {
				t.Fatal(err)
			}
			if dir != tc.want {
				t.Fatalf("Dir = %q, want %q", dir, tc.want)
			}
		})
	}
}

// C03: no absolute HOME is an error that names HOME, never a fall back to Go's default.
func TestDirRefusesWithoutAnAbsoluteHome(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		env  []string
	}{
		{name: "absent", env: []string{"PATH=/usr/bin", "XDG_CACHE_HOME=/var/xdg"}},
		{name: "empty", env: []string{"HOME="}},
		{name: "relative", env: []string{"HOME=home/agent", "XDG_CACHE_HOME=/var/xdg"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir, err := Dir(tc.env)
			if err == nil {
				t.Fatalf("Dir = %q, want an error", dir)
			}
			if !strings.Contains(err.Error(), "HOME") {
				t.Fatalf("error = %q, want it to name HOME", err)
			}
		})
	}
}

// C04: apply replaces an existing entry rather than leaving the ambient one in place.
func TestApplyReplacesAnExistingEntry(t *testing.T) {
	t.Parallel()
	applied, err := Apply([]string{"GOCACHE=/ambient/cache", "HOME=/home/agent", "PATH=/usr/bin", "GOCACHE=/second/cache"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"HOME=/home/agent", "PATH=/usr/bin", "GOCACHE=/home/agent/.cache/bench/go-build"}
	if !slices.Equal(applied, want) {
		t.Fatalf("Apply = %#v, want %#v", applied, want)
	}
}

func TestApplyReportsADerivationRefusal(t *testing.T) {
	t.Parallel()
	applied, err := Apply([]string{"GOCACHE=/ambient/cache"})
	if err == nil {
		t.Fatalf("Apply = %#v, want an error", applied)
	}
	if !strings.Contains(err.Error(), "HOME") {
		t.Fatalf("error = %q, want it to name HOME", err)
	}
}
