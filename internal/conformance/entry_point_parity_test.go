package conformance

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

// parityShimDirs are the two directories the check enumerates. Everything else in the
// table is a fixed path, because no directory holds a family of front doors.
var parityShimDirs = []string{".bench/hooks", ".bench/adapters"}

// checkEntryPointParity holds every entry point to one core. The runtime rows prove that a
// shim reaches the registry and answers like the direct verb. The static rows prove the
// same route where a run would need a pool, a release, or a CI runner. A root without the
// wrapper is not graded at runtime, because a shim with no core on PATH takes its own rim.
func checkEntryPointParity(root string) []string {
	var diags []string
	diags = append(diags, parityEnumerationDiags(root)...)
	diags = append(diags, parityStaticDiags(root)...)
	diags = append(diags, parityCoverageDiags(root)...)
	diags = append(diags, parityRuntimeDiags(root)...)
	return diags
}

// parityEnumerationDiags names each shim on disk the table does not carry. Only a regular
// file is an entry: a link or a FIFO at a shim path is not something to run.
func parityEnumerationDiags(root string) []string {
	var diags []string
	for _, dir := range parityShimDirs {
		entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.Type().IsRegular() {
				continue
			}
			rel := dir + "/" + entry.Name()
			if _, listed := entryPointParity[rel]; !listed {
				diags = append(diags, fmt.Sprintf("entry-point-parity: %s reaches no parity row; give it a row naming the registry command it calls", rel))
			}
		}
	}
	sort.Strings(diags)
	return diags
}

// parityStaticDiags grades every static row's text. An absent or refused entry is red
// rather than skipped: the row is the only grader of that route, so a silent skip would
// leave the route ungraded the moment the file is deleted or replaced by a link.
func parityStaticDiags(root string) []string {
	var diags []string
	for _, rel := range parityRelPaths() {
		row := entryPointParity[rel]
		if !row.isStatic() {
			continue
		}
		text, readable := parityRead(filepath.Join(root, filepath.FromSlash(rel)))
		if !readable {
			diags = append(diags, fmt.Sprintf("entry-point-parity: %s is absent or unreadable, and the registry command %q is therefore ungraded", rel, row.command))
			continue
		}
		if !row.must.MatchString(text) {
			diags = append(diags, fmt.Sprintf("entry-point-parity: %s does not name the registry command %q", rel, row.command))
		}
	}
	return diags
}

// parityCoverageDiags grades the registry's internal commands against the table. A public
// command answers to an agent who reads the help inventory; an internal one is reached
// only through a shim, so it needs a row or a stated reason.
func parityCoverageDiags(root string) []string {
	path := filepath.Join(root, "cmd", "bench", "main.go")
	body := readIfExists(path)
	if body == "" {
		return nil
	}
	entries, err := parseCommandRegistry(path, body)
	if err != nil {
		return []string{"entry-point-parity: cmd/bench/main.go cannot be parsed for the command registry: " + err.Error()}
	}
	reached := map[string]bool{}
	for _, row := range entryPointParity {
		reached[row.command] = true
	}
	var diags []string
	for _, entry := range entries {
		if !parityInternalCommand(entry) {
			continue
		}
		if reached[entry.name] || parityExemptCommands[entry.name] != "" {
			continue
		}
		diags = append(diags, fmt.Sprintf("entry-point-parity: registry command %q is reached by no parity row and carries no exemption reason", entry.name))
	}
	return diags
}

func parityInternalCommand(entry commandRegistryEntry) bool {
	fields := entry.fields["Inventory"]
	if len(fields) != 1 {
		return false
	}
	value, ok := fields[0].(*ast.Ident)
	return ok && value.Name == "internalInventory"
}

// parityRead returns an entry's text through the no-follow classifier. A link, a FIFO, or
// bytes outside the bounded read is refused rather than followed. An absent entry reads
// the same way a refused one does, and the caller decides what that means for its row.
func parityRead(path string) (string, bool) {
	classified := bounds.ClassifyNoFollow(path)
	if classified.State != bounds.StateParsed {
		return "", false
	}
	return string(classified.Data), true
}

func parityRelPaths() []string {
	rels := make([]string, 0, len(entryPointParity))
	for rel := range entryPointParity {
		rels = append(rels, rel)
	}
	sort.Strings(rels)
	return rels
}

// parityRuntimeDiags runs each runtime entry through the stub wrapper under
// BENCH_COMMAND_OBSERVE=1. The runs share one routed temp repository, so every row reads
// the same binding and no row depends on the graded root's own state.
func parityRuntimeDiags(root string) []string {
	realBench := filepath.Join(root, "bin", "bench.sh")
	if !exists(realBench) {
		return nil
	}
	bindir, cleanup, err := adapterStubDir(realBench)
	if err != nil {
		return []string{"entry-point-parity setup failed: " + err.Error()}
	}
	defer cleanup()
	repo, cleanupRepo, err := tempGitRepoWithLines(openCodeBinding)
	if err != nil {
		return []string{"entry-point-parity setup failed: " + err.Error()}
	}
	defer cleanupRepo()

	wrapper := filepath.Join(bindir, "bench")
	base, cleanupEnv, err := conformanceSubprocessEnv(bindir)
	if err != nil {
		return []string{"entry-point-parity setup failed: " + err.Error()}
	}
	defer cleanupEnv()
	base = append(base, "BENCH_COMMAND_OBSERVE=1")
	var diags []string
	for _, rel := range parityRelPaths() {
		row := entryPointParity[rel]
		entry := filepath.Join(root, filepath.FromSlash(rel))
		if row.isStatic() || !exists(entry) {
			continue
		}
		diags = append(diags, parityRowDiags(repo, rel, row, entry, wrapper, append(base, row.env...))...)
	}
	return diags
}

func parityRowDiags(repo, rel string, row parityRow, entry, wrapper string, env []string) []string {
	shim := runWithInputEnv(repo, env, row.input, "bash", entry)
	if shim == nil {
		return []string{"entry-point-parity: " + rel + " did not run"}
	}
	var diags []string
	if !parityObserved(shim.Stderr, row.command) {
		diags = append(diags, fmt.Sprintf("entry-point-parity: %s does not reach the registry: no command-registry:%s on stderr", rel, row.command))
	}
	if !row.compare {
		return diags
	}
	argv := make([]string, 0, len(row.direct)+1)
	argv = append(argv, wrapper)
	for _, arg := range row.direct {
		if arg == parityWrapperToken {
			arg = wrapper
		}
		argv = append(argv, arg)
	}
	direct := runWithInputEnv(repo, env, row.input, argv...)
	if direct == nil {
		return append(diags, "entry-point-parity: direct "+row.command+" did not run")
	}
	if shim.ExitCode != direct.ExitCode {
		diags = append(diags, fmt.Sprintf("entry-point-parity: %s exits %d and the direct %s exits %d", rel, shim.ExitCode, row.command, direct.ExitCode))
	}
	if !strings.HasSuffix(shim.Stdout, direct.Stdout) {
		diags = append(diags, fmt.Sprintf("entry-point-parity: %s stdout does not end with the direct %s stdout", rel, row.command))
	}
	return diags
}

// parityObserved reads only the observed id line. A shim's own warnings share stderr, and
// a check that matched the whole stream would accept a warning that quoted the id.
func parityObserved(stderr, command string) bool {
	want := "command-registry:" + command
	for _, line := range strings.Split(stderr, "\n") {
		if strings.TrimRight(line, "\r") == want {
			return true
		}
	}
	return false
}

// parityHonestStatics is the shortest text each static row accepts. A fixture that means
// to grade something else plants all of them, so the only red it can take is its own.
var parityHonestStatics = map[string]string{
	".bench/hooks/worktree-lifecycle.sh":   "#!/usr/bin/env bash\ncmd=x\nexec \"$cmd\" worktree-hook \"$@\"\n",
	"scripts/release-preflight.sh":         "#!/usr/bin/env bash\nbinary=x\nexec \"$binary\" release-preflight \"$@\"\n",
	".github/workflows/native-runtime.yml": "jobs:\n  run:\n    steps:\n      - run: bash scripts/release-preflight.sh\n",
	".agents/commands/bench.md":            "# bench\n\nRun `bench status --route` and take its one row.\n",
}

// TestEntryPointParityNamesAnAbsentStaticEntry proves a static row reds when its subject
// leaves the tree. A skip would retire the row's whole route in silence.
func TestEntryPointParityNamesAnAbsentStaticEntry(t *testing.T) {
	diags := checkEntryPointParity(throwawayRoot{}.build(t))
	want := `entry-point-parity: scripts/release-preflight.sh is absent or unreadable, and the registry command "release-preflight" is therefore ungraded`
	if !containsDiagnostic(diags, want) {
		t.Fatalf("an absent static entry did not bite with %q:\n%s", want, strings.Join(diags, "\n"))
	}
}

func TestEntryPointParityNamesAShimOutsideTheTable(t *testing.T) {
	root := throwawayRoot{files: map[string]string{".bench/hooks/extra.sh": "#!/usr/bin/env bash\nexit 0\n"}}.build(t)

	diags := checkEntryPointParity(root)

	if !containsDiagnostic(diags, ".bench/hooks/extra.sh reaches no parity row") {
		t.Fatalf("a shim outside the table was not named:\n%s", strings.Join(diags, "\n"))
	}
}

// TestEntryPointParityStaticRowsBite grades the two rows whose subject is text. Each case
// keeps the honest form green, so the mutation, not the row's shape, is what reds.
func TestEntryPointParityStaticRowsBite(t *testing.T) {
	if diags := checkEntryPointParity(parityStaticRoot(t, "", "")); len(diags) != 0 {
		t.Fatalf("the honest static tree is not green:\n%s", strings.Join(diags, "\n"))
	}
	tests := []struct {
		name, rel, from, to, want string
	}{
		{
			name: "CI exec line", rel: "scripts/release-preflight.sh",
			from: "release-preflight ", to: "release-preflight-2 ",
			want: `scripts/release-preflight.sh does not name the registry command "release-preflight"`,
		},
		{
			name: "front-door verb", rel: ".agents/commands/bench.md",
			from: "--route", to: "--routes",
			want: `.agents/commands/bench.md does not name the registry command "status"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mutated := strings.Replace(parityHonestStatics[tt.rel], tt.from, tt.to, 1)
			if mutated == parityHonestStatics[tt.rel] {
				t.Fatalf("the mutation of %s changed nothing", tt.rel)
			}
			diags := checkEntryPointParity(parityStaticRoot(t, tt.rel, mutated))
			if !containsDiagnostic(diags, tt.want) {
				t.Fatalf("the mutated %s did not bite with %q:\n%s", tt.rel, tt.want, strings.Join(diags, "\n"))
			}
		})
	}
}

// parityStaticRoot plants every honest static entry, then overwrites the one a case
// mutates. A root that planted only the subject would red on the other rows too, and the
// case could no longer prove which fault it caught.
func parityStaticRoot(t *testing.T, rel, content string) string {
	t.Helper()
	files := map[string]string{}
	for staticRel, honest := range parityHonestStatics {
		files[staticRel] = honest
	}
	if rel != "" {
		files[rel] = content
	}
	return throwawayRoot{files: files}.build(t)
}

// parityEnvPrefixedPreflight is the second honest form of the CI exec line: the script
// hands the binary it built to the child phases. The row accepts both forms, so this tree
// grades the alternation rather than leave it to the live root alone.
const parityEnvPrefixedPreflight = "#!/usr/bin/env bash\nbinary=x\nexec env BENCH_RUN_BINARY=\"$binary\" \"$binary\" release-preflight \"$@\"\n"

// TestEntryPointParityAcceptsTheEnvPrefixedPreflight proves the env prefix stays green and
// the verb name still carries the verdict.
func TestEntryPointParityAcceptsTheEnvPrefixedPreflight(t *testing.T) {
	rel := "scripts/release-preflight.sh"
	if diags := checkEntryPointParity(parityStaticRoot(t, rel, parityEnvPrefixedPreflight)); len(diags) != 0 {
		t.Fatalf("the env-prefixed exec line is not green:\n%s", strings.Join(diags, "\n"))
	}

	mutated := strings.Replace(parityEnvPrefixedPreflight, "release-preflight ", "release-preflight-2 ", 1)
	want := `scripts/release-preflight.sh does not name the registry command "release-preflight"`
	if diags := checkEntryPointParity(parityStaticRoot(t, rel, mutated)); !containsDiagnostic(diags, want) {
		t.Fatalf("the mutated env-prefixed exec line did not bite with %q:\n%s", want, strings.Join(diags, "\n"))
	}
}
