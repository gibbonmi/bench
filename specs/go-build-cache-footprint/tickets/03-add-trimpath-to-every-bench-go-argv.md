# Add -trimpath to every Bench-owned Go argv

Blocked by: 01-resolve-test-roots-without-runtime-caller.md

Writes: internal/gate/phases.go, internal/gate/gate_go.go, internal/gate/lane.go, internal/testreport/testreport.go, internal/preprelease/preprelease.go, internal/conformance/ordinary_build_census_test.go, internal/gate/lane_test.go, internal/gate/branch_native_phases_test.go, internal/preprelease/preprelease_test.go, internal/testreport/testreport_test.go, CHANGELOG.md

## What to build

Without `-trimpath` Go hashes the package directory into every compile action
ID, so each checkout writes a complete new set of archives. Put `-trimpath` on
every Bench-owned `go test`, `go vet`, and `go build` argv, and keep `-count=1`
on every test argv.

One exported producer in the gate package owns the base test argv. It returns
`go test -trimpath -count=1` before its caller's arguments. Six forms compose
it: the `test`, `race`, and `system` phases, the focused `bench test` argv, the
release `coreTestStep` argv, and the ship conformance step. The release
package's own copy at `internal/preprelease/preprelease.go:131` goes away, and
the ship conformance step composes the exported gate producer instead.

One shared flag owner in the gate package supplies `-trimpath` to the non-test
argvs. The `vet` phase in `internal/gate/phases.go` and the kit lane's `vet`
and `build` entries in `internal/gate/lane.go` take the flag from that owner.
That owner stays separate from the release build flags in
`internal/releaseevidence/requirements.json`, because those flags state the
shipped binary's evidence contract, not gate policy.

The census literals in `internal/conformance/ordinary_build_census_test.go:36-48`
are the seam that moves for T01. They pin the `test` argv, the `race` prefix
`go test -race -count=1 -v`, and the `system` argv as literal values.
`internal/gate/branch_native_phases_test.go:108` compares the race argv against
`raceDriverArgv()` itself, so that comparison passes whatever the producer
returns and catches no missing flag.

No test asserts the `coreTestStep` argv or the `bench test` argv today, and
`internal/preprelease/preprelease_test.go:74` asserts `gate.GateGoArgv` rather
than a test argv. T04, T05, and T07 therefore add new literal expectations
rather than move existing ones.

A linked repository's declared phase or lane argv keeps its own flags, because
the project owns its manifest.

## Acceptance

- [ ] T01 — The `test`, `race`, and `system` phase argvs each begin with `go test -trimpath -count=1`.
- [ ] T02 — The `vet` phase argv is `go -C <root> vet -trimpath ./...`.
- [ ] T03 — The kit lane's `vet` and `build` argvs each carry `-trimpath`.
- [ ] T04 — The `bench test` argv carries `-trimpath` and `-count=1`.
- [ ] T05 — The release `coreTestStep` argv carries `-trimpath` and `-count=1`.
- [ ] T07 — The ship conformance step argv begins with `go -C <kit> test -trimpath -count=1`.

Delivered outcome: two checkouts with identical content share one action ID, so
a second worktree stops writing a full archive set.
