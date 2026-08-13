# Expose the frozen source range to diff and preflight

Blocked by: none
Writes: internal/diff, internal/preflight, internal/axi, bin/bench.sh

## What to build

Add `--base <commit>` to live `bench diff [--full]` and to
`bench preflight build|review <slug>`. One typed source-range derivation resolves
the full base and source-tip identities, requires base ancestry, supplies the
committed path inventory to both callers, and detects movement. Live diff and
build preflight add index, tracked-worktree, and untracked state; review preflight
requires a clean source. `--base` and `--commit` are mutually exclusive, accepted
explicit-base operations never write Git config, and all existing bare and
`--commit` behavior stays compatible.

The shared ancestry derivation treats a base equal to the source tip as an
inclusive ancestor with an empty committed range; it refuses only a base that is
not an ancestor. Preflight may still grade that empty range red under its own
nonempty-build policy, while reauthorization consumes the inclusive ancestry
fact directly.

The public proof includes the retained FT198 graph shape: review from `0924e02e`
through its source tip, retain every source-authored path, and exclude later
destination-only phase state. This ticket lands green without changing the
landing command.

## Acceptance

- [ ] Explicit-base diff and preflight report the complete live source snapshot,
      full identities, and unchanged Git config bytes (covers PL1).
- [ ] The FT198-shaped public preflight owns exactly the reviewed source range and
      excludes destination-only handoff history (covers PL2).
- [ ] Missing, malformed, ambiguous, non-ancestor, and mutually exclusive bases
      refuse without config or ref writes (covers PL3).
- [ ] Clean review pins the exact tip and refuses a dirty source (covers the
      review half of PL4); one retry followed by the existing drift refusal
      covers repository movement (covers PL5).
