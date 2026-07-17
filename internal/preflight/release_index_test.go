package preflight

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleaseEvidenceIsDeterministicBoundAndIdempotent(t *testing.T) {
	root := preflightRepo(t)
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("first verify exit=%d stderr=%s", code, stderr.String())
	}
	indexPath := filepath.Join(root, "dist", "preflight", "release-index.json")
	sumsPath := filepath.Join(root, "dist", "preflight", "SHA256SUMS")
	firstIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	firstSums, err := os.ReadFile(sumsPath)
	if err != nil {
		t.Fatal(err)
	}
	var index releaseIndex
	if err := json.Unmarshal(firstIndex, &index); err != nil {
		t.Fatal(err)
	}
	goSumData, err := os.ReadFile(filepath.Join(root, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	goSumBound := false
	for _, input := range index.Inputs {
		if input.Path == "go.sum" {
			goSumBound = input.SHA256 == sha256Hex(goSumData)
		}
	}
	if !goSumBound {
		t.Fatal("release index does not bind go.sum")
	}
	wantVersions := map[string]string{}
	for _, command := range [][]string{{"go", "env", "GOVERSION"}, {"node", "--version"}, {"npm", "--version"}} {
		out, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			t.Fatal(err)
		}
		wantVersions[command[0]] = strings.TrimSpace(string(out))
	}
	if len(index.Toolchains) != 3 {
		t.Fatalf("release index toolchains=%+v, want exact Go, Node, and npm observations", index.Toolchains)
	}
	for _, toolchain := range index.Toolchains {
		if toolchain.Version != wantVersions[toolchain.Name] {
			t.Fatalf("release index toolchain %s version=%q, want %q", toolchain.Name, toolchain.Version, wantVersions[toolchain.Name])
		}
		if (toolchain.Name == "go" || toolchain.Name == "npm") && len(toolchain.Flags) == 0 {
			t.Fatalf("release index toolchain %s omits release flags", toolchain.Name)
		}
	}
	sumByName := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(firstSums)), "\n") {
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			t.Fatalf("malformed checksum line %q", line)
		}
		sumByName[parts[1]] = parts[0]
	}
	if len(sumByName) != len(index.Artifacts) {
		t.Fatalf("checksum count=%d artifacts=%d", len(sumByName), len(index.Artifacts))
	}
	for _, artifact := range index.Artifacts {
		data, err := os.ReadFile(filepath.Join(root, "dist", "artifacts", artifact.Name))
		if err != nil {
			t.Fatal(err)
		}
		if artifact.SHA256 != sha256Hex(data) || sumByName[artifact.Name] != artifact.SHA256 {
			t.Fatalf("artifact digest binding failed for %s", artifact.Name)
		}
	}
	artifactDir := filepath.Join(root, "dist", "artifacts")
	staging := t.TempDir()
	entries, err := os.ReadDir(artifactDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if err := os.Rename(filepath.Join(artifactDir, entry.Name()), filepath.Join(staging, entry.Name())); err != nil {
			t.Fatal(err)
		}
	}
	for i := len(entries) - 1; i >= 0; i-- {
		if err := os.Rename(filepath.Join(staging, entries[i].Name()), filepath.Join(artifactDir, entries[i].Name())); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("LANG", "C")
	t.Setenv("TZ", "UTC")
	stderr.Reset()
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("second verify exit=%d stderr=%s", code, stderr.String())
	}
	secondIndex, _ := os.ReadFile(indexPath)
	secondSums, _ := os.ReadFile(sumsPath)
	if string(secondIndex) != string(firstIndex) || string(secondSums) != string(firstSums) {
		t.Fatal("release evidence changed with enumeration order or environment")
	}

	failing := filepath.Join(root, "fail-gate")
	if err := os.WriteFile(failing, []byte("#!/bin/sh\nexit 9\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	greenGate := os.Getenv("BENCH_PREFLIGHT_GATE")
	if err := os.Setenv("BENCH_PREFLIGHT_GATE", failing); err != nil {
		t.Fatal(err)
	}
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("red rerun exit=%d stderr=%s", code, stderr.String())
	}
	redIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	var red releaseIndex
	if err := json.Unmarshal(redIndex, &red); err != nil || red.Status != StatusRed {
		t.Fatalf("red rerun index=%s err=%v", redIndex, err)
	}
	if err := os.Setenv("BENCH_PREFLIGHT_GATE", greenGate); err != nil {
		t.Fatal(err)
	}
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("green-after-red exit=%d stderr=%s", code, stderr.String())
	}
	var final releaseIndex
	finalIndex, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(finalIndex, &final); err != nil || final.Status != StatusGreen {
		t.Fatalf("green-after-red index=%s err=%v", finalIndex, err)
	}
	files, err := os.ReadDir(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(PhaseNames(ModeVerify))+3 {
		t.Fatalf("promoted evidence file count=%d, want %d", len(files), len(PhaseNames(ModeVerify))+3)
	}
}

func TestReleasePolicyFailureClassesAreRed(t *testing.T) {
	root := preflightRepo(t)
	r := &runner{root: root, mode: ModePublish, binaryVersion: "0.2.0", stderr: &bytes.Buffer{}}
	if err := r.populateBaseIdentity(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_REF", "refs/tags/v0.2.0-beta.1")
	if err := r.checkIdentity(context.Background()); err == nil {
		t.Fatal("prerelease tag passed")
	}
	if err := r.checkAncestry(context.Background()); err == nil {
		t.Fatal("missing origin/main ancestry passed")
	}
	tag := "v0.2.0"
	r.identity.Tag = &tag
	for name, body := range map[string]string{
		"duplicate":          "## [0.2.0] - 2026-07-16\n## [0.2.0] - 2026-07-17\n",
		"legacy":             "## v0.2.0 (2026-07-16)\n",
		"invalid date":       "## [Unreleased]\n\n## [0.2.0] - 2026-02-30\n",
		"missing unreleased": "## [0.2.0] - 2026-07-16\n",
	} {
		if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := r.checkChangelog(); err == nil {
			t.Fatalf("%s changelog passed", name)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ntoolchain go1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readToolchain(root); err == nil {
		t.Fatal("non-patch toolchain passed")
	}
}
