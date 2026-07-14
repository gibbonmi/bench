package runtime

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestRuntimeShiftAdapterContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "bench shift adapter preflight contract", testShiftAdapterPreflight)
	contract.RunParallel(t, "bench shift adapter single-argument contract", testShiftAdapterSingleArgument)
	contract.RunParallel(t, "reference adapter files contract", testReferenceAdapterFiles)
	contract.RunParallel(t, "bench worktree lifecycle surface parity contract", testWorktreeLifecycleSurfaceParity)
	contract.RunParallel(t, "claude worktree event lifecycle contract", testClaudeWorktreeEventLifecycle)
	contract.RunParallel(t, "claude oversized worktree create event contract", testClaudeOversizedWorktreeCreateEvent)
	contract.RunParallel(t, "claude worktree removal failure stays locked contract", testClaudeWorktreeRemovalFailureStaysLocked)
	contract.RunParallel(t, "worktree lifecycle safety family contract", testWorktreeLifecycleSafetyFamily)
}

func testClaudeOversizedWorktreeCreateEvent(t *testing.T) {
	f := onMainFixture(t)
	f.Bench("link").RequireExit(0)
	createCommand := configuredClaudeWorktreeCommand(t, f, "WorktreeCreate")
	env := map[string]string{
		"BENCH_HOME":         filepath.Join(f.Root, ".bench-home"),
		"CLAUDE_PROJECT_DIR": f.Root,
	}
	createJSON := `{"session_id":"oversized-session","cwd":` + string(mustJSON(t, f.Root)) + `,"name":"must not create"}` + "\n"
	oversized := createJSON + strings.Repeat(" ", (1<<20)-len(createJSON)) + "malformed"
	created := contract.RunAtWithInput(t, f, f.Root, env, oversized, createCommand.Command, createCommand.Args...)
	if created.ExitCode == 0 {
		t.Fatalf("oversized Claude WorktreeCreate event succeeded\nstdout:\n%s\nstderr:\n%s", created.Stdout, created.Stderr)
	}
	if created.Stdout != "" {
		t.Fatalf("oversized Claude WorktreeCreate returned a path: %q", created.Stdout)
	}
	registrations, err := benchgit.Worktrees(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	if len(registrations) != 1 || registrations[0].Path != f.Root {
		t.Fatalf("oversized Claude WorktreeCreate changed registrations: %#v", registrations)
	}
	assignments, err := intent.Assignments(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	if len(assignments) != 0 {
		t.Fatalf("oversized Claude WorktreeCreate wrote assignments: %#v", assignments)
	}
	if branches := f.Git("for-each-ref", "--format=%(refname)", "refs/heads/bench/assign/").Stdout; branches != "" {
		t.Fatalf("oversized Claude WorktreeCreate wrote assignment branches: %q", branches)
	}
	common := strings.TrimSpace(f.Git("rev-parse", "--path-format=absolute", "--git-common-dir").Stdout)
	markers, err := filepath.Glob(filepath.Join(common, "worktrees", "*", "bench-owner"))
	if err != nil {
		t.Fatal(err)
	}
	if len(markers) != 0 {
		t.Fatalf("oversized Claude WorktreeCreate wrote owner markers: %#v", markers)
	}
}

func testWorktreeLifecycleSurfaceParity(t *testing.T) {
	wantLifecycle := normalizedLifecycleState{true, true, true, true, true, true}
	wantCleanup := normalizedCleanupState{true, true, true, true}
	var gotLifecycle []normalizedLifecycleState
	var gotCleanup []normalizedCleanupState

	interactive := onMainFixture(t)
	interactiveHome := filepath.Join(interactive.Root, ".bench-home")
	ready, proceed := filepath.Join(t.TempDir(), "ready"), filepath.Join(t.TempDir(), "proceed")
	shell := filepath.Join(t.TempDir(), "surface-shell")
	contract.WriteExecutableAbs(t, shell, "#!/usr/bin/env bash\nprintf '%s\\n' \"$PWD\" > \"$BENCH_SURFACE_READY\"\nwhile [[ ! -e \"$BENCH_SURFACE_PROCEED\" ]]; do sleep 0.01; done\n")
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(benchPath(t), "worktree", "surface parity")
	cmd.Dir = interactive.Root
	cmd.Env = surfaceEnv(interactive, map[string]string{
		"BENCH_HOME": interactiveHome, "SHELL": shell,
		"BENCH_SURFACE_READY": ready, "BENCH_SURFACE_PROCEED": proceed,
	})
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	path := waitForSurfacePath(t, ready, cmd)
	state, assignment := inspectLifecycleState(t, interactive, path)
	gotLifecycle = append(gotLifecycle, state)
	if err := os.WriteFile(proceed, []byte("continue\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("interactive lifecycle failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	gotCleanup = append(gotCleanup, inspectCleanupState(t, interactive, assignment))

	for _, tc := range []struct {
		name   string
		linked bool
	}{
		{name: "real CLI"},
		{name: "linked by-path CLI", linked: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := onMainFixture(t)
			launcher := benchPath(t)
			if tc.linked {
				f.Bench("link").RequireExit(0)
				launcher = filepath.Join(f.Root, ".bench", "bin", "bench.sh")
			}
			env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
			created := contract.RunAt(t, f, f.Root, env, launcher, "worktree", "create", "--request", "surface-request", "--label", "surface parity")
			created.RequireExit(0)
			path := worktreeCreatePath(t, created.Stdout)
			state, assignment := inspectLifecycleState(t, f, path)
			gotLifecycle = append(gotLifecycle, state)
			released := contract.RunAt(t, f, f.Root, env, launcher, "worktree", "release", "--request", "surface-request", path)
			released.RequireExit(0)
			contract.RequireContains(t, released.Stdout, "complete")
			gotCleanup = append(gotCleanup, inspectCleanupState(t, f, assignment))
		})
	}

	for i, got := range gotLifecycle {
		if !reflect.DeepEqual(got, wantLifecycle) {
			t.Fatalf("surface %d lifecycle state = %#v, want normalized marker/assignment/lock/branch parity %#v", i, got, wantLifecycle)
		}
	}
	for i, got := range gotCleanup {
		if !reflect.DeepEqual(got, wantCleanup) {
			t.Fatalf("surface %d cleanup state = %#v, want normalized release parity %#v", i, got, wantCleanup)
		}
	}
}

func testClaudeWorktreeEventLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts []contract.FixtureOption
	}{
		{name: "ordinary project root"},
		{name: "project root with spaces", opts: []contract.FixtureOption{contract.WithSpacePath()}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := onMainFixture(t, tc.opts...)
			f.Bench("link").RequireExit(0)
			requireWorktreeLifecycleSharedResolver(t, f)
			createCommand := configuredClaudeWorktreeCommand(t, f, "WorktreeCreate")
			removeCommand := configuredClaudeWorktreeCommand(t, f, "WorktreeRemove")
			env := map[string]string{
				"BENCH_HOME":         filepath.Join(f.Root, ".bench-home"),
				"CLAUDE_PROJECT_DIR": f.Root,
			}
			createJSON := `{"session_id":"session-22","cwd":` + string(mustJSON(t, f.Root)) + `,"name":"official event label","hook_event_name":"WorktreeCreate"}` + "\n"
			created := contract.RunAtWithInput(t, f, f.Root, env, createJSON, createCommand.Command, createCommand.Args...)
			created.RequireExit(0)
			path := strings.TrimSuffix(created.Stdout, "\n")
			if !filepath.IsAbs(path) || created.Stdout != path+"\n" || strings.Contains(path, "\n") {
				t.Fatalf("Claude WorktreeCreate stdout = %q, want exactly one absolute path plus newline", created.Stdout)
			}
			state, assignment := inspectLifecycleState(t, f, path)
			wantLifecycle := normalizedLifecycleState{true, true, true, true, true, true}
			if !reflect.DeepEqual(state, wantLifecycle) {
				t.Fatalf("Claude WorktreeCreate lifecycle state = %#v, want %#v", state, wantLifecycle)
			}
			removeJSON := `{"session_id":"session-22","worktree_path":` + string(mustJSON(t, path)) + `,"hook_event_name":"WorktreeRemove"}` + "\n"
			removed := contract.RunAtWithInput(t, f, f.Root, env, removeJSON, removeCommand.Command, removeCommand.Args...)
			removed.RequireExit(0)
			wantCleanup := normalizedCleanupState{true, true, true, true}
			if got := inspectCleanupState(t, f, assignment); !reflect.DeepEqual(got, wantCleanup) {
				t.Fatalf("Claude WorktreeRemove cleanup state = %#v, want routed release %#v", got, wantCleanup)
			}
			ledgerBefore, err := intent.Read(f.Root)
			if err != nil {
				t.Fatal(err)
			}
			refsBefore := f.Git("for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/", "refs/heads/bench/assign/").Stdout
			replayed := contract.RunAtWithInput(t, f, f.Root, env, removeJSON, removeCommand.Command, removeCommand.Args...)
			replayed.RequireExit(0)
			ledgerAfter, err := intent.Read(f.Root)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(ledgerAfter.CleanupReceipts, ledgerBefore.CleanupReceipts) {
				t.Fatalf("identical completed WorktreeRemove replay changed receipts\nbefore: %#v\nafter:  %#v", ledgerBefore.CleanupReceipts, ledgerAfter.CleanupReceipts)
			}
			refsAfter := f.Git("for-each-ref", "--format=%(refname) %(objectname)", "refs/bench/", "refs/heads/bench/assign/").Stdout
			if refsAfter != refsBefore {
				t.Fatalf("identical completed WorktreeRemove replay changed refs\nbefore: %q\nafter:  %q", refsBefore, refsAfter)
			}
		})
	}
}

func requireWorktreeLifecycleSharedResolver(t *testing.T, f contract.Fixture) {
	t.Helper()
	hook := f.ReadFile(".bench/hooks/worktree-lifecycle.sh")
	for _, needle := range []string{"../lib/resolve-bench.sh", `. "$lib"`, "bench_resolve_wrapper"} {
		if !strings.Contains(hook, needle) {
			t.Fatalf("linked worktree lifecycle hook does not use the shared wrapper resolver: missing %q", needle)
		}
	}
	for _, duplicate := range []string{"for candidate in", `"$root/.bench/bin/bench.sh"`, `"$root/bin/bench.sh"`} {
		if strings.Contains(hook, duplicate) {
			t.Fatalf("linked worktree lifecycle hook duplicates wrapper resolution instead of using the shared resolver: found %q", duplicate)
		}
	}
}

func testClaudeWorktreeRemovalFailureStaysLocked(t *testing.T) {
	f := onMainFixture(t)
	f.Bench("link").RequireExit(0)
	createCommand := configuredClaudeWorktreeCommand(t, f, "WorktreeCreate")
	removeCommand := configuredClaudeWorktreeCommand(t, f, "WorktreeRemove")
	env := map[string]string{
		"BENCH_HOME":         filepath.Join(f.Root, ".bench-home"),
		"CLAUDE_PROJECT_DIR": f.Root,
	}
	createJSON := `{"session_id":"session-23","cwd":` + string(mustJSON(t, f.Root)) + `,"name":"failure retention"}` + "\n"
	created := contract.RunAtWithInput(t, f, f.Root, env, createJSON, createCommand.Command, createCommand.Args...)
	created.RequireExit(0)
	path := strings.TrimSpace(created.Stdout)
	_, assignment := inspectLifecycleState(t, f, path)
	contract.WriteFileAbs(t, filepath.Join(path, "must-recover.txt"), "preserve me\n")
	common := strings.TrimSpace(f.Git("rev-parse", "--path-format=absolute", "--git-common-dir").Stdout)
	recoveryBlock := filepath.Join(common, "refs", "bench", "recovery")
	if err := os.MkdirAll(filepath.Dir(recoveryBlock), 0o755); err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, recoveryBlock, "blocks nested recovery refs\n")
	removeJSON := `{"session_id":"session-23","worktree_path":` + string(mustJSON(t, path)) + `}` + "\n"
	removed := contract.RunAtWithInput(t, f, f.Root, env, removeJSON, removeCommand.Command, removeCommand.Args...)
	if removed.ExitCode == 0 {
		t.Fatalf("Claude WorktreeRemove reported success after recovery-ref failure\nstdout:\n%s\nstderr:\n%s", removed.Stdout, removed.Stderr)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("failed Claude removal lost checkout: %v", err)
	}
	if got := registrationLockReason(t, f, path); !strings.Contains(got, "owner="+assignment.OwnerID) || !strings.Contains(got, "assignment="+assignment.ID) {
		t.Fatalf("failed Claude removal changed exact Bench lock: %q", got)
	}
	if err := os.Remove(filepath.Join(path, "must-recover.txt")); err != nil {
		t.Fatal(err)
	}
	ordinary := f.GitAllow("worktree", "remove", path)
	if ordinary.ExitCode == 0 {
		t.Fatal("ordinary git worktree remove bypassed the retained Bench lock")
	}
	ordinary.RequireContains(ordinary.Stderr, "locked")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("ordinary removal erased locked checkout: %v", err)
	}
	contract.WriteFileAbs(t, filepath.Join(path, "must-recover.txt"), "preserve me\n")
	if err := os.Remove(recoveryBlock); err != nil {
		t.Fatal(err)
	}
	retried := contract.RunAtWithInput(t, f, f.Root, env, removeJSON, removeCommand.Command, removeCommand.Args...)
	retried.RequireExit(0)
	wantCleanup := normalizedCleanupState{true, true, true, false}
	if got := inspectCleanupState(t, f, assignment); !reflect.DeepEqual(got, wantCleanup) {
		t.Fatalf("retry after recovery-ref restoration = %#v, want %#v", got, wantCleanup)
	}
	refs := f.Git("for-each-ref", "--format=%(refname)", "refs/bench/recovery/").Stdout
	contract.RequireContains(t, refs, "refs/bench/recovery/"+assignment.OwnerID+"/"+assignment.ID+"/")
	assignments, err := intent.Assignments(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	for _, current := range assignments {
		if current.ID == assignment.ID && current.State != intent.StateRecovered {
			t.Fatalf("retried assignment state = %q, want recovered", current.State)
		}
	}
}

func testWorktreeLifecycleSafetyFamily(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	contract.NoteContractFailure(t, "worktree lifecycle safety family contract failed")
	f := onMainFixture(t)
	env := map[string]string{"BENCH_HOME": filepath.Join(f.Root, ".bench-home")}
	var bypasses []string

	foreign := filepath.Join(t.TempDir(), "foreign")
	f.Git("worktree", "add", "-q", "-b", "foreign-safety-family", foreign, "HEAD")
	resume := f.BenchEnv(env, "resume-clean")
	if resume.ExitCode != 0 {
		bypasses = append(bypasses, "foreign classification failed")
	}
	if _, err := os.Stat(foreign); err != nil {
		bypasses = append(bypasses, "foreign worktree was removed")
	} else {
		f.Git("worktree", "remove", "--force", foreign)
	}

	created := f.BenchEnv(env, "worktree", "create", "--request", "safety-family", "--label", "safety family")
	if created.ExitCode != 0 {
		t.Fatalf("safety-family create failed before preservation probe\nstdout:\n%s\nstderr:\n%s", created.Stdout, created.Stderr)
	}
	path := worktreeCreatePath(t, created.Stdout)
	beforeLock := registrationLockReason(t, f, path)
	if beforeLock == "" {
		bypasses = append(bypasses, "creation lacked fail-closed lock")
	}
	dirty := filepath.Join(path, "must-recover.txt")
	contract.WriteFileAbs(t, dirty, "preserve before unlock\n")
	common := strings.TrimSpace(f.Git("rev-parse", "--path-format=absolute", "--git-common-dir").Stdout)
	recoveryBlock := filepath.Join(common, "refs", "bench", "recovery")
	if err := os.MkdirAll(filepath.Dir(recoveryBlock), 0o755); err != nil {
		t.Fatal(err)
	}
	contract.WriteFileAbs(t, recoveryBlock, "blocks recovery ref\n")
	release := f.BenchEnv(env, "worktree", "release", "--request", "safety-family", path)
	if release.ExitCode == 0 {
		bypasses = append(bypasses, "release succeeded without durable recovery")
	}
	if _, err := os.Stat(path); err != nil {
		bypasses = append(bypasses, "checkout disappeared before recovery")
	} else {
		afterLock := registrationLockReason(t, f, path)
		if beforeLock == "" || afterLock != beforeLock {
			bypasses = append(bypasses, "exact lock changed before recovery")
		}
		if err := os.Remove(dirty); err != nil {
			t.Fatal(err)
		}
		if ordinary := f.GitAllow("worktree", "remove", path); ordinary.ExitCode == 0 {
			bypasses = append(bypasses, "ordinary Git removal bypassed lock")
		}
	}
	if len(bypasses) != 0 {
		t.Fatalf("ownership/recovery/lock safety family bypassed: %s", strings.Join(bypasses, "; "))
	}
}

func testShiftAdapterPreflight(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	unset := f.BenchEnvSpec(contract.Env{"BENCH_AGENT": nil, "BENCH_HOME": strPtr(home)}, "shift", "probe")
	if unset.ExitCode == 0 {
		t.Fatal("shift with no BENCH_AGENT succeeded; should error")
	}
	unset.RequireContains(unset.Stderr, "BENCH_AGENT")
	if !regexp.MustCompile(`(?i)configure.*adapter|adapter.*configure`).MatchString(unset.Stderr) {
		t.Fatalf("unconfigured-adapter error is not a configure-your-adapter message:\n%s", unset.Stderr)
	}
	unset.RequireNotContains(unset.Stdout, "iteration 1/")

	empty := f.BenchEnv(map[string]string{"BENCH_AGENT": "", "BENCH_HOME": home}, "shift", "probe")
	if empty.ExitCode == 0 {
		t.Fatal("shift with empty BENCH_AGENT succeeded; should error")
	}
	empty.RequireContains(empty.Stderr, "BENCH_AGENT")
	if !regexp.MustCompile(`(?i)configure.*adapter|adapter.*configure`).MatchString(empty.Stderr) {
		t.Fatalf("empty-adapter error is not a configure-your-adapter message:\n%s", empty.Stderr)
	}
	empty.RequireNotContains(empty.Stdout, "iteration 1/")

	missing := f.BenchEnv(map[string]string{"BENCH_AGENT": "/no/such/adapter", "BENCH_HOME": home}, "shift", "probe")
	if missing.ExitCode == 0 {
		t.Fatal("shift with a missing adapter path succeeded; should error")
	}
	missing.RequireContains(missing.Stderr, "not executable")
	missing.RequireNotContains(missing.Stdout, "iteration 1/")

	keyword := f.BenchEnv(map[string]string{"BENCH_AGENT": "if", "BENCH_HOME": home}, "shift", "probe")
	if keyword.ExitCode == 0 {
		t.Fatal("shift with a shell-keyword adapter succeeded; should error")
	}
	keyword.RequireContains(keyword.Stderr, "not executable")
}

func testShiftAdapterSingleArgument(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("adapter", `#!/usr/bin/env bash
{
  printf 'argc=%s\n' "$#"
  printf 'shift_env=%s\n' "${BENCH_SHIFT:-unset}"
  printf '%s\n@@@@\n' "$1"
} >> "$BENCH_TEST_RECORD"
`)
	f.CommitAll("adapter")
	home := t.TempDir()
	record := filepath.Join(t.TempDir(), "record.txt")

	// The probe adapter only records its invocation; it never mutates the tree, so the
	// honest taxonomy is no-op/4, not complete/0.
	f.BenchEnv(map[string]string{"BENCH_TEST_RECORD": record, "BENCH_AGENT": filepath.Join(f.Root, "adapter"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "adapter-arg-probe").RequireExit(4)

	data, err := os.ReadFile(record)
	if err != nil {
		t.Fatalf("adapter was never invoked: %v", err)
	}
	text := string(data)
	for _, needle := range []string{
		"argc=1",
		"shift_env=1",
		"adapter-arg-probe",
		"You are one iteration of a Bench shift",
		"decides if it counts",
	} {
		if !strings.Contains(text, needle) {
			t.Fatalf("adapter record missing %q:\n%s", needle, text)
		}
	}
	if regexp.MustCompile(`(?m)^-p$`).MatchString(text) {
		t.Fatal("loop still passes the Claude-specific -p flag")
	}
}

func testReferenceAdapterFiles(t *testing.T) {
	root := contract.KitRoot(t)
	for _, adapter := range []string{"claude", "codex", "opencode"} {
		path := filepath.Join(root, ".bench", "adapters", adapter)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("reference adapter missing: .bench/adapters/%s", adapter)
		}
		if info.Mode()&0o111 == 0 {
			t.Fatalf("reference adapter not executable: .bench/adapters/%s", adapter)
		}
		probe := contract.NewFixture(t, contract.WithNoRepo()).Run("bash", "-n", path)
		probe.RequireExit(0)
		text, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read adapter %s: %v", adapter, err)
		}
		if !regexp.MustCompile(`(?m)^exec `).Match(text) {
			t.Fatalf("reference adapter %s does not exec its harness", adapter)
		}
		if !strings.Contains(string(text), `"$1"`) {
			t.Fatalf("reference adapter %s does not pass the prompt as $1", adapter)
		}
	}
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "claude"), `claude -p -- "$1"`, "claude adapter does not map the prompt to claude -p")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "codex"), `codex exec --sandbox workspace-write -m "$model" -- "$1"`, "routed codex adapter does not select workspace-write while preserving model and prompt")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "codex"), `codex exec --sandbox workspace-write -- "$1"`, "unrouted codex adapter does not select workspace-write while preserving the prompt")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `model="$("$_cmd" resolve-model --provider-model)"`, "opencode adapter does not request provider/model compatibility from the resolver")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `opencode run --model "$model" -- "$1"`, "routed opencode adapter does not preserve model and prompt")
	requireFileContains(t, filepath.Join(root, ".bench", "adapters", "opencode"), `opencode run -- "$1"`, "opencode adapter does not map the prompt to opencode run")
}
