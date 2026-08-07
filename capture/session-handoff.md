# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `e1c24b8`, 4 dirty paths, 19 unpushed commits
Spec: none staged.
Gate: green at `f4fd683` — stale, work tree `ea46b76`

## State

**Phase reached: pre-push-guard-visibility promoted terminal; capture drain and spec retirement are next.**

Run `9975d21b` promoted on green: candidate `a57e2334`, promotion commit `e1c24b83`, retained evidence `v1:4cc9bdc6…`, `Status: implemented`. Nineteen assignments across three review rounds; every accepted finding closed, including the spec-authored clean-skip defect (predicate amended to prospective bytes-and-mode match, amendment landed in the closing commit after promotion because the run pins staged spec content). The closing commit also carries the dead `manifestOwnedClean` deletion, three learnings entries (gate-free repair bookkeeping, SpecTip pin, ticket-fence enforcement), and `capture/retros/pre-push-guard-visibility.md`.

Open items for the drain: the retro plus the three learnings entries; two contestable unpinned branches recorded in the round-three review receipt (permission comparison in `convergedFingerprint`, skip-path manifest-row write) as follow-up ticket candidates; `internal/contract/surface/link_lifecycle_test.go` is over the 400-line structure budget without a grant — reviewer's grant or split. Spec retirement (`bench spec retire pre-push-guard-visibility`) was deferred from this close so the reviewer can decide what durable content to promote first. A parked light-path task exists: add the repair-ticket rationale prose to `.agents/commands/bench-implement-spec.md` (prompt already delivered to the reviewer). The cross-harness falsification pass was explicitly skipped by the reviewer; all required lifecycle reviews and promotion ran.

## Next command

`/bench-what-next`

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
