# Bench handoff

This is a pickup note for the next kit-development session. Durable product facts
live in `README.md`, `AGENTS.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`,
`CONTEXT.md`, and `projects/benchkit.md`; do not treat this file as a second
inventory of commands, skills, hooks, or CLI surfaces.

## Current state

Active phase: `$bench-integrate-learnings`. The work is paused after the reviewer
approved and the session applied the line-governance promotion set. The working
tree is intentionally dirty and uncommitted.

The gate dogfood loop is still pending because a separate session was diagnosing a
known gate problem when this pass started. Do not claim synthesis complete until the
gate question is settled and the appropriate final check runs.

## Current diff

Promoted in the dirty tree:

- The line declaration now uses model + effort + iteration cap. Token estimates are
  only optional sizing notes, not the stop condition.
- `craft-line` keeps venue routing right-sized: delegate when
  isolation/parallelism earns its cost, stay inline for tiny slices or atomic
  diffs, and flag collapsed per-story lines in the exit report.
- `/bench-write-spec` keeps mid-tier spec authoring as the default, with a top-tier
  exception only when the Handoff carries uncertainty flags and the reviewer
  approves.
- `craft-delegate` prefers a map Handoff plus line-ranged excerpts over whole-file
  read lists when charging a delegate.
- `/bench-debug` allows direct fix-and-gate for small single-seam fixes instead of
  forcing `bench shift`.

Recorded/pruned in the dirty tree:

- `CHANGELOG.md` has the 2026-07-04 line-governance learnings entry.
- `.bench/learnings.md` was pruned from 14 open entries to 8.
- `ROADMAP.md` no longer carries the iteration-cap item.
- `ASSESSMENT.md` marks line governance and F3 closed, and recommends continuing
  with the remaining learnings clusters.

Touched files at pause:

```
.agents/commands/bench-debug.md
.agents/commands/bench-implement-spec.md
.agents/commands/bench-write-spec.md
.agents/skills/bench-craft-delegate/SKILL.md
.agents/skills/bench-craft-line/SKILL.md
.bench/BENCH.md
.bench/learnings.md
ASSESSMENT.md
CHANGELOG.md
HANDOFF.md
ROADMAP.md
projects/benchkit.md
```

## Verification already run

- `git diff --check`
- `BENCH_CONFORMANCE_ROOT=/home/devuser/workspace/bench go test -count=1 ./internal/conformance -run '^TestRootConformance$'`
- `go test -count=1 ./...`
- Stale-wording grep over edited guidance surfaces found no live `token cap`,
  mandatory-delegation, or old top-tier-ban wording outside historical changelog.

Not run: `bench gate`, for the gate-diagnosis reason above.

## Remaining learnings

`bench learnings` currently reports 8 open entries:

- Session-start stale gate: split benign drift from real drift.
- Skill/command dogfood gotcha: trigger edits take effect in the next session.
- Review implementation: live-session verification of returned findings.
- Linked-repo by-path CLI edge inventory.
- Review findings persistence.
- External-format/library decision: official library vs private dialect.
- Runnable probes for byte/wire compatibility claims.
- Codex/no-subagent review-axis fallback.

Recommended next decision: review-findings persistence, because it affects whether
the rest of `/bench-integrate-learnings` has a durable pickup surface.
