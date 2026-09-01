# Show the single-file example in gate-prose help

Blocked by: none
Writes: internal/gate/gate_prose.go, internal/gate/gate_prose_test.go, CHANGELOG.md
Covers: none

## What to build

`bench gate-prose --help` shows the grammar line and one example of the single-file
form. A caller who names one file as the first operand gets a refusal, because that
operand is the root. The help text does not show the working form, so the caller
learns the shape only by a refusal. The usage text becomes two lines: the grammar
line stays as it is, and a second line reads `example: bench gate-prose . -- <path>`.
Each site that prints the usage text — the `--help` path, the usage-error path, and
the non-directory-root path — prints both lines.

## Acceptance

- [ ] `bench gate-prose --help` prints the grammar line and the example line to stdout at exit 0.
- [ ] A usage error prints both lines to stderr at exit 2.
- [ ] The `internal/gate` gate-prose tests stay green.
