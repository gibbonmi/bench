# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-integrate-learnings` reviews the open entries, promotes the
generalizable ones into the kit with sign-off, and prunes them: a resolved entry
leaves this file, and its verdict (promoted or dismissed, one line of why) is
recorded in the integration commit and CHANGELOG. The journal holds open entries
only; history lives in git. Never rewrite a kit rule yourself — that is the whole
point of capturing here instead.

Format per entry:

## <date> — <short title>  [open]
- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

An entry leaves this file only via /bench-integrate-learnings.

<!-- entries below -->

## 2026-07-03 — session-start stale gate confuses without a benign/real split  [open]
- **What happened:** Reviewer flagged that the gate "is almost always stale in a
  new session," which reads as alarming. Diagnosis: the verdict is content-addressed,
  so it goes stale the instant the tree moves past the last green — and sessions
  routinely end with a change after the last gate run (a manual commit that wasn't
  re-gated, or a `bench idea` park that dirties ROADMAP.md). So new sessions almost
  always open stale, but the drift is often benign (capture-scratch like ROADMAP.md /
  .bench-notes.md, which no gate check reads) rather than unverified code.
- **Right behavior:** At session start, when the gate reads stale, tell the user
  *why* — split "benign drift only (e.g. a parked idea) → just a reminder to re-run
  the gate" from "committed code moved since the last green → real, re-run before you
  trust it." The bare word "stale → re-run the gate" hides that distinction and reads
  as an error even when it's harmless.
- **Proposed rule change:** Consider having `bench status` classify a stale verdict:
  if the diff from the gated tree is confined to capture-scratch paths, word it as a
  benign reminder; otherwise flag it as real. (Distinct from, and a lighter-weight
  alternative to, the parked idea of carving capture-scratch out of `gate_tree_hash`
  entirely — that changes the oracle's key and is sensitive on the tripwire branch.)
