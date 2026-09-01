# State the refusal shape in the reference doc

Blocked by: declare-the-landing-refusal-registry.md
Writes: .bench/BENCH-reference.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, tests/canary/workflow-guidance-anchors, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape (new), tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/benchref-imported, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LRS19

## What to build

`.bench/BENCH-reference.md` states the refusal shape the code enforces, so the
doc and the code agree. The anchors registry requires that sentence. A new
canary fixture bites when a doc edit drops the sentence.

The doc edit keeps the anchored sentences verbatim. These are "The spec is
optional on the landing and on its resume" and the fast-lane landing-shape
marker. The doc edit moves lines that fixtures under
`tests/canary/workflow-guidance-anchors/` compare, so the fixture bite test
belongs in your focused checks.

The sentence states the contract that the registry constructor enforces. Every
landing face constructs through that constructor, and the constructor takes the
route as a required argument. Every ticket that writes
`internal/worktree/land_refusal.go` builds serially after the registry ticket,
one at a time.

## Acceptance

- [ ] LRS19 — the anchors registry requires the refusal-shape sentence in
      `.bench/BENCH-reference.md`.
- [ ] The new fixture reds when the doc drops that sentence.
- [ ] The two existing anchored sentences stay verbatim.
- [ ] The doc sentence states that every landing face constructs through the
      registry constructor.

## Delegate charge

You work in the Bench repo on the `landing-refusal-standard` spec. Line: opus /
medium. Effort: medium, at most 3 iterations.

Read `specs/landing-refusal-standard/spec.md` first. Read
`internal/anchors/registry_data.go` and
`internal/anchors/registry_data_test.go`. Read
`tests/canary/workflow-guidance-anchors/reference-gate-authority/` as the
fixture prior art. Read the landing section of `.bench/BENCH-reference.md`.

Your blocker landed the registry constructor in
`internal/worktree/land_refusal.go`. State the shape that constructor enforces;
do not invent a second shape.

Keep the anchor needle on one physical line. Run the canary fixture bite test
in your focused checks, beside the anchors package tests. LRS19 rides in
`TestWorktreeRuleAnchorsRedOnRemoval`.

The `workflow-guidance-anchors`, `skills-index-command-adapters`, and
`docs-currency-token-diet` fixtures pin `.bench/BENCH-reference.md`. Your
`Writes` names each of them. The five `cmd/bench` and
`internal/conformance` entries are the
registry closure for the `internal/anchors` package. Edit any of them only if
your change reaches it.

Coverage rows: LRS19. Show LRS19 red before your edit. Show LRS19 green after.
Return the red-to-green log.

Self-probe with an omission mutation. Delete the refusal-shape sentence from
the doc and report the observed result. If the mutation returns green, add the
missing row.

Run `bench worktree exec "<label>" -- go test -parallel 2 ./internal/anchors/ ./internal/conformance/`.
The exec form is the only command form. Do not use `cd`. Do not commit. Do not
edit the spec.
