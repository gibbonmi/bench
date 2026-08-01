package gate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExecuteTreeBuildsExactUnpublishedBenchkitSource(t *testing.T) {
	kit := kitRootForTest(t)
	root := gateTestRepo(t, string(mustReadGateTestFile(t, filepath.Join(kit, ".bench", "gate.sh"))), `{"schema":1,"closure":"local","environment":["HOME"],"paths":[],"tools":[]}`)
	for _, rel := range []string{
		".bench/gate-prospective.sh", "scripts/go-build.sh", "scripts/go-build.inputs",
		"package.json", "internal/releaseevidence/requirements.json", "internal/freshness/freshness.go",
		"internal/freshness/check/main.go", "internal/freshness/cmd/main.go",
	} {
		mode := os.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") {
			mode = 0o755
		}
		writeGateTestFile(t, root, rel, string(mustReadGateTestFile(t, filepath.Join(kit, filepath.FromSlash(rel)))), mode)
	}
	writeGateTestFile(t, root, ".gitignore", "dist/\n", 0o644)
	writeGateTestFile(t, root, "go.mod", "module github.com/gibbonmi/bench\n\ngo 1.25\n", 0o644)
	writeGateTestFile(t, root, "cmd/bench/main.go", prospectiveBenchMain("source A"), 0o644)
	gitRun(t, root, "add", ".")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "source A")
	writeGateTestFile(t, root, "cmd/bench/main.go", prospectiveBenchMain("source B"), 0o644)
	gitRun(t, root, "add", "cmd/bench/main.go")
	tree := gitOutput(t, root, "write-tree")
	gitRun(t, root, "reset", "--hard", "HEAD")

	direct := exec.Command(filepath.Join(root, ".bench", "gate.sh"))
	direct.Dir = root
	direct.Env = append(os.Environ(), "BENCH_GATE_PROSPECTIVE=1", "GOCACHE="+t.TempDir())
	directOutput, directErr := direct.CombinedOutput()
	if directErr == nil || !strings.Contains(string(directOutput), "rebuild with") {
		t.Fatalf("ordinary real wrapper with ambient marker = %v, output=%q; want missing-artifact freshness refusal", directErr, directOutput)
	}

	var stdout, stderr bytes.Buffer
	if got := ExecuteTree(context.Background(), root, tree, &stdout, &stderr); got.ActionExit != 0 {
		t.Fatalf("prospective benchkit execution = %+v, want green; stdout=%q stderr=%q", got, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "source B") || strings.Contains(stdout.String(), "source A") {
		t.Fatalf("prospective output = %q, want only unpublished source B", stdout.String())
	}
	if _, err := os.Lstat(filepath.Join(root, "dist")); !os.IsNotExist(err) {
		t.Fatalf("prospective execution populated ordinary checkout dist: %v", err)
	}
	writeGateTestFile(t, root, "cmd/bench/main.go", prospectiveBenchMain("source B"), 0o644)
	gitRun(t, root, "add", "cmd/bench/main.go")
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "source B")
	if committed := gitOutput(t, root, "show", "-s", "--format=%T", "HEAD"); committed != tree {
		t.Fatalf("committed source B tree = %s, want gated tree %s", committed, tree)
	}
	branch := gitOutput(t, root, "branch", "--show-current")
	gitRun(t, root, "update-ref", "refs/bench/green/"+branch, "HEAD")
	if got := ValidateProjectGreen(root, branch); !got.ReusableGreen {
		t.Fatalf("prospective green did not validate committed project-green: %+v", got)
	}
}

func TestOrdinaryGreenRemainsProspectiveBootstrapEvidence(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	writeGateTestFile(t, root, prospectiveGatePath, "#!/usr/bin/env bash\nexit 97\n", 0o755)
	gitRun(t, root, "add", prospectiveGatePath)
	gitRun(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-m", "prospective hook")
	tree := gitOutput(t, root, "write-tree")

	if got := Execute(context.Background(), root, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("ordinary execution = %+v, want green", got)
	}
	if got := InspectTree(root, tree); !got.ReusableGreen {
		t.Fatalf("ordinary green did not remain prospective bootstrap evidence: %+v", got)
	}
	if got := ExecuteTree(context.Background(), root, tree, io.Discard, io.Discard); got.ActionExit != 0 {
		t.Fatalf("prospective bootstrap did not reuse ordinary green: %+v", got)
	}
}

func TestPolicyVersionMismatchInvalidatesGreen(t *testing.T) {
	root := reusableEvidenceRepo(t, 0)
	old, err := buildSubjectForPolicy(root, root, "oracle-v1/freshness-v1")
	if err != nil {
		t.Fatal(err)
	}
	current, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	if old.Oracle == current.Oracle {
		t.Fatal("policy version did not change oracle identity")
	}
	gitdir := gitOutput(t, root, "rev-parse", "--absolute-git-dir")
	if err := durableReplace(gitdir, verdictRecord{Schema: 1, State: Ready, Status: "green", Tree: old.Tree, Oracle: old.Oracle, RecordedAt: time.Now().UTC().Format(time.RFC3339)}); err != nil {
		t.Fatal(err)
	}
	if got := Inspect(root); got.ReusableGreen || got.Reason != "oracle changed" {
		t.Fatalf("policy-mismatched green = %+v, want oracle-changed refusal", got)
	}
}

func prospectiveBenchMain(sentinel string) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

var version string

func main() {
	if len(os.Args) != 3 {
		os.Exit(90)
	}
	switch os.Args[1] {
	case "freshness-check":
		return
	case "gate-phases":
		fmt.Println(%q)
	default:
		os.Exit(91)
	}
}
`, sentinel)
}

func mustReadGateTestFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
