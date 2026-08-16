# Pocock engineering skills vs. Bench's spec-and-ticket fences

Read 2026-08-16 from `~/workspace/skills/skills/engineering/` (outside this
repository, not vendored). Read in full: `to-spec/SKILL.md`,
`to-tickets/SKILL.md`, `implement/SKILL.md`, `tdd/SKILL.md`, `tdd/tests.md`,
`tdd/mocking.md`. Not read: the `agents/openai.yaml` file beside each skill, and
every other skill in that tree (`code-review`, `codebase-design`,
`domain-modeling`, `grill-with-docs`, `prototype`, `research`, `triage`, and the
rest).

## Surface comparison

| subject | Pocock | Bench |
|---|---|---|
| spec authoring | `to-spec`, ~90 lines, 8-section template | `bench-write-spec` 190 + `craft-spec` 71 + `craft-seams` 98 + `craft-domain` 70 |
| testing section | prose: what a good test is, which modules, prior art | acceptance coverage map (6 columns, checked by `bench coverage --check`), red-signal grammar, edge inventory, `Won't handle` lines, seam diagram, ownership fences, per-group `Line:` |
| review before code | none; the `to-tickets` quiz is the only gate | 2 rounds x 2 artifacts, delegated cross-family |
| tickets | `to-tickets`: vertical slice, demoable alone, one context window, blocking edges, quiz | `craft-tickets`: same, plus "smallest independently-green", write-delegate charge, `Writes:` |
| implement | 6 lines: TDD at agreed seams, typecheck, full suite once, code-review, commit | 60 lines plus worktree / preflight / land choreography |
| TDD | 38 + 77 + 59 lines | 117 + 50 + 45 lines |

## Where the two agree

Seam confirmation before any test is load-bearing in both: `to-spec` step 2
sketches seams and confirms them with the user; `tdd` refuses a test at an
unconfirmed seam; `craft-seams` and `craft-tdd` say the same. The TDD skills are
near-identical in substance — Bench's additions are the over-fit guard, the
compile-error-red rule, and the honest-stub rule, all of which are genuine.

The divergence is entirely upstream of the implementation loop.

## Evidence from this repository

- FT210's spec review blocked on 10 findings; five were false claims about
  current code written into coverage cells, and one reopened a closed decision
  built on such a premise. Those cells exist only because the map demands a
  red-signal prediction per row before any test runs.
- FT210 ran four delegated rounds (spec x2, tickets x2). Both loops ended by
  folding partials at the two-round cap, not by reaching accept.
- FT210's spec is 365 lines with 37 stories and 20 coverage rows, over a
  263-line compiled decision map, for a feature that sliced to 5 tickets.
- A literal reading of `craft-tickets` produced 17 tickets for that spec where
  `to-tickets` produced 5 with the same dependency spine; the reviewer approved
  the 5.
- 26 of 34 current `specs/` directories are `light-path-*` — one ticket file, no
  spec at all.

## Conclusion carried into the map

The harm is specific, not general: Bench's spec artifact makes predictions about
test behavior that then require their own review loop to verify. `to-spec` makes
no claim that can be false about the codebase. Ownership fences, the `Line:`
declaration, and the gate-as-oracle are Bench capabilities Pocock's
single-session flow does not have, and they stay.
