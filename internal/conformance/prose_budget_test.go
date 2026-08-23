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
	// proseBudgetProfile holds the reviewer-owned budget table. Only this file changes
	// when the reviewer raises a budget.
	proseBudgetProfile = "projects/benchkit.md"
	// proseBudgetSection anchors the parse. This stops the profile's other tables from
	// becoming budget policy.
	proseBudgetSection = "Guidance prose budgets"
	// proseBudgetSkillsDir is the guidance tree whose SKILL.md files join the universe
	// structurally. The table's own glob rows cannot enumerate this tree. If the table
	// loses its default row, it enumerates nothing. A skill added next month then passes
	// with no budget applied. The check does not report it as unclassified.
	proseBudgetSkillsDir = ".agents/skills"
	// proseBudgetSkillFile is the one file per skill directory that the budget grades. A
	// skill's references and assets stay outside the reviewed universe.
	proseBudgetSkillFile = "SKILL.md"
)

// proseBudgetPolicy is the profile's table. It holds exact subject rows and glob rows that
// classify whatever the enumeration finds.
type proseBudgetPolicy struct {
	exact    map[string]int
	patterns []proseBudgetPattern
}

type proseBudgetPattern struct {
	glob  string
	limit int
}

// limitFor resolves one subject's budget. The check consults exact rows before glob rows, so
// the reviewer's specific number always beats the default that also matches the same subject.
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
		// This is one classification, not a second reader for the file mode. The shared
		// producer classifier owns shape and bounded bytes; this check owns only the cost
		// each state adds to the budget report.
		subject := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(rel)))
		switch {
		case subject.State == bounds.StateAbsent:
			diags = append(diags, "prose-budget subject missing: "+rel+" is named by the "+proseBudgetProfile+" budget table but absent from the tree")
			continue
		case subject.State == bounds.StateWrongType:
			diags = append(diags, "prose-budget subject refused: "+rel+" is not a regular file ("+subject.Reason+")")
			continue
		}
		limit, classified := policy.limitFor(rel)
		if !classified {
			diags = append(diags, "prose-budget subject unclassified: "+rel+" matches no row in the "+proseBudgetProfile+" budget table")
			continue
		}
		if subject.State.Failed() {
			diags = append(diags, "prose-budget subject unreadable: "+rel+" could not be read ("+subject.Reason+")")
			continue
		}
		if count := proseBudgetLineCount(subject.Data); count > limit {
			diags = append(diags, fmt.Sprintf("prose-budget exceeded: %s is %d lines, over its %d-line budget", rel, count, limit))
		}
	}
	return diags
}

// parseProseBudgetPolicy reads the reviewer's table out of the profile. Any policy fault
// returns diagnostics and no policy. A table nobody can read cannot say which subject is
// over budget. If the check grades a partial parse, it reports a clean tree for a broken
// policy.
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
			// The table ends at the first line that is not a row. Prose after that line
			// belongs to the section, not to the policy.
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

// proseBudgetSubjects lists the canonical universe in stable order: every exact subject the
// table names, plus every skill's SKILL.md. It also returns one diagnostic per skill directory
// that is a symbolic link. The check refuses that entry instead of descending into it, because
// `.claude/skills` is a tree of links to these same files. Following a link would pull an
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
	entries, rootDiag := proseBudgetSkillEntries(root)
	if rootDiag != "" {
		diags = append(diags, rootDiag)
	}
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
		// A skill directory with no SKILL.md is a fact for a different check. Only a
		// subject the table names must exist.
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(skill))); err == nil {
			add(skill)
		}
	}
	sort.Strings(subjects)
	return subjects, diags
}

// proseBudgetSkillEntries classifies the skills root before anything reads through it. The
// check classifies the root for the same reason it classifies each child. A linked
// `.agents/skills` would enumerate whatever it points at under the canonical path. An absent
// root yields no entries and no diagnostic, because only a subject the table names must
// exist.
func proseBudgetSkillEntries(root string) ([]os.DirEntry, string) {
	dir := filepath.Join(root, filepath.FromSlash(proseBudgetSkillsDir))
	info, err := os.Lstat(dir)
	switch {
	case err != nil:
		return nil, ""
	case info.Mode()&os.ModeSymlink != 0:
		return nil, "prose-budget subject refused: " + proseBudgetSkillsDir + " is a symbolic link, not a regular directory"
	case !info.IsDir():
		return nil, "prose-budget subject refused: " + proseBudgetSkillsDir + " is not a directory"
	}
	// os.ReadDir sorts by filename, so the enumeration reports in one order run to run.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "prose-budget subject unreadable: " + proseBudgetSkillsDir + ": " + err.Error()
	}
	return entries, ""
}

// proseBudgetLineCount counts lines the way a reader does, so a file ending with a newline
// and the same file without one carry the same count.
func proseBudgetLineCount(body []byte) int {
	if len(body) == 0 {
		return 0
	}
	return bytes.Count(bytes.TrimSuffix(body, []byte("\n")), []byte("\n")) + 1
}

// proseBudgetTable renders the profile subsection that the checker parses. It takes the rows
// verbatim, so a test case can corrupt the header or a cell without a second table author.
func proseBudgetTable(header string, rows ...string) string {
	body := "# benchkit\n\n## Gate\n\n### " + proseBudgetSection + "\n\n" + header + "\n|---|---|\n"
	for _, row := range rows {
		body += row + "\n"
	}
	return body + "\n## Notes for cold sessions\n\ntail\n"
}

// proseBudgetHeader is the header row the parser anchors the table on.
const proseBudgetHeader = "| subject | limit |"

// proseBudgetRows is a healthy policy. It holds one exact row outside the skills tree, one
// exact row inside it, and the all-skills default that the exact row must beat.
var proseBudgetRows = []string{
	"| `.bench/BENCH.md` | 150 |",
	"| `.agents/skills/bench-craft-tickets/SKILL.md` | 100 |",
	"| `.agents/skills/*/SKILL.md` | 120 |",
}

// proseBudgetLines renders count logical lines of filler prose, ending with a newline
// unless trailing is false.
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

// healthyProseBudgetFiles renders every subject that proseBudgetRows names, each one line
// under its limit. This lets a test case push exactly the subject it is about over the edge.
func healthyProseBudgetFiles() map[string]string {
	return map[string]string{
		".bench/BENCH.md": proseBudgetLines(149, true),
		".agents/skills/bench-craft-tickets/SKILL.md": proseBudgetLines(99, true),
		".agents/skills/bench-craft-line/SKILL.md":    proseBudgetLines(119, true),
	}
}

// TestGuidanceProseBudgetsComeFromTheProfileTable is the single-source row. The limit the
// checker enforces moves when the reviewer moves the cell; a hard-coded constant cannot do
// this. The exact craft-tickets row must also beat the all-skills default it matches, so the
// test checks both in the same parse.
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

	// The default row allows 120, but the exact row grades this subject.
	over := healthyProseBudgetFiles()
	over[".agents/skills/bench-craft-tickets/SKILL.md"] = proseBudgetLines(110, true)
	diags = checkGuidanceProseBudgets(writeProseBudgetRoot(t, healthy, over))
	if !containsDiagnostic(diags, "over its 100-line budget") {
		t.Fatalf("the all-skills default overrode the exact craft-tickets row:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetPolicyFaultsFailClosed covers the four ways the reviewer's table
// can stop being readable. Each fault stops the run instead of grading a partial parse. A
// policy nobody can read cannot say which subject is over budget. A clean verdict from a
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
// tree, not to the table's own rows. The default row grades a skill added next month with
// no editor change. With no default row, the check reports the skill as unclassified. A
// table-driven enumeration cannot give this answer; it never looks at the new directory.
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
// fails by a deadline, not by a wrong answer from an implementation that opens the file
// first. The symlink case points at a subject far over budget. A checker that follows it
// reports the target's line count under the link's path.
func TestGuidanceProseBudgetRefusesNonRegularSubjects(t *testing.T) {
	for _, tc := range []struct {
		kind, want string
		plant      func(t *testing.T, path string)
	}{
		{
			kind: "symlink",
			want: "prose-budget subject refused: .agents/skills/bench-craft-linked/SKILL.md is not a regular file",
			plant: func(t *testing.T, path string) {
				target := filepath.Join(filepath.Dir(path), "..", "bench-craft-line", "SKILL.md")
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
		},
		{
			kind: "fifo",
			want: "prose-budget subject refused: .agents/skills/bench-craft-linked/SKILL.md is not a regular file",
			plant: func(t *testing.T, path string) {
				if err := syscall.Mkfifo(path, 0o644); err != nil {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
				}
			},
		},
		{
			kind: "socket",
			want: "prose-budget subject refused: .agents/skills/bench-craft-linked/SKILL.md is not a regular file",
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
// level too. `.claude/skills` is a tree of symlinks to these same files. A link planted
// beside the canonical skills would pull an adapter surface into the budget universe under
// a canonical path.
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

// TestGuidanceProseBudgetRefusesASymlinkedSkillsRoot keeps the refusal at the root of the
// guidance tree, where a single link would redirect the whole enumeration. The linked tree
// holds a skill far over budget. A checker that reads through the link reports that subject
// under a canonical path, instead of refusing the root.
func TestGuidanceProseBudgetRefusesASymlinkedSkillsRoot(t *testing.T) {
	root := writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows[0], proseBudgetRows[2]), map[string]string{
		".bench/BENCH.md":                          proseBudgetLines(149, true),
		"payload/skills/bench-craft-line/SKILL.md": proseBudgetLines(500, true),
	})
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "payload", "skills"), filepath.Join(root, ".agents", "skills")); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
	diags := checkGuidanceProseBudgets(root)
	if !containsDiagnostic(diags, "prose-budget subject refused: .agents/skills is a symbolic link, not a regular directory") {
		t.Fatalf("a symlinked skills root was not refused:\n%s", strings.Join(diags, "\n"))
	}
	if containsDiagnostic(diags, "prose-budget exceeded") {
		t.Fatalf("the check enumerated through the symlinked root:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetReportsAnAbsentTableSubject pins the missing-subject diagnostic. A
// row the reviewer keeps for a file nobody ships marks a stale table, not a clean tree.
func TestGuidanceProseBudgetReportsAnAbsentTableSubject(t *testing.T) {
	files := healthyProseBudgetFiles()
	delete(files, ".agents/skills/bench-craft-tickets/SKILL.md")
	diags := checkGuidanceProseBudgets(writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), files))
	if !containsDiagnostic(diags, "prose-budget subject missing: .agents/skills/bench-craft-tickets/SKILL.md is named by the projects/benchkit.md budget table but absent from the tree") {
		t.Fatalf("a table row naming an absent subject was accepted:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetReportsAnUnreadableSubject pins the unreadable diagnostic. When the
// check cannot get the bytes of a regular, classified subject, it reports this instead of
// counting the subject as within budget.
func TestGuidanceProseBudgetReportsAnUnreadableSubject(t *testing.T) {
	root := writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), healthyProseBudgetFiles())
	subject := filepath.Join(root, ".agents", "skills", "bench-craft-tickets", "SKILL.md")
	if err := os.Chmod(subject, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(subject, 0o644) })
	if _, err := os.ReadFile(subject); err == nil {
		capability.Capability(t, capability.Privilege, "the test process reads mode 0000 files, so an unreadable subject cannot be planted")
	}
	diags := checkGuidanceProseBudgets(root)
	if !containsDiagnostic(diags, "prose-budget subject unreadable: .agents/skills/bench-craft-tickets/SKILL.md could not be read") {
		t.Fatalf("an unreadable subject was accepted:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetReportsAnUnreadableSkillsRoot pins the fail-closed diagnostic for a
// skills root that Lstat classifies as a real directory, but os.ReadDir cannot enumerate. A
// permission fault must not silently drop every wildcard subject under this root.
func TestGuidanceProseBudgetReportsAnUnreadableSkillsRoot(t *testing.T) {
	root := writeProseBudgetRoot(t, proseBudgetTable(proseBudgetHeader, proseBudgetRows...), healthyProseBudgetFiles())
	dir := filepath.Join(root, ".agents", "skills")
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	if _, err := os.ReadDir(dir); err == nil {
		capability.Capability(t, capability.Privilege, "the test process reads mode 0000 directories, so an unreadable skills root cannot be planted")
	}
	diags := checkGuidanceProseBudgets(root)
	if !containsDiagnostic(diags, "prose-budget subject unreadable: .agents/skills: ") {
		t.Fatalf("an unreadable skills root was accepted:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetReportsEverySubjectInOneRun pins the quantifier. A checker that
// stops at its first hit sends the reviewer back for another run per violation. A test with
// only one violation cannot see the second and third paths.
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

// TestGuidanceProseBudgetBoundaryIsNewlineIndifferent walks the edge from both sides, under
// both file endings. A count that reads the trailing newline as an extra line moves the
// whole boundary by one. The test asserts exactly-at-budget and one-over together.
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

// TestGuidanceProseBudgetsHoldOnTheLiveTree is the check's live-tree assertion. The kit's
// own guidance sits inside the budgets its profile publishes.
func TestGuidanceProseBudgetsHoldOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	if diags := checkGuidanceProseBudgets(h.KitRoot); len(diags) != 0 {
		t.Fatalf("the kit's guidance is over its published prose budgets:\n%s", strings.Join(diags, "\n"))
	}
}

// TestGuidanceProseBudgetCanaryFixtureBites is the fixture owner. The family's mutation
// pushes a classified guidance file over its limit, and the bound check raises the
// fixture's own diagnostic. Restoring the fixture removes that diagnostic. The red belongs
// to the mutation, not to ambient state in the materialized tree.
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
