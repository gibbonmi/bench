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

// landRouteVersion is the package version that every fabricated installation stamps.
const landRouteVersion = "9.9.9"

type landRouteInstall struct {
	root, wrapper, broker, manifest, stub string
}

// bindBroker writes body as the install's broker and binds it at landRouteVersion.
func (i landRouteInstall) bindBroker(t *testing.T, body []byte) {
	t.Helper()
	if err := os.WriteFile(i.broker, body, 0o755); err != nil {
		t.Fatal(err)
	}
	i.writeManifest(t, landRouteManifest(i.broker, landRouteVersion, brokerPlatformSuffix(), fileDigest(t, i.broker)))
}

// writeManifest replaces the install's broker manifest with body, so one case can
// state exactly the binding the land route reads.
func (i landRouteInstall) writeManifest(t *testing.T, body string) {
	t.Helper()
	if err := os.WriteFile(i.manifest, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// plantStampedBuild writes the install's stamped build script. Each run appends one
// line to the install's build-run log and republishes the broker manifest at version,
// which is what the real subject-mode build does inside its publication transaction.
// A case states the version the rebuild publishes, so the route meets a recovered
// manifest or a still-unstamped one, plus any shell line the publication runs after it.
func (i landRouteInstall) plantStampedBuild(t *testing.T, version string, after ...string) {
	t.Helper()
	script := filepath.Join(i.root, "scripts", "go-build.sh")
	if err := os.MkdirAll(filepath.Dir(script), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "#!/bin/sh\necho run >> '" + i.buildLog() + "'\ncat > '" + i.manifest + "' <<'MANIFEST'\n" +
		landRouteManifest(i.broker, version, brokerPlatformSuffix(), fileDigest(t, i.broker)) + "MANIFEST\n"
	for _, line := range after {
		body += line + "\n"
	}
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func (i landRouteInstall) buildLog() string { return filepath.Join(i.root, "build-runs") }

// buildRuns counts how many times the planted build script ran. The recovery is
// single-pass by decision, so this count is what separates one rebuild from a loop.
func (i landRouteInstall) buildRuns(t *testing.T) int {
	t.Helper()
	body, err := os.ReadFile(i.buildLog())
	if errors.Is(err, os.ErrNotExist) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(body), "\n")
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
	bin, broker := filepath.Join(install, "bin"), filepath.Join(install, "libexec", "broker")
	for _, dir := range []string{bin, filepath.Dir(broker)} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	source, err := os.ReadFile(filepath.Join(owner.kit, "bin", "bench.sh"))
	if err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(bin, "bench.sh")
	if err := os.WriteFile(wrapper, source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(install, "package.json"), []byte("{\n  \"version\": \""+landRouteVersion+"\"\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	created := landRouteInstall{root: install, wrapper: wrapper, broker: broker, manifest: filepath.Join(bin, "bench-broker.manifest"),
		stub: "#!/bin/sh\nprintf 'broker=%s\\n' \"$0\"\nprintf 'argv=%s\\n' \"$*\"\n"}
	created.bindBroker(t, []byte(created.stub))
	return created
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
	repo, err := owner.journeyRepository(name)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}
