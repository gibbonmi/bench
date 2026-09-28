# 3. Give each conformance probe its own Bench home

Blocked by: none
Writes: internal/conformance/checks_test.go, internal/conformance/conformance_env_test.go (new), internal/conformance/entry_point_parity_test.go, internal/conformance/entry_point_parity_bite_test.go, internal/conformance/cross_compile_stress_test.go, internal/conformance/harness_test.go, internal/conformance/fixture_bite_test.go, internal/conformance/line_routing_exec_test.go
Covers: TD20, TD21

## What to build

Chunk: TD-C2.

The conformance subprocess environment gives each call its own Bench home under a directory that the call owns. The fixed shared name under `TMPDIR` goes away. The npm cache stays the one declared shared cache, under its fixed name, when the base environment names no npm cache.

Each caller handles allocation failure and owns cleanup of its private home.
The environment helper owns optional PATH composition.
Environment helpers and their tests stay together.
The existing parity omission test moves to the parity fixture file to preserve the file growth bound.
The stress-tagged caller also compiles and runs its existing matrix test.

## Acceptance

- [ ] Two calls give two different Bench home values, and neither value is the old fixed name.
- [ ] With no npm cache in the base environment, the call sets the npm cache to the fixed shared name under `TMPDIR`.
