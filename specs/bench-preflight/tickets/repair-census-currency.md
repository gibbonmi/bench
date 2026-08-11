# Repair: census carries the new preflight domain and its keys resolve

Blocked by: none
Ownership fence: `internal/conformance/ordinary_build_census_test.go`
Integration surfaces: census map→`internal/conformance/ordinary_build_census_test.go`; new-domain tests→existing `internal/preflight/decision_test.go` (read, not written)
Contracts: none crosses — the census file both declares and checks the entries; RC1 asserts the resolution rule inside that one owner
Closure: RC1/new-domain-entry, RC1/keys-resolve

## What to build

Review finding S3 (+S5): the rename moved the census key
`internal/preflight/decision_test.go` → `internal/releasepreflight/...`
without re-adding the new decision domain's entry, and reverting a census key
to a nonexistent path leaves the census silently green — the currency
assertion the rename ticket's add-the-assertion-or-block rule demanded never
landed for map keys. Add the new domain's `internal/preflight/decision_test.go`
entry, and extend the census with the same currency rule the
`ReleaseOnlyPackages` check got: every census map key must resolve to an
existing file in the tree, red naming the stale key otherwise. Record the
break-it command for both this and the `ReleaseOnlyPackages` check in the
landing commit body (S5's missing bite record).

## Acceptance

- [ ] [RC1] (covers local) the census map names `internal/preflight/decision_test.go`, and a census key naming a file absent from the tree turns `go test ./internal/conformance -run Census` red naming that key.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RC1/new-domain-entry | remove the new `internal/preflight/decision_test.go` entry | the census key-resolution assertion paired with the census's own file enumeration | `go test ./internal/conformance -run Census`, expect the unlisted-file or stale-key failure |
| RC1/keys-resolve | point one census key at `internal/preflight/decision_test_gone.go` | the new key-resolution assertion | same run, expect the stale-key failure naming the path |
