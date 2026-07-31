# Teach and enforce spec-build cadence

Blocked by: Expose spec-build porcelain

Ownership fence: `.agents/commands/bench-implement-spec.md`, `.agents/skills/bench-craft-delegate/SKILL.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`, `bin/bench.sh`, `projects/benchkit.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/registry/registry.go`, `tests/canary/workflow-guidance-anchors`, `CHANGELOG.md`
Assumptions: all eight porcelain commands and their runtime contracts are green at `925d64e`; `docs_workflow_helpers_test.go` owns the workflow anchors and the registry already maps `workflow-guidance-anchors` to its conformance owner
Line: approved top-tier implementation on `gpt-5.6-sol` at high effort

## What to build

Update the canonical implementation command, delegation skill, shared workflow
and inventories, shell help, project profile, changelog, and conformance owners
as one consistent leverage cut. The workflow must continuously fill every
ownership-safe frontier slot through public porcelain, submit focused evidence
plus a coordinator-owned probe, review the exact composition before promotion,
and route accepted repairs through the same provisional lifecycle. Light-path
work, `bench shift`, and ordinary `bench commit` keep commit-on-green behavior.

Before editing, enumerate the live harness capacity and the complete ready
frontier. For every unused slot, retain one structured reason from this closed
set: dependency, ownership-fence overlap, unavailable capacity, or measured
resource constraint. Recompute the frontier after every integration or release;
do not rely on ticket-time readiness claims.

## Surface and mutation matrix

| Owner | Fact it must carry | Independent red mutation |
|---|---|---|
| `.agents/commands/bench-implement-spec.md` | full runs start/refill the safe frontier, use public porcelain, review the exact composition, and route accepted repairs through a new provisional ticket | replace refill-after-integration with drain-then-refill, or route a repair directly to the working branch |
| `.agents/skills/bench-craft-delegate/SKILL.md` | delegates receive owned worktrees, focused evidence, and a coordinator probe whose mutation kind differs from the author probe | swap the probe rule for a second instance of the author's mutation kind |
| `.bench/BENCH.md` | full spec builds use provisional cadence while light-path work, `bench shift`, and ordinary `bench commit` remain commit-on-green | broaden provisional cadence to an ordinary commit path |
| `.bench/BENCH-reference.md` | the lifecycle lookup names all eight public operations and their routing purpose without duplicating the canonical workflow | delete or swap one operation-to-purpose entry |
| `bin/bench.sh` | help exposes copy-paste grammar for all eight operations, every evidence flag, and abandon plan/apply | move a flag value into the positional slot or delete one grammar |
| `projects/benchkit.md` | the Bench-kit profile names the approved line, structure preflight, dogfood traces, and final composed-gate boundary | drop the structure preflight or replace the approved line |
| `CHANGELOG.md` | the user-visible entry distinguishes provisional full-spec cadence from unchanged light/shift/ordinary-commit cadence | delete the unchanged-path control |
| `internal/conformance` | existing workflow owners check every fact above and keep one source per fact | point a check at a copied fixture instead of the canonical surface |
| `tests/canary/workflow-guidance-anchors` | each owned check bites through a deletion, swap, or routing mutation, not merely another token deletion | swap frontier order, probe kind, or repair route while leaving all tokens present |

## Dogfood and regression matrix

- A three-ticket trace fills three ownership-safe slots concurrently, integrates one ticket, and records why every remaining slot is or is not fillable.
- A four-ticket trace integrates one frontier ticket and dispatches the newly eligible fourth ticket while another delegate remains active.
- An accepted review finding becomes a new repair ticket and traverses assign, checkpoint, integrate, fresh composed review, and promote; no direct working-branch repair is permitted.
- Static and runtime checks reject lifecycle Git plumbing synthesized outside the eight public operations.
- Existing shift and ordinary-commit gate-call contracts run unchanged as positive controls.

## Structure preflight

`internal/conformance` is already structurally wide and the repository currently
has 52 inherited structure findings. Extend the named existing conformance files
and the existing `workflow-guidance-anchors` family; do not add another Go file or
new canary family without a `craft-seams` split or explicit grant. The ticket must
finish with no new structure finding.

## Acceptance

- [ ] [R41-R43] Initial and newly eligible frontiers fill harness capacity; every unused slot names dependency, fence overlap, unavailable capacity, or a measured resource constraint.
- [ ] [R44-R45] Harnesses use all eight public commands, never synthesize Git lifecycle plumbing, and review the exact composition before promote with repairs re-entering the provisional path.
- [ ] [R46-R50] Each enumerated command, skill, shared guide/reference, help, and profile surface independently carries its owned lifecycle or routing facts and its conformance mutation.
- [ ] [R51] Existing shift and ordinary commit gate-call contracts remain green as regression controls.
- [ ] A three-ticket dogfood trace proves concurrent frontier dispatch, and a four-ticket trace proves newly eligible dispatch while another delegate remains active.
