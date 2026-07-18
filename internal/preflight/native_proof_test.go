package preflight

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNativeProofAggregationRejectsIncompleteAndRedProofs(t *testing.T) {
	root := preflightRepo(t)
	proofs := filepath.Join(root, "dist", "native-proofs")
	writeNativeProofs(t, proofs, nil)
	run := func() (string, error) {
		command := exec.Command("bash", filepath.Join(root, "scripts", "aggregate-native-proofs.sh"), proofs)
		command.Dir = root
		output, err := command.CombinedOutput()
		return string(output), err
	}
	if output, err := run(); err != nil || !strings.Contains(output, "four canonical targets verified") {
		t.Fatalf("valid four-target aggregation failed: %v\n%s", err, output)
	}
	if err := os.Remove(filepath.Join(proofs, "linux-x64.json")); err != nil {
		t.Fatal(err)
	}
	if output, err := run(); err == nil || !strings.Contains(output, "native proof set is missing linux/x64") {
		t.Fatalf("missing proof did not fail closed: %v\n%s", err, output)
	}
	writeNativeProofs(t, proofs, nil)
	if err := os.WriteFile(filepath.Join(proofs, "darwin-arm64.json"), []byte("not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if output, err := run(); err == nil || !strings.Contains(output, "native proof is incomplete or red for darwin/arm64") {
		t.Fatalf("corrupt proof did not fail closed: %v\n%s", err, output)
	}
	for label, mutate := range map[string]func(map[string]any){
		"wrong runner":  func(proof map[string]any) { proof["runner"] = "wrong-runner" },
		"tool failure":  func(proof map[string]any) { proof["tools_status"] = "red" },
		"no-op failure": func(proof map[string]any) { proof["operations_status"] = "red" },
	} {
		t.Run(label, func(t *testing.T) {
			writeNativeProofs(t, proofs, mutate)
			if output, err := run(); err == nil || !strings.Contains(output, "native proof is incomplete or red for darwin/arm64") {
				t.Fatalf("%s did not fail closed: %v\n%s", label, err, output)
			}
		})
	}
}

func writeNativeProofs(t *testing.T, directory string, mutate func(map[string]any)) {
	t.Helper()
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, target := range []struct {
		name, runner, musl string
	}{
		{"darwin-arm64", "macos-15", "not_applicable"},
		{"darwin-x64", "macos-15-intel", "not_applicable"},
		{"linux-arm64", "ubuntu-24.04-arm", "green"},
		{"linux-x64", "ubuntu-24.04", "green"},
	} {
		proof := map[string]any{
			"schema_version":    1,
			"target":            target.name,
			"runner":            target.runner,
			"status":            "green",
			"rebuilt_sha256":    "rebuilt",
			"binary_sha256":     "binary",
			"package_sha256":    "package",
			"archive_sha256":    "archive",
			"musl_status":       target.musl,
			"operations_status": "green",
			"strip_status":      "green",
			"tools_status":      "green",
		}
		if target.name == "darwin-arm64" && mutate != nil {
			mutate(proof)
		}
		data, err := json.Marshal(proof)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, target.name+".json"), append(data, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
