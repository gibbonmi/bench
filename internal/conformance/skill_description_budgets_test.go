package conformance

import (
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
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/skillsindex"
)

const (
	// skillDescriptionProfile holds the reviewer-owned budget table. Only this file
	// changes when the reviewer raises or lowers a description budget.
	skillDescriptionProfile = "projects/benchkit.md"
	// skillDescriptionSection anchors the parse, so the profile's other subject-and-limit
	// table cannot become description policy.
	skillDescriptionSection = "Skill description budgets"
	// skillDescriptionKey is the frontmatter key the harness listing renders.
	skillDescriptionKey = "description"
	// skillDescriptionSkillsDir and skillDescriptionCommandDir are the two trees whose
	// files enter the listing a session pays for on every turn.
	skillDescriptionSkillsDir  = ".agents/skills"
	skillDescriptionSkillFile  = "SKILL.md"
	skillDescriptionCommandDir = ".agents/commands"
)

// checkSkillDescriptionBudgets grades every shipped skill and command description against
// the profile's own numbers. The budget resolver is the prose budget's, because one subject
// pattern resolves to one limit the same way for both tables.
func checkSkillDescriptionBudgets(root string) []string {
	policy, diags := parseSkillDescriptionPolicy(readIfExists(filepath.Join(root, filepath.FromSlash(skillDescriptionProfile))))
	if len(diags) != 0 {
		return diags
	}
	subjects, diags := skillDescriptionSubjects(root)
	for _, rel := range subjects {
		// Classification precedes every read. A link is refused rather than followed,
		// and a FIFO cannot block the gate in open(2).
		subject := bounds.ClassifyNoFollow(filepath.Join(root, filepath.FromSlash(rel)))
		if subject.State != bounds.StateParsed {
			diags = append(diags, "skill-description subject refused: "+rel+" is not a readable regular file ("+subject.Reason+")")
			continue
		}
		limit, classified := policy.limitFor(rel)
		if !classified {
			diags = append(diags, "skill-description subject unclassified: "+rel+" matches no row in the "+skillDescriptionProfile+" description budget table")
			continue
		}
		// The frontmatter reader is the skills index's own. It returns the first value
		// line, which is exactly what a harness listing renders.
		value := skillsindex.FrontmatterField(filepath.Join(root, filepath.FromSlash(rel)), skillDescriptionKey)
		if strings.TrimSpace(value) == "" {
			diags = append(diags, "skill-description missing: "+rel+" declares no "+skillDescriptionKey+" value in its frontmatter")
			continue
		}
		if skillDescriptionFolded(string(subject.Data)) {
			diags = append(diags, "skill-description folded: "+rel+" continues its description on an indented second line, so the listing reads the first line alone")
		}
		if count := skillDescriptionLength(value); count > limit {
			diags = append(diags, fmt.Sprintf("skill-description exceeded: %s is %d characters, over its %d-character budget", rel, count, limit))
		}
	}
	return diags
}

// skillDescriptionLength counts the description the way a listing renders it: every run of
// white space becomes one space, and the count is runes rather than bytes. A byte count
// would red a correct description that carries an em dash.
func skillDescriptionLength(value string) int {
	value = strings.Join(strings.Fields(value), " ")
	return utf8.RuneCountInString(value)
}

// skillDescriptionFolded reports whether the description value continues on an indented
// line. The frontmatter reader returns the first value line, so a folded description would
// otherwise pass on a count the listing never renders. This reads the line after the key
// rather than parsing the block a second time.
func skillDescriptionFolded(body string) bool {
	lines := strings.Split(body, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return false
	}
	prefix := skillDescriptionKey + ":"
	for i, line := range lines[1:] {
		if line == "---" {
			return false
		}
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		if i+2 >= len(lines) {
			return false
		}
		next := lines[i+2]
		return next != "---" && strings.TrimSpace(next) != "" && (strings.HasPrefix(next, " ") || strings.HasPrefix(next, "\t"))
	}
	return false
}

// parseSkillDescriptionPolicy reads the reviewer's table out of the profile with the
// section and row helpers the prose budget check uses. Any policy fault returns
// diagnostics and no policy: a table nobody can read cannot say which description is over
// budget, and a clean verdict from a broken table is green by omission.
func parseSkillDescriptionPolicy(profile string) (proseBudgetPolicy, []string) {
	section, anchored := profileSection(profile, skillDescriptionSection)
	if !anchored {
		return proseBudgetPolicy{}, []string{"skill-description policy missing: " + skillDescriptionProfile + " has no '" + skillDescriptionSection + "' heading holding the budget table"}
	}
	policy := proseBudgetPolicy{exact: map[string]int{}}
	var diags []string
	header, rows := false, 0
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
		limit, err := strconv.Atoi(cell)
		if err != nil || limit <= 0 {
			diags = append(diags, "skill-description policy malformed row: "+skillDescriptionProfile+" gives subject '"+subject+"' the limit '"+cell+"', which is not a positive character count")
			continue
		}
		if _, err := path.Match(subject, ""); err != nil {
			diags = append(diags, "skill-description policy malformed row: "+skillDescriptionProfile+" gives subject '"+subject+"', which is not a valid path pattern")
			continue
		}
		rows++
		if strings.ContainsAny(subject, "*?[") {
			policy.patterns = append(policy.patterns, proseBudgetPattern{glob: subject, limit: limit})
			continue
		}
		policy.exact[subject] = limit
	}
	switch {
	case !header:
		return proseBudgetPolicy{}, []string{"skill-description policy malformed: " + skillDescriptionProfile + " renders no '| subject | limit |' table under its '" + skillDescriptionSection + "' heading"}
	case len(diags) != 0:
		return proseBudgetPolicy{}, diags
	case rows == 0:
		return proseBudgetPolicy{}, []string{"skill-description policy empty: " + skillDescriptionProfile + "'s description budget table names no subject"}
	}
	return policy, nil
}

// skillDescriptionSubjects lists the listing's universe in stable order: one SKILL.md per
// skill directory, and every Markdown file in the command tree. It also returns one
// diagnostic per directory entry that is a symbolic link, because `.claude/skills` is a
// tree of links to these same files and following one would pull an adapter surface into
// the universe under a canonical path.
func skillDescriptionSubjects(root string) (subjects, diags []string) {
	entries, diag := skillDescriptionEntries(root, skillDescriptionSkillsDir)
	if diag != "" {
		diags = append(diags, diag)
	}
	for _, entry := range entries {
		rel := path.Join(skillDescriptionSkillsDir, entry.Name())
		if entry.Type()&os.ModeSymlink != 0 {
			diags = append(diags, "skill-description subject refused: "+rel+" is a symbolic link, not a regular directory")
			continue
		}
		if !entry.IsDir() {
			continue
		}
		skill := path.Join(rel, skillDescriptionSkillFile)
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(skill))); err == nil {
			subjects = append(subjects, skill)
		}
	}
	commands, diag := skillDescriptionEntries(root, skillDescriptionCommandDir)
	if diag != "" {
		diags = append(diags, diag)
	}
	for _, entry := range commands {
		rel := path.Join(skillDescriptionCommandDir, entry.Name())
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			diags = append(diags, "skill-description subject refused: "+rel+" is a symbolic link, not a regular directory")
			continue
		}
		subjects = append(subjects, rel)
	}
	sort.Strings(subjects)
	return subjects, diags
}

// skillDescriptionEntries classifies a listing tree's root before anything reads through
// it. A linked root would enumerate whatever it points at under the canonical path. An
// absent root yields no entries and no diagnostic, because a consumer repo ships neither
// tree until it links.
func skillDescriptionEntries(root, rel string) ([]os.DirEntry, string) {
	dir := filepath.Join(root, filepath.FromSlash(rel))
	info, err := os.Lstat(dir)
	switch {
	case err != nil:
		return nil, ""
	case info.Mode()&os.ModeSymlink != 0:
		return nil, "skill-description subject refused: " + rel + " is a symbolic link, not a regular directory"
	case !info.IsDir():
		return nil, "skill-description subject refused: " + rel + " is not a directory"
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, "skill-description subject unreadable: " + rel + ": " + err.Error()
	}
	return entries, ""
}

// skillDescriptionTable renders the profile subsection the checker parses. It takes the
// rows verbatim, so a test case can corrupt the header or a cell without a second table
// author.
func skillDescriptionTable(header string, rows ...string) string {
	body := "# benchkit\n\n## Gate\n\n### " + skillDescriptionSection + "\n\n" + header + "\n|---|---|\n"
	for _, row := range rows {
		body += row + "\n"
	}
	return body + "\n## Notes for cold sessions\n\ntail\n"
}

// skillDescriptionHeader is the header row the parser anchors the table on, and
// skillDescriptionRows is a healthy policy: one glob row per listing tree.
const skillDescriptionHeader = "| subject | limit |"

var skillDescriptionRows = []string{
	"| `.agents/skills/*/SKILL.md` | 40 |",
	"| `.agents/commands/*.md` | 30 |",
}

// skillDescriptionFile renders a frontmatter block carrying one description value.
func skillDescriptionFile(description string) string {
	return "---\nname: fixture\ndescription: " + description + "\n---\n\nbody\n"
}

// healthySkillDescriptionFiles renders one subject per glob row, each inside its limit.
func healthySkillDescriptionFiles() map[string]string {
	return map[string]string{
		".agents/skills/bench-craft-fixture/SKILL.md": skillDescriptionFile("A fixture skill within its budget."),
		".agents/commands/bench-fixture.md":           skillDescriptionFile("A fixture command in budget."),
	}
}

// writeSkillDescriptionRoot plants a profile and the named listing files in a throwaway
// root.
func writeSkillDescriptionRoot(t *testing.T, profile string, files map[string]string) string {
	t.Helper()
	planted := map[string]string{skillDescriptionProfile: profile}
	for rel, content := range files {
		planted[rel] = content
	}
	return throwawayRoot{files: planted}.build(t)
}

// TestSkillDescriptionBudgetsComeFromTheProfileTable is the single-source row. The limit
// the checker enforces moves when the reviewer moves the cell, which a hard-coded constant
// cannot do. The lowered cell is named independently of the table the healthy case uses,
// so a checker that ignores the table cannot satisfy both halves.
func TestSkillDescriptionBudgetsComeFromTheProfileTable(t *testing.T) {
	healthy := skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...)
	if diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, healthy, healthySkillDescriptionFiles())); len(diags) != 0 {
		t.Fatalf("a tree inside every budget got diagnostics:\n%s", strings.Join(diags, "\n"))
	}
	lowered := strings.Replace(healthy, "SKILL.md` | 40 |", "SKILL.md` | 20 |", 1)
	diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, lowered, healthySkillDescriptionFiles()))
	if !containsDiagnostic(diags, "skill-description exceeded: .agents/skills/bench-craft-fixture/SKILL.md is 34 characters, over its 20-character budget") {
		t.Fatalf("lowering the skill cell to 20 did not move the enforced limit:\n%s", strings.Join(diags, "\n"))
	}
	over := healthySkillDescriptionFiles()
	over[".agents/commands/bench-fixture.md"] = skillDescriptionFile("A fixture command description that runs past the command budget.")
	diags = checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, healthy, over))
	if !containsDiagnostic(diags, "skill-description exceeded: .agents/commands/bench-fixture.md is 64 characters, over its 30-character budget") {
		t.Fatalf("the command row did not grade the command tree:\n%s", strings.Join(diags, "\n"))
	}
}

// TestSkillDescriptionBudgetCountsCollapsedRunes pins the two halves of the count. The
// multibyte case is one rune per character, so a byte count reds a description inside its
// budget. The padded case carries 34 characters of white space runs that collapse to 34
// rendered characters minus the padding, so a count taken before the collapse reports a
// subject the listing never renders that long.
func TestSkillDescriptionBudgetCountsCollapsedRunes(t *testing.T) {
	files := healthySkillDescriptionFiles()
	// Thirty-eight runes, and more bytes than that: each em dash costs three bytes.
	files[".agents/skills/bench-craft-fixture/SKILL.md"] = skillDescriptionFile("Rødgrød — æblegrød — 中文 — fixture line")
	if diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...), files)); len(diags) != 0 {
		t.Fatalf("a multibyte description inside its budget was counted as bytes:\n%s", strings.Join(diags, "\n"))
	}
	// Thirty-four rendered characters, padded to fifty-one raw ones.
	padded := healthySkillDescriptionFiles()
	padded[".agents/skills/bench-craft-fixture/SKILL.md"] = skillDescriptionFile("A     fixture     skill     within     its     budget.")
	if diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...), padded)); len(diags) != 0 {
		t.Fatalf("white space runs were counted instead of collapsed:\n%s", strings.Join(diags, "\n"))
	}
	lowered := strings.Replace(skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...), "SKILL.md` | 40 |", "SKILL.md` | 20 |", 1)
	diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, lowered, padded))
	if !containsDiagnostic(diags, "skill-description exceeded: .agents/skills/bench-craft-fixture/SKILL.md is 34 characters, over its 20-character budget") {
		t.Fatalf("the padded description did not report its collapsed count:\n%s", strings.Join(diags, "\n"))
	}
}

// TestSkillDescriptionBudgetPolicyFaultsFailClosed covers the ways the reviewer's table
// can stop being readable. Each fault stops the run instead of grading a partial parse.
func TestSkillDescriptionBudgetPolicyFaultsFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name, profile, want string
	}{
		{
			"missing section",
			"# benchkit\n\n## Gate\n\nNo description budget table lives here.\n",
			"skill-description policy missing: projects/benchkit.md has no 'Skill description budgets' heading",
		},
		{
			"malformed header",
			skillDescriptionTable("| path | maximum |", skillDescriptionRows...),
			"skill-description policy malformed: projects/benchkit.md renders no '| subject | limit |' table",
		},
		{
			"malformed row",
			skillDescriptionTable(skillDescriptionHeader, "| `.agents/skills/*/SKILL.md` | abc |"),
			"skill-description policy malformed row: projects/benchkit.md gives subject '.agents/skills/*/SKILL.md' the limit 'abc'",
		},
		{
			"empty table",
			skillDescriptionTable(skillDescriptionHeader),
			"skill-description policy empty: projects/benchkit.md's description budget table names no subject",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, tc.profile, healthySkillDescriptionFiles()))
			if !containsDiagnostic(diags, tc.want) {
				t.Fatalf("want %q, got:\n%s", tc.want, strings.Join(diags, "\n"))
			}
			if containsDiagnostic(diags, "skill-description exceeded") {
				t.Fatalf("a broken policy still graded a subject:\n%s", strings.Join(diags, "\n"))
			}
		})
	}
}

// TestSkillDescriptionBudgetReportsMissingAndFoldedValues pins the two shape diagnostics.
// A file with no description is not a lean one, and a folded value counts a line the
// listing never renders whole.
func TestSkillDescriptionBudgetReportsMissingAndFoldedValues(t *testing.T) {
	files := healthySkillDescriptionFiles()
	files[".agents/skills/bench-craft-fixture/SKILL.md"] = "---\nname: fixture\n---\n\nbody\n"
	files[".agents/commands/bench-fixture.md"] = "---\nname: fixture\ndescription: A fixture command\n  that folds onto a second line.\n---\n\nbody\n"
	diags := checkSkillDescriptionBudgets(writeSkillDescriptionRoot(t, skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...), files))
	for _, want := range []string{
		"skill-description missing: .agents/skills/bench-craft-fixture/SKILL.md declares no description value in its frontmatter",
		"skill-description folded: .agents/commands/bench-fixture.md continues its description on an indented second line",
	} {
		if !containsDiagnostic(diags, want) {
			t.Fatalf("want %q in one run, got:\n%s", want, strings.Join(diags, "\n"))
		}
	}
}

// TestSkillDescriptionBudgetRefusesNonRegularSubjects covers the entry kinds a listing
// tree can hold besides a file. Classification precedes every read, so the FIFO case fails
// by a deadline rather than by a wrong answer from an implementation that opens first.
func TestSkillDescriptionBudgetRefusesNonRegularSubjects(t *testing.T) {
	for _, tc := range []struct {
		kind  string
		plant func(t *testing.T, path string)
	}{
		{
			kind: "symlink",
			plant: func(t *testing.T, path string) {
				target := filepath.Join(filepath.Dir(path), "..", "bench-craft-fixture", skillDescriptionSkillFile)
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
		},
		{
			kind: "fifo",
			plant: func(t *testing.T, path string) {
				if err := syscall.Mkfifo(path, 0o644); err != nil {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
				}
			},
		},
		{
			kind: "socket",
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
			root := writeSkillDescriptionRoot(t, skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...), healthySkillDescriptionFiles())
			dir := filepath.Join(root, ".agents", "skills", "bench-craft-planted")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			tc.plant(t, filepath.Join(dir, skillDescriptionSkillFile))
			done := make(chan []string, 1)
			go func() { done <- checkSkillDescriptionBudgets(root) }()
			select {
			case diags := <-done:
				if !containsDiagnostic(diags, "skill-description subject refused: .agents/skills/bench-craft-planted/SKILL.md is not a readable regular file") {
					t.Fatalf("a %s subject was not refused:\n%s", tc.kind, strings.Join(diags, "\n"))
				}
			case <-time.After(bounds.TestDeadline(0)):
				t.Fatalf("the check blocked on a %s, so it opened the subject before classifying it", tc.kind)
			}
		})
	}
}

// TestSkillDescriptionBudgetRefusesASymlinkedSkillDirectory keeps the refusal at the
// directory level. A link planted beside the canonical skills would pull an adapter
// surface into the graded universe under a canonical path.
func TestSkillDescriptionBudgetRefusesASymlinkedSkillDirectory(t *testing.T) {
	root := writeSkillDescriptionRoot(t, skillDescriptionTable(skillDescriptionHeader, skillDescriptionRows...), healthySkillDescriptionFiles())
	link := filepath.Join(root, ".agents", "skills", "bench-craft-adapter")
	if err := os.Symlink(filepath.Join(root, ".agents", "skills", "bench-craft-fixture"), link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
	diags := checkSkillDescriptionBudgets(root)
	if !containsDiagnostic(diags, "skill-description subject refused: .agents/skills/bench-craft-adapter is a symbolic link, not a regular directory") {
		t.Fatalf("a symlinked skill directory was followed:\n%s", strings.Join(diags, "\n"))
	}
}

// TestSkillDescriptionBudgetsHoldOnTheLiveTree is the check's live-tree assertion. Every
// shipped description sits inside the budget the kit's profile publishes.
func TestSkillDescriptionBudgetsHoldOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	if diags := checkSkillDescriptionBudgets(h.KitRoot); len(diags) != 0 {
		t.Fatalf("the kit's descriptions are over their published budgets:\n%s", strings.Join(diags, "\n"))
	}
}
