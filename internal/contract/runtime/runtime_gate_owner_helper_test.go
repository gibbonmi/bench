package runtime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

type story5GateOwner struct {
	cmd      *exec.Cmd
	gatePGID int
	exitCh   chan error
}

func startStory5GateOwner(t *testing.T, f contract.Fixture) *story5GateOwner {
	t.Helper()
	cmd := exec.Command("bash", benchPath(t), "gate")
	cmd.Dir = f.Root
	cmd.Env = surfaceEnv(f, nil)
	var childOut bytes.Buffer
	cmd.Stdout = &childOut
	cmd.Stderr = &childOut
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	// This goroutine is the single owner of cmd.Wait; stop() below only
	// receives from exitCh, it never calls Wait itself. A second Wait on
	// the same process is undefined behavior once the first has reaped it.
	exitCh := make(chan error, 1)
	exitReady := make(chan struct{})
	go func() {
		exitCh <- cmd.Wait()
		close(exitReady)
	}()
	miss := contract.WaitForTwoLegMarkers(
		filepath.Join(gitDir(t, f), "bench-gate-owner"),
		filepath.Join(gitDir(t, f), "story5-owner-started"),
		5*time.Second, 60*time.Second, os.Stat, exitReady, time.Now, time.Sleep)
	if miss == "" {
		var pgid int
		if _, err := fmt.Sscanf(strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(gitDir(t, f), "story5-gate-pgid")))), "%d", &pgid); err != nil {
			t.Fatal(err)
		}
		return &story5GateOwner{cmd: cmd, gatePGID: pgid, exitCh: exitCh}
	}
	if miss == contract.MarkerWaitExited {
		select {
		case err := <-exitCh:
			t.Fatalf("gate owner exited before reaching pending state: err=%v state=%v\nchild output:\n%s", err, cmd.ProcessState, childOut.String())
		default:
			t.Fatal("gate owner exit signal was not available")
		}
	}
	var stateMsg string
	select {
	case err := <-exitCh:
		stateMsg = fmt.Sprintf("child EXITED early: err=%v state=%v", err, cmd.ProcessState)
	default:
		pgid := cmd.Process.Pid
		// Per-thread dump: in a Go child the main thread futex-waits while the
		// thread doing the real blocking syscall is elsewhere, so per-process
		// WCHAN cannot name the blocked site.
		tree, _ := exec.Command("bash", "-c", fmt.Sprintf(
			"ps -eLo pid,tid,pgid,stat,wchan:30,etime,args | awk 'NR==1||$3==%d'", pgid)).CombinedOutput()
		lock, lockErr := os.Stat(filepath.Join(gitDir(t, f), "bench-gate.lock"))
		// SIGQUIT makes the Go runtime print every goroutine stack to stderr
		// (captured in childOut) and exit; wait for the reap before reading
		// childOut so the pipe copier has finished writing it.
		_ = syscall.Kill(-pgid, syscall.SIGQUIT)
		var quitMsg string
		select {
		case err := <-exitCh:
			quitMsg = fmt.Sprintf("child reaped after SIGQUIT: err=%v", err)
		case <-time.After(3 * time.Second):
			// childOut is written by the exec pipe copiers until Wait
			// returns, so every read of it must happen after an exitCh
			// receive; force the reap with a group kill (the direct-child
			// Kill below can't reach group members holding the pipes open)
			// rather than read childOut unreaped.
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
			select {
			case err := <-exitCh:
				quitMsg = fmt.Sprintf("child ignored SIGQUIT and was SIGKILLed: err=%v", err)
			case <-time.After(3 * time.Second):
				t.Fatalf("gate owner did not reach pending state\nchild STILL ALIVE after SIGQUIT and SIGKILL\nlock=%v err=%v\nprocess group (per-thread, pre-SIGQUIT):\n%s\nchild output withheld: child was never reaped, so the exec pipe copiers may still be writing it",
					lock != nil, lockErr, tree)
			}
		}
		stateMsg = fmt.Sprintf("child STILL ALIVE at deadline\nlock=%v err=%v\n%s\nprocess group (per-thread, pre-SIGQUIT):\n%s",
			lock != nil, lockErr, quitMsg, tree)
	}
	_ = cmd.Process.Kill()
	t.Fatalf("gate owner did not reach pending state (%s marker missed)\n%s\nchild output:\n%s", miss, stateMsg, childOut.String())
	return nil
}

func (o *story5GateOwner) stop(t *testing.T) {
	t.Helper()
	_ = syscall.Kill(-o.gatePGID, syscall.SIGKILL)
	_ = syscall.Kill(-o.cmd.Process.Pid, syscall.SIGKILL)
	select {
	case <-o.exitCh: // the single Wait lives in startStory5GateOwner's goroutine
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out reaping Story5 gate owner pgid=%d", o.gatePGID)
	}
}
