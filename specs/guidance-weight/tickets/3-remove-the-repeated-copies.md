# Remove the repeated rule copies and move the staleness exception into its pass

Blocked by: 2-repoint-the-build-rule-pointers.md
Writes: .agents/skills/bench-implement-spec/references/staleness-pass.md, .bench/BENCH.md, AGENTS.md, CHANGELOG.md, cmd/bench/anchor_help_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/testdata/anchors/fixture-repo/AGENTS.md, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/dogfood-referent-shipped, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/workflow-guidance-anchors/agents-system-suite-route, tests/canary/workflow-guidance-anchors/capture-sink-anchor, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer, tests/canary/workflow-guidance-anchors/core-absent-category-context, tests/canary/workflow-guidance-anchors/core-absent-category-oracle, tests/canary/workflow-guidance-anchors/core-absent-category-setup, tests/canary/workflow-guidance-anchors/core-absent-category-work, tests/canary/workflow-guidance-anchors/core-absent-gate-authority, tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship, tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback, tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source, tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner, tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner, tests/canary/workflow-guidance-anchors/core-absent-skills-guidance, tests/canary/workflow-guidance-anchors/core-absent-upgrade-route, tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response, tests/canary/workflow-guidance-anchors/core-cli-call-pair, tests/canary/workflow-guidance-anchors/core-progressive-loading-term, tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review, tests/canary/workflow-guidance-anchors/delegated-green-predecessor, tests/canary/workflow-guidance-anchors/delegated-opt-in-entry, tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop, tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint, tests/canary/workflow-guidance-anchors/delegated-repair-ownership, tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated, tests/canary/workflow-guidance-anchors/light-path-inline, tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor, tests/canary/workflow-guidance-anchors/ticket-light-path-anchor
Covers: GW16, GW17, GW18, GW19, GW20, GW21, GW22, GW51

## What to build

Every session loads `.bench/BENCH.md` and `AGENTS.md`, so each sentence there that an on-demand owner also states costs every session and adds no rule.
This ticket removes two such copies with their anchors.
They are the owned-red convergence bullet B12 of `.bench/BENCH.md` and the system-suite route bullet A17 of `AGENTS.md`.
It also moves the staleness-pass exception B10 of the spec-contradiction predicate, with its flag duty, into `.agents/skills/bench-implement-spec/references/staleness-pass.md`.
The staleness pass is the only phase that acts on that exception.
Each kept copy keeps its own anchor, and the moved sentence gets a new one, so a dropped rule still reds the gate.

This ticket reads three contracts from chunk GW1.
Ticket `1-move-the-build-rules-into-their-own-file.md` removed the acceptance-shortfall bullet and changed the lead-in to "Two predicates ride with them:".
Ticket 1 also added the section `### Always-loaded guidance weight` to `CHANGELOG.md`, with one bullet for chunk GW1.
Ticket `2-repoint-the-build-rule-pointers.md` repointed `staleness-pass.md:52`, so this ticket edits that file after it.
Line citations below are at base `bc2ea0ff`.
If a landing moves a line, find the text by its quotation.

**The operating guide.** Make these three edits to `.bench/BENCH.md`, and change no other line:

1. Change the lead-in "Two predicates ride with them:" to "One predicate rides with them:".
2. End the spec-contradiction bullet at "Ask about a behavioral one.". Line 62 ends at that sentence, and line 63 goes.
3. Remove the owned-red convergence bullet at lines 66-67.

The predicates then read as below:

```markdown
One predicate rides with them:

- If you find a **non-behavioral spec contradiction**, follow the current tree
  convention and flag it for reviewer veto. Ask about a behavioral one.
```

The row `internal/anchors/registry_data.go:376` keeps its needle, which sits on the kept lines 61-62.
The file shrinks, so its line row in `projects/benchkit.md:490` stays green.

**The staleness pass.** In `## Orchestrator review` of `staleness-pass.md`, the last sentence of the first paragraph (`:45-46`) becomes the two sentences below:

```markdown
It decides each finding by the spec-contradiction predicate of `.bench/BENCH.md`, a behavioral one included, and flags it.
It does not ask the reviewer.
```

Put the first sentence on one physical line, because the new row pins it whole.
The two sentences before it stay as they are.

**The new row.** Add one `RequireInSection` row in the group `AfterImplementSpec` to `chunkChainAnchors` in `internal/anchors/registry_chunk_chain.go`.
Put it directly after the `stalenessPass` row at `:21`.

| `File` | section | needle | diagnostic |
|---|---|---|---|
| `stalenessPass` | `Orchestrator review` | the first moved sentence | `stalenessPass + " dropped the orchestrator's decide-and-flag rule for each spec contradiction, a behavioral one included"` |

The diagnostic opens with the path, so none of the count-pin prefixes that ticket 1 lists reads it.
Correct the doc comment of `chunkChainAnchors` (`:3-17`) and the comment of `stalenessPass` (`:76`), so that both name the new row.

**The new test.** Add `TestStalenessExceptionAnchorRedsOnRemoval` to `internal/anchors/registry_chunk_chain_test.go`.
Write the path, the section, the needle, and the diagnostic by hand, as independent expectations.
Use `anchorHarness` with the group `AfterImplementSpec` and one rule, as `TestChunkChainAnchors` does (`:6-63`).
The harness passes a tree that holds the sentence, and it requires the diagnostic for a tree without it (`internal/anchors/anchor_harness_test.go:161-173`).

**The retired rows.** Remove these two rows from `internal/anchors/registry_data.go`:

- `:378`, the owned-red convergence pointer of `.bench/BENCH.md`. The `craft-line` copy at `.agents/skills/bench-craft-line/SKILL.md:80-81` keeps its row `:191`.
- `:383`, the system-suite route of `AGENTS.md`. The `projects/benchkit.md:676-677` copy keeps its row `:385` and its canary `benchkit-system-suite-route`.

Do not edit `projects/benchkit.md` or `.agents/skills/bench-craft-line/SKILL.md`.

**The working agreement.** Remove the last bullet of `AGENTS.md`, at lines 85-86.
The kit's `CLAUDE.md` imports `AGENTS.md`, so it needs no edit.

**The A17 dependents.** These three dependents of the row `internal/anchors/registry_data.go:383` retire with it:

- The `AGENTS.md` rule of `TestSystemSuiteRouteAnchorsRedOnRemoval` at `internal/anchors/registry_data_test.go:1244-1248`. The `projects/benchkit.md` rule stays. Rewrite the doc comment at `:1236-1239` so that it names only the profile, in no more lines.
- The canary `tests/canary/workflow-guidance-anchors/agents-system-suite-route`. Delete the directory with its `BASE`, `EXPECT`, and `MUTATE.json`.
- The plant at `cmd/bench/testdata/anchors/fixture-repo/AGENTS.md:11` and the blank line after it. Each later plant then sits two lines after the one before it.

`cmd/bench/anchor_help_test.go` derives its rows and its lines from the registry and the fixture (`:23-40`, `:82-89`), so it needs no edit.
`cmd/bench/anchors_dir_test.go` and `cmd/bench/command_registry_test.go:303` also read the fixture, and they count from the registry too.
`internal/handoffdoc/testdata/legacy-session-handoff.md:15` names the canary inside a frozen record, so it stays.

**The changelog.** Add one bullet for chunk GW2 to `### Always-loaded guidance weight` in `CHANGELOG.md`, after the GW1 bullet.
This is the proposed text:

```markdown
- Changed `.bench/BENCH.md` and the kit's `AGENTS.md` so that they no longer repeat three rules that an on-demand owner states. The owned-red convergence rule stays in `craft-line`, and the hand-run route of the system suite stays in `projects/benchkit.md`. The staleness-pass exception for a behavioral spec contradiction moved into `.agents/skills/bench-implement-spec/references/staleness-pass.md`, which now states that the orchestrator decides and flags each finding. Before, every session loaded these sentences.
```

**The structure limits.** `internal/anchors/registry_data.go` has 478 lines and `internal/anchors/registry_data_test.go` has 1,286 lines.
Both are over the 400-line cap with no grant (`internal/structure/structure.go:308-331`).
This ticket only removes lines from them, so neither gains a line.
The new row and the new test go in `internal/anchors/registry_chunk_chain.go` and `internal/anchors/registry_chunk_chain_test.go`, which stay under the cap.

**The Writes line.** The line holds each research path and the five bound files of `internal/anchors` and `cmd/bench` (`internal/tickets/registry_data.go:21-27`, `:44-45`).
The closure adds four anchor registry files that hold a literal equal to `.bench/BENCH.md` or `AGENTS.md` (`internal/anchors/references.go:31-59`).
It also adds the 40 other fixture directories that pin `.bench/BENCH.md`, `AGENTS.md`, or `CHANGELOG.md` through `BASE`, an overlay, or a mutation.
This ticket drops no research path.

`internal/anchors/registry_chunk_chain_test.go` holds `.bench/BENCH.md` only inside longer needles, so the closure does not name it.
This ticket names it because the new test goes there.
`internal/anchors/registry_chunk_chain.go` holds the staleness-pass path as the literal of `stalenessPass`, so the closure names it for that path too.

## Acceptance

- [ ] `.bench/BENCH.md` holds no owned-red convergence bullet, and `rg "diff-owned reds" .bench/BENCH.md` finds no line (GW16).
- [ ] The `craft-line` row `internal/anchors/registry_data.go:191` still binds its kept copy, and `TestRootConformance` passes (GW17).
- [ ] `projects/benchkit.md` keeps the hand-run route, and the kept rule of `TestSystemSuiteRouteAnchorsRedOnRemoval` passes (GW18).
- [ ] `AGENTS.md` and the fixture repo hold no A17 line, and `TestAnchorsReportsNeedleLines` passes (GW19).
- [ ] The canary `agents-system-suite-route` is gone, and `TestEveryRetainedFixtureBitesThroughRegisteredOwner` passes (GW20).
- [ ] `## Orchestrator review` of `staleness-pass.md` states the moved sentence, and `TestStalenessExceptionAnchorRedsOnRemoval` passes and reds without the new row (GW21).
- [ ] The spec-contradiction bullet of `.bench/BENCH.md` ends at "Ask about a behavioral one.", and `rg -i staleness .bench/BENCH.md` finds no line (GW22).
- [ ] The lead-in reads "One predicate rides with them:" above one bullet, which the GW2 review reads for GW14.
- [ ] The section `### Always-loaded guidance weight` of `CHANGELOG.md` holds the GW2 bullet after the GW1 bullet (GW51).

## Verification

Record each red in the return, with its command and its failure line.

- Remove the row at `internal/anchors/registry_data.go:383` first, and run `bench test --package ./internal/anchors --run TestSystemSuiteRouteAnchorsRedOnRemoval`. The `AGENTS.md` rule reds, because no row raises its diagnostic. Then remove that rule, and record the green run.
- With that row gone, run `bench test --package ./cmd/bench --run TestAnchorsReportsNeedleLines`. The test reds, because the old plant moves each later plant down two lines. Then remove the plant, and record the green run.
- Remove the bullet from `AGENTS.md`, and run `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner`. The canary `agents-system-suite-route` reds, because its `old` text is absent. Then delete the canary, and record the green run.
- Write the new test before its row, and run `bench test --package ./internal/anchors --run TestStalenessExceptionAnchorRedsOnRemoval`. The test reds, because no row raises its diagnostic. Then add the row and the moved sentence, and record the green run.

After the edit, record one observed red through `bench probe` for the new row:

- Omit the row with `--package ./internal/anchors --run TestStalenessExceptionAnchorRedsOnRemoval`. The test reds.
- Swap `a behavioral one included, and flags it.` in `staleness-pass.md` for `and does not flag it.` with `--check docs-currency-workflow`. The check reds with the new diagnostic. Run `bench worktree build <target>` first, because a named check compiles from the run binary's source.

Each probe reports a restored subject.
For GW16 and GW22, run the two `rg` searches of the acceptance rows, and record that each finds no line.
For GW19, run `rg -F "bench test --check system"` over `AGENTS.md` and the fixture repo, and record that it finds no line.
Then run `bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench`.
Also run the checks `docs-currency-workflow`, `prose-mechanics`, `guidance-prose-budgets`, and `root-conformance`.
