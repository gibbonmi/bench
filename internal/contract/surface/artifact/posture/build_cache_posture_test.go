package posture

import fixture "github.com/gibbonmi/bench/internal/contract/surface/artifact/internal/fixture"

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestMain puts the whole package on the dev-tier build posture; contract.ProcessEnv and
// the fixture-driven mergeEnv both start from os.Environ(), so it reaches every subprocess
// the package spawns. The rows asserting the hermetic default strip the token back out
// through ambientBuildEnv. It also removes the package-owned shared set after every run.

func ambientBuildEnv(extra []string, remove ...string) []string {
	dropped := make(map[string]bool, len(remove))
	for _, key := range remove {
		dropped[key] = true
	}
	ambient := os.Environ()
	env := make([]string, 0, len(ambient)+len(extra))
	for _, entry := range ambient {
		if key, _, ok := strings.Cut(entry, "="); ok && !dropped[key] {
			env = append(env, entry)
		}
	}
	return append(env, extra...)
}

func TestArtifactBuilderHonorsHermeticDefault(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	for _, test := range []struct {
		name  string
		token []string
	}{
		{name: "absent"},
		{name: "unrecognized value", token: []string{contract.SharedBuildCacheEnv + "=yes"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := fixture.CommittedHostileArtifactSource(t, root)
			output := filepath.Join(t.TempDir(), "hermetic artifact output [*]")
			build := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
			build.Env = ambientBuildEnv(test.token, contract.SharedBuildCacheEnv)
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
		{name: "hermetic default refuses offline", env: ambientBuildEnv([]string{"GOPROXY=off"}, contract.SharedBuildCacheEnv)},
		{name: "opt-in reuses passed caches", env: ambientBuildEnv([]string{"GOPROXY=off", "GOCACHE=" + goCache, "GOMODCACHE=" + goModCache}, "GOCACHE", "GOMODCACHE"), wantBuild: true},
		{name: "opt-in resolves caches before the HOME override", env: ambientBuildEnv([]string{"GOPROXY=off"}, "GOCACHE", "GOMODCACHE"), wantBuild: true},
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
		{name: "populated", body: staleReproducibilityRecord},
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

const staleReproducibilityRecord = "{\"schema_version\":1,\"status\":\"green\",\"builds\":2}\n"
