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
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
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
		attachmentDirect: {"check-agent-line", "commit", "gate-go", "gate-prose", "guard-bench-follow-on", "guard-git", "resume-clean", "session-inspect", "shift", "spec", "version", "worktree"},
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
	// recordBacked marks a query that projects a record compiled into the binary. Such a
	// query reads no disk and needs no repository. So it has no empty projection and no
	// path it can refuse, and the definitive-empty and structured-refusal subtests do not
	// apply to it. Every other subtest still runs. A query that reads the repository must
	// leave this field false.
	recordBacked bool
}

// resultBlock names the table whose rows differ between the success and empty cases.
// The name comes from the identifying marker, so the two stay in step.
func (tc axiEnvelopeCase) resultBlock(t *testing.T) string {
	t.Helper()
	name, _, found := strings.Cut(tc.successMarker, "[")
	if !found || name == "" {
		t.Fatalf("success marker %q names no result table", tc.successMarker)
	}
	return name
}

// axiEnvelopeRows decodes a member's whole stdout as one TOON document and returns its
// result rows. The decode fails on any byte before, between, or after the blocks. The
// recovered block order pins the order the command emitted, and the help table must be
// schema-correct and terminal.
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

		// The empty and refusal subtests do not apply to a record-backed query, so they
		// are not registered for one. A registered-and-skipped subtest would report a
		// missing capability, and the fixture reports a declared shape instead.
		if !tc.recordBacked {
			t.Run(name+"/definitive-empty", func(t *testing.T) {
				// A fixture that declares no empty case, and does not declare why, is an
				// incomplete fixture. The check names it rather than crashing on the nil.
				if tc.setupEmpty == nil {
					t.Fatalf("fixture %q declares no empty case; set recordBacked if the query has none", name)
				}
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
		}

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
		"consumers": {
			route: []string{"consumers"}, successArgv: []string{"consumers", "target.Symbol"}, emptyArgv: []string{"consumers", "target.Symbol"},
			blocks:        []string{"consumers", "meta", "citation", "help"},
			successMarker: "consumers[1]{file,line,via,enclosing}:\n", emptyMarker: "consumers[0]{file,line,via,enclosing}:\n", usage: "usage: bench consumers", setupSuccess: setupAXIConsumers, setupEmpty: setupAXIEmptyConsumers,
		},
		"coverage": {
			route: []string{"coverage"}, successArgv: []string{"coverage", "fixture"}, emptyArgv: []string{"coverage", "empty"},
			blocks:        []string{"spec", "state", "rows", "help"},
			successMarker: "rows[1]{story,behavior,seam}:\n", emptyMarker: "state: mapped\nrows[0]{story,behavior,seam}:\n", usage: "usage: bench coverage", setupSuccess: setupAXICoverage, setupEmpty: setupAXIEmptyCoverage,
		},
		"worktree list": {
			route: []string{"worktree", "list"}, successArgv: []string{"worktree", "list"}, emptyArgv: []string{"worktree", "list"},
			blocks:        []string{"worktrees", "help"},
			successMarker: "worktrees[1]{id,label,request,state,source,tree,lease,landed,ignored}:\n", emptyMarker: "worktrees[0]{id,label,request,state,source,tree,lease,landed,ignored}:\n", usage: "usage: bench worktree list", setupSuccess: setupAXIWorktree, setupEmpty: noSetup,
		},
		// harnesses projects a record compiled into the binary, so it declares no empty
		// case and no refusal case. See recordBacked.
		"harnesses": {
			route: []string{"harnesses"}, successArgv: []string{"harnesses"},
			blocks:        []string{"schema", "harnesses", "help"},
			successMarker: "harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}:\n", usage: "usage: bench harnesses", setupSuccess: noSetup, recordBacked: true,
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

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| FX1 | 1 | fixture | command | catches omission |
`
	writeAXIFixture(t, filepath.Join(root, "specs", "fixture", "spec.md"), spec)
}

func setupAXIEmptyCoverage(t *testing.T, root string) {
	const spec = `# Empty

## User stories

1. As a caller, I get no rows.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
`
	writeAXIFixture(t, filepath.Join(root, "specs", "empty", "spec.md"), spec)
}

// setupAXIConsumers plants the smallest real Go module the resolver can answer over: one
// declaration and one reference to it. The loader seam is package-internal, so this
// surface test drives the real go/packages path rather than a stub.
func setupAXIConsumers(t *testing.T, root string) {
	setupAXIEmptyConsumers(t, root)
	writeAXIFixture(t, filepath.Join(root, "consumer", "consumer.go"),
		"package consumer\n\nimport \"example.com/axifixture/target\"\n\nfunc Direct() { target.Symbol() }\n")
}

// setupAXIEmptyConsumers plants the same module without a consumer, so the resolved
// symbol has zero references and the response is the definitive empty table.
func setupAXIEmptyConsumers(t *testing.T, root string) {
	writeAXIFixture(t, filepath.Join(root, "go.mod"), "module example.com/axifixture\n\ngo 1.25\n")
	writeAXIFixture(t, filepath.Join(root, "target", "target.go"), "package target\n\nfunc Symbol() {}\n")
}

func setupAXIWorktree(t *testing.T, root string) {
	linked := filepath.Join(t.TempDir(), "linked")
	runAXIGit(t, "-C", root, "worktree", "add", "-q", "-b", "fixture-linked", linked)
}

// setupAXIRoadmap writes the split board the roadmap surfaces read: one index heading
// line and the row file that owns its detail.
func setupAXIRoadmap(t *testing.T, root string) {
	const heading = "**FT1 — fixture.**"
	roadmaptest.WriteSplitBoard(t, root, heading+"\n", map[string]string{"FT1.md": heading + "\n"})
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
	return runAXICommandAsAt(t, cwd, "bench", argv)
}

func runAXICommandAsAt(t *testing.T, cwd, executable string, argv []string) axiCommandResult {
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
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: executable}.Run(argv)
	return axiCommandResult{stdout: stdout.String(), stderr: stderr.String(), code: code}
}

// TestWorktreeLandNeverConsultsTheInvokedExecutable drives the real dispatcher in a
// repository that declares Go build inputs. The stable-owner landing runs entirely
// under the invoked process: the registry seam must hand the landing no executable
// proof, so the refusals are repository proofs, the sentinel is never named, and the
// planted executable never runs. A dispatcher that reintroduced the verification or
// rebuild path would name or execute the sentinel and go red here.
func TestWorktreeLandNeverConsultsTheInvokedExecutable(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "scripts", "go-build.inputs"), "build_script=scripts/go-build.sh\n")
	ran := filepath.Join(t.TempDir(), "sentinel-ran")
	sentinel := filepath.Join(t.TempDir(), "invoked-bench")
	writeAXIFixture(t, sentinel, "#!/bin/sh\n: > "+ran+"\n")
	if err := os.Chmod(sentinel, 0o755); err != nil {
		t.Fatal(err)
	}
	argv := []string{"worktree", "land", "--request", "r", "--base", "b", "--source-tip", "s", "--spec", "x", "-m", "land", root}
	result := runAXICommandAsAt(t, root, sentinel, argv)
	if result.code != 1 || result.stderr != "" || !strings.HasPrefix(result.stdout, "refused{detail=") {
		t.Fatalf("land = stdout=%q stderr=%q exit=%d, want repository-proof refusals", result.stdout, result.stderr, result.code)
	}
	if strings.Contains(result.stdout, sentinel) || strings.Contains(result.stdout, "untrusted") {
		t.Fatalf("land consulted the invoked executable: %q", result.stdout)
	}
	if _, err := os.Stat(ran); !os.IsNotExist(err) {
		t.Fatalf("land executed the invoked repository executable: %v", err)
	}
}

// keptRoutes names the surface a removal must not take with it. This list stands apart
// from commandRegistry, so deleting a route there does not also delete the expectation
// here; only an enumeration against the reviewer's keep list turns an over-broad deletion
// red. Each entry drives the real dispatcher and asserts that the route answers its own
// grammar at exit 0, so a surviving but misrouted verb fails as loudly as a deleted one.
var keptRoutes = []struct {
	argv []string
	help string
}{
	{[]string{"worktree", "--help"}, "usage: bench worktree"},
	{[]string{"worktree", "create", "--help"}, "usage: bench worktree create"},
	{[]string{"worktree", "reauthorize", "--help"}, "usage: bench worktree reauthorize"},
	{[]string{"worktree", "merge", "--help"}, "usage: bench worktree merge"},
	{[]string{"gate", "--help"}, "usage: bench gate"},
	{[]string{"commit", "--help"}, "usage: bench commit"},
	{[]string{"status", "--help"}, "usage: bench status"},
	{[]string{"guards", "--help"}, "usage: bench guards"},
	{[]string{"idea", "--help"}, "usage: bench idea"},
	{[]string{"learning", "--help"}, "usage: bench learning"},
	{[]string{"retro", "--help"}, "usage: bench retro"},
	{[]string{"roadmap", "--help"}, "usage: bench roadmap"},
	{[]string{"spec", "retire", "--help"}, "usage: bench spec retire"},
	{[]string{"spec", "history", "--help"}, "usage: bench spec history"},
}

// keptWorktreeGrammars are the pool operations the worktree family help must keep naming.
// The surviving family route says nothing about the operations under it. Each operation is
// reachable only through that dispatcher, so the grammar line is where its survival shows.
var keptWorktreeGrammars = []string{
	// WF42: the full create grammar, so the family help pins `--from` rather than the
	// bare verb.
	usage.WorktreeCreate,
	"bench worktree path",
	"bench worktree exec",
	usage.WorktreeShow,
	usage.WorktreeBuild,
	"bench worktree release",
	"bench worktree clean",
	"bench worktree reauthorize",
	"bench worktree merge",
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

func TestWorktreeHelpNamesCleanGrammar(t *testing.T) {
	out, code := runKeptRoute([]string{"worktree", "--help"})
	const grammar = "bench worktree clean [--discard-ignored] [--discard-branch] [--full] (<path> | --landed) [--apply <fingerprint>] | bench worktree clean --discard-branch --unclaimed [--apply <fingerprint> | --apply-current]"
	if code != 0 || !strings.Contains(out, grammar) {
		t.Fatalf("worktree --help exit=%d output=%q, want %q", code, out, grammar)
	}
}

func TestBareWorktreeRefusesBeforeItAcquiresAssignment(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	t.Setenv("BENCH_HOME", t.TempDir())
	before, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}

	result := runAXICommandAt(t, root, []string{"worktree"})
	if result.stdout != usage.WorktreeUsage() || result.stderr != "" || result.code != 2 {
		t.Fatalf("bare worktree = stdout=%q stderr=%q exit=%d, want usage on stdout and exit 2", result.stdout, result.stderr, result.code)
	}
	after, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("bare worktree assignments = %#v, want %#v", after, before)
	}
}

func TestUnknownWorktreeSubcommandRefusesBeforeItAcquiresAssignment(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	before, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}

	result := runAXICommandAt(t, root, []string{"worktree", "unknown"})
	want := toon.Usage("bench worktree", "unknown") + "\n"
	if result.stdout != "" || result.stderr != want || result.code != 2 {
		t.Fatalf("unknown worktree subcommand = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 2", result.stdout, result.stderr, result.code, want)
	}
	after, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("unknown worktree subcommand assignments = %#v, want %#v", after, before)
	}
}

func TestUnknownWorktreeFlagRefusesBeforeItAcquiresAssignment(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	before, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}

	result := runAXICommandAt(t, root, []string{"worktree", "--unknown"})
	want := toon.Usage("bench worktree", "--unknown") + "\n"
	if result.stdout != "" || result.stderr != want || result.code != 2 {
		t.Fatalf("unknown worktree flag = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 2", result.stdout, result.stderr, result.code, want)
	}
	after, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("unknown worktree flag assignments = %#v, want %#v", after, before)
	}
}

func TestWorktreeCreateKeepsParserFirstDispatch(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	before, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}

	result := runAXICommandAt(t, root, []string{"worktree", "create", "--unknown"})
	want := toon.Usage("bench worktree create", "--unknown") + "\n"
	if result.stdout != "" || result.stderr != want || result.code != 2 {
		t.Fatalf("worktree create unknown flag = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 2", result.stdout, result.stderr, result.code, want)
	}
	after, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("worktree create unknown flag assignments = %#v, want %#v", after, before)
	}
}

// removedGrammars are the verbs the lifecycle removal took out. This list stands apart
// from the registry, because a restored route answers its own grammar instead of its
// family's refusal, and a check derived from the registry would miss such a restoration.
var removedGrammars = []struct {
	argv  []string
	usage string
}{
	{[]string{"spec", "build", "start", "x"}, "usage: bench spec (unknown argument: build)"},
	{[]string{"spec", "build"}, "usage: bench spec (unknown argument: build)"},
	{[]string{"worktree", "recovery", "x"}, "usage: bench worktree (unknown argument: recovery)"},
}

// TestRemovedGrammarsRefuseThroughTheirFamily pins both halves of a removal. The verb
// refuses at its family's unknown-argument error, and the family help no longer advertises
// it. The worktree family's fallback is a free-form objective, so a route that merely
// stopped being routed would instead open a subshell named for the removed verb.
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

// TestSkillsIndexRoutesThroughDispatch drives the verb where the wrapper sends it. It
// covers the registry route, the module's grammar, and the reference-file bytes a later
// invocation reads back rather than a value held over from the write.
func TestSkillsIndexRoutesThroughDispatch(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, ".agents", "skills", "alpha", "SKILL.md"), "---\nname: alpha\nindex: doing alpha things\n---\n")
	// The markers are the fixture's own text. cmd/bench grades routing, so it seeds a
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

// runKeptRoute joins both sinks. Help lands on stdout for some grammars and stderr for
// others, and callers do not grade which sink a route picked.
func runKeptRoute(argv []string) (string, int) {
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	return stdout.String() + stderr.String(), code
}

// TestSkillsIndexDistinguishesMissingGitFromOutsideRepository grades the two ways
// repository discovery fails as one partition. An unlaunchable `git` and a `git` executed
// outside a work tree must name different recovery actions. The test asserts both
// together, so a missing-tool special case cannot regress the outside-repository line.
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

// TestHelpSpecRowsNameRetireAndHistoryOnly pins FA9: the retired subcommand keeps no
// help row, while the two surviving spec rows stay advertised.
func TestHelpSpecRowsNameRetireAndHistoryOnly(t *testing.T) {
	help := renderCommandHelp()
	if strings.Contains(help, "bench spec implemented") {
		t.Errorf("bench help still advertises a retired subcommand:\n%s", help)
	}
	for _, want := range []string{"bench spec retire <slug>", "bench spec history <slug>"} {
		if !strings.Contains(help, want) {
			t.Errorf("bench help is missing %q:\n%s", want, help)
		}
	}
}

func TestTestHelpNamesOnlyRunnableFocusedForms(t *testing.T) {
	help := renderCommandHelp()
	if !strings.Contains(help, "bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name>") {
		t.Errorf("bench help is missing the focused test grammar:\n%s", help)
	}
}

func TestFocusedHelpNamesCheckAsEvidenceWithoutVerdict(t *testing.T) {
	help := renderCommandHelp()
	if !strings.Contains(help, "bench test [--full] --check <name>") {
		t.Fatalf("bench help is missing the named-check form:\n%s", help)
	}
	if strings.Contains(strings.ToLower(help), "test --check <name>  gate") {
		t.Fatalf("bench help presents a focused check as a verdict:\n%s", help)
	}
}

// TestWorktreeRoutesKeepTheirBytesFromASubdirectory grades WF11. The four verbs that now
// receive an explicit root must answer the same bytes from a subdirectory as from the
// repository root, because the boundary resolves the root and the verb no longer reads
// the working directory.
//
// The two invocations run against a fixture repository, not the live checkout.
// The population the answer names is the one registered in the root the boundary
// resolves, so a worktree that arrives in the real repository between the runs
// cannot reach the fixture root and cannot flip the comparison. The "a worktree
// arrives between the runs" row proves that root scope: it registers a linked
// worktree in a second fixture repository between the two runs, and the bytes
// stay the same.
func TestWorktreeRoutesKeepTheirBytesFromASubdirectory(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	subdirectory := filepath.Join(root, "nested", "deep")
	for _, argv := range [][]string{{"worktree", "clean", "--help"}, {"worktree", "list"}} {
		for _, row := range []struct {
			name    string
			between func(*testing.T)
		}{
			{name: "quiet", between: func(*testing.T) {}},
			{name: "a worktree arrives between the runs", between: func(t *testing.T) {
				setupAXIWorktree(t, newAXIEnvelopeRepo(t))
			}},
		} {
			t.Run(strings.Join(argv, " ")+"/"+row.name, func(t *testing.T) {
				atRoot := runAXICommandAt(t, root, argv)
				row.between(t)
				below := runAXICommandAt(t, subdirectory, argv)
				if atRoot != below {
					t.Fatalf("%v from %s = %+v, from %s = %+v, want the same bytes and exit code", argv, root, atRoot, subdirectory, below)
				}
			})
		}
	}
}

func TestWorktreeCleanHelpUsesItsGrammar(t *testing.T) {
	result := runAXICommandAt(t, t.TempDir(), []string{"worktree", "clean", "--help"})
	want := "usage: " + usage.WorktreeClean + "\n"
	if result.stdout != want || result.stderr != "" || result.code != 0 {
		t.Fatalf("worktree clean --help = stdout=%q stderr=%q exit=%d, want stdout=%q stderr=\"\" exit=0", result.stdout, result.stderr, result.code, want)
	}
}

// TestWorktreeHelpNamesTheExecGateForm pins WF14. The worktree usage trailer names the
// gate as a worktree-native exec form, and it names no raw wrapper path. An appended
// line beside the old trailer would keep the raw path alive, so the assertion demands
// the trailer be the last line and demands the whole text hold no bin/bench.sh.
func TestWorktreeHelpNamesTheExecGateForm(t *testing.T) {
	out, code := runKeptRoute([]string{"worktree", "--help"})
	if code != 0 {
		t.Fatalf("worktree --help exit = %d, want 0; output=%q", code, out)
	}
	const trailer = "bench worktree exec <target> -- bench gate"
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if last := strings.TrimSpace(lines[len(lines)-1]); last != trailer {
		t.Errorf("worktree help last line = %q, want %q", last, trailer)
	}
	if strings.Contains(out, "bin/bench.sh") {
		t.Errorf("worktree help = %q, want no bin/bench.sh", out)
	}
}
