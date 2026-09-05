# Retro: ledger-settle-policy

## Outcome

The landing `3ed9140f` published the spec on 2026-09-05 over the frozen pair `a928bebc`
to `8827e12e`. One ledger transaction now owns the lock, the read mode, and the atomic
replacement for the seven intent mutators. The ledger schema lives in the leaf package
`internal/intent/ledger` with aliases in `internal/intent`, so no importer moved. The
admission rules live in `internal/intent/admissionpolicy` with typed liveness facts, and
the settle decision lives in `internal/landing/settlepolicy`. A conflicted landing names
why a settle refused, and a plain conflict keeps its bare kind. `CONTEXT.md` gained
**ledger transaction** and **settle verdict**.

The build found that the widened policy census reds `lifecyclepolicy`, which now imports
the leaf. It also found that a plain-conflict journey pins the bare refusal text, so a
reason surfaces only when a capture path engages the policy. The spec stays
`implemented` as the veto surface for fourteen recorded decisions and seven open
review notes in `reviews/ledger-settle-policy.md`. FT302 stays open until the reviewer
retires this spec.

## Gate-stage timings

| stage | landing gate |
| --- | --- |
| gofmt | 112 ms |
| vet | 1021 ms |
| race | 2423 ms |
| test | 67146 ms |
| system | 25942 ms |
| shellcheck | 596 ms |

Nine fold gates ran the test stage between 66624 ms and 83977 ms. Two early folds ran red
on a landing test that refused as `infrastructure`. A delegate's tests ran on the same
machine at the time, and each red cost one full gate. One fold ran red on a raw `t.Skip`, and one
on a coverage-row citation shape.

## Ticket-versus-spec-slice and delegate performance

Eight ticket charges ran on `opus`: two at low and six at medium. Five landed first-pass
on behavior (the census pattern, the leaf move with its fence extension, the admission
policy, the settle policy, the surfaced reason). The transaction ticket left a raw
`t.Skip` that the conformance check caught. The mutator ticket carried two comments that
narrated the migration. The adapter ticket deleted two union journeys that no policy
table covered, against its own "delete no partition" line.

Two repair charges on `opus` at low and medium landed first-pass. Three review axes and one
scoped re-review ran on `opus` at medium and returned seventeen findings, of which seven
became repairs. Every delegate probe bit. The coordinator ran eleven probes, and nine bit;
the two silent greens were inherited gaps at the base, not missing rows.

## Coordinator catches

- The spec said the three existing policy children grade against the widened census
  rule. `lifecyclepolicy` imports `internal/intent`, so the first fold ran red.
- A plain-conflict journey asserts `detail=composition conflict: textual,next=`, so the
  reason needed an engagement rule before ticket 08 started.
- Two fold reds were load flakes, proved green in isolation, and the third fold on a quiet
  machine passed.
- Ticket 02's Writes line named a deleted file and missed a bound package's registry
  closure, and the review preflight ran red on both.
- A coverage row cited a file without its name list, and the conformance check ran red.
- The tolerant read decoded the entries on every read, a behavior change the delegate
  reported and the spec now records with a row.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-forbid-the-new-parent-adapters-in-the-policy-census | 1 | spec-row |
| 02-extract-the-leaf-ledger-package | 0 | none |
| 03-add-the-ledger-transaction | 1 | delegate-error |
| 04-add-the-admission-policy-child | 0 | none |
| 05-migrate-the-assignment-mutators | 1 | delegate-error |
| 06-add-the-settle-policy-child | 1 | spec-row |
| 07-migrate-the-composition-adapter-onto-the-settle-policy | 1 | delegate-error |
| 08-surface-the-settle-refusal-reason | 0 | none |
| 09-repair-the-intent-review-findings | 0 | none |
| 10-repair-the-landing-review-findings | 0 | none |

## Agent-experience improvements

### Bench CLI

- Add `bench structure --path <prefix>` so that a fenced ticket reads its own budget rows
  instead of the whole debt census. The census entry `ft302-c34-spec landed with 23 raw
  calls` proposes it.
  Feeds: new
- Let `bench worktree path` print a one-line note that the path serves the file tools
  only, because the follow-on hook refuses that path in a Bash command.
  Feeds: new
- Let the exec child receive the worktree as `PWD`, because a script that trusts `PWD`
  resolves paths into the primary checkout.
  Feeds: new

### Skills

- Make `craft-spec` require a row that widens a forbidden-import pattern to name the
  enumeration of its current importers across every graded package.
  Feeds: new
- Make `craft-tickets` give a deleted file's Writes entry a marker the preflight accepts.
  The plain path stops resolving after the move.
  Feeds: new

### Process

- Run a fold or a whole-tree gate only when no delegate test run is live. Two landing
  tests refused as `infrastructure` under load and proved green in isolation.
  Feeds: none
- Make a spec that changes a rendered message enumerate the exact-match tests on that
  text. A closure-headroom claim of "no edit" hid one exact match.
  Feeds: none
