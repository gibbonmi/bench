package gitguard

import "testing"

// Test checkers stand in for repository truth so classification is exercised without a
// real repo. refYes: every ref resolves and no branch exists (the "target is a valid
// ref / creation is fresh" world). refNo: nothing resolves and every branch exists (the
// "unknown target / forced creation would clobber" world).
var (
	refYes = Checker{RefResolves: func(string) bool { return true }, BranchExists: func(string) bool { return false }}
	refNo  = Checker{RefResolves: func(string) bool { return false }, BranchExists: func(string) bool { return true }}
)

// TestClassifyVerdicts walks each verb's option and carve-out matrix through the full
// tokenize→scan→classify path, pinning what the 79-case shell matrix samples.
func TestClassifyVerdicts(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
		chk  Checker
		want string
	}{
		// unconditional destructive verbs
		{"push", "git push", refYes, "git push"},
		{"rebase", "git rebase main", refYes, "history rewrite"},
		{"reset --hard", "git reset --hard", refYes, "git reset --hard"},
		{"reset --soft allowed", "git reset --soft HEAD~1", refYes, ""},
		{"reset mixed allowed", "git reset HEAD~1", refYes, ""},
		{"clean -fd", "git clean -fd", refYes, "git clean -f"},
		{"clean --force", "git clean --force", refYes, "git clean -f"},
		{"amend", "git commit --amend -m x", refYes, "git commit --amend"},
		{"commit allowed", "git commit -m checkout-notes", refYes, ""},
		{"update-ref -d", "git update-ref -d refs/heads/x", refYes, "git update-ref -d"},
		{"tag -d", "git tag -d v1", refYes, "git tag -d"},
		{"tag --delete", "git tag --delete v1", refYes, "git tag -d"},
		{"reflog expire", "git reflog expire --expire=now --all", refYes, "git reflog expire"},

		// branch delete / force + worktree- carve-out
		{"branch -D other", "git branch -D old-work", refYes, "git branch -D"},
		{"branch -f move", "git branch -f main HEAD~1", refYes, "git branch -f"},
		{"branch -D delegate allowed", "git branch -D worktree-agent-x", refYes, ""},
		{"branch -d two delegates allowed", "git branch -d worktree-a worktree-b", refYes, ""},
		{"branch -D mixed blocks", "git branch -D worktree-a main", refYes, "git branch -D"},
		{"branch -D bare blocks", "git branch -D", refYes, "git branch -D"},
		{"branch -f delegate still blocks (force-move)", "git branch -f worktree-a HEAD~1", refYes, "git branch -f"},
		{"branch list allowed", "git branch", refYes, ""},

		// checkout
		{"checkout pathspec --", "git checkout -- README.md", refYes, "git checkout path"},
		{"checkout file blocks", "git checkout README.md", refNo, "git checkout path"},
		{"checkout HEAD file blocks (two free args)", "git checkout HEAD README.md", refYes, "git checkout path"},
		{"checkout -f", "git checkout -f main", refYes, "git checkout path"},
		{"checkout pathspec-from-file", "git checkout --pathspec-from-file=paths.txt", refYes, "git checkout path"},
		{"checkout resolving ref allowed", "git checkout main", refYes, ""},
		{"checkout unknown ref blocks", "git checkout not-a-ref", refNo, "git checkout path"},
		{"checkout -b allowed", "git checkout -b feature", refYes, ""},
		{"checkout -B existing blocks", "git checkout -B main", refNo, "git checkout path"}, // refNo → BranchExists true
		{"checkout -B fresh allowed", "git checkout -B fresh", refYes, ""},                  // refYes → BranchExists false
		{"xargs checkout blocks", "xargs git checkout", refYes, "git checkout path"},

		// switch
		{"switch -f", "git switch -f main", refYes, "git switch --force"},
		{"switch --discard-changes", "git switch --discard-changes main", refYes, "git switch --force"},
		{"switch -C existing blocks", "git switch -C main", refNo, "git switch --force"}, // refNo → BranchExists true
		{"switch allowed", "git switch main", refYes, ""},
		{"switch -c allowed", "git switch -c feature", refYes, ""},

		// restore
		{"restore path", "git restore README.md", refYes, "git restore path"},
		{"restore dot", "git restore .", refYes, "git restore path"},
		{"restore --staged allowed", "git restore --staged README.md", refYes, ""},
		{"restore pathspec-from-file", "git restore --pathspec-from-file=p", refYes, "git restore path"},
		{"xargs restore blocks", "xargs git restore", refYes, "git restore path"},

		// stash / worktree
		{"stash drop", "git stash drop", refYes, "git stash drop"},
		{"stash clear", "git stash clear", refYes, "git stash drop"},
		{"stash pop allowed", "git stash pop", refYes, ""},
		{"stash push allowed", "git stash push -m wip", refYes, ""},
		{"worktree remove --force other", "git worktree remove --force wt", refYes, "git worktree remove --force"},
		{"worktree remove delegate allowed", "git worktree remove --force .claude/worktrees/agent-x", refYes, ""},
		{"worktree remove traversal blocks", "git worktree remove --force .claude/worktrees/../src", refYes, "git worktree remove --force"},
		{"worktree remove mixed blocks", "git worktree remove --force .claude/worktrees/a wt2", refYes, "git worktree remove --force"},
		{"worktree remove bare blocks", "git worktree remove --force", refYes, "git worktree remove --force"},

		// non-command git words, prefixes, wrappers
		{"echo git push allowed", "echo git push", refYes, ""},
		{"env prefix", "env git push", refYes, "git push"},
		{"timeout prefix", "timeout 5 git reset --hard", refYes, "git reset --hard"},
		{"command prefix", "command git push", refYes, "git push"},
		{"nohup prefix", "nohup git push", refYes, "git push"},
		{"wrapper -lc", "bash -lc 'git push'", refYes, "git push"},
		{"later separator blocks", "git status && git push", refYes, "git push"},
		{"newline block blocks", "git add -A\ngit commit -m wip\ngit push origin main", refYes, "git push"},
		{"clean newline flow allowed", "git add -A\ngit status --short\ngit commit -m wip", refYes, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, c.chk); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}
