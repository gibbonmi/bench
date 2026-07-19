package preflight

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type canonicalTarget struct {
	os, arch, goos, goarch, runner string
}

func canonicalReleaseTargets(t *testing.T, root string) []canonicalTarget {
	t.Helper()
	plan := exec.Command("node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "targets")
	data, err := plan.Output()
	if err != nil {
		t.Fatalf("read canonical release targets: %v", err)
	}
	var targets []canonicalTarget
	for _, row := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fields := strings.Split(row, "\t")
		if len(fields) != 5 {
			t.Fatalf("malformed canonical release target: %q", row)
		}
		targets = append(targets, canonicalTarget{fields[0], fields[1], fields[2], fields[3], fields[4]})
	}
	return targets
}

func TestNativeProofAggregationRejectsIncompleteAndRedProofs(t *testing.T) {
	root := preflightRepo(t)
	proofs := filepath.Join(root, "dist", "native-proofs")
	writeNativeProofs(t, root, proofs, nil)
	run := func() (string, error) {
		command := exec.Command("bash", filepath.Join(root, "scripts", "aggregate-native-proofs.sh"), proofs)
		command.Dir = root
		output, err := command.CombinedOutput()
		return string(output), err
	}
	if output, err := run(); err != nil || !strings.Contains(output, "canonical targets verified") {
		t.Fatalf("valid four-target aggregation failed: %v\n%s", err, output)
	}
	if err := os.Remove(filepath.Join(proofs, "linux-x64.json")); err != nil {
		t.Fatal(err)
	}
	if output, err := run(); err == nil || !strings.Contains(output, "native proof set is missing linux/x64") {
		t.Fatalf("missing proof did not fail closed: %v\n%s", err, output)
	}
	writeNativeProofs(t, root, proofs, nil)
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
			writeNativeProofs(t, root, proofs, mutate)
			if output, err := run(); err == nil || !strings.Contains(output, "native proof is incomplete or red for darwin/arm64") {
				t.Fatalf("%s did not fail closed: %v\n%s", label, err, output)
			}
		})
	}
}

func writeNativeProofs(t *testing.T, root, directory string, mutate func(map[string]any)) {
	t.Helper()
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, target := range canonicalReleaseTargets(t, root) {
		name, runner, musl := target.os+"-"+target.arch, target.runner, "not_applicable"
		if target.os == "linux" {
			musl = "green"
		}
		proof := map[string]any{
			"schema_version":    1,
			"target":            name,
			"runner":            runner,
			"status":            "green",
			"rebuilt_sha256":    "rebuilt",
			"binary_sha256":     "binary",
			"package_sha256":    "package",
			"archive_sha256":    "archive",
			"musl_status":       musl,
			"operations_status": "green",
			"strip_status":      "green",
			"tools_status":      "green",
		}
		if name == "darwin-arm64" && mutate != nil {
			mutate(proof)
		}
		data, err := json.Marshal(proof)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, name+".json"), append(data, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
