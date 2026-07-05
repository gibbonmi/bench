package contract

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGoRoutingContracts(t *testing.T) {
	t.Parallel()
	skipIfSubjectBenchMissing(t)
	runParallel(t, "bench version output contract", testGoRoutingVersionOutput)
	runParallel(t, "bench version failed outside a git repo", testGoRoutingVersionOutsideRepo)
	runParallel(t, "version-routing seam contract failed", testGoRoutingFabricatedVersionRouting)
	runParallel(t, "linked-worktree binary-resolution contract failed", testGoRoutingLinkedWorktreeBinaryResolution)
}

func testGoRoutingVersionOutput(t *testing.T) {
	f := NewFixture(t)
	version := goRoutingPackageVersion(t)
	want := "benchkit " + version + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")"

	out := f.Bench("version")

	out.RequireExit(0)
	got := strings.TrimSuffix(out.Stdout, "\n")
	if got != want {
		t.Fatalf("bench version output %q != expected %q\nstderr:\n%s", got, want, out.Stderr)
	}
	if got == "" || strings.Contains(got, "\n") || !strings.HasSuffix(out.Stdout, "\n") {
		t.Fatalf("bench version was not exactly one line\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
}

func testGoRoutingVersionOutsideRepo(t *testing.T) {
	f := NewFixture(t, WithNoRepo())

	out := f.Bench("version")

	out.RequireExit(0)
}

func testGoRoutingFabricatedVersionRouting(t *testing.T) {
	f := NewFixture(t, WithNoRepo())
	kit := filepath.Join(f.Root, "a b", "kit")
	goRoutingCopyTree(t, filepath.Join(SubjectRoot(t), "bin"), filepath.Join(kit, "bin"))
	hostPackage := strings.TrimSpace(goRoutingNode(t, "process.stdout.write('@benchkit/'+process.platform+'-'+process.arch)"))
	run := func(env map[string]string) Probe {
		return f.RunEnv(env, "bash", filepath.Join(kit, "bin", "bench.sh"), "version")
	}

	goRoutingWriteStub(t, filepath.Join(kit, "dist", "bench"), "devbuild")
	goRoutingWriteStub(t, filepath.Join(kit, "node_modules", hostPackage, "bin", "bench"), "bundled")
	goRoutingWriteStub(t, filepath.Join(kit, "..", hostPackage, "bin", "bench"), "hoisted")
	if got := strings.TrimSpace(run(nil).Stdout); got != "devbuild" {
		t.Fatalf("dev build not preferred over bundled: got %q", got)
	}
	goRoutingRemove(t, filepath.Join(kit, "dist", "bench"))
	if got := strings.TrimSpace(run(nil).Stdout); got != "bundled" {
		t.Fatalf("bundled not preferred over hoisted: got %q", got)
	}
	goRoutingRemove(t, filepath.Join(kit, "node_modules", hostPackage, "bin", "bench"))
	if got := strings.TrimSpace(run(nil).Stdout); got != "hoisted" {
		t.Fatalf("hoisted sibling not resolved: got %q", got)
	}
	goRoutingRemove(t, filepath.Join(kit, "..", hostPackage, "bin", "bench"))

	out := run(nil)
	out.RequireExit(127)
	out.RequireContains(out.Stderr+out.Stdout, hostPackage)
	out.RequireContains(out.Stderr+out.Stdout, "npm install")

	if err := os.WriteFile(filepath.Join(kit, "dist", "bench"), nil, 0o755); err != nil {
		t.Fatalf("write empty dist/bench: %v", err)
	}
	out = run(nil)
	out.RequireExit(127)

	if err := os.WriteFile(filepath.Join(kit, "dist", "bench"), []byte("#!/bin/sh\necho nope\n"), 0o644); err != nil {
		t.Fatalf("write non-executable dist/bench: %v", err)
	}
	if err := os.Chmod(filepath.Join(kit, "dist", "bench"), 0o644); err != nil {
		t.Fatalf("chmod non-executable dist/bench: %v", err)
	}
	out = run(nil)
	out.RequireExit(127)
	goRoutingRemove(t, filepath.Join(kit, "dist", "bench"))

	goRoutingWriteStub(t, filepath.Join(kit, "dist", "bench"), "devbuild")
	link := filepath.Join(f.Root, "benchlink")
	if err := os.Symlink(filepath.Join(kit, "bin", "bench.sh"), link); err != nil {
		t.Fatalf("symlink bench wrapper: %v", err)
	}
	if got := strings.TrimSpace(f.Run("bash", link, "version").Stdout); got != "devbuild" {
		t.Fatalf("symlink invocation did not resolve the kit root: got %q", got)
	}
	goRoutingRemove(t, filepath.Join(kit, "dist", "bench"))

	fakebin := filepath.Join(f.Root, "fakebin")
	f.WriteExecutable("fakebin/uname", "#!/bin/sh\ncase \"$1\" in -s) echo Plan9;; -m) echo sparc64;; *) exec /usr/bin/uname \"$@\";; esac\n")
	out = run(map[string]string{"PATH": fakebin + string(os.PathListSeparator) + os.Getenv("PATH")})
	out.RequireExit(2)
	out.RequireContains(strings.ToLower(out.Stderr+out.Stdout), "unsupported platform")
	if strings.Contains(out.Stderr+out.Stdout, "@benchkit/") {
		t.Fatalf("unsupported platform named a package\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
}

func testGoRoutingLinkedWorktreeBinaryResolution(t *testing.T) {
	f := NewFixture(t, WithNoRepo())
	main := filepath.Join(f.Root, "main")
	linked := filepath.Join(f.Root, "linked")
	if err := os.MkdirAll(main, 0o755); err != nil {
		t.Fatalf("mkdir main repo: %v", err)
	}
	repo := linkFixtureAt(t, main, f.Env)
	repo.Git("init", "-q")
	goRoutingCopyTree(t, filepath.Join(SubjectRoot(t), "bin"), filepath.Join(main, "bin"))
	repo.Git("add", "-A")
	repo.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "init")
	goRoutingWriteStub(t, filepath.Join(main, "dist", "bench"), "mainbuild")
	repo.Git("worktree", "add", "-q", "--detach", linked, "HEAD")

	out := f.Run("bash", filepath.Join(linked, "bin", "bench.sh"), "version")
	out.RequireExit(0)
	if got := strings.TrimSpace(out.Stdout); got != "mainbuild" {
		t.Fatalf("worktree wrapper did not resolve the main tree's binary: got %q", got)
	}

	goRoutingWriteStub(t, filepath.Join(linked, "dist", "bench"), "localbuild")
	out = f.Run("bash", filepath.Join(linked, "bin", "bench.sh"), "version")
	out.RequireExit(0)
	if got := strings.TrimSpace(out.Stdout); got != "localbuild" {
		t.Fatalf("worktree-local build not preferred over the main tree's: got %q", got)
	}
}

func goRoutingPackageVersion(t testing.TB) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(SubjectRoot(t), "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatalf("parse package.json: %v", err)
	}
	return pkg.Version
}

func goRoutingNode(t testing.TB, script string) string {
	t.Helper()
	out, err := exec.Command("node", "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("node script failed: %v\n%s", err, out)
	}
	return string(out)
}

func goRoutingWriteStub(t testing.TB, path, output string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho "+output+"\n"), 0o755); err != nil {
		t.Fatalf("write stub %s: %v", path, err)
	}
}

func goRoutingCopyTree(t testing.TB, src, dst string) {
	t.Helper()
	if err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	}); err != nil {
		t.Fatalf("copy %s to %s: %v", src, dst, err)
	}
}

func goRoutingRemove(t testing.TB, path string) {
	t.Helper()
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove %s: %v", path, err)
	}
}
