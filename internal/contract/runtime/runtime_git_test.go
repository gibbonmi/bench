package runtime

import (
	"encoding/json"
	"errors"
	"github.com/gibbonmi/bench/internal/contract"
	"os"
	"path/filepath"
	"testing"
)

func TestRuntimeGitContracts(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	contract.RunParallel(t, "block-dangerous-git allow/block matrix contract failed", testBlockDangerousGitMatrix)
}

func testBlockDangerousGitMatrix(t *testing.T) {
	f := runtimeGitFixture(t)

	for _, c := range runtimeGitMatrix() {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			probe := runGitGuard(t, f, c)
			if c.block {
				probe.RequireExit(2)
				probe.RequireContains(probe.Stdout+probe.Stderr, "BLOCKED:")
				return
			}
			probe.RequireExit(0)
		})
	}
}

type runtimeGitCase struct {
	name        string
	command     string
	block       bool
	outsideRepo bool
}

func runtimeGitMatrix() []runtimeGitCase {
	return []runtimeGitCase{
		{name: "blocks git push", command: "git push", block: true},
		{name: "blocks git -C push", command: "git -C . push", block: true},
		{name: "blocks git -C reset hard", command: "git -C /tmp reset --hard", block: true},
		{name: "blocks forced clean", command: "git -C . clean -fd", block: true},
		{name: "blocks branch delete", command: "git branch -D old-work", block: true},
		{name: "blocks rebase", command: "git rebase main", block: true},
		{name: "blocks checkout explicit pathspec", command: "git checkout -- README.md", block: true},
		{name: "blocks checkout pathspec-from-file equals", command: "git checkout --pathspec-from-file=paths.txt", block: true},
		{name: "blocks checkout pathspec-from-file separate", command: "git checkout --pathspec-from-file paths.txt", block: true},
		{name: "blocks restore path", command: "git restore README.md", block: true},
		{name: "blocks restore pathspec-from-file equals", command: "git restore --pathspec-from-file=paths.txt", block: true},
		{name: "blocks restore pathspec-from-file separate", command: "git restore --pathspec-from-file paths.txt", block: true},
		{name: "blocks checkout file", command: "git checkout README.md", block: true},
		{name: "blocks checkout file with redirection", command: "git checkout README.md > checkout.log", block: true},
		{name: "blocks checkout commit plus file", command: "git checkout HEAD README.md", block: true},
		{name: "blocks forced checkout", command: "git checkout -f main", block: true},
		{name: "blocks forced switch", command: "git switch -f main", block: true},
		{name: "blocks switch discard changes", command: "git switch --discard-changes main", block: true},
		{name: "blocks stash drop", command: "git stash drop", block: true},
		{name: "blocks stash clear", command: "git stash clear", block: true},
		{name: "blocks commit amend", command: "git commit --amend -m x", block: true},
		{name: "blocks branch force move", command: "git branch -f main HEAD~1", block: true},
		{name: "blocks update-ref delete", command: "git update-ref -d refs/heads/x", block: true},
		{name: "blocks tag delete short", command: "git tag -d v1", block: true},
		{name: "blocks tag delete long", command: "git tag --delete v1", block: true},
		{name: "blocks reflog expire", command: "git reflog expire --expire=now --all", block: true},
		{name: "blocks forced worktree remove", command: "git worktree remove --force wt", block: true},
		{name: "blocks checkout existing branch recreate", command: "git checkout -B main", block: true},
		{name: "blocks switch existing branch recreate", command: "git switch -C main", block: true},
		{name: "allows delegate relative worktree removal", command: "git worktree remove --force .claude/worktrees/agent-abc123"},
		{name: "allows delegate absolute worktree removal", command: "git worktree remove --force /home/u/repo/.claude/worktrees/agent-abc123"},
		{name: "allows delegate branch force delete", command: "git branch -D worktree-agent-abc123"},
		{name: "allows delegate branch delete multiple", command: "git branch -d worktree-agent-abc123 worktree-agent-def456"},
		{name: "blocks delegate worktree traversal", command: "git worktree remove --force .claude/worktrees/../../src", block: true},
		{name: "blocks mixed delegate worktree targets", command: "git worktree remove --force .claude/worktrees/agent-a wt2", block: true},
		{name: "blocks bare forced worktree remove", command: "git worktree remove --force", block: true},
		{name: "blocks mixed delegate branch delete", command: "git branch -D worktree-agent-a main", block: true},
		{name: "blocks bare branch delete", command: "git branch -D", block: true},
		{name: "blocks delegate branch force move", command: "git branch -f worktree-agent-a HEAD~1", block: true},
		{name: "blocks newline-separated push", command: "git add -A\ngit commit -m wip\ngit push origin main", block: true},
		{name: "blocks newline-separated reset hard", command: "git status\ngit reset --hard HEAD~1", block: true},
		{name: "blocks wrapper string with newline push", command: "bash -c 'cd repo\ngit push'", block: true},
		{name: "allows newline-separated harmless flow", command: "git add -A\ngit status --short\ngit commit -m wip"},
		{name: "blocks later push after and separator", command: "git status && git push", block: true},
		{name: "blocks later push after semicolon", command: "git add -A; git push origin main", block: true},
		{name: "blocks bash wrapper push", command: "bash -c 'git push'", block: true},
		{name: "blocks bash login wrapper push", command: "bash -lc 'git push'", block: true},
		{name: "blocks sh wrapper reset hard", command: "sh -c 'git reset --hard'", block: true},
		{name: "blocks zsh wrapper stash clear", command: "zsh -c 'git stash clear'", block: true},
		{name: "blocks env prefix push", command: "env git push", block: true},
		{name: "blocks env assignment push", command: "GIT_TRACE=1 git push", block: true},
		{name: "blocks timeout prefix reset hard", command: "timeout 5 git reset --hard", block: true},
		{name: "blocks xargs checkout", command: "xargs git checkout", block: true},
		{name: "blocks xargs restore", command: "xargs git restore", block: true},
		{name: "blocks command prefix push", command: "command git push", block: true},
		{name: "blocks nohup prefix push", command: "nohup git push", block: true},
		{name: "blocks non-ref checkout target", command: "git checkout not-a-ref-anywhere", block: true},
		{name: "allows checkout branch", command: "git checkout main"},
		{name: "allows switch branch", command: "git switch main"},
		{name: "allows checkout new branch", command: "git checkout -b feature-fresh"},
		{name: "allows switch new branch", command: "git switch -c feature-fresh2"},
		{name: "allows checkout fresh branch recreate", command: "git checkout -B feature-fresh3"},
		{name: "allows switch fresh branch recreate", command: "git switch -C feature-fresh4"},
		{name: "allows checkout branch with stderr redirection", command: "git checkout main 2>/dev/null"},
		{name: "allows checkout branch with stdout redirection", command: "git checkout main > checkout.log"},
		{name: "blocks bare stash", command: "git stash", block: true},
		{name: "blocks stash pop", command: "git stash pop", block: true},
		{name: "blocks stash apply", command: "git stash apply", block: true},
		{name: "blocks stash push", command: "git stash push -m wip", block: true},
		{name: "blocks stash message spelled like a read-only verb", command: "git stash -m list", block: true},
		{name: "blocks stash pathspec spelled like a read-only verb", command: "git stash -- list", block: true},
		{name: "allows stash list", command: "git stash list"},
		{name: "allows stash show", command: "git stash show"},
		{name: "allows index-only restore", command: "git restore --staged README.md"},
		{name: "allows soft reset", command: "git reset --soft HEAD~1"},
		{name: "allows mixed reset", command: "git reset HEAD~1"},
		{name: "allows status with git -C", command: "git -C . status --short"},
		{name: "allows commit containing checkout word", command: "git commit -m checkout-notes"},
		{name: "allows non-command echo git push", command: "echo git push"},
		{name: "blocks checkout outside repo", command: "git checkout main", block: true, outsideRepo: true},
	}
}

func runtimeGitFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := contract.NewFixture(t)
	installFixtureBenchWrapper(t, f)
	f.Git("-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", "init")
	f.GitAllow("branch", "-q", "main")
	f.WriteFile("README.md", "readme\n")
	return f
}

func installFixtureBenchWrapper(t *testing.T, f contract.Fixture) {
	t.Helper()
	wrapper := filepath.Join(f.Root, ".bench", "bin", "bench.sh")
	if err := os.MkdirAll(filepath.Dir(wrapper), 0o755); err != nil {
		t.Fatalf("create fixture bench wrapper dir: %v", err)
	}
	target := filepath.Join(contract.SubjectRoot(t), "bin", "bench.sh")
	if err := os.Symlink(target, wrapper); err != nil {
		if !errors.Is(err, os.ErrExist) {
			t.Fatalf("symlink fixture bench wrapper: %v", err)
		}
	}
	bin := filepath.Join(f.Root, ".bench-contract-env", "path-bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatalf("create fixture PATH bin: %v", err)
	}
	global := filepath.Join(bin, "bench")
	if err := os.Symlink(target, global); err != nil {
		if !errors.Is(err, os.ErrExist) {
			t.Fatalf("symlink fixture global bench wrapper: %v", err)
		}
	}
	f.Env["PATH"] = bin + string(os.PathListSeparator) + os.Getenv("PATH")
}

func runGitGuard(t *testing.T, f contract.Fixture, c runtimeGitCase) contract.Probe {
	t.Helper()
	envelope, err := json.Marshal(struct {
		ToolInput struct {
			Command string `json:"command"`
		} `json:"tool_input"`
	}{ToolInput: struct {
		Command string `json:"command"`
	}{Command: c.command}})
	if err != nil {
		t.Fatalf("marshal guard envelope: %v", err)
	}

	dir := f.Root
	if c.outsideRepo {
		dir = t.TempDir()
	}
	return contract.RunAtWithInput(t, f, dir, nil, string(append(envelope, '\n')), "bash", filepath.Join(contract.SubjectRoot(t), ".bench", "hooks", "block-dangerous-git.sh"))
}
