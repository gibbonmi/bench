package conformance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type workflowTriggerShape struct {
	pullRequest, pushBranches, mainBranch bool
}

func checkReleasePreflight(root string) []string {
	var diags []string
	data, err := os.ReadFile(filepath.Join(root, "internal", "preflight", "registry.json"))
	if err != nil {
		return []string{"release preflight registry is missing or unreadable"}
	}
	var registry struct {
		Verify      []string `json:"verify"`
		PublishOnly []string `json:"publish_only"`
	}
	if err == nil && json.Unmarshal(data, &registry) != nil {
		return []string{"release preflight registry is missing or unreadable"}
	}
	requirementData, err := os.ReadFile(filepath.Join(root, "internal", "preflight", "requirements.json"))
	if err != nil {
		diags = append(diags, "release requirement registry is missing or unreadable")
	}
	var requirementRegistry struct {
		SchemaVersion int `json:"schema_version"`
		Records       []struct {
			Key          string   `json:"key"`
			Owner        string   `json:"owner"`
			Schema       string   `json:"schema"`
			Profiles     []string `json:"profiles"`
			Requiredness string   `json:"requiredness"`
			Path         string   `json:"path"`
		} `json:"records"`
	}
	if err == nil && (json.Unmarshal(requirementData, &requirementRegistry) != nil || requirementRegistry.SchemaVersion != 1) {
		diags = append(diags, "release requirement registry is missing or unreadable")
	}
	if err != nil {
		requirementData = nil
	}
	public, bank := map[string]bool{}, map[string]bool{}
	for _, record := range requirementRegistry.Records {
		if record.Key == "" || record.Owner == "" || record.Schema == "" || record.Path == "" || record.Requiredness != "required" && record.Requiredness != "conditional" || len(record.Profiles) == 0 {
			return append(diags, "release requirement registry has incomplete schema")
		}
		for _, profile := range record.Profiles {
			if profile == "public" {
				public[record.Key] = true
			}
			if profile == "bank" {
				bank[record.Key] = true
			}
		}
	}
	for key := range public {
		if !bank[key] {
			diags = append(diags, "release requirement registry bank profile is not a strict public superset")
			break
		}
	}
	for _, key := range []string{"public.ft88.data_handling", "public.ft87.offline_network_control", "bank.ft71.local_event"} {
		if !containsKey(requirementRegistry.Records, key) {
			diags = append(diags, "release requirement registry omits "+key)
		}
	}
	assets := readIfExists(filepath.Join(root, "scripts", "wrapper-assets.json"))
	for _, asset := range []string{"LICENSE", "governance"} {
		if !strings.Contains(assets, `"source": "`+asset+`"`) {
			diags = append(diags, "release package evidence allowlist omits "+asset)
		}
	}
	releaseEvidence := readIfExists(filepath.Join(root, "internal", "preflight", "release_evidence.go"))
	if !strings.Contains(releaseEvidence, `component_manifest_sha256`) {
		diags = append(diags, "release index does not bind component manifest digests")
	}
	if !reflect.DeepEqual(registry.Verify, []string{"gate", "race", "vet", "vulnerability", "artifacts", "smoke"}) {
		diags = append(diags, "release preflight verify registry omits or reorders a required phase class")
	}
	if !reflect.DeepEqual(registry.PublishOnly, []string{"identity", "ancestry", "changelog"}) {
		diags = append(diags, "release preflight publish registry omits or reorders a required phase class")
	}
	native := readIfExists(filepath.Join(root, ".github", "workflows", "native-runtime.yml"))
	release := readIfExists(filepath.Join(root, ".github", "workflows", "release.yml"))
	installer := readIfExists(filepath.Join(root, "scripts", "install-govulncheck.sh"))
	if strings.Count(native, "bash scripts/install-govulncheck.sh") != 1 || strings.Count(release, "bash scripts/install-govulncheck.sh") != 1 || !regexp.MustCompile(`govulncheck@v[0-9]+\.[0-9]+\.[0-9]+`).MatchString(installer) {
		diags = append(diags, "release workflows do not consume the repository-owned govulncheck setup")
	}
	if strings.Contains(native, "govulncheck@") || strings.Contains(release, "govulncheck@") {
		diags = append(diags, "release workflows duplicate the govulncheck version pin")
	}
	for message, anchor := range map[string]string{"native verification bypasses full release preflight": "scripts/release-preflight.sh --mode verify\n", "native verification does not upload preflight evidence": "verify-preflight-evidence", "native runner matrix bypasses focused smoke": "--mode verify --phase smoke"} {
		if strings.Count(native, anchor) != 1 {
			diags = append(diags, message)
		}
	}
	for message, anchor := range map[string]string{"tag publication bypasses full release preflight": "scripts/release-preflight.sh --mode publish --profile public\n", "tag publication does not upload preflight evidence": "publish-preflight-evidence", "tag runner matrix bypasses focused smoke": "--mode publish --profile public --phase smoke", "tag evidence does not request repository maximum retention": "retention-days: ${{ github.retention_days }}", "publication does not wait for preflight and every native smoke row": "needs: [preflight, smoke]"} {
		if !strings.Contains(release, anchor) {
			diags = append(diags, message)
		}
	}
	platform, wrapper := strings.Index(release, "name: Publish platform packages"), strings.Index(release, "name: Publish wrapper")
	if platform < 0 || wrapper < platform {
		diags = append(diags, "release publication is not platform-first and wrapper-last")
	}
	if !regexp.MustCompile(`(?m)^toolchain go[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(readIfExists(filepath.Join(root, "go.mod"))) {
		diags = append(diags, "release preflight requires an exact Go patch toolchain")
	}
	return diags
}

func containsKey(records []struct {
	Key          string   `json:"key"`
	Owner        string   `json:"owner"`
	Schema       string   `json:"schema"`
	Profiles     []string `json:"profiles"`
	Requiredness string   `json:"requiredness"`
	Path         string   `json:"path"`
}, want string) bool {
	for _, record := range records {
		if record.Key == want {
			return true
		}
	}
	return false
}

func nativeWorkflowTriggers(text string) workflowTriggerShape {
	var shape workflowTriggerShape
	inOn, inPush, inBranches := false, false, false
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimRight(raw, " \t\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inOn = line == "on:"
			inPush, inBranches = false, false
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
