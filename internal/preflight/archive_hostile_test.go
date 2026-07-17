package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"strings"
	"testing"
)

func TestArchiveRejectsMemberLargerThanInspectionLimit(t *testing.T) {
	t.Cleanup(setArchiveMemberLimitForTesting(4))
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "package/oversize", Mode: 0o644, Size: 5, Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte("12345")); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	err := validateTarballForTesting(compressed.Bytes())
	if err == nil || !strings.Contains(err.Error(), "exceeds inspection limit") {
		t.Fatalf("oversize member error = %v", err)
	}
}
