# Repair the FT78 symlink proof's prospective fixture

Blocked by: adopt-exact-landing-in-commit.md
Ownership fence: `internal/gate/story3_proof_helpers_test.go`
Integration surfaces: hostile symlink proof setup→`internal/gate/story3_proof_helpers_test.go`; path-scoped `bench commit` prospective authorization→the fixture's tracked gate and manifest
Contracts: committed gate and manifest bytes, the ignored declared target, and the common-Git-dir run marker cross fixture setup→the real compiled `bench commit` invocation in `internal/gate/story3_proof_helpers_test.go`, asserted by PF1-PF2 against the detached prospective checkout and repository marker

## What to build

Repair only the retained FT78 symlink-chain fixture. Before its first green seed,
make the executable gate and gate-input manifest part of the fixture's committed
prospective subject while leaving `work.txt` as the named uncommitted landing delta
and leaving `inputs/` ignored. Make the fixture gate append its run marker through
`git rev-parse --git-common-dir`, because `.git` is a file in the detached worktree
used for prospective authorization.

Do not change production gate, subject, reuse, landing, or commit behavior. Do not
weaken the existing oracle-change, two-run, red-verdict, or unchanged-HEAD assertions.

## Acceptance

- [ ] [PF1] `TestFT78Story3ProofLedger/R8/symlink-chain` reaches the prospective gate after the ignored resolved target changes, records one additional red run in the common Git directory, changes oracle identity, and leaves HEAD unchanged.
- [ ] [PF2] The fixture's prospective subject contains its executable gate and manifest while the named `work.txt` delta and ignored symlink chain retain their original attribution roles.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PF1 | leave the fixture gate or manifest outside the committed prospective subject | existing `TestFT78Story3ProofLedger/R8/symlink-chain` | seed green, mutate the ignored resolved target, run `bench commit work.txt`, and observe infrastructure refusal before the expected refreshed red verdict |
| PF2 | write the run marker through the detached checkout's `.git` path instead of the common Git directory | existing `TestFT78Story3ProofLedger/R8/symlink-chain` | execute the prospective gate and observe the marker write fail or the required second run remain uncounted |
