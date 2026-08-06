# Preserve executable spec mode

Blocked by: build-exact-landing-owner.md
Ownership fence: `internal/landing/landing.go`, `internal/landing/landing_test.go`
Integration surfaces: transitioned-entry mode derivation and reconciliation→`internal/landing/landing.go`; executable and non-executable spec-transition coverage→`internal/landing/landing_test.go`; implemented-bytes producer→existing `internal/spec/spec.go` plus MR1; commit and index mode publication→existing `internal/git` Git invocation surface plus MR1
Contracts: the transitioned spec's implemented bytes and derived Git regular-file mode (`100644` or `100755`) cross `internal/spec`→`internal/landing/landing.go`, asserted by MR1 and MR2 against the real staged-spec transformer and real Git commit, index, and worktree observations

Accepted exact-candidate review repair for finding P1. Prospective composition
currently hardcodes a transitioned spec entry to Git mode `100644`, while
reconciliation preserves the invoking checkout's executable bit before staging.
An executable staged spec can therefore publish one mode and leave the real index
at another.

## What to build

Derive the transitioned spec entry's Git mode from the attributed filesystem entry,
using Git's regular-file modes (`100644` or `100755`) rather than carrying arbitrary
permission bits. Keep the transformed bytes and selected mode identical across the
authorized tree, published commit, and reconciled index/worktree.

## Acceptance

- [ ] [MR1] A staged executable spec is composed and published as `100755`; authorization observes implemented bytes, and the invoking index/worktree is clean after reconciliation with the executable bit preserved.
- [ ] [MR2] A staged non-executable spec remains `100644`; its narrower filesystem permissions remain preserved in the worktree, and existing red/CAS-loss behavior remains unchanged.
- [ ] [MR3] Composition and reconciliation use the same single mode-derivation helper: with `core.filemode=false` set in the invoking repository, the published commit mode and the reconciled index mode still match the derived mode for both the `100644` and `100755` transitions.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MR1 | hardcode the transitioned entry to `100644` again | executable landing spec-transition test | stage an executable spec, land it, then assert commit/index mode identity and clean reconciliation; expect mode drift |
| MR2 | force the transitioned entry to `100755` unconditionally | non-executable landing spec-transition test | stage a non-executable spec, land it, then assert `100644` in commit and index with narrower filesystem permissions preserved; expect mode drift |
| MR3 | omit the explicit `update-index --cacheinfo` mode correction and rely on `git add` under `core.filemode=false` | `TestLandPublishesExecutableSpecModeAndReconcilesClean` | `go test ./internal/landing -run '^TestLandPublishesExecutableSpecModeAndReconcilesClean$' -count=1` |
