# Reclaim prior terminal runs and report a partial apply

Blocked by: none
Ownership fence: `internal/specbuild`
Contracts: the set of run records the enumeration answers for crosses `internal/specbuild/state.go`→`internal/specbuild/reclaim.go` and is asserted by RH1 against a slug restarted through the real lifecycle rather than a hand-written history entry; the applied-ref report crosses `ApplyReclaim`→its receipt in the same package and is asserted by RH4 over a deletion that fails partway
Assumptions: `provisionalResidue` stays the one enumeration shared by promotion's reclamation and the maintainer pass; a history entry is already validated as terminal by the record loader so this ticket consumes that verdict rather than re-deriving terminality

## What to build

Two defects the authoritative review found in candidate `2f486147`.

**Prior terminal runs are invisible.** `record.History` holds earlier terminal runs for the
same slug, and the loader already validates each as terminal. But `provisionalResidue`
reads only the live record's assignments, checkpoint refs, and candidate. So after a slug
restarts, the prior run's candidate and checkpoint refs never appear in the plan at all,
and its assignment branches surface only as `unclassified` and are retained. That is a
story 4 violation rather than its "no owning record" clause: the history entry **is** the
owning record, and it says terminal.

Extend the enumeration to answer for prior terminal runs as well as the live one, keeping
it the single shared source. Terminality still comes from each record, never from a name.

**A partial apply reports nothing.** `ApplyReclaim` validates one fingerprint and then
deletes refs one at a time, returning on the first compare-and-swap failure. If the first
ref deletes and the second has moved, the first is gone, no applied receipt comes back, and
retrying the original fingerprint refuses because the inventory changed. The operator is
told nothing happened while something did.

Do not invent a transaction Git cannot give you. Make the outcome honest and convergent:
report which refs were deleted before the failure, and make a fresh plan plus its
fingerprint complete the remainder. The convergence must be tested, not asserted in prose.

## Acceptance

- [ ] [RH1] a slug whose prior run is terminal has that run's candidate and checkpoint refs reported in the plan, with its assignment branches classified from the history record rather than left unclassified.
- [ ] [RH2] the live record's own refs keep their existing classifications, so extending the enumeration does not reclassify anything it already answered for.
- [ ] [RH3] a prior run whose record is not terminal keeps its refs retained, so history is not a blanket licence to delete.
- [ ] [RH4] an apply interrupted by a compare-and-swap failure reports the refs it did delete and returns a non-success outcome naming the drift.
- [ ] [RH5] re-planning after such a partial apply yields a fingerprint that applies cleanly and removes the remainder.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RH1 | enumerate the live record only, as the candidate does today | the history-residue test | drop the history arm, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the prior-run candidate assertion to fail |
| RH2 | classify every history-located ref as reclaimable regardless of which record owns it | the live-record regression test | collapse the per-record verdict, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the live-classification assertion to fail |
| RH3 | treat presence in history as terminal without reading the record's flag | the non-terminal history test | drop the flag check on history entries, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the retained assertion to fail |
| RH4 | return the bare error on a mid-loop failure, discarding the deleted-ref list | the partial-apply report test | drop the accumulated list from the returned plan, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the reported-deletions assertion to fail |
| RH5 | keep the pre-failure inventory in the re-plan so the stale fingerprint still matches | the convergence test | cache the inventory across plans, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the fresh-fingerprint assertion to fail |
