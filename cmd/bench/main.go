// Command bench is the compiled core of the Bench kit. The shell CLI routes every
// compiled subcommand through the production registry below.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/dashboard"
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/guards"
	"github.com/gibbonmi/bench/internal/handoff"
	"github.com/gibbonmi/bench/internal/harness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/models"
	"github.com/gibbonmi/bench/internal/outline"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preprelease"
	"github.com/gibbonmi/bench/internal/publication"
	"github.com/gibbonmi/bench/internal/releasepreflight"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/roadmapflow"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/sessioninspect"
	"github.com/gibbonmi/bench/internal/shift"
	"github.com/gibbonmi/bench/internal/skillsindex"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/stophook"
	"github.com/gibbonmi/bench/internal/structure"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
)

// version is stamped at build time via -ldflags "-X main.version=<pkg.json version>";
// scripts/go-build.sh is the one source of build flags. An unstamped build prints "dev",
// which tells the reader that the binary did not come from the gate or the release
// workflow.
var version = "dev"

func main() {
	// The wrapper's implicit-repair grant is spent once it execs this binary. Scrub it so
	// gate phases and their fixtures never inherit an invocation-dependent privilege.
	os.Unsetenv("BENCH_ALLOW_IMPLICIT_REPAIR")
	var observation io.Writer
	if os.Getenv("BENCH_COMMAND_OBSERVE") == "1" {
		observation = os.Stderr
	}
	os.Exit(Command{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr, Executable: os.Args[0], Observe: observation}.Run(os.Args[1:]))
}

var commandRegistry = []commandDefinition{
	{Name: "anchors", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 15, Suffix: " <path>", Description: "anchors pinning a repo-relative path as TOON (kind, section, needle)"}), Run: outputCommand(anchorsCommand)},
	{Name: "learnings", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 16, Description: "open journal entries as a TOON table (date, title)"}), Run: outputCommand(learnings.Command)},
	{Name: "maps", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 17, Description: "unresolved decision-map tickets as TOON (map, ticket, type, state)"}), Run: outputCommand(maps.Command)},
	{Name: "guards", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 18, Description: "every guard's deny surface as TOON (guard, boundary, denies)"}), Run: outputCommand(guards.Command)},
	{Name: "diff", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 19, Description: "review base + changed files as TOON (--full appends log + diff body; --base freezes source)"}), Run: outputCommand(diff.Command)},
	{Name: "preflight", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 21, Suffix: " review|build <slug>", Description: "phase-entry checks that a spec's artifacts agree with the tree, one verdict row per check"}), Run: outputCommand(preflight.Command)},
	{Name: "coverage", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 20, Suffix: " <spec>", Description: "acceptance-coverage state and rows as TOON (--check to validate)"}), Run: outputCommand(coverage.Command)},
	{Name: "status", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 10, Description: "ambient dashboard: what needs attention + the next action"}), Run: outputCommand(status.Command)},
	{Name: "handoff", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 11, Suffix: " [--harness <name>] [--next <command>]", Description: "print the cold-start pin block and rewrite capture/session-handoff.md"}), Run: outputCommand(handoff.Command)},
	{Name: "commands", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 12, Suffix: " --brief", Description: "print the direct, read-only command probe"}), Run: outputCommand(commandsCommand)},
	{Name: "dashboard", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 13, Suffix: " [--stdout]", Gap: 1, Description: "write a self-contained HTML snapshot of the board (--stdout emits it)"}), Run: outputCommand(dashboard.Command)},
	{Name: "structure", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 6, Description: "flag oversized files + crowded dirs (wire into the gate)"}), Run: outputCommand(structure.Command)},
	{Name: "models", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 5, Description: "list advisory model-id candidates for the line binding"}), Run: outputCommand(models.Command)},
	{Name: "outline", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 23, Suffix: " [path] [--full]", Description: "top-level directory symbol counts as TOON; a path or --full locates candidate seams (file:line), never the project's blessed seams"}), Run: outputCommand(outline.Command)},
	{Name: "idea", AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 8, Suffix: " \"<text>\"", Description: "park an out-of-scope idea in capture/IDEAS.md (commit to nothing)"}), Run: outputCommand(roadmap.IdeaCommand)},
	{Name: "roadmap", AXI: axiApprovedRoot, Inventory: publicInventory(helpRow{Order: 9, Description: "show the top 10 roadmap rows + drain state"}), Run: outputCommand(roadmapCommand)},
	{Name: "skills-index", AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 7, Suffix: " [--check|--write]", Description: "print skills-index drift (default) or regenerate it"}), Run: outputCommand(skillsindex.Command)},
	{Name: "tree-hash", AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: outputCommand(treeHash)},
	{Name: "resolve-model", AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: outputCommand(resolveModel)},
	{Name: "worktree-pool", AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: outputCommand(worktree.PoolCommand)},
	{Name: "worktree-lease-file", AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: outputCommand(worktree.LeaseFileCommand)},
	{Name: "test", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 22, Suffix: " [--full] [package]", Description: "run fresh Go tests and render package, failure, and skip evidence as TOON"}), Run: outputCommand(testCommand)},
	{Name: "help", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(), Kind: commandHelp},
	{Name: "repair", AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 25, Suffix: " [--prune]", Description: "explicitly install the pinned platform binary or prune stale cache entries"}), WrapperOnly: true},

	{Name: "version", Attachment: attachmentDirect, AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 43, Description: "print the installed Bench version (os/arch)"}), Run: versionCommand},
	{Name: "worktree", Attachment: attachmentDirect, AXI: axiApprovedChildren("list"), Inventory: publicInventory(
		helpRow{Order: 31, Suffix: " [--refresh] [objective]", Gap: 1, Description: "create an owned worktree subshell and release it on exit"},
		helpRow{Order: 32, Suffix: " list", Description: "list assignments and registered worktrees as TOON"},
		helpRow{Order: 33, Suffix: " path <target>", Description: "print one active owned worktree's absolute path"},
		helpRow{Order: 34, Suffix: " exec <target> -- <command> [args...]", Description: "run a child directly in an active owned worktree"},
		helpRow{Order: 35, Suffix: " reauthorize --assignment <id> --request <token> --base <commit> --source-tip <commit> <path>", Description: "replace one lost request token after identity proof"},
		helpRow{Order: 36, Suffix: " --help", Description: "show exact list, path, exec, create, release, clean, reclaim, and reauthorize grammar"},
	), Run: worktreeCommand},
	{Name: "resume-clean", Attachment: attachmentDirect, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return worktree.ResumeCleanCommand(args, c.Stdout, c.Stderr) }},
	{Name: "session-inspect", Attachment: attachmentDirect, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return sessioninspect.Command(args, c.Stdout, c.Stderr) }},
	{Name: "shift", Attachment: attachmentDirect, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 38, Suffix: " [--refresh] \"<objective>\"", Gap: 1, Description: "gated loop in a pooled worktree; commit on green"}), Run: func(c Command, args []string) int { return shift.Command(args, c.Stdout, c.Stderr) }},
	{Name: "commit", Attachment: attachmentDirect, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 39, Suffix: " -m <msg> <path>...", Description: "gate, then commit named paths on green"}), Run: func(c Command, args []string) int { return commit.Command(args, c.Stdout, c.Stderr) }},
	{Name: "spec", Attachment: attachmentDirect, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(
		helpRow{Order: 41, Suffix: " retire <slug>", Description: "delete a merged spec + its review pickup (validated)"},
		helpRow{Order: 42, Suffix: " history <slug>", Description: "retire/delete commits for a spec, newest first (TOON)"},
	), Run: outputCommand(spec.Command)},
	{Name: "gate-go", Attachment: attachmentDirect, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return gate.GateGoCommand(args, c.Stdout, c.Stderr) }},
	{Name: "guard-git", Attachment: attachmentDirect, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return guardGit(args, c.Stdin, c.Stdout, c.Stderr) }},
	{Name: "guard-bench-follow-on", Attachment: attachmentDirect, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return guardBenchFollowOn(args, c.Stdin, c.Stdout, c.Stderr) }},
	{Name: "check-agent-line", Attachment: attachmentDirect, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return checkAgentLine(args, c.Stdin, c.Stdout, c.Stderr) }},

	{Name: "setup", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 0, Suffix: " [--plan|--yes]", Description: "inspect, preview, and converge the current repository"}), Run: adoptCommand("setup")},
	{Name: "link", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 1, Suffix: " [copy|symlink]", Description: "safely wire the kit into this repo for every harness"}), Run: adoptCommand("link")},
	{Name: "init", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 2, Description: "scaffold .bench/gate.sh in the current repo"}), Run: adoptCommand("init")},
	{Name: "doctor", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 24, Suffix: " [--fix]", Description: "report (and repair) the PATH shim under a node version manager"}), Run: adoptCommand("doctor")},
	{Name: "unlink", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 3, Suffix: " [--dry-run]", Description: "remove the per-repo Bench footprint the manifest records"}), Run: adoptCommand("unlink")},
	{Name: "upgrade", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(helpRow{Order: 4, Suffix: " [--check] [--force]", Description: "plan and apply a relink onto the installed kit version"}), Run: adoptCommand("upgrade")},
	{Name: "worktree-hook", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return harness.WorktreeCommand(args, c.Stdin, c.Stdout, c.Stderr) }},
	{Name: "gate", Attachment: attachmentSystem, AXI: axiExempt(axiReasonMutation), Inventory: publicInventory(
		helpRow{Order: 26, Suffix: " [--fresh]", Description: "run the project gate (the oracle; --fresh ignores a reusable green)"},
		helpRow{Order: 30, Suffix: " pin", Description: "pin HEAD's .bench tree for pre-push verification"},
		helpRow{Order: 37, Prefix: "bash bin/bench.sh", Suffix: " --fresh", Description: "run the current worktree's gate"},
	), Run: func(c Command, args []string) int { return gate.Command(args, c.Stdin, c.Stdout, c.Stderr) }},
	{Name: "gate-run", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return gate.RunCommand(args, c.Stdout, c.Stderr) }},
	{Name: "gate-pin", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return gate.PinCommand(args, c.Stdin, c.Stdout, c.Stderr) }},
	{Name: "gate-phases", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return gatePhasesCommand(args, c.Stdout, c.Stderr) }},
	{Name: "freshness-check", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return freshnessCheck(args, c.Executable, c.Stderr) }},
	{Name: "freshness-publish", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return freshnessPublish(args, c.Executable, c.Stderr) }},
	{Name: "canary", Attachment: attachmentSystem, AXI: axiExempt(axiReasonOperational), Inventory: publicInventory(helpRow{Order: 14, Suffix: " [root]", Description: "validate fixture inventory"}), Run: func(c Command, args []string) int { return canary.Run(args, c.Stdout, c.Stderr) }},
	{Name: "stop-verdict", Attachment: attachmentSystem, AXI: axiExempt(axiReasonPlumbing), Inventory: internalInventory, Run: func(c Command, args []string) int { return stopVerdict(args, c.Stdin, c.Stderr) }},

	{Name: "release-preflight", Attachment: attachmentShip, AXI: axiExempt(axiReasonRelease), Inventory: publicInventory(helpRow{Order: 28, Suffix: " --mode verify|publish [--profile public|bank] [--phase name]", Description: "run repository release authorization"}), Run: func(c Command, args []string) int { return releasepreflight.Command(args, version, c.Stderr) }},
	{Name: "prep-release", Attachment: attachmentShip, AXI: axiExempt(axiReasonRelease), Inventory: publicInventory(helpRow{Order: 27, Description: "ship-tier rehearsal: artifacts, cross-compile, preflight verify, ship canary"}), Run: func(c Command, args []string) int { return preprelease.Command(args, c.Stdout, c.Stderr) }},
	{Name: "release", Attachment: attachmentShip, AXI: axiExempt(axiReasonRelease), Inventory: publicInventory(helpRow{Order: 29, Suffix: " prepare|submit|promote|rollback|status --version <v> [--profile public|bank] [--root dir] [--registry url] [--path first|staged] [--adapter npm|fixture] [--provenance] [--message text]", Description: "governed npm publication"}), Run: func(c Command, args []string) int { return publication.Command(args, c.Stdout, c.Stderr) }},
}

func outputCommand(fn func([]string) (string, int)) commandHandler {
	return func(c Command, args []string) int {
		out, code := fn(args)
		fmt.Fprint(c.Stdout, out)
		return code
	}
}

func adoptCommand(name string) commandHandler {
	return func(c Command, args []string) int {
		return adopt.Run(append([]string{name}, args...), c.Stdout, c.Stderr, version)
	}
}

func versionCommand(c Command, _ []string) int {
	fmt.Fprintln(c.Stdout, versionLine(version, runtime.GOOS, runtime.GOARCH))
	return 0
}

func helpCommand(c Command, args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(c.Stdout, toon.Usage("bench help", args[0]))
		return 2
	}
	fmt.Fprint(c.Stdout, renderCommandHelp())
	return 0
}

var anchorsGrammar = usage.Grammar{
	Cmd:     "bench anchors",
	Help:    "usage: bench anchors <path>",
	MinArgs: 1,
	MaxArgs: 1,
}

func anchorsCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(anchorsGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	path := anchorQueryPath(root, parsed.Positionals[0])
	var rows [][]string
	for _, anchor := range anchors.Entries() {
		if anchor.File == path {
			rows = append(rows, []string{anchorKindName(anchor.Kind), anchor.Section, anchor.Needle})
		}
	}
	out, err := toon.Table("anchors", []string{"kind", "section", "needle"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	help, err := axi.RenderHelp(nil)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out += help
	return out, 0
}

func anchorQueryPath(root, arg string) string {
	candidate := arg
	if !filepath.IsAbs(candidate) {
		if cwd, err := os.Getwd(); err == nil {
			candidate = filepath.Join(cwd, candidate)
		}
	}
	if _, err := os.Lstat(candidate); err == nil {
		if relative, err := filepath.Rel(root, candidate); err == nil {
			return filepath.ToSlash(filepath.Clean(relative))
		}
	}
	return filepath.ToSlash(filepath.Clean(arg))
}

func anchorKindName(kind anchors.Kind) string {
	switch kind {
	case anchors.Require:
		return "require"
	case anchors.Forbid:
		return "forbid"
	case anchors.RequireInSection:
		return "require-in-section"
	case anchors.ForbidInSection:
		return "forbid-in-section"
	default:
		return "unknown"
	}
}

func testCommand(args []string) (string, int) {
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return testreport.Command(root, args)
}

// commandsGrammar is the declared argument shape usage.Parse enforces for `bench
// commands`. Arity, flag recognition, `--`, and help all come from there rather than a
// local switch. `--brief` is the only form that lists anything, so an invocation without
// it keeps its exit-2 usage answer.
var commandsGrammar = usage.Grammar{
	Cmd:   "bench commands",
	Help:  "usage: bench commands --brief",
	Flags: []usage.Flag{{Name: "--brief"}},
}

func commandsCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(commandsGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if _, brief := parsed.Flags["--brief"]; !brief {
		return commandsGrammar.Help + "\n", 2
	}
	return "version\ncommands --brief\nstatus\n", 0
}

var gatePhasesCommand = gate.PhasesCommand

func roadmapCommand(args []string) (string, int) {
	if len(args) == 0 || len(args) == 1 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help") {
		return roadmap.RoadmapCommand(args)
	}
	// --flow is a mode selector like --context, so it is routed on its leading position
	// and its own grammar reports any misuse that follows.
	if args[0] == "--flow" {
		return roadmapflow.Command(args)
	}
	return roadmap.ContextCommand(args, func(root string) roadmap.GateCacheFact {
		g := status.GateVerdict(root)
		return roadmap.GateCacheFact{Present: g.Present, State: g.State, PendingStatus: g.PendingStatus, Status: g.Status, CachedTree: g.CachedTree, WorkTree: g.WorkTree, Timestamp: g.Timestamp, Stale: g.Stale, CacheBytes: g.CacheBytes}
	})
}

// linesEnv resolves the repo's .bench/lines.env for the two binding consumers,
// resolve-model and check-agent-line. A cwd outside a repo reads as no binding, so the
// verdicts take their unrouted branch rather than denying against an absent oracle. A
// file that is present but fails to read is reported as unreadable rather than absent,
// because a corrupt oracle must announce itself instead of silently disabling
// enforcement.
func linesEnv() lines.Source {
	root, err := git.Root()
	if err != nil {
		return lines.Source{}
	}
	src := lines.Source{Path: filepath.Join(root, ".bench", "lines.env")}
	data, err := os.ReadFile(src.Path)
	if err != nil {
		src.Unreadable = !os.IsNotExist(err)
		src.Exists = src.Unreadable
		return src
	}
	src.Exists = true
	src.Content = data
	return src
}

// harnessFlag is the flag every binding consumer takes: the caller names its own harness
// so the matrix resolves one column and a denial advises in tokens that harness can pass.
var harnessFlag = usage.Flag{Name: "--harness", HasValue: true, NoEmptyValue: true}

var resolveModelGrammar = usage.Grammar{
	Cmd:   "bench resolve-model",
	Help:  "usage: bench resolve-model --harness <" + strings.Join(lines.Harnesses, "|") + ">",
	Flags: []usage.Flag{harnessFlag},
}

var checkAgentLineGrammar = usage.Grammar{
	Cmd:   "bench check-agent-line",
	Help:  "usage: bench check-agent-line --harness <" + strings.Join(lines.Harnesses, "|") + ">",
	Flags: []usage.Flag{harnessFlag},
}

// resolveModel is the `bench resolve-model` plumbing subcommand for the shift adapters.
// It prints the model to pass via the harness --model flag, empty for passthrough, to
// stdout and returns an exit code. Any warning or error goes to os.Stderr directly: the
// map signature carries only stdout, and the adapter captures stdout as the model, so a
// warning must never ride there. BENCH_MODEL names a tier and --harness names the column.
// In a routed repo an unset or unbound tier exits 1 and the adapter refuses to launch.
// The verdict lives in internal/lines, so it is unit-tested without a repo.
func resolveModel(args []string) (string, int) {
	harness, line, code := parseHarness(resolveModelGrammar, args)
	if line != "" {
		fmt.Fprintln(os.Stderr, line)
		return "", code
	}
	benchModel, set := os.LookupEnv("BENCH_MODEL")
	model, code, stderr := lines.ResolveModelVerdict(harness, benchModel, set, linesEnv())
	if stderr != "" {
		fmt.Fprintln(os.Stderr, stderr)
	}
	if model == "" {
		return "", code
	}
	return model + "\n", code
}

// parseHarness applies g to args and returns the named harness. A missing --harness is a
// misuse like any other usage error, so both binding consumers answer it the same way
// rather than guessing a column; the returned line is non-empty exactly when the caller
// must print it and exit with code.
func parseHarness(g usage.Grammar, args []string) (harness, line string, code int) {
	parsed, line, code := usage.Parse(g, args)
	if line != "" {
		return "", line, code
	}
	harness, named := parsed.Flags[harnessFlag.Name]
	if !named {
		return "", g.Help, 2
	}
	return harness, "", 0
}

// checkAgentLine is the delegation guard subcommand. It reads the Agent PreToolUse
// envelope on stdin, reads the binding through internal/lines, and yields the verdict as
// an exit code: 0 allow, including a degraded warn-and-allow with its WARNING on stderr,
// or 2 deny, with the DENIED message on stderr. The deferred recover maps any panic to 3,
// so exit 2 means only an intentional deny, and the shim's fail-open rim catches a crash.
func checkAgentLine(args []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	harness, line, usageCode := parseHarness(checkAgentLineGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return usageCode
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		data = nil // unreadable stdin reads as unparseable → fail open
	}
	exit, msg := lines.AgentLineVerdict(data, harness, linesEnv())
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
	// objective is a capture-quality guard only. The ledger stores the entry key, not the
	// text, so a missing objective still skips a meaningless capture.
	objective := envelope.ToolInput.Description
	if objective == "" {
		objective = sanitize.Preview(envelope.ToolInput.Prompt)
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
	entry := intent.NewEntry(intent.KindClaudeAgent)
	entry.Key = envelope.ToolUseID
	if err := intent.Upsert(root, entry); err != nil {
		fmt.Fprintf(stderr, "WARNING: check-agent-line: intent capture failed: %v\n", err)
	}
}

// stopVerdict is the completion-oracle subcommand. It reads the Stop envelope on stdin
// and takes the resolved wrapper as args[0], because the shim passes it so gate
// resolution stays in bin/bench.sh. It orchestrates the verdict through internal/stophook:
// honoring stop_hook_active, enforcing only when BENCH_SHIFT=1, running `<wrapper> gate`,
// writing the verdict cache, and returning 0 allow or 2 block. A panic maps to 3, which
// the shim treats as a core error and fails open, with no forged verdict, like a missing core.
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
// the gate subject's tree identity, kept as a standalone subcommand for external callers.
// The gate owner calls git.TreeHash directly, so the hash still has one source. Root is
// args[0] when given, else the cwd's repo. It prints the content hash or `none`; a
// non-hash tells the caller to fail safe rather than forge.
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

// guardGit is the destructive-git guard subcommand. It reads the PreToolUse envelope on
// stdin, classifies through internal/gitguard, and yields the verdict as an exit code:
// 0 allow, 2 block, with the `BLOCKED:` message on stderr, or 3 a genuine failure to run.
// The deferred recover maps any panic to 3, not Go's default exit-2, so exit 2 means
// only an intentional block, and the shim can trust it.
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

func guardBenchFollowOn(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if recover() != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	command, err := benchguard.CommandFromEnvelope(data)
	if err != nil {
		fmt.Fprintln(stderr, "WARNING: block-bench-follow-on: unreadable command field — allowing Bash.")
		return 0
	}
	if !benchguard.Classify(command, benchguard.DefaultResolver()) {
		return 0
	}
	fmt.Fprintln(stderr, benchguard.BlockMessage())
	return 2
}

func worktreeCommand(c Command, args []string) int {
	if len(args) > 0 && args[0] == "exec" {
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return 1
		}
		return worktree.ExecCommand(root, args[1:], c.Stdin, c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "path" {
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return 1
		}
		return worktree.PathCommand(root, args[1:], c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "list" {
		out, code := worktree.ListCommand(args[1:])
		fmt.Fprint(c.Stdout, out)
		return code
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		fmt.Fprint(c.Stdout, usage.WorktreeUsage())
		return 0
	}
	if len(args) > 0 && args[0] == "create" {
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return 1
		}
		return worktree.CreateCommand(root, args[1:], c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "release" {
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return 1
		}
		return worktree.ReleaseCommand(root, args[1:], c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "clean" {
		return worktree.CleanCommand(args[1:], c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "reclaim" {
		return worktree.ReclaimCommand(args[1:], c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "reauthorize" {
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return 1
		}
		return worktree.ReauthorizeCommand(root, args[1:], c.Stdout, c.Stderr)
	}
	if len(args) > 0 && args[0] == "land" {
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return 1
		}
		return worktree.LandCommand(root, c.Executable, args[1:], c.Stdout, c.Stderr)
	}
	// `recovery` is no longer a worktree subcommand. This family's fallback is a
	// free-form objective, so naming it here reports the removed verb instead of
	// opening a subshell called "recovery".
	if len(args) > 0 && args[0] == "recovery" {
		fmt.Fprintln(c.Stderr, toon.Usage("bench worktree", args[0]))
		return 2
	}
	return worktree.Subshell(args, c.Stdin, c.Stdout, c.Stderr)
}

// versionLine renders the single line `bench version` prints. Kept as a pure
// function so the format has one source and a table test can pin it without a
// process boundary.
func versionLine(v, goos, goarch string) string {
	return fmt.Sprintf("bench %s (%s/%s)", v, goos, goarch)
}
