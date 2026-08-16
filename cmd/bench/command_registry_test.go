package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	gitpkg "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

func TestCommandRunsVersionInProcess(t *testing.T) {
	var stdout, stderr, observation bytes.Buffer
	command := Command{Stdout: &stdout, Stderr: &stderr, Executable: "/selected/bench", Observe: &observation}
	if code := command.Run([]string{"version"}); code != 0 {
		t.Fatalf("version exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	want := versionLine(version, runtime.GOOS, runtime.GOARCH) + "\n"
	if stdout.String() != want {
		t.Fatalf("version stdout = %q, want %q", stdout.String(), want)
	}
	if want := "command-registry:version\n"; observation.String() != want {
		t.Fatalf("version implementation observation = %q, want %q", observation.String(), want)
	}
}

func TestCommandDispositionsAreComplete(t *testing.T) {
	want := map[processAttachment][]string{
		attachmentDirect: {"check-agent-line", "commit", "gate-go", "guard-git", "resume-clean", "session-inspect", "shift", "spec", "version", "worktree"},
		attachmentSystem: {"canary", "doctor", "freshness-check", "freshness-publish", "gate", "gate-phases", "gate-pin", "gate-run", "init", "link", "setup", "stop-verdict", "unlink", "upgrade", "worktree-hook"},
		attachmentShip:   {"prep-release", "release", "release-preflight"},
	}
	got := map[processAttachment][]string{}
	seen := map[string]bool{}
	for _, disposition := range commandDispositions() {
		if seen[disposition.Name] {
			t.Fatalf("command disposition repeats %q", disposition.Name)
		}
		seen[disposition.Name] = true
		got[disposition.Attachment] = append(got[disposition.Attachment], disposition.Name)
	}
	for attachment := range got {
		sort.Strings(got[attachment])
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command dispositions = %#v, want %#v", got, want)
	}
}

func TestCommandRegistryAXIDispositionsAreComplete(t *testing.T) {
	_ = axiRegistryMemberNames(t)
}

func axiRegistryMemberNames(t *testing.T) []string {
	t.Helper()
	var members []string
	seen := map[string]bool{}
	for _, definition := range commandRegistry {
		if seen[definition.Name] {
			t.Fatalf("commandRegistry repeats %q", definition.Name)
		}
		seen[definition.Name] = true
		disposition := definition.AXI
		switch {
		case disposition.root:
			if len(disposition.children) != 0 || disposition.exemption != "" {
				t.Fatalf("command %q has a conflicting AXI root disposition: %#v", definition.Name, disposition)
			}
			members = append(members, definition.Name)
		case len(disposition.children) > 0:
			if disposition.root || disposition.exemption != "" {
				t.Fatalf("command %q has a conflicting AXI child disposition: %#v", definition.Name, disposition)
			}
			children := map[string]bool{}
			for _, child := range disposition.children {
				if child == "" || strings.ContainsAny(child, " \t\r\n") || children[child] {
					t.Fatalf("command %q has invalid AXI children %q", definition.Name, disposition.children)
				}
				children[child] = true
				members = append(members, definition.Name+" "+child)
			}
		case disposition.exemption != "":
			if disposition.root || len(disposition.children) != 0 || strings.TrimSpace(disposition.exemption) == "" {
				t.Fatalf("command %q has a conflicting AXI exemption: %#v", definition.Name, disposition)
			}
		default:
			t.Fatalf("command %q has no AXI disposition", definition.Name)
		}
	}
	sort.Strings(members)
	return members
}

type axiEnvelopeCase struct {
	successMarker, emptyMarker, usage              string
	blocks                                         []string
	route, successArgv, deepSuccessArgv, emptyArgv []string
	setupSuccess, setupEmpty                       func(*testing.T, string)
}

// resultBlock is the table whose rows the member's success and empty cases differ over —
// the block its identifying marker names, so the two do not drift apart.
func (tc axiEnvelopeCase) resultBlock(t *testing.T) string {
	t.Helper()
	name, _, found := strings.Cut(tc.successMarker, "[")
	if !found || name == "" {
		t.Fatalf("success marker %q names no result table", tc.successMarker)
	}
	return name
}

// axiEnvelopeRows decodes a member's whole stdout as one TOON document and returns its
// result rows. Decoding the complete output is the envelope claim: bytes before, between,
// or after the blocks fail the decode, the recovered block order pins what the command
// emitted and in what order, and the help table has to be schema-correct and terminal.
func axiEnvelopeRows(t *testing.T, tc axiEnvelopeCase, result axiCommandResult) []any {
	t.Helper()
	document, err := axitest.DecodeDocument(result.stdout)
	if err != nil {
		t.Fatalf("stdout = %q, want structured TOON: %v", result.stdout, err)
	}
	if !reflect.DeepEqual(document.Blocks, tc.blocks) {
		t.Fatalf("stdout blocks = %q, want %q; stdout=%q", document.Blocks, tc.blocks, result.stdout)
	}
	if _, err := document.HelpActions(); err != nil {
		t.Fatalf("stdout = %q, want a terminal help[N]{cmd,why} envelope: %v", result.stdout, err)
	}
	rows, err := document.Rows(tc.resultBlock(t))
	if err != nil {
		t.Fatalf("stdout = %q: %v", result.stdout, err)
	}
	return rows
}

type axiCommandResult struct {
	stdout, stderr string
	code           int
}

func TestAXIRegistryBindsEachRealCommandEnvelope(t *testing.T) {
	cases := axiEnvelopeCases()
	members := axiRegistryMemberNames(t)
	if caseNames := axiEnvelopeCaseNames(cases); !reflect.DeepEqual(caseNames, members) {
		t.Fatalf("AXI envelope fixtures = %q, registry members = %q", caseNames, members)
	}
	for _, name := range members {
		tc := cases[name]
		t.Run(name+"/structured-success", func(t *testing.T) {
			root := newAXIEnvelopeRepo(t)
			tc.setupSuccess(t, root)
			result := runAXICommandAt(t, root, tc.successArgv)
			if result.code != 0 || result.stderr != "" || !strings.Contains(result.stdout, tc.successMarker) {
				t.Fatalf("result = stdout=%q stderr=%q exit=%d; want non-empty %q TOON on stdout/0", result.stdout, result.stderr, result.code, tc.successMarker)
			}
			if rows := axiEnvelopeRows(t, tc, result); len(rows) == 0 {
				t.Fatalf("stdout = %q, want at least one result row", result.stdout)
			}
		})

		t.Run(name+"/definitive-empty", func(t *testing.T) {
			root := newAXIEnvelopeRepo(t)
			tc.setupEmpty(t, root)
			result := runAXICommandAt(t, root, tc.emptyArgv)
			if result.code != 0 || result.stderr != "" || !strings.Contains(result.stdout, tc.emptyMarker) {
				t.Fatalf("result = stdout=%q stderr=%q exit=%d; want definitive %q on stdout/0", result.stdout, result.stderr, result.code, tc.emptyMarker)
			}
			if rows := axiEnvelopeRows(t, tc, result); len(rows) != 0 {
				t.Fatalf("stdout = %q, want a zero-row result table", result.stdout)
			}
		})

		t.Run(name+"/structured-refusal", func(t *testing.T) {
			result := runAXICommandAt(t, t.TempDir(), tc.successArgv)
			if result.code != 1 || result.stderr != "" || !strings.HasPrefix(result.stdout, "error:") {
				t.Fatalf("result = stdout=%q stderr=%q exit=%d; want structured stdout refusal/1", result.stdout, result.stderr, result.code)
			}
		})

		t.Run(name+"/unknown-flag", func(t *testing.T) {
			argv := append(append([]string(nil), tc.route...), "--unknown-axi-probe")
			result := runAXICommandAt(t, t.TempDir(), argv)
			if result.code != 2 || result.stderr != "" || !strings.HasPrefix(result.stdout, tc.usage) {
				t.Fatalf("result = stdout=%q stderr=%q exit=%d; want stdout usage/2 beginning %q", result.stdout, result.stderr, result.code, tc.usage)
			}
		})

		t.Run(name+"/help-spellings", func(t *testing.T) {
			for _, spelling := range []string{"--help", "-h", "help"} {
				argv := append(append([]string(nil), tc.route...), spelling)
				result := runAXICommandAt(t, t.TempDir(), argv)
				if result.code != 0 || result.stderr != "" || !strings.HasPrefix(result.stdout, tc.usage) {
					t.Errorf("%s = stdout=%q stderr=%q exit=%d; want stdout usage/0 beginning %q", spelling, result.stdout, result.stderr, result.code, tc.usage)
				}
			}
		})

		t.Run(name+"/deep-cwd", func(t *testing.T) {
			root := newAXIEnvelopeRepo(t)
			tc.setupSuccess(t, root)
			fromRoot := runAXICommandAt(t, root, tc.successArgv)
			deepArgv := tc.successArgv
			if tc.deepSuccessArgv != nil {
				deepArgv = tc.deepSuccessArgv
			}
			fromDeep := runAXICommandAt(t, filepath.Join(root, "nested", "deep"), deepArgv)
			if fromRoot.code != 0 || fromRoot.stderr != "" || !strings.Contains(fromRoot.stdout, tc.successMarker) {
				t.Fatalf("root control = %#v, want successful %q envelope", fromRoot, tc.successMarker)
			}
			if fromDeep != fromRoot {
				t.Fatalf("deep cwd result = %#v, root result = %#v", fromDeep, fromRoot)
			}
		})
	}
}

func axiEnvelopeCases() map[string]axiEnvelopeCase {
	noSetup := func(*testing.T, string) {}
	return map[string]axiEnvelopeCase{
		"anchors": {
			route: []string{"anchors"}, successArgv: []string{"anchors", ".bench/BENCH.md"}, deepSuccessArgv: []string{"anchors", "../../.bench/BENCH.md"}, emptyArgv: []string{"anchors", "unregistered.md"},
			blocks:        []string{"anchors", "help"},
			successMarker: "anchors[", emptyMarker: "anchors[0]{kind,section,needle}:\n", usage: "usage: bench anchors", setupSuccess: setupAXIAnchors, setupEmpty: noSetup,
		},
		"learnings": {
			route: []string{"learnings"}, successArgv: []string{"learnings"}, emptyArgv: []string{"learnings"},
			blocks:        []string{"learnings", "help"},
			successMarker: "learnings[1]{date,title}:\n", emptyMarker: "learnings[0]{date,title}:\n", usage: "usage: bench learnings", setupSuccess: setupAXILearnings, setupEmpty: noSetup,
		},
		"maps": {
			route: []string{"maps"}, successArgv: []string{"maps"}, emptyArgv: []string{"maps"},
			blocks:        []string{"maps", "help"},
			successMarker: "maps[1]{map,title,type,state,blockers}:\n", emptyMarker: "maps[0]{map,title,type,state,blockers}:\n", usage: "usage: bench maps", setupSuccess: setupAXIMap, setupEmpty: noSetup,
		},
		"guards": {
			route: []string{"guards"}, successArgv: []string{"guards"}, emptyArgv: []string{"guards"},
			blocks:        []string{"guards", "guard_scan", "help"},
			successMarker: "guards[1]{guard,boundary,denies,branch,provenance,currency,wired}:\n", emptyMarker: "guards[0]{guard,boundary,denies,branch,provenance,currency,wired}:\n", usage: "usage: bench guards", setupSuccess: noSetup, setupEmpty: setupAXIEmptyGuards,
		},
		"diff": {
			route: []string{"diff"}, successArgv: []string{"diff"}, emptyArgv: []string{"diff"},
			blocks:        []string{"revision", "aggregate", "files", "checkout", "whitespace", "help"},
			successMarker: "files[1]{status,path,kind}:\n", emptyMarker: "files[0]{status,path,kind}:\n", usage: "usage: bench diff", setupSuccess: setupAXIDiff, setupEmpty: noSetup,
		},
		"coverage": {
			route: []string{"coverage"}, successArgv: []string{"coverage", "fixture"}, emptyArgv: []string{"coverage", "empty"},
			blocks:        []string{"spec", "state", "rows", "help"},
			successMarker: "rows[1]{story,seam,red_signal}:\n", emptyMarker: "state: mapped\nrows[0]{story,seam,red_signal}:\n", usage: "usage: bench coverage", setupSuccess: setupAXICoverage, setupEmpty: setupAXIEmptyCoverage,
		},
		"worktree list": {
			route: []string{"worktree", "list"}, successArgv: []string{"worktree", "list"}, emptyArgv: []string{"worktree", "list"},
			blocks:        []string{"worktrees", "help"},
			successMarker: "worktrees[1]{id,label,state,source,tree,lease,landed,ignored}:\n", emptyMarker: "worktrees[0]{id,label,state,source,tree,lease,landed,ignored}:\n", usage: "usage: bench worktree list", setupSuccess: setupAXIWorktree, setupEmpty: noSetup,
		},
		"roadmap": {
			route: []string{"roadmap"}, successArgv: []string{"roadmap"}, emptyArgv: []string{"roadmap"},
			blocks:        []string{"roadmap", "board", "sequence", "drain", "help"},
			successMarker: "roadmap[", emptyMarker: "roadmap[0]{id,title,spec,spec_status,external_trigger,occurrence_count,occurrence_keys}:\n", usage: "usage: bench roadmap", setupSuccess: setupAXIRoadmap, setupEmpty: noSetup,
		},
	}
}

func axiEnvelopeCaseNames(cases map[string]axiEnvelopeCase) []string {
	names := make([]string, 0, len(cases))
	for name := range cases {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func newAXIEnvelopeRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	writeAXIFixture(t, filepath.Join(root, "nested", "deep", ".keep"), "fixture\n")
	runAXIGit(t, "-C", root, "add", ".")
	runAXIGit(t, "-C", root, "commit", "-q", "-m", "fixture base")
	base := strings.TrimSpace(runAXIGit(t, "-C", root, "rev-parse", "HEAD"))
	runAXIGit(t, "-C", root, "config", "branch.main.benchBase", base)
	return root
}

func setupAXIAnchors(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, ".bench", "BENCH.md"), "# Fixture\n")
}

func setupAXILearnings(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "capture", "learnings.md"), "# Learnings — usage journal\n\n## 2026-01-01 — fixture [open]\n")
}

func setupAXIMap(t *testing.T, root string) {
	const document = `# Fixture

Status: shaping

## Destination

Settle it.

## #1: Decide

Blocked by: none
Type: Research

### Question

What now?

### Answer

— (open)

## Not yet specified

## Spec-writer discretion

## Out of scope

## Sources
`
	writeAXIFixture(t, filepath.Join(root, "decisions", "fixture.md"), document)
}

func setupAXIEmptyGuards(t *testing.T, root string) {
	const hook = `#!/bin/sh
# bench:managed-pre-push
# name: pre-push
# boundary: pre-push
# denies: nothing (informational)
# why: fixture has no deny-capable guards
`
	writeAXIFixture(t, filepath.Join(root, ".git", "hooks", "pre-push"), hook)
}

func setupAXIDiff(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "dirty.txt"), "dirty\n")
}

func setupAXICoverage(t *testing.T, root string) {
	const spec = `# Fixture

## User stories

1. As a caller, I get one row.

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| FX1 | 1 | fixture | command | observed red | catches omission |
`
	writeAXIFixture(t, filepath.Join(root, "specs", "fixture", "spec.md"), spec)
}

func setupAXIEmptyCoverage(t *testing.T, root string) {
	const spec = `# Empty

## User stories

1. As a caller, I get no rows.

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
`
	writeAXIFixture(t, filepath.Join(root, "specs", "empty", "spec.md"), spec)
}

func setupAXIWorktree(t *testing.T, root string) {
	linked := filepath.Join(t.TempDir(), "linked")
	runAXIGit(t, "-C", root, "worktree", "add", "-q", "-b", "fixture-linked", linked)
}

func setupAXIRoadmap(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "ROADMAP.md"), "**FT1 — fixture.**\n")
}

func writeAXIFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir fixture parent: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}

func runAXIGit(t *testing.T, args ...string) string {
	t.Helper()
	out, err := gitpkg.Raw(args...)
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func runAXICommandAt(t *testing.T, cwd string, argv []string) axiCommandResult {
	t.Helper()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	}()
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	return axiCommandResult{stdout: stdout.String(), stderr: stderr.String(), code: code}
}

// keptRoutes is the surface a removal may not take with it, written down rather than
// derived. Every other routing check reads commandRegistry, so deleting a route deletes it
// from the expectation in the same edit and the check stays green; an enumeration authored
// against the reviewer's keep list is the only thing an over-broad deletion turns red. Each
// entry drives the real dispatcher and asserts the route answers its own grammar at exit 0,
// so a surviving-but-misrouted verb fails as loudly as a deleted one.
var keptRoutes = []struct {
	argv []string
	help string
}{
	{[]string{"worktree", "--help"}, "usage: bench worktree"},
	{[]string{"worktree", "reauthorize", "--help"}, "usage: bench worktree reauthorize"},
	{[]string{"gate", "--help"}, "usage: bench gate"},
	{[]string{"commit", "--help"}, "usage: bench commit"},
	{[]string{"status", "--help"}, "usage: bench status"},
	{[]string{"guards", "--help"}, "usage: bench guards"},
	{[]string{"idea", "--help"}, "usage: bench idea"},
	{[]string{"roadmap", "--help"}, "usage: bench roadmap"},
	{[]string{"spec", "implemented", "--help"}, "usage: bench spec implemented"},
	{[]string{"spec", "retire", "--help"}, "usage: bench spec retire"},
	{[]string{"spec", "history", "--help"}, "usage: bench spec history"},
}

// keptWorktreeGrammars are the pool operations the worktree family help has to keep naming.
// The family route surviving says nothing about the operations under it, and each one is
// reached only through that dispatcher, so the grammar line is where their survival shows.
var keptWorktreeGrammars = []string{
	"bench worktree create",
	"bench worktree path",
	"bench worktree exec",
	"bench worktree release",
	"bench worktree clean",
	"bench worktree reauthorize",
	"bench worktree land",
}

func TestKeptRoutesAnswerTheirOwnHelp(t *testing.T) {
	for _, kept := range keptRoutes {
		name := strings.Join(kept.argv, " ")
		t.Run(name, func(t *testing.T) {
			out, code := runKeptRoute(kept.argv)
			if code != 0 {
				t.Fatalf("%s exit = %d, want 0; output=%q", name, code, out)
			}
			if !strings.Contains(out, kept.help) {
				t.Fatalf("%s output = %q, want it to name %q", name, out, kept.help)
			}
		})
	}
}

func TestKeptWorktreeOperationsKeepTheirGrammar(t *testing.T) {
	out, code := runKeptRoute([]string{"worktree", "--help"})
	if code != 0 {
		t.Fatalf("worktree --help exit = %d, want 0; output=%q", code, out)
	}
	for _, grammar := range keptWorktreeGrammars {
		if !strings.Contains(out, grammar) {
			t.Errorf("worktree help = %q, want it to name %q", out, grammar)
		}
	}
}

// removedGrammars are the verbs the lifecycle removal took out, written down rather than
// derived: a restored route answers its own grammar instead of its family's refusal, and
// the registry a derived check would read is exactly what such a restoration changes.
var removedGrammars = []struct {
	argv  []string
	usage string
}{
	{[]string{"spec", "build", "start", "x"}, "usage: bench spec (unknown argument: build)"},
	{[]string{"spec", "build"}, "usage: bench spec (unknown argument: build)"},
	{[]string{"worktree", "recovery", "x"}, "usage: bench worktree (unknown argument: recovery)"},
}

// TestRemovedGrammarsRefuseThroughTheirFamily pins both halves of a removal: the verb
// refuses at its family's unknown-argument error, and the family help no longer advertises
// it. The worktree family's fallback is a free-form objective, so a route that merely
// stopped being routed would open a subshell named for the removed verb instead.
func TestRemovedGrammarsRefuseThroughTheirFamily(t *testing.T) {
	for _, removed := range removedGrammars {
		name := strings.Join(removed.argv, " ")
		t.Run(name, func(t *testing.T) {
			out, code := runKeptRoute(removed.argv)
			if code != 2 {
				t.Fatalf("%s exit = %d, want 2; output=%q", name, code, out)
			}
			if !strings.Contains(out, removed.usage) {
				t.Fatalf("%s output = %q, want it to name %q", name, out, removed.usage)
			}
		})
	}
	out, code := runKeptRoute([]string{"worktree", "--help"})
	if code != 0 {
		t.Fatalf("worktree --help exit = %d, want 0; output=%q", code, out)
	}
	if strings.Contains(out, "recovery") {
		t.Fatalf("worktree help = %q, want it to name no recovery grammar", out)
	}
}

// TestSkillsIndexRoutesThroughDispatch drives the verb where the wrapper sends it: the
// registry route, the module's grammar, and the reference-file bytes a later invocation
// reads back rather than a value held over from the write.
func TestSkillsIndexRoutesThroughDispatch(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, ".agents", "skills", "alpha", "SKILL.md"), "---\nname: alpha\nindex: doing alpha things\n---\n")
	// The markers are the fixture's own text: cmd/bench grades routing, so it seeds a
	// reference file rather than reaching into the module for the block's shape.
	writeAXIFixture(t, filepath.Join(root, ".bench", "BENCH-reference.md"),
		"# Reference\n\n<!-- bench:skills-index:start -->\n<!-- bench:skills-index:end -->\n")

	drift := "skills index missing entry for skill 'alpha' (regenerate: bench skills-index --write)\n"
	for _, tc := range []struct {
		name   string
		argv   []string
		stdout string
		code   int
	}{
		{"default check", []string{"skills-index"}, drift, 1},
		{"explicit check", []string{"skills-index", "--check"}, drift, 1},
		{"write", []string{"skills-index", "--write"}, "", 0},
		{"check after write", []string{"skills-index", "--check"}, "", 0},
		{"help", []string{"skills-index", "--help"}, "usage: bench skills-index [--check|--write]\n", 0},
		{"unknown flag", []string{"skills-index", "--bogus"}, "usage: bench skills-index (unknown argument: --bogus)\n", 2},
		{"conflicting modes", []string{"skills-index", "--check", "--write"}, "usage: bench skills-index [--check|--write] (--check and --write are mutually exclusive)\n", 2},
	} {
		result := runAXICommandAt(t, root, tc.argv)
		if result.stdout != tc.stdout || result.stderr != "" || result.code != tc.code {
			t.Fatalf("%s: %v = stdout=%q stderr=%q exit=%d, want stdout=%q exit=%d", tc.name, tc.argv, result.stdout, result.stderr, result.code, tc.stdout, tc.code)
		}
	}
}

// runKeptRoute joins both sinks: help lands on stdout for some grammars and stderr for
// others, and which sink a route picked is not what its callers are grading.
func runKeptRoute(argv []string) (string, int) {
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	return stdout.String() + stderr.String(), code
}

// TestSkillsIndexDistinguishesMissingGitFromOutsideRepository grades the two ways
// repository discovery fails as one partition: an unlaunchable `git` and an executed
// `git` outside a work tree must name different recovery actions, and they are asserted
// together so a missing-tool special case cannot regress the outside-repository line.
func TestSkillsIndexDistinguishesMissingGitFromOutsideRepository(t *testing.T) {
	for _, tc := range []struct {
		name   string
		path   string
		stdout string
	}{
		{"git cannot launch", "", "error: required tool is missing or not executable: git — install git and re-run\n"},
		{"git runs outside a repository", os.Getenv("PATH"), "error: not in a git repository — run inside a Bench-linked repo\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PATH", tc.path)
			result := runAXICommandAt(t, t.TempDir(), []string{"skills-index"})
			if result.stdout != tc.stdout || result.stderr != "" || result.code != 1 {
				t.Fatalf("skills-index with PATH=%q = stdout=%q stderr=%q exit=%d, want stdout=%q exit=1", tc.path, result.stdout, result.stderr, result.code, tc.stdout)
			}
		})
	}
}
