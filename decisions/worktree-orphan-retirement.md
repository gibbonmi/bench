# Worktree orphan retirement (FT148)

Status: ready

## Destination

A worktree cut by a session that later dies must stop being immortal. Today
nothing retires it: `bench worktree release` matches only the exact plaintext
request string that created the assignment, the ledger stores a one-way digest
of that string, and the harness hook derives it from the session id. So once
the creating session is gone its worktrees are structurally unreleasable. The
pool accreted from 2026-07-09 to 2026-07-27, and every entry was re-preserved
at every resume sweep. Draining it by hand took a staged script and a full
session.

## Provenance

This map was written in the same session as the spec it compiles — the
highest-bias path in this workflow. Every ticket uses the canonical `Grill`
type; its decision-specific provenance is:

- **#1 and #4:** Closed by reviewer, 2026-07-27 (roadmap row), signed off
  before this session and recorded in `ROADMAP.md`'s FT148 row.
- **#2:** Closed by reviewer, 2026-07-27 (spec-authoring session), replacing
  a rejected first answer.
- **#3:** Closed by reviewer, 2026-07-27 (spec-authoring session).
- **#5:** Decided by the author, then put to the reviewer and **approved
  2026-07-27**; the roadmap row posed it but did not decide it. The spec flags
  this state-destroying behavior for veto.
- **#6:** Closed by reviewer, 2026-07-27 (spec-authoring session). This scope
  addition goes beyond the roadmap row's split and is flagged in the spec for
  veto.
- **#7:** Closed by reviewer, 2026-07-27 (roadmap row) for content; its seam
  was decided by the author.

A mid-tier falsification pass on the first draft found three faults. The
original #2 (a lease conjunct) was unimplementable, #5 carried a sign-off the
roadmap row does not give, and several assertables were unobservable. Those
findings were verified against the tree and are folded in below.

## Notes

## Decisions so far

- [Which command retires an orphan?](worktree-orphan-retirement/tickets/1.md): `bench worktree clean`.
- [What makes an assignment orphaned?](worktree-orphan-retirement/tickets/2.md): **Age alone**: the assignment is in state `active`, and it is older than `bounds.AssignmentStale`, a fixed **7-day** constant beside `LeaseStale`.
- [How does an unstamped record age?](worktree-orphan-retirement/tickets/3.md): Absent = aged.
- [Does the resume sweep clean orphans, or only report them?](worktree-orphan-retirement/tickets/4.md): Report only.
- [What happens to a ledger row whose tree is already gone?](worktree-orphan-retirement/tickets/5.md): Extend the same sweep: an `active` record that is orphaned by #2.
- [The preserved wall stays after this build — what does the reviewer see?](worktree-orphan-retirement/tickets/6.md): Bound it.
- [Where does the prose half land, and how does the gate see it?](worktree-orphan-retirement/tickets/7.md): The `require(<file>, <phrase>)` registry in `checkWorkflowAnchors` (`internal/conformance/docs_workflow_helpers_test.go`).

## Not yet specified

## Spec-writer discretion

## Out of scope

- Draining recovered rows and their recovery references; FT98 owns that retained payload work.
- Any gate, release-path, or worktree behavior outside FT148's orphan retirement and reported cleanup scope.

## Sources
