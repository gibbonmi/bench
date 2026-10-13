package testreport

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/runbinary"
)

func TestHeldRunsOwnAFreshBuildWhenStale(t *testing.T) {
	tree := newSealedTree(t)
	held := tree.hold(t, context.Background(), tree.factory())
	runHeld(t, held, OutcomePassed)
	tree.write(t, "v2")
	runHeld(t, held, OutcomePassed)
	if want := []string{"v1", "v2"}; !slices.Equal(tree.read, want) {
		t.Fatalf("builds read trees %q, want %q", tree.read, want)
	}
	if got := held.Builds(); got != len(tree.read) {
		t.Fatalf("Builds() = %d, want %d", got, len(tree.read))
	}
	if got := readTestReportFile(t, tree.marker); got != tree.outputs[1] {
		t.Fatalf("stale run selected %q, want the fresh build %q", got, tree.outputs[1])
	}
}

func TestHeldRunsReuseOnAFreshVerify(t *testing.T) {
	tree := newSealedTree(t)
	held := tree.hold(t, context.Background(), tree.factory())
	for range 3 {
		runHeld(t, held, OutcomePassed)
	}
	if want := []string{"v1"}; !slices.Equal(tree.read, want) {
		t.Fatalf("builds read trees %q, want %q", tree.read, want)
	}
	if got := held.Builds(); got != len(tree.read) {
		t.Fatalf("Builds() = %d, want %d", got, len(tree.read))
	}
	// Own verifies its output once, and each later run verifies the held path.
	path := tree.outputs[0]
	if want := []string{path, path, path}; !slices.Equal(tree.verified, want) {
		t.Fatalf("verified %q, want %q", tree.verified, want)
	}
	if got := readTestReportFile(t, tree.marker); got != path {
		t.Fatalf("third run selected %q, want the held build %q", got, path)
	}
}

func TestHeldRunsRefuseAStaleInheritedBinary(t *testing.T) {
	tree := newSealedTree(t)
	held := tree.hold(t, context.Background(), tree.factory())
	selected := tree.inherit(t)
	runHeld(t, held, OutcomePassed)
	tree.write(t, "v2")
	runHeld(t, held, OutcomeRefused)
	if len(tree.read) != 0 || held.Builds() != 0 {
		t.Fatalf("builds read trees %q and Builds() = %d, want no build", tree.read, held.Builds())
	}
	if want := []string{selected, selected}; !slices.Equal(tree.verified, want) {
		t.Fatalf("verified %q, want %q", tree.verified, want)
	}
}

func TestHeldRunsRemoveEachOwnedBuild(t *testing.T) {
	tree := newSealedTree(t)
	held := tree.hold(t, context.Background(), tree.factory())
	runHeld(t, held, OutcomePassed)
	tree.write(t, "v2")
	runHeld(t, held, OutcomePassed)
	if len(tree.outputs) != 2 {
		t.Fatalf("builds wrote %q, want a held and a fresh build", tree.outputs)
	}
	if got, want := tree.runDirs(t), []string{filepath.Dir(tree.outputs[0])}; !slices.Equal(got, want) {
		t.Fatalf("run directories after the stale run = %q, want only the held %q", got, want)
	}
	if err := held.Close(); err != nil {
		t.Fatal(err)
	}
	if got := tree.runDirs(t); len(got) != 0 {
		t.Fatalf("run directories after Close = %q, want none", got)
	}
}

func TestHeldRunsInterruptACancelledBuild(t *testing.T) {
	tree := newSealedTree(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	factory := tree.factory()
	seal, park := factory.Build, parkingBuild(parkingBuilderSource(t))
	parked := ""
	factory.Build = func(ctx context.Context, source, output string) error {
		if len(tree.read) == 0 {
			return seal(ctx, source, output)
		}
		parked = output
		return park(ctx, source, output)
	}
	pidFile := filepath.Join(t.TempDir(), "builder-descendant")
	t.Setenv(cancelGoPIDEnv, pidFile)
	t.Setenv(cancelParkingPath, os.Args[0])
	held := tree.hold(t, ctx, factory)
	runHeld(t, held, OutcomePassed)
	tree.write(t, "v2")

	outcomes := make(chan Outcome, 1)
	go func() {
		outcome, _, _ := held.Run()
		outcomes <- outcome
	}()
	descendant := waitForPID(t, pidFile)
	t.Cleanup(func() { _ = syscall.Kill(descendant, syscall.SIGKILL) })
	cancel()
	select {
	case outcome := <-outcomes:
		if outcome.Kind != OutcomeInterrupted {
			t.Fatalf("cancelled build = %q, want %q", outcome.Kind, OutcomeInterrupted)
		}
	case <-time.After(deadline()):
		t.Fatalf("cancelled held run did not return within %s", deadline())
	}
	requireProcessGone(t, descendant)
	if _, err := os.Stat(filepath.Dir(parked)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("parked build directory after cancel = %v, want removed", err)
	}
	if got := held.Builds(); got != len(tree.read) {
		t.Fatalf("Builds() = %d, want %d", got, len(tree.read))
	}
}

func TestHeldRunsInterruptACancelledValidation(t *testing.T) {
	tree := newSealedTree(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	factory := tree.factory()
	verify := factory.Verify
	factory.Verify = func(source, executable string) error {
		if len(tree.verified) == 1 {
			cancel()
		}
		return verify(source, executable)
	}
	held := tree.hold(t, ctx, factory)
	tree.inherit(t)
	runHeld(t, held, OutcomePassed)
	tree.write(t, "v2")
	runHeld(t, held, OutcomeInterrupted)
	if len(tree.read) != 0 {
		t.Fatalf("builds read trees %q, want no build", tree.read)
	}
}

// heldTreeFile is the one tree file that a sealed build reads.
const heldTreeFile = "tree"

// sealedTree is a run-binary factory over one tree file. Its builder seals the tree bytes
// into its output, and its verifier refuses a seal that differs from the current tree.
type sealedTree struct {
	root     string
	marker   string
	tempRoot string
	read     []string
	outputs  []string
	verified []string
}

func newSealedTree(t *testing.T) *sealedTree {
	t.Helper()
	root, marker := selectedPathTestModule(t)
	tree := &sealedTree{root: root, marker: marker, tempRoot: t.TempDir()}
	tree.write(t, "v1")
	return tree
}

func (s *sealedTree) write(t *testing.T, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(s.root, heldTreeFile), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (s *sealedTree) factory() runbinary.Factory {
	return runbinary.Factory{
		TempRoot: s.tempRoot,
		Build: func(_ context.Context, source, output string) error {
			tree, err := os.ReadFile(filepath.Join(source, heldTreeFile))
			if err != nil {
				return err
			}
			s.read = append(s.read, string(tree))
			s.outputs = append(s.outputs, output)
			return os.WriteFile(output, tree, 0o755)
		},
		Verify: func(source, executable string) error {
			s.verified = append(s.verified, executable)
			tree, err := os.ReadFile(filepath.Join(source, heldTreeFile))
			if err != nil {
				return err
			}
			seal, err := os.ReadFile(executable)
			if err != nil {
				return err
			}
			if !bytes.Equal(seal, tree) {
				return errors.New("the seal does not match the tree")
			}
			return nil
		},
	}
}

// hold installs factory and holds the default request of the tree's module under ctx.
func (s *sealedTree) hold(t *testing.T, ctx context.Context, factory runbinary.Factory) *HeldRuns {
	t.Helper()
	installTestSelectionFactory(t, factory)
	request, line, code := Prepare(s.root, nil)
	if line != "" {
		t.Fatalf("Prepare = (%q, %d), want an accepted selection", line, code)
	}
	held := Hold(ctx, s.root, request)
	t.Cleanup(func() { _ = held.Close() })
	return held
}

// inherit names an inherited run binary sealed from the current tree. It follows hold,
// because the factory install clears the inherited selection.
func (s *sealedTree) inherit(t *testing.T) string {
	t.Helper()
	tree, err := os.ReadFile(filepath.Join(s.root, heldTreeFile))
	if err != nil {
		t.Fatal(err)
	}
	selected := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(selected, tree, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, selected)
	return selected
}

// runDirs lists the factory's private temp root, which holds only run directories.
func (s *sealedTree) runDirs(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(s.tempRoot)
	if err != nil {
		t.Fatal(err)
	}
	dirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		dirs = append(dirs, filepath.Join(s.tempRoot, entry.Name()))
	}
	return dirs
}

func runHeld(t *testing.T, held *HeldRuns, want OutcomeKind) {
	t.Helper()
	if outcome, out, _ := held.Run(); outcome.Kind != want {
		t.Fatalf("held run = %q, want %q\n%s", outcome.Kind, want, out)
	}
}
