# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `9674d6ec`; the rescoped `specs/axi-query-disclosure/spec.md` and this handoff refresh are uncommitted pending reviewer sign-off
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged, rescope landed at `00583fa9`), `specs/axi-query-disclosure/spec.md` (Status: staged, rescoped, sign-off pending), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `9674d6ec` via path-scoped `bench commit`; commits unpushed

## State

Both FT173 rescopes are done. `axi-coherent-diff` (approved, landed `00583fa9`)
introduces the typed-action owner and `help[N]{cmd,why}` renderer in new
`internal/axi` and makes `bench diff` the coherent snapshot. `axi-query-disclosure`
is rescoped on top of it: the spec-build family is out of the approved set (six
root queries plus `worktree list`), the harness-phase action kind lands there as
its only additive `internal/axi` change, and its codex `gpt-5.6-sol`/xhigh
falsification loop accepted on iteration 3 after two repair rounds (QD6
additive-compat paired fixtures incl. `anchors`, registry-declared AXI
membership on the production entries in `cmd/bench/main.go`, seven exact
conformance predicates, a guidance-versus-registry cross-check, both inventoried
worktree actions, blocked ticket order). `bench coverage --check` is green at
6 rows; reviewer sign-off on that rescope is the open hard stop. The decision
map `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md` stays
shared top-level because two staged specs consume it.

## Next command

Reviewer sign-off on `specs/axi-query-disclosure/spec.md`, then commit it and run `/bench-implement-spec specs/axi-coherent-diff/spec.md` in a fresh mid-tier session.

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
