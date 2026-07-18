package conformance

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

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
		if job := workflowJob(native, "evidence"); !strings.Contains(job, "needs: [preflight, native-proof]") || !strings.Contains(job, "scripts/release-preflight.sh --mode verify") || !strings.Contains(job, "verify-preflight-evidence") {
			diags = append(diags, "native verification does not finalize full release evidence after native proofs")
		}
		if job := workflowJob(native, "smoke"); !strings.Contains(job, "needs: [preflight, evidence]") || !strings.Contains(job, "verify-preflight-evidence") || !strings.Contains(job, "scripts/smoke-artifacts.sh") {
			diags = append(diags, "native runner matrix does not consume finalized release evidence")
		}
	}
	if release != "" {
		for message, anchor := range map[string]string{"tag publication bypasses full release preflight": "scripts/release-preflight.sh --mode publish --profile public\n", "tag publication does not upload preflight evidence": "publish-preflight-evidence", "tag evidence does not request repository maximum retention": "retention-days: ${{ github.retention_days }}", "publication does not wait for publish authorization": "needs: [authorize]"} {
			if !strings.Contains(release, anchor) {
				diags = append(diags, message)
			}
		}
		if job := workflowJob(release, "evidence"); !strings.Contains(job, "needs: [preflight, native-proof]") || !strings.Contains(job, "scripts/release-preflight.sh --mode verify") || !strings.Contains(job, "release-native-proof-evidence") {
			diags = append(diags, "tag evidence does not finalize complete native proofs")
		}
		if job := workflowJob(release, "smoke"); !strings.Contains(job, "needs: [preflight, evidence]") || !strings.Contains(job, "verify-preflight-evidence") {
			diags = append(diags, "tag smoke does not consume finalized release evidence")
		}
		if job := workflowJob(release, "authorize"); !strings.Contains(job, "needs: [preflight, evidence, smoke]") || !strings.Contains(job, "release-native-proof-evidence") {
			diags = append(diags, "publication does not wait for finalized evidence and every native proof row")
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
	if comparison != "" && (!strings.Contains(comparison, `cmp -s "$left/$name" "$right/$name"`) || !strings.Contains(comparison, "reproducibility mismatch:")) {
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
	var plan struct {
		Targets []struct {
			OS     string `json:"os"`
			Arch   string `json:"arch"`
			Runner string `json:"runner"`
		} `json:"targets"`
	}
	if err := readJSONAt(root, filepath.Join("scripts", "release-plan.json"), &plan); err != nil {
		return []string{"release package evidence matrix is unreadable"}
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := readJSONAt(root, filepath.Join("package.json"), &pkg); err != nil {
		return []string{"release package evidence package identity is unreadable"}
	}
	artifacts := append([]string{"redbench-" + pkg.Version + ".tgz"}, func() []string {
		out := make([]string, 0, len(plan.Targets))
		for _, target := range plan.Targets {
			out = append(out, fmt.Sprintf("redbench-%s-%s-%s.tgz", target.OS, target.Arch, pkg.Version))
		}
		return out
	}()...)
	for _, name := range artifacts {
		files, err := archiveFiles(filepath.Join(root, "dist", "artifacts", name))
		if err != nil {
			return []string{"release package evidence artifact is unreadable: " + name}
		}
		for _, evidence := range packageEvidence {
			if evidence.PackageMode == "" {
				continue
			}
			if _, ok := files["package/"+evidence.Path]; !ok {
				return []string{"release package evidence allowlist omits " + evidence.Path}
			}
		}
	}
	proofDir := filepath.Join(root, "dist", "native-proofs")
	if err := os.MkdirAll(proofDir, 0o755); err != nil {
		return []string{"release evidence probe could not prepare native proof directory"}
	}
	digest := func(data []byte) string { sum := sha256.Sum256(data); return fmt.Sprintf("%x", sum) }
	for _, target := range plan.Targets {
		platform := fmt.Sprintf("redbench-%s-%s-%s.tgz", target.OS, target.Arch, pkg.Version)
		archive := fmt.Sprintf("redbench-%s-%s-%s.tar.gz", pkg.Version, target.OS, target.Arch)
		files, err := archiveFiles(filepath.Join(root, "dist", "artifacts", platform))
		binary, ok := files["package/bin/bench"]
		if err != nil || !ok {
			return []string{"release evidence probe could not inspect native proof binary"}
		}
		packageBytes, packageErr := os.ReadFile(filepath.Join(root, "dist", "artifacts", platform))
		archiveBytes, archiveErr := os.ReadFile(filepath.Join(root, "dist", "artifacts", archive))
		if packageErr != nil || archiveErr != nil {
			return []string{"release evidence probe could not read native proof artifacts"}
		}
		proof, marshalErr := json.Marshal(map[string]any{"schema_version": 1, "target": target.OS + "-" + target.Arch, "runner": target.Runner, "status": "green", "rebuilt_sha256": digest(binary), "binary_sha256": digest(binary), "package_sha256": digest(packageBytes), "archive_sha256": digest(archiveBytes), "musl_status": map[bool]string{true: "green", false: "not_applicable"}[target.OS == "linux"], "operations_status": "green", "strip_status": "green", "tools_status": "green"})
		if marshalErr != nil || os.WriteFile(filepath.Join(proofDir, target.OS+"-"+target.Arch+".json"), append(proof, '\n'), 0o644) != nil {
			return []string{"release evidence probe could not write native proof"}
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

func archiveFiles(path string) (map[string][]byte, error) {
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
	files := map[string][]byte{}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return files, nil
		}
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		files[header.Name] = body
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
