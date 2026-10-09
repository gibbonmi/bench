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
			block, renderErr := renderVerdicts("conflicts", []lifecycleVerdict{{".bench/link-manifest.tsv", "modified-managed"}})
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
	// unconditionally below (validateAgentsPath, stagedAgents, stagedClaude), so the
	// special-file check runs first and routes straight to a conflict instead.
	agentsPath := filepath.Join(root, "AGENTS.md")
	agentsSpecial := isSpecialFile(agentsPath)
	if agentsSpecial {
		conflicts = append(conflicts, lifecycleVerdict{"AGENTS.md", "project-owned"})
	} else if err := validateAgentsPath(agentsPath); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	claudeSpecial := isSpecialFile(filepath.Join(root, "CLAUDE.md"))
	if claudeSpecial {
		conflicts = append(conflicts, lifecycleVerdict{"CLAUDE.md", "project-owned"})
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
		lifecycle := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: "link", Observations: []LifecycleObservation{o}}})
		if lifecycle.Refusal != "" {
			fmt.Fprintln(stderr, lifecycle.Refusal)
			return 1, false
		}
		for _, d := range lifecycle.Decisions {
			if d.ManifestAction != manifestOmit {
				rows[d.Path] = d.Fingerprint
			}
			switch {
			case d.Kind == lifecycleAdd || d.Kind == lifecycleChange:
				// A seed is staged only once the owner installs it.
				if staged == "" {
					if staged, err = stagePlanEntry(stage, e, mode); err != nil {
						fmt.Fprintln(stderr, err)
						return 1, false
					}
				}
				changes = append(changes, stagedChange{rel: d.Path, stage: staged})
			case d.Reason == "modified-managed" || d.Reason == "project-owned":
				conflicts = append(conflicts, lifecycleVerdict{d.Path, d.Reason})
			}
		}
	}
	// A dropped old row leaves only when it is still clean. Modified rows remain owned.
	for _, row := range old.Rows() {
		// CLAUDE.md is conditionally owned and is reconciled by stagedClaude below.
		if row.rel == "CLAUDE.md" {
			full := filepath.Join(root, row.rel)
			if fp, err := fingerprintPath(full); err == nil && fp != row.hash && !reclaimableClaude(full) {
				rows[row.rel] = row.hash
				conflicts = append(conflicts, lifecycleVerdict{row.rel, "modified-managed"})
			}
			continue
		}
		// The lifecycle decision above already chose the row of every planned path.
		if planned[row.rel] {
			continue
		}
		if _, present := rows[row.rel]; present {
			continue
		}
		full, ok := resolveInside(root, row.rel)
		if !ok || hasSymlinkParent(root, row.rel) {
			fmt.Fprintln(stderr, "refused manifest path: "+row.rel)
			return 1, false
		}
		if info, err := os.Lstat(full); err == nil && !info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0 {
			fmt.Fprintln(stderr, "refused non-regular dropped asset: "+row.rel)
			return 1, false
		}
		if _, err := os.Lstat(full); err == nil {
			fp, err := fingerprintPath(full)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1, false
			}
			if fp == row.hash {
				changes = append(changes, stagedChange{rel: row.rel})
			} else {
				rows[row.rel] = row.hash
				conflicts = append(conflicts, lifecycleVerdict{row.rel, "kept-modified-removed"})
			}
		}
	}
	if !agentsSpecial {
		agents, err := stagedRepairAgents(stage, root, repair)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		if agents != "" {
			changes = append(changes, stagedChange{rel: "AGENTS.md", stage: agents})
		}
	}
	if !claudeSpecial {
		claude, managed, err := stagedClaude(stage, root)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		if managed {
			fp, _ := fingerprintPath(claude)
			rows["CLAUDE.md"] = fp
			changes = append(changes, stagedChange{rel: "CLAUDE.md", stage: claude})
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
			block, err := renderVerdicts("conflicts", []lifecycleVerdict{{hook.Path, "modified-managed"}})
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
	changes = append(changes, stagedChange{rel: "pre-push", dest: hook.Path, stage: hookStage})
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
