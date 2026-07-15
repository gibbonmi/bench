// Package sessioninspect runs the bounded SessionStart inspection sequence.
package sessioninspect

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/worktree"
)

type phase func(context.Context, io.Writer, io.Writer, string) int

const commandTimeout = 10 * time.Second

type stderrKey struct{}

var phases = []phase{resumePhase, statusPhase, guardsPhase}
var runInspect = Inspect

func Inspect(ctx context.Context, w io.Writer, root string) int {
	stderr, _ := ctx.Value(stderrKey{}).(io.Writer)
	if stderr == nil {
		stderr = w
	}
	for _, run := range phases {
		done := make(chan int, 1)
		go func(run phase) { done <- run(ctx, w, stderr, root) }(run)
		select {
		case <-done:
		case <-ctx.Done():
			fmt.Fprintln(stderr, "warning: bench session-inspect: deadline exceeded; session inspection stopped")
			return 0
		}
	}
	return 0
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
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	ctx = context.WithValue(ctx, stderrKey{}, stderr)
	return runInspect(ctx, stdout, root)
}

func resumePhase(ctx context.Context, stdout, stderr io.Writer, root string) int {
	_ = ctx
	_ = root
	var phaseErr bytes.Buffer
	code := worktree.ResumeCleanCommand(nil, stdout, &phaseErr)
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
