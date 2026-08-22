# Roadmap flow: the board stops growing faster than it retires

Status: ready

## Destination

Every drain reports the board's flow — rows opened, fed, retired, and open
mass — from one CLI source, and the drain's verdicts, the retro's improvement
items, and the row grammar carry the markers that keep the net flow at or
below zero over a window. The gate checks the markers. A one-time reviewed
pass brings the existing rows under the grammar.

## #1: What is the measured baseline?

Blocked by: none
Type: Research

### Question

Per drain, how many rows open, feed, and retire, and how does open mass move?

### Answer

See `decisions/assets/roadmap-flow-baseline.md`. Open mass rose 66 → 72 over
two weeks; a drain opens about 1.5 rows, feeds about 5, and retires under 1
unless it runs `--restructure`.

## #2: What outcome, and where does the pressure live?

Blocked by: none
Type: Grill

### Question

Count cap, retirement speed, or flow? Drain, row, or landing?

### Answer

Flow: each drain's exit reports opened, fed, retired, and open mass, and the
target is a net delta (opened minus retired) at or below zero over the window
in #6. The pressure lives in the drain's verdicts, and the landing side (retro
items) is in scope. The numbers come from a `bench roadmap` projection derived
from git history — one source per derived count; the drain quotes it.

## #3: What do the drain verdicts require?

Blocked by: #2
Type: Grill

### Question

When does an entry feed a row, open a row, or get dismissed?

### Answer

An entry feeds a row only when it changes that row's priority, scope, or
`Next:`; an entry that only adds an occurrence line is dismissed with one line
of why. A new row must carry a `Next:` token and a class (fix, feature,
decision-only); without a `Next:` the entry is dismissed or parked in
`capture/IDEAS.md`. A drained item that meets the light-path observables is
built in the drain session by default; the reviewer may decline. When the
window's net delta is positive, the next drain must propose reducing moves
(merge, fold, prune) in its batch diff; otherwise `--restructure` stays opt-in.

## #4: What does the retro require?

Blocked by: #2
Type: Grill

### Question

Cap the retro's recommendations, or change their test?

### Answer

No count cap. Each item under `## Agent-experience improvements` names the row
it feeds (`Feeds: FT<n> | new | none`) and the one sentence it changes in that
row's priority, scope, or `Next:`. An item that cannot is written as an
occurrence candidate, and the drain dismisses it by default.

## #5: What does the gate check, and what grammar does that need?

Blocked by: #3, #4
Type: Grill

### Question

Advisory or gate-enforced; which predicates; how does it relate to FT172?

### Answer

Gate-enforced from the start, on mechanical markers only: (a) every
`ROADMAP.md` row carries a `Next:` line whose value is one of
`shape | spec | ticket | decide | kit-edit`, mapping to `/bench-shape-idea`,
`/bench-write-spec`, a light-path ticket, a reviewer decision, and a
`craft-synthesis` edit; (b) every retro improvement item carries a `Feeds:`
marker. The "changes the row" judgment and the restructure trigger stay
advisory phase text. This work absorbs FT172's grammar half: it documents and
tests the row grammar as a contract; FT172 keeps only its `roadmap_id` half.

## #6: What window and target, and which commits count?

Blocked by: #1
Type: Grill

### Question

How wide is the window, and what is a flow event?

### Answer

A flow event is any commit that adds or deletes a `roadmap/FT<n>.md` file; the
report derives from file adds and deletes, not commit subjects. A drain commit
is a flow event that adds at least one row. The window is the commit range
spanning the last three drain commits; the target is opened minus retired, at
or below zero, summed over every flow event in that range. A positive sum
forces the #3 restructure proposal on the next drain.

## #7: What happens to the 72 existing rows?

Blocked by: #5
Type: Grill

### Question

Grandfather, or migrate?

### Answer

One reviewed `--restructure` drain proposes a `Next:` for each row; a row with
no honest next action moves to the parked section. The reviewer approves the
batch diff.

## #8: One spec or a split?

Blocked by: #3, #4, #5, #7
Type: Grill

### Question

Four separable behaviors: (A) flow report CLI, (B) drain verdicts, row grammar,
and gate markers, (C) retro change test, (D) migration pass.

### Answer

One spec, tickets in A → B → C → D order; they share one gate change and one
grammar.

## Not yet specified

## Spec-writer discretion

- The mechanical form of the flow report (`bench roadmap --flow` or a column in
  an existing projection) and its TOON shape, within `craft-cli`.
- Where the row-grammar contract text lives, as long as the parser and its doc
  share one source.

## Out of scope

- A hard cap on open rows.
- A shelf life or automatic expiry on rows.
- Limits on `capture/learnings.md` entries per landing; the journal keeps its
  existing dismiss path.
- FT172's `roadmap_id` half.

## Sources

- Path: `decisions/assets/roadmap-flow-baseline.md`
  Supports: #1 baseline and #6 window and target.
  Drift: re-measure after the flow report ships; the CLI replaces this hand count.
