package conformance

// This file grades the profile's fast-lane advertisement. The kit's built-in lane is
// the single source for what a worktree commit runs. The profile table is an
// advertisement of that source, so a drifted row reds rather than misleads a reader.

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
)

// laneProfileRootToken is how the profile spells the graded root inside an argv cell,
// and laneProfileBinary is how it spells the selected Bench executable. Both are
// rendering choices this check owns, because the built-in argv carries placeholders no
// reader would recognize.
const (
	laneProfileRootToken   = "<root>"
	laneProfileBinary      = "bench"
	laneProfileMarkdown    = "<named Markdown>"
	laneProfileTableHeader = "check"
)

func checkProfileLaneTable(root string) []string {
	profile := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	if profile == "" {
		return nil
	}
	names, argv, diags := parseProfileLaneTable(profile)
	if len(diags) != 0 {
		return diags
	}
	wantNames, wantArgv, err := builtInLaneRows(root)
	if err != "" {
		return []string{err}
	}
	if !slices.Equal(names, wantNames) {
		return append(diags, fmt.Sprintf("profile lane rows = %v, want %v", names, wantNames))
	}
	for i, name := range names {
		if argv[i] != wantArgv[i] {
			diags = append(diags, fmt.Sprintf("profile lane row stale: %s renders %q, the kit lane runs %q", name, argv[i], wantArgv[i]))
		}
	}
	return diags
}

// builtInLaneRows renders the kit's lane as the profile spells it. The run-binary
// placeholder is read back from the phase table, so this check never restates a token
// the gate owns.
func builtInLaneRows(root string) (names, argv []string, diagnostic string) {
	token := runBinaryToken(root)
	if token == "" {
		return nil, nil, "kit phase table declares no Bench-owned phase, so the lane run binary is underivable"
	}
	for _, check := range gate.BenchkitLane(laneProfileRootToken, root) {
		names = append(names, check.Name)
		rendered := make([]string, 0, len(check.Argv))
		for _, arg := range check.Argv {
			switch arg {
			case token:
				rendered = append(rendered, laneProfileBinary)
			case gate.LaneNamedMarkdownToken:
				rendered = append(rendered, laneProfileMarkdown)
			default:
				rendered = append(rendered, arg)
			}
		}
		argv = append(argv, strings.Join(rendered, " "))
	}
	return names, argv, ""
}

// runBinaryToken answers the placeholder the gate substitutes with the selected Bench
// executable. It comes from the phase table's own first Bench-owned phase.
func runBinaryToken(root string) string {
	for _, phase := range gate.BenchkitPhases(root, root) {
		if len(phase.Argv) != 0 && strings.HasPrefix(phase.Argv[0], "<") {
			return phase.Argv[0]
		}
	}
	return ""
}

// parseProfileLaneTable reads the lane table's rendered rows in document order. It
// finds the table by its header cells, so the neighbouring phase table, whose argv
// column is spelled the same way, is never read in its place.
func parseProfileLaneTable(profile string) (names, argv, diags []string) {
	found := false
	seen := map[string]bool{}
	for _, line := range strings.Split(profile, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			if found {
				break
			}
			continue
		}
		cells := strings.Split(strings.Trim(trimmed, "|"), "|")
		if len(cells) < 2 {
			continue
		}
		name := strings.Trim(strings.TrimSpace(cells[0]), "`")
		rendered := strings.Trim(strings.TrimSpace(cells[1]), "`")
		if !found {
			found = name == laneProfileTableHeader && rendered == "authoritative argv"
			continue
		}
		if strings.Trim(name, "-:") == "" {
			continue
		}
		if seen[name] {
			diags = append(diags, "profile lane row duplicated: "+name)
			continue
		}
		seen[name] = true
		names = append(names, name)
		argv = append(argv, rendered)
	}
	if !found {
		return nil, nil, []string{"profile lane table missing"}
	}
	return names, argv, diags
}
