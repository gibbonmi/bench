package adopt

import (
	"errors"
	"fmt"
	"github.com/gibbonmi/bench/internal/adopt/transaction"
	"os"
	"path/filepath"
	"strings"
)

const (
	benchStartMarker = "<!-- bench:start -->"
	benchEndMarker   = "<!-- bench:end -->"
)

// BenchAgentsBlock is the managed block `bench link` owns inside AGENTS.md.
func BenchAgentsBlock() string {
	return strings.Join([]string{
		benchStartMarker,
		"## Bench",
		"",
		"Bench is installed in this repo.",
		"",
		"- The gate is the oracle: run `bench gate`; done means it exits zero.",
		"- Full operating guide: `.bench/BENCH.md`.",
		"- Portable commands: `.agents/commands/`.",
		"- Portable skills: `.agents/skills/`.",
		"- Project profile: `projects/<name>.md` when present.",
		"- The reviewer owns merge decisions.",
		benchEndMarker,
		"",
	}, "\n")
}

type markerScan struct {
	starts, ends       []int
	unbalancedFence    bool
	benchTextInContent bool
}

func scanMarkers(content string) markerScan {
	lines := splitLines(content)
	inFence := false
	var scan markerScan
	for i, line := range lines {
		lineNo := i + 1
		lineIsFence := strings.HasPrefix(line, "```")
		if lineIsFence {
			inFence = !inFence
		}
		if strings.Contains(line, "bench:") {
			scan.benchTextInContent = true
		}
		if lineIsFence || inFence {
			continue
		}
		if strings.Contains(line, benchStartMarker) {
			scan.starts = append(scan.starts, lineNo)
		}
		if strings.Contains(line, benchEndMarker) {
			scan.ends = append(scan.ends, lineNo)
		}
	}
	scan.unbalancedFence = inFence
	return scan
}

func validateAgentsContent(content string) error {
	scan := scanMarkers(content)
	if scan.unbalancedFence && scan.benchTextInContent {
		return errors.New("conflict: AGENTS.md has an unclosed code fence around Bench markers; marker detection cannot be trusted")
	}
	if len(scan.starts) == 0 && len(scan.ends) == 0 {
		return nil
	}
	if len(scan.starts) == 1 && len(scan.ends) == 1 && scan.starts[0] < scan.ends[0] {
		return nil
	}
	return errors.New("conflict: AGENTS.md has malformed Bench managed block markers")
}

// RewriteAgentsBlock appends or replaces the unfenced Bench managed block while
// leaving fenced examples and project-owned text intact.
func RewriteAgentsBlock(content string) (string, error) {
	if err := validateAgentsContent(content); err != nil {
		return "", err
	}
	block := BenchAgentsBlock()
	scan := scanMarkers(content)
	if len(scan.starts) == 0 {
		if content == "" {
			return block, nil
		}
		return content + "\n" + block, nil
	}

	lines := splitLines(content)
	var out []string
	inFence := false
	skip := false
	done := false
	for _, line := range lines {
		lineIsFence := strings.HasPrefix(line, "```")
		if lineIsFence {
			inFence = !inFence
		}
		if !lineIsFence && !inFence && strings.Contains(line, benchStartMarker) {
			if !done {
				out = append(out, strings.Split(strings.TrimSuffix(block, "\n"), "\n")...)
				done = true
			}
			skip = true
			continue
		}
		if !lineIsFence && !inFence && strings.Contains(line, benchEndMarker) {
			skip = false
			continue
		}
		if !skip {
			out = append(out, line)
		}
	}
	result := strings.Join(out, "\n")
	if strings.HasSuffix(content, "\n") || !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result, nil
}

// StripAgentsBlock removes the unfenced Bench managed block from AGENTS.md while leaving
// fenced marker examples and project-owned prose intact. It is the reverse of
// RewriteAgentsBlock, sharing its fence-aware scan so the two agree on what "the managed
// block" is. Content with no managed block returns unchanged. unlink reads a resulting
// whitespace-only file as the signal to remove an AGENTS.md that link created with no
// user content.
func StripAgentsBlock(content string) (string, error) {
	if err := validateAgentsContent(content); err != nil {
		return "", err
	}
	scan := scanMarkers(content)
	if len(scan.starts) == 0 {
		return content, nil
	}
	lines := splitLines(content)
	var out []string
	inFence := false
	skip := false
	for _, line := range lines {
		lineIsFence := strings.HasPrefix(line, "```")
		if lineIsFence {
			inFence = !inFence
		}
		if !lineIsFence && !inFence && strings.Contains(line, benchStartMarker) {
			skip = true
			continue
		}
		if !lineIsFence && !inFence && strings.Contains(line, benchEndMarker) {
			skip = false
			continue
		}
		if !skip {
			out = append(out, line)
		}
	}
	result := strings.Join(out, "\n")
	if result != "" && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result, nil
}

func splitLines(content string) []string {
	if content == "" {
		return nil
	}
	content = strings.TrimSuffix(content, "\n")
	return strings.Split(content, "\n")
}

func benchClaudeMD() string {
	return "# Bench\n\nCanonical agreement in AGENTS.md; platform rules in .bench/BENCH.md.\n\n@AGENTS.md\n@.bench/BENCH.md\n"
}

func legacyClaudeMD() string {
	return "# Bench\n\nCanonical agreement in AGENTS.md.\n\n@AGENTS.md\n"
}

// claudeImportLines returns the marker-owned import lines every reclaimable CLAUDE.md
// form must carry. It derives them from legacyClaudeMD, the minimal form, rather than
// re-declaring the "@AGENTS.md" literal a second time. The canonical form, benchClaudeMD,
// is legacy plus one more import line, so the legacy set alone is the single check that
// accepts both forms.
func claudeImportLines() []string {
	var lines []string
	for _, line := range strings.Split(legacyClaudeMD(), "\n") {
		if strings.HasPrefix(line, "@") {
			lines = append(lines, line)
		}
	}
	return lines
}

// claudeHasImports reports whether content carries every marker-owned import line. It is
// true for both the canonical and legacy forms, and false for a preserved CLAUDE.md whose
// imports were stripped (doctor's row-11 red cell).
func claudeHasImports(content string) bool {
	for _, line := range claudeImportLines() {
		if !strings.Contains(content, line) {
			return false
		}
	}
	return true
}

// unlinkAgents is the AGENTS.md transform that unlink prepared: the staged stripped bytes,
// or the cause of a failed transform.
type unlinkAgents struct{ stage, failure string }

// observeUnlinkAgents records the facts of AGENTS.md for unlink and prepares the strip of
// the fenced Bench block, which keeps the user's surrounding prose. A strip that leaves the
// file whitespace-only, the case where link created it with no user content, prepares the
// removal of the file, mirroring link's create-if-absent symmetry. A special file or a
// malformed managed block fails the transform, so the file stays and the manifest survives
// for a manual fix. A real run stages the stripped bytes in stage. A dry run passes an
// empty stage and only fingerprints them. AGENTS.md is bespoke, not a manifest row.
func observeUnlinkAgents(root, stage string) (LifecycleObservation, unlinkAgents) {
	o := LifecycleObservation{Path: "AGENTS.md", PlanKind: bespokePlanKind}
	path := filepath.Join(root, o.Path)
	if isSpecialFile(path) {
		o.Present, o.Unprepared = true, true
		return o, unlinkAgents{failure: "not a regular file"}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return o, unlinkAgents{}
	}
	o.Present = true
	stripped, err := StripAgentsBlock(string(content))
	switch {
	case err != nil:
		o.Unprepared = true
		return o, unlinkAgents{failure: "its managed block could not be parsed"}
	case stripped == string(content):
		return o, unlinkAgents{}
	case strings.TrimSpace(stripped) == "":
		o.StagedKind = fileAbsent
		return o, unlinkAgents{}
	}
	o.Staged, o.StagedKind = hashBytes([]byte(stripped)), fileRegular
	if stage == "" {
		return o, unlinkAgents{}
	}
	info, err := os.Stat(path)
	if err != nil {
		o.Unprepared = true
		return o, unlinkAgents{failure: "mode unavailable"}
	}
	staged, err := stageBytes(stage, "agents", []byte(stripped), info.Mode().Perm())
	if err != nil {
		o.Unprepared = true
		return o, unlinkAgents{failure: "staging failed"}
	}
	return o, unlinkAgents{stage: staged}
}

// report is the unlink report line for the AGENTS.md decision d, empty when d leaves the
// file alone.
func (a unlinkAgents) report(d LifecycleDecision) string {
	switch {
	case d.Reason == reasonRefused:
		return "kept AGENTS.md (" + a.failure + ")"
	case d.Kind == lifecycleRemove:
		return "AGENTS.md (removed - no user prose remained)"
	case d.Kind == lifecycleChange:
		return "AGENTS.md (managed block stripped, prose kept)"
	}
	return ""
}

func stagedRepairAgents(stage, root string, repair *repairRun) (string, error) {
	if repair == nil {
		return stagedAgents(stage, root)
	}
	path := filepath.Join(root, "AGENTS.md")
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return stagedAgents(stage, root)
	}
	if err != nil {
		return "", err
	}
	if err := validateAgentsContent(string(content)); err != nil {
		return "", err
	}
	scan := scanMarkers(string(content))
	if len(scan.starts) > 0 {
		lines := strings.SplitAfter(string(content), "\n")
		start, end := 0, 0
		for i, line := range lines {
			if i < scan.starts[0]-1 {
				start += len(line)
			}
			if i < scan.ends[0] {
				end += len(line)
			}
		}
		if strings.TrimSuffix(string(content[start:end]), "\n") != strings.TrimSuffix(BenchAgentsBlock(), "\n") {
			return "", fmt.Errorf("conflict: AGENTS.md has a modified managed block")
		}
		return "", nil
	}
	next, err := RewriteAgentsBlock(string(content))
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	result, err := stageBytes(stage, "agents", []byte(next), info.Mode().Perm())
	if err != nil {
		return "", err
	}
	if repair.spans == nil {
		repair.spans = map[string]*transaction.Span{}
	}
	repair.spans[path] = &transaction.Span{Offset: len(content), AfterLength: len(next) - len(content)}
	return result, nil
}
