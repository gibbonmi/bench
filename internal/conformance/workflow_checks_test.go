package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type workflowTriggerShape struct {
	pullRequest, pushBranches, mainBranch bool
}

// The workflow seam is the job body: these checks parse job ownership before
// asserting that evidence follows all native proof rows.
func checkReleaseWorkflow(root string) []string {
	if !exists(filepath.Join(root, "scripts", "release-plan.json")) {
		return nil
	}
	wf := filepath.Join(root, ".github", "workflows", "release.yml")
	if !exists(wf) {
		return []string{"release workflow missing (.github/workflows/release.yml)"}
	}
	text := readIfExists(wf)
	var diags []string
	if !regexp.MustCompile(`(?m)^\s*tags:`).MatchString(text) {
		diags = append(diags, "release workflow does not trigger on tags")
	}
	if !strings.Contains(readIfExists(filepath.Join(root, "scripts", "build-artifacts.sh")), "scripts/release-plan.mjs") {
		diags = append(diags, "artifact builder does not derive targets from the canonical release plan")
	}
	for message, anchor := range map[string]string{"release workflow does not run full publish preflight": "scripts/release-preflight.sh --mode publish", "release workflow does not publish to npm": "npm publish", "release workflow does not publish with provenance": "provenance"} {
		if !strings.Contains(text, anchor) {
			diags = append(diags, message)
		}
	}
	return diags
}

func checkNativeRuntimeWorkflow(root string) []string {
	workflow := filepath.Join(root, ".github", "workflows", "native-runtime.yml")
	if !exists(workflow) {
		return []string{"native verification workflow missing (.github/workflows/native-runtime.yml)"}
	}
	text, diags := readIfExists(workflow), []string{}
	triggers := nativeWorkflowTriggers(text)
	for label, present := range map[string]bool{"pull requests": triggers.pullRequest, "default branch pushes": triggers.pushBranches, "main branch": triggers.mainBranch} {
		if !present {
			diags = append(diags, "native verification workflow does not include "+label)
		}
	}
	for label, anchor := range map[string]string{"canonical matrix": "scripts/release-plan.mjs", "derived matrix": "fromJSON(needs.preflight.outputs.matrix)", "full preflight": "scripts/release-preflight.sh --mode verify", "native proof builder": "scripts/native-proof.sh", "native proof aggregate": "scripts/aggregate-native-proofs.sh", "native proof evidence": "dist/native-proofs"} {
		if !strings.Contains(text, anchor) {
			diags = append(diags, "native verification workflow does not include "+label)
		}
	}
	if job := workflowJob(text, "evidence"); !strings.Contains(job, "needs: [preflight, native-proof]") || !strings.Contains(job, "scripts/release-preflight.sh --mode verify") {
		diags = append(diags, "native verification does not finalize evidence after every native proof")
	}
	if job := workflowJob(text, "smoke"); !strings.Contains(job, "needs: [preflight, evidence]") || !strings.Contains(job, "verify-preflight-evidence") || !strings.Contains(job, "scripts/smoke-artifacts.sh") {
		diags = append(diags, "native verification does not run smoke from finalized evidence")
	}
	if proof := readIfExists(filepath.Join(root, "scripts", "native-proof.sh")); proof != "" && !strings.Contains(proof, "docker run --rm --network none") {
		diags = append(diags, "native proof does not isolate the Linux non-glibc execution")
	}
	data, err := os.ReadFile(filepath.Join(root, "scripts", "release-plan.json"))
	var plan struct {
		Targets []struct {
			Runner string `json:"runner"`
		} `json:"targets"`
	}
	if err != nil || json.Unmarshal(data, &plan) != nil {
		return append(diags, "native verification matrix is unreadable")
	}
	want, got := []string{"macos-15", "macos-15-intel", "ubuntu-24.04-arm", "ubuntu-24.04"}, make([]string, 0, len(plan.Targets))
	for _, row := range plan.Targets {
		got = append(got, row.Runner)
	}
	if !reflect.DeepEqual(got, want) {
		diags = append(diags, fmt.Sprintf("native verification runner labels = %v, want %v", got, want))
	}
	return diags
}

func workflowJob(workflow, name string) string {
	needle := "  " + name + ":\n"
	start := strings.Index(workflow, needle)
	if start < 0 {
		return ""
	}
	rest, end := workflow[start+len(needle):], len(workflow[start+len(needle):])
	for _, candidate := range []string{"\n  preflight:\n", "\n  native-proof:\n", "\n  evidence:\n", "\n  smoke:\n", "\n  authorize:\n", "\n  publish:\n"} {
		if offset := strings.Index(rest, candidate); offset >= 0 && offset < end {
			end = offset
		}
	}
	return rest[:end]
}

func nativeWorkflowTriggers(text string) workflowTriggerShape {
	var shape workflowTriggerShape
	inOn, inPush, inBranches := false, false, false
	for _, raw := range strings.Split(text, "\n") {
		line, trimmed := strings.TrimRight(raw, " \t\r"), strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inOn, inPush, inBranches = line == "on:", false, false
			continue
		}
		if !inOn {
			continue
		}
		if indent == 2 {
			inPush, inBranches = trimmed == "push:", false
			if trimmed == "pull_request:" {
				shape.pullRequest = true
			}
			continue
		}
		if inPush && indent == 4 && trimmed == "branches:" {
			shape.pushBranches, inBranches = true, true
			continue
		}
		if inBranches && indent == 6 && trimmed == "- main" {
			shape.mainBranch = true
		}
	}
	return shape
}

func TestNativeWorkflowEvidenceEdgeBites(t *testing.T) {
	kit, err := findKitRoot()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	for _, relative := range []string{"scripts/release-plan.json", "scripts/native-proof.sh"} {
		data, err := os.ReadFile(filepath.Join(kit, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	workflow, err := os.ReadFile(filepath.Join(kit, ".github", "workflows", "native-runtime.yml"))
	if err != nil {
		t.Fatal(err)
	}
	broken := strings.Replace(string(workflow), "needs: [preflight, evidence]", "needs: preflight", 1)
	if broken == string(workflow) {
		t.Fatal("native workflow mutation did not remove the evidence edge")
	}
	path := filepath.Join(root, ".github", "workflows", "native-runtime.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	if diagnostics := strings.Join(checkNativeRuntimeWorkflow(root), "\n"); !strings.Contains(diagnostics, "native verification does not run smoke from finalized evidence") {
		t.Fatalf("removed evidence edge did not bite:\n%s", diagnostics)
	}
}
