# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — last landed commit `135e55f3`; ticket 4's prose-repair diff is
in the working tree awaiting the coordinator's landing commit
Spec: `specs/remove-spec-build-lifecycle/spec.md` (Spec A) — implemented
pending review; 4/4 tickets landed

## State

Spec A (lifecycle removal) is built: the `bench spec build` grammar and
`bench worktree recovery` are gone from the binary, `bench resume`'s reconcile
deletes the two lifecycle ref namespaces and purges lifecycle ledger
assignments, `bench commit --spec` is the sole staged→implemented author with
its semantics in `--help`, the ticket parser is deleted (`Blocked by:` stays a
documented convention), and every kit-prose reference is repaired with a new
standing removed-verb sweep in `internal/conformance` that was observed red
before the repair. The reviewer authorized direct-to-main commits for this
program. Closed decisions stay closed: full lifecycle removal with zero
backwards compatibility, serial commit-on-green as the spec-backed cadence,
only `Blocked by:` machine-read in tickets. The doctrine rewrite (skill diets,
`bench-craft-domain`) is Spec C, not this build.

Flagged for reviewer veto:

- The removed-verb sweep exempts `CHANGELOG.md`, deviating from RM8's literal
  file list — it is append-only history and scrubbing it would falsify the
  record; the changelog carries the one sanctioned removal entry.
- `bench worktree clean --apply` still preserves dirty payloads into
  `refs/bench/recovery/`, a namespace resume now sweeps; recommend a
  refuse-dirty follow-up so clean never writes into a swept namespace.
- RM11's keep-list predicate was corrected during the build (the enumerated
  kept-verb `--help` guard).

## Next command

`/bench-review-implementation` for `specs/remove-spec-build-lifecycle/spec.md`
over the composed Spec A diff.

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
