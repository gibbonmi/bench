package conformance

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type workflowTriggerShape struct {
	pullRequest, pushBranches, mainBranch bool
}

type requirementRecord struct {
	Key          string   `json:"key"`
	Owner        string   `json:"owner"`
	Schema       string   `json:"schema"`
	Profiles     []string `json:"profiles"`
	Requiredness string   `json:"requiredness"`
	Path         string   `json:"path"`
	PackageMode  string   `json:"package_mode"`
}

func checkReleasePreflight(root string) []string {
	var diags []string
	data, err := os.ReadFile(filepath.Join(root, "internal", "releaseevidence", "registry.json"))
	if err != nil {
		if os.IsNotExist(err) {
			diags = append(diags, "release preflight registry is absent")
		} else {
			diags = append(diags, "release preflight registry is unreadable")
		}
	}
	registryUsable := err == nil
	var registry struct {
		Verify      []string `json:"verify"`
		PublishOnly []string `json:"publish_only"`
	}
	if registryUsable {
		if decodeErr := json.Unmarshal(data, &registry); decodeErr != nil {
			diags = append(diags, "release preflight registry is malformed")
			registryUsable = false
		}
	}
	requirementData, err := os.ReadFile(filepath.Join(root, "internal", "releaseevidence", "requirements.json"))
	if err != nil {
		if os.IsNotExist(err) {
			diags = append(diags, "release requirement registry is absent")
		} else {
			diags = append(diags, "release requirement registry is unreadable")
		}
	}
	requirementUsable := err == nil
	var requirementRegistry struct {
		SchemaVersion int                 `json:"schema_version"`
		Records       []requirementRecord `json:"records"`
	}
	if requirementUsable {
		if decodeErr := json.Unmarshal(requirementData, &requirementRegistry); decodeErr != nil {
			diags = append(diags, "release requirement registry is malformed")
			requirementUsable = false
		} else if requirementRegistry.SchemaVersion != 1 {
			diags = append(diags, "release requirement registry has unsupported schema version")
			requirementUsable = false
		}
	}
	if err != nil {
		requirementData = nil
	}
	public, bank := map[string]bool{}, map[string]bool{}
	if requirementUsable {
		for _, record := range requirementRegistry.Records {
			if record.Key == "" || record.Owner == "" || record.Schema == "" || record.Path == "" || record.Requiredness != "required" && record.Requiredness != "conditional" || len(record.Profiles) == 0 {
				diags = append(diags, "release requirement registry has incomplete schema: "+record.Key)
				continue
			}
			for _, profile := range record.Profiles {
				if profile == "public" {
					public[record.Key] = true
				}
				if profile == "bank" {
					bank[record.Key] = true
				}
			}
		}
		for key := range public {
			if !bank[key] {
				diags = append(diags, "release requirement registry bank profile is not a strict public superset")
				break
			}
		}
		for _, key := range []string{"public.ft88.data_handling", "public.ft87.offline_network_control", "bank.ft71.local_event"} {
			if !containsKey(requirementRegistry.Records, key) {
				diags = append(diags, "release requirement registry omits "+key)
			}
		}
	}
	if registryUsable && !reflect.DeepEqual(registry.Verify, []string{"gate", "race", "vet", "vulnerability", "artifacts", "smoke"}) {
		diags = append(diags, "release preflight verify registry omits or reorders a required phase class")
	}
	if registryUsable && !reflect.DeepEqual(registry.PublishOnly, []string{"identity", "ancestry", "changelog"}) {
		diags = append(diags, "release preflight publish registry omits or reorders a required phase class")
	}
	native := readIfExists(filepath.Join(root, ".github", "workflows", "native-runtime.yml"))
	release := readIfExists(filepath.Join(root, ".github", "workflows", "release.yml"))
	installer := readIfExists(filepath.Join(root, "scripts", "install-govulncheck.sh"))
	if strings.Count(native, "bash scripts/install-govulncheck.sh") != 1 || strings.Count(release, "bash scripts/install-govulncheck.sh") != 2 || !regexp.MustCompile(`govulncheck@v[0-9]+\.[0-9]+\.[0-9]+`).MatchString(installer) {
		diags = append(diags, "release workflows do not consume the repository-owned govulncheck setup")
	}
	if strings.Contains(native, "govulncheck@") || strings.Contains(release, "govulncheck@") {
		diags = append(diags, "release workflows duplicate the govulncheck version pin")
	}
	if native != "" {
		for message, anchor := range map[string]string{"native verification bypasses full release preflight": "scripts/release-preflight.sh --mode verify\n", "native verification does not upload preflight evidence": "verify-preflight-evidence", "native runner matrix bypasses focused smoke": "--mode verify --phase smoke"} {
			if strings.Count(native, anchor) != 1 {
				diags = append(diags, message)
			}
		}
	}
	if release != "" {
		for message, anchor := range map[string]string{"tag publication bypasses full release preflight": "scripts/release-preflight.sh --mode publish --profile public\n", "tag publication does not upload preflight evidence": "publish-preflight-evidence", "tag runner matrix bypasses focused smoke": "--mode verify --phase smoke", "tag evidence does not request repository maximum retention": "retention-days: ${{ github.retention_days }}", "publication does not wait for preflight and every native proof row": "needs: [preflight, smoke, native-evidence]", "publication does not wait for publish authorization": "needs: [authorize]"} {
			if !strings.Contains(release, anchor) {
				diags = append(diags, message)
			}
		}
	}
	if release != "" {
		platform, wrapper := strings.Index(release, "name: Publish platform packages"), strings.Index(release, "name: Publish wrapper")
		if platform < 0 || wrapper < platform {
			diags = append(diags, "release publication is not platform-first and wrapper-last")
		}
	}
	if !regexp.MustCompile(`(?m)^toolchain go[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(readIfExists(filepath.Join(root, "go.mod"))) {
		diags = append(diags, "release preflight requires an exact Go patch toolchain")
	}
	if registryUsable && requirementUsable {
		diags = append(diags, runReleaseEvidenceProbe(root, requirementRegistry.Records)...)
	}
	archiveEvidence := readIfExists(filepath.Join(root, "internal", "releaseevidence", "artifact_evidence.go"))
	if archiveEvidence != "" && !strings.Contains(archiveEvidence, "approved npm tarball bytes") {
		diags = append(diags, "release evidence does not bind archive digest bytes")
	}
	comparison := readIfExists(filepath.Join(root, "scripts", "compare-artifacts.sh"))
	if comparison != "" && (!strings.Contains(comparison, "cmp -s") || !strings.Contains(comparison, "reproducibility mismatch:")) {
		diags = append(diags, "reproducibility comparator does not require exact byte equality")
	}
	offlineSmoke := readIfExists(filepath.Join(root, "scripts", "smoke-offline.sh"))
	if offlineSmoke != "" && (strings.Count(offlineSmoke, "BENCH_NO_REPAIR=1") != 7 || strings.Count(offlineSmoke, "npm_config_offline=true") != 3) {
		diags = append(diags, "offline smoke permits repair or network fallback")
	}
	if offlineSmoke != "" && (!strings.Contains(offlineSmoke, "printf 'offline smoke: loopback registry fixture did not start\\n' >&2\n  exit 1") || strings.Contains(offlineSmoke, "local-fixture")) {
		diags = append(diags, "offline registry smoke does not fail closed")
	}
	return diags
}

func TestReleasePreflightDiagnosticsDistinguishAndAggregate(t *testing.T) {
	missing := t.TempDir()
	if got := strings.Join(checkReleasePreflight(missing), "\n"); !strings.Contains(got, "release preflight registry is absent") || !strings.Contains(got, "release requirement registry is absent") {
		t.Fatalf("absent diagnostics = %q", got)
	}
	malformed := t.TempDir()
	if err := os.MkdirAll(filepath.Join(malformed, "internal", "releaseevidence"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(malformed, "internal", "releaseevidence", "registry.json"), []byte("{\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(malformed, "internal", "releaseevidence", "requirements.json"), []byte("{\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := strings.Join(checkReleasePreflight(malformed), "\n")
	if !strings.Contains(got, "release preflight registry is malformed") || !strings.Contains(got, "release requirement registry is malformed") {
		t.Fatalf("malformed diagnostics = %q", got)
	}
}

func runReleaseEvidenceProbe(root string, packageEvidence []requirementRecord) []string {
	probeMain := filepath.Join(root, "cmd", "bench", "main.go")
	probeDirCreated := false
	if info, err := os.Stat(probeMain); err == nil {
		if !info.Mode().IsRegular() {
			return []string{"release evidence probe command is not a regular file"}
		}
	} else if os.IsNotExist(err) {
		if info, dirErr := os.Stat(filepath.Dir(probeMain)); dirErr == nil {
			if !info.IsDir() {
				return []string{"release evidence probe command directory is not a directory"}
			}
		} else if os.IsNotExist(dirErr) {
			probeDirCreated = true
		} else {
			return []string{"release evidence probe command directory is unreadable"}
		}
		if err := os.MkdirAll(filepath.Dir(probeMain), 0o755); err != nil {
			return []string{"release evidence probe could not prepare preflight command"}
		}
		const source = `package main

import (
	"os"

	"github.com/gibbonmi/bench/internal/preflight"
)

var version string

func main() {
	if len(os.Args) < 2 || os.Args[1] != "release-preflight" {
		os.Exit(2)
	}
	os.Exit(preflight.Command(os.Args[2:], version, os.Stderr))
}
`
		if err := os.WriteFile(probeMain, []byte(source), 0o644); err != nil {
			return []string{"release evidence probe could not prepare preflight command"}
		}
		defer func() {
			_ = os.Remove(probeMain)
			if probeDirCreated {
				_ = os.Remove(filepath.Dir(probeMain))
				_ = os.Remove(filepath.Dir(filepath.Dir(probeMain)))
			}
		}()
	} else {
		return []string{"release evidence probe command is unreadable"}
	}
	build := exec.Command("bash", filepath.Join(root, "scripts", "build-artifacts.sh"), root, filepath.Join(root, "dist", "artifacts"))
	build.Dir = root
	if output, err := build.CombinedOutput(); err != nil {
		return []string{"release package evidence build failed: " + strings.TrimSpace(string(output))}
	}
	var matrix []struct {
		OS   string `json:"os"`
		Arch string `json:"arch"`
	}
	if err := readJSONAt(root, filepath.Join("scripts", "platforms.json"), &matrix); err != nil {
		return []string{"release package evidence matrix is unreadable"}
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := readJSONAt(root, filepath.Join("package.json"), &pkg); err != nil {
		return []string{"release package evidence package identity is unreadable"}
	}
	artifacts := append([]string{"redbench-" + pkg.Version + ".tgz"}, func() []string {
		out := make([]string, 0, len(matrix))
		for _, target := range matrix {
			out = append(out, fmt.Sprintf("redbench-%s-%s-%s.tgz", target.OS, target.Arch, pkg.Version))
		}
		return out
	}()...)
	for _, name := range artifacts {
		names, err := archiveNames(filepath.Join(root, "dist", "artifacts", name))
		if err != nil {
			return []string{"release package evidence artifact is unreadable: " + name}
		}
		for _, evidence := range packageEvidence {
			if evidence.PackageMode == "" {
				continue
			}
			if !names["package/"+evidence.Path] {
				return []string{"release package evidence allowlist omits " + evidence.Path}
			}
		}
	}
	fake, err := os.CreateTemp("", "bench-preflight-probe-*")
	if err != nil {
		return []string{"release evidence probe could not create phase fixture"}
	}
	fakePath := fake.Name()
	defer os.Remove(fakePath)
	if _, err := fake.WriteString("#!/bin/sh\nprintf '{}\\n'\nexit 0\n"); err != nil || fake.Chmod(0o755) != nil || fake.Close() != nil {
		return []string{"release evidence probe could not prepare phase fixture"}
	}
	env := append([]string{}, os.Environ()...)
	for _, phase := range []string{"gate", "race", "vet", "vulnerability", "artifacts", "smoke"} {
		env = append(env, "BENCH_PREFLIGHT_"+strings.ToUpper(phase)+"="+fakePath)
	}
	preflightBinary := filepath.Join(root, "dist", "bench-preflight")
	if err := os.Remove(preflightBinary); err != nil && !os.IsNotExist(err) {
		return []string{"release evidence probe could not reset the preflight command"}
	}
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "release-preflight.sh"), "--mode", "verify")
	cmd.Dir, cmd.Env = root, env
	if output, err := exec.Command("git", "-C", root, "rev-parse", "--verify", "HEAD").CombinedOutput(); err != nil || strings.TrimSpace(string(output)) == "" {
		commit := exec.Command("git", "-C", root, "-c", "user.name=release-evidence-probe", "-c", "user.email=release-evidence-probe@example.invalid", "commit", "--quiet", "--allow-empty", "-m", "release-evidence-probe")
		if output, err := commit.CombinedOutput(); err != nil {
			return []string{"release evidence probe could not establish source HEAD: " + strings.TrimSpace(string(output))}
		}
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		return []string{"release evidence probe failed: " + strings.TrimSpace(string(output))}
	}
	indexData, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "release-index.json"))
	if err != nil {
		return []string{"release evidence probe did not generate release-index.json"}
	}
	var index struct {
		Artifacts []struct {
			ComponentDigest string `json:"component_manifest_sha256"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(indexData, &index); err != nil {
		return []string{"release evidence probe generated malformed release-index.json"}
	}
	for _, artifact := range index.Artifacts {
		if artifact.ComponentDigest == "" {
			return []string{"release index does not bind component manifest digests"}
		}
	}
	return nil
}

func readJSONAt(root, rel string, value any) error {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func archiveNames(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	reader := tar.NewReader(gz)
	names := map[string]bool{}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return names, nil
		}
		if err != nil {
			return nil, err
		}
		names[header.Name] = true
	}
}

func containsKey(records []requirementRecord, want string) bool {
	for _, record := range records {
		if record.Key == want {
			return true
		}
	}
	return false
}

func nativeWorkflowTriggers(text string) workflowTriggerShape {
	var shape workflowTriggerShape
	inOn, inPush, inBranches := false, false, false
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimRight(raw, " \t\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inOn = line == "on:"
			inPush, inBranches = false, false
			continue
		}
		if !inOn {
			continue
		}
		if indent == 2 {
			inPush, inBranches = trimmed == "push:", false
			if trimmed == "pull_request:" {
				shape.pullRequest = true
			}
			continue
		}
		if inPush && indent == 4 && trimmed == "branches:" {
			shape.pushBranches, inBranches = true, true
			continue
		}
		if inBranches && indent == 6 && trimmed == "- main" {
			shape.mainBranch = true
		}
	}
	return shape
}
