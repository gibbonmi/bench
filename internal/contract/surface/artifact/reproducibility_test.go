package artifact

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/compare-artifacts.sh")
	var plan struct {
		Targets []artifactPlatform `json:"targets"`
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
	record := filepath.Join(t.TempDir(), "reproducibility.json")
	run := func() contract.Probe {
		return contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(root, "scripts", "compare-artifacts.sh"), left, right, record, root, rightRoot)
	}
	run().RequireExit(0)
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

func TestReleaseArtifactVerifierAcceptsValidEvidenceAndRejectsChecksumMutation(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/verify-release-artifact.mjs")
	dir := t.TempDir()
	artifact := filepath.Join(dir, "redbench-0.1.0-linux-x64.tar.gz")
	data := []byte("exact artifact bytes\n")
	if err := os.WriteFile(artifact, data, 0o644); err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(data))
	index := []byte(fmt.Sprintf("{\"artifacts\":[{\"name\":\"%s\",\"sha256\":\"%s\"}]}\n", filepath.Base(artifact), digest))
	indexPath, sumsPath := filepath.Join(dir, "release-index.json"), filepath.Join(dir, "SHA256SUMS")
	if err := os.WriteFile(indexPath, index, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sumsPath, []byte(digest+"  "+filepath.Base(artifact)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run := func() contract.Probe {
		return contract.NewExecFixtureAt(t, root).Run("node", filepath.Join(root, "scripts", "verify-release-artifact.mjs"), indexPath, sumsPath, artifact)
	}
	run().RequireExit(0)
	if err := os.WriteFile(sumsPath, []byte(strings.Repeat("0", 64)+"  "+filepath.Base(artifact)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	probe := run()
	if probe.ExitCode == 0 {
		t.Fatal("offline artifact verifier accepted a checksum mutation")
	}
	probe.RequireContains(probe.Stderr, "offline smoke: supplied release evidence does not bind artifact bytes")
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
