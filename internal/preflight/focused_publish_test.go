package preflight

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublishIdentityFailurePreventsArtifactConstruction(t *testing.T) {
	root := preflightRepo(t)
	marker := filepath.Join(t.TempDir(), "artifact-construction-ran")
	phase := filepath.Join(t.TempDir(), "artifact-phase")
	if err := os.WriteFile(phase, []byte("#!/bin/sh\n: > \"$BENCH_ORDER_MARKER\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_REF", "not-a-release-tag")
	t.Setenv("BENCH_PREFLIGHT_ARTIFACTS", phase)
	t.Setenv("BENCH_ORDER_MARKER", marker)
	results := (&runner{root: root, mode: ModePublish, binaryVersion: "0.2.0"}).run(context.Background(), "")
	if len(results) == 0 || results[0].Name != "identity" || results[0].Status != StatusRed {
		t.Fatalf("publish did not fail first at identity: %+v", results)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("artifact construction ran before release identity authorization: %v", err)
	}
}

func TestBuiltCommandFocusedPublishRunsDiagnosticWithoutAuthorizing(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(projectRoot(t), "scripts", "go-build.sh"), projectRoot(t), binary)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	root := preflightRepo(t)
	full := exec.Command(binary, "release-preflight", "--mode", "verify")
	full.Dir = root
	if output, err := full.CombinedOutput(); err != nil {
		t.Fatalf("initial full verify: %v\n%s", err, output)
	}
	prior, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "focused-publish-ran")
	phase := filepath.Join(root, "focused-publish")
	if err := os.WriteFile(phase, []byte("#!/bin/sh\nprintf 'ran\\n' > \"$BENCH_FOCUSED_MARKER\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binary, "release-preflight", "--mode", "publish", "--profile", "public", "--phase", "gate")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "BENCH_PREFLIGHT_GATE="+phase, "BENCH_FOCUSED_MARKER="+marker)
	output, err := cmd.CombinedOutput()
	if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		t.Fatalf("focused publish exit = %v, want non-authorizing exit 1\n%s", err, output)
	}
	if !strings.Contains(string(output), "focused publish runs cannot authorize publication") {
		t.Fatalf("focused publish output does not explain non-authorization:\n%s", output)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "ran\n" {
		t.Fatalf("focused publish diagnostic did not run: %q %v", data, err)
	}
	after, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(prior, after) {
		t.Fatal("focused publish replaced the prior complete authoritative generation")
	}
}
