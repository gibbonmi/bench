# Refuse a ticket whose red-mutation probes contradict its fence

Blocked by: none
Ownership fence: `internal/specbuild`, `.agents/skills/bench-craft-tickets/SKILL.md`
Contracts: the extracted probe packages and the agreement verdict cross `ParseTicket`→the assign refusal inside `internal/specbuild`, asserted by FA1 and FA2 against real ticket markdown parsed by the real function; the exported validation crosses `internal/specbuild`→the conformance sweep of `sweep-staged-ticket-agreement.md`, asserted there against this ticket's real export

## What to build

For every `go test` package argument in a red-mutation row's operation-sequence
cell, the row must have somewhere its red can live: a fence entry that is a
segment-aware package prefix of the probed package, a fenced `_test.go` file
whose directory is exactly the probed package, or an owner cell naming a
backticked `_test.go` path that exists in the tree with its directory exactly
the probed package. A ticket failing the rule is refused at
`bench spec build assign`, before any delegate is charged.

Extraction reads only the operation-sequence cell — a `go test` in the mutation
column or in prose is not a probe — and takes every package argument of every
invocation in the cell, skipping flags. A `./...` argument is refused outright:
no fence can own it. A `_test.go` in a subdirectory of the probed package is a
different Go package and does not count. An owner-cell path that does not exist
falls back to the fence rule rather than crashing.

The validation is one exported function beside `ParseTicket`, because the
conformance sweep in the next ticket must consume the same derivation. This
ticket also documents the owner-cell convention in `craft-tickets`' red-mutations
field bullet — the pre-existing-owner escape for regression rows that author
nothing — so the enforcement and its advertisement land in one change.

## Acceptance

- [ ] [FA1] assign refuses a ticket whose operation-sequence cell probes a package with no fence prefix, no fenced same-directory `_test.go`, and no existing owner-cell test path, naming the row and the package.
- [ ] [FA2] a ticket whose fence names the probed package's own `_test.go` file assigns, and a regression row whose owner cell names an existing `_test.go` in the probed package assigns without fence coverage.
- [ ] [FA3] an owner cell naming a `_test.go` path that does not exist falls back to the fence rule, so a dangling path cannot license a probe.
- [ ] [FA4] a `./...` package argument is refused.
- [ ] [FA5] containment is segment-aware in both directions and a subdirectory `_test.go` does not satisfy the rule.
- [ ] [FA6] `go test -count=1 ./a ./b` extracts both packages and skips the flag, and a `go test` outside the operation-sequence cell is not extracted.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FA1 | drop the refusal so every parsed ticket assigns | the probe-refusal test | remove the agreement call from the assign path, run `go test ./internal/specbuild -run Agreement -timeout 180s`, expect a nonzero matched-test count and the refusal assertion to fail |
| FA2 | require a package-prefix fence entry for every probe | the legitimate-shapes regression test | drop the fenced-test and owner-cell arms, run `go test ./internal/specbuild -run Agreement -timeout 180s`, expect a nonzero matched-test count and the stays-assignable assertion to fail |
| FA3 | accept any backticked owner path without checking existence | the dangling-owner test | skip the existence stat, run `go test ./internal/specbuild -run Agreement -timeout 180s`, expect a nonzero matched-test count and the fallback-refusal assertion to fail |
| FA4 | treat `./...` as an ordinary package argument | the wildcard-refusal test | drop the wildcard branch, run `go test ./internal/specbuild -run Agreement -timeout 180s`, expect a nonzero matched-test count and the wildcard assertion to fail |
| FA5 | compare paths with a raw string prefix | the segment-boundary test | replace the segment comparison with `strings.HasPrefix`, run `go test ./internal/specbuild -run Agreement -timeout 180s`, expect a nonzero matched-test count and the `internal/spec`-vs-`internal/specbuild` assertion to fail |
| FA6 | extract only the first package argument of the first invocation | the extraction-scope test | return after the first token, run `go test ./internal/specbuild -run Agreement -timeout 180s`, expect a nonzero matched-test count and the multi-package assertion to fail |
