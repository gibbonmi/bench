package conformance

import "regexp"

// parityRow states how one entry point reaches the command registry. A row with a `must`
// pattern is static: the check grades its text, because a run needs a live pool, a
// release, or a CI runner. Every other row is runtime: the check runs the entry with one
// canned input and reads the observed registry id. A runtime row with `compare` also runs
// the direct verb, so a shim that holds a second opinion is red.
type parityRow struct {
	command string
	must    *regexp.Regexp
	direct  []string
	input   string
	env     []string
	compare bool
}

func (r parityRow) isStatic() bool { return r.must != nil }

// parityVerb builds the pattern that matches a verb and nothing that only starts with it.
// A trailing word break would accept `release-preflight-2`, because a dash is a word
// break. Every static row that names a verb builds its pattern here.
func parityVerb(prefix, verb string) *regexp.Regexp {
	return regexp.MustCompile(prefix + regexp.QuoteMeta(verb) + `([^0-9A-Za-z-]|$)`)
}

// parityWrapperToken stands for the resolved wrapper path in a direct argv. stop.sh passes
// the wrapper it resolved to the core, so the direct call must pass the same one.
const parityWrapperToken = "{{wrapper}}"

// parityBenignBash is a Bash envelope both guards allow. It proves the shim reaches the
// core, not that the core denies; the guards own their own verdict checks.
const parityBenignBash = `{"tool_name":"Bash","tool_input":{"command":"ls"}}`

// parityBenignWrite is a file-tool envelope the write guard allows. The path sits outside
// every repository, so the row proves the shim reaches the core without depending on the
// graded root's own checkout kind.
const parityBenignWrite = `{"tool_name":"Write","tool_input":{"file_path":"/nonexistent/parity/x.txt"},"cwd":"/"}`

// entryPointParity maps every shim, adapter, front door, and CI entry to the one registry
// command it must reach. The check enumerates the shim and adapter directories from disk,
// so an entry outside this table is red rather than ungraded.
var entryPointParity = map[string]parityRow{
	".bench/hooks/block-dangerous-git.sh": {
		command: "guard-git", direct: []string{"guard-git"}, input: parityBenignBash, compare: true,
	},
	".bench/hooks/block-bench-follow-on.sh": {
		command: "guard-bench-follow-on", direct: []string{"guard-bench-follow-on"}, input: parityBenignBash, compare: true,
	},
	".bench/hooks/block-primary-file-write.sh": {
		command: "guard-file-write", direct: []string{"guard-file-write"}, input: parityBenignWrite, compare: true,
	},
	".bench/hooks/block-unsafe-rg.sh": {
		command: "guard-rg", direct: []string{"guard-rg"}, input: parityBenignBash, compare: true,
	},
	".bench/hooks/check-agent-line.sh": {
		command: "check-agent-line",
		direct:  []string{"check-agent-line", "--harness", "claude"},
		input:   `{"tool_name":"Agent","tool_input":{"prompt":"x","model":"opus-4-8"}}`,
		compare: true,
	},
	// stop_hook_active short-circuits the verdict, so the row costs no gate run.
	".bench/hooks/stop.sh": {
		command: "stop-verdict",
		direct:  []string{"stop-verdict", parityWrapperToken},
		input:   `{"stop_hook_active":true}`,
		compare: true,
	},
	// session-start.sh prints one advisory line before the core's output, so the comparison
	// is a suffix rather than an equality.
	".bench/hooks/session-start.sh": {
		command: "session-inspect", direct: []string{"session-inspect"}, compare: true,
	},
	// A create needs a live pool, so this row is static.
	".bench/hooks/worktree-lifecycle.sh": {
		command: "worktree-hook", must: parityVerb(`(?m)^exec "\$cmd" `, "worktree-hook"),
	},
	// An adapter launches its harness, so its stdout is the harness's. Only the observed id
	// is comparable.
	".bench/adapters/claude": {
		command: "resolve-model", env: []string{"BENCH_MODEL=mid"},
	},
	".bench/adapters/codex": {
		command: "resolve-model", env: []string{"BENCH_MODEL=mid"},
	},
	".bench/adapters/opencode": {
		command: "resolve-model", env: []string{"BENCH_MODEL=mid"},
	},
	// The wrapper's bare invocation is the second front door. It must print what the routed
	// verb prints, byte for byte.
	"bin/bench.sh": {
		command: "status", direct: []string{"status", "--route"}, compare: true,
	},
	// The CI chain is two static halves: the script execs the verb, and the workflow runs
	// the script. A staged binary would prove no more than the exec line does. The exec line
	// may carry an env prefix, because the script passes the binary it built to the child
	// phases through BENCH_RUN_BINARY.
	"scripts/release-preflight.sh": {
		command: "release-preflight", must: parityVerb(`(?m)^exec (env BENCH_RUN_BINARY="\$binary" )?"\$binary" `, "release-preflight"),
	},
	".github/workflows/native-runtime.yml": {
		command: "release-preflight", must: regexp.MustCompile(`bash scripts/release-preflight\.sh\b`),
	},
	// Both phase adapters route through this file, so the verb it names must be the exact
	// one the registry owns.
	".agents/commands/bench.md": {
		command: "status", must: parityVerb("bench status --", "route"),
	},
}

// Each reason covers a whole class of plumbing verbs, so the file states it once.
const (
	whyWrapperVerb = "bin/bench.sh calls it, and bench-sh-routes grades that route"
	whyGatePhase   = "the gate runner starts it as a phase, so the gate's own verdict grades the route"
	whyBuildScript = "scripts/go-build.sh calls it after it links the binary"
)

// parityExemptCommands names each internal registry command no parity row reaches, with
// the reason it needs none. A command in neither the rows nor this table is red, so the
// next plumbing verb cannot escape the table.
var parityExemptCommands = map[string]string{
	"tree-hash":           whyWrapperVerb,
	"worktree-pool":       whyWrapperVerb,
	"worktree-lease-file": whyWrapperVerb,
	"resume-clean":        whyWrapperVerb,
	"gate-run":            whyWrapperVerb,
	"gate-background":     whyWrapperVerb,
	"freshness-check":     whyWrapperVerb,
	"gate-go":             whyGatePhase,
	"gate-prose":          whyGatePhase,
	"gate-phases":         whyGatePhase,
	"freshness-publish":   whyBuildScript,
}
