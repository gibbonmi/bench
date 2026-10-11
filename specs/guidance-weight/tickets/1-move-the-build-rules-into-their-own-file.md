# Move the build rules into their own reference file

Blocked by: none
Writes: .agents/skills/bench-implement-spec/references/build-rules.md (new), .bench/BENCH.md, CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/implementation_continuation_test.go, internal/conformance/retained_workflow_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/docs-currency-token-diet/dogfood-referent-shipped, tests/canary/docs-currency-token-diet/missing-cli-inventory, tests/canary/docs-currency-token-diet/stale-cli-doc-reference, tests/canary/load-validity-metadata/readme-shared-rule-drift, tests/canary/load-validity-metadata/shared-rule-drift, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/workflow-guidance-anchors/capture-sink-anchor, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer, tests/canary/workflow-guidance-anchors/core-absent-category-context, tests/canary/workflow-guidance-anchors/core-absent-category-oracle, tests/canary/workflow-guidance-anchors/core-absent-category-setup, tests/canary/workflow-guidance-anchors/core-absent-category-work, tests/canary/workflow-guidance-anchors/core-absent-gate-authority, tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship, tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback, tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source, tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner, tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner, tests/canary/workflow-guidance-anchors/core-absent-skills-guidance, tests/canary/workflow-guidance-anchors/core-absent-upgrade-route, tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response, tests/canary/workflow-guidance-anchors/core-cli-call-pair, tests/canary/workflow-guidance-anchors/core-progressive-loading-term, tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review, tests/canary/workflow-guidance-anchors/delegated-green-predecessor, tests/canary/workflow-guidance-anchors/delegated-opt-in-entry, tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop, tests/canary/workflow-guidance-anchors/delegated-per-ticket-author, tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint, tests/canary/workflow-guidance-anchors/delegated-repair-ownership, tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated, tests/canary/workflow-guidance-anchors/light-path-inline, tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor, tests/canary/workflow-guidance-anchors/ticket-light-path-anchor
Covers: GW01, GW02, GW03, GW04, GW05, GW06, GW07, GW10, GW11, GW13, GW14, GW15, GW51, GW52

## What to build

A session outside a build no longer loads the build-orchestration clauses.
`/bench-implement-spec`, `/bench-review-implementation`, and `/bench-final-check` read them in `.agents/skills/bench-implement-spec/references/build-rules.md`.
`.bench/BENCH.md` keeps one pointer to that file.
Each anchor, mirror, and canary of a moved clause follows the clause, so a dropped or reworded rule still reds the gate.

This ticket leaves the pointer lines of the phase files and craft skills as they are.
They still name `.bench/BENCH.md`, and its new pointer leads to the build rules.
Ticket `2-repoint-the-build-rule-pointers.md` repoints those lines next.
Line citations below are at base `bc2ea0ff`.
If a landing moves a line, find the text by its quotation.

**The build rules file.** Create `.agents/skills/bench-implement-spec/references/build-rules.md` in this order:

1. The heading `# Build rules`.
2. The lead paragraph, which is this one sentence: "`/bench-implement-spec`, `/bench-review-implementation`, and `/bench-final-check` load these rules for each spec-backed build."
3. The six sections of the table below, in table order.

| section | clauses | source in `.bench/BENCH.md` |
|---|---|---|
| `## Acceptance shortfall` | B11 | `:64-65` |
| `## Ticket authors` | B26 | `:116` |
| `## Chunks and review` | B27, B28 | `:118`, `:120` |
| `## Repairs` | B29 | `:122` |
| `## Plan expansion` | B30, B31 | `:124`, `:126` |
| `## Delegated runs` | B32, B33 | `:128`, `:130` |

The B11 bullet becomes this paragraph: "If a build cannot meet an acceptance row, that is a **material acceptance shortfall**: the build exits and reports. It does not land a silent partial."
Every other clause moves byte for byte, with its bold label and its reviewer voice ("me", "my").
Each clause stays one paragraph, so each paragraph keeps its sentence count and passes `prose-mechanics` as it does today.

**The operating guide.** Make these four edits to `.bench/BENCH.md`, and change no other line:

1. Change the lead-in at line 59 to "Two predicates ride with them:".
2. Remove the B11 bullet at lines 64-65.
3. Remove the eight paragraphs at lines 116-130.
4. Put the pointer paragraph below after the numbered phase list that ends at line 114.

```markdown
`.agents/skills/bench-implement-spec/references/build-rules.md` holds the build rules for authors, chunks, repairs, plan expansion, and delegated runs.
```

The pointer then sits between the phase list and the paragraph that opens "**Right-size the process; ask before deviating.**".
The file shrinks, so its line row in `projects/benchkit.md:490` stays green.

**The moved rows.** Each row below keeps its group, kind, and needle, and its `File` becomes the build rules.
Spell the path once, as a named constant beside `stalenessPass` in `internal/anchors/registry_chunk_chain.go:79-84`.
The 21 `Require` rows sit in three registry files:

- `internal/anchors/registry_data.go:377`, the B11 row.
- `internal/anchors/registry_ft311_review_dispatch.go:39`, `:40`, `:45`, `:46`, `:47`, `:48`, `:51`, and `:52`.
- `internal/anchors/registry_retained_workflow.go:213`, `:214`, `:215`, `:216`, `:237`, `:238`, `:239`, `:240`, `:241`, `:242`, `:294`, and `:341`.

The row at `internal/anchors/registry_retained_workflow.go:341` is `RequireInSection`, and its section changes from `Workflow` to `Ticket authors`.
The five moved `Forbid` rows are `internal/anchors/registry_ft311_review_dispatch.go:41`, `:49`, `:50`, and `internal/anchors/registry_retained_workflow.go:217`, `:270`.
The `Forbid` rows at `internal/anchors/registry_retained_workflow.go:295-302` keep their phase files.
The spec keeps GW01 review-owned, so add no `Forbid` row that binds `.bench/BENCH.md` for the moved text.

**The diagnostics.** In each moved row, "operating guide" becomes "build rules", and the verb takes the plural form where it needs it.
For example, the row at `internal/anchors/registry_ft311_review_dispatch.go:47` reads "build rules allow advancement before three-axis chunk review".
The row at `internal/anchors/registry_data.go:377` opens with `.bench/BENCH.md`, so it opens with the build rules path instead.
Each diagnostic keeps its family prefix: `retained workflow: `, `fresh ticket author: `, or `implementation continuation: `.
`TestRetainedWorkflow`, `TestFreshTicketAuthors`, and `TestImplementationContinuation` select their rows by these prefixes.

These unmoved rows name the operating guide as the owner of a moved rule, so their diagnostics name the build rules too:

- `internal/anchors/registry_retained_workflow.go:46`, the craft-spec copy of the implementation-chunk contract.
- `internal/anchors/registry_retained_workflow.go:245`, the craft-delegate copy of the fresh author session rule.
- `internal/anchors/registry_retained_workflow.go:295-302`, the phase-file copies of the chunk-review start, review-repeat, final reconciliation, and chunk-tip review fence rules.
- `internal/anchors/registry_retained_workflow.go:345`, the craft-delegate copy of the diagnostic helper boundary.

"that the operating guide owns" becomes "that the build rules own", and "the operating guide's" becomes "the build rules'".
Each of these rows keeps its `File`, kind, and needle.
The pointer rows at `internal/anchors/registry_retained_workflow.go:47` and `:303` keep their diagnostics here, because their pointers change in ticket 2.

**The owner comments.** Correct each comment that names `.bench/BENCH.md` or the operating guide as the owner of a moved sentence.
The research found them at `internal/anchors/registry_ft311_review_dispatch.go:72-74` and `internal/anchors/registry_retained_workflow.go:319-321`.
Also read `internal/anchors/registry_retained_workflow.go:224-235`.
The doc comment of `chunkChainAnchors` at `internal/anchors/registry_chunk_chain.go:3-17` describes the two new rows.

**The new rows.** Add two `Require` rows in the group `AfterImplementSpec` to `chunkChainAnchors`:

| row | `File` | kind and section | needle | diagnostic |
|---|---|---|---|---|
| GW03 | `.bench/BENCH.md` | `RequireInSection`, `Workflow` | the pointer sentence | `.bench/BENCH.md Workflow section dropped its pointer to the build rules` |
| GW04 | the build rules | `Require` | the lead sentence | `.agents/skills/bench-implement-spec/references/build-rules.md is absent or dropped its build-phase lead` |

No new diagnostic opens with a prefix that a count pin reads.
Those prefixes are `prepared build guidance: `, `prepared review dispatch: `, `consumer context: `, and `workflow integration source: `.
The pins are at `internal/conformance/ft311_preparation_test.go:22-23`, `:40-41`, `:124-125`, and `internal/conformance/docs_workflow_helpers_test.go:521-522`.

**The new tests.** Add both tests to `internal/anchors/registry_chunk_chain_test.go`, and write each path, needle, and diagnostic by hand:

- `TestBuildRulesPointerAnchorsRedOnRemoval` passes a tree whose `## Workflow` holds the pointer. It requires the GW03 diagnostic for a tree without the pointer, and for a tree whose pointer sits outside `## Workflow`.
- `TestBuildRulesReferenceAnchorsRedOnAbsence` follows `TestReferenceFileAnchorsRedOnAbsence` at `internal/anchors/registry_data_test.go:571-610`. An absent file raises `acceptance coverage anchor file missing: <path>`. A file without its lead raises the GW04 diagnostic. A file with its lead raises neither.

The second test can reuse `writeReferences` at `internal/anchors/registry_data_test.go:614`.

**The mirrors.** These independent expectations take the moved rows:

- `internal/conformance/retained_workflow_test.go:56-60`, `:64-67`, `:72-75`, and `:88` in `wantDiagnostics` take the renamed diagnostics.
- `internal/conformance/retained_workflow_test.go:118-137`, `:151-166`, and `:183-202` in `wantPredicates` take the renamed keys and the build rules path.
- `internal/conformance/implementation_continuation_test.go:81-84` takes the renamed key, the build rules path, and the section `Ticket authors`.
- `internal/conformance/implementation_continuation_test.go:97` takes the renamed diagnostic of the row at `internal/anchors/registry_retained_workflow.go:345`.

**The canaries.** Edit these fixtures under `tests/canary/workflow-guidance-anchors/`:

- `delegated-chunk-tip-review`, `delegated-green-predecessor`, `delegated-opt-in-entry`, `delegated-pending-predecessor-stop`, `delegated-prerequisite-checkpoint`, and `delegated-repair-ownership`. In each, `BASE` and the `MUTATE.json` path name the build rules, and `EXPECT` takes the renamed diagnostic. The `old` and `new` strings stay.
- `delegated-per-ticket-author/EXPECT` takes the renamed diagnostic of the row at `internal/anchors/registry_retained_workflow.go:245`.

**The changelog.** Add the section `### Always-loaded guidance weight` as the first section under `## [Unreleased]` of `CHANGELOG.md`.
Give it one bullet for chunk GW1, in the form of the current entries.
Ticket 2 extends this bullet, and GW2 to GW5 add their own bullets.
This is the proposed text:

```markdown
### Always-loaded guidance weight

- Changed `.bench/BENCH.md` so that it no longer holds the build rules. The acceptance-shortfall predicate and the rules for ticket authors, chunks and review, repairs, plan expansion, and delegated runs moved word for word into `.agents/skills/bench-implement-spec/references/build-rules.md`. The build, review, and final-check phases load that file, and `.bench/BENCH.md` keeps a one-line pointer to it. Before, every session loaded these rules.
```

**The structure limits.** `internal/anchors/registry_data.go` has 478 lines and `internal/anchors/registry_data_test.go` has 1,286 lines.
Both are over the 400-line cap with no grant, so neither file gains a line (`internal/structure/structure.go:308-331`).
`internal/anchors/registry_retained_workflow.go` has 397 lines, so its edits stay in place.
The build rules file has no budget row, because reference files stay outside the budget universe.

**The Writes line.** The line holds each research path that this design can touch, and the five bound files of `internal/anchors` (`internal/tickets/registry_data.go:21-27`, `:45`).
The closure adds `internal/anchors/registry_commitment.go` and `internal/anchors/anchor_harness_diagnostics_test.go`, which name `.bench/BENCH.md`.
It also adds the 34 other fixture directories that pin `.bench/BENCH.md` through `BASE`, an overlay, or a mutation.
Ticket 2 holds the research paths that only the pointer edits can touch, and this ticket drops no research path.

## Acceptance

- [ ] `.bench/BENCH.md` holds no needle of B11 or B26 to B33 (GW01).
- [ ] The build rules hold B11 and B26 to B33 word for word, in base order, under the six sections (GW02).
- [ ] `## Workflow` of `.bench/BENCH.md` holds the pointer, and a tree without it raises the GW03 diagnostic (GW03).
- [ ] The build rules open with the lead, and an absent or leadless file raises its diagnostic (GW04).
- [ ] Each of the 21 moved `Require` rows binds the build rules with its base needle, and `TestRootConformance` passes (GW05).
- [ ] The continuation row names the section `Ticket authors`, and `TestImplementationContinuation` passes (GW06).
- [ ] The `Forbid` row at `internal/anchors/registry_retained_workflow.go:217` binds the build rules, and `TestRetainedWorkflow` passes (GW07).
- [ ] The four fresh-author `Forbid` rows bind the build rules (GW52).
- [ ] The mirrors name the build rules for each moved row (GW10).
- [ ] Each of the six delegated canaries mutates the build rules and bites its row (GW11).
- [ ] No moved row's diagnostic names the operating guide, and no listed unmoved row names it as an owner (GW13).
- [ ] The lead-in reads "Two predicates ride with them:" above two bullets (GW14).
- [ ] The build rules pass `prose-mechanics` (GW15).
- [ ] `CHANGELOG.md` holds the new section with the GW1 bullet (GW51).

## Verification

Record each red in the return, with its command and its failure line.

- Write both new tests before the rows, and run `bench test --package ./internal/anchors --run TestBuildRules`. Both tests fail, because no row raises their diagnostics. Then add the rows, and record the green run.
- Repoint the mirrors before the rows, and run `bench test --package ./internal/conformance --run 'TestRetainedWorkflow|TestImplementationContinuation'`. Each repointed mirror reports its row as absent or drifted. Then repoint the rows, and record the green run.

After the edit, record one observed red through `bench probe` for each repointed mirror and each new row:

- For each repointed mirror, swap the `File` of its row back to `.bench/BENCH.md`. For the row at `:341`, swap the section back to `Workflow`. Run the probe with `--package ./internal/conformance --run TestRetainedWorkflow` or `--run TestImplementationContinuation`. The mirror reds.
- For each new row, omit the row with `--package ./internal/anchors --run TestBuildRules`. Its test reds.
- Omit the pointer sentence from `.bench/BENCH.md` with `--check docs-currency-workflow`. The check reds with the GW03 diagnostic. Run `bench worktree build <target>` first, because a named check compiles from the run binary's source.

Each probe reports a restored subject.
Then run `bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench`.
Also run the checks `docs-currency-workflow`, `prose-mechanics`, `guidance-prose-budgets`, `skills-index-command-adapters`, and `root-conformance`.
