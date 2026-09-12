package main

import (
	"bytes"
	"encoding/json"
	"github.com/gibbonmi/bench/internal/assessment"
	"github.com/gibbonmi/bench/internal/poolkey"
	"os"
	"path/filepath"
	"testing"
)

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
  bench retro <slug> (--body <markdown> | --scaffold)  draft, or validate and create, one primary-local implementation retrospective
  bench roadmap              show the top 10 roadmap rows + drain state
  bench status               ambient dashboard: what needs attention + the next action
  bench handoff [--harness <name>] [--next <command>] [--state-file <path>]  print the cold-start pin block and rewrite capture/session-handoff.md
  bench commands --brief     print the direct, read-only command probe
  bench dashboard [--stdout] write a self-contained HTML snapshot of the board (--stdout emits it)
  bench canary [root]        validate fixture inventory
  bench anchors <path>       anchors pinning a repo-relative path as TOON (kind, section, needle, line)
  bench learnings            open journal entries as a TOON table (date, title)
  bench maps                 unresolved decision-map tickets as TOON (map, ticket, type, state)
  bench guards               every guard's deny surface as TOON (guard, boundary, denies)
  bench diff                 review base + changed files as TOON (--full appends log + diff body; --base freezes source)
  bench harnesses [<harness>] [--record <path> --format <source-id>]  the harness record as TOON; one name prints that harness's cells; both flags observe one named session record
  bench assessment list | show <run-id> | record --input <file> | compare --plan <file> --runs <id,...>  store and inspect local workflow cost and quality
  bench coverage <spec>      acceptance-coverage state and rows as TOON (--check to validate)
  bench preflight review|build <slug>  phase-entry checks that a spec's artifacts agree with the tree, one verdict row per check
  bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name>  run focused Go-test or named-check evidence as TOON; no gate verdict
  bench probe <file> (--swap <old> --with <new> | --omit <old>) (--package <expr> [--run <go-regex>] | --check <name>) [--full]  mutate one file once, run one focused test or check, restore the file, and report bit, silent, invalid, or restore-failed
  bench outline [path] [--full]  top-level directory symbol counts as TOON; a path or --full locates candidate seams (file:line), never the project's blessed seams
  bench consumers <qualified-symbol> [--full]  every resolved Go reference edge as TOON (file:line, via, enclosing); identifies edges, never blessed seams
  bench doctor [--fix]       report (and repair) the PATH shim under a node version manager
  bench repair [--prune]     explicitly install the pinned platform binary or prune stale cache entries
  bench gate [--fresh] [--checkpoint <spec-path> (--chunk <id> | --complete)]  run the project gate (the oracle; --fresh ignores a reusable green)
  bench prep-release         ship-tier rehearsal: artifacts, cross-compile, preflight verify, ship canary
  bench release-preflight --mode verify|publish [--profile public|bank] [--phase name]  run repository release authorization
  bench release prepare|submit|promote|rollback|status --version <v> [--profile public|bank] [--root dir] [--registry url] [--path first|staged] [--adapter npm|fixture] [--provenance] [--message text]  governed npm publication
  bench worktree shell [--refresh] [objective] create an owned worktree subshell and release it on exit
  bench worktree list        list assignments and registered worktrees as TOON
  bench worktree path <target>  print one active owned worktree's absolute path for the file tools
  bench worktree exec <target> [--env KEY=VALUE]... -- <command> [args...]  run a child directly in an active owned worktree
  bench worktree show <target> <rev>:<path>  print one blob from a revision of an active owned worktree
  bench worktree build <target>  build an active owned worktree's tree into its own dist/bench
  bench worktree exec <target> -- bench gate  run one active owned worktree's gate
  bench worktree reauthorize --assignment <id> --request <token> --base <commit> --source-tip <commit> <path>  replace one lost request token after identity proof
  bench worktree merge --from <commit|target> <target>  merge a default-branch commit or a sibling's tip into an owned worktree
  bench worktree reset (--to <commit> | --restore <ref>) <target> [--apply <fingerprint>]  plan or apply a recoverable reset or restore
  bench worktree --help      show exact list, path, exec, show, build, create, release, clean, reclaim, reauthorize, merge, and reset grammar
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

func assessmentEnvelopeCases() map[string]axiEnvelopeCase {
	return map[string]axiEnvelopeCase{
		"assessment compare": {route: []string{"assessment", "compare"}, successArgv: []string{"assessment", "compare", "--plan", "comparison.json", "--runs", "fixture"}, deepSuccessArgv: []string{"assessment", "compare", "--plan", "../../comparison.json", "--runs", "fixture"}, emptyArgv: []string{"assessment", "compare", "--plan", "comparison.json", "--runs", ""}, blocks: []string{"comparison", "runs", "outcomes", "roles", "usage", "costs", "conditions", "variation", "reasons", "help"}, successMarker: "runs[1]", emptyMarker: "runs[0]", usage: "usage: bench assessment", setupSuccess: setupAssessmentComparison, setupEmpty: setupAssessmentComparison},
		"assessment list":    {route: []string{"assessment", "list"}, successArgv: []string{"assessment", "list"}, emptyArgv: []string{"assessment", "list"}, blocks: []string{"runs", "help"}, successMarker: "runs[1]", emptyMarker: "runs[0]", usage: "usage: bench assessment", setupSuccess: setupAssessment, setupEmpty: func(t *testing.T, _ string) { t.Setenv("BENCH_HOME", t.TempDir()) }},
		"assessment show":    {route: []string{"assessment", "show"}, successArgv: []string{"assessment", "show", "fixture"}, blocks: []string{"attempts", "summary", "record", "help"}, successMarker: "record[", usage: "usage: bench assessment", setupSuccess: setupAssessment, recordBacked: true},
	}
}
func setupAssessment(t *testing.T, root string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("BENCH_HOME", home)
	r := assessment.Run{Version: 1, RunID: "fixture", RepoKey: poolkey.Key(root), Source: "synthetic", Condition: "bench", TaskID: "fixture", State: "running"}
	if err := (assessment.Store{Home: home, Root: root}).Record(r); err != nil {
		t.Fatal(err)
	}
}

func setupAssessmentComparison(t *testing.T, root string) {
	t.Helper()
	setupAssessment(t, root)
	zero := 0.0
	p := assessment.Plan{Version: 1, ID: "fixture-plan", Purpose: "descriptive", Approval: assessment.Reference{Producer: "synthetic", Native: "fixture approval"}, Budget: assessment.Budget{Amount: &zero, Currency: "USD"}, Tasks: []assessment.PlanTask{{ID: "fixture", Repetitions: 1}}, Variable: "capability", QualityTolerance: &assessment.Tolerance{MaxFailureRate: &zero}, Conditions: []assessment.Condition{{ID: "bench", Revision: "synthetic", Harness: "fixture-v1", Lines: map[string]assessment.Line{"implementation": {Model: "synthetic", Effort: "high"}}, Limits: map[string]string{"iterations": "1"}}}}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, "comparison.json"), data, 0600); err != nil {
		t.Fatal(err)
	}
}
