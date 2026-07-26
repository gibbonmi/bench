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

		// worktree (stash has its own tests below)
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

// TestClassifyStashMutations walks both stash deny classes. Expectations read their labels
// from the deny table rather than repeating the strings, so rewording a row cannot leave a
// stale copy passing here. The two labels must differ: the blocked agent reads the label to
// learn which hazard it hit.
func TestClassifyStashMutations(t *testing.T) {
	worktree := denyLabels["stash-push"]
	history := denyLabels["stash"]
	if worktree == "" || history == "" {
		t.Fatalf("stash deny classes missing from the table: stash-push=%q stash=%q", worktree, history)
	}
	if worktree == history {
		t.Fatalf("both stash classes carry label %q, so the refusal names the wrong hazard", worktree)
	}
	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"bare stash is push", "git stash", worktree},
		{"push", "git stash push", worktree},
		{"save", "git stash save wip", worktree},
		{"pop", "git stash pop", worktree},
		{"apply", "git stash apply", worktree},
		{"branch", "git stash branch recovered", worktree},
		{"unrecognized verb fails closed", "git stash create", worktree},

		// drop and clear destroy stash history, a different hazard from cross-applying
		// working-tree state, and keep their own label.
		{"drop", "git stash drop", history},
		{"clear", "git stash clear", history},

		{"quoted multi-word message", `git stash push -m "wip thing"`, worktree},
		{"pathspec past the separator", `git stash push -- "path/with space"`, worktree},
		{"one wrapper level", "bash -c 'git stash pop'", worktree},
		{"xargs", "xargs git stash", worktree},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, refYes); got != c.want {
				t.Errorf("Classify(%q) = %q, want %q", c.cmd, got, c.want)
			}
		})
	}
}

// TestClassifyStashReadOnly pins the allow set. Without it, refusing every `git stash`
// would satisfy TestClassifyStashMutations while taking away the inspection an agent
// legitimately needs.
func TestClassifyStashReadOnly(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
	}{
		{"list", "git stash list"},
		{"list with args", "git stash list --stat"},
		{"show", "git stash show"},
		{"show with args", "git stash show stash@{0} --stat"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, refYes); got != "" {
				t.Errorf("Classify(%q) = %q, want allow", c.cmd, got)
			}
		})
	}
}
