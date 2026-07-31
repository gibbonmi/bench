# Teach and enforce spec-build cadence

Blocked by: Expose spec-build porcelain

Ownership fence: `.agents/commands/bench-implement-spec.md`, `.agents/commands/bench-review-implementation.md`, `.agents/commands/bench-final-check.md`, `.agents/skills/bench-craft-delegate/SKILL.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`, `bin/bench.sh`, `projects/benchkit.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/registry/registry.go`, `internal/contract/runtime/runtime_spec_build_test.go`, `tests/canary/workflow-guidance-anchors`, `CHANGELOG.md`
Assumptions: the porcelain commit that supplies all eight public commands and their runtime contracts is an ancestor of the current clean assignment base; `docs_workflow_helpers_test.go` owns the workflow anchors and the registry already maps `workflow-guidance-anchors` to its conformance owner
Line: approved top-tier implementation on `gpt-5.6-sol` at high effort

## What to build

Update the canonical implementation command, delegation skill, shared workflow
and inventories, shell help, project profile, changelog, and conformance owners
as one consistent leverage cut. The workflow must continuously fill every
ownership-safe frontier slot through public porcelain, submit focused evidence
plus a coordinator-owned probe, review the exact composition before promotion,
and route accepted repairs through the same provisional lifecycle. The review
must bind the exact candidate, `promote` must be the sole final gate, status, and
commit author, and terminal final-check must only report retained evidence and
write the retro. Light-path work, `bench shift`, and ordinary `bench commit` keep
commit-on-green behavior.

Before editing, enumerate the live harness capacity and the complete ready
frontier. For every unused slot, retain one structured reason from this closed
set: dependency, ownership-fence overlap, unavailable capacity, or measured
resource constraint. Recompute the frontier after every integration or release;
do not rely on ticket-time readiness claims.

## Surface and mutation matrix

| Owner | Fact it must carry | Independent red mutation |
|---|---|---|
| `.agents/commands/bench-implement-spec.md` | full runs start/refill the safe frontier, use public porcelain, review the exact composition, and route accepted repairs through a new provisional ticket | replace refill-after-integration with drain-then-refill, or route a repair directly to the working branch |
| `.agents/commands/bench-review-implementation.md` | active spec builds bind review to the exact candidate and route accepted findings through provisional repair before promotion | permit a direct working-branch repair or an unbound review |
| `.agents/commands/bench-final-check.md` | a promoted spec build reports retained evidence and writes the retro without authoring another gate, status transition, or commit | restore `bench commit --spec` or `bench spec implemented` after promote |
| `.agents/skills/bench-craft-delegate/SKILL.md` | delegates receive owned worktrees, focused evidence, and a coordinator probe whose mutation kind differs from the author probe | swap the probe rule for a second instance of the author's mutation kind |
| `.bench/BENCH.md` | full spec builds use provisional cadence while light-path work, `bench shift`, and ordinary `bench commit` remain commit-on-green | broaden provisional cadence to an ordinary commit path |
| `.bench/BENCH-reference.md` | the lifecycle lookup names all eight public operations and their routing purpose without duplicating the canonical workflow | delete or swap one operation-to-purpose entry |
| `bin/bench.sh` | help exposes copy-paste grammar for all eight operations, every evidence flag, and abandon plan/apply | move a flag value into the positional slot or delete one grammar |
| `projects/benchkit.md` | the Bench-kit profile names the approved line, structure preflight, dogfood traces, and final composed-gate boundary | drop the structure preflight or replace the approved line |
| `CHANGELOG.md` | the user-visible entry distinguishes provisional full-spec cadence from unchanged light/shift/ordinary-commit cadence | delete the unchanged-path control |
| `internal/conformance` | existing workflow owners check every fact above and keep one source per fact | point a check at a copied fixture instead of the canonical surface |
| `internal/contract/runtime/runtime_spec_build_test.go` | runtime dogfood covers the three-to-four-ticket refill and retained red/repair/green evidence | remove the fourth-ticket refill while leaving the three-ticket trace green |
| `tests/canary/workflow-guidance-anchors` | each owned check bites through a deletion, swap, routing, or additive-contradiction mutation, not merely another token deletion | add permission for a direct repair while leaving every required token present |

## Dogfood and regression matrix

- A three-ticket trace fills three ownership-safe slots concurrently, integrates one ticket, and records why every remaining slot is or is not fillable.
- A four-ticket trace integrates one frontier ticket and dispatches the newly eligible fourth ticket while another delegate remains active.
- An accepted review finding becomes a new repair ticket and traverses assign, checkpoint, integrate, fresh composed review, and promote; no direct working-branch repair is permitted.
- Static and runtime checks reject lifecycle Git plumbing synthesized outside the eight public operations.
- Existing shift and ordinary-commit gate-call contracts run unchanged as positive controls.

## Structure preflight

`internal/conformance` is already structurally wide and the repository currently
has 52 inherited structure findings. At assignment entry,
`docs_workflow_helpers_test.go` is 758/660 lines,
`docs_workflow_checks_test.go` is 425/400 lines, and `fixture_bite_test.go` is
573/400 lines. Extend the named existing conformance files and the existing
`workflow-guidance-anchors` family, but first compact table-driven helpers so no
named owner grows beyond its entry size and the total structure finding count
does not increase. Do not add another Go file or new canary family without a
`craft-seams` split or explicit grant.

## Acceptance

- [ ] [R41-R43] Initial and newly eligible frontiers fill harness capacity; every unused slot names dependency, fence overlap, unavailable capacity, or a measured resource constraint.
- [ ] [R44-R45] Harnesses use all eight public commands, never synthesize Git lifecycle plumbing, and review the exact composition before promote with repairs re-entering the provisional path.
- [ ] [R46-R50] Each enumerated command, skill, shared guide/reference, help, and profile surface independently carries its owned lifecycle or routing facts and its conformance mutation.
- [ ] [R51] Existing shift and ordinary commit gate-call contracts remain green as regression controls.
- [ ] A three-ticket dogfood trace proves concurrent frontier dispatch, and a four-ticket trace proves newly eligible dispatch while another delegate remains active.
- [ ] Review binds the exact candidate; `promote` is the sole final gate, status, and commit author; terminal final-check only reports retained evidence and writes the retro.
- [ ] A coordinator-authored additive contradiction fails before repair and passes after repair, using a mutation kind different from the implementor's deletion or replacement probe.
