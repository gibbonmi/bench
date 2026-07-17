package preflight

import (
	"bytes"
	"encoding/json"
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
	cmd := exec.Command(binary, "release-preflight", "--mode", "verify")
	cmd.Dir = root
	cmd.Env = env
	return cmd.CombinedOutput()
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
