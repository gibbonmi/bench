package conformance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gibbonmi/bench/internal/lines"
)

func checkLineRouting(root string) []string {
	var diags []string
	diags = append(diags, checkLineBinding(root)...)
	diags = append(diags, checkClaudeAgentHookWiring(root)...)
	diags = append(diags, checkAgentHookBehavior(root)...)
	diags = append(diags, checkAdapterLineGuards(root)...)
	return diags
}

func checkLineBinding(root string) []string {
	path := filepath.Join(root, ".bench", "lines.env")
	if !exists(path) {
		return []string{"lines.env missing: .bench/lines.env is the tier binding enforcement reads"}
	}
	content := []byte(readIfExists(path))
	binding := lines.ParseBinding(content)
	var diags []string
	tiers := []struct {
		label string
		key   string
		value string
	}{
		{"top", "BENCH_TIER_TOP", binding.Top},
		{"mid", "BENCH_TIER_MID", binding.Mid},
		{"cheap", "BENCH_TIER_CHEAP", binding.Cheap},
	}
	modelID := regexp.MustCompile(`^claude-[a-z0-9][a-z0-9.-]*$`)
	for _, tier := range tiers {
		if tier.value == "" {
			diags = append(diags, fmt.Sprintf("lines.env tier unset: %s has no value in .bench/lines.env", tier.key))
		} else if !modelID.MatchString(tier.value) {
			diags = append(diags, fmt.Sprintf("lines.env tier malformed: %s='%s' is not a model id", tier.label, tier.value))
		}
	}
	aliases := []struct {
		key   string
		value string
	}{
		{"BENCH_ALIAS_TOP", binding.AliasTop},
		{"BENCH_ALIAS_MID", binding.AliasMid},
		{"BENCH_ALIAS_CHEAP", binding.AliasCheap},
	}
	aliasRe := regexp.MustCompile(`^[a-z0-9-]+$`)
	for _, alias := range aliases {
		if !regexp.MustCompile(`(?m)^[ \t]*` + regexp.QuoteMeta(alias.key) + `=`).Match(content) {
			continue
		}
		if !aliasRe.MatchString(alias.value) {
			diags = append(diags, fmt.Sprintf("lines.env alias malformed: %s='%s' is not a bare alias", alias.key, alias.value))
		}
	}

	profile := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	if profile != "" {
		for _, tier := range tiers {
			if tier.value == "" {
				continue
			}
			if !strings.Contains(profile, tier.value) {
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not name bound model id '%s' (%s in lines.env)", tier.value, tier.key))
			}
		}
		for _, alias := range aliases {
			if alias.value == "" {
				continue
			}
			want := alias.key + "=" + alias.value
			if !strings.Contains(profile, want) {
				diags = append(diags, fmt.Sprintf("profile Lines prose stale: projects/benchkit.md does not carry alias declaration %s", want))
			}
		}
	}
	return diags
}

func checkClaudeAgentHookWiring(root string) []string {
	path := filepath.Join(root, ".claude", "settings.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg struct {
		Hooks map[string][]struct {
			Matcher string `json:"matcher"`
			Hooks   []struct {
				Command string `json:"command"`
			} `json:"hooks"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	for _, group := range cfg.Hooks["PreToolUse"] {
		if group.Matcher != "Agent" {
			continue
		}
		for _, hook := range group.Hooks {
			if strings.Contains(hook.Command, ".bench/hooks/check-agent-line.sh") {
				return nil
			}
		}
	}
	return []string{"claude settings.json PreToolUse Agent matcher missing or does not run .bench/hooks/check-agent-line.sh"}
}

func checkAgentHookBehavior(root string) []string {
	hook := filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")
	realBench := filepath.Join(root, "bin", "bench.sh")
	if !exists(hook) {
		return nil
	}
	if !exists(realBench) {
		probe := runWithInput(root, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "bash", hook)
		if probe != nil && probe.ExitCode == 0 {
			return []string{"check-agent-line.sh does not deny an unbound model"}
		}
		return nil
	}

	bindir, cleanup, err := wrapperStubDir(realBench)
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanup()
	env := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
	var diags []string

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\nBENCH_ALIAS_MID=opus\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupRouted()

	hookCase := func(label, cwd, input, wantErr string, wantExit int) {
		probe := runWithInputEnv(cwd, env, input, "bash", hook)
		if probe == nil || probe.ExitCode != wantExit {
			got := -1
			if probe != nil {
				got = probe.ExitCode
			}
			diags = append(diags, fmt.Sprintf("check-agent-line.sh %s exit %d (want %d)", label, got, wantExit))
			return
		}
		if wantErr != "" && !strings.Contains(probe.Stderr, wantErr) && !strings.Contains(probe.Stdout, wantErr) {
			diags = append(diags, fmt.Sprintf("check-agent-line.sh %s did not warn with %q", label, wantErr))
		}
	}
	hookCase("denies a bound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-opus-4-8"}}`, "", 0)
	hookCase("denies a declared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus"}}`, "", 0)
	hookCase("does not deny an undeclared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"sonnet"}}`, "", 2)
	hookCase("does not deny an unbound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "", 2)
	hookCase("does not fail open on malformed stdin", routed, `not json at all`, "not parseable as JSON", 0)
	hookCase("does not fail open on a missing model field", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "no resolvedModel/model field", 0)

	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	hookCase("does not fail open without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "no .bench/lines.env", 0)

	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupPartial()
	hookCase("does not fail open on an incomplete binding", partial, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"claude-nonexistent-9"}}`, "unset or empty", 0)
	return diags
}

func checkAdapterLineGuards(root string) []string {
	if !exists(filepath.Join(root, ".bench", "adapters")) {
		return nil
	}
	var diags []string
	realBench := filepath.Join(root, "bin", "bench.sh")
	bindir := ""
	cleanup := func() {}
	if exists(realBench) {
		var err error
		bindir, cleanup, err = adapterStubDir(realBench)
		if err != nil {
			return []string{"adapter line guard setup failed: " + err.Error()}
		}
		defer cleanup()
	}
	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			diags = append(diags, "adapter missing from .bench/adapters: "+name)
			continue
		}
		text := readIfExists(path)
		hasResolveCall := regexp.MustCompile(`(?m)^[[:space:]]*model="\$\([^)]*resolve-model`).MatchString(text)
		if !hasResolveCall {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
	}
	if bindir == "" {
		return diags
	}

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=claude-opus-4-8\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupRouted()
	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=claude-fable-5\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=claude-sonnet-4-6\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupPartial()

	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			continue
		}
		envBase := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
		bound := runAtEnv(routed, append(envBase, "BENCH_MODEL=claude-opus-4-8"), "bash", path, "--line probe prompt")
		if bound == nil || bound.ExitCode != 0 || !strings.Contains(bound.Stdout, "claude-opus-4-8") || !strings.Contains(bound.Stdout, "--line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass BENCH_MODEL and a dash-leading prompt to the harness in a routed repo", name))
		}
		unset := runAtEnv(routed, envBase, "bash", path, "line probe prompt")
		if unset != nil && unset.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
		}
		unbound := runAtEnv(routed, append(envBase, "BENCH_MODEL=claude-nonexistent-9"), "bash", path, "line probe prompt")
		if unbound != nil && unbound.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
		pass := runAtEnv(unrouted, envBase, "bash", path, "line probe prompt")
		if pass == nil || pass.ExitCode != 0 || !strings.Contains(pass.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass through in an unrouted repo", name))
		}
		explicit := runAtEnv(unrouted, append(envBase, "BENCH_MODEL=claude-anything-7"), "bash", path, "line probe prompt")
		if explicit == nil || explicit.ExitCode != 0 || !strings.Contains(explicit.Stdout, "claude-anything-7") || !strings.Contains(explicit.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass an explicit BENCH_MODEL through in an unrouted repo", name))
		}
		partialProbe := runAtEnv(partial, append(envBase, "BENCH_MODEL=claude-anything-7"), "bash", path, "line probe prompt")
		if partialProbe == nil || partialProbe.ExitCode != 0 || !strings.Contains(partialProbe.Stdout, "claude-anything-7") || !strings.Contains(partialProbe.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not fall back to passthrough on an incomplete binding", name))
		}
	}
	return diags
}
