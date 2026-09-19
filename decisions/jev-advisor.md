# jev advisor: typed answers for routine Bench decisions

Status: shaping

## Destination

This map decides whether Bench uses the jev decision service, and where.
jev is a hosted TypeSafe model. A caller sends state and typed questions, and
jev returns Choice, Score, and Noul answers with probabilities. jev generates
no text.

The map holds two outcomes. The first outcome is an inventory of the routine
decisions in the Bench workflow, with a fit verdict for each one. The second
outcome is the set of candidates that a first spec builds. The compaction
advisor is the pilot candidate.

## Notes

Domain: the Bench workflow decisions and the jev service. Charge
`bench-craft-domain` for each session. Charge `bench-craft-research` for the
Research tickets and the `prototype` skill for the Prototype ticket.

The terms stay apart:

- A **jev call** is one request from Bench or a harness hook to a jev
  endpoint. Do not write "LLM call" or "model call" for it.
- The **hosted jev** is the TypeSafe service. **local-jev** is the planned
  loopback server with the same request shape and a different model. Do not
  write "jev" alone where the endpoint changes the meaning.
- A **candidate** is one routine Bench decision that the inventory examines.
- The **paper screen** is the Research verdict on a candidate. The
  **replay probe** is a Prototype run of a candidate against recorded Bench
  decisions. Do not write "test" or "evaluation" for one of them.
- The **pilot** is the compaction advisor candidate.
- A **pre-answer** is a jev answer on a reviewer decision that the reviewer
  confirms or rejects.

Two standing contracts bind each ticket. `DATA_HANDLING.md` states that Bench
uploads no prompt, objective, or environment data. The dependency standard in
`AGENTS.md` excludes a runtime service unless the reviewer decides otherwise.

A map-owned asset stays in the map's assets folder,
decisions/jev-advisor/assets/.

## Decisions so far

- [Which egress posture does a jev call have?](jev-advisor/tickets/1.md): a jev call occurs only with a committed opt-in plus a key, and `BENCH_OFFLINE` suppresses it.
- [Which authority ceilings hold for every candidate?](jev-advisor/tickets/2.md): a jev answer never changes a gate verdict, and a reviewer decision gets at most a pre-answer.
- [Which decisions does the inventory cover?](jev-advisor/tickets/3.md): agent judgment calls, heuristics in the Bench executable, and reviewer decisions.
- [Does the compaction advisor stay in this map?](jev-advisor/tickets/4.md): it stays as the pilot.
- [What evidence proves that a candidate fits jev?](jev-advisor/tickets/5.md): the paper screen first, then a replay probe before a spec includes the candidate.
- [How does local-jev enter the map?](jev-advisor/tickets/15.md): ticket #1 stays, and local-jev is the second arm of the replay probe.

## Not yet specified

- The reviewer stated that Bench took ideas from jev before. The tree names
  no jev, TypeSafe, or System One concept. The map must name those ideas so
  that one concept does not get two names.
- compact-adviser is hint-only on Codex and automatic on Claude Code. The
  map does not yet state the harness parity that Bench requires.

## Spec-writer discretion

## Out of scope

- A jev answer that changes or replaces a gate verdict.
- A jev answer that closes a reviewer decision without the reviewer's
  confirmation.
- A jev call in a default install with no opt-in.
- Text generation by jev. The service does not generate text.

## Sources

- URL: https://docs.typesafe.ai/introduction
  Supports: the destination: the three question types and the no-text limit.
  Drift: a new question type or a change to the answer fields.
- URL: https://docs.typesafe.ai/api.md
  Supports: ticket #1: the endpoint is hosted and needs a bearer key.
  Drift: a self-hosted or local deployment option.
- URL: https://docs.typesafe.ai/model-jaggedness/jev-1.13.md
  Supports: tickets #5 and #7: the limitation list that the paper screen applies.
  Drift: a newer model limitation page.
- URL: https://github.com/kunchenguid/compact-adviser
  Supports: tickets #4 and #9: the pilot mechanism, its inputs, and its harness modes.
  Drift: a change to the two questions, the threshold curve, or the harness list.
- URL: https://github.com/kunchenguid/local-jev/blob/main/IMPLEMENTATION-PLAN.md
  Supports: tickets #15, #16, and #17: the loopback endpoint, the ModernBERT model, and the omitted `confidence` field.
  Drift: a first release, a second commit, or a change to the answer fields.
- Path: `DATA_HANDLING.md`
  Supports: ticket #1: the Network section lists three egress cases and promises no content upload.
  Drift: a change to the Network section.
- Path: `AGENTS.md`
  Supports: tickets #1 and #11: the dependency standard excludes a runtime service by default.
  Drift: a change to the dependency standard paragraph.
- Path: `.bench/BENCH.md`
  Supports: ticket #2: invariant 1 and the Roles section fix the two ceilings.
  Drift: a change to invariant 1 or to the Roles section.
