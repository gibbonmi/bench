# Add the setup and specs-staged signals (and the flagged decisions-ready state)

Blocked by: 02-normalize-actions.md
Writes: internal/status, internal/maps

## What to build

`setup` row: severity 0, appended before the gate row, when the git root has no
`.bench/` directory — detail `no .bench/`, command `bench setup`. `specs` staged state:
severity 4, appended before the drain row, counting `Status: staged` through the spec
Status reader; command `/bench-implement-spec specs/<slug>/spec.md` for one, bare for
several; a paired `reviews/<slug>.md` changes nothing; retirement keeps its own row.
Reviewer-flagged: `decisions` row gains `N ready map(s)` → `/bench-write-spec
decisions/<map>.md` (one) / `/bench-write-spec` (several) only when no shaping or invalid
map is active — build it behind its own commit so the reviewer can strike it.

Covers: R14, R15, R16, R17, R18, R19, R20, R21, R22, R23

## Acceptance

- [ ] Un-adopted fixture: board lead is `▶ bench setup  (setup)`; adopted fixture: no setup row.
- [ ] One staged spec renders `specs  1 staged spec(s)  → /bench-implement-spec specs/<s>/spec.md`; two render the bare command; staged+implemented render two rows, staged first; staged orders after guards and before drain.
- [ ] Spec dir without spec.md, fenced Status, and a FIFO at `specs/` count zero and return.
- [ ] Ready-only decisions render `1 ready map(s)` → `/bench-write-spec decisions/<m>.md`; ready+shaping keeps `/bench-shape-idea`.
- [ ] Gate green.
