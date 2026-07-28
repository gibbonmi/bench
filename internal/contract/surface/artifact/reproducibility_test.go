package artifact

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/testrepo"
)

func TestGoBuildIgnoresCheckoutTopology(t *testing.T) {
	root := contract.SubjectRoot(t)
	clone := filepath.Join(t.TempDir(), "isolated-source")
	command := exec.Command("git", "clone", "--quiet", "--no-hardlinks", root, clone)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("clone isolated source: %v\n%s", err, output)
	}
	// Mirror the working tree onto the clone so uncommitted DELETIONS (a file removed in
	// the source but still tracked at HEAD) are reflected too — a tar overlay can only add
	// or modify, leaving a deleted source file resurrected in the clone and drifting the
	// binary. --delete makes the clone tree a true mirror; .git and dist stay untouched.
	if _, err := exec.LookPath("rsync"); err != nil {
		capability.Capability(t, capability.Tool, fmt.Sprintf("reproducibility probe needs rsync on PATH: %v", err))
	}
	overlay := exec.Command("rsync", "-a", "--delete", "--exclude=/.git", "--exclude=/dist", root+"/", clone)
	if output, err := overlay.CombinedOutput(); err != nil {
		t.Fatalf("overlay source snapshot: %v\n%s", err, output)
	}
	outputs := []string{filepath.Join(t.TempDir(), "worktree-binary"), filepath.Join(t.TempDir(), "clone-binary")}
	for index, source := range []string{root, clone} {
		build := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), source, outputs[index])
		build.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=darwin", "GOARCH=arm64")
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build checkout %d: %v\n%s", index, err, output)
		}
	}
	first, err := os.ReadFile(outputs[0])
	if err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(outputs[1])
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(first) != sha256.Sum256(second) {
		t.Fatal("checkout topology changed release binary bytes")
	}
}

func TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/compare-artifacts.sh")
	var plan struct {
		Targets []contract.ReleaseTarget `json:"targets"`
	}
	var wrapper struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "release-plan.json"), &plan)
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &wrapper)
	names := []string{"redbench-" + wrapper.Version + ".tgz"}
	for _, target := range plan.Targets {
		names = append(names, fmt.Sprintf("redbench-%s-%s-%s.tgz", target.OS, target.Arch, wrapper.Version), fmt.Sprintf("redbench-%s-%s-%s.tar.gz", wrapper.Version, target.OS, target.Arch))
	}
	left, right := filepath.Join(t.TempDir(), "first"), filepath.Join(t.TempDir(), "second")
	for _, directory := range []string{left, right} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, name := range names {
			if err := os.WriteFile(filepath.Join(directory, name), []byte("identical "+name+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	rightRoot := copyReleaseBoundEvidenceRoot(t, root)
	finalLeft, finalRight := filepath.Join(t.TempDir(), "first-final"), filepath.Join(t.TempDir(), "second-final")
	for _, directory := range []string{finalLeft, finalRight} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"release-index.json", "SHA256SUMS", "run-manifest.json", "phase-artifacts.json"} {
			if err := os.WriteFile(filepath.Join(directory, name), []byte("independently generated "+name+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	record := filepath.Join(t.TempDir(), "reproducibility.json")
	comparator := filepath.Join(root, "scripts", "compare-artifacts.sh")
	externalComparator, cleanupComparator, err := testrepo.TwoHopRelativeSymlink(comparator)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupComparator()
	runWith := func(script string) contract.Probe {
		return contract.NewExecFixtureAt(t, root).Run("bash", script, left, right, record, root, rightRoot, finalLeft, finalRight)
	}
	run := func() contract.Probe {
		return runWith(comparator)
	}
	sameRoot := contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(root, "scripts", "compare-artifacts.sh"), left, right, record, root, root)
	assertComparatorRed(t, sameRoot, "reproducibility comparison requires isolated source roots")
	run().RequireExit(0)
	runWith(externalComparator).RequireExit(0)
	var green struct {
		Artifacts []struct {
			Match bool `json:"match"`
		} `json:"artifacts"`
		Evidence []struct {
			Match bool `json:"match"`
		} `json:"evidence"`
	}
	data, err := os.ReadFile(record)
	if err != nil || json.Unmarshal(data, &green) != nil || len(green.Artifacts) != 9 || len(green.Evidence) == 0 {
		t.Fatalf("valid reproducibility record is incomplete: %s, %v", data, err)
	}
	for _, item := range append(green.Artifacts, green.Evidence...) {
		if !item.Match {
			t.Fatal("valid reproducibility record fabricated a non-match")
		}
	}
	first := names[0]
	if err := os.WriteFile(filepath.Join(right, first), []byte("different\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	assertComparatorRed(t, run(), "reproducibility mismatch: "+first)
	if err := os.WriteFile(filepath.Join(right, first), []byte("identical "+first+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(right, first)); err != nil {
		t.Fatal(err)
	}
	assertComparatorRed(t, run(), "reproducibility comparison missing second-build artifact: "+first)
	if err := os.WriteFile(filepath.Join(right, first), []byte("identical "+first+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(right, "unexpected-artifact"), []byte("extra\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	assertComparatorRed(t, run(), "reproducibility comparison found unexpected artifact: unexpected-artifact")
	if err := os.Remove(filepath.Join(right, "unexpected-artifact")); err != nil {
		t.Fatal(err)
	}
	evidencePath := filepath.Join(rightRoot, "LICENSE")
	if err := os.WriteFile(evidencePath, []byte("different evidence\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	assertComparatorRed(t, run(), "reproducibility release-evidence mismatch: LICENSE")
	if err := os.Remove(evidencePath); err != nil {
		t.Fatal(err)
	}
	assertComparatorRed(t, run(), "reproducibility comparison missing second-build release evidence: LICENSE")
	license, err := os.ReadFile(filepath.Join(root, "LICENSE"))
	if err != nil || os.WriteFile(evidencePath, license, 0o644) != nil {
		t.Fatalf("restore isolated evidence root: %v", err)
	}
	if err := os.WriteFile(filepath.Join(finalRight, "release-index.json"), []byte("drifted final index\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	assertComparatorRed(t, run(), "reproducibility final-evidence mismatch: release-index.json")
}

func TestOfflineSmokeRequiresApprovedReleaseEvidence(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/smoke-offline.sh")
	empty := t.TempDir()
	probe := contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(root, "scripts", "smoke-offline.sh"), empty, filepath.Join(empty, "missing-evidence"))
	if probe.ExitCode == 0 {
		t.Fatal("offline smoke accepted missing approved release evidence")
	}
	probe.RequireContains(probe.Stderr, "offline smoke: approved release evidence is missing or unsafe")
}

func TestReleaseArtifactVerifierRequiresFullyApprovedEvidence(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/verify-release-artifact.mjs")
	dir := t.TempDir()
	artifact := filepath.Join(dir, "redbench-0.1.0-linux-x64.tar.gz")
	data := []byte("exact artifact bytes\n")
	if err := os.WriteFile(artifact, data, 0o644); err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(data))
	indexPath, sumsPath := filepath.Join(dir, "release-index.json"), filepath.Join(dir, "SHA256SUMS")
	index, sums := approvedReleaseEvidenceFixture(t, root, "0.1.0", filepath.Base(artifact), digest, int64(len(data)))
	write := func() {
		encoded, err := json.Marshal(index)
		if err != nil || os.WriteFile(indexPath, append(encoded, '\n'), 0o644) != nil || os.WriteFile(sumsPath, []byte(sums), 0o644) != nil {
			t.Fatalf("write evidence fixture: %v", err)
		}
	}
	write()
	run := func() contract.Probe {
		return contract.NewExecFixtureAt(t, root).Run("node", filepath.Join(root, "scripts", "verify-release-artifact.mjs"), indexPath, sumsPath, artifact)
	}
	run().RequireExit(0)
	for _, mutation := range []struct {
		name, diagnostic string
		apply            func()
	}{
		{"approval status", "release evidence is not a fully approved preflight", func() { index["status"] = "red" }},
		{"artifact completeness", "release phase or artifact inventory is incomplete", func() {
			artifacts := index["artifacts"].([]map[string]any)
			index["artifacts"] = artifacts[:len(artifacts)-1]
		}},
		{"phase completeness", "release phase or artifact inventory is incomplete", func() {
			phases := index["phases"].([]map[string]any)
			index["phases"] = phases[:len(phases)-1]
		}},
		{"target matrix os", "release target evidence is malformed", func() { index["targets"].([]map[string]any)[0]["os"] = "index-owned-os" }},
		{"target matrix arch", "release target evidence is malformed", func() { index["targets"].([]map[string]any)[0]["arch"] = "index-owned-arch" }},
		{"target matrix goos", "release target evidence is malformed", func() { index["targets"].([]map[string]any)[0]["goos"] = "index-owned-goos" }},
		{"target matrix goarch", "release target evidence is malformed", func() { index["targets"].([]map[string]any)[0]["goarch"] = "index-owned-goarch" }},
		{"target matrix runner", "release target evidence is malformed", func() { index["targets"].([]map[string]any)[0]["runner"] = "index-owned-runner" }},
		{"target matrix order", "release target evidence is malformed", func() {
			targets := index["targets"].([]map[string]any)
			targets[0], targets[1] = targets[1], targets[0]
		}},
		{"target matrix cardinality", "release target evidence is malformed", func() {
			targets := index["targets"].([]map[string]any)
			index["targets"] = targets[:len(targets)-1]
		}},
		{"requirement completeness", "release requirements are incomplete or unsatisfied", func() {
			requirements := index["requirements"].([]map[string]any)
			index["requirements"] = requirements[:len(requirements)-1]
		}},
		{"reproducibility", "reproducibility evidence is incomplete or mismatched", func() { index["reproducibility"].(map[string]any)["status"] = "red" }},
		{"native proof", "native proof evidence is incomplete or red", func() { index["native_proofs"].([]map[string]any)[0]["strip_status"] = "red" }},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			index, sums = approvedReleaseEvidenceFixture(t, root, "0.1.0", filepath.Base(artifact), digest, int64(len(data)))
			mutation.apply()
			write()
			probe := run()
			if probe.ExitCode == 0 {
				t.Fatalf("offline artifact verifier accepted %s mutation", mutation.name)
			}
			probe.RequireContains(probe.Stderr, "offline smoke: "+mutation.diagnostic)
		})
	}
	index, sums = approvedReleaseEvidenceFixture(t, root, "0.1.0", filepath.Base(artifact), digest, int64(len(data)))
	write()
	mutatedSums := strings.Replace(sums, digest+"  "+filepath.Base(artifact), strings.Repeat("0", 64)+"  "+filepath.Base(artifact), 1)
	if err := os.WriteFile(sumsPath, []byte(mutatedSums), 0o644); err != nil {
		t.Fatal(err)
	}
	probe := run()
	if probe.ExitCode == 0 {
		t.Fatal("offline artifact verifier accepted a checksum mutation")
	}
	probe.RequireContains(probe.Stderr, "offline smoke: SHA256SUMS does not derive exactly from the release index")
}

func approvedReleaseEvidenceFixture(t *testing.T, root, version, selected, selectedDigest string, selectedSize int64) (map[string]any, string) {
	t.Helper()
	hex := fmt.Sprintf("%064x", 1)
	targets := []map[string]any{}
	proofs := []map[string]any{}
	artifacts := []map[string]any{}
	var plan struct {
		Targets []struct {
			OS, Arch, GOOS, GOARCH, Runner string
		} `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "release-plan.json"), &plan)
	for _, target := range plan.Targets {
		name := target.OS + "-" + target.Arch
		targets = append(targets, map[string]any{"os": target.OS, "arch": target.Arch, "goos": target.GOOS, "goarch": target.GOARCH, "runner": target.Runner})
		proofs = append(proofs, map[string]any{"schema_version": 1, "target": name, "runner": target.Runner, "status": "green", "rebuilt_sha256": hex, "binary_sha256": hex, "package_sha256": hex, "archive_sha256": hex, "musl_status": map[bool]string{true: "green", false: "not_applicable"}[target.OS == "linux"], "operations_status": "green", "strip_status": "green", "tools_status": "green"})
	}
	command := exec.Command("node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "artifact-records", version)
	data, err := command.Output()
	if err != nil {
		t.Fatalf("read canonical artifact records: %v\n%s", err, data)
	}
	var records []struct{ Name, Target, Kind string }
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("decode canonical artifact records: %v\n%s", err, data)
	}
	artifactsByTargetKind := make(map[string]map[string]any, len(records))
	for _, record := range records {
		artifact := map[string]any{"name": record.Name, "target": record.Target, "size": int64(1), "sha256": hex, "component_manifest_sha256": hex, "sbom_sha256": hex, "inventory_sha256": hex}
		artifacts = append(artifacts, artifact)
		artifactsByTargetKind[record.Target+"/"+record.Kind] = artifact
	}
	for _, artifact := range artifacts {
		if artifact["name"] == selected {
			artifact["sha256"], artifact["size"] = selectedDigest, selectedSize
		}
	}
	for _, proof := range proofs {
		target := proof["target"].(string)
		proof["package_sha256"] = artifactsByTargetKind[target+"/platform"]["sha256"]
		proof["archive_sha256"] = artifactsByTargetKind[target+"/archive"]["sha256"]
	}
	repro := []map[string]any{}
	var sums strings.Builder
	for _, artifact := range artifacts {
		repro = append(repro, map[string]any{"name": artifact["name"], "size": artifact["size"], "sha256": artifact["sha256"], "match": true})
		fmt.Fprintf(&sums, "%s  %s\n", artifact["sha256"], artifact["name"])
	}
	var registry struct {
		Verify        []string `json:"verify"`
		PublishBefore []string `json:"publish_before"`
		PublishOnly   []string `json:"publish_only"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "internal", "releaseevidence", "registry.json"), &registry)
	phases := []map[string]any{}
	for _, name := range append(append(append([]string{}, registry.PublishBefore...), registry.Verify...), registry.PublishOnly...) {
		phases = append(phases, map[string]any{"name": name, "status": "green", "record_sha256": hex})
	}
	var requirementRegistry struct {
		Records []struct {
			Key          string `json:"key"`
			Owner        string `json:"owner"`
			Schema       string `json:"schema"`
			Requiredness string `json:"requiredness"`
		} `json:"records"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "internal", "releaseevidence", "requirements.json"), &requirementRegistry)
	requirements := []map[string]any{}
	for _, item := range requirementRegistry.Records {
		requirements = append(requirements, map[string]any{"key": item.Key, "owner": item.Owner, "schema": item.Schema, "requiredness": item.Requiredness, "applicable": true, "status": "satisfied"})
	}
	index := map[string]any{"schema_version": 1, "mode": "publish", "scope": "preflight", "profile": "public", "status": "green", "identity": map[string]any{"package_version": version, "source_commit": hex, "binary_version": version, "toolchain": "go"}, "targets": targets, "artifacts": artifacts, "phases": phases, "requirements": requirements, "native_proofs": proofs, "reproducibility": map[string]any{"schema_version": 1, "status": "green", "builds": 2, "artifacts": repro, "evidence": []map[string]any{{"name": "LICENSE", "size": int64(1), "sha256": hex, "match": true}}}}
	return index, sums.String()
}

func assertComparatorRed(t *testing.T, probe contract.Probe, message string) {
	t.Helper()
	if probe.ExitCode == 0 {
		t.Fatalf("mutation passed reproducibility comparison: %s", message)
	}
	probe.RequireContains(probe.Stderr, message)
}

func copyReleaseBoundEvidenceRoot(t *testing.T, root string) string {
	t.Helper()
	var requirements struct {
		Records []struct {
			Path        string `json:"path"`
			PackageMode string `json:"package_mode"`
		} `json:"records"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "internal", "releaseevidence", "requirements.json"), &requirements)
	destination := t.TempDir()
	paths := []string{"internal/releaseevidence/requirements.json", "scripts/release-plan.json"}
	for _, requirement := range requirements.Records {
		if requirement.PackageMode != "" {
			paths = append(paths, requirement.Path)
		}
	}
	for _, relative := range paths {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(destination, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return destination
}
