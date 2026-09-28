# Cover symbolic-ref boundaries

Blocked by: none
Writes: projects/benchkit.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading
Covers: none

## What to build

Refine the symbolic-ref checklist entry. Cover each transaction ref at every read, write, and delete boundary with direct refs and symbolic refs.

## Acceptance

- [ ] The checklist covers each transaction ref at every read, write, and delete boundary.
- [ ] Each boundary uses direct refs and symbolic refs.
- [ ] The checklist keeps the target-ref deletion requirement and its existing anchors.
