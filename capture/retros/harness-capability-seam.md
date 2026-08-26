# Retro: harness-capability-seam

## Outcome

The spec landed as `1100d6bd` on 2026-08-26 from the reviewed pair
`c35f1516..a512bb1b`, with the whole-project gate green in the landing run.
Nine tickets, one fence amendment, two review repair tickets, one fixture
repair, and one CI-parity fix committed serially on one integration source.
`internal/harnesses` now owns one record with one row per harness, and
`lines`, `status`, `guards`, and the conformance package derive their harness
lists from it. `bench harnesses` projects the record, `bench status --route
--harness none` routes past every phase action, and the `harness-record` and
`entry-point-parity` checks grade the record and the shims against the tree.

The landed record carries `census=55`. The final check wrote one census
learning for the landing.

## Gate-stage timings

| phase | elapsed |
| --- | --- |
| gofmt | 73 ms |
| vet | 1.4 s |
| test | 42.1 s |
| race | 4.4 s |
| system | 10.1 s |
| shellcheck | 0.5 s |

The landing went green on the first attempt. The whole gate took 65.9 s.

## Ticket-versus-spec-slice and delegate performance

Every charge was ticket-sized; no charge received a spec slice. Opus/high
carried tickets 01 and 08 and each landed first-pass on behavior. The reviewer
then moved every later charge to Opus at low or medium. Opus/medium carried
tickets 02, 03, 05, 06, 07, 09, the two repair tickets, and the three review
axes. Opus/low carried ticket 04 and the fixture repair.

Every ticket landed first-pass on behavior. Ticket 06 stopped twice for a
reviewer decision: the `cells[13]` shape and the AXI envelope opt-out. Both
stops were correct.

## Coordinator catches

- Ticket 06 shipped HC28 with no test; the delegate had verified the wrapper by hand. The coordinator's label probe stayed green and exposed both the missing test and the hand-kept `benchShRoutes` list.
- A `main` change to `scripts/release-preflight.sh` during the build would have redded the new CI-parity row at landing. The coordinator found it by comparing `main` against the frozen base before the landing.
- The re-review found the live-symlink adapter fixture planted two faults through the shared planter; the delegate had reported it as a note, not a defect.
- One repair delegate reported three live-root reds from its own stale checkout; the coordinator attributed them by re-running on the source.
- Every coordinator probe bit at a site and kind distinct from the delegate's, except one vacuous wrapper-label probe that led to the two findings above.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| 01-add-the-harness-record-package | 0 | none |
| 02-derive-the-binding-matrix-from-the-record | 0 | none |
| 03-derive-the-status-harness-grammar-from-the-record | 0 | none |
| 04-derive-the-guard-and-conformance-harness-loops | 0 | none |
| 05-route-past-a-formless-harness | 0 | none |
| 06-add-the-bench-harnesses-verb | 3 | spec-row, spec-row, ticket-slicing |
| 07-add-the-harness-record-conformance-check | 1 | spec-row |
| 08-add-the-entry-point-parity-conformance-check | 2 | spec-row, tree-drift |
| 09-document-the-harness-record | 1 | spec-row |
| 10-repair-the-two-checks-after-review | 1 | delegate-error |
| 11-repair-the-record-and-its-readers-after-review | 0 | none |

Ticket 06's rounds are the `cells[12]` count, the AXI envelope fixture the
spec never named, and the HC28 test the ticket did not require. Ticket 07's
round is the live-symlink adapter refusal the edge inventory named and the
check did not do. Ticket 08's rounds are the silent absent static entry and
the `main` change to the CI script. Ticket 09's round is the one-source claim
about the adapter list. Ticket 10's round is the two-fault fixture.

## Agent-experience improvements

### Bench CLI

- Add `--env KEY=VALUE` to `bench worktree exec`, so a delegate never sets a pool-path shell variable; the census learning `hcs-integration census 55` records 74 such calls.
  Feeds: new
- Add a `bench probe` verb that copies a file aside, applies one named mutation, runs the focused test, reports the expected red, and restores. Each probe cost four to five calls, and the console shows a by-design failure.
  Feeds: new
- Add `bench gate --check <name>` so one conformance check runs over the live root without a hand-built inherited list; three delegates asked for it.
  Feeds: new
- Make `bench worktree clean --landed --apply` remove every landed tree, dirty ones through recover-remove, under one fingerprint; nine trees cost about twenty-five serial calls.
  Feeds: new
- Print the usage line at exit 0 for `bench worktree release --help`, and keep a retained assignment reachable by `bench worktree exec`.
  Feeds: new

### Skills

- Make `craft-tickets` require an AXI-verb ticket to name every fixture registry in `cmd/bench`, including the help inventory and the envelope cases. Ticket 06 named five of eight.
  Feeds: new
- Make `craft-spec` require a coverage row's inline TOON header to match the design text's shape; HC25 said `cells[12]` where the design said `cells[N]`.
  Feeds: new

### Process

- Compare `main` against the frozen base before a landing and grade any file a new oracle check reads; the CI script changed on `main` mid-build.
  Feeds: none
- Run per-path worktree cleans serially; the plan fingerprint binds pool-wide state, so parallel applies stale each other.
  Feeds: none
- A hand-verified acceptance row is not closed; the coordinator asks for the test name before it folds a done-claim.
  Feeds: none
