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
	root, wrapper, broker, manifest, stub string
}

// writeManifest replaces the install's broker manifest with body, so one case can
// state exactly the binding the land route reads.
func (i landRouteInstall) writeManifest(t *testing.T, body string) {
	t.Helper()
	if err := os.WriteFile(i.manifest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
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
	manifest := landRouteManifest(broker, "9.9.9", brokerPlatformSuffix(), fileDigest(t, broker))
	manifestPath := filepath.Join(bin, "bench-broker.manifest")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return landRouteInstall{root: install, wrapper: wrapper, broker: broker, manifest: manifestPath, stub: stub}
}

// landRouteManifest is the one spelling of the broker manifest the land route reads:
// the bound path, version, platform, and executable digest, one field to a row.
func landRouteManifest(broker, version, platform, digest string) string {
	return "path\t" + broker + "\nversion\t" + version + "\nplatform\t" + platform + "\nsha256\t" + digest + "\n"
}

func fileDigest(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
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

func TestWorktreeLandRouteGivesRecoveredGoPathToBroker(t *testing.T) {
	install := newLandRouteInstall(t)
	if err := os.WriteFile(filepath.Join(install.root, "go.mod"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	buildScript := filepath.Join(install.root, "scripts", "go-build.sh")
	if err := os.MkdirAll(filepath.Dir(buildScript), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(buildScript, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	goDir := filepath.Join(t.TempDir(), "recovered go", "bin")
	if err := os.MkdirAll(goDir, 0o755); err != nil {
		t.Fatal(err)
	}
	goExecutable := filepath.Join(goDir, "go")
	if err := os.WriteFile(goExecutable, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	profile := "export PATH='" + goDir + "':\"$PATH\"\n"
	if err := os.WriteFile(filepath.Join(home, ".bash_profile"), []byte(profile), 0o644); err != nil {
		t.Fatal(err)
	}
	broker := "#!/bin/sh\nprintf 'go=%s\\n' \"$(command -v go)\"\n"
	if err := os.WriteFile(install.broker, []byte(broker), 0o755); err != nil {
		t.Fatal(err)
	}
	install.writeManifest(t, landRouteManifest(install.broker, "9.9.9", brokerPlatformSuffix(), fileDigest(t, install.broker)))

	repo := landRouteRepo(t, "land-route [go recovery]-")
	shimDir, marker := landRouteGitShim(t)
	path := shimDir + string(os.PathListSeparator) + privateToolPath(t,
		"awk", "bash", "basename", "dirname", "env", "readlink", "sha256sum", "timeout", "tr")
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	result := owner.runAt(repo, landRouteEnv(shimDir, marker,
		"PATH="+path,
		"HOME="+home), "bash", install.wrapper,
		"worktree", "land", "--request", "r", "--base", "b", "--source-tip", "t", "-m", "m", ".")
	if result.code != 0 || result.stdout != "go="+goExecutable+"\n" {
		t.Fatalf("land route recovered Go = (%d, %q, %q), want broker Go %q", result.code, result.stdout, result.stderr, goExecutable)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("land route read the repository before broker selection: %v", err)
	}
}

// TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads is SOL03. Each
// inherited routing override makes the public landing refuse before any repository
// read: no broker starts, no git runs, and the message names the variable.
func TestWorktreeLandRouteRefusesInheritedRoutingBeforeRepositoryReads(t *testing.T) {
	install := newLandRouteInstall(t)
	primary := landRouteRepo(t, "land-route [override]-")
	plantRepositoryDecoys(t, primary)
	// Each variable appears twice: carrying a value, and set but empty. A guard that
	// tests the value rather than the presence accepts the empty spelling, and the
	// caller has then bypassed the promotion owner with an override the route reads.
	for _, row := range []struct{ name, override string }{
		{"BENCH_KIT", "BENCH_KIT=" + primary},
		{"BENCH_RUN_BINARY", "BENCH_RUN_BINARY=" + owner.selected.path},
		{"BENCH_WRAPPER", "BENCH_WRAPPER=" + install.wrapper},
		{"BENCH_KIT empty", "BENCH_KIT="},
		{"BENCH_RUN_BINARY empty", "BENCH_RUN_BINARY="},
		{"BENCH_WRAPPER empty", "BENCH_WRAPPER="},
	} {
		override := row.override
		name, _, _ := strings.Cut(override, "=")
		t.Run(row.name, func(t *testing.T) {
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

// TestWorktreeLandRouteDoesNotRederiveTheBrokerPlatform pins the single source of the
// broker platform fact. The installer writes that value; the wrapper never derives its
// own copy to compare against, so a manifest platform the host does not spell still
// routes to the broker the digest binds.
func TestWorktreeLandRouteDoesNotRederiveTheBrokerPlatform(t *testing.T) {
	install := newLandRouteInstall(t)
	install.writeManifest(t, landRouteManifest(install.broker, "9.9.9", "plan9-vax", fileDigest(t, install.broker)))
	primary := landRouteRepo(t, "land-route [platform]-")
	plantRepositoryDecoys(t, primary)

	shimDir, marker := landRouteGitShim(t)
	result := owner.runAt(primary, landRouteEnv(shimDir, marker), "bash", install.wrapper,
		"worktree", "land", "--request", "r", "--base", "b", "--source-tip", "t", "-m", "m", ".")
	if result.code != 0 || !strings.Contains(result.stdout, "broker="+install.broker+"\n") {
		t.Fatalf("foreign manifest platform = (%d, %q, %q), want the installed broker %s", result.code, result.stdout, result.stderr, install.broker)
	}
}

// TestWorktreeLandRouteRefusesEveryUnauthenticatedBroker is C1. Each row breaks one
// binding the installation manifest carries. The land route is fail-closed on all of
// them: it exits 127 with the diagnostic that names the broken binding, and neither
// the planted broker nor a repository decoy ever runs.
func TestWorktreeLandRouteRefusesEveryUnauthenticatedBroker(t *testing.T) {
	for _, row := range []struct {
		name, want string
		breakIt    func(*testing.T, landRouteInstall)
	}{
		{
			name: "missing manifest",
			want: "no promotion-broker manifest at",
			breakIt: func(t *testing.T, i landRouteInstall) {
				if err := os.Remove(i.manifest); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "incomplete manifest",
			want: "is incomplete",
			breakIt: func(t *testing.T, i landRouteInstall) {
				i.writeManifest(t, "path\t"+i.broker+"\nversion\t9.9.9\nplatform\t"+brokerPlatformSuffix()+"\n")
			},
		},
		{
			name: "version mismatch",
			want: "does not match installed package",
			breakIt: func(t *testing.T, i landRouteInstall) {
				if err := os.WriteFile(filepath.Join(i.root, "package.json"), []byte("{\n  \"version\": \"1.0.0\"\n}\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "absent broker",
			want: "is not a regular executable",
			breakIt: func(t *testing.T, i landRouteInstall) {
				if err := os.Remove(i.broker); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "symlinked broker",
			want: "is not a regular executable",
			breakIt: func(t *testing.T, i landRouteInstall) {
				target := filepath.Join(i.root, "libexec", "real-broker")
				if err := os.WriteFile(target, []byte(i.stub), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(i.broker); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(target, i.broker); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "empty broker",
			want: "is not a regular executable",
			breakIt: func(t *testing.T, i landRouteInstall) {
				if err := os.Truncate(i.broker, 0); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "non-executable broker",
			want: "is not a regular executable",
			breakIt: func(t *testing.T, i landRouteInstall) {
				if err := os.Chmod(i.broker, 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "digest mismatch",
			want: "does not match its manifest digest",
			breakIt: func(t *testing.T, i landRouteInstall) {
				if err := os.WriteFile(i.broker, []byte(i.stub+"printf 'tampered\\n'\n"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
		},
	} {
		t.Run(row.name, func(t *testing.T) {
			install := newLandRouteInstall(t)
			row.breakIt(t, install)
			repo := landRouteRepo(t, "land-route [refusal]-")
			plantRepositoryDecoys(t, repo)

			shimDir, marker := landRouteGitShim(t)
			result := owner.runAt(repo, landRouteEnv(shimDir, marker), "bash", install.wrapper,
				"worktree", "land", "--request", "r", "--base", "b", "--source-tip", "t", "-m", "m", ".")
			if result.code != 127 || !strings.Contains(result.stderr, row.want) {
				t.Fatalf("%s = (%d, %q, %q), want exit 127 naming %q", row.name, result.code, result.stdout, result.stderr, row.want)
			}
			if strings.Contains(result.stdout, "broker=") || strings.Contains(result.stdout, "decoy=") {
				t.Fatalf("%s still executed a candidate owner: %q", row.name, result.stdout)
			}
			if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("%s read the repository before refusing: %v", row.name, err)
			}
		})
	}
}
