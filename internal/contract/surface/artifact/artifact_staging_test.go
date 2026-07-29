package artifact

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestArtifactSourceStagesCommittedHostPlan(t *testing.T) {
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	staged := shared.sourceRoot
	var plan struct {
		Targets []contract.ReleaseTarget `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(staged, "scripts", "release-plan.json"), &plan)
	hostOS, hostArch := contract.GoEnvPair(t, "GOOS", "GOARCH")
	if len(plan.Targets) != 1 || plan.Targets[0].OS != hostOS || plan.Targets[0].GOArch != hostArch {
		t.Fatalf("staged release plan targets = %+v, want exactly host %s/%s", plan.Targets, hostOS, hostArch)
	}
	out := shared.outputDir
	assertPlannedArtifactNames(t, staged, out)
}

func TestArtifactSourceSkipsWhenHostTargetIsAbsent(t *testing.T) {
	const child = "BENCH_TEST_HOSTLESS_ARTIFACT_PLAN"
	if os.Getenv(child) == "1" {
		committedHostileArtifactSource(t, contract.SubjectRoot(t), stageHostlessArtifactPlan)
		t.Fatal("hostless staged plan did not skip")
	}
	hostOS, hostArch := contract.GoEnvPair(t, "GOOS", "GOARCH")
	command := exec.Command(os.Args[0], "-test.run=^TestArtifactSourceSkipsWhenHostTargetIsAbsent$", "-test.v")
	command.Env = append(os.Environ(), child+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("hostless artifact staging subprocess failed: %v\n%s", err, output)
	}
	want := "artifact contract tests require release plan target for host " + hostOS + "/" + hostArch
	if !strings.Contains(string(output), "--- SKIP: TestArtifactSourceSkipsWhenHostTargetIsAbsent") || !strings.Contains(string(output), want) {
		t.Fatalf("hostless artifact staging did not report named skip %q:\n%s", want, output)
	}
}
