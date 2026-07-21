package releaseevidence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"testing"
)

func TestValidateTarballAcceptsBufferedGzipTrailer(t *testing.T) {
	// Keep the compressed stream across the gzip reader's internal read boundary,
	// leaving bytes from the active member unread when tar reaches its end markers.
	payload := make([]byte, 95_348)
	state := uint32(1)
	for i := range payload {
		state ^= state << 13
		state ^= state >> 17
		state ^= state << 5
		payload[i] = byte(state)
	}
	manifest := []byte(`{"schema_version":1,"component":{"name":"fixture","version":"1.0.0","target":{"os":"all","arch":"all"}},"files":[]}`)

	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	for _, file := range []struct {
		name string
		data []byte
	}{
		{name: "package/payload", data: payload},
		{name: "package/component-manifest.json", data: manifest},
	} {
		if err := tw.WriteHeader(&tar.Header{Name: file.name, Mode: 0o644, Size: int64(len(file.data)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(file.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	if err := ValidateTarballForTesting(archive.Bytes()); err != nil {
		t.Fatalf("valid tarball rejected: %v", err)
	}
}
