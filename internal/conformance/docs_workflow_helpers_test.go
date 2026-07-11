package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/coverage"
)

func checkWorkflowAnchors(root string) []string {
	var diags []string
	require := func(rel, needle string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !exists(path) {
			diags = append(diags, "acceptance coverage anchor file missing: "+rel)
			return
		}
		if !strings.Contains(readIfExists(path), needle) {
			diags = append(diags, fmt.Sprintf("%s missing acceptance coverage anchor: %s", rel, needle))
		}
	}

	require(".agents/commands/bench-write-spec.md", "acceptance coverage map")
	require(".agents/commands/bench-write-spec.md", "why it catches the failure")
	require(".agents/commands/bench-write-spec.md", "red signal")
	require(".agents/skills/bench-craft-spec/SKILL.md", "why it catches the failure")
	require(".agents/skills/bench-craft-spec/SKILL.md", "re-run idempotency")
	require(".agents/skills/bench-craft-spec/SKILL.md", "separate capability")
	require(".agents/commands/bench-write-spec.md", "bench-craft-spec")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "bench-craft-spec")
	require(".agents/skills/bench-craft-review/SKILL.md", "bench-craft-spec")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "acceptance row")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "not TDD-able")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "call count")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "row schema and the red-signal definition are")
	require(".agents/commands/bench-implement-spec.md", "coverage table")
	require(".agents/commands/bench-implement-spec.md", "already covered")
	require(".agents/commands/bench-implement-spec.md", "turning red-to-green")
	require(".agents/commands/bench-implement-spec.md", "bench coverage <spec>")
	require(".agents/commands/bench-review-implementation.md", "acceptance coverage map")
	require(".agents/commands/bench-review-implementation.md", "mapped behavior")
	require(".agents/commands/bench-review-implementation.md", "bench diff --full")
	require(".agents/commands/bench-review-implementation.md", "bench diff --full --commit")
	require(".agents/commands/bench-final-check.md", ".bench/gate.sh")
	require(".agents/commands/bench-final-check.md", "BENCH_GATE")
	require(".agents/commands/bench-write-spec.md", "seam diagram")
	require(".agents/commands/bench-write-spec.md", "tests attach here")
	require(".agents/commands/bench-write-spec.md", "edge inventory")
	require(".agents/commands/bench-write-spec.md", "Won't handle")
	require(".agents/commands/bench-write-spec.md", "hostile-input checklist")
	require(".agents/skills/bench-craft-tdd/SKILL.md", "floor, not the ceiling")
	require(".agents/skills/bench-craft-seams/SKILL.md", "failure modes")
	require(".agents/skills/bench-craft-seams/SKILL.md", "structure.budgets")
	require(".agents/commands/bench-review-implementation.md", "## Coverage")
	require(".agents/commands/bench-review-implementation.md", "Coverage axis")
	require(".agents/commands/bench-setup-repo.md", "hostile-input checklist")
	require("projects/benchkit.md", "hostile-input checklist")
	require(".agents/commands/bench-setup-repo.md", "craft-gate")
	require(".agents/commands/bench-final-check.md", "craft-gate")
	require(".agents/commands/bench-review-implementation.md", "craft-review")
	require(".agents/skills/bench-craft-review/SKILL.md", "an edge nobody decided")
	require(".agents/commands/bench-review-implementation.md", "craft-delegate")
	require(".agents/skills/bench-craft-delegate/SKILL.md", "a claim, not a result")
	require(".agents/commands/bench-implement-spec.md", "When the build stops short")
	require(".agents/commands/bench-write-spec.md", "Superseded by")
	require(".agents/commands/bench-debug.md", "before launching the shift")
	require(".agents/commands/bench-shape-idea.md", "## Handoff")
	require(".agents/commands/bench-shape-idea.md", "Hostile-input owner")
	require(".agents/commands/bench-shape-idea.md", "Dependency order")
	require(".agents/commands/bench-shape-idea.md", "n/a \u2014")
	require(".agents/commands/bench-write-spec.md", "map's Handoff")
	require(".agents/commands/bench-write-spec.md", "spec-retire:")
	require(".agents/commands/bench-write-spec.md", "Status: staged")
	require(".agents/commands/bench-write-spec.md", "new session on the mid tier")
	require(".agents/commands/bench-implement-spec.md", "Status: implemented")
	require(".agents/commands/bench-debug.md", "diff-filter=D")
	require(".agents/commands/bench-review-implementation.md", "reviews/<spec-slug>.md")
	require(".agents/commands/bench-review-implementation.md", "same session that writes it")
	require(".agents/commands/bench-implement-spec.md", "reviews/<spec-slug>.md")
	require(".agents/commands/bench-implement-spec.md", "names and deletes the file")
	require(".agents/commands/bench-final-check.md", "not outlive the decision it captured")
	require(".agents/commands/bench-implement-spec.md", "bench commit -m")
	require(".agents/commands/bench-final-check.md", "bench commit -m")
	require(".agents/commands/bench-review-implementation.md", "actionable findings")
	require(".agents/commands/bench-review-implementation.md", "writes no artifact")
	require(".agents/commands/bench-review-implementation.md", "same green fix commit")
	require(".agents/commands/bench-review-implementation.md", "section per axis")
	require(".agents/commands/bench-review-implementation.md", "finding count, its worst issue")
	require(".agents/commands/bench-review-implementation.md", "doc citation its axis supplied")
	require(".agents/commands/bench-what-next.md", "Reconcile first")
	require(".agents/commands/bench-what-next.md", "through the accused command")
	require(".agents/commands/bench-debug.md", "through the accused command")
	require(".agents/commands/bench-what-next.md", "empties to zero")
	require(".agents/commands/bench-what-next.md", "verdict in the batch diff")
	require(".agents/commands/bench-what-next.md", "one uncommitted batch diff")
	require(".agents/commands/bench-what-next.md", "commit on green")
	require(".agents/commands/bench-what-next.md", "## Recommended sequence")

	require(".agents/commands/bench-assess.md", "verify the previous assessment's backlog landed")
	require(".agents/commands/bench-assess.md", "read-only area sweeps on the mid tier")
	require(".agents/commands/bench-assess.md", "synthesize adversarially on the top tier")
	require(".agents/commands/bench-assess.md", "replaces its predecessor")
	require(".agents/commands/bench-assess.md", "/bench-what-next")

	requireCollapsed := func(rel, needle, diag string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !exists(path) {
			diags = append(diags, "acceptance coverage anchor file missing: "+rel)
			return
		}
		if !strings.Contains(collapseSpace(readIfExists(path)), needle) {
			diags = append(diags, diag)
		}
	}
	requireCollapsed(".agents/commands/bench-implement-spec.md", "apply `craft-seams`' split-or-grant rule",
		".agents/commands/bench-implement-spec.md dropped the craft-seams split-or-grant pointer")
	requireCollapsed(".agents/skills/bench-craft-seams/SKILL.md", "check both the file-length budget and the directory's file-count headroom",
		".agents/skills/bench-craft-seams/SKILL.md dropped the structure split-vs-grant headroom rule")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "pins every file-tool path to that root",
		".agents/skills/bench-craft-delegate/SKILL.md dropped the shared-worktree file-tool path pin")
	requireCollapsed(".agents/skills/bench-craft-delegate/SKILL.md", "names a commit-specific sentinel",
		".agents/skills/bench-craft-delegate/SKILL.md dropped the fix-pass snapshot sentinel precondition")
	requireCollapsed(".agents/commands/bench-write-spec.md", "stale-command-reference sweep remains fail-closed across staged specs",
		".agents/commands/bench-write-spec.md dropped the staged-spec fail-closed command sweep posture")
	requireCollapsed(".agents/commands/bench-what-next.md", "use `bench spec history <slug>` for the shipped-row check",
		".agents/commands/bench-what-next.md dropped the bench spec history shipped-row check")
	requireCollapsed(".bench/BENCH.md", "Parked ideas land in `IDEAS.md`",
		".bench/BENCH.md Capture section does not name IDEAS.md as the capture sink")
	requireCollapsed(".bench/BENCH.md", "append the dated line (`- YYYY-MM-DD <text>`) to `IDEAS.md`",
		".bench/BENCH.md Capture section lost the no-PATH fallback append to IDEAS.md")
	requireCollapsed(".agents/commands/bench-write-spec.md", "promote-then-delete commit removes the spec's `ROADMAP.md` row",
		".agents/commands/bench-write-spec.md does not remove the spec's ROADMAP.md row in the promote-then-delete commit (row presence is status)")
	requireCollapsed(".agents/commands/bench-shape-idea.md", "never pause for permission or a re-prompt",
		".agents/commands/bench-shape-idea.md dropped the resume-mode grill continuation rule; a running grill carries into newly-unblocked tickets without pausing for a re-prompt")
	requireCollapsed(".agents/commands/bench-review-implementation.md", "Integrate the findings accepted for this round, run focused checks for the changed behavior, then run one final gate and stop. Open another semantic review round only when that gate fails or the reviewer requests one.",
		".agents/commands/bench-review-implementation.md dropped the terminal repair-pass bound")

	shapeIdeaPath := filepath.Join(root, ".agents", "commands", "bench-shape-idea.md")
	if exists(shapeIdeaPath) && strings.Contains(collapseSpace(readIfExists(shapeIdeaPath)), "straight to `/bench-write-spec`") {
		diags = append(diags, ".agents/commands/bench-shape-idea.md reintroduces the skip-to-spec bypass fragment; every idea must yield a map with a Handoff before a spec")
	}
	writeSpecPath := filepath.Join(root, ".agents", "commands", "bench-write-spec.md")
	if exists(writeSpecPath) && !strings.Contains(collapseSpace(readIfExists(writeSpecPath)), "refuses to run without") {
		diags = append(diags, ".agents/commands/bench-write-spec.md dropped the map-required entry contract (refuses to run without a complete map)")
	}

	readme := readIfExists(filepath.Join(root, "README.md"))
	if readme != "" {
		if !strings.Contains(readme, "session-start.sh") {
			diags = append(diags, "README layout omits .bench/hooks/session-start.sh")
		}
		if !strings.Contains(readme, "bench.sh") {
			diags = append(diags, "README layout omits the real bin/bench.sh filename")
		}
		if !strings.Contains(readme, "benchkit.md") {
			diags = append(diags, "README layout omits projects/benchkit.md")
		}
		if strings.Contains(readme, "\u2502   \u2514\u2500\u2500 bench                 #") {
			diags = append(diags, "README layout still names bin/bench instead of bin/bench.sh")
		}
	}

	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-implement-spec.md")); text != "" && !strings.Contains(text, "craft-line") {
		diags = append(diags, "bench-implement-spec.md does not reference craft-line")
	}
	if text := readIfExists(filepath.Join(root, ".agents", "commands", "bench-write-spec.md")); text != "" {
		if !strings.Contains(text, "craft-line") {
			diags = append(diags, "bench-write-spec.md does not reference craft-line")
		}
		if !strings.Contains(text, "model and effort") {
			diags = append(diags, "bench-write-spec.md does not mandate per-story model and effort")
		}
	}
	if text := readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md")); text != "" && !strings.Contains(text, "BENCH_MODEL") {
		diags = append(diags, "BENCH-reference.md adapter contract does not document BENCH_MODEL")
	}
	return diags
}

func checkSkillsIndexGenerateVerify(root, kitRoot string) []string {
	if !exists(filepath.Join(kitRoot, ".bench", "skills-index.sh")) {
		return nil
	}
	tmp, err := os.MkdirTemp("", "bench-skills-index-*")
	if err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	defer os.RemoveAll(tmp)
	if err := os.MkdirAll(filepath.Join(tmp, ".bench"), 0o755); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	if err := os.MkdirAll(filepath.Join(tmp, ".agents", "skills", "zeta-skill"), 0o755); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	if err := os.WriteFile(filepath.Join(tmp, ".agents", "skills", "zeta-skill", "SKILL.md"), []byte("---\nname: zeta-skill\ndescription: d\nindex: doing zeta things\n---\n"), 0o644); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	if err := os.WriteFile(filepath.Join(tmp, ".bench", "BENCH-reference.md"), []byte("# Reference\n\n<!-- bench:skills-index:start -->\n<!-- bench:skills-index:end -->\n"), 0o644); err != nil {
		return []string{"skills-index generate/verify contract setup failed: " + err.Error()}
	}
	script := filepath.Join(kitRoot, ".bench", "skills-index.sh")
	if probe := runAt(tmp, "bash", script, "--check"); probe == nil || probe.ExitCode == 0 {
		return []string{"skills-index generate/verify contract failed: check passed on an empty index block"}
	}
	if probe := runAt(tmp, "bash", script, "--write"); probe == nil || probe.ExitCode != 0 {
		return []string{"skills-index generate/verify contract failed: --write failed"}
	}
	generated := readIfExists(filepath.Join(tmp, ".bench", "BENCH-reference.md"))
	if !strings.Contains(generated, "- doing zeta things \u2192 `.agents/skills/zeta-skill/SKILL.md`") {
		return []string{"skills-index generate/verify contract failed: --write did not generate the entry from frontmatter"}
	}
	if probe := runAt(tmp, "bash", script, "--check"); probe == nil || probe.ExitCode != 0 {
		return []string{"skills-index generate/verify contract failed: check red right after --write"}
	}
	before := readIfExists(filepath.Join(tmp, ".bench", "BENCH-reference.md"))
	if probe := runAt(tmp, "bash", script, "--write"); probe == nil || probe.ExitCode != 0 {
		return []string{"skills-index generate/verify contract failed: second --write failed"}
	}
	if before != readIfExists(filepath.Join(tmp, ".bench", "BENCH-reference.md")) {
		return []string{"skills-index generate/verify contract failed: --write is not idempotent"}
	}
	return nil
}

func checkCoverageMaps(root string) []string {
	specsDir := filepath.Join(root, "specs")
	if !exists(specsDir) {
		return nil
	}
	matches, _ := filepath.Glob(filepath.Join(specsDir, "*.md"))
	sort.Strings(matches)
	var diags []string
	for _, path := range matches {
		out, code := coverage.Command([]string{"--check", path})
		if code == 0 {
			continue
		}
		if strings.TrimSpace(out) == "" {
			diags = append(diags, fmt.Sprintf("%s coverage --check failed (exit %d) with no message", slashRel(root, path), code))
			continue
		}
		for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			diags = append(diags, strings.TrimPrefix(line, "error: "))
		}
	}
	return diags
}

func collapseSpace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}
