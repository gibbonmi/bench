package main

import (
	"bytes"
	"encoding/json"
	"github.com/gibbonmi/bench/internal/assessment"
	"github.com/gibbonmi/bench/internal/poolkey"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
  bench capture drain [show|commit|abort] [<drain-id>]  seal or finish one concurrent-safe capture drain
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
  bench harnesses [<harness> [--record <path> --format <source-id>]]  the harness record as TOON; one name prints that harness's cells; both flags observe one named session record
  bench assessment list | show <run-id> | record --input <file> | compare --plan <file> --runs <id,...>  store and inspect local workflow cost and quality
  bench coverage <spec>      acceptance-coverage state and rows as TOON (--check to validate)
  bench preflight review <slug> [--base <commit>] [--source-tip <commit>]  review-entry checks that a spec's artifacts agree with the tree, one verdict row per check
  bench preflight review <slug> --charge --base <commit> --source-tip <commit> [--max-store-bytes <n>]  prepare one immutable review evidence artifact and print its bounded orientation
  bench preflight build <slug> [--base <commit>] [--source-tip <commit>]  build-entry checks that a spec's artifacts agree with the tree, one verdict row per check
  bench preflight build <slug> --charge --ticket <basename> --base <commit> --source-tip <commit> [--max-store-bytes <n>]  prepare one immutable build evidence artifact and print its bounded orientation
  bench preflight build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>  propose one ticket's Writes: entries from the pinned source
  bench preflight evidence <id> [--cursor <cursor>]  print one bounded fragment of a prepared evidence artifact and its exact successor
  bench preflight evidence <id> --source <source-id> [--cursor <cursor>]  print one bounded fragment of one declared source stream and its exact successor
  bench preflight evidence <id> --verify  verify every stored page and source digest of a prepared evidence artifact
  bench preflight evidence <id> --check-current  bind a prepared evidence artifact to the current assignment and source pair
  bench repair-pilot activate | report [--full]  collect and report attributed repair evidence for an explicit local pilot
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
  bench worktree land --request <opaque-id> --base <commit> --source-tip <commit> [--spec <slug>] -m <message> <path>  compose, gate, and publish one owned worktree
  bench worktree land --resume <published-commit> --request <opaque-id> --base <commit> --source-tip <commit> [--spec <slug>] <path>  resume incomplete post-publication landing work
  bench worktree --help      show exact list, path, exec, show, build, create, release, clean, reclaim, reauthorize, merge, reset, and land grammar
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

// TestEvidenceHelpInventory is CE152. TestHelpInventoryIsComplete owns the independent
// root inventory; these cases prove that `bench preflight --help` advertises exactly the
// preflight forms root help advertises, and that neither advertises a later operation.
func TestEvidenceHelpInventory(t *testing.T) {
	var root bytes.Buffer
	if code := (Command{Stdout: &root}).Run([]string{"help"}); code != 0 {
		t.Fatalf("help exit = %d", code)
	}
	var preflightHelp bytes.Buffer
	if code := (Command{Stdout: &preflightHelp}).Run([]string{"preflight", "--help"}); code != 0 {
		t.Fatalf("preflight help exit = %d", code)
	}
	var rootForms, preflightForms []string
	for _, line := range strings.Split(root.String(), "\n") {
		if form, ok := strings.CutPrefix(line, "  bench preflight "); ok {
			command, _, _ := strings.Cut(form, "  ")
			rootForms = append(rootForms, command)
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(preflightHelp.String()), "\n") {
		form := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "usage:"))
		preflightForms = append(preflightForms, strings.TrimPrefix(form, "bench preflight "))
	}
	t.Run("same forms", func(t *testing.T) {
		if len(rootForms) == 0 || strings.Join(rootForms, "\n") != strings.Join(preflightForms, "\n") {
			t.Fatalf("root preflight forms:\n%s\npreflight help forms:\n%s", strings.Join(rootForms, "\n"), strings.Join(preflightForms, "\n"))
		}
	})
	t.Run("preparation and every implemented read", func(t *testing.T) {
		for _, want := range []string{"build <slug> --charge --ticket <basename> --base <commit> --source-tip <commit> [--max-store-bytes <n>]",
			"evidence <id> [--cursor <cursor>]", "evidence <id> --source <source-id> [--cursor <cursor>]",
			"evidence <id> --verify", "evidence <id> --check-current"} {
			if !strings.Contains(strings.Join(preflightForms, "\n"), want) {
				t.Errorf("preflight help omits %q", want)
			}
		}
	})
	t.Run("no later operation", func(t *testing.T) {
		for _, later := range []string{"evidence-clean"} {
			if strings.Contains(root.String(), later) || strings.Contains(preflightHelp.String(), later) {
				t.Errorf("help advertises the later operation %q", later)
			}
		}
	})
}

func TestRepairPilotRoute(t *testing.T) {
	t.Run("dispatch", func(t *testing.T) {
		root := newAXIEnvelopeRepo(t)
		t.Setenv("BENCH_HOME", t.TempDir())
		t.Setenv("BENCH_KIT", root)
		for _, row := range []struct {
			argv []string
			want string
		}{
			{argv: []string{"repair-pilot", "report"}, want: "inactive"},
			{argv: []string{"repair-pilot", "activate"}, want: "active"},
			{argv: []string{"assessment", "--help"}, want: "usage: bench assessment"},
		} {
			result := runAXICommandAt(t, root, row.argv)
			if result.code != 0 || !strings.Contains(result.stdout, row.want) {
				t.Fatalf("%v = stdout %q, stderr %q, exit %d; want %q at exit 0", row.argv, result.stdout, result.stderr, result.code, row.want)
			}
		}
		input := map[string]any{
			"version": 1,
			"observation": map[string]any{
				"id": "public-route", "observed_at": time.Now().UTC().Format(time.RFC3339Nano),
				"sequence":      map[string]string{"source": "public-source", "spec": "public-spec", "chunk": "public-chunk"},
				"assignment_id": "public-assignment", "session_id": "public-session", "source_revision": "public-revision",
				"stage": "pre-review", "kind": "failure", "failure_completeness": "complete",
				"failures": []map[string]any{{"check": "unit", "identity": "REQ-1", "diagnostic": "fixture blocker", "ownership": "diff-owned", "blocking": true,
					"reference": map[string]string{"producer": "public-route-test", "native": "native:public-failure"}}},
				"references": []map[string]string{{"producer": "public-route-test", "native": "native:public-observation"}},
			},
		}
		data, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		inputPath := filepath.Join(t.TempDir(), "record.json")
		if err := os.WriteFile(inputPath, append(data, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
		result := runAXICommandAt(t, root, []string{"repair-pilot", "record", "--input", inputPath})
		if result.code != 0 || !strings.Contains(result.stdout, "active") {
			t.Fatalf("public record route = stdout %q, stderr %q, exit %d", result.stdout, result.stderr, result.code)
		}
	})
	t.Run("worktrees", func(t *testing.T) {
		root := newAXIEnvelopeRepo(t)
		linked := filepath.Join(t.TempDir(), "linked")
		runAXIGit(t, "-C", root, "worktree", "add", "-q", "-b", "repair-pilot-linked", linked)
		home := t.TempDir()
		t.Setenv("BENCH_HOME", home)
		t.Setenv("BENCH_KIT", linked)
		activated := runAXICommandAt(t, linked, []string{"repair-pilot", "activate"})
		if activated.code != 0 || !strings.Contains(activated.stdout, "active") {
			t.Fatalf("worktree activation = stdout %q, stderr %q, exit %d", activated.stdout, activated.stderr, activated.code)
		}
		t.Setenv("BENCH_KIT", root)
		reported := runAXICommandAt(t, root, []string{"repair-pilot", "report"})
		if reported.code != 0 || !strings.Contains(reported.stdout, "active") {
			t.Fatalf("primary report = stdout %q, stderr %q, exit %d", reported.stdout, reported.stderr, reported.code)
		}
		pilot := filepath.Join(home, "repair-pilot", poolkey.Key(root), "pilot.json")
		if _, err := os.Stat(pilot); err != nil {
			t.Fatalf("canonical pilot document: %v", err)
		}
	})
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
