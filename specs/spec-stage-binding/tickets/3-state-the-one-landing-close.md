# State the one-landing close in the guidance

Blocked by: 1-bind-a-new-deliverable-from-its-publishing-tree.md, 2-name-the-binding-state-in-the-handoff.md
Writes: .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, .bench/BENCH.md, AGENTS.md, CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, cmd/bench/testdata/anchors/fixture-repo/AGENTS.md, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_debug_loop.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/commitment_guidance_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: SB31, SB32, SB33, SB34, SB35, SB36, SB37, SB38, SB39, SB40, SB41, SB42, SB43

## What to build

A spec author reads one close order in the Exit handoff of `/bench-write-spec`.
After sign-off, the close commits the binding beside the spec, lands both in one spec-less landing, and passes the start command to the handoff.
`.bench/BENCH.md` lets a spec-stage close bind its spec on the reviewer's sign-off.
`AGENTS.md` pins the binding state in a phase-close handoff and runs the spec-stage handoff from the primary checkout.
The slicing checks put `.bench/commitment.json` in the fence and in one ticket's `Writes:` line.
The `craft-spec` template names the binding that its `Roadmap:` line feeds.

This ticket reads these contracts from tickets 1 and 2:

- Ticket 1: one spec-less landing of a planning worktree publishes a staged spec and its approved binding.
- Ticket 1: plan and approve in that worktree bind a staged spec that the worktree HEAD holds and main lacks.
- Ticket 1: after that landing, `bench commitment show` prints the start command with the bound spec path. `bench status` keeps the same command (`internal/status/delivery.go:19-41`).
- Ticket 1: start, readiness, and the outlook read only the default-branch policy.
- Ticket 2: `bench handoff` writes one `Spec binding: <value>` line directly after each `Spec status` line. `handoffdoc.LabelSpecBinding` is the one spelling of the label.
- Ticket 2: the value reads the default-branch policy, so a spec-stage handoff shows the binding only after the landing puts it on main.

The guidance describes the shipped behavior of both tickets, so this ticket waits for both checkpoints.
All three tickets also name the five command registry files.

**The close order.** Edit the Exit handoff of `.agents/commands/bench-write-spec.md`.
Keep its paragraph at lines 15-16 byte-identical, because it already holds five of its six sentences.
Put this block after that paragraph and before `## Entry contract`, with one blank line on each side.
The block is one intro line and six steps, with no blank line between them.

```markdown
After sign-off, close the spec stage in this order, so one spec-less landing carries the staged spec and its binding.
1. Write the verification log into the spec and the scorecard updates, and commit them.
2. After the last spec edit, copy the spec's `bench commitment inventory` row into one new deliverable binding of the current policy. Put it on the outcome whose sources hold the row's scope, with that scope as the obligation. A spec with no `Roadmap:` line binds to the outcome that `bench status` names.
3. Run `bench commitment plan` and `bench commitment approve` with the sign-off as the decision, and commit the binding.
4. Run `bench preflight build <slug> --plan-only` at that tip.
5. Land spec-less once.
6. From the primary checkout, run `bench handoff --next` with the start command that `bench status` prints.
```

The intro line carries SB31, step 2 carries SB32, and step 6 carries SB33.
The six steps are the steps of the spec section "The close order".
Step 2 adds "After the last spec edit," and the no-`Roadmap:` sentence, because SB32 and that spec section state both.
A list marker starts a new paragraph for `prose-mechanics` (`internal/prose/parse.go:181-183`), so each step is its own paragraph.
The block keeps the staged-spec sentence of line 16, which `internal/anchors/registry_commitment.go:42` pins (SB43).

**The prose budgets.** `.agents/commands/bench-write-spec.md` holds 87 of its 87 lines (`projects/benchkit.md:492`).
The block adds eight lines: the intro line, six steps, and one blank line.
Pay for them with these joins of the base file.
Each join makes one physical line with one space at each seam, and it changes no word.

- Lines 9-11 free two lines.
- Lines 20-21 free one line.
- Lines 22-23 free one line.
- Lines 28-29 free one line, and the join drops the indent of line 29.
- Lines 51-52 free one line.
- Lines 63-65 free two lines, and the join drops the indent of lines 64 and 65.

Do not join lines 21-22, 32-35, 50-51, 61-62, 68-69, 79-81, or 86-87.
A canary `MUTATE.json` `old` string or `internal/conformance/fixture_bite_test.go:368-369` holds each of those newlines byte-exact.
Anchors collapse whitespace (`internal/anchors/locate.go:177`), so each planned join keeps every anchor.
The command ends at 87 lines or fewer.

`.bench/BENCH.md` holds 192 of 192 lines (`projects/benchkit.md:490`), so its grant joins the paragraph of line 149 and adds no line.
`.agents/skills/bench-craft-spec/SKILL.md` holds 155 of 155 lines (`projects/benchkit.md:495`), so the template comment stays on its one line.
`AGENTS.md` (123 lines) and the slicing checks (73 lines) have no budget row.
Each new sentence still obeys the prose bounds: 25 words for each sentence and six sentences for each paragraph (`internal/prose/prose.go:17-18`).

**The binding grant.** Put this sentence into the paragraph of `.bench/BENCH.md:149`, directly before its binding-only limit sentence:

```markdown
A spec-stage close binds its approved spec the same way and cites my sign-off as the decision.
```

That paragraph grows from three sentences to four.
The sentences of lines 147 and 149 stay byte-identical, so the existing `binding-only limit` mirror stays green (SB35).

**The handoff pin.** Change two sentences of the Phase-close handoff of `AGENTS.md`.
Replace the sentence at lines 68-70 that opens "It must also pin the spec path and status" with this sentence:

```markdown
It must also pin the spec path, the spec status, the binding state, the decisions that stay closed, and the exact harness-native next command.
```

Add this sentence on a new line directly after the drain sentence of line 82:

```markdown
A spec-stage close also runs it from the primary checkout after its landing, because the landing releases the worktree.
```

Keep "A drain runs it from the primary checkout after its landing." byte-identical on its own physical line, because `internal/anchors/registry_retained_workflow.go:100` pins it.
Rewrap only the lines of the changed sentence.
No anchor pins the pin list or the venue, so SB38 and SB39 stay review-owned, and this ticket adds no `AGENTS.md` anchor row.

**The slicing rule.** Add this bullet as the last bullet of `## Slicing rules` in `.agents/skills/bench-craft-tickets/references/slicing-checks.md`:

```markdown
- The ownership fence and exactly one ticket's `Writes:` line name `.bench/commitment.json`. The spec-stage close commits the binding inside the range that the plan-only preflight grades.
```

The text is the proposed text of the spec.
The pinned needles at `internal/anchors/registry_ticket_passes_test.go:14-43` stay byte-identical, so the rule takes a new bullet.
The rule goes under `## Slicing rules`, because no preflight row enforces it (spec D5).

**The template comment.** Replace `.agents/skills/bench-craft-spec/SKILL.md:96` with this one line:

```markdown
Roadmap: <optional — the FT<n> row this spec settles; the spec-stage close binds this spec to the outcome whose sources hold that row, and `bench spec retire` names that row's detail file>
```

The line sits inside the fenced template that opens at line 91, so `prose-mechanics` does not grade it.

**The anchors.** Each new pinned sentence takes one anchor row and one mirror.
Each mirror removes its sentence and requires the anchor diagnostic.
Write each mirror needle and diagnostic by hand, independently of the registry.
The commitment rows use the prefix `CommitmentDiagnosticPrefix`, which is `commitment guidance: ` (`internal/anchors/registry_commitment.go:8`).

| row | registry file | kind and section | needle | diagnostic | mirror |
|---|---|---|---|---|---|
| SB31 | `internal/anchors/registry_commitment.go` | `RequireInSection`, `Exit handoff` | the intro line of the block | `commitment guidance: spec authoring dropped the one-landing close` | a `remove` case in `TestCommitmentGuidance` |
| SB32 | `internal/anchors/registry_commitment.go` | `RequireInSection`, `Exit handoff` | the first sentence of step 2 | `commitment guidance: spec authoring dropped the binding from the inventory row after the last spec edit` | a `remove` case in `TestCommitmentGuidance` |
| SB33 | `internal/anchors/registry_commitment.go` | `RequireInSection`, `Exit handoff` | the sentence of step 6 | `commitment guidance: spec authoring dropped the start command as the handoff next command` | a `remove` case in `TestCommitmentGuidance` |
| SB34 | `internal/anchors/registry_commitment.go` | `RequireInSection`, `Workflow` | the grant sentence | `commitment guidance: operating guide dropped the spec-stage binding grant` | a `remove` case in `TestCommitmentGuidance` |
| SB36 | `internal/anchors/registry_ticket_passes.go` | `RequireInSection`, `Slicing rules` | the first sentence of the slicing rule | `ticket slicing: the fence and one Writes line own the commitment policy` | one `anchorRule` in `TestTicketSlicingPasses` |

Each `remove` case names the live file and the exact live sentence, and that sentence occurs once in the file.
The SB34 case keeps a leading space, as the `binding-only limit` case does (`internal/conformance/commitment_guidance_test.go:56`).
`TestRegistryNeedlesDistinct` refuses a repeated needle for one file and kind (`internal/anchors/anchor_harness_diagnostics_test.go:202-224`).
No current needle repeats these five needles.

**The needles to avoid.** The command edit keeps the two `Require` sentences that `internal/anchors/registry_data.go:336-337` pin on line 16.
The new command text matches none of these:

- "## Handoff" and "map's Handoff" (`internal/anchors/registry_data.go:358-359`).
- The retired sign-off route at `internal/anchors/registry_data.go:335`, which opens "Stop for sign-off, then continue into `/bench-implement-spec`".
- "bench preflight evidence" (`internal/anchors/registry_ft311_preparation.go:50`).
- "chunk ticket review line" (`internal/anchors/registry_ticket_passes.go:49`), "decisions/assets/" (`internal/anchors/registry_decision_maps.go:13`), and the falsification sentence at `internal/anchors/registry_data.go:41`.
- "A fresh session" with a newline and then "builds after ticket approval", and "fresh mid-tier build session" (`internal/conformance/retained_workflow_test.go:231-233`).

The grant matches no `.bench/BENCH.md` forbid row.
Those rows are at `internal/anchors/registry_commitment.go:60-64`, `internal/anchors/registry_retained_workflow.go:88`, `:90`, `:217`, and `:270`.
More rows are at `internal/anchors/registry_ft311_review_dispatch.go:41` and `:49-50`, and at `internal/anchors/registry_data.go:166-167`, `:186`, and `:418-431`.
The row at `internal/anchors/registry_data.go:167` forbids the word "thorough".

The two `AGENTS.md` sentences match none of these needles:

- `anchors.CommitmentRuleMarker`, "**Deliver only the committed outcome.**", and each other shared-rule marker. A copy reds `internal/conformance/commitment_guidance_test.go:112-121` and the check at `internal/conformance/validity_checks_test.go:251-271`.
- "`bench handoff` rewrites only the calling worktree's assignment section." (`handoffSectionRule`, `internal/anchors/registry_data.go:384`).
- "The primary checkout owns the `main` section." (`handoffMainOwner`, `internal/anchors/registry_retained_workflow.go:97`).
- "The verb keeps a non-empty Next command." (`handoffNextRule`, `internal/anchors/registry_retained_workflow.go:98`).
- "refuses a State that pins a commit outside the tip's ancestry" (`handoffAncestryRefusal`, `internal/anchors/registry_retained_workflow.go:99`).
- "one gate-priced commit" (`:101`) and "Give the drafted State to `bench handoff --state-file <path>`." (`:105`) in `internal/anchors/registry_retained_workflow.go`.
- "the four invariants, how the pieces fit" and "the communication rules, and the skills index" (`internal/anchors/registry_retained_workflow.go:93-94`).

The template comment matches no forbid needle of the skill.
Those needles are at `internal/anchors/registry_retained_workflow.go:43`, `:45-46`, and `:149`, at `internal/anchors/registry_debug_loop.go:43-44`, and at `internal/conformance/retained_workflow_test.go:250-251`.

**The changelog entry.** Add one section under `## [Unreleased]` of `CHANGELOG.md`, as its first section.
It holds one `Changed` bullet for the one-landing close and one `Added` bullet for the `Spec binding` line, in the form of the current entries.
This is the proposed text:

```markdown
### Spec stage binding

- Changed the spec-stage close so that one spec-less landing publishes the staged spec and its deliverable binding. `bench commitment plan` and `bench commitment approve` read a new deliverable at the worktree's HEAD, and `bench commit` and `bench worktree land` read it in the candidate tree. `/bench-write-spec` now states the close order, and `.bench/BENCH.md` lets a spec-stage close bind its spec on the reviewer's sign-off. The slicing checks put `.bench/commitment.json` in the fence and in one ticket's `Writes:` line, so the plan-only preflight stays green. After the landing, `bench status` prints the complete `bench commitment start` command with the spec path. Before, the binding needed a second landing and a second gate.
- Added a `Spec binding` line after each `Spec status` line of a `bench handoff` section. The line names each outcome that the default-branch policy binds to the spec, or it reads `unbound`. When the policy does not parse, the line reads `binding unknown (policy unreadable)`, and the handoff still writes. Before, the handoff named no binding, so a handoff could call a spec ready while main held no binding.
```

Keep the `CHANGELOG.md` sentences that `internal/anchors/registry_data.go:334`, `:371`, and `:372` pin, and the two canary `old` strings.

**The red route of each row.** Record each red in the return, with its command and its failure line:

- SB31, SB32, SB33, and SB34: write the guidance sentence and its `remove` case first, and run `bench test --package ./internal/conformance --run TestCommitmentGuidance`. Each new case fails with `mutated <file> raised []` until its anchor row exists. Then add the anchor rows, and record the green run with its `observed red:` log line.
- SB36: add the `anchorRule` to `TestTicketSlicingPasses` before the anchor row, and run `bench test --package ./internal/anchors --run TestTicketSlicingPasses`. The rule fails with `tree broken at`, because no anchor raises its diagnostic. Then add the row.
- SB35 and SB43: the existing `binding-only limit` case (`internal/conformance/commitment_guidance_test.go:56`) and `staged spec` case (`:69`) stay green with no edit. Each `old` string still occurs once in its live file.
- SB41: add the block before the joins, and run `bench test --check guidance-prose-budgets`. `TestGuidanceProseBudgetsHoldOnTheLiveTree` reds at 95 lines against the 87-line budget. Then make the joins, and record the green run.
- SB37: review-owned. Ticket 1's `Writes:` line holds `.bench/commitment.json` today, and this ticket does not name it. The review reads the spec's fence section, and the plan-only preflight `fence-writes` row grades the union.
- SB38 and SB39: review-owned. The review reads the two Phase-close handoff paragraphs of `AGENTS.md`, because no anchor pins them.
- SB40: review-owned. The review reads `.agents/skills/bench-craft-spec/SKILL.md:96`.
- SB42: review-owned. The review reads the five new anchor rows, their five mirrors, and the red log of the author for each one.

After the edit, record one red through `bench probe` for each probe below:

- Omit the SB33 anchor row from `internal/anchors/registry_commitment.go`, with `--package ./internal/conformance --run TestCommitmentGuidance`. The SB33 case fails.
- Omit the SB36 anchor row from `internal/anchors/registry_ticket_passes.go`, with `--package ./internal/anchors --run TestTicketSlicingPasses`. The new rule fails.
- Swap the SB34 needle in `internal/anchors/registry_commitment.go` for a reworded grant, with `--package ./internal/conformance --run TestCommitmentGuidance`. The live `.bench/BENCH.md` then raises the grant diagnostic, and the SB34 case fails.

Each probe reports a restored subject.
Then run `bench test --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance`.
Also run the checks `docs-currency-workflow`, `guidance-prose-budgets`, `prose-mechanics`, `skills-index-command-adapters`, and `root-conformance`.

**The Writes line.** The five guidance files and `CHANGELOG.md` are the edit sites.
Each anchor registry file that names a guidance path joins the line under the anchor closure rule: fourteen files in `internal/anchors`.
`internal/anchors` is a bound package, so the five command registry files join the line.
This ticket does not change those five files, and four of them hold more than 400 lines.
`internal/conformance/commitment_guidance_test.go` holds the commitment mirrors.

Each fixture directory that pins a guidance file joins the line, as `--propose-writes` lists them.
The overlay copies of the guidance files under those directories should need no edit.
Each canary `EXPECT` names one diagnostic, such as `tests/canary/workflow-guidance-anchors/command-handoff-anchor/EXPECT:1`.
No overlay holds the staged-spec sentence that `internal/anchors/registry_commitment.go:42` requires today, and those canaries pass.
If a canary turns red, update its overlay and report it.
The two canary units `changelog-reduced-schema-columns` and `changelog-ticket-vocabulary` pin `CHANGELOG.md`.

`cmd/bench/testdata/anchors/fixture-repo/AGENTS.md` plants each `AGENTS.md` `Require` needle (`cmd/bench/anchor_help_test.go:17-40`).
It changes only if this ticket adds an `AGENTS.md` `Require` row, and the plan adds none.
`internal/conformance/fixture_bite_test.go:368-369` holds two byte-exact command needles, and the planned joins leave them unchanged.

These paths stay off the line:

- `.bench/commitment.json`, because ticket 1 owns it under decision D5, and exactly one ticket names it.
- `.claude/commands/bench-write-spec.md`, because `.claude/commands` is a symbolic link to `../.agents/commands`.
- `internal/conformance/retained_workflow_test.go`, because its list at `:227-262` holds only forbidden phrases, and an edit there weakens a check.
- `.bench/BENCH-reference.md`, because decision D6 gives the pin list to `AGENTS.md`, and the reference names no handoff pin (`:361`).

Keep every file that this ticket writes at 400 lines or fewer.
A file that already holds more than 400 lines takes no growth.
These include `internal/anchors/registry_data.go` (478), `internal/anchors/registry_data_test.go` (1286), and `internal/conformance/fixture_bite_test.go` (872).
The lane runs `bench structure --growth`, and a file over 400 lines that gains a line turns it red.

## Acceptance

- [ ] The write-spec Exit handoff states that one spec-less landing carries the staged spec and its binding.
- [ ] The write-spec Exit handoff binds from the spec's `bench commitment inventory` row after the last spec edit.
- [ ] The write-spec Exit handoff passes the start command that `bench status` prints to `bench handoff --next`.
- [ ] `.bench/BENCH.md` lets a spec-stage close bind its approved spec with the sign-off as the decision.
- [ ] `.bench/BENCH.md` keeps the binding-only limit sentence after the new grant.
- [ ] The slicing checks require the fence and exactly one ticket `Writes:` line to name `.bench/commitment.json`.
- [ ] The ownership fence of this spec and the `Writes:` line of its first ticket name `.bench/commitment.json`.
- [ ] `AGENTS.md` requires the phase-close handoff to pin the binding state.
- [ ] `AGENTS.md` runs the spec-stage handoff from the primary checkout after its landing.
- [ ] The `Roadmap:` comment of the `craft-spec` template states that the close binds the spec to the outcome whose sources hold that row.
- [ ] Each edited guidance file stays within its `guidance-prose-budgets` limit.
- [ ] Each new pinned sentence has one anchor row and one mutation mirror that logs its observed red.
- [ ] The write-spec Exit handoff keeps the planning-only staged-spec sentence.
- [ ] `CHANGELOG.md` holds one typed section under `## [Unreleased]` that states the one-landing close and the `Spec binding` line.
- [ ] `.agents/commands/bench-write-spec.md` holds 87 lines or fewer, `.bench/BENCH.md` 192 or fewer, and `.agents/skills/bench-craft-spec/SKILL.md` 155 or fewer.
