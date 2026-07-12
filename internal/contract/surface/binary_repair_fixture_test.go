package surface

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

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
