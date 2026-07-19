package preflight

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/testrepo"
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
	aggregator := filepath.Join(root, "scripts", "aggregate-native-proofs.sh")
	externalAggregator, cleanupAggregator, err := testrepo.TwoHopRelativeSymlink(aggregator)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupAggregator()
	runWith := func(script string) (string, error) {
		command := exec.Command("bash", script, proofs)
		command.Dir = root
		output, err := command.CombinedOutput()
		return string(output), err
	}
	run := func() (string, error) { return runWith(aggregator) }
	if output, err := run(); err != nil || !strings.Contains(output, "canonical targets verified") {
		t.Fatalf("valid four-target aggregation failed: %v\n%s", err, output)
	}
	if output, err := runWith(externalAggregator); err != nil || !strings.Contains(output, "canonical targets verified") {
		t.Fatalf("external-symlink aggregation failed: %v\n%s", err, output)
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
	for label, mutate := range map[string]func(map[string]any, canonicalTarget){
		"wrong runner":  func(proof map[string]any, _ canonicalTarget) { proof["runner"] = "wrong-runner" },
		"tool failure":  func(proof map[string]any, _ canonicalTarget) { proof["tools_status"] = "red" },
		"no-op failure": func(proof map[string]any, _ canonicalTarget) { proof["operations_status"] = "red" },
		"digest binding": func(proof map[string]any, target canonicalTarget) {
			if target.os == "darwin" && target.arch == "arm64" {
				proof["binary_sha256"] = strings.Repeat("2", 64)
			}
		},
		"Linux musl": func(proof map[string]any, target canonicalTarget) {
			if target.os == "linux" && target.arch == "x64" {
				proof["musl_status"] = "red"
			}
		},
		"strip status": func(proof map[string]any, target canonicalTarget) {
			if target.os == "darwin" && target.arch == "arm64" {
				proof["strip_status"] = "red"
			}
		},
	} {
		t.Run(label, func(t *testing.T) {
			writeNativeProofs(t, root, proofs, mutate)
			if output, err := run(); err == nil {
				if label == "digest binding" {
					t.Fatalf("digest binding mutation passed native proof aggregation\n%s", output)
				}
				t.Fatalf("%s did not fail closed: %v\n%s", label, err, output)
			} else if !strings.Contains(output, "native proof is incomplete or red for") {
				t.Fatalf("%s did not fail closed: %v\n%s", label, err, output)
			}
		})
	}
}

func writeNativeProofs(t *testing.T, root, directory string, mutate func(map[string]any, canonicalTarget)) {
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
			"rebuilt_sha256":    strings.Repeat("1", 64),
			"binary_sha256":     strings.Repeat("1", 64),
			"package_sha256":    strings.Repeat("1", 64),
			"archive_sha256":    strings.Repeat("1", 64),
			"musl_status":       musl,
			"operations_status": "green",
			"strip_status":      "green",
			"tools_status":      "green",
		}
		if mutate != nil {
			mutate(proof, target)
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
