package backgroundstate_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

const (
	readyEnv  = "BACKGROUNDSTATE_OWNER_READY"
	exitedEnv = "BACKGROUNDSTATE_OWNER_EXITED"
)

func TestMain(m *testing.M) {
	processtest.Main(m, map[string]func() int{"owner": ownerHelper})
}

// ownerHelper is a grade owner with no child. It publishes its pid once its SIGTERM handler is
// set. On SIGTERM it waits, writes the exited marker, and exits, so a stop that returns at the
// signal finds no marker.
func ownerHelper() int {
	terms := make(chan os.Signal, 1)
	signal.Notify(terms, syscall.SIGTERM)
	ready := os.Getenv(readyEnv)
	if os.WriteFile(ready+".partial", []byte(strconv.Itoa(os.Getpid())), 0o600) != nil || os.Rename(ready+".partial", ready) != nil {
		return 3
	}
	<-terms
	time.Sleep(300 * time.Millisecond)
	if os.WriteFile(os.Getenv(exitedEnv), nil, 0o600) != nil {
		return 3
	}
	return 0
}

// owner starts one owner helper and answers its pid and its exited marker. The test reaps the
// owner, so its exit removes its pid.
func owner(t *testing.T) (pid int, exited string) {
	t.Helper()
	dir := t.TempDir()
	ready, exited := filepath.Join(dir, "ready"), filepath.Join(dir, "exited")
	cmd := processtest.Helper(t, "owner", readyEnv+"="+ready, exitedEnv+"="+exited)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { _ = cmd.Wait() }()
	return processtest.WaitPID(t, ready), exited
}

// TestStopWaitsForTheOwnerExit is BG113. The stop of a live claim returns only after its owner
// has exited. The stop of a claim of another boot sends nothing, so its pid keeps running.
func TestStopWaitsForTheOwnerExit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), bounds.TestDeadline(0))
	defer cancel()
	t.Run("live claim", func(t *testing.T) {
		pid, exited := owner(t)
		claim := backgroundstate.Claim{Key: key, Worktree: "worktree", Owner: pid, Boot: currentBoot(t)}
		if err := backgroundstate.Stop(ctx, claim); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(exited); err != nil {
			t.Fatalf("Stop returned before the owner exited: %v", err)
		}
		if alive, err := subprocess.Alive(pid); err != nil || alive {
			t.Fatalf("Alive after Stop = %v, %v; want an exited owner", alive, err)
		}
	})
	t.Run("claim of another boot", func(t *testing.T) {
		pid, exited := owner(t)
		claim := backgroundstate.Claim{Key: key, Worktree: "worktree", Owner: pid, Boot: "another-boot"}
		if err := backgroundstate.Stop(ctx, claim); err != nil {
			t.Fatal(err)
		}
		if alive, err := subprocess.Alive(pid); err != nil || !alive {
			t.Fatalf("Alive after the stop of another boot's claim = %v, %v; want the pid untouched", alive, err)
		}
		if _, err := os.Stat(exited); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("the pid of another boot's claim exited: %v", err)
		}
	})
}
