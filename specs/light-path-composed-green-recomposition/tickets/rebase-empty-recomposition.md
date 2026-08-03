# Rebase empty recomposition

Blocked by: accept-composed-green-at-lifecycle-entry.md
Ownership fence: `internal/specbuild`
Assumptions: zero integrations make CandidateTip equal Base; recomposition bootstrap still authorizes the moved subject tip; one-integration replay remains the positive control; claims re-derived from the tree after the blocker lands

## What to build

Promotion recomposition with zero integrated checkpoints rebases the active run directly onto the moved subject tip without sending an empty patch to Git, while non-empty checkpoint replay remains unchanged.

## Acceptance

- [ ] [ER1] A run with zero integrations recomposes successfully after an unrelated subject-tip commit and records Base and CandidateTip as that new tip.
- [ ] [ER2] Empty-boundary recomposition still enforces gate bootstrap before mutating durable run state.
- [ ] [ER3] Recomposition with one integration continues to replay the checkpoint patch onto the moved tip.
- [ ] [ER4] The empty-patch `git apply --cached` error is unreachable from promotion recomposition.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| ER1 | always replay Base-to-CandidateTip even when both commits match | zero-integration promotion regression test | start a run, move the subject tip, promote, run `go test ./internal/specbuild`, expect `No valid patches in input` instead of success |
| ER2 | update Base and CandidateTip before bootstrap succeeds | bootstrap-failure recomposition test | inject a rejecting gate, promote the zero-integration run, reload status, expect the original base and candidate to remain |
| ER3 | treat every recomposition as empty and discard a real integration | existing one-integration recomposition control | integrate one checkpoint, move the tip, promote, run `go test ./internal/specbuild`, expect the integrated content to survive |
| ER4 | call replayCheckpoint with an empty patch | zero-integration promotion regression test | drive the public promotion path and assert success, with the Git empty-patch error absent |
