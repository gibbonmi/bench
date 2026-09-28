package env

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

func kitRunBase(t *testing.T) []string {
	t.Helper()
	home := t.TempDir()
	tmp := filepath.Join(t.TempDir(), "temporary space")
	if err := os.Mkdir(tmp, 0o700); err != nil {
		t.Fatal(err)
	}
	global := filepath.Join(home, "global")
	system := filepath.Join(home, "system")
	goenv := filepath.Join(home, "goenv")
	for path, contents := range map[string]string{
		global: "[probe]\nmarker = visible\n",
		system: "[probe]\nsystem = visible\n",
		goenv:  "GOPATH=" + filepath.Join(home, "operator-go") + "\nGOMODCACHE=" + filepath.Join(home, "operator-modules") + "\n",
	} {
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	base := []string{}
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "GIT_CONFIG_") && !assignsAny(entry, []string{"HOME", "TMPDIR", "GOPATH", "GOMODCACHE", "GOENV", "XDG_CONFIG_HOME"}) {
			base = append(base, entry)
		}
	}
	return append(base, "HOME="+home, "TMPDIR="+tmp, "GOENV="+goenv,
		"GIT_CONFIG_GLOBAL="+global, "GIT_CONFIG_SYSTEM="+system, "GIT_CONFIG_NOSYSTEM=0",
		"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=maintenance.auto", "GIT_CONFIG_VALUE_0=true")
}

func TestKitTestRunIsolationAndGoSettings(t *testing.T) {
	for _, configured := range []bool{false, true} {
		name := "default Go settings"
		if configured {
			name = "configured Go settings"
		}
		t.Run(name, func(t *testing.T) { checkKitRunIsolation(t, configured) })
	}
}

func checkKitRunIsolation(t *testing.T, configured bool) {
	t.Helper()
	base := kitRunBase(t)
	if !configured {
		base = append(base, "GOENV=")
	}
	for _, args := range [][]string{{"config", "--global", "--get", "probe.marker"}, {"config", "--get", "probe.system"}} {
		cmd := exec.Command("git", args...)
		cmd.Env = base
		if out, err := cmd.CombinedOutput(); err != nil || string(out) != "visible\n" {
			t.Fatalf("operator marker = %q, %v", out, err)
		}
	}
	run, err := OpenKitTestRun(base)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := run.Close(); err != nil {
			t.Error(err)
		}
	})
	entries := run.Entries()
	for _, name := range []string{"HOME", "TMPDIR"} {
		dir := kitEnvValue(entries, name)
		if filepath.Dir(dir) != run.dir || filepath.Dir(run.dir) != kitEnvValue(base, "TMPDIR") {
			t.Fatalf("%s=%q is outside the run under TMPDIR", name, dir)
		}
	}
	if growth := len(kitEnvValue(entries, "TMPDIR")) - len(kitEnvValue(base, "TMPDIR")); growth > 16 {
		t.Fatalf("private TMPDIR grew by %d bytes, want at most 16", growth)
	}
	if kitEnvValue(entries, "GOCACHE") != "" {
		t.Fatal("run overrides the caller's build cache")
	}
	probe := gittest.KitRunProbe(t, base)
	cmd := exec.Command(probe[0], probe[1:]...)
	cmd.Env = append(base, entries...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("kit run probe: %v: %s", err, out)
	}
	entries[0] = "HOME=changed"
	if run.Entries()[0] == entries[0] {
		t.Fatal("caller changed the run's entries")
	}
}

func TestKitTestRunCloseRestoresDirectoryAccess(t *testing.T) {
	run, err := OpenKitTestRun(kitRunBase(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = run.Close() })
	locked := filepath.Join(kitEnvValue(run.Entries(), "TMPDIR"), "locked")
	if err := os.Mkdir(locked, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(locked, "file"), []byte("held"), 0o600); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(locked, "outside")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(locked, 0); err != nil {
		t.Fatal(err)
	}
	if err := run.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(run.dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("run directory remains after close: %v", err)
	}
	if _, err := os.Stat(outside); err != nil {
		t.Fatalf("close removed a symlink target: %v", err)
	}
}

func TestKitTestRunIgnoresInheritedBashStartup(t *testing.T) {
	const marker = "inherited-startup"
	startup := filepath.Join(t.TempDir(), "startup")
	if err := os.WriteFile(startup, []byte("printf '"+marker+"\\n'\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	base := append(kitRunBase(t), "BASH_ENV="+startup)
	bashOutput := func(environment []string) string {
		t.Helper()
		cmd := exec.Command("bash", "-c", "printf child")
		cmd.Env = environment
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Bash child: %v: %s", err, out)
		}
		return string(out)
	}
	if got := bashOutput(base); got != marker+"\nchild" {
		t.Fatalf("base startup control = %q", got)
	}
	run, err := OpenKitTestRun(base)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := run.Close(); err != nil {
			t.Error(err)
		}
	})
	if got := bashOutput(append(base, run.Entries()...)); got != "child" {
		t.Fatalf("kit child sourced the inherited startup file: %q", got)
	}
}

func TestKitTestRunRejectsInvalidHomeBeforeCreation(t *testing.T) {
	for _, home := range []string{"absent", "", "relative"} {
		t.Run(home, func(t *testing.T) {
			t.Chdir(t.TempDir())
			tmp := t.TempDir()
			base := []string{"TMPDIR=" + tmp}
			if home != "absent" {
				base = append(base, "HOME="+home)
			}
			if _, err := OpenKitTestRun(base); err == nil || !strings.Contains(err.Error(), "HOME") {
				t.Fatalf("invalid HOME: %v", err)
			}
			entries, err := os.ReadDir(tmp)
			if err != nil || len(entries) != 0 {
				t.Fatalf("invalid HOME created a run: %v, %v", entries, err)
			}
		})
	}
}

func TestKitTestRunRejectsInvalidTmpdir(t *testing.T) {
	base := kitRunBase(t)
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	base = append(base, "TMPDIR="+file)
	if _, err := OpenKitTestRun(base); err == nil || !strings.Contains(err.Error(), "TMPDIR") {
		t.Fatalf("invalid TMPDIR: %v", err)
	}
}
