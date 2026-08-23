package adopt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

// Setup is the `bench setup` entry point. It composes the existing init, link, and
// transaction seams into one inspect -> preview -> confirm -> write flow. It adds no
// second transaction and no second asset installer.
func Setup(args []string, stdin io.Reader, stdout, stderr io.Writer, version string) int {
	return setup(args, stdin, stdout, stderr, version, terminal.IsTerminal)
}

func setup(args []string, stdin io.Reader, stdout, stderr io.Writer, version string, isTerminal func(io.Reader) bool) int {
	planOnly, autoYes, ok := parseSetupArgs(args)
	if !ok {
		fmt.Fprintln(stderr, "usage: bench setup [--plan|--yes]")
		return 2
	}

	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "bench setup: not a git repository - run 'git init' first, then re-run bench setup")
		return 1
	}
	kit := kitDir()
	facts := inspectRepo(root, kit)
	fmt.Fprint(stdout, renderSetupPreview(root, facts))

	if planOnly {
		return 0
	}

	tty := isTerminal(stdin)
	if !autoYes && !tty {
		fmt.Fprintln(stderr, "bench setup: non-interactive stdin with no confirmation flag - pass --yes to auto-confirm an ambiguity-free plan, or --plan to preview only and write nothing")
		return 1
	}

	// One shared bufio.Reader carries any interactive I/O. A bufio.Reader may buffer more of
	// stdin than one ReadString('\n') consumes. A second bufio.Reader created later for the
	// confirm prompt would strand an already-typed answer in the first reader's internal
	// buffer.
	reader := bufio.NewReader(stdin)

	if len(facts.openQuestions) > 0 {
		if autoYes {
			for _, q := range facts.openQuestions {
				fmt.Fprintln(stderr, "bench setup: open question - "+q)
			}
			fmt.Fprintln(stderr, "bench setup: --yes refuses an ambiguous plan; resolve the open question(s) above, then re-run")
			return 1
		}
		// Interactive one-at-a-time question sequencing asks each open question in order over
		// the shared reader. It then falls through to the same single confirm an ambiguity-free
		// plan uses.
		resolved, err := askSetupQuestions(facts, reader, stdout)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		facts = resolved
	}

	if !autoYes {
		fmt.Fprint(stdout, "Proceed with this plan? [y/N] ")
		line, readErr := reader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			fmt.Fprintln(stderr, "bench setup: could not read confirmation")
			return 1
		}
		if !isAffirmative(line) {
			fmt.Fprintln(stderr, "bench setup: declined; no writes")
			return 1
		}
	}

	return convergeSetup(root, kit, version, facts, stdout, stderr)
}

func parseSetupArgs(args []string) (planOnly, autoYes, ok bool) {
	for _, a := range args {
		switch a {
		case "--plan":
			planOnly = true
		case "--yes":
			autoYes = true
		default:
			return false, false, false
		}
	}
	if planOnly && autoYes {
		return false, false, false
	}
	return planOnly, autoYes, true
}

func isAffirmative(line string) bool {
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true
	}
	return false
}

// setupFacts is every inspected, preview-observable fact this slice acts on. Facts
// beyond the gate table (harness-file matrix, foreign-gate detail, etc.) are the later
// stories' territory; this stays the honest subset stories 1/2/6/7/9 need.
type setupFacts struct {
	agentsExists     bool
	claudeExists     bool
	gateExists       bool
	gateInputsExists bool
	profileExists    bool
	profileName      string
	gateCandidates   []gateCandidate
	gateCommand      string // resolved command when exactly one candidate matched
	zeroSignal       bool
	openQuestions    []string
	ignoredSignals   []string // preview-only fact lines for a candidate a read/parse error excluded
}

type gateCandidate struct {
	source, command string
}

// label is the one "<command> (<source>)" rendering every candidate listing composes
// from. The open-question text in inspectRepo and the interactive prompt in
// setup_prompt.go both call this instead of each formatting the pair itself.
func (c gateCandidate) label() string {
	return fmt.Sprintf("%s (%s)", c.command, c.source)
}

func inspectRepo(root, kit string) setupFacts {
	name := filepath.Base(root)
	profileRel := filepath.Join("projects", name+".md")
	facts := setupFacts{
		agentsExists:     fileExists(filepath.Join(root, "AGENTS.md")),
		claudeExists:     fileExists(filepath.Join(root, "CLAUDE.md")),
		gateExists:       fileExists(filepath.Join(root, ".bench", "gate.sh")),
		gateInputsExists: fileExists(filepath.Join(root, ".bench", "gate-inputs.json")),
		profileExists:    fileExists(filepath.Join(root, profileRel)),
		profileName:      name,
	}
	facts.gateCandidates, facts.ignoredSignals = detectGateCandidates(root)
	switch len(facts.gateCandidates) {
	case 0:
		facts.zeroSignal = true
	case 1:
		facts.gateCommand = facts.gateCandidates[0].command
	default:
		labels := make([]string, len(facts.gateCandidates))
		for i, c := range facts.gateCandidates {
			labels[i] = c.label()
		}
		facts.openQuestions = append(facts.openQuestions,
			"multiple test commands detected - "+strings.Join(labels, ", ")+" - which one should .bench/gate.sh run?")
	}
	return facts
}

// detectGateCandidates is the one source for the gate-inference table. It drives both the
// plan preview and the written gate.sh, so detection and proposal never drift apart. The
// second return value carries one preview-only fact line per candidate a read/parse error
// excluded. An unreadable or malformed package.json/Makefile is never just dropped from
// the table with no trace: nothing is acted on silently.
func detectGateCandidates(root string) ([]gateCandidate, []string) {
	var out []gateCandidate
	var ignored []string
	if fileExists(filepath.Join(root, "go.mod")) {
		out = append(out, gateCandidate{"go.mod", "go test ./..."})
	}
	if ok, note := packageJSONTestScript(root); note != "" {
		ignored = append(ignored, note)
	} else if ok {
		out = append(out, gateCandidate{"package.json test script", "npm test"})
	}
	if fileExists(filepath.Join(root, "Cargo.toml")) {
		out = append(out, gateCandidate{"Cargo.toml", "cargo test"})
	}
	if ok, note := makefileTestTarget(root); note != "" {
		ignored = append(ignored, note)
	} else if ok {
		out = append(out, gateCandidate{"Makefile test target", "make test"})
	}
	return out, ignored
}

// packageJSONTestScript reports whether package.json declares a non-empty "test" script.
// A non-empty note means the file exists but could not be inspected, because it is
// unreadable or malformed JSON. The candidate is excluded either way, but the note
// carries why, so the preview names the exclusion instead of silently dropping the
// signal.
func packageJSONTestScript(root string) (bool, string) {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if os.IsNotExist(err) {
		return false, ""
	}
	if err != nil {
		return false, "package.json unreadable -> ignored as a gate signal"
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false, "package.json malformed JSON -> ignored as a gate signal"
	}
	return strings.TrimSpace(pkg.Scripts["test"]) != "", ""
}

// makefileTestTarget reports whether the Makefile declares a "test:" target, with the
// same exists-but-unreadable note convention as packageJSONTestScript.
func makefileTestTarget(root string) (bool, string) {
	data, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if os.IsNotExist(err) {
		return false, ""
	}
	if err != nil {
		return false, "Makefile unreadable -> ignored as a gate signal"
	}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimRight(line, "\r")
		if trimmed == "test:" || strings.HasPrefix(trimmed, "test:") || strings.HasPrefix(trimmed, "test :") {
			return true, ""
		}
	}
	return false, ""
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// renderSetupPreview lists every inferred fact with its consequence. Nothing here is
// acted on until confirm/--yes; a run that stops at --plan leaves the tree untouched.
func renderSetupPreview(root string, facts setupFacts) string {
	var b strings.Builder
	fmt.Fprintf(&b, "bench setup: plan for %s\n", root)
	fmt.Fprintln(&b, "  git repository detected -> proceeding")

	if facts.agentsExists {
		fmt.Fprintln(&b, "  AGENTS.md exists -> the managed Bench block will be added or updated, project content preserved")
	} else {
		fmt.Fprintln(&b, "  AGENTS.md absent -> will be created with the managed Bench block")
	}
	if facts.claudeExists {
		fmt.Fprintln(&b, "  CLAUDE.md exists -> the managed import lines will be added or updated, project content preserved")
	} else {
		fmt.Fprintln(&b, "  CLAUDE.md absent -> will be created with the managed @AGENTS.md / @.bench/BENCH.md import lines")
	}

	switch {
	case facts.zeroSignal:
		fmt.Fprintln(&b, "  .bench/gate.sh: no build system detected -> a fail-closed stub will be written (red until configured)")
	case len(facts.openQuestions) > 0:
		for _, q := range facts.openQuestions {
			fmt.Fprintln(&b, "  .bench/gate.sh: "+q)
		}
	default:
		fmt.Fprintf(&b, "  .bench/gate.sh: %s detected -> proposing `%s`\n", facts.gateCandidates[0].source, facts.gateCommand)
	}
	for _, note := range facts.ignoredSignals {
		fmt.Fprintln(&b, "  "+note)
	}
	if facts.gateExists {
		fmt.Fprintln(&b, "  .bench/gate.sh already exists -> converged if bench-managed, preserved as a conflict otherwise")
	}

	if facts.gateInputsExists {
		fmt.Fprintln(&b, "  .bench/gate-inputs.json exists -> left as-is (reviewer-owned content)")
	} else {
		fmt.Fprintln(&b, "  .bench/gate-inputs.json absent -> will be seeded declaring BENCH_HOME and HOME")
	}

	if facts.profileExists {
		fmt.Fprintf(&b, "  projects/%s.md exists -> left as-is (reviewer-owned content)\n", facts.profileName)
	} else {
		fmt.Fprintf(&b, "  projects/%s.md absent -> will be scaffolded as a starting profile\n", facts.profileName)
	}

	fmt.Fprintln(&b, "  platform assets (.bench/, .agents/, .claude/, .codex/) -> installed or converged via the link lifecycle")
	return b.String()
}

// convergeSetup writes the confirmed plan. It composes buildLinkPlan (link's asset list)
// plus one inline gate.sh entry plus one seed-if-absent profile entry, all staged and
// promoted by the single transaction. There is no second transaction, no second asset
// installer, and no write that lands outside it. The profile uses the "seed" planEntry
// kind rather than "inline". A reviewer-authored profile already at that path is left
// untouched and never recorded in the manifest. A later hand-edit can never read back as
// a modified-managed conflict.
func convergeSetup(root, kit, version string, facts setupFacts, stdout, stderr io.Writer) int {
	plan, err := buildLinkPlan(kit)
	if err != nil {
		fmt.Fprintln(stderr, toon.Errorf("consumer payload allowlist rejected", err.Error()))
		return 1
	}
	plan = append(plan, planEntry{rel: ".bench/gate.sh", kind: "inline-exec", content: setupGateScript(facts)})
	profileRel := filepath.Join("projects", facts.profileName+".md")
	plan = append(plan, planEntry{rel: filepath.ToSlash(profileRel), kind: "seed", content: scaffoldProfile(facts)})
	plan = append(plan, planEntry{rel: ".bench/gate-inputs.json", kind: "seed", content: scaffoldGateInputs()})

	result, changed := transactionalLink(root, kit, "copy", version, plan, stdout, stderr)
	if result != 0 && result != 3 {
		return result
	}

	if facts.zeroSignal {
		fmt.Fprintln(stdout, "bench setup: zero-signal repository - .bench/gate.sh is a fail-closed stub; configure it before this repo can go green")
	} else if result == 0 {
		if changed {
			fmt.Fprintf(stdout, "bench setup: converged %s\n", root)
		} else {
			fmt.Fprintf(stdout, "bench setup: already converged %s\n", root)
		}
	}
	return finishSetup(stdout, facts, result == 3)
}

// setupGateScript embeds the proposed command for a detected ecosystem, or extends init's
// fail-closed scaffoldGate() shape when nothing was detected. A fabricated green gate is
// never written. Both branches compose the one gateScriptPreamble (init.go) rather than
// each re-authoring the shebang/set/git-root-guard lines.
func setupGateScript(facts setupFacts) string {
	if facts.zeroSignal {
		return scaffoldGate()
	}
	return gateScriptPreamble("# bench setup wrote this gate from the gate-inference table. Run /bench-setup-repo\n# to refine it.\n") + facts.gateCommand + "\n"
}

// scaffoldGateInputs is the seeded gate input manifest every adopted repository starts
// from. The gate launches its script with PATH plus only the names declared here. The
// installed wrapper's first statement reads HOME, to derive BENCH_HOME, under set -u. A
// repository with no manifest cannot run its own gate at all. The declared tools are the
// ones that wrapper and the scaffolded gate invoke. paths stays empty, because the seeded
// gate closes over no tracked file beyond the script the gate already reads itself.
func scaffoldGateInputs() string {
	return `{
  "schema": 1,
  "closure": "local",
  "environment": ["BENCH_HOME", "HOME"],
  "paths": [],
  "tools": ["bash", "basename", "dirname", "git", "readlink", "uname"]
}
`
}

func scaffoldProfile(facts setupFacts) string {
	gateNote := "no build system was detected; .bench/gate.sh is a fail-closed stub until configured."
	if !facts.zeroSignal && len(facts.openQuestions) == 0 {
		gateNote = fmt.Sprintf("`.bench/gate.sh` runs `%s`, proposed from %s.", facts.gateCommand, facts.gateCandidates[0].source)
	}
	return fmt.Sprintf(`# Project: %s

Scaffolded by `+"`bench setup`"+`. Fill in the judgment content below through the
`+"`/bench-setup-repo`"+` conversation - seams worth testing, the lines (model +
effort routing), and notes for a cold session.

## Gate (`+"`.bench/gate.sh`"+`)

    .bench/gate.sh

%s

## Lines (model + effort routing)

Not yet set - `+"`/bench-setup-repo`"+` binds tiers and cached routings for this
project.

## Notes for cold sessions

Read `+"`AGENTS.md`"+` first, then `+"`.bench/BENCH.md`"+` for the shared platform
rules.
`, facts.profileName, gateNote)
}
