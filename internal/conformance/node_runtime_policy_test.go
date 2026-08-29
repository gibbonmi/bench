package conformance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/conformance/registry"
)

var (
	nodeRuntimeVersionRE   = regexp.MustCompile(`^([0-9]+)\.[0-9]+\.[0-9]+$`)
	setupNodeVersionFileRE = regexp.MustCompile(`(?m)^[ \t]+node-version-file:[ \t]*['"]\.node-version['"][ \t]*(?:#.*)?$`)
)

func checkNodeRuntimePolicy(root string) []string {
	publication := filepath.Join(root, "internal", "publication", "npm_registry.go")
	packagePath := filepath.Join(root, "package.json")
	if !exists(publication) || !exists(packagePath) {
		return nil
	}

	versionData, err := os.ReadFile(filepath.Join(root, ".node-version"))
	if err != nil {
		return []string{".node-version is missing or unreadable"}
	}
	version := strings.TrimSpace(string(versionData))
	match := nodeRuntimeVersionRE.FindStringSubmatch(version)
	if match == nil {
		return []string{".node-version must contain a full semantic version"}
	}
	major := match[1]

	var diags []string
	var pkg struct {
		Engines struct {
			Node string `json:"node"`
		} `json:"engines"`
	}
	packageData, packageErr := os.ReadFile(packagePath)
	if packageErr != nil || json.Unmarshal(packageData, &pkg) != nil || pkg.Engines.Node != ">="+major {
		diags = append(diags, "package.json engines.node does not project .node-version")
	}

	for _, rel := range []string{".github/workflows/release.yml", ".github/workflows/native-runtime.yml"} {
		diags = append(diags, checkSetupNodeVersionFile(rel, readIfExists(filepath.Join(root, filepath.FromSlash(rel))))...)
	}
	if !strings.Contains(readIfExists(publication), `minStagedNodeVersion = "`+version+`"`) {
		diags = append(diags, "staged publication Node floor does not project .node-version")
	}
	if !strings.Contains(readIfExists(filepath.Join(root, "docs", "release-runbook.md")), "Node "+major+"+") {
		diags = append(diags, "release guidance Node floor does not project .node-version")
	}
	return diags
}

func checkSetupNodeVersionFile(rel, workflow string) []string {
	if workflow == "" {
		return nil
	}
	parts := strings.Split(workflow, "uses: actions/setup-node@")
	if len(parts) == 1 {
		return []string{rel + " has no setup-node step"}
	}
	var diags []string
	for index, rest := range parts[1:] {
		if end := strings.Index(rest, "\n      - "); end >= 0 {
			rest = rest[:end]
		}
		if !setupNodeVersionFileRE.MatchString(rest) {
			diags = append(diags, rel+" setup-node step "+strconv.Itoa(index+1)+" does not read .node-version")
		}
	}
	return diags
}

func TestNodeRuntimePolicyProjectionsBite(t *testing.T) {
	root := t.TempDir()
	originals := map[string]string{
		".node-version":                        "24.0.0\n",
		"package.json":                         "{\"engines\":{\"node\":\">=24\"}}\n",
		"internal/publication/npm_registry.go": "package publication\nconst minStagedNodeVersion = \"24.0.0\"\n",
		"docs/release-runbook.md":              "The operator uses Node 24+.\n",
		".github/workflows/release.yml":        "jobs:\n  release:\n    steps:\n      - uses: actions/setup-node@digest\n        with:\n          node-version-file: '.node-version'\n",
		".github/workflows/native-runtime.yml": "jobs:\n  verify:\n    steps:\n      - uses: actions/setup-node@digest\n        with:\n          node-version-file: '.node-version'\n",
	}
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for rel, body := range originals {
		write(rel, body)
	}

	for _, mutation := range []struct {
		name, rel, body, want string
	}{
		{"invalid canonical version", ".node-version", "24\n", ".node-version must contain a full semantic version"},
		{"stale package engine", "package.json", "{\"engines\":{\"node\":\">=23\"}}\n", "package.json engines.node does not project .node-version"},
		{"stale publication floor", "internal/publication/npm_registry.go", "package publication\nconst minStagedNodeVersion = \"23.0.0\"\n", "staged publication Node floor does not project .node-version"},
		{"stale release guidance", "docs/release-runbook.md", "The operator uses Node 23+.\n", "release guidance Node floor does not project .node-version"},
		{"literal release workflow version", ".github/workflows/release.yml", "jobs:\n  release:\n    steps:\n      - uses: actions/setup-node@digest\n        with:\n          node-version: '24'\n", ".github/workflows/release.yml setup-node step 1 does not read .node-version"},
		{"literal native workflow version", ".github/workflows/native-runtime.yml", "jobs:\n  verify:\n    steps:\n      - uses: actions/setup-node@digest\n        with:\n          node-version: '24'\n", ".github/workflows/native-runtime.yml setup-node step 1 does not read .node-version"},
		{"commented canonical workflow version", ".github/workflows/release.yml", "jobs:\n  release:\n    steps:\n      - uses: actions/setup-node@digest\n        with:\n          node-version: '23'\n          # node-version-file: '.node-version'\n", ".github/workflows/release.yml setup-node step 1 does not read .node-version"},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			write(mutation.rel, mutation.body)
			if diags := checkPackageCoreAndGuards(root, registry.Dev); !containsDiagnostic(diags, mutation.want) {
				t.Fatalf("mutation did not bite with %q:\n%s", mutation.want, strings.Join(diags, "\n"))
			}
			write(mutation.rel, originals[mutation.rel])
			if diags := checkPackageCoreAndGuards(root, registry.Dev); containsDiagnostic(diags, mutation.want) {
				t.Fatalf("restored projection retained %q:\n%s", mutation.want, strings.Join(diags, "\n"))
			}
		})
	}
}
