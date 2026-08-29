# Retire the raw-Git route in guidance

Blocked by: add-worktree-merge-verb.md
Writes: internal/worktree/land_refusal.go, internal/worktree/land_surface_test.go, .agents/skills/bench-craft-delegate/SKILL.md, internal/anchors/registry_data.go, internal/conformance/fixture_bite_test.go, .bench/BENCH-reference.md, docs/adr/0014-main-receives-writes-only-through-landings.md, projects/benchkit.md, CONTEXT.md

## What to build

Every place that advertises the raw route or claims no verb exists names the
verb instead. The landing's conflict `next=` keeps `git -C <path> merge
<destination>` as the hand resolution, states that `bench worktree merge`
refuses the same conflict, and drops the parenthetical claim.

The `craft-delegate` stale-base sentence names
`bench worktree merge --from main <target>` for the coordinator, and its
anchor needle and fixture-bite case move with it in the same commit. No
canary fixture carries the sentence. The skill file sits exactly at its line
budget. The replacement is line-neutral, or the same commit raises the
profile's budget row by the exact count.

The reference's repair list, ADR 0014's
consequence, the profile's CLI seam paragraph, and the glossary describe the
verb as the current decided state. The glossary gains `sibling assignment` and
`worktree merge`, each with the Avoid list the other entries carry.

## Acceptance

- [ ] WM27: the landing conflict `next=` names `git -C <path> merge
      <destination>` as the hand resolution, states that `bench worktree merge`
      refuses the same conflict, and no longer claims that no verb exists.
- [ ] WM28: the anchor registry needle and the fixture-bite case carry the
      merge-verb sentence, and the mutation that hands the merge to the
      delegate goes red.
- [ ] The reference, ADR 0014, the profile, and `CONTEXT.md` name the verb,
      and the gate's prose check is green on each edited file.
- [ ] The `guidance-prose-budgets` check is green after the skill edit.
