package surface

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/contract"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
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

type binaryRepairRegistry struct {
	*httptest.Server
	hits atomic.Int64
}

type binaryRepairRegistryConfig struct {
	integrity string
	tgz       []byte
}

type binaryRepairRegistryOption func(*binaryRepairRegistryConfig)

func binaryRepairIntegrity(value string) binaryRepairRegistryOption {
	return func(cfg *binaryRepairRegistryConfig) { cfg.integrity = value }
}

func binaryRepairTarballBytes(tgz []byte) binaryRepairRegistryOption {
	return func(cfg *binaryRepairRegistryConfig) { cfg.tgz = tgz }
}

func newBinaryRepairRegistry(t testing.TB, version, binary string, opts ...binaryRepairRegistryOption) *binaryRepairRegistry {
	t.Helper()
	cfg := binaryRepairRegistryConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	tgz := cfg.tgz
	if tgz == nil {
		tgz = binaryRepairTarball(t, binary)
	}
	sum := sha512.Sum512(tgz)
	integrity := "sha512-" + base64.StdEncoding.EncodeToString(sum[:])
	if cfg.integrity != "" {
		integrity = cfg.integrity
	}
	reg := &binaryRepairRegistry{}
	mux := http.NewServeMux()
	mux.HandleFunc("/@redbench%2f"+binaryRepairPlatformSuffix(t), func(w http.ResponseWriter, r *http.Request) {
		reg.hits.Add(1)
		meta := map[string]any{
			"versions": map[string]any{
				version: map[string]any{
					"dist": map[string]string{
						"tarball":   reg.URL + "/bench.tgz",
						"integrity": integrity,
					},
				},
			},
		}
		if err := json.NewEncoder(w).Encode(meta); err != nil {
			t.Fatalf("encode registry metadata: %v", err)
		}
	})
	mux.HandleFunc("/bench.tgz", func(w http.ResponseWriter, r *http.Request) {
		reg.hits.Add(1)
		_, _ = w.Write(tgz)
	})
	reg.Server = httptest.NewServer(mux)
	t.Cleanup(reg.Close)
	return reg
}

func (r *binaryRepairRegistry) Hits() int {
	return int(r.hits.Load())
}

func (r *binaryRepairRegistry) ResetHits() {
	r.hits.Store(0)
}

func binaryRepairFixtureKit(t testing.TB, opts ...contract.FixtureOption) (contract.Fixture, string) {
	t.Helper()
	return binaryRepairFixtureKitVersion(t, "9.8.7", opts...)
}

func binaryRepairFixtureKitVersion(t testing.TB, version string, opts ...contract.FixtureOption) (contract.Fixture, string) {
	t.Helper()
	f := contract.NewFixture(t, opts...)
	kit := filepath.Join(f.Root, "kit")
	goRoutingCopyTree(t, filepath.Join(contract.SubjectRoot(t), "bin"), filepath.Join(kit, "bin"))
	contract.WriteFileAbs(t, filepath.Join(kit, "package.json"), `{"name":"benchkit","version":"`+version+`"}`+"\n")
	return f, kit
}

func binaryRepairTarball(t testing.TB, binary string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	data := []byte(binary)
	header := &tar.Header{
		Name: "package/bin/bench",
		Mode: 0o755,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("write tar header: %v", err)
	}
	if _, err := tw.Write(data); err != nil {
		t.Fatalf("write tar payload: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

func binaryRepairMalformedTarball(t testing.TB) []byte {
	t.Helper()
	var tarBuf bytes.Buffer
	header := make([]byte, 512)
	copy(header[0:], "package/bin/bench")
	copy(header[100:], []byte("0000755\x00"))
	copy(header[124:], []byte("00000000020\x00"))
	header[156] = '0'
	tarBuf.Write(header)
	tarBuf.WriteString("#!/bin/sh\n")
	var gzBuf bytes.Buffer
	gz := gzip.NewWriter(&gzBuf)
	if _, err := gz.Write(tarBuf.Bytes()); err != nil {
		t.Fatalf("write malformed tar gzip: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close malformed tar gzip: %v", err)
	}
	return gzBuf.Bytes()
}

func binaryRepairCachePath(t testing.TB, f contract.Fixture, version string) string {
	t.Helper()
	return filepath.Join(f.Env["BENCH_HOME"], "cache", "bin", version, binaryRepairPlatformSuffix(t), "bench")
}

func binaryRepairPlatformSuffix(t testing.TB) string {
	t.Helper()
	script := fmt.Sprintf(
		"source <(sed '/^case \"${1-help}\" in/,$d' %q); platform_pkg",
		filepath.Join(contract.KitRoot(t), "bin", "bench.sh"),
	)
	out, err := exec.Command("bash", "-c", script).CombinedOutput()
	if err != nil {
		t.Skipf("launcher platform_pkg unavailable: %v\n%s", err, out)
	}
	return strings.TrimPrefix(strings.TrimSpace(string(out)), "@redbench/")
}

// requireInterim127Remedy is the single source for the pre-publish 127 contract:
// the missing-binary error prints the clone/git build remedy and no npm line (the
// npm remedy would name unpublished @redbench/* packages; it rejoins publish-day).
func requireInterim127Remedy(t *testing.T, stderr string) {
	t.Helper()
	if !strings.Contains(stderr, "scripts/go-build.sh") {
		t.Fatalf("127 error did not name the build remedy scripts/go-build.sh\nstderr:\n%s", stderr)
	}
	if strings.Contains(stderr, "npm install") {
		t.Fatalf("127 error named an npm install remedy before publish\nstderr:\n%s", stderr)
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
