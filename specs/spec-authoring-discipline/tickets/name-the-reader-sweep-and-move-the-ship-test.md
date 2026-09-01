# Name the reader sweep and move the ship-test question

Blocked by: require-a-row-for-each-in-scope-promise.md
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-spec/references/map-discipline.md, .agents/commands/bench-write-spec.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/craft-spec-reader-sweep-name (new), tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers (new), tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers (new), tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound (new), tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence (new), tests/canary/workflow-guidance-anchors/reader-sweep-term (new), tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence (new), tests/canary/workflow-guidance-anchors/write-spec-branch-split (new), tests/canary/skills-index-command-adapters, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SAD22, SAD23, SAD24, SAD25, SAD26, SAD27, SAD28, SAD29, SAD30, SAD31, SAD32, SAD33

## What to build

The existing sweep sentence in `.agents/skills/bench-craft-spec/SKILL.md` gains
the name **reader sweep**. That file has two spare lines, so only the name
extension lands there. The sweep rules land in
`.agents/skills/bench-craft-spec/references/map-discipline.md`, which carries no
line budget.

Four sweep rules land in the reference file. The sweep lists each named consumer
of the decision fact. The sweep lists each helper a named consumer calls
directly. A deeper callee joins the sweep only when the callee reads the fact
itself. Each shared reader takes an exact ownership fence.

The kit writes the two-word term "reader sweep" and never "reader census". The
word "census" still means the raw-call census in `CONTEXT.md` or the architecture
census in `projects/benchkit.md`. A forbid tuple holds this property, because
absence is what the rule states. The glossary ticket
`define-the-map-glossary-terms.md` writes the same two-word term into
`CONTEXT.md`. Both tickets spell the needle "reader sweep" the same way, so the
forbid tuple and the glossary entry cannot disagree.

`.agents/commands/bench-write-spec.md` sequences the reader sweep before the
`craft-tickets` charge. The phase file names the sweep and states no sweep rule
of its own, because `craft-spec` owns the rules. Row SAD29 is review-owned on the
Standards axis, and it grades this ticket's phase-file diff: only a reader sees
a restated rule.

The ship-test question moves down beside the ticket-slicing step. The move keeps
the needle "could a narrower capability ship on its own gate" byte for byte, and
a whole-file require anchor holds that needle. The moved sentence carries the
split arm, which states that the ticket graph splits where a consumer branch
lands green alone. The review-round paragraph keeps no copy of the question.

The anchor mechanism scopes a sectioned anchor to an H2 section. This file holds
the slicing step and the review-round paragraph under one `## Process` heading.
So no section-scoped anchor can express the placement. Rows SAD30 and
SAD31 are review-owned on the Standards axis for that reason.

The phase file sits at 73 of 73 lines and must net zero. The removal from the
review-round paragraph frees two lines, and the relocated question with its split
arm takes two. The sweep sequencing clause edits an existing step sentence in
place, so it takes no line. Fixtures under
`tests/canary/workflow-guidance-anchors/` pin exact bytes near the slicing step,
so run the fixture bite test in your focused checks.

The five `cmd/bench` and `internal/conformance` entries in `Writes:` are the
registry closure of the `internal/anchors` package. Edit one only when your
change reaches it. `tests/canary/skills-index-command-adapters` pins
`.agents/commands/bench-write-spec.md`, so that family is in your closure too.

## Acceptance

- [ ] SAD22 — the craft-spec sweep sentence names the reader sweep.
- [ ] SAD23 — the reference requires each named consumer of the decision fact.
- [ ] SAD24 — the reference requires each helper a named consumer calls directly.
- [ ] SAD25 — the reference admits a deeper callee only when the callee reads the fact.
- [ ] SAD26 — the reference requires an exact ownership fence for each shared reader.
- [ ] SAD27 — no kit file writes "reader census" for this sweep.
- [ ] SAD28 — the phase file sequences the reader sweep before the craft-tickets charge.
- [ ] SAD29 — the phase file names the sweep and states no sweep rule of its own.
- [ ] SAD30 — the ship-test question sits inside the ticket-slicing step (review-owned).
- [ ] SAD31 — the review-round paragraph holds no copy of the question (review-owned).
- [ ] SAD32 — the moved sentence carries the consumer-branch split arm.
- [ ] SAD33 — the phase file stays at 73 lines or fewer.
- [ ] The needle "could a narrower capability ship on its own gate" survives byte for byte.
- [ ] Each existing `workflow-guidance-anchors` fixture still bites.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/spec-authoring-discipline/spec.md` first. Then read
`.agents/commands/bench-write-spec.md` in full. Read the phase-file tuples at
lines 25 to 47 in `internal/anchors/registry_data.go`. Read the
narrower-capability needle at line 132. Read the ticket-slicing needles at lines
181 to 184.

Read the sweep sentence in `.agents/skills/bench-craft-spec/SKILL.md`. Read
`.agents/skills/bench-craft-spec/references/map-discipline.md`, which your
blocker changed.

Move the ship-test sentence pair without one byte of change. Add the sweep
sequencing sentence. Keep the phase file at 73 lines or fewer.

Write "reader sweep" everywhere. Never write "reader census". The glossary ticket
spells the same term, so do not vary it.

Add a forbid tuple for the census spelling. Do not add a section-scoped tuple for
the ship-test question, because both positions share one `## Process` heading.

Coverage rows: SAD22 through SAD33. Show each anchored row red before your edit. Show
each row green after. Return the red-to-green log. Rows SAD30 and SAD31 are
review-owned, so report the placement in prose instead.

Self-probe with an omission mutation. Drop the split arm from the moved sentence
and report the observed result. If the mutation returns green, add the missing
row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/anchors/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
