package shift

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Validation bounds, pinned by the spec: both iteration caps share [1,100]; the wall
// deadline is a Go duration in (0,24h]. Defaults apply whenever the variable is unset
// or set to the empty string — the two read identically through os.Getenv, so no
// separate unset-vs-empty branch is needed to honor "empty string = unset".
const (
	itersMin             = 1
	itersMax             = 100
	maxItersDefault      = 12
	refactorItersDefault = 4
	maxWallDefault       = 2 * time.Hour
	maxWallCap           = 24 * time.Hour
)

// hasControlByte reports whether s carries any byte below 0x20 or the DEL byte 0x7f.
// This is a stricter, single-purpose check than toon.Representable's cell-escaping
// notion (which tolerates \n/\r/\t inside an already-escaped TOON cell): a shift
// objective is one line of operator intent, never a pre-escaped cell, so every control
// byte is rejected outright, no exceptions.
func hasControlByte(s string) bool {
	for i := 0; i < len(s); i++ {
		if b := s[i]; b < 0x20 || b == 0x7f {
			return true
		}
	}
	return false
}

// validateObjective rejects an empty objective — the "improve the codebase" default is
// gone, an operator must state one — or one carrying a control byte, so hostile or
// accidental text is refused at entry, before it can reach the intent ledger or the
// TOON emitter.
func validateObjective(objective string) error {
	if objective == "" {
		return errors.New("objective required: bench shift <objective...>")
	}
	if hasControlByte(objective) {
		return errors.New("objective contains a control byte")
	}
	return nil
}

// parseBoundedInt reads name from the environment as an integer in [min,max],
// returning def when the variable is unset or empty. A set-but-invalid value
// (non-integer, or out of range) is an error naming the variable and the accepted
// range — replacing envInt's silent fallback, which ran a shift the operator never
// authorized.
func parseBoundedInt(name string, def, min, max int) (int, error) {
	v := os.Getenv(name)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < min || n > max {
		return 0, fmt.Errorf("%s must be an integer in [%d,%d]", name, min, max)
	}
	return n, nil
}

// parseWallDuration reads name from the environment as a Go duration in (0,max],
// returning def when the variable is unset or empty. A set-but-invalid value is an
// error naming the variable and the accepted range. This slice validates
// BENCH_MAX_WALL but does not yet wire a wall timer — that is the next slice.
func parseWallDuration(name string, def, max time.Duration) (time.Duration, error) {
	v := os.Getenv(name)
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 || d > max {
		return 0, fmt.Errorf("%s must be a duration greater than 0 and at most %s", name, max)
	}
	return d, nil
}
