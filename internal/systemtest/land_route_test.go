//go:build system

package systemtest

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWorktreeLandRouteSelectsOneOwnerFromEveryDirectory is SOL02. Primary-root,
// source-root, nested, and outside-repository invocations must all select the same
// installed promotion broker, with no repository read and no repository executable.
func TestWorktreeLandRouteSelectsOneOwnerFromEveryDirectory(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	install.bindBroker(t, []byte("#!/bin/sh\nprintf 'go=%s\\n' \"$(command -v go)\"\n"))

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
	t.Parallel()
	install := newLandRouteInstall(t)
	primary := landRouteRepo(t, "land-route [override]-")
	plantRepositoryDecoys(t, primary)
	// The install carries the one state the route would rebuild for. The override
	// refusals stand ahead of that recovery, so no build runs on an inherited route.
	install.plantStampedBuild(t, landRouteVersion)
	install.writeManifest(t, landRouteManifest(install.broker, "dev", brokerPlatformSuffix(), fileDigest(t, install.broker)))
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
			// LRS14: the refusal hands back the command to re-run, not the variable alone.
			if !strings.Contains(result.stderr, "'bench worktree land'") {
				t.Fatalf("override %s stderr = %q, want the exact 'bench worktree land' command", name, result.stderr)
			}
			// LRS15: the refusal states the venue that re-run uses.
			if !strings.Contains(result.stderr, "outside 'bench worktree exec'") {
				t.Fatalf("override %s stderr = %q, want the venue sentence naming 'bench worktree exec'", name, result.stderr)
			}
			if result.stdout != "" {
				t.Fatalf("override %s still started an executable: %q", name, result.stdout)
			}
			if got := install.buildRuns(t); got != 0 {
				t.Fatalf("override %s ran the stamped build %d times, want none before the refusal", name, got)
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
	t.Parallel()
	install := newLandRouteInstall(t)
	install.writeManifest(t, landRouteManifest(install.broker, landRouteVersion, "plan9-vax", fileDigest(t, install.broker)))
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
// The unstamped rows are the exception the table carries: a `dev` version is the
// build's omission, not a tamper, so it earns one stamped rebuild and one re-read.
// The recovered row lands; the row the rebuild leaves unstamped refuses like any
// other version mismatch. wantBuilds pins the pass count on every row, so a rebuild
// that reaches a tamper row, or that runs twice, reds here.
func TestWorktreeLandRouteRefusesEveryUnauthenticatedBroker(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, want string
		recovers   bool
		wantBuilds int
		breakIt    func(*testing.T, landRouteInstall)
	}{
		{
			name:       "unstamped manifest recovers",
			recovers:   true,
			wantBuilds: 1,
			breakIt: func(t *testing.T, i landRouteInstall) {
				i.plantStampedBuild(t, landRouteVersion)
				i.writeManifest(t, landRouteManifest(i.broker, "dev", brokerPlatformSuffix(), fileDigest(t, i.broker)))
			},
		},
		{
			name:       "unstamped manifest after the rebuild",
			want:       "does not match installed package",
			wantBuilds: 1,
			breakIt: func(t *testing.T, i landRouteInstall) {
				i.plantStampedBuild(t, "dev")
				i.writeManifest(t, landRouteManifest(i.broker, "dev", brokerPlatformSuffix(), fileDigest(t, i.broker)))
			},
		},
		{
			name:       "unstamped manifest with a tampered broker",
			want:       "does not match its manifest digest",
			wantBuilds: 1,
			breakIt: func(t *testing.T, i landRouteInstall) {
				// The rebuild republishes the manifest, then leaves the executable it
				// bound tampered. A recovered version never authorizes bytes the digest
				// does not bind.
				i.plantStampedBuild(t, landRouteVersion, "printf 'tampered\\n' >> '"+i.broker+"'")
				i.writeManifest(t, landRouteManifest(i.broker, "dev", brokerPlatformSuffix(), fileDigest(t, i.broker)))
			},
		},
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
				i.writeManifest(t, "path\t"+i.broker+"\nversion\t"+landRouteVersion+"\nplatform\t"+brokerPlatformSuffix()+"\n")
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
			if row.recovers {
				if result.code != 0 || !strings.Contains(result.stdout, "broker="+install.broker+"\n") {
					t.Fatalf("%s = (%d, %q, %q), want the recovered broker %s", row.name, result.code, result.stdout, result.stderr, install.broker)
				}
			} else {
				if result.code != 127 || !strings.Contains(result.stderr, row.want) {
					t.Fatalf("%s = (%d, %q, %q), want exit 127 naming %q", row.name, result.code, result.stdout, result.stderr, row.want)
				}
				// Every refusal hands back the one repair that republishes the broker.
				if !strings.Contains(result.stderr, "bench doctor --fix") {
					t.Fatalf("%s stderr = %q, want the repair advice", row.name, result.stderr)
				}
				if strings.Contains(result.stdout, "broker=") || strings.Contains(result.stdout, "decoy=") {
					t.Fatalf("%s still executed a candidate owner: %q", row.name, result.stdout)
				}
			}
			if got := install.buildRuns(t); got != row.wantBuilds {
				t.Fatalf("%s ran the stamped build %d times, want %d", row.name, got, row.wantBuilds)
			}
			if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("%s read the repository before refusing: %v", row.name, err)
			}
		})
	}
}
