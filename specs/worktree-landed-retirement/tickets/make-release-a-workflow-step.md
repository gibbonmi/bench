# Make release a workflow step

Blocked by: apply-the-landed-plan-and-settle-records.md
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/commands/bench-final-check.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors/, CHANGELOG.md

## What to build

`craft-delegate`'s acceptance sentence in `## Verifying the done-claim` becomes
"Acceptance closes an independent worktree after its slice lands: the coordinator runs
`bench worktree release --request <opaque-id> <path>` for it." `/bench-final-check`'s
post-merge tail in `## Exit handoff` replaces "leftover worktrees and scratch branches go
through `bench worktree clean`" with "scratch branches go through `bench worktree clean`;
leftover worktrees are retired by `bench worktree clean --landed`: run the plan, apply it,
and carry the plan and apply result in the landing report". Three anchors join the
registry in the `AfterImplementSpec` group — `RequireInSection` (craft-delegate,
`Verifying the done-claim`, that needle, diagnostic ".agents/skills/bench-craft-delegate/SKILL.md
dropped release-at-acceptance"); `RequireInSection` (final-check, `Exit handoff`, the
new sentence, diagnostic ".agents/commands/bench-final-check.md post-merge tail dropped
the landed worktree sweep step"); `ForbidInSection` (final-check, `Exit handoff`,
"leftover worktrees and scratch branches go through", diagnostic
".agents/commands/bench-final-check.md still routes leftover worktrees to a bare per-path
clean") — each with one canary fixture (`BASE`, `MUTATE.json`, `EXPECT`) under
`tests/canary/workflow-guidance-anchors/`. Because it lands last, this ticket writes the
one `Added` `CHANGELOG.md` entry under `[Unreleased]` naming the `landed` classification
and `bench worktree clean --landed`. Demo: remove either sentence, gate goes red.

## Acceptance

- [ ] `(covers LR15)` All three tuples are registered exactly once and the live command and skill files satisfy them on the clean tree.
- [ ] `(covers LR34)` Each tuple's canary fixture mutation emits exactly its diagnostic and restoration clears it.
- [ ] `(covers LR16)` `CHANGELOG.md` `[Unreleased]` carries exactly one `Added` entry naming both.
