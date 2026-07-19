package conformance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/testrepo"
)

func verifyAuthoritativeNativeProofMutations(root, proofPath string, env []string) []string {
	originalProof, err := os.ReadFile(proofPath)
	if err != nil {
		return []string{"release evidence probe could not read authoritative native proof"}
	}
	for _, mutation := range []struct {
		name, field string
		value       any
	}{
		{name: "digest", field: "binary_sha256", value: strings.Repeat("0", 64)},
		{name: "Linux musl", field: "musl_status", value: "red"},
		{name: "strip", field: "strip_status", value: "red"},
	} {
		var proof map[string]any
		if json.Unmarshal(originalProof, &proof) != nil {
			return []string{"release evidence probe could not decode authoritative native proof"}
		}
		proof[mutation.field] = mutation.value
		mutated, _ := json.Marshal(proof)
		if os.WriteFile(proofPath, append(mutated, '\n'), 0o644) != nil {
			return []string{"release evidence probe could not mutate authoritative native proof"}
		}
		mutationCommand := exec.Command("bash", filepath.Join(root, "scripts", "release-preflight.sh"), "--mode", "verify")
		mutationCommand.Dir, mutationCommand.Env = root, env
		output, runErr := mutationCommand.CombinedOutput()
		if restoreErr := os.WriteFile(proofPath, originalProof, 0o644); restoreErr != nil {
			return []string{"release evidence probe could not restore authoritative native proof"}
		}
		if runErr == nil {
			return []string{"authoritative native proof " + mutation.name + " mutation passed"}
		}
		if !strings.Contains(string(output), "does not match inspected artifacts") {
			return []string{"authoritative native proof " + mutation.name + " mutation was not attributed"}
		}
	}
	return nil
}

func validateReleaseIndexComponentDigests(indexData []byte) string {
	var index struct {
		Artifacts []struct {
			ComponentDigest string `json:"component_manifest_sha256"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(indexData, &index); err != nil {
		return "release evidence probe generated malformed release-index.json"
	}
	for _, artifact := range index.Artifacts {
		if artifact.ComponentDigest == "" {
			return "release index does not bind component manifest digests"
		}
	}
	return ""
}

func materializeAuthenticatedReleaseProbe(source string) (string, func(), error) {
	root, err := os.MkdirTemp("", "bench-release-probe-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() { _ = os.RemoveAll(root) }
	if err := testrepo.CommitWorkingTree(source, root); err != nil {
		cleanup()
		return "", func() {}, err
	}
	return root, cleanup, nil
}
