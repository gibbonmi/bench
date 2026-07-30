package prepared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func writeAggregateProof(t *testing.T, directory, osName, arch, runner string, mismatch bool) {
	t.Helper()
	digest, binary := strings.Repeat("1", 64), strings.Repeat("1", 64)
	if mismatch {
		binary = strings.Repeat("2", 64)
	}
	musl := "not_applicable"
	if osName == "linux" {
		musl = "green"
	}
	name := osName + "-" + arch
	proof := map[string]any{"schema_version": 1, "target": name, "runner": runner, "status": "green", "rebuilt_sha256": digest, "binary_sha256": binary, "package_sha256": digest, "archive_sha256": digest, "musl_status": musl, "operations_status": "green", "strip_status": "green", "tools_status": "green"}
	data, err := json.Marshal(proof)
	if err != nil || os.WriteFile(filepath.Join(directory, name+".json"), append(data, '\n'), 0o644) != nil {
		t.Fatalf("write native proof: %v", err)
	}
}

func TestOfflineArchiveProjection(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	buildRoot, npmOut := shared.sourceRoot, shared.outputDir

	var plan struct {
		Targets []contract.ReleaseTarget `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(buildRoot, "scripts", "release-plan.json"), &plan)
	matrix := plan.Targets
	var wrapper struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &wrapper)
	offlineOut := filepath.Join(t.TempDir(), "offline artifacts [hostile]")
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-offline-archives.sh"), npmOut, offlineOut).RequireExit(0)
	assertOfflineArchiveSet(t, npmOut, offlineOut, wrapper.Version, matrix)
}

// The offline builder replaces its output directory wholesale — it moves the old
// one aside and deletes it. Anything in there that the build cannot account for
// is somebody's work, so the refusal has to come before the archives are built
// rather than beside the swap that would destroy it.

func assembleOfflineArchiveStage(t *testing.T, scriptRoot, planRoot, evidenceRoot string) string {
	t.Helper()
	stage := t.TempDir()
	npmDir, archiveDir := filepath.Join(stage, "packages"), filepath.Join(stage, "archive")
	wrapperExtract, platformExtract := filepath.Join(stage, "wrapper"), filepath.Join(stage, "platform")
	for _, directory := range []string{npmDir, archiveDir, filepath.Join(wrapperExtract, "package"), filepath.Join(platformExtract, "package")} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(path string, data []byte, mode os.FileMode) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, mode); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(stage, "bench"), []byte("fixture binary\n"), 0o755)
	write(filepath.Join(npmDir, "redbench-0.1.0.tgz"), []byte("wrapper\n"), 0o644)
	write(filepath.Join(npmDir, "redbench-linux-x64-0.1.0.tgz"), []byte("platform\n"), 0o644)
	write(filepath.Join(wrapperExtract, "package", "component-manifest.json"), []byte("{}\n"), 0o644)
	write(filepath.Join(platformExtract, "package", "component-manifest.json"), []byte("{}\n"), 0o644)
	var requirements struct {
		Records []struct {
			Path        string `json:"path"`
			PackageMode string `json:"package_mode"`
		} `json:"records"`
	}
	contract.ReadJSONFile(t, filepath.Join(planRoot, "internal", "releaseevidence", "requirements.json"), &requirements)
	for _, record := range requirements.Records {
		if record.PackageMode != "" {
			data, err := os.ReadFile(filepath.Join(evidenceRoot, filepath.FromSlash(record.Path)))
			if err != nil {
				t.Fatal(err)
			}
			write(filepath.Join(wrapperExtract, "package", filepath.FromSlash(record.Path)), data, 0o644)
		}
	}
	contract.NewExecFixtureAt(t, scriptRoot).Run("node", filepath.Join(scriptRoot, "scripts", "assemble-offline-archive.mjs"), planRoot, npmDir, archiveDir, "linux-x64", "0.1.0", filepath.Join(stage, "bench"), wrapperExtract, platformExtract).RequireExit(0)
	return archiveDir
}

func assertOfflineArchiveSet(t *testing.T, npmArtifacts, output, version string, matrix []contract.ReleaseTarget) {
	t.Helper()
	files, err := os.ReadDir(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(matrix) {
		t.Fatalf("offline archive count = %d, want release-plan count %d", len(files), len(matrix))
	}
	var requirements struct {
		Records []struct {
			Path        string `json:"path"`
			PackageMode string `json:"package_mode"`
		} `json:"records"`
	}
	contract.ReadJSONFile(t, filepath.Join(contract.SubjectRoot(t), "internal", "releaseevidence", "requirements.json"), &requirements)
	for _, platform := range matrix {
		root := fmt.Sprintf("redbench-%s-%s-%s", version, platform.OS, platform.Arch)
		name := root + ".tar.gz"
		archive := filepath.Join(output, name)
		info, err := os.Stat(archive)
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			t.Fatalf("offline archive %s is missing or invalid: %v", name, err)
		}
		entries := contract.ReadTarball(t, archive)
		wrapperName := "redbench-" + version + ".tgz"
		platformName := "redbench-" + platform.OS + "-" + platform.Arch + "-" + version + ".tgz"
		want := map[string]int64{
			root + "/bin/bench":                                            0o755,
			root + "/packages/" + wrapperName:                              0o644,
			root + "/packages/" + platformName:                             0o644,
			root + "/OFFLINE.md":                                           0o644,
			root + "/evidence/component-manifest.json":                     0o644,
			root + "/evidence/components/wrapper-component-manifest.json":  0o644,
			root + "/evidence/components/platform-component-manifest.json": 0o644,
		}
		for _, record := range requirements.Records {
			if record.PackageMode != "" {
				want[root+"/evidence/"+record.Path] = 0o644
			}
		}
		if len(entries) != len(want) {
			t.Fatalf("offline archive %s entry count = %d, want %d", name, len(entries), len(want))
		}
		for path, mode := range want {
			entry, ok := entries[path]
			if !ok || entry.Mode&0o777 != mode || len(entry.Data) == 0 {
				t.Fatalf("offline archive %s entry %s is absent, unsafe, or empty", name, path)
			}
		}
		wrapperBytes, err := os.ReadFile(filepath.Join(npmArtifacts, wrapperName))
		if err != nil || !bytes.Equal(entries[root+"/packages/"+wrapperName].Data, wrapperBytes) {
			t.Fatalf("offline archive %s wrapper package differs from the approved tarball", name)
		}
		platformBytes, err := os.ReadFile(filepath.Join(npmArtifacts, platformName))
		if err != nil || !bytes.Equal(entries[root+"/packages/"+platformName].Data, platformBytes) {
			t.Fatalf("offline archive %s platform package differs from the approved tarball", name)
		}
		platformPackage := contract.ReadTarball(t, filepath.Join(npmArtifacts, platformName))
		if !bytes.Equal(entries[root+"/bin/bench"].Data, platformPackage["package/bin/bench"].Data) {
			t.Fatalf("offline archive %s binary differs from its platform package", name)
		}
		instructions := entries[root+"/OFFLINE.md"].Data
		checksumCommand := fmt.Sprintf(`awk '$2 == "%s" { print }' SHA256SUMS | sha256sum -c -`, name)
		if !bytes.HasSuffix(instructions, []byte("\n")) || !strings.Contains(string(instructions), "--offline") || !strings.Contains(string(instructions), checksumCommand) || strings.Contains(string(instructions), "sha256sum -c SHA256SUMS") || !strings.Contains(string(instructions), "npm publish ./packages/") || !strings.Contains(string(instructions), "## Removal") {
			t.Fatalf("offline archive %s instructions are incomplete or not LF-terminated", name)
		}
	}
}
