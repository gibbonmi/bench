// Package terminal owns the shared stdin terminal predicate for human-attended
// commands.
package terminal

import (
	"io"
	"os"
)

func IsTerminal(stdin io.Reader) bool {
	file, ok := stdin.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
