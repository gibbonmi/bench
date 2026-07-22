package adopt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

type lifecycleVerdict struct{ rel, reason string }

func transactionalLink(root, kit, mode, version string, plan []planEntry, stdout, stderr io.Writer) int {
	old, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	accepted := make([]planEntry, 0, len(plan))
	planned := make(map[string]bool, len(plan))
	conflicts := []lifecycleVerdict{}
	if err := validateAgentsPath(filepath.Join(root, "AGENTS.md")); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	prepush := filepath.Join(hooksDir(root), "pre-push")
	if content, err := os.ReadFile(prepush); err == nil && !strings.Contains(string(content), PrePushMarker) {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", prepush)
		return 1
	}
	for _, e := range plan {
		planned[e.rel] = true
		if e.kind != "inline" {
			if info, err := os.Stat(e.src); err != nil || !info.Mode().IsRegular() {
				fmt.Fprintf(stderr, "conflict: kit asset missing: %s\n", e.src)
				return 1
			}
		}
		if hasSymlinkParent(root, e.rel) {
			fmt.Fprintf(stderr, "conflict: %s has a symlink parent directory\n", e.rel)
			return 1
		}
		parent := filepath.Join(root, filepath.Dir(e.rel))
		if info, err := os.Stat(parent); err == nil && !info.IsDir() {
			fmt.Fprintf(stderr, "conflict: parent path for %s is not a directory\n", e.rel)
			return 1
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
		return 1
	}
	defer os.RemoveAll(stage)
	changes := []stagedChange{}
	rows := map[string]string{}
	for _, e := range accepted {
		p, err := stagePlanEntry(stage, e, mode)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fp, err := fingerprintPath(p)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		rows[e.rel] = fp
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
			return 1
		}
		if info, err := os.Lstat(full); err == nil && !info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0 {
			fmt.Fprintln(stderr, "refused non-regular dropped asset: "+row.rel)
			return 1
		}
		if _, err := os.Lstat(full); err == nil {
			fp, err := fingerprintPath(full)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return 1
			}
			if fp == row.hash {
				changes = append(changes, stagedChange{rel: row.rel, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
			} else {
				rows[row.rel] = row.hash
				conflicts = append(conflicts, lifecycleVerdict{row.rel, "kept-modified-removed"})
			}
		}
	}
	agents, err := stagedAgents(stage, root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if agents != "" {
		changes = append(changes, stagedChange{rel: "AGENTS.md", stage: agents, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
	}
	claude, managed, err := stagedClaude(stage, root)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if managed {
		fp, _ := fingerprintPath(claude)
		rows["CLAUDE.md"] = fp
		changes = append(changes, stagedChange{rel: "CLAUDE.md", stage: claude, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
	}
	hookDest := filepath.Join(hooksDir(root), "pre-push")
	hookStage, hookDirs, err := stageBeside(hookDest, []byte(strings.ReplaceAll(prePushTemplate, prePushBranchToken, hookBranch(root))), 0o755)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
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
		return 1
	}
	changes = append(changes, stagedChange{rel: ".bench/link-manifest.tsv", stage: manifest, backup: filepath.Join(stage, fmt.Sprintf("backup-%d", len(changes)))})
	if len(conflicts) > 0 {
		if _, err := renderVerdicts("conflicts", conflicts); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	if err := promoteAll(root, changes); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	promoted = true
	if len(conflicts) > 0 {
		block, _ := renderVerdicts("conflicts", conflicts)
		fmt.Fprint(stdout, block)
		return 3
	}
	return 0
}

func stagePlanEntry(dir string, e planEntry, mode string) (string, error) {
	name := fmt.Sprintf("asset-%x", hashBytes([]byte(e.rel)))
	if e.kind == "adapter" {
		target, ok := AdapterTarget(e.rel)
		if !ok {
			return "", fmt.Errorf("adapter target unavailable for %s", e.rel)
		}
		return stageSymlink(dir, name, target)
	}
	if mode == "symlink" && e.kind != "inline" {
		return stageSymlink(dir, name, e.src)
	}
	if e.kind == "inline" {
		return stageBytes(dir, name, []byte(e.content), 0o644)
	}
	b, err := os.ReadFile(e.src)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(e.src)
	if err != nil {
		return "", err
	}
	return stageBytes(dir, name, b, info.Mode().Perm())
}
func stagedAgents(stage, root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if os.IsNotExist(err) {
		return stageBytes(stage, "agents", []byte(BenchAgentsBlock()), 0o644)
	}
	if err != nil {
		return "", err
	}
	next, err := RewriteAgentsBlock(string(b))
	if err != nil {
		return "", err
	}
	if next == string(b) {
		return "", nil
	}
	return stageBytes(stage, "agents", []byte(next), 0o644)
}
func stagedClaude(stage, root string) (string, bool, error) {
	b, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if os.IsNotExist(err) || string(b) == legacyClaudeMD() || string(b) == benchClaudeMD() {
		p, e := stageBytes(stage, "claude", []byte(benchClaudeMD()), 0o644)
		return p, true, e
	}
	return "", false, err
}
func reclaimableClaude(path string) bool {
	b, err := os.ReadFile(path)
	return err == nil && (string(b) == legacyClaudeMD() || string(b) == benchClaudeMD())
}
func hookBranch(root string) string {
	if out, err := git.Output("-C", root, "ls-remote", "--symref", "origin", "HEAD"); err == nil {
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "ref: refs/heads/") {
				fields := strings.Fields(line)
				if len(fields) > 1 {
					return strings.TrimPrefix(fields[1], "refs/heads/")
				}
			}
		}
	}
	return git.DefaultBranch(root)
}
func stageManifest(stage, version string, rows []manifestRow) (string, error) {
	return stageBytes(stage, "manifest", manifestBytes(version, rows), 0o644)
}
func renderVerdicts(name string, vs []lifecycleVerdict) (string, error) {
	rows := make([][]string, len(vs))
	for i, v := range vs {
		rows[i] = []string{v.rel, v.reason}
	}
	return toon.Table(name, []string{"path", "reason"}, rows)
}

func stageBeside(dest string, data []byte, mode os.FileMode) (string, []string, error) {
	dir := filepath.Dir(dest)
	created := missingDirs(dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", nil, err
	}
	f, err := os.CreateTemp(dir, ".bench-link-stage-")
	if err != nil {
		removeEmptyDirs(created)
		return "", nil, err
	}
	path := f.Name()
	err = f.Chmod(mode)
	if err == nil {
		err = writeSyncClose(path, f, data)
	}
	if err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		removeEmptyDirs(created)
		return "", nil, err
	}
	return path, created, nil
}

func stageSymlink(dir, name, target string) (string, error) {
	path := filepath.Join(dir, name)
	if err := os.Symlink(target, path); err != nil {
		return "", err
	}
	if err := syncDirectory(dir); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("sync staged symlink directory %s: %w", dir, err)
	}
	return path, nil
}

func missingDirs(dir string) []string {
	var dirs []string
	for current := dir; ; current = filepath.Dir(current) {
		if _, err := os.Lstat(current); err == nil {
			break
		}
		dirs = append(dirs, current)
		if parent := filepath.Dir(current); parent == current {
			break
		}
	}
	return dirs
}

func removeEmptyDirs(dirs []string) {
	for _, dir := range dirs {
		_ = os.Remove(dir)
	}
}
