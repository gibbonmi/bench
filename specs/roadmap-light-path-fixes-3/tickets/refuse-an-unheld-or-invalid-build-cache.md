# Refuse an unheld or invalid build cache

Blocked by: derive-the-canonical-path-in-one-leaf-package.md
Writes: internal/gocache/gocache.go, internal/gocache/gocache_test.go, internal/gate/run_transaction.go, internal/gate/lane.go, internal/gate/subject.go, internal/gate/report.go, internal/testreport/command.go
Covers: LQ17, LQ18, LQ19, LQ20, LQ28

## What to build

Verify the premise first. The three `gocache.Hold` callers in
internal/gate/run_transaction.go, internal/gate/lane.go, and
internal/testreport/command.go run on `err != nil`. `Apply` strips every
inbound `GOCACHE` entry and appends the home derivation, so a relative inbound
value never reaches a child. `FromEnv` returns the inbound entry verbatim to
the footprint report, and `HoldDir` creates the derived directory and opens the
lock file inside it.

Then make each `Hold` caller refuse on a hold error with the adopt
`toon.Errorf` shape, exit 1, the error text, and the sanitized cache path.
Rewrite each site's comment, because each states the old policy. A hold error
already covers an unwritable derived directory, because the lock-file open
fails there. Make `FromEnv` return a refusal for a relative inbound entry, and
make the footprint report print it. Leave the empty-value fallthrough and the
`HoldDir` create in place: an absent derived directory is created, not refused.

The gate subject closure in internal/gate/subject.go drops the entry on an
`Apply` error today. Make it propagate the error, so a refused cache cannot
change the oracle hash in silence.

## Acceptance

- [ ] An unwritable derived cache directory refuses the gate transaction, the lane, and the focused run, each printing the error and the path.
- [ ] `TestFromEnvFallsBackToTheHomeDerivation` still passes for `GOCACHE=`.
- [ ] A new `FromEnv` row refuses `GOCACHE=cache` and the footprint report prints `cache`.
- [ ] `TestHolderCreatesTheDirectoryAndTheLockFile` still passes, so an absent derived directory is created.
- [ ] The subject closure returns the `Apply` error instead of dropping the entry.
- [ ] Self-probe: restore `err == nil` at one call site, and report that site's refusal test red.
