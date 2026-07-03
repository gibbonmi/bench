# Grill-to-spec handoff and spec lifecycle

Graduated from the roadmap (seam-output checklist for the grill-to-spec handoff,
parked 2026-07-03). The bootstrap grill grew the scope to spec slicing, spec
status, and folder hygiene — all overlapping on the same seams, so this map
yields **one spec** (reviewer call, #9). All tickets resolved 2026-07-03.

## #1: What must a closed map hand the spec-writer?

Blocked by: —
Type: Grill

### Question
The mid-tier spec-writer must read seams off the map instead of inventing them
or escalating. What does the map name?

### Answer
Nine items, recorded in the map's `## Handoff` section (#2):

1. Module boundaries plus the responsibility of each unit — inside vs outside.
2. Contract per boundary: inputs, outputs, exit codes, error posture (the
   observable interface).
3. Deep vs thin — which units hide complexity so the seam attaches at the
   interface; which are pass-throughs with no seam of their own.
4. What a black-box test can assert at each seam (exit code, TOON stdout,
   file/git state) without reaching inside.
5. Gate/oracle attachment — which seam the gate observes and how; flag any seam
   the gate cannot see (needs TDD or manual verify).
6. Hostile-input owner per seam: map each class from the profile's
   hostile-input checklist to the seam that owns it.
7. Uncertainty flag — any seam the grill could NOT settle, so the spec-writer
   escalates per the `craft-line` ladder instead of guessing on mid.
8. Rejected alternatives, so the spec-writer does not reopen closed decisions.
9. Domain watch-outs — hazards stated as domain facts for ANY reader
   (invariant #3), never model-addressed coaching; operating lessons go through
   `.bench/learnings.md` and `/bench-integrate-learnings`, never per-spec notes.

Framing: "here is the idea, the settled questions, and the seams — go build the
spec." On spec completion the spec-writer recommends a NEW session on the mid
tier to orchestrate the build — not the top model (no big-context top-model
iteration; escalate per-stage only). Kit prose names tiers abstractly; the
model binding stays in `projects/<name>.md`.

## #2: Where does the checklist live, and when is it required?

Blocked by: —
Type: Grill

### Question
Command exit prose, or map-side structure? Every map, or only code builds?

### Answer
A named `## Handoff` section in the decision-map template, defined once in
`/bench-shape-idea`'s Structure block; `/bench-write-spec` points at it (one
source per fact). The map is the only artifact the spec session loads, so the
seams must live in the map. Required on **every** close, with explicit
`n/a — <one clause>` per item that doesn't apply — a conditional rule is
lawyerable; explicit n/a keeps exclusions as decisions on the page.

## #3: Is the Handoff enforced?

Blocked by: —
Type: Grill

### Question
A required-but-unchecked section is a prose promise without a check.

### Answer
Yes — `bench maps` extends: a map with zero open tickets but a missing
`## Handoff` section, or handoff items still holding placeholders, keeps
showing a row. The exit's existing "refuse to close while the map shows a row"
loop then covers it; no new command or gate fragment. The check fires only on
zero-open-ticket maps — a map with open tickets is not a close candidate.

## #4: How does `/bench-write-spec` consume the Handoff?

Blocked by: —
Type: Grill

### Question
The producer side is pointless if the consumer isn't told to read it.

### Answer
Step 2 (Pick the seams) rewires: read seams off the map's Handoff first, verify
them against the current repo, invent seams only where the map is silent.
Map-sourced seams are **pre-agreed** (approved at map close) — the spec pauses
only on deviations, uncertainty-flagged seams (escalate per the `craft-line`
ladder), or seams it had to invent.

## #5: How do multiple outcomes slice into specs?

Blocked by: —
Type: Grill

### Question
When a map yields more than one buildable outcome, what bounds the slices?

### Answer
The Handoff's seam list is the menu slices are cut from; slicing stays the
reviewer's call. Hard rule: **one seam, one owner** — no two live specs claim
the same seam (shared seams produce coupled builds whose coverage maps and
shifts collide). Outcomes overlapping on a seam **merge into one spec** by
default — same-seam halves fail the separate-capability test, and splitting
them is the deferral the "just build it now" rule blocks. When multiple slices
do come out, the Handoff records a recommended dependency order — a
recommendation, not a decision. A genuinely separate capability sharing a seam
is sequenced instead: not staged until the owning spec implements.

## #6: Where does spec status live?

Blocked by: —
Type: Grill

### Question
Staging/implemented folders were proposed for findability in a crowded `specs/`.

### Answer
Flat `specs/`, with a `Status: staged | implemented` line in the file.
Location-as-status costs a file move per spec and rots every reference written
during shaping; the kit already keeps spec state in-file. Flip owner:
`/bench-implement-spec` flips staged → implemented in the green-gate exit
commit, where `implemented` means awaiting review/merge (preserving the spec's
post-hoc-veto role under a batch approval). Findability comes from the handoff
naming the spec path.

## #7: What is the terminal rule for maps and specs?

Blocked by: —
Type: Grill

### Question
Both `decisions/` (15 maps, all closed) and `specs/` (39 files, 14 already
retired-in-place) accumulate history in the working tree, against invariant #3.

### Answer
**Promote, then delete — no archive folder.** Anything durable is promoted
first (decision → ADR; hostile edge → the profile checklist; seam →
the profile's seam list), then the file is deleted; git is the archive. Only
the trigger differs: a map deletes once every slice it produced is staged as a
spec; a spec deletes once merged to the default branch. Superseded/retired
specs delete rather than carry markers (revises `/bench-write-spec` step 8).
Deletion commits use a `spec-retire: <name>` convention so
`git log --grep=spec-retire` / `git log --diff-filter=D -- specs/` is the
searchable archive, and `/bench-debug` gains a pointer telling defect sessions
to check that history. An archive folder was rejected: relocated crowding plus
silent rot, duplicating what git does. A `bench specs --retired` wrapper is
parked on the roadmap behind the second-wave parser evidence bar.

## #8: Who notices and executes retirement after merge?

Blocked by: —
Type: Grill

### Question
A terminal rule with no owner rots — the 14 retired-but-present specs prove it.

### Answer
An ambient `bench status` signal, **full code — no LLM judgment**: on the
default branch, any spec still carrying `Status: implemented` is by definition
merged-but-not-retired (implemented means awaiting merge). One low-severity
row: "N merged specs awaiting retirement → promote-then-delete." The next
session executes the sweep as routine hygiene; the signal self-clears.
Detection keys on the positive marker only — a spec with no `Status:` line
stays silent, so pre-convention specs and consumer repos never false-positive.

## #9: What does the one-time cleanup cover, and how does it ship?

Blocked by: —
Type: Grill

### Question
Does the existing backlog ride in this spec or a separate chore?

### Answer
In this spec (reviewer call). Scope: the ~14 retired/superseded specs delete
outright; the remaining implemented specs (all merged) each get a promotion
read then delete; the 15 closed maps promote-then-delete (this map lives until
its spec is staged). Every deletion trips the gate's stale-reference sweep —
reference cleanup is part of each delete, not a follow-up. The promotion reads
are the expensive story in the spec.

## Handoff

Idea: a closed decision map must hand the spec-writer the seams (the nine items
in #1), and spec/map artifacts get a full lifecycle: staged → implemented →
promoted-and-deleted, with `bench status` owning the retirement nag. One spec;
includes the one-time cleanup (#9).

1. **Module boundaries.** `bench maps` close-readiness check (owns map-format
   knowledge: open placeholders + Handoff presence on zero-open maps) ·
   `bench status` retirement signal (owns the merged-but-present rule; ranks on
   the existing severity ladder) · command prose edits: `/bench-shape-idea`
   (Handoff template + exit), `/bench-write-spec` (step 2 rewire, step 8
   delete-not-retire, `Status:` line in template, new-session-on-mid exit
   recommendation), `/bench-implement-spec` (status flip at green-gate exit),
   `/bench-debug` (git-archaeology pointer) · one-time cleanup sweep.
2. **Contracts.** Both CLI surfaces keep the AXI hybrid contract: flat TOON
   stdout, definitive empty states, structured stdout errors, exit 0/1/2.
   `bench maps` gains a row state for handoff-missing; `bench status` gains one
   low-severity signal row that never leads and respects the five-row budget
   and show-only-on-signal.
3. **Deep vs thin.** The `bench maps` parser is the deep unit — all map-format
   knowledge stays inside it (the gate and the exit consume its rows, never
   re-parse). The status signal is a thin rule plugged into the existing
   ladder — no seam of its own beyond the `bench status` contract. Command
   prose is not a code seam.
4. **Black-box assertables.** In throwaway fixture repos: a zero-open-ticket
   map without `## Handoff` → `maps` emits the row / with a filled Handoff →
   silent; a `Status: implemented` spec on the default branch → `status` emits
   the row / no `Status:` line → silent. Assert exit codes, TOON stdout, and
   (for the cleanup) file/git state.
5. **Gate attachment.** AXI contract fragments (gate layer 6) plus canary
   fixtures (layer 7) for both CLI changes; layer-3 conformance anchors for the
   command-prose promises; the stale-reference sweep observes the cleanup.
   Gate-blind: the semantic quality of a written Handoff, and the promotion
   judgment during cleanup — both land on review/reviewer veto, not the gate.
6. **Hostile-input owners.** `bench maps` owns: Handoff heading absent vs
   present-but-empty items, no trailing newline, `decisions/` absent vs empty.
   `bench status` owns: spec paths with spaces/globs, missing `Status:` line
   (must stay silent), cwd deeper than repo root. The cleanup owns: re-run
   idempotency (a second sweep is a no-op).
7. **Uncertainty flags.** None — no seam left unsettled by the grill.
8. **Rejected alternatives.** Exit-prose-only checklist (never reaches the spec
   session) · Handoff required only for code builds (lawyerable) ·
   folders-as-status and an archive folder (path churn, rotted refs, silent
   rot) · retire-in-place markers (superseded by delete) · LLM-judged
   retirement detection (full code decided) · building `bench specs --retired`
   now (parked behind the evidence bar).
9. **Domain watch-outs.** Kit prose names tiers, never model ids — the binding
   lives in `projects/<name>.md`. Retirement detection must key on the positive
   `Status: implemented` marker, never on absence. The Handoff check fires only
   on zero-open-ticket maps. Every deletion trips the stale-reference sweep, so
   reference cleanup ships inside each delete. Canary fixture repos hide
   dot-dirs behind the `dot-` prefix; new fixtures follow it.

Dependency order: single spec — n/a.
