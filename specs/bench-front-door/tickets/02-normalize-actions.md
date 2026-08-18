# Normalize every board action to a command or empty

Blocked by: 01-extract-route-owner.md
Writes: internal/status, internal/dashboard

## What to build

Apply the spec's normalization table to every producer in the status package: gate
states, git precedence (dirty → `/bench-final-check`, else unpushed or unique branch →
`git push`, unavailable → `git status`), intent → `bench status --all`, worktree states
(`bench worktree list`, `git worktree list`, `bench worktree clean <path>`; typed
failure text moves into the detail cell), structure → `bench structure`, decisions
scan failure → `bench maps`, reviews orphan → empty, locked-pending gate → empty. The
drain and roadmap rows keep `/bench-what-next` until ticket 05. Add one table-driven
test that walks every producible row and asserts the owner's predicate: invocable or
empty. Dashboard renders an empty action as an empty cell.

Covers: R24, R26, R27, R28, R29, R30, R31, R32

## Acceptance

- [ ] The all-rows table test is green, and deliberately leaving one prose action makes exactly that case red (record the probe in the commit message).
- [ ] Gate fixtures red/timeout/invalid/unavailable/drifted/exact-tip-partial/locked-pending render `/bench-debug`, `bench gate --fresh`, `bench gate`, `bench gate --fresh`, `bench gate`, `bench gate --fresh`, empty.
- [ ] Dirty tree renders `/bench-final-check` with a detail listing every clause; clean+unpushed renders `git push`.
- [ ] Dashboard test with an empty-action row renders `<td></td>` and exits 0.
- [ ] Gate green.
