# Decision map — destructive-git guard rework

The PreToolUse guard (`.bench/hooks/block-dangerous-git.sh`) has verified holes
(bare `git checkout <file>`, `git switch --discard-changes`, `stash drop/clear`,
`commit --amend`, wrapper one-liners) and false positives (`git restore --staged`,
non-command `git` tokens). This map records the design decisions for the rework.

## #1: What is the guard's threat model?

Type: Grill

### Question
Is the guard an honest-mistake layer or an evasion-resistant boundary? The answer
decides parsing depth, the bypass-handling scope, and what the hook header may
honestly claim.

### Answer
Honest-mistake layer. It stops a well-meaning agent from reflexively running
destructive git; it does not resist deliberate evasion. The pre-push git hook and
pooled-worktree isolation are the backstops for a misaligned agent. The residual
gap is documented in the hook header, not pretended away.

## #2: How are ambiguous commands and known false positives handled?

Blocked by: #1
Type: Grill

### Question
When the analyzer cannot classify a command, does it fail open or closed, and
which current false positives are fixed?

### Answer
Fail closed on genuine ambiguity — a wrongly blocked agent loses a turn; a
wrongly allowed checkout loses work. Two provable false positives are fixed
precisely: `git restore --staged <path>` is allowed (index-only, loses nothing),
and only tokens in command position count as `git` invocations (no more blocking
`echo git push`).

## #3: What is the blocklist?

Blocked by: #1
Type: Grill

### Question
Which subcommands join the blocked set, and is `commit --amend` the agent's or
the reviewer's?

### Answer
Rule: block what silently destroys uncommitted/stashed work or rewrites history;
allow anything git itself refuses to do unsafely.

Blocked (added): `checkout` with any pathspec (bare or `--`) or `-f`,
`switch --discard-changes`/`-f`, `stash drop`/`stash clear`, `branch -f`,
`commit --amend`, `reflog expire`, `update-ref -d`, `tag -d`,
`worktree remove --force`. Existing blocks (`push`, `reset --hard`, `clean -f`,
`branch -D/-d/--delete`, `rebase`, `restore <pathspec>`) stay.

Allowed: plain `checkout <branch>` / `switch <branch>`, `stash`
push/pop/apply/list, `restore --staged <path>`, `reset` without `--hard`.

`commit --amend` is blocked: in a shift the loop commits on green, never the
agent, and all history rewrites are the reviewer's.

## #4: How far does wrapper handling go?

Blocked by: #1
Type: Grill

### Question
Does the guard look inside `bash -c` strings, `xargs`, `env` prefixes — and how
deep before it becomes the evasion game #1 declined?

### Answer
One level, string-scan only. A `sh`/`bash`/`zsh` token with `-c` gets its command
string re-run through the same analyzer; `env`, `command`, `nohup`, `timeout`,
and `xargs` are transparent prefixes (skipped, scanning continues for `git`). No
deeper recursion. This catches the wrappers honest agents actually produce; the
rest is the documented residual gap.

## #5: How is the rework verified?

Blocked by: #2, #3, #4
Type: Research

### Answer
Extend `.bench/gate-runtime-git-contracts.sh` with a full allow/block matrix: one
case per blocklist row, one per allowed row (regressions for the fixed false
positives), and the wrapper cases from #4. Every verdict asserted both ways —
blocked commands exit 2 with a `BLOCKED:` message, allowed commands exit 0.

---

All tickets resolved. Ready for /bench-write-spec.
