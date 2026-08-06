# Retro: go-build-cache-footprint (FT195)

## Outcome

Promoted terminal: squash `af02997` on main, candidate `caaa6fb0`, run
`1af28d4…`. The composition: three original tickets (single-`go build` sealed
subject, artifact mode without execution, `-count=1` gate test verdicts) plus
eight repair tickets across four review/gate rounds — publication-transaction
ownership in the Go child, all-package topology enumeration (settled at
one entry after the gate's attested publisher proved production-unreachable),
single-sourced gate rationale, fake-`go` stub dedup, handled-signal temp-litter
cleanup, artifact stale-seal cleanup in the EXIT trap, and the artifact
topology registry sync. Three flagged residuals remain on receipts for the
reviewer: overlapping malformed-selector tables, artifact-only hostile-output
coverage, and a cosmetic ticket-prose reflow.

## Gate-stage timings

Ticket-commit gates on main ran ~4–5 min each (five landings). The promote
prospective gate ran twice: one red (contract phase, `TestSubjectPackageTopology`),
one green. Debug reproduction of the red paid one full gate on the exported
prospective tree (~6 min; contract phase red, all others green) plus targeted
package runs: `internal/gate` 156–190 s, `internal/preflight` 168 s,
conformance suite 103 s, race over the two changed packages ~64 s. The
`-count=1` change this spec ships is itself the reason test verdicts were all
fresh executions.

## Ticket-versus-spec-slice and delegate performance

Eight repair delegates, all ticket-sized with receipt-derived fences: every one
returned in-fence and green in one pass, none hit its iteration cap, and the
two cheap-tier (sonnet) charges — proven-design topology application and the
registry sync — performed identically to mid-tier on their tight scopes. The
one wide ticket (transaction: five entries, four directories) was the only one
that generated follow-on findings, and a retro split showed it decomposed into
two independently-green two-file tickets — direct evidence for the
thinnest-green-slice simplification now queued against craft-tickets. One
delegate deviated from its ticket's stated mechanism (post-`mv` promotion flag)
with probe-verified justification (bash trap boundary semantics) and shipped a
better record (staged-entry consumption); charges should state the observable,
not the mechanism.

## Coordinator catches

- A coordinator probe aimed at the wrong owner (signal test) let an input
  mutation survive; re-aiming at the normal-path contract caught it — probe
  site selection needs the same care as mutation kind.
- The promote red itself: no whole gate runs on a candidate until promote, so
  three tickets' unregistered tests in a classified family surfaced rounds
  late; the tickets' `Integration surfaces:` lines had missed the classifier
  registry (`fixture/topology_test.go`).
- `git archive` + `add -A` re-exports silently drop tracked-but-gitignored
  files (`projects/gl-axi.md`, `projects/regroup.md`), breaking posture
  fixtures in disposable trees.
- The review's reachability finding (attested publisher never called in
  production) exposed that the earlier two-publisher decision rested on
  source presence, not call-graph reachability.

## Agent-experience improvements

### Bench CLI

- `bench spec build promote` retains only a one-line red attribution and an
  evidence digest; the debugging session had to rebuild the prospective tree
  and re-run the whole gate to learn the failing phase. Retaining the failing
  phase name and its first failing-test lines with the run state would cut a
  ~15-minute reproduction to zero.
- `bench spec build status <slug> --full` prints all historical assignments;
  a run with rounds of repairs buries the active frontier. A `--active` cut
  or terminal/active grouping would help cold pickup.
- A `promote --check` (precondition dry-run: clean checkout, recomposition
  needed) would have caught the dirty-`capture/learnings.md` refusal and the
  contract-anchor ticket refusal before paying gates.

### Skills

- craft-tickets: the sizing simplification (thinnest independently-green
  slice; keep-whole must name the stranded red) is decided and queued — this
  run is its evidence base.
- craft-tickets' classified-family rule worked when applied and bit when
  skipped; the charge template the coordinator uses should carry an explicit
  "sibling-family classifier search" line so the discovery step is visible in
  every ticket, not just remembered.

### Process

- Whole-gate feedback arrives only at promote; for compositions that add tests
  to contract-classified families, a cheap coordinator pre-promote check
  (run the classifier fixture package against the candidate) would catch
  registry drift rounds earlier without duplicating the gate.
- Capture appends must match `bench learnings` grammar (`[open]` suffix);
  the malformed-file state was silent until `bench status` flagged it.
