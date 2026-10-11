# Name the binding state in the handoff

Blocked by: 1-bind-a-new-deliverable-from-its-publishing-tree.md
Writes: CHANGELOG.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/commitment/commitmenttest/repo.go, internal/commitment/delivery.go, internal/commitment/delivery_test.go, internal/commitment/repository/outlook.go, internal/conformance/axi_query_registry_test.go, internal/conformance/handoff_single_source_test.go, internal/conformance/subcommand_routing_table_test.go, internal/handoff/facts.go, internal/handoff/facts_test.go, internal/handoff/handoff.go, internal/handoff/render.go, internal/handoff/render_test.go, internal/handoff/sections_test.go, internal/handoff/text.go, internal/handoffdoc/document.go, internal/handoffdoc/document_test.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary
Covers: SB23, SB24, SB25, SB26, SB27, SB28, SB29, SB30, SB44

## What to build

The next session reads the binding state of each live spec in its handoff section.
`bench handoff` writes one `Spec binding` line directly after the `Spec status` line of each live spec.
The value comes from the default-branch policy, so a binding that only a worktree holds reads `unbound`.
The handoff still writes when the policy is absent or unreadable.

The spec term applies here.
The binding state of a spec path is the ids of the outcomes whose undelivered deliverable binding in the default-branch policy names that path.

This ticket reads this contract from ticket 1:

- Start, readiness, and the outlook read only the default-branch policy.
- After one spec-less landing publishes a staged spec and its approved binding, the default-branch policy binds that spec.
- Ticket 1 can add a policy rule beside `Unsettled` in `internal/commitment/delivery.go`. If it does, reuse that rule, and do not restate a settlement rule.

The value of the line is one of these:

- The ids of the binding outcomes, in policy order, with one space between two ids.
- `unbound`, when no outcome binds the path.
- `unbound`, when the default branch holds no policy, or when no default branch resolves.
- `binding unknown (policy unreadable)`, when the policy does not read or does not parse.

Make these edits:

- In `internal/handoffdoc/document.go`, add the label constant `LabelSpecBinding` with the value `Spec binding`. The parser reads any label line as a generic field (`document.go:381`), so the parser takes no other edit.
- In `internal/commitment/delivery.go`, add the rule that maps a path to its binding outcomes, beside `Unsettled`. The rule reuses the settlement of `Unsettled`, so a delivered binding does not count.
- In `internal/commitment/repository/outlook.go`, export one test that reports an unresolved default branch. `errDefaultUnresolved` stays private (`outlook.go:10-11`), and `Store.Outlook` calls the exported test in place of its own `errors.Is` check.
- In `internal/handoff`, read the policy once per run through `commitrepo.Store.Policy`, and give the result to `pins`. The handoff calls the same exported test to find an unresolved default branch.
- Add `binding unknown (policy unreadable)` to the unknown phrases at `internal/handoff/facts.go:44-53`.
- In `internal/handoff/text.go`, name the `Spec binding` label in the Shape text beside `Spec` and `Spec status`.
- In `internal/conformance/handoff_single_source_test.go`, add `handoffdoc.LabelSpecBinding` to the token list of `TestHandoffShapeNamesTheSectionGrammar` (`:171-180`).

The new import edges are `internal/handoff` to `internal/commitment` and to `internal/commitment/repository`.
`internal/handoff` already depends on both packages through `internal/status`, and neither package depends on `internal/handoff`.
`internal/status/handoff.go` reads only the section key and the `Worktree tip` field, so it takes no edit.

`internal/conformance/refusal_route_bypass_test.go:23-46` scans each string literal of the production files in `internal/commitment` and `internal/commitment/repository`.
Each new literal there avoids the route needles `next=`, `next[`, `run bench `, and `; then `, and no literal is the bare label `next`.
Keep the two value phrases in `internal/handoff`, which the scan does not read.

Obey these size limits:

- `internal/handoffdoc/document.go` holds 400 of 400 lines. Pay for the new constant with an equal cut. The stale "six pins" comment at `document.go:35-37` predates the seventh pin, so a shorter comment can free the line.
- `internal/conformance/handoff_single_source_test.go` holds 399 lines, so it takes at most one more line.
- `internal/commitment/repository/repository.go` holds 400 lines and takes no growth. So put the exported test in `outlook.go`, which holds 32 lines.
- Keep every file that this ticket writes at 400 lines or fewer. The lane runs `bench structure --growth`, and a file over 400 lines that gains a line turns it red.

Put SB23, SB24, SB25, and SB26 in the planned `TestSectionRendersTheSpecBinding` in `internal/handoff/render_test.go`.
Put SB27, SB28, and SB44 in the planned `TestSectionRendersAnUnavailableBinding` in the same file.
`TestSectionRendersTheAssignmentPins` and `TestSectionRendersOnePairPerLiveSpec` are the precedents.
Each case builds its repository with `benchRepo` and its worktree with `activeAssignment`, and commits the policy on main through `commitmenttest`.
Each expected outcome id comes from the fixture policy, and each spec path comes from the fixture writes.
No expectation calls the code under test.

These fixtures prove the rows:

- SB23: the policy on main binds the live spec to one outcome. The line after `Spec status` is `Spec binding: <id>`.
- SB24: the policy on main binds a different path. The line reads `Spec binding: unbound`.
- SB25: the policy on main does not bind the spec, and the worktree commits a policy that binds it. The line reads `Spec binding: unbound`.
- SB26: two outcomes of the policy on main bind the spec. The line names both ids in policy order with one space between them.
- SB27: main holds no policy. The line reads `Spec binding: unbound`.
- SB28: main holds a policy file that does not parse, such as an empty file. The handoff exits zero, and the line reads `Spec binding: binding unknown (policy unreadable)`.
- SB44: the primary checkout has two local branches, neither one named `main`, and no `origin/HEAD`, so no default branch resolves (`internal/git/default_branch.go:18-30`). A run with `--next` writes the main section, and the line reads `Spec binding: unbound`, because the outlook reads that state as no policy.

SB29 is `TestHandoffShapeNamesTheSectionGrammar`, and SB30 is `TestSectionRendersOnePairPerLiveSpec`.
`TestSectionRendersOnePairPerLiveSpec` counts `Spec: ` lines, and a `Spec binding: ` line does not match that count.
SB44 is the only test that observes the exported test of `outlook.go`.
If the handoff refuses in the SB44 repository for a reason outside this ticket, stop and report the refusal with its detail.

Record these reds before the production edit:

- SB23, SB24, SB25, SB26, SB27, SB28, and SB44 fail, because the planned tests find no `Spec binding` line.
- SB29 fails after the label constant and its token join the test, and before the Shape edit.
- SB30 passes before the edit and stays green with no edit, so its red comes from a probe.

After the edit, record one red through `bench probe` for each probe below:

- Central probe: make the handoff read the policy file of the assignment worktree in place of the default-branch policy. SB25 fails.
- Write the `Spec binding` line between the `Spec` line and the `Spec status` line. SB30 fails.
- Make the rule return only the first binding outcome. SB26 fails.
- Give `binding unknown (policy unreadable)` when the default branch holds no policy. SB27 fails.
- Return the policy error from the handoff run in place of the phrase. SB28 fails.
- Make the exported test of `outlook.go` report no unresolved default branch. SB44 fails.

Each probe reports a restored subject.

The binding registry binds the five command registry files to `internal/commitment` and `internal/anchors`, so they join this `Writes:` line.
This ticket does not change them, and four of them hold more than 400 lines.

The `Writes:` line names `CHANGELOG.md` as an upper bound, because a research return named it for this chunk.
This ticket does not write the `CHANGELOG.md` entry, because ticket 3 writes the one entry of this spec.
The anchor registry files `internal/anchors/registry_data.go` and `internal/anchors/registry_data_test.go` name `CHANGELOG.md`, and two canary units pin it.
So these four files join the `Writes:` line, and this ticket does not change them.

The other paths of the `Writes:` line are possible sites:

- `internal/handoff/handoff.go` calls `collect` at `:76`, so it can own the one policy read of a run.
- `internal/handoff/render.go` and `internal/handoff/facts_test.go` can change when the policy result joins the facts.
- `internal/handoff/sections_test.go` holds `benchRepo` (`:73-85`), which commits no policy, so a test helper can extend it.
- `internal/commitment/commitmenttest/repo.go` holds `WritePolicy` and `CommitPolicy`, so a fixture helper can join it.
- `internal/commitment/delivery_test.go` can hold a unit test of the new rule.
- `internal/handoffdoc/document_test.go` names "six pins" at `:12`, and its round-trip fixture can carry the new label.

A search for the rendered `Spec status` text finds `text.go`, `facts.go`, `document.go`, `document_test.go`, `render_test.go`, and `handoff_single_source_test.go`.
Each one is on this `Writes:` line.
`AGENTS.md:69` belongs to ticket 3, and the capture note `capture/codex-workflow-companion/README.md:35` is hand-written text, not a rendered output.
This ticket does not name `.bench/commitment.json`, because ticket 1 owns it under decision D5.

This ticket gives ticket 3 this contract:

- `bench handoff` writes one `Spec binding: <value>` line directly after each `Spec status` line, with the values above.
- The value reads the default-branch policy, so a spec-stage handoff shows the binding only after the landing puts it on main.
- `handoffdoc.LabelSpecBinding` is the one spelling of the label, and the Shape text names it.

## Acceptance

- [ ] With a binding on main, `bench handoff` writes `Spec binding: <outcome>` on the line after the spec's `Spec status` line.
- [ ] `bench handoff` writes `Spec binding: unbound` for a live spec that the default-branch policy does not bind.
- [ ] `bench handoff` writes `Spec binding: unbound` when only the worktree's policy binds the spec.
- [ ] `bench handoff` writes both outcome ids in policy order with one space between them when two outcomes bind the spec.
- [ ] `bench handoff` writes `Spec binding: unbound` in a repository whose default branch holds no policy.
- [ ] With an unparseable policy on main, `bench handoff` writes `Spec binding: binding unknown (policy unreadable)`.
- [ ] `bench handoff` writes `Spec binding: unbound` in a repository whose default branch does not resolve.
- [ ] The Shape text names the `Spec binding` label.
- [ ] Each live spec keeps its `Spec` line directly before its `Spec status` line.
- [ ] `internal/handoffdoc/document.go` and `internal/conformance/handoff_single_source_test.go` hold 400 lines or fewer.
