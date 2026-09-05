# Move the kit-source predicate into gate

Blocked by: none
Writes: internal/gate/kit_source.go (new), internal/gate/kit_source_test.go (new), internal/adopt/doctor.go, internal/adopt/doctor_rows.go, internal/adopt/link.go, internal/adopt/adopt.go, internal/adopt/init.go, internal/adopt/setup.go, internal/adopt/upgrade.go, internal/worktree/land.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SR20, SR21, SR22, SR23

## What to build

The line is opus / low. Move the kit-source predicate, its symlink-only path
compare, and the kit-directory derivation from `internal/adopt` to
`internal/gate`. They land beside the gate's kit-root derivation, in the new
file `internal/gate/kit_source.go`, per decision (o). The adopt callers and
the landing's default joins name the moved symbols.

`internal/worktree` drops its `internal/adopt` import, and `internal/adopt`
gains an `internal/gate` import. The gate's dependency closure holds no adopt
package, so no cycle forms. The two kit fallbacks stay distinct. The gate's
kit root falls back to the graded root. The adopt kit directory falls back to
the executable's parent, then the current directory. One comment beside them
names that difference.

The compare keeps its symlink-only resolution and calls no `filepath.Abs`,
per decision (d). `bench link`, the doctor rows, and the broker-change notice
answer exactly as today.

The exit proof for this ticket is the pre-existing suite, green with its test
logic unchanged. A mechanical rename is permitted. A needed assertion change
stops the ticket and reports.

## Acceptance

- [ ] `go list -deps ./internal/worktree` names no `internal/adopt`.
- [ ] `bench consumers gate.KitSourceCheckout` lists the landing's default joins.
- [ ] `bench link` refuses in the kit checkout.
- [ ] The doctor rows for the kit checkout read green, and the absent pre-push row routes to the fix.
- [ ] The landing ignores a forged primary executable and seal with the kit-source seam swapped.
- [ ] The moved compare calls no `filepath.Abs`, and `bench test --check canonical-path-owner` stays green.
- [ ] Self-probe: swap the graded-root fallback into the adopt kit directory, and report the observed red.
