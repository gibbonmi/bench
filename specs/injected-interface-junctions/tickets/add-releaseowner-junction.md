# Add the ReleaseOwner junction tests

Blocked by: none
Ownership fence: `internal/specbuild/release_junction_test.go`
Contracts: ReleaseEvidence and the release verdict cross `internal/specbuild`→`internal/worktree` through the real `worktree.ReleaseProvisional`, asserted by RJ1 and RJ2 against the real producer
Assumptions: production is presumed correct — a red against unmodified production is a stop-and-surface finding, never forced green; the existing `realOwner` fake with its bare-nil `Release` stays for the fast path; claims re-derived from the tree at pickup

## What to build

Service-level release coverage where the real `worktree.ReleaseProvisional`
sits behind the specbuild seam: one integrate/release path that succeeds on
durable evidence, and one that surfaces the producer's real refusal.

## Acceptance

- [ ] [RJ1] a Service-driven release with durable checkpoint/integration evidence completes against the real `worktree.ReleaseProvisional` (a test owner delegating to it, mirroring `realOwner.Create`'s pattern).
- [ ] [RJ2] a Service-driven release whose request does not match the assignment surfaces the producer's exact refusal (`provisional release request, assignment, or path mismatch; checkout retained`) through the consuming surface.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RJ1 | drop the CheckpointRef from the evidence the test owner forwards | the junction test | apply the mutation, run `go test ./internal/specbuild -run <RJ1 test>`, expect the release refusal to fail the test |
| RJ2 | swap the mismatch fixture's request for the assignment's real request | the junction test | apply the mutation, run `go test ./internal/specbuild -run <RJ2 test>`, expect the asserted refusal string to go missing |
