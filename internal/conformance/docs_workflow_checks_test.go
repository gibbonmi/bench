package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/coverage"
	"github.com/gibbonmi/bench/internal/roadmap"
)

func checkDocsCurrencyAndWorkflow(root, kitRoot string) []string {
	var diags []string
	diags = append(diags, checkStaleCommandReferences(root)...)
	diags = append(diags, checkColdPickupCLILists(root)...)
	diags = append(diags, checkAXIProfileAnchors(root)...)
	diags = append(diags, checkBenchReferenceTokenDiet(root)...)
	diags = append(diags, checkShippedDogfoodReferents(root)...)
	diags = append(diags, checkSignalVocabulary(root)...)
	diags = append(diags, checkCommandFirstAnchors(root)...)
	diags = append(diags, checkWorkflowAnchors(root)...)
	diags = append(diags, checkOccurrenceLedgerAndMaintenance(root)...)
	diags = append(diags, checkPrePushREADMEClaim(root)...)
	diags = append(diags, checkSkillsIndexGenerateVerify(root, kitRoot)...)
	diags = append(diags, checkCoverageMaps(root)...)
	diags = append(diags, checkRemovedVerbSweep(root)...)
	return diags
}

// checkRemovedVerbSweep keeps the deleted provisional lifecycle out of kit
// prose: the literal command tokens `spec build` and `worktree recovery` must
// not reappear in any surface an agent reads for guidance. The stale-command
// and cold-pickup sweeps match only slash-command and single-word `bench <cmd>`
// tokens, so a leftover `bench spec build promote` passes both. Three surfaces
// are exempt because their job is history, not guidance: CHANGELOG.md
// (append-only — scrubbing it would falsify the record), capture/ (journal),
// and specs/remove-spec-build-lifecycle/ (the removal's own record). Retired
// spec residue without a `Status: staged` spec.md is history too; only staged
// specs are guidance a build will act on.
func checkRemovedVerbSweep(root string) []string {
	var files []string
	for _, rel := range []string{"README.md", "ROADMAP.md"} {
		if path := filepath.Join(root, rel); exists(path) {
			files = append(files, path)
		}
	}
	for _, rel := range []string{".bench", ".agents", "projects"} {
		for _, path := range walkConformanceDocs(filepath.Join(root, rel)) {
			if strings.HasSuffix(path, ".md") {
				files = append(files, path)
			}
		}
	}
	stagedRE := regexp.MustCompile(`(?m)^Status: staged$`)
	specFiles, _ := filepath.Glob(filepath.Join(root, "specs", "*", "spec.md"))
	sort.Strings(specFiles)
	for _, specPath := range specFiles {
		dir := filepath.Dir(specPath)
		if filepath.Base(dir) == "remove-spec-build-lifecycle" || !stagedRE.MatchString(readIfExists(specPath)) {
			continue
		}
		for _, path := range walkConformanceDocs(dir) {
			if strings.HasSuffix(path, ".md") {
				files = append(files, path)
			}
		}
	}
	files = uniqueSorted(files)

	var diags []string
	for _, file := range files {
		rel := slashRel(root, file)
		for i, line := range strings.Split(readIfExists(file), "\n") {
			for _, token := range []string{"spec build", "worktree recovery"} {
				if strings.Contains(line, token) {
					diags = append(diags, fmt.Sprintf("%s:%d carries removed command token %q; the spec-build lifecycle and the recovery verb are deleted", rel, i+1, token))
				}
			}
		}
	}
	return diags
}

func TestRemovedVerbSweepBites(t *testing.T) {
	write := func(t *testing.T, root, rel, content string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		name string
		rel  string
		text string
		want bool
	}{
		{"operating guide token", ".bench/BENCH.md", "run `bench spec build promote` to land it\n", true},
		{"skill recovery token", ".agents/skills/bench-craft-delegate/SKILL.md", "retire it with bench worktree recovery\n", true},
		{"roadmap token", "ROADMAP.md", "resolve the run via `bench spec build abandon`\n", true},
		{"staged spec ticket token", "specs/some-feature/tickets/one.md", "submit through `bench spec build checkpoint`\n", true},
		{"removal record exempt", "specs/remove-spec-build-lifecycle/spec.md", "Status: staged\ndelete the `bench spec build` grammar\n", false},
		{"changelog exempt", "CHANGELOG.md", "- Removed `bench spec build` and `bench worktree recovery`.\n", false},
		{"capture exempt", "capture/learnings.md", "ran `bench spec build abandon` and it refused\n", false},
		{"hyphenated prose is not the token", ".bench/BENCH.md", "the spec-build lifecycle was removed\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			// A staged spec.md makes the folder's files sweep subjects; the
			// removal record's own folder must stay exempt even when staged.
			if strings.HasPrefix(tt.rel, "specs/") {
				dir := filepath.Dir(strings.TrimPrefix(tt.rel, "specs/"))
				for strings.Contains(dir, "/") {
					dir = filepath.Dir(dir)
				}
				write(t, root, "specs/"+dir+"/spec.md", "Status: staged\n")
			}
			write(t, root, tt.rel, tt.text)
			got := len(checkRemovedVerbSweep(root)) > 0
			if got != tt.want {
				t.Fatalf("sweep diagnostics present = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("unstaged spec residue exempt", func(t *testing.T) {
		root := t.TempDir()
		write(t, root, "specs/old-run/spec.md", "Status: implemented\n")
		write(t, root, "specs/old-run/tickets/one.md", "land through `bench spec build integrate`\n")
		write(t, root, "specs/residue/tickets/one.md", "bench worktree recovery cleans it\n")
		if diags := checkRemovedVerbSweep(root); len(diags) != 0 {
			t.Fatalf("unstaged spec residue was swept:\n%s", strings.Join(diags, "\n"))
		}
	})
}

func checkPrePushREADMEClaim(root string) []string {
	readme := strings.ReplaceAll(readIfExists(filepath.Join(root, "README.md")), "`", "")
	var diags []string
	if !strings.Contains(readme, "pre-push hook protects the branch it resolves") {
		diags = append(diags, "README.md does not say the pre-push hook protects the branch it resolves")
	}
	if strings.Contains(readme, "pre-push hook protects the default branch") {
		diags = append(diags, "README.md claims the pre-push hook protects the default branch without qualification")
	}
	return diags
}

func TestPrePushREADMEClaimBites(t *testing.T) {
	for _, test := range []struct {
		name   string
		readme string
	}{
		{"missing resolved branch", "the pre-push hook protects a branch"},
		{"unqualified default branch", "the pre-push hook protects the branch it resolves\nthe pre-push hook protects the default branch"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(test.readme), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkPrePushREADMEClaim(root); len(diags) == 0 {
				t.Fatal("pre-push README claim mutation passed")
			}
		})
	}
}

func checkOccurrenceLedgerMigration(root string) []string {
	data, err := os.ReadFile(filepath.Join(root, "ROADMAP.md"))
	if err != nil {
		return []string{"ROADMAP.md unavailable for occurrence-ledger migration check"}
	}
	doc, failures := roadmap.ParseDocument(data, nil, true)
	if len(failures) != 0 {
		return []string{"ROADMAP.md occurrence-ledger migration has malformed rows"}
	}
	want := map[string]int{"FT71": 1, "FT158": 3, "FT98": 3, "FT169": 1, "FT141": 1, "FT94": 1, "FT125": 1}
	got := map[string]int{}
	for _, row := range doc.Rows {
		if _, ok := want[row.ID]; ok {
			got[row.ID] = row.OccurrenceCount
		}
	}
	var diags []string
	for id, count := range want {
		if got[id] != count {
			diags = append(diags, "ROADMAP.md occurrence-ledger migration count for "+id+" is wrong")
		}
	}
	if strings.Contains(strings.ToLower(string(data)), "evidence supplied") {
		diags = append(diags, "ROADMAP.md retains a legacy evidence-supplied heading count")
	}
	return diags
}

func TestOccurrenceLedgerMigrationCheckBites(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte("**FT71 (HIGH, evidence supplied) — title.**\nOccurrences: baseline-01\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if diags := checkOccurrenceLedgerMigration(root); len(diags) == 0 {
		t.Fatal("legacy heading mutation passed occurrence-ledger migration check")
	}
}

func checkSpecAuthorizationContract(root string) []string {
	path := filepath.Join(root, ".agents", "commands", "bench-write-spec.md")
	if !exists(path) {
		return nil
	}
	entryContract, ok := markdownH2Sections(stripHTMLComments(readIfExists(path)), "Entry contract")
	if ok != 1 {
		return []string{".agents/commands/bench-write-spec.md dropped the exactly-three authorization boundary"}
	}
	declaration := strings.ToLower(collapseSpace(entryContract))
	if marker := strings.Index(declaration, "- **"); marker >= 0 {
		declaration = declaration[:marker]
	}
	var diags []string
	for _, authorization := range []struct{ anchor, diag string }{
		{"a ready compiled map", ".agents/commands/bench-write-spec.md dropped the ready compiled map authorization"},
		{"reviewer-confirmed current conversation", ".agents/commands/bench-write-spec.md dropped the reviewer-confirmed current conversation authorization"},
		{"named reviewed artifact", ".agents/commands/bench-write-spec.md dropped the named reviewed artifact authorization"},
	} {
		if !strings.Contains(declaration, authorization.anchor) {
			diags = append(diags, authorization.diag)
		}
	}
	if !strings.Contains(declaration, "accept exactly one of three decision sources") ||
		!strings.Contains(declaration, "no unnamed memory, unreviewed note, or fourth override authorizes a draft") {
		diags = append(diags, ".agents/commands/bench-write-spec.md dropped the exactly-three authorization boundary")
	}
	return diags
}

func checkCoverageMaps(root string) []string {
	specsDir := filepath.Join(root, "specs")
	if !exists(specsDir) {
		return nil
	}
	matches, _ := filepath.Glob(filepath.Join(specsDir, "*", "spec.md"))
	flat, _ := filepath.Glob(filepath.Join(specsDir, "*.md"))
	sort.Strings(matches)
	sort.Strings(flat)
	var diags []string
	for _, path := range flat {
		diags = append(diags, fmt.Sprintf("stray flat live spec: %s; move it to a folder containing spec.md", slashRel(root, path)))
	}
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

func checkStaleCommandReferences(root string) []string {
	commandsDir := filepath.Join(root, ".agents", "commands")
	validSlash := map[string]bool{"/model": true}
	commandFiles, _ := filepath.Glob(filepath.Join(commandsDir, "*.md"))
	for _, file := range commandFiles {
		validSlash["/"+strings.TrimSuffix(filepath.Base(file), ".md")] = true
	}
	validCodex := map[string]bool{}
	for slash := range validSlash {
		if strings.HasPrefix(slash, "/bench-") {
			validCodex["$"+strings.TrimPrefix(slash, "/")] = true
		}
	}

	var files []string
	for _, rel := range []string{
		"README.md",
		"AGENTS.md",
		".bench/BENCH.md",
		".bench/BENCH-reference.md",
		"capture/learnings.md",
		"CONTEXT.md",
		"HANDOFF.md",
		"CHANGELOG.md",
	} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if exists(path) {
			files = append(files, path)
		}
	}
	for _, rel := range []string{"specs", "decisions", ".agents"} {
		files = append(files, walkConformanceDocs(filepath.Join(root, filepath.FromSlash(rel)))...)
	}
	files = uniqueSorted(files)

	knownStale := map[string]bool{
		"/resynthesize":   true,
		"/spec":           true,
		"/grill":          true,
		"/start-ideation": true,
		"/setup":          true,
		"/build":          true,
		"/prep-shift":     true,
		"/fix-bug":        true,
		"/verify-gate":    true,
		"/map":            true,
		"/diagnose":       true,
		"/review":         true,
		"/verify":         true,
		"/shift":          true,
	}
	historicalMarker := regexp.MustCompile(`(?m)^<!-- command-currency: historical -->$`)
	slashRef := regexp.MustCompile(`(^|[\s([` + "`" + `"'])/([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])`)
	codexRef := regexp.MustCompile(`(^|[\s([` + "`" + `"'])\$([A-Za-z][A-Za-z0-9_-]*[A-Za-z0-9])`)

	var diags []string
	for _, file := range files {
		text := readIfExists(file)
		rel := slashRel(root, file)
		if historicalMarker.MatchString(text) {
			continue
		}
		switch rel {
		case "capture/learnings.md":
			text = strings.Split(text, "<!-- entries below -->")[0]
		case "CHANGELOG.md":
			if idx := strings.Index(text, "\n## "); idx >= 0 {
				text = text[:idx]
			}
		}
		for i, line := range strings.Split(text, "\n") {
			for _, match := range slashRef.FindAllStringSubmatch(line, -1) {
				token := "/" + match[2]
				if !validSlash[token] && (strings.HasPrefix(token, "/bench-") || knownStale[token]) {
					diags = append(diags, fmt.Sprintf("stale command reference %s in %s:%d", token, rel, i+1))
				}
			}
			for _, match := range codexRef.FindAllStringSubmatch(line, -1) {
				token := "$" + match[2]
				if strings.HasPrefix(token, "$bench-") && !validCodex[token] {
					diags = append(diags, fmt.Sprintf("stale Codex adapter reference %s in %s:%d", token, rel, i+1))
				}
			}
		}
	}
	return diags
}

func checkColdPickupCLILists(root string) []string {
	bench := readIfExists(filepath.Join(root, "bin", "bench.sh"))
	if bench == "" {
		return nil
	}
	cmdRE := regexp.MustCompile(`(?m)^  ([a-z][a-z-]*)\)\s`)
	var commands []string
	for _, match := range cmdRE.FindAllStringSubmatch(bench, -1) {
		commands = append(commands, match[1])
	}
	sort.Strings(commands)
	known := map[string]bool{}
	for _, command := range commands {
		known[command] = true
	}
	docRef := regexp.MustCompile("`bench ([a-z][a-z-]*)\\b")
	var diags []string
	// The documented surface is the always-loaded inventory plus the reference file's
	// plumbing enumeration: sessions read BENCH.md, hooks-and-adapters plumbing is
	// deliberately demoted to BENCH-reference.md, and a route documented in neither
	// is still a cold-pickup hole.
	guide := readIfExists(filepath.Join(root, ".bench", "BENCH.md"))
	if guide != "" {
		guide += "\n" + readIfExists(filepath.Join(root, ".bench", "BENCH-reference.md"))
		for _, command := range commands {
			if !strings.Contains(guide, "bench "+command) {
				diags = append(diags, fmt.Sprintf(".bench/BENCH.md or .bench/BENCH-reference.md does not list CLI command 'bench %s'", command))
			}
		}
	}
	// Reverse check: a `bench <cmd>` reference that names no route is a dead pointer.
	// Commands and skills alike route the reader to the CLI in prose, so the whole
	// .agents markdown tree is in scope, and its file set comes from the same walk
	// checkStaleCommandReferences uses, so both sweeps discover the same files. Two
	// asymmetries are deliberate: the walk also yields shell, JSON, and YAML, where a
	// regex written for prose invites a false positive, so only .md is scanned; and the
	// `<!-- command-currency: historical -->` marker exempts a file from the stale-command
	// sweep but not from this one, because a canary fixture can only assert red, so an
	// exemption that turns a diagnostic green cannot be proven to work.
	refFiles := []string{"HANDOFF.md", ".bench/BENCH.md", ".bench/BENCH-reference.md"}
	agentDocs := walkConformanceDocs(filepath.Join(root, ".agents"))
	sort.Strings(agentDocs)
	for _, abs := range agentDocs {
		if strings.HasSuffix(abs, ".md") {
			refFiles = append(refFiles, slashRel(root, abs))
		}
	}
	for _, rel := range refFiles {
		path := filepath.Join(root, filepath.FromSlash(rel))
		text := readIfExists(path)
		if text == "" {
			continue
		}
		for _, match := range docRef.FindAllStringSubmatch(text, -1) {
			if !known[match[1]] {
				diags = append(diags, fmt.Sprintf("%s documents unknown CLI command 'bench %s' (removed or renamed in bin/bench.sh?)", rel, match[1]))
			}
		}
	}
	return diags
}

func checkAXIProfileAnchors(root string) []string {
	text := readIfExists(filepath.Join(root, "projects", "benchkit.md"))
	var diags []string
	if text == "" {
		return nil
	}
	if !strings.Contains(text, "bench diff") {
		diags = append(diags, "projects/benchkit.md does not name bench diff on the AXI seam")
	}
	if !strings.Contains(text, "bench coverage") {
		diags = append(diags, "projects/benchkit.md does not name bench coverage on the AXI seam")
	}
	return diags
}

func checkBenchReferenceTokenDiet(root string) []string {
	benchPath := filepath.Join(root, ".bench", "BENCH.md")
	if !exists(benchPath) {
		return nil
	}
	refRel := ".bench/BENCH-reference.md"
	refPath := filepath.Join(root, filepath.FromSlash(refRel))
	if !exists(refPath) {
		return []string{refRel + " missing: the token-diet reference file the operating guide points to"}
	}
	var diags []string
	bench := readIfExists(benchPath)
	ref := readIfExists(refPath)
	if !strings.Contains(bench, "BENCH-reference.md") {
		diags = append(diags, ".bench/BENCH.md does not point to .bench/BENCH-reference.md (agents can't find the moved lookup sections)")
	}
	claude := readIfExists(filepath.Join(root, "CLAUDE.md"))
	for _, line := range strings.Split(claude, "\n") {
		if strings.TrimSpace(line) == "@.bench/BENCH-reference.md" {
			diags = append(diags, ".bench/BENCH-reference.md is @-imported by CLAUDE.md; it must stay on-demand (referenced by path, never imported, or the token diet regresses)")
			break
		}
	}
	benchHeads := h2Headings(bench)
	refHeads := h2Headings(ref)
	var dup []string
	for head := range benchHeads {
		if refHeads[head] {
			dup = append(dup, head)
		}
	}
	sort.Strings(dup)
	if len(dup) > 0 {
		diags = append(diags, "section heading present in both .bench/BENCH.md and .bench/BENCH-reference.md (single-source violation): "+strings.Join(dup, "|")+"|")
	}
	return diags
}

func checkShippedDogfoodReferents(root string) []string {
	// Platform files installed verbatim into consumer repos must stay
	// consumer-generic: a dogfood-only referent shipped here lands in every
	// linked repo. AGENTS.md is exempt — consumers keep their own.
	needles := []string{"projects/benchkit"}
	var files []string
	for _, rel := range []string{".bench/BENCH.md", ".bench/BENCH-reference.md"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if exists(path) {
			files = append(files, path)
		}
	}
	files = append(files, walkConformanceDocs(filepath.Join(root, ".agents"))...)
	files = uniqueSorted(files)

	var diags []string
	for _, file := range files {
		rel := slashRel(root, file)
		for i, line := range strings.Split(readIfExists(file), "\n") {
			for _, needle := range needles {
				if strings.Contains(line, needle) {
					diags = append(diags, fmt.Sprintf("shipped platform file %s:%d carries dogfood referent %q — use projects/<name>.md or move the fact into the profile", rel, i+1, needle))
				}
			}
		}
	}
	return diags
}

func checkCommandFirstAnchors(root string) []string {
	var diags []string
	readme := readIfExists(filepath.Join(root, "README.md"))
	if readme == "" {
		diags = append(diags, "README.md missing")
	} else {
		firstH2 := ""
		for _, line := range strings.Split(readme, "\n") {
			if strings.HasPrefix(line, "## ") {
				firstH2 = line
				break
			}
		}
		if firstH2 != "## Reviewer quick start" {
			if firstH2 == "" {
				firstH2 = "(none)"
			}
			diags = append(diags, fmt.Sprintf("README first H2 is '%s'; expected '## Reviewer quick start'", firstH2))
		}
	}

	commandFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "commands", "*.md"))
	sort.Strings(commandFiles)
	for _, file := range commandFiles {
		rel := slashRel(root, file)
		text := readIfExists(file)
		if !regexp.MustCompile(`(?m)^## Entry orientation$`).MatchString(text) {
			diags = append(diags, rel+" missing Entry orientation")
		}
		if !regexp.MustCompile(`(?m)^## Exit handoff$`).MatchString(text) {
			diags = append(diags, rel+" missing Exit handoff")
		}
	}
	return diags
}

func walkConformanceDocs(dir string) []string {
	var files []string
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == "SKILL.md" ||
			strings.HasSuffix(name, ".md") ||
			strings.HasSuffix(name, ".yaml") ||
			strings.HasSuffix(name, ".yml") ||
			strings.HasSuffix(name, ".json") ||
			strings.HasSuffix(name, ".sh") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func h2Headings(text string) map[string]bool {
	out := map[string]bool{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			out[line] = true
		}
	}
	return out
}
