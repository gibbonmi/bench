# Make the full-run options default

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, .agents/skills/bench-craft-line/SKILL.md, .bench/BENCH.md
Covers: none

## What to build

A `/bench-implement-spec --full <spec>` run takes `--delegate` and `--auto-approve` by default. The call turns one off with `--no-delegate` or `--no-auto-approve`.

`craft-line`'s "Delegated author lines" owns the delegated defaults:

- The default author limit is 3.
- Each review axis uses the cheap binding at high effort.
- Each diagnostic consultation goes to the mid binding and the cheap binding, both at high effort. A top-binding consultation at high effort breaks a disagreement.

A `--reviewer`, `--consultant`, or author-limit value in the call replaces the matching default. The defaults name tiers, not model ids, so each harness resolves them through `.bench/lines.env`.

`.bench/BENCH.md` states that a full run delegates by default.

## Acceptance

- [ ] The three files state the defaults, and only `craft-line` states the default lines and the author limit.
- [ ] `bench gate` passes at the landing.
