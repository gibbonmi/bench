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
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	staged := committedHostileArtifactSource(t, root)
	var plan struct {
		Targets []artifactPlatform `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(staged, "scripts", "release-plan.json"), &plan)
	goEnv, err := exec.Command("go", "env", "GOOS", "GOARCH").Output()
	if err != nil {
		t.Fatalf("read host Go target: %v", err)
	}
	host := strings.Fields(string(goEnv))
	if len(plan.Targets) != 1 || plan.Targets[0].OS != host[0] || plan.Targets[0].GOArch != host[1] {
		t.Fatalf("staged release plan targets = %+v, want exactly host %s/%s", plan.Targets, host[0], host[1])
	}
	out := filepath.Join(t.TempDir(), "host-only artifact output [*]")
	contract.NewExecFixtureAt(t, root).Run("bash", filepath.Join(staged, "scripts", "build-artifacts.sh"), staged, out).RequireExit(0)
	assertPlannedArtifactNames(t, staged, out)
}

func TestArtifactSourceSkipsWhenHostTargetIsAbsent(t *testing.T) {
	const child = "BENCH_TEST_HOSTLESS_ARTIFACT_PLAN"
	if os.Getenv(child) == "1" {
		committedHostileArtifactSource(t, contract.SubjectRoot(t), stageHostlessArtifactPlan)
		t.Fatal("hostless staged plan did not skip")
	}
	goEnv, err := exec.Command("go", "env", "GOOS", "GOARCH").Output()
	if err != nil {
		t.Fatalf("read host Go target: %v", err)
	}
	host := strings.Fields(string(goEnv))
	command := exec.Command(os.Args[0], "-test.run=^TestArtifactSourceSkipsWhenHostTargetIsAbsent$", "-test.v")
	command.Env = append(os.Environ(), child+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("hostless artifact staging subprocess failed: %v\n%s", err, output)
	}
	want := "artifact contract tests require release plan target for host " + host[0] + "/" + host[1]
	if !strings.Contains(string(output), "--- SKIP: TestArtifactSourceSkipsWhenHostTargetIsAbsent") || !strings.Contains(string(output), want) {
		t.Fatalf("hostless artifact staging did not report named skip %q:\n%s", want, output)
	}
}
