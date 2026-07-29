package artifact

import (
	"bytes"
	"fmt"
	"maps"
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
func TestMain(m *testing.M) {
	if err := os.Setenv(contract.SharedBuildCacheEnv, "1"); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := packageSharedArtifactSet.cleanup(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "clean up shared artifact set: %v\n", err)
		code = 1
	}
	os.Exit(code)
}

// ambientBuildEnv is the real process environment plus extra, minus remove. These rows turn
// on what `go env` resolves against the ambient HOME, so they cannot use
// contract.NewExecFixtureAt, whose isolated HOME is a fresh temp directory.
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
			source := committedHostileArtifactSource(t, root)
			output := filepath.Join(t.TempDir(), "hermetic artifact output [*]")
			build := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
			build.Env = ambientBuildEnv(test.token, contract.SharedBuildCacheEnv)
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("hermetic artifact build failed: %v\n%s", err, out)
			}
			assertPromotedReproducibility(t, output)
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
			source := committedHostileArtifactSource(t, root)
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

func TestSharedCacheBuildPromotesNoRecord(t *testing.T) {
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	output := requireSharedArtifactSet(t).outputDir
	if _, err := os.Stat(filepath.Join(filepath.Dir(output), "reproducibility.json")); !os.IsNotExist(err) {
		t.Fatalf("shared-cache build promoted a reproducibility record: %v", err)
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
			source := committedHostileArtifactSource(t, root)
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

func TestSharedCacheBuildRestoresRecordOnInterruptedPromotion(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	source := committedHostileArtifactSource(t, root)
	prepared := copyPreparedArtifactGeneration(t, shared.outputDir)
	expected := shared.entryCount
	for _, test := range []struct {
		name  string
		abort func(*testing.T, string, string, string)
	}{
		{name: "signal at the seam never moves the record", abort: interruptArtifactPromotion},
		{name: "promotion failure after the record moves restores it", abort: failArtifactPromotionAfterRecordMove},
	} {
		t.Run(test.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "aborted shared cache output [*]")
			seedReady := filepath.Join(t.TempDir(), "seed-ready")
			seed := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
			seed.Env = promotionTestEnv(prepared, seedReady)
			if out, err := runArtifactBuildThroughPromotionSeam(t, seed, seedReady); err != nil {
				t.Fatalf("seed generation failed: %v\n%s", err, out)
			}
			record := filepath.Join(filepath.Dir(output), "reproducibility.json")
			want := []byte(staleReproducibilityRecord)
			if err := os.WriteFile(record, want, 0o644); err != nil {
				t.Fatal(err)
			}
			prior := promotedArtifactDigests(t, output)
			test.abort(t, source, prepared, output)
			if got, err := os.ReadFile(record); err != nil || !bytes.Equal(got, want) {
				t.Fatalf("aborted shared-cache promotion lost the prior record: %q, %v", got, err)
			}
			after := promotedArtifactDigests(t, output)
			if len(after) != expected || !maps.Equal(after, prior) {
				t.Fatalf("aborted shared-cache promotion changed the prior generation: got=%v want=%v", after, prior)
			}
		})
	}
}

// failArtifactPromotionAfterRecordMove deletes the staged generation while the build is
// parked at the promotion seam, so `mv "$artifacts" "$output"` fails only after the prior
// artifacts and the reproducibility record have already been moved aside. A signal delivered
// at that same seam exits above both moves, which is why only this abort reaches cleanup's
// restore branch.
func failArtifactPromotionAfterRecordMove(t *testing.T, source, prepared, output string) {
	t.Helper()
	ready := filepath.Join(t.TempDir(), "promotion-ready")
	cmd := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
	cmd.Env = promotionTestEnv(prepared, ready)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	awaitArtifactPromotionSeam(t, cmd, ready)
	removeStagedArtifactGeneration(t, filepath.Dir(output))
	if err := os.Remove(ready); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err == nil {
		t.Fatal("artifact builder promoted an absent staged generation")
	}
}

// removeStagedArtifactGeneration deletes the generation the parked build is about to
// promote. The prior-generation backup shares the .bench-artifacts. prefix but is created
// only after the seam releases, so while the build is parked the directory holding an
// "artifacts" child is unambiguously the live stage.
func removeStagedArtifactGeneration(t *testing.T, parent string) {
	t.Helper()
	stages, err := filepath.Glob(filepath.Join(parent, ".bench-artifacts.*"))
	if err != nil {
		t.Fatal(err)
	}
	staged := ""
	for _, stage := range stages {
		candidate := filepath.Join(stage, "artifacts")
		if info, err := os.Stat(candidate); err != nil || !info.IsDir() {
			continue
		}
		if staged != "" {
			t.Fatalf("ambiguous artifact staging directories under %s: %v", parent, stages)
		}
		staged = candidate
	}
	if staged == "" {
		t.Fatalf("no artifact staging directory under %s: %v", parent, stages)
	}
	if err := os.RemoveAll(staged); err != nil {
		t.Fatal(err)
	}
}

const staleReproducibilityRecord = "{\"schema_version\":1,\"status\":\"green\",\"builds\":2}\n"
