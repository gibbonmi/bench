## Outcome

FT311 prepares pinned build charges, proposes ownership closure, gathers one shared review packet, and routes approved work through native phases.

Landing `4d0354d1784e8f5959759dc074b6ace1aeb98462` completed on 2026-09-09 UTC, or 2026-09-08 locally. The feature spec status is implemented.

The first source gate at `6f8baa77` stopped on one stale removed command token. The other five phases passed.

The Codex root had misread the user stop and made that commit without checks. The user clarified that only cross-harness checks had stopped.

Commit `4d870d18` added `CONTEXT.md` and `docs/research/ft311-session-orchestration.md` to the ownership fence. The first landing refused these pre-reviewed, user-approved additions.

Three independent Terra/high axes returned zero findings on the final source. They later acknowledged the closure and final exact-fence deltas with zero findings.

Prose checks passed for the fence extension. Linked dogfood commit `6c95593` completed with exit 0.

The dogfood branch was `bench/shift-20260908-215038`. It produced one commit in two iterations without recovery.

Landing returned exit 3 only because 17 ignored source-release residuals remained. The complete source was preserved at `/tmp/ft311-landed-source-preserved-20260909`.

Cleanup removed those files with fingerprint `ca899282884e069a5dc775fb964d5e324ffb1a1686843eec770d2e0cfb411d08`. The installed CLI was rebuilt, and `doctor --fix` published the current broker.

A follow-up `land --resume` refused `missing-terminal-receipt` because cleanup used `clean`, rather than `release`. The landed publication remains intact.

This required no product repair or new gate.

## Gate-stage timings

| run | stage | elapsed | result |
| --- | --- | --- | --- |
| source gate at `6f8baa77` | gofmt | 100 ms | green |
| source gate at `6f8baa77` | go vet | 1,164 ms | green |
| source gate at `6f8baa77` | go test | 86,615 ms | red: stale token at spec line 411 |
| source gate at `6f8baa77` | race | 6,208 ms | green |
| source gate at `6f8baa77` | system | 37,934 ms | green |
| source gate at `6f8baa77` | shellcheck | 617 ms | green |
| source gate at `2e5de3ca` | gofmt | 146 ms | green |
| source gate at `2e5de3ca` | go vet | 1,429 ms | green |
| source gate at `2e5de3ca` | go test | 81,472 ms | green |
| source gate at `2e5de3ca` | race | 3,566 ms | green |
| source gate at `2e5de3ca` | system | 36,350 ms | green |
| source gate at `2e5de3ca` | shellcheck | 596 ms | green |
| linked dogfood | first gate | 2,358 ms | green |
| linked dogfood | second gate | 99 ms | green |
| final landing gate | gofmt | 101 ms | green |
| final landing gate | go vet | 1,383 ms | green |
| final landing gate | go test | 75,883 ms | green |
| final landing gate | race | 3,044 ms | green |
| final landing gate | system | 39,011 ms | green |
| final landing gate | shellcheck | 571 ms | green |

The final landing gate reported six optional capability skips and zero environment skips.

## Ticket-versus-spec-slice and delegate performance

| ticket | model / effort | result |
| --- | --- | --- |
| 1. Prepare one build charge | Sol / high | Produced the pinned build form. Review later required packet composition and charge-cell repairs. |
| 2. Propose ticket ownership closure | Sol / high | Needed one repair for DP10-DP14 coverage and shared orchestration. The repair supplied exact-baseline reds, green tests, and restored mutations. |
| 3. Prepare shared review evidence | unknown / unknown | Produced the shared review packet. Review later required one packet emitter and explicit cross-package cell grading. |
| 4. Consume charges and triage repairs | Terra / high | Returned an implementation checkpoint with live exercises open. Review found that guidance lacked the triage input contract. |
| 5. Dispatch prepared native reviews | Opus / high | Added native dispatch and capable-harness handoff guidance. Review found that the phase did not consume the shared packet. |
| 6. Repair reviewed metadata | unknown / unknown | Repaired DP7 and DP8 metadata without changing behavior. |

The ticket charges kept implementation ownership bounded. Later review showed that shared charge behavior and phase guidance needed complete-slice checks across tickets.

## Coordinator catches

- Ticket 2 returned because its public-command tests missed DP10-DP14 partitions and charge and proposal had separate retry orchestration.
- Review collapsed 20 raw findings into 16 repair targets.
- The largest targets covered packet rendering, charge-cell assertions, prepared-evidence consumption, and stored live evidence.
- Linked repositories do not carry the kit's compiled ticket-binding registry. The coordinator rejected a target-repository source requirement.
- The source gate found the stale command token after the Codex root made the unchecked local commit.
- Terra/high repaired the tier helper and changelog, then repaired the stale-token wording in `676484a`.
- DP20 ran through an actual prepared native dispatch.
- DP21 withheld Coverage while two axes returned. The later Coverage return restored a complete zero-finding result.
- The coordinator rejected a false README seal diagnosis because README is not a build input. It routed dogfood to a linked repository.
- Failed dogfood setups did not count as acceptance evidence.
- Rejected setups used the wrong pool, mismatched the installed kit, lacked inputs, or placed cache inside the graded root.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | ---: | --- |
| 1.md | unknown | unknown |
| 2.md | 1 | delegate-error |
| 3.md | unknown | unknown |
| 4.md | unknown | unknown |
| 5.md | unknown | unknown |
| 6.md | unknown | unknown |

## Agent-experience improvements

### Bench CLI

- Add `bench worktree exec <label> -- rg --files` to census guidance for raw `ls` calls; the confirmed census is two calls. The captured learning is `ft311-log-research: 2 raw calls`.
  Feeds: new

### Skills

- Require repair charges to replay acceptance tests against exact pre-change production when the original ticket lacks retained pre-edit red evidence.
  Feeds: new

### Process

- Record whether a user stop covers cross-harness work or all verification before changing the completion plan.
  Feeds: new
