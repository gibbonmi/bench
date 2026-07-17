package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type fixtureArchiveFile struct {
	data []byte
	mode int64
}

type fixtureManifestFile struct {
	Path   string `json:"path"`
	Mode   string `json:"mode"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

func seedEvidenceFixture(t *testing.T, root string) {
	t.Helper()
	source := projectRoot(t)
	for _, rel := range []string{"LICENSE", "governance/THIRD_PARTY_NOTICES.txt", "governance/sbom.spdx.json", "governance/policies/supported-versions.json", "governance/policies/security-response.json", "governance/policies/dependency-license-change.json", "governance/policies/threat-model.json", "governance/policies/recovery-rollback.json", "governance/policies/support.json"} {
		data, err := os.ReadFile(filepath.Join(source, rel))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	commit := gitFixtureOutput(t, root, "rev-parse", "HEAD")
	for _, requirement := range []struct{ key, owner, path string }{
		{"public.ft88.data_handling", "FT88", "release-evidence/ft88-data-handling.json"},
		{"public.ft87.offline_network_control", "FT87", "release-evidence/ft87-offline-network-control.json"},
		{"bank.ft71.local_event", "FT71", "release-evidence/ft71-local-event.json"},
	} {
		payload := json.RawMessage(`{"fixture":true}`)
		record := struct {
			SchemaVersion int    `json:"schema_version"`
			Key           string `json:"key"`
			Owner         string `json:"owner"`
			Identity      struct {
				SourceCommit   string `json:"source_commit"`
				PackageVersion string `json:"package_version"`
			} `json:"identity"`
			Status  string          `json:"status"`
			Reason  string          `json:"reason"`
			Payload json.RawMessage `json:"payload"`
			Digest  string          `json:"digest"`
		}{SchemaVersion: 1, Key: requirement.key, Owner: requirement.owner, Status: "satisfied", Payload: payload, Digest: sha256Hex(payload)}
		record.Identity.SourceCommit, record.Identity.PackageVersion = commit, "0.2.0"
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, filepath.FromSlash(requirement.path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	artifactDir := filepath.Join(root, "dist", "artifacts")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	governance := map[string][]byte{}
	for _, rel := range []string{"governance/THIRD_PARTY_NOTICES.txt", "governance/sbom.spdx.json", "governance/policies/supported-versions.json", "governance/policies/security-response.json", "governance/policies/dependency-license-change.json", "governance/policies/threat-model.json", "governance/policies/recovery-rollback.json", "governance/policies/support.json"} {
		governance[rel] = mustFixtureFile(t, filepath.Join(root, rel))
	}
	for _, item := range []struct{ name, target string }{{"redbench", "all-all"}, {"redbench-darwin-arm64", "darwin-arm64"}, {"redbench-darwin-x64", "darwin-x64"}, {"redbench-linux-arm64", "linux-arm64"}, {"redbench-linux-x64", "linux-x64"}} {
		archiveName, componentName := item.name, item.name
		if item.target != "all-all" {
			componentName = "@redbench/" + item.target
		}
		files := map[string]fixtureArchiveFile{
			"LICENSE":      {data: mustFixtureFile(t, filepath.Join(root, "LICENSE")), mode: 0o644},
			"package.json": {data: []byte(fmt.Sprintf("{\"name\":\"%s\",\"version\":\"0.2.0\"}\n", componentName)), mode: 0o644},
		}
		for rel, data := range governance {
			files[rel] = fixtureArchiveFile{data: data, mode: 0o644}
		}
		if item.target != "all-all" {
			files["bin/bench"] = fixtureArchiveFile{data: []byte("#!/bin/sh\n"), mode: 0o755}
		}
		paths := make([]string, 0, len(files))
		for rel := range files {
			paths = append(paths, rel)
		}
		sort.Strings(paths)
		manifestFiles := make([]fixtureManifestFile, 0, len(files))
		for _, rel := range paths {
			file := files[rel]
			manifestFiles = append(manifestFiles, fixtureManifestFile{Path: rel, Mode: fmt.Sprintf("%o", file.mode), Size: int64(len(file.data)), SHA256: sha256Hex(file.data)})
		}
		osName, arch := "all", "all"
		if item.target != "all-all" {
			parts := strings.Split(item.target, "-")
			osName, arch = parts[0], parts[1]
		}
		manifest := struct {
			SchemaVersion int `json:"schema_version"`
			Component     struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Target  struct {
					OS   string `json:"os"`
					Arch string `json:"arch"`
				} `json:"target"`
			} `json:"component"`
			Files []fixtureManifestFile `json:"files"`
		}{SchemaVersion: 1, Files: manifestFiles}
		manifest.Component.Name, manifest.Component.Version = componentName, "0.2.0"
		manifest.Component.Target.OS, manifest.Component.Target.Arch = osName, arch
		manifestData, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		files["component-manifest.json"] = fixtureArchiveFile{data: append(manifestData, '\n'), mode: 0o644}
		if err := os.WriteFile(filepath.Join(artifactDir, archiveName+"-0.2.0.tgz"), fixtureTarball(t, files), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func fixtureTarball(t *testing.T, files map[string]fixtureArchiveFile) []byte {
	t.Helper()
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	tarWriter := tar.NewWriter(gz)
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		file := files[path]
		if err := tarWriter.WriteHeader(&tar.Header{Name: "package/" + path, Mode: file.mode, Size: int64(len(file.data)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tarWriter.Write(file.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return compressed.Bytes()
}

func sha256Hex(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func mustFixtureFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func gitFixtureOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}
