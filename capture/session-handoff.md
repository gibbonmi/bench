# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — continuation baseline `ee9c4c10`; the staged-spec commit follows it
Spec: `specs/inherited-toolchain-environment/spec.md` — Status: staged, Terra/medium review accepted in 3 iterations
Gate: the atomic staged-spec commit owns the full development-gate result

## State

FT242 is specified as two independently green tickets on one frontier. The
gate slice makes a Go module fail closed when the built-in phase table cannot
resolve Go. The SessionStart slice adds bounded clean-login discovery,
PATH-preserving recovery, descendant teardown evidence, and a manual real
Codex-client/CLI WSL evidence gate before any portability claim. Terra/medium
accepted the combined spec-and-ticket review after three iterations; the
required review learning is open for the next `/bench-drain`.

## Next command

`$bench-implement-spec`

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
