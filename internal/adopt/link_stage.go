package adopt

import (
	"bytes"
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
	// inline-exec carries generated content, such as bench setup's gate.sh, that must land
	// executable. It follows the same staged-write path as "inline" with a different mode
	// bit, not a second write mechanism.
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
// already one of the known bench-generated forms. Any other existing content, including a
// pre-existing empty file, is project-owned and left untouched. bench never writes into a
// CLAUDE.md path a user already claimed, even with zero bytes; see the
// pre-existing-empty-CLAUDE.md regression guard in the link/unlink surface contracts.
// bench never injects the import lines into prose it did not write; see the "relink
// injected an import into a project-owned CLAUDE.md" guard.
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
	return stageBytes(stage, "manifest", manifestBytes(version, rows), manifestMode)
}

func stageManagedPrePush(root string, health PrePushHealth) (string, []string, error) {
	return stageBeside(health.Path, renderedPrePush(root), prePushMode)
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

// isSpecialFile reports whether path exists and is something other than a regular file or
// symlink: a FIFO, socket, device node, or directory. Lstat alone never blocks, unlike
// opening the path for read. A caller uses this to route the path straight to a conflict
// instead of ever attempting to read it. Every call site, link_transaction.go's
// AGENTS.md/CLAUDE.md guards and doctor_rows.go's per-row checks, fixes one of those two
// instruction-file paths. It never classifies a general plan entry; those go through
// their own os.Stat(parent).IsDir() check above. A directory sitting where
// AGENTS.md/CLAUDE.md belongs gets the same preserved conflict a FIFO already gets,
// instead of a raw read error.
func isSpecialFile(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode&os.ModeSymlink == 0 && !mode.IsRegular()
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
func ownedUnmodified(dest, owned, staged string, strict bool) bool {
	if owned == "" {
		return false
	}
	fp, err := fingerprintPath(dest)
	if err != nil || fp != owned {
		return false
	}
	if strict {
		before, beforeErr := os.Lstat(dest)
		after, afterErr := os.Lstat(staged)
		if beforeErr != nil || afterErr != nil {
			return false
		}
		if before.Mode().IsRegular() && after.Mode().IsRegular() && before.Mode().Perm() != after.Mode().Perm() {
			return false
		}
	}
	return true
}
