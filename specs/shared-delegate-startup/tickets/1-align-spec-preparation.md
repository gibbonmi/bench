# Align spec preparation with shared rules

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/SKILL.md, .agents/skills/bench-craft-delegate/references/shared-preparation.md (new), .agents/skills/bench-craft-line/SKILL.md, .agents/commands/bench-write-spec.md, .bench/BENCH-reference.md, CONTEXT.md, reviews/shared-delegate-startup.md (new), tests/canary/claude-agent-definitions/agent-unnamed-in-skill, tests/canary/claude-agent-definitions/skill-names-missing-agent, tests/canary/docs-currency-token-diet/signal-vocabulary-drift, tests/canary/skills-index-command-adapters/adapter-inert-invocation-key, tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy, tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted, tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary, tests/canary/workflow-guidance-anchors/context-coverage-map-term, tests/canary/workflow-guidance-anchors/context-coverage-row-parts, tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary, tests/canary/workflow-guidance-anchors/context-decision-map-term, tests/canary/workflow-guidance-anchors/context-reader-sweep-term, tests/canary/workflow-guidance-anchors/context-ticket-vocabulary, tests/canary/workflow-guidance-anchors/decision-map-asset-path, tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package, tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap, tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge, tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green, tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller, tests/canary/workflow-guidance-anchors/delegate-model-id-escalation, tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface, tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance, tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents, tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row, tests/canary/workflow-guidance-anchors/delegated-author-limit, tests/canary/workflow-guidance-anchors/delegated-mid-review-route, tests/canary/workflow-guidance-anchors/delegated-per-ticket-author, tests/canary/workflow-guidance-anchors/delegated-tier-authorization, tests/canary/workflow-guidance-anchors/delegated-unbound-model-stop, tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route, tests/canary/workflow-guidance-anchors/prepared-triage-bounds, tests/canary/workflow-guidance-anchors/reference-agent-push-rule, tests/canary/workflow-guidance-anchors/story-line-anchor-missing, tests/canary/workflow-guidance-anchors/write-spec-artifact-authorization, tests/canary/workflow-guidance-anchors/write-spec-authorization-boundary, tests/canary/workflow-guidance-anchors/write-spec-branch-split, tests/canary/workflow-guidance-anchors/write-spec-conversation-authorization, tests/canary/workflow-guidance-anchors/write-spec-conversation-fork, tests/canary/workflow-guidance-anchors/write-spec-decision-source, tests/canary/workflow-guidance-anchors/write-spec-late-uncertainty, tests/canary/workflow-guidance-anchors/write-spec-map-sources, tests/canary/workflow-guidance-anchors/write-spec-phase-ownership, tests/canary/workflow-guidance-anchors/write-spec-reader-sweep-sequence, tests/canary/workflow-guidance-anchors/write-spec-review-made-conditional, tests/canary/workflow-guidance-anchors/write-spec-review-tier-escalated, tests/canary/workflow-guidance-anchors/write-spec-slicing-sign-off, tests/canary/workflow-guidance-anchors/write-spec-ticket-approval-fields, tests/canary/workflow-guidance-anchors/write-spec-ticket-breakdown-charge, tests/canary/workflow-guidance-anchors/write-spec-verification-learning, tests/canary/workflow-guidance-anchors/write-spec-verification-learning-contents, tests/canary/workflow-guidance-anchors/write-spec-verification-log, tests/canary/claude-agent-definitions/model-declared, tests/canary/claude-agent-definitions/name-mismatch, tests/canary/claude-agent-definitions/shell-tool-absent, tests/canary/claude-agent-definitions/spawning-tool, tests/canary/claude-agent-definitions/tools-absent, tests/canary/docs-currency-token-diet/benchref-imported, tests/canary/docs-currency-token-diet/benchref-pointer-dropped, tests/canary/docs-currency-token-diet/benchref-section-duplicated, tests/canary/skills-index-command-adapters/dangling-index, tests/canary/skills-index-command-adapters/missing-index-field, tests/canary/skills-index-command-adapters/stale-index-wording, tests/canary/skills-index-command-adapters/unindexed-skill, tests/canary/workflow-guidance-anchors/acceptance-coverage-anchor, tests/canary/workflow-guidance-anchors/command-handoff-anchor, tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer, tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor, tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor, tests/canary/workflow-guidance-anchors/edge-inventory-anchor, tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor, tests/canary/workflow-guidance-anchors/reference-bench-operational-layer, tests/canary/workflow-guidance-anchors/reference-category-context, tests/canary/workflow-guidance-anchors/reference-category-oracle, tests/canary/workflow-guidance-anchors/reference-category-setup, tests/canary/workflow-guidance-anchors/reference-category-work, tests/canary/workflow-guidance-anchors/reference-gate-authority, tests/canary/workflow-guidance-anchors/reference-kit-only-ship, tests/canary/workflow-guidance-anchors/reference-no-path-fallback, tests/canary/workflow-guidance-anchors/reference-progressive-loading-term, tests/canary/workflow-guidance-anchors/reference-refusal-route-shape, tests/canary/workflow-guidance-anchors/reference-retro-capture-owner, tests/canary/workflow-guidance-anchors/reference-retro-drain-owner, tests/canary/workflow-guidance-anchors/reference-skills-guidance, tests/canary/workflow-guidance-anchors/reference-upgrade-route, tests/canary/workflow-guidance-anchors/shared-worktree-path-pin, tests/canary/workflow-guidance-anchors/slicing-in-write-spec, tests/canary/workflow-guidance-anchors/spec-retire-roadmap-row, tests/canary/workflow-guidance-anchors/staged-command-sweep-anchor, tests/canary/workflow-guidance-anchors/ticket-decision-map-lifecycle-anchor, tests/canary/workflow-guidance-anchors/ticket-stage-routing-anchor, tests/canary/workflow-guidance-anchors/write-spec-frozen-base-and-tip-review, tests/canary/workflow-guidance-anchors/write-spec-handoff-anchor, tests/canary/workflow-guidance-anchors/write-spec-post-slicing-handoff, tests/canary/workflow-guidance-anchors/write-spec-ready-map-authorization
Covers: SP2, SP3, SP4, SP5, SP6, SP7, SP8, SP9, SP10, SP11, SP12, SP13, SP16, SP18, SP19, SP20, SP35, SP36, SP37, SP38, SP39, SP42, SP45, SP46, SP47, SP48, SP49, SP50, SP51, SP52

## What to build

Complete the optional route from an explicit instruction through the existing spec author fork and cost evidence.
Create the shared procedure in the focused delegation reference.
Link its owner from the delegation skill and spec phase.
Qualify the fresh-delegate charge advice so it cannot contradict inherited context.
Keep one approved-context author for the spec and tickets by default.

Permit the invoking agent to select a sequential ticket fork for substantial distinct spec and ticket jobs.
Require its stated reason, without another user permission turn or preparation session.
The second fork starts after the completed spec returns and the first writer stops.
It receives the completed spec and approved source with the inherited invoking line and existing assignment isolation.
If the second fork is unsuitable or unavailable, report the fallback and resume the original author for tickets.
One independent review covers the completed pair.

Keep line selection in the line skill and cost provenance in the existing assessment reference.
Define the four proposed terms in the glossary.
The shared reference links to those owners instead of restating their policy.
The next tickets connect implementation and review consumers to this same procedure.
Do not advertise combined adoption until both consumers are complete.

Before edits, compare every owned case with the current source.
Record the old and new results in the review pickup when real evidence exists.
Keep raw native observations separate from source-based expectations.
Preserve the pinned fixtures co-named in Writes.
Reclaim required headroom within each changed guidance file in this ticket.
Do not raise its budget or move its anchored duties.

## Acceptance

- [ ] Explicit spec activation uses the existing approved-context author fork and shared procedure.
- [ ] Standalone spec activation preserves existing approvals, the ticket graph, and the normal route when absent.
- [ ] Eligibility uses the active native contract, required context, inherited line, isolation, and parent capability.
- [ ] Full forks omit conflicting model or effort overrides and report any permitted fresh-delegate fallback.
- [ ] A changed line requires a documented lower total with preparation, delegates, output, and expected repairs.
- [ ] Unknown required inputs and tied or higher totals retain the existing line.
- [ ] Each authorized spec-source kind retains the existing author fork and capable-session handoff.
- [ ] The default route retains one author, and neither route adds a preparation session.
- [ ] The invoking agent can select the sequential option for substantial distinct jobs with a stated reason and no further permission turn.
- [ ] The second writer starts after the completed spec returns and the first writer stops.
- [ ] The ticket fork receives the completed spec and approved source with the inherited line and existing assignment isolation.
- [ ] The estimate includes both authors, preparation, output, and expected repairs.
- [ ] An unsuitable or unavailable second fork causes a reported return to the original author for tickets.
- [ ] One independent review covers the completed pair.
- [ ] Existing assessment records distinguish estimates, actual charges, inherited baselines, and unknown child measurements.
- [ ] Successful dispatch and fewer reads alone produce no measured savings claim.
- [ ] Live compatibility claims require observed context, line, and assignment evidence.
- [ ] Owned semantic cases have cited old/new evidence, with unresolved cases reported rather than marked passed.

Run the existing workflow and guidance-budget checks, then prose checks for changed Markdown.
The semantic cases remain review-owned; those commands do not establish native runtime behavior.
