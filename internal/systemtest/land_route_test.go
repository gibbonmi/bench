//go:build system

package systemtest

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// brokerPlatformSuffix derives the npm-spelled platform suffix the wrapper computes
// from uname, so the fixture manifest matches the host the wrapper reports.
func brokerPlatformSuffix() string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	return runtime.GOOS + "-" + arch
}

type landRouteInstall struct {
	wrapper, broker string
}

// newLandRouteInstall fabricates an installed distribution: the kit's wrapper, a
// package version, a stub promotion broker that prints its own identity, and the
// broker manifest binding path, version, platform, and executable digest.
func newLandRouteInstall(t *testing.T) landRouteInstall {
	t.Helper()
	install, err := os.MkdirTemp(owner.root, "land-route [install]-")
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(install, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile(filepath.Join(owner.kit, "bin", "bench.sh"))
	if err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(bin, "bench.sh")
	if err := os.WriteFile(wrapper, source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(install, "package.json"), []byte("{\n  \"version\": \"9.9.9\"\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	broker := filepath.Join(install, "libexec", "broker")
	if err := os.MkdirAll(filepath.Dir(broker), 0o755); err != nil {
		t.Fatal(err)
	}
	stub := "#!/bin/sh\nprintf 'broker=%s\\n' \"$0\"\nprintf 'argv=%s\\n' \"$*\"\n"
	if err := os.WriteFile(broker, []byte(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(stub))
	manifest := "path\t" + broker + "\nversion\t9.9.9\nplatform\t" + brokerPlatformSuffix() + "\nsha256\t" + hex.EncodeToString(sum[:]) + "\n"
	if err := os.WriteFile(filepath.Join(bin, "bench-broker.manifest"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return landRouteInstall{wrapper: wrapper, broker: broker}
}

// landRouteGitShim puts a failing git first on PATH that records every invocation.
// The land route must select its owner before any repository read, so the marker
// staying absent is the before-repository-reads proof.
func landRouteGitShim(t *testing.T) (shimDir, marker string) {
	t.Helper()
	shimDir = t.TempDir()
	marker = filepath.Join(t.TempDir(), "git-ran")
	shim := "#!/bin/sh\n: > \"$LAND_ROUTE_GIT_MARKER\"\nexit 1\n"
	if err := os.WriteFile(filepath.Join(shimDir, "git"), []byte(shim), 0o755); err != nil {
		t.Fatal(err)
	}
	return shimDir, marker
}

// landRouteEnv strips every inherited routing override and fronts PATH with the git
// probe shim. Extra entries layer on top.
func landRouteEnv(shimDir, marker string, extra ...string) []string {
	base := []string{
		"BENCH_KIT", "BENCH_RUN_BINARY", "BENCH_WRAPPER",
		"PATH=" + shimDir + string(os.PathListSeparator) + os.Getenv("PATH"),
		"LAND_ROUTE_GIT_MARKER=" + marker,
	}
	return append(base, extra...)
}

// plantRepositoryDecoys plants the executables a current-directory resolver would
// select: a repo-local dev build and a bundled platform package binary.
func plantRepositoryDecoys(t *testing.T, root string) {
	t.Helper()
	decoy := "#!/bin/sh\nprintf 'decoy=%s\\n' \"$0\"\n"
	for _, rel := range []string{
		filepath.Join("dist", "bench"),
		filepath.Join("node_modules", "@redbench", brokerPlatformSuffix(), "bin", "bench"),
	} {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(decoy), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func landRouteRepo(t *testing.T, name string) string {
	t.Helper()
	repo, err := os.MkdirTemp(owner.root, name)
	if err != nil {
		t.Fatal(err)
	}
	if result := owner.runAt(repo, nil, "git", "init", "-q", "-b", "main"); result.code != 0 {
		t.Fatalf("git init %s = (%d, %q)", name, result.code, result.stderr)
	}
	return repo
}

// TestWorktreeLandRouteSelectsOneOwnerFromEveryDirectory is SOL02. Primary-root,
// source-root, nested, and outside-repository invocations must all select the same
// installed promotion broker, with no repository read and no repository executable.
func TestWorktreeLandRouteSelectsOneOwnerFromEveryDirectory(t *testing.T) {
	install := newLandRouteInstall(t)
	primary := landRouteRepo(t, "land-route [primary]-")
	source := landRouteRepo(t, "land-route [source]-")
	plantRepositoryDecoys(t, primary)
	plantRepositoryDecoys(t, source)
	nested := filepath.Join(primary, "nested", "deep")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	outside, err := os.MkdirTemp(owner.root, "land-route [outside]-")
	if err != nil {
		t.Fatal(err)
	}

	first := ""
	for _, cwd := range []string{primary, source, nested, outside} {
		shimDir, marker := landRouteGitShim(t)
		result := owner.runAt(cwd, landRouteEnv(shimDir, marker), "bash", install.wrapper,
			"worktree", "land", "--request", "r", "--base", "b", "--source-tip", "t", "-m", "m", ".")
		if result.code != 0 || !strings.Contains(result.stdout, "broker="+install.broker+"\n") {
			t.Fatalf("land route from %s = (%d, %q, %q), want the installed broker %s", cwd, result.code, result.stdout, result.stderr, install.broker)
		}
		if !strings.Contains(result.stdout, "argv=worktree land --request r --base b --source-tip t -m m .\n") {
			t.Fatalf("land route from %s forwarded %q", cwd, result.stdout)
		}
		if strings.Contains(result.stdout, "decoy=") {
			t.Fatalf("land route from %s selected a repository executable: %q", cwd, result.stdout)
		}
		if first == "" {
			first = result.stdout
		} else if result.stdout != first {
			t.Fatalf("owner selection changed with the current directory: first=%q now=%q", first, result.stdout)
		}
		if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("land route from %s read the repository before selecting its owner: %v", cwd, err)
		}
	}
}

// TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads is SOL03. Each
// inherited routing override makes the public landing refuse before any repository
// read: no broker starts, no git runs, and the message names the variable.
func TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads(t *testing.T) {
	install := newLandRouteInstall(t)
	primary := landRouteRepo(t, "land-route [override]-")
	plantRepositoryDecoys(t, primary)
	for _, override := range []string{
		"BENCH_KIT=" + primary,
		"BENCH_RUN_BINARY=" + owner.selected.path,
		"BENCH_WRAPPER=" + install.wrapper,
	} {
		name, _, _ := strings.Cut(override, "=")
		t.Run(name, func(t *testing.T) {
			shimDir, marker := landRouteGitShim(t)
			result := owner.runAt(primary, landRouteEnv(shimDir, marker, override), "bash", install.wrapper,
				"worktree", "land", "--request", "r", "--base", "b", "--source-tip", "t", "-m", "m", ".")
			if result.code != 1 || !strings.Contains(result.stderr, name) {
				t.Fatalf("override %s = (%d, %q, %q), want a refusal naming it", name, result.code, result.stdout, result.stderr)
			}
			if result.stdout != "" {
				t.Fatalf("override %s still started an executable: %q", name, result.stdout)
			}
			if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("override %s read the repository before refusing: %v", name, err)
			}
		})
	}
}
