# `bench worktree clean` + honest status actions

Status: staged

## Problem
`bench status` tells a reviewer with out-of-pool worktrees to "resume or clean
up (bench worktree)", but `bench worktree clean` is not a verb today:
`bench worktree` ignores every extra argument and opens the pooled subshell. The
dashboard points at an action that cannot clear the signal, and an unknown
worktree argument exits successfully while doing the wrong thing.

## Solution
Add a human-attended `bench worktree clean` verb that removes only clean
out-of-pool worktrees with non-force `git worktree remove`, then runs
`git worktree prune`. Reject unknown worktree arguments with usage exit 2, keep
bare `bench worktree` as the pooled subshell, and split the status worktree
signal so each class names an action that can actually affect it.

## User stories
1. As a reviewer, I want `bench worktree clean` to list clean out-of-pool
   worktrees, require typed confirmation on a TTY, remove confirmed candidates
   without `--force`, and prune missing registrations, so that status cleanup is
   deliberate and recoverable. Line: gpt-5.5 / medium. The reviewer requested
   GPT-5.5 for FT11, and the work is gate-observable but spans git worktree
   state and human-attended command posture.
2. As a reviewer, I want `bench worktree clean` to refuse non-TTY stdin before
   removing anything, so that headless agents cannot sweep worktrees. Line:
   gpt-5.5 / medium. The reviewer requested GPT-5.5 for FT11, and the behavior
   is safety-critical even though the seam is known.
3. As a reviewer, I want dirty out-of-pool worktrees, leased pool entries, the
   repo root, and warm unleased pool entries to remain untouched by `clean`, so
   that live or user-owned work is not destroyed. Line: gpt-5.5 / medium. The
   reviewer requested GPT-5.5 for FT11, and the gate can observe the negative
   cases but the candidate classification must not duplicate status logic.
4. As a reviewer, I want `bench worktree <unknown>` to exit 2 with usage naming
   valid forms, while bare `bench worktree` still opens the pooled subshell, so
   that typos never silently run a different command. Line: gpt-5.5 / medium.
   The reviewer requested GPT-5.5 for FT11, and this is a small dispatch change
   whose wrong behavior is directly testable.
5. As a reviewer, I want `bench status` to report out-of-pool worktrees,
   leased pool entries, and orphaned scratch branches as separate rows when
   present, with action text matched to each class, so that the ambient
   dashboard remains honest. Line: gpt-5.5 / medium. The reviewer requested
   GPT-5.5 for FT11, and the existing status seam can observe the full output
   but the semantics are user-facing.
6. As a maintainer, I want worktree classification and TTY detection to have one
   shared source each, so that cleanup and status do not drift into parallel
   implementations of the same facts. Line: gpt-5.5 / medium. The reviewer
   requested GPT-5.5 for FT11, and this is a standards-sensitive kit change
   where duplicated knowledge is the main failure mode.

## Implementation decisions
- Worktree command dispatch changes only at the command router: no args keep the
  existing subshell, `clean` routes to a new cleanup command, and any other arg
  returns usage exit 2.
- The cleanup command owns the full flow: resolve the repo root, enumerate
  registered worktrees, classify them, filter clean out-of-pool candidates,
  print the candidate list, require a typed confirmation token on a TTY, remove
  candidates with non-force `git worktree remove`, run `git worktree prune`, and
  report removed and refused paths.
- Removal never passes `--force`. Dirty candidates, candidates dirtied between
  listing and removal, and the current worktree are refused by git and reported,
  not destroyed.
- Worktree classification lives once in `internal/worktree` and tags registered
  worktrees as root, pool warm, pool leased, or out-of-pool. Status and cleanup
  both consume that classifier.
- TTY detection moves out of the gate-pin implementation into one shared helper
  with injectable terminal detection, so both human-attended commands use the
  same `os.ModeCharDevice` fact.
- `bench status` emits per-class worktree rows instead of one combined active
  worktree count. Out-of-pool rows name `bench worktree clean`; leased pool rows
  do not name `clean`; orphaned scratch branch rows keep the existing delete
  branch action.
- The literal interactive prompt is not covered by a PTY helper. As with
  `bench gate pin`, tests cover non-TTY refusal and forced-terminal confirm paths.

## Testing decisions
- Good tests exercise external behavior through the runtime CLI and status
  seams, using throwaway fixture repos and real git worktree state.
- The cleanup implementation may have focused package tests for candidate
  classification and TTY helper reuse, but acceptance coverage attaches at the
  runtime seam because the user-visible contract is command behavior.
- The feature must pass `.bench/gate.sh`.

### Seam diagram
Runtime worktree cleanup seam:

    trigger: reviewer runs bench worktree clean
        |
        v
    CLI args + stdin + git worktree state  -->  [ bench worktree clean ]  -->  stdout, stderr, exit code, git worktree registry
    fixture repo + forced terminal input   -->  [                      ]
                                             <- tests attach here: runtime contract drives the CLI and inspects git state

Runtime status seam:

    trigger: reviewer runs bench status
        |
        v
    git repo + worktree registry + lease files  -->  [ bench status ]  -->  ambient dashboard rows
    fixture worktrees and leases                -->  [              ]
                                                   <- tests attach here: runtime contract inspects rendered status text

Shared classification seam:

    trigger: cleanup command or status renderer asks for registered worktrees
        |
        v
    repo root + BENCH_HOME + git worktree list  -->  [ worktree classifier ]  -->  root, pool warm, pool leased, out-of-pool tags
    fixture repo and pool layout                -->  [                     ]
                                                 <- tests attach here: package tests assert tags for representative layouts

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Clean out-of-pool worktree is listed, confirmed through forced-terminal input, removed with non-force git worktree remove, and followed by prune. | Runtime worktree cleanup seam | Add `testRuntimeWorktreeCleanRemovesConfirmedOutOfPool`; it fails on current code because `cmd/bench/main.go` sends `clean` to `worktree.Subshell` instead of a cleanup verb. | A subshell-only implementation never removes the registered out-of-pool worktree, so the fixture still appears in `git worktree list --porcelain`. |
| 1 | No candidates prints a nothing-to-clean result and exits 0. | Runtime worktree cleanup seam | Add `testRuntimeWorktreeCleanNoCandidates`; it fails on current code because `bench worktree clean` enters the subshell path instead of reporting cleanup state. | A command that cannot distinguish no candidates from subshell launch leaves no cleanup-specific output for status users. |
| 2 | Non-TTY stdin refuses before removal and exits non-zero. | Runtime worktree cleanup seam | Add `testRuntimeWorktreeCleanRejectsNonTTY`; it fails on current code because `clean` is swallowed by the worktree subshell dispatcher. | The fixture keeps the worktree registered and asserts the error posture, so a headless remover or swallowed arg fails. |
| 3 | Dirty out-of-pool worktree is not removed and is reported as refused. | Runtime worktree cleanup seam | Add `testRuntimeWorktreeCleanLeavesDirtyOutOfPool`; it fails until cleanup filters or non-force removal preserves the dirty candidate. | The dirty checkout remains registered after the command; a force remove or bad clean predicate fails the state assertion. |
| 3 | Leased pool entry, warm unleased pool entry, and repo root are not offered as cleanup candidates. | Shared classification seam | Add `TestClassifyRegisteredWorktrees`; it fails until one classifier tags root, pool warm, pool leased, and out-of-pool cases distinctly. | Status and cleanup must derive from one classifier, and the tags make the negative candidate set explicit. |
| 4 | `bench worktree badverb` exits 2 and names valid forms. | Runtime worktree cleanup seam | Add `testRuntimeWorktreeRejectsUnknownArg`; it fails on current code because unknown args are ignored and the subshell path runs. | A typo cannot pass as successful cleanup when the exit code and usage text are asserted. |
| 4 | Bare `bench worktree` still opens and releases the pooled subshell. | Runtime worktree cleanup seam | Existing `testRuntimeWorktreeLeaseReuse` remains green; add a focused assertion if dispatch refactor risks bypassing `Subshell`. | This guards the compatibility surface while the new arg dispatch is introduced. |
| 5 | Out-of-pool status row names `bench worktree clean`. | Runtime status seam | Extend `testRuntimeStatusWarmPool`; it fails on current code because `appendWorktree` emits `resume or clean up (bench worktree)`. | The rendered dashboard must point at the verb that can clear out-of-pool worktrees. |
| 5 | Leased-only status row does not name `bench worktree clean`. | Runtime status seam | Extend `testRuntimeStatusWarmPool`; it fails until leased pool entries get their own action text. | The row stays honest by not promising that cleanup touches live or reclaimable pool state. |
| 5 | Orphaned scratch branch action remains `delete scratch branch` even when other worktree classes are present. | Runtime status seam | Extend `testRuntimeStatusWarmPool`; it fails if the split rows collapse orphaned branches into a combined cleanup action. | The orphaned-branch class has a different remediation and must not inherit the cleanup verb. |
| 6 | Cleanup and status both consume the shared classifier rather than separately deriving out-of-pool state. | Shared classification seam | Add a package-level test around classifier output and update status tests to build rows from those tags; code review verifies no second pool-prefix plus lease-file derivation remains. | The classifier output is the one observable fact both callers share, so duplicated classification logic has no local test to hide behind. |
| edge of 1 | Out-of-pool path containing spaces or glob characters is printed intact and removed by argument vector, not shell splitting. | Runtime worktree cleanup seam | Add a space-path fixture variant to the confirmed removal test. | Git receives the path as one arg; a shell-split implementation fails to remove or prints a mangled candidate. |
| edge of 1 | Registered worktree directory already missing is handled by prune without failing cleanup. | Runtime worktree cleanup seam | Add `testRuntimeWorktreeCleanPrunesMissingRegistration`. | The registry no longer lists the missing worktree after prune, proving clean handles stale git metadata. |
| edge of 1 | Second `bench worktree clean` after successful cleanup finds nothing and exits 0. | Runtime worktree cleanup seam | Extend confirmed removal test to rerun clean. | Cleanup is idempotent, so repeated status advice cannot create a false failure. |
| edge of 1 | Running from a nested cwd resolves the repo root and the same pool. | Runtime worktree cleanup seam | Add nested-cwd execution to one cleanup contract. | A cwd-sensitive implementation misses candidates or classifies the current directory incorrectly. |
| edge of 1 | Candidate that becomes dirty after listing is not forced and is reported as refused. | Runtime worktree cleanup seam | Use an injected remover or package-level cleanup seam test to dirty the candidate before remove. | The non-force safety floor is preserved even across a race between listing and removal. |
| edge of 2 | Interrupt before confirmation removes nothing. | Runtime worktree cleanup seam | Covered by non-TTY refusal and decline/no-confirmation path; no PTY signal helper is built for this spec. | The command does no removal before confirmation, so the tested pre-confirm exits pin the safety invariant. |
| edge of 4 | Unquoted multi-word worktree arguments are rejected as unknown args, not collapsed or split by the command. | Runtime worktree cleanup seam | Add `bench worktree clean extra` and `bench worktree bad verb` usage cases. | The dispatcher consumes the args slice directly and never silently falls through. |

### Edge inventory
Error path, empty/absent input, boundary values, malformed input,
interrupted/partial state, re-run idempotency, hostile environment, and the
profile's shell-CLI hostile inputs were walked. Each in-scope edge is mapped
above.

Won't handle: literal PTY prompt rendering - raw terminal interaction stays
manual-verify because non-TTY refusal and injected forced-terminal confirmation
cover the write/no-write behavior and no reusable PTY helper exists.

Won't handle: leased-entry owner liveness - reclaim remains owned by pool
acquire, and status only needs to avoid naming cleanup for leased entries.

Won't handle: a separate `bench worktree prune` verb - prune is folded into
`bench worktree clean`, and a standalone prune command is a separate capability.

## Out of scope
- **Stale leased-entry owner detection** - this is a separate pool-liveness
  capability because `Acquire` already owns reclaim and FT11 only makes status
  actions honest - estimate: 2 edits, 1 gate run.
- **Standalone `bench worktree prune`** - this is a separate command surface
  because FT11 folds prune into the human-attended cleanup verb - estimate:
  2 edits, 1 gate run.
- **Reusable PTY test helper** - this is a separate testing-infrastructure
  capability because this spec can assert refusal and forced-confirm paths
  without raw terminal automation - estimate: 3 edits, 1 gate run.
