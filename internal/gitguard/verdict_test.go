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
		{"branch -d other", "git branch -d old-work", refYes, "git branch -d"},
		{"branch --delete other", "git branch --delete old-work", refYes, "git branch -d"},
		{"branch mixed delete flags favors force", "git branch -d -D old-work", refYes, "git branch -D"},
		{"branch safe delete with short force favors force", "git branch -d -f old-work", refYes, "git branch -D"},
		{"branch safe delete with long force favors force", "git branch --delete --force old-work", refYes, "git branch -D"},
		{"branch reversed long force delete favors force", "git branch --force --delete old-work", refYes, "git branch -D"},
		{"branch -f move", "git branch -f main HEAD~1", refYes, "git branch -f"},
		{"branch -D delegate allowed", "git branch -D worktree-agent-x", refYes, ""},
		{"branch reversed long force delete delegate allowed", "git branch --force --delete worktree-agent-x", refYes, ""},
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

		// stash discard, filter-branch, rm -rf
		{"stash drop", "git stash drop", refYes, "git stash drop"},
		{"stash drop with a flag", "git stash drop -q", refYes, "git stash drop"},
		{"stash clear", "git stash clear", refYes, "git stash clear"},
		{"stash drop only as a flag value", "git stash -m drop", refYes, "git stash drop"},
		{"stash clear only after a double dash", "git stash -- clear", refYes, "git stash clear"},
		{"stash bare allowed", "git stash", refYes, ""},
		{"stash list allowed", "git stash list", refYes, ""},
		{"stash pop allowed", "git stash pop", refYes, ""},
		{"filter-branch bare", "git filter-branch", refYes, "git filter-branch"},
		{"filter-branch with filters", "git filter-branch --index-filter x HEAD", refYes, "git filter-branch"},
		{"rm -rf cluster", "git rm -rf build", refYes, "git rm -rf"},
		{"rm -r -f split", "git rm -r -f build", refYes, "git rm -rf"},
		{"rm --force -r", "git rm --force -r build", refYes, "git rm -rf"},
		{"rm path allowed", "git rm stale.txt", refYes, ""},
		{"rm -r allowed", "git rm -r vendor", refYes, ""},
		{"rm --cached allowed", "git rm --cached secrets.env", refYes, ""},
		{"rm --force without recursion allowed", "git rm --force stale.txt", refYes, ""},

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

// TestClassifyStashOperationsAllowed pins the working half of the verb. `drop` and
// `clear` discard work and are denied in TestClassifyVerdicts; every other stash form
// stays usable, including one whose args merely mention a safe operation name.
func TestClassifyStashOperationsAllowed(t *testing.T) {
	commands := []string{
		"git stash",
		"git stash push",
		"git stash save wip",
		"git stash pop",
		"git stash apply",
		"git stash branch recovered",
		"git stash create",
		`git stash push -m "wip thing"`,
		`git stash push -- "path/with space"`,
		"git stash -m list",
		"git stash -- list",
		"git stash --message show",
		"git stash -u -m list",
		"git -c foo=bar stash pop",
		"bash -c 'git stash pop'",
		"xargs git stash",
	}
	for _, command := range commands {
		if got := Classify(command, refYes); got != "" {
			t.Errorf("Classify(%q) = %q, want allow", command, got)
		}
	}
}

func TestClassifyStashReadOnly(t *testing.T) {
	cases := []struct {
		name string
		cmd  string
	}{
		{"list", "git stash list"},
		{"list with args", "git stash list --stat"},
		{"show", "git stash show"},
		{"show with args", "git stash show stash@{0} --stat"},
		{"list with a short flag", "git stash list -p"},
		{"show with a long flag", "git stash show --stat"},
		{"list redirected to a file", "git stash list > out.txt"},
		{"list after a global option", "git --git-dir=x stash list"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Classify(c.cmd, refYes); got != "" {
				t.Errorf("Classify(%q) = %q, want allow", c.cmd, got)
			}
		})
	}
}
