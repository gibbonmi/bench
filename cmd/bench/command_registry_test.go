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

	"github.com/gibbonmi/bench/internal/intent"
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

// TestWorktreeLandNeverConsultsTheInvokedExecutable drives the real dispatcher in a
// repository that declares Go build inputs. The stable-owner landing runs entirely
// under the invoked process: the registry seam must hand the landing no executable
// proof, so the refusals are repository proofs, the sentinel is never named, and the
// planted executable never runs. A dispatcher that reintroduced the verification or
// rebuild path would name or execute the sentinel and go red here.
func TestWorktreeLandNeverConsultsTheInvokedExecutable(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	declareBuildInputs(t, root)
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
	{[]string{"worktree", "reset", "--help"}, "usage: bench worktree reset"},
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
	const grammar = "bench worktree clean [--discard-ignored] [--discard-branch] [--full] (<path> | --landed | --target <target>...) [--apply <fingerprint>] | bench worktree clean --discard-branch --unclaimed [--apply <fingerprint> | --apply-current]"
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

func TestWorktreeShellRunsAndReleasesItsAssignment(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	t.Setenv("BENCH_HOME", t.TempDir())
	t.Setenv("SHELL", "true")

	result := runAXICommandAt(t, root, []string{"worktree", "shell"})
	if result.code != 0 || !strings.Contains(result.stderr, "🪵 worktree: ") {
		t.Fatalf("worktree shell = stdout=%q stderr=%q exit=%d, want successful shell journey", result.stdout, result.stderr, result.code)
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(assignments) != 0 {
		t.Fatalf("worktree shell left assignments = %#v, want none after normal exit", assignments)
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
		if withoutTreeTargetLine("skills-index", result.stdout) != tc.stdout || result.stderr != "" || result.code != tc.code {
			t.Fatalf("%s: %v = stdout=%q stderr=%q exit=%d, want stdout=%q exit=%d", tc.name, tc.argv, result.stdout, result.stderr, result.code, tc.stdout, tc.code)
		}
	}
}

// runKeptRoute joins both sinks. Help lands on stdout for some grammars and stderr for
// others, and callers do not grade which sink a route picked.
func runKeptRoute(argv []string) (string, int) {
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	return withoutTreeRow(stdout.String()) + stderr.String(), code
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
