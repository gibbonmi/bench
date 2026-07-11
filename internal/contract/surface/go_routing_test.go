package surface

import (
	"encoding/json"
	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestGoRoutingContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "bench version output contract", testGoRoutingVersionOutput)
	contract.RunParallel(t, "bench version failed outside a git repo", testGoRoutingVersionOutsideRepo)
	contract.RunParallel(t, "version-routing seam contract failed", testGoRoutingFabricatedVersionRouting)
	contract.RunParallel(t, "linked-worktree binary-resolution contract failed", testGoRoutingLinkedWorktreeBinaryResolution)
	contract.RunParallel(t, "repo-local wrapper forwarding contract failed", testGoRoutingRepoLocalWrapperForwarding)
	contract.RunParallel(t, "repo-local linked-worktree forwarding contract failed", testGoRoutingShimLinkedWorktree)
	contract.RunParallel(t, "unknown subcommand exits 2 on stderr", testGoRoutingUnknownSubcommandExits2OnStderr)
	contract.RunParallel(t, "help variants stay on stdout at exit 0", testGoRoutingHelpVariantsStayOnStdoutExit0)
	contract.RunParallel(t, "--version routes to the version subcommand", testGoRoutingVersionFlagMatchesVersionSubcommand)
	contract.RunParallel(t, "resolve_script_path capped a symlink cycle instead of hanging", testGoRoutingResolveScriptPathHopCap)
}

func testGoRoutingShimLinkedWorktree(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	main := filepath.Join(f.Root, "main")
	linked := filepath.Join(f.Root, "linked")
	if err := os.MkdirAll(main, 0755); err != nil {
		t.Fatal(err)
	}
	repo := contract.NewFixtureAt(t, main, f.Env)
	repo.Git("init", "-q")
	repo.WriteFile("seed", "x")
	repo.Git("add", "seed")
	repo.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "seed")
	repo.Git("worktree", "add", "-q", "--detach", linked, "HEAD")
	local := filepath.Join(main, ".bench", "bin", "bench.sh")
	target := filepath.Join(f.Root, "global")
	shim := filepath.Join(f.Root, "bench")
	contract.WriteExecutableAbs(t, local, "#!/bin/sh\nprintf 'main-local:%s|%s\\n' \"$1\" \"$2\"\n")
	contract.WriteExecutableAbs(t, target, "#!/bin/sh\necho global\nexit 7\n")
	contract.WriteExecutableAbs(t, shim, adopt.ShimContent(target)+"\n")
	p := contract.RunAt(t, f, linked, nil, shim, "--context", "--full")
	p.RequireExit(0)
	if p.Stdout != "main-local:--context|--full\n" {
		t.Fatalf("linked-worktree routing = %q", p.Stdout)
	}
}

func testGoRoutingRepoLocalWrapperForwarding(t *testing.T) {
	f := contract.NewFixture(t)
	target := filepath.Join(f.Root, "global target")
	local := filepath.Join(f.Root, ".bench", "bin", "bench.sh")
	shim := filepath.Join(f.Root, "top bench")
	contract.WriteExecutableAbs(t, target, "#!/bin/sh\nprintf 'global:%s\\n' \"$*\"\nexit 7\n")
	contract.WriteExecutableAbs(t, local, "#!/bin/sh\nprintf 'local:%s|%s\\n' \"$1\" \"$2\"\nexit 2\n")
	contract.WriteExecutableAbs(t, shim, adopt.ShimContent(target)+"\n")

	out := f.Run(shim, "--context", "--full")
	out.RequireExit(2)
	if out.Stdout != "local:--context|--full\n" || out.Stderr != "" {
		t.Fatalf("repo-local forwarding changed argv/bytes\nstdout:%q\nstderr:%q", out.Stdout, out.Stderr)
	}

	if err := os.Remove(local); err != nil {
		t.Fatal(err)
	}
	out = f.Run(shim, "--context", "--full")
	out.RequireExit(7)
	if out.Stdout != "global:--context --full\n" {
		t.Fatalf("missing-local fallback = %q", out.Stdout)
	}

	if err := os.Symlink(target, local); err != nil {
		t.Fatal(err)
	}
	out = f.Run(shim, "same")
	out.RequireExit(7)
	if out.Stdout != "global:same\n" {
		t.Fatalf("self-resolution did not fall back: %q", out.Stdout)
	}
}

// testGoRoutingUnknownSubcommandExits2OnStderr pins the typo case: an unrecognized
// token must be distinguishable from success (exit 2, message on stderr), not the
// former exit-0/stdout help fallthrough that made a typo indistinguishable from a
// legitimate no-op to a wrapping script. It also pins flag-shaped tokens, the empty
// string (the shell wrapper used to swallow "" into help via `${1:-help}`'s
// unset-or-empty default), and a multi-word token — the unknown-subcommand contract
// covers more than the plain typo shape.
func testGoRoutingUnknownSubcommandExits2OnStderr(t *testing.T) {
	for _, tok := range []string{"frobnicate", "--frobnicate", "-x", "", "foo bar"} {
		t.Run(tok, func(t *testing.T) {
			f := contract.NewFixture(t)

			out := f.Bench(tok)

			out.RequireExit(2)
			if strings.TrimSpace(out.Stdout) != "" {
				t.Fatalf("bench %q wrote to stdout, want it silent there:\nstdout:\n%s", tok, out.Stdout)
			}
			if strings.TrimSpace(out.Stderr) == "" {
				t.Fatalf("bench %q wrote nothing to stderr, want a message there", tok)
			}
		})
	}
}

// testGoRoutingHelpVariantsStayOnStdoutExit0 pins that splitting the typo case out of
// the catch-all does not cannibalize the legitimate help request: bare invocation,
// `help`, `--help`, and `-h` all still print help to stdout at exit 0.
func testGoRoutingHelpVariantsStayOnStdoutExit0(t *testing.T) {
	f := contract.NewFixture(t)

	for _, args := range [][]string{{}, {"help"}, {"--help"}, {"-h"}} {
		out := f.Bench(args...)
		out.RequireExit(0)
		out.RequireContains(out.Stdout, "bench link")
		if strings.TrimSpace(out.Stderr) != "" {
			t.Fatalf("bench %v wrote to stderr, want help silent there:\nstderr:\n%s", args, out.Stderr)
		}
	}
}

// testGoRoutingVersionFlagMatchesVersionSubcommand pins that `--version` reaches the
// same implementation `bench version` does, rather than falling into the help
// catch-all it shared with a typo before this split.
func testGoRoutingVersionFlagMatchesVersionSubcommand(t *testing.T) {
	f := contract.NewFixture(t)

	want := f.Bench("version")
	want.RequireExit(0)

	got := f.Bench("--version")
	got.RequireExit(0)
	if got.Stdout != want.Stdout {
		t.Fatalf("bench --version stdout = %q, want %q (same as bench version)", got.Stdout, want.Stdout)
	}
}

func testGoRoutingVersionOutput(t *testing.T) {
	f := contract.NewFixture(t)
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
	f := contract.NewFixture(t, contract.WithNoRepo())

	out := f.Bench("version")

	out.RequireExit(0)
}

func testGoRoutingFabricatedVersionRouting(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())
	kit := filepath.Join(f.Root, "a b", "kit")
	goRoutingCopyTree(t, filepath.Join(contract.SubjectRoot(t), "bin"), filepath.Join(kit, "bin"))
	hostPackage := strings.TrimSpace(goRoutingNode(t, "process.stdout.write('@redbench/'+process.platform+'-'+process.arch)"))
	run := func(env map[string]string) contract.Probe {
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
	requireInterim127Remedy(t, out.Stderr+out.Stdout)

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
	f := contract.NewFixture(t, contract.WithNoRepo())
	main := filepath.Join(f.Root, "main")
	linked := filepath.Join(f.Root, "linked")
	if err := os.MkdirAll(main, 0o755); err != nil {
		t.Fatalf("mkdir main repo: %v", err)
	}
	repo := linkFixtureAt(t, main, f.Env)
	repo.Git("init", "-q")
	goRoutingCopyTree(t, filepath.Join(contract.SubjectRoot(t), "bin"), filepath.Join(main, "bin"))
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
	data, err := os.ReadFile(filepath.Join(contract.SubjectRoot(t), "package.json"))
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

// testGoRoutingResolveScriptPathHopCap drives the wrapper's real, unmodified
// resolve_script_path function into a genuine symlink cycle and asserts it fails
// fast with a structured error rather than hanging.
//
// A cyclic path can never be the literal `bash <path>` invocation target: opening
// the entry point to read it as a script requires the kernel to fully resolve the
// symlink chain, and a real cycle ELOOPs there before resolve_script_path's own
// loop ever runs — that would just be exercising the OS's own bound, not the
// wrapper's. The actual gap is different: resolve_script_path's loop walks
// $source with `readlink`/`[[ -L ]]` only, which (unlike open()) never needs to
// fully dereference a symlink's target, so it never ELOOPs on a cycle by itself.
// To reach that loop with $source pointed at a genuine cycle, this test exploits
// how bash reports a function's origin: BASH_SOURCE[0] inside any function is
// always non-empty — either the file that defined it, or, when the function was
// eval'd from a string instead of sourced from a file, the literal placeholder
// "environment". resolve_script_path's own fallback, `${BASH_SOURCE[0]:-$0}`,
// then resolves to that placeholder, which the loop treats as a plain filename
// relative to the working directory. So eval-ing the wrapper's real source (not a
// reimplementation) and naming a self-referential symlink "environment" in the
// working directory drives the production loop into exactly the cycle the hop cap
// defends against.
func testGoRoutingResolveScriptPathHopCap(t *testing.T) {
	f := contract.NewFixture(t, contract.WithNoRepo())

	if err := os.Symlink("environment", filepath.Join(f.Root, "environment")); err != nil {
		t.Fatalf("create self-referential symlink: %v", err)
	}

	wrapper := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	src, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatalf("read bin/bench.sh: %v", err)
	}

	// Clear the driver's own positional params before eval-ing the wrapper's full
	// text, so its bottom `case "${1-help}"` dispatch defaults to the harmless
	// help branch instead of routing on whatever happens to occupy $1.
	driver := `src="$1"; set --; eval "$src" >/dev/null 2>&1; resolve_script_path`

	out := contract.RunAtWithTimeout(t, f, f.Root, nil, 10*time.Second, "bash", "-c", driver, "probe", string(src))

	if out.TimedOut {
		t.Fatal("resolve_script_path hung chasing a self-referential symlink cycle; want a bounded hop cap")
	}
	if out.ExitCode == 0 {
		t.Fatalf("resolve_script_path resolved a symlink cycle instead of erroring\nstdout:\n%s", out.Stdout)
	}
	out.RequireContains(out.Stderr, "symlink cycle")
}
