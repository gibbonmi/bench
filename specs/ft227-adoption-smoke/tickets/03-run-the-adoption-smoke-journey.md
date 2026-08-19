# Run the adoption smoke journey under the system phase

Blocked by: 01-seed-the-gate-input-manifest.md, 02-guard-the-scaffolded-inventory-call.md
Writes: internal/systemtest/adoption_test.go, projects/benchkit.md

## What to build

One tagged (`//go:build system`) journey in a new file of `internal/systemtest`
that adopts `owner.repos[1]` with `bench setup --yes` on the selected executable
(with `BENCH_KIT` at the real kit) and then runs the scaffolded gate through the
installed `.bench/bin/bench.sh` with `BENCH_RUN_BINARY` at the selected
executable. Every launch — setup and each gate leg — binds one private
`BENCH_HOME` under the test's temporary directory. `runSelected` and
`runWrapper` carry fixed environments with no `BENCH_HOME` override, so each
launch goes through `owner.observeSelected()` plus `owner.runAt(...)` with its
own environment, the pattern `TestWorktreeReauthorizeJourney` already uses.
The stub is not a phase-table gate, so `gate-run` injects no `BENCH_RUN_BINARY`
into it: inside the gate the environment is `PATH` plus the two declared names,
and the stub's wrapper resolves the `.bench/dist/bench` copy setup installed.

A green leg asserts the exact `gate: green` line, never the substring — a
reused verdict prints `gate: green (fresh verdict reused for this tree)`.

Legs, in order; the first gate run is plain `gate`, every later one `gate --fresh`:

1. setup exits 3, prints the red sentinel doctor row, and leaves `.bench/gate.sh`,
   `.bench/gate-inputs.json`, `.bench/bin/bench.sh`, `.bench/dist/bench` on disk.
2. the untouched stub: gate exits 1, stderr names the sentinel remedy.
3. delete the one line carrying the exported sentinel marker: gate exits 0,
   stdout has the exact line `gate: green`, no `tests/canary` exists.
4. create `tests/canary/<family>/<fixture>/` with one file: gate exits 0, stdout
   has `canary inventory ok (1 fixture bindings)` and the exact line `gate: green`.
5. remove `.bench/gate-inputs.json`: gate exits 1, stderr contains
   `HOME: unbound variable`.
6. restore the manifest, empty `tests/canary` to a bare directory: gate exits 1,
   stderr contains `canary fixture inventory is empty`.
7. the private `BENCH_HOME` directory is still empty.

The journey records the selected executable through the owner's observation
helper on every launch it makes, so the identity ledger and the three-repository
budget hold. `projects/benchkit.md`'s system-package paragraph names the adoption
journey beside the stripped-distribution one.

## Acceptance

- [ ] setup in `owner.repos[1]` exits 3 with the red sentinel row and installs the four paths (SM1).
- [ ] the untouched stub is red naming the sentinel remedy (SM2).
- [ ] with the sentinel line gone and no `tests/canary`, the gate is green through the installed wrapper (SM3).
- [ ] with one project fixture, the gate is green and reports one binding (SM4).
- [ ] with the manifest removed, the gate is red naming `HOME` (SM5).
- [ ] with an empty `tests/canary`, the gate is red naming the empty inventory (SM6).
- [ ] every launch binds the private `BENCH_HOME`, and it stays empty (SM7).
- [ ] `projects/benchkit.md` names the adoption journey.
