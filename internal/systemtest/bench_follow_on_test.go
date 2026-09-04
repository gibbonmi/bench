//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/freshness"
)

// plantWrapperAnswer builds a repository whose Bench wrapper answers one fixed stderr
// line at one fixed exit code. Both hook tests stand a core in this way, because a stale
// core and a genuine refusal differ only in that line and in nothing else the shim sees.
// The line must hold no single quote.
func plantWrapperAnswer(t *testing.T, line string, code int) string {
	t.Helper()
	repo := t.TempDir()
	if result := owner.runAt(repo, nil, "git", "init", "-q"); result.code != 0 {
		t.Fatal(result.stderr)
	}
	wrapper := filepath.Join(repo, ".bench", "bin", "bench.sh")
	if err := os.MkdirAll(filepath.Dir(wrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	script := "#!/bin/sh\ncat >/dev/null\necho '" + line + "' >&2\nexit " + strconv.Itoa(code) + "\n"
	if err := os.WriteFile(wrapper, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return repo
}

// staleCoreAnswer is what a binary built before the guard subcommand existed says.
const staleCoreAnswer = `bench: unknown subcommand: "guard-bench-follow-on"`

func TestBenchFollowOnHookProcess(t *testing.T) {
	repo := owner.repos[0]
	installWrapper(t, repo)
	hook := filepath.Join(owner.kit, ".bench", "hooks", "block-bench-follow-on.sh")
	if _, err := os.Stat(hook); err != nil {
		t.Fatal(err)
	}

	run := func(dir, command string) processResult {
		t.Helper()
		shell, err := exec.LookPath("bash")
		if err != nil {
			t.Fatal(err)
		}
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		envelope := `{"tool_name":"Bash","tool_input":{"command":` + shellQuoteJSON(command) + `}}`
		return owner.runWithInput(dir, []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}, envelope, shell, hook)
	}
	// blocked asserts the fixed refusal sentence on every refusal. (Coverage row G11.)
	blocked := func(name, dir, command string) {
		t.Helper()
		marker := filepath.Join(t.TempDir(), "marker")
		result := run(dir, strings.ReplaceAll(command, "<marker>", marker))
		if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: Bench response is bounded, complete, and self-contained") {
			t.Errorf("%s = (%d, %q, %q), want refusal", name, result.code, result.stdout, result.stderr)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Errorf("%s marker exists: %v", name, err)
		}
	}

	if result := run(repo, "bench gate --fresh"); result.code != 0 || strings.Contains(result.stderr, "BLOCKED:") {
		t.Fatalf("FOG01 bare Bench = (%d, %q, %q), want allowed", result.code, result.stdout, result.stderr)
	}
	for _, tc := range []struct{ name, command string }{
		{"FOG02 pipeline", "bench help | touch <marker>"},
		{"FOG03 stderr pipeline", "bench help |& touch <marker>"},
		{"FOG06 and", "bench help && touch <marker>"},
		{"FOG07 semicolon", "bench help; touch <marker>"},
		{"FOG08 newline", "bench help\ntouch <marker>"},
		{"FOG26 or", "bench help extra || touch <marker>"},
		{"FOG27 env prefix", "env X=1 bench help | touch <marker>"},
		{"FOG40 assignment prefix", "X=1 bench help | touch <marker>"},
		{"leading descriptor redirect", "2>/dev/null bench help | touch <marker>"},
		{"env options and end marker", "env -i -- X=1 bench help | touch <marker>"},
		{"env unset option", "env -u X bench help | touch <marker>"},
		{"command end marker", "command -- bench help | touch <marker>"},
		{"nohup end marker", "nohup -- bench help | touch <marker>"},
		{"timeout options and end marker", "timeout -- 5 bench help | touch <marker>"},
		{"timeout signal and kill after", "timeout -s KILL -k 1 5 bench help | touch <marker>"},
		{"xargs end marker", "xargs -- bench help | touch <marker>"},
		{"xargs max args", "xargs -n 1 bench help | touch <marker>"},
		{"FOG34 deeper relative", "../../bin/bench.sh help | touch <marker>"},
		{"FOG35 absolute", filepath.Join(owner.kit, "bin", "bench.sh") + " help | touch <marker>"},
		{"FOG37 output redirect", "bench help > <marker>"},
		{"wrapper outer output redirect", "bash -lc 'bench help' > <marker>"},
	} {
		dir := repo
		if tc.name == "FOG34 deeper relative" {
			dir = filepath.Join(owner.kit, "internal", "systemtest")
		}
		blocked(tc.name, dir, tc.command)
	}
	for _, tc := range []struct{ name, command string }{
		{"FOG04 descriptor redirect", "bench gate 2>&1"},
		{"FOG05 input redirect", "bench gate </dev/null"},
		{"FOG09 worktree redirect", "bench worktree exec label -- bench gate --fresh > /tmp/bench-follow-on-output"},
		{"FOG10 worktree pipeline", "bench worktree exec label -- bench gate --fresh | cat"},
		{"FOG18 wrapper", "bash -lc 'bench help | cat'"},
		{"FOG18 wrapper outer pipeline", "bash -lc 'bench help' | cat"},
		{"wrapper outer and", "bash -lc 'bench help' && cat"},
		// FOG38 is coverage row G6: a heredoc on any head but exec still refuses.
		{"FOG38 heredoc redirect", "bench gate <<'EOF'\ninput\nEOF"},
	} {
		result := run(repo, tc.command)
		if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: Bench response is bounded, complete, and self-contained") {
			t.Errorf("%s = (%d, %q, %q), want refusal", tc.name, result.code, result.stdout, result.stderr)
		}
	}
	// The refusal line names the Bench segment and the operator that caused it, so
	// an agent reads the cause at the hook seam. (Coverage rows G8 and G9.)
	for _, tc := range []struct{ row, command, tail string }{
		{"G8", "cat a && echo x; bench maps", "segment=bench maps operator=;"},
		{"G9", "bench gate 2>&1", "segment=bench gate operator=2>&1"},
	} {
		result := run(repo, tc.command)
		if result.code != 2 || !strings.HasSuffix(strings.TrimRight(result.stderr, "\n"), tc.tail) {
			t.Errorf("%s = (%d, %q, %q), want stderr ending with %q", tc.row, result.code, result.stdout, result.stderr, tc.tail)
		}
	}

	for _, envelope := range []string{
		`{"tool_name":"Bash","tool_input":{"command":"bench help \u0026\u0026 touch nowhere"}}`,
		`{"tool_name":"Bash","tool_input":{"command":"bench help \u003e nowhere"}}`,
	} {
		result := owner.runWithInput(repo, []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}, envelope, shellPath(t), hook)
		if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: Bench response is bounded, complete, and self-contained") {
			t.Errorf("Unicode envelope %q = (%d, %q, %q), want refusal", envelope, result.code, result.stdout, result.stderr)
		}
	}
	// A route into the worktree pool refuses with its own message, because the one command
	// form for a Bench worktree is `bench worktree exec`. A `cd`, an assignment, and a git
	// repository option are the three routes. The denial runs before the follow-on
	// verdict, so a pool reference with a follow-on names the reference. A read of a file
	// under the pool, and a pool path the `bench worktree` grammar takes, stay allowed.
	poolHome := t.TempDir()
	target := filepath.Join(poolHome, "worktrees", "bench-1", strings.Repeat("a", 32)+"-"+strings.Repeat("b", 32))
	poolEnv := []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit, "BENCH_HOME=" + poolHome}
	for _, tc := range []struct{ name, command string }{
		{"FOG41 bare cd into the pool", "cd " + target},
		{"FOG42 cd into the pool then a follow-on", "cd " + target + " && go test ./..."},
		{"FOG43 cd into the pool then a Bench call", "cd " + target + " && bench gate"},
		{"FOG44 assignment then git -C the variable", "W=" + target + "; git -C $W log"},
		{"FOG45 git -C into the pool", "git -C " + target + " log"},
		{"FOG46 git --work-tree into the pool", "git --work-tree " + target + " status"},
		{"FOG47 export of the pool path", "export W=" + target},
	} {
		envelope := `{"tool_name":"Bash","tool_input":{"command":` + shellQuoteJSON(tc.command) + `}}`
		result := owner.runWithInput(repo, poolEnv, envelope, shellPath(t), hook)
		if result.code != 2 || !strings.Contains(result.stderr, `never cd, assign, or git -C into the pool path. target=`+target) {
			t.Errorf("%s = (%d, %q, %q), want the pool refusal naming the target", tc.name, result.code, result.stdout, result.stderr)
		}
		if strings.Contains(result.stderr, "BLOCKED: Bench response is bounded") {
			t.Errorf("%s = %q, want the pool denial rather than the follow-on refusal", tc.name, result.stderr)
		}
	}
	for _, command := range []string{"cd /tmp", `cd "$W"`, "bench worktree exec \"x\" -- sh -c 'cd sub && go test ./...'", "cat " + target + "/x", "bench worktree release --request r " + target} {
		envelope := `{"tool_name":"Bash","tool_input":{"command":` + shellQuoteJSON(command) + `}}`
		if result := owner.runWithInput(repo, poolEnv, envelope, shellPath(t), hook); result.code != 0 || strings.Contains(result.stderr, "BLOCKED:") {
			t.Errorf("allowed command %q = (%d, %q, %q)", command, result.code, result.stdout, result.stderr)
		}
	}

	alias := filepath.Join(repo, "kit-command")
	if err := os.Symlink(filepath.Join(repo, "bin", "bench.sh"), alias); err != nil {
		t.Fatal(err)
	}
	blocked("FOG36 live symlink", repo, "./kit-command help | touch <marker>")
	// The last two rows are the exec exception: a non-Bench step after the exec
	// child, and a heredoc that feeds the child. (Coverage row G12.)
	for _, command := range []string{
		"bench gate --fresh",
		"bench worktree exec label -- bash -lc 'go test && go vet'",
		"rg bench AGENTS.md",
		"printf hi > bench",
		"cat <<'EOF'\nbench gate | tail -20\nEOF",
		"command -V bench | cat",
		"command -v bench | cat",
		"bench worktree exec label -- cp a b; cp b a",
		"bench worktree exec label -- cat <<'EOF'\nbody line\nEOF",
	} {
		if result := run(repo, command); result.code != 0 || strings.Contains(result.stderr, "BLOCKED:") {
			t.Errorf("allowed command %q = (%d, %q, %q)", command, result.code, result.stdout, result.stderr)
		}
	}

	// A core that refuses genuinely also exits 2, and its refusal must reach the caller
	// unchanged. A shim that read every exit 2 as a stale answer would fail open on the
	// pool denial, which is the one refusal a non-Bench call can earn. (Coverage row BF24.)
	refusingRepo := plantWrapperAnswer(t, benchguard.PoolReferenceMessage("/pool/target"), 2)
	refused := owner.runWithInput(refusingRepo, []string{"BENCH_KIT=" + owner.kit}, `{"tool_input":{"command":"ls"}}`, shellPath(t), hook)
	if refused.code != 2 || !strings.Contains(refused.stderr, "BLOCKED:") {
		t.Errorf("BF24 genuine refusal = (%d, %q, %q), want exit 2 with the refusal forwarded", refused.code, refused.stdout, refused.stderr)
	}
	if strings.Contains(refused.stderr, "bash scripts/go-build.sh") {
		t.Errorf("BF24 genuine refusal = %q, want no stale-binary reading", refused.stderr)
	}
}

func TestBenchFollowOnHookDegradedRim(t *testing.T) {
	hook := filepath.Join(owner.kit, ".bench", "hooks", "block-bench-follow-on.sh")
	run := func(hook, dir string, overrides []string, envelope string) processResult {
		t.Helper()
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		return owner.runWithInput(dir, overrides, envelope, shellPath(t), hook)
	}
	repo := owner.repos[1]
	installWrapper(t, repo)
	for _, envelope := range []string{
		`{`,
		`{"tool_input":{}}`,
		`{"tool_input":{"command":1}}`,
		`{"tool_input":{"command":""}}`,
		`{"tool_input":{"command":"\u0000"}}`,
	} {
		result := run(hook, repo, []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}, envelope)
		if result.code != 0 || !strings.Contains(result.stderr, "WARNING: block-bench-follow-on: unreadable command field") {
			t.Errorf("FOG28-32 %q = (%d, %q, %q), want warning allow", envelope, result.code, result.stdout, result.stderr)
		}
	}
	missingRoot := t.TempDir()
	missingHook := filepath.Join(missingRoot, "hooks", "block-bench-follow-on.sh")
	data, err := os.ReadFile(hook)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(missingHook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(missingHook, data, 0o755); err != nil {
		t.Fatal(err)
	}
	result := run(missingHook, t.TempDir(), []string{"PATH=" + privateToolPath(t, "git", "cat")}, `{"tool_input":{"command":"bench help | cat"}}`)
	if result.code != 0 || !strings.Contains(result.stderr, "wrapper resolver missing") {
		t.Errorf("FOG19 missing resolver = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	result = run(hook, t.TempDir(), []string{"PATH=" + privateToolPath(t, "git", "cat", "dirname")}, `{"tool_input":{"command":"bench help | cat"}}`)
	if result.code != 0 || !strings.Contains(result.stderr, "bench core not found") {
		t.Errorf("FOG19/33 missing wrapper or binary = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}

	brokenRepo := t.TempDir()
	if result := owner.runAt(brokenRepo, nil, "git", "init", "-q"); result.code != 0 {
		t.Fatal(result.stderr)
	}
	badWrapper := filepath.Join(brokenRepo, ".bench", "bin", "bench.sh")
	if err := os.MkdirAll(filepath.Dir(badWrapper), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, exit string
	}{
		{"FOG20 core error", "7"},
		{"FOG33 missing platform binary", "127"},
	} {
		if err := os.WriteFile(badWrapper, []byte("#!/bin/sh\nexit "+tc.exit+"\n"), 0o755); err != nil {
			t.Fatal(err)
		}
		result = run(hook, brokenRepo, nil, `{"tool_input":{"command":"bench help | cat"}}`)
		if result.code != 0 || !strings.Contains(result.stderr, "bench core errored (exit "+tc.exit+")") {
			t.Errorf("%s = (%d, %q, %q)", tc.name, result.code, result.stdout, result.stderr)
		}
	}

	// A stale core answers `unknown subcommand` at exit 2, and the shim splits its
	// posture on that answer. A non-Bench call passes, so the shell stays usable while
	// the binary is rebuilt; a Bench call refuses, so no verb runs stale. Both carry
	// the one rebuild sentence, which the test reads from its own owner rather than
	// from a second copy. (Coverage rows BF22 and BF23.)
	staleRepo := plantWrapperAnswer(t, staleCoreAnswer, 2)
	action := freshness.RebuildAction(owner.kit)
	for _, tc := range []struct {
		row, command string
		code         int
		refuses      bool
	}{
		{"BF22 stale core passes a non-Bench call", "ls", 0, false},
		{"BF23 stale core refuses a Bench call", "bench gate", 2, true},
		{"BF23 stale core refuses a quoted Bench call", "'bench' gate", 2, true},
	} {
		envelope := `{"tool_input":{"command":` + shellQuoteJSON(tc.command) + `}}`
		result := run(hook, staleRepo, []string{"BENCH_KIT=" + owner.kit}, envelope)
		if result.code != tc.code || !strings.Contains(result.stderr, action) {
			t.Errorf("%s = (%d, %q, %q), want exit %d and the rebuild sentence %q", tc.row, result.code, result.stdout, result.stderr, tc.code, action)
		}
		if refused := strings.Contains(result.stderr, "BLOCKED:"); refused != tc.refuses {
			t.Errorf("%s stderr = %q, want a refusal %t", tc.row, result.stderr, tc.refuses)
		}
	}

	// The shim reads the command out of the envelope's tool_input object, so a field
	// that follows that object never joins the command, and a JSON-escaped operator
	// classifies exactly as its literal spelling does. Both rows exercise the shim's
	// own reader, which decides only under a stale core. (Coverage rows BF22 and BF23.)
	for _, tc := range []struct {
		row, envelope string
		code          int
		refuses       bool
	}{
		{"BF22 a trailing cwd stays out of the command", `{"tool_name":"Bash","tool_input":{"command":"ls"},"cwd":"/home/u/bench"}`, 0, false},
		{"BF23 an escaped operator refuses like its literal spelling", `{"tool_input":{"command":"ls \u0026\u0026 bench gate"}}`, 2, true},
	} {
		result := run(hook, staleRepo, []string{"BENCH_KIT=" + owner.kit}, tc.envelope)
		if result.code != tc.code || !strings.Contains(result.stderr, action) {
			t.Errorf("%s = (%d, %q, %q), want exit %d and the rebuild sentence %q", tc.row, result.code, result.stdout, result.stderr, tc.code, action)
		}
		if refused := strings.Contains(result.stderr, "BLOCKED:"); refused != tc.refuses {
			t.Errorf("%s stderr = %q, want a refusal %t", tc.row, result.stderr, tc.refuses)
		}
	}
}

func shellPath(t *testing.T) string {
	t.Helper()
	shell, err := exec.LookPath("bash")
	if err != nil {
		t.Fatal(err)
	}
	return shell
}

func shellQuoteJSON(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", `"`, `\\"`, "\n", "\\n")
	return `"` + replacer.Replace(value) + `"`
}
