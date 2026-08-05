# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `e176eb9`, 1 dirty path, 6 unpushed commits
Spec: `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: red at `2392bdf` — current

## State

**Phase reached: implementation remains active; two out-of-fence prerequisites are integrated, the adapter is refreshed and focused-green, and a third fixture prerequisite stopped the full runtime package at the declared cap.**

Run `869a2d33cab96f962882762a9ea56fd21c952976248e64fa917f6da6c48dbca3` has candidate `30a64752f82fcd330dfe217d833d80e33eb9d939`. Integrated repair assignments:

- `62676f7c03ae6c96c024211994b82e3d`: writer-aware prospective authorization and landing output, full scoped packages green.
- `383e2781a83a5743e1d5251148bf4450`: exact green reuse before gate-lock acquisition, full `internal/gate` green in 144.947s.

Adapter assignment `daa1f4183dcb0dd159cca7b738d0fe7c` is active on base `30a64752f82fcd330dfe217d833d80e33eb9d939`. Refresh preserved its three in-fence dirty paths: `internal/commit/commit.go`, `internal/contract/runtime/runtime_commit_test.go`, and `CHANGELOG.md`. AC4's public prospective-checkout inspection and AC7's deterministic two-process CAS/rerun journey are focused-green; the wrong-base mutation is red and restored.

The full `go test -count=1 ./internal/commit ./internal/contract/runtime` timed out only at `TestFT78Story5ProofLedger/R14/interrupted-pending-inspection`. Tight repro:

`timeout 20s go test -count=1 ./internal/contract/runtime -run 'TestFT78Story5ProofLedger/R14/interrupted-pending-inspection' -v`

Exit `124`, combined-output SHA-256 `d15278907b6545b4b1e2f0e60c0f371b23ba5687c35f2fbdd20965e6800d7a0e`, produced `2026-08-05T19:23:42.151084233Z`. Confirmed cause: the fixture gate script addresses `.git/story5-*` as a directory; prospective checkouts carry `.git` as a file, so its one-shot sentinel is never found and the script enters its infinite first-run loop. The repair must own `internal/contract/runtime/runtime_gate_action_proof_test.go` and resolve the common Git directory through Git, matching the existing prospective-safe fixture pattern. The adapter remains resumable afterward.

Closed decisions: use `assign --refresh`, never abandon, for this recoverable run; preserve the original adapter request `7c2e9a1d4f6b4830a5c7e1d9b3f60842`; no review or promotion occurs until the adapter checkpoints and integrates; the next repair gets a top-tier line because the mid-tier adapter missed named public rows.

## Next command

`$bench-implement-spec --full specs/exact-prospective-landing/spec.md`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
