# Bind the proof scripts to the proven targets

Blocked by: add-proven-target-field.md
Writes: scripts/native-proof.sh, scripts/aggregate-native-proofs.sh, internal/conformance

## What to build

The two proof scripts read the proven list rather than the shipped list.

`scripts/native-proof.sh` calls `proof-target` in place of `target`, so an unproven target exits non-zero on the existing matrix message. The script's Darwin branch goes with it, because no proven target reaches that branch. The Linux branch stays whole, and `musl_status` keeps both of its values.

`scripts/aggregate-native-proofs.sh` iterates the proven targets and compares the directory against the proven proof file set. A missing proven proof stays red, and a proof file for an unproven target stays red.

This ticket opens the new script execution seam. It adds a test in `internal/conformance` that runs each script against a temporary proof directory and reads the exit status and the message.

## Acceptance

- [ ] The aggregator fails when a proven target has no proof file (row B6).
- [ ] The aggregator fails when the directory holds a proof file for an unproven target (row B7).
- [ ] `native-proof.sh` exits non-zero when it is called for an unproven target (row B8).
- [ ] `native-proof.sh` contains no Darwin branch and no `nm` assertion (row B11).
- [ ] The aggregator still reports the verified count on a complete proven set.
