# Staleness pass

`/bench-implement-spec` charges this pass after the build preflight. The pass
is a light assessment of an approved spec and its tickets against the current
`main` tree. Every orchestrator sends the same charge, so the cost stays small
and predictable.

## The charge

- Line: the mid tier, medium effort, one iteration.
- Budget: 15 tool calls.
- Access: read-only.
- Inputs: the spec path, the reviewed graph commit, and each preflight red row that the pass took.

## Procedure

1. List the drift. Run `git diff --name-only <reviewed graph commit> main`.
   Keep each path that a ticket `Writes:` line names. Keep each file that
   holds a symbol that the spec or a ticket cites. Add each file that a taken
   preflight red row names.
2. If the drift list is empty, return no contradiction and stop.
3. Read only the drifted files and the spec and ticket claims that name them.
   A claim is a current-code claim under the claim rules of `craft-spec`.
4. Classify each contradiction as behavioral or non-behavioral under the
   spec-contradiction predicate of `.bench/BENCH.md`.

A contradiction is blocking when a ticket author who follows the claim
builds against a false fact. Leave counts, coverage rows, fences, and caller
lists to the ticket authors and the chunk reviews, because they own those
checks.

## Return

Return one row for each blocking contradiction, and return no other finding.
Each row names the claim, the drifted file, and the class. The orchestrator
routes each row as the implementation command states.
