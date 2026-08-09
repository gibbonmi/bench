# Spec-ticket handoff contract

Status: staged

Decision source: reviewer-confirmed current conversation on 2026-08-09

## Problem

A reviewer can approve a Bench spec whose coverage, seams, and writer limits are
clear in the spec but become optional or implicit when implementation tickets are
derived. The default spec template omits coverage-row IDs and ownership fences,
so the ordinary authored artifact does not activate the lifecycle's existing
row-totality check and gives ticketing no canonical writer envelope. Ticket
authors must also invent how an observed red, an already-covered row, or a
not-TDD-able row becomes ticket acceptance and mutation evidence. The gap is
discovered during slicing or implementation instead of being visible at spec
approval.

## Solution

Make the spec-to-ticket handoff an explicit reviewed contract. New specs default
to identified coverage rows, record exact spec-time ownership fences, and show
both in the approval surface. Ticket breakdown derives one accountable ticket
disposition for every spec row, distinguishes spec-time red evidence from
post-implementation mutation evidence, accounts for every story, seam, edge
disposition, and fence before assignment, and routes any required fence widening
back to spec approval. Existing five-column specs remain valid; the change fixes
what the authoring phase produces rather than retroactively invalidating older
artifacts.

## User stories

1. As a reviewer, I can approve coverage rows with stable IDs and exact ownership
   fences in the spec itself, so the build cannot silently drop either when it
   creates tickets.

   Line: `gpt-5.6-sol / high`. This changes high-leverage kit guidance whose
   semantic failures are only partly gate-observable.

2. As a ticket author, I can deterministically disposition every spec row into
   ticket acceptance and mutation evidence without treating `covers local` as an
   escape from reviewed scope.

   Line: `gpt-5.6-sol / high`. The handoff rule spans semantic prose and existing
   lifecycle enforcement, so a confident but wrong formulation can remain green.

3. As a build reviewer, I see a complete handoff ledger before assignment that
   accounts for every story, row, seam, edge disposition, and ownership fence and
   stops on fence drift, so omissions surface before provisional work begins.

   Line: `gpt-5.6-sol / high`. This is reviewer-facing workflow guidance with a
   cross-phase completion criterion.

## Implementation decisions

- Origin: **learnings**. The external `to-spec` comparison prompted the audit,
  but the candidate is justified by observed gaps in Bench's own spec and ticket
  contracts; no upstream tracker or template behavior is imported.
- `craft-spec` remains the single owner of the acceptance-map schema and
  spec-time ownership fences. Its default authored shape uses the leading `row`
  column. A fence entry is an exact file or path prefix, never a glob or an
  implementation ticket. The coverage parser continues accepting the legacy
  five-column shape.
- `/bench-write-spec` renders the six-column map in its template, adds an
  `## Ownership fences` template section, and includes both row IDs and fences in
  the approval table. It points to `craft-spec` for their meaning rather than
  restating their grammar.
- `craft-tickets` owns the handoff derivation. The handoff ledger names one
  accountable ticket for each spec coverage row: the first claimant in the
  approved build-plan order, which already places blockers before consumers.
  Later tickets may name the same row as defense in
  depth, matching existing assign/promote semantics, but they do not change the
  ledger's accountable owner. Additional acceptance discovered from the current
  tree uses `(covers local)` only when it is genuinely absent from the approved
  map or is a repair row.
- The spec red signal and ticket red mutation remain different evidence. An
  observed-red row carries its public failing operation into ticket acceptance;
  the ticket adds an independent post-implementation subject mutation. An
  `already covered` row keeps the existing control as its positive oracle and
  adds a subject mutation that proves the changed route reaches that control. A
  `not TDD-able` row names its blocker, becomes ticket acceptance on the first
  ticket where the seam exists, and receives its mutation there. No ticket copies
  a spec absence probe as a mutation after the absence has ceased to exist.
- Before assignment, the existing ticket-breakdown review emits a handoff ledger
  derived from the spec and ticket files. It accounts for every story, coverage
  row, named seam, edge row or `Won't handle` disposition, and spec ownership
  fence. Every mapped row has an accountable ticket, every ticket path stays
  inside a spec fence, every used seam appears in ticket acceptance, and every
  fence has a ticket owner or an explicit unused disposition.
- Contract discovery still re-derives value traffic and integration surfaces
  from the current tree after predecessors land. That refresh may narrow a ticket
  inside an approved fence. A required path outside all approved spec fences is
  spec drift: ticketing stops and routes the exact path and reason back through
  `/bench-write-spec`; it never widens the fence itself.
- `/bench-implement-spec` points fence drift to the same spec-repair route as seam
  drift and requires the approved handoff ledger before `bench spec build start`.
  It does not duplicate the ledger schema owned by `craft-tickets`.
- The workflow-guidance anchor registry and its mutation fixtures prove the new
  default map, fence section, derivation branches, handoff-ledger totality, and
  fence-drift route remain present and bite when removed. The implementation adds
  no new lifecycle command, status, or promotion authority.
- The change is learning-sourced kit guidance and receives one concise
  `CHANGELOG.md` entry.
- The three stories remain one deliberate spec because they share the same
  authored map/fence contract, the same ticket-breakdown handoff, and the same
  workflow-guidance oracle. Splitting them would make either the producer or the
  consumer green without an end-to-end handoff.
- `Bootstrap authority before execution` is not applicable. No story authorizes
  or authenticates an executable chain, the fence-drift stop is reviewer-owned
  guidance, and the change introduces no lifecycle authority.

## Handoff contract demonstrated by this spec

This section is reviewer explanation, not a second normative schema. It shows the
new information in this spec that the current `/bench-write-spec` template does
not require, and names the proposed owner of each rule.

| Current template omission | Present in this spec | Proposed canonical owner |
|---|---|---|
| The map starts with `story`, so ticket `covers` traceability is optional. | Every row below has a stable `SH` ID. | `craft-spec` owns the row schema; `/bench-write-spec` renders its default. |
| There is no ownership-fence section or fence approval item. | `## Ownership fences` below gives exact writer envelopes and the approval table dispositions them. | `craft-spec` owns fence meaning; `/bench-write-spec` owns the artifact slot. |
| No rule transforms the three red-signal classifications into ticket evidence. | Implementation decisions define observed-red, already-covered, and not-TDD-able routes separately. | `craft-tickets` owns the handoff derivation. |
| Breakdown review checks ticket-local honesty but not complete spec handoff. | Story, row, seam, edge, and fence totality is an acceptance behavior. | `craft-tickets` owns the handoff ledger; `/bench-implement-spec` requires it before start. |
| Ticket contract discovery can imply a wider path without a named authority rule. | Any path outside the approved fences is an explicit return to spec approval. | `craft-tickets` owns discovery; `/bench-implement-spec` owns the phase route. |

## Testing decisions

- A good test mutates one new required clause in a workflow-guidance fixture and
  observes the exact conformance diagnostic through the registered owner. A
  substring-presence assertion alone is insufficient when relocation would make
  the rule inert; section-sensitive clauses use section-scoped anchors or fixture
  mutations that preserve the words while breaking their placement.
- Existing coverage parser and spec-build covers tests are positive controls:
  six-column maps activate assign/promote coverage enforcement, while legacy
  five-column maps continue validating. The implementation does not duplicate
  those parsers or counts in guidance tests.
- New or extended tests attach at `internal/anchors` and
  `internal/conformance`, following `fixture_bite_test.go`,
  `docs_workflow_helpers_test.go`, and
  `tests/canary/workflow-guidance-anchors` as prior art.
- The gate seam is the branch-native `test` phase over the anchor/conformance
  owners, with the retained workflow-guidance canary proving the omission
  mutations bite.
- Fresh-session dogfood is a `craft-synthesis` verification exception rather
  than a coverage row. It grades whether an already-loaded guidance change is
  followed by a new agent, not behavior an implementation ticket can prove with
  an independent subject mutation. After the prose lands, a fresh session authors
  a small spec with IDs and fences, derives tickets and the handoff ledger, and
  stops before implementation.

### Seam diagram

    trigger: bench gate validates the changed kit guidance
        │
        ▼
    command/skill Markdown ──▶ [ workflow-guidance anchor owner ] ──▶ exact diagnostic or green
                                      ◀ tests attach: remove or relocate one required handoff clause

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| SH1 | 1 | The authored spec template leads its acceptance map with unique row IDs. | `/bench-write-spec` template in its `Template` section | observed red: `rg -n '^\| row \| story \| behavior \| seam \| red signal \| why it catches the failure \|$' .agents/commands/bench-write-spec.md` exited 1 | A recommendation outside the canonical template leaves ordinary new specs outside assign/promote totality. |
| SH2 | 1 | The existing coverage parser continues accepting both valid six-column maps and legacy five-column maps after the authored default changes. | exported coverage parser | already covered by `go test -run '^TestParseSpecOptIn$' ./internal/coverage`, whose legacy assertion drives the five-column path after the six-column path | The same exported entry point exercises the new default shape and the retained compatibility path, so deleting legacy support makes the control red. |
| SH3 | 1 | The spec template contains an `## Ownership fences` section whose entries point to `craft-spec` for their meaning. | `/bench-write-spec` `Template` section | observed red: `rg -n '^## Ownership fences$' .agents/commands/bench-write-spec.md` exited 1 | A fence rule left only in `craft-spec` can still disappear from the authored artifact. |
| SH4 | 1 | The reviewer approval table requires an explicit disposition for ownership fences alongside stories, seams, coverage, edges, and scope. | `/bench-write-spec` approval paragraph | observed red: `rg -n -U '(?s)Before a build starts.{0,400}ownership fences' .agents/commands/bench-write-spec.md` exited 1 | Adding the template section without adding it to reviewer veto surface leaves an unapproved writer envelope. |
| SH5 | 2 | An observed-red spec row carries its failing public operation into the accountable ticket, which adds a distinct post-implementation subject mutation. | `craft-tickets` row-derivation section | observed red: `rg -n -F 'An observed-red spec row carries its failing public operation into ticket acceptance' .agents/skills/bench-craft-tickets/SKILL.md` exited 1 | The exact branch rejects copying an obsolete absence probe as the finished mutation. |
| SH6 | 2 | An already-covered spec row keeps its named control and adds a subject mutation proving the changed production route reaches that control. | `craft-tickets` row-derivation section | observed red: `rg -n -F 'An already-covered spec row keeps its named control and adds a subject mutation' .agents/skills/bench-craft-tickets/SKILL.md` exited 1 | The exact branch prevents a pre-existing control from being cited without exercising the new route. |
| SH7 | 2 | A not-TDD-able spec row names its blocker, maps to the first ticket where the seam exists, and receives its subject mutation there. | `craft-tickets` row-derivation section | observed red: `rg -n -F 'A not-TDD-able spec row names its blocker and maps to the first ticket where the seam exists' .agents/skills/bench-craft-tickets/SKILL.md` exited 1 | The exact branch prevents permanent exemption from red evidence after the seam becomes available. |
| SH8 | 3 | The handoff ledger accounts for every user story, coverage row, named seam, edge row or `Won't handle`, and spec ownership fence before assignment. | `craft-tickets` `Review the breakdown before assignment` section | observed red: `rg -n -F 'The handoff ledger accounts for every user story, coverage row, named seam, edge row or' .agents/skills/bench-craft-tickets/SKILL.md` exited 1 | One conjunctive completion sentence cannot be satisfied by mentioning only a row or a fence. |
| SH9 | 3 | A ticket path outside every approved spec fence stops ticketing and returns the exact path and reason to `/bench-write-spec`. | `craft-tickets` contract-discovery section | observed red: `rg -n -F 'A required path outside every approved spec fence stops ticketing' .agents/skills/bench-craft-tickets/SKILL.md` exited 1 | Current-tree discovery remains authoritative without granting itself wider writer authority. |
| SH10 | 3 | `/bench-implement-spec` treats fence drift like seam drift and requires a repaired spec approval plus a complete handoff ledger before lifecycle start. | `/bench-implement-spec` pre-build and ticket-breakdown route | observed red: `rg -n -F 'Fence drift takes the same route and the lifecycle does not start without a complete handoff ledger' .agents/commands/bench-implement-spec.md` exited 1 | A ticket-only stop can be bypassed when the implementation coordinator starts the lifecycle from incomplete files. |
| SH11 | 1–3 | Each of SH1 and SH3–SH10 has its own registered section-sensitive omission or relocation mutation and clause-specific diagnostic. | anchor registry and workflow-guidance mutation fixtures | not TDD-able until the normative clauses exist; current `internal/anchors/registry_data.go` contains no owner for any of their exact predicates | Per-clause mutations prove each requirement remains in the section where a fresh agent acts on it. |

### Degenerate implementations the map rejects

- Story 1's cheapest wrong implementation adds the map and fence headings but
  omits legacy compatibility and reviewer fence approval. SH2 and SH4 remain red.
- Story 2's cheapest wrong implementation explains only observed reds and treats
  the other classifications as ticket-author discretion. SH6 and SH7 remain red.
- Story 3's composition degenerate registers anchors whose words match in harmless
  sections, so the oracle fence is green while the guidance fence remains inert.
  SH11 requires section-sensitive relocation mutations over the real clauses.

### Edge inventory

- Error path — malformed or duplicate row IDs retain the existing coverage
  diagnostics; a missing handoff disposition blocks breakdown approval rather
  than being silently classified local.
- Empty or absent input — a spec with no coverage map remains outside this normal
  non-trivial path and follows the existing explicit fallback; a spec with an
  empty ownership-fence section is incomplete, not an unrestricted build.
- Boundary values — one story/row/fence and many stories/rows/fences both receive
  totality accounting; one ticket may own several rows, while later tickets may
  reinforce the accountable claimant.
- Malformed input — malformed row IDs and unresolved ticket names retain existing
  parser refusals. Breakdown review grades glob-shaped fences and paths outside
  the repo as invalid writer authority; the lifecycle parser does not enforce
  fence shape.
- Interrupted or partial state — a session ending after ticket files are written
  re-derives the ledger from the spec and files before assignment; conversation
  memory is not evidence.
- Re-run idempotency — re-running breakdown review over unchanged spec and ticket
  files produces the same dispositions and no duplicate artifact.
- Process-boundary lifecycle — assign and promote re-read the staged spec and
  ticket files; the handoff does not depend on in-memory author state.
- Hostile environment — paths with spaces or glob characters remain literal
  backticked fence entries; glob syntax is not interpreted as authority.
- Command self-observation — the read-only breakdown review does not edit the spec
  it grades; repairs occur before lifecycle start and are reviewed again.
- Special files and dangling symlinks — existing spec/ticket discovery refusal
  remains authoritative; this guidance change adds no new file reader.

### Won't handle

- Turning semantic story, seam, or fence totality into a new lifecycle parser —
  the existing covers enforcement remains mechanical, while the handoff ledger is
  reviewer-visible breakdown evidence.
- Replacing contract re-derivation with frozen spec claims — current-tree refresh
  remains necessary after predecessor tickets land.
- Enforcing one claimant per spec row — later tickets may reinforce a mapped
  behavior; the handoff ledger identifies the accountable ticket while existing
  promotion requires at least one integrated claimant.

## Ownership fences

- Handoff-contract owner: `.agents/commands/bench-write-spec.md`,
  `.agents/commands/bench-implement-spec.md`,
  `.agents/skills/bench-craft-spec/SKILL.md`,
  `.agents/skills/bench-craft-tickets/SKILL.md`, `projects/benchkit.md`,
  `CHANGELOG.md`, `internal/anchors/registry_data.go`,
  `internal/conformance/fixture_bite_test.go`,
  `internal/conformance/docs_workflow_helpers_test.go`,
  `tests/canary/workflow-guidance-anchors`.
- Ticket-artifact owner: `specs/spec-ticket-handoff-contract/tickets`.

One handoff-contract writer owns each normative clause with its independent
omission oracle so every tracer ticket can land green. A required path outside
these fences returns to spec approval. `specs/spec-ticket-handoff-contract/spec.md`,
`ROADMAP.md`, `capture/**`, the concurrent AXI reslice, and every other spec remain
foreign during implementation.

## Out of scope

- Cross-harness portable story-line notation — this is a separate routing
  capability requiring `craft-line` and profile decisions: 3 edits, 1 gate run.
- Reducing `/bench-write-spec`'s mandatory skill preload — context economy should
  be measured after the artifact handoff is complete: 3 edits, 1 gate run.
- A general machine-readable story/seam/fence graph in the spec-build lifecycle —
  this would introduce a new parser and durable schema rather than repair the
  current guidance handoff: 10 edits, 1 promotion gate run.
