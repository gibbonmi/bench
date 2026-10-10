# Fork the spec author after a research fan-out

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, projects/benchkit.md, docs/field-guide.html, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_craft_research.go, internal/anchors/registry_craft_research_test.go, tests/canary/workflow-guidance-anchors/write-spec-conversation-fork/MUTATE.json, tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence/MUTATE.json, CHANGELOG.md
Covers: none

## What to build

The reviewer directed on 2026-10-10 that the `/bench-write-spec` fork happens after a fan-out by research delegates. Today `## Who runs this phase` says that the invoking session immediately forks the current conversation. The fork then does the evidence reads alone in step 1.

`## Who runs this phase` gets the fan-out first. The invoking session charges read-only research delegates over the one decision source, under `craft-delegate`, `craft-line`, and `craft-research`. `craft-research` owns the round, the number of readers, the citations, and the not-verified list, and the command restates none of these rules. The delegates read the enforcement surface that the rows touch, `docs/known-issues.md` when it exists, and `craft-spec`'s reader sweep. They also read the code sites, tests, fixtures, and anchors of each decision-source item. The research delegates write nothing.

After every research delegate returns, the invoking session forks. The charge gives the fork the findings as compressed inputs, and the fork does not repeat their reads. The fork keeps every spec-and-ticket write in the phase worktree. The invoking session coordinates the phase and presents the approval table. The line rules for the fork stay.

Step 1 `Author` consumes the findings. It reads only what the findings leave open, before the row lock. It keeps `bench coverage --check`, the check floor, and every other obligation.

The reviewer changed the fork rule and the step-1 reader sweep rule, so their anchors change with them. The fork anchor pins the new fork sentence, and its canary mutates that sentence. The reader-sweep anchor pins the new step-1 sentence, and its canary mutates that sentence. A new require anchor pins the fan-out sentence. An anchor harness test removes each of the three sentences and expects its own diagnostic.

`projects/benchkit.md` and `docs/field-guide.html` restate the fork timing, so they change in the same ticket. The `guidance-prose-budgets` row for the command moves to the new line count only if the command grows.

## Acceptance

- [ ] `## Who runs this phase` charges read-only research delegates over the decision source under `craft-delegate`, `craft-line`, and `craft-research` before the fork.
- [ ] The fan-out covers the enforcement surface, `docs/known-issues.md`, `craft-spec`'s reader sweep, and the code sites, tests, fixtures, and anchors of each decision-source item.
- [ ] The command states that the research delegates write nothing and restates no `craft-research` rule.
- [ ] The fork starts only after every research delegate returns, and it takes the findings as compressed inputs without a repeat of their reads.
- [ ] Step 1 consumes the findings, reads only what they leave open, and keeps `bench coverage --check` and the check floor.
- [ ] `projects/benchkit.md` and `docs/field-guide.html` state the fork after the research fan-out.
- [ ] The anchor harness goes red when the fan-out sentence, the fork sentence, or the step-1 sentence is removed.
- [ ] `CHANGELOG.md` records the change.
