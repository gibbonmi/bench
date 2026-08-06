# Remove publication temp litter on handled signals

Blocked by: none
Ownership fence: `internal/freshness/freshness.go`, `internal/freshness/freshness_test.go`
Integration surfaces: publication temporary name families→`internal/freshness/freshness.go`; residue enumeration in interrupt and failure tests→`internal/freshness/freshness_test.go`
Contracts: the enumerated publication-temporary name families cross `internal/freshness/freshness.go`→`internal/freshness/freshness_test.go`, asserted by TL2 against the real publisher's own naming

## What to build

Repair for review findings `spec-04-seal-temp-survives-handled-signal`,
`spec-05-backup-pair-leaks-on-post-seal-signal`, and
`coverage-01-residue-glob-blind-to-seal-temp-names` — one temp-litter behavior
proven from both sides, which is why three findings ride two rows. The signal
handler runs `publicationExit` unconditionally, so it preempts both `writeSeal`'s
temp cleanup (leaking `<seal>.tmp-*` when the signal lands mid-seal) and the
transaction's deferred `close()` (leaking both `.bench-publish-backup-*` files
when the signal lands after the seal resolves). Make the handler remove every
publication temporary it owns — backups and any in-flight seal temp, with the
temp's name owned by the publication value, not re-derived — before exiting, in
both the restored and already-sealed outcomes. The residue assertions currently
glob only `.bench-*` and cannot see `bench.seal.tmp-*` names; enumerate every
temporary name family the publisher can create so new litter is red.

## Acceptance

- [ ] [TL1] A handled termination at any point of sealed publication leaves no publication temporary — no backup file and no seal temp — beside the destination, in both the restored-pair and the already-sealed outcomes.
- [ ] [TL2] The residue assertions enumerate every publication-temporary name family the publisher can create, derived from the publisher's own naming, so a planted `bench.seal.tmp-*` or backup leftover is red.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TL1 | drop the handler's temporary cleanup before exit | the interrupt residue contract | block sealed publication in the install-to-seal window under a deadline, deliver the handled signal, expect the leftover enumeration to fail on the surviving temporaries |
| TL2 | narrow the residue enumeration back to the `.bench-*` glob alone | a litter-planting check independent of the enumeration under test | plant a `bench.seal.tmp-` file beside the destination, run the residue assertion, expect the planted litter to be reported |
