# Repair post-review evidence closure

Blocked by: repair-registry-ast-single-source.md, repair-disclosure-value-proof-seam.md, repair-guards-enumeration-compatibility.md, repair-worktree-review-closure.md
Writes: `reviews/axi-query-disclosure.md`, `capture/session-handoff.md`

## What to build

Verify every committed review finding against the landed repair commits, delete the transient pickup only when every finding is closed, and regenerate the phase handoff from live repository state.

## Acceptance

- [ ] [PE1] (covers local) every finding in `reviews/axi-query-disclosure.md` names one landed repair and has a focused green check plus its required observed-red mutation evidence.
- [ ] [PE2] (covers local) the pickup file is deleted only after the four repair tickets are independently green on main; no resolved review artifact remains.
- [ ] [PE3] (covers local) `capture/session-handoff.md` is regenerated with `bench handoff`, pins the resulting candidate and review-ready state, and names `$bench-review-implementation specs/axi-query-disclosure/spec.md --full` as the exact next command.
