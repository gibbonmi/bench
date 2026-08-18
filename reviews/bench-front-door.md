## Standards

Finding count: 3. Worst: ask-user.

1. **ask-user** — The status action parser duplicates the command strings and grammars emitted by the producers, violating the one-source-per-fact standard: `internal/status/status.go:83-133` re-derives actions produced at `internal/status/status.go:235`, `internal/status/status.go:252`, `internal/status/status.go:364-391`, `internal/status/status.go:453-472`, `internal/status/status.go:496-573`, `internal/status/status.go:605`, and `internal/status/status.go:663-778`.
2. **ask-user** — The executable registry definitions at `cmd/bench/main.go:70-92` and `cmd/bench/main.go:144-173` duplicate command names and grammars in the help literal at `cmd/bench/main.go:93-142`, contrary to the one-source-per-fact standard.
3. **auto-fix** — `.bench/BENCH-reference.md:103-106` still claims `.bench/BENCH.md` owns the canonical CLI inventory, but `.bench/BENCH.md:36-42` and `projects/benchkit.md:25-28` make `bench help`, rendered from `commandRegistry`, canonical.

## Spec

Finding count: 1. Worst: auto-fix.

1. **auto-fix** — Story 33 requires direct binary `help`, `--help`, and `-h` to print the same inventory (`specs/bench-front-door/spec.md:102`), and R38 requires binary/wrapper inventory equality (`specs/bench-front-door/spec.md:236`), but the direct dispatcher defaults only an empty argument list at `cmd/bench/command_registry.go:74` and registers only the literal `help` name at `cmd/bench/main.go:94`; `cmd/bench/main_test.go:124-143` invokes only literal `help` on the binary and exercises the aliases only through the wrapper.

## Coverage

Finding count: 3. Worst: high.

1. **high · auto-fix** — Direct binary `--help` and `-h` are untested even though Story 33 and R38 require them (`specs/bench-front-door/spec.md:102`, `specs/bench-front-door/spec.md:236`): `cmd/bench/main_test.go:129-143` invokes the binary only with literal `help`, so either alias can return unknown-subcommand without a red test. This is the same repair target as the Spec finding, not a second repair.
2. **high · auto-fix** — Raw control bytes in a staged-spec slug or ready-map path can enter runner-up actions through the producer paths at `internal/status/status.go:249` and `internal/status/status.go:708`, then reach the raw `also:` formatter at `internal/status/status.go:314` without sanitization; no test drives those producer-derived inputs and asserts a safe refusal or escaped rendering.
3. **medium · ask-user** — Help completeness has no public-route classification or deletion oracle: `internal/conformance/package_shipped_surface_test.go:115` anchors only the inventory title, while `cmd/bench/main_test.go:57` proves rendering from a supplied help literal. Deleting one public command row from that literal can therefore remain green unless the reviewer first decides which registry routes belong in the public inventory.
