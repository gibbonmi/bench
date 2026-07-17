package artifact

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestOfflineArchiveProjection(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	buildRoot := committedHostileArtifactSource(t, root)
	npmOut := filepath.Join(t.TempDir(), "npm artifacts [hostile]")
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-artifacts.sh"), buildRoot, npmOut).RequireExit(0)

	var matrix []artifactPlatform
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "platforms.json"), &matrix)
	var wrapper struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &wrapper)
	offlineOut := filepath.Join(t.TempDir(), "offline artifacts [hostile]")
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(buildRoot, "scripts", "build-offline-archives.sh"), npmOut, offlineOut).RequireExit(0)
	assertOfflineArchiveSet(t, offlineOut, wrapper.Version, matrix)
}

func assertOfflineArchiveSet(t *testing.T, output, version string, matrix []artifactPlatform) {
	t.Helper()
	files, err := os.ReadDir(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 4 {
		t.Fatalf("offline archive count = %d, want exactly four", len(files))
	}
	for _, platform := range matrix {
		root := fmt.Sprintf("redbench-%s-%s-%s", version, platform.OS, platform.Arch)
		name := root + ".tar.gz"
		archive := filepath.Join(output, name)
		info, err := os.Stat(archive)
		if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
			t.Fatalf("offline archive %s is missing or invalid: %v", name, err)
		}
		entries := contract.ReadTarball(t, archive)
		for _, packageName := range []string{
			root + "/packages/redbench-" + version + ".tgz",
			root + "/packages/redbench-" + platform.OS + "-" + platform.Arch + "-" + version + ".tgz",
		} {
			if entry, ok := entries[packageName]; !ok || len(entry.Data) == 0 {
				t.Fatalf("offline archive %s does not carry non-empty %s", name, packageName)
			}
		}
	}
}
