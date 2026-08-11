# Ticket bundle refusal

Status: staged

Decision source: reviewer-confirmed conversation, 2026-08-10 (the AXI one-ticket
post-mortem: remove the author-asserted keep-together exception inside the
spec-build lifecycle, enforce ticket size mechanically at assign with a
reviewer-owned override, and disclose an atomic landing's ticket cost at spec
sign-off), repaired 2026-08-10 against the eight accepted findings of the Sol
falsification review of this spec's first draft.

## Problem

The ticket discipline says "smallest independently-green unit," but the rule
carries an author-asserted escape hatch: a ticket stays whole if its author
names a red the thinner cut would strand. The AXI build exercised it — one
ticket with 8 acceptance rows and ~80 closure tokens — and the claim was never
re-derived: the breakdown review is advisory and retains no evidence, and
`assign` checks only that the declared closure graph is closed, never its
magnitude. The retained lifecycle journal is the observed failure: operation
`c41def0c` records `assign` leasing that ticket without complaint. Inside the
spec-build lifecycle the exception's premise is void anyway: checkpoints carry
focused evidence and `promote` is the sole full-gate boundary, so a mid-run
ticket can never strand a project-gate red — sequencing replaces bundling. The
reviewer's standing small-ticket preference had no refusal behind it, so a
persuasive paragraph overrode it.

## Solution

Ticket size becomes reviewer-owned and mechanically enforced. `bench spec build
assign` refuses any ticket — whatever its grammar generation — whose acceptance
rows or closure tokens exceed a named bound, unless the ticket's header block
carries an explicit `Bundle-approved:` line: the reviewer's marker, never the
author's argument. `craft-tickets` drops the author-asserted keep-together
exception for lifecycle builds — including its breakdown-review charge item —
routing an apparent stranded red through the existing structural cuts
(prefactor, junction, expand–migrate–contract) with totality oracles landing in
the terminal ticket. `craft-spec` requires a spec that mandates a single-ticket
landing to state that ticket's size, and `/bench-write-spec`'s approval table
carries the disclosure, so "atomic migration" can no longer be approved without
the reviewer seeing what it costs.

This spec's own disclosure: no single-ticket landing is mandated; the implied
breakdown is four tickets — the refusal core (TB1 + TB3: 2 rows, ~10 tokens),
the override and its anchoring (TB2 + TB4: 2 rows, ~8 tokens, blocked by the
core), the `craft-tickets` edit (TB5: 1 row), and the `craft-spec` +
`/bench-write-spec` edit (TB6: 1 row) — each inside both bounds, none needing
`Bundle-approved:`. The round-2 review measured the draft's single code ticket
at four rows and at least 16 tokens; it is split rather than overridden,
which is the behavior the discipline prescribes.

## User stories

1. As the reviewer, I want `assign` to refuse an oversized ticket that lacks my
   explicit `Bundle-approved:` marker — for every size dimension, whatever the
   ticket's grammar generation — naming the exceeded dimension, the measured
   count, the bound, and the override form, so a bundle reaches the lifecycle
   only by my decision. Line: gpt-5.6-terra / medium. The refusal is fully
   gate-observable but sits on the lifecycle authority path.
2. As the reviewer, I want `craft-tickets` to route every keep-together claim
   through a structural cut instead of an author-named red, with no surviving
   reference teaching the old exception, so the next build cannot argue a
   bundle past the default the way the AXI build did. Line: gpt-5.6-terra /
   medium. Prose semantics the gate cannot grade.
3. As the reviewer, I want a spec that mandates a single-ticket landing to
   disclose that ticket's size, and the `/bench-write-spec` approval table to
   surface the disclosure, so the bundling decision is made by me at sign-off
   rather than inherited by ticketing. Line: gpt-5.6-terra / medium. Prose
   semantics the gate cannot grade.

## Implementation decisions

- Two bounds, named constants owned by `internal/specbuild`, one source:
  acceptance rows > 5 and closure tokens > 15. Rows are counted
  deflation-resistantly: `R`-ranges count by their expansion as today, and an
  acceptance ID of any other tag written in range shape (`[A1-A8]`) counts by
  its span for the bound even though it stays one literal row semantically —
  the parser's documented tolerance for that shape must not become a
  compression that reads an eight-row bundle as one.
  Fence-entry count is deliberately **not** a dimension: it measures spelling,
  not scope — five exact `internal/gate/*` entries collapse to one broader
  `internal/gate/` prefix, so a fence bound is author-avoidable in the wrong
  direction and wrongly refuses a legitimate one-row tracer
  (`own-gate-run-binary`: 1 row, 8 tokens, 9 exact entries). Fence breadth
  stays review's question. Calibration on the retained corpus: 8/79 and 7/7
  refuse; 4/4, 1/1, and 1/8 pass. The values are reviewer-approved at this
  spec's sign-off and change only by reviewer edit.
- The refusal is grammar-independent: it applies to every ticket `assign`
  parses, modern or not, because "modern" is inferred from author-controlled
  content and a freshly written legacy-shaped ticket must not dodge the bound
  the way it already dodges `requireClosure`. Acceptance rows parse under
  every grammar generation; a ticket with no `Closure:` line has zero tokens,
  so that dimension is vacuous for it rather than exempting it. A genuinely
  historical oversized ticket, should one surface, takes a one-line
  `Bundle-approved:` addition — reviewer-priced, strands nothing.
- `Bundle-approved:` is recognized **only in the ticket's header block**: the
  region before the first `##` section heading, which in every corpus ticket
  holds the H1 title and the `Blocked by:` sibling fields. The H1 title does
  not terminate the block — the round-2 review showed a "first `#` line" rule
  would make the canonical field position inert in all 65 corpus tickets. An
  occurrence at or after the first `##` line (body prose, a fenced example, a
  quoted grammar token) is inert, so quoting the field cannot lift the
  refusal; this is the hostile-input checklist's quoted-grammar-token class
  applied to an authority-bearing field. Mechanically, only a non-empty value lifts the
  refusal — for every dimension at once; the value's honesty (who approved,
  when) is review's question, exactly like `covers local`. The existing
  skipped-line rule keeps sibling tickets parsable; the existing one-line rule
  means a wrapped continuation reads as its first line only.
- The check runs in `assign` beside `requireClosure`, after parse and before
  any operation begins. The refusal is an ordinary lifecycle error in the
  current `assign` idiom (prose `fmt.Errorf`, exit 1 at the CLI). The
  typed-remedy envelope is the staged AXI family's concern; this spec does not
  couple to an unpromoted candidate, and the AXI matrix's constructor closure
  picks the new refusal up as one more class when that build's derivation
  next runs.
- `craft-tickets` edits: the "Draft the breakdown" escape hatch (name-a-red to
  keep a group whole) is removed for spec-build lifecycle work; its
  signal-table row routes an apparent stranded red to the structural cuts the
  skill already owns, with `Blocked by:` sequencing and totality oracles in
  the terminal ticket, because promote is the only full-gate boundary. The
  breakdown-review charge item that re-derives the keep-together claim by
  attempting the split is replaced in the same edit: the pass instead verifies
  that every grouping routes through a structural cut and that any
  `Bundle-approved:` line traces to a recorded reviewer sign-off — no orphan
  reference to the removed exception survives anywhere in the skill. The
  commit-on-green paths (light path, `bench shift`) keep a narrow exception
  that is reviewer-owned: explicit approval on the session's approval surface,
  never author prose. The ticket-shape section documents the header-block
  `Bundle-approved:` line and the fact that assign refuses an oversized
  ticket without it; the enforcement semantics stay owned by the code.
- `craft-spec` edit: in `Story sizing and scope cuts`, a Solution or
  implementation decision that mandates a single-ticket landing states that
  ticket's size (rows and closure tokens); an above-bounds bundle is a
  reviewer sign-off decision recorded on the resulting ticket as its
  `Bundle-approved:` line. `/bench-write-spec`'s approval-table enumeration
  adds the disclosure by pointer to `craft-spec`'s rule, because the command
  owns the table's contents and a skill-only edit would leave the phase
  emitting tables without it.
- Skills are consumed via symlinks (`.claude/skills` → `.agents/skills`), so
  the prose fences are under `.agents/` alone.
- The stories partition into a code fence and a prose fence, and the bundle is
  chosen, not defaulted: the `Bundle-approved:` field name and the existence
  of the refusal are one fact advertised in prose and enforced in code, so
  shipping either half alone would advertise an enforcement that does not
  exist or enforce one nothing documents. Chosen bundle ≠ single-ticket
  landing: the three tickets above land independently green.

## Testing decisions

- TDD attaches at the public `ParseTicket` + assign-refusal seam in
  `internal/specbuild`, the same seam `requireClosure`'s tests already use:
  write a ticket fixture file, call the public operation, assert the exact
  refusal or acceptance.
- The gate observes story 1 through `go test ./internal/specbuild`. Stories 2
  and 3 are guidance prose the gate cannot grade; their deltas are checkable
  by grep and their semantics are review-graded — classified honestly below.
- Boundary discipline: each bound is tested at the bound (accepted) and one
  past it (refused), per dimension; the override is tested per dimension and
  with all dimensions exceeded at once.

### Seam diagram

    trigger: bench spec build assign <slug> --ticket <file>
        │
        ▼
    ticket file ──▶ [ ParseTicket → bundle-bound check → requireClosure ] ──▶ lease, or named refusal
                          ◀ tests attach here: fixture ticket files driven
                            through the public assign path, asserting the
                            exact refusal text and the acceptance cases

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| TB1 | 1 | Assign refuses a ticket of any grammar generation exceeding either bound — rows > 5 post-expansion, closure tokens > 15 — when the header block carries no non-empty `Bundle-approved:`; the message names the exceeded dimension, measured count, bound, and the override line form, one refusal test per dimension. | public assign path in `internal/specbuild` | observed red: the retained lifecycle journal (operation `c41def0c`, run digest state) records `assign` leasing the 8-row/79-token AXI ticket; `rg -n "Bundle-approved" internal .agents` additionally exited 1 | The retained journal is a real observation of the exact wrong behavior — assign leased the bundle this row refuses — not an absence claim. |
| TB2 | 1 | A non-empty header-block `Bundle-approved:` lifts the refusal for **each** dimension independently — a rows-only overflow, a tokens-only overflow, and a both-dimensions overflow are each accepted with the marker — and a ticket at both bounds exactly (5 rows, 15 tokens) is accepted without one. | public assign path | not TDD-able until TB1 supplies the check; today every case passes vacuously | Enumerating the override per dimension rejects an implementation that honors the marker for row overflow while still refusing approved token overflow; the boundary cases keep the refusal from shrinking legitimate tickets. |
| TB3 | 1 | A legacy-shaped ticket (no `Closure:`, no `## Red mutations`) with more than 5 acceptance rows is refused like any other; one with 5 or fewer rows keeps its existing assignability, `requireClosure` exemption included. | public assign path | not TDD-able until TB1 supplies the check; today the existing legacy-ticket test proves such content reaches assign ungated | "Modern" is author-controlled content, so a modern-only check is an authoring-time bypass; grammar-independence closes it while the vacuous-token rule strands no genuinely small legacy ticket. |
| TB4 | 1 | A `Bundle-approved:` line outside the header block — in body prose or inside a fenced code example — does not lift the refusal. | public assign path | not TDD-able until TB1 supplies the check; the current parser scans every line globally (`ParseTicket`'s single trimmed-line loop), so position-anchoring has no owner yet | An authority-bearing field recognized anywhere in the file is forgeable by quotation — the checklist's quoted-grammar-token class; anchoring makes the forgery inert and this row proves it. |
| TB5 | 2 | `craft-tickets` carries no reference to the author-asserted exception anywhere in the skill: the name-a-red sentence and its signal-table row are replaced by the structural routing, the breakdown-review charge item no longer instructs reproducing a keep-together red and instead verifies structural routing plus `Bundle-approved:` provenance, the commit-on-green exception is marked reviewer-owned, and the ticket-shape section documents the header-block field and the assign refusal. | `.agents/skills/bench-craft-tickets/SKILL.md` | reviewer-graded prose; the mechanical delta is checkable — `rg -n -i -e "specific red" -e "thinner cut" -e "keep.{0,30}together" .agents/skills/bench-craft-tickets/SKILL.md` observed matching lines 53, 59, and 392–393 today and must return no exception-teaching match after | The falsification review proved the sentence-only probe green while line 392's breakdown-review charge still taught the exception; the whole-skill sweep is the delta that catches the orphan. |
| TB6 | 3 | `craft-spec`'s `Story sizing and scope cuts` requires a single-ticket-landing Solution to state that ticket's rows and closure tokens, with an above-bounds bundle recorded as the ticket's `Bundle-approved:` line; `/bench-write-spec`'s approval-table enumeration names the disclosure by pointer to that rule. | `.agents/skills/bench-craft-spec/SKILL.md` and `.agents/commands/bench-write-spec.md` | reviewer-graded prose; delta checkable — `rg -n "Bundle-approved" .agents/skills/bench-craft-spec/SKILL.md .agents/commands/bench-write-spec.md` observed exiting 1 today | The approval table's contents are enumerated by the command, not the skill, so a skill-only edit leaves the phase emitting tables without the disclosure — the cheapest wrong implementation the review named. |

### Edge inventory

- Error path — TB1 is the refusal itself; a refused assign leaves no operation
  journal entry, matching every other pre-operation refusal in assign.
- Empty or absent input — `Bundle-approved:` absent, present-but-empty, and
  whitespace-only all refuse (TB1); an empty `Closure:` still refuses modern
  tickets under `requireClosure`, unchanged.
- Boundary values — TB2 pins at-bound acceptance for both dimensions jointly;
  TB1 pins one past each bound separately.
- Malformed input — a wrapped `Bundle-approved:` continuation reads as its
  first line only (existing one-line grammar rule); a value on the next line
  is dropped and the field reads empty, refusing.
- Interrupted or partial state — the check is pure and pre-operation; nothing
  new is journaled, so interruption semantics are unchanged.
- Re-run idempotency — a refused assign is repeatable with identical output;
  an accepted at-bound assign follows the existing request-replay contract.
- Process-boundary lifecycle — the ticket is re-read and re-parsed on every
  assign; no new durable state exists to reload.
- Hostile environment — the quoted-grammar-token class is TB4 (a fenced or
  quoted `Bundle-approved:` must be inert); counted fields are parsed lists;
  the override value passes through no renderer on the refusal path and is
  never echoed into a command; existing sanitize rules on error text apply.
- Command self-observation — the check reads the ticket file only; it mutates
  nothing it reports.
- Special files and dangling symlinks — `ParseTicket` already Lstats and
  refuses non-regular ticket files; no new discovered-path reader is added.

**Won't handle:** the honesty of a `Bundle-approved:` value (who approved,
whether they meant it) — review's question, like `covers local`; TB5's edited
breakdown-review charge is where its provenance is verified. Batching *below*
the bounds (the four-finding repair shape) — a size proxy cannot see
independence; the checkpoint-review capability owns making that visible.
Re-checking size at checkpoint or integrate — the assign-time ticket digest
already pins the graded content. A fence-entry bound — measures spelling, not
scope, and rewards broader prefixes; fence breadth stays with review.

## Ownership fences

- `internal/specbuild/`
- `.agents/skills/bench-craft-tickets/SKILL.md`
- `.agents/skills/bench-craft-spec/SKILL.md`
- `.agents/commands/bench-write-spec.md`

## Out of scope

- Checkpoint-scoped advisory three-axis review — its own spec, next in this
  batch (~10 edits, 1 promotion gate).
- Typed refusal envelope and `help[]` action for the new refusal class — the
  staged AXI family build owns refusal rendering (0 edits here; the class
  joins its derived matrix when that build lands).
- Repair-entry cost around review findings — FT200 (roadmap-owned).
- Position-anchoring the *existing* header fields (`Blocked by:`,
  `Contracts:`, …) — same weakness, but none is authority-bearing; a separate
  hardening pass if wanted (~4 edits, 1 gate run).
