package landing

import (
	"strings"
	"testing"
)

func TestLocalCapturePathRecognizesTheBuildLog(t *testing.T) {
	for _, tc := range []struct {
		path string
		want bool
	}{
		{path: "capture/build-log/", want: true},
		{path: "capture/build-log", want: true},
		{path: "capture/build-log/build-directions.log", want: true},
		{path: "capture/build-log-foreign/x.log"},
		{path: "capture/build-log/../foreign"},
		{path: "capture/build-log//build-directions.log"},
		{path: "capture/build-log/\x1b"},
	} {
		t.Run(strings.ReplaceAll(tc.path, "\x1b", "escape"), func(t *testing.T) {
			if got := LocalCapturePath(tc.path); got != tc.want {
				t.Fatalf("LocalCapturePath(%q) = %t, want %t", tc.path, got, tc.want)
			}
		})
	}
}

func TestCheckoutFingerprintIgnoresTheBuildLog(t *testing.T) {
	before, after, err := ignoredWriteFingerprints(t, "capture/build-log/", "capture/build-log/build-directions.log")
	if err != nil || after != before {
		t.Fatalf("build-log fingerprint = (%q, %v), want %q", after, err, before)
	}
}

// ignoredWriteFingerprints commits a .gitignore that holds the ignore line in a new
// checkout, then writes one file at path. It returns the fingerprint before the write and
// the fingerprint and error after it.
func ignoredWriteFingerprints(t *testing.T, ignore, path string) (before, after string, err error) {
	t.Helper()
	root := fixture(t)
	write(t, root, ".gitignore", ignore+"\n")
	git(t, root, "add", ".gitignore")
	git(t, root, "commit", "-qm", "ignore "+ignore)
	before, err = CheckoutFingerprint(root)
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, path, "ignored\n")
	after, err = CheckoutFingerprint(root)
	return before, after, err
}
