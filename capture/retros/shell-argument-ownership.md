## Outcome

FT354 landed at `ff5b1b1e1`. Every printed shell argument in the enumerated producers now quotes through `sanitize.ShellQuote`. The AXI known-argument renderer switched to that owner in one green ticket with its complete fixture closure, and `axi.ShellQuote` and `freshness.shellQuote` are removed. The recovery line renders through the help-row renderer, and its refused-path fallback quotes each fixed word and keeps the request slot bare. All 44 acceptance rows reconciled, and S34 closed through the review-owned definition census. The seven final verifications passed, and the landing gate was green.

The build ran across several sessions. Chunks S-A, S-B1, S-B2, and S-B4 closed before this session. This session closed S-B4 at its checkpoint, then built and reviewed S-D1 and S-E, and landed.

## Cost

- landing: ff5b1b1e1d6fdd525cd00661c05d08b5130bfd6c, range ff5b1b1e1d6fdd525cd00661c05d08b5130bfd6c^1..ff5b1b1e1d6fdd525cd00661c05d08b5130bfd6c^2
- wall: 2913 min
- commits: 62 total, 50 bookkeeping, 12 change
- graph: 7 tickets, chain 5, 6 chunks, 4 serial checkpoints
- dispatch: 7 authors, 3 repairs
- trace: commit 67 spans in 663741 ms, lane 0 spans in 0 ms, landing 2 spans in 347801 ms
- wall per ticket: 416.2 min against 60 min, over budget
- bookkeeping per ticket: 7.1 against 7.1

The wall time spans several sessions and the idle time between them, so the per-ticket wall figure overstates the work. Delegate tokens in this session, from the native returns: the ticket 06 author at 87k, the ticket 04 author at 89k, the ticket 07 author at 376k over three attempts, the repair session at 172k over two commits, and four review sessions at 62k to 137k each. Provider cost is unknown.

## Gate-stage timings

- landing: commit ff5b1b1e1d6fdd525cd00661c05d08b5130bfd6c, trace 4b95404c29c79394075d190099b24cce
- gofmt: 307 ms
- vet: 4611 ms
- test: 245767 ms
- race: 9945 ms
- system: 48049 ms
- shellcheck: 593 ms

The S-B4 and S-D1 checkpoint gates were green, with the system phase at 49.2 s and 50.2 s. The final system verification took 50.4 s over 288 tests.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| unknown | unknown | unknown | unknown | unknown | unknown |

The orchestrator declared no expected repair rounds and no confidence, so this build has no labeled pair. Brier mean: unknown, over 0 pairs, with 0 abstentions.

Ticket 06 met its rows in one pass, and its four plan probes went red and passed after restoration. Ticket 04 met its row in one pass and added the shared sentinel helper that the S-A advice asked for. Ticket 07 used three attempts. The first two attempts stopped correctly on six exact-byte readers outside the fence, and the third attempt committed the full closure with sixteen red-then-restored mutants. The S-E repair session closed four findings in one commit, and it flagged a duplicated placeholder that its own repair added.

## Coordinator catches

- The ticket 07 fence stops were real. All six readers sat inside the directory of a signed-off `Writes:` entry, so each widen took a plan commit and no reviewer stop.
- The repair session flagged its own repeated `<request>` literal and string compare. The orchestrator widened ticket 07 with `path.go` and kept the fix inside the same repair cycle, so no second review round was necessary.
- The ticket 04 assignment came from before ticket 04 joined S-E. The verb refused a new occurrence without a trigger, so ticket 04 kept opus / medium, and the orchestrator flagged the effort for veto.
- The orchestrator charged the record commit `4ac0f5b42` as the S-D1 chunk base. The review preflight ran green on the correct source tip `18517a9a7`, and the orchestrator recorded the error with `bench learning`.
- The first landing refused because the record had no final verification. The orchestrator ran the seven final entries, recorded them and the completion, and landed again.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-quote-tree-build-arguments.md | 0 | none |
| 02-quote-selected-history.md | 0 | none |
| 03-quote-citation-commands.md | 1 | delegate-error |
| 04-quote-recovery-operands.md | 0 | none |
| 05-quote-evidence-commands.md | 1 | delegate-error |
| 06-share-rebuild-and-fixture-quotes.md | 0 | none |
| 07-switch-axi-output-atomically.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Let `bench preflight build --charge` accept a checkout whose only dirty paths are inside the assigned ticket fence, so a mid-ticket widen does not force a stash and restore cycle; this build did two.
  Feeds: new
- Make `bench spec widen` also add the widened path to the spec's ownership fence union, or derive that union from the ticket `Writes:` lines; S-E-SPEC-2 was this drift.
  Feeds: new
- Make `bench worktree land` refuse before the gate when the record has no final verification or completion entry, and print the seven final requirement ids; the first landing ran a gate before it refused.
  Feeds: new
- Census account, 17 raw calls. Owned gap, 1 `for` call: the orchestrator looped `bench record verification --final` over seven requirements. Add a form that records several final requirements from one result file.
  Feeds: new
- Census account, missed form, 1 `sed` call: the orchestrator changed a count in the spec and a ticket with `sed -i`, where the Edit tool does the job. Add a line to the orchestrator guidance that prose edits use the file tools.
  Feeds: none
- Census account, sanctioned reads: 12 `rg` calls need no Bench form. Unknown purpose: 3 `cat` calls from delegates; no change.
  Feeds: none

### Skills

- In craft-spec, make the build-entry census of rendered-command readers search testdata goldens and partial-command substrings, such as `--apply` and fingerprint fragments, as well as full command literals; the census missed six readers.
  Feeds: new
- In craft-tickets, give the switch ticket of an output change an attempt cap that counts only failed attempts, so a correct fence stop does not consume the cap.
  Feeds: new

### Process

- Derive each chunk base from the predecessor chunk entry's recorded tip in the review record, never from the current head after the predecessor's records.
  Feeds: none
- When the plan probes are the witnesses, write the witness section in the review-record commit or the next assignment commit, not in a separate commit; this build paid three such commits.
  Feeds: none