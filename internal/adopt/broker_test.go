package adopt

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestWriteBrokerManifestBindsTheRunningBrokerIdentity pins the published binding:
// the manifest lands beside the wrapper and binds the running executable's path, the
// given version, this host's platform, and the executable's content digest.
func TestWriteBrokerManifestBindsTheRunningBrokerIdentity(t *testing.T) {
	dir := t.TempDir()
	wrapper := filepath.Join(dir, "bin", "bench.sh")
	if err := os.MkdirAll(filepath.Dir(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_WRAPPER", wrapper)

	path, broker, err := WriteBrokerManifest("1.2.3")
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, "bin", BrokerManifestName) {
		t.Fatalf("manifest path = %s, want beside the wrapper", path)
	}
	fields, err := ReadBrokerManifest(path)
	if err != nil {
		t.Fatal(err)
	}
	if fields["path"] != broker || !filepath.IsAbs(broker) {
		t.Fatalf("manifest path field = %q, want the absolute broker %q", fields["path"], broker)
	}
	running, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if resolved, err := filepath.EvalSymlinks(running); err == nil {
		running = resolved
	}
	if broker != running {
		t.Fatalf("bound broker = %q, want the running executable %q", broker, running)
	}
	if fields["version"] != "1.2.3" {
		t.Fatalf("manifest version = %q", fields["version"])
	}
	if !regexp.MustCompile(`^(linux|darwin)-(x64|arm64)$`).MatchString(fields["platform"]) {
		t.Fatalf("manifest platform = %q", fields["platform"])
	}
	digest, err := fingerprintPath(broker)
	if err != nil {
		t.Fatal(err)
	}
	if fields["sha256"] != digest {
		t.Fatalf("manifest digest = %q, want %q", fields["sha256"], digest)
	}
}
