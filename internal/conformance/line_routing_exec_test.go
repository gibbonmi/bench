package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

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

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=gpt-5.3-codex-spark\nBENCH_TIER_CHEAP=openai/gpt-5\nBENCH_ALIAS_MID=opus\n")
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
	hookCase("allows a bound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-5.3-codex-spark"}}`, "", 0)
	hookCase("allows a declared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus"}}`, "", 0)
	hookCase("denies an undeclared alias", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"sonnet"}}`, "", 2)
	hookCase("denies an unbound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-9"}}`, "", 2)
	hookCase("does not fail open on malformed stdin", routed, `not json at all`, "not parseable as JSON", 0)
	// Ratified posture flip (enforcement-verification): in a routed, completely-bound repo a
	// missing model DENIES (exit 2) rather than warning — an omitted model inherits the
	// session's model, the silent-escalation path invariant #2 exists to stop.
	hookCase("denies a missing model field in a routed repo", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "bound alias", 2)

	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	hookCase("does not fail open without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-9"}}`, "no .bench/lines.env", 0)
	// The missing-model deny is gated on routing: with no binding to enforce, a missing
	// model keeps the fail-open rim (the residual the decision deliberately preserves).
	hookCase("does not fail open on a missing model without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "no resolvedModel/model field", 0)

	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=openai/gpt-5\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupPartial()
	hookCase("does not fail open on an incomplete binding", partial, `{"tool_name":"Agent","tool_input":{"prompt":"x","resolvedModel":"gpt-9"}}`, "unset or empty", 0)
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

	routed, cleanupRouted, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=gpt-5.3-codex-spark\nBENCH_TIER_CHEAP=openai/gpt-5\n")
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
	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_TIER_TOP=gpt-5.4\nBENCH_TIER_MID=\nBENCH_TIER_CHEAP=openai/gpt-5\n")
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
		bound := runAtEnv(routed, append(envBase, "BENCH_MODEL=gpt-5.3-codex-spark"), "bash", path, "--line probe prompt")
		if bound == nil || bound.ExitCode != 0 || !strings.Contains(bound.Stdout, "gpt-5.3-codex-spark") || !strings.Contains(bound.Stdout, "--line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass BENCH_MODEL and a dash-leading prompt to the harness in a routed repo", name))
		}
		unset := runAtEnv(routed, envBase, "bash", path, "line probe prompt")
		if unset != nil && unset.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
		}
		unbound := runAtEnv(routed, append(envBase, "BENCH_MODEL=gpt-9"), "bash", path, "line probe prompt")
		if unbound != nil && unbound.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
		pass := runAtEnv(unrouted, envBase, "bash", path, "line probe prompt")
		if pass == nil || pass.ExitCode != 0 || !strings.Contains(pass.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass through in an unrouted repo", name))
		}
		explicit := runAtEnv(unrouted, append(envBase, "BENCH_MODEL=gpt-anything-7"), "bash", path, "line probe prompt")
		if explicit == nil || explicit.ExitCode != 0 || !strings.Contains(explicit.Stdout, "gpt-anything-7") || !strings.Contains(explicit.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass an explicit BENCH_MODEL through in an unrouted repo", name))
		}
		partialProbe := runAtEnv(partial, append(envBase, "BENCH_MODEL=gpt-anything-7"), "bash", path, "line probe prompt")
		if partialProbe == nil || partialProbe.ExitCode != 0 || !strings.Contains(partialProbe.Stdout, "gpt-anything-7") || !strings.Contains(partialProbe.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not fall back to passthrough on an incomplete binding", name))
		}
	}
	return diags
}
