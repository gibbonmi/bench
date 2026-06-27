# Project: Regroup

Hockey video-analysis tool. Core loop: an assistant coach tagging events during
film review. The quality metric is **coding velocity with eyes on video** — a
practiced user tags a complete event (type + outcome + coordinate) in under ~2s
without looking away from playback. Stack: DuckDB + Parquet + Pydantic.

## Seams (test here; everything else is free to change)

- **Phase/game-state machine.** The six-state phase taxonomy with possessing-team
  frames and two-axis battle outcomes (strength as modifier) is the domain core.
  Test through the state-transition interface — feed event sequences, assert
  resulting phases and outcomes. Highest-value seam; a wrong transition is the
  worst class of bug.
- **`CoordinateProvider`.** The rink-coordinate abstraction. Test through its
  interface so capture surface and storage format can both vary behind it.
- **Event store.** Append-only over DuckDB/Parquet; corrections are new events.
  Test ingestion and query at the repository interface, not against raw SQL.

UI is **not** a unit-test seam. It is gated by the design system, the
`regroup-ui` screenshot loop, and the interaction-state review instead.

## Design system (separate repo)

Visual decisions live in the **separate Regroup design repo** (tokens + canonical
component inventory), consumed here as a submodule / package / pinned path. See the
`design-system` skill. It is the source of truth for every color, spacing, type,
and motion value, and for every component. UI work **composes** from it and
**references tokens** — it never hardcodes a value or regenerates a canonical
component. A missing token or variant is a design decision: add it in the design
repo (via Claude Design when you're in a Claude session, or by editing the repo
directly in any harness), commit, re-pin, then resume the build. Because the source
of truth is a repo of plain artifacts, this works identically under Claude Code,
Codex, or any other harness.

Regroup-specific conformance rules the `design-conformance` gate enforces, on top of
the generic raw-value/canonical-import scan: no `#FFFFFF` (or near-white) adjacent to
the player surface; the shuttle slider, event chips, and timeline markers are
canonical components composed, never restyled; the `regroup-ui` interaction skill and
its screenshot loop are part of UI done-ness alongside this visual layer.

## Gate (`.bench/gate`)

Logic shifts:
```
mypy regroup && pytest -q && ruff check regroup
```

UI shifts — the logic gate is necessary but **not sufficient**. Also required:
- **design-conformance** — no raw hex / hardcoded px / literal durations (values
  must reference tokens), no `#FFFFFF` near the player, new components import from
  the canonical inventory rather than duplicating it. (You implement this lint,
  like gl-axi's `axi-conformance`.)
- **`regroup-ui` screenshot loop** — the two chasms (decision chasm, chasm of
  understanding) and the five interaction states (hover/press/in-flight/success/
  failure) are review gates a green suite can't see.

Never regenerate the canonical shuttle slider — compose it.

## Lines (model + effort routing)

- **Phase/possession ontology, battle-outcome logic** → top model, high effort.
  The genuinely uncertain seam; spend here.
- **CoordinateProvider plumbing, Pydantic models, store wiring** → cheap model,
  low–medium effort at the known seam.
- **UI components** → mid model with `design-system` + `regroup-ui` + the
  screenshot loop. The loop and the conformance gate, not raw model strength, are
  what catch the failures here.

## Notes for cold sessions

- Read the **design repo** (tokens + canonical components) before generating any
  UI; read `decisions/` for the current
  domain model before proposing changes to it.
- States before pixels: enumerate the interaction state machine before writing any
  interactive component (idle → coding-armed → coordinate-pending → committed →
  undo-window).
- The phase taxonomy is a closed decision; don't reopen it inside a build.
