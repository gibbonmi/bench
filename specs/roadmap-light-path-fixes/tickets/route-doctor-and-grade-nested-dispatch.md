# Route doctor and grade nested dispatch

Blocked by: none
Writes: internal/adopt/doctor.go, internal/adopt/adopt_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF7

## What to build

Route bench doctor through usage.Parse. Correct whyNested so it describes
leaf-owned grammars, and grade the doctor leaf through that routing census.

## Acceptance

- [ ] Doctor help exits zero on stdout.
- [ ] Doctor rejects unknown flags through parser-standard stderr.
- [ ] The routing census uses the corrected nested-dispatch reason.

