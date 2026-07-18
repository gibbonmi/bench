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
		if !bytes.HasSuffix(instructions, []byte("\n")) || !strings.Contains(string(instructions), "--offline") || !strings.Contains(string(instructions), "sha256sum -c SHA256SUMS") || !strings.Contains(string(instructions), "npm publish ./packages/") || !strings.Contains(string(instructions), "## Removal") {
			t.Fatalf("offline archive %s instructions are incomplete or not LF-terminated", name)
		}
	}
}
