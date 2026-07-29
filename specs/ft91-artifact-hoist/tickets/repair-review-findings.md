# Repair concrete semantic-review findings

Blocked by: Build the package-scoped shared artifact set

## What to build

Keep the promoted artifact count single-sourced from the fingerprint traversal,
and make shared-set consumers follow the repository's existing privileged-runner
capability posture when UID 0 would bypass the read-only belt.

## Acceptance

- [ ] The shared set derives its promoted entry count from the fingerprint map without a second directory traversal.
- [ ] A privileged runner capability-skips shared-set consumers before staging because mode-bit write protection is unavailable.
- [ ] Non-privileged focused singleton and sharer tests remain green.
