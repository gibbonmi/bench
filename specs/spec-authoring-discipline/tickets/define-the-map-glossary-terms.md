# Define the two map terms and the reader sweep in CONTEXT.md

Blocked by: add-the-decision-map-authoring-steps.md
Writes: CONTEXT.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/context-coverage-map-term (new), tests/canary/workflow-guidance-anchors/context-decision-map-term (new), tests/canary/workflow-guidance-anchors/context-reader-sweep-term (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SAD39, SAD40, SAD41, SAD42

## What to build

`CONTEXT.md` defines "coverage map" with an Avoid list. The Avoid list names the
bare word "map". `CONTEXT.md` keeps the "decision map" entry with its own Avoid
list, so the two map artifacts stay distinct. `CONTEXT.md` defines "reader sweep"
with an Avoid list that reserves "census".

The needle "reader sweep" is the term the ticket
`name-the-reader-sweep-and-move-the-ship-test.md` writes into the craft-spec
skill and its reference file. That ticket also adds a forbid tuple for the
census spelling. This ticket spells the same two-word term, so the glossary entry
and the forbid tuple cannot disagree.

Each entry ships in the FT144 shape: the entry text, one anchor tuple in
`internal/anchors/registry_data.go`, one red-on-removal registry test, and one
canary fixture under `tests/canary/workflow-guidance-anchors/`.

Row SAD42 is review-owned on the Standards axis, and it grades every new sentence
this build wrote. This ticket is the last one in the graph, so it carries that
row. A bare "map" reads as correct in isolation, so a reader must check the whole
build's new prose for the two-word terms.

This ticket is also the last one that writes `internal/anchors` and
`tests/canary/workflow-guidance-anchors`, so it carries the invariant of both.
Every anchor tuple must still red on removal, and every fixture in the family
must still bite.

The five `cmd/bench` and `internal/conformance` entries in `Writes:` are the
registry closure of the `internal/anchors` package. Edit one only when your
change reaches it.

## Acceptance

- [ ] SAD39 — CONTEXT.md defines "coverage map" and lists the terms to avoid.
- [ ] SAD40 — CONTEXT.md keeps the "decision map" entry with its Avoid list.
- [ ] SAD41 — CONTEXT.md defines "reader sweep" and reserves "census".
- [ ] SAD42 — each new sentence in this build writes a two-word map term.
- [ ] The three new fixtures red on their mutations and green on their restores.
- [ ] Every fixture under `tests/canary/workflow-guidance-anchors` still bites.
- [ ] Every anchor tuple in `internal/anchors` still reds on removal.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 2 iterations. Story 42 audits every sentence this
build adds or moves, so the audit needs the whole diff.

Read `specs/spec-authoring-discipline/spec.md` first. Then read `CONTEXT.md` in
full. Read the existing "decision map" entry and its Avoid list. Read
`tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary/` as the
fixture prior art. Read the anchor tuples your blockers added in
`internal/anchors/registry_data.go`.

Add exactly three glossary entries. Take the Avoid terms from the spec. Do not
invent a fourth entry.

Write "reader sweep" the same way the craft-spec skill writes it. Never write
"reader census".

Read the whole build's new prose for row SAD42. Report each bare "map" you find in a
new sentence.

Coverage rows: SAD39, SAD40, SAD41, SAD42. Show each anchored row red before your edit. Show
each row green after. Return the red-to-green log.

Self-probe with an omission mutation. Delete the "coverage map" entry and report
the observed result. If the mutation returns green, add the missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/anchors/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
