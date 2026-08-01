# Reduce the gate run and inherit from a full-green ancestor

Blocked by: Compute the stripped subject identity, Run excludable phases against a stripped worktree, Expand the verdict record with the reduced shape

Ownership fence: `internal/gate/gate.go`, `internal/gate/reduced_run_test.go`, `internal/contract/runtime/runtime_gate_reduced_test.go`
Assumptions: the stripped identity, the stripped-worktree construction, and the reduced record shape all exist; the existing execution lock and fresh-green reuse path are unchanged and no new lock is introduced

## What to build

The execution path that actually reduces. When the changeset is confined to the
allowlist and a full-green ancestor's stripped identity matches the current tree's,
run only the included phases and write a reduced verdict naming that ancestor.

Ancestry is the part with the sharp edge. Consecutive capture commits all inherit
from the same full green, never from each other: a reduced verdict is not itself a
valid ancestor, because evidence attributed to a run that graded nothing is worse
than a slow gate. The naive single-slot implementation loses the ancestor on the
first reduced write and then either falls back to a full run on the second commit —
the motivating scenario, unserved — or chains, which is unsound. Fail closed at
every state where there is nothing sound to inherit: no ancestor at all, or an
ancestor older than the existing freshness window, both run the full gate.

The operator has to see it. A reduced run announces which phases it skipped and
which full-green run supplied their evidence, in the same spirit as the existing
reused-verdict announcement — a skipped run that says nothing reads as a gate that
never ran.

## Acceptance

- [ ] [R14] A capture-only changeset executes the included phases only, observable in the recorded phase list.
- [ ] [R17] A second consecutive capture commit inherits from the same full green as the first rather than falling back to a full run.
- [ ] [R18] A reduced verdict is never accepted as an ancestor, so no chain of reduced records can form.
- [ ] [R19] An ancestor older than the freshness window forces a full run, with inherited evidence never re-stamped to the current time.
- [ ] [R20] An allowlist-confined change with no ancestor runs the full gate rather than emitting a reduced verdict with an empty ancestor.
- [ ] [R24] The reduced run announces the phases it skipped and the ancestor whose evidence it inherited.
