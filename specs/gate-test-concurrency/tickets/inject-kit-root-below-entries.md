# Inject the kit root below the gate entries

Blocked by: none
Ownership fence: `internal/gate`
Integration surfaces: exported entry signatures (`RunCommand`, `Execute`, `ExecuteReusingFreshGreen`, `ExecuteTree`, `PhasesCommand`, `ComposedGreen`) unchanged→existing external callers (`cmd/bench`, `internal/status`) exercised by IK2; kit-taking execution boundary→retire-fixture-kit-pins.md; kit-root carrier→`internal/gate`
Contracts: none crosses — the carrier, the kit-taking boundary, and every consumer live inside the one fence, and external callers keep today's entry signatures
Closure: IK1/phase-table, IK1/component-identity, IK1/component-scoping, IK2/env-set, IK2/empty-fallback, IK2/kit-ne-root, IK3/subject-strip

## What to build

Ambient `BENCH_KIT` is resolved exactly once, at the exported gate entries,
under today's rule (set and non-empty wins, else the graded root). The entries
are the gate run command, the execute and reuse-execute paths, the
tree-execute path, the phases plumbing command, and the composed-green query
`internal/status` consults — all with unchanged signatures. Below the entries
the kit root travels as an explicit value — parameter or evaluation-scoped
field, implementer's choice — and all three current consumers move onto it:
phase-table resolution, component-identity resolution
(`resolveComponentIdentities`), and component scoping
(`scopeComponentsForIdentityGenerations`, reached from evaluation and from the
composed-green query alike). No production code in `internal/gate` reads the
variable after this ticket; `kitRoot`'s single-derivation comment discipline
transfers to the entry-time resolution.

This ticket also exposes the package-internal execution boundary that takes
the resolved kit explicitly (the seam beneath the exported entries that the
execute paths already share). It is what the dependent ticket migrates
entry-driving fixture tests onto; here it only needs to exist and carry the
injected kit, proven by the guard test below. Existing tests keep their pins
and their passing state — no test-side migration lands in this ticket.

A new serial ambient-guard test (it pins env deliberately, so it never adopts
`t.Parallel`) exports a foreign `BENCH_KIT`, constructs a fixture with the kit
injected, and asserts all three consumers resolve from the fixture through the
kit-taking boundary: the phase table carries the fixture's phases, component
identities resolve, and scoping reports the fixture eligible with the fixture
as runner root. A consumer left on an ambient-fallback read is correct
whenever kit equals root — every fixture — so only this hostile-ambient
assertion distinguishes it; that is why the guard names each consumer
separately. A second new serial assertion pins `BENCH_KIT=""` at an exported
entry and asserts the graded-root fallback — no test covers set-but-empty
today.

## Acceptance

- [ ] [IK1] (covers KC1) with ambient `BENCH_KIT` naming a foreign path and no per-test pin, evaluation through the kit-taking boundary over an injected fixture root resolves the fixture's own phase table, component identities, and component scoping.
- [ ] [IK2] (covers KC2) a wrapper-routed entry with `BENCH_KIT` set resolves the kit from the environment at entry, including kit ≠ root; a new serial assertion proves set-but-empty falls back to the graded root.
- [ ] [IK3] (covers KC4) the closed subject environment still strips `BENCH_KIT`/`BENCH_WRAPPER` from the gate script's environment.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| IK1/phase-table | make the phase-table consumer prefer ambient `BENCH_KIT` over the injected value | the serial ambient-guard test | apply, run `go test -count=1 -run <guard> ./internal/gate`, expect the table assertion red |
| IK1/component-identity | make the identity consumer re-read ambient `BENCH_KIT` | the serial ambient-guard test | apply, run the guard, expect the identity assertion red |
| IK1/component-scoping | make the scoping consumer re-read ambient `BENCH_KIT` | the serial ambient-guard test | apply, run the guard, expect the eligibility/runner-root assertion red |
| IK2/env-set | make entry resolution ignore the environment and always use the graded root | the existing entry tests that pin `BENCH_KIT` | apply, run `go test -count=1 -run <entry tests> ./internal/gate`, expect red |
| IK2/empty-fallback | make entry resolution accept an empty `BENCH_KIT` as the kit | the new set-but-empty entry assertion | apply, run its focused test, expect the fallback assertion red |
| IK2/kit-ne-root | make entry resolution return the graded root when the environment names another directory | the kit≠root manifest-resolution test | apply, run its focused test, expect red |
| IK3/subject-strip | stop stripping `BENCH_KIT` from the subject environment | the existing subject-env stripping test | apply, run its focused test, expect red |
