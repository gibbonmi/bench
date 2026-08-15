package shift

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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
// returning def when the variable is unset or empty. A zero default means no timer;
// a set-but-invalid value is an error naming the variable and the accepted range.
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
