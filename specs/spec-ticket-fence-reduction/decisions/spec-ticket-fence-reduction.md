# Spec and ticket fence reduction

Status: ready

## Destination

Cut the spec-and-ticket authoring fences down to what earns its keep, remaking
`craft-spec` on Pocock's `to-spec` and `craft-tickets` on `to-tickets`, in one
bundled spec. Concretely: the acceptance coverage map loses its `red signal`
column, the four-round pre-implementation review loop becomes one same-family
round over the spec-and-tickets pair, the edge-class walk moves to the TDD loop,
the spec template moves into `craft-spec` and `bench-write-spec.md` shrinks under
a new 60-line prose budget, and `craft-tickets` takes `to-tickets`' sizing rule.
`internal/coverage` lands the five-column header expand-contract. The light path,
the four invariants, the gate, and `/bench-implement-spec`'s worktree machinery
are untouched.

## #1: Does this map decide the spec artifact alone or the spec-and-ticket fence set?

Blocked by: none
Type: Grill

### Question

`craft-tickets`' remake is already an open learning (17-vs-5 on FT210). One map
or two?

### Answer

One map covering `craft-spec`, `craft-tickets`, `bench-write-spec`, and their
shared enforcement. Both edit the same anchor rows in
`internal/anchors/registry_data.go` and the same README/CHANGELOG vocabulary, so
splitting means editing those rows twice. `/bench-implement-spec` is out of scope.

## #2: Does the machine-checked acceptance coverage map survive?

Blocked by: none
Type: Grill

### Question

Keep the table and `bench coverage --check`, or replace it with `to-spec`'s
Testing Decisions prose?

### Answer

It survives, reduced. The orphan-story check — every declared story is referenced
by a row or carries a `Not covered: story <n> — <reason>` line — is the breadth
floor's only enforcement and prose cannot do it. `bench coverage --check` keeps
that check, the one-predicate-per-behavior check, and the four-story fan-out
bound.

## #3: What replaces the two-rounds-per-artifact pre-implementation review loop?

Blocked by: none
Type: Grill

### Question

FT210 ran four delegated rounds (spec ×2, tickets ×2) and terminated on the
two-round cap with partials folded, not on convergence.

### Answer

One round over the spec-and-tickets pair, with `craft-tickets`' granularity /
blocking-edges / merge-split quiz as the approval round. No second
fix-verification round.

## #4: Where does the edge inventory live?

Blocked by: none
Type: Grill

### Question

Walk the edge classes at spec time, or at the seam in the TDD loop?

### Answer

The walk moves to the TDD loop, which `craft-tdd` already describes at the seam
where the classes are visible. The spec keeps the project profile's
hostile-input checklist attachment. The requirement that every walked edge class
produce either a coverage row or a `**Won't handle**` line is retired.

## #5: Which columns survive in the reduced coverage map?

Blocked by: #2
Type: Grill

### Question

Which of the six columns are cut?

### Answer

The header becomes exactly
`| row | story | behavior | seam | why it catches the failure |`. Only the
`red signal` cell is cut: it is the one column that asserts something about
current code before any test runs, and it produced five of FT210's ten blocking
findings. `why it catches the failure` stays as the anti-degenerate check; the
`row` ID column stays, keeping ticket-covers traceability.

## #6: How does the schema change land against specs in flight?

Blocked by: #5
Type: Grill

### Question

`specs/worktree-landed-retirement/spec.md` is staged with 20 six-column rows.

### Answer

Expand-contract. The parser accepts both the five-column and six-column headers;
the six-column form drains as specs retire; a contract ticket blocked by that
drain deletes the six-column branch. No staged spec is rewritten mid-build.

## #7: What sizing rule replaces "smallest independently-green" in craft-tickets?

Blocked by: #1
Type: Grill

### Question

The literal reading produced 17 tickets where `to-tickets` produced 5 with the
same dependency spine; the reviewer approved the 5.

### Answer

`to-tickets`' rule verbatim — a complete vertical slice, demoable alone, sized to
one fresh context window, prefactoring first — plus the learning's clause: a
coverage row that only adds a test to a seam its parent slice already opened is
that slice's acceptance row, not a ticket. The splitting sentence and its anchors
are retired in the same change. The phrase "independently-green implementation
tickets" survives as README/CHANGELOG/CONTEXT vocabulary, because it is what
distinguishes implementation tickets from shaping decision tickets.

## #8: Where does the spec template live?

Blocked by: #1
Type: Grill

### Question

`bench-write-spec.md` is 190 lines and carries the template; `craft-spec` is 71
and owns the discipline.

### Answer

The template moves into `craft-spec`. `bench-write-spec.md` keeps only the entry
contract, ownership, and exit handoff. The anchor rows pinning coverage
vocabulary and the template's header line to the command retarget to the skill.

## #9: One bundled spec, or two?

Blocked by: #1, #2
Type: Grill

### Question

The guidance remake and the `internal/coverage` schema change are independently
useful behaviors.

### Answer

One bundled spec. The anchor rows and the `coverage-map-validation` canaries pin
the column list to the prose, so splitting forces either a red gate between the
two specs or a temporary anchor exemption.

## #10: Does the red-signal grammar survive anywhere in the spec artifact?

Blocked by: #5
Type: Grill

### Question

`craft-spec` defines four forms (`observed red:`, `not observed:`,
`already covered:`, `not TDD-able:`), unenforced by the gate.

### Answer

No. The four-form grammar is deleted from `craft-spec`. `craft-tdd`'s
"Acceptance rows" section is the single source: it classifies a row as
`already covered` or `not TDD-able` at the moment the row runs, where the
classification is observed rather than predicted.

## #11: Does **Won't handle** survive?

Blocked by: #4
Type: Grill

### Question

With the edge walk in the TDD loop, does the spec keep an edge-level exclusion
construct?

### Answer

Yes, for reviewed exclusions only: an edge the author deliberately excludes gets
a `**Won't handle**` line as the reviewer's veto surface. The requirement that
every walked class produce a row or a line is what is retired (see #4).

## #12: Does the surviving review round stay cross-family?

Blocked by: #3
Type: Grill

### Question

FT210's rounds ran `gpt-5.6-sol`/high through `codex exec`.

### Answer

No. The one round runs as a same-family delegate through the invoking harness's
native agent surface. `bench-write-spec` drops its cross-family branch and the
unbound-model-id refusal rules.

## #13: Does references/cross-harness-reviewers.md survive?

Blocked by: #12
Type: Grill

### Question

The file is a general `craft-delegate` capability, referenced from
`craft-delegate`, `bench-write-spec`, `craft-skills`, and two canaries.

### Answer

It survives, unowned by write-spec. `craft-delegate` keeps the recipes for any
cross-family delegation; only `bench-write-spec`'s cross-family clause and its
resolution rules go. The `cross-harness-reviewer-recipes` and
`delegate-cross-harness-reviewer-pointer` canaries stay green.

## #14: Does --reviewer survive on /bench-write-spec?

Blocked by: #12
Type: Grill

### Question

The flag currently resolves a tier or a bound model id, own-family or
cross-family.

### Answer

It survives as a same-family tier-or-effort override: `--reviewer <tier> [effort]`
resolves through `.bench/lines.env` for the invoking harness's own column. One
resolution path; the cross-family branch is removed.

## #15: What tier does the single review round get by default?

Blocked by: #3
Type: Grill

### Question

Four rounds become one; where does the line sit?

### Answer

Mid tier at high effort. The round now hunts a smaller artifact carrying no
current-code claims, so the falsification job is narrower than the loop it
replaces. Top stays a reviewer-approved escalation per `projects/benchkit.md`.

Superseding note, 2026-08-16: for this map's own spec the reviewer named `fable`
(`BENCH_CLAUDE_TOP`) at high effort — an approved top-tier escalation consistent
with the profile's doc-authoring leverage override, since the artifact under
review is kit guidance prose. This fixes the tier for this build only; the
default the spec writes into `bench-write-spec` remains mid.

## #16: Does the light-path threshold move?

Blocked by: #1
Type: Grill

### Question

26 of 34 current `specs/` directories are `light-path-*`.

### Answer

No. The right-size table in `.bench/BENCH.md` is unchanged. Moving the threshold
in the same change would confound the measurement of whether the cheaper heavy
path is actually used.

## #17: Does the bench-write-spec.md shrink get a prose-budget row?

Blocked by: #8
Type: Grill

### Question

`projects/benchkit.md` budgets `.bench/BENCH.md` (180),
`bench-implement-spec.md` (60), `craft-tickets` (100), and a 120-line skill glob.
`bench-write-spec.md` is unbudgeted at 190 lines.

### Answer

Yes — add a `.agents/commands/bench-write-spec.md | 60` row, matching
`bench-implement-spec.md`. The `guidance-prose-budgets` check then makes the
shrink permanent rather than a one-time diff that re-accretes.

## #18: How do the two open capture/learnings.md entries close?

Blocked by: #1
Type: Grill

### Question

Both 2026-08-16 entries (FT210's false coverage cells; craft-tickets 17-vs-5) are
answered by decisions in this map.

### Answer

This map is their verdict. The next `/bench-what-next` drain marks both resolved
against it and adds the roadmap row, rather than re-deliberating them.

## Not yet specified

## Spec-writer discretion

- Which exact anchor rows in `internal/anchors/registry_data.go` retarget versus
  retire, and which canary fixtures under `tests/canary/coverage-map-validation`
  and `tests/canary/workflow-guidance-anchors` update — `craft-spec`'s
  read-the-enforcement-surface rule owns this inventory at authoring time.
- The mechanical follow-ons of the column cut: `craft-delegate`'s write-delegation
  charge line ("behavior, seam, red signal"), `craft-review`'s Coverage axis
  wording, and `bench coverage <spec>` task seeding in `/bench-implement-spec`.
- `CONTEXT.md` glossary entries for **coverage row** and **acceptance row** land
  in the same green change as the behavior, not before — a glossary written ahead
  of the change would describe a tree that does not exist yet.
- Ticket slicing and `Blocked by:` edges for the bundled spec.

## Out of scope

- `/bench-implement-spec` and its worktree, preflight, and `bench worktree land`
  choreography (#1). The worktree machinery is a real capability, not ceremony.
- The right-size table's light-path threshold (#16).
- The four invariants, `bench gate`, and `/bench-review-implementation`'s
  three-axis structure.
- Rewriting `specs/worktree-landed-retirement/spec.md`'s existing rows (#6).
- Deleting `references/cross-harness-reviewers.md` (#13).

## Sources

- Path: `specs/spec-ticket-fence-reduction/decisions/pocock-skills-comparison.md`
  Supports: #1, #3, #4, #5, #7 — the surface comparison of Pocock's `to-spec`, `to-tickets`, `implement`, and `tdd` against Bench's fences, read 2026-08-16 from outside this repository, with what was and was not read named in the asset.
  Drift: re-read the upstream skills if they change; re-verify the FT210 counts if `specs/worktree-landed-retirement/` is retired or re-sliced, and the light-path ratio if `specs/` turnover is material.
- Path: `internal/coverage/coverage.go`
  Supports: #2, #5, #6, #10 — `Check` validates the five- or six-column header, non-empty cells, one predicate per behavior, the `bounds.CoverageRowStories` fan-out bound, and orphan stories; the four-form red-signal grammar is prose-only and unenforced.
  Drift: re-read before spec authoring if the coverage parser, `internal/bounds`, or the header constants change.
- Path: `internal/anchors/registry_data.go`
  Supports: #7, #8, #9, #11 — six rows pin `acceptance coverage map`, `red signal`, `edge inventory`, `Won't handle`, the six-column template header, and the falsification-question sentence to `.agents/commands/bench-write-spec.md`; further rows pin `independently-green` vocabulary to `.bench/BENCH.md`, `README.md`, and `CHANGELOG.md`.
  Drift: re-enumerate at authoring time — `craft-spec`'s read-the-enforcement-surface rule owns the exact retarget-versus-retire inventory, and the registry moves with every landed spec.
- Path: `capture/learnings.md`
  Supports: #5, #7, #18 — the two open 2026-08-16 entries (FT210's false coverage cells; `craft-tickets` 17-vs-5).
  Drift: re-read after any `/bench-what-next` drain, which may close or reword both entries.
- Path: `projects/benchkit.md`
  Supports: #15, #17 — the guidance-prose-budget table is the one source the `guidance-prose-budgets` check parses, and `.agents/commands/bench-write-spec.md` has no row; top tier remains a reviewer-approved escalation.
  Drift: re-read if the budget table, the tier bindings, or the line-routing section changes.
