# 3. Give each conformance probe its own Bench home

Blocked by: none
Writes: internal/conformance/checks_test.go, internal/conformance/conformance_env_test.go (new)
Covers: TD20, TD21

## What to build

Chunk: TD-C2.

The conformance subprocess environment gives each call its own Bench home under a directory that the call owns. The fixed shared name under `TMPDIR` goes away. The npm cache stays the one declared shared cache, under its fixed name, when the base environment names no npm cache.

## Acceptance

- [ ] Two calls give two different Bench home values, and neither value is the old fixed name.
- [ ] With no npm cache in the base environment, the call sets the npm cache to the fixed shared name under `TMPDIR`.
