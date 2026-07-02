---
name: craft-cli
description: Design standards for building CLI tools that agents drive through the shell — TOON output, minimal schemas, structured errors, ambient context. Use whenever building, modifying, or reviewing any agent-facing CLI whose project declares AXI conformance. Not for tools with their own declared output contract (bench itself: plain text, stderr errors, documented exit codes). Reach for this on any command output, flag, error message, or exit code in an AXI-conformant tool.
---

# AXI — Agent eXperience Interface

Standards for a CLI an agent uses by shell execution. The goal is higher accuracy
at lower token cost than either a raw CLI or an MCP server. In a project that
declares AXI conformance, treat these as the conformance target the gate checks.
A tool with its own documented output contract (bench's plain-text
stderr/exit-code contract, for one) is out of this skill's scope; do not "fix"
such a tool toward these rules. Full spec: https://axi.md

## The principles

1. **Token-efficient output.** Emit [TOON](https://toonformat.dev) on stdout
   (~40% fewer tokens than equivalent JSON). Keep internal logic on JSON; convert
   at the output boundary only.

2. **Minimal default schemas.** Every field costs tokens times row count. Default
   list views to 3–4 fields (id, title, status), not ten. Default limits high
   enough to cover the common case in one call. Long bodies go in detail views.
   Offer `--fields` for explicit extras.

3. **Truncate, don't omit.** In detail views, show a preview of large fields plus
   the total size and the escape hatch (`--full`). Never silently drop a field;
   never dump the whole thing by default.

4. **Pre-compute aggregates.** The expensive cost is the follow-up call. Include
   total counts ("30 of 847"), and cheap derived status ("checks: 3/3 passed")
   inline, so the agent doesn't have to ask again.

5. **Definitive empty states.** Say "0 closed issues found," not silence. Make it
   clear the command succeeded and the absence is the answer.

6. **Structured errors and honest exit codes.** Errors go to **stdout** in the
   same structured format, with an actionable suggestion that references your
   CLI's own commands — never a leaked dependency stack trace. Mutations are
   idempotent (closing an already-closed thing is a no-op, exit 0). Reserve
   nonzero for genuinely unsatisfiable intent. No interactive prompts — every
   operation completes from flags alone. stdout = data the agent reads; stderr =
   progress/debug it ignores; exit codes: 0 ok (incl. no-op), 1 error, 2 usage.

7. **Ambient context.** Offer a session-start hook that runs the tool and prints a
   compact dashboard of current state, so a fresh session can act without a
   discovery round-trip.

## Conformance is a gate check

In an AXI-conformant project these aren't style preferences — they're testable.
The project gate asserts TOON-shaped stdout, minimal default schemas, structured
stdout errors, and correct exit codes, and can run a paired-delta harness against
the raw tool the CLI wraps. A change that regresses ergonomics fails the gate the
same as a broken test. That's the external oracle pointed at the thing being
built.
