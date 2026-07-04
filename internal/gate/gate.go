// Package gate is the oracle's selection logic in one Go home: the ordered
// resolution chain (`.bench/gate.sh` beats `$BENCH_GATE` beats auto-detect), the gate
// run from the repo root, and the verdict-cache record keyed to git.TreeHash. Both the
// standalone `bench gate` (via the shell's one-glance run_gate → `bench gate-run`) and
// the in-process shift loop read this package, so gate resolution and the cache-write
// format each live in exactly one place — a second live resolver, or a second cache
// writer, is the worst class of bug in a kit whose premise is "the gate is the oracle".
package gate

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/git"
)

// Kind names the resolved gate. The zero value None is the no-gate case (exit 3,
// nothing recorded); the rest map to a command run from the repo root.
type Kind int

const (
	None Kind = iota
	GateSh
	BenchGate
	Pnpm
	Npm
	Pyproject
	Cargo
)

// Resolution is the chosen gate: its Kind and, for BenchGate, the command string the
// `$BENCH_GATE` env var carried.
type Resolution struct {
	Kind    Kind
	Command string
}

// treeHashRE is the shape a real git tree hash must match before it is written to the
// verdict cache. Anything else (notably git.TreeHash's "none" on failure) is refused,
// so Record never forges a tree — the no-forged-verdict guarantee shared with the Stop
// hook, whose recordGate delegates here.
var treeHashRE = regexp.MustCompile(`^[0-9a-f]+$`)

// FS injects the two filesystem probes Resolve needs — `-x` for the executable
// `.bench/gate.sh` and `-f` for the auto-detect lockfiles — so the resolution
// precedence is a pure function unit-testable without a real tree.
type FS struct {
	Executable func(path string) bool
	Exists     func(path string) bool
}

// RealFS is the production probe set: a regular executable file for Executable, a
// regular (non-directory) file for Exists.
func RealFS() FS {
	return FS{
		Executable: func(p string) bool {
			info, err := os.Stat(p)
			return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
		},
		Exists: func(p string) bool {
			info, err := os.Stat(p)
			return err == nil && !info.IsDir()
		},
	}
}

// Resolve is the ordered chain as a pure function: an executable `.bench/gate.sh`
// wins, then a non-empty `$BENCH_GATE`, then the first auto-detect lockfile in the
// fixed order pnpm → npm → pyproject → cargo, then None. A reordered chain would
// silently run the wrong oracle; this is the precedence the NEW resolution-order
// contract and the table test both pin.
func Resolve(root, benchGate string, fs FS) Resolution {
	if fs.Executable(filepath.Join(root, ".bench", "gate.sh")) {
		return Resolution{Kind: GateSh}
	}
	if benchGate != "" {
		return Resolution{Kind: BenchGate, Command: benchGate}
	}
	for _, d := range []struct {
		file string
		kind Kind
	}{
		{"pnpm-lock.yaml", Pnpm},
		{"package.json", Npm},
		{"pyproject.toml", Pyproject},
		{"Cargo.toml", Cargo},
	} {
		if fs.Exists(filepath.Join(root, d.file)) {
			return Resolution{Kind: d.kind}
		}
	}
	return Resolution{Kind: None}
}

// command builds the shell command a resolution runs from the repo root. None has no
// command (handled by the caller). The auto-detect strings mirror bench.sh's
// best-effort defaults byte-for-byte.
func (r Resolution) command(root string) *exec.Cmd {
	switch r.Kind {
	case GateSh:
		return exec.Command(filepath.Join(root, ".bench", "gate.sh"))
	case BenchGate:
		return exec.Command("bash", "-c", r.Command)
	case Pnpm:
		return exec.Command("bash", "-c", "pnpm -s typecheck && pnpm -s test && pnpm -s lint")
	case Npm:
		return exec.Command("bash", "-c", "npm run -s typecheck && npm test --silent && npm run -s lint")
	case Pyproject:
		return exec.Command("bash", "-c", "mypy . && pytest -q && ruff check .")
	case Cargo:
		return exec.Command("bash", "-c", "cargo test --quiet && cargo clippy -q -- -D warnings")
	default:
		return nil
	}
}

// gateEnv returns the caller's environment with wrapper-routing internals removed.
// `bench gate` reaches this package through bin/bench.sh -> route_binary, which sets
// BENCH_KIT/BENCH_WRAPPER so the binary can find its assets. Those are not part of the
// project gate's contract; leaking them into the gate makes fixture wrappers resolve the
// live kit instead of their own fabricated layout.
func gateEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "BENCH_KIT=") || strings.HasPrefix(kv, "BENCH_WRAPPER=") {
			continue
		}
		env = append(env, kv)
	}
	return env
}

// Run executes the resolved gate from the repo root and returns its exit code, with
// the gate's own output streamed to stdout/stderr. The gate is run from the working
// tree by design (an agent can edit the file it is graded by; the canary tripwire, not
// this call site, keeps that safe). None must not reach here — the caller handles the
// no-gate exit-3-nothing-recorded case.
func Run(root string, res Resolution, stdout, stderr io.Writer) int {
	cmd := res.command(root)
	if cmd == nil {
		return 3
	}
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = gateEnv()
	if err := cmd.Run(); err != nil {
		if cmd.ProcessState != nil {
			if code := cmd.ProcessState.ExitCode(); code > 0 {
				return code
			}
		}
		return 1 // failed to start, or a signal death: treat as red
	}
	return 0
}

// RunContext executes the resolved gate like Run, but puts the gate in its own process
// group and kills that group before returning when ctx is canceled. Shift uses this
// path so an interrupt cannot release the pooled worktree while a gate child keeps
// running and writing into it. Standalone `bench gate` uses Run, preserving normal
// foreground-process signal delivery.
func RunContext(ctx context.Context, root string, res Resolution, stdout, stderr io.Writer) int {
	cmd := res.command(root)
	if cmd == nil {
		return 3
	}
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = gateEnv()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return 1
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			return 0
		}
		if cmd.ProcessState != nil {
			if code := cmd.ProcessState.ExitCode(); code > 0 {
				return code
			}
		}
		return 1 // failed to start, or a signal death: treat as red
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return 130
	}
}

// Record writes the verdict cache "<status> <tree hash> <iso8601>" to
// <git-dir>/bench-last-gate, keyed to the content hash of the tested tree (not the
// commit sha, so a commit-on-green does not stale the verdict that authorized it). rc
// == 0 records green, else red. A tree hash that is not a real hash writes nothing and
// warns — the no-forged-verdict guarantee. This is the ONE verdict-cache writer; the
// Stop hook reaches it through the gate path it runs, never a second implementation.
func Record(root string, rc int, stderr io.Writer) {
	gitdir, err := git.Output("-C", root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return
	}
	status := "green"
	if rc != 0 {
		status = "red"
	}
	tree := git.TreeHash(root)
	if !treeHashRE.MatchString(tree) {
		fmt.Fprintln(stderr, "WARNING: bench tree-hash unavailable — not recording a gate verdict (no forged tree).")
		return
	}
	line := status + " " + tree + " " + time.Now().UTC().Format("2006-01-02T15:04:05Z") + "\n"
	_ = os.WriteFile(filepath.Join(gitdir, "bench-last-gate"), []byte(line), 0o644)
}

// RunAndRecord resolves, runs, and records the gate for root, returning its exit code.
// The no-gate case exits 3 and records nothing (the chain resolved to None); every
// resolved gate runs and then records its verdict. Shared by the `gate-run` subcommand
// and the in-process shift loop, so neither carries its own resolve-run-record chain.
func RunAndRecord(root string, stdout, stderr io.Writer) int {
	res := Resolve(root, os.Getenv("BENCH_GATE"), RealFS())
	if res.Kind == None {
		fmt.Fprintln(stderr, "no gate found: add an executable .bench/gate.sh or set BENCH_GATE")
		return 3
	}
	rc := Run(root, res, stdout, stderr)
	Record(root, rc, stderr)
	return rc
}

// RunAndRecordContext is RunAndRecord with cancellation for in-process callers that
// own teardown. A canceled gate is not recorded as red because the oracle did not
// finish judging the tree.
func RunAndRecordContext(ctx context.Context, root string, stdout, stderr io.Writer) int {
	res := Resolve(root, os.Getenv("BENCH_GATE"), RealFS())
	if res.Kind == None {
		fmt.Fprintln(stderr, "no gate found: add an executable .bench/gate.sh or set BENCH_GATE")
		return 3
	}
	rc := RunContext(ctx, root, res, stdout, stderr)
	if ctx.Err() != nil {
		return rc
	}
	Record(root, rc, stderr)
	return rc
}

// RunCommand is the `bench gate-run [root]` plumbing subcommand: the shell's one-glance
// run_gate forwards here so gate resolution lives in exactly one place. Root is args[0]
// when the shell passes the resolved repo root, else the cwd's repo — resolved so the
// gate always runs from the top level even when invoked from a subdirectory.
func RunCommand(args []string, stdout, stderr io.Writer) int {
	var root string
	if len(args) > 0 && args[0] != "" {
		root = args[0]
	} else {
		r, err := git.Root()
		if err != nil {
			fmt.Fprintln(stderr, "not in a git repo")
			return 1
		}
		root = r
	}
	return RunAndRecord(root, stdout, stderr)
}
