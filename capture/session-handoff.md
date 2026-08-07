# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `fd0eb8f`, 1 dirty path, 37 unpushed commits
Spec: `specs/conformance-harness-scope/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `dadc792` — stale, work tree `7b7f77a`

## State

**Phase reached: conformance-harness-scope spec staged and reviewer-approved.**

`specs/conformance-harness-scope/spec.md` is the uncommitted staged spec. Its
reviewer-confirmed 2026-08-07 source authorizes two test-only outcomes: direct
conformance fixture bites run only the check resolved from the executable
registry, and freshness failure classes move to `Verify`/`Check` while
representative shell-entry composition remains. `bench coverage --check`
validates all 13 acceptance rows. The required read-only falsification pass
found a duplicated-registry false green and an overbroad fixture-family
quantifier; both are repaired in the draft.

Implementation has not started and no spec-build run exists. The closed fences
permit only the named conformance fixture-bite test files and the freshness plus
gate-entry test files. Production registry contents, selection semantics,
freshness policy, gate routing, canary fixtures, scheduler work, and the
post-reduction census remain out of scope. The direct-bite quantifier is the ten
families enumerated in the spec; `injected-ports` stays canary-sweep-owned. A
registry-source rebind mutation must prove the helper consumes
`registry.FamilyCheck` rather than a private copy.

`specs/pre-push-guard-visibility/spec.md` remains a separate staged spec and is
not part of this lifecycle. The only uncommitted paths owned by this phase are
the new conformance-harness-scope spec and this rewritten handoff.

## Next command

`$bench-implement-spec --full conformance-harness-scope`

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
