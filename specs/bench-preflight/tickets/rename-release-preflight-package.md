# Rename internal/preflight to internal/releasepreflight

Blocked by: none
Ownership fence: `internal/preflight/`, `internal/releasepreflight/`, `internal/preprelease/`, `internal/publication/approved.go`, `internal/conformance/registry/packages.go`, `internal/conformance/ordinary_build_census_test.go`, `internal/conformance/native_workflow_test.go`, `cmd/bench/main.go`, `.bench/BENCH-reference.md`, `.bench/structure-accept`, `tests/canary/package-core-guard/`
Integration surfaces: Go imports→`cmd/bench/main.go`, `internal/conformance/native_workflow_test.go`; release-only package string→`internal/conformance/registry/packages.go`; build-census map key→`internal/conformance/ordinary_build_census_test.go`; structure accept entry→`.bench/structure-accept`; canary probe baseline→`tests/canary/package-core-guard/`; prose references→`internal/preprelease/preprelease.go`, `internal/preprelease/preprelease_test.go`, `internal/publication/approved.go`, `.bench/BENCH-reference.md`; freed package path→implement-preflight-review.md
Contracts: the package import path (old literal `internal/preflight`, new literal `internal/releasepreflight`; domain: exactly one live spelling tree-wide outside the exempt historical set; absence of the old literal is the post-condition) crosses `internal/releasepreflight/`→every surface named in Integration surfaces, asserted per member by the RN2 mutation rows against the real tree
Closure: RN1/verb-dispatches, RN1/usage-line-unchanged, RN1/exported-surface-stable, RN2/main-import, RN2/native-workflow-import, RN2/release-only-string, RN2/census-key, RN2/structure-accept-entry, RN2/canary-baseline-paths, RN2/prose-references

## What to build

Move the release-authorization decision domain from `internal/preflight` to
`internal/releasepreflight` so the `preflight` package path is free for the new
start-oracle domain (verb parity: package name matches its CLI verb). The CLI
verb `release-preflight`, its usage line, and every exported symbol keep their
names — no caller-visible behavior changes. Surfaces: two real Go imports
(`cmd/bench/main.go:33`, `internal/conformance/native_workflow_test.go:216`);
the `ReleaseOnlyPackages` string in `internal/conformance/registry/packages.go:13`;
the census map key in `internal/conformance/ordinary_build_census_test.go:104`;
the accept-list entry at `.bench/structure-accept:14`; the file paths in
`tests/canary/package-core-guard/release-evidence-probe-base.txt`; prose
comments in `internal/preprelease/preprelease.go:103`,
`internal/preprelease/preprelease_test.go:63`, `internal/publication/approved.go:96`;
any `.bench/BENCH-reference.md` mention. Exempt (reviewer-approved 2026-08-11,
recorded in PF24): `specs/`, `CHANGELOG.md`, `capture/`, `ROADMAP.md`,
`decisions/`.

## Acceptance

- [ ] [RN1] (covers PF24) `bench release-preflight` still dispatches with its usage line byte-unchanged, the package's exported identifier list (`go doc github.com/gibbonmi/bench/internal/releasepreflight` versus the pre-move `go doc .../internal/preflight`) is name-for-name identical, and `go test ./internal/releasepreflight ./internal/conformance ./cmd/...` is green with the package at its new path.
- [ ] [RN2] (covers PF24) a repo-wide word-boundary search for the literal `internal/preflight`, excluding only the PF24 exempt set, returns zero hits — every surface in the Integration surfaces line names `internal/releasepreflight`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RN1/verb-dispatches | drop the renamed package's row from `commandRegistry` | the built binary | `go run ./cmd/bench release-preflight`, expect `unknown subcommand` |
| RN1/usage-line-unchanged | reword the usage string during the move | byte comparison of the usage line before and after the move | run `bench release-preflight` with no args pre- and post-move, expect the recorded pre-move line; a mismatch is the red |
| RN1/exported-surface-stable | rename one exported identifier (e.g. `Decide`) during the move, updating its callers | the pre/post `go doc` exported-identifier comparison | capture `go doc` for the package pre-move, apply the mutation, re-capture, expect the name diff as the red |
| RN2/main-import | leave `cmd/bench/main.go`'s import at the old path | the Go compiler | `go build ./...`, expect the unresolved-import failure |
| RN2/native-workflow-import | leave `native_workflow_test.go`'s import at the old path | the Go compiler over test packages | `go vet ./internal/conformance`, expect the unresolved-import failure |
| RN2/release-only-string | leave `ReleaseOnlyPackages` naming the old path | the conformance census over the moved tree | `go test ./internal/conformance -run 'Census\|ReleaseOnly'`, expect the stale-package failure; if it is silently green, add the missing currency assertion inside this ticket's fence so the mutation reds, or return a blocked report naming the enforcement gap — silence is never acceptance |
| RN2/census-key | leave the census map key at the old path | `go test ./internal/conformance -run Census` | expect the unknown-file or stale-key failure; same add-the-assertion-or-block rule as above |
| RN2/structure-accept-entry | leave `.bench/structure-accept:14` naming the old package | `bench structure` over the moved tree | run `bench structure`, expect the stale-accept or reappeared-issue signal; if it tolerates the stale entry silently, same add-the-assertion-or-block rule — `.bench/structure-accept` is in this fence |
| RN2/canary-baseline-paths | leave the probe baseline enumerating old paths | the package-core-guard canary | run the canary family (`bench canary` scoped to package-core-guard), expect the baseline mismatch; same add-the-assertion-or-block rule |
| RN2/prose-references | leave one prose comment naming the old path | the RN2 repo-wide grep run by the verifying coordinator | run the RN2 search, expect the surviving hit as the red |
