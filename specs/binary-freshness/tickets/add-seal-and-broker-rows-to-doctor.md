# Add seal and broker rows to doctor

Blocked by: none
Writes: internal/adopt/doctor_rows.go, internal/adopt/doctor.go, internal/adopt/broker.go, internal/adopt/broker_test.go, internal/adopt/setup_test.go, internal/brokermanifest (new), internal/conformance/broker_reasons_test.go (new)
Covers: BF3, BF4

## What to build

Verify the premise first: `doctorRows` in internal/adopt/doctor_rows.go lists
nine rows, and none reads a seal or the broker manifest. Then move
`WriteBrokerManifest` and `ReadBrokerManifest` into a new leaf package
`internal/brokermanifest`, imported by `internal/adopt`. The publish ticket
then imports it from `internal/freshness` without a cycle. Export
`kitSourceCheckout` as `KitSourceCheckout` for the landing-line ticket.

Add two rows. The seal row calls `freshness.Verify` on the kit's `dist/bench`
and prints the verdict with `RebuildAction` on a mismatch. The broker row reads
the manifest and grades it with the five predicates the land route applies.
They are present, complete, version equals the installed package, path is a
regular non-symlink executable, and digest matches the bytes. Each reason
string matches the land route's wording.

Add one authored conformance expectation that lists the five reason strings
and asserts both `bin/bench.sh` and the doctor row carry each one. The route
is shell before any binary is trusted, so the two derivations are necessary,
and the expectation is what keeps them one.

## Acceptance

- [ ] `bench doctor` on a kit with a stale seal prints the seal row with the source-digest reason and the sentence.
- [ ] `bench doctor` on a kit with a `dev` manifest prints the broker row with the version reason.
- [ ] The reason-list expectation passes on the live tree and reds when one reason is removed from either side.
- [ ] Self-probe: drop the digest predicate from the broker row, and report the expectation red.
