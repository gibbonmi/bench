# Agent performance scorecards

These working documents guide model and effort routing from observed landing
quality, cost, and coordination load. They record current aggregates, not a run
log. Provider-specific evidence lives in `open-ai-models.md` and
`claude-models.md`.

## Update contract

- During every implementation retro, read both provider scorecards and refresh
  each provider whose models served as implementer, reviewer, or orchestrator.
  Leave an uninvolved provider unchanged.
- Replace the provider's last-incorporated-landing line and rewrite affected rows
  in place.
- Write `observed quality` as one or two sentences that restate this landing's
  summary for that model, effort, and role. Replace the previous text; do not
  append to it.
- Rewrite `## Current decisions` in full on every run. Keep a decision that still
  steers routing, drop one a later run or the tree superseded, and add one this
  landing earned.
- Write every sentence in ASD-STE100; the reference is
  `.agents/skills/bench-craft-spec/references/ste-prose.md`.
- Aggregate at most the latest 10 comparable assignments per model/effort/role.
- Keep at most 6 routing rows and 5 representative-evidence rows per provider.
- Replace the least decision-relevant evidence when a stronger example arrives.
- Keep each provider file at or below 120 lines; shorten or remove evidence before
  adding a section.
- Attribute rework to one origin: delegate, spec/ticket, tree/tooling, reviewer,
  or orchestrator. Do not charge upstream omissions to the implementer.
- Compare dollars separately from tokens and latency. Unknown measurements stay
  unknown; qualitative impressions never become invented token counts.

Good: fold a landing into its model/effort/role aggregate and replace weaker
evidence. Bad: append a dated narrative or add one row per run.

## Measures

| measure | meaning |
| --- | --- |
| first-pass accepted | Done-claim needed no coordinator-authored correction before commit |
| coordinator catches | Material completeness, fence, correctness, or evidence misses |
| repair rounds | Returned write pass after an accepted finding; in-pass feedback is still a catch |
| mutation quality | Required mutations bit and production was restored exactly |
| terminal quality | Focused checks, full gate, and exact-tip review outcome |
| efficiency | Token evidence when available, relative dollar input, and wall-clock churn |
