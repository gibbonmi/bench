# Prove every retained fixture through its owner

Blocked by: none
Writes: `internal/conformance`

## What to build

Replace the partial hand-maintained fixture-bite coverage with one test-internal runner that derives cases from `canary.Fixtures`, resolves each fixture's exact registered non-meta owner and declared tier, and proves the fixture expectation appears after materialization and disappears after restoration. Keep the runner and live wrapper in the architecture-censused fixture proof file; do not add a production journey, process constructor, nested sweep, or focused FT168 surface.

## Acceptance

- [ ] (covers PB1) All 184 retained fixtures receive direct attributed subtests through their exact owner: 182 at dev and the two `release-evidence-probe` fixtures at ship, with the enumerated 31-member gap closed and no retirement.
- [ ] (covers PB2) A synthetic canary directory containing a registered-family dev fixture and an explicit-`CHECK` ship fixture drives both through the same producer-derived runner; stored names and tier or binding-form filters are red.
- [ ] (covers PB6) The live test calls the runner once, records identity only after `runFixtureBite` returns, compares the observed set with a fresh `canary.Fixtures` key set, and an AST grader rejects fixture-name literals, manual sources, a second proof loop, or record-before-proof.
- [ ] (covers PB3) Each fixture's normalized non-empty `EXPECT` appears after mutation and is absent after restoration through the same registered owner.
- [ ] (covers PB4) Missing fixtures; absent, zero-byte, ASCII-whitespace, or non-ASCII-whitespace expectations; and unbound, unregistered, or meta checks fail with fixture attribution, while valid owners run at their declared dev or ship tier.
- [ ] (covers PB5) `fixture_bite_test.go` is an exact `directArchitectureTests` member and admits no generic repository/process constructor, gate, wrapper, `go test`, or `go run`; native ship-owner implementation files remain outside that harness census.
