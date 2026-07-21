package conformance

import (
	"strings"
	"testing"
)

func offlineSmokeDeniesRepairAndEgress(smoke string) bool {
	requiredOnce := []string{
		`repair_disabled="${BENCH_OFFLINE_REPAIR_DISABLED:-1}"`,
		`offline_mode=true`,
		`[[ "$repair_disabled" == 1 && "$offline_mode" == true ]]`,
		`BENCH_OFFLINE_ALLOWED_ORIGIN="$registry_origin"`,
		`go|cargo|cc|gcc|clang|make)`,
		`bwrap --unshare-net`,
		`sandbox-exec -p '(version 1) (allow default) (deny network*)'`,
	}
	for _, anchor := range requiredOnce {
		if strings.Count(smoke, anchor) != 1 {
			return false
		}
	}
	return strings.Count(smoke, `BENCH_NO_REPAIR="$repair_disabled"`) == 2 &&
		strings.Count(smoke, "BENCH_NO_REPAIR=1") == 8 &&
		strings.Count(smoke, "BENCH_OFFLINE=1") == 18 &&
		strings.Count(smoke, `npm_config_offline="$offline_mode"`) == 1 &&
		strings.Count(smoke, "npm_config_offline=true") == 2 &&
		strings.Count(smoke, `NODE_OPTIONS="--require=$root/scripts/offline-network-sentinel.cjs"`) == 4 &&
		strings.Count(smoke, `[[ ! -s "$egress_log" ]]`) == 4 &&
		strings.Count(smoke, "BENCH_OFFLINE_ALLOWED_ORIGIN=") == 1
}

func offlineSmokeEnumeratesSliceOneSuppressions(smoke string) bool {
	requiredExecutableProof := []string{
		`printf '#!/bin/sh\n: > "$BENCH_OFFLINE_REPAIR_MARKER"\nexit 97\n' > "$probe_bin/node"`,
		`repair_output="$(HOME="$local_home" BENCH_HOME="$local_home/.bench" BENCH_OFFLINE=1 BENCH_OFFLINE_REPAIR_MARKER="$repair_marker" PATH="$probe_bin:$PATH" bash "$installed" models 2>&1)"`,
		`[[ "$repair_exit" == 127 && "$repair_output" == *"repair suppressed by BENCH_OFFLINE=1"* && ! -e "$repair_marker" ]]`,
		`printf '#!/bin/sh\n: > "$BENCH_OFFLINE_CODEX_MARKER"\nexit 97\n' > "$probe_bin/codex"`,
		`BENCH_OFFLINE=1 OPENAI_API_KEY=sentinel ANTHROPIC_API_KEY=sentinel BENCH_OFFLINE_CODEX_MARKER="$codex_marker"`,
		`for source in codex openai anthropic; do`,
		`grep -q "^  ${source},offline,offline,BENCH_OFFLINE=1$" "$models_output"`,
		`[[ ! -e "$codex_marker" ]]`,
		`if [ "${1:-}" = -C ] && [ "${3:-}" = fetch ]; then : > "$BENCH_OFFLINE_GIT_MARKER"; fi`,
		`BENCH_OFFLINE=1 BENCH_OFFLINE_GIT_MARKER="$git_marker" BENCH_OFFLINE_REAL_GIT="$(command -v git)"`,
		`worktree create --refresh --request offline-smoke --label offline-smoke`,
		`grep -q '^  offline,BENCH_OFFLINE=1$' "$refresh_output"`,
		`[[ ! -e "$git_marker" ]]`,
		`printf 'offline suppression: repair,git_refresh,codex_live,codex_bundled,openai_http,anthropic_http BENCH_OFFLINE=1 zero_attempts\n'`,
	}
	for _, anchor := range requiredExecutableProof {
		if strings.Count(smoke, anchor) != 1 {
			return false
		}
	}
	return true
}

func TestOfflineSmokeSliceOneProofIsExecutableNotTokenOnly(t *testing.T) {
	smoke := NewHarness(t).ReadRootFile("scripts", "smoke-offline.sh")
	if !offlineSmokeEnumeratesSliceOneSuppressions(smoke) {
		t.Fatal("real offline smoke does not carry the executable six-operation proof")
	}
	tokenOnly := `printf 'offline suppression: repair,git_refresh,codex_live,codex_bundled,openai_http,anthropic_http BENCH_OFFLINE=1 zero_attempts\n'`
	if offlineSmokeEnumeratesSliceOneSuppressions(tokenOnly) {
		t.Fatal("enumeration token alone passed without executable suppression probes")
	}
	mutated := strings.Replace(smoke, "openai_http,anthropic_http", "openai_http", 1)
	if offlineSmokeEnumeratesSliceOneSuppressions(mutated) {
		t.Fatal("omitting the Anthropic operation from the six-operation enumeration did not bite")
	}
}
