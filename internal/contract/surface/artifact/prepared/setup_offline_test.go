package prepared

import fixture "github.com/gibbonmi/bench/internal/contract/surface/artifact/internal/fixture"

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestPackedArtifactRunsSetupOfflineFromASpacedPrefix is FT76 story 12 / coverage row
// 12 (packed-artifact cold run, Seam 2): stage the real wrapper+native tarballs the
// same way assertInstalledArtifactLifecycle does, install them offline into an npm
// prefix whose path contains a space, then run `bench setup --yes` through the
// installed wrapper against an empty target git repo. It asserts the full convergence
// set and, critically, that every durable asset lands repo-local — nothing durable is
// left behind in the npm-managed prefix, which is ephemeral cache state.
func TestPackedArtifactRunsSetupOfflineFromASpacedPrefix(t *testing.T) {
	root := contract.SubjectRoot(t)
	contract.SkipIfSubjectFileMissing(t, "scripts/build-artifacts.sh")
	shared := requireSharedArtifactSet(t)
	out := shared.outputDir

	var wrapper struct {
		Version string `json:"version"`
	}
	contract.ReadJSONFile(t, filepath.Join(root, "package.json"), &wrapper)
	target := runtime.GOOS + "-" + runtime.GOARCH
	if runtime.GOARCH == "amd64" {
		target = runtime.GOOS + "-x64"
	}
	wrapperTar := filepath.Join(out, "redbench-"+wrapper.Version+".tgz")
	nativeTar := filepath.Join(out, "redbench-"+target+"-"+wrapper.Version+".tgz")

	tmp := t.TempDir()
	// The spec's row-12 edge is a prefix path containing a space: npm and the wrapper
	// must both survive it, not just tolerate plain paths.
	prefix := filepath.Join(tmp, "npm prefix [with a space]")
	home := filepath.Join(tmp, "home [with a space]")
	repo := filepath.Join(tmp, "target repo")
	for _, dir := range []string{prefix, home, repo} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(prefix, "package.json"), []byte("{\"private\":true}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Same offline-install recipe assertInstalledArtifactLifecycle already uses for
	// link/init - reused here rather than re-derived, just against a spaced prefix.
	npmEnv := map[string]string{
		"npm_config_audit":           "false",
		"npm_config_cache":           filepath.Join(tmp, "npm-cache"),
		"npm_config_fund":            "false",
		"npm_config_offline":         "true",
		"npm_config_registry":        "http://127.0.0.1:9",
		"npm_config_update_notifier": "false",
	}
	fixture.RunLifecycle(t, prefix, npmEnv, "npm", "install", "--ignore-scripts", "--omit=optional", "--prefix", prefix, wrapperTar, nativeTar)
	installed := filepath.Join(prefix, "node_modules", "redbench", "bin", "bench.sh")
	if info, err := os.Stat(installed); err != nil || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("offline install into spaced prefix did not produce an executable wrapper: %v, %v", info, err)
	}

	fixture.RunLifecycle(t, repo, nil, "git", "init", "-q")
	fixture.RunLifecycle(t, repo, nil, "git", "config", "user.email", "bench@local")
	fixture.RunLifecycle(t, repo, nil, "git", "config", "user.name", "bench")
	// A single unambiguous gate signal keeps this leg's setup run fully green (exit 0)
	// so the assertions below stay about convergence and locality, not the separate
	// zero-signal/ambiguity paths rows 3 and 9 already cover.
	if err := os.WriteFile(filepath.Join(repo, "go.mod"), []byte("module target\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fixture.RunLifecycle(t, repo, nil, "git", "add", "go.mod")
	fixture.RunLifecycle(t, repo, nil, "git", "commit", "-qm", "seed")

	env := map[string]string{
		"BENCH_HOME": home,
		"HOME":       home,
		"PATH":       "/usr/bin:/bin",
	}
	setupOut := fixture.RunLifecycle(t, repo, env, "bash", installed, "setup", "--yes")
	assertPackedSetupConverged(t, repo, prefix, setupOut)
}

// assertPackedSetupConverged is the row-12 oracle: every asset the spec names by name
// (the AGENTS.md marker pair, the CLAUDE.md import lines, gate.sh, the profile, the
// link manifest, and the repo-local launcher) plus the locality claim - the npm
// prefix's copy of the wrapper package is untouched by setup's writes, because it is
// ephemeral install cache, not where durable state belongs.
func assertPackedSetupConverged(t *testing.T, repo, prefix, setupOut string) {
	t.Helper()
	agents, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil || !strings.Contains(string(agents), "<!-- bench:start -->") || !strings.Contains(string(agents), "<!-- bench:end -->") {
		t.Fatalf("packed setup did not converge the AGENTS.md marker pair: %q, %v", agents, err)
	}
	claude, err := os.ReadFile(filepath.Join(repo, "CLAUDE.md"))
	if err != nil || !strings.Contains(string(claude), "@AGENTS.md") || !strings.Contains(string(claude), "@.bench/BENCH.md") {
		t.Fatalf("packed setup did not converge the CLAUDE.md import lines: %q, %v", claude, err)
	}
	gate, err := os.ReadFile(filepath.Join(repo, ".bench", "gate.sh"))
	if err != nil || !strings.Contains(string(gate), "go test ./...") {
		t.Fatalf("packed setup did not converge the inferred gate.sh: %q, %v", gate, err)
	}
	profile := filepath.Join(repo, "projects", filepath.Base(repo)+".md")
	if _, err := os.Stat(profile); err != nil {
		t.Fatalf("packed setup did not scaffold the profile: %v", err)
	}
	manifest, err := os.ReadFile(filepath.Join(repo, ".bench", "link-manifest.tsv"))
	if err != nil || !strings.Contains(string(manifest), "#kit\t") {
		t.Fatalf("packed setup did not write the link manifest: %q, %v", manifest, err)
	}
	launcher := filepath.Join(repo, ".bench", "bin", "bench.sh")
	if info, err := os.Stat(launcher); err != nil || info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("packed setup did not install an executable repo-local launcher: %v, %v", info, err)
	}
	if !strings.Contains(setupOut, "converged "+repo) {
		t.Fatalf("packed setup did not report convergence: %q", setupOut)
	}

	// Locality: the npm prefix's own copy of the wrapper package must be exactly what
	// npm installed - setup's writes belong to the target repo, never back into the
	// install-cache copy that a subsequent `npm prune`/reinstall can discard freely.
	installed := filepath.Join(prefix, "node_modules", "redbench")
	for _, durable := range []string{
		filepath.Join(installed, "AGENTS.md"), filepath.Join(installed, "CLAUDE.md"),
		filepath.Join(installed, ".bench", "gate.sh"), filepath.Join(installed, ".bench", "link-manifest.tsv"),
		filepath.Join(installed, ".bench", "bin", "bench.sh"),
		filepath.Join(installed, "projects", filepath.Base(repo)+".md"),
	} {
		if _, err := os.Stat(durable); !os.IsNotExist(err) {
			t.Fatalf("packed setup left durable state in the ephemeral npm-cache copy: %s (err=%v)", durable, err)
		}
	}
}
