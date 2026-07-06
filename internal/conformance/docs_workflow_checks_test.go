package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func checkDocsCurrencyAndWorkflow(root, kitRoot string) []string {
	var diags []string
	diags = append(diags, checkStaleCommandReferences(root)...)
	diags = append(diags, checkColdPickupCLILists(root)...)
	diags = append(diags, checkAXIProfileAnchors(root)...)
	diags = append(diags, checkBenchReferenceTokenDiet(root)...)
	diags = append(diags, checkShippedDogfoodReferents(root)...)
	diags = append(diags, checkCommandFirstAnchors(root)...)
	diags = append(diags, checkWorkflowAnchors(root)...)
	diags = append(diags, checkSkillsIndexGenerateVerify(root, kitRoot)...)
	diags = append(diags, checkCoverageMaps(root)...)
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
		".bench/learnings.md",
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
		case ".bench/learnings.md":
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
	for _, rel := range []string{".bench/BENCH.md"} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		text := readIfExists(path)
		if text == "" {
			continue
		}
		for _, command := range commands {
			if !strings.Contains(text, "bench "+command) {
				diags = append(diags, fmt.Sprintf("%s does not list CLI command 'bench %s'", rel, command))
			}
		}
	}
	// Reverse check: a `bench <cmd>` reference that names no route is a dead pointer.
	// The command phase files are in scope so the collapsed commit-discipline prose
	// (pointing at `bench commit` / `bench spec implemented`) is held to a routing token.
	refFiles := []string{"HANDOFF.md", ".bench/BENCH.md", ".bench/BENCH-reference.md"}
	cmdFiles, _ := filepath.Glob(filepath.Join(root, ".agents", "commands", "*.md"))
	sort.Strings(cmdFiles)
	for _, abs := range cmdFiles {
		refFiles = append(refFiles, slashRel(root, abs))
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
