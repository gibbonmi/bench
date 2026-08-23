# The default branch receives writes only through landings, and a landing composes by merge

Every Bench phase runs in a bench worktree and lands through the landing verb. A drain, a decision map, a bug fix, a light-path ticket, and a spec-backed build all land the same way. The landing verb therefore accepts a source with no spec. Without a spec it skips only the spec transition and the ownership-fence check. It still proves the destination, the assignment, the identity, the reviewed range, and the gate. A spec that names a tickets-only folder closes that folder on the landing, so a light-path ticket retires through the verb.

The landing composes the reviewed source onto the moved default branch by merge. Rebase is rejected, because it rewrites the reviewed tip and breaks the identity the review graded. The published commit keeps the destination as its first parent and the reviewed source tip as its second.

A conflicted phase-owned capture file composes by one rule table with three verbs. The session handoff takes the source side, because the closing session's state is the handoff. The two append-only journals take the union of both sides, so no appended entry is lost. Those journals are the learnings journal and the ideas inbox. Every other capture file takes the destination side.

A path the table does not name, a board file, or a non-regular object kind refuses. The refusal names every conflicted path and the source repair in order. That repair is: merge the destination into the source worktree, commit, review the new range, and re-run the landing.

The rule is guidance, not enforcement. No hook refuses a commit on the default branch, and the ordinary commit verb keeps working on any branch. The handoff is written in the phase's worktree and lands with it, so the handoff rule and the merge rule agree.

## Considered options

- **Rebase the source onto the moved default branch.** Rejected. The review graded an exact tip; a rebase produces new commits nobody reviewed.
- **Take one side for every conflicted capture file.** Rejected. The journals are append-only; a one-side rule silently discards the other session's entries.
- **A hook that refuses commits on the default branch.** Deferred. It needs a reviewer enforcement decision, and the commit verb is still the caller for some flows.
- **A merge rule for the board.** Deferred. A board conflict refuses with the repair named; the drain session repairs by merge.

## Consequences

- Two sessions can land in either order. The second composes onto the moved default branch, and only a genuine conflict outside the table stops it.
- Hand landings on the default branch stop being the writes that move it under everyone else.
- A worktree that loses its request token cannot land until it is reauthorized; the refusal names that verb.
- No Bench verb yet moves a retained worktree onto the destination; the repair names raw Git until one exists.
