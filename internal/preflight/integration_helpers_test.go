package preflight

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

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
	cmd := exec.Command("git", "tag", "v0.2.0")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git tag: %v %s", err, output)
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
	for _, rel := range []string{"bin/bench.sh", ".bench/gate.sh", "internal/releaseevidence/registry.json", "internal/releaseevidence/requirements.json", "scripts/build-artifacts.sh", "scripts/build-offline-archives.sh", "scripts/write-deterministic-archive.mjs", "scripts/compare-artifacts.sh", "scripts/native-proof.sh", "scripts/aggregate-native-proofs.sh", "scripts/build-release-evidence.mjs", "scripts/smoke-artifacts.sh", "scripts/smoke-offline.sh", "scripts/offline-registry.mjs", "scripts/go-build.sh", "scripts/platforms.json", "scripts/wrapper-assets.json", "package.json"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		body := "{}\n"
		if rel == "internal/releaseevidence/registry.json" || rel == "internal/releaseevidence/requirements.json" || rel == "scripts/build-release-evidence.mjs" || rel == "scripts/build-offline-archives.sh" || rel == "scripts/write-deterministic-archive.mjs" {
			data, err := os.ReadFile(filepath.Join(projectRoot(t), filepath.FromSlash(rel)))
			if err != nil {
				t.Fatal(err)
			}
			body = string(data)
		}
		if (filepath.Ext(path) == ".sh" && rel != "scripts/build-offline-archives.sh") || rel == "bin/bench.sh" {
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
