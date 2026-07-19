package artifact

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestOfflineArchiveProjection(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	buildRoot := committedHostileArtifactSource(t, root)
	npmOut := filepath.Join(t.TempDir(), "npm artifacts [hostile]")
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-artifacts.sh"), buildRoot, npmOut).RequireExit(0)

	var plan struct {
		Targets []artifactPlatform `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "release-plan.json"), &plan)
	matrix := plan.Targets
	var wrapper struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &wrapper)
	offlineOut := filepath.Join(t.TempDir(), "offline artifacts [hostile]")
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-offline-archives.sh"), npmOut, offlineOut).RequireExit(0)
	assertOfflineArchiveSet(t, npmOut, offlineOut, wrapper.Version, matrix)
}

func TestOfflineInstructionsVerifyOnlyTargetArchive(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/assemble-offline-archive.mjs")
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
	contract.ReadJSONFile(t, filepath.Join(root, "internal", "releaseevidence", "requirements.json"), &requirements)
	for _, record := range requirements.Records {
		if record.PackageMode != "" {
			data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(record.Path)))
			if err != nil {
				t.Fatal(err)
			}
			write(filepath.Join(wrapperExtract, "package", filepath.FromSlash(record.Path)), data, 0o644)
		}
	}
	contract.NewExecFixtureAt(t, root).Run("node", filepath.Join(root, "scripts", "assemble-offline-archive.mjs"), root, npmDir, archiveDir, "linux-x64", "0.1.0", filepath.Join(stage, "bench"), wrapperExtract, platformExtract).RequireExit(0)
	instructions, err := os.ReadFile(filepath.Join(archiveDir, "OFFLINE.md"))
	if err != nil {
		t.Fatal(err)
	}
	const targeted = `awk '$2 == "redbench-0.1.0-linux-x64.tar.gz" { print }' SHA256SUMS | sha256sum -c -`
	if !bytes.Contains(instructions, []byte(targeted)) || bytes.Contains(instructions, []byte("sha256sum -c SHA256SUMS")) {
		t.Fatalf("offline instructions do not verify only the selected target archive:\n%s", instructions)
	}
}

func assertOfflineArchiveSet(t *testing.T, npmArtifacts, output, version string, matrix []artifactPlatform) {
	t.Helper()
	files, err := os.ReadDir(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("offline archive count = %d, want exactly four", len(files))
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
