// Package env is the deep module behind the subprocess environment Bench
// constructs for the harness adapter a shift launches. That launch builds its
// environment from a documented passlist, instead of inheriting the parent
// process's environment verbatim. A committed `.bench/env.allow` is the
// only way to widen it. Build is the one call: policy, opt-in parsing, and
// construction all sit behind it, so a second ad hoc filter never drifts from
// this one.
//
// This package serves the harness-adapter (agent) launch only. The project
// gate's environment is a separate, already-closed subject: FT78's
// manifest-declared closure, owned by internal/gate. That closure launches the
// gate script with PATH plus only the names declared under `environment` in
// `.bench/gate-inputs.json`. This feature pins that closure with a sentinel
// contract rather than rebuilding it here, so there is no gate passlist in this
// package.
//
// The default passlist is exported as data (SharedBasics, AgentPasslist). This
// lets a conformance check read the exact values the enforcement uses, rather
// than a second transcription of them. That is the single-source requirement
// this package exists to satisfy.
package env

import (
	"os"
	"path/filepath"
	"strings"
)

// SharedBasics is the process-basics set the adapter carries: enough for a
// subprocess to resolve paths, locate its home directory, and render sanely in
// a terminal. HOME and XDG_CONFIG_HOME are load-bearing for git's own config
// resolution; LC_* is a glob because exact-name matching breaks real-system
// locales.
var SharedBasics = []string{
	"PATH", "HOME", "USER", "LOGNAME", "SHELL", "TMPDIR", "TERM", "COLORTERM",
	"LANG", "LC_*", "XDG_*",
}

// benchFamily is the BENCH_* glob the adapter carries.
const benchFamily = "BENCH_*"

// agentAdditions is what the agent class adds beyond SharedBasics: its own
// BENCH_* namespace, and the shipped adapters' documented harness variables.
// These cover Claude Code, Codex, opencode, and the cloud credential chains
// Claude Code documents for Bedrock and Vertex routing. Every name here is
// cited to an official page in DATA_HANDLING.md.
//
// A default glob must not straddle two families. Every glob here covers a
// namespace a single owner controls, so a name whose prefix is shared with a
// foreign family is enumerated exactly instead. This is the rule that kept a
// `GO*` glob out of the retired gate-class draft. That glob would also have
// matched GOOGLE_APPLICATION_CREDENTIALS and handed a subprocess a cloud
// credential, and it binds every future addition to this list.
//
// A unit-test edge row checks each default glob against a fixture of foreign
// names. So an exact set cannot later be "simplified" into a wider glob.
var agentAdditions = []string{
	benchFamily,
	"ANTHROPIC_*", "CLAUDE_CODE_*", "CLAUDE_CONFIG_DIR", "API_TIMEOUT_MS",
	"CODEX_*", "RUST_LOG", "SSL_CERT_FILE",
	"OPENCODE_*", "OPENAI_API_KEY",
	"AWS_*", "GOOGLE_*", "GCLOUD_PROJECT", "CLOUD_ML_REGION", "VERTEX_LOCATION",
}

// AgentPasslist is the default name patterns (exact names and PREFIX* globs)
// the adapter admits, before any .bench/env.allow [agent] additions.
var AgentPasslist = concat(SharedBasics, agentAdditions)

func concat(lists ...[]string) []string {
	var out []string
	for _, l := range lists {
		out = append(out, l...)
	}
	return out
}

// Build returns the ordered environment the harness adapter should launch with:
// every parent-environment name matching the default passlist, plus any matching
// .bench/env.allow [agent] entry. Both come in the parent's own order, with
// values passed through byte for byte. repoRoot is the resolved repo root; this
// package never discovers it. A malformed .bench/env.allow fails closed. Build
// returns a nil slice and an error naming the offending line and reason. It
// never degrades to defaults.
func Build(repoRoot string) ([]string, error) {
	allow, err := parseAllowFile(filepath.Join(repoRoot, ".bench", "env.allow"))
	if err != nil {
		return nil, err
	}

	patterns := append([]string(nil), AgentPasslist...)
	patterns = append(patterns, allow.agent...)

	var result []string
	for _, kv := range os.Environ() {
		name, _, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		if matchesAny(name, patterns) {
			result = append(result, kv)
		}
	}
	return result, nil
}

// matchesAny reports whether name matches any pattern: an exact name, or a
// PREFIX* glob matched by prefix. Patterns never match on value — only on
// name — so a passlisted variable's value passes through unaltered including a
// multi-line or very large one.
func matchesAny(name string, patterns []string) bool {
	for _, p := range patterns {
		if strings.HasSuffix(p, "*") {
			if strings.HasPrefix(name, p[:len(p)-1]) {
				return true
			}
			continue
		}
		if name == p {
			return true
		}
	}
	return false
}
