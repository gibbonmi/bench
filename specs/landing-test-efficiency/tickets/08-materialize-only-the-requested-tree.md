# 08 Materialize only the requested prospective tree

Blocked by: none
Writes: internal/gate/prospectiveartifact/prospectiveartifact.go, internal/gate/prospectiveartifact/materialize_test.go (new)
Covers: LTE35, LTE36, LTE37

## What to build

Review chunk: LTE-C7.

Change Owner.Materialize to populate the requested tree without first populating HEAD. Preserve the owner record, Git registration, index, and teardown contracts.
Use a focused implementation probe to choose and verify the Git operation. Do not assume an untested flag supplies the required behavior.
The real-Git proof compares the resulting index tree and filesystem entries with the requested tree.

Include additions, deletions, content changes, executable modes, symlinks, and paths with spaces or glob characters.
Count tracked-file population through an observed Git-operation trace. A successful final checkout alone cannot prove that duplicate work disappeared.
Keep Materialize's public signature so standalone execution, inspection, and lane callers need no migration.

## Acceptance

- [ ] The resulting index names exactly the requested tree.
- [ ] Requested additions, deletions, and changed contents appear in the filesystem.
- [ ] Executable modes and symlink targets match the requested tree.
- [ ] Hostile path bytes reach Git as arguments without shell interpretation.
- [ ] The trace contains one tracked-file population, and a preliminary HEAD population fails the test.
- [ ] Existing registration ordering, owner-record publication, and close-confinement tests still pass.

## Verification

Run the prospectiveartifact package and the focused real-Git probe. Demonstrate a HEAD-only result failure and a double-population failure. Retain existing malformed, foreign, non-private, dead-owner, and live-owner cases.

