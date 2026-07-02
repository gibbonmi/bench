# Runtime contracts for the destructive-git guard hook: the full allow/block
# matrix from specs/git-guard-rework.md. Every verdict is asserted both ways —
# blocked commands exit 2 with a BLOCKED: message, allowed commands exit 0.

tmp="$(mktemp -d)"
(
  set -u; cd "$tmp"
  guard="$root/.bench/hooks/block-dangerous-git.sh"
  fails=0
  verdict_out=""
  run_verdict() {
    local command="$1" rc
    verdict_out="$(printf '{"tool_input":{"command":"%s"}}\n' "$command" | bash "$guard" 2>&1)" && rc=0 || rc=$?
    return "$rc"
  }
  expect_block() {
    local command="$1" rc
    run_verdict "$command" && rc=0 || rc=$?
    [ "$rc" = "2" ] || { echo "guard allowed: $command (exit $rc)"; fails=$((fails+1)); return; }
    grep -qF 'BLOCKED:' <<<"$verdict_out" || { echo "guard block output not actionable for: $command"; fails=$((fails+1)); }
  }
  expect_allow() {
    local command="$1" rc
    run_verdict "$command" && rc=0 || rc=$?
    [ "$rc" = "0" ] || { echo "guard blocked harmless: $command (exit $rc)"; fails=$((fails+1)); }
  }

  # The repo the guard classifies refs against: 'main' resolves, README.md exists.
  git init -q repo
  cd repo
  git -c user.email=bench@local -c user.name=bench commit -q --allow-empty -m init
  git branch -q main 2>/dev/null || true
  printf 'readme\n' > README.md

  # -- push, hard reset, forced clean, branch delete, rebase, pathspec forms block (story 7)
  expect_block 'git push'
  expect_block 'git -C . push'
  expect_block 'git -C /tmp reset --hard'
  expect_block 'git -C . clean -fd'
  expect_block 'git branch -D old-work'
  expect_block 'git rebase main'
  expect_block 'git checkout -- README.md'
  expect_block 'git checkout --pathspec-from-file=paths.txt'
  expect_block 'git checkout --pathspec-from-file paths.txt'
  expect_block 'git restore README.md'
  expect_block 'git restore --pathspec-from-file=paths.txt'
  expect_block 'git restore --pathspec-from-file paths.txt'

  # -- checkout with a pathspec blocks in every form (story 1); forced switch blocks (story 2)
  expect_block 'git checkout README.md'
  expect_block 'git checkout README.md > checkout.log'
  expect_block 'git checkout HEAD README.md'
  expect_block 'git checkout -f main'
  expect_block 'git switch -f main'
  expect_block 'git switch --discard-changes main'

  # -- stash destruction (story 4), amend (story 5), ref surgery (story 6)
  expect_block 'git stash drop'
  expect_block 'git stash clear'
  expect_block 'git commit --amend -m x'
  expect_block 'git branch -f main HEAD~1'
  expect_block 'git update-ref -d refs/heads/x'
  expect_block 'git tag -d v1'
  expect_block 'git tag --delete v1'
  expect_block 'git reflog expire --expire=now --all'
  expect_block 'git worktree remove --force wt'
  expect_block 'git checkout -B main'
  expect_block 'git switch -C main'

  # -- delegate-scratch carve-out: harness worktrees under .claude/worktrees/ and
  #    their worktree-* branches are agent-created scratch, cleanable by the agent.
  #    Everything adjacent — traversal out of the directory, mixed targets, bare
  #    force-remove, branch force-move — still blocks.
  expect_allow 'git worktree remove --force .claude/worktrees/agent-abc123'
  expect_allow 'git worktree remove --force /home/u/repo/.claude/worktrees/agent-abc123'
  expect_allow 'git branch -D worktree-agent-abc123'
  expect_allow 'git branch -d worktree-agent-abc123 worktree-agent-def456'
  expect_block 'git worktree remove --force .claude/worktrees/../../src'
  expect_block 'git worktree remove --force .claude/worktrees/agent-a wt2'
  expect_block 'git worktree remove --force'
  expect_block 'git branch -D worktree-agent-a main'
  expect_block 'git branch -D'
  expect_block 'git branch -f worktree-agent-a HEAD~1'

  # -- wrapper strings and prefixed invocations are classified like bare ones (stories 10, 11, 12)
  expect_block "bash -c 'git push'"
  expect_block "bash -lc 'git push'"
  expect_block "sh -c 'git reset --hard'"
  expect_block "zsh -c 'git stash clear'"
  expect_block 'env git push'
  expect_block 'GIT_TRACE=1 git push'
  expect_block 'timeout 5 git reset --hard'
  expect_block 'xargs git checkout'
  expect_block 'xargs git restore'
  expect_block 'command git push'
  expect_block 'nohup git push'

  # -- non-ref checkout target fails closed (story 13)
  expect_block 'git checkout not-a-ref-anywhere'

  # -- branch movement and creation are allowed, redirections included (story 3)
  expect_allow 'git checkout main'
  expect_allow 'git switch main'
  expect_allow 'git checkout -b feature-fresh'
  expect_allow 'git switch -c feature-fresh2'
  expect_allow 'git checkout -B feature-fresh3'
  expect_allow 'git switch -C feature-fresh4'
  expect_allow 'git checkout main 2>/dev/null'
  expect_allow 'git checkout main > checkout.log'

  # -- stash flow (story 4), index-only restore (story 8), soft/mixed resets (story 14) are allowed
  expect_allow 'git stash'
  expect_allow 'git stash pop'
  expect_allow 'git stash apply'
  expect_allow 'git stash push -m wip'
  expect_allow 'git stash list'
  expect_allow 'git restore --staged README.md'
  expect_allow 'git reset --soft HEAD~1'
  expect_allow 'git reset HEAD~1'

  # -- non-command git words and everyday commands (story 9)
  expect_allow 'git -C . status --short'
  expect_allow 'git commit -m checkout-notes'
  expect_allow 'echo git push'

  # -- outside a repo, ref-ness is unknowable: checkout free-arg fails closed (story 13)
  cd "$tmp"
  expect_block 'git checkout main'

  [ "$fails" -eq 0 ] || { echo "$fails guard matrix case(s) failed"; exit 1; }
) || err "block-dangerous-git allow/block matrix contract failed"
rm -rf "$tmp"
