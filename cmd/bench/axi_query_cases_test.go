package main

import (
	"path/filepath"
	"testing"
)

func axiEnvelopeCases() map[string]axiEnvelopeCase {
	noSetup := func(*testing.T, string) {}
	cases := map[string]axiEnvelopeCase{
		"anchors": {
			route: []string{"anchors"}, successArgv: []string{"anchors", "AGENTS.md"}, deepSuccessArgv: []string{"anchors", "../../AGENTS.md"}, emptyArgv: []string{"anchors", "unregistered.md"},
			blocks:        []string{"anchors", "help"},
			successMarker: "anchors[", emptyMarker: "anchors[0]{kind,section,needle,line}:\n", usage: "usage: bench anchors", setupSuccess: setupAXIAnchors, setupEmpty: noSetup,
		},
		"learnings": {
			route: []string{"learnings"}, successArgv: []string{"learnings"}, emptyArgv: []string{"learnings"},
			blocks:        []string{"learnings", "help"},
			successMarker: "learnings[1]{date,title}:\n", emptyMarker: "learnings[0]{date,title}:\n", usage: "usage: bench learnings", setupSuccess: setupAXILearnings, setupEmpty: noSetup,
		},
		"maps": {
			route: []string{"maps"}, successArgv: []string{"maps"}, emptyArgv: []string{"maps"},
			blocks:        []string{"maps", "help"},
			successMarker: "maps[1]{map,title,type,state,blockers,path}:\n", emptyMarker: "maps[0]{map,title,type,state,blockers,path}:\n", usage: "usage: bench maps", setupSuccess: setupAXIMap, setupEmpty: noSetup,
		},
		"guards": {
			route: []string{"guards"}, successArgv: []string{"guards"}, emptyArgv: []string{"guards"},
			blocks:        []string{"guards", "guard_scan", "help"},
			successMarker: "guards[1]{guard,boundary,denies,branch,provenance,currency,wired}:\n", emptyMarker: "guards[0]{guard,boundary,denies,branch,provenance,currency,wired}:\n", usage: "usage: bench guards", setupSuccess: noSetup, setupEmpty: setupAXIEmptyGuards,
		},
		"diff": {
			route: []string{"diff"}, successArgv: []string{"diff"}, emptyArgv: []string{"diff"},
			blocks:        []string{"revision", "aggregate", "files", "checkout", "whitespace", "help"},
			successMarker: "files[1]{status,path,kind}:\n", emptyMarker: "files[0]{status,path,kind}:\n", usage: "usage: bench diff", setupSuccess: setupAXIDiff, setupEmpty: noSetup,
		},
		"consumers": {
			route: []string{"consumers"}, successArgv: []string{"consumers", "target.Symbol"}, emptyArgv: []string{"consumers", "target.Symbol"},
			blocks:        []string{"consumers", "meta", "citation", "help"},
			successMarker: "consumers[1]{file,line,via,enclosing}:\n", emptyMarker: "consumers[0]{file,line,via,enclosing}:\n", usage: "usage: bench consumers", setupSuccess: setupAXIConsumers, setupEmpty: setupAXIEmptyConsumers,
		},
		"coverage": {
			route: []string{"coverage"}, successArgv: []string{"coverage", "fixture"}, emptyArgv: []string{"coverage", "empty"},
			blocks:        []string{"spec", "state", "rows", "help"},
			successMarker: "rows[1]{story,behavior,seam}:\n", emptyMarker: "state: mapped\nrows[0]{story,behavior,seam}:\n", usage: "usage: bench coverage", setupSuccess: setupAXICoverage, setupEmpty: setupAXIEmptyCoverage,
		},
		"read": {
			route: []string{"read"}, successArgv: []string{"read", "fixture.txt"}, deepSuccessArgv: []string{"read", "../../fixture.txt"}, emptyArgv: []string{"read", "empty.txt"},
			blocks:        []string{"page", "content", "help"},
			successMarker: "content[1]{data}:\n", emptyMarker: "content[0]{data}:\n", usage: "usage: bench read",
			setupSuccess: func(t *testing.T, root string) { writeAXIFixture(t, filepath.Join(root, "fixture.txt"), "fixture\n") },
			setupEmpty:   func(t *testing.T, root string) { writeAXIFixture(t, filepath.Join(root, "empty.txt"), "") },
		},
		"worktree list": {
			route: []string{"worktree", "list"}, successArgv: []string{"worktree", "list"}, emptyArgv: []string{"worktree", "list"},
			blocks:        []string{"worktrees", "help"},
			successMarker: "worktrees[1]{id,label,request,state,source,tree,lease,landed,ignored}:\n", emptyMarker: "worktrees[0]{id,label,request,state,source,tree,lease,landed,ignored}:\n", usage: "usage: bench worktree list", setupSuccess: setupAXIWorktree, setupEmpty: noSetup,
		},
		// harnesses projects a record compiled into the binary, so it declares no empty
		// case and no refusal case. See recordBacked.
		"harnesses": {
			route: []string{"harnesses"}, successArgv: []string{"harnesses"},
			blocks:        []string{"schema", "harnesses", "help"},
			successMarker: "harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}:\n", usage: "usage: bench harnesses", setupSuccess: noSetup, recordBacked: true,
		},
		"roadmap": {
			route: []string{"roadmap"}, successArgv: []string{"roadmap"}, emptyArgv: []string{"roadmap"},
			blocks:        []string{"roadmap", "board", "commitment_outlook", "commitment_blockers", "sequence", "drain", "help"},
			successMarker: "roadmap[", emptyMarker: "roadmap[0]{id,title,spec,spec_status,external_trigger,occurrence_count,occurrence_keys}:\n", usage: "usage: bench roadmap", setupSuccess: setupAXIRoadmap, setupEmpty: noSetup,
		},
	}
	for name, item := range assessmentEnvelopeCases() {
		cases[name] = item
	}
	return cases
}
