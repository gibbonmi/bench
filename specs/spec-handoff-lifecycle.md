# Grill-to-spec handoff and spec lifecycle

Status: staged

Source map: `decisions/spec-handoff-lifecycle.md` (all nine tickets resolved;
the map deletes in story 11 once this spec is approved — this spec is its only
slice, and the spec carries every decision forward).

## Problem

A closed decision map hands the spec-writer prose conclusions but not the
structure — so a mid-tier spec session re-derives seams, invents boundaries the
grill already settled, or escalates to the top model for answers the map
already holds. Meanwhile `decisions/` (15 maps, all closed) and `specs/`
(39 files, 14 retired-in-place) accumulate history in the working tree against
invariant #3, and nothing owns noticing that a merged spec is still sitting
there.

## Solution

Every map closes with a nine-item `## Handoff` section the spec-writer reads
seams off, enforced by `bench maps` through the existing refuse-to-close loop.
Specs carry an in-file `Status: staged | implemented` line; merge is the
terminal trigger — promote what's durable, delete the file, git is the archive.
`bench status` owns the retirement nag with a full-code detector. The existing
backlog (all closed maps, all merged specs) is cleaned up in this build.

## User stories

1. As a reviewer closing a map, I want `/bench-shape-idea`'s Structure block to
   define the `## Handoff` section — the nine items from map ticket #1, each
   allowing `n/a — <one clause>`, plus a dependency-order line when a map
   yields multiple slices, required on every close — so the spec-writer reads
   seams off the map instead of inventing them.
   Line: claude-opus-4-8 / high. The content is fully decided in the map, but
   guidance prose is gate-blind beyond an anchor check, so it gets the mid
   model at high effort rather than cheap transcription.

2. As a closing session, I want `bench maps` to keep emitting a row for a map
   with zero open tickets whose `## Handoff` is missing or still carries a
   placeholder marker, so the exit's existing refuse-to-close loop physically
   blocks an unhandled close.
   Line: claude-sonnet-4-6 / medium. This is shell plumbing at a known parser
   seam whose exact rows the new AXI contract observes.

3. As the gate, I want an AXI contract proving the handoff rows: a zero-open
   map without `## Handoff` rows as missing, a filled Handoff stays silent, an
   open-ticket map emits no handoff row, and a fenced `## Handoff` example
   does not count — so the check can never rot into an always-pass.
   Line: claude-opus-4-8 / medium. Contract fragments are oracle content and
   the profile routes oracle correctness to mid effort.

4. As a session on the default branch, I want `bench status` to emit one
   low-severity row — "N merged spec(s) awaiting retirement →
   promote-then-delete (spec-retire)" — counting specs whose fence-skipped
   body carries a line-start `Status: implemented`, silent on any other
   branch and on specs with no Status line, so retirement has an owner that
   is full code with no LLM judgment.
   Line: claude-sonnet-4-6 / medium. It is one detector plus one ladder row at
   the existing status seam, fully observed by its contract.

5. As the gate, I want an AXI contract proving the retirement signal: fires on
   the default branch with a `Status: implemented` spec, silent off the
   default branch, silent when the marker is absent or only in a fenced
   block, and self-clears when the spec is deleted.
   Line: claude-opus-4-8 / medium. Same oracle-content reasoning as story 3.

6. As the spec-writer, I want `/bench-write-spec` rewired: step 2 reads seams
   off the map's Handoff (map-sourced seams are pre-agreed; pause only on
   deviations, uncertainty-flagged seams, or invented ones), step 8 becomes
   promote-then-delete with the `spec-retire: <name>` commit convention
   (replacing retire-in-place markers), the template gains the
   `Status: staged` line, and the exit recommends a NEW session on the mid
   tier to orchestrate the build — so the consumer side of the handoff and
   the spec lifecycle are in the command that runs them.
   Line: claude-opus-4-8 / high. This is the load-bearing guidance prose of
   the feature and its semantics are exactly what the gate cannot grade.

7. As the implementing session, I want `/bench-implement-spec`'s exit to flip
   the spec's `Status: staged` line to `Status: implemented` in the green-gate
   commit, so "implemented" honestly means awaiting review/merge.
   Line: claude-opus-4-8 / medium. A small prose promise, but its wording
   defines the lifecycle state machine, so it stays on mid.

8. As a defect session, I want `/bench-debug` to point at the git archive —
   check `git log --diff-filter=D -- specs/` (or `--grep=spec-retire`) for the
   feature's retired spec — so deleted specs stay discoverable cold.
   Line: claude-sonnet-4-6 / low. One anchored pointer line with an exact
   decided wording.

9. As the gate, I want layer-3 conformance anchors (plus a red canary each)
   for the four new prose promises — the Handoff template in shape-idea, the
   Handoff-consumption step in write-spec, the status flip in implement-spec,
   and the archaeology pointer in bench-debug — so the prose cannot silently
   drop out of the commands.
   Line: claude-opus-4-8 / medium. Anchor checks and canary fixtures are gate
   logic; a wrong anchor is an always-pass.

10. As the repo, I want the mechanical half of the cleanup: the ~14 specs
    already marked retired/superseded deleted outright under `spec-retire:`
    commits, with every dangling reference fixed in the same commits, so the
    gate's stale-reference sweep stays green through each step.
    Line: claude-sonnet-4-6 / low. Deletion plus reference fixes are fully
    observed by the existing sweep.

11. As the repo, I want the judgment half of the cleanup: each remaining
    implemented spec and each closed map (including
    `decisions/spec-handoff-lifecycle.md`) gets a promotion read — durable
    decision → ADR, hostile edge → the profile checklist, seam → the profile
    seam list, most reads ending in plain deletion because the kit already
    records its decisions — then deletes under the same commit convention;
    the parked `bench specs --retired` roadmap line is repointed from the map
    to this spec in the same sweep.
    Line: claude-opus-4-8 / medium. What counts as durable is semantic
    judgment the gate cannot see, so it runs on mid with each commit as
    reviewer veto surface.

12. As a session landing mid-transition, I want stories 10–11 sequenced before
    story 2 in the build order, so `bench maps` never spends the build nagging
    about 14 legacy closed maps that are about to be deleted.
    Line: claude-sonnet-4-6 / low. It is a build-ordering constraint, not new
    code.

## Implementation decisions

- **One detection core per fact.** The Handoff check extends the existing
  shared awk prelude discipline in the maps parser: `maps_rows` and
  `maps_unresolved_count` grow the same handoff rule (a zero-open-ticket map
  file missing a line-start `## Handoff` heading, or carrying a placeholder
  marker after it), so the `maps` listing and the `status` count cannot drift.
  The handoff row reuses the existing four-field schema:
  `<map>,handoff,handoff,missing` (or the placeholder's own state). No new
  subcommand, no schema change.
- **Retirement detector is fence-aware awk, positive-marker only.** It matches
  `Status: implemented` at line start outside ``` fences in `specs/*.md`, only
  when the current branch equals `default_branch()`. Absence of a Status line
  is silence by design — pre-convention specs and consumer repos never
  false-positive. Severity slots below every existing signal; the row self-
  clears on deletion.
- **Lifecycle states are two, in-file:** `staged` (written, approved, building)
  → `implemented` (green gate, awaiting review/merge) → file deleted on merge
  (promote-then-delete). No folders, no archive; `spec-retire: <name>` commits
  make `git log` the searchable archive.
- **Kit prose names tiers, never model ids.** The write-spec exit says "mid
  tier"; the binding stays in `projects/<name>.md`.
- **The gate checks presence, review checks meaning.** Layer-3 anchors assert
  the prose promises exist; the semantic quality of a written Handoff and the
  promotion judgment in story 11 land on the review axis and reviewer veto,
  not the gate (decided in the map, Handoff item 5).
- **Transition sequencing:** cleanup (stories 10–11) lands before the maps
  extension (story 2) so the new check meets a folder containing only maps
  that satisfy it.

## Testing decisions

- A good test here drives the CLI black-box in a throwaway fixture repo and
  asserts exit code, TOON stdout, and file/git state — never parser internals.
  Prior art: the "AXI maps unresolved-ticket contract" family in the AXI
  contract fragments, and the canary fixtures under `tests/canary/`.
- Seams tested: the `bench maps` contract, the `bench status` contract, and
  the kit-content conformance layer (gate layer 3 anchors). The cleanup is
  observed by the existing stale-reference sweep, not a new seam.
- Gate command: `.bench/gate.sh` (the project gate), which already runs the
  AXI fragments and canaries.

### Seam diagram — `bench maps` (handoff rows)

    trigger: closing session's exit loop / AXI contract in fixture repo
        │
        ▼
    decisions/*.md  ──▶  [ maps parser (shared prelude:      ]  ──▶  TOON rows
    (tickets +           [   fence-skip, CR-strip, markers,  ]       exit 0/1/2
     ## Handoff)         [   NEW: handoff-at-close rule)     ]
                              ◀ tests attach here: run `bench maps` against
                                fixture maps; assert exact handoff rows,
                                header count, and silence for filled Handoffs

### Seam diagram — `bench status` (retirement signal)

    trigger: SessionStart hook / reviewer / AXI contract in fixture repo
        │
        ▼
    specs/*.md ──▶  [ status renderer: severity ladder    ]  ──▶  ranked rows,
    git branch  ──▶ [   + NEW fence-aware implemented-    ]       ≤5 shown,
                    [     on-default-branch detector      ]       lead action
                              ◀ tests attach here: fixture repo on/off the
                                default branch; assert the row text, its
                                absence, and self-clearing after deletion

### Seam diagram — kit-content conformance (prose anchors)

    trigger: every gate run (layer 3) / canary run (layer 7)
        │
        ▼
    command files  ──▶  [ conformance checks: anchor     ]  ──▶  err() lines,
    (.agents/          [   present in shape-idea,        ]       gate red/green
     commands/*.md)    [   write-spec, implement-spec,   ]
                       [   bench-debug ]
                              ◀ tests attach here: canary fixtures drop each
                                anchor; the gate must go red with the
                                targeted substring

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | shape-idea Structure defines the nine-item Handoff, n/a rule, dependency-order line | conformance anchors | new layer-3 anchor check run before the prose lands — red | the anchor greps for the Handoff heading and item list; absent prose means red gate |
| 2 | zero-open map without `## Handoff` emits `<map>,handoff,handoff,missing`; filled Handoff silent | `bench maps` contract | story-3 contract run before the parser change — red | the contract asserts the exact row; today's parser emits nothing for such a map |
| 3 | open-ticket map emits no handoff row; fenced `## Handoff` ignored | `bench maps` contract | same contract, over-match assertions — red | asserts silence, so an over-eager rule (firing on open maps or fenced text) fails it |
| 4 | on default branch, spec with line-start `Status: implemented` yields the retirement row | `bench status` contract | story-5 contract run before the detector lands — red | today's status has no such row; the contract greps for its exact text |
| 5 | off default branch, or marker absent/fenced: no row; row gone after deletion | `bench status` contract | same contract, silence + self-clear assertions — red | a detector keyed on absence or ignoring fences/branch fails these assertions |
| 6 | write-spec carries Handoff-consumption step, spec-retire step 8, Status line in template, mid-tier exit | conformance anchors | new anchor check before the prose lands — red | each promise is a named anchor; dropping any one reds the gate |
| 7 | implement-spec exit flips staged → implemented in the green-gate commit | conformance anchors | new anchor check before the prose lands — red | anchor greps the flip instruction; silent removal reds the gate |
| 8 | bench-debug names the git-archaeology recipe | conformance anchors | new anchor check before the prose lands — red | anchor greps the `--diff-filter=D` pointer |
| 9 | each new anchor has a canary fixture that reds the gate with its targeted substring | canary (gate layer 7) | canary run with the fixture in place, check absent — red | proves each anchor bites; an always-pass anchor fails its canary |
| 10–11 | backlog deleted, references fixed, roadmap line repointed | stale-reference sweep + `bench status` | not red-first TDD — deletion has no failing pre-test; verified by the sweep staying green through each `spec-retire:` commit and the retirement row reaching zero | a missed reference reds the existing sweep; a missed spec keeps the status row nonzero |
| 12 | cleanup lands before the maps extension | build order | not TDD-able — ordering constraint; verified by commit sequence in review | out-of-order landing shows as a `bench maps` nag storm in the transcript |

### Edge inventory

- Handoff heading absent vs present-but-empty items — both covered (story 2/3
  rows; an empty section still carries placeholder markers or fails item
  presence).
- Placeholder (`— (open`) inside the Handoff — covered: the existing marker
  rule already applies once the file is scanned; the contract asserts the row.
- Fenced `## Handoff` or fenced `Status: implemented` example — covered
  (stories 3, 5): both detectors reuse fence-skipping.
- CRLF and missing trailing newline in maps/specs — covered by the shared
  prelude (CR-strip) and awk's last-line handling; existing CRLF contract is
  prior art.
- `decisions/` or `specs/` absent vs present-but-empty — covered: existing
  definitive-empty contracts; detectors return silence/empty table.
- Spec paths with spaces/globs — covered: existing space-path contract
  pattern extends to the status fixture.
- cwd deeper than repo root — covered by existing subdirectory
  root-resolution contract; both commands resolve the root already.
- Re-run idempotency of the cleanup — covered (story 10–11 row): a second
  sweep finds nothing to delete; the status row stays zero.
- Interrupted cleanup — covered by commit granularity: each `spec-retire:`
  commit is green-gated, so an interrupt leaves a green tree and a smaller
  backlog.
- **Won't handle:** `default_branch()` fallback robustness — pre-existing debt
  already parked on the roadmap (quality finding, 2026-07-02); the signal
  degrades to silence, never to a false gate verdict.
- **Won't handle:** semantic validation of Handoff item quality — gate-blind
  by design; the review axis owns it (map Handoff item 5).
- **Won't handle:** migrating consumer repos' existing specs to the Status
  convention — positive-marker keying keeps them silent; adoption is opt-in
  per repo.

## Out of scope

- `bench specs --retired` (git-archaeology wrapper) — a separate query
  capability parked on the roadmap behind the second-wave parser evidence
  bar — ~4 edits, ~2 gate runs.
- Adversarial gate pinning — separate threat model, already parked on the
  roadmap — ~6 edits, ~4 gate runs.
- README/HANDOFF inventory rot and the other 2026-07-02 quality findings —
  separate audit remediation already parked — ~10 edits, ~4 gate runs.
