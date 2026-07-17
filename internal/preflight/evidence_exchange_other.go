//go:build !linux && !darwin

package preflight

import (
	"fmt"
	"runtime"
)

func atomicExchangeDirs(_, _ string) error {
	return fmt.Errorf("atomic evidence replacement is unsupported on %s", runtime.GOOS)
}
