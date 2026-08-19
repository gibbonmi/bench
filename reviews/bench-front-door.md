## Standards

Finding count: 1. Worst: auto-fix.

1. **auto-fix** — `.bench/BENCH-reference.md:117` duplicates command visibility with a literal plumbing list that disagrees with `cmd/bench/main.go`'s registry-owned `Inventory` classifications. Remove the enumeration and classify the board-routed `setup` definition public.

## Spec

Finding count: 3. Worst: auto-fix.

1. **auto-fix** — Story 29 and the adapter decision (`specs/bench-front-door/spec.md:92,155-157`) require Codex `$bench-*` routes to load their named phase, but `.agents/commands/bench.md:11-14` recognizes only `/bench-*`.
2. **auto-fix** — R11 (`specs/bench-front-door/spec.md:53,209`) requires malformed `status --route` forms to print the exact grammar, but `internal/status/status.go:397-412` returns generic usage text.
3. **auto-fix** — Story 34/R39 and the scope boundary (`specs/bench-front-door/spec.md:103,237,288`) still describe a literal non-registry inventory, contrary to the reviewer-approved registry projection in `cmd/bench/command_registry.go`.

## Coverage

Finding count: 4. Worst: high.

1. **auto-fix** — R16/R21 lead routes with control-bearing staged-spec or ready-map paths reach `toon.Table` unsanitized (`internal/status/status.go:424-427`); existing tests cover only runner-ups. Add lead-path hostile cases and sanitize or refuse safely.
2. **auto-fix** — R39 still lacks an inventory-completeness deletion oracle for every public row: `bench setup` is board-routed but internal (`cmd/bench/main.go:118`), and deleting `worktree exec` leaves current checks green. Classify setup and add the independently authored full inventory/deletion proof.
3. **auto-fix** — Story 7/R9 lacks a command-boundary Codex handoff test: `internal/handoff/render_test.go` covers only the Claude fallback.
4. **auto-fix** — Story 19/R25 lacks an oracle for the generic `/bench-*` invocability arm; `otherPhaseAction` is skipped by the round-trip test at `internal/status/route_test.go:38-43`.
