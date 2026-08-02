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

## 2026-08-01 — light-path ticket receipts have no terminal lifecycle  [open]

- **What happened:** Five shipped light-path changes still have ticket-only
  folders under `specs/`, but no `spec.md`. Their implementation commits landed,
  while `bench spec history` returns no record for their slugs and `bench spec
  retire` cannot target them. The folders therefore persist as committed receipts
  without being classified as active specs, retained history, or retireable
  lifecycle state.
- **Right behavior:** A light-path close gives its ticket receipt one explicit
  terminal disposition. If receipts remain under `specs/`, a canonical reader
  identifies why they remain and distinguishes them from active specs; otherwise
  the close removes them after promoting any durable content. Sessions should not
  have to infer lifecycle state from the absence of `spec.md`.
- **Proposed rule change:** Make the light-path final-check path either retain and
  index ticket-only receipts through one existing spec reader, or remove their
  folders at close. Choose one owner and one policy rather than adding a second
  ad-hoc archive convention.
