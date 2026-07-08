package surface

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestGitInstallProbe is seam C of the npm-identity spec: before the packages are
// published, a consumer installs Bench straight from the git repo, so package.json's
// prepare script must Go-build the core at install time and dist/ must ride in files[]
// for the packed tree to carry a runnable binary. The probe stages a git repo from the
// working tree (uncommitted changes included, so it exercises the package.json under
// test), installs it over the git protocol into a throwaway prefix, and runs the
// installed `bench version`. Pre-change — no prepare, no dist/ in files[] — the install
// carries no binary and the run is non-zero; the prepare + dist/ change is what turns
// it green.
func TestGitInstallProbe(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	requireExecutables(t, "git", "npm", "node", "go")
	subject := contract.SubjectRoot(t)

	repo := stageWorkingTreeRepo(t, subject)
	prefix := t.TempDir()
	env := gitInstallEnv(t, t.TempDir(), t.TempDir())
	f := contract.NewFixtureAt(t, prefix, env)

	install := f.Run("npm", "install", "--prefix", prefix, "--no-fund", "--no-audit", "git+file://"+repo)
	install.RequireExit(0)

	benchBin := filepath.Join(prefix, "node_modules", ".bin", "bench")
	out := f.Run(benchBin, "version")
	out.RequireExit(0)
	if !strings.HasPrefix(strings.TrimSpace(out.Stdout), "benchkit ") {
		t.Fatalf("git-installed bench version did not print the version line\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
}

func requireExecutables(t *testing.T, names ...string) {
	t.Helper()
	for _, name := range names {
		if _, err := exec.LookPath(name); err != nil {
			t.Skipf("git-install probe needs %s on PATH: %v", name, err)
		}
	}
}

// stageWorkingTreeRepo copies the subject's tracked files at their working-tree content
// (not HEAD, so the package.json under test is what gets installed) into a fresh temp
// repo and commits them, yielding a git+file:// install source.
func stageWorkingTreeRepo(t *testing.T, subject string) string {
	t.Helper()
	repo := t.TempDir()
	lsOut, err := exec.Command("git", "-C", subject, "ls-files", "-z").Output()
	if err != nil {
		t.Fatalf("git ls-files: %v", err)
	}
	for _, rel := range strings.Split(strings.TrimRight(string(lsOut), "\x00"), "\x00") {
		if rel == "" {
			continue
		}
		src := filepath.Join(subject, filepath.FromSlash(rel))
		info, err := os.Lstat(src)
		if err != nil {
			continue // tracked but removed in the working tree
		}
		dst := filepath.Join(repo, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", filepath.Dir(dst), err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(src)
			if err != nil {
				t.Fatalf("readlink %s: %v", src, err)
			}
			if err := os.Symlink(target, dst); err != nil {
				t.Fatalf("symlink %s: %v", dst, err)
			}
			continue
		}
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, info.Mode().Perm()); err != nil {
			t.Fatalf("write %s: %v", dst, err)
		}
	}
	git := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=probe", "GIT_AUTHOR_EMAIL=probe@example.com",
			"GIT_COMMITTER_NAME=probe", "GIT_COMMITTER_EMAIL=probe@example.com")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q")
	git("add", "-A")
	git("commit", "-q", "-m", "probe")
	return repo
}

// gitInstallEnv isolates npm's HOME and cache to throwaway dirs while pinning Go's caches
// to their ambient locations, so the prepare build reuses the real module cache instead
// of re-downloading into (read-only) temp dirs the test then cannot clean up. npm_config_omit
// propagates to the nested git-dependency build, keeping the whole install off the network:
// the @redbench/* platform packages are unpublished until publish day and the git build
// path does not need them (prepare builds dist/bench). Their published-state tolerance is a
// release-smoke concern, not this gate's.
func gitInstallEnv(t *testing.T, home, npmCache string) map[string]string {
	t.Helper()
	goEnv := func(key string) string {
		out, err := exec.Command("go", "env", key).Output()
		if err != nil {
			t.Fatalf("go env %s: %v", key, err)
		}
		return strings.TrimSpace(string(out))
	}
	return map[string]string{
		"HOME":             home,
		"npm_config_cache": npmCache,
		"npm_config_omit":  "optional",
		"GOMODCACHE":       goEnv("GOMODCACHE"),
		"GOCACHE":          goEnv("GOCACHE"),
		"GOPATH":           goEnv("GOPATH"),
	}
}
