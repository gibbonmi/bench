# Add the guard-rg verb and its fail-open shim

Blocked by: 02-classify-rg-calls.md
Writes: cmd/bench/main.go, cmd/bench/guards.go, cmd/bench/guard_rg_test.go (new), .bench/hooks/block-unsafe-rg.sh (new), internal/otelrecord/registry.go, internal/conformance/subcommand_routing_table_test.go, internal/conformance/entry_point_parity_test.go, internal/conformance/entry_point_parity_rows_test.go (new), internal/conformance/package_core_checks_test.go, internal/conformance/guard_header_manifests_test.go (new), internal/systemtest/rg_guard_hook_test.go (new), internal/packagesurface/assets.go, internal/packagesurface/assets_test.go
Covers: R41, R42, R43, R44, R45, R46, R50, R56, R57

## What to build

A PreToolUse hook can call `bench guard-rg`, and the verb refuses each Bash command that `rgguard.Classify` blocks. This ticket adds the verb, the shim, and their registrations. `04-wire-rg-guard-hook.md` wires the shim into the harnesses.

Add the registry row `guard-rg` beside the guard rows at `cmd/bench/main.go:123-126`, with `Hook: true`, `internalInventory`, and `axiExempt(axiReasonPlumbing)`.
Add its handler in `cmd/bench/guards.go`, after the shape of `guardBenchFollowOn` (`cmd/bench/guards.go:152-189`). It reads the command through `rgguard.CommandFromEnvelope` from `02-classify-rg-calls.md` and decodes no envelope itself. It recovers a panic to exit 3. On an unreadable envelope it warns on stderr and exits 0. On a blocked verdict it prints `Verdict.Message()` on stderr and exits 2. Otherwise it exits 0.

Add the shim `.bench/hooks/block-unsafe-rg.sh` with mode 755 and the four manifest keys `name`, `boundary`, `denies`, and `why` (`internal/guards/guards.go:94`). Its `# denies:` line states the rule: an `rg` call with no path operand and no stdin supplier, and an `rg` call with a live command substitution. That line is the shipped advertisement of the rule.
The shim copies the fail-open posture of `.bench/hooks/block-bench-follow-on.sh` (`:20`, `:24`, `:42-57`): a missing shared library, a missing core, a core that does not know `guard-rg`, and any core exit other than 0 or 2 each warn on stderr and exit 0. The shim must lint clean under the gate's shellcheck phase at `-S warning`.

Register the verb:

- a `hook.guard-rg` seam in `internal/otelrecord/registry.go:36-42`;
- `"guard-rg": exempt(whyPlumbing)` in `internal/conformance/subcommand_routing_table_test.go:54-67`;
- an `entryPointParity` row for the shim with `command` `guard-rg`, a benign allowed Bash input such as `rg foo .`, and `compare: true`;
- the shim in the `checkGuardHeaderManifests` list (`internal/conformance/package_core_checks_test.go:312-323`).

Optionally add the shim to `RequiredPackAssets` (`internal/packagesurface/assets.go:13-36`) and its test. The consumer payload ships `.bench/hooks` as a tree (`.bench/consumer-payload.json:18`), so the payload takes no edit.
Do not edit `bin/bench.sh`. It has 538 lines over its 400-line budget, and `guard-rg` already reaches the binary through the default route at `bin/bench.sh:537`.

`internal/conformance/entry_point_parity_test.go` (444 lines) and `internal/conformance/package_core_checks_test.go` (451 lines) are over budget with no grant, and the gate lane runs `structure --growth` (`internal/gate/lane.go:94`). A file that grows past its base count reds that lane (`internal/structure/structure.go:309-331`). In each file, move at least as many lines as you add into the named new sibling file (`entry_point_parity_rows_test.go`, `guard_header_manifests_test.go`). Move a cohesive block such as the row table or the manifest check, and change no behavior.

Test the verb in `cmd/bench/guard_rg_test.go` (`TestGuardRgVerb`). Test the shim through a real envelope pipe in `internal/systemtest/rg_guard_hook_test.go` (`TestRgGuardHookFailsOpen`, `TestRgGuardHookRefuses`), after `internal/systemtest/bench_follow_on_test.go` and its `staleCoreAnswer` (`:40`).
Each hook child gets `BENCH_KIT=` the owner kit and `BENCH_RUN_BINARY=` the selected binary, as `bench_follow_on_test.go:60` does. BENCH_KIT is supplied by `bench test --check system`; never run the tagged suite directly.

## Acceptance

- [ ] `checkGuardHeaderManifests` reds when the shim lacks a manifest key (R41).
- [ ] The shim with no resolvable core exits 0 and prints a warning (R42).
- [ ] The shim over a core without `guard-rg` exits 0 (R43).
- [ ] `bench guard-rg` with an unreadable envelope exits 0 and warns (R44).
- [ ] The shim with an `rg foo` envelope exits 2 and prints the pathless message on stderr (R45).
- [ ] `TestHookSeamsMatchTheHookDispatchRows` passes with the `hook.guard-rg` row and reds without it (R46).
- [ ] `bench test --check entry-point-parity` passes with the shim row and reds without it (R56).
- [ ] `bench test --check subcommand-routing` passes with the `guard-rg` row and reds without it (R57).
- [ ] The link plan copies `.bench/hooks/block-unsafe-rg.sh` through the hooks tree row (R50, review-owned).
- [ ] `bench structure --growth` reports no grown file.
