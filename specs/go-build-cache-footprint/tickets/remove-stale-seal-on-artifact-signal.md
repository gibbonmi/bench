# Remove the stale seal on a signalled artifact install

Blocked by: single-source-fake-go-builder-stub.md
Ownership fence: `scripts/go-build.sh`, `internal/contract/surface/artifact/posture`
Integration surfaces: install-then-cleanup ordering→`scripts/go-build.sh`; signal-window observation→`internal/contract/surface/artifact/posture`
Contracts: the promoted-install marker crosses `scripts/go-build.sh`'s install line→its own cleanup trap, asserted by AS1 against the real builder under a delivered signal

## What to build

Repair for review finding
`spec-06-artifact-signal-between-install-and-seal-removal`. A handled
SIGINT/SIGTERM/SIGHUP arriving between artifact mode's `mv -- "$staged" "$out"`
and `rm -f -- "$out.seal"` exits through the traps leaving the new artifact
beside the retired subject's stale seal. Two directory entries cannot change
atomically, but the handled-signal window can close: record that the install
promoted (a flag set immediately after the `mv`), and make the EXIT-trap
cleanup finish the transaction — removing the stale destination seal whenever
the install promoted — so every handled exit path converges on the promised
unsealed artifact. Only an uncatchable kill between the two operations remains,
which is the edge inventory's won't-handle. The builder still holds no rollback
state; this is completion, not rollback.

## Acceptance

- [ ] [AS1] A handled termination delivered after the artifact install promotes but before the seal removal runs leaves the destination with the new artifact and no seal, under the real builder's own traps.
- [ ] [AS2] The ordinary artifact journey is unchanged: successful installs end with no destination seal, and pre-install failures still leave the prior output and prior seal byte-for-byte.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AS1 | drop the promoted-install seal removal from the cleanup trap | the signal-window contract with a bounded marker-controlled blocker | block after the install promotes under a deadline, deliver the handled signal, expect the surviving stale seal to fail the assertion |
| AS2 | remove the stale seal before the install instead of after | the pre-install failure table with fingerprints | fail the build stage before install, expect the missing prior seal to fail the byte-for-byte comparison |
