package subprocess_test

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

const (
	markerEnv  = "SUBPROCESS_HELPER_MARKER"
	gatePIDEnv = "SUBPROCESS_HELPER_GATE_PID"
)

func TestMain(m *testing.M) {
	processtest.Main(m, map[string]func() int{"actor": actorHelper})
}

// gateAdapter is the production directory adapter of the resource that a gated launch registers.
var gateAdapter = protection.Directory{Kind: protection.KindKitTest, Noun: "run"}

func gateProtocol() *protection.Protocol {
	return &protection.Protocol{Validators: protection.Validators{gateAdapter.Kind: gateAdapter.Validate}, ReadLimit: bounds.ControlRecordLimit}
}

// gateResource opens one real resource generation under p in a private directory.
func gateResource(t *testing.T, p *protection.Protocol) protection.Generation {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "run")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	g, err := gateAdapter.Open(p, dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = g.Discard() })
	return g
}

// TestStartGateActorDeathBeforeOpen is PL128. The actor dies after its gate starts and before
// the open signal, so the gate reads end of file and the program never runs.
func TestStartGateActorDeathBeforeOpen(t *testing.T) {
	g := gateResource(t, gateProtocol())
	dir := t.TempDir()
	marker, gatePID := filepath.Join(dir, "started"), filepath.Join(dir, "gate.pid")
	actor := processtest.Helper(t, "actor", protection.Environ(nil, g.Resources())[0], markerEnv+"="+marker, gatePIDEnv+"="+gatePID)
	// The held stdin keeps the actor inside its Start until the test kills it.
	hold, release, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = hold.Close(), release.Close() })
	actor.Stdin = hold
	if err := actor.Start(); err != nil {
		t.Fatal(err)
	}
	pid := processtest.WaitPID(t, gatePID)
	_ = actor.Process.Kill()
	_ = actor.Wait()
	processtest.WaitAbsent(t, pid)
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("the program ran after its actor died before the open signal: %v", err)
	}
}

// TestStartGateExecFailure is PL129. A registered program that cannot execute fails behind its
// gate, and the launch reports the program's start error, not a child exit. A gate that cannot
// enter the command's directory reports the program's start error too.
func TestStartGateExecFailure(t *testing.T) {
	p := gateProtocol()
	plain := filepath.Join(t.TempDir(), "plain")
	if err := os.WriteFile(plain, []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	homeless := exec.Command("true")
	homeless.Dir = "/definitely/missing"
	for name, c := range map[string]struct {
		cmd  *exec.Cmd
		want error
	}{
		"missing":        {exec.Command("/definitely/missing/bench-child"), fs.ErrNotExist},
		"not executable": {exec.Command(plain), fs.ErrPermission},
		"missing dir":    {homeless, fs.ErrNotExist},
	} {
		g := gateResource(t, p)
		registration := subprocess.Registration{Protocol: p, Resources: g.Resources()}
		o, _ := waitWithin(t, started(t, subprocess.NewLifetime(c.cmd, policy(0, final, 0), subprocess.Streams{}, registration)), final)
		var pathErr *fs.PathError
		if o.Process.Started || o.Process.Waited || !errors.As(o.Process.StartErr, &pathErr) || pathErr.Path != c.cmd.Path || !errors.Is(o.Process.StartErr, c.want) {
			t.Fatalf("%s: process=%+v, want a start error that names %s and %v", name, o.Process, c.cmd.Path, c.want)
		}
		if err := g.Release("run", nil); o.Protection.State != protection.Safe || err != nil {
			t.Fatalf("%s: protection=%+v release=%v, want no remaining user", name, o.Protection, err)
		}
	}
}

// TestStartGateStartedFailure is the Started half of PL128. The store refuses the Started write
// after the gate starts, so the gate stays closed and the launch reports an uncertain user. The
// held group signal gives an opened gate time to run the program before the KILL lands.
func TestStartGateStartedFailure(t *testing.T) {
	g := gateResource(t, gateProtocol())
	marker := filepath.Join(t.TempDir(), "started")
	l := subprocess.NewLifetime(exec.Command("touch", marker), policy(0, final, 0), subprocess.Streams{}, subprocess.Registration{Protocol: gateProtocol(), Resources: g.Resources()})
	subprocess.Hold(l, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error {
		if err := real(cmd); err != nil {
			return err
		}
		for _, d := range g.Resources() {
			sealStore(t, d.Store)
		}
		return nil
	}, nil)
	subprocess.HoldSignal(l, func(real func(int, syscall.Signal) error, target int, sig syscall.Signal) error {
		deadline := time.Now().Add(grace)
		for time.Now().Before(deadline) {
			if _, err := os.Stat(marker); err == nil {
				break
			}
			time.Sleep(time.Millisecond)
		}
		return real(target, sig)
	})
	o, _ := waitWithin(t, started(t, l), grace+final)
	if _, err := os.Stat(marker); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("the program ran after its Started write failed: %v", err)
	}
	if o.Protection.State != protection.Uncertain || o.Process.Signal != syscall.SIGKILL {
		t.Fatalf("process=%+v protection=%+v, want a killed gate and an uncertain user", o.Process, o.Protection)
	}
}

// TestRunInheritedRegistersOnLocalProtocol is RS16. The inherited vector is empty and the test
// installs no protocol, so only the protocol that the registration names can register the
// launch, and the child environment then names the local resource.
func TestRunInheritedRegistersOnLocalProtocol(t *testing.T) {
	t.Setenv(protection.Env, "")
	p := gateProtocol()
	g := gateResource(t, p)
	cmd := exec.Command("sh", "-c", `printf %s "$`+protection.Env+`"`)
	streams := subprocess.Streams{Stdout: subprocess.Output{Capture: true}}
	o, err := subprocess.RunInherited(context.Background(), cmd, policy(0, final, 0), streams, subprocess.Interrupt(syscall.SIGKILL), subprocess.Registration{Protocol: p, Resources: g.Resources()})
	if err != nil {
		t.Fatal(err)
	}
	vector, err := protection.ParseVector(string(o.Streams.Stdout.Data))
	if err != nil || !slices.Equal(vector, g.Resources()) {
		t.Fatalf("child vector = %+v, %v; want the local resource %+v", vector, err, g.Resources())
	}
}

// sealStore makes every directory under store read-only until the test ends.
func sealStore(t *testing.T, store string) {
	t.Helper()
	var dirs []string
	if err := filepath.WalkDir(store, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.IsDir() {
			dirs = append(dirs, path)
		}
		return err
	}); err != nil {
		t.Fatalf("walk %s: %v", store, err)
	}
	for _, dir := range dirs {
		if err := os.Chmod(dir, 0o500); err != nil {
			t.Error(err)
		}
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	}
}

// actorHelper registers a gated launch of a marker program and publishes the gate's PID. Its
// Start then blocks on stdin, before Started and the open signal, until the test kills it.
func actorHelper() int {
	inherited, err := protection.Inherited()
	if err != nil {
		return 2
	}
	cmd := exec.Command("sh", "-c", `touch "$0"`, os.Getenv(markerEnv))
	l := subprocess.NewLifetime(cmd, subprocess.Policy{}, subprocess.Streams{}, subprocess.Registration{Protocol: gateProtocol(), Resources: inherited})
	subprocess.Hold(l, func(real func(*exec.Cmd) error, cmd *exec.Cmd) error {
		if err := real(cmd); err != nil {
			return err
		}
		path := os.Getenv(gatePIDEnv)
		if os.WriteFile(path+".partial", []byte(strconv.Itoa(cmd.Process.Pid)), 0o600) == nil && os.Rename(path+".partial", path) == nil {
			_, _ = io.Copy(io.Discard, os.Stdin)
		}
		return nil
	}, nil)
	if l.Start() != nil {
		return 2
	}
	l.Wait()
	return 0
}
