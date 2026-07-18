package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestReleaseEvidenceIsDeterministicBoundAndIdempotent(t *testing.T) {
	root := preflightRepo(t)
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("first verify exit=%d stderr=%s", code, stderr.String())
	}
	indexPath := filepath.Join(root, "dist", "preflight", "release-index.json")
	sumsPath := filepath.Join(root, "dist", "preflight", "SHA256SUMS")
	firstIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	firstSums, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	var index releaseIndex
	if err := json.Unmarshal(firstIndex, &index); err != nil {
		t.Fatal(err)
	}
	goSumData, err := os.ReadFile(filepath.Join(root, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	goSumBound := false
	for _, input := range index.Inputs {
		if input.Path == "go.sum" {
			goSumBound = input.SHA256 == sha256Hex(goSumData)
		}
	}
	if !goSumBound {
		t.Fatal("release index does not bind go.sum")
	}
	assertReleaseIndexRelationships(t, root, index)
	wantVersions := map[string]string{}
	for _, command := range [][]string{{"go", "env", "GOVERSION"}, {"node", "--version"}, {"npm", "--version"}} {
		out, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			t.Fatal(err)
		}
		wantVersions[command[0]] = strings.TrimSpace(string(out))
	}
	if len(index.Toolchains) != 3 {
		t.Fatalf("release index toolchains=%+v, want exact Go, Node, and npm observations", index.Toolchains)
	}
	for _, toolchain := range index.Toolchains {
		if toolchain.Version != wantVersions[toolchain.Name] {
			t.Fatalf("release index toolchain %s version=%q, want %q", toolchain.Name, toolchain.Version, wantVersions[toolchain.Name])
		}
		if (toolchain.Name == "go" || toolchain.Name == "npm") && len(toolchain.Flags) == 0 {
			t.Fatalf("release index toolchain %s omits release flags", toolchain.Name)
		}
	}
	sumByName := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(firstSums)), "\n") {
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			t.Fatalf("malformed checksum line %q", line)
		}
		sumByName[parts[1]] = parts[0]
	}
	if len(sumByName) != len(index.Artifacts) {
		t.Fatalf("checksum count=%d artifacts=%d", len(sumByName), len(index.Artifacts))
	}
	for _, artifact := range index.Artifacts {
		data, err := os.ReadFile(filepath.Join(root, "dist", "artifacts", artifact.Name))
		if err != nil {
			t.Fatal(err)
		}
		if artifact.SHA256 != sha256Hex(data) || sumByName[artifact.Name] != artifact.SHA256 {
			t.Fatalf("artifact digest binding failed for %s", artifact.Name)
		}
	}
	artifactDir := filepath.Join(root, "dist", "artifacts")
	staging := t.TempDir()
	entries, err := os.ReadDir(artifactDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if err := os.Rename(filepath.Join(artifactDir, entry.Name()), filepath.Join(staging, entry.Name())); err != nil {
			t.Fatal(err)
		}
	}
	for i := len(entries) - 1; i >= 0; i-- {
		if err := os.Rename(filepath.Join(staging, entries[i].Name()), filepath.Join(artifactDir, entries[i].Name())); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("LANG", "C")
	t.Setenv("TZ", "UTC")
	for key, directory := range map[string]string{
		"HOME":             filepath.Join(t.TempDir(), "isolated home"),
		"npm_config_cache": filepath.Join(t.TempDir(), "isolated npm cache"),
		"GOCACHE":          filepath.Join(t.TempDir(), "isolated go cache"),
		"GOMODCACHE":       filepath.Join(t.TempDir(), "isolated module cache"),
		"TMPDIR":           filepath.Join(t.TempDir(), "isolated tmp root"),
	} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Setenv(key, directory)
	}
	stderr.Reset()
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("second verify exit=%d stderr=%s", code, stderr.String())
	}
	secondIndex, _ := os.ReadFile(indexPath)
	secondSums, _ := os.ReadFile(sumsPath)
	if string(secondIndex) != string(firstIndex) || string(secondSums) != string(firstSums) {
		t.Fatal("release evidence changed with enumeration order or environment")
	}
	if err := os.WriteFile(filepath.Join(root, "dist", "native-proofs", "unexpected.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("unexpected native proof exit=%d stderr=%s", code, stderr.String())
	}
	if err := os.Remove(filepath.Join(root, "dist", "native-proofs", "unexpected.json")); err != nil {
		t.Fatal(err)
	}

	failing := filepath.Join(root, "fail-gate")
	if err := os.WriteFile(failing, []byte("#!/bin/sh\nexit 9\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	greenGate := os.Getenv("BENCH_PREFLIGHT_GATE")
	if err := os.Setenv("BENCH_PREFLIGHT_GATE", failing); err != nil {
		t.Fatal(err)
	}
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("red rerun exit=%d stderr=%s", code, stderr.String())
	}
	redIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	var red releaseIndex
	if err := json.Unmarshal(redIndex, &red); err != nil || red.Status != StatusRed {
		t.Fatalf("red rerun index=%s err=%v", redIndex, err)
	}
	if err := os.Setenv("BENCH_PREFLIGHT_GATE", greenGate); err != nil {
		t.Fatal(err)
	}
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("green-after-red exit=%d stderr=%s", code, stderr.String())
	}
	var final releaseIndex
	finalIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(finalIndex, &final); err != nil || final.Status != StatusGreen {
		t.Fatalf("green-after-red index=%s err=%v", finalIndex, err)
	}
	files, err := os.ReadDir(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(PhaseNames(ModeVerify))+3 {
		t.Fatalf("promoted evidence file count=%d, want %d", len(files), len(PhaseNames(ModeVerify))+3)
	}
}

func assertReleaseIndexRelationships(t *testing.T, root string, index releaseIndex) {
	t.Helper()
	commit := gitFixtureOutput(t, root, "rev-parse", "HEAD")
	if index.Identity.SourceCommit != commit || index.Identity.PackageVersion != "0.2.0" || index.Identity.BinaryVersion != "" || index.Identity.Tag != "" || index.Identity.ChangelogHeading != "" || index.Identity.Toolchain != "1.25.0" {
		t.Fatalf("release identity does not bind fixture facts: %+v", index.Identity)
	}
	wantTargets := []string{"darwin/arm64/darwin/arm64/macos-15", "darwin/x64/darwin/amd64/macos-15-intel", "linux/arm64/linux/arm64/ubuntu-24.04-arm", "linux/x64/linux/amd64/ubuntu-24.04"}
	gotTargets := make([]string, 0, len(index.Targets))
	for _, target := range index.Targets {
		gotTargets = append(gotTargets, fmt.Sprintf("%s/%s/%s/%s/%s", target.OS, target.Arch, target.GOOS, target.GOArch, target.Runner))
	}
	if !reflect.DeepEqual(gotTargets, wantTargets) {
		t.Fatalf("release target bindings = %v, want %v", gotTargets, wantTargets)
	}

	artifactDir := filepath.Join(root, "dist", "artifacts")
	for _, artifact := range index.Artifacts {
		data := mustFixtureFile(t, filepath.Join(artifactDir, artifact.Name))
		if strings.HasSuffix(artifact.Name, ".tar.gz") {
			if artifact.Size != int64(len(data)) || artifact.SHA256 != sha256Hex(data) || artifact.ComponentDigest == "" || artifact.SBOMDigest == "" || artifact.InventoryDigest == "" {
				t.Fatalf("offline artifact evidence relationship is incomplete for %s: %+v", artifact.Name, artifact)
			}
			continue
		}
		files := independentlyReadArchive(t, data)
		manifestBytes := files["component-manifest.json"]
		var manifest struct {
			Files []fixtureManifestFile `json:"files"`
		}
		if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
			t.Fatal(err)
		}
		inventory := independentInventoryBytes(t, manifest.Files)
		if artifact.Size != int64(len(data)) || artifact.SHA256 != sha256Hex(data) || artifact.ComponentDigest != sha256Hex(manifestBytes) || artifact.SBOMDigest != sha256Hex(files["governance/sbom.spdx.json"]) || artifact.InventoryDigest != sha256Hex(inventory) {
			t.Fatalf("artifact/component/SBOM/inventory relationship is incomplete for %s: %+v", artifact.Name, artifact)
		}
	}

	wantInputs := map[string]string{}
	for _, path := range []string{"LICENSE", "go.mod", "go.sum", "governance/THIRD_PARTY_NOTICES.txt", "governance/policies/dependency-license-change.json", "governance/policies/recovery-rollback.json", "governance/policies/security-response.json", "governance/policies/support.json", "governance/policies/supported-versions.json", "governance/policies/threat-model.json", "governance/sbom.spdx.json", "internal/releaseevidence/registry.json", "internal/releaseevidence/requirements.json", "release-evidence/ft87-offline-network-control.json", "release-evidence/ft88-data-handling.json", ".bench/gate.sh", "bin/bench.sh", "package.json", "scripts/build-artifacts.sh", "scripts/build-offline-archives.sh", "scripts/assemble-offline-archive.mjs", "scripts/release-plan.mjs", "scripts/release-plan.json", "scripts/write-deterministic-archive.mjs", "scripts/compare-artifacts.sh", "scripts/native-proof.sh", "scripts/aggregate-native-proofs.sh", "scripts/build-release-evidence.mjs", "scripts/go-build.sh", "scripts/smoke-artifacts.sh", "scripts/smoke-offline.sh", "scripts/offline-registry.mjs", "scripts/wrapper-assets.json"} {
		wantInputs[path] = sha256Hex(mustFixtureFile(t, filepath.Join(root, filepath.FromSlash(path))))
	}
	gotInputs := map[string]string{}
	for _, input := range index.Inputs {
		gotInputs[input.Path] = input.SHA256
	}
	if !reflect.DeepEqual(gotInputs, wantInputs) {
		t.Fatalf("policy/input bindings differ\ngot: %v\nwant: %v", gotInputs, wantInputs)
	}

	wantRequirements := map[string]string{}
	for _, requirement := range Requirements() {
		status, reason := "satisfied", ""
		if requirement.Key == "bank.ft71.local_event" {
			status, reason = "not_applicable", "requirement is not applicable to selected profile"
		}
		digest := ""
		if status == "satisfied" {
			digest = sha256Hex(mustFixtureFile(t, filepath.Join(root, filepath.FromSlash(requirement.Path))))
		}
		wantRequirements[requirement.Key] = fmt.Sprintf("%s|%s|%s|%s|%t|%s|%s", requirement.Owner, requirement.Schema, requirement.Requiredness, status, status == "satisfied", reason, digest)
	}
	gotRequirements := map[string]string{}
	for _, status := range index.Requirements {
		gotRequirements[status.Key] = fmt.Sprintf("%s|%s|%s|%s|%t|%s|%s", status.Owner, status.Schema, status.Requiredness, status.Status, status.Applicable, status.Reason, status.Digest)
	}
	if !reflect.DeepEqual(gotRequirements, wantRequirements) {
		t.Fatalf("requirement bindings differ\ngot: %v\nwant: %v", gotRequirements, wantRequirements)
	}

	wantPhases := []string{"gate", "race", "vet", "vulnerability", "artifacts", "smoke"}
	for i, phase := range index.Phases {
		if i >= len(wantPhases) || phase.Name != wantPhases[i] || phase.Status != StatusGreen {
			t.Fatalf("phase ordering/status differs: %+v", index.Phases)
		}
		record := fmt.Sprintf("{\"schema_version\":1,\"phase\":%q,\"mode\":\"verify\",\"status\":\"green\",\"exit_code\":0,\"error\":null}\n", phase.Name)
		if phase.Digest != sha256Hex([]byte(record)) {
			t.Fatalf("phase %s record digest is not bound", phase.Name)
		}
	}
	if len(index.Phases) != len(wantPhases) {
		t.Fatalf("phase bindings = %+v", index.Phases)
	}
	wantFlags := []string{"gate={root}/bin/bench.sh gate", "race=go test -race -count=1 ./...", "vet=go vet ./...", "vulnerability=", "artifacts=bash {root}/scripts/build-artifacts.sh {root} {root}/dist/artifacts", "smoke=bash {root}/scripts/smoke-artifacts.sh {root}/dist/artifacts"}
	if !reflect.DeepEqual(index.Flags, wantFlags) {
		t.Fatalf("phase flags = %v, want %v", index.Flags, wantFlags)
	}
}

func independentlyReadArchive(t *testing.T, data []byte) map[string][]byte {
	t.Helper()
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	reader := tar.NewReader(gz)
	files := map[string][]byte{}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		files[strings.TrimPrefix(header.Name, "package/")] = body
	}
	return files
}

func independentInventoryBytes(t *testing.T, files []fixtureManifestFile) []byte {
	t.Helper()
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	var out bytes.Buffer
	out.WriteByte('[')
	for i, file := range files {
		if i > 0 {
			out.WriteByte(',')
		}
		pathJSON, _ := json.Marshal(file.Path)
		modeJSON, _ := json.Marshal(file.Mode)
		digestJSON, _ := json.Marshal(file.SHA256)
		fmt.Fprintf(&out, "{\"mode\":%s,\"path\":%s,\"sha256\":%s,\"size\":%d}", modeJSON, pathJSON, digestJSON, file.Size)
	}
	out.WriteString("]\n")
	return out.Bytes()
}

func TestReleasePolicyFailureClassesAreRed(t *testing.T) {
	root := preflightRepo(t)
	r := &runner{root: root, mode: ModePublish, binaryVersion: "0.2.0", stderr: &bytes.Buffer{}}
	if err := r.populateBaseIdentity(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_REF", "refs/tags/v0.2.0-beta.1")
	if err := r.checkIdentity(context.Background()); err == nil {
		t.Fatal("prerelease tag passed")
	}
	if err := r.checkAncestry(context.Background()); err == nil {
		t.Fatal("missing origin/main ancestry passed")
	}
	tag := "v0.2.0"
	r.identity.Tag = &tag
	for name, body := range map[string]string{
		"duplicate":          "## [0.2.0] - 2026-07-16\n## [0.2.0] - 2026-07-17\n",
		"legacy":             "## v0.2.0 (2026-07-16)\n",
		"invalid date":       "## [Unreleased]\n\n## [0.2.0] - 2026-02-30\n",
		"missing unreleased": "## [0.2.0] - 2026-07-16\n",
	} {
		if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := r.checkChangelog(); err == nil {
			t.Fatalf("%s changelog passed", name)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ntoolchain go1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readToolchain(root); err == nil {
		t.Fatal("non-patch toolchain passed")
	}
}
