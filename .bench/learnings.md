# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")


## 2026-07-27 — the Next-command field is stated twice and enforced zero times  [open]

- **What happened:** I closed `/bench-write-spec` by writing `## Next command` as
  a paragraph — the invocation plus "in a fresh mid-tier session", the
  interactive-not-`shift` rationale, reading order, and a `bench gate pin` note.
  The reviewer caught it and said it is not the first time. The field is supposed
  to hold the exact harness-native invocation and nothing else; imperative prose
  there is aimed at a session that has already started, so it can only misfire —
  a cold session can read "start a fresh session" as work to do.
- **Right behavior:** `## Next command` holds one backticked invocation. Rationale,
  tier, and venue belong in `## State`; a gotcha scoped to one build belongs in
  that spec's testing decisions, per the file's own Shape section.
- **Proposed rule change:** none to the prose — the rule is already written in two
  places (`AGENTS.md`'s Phase-close handoff paragraph: "the exact harness-native
  next command", and the Shape section `internal/handoff/text.go` writes: "the
  exact harness-native invocation, not a description of it"). A third statement
  would be a third derivation of a fact the kit already carries twice and
  enforces nowhere, which is the defect this repo's code standard names. What is
  missing is authority. Codification candidate: grade the field — the `## Next
  command` body must be a single line holding exactly one backticked invocation,
  with prose refused. Two possible owners, reviewer's call: a conformance check,
  which makes it a permanent tripwire on any root carrying the file, or
  validation inside `bench handoff` at write time, which catches it earlier but
  only on the surface that writes through the CLI. The conformance check is the
  one I would pick, because the field is most often hand-written.
