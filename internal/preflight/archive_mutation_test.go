package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuiltCommandRejectsConcatenatedGzipArchiveAndPreservesPriorGeneration(t *testing.T) {
	binary := buildPreflightBinary(t)
	root := preflightRepo(t)
	prior := priorGeneration(t, binary, root)
	archive := filepath.Join(root, "dist", "artifacts", "redbench-0.2.0-darwin-arm64.tar.gz")
	data, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	var trailing bytes.Buffer
	gz := gzip.NewWriter(&trailing)
	if _, err := gz.Write([]byte("not an archive\n")); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archive, append(data, trailing.Bytes()...), 0o644); err != nil {
		t.Fatal(err)
	}
	output, err := runBuilt(t, binary, root, os.Environ())
	if err == nil || !strings.Contains(string(output), "concatenated gzip members") {
		t.Fatalf("error=%v\n%s", err, output)
	}
	assertPriorGeneration(t, root, prior)
}

func TestBuiltCommandRejectsHostilePackageArchivesAndPreservesPriorGeneration(t *testing.T) {
	binary := buildPreflightBinary(t)
	for _, test := range []struct {
		name, want string
		headers    []*tar.Header
	}{
		{"duplicate", "archive contains duplicate path", []*tar.Header{{Name: "package/duplicate", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg}, {Name: "package/duplicate", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg}}},
		{"traversal", "unsafe archive path", []*tar.Header{{Name: "package/../escape", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg}}},
		{"special", "archive contains special file", []*tar.Header{{Name: "package/link", Mode: 0o777, Typeflag: tar.TypeSymlink, Linkname: "target"}}},
		{"unsafe-mode", "archive contains unsafe mode", []*tar.Header{{Name: "package/private", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg}}},
		{"empty", "archive contains empty member", []*tar.Header{{Name: "package/empty", Mode: 0o644, Size: 0, Typeflag: tar.TypeReg}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			assertHostilePackageArchive(t, binary, test.headers, test.want)
		})
	}
}

func assertHostilePackageArchive(t *testing.T, binary string, headers []*tar.Header, want string) {
	t.Helper()
	root := preflightRepo(t)
	prior := priorGeneration(t, binary, root)
	archive := filepath.Join(root, "dist", "artifacts", "redbench-darwin-arm64-0.2.0.tgz")
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	writer := tar.NewWriter(gz)
	for _, header := range headers {
		if err := writer.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if header.Size > 0 {
			if _, err := writer.Write(bytes.Repeat([]byte("x"), int(header.Size))); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(archive, body.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	output, err := runBuilt(t, binary, root, os.Environ())
	if err == nil || !strings.Contains(string(output), want) {
		t.Fatalf("error=%v\n%s", err, output)
	}
	assertPriorGeneration(t, root, prior)
}
