package releaseevidence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestOfflineNetworkControlEnvelopeWriter(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Clean(filepath.Join(cwd, "..", ".."))
	dir := t.TempDir()
	payloadFile, evidenceRoot := filepath.Join(dir, "payload.json"), filepath.Join(dir, "evidence")
	payload := map[string]any{"flag": "BENCH_OFFLINE=1", "journeys": []string{"direct", "local-npm", "loopback-registry"}, "operations": []map[string]any{
		{"class": "wrapper_binary_repair", "observed_attempts": 0},
		{"class": "worktree_git_refresh", "observed_attempts": 0},
		{"class": "codex_discovery_subprocess_and_bundled_fallback", "observed_attempts": 0},
		{"class": "openai_models_request", "observed_attempts": 0},
		{"class": "anthropic_models_request", "observed_attempts": 0},
	}}
	data, _ := json.Marshal(payload)
	if err := os.WriteFile(payloadFile, data, 0o644); err != nil {
		t.Fatal(err)
	}
	record := ft87Requirement(t)
	cmd := exec.Command("node", filepath.Join(root, "scripts", "write-producer-envelope.mjs"), root, evidenceRoot, record.Key, "abc123", "0.2.0", "satisfied", payloadFile)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("write envelope: %v\n%s", err, out)
	}
	bytes, err := os.ReadFile(filepath.Join(evidenceRoot, filepath.FromSlash(record.Path)))
	if err != nil {
		t.Fatal(err)
	}
	commit, version := "abc123", "0.2.0"
	if err := validateRequirementBytes(record, bytes, Identity{SourceCommit: &commit, PackageVersion: &version}); err != nil {
		t.Fatalf("valid offline envelope rejected: %v", err)
	}
	var envelope struct {
		Schema, Status string
		Payload        struct {
			Operations []struct {
				Class    string `json:"class"`
				Attempts int    `json:"observed_attempts"`
			} `json:"operations"`
		}
	}
	if err := json.Unmarshal(bytes, &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Schema != record.Schema || envelope.Status != "satisfied" || len(envelope.Payload.Operations) != 5 {
		t.Fatalf("offline envelope shape = %+v", envelope)
	}
	for _, operation := range envelope.Payload.Operations {
		if operation.Attempts != 0 {
			t.Fatalf("passing operation = %+v", operation)
		}
	}
}

func TestOfflineNetworkControlNonzeroSentinelIsNotPass(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Clean(filepath.Join(cwd, "..", ".."))
	dir := t.TempDir()
	payload, evidenceRoot := filepath.Join(dir, "payload.json"), filepath.Join(dir, "evidence")
	if err := os.WriteFile(payload, []byte(`{"operations":[{"class":"wrapper_binary_repair","observed_attempts":1}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	record := ft87Requirement(t)
	cmd := exec.Command("node", filepath.Join(root, "scripts", "write-producer-envelope.mjs"), root, evidenceRoot, record.Key, "abc123", "0.2.0", "auto", payload)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("write nonpass envelope: %v\n%s", err, out)
	}
	var envelope struct {
		Status string `json:"status"`
	}
	data, _ := os.ReadFile(filepath.Join(evidenceRoot, filepath.FromSlash(record.Path)))
	_ = json.Unmarshal(data, &envelope)
	if envelope.Status == "satisfied" || envelope.Status == "" {
		t.Fatalf("nonzero sentinel status = %q", envelope.Status)
	}
}

func ft87Requirement(t testing.TB) Requirement {
	t.Helper()
	for _, record := range Requirements() {
		if record.Key == "public.ft87.offline_network_control" {
			return record
		}
	}
	t.Fatal("FT87 offline requirement is absent")
	return Requirement{}
}

func TestValidateTarballAcceptsBufferedGzipTrailer(t *testing.T) {
	if err := ValidateTarballForTesting(bufferedPackageTarball(t)); err != nil {
		t.Fatalf("valid tarball rejected: %v", err)
	}
}

func TestValidateTarballRejectsCorruptActiveMemberTrailer(t *testing.T) {
	archive := bufferedPackageTarball(t)
	archive[len(archive)-1] ^= 0xff
	if err := ValidateTarballForTesting(archive); err == nil || !strings.Contains(err.Error(), "gzip: invalid checksum") {
		t.Fatalf("corrupt active-member trailer error = %v, want attributed gzip checksum rejection", err)
	}
}

func bufferedPackageTarball(t *testing.T) []byte {
	t.Helper()
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
	return archive.Bytes()
}
