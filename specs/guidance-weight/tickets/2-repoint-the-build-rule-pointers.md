# Repoint the build-rule pointers to the build rules

Blocked by: 1-move-the-build-rules-into-their-own-file.md
Writes: .agents/commands/bench-debug.md, .agents/commands/bench-drain.md, .agents/commands/bench-final-check.md, .agents/commands/bench-implement-spec.md, .agents/commands/bench-review-implementation.md, .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/skills/bench-craft-gate/SKILL.md, .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-tickets/SKILL.md, .agents/skills/bench-implement-spec/references/staleness-pass.md, .bench/BENCH-reference.md, CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, docs/field-guide.html, internal/anchors/registry_chunk_chain.go, internal/anchors/registry_chunk_chain_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/conformance/axi_query_registry_test.go, internal/conformance/recurrence_maintenance_contract_test.go, internal/conformance/retained_workflow_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/delegated-serial-ticket-checkpoint, tests/canary/workflow-guidance-anchors/review-repair-ticket-covers, tests/canary/workflow-guidance-anchors/review-repair-ticket-owner
Covers: GW08, GW09, GW12, GW51

## What to build

A reader of a phase file or a craft skill follows each build-rule pointer straight to `.agents/skills/bench-implement-spec/references/build-rules.md`.
No such line names `.bench/BENCH.md` as the owner of a moved rule.
Each anchor, mirror, and canary that quotes a repointed line takes the new text, so it stays green and still bites.

This ticket reads one contract from ticket `1-move-the-build-rules-into-their-own-file.md`.
The build rules file exists with B11 and B26 to B33, and `.bench/BENCH.md` holds its one-line pointer.
A pointer to a file that does not exist yet would mislead, so this ticket waits for that checkpoint.
Line citations below are at base `bc2ea0ff`.
If a landing moves a line, find the text by its quotation.

**The pointer lines.** On each line below, replace `.bench/BENCH.md` with `.agents/skills/bench-implement-spec/references/build-rules.md`.
Keep the code span, and change no other text.
Each edit stays on its physical line, so no file gains a line.

- `.agents/commands/bench-implement-spec.md:21`, `:23`, `:33`, `:54`, `:60`, `:90`
- `.agents/commands/bench-review-implementation.md:34`, `:44`, `:50`, `:52`
- `.agents/commands/bench-final-check.md:194`, `.agents/commands/bench-write-spec.md:50`, `.agents/commands/bench-debug.md:157`
- `.agents/skills/bench-craft-spec/SKILL.md:68`, `:70`, `.agents/skills/bench-craft-delegate/SKILL.md:13`, `:122`, `.agents/skills/bench-craft-gate/SKILL.md:28`
- `.agents/skills/bench-craft-tickets/SKILL.md:67`, `:95`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:10`, `:39`, `:94`, `:105`
- `.agents/skills/bench-implement-spec/references/staleness-pass.md:52`
- `docs/field-guide.html:1333`, inside its `<code>` element
- `.bench/BENCH-reference.md:279`, in "under the repair rule of `.bench/BENCH.md`"

The line `.agents/skills/bench-craft-tickets/SKILL.md:95` says "the operating guide's delegated policy" with no path.
It becomes "the delegated policy of `.agents/skills/bench-implement-spec/references/build-rules.md`".
Five command and skill files sit at their line rows (`projects/benchkit.md:491-497`), and the single-line edits keep them there.
On `.agents/skills/bench-craft-spec/SKILL.md:68`, change only the first sentence, because the GW5 ticket owns its budget sentence.

**The drain line.** `.agents/commands/bench-drain.md:72` stays as it is.
Its implement-now delegates are the light-path fixes that `.agents/commands/bench-drain.md:217` and `:222` dispatch.
Their authorship rule is the light-path fix rule at `.bench/BENCH.md:148`, which stays in `.bench/BENCH.md`.
So that line still names its owner, and its mirrors at `internal/conformance/recurrence_maintenance_contract_test.go:91` and `:217` stay too.

The `Writes:` line keeps both files, so a reviewer veto of this call needs no plan expansion.
Under that veto, the line takes the build rules path, and both mirrors take the new text.

**The anchors.** Each anchor that quotes a repointed line takes the new text in its needle:

- `internal/anchors/registry_retained_workflow.go:47`, `:212`, `:218`, `:219`, `:220`, `:280`, `:303`, `:307`, `:310`, `:316`
- `internal/anchors/registry_data.go:136`, `:140`, `:442`
- `internal/anchors/registry_chunk_chain.go:62`
- `internal/anchors/registry_ft311_review_dispatch.go:42`

Use the build rules constant of ticket 1 where a needle concatenates the path, or spell the needle whole.
The rows at `internal/anchors/registry_retained_workflow.go:47` and `:303` name the operating guide as the owner in their diagnostics.
In both, "the operating guide's" becomes "the build rules'", and each keeps its family prefix.

**The mirrors.** These independent expectations take the new text:

- `internal/conformance/retained_workflow_test.go:138-141`, `:146-150`, and `:167-178` in `wantPredicates`
- `internal/anchors/registry_data_test.go:349` and `:892`
- `internal/anchors/registry_chunk_chain_test.go:51`

**The canaries.** Edit these fixtures under `tests/canary/workflow-guidance-anchors/`:

- `review-repair-ticket-owner/MUTATE.json`: the `old` string and the first line of the `new` string take the new text of `.agents/commands/bench-review-implementation.md:52`.
- `review-repair-ticket-covers/MUTATE.json`: the `old` string takes the new text, and the `new` string takes the build rules path too.
- `delegated-serial-ticket-checkpoint/MUTATE.json`: the `old` string takes the new text of `.agents/skills/bench-craft-tickets/SKILL.md:95`.

Each `old` string occurs exactly once in its live file, and each `EXPECT` stays.

**The changelog.** Ticket 1 wrote the GW1 bullet.
Put the sentence below into that bullet, before its "Before," sentence, so the section keeps one bullet for each chunk.

```markdown
Each phase file and craft skill that cited `.bench/BENCH.md` for these rules now cites the new file.
```

**The structure limits.** `internal/anchors/registry_data.go` (478 lines) and `internal/anchors/registry_data_test.go` (1,286 lines) gain no line (`internal/structure/structure.go:308-331`).
Each needle edit stays in place in `internal/anchors/registry_retained_workflow.go`.

## Acceptance

- [ ] No line under `.agents`, `docs`, or `.bench/BENCH-reference.md` names `.bench/BENCH.md` as the owner of a moved rule (GW08).
- [ ] Each anchor that quotes a repointed line quotes its new text, and `TestRootConformance` passes (GW09).
- [ ] The three pointer canaries quote the new text in `MUTATE.json`, and each bites its row (GW12).
- [ ] The GW1 bullet of `CHANGELOG.md` names the repointed pointers (GW51).
- [ ] `.agents/commands/bench-drain.md:72` still names `.bench/BENCH.md` for the light-path delegate authorship.

## Verification

Record each red in the return, with its command and its failure line.

- Repoint the mirrors before the anchors, and run `bench test --package ./internal/conformance --run TestRetainedWorkflow`. Each repointed mirror reports its row as drifted. Then repoint the anchors and the guidance lines, and record the green run.
- Do the same for `bench test --package ./internal/anchors --run 'TestChunkChainAnchors|TestCraftSpecMapDisciplineAnchorsRedOnRemoval|TestRepairTicketOwnerAnchorsRedOnRemoval'`.

After the edit, record one observed red through `bench probe` for each repointed mirror:

- Swap the needle of its row back to the base text, with `--package ./internal/conformance --run TestRetainedWorkflow` or the matching `./internal/anchors` test. The mirror reds.
- Swap one repointed line of `.agents/commands/bench-review-implementation.md` back to its base text with `--check docs-currency-workflow`. The check reds with that row's diagnostic. Run `bench worktree build <target>` first.

For GW08, run `rg -F "BENCH.md"` over `.agents`, `docs`, and `.bench/BENCH-reference.md`, and read each hit against the list above.
Each probe reports a restored subject.
Then run `bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench`.
Also run the checks `docs-currency-workflow`, `prose-mechanics`, `guidance-prose-budgets`, and `root-conformance`.
