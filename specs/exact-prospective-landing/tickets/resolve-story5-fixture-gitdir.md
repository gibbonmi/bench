# Resolve Story 5 fixture Git directory

Blocked by: none
Ownership fence: `internal/contract/runtime/runtime_gate_action_proof_test.go`
Integration surfaces: Story 5 gate-script markers→existing `internal/contract/runtime/runtime_gate_owner_helper_test.go` plus RF1; prospective checkout creation→existing `internal/landing` exercised by RF1-RF2; adapter package completion→adopt-exact-landing-in-commit.md
Contracts: the repository-common Git directory crosses the shell gate scripts in `internal/contract/runtime/runtime_gate_action_proof_test.go`→their Go marker readers, asserted by RF1-RF2 against the real prospective checkout and Git process

## What to build

Make every Story 5 gate-action fixture script resolve its repository-common Git directory through Git before reading or writing coordination files, so the same fixture works from an ordinary checkout and from the prospective checkout where `.git` is a file.

## Acceptance

- [ ] [RF1] The interrupted-pending proof creates and reuses its one-shot owner sentinel and marker files in the common Git directory, so the bounded R14 repro completes green from the prospective checkout instead of re-entering the infinite first-run loop.
- [ ] [RF2] The open-green, stale-green, locked-pending, interrupted, and cancellation gate scripts exchange every `story5-*` coordination file with their existing Go readers through the common Git directory, and the full Story 5 proof ledger completes green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RF1 | address the one-shot owner sentinel through `.git/story5-owner-once` again | existing R14 interrupted-pending proof | run `timeout 20s go test -count=1 ./internal/contract/runtime -run 'TestFT78Story5ProofLedger/R14/interrupted-pending-inspection' -v`, expect bounded exit 124 at the interrupted-pending row |
| RF2 | address one remaining gate-script coordination file through `.git/story5-*` instead of the resolved common Git directory | existing Story 5 proof-ledger rows and Go marker readers | run `go test -count=1 ./internal/contract/runtime -run '^TestFT78Story5ProofLedger$'`, expect the owning row to fail or hit its bounded marker deadline |
