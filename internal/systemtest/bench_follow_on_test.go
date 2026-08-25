//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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
		{"command end marker", "command -- bench help | touch <marker>"},
		{"nohup end marker", "nohup -- bench help | touch <marker>"},
		{"timeout options and end marker", "timeout -- 5 bench help | touch <marker>"},
		{"xargs end marker", "xargs -- bench help | touch <marker>"},
		{"FOG34 deeper relative", "../../bin/bench.sh help | touch <marker>"},
		{"FOG35 absolute", filepath.Join(owner.kit, "bin", "bench.sh") + " help | touch <marker>"},
		{"FOG37 output redirect", "bench help > <marker>"},
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
		{"FOG38 heredoc redirect", "bench gate <<'EOF'\ninput\nEOF"},
	} {
		result := run(repo, tc.command)
		if result.code != 2 || !strings.Contains(result.stderr, "BLOCKED: Bench response is bounded, complete, and self-contained") {
			t.Errorf("%s = (%d, %q, %q), want refusal", tc.name, result.code, result.stdout, result.stderr)
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
	alias := filepath.Join(repo, "kit-command")
	if err := os.Symlink(filepath.Join(repo, "bin", "bench.sh"), alias); err != nil {
		t.Fatal(err)
	}
	blocked("FOG36 live symlink", repo, "./kit-command help | touch <marker>")
	for _, command := range []string{
		"bench gate --fresh",
		"bench worktree exec label -- bash -lc 'go test && go vet'",
		"rg bench AGENTS.md",
		"printf hi > bench",
		"cat <<'EOF'\nbench gate | tail -20\nEOF",
	} {
		if result := run(repo, command); result.code != 0 || strings.Contains(result.stderr, "BLOCKED:") {
			t.Errorf("allowed command %q = (%d, %q, %q)", command, result.code, result.stdout, result.stderr)
		}
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
