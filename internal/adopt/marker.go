package adopt

import (
	"errors"
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
