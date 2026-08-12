package conformance

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
)

const (
	// proseBudgetProfile holds the reviewer-owned budget table. Raising a budget is an
	// edit to this file and nowhere else.
	proseBudgetProfile = "projects/benchkit.md"
	// proseBudgetSection anchors the parse, so the profile's other tables can never be
	// read as budget policy.
	proseBudgetSection = "Guidance prose budgets"
	// proseBudgetSkillsDir is the guidance tree whose SKILL.md files join the universe
	// structurally. Enumeration cannot come from the table's own glob rows: a table that
	// lost its default row would then enumerate nothing, and the skill somebody adds next
	// month would pass unbudgeted instead of being reported as unclassified.
	proseBudgetSkillsDir = ".agents/skills"
	// proseBudgetSkillFile is the one file per skill directory the budget grades; a
	// skill's references and assets are outside the reviewed universe.
	proseBudgetSkillFile = "SKILL.md"
)

// proseBudgetPolicy is the profile's table: exact subject rows, and glob rows that classify
// whatever the enumeration finds.
type proseBudgetPolicy struct {
	exact    map[string]int
	patterns []proseBudgetPattern
}

type proseBudgetPattern struct {
	glob  string
	limit int
}

// limitFor resolves one subject's budget. Exact rows are consulted before glob rows, so the
// reviewer's specific number always beats the default the same subject also matches.
func (p proseBudgetPolicy) limitFor(rel string) (int, bool) {
	if limit, found := p.exact[rel]; found {
		return limit, true
	}
	for _, pattern := range p.patterns {
		if matched, err := path.Match(pattern.glob, rel); err == nil && matched {
			return pattern.limit, true
		}
	}
	return 0, false
}

// checkGuidanceProseBudgets grades every guidance file the profile budgets against the
// profile's own numbers.
func checkGuidanceProseBudgets(root string) []string {
	policy, diags := parseProseBudgetPolicy(readIfExists(filepath.Join(root, filepath.FromSlash(proseBudgetProfile))))
	if len(diags) != 0 {
		return diags
	}
	subjects, diags := proseBudgetSubjects(root, policy)
	for _, rel := range subjects {
		absolute := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Lstat(absolute)
		switch {
		case err != nil:
			diags = append(diags, "prose-budget subject missing: "+rel+" is named by the "+proseBudgetProfile+" budget table but absent from the tree")
			continue
		case info.Mode()&os.ModeSymlink != 0:
			diags = append(diags, "prose-budget subject refused: "+rel+" is a symbolic link, not a regular file")
			continue
		case !info.Mode().IsRegular():
			diags = append(diags, "prose-budget subject refused: "+rel+" is a named pipe or other special file, not a regular file")
			continue
		}
		limit, classified := policy.limitFor(rel)
		if !classified {
			diags = append(diags, "prose-budget subject unclassified: "+rel+" matches no row in the "+proseBudgetProfile+" budget table")
			continue
		}
		body, err := os.ReadFile(absolute)
		if err != nil {
			diags = append(diags, "prose-budget subject unreadable: "+rel+" could not be read ("+err.Error()+")")
			continue
		}
		if count := proseBudgetLineCount(body); count > limit {
			diags = append(diags, fmt.Sprintf("prose-budget exceeded: %s is %d lines, over its %d-line budget", rel, count, limit))
		}
	}
	return diags
}

// parseProseBudgetPolicy reads the reviewer's table out of the profile. Any policy fault
// returns diagnostics and no policy: a table nobody can read cannot say which subject is
// over, so grading a partial parse would report a clean tree for a broken policy.
func parseProseBudgetPolicy(profile string) (proseBudgetPolicy, []string) {
	section, anchored := profileSection(profile, proseBudgetSection)
	if !anchored {
		return proseBudgetPolicy{}, []string{"prose-budget policy missing: " + proseBudgetProfile + " has no '" + proseBudgetSection + "' heading holding the budget table"}
	}
	policy := proseBudgetPolicy{exact: map[string]int{}}
	seen := map[string]bool{}
	var diags []string
	header := false
	for _, line := range strings.Split(section, "\n") {
		row, isRow := markdownRow(line)
		if !isRow {
			// The table ends at the first line that is not a row; prose after it belongs
			// to the section, not to the policy.
			if header {
				break
			}
			continue
		}
		if len(row) < 2 {
			continue
		}
		if !header {
			header = strings.ToLower(row[0]) == "subject" && strings.ToLower(row[1]) == "limit"
			continue
		}
		if isRuleRow(row) || row[0] == "" {
			continue
		}
		subject, cell := row[0], row[1]
		if seen[subject] {
			diags = append(diags, "prose-budget policy duplicate row: "+proseBudgetProfile+" names subject '"+subject+"' more than once")
			continue
		}
		seen[subject] = true
		limit, err := strconv.Atoi(cell)
		if err != nil || limit <= 0 {
			diags = append(diags, "prose-budget policy invalid limit: "+proseBudgetProfile+" gives subject '"+subject+"' the limit '"+cell+"', which is not a positive line count")
			continue
		}
		if !strings.ContainsAny(subject, "*?[") {
			policy.exact[subject] = limit
			continue
		}
		if _, err := path.Match(subject, ""); err != nil {
			diags = append(diags, "prose-budget policy malformed: "+proseBudgetProfile+" gives subject '"+subject+"', which is not a valid path pattern")
			continue
		}
		policy.patterns = append(policy.patterns, proseBudgetPattern{glob: subject, limit: limit})
	}
	switch {
	case !header:
		return proseBudgetPolicy{}, []string{"prose-budget policy malformed: " + proseBudgetProfile + " renders no '| subject | limit |' table under its '" + proseBudgetSection + "' heading"}
	case len(diags) != 0:
		return proseBudgetPolicy{}, diags
	case len(seen) == 0:
		return proseBudgetPolicy{}, []string{"prose-budget policy missing: " + proseBudgetProfile + "'s budget table names no subject"}
	}
	return policy, nil
}

// proseBudgetSubjects lists the canonical universe in stable order — every exact subject the
// table names, plus every skill's SKILL.md — along with one diagnostic per skill directory
// that is a symbolic link. That entry is refused rather than descended into because
// `.claude/skills` is a tree of links to these very files, and following one would pull an
// adapter surface into the universe under a canonical path.
func proseBudgetSubjects(root string, policy proseBudgetPolicy) (subjects, diags []string) {
	seen := map[string]bool{}
	add := func(rel string) {
		if !seen[rel] {
			seen[rel] = true
			subjects = append(subjects, rel)
		}
	}
	for subject := range policy.exact {
		add(subject)
	}
	// os.ReadDir sorts by filename, so the enumeration reports in one order run to run.
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(proseBudgetSkillsDir)))
	if err == nil {
		for _, entry := range entries {
			rel := path.Join(proseBudgetSkillsDir, entry.Name())
			if entry.Type()&os.ModeSymlink != 0 {
				diags = append(diags, "prose-budget subject refused: "+rel+" is a symbolic link, not a regular directory")
				continue
			}
			if !entry.IsDir() {
				continue
			}
			skill := path.Join(rel, proseBudgetSkillFile)
			// A skill directory with no SKILL.md at all is some other check's fact; only
			// a subject the table names is required to exist.
			if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(skill))); err == nil {
				add(skill)
			}
		}
	}
	sort.Strings(subjects)
	return subjects, diags
}

// proseBudgetLineCount counts lines the way a reader does, so a file ending with a newline
// and the same file without one carry the same count.
func proseBudgetLineCount(body []byte) int {
	if len(body) == 0 {
		return 0
	}
	return bytes.Count(bytes.TrimSuffix(body, []byte("\n")), []byte("\n")) + 1
}

// proseBudgetTable renders the profile subsection the checker parses, taking its rows
// verbatim so a case can corrupt the header or a cell without a second table author.
func proseBudgetTable(header string, rows ...string) string {
	body := "# benchkit\n\n## Gate\n\n### " + proseBudgetSection + "\n\n" + header + "\n|---|---|\n"
	for _, row := range rows {
		body += row + "\n"
	}
	return body + "\n## Notes for cold sessions\n\ntail\n"
}

// proseBudgetHeader is the header row the parser anchors the table on.
const proseBudgetHeader = "| subject | limit |"

// proseBudgetRows is a healthy policy: one exact row outside the skills tree, one exact
// row inside it, and the all-skills default the exact row has to beat.
var proseBudgetRows = []string{
	"| `.bench/BENCH.md` | 150 |",
	"| `.agents/skills/bench-craft-tickets/SKILL.md` | 100 |",
	"| `.agents/skills/*/SKILL.md` | 120 |",
}

// lines renders count logical lines, ending with a newline unless trailing is false.
func proseBudgetLines(count int, trailing bool) string {
	if count == 0 {
		return ""
	}
	body := strings.Repeat("prose\n", count)
	if !trailing {
		body = strings.TrimSuffix(body, "\n")
	}
	return body
}

// writeProseBudgetRoot plants a profile and the named guidance files in a throwaway root.
func writeProseBudgetRoot(t *testing.T, profile string, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	write := func(rel, content string) {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(proseBudgetProfile, profile)
	for rel, content := range files {
		write(rel, content)
	}
	return root
}

// healthyProseBudgetFiles renders every subject proseBudgetRows names, each one line under
// its limit, so a case can push exactly the subject it is about over the edge.
func healthyProseBudgetFiles() map[string]string {
	return map[string]string{
		".bench/BENCH.md": proseBudgetLines(149, true),
		".agents/skills/bench-craft-tickets/SKILL.md": proseBudgetLines(99, true),
		".agents/skills/bench-craft-line/SKILL.md":    proseBudgetLines(119, true),
	}
}

// TestGuidanceProseBudgetsComeFromTheProfileTable is the single-source row: the limit the
// checker enforces moves when the reviewer moves the cell, which a hard-coded constant
// cannot do. The exact craft-tickets row also has to beat the all-skills default it
// matches, so the two are checked in the same parse.
func TestGuidanceProseBudgetsComeFromTheProfileTable(t *testing.T) {
	files := healthyProseBudgetFiles()
	healthy := proseBudgetTable(proseBudgetHeader, proseBudgetRows...)
	if diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, healthy, files)); len(diags) != 0 {
		t.Fatalf("a tree inside every budget got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	lowered := strings.Replace(healthy, "SKILL.md` | 100 |", "SKILL.md` | 90 |", 1)
	diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, lowered, files))
	if !containsDiagnostic(diags, "prose-budget exceeded: .agents/skills/bench-craft-tickets/SKILL.md is 99 lines, over its 90-line budget") {
		t.Fatalf("lowering the craft-tickets cell to 90 did not move the enforced limit:\n%s", strings.Join(diags, "\n"))
	}

	// The default row would allow 120; the exact row is what the subject is graded by.
	over := healthyProseBudgetFiles()
	over[".agents/skills/bench-craft-tickets/SKILL.md"] = proseBudgetLines(110, true)
	diags = checkGuidanceProseBudgets(writeProseBudgetRoot(t, healthy, over))
	if !containsDiagnostic(diags, "over its 100-line budget") {
		t.Fatalf("the all-skills default overrode the exact craft-tickets row:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetPolicyFaultsFailClosed covers the four ways the reviewer's table
// can stop being readable. Each one stops the run rather than grading a partial parse: a
// policy nobody can read cannot say which subject is over, and a clean verdict from a
// broken table is the green-by-omission this check exists to prevent.
func TestGuidanceProseBudgetPolicyFaultsFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name, profile, want string
	}{
		{
			"malformed header",
			proseBudgetTable("| path | maximum |", proseBudgetRows...),
			"prose-budget policy malformed",
		},
		{
			"missing section",
			"# benchkit\n\n## Gate\n\nNo budget table lives here.\n",
			"prose-budget policy missing",
		},
		{
			"duplicate row",
			proseBudgetTable(proseBudgetHeader, append(proseBudgetRows, "| `.bench/BENCH.md` | 200 |")...),
			"prose-budget policy duplicate row: projects/benchkit.md names subject '.bench/BENCH.md' more than once",
		},
		{
			"invalid limit",
			proseBudgetTable(proseBudgetHeader, "| `.bench/BENCH.md` | abc |"),
			"prose-budget policy invalid limit: projects/benchkit.md gives subject '.bench/BENCH.md' the limit 'abc'",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, tc.profile, healthyProseBudgetFiles()))
			if !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q, got:\n%s", tc.want, strings.Join(diags, "\n"))
			}
		})
	}
}

// TestGuidanceProseBudgetClassifiesEverySkillItFinds pins the enumeration to the guidance
// tree rather than to the table's own rows. A skill added next month is graded by the
// default row without anybody editing the checker; with no default row it is reported
// unclassified, which is the answer a table-driven enumeration structurally cannot give —
// it would simply never look at the new directory.
func TestGuidanceProseBudgetClassifiesEverySkillItFinds(t *testing.T) {
	files := healthyProseBudgetFiles()
	files[".agents/skills/bench-craft-new/SKILL.md"] = proseBudgetLines(121, true)
	withDefault := checkGuidanceProseBudgets(writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), files))
	if !containsDiagnostic(withDefault, "prose-budget exceeded: .agents/skills/bench-craft-new/SKILL.md is 121 lines, over its 120-line budget") {
		t.Fatalf("a newly added skill was not auto-classified under the default row:\n%s", strings.Join(withDefault, "\n"))
	}

	noDefault := proseBudgetTable(proseBudgetHeader, proseBudgetRows[0], proseBudgetRows[1])
	diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, noDefault, files))
	if !containsDiagnostic(diags, "prose-budget subject unclassified: .agents/skills/bench-craft-new/SKILL.md matches no row") {
		t.Fatalf("a skill no row classifies was accepted:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetRefusesNonRegularSubjects covers the entry kinds the canonical
// universe can hold besides a file. Classification precedes every read, so the FIFO case
// fails by expiring its deadline rather than by a wrong answer when an implementation opens
// first, and the symlink case points at a subject far over budget: a checker that followed
// it would report the target's line count under the link's path.
func TestGuidanceProseBudgetRefusesNonRegularSubjects(t *testing.T) {
	for _, tc := range []struct {
		kind, want string
		plant      func(t *testing.T, path string)
	}{
		{
			kind: "symlink",
			want: "prose-budget subject refused: .agents/skills/bench-craft-linked/SKILL.md is a symbolic link, not a regular file",
			plant: func(t *testing.T, path string) {
				target := filepath.Join(filepath.Dir(path), "..", "bench-craft-line", "SKILL.md")
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
		},
		{
			kind: "fifo",
			want: "prose-budget subject refused: .agents/skills/bench-craft-linked/SKILL.md is a named pipe or other special file, not a regular file",
			plant: func(t *testing.T, path string) {
				if err := syscall.Mkfifo(path, 0o644); err != nil {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
				}
			},
		},
		{
			kind: "socket",
			want: "prose-budget subject refused: .agents/skills/bench-craft-linked/SKILL.md is a named pipe or other special file, not a regular file",
			plant: func(t *testing.T, path string) {
				listener, err := net.Listen("unix", path)
				if err != nil {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("unix sockets unavailable on this filesystem: %v", err))
				}
				t.Cleanup(func() { listener.Close() })
			},
		},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			files := healthyProseBudgetFiles()
			files[".agents/skills/bench-craft-line/SKILL.md"] = proseBudgetLines(119, true)
			root := writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), files)
			dir := filepath.Join(root, ".agents", "skills", "bench-craft-linked")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			tc.plant(t, filepath.Join(dir, "SKILL.md"))
			done := make(chan []string, 1)
			go func() { done <- checkGuidanceProseBudgets(root) }()
			select {
			case diags := <-done:
				if !containsDiagnostic(diags, tc.want) {
					t.Fatalf("a %s subject was not refused:\n%s", tc.kind, strings.Join(diags, "\n"))
				}
			case <-time.After(bounds.TestDeadline(0)):
				t.Fatalf("the check blocked on a %s, so it opened the subject before classifying it", tc.kind)
			}
		})
	}
}

// TestGuidanceProseBudgetRefusesASymlinkedSkillDirectory keeps the refusal at the directory
// level too. `.claude/skills` is a tree of symlinks to these very files, so a link planted
// beside the canonical skills is the shape that would pull an adapter surface into the
// budget universe under a canonical path.
func TestGuidanceProseBudgetRefusesASymlinkedSkillDirectory(t *testing.T) {
	root := writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), healthyProseBudgetFiles())
	link := filepath.Join(root, ".agents", "skills", "bench-craft-adapter")
	if err := os.Symlink(filepath.Join(root, ".agents", "skills", "bench-craft-line"), link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
	diags := checkGuidanceProseBudgets(root)
	if !containsDiagnostic(diags, "prose-budget subject refused: .agents/skills/bench-craft-adapter is a symbolic link, not a regular directory") {
		t.Fatalf("a symlinked skill directory was followed:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetReportsEverySubjectInOneRun pins the quantifier. A checker that
// returns at its first hit sends a reviewer back for another run per violation, and the
// second and third paths are exactly what a one-violation test cannot see.
func TestGuidanceProseBudgetReportsEverySubjectInOneRun(t *testing.T) {
	files := map[string]string{
		".bench/BENCH.md": proseBudgetLines(151, true),
		".agents/skills/bench-craft-tickets/SKILL.md": proseBudgetLines(101, true),
		".agents/skills/bench-craft-line/SKILL.md":    proseBudgetLines(121, true),
	}
	diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), files))
	for _, want := range []string{
		"prose-budget exceeded: .bench/BENCH.md is 151 lines, over its 150-line budget",
		"prose-budget exceeded: .agents/skills/bench-craft-tickets/SKILL.md is 101 lines, over its 100-line budget",
		"prose-budget exceeded: .agents/skills/bench-craft-line/SKILL.md is 121 lines, over its 120-line budget",
	} {
		if !containsDiagnostic(diags, want) {
			t.Fatalf("want %q in one run, got:\n%s", want, strings.Join(diags, "\n"))
		}
	}
}

// TestGuidanceProseBudgetBoundaryIsNewlineIndifferent walks the edge from both sides under
// both file endings. A count that reads the trailing newline as a fourth line would move
// the whole boundary by one, so exactly-at-budget and one-over are asserted together.
func TestGuidanceProseBudgetBoundaryIsNewlineIndifferent(t *testing.T) {
	for _, trailing := range []bool{true, false} {
		name := "with trailing newline"
		if !trailing {
			name = "without trailing newline"
		}
		t.Run(name, func(t *testing.T) {
			at := healthyProseBudgetFiles()
			at[".agents/skills/bench-craft-tickets/SKILL.md"] = proseBudgetLines(100, trailing)
			if diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), at)); len(diags) != 0 {
				t.Fatalf("a subject at exactly its limit got diagnostics:\n%s", strings.Join(diags, "\n"))
			}
			over := healthyProseBudgetFiles()
			over[".agents/skills/bench-craft-tickets/SKILL.md"] = proseBudgetLines(101, trailing)
			diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), over))
			if !containsDiagnostic(diags, "prose-budget exceeded: .agents/skills/bench-craft-tickets/SKILL.md is 101 lines, over its 100-line budget") {
				t.Fatalf("one line over the limit was accepted:\n%s", strings.Join(diags, "\n"))
			}
		})
	}
}

// TestGuidanceProseBudgetsHoldOnTheLiveTree is the check's live-tree assertion: the kit's
// own guidance sits inside the budgets its profile publishes.
func TestGuidanceProseBudgetsHoldOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	if diags := checkGuidanceProseBudgets(h.KitRoot); len(diags) != 0 {
		t.Fatalf("the kit's guidance is over its published prose budgets:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetCanaryFixtureBites is the fixture owner: the family's mutation
// pushes a classified guidance file over its limit, the bound check raises the fixture's
// own diagnostic, and restoring the fixture takes that diagnostic away — so the red belongs
// to the mutation rather than to ambient state in the materialized tree.
func TestGuidanceProseBudgetCanaryFixtureBites(t *testing.T) {
	const fixtureName = "over-budget-skill"
	h := NewHarness(t)
	owner, bound := conformanceChecks["guidance-prose-budgets"]
	if !bound {
		t.Fatal("guidance-prose-budgets conformance owner is not bound")
	}
	fixture := h.KitPath("tests", "canary", "guidance-prose-budgets", fixtureName)
	want := strings.TrimSpace(readIfExists(filepath.Join(fixture, "EXPECT")))
	if want == "" {
		t.Fatalf("fixture %s has no EXPECT diagnostic", fixtureName)
	}
	root := t.TempDir()
	if err := canary.MaterializeMutationFixture(h.KitRoot, fixture, root); err != nil {
		t.Fatalf("materialize %s: %v", fixtureName, err)
	}
	if diags := owner.run(root, h.KitRoot, registry.Dev); !containsDiagnostic(diags, want) {
		t.Fatalf("the mutated fixture did not raise %q:\n%s", want, strings.Join(diags, "\n"))
	}
	if err := canary.RestoreMutationFixture(h.KitRoot, fixture, root); err != nil {
		t.Fatalf("restore %s: %v", fixtureName, err)
	}
	if diags := owner.run(root, h.KitRoot, registry.Dev); containsDiagnostic(diags, want) {
		t.Fatalf("the restored fixture still raises %q:\n%s", want, strings.Join(diags, "\n"))
	}
}
