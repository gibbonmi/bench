package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type lifecycleVerdict struct{ rel, reason string }

// transactionalLink stages and promotes plan into root as one FT84 transaction, and
// reports whether anything on disk actually changed (the second return) alongside the
// usual 0/1/2/3 result - a caller that wants to distinguish "converged, nothing to do"
// from "converged, wrote something" (bench setup's already-converged report) reads the
// bool; bench link ignores it.
func transactionalLink(root, kit, mode, version string, plan []planEntry, stdout, stderr io.Writer) (int, bool) {
	old, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	accepted := make([]planEntry, 0, len(plan))
	planned := make(map[string]bool, len(plan))
	conflicts := []lifecycleVerdict{}
	// A FIFO/socket/device at AGENTS.md or CLAUDE.md must never be opened for read —
	// os.ReadFile on a FIFO with no writer on the other end blocks forever. Both are
	// read unconditionally below (validateAgentsPath, stagedAgents, stagedClaude), so
	// the special-file check runs first and routes straight to a conflict instead.
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
	prepush := filepath.Join(hooksDir(root), "pre-push")
	if content, err := os.ReadFile(prepush); err == nil && !strings.Contains(string(content), PrePushMarker) {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", prepush)
		return 1, false
	}
	for _, e := range plan {
		planned[e.rel] = true
		// "seed" (seed-if-absent) entries - bench setup's profile - are neither a
		// managed/converged asset nor a conflict candidate: an existing file at the
		// path is reviewer-owned judgment content and is skipped silently (no
		// conflict, no manifest row); an absent one is staged and promoted with the
		// rest of the transaction like any other write.
		if e.kind == "seed" {
			if hasSymlinkParent(root, e.rel) {
				fmt.Fprintf(stderr, "conflict: %s has a symlink parent directory\n", e.rel)
				return 1, false
			}
			if _, err := os.Lstat(filepath.Join(root, e.rel)); err == nil {
				continue
			}
			accepted = append(accepted, e)
			continue
		}
		if e.kind != "inline" && e.kind != "inline-exec" {
			if info, err := os.Stat(e.src); err != nil || !info.Mode().IsRegular() {
				fmt.Fprintf(stderr, "conflict: kit asset missing: %s\n", e.src)
				return 1, false
			}
		}
		if hasSymlinkParent(root, e.rel) {
			fmt.Fprintf(stderr, "conflict: %s has a symlink parent directory\n", e.rel)
			return 1, false
		}
		parent := filepath.Join(root, filepath.Dir(e.rel))
		if info, err := os.Stat(parent); err == nil && !info.IsDir() {
			fmt.Fprintf(stderr, "conflict: parent path for %s is not a directory\n", e.rel)
			return 1, false
		}
		if _, err := os.Lstat(filepath.Join(root, e.rel)); err == nil && !manifestOwnedClean(root, e.rel) {
			reason := "project-owned"
			if old.Hash(e.rel) != "" {
				reason = "modified-managed"
			}
			conflicts = append(conflicts, lifecycleVerdict{e.rel, reason})
			continue
		}
		accepted = append(accepted, e)
	}
	stage, err := os.MkdirTemp(root, ".bench-link-")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	defer os.RemoveAll(stage)
	changes := []stagedChange{}
	rows := map[string]string{}
	for _, e := range accepted {
		p, err := stagePlanEntry(stage, e, mode)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		// A seed entry is never recorded as a managed row: recording it would make a
		// later reviewer hand-edit read back as a modified-managed conflict on the
		// next run, which defeats "seed-if-absent, then reviewer-owned".
		if e.kind != "seed" {
			fp, err := fingerprintPath(p)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1, false
			}
			rows[e.rel] = fp
		}
		changes = append(changes, stagedChange{rel: e.rel, stage: p, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
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
		if planned[row.rel] {
			// A current-plan collision was preserved, so retain its previous ownership
			// row rather than treating it as a removed asset.
			if _, accepted := rows[row.rel]; !accepted {
				rows[row.rel] = row.hash
			}
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
				changes = append(changes, stagedChange{rel: row.rel, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
			} else {
				rows[row.rel] = row.hash
				conflicts = append(conflicts, lifecycleVerdict{row.rel, "kept-modified-removed"})
			}
		}
	}
	if !agentsSpecial {
		agents, err := stagedAgents(stage, root)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		if agents != "" {
			changes = append(changes, stagedChange{rel: "AGENTS.md", stage: agents, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
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
			changes = append(changes, stagedChange{rel: "CLAUDE.md", stage: claude, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
		}
	}
	hookDest := filepath.Join(hooksDir(root), "pre-push")
	hookStage, hookDirs, err := stageBeside(hookDest, []byte(strings.ReplaceAll(prePushTemplate, prePushBranchToken, hookBranch(root))), 0o755)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	promoted := false
	defer func() {
		_ = os.Remove(hookStage)
		if !promoted {
			removeEmptyDirs(hookDirs)
		}
	}()
	changes = append(changes, stagedChange{rel: "pre-push", dest: hookDest, stage: hookStage, backup: hookStage + ".backup"})
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
	changes = append(changes, stagedChange{rel: ".bench/link-manifest.tsv", stage: manifest, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
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
	if err := promoteAll(root, changes); err != nil {
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
