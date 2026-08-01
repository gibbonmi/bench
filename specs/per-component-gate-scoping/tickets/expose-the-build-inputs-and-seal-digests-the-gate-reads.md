# Expose the build inputs and seal digests the gate reads

Blocked by: none
Ownership fence: `internal/freshness/freshness.go`, `internal/freshness/freshness_test.go`
Assumptions: `buildInputs`, `parseSeal`, and `secureContents` are unexported in
`internal/freshness` at base `cce51af`; `Digest` hashes exactly `buildInputs`'
paths. Re-derive both facts from the tree at pickup rather than from this line.

## What to build

The gate's build-component declaration and its build attestation both need facts
that today only live inside `internal/freshness`: which files the binary's build
closure covers, and what a published seal claims. Export them so neither
consumer re-derives a second, drifting copy — the build declaration takes its
input set from `go list -deps ./cmd/bench` through this accessor, and the
attestation compares against the seal this accessor parses.

`BuildInputs(root)` returns the repository-relative, slash-separated, sorted
paths `Digest` hashes. `SealDigests(executable)` returns the seal's source and
executable digests through the same strict read `Verify` performs, so a
symlinked sidecar, an irregular file, or malformed contents refuse rather than
resolve.

## Acceptance

- [ ] PS1 — `BuildInputs` returns exactly the paths `Digest` hashes for a module, repository-relative and sorted, including `go.mod`, `go.sum`, and every auxiliary manifest entry.
- [ ] PS2 — `SealDigests` returns a published seal's two digests, and refuses a symlinked sidecar, an irregular sidecar, and malformed seal contents.
- [ ] PS3 — a `BuildInputs` failure to resolve the closure is returned as an error, never as an empty set.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS1 | drop the `go.mod`/auxiliary-manifest additions from `BuildInputs`' return, keeping them in `Digest` | `TestBuildInputsCoverEveryDigestedPath` | publish a fixture module, call `BuildInputs`, rehash its paths with the framing `Digest` uses, compare against `Digest(root)` |
| PS2 | replace `SealDigests`' `secureContents` read with `os.ReadFile` | `TestSealDigestsRefuseAnUntrustedSidecar` | publish a binary, replace `dist/bench.seal` with a symlink to a valid seal, call `SealDigests`, expect a refusal |
| PS3 | return `nil, nil` when `go list` exits nonzero | `TestBuildInputsErrorOnUnresolvableClosure` | point `BuildInputs` at a root with no `go.mod`, expect a non-nil error and no paths |
