// Package runbinary owns the Bench executable selected for one top-level run.
package runbinary

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/subprocess"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// Env carries the exact Bench executable selected by the current top-level run.
const Env = "BENCH_RUN_BINARY"

// BuilderCancelGrace is how long a cancelled builder group has to exit on SIGTERM
// before it is killed. This constant is exported, so a test waiting out the drain
// derives its own deadline from this window, instead of guessing a literal.
const BuilderCancelGrace = bounds.BuilderCancelGrace

type Builder func(context.Context, string, string) error
type Verifier func(string, string) error

// Factory supplies the authoring seams used by tests. Zero values select the private
// operating-system temp directory, canonical builder, and freshness verifier. RepairRoot
// is the checkout a refusal tells the operator to rebuild in; its zero value names the
// source root the selection was graded against. They part only where the digest root is
// a composed tree the operator cannot rebuild in. Protection is the protocol of the
// selection's private generation, and its zero value selects the installed protocol. A
// process with no installed protocol opens no generation.
type Factory struct {
	TempRoot   string
	RepairRoot string
	Build      Builder
	Verify     Verifier
	Protection *protection.Protocol
}

// Selection is the exact source-bound executable used by one run. Close removes bytes
// only for an owner-authored selection. An inherited consumer never owns its lifetime.
type Selection struct {
	Path       string
	SourceRoot string

	dir        string
	protocol   *protection.Protocol
	generation protection.Generation
	closeOnce  sync.Once
	closeErr   error
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
	raw, present := os.LookupEnv(Env)
	return f.Select(ctx, sourceRoot, raw, present)
}

// Select applies the inherited-selection rule to an Env value that the caller read: a
// present raw value is inherited or refused, and an absent one makes the factory build.
func (f Factory) Select(ctx context.Context, sourceRoot, raw string, present bool) (*Selection, error) {
	if present {
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

// Own creates one private selection through the factory's builder. Its private generation
// opens before the build. A failed selection returns its own Close error too.
func (f Factory) Own(ctx context.Context, sourceRoot string) (_ *Selection, err error) {
	source, err := canonicalSourceRoot(sourceRoot)
	if err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp(f.TempRoot, runDirPrefix)
	if err != nil {
		return nil, fmt.Errorf("create private Bench run directory: %w", err)
	}
	selection := &Selection{Path: filepath.Join(dir, "bench"), SourceRoot: source, dir: dir, protocol: f.Protection}
	if selection.protocol == nil {
		selection.protocol = protection.Installed()
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, selection.Close())
		}
	}()
	if selection.protocol != nil {
		if selection.generation, err = selectionDirectory.Open(selection.protocol, dir); err != nil {
			return nil, fmt.Errorf("protect private Bench run directory: %w", err)
		}
	}
	build := f.Build
	if build == nil {
		build = selection.canonicalBuild
	}
	if err := build(ctx, source, selection.Path); err != nil {
		return nil, fmt.Errorf("build private Bench executable: %w", err)
	}
	if err := f.Validate(source, selection.Path); err != nil {
		return nil, err
	}
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
	if err := f.Validate(source, value); err != nil {
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

// Validate accepts a cleaned absolute regular executable that the factory's verifier matches
// to sourceRoot. Own and Inherit apply it once, and a holder of a selection applies it before
// each reuse.
func (f Factory) Validate(sourceRoot, executable string) error {
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

// Owned reports whether the factory built this selection. An inherited selection reads false.
func (s *Selection) Owned() bool {
	return s != nil && s.dir != ""
}

// Resources names the selection's private generation for a launch that uses the executable.
// An inherited or unprotected selection names none.
func (s *Selection) Resources() []protection.Descriptor {
	if s == nil {
		return nil
	}
	return s.generation.Resources()
}

// Close releases an owner-authored private directory after its children have returned. A
// protected generation retains the directory, and Close returns that unresolved release.
func (s *Selection) Close() error {
	if !s.Owned() {
		return nil
	}
	s.closeOnce.Do(func() { s.closeErr = s.release() })
	return s.closeErr
}

// release removes the directory through the generation's safe removal.
func (s *Selection) release() error {
	return s.generation.Remove("selection "+s.dir, nil, func() error { return os.RemoveAll(s.dir) })
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
	return canonicalBuild(ctx, sourceRoot, output, subprocess.Registration{})
}

// BuildSubject publishes one Bench executable from sourceRoot at output under the build
// script's own manifest directory. It is the form a caller takes when output is the
// checkout's own published executable, so the manifest lands where that checkout's
// wrapper reads it rather than beside a throwaway.
func BuildSubject(ctx context.Context, sourceRoot, output string) error {
	return runBuildScript(ctx, sourceRoot, output, "", subprocess.Registration{})
}

// canonicalBuild is the selection's default builder. Its launch registers on the selection's
// private generation.
func (s *Selection) canonicalBuild(ctx context.Context, sourceRoot, output string) error {
	return canonicalBuild(ctx, sourceRoot, output, subprocess.Registration{Protocol: s.protocol, Resources: s.Resources()})
}

func canonicalBuild(ctx context.Context, sourceRoot, output string, local subprocess.Registration) error {
	// The private selection is deleted when the run owner closes it, so its manifest goes
	// beside it rather than beside the source checkout's wrapper. A manifest left in the
	// wrapper's directory would outlive this executable and refuse the next landing.
	return runBuildScript(ctx, sourceRoot, output, filepath.Dir(output), local)
}

// BuildError is a build script that ran and exited with a failure. A cancelled build and an
// absent Go are other errors. Output is everything the script printed, so a reader can name
// the compile error.
type BuildError struct {
	Err    error
	Output string
}

func (e BuildError) Error() string {
	return fmt.Sprintf("%v: %s", e.Err, strings.TrimSpace(e.Output))
}

func (e BuildError) Unwrap() error { return e.Err }

// builderPolicy kills and observes the remaining group after a normal exit. A cancellation
// sends TERM, and KILL follows the builder grace.
var builderPolicy = subprocess.Policy{
	Normal: subprocess.Drain,
	Grace:  bounds.ProcessWindow(bounds.FixedWindow(BuilderCancelGrace)),
	Final:  bounds.ProcessWindow(bounds.FixedWindow(bounds.ProcessFinalWindow)),
	Poll:   bounds.ProcessPoll,
}

// builderCancelRequest answers a deadline for an expired context and an interruption
// otherwise, each with TERM.
func builderCancelRequest(ctx context.Context) (subprocess.Cause, syscall.Signal) {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return subprocess.CauseDeadline, syscall.SIGTERM
	}
	return subprocess.CauseInterruption, syscall.SIGTERM
}

// runBuildScript runs the sanctioned build script once. An empty manifestDir passes no
// manifest-directory flag, which leaves the script its own default; a named one overrides
// it. The callers differ only in that operand and their local resources, so they share
// this one exec. The launch registers on the inherited union plus local through RunInherited.
func runBuildScript(ctx context.Context, sourceRoot, output, manifestDir string, local subprocess.Registration) error {
	if _, err := exec.LookPath("go"); err != nil {
		return errors.New("Go is absent from PATH; prepend an executable Go toolchain directory to PATH and retry")
	}
	script := []string{filepath.Join(sourceRoot, "scripts", "go-build.sh")}
	if manifestDir != "" {
		script = append(script, "--manifest-dir", manifestDir)
	}
	cmd := exec.Command("bash", append(script, sourceRoot, output)...)
	cmd.Dir = sourceRoot
	buildEnv, err := buildEnvironment(os.Environ())
	if err != nil {
		return err
	}
	cmd.Env = buildEnv
	cmd.SysProcAttr = builderProcAttr()
	streams := subprocess.Streams{Stdout: subprocess.Output{Capture: true}, MergeStderr: true}
	outcome, err := subprocess.RunInherited(ctx, cmd, builderPolicy, streams, builderCancelRequest, local)
	if err != nil {
		return err
	}
	switch {
	case outcome.Cancellation.Cause != subprocess.CauseNone:
		return context.Cause(ctx)
	case outcome.Process.StartErr != nil:
		return outcome.Process.StartErr
	case outcome.Process.WaitErr != nil:
		return BuildError{Err: outcome.Process.WaitErr, Output: string(outcome.Streams.Stdout.Data)}
	}
	if err := outcome.Uncertified(); err != nil {
		return fmt.Errorf("builder %w", err)
	}
	return nil
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
