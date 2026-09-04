# State the hook rule in the guidance

Blocked by: none
Writes: .bench/BENCH-reference.md, docs/adr/0001-working-tree-gate-tripwire.md, roadmap/FT141.md, docs/reporesident-distillation.md, tests/canary/workflow-guidance-anchors/distillation-reduced-schema-refactor-shape, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route
Covers: PR15, PR16, PR17, PR18, PR19

## What to build

Verify the premise first: read the `Hook Layers` list in
.bench/BENCH-reference.md, whose first bullet names `bench gate pin`, the
drift clause, and the disarmed drift check. Read the anchor row in
internal/anchors/registry_data.go whose needle is the push-rule sentence,
and the fixture under
tests/canary/workflow-guidance-anchors/reference-agent-push-rule. Read
docs/adr/0001-working-tree-gate-tripwire.md, roadmap/FT141.md, and line 187
of docs/reporesident-distillation.md.

Rewrite the first `Hook Layers` bullet. It states that the hook blocks a
direct push to the default branch. It states that the reviewer lifts the
clause with `git config bench.allowProtectedPush true`. It states that
guard discovery reports a static deny surface while enforcement stays live.
Keep the push-rule sentence byte-identical on its own line. Name no
repo-only path beside a claim word.

Rewrite ADR 0001 under its current file name. The title and the body record
that planted-reason proofs defend a working-tree gate change and that the
reviewer's merge review is the human check. Remove the pin paragraph and
keep the threat-model sentence. Write no file path and no code snippet.

Replace the first paragraph of roadmap/FT141.md. It says nothing is built,
the name the row proposes is free, and the reviewer decides the name and
whether the verb is new. Quote no name the row does not spell. Keep `Next: decide` and every ledger line.
Edit the distillation doc's parenthetical so it names `internal/git` alone.

Write every sentence in ASD-STE100, and run `bench gate-prose` over the four
files before the commit.

## Acceptance

- [ ] The reference guide's first hook-layer bullet holds no `gate pin`, `drift`, or `pinned`.
- [ ] The `reference-agent-push-rule` fixture reports its diagnostic only when its mutation runs.
- [ ] ADR 0001 holds no `pin`, and its title ends with `planted-reason proofs`.
- [ ] The first paragraph of the FT141 detail holds the word `free`, and the file keeps `Next: decide`.
- [ ] The distillation doc names `internal/git` alone as the tree-hash owner.
- [ ] Self-probe: drop the push-rule sentence, and report the fixture bite red.
