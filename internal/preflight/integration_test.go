package preflight

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReleasePreflightScriptBootstrapsBuiltFullAndFocusedCommands(t *testing.T) {
	source := projectRoot(t)
	root := t.TempDir()
	listed := exec.Command("git", "-C", source, "ls-files", "-z")
	out, err := listed.Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00") {
		from, to := filepath.Join(source, rel), filepath.Join(root, rel)
		info, err := os.Lstat(from)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			t.Fatal(err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(from)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, to); err != nil {
				t.Fatal(err)
			}
			continue
		}
		data, err := os.ReadFile(from)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(to, data, info.Mode().Perm()); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "test@example.invalid"}, {"config", "user.name", "Test"}, {"add", "."}, {"commit", "-qm", "clean checkout"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, output)
		}
	}
	fake := filepath.Join(root, "phase")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nprintf '{\"config\":{}}\\n'\n"), fs.FileMode(0o755)); err != nil {
		t.Fatal(err)
	}
	env := append([]string{}, os.Environ()...)
	for _, phase := range PhaseNames(ModeVerify) {
		env = append(env, "BENCH_PREFLIGHT_"+strings.ToUpper(phase)+"="+fake)
	}
	for _, args := range [][]string{{"--mode", "verify"}, {"--mode", "verify", "--phase", "smoke"}} {
		cmd := exec.Command("bash", append([]string{filepath.Join(root, "scripts", "release-preflight.sh")}, args...)...)
		cmd.Dir, cmd.Env = root, env
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("built preflight %v: %v\n%s", args, err, output)
		}
	}
	if info, err := os.Stat(filepath.Join(root, "dist", "bench-preflight")); err != nil || info.Mode()&0o111 == 0 {
		t.Fatalf("compiled command was not bootstrapped: %v", err)
	}
}

func TestBuiltCommandReleasePolicyFailuresAreRed(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(projectRoot(t), "scripts", "go-build.sh"), projectRoot(t), binary)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}

	t.Run("tag version mismatch", func(t *testing.T) {
		root := preflightRepo(t)
		tagRelease(t, root, false)
		if err := os.WriteFile(filepath.Join(root, "package.json"), []byte(`{"version":"0.2.1"}`), 0o644); err != nil {
			t.Fatal(err)
		}
		assertBuiltRed(t, binary, root, []string{"--mode", "publish", "--phase", "identity"}, "must agree")
	})
	t.Run("stranded changelog", func(t *testing.T) {
		root := preflightRepo(t)
		tagRelease(t, root, true)
		body := "# Changelog\n\n## [Unreleased]\n\n- stranded\n\n## [0.2.0] - 2026-07-16\n\n- release\n"
		if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		assertBuiltRed(t, binary, root, []string{"--mode", "publish"}, "stranded content")
	})
	t.Run("unrelated ancestry cannot be proven", func(t *testing.T) {
		root := preflightRepo(t)
		tagRelease(t, root, false)
		treeCmd := exec.Command("git", "rev-parse", "HEAD^{tree}")
		treeCmd.Dir = root
		tree, err := treeCmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		commitCmd := exec.Command("git", "commit-tree", strings.TrimSpace(string(tree)), "-m", "unrelated")
		commitCmd.Dir = root
		unrelated, err := commitCmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		refCmd := exec.Command("git", "update-ref", "refs/remotes/origin/main", strings.TrimSpace(string(unrelated)))
		refCmd.Dir = root
		if output, err := refCmd.CombinedOutput(); err != nil {
			t.Fatalf("origin: %v %s", err, output)
		}
		assertBuiltRed(t, binary, root, []string{"--mode", "publish"}, "ancestry")
	})
	for _, policy := range []struct{ name, body, want string }{
		{"malformed exceptions", `{`, "malformed"},
		{"expired exceptions", `[{"id":"GO-1","reason":"temporary","expires":"2026-07-15"}]`, "expired"},
	} {
		t.Run(policy.name, func(t *testing.T) {
			root := preflightRepo(t)
			scanner := filepath.Join(root, "scanner")
			if err := os.WriteFile(scanner, []byte("#!/bin/sh\nprintf '{\"finding\":{\"osv\":\"GO-1\",\"trace\":[{}]}}\\n'\n"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "scripts", "vuln-exceptions.json"), []byte(policy.body), 0o644); err != nil {
				t.Fatal(err)
			}
			t.Setenv("BENCH_PREFLIGHT_VULNERABILITY", scanner)
			t.Setenv("BENCH_PREFLIGHT_DATE", "2026-07-16")
			assertBuiltRed(t, binary, root, []string{"--mode", "verify", "--phase", "vulnerability"}, policy.want)
		})
	}
	t.Run("promotion fault preserves hostile target", func(t *testing.T) {
		root := preflightRepo(t)
		target := filepath.Join(root, "dist", "preflight")
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte("prior"), 0o644); err != nil {
			t.Fatal(err)
		}
		assertBuiltRed(t, binary, root, []string{"--mode", "verify", "--phase", "gate"}, "not a real directory")
		data, err := os.ReadFile(target)
		if err != nil || string(data) != "prior" {
			t.Fatalf("prior target changed: %q %v", data, err)
		}
	})
}

func assertBuiltRed(t *testing.T, binary, root string, args []string, want string) {
	t.Helper()
	cmd := exec.Command(binary, append([]string{"release-preflight"}, args...)...)
	cmd.Dir = root
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("command passed:\n%s", output)
	}
	if !strings.Contains(string(output), want) {
		t.Fatalf("output does not contain %q:\n%s", want, output)
	}
}

func tagRelease(t *testing.T, root string, withOrigin bool) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(releaseChangelog("2026-07-16")), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"tag", "v0.2.0"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, output)
		}
	}
	if withOrigin {
		cmd := exec.Command("git", "update-ref", "refs/remotes/origin/main", "HEAD")
		cmd.Dir = root
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("origin: %v %s", err, output)
		}
	}
	t.Setenv("BENCH_PREFLIGHT_REF", "refs/tags/v0.2.0")
}

func projectRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve test source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func preflightRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, rel := range []string{"bin/bench.sh", ".bench/gate.sh", "scripts/build-artifacts.sh", "scripts/smoke-artifacts.sh", "scripts/platforms.json", "scripts/wrapper-assets.json", "package.json"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		body := "{}\n"
		if filepath.Ext(path) == ".sh" || rel == "bin/bench.sh" {
			body = "#!/bin/sh\nexit 0\n"
		}
		if rel == "package.json" {
			body = `{"version":"0.2.0"}`
		}
		if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ngo 1.25\ntoolchain go1.25.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	fake := filepath.Join(root, "phase")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nprintf '{\"config\":{}}\\n'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range PhaseNames(ModePublish) {
		t.Setenv("BENCH_PREFLIGHT_"+strings.ToUpper(name), fake)
	}
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "test@example.invalid"}, {"config", "user.name", "Test"}, {"add", "."}, {"commit", "-qm", "fixture"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	return root
}

func intPtr(v int) *int { return &v }
