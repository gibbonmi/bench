# Grill continuation across tickets in shape-idea resume mode

Status: draft

## Problem

FT41, from the learnings journal: a session running `/bench-shape-idea` in
resume mode resolved the one named ticket, then stopped and asked the reviewer
to re-prompt before grilling the next ticket — a one-question ticket the same
sitting could have carried straight through. The reviewer was present and
answering; the pause bought nothing and cost a round-trip.

The root cause is the wording in `.agents/commands/bench-shape-idea.md`. The
**Resume** mode says:

> **Resume** (map + ticket number in): load the whole map, resolve that one
> ticket, record the answer in its body (current state only), add any
> newly-discovered tickets with correct `Blocked by` edges, **then stop**. If a
> resolution invalidates other tickets, update or delete them.

The flat "resolve that one ticket … then stop" reads as one-ticket-per-invocation
even when resolving that ticket unblocks another whose single question the
present reviewer could answer in the same breath. The command's Exit handoff
reinforces it, recommending "another focused `/bench-shape-idea` resume on the
next open ticket" as the next step.

The `craft-grill` skill already carries the correct stop rule — "Stop when the
fog is gone: when the remaining questions no longer change what gets built" —
so it is *aligned*, not conflicting. The defect lives only in the shape-idea
command's resume prose.

## Solution

Soften the Resume-mode stop language so a running grill continues within the
same sitting. Concretely, replace the flat "then stop" with a continuation rule
while keeping both legitimate stop conditions:

- **Continue while the reviewer is present and answering.** After resolving the
  named ticket, carry the grill straight into any ticket the resolution just
  unblocked whose question the present reviewer can settle now — question after
  question, never pausing for permission or a re-prompt.
- **Preserve the two real stops.** Stop when no unresolved question the present
  reviewer can answer remains (the fog is gone — `craft-grill`'s existing rule),
  or when the reviewer stops answering (goes AFK / leaves the sitting). A ticket
  that needs external evidence — a Research probe rather than a reviewer
  answer — is not force-grilled: record it open with its `Blocked by` edges and
  move to another answerable ticket, or stop.
- **Everything else in Resume mode stands.** Record each answer as current
  state only, add newly-discovered tickets with correct edges,
  invalidate/update/delete tickets a resolution supersedes, and run
  `bench maps` before declaring the map closed.

Bootstrap mode is untouched: "do not also resolve tickets" is a deliberate
one-session-builds-the-map rule, not the defect. `craft-grill` is untouched: it
already stops on fog-gone.

This is leverage kit prose loaded by every shape-idea resume session, so it
follows `craft-synthesis` (respect closed decisions, three quality loops,
propose don't merge) and takes the `craft-line` leverage override.

## User stories

1. As a session running `/bench-shape-idea` in resume mode with the reviewer
   present and answering, I want the command's Resume-mode prose to tell me to
   continue question-after-question into newly-unblocked tickets — never
   pausing for permission or a re-prompt — and to stop only when no unresolved
   question the present reviewer can answer remains or the reviewer stops
   answering, so one sitting carries the grill as far as the reviewer's
   presence allows instead of stalling for a re-prompt after every ticket.
   Line: claude-fable-5 / high. This is guidance prose loaded by every
   shape-idea resume session, so it takes `craft-line`'s leverage override of
   top model and high effort even though the diff is small.

## Implementation decisions

- **One file, one prose edit: `.agents/commands/bench-shape-idea.md`, the
  `## Two modes` → **Resume** paragraph.** Replace the trailing "then stop"
  clause with the continuation rule and its two preserved stops. Keep the
  Bootstrap paragraph and its "do not also resolve tickets" line exactly as
  written. No `craft-grill` edit — its fog-gone stop already matches.
- **Enforce the new rule with a gate anchor, matching the file's existing
  anchor pattern.** The shape-idea command already carries four positive
  conformance anchors (`## Handoff`, `Hostile-input owner`, `Dependency order`,
  `n/a —`) plus a negative bypass check, each pinned by a canary. Add one
  anchor pinning a distinctive, load-bearing phrase from the new prose (e.g.
  `never pause for permission or a re-prompt`) via `checkWorkflowAnchors` in
  `internal/conformance/docs_workflow_helpers_test.go`, with a meaningful
  diagnostic (the phrase alone is generic; use a `requireCollapsed`-style diag
  naming the continuation rule). The rule is exactly the leverage prose whose
  *absence* caused FT41, so pinning it against silent future deletion is
  proportionate — and it composes the file's existing seam rather than
  inventing one.
- **The anchor's bite proof is a canary fixture.** Add
  `tests/canary/workflow-guidance-anchors/shape-idea-grill-continuation/` — a
  mutated copy of `bench-shape-idea.md` with the continuation phrase dropped —
  plus its `EXPECT` file carrying the new diagnostic string, registered in
  `internal/conformance/registry_test.go`
  (`conformanceFixture(".bench/gate-docs-contracts.sh")`) and added to the bite
  list in `TestDocsCurrencyTokenDietAndWorkflowFixturesBite`
  (`internal/conformance/fixture_bite_test.go`). This is a test-only addition;
  no CLI, hook, or adapter code changes.
- **No existing anchor moves.** The edit is in the Resume paragraph, which no
  current anchor touches; every pinned phrase and the bypass negative-check
  survive unchanged.
- **Reviewer down-scope option (recorded, not chosen here):** if the reviewer
  judges the new anchor as sediment, the story can ship prose-only, verified by
  `rg -F` grep plus `craft-skills` review, with the existing shape-idea anchors
  as the collateral guard — the same posture the retired `kit-rule-edits` spec
  took. The recommended path adds the anchor because it converts a
  not-TDD-able prose change into one with a genuine red signal.

## Testing decisions

- **What a good test is here.** The artifact is command prose; the observable
  is the phrase's presence in the live file, read exactly as the docs
  conformance scan reads it, plus a canary that goes red when the phrase is
  absent. Assert the continuation rule's presence — do not re-derive its
  wording, and do not pretend prose presence is behavioral TDD of the grill
  loop.
- **Which seams get tested, and prior art.** The seam is the docs conformance
  scan (`internal/conformance`, `RunConformance` → `checkWorkflowAnchors`,
  which `strings.Contains`-checks anchors like `## Handoff` and the bypass
  fragment on this same file) plus its canary bite test. Prior art: the
  `shape-idea-handoff-anchor` and `shape-idea-bypass` fixtures under
  `tests/canary/workflow-guidance-anchors/` — the new fixture follows their
  shape (a `files/dot-agents/commands/bench-shape-idea.md` copy + an `EXPECT`
  diagnostic).
- **Gate command.** `.bench/gate.sh`. Adding a conformance anchor touches the
  gate layer, so per `craft-synthesis`'s dogfood proportionality this is not a
  prose-only change — it dogfoods: a real `/bench-shape-idea` resume on the
  changed kit plus a green gate, not just a read.

### Seam diagram

    trigger: docs conformance scan (RunConformance → checkWorkflowAnchors)
             + canary bite test + every /bench-shape-idea resume session
        │
        ▼
    edit  ──▶  [ .agents/commands/bench-shape-idea.md  ]  ──▶  Resume-mode continuation rule
    fixture ─▶ [   (## Two modes, Resume paragraph)     ]  ──▶  pinned phrase present
                      ◀ tests attach here: new require() reads the live file;
                        canary fixture (phrase dropped) asserts the diag under bite test

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Resume-mode prose carries the continuation rule (continue into newly-unblocked tickets, never pause for a re-prompt) as a gate-pinned anchor | `checkWorkflowAnchors` in `docs_workflow_helpers_test.go` + canary bite | New `shape-idea-grill-continuation` fixture (phrase dropped) added to `TestDocsCurrencyTokenDietAndWorkflowFixturesBite` first goes red — `RunConformance` emits no diagnostic until the `require()` anchor is added; adding the anchor greens it. Genuine red→green | With no anchor, a future edit could silently drop the continuation rule and the gate would stay green; the fixture proves the anchor bites when the phrase is absent, and the live-file check keeps the rule on the page |
| 1 | The edit breaks no existing shape-idea enforcement | `.bench/gate.sh` (docs conformance scan + canary meta) | `.bench/gate.sh` exits 0 today and must stay 0 — the `## Handoff`, `Hostile-input owner`, `Dependency order`, `n/a —` anchors and the `straight to /bench-write-spec` bypass negative-check all remain satisfied | The edit sits in the Resume paragraph, which none of those anchors touch; this row guards against an edit that accidentally reintroduces the bypass fragment or drops a pinned phrase, which would flip the exit code |

Degenerate-implementation check: the cheapest wrong build ships the prose but
omits the anchor+fixture, or omits one preserved stop condition. Row 1's bite
goes red if the anchor is missing; the preserved-stop clauses (AFK,
external-evidence, fog-gone) are not gate-observable and are confirmed by
`craft-skills`/`craft-review` reading — stated here, not pretended into a red
signal.

### Edge inventory

Walked per the story's behavior:

- **Reviewer goes AFK mid-grill** → coverage by prose (the "reviewer stops
  answering / leaves" stop). Continuation is conditioned on the reviewer being
  present and answering, so an AFK reviewer ends the grill cleanly. Semantics —
  verified by review, not gate.
- **Next unblocked ticket needs external evidence (a Research probe, not a
  reviewer answer)** → coverage by prose. Such a ticket is recorded open with
  its `Blocked by` edges, not force-grilled; continuation applies only to
  questions the present reviewer can settle now. Verified by review.
- **All tickets resolved** → coverage by preserved fog-gone stop. No unresolved
  question remains, so continuation ends and the map becomes closeable (still
  gated on `bench maps` before close). Verified by review.
- **Re-run idempotency of the prose edit** — Won't handle: a single add-only
  text change with no runtime surface; a double-insert is caught by review, not
  gate.
- **Bootstrap mode altered by proximity** — Won't handle: the edit is scoped to
  the Resume paragraph; Bootstrap's "do not also resolve tickets" line is left
  verbatim.
- **`craft-grill` drift** — Won't handle: the skill already carries the
  fog-gone stop, so it needs no edit; nothing to keep in sync.

## Out of scope

- **Any behavior change to Bootstrap mode or to `craft-grill`** — not a
  separate capability, simply not the defect: Bootstrap's one-session rule and
  the skill's fog-gone stop are already correct. Excluded, no estimate.
- **A CLI or `bench maps` change to detect/annotate a mid-sitting stop** — a
  distinct future capability (tooling that reasons about grill continuation
  state), not the rest of this prose fix; build later as ~2-3 edits + a bite
  fixture, ~2 gate runs. Not needed for the rule to hold.
- **Prose-only variant without the new anchor** — the recorded reviewer
  down-scope, not a separate capability; if chosen it *shrinks* this spec (drop
  the fixture and `require()` line, verify by grep + review), so it lives in
  Implementation decisions as the alternative, not as deferred future work.
