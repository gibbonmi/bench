// Package gocache owns the Bench Go build cache location. It is the one source of
// the directory Bench hands to a Go toolchain child, and the one place that writes
// the child's GOCACHE entry.
package gocache

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// Env is the environment name a Go toolchain child reads the build cache from.
const Env = "GOCACHE"

// homeEnv is the one environment name the derivation reads. XDG_CACHE_HOME stays out:
// the gate closure declares no XDG name, so a value the closure cannot see must never
// steer the cache it hands its children.
const homeEnv = "HOME"

// Dir returns the Bench build cache directory for the given environment slice. The
// derivation reads the slice's HOME alone; it never runs `go env` and never reads the
// process environment. A relative HOME counts as absent.
func Dir(env []string) (string, error) {
	home := value(env, homeEnv)
	if !filepath.IsAbs(home) {
		return "", errors.New("HOME is absent or relative; the Bench build cache needs an absolute HOME")
	}
	return filepath.Join(home, ".cache", "bench", "go-build"), nil
}

// Apply returns env with every GOCACHE entry removed and the Bench entry appended, so
// a child is handed one value and an ambient entry never survives.
func Apply(env []string) ([]string, error) {
	dir, err := Dir(env)
	if err != nil {
		return nil, err
	}
	applied := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if name, _, ok := strings.Cut(entry, "="); ok && name == Env {
			continue
		}
		applied = append(applied, entry)
	}
	return append(applied, Env+"="+dir), nil
}

// value answers the last entry for name, which is the value a child's getenv answers
// with when a slice carries a name twice.
func value(env []string, name string) string {
	found := ""
	for _, entry := range env {
		key, val, ok := strings.Cut(entry, "=")
		if ok && key == name {
			found = val
		}
	}
	return found
}

// FromEnv returns the build cache directory the given environment slice carries. A
// process that Apply prepared reads its own directory back through this function, so no
// caller re-derives a path that Apply already decided. A slice without the entry falls
// back to the HOME derivation, because a runner launched from a plain shell carries no
// GOCACHE entry and still reads the same directory. A slice with neither the entry nor an
// absolute HOME answers an empty string and no error: that slice names no directory, and a
// reader with nothing to name has nothing to refuse.
//
// A relative entry is the one refusal. Apply never hands an inbound value to a child, so a
// relative entry reaches this function alone, and a reader that answered it verbatim would
// name a directory that moves with the reader's own working directory.
func FromEnv(env []string) (string, error) {
	if entry := value(env, Env); entry != "" {
		if !filepath.IsAbs(entry) {
			return "", errors.New(Env + " is relative: " + entry + "; the Bench build cache needs an absolute path")
		}
		return entry, nil
	}
	dir, err := Dir(env)
	if err != nil {
		return "", nil
	}
	return dir, nil
}

// Declared reports whether env names the home the derivation reads. A caller that must
// tell "this environment names no cache directory" from "this environment names one the
// derivation refuses" asks here, because the derivation's one input is this package's own
// decision and a second reading of the name would drift from it.
func Declared(env []string) bool { return value(env, homeEnv) != "" }

// Refusal is the one line a Hold caller prints when the hold fails. The gate run, the lane,
// and the focused run each refuse with this text, so an operator reads the same error and
// the same cache path wherever a build stopped. The error and the path pass the control-rune
// filter, because the path comes from HOME, which the operator owns, and no byte of it
// reaches a terminal raw. A hold that failed in the derivation names no path, so the line
// carries the derivation's own error alone.
func Refusal(env []string, err error) string {
	hint := sanitize.Controls(err.Error())
	if dir, dirErr := Dir(env); dirErr == nil {
		hint += " at " + sanitize.Controls(dir)
	}
	return toon.Errorf("cache lock unavailable", hint)
}
