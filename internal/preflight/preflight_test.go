package preflight

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
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

func TestBankPublishRequiresFT71Evidence(t *testing.T) {
	root := preflightRepo(t)
	if err := os.Remove(filepath.Join(root, "release-evidence", "ft71-local-event.json")); err != nil {
		t.Fatal(err)
	}
	tagRelease(t, root, true)
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "publish", "--profile", "bank"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("bank publish exit=%d stderr=%s", code, stderr.String())
	}
	data, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "release-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var index releaseIndex
	if err := json.Unmarshal(data, &index); err != nil {
		t.Fatal(err)
	}
	if index.Status != StatusRed || requirementIndexStatus(index.Requirements, "bank.ft71.local_event") != "missing" {
		t.Fatalf("bank index = %+v", index)
	}
}

func TestGovernanceSchemaRejectsUnknownFieldsAndVersions(t *testing.T) {
	for _, test := range []struct {
		name string
		body string
	}{
		{name: "unknown field", body: `{"schema_version":1,"policy":"support","route":"GitHub Issues","personal_email":false,"non_personal":true,"unexpected":true}` + "\n"},
		{name: "unknown version", body: `{"schema_version":2,"policy":"support","route":"GitHub Issues","personal_email":false,"non_personal":true}` + "\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := preflightRepo(t)
			if err := os.WriteFile(filepath.Join(root, "governance", "policies", "support.json"), []byte(test.body), 0o644); err != nil {
				t.Fatal(err)
			}
			old, _ := os.Getwd()
			if err := os.Chdir(root); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chdir(old) })
			var stderr bytes.Buffer
			if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
				t.Fatalf("exit=%d stderr=%s", code, stderr.String())
			}
			data, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "release-index.json"))
			if err != nil {
				t.Fatal(err)
			}
			var index releaseIndex
			if err := json.Unmarshal(data, &index); err != nil {
				t.Fatal(err)
			}
			if requirementIndexStatus(index.Requirements, "core.policy.support") != "invalid" {
				t.Fatalf("support status = %q, want invalid", requirementIndexStatus(index.Requirements, "core.policy.support"))
			}
		})
	}
}

func TestFinalTarRejectsHostileEvidenceWithConsistentInventory(t *testing.T) {
	for _, test := range []struct {
		name string
		path string
		data []byte
		mode int64
		want string
	}{
		{name: "empty notices", path: "governance/THIRD_PARTY_NOTICES.txt", data: []byte{}, mode: 0o644, want: "empty or missing"},
		{name: "unsafe mode", path: "LICENSE", data: []byte("license\n"), mode: 0o777, want: "unsafe mode"},
		{name: "malformed SPDX", path: "governance/sbom.spdx.json", data: []byte(`{"SPDXVersion":"SPDX-2.3","unexpected":true}` + "\n"), mode: 0o644, want: "malformed SPDX JSON"},
		{name: "malformed policy", path: "governance/policies/support.json", data: []byte(`{"schema_version":1,"policy":"support","route":"GitHub Issues","personal_email":false,"non_personal":true,"unexpected":true}` + "\n"), mode: 0o644, want: "malformed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := preflightRepo(t)
			archive := filepath.Join(root, "dist", "artifacts", "redbench-0.2.0.tgz")
			rewriteFixtureTarball(t, archive, func(files map[string]fixtureArchiveFile) {
				file := files[test.path]
				file.data, file.mode = test.data, test.mode
				files[test.path] = file
			})
			old, _ := os.Getwd()
			if err := os.Chdir(root); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chdir(old) })
			var stderr bytes.Buffer
			if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
				t.Fatalf("exit=%d stderr=%s", code, stderr.String())
			}
			if !strings.Contains(stderr.String(), test.want) {
				t.Fatalf("output does not contain %q:\n%s", test.want, stderr.String())
			}
		})
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

func TestFullRunPhaseRedPromotesCompleteReleaseEvidence(t *testing.T) {
	root := preflightRepo(t)
	failing := filepath.Join(root, "fail-gate")
	if err := os.WriteFile(failing, []byte("#!/bin/sh\nexit 9\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_PREFLIGHT_GATE", failing)
	old, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	indexData, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "release-index.json"))
	if err != nil {
		t.Fatalf("complete red release index missing: %v\n%s", err, stderr.String())
	}
	var index releaseIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		t.Fatal(err)
	}
	if index.Status != StatusRed || len(index.Requirements) != len(Requirements()) || len(index.Artifacts) != 5 {
		t.Fatalf("red index = %+v", index)
	}
	if sums, err := os.ReadFile(filepath.Join(root, "dist", "preflight", "SHA256SUMS")); err != nil || len(sums) == 0 {
		t.Fatalf("complete red checksums missing: %v", err)
	}
}

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
