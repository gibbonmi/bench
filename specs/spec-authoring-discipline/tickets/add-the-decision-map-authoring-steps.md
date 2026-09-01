# Add the decision-map authoring steps and the asset convention

Blocked by: add-a-sources-example-to-the-map-template.md, split-the-wrapped-field-refusal.md, name-the-reader-sweep-and-move-the-ship-test.md
Writes: .agents/commands/bench-shape-idea.md, .agents/commands/bench-write-spec.md, internal/maps/schema.go, internal/maps/maps_parse_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/shape-idea-read-ready-sources (new), tests/canary/workflow-guidance-anchors/shape-idea-skeleton-checks (new), tests/canary/workflow-guidance-anchors/prose-preflight-term (new), tests/canary/workflow-guidance-anchors/decision-map-asset-path (new), tests/canary/skills-index-command-adapters, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SAD34, SAD35, SAD36, SAD37, SAD38

## What to build

`.agents/commands/bench-shape-idea.md` owns the decision-map authoring steps. The
phase file tells the author to read one ready decision map's Sources block before
the first write. The phase file then tells the author to run `bench maps` and
`bench gate-prose` on the first skeleton. The file must name both verbs, because
one verb alone leaves one lane unrun.

No kit file contains the term "prose preflight". A forbid tuple holds that
property. The named verb `bench gate-prose` stays the one handle.

Three files each spell the path `decisions/assets/`:
`.agents/commands/bench-write-spec.md`, `.agents/commands/bench-shape-idea.md`,
and the decision-map template in `internal/maps/schema.go`. The three copies must
agree exactly. A map-owned asset lives under that path, so the candidate scanner
never reads an asset as a map.

The scanner keeps its current skip behavior. It lists only the direct children of
a decisions directory, so it already skips the nested asset directory.
`TestDiscoverDecisionMapCandidatesDirectChildren` at line 205 in
`internal/maps/maps_parse_test.go` already writes `decisions/assets/nested.md`
and expects the skip. Add no Go change to the scanner.

This ticket is the last one that writes `internal/maps`, so it carries the
invariant of that whole package. The rendered template must still validate with
zero diagnostics after your asset-path line joins it. The eleven single-count
answer replacements must stay green.

The five `cmd/bench` and `internal/conformance` entries in `Writes:` are the
registry closure of the `internal/maps` and `internal/anchors` packages. Edit one
only when your change reaches it.
`tests/canary/skills-index-command-adapters` pins
`.agents/commands/bench-write-spec.md`, so that family is in your closure too.

## Acceptance

- [ ] SAD34 — the shape-idea phase file names the ready-map Sources read before the first write.
- [ ] SAD35 — the shape-idea phase file names `bench maps` and `bench gate-prose` for the first skeleton.
- [ ] SAD36 — no kit file contains the term "prose preflight".
- [ ] SAD37 — the two phase files and the template each spell `decisions/assets/`.
- [ ] SAD38 — the scanner still skips a map file nested under a decisions asset directory.
- [ ] The rendered template still validates with zero diagnostics.
- [ ] Each existing `workflow-guidance-anchors` fixture still bites.

## Delegate charge

You work in the Bench repo on the `spec-authoring-discipline` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/spec-authoring-discipline/spec.md` first. Then read
`.agents/commands/bench-shape-idea.md` in full. Read
`.agents/commands/bench-write-spec.md`, which your blocker changed. Read
`DecisionMapTemplate` in `internal/maps/schema.go`, which your blocker changed.
Read `TestDiscoverDecisionMapCandidatesDirectChildren` at line 205 in
`internal/maps/maps_parse_test.go`. Read
`tests/canary/workflow-guidance-anchors/shape-idea-map-template/` as the fixture
prior art.

`.agents/commands/bench-shape-idea.md` carries no line budget.
`.agents/commands/bench-write-spec.md` sits at 73 of 73 lines, so keep its line
count.

Spell `decisions/assets/` the same way in all three files. Add no Go change to
the scanner.

Coverage rows: SAD34, SAD35, SAD36, SAD37, SAD38. Show each row red before your edit. Show each
row green after. Return the red-to-green log.

Self-probe with an omission mutation. Spell a different asset path in the
template and report the observed result. If the mutation returns green, add the
missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/maps/ ./internal/anchors/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
