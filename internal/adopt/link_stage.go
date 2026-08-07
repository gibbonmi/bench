package adopt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// link_stage.go holds transactionalLink's staging helpers - split out of
// link_transaction.go (which keeps the transaction's own conflict/reconcile logic) to
// stay under the repo's per-file line budget.

func stagePlanEntry(dir string, e planEntry, mode string) (string, error) {
	name := fmt.Sprintf("asset-%x", hashBytes([]byte(e.rel)))
	if e.kind == "adapter" {
		target, ok := AdapterTarget(e.rel)
		if !ok {
			return "", fmt.Errorf("adapter target unavailable for %s", e.rel)
		}
		return stageSymlink(dir, name, target)
	}
	if mode == "symlink" && e.kind != "inline" && e.kind != "inline-exec" && e.kind != "seed" {
		return stageSymlink(dir, name, e.src)
	}
	if e.kind == "inline" || e.kind == "seed" {
		return stageBytes(dir, name, []byte(e.content), 0o644)
	}
	// inline-exec carries generated content (bench setup's gate.sh) that must land
	// executable, same staged-write path as "inline" with a different mode bit — not
	// a second write mechanism.
	if e.kind == "inline-exec" {
		return stageBytes(dir, name, []byte(e.content), 0o755)
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

// stagedClaude converges CLAUDE.md to the canonical Bench form when it is absent or
// already one of the known bench-generated forms. Any other existing content -
// including a pre-existing empty file - is project-owned and left untouched: bench
// never writes into a CLAUDE.md path a user already claimed, even with zero bytes
// (see the pre-existing-empty-CLAUDE.md regression guard in the link/unlink surface
// contracts), and never injects the import lines into prose it didn't write (see the
// "relink injected an import into a project-owned CLAUDE.md" guard).
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
	return protectedBranch(root)
}
func stageManifest(stage, version string, rows []manifestRow) (string, error) {
	return stageBytes(stage, "manifest", manifestBytes(version, rows), 0o644)
}

func stageManagedPrePush(root string, health PrePushHealth) (string, []string, error) {
	return stageBeside(health.Path, renderedPrePush(root), 0o755)
}

func renderedPrePush(root string) []byte {
	return []byte(renderPrePush(root))
}

func renderPrePush(root string) string {
	return renderPrePushBranch(hookBranch(root))
}

func renderPrePushBranch(branch string) string {
	return strings.ReplaceAll(prePushTemplate, prePushBranchToken, branch)
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

// isSpecialFile reports whether path exists and is something other than a regular
// file or symlink — a FIFO, socket, device node, or directory. Lstat alone never
// blocks, unlike opening the path for read, so callers can use this to route the path
// straight to a conflict instead of ever attempting to read it. Every call site
// (link_transaction.go's AGENTS.md/CLAUDE.md guards, doctor_rows.go's per-row checks)
// fixes one of those two instruction-file paths; it never classifies a general plan
// entry (those go through their own os.Stat(parent).IsDir() check above), so treating
// a directory as "special" here cannot change how a legitimate directory plan entry is
// handled - it only turns a directory sitting where AGENTS.md/CLAUDE.md belongs into
// the same preserved conflict a FIFO already gets, instead of a raw read error.
func isSpecialFile(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode&os.ModeSymlink == 0 && !mode.IsRegular()
}
