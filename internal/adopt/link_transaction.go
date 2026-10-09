package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type lifecycleVerdict struct{ rel, reason string }

// transactionalLink stages and promotes plan into root as one transaction. It reports
// whether anything on disk actually changed (the second return) alongside the usual
// 0/1/2/3 result. bench setup's already-converged report reads the bool to distinguish
// "converged, nothing to do" from "converged, wrote something". bench link ignores it.
func transactionalLink(root, kit, mode, version string, plan []planEntry, stdout, stderr io.Writer) (int, bool) {
	return transactionalRepair(root, kit, mode, version, plan, stdout, stderr, nil)
}

func transactionalRepair(root, kit, mode, version string, plan []planEntry, stdout, stderr io.Writer, repair *repairRun) (int, bool) {
	lease, err := lockAdoption(root, plan)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	defer lease.Close()
	lease.repair = repair
	old, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	if repair != nil {
		if err := validateRepairManifest(filepath.Join(root, ".bench", "link-manifest.tsv"), old); err != nil {
			block, renderErr := renderVerdicts("conflicts", []lifecycleVerdict{{".bench/link-manifest.tsv", reasonModifiedManaged}})
			if renderErr != nil {
				fmt.Fprintln(stderr, renderErr)
			} else {
				fmt.Fprint(stdout, block)
			}
			fmt.Fprintln(stderr, err)
			return 1, false
		}
	}
	planned := make(map[string]bool, len(plan))
	conflicts := []lifecycleVerdict{}
	// A FIFO/socket/device at AGENTS.md or CLAUDE.md must never be opened for read.
	// os.ReadFile on a FIFO with no writer on the other end blocks forever. Both are read
	// below (validateAgentsPath, stagedRepairAgents, observeClaude), so the special-file
	// check runs first and routes straight to a conflict instead.
	agentsPath := filepath.Join(root, "AGENTS.md")
	agentsSpecial := isSpecialFile(agentsPath)
	if agentsSpecial {
		conflicts = append(conflicts, lifecycleVerdict{"AGENTS.md", reasonProjectOwned})
	} else if err := validateAgentsPath(agentsPath); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	claudeSpecial := isSpecialFile(filepath.Join(root, "CLAUDE.md"))
	if claudeSpecial {
		conflicts = append(conflicts, lifecycleVerdict{"CLAUDE.md", reasonProjectOwned})
	}
	if repair == nil {
		populateOriginHead(root)
	}
	// The hooks directory is resolved through the reader before InspectPrePush, and before
	// any file is staged, so an unresolved answer refuses the whole transaction rather than
	// leaving InspectPrePush's own absent-state fallback silently skip the hook.
	if _, err := hooksDir(root); err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", root, err)
		return 1, false
	}
	hook := InspectPrePush(root)
	_, hookPresent := os.Lstat(hook.Path)
	if hook.State == PrePushForeign || (hook.State == PrePushDiverted && hookPresent == nil) {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", hook.Path)
		return 1, false
	}
	stage, err := os.MkdirTemp(root, ".bench-link-")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	defer os.RemoveAll(stage)
	rows := map[string]string{}
	changes := []stagedChange{}
	// Each entry runs its hard checks, then takes its lifecycle decision, before the next
	// entry is inspected. A refusal therefore wins over every later entry's checks.
	for order, e := range plan {
		planned[e.rel] = true
		dest := filepath.Join(root, e.rel)
		o := LifecycleObservation{Path: e.rel, Order: order, PlanKind: e.kind, Recorded: old.Hash(e.rel), Strict: repair != nil}
		staged := ""
		if e.kind == "seed" {
			_, err := os.Lstat(dest)
			o.Present = err == nil
		} else {
			if e.kind != "inline" && e.kind != "inline-exec" {
				if info, err := os.Stat(e.src); err != nil || !info.Mode().IsRegular() {
					fmt.Fprintf(stderr, "conflict: kit asset missing: %s\n", e.src)
					return 1, false
				}
			}
			parent := filepath.Join(root, filepath.Dir(e.rel))
			if info, err := os.Stat(parent); err == nil && !info.IsDir() {
				fmt.Fprintf(stderr, "conflict: parent path for %s is not a directory\n", e.rel)
				return 1, false
			}
			if staged, err = stagePlanEntry(stage, e, mode); err == nil {
				err = observeDestination(&o, dest, staged)
			}
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1, false
			}
		}
		o.SymlinkParent = hasSymlinkParent(root, e.rel)
		d, ok := decideLink(o, rows, stderr)
		if !ok {
			return 1, false
		}
		switch {
		case publishes(d):
			// A seed is staged only once the owner installs it.
			if staged == "" {
				if staged, err = stagePlanEntry(stage, e, mode); err != nil {
					fmt.Fprintln(stderr, err)
					return 1, false
				}
			}
			changes = append(changes, stagedChange{rel: d.Path, stage: staged})
		case d.Reason == reasonModifiedManaged || d.Reason == reasonProjectOwned:
			conflicts = append(conflicts, lifecycleVerdict{d.Path, d.Reason})
		}
	}
	// A dropped row is an old row that the plan no longer names. CLAUDE.md is conditionally
	// owned, so the conditional owner below decides its row instead of the withdrawal rule.
	for order, row := range old.Rows() {
		if planned[row.rel] || row.rel == "CLAUDE.md" {
			continue
		}
		o, err := observeWithdrawal(root, row, len(plan)+order)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		d, ok := decideLink(o, rows, stderr)
		if !ok {
			return 1, false
		}
		switch {
		case d.Kind == lifecycleRemove:
			changes = append(changes, stagedChange{rel: d.Path})
		case d.Reason == reasonKeptModifiedRemoved:
			conflicts = append(conflicts, lifecycleVerdict{d.Path, d.Reason})
		}
	}
	if !agentsSpecial {
		agents, err := stagedRepairAgents(stage, root, repair)
		o := LifecycleObservation{Path: "AGENTS.md", PlanKind: bespokePlanKind}
		if err == nil && agents != "" {
			err = observeDestination(&o, agentsPath, agents)
		}
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		d, ok := decideLink(o, rows, stderr)
		if !ok {
			return 1, false
		}
		if publishes(d) {
			changes = append(changes, stagedChange{rel: d.Path, stage: agents})
		}
	}
	if !claudeSpecial {
		o := LifecycleObservation{Path: "CLAUDE.md", PlanKind: claudePlanKind, Recorded: old.Hash("CLAUDE.md")}
		claude, err := observeClaude(&o, stage, root)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		d, ok := decideLink(o, rows, stderr)
		if !ok {
			return 1, false
		}
		switch {
		case publishes(d):
			changes = append(changes, stagedChange{rel: d.Path, stage: claude})
		case d.Reason == reasonModifiedManaged:
			// The report omits a project-owned CLAUDE.md. Only the special-file guard
			// above reports one.
			conflicts = append(conflicts, lifecycleVerdict{d.Path, d.Reason})
		}
	}
	hookStage, hookDirs, err := stageManagedPrePush(root, hook)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	if repair != nil && hook.State == PrePushManaged {
		if err := validateRepairHook(hook); err != nil {
			_ = os.Remove(hookStage)
			removeEmptyDirs(hookDirs)
			block, err := renderVerdicts("conflicts", []lifecycleVerdict{{hook.Path, reasonModifiedManaged}})
			if err != nil {
				fmt.Fprintln(stderr, err)
			} else {
				fmt.Fprint(stdout, block)
			}
			return 1, false
		}
	}
	promoted := false
	defer func() {
		_ = os.Remove(hookStage)
		if !promoted {
			removeEmptyDirs(hookDirs)
		}
	}()
	o := LifecycleObservation{Path: "pre-push", PlanKind: bespokePlanKind}
	if err := observeDestination(&o, hook.Path, hookStage); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	d, ok := decideLink(o, rows, stderr)
	if !ok {
		return 1, false
	}
	if publishes(d) {
		changes = append(changes, stagedChange{rel: d.Path, dest: hook.Path, stage: hookStage})
	}
	manifestRows := make([]manifestRow, 0, len(rows))
	for rel, hash := range rows {
		manifestRows = append(manifestRows, manifestRow{rel, hash})
	}
	sort.Slice(manifestRows, func(i, j int) bool { return manifestRows[i].rel < manifestRows[j].rel })
	manifest, err := stageManifest(stage, version, manifestRows)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	changes = append(changes, stagedChange{rel: ".bench/link-manifest.tsv", stage: manifest})
	if len(conflicts) > 0 {
		if _, err := renderVerdicts("conflicts", conflicts); err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
	}
	changed, err := changesModifyTree(root, changes)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	if err := promoteWithLease(root, changes, lease); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	promoted = true
	if len(conflicts) > 0 {
		block, _ := renderVerdicts("conflicts", conflicts)
		fmt.Fprint(stdout, block)
		return 3, changed
	}
	return 0, changed
}

// decideLink submits one link observation to the lifecycle owner and keeps the manifest
// row that its decision selects. It prints the owner's hard refusal and reports false.
func decideLink(o LifecycleObservation, rows map[string]string, stderr io.Writer) (LifecycleDecision, bool) {
	lifecycle := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: []LifecycleObservation{o}}})
	if lifecycle.Refusal != "" {
		fmt.Fprintln(stderr, lifecycle.Refusal)
		return LifecycleDecision{}, false
	}
	d := lifecycle.Decisions[0]
	if d.ManifestAction != manifestOmit {
		rows[d.Path] = d.Fingerprint
	}
	return d, true
}

// publishes reports whether a decision selects the write of its prepared replacement.
func publishes(d LifecycleDecision) bool {
	return d.Kind == lifecycleAdd || d.Kind == lifecycleChange
}
