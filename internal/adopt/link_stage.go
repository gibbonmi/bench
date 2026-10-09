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

// observeClaude stages the canonical CLAUDE.md form and records the facts that the
// conditional owner decides from. It returns the staged replacement. The known-form fact
// comes from the bytes that the path reads, so a path that reads as absent, such as a
// dangling link, is absent, and any other read error is returned. Content outside the
// known forms, including a pre-existing empty file, is project content: bench never writes
// into a CLAUDE.md path a user already claimed, and never injects the import lines into
// prose it did not write. The caller observes a special file through isSpecialFile instead,
// because reading a writerless FIFO never returns.
func observeClaude(o *LifecycleObservation, stage, root string) (string, error) {
	staged, err := stageBytes(stage, "claude", []byte(benchClaudeMD()), 0o644)
	if err != nil {
		return "", err
	}
	dest := filepath.Join(root, "CLAUDE.md")
	if err := observeDestination(o, dest, staged); err != nil {
		return "", err
	}
	b, err := os.ReadFile(dest)
	if os.IsNotExist(err) {
		o.Present = false
		return staged, nil
	}
	o.KnownForm = string(b) == legacyClaudeMD() || string(b) == benchClaudeMD()
	return staged, err
}

// observeWithdrawal records the facts of a manifest row that the plan no longer names. It
// inspects no path that escapes the repository or crosses a symlink parent, and it
// fingerprints only a regular file or a symlink, so it never opens a special file. A leaf
// that Lstat cannot inspect is absent. A fingerprint failure is returned, because an
// unreadable dropped asset is a hard failure.
func observeWithdrawal(root string, row manifestRow, order int) (LifecycleObservation, error) {
	o := LifecycleObservation{Path: row.rel, Order: order, Recorded: row.hash}
	full, inside := resolveInside(root, row.rel)
	o.Escapes = !inside
	o.SymlinkParent = inside && hasSymlinkParent(root, row.rel)
	if o.Escapes || o.SymlinkParent {
		return o, nil
	}
	info, err := os.Lstat(full)
	if err != nil {
		return o, nil
	}
	o.Present, o.FileKind = true, fileKind(info.Mode())
	if o.FileKind == "other" {
		return o, nil
	}
	o.Observed, err = fingerprintPath(full)
	return o, err
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

// observeDestination records the staged entry and the destination facts that the
// lifecycle owner compares. A fingerprint covers content only, so a present destination
// also records both modes: a kit asset can change its executable bit without changing a
// byte. An unreadable destination keeps an empty Observed fingerprint.
//
// A staged symlink also records what dest resolves to. A repo may satisfy a whole adapter
// directory with one directory-level symlink (.claude/commands -> ../.agents/commands).
// That symlink leaves dest resolving through its parent to the very file the staged link
// names, so a reader of dest sees the same bytes as through the staged link.
func observeDestination(o *LifecycleObservation, dest, staged string) error {
	stagedInfo, err := os.Lstat(staged)
	if err != nil {
		return err
	}
	o.StagedKind = fileKind(stagedInfo.Mode())
	if o.Staged, err = fingerprintPath(staged); err != nil {
		return err
	}
	destInfo, err := os.Lstat(dest)
	if err != nil {
		return nil
	}
	o.Present, o.FileKind = true, fileKind(destInfo.Mode())
	o.Perm, o.StagedPerm = destInfo.Mode().Perm(), stagedInfo.Mode().Perm()
	o.Observed, _ = fingerprintPath(dest)
	if target, ok := stagedTarget(dest, staged); ok {
		o.CanonicalAdapter = sameRegularFile(dest, target)
		o.ResolvedEquivalent = sameRegularContent(dest, target)
	}
	return nil
}

// stagedTarget returns the file that a staged symlink names once it is promoted to dest.
// stageSymlink writes the link inside the transaction's stage directory, so a relative
// target resolves against dest's own directory.
func stagedTarget(dest, staged string) (string, bool) {
	target, err := os.Readlink(staged)
	if err != nil {
		return "", false
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(dest), target)
	}
	return target, true
}

// sameRegularFile reports whether two paths resolve to one regular file.
func sameRegularFile(a, b string) bool {
	aInfo, err := os.Stat(a)
	if err != nil {
		return false
	}
	bInfo, err := os.Stat(b)
	return err == nil && aInfo.Mode().IsRegular() && bInfo.Mode().IsRegular() && os.SameFile(aInfo, bInfo)
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
