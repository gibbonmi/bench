package gate

import (
	"context"
	"os/exec"
	"syscall"
	"time"
)

const processGroupCancelGrace = 2 * time.Second

type processGroupResult struct {
	Code     int
	StartErr error
}

func runProcessGroupCommand(ctx context.Context, cmd *exec.Cmd) processGroupResult {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	if err := cmd.Start(); err != nil {
		if ctx.Err() != nil {
			return processGroupResult{Code: 130, StartErr: err}
		}
		return processGroupResult{Code: 1, StartErr: err}
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return processGroupResult{Code: processExitCode(cmd, err)}
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case <-done:
		case <-time.After(processGroupCancelGrace):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return processGroupResult{Code: 130}
	}
}

func processExitCode(cmd *exec.Cmd, err error) int {
	if err == nil {
		return 0
	}
	if cmd.ProcessState != nil {
		if code := cmd.ProcessState.ExitCode(); code > 0 {
			return code
		}
	}
	return 1
}
