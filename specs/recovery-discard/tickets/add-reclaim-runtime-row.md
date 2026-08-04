# Drive reclaim end to end through the runtime fixture

Blocked by: add-spec-build-reclaim.md
Ownership fence: `internal/contract/runtime/runtime_spec_build_test.go`
Contracts: the reclaim CLI grammar and its TOON receipt cross `internal/spec/build.go` and `cmd/bench/specbuild.go`→`internal/contract/runtime/runtime_spec_build_test.go` and are asserted by RR1 and RR3 against the real built `bench` binary rather than a service call
Assumptions: `add-spec-build-reclaim.md` has landed so the verb and its grammar row and its renderer all exist; the fixture drives the real subject binary through the package's existing spec-build helpers; the promotion-reclaims-refs row already in this file stays as the forward-path positive control; every claim is re-derived from the tree at pickup

## What to build

Without a runtime row the reclaim CLI shape is a delegate's guess: the package
rows in `internal/specbuild` grade the service, and none of them can see the
grammar, the dispatch arm, or the receipt the operator actually reads.

Drive a full spec-build run through the fixture to a terminal record, arrange for
its provisional refs to still exist, then run `bench spec build reclaim <slug>`
and assert the plan output, and `--apply <fingerprint>` and assert the refs are
gone. Assert the drifted-fingerprint refusal through the CLI too, because that is
the boundary where a maintainer supplies a stale value.

## Acceptance

- [ ] [RR1] `bench spec build reclaim <slug>` prints a plan with a fingerprint and mutates no ref.
- [ ] [RR2] `bench spec build reclaim <slug> --apply <fingerprint>` with the exact planned fingerprint exits zero and the run's reclaimable refs no longer resolve.
- [ ] [RR3] `--apply` with a stale fingerprint exits non-zero, names the fresh-plan action, and leaves every ref in place.
- [ ] [RR4] `reclaim` with no slug exits 2 with the usage line.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RR1 | delete the reclaimable refs during the plan dispatch arm | the reclaim runtime row | apply the plan in the plan-only branch of `cmd/bench/specbuild.go`, rebuild the subject binary, run `go test ./internal/contract/runtime -run SpecBuild -timeout 300s`, expect the unchanged-refs assertion to fail |
| RR2 | drop the `--apply` case from the reclaim dispatch so it always plans | the reclaim runtime row | remove the flag branch, rebuild the subject binary, run `go test ./internal/contract/runtime -run SpecBuild -timeout 300s`, expect the refs-gone assertion to fail |
| RR3 | exit zero on a stale fingerprint | the reclaim runtime row | swallow the refusal error in the dispatch arm, rebuild the subject binary, run `go test ./internal/contract/runtime -run SpecBuild -timeout 300s`, expect the non-zero-exit assertion to fail |
| RR4 | set the reclaim grammar's minimum positional count to zero | the reclaim runtime row | relax the grammar, rebuild the subject binary, run `go test ./internal/contract/runtime -run SpecBuild -timeout 300s`, expect the exit-2 assertion to fail |
