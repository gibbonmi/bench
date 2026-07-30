package offline

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
)

func TestReleasePlanProjectsDerivedArchiveInventory(t *testing.T) {
	root := contract.SubjectRoot(t)
	command := exec.Command("node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "archive-inventory", "linux-x64", "0.1.0")
	data, err := command.Output()
	if err != nil {
		t.Fatalf("canonical release plan did not project archive inventory: %v", err)
	}
	var inventory map[string]int64
	if json.Unmarshal(data, &inventory) != nil || inventory["bin/bench"] != 0o755 || inventory["packages/redbench-0.1.0.tgz"] != 0o644 {
		t.Fatalf("canonical archive inventory is incomplete: %s", data)
	}
}

func TestReleasePlanProjectsRelocatedPackageEvidence(t *testing.T) {
	root := contract.SubjectRoot(t)
	fixture := t.TempDir()
	for _, relative := range []string{"scripts/release-plan.mjs", "scripts/release-plan.json", "internal/releaseevidence/requirements.json"} {
		data, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatal(err)
		}
		if relative == "scripts/release-plan.json" {
			data = bytes.Replace(data, []byte("evidence/{package_evidence}"), []byte("relocated/{package_evidence}"), 1)
		}
		target := filepath.Join(fixture, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil || os.WriteFile(target, data, 0o644) != nil {
			t.Fatalf("copy relocated release-plan fixture: %v", err)
		}
	}
	var requirements struct {
		Records []struct {
			Key         string `json:"key"`
			Path        string `json:"path"`
			PackageMode string `json:"package_mode"`
		} `json:"records"`
	}
	contract.ReadJSONFile(t, filepath.Join(fixture, "internal", "releaseevidence", "requirements.json"), &requirements)
	for _, record := range requirements.Records {
		if record.PackageMode == "" {
			continue
		}
		command := exec.Command("node", filepath.Join(fixture, "scripts", "release-plan.mjs"), fixture, "archive-evidence-path", record.Key, "linux-x64", "0.1.0")
		data, err := command.Output()
		if err != nil || strings.TrimSpace(string(data)) != "relocated/"+record.Path {
			t.Fatalf("relocated package evidence %s was not projected from the release plan: %v\n%s", record.Key, err, data)
		}
	}
	archiveDir := assembleOfflineArchiveStage(t, root, fixture, root)
	for _, record := range requirements.Records {
		if record.PackageMode == "" {
			continue
		}
		want, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(record.Path)))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(archiveDir, "relocated", filepath.FromSlash(record.Path)))
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("relocated package evidence %s was not assembled from %s: %v", record.Key, record.Path, err)
		}
	}
}

func TestNativeProofAggregatorRejectsOneTargetOmission(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/aggregate-native-proofs.sh")
	var plan struct {
		Targets []struct {
			OS     string `json:"os"`
			Arch   string `json:"arch"`
			Runner string `json:"runner"`
		} `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "release-plan.json"), &plan)
	if len(plan.Targets) < 2 {
		t.Fatal("release plan needs multiple targets for omission coverage")
	}
	proofs := t.TempDir()
	for _, target := range plan.Targets[:len(plan.Targets)-1] {
		writeAggregateProof(t, proofs, target.OS, target.Arch, target.Runner, false)
	}
	command := exec.Command("bash", filepath.Join(root, "scripts", "aggregate-native-proofs.sh"), proofs)
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("one-target omission passed native proof aggregation:\n%s", output)
	}
	missing := plan.Targets[len(plan.Targets)-1]
	if !strings.Contains(string(output), "native proof set is missing "+missing.OS+"/"+missing.Arch) {
		t.Fatalf("one-target omission was not attributed:\n%s", output)
	}
}

func TestNativeProofAggregatorRejectsDigestMismatch(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/aggregate-native-proofs.sh")
	var plan struct {
		Targets []struct{ OS, Arch, Runner string } `json:"targets"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "scripts", "release-plan.json"), &plan)
	proofs := t.TempDir()
	for index, target := range plan.Targets {
		writeAggregateProof(t, proofs, target.OS, target.Arch, target.Runner, index == 0)
	}
	command := exec.Command("bash", filepath.Join(root, "scripts", "aggregate-native-proofs.sh"), proofs)
	if output, err := command.CombinedOutput(); err == nil {
		t.Fatalf("digest binding mutation passed native proof aggregation:\n%s", output)
	}
}

func TestAuthoritativeNativeProofBehaviorCanary(t *testing.T) {
	root := contract.SubjectRoot(t)
	marker := filepath.Join(root, "internal", "releaseevidence", "native_canary_test.go")
	if _, err := os.Stat(marker); os.IsNotExist(err) {
		capability.Environment(t, "authoritative native-proof canary fixture is not materialized")
	}
	command := exec.Command("go", "test", "-count=1", "-tags=bench_canary_native_proof", "-run", "^TestNativeProofAuthorizationCanary$", "./internal/releaseevidence")
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("authoritative native proof canary subprocess failed: %v\n%s", err, output)
	}
}

func writeAggregateProof(t *testing.T, directory, osName, arch, runner string, mismatch bool) {
	t.Helper()
	digest, binary := strings.Repeat("1", 64), strings.Repeat("1", 64)
	if mismatch {
		binary = strings.Repeat("2", 64)
	}
	musl := "not_applicable"
	if osName == "linux" {
		musl = "green"
	}
	name := osName + "-" + arch
	proof := map[string]any{"schema_version": 1, "target": name, "runner": runner, "status": "green", "rebuilt_sha256": digest, "binary_sha256": binary, "package_sha256": digest, "archive_sha256": digest, "musl_status": musl, "operations_status": "green", "strip_status": "green", "tools_status": "green"}
	data, err := json.Marshal(proof)
	if err != nil || os.WriteFile(filepath.Join(directory, name+".json"), append(data, '\n'), 0o644) != nil {
		t.Fatalf("write native proof: %v", err)
	}
}

func TestOfflineRegistryDerivesAcceptedTargetsFromReleasePlan(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/offline-registry.mjs")
	fixture := t.TempDir()
	for _, relative := range []string{"scripts/offline-registry.mjs", "scripts/release-plan.mjs"} {
		data, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(fixture, relative)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil || os.WriteFile(target, data, 0o644) != nil {
			t.Fatalf("copy registry fixture: %v", err)
		}
	}
	plan := `{"schema_version":1,"targets":[{"os":"haiku","arch":"riscv64","goos":"haiku","goarch":"riscv64","runner":"changed-runner"}],"archive_entries":[]}` + "\n"
	if err := os.WriteFile(filepath.Join(fixture, "scripts", "release-plan.json"), []byte(plan), 0o644); err != nil {
		t.Fatal(err)
	}
	store := filepath.Join(fixture, "store")
	if err := os.MkdirAll(store, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"redbench-haiku-riscv64-1.2.3.tgz", "redbench-darwin-x64-1.2.3.tgz"} {
		if err := os.WriteFile(filepath.Join(store, name), []byte(name+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	portFile, requestFile := filepath.Join(fixture, "port"), filepath.Join(fixture, "requests")
	server := exec.Command("node", filepath.Join(fixture, "scripts", "offline-registry.mjs"), store, portFile, requestFile)
	if err := server.Start(); err != nil {
		t.Fatalf("start registry: %v", err)
	}
	t.Cleanup(func() {
		_ = server.Process.Signal(os.Interrupt)
		_ = server.Wait()
	})
	var port string
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(portFile); err == nil {
			port = strings.TrimSpace(string(data))
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if port == "" {
		t.Fatal("offline registry did not publish its port")
	}
	client := &http.Client{Timeout: 2 * time.Second}
	for _, test := range []struct {
		name string
		want int
	}{{"haiku-riscv64", http.StatusOK}, {"darwin-x64", http.StatusNotFound}} {
		response, err := client.Get("http://127.0.0.1:" + port + "/%40redbench%2F" + test.name)
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()
		if response.StatusCode != test.want {
			t.Fatalf("registry target %s status = %d, want %d", test.name, response.StatusCode, test.want)
		}
	}
}

func TestOfflineArchiveBuildRefusesOutputItCannotAccountFor(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-offline-archives.sh")
	script := filepath.Join(root, "scripts", "build-offline-archives.sh")
	const refusal = "did not produce"
	for _, test := range []struct {
		name    string
		seed    func(t *testing.T, dir string)
		refused bool
	}{
		{"a live checkout", func(t *testing.T, dir string) {
			if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "KEEP.txt"), []byte("committed work\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}, true},
		{"the npm tarballs it consumes", func(t *testing.T, dir string) {
			if err := os.WriteFile(filepath.Join(dir, "redbench-0.2.0.tgz"), []byte("wrapper\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "offline out [hostile]")
			if err := os.MkdirAll(output, 0o755); err != nil {
				t.Fatal(err)
			}
			test.seed(t, output)
			probe := contract.NewExecFixtureAt(t, root).Run("bash", script, t.TempDir(), output)
			probe.RequireExit(1)
			if test.refused {
				probe.RequireContains(probe.Stderr, refusal)
			} else {
				probe.RequireNotContains(probe.Stderr, refusal)
			}
			entries, err := os.ReadDir(output)
			if err != nil {
				t.Fatalf("read output directory after refusal: %v", err)
			}
			if len(entries) == 0 {
				t.Fatal("the build emptied an output directory it did not complete")
			}
		})
	}
}

func TestOfflineInstructionsVerifyOnlyTargetArchive(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/assemble-offline-archive.mjs")
	archiveDir := assembleOfflineArchiveStage(t, root, root, root)
	instructions, err := os.ReadFile(filepath.Join(archiveDir, "OFFLINE.md"))
	if err != nil {
		t.Fatal(err)
	}
	const targeted = `awk '$2 == "redbench-0.1.0-linux-x64.tar.gz" { print }' SHA256SUMS | sha256sum -c -`
	if !bytes.Contains(instructions, []byte(targeted)) || bytes.Contains(instructions, []byte("sha256sum -c SHA256SUMS")) {
		t.Fatalf("offline instructions do not verify only the selected target archive:\n%s", instructions)
	}
}

func assembleOfflineArchiveStage(t *testing.T, scriptRoot, planRoot, evidenceRoot string) string {
	t.Helper()
	stage := t.TempDir()
	npmDir, archiveDir := filepath.Join(stage, "packages"), filepath.Join(stage, "archive")
	wrapperExtract, platformExtract := filepath.Join(stage, "wrapper"), filepath.Join(stage, "platform")
	for _, directory := range []string{npmDir, archiveDir, filepath.Join(wrapperExtract, "package"), filepath.Join(platformExtract, "package")} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(path string, data []byte, mode os.FileMode) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, mode); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(stage, "bench"), []byte("fixture binary\n"), 0o755)
	write(filepath.Join(npmDir, "redbench-0.1.0.tgz"), []byte("wrapper\n"), 0o644)
	write(filepath.Join(npmDir, "redbench-linux-x64-0.1.0.tgz"), []byte("platform\n"), 0o644)
	write(filepath.Join(wrapperExtract, "package", "component-manifest.json"), []byte("{}\n"), 0o644)
	write(filepath.Join(platformExtract, "package", "component-manifest.json"), []byte("{}\n"), 0o644)
	var requirements struct {
		Records []struct {
			Path        string `json:"path"`
			PackageMode string `json:"package_mode"`
		} `json:"records"`
	}
	contract.ReadJSONFile(t, filepath.Join(planRoot, "internal", "releaseevidence", "requirements.json"), &requirements)
	for _, record := range requirements.Records {
		if record.PackageMode != "" {
			data, err := os.ReadFile(filepath.Join(evidenceRoot, filepath.FromSlash(record.Path)))
			if err != nil {
				t.Fatal(err)
			}
			write(filepath.Join(wrapperExtract, "package", filepath.FromSlash(record.Path)), data, 0o644)
		}
	}
	contract.NewExecFixtureAt(t, scriptRoot).Run("node", filepath.Join(scriptRoot, "scripts", "assemble-offline-archive.mjs"), planRoot, npmDir, archiveDir, "linux-x64", "0.1.0", filepath.Join(stage, "bench"), wrapperExtract, platformExtract).RequireExit(0)
	return archiveDir
}
