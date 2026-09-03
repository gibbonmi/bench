package gocache

import (
	"errors"
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

// R18: FromEnv answers the entry when the slice carries one, falls back to the HOME
// derivation when it does not, and answers an empty string when the slice names neither.
// An empty entry is absent, so it reaches the derivation rather than a refusal.
func TestFromEnvFallsBackToTheHomeDerivation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		env  []string
		want string
	}{
		{"entry", []string{"HOME=/home/agent", "GOCACHE=/ambient/cache"}, "/ambient/cache"},
		{"home derivation", []string{"HOME=/home/agent", "PATH=/usr/bin"}, "/home/agent/.cache/bench/go-build"},
		{"empty entry", []string{"HOME=/home/agent", "GOCACHE="}, "/home/agent/.cache/bench/go-build"},
		{"neither", []string{"PATH=/usr/bin"}, ""},
		{"relative home", []string{"HOME=agent"}, ""},
	} {
		got, err := FromEnv(tc.env)
		if err != nil {
			t.Errorf("%s: FromEnv = %v, want no refusal", tc.name, err)
			continue
		}
		if got != tc.want {
			t.Errorf("%s: FromEnv = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// C20: a relative entry is a refusal that names the value. A verbatim answer would report
// a directory that moves with the reader's own working directory.
func TestFromEnvRefusesARelativeEntry(t *testing.T) {
	t.Parallel()
	dir, err := FromEnv([]string{"HOME=/home/agent", "GOCACHE=cache"})
	if err == nil {
		t.Fatalf("FromEnv = %q, want a refusal", dir)
	}
	if dir != "" {
		t.Errorf("FromEnv directory = %q, want an empty string beside the refusal", dir)
	}
	if !strings.Contains(err.Error(), "cache") {
		t.Errorf("refusal = %q, want it to name the entry value", err)
	}
}

// C21: Declared separates an environment that names no home from one whose home the
// derivation refuses, which is the distinction the gate closure reads.
func TestDeclaredReportsTheHomeName(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		env  []string
		want bool
	}{
		{"absolute", []string{"HOME=/home/agent"}, true},
		{"relative", []string{"HOME=agent"}, true},
		{"empty", []string{"HOME="}, false},
		{"absent", []string{"PATH=/usr/bin"}, false},
	} {
		if got := Declared(tc.env); got != tc.want {
			t.Errorf("%s: Declared = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// C22: the shared refusal line names the hold error and the derived path, and it names the
// error alone when the derivation itself failed. A control byte in either reaches no
// terminal raw.
func TestRefusalNamesTheHoldErrorAndThePath(t *testing.T) {
	t.Parallel()
	held := Refusal([]string{"HOME=/home/ag\x1bent"}, errors.New("permission denied"))
	if !strings.Contains(held, "permission denied") || !strings.Contains(held, "/home/ag\\u001bent/.cache/bench/go-build") {
		t.Errorf("refusal = %q, want the error and the escaped path", held)
	}
	if strings.Contains(held, "\x1b") {
		t.Errorf("refusal = %q, want no control byte", held)
	}
	underived := Refusal([]string{"PATH=/usr/bin"}, errors.New("HOME is absent"))
	if !strings.Contains(underived, "HOME is absent") || strings.Contains(underived, " at ") {
		t.Errorf("refusal = %q, want the derivation error alone", underived)
	}
}
