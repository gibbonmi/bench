---
name: craft-cli
description: Design standards for building CLI tools that agents drive through the shell — TOON output, minimal schemas, structured errors, ambient context. Use whenever building, modifying, or reviewing any agent-facing CLI whose project declares AXI conformance, or when deciding whether a new agent-facing CLI should conform.
index: building an agent-facing CLI
---

# AXI — Agent eXperience Interface

AXI is the output contract for a CLI an agent drives through the shell. It aims
for accurate action at low token cost. Apply it per surface: a declared query may
conform while operational siblings keep their own documented contract. Full spec:
https://axi.md

## The principles

1. **Token-efficient output.** Emit compact [TOON](https://toonformat.dev) on
   stdout. Keep internal representations private and convert once at the output
   boundary.

2. **Minimal default schemas.** Return only fields needed for the usual decision,
   normally three or four per list row. Put long bodies in a detail view and set
   the default limit high enough for the common case.

   One contrast shows the target:

   ```text
   issues[2]{id,title,status}:
     41,Fix login redirect,open
     42,Rate-limit webhooks,closed
   help[0]{cmd,why}:
   ```

   Good: the complete answer and its terminal state are visible at once.

   ```json
   [{"id":41,"iid":41,"project_id":7,"title":"Fix login redirect",
     "state":"open","labels":[],"author":{"name":"..."},"web_url":"..."}]
   ```

   Bad: unrelated fields multiply by every row before the agent can decide.

3. **Content truncation.** Preview large values with their total size and an
   explicit way to request the complete value. Never silently omit content or
   dump an unbounded body by default.

4. **Pre-computed aggregates.** Include cheap counts and derived verdicts such as
   `30 of 847` or `checks: 3/3 passed`; do not make the caller reconstruct them.

5. **Definitive empty states.** Emit a typed zero-row result. Silence is not an
   answer, and a successful absence exits zero.

6. **Structured errors and honest exit codes.** Put actionable, structured
   refusals on stdout, never dependency traces. Use exit 0 for success and
   idempotent no-ops, 1 for unsatisfied intent, and 2 for usage. Reserve stderr
   for progress or debugging; never prompt interactively.

7. **Ambient context.** Let a session-start hook print a compact current-state
   dashboard so a fresh agent can act without a discovery call.

8. **Content first.** Start stdout with the requested data, empty state, refusal,
   or usage. Append navigation only after that answer.

9. **Contextual disclosure.** Append actions earned by the returned state, never
   a static command catalog.

10. **Consistent way to get help.** All approved queries accept `--help`, `-h`,
    and bare `help`; each prints usage on stdout and exits zero.

## Bench application

This table is the complete approved Bench AXI query set and says how the contract
applies to each surface. Commands not listed retain their own contracts.

| approved query | contextual disclosure |
| --- | --- |
| `bench anchors` | Report anchor matches or a definitive empty result; terminal checks offer no repair busywork. |
| `bench learnings` | Offer the learnings drain for each distinct open entry. |
| `bench maps` | Offer shaping or template repair for each decision state that needs it. |
| `bench guards` | Offer the applicable repair for each stale or unwired guard. |
| `bench diff` | Offer full inspection or retry only when snapshot state warrants it. |
| `bench coverage` | Successful default extraction offers one check action per mapped coverage row: `bench coverage --check <spec>`. For repairable mapped rows, that exact command is the retry with why `retry after repairing coverage map`. `coverage --check` and every refusal retain their error contracts and append no disclosure. |
| `bench roadmap` | Default index omits bodies; request selected complete rows with `bench roadmap --context --row <ID,...>` or the complete snapshot with `bench roadmap --context --full`. |
| `bench worktree list` | Offer inspect, execute, or clean actions according to each worktree state. |

Every approved result ends with `help[N]{cmd,why}:`; an honest empty envelope is
`help[0]{cmd,why}:`. Derive one state-derived action per matching row. Collapse
exact duplicates while retaining stable source order. Carry every known argument
literally, with placeholders only for unknown future-input slots. Terminal results
offer no busywork.

## Conformance

The project gate derives the approved set from the production command registry
and compares both membership directions with the table above. It also grades the
ten ordered principles, output envelopes, help spellings, and executable behavior.
