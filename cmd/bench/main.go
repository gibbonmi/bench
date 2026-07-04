// Command bench is the compiled core of the Bench kit — the strangler target the
// shell CLI routes ported subcommands into. Dispatch is a `commands` map of the ported
// AXI query subcommands (learnings, maps, guards, diff, coverage), each resolving repo
// state and returning its stdout plus an exit code, plus a direct `version` case that
// needs the build-time GOOS/GOARCH rather than repo state. Every later slice adds names
// to that map; the shell router (bin/bench.sh) grows names, not mechanisms.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/models"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/shift"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/stophook"
	"github.com/gibbonmi/bench/internal/structure"
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
	"structure":           structure.Command,
	"models":              models.Command,
	"idea":                roadmap.IdeaCommand,
	"roadmap":             roadmap.RoadmapCommand,
	"tree-hash":           treeHash,
	"resolve-model":       resolveModel,
	"worktree-pool":       worktree.PoolCommand,
	"worktree-lease-file": worktree.LeaseFileCommand,
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
// 1 and the adapter refuses to launch. Aliases do not bind here (only tier ids do); the
// verdict itself lives in internal/lines so it is unit-tested without a repo.
func resolveModel(args []string) (string, int) {
	benchModel, set := os.LookupEnv("BENCH_MODEL")
	path, exists, content := linesEnv()
	model, code, stderr := lines.ResolveModelVerdict(benchModel, set, exists, path, content)
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
// `--describe-binding` emits the live denies clause to stdout without reading stdin,
// feeding the shim's `--describe`.
func checkAgentLine(args []string, stdin io.Reader, stdout, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	if len(args) > 0 && args[0] == "--describe-binding" {
		_, exists, content := linesEnv()
		fmt.Fprintln(stdout, lines.DescribeBinding(exists, content))
		return 0
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		data = nil // unreadable stdin reads as unparseable → fail open
	}
	_, exists, content := linesEnv()
	exit, msg := lines.AgentLineVerdict(data, exists, content)
	if msg != "" {
		fmt.Fprintln(stderr, msg)
	}
	return exit
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
// the gate verdict-cache key, kept as a standalone subcommand for external callers. The
// in-process verdict recorders (internal/gate.Record and the Stop hook) call git.TreeHash
// directly, so the hash still has one source. Root is args[0] when given, else the cwd's
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
// only an intentional block and the shim can trust it. `--describe-classes` prints the
// deny surface to stdout without reading stdin, feeding the shim's `--describe`.
func guardGit(args []string, stdin io.Reader, stdout, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	if len(args) > 0 && args[0] == "--describe-classes" {
		fmt.Fprintln(stdout, gitguard.DescribeClasses())
		return 0
	}
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
	case "link", "init", "doctor":
		return adopt.Run(args, stdout, stderr, version)
	case "version":
		fmt.Fprintln(stdout, versionLine(version, runtime.GOOS, runtime.GOARCH))
		return 0
	case "worktree":
		return worktree.Subshell(os.Stdin, stdout, stderr)
	case "shift":
		return shift.Command(args[1:], os.Stdin, stdout, stderr)
	case "gate-run":
		return gate.RunCommand(args[1:], stdout, stderr)
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
