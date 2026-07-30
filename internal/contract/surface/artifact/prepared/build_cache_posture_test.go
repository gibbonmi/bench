package prepared

import fixture "github.com/gibbonmi/bench/internal/contract/surface/artifact/internal/fixture"

import (
	"bytes"
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

func TestSharedCacheBuildPromotesNoRecord(t *testing.T) {
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	output := requireSharedArtifactSet(t).outputDir
	if _, err := os.Stat(filepath.Join(filepath.Dir(output), "reproducibility.json")); !os.IsNotExist(err) {
		t.Fatalf("shared-cache build promoted a reproducibility record: %v", err)
	}
}

func TestSharedCacheBuildRestoresRecordOnInterruptedPromotion(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	source := fixture.CommittedHostileArtifactSource(t, root)
	prepared := fixture.CopyPreparedArtifactGeneration(t, shared.outputDir)
	expected := shared.entryCount
	for _, test := range []struct {
		name  string
		abort func(*testing.T, string, string, string)
	}{
		{name: "signal at the seam never moves the record", abort: fixture.InterruptArtifactPromotion},
		{name: "promotion failure after the record moves restores it", abort: failArtifactPromotionAfterRecordMove},
	} {
		t.Run(test.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "aborted shared cache output [*]")
			seedReady := filepath.Join(t.TempDir(), "seed-ready")
			seed := exec.Command("bash", filepath.Join(source, "scripts", "build-artifacts.sh"), source, output)
			seed.Env = fixture.PromotionTestEnv(prepared, seedReady)
			if out, err := fixture.RunArtifactBuildThroughPromotionSeam(t, seed, seedReady); err != nil {
				t.Fatalf("seed generation failed: %v\n%s", err, out)
			}
			record := filepath.Join(filepath.Dir(output), "reproducibility.json")
			want := []byte(staleReproducibilityRecord)
			if err := os.WriteFile(record, want, 0o644); err != nil {
				t.Fatal(err)
			}
			prior := fixture.PromotedArtifactDigests(t, output)
			test.abort(t, source, prepared, output)
			if got, err := os.ReadFile(record); err != nil || !bytes.Equal(got, want) {
				t.Fatalf("aborted shared-cache promotion lost the prior record: %q, %v", got, err)
			}
			after := fixture.PromotedArtifactDigests(t, output)
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
	cmd.Env = fixture.PromotionTestEnv(prepared, ready)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	fixture.AwaitArtifactPromotionSeam(t, cmd, ready)
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
