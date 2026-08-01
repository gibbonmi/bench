# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — 5 unpushed commits
Spec: `specs/reduced-gate-phase-set/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)

## State

- Phase reached: **`/bench-implement-spec --full` — tickets derived, build not yet
  started.** Ten tickets under `specs/reduced-gate-phase-set/tickets/`, blocker graph
  verified by hand. `tickets/README.md` holds the `[R01]`–`[R26]` → (story, behavior)
  map, assigned by hand because `bench coverage` emits no stable row identity and this
  spec has rows sharing a story and a seam string. Next action is
  `bench spec build start reduced-gate-phase-set`, then assign the frontier.
- Reviewer ruling landed this session: **`specs/` joins the path allowlist** alongside
  `capture/`, `ROADMAP.md`, and `.bench-notes.md`. It is all formatted documents, and
  the phases that grade it — conformance and conformance-suite — are the included set,
  so a spec edit stays graded on a reduced run. The spec's enumeration and its
  membership paragraph were edited to record it.
- Unverified by reading, and deliberately so: whether `test`, `contract`, or `canary`
  touch the real `specs/` rather than a fixture tree. Story 4's stripped construction
  is the thing that answers it, and it reds loudly rather than silently if one does.
- One spec ambiguity ruled on and flagged for veto: row R04 says "a nested or
  sibling-prefixed path is not" a member, contradicting R02, which requires every
  co-located capture surface to be covered — `capture/retros/<slug>.md` is nested.
  Read as descendant containment with the path boundary respected, so a sibling
  prefix (`capture-old/x.md`) never matches. Written into ticket 1.
- The build's central design, because a wrong summary sends it the wrong way:
  excludable phases run on a full gate against a materialized stripped worktree *with
  capabilities required*. Stripping alone is not enforcement — `skipIfSubjectFileMissing`
  turns an absent subject file into a capability skip, informational in the dev tier,
  so a check whose file vanished would go permanently green. That was the first draft's
  fatal error.
- Reviewer decisions this session, both closed: ticket 7 (ancestor selection) runs at
  `fable`/high rather than the spec's `opus`; every other ticket keeps its spec line.
  A Codex CLI falsification pass at the top binding runs over the finished diff before
  promotion, charged to refute rather than to grade.
- Seam call the spec left open: the declaration lives in `internal/gate` (exported),
  not a new package — `internal/status` already imports `internal/gate`, so there is no
  cycle, and the coverage rows name `internal/gate` as the unit seam.
- Two open learnings and one parked idea await the next drain; not this build's work.

## Next command

`/bench-implement-spec --full specs/reduced-gate-phase-set/spec.md`

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
