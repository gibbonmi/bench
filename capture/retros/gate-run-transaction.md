## Outcome

Gate-run-transaction landed as `379035dfab8865edff0531f1f268d4ab4432d3d9`
from reviewed pair `59faeec629865eda410b00f746a95280556ac125` →
`b7d38bd1f06562f66812c7bbbf0dac22803979c1`. The run transaction now owns
locking, pending and terminal persistence, evidence retention and invalidation,
dispatch, drift, and timeout outcomes; verdict loading uses one five-class
registry. Retirement is `dac455bc14caefff32e57b8a196a0acf54887adb`.

## Gate-stage timings

Landing gate `20260815T173305.118245547Z-241672` was green in 58.589s:
gofmt 0.076s; vet 1.507s; test 48.305s; race 4.243s; system 2.384s;
shellcheck 0.277s.

## Ticket-versus-spec-slice and delegate performance

Six serial tickets stayed independently green: two behavior-characterization
slices, transaction extraction, engine-seam deletion, record-class registry,
and registry-derived narrow reuse. Luna produced useful biting tests and compact
implementations at xhigh/max, but every substantive Luna lane required coordinator
scope or completeness inspection; Terra/high semantic review found six actionable
issues across two cycles and then returned two consecutive zero-finding exact-tip
cycles plus a zero-finding reconciliation review. Per-agent token counts were not
available, so the cost conclusion remains qualitative: at Luna's stated 0.2×
per-token price, the observed rework likely remained dollar-positive while losing
token and wall-clock efficiency to Terra.

## Coordinator catches

- Replaced a direct process constructor that violated the architecture census.
- Corrected two unreachable spec mutations: pre-lock reuse for GC5 and the
  ineffective pending rewrite for GT5.
- Required retained public registry tests after a delegate returned probes only.
- Narrowed a registry-reason draft back to its ticket write fence.
- Restored a base-owned verification log that a repair delegate over-deleted.
- Turned semantic-review gaps into cancellation, pending-write, owner-write, and
  single-sourced fixture coverage, each with an independent mutation receipt.
- Reconciled staged-spec bytes through main and refreshed exact review before
  landing; fingerprinted cleanup removed 13 ignored logs after publication.
- A mistaken `bench worktree clean --help` probe exposed missing inventory grammar;
  the canonical command list now carries the plan/apply syntax in `2d2bcbeb`.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | ---: | --- |
| characterize-gate-run-outcomes | 2 | delegate-error; delegate-error |
| characterize-gate-contention-and-persistence-failures | 2 | shaping-ambiguity; spec-row |
| extract-gate-run-transaction | 0 | none |
| delete-gate-engine-seam | 0 | none |
| register-verdict-record-classes | 1 | delegate-error |
| derive-narrow-verdict-reason | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

The canonical inventory now spells the fingerprinted `bench worktree clean`
grammar, including ignored-residual disposal. Landing should additionally emit the
exact staged-spec reconciliation route when destination and source bytes differ.

### Skills

Cap-accepted spec verification needs a reachability probe for folded red signals.
Transaction-shaped specs need a failure matrix with pre-oracle persistence,
in-oracle interruption, and terminal persistence as separate coverage classes.

### Process

Maintain a bounded model/effort scorecard that replaces representative evidence
and rolls aggregates forward instead of appending a session diary. Separate
delegate-caused rework from spec-origin review findings so routing decisions do not
penalize an implementer for upstream omissions.
