# Sweep staged specs for fence/probe contradictions at the gate

Blocked by: enforce-fence-probe-agreement.md
Ownership fence: `internal/conformance`, `tests/canary`, `projects/benchkit.md`
Contracts: the agreement verdict crosses `internal/specbuild`→the sweep in `internal/conformance`, asserted by SW1 against the real exported validation over real staged ticket files, never a second derivation of the rule; the check's registration crosses the conformance registry→the profile's per-check table in `projects/benchkit.md`, asserted by SW3 through the real gate

## What to build

A conformance check sweeps every ticket of every spec whose `Status:` line reads
staged and runs the exported agreement validation from the previous ticket over
each. A violating ticket turns the gate red naming the file and the row.
Implemented specs are skipped: their tickets are history awaiting
promote-then-delete retirement, and no build will charge them again.

Register the check in the conformance registry with its input source, add its
row to the profile's per-check table, and give it a canary fixture whose broken
ticket proves the check bites with a targeted substring. If any currently staged
spec violates the rule when this lands — `specs/recovery-discard/tickets/add-spec-build-reclaim.md`'s
RM5 row, should that run end abandoned rather than promoted — repair the
artifact in this same green change rather than landing a red gate.

Fence spans three directories: one check advertised in the registry, the canary,
and the profile table, and a red gate between any pair is not a landable
intermediate state.

## Acceptance

- [ ] [SW1] the sweep reds a staged spec carrying a ticket whose probe contradicts its fence, naming the ticket file, and consumes the exported validation rather than its own copy of the rule.
- [ ] [SW2] a spec whose `Status:` line is not staged is skipped, and a staged spec with no tickets sweeps green.
- [ ] [SW3] a canary fixture with a deliberately contradictory staged ticket goes red with the check's targeted substring.
- [ ] [SW4] no staged spec in the tree violates the rule after this ticket lands.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SW1 | hardcode the sweep's verdict green | the canary sweep | neuter the check body, run `bench gate`, expect the canary phase to red on the fixture's missing failure substring |
| SW2 | sweep every spec regardless of status | the status-scope test | drop the staged filter, run `go test ./internal/conformance -run Sweep -timeout 300s`, expect a nonzero matched-test count and the implemented-spec-skipped assertion to fail |
| SW3 | blank the fixture's EXPECT | the canary baseline | empty the EXPECT file, run `bench gate`, expect the vacuous-EXPECT rejection to red the canary phase |
| SW4 | reintroduce a contradictory row into a staged ticket | the sweep itself | add a probe outside the fence to any staged ticket, run `bench gate`, expect the conformance phase to red naming that file |
