package main

import (
	"bytes"
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
  bench retro <slug> --body <markdown>  validate and create one primary-local implementation retrospective
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
  bench harnesses [<harness>]  the harness record as TOON (harness, provider, phase_form, hooks, delegation_guard); one name prints that harness's cells
  bench coverage <spec>      acceptance-coverage state and rows as TOON (--check to validate)
  bench preflight review|build <slug>  phase-entry checks that a spec's artifacts agree with the tree, one verdict row per check
  bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name>  run focused Go-test or named-check evidence as TOON; no gate verdict
  bench probe <file> (--swap <old> --with <new> | --omit <old>) (--package <expr> [--run <go-regex>] | --check <name>) [--full]  mutate one file once, run one focused test or check, restore the file, and report bit, silent, invalid, or restore-failed
  bench outline [path] [--full]  top-level directory symbol counts as TOON; a path or --full locates candidate seams (file:line), never the project's blessed seams
  bench consumers <qualified-symbol> [--full]  every resolved Go reference edge as TOON (file:line, via, enclosing); identifies edges, never blessed seams
  bench doctor [--fix]       report (and repair) the PATH shim under a node version manager
  bench repair [--prune]     explicitly install the pinned platform binary or prune stale cache entries
  bench gate [--fresh]       run the project gate (the oracle; --fresh ignores a reusable green)
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
