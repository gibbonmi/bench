# Prune duplicated and change-narrating prose

Blocked by: none
Writes: README.md, projects/benchkit.md, .bench/BENCH.md, .bench/BENCH-reference.md, .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-skills/SKILL.md, scripts/gen-platform-packages.sh, scripts/build-artifacts.sh, internal/worktree/reconcile.go, internal/shift/session.go, internal/lines/lines.go, internal/gate/verdict.go, internal/toon/toon_test.go, internal/conformance/package_shipped_surface_test.go, internal/dashboard/dashboard_test.go, capture/learnings.md

## What to build

The README points to the reference for the file map and harness details. It
states the current adapter contract and keeps its public story and provenance.
The review command states authority once. The skill-writing guidance keeps
each pruning rule in one section. Comments state timeless constraints instead
of change history. Stray export tags and broken source wraps leave the root
documents.

## Acceptance

- [ ] The README contains no second directory tree and states that adapters read the prompt from stdin.
- [ ] The review command states that the gate decides done and does not repeat that boundary at exit.
- [ ] `craft-skills` defines duplication, no-op prose, and sediment only under `Pruning`.
- [ ] The nine approved comments describe the current code without change narration.
- [ ] No authored document ends with a stray `</content>` tag.
- [ ] The prose-only batch passes the project gate through `bench commit`.
