package preflight

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	for _, evidence := range packageEvidenceRegistry() {
		rel := evidence.Path
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
	for _, requirement := range Requirements() {
		if !requirement.Producer {
			continue
		}
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
		}{SchemaVersion: 1, Key: requirement.Key, Owner: requirement.Owner, Status: "satisfied", Payload: payload, Digest: sha256Hex(payload)}
		record.Identity.SourceCommit, record.Identity.PackageVersion = commit, "0.2.0"
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, filepath.FromSlash(requirement.Path))
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
	for _, evidence := range packageEvidenceRegistry() {
		rel := evidence.Path
		if rel == "LICENSE" {
			continue
		}
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
			if rel == "governance/sbom.spdx.json" {
				data = fixtureSBOM(t, data, componentName)
			}
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

func rewriteFixtureTarball(t *testing.T, path string, mutate func(map[string]fixtureArchiveFile)) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	tarReader := tar.NewReader(reader)
	files := map[string]fixtureArchiveFile{}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(tarReader)
		if err != nil {
			t.Fatal(err)
		}
		files[strings.TrimPrefix(header.Name, "package/")] = fixtureArchiveFile{data: body, mode: header.Mode & 0o777}
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	mutate(files)
	var manifest struct {
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
	}
	if err := json.Unmarshal(files["component-manifest.json"].data, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest.Files = nil
	paths := make([]string, 0, len(files)-1)
	for rel := range files {
		if rel != "component-manifest.json" {
			paths = append(paths, rel)
		}
	}
	sort.Strings(paths)
	for _, rel := range paths {
		file := files[rel]
		manifest.Files = append(manifest.Files, fixtureManifestFile{Path: rel, Mode: fmt.Sprintf("%o", file.mode), Size: int64(len(file.data)), SHA256: sha256Hex(file.data)})
	}
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	files["component-manifest.json"] = fixtureArchiveFile{data: append(manifestData, '\n'), mode: 0o644}
	if err := os.WriteFile(path, fixtureTarball(t, files), 0o644); err != nil {
		t.Fatal(err)
	}
}

func sha256Hex(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }

func fixtureSBOM(t *testing.T, source []byte, name string) []byte {
	t.Helper()
	var document struct {
		SPDXID            string `json:"SPDXID"`
		SPDXVersion       string `json:"SPDXVersion"`
		CreationInfo      any    `json:"creationInfo"`
		DataLicense       string `json:"dataLicense"`
		Name              string `json:"name"`
		DocumentNamespace string `json:"documentNamespace"`
		Packages          []struct {
			SPDXID           string `json:"SPDXID"`
			Name             string `json:"name"`
			VersionInfo      string `json:"versionInfo"`
			DownloadLocation string `json:"downloadLocation"`
			LicenseConcluded string `json:"licenseConcluded"`
			LicenseDeclared  string `json:"licenseDeclared"`
		} `json:"packages"`
		Relationships []struct {
			SPDXElementID      string `json:"spdxElementId"`
			RelationshipType   string `json:"relationshipType"`
			RelatedSPDXElement string `json:"relatedSpdxElement"`
		} `json:"relationships"`
	}
	if err := json.Unmarshal(source, &document); err != nil {
		t.Fatal(err)
	}
	document.Name = name + "-release"
	document.DocumentNamespace = "https://github.com/gibbonmi/bench/releases/sbom/" + name
	document.Packages[0].Name = name
	document.Packages[0].VersionInfo = "0.2.0"
	document.Packages[0].SPDXID = "SPDXRef-Package-" + strings.NewReplacer("@", "-", "/", "-").Replace(name)
	document.Relationships[0].RelatedSPDXElement = document.Packages[0].SPDXID
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

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
