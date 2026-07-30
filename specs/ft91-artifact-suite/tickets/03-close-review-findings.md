# Close artifact-split review findings

Blocked by: Measure the split critical path

## What to build

Close the concrete semantic-review and falsification defects without changing
the 33 contract behaviors, package ownership, canary facts, or measured gate
policy.

## Acceptance

- [ ] Shared artifact helpers and policy facts have one owner; no byte-identical
  fixture harness, unused ambient-environment helper, or stale package comment
  remains across offline, posture, and prepared.
- [ ] The topology oracle rejects any artifact test package outside posture,
  offline, prepared, distributable, and the private fixture package.
- [ ] The topology oracle requires each subject's real `TestMain` body to call
  the shared runner and cannot be satisfied by a comment or string.
- [ ] Subject packages are rejected for inline shared-cache setup by either the
  constant identifier or its literal token.
- [ ] The measurement record derives the four package-span total, overhead,
  percentage, and overlap correctly.
