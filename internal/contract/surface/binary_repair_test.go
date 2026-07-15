package surface

import (
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBinaryRepairContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "binary repair download-and-run contract failed", testRepairDownloadsAndRuns)
	contract.RunParallel(t, "binary repair bad-integrity contract failed", testRepairRefusesBadHash)
	contract.RunParallel(t, "binary repair malformed-tar contract failed", testRepairRefusesMalformedTar)
	contract.RunParallel(t, "binary repair offline fallback contract failed", testRepairOffline)
	contract.RunParallel(t, "binary repair plumbing exclusion contract failed", testRepairSkipsPlumbing)
	contract.RunParallel(t, "binary repair announcement contract failed", testRepairAnnounces)
	contract.RunParallel(t, "binary repair idempotency contract failed", testRepairIdempotent)
	contract.RunParallel(t, "binary repair version-keyed cache contract failed", testRepairVersionKeyed)
	contract.RunParallel(t, "binary repair no-node contract failed", testRepairNoNode)
	contract.RunParallel(t, "binary repair disabled contract failed", testRepairDisabled)
	contract.RunParallel(t, "binary repair torn-cache contract failed", testRepairReplacesTornCache)
	contract.RunParallel(t, "linked manifest repair contract failed", testRepairReadsLinkedManifestWithoutNewline)
}

func testRepairReadsLinkedManifestWithoutNewline(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t, contract.WithSpacePath())
	version := "9.8.7"
	if err := os.Remove(filepath.Join(kit, "package.json")); err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, filepath.Join(kit, "link-manifest.tsv"), "#kit\t"+version)
	registry := newBinaryRepairRegistry(t, version, "#!/bin/sh\nprintf 'manifest:%s\\n' \"$1\"\n")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	out := f.BenchEnv(env, "version")

	out.RequireExit(0)
	if strings.TrimSpace(out.Stdout) != "manifest:version" {
		t.Fatalf("linked manifest repair selected wrong runtime: %q", out.Stdout)
	}
	requireFileExecutable(t, binaryRepairCachePath(t, f, version), "manifest repair did not promote exact cache binary")
}

func testRepairDownloadsAndRuns(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t, contract.WithSpacePath())
	version := "9.8.7"
	binary := "#!/bin/sh\necho repaired:$1\n"
	registry := newBinaryRepairRegistry(t, version, binary)
	env := map[string]string{
		"BENCH_KIT":          kit,
		"BENCH_NPM_REGISTRY": registry.URL,
	}

	out := f.BenchEnv(env, "version")

	out.RequireExit(0)
	if strings.TrimSpace(out.Stdout) != "repaired:version" {
		t.Fatalf("repaired binary did not run\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
	cachePath := binaryRepairCachePath(t, f, version)
	requireFileExecutable(t, cachePath, "repair did not install an executable cache binary")
	out.RequireContains(out.Stderr, "bench: created "+filepath.Dir(cachePath))
	out.RequireContains(out.Stderr, "bench: wrote "+cachePath)
	out.RequireContains(out.Stderr, "bench: wrote "+filepath.Dir(cachePath)+string(os.PathSeparator)+".bench-")
	if got := registry.Hits(); got != 2 {
		t.Fatalf("registry hits = %d, want metadata + tarball", got)
	}
}

func testRepairRefusesBadHash(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	version := "9.8.7"
	registry := newBinaryRepairRegistry(t, version, "#!/bin/sh\necho bad\n", binaryRepairIntegrity("sha512-deadbeef"))
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	out := f.BenchEnv(env, "version")

	out.RequireExit(127)
	out.RequireContains(out.Stderr, "integrity mismatch")
	requireInterim127Remedy(t, out.Stderr)
	requirePathAbsent(t, binaryRepairCachePath(t, f, version), "bad hash installed a binary")
}

func testRepairRefusesMalformedTar(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	version := "9.8.7"
	registry := newBinaryRepairRegistry(t, version, "", binaryRepairTarballBytes(binaryRepairMalformedTarball(t)))
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	out := f.BenchEnv(env, "version")

	out.RequireExit(127)
	out.RequireContains(out.Stderr, "truncated tar entry")
	requireInterim127Remedy(t, out.Stderr)
	requirePathAbsent(t, binaryRepairCachePath(t, f, version), "malformed tar installed a binary")
}

func testRepairOffline(t *testing.T) {
	// WithSpacePath: the printed build remedy must stay quote-safe when the kit path
	// contains spaces (spec row 6 edge inventory); this is the 127-exit test that backs
	// that claim.
	f, kit := binaryRepairFixtureKit(t, contract.WithSpacePath())
	version := "9.8.7"
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": "http://127.0.0.1:1"}

	out := f.BenchEnv(env, "version")

	out.RequireExit(127)
	out.RequireContains(out.Stderr, "repair failed")
	requireInterim127Remedy(t, out.Stderr)
	requirePathAbsent(t, binaryRepairCachePath(t, f, version), "offline repair left a cache binary")
}

func testRepairSkipsPlumbing(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	registry := newBinaryRepairRegistry(t, "9.8.7", "#!/bin/sh\necho should-not-run\n")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	out := f.BenchEnv(env, "tree-hash")

	out.RequireExit(127)
	if got := registry.Hits(); got != 0 {
		t.Fatalf("plumbing command hit registry %d time(s), want 0", got)
	}
}

func testRepairAnnounces(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	version := "9.8.7"
	registry := newBinaryRepairRegistry(t, version, "#!/bin/sh\necho announced\n")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	out := f.BenchEnv(env, "version")

	out.RequireExit(0)
	cachePath := binaryRepairCachePath(t, f, version)
	announcement := "bench: installing @redbench/" + binaryRepairPlatformSuffix(t) + "@" + version + " sha512:"
	out.RequireContains(out.Stderr, announcement)
	if strings.Index(out.Stderr, announcement) > strings.Index(out.Stderr, "bench: created ") {
		t.Fatalf("repair announcement came after installation began\nstderr:\n%s", out.Stderr)
	}
	out.RequireContains(out.Stderr, "bench: created "+filepath.Dir(cachePath))
	out.RequireContains(out.Stderr, "bench: wrote "+filepath.Dir(cachePath)+string(os.PathSeparator)+".bench-")
	out.RequireContains(out.Stderr, "bench: wrote "+cachePath)
}

func testRepairIdempotent(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	version := "9.8.7"
	registry := newBinaryRepairRegistry(t, version, "#!/bin/sh\necho cached:$1\n")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	first := f.BenchEnv(env, "version")
	first.RequireExit(0)
	registry.ResetHits()
	second := f.BenchEnv(env, "version")

	second.RequireExit(0)
	if strings.TrimSpace(second.Stdout) != "cached:version" {
		t.Fatalf("cached binary did not run on second invocation\nstdout:\n%s\nstderr:\n%s", second.Stdout, second.Stderr)
	}
	if got := registry.Hits(); got != 0 {
		t.Fatalf("second invocation hit registry %d time(s), want 0", got)
	}
}

func testRepairVersionKeyed(t *testing.T) {
	f, kit := binaryRepairFixtureKitVersion(t, "9.8.7")
	registryA := newBinaryRepairRegistry(t, "9.8.7", "#!/bin/sh\necho v987\n")
	envA := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registryA.URL}
	f.BenchEnv(envA, "version").RequireExit(0)
	oldCache := binaryRepairCachePath(t, f, "9.8.7")
	requireFileExecutable(t, oldCache, "first repair did not populate old version cache")
	contract.WriteFileAbs(t, filepath.Join(kit, "package.json"), `{"name":"benchkit","version":"9.8.8"}`+"\n")
	registryB := newBinaryRepairRegistry(t, "9.8.8", "#!/bin/sh\necho v988\n")
	envB := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registryB.URL}

	out := f.BenchEnv(envB, "version")

	out.RequireExit(0)
	if strings.TrimSpace(out.Stdout) != "v988" {
		t.Fatalf("version-keyed cache did not fetch new version\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
	requireFileExecutable(t, oldCache, "new version removed old cache")
	requireFileExecutable(t, binaryRepairCachePath(t, f, "9.8.8"), "new version cache not populated")
	if got := registryB.Hits(); got != 2 {
		t.Fatalf("new version registry hits = %d, want metadata + tarball", got)
	}
}

func testRepairNoNode(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	fakebin := filepath.Join(f.Root, "fakebin")
	for _, name := range []string{"dirname", "uname", "tr"} {
		target, err := exec.LookPath(name)
		if err != nil {
			t.Fatalf("look up %s: %v", name, err)
		}
		if err := os.MkdirAll(fakebin, 0o755); err != nil {
			t.Fatalf("mkdir fakebin: %v", err)
		}
		if err := os.Symlink(target, filepath.Join(fakebin, name)); err != nil {
			t.Fatalf("symlink %s: %v", name, err)
		}
	}
	env := map[string]string{"BENCH_KIT": kit, "PATH": fakebin}

	out := f.BenchEnv(env, "version")

	out.RequireExit(127)
	out.RequireContains(out.Stderr, "repair skipped because node is not on PATH")
	requireInterim127Remedy(t, out.Stderr)
}

func testRepairDisabled(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	registry := newBinaryRepairRegistry(t, "9.8.7", "#!/bin/sh\necho should-not-run\n")
	env := map[string]string{
		"BENCH_KIT":          kit,
		"BENCH_NPM_REGISTRY": registry.URL,
		"BENCH_NO_REPAIR":    "1",
	}

	out := f.BenchEnv(env, "version")

	out.RequireExit(127)
	out.RequireContains(out.Stderr, "bench: repair disabled by BENCH_NO_REPAIR")
	requireInterim127Remedy(t, out.Stderr)
	if got := registry.Hits(); got != 0 {
		t.Fatalf("disabled repair hit registry %d time(s), want 0", got)
	}
}

func testRepairReplacesTornCache(t *testing.T) {
	f, kit := binaryRepairFixtureKit(t)
	version := "9.8.7"
	cachePath := binaryRepairCachePath(t, f, version)
	contract.WriteFileAbs(t, cachePath, "")
	registry := newBinaryRepairRegistry(t, version, "#!/bin/sh\necho healed\n")
	env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}

	out := f.BenchEnv(env, "version")

	out.RequireExit(0)
	if strings.TrimSpace(out.Stdout) != "healed" {
		t.Fatalf("torn cache was not repaired\nstdout:\n%s\nstderr:\n%s", out.Stdout, out.Stderr)
	}
	requireFileExecutable(t, cachePath, "torn cache replacement is not executable")
}

// requireInterim127Remedy is the single source for the fail-closed 127 contract:
// the missing-binary error points to reinstall or repair rather than a maintainer
// build script that a linked clone does not ship.
func requireInterim127Remedy(t *testing.T, stderr string) {
	t.Helper()
	if !strings.Contains(stderr, "reinstall redbench") || !strings.Contains(stderr, "repair") {
		t.Fatalf("127 error did not name reinstall/repair remedies\nstderr:\n%s", stderr)
	}
	if strings.Contains(stderr, "scripts/go-build.sh") {
		t.Fatalf("127 error named the maintainer-only build script\nstderr:\n%s", stderr)
	}
}

func requireFileExecutable(t testing.TB, path, msg string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
	if info.Size() == 0 || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("%s: mode=%v size=%d", msg, info.Mode().Perm(), info.Size())
	}
}
