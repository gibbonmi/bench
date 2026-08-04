# Anchor the reclaim prose so a silent revert reds the gate

Blocked by: none
Ownership fence: `internal/conformance/docs_workflow_helpers_test.go`
Contracts: the reclaim sentence in each of the three always-loaded documents crosses those documents→the anchor table in `internal/conformance/docs_workflow_helpers_test.go` and is asserted by AC1 by reverting each document independently rather than checking one of the three
Assumptions: the three documents already carry correct reclaim wording and this ticket adds no prose; the anchor table's existing rows stay untouched so the eight-operation lifecycle keeps its own anchors

## What to build

The prose ticket in this build claimed its acceptance row was enforced by "enumerating the
parser's operation set against every advertisement". No such check exists. The anchor table
still pins only the old eight-operation lifecycle, so deleting the reclaim sentence from
any of the three always-loaded documents leaves the gate green. The wording is correct and
completely unguarded — which is the state the one-source-per-fact rule exists to prevent,
because the next edit to any of those files can silently drop it.

Add anchors covering the reclaim sentence in `.bench/BENCH.md`,
`.agents/commands/bench-implement-spec.md`, and `projects/benchkit.md`, each with a
diagnostic naming what was lost. Follow the existing row shape in the table rather than
inventing a second mechanism; the table is already the one place this kind of fact is
pinned.

Anchor the claim that carries the meaning — that `reclaim` is maintainer-run and sits
outside the lifecycle a build harness drives — not an incidental phrase that a legitimate
rewording would break for no reason.

## Acceptance

- [ ] [AC1] deleting the reclaim sentence from any one of the three documents turns conformance red, asserted independently for each of the three rather than for one standing in for the rest.
- [ ] [AC2] each new anchor's diagnostic names the document and the lost claim, so the failure says what to restore.
- [ ] [AC3] the pre-existing anchors still pass unchanged, so the eight-operation lifecycle keeps the protection it already had.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AC1 | add anchors for only `.bench/BENCH.md` and leave the other two documents unpinned | the per-document anchor test | omit the two rows, revert the `projects/benchkit.md` sentence, run `BENCH_CONFORMANCE_ROOT="$PWD" go test ./internal/conformance -run TestRootConformance -count=1 -timeout 300s`, expect it to stay green when it must go red |
| AC2 | give the new anchors an empty diagnostic string | the diagnostic-text test | blank the `diag` field, revert one sentence, run the same command, expect the failure to name no document |
| AC3 | replace an existing lifecycle anchor's needle with the reclaim sentence | the existing-anchors test | swap the needle, run the same command against the unmodified tree, expect the eight-operation anchor to fail |
