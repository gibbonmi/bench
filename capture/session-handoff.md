# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2c398ee5`; the rescoped `specs/axi-coherent-diff/spec.md` is uncommitted pending reviewer sign-off
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged, rescoped), `specs/axi-query-disclosure/spec.md` (Status: staged, still carries the retired spec-build prerequisite — its own rescope is next in FT173), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: not run this session; only spec prose changed

## State

FT173's `axi-coherent-diff` rescope is written: the retired spec-build/`axi.Action`
prerequisite is removed, the behavior-first ruling preserved, and the slice now
introduces the typed-action owner and `help[]` renderer itself (new `internal/axi`,
CD8, one tracer with CD7). A reviewer-mandated codex `gpt-5.6-sol`/xhigh
falsification loop accepted on iteration 4 after driving three repair rounds:
an old-to-new compatibility row (CD9), owner-bypass proofing (no-own-porcelain
probe, production-path help fixtures), the `--commit` help contradiction, exact
schemas for every public block, divergence in the revision row, and
default-ref/recorded-base identity dimensions with per-attempt resolution
semantics. `bench coverage --check` is green at 9 rows. Reviewer sign-off on the
rescoped spec is the open hard stop; the decision map
`decisions/byte-preserving-axi-foundation/ft173-axi-contract.md` stays shared
top-level because two staged specs consume it.

## Next command

Reviewer sign-off, then commit the rescope and run `/bench-implement-spec specs/axi-coherent-diff/spec.md` in a fresh mid-tier session.

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
