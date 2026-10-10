# Review the spec, then slice tickets by chunk

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, .agents/skills/bench-craft-line/SKILL.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, projects/benchkit.md, docs/field-guide.html, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, internal/anchors/registry_commitment.go, internal/anchors/registry_calibration.go, internal/anchors/registry_calibration_test.go, internal/anchors/registry_decision_maps.go, internal/anchors/registry_decision_maps_test.go, internal/anchors/registry_ft311_preparation.go, internal/anchors/registry_ft311_review_dispatch.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/write-spec-conversation-fork/, tests/canary/workflow-guidance-anchors/write-spec-review-made-conditional/, tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated/, tests/canary/workflow-guidance-anchors/write-spec-verification-learning/, tests/canary/workflow-guidance-anchors/write-spec-ticket-breakdown-charge/, tests/canary/workflow-guidance-anchors/benchkit-review-round-owner/, tests/canary/workflow-guidance-anchors/benchkit-review-round-routing/, tests/canary/workflow-guidance-anchors/benchkit-spec-ownership/, CHANGELOG.md
Covers: none

## What to build

The reviewer directed on 2026-10-10 that `/bench-write-spec` slices tickets one chunk at a time after a spec-only review. Today the fork authors the spec and slices every ticket. One review at reviewer sign-off then covers the spec and the tickets together.

The research fan-out stays first. The fork then authors the spec and its chunk table only. One independent review round grades the spec alone under `craft-spec`'s review rubric, its chunk table included, and returns accept or reject. Slicing starts only after that round accepts.

The invoking session slices the tickets one chunk at a time, in the order of the accepted chunk table. For each chunk, read-only research delegates return each likely fence expansion with its evidence. The invoking session verifies the findings, writes the chunk's tickets under `craft-tickets`, and compares each `Writes:` line with `bench preflight build --propose-writes`. A cheap-tier, high-effort delegate then grades the chunk under the `craft-tickets` quiz and returns accept or reject. A chunk review can run while the next chunk's research runs, but the writes stay serial.

After the last chunk is accepted, the invoking session writes the fence section, and `bench preflight build <slug> --plan-only` grades the whole ticket graph. The spec and the tickets then receive one sign-off.

Each fact has one owner:

- `craft-tickets`' slicing checks own the fence-expansion evidence classes, and the command points to them.
- `craft-line`'s stage defaults own the tiers of the spec review and the chunk ticket review, and the command and `projects/benchkit.md` point to them.
- The command owns the iteration count of the verification log. The count sums the iterations of the spec round and of every chunk round. Each round that takes more than one iteration gets one learning entry.

The direction changes these rules, so their anchors and canaries change with them:

- the fork rule and the review-round rule;
- the learning-entry rule and the ticket-write rule;
- the benchkit authoring and review routes;
- the slicing-check preflight rule.

New require anchors pin the spec-only review, the per-chunk slicing, and the fence-expansion research. They also pin the per-chunk review, its cheap tier, the `--plan-only` check, and the iteration count. `docs/field-guide.html` restates the phase order, so it changes too. The `guidance-prose-budgets` rows move only for the subjects that grow.

## Acceptance

- [ ] The fork authors the spec and its chunk table only, and the invoking session slices the tickets.
- [ ] One independent review round grades the spec alone before any slicing and returns accept or reject.
- [ ] Each chunk's research returns its likely fence expansions in the classes that the slicing checks list, each with evidence.
- [ ] The invoking session writes each chunk's tickets and compares each `Writes:` line with `bench preflight build --propose-writes`.
- [ ] A delegate on `craft-line`'s cheap-tier, high-effort chunk ticket review line accepts or rejects each chunk, and a rejected chunk is repaired and reviewed again alone.
- [ ] After the last chunk, `bench preflight build <slug> --plan-only` grades the whole ticket graph, and one sign-off follows.
- [ ] The verification log states one iteration-count rule over the spec round and every chunk round.
- [ ] The anchor harness goes red when each new key sentence is removed.
- [ ] `CHANGELOG.md` records the change.
