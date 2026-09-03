package runbinary

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gocache"
)

func TestFactoryOwnBuildsOnePrivateAbsoluteSelectionAndCleansIt(t *testing.T) {
	tempRoot := t.TempDir()
	source := t.TempDir()
	builds := 0
	factory := Factory{
		TempRoot: tempRoot,
		Build: func(_ context.Context, gotSource, output string) error {
			builds++
			if gotSource != source {
				t.Fatalf("builder source = %q, want %q", gotSource, source)
			}
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(gotSource, executable string) error {
			if gotSource != source {
				t.Fatalf("verifier source = %q, want %q", gotSource, source)
			}
			info, err := os.Lstat(executable)
			if err != nil || !info.Mode().IsRegular() {
				t.Fatalf("selected executable = %v, %v", info, err)
			}
			return nil
		},
	}

	t.Setenv(Env, filepath.Join(tempRoot, "hostile"))
	selection, err := factory.Own(context.Background(), source)
	if err != nil {
		t.Fatalf("Own: %v", err)
	}
	if builds != 1 {
		t.Fatalf("builds = %d, want 1", builds)
	}
	if !filepath.IsAbs(selection.Path) || filepath.Clean(selection.Path) != selection.Path {
		t.Fatalf("selection path = %q, want cleaned absolute path", selection.Path)
	}
	if selection.Path == os.Getenv(Env) || !strings.HasPrefix(filepath.Dir(selection.Path), tempRoot+string(os.PathSeparator)) {
		t.Fatalf("selection path = %q, want private path under %q and not ambient", selection.Path, tempRoot)
	}
	parent := filepath.Dir(selection.Path)
	if err := selection.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := os.Stat(parent); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("run directory after Close = %v, want absent", err)
	}
}

func TestFactoryOwnRemovesPartialOutputOnBuilderFailure(t *testing.T) {
	tempRoot := t.TempDir()
	factory := Factory{
		TempRoot: tempRoot,
		Build: func(_ context.Context, _, output string) error {
			if err := os.WriteFile(output, []byte("partial"), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(output+".seal", []byte("partial"), 0o600); err != nil {
				return err
			}
			return errors.New("builder red")
		},
		Verify: func(string, string) error { return nil },
	}
	if _, err := factory.Own(context.Background(), t.TempDir()); err == nil {
		t.Fatal("Own builder failure = nil")
	}
	entries, err := os.ReadDir(tempRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("partial run entries = %v, want none", entries)
	}
}

// parkingBuilder is one source tree whose build script publishes the selection and then
// parks a descendant. The descendant traps INT and TERM, so only a group SIGKILL removes
// it, and it records whether the selection went away before that kill arrived. The fields
// name the paths a test reads back.
type parkingBuilder struct {
	source     string
	builderPID string
	childPID   string
	premature  string
}

// newParkingBuilder writes that script under a fresh source root. When hold is true the
// builder child waits on its descendant, so the builder child is still alive when its
// owner dies. The drain test needs the opposite: a builder that returns while the
// descendant stays parked.
func newParkingBuilder(t *testing.T, hold bool) parkingBuilder {
	t.Helper()
	source := t.TempDir()
	builder := parkingBuilder{
		source:     source,
		builderPID: filepath.Join(source, "builder-self"),
		childPID:   filepath.Join(source, "builder-child"),
		premature:  filepath.Join(source, "selection-removed-first"),
	}
	if err := os.MkdirAll(filepath.Join(source, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	script := `#!/usr/bin/env bash
output="$2"
printf selected > "$output"
chmod 755 "$output"
echo $$ > "` + builder.builderPID + `"
(
  trap '' INT TERM
  exec >/dev/null 2>&1
  while test -e "$output"; do sleep .01; done
  printf premature > "` + builder.premature + `"
) &
echo $! > "` + builder.childPID + `"
`
	if hold {
		script += "wait\n"
	}
	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return builder
}

func TestCanonicalBuilderDrainsDescendantsBeforeReturningSelection(t *testing.T) {
	builder := newParkingBuilder(t, false)
	factory := Factory{TempRoot: t.TempDir(), Verify: func(string, string) error { return nil }}
	selection, err := factory.Own(context.Background(), builder.source)
	if err != nil {
		t.Fatal(err)
	}
	child := readPID(t, builder.childPID)
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	requireProcessExit(t, child)
	if err := selection.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(builder.premature); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("builder descendant observed premature cleanup: %v", err)
	}
}

// TestFactoryValidateRefusesAStaleInheritedSeal is BF7. An inherited executable carries
// an intact seal pair of its own, so the pair check alone accepts a binary its sources
// have moved past. A named source root grades the seal's source digest against the tree
// the caller named, and only the caller that names no root keeps the narrower check.
func TestFactoryValidateRefusesAStaleInheritedSeal(t *testing.T) {
	source := writeInheritedSourceFixture(t)
	dir := t.TempDir()
	staged := filepath.Join(dir, "staged-bench")
	if err := os.WriteFile(staged, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(dir, "bench")
	if err := freshness.Publish(source, staged, executable, "1.2.3"); err != nil {
		t.Fatal(err)
	}
	factory := Factory{}
	if _, err := factory.Inherit(source, executable); err != nil {
		t.Fatalf("freshly sealed inherited selection: %v", err)
	}

	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte("#!/usr/bin/env bash\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := factory.Inherit(source, executable)
	if err == nil {
		t.Fatal("inherited selection with a stale source digest = nil, want refusal")
	}
	if !strings.Contains(err.Error(), freshness.RebuildAction(source)) {
		t.Fatalf("stale inherited refusal = %q, want the rebuild action %q", err, freshness.RebuildAction(source))
	}
	if _, err := factory.Inherit("", executable); err != nil {
		t.Fatalf("inherited selection with no named source root: %v", err)
	}
}

// writeInheritedSourceFixture writes the smallest tree freshness.Digest resolves: one
// command package, the module file, and the auxiliary manifest with the build script it
// names.
func writeInheritedSourceFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range map[string]string{
		"go.mod":                  "module example.com/runbinaryfixture\n\ngo 1.25\n",
		"cmd/bench/main.go":       "package main\n\nfunc main() {}\n",
		"scripts/go-build.sh":     "#!/usr/bin/env bash\nexit 0\n",
		"scripts/go-build.inputs": "build_script=scripts/go-build.sh\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestFactoryInheritRequiresCleanAbsoluteRegularExecutable(t *testing.T) {
	source := t.TempDir()
	dir := t.TempDir()
	executable := filepath.Join(dir, "bench")
	if err := os.WriteFile(executable, []byte("selected"), 0o755); err != nil {
		t.Fatal(err)
	}
	factory := Factory{Verify: func(gotSource, gotExecutable string) error {
		if gotSource != source || gotExecutable != executable {
			t.Fatalf("Verify(%q, %q), want (%q, %q)", gotSource, gotExecutable, source, executable)
		}
		return nil
	}}
	selection, err := factory.Inherit(source, executable)
	if err != nil {
		t.Fatalf("Inherit valid selection: %v", err)
	}
	if selection.Path != executable {
		t.Fatalf("selection path = %q, want %q", selection.Path, executable)
	}

	for name, value := range map[string]string{
		"missing":  "",
		"relative": "bench",
		"unclean":  dir + string(os.PathSeparator) + "sub" + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "bench",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := factory.Inherit(source, value); err == nil {
				t.Fatalf("Inherit(%q) = nil, want refusal", value)
			}
		})
	}
	symlink := filepath.Join(dir, "link")
	if err := os.Symlink(executable, symlink); err != nil {
		t.Fatal(err)
	}
	if _, err := factory.Inherit(source, symlink); err == nil {
		t.Fatal("Inherit symlink = nil, want refusal")
	}
}

func TestWithEnvReplacesEveryAmbientSelection(t *testing.T) {
	got := WithEnv([]string{"A=1", Env + "=old", Env + "=older"}, "/private/bench")
	want := []string{"A=1", Env + "=/private/bench"}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("WithEnv = %q, want %q", got, want)
	}
}

func readPID(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var pid int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &pid); err != nil || pid <= 0 {
		t.Fatalf("read pid %q: %v", data, err)
	}
	return pid
}

// requireProcessExit waits out the grace a cancelled builder group has to exit, so the
// window derives from that grace rather than from a literal a loaded machine can beat.
func requireProcessExit(t *testing.T, pid int) {
	t.Helper()
	window := bounds.TestDeadline(BuilderCancelGrace)
	deadline := time.Now().Add(window)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal(bounds.TestTimeoutVerdict(fmt.Sprintf("process %d to exit after the builder returned", pid), window))
}

// C09: the private build's environment carries the Bench build cache entry, so the
// builder writes to the one Bench-owned directory instead of the ambient one.
func TestBuildEnvironmentCarriesTheBenchBuildCache(t *testing.T) {
	t.Parallel()
	built, err := buildEnvironment([]string{"HOME=/home/agent", "GOCACHE=/ambient/cache", "GOOS=plan9"})
	if err != nil {
		t.Fatal(err)
	}
	want, err := gocache.Dir([]string{"HOME=/home/agent"})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(built, []string{"HOME=/home/agent", gocache.Env + "=" + want}) {
		t.Fatalf("buildEnvironment = %#v, want HOME and %s=%s only", built, gocache.Env, want)
	}
}

func TestBuildEnvironmentRefusesWithoutAnAbsoluteHome(t *testing.T) {
	t.Parallel()
	built, err := buildEnvironment([]string{"PATH=/usr/bin"})
	if err == nil {
		t.Fatalf("buildEnvironment = %#v, want an error", built)
	}
	if !strings.Contains(err.Error(), "HOME") {
		t.Fatalf("error = %q, want it to name HOME", err)
	}
}
