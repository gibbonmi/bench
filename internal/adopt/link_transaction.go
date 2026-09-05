package adopt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type lifecycleVerdict struct{ rel, reason string }

// stagedEntry pairs a plan entry with the staged file promotion will rename into place.
// Classification stages every entry it inspects because the skip decision compares the
// destination against exactly the bytes and mode this plan would write. Staging is the
// only definition of those.
type stagedEntry struct {
	entry planEntry
	stage string
}

// convergedFingerprint returns dest's fingerprint when dest already holds exactly what
// promoting staged would leave there, and "" when the entry still needs a write. The
// permission bits are compared alongside the fingerprint because a fingerprint covers
// content only. A kit asset can change its executable bit without changing a byte.
func convergedFingerprint(dest, staged string) string {
	destInfo, err := os.Lstat(dest)
	if err != nil {
		return ""
	}
	stagedInfo, err := os.Lstat(staged)
	if err != nil {
		return ""
	}
	if stagedInfo.Mode()&os.ModeSymlink != 0 {
		return convergedSymlinkFingerprint(dest, staged)
	}
	if destInfo.Mode()&os.ModeSymlink == 0 && destInfo.Mode().Perm() != stagedInfo.Mode().Perm() {
		return ""
	}
	destPrint, err := fingerprintPath(dest)
	if err != nil {
		return ""
	}
	stagedPrint, err := fingerprintPath(staged)
	if err != nil || destPrint != stagedPrint {
		return ""
	}
	return destPrint
}

// convergedSymlinkFingerprint answers convergedFingerprint for a staged symlink, whose
// own permission bits carry nothing to compare. An identical link at dest is not the only
// converged shape. A repo may satisfy a whole adapter directory with one directory-level
// symlink (.claude/commands -> ../.agents/commands). That symlink leaves dest resolving
// through its parent to the very file the staged link names. Both shapes are converged,
// because a reader of dest sees the same bytes either way. Refusing the second shape
// would send an untouched repo into the symlink-parent refusal on every entry.
func convergedSymlinkFingerprint(dest, staged string) string {
	destPrint, err := fingerprintPath(dest)
	if err != nil {
		return ""
	}
	if stagedPrint, err := fingerprintPath(staged); err == nil && stagedPrint == destPrint {
		return destPrint
	}
	target, err := os.Readlink(staged)
	if err != nil {
		return ""
	}
	// stageSymlink writes the link inside the transaction's stage directory, so a relative
	// target only names its file once promoted: resolve it against dest's own directory.
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(dest), target)
	}
	if !sameRegularContent(dest, target) {
		return ""
	}
	return destPrint
}

func sameAdapterTarget(dest, staged string) bool {
	target, err := os.Readlink(staged)
	if err != nil {
		return false
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(dest), target)
	}
	destInfo, err := os.Stat(dest)
	if err != nil {
		return false
	}
	targetInfo, err := os.Stat(target)
	return err == nil && destInfo.Mode().IsRegular() && targetInfo.Mode().IsRegular() && os.SameFile(destInfo, targetInfo)
}

// sameRegularContent reports whether two paths resolve to regular files holding the same
// bytes. Each is stat'd through its links first, because a FIFO or device reached by
// either path would block the read forever.
func sameRegularContent(a, b string) bool {
	for _, path := range []string{a, b} {
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			return false
		}
	}
	first, err := os.ReadFile(a)
	if err != nil {
		return false
	}
	second, err := os.ReadFile(b)
	return err == nil && bytes.Equal(first, second)
}

// ownedUnmodified reports whether dest still carries the exact bytes recorded for it in
// the previous manifest. owned is that manifest's hash, and "" means unowned.
func ownedUnmodified(dest, owned string) bool {
	if owned == "" {
		return false
	}
	fp, err := fingerprintPath(dest)
	return err == nil && fp == owned
}

// transactionalLink stages and promotes plan into root as one transaction. It reports
// whether anything on disk actually changed (the second return) alongside the usual
// 0/1/2/3 result. bench setup's already-converged report reads the bool to distinguish
// "converged, nothing to do" from "converged, wrote something". bench link ignores it.
func transactionalLink(root, kit, mode, version string, plan []planEntry, stdout, stderr io.Writer) (int, bool) {
	old, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	accepted := make([]stagedEntry, 0, len(plan))
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
	populateOriginHead(root)
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
	for _, e := range plan {
		planned[e.rel] = true
		// "seed" (seed-if-absent) entries, such as bench setup's profile, are neither a
		// managed/converged asset nor a conflict candidate. An existing file at the path is
		// reviewer-owned judgment content and is skipped silently, with no conflict and no
		// manifest row. An absent one is staged and promoted with the rest of the transaction
		// like any other write.
		if e.kind == "seed" {
			if hasSymlinkParent(root, e.rel) {
				fmt.Fprintf(stderr, "conflict: %s has a symlink parent directory\n", e.rel)
				return 1, false
			}
			if _, err := os.Lstat(filepath.Join(root, e.rel)); err == nil {
				continue
			}
			p, err := stagePlanEntry(stage, e, mode)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1, false
			}
			accepted = append(accepted, stagedEntry{e, p})
			continue
		}
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
		// A manifest-owned destination that already holds what this plan would write needs no
		// write at all. It leaves the transaction here, never renamed over and never a conflict
		// candidate, carrying its own fingerprint into the manifest. The comparison uses the
		// incoming kit bytes rather than the recorded hash, so a locally stale kit asset is
		// still caught. Keying the skip on the old hash would make every release a content
		// no-op instead. Ownership still gates the skip, except for a byte-identical adapter
		// resolved through a symlink parent to its canonical target. Everything past this point
		// wants to write, so a symlinked parent is a hard refusal that outranks the soft
		// per-entry conflict report.
		dest := filepath.Join(root, e.rel)
		_, statErr := os.Lstat(dest)
		exists := statErr == nil
		owned := old.Hash(e.rel)
		staged, err := stagePlanEntry(stage, e, mode)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1, false
		}
		parentSymlink := hasSymlinkParent(root, e.rel)
		adoptConvergedAdapter := owned == "" && e.kind == "adapter" && parentSymlink && sameAdapterTarget(dest, staged)
		if exists && (owned != "" || adoptConvergedAdapter) {
			if fp := convergedFingerprint(dest, staged); fp != "" {
				rows[e.rel] = fp
				continue
			}
		}
		if parentSymlink {
			fmt.Fprintf(stderr, "conflict: %s has a symlink parent directory\n", e.rel)
			return 1, false
		}
		// An owned destination still carrying the bytes the manifest recorded is the consumer's
		// untouched copy of an older kit. It is rewritten rather than reported. Only a
		// destination that answers to neither the manifest nor the incoming kit is someone's
		// local edit to preserve.
		if exists && !ownedUnmodified(dest, owned) {
			reason := "project-owned"
			if owned != "" {
				reason = "modified-managed"
			}
			conflicts = append(conflicts, lifecycleVerdict{e.rel, reason})
			continue
		}
		accepted = append(accepted, stagedEntry{e, staged})
	}
	changes := []stagedChange{}
	for _, a := range accepted {
		// A seed entry is never recorded as a managed row. A recorded row would make a later
		// reviewer hand-edit read back as a modified-managed conflict on the next run, which
		// defeats "seed-if-absent, then reviewer-owned".
		if a.entry.kind != "seed" {
			fp, err := fingerprintPath(a.stage)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1, false
			}
			rows[a.entry.rel] = fp
		}
		changes = append(changes, stagedChange{rel: a.entry.rel, stage: a.stage, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
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
	hookStage, hookDirs, err := stageManagedPrePush(root, hook)
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
	changes = append(changes, stagedChange{rel: "pre-push", dest: hook.Path, stage: hookStage, backup: hookStage + ".backup"})
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
