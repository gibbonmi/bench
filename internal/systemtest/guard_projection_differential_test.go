//go:build system

package systemtest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// differentialGolden holds the tuples that the pre-migration executable and kit produced.
const differentialGolden = "testdata/guard-projection-family-d.golden"

// differentialCommands are the Bash commands of family D. Each one runs through the Git and
// follow-on core entries and through their hooks. <pool> names a worktree pool target.
var differentialCommands = []struct{ family, command string }{
	{"P", "env -u X bench help | cat"},
	{"P", "timeout -s KILL -k 1 5 bench help | cat"},
	{"P", "xargs -n 1 bench help | cat"},
	{"P", "command -V bench | cat"},
	{"P", "command -v git"},
	{"P", "X=1 git reset --hard"},
	{"P", "env X=1 git commit -m x"},
	{"P", "nohup git reset --hard"},
	{"P", "timeout 5 git reset --hard"},
	{"P", "xargs git push"},
	{"W", "sh -c 'git reset --hard'"},
	{"W", "bash -c 'git reset --hard'"},
	{"W", "zsh -c 'git reset --hard'"},
	{"W", "/bin/bash -c 'git reset --hard'"},
	{"W", "/usr/bin/zsh -lc 'git reset --hard'"},
	{"W", "bash -lc 'bench help | cat'"},
	{"W", "bash -c"},
	{"W", "bash -c -- 'git reset --hard'"},
	{"W", "bash -c -c 'git status'"},
	{"W", "bash -c -c 'git reset --hard'"},
	{"W4", "bash -c 'echo ok' -c 'git reset --hard'"},
	{"W4", "bash -c 'echo ok' -c 'bench help | cat'"},
	{"W4", `bash -c "bash -c 'git reset --hard'"`},
	{"X", "bench worktree exec X --env A=1 -- git commit -m 'two words'"},
	{"X", "bench worktree exec X git reset --hard"},
	{"X", "bench worktree exec X --"},
	{"X", "./dist/bench worktree exec X -- git reset --hard"},
	{"X", "mybench worktree exec X -- git reset --hard"},
	{"X", "bench worktree exec X -- git status"},
	{"X3", "bench worktree exec X -- bash -c 'git commit -m x'"},
	{"X3", "bash -c 'bench worktree exec X -- git commit -m x'"},
	{"X3", "echo bench worktree exec X -- git commit"},
	{"X3", `bash -c "bash -c 'git commit -m x'"`},
	{"R", "bench commitment override --assignment x --reason y"},
	{"R", "bench commitment --in primary override"},
	{"R", "bench worktree exec t -- bench commitment override"},
	{"R", "bench worktree exec t -- bash -lc 'bench commitment override --assignment x --reason y'"},
	{"R", "bash -c 'bench commitment override --assignment x --reason y'"},
	{"R", "env X=1 ./bin/bench.sh commitment override --assignment x --reason y"},
	{"R", "bench commitment show"},
	{"R", "bench commitment block --outcome A --reason override"},
	{"R", "rg 'bench commitment override' ."},
	{"R", "echo commitment override"},
	{"R", `bash -c "bash -c 'bench commitment override'"`},
	{"R", "bench worktree exec t -- bench worktree exec u -- bench commitment override"},
	{"R", "bench worktree exec t bench commitment override"},
	{"R", "x worktree exec t -- bench commitment override"},
	{"A", "printf hi;git reset --hard"},
	{"A", "printf hi&&git reset --hard"},
	{"A", "printf hi\ngit reset --hard"},
	{"A", occurrenceCommand},
	{"A", "printf hi;git push --force\necho \"x"},
	{"A", "printf hi;git push --force\nprintf x\\"},
	{"G", "git -C /other push origin topic"},
	{"G", "git --git-dir=/x push origin topic"},
	{"G", "git --work-tree=/x push origin topic"},
	{"G", "git -C /other commit -m x"},
	{"G", "git --git-dir /x reset --hard"},
	{"G", "git -c core.x=y reset --hard"},
	{"G", "git --git-dir=x stash list"},
	{"D", "git status"},
	{"D", "bench gate --fresh"},
	{"D", "./kit-command help | cat"},
	{"D", "cd <pool>"},
	{"D", "cd <pool> && bench gate"},
	{"D", "git -C <pool> log"},
	{"D", "bench gate <<'EOF'\ninput\nEOF"},
	{"D", "cat <<'EOF'\ngit reset --hard\nEOF"},
	{"D", "bench worktree exec label -- cat <<'EOF'\nbody line\nEOF"},
	{"D", "bench gate 2>&1"},
	{"D", "bench help > out"},
	{"D", "git status > out"},
	{"D", "echo 'a && git reset --hard'"},
	{"D", `echo "bench help | cat"`},
	{"D", "'bench' gate"},
}

// differentialEnvelopes are ordered raw envelope bytes. Each one runs through all three core
// entries and all three hooks. /repo names the primary checkout, and /outside names no repository.
var differentialEnvelopes = []struct{ family, envelope string }{
	{"E0", ``},
	{"E0", `not json`},
	{"E0", `{}`},
	{"E0", `{"tool_input":null}`},
	{"E0", `{"tool_input":{}}`},
	{"E0", `{"tool_input":{"command":null}}`},
	{"E0", `{"tool_input":{"command":7}}`},
	{"E0", `{"tool_input":{"command":["ls"]}}`},
	{"E0", `{"tool_input":{"command":""}}`},
	{"E0", `{"tool_input":{"command":"\u0000"}}`},
	{"EP", `{"tool_input":{"COMMAND":"git reset --hard"}}`},
	{"EP", `{"tool_input":{"command":"git reset --hard","command":null}}`},
	{"EP", `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":null,"cwd":"/repo"}`},
	{"EP", `{"tool_input":{"command":7,"command":"git reset --hard"}}`},
	{"EP", `{"cwd":"/repo","cwd":null,"tool_input":{"file_path":"README.md"}}`},
	{"EP", `{"TOOL_INPUT":{"command":"git reset --hard"}}`},
	{"EP", `{"CWD":"/repo","TOOL_INPUT":{"file_path":"README.md"}}`},
	{"EP", `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"TOOL_INPUT":{"other":1},"cwd":"/repo"}`},
	{"EP", `{"tool_input":{"command":"git reset --hard","file_path":"README.md"},"tool_input":{},"cwd":"/repo"}`},
	{"EP", `{"tool_input":{"command":"git status","file_path":"old.md"},"tool_input":{"command":"git reset --hard","file_path":"README.md"},"cwd":"/repo"}`},
	{"EP", `{"cwd":7,"cwd":"/repo","tool_input":{"file_path":"README.md"}}`},
	{"F", `{"tool_input":{"file_path":"/repo/README.md"}}`},
	{"F", `{"cwd":"/repo","tool_input":{"file_path":"README.md"}}`},
	{"F", `{"tool_input":{"file_path":"README.md"}}`},
	{"F", `{"cwd":7,"tool_input":{"file_path":"README.md"}}`},
	{"F", `{"cwd":null,"tool_input":{"file_path":"README.md"}}`},
	{"F", `{"tool_input":{"file_path":"/repo/new.md"}}`},
	{"F", `{"tool_input":{"file_path":"/repo/missing/x.md"}}`},
	{"F", `{"tool_input":{"file_path":"/outside/x.md"}}`},
	{"F", `{"tool_input":{"file_path":"/outside/missing/x.md"}}`},
}

// differentialEntry is one actual guard entry that reads an envelope on stdin.
type differentialEntry struct {
	name string
	run  func(envelope string) processResult
}

// TestGuardProjectionDifferentialFamilies is GP46. Family D and the decoder matrix keep the
// exit, stdout, and stderr that the pre-migration executable recorded in the golden.
func TestGuardProjectionDifferentialFamilies(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(differentialGolden)
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	got := differentialTuples(t)
	if len(got) != len(want) {
		t.Errorf("GP46 tuple count = %d, want %d", len(got), len(want))
	}
	// The first failure line names every differing member, so a bounded report shows each one.
	var members, tuples []string
	for i := range min(len(got), len(want)) {
		if got[i] != want[i] {
			// The quoted input ends the member name, and an input can hold its own " | ".
			member, _, _ := strings.Cut(want[i], `" | `)
			members = append(members, member+`"`)
			tuples = append(tuples, "got  "+got[i], "want "+want[i])
		}
	}
	if len(members) > 0 {
		t.Errorf("GP46 members differ from the golden: %s\n%s", strings.Join(members, "; "), strings.Join(tuples, "\n"))
	}
}

// differentialTuples runs every member through its entries. It returns one normalized line for
// each member, which holds the member's entry tuples in a fixed order.
func differentialTuples(t *testing.T) []string {
	t.Helper()
	scratch := t.TempDir()
	primary := differentialPrimary(t)
	pool := filepath.Join(scratch, "pool")
	outside := filepath.Join(scratch, "outside")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(pool, "worktrees", "bench-1", strings.Repeat("a", 32)+"-"+strings.Repeat("b", 32))
	place := strings.NewReplacer("/repo", primary, "/outside", outside, "<pool>", target)
	// The kit differs between the recording and the candidate, and each fixture root is temporary.
	normalize := strings.NewReplacer(owner.kit, "<kit>", primary, "<fixture>", scratch, "<fixture>", owner.root, "<fixture>", filepath.Dir(scratch), "<fixture>")
	live := []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit, "BENCH_HOME=" + pool}
	core := func(guard string) differentialEntry {
		return differentialEntry{"core " + guard, func(envelope string) processResult {
			if err := owner.observeSelected(); err != nil {
				t.Fatal(err)
			}
			return owner.runWithInput(primary, live, envelope, owner.selected.path, guard)
		}}
	}
	hook := func(name, path string) differentialEntry {
		return differentialEntry{"hook " + name, func(envelope string) processResult {
			return projectionRunHook(t, path, primary, live, envelope)
		}}
	}
	// Each guard pairs its hook with its core entry and with the envelope its degraded hook reads.
	guards := []differentialGuard{
		{"git", gitGuardHookPath(), "guard-git", hookEnvelope(t, "Bash", "git reset --hard")},
		{"follow-on", filepath.Join(owner.kit, ".bench", "hooks", "block-bench-follow-on.sh"), "guard-bench-follow-on", hookEnvelope(t, "Bash", "bench help | cat")},
		{"write", projectionWriteHookPath(), "guard-file-write", place.Replace(`{"cwd":"/repo","tool_input":{"file_path":"README.md"}}`)},
	}
	commandEntries := []differentialEntry{core(guards[0].guard), core(guards[1].guard), hook(guards[0].name, guards[0].hook), hook(guards[1].name, guards[1].hook)}
	envelopeEntries := append(append([]differentialEntry{}, commandEntries...), core(guards[2].guard), hook(guards[2].name, guards[2].hook))
	var members []string
	tuples := map[string][]string{}
	record := func(family, input, entry string, result processResult) {
		member := fmt.Sprintf("%s %q", family, input)
		if _, seen := tuples[member]; !seen {
			members = append(members, member)
		}
		tuples[member] = append(tuples[member], fmt.Sprintf("%s exit=%d stdout=%q stderr=%q", entry, result.code, result.stdout, result.stderr))
	}
	for _, member := range differentialCommands {
		envelope := hookEnvelope(t, "Bash", place.Replace(member.command))
		for _, entry := range commandEntries {
			record(member.family, member.command, entry.name, entry.run(envelope))
		}
	}
	for _, member := range differentialEnvelopes {
		for _, entry := range envelopeEntries {
			record(member.family, member.envelope, entry.name, entry.run(place.Replace(member.envelope)))
		}
	}
	for _, line := range differentialDegradedHooks(t, guards) {
		record("H", line.name, line.hook, line.result)
	}
	lines := make([]string, 0, len(members))
	for _, member := range members {
		lines = append(lines, normalize.Replace(member+" | "+strings.Join(tuples[member], " | ")))
	}
	return lines
}

// differentialPrimary builds a real primary checkout with a tracked README.md, the repository
// wrapper, and a symlink alias of that wrapper.
func differentialPrimary(t *testing.T) string {
	t.Helper()
	primary := journeyRepo(t, "differential [primary]-")
	if err := os.WriteFile(filepath.Join(primary, "README.md"), []byte("readme\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	systemGit(t, primary, "add", "README.md")
	systemGit(t, primary, systemCommitArgs("-qm", "readme")...)
	installWrapper(t, primary)
	if err := os.Symlink(filepath.Join(primary, "bin", "bench.sh"), filepath.Join(primary, "kit-command")); err != nil {
		t.Fatal(err)
	}
	return primary
}

// differentialGuard is one guard: its name, its hook, its core entry, and its degraded envelope.
type differentialGuard struct{ name, hook, guard, envelope string }

// differentialDegraded is one degraded hook result.
type differentialDegraded struct {
	name, hook string
	result     processResult
}

// differentialDegradedHooks runs each hook with a missing library, a missing core, an erroring
// core, a stale core, and a missing platform binary with a recovery read.
func differentialDegradedHooks(t *testing.T, guards []differentialGuard) []differentialDegraded {
	t.Helper()
	private := []string{"PATH=" + privateToolPath(t, "git", "cat", "dirname")}
	outside := t.TempDir()
	recovery := hookEnvelope(t, "Bash", "cat .github/workflows/gate.yml")
	var results []differentialDegraded
	for _, h := range guards {
		missingBinary := plantWrapperAnswer(t, "no installed binary for this platform", 127)
		for _, c := range []struct {
			name, hook, dir, envelope string
			env                       []string
		}{
			{"missing library", hookWithoutLibrary(t, h.hook), outside, h.envelope, private},
			{"missing core", h.hook, outside, h.envelope, private},
			{"core error", h.hook, plantWrapperAnswer(t, "core analyzer failed", 3), h.envelope, nil},
			{"stale core", h.hook, plantWrapperAnswer(t, `bench: unknown subcommand: "`+h.guard+`"`, 2), h.envelope, []string{"BENCH_KIT=" + owner.kit}},
			{"missing binary", h.hook, missingBinary, h.envelope, nil},
			{"missing binary recovery", h.hook, missingBinary, recovery, nil},
		} {
			results = append(results, differentialDegraded{c.name, "hook " + h.name, projectionRunHook(t, c.hook, c.dir, c.env, c.envelope)})
		}
	}
	return results
}
