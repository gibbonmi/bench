package runbinary

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
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

func TestCanonicalBuilderDrainsDescendantsBeforeReturningSelection(t *testing.T) {
	source := t.TempDir()
	pids := filepath.Join(source, "builder-child")
	premature := filepath.Join(source, "selection-removed-first")
	if err := os.MkdirAll(filepath.Join(source, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	script := `#!/usr/bin/env bash
output="$2"
printf selected > "$output"
chmod 755 "$output"
(
  trap '' INT TERM
  exec >/dev/null 2>&1
  while test -e "$output"; do sleep .01; done
  printf premature > "` + premature + `"
) &
echo $! > "` + pids + `"
`
	if err := os.WriteFile(filepath.Join(source, "scripts", "go-build.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	factory := Factory{TempRoot: t.TempDir(), Verify: func(string, string) error { return nil }}
	selection, err := factory.Own(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	child := readPID(t, pids)
	t.Cleanup(func() { _ = syscall.Kill(child, syscall.SIGKILL) })
	requireProcessExit(t, child)
	if err := selection.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(premature); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("builder descendant observed premature cleanup: %v", err)
	}
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

func requireProcessExit(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if err := syscall.Kill(pid, 0); errors.Is(err, syscall.ESRCH) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("process %d survived builder return", pid)
}
