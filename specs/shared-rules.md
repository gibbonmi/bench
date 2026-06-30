# single-source the shared platform rules

## Problem

The four invariants and the communication rules live in two hand-maintained places at
once: full in benchkit's `AGENTS.md` ("The four invariants", "How to talk to me") and
condensed in `.bench/BENCH.md` (the guide `bench link` ships to consumers). Edit one and
you must remember to mirror the other — which already failed this session: I edited
communication in `AGENTS.md` and had to hand-copy it to `BENCH.md`. The rules that define
how bench works can silently diverge between the dev repo and every consumer.

## Solution

Make `.bench/BENCH.md` the single canonical home of the four invariants and the
communication rules, at full richness. `AGENTS.md` stops restating them and instead
references `BENCH.md` as the source. One edit, one place, and the gate forbids
re-duplication. Consumers already point at `BENCH.md`, so they inherit the same rules the
dev repo runs on — benchkit dogfoods its own shipped guide.

## User stories

1. As a maintainer, I want the four invariants and the communication rules in exactly one
   place, so a single edit reaches both the dev repo and every consumer.
2. As a maintainer, I want `.bench/BENCH.md` to hold those rules at full richness — the
   invariant explanations move intact, not condensed away.
3. As a maintainer, I want `AGENTS.md` to reference `BENCH.md` for those rules instead of
   restating them, so the two can't drift.
4. As an agent in any AGENTS.md harness, I want `AGENTS.md` to state plainly that the
   shared rules are canonical in `.bench/BENCH.md` and that I must read it.
5. As an agent in benchkit's repo on Claude Code, I want the canonical rules to still load
   into context — `CLAUDE.md` imports `.bench/BENCH.md` alongside `AGENTS.md` — so
   single-sourcing doesn't cost me the rules.
6. As a consumer who ran `bench link`, I want the full invariants + communication in
   `BENCH.md` (which I already point at), so I inherit the dev repo's actual rules.
7. As a maintainer, I want the gate to fail if `AGENTS.md` re-duplicates the invariants or
   communication, so reintroduced drift can't ship.
8. As a maintainer, I want the gate to fail if `BENCH.md` is missing the canonical rules,
   so the single source can't be gutted.
9. As a maintainer, I want the gate to fail if `AGENTS.md` drops its pointer to
   `BENCH.md`, so the rules can't become unreachable.
10. As a maintainer, I want a canary fixture proving the drift check bites.
11. As a maintainer, I want stale cross-references updated — e.g. `CONTEXT.md` currently
    says invariants live "in AGENTS.md" — so no doc points at the old home.

## Implementation decisions

- **`.bench/BENCH.md` is canonical.** Its current condensed "Core Rules" becomes the full
  four invariants (migrated verbatim from `AGENTS.md`), and "Communication" holds the
  full communication rules. This is the one source.
- **`AGENTS.md` references, doesn't restate.** Replace its "## The four invariants" and
  "## How to talk to me" sections with one short pointer section: these shared platform
  rules are canonical in `.bench/BENCH.md`; read it; don't copy them back here. Keep the
  harness-neutral prose pointer — no Claude-only `@import` in `AGENTS.md`.
- **`CLAUDE.md` loads both.** It already imports `@AGENTS.md`; add `@.bench/BENCH.md` so
  benchkit's own Claude sessions still get the canonical rules in context. `CLAUDE.md` is
  the Claude adapter, so the import lives there, not in the neutral `AGENTS.md`.
- **Consumer path is unchanged.** The link-installed managed block in a consumer's
  `AGENTS.md` already points to `.bench/BENCH.md`; it now simply finds fuller rules there.
  No change to `bench link`.
- **Drift is gate-enforced.** A new kit-conformance check encodes the single-source
  property directly: each canonical marker (the four invariant headlines + a communication
  marker) must appear in `BENCH.md` and **not** in `AGENTS.md`, and `AGENTS.md` must carry
  the `.bench/BENCH.md` pointer.
- **Fix the cross-reference.** Update `CONTEXT.md`'s "invariant" definition (and any like
  it) that names `AGENTS.md` as the invariants' home.

## Testing decisions

- **A good test** asserts the single-source structural property on the repo's real
  `AGENTS.md` + `BENCH.md` — external structure, not internals.
- **Seam:** a kit-conformance check in `.bench/gate.sh`, grepping the two files directly.
  **Prior art:** the existing 5a/5b/5c checks already grep `AGENTS.md` for skill/command
  index↔disk sync — same shape, same section.
- **The check asserts:** (a) each canonical marker is present in `BENCH.md`; (b) no
  canonical marker appears in `AGENTS.md` (the anti-duplication / drift guard); (c)
  `AGENTS.md` contains the `.bench/BENCH.md` pointer. Markers are the four invariant
  headlines ("The gate is the oracle…", "Declare the line…", "Document for the
  teammate…", "One small change at a time…") plus a stable communication phrase.
- **Canary fixture** (`tests/canary/`): an `AGENTS.md` that copies an invariant back in
  must make the gate go red, proving the drift guard bites. Follow the existing fixture
  shape; targeted `EXPECT` substring.
- **Gate command:** the project gate — `.bench/gate.sh` (`bench gate`).

## Out of scope

- **Single-sourcing the other overlapping sections** — "How the pieces fit" / "Process
  proportionality" (AGENTS.md) and "Files / Workflow / Commands / Hook Layers"
  (BENCH.md). A separate capability (more sections, possibly different granularity calls)
  with its own future spec; you scoped this to invariants + communication. Est ~30–45 min.
- **A `bench sync` command / managed-block embedding** to auto-propagate identical text
  into both files (the rejected mechanism A). Option C needs none of it — there's one
  copy. Not "the rest of this feature"; it's a different mechanism we chose against.
