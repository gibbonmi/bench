package subprocess

import (
	"fmt"
	"os/exec"
	"syscall"
)

// Replace replaces the image of the running process with the program name, run with args
// and env. It resolves name as exec.Command does. The process keeps its PID, its group, and
// its open descriptors, so a parent that waits on it waits on the program. Replace returns
// only when the replacement fails.
func Replace(name string, args, env []string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return err
	}
	if err := syscall.Exec(path, append([]string{name}, args...), env); err != nil {
		return fmt.Errorf("replace the process image with %s: %w", path, err)
	}
	return nil
}
