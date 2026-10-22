package subprocess

import "os/exec"

// Hold wraps the real Start and Wait calls of a lifetime before it starts. A test uses it
// to pause a real launch or a real waiter and to count the waiters it creates.
func Hold(l *Lifetime, start, wait func(real func(*exec.Cmd) error, cmd *exec.Cmd) error) {
	realStart, realWait := l.start, l.wait
	if start != nil {
		l.start = func(cmd *exec.Cmd) error { return start(realStart, cmd) }
	}
	if wait != nil {
		l.wait = func(cmd *exec.Cmd) error { return wait(realWait, cmd) }
	}
}
