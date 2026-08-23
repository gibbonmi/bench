package conformance

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// matrixBinding is the routed fixture every runtime check resolves against. codex and
// claude hold their own columns, and opencode is deliberately unadopted. No cell is
// shared between columns, so a surface reading the wrong column names the wrong token.
const matrixBinding = "BENCH_CODEX_TOP=gpt-5.4\n" +
	"BENCH_CODEX_MID=gpt-5.3-codex-spark\n" +
	"BENCH_CODEX_CHEAP=openai/gpt-5\n" +
	"BENCH_CLAUDE_TOP=fable-5\n" +
	"BENCH_CLAUDE_MID=opus-4-8\n" +
	"BENCH_CLAUDE_CHEAP=sonnet-5\n"

// openCodeBinding adds the provider-qualified opencode column, the one shape that harness's
// cells accept.
const openCodeBinding = matrixBinding +
	"BENCH_OPENCODE_TOP=openai/gpt-5.6-sol\n" +
	"BENCH_OPENCODE_MID=openai/gpt-5.6-terra\n" +
	"BENCH_OPENCODE_CHEAP=openai/gpt-5.6-luna\n"

// coreless returns env with a PATH carrying neither a bench wrapper nor the stub dir. The
// shims' wrapper search then comes up empty, and each one takes its own missing-core rim.
func coreless(env []string) []string {
	out := make([]string, 0, len(env)+1)
	for _, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			continue
		}
		out = append(out, kv)
	}
	return append(out, "PATH=/usr/bin:/bin")
}

func checkAgentHookBehavior(root string) []string {
	hook := filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")
	realBench := filepath.Join(root, "bin", "bench.sh")
	if !exists(hook) {
		return nil
	}
	if !exists(realBench) {
		probe := runWithInput(root, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"claude-nonexistent-9"}}`, "bash", hook)
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

	routed, cleanupRouted, err := tempGitRepoWithLines(matrixBinding)
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
	hookCase("allows a bound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus-4-8"}}`, "", 0)
	// Enforcement stays permissive across the whole matrix. Another harness's bound cell is
	// still a declared line, even though the advice narrows to the asking harness.
	hookCase("allows another harness's bound cell", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-5.3-codex-spark"}}`, "", 0)
	hookCase("captures allowed intent", routed, `{"tool_name":"Agent","tool_use_id":"allowed-1","tool_input":{"description":"ship safely","prompt":"fallback","model":"opus-4-8"}}`, "", 0)
	ledgerPath := filepath.Join(routed, ".git", "bench-intent.json")
	ledger, err := os.ReadFile(ledgerPath)
	if err != nil || !strings.Contains(string(ledger), `"key":"allowed-1"`) || strings.Contains(string(ledger), `"objective"`) {
		diags = append(diags, fmt.Sprintf("check-agent-line intent capture contract failed: allowed ledger=%q err=%v", ledger, err))
	}
	hookCase("replays allowed intent idempotently", routed, `{"tool_name":"Agent","tool_use_id":"allowed-1","tool_input":{"description":"ship safely","prompt":"fallback","model":"opus-4-8"}}`, "", 0)
	replayed, replayErr := os.ReadFile(ledgerPath)
	if replayErr != nil || !bytes.Equal(ledger, replayed) {
		diags = append(diags, fmt.Sprintf("check-agent-line intent capture contract failed: replay changed ledger bytes: before=%q after=%q err=%v", ledger, replayed, replayErr))
	}
	// A tier name is not a token the Agent tool can pass. Only bound cells allow it.
	hookCase("denies a tier name", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"mid"}}`, "", 2)
	hookCase("denies an unbound model", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-9"}}`, "harness claude binds top=fable-5", 2)
	// The advice is the asking harness's own column and nothing else. A Claude session never
	// reads a recovery instruction naming ids it cannot pass.
	if probe := runWithInputEnv(routed, env, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-9"}}`, "bash", hook); probe != nil && strings.Contains(probe.Stderr, "gpt-5.3-codex-spark") {
		diags = append(diags, "check-agent-line.sh deny names another harness's column: "+probe.Stderr)
	}
	hookCase("denied intent is not captured", routed, `{"tool_name":"Agent","tool_use_id":"denied-1","tool_input":{"description":"must not persist","model":"gpt-9"}}`, "", 2)
	if ledger, err := os.ReadFile(ledgerPath); err != nil || strings.Contains(string(ledger), "denied-1") {
		diags = append(diags, fmt.Sprintf("check-agent-line intent capture contract failed: denied call changed ledger=%q err=%v", ledger, err))
	}
	hookCase("does not fail open on malformed stdin", routed, `not json at all`, "not parseable as JSON", 0)
	// Ratified posture flip (enforcement-verification): in a routed repo whose claude column
	// is complete, a missing model DENIES (exit 2) rather than warning. An omitted model
	// inherits the session's model, the silent-escalation path invariant #2 exists to stop.
	hookCase("denies a missing model field in a routed repo", routed, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "bound tier token", 2)

	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	hookCase("does not fail open without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-9"}}`, "no .bench/lines.env", 0)
	// The missing-model deny is gated on routing. With no binding to enforce, a missing
	// model keeps the fail-open rim, the residual the decision deliberately preserves.
	hookCase("does not fail open on a missing model without lines.env", unrouted, `{"tool_name":"Agent","tool_input":{"prompt":"x"}}`, "no resolvedModel/model field", 0)

	partial, cleanupPartial, err := tempGitRepoWithLines("BENCH_CODEX_TOP=gpt-5.4\nBENCH_CODEX_MID=\nBENCH_CODEX_CHEAP=openai/gpt-5\n")
	if err != nil {
		return []string{"check-agent-line.sh setup failed: " + err.Error()}
	}
	defer cleanupPartial()
	hookCase("does not fail open on an incomplete claude column", partial, `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-9"}}`, "claude column is incomplete", 0)

	// This is the guard's own rim. With no reachable core, the shim warns and ALLOWS,
	// because a guard that dies with its core would block every delegation.
	corelessProbe := runWithInputEnv(routed, coreless(env), `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-9"}}`, "bash", hook)
	if corelessProbe == nil || corelessProbe.ExitCode != 0 || !strings.Contains(corelessProbe.Stderr, "bench core not found") {
		diags = append(diags, fmt.Sprintf("check-agent-line.sh does not fail open when the core is unreachable: %+v", corelessProbe))
	}
	return diags
}

// harnessOf maps each adapter to the harness column it must ask the core for. An adapter
// that asks for another harness's column launches the wrong family.
var harnessOf = map[string]string{"claude": "claude", "codex": "codex", "opencode": "opencode"}

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
			continue
		}
		if !strings.Contains(text, `resolve-model --harness `+harnessOf[name]) {
			diags = append(diags, fmt.Sprintf("adapter %s does not name its own harness when resolving the line", name))
		}
	}
	if bindir == "" {
		return diags
	}

	routed, cleanupRouted, err := tempGitRepoWithLines(matrixBinding)
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupRouted()
	openCodeRouted, cleanupOpenCodeRouted, err := tempGitRepoWithLines(openCodeBinding)
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupOpenCodeRouted()
	unrouted, cleanupUnrouted, err := tempGitRepoWithLines("")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupUnrouted()
	os.Remove(filepath.Join(unrouted, ".bench", "lines.env"))
	unbound, cleanupUnbound, err := tempGitRepoWithLines("# no cell is bound here\n")
	if err != nil {
		return append(diags, "adapter line guard setup failed: "+err.Error())
	}
	defer cleanupUnbound()

	for _, name := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", name)
		if !exists(path) {
			continue
		}
		harness := harnessOf[name]
		envBase := append(conformanceSubprocessEnv(), "PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"))
		// opencode's column is unbound in the shared fixture, so its launch cases run
		// against the provider-qualified fixture instead.
		boundRepo := routed
		if name == "opencode" {
			boundRepo = openCodeRouted
		}
		// The adapter must launch on exactly what the core resolves for its own harness and
		// tier. A shim that recomputes resolution passes a "calls resolve-model" check while
		// drifting from the core. The test therefore compares the values directly.
		core := runWithInputEnv(boundRepo, append(envBase, "BENCH_MODEL=mid"), "", filepath.Join(bindir, "bench"), "resolve-model", "--harness", harness)
		if core == nil || core.ExitCode != 0 || strings.TrimSpace(core.Stdout) == "" {
			diags = append(diags, fmt.Sprintf("adapter %s core comparison failed: bench resolve-model --harness %s produced %+v", name, harness, core))
			continue
		}
		coreModel := strings.TrimSpace(core.Stdout)
		bound := runWithInputEnv(boundRepo, append(envBase, "BENCH_MODEL=mid"), "--line probe prompt", "bash", path)
		if bound == nil || bound.ExitCode != 0 || !strings.Contains(bound.Stdout, coreModel) || !strings.Contains(bound.Stdout, "--line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not launch on the core's model %s with a dash-leading prompt", name, coreModel))
		}
		if name == "codex" && (bound == nil || !strings.Contains(bound.Stdout, "--sandbox\nworkspace-write")) {
			diags = append(diags, "adapter codex routed path does not select the workspace-write sandbox")
		}
		if name == "opencode" {
			// An unadopted harness fails closed. It has no fallback to another harness's column.
			unboundColumn := runWithInputEnv(routed, append(envBase, "BENCH_MODEL=mid"), "line probe prompt", "bash", path)
			if unboundColumn == nil || unboundColumn.ExitCode == 0 || !strings.Contains(unboundColumn.Stderr, "opencode column is unbound") {
				diags = append(diags, fmt.Sprintf("adapter opencode does not refuse to launch while its column is unbound: %+v", unboundColumn))
			}
		}
		unset := runWithInputEnv(boundRepo, envBase, "line probe prompt", "bash", path)
		if unset != nil && unset.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse undeclared BENCH_MODEL in a routed repo", name))
		}
		unboundTier := runWithInputEnv(boundRepo, append(envBase, "BENCH_MODEL=gpt-9"), "line probe prompt", "bash", path)
		if unboundTier != nil && unboundTier.ExitCode == 0 {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse an unbound BENCH_MODEL in a routed repo", name))
		}
		pass := runWithInputEnv(unrouted, envBase, "line probe prompt", "bash", path)
		if pass == nil || pass.ExitCode != 0 || !strings.Contains(pass.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass through in an unrouted repo", name))
		}
		if name == "codex" && (pass == nil || !strings.Contains(pass.Stdout, "--sandbox\nworkspace-write")) {
			diags = append(diags, "adapter codex unrouted path does not select the workspace-write sandbox")
		}
		explicit := runWithInputEnv(unrouted, append(envBase, "BENCH_MODEL=gpt-anything-7"), "line probe prompt", "bash", path)
		if explicit == nil || explicit.ExitCode != 0 || !strings.Contains(explicit.Stdout, "gpt-anything-7") || !strings.Contains(explicit.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not pass an explicit BENCH_MODEL through in an unrouted repo", name))
		}
		unboundProbe := runWithInputEnv(unbound, append(envBase, "BENCH_MODEL=gpt-anything-7"), "line probe prompt", "bash", path)
		if unboundProbe == nil || unboundProbe.ExitCode != 0 || !strings.Contains(unboundProbe.Stdout, "gpt-anything-7") || !strings.Contains(unboundProbe.Stdout, "line probe prompt") {
			diags = append(diags, fmt.Sprintf("adapter %s does not fall back to passthrough on a binding that binds no cell", name))
		}
		// The adapters' rim is the mirror of the guard's. Where the hook fails open, an adapter
		// refuses, because an unguarded passthrough in a routed repo is silent de-enforcement.
		corelessProbe := runWithInputEnv(boundRepo, coreless(append(envBase, "BENCH_MODEL=mid")), "line probe prompt", "bash", path)
		if corelessProbe == nil || corelessProbe.ExitCode != 1 || !strings.Contains(corelessProbe.Stderr, "bench core not found") {
			diags = append(diags, fmt.Sprintf("adapter %s does not refuse to launch when the core is unreachable: %+v", name, corelessProbe))
		}
	}
	return diags
}

// checkLineHarnessSurfaces holds every runtime surface to one core. Each surface names
// its own harness and resolves that harness's column. The kit CLI, a linked repo's
// wrapper, the guard, and the three adapters therefore cannot drift apart. Cases run per
// surface, so an omitted one names itself.
func checkLineHarnessSurfaces(root string) []string {
	realBench := filepath.Join(root, "bin", "bench.sh")
	hook := filepath.Join(root, ".bench", "hooks", "check-agent-line.sh")
	if !exists(realBench) || !exists(hook) {
		return nil
	}
	bindir, cleanup, err := adapterStubDir(realBench)
	if err != nil {
		return []string{"line harness surface setup failed: " + err.Error()}
	}
	defer cleanup()
	routed, cleanupRouted, err := tempGitRepoWithLines(openCodeBinding)
	if err != nil {
		return []string{"line harness surface setup failed: " + err.Error()}
	}
	defer cleanupRouted()
	env := append(conformanceSubprocessEnv(),
		"PATH="+bindir+string(os.PathListSeparator)+os.Getenv("PATH"),
		"BENCH_MODEL=mid")

	var diags []string
	for _, surface := range []struct {
		name, want string
		args       []string
		input      string
		wantExit   int
	}{
		// This case is the kit's own CLI, invoked from the kit checkout the way a session runs
		// it.
		{name: "kit CLI", want: "gpt-5.3-codex-spark", args: []string{"bash", realBench, "resolve-model", "--harness", "codex"}},
		// A linked repo reaches the same core through the resolved wrapper on PATH.
		{name: "linked-repo CLI", want: "opus-4-8", args: []string{"bash", filepath.Join(bindir, "bench"), "resolve-model", "--harness", "claude"}},
		// The guard names claude, so its denial advises in the claude column.
		{name: "claude hook", want: "harness claude binds top=fable-5", args: []string{"bash", hook}, input: `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"gpt-9"}}`, wantExit: 2},
		{name: "codex adapter", want: "gpt-5.3-codex-spark", args: []string{"bash", filepath.Join(root, ".bench", "adapters", "codex")}, input: "prompt"},
		{name: "claude adapter", want: "opus-4-8", args: []string{"bash", filepath.Join(root, ".bench", "adapters", "claude")}, input: "prompt"},
		{name: "opencode adapter", want: "openai/gpt-5.6-terra", args: []string{"bash", filepath.Join(root, ".bench", "adapters", "opencode")}, input: "prompt"},
	} {
		probe := runWithInputEnv(routed, env, surface.input, surface.args...)
		if probe == nil || probe.ExitCode != surface.wantExit {
			got := -1
			if probe != nil {
				got = probe.ExitCode
			}
			diags = append(diags, fmt.Sprintf("line surface %s exit %d (want %d)", surface.name, got, surface.wantExit))
			continue
		}
		if !strings.Contains(probe.Stdout, surface.want) && !strings.Contains(probe.Stderr, surface.want) {
			diags = append(diags, fmt.Sprintf("line surface %s does not resolve its own harness column (want %q)", surface.name, surface.want))
		}
	}
	return diags
}
