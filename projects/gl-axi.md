# Project: gl-axi

An agent-optimized GitLab CLI: an AXI-conformant wrapper over `glab`, designed so
an agent gets higher accuracy at lower token cost than driving raw `glab`. The
`axi` skill is the design spec; this profile makes conformance a gate check.

## Seams (test here)

- **Output boundary.** The single place internal JSON becomes TOON on stdout. Test
  through it: given a backend payload, assert TOON shape, minimal default schema,
  truncation-with-total, and inline aggregates. This is the AXI surface and the
  highest seam — keep the conversion in one module so everything upstream stays on
  JSON and free to change.
- **`glab` adapter.** The wrapper that calls `glab` and normalizes results/errors.
  Test that dependency errors are translated to structured stdout errors (no
  leaked `glab` stack traces) and that mutations are idempotent.
- **Command surface.** Flag parsing and exit codes. Test that every operation
  completes from flags alone (no interactive prompts) and that exit codes follow
  0/1/2 (incl. no-op = 0).

## Gate (`.bench/gate`)

```
pytest -q && axi-conformance ./gl-axi && bench-glab-delta --assert
```

Three layers of oracle, in order of authority:

1. **Tests** — behavior at the seams above.
2. **`axi-conformance`** — asserts the seven AXI principles mechanically: TOON
   stdout, ≤4-field default list schemas, structured stdout errors, correct exit
   codes, definitive empty states. A regression here fails the build like a broken
   test. (This is your paired-delta instinct turned into a conformance gate.)
3. **`bench-glab-delta --assert`** — your paired per-task harness vs raw `glab`,
   with deterministic assertions and a non-regression threshold on the token/
   accuracy delta. The build is not done if gl-axi got worse than `glab` on a task.

Prefer deterministic assertions over an LLM judge throughout — a TOON-shape check
is a parser, not a model.

## Lines (model + effort routing)

- **TOON serialization / field-projection logic** → mid–high effort. Get the
  deterministic field projection right *first*, then serialize; that ordering is a
  closed decision.
- **New command wrappers, flag plumbing** → cheap model, low effort. These are
  mechanical once the output boundary exists.
- **Conformance-harness work** → mid effort; correctness of the oracle matters
  more than speed.

## Notes for cold sessions

- The pipeline is deterministic field projection → TOON serialization, in that
  order. Don't serialize the full object and trim after.
- gl-axi ships as a skill plus an optional session-start hook (ambient repo
  dashboard), mirroring the AXI family pattern.
