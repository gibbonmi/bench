# Destructive-git guard rework

Decisions: `decisions/git-guard-rework.md` (all tickets resolved 2026-07-01).

## Problem

The PreToolUse guard is meant to deny the agent destructive git authority, but its
current analyzer misses most of the destructive surface (bare `checkout <file>`,
`switch --discard-changes`, `stash drop/clear`, `commit --amend`, force ref moves,
one-level shell wrappers) while blocking two harmless operations
(`restore --staged`, any command merely containing the words `git push`). The
reviewer cannot trust the layer, and the agent loses turns to false positives.

## Solution

The guard becomes an honest-mistake layer with a coherent rule: block what
silently destroys uncommitted/stashed work or rewrites history; allow what git
itself refuses to do unsafely. It looks one level into shell wrappers, fails
closed on ambiguity, fixes the two provable false positives, and documents its
residual gap instead of pretending to resist evasion.

## User stories

1. As the reviewer, I want `git checkout <pathspec>` blocked in every form — bare,
   `--`, `--pathspec-from-file`, after a ref (`git checkout HEAD <file>`) — so an
   agent cannot silently discard uncommitted work.
2. As the reviewer, I want `git checkout -f` and `git switch -f` /
   `git switch --discard-changes` blocked, so force-switching cannot clobber work.
3. As the agent, I want plain `git checkout <branch>` and `git switch <branch>`
   (including `-b`/`-c` creation forms) allowed, so ordinary branch movement is
   not penalized.
4. As the reviewer, I want `git stash drop` and `git stash clear` blocked, while
   `stash` push/pop/apply/list stay allowed, so stashed work cannot be destroyed
   but normal stash flow works.
5. As the reviewer, I want `git commit --amend` blocked, so history rewrites stay
   mine even on unpushed commits.
6. As the reviewer, I want force ref surgery blocked: `branch -f`,
   `update-ref -d`, `tag -d`, `reflog expire`, `worktree remove --force`, and
   `checkout -B` / `switch -C` when the named branch already exists (they
   force-move the ref exactly like `branch -f`; with a fresh name they are
   ordinary creation and stay allowed).
7. As the reviewer, I want the existing blocks preserved: `push`, `reset --hard`,
   `clean -f`, `branch -D/-d/--delete`, `rebase`, `restore <pathspec>`.
8. As the agent, I want `git restore --staged <path>` allowed, because unstaging
   touches only the index and loses nothing.
9. As the agent, I want non-command `git` words ignored (`echo git push`,
   `grep 'git push' docs/`), so mentioning git is not running git.
10. As the reviewer, I want a `sh`/`bash`/`zsh` `-c` command string re-scanned by
    the same analyzer, so `bash -c 'git push'` is caught.
11. As the reviewer, I want `env`, `command`, `nohup`, `timeout`, and `xargs`
    treated as transparent prefixes, so `env git push` stays blocked and
    `timeout 5 git reset --hard` is caught.
12. As the reviewer, I want `xargs git <pathspec-sensitive-subcommand>` (checkout,
    restore) blocked even with no visible pathspec, because xargs exists to append
    arguments the guard cannot see — fail closed.
13. As the reviewer, I want ambiguity to fail closed: a checkout target that is
    not a resolvable ref is treated as a pathspec; outside a repo, or when
    ref-ness cannot be determined, checkout with a free argument is blocked.
14. As the reviewer, I want `git reset` without `--hard` allowed (soft/mixed lose
    nothing), so the block stays aimed at destruction, not workflow.
15. As a future maintainer, I want the hook header to state the threat model
    (honest-mistake layer) and its residual gap (no deep evasion resistance;
    backstops are the pre-push hook and worktree isolation), so nobody mistakes
    the layer for a boundary.
16. As the reviewer, I want every verdict in this spec asserted in the gate's
    contract matrix — blocked commands exit 2 with a `BLOCKED:` message, allowed
    commands exit 0 — so the guard cannot silently regress.

## Implementation decisions

- One analyzer, one verdict function: the Python block inside the hook keeps its
  shape (tokenize, find git invocations, classify), extended — no second parser,
  no external dependency.
- Command-position detection: a token counts as `git` only when it starts a
  command (first token, or immediately after `;`, `&&`, `||`, `|`, `&`, or a
  transparent prefix). This kills the word-match false positives.
- Wrapper handling is exactly one level: `-c` string re-scan for sh/bash/zsh,
  prefix-skipping for env/command/nohup/timeout/xargs. No recursion beyond that;
  the residual gap is documented, not chased.
- Checkout target classification shells out to
  `git rev-parse --verify --quiet <arg>` in the hook's working directory: a
  resolvable ref with no trailing free args is a branch switch (allowed); anything
  else is a pathspec (blocked). Unresolvable environment → blocked (fail closed).
- Subcommand option tables (`-b`/`-B`/`--orphan` for checkout, `-c`/`-C` for
  switch, `-m` etc. for commit) are extended only as far as verdicts require —
  the analyzer classifies, it does not fully parse git. The force-creation forms
  `-B`/`-C` take the branch name they target: if
  `git rev-parse --verify --quiet refs/heads/<name>` resolves, the form is a ref
  force-move and blocks; an unresolvable name is ordinary creation and passes
  (git itself fails outside a repo, so nothing is lost there).
- The contract matrix lives in `.bench/gate-runtime-git-contracts.sh` beside the
  existing cases, with an `expect_allow`/`expect_block` helper pair replacing the
  block-only `run_guard`.

## Testing decisions

- A good test here feeds the hook a JSON payload on stdin and asserts the exit
  code and, for blocks, the `BLOCKED:` stderr message — external behavior at the
  process seam, no reaching into the Python.
- Seam: the hook invocation boundary (stdin JSON → exit code + stderr). It is the
  existing, gate-tested seam; prior art is `run_guard` in
  `.bench/gate-runtime-git-contracts.sh`.
- Gate: `bench gate` (the contract file runs inside it). Done = gate green with
  the full matrix present.
- Line (declared): **mid tier — Sonnet 5 (`claude-sonnet-5`), medium effort,
  ~150k tokens** — guard logic is oracle-adjacent conformance code at a known
  seam (profile: gate/conformance → mid); no silent escalation.

### Acceptance coverage map

Red signals observed live on 2026-07-01 against the current hook (exit codes as
listed; want = post-build verdict).

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `git checkout README.md` → block | hook stdin/exit | observed exit 0, want 2 | exits 0 exactly while the destructive form is allowed |
| 1 | `git checkout HEAD README.md` → block | hook stdin/exit | observed exit 0, want 2 | same — ref-then-pathspec form slips the `--`-only detection |
| 2 | `git checkout -f main`, `git switch -f main`, `git switch --discard-changes main` → block | hook stdin/exit | observed exit 0 (all three), want 2 | force-switch destruction is currently invisible to the analyzer |
| 3 | `git checkout main`, `git switch main` → allow | hook stdin/exit | already green (exit 0) — regression rows | would flip to 2 if pathspec blocking overreaches into branch switching |
| 4 | `git stash drop`, `git stash clear` → block; `git stash pop` → allow | hook stdin/exit | observed exit 0 on drop/clear, want 2; pop already green | stash destruction has no case in the dispatch today |
| 5 | `git commit --amend -m x` → block | hook stdin/exit | observed exit 0, want 2 | amend is a history rewrite the guard claims to own |
| 6 | `branch -f`, `update-ref -d`, `tag -d`, `reflog expire`, `worktree remove --force` → block | hook stdin/exit | observed exit 0 (all five), want 2 | none of these subcommand/flag pairs exist in the dispatch |
| 6 | `checkout -B <existing>` / `switch -C <existing>` → block; fresh-name forms → allow | hook stdin/exit | observed exit 0 on `-B main`/`-C main`, want 2; fresh-name forms already green as regression rows | force-creation onto an existing branch is a ref force-move the dispatch cannot see today |
| 7 | existing blocks (`push`, `reset --hard`, `clean -fd`, `branch -D`, `rebase`, `restore <path>`) → block | hook stdin/exit | already covered (existing contract cases) | regression guard for the rework |
| 8 | `git restore --staged README.md` → allow | hook stdin/exit | observed exit 2, want 0 | the false positive fires today; row flips when fixed |
| 9 | `echo git push` → allow | hook stdin/exit | observed exit 2, want 0 | word-match false positive fires today |
| 10 | `bash -c "git push"`, `sh -c "git reset --hard"` → block | hook stdin/exit | observed exit 0 (both), want 2 | one-level wrapper strings are unscanned today |
| 11 | `env git push` → block (stays), `timeout 5 git reset --hard` → block | hook stdin/exit | env form already covered; timeout form not TDD-able red first (current basename scan may pass it) — assert in matrix | prefix transparency must survive the command-position rework, which would otherwise stop matching `git` after a prefix |
| 12 | `xargs git checkout` → block | hook stdin/exit | observed exit 0, want 2 | xargs-fed pathspecs are invisible; fail-closed row |
| 13 | checkout free-arg that is not a ref, outside-repo checkout free-arg → block | hook stdin/exit | not TDD-able as a distinct red (superset of story 1 rows until ref resolution exists); asserted in matrix with a non-ref token | proves classification fails closed rather than defaulting to allow |
| 14 | `git reset --soft HEAD~1` → allow | hook stdin/exit | already green — regression row | would flip if reset blocking overreaches past `--hard` |
| 15 | header states threat model + residual gap | file content | not TDD-able (prose); reviewer check at review phase | — |
| 16 | full matrix present both ways | gate | matrix rows are the signal | any silent regression turns a row red in `bench gate` |

## Out of scope

- **Evasion-resistant command analysis** (deep recursion, alias resolution,
  interpreter tracing) — a different capability declined by decision #1; the
  backstops own that threat. Estimate if ever built: 1–2 days agent time, and it
  still cannot close the channel.
- **Codex hook layer repair** (`.codex/hooks.json` likely inert) — separate
  adapter capability with its own decision map to come (`/bench-shape-idea`
  queued). Estimate: research spike ~1h + adapter build unknown until the spike.
