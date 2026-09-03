// Package runbinary owns the Bench executable selected for one top-level run.
package runbinary

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gocache"
)

// Env carries the exact Bench executable selected by the current top-level run.
const Env = "BENCH_RUN_BINARY"

// BuilderCancelGrace is how long a cancelled builder group has to exit on SIGTERM
// before it is killed. This constant is exported, so a test waiting out the drain
// derives its own deadline from this window, instead of guessing a literal.
const BuilderCancelGrace = 2 * time.Second

type Builder func(context.Context, string, string) error
type Verifier func(string, string) error

// Factory supplies the authoring seams used by tests. Zero values select the private
// operating-system temp directory, canonical builder, and freshness verifier. RepairRoot
// is the checkout a refusal tells the operator to rebuild in; its zero value names the
// source root the selection was graded against. They part only where the digest root is
// a composed tree the operator cannot rebuild in.
type Factory struct {
	TempRoot   string
	RepairRoot string
	Build      Builder
	Verify     Verifier
}

// Selection is the exact source-bound executable used by one run. Close removes bytes
// only for an owner-authored selection. An inherited consumer never owns its lifetime.
type Selection struct {
	Path       string
	SourceRoot string

	dir       string
	closeOnce sync.Once
	closeErr  error
}

// Own creates one private host-target Bench executable from sourceRoot.
func Own(ctx context.Context, sourceRoot string) (*Selection, error) {
	return (Factory{}).Own(ctx, sourceRoot)
}

// ReuseOrOwn consumes an explicitly inherited selection, or authors one when the
// variable is absent. A present invalid value refuses instead of becoming permission
// to build.
func ReuseOrOwn(ctx context.Context, sourceRoot string) (*Selection, error) {
	return (Factory{}).ReuseOrOwn(ctx, sourceRoot)
}

// ReuseOrOwn applies the factory's seams while preserving the inherited-selection rule.
func (f Factory) ReuseOrOwn(ctx context.Context, sourceRoot string) (*Selection, error) {
	if raw, present := os.LookupEnv(Env); present {
		if err := requireSourceRoot(sourceRoot); err != nil {
			return nil, err
		}
		return f.Inherit(sourceRoot, raw)
	}
	return f.Own(ctx, sourceRoot)
}

// Inherit validates the selection supplied by an existing run owner.
func Inherit(sourceRoot string) (*Selection, error) {
	raw, present := os.LookupEnv(Env)
	if !present {
		return nil, fmt.Errorf("%s is required for an inherited Bench run", Env)
	}
	if err := requireSourceRoot(sourceRoot); err != nil {
		return nil, err
	}
	return (Factory{}).Inherit(sourceRoot, raw)
}

// requireSourceRoot keeps a run that resolves its own kit fail-closed. Such a run always
// has a tree to grade its selection against, so an empty root is a resolution defect, not
// permission to fall back to the narrower check.
func requireSourceRoot(sourceRoot string) error {
	if sourceRoot == "" {
		return errors.New("a Bench source root is required to grade an inherited selection")
	}
	return nil
}

// Own creates one private selection through the factory's builder.
func (f Factory) Own(ctx context.Context, sourceRoot string) (*Selection, error) {
	source, err := canonicalSourceRoot(sourceRoot)
	if err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp(f.TempRoot, "bench-run-")
	if err != nil {
		return nil, fmt.Errorf("create private Bench run directory: %w", err)
	}
	selection := &Selection{Path: filepath.Join(dir, "bench"), SourceRoot: source, dir: dir}
	complete := false
	defer func() {
		if !complete {
			_ = selection.Close()
		}
	}()
	build := f.Build
	if build == nil {
		build = canonicalBuild
	}
	if err := build(ctx, source, selection.Path); err != nil {
		return nil, fmt.Errorf("build private Bench executable: %w", err)
	}
	if err := f.validate(source, selection.Path); err != nil {
		return nil, err
	}
	complete = true
	return selection, nil
}

// Inherit validates value without taking ownership of its containing directory. An empty
// sourceRoot names no tree to grade against, which keeps the seal-pair check; every
// entry point that resolves a run's own kit names one.
func (f Factory) Inherit(sourceRoot, value string) (*Selection, error) {
	source := ""
	if sourceRoot != "" {
		canonical, err := canonicalSourceRoot(sourceRoot)
		if err != nil {
			return nil, err
		}
		source = canonical
	}
	if err := f.validate(source, value); err != nil {
		return nil, err
	}
	return &Selection{Path: value, SourceRoot: source}, nil
}

// canonicalVerify is the freshness check a factory with no verification seam applies. A
// named source root proves the seal's source digest, so a self-consistent stale seal
// cannot pass; a caller that names no root keeps the narrower seal-pair check, which is
// all an inherited consumer with no tree of its own can prove.
func (f Factory) canonicalVerify(sourceRoot, executable string) error {
	if sourceRoot == "" {
		return freshness.VerifyExecutable(executable)
	}
	repairRoot := f.RepairRoot
	if repairRoot == "" {
		repairRoot = sourceRoot
	}
	return freshness.VerifyAt(sourceRoot, repairRoot, executable)
}

func (f Factory) validate(sourceRoot, executable string) error {
	if executable == "" {
		return errors.New("selected Bench executable path is empty")
	}
	if !filepath.IsAbs(executable) || filepath.Clean(executable) != executable {
		return fmt.Errorf("selected Bench executable path %q is not a cleaned absolute path", executable)
	}
	info, err := os.Lstat(executable)
	if err != nil {
		return fmt.Errorf("selected Bench executable %q is unavailable: %w", executable, err)
	}
	if !info.Mode().IsRegular() || info.Mode()&0o111 == 0 {
		return fmt.Errorf("selected Bench executable %q is not a regular executable", executable)
	}
	resolved, err := filepath.EvalSymlinks(executable)
	if err != nil || filepath.Clean(resolved) != executable {
		return fmt.Errorf("selected Bench executable %q traverses a symlink", executable)
	}
	verify := f.Verify
	if verify == nil {
		verify = f.canonicalVerify
	}
	if err := verify(sourceRoot, executable); err != nil {
		return fmt.Errorf("selected Bench executable %q does not match source %q: %w", executable, sourceRoot, err)
	}
	return nil
}

// Close releases an owner-authored private directory after its children have returned.
func (s *Selection) Close() error {
	if s == nil || s.dir == "" {
		return nil
	}
	s.closeOnce.Do(func() { s.closeErr = os.RemoveAll(s.dir) })
	return s.closeErr
}

// WithEnv replaces every ambient selection with path.
func WithEnv(base []string, path string) []string {
	env := make([]string, 0, len(base)+1)
	for _, item := range base {
		if !strings.HasPrefix(item, Env+"=") {
			env = append(env, item)
		}
	}
	return append(env, Env+"="+path)
}

func canonicalSourceRoot(root string) (string, error) {
	resolved, err := canonicalpath.Resolve(root)
	if err != nil {
		return "", fmt.Errorf("resolve Bench source root %q: %w", root, err)
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("Bench source root %q is not a directory", resolved)
	}
	return filepath.Clean(resolved), nil
}

// Build publishes one Bench executable from sourceRoot at output through the sanctioned
// build script, the same build the run owner uses for its private selection. It is the
// one Go seam onto that script.
func Build(ctx context.Context, sourceRoot, output string) error {
	return canonicalBuild(ctx, sourceRoot, output)
}

func canonicalBuild(ctx context.Context, sourceRoot, output string) error {
	if _, err := exec.LookPath("go"); err != nil {
		return errors.New("Go is absent from PATH; prepend an executable Go toolchain directory to PATH and retry")
	}
	cmd := exec.Command("bash", filepath.Join(sourceRoot, "scripts", "go-build.sh"), sourceRoot, output)
	cmd.Dir = sourceRoot
	buildEnv, err := buildEnvironment(os.Environ())
	if err != nil {
		return err
	}
	cmd.Env = buildEnv
	cmd.SysProcAttr = builderProcAttr()
	var outputBytes bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outputBytes, &outputBytes
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		drainBuilderGroup(cmd.Process.Pid)
		if err != nil {
			return fmt.Errorf("%w: %s", err, strings.TrimSpace(outputBytes.String()))
		}
		return nil
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		select {
		case <-done:
			drainBuilderGroup(cmd.Process.Pid)
		case <-time.After(BuilderCancelGrace):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
			drainBuilderGroup(cmd.Process.Pid)
		}
		return context.Cause(ctx)
	}
}

func drainBuilderGroup(pgid int) {
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	for {
		if err := syscall.Kill(-pgid, 0); errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// buildEnvironment returns the environment the build script runs under. The cross
// compilation names and the selection name go, and the Bench build cache entry
// arrives, so the private build shares the archives the gate's phases read.
func buildEnvironment(base []string) ([]string, error) {
	env := make([]string, 0, len(base))
	for _, item := range base {
		if strings.HasPrefix(item, "GOOS=") || strings.HasPrefix(item, "GOARCH=") || strings.HasPrefix(item, Env+"=") {
			continue
		}
		env = append(env, item)
	}
	return gocache.Apply(env)
}
