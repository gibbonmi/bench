package adopt

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
)

type planEntry struct {
	src, rel, kind string
}

func Link(args []string, stdout, stderr io.Writer, version string) int {
	mode := "copy"
	if len(args) > 0 {
		mode = args[0]
	}
	if mode != "copy" && mode != "symlink" {
		fmt.Fprintln(stderr, "usage: bench link [copy|symlink]")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, "not in a git repo")
		return 1
	}
	kit := kitDir()
	if isEphemeralKit(kit) && mode == "symlink" {
		mode = "copy"
		fmt.Fprintln(stdout, "(running from an ephemeral package cache - using copy mode so files don't dangle)")
	}
	plan := buildLinkPlan(kit)
	if !preflightLink(root, plan, stderr) {
		return 1
	}
	if err := writeAgentsFile(root); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := installClaudeMD(root); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := installPlan(root, mode, version, plan); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := installGitHook(root, stderr); err != nil {
		return 1
	}
	fmt.Fprintf(stdout, "linked Bench into %s (mode: %s).\n", root, mode)
	fmt.Fprintln(stdout, "  instructions: AGENTS.md managed block -> .bench/BENCH.md")
	fmt.Fprintln(stdout, "  portable surface: .agents/{skills,commands}")
	fmt.Fprintln(stdout, "  adapters: .claude/settings.json, .codex/hooks.json, .claude/{skills,commands} -> .agents")
	fmt.Fprintln(stdout, "  enforcement: shared .bench/hooks + git pre-push guard + the bench shift loop")
	fmt.Fprintln(stdout, "Run 'bench init' next to scaffold .bench/gate.sh.")
	return 0
}

func isEphemeralKit(kit string) bool {
	return strings.Contains(kit, "_npx") ||
		strings.Contains(kit, "/dlx-") ||
		strings.Contains(kit, "npm-cache") ||
		strings.Contains(kit, "/.npm/_cacache/")
}

func buildLinkPlan(kit string) []planEntry {
	var plan []planEntry
	appendTreeToPlan(&plan, filepath.Join(kit, "bin"), ".bench/bin", "file")
	if bin := currentExecutablePath(); bin != "" {
		plan = append(plan, planEntry{bin, ".bench/dist/bench", "file"})
	}
	for _, e := range []planEntry{
		{filepath.Join(kit, ".bench", "BENCH.md"), ".bench/BENCH.md", "file"},
		{filepath.Join(kit, ".bench", "BENCH-reference.md"), ".bench/BENCH-reference.md", "file"},
		{filepath.Join(kit, ".claude", "README.md"), ".claude/README.md", "file"},
		{filepath.Join(kit, ".claude", "settings.json"), ".claude/settings.json", "file"},
		{filepath.Join(kit, ".codex", "hooks.json"), ".codex/hooks.json", "file"},
	} {
		plan = append(plan, e)
	}
	appendTreeToPlan(&plan, filepath.Join(kit, ".bench", "hooks"), ".bench/hooks", "file")
	appendTreeToPlan(&plan, filepath.Join(kit, ".bench", "adapters"), ".bench/adapters", "file")
	appendTreeToPlan(&plan, filepath.Join(kit, ".bench", "lib"), ".bench/lib", "file")
	appendTreeToPlan(&plan, filepath.Join(kit, ".agents", "commands"), ".agents/commands", "file")
	appendTreeToPlan(&plan, filepath.Join(kit, ".agents", "skills"), ".agents/skills", "file")
	appendTreeToPlan(&plan, filepath.Join(kit, ".agents", "commands"), ".claude/commands", "adapter")
	appendClaudeSkillsToPlan(&plan, kit)
	return plan
}

func appendTreeToPlan(plan *[]planEntry, srcRoot, destRoot, kind string) {
	info, err := os.Stat(srcRoot)
	if err != nil || !info.IsDir() {
		return
	}
	var entries []planEntry
	_ = filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return nil
		}
		entries = append(entries, planEntry{path, filepath.ToSlash(filepath.Join(destRoot, rel)), kind})
		return nil
	})
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })
	*plan = append(*plan, entries...)
}

func appendClaudeSkillsToPlan(plan *[]planEntry, kit string) {
	srcRoot := filepath.Join(kit, ".agents", "skills")
	info, err := os.Stat(srcRoot)
	if err != nil || !info.IsDir() {
		return
	}
	var entries []planEntry
	_ = filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return nil
		}
		top := strings.Split(filepath.ToSlash(rel), "/")[0]
		if _, err := os.Stat(filepath.Join(kit, ".agents", "commands", top+".md")); err == nil {
			return nil
		}
		entries = append(entries, planEntry{path, filepath.ToSlash(filepath.Join(".claude/skills", rel)), "adapter"})
		return nil
	})
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })
	*plan = append(*plan, entries...)
}

func currentExecutablePath() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() || info.Size() == 0 || info.Mode().Perm()&0o111 == 0 {
		return ""
	}
	return path
}

func preflightLink(root string, plan []planEntry, stderr io.Writer) bool {
	conflicts := 0
	if err := validateAgentsPath(filepath.Join(root, "AGENTS.md")); err != nil {
		fmt.Fprintln(stderr, err)
		conflicts++
	}
	hooks := hooksDir(root)
	prepush := filepath.Join(hooks, "pre-push")
	if content, err := os.ReadFile(prepush); err == nil && !strings.Contains(string(content), "bench:managed-pre-push") {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", prepush)
		conflicts++
	}
	for _, e := range plan {
		if e.rel == "" {
			continue
		}
		if info, err := os.Stat(e.src); err != nil || !info.Mode().IsRegular() {
			fmt.Fprintf(stderr, "conflict: kit asset missing: %s\n", e.src)
			conflicts++
			continue
		}
		if hasSymlinkParent(root, e.rel) {
			fmt.Fprintf(stderr, "conflict: %s has a symlink parent directory\n", e.rel)
			conflicts++
			continue
		}
		parent := filepath.Join(root, filepath.Dir(e.rel))
		if info, err := os.Stat(parent); err == nil && !info.IsDir() {
			fmt.Fprintf(stderr, "conflict: parent path for %s is not a directory\n", e.rel)
			conflicts++
			continue
		}
		dest := filepath.Join(root, e.rel)
		if _, err := os.Lstat(dest); err == nil {
			if !manifestOwnedClean(root, e.rel) {
				old := manifestHash(root, e.rel)
				if old != "" {
					fmt.Fprintf(stderr, "conflict: modified Bench-managed file: %s\n", e.rel)
				} else {
					fmt.Fprintf(stderr, "conflict: project-owned file exists: %s\n", e.rel)
				}
				conflicts++
			}
		}
	}
	return conflicts == 0
}

func validateAgentsPath(path string) error {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return validateAgentsContent(string(content))
}

func manifestHash(root, rel string) string {
	m, err := ReadManifest(filepath.Join(root, ".bench", "link-manifest.tsv"))
	if err != nil {
		return ""
	}
	return m.Hash(rel)
}

func manifestOwnedClean(root, rel string) bool {
	old := manifestHash(root, rel)
	if old == "" {
		return false
	}
	path := filepath.Join(root, rel)
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		return true
	}
	now, err := fingerprintPath(path)
	return err == nil && now == old
}

func hasSymlinkParent(root, rel string) bool {
	dir := filepath.Dir(rel)
	if dir == "." {
		return false
	}
	path := root
	for _, part := range strings.Split(filepath.ToSlash(dir), "/") {
		path = filepath.Join(path, part)
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func writeAgentsFile(root string) error {
	path := filepath.Join(root, "AGENTS.md")
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return os.WriteFile(path, []byte(BenchAgentsBlock()), 0o644)
	}
	if err != nil {
		return err
	}
	next, err := RewriteAgentsBlock(string(content))
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(next), 0o644)
}

func installClaudeMD(root string) error {
	path := filepath.Join(root, "CLAUDE.md")
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) || string(content) == legacyClaudeMD() {
		return os.WriteFile(path, []byte(benchClaudeMD()), 0o644)
	}
	return err
}

func installPlan(root, mode, version string, plan []planEntry) error {
	var rows []manifestRow
	for _, e := range plan {
		if err := installPlannedFile(root, mode, e); err != nil {
			return err
		}
		fp, err := fingerprintPath(filepath.Join(root, e.rel))
		if err != nil {
			return err
		}
		rows = append(rows, manifestRow{e.rel, fp})
	}
	return writeManifest(filepath.Join(root, ".bench", "link-manifest.tsv"), version, rows)
}

func installPlannedFile(root, mode string, e planEntry) error {
	dest := filepath.Join(root, e.rel)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	_ = os.Remove(dest)
	switch {
	case e.kind == "adapter":
		target, ok := AdapterTarget(e.rel)
		if !ok {
			return fmt.Errorf("adapter target unavailable for %s", e.rel)
		}
		return os.Symlink(target, dest)
	case mode == "symlink":
		return os.Symlink(e.src, dest)
	default:
		return copyFile(e.src, dest)
	}
}

func AdapterTarget(rel string) (string, bool) {
	var sourceRel string
	switch {
	case strings.HasPrefix(rel, ".claude/commands/"):
		sourceRel = ".agents/commands/" + strings.TrimPrefix(rel, ".claude/commands/")
	case strings.HasPrefix(rel, ".claude/skills/"):
		sourceRel = ".agents/skills/" + strings.TrimPrefix(rel, ".claude/skills/")
	default:
		return "", false
	}
	destDir := filepath.ToSlash(filepath.Dir(rel))
	parts := strings.Split(destDir, "/")
	ups := strings.Repeat("../", len(parts))
	return ups + sourceRel, true
}

func copyFile(src, dest string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, data, info.Mode().Perm())
}

func hooksDir(root string) string {
	out, err := git.Output("-C", root, "rev-parse", "--git-path", "hooks")
	if err != nil || out == "" {
		return filepath.Join(root, ".git", "hooks")
	}
	if filepath.IsAbs(out) {
		return out
	}
	return filepath.Join(root, out)
}

func installGitHook(root string, stderr io.Writer) error {
	if gitOK("-C", root, "remote", "get-url", "origin") &&
		!gitOK("-C", root, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD") {
		_ = exec.Command("git", "-C", root, "remote", "set-head", "origin", "--auto").Run()
	}
	def := git.DefaultBranch(root)
	hooks := hooksDir(root)
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}
	prepush := filepath.Join(hooks, "pre-push")
	if content, err := os.ReadFile(prepush); err == nil && !strings.Contains(string(content), "bench:managed-pre-push") {
		fmt.Fprintf(stderr, "conflict: %s exists and is not Bench-managed\n", prepush)
		return fmt.Errorf("foreign pre-push")
	}
	text := fmt.Sprintf(`#!/usr/bin/env bash
# bench:managed-pre-push
# Installed by 'bench link'. The merge is the human's; agents don't push %[1]s.
if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: pre-push\n'
  printf 'boundary: pre-push\n'
  printf 'denies: direct push to %[1]s\n'
  printf 'why: the merge belongs to the reviewer; agents open a PR instead of pushing %[1]s\n'
  exit 0
fi
while read -r local_ref local_oid remote_ref remote_oid; do
  if [[ "$remote_ref" == "refs/heads/%[1]s" ]]; then
    echo "blocked: direct push to %[1]s. Open a PR or merge it yourself." >&2
    exit 1
  fi
done
exit 0
`, def)
	return os.WriteFile(prepush, []byte(text), 0o755)
}

func gitOK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}
