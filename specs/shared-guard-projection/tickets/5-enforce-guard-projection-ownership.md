# Enforce guard projection ownership

Blocked by: 1-share-git-command-projections.md, 2-add-envelope-owner-for-git.md, 3-migrate-bench-and-write-readers.md, 4-migrate-bench-projection-consumers.md
Writes: internal/conformance/guard_projection_owner_test.go (new), internal/conformance/checks_test.go, internal/conformance/tier_live_tree_test.go, internal/conformance/guard_scan_owner_test.go, internal/conformance/guard_owner_walk_test.go (new)
Covers: GP30, GP31, GP32, GP33, GP42, GP67

## What to build

Make the registered conformance check refuse any second owner of the guard grammar.
Tickets 1–4 supply the shared projection and envelope owner and retire every private copy in gitguard, benchguard, and writeguard, so the live tree is clean when this check starts.

Add `checkGuardProjectionOwner` beneath the registered package-core-guard check.
It grades the production Go AST in shellcommand, gitguard, benchguard (`benchguard.go` and `reviewer.go`), and writeguard through real import bindings, unused function bodies included.
Outside shellcommand, it refuses hook-envelope JSON decoding, the shell-name predicate, the wrapper `-c` recognizer, and adjacent worktree-exec child extraction.
It ignores comments and fixture strings. A consumer root without kit guard source passes; a kit root with a missing or unreadable owner gets a named diagnostic.
The subcheck engages only on a root that holds `internal/conformance` and at least one governed guard package. The `guard-classifier-table/word-test-drops-xargs` canary tree holds no guard package, so it must get no new diagnostic.
Keep the registry name, tier, subject, and input binding. Classify every test in `guard_projection_owner_test.go` that reads live sources in `tier_live_tree_test.go`.
Keep the resolver, classifier-table, and entry-point-parity checks on their registered execution paths.
The check must not flag Git's own `worktree` subcommand handling (`gitguard/verdict.go:70`, `gitguard.go:83`); it refuses only the adjacent `worktree exec … --` child extraction.
To stay within about 200k context, read only the migrated call sites in the guard packages, and read the four retired functions through `git show <build-entry commit>:<path>` for the mutation tests.

The guard-scan-loop outcome moved `checkPackageCoreAndGuards` to `checks_test.go:83-98`. Extend that sole declaration after `checkGuardScanLoopOwner` (`:97`), and keep that subcheck.
Do not edit package_core_checks_test.go.
Keep checks_test.go within its 739-line grant (660 lines at 1228eedff), and do not add the banned skills-index literals (`skills_index_checks_test.go:375-387`).
Do not call `Skip`, `Skipf`, or `SkipNow` in the new test file (`internal/conformance/skip_ownership_test.go:30-64`). The new file stays within 400 lines. Add no grant.

## Acceptance

- [ ] A renamed unused copy of each retired recognizer or envelope decoder in gitguard, `benchguard.go`, or writeguard makes package-core-guard red (GP30).
- [ ] A renamed private wrapper-flag or exec-child copy in `reviewer.go` makes package-core-guard red (GP67).
- [ ] A consumer tree passes, and a kit tree with a missing owner gets the named diagnostic (GP31, GP32).
- [ ] The live-tree owner test runs through its registered binding and its classification (GP33).
- [ ] The resolver, classifier-table, and parity checks keep their registered execution paths (GP42).
- [ ] Every retained canary fixture keeps its expected diagnostics, with no new projection-owner diagnostic.
- [ ] Omitting only the registered `checkGuardProjectionOwner` call makes `TestGuardProjectionOwnerBitesOnPrivateCopies` red, and restoration returns the same green.
