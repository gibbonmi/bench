# Enforce guard projection ownership

Blocked by: 1-share-git-command-projections.md, 2-migrate-guard-envelope-consumers.md, 3-migrate-bench-projection-consumers.md
Writes: internal/conformance/guard_projection_owner_test.go (new), internal/conformance/package_core_checks_test.go, internal/conformance/checks_test.go, internal/conformance/tier_live_tree_test.go, internal/systemtest/guard_projection_test.go (new), cmd/bench/guard_projection_test.go (new)
Covers: GP30, GP31, GP32, GP33, GP34, GP35, GP36, GP37, GP38, GP39, GP40, GP41, GP42, GP44, GP45, GP46, GP67, GP71

## What to build

Enforce one owner for the guard grammar, and prove through the actual hook and core entries that the migration changed no supported behavior.
Add `checkGuardProjectionOwner` beneath the registered package-core-guard check.
It grades the production Go AST in shellcommand, gitguard, benchguard (`benchguard.go` and `reviewer.go`), and writeguard through real import bindings, unused function bodies included.
Outside shellcommand, it refuses hook-envelope JSON decoding, the shell-name predicate, the wrapper `-c` recognizer, and adjacent worktree-exec child extraction.
It ignores comments and fixture strings. A consumer root without kit guard source passes; a kit root with a missing or unreadable owner gets a named diagnostic.
The subcheck engages only on a root that holds `internal/conformance` and at least one governed guard package. The `guard-classifier-table/word-test-drops-xargs` canary tree holds no guard package, so it must get no new diagnostic.
Keep the registry name, tier, subject, and input binding. Classify every test in `guard_projection_owner_test.go` that reads live sources in `tier_live_tree_test.go`.

If `checkPackageCoreAndGuards` is still in package_core_checks_test.go, move only that 16-line declaration to checks_test.go and add the subcheck there.
If the guard-scan-loop outcome moved it first, extend that sole declaration and keep its subcheck.
Keep the npm producers and their formatter call in package_core_checks_test.go. package_core_checks_test.go must not grow.
Keep checks_test.go within its 739-line grant after either outcome order, and do not add the banned skills-index literals (`skills_index_checks_test.go:375-387`).

Run every actual hook and selected-core witness through the existing system process owner under `BENCH_KIT`, as the spec's Testing decisions state.
GP46 compares the pre-migration reference executable, captured before ticket 1, with the candidate over family D and the full decoder matrix, without elapsed-time equality.
This final regression repeats earlier cases; it does not supply a missing earlier test or repair.
Keep the resolver, classifier-table, entry-point-parity, canary, and census execution paths.
Do not call `Skip`, `Skipf`, or `SkipNow` in a new test outside the skip-owner directory (`cmd/bench/skip_ownership_test.go:30-64`).
Give the new system test file the `//go:build system` tag, and write no literal wait deadline in a new test (`wait_deadline_literal_test.go:33-39`).
Each new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] A renamed unused copy of each retired recognizer or envelope decoder in gitguard, `benchguard.go`, or writeguard makes package-core-guard red (GP30).
- [ ] A renamed private wrapper-flag or exec-child copy in `reviewer.go` makes package-core-guard red (GP67).
- [ ] A consumer tree passes, and a kit tree with a missing owner gets the named diagnostic (GP31, GP32).
- [ ] Every retained canary fixture keeps its expected diagnostics, with no new projection-owner diagnostic.
- [ ] The live-tree owner test runs through its registered binding and its classification (GP33).
- [ ] The landed degraded-guard witnesses for family A, parse advice, missing library, missing core, core error, and recovery stay unchanged and pass (GP34–GP37, GP39, GP40).
- [ ] The file-write hook keeps its warning-allow tuple without its library or core, and the follow-on hook keeps its warning-allow tuple without its resolver, wrapper, or core (GP38, GP71).
- [ ] Stale-core Bench and non-Bench calls keep their different refusal tuples (GP41).
- [ ] The resolver, classifier-table, and parity checks keep their registered execution paths (GP42).
- [ ] A supported Bench invocation makes no raw census record (GP44).
- [ ] An expansion-shaped command leaves the sentinel absent through the actual entry (GP45).
- [ ] Reference and candidate executables give identical tuples over family D and the decoder matrix (GP46).
- [ ] Omitting only the registered `checkGuardProjectionOwner` call makes `TestGuardProjectionOwnerBitesOnPrivateCopies` red, and restoration returns the same green.
