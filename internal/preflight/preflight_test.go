package preflight

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestRegistryDerivesPublishFromVerify(t *testing.T) {
	wantVerify := []string{"gate", "race", "vet", "vulnerability", "artifacts", "smoke"}
	wantPublish := append(append([]string{}, wantVerify...), "identity", "ancestry", "changelog")
	if got := PhaseNames(ModeVerify); !reflect.DeepEqual(got, wantVerify) {
		t.Fatalf("verify phases = %v, want %v", got, wantVerify)
	}
	if got := PhaseNames(ModePublish); !reflect.DeepEqual(got, wantPublish) {
		t.Fatalf("publish phases = %v, want %v", got, wantPublish)
	}
}

func TestEvidenceReplacesPriorCompleteRun(t *testing.T) {
	root := t.TempDir()
	old := filepath.Join(root, "dist", "preflight")
	if err := os.MkdirAll(old, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "stale.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	results := []Result{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}
	manifest := Manifest{SchemaVersion: 1, Mode: ModeVerify, Scope: ScopeFocused, Status: StatusGreen, Identity: Identity{}, Phases: phaseSummaries(results)}
	if err := PromoteEvidence(root, ModeVerify, results, manifest); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(old, "stale.json")); !os.IsNotExist(err) {
		t.Fatalf("stale record survived replacement: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(old, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 6 || got["scope"] != string(ScopeFocused) {
		t.Fatalf("manifest = %#v", got)
	}
}

func TestVulnerabilityPolicyRejectsUnusedAndUncovered(t *testing.T) {
	now := "2026-07-16"
	valid := []byte(`[{"id":"GO-1","reason":"temporary","expires":"2026-07-16"}]`)
	if err := ValidateVulnerabilityPolicy([]string{"GO-1"}, valid, true, now); err != nil {
		t.Fatal(err)
	}
	if err := ValidateVulnerabilityPolicy(nil, valid, true, now); err == nil {
		t.Fatal("unused exception passed")
	}
	if err := ValidateVulnerabilityPolicy([]string{"GO-2"}, valid, true, now); err == nil {
		t.Fatal("uncovered finding passed")
	}
}

func TestCommandFullVerifyWritesSixGreenRecords(t *testing.T) {
	root := preflightRepo(t)
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(filepath.Join(root, "nested")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	data, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Status != StatusGreen || manifest.Scope != ScopePreflight || len(manifest.Phases) != 6 {
		t.Fatalf("manifest=%+v", manifest)
	}
}

func TestReleaseProfilesStayPendingInVerifyAndRedInPublish(t *testing.T) {
	root := preflightRepo(t)
	for _, rel := range []string{"release-evidence/ft88-data-handling.json", "release-evidence/ft87-offline-network-control.json"} {
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatal(err)
		}
	}
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("verify exit=%d stderr=%s", code, stderr.String())
	}
	indexData, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "release-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var index releaseIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"public.ft88.data_handling", "public.ft87.offline_network_control"} {
		if status := requirementIndexStatus(index.Requirements, key); status != "pending" {
			t.Fatalf("verify status for %s = %q", key, status)
		}
	}
	tagRelease(t, root, true)
	stderr.Reset()
	if code := Command([]string{"--mode", "publish", "--profile", "public"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("publish exit=%d stderr=%s", code, stderr.String())
	}
	indexData, err = os.ReadFile(filepath.Join(root, "dist", "preflight", "release-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(indexData, &index); err != nil {
		t.Fatal(err)
	}
	if index.Status != StatusRed || requirementIndexStatus(index.Requirements, "public.ft88.data_handling") != "missing" {
		t.Fatalf("publish index = %+v", index)
	}
}

func requirementIndexStatus(statuses []requirementStatus, key string) string {
	for _, status := range statuses {
		if status.Key == key {
			return status.Status
		}
	}
	return ""
}

func TestCommandPublishIsGreenStrictSuperset(t *testing.T) {
	root := preflightRepo(t)
	if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte(releaseChangelog("2026-07-16")), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"tag", "v0.2.0"}, {"update-ref", "refs/remotes/origin/main", "HEAD"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	t.Setenv("BENCH_PREFLIGHT_REF", "refs/tags/v0.2.0")
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "publish", "--profile", "public"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	data, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Phases) != 9 || manifest.Identity.Tag == nil || *manifest.Identity.Tag != "v0.2.0" || manifest.Identity.ChangelogHeading == nil {
		t.Fatalf("manifest=%+v", manifest)
	}
}

func TestArtifactFailureRecordsSmokeNotRun(t *testing.T) {
	root := preflightRepo(t)
	failing := filepath.Join(root, "fail-artifacts")
	if err := os.WriteFile(failing, []byte("#!/bin/sh\nexit 9\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_ARTIFACTS", failing)
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	data, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "smoke.json"))
	if err != nil {
		t.Fatal(err)
	}
	var record Record
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	if record.Status != StatusNotRun || record.ExitCode != nil {
		t.Fatalf("smoke record=%+v", record)
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
	if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte("## [0.2.0] - 2026-07-16\n## [0.2.0] - 2026-07-17\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.checkChangelog(); err == nil {
		t.Fatal("duplicate changelog heading passed")
	}
	if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte("## v0.2.0 (2026-07-16)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.checkChangelog(); err == nil {
		t.Fatal("legacy changelog heading passed")
	}
	if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte("## [Unreleased]\n\n## [0.2.0] - 2026-02-30\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.checkChangelog(); err == nil {
		t.Fatal("invalid changelog date passed")
	}
	if err := os.WriteFile(filepath.Join(root, "CHANGELOG.md"), []byte("## [0.2.0] - 2026-07-16\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.checkChangelog(); err == nil {
		t.Fatal("missing Unreleased heading passed")
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n\ntoolchain go1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := readToolchain(root); err == nil {
		t.Fatal("non-patch toolchain passed")
	}
}

func releaseChangelog(date string) string {
	return "# Changelog\n\n## [Unreleased]\n\n## [0.2.0] - " + date + "\n\n- release\n"
}

func TestCancellationAndUnsafePromotionFailClosed(t *testing.T) {
	root := preflightRepo(t)
	blocking := filepath.Join(root, "blocking")
	if err := os.WriteFile(blocking, []byte("#!/bin/sh\nsleep 30\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_GATE", blocking)
	r := &runner{root: root, mode: ModeVerify, stderr: &bytes.Buffer{}}
	if err := r.populateBaseIdentity(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(50 * time.Millisecond); cancel() }()
	started := time.Now()
	result := r.runPhase(ctx, "gate")
	if result.Status != StatusInterrupted || time.Since(started) > 2*time.Second {
		t.Fatalf("interrupted result=%+v elapsed=%v", result, time.Since(started))
	}
	escape := filepath.Join(t.TempDir(), "escape")
	if err := os.MkdirAll(filepath.Join(root, "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(escape, filepath.Join(root, "dist", "preflight")); err != nil {
		t.Fatal(err)
	}
	manifest := Manifest{SchemaVersion: 1, Mode: ModeVerify, Scope: ScopeFocused, Status: StatusGreen, Identity: Identity{}, Phases: []PhaseSummary{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}}
	if err := PromoteEvidence(root, ModeVerify, []Result{{Name: "gate", Status: StatusGreen, ExitCode: intPtr(0)}}, manifest); err == nil {
		t.Fatal("symlink promotion target passed")
	}
	if _, err := os.Stat(escape); !os.IsNotExist(err) {
		t.Fatalf("promotion escaped output root: %v", err)
	}
}
