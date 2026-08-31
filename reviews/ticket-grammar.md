# Review pickup — ticket-grammar

Base `5ae3a92e`, reviewed tip `f171f3be`. Raw findings: Standards 7, Spec 8,
Coverage 8. De-duplicated repair targets: 4 (R1–R4). Reviewer-pending items
carry no repair and await the reviewer.

## Standards

Count: 7. Worst: the ticket enumeration is derived in two places, and the two
copies already disagree.

- S1 (high, ask-user → accepted as R1): `internal/preflight/gather.go:330-378`
  and `internal/conformance/ticket_grammar_test.go:108-144` each derive the
  enumeration (`.md`-only, recursion, lstat-first, basename identity). The
  copies diverge on duplicate identity and on special-file handling.
- S2 (medium, auto-fix → R1): the spec-tag rule is derived in
  `internal/preflight/decision.go:565` and `ticket_grammar_test.go:149-161`.
- S3 (medium, auto-fix → R4): the seed-owner triple is stated in
  `registry_data.go:31-51`, `registry_data_test.go:41`, and
  `ticket_grammar_test.go:25`.
- S4 (low, auto-fix → R4): provenance narration in comments at
  `registry_data.go:29-30`, `registry_data_test.go:38-39`,
  `ticket_grammar_test.go:23-24`.
- S5, S6 (low, no-op): `holdsString` repetition and fixture-builder repetition
  are incidental text; an abstraction would be worse.
- S7 (low, auto-fix → R4): the SKILL.md "three rules" numeral is a derived
  count; drop the numeral.

## Spec

Count: 8. Worst: the canary family substituted away a class the mutation
harness can produce. Rows fully met: 36 of 39; TG25, TG29, TG37 partial.

- P1 (medium, ask-user → accepted as R3): `MUTATE.json` can rewrite
  `registry_data.go` to name an absent file, so the
  absent-bound-registry-file class can red a real gate run. The
  uncovered-declared-row half of the substitution stays impossible (no-op).
- P2 (medium, ask-user, reviewer-pending): in a tickets-only folder only the
  foreign-tag check is skipped; `missing Covers`, `duplicate`, and `malformed`
  still bite. A light-path ticket must carry `Covers: none`. The behavior is
  coherent and the template advertises it, but it is narrower than TG25's
  wording. Decide: relax under the empty tag, or amend TG25.
- P3 (medium, auto-fix → R1): sweep diagnostics name the folder, not the
  ticket file, for every non-edge class; preflight names the file.
- P4 (low, ask-user → folded into R1): preflight classifies before the
  extension test. The sweep skips non-`.md` first. R1 adopts the preflight
  order in both venues.
- P5 (medium, ask-user, reviewer-pending): the closures force `Writes:`
  entries the ticket never edits (six fixture directories, five registry
  files). The parked closure-ergonomics idea covers this.
- P6, P7 (low, no-op / reviewer-pending): the example-test fence stripping is
  fragile-later only. The dangling-blocker row-name split is consistent with
  the map and stays a naming question.
- P8 (low, reviewer-pending): the mid-build fence amendments
  (`reviews/ticket-grammar.md`, `internal/systemtest/owner_land_race_test.go`).

## Coverage

Count: 8. Worst: an unguarded control byte in a blocker basename reaches a
TOON detail cell.

- C1 (high, ask-user → accepted as R2): `tickets.go:174` formats blocker
  targets with `%s`. `gather.go:423` and `decision.go:359` carry the byte into
  a `tickets-parse` detail cell. `Writes:` entries get `toon.Representable`.
  Blocker targets get nothing. No test plants a control byte in a field value.
- C2 (high → R1), C3 (medium → R1), C4 (low-medium → R1): the sweep diverges
  from preflight three times. It drops a duplicate basename silently. It skips
  special files named non-`.md`. It reads bytes with no size bound.
- C5 (low → R1): a digit-leading row ID yields the empty tag and silently
  disables the Covers checks; add the degenerate-tag case at the shared
  derivation.
- C6 (low, auto-fix → R2): no table row exercises a lax `Covers:` separator
  (`TG1,TG2` and `TG1,`); pin the fail-closed behavior.
- C7, C8 (low, no-op): `Cycles` over an empty set is nil-safe; a directory at
  a bound file path stays a non-goal.

## Repair targets

- R1 — single-source the ticket enumeration and the spec-tag derivation.
  The sweep gains classify-first order, the duplicate-identity diagnostic,
  the preflight size bound, and per-file diagnostic prefixes. Update the
  affected canary EXPECT lines. Add the degenerate-tag case.
- R2 — guard blocker basenames for representability like `Writes:` entries.
  Plant a control byte in a field value in both suites. Add the lax-separator
  `Covers:` table rows.
- R3 — add the absent-bound-registry-file canary fixture through
  `MUTATE.json` over `registry_data.go`.
- R4 — single-source the seed-owner list; reword the provenance comments
  timelessly; drop the SKILL.md rule numeral.
