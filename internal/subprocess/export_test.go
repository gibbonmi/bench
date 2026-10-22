package subprocess

import (
	"os/exec"
	"syscall"
)

// HoldSignal wraps the real group signal of a lifetime before it starts. A real group cannot
// refuse a signal on demand, so a test fails it here.
func HoldSignal(l *Lifetime, kill func(real func(int, syscall.Signal) error, target int, sig syscall.Signal) error) {
	realKill := l.kill
	l.kill = func(target int, sig syscall.Signal) error { return kill(realKill, target, sig) }
}

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
