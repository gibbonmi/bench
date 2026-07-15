package surface

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func assertPackedSetupForwarding(t *testing.T, dir, wrapper, shim string, env map[string]string) {
	t.Helper()
	real := wrapper + ".real"
	if err := os.Rename(wrapper, real); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Rename(real, wrapper) })
	stub := "#!/bin/sh\nprintf 'installed-setup:%s:%s\\n' \"$1\" \"$2\"\nexit 23\n"
	if err := os.WriteFile(wrapper, []byte(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	probe := exec.Command(shim, "setup", "a b")
	probe.Dir, probe.Env = dir, lifecycleEnv(env)
	out, err := probe.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 23 || string(out) != "installed-setup:setup:a b\n" {
		t.Fatalf("stable shim setup forwarding = %v, %q", err, out)
	}
	if err := os.Remove(wrapper); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(real, wrapper); err != nil {
		t.Fatal(err)
	}
}

func assertPackedEntrySurfaceIdentity(t *testing.T, repo string, env map[string]string, versionOut string) {
	t.Helper()
	hookOut := runLifecycle(t, repo, env, "bash", filepath.Join(repo, ".bench", "hooks", "session-start.sh"))
	if !strings.Contains(hookOut, filepath.Join(repo, ".bench", "bin", "bench.sh")) {
		t.Fatalf("session hook did not identify linked launcher:\n%s", hookOut)
	}
	stop := exec.Command("bash", filepath.Join(repo, ".bench", "hooks", "stop.sh"))
	stop.Dir, stop.Env, stop.Stdin = repo, lifecycleEnv(env), strings.NewReader("{}\n")
	if out, err := stop.CombinedOutput(); err != nil || strings.Contains(string(out), "GLOBAL-RUNTIME") {
		t.Fatalf("stop hook escaped linked launcher: %v\n%s", err, out)
	}
	stubDir := filepath.Join(repo, "provider stubs [*]")
	if err := os.MkdirAll(stubDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, provider := range []string{"claude", "codex", "opencode"} {
		body := "#!/bin/sh\nprintf 'provider:" + provider + ":%s\\n' \"$*\"\n"
		if err := os.WriteFile(filepath.Join(stubDir, provider), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
		adapterEnv := cloneEnv(env)
		adapterEnv["PATH"] = stubDir + string(os.PathListSeparator) + os.Getenv("PATH")
		out := runLifecycle(t, repo, adapterEnv, "bash", filepath.Join(repo, ".bench", "adapters", provider), "identity prompt")
		if !strings.Contains(out, "provider:"+provider+":") {
			t.Fatalf("packed %s adapter escaped linked launcher/provider: %s", provider, out)
		}
	}
	if !strings.Contains(versionOut, runtime.GOOS+"/") {
		t.Fatalf("selected target identity missing from installed version: %s", versionOut)
	}
}

func runPackedFreshClone(t *testing.T, repo, _, shim, version string) {
	t.Helper()
	runLifecycle(t, repo, nil, "git", "add", "-A")
	runLifecycle(t, repo, nil, "git", "commit", "-qm", "linked state")
	clone := filepath.Join(t.TempDir(), "committed fresh clone [*]")
	runLifecycle(t, filepath.Dir(clone), nil, "git", "clone", "-q", repo, clone)
	home := filepath.Join(t.TempDir(), "empty cache [*]")
	bin := filepath.Join(home, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{"HOME": home, "BENCH_HOME": home, "PATH": bin + string(os.PathListSeparator) + "/usr/bin:/bin"}
	runLifecycle(t, clone, env, shim, "link")
	runLifecycle(t, clone, env, shim, "init")
	manifest, err := os.ReadFile(filepath.Join(clone, ".bench", "link-manifest.tsv"))
	if err != nil || !bytes.Contains(manifest, []byte("#kit\t"+version+"\n")) {
		t.Fatalf("fresh clone maintenance lost installed identity: %q, %v", manifest, err)
	}
	if entries, err := os.ReadDir(filepath.Join(home, "cache")); err == nil && len(entries) != 0 {
		t.Fatalf("maintenance unexpectedly populated empty runtime cache: %v", entries)
	}
}

func runLifecycle(t *testing.T, dir string, overrides map[string]string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir, cmd.Env = dir, lifecycleEnv(overrides)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", name, args, err, out)
	}
	return string(out)
}

func lifecycleEnv(overrides map[string]string) []string {
	env := map[string]string{}
	for _, pair := range os.Environ() {
		if key, value, ok := strings.Cut(pair, "="); ok {
			env[key] = value
		}
	}
	for key, value := range overrides {
		env[key] = value
	}
	out := make([]string, 0, len(env))
	for key, value := range env {
		out = append(out, key+"="+value)
	}
	return out
}

func cloneEnv(env map[string]string) map[string]string {
	clone := map[string]string{}
	for key, value := range env {
		clone[key] = value
	}
	return clone
}

func committedHostileArtifactSource(t *testing.T, root string) string {
	t.Helper()
	origin := filepath.Join(t.TempDir(), "committed origin [*]")
	out, err := exec.Command("git", "-C", root, "ls-files", "-z", "--cached", "--others", "--exclude-standard").Output()
	if err != nil {
		t.Fatal(err)
	}
	rels := strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
	rels = append(rels, "projects/gl-axi.md", "projects/regroup.md")
	for _, rel := range rels {
		src, dst := filepath.Join(root, rel), filepath.Join(origin, rel)
		info, err := os.Lstat(src)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			t.Fatal(err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(src)
			if err == nil {
				err = os.Symlink(target, dst)
			}
			if err != nil {
				t.Fatal(err)
			}
			continue
		}
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, data, info.Mode().Perm()); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "bench@local"}, {"config", "user.name", "bench"}, {"add", "-f", "."}, {"commit", "-qm", "artifact source"}} {
		cmd := exec.Command("git", append([]string{"-C", origin}, args...)...)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, output)
		}
	}
	clone := filepath.Join(t.TempDir(), "fresh source clone [*]")
	if output, err := exec.Command("git", "clone", "-q", origin, clone).CombinedOutput(); err != nil {
		t.Fatalf("clone committed source: %v\n%s", err, output)
	}
	return clone
}

func assertInterruptedArtifactPromotion(t *testing.T, source, output string, wantFiles int) {
	t.Helper()
	ready := filepath.Join(t.TempDir(), "promotion-ready")
	cmd := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
	cmd.Env = append(os.Environ(), "BENCH_TEST_PROMOTION_READY_FILE="+ready)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(90 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			t.Fatal("builder did not reach deterministic promotion seam")
		}
		time.Sleep(25 * time.Millisecond)
	}
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err == nil {
		t.Fatal("interrupted builder exited successfully")
	}
	files, err := os.ReadDir(output)
	if err != nil || len(files) != wantFiles {
		t.Fatalf("promotion interruption left partial/absent set: files=%d err=%v", len(files), err)
	}
}
