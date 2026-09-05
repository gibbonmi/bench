// Package sessioninspect runs the bounded SessionStart inspection sequence.
package sessioninspect

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/worktree"
)

type phase func(context.Context, io.Writer, io.Writer, string) int

type stderrKey struct{}

var phases = []phase{environmentPhase, resumePhase, statusPhase, guardsPhase}
var runInspect = Inspect

func Inspect(ctx context.Context, w io.Writer, root string) int {
	stderr, _ := ctx.Value(stderrKey{}).(io.Writer)
	if stderr == nil {
		stderr = w
	}
	for _, run := range phases {
		done := make(chan int, 1)
		go func(run phase) { done <- run(ctx, w, stderr, root) }(run)
		if !phaseFinished(ctx, done) {
			fmt.Fprintln(stderr, "warning: bench session-inspect: deadline exceeded; session inspection stopped")
			return 0
		}
	}
	return 0
}

func phaseFinished(ctx context.Context, done <-chan int) bool {
	select {
	case <-done:
		return ctx.Err() == nil
	case <-ctx.Done():
		return false
	}
}

func Command(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: bench session-inspect")
		return 0
	}
	root, err := git.Root()
	if err != nil {
		return 0
	}
	ctx, cancel := bounds.Context(context.Background(), bounds.ProviderTimeout)
	defer cancel()
	ctx = context.WithValue(ctx, stderrKey{}, stderr)
	return runInspect(ctx, stdout, root)
}

func environmentPhase(ctx context.Context, stdout, _ io.Writer, root string) int {
	info, err := os.Lstat(filepath.Join(root, "go.mod"))
	if err != nil || !info.Mode().IsRegular() {
		return 0
	}
	if _, err := exec.LookPath("go"); err == nil {
		return 0
	}
	command := exec.Command("bash", "-c", "exec bash -lc 'command -v go' 2>/dev/null")
	command.Dir = root
	command.Env = capability.WithoutEnvironment(capability.WithoutEnvironment(os.Environ(), "ENVMAN_LOAD"), "BASH_ENV")
	result := bounds.Run(ctx, bounds.EnvironmentDiscoveryTimeout, command)
	executable, valid := discoveredExecutable(result)
	if !valid {
		fmt.Fprintln(stdout, "bench: Go is absent from PATH and the clean Bash login did not resolve an executable Go toolchain.")
		return 0
	}
	fmt.Fprintf(stdout, "bench: environment closure is partial: Go is absent from PATH, but the clean Bash login resolves %s (ENVMAN_LOAD=%q).\n", executable, os.Getenv("ENVMAN_LOAD"))
	fmt.Fprintf(stdout, "bench: recover without replacing harness tools: export PATH=%s:\"$PATH\"\n", sanitize.ShellQuote(filepath.Dir(executable)))
	return 0
}

func discoveredExecutable(result bounds.ProcessResult) (string, bool) {
	if result.Status != bounds.ProcessComplete {
		return "", false
	}
	line := string(result.Output)
	if strings.HasSuffix(line, "\n") {
		line = strings.TrimSuffix(line, "\n")
	}
	if line == "" || !utf8.ValidString(line) || !filepath.IsAbs(line) {
		return "", false
	}
	for _, char := range line {
		if unicode.IsControl(char) {
			return "", false
		}
	}
	info, err := os.Stat(line)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		return "", false
	}
	return line, true
}

func resumePhase(ctx context.Context, stdout, stderr io.Writer, root string) int {
	_ = ctx
	code := worktree.ResumeCleanCommand(root, worktree.Home(), nil, stdout, stderr)
	if code != 0 {
		fmt.Fprintln(stderr, "warning: bench session-start: resume-clean failed; inspect retained worktree state")
	}
	return code
}

func statusPhase(ctx context.Context, stdout, _ io.Writer, root string) int {
	_ = ctx
	_ = root
	out, code := status.Command(nil)
	fmt.Fprint(stdout, out)
	return code
}

func guardsPhase(ctx context.Context, stdout, _ io.Writer, root string) int {
	_ = ctx
	_ = root
	out, code := guards.Command([]string{"--brief"})
	fmt.Fprint(stdout, out)
	return code
}
