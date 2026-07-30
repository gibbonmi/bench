package axi

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/contract/internal/freshnessfixture"
	"github.com/gibbonmi/bench/internal/freshness"
)

const axiFreshnessChildEnv = "BENCH_FT131_AXI_FRESHNESS_CHILD"

func TestRunBenchInDirRefusesStaleSubject(t *testing.T) {
	if os.Getenv(axiFreshnessChildEnv) == "1" {
		f := contract.NewFixture(t)
		probe := runBenchInDir(t, f, f.Root, "learnings")
		probe.RequireExit(0)
		probe.RequireContains(probe.Stdout, "otherwise acceptable AXI output")
		return
	}

	root := freshnessfixture.StaleSubject(t, "otherwise acceptable AXI output")
	command := exec.Command(os.Args[0], "-test.run=^TestRunBenchInDirRefusesStaleSubject$")
	command.Env = append(os.Environ(), axiFreshnessChildEnv+"=1", canary.SubjectRootEnv+"="+root)
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("stale AXI subject satisfied its assertion:\n%s", output)
	}
	if !strings.Contains(string(output), "rebuild with "+freshness.RebuildAction(root)) {
		t.Fatalf("stale AXI subject did not report freshness:\n%s", output)
	}
	if strings.Contains(string(output), "otherwise acceptable AXI output") {
		t.Fatalf("stale AXI subject ran before freshness refusal:\n%s", output)
	}
}
