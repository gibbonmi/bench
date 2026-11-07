## Outcome

FT446 landed at 6198ab5. A retired protection store now keeps a retired record that names its resource. The nested sweep and the last uncontended cache release delete each retired store whose resource is gone. `RunInherited` is the one inherited registration owner for both runners, and the closure check refuses a hand-built copy. All 31 coverage rows reconciled. The final package run passed 143 packages, and the system suite passed 288 tests.

## Cost

- landing: 6198ab5e50c65d51e0f4187c02cf9af4167759f9, range 6198ab5e50c65d51e0f4187c02cf9af4167759f9^1..6198ab5e50c65d51e0f4187c02cf9af4167759f9^2
- wall: 118 min
- commits: 31 total, 24 bookkeeping, 7 change
- graph: 5 tickets, chain 3, 4 chunks, 3 serial checkpoints
- dispatch: 5 authors, 2 repairs
- trace: unknown
- wall per ticket: 23.6 min against 60 min
- bookkeeping per ticket: 4.8 against 7.1

## Gate-stage timings

- landing: commit 6198ab5e50c65d51e0f4187c02cf9af4167759f9, trace cc498cc742acd012ad1bbf72a6bf5ae2
- Each chunk checkpoint ran the system stage in 47.8 to 49.5 seconds and shellcheck in about 0.5 seconds.
- The complete checkpoint ran the system stage in 49.0 seconds. The landing reused that fresh green.
- The ticket authors took 8.2, 10.3, 9.8, 8.6, and 4.9 minutes. The unified reviews took 4.3, 2.9, 4.7, and 5.2 minutes.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| 01-record-retired-store-path.md | RS1 | verified | 9 | held | opus / high / implementation |
| 01-record-retired-store-path.md | RS2 | verified | 9 | held | opus / high / implementation |
| 01-record-retired-store-path.md | RS3 | verified | 8 | held | opus / high / implementation |
| 01-record-retired-store-path.md | RS25 | verified | 9 | held | opus / high / implementation |
| 01-record-retired-store-path.md | RS26 | verified | 8 | held | opus / high / implementation |
| 02-sweep-retired-nested-stores.md | RS5 | verified | 9 | held | opus / high / implementation |
| 02-sweep-retired-nested-stores.md | RS6 | verified | 9 | held | opus / high / implementation |
| 02-sweep-retired-nested-stores.md | RS7 | verified | 9 | held | opus / high / implementation |
| 02-sweep-retired-nested-stores.md | RS8 | verified | 9 | held | opus / high / implementation |
| 02-sweep-retired-nested-stores.md | RS31 | verified | 9 | held | opus / high / implementation |
| 02-sweep-retired-nested-stores.md | RS4 | verified | 8 | refuted | opus / high / implementation |
| 03-discard-gone-retired-caches.md | RS9 | verified | 9 | held | opus / medium / implementation |
| 03-discard-gone-retired-caches.md | RS10 | verified | 9 | held | opus / medium / implementation |
| 03-discard-gone-retired-caches.md | RS30 | verified | 9 | held | opus / medium / implementation |
| 03-discard-gone-retired-caches.md | RS12 | verified | 9 | held | opus / medium / implementation |
| 03-discard-gone-retired-caches.md | RS13 | verified | 9 | held | opus / medium / implementation |
| 03-discard-gone-retired-caches.md | RS28 | verified | 9 | held | opus / medium / implementation |
| 03-discard-gone-retired-caches.md | RS11 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS16 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS17 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS19 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS20 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS21 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS22 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS23 | verified | 9 | held | opus / medium / implementation |
| 04-run-inherited-under-local-protocol.md | RS24 | verified | 9 | held | opus / medium / implementation |
| 05-refuse-hand-built-inherited-registrations.md | RS18 | verified | 9 | held | opus / medium / implementation |
| 05-refuse-hand-built-inherited-registrations.md | uncounted readers green | verified | 8 | held | opus / medium / implementation |
| 05-refuse-hand-built-inherited-registrations.md | RS29 | verified | 9 | held | opus / medium / implementation |
| 05-refuse-hand-built-inherited-registrations.md | copy in either runner reds | verified | 9 | held | opus / medium / implementation |
| 05-refuse-hand-built-inherited-registrations.md | RS24 | verified | 9 | held | opus / medium / implementation |

Brier mean: 0.033 over 31 pairs. Abstentions: 0.

Each ticket stayed inside its planned chunk. Ticket 03 needed one widen in the directory of its fence. No author used more than one attempt.

## Coordinator catches

- The staleness pass found one ticket 01 claim that could not hold: an export test file is not visible outside its package. The amendment named the literal file name in the checkout test.
- The RS-C3 reviewer returned a duplicated fixture as advice. The orchestrator accepted it as a finding under the one-source rule.
- The RS-C3 repair session found a mis-staged probe excerpt by its byte count, and restaged it.
- The orchestrator read the empty registration literal in the bounds runner as allowed under D4, and flagged that reading for reviewer veto.
- The RS-C2 repair session found that the live-generation branch of the nested sweep decodes a record without its validity check. That defect is outside the spec, and it went to the learnings journal.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-record-retired-store-path.md | 0 | none |
| 02-sweep-retired-nested-stores.md | 1 | delegate-error |
| 03-discard-gone-retired-caches.md | 1 | one-source |
| 04-run-inherited-under-local-protocol.md | 0 | none |
| 05-refuse-hand-built-inherited-registrations.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Exclude the spec, ticket, and review record paths from the test-only predicate of `bench record confirmation` and `bench record ticket`, so a recorded `Writes:` widen does not force a three-axis confirming round.
  Feeds: new
- Make `bench record excerpt` match the probe subject and its expected test to the named requirement, or print the staged mutation, so a mis-staged excerpt cannot pass silently.
  Feeds: new
- Print the inner response path from `bench worktree exec`, or let `bench record confirmation` take the last probe of the assignment, because the outer spill file refuses as a spilled frame.
  Feeds: new
- Let `bench record assignment` accept a dirty spec that holds only the version 2 execution block, so the block and the first assignment share one plan commit.
  Feeds: new
- Add one unified form of `bench record review` that reads one native report and splits its three axis sections, so a unified review takes one call in place of three calls and three excerpt files.
  Feeds: new
- Let `bench preflight review --charge` accept the frozen chunk tip when only record commits follow it, because the first call with the chunk tip always refuses.
  Feeds: new
- Census account, sanctioned read: 18 `rg` calls were text searches that the repository rules choose. No Bench form is proposed.
  Feeds: none
- Census account, owned gap: 10 `cat` calls appended orchestrator prose to the review record and wrote handoff state drafts. Add a `bench record note` form that appends one prose section to the review record outside the fence.
  Feeds: new
- Census account, missed form: 3 `sed` calls read or rewrote file lines in delegate sessions. Route a line read to `bench worktree read --lines`, and keep an edit on the file tools.
  Feeds: none

### Skills

- State in the author charge of `/bench-implement-spec` that plan probes on one subject file run as one `bench probe --mutants` set, with one `bench record excerpt --mutant` call for each earlier row, because RS-C2 cut two probe baselines that way.
  Feeds: new
- State in the review charge the citation, severity, and disposition rules directly, because two of four reviewers did not read the review skill and the repair policy.
  Feeds: new
- Name each hostile-input class of the project checklist as its own witness when a coverage row cites that checklist, because the RS4 witness cap let the validity class go untested.
  Feeds: new

### Process

- Let `write-spec` author the `bench-planning-facts` fence, or let the staleness pass omit itself when no code changed since the spec's proof commit, because the omission rule cannot fire today and this build paid three delegates for no drift.
  Feeds: new
- Let a ticket author report a same-directory test-file widen before its commit, so the orchestrator widens before review and a two-line fixture move does not take a review round and a repair round.
  Feeds: new
- Let an independent chunk author run in a sibling worktree while the predecessor chunk is under review, because RS-C4 had no dependency and each review left the author slot idle for about five minutes.
  Feeds: new