# Pin the stable promotion owner

Blocked by: none
Writes: bin/bench.sh, bin/bench-postinstall.sh, internal/adopt, internal/runbinary, internal/worktree/land.go, internal/worktree/land_freshness_test.go, internal/worktree/land_identity.go, internal/worktree/land_identity_test.go, internal/systemtest, cmd/bench/command_registry_test.go

## What to build

Give public landing a dedicated wrapper route to the installed promotion broker.
Authenticate the broker through the installation manifest.
Resolve the destination and assignment as separate subjects.
Remove current-directory selection and inherited routing from this command.

## Acceptance

- [ ] SOL01 proves that candidate landing code does not run during its own promotion.
- [ ] SOL02 selects one owner from every supported current directory.
- [ ] SOL03 refuses each inherited routing override before repository reads.
- [ ] SOL04 keeps one owner process through publication and release.
- [ ] SOL05 invalidates a changed request before composition.
- [ ] SOL06 invalidates a changed review base before composition.
- [ ] SOL07 invalidates a changed source tip before composition.
- [ ] SOL08 invalidates a changed source fingerprint before the gate.
- [ ] SOL16 reports the release or repair action for a broker change.
- [ ] SOL17 ignores a forged primary executable and seal.
