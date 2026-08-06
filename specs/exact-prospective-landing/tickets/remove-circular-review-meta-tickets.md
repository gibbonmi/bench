# Remove circular review meta-tickets

Blocked by: close-adapter-blocker-metadata.md, preserve-executable-spec-mode.md
Ownership fence: `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`, `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`, `specs/exact-prospective-landing/tickets/normalize-review-repair-ticket-metadata.md`, `specs/exact-prospective-landing/tickets/split-review-repair-mutation-rows.md`
Integration surfaces: producer adapter-consumer declarations -> `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`; executable-mode regression owner -> `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`; verification-only repair records -> removal of `specs/exact-prospective-landing/tickets/normalize-review-repair-ticket-metadata.md` and `specs/exact-prospective-landing/tickets/split-review-repair-mutation-rows.md`
Contracts: the four producer tickets' declared adapter-consumer basenames cross into `specs/exact-prospective-landing/tickets/close-adapter-blocker-metadata.md`, asserted by SR1 against those real producer records; the `core.filemode=false` mode-identity behavior crosses the existing landing implementation into `specs/exact-prospective-landing/tickets/preserve-executable-spec-mode.md`, asserted by SR2 against `TestLandPublishesExecutableSpecModeAndReconcilesClean`

## What to build

Fold the useful factual corrections into the two substantive repair tickets. Keep
only the blocker-set behavior in the adapter metadata ticket, and bind the
executable-mode ticket's MR3 to its real `core.filemode=false` landing test. Delete
the two verification-only meta-ticket files: the tree has no lifecycle or
conformance owner for their claimed review/shape probes, so retaining them would
preserve circular or invented evidence. Do not add a new harness or gate.

## Acceptance

- [ ] [SR1] `close-adapter-blocker-metadata.md` carries one behavioral BR1 row: the adopt ticket lists exactly the four producer basenames named by the real producer tickets' adapter-consumer surfaces, with the producer declarations as the independent mutation owner.
- [ ] [SR2] `preserve-executable-spec-mode.md` describes MR3 as the `core.filemode=false` commit/index mode-identity behavior and names `TestLandPublishesExecutableSpecModeAndReconcilesClean` plus its exact focused `go test` invocation as the independent mutation owner and public red operation.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SR1 | omit or corrupt one producer basename in the adopt ticket's `Blocked by:` list | the four producer tickets' `Integration surfaces:` declarations | read the four producer declarations and compare their consumer basename set with `adopt-exact-landing-in-commit.md`; expect the sets to differ |
| SR2 | omit the explicit index mode correction and rely on `git add` under `core.filemode=false` | `TestLandPublishesExecutableSpecModeAndReconcilesClean` | run `go test ./internal/landing -run '^TestLandPublishesExecutableSpecModeAndReconcilesClean$' -count=1`; expect the reconciled index mode mismatch |
