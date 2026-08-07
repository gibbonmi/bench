# Repair clean-skip propagation

Blocked by: none
Ownership fence: `internal/adopt/link_transaction.go`, `internal/contract/surface/link_lifecycle_test.go`
Integration surfaces: clean-skip predicate→`internal/adopt/link_transaction.go` + CSP1; link lifecycle contract→`internal/contract/surface/link_lifecycle_test.go` + CSP1
Contracts: none crosses
Closure: CSP1/relink-propagates, CSP2/upgrade-propagates

## What to build

Make the clean-skip predicate prospective. The plan loop currently skips any
existing destination that matches the *old* manifest hash, so a kit asset whose
bytes changed between releases while the destination stayed untouched is never
rewritten: relink and `bench upgrade` become content no-ops for every clean
entry, and the retained manifest row diverges from the shipped kit permanently.
Skip an entry only when the destination already matches the bytes and mode this
plan would write from its kit source; a destination matching the old manifest
but not the incoming source needs a write and proceeds through the unchanged
symlink-parent and conflict semantics. Keep the genuinely-unchanged skip
observable by inode stability (the existing clean-entry fixture must stay
green), and keep both symlink-parent aborts green. Add the two-kit fixture the
coverage row demands: kits differing in one shared managed file, asserting
destination content and manifest hash after relink and after upgrade.

## Acceptance

- [ ] [CSP1] (covers local) `bench link` over a destination that matches the old manifest but not the incoming kit bytes rewrites the file and records the new manifest hash.
- [ ] [CSP2] (covers local) `bench upgrade` between two kits differing in one shared managed file propagates the new bytes and manifest hash to an untouched destination.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CSP1/relink-propagates | key the skip on the old manifest hash (`manifestOwnedClean`) instead of the prospective source bytes | link lifecycle contract | relink kitB over a kitA-converged tree, expect the content and manifest-hash assertions red |
| CSP2/upgrade-propagates | retain the old manifest row for a skipped-but-changed entry | link lifecycle contract | upgrade kitA→kitB with one shared file changed, expect the manifest-hash assertion red |
