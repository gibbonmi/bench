package preflight

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestReleasePreflightScriptBootstrapsBuiltFullAndFocusedCommands(t *testing.T) {
	source := projectRoot(t)
	root := t.TempDir()
	listed := exec.Command("git", "-C", source, "ls-files", "--cached", "--others", "--exclude-standard", "-z")
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
		if phase == "artifacts" {
			continue
		}
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
		assertBuiltRed(t, binary, root, []string{"--mode", "publish", "--profile", "public"}, "must agree")
	})
	t.Run("stranded changelog", func(t *testing.T) {
		root := preflightRepo(t)
		tagRelease(t, root, true)
		body := "# Changelog\n\n## [Unreleased]\n\n- stranded\n\n## [0.2.0] - 2026-07-16\n\n- release\n"
		if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		assertBuiltRed(t, binary, root, []string{"--mode", "publish", "--profile", "public"}, "stranded content")
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
		assertBuiltRed(t, binary, root, []string{"--mode", "publish", "--profile", "public"}, "ancestry")
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

func TestBuiltCommandCancellationPreservesPriorCompleteEvidence(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("SIGINT process control is POSIX")
	}
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(projectRoot(t), "scripts", "go-build.sh"), projectRoot(t), binary)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	root := preflightRepo(t)
	first := exec.Command(binary, "release-preflight", "--mode", "verify")
	first.Dir = root
	if output, err := first.CombinedOutput(); err != nil {
		t.Fatalf("initial verify: %v\n%s", err, output)
	}
	indexPath := filepath.Join(root, "dist", "preflight", "release-index.json")
	sumsPath := filepath.Join(root, "dist", "preflight", "SHA256SUMS")
	oldIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	oldSums, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(t.TempDir(), "gate-ready")
	blocking := filepath.Join(root, "blocking-gate")
	if err := os.WriteFile(blocking, []byte("#!/bin/sh\nprintf ready > \"$BENCH_BLOCK_READY\"\nsleep 30\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_GATE", blocking)
	t.Setenv("BENCH_BLOCK_READY", ready)
	interrupted := exec.Command(binary, "release-preflight", "--mode", "verify")
	interrupted.Dir = root
	if err := interrupted.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = interrupted.Process.Kill()
			t.Fatal("blocking phase did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := interrupted.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	if err := interrupted.Wait(); err == nil {
		t.Fatal("interrupted preflight exited successfully")
	}
	newIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	newSums, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(newIndex) != string(oldIndex) || string(newSums) != string(oldSums) {
		t.Fatal("interrupted preflight replaced the prior complete evidence generation")
	}
}

func TestBuiltCommandInputDriftPreservesPriorCompleteEvidence(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(projectRoot(t), "scripts", "go-build.sh"), projectRoot(t), binary)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	root := preflightRepo(t)
	first := exec.Command(binary, "release-preflight", "--mode", "verify")
	first.Dir = root
	if output, err := first.CombinedOutput(); err != nil {
		t.Fatalf("initial verify: %v\n%s", err, output)
	}
	indexPath := filepath.Join(root, "dist", "preflight", "release-index.json")
	sumsPath := filepath.Join(root, "dist", "preflight", "SHA256SUMS")
	oldIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	oldSums, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	ready := filepath.Join(t.TempDir(), "evidence-ready")
	t.Setenv("BENCH_PREFLIGHT_EVIDENCE_READY_FILE", ready)
	drifted := exec.Command(binary, "release-preflight", "--mode", "verify")
	drifted.Dir = root
	outputFile, err := os.Create(filepath.Join(t.TempDir(), "drift-output"))
	if err != nil {
		t.Fatal(err)
	}
	drifted.Stdout, drifted.Stderr = outputFile, outputFile
	if err := drifted.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = drifted.Process.Kill()
			t.Fatal("evidence assembly did not reach the synchronization seam")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := os.WriteFile(filepath.Join(root, "go.sum"), []byte("fixture.example/module v0.0.0 h1:changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(ready); err != nil {
		t.Fatal(err)
	}
	waitErr := drifted.Wait()
	_ = outputFile.Close()
	output, _ := os.ReadFile(outputFile.Name())
	if waitErr == nil {
		t.Fatal("input drift unexpectedly passed")
	}
	if !strings.Contains(string(output), "drift") {
		t.Fatalf("input drift output lacks attribution:\n%s", output)
	}
	newIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	newSums, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(newIndex) != string(oldIndex) || string(newSums) != string(oldSums) {
		t.Fatal("input drift replaced the prior complete evidence generation")
	}
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
	for _, rel := range []string{"bin/bench.sh", ".bench/gate.sh", "scripts/build-artifacts.sh", "scripts/smoke-artifacts.sh", "scripts/go-build.sh", "scripts/platforms.json", "scripts/wrapper-assets.json", "package.json"} {
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
		if rel == "scripts/platforms.json" {
			body = `[{"os":"darwin","arch":"arm64","goos":"darwin","goarch":"arm64","runner":"macos-14"},{"os":"darwin","arch":"x64","goos":"darwin","goarch":"amd64","runner":"macos-13"},{"os":"linux","arch":"arm64","goos":"linux","goarch":"arm64","runner":"ubuntu-24.04"},{"os":"linux","arch":"x64","goos":"linux","goarch":"amd64","runner":"ubuntu-24.04"}]` + "\n"
		}
		if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ngo 1.25\ntoolchain go1.25.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.sum"), []byte("fixture.example/module v0.0.0 h1:fixture\n"), 0o644); err != nil {
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
	seedEvidenceFixture(t, root)
	return root
}

func intPtr(v int) *int { return &v }
