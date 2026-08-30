package conformance

// This file grades the profile's fast-lane advertisement. The kit's built-in lane is
// the single source for what a worktree commit runs. The profile table is an
// advertisement of that source, so a drifted row reds rather than misleads a reader.

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

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
	rows, diags := parseProfileLaneTable(profile)
	if len(diags) != 0 {
		return diags
	}
	want, err := builtInLaneRows(root)
	if err != "" {
		return []string{err}
	}
	names, wantNames := laneRowNames(rows), laneRowNames(want)
	if !slices.Equal(names, wantNames) {
		return append(diags, fmt.Sprintf("profile lane rows = %v, want %v", names, wantNames))
	}
	for i, row := range rows {
		if row.argv != want[i].argv {
			diags = append(diags, fmt.Sprintf("profile lane row stale: %s renders %q, the kit lane runs %q", row.name, row.argv, want[i].argv))
		}
		if canonicalLaneClasses(row.selectedBy) != canonicalLaneClasses(want[i].selectedBy) {
			diags = append(diags, fmt.Sprintf("profile lane row stale: %s is selected by %q, the kit class table selects it by %q", row.name, row.selectedBy, want[i].selectedBy))
		}
	}
	return diags
}

// profileLaneRow is one rendered row of the lane table: the check, the argv the reader
// is promised, and the classes that select the check.
type profileLaneRow struct {
	name       string
	argv       string
	selectedBy string
}

func laneRowNames(rows []profileLaneRow) []string {
	names := make([]string, len(rows))
	for i, row := range rows {
		names[i] = row.name
	}
	return names
}

// canonicalLaneClasses reduces a `selected by` cell to its class list, so the comparison
// grades the names and not the spacing or the backticks around them.
func canonicalLaneClasses(cell string) string {
	var names []string
	for _, name := range strings.Split(cell, ",") {
		if name = strings.Trim(strings.TrimSpace(name), "`"); name != "" {
			names = append(names, name)
		}
	}
	return strings.Join(names, ",")
}

// laneSelectedBy renders the classes that select one check, in the class table's own
// order. The kit table is the one source, so a profile cell can only advertise it.
func laneSelectedBy(name string) string {
	var names []string
	for _, class := range gate.LaneClasses() {
		if slices.Contains(class.Checks, name) {
			names = append(names, class.Name)
		}
	}
	return strings.Join(names, ", ")
}

// builtInLaneRows renders the kit's lane as the profile spells it. The run-binary
// placeholder is read back from the phase table, so this check never restates a token
// the gate owns.
func builtInLaneRows(root string) (rows []profileLaneRow, diagnostic string) {
	token := runBinaryToken(root)
	if token == "" {
		return nil, "kit phase table declares no Bench-owned phase, so the lane run binary is underivable"
	}
	for _, check := range gate.BenchkitLane(laneProfileRootToken, root) {
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
		rows = append(rows, profileLaneRow{
			name:       check.Name,
			argv:       strings.Join(rendered, " "),
			selectedBy: laneSelectedBy(check.Name),
		})
	}
	return rows, ""
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
func parseProfileLaneTable(profile string) (rows []profileLaneRow, diags []string) {
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
		selectedBy := ""
		if len(cells) > 2 {
			selectedBy = strings.TrimSpace(cells[2])
		}
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
		rows = append(rows, profileLaneRow{name: name, argv: rendered, selectedBy: selectedBy})
	}
	if !found {
		return nil, []string{"profile lane table missing"}
	}
	return rows, diags
}

// TestProfileLaneTableRedsAStaleSelectedByCell is PL36. The fixture profile renders the
// kit's own rows, so the one planted defect is the `gofmt` row's selectors. The check
// therefore reds exactly once, and it names the check the reader must correct.
func TestProfileLaneTableRedsAStaleSelectedByCell(t *testing.T) {
	root := t.TempDir()
	writeProfileFixture(t, root, "go.mod", "module example.com/x\n\ngo 1.24\n")
	rows, diagnostic := builtInLaneRows(root)
	if diagnostic != "" {
		t.Fatal(diagnostic)
	}
	table := []string{"| check | authoritative argv | selected by |", "|---|---|---|"}
	for _, row := range rows {
		selectedBy := row.selectedBy
		if row.name == "gofmt" {
			selectedBy = "markdown"
		}
		table = append(table, fmt.Sprintf("| `%s` | `%s` | %s |", row.name, row.argv, selectedBy))
	}
	writeProfileFixture(t, root, filepath.Join("projects", "benchkit.md"),
		"# Profile\n\n"+strings.Join(table, "\n")+"\n")

	diags := checkProfileLaneTable(root)
	if len(diags) != 1 {
		t.Fatalf("diagnostics = %v, want exactly one", diags)
	}
	if !strings.Contains(diags[0], "gofmt") {
		t.Errorf("diagnostic = %q, want the stale row's check named", diags[0])
	}
}

func writeProfileFixture(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
