# Record the hoist evidence

Blocked by: Build the package-scoped shared artifact set

## What to build

Verify that the gate-critical-path decision map and its three assets live only
under the FT91 spec with no stale references, then measure the post-change
artifact package and full dev gate and record both results in the spec-local
timeline asset.

## Acceptance

- [ ] The decision map and all three map-owned assets exist only under `specs/ft91-artifact-hoist/decisions/`, and every reference resolves there.
- [ ] A fresh artifact-package timing and a fresh full dev-gate timing are recorded in the timeline asset.
- [ ] The recorded evidence states the FT91 stop-rule result against the 60-second threshold.
