package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// This is the idiom-setting table test for the module: pure logic, table-driven, no
// process boundary. Acceptance for the version line lives at the shell seam, the gate's
// version-routing contract, and this test pins the format one layer down.
func TestVersionLine(t *testing.T) {
	cases := []struct {
		v, goos, goarch, want string
	}{
		{"0.2.0", "linux", "amd64", "bench 0.2.0 (linux/amd64)"},
		{"dev", "darwin", "arm64", "bench dev (darwin/arm64)"},
		{"1.0.0", "linux", "arm64", "bench 1.0.0 (linux/arm64)"},
	}
	for _, c := range cases {
		if got := versionLine(c.v, c.goos, c.goarch); got != c.want {
			t.Errorf("versionLine(%q,%q,%q) = %q, want %q", c.v, c.goos, c.goarch, got, c.want)
		}
	}
}

func TestRunVersionExits0(t *testing.T) {
	if rc := (Command{}).Run([]string{"version"}); rc != 0 {
		t.Errorf("run version exit = %d, want 0", rc)
	}
}

func TestRunUnknownExits2(t *testing.T) {
	if rc := (Command{}).Run([]string{"nope"}); rc != 2 {
		t.Errorf("run nope exit = %d, want 2", rc)
	}
}

func TestRunGateRejectsBriefUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"gate", "--brief"})
	if code != 2 || stdout.Len() != 0 || stderr.String() != "usage: bench gate [--fresh|pin]\n" {
		t.Fatalf("gate --brief = (stdout %q, stderr %q, exit %d), want usage on stderr and exit 2", stdout.String(), stderr.String(), code)
	}
}

func TestRunStatusRouteEmitsOneNextRow(t *testing.T) {
	stdout := tempFile(t)
	if code := (Command{Stdout: stdout}).Run([]string{"status", "--route"}); code != 0 {
		t.Fatalf("status --route exit = %d, want 0", code)
	}
	if got := readFile(t, stdout); !strings.HasPrefix(got, "next[1]{state,why,command}:\n") {
		t.Fatalf("status --route = %q, want one next row", got)
	}
}

func TestGuardsQueryReportsFollowOnGuardThroughCommandSeam(t *testing.T) {
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"guards"})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("bench guards = stdout=%q stderr=%q exit=%d, want stdout/0", stdout.String(), stderr.String(), code)
	}
	for _, want := range []string{
		"block-bench-follow-on",
		"PreToolUse:Bash",
		"Bench shell follow-ons",
		"claude,codex",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("bench guards = %q, want %q", stdout.String(), want)
		}
	}
}

func TestHarnessesQueryProjectsTheRecordThroughCommandSeam(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"harnesses"})
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("bench harnesses = stdout=%q stderr=%q exit=%d, want stdout/0", stdout.String(), stderr.String(), code)
	}
	for _, want := range []string{
		"schema: 1\n",
		"harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}:\n",
		"help[0]{cmd,why}:\n",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("bench harnesses = %q, want %q", stdout.String(), want)
		}
	}
}

func TestHelpRendersPublicCommandRegistryRows(t *testing.T) {
	old := commandRegistry
	t.Cleanup(func() { commandRegistry = old })
	commandRegistry = []commandDefinition{
		{Name: "help", Inventory: publicInventory(), Kind: commandHelp},
		{
			Name:      "status",
			Inventory: publicInventory(helpRow{Order: 1, Description: "prove the public command owns this row"}),
		},
	}

	for _, spelling := range []string{"help", "--help", "-h"} {
		t.Run(spelling, func(t *testing.T) {
			var stdout bytes.Buffer
			if code := (Command{Stdout: &stdout}).Run([]string{spelling}); code != 0 {
				t.Fatalf("%s exit = %d, want 0", spelling, code)
			}
			if want := helpInventoryTitle + "\n  bench status               prove the public command owns this row\n"; stdout.String() != want {
				t.Fatalf("%s stdout = %q, want registry rows %q", spelling, stdout.String(), want)
			}
		})
	}
}

func TestHelpInventoryIsComplete(t *testing.T) {
	// The independently authored expectation is the omission oracle. Deriving it from
	// commandRegistry would let a deleted public or child row disappear from both sides.
	const want = `bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants.
  bench setup [--plan|--yes]  inspect, preview, and converge the current repository
  bench link [copy|symlink]  safely wire the kit into this repo for every harness
  bench init                 scaffold .bench/gate.sh in the current repo
  bench unlink [--dry-run]   remove the per-repo Bench footprint the manifest records
  bench upgrade [--check] [--force]  plan and apply a relink onto the installed kit version
  bench models               list advisory model-id candidates for the line binding
  bench structure            flag oversized files + crowded dirs (wire into the gate)
  bench cache                report the Bench Go build cache footprint (bytes, files, last trim)
  bench cache clean          take the cache lock and empty the Bench Go build cache (refuses under a live run)
  bench skills-index [--check|--write]  print skills-index drift (default) or regenerate it
  bench idea "<text>"        park an out-of-scope idea in capture/IDEAS.md (commit to nothing)
  bench learning "<title>" --what --right [--rule]  append one open entry to capture/learnings.md (the drain verdicts it)
  bench retro <slug> --body <markdown>  validate and create one primary-local implementation retrospective
  bench roadmap              show the top 10 roadmap rows + drain state
  bench status               ambient dashboard: what needs attention + the next action
  bench handoff [--harness <name>] [--next <command>]  print the cold-start pin block and rewrite capture/session-handoff.md
  bench commands --brief     print the direct, read-only command probe
  bench dashboard [--stdout] write a self-contained HTML snapshot of the board (--stdout emits it)
  bench canary [root]        validate fixture inventory
  bench anchors <path>       anchors pinning a repo-relative path as TOON (kind, section, needle)
  bench learnings            open journal entries as a TOON table (date, title)
  bench maps                 unresolved decision-map tickets as TOON (map, ticket, type, state)
  bench guards               every guard's deny surface as TOON (guard, boundary, denies)
  bench diff                 review base + changed files as TOON (--full appends log + diff body; --base freezes source)
  bench harnesses [<harness>]  the harness record as TOON (harness, provider, phase_form, hooks, delegation_guard); one name prints that harness's cells
  bench coverage <spec>      acceptance-coverage state and rows as TOON (--check to validate)
  bench preflight review|build <slug>  phase-entry checks that a spec's artifacts agree with the tree, one verdict row per check
  bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name>  run focused Go-test or named-check evidence as TOON; no gate verdict
  bench outline [path] [--full]  top-level directory symbol counts as TOON; a path or --full locates candidate seams (file:line), never the project's blessed seams
  bench consumers <qualified-symbol> [--full]  every resolved Go reference edge as TOON (file:line, via, enclosing); identifies edges, never blessed seams
  bench doctor [--fix]       report (and repair) the PATH shim under a node version manager
  bench repair [--prune]     explicitly install the pinned platform binary or prune stale cache entries
  bench gate [--fresh]       run the project gate (the oracle; --fresh ignores a reusable green)
  bench prep-release         ship-tier rehearsal: artifacts, cross-compile, preflight verify, ship canary
  bench release-preflight --mode verify|publish [--profile public|bank] [--phase name]  run repository release authorization
  bench release prepare|submit|promote|rollback|status --version <v> [--profile public|bank] [--root dir] [--registry url] [--path first|staged] [--adapter npm|fixture] [--provenance] [--message text]  governed npm publication
  bench gate pin             pin HEAD's .bench tree for pre-push verification
  bench worktree shell [--refresh] [objective] create an owned worktree subshell and release it on exit
  bench worktree list        list assignments and registered worktrees as TOON
  bench worktree path <target>  print one active owned worktree's absolute path
  bench worktree exec <target> [--env KEY=VALUE]... -- <command> [args...]  run a child directly in an active owned worktree
  bench worktree show <target> <rev>:<path>  print one blob from a revision of an active owned worktree
  bench worktree build <target>  build an active owned worktree's tree into its own dist/bench
  bench worktree exec <target> -- bench gate  run one active owned worktree's gate
  bench worktree reauthorize --assignment <id> --request <token> --base <commit> --source-tip <commit> <path>  replace one lost request token after identity proof
  bench worktree merge --from <commit|target> <target>  merge a default-branch commit or a sibling's tip into an owned worktree
  bench worktree --help      show exact list, path, exec, show, build, create, release, clean, reclaim, reauthorize, and merge grammar
  bench shift [--refresh] "<objective>" gated loop in a pooled worktree; commit on green
  bench commit -m <msg> <path>...  gate, then commit named paths on green
  bench spec retire <slug>   delete a merged spec + its review pickup (validated)
  bench spec history <slug>  retire/delete commits for a spec, newest first (TOON)
  bench version              print the installed Bench version (os/arch)
`

	var stdout bytes.Buffer
	if code := (Command{Stdout: &stdout}).Run([]string{"help"}); code != 0 {
		t.Fatalf("help exit = %d, want 0", code)
	}
	if stdout.String() != want {
		t.Fatalf("help inventory:\n%s\nwant complete public inventory:\n%s", stdout.String(), want)
	}
}

func TestRootAndHelpAlignWrapperAndBinary(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	// The three root reads are byte-compared, so they must observe one repository and one
	// pool that nothing else can move between reads. A concurrent session's worktree or
	// commit in the live kit tree would otherwise change the route table under the test.
	observed := gittest.RepoOnBranch(t, "main")
	home := filepath.Join(t.TempDir(), "bench-home")
	t.Setenv("BENCH_HOME", home)
	t.Chdir(observed)

	var directRoot bytes.Buffer
	if code := (Command{Stdout: &directRoot}).Run(nil); code != 0 {
		t.Fatalf("in-process root exit = %d, want 0", code)
	}
	if !strings.HasPrefix(directRoot.String(), "next[1]{state,why,command}:\n") {
		t.Fatalf("in-process root = %q, want next route table", directRoot.String())
	}

	binary := filepath.Join(t.TempDir(), "bench")
	build := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, binary)
	cleanEnv := append(capability.WithoutEnvironment(
		capability.WithoutEnvironment(os.Environ(), runbinary.Env), "BENCH_HOME"), "BENCH_HOME="+home)
	build.Dir = root
	build.Env = cleanEnv
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build bench: %v\n%s", err, out)
	}

	command := func(path string, args ...string) *exec.Cmd {
		t.Helper()
		cmd := exec.Command(path, args...)
		cmd.Dir = observed
		cmd.Env = cleanEnv
		if path != binary {
			cmd.Env = append(capability.WithoutEnvironment(runbinary.WithEnv(cleanEnv, binary), "BENCH_KIT"), "BENCH_KIT="+root)
		}
		return cmd
	}
	run := func(path string, args ...string) string {
		t.Helper()
		cmd := command(path, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", path, args, err, out)
		}
		return string(out)
	}

	binaryRoot := run(binary)
	if !strings.HasPrefix(binaryRoot, "next[1]{state,why,command}:\n") {
		t.Fatalf("binary root = %q, want next route table", binaryRoot)
	}
	if binaryRoot != directRoot.String() {
		t.Fatalf("binary root = %q, in-process root = %q", binaryRoot, directRoot.String())
	}
	wrapper := filepath.Join(root, "bin", "bench.sh")
	if wrapperRoot := run(wrapper); wrapperRoot != binaryRoot {
		t.Fatalf("wrapper root = %q, binary root = %q", wrapperRoot, binaryRoot)
	}

	binaryHelp := run(binary, "help")
	if !strings.HasPrefix(binaryHelp, "bench — Pocock"+" pipeline") {
		t.Fatalf("binary help = %q, want inventory", binaryHelp)
	}
	if !strings.Contains(binaryHelp, "bench canary [root]        validate fixture inventory") {
		t.Fatalf("binary help missing canary inventory wording:\n%s", binaryHelp)
	}
	if strings.Contains(binaryHelp, "run the gate against known-broken fixtures") {
		t.Fatalf("binary help retained stale canary execution wording:\n%s", binaryHelp)
	}
	for _, spelling := range []string{"help", "--help", "-h"} {
		if spellingHelp := run(binary, spelling); spellingHelp != binaryHelp {
			t.Errorf("binary %s = %q, binary help = %q", spelling, spellingHelp, binaryHelp)
		}
		if wrapperHelp := run(wrapper, spelling); wrapperHelp != binaryHelp {
			t.Errorf("wrapper %s = %q, binary help = %q", spelling, wrapperHelp, binaryHelp)
		}
	}

	// bench harnesses projects a record compiled into the binary, so the wrapper and the
	// binary must agree byte for byte on both projections. A wrapper that reshapes the TOON
	// on the way through breaks every agent that parses it.
	for _, probe := range []struct {
		argv   []string
		header string
	}{
		{argv: []string{"harnesses"}, header: "schema: 1\nharnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}:\n"},
		{argv: []string{"harnesses", "codex"}, header: "schema: 1\ncells[13]{field,value,source,checked}:\n"},
	} {
		binaryHarnesses := run(binary, probe.argv...)
		if !strings.HasPrefix(binaryHarnesses, probe.header) {
			t.Errorf("binary %v = %q, want the %q header", probe.argv, binaryHarnesses, probe.header)
		}
		if wrapperHarnesses := run(wrapper, probe.argv...); wrapperHarnesses != binaryHarnesses {
			t.Errorf("wrapper %v = %q, binary %v = %q", probe.argv, wrapperHarnesses, probe.argv, binaryHarnesses)
		}
	}

	cmd := command(wrapper, "help", "extra")
	out, err := cmd.CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 2 || string(out) != "usage: bench help (unknown argument: extra)\n" {
		t.Fatalf("wrapper help extra = (output %q, error %v), want help usage and exit 2", out, err)
	}
}

func TestRunHelpRejectsTrailingArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"help", "extra"})
	if code != 2 || stdout.String() != "usage: bench help (unknown argument: extra)\n" || stderr.Len() != 0 {
		t.Fatalf("help extra = (stdout %q, stderr %q, exit %d), want help usage on stdout and exit 2", stdout.String(), stderr.String(), code)
	}
}

func TestHelpKeepsStatusPublicRoute(t *testing.T) {
	t.Run("inventory", func(t *testing.T) {
		var stdout bytes.Buffer
		if code := (Command{Stdout: &stdout}).Run([]string{"help"}); code != 0 {
			t.Fatalf("help exit = %d, want 0", code)
		}
		const row = "  bench status               ambient dashboard: what needs attention + the next action\n"
		if !strings.Contains(stdout.String(), row) {
			t.Fatalf("help omitted independently required public status row %q", row)
		}
	})
	t.Run("dispatch", func(t *testing.T) {
		var stdout bytes.Buffer
		if code := (Command{Stdout: &stdout}).Run([]string{"status", "--help"}); code != 0 {
			t.Fatalf("status --help exit = %d, want 0", code)
		}
		if !strings.HasPrefix(stdout.String(), "usage: bench status") {
			t.Fatalf("status --help = %q, want status grammar", stdout.String())
		}
	})
}

// TestResolveModelHarnessFlag drives the CLI's argument surface. --harness selects the
// column, and the retired --alias and --provider-model spellings are rejected rather than
// quietly resolving a model, so there is only one way to ask the binding a question.
func TestResolveModelHarnessFlag(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	binding := "BENCH_CODEX_TOP=gpt-5.6-sol\nBENCH_CODEX_MID=gpt-5.6-terra\nBENCH_CODEX_CHEAP=gpt-5.6-luna\n" +
		"BENCH_CLAUDE_TOP=fable\nBENCH_CLAUDE_MID=opus\nBENCH_CLAUDE_CHEAP=sonnet\n"
	if err := os.WriteFile(filepath.Join(root, ".bench", "lines.env"), []byte(binding), 0o644); err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("BENCH_MODEL", "cheap")

	for _, tt := range []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{"codex column", []string{"--harness", "codex"}, "gpt-5.6-luna\n", 0},
		{"claude column", []string{"--harness", "claude"}, "sonnet\n", 0},
		{"unbound column", []string{"--harness", "opencode"}, "", 1},
		{"unknown harness", []string{"--harness", "gemini"}, "", 1},
		{"missing harness", nil, "", 2},
		{"retired alias flag", []string{"--alias"}, "", 2},
		{"retired provider-model flag", []string{"--provider-model"}, "", 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out, code := resolveModel(tt.args)
			if out != tt.wantOut || code != tt.wantCode {
				t.Fatalf("resolveModel(%v) = (%q, %d), want (%q, %d)", tt.args, out, code, tt.wantOut, tt.wantCode)
			}
		})
	}
}

// TestCheckAgentLineHarnessFlag pins the guard's own argument surface. It takes the same
// --harness flag, and a retired flag is a usage error rather than a silent allow.
func TestCheckAgentLineHarnessFlag(t *testing.T) {
	for _, tt := range []struct {
		name     string
		args     []string
		wantCode int
	}{
		{"missing harness", nil, 2},
		{"retired alias flag", []string{"--alias"}, 2},
		{"unknown harness", []string{"--harness", "gemini"}, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			code := checkAgentLine(tt.args, strings.NewReader(`{"tool_input":{"model":"opus"}}`), nil, &stderr)
			if code != tt.wantCode {
				t.Fatalf("checkAgentLine(%v) = %d, want %d (stderr=%q)", tt.args, code, tt.wantCode, stderr.String())
			}
		})
	}
}

func TestRunCanaryRetainsPositionalGrammar(t *testing.T) {
	for _, help := range []string{"help", "--help", "-h"} {
		stdout, stderr := tempFile(t), tempFile(t)
		if rc := (Command{Stdout: stdout, Stderr: stderr}).Run([]string{"canary", help}); rc != 0 {
			t.Errorf("canary %s exit = %d, want 0", help, rc)
		}
		if got := readFile(t, stdout); got != "usage: bench canary [root]\n" {
			t.Errorf("canary %s stdout = %q, want exact usage", help, got)
		}
	}

	stderr := tempFile(t)
	if rc := (Command{Stderr: stderr}).Run([]string{"canary", "one", "two"}); rc != 2 {
		t.Fatalf("run canary too-many-arguments exit = %d, want 2", rc)
	}
	if got := readFile(t, stderr); !strings.Contains(got, "usage: bench canary") || !strings.Contains(got, "unknown argument: two") {
		t.Fatalf("canary too-many-arguments stderr = %q, want usage and offending argument", got)
	}

	stderr = tempFile(t)
	if rc := (Command{Stderr: stderr}).Run([]string{"canary", filepath.Join(t.TempDir(), "missing")}); rc != 1 {
		t.Fatalf("run canary invalid root exit = %d, want 1", rc)
	}
}

func TestRunGatePhasesDispatchesToCommand(t *testing.T) {
	old := gatePhasesCommand
	t.Cleanup(func() { gatePhasesCommand = old })
	var gotArgs []string
	gatePhasesCommand = func(args []string, stdout, stderr io.Writer) int {
		gotArgs = append([]string(nil), args...)
		return 37
	}

	stdout := tempFile(t)
	stderr := tempFile(t)
	rc := (Command{Stdout: stdout, Stderr: stderr}).Run([]string{"gate-phases", "/tmp/root"})

	if rc != 37 {
		t.Fatalf("run gate-phases exit = %d, want injected exit 37", rc)
	}
	if want := []string{"/tmp/root"}; !reflect.DeepEqual(gotArgs, want) {
		t.Fatalf("gate-phases args = %#v, want %#v", gotArgs, want)
	}
}

func TestShellWrapperRoutesGatePhasesWithNonRootKit(t *testing.T) {
	root := t.TempDir()
	kit := filepath.Join(root, "kit")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	argvFile := filepath.Join(root, "argv")
	writeExecutable(t, filepath.Join(kit, "dist", "bench"), `#!/usr/bin/env bash
printf '%s\n' "$BENCH_KIT" "$@" > "$BENCH_TEST_ARGV"
`)

	cmd := exec.Command("bash", filepath.Join(kit, "bin", "bench.sh"), "gate-phases", "/tmp/repo root")
	cmd.Env = append(capability.WithoutEnvironment(os.Environ(), runbinary.Env),
		"BENCH_TEST_ARGV="+argvFile, "BENCH_HOME="+filepath.Join(root, "home"),
		"BENCH_KIT="+kit, "BENCH_WRAPPER=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bench.sh gate-phases failed: %v\n%s", err, out)
	}
	got := strings.Split(strings.TrimSpace(readPath(t, argvFile)), "\n")
	want := []string{kit, "gate-phases", "/tmp/repo root"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("wrapper routed argv = %#v, want %#v\noutput:\n%s", got, want, out)
	}
}

// repairUsageLine is the grammar the wrapper-only repair verb answers with. It is
// written here independently of bin/bench.sh, so a reworded or deleted help arm reds this
// test rather than passing against a re-derived expectation.
const repairUsageLine = "usage: bench repair [--prune]\n"

// TestShellWrapperRepairAnswersItsOwnHelp holds repair to the help contract every
// Go-implemented verb keeps: `--help` prints the grammar on stdout at exit 0. The
// unknown-argument leg runs beside it, because the help arm sits in front of the
// usage-error arm and must not take that arm's cases with it. Neither leg resolves a
// binary, so the case needs no dist/bench.
func TestShellWrapperRepairAnswersItsOwnHelp(t *testing.T) {
	root := t.TempDir()
	kit := filepath.Join(root, "kit")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))

	run := func(args ...string) (string, string, int) {
		t.Helper()
		cmd := exec.Command("bash", append([]string{filepath.Join(kit, "bin", "bench.sh"), "repair"}, args...)...)
		env := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
		for _, name := range []string{"BENCH_KIT", "BENCH_WRAPPER"} {
			env = capability.WithoutEnvironment(env, name)
		}
		cmd.Env = append(env, "BENCH_HOME="+filepath.Join(root, "home"), "BENCH_KIT="+kit)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		code := 0
		if err := cmd.Run(); err != nil {
			exit, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("bench.sh repair %v: %v", args, err)
			}
			code = exit.ExitCode()
		}
		return stdout.String(), stderr.String(), code
	}

	for _, flag := range []string{"--help", "-h", "help"} {
		stdout, stderr, code := run(flag)
		if code != 0 || stdout != repairUsageLine || stderr != "" {
			t.Errorf("repair %s = (exit %d, stdout %q, stderr %q), want exit 0 and %q on stdout", flag, code, stdout, stderr, repairUsageLine)
		}
	}

	stdout, stderr, code := run("--bogus")
	if code != 2 || stdout != "" || stderr != repairUsageLine {
		t.Fatalf("repair --bogus = (exit %d, stdout %q, stderr %q), want exit 2 and %q on stderr", code, stdout, stderr, repairUsageLine)
	}
}

func TestShellWrapperOfflineMissingBinarySuppressesRepairBeforeStartingNode(t *testing.T) {
	root := t.TempDir()
	kit := filepath.Join(root, "kit")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	marker := filepath.Join(root, "node-started")
	nodeDir := filepath.Join(root, "bin")
	writeExecutable(t, filepath.Join(nodeDir, "node"), "#!/usr/bin/env bash\ntouch \"$BENCH_TEST_NODE_MARKER\"\n")

	cmd := exec.Command("bash", filepath.Join(kit, "bin", "bench.sh"), "models")
	env := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	for _, name := range []string{"BENCH_KIT", "BENCH_WRAPPER", "BENCH_REPAIR", "BENCH_OFFLINE"} {
		env = capability.WithoutEnvironment(env, name)
	}
	cmd.Env = append(env,
		"BENCH_HOME="+filepath.Join(root, "home"), "BENCH_KIT="+kit, "BENCH_OFFLINE=1",
		"BENCH_TEST_NODE_MARKER="+marker, "PATH="+nodeDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	out, err := cmd.CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	_, markerErr := os.Stat(marker)
	markerStarted := markerErr == nil
	if !ok || exit.ExitCode() != 127 || !strings.Contains(string(out), "bench: repair suppressed by BENCH_OFFLINE=1") || markerStarted {
		t.Fatalf("offline missing binary = (exit %v, output %q, node started %t), want exit 127, offline suppression, and no node start", err, out, markerStarted)
	}
}

// followOnEnvelope builds a raw-call envelope whose command names a path under the
// given root's worktree pool, keyed to a fixed owner and assignment id.
func followOnEnvelope(home, root string) (string, string) {
	owner := strings.Repeat("a", 32)
	id := strings.Repeat("b", 32)
	segment := poolkey.AssignmentSegment(owner, id)
	path := filepath.Join(poolkey.Pool(home, root), segment, "x")
	return `{"tool_input":{"command":"sed -n 1p ` + path + `"}}`, id
}

// TestGuardBenchFollowOnRecordsRawCall covers EC09: a raw-call envelope whose command
// names a path under the pool makes the verb exit 0 and write one census record.
// The verdict must change if this record is skipped, since a record's own error can
// never surface here.
func TestGuardBenchFollowOnRecordsRawCall(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	home := filepath.Join(t.TempDir(), "home")
	t.Setenv("BENCH_HOME", home)
	envelope, id := followOnEnvelope(home, root)
	var errb bytes.Buffer
	if code := guardBenchFollowOn(nil, strings.NewReader(envelope), io.Discard, &errb); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if errb.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", errb.String())
	}
	if _, err := os.Stat(filepath.Join(census.Dir(home, root), id)); err != nil {
		t.Fatalf("census record not written: %v", err)
	}
}

// TestGuardBenchFollowOnUnwritableHomeStaysSilent covers EC10: a home the record
// cannot write under still exits 0 with empty stderr, because a record error never
// surfaces as a diagnostic or a block.
func TestGuardBenchFollowOnUnwritableHomeStaysSilent(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	home := filepath.Join(t.TempDir(), "home-is-a-file")
	if err := os.WriteFile(home, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_HOME", home)
	envelope, _ := followOnEnvelope(home, root)
	var errb bytes.Buffer
	if code := guardBenchFollowOn(nil, strings.NewReader(envelope), io.Discard, &errb); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if errb.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", errb.String())
	}
}

// TestGuardBenchFollowOnNoCommandFieldRecordsNothing keeps the warning path: an
// envelope with no command field records nothing, because the decode fails before
// the record call.
func TestGuardBenchFollowOnNoCommandFieldRecordsNothing(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	home := filepath.Join(t.TempDir(), "home")
	t.Setenv("BENCH_HOME", home)
	var errb bytes.Buffer
	if code := guardBenchFollowOn(nil, strings.NewReader(`{"tool_input":{}}`), io.Discard, &errb); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(errb.String(), "WARNING") {
		t.Fatalf("stderr = %q, want WARNING", errb.String())
	}
	if _, err := os.Stat(filepath.Join(home, "census")); !os.IsNotExist(err) {
		t.Fatalf("census dir should not exist, stat err = %v", err)
	}
}

// TestGuardBenchFollowOnWithoutPoolPrefixResolvesNoRoot proves the substring test runs
// before any root resolution: outside a git repository, a plain command still exits 0,
// so the verb never called git.Root for it.
func TestGuardBenchFollowOnWithoutPoolPrefixResolvesNoRoot(t *testing.T) {
	dir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	home := filepath.Join(t.TempDir(), "home")
	t.Setenv("BENCH_HOME", home)
	var errb bytes.Buffer
	envelope := `{"tool_input":{"command":"echo hi"}}`
	if code := guardBenchFollowOn(nil, strings.NewReader(envelope), io.Discard, &errb); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if _, err := os.Stat(filepath.Join(home, "census")); !os.IsNotExist(err) {
		t.Fatalf("census dir should not exist, stat err = %v", err)
	}
}

func TestGuardGitBlockAllow(t *testing.T) {
	var errb bytes.Buffer
	block := `{"tool_input":{"command":"git push"}}`
	if code := guardGit(nil, strings.NewReader(block), io.Discard, &errb); code != 2 {
		t.Errorf("block exit = %d, want 2", code)
	}
	if !strings.Contains(errb.String(), "BLOCKED:") {
		t.Errorf("block did not emit BLOCKED on stderr: %q", errb.String())
	}
	for _, in := range []string{`{"tool_input":{"command":"git status"}}`, "not json", `{"tool_input":{"command":""}}`} {
		if code := guardGit(nil, strings.NewReader(in), io.Discard, io.Discard); code != 0 {
			t.Errorf("allow exit for %q = %d, want 0", in, code)
		}
	}
}

func TestCaptureClaudeAgentIntentReplayIsByteIdempotent(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", root).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	envelope := []byte(`{"tool_name":"Agent","tool_use_id":"same-id","tool_input":{"description":"same objective"}}`)
	captureClaudeAgentIntent(envelope, io.Discard)
	path := filepath.Join(root, ".git", "bench-intent.json")
	before := readPath(t, path)
	captureClaudeAgentIntent(envelope, io.Discard)
	after := readPath(t, path)
	if after != before {
		t.Fatalf("replayed Claude capture changed ledger bytes:\nbefore=%s\nafter=%s", before, after)
	}
}

// panicReader forces guardGit's stdin read to panic, exercising the recover-to-exit-3
// rim, so a crash can never masquerade as an exit-2 block.
type panicReader struct{}

func (panicReader) Read([]byte) (int, error) { panic("boom") }

func TestGuardGitRecoversToExit3(t *testing.T) {
	if code := guardGit(nil, panicReader{}, io.Discard, io.Discard); code != 3 {
		t.Errorf("panic mapped to exit %d, want 3", code)
	}
}

func tempFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "out-*")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func readFile(t *testing.T, f *os.File) string {
	t.Helper()
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func readPath(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func copyExecutable(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	writeExecutable(t, dst, string(data))
}

func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// gitInitRepo makes <dir> a git repository with one commit, so `git worktree add` has a
// commit to branch from.
func gitInitRepo(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"config", "user.email", "probe@example.test"},
		{"config", "user.name", "probe"},
		{"add", "-A"},
		{"commit", "-qm", "seed", "--allow-empty"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v in %s: %v\n%s", args, dir, err, out)
		}
	}
}

// gitlessPATH builds a PATH holding the externals bin/bench.sh needs, with no `git`, so
// the wrapper's same-repository probe hits its no-git rim rather than an unrelated failure.
func gitlessPATH(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, tool := range []string{"bash", "env", "uname", "dirname", "basename", "readlink", "tr"} {
		src, err := exec.LookPath(tool)
		if err != nil {
			capability.Capability(t, capability.Tool, tool+" not on PATH")
		}
		if err := os.Symlink(src, filepath.Join(dir, tool)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := exec.LookPath("git"); err != nil {
		capability.Capability(t, capability.Tool, "git not on PATH")
	}
	return dir
}

// The wrapper resolves the kit for this invocation. A CWD inside a worktree of the
// wrapper's own repository moves the kit to that worktree. Every other CWD, a different
// repository (the linked-project case), no repository at all, or an environment without
// git, keeps the wrapper's own tree. The same-repository test compares the git common
// directory's identity, so a linked project repo, which a path-prefix test could not
// tell from a worktree, still gets the wrapper's tree.
func TestShellWrapperResolvesKitForTheInvocation(t *testing.T) {
	base := t.TempDir()
	kit := filepath.Join(base, "kit")
	argvFile := filepath.Join(base, "argv")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(kit, "bin", "bench.sh"))
	gitInitRepo(t, kit)
	// dist/bench stays untracked, exactly as in a real checkout. The worktree case reaches
	// it only because main_tree_kit re-anchors the binary lookup at the main tree.
	writeExecutable(t, filepath.Join(kit, "dist", "bench"), `#!/usr/bin/env bash
printf '%s\n' "$BENCH_KIT" > "$BENCH_TEST_ARGV"
`)
	worktree := filepath.Join(base, "wt")
	add := exec.Command("git", "worktree", "add", "-q", "-b", "probe", worktree)
	add.Dir = kit
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("git worktree add: %v\n%s", err, out)
	}
	project := filepath.Join(base, "project")
	gitInitRepo(t, project)
	// In the adopted-repo layout, the kit is a subdirectory of the project repo. The
	// project's own CWD shares its common dir with the wrapper's tree, so only the
	// "kit is its tree's top level" condition keeps the wrapper naming <repo>/.bench.
	adopted := filepath.Join(base, "adopted")
	adoptedKit := filepath.Join(adopted, ".bench")
	copyExecutable(t, filepath.Join("..", "..", "bin", "bench.sh"), filepath.Join(adoptedKit, "bin", "bench.sh"))
	writeExecutable(t, filepath.Join(adoptedKit, "dist", "bench"), `#!/usr/bin/env bash
printf '%s\n' "$BENCH_KIT" > "$BENCH_TEST_ARGV"
`)
	gitInitRepo(t, adopted)
	plain := filepath.Join(base, "plain")
	if err := os.MkdirAll(plain, 0o755); err != nil {
		t.Fatal(err)
	}

	wrapper := filepath.Join(kit, "bin", "bench.sh")
	cleanEnv := capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	cleanEnv = capability.WithoutEnvironment(capability.WithoutEnvironment(cleanEnv, "BENCH_KIT"), "BENCH_WRAPPER")
	cleanEnv = append(cleanEnv, "BENCH_TEST_ARGV="+argvFile, "BENCH_HOME="+filepath.Join(base, "home"))

	for _, tc := range []struct {
		name    string
		wrapper string
		cwd     string
		env     []string
		want    string
	}{
		{name: "worktree of the wrapper's repository", cwd: worktree, want: worktree},
		{name: "an adopted repo's own kit subdirectory", wrapper: filepath.Join(adoptedKit, "bin", "bench.sh"), cwd: adopted, want: adoptedKit},
		{name: "the wrapper's own main checkout", cwd: kit, want: kit},
		{name: "a different repository (linked project)", cwd: project, want: kit},
		{name: "outside any repository", cwd: plain, want: kit},
		{name: "no git on PATH", cwd: worktree, env: []string{"PATH=" + gitlessPATH(t)}, want: kit},
		{name: "an explicit BENCH_KIT wins", cwd: worktree, env: []string{"BENCH_KIT=" + kit}, want: kit},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Remove(argvFile); err != nil && !os.IsNotExist(err) {
				t.Fatal(err)
			}
			invoked := wrapper
			if tc.wrapper != "" {
				invoked = tc.wrapper
			}
			cmd := exec.Command("bash", invoked, "gate-phases", "/tmp/root")
			cmd.Dir = tc.cwd
			cmd.Env = append(append([]string(nil), cleanEnv...), tc.env...)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("bench.sh from %s: %v\n%s", tc.cwd, err, out)
			}
			if got := strings.TrimSpace(readPath(t, argvFile)); got != tc.want {
				t.Fatalf("BENCH_KIT = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestRoadmapFlagsRouteToTheirOwners covers RF10. The flow flag adds a surface and moves
// none. The bare form is byte-compared against its owner, and the context form is pinned
// to the schema `/bench-drain` accepts.
func TestRoadmapFlagsRouteToTheirOwners(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	const heading = "**FT1 — fixture.**"
	roadmaptest.WriteSplitBoard(t, root, heading+"\n", map[string]string{"FT1.md": heading + "\n"})
	t.Chdir(root)

	wantBare, wantCode := roadmap.RoadmapCommand(nil)
	gotBare, gotCode := roadmapCommand(nil)
	if gotBare != wantBare || gotCode != wantCode {
		t.Fatalf("bare roadmap = %q, %d; want its owner's %q, %d", gotBare, gotCode, wantBare, wantCode)
	}

	context, code := roadmapCommand([]string{"--context"})
	if code != 0 || !strings.Contains(context, "context[1]{schema,full,sequence_trusted}:\n  4,") {
		t.Fatalf("roadmap --context = %q, %d; want the schema-4 snapshot on exit 0", context, code)
	}

	flow, code := roadmapCommand([]string{"--flow"})
	if code != 0 || !strings.HasPrefix(flow, "flow[") {
		t.Fatalf("roadmap --flow = %q, %d; want the flow block on exit 0", flow, code)
	}
}

// TestCommandsBriefGradesTheRunningExecutable covers BF1 and BF2. The probe is the one
// verb a session runs to learn whether Bench answers at all, so a stale binary answering
// it is a false pass nothing else in the loop catches. The rows run out of process,
// because the verdict is about the executable that is running, and a test binary is not
// the checkout's publication. Both rows read one real subject-mode build, so the
// published pair under grading is the one the builder writes.
func TestCommandsBriefGradesTheRunningExecutable(t *testing.T) {
	root := commandsBriefCheckout(t)
	published := commandsBriefPublication(t, root)

	t.Run("a matching seal keeps the three-line answer", func(t *testing.T) {
		binary := commandsBriefCopy(t, root, published)
		out, code := commandsBriefProbe(t, root, binary)
		if code != 0 {
			t.Fatalf("commands --brief under a matching seal = exit %d:\n%s", code, out)
		}
		if want, _ := commandsCommand([]string{"--brief"}); out != want {
			t.Fatalf("commands --brief = %q, want the probe answer %q", out, want)
		}
	})

	t.Run("a mismatched seal refuses with the rebuild action", func(t *testing.T) {
		binary := commandsBriefCopy(t, root, published)
		commandsBriefBreakSeal(t, binary)
		out, code := commandsBriefProbe(t, root, binary)
		if code == 0 {
			t.Fatalf("commands --brief under a mismatched seal exited 0:\n%s", out)
		}
		if action := freshness.RebuildAction(root); !strings.Contains(out, action) {
			t.Fatalf("commands --brief refusal = %q, want the rebuild action %q", out, action)
		}
	})
}

// commandsBriefCheckout names the checkout exactly as the probe's own handler will,
// because the rebuild action the refusal prints is built from that spelling.
func commandsBriefCheckout(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := git.RootAt(root)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

// commandsBriefPublication builds the checkout's own binary into the checkout, which is
// where the wrapper's binary route resolves it. A build outside the tree is not this
// checkout's publication, so only an in-tree one exercises the grading branch.
func commandsBriefPublication(t *testing.T, root string) string {
	t.Helper()
	dist := filepath.Join(root, "dist")
	if err := os.MkdirAll(dist, 0o755); err != nil {
		t.Fatal(err)
	}
	dir, err := os.MkdirTemp(dist, "commands-brief-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	out := filepath.Join(dir, "bench")
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "go-build.sh"), root, out)
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("subject build: %v\n%s", err, output)
	}
	return out
}

// commandsBriefCopy gives one row its own published pair, so neither row depends on the
// order the other ran in.
func commandsBriefCopy(t *testing.T, root, published string) string {
	t.Helper()
	dir, err := os.MkdirTemp(filepath.Join(root, "dist"), "commands-brief-row-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	binary := filepath.Join(dir, "bench")
	for _, pair := range [][2]string{{published, binary}, {published + ".seal", binary + ".seal"}} {
		body, err := os.ReadFile(pair[0])
		if err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0o644)
		if pair[0] == published {
			mode = 0o755
		}
		if err := os.WriteFile(pair[1], body, mode); err != nil {
			t.Fatal(err)
		}
	}
	return binary
}

// commandsBriefBreakSeal rewrites the source digest the seal records, which is the exact
// state a rebuilt tree beside an unrebuilt binary leaves behind.
func commandsBriefBreakSeal(t *testing.T, binary string) {
	t.Helper()
	body, err := os.ReadFile(binary + ".seal")
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(body, &fields); err != nil {
		t.Fatal(err)
	}
	sources, ok := fields["sources"].(string)
	if !ok || sources == "" {
		t.Fatalf("seal records no source digest: %s", body)
	}
	fields["sources"] = strings.Repeat("0", len(sources))
	edited, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary+".seal", edited, 0o644); err != nil {
		t.Fatal(err)
	}
}

func commandsBriefProbe(t *testing.T, root, binary string) (string, int) {
	t.Helper()
	cmd := exec.Command(binary, "commands", "--brief")
	cmd.Dir = root
	cmd.Env = capability.WithoutEnvironment(os.Environ(), runbinary.Env)
	output, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		exit := &exec.ExitError{}
		if !errors.As(err, &exit) {
			t.Fatalf("run %s commands --brief: %v\n%s", binary, err, output)
		}
		code = exit.ExitCode()
	}
	return string(output), code
}
