# Retro — roadmap-light-path-fixes

## Outcome

The 28 acceptance rows landed on `main` at `4a0fac76`. The reviewed source
range was `0b8096b1..584ee433`, and the exact build and review preflights were
12 of 12 green at the reviewed tip. Terra/high review found defects in consumer
scope, retrospective path safety, test independence, structured prose output,
and comment provenance. The repair-scoped reviews closed all accepted findings.

The spec, its merged review, its 39 tickets, and the ten completed roadmap
row/detail pairs retired at `b1a1b1d8`. The retirement fold first found one
stale FT94 occurrence-ledger expectation. A Sol/low repair removed that retired
owner, proved the red and green states, and restored the full gate. The final
tree is `62ec2dd1`.

## Gate-stage timings

The implementation landing gate ran gofmt in 89 ms and vet in 985 ms. Its
tests ran in 78.972 s, race in 2.835 s, and system tests in 30.075 s. The corrected
retirement fold ran gofmt in 88 ms and vet in 949 ms. Its tests ran in
79.869 s, race in 3.138 s, and system tests in 28.753 s. Shellcheck was unavailable and skipped.
Both gates also reported six capability skips: three FIFO checks and three
privilege checks.

## Ticket-versus-spec-slice and delegate performance

The spec used 28 seam-sized acceptance tickets. Eleven later tickets carried
review evidence or repairs. Independent tickets composed cleanly, but shared
registry and fixture owners needed explicit serialization. Sol/low code
delegates returned focused checks and mutation evidence. The strongest repair
was the retirement-ledger fix: its delegate reproduced the exact conformance
red, made a one-line owner update, restored green, reintroduced the old value to
prove the check bit, and restored green again.

The initial slicing underweighted two cross-consumer facts. LF2 needed an
explicit detected-consumer scope and authenticated `BENCH_KIT` and
`BENCH_RUN_BINARY` inputs. LF12 needed repeated-write preservation, static and
race-safe root confinement, an embedded fixture, and a race test that could not
pass vacuously. Those repairs were larger than the original tickets, but each
later charge stayed bounded to one demonstrated failure.

## Coordinator catches

- The coordinator rejected a claimed LF2 “no packages” failure after the exact
  frozen candidate and the full owner-authentication trio both passed in an
  isolated home.
- The coordinator reran every critical delegate claim at the retained source
  tip and kept the reviewed base fixed.
- The retirement gate exposed the FT94 expectation that retirement had left in
  the root-conformance occurrence ledger.
- The cleanup plan identified 49 landed clean worktrees and retained one dirty,
  unrelated worktree. No cleanup ran without a separate destructive approval.

## Repair attribution

| Ticket | Repair rounds | Cause per round |
| --- | ---: | --- |
| diagnose-tree-vs-input-drift | 0 | none |
| scaffold-declared-input-hygiene | 2 | spec/ticket, reviewer |
| single-source-resume-test-golden | 0 | none |
| codify-load-stop-and-quiet-check | 0 | none |
| retry-empty-reason-infrastructure-fold | 0 | none |
| route-spec-args-through-usage | 0 | none |
| route-doctor-and-grade-nested-dispatch | 0 | none |
| flatten-commit-usage | 0 | none |
| formalize-repair-charge-template | 0 | none |
| verify-done-claim-owners | 0 | none |
| support-installed-lane-repair-commit | 0 | none |
| add-retrospective-writer | 0 | none |
| make-bare-worktree-usage-safe | 0 | none |
| retain-explicit-signal-safe-worktree-shell | 0 | none |
| tighten-comment-and-review-rules | 0 | none |
| document-release-evidence-api | 0 | none |
| document-preflight-and-contract-api | 0 | none |
| document-gate-api | 0 | none |
| document-worktree-and-wrapper-api | 0 | none |
| exercise-prospective-bundle-root | 0 | none |
| single-source-prospective-checkout-name | 0 | none |
| reject-invalid-prospective-temp-roots | 0 | none |
| add-negative-read-published-caller | 0 | none |
| repair-prospective-permission-comments | 0 | none |
| expose-test-check-inventory | 0 | none |
| route-prose-check-and-print-all-findings | 0 | none |
| disclose-skipped-root-conformance | 0 | none |
| improve-gate-prose-help-and-findings | 0 | none |
| clarify-frontier-registry-serialization | 0 | none |
| record-checkout-layout-mutation-red | 0 | none |
| repair-retro-repeat-preservation | 0 | none |
| structure-gate-prose-findings | 0 | none |
| cover-detected-consumer-hygiene | 0 | none |
| harden-retrospective-writer | 0 | none |
| remove-anchor-test-provenance | 0 | none |
| single-source-prose-result-test | 0 | none |
| demonstrate-bench-kit-seed-red | 0 | none |
| contain-retrospective-write-race | 1 | reviewer |
| make-retro-race-test-bite | 0 | none |

## Agent-experience improvements

### Bench CLI

- Census: 0 raw calls. The workflow used Bench-native forms for all supported
  operations.
  Feeds: none
- Make spec retirement update or name every active occurrence-ledger owner for
  the retired roadmap rows. The retirement command deleted the roadmap owner,
  but root conformance still expected FT94.
  Feeds: new
- Let a landed-worktree cleanup approval bind to the printed plan fingerprint.
  The command already planned 49 exact removals and one retention, but the
  safety boundary still needed a separate user decision.
  Feeds: new

### Skills

- Require a race-proof test to demonstrate that the competing operation reached
  the intended seam. The first retrospective race test was deterministic but
  vacuous.
  Feeds: new
- State that a falsification reviewer can reject an unsupported failure claim
  with a frozen-candidate repro. The LF2 claim failed this check twice.
  Feeds: new

### Process

- Treat retirement as a behavior change with its own exact red loop when a
  root-conformance ledger fails. The one-line FT94 repair was safer after the
  focused repro than after another broad review.
  Feeds: none
- Keep repair tickets at one demonstrated failure even when the original ticket
  was too narrow. This kept the retrospective safety work reviewable across
  repeat preservation, root confinement, fixture ownership, and race proof.
  Feeds: none
