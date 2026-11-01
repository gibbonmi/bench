package releaseevidence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/jsonfile/jsonfiletest"
)

func TestOfflineNetworkControlEnvelopeWriter(t *testing.T) {
	root := repositoryRoot(t)
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

func TestBuildReleaseEvidenceIncludesRegisteredPackageEvidence(t *testing.T) {
	root := repositoryRoot(t)
	staging := t.TempDir()
	wrapper, packages := filepath.Join(staging, "wrapper"), filepath.Join(staging, "packages")
	cmd := exec.Command("node", filepath.Join(root, "scripts", "build-release-evidence.mjs"), root, wrapper, packages)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build release evidence: %v\n%s", err, out)
	}
	for _, requirement := range Requirements() {
		if requirement.PackageMode == "" {
			continue
		}
		path := filepath.Join(wrapper, filepath.FromSlash(requirement.Path))
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("registered package evidence %q: %v", requirement.Path, err)
		}
		if !info.Mode().IsRegular() || info.Size() == 0 {
			t.Fatalf("registered package evidence %q is not a non-empty regular file", requirement.Path)
		}
	}
}

func TestOfflineNetworkControlNonzeroSentinelIsNotPass(t *testing.T) {
	root := repositoryRoot(t)
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
	return requirementByKey(t, "public.ft87.offline_network_control")
}

func requirementByKey(t testing.TB, key string) Requirement {
	t.Helper()
	for _, record := range Requirements() {
		if record.Key == key {
			return record
		}
	}
	t.Fatalf("requirement %s is absent", key)
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
	return packageTarball(t, packageMember{"payload", payload}, packageMember{requirements.ComponentManifest.Path, manifest})
}

// packageMember is one regular file of a fixture package archive. Its name is relative to the
// package directory.
type packageMember struct {
	name string
	data []byte
}

// fixtureMode is the mode of every fixture archive member.
const fixtureMode = 0o644

// packageTarball writes members, in order, as one gzip tar package archive.
func packageTarball(t *testing.T, members ...packageMember) []byte {
	t.Helper()
	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)
	for _, member := range members {
		if err := tw.WriteHeader(&tar.Header{Name: "package/" + member.name, Mode: fixtureMode, Size: int64(len(member.data)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(member.data); err != nil {
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

// fixtureVersion is the package version of the fixture wrapper package.
const fixtureVersion = "1.0.0"

// identityPackage is the package.json identity of the fixture wrapper package.
const identityPackage = `{"name":"redbench","version":"` + fixtureVersion + `"}`

// componentManifestJSON is the wrapper component manifest that inventories files. Its keys are
// the registry's field labels.
func componentManifestJSON(t testing.TB, files map[string][]byte) string {
	t.Helper()
	labels := requirements.ComponentManifest
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	inventory := make([]map[string]any, 0, len(paths))
	for _, path := range paths {
		inventory = append(inventory, map[string]any{labels.FileFields.Path: path, labels.FileFields.Mode: fmt.Sprintf("%o", fixtureMode), labels.FileFields.Size: len(files[path]), labels.FileFields.SHA256: digest(files[path])})
	}
	target := map[string]any{labels.TargetFields.OS: "all", labels.TargetFields.Arch: "all"}
	component := map[string]any{labels.ComponentFields.Name: "redbench", labels.ComponentFields.Version: fixtureVersion, labels.ComponentFields.Target: target}
	data, err := json.Marshal(map[string]any{labels.RootFields.SchemaVersion: 1, labels.RootFields.Component: component, labels.RootFields.Files: inventory})
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// identityTarball is a wrapper package archive whose package.json is pkg. Its component manifest
// inventories that file.
func identityTarball(t *testing.T, pkg string) []byte {
	t.Helper()
	files := map[string][]byte{"package.json": []byte(pkg)}
	return packageTarball(t, packageMember{"package.json", files["package.json"]}, packageMember{requirements.ComponentManifest.Path, []byte(componentManifestJSON(t, files))})
}

// withoutPackageEvidence is the requirement registry without package evidence records. A
// fixture package then needs only its identity and its component manifest.
func withoutPackageEvidence() Registry {
	value := RequirementsRegistry()
	value.Records = nil
	for _, record := range RequirementsRegistry().Records {
		if record.PackageMode == "" {
			value.Records = append(value.Records, record)
		}
	}
	return value
}

func TestPackageIdentityJSONCompatibility(t *testing.T) {
	restore := SetRequirementsForTesting(withoutPackageEvidence())
	defer restore()
	validate := func(t *testing.T, pkg string) error {
		files, manifest, err := readTarball(identityTarball(t, pkg))
		if err != nil {
			t.Fatalf("fixture archive: %v", err)
		}
		return validatePackageEvidence(files, manifest, "wrapper", fixtureVersion)
	}
	identity := strings.TrimSuffix(strings.TrimPrefix(identityPackage, "{"), "}")
	if err := validate(t, identityPackage); err != nil {
		t.Fatalf("fixture package identity: error = %v, want it accepted", err)
	}

	t.Run("unrelated fields", func(t *testing.T) {
		pkg := `{` + identity + `,"description":"fixture","scripts":{"test":"true"},"files":["bin"]}`
		if err := validate(t, pkg); err != nil {
			t.Errorf("package.json with unrelated fields: error = %v, want it accepted", err)
		}
	})
	for _, refused := range []struct{ name, pkg string }{
		{"duplicate identity key", `{"name":"other",` + identity + `}`},
		{"duplicate inside an unrelated field", `{` + identity + `,"scripts":{"test":"a","test":"b"}}`},
	} {
		t.Run(refused.name, func(t *testing.T) {
			if err := validate(t, refused.pkg); err == nil || err.Error() != "package identity is malformed" {
				t.Errorf("package.json %s: error = %v, want the malformed-identity refusal", refused.pkg, err)
			}
		})
	}

	overflow := `{` + identity + `,"size":` + overflowNumber + `}`
	current := map[string]string{}
	for label, pkg := range map[string]string{
		"raw overflow":                overflow,
		"duplicate plus raw overflow": strings.Replace(overflow, "{", `{"name":"other",`, 1),
	} {
		current[label] = resultText(validate(t, pkg))
	}
	var stored releaseJSONResults
	jsonfiletest.ReadBaseline(t, releaseJSONBaseline, &stored)
	t.Run("package baseline", func(t *testing.T) { jsonfiletest.Compare(t, current, stored.Package) })
}
