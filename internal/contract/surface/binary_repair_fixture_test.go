package surface

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

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

type binaryRepairRegistry struct {
	*httptest.Server
	hits atomic.Int64
}

type binaryRepairRegistryConfig struct {
	integrity    string
	pin          string
	tgz          []byte
	hangMetadata bool
}

func binaryRepairHangMetadata() binaryRepairRegistryOption {
	return func(cfg *binaryRepairRegistryConfig) { cfg.hangMetadata = true }
}

type binaryRepairRegistryOption func(*binaryRepairRegistryConfig)

func binaryRepairIntegrity(value string) binaryRepairRegistryOption {
	return func(cfg *binaryRepairRegistryConfig) { cfg.integrity = value }
}

func binaryRepairPinIntegrity(value string) binaryRepairRegistryOption {
	return func(cfg *binaryRepairRegistryConfig) { cfg.pin = value }
}

func binaryRepairTarballBytes(tgz []byte) binaryRepairRegistryOption {
	return func(cfg *binaryRepairRegistryConfig) { cfg.tgz = tgz }
}

func newBinaryRepairRegistry(t testing.TB, kit, version, binary string, opts ...binaryRepairRegistryOption) *binaryRepairRegistry {
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
	pin := integrity
	if cfg.integrity != "" {
		integrity = cfg.integrity
	}
	if cfg.pin != "" {
		pin = cfg.pin
	}
	manifest := map[string]any{"schema_version": 1, "pins": []map[string]string{{
		"name": "@redbench/" + binaryRepairPlatformSuffix(t), "version": version, "integrity": pin,
	}}}
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, filepath.Join(kit, "binary-pins.json"), string(manifestBytes)+"\n")
	reg := &binaryRepairRegistry{}
	mux := http.NewServeMux()
	mux.HandleFunc("/@redbench%2f"+binaryRepairPlatformSuffix(t), func(w http.ResponseWriter, r *http.Request) {
		reg.hits.Add(1)
		if cfg.hangMetadata {
			<-r.Context().Done()
			return
		}
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
	// Existing repair cases opt into the implicit path. Explicit-default tests
	// override this to the empty string so the same fixture covers both postures.
	f.Env["BENCH_REPAIR"] = "1"
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

func binaryRepairGunzip(t testing.TB, tgz []byte) []byte {
	t.Helper()
	r, err := gzip.NewReader(bytes.NewReader(tgz))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func binaryRepairGzipBomb(t testing.TB, expandedSize int) []byte {
	t.Helper()
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	chunk := make([]byte, 32*1024)
	for remaining := expandedSize; remaining > 0; {
		writeSize := min(remaining, len(chunk))
		if _, err := writer.Write(chunk[:writeSize]); err != nil {
			t.Fatal(err)
		}
		remaining -= writeSize
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return compressed.Bytes()
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
		capability.Environment(t, fmt.Sprintf("launcher platform_pkg unavailable: %v\n%s", err, out))
	}
	return strings.TrimPrefix(strings.TrimSpace(string(out)), "@redbench/")
}

func testRepairOptInExact(t *testing.T) {
	for _, value := range []string{"1", "", "0", "true"} {
		t.Run("value="+value, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho repaired\n")
			out := f.BenchEnv(map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL, "BENCH_REPAIR": value}, "version")
			if value == "1" {
				out.RequireExit(0)
				return
			}
			out.RequireExit(127)
			if registry.Hits() != 0 {
				t.Fatalf("BENCH_REPAIR=%q started repair", value)
			}
		})
	}
}

func testRepairSuppressionPrecedence(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		message string
	}{
		{name: "offline", env: map[string]string{"BENCH_OFFLINE": "1"}, message: "repair suppressed by BENCH_OFFLINE=1"},
		{name: "no-repair", env: map[string]string{"BENCH_NO_REPAIR": "1"}, message: "repair disabled by BENCH_NO_REPAIR"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho should-not-run\n")
			tc.env["BENCH_KIT"], tc.env["BENCH_NPM_REGISTRY"], tc.env["BENCH_REPAIR"] = kit, registry.URL, "1"
			out := f.BenchEnv(tc.env, "version")
			out.RequireExit(127)
			out.RequireContains(out.Stderr, tc.message)
			if registry.Hits() != 0 {
				t.Fatalf("suppressed repair hit registry")
			}
		})
	}
}

func testRepairBenchOfflineExact(t *testing.T) {
	for _, value := range []string{"1", "", "0", "true"} {
		t.Run("value="+value, func(t *testing.T) {
			f, kit := binaryRepairFixtureKit(t)
			registry := newBinaryRepairRegistry(t, kit, "9.8.7", "#!/bin/sh\necho repaired\n")
			env := map[string]string{"BENCH_KIT": kit, "BENCH_NPM_REGISTRY": registry.URL}
			if value != "" {
				env["BENCH_OFFLINE"] = value
			}
			out := f.BenchEnv(env, "version")
			if value == "1" {
				out.RequireExit(127)
				out.RequireContains(out.Stderr, "repair suppressed by BENCH_OFFLINE=1")
				if registry.Hits() != 0 {
					t.Fatalf("offline repair hit registry")
				}
				return
			}
			out.RequireExit(0)
			if registry.Hits() == 0 {
				t.Fatalf("BENCH_OFFLINE=%q incorrectly enabled offline mode", value)
			}
		})
	}
}
