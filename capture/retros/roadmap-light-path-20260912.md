# Roadmap light-path fixes retrospective

## Outcome

Eleven bounded fixes and the roadmap cleanup passed through isolated worktrees and serial landings.
The incoming delegate-c3 worktree landed before this batch changed overlapping kit content.
FT311 closed when automatic cleanup began retaining branches with no proven contribution.
Other roadmap owners retain their separate residual requirements.

The fixes cover child PWD, compiler attribution, root help, byte-safe filenames, and empty sibling preservation.
They also cover complete diagnostics, anchor attribution, private command fixtures, derived waits, ticket verification passes, and timeout recovery counts.
The ten original fix reviews ended clear across Standards, Spec, and Coverage.
The user stopped further reviews before the timeout repair received a result.
The roadmap review found one stale FT311 dependency sentence, which the coordinator removed.

Each full gate reported seven capability skips and no environment skips.
Optional shellcheck was unavailable. This batch did not perform release qualification.

## Gate-stage timings

Timings are milliseconds. The failed wait attempt is retained.

| operation | format | vet | tests | race | system | shellcheck | result |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| child PWD | 137 | 1320 | 136928 | 3417 | 30878 | 54 | green |
| filtered compiler | 201 | 1806 | 253905 | 6624 | 31800 | 43 | green |
| root help | 128 | 1074 | 137284 | 3477 | 31496 | 31 | green |
| NUL paths | 103 | 1022 | 137225 | 3313 | 31142 | 30 | green |
| empty sibling | 115 | 1014 | 131120 | 6077 | 32105 | 32 | green |
| full failures | 114 | 999 | 142457 | 5767 | 29379 | 39 | green |
| anchor merge | 113 | 1041 | 133439 | 3971 | 30737 | 38 | green |
| anchor attribution | 105 | 1037 | 133536 | 3596 | 37987 | 86 | green |
| command fixtures | 202 | 1700 | 225821 | 4356 | 46587 | 53 | green |
| derived waits — failed attempt | 158 | 1712 | 143708 | 3415 | 32681 | 40 | red |
| timeout recovery count | 121 | 1167 | 138523 | 3426 | 31306 | 43 | green |
| derived waits | 146 | 1065 | 142165 | 3551 | 32047 | 49 | green |
| ticket guidance | 106 | 1064 | 150878 | 3958 | 37124 | 51 | green |

The final cleanup gate follows this capture commit. Its result remains in the landing trace.

## Ticket-versus-spec-slice and delegate performance

Astra/high implemented the child-PWD ticket before the user moved substantive implementation into the main session.
The main session implemented five fixes and retained their tests and repairs.
Sol/high implemented five bounded fixes: filenames, help, command fixtures, waits, and timeout recovery counts.
One Sol/high delegate performed all three review axes for each original fix.

The first filename review found missing modified, copied, and renamed input cases.
The full-output review found a changelog placement defect and an escaping coverage gap.
Their authors repaired each accepted finding, and the same reviewer checked the repair.
The other fix reviews reported no findings.
The deadline landing exposed an existing timeout-test flake, which received its own ticket and author.

One reviewer attempted a nested dispatch early in the batch; the coordinator stopped it before using its output.
Three later review dispatches met the active-agent limit and performed no work.
A completed author turn released capacity for the final fresh reviewer.
Native token counts, prices, and comparable review latency remain unknown.
The results support another combined-review trial, but establish no cost or speed advantage.

## Coordinator catches

Every write delegate returned a behavioral mutation, and the coordinator used a different mutation kind at a different site.
Compiler errors from invalid probes were repaired and never counted as behavioral red evidence.
The coordinator checked source identities, tree cleanliness, and the final diff after each merge.
Existing ticket fixture anchors remained intact after the guidance was shortened to its unchanged 100-line budget.

Four changelog conflicts needed a native merge after Bench refused composition.
Bench commit graded the resolved content but left MERGE_HEAD pending behind a single-parent commit.
A native commit completed each merge; the coordinator verified identical tree content and current-main ancestry.
A conflict-free merge started another full gate. Further writes paused after that gate became visible.
A closeout assignment was created before the coordinator recognized that additional gate; its source was not edited during the gate.

The command-fixture author made two early shell follow-on calls outside the intended wrapper.
One Go call failed, and the other calls only read primary conformance files.
The author reran the steps through the required worktree venue.

The benchmark directory contains only its two-byte ignore file after the earlier authorized cleanup.
This batch kept compact assessment records and logs instead of new benchmark bundles.
Automatic benchmark expiry remains unimplemented.

The timeout-test diagnosis first proposed a cause before a repeatable short repro existed.
The coordinator recorded this ordering error, then reproduced the exact counter failure twice with a controlled startup delay.
Instrumentation measured one run before recovery and two afterward, confirming a fresh execution despite the failing fixed total.
The temporary instrumentation was restored before delegation.

The coordinator also caught a cleanup placement that could leak the short timeout between failing subtests.
The author moved cleanup into each subtest before the final focused checks.
The user stopped reviews and requested immediate landing of non-blocking work.

For the next drain, the user requested a discussion of prose enforcement that limits repair loops in casual implementations.
Consider an explicit repair cap, blocker criteria, and a stop rule for minor prose findings.
Such enforcement must retain behavioral guarantees and required gates without adding another paperwork loop.

## Repair attribution

Rounds count returned author passes after accepted semantic-review findings.
In-pass compiler, fixture, and integration repairs remain in the evidence without adding review rounds.

| ticket | rounds | causes |
| --- | --- | --- |
| exec-pwd | 0 | none |
| ft290-filtered-compile-diagnostic | 0 | none |
| ft89-root-land-help | 0 | none |
| ft307-nul-touched-paths | 1 | delegate-error |
| ft311-zero-contribution-cleanup | 0 | none |
| ft290-complete-failure-output | 1 | other |
| ft120-anchor-diagnostic-attribution | 0 | none |
| ft120-isolated-command-fixtures | 0 | none |
| ft115-derived-outer-waits | 0 | none |
| ft300-ticket-relocation-passes | 0 | none |
| timeout-recovery-count | 0 | none |
| roadmap closeout | 1 | other |

## Agent-experience improvements

### Bench CLI

- Complete merge ancestry with the lane-graded content so the suggested conflict-repair command leaves no pending MERGE_HEAD.
  Feeds: FT258
- Make the full-gate cost of conflict-free worktree merge visible before execution.
  Feeds: FT254
- Keep the exact landing census beside each source identity in the local assessment.
  Feeds: none

### Skills

- Preserve existing immutable mutation anchors when reducing guidance prose to its current budget.
  Feeds: none

### Process

- Define a bounded repair policy for casual work that distinguishes blockers from minor prose issues and stops review loops after required checks pass.
  Feeds: new

- Continue the user-requested single-reviewer trial with separate axis reports and independent coordinator mutations.
  Feeds: none
- End completed author turns before dispatching fresh review contexts.
  Feeds: none
- Record unavailable token counts, costs, and comparative latency as unknown.
  Feeds: none
