## Outcome

FT363 landed at `b062b00b2`. The four cleanup adapters (explicit recorded, explicit unrecorded, landed and folded-sibling, unclaimed) visit their members through one package-local ordered traversal in `internal/worktree/cleanup_members.go`. A frozen all-mode baseline of 31 records, captured from the original code, compares each mode. A copy-survival census fails when an original member loop remains. One production owner builds the options input and runs the versioned hash loop for the explicit, landed, and unclaimed fingerprints, and known-answer digests pin the bytes. All 29 acceptance rows reconciled, and the complete checkpoint was green.

The reviewer directed opus / medium for every author and every review axis. The orchestrator resliced twice inside the approved behavior: ticket 00 took the baseline capture from ticket 01 at build entry, and ticket 05 took the census, the fingerprint owner, and the R26 completion from ticket 04 at the C3 checkpoint.

## Cost

- landing: b062b00b2414eb923d23d077e9a313f8f56ca821, range b062b00b2414eb923d23d077e9a313f8f56ca821^1..b062b00b2414eb923d23d077e9a313f8f56ca821^2
- wall: 211 min
- commits: 60 total, 51 bookkeeping, 9 change
- graph: 6 tickets, chain 6, 4 chunks, 4 serial checkpoints
- dispatch: 6 authors, 3 repairs
- trace: commit 61 spans in 633743 ms, lane 0 spans in 0 ms, landing 1 spans in 12663 ms
- wall per ticket: 35.3 min against 60 min
- bookkeeping per ticket: 8.5 against 7.1, over budget

Delegate tokens, from the native returns: 3 staleness slices at 44k to 66k each, 6 authors at 110k to 230k each, 4 repairs at 66k to 110k each, and 8 unified review sessions at 36k to 119k each. Provider cost is unknown.

## Gate-stage timings

- landing: commit b062b00b2414eb923d23d077e9a313f8f56ca821, trace 18673ac0bd1d1535cf2f0268e643b21b
- The system phase took 66.2 s at the C1 checkpoint, 51.3 s at C2, 53.9 s at C3, and 49.5 s at the complete checkpoint. The landing reused the fresh green verdict.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| C1 | the chunk takes exactly one repair round | met | 6 | true | opus / medium / orchestrator |
| C2 | the chunk takes exactly one repair round | not met, zero rounds | 6 | false | opus / medium / orchestrator |
| C3 | the chunk takes exactly one repair round | met | 6 | true | opus / medium / orchestrator |
| C4 | the chunk takes exactly one repair round | met | 5 | true | opus / medium / orchestrator |

Brier mean: 0.233, over 4 pairs, with 0 abstentions.

Ticket 00 stopped twice before its commit, each time correctly, on a fence gap: the package test-count pin and the 400-line structure bound. Ticket 01 met its rows in one pass, but its review found two gaps that the ticket text caused. Tickets 02 and 03 met their rows in one pass, and each reported a true spec-versus-code divergence. Ticket 04 met its rows in one pass. Ticket 05 met its rows, but its census had a closure bypass and its digest check had no known answer.

## Coordinator catches

- The ticket 00 fence stops were real, and both widens took no reviewer stop.
- The ticket 02 author showed that the unrecorded adapter reads no run; the orchestrator corrected the spec sentence.
- The C3 repair reported a new duplication of the options input and the hash loop; the orchestrator routed it to ticket 05 instead of a second C3 repair.
- The C4 axis graded the census bypass and the missing known-answer digest as advice; the orchestrator promoted both, because each fails a binding spec requirement.
- The C4 test-only confirmation refused, because a prose-only spec correction sat in the repair window; one short unified round replaced it.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 00-capture-cleanup-baseline.md | 0 | none |
| 01-share-recorded-traversal.md | 1 | one-source |
| 02-migrate-unrecorded-traversal.md | 0 | none |
| 03-migrate-landed-traversal.md | 1 | one-source |
| 04-migrate-unclaimed-and-contract.md | 0 | none |
| 05-contract-traversal-owner.md | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Make `bench record confirmation` treat a prose-only spec correction outside the plan fence as evidence-only, so a test-only repair keeps its coordinator route; this build paid one confirming round for it.
  Feeds: new
- Make an assignment-only plan change leave the chunk freeze current, so `bench record amendment` stops running after every assignment; this build ran it eight times.
  Feeds: new
- Let `bench preflight build --charge` accept a checkout whose only dirty paths are untracked files inside the assigned ticket fence, so a mid-ticket widen does not force a copy-aside and restore cycle; this build did two.
  Feeds: new
- Make `bench record ticket` stage the last probe response for a mutant set, or name the excerpt step in its refusal; three authors met the refusal.
  Feeds: new
- Census account, 56 raw calls. Owned gap, 21 `cat` and 7 `python3` calls: the orchestrator appended review pickup prose and plan-fence edits by hand. Add a record form that renders witness and round sections from the record entries, and a plan verb that adds a ticket to a chunk and moves requirements between tickets.
  Feeds: new
- Census account, sanctioned reads: 17 `rg` calls and 1 `mkdir` call need no Bench form. Unknown purpose: 4 `cd`, 3 `sed`, 2 `ls`, and 1 `git log` calls from delegates; no change.
  Feeds: none

### Skills

- In craft-tickets, put the package test-count pin file in the `Writes:` line of each ticket that adds a top-level test in a pinned package, and pre-fence one overflow test file when a planned test file is near the 400-line bound.
  Feeds: new
- In craft-tickets, quote each must-fail clause of a covered spec row in the ticket, and require one witness that drives each migrated adapter for each mode-agnostic row.
  Feeds: new

### Process

- At the first chunk gap, widen every successor ticket for the shared pin and overflow files in one plan commit; C2 to C4 then took zero widen stops.
  Feeds: none
- Make the spec stage declare the implementation line as a tier, not as a model id of one harness; this spec named an unbound Codex id, and the build needed reviewer direction before its first dispatch.
  Feeds: new