# 07 Share skip scans and prove reuse through real dispatch

Blocked by: 04-share-source-observations-for-git-policy.md, 05-share-the-architecture-source-scan.md, 06-share-the-bounds-source-scan.md
Writes: internal/conformance/skip_ownership_test.go, internal/conformance/check_bindings_test.go, internal/conformance/tier_live_tree_test.go, internal/conformance/fixture_bite_test.go
Covers: LTE30, LTE31, LTE34, LTE46

## What to build

Review chunk: LTE-C6.

Move skip ownership onto the accepted snapshot interface. Preserve root files, tests, owner exemptions, renamed receivers, text-only mentions, and the no-Skip parse bypass.
Complete composed verification through RunConformanceSelection and the four real bindings. Use a source file that the visitors' domains overlap.
Prove one directory observation, read, and parse for equal subjects. Prove separate observations and diagnostics for distinct root and kit subjects.

Detect surviving direct source traversal, reads, and parsing in migrated visitors and their source helpers. Derive the visitor identities from the executable bindings.
Record explicit metadata-read exceptions that are outside Go scanning. Do not permit an exception to cover a duplicate source walker.
Demonstrate one red for per-binding snapshots and another for a restored old walker.

Retain fresh observations through the registered fixture mutation and restoration path. This ticket owns the combined proof because all four consumers now exist.

## Acceptance

- [ ] Skip diagnostics equal the prior visitor for the complete declared domain and error family.
- [ ] A malformed source without Skip retains its existing bypass.
- [ ] Different root and kit subjects produce their respective diagnostics.
- [ ] The four real bindings share counted observations through the dispatcher.
- [ ] Per-binding snapshots make the composed proof fail.
- [ ] Restoring a direct old walker makes the surviving-call proof fail.
- [ ] Registered fixture mutation and restoration each observe current bytes.

## Verification

Run skip equivalence, composed sharing, subject routing, metadata, and the retained fixture universe. Execute both named sharing mutations and restore each exactly. Run the complete conformance package at the chunk checkpoint.

