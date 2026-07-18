package contract

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

type TarEntry struct {
	Mode int64
	Data []byte
}

func ReadJSONFile(t testing.TB, path string, dst any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, dst); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
}

func ReadTarball(t testing.TB, path string) map[string]TarEntry {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	entries := map[string]TarEntry{}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		name := strings.TrimPrefix(h.Name, "./")
		if strings.Contains(name, "\\") || strings.HasPrefix(name, "/") || strings.Contains(name, "../") || name == "" {
			t.Fatalf("unsafe tarball path %q", h.Name)
		}
		if h.Typeflag == tar.TypeDir {
			continue
		}
		if h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeRegA {
			t.Fatalf("tarball contains special member %s", name)
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		if len(data) == 0 {
			t.Fatalf("tarball contains empty member %s", name)
		}
		if _, exists := entries[name]; exists {
			t.Fatalf("tarball contains duplicate member %s", name)
		}
		entries[name] = TarEntry{Mode: h.Mode, Data: data}
	}
	return entries
}
