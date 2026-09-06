# Roadmap flow: the board stops growing faster than it retires

Status: ready

## Destination

Every drain reports the board's flow — rows opened, fed, retired, and open
mass — from one CLI source. The drain's verdicts, the retro's improvement
items, and the row grammar carry the markers. The markers keep the net flow at
or below zero over a window. The gate checks the markers. A one-time reviewed
pass brings the existing rows under the grammar.

## Notes

## Decisions so far

- [What is the measured baseline?](roadmap-flow/tickets/1.md): See `decisions/roadmap-flow/assets/roadmap-flow-baseline.md`.
- [What outcome, and where does the pressure live?](roadmap-flow/tickets/2.md): Flow: each drain's exit reports opened, fed, retired, and open mass.
- [What do the drain verdicts require?](roadmap-flow/tickets/3.md): An entry feeds a row only when it changes that row's priority, scope, or `Next:`.
- [What does the retro require?](roadmap-flow/tickets/4.md): No count cap.
- [What does the gate check, and what grammar does that need?](roadmap-flow/tickets/5.md): Gate-enforced from the start, on mechanical markers only.
- [What window and target, and which commits count?](roadmap-flow/tickets/6.md): A flow event is any commit that adds or deletes a `roadmap/FT<n>.md` file.
- [What happens to the 72 existing rows?](roadmap-flow/tickets/7.md): One reviewed `--restructure` drain proposes a `Next:` for each row.
- [One spec or a split?](roadmap-flow/tickets/8.md): One spec, tickets in A → B → C → D order; they share one gate change and one grammar.

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

- Path: `decisions/roadmap-flow/assets/roadmap-flow-baseline.md`
  Supports: #1 baseline and #6 window and target.
  Drift: re-measure after the flow report ships; the CLI replaces this hand count.
