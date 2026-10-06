# 06 Share source observations with bounds policy checks

Blocked by: 04-share-source-observations-for-git-policy.md
Writes: internal/conformance/bounds_policy_test.go
Covers: LTE29

## What to build

Review chunk: LTE-C6.

Move bounds caller traversal and parsing onto the accepted snapshot interface. Preserve registry, read-seam, wait-policy, and offline-policy inputs.
Retain non-test caller coverage under cmd and internal, excluding the bounds owner. Preserve current registry and source failure diagnostics.
Reuse the snapshot's parse result and positions without changing expression formatting or declaration analysis.

The predecessor supplies observation ownership only. Bounds retains policy decisions and independent omission expectations.
Compare exact old and new diagnostic slices for registry, caller, read-seam, wait, and malformed-source cases.

## Acceptance

- [ ] Each existing bounds violation retains its diagnostic after migration.
- [ ] Registry absence and malformed registry source remain failures.
- [ ] Read-seam discovery retains all currently derived exported limit arguments.
- [ ] Test files and the bounds owner retain their caller exemptions.
- [ ] The existing injected-wait and elapsed-wait classifications remain unchanged.
- [ ] Caller scanning no longer owns a separate traversal or parser.

## Verification

Run bounds policy tests and the exact diagnostic comparison family. Demonstrate a bounds bypass red, then restore it. This ticket must pass while the skip visitor is still unmigrated.

