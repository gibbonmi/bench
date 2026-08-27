// Package gocache owns the Bench Go build cache location. It is the one source of
// the directory Bench hands to a Go toolchain child, and the one place that writes
// the child's GOCACHE entry.
package gocache

import (
	"errors"
	"path/filepath"
	"strings"
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
// caller re-derives a path that Apply already decided. An environment without the entry
// answers an empty string.
func FromEnv(env []string) string { return value(env, Env) }
