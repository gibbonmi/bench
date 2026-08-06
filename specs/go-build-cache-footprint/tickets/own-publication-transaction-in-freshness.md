# Own the publication transaction in the freshness publisher

Blocked by: none
Ownership fence: `scripts/go-build.sh`, `internal/freshness/freshness.go`, `internal/freshness/freshness_test.go`, `cmd/bench/freshness_publish.go`, `internal/contract/surface/artifact/posture/mode_test.go`
Integration surfaces: sealed publication invocation→`scripts/go-build.sh`; transaction owner→`internal/freshness/freshness.go`; failure and interrupt observation→`internal/contract/surface/artifact/posture/mode_test.go`; command-entry signal wiring→`cmd/bench/freshness_publish.go`; gate attested-build publisher as unchanged `Publish` caller→existing path `internal/gate/build_attestation.go` + PT2
Contracts: the source root and literal destination cross `scripts/go-build.sh`→`internal/freshness/freshness.go` through the unchanged staged `freshness-publish` grammar, asserted by PT1 and PT2 against the real builder

## What to build

Repair for review findings `spec-01-existing-output-replacement-not-atomic` and
`spec-02-sigint-after-publish-does-not-restore-pair`, under the closed reviewer
decision recorded in `/tmp/bench-ft195-review-debug-receipt.md`: the Go
publication child owns the sealed transaction. `freshness.Publish` holds backup,
atomic executable+seal pair replacement, and signal-safe restore entirely inside
the publication operation, so a handled termination between executable install
and seal write restores the prior pair without any shell cooperation. The
builder's sealed path stops owning rollback state. Artifact mode cannot execute
the staged output, so its unsealed install becomes a single atomic rename with
no move-aside — no window in which the prior output is absent. The fence spans
four directories because the transaction is one observable outcome across the
builder invocation, the command entry, the publication owner, and the contract
that observes them; any thinner slice is a horizontal layer that cannot land
green alone. Two acceptance rows carry two findings because both prove the same
single-owner transaction from its two failure sides.

## Acceptance

- [ ] [PT1] Replacing an existing output in either mode leaves no observable window in which the prior executable (or, sealed mode, the prior consistent pair) is absent: no move-aside of the live output, and a failure at any stage before completion leaves the prior state byte-for-byte.
- [ ] [PT2] A handled termination signal delivered to the sealed publication between executable install and seal write restores the prior executable+seal pair byte-for-byte, and the restore is owned by the publication operation itself — the invoking shell holds no sealed-path rollback state.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PT1 | restore the artifact-mode move-aside of the existing output before installing the staged one | the builder failure table with fingerprints and a bounded marker-controlled blocker | fingerprint the existing output, drive the real builder to fail or block at the install stage under a deadline, expect the absence window or lost prior output to fail the table |
| PT2 | drop the publication operation's termination restore path | the interrupt contract with a bounded marker-controlled blocker | fingerprint the prior pair, block the sealed publication between executable install and seal write under a deadline, deliver the handled signal, expect the unrestored pair to fail the fingerprint comparison |
