package preflight

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func buildPreflightBinary(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(projectRoot(t), "scripts", "go-build.sh"), projectRoot(t), binary)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	return binary
}

func runBuilt(t *testing.T, binary, root string, env []string) ([]byte, error) {
	t.Helper()
	return runBuiltArgs(t, binary, root, env, "--mode", "verify")
}

func runBuiltArgs(t *testing.T, binary, root string, env []string, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command(binary, append([]string{"release-preflight"}, args...)...)
	cmd.Dir = root
	cmd.Env = env
	return cmd.CombinedOutput()
}

func TestBuiltCommandProfileAcceptanceMatrix(t *testing.T) {
	binary := buildPreflightBinary(t)
	for _, test := range []struct {
		name string
		args []string
	}{
		{name: "missing profile", args: []string{"--mode", "publish"}},
		{name: "unknown profile", args: []string{"--mode", "publish", "--profile", "unknown"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			output, err := runBuiltArgs(t, binary, preflightRepo(t), os.Environ(), test.args...)
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 || !strings.Contains(string(output), `"kind":"usage"`) {
				t.Fatalf("error=%v, want usage exit 2\n%s", err, output)
			}
		})
	}
	for _, profile := range []string{"public", "bank"} {
		t.Run("green "+profile, func(t *testing.T) {
			root := preflightRepo(t)
			tagRelease(t, root, true)
			output, err := runBuiltArgs(t, binary, root, os.Environ(), "--mode", "publish", "--profile", profile)
			if err != nil {
				t.Fatalf("green %s publish: %v\n%s", profile, err, output)
			}
		})
	}
}

func TestConditionalRecordReasonFlowsThroughAuthoritativeEvidence(t *testing.T) {
	root := preflightRepo(t)
	registry := requirements
	registry.Records = append([]Requirement(nil), requirements.Records...)
	for i := range registry.Records {
		if registry.Records[i].Key == "bank.ft71.local_event" {
			registry.Records[i].Requiredness = "conditional"
		}
	}
	t.Cleanup(setRequirementsForTesting(registry))
	recordPath := filepath.Join(root, "release-evidence", "ft71-local-event.json")
	data, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	record["status"] = "not_applicable"
	record["reason"] = "fixture has no local-event capability"
	writeRecord := func() {
		encoded, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(recordPath, append(encoded, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeRecord()
	tagRelease(t, root, true)
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "publish", "--profile", "bank"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("conditional bank publish exit=%d\n%s", code, stderr.String())
	}
	indexData := mustFixtureFile(t, filepath.Join(root, "dist", "preflight", "release-index.json"))
	var index releaseIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		t.Fatal(err)
	}
	foundReason := false
	for _, status := range index.Requirements {
		if status.Key == "bank.ft71.local_event" && status.Status == "not_applicable" && status.Reason == "fixture has no local-event capability" {
			foundReason = true
		}
	}
	if !foundReason {
		t.Fatalf("conditional reason was not bound in release index: %+v", index.Requirements)
	}
	record["reason"] = ""
	writeRecord()
	stderr.Reset()
	if code := Command([]string{"--mode", "publish", "--profile", "bank"}, "0.2.0", &stderr); code != 1 || !strings.Contains(stderr.String(), "mismatched schema, key, owner, or status") {
		t.Fatalf("missing conditional reason exit=%d\n%s", code, stderr.String())
	}
}

func priorGeneration(t *testing.T, binary, root string) []byte {
	t.Helper()
	if output, err := runBuilt(t, binary, root, os.Environ()); err != nil {
		t.Fatalf("initial verify: %v\n%s", err, output)
	}
	prior, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	return prior
}

func assertPriorGeneration(t *testing.T, root string, prior []byte) {
	t.Helper()
	after, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(prior, after) {
		t.Fatal("failure replaced prior complete evidence")
	}
}

func TestBuiltCommandRejectsProducerRecordMutationsDistinctly(t *testing.T) {
	binary := buildPreflightBinary(t)
	for _, test := range []struct {
		name, mutation, want string
	}{
		{"empty", "", "record is empty"},
		{"malformed syntax", "", "is malformed"},
		{"unknown version", `"schema_version": 1`, "mismatched schema"},
		{"duplicate key", `"key":"public.ft88.data_handling"`, "duplicate JSON key"},
		{"mismatched identity", `"source_commit": "`, "identity does not match release"},
		{"digest mismatch", `"digest": "`, "digest does not match payload"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := preflightRepo(t)
			path := filepath.Join(root, "release-evidence", "ft88-data-handling.json")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			switch test.name {
			case "empty":
				data = nil
			case "malformed syntax":
				data = []byte("{\n")
			case "unknown version":
				data = bytes.Replace(data, []byte(test.mutation), []byte(`"schema_version": 2`), 1)
			case "duplicate key":
				data = bytes.Replace(data, []byte("{\n"), []byte("{\n  "+test.mutation+",\n"), 1)
			case "mismatched identity":
				data = bytes.Replace(data, []byte(test.mutation), []byte(`"source_commit": "wrong-`), 1)
			case "digest mismatch":
				data = bytes.Replace(data, []byte(test.mutation), []byte(`"digest": "00`), 1)
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
			output, err := runBuilt(t, binary, root, os.Environ())
			if err == nil || !strings.Contains(string(output), test.want) {
				t.Fatalf("error=%v, want %q\n%s", err, test.want, output)
			}
		})
	}
}

func TestBuiltCommandMissingNodeOrNPMPreservesPriorGeneration(t *testing.T) {
	binary := buildPreflightBinary(t)
	for _, missing := range []string{"node", "npm"} {
		t.Run(missing, func(t *testing.T) {
			root := preflightRepo(t)
			prior := priorGeneration(t, binary, root)
			pathDir := t.TempDir()
			for _, name := range []string{"git", "go", "node", "npm"} {
				if name == missing {
					continue
				}
				resolved, err := exec.LookPath(name)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(resolved, filepath.Join(pathDir, name)); err != nil {
					t.Fatal(err)
				}
			}
			env := append(os.Environ(), "PATH="+pathDir)
			started := time.Now()
			output, err := runBuilt(t, binary, root, env)
			if err == nil || time.Since(started) > 5*time.Second || !strings.Contains(string(output), "toolchain "+missing+" version is unavailable") {
				t.Fatalf("error=%v elapsed=%v\n%s", err, time.Since(started), output)
			}
			assertPriorGeneration(t, root, prior)
		})
	}
}

func TestBuiltCommandControlByteArchivePathPreservesPriorGeneration(t *testing.T) {
	binary := buildPreflightBinary(t)
	root := preflightRepo(t)
	prior := priorGeneration(t, binary, root)
	archive := filepath.Join(root, "dist", "artifacts", "redbench-0.2.0.tgz")
	rewriteFixtureTarball(t, archive, func(files map[string]fixtureArchiveFile) {
		files["hostile\x1bpath"] = fixtureArchiveFile{data: []byte("hostile\n"), mode: 0o644}
	})
	output, err := runBuilt(t, binary, root, os.Environ())
	if err == nil || !strings.Contains(string(output), "unsafe archive path") {
		t.Fatalf("error=%v\n%s", err, output)
	}
	assertPriorGeneration(t, root, prior)
}

func TestBuiltCommandArchiveBudgetsFailPromptlyAndPreservePriorGeneration(t *testing.T) {
	binary := buildPreflightBinary(t)
	for _, test := range []struct {
		name  string
		want  string
		write func(*testing.T, string)
	}{
		{name: "compressed bytes", want: "compressed size exceeds inspection limit", write: func(t *testing.T, path string) {
			if err := os.Truncate(path, (128<<20)+1); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "member count", want: "member count exceeds inspection limit", write: func(t *testing.T, path string) {
			if err := os.WriteFile(path, hostileArchive(t, 10_001, 1), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "expanded bytes", want: "expanded size exceeds inspection limit", write: func(t *testing.T, path string) {
			if err := os.WriteFile(path, hostileArchive(t, 65, 1<<20), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := preflightRepo(t)
			prior := priorGeneration(t, binary, root)
			test.write(t, filepath.Join(root, "dist", "artifacts", "redbench-0.2.0.tgz"))
			started := time.Now()
			output, err := runBuilt(t, binary, root, os.Environ())
			if err == nil || time.Since(started) > 5*time.Second || !strings.Contains(string(output), test.want) {
				t.Fatalf("error=%v elapsed=%v, want %q\n%s", err, time.Since(started), test.want, output)
			}
			assertPriorGeneration(t, root, prior)
		})
	}
}

func TestIndexEncodingFailurePreservesPriorGeneration(t *testing.T) {
	root := preflightRepo(t)
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stderr bytes.Buffer
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 0 {
		t.Fatalf("initial verify exit=%d\n%s", code, stderr.String())
	}
	prior, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(setIndexEncoderForTesting(func(any) ([]byte, error) {
		return nil, errors.New("deterministic encoding fault")
	}))
	stderr.Reset()
	if code := Command([]string{"--mode", "verify"}, "0.2.0", &stderr); code != 1 || !strings.Contains(stderr.String(), "release index encoding failed: deterministic encoding fault") {
		t.Fatalf("encoding failure exit=%d\n%s", code, stderr.String())
	}
	after, err := snapshotTree(filepath.Join(root, "dist", "preflight"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, prior) {
		t.Fatal("encoding failure replaced the prior complete generation")
	}
}

func TestBuiltCommandSourceIdentityDriftPreservesPriorGeneration(t *testing.T) {
	binary := buildPreflightBinary(t)
	root := preflightRepo(t)
	prior := priorGeneration(t, binary, root)
	ready := filepath.Join(t.TempDir(), "evidence-ready")
	env := append(os.Environ(), "BENCH_PREFLIGHT_EVIDENCE_READY_FILE="+ready)
	cmd := exec.Command(binary, "release-preflight", "--mode", "verify")
	cmd.Dir, cmd.Env = root, env
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			t.Fatal("evidence probe did not become ready")
		}
		time.Sleep(10 * time.Millisecond)
	}
	tree := strings.TrimSpace(gitFixtureOutput(t, root, "rev-parse", "HEAD^{tree}"))
	commit := exec.Command("git", "commit-tree", tree, "-m", "byte-identical source drift")
	commit.Dir = root
	newCommit, err := commit.Output()
	if err != nil {
		t.Fatal(err)
	}
	update := exec.Command("git", "reset", "--soft", strings.TrimSpace(string(newCommit)))
	update.Dir = root
	if output, err := update.CombinedOutput(); err != nil {
		t.Fatalf("advance HEAD: %v %s", err, output)
	}
	if err := os.Remove(ready); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err == nil || !strings.Contains(output.String(), "source identity drift") {
		t.Fatalf("error=%v\n%s", err, output.String())
	}
	assertPriorGeneration(t, root, prior)
}

func TestRegistryToolchainOperationsAreExecutionArguments(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(projectRoot(t), "internal", "releaseevidence", "requirements.json"))
	if err != nil {
		t.Fatal(err)
	}
	var registry struct {
		Toolchains []json.RawMessage `json:"toolchains"`
	}
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(registry.Toolchains)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"build":["-trimpath"`, `"pack":["--ignore-scripts","--silent"]`, `"install":["--ignore-scripts","--omit=optional"]`} {
		if !bytes.Contains(encoded, []byte(want)) {
			t.Fatalf("toolchain operations omit %s: %s", want, encoded)
		}
	}
}
