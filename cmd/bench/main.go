// Command bench is the compiled core of the Bench kit — the strangler target the
// shell CLI routes ported subcommands into. Dispatch is a `commands` map of the ported
// AXI query subcommands (learnings, maps, guards, diff, coverage, worktree list),
// each resolving repo state and returning its stdout plus an exit code, plus a direct
// `version` case that needs the build-time GOOS/GOARCH rather than repo state. Every
// later slice adds names to that map; the shell router (bin/bench.sh) grows names,
// not mechanisms.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/dashboard"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/harness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/models"
	"github.com/gibbonmi/bench/internal/outline"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/publication"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/sessioninspect"
	"github.com/gibbonmi/bench/internal/shift"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/stophook"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
)

// version is stamped at build time via -ldflags "-X main.version=<pkg.json version>"
// (see scripts/go-build.sh — the one source of build flags). An unstamped build
// prints "dev", which is the tell that the binary was not produced by the gate or
// the release workflow.
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// commands are the ported AXI query subcommands: each resolves arguments and repo
// state, then returns its complete stdout plus an exit code — the dispatch below just
// writes and exits. This map is the one seam that grows per slice; version stays a
// direct case because it needs the build-time GOOS/GOARCH, not repo state.
var commands = map[string]func([]string) (string, int){
	"learnings":           learnings.Command,
	"maps":                maps.Command,
	"guards":              guards.Command,
	"diff":                diff.Command,
	"coverage":            coverage.Command,
	"status":              status.Command,
	"commands":            commandsCommand,
	"dashboard":           dashboard.Command,
	"structure":           structure.Command,
	"models":              models.Command,
	"outline":             outline.Command,
	"idea":                roadmap.IdeaCommand,
	"roadmap":             roadmapCommand,
	"tree-hash":           treeHash,
	"resolve-model":       resolveModel,
	"worktree-pool":       worktree.PoolCommand,
	"worktree-lease-file": worktree.LeaseFileCommand,
}

func commandsCommand(args []string) (string, int) {
	if len(args) != 1 || args[0] != "--brief" {
		return "usage: bench commands --brief\n", 2
	}
	return "version\ncommands --brief\nstatus\n", 0
}

var gatePhasesCommand = gate.PhasesCommand

func roadmapCommand(args []string) (string, int) {
	if len(args) == 0 {
		return roadmap.RoadmapCommand(nil)
	}
	return roadmap.ContextCommand(args, func(root string) roadmap.GateCacheFact {
		g := status.GateVerdict(root)
		return roadmap.GateCacheFact{Present: g.Present, State: g.State, PendingStatus: g.PendingStatus, Status: g.Status, CachedTree: g.CachedTree, WorkTree: g.WorkTree, Timestamp: g.Timestamp, Stale: g.Stale, CacheBytes: g.CacheBytes}
	})
}

// linesEnv resolves the repo's .bench/lines.env — its path, whether it exists, and its
// content — for the two binding consumers (resolve-model and check-agent-line). A cwd
// outside a repo, or an unreadable file, reads as no binding (exists=false): the
// verdicts then take their unrouted / fail-open branch, never denying against an
// absent oracle.
func linesEnv() (path string, exists bool, content []byte) {
	root, err := git.Root()
	if err != nil {
		return "", false, nil
	}
	path = filepath.Join(root, ".bench", "lines.env")
	data, err := os.ReadFile(path)
	if err != nil {
		return path, false, nil
	}
	return path, true, data
}

// resolveModel is the `bench resolve-model` plumbing subcommand for the shift adapters:
// it prints the model to pass via the harness --model flag (empty for passthrough) to
// stdout and returns an exit code. Any warning/error goes to os.Stderr directly — the
// map signature carries only stdout, and the adapter captures stdout AS the model, so a
// warning must never ride there. In a routed repo an unset or unbound BENCH_MODEL exits
// 1 and the adapter refuses to launch. The optional modes validate the exact tier id,
// then either return its corresponding BENCH_ALIAS_* value or require provider/model
// compatibility; the verdicts live in internal/lines so they are unit-tested without a
// repo.
func resolveModel(args []string) (string, int) {
	resolve := lines.ResolveModelVerdict
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "usage: bench resolve-model [--alias | --provider-model]")
		return "", 2
	}
	if len(args) == 1 {
		switch args[0] {
		case "--alias":
			resolve = lines.ResolveModelAliasVerdict
		case "--provider-model":
			resolve = lines.ResolveProviderModelVerdict
		default:
			fmt.Fprintln(os.Stderr, "usage: bench resolve-model [--alias | --provider-model]")
			return "", 2
		}
	}
	benchModel, set := os.LookupEnv("BENCH_MODEL")
	path, exists, content := linesEnv()
	model, code, stderr := resolve(benchModel, set, exists, path, content)
	if stderr != "" {
		fmt.Fprintln(os.Stderr, stderr)
	}
	if model == "" {
		return "", code
	}
	return model + "\n", code
}

// checkAgentLine is the delegation guard subcommand: it reads the Agent PreToolUse
// envelope on stdin, reads the binding through internal/lines, and yields the verdict as
// an exit code — 0 allow (or any degraded warn-and-allow, with its WARNING on stderr), 2
// deny (with the DENIED message on stderr). The deferred recover maps any panic to 3, so
// exit 2 means only an intentional deny and the shim's fail-open rim catches a crash.
func checkAgentLine(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		data = nil // unreadable stdin reads as unparseable → fail open
	}
	_, exists, content := linesEnv()
	exit, msg := lines.AgentLineVerdict(data, exists, content)
	if msg != "" {
		fmt.Fprintln(stderr, msg)
	}
	if exit == 0 {
		captureClaudeAgentIntent(data, stderr)
	}
	return exit
}

func captureClaudeAgentIntent(data []byte, stderr io.Writer) {
	var envelope struct {
		ToolName  string `json:"tool_name"`
		ToolUseID string `json:"tool_use_id"`
		ToolInput struct {
			Description string `json:"description"`
			Prompt      string `json:"prompt"`
		} `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		fmt.Fprintf(stderr, "WARNING: check-agent-line: intent capture skipped: malformed envelope: %v\n", err)
		return
	}
	if envelope.ToolName != "" && envelope.ToolName != "Agent" {
		return
	}
	objective := envelope.ToolInput.Description
	if objective == "" {
		objective = intent.Preview(envelope.ToolInput.Prompt)
	}
	if envelope.ToolUseID == "" || objective == "" {
		fmt.Fprintln(stderr, "WARNING: check-agent-line: intent capture skipped: missing tool_use_id or objective")
		return
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "WARNING: check-agent-line: intent capture skipped outside a repository")
		return
	}
	entry := intent.NewEntry(intent.KindClaudeAgent, objective)
	entry.Key = envelope.ToolUseID
	if err := intent.Upsert(root, entry); err != nil {
		fmt.Fprintf(stderr, "WARNING: check-agent-line: intent capture failed: %v\n", err)
	}
}

// stopVerdict is the completion-oracle subcommand: it reads the Stop envelope on stdin,
// takes the resolved wrapper as args[0] (the shim passes it so gate resolution stays in
// bin/bench.sh), and orchestrates the verdict through internal/stophook — honoring
// stop_hook_active, enforcing only when BENCH_SHIFT=1, running `<wrapper> gate`, writing
// the verdict cache, and returning 0 allow / 2 block. A panic maps to 3, which the shim
// treats as a core error and fails open (no forged verdict), exactly like a missing core.
func stopVerdict(args []string, stdin io.Reader, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	if len(args) == 0 || args[0] == "" {
		fmt.Fprintln(stderr, "bench stop-verdict: missing wrapper argument")
		return 3
	}
	data, _ := io.ReadAll(stdin)
	armed := os.Getenv("BENCH_SHIFT") == "1"
	return stophook.Run(data, args[0], armed, stderr)
}

// treeHash exposes git.TreeHash as the `bench tree-hash [root]` plumbing subcommand:
// gate subject's tree identity, kept as a standalone subcommand for external callers.
// The gate owner calls git.TreeHash directly, so the hash still has one source. Root is
// args[0] when given, else the cwd's
// repo. It prints the content hash or `none`; a non-hash tells the caller to fail safe
// rather than forge.
func treeHash(args []string) (string, int) {
	var root string
	if len(args) > 0 && args[0] != "" {
		root = args[0]
	} else {
		r, err := git.Root()
		if err != nil {
			return "none\n", 0
		}
		root = r
	}
	return git.TreeHash(root) + "\n", 0
}

// guardGit is the destructive-git guard subcommand: it reads the PreToolUse envelope on
// stdin, classifies through internal/gitguard, and yields the verdict as an exit code —
// 0 allow, 2 block (with the `BLOCKED:` message on stderr), 3 a genuine failure to run.
// The deferred recover maps any panic to 3, not Go's default exit-2, so exit 2 means
// only an intentional block and the shim can trust it.
func guardGit(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	command := gitguard.CommandFromEnvelope(data)
	if command == "" {
		return 0
	}
	chk := gitguard.Checker{RefResolves: git.RefResolves, BranchExists: git.BranchExists}
	label := gitguard.Classify(command, chk)
	if label == "" {
		return 0
	}
	fmt.Fprintln(stderr, gitguard.BlockMessage(label))
	return 2
}

func run(args []string, stdout, stderr *os.File) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "bench: no subcommand")
		return 2
	}
	if fn, ok := commands[args[0]]; ok {
		out, code := fn(args[1:])
		fmt.Fprint(stdout, out)
		return code
	}
	switch args[0] {
	case "link", "init", "doctor", "unlink":
		return adopt.Run(args, stdout, stderr, version)
	case "version":
		fmt.Fprintln(stdout, versionLine(version, runtime.GOOS, runtime.GOARCH))
		return 0
	case "worktree":
		if len(args) > 1 && args[1] == "list" {
			out, code := worktree.ListCommand(args[2:])
			fmt.Fprint(stdout, out)
			return code
		}
		if len(args) == 2 && (args[1] == "help" || args[1] == "--help" || args[1] == "-h") {
			fmt.Fprint(stdout, usage.WorktreeUsage())
			return 0
		}
		if len(args) > 1 && args[1] == "create" {
			root, err := git.Root()
			if err != nil {
				fmt.Fprintln(stderr, toon.NotInRepo())
				return 1
			}
			return worktree.CreateCommand(root, args[2:], stdout, stderr)
		}
		if len(args) > 1 && args[1] == "release" {
			root, err := git.Root()
			if err != nil {
				fmt.Fprintln(stderr, toon.NotInRepo())
				return 1
			}
			return worktree.ReleaseCommand(root, args[2:], stdout, stderr)
		}
		if len(args) > 1 && args[1] == "clean" {
			return worktree.CleanCommand(args[2:], stdout, stderr)
		}
		if len(args) > 1 && args[1] == "recovery" {
			return worktree.RecoveryCommand(args[2:], stdout, stderr)
		}
		return worktree.Subshell(args[1:], os.Stdin, stdout, stderr)
	case "resume-clean":
		return worktree.ResumeCleanCommand(args[1:], stdout, stderr)
	case "session-inspect":
		return sessioninspect.Command(args[1:], stdout, stderr)
	case "worktree-hook":
		return harness.WorktreeCommand(args[1:], os.Stdin, stdout, stderr)
	case "shift":
		return shift.Command(args[1:], os.Stdin, stdout, stderr)
	case "commit":
		return commit.Command(args[1:], stdout, stderr)
	case "spec":
		out, code := spec.Command(args[1:])
		fmt.Fprint(stdout, out)
		return code
	case "gate-run":
		return gate.RunCommand(args[1:], stdout, stderr)
	case "gate-pin":
		return gate.PinCommand(args[1:], os.Stdin, stdout, stderr)
	case "gate-phases":
		return gatePhasesCommand(args[1:], stdout, stderr)
	case "release-preflight":
		return preflight.Command(args[1:], version, stderr)
	case "release":
		return publication.Command(args[1:], stdout, stderr)
	case "canary":
		return canary.Run(args[1:], stdout, stderr)
	case "guard-git":
		return guardGit(args[1:], os.Stdin, stdout, stderr)
	case "check-agent-line":
		return checkAgentLine(args[1:], os.Stdin, stdout, stderr)
	case "stop-verdict":
		return stopVerdict(args[1:], os.Stdin, stderr)
	default:
		fmt.Fprintf(stderr, "bench: unknown subcommand: %q\n", args[0])
		return 2
	}
}

// versionLine renders the single line `bench version` prints. Kept as a pure
// function so the format has one source and a table test can pin it without a
// process boundary.
func versionLine(v, goos, goarch string) string {
	return fmt.Sprintf("benchkit %s (%s/%s)", v, goos, goarch)
}
