package posture

import fixture "github.com/gibbonmi/bench/internal/contract/surface/artifact/internal/fixture"

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactBuilderHonorsHermeticDefault(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	for _, test := range []struct {
		name  string
		token []string
	}{
		{name: "absent"},
		{name: "unrecognized value", token: []string{fixture.InvalidSharedCacheToken()}},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := fixture.CommittedHostileArtifactSource(t, root)
			output := filepath.Join(t.TempDir(), "hermetic artifact output [*]")
			build := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
			build.Env = fixture.AmbientBuildEnvWithoutSharedCache(test.token)
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("hermetic artifact build failed: %v\n%s", err, out)
			}
			fixture.AssertPromotedReproducibility(t, output)
		})
	}
}

// TestBuildCachePostureUnderGoproxyOff separates the two postures by whether the build can
// run offline: only a build that reused the ambient module cache resolves
// github.com/toon-format/toon-go without a proxy.
func TestBuildCachePostureUnderGoproxyOff(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	goCache, goModCache := contract.GoEnvPair(t, "GOCACHE", "GOMODCACHE")
	for _, test := range []struct {
		name      string
		env       []string
		wantBuild bool
	}{
		{name: "hermetic default refuses offline", env: fixture.AmbientBuildEnvWithoutSharedCache([]string{"GOPROXY=off"})},
		{name: "opt-in reuses passed caches", env: fixture.AmbientBuildEnv([]string{"GOPROXY=off", "GOCACHE=" + goCache, "GOMODCACHE=" + goModCache}, "GOCACHE", "GOMODCACHE"), wantBuild: true},
		{name: "opt-in resolves caches before the HOME override", env: fixture.AmbientBuildEnv([]string{"GOPROXY=off"}, "GOCACHE", "GOMODCACHE"), wantBuild: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := fixture.CommittedHostileArtifactSource(t, root)
			output := filepath.Join(t.TempDir(), "offline artifact output [*]")
			build := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
			build.Env = test.env
			out, err := build.CombinedOutput()
			if test.wantBuild && err != nil {
				t.Fatalf("shared-cache build did not reuse the ambient module cache offline: %v\n%s", err, out)
			}
			if !test.wantBuild && err == nil {
				t.Fatalf("hermetic build succeeded under GOPROXY=off:\n%s", out)
			}
		})
	}
}

func TestSharedCacheBuildRemovesStaleReproducibilityRecord(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	for _, test := range []struct {
		name string
		body string
	}{
		{name: "populated", body: fixture.StaleReproducibilityRecord},
		{name: "zero byte", body: ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := fixture.CommittedHostileArtifactSource(t, root)
			output := filepath.Join(t.TempDir(), "stale record artifact output [*]")
			record := filepath.Join(filepath.Dir(output), "reproducibility.json")
			if err := os.WriteFile(record, []byte(test.body), 0o644); err != nil {
				t.Fatal(err)
			}
			contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output).RequireExit(0)
			if _, err := os.Stat(record); !os.IsNotExist(err) {
				t.Fatalf("shared-cache build left a stale reproducibility record: %v", err)
			}
		})
	}
}
