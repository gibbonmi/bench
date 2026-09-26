## Outcome

FT337 has a staged spec and two serial implementation tickets.
Opus/high accepted the spec after two review passes and the supported nonblocking edits.
Human approval remains pending. Neither implementation ticket has started.
The accepted author tip is `4e461ebbb5d1c3033f37764de82c8d59759d14f4`.
The phase includes the debug landing `3119dc72d48ad44a0a382bdf13a437a877b3e538`.

The spec defines complete file evidence, stable digests, source identity, and the author-record order.
A real narrow/control comparison remains an implementation checkpoint obligation.
The spec lane, build, coverage, prose, and preflight checks passed.
This capture precedes the spec staging gate and landing. Those outcomes remain pending.

## Gate-stage timings

The scaffold selected the latest landing span, which belongs to the debug prerequisite.
These timings do not grade the spec staging candidate.

| stage | debug prerequisite |
|---|---|
| gofmt | 124 ms |
| vet | 1202 ms |
| test | 130850 ms |
| race | 4022 ms |
| system | 44879 ms |
| shellcheck | 591 ms |

Debug landing: `3119dc72d48ad44a0a382bdf13a437a877b3e538`.
Trace: `b8bd767b0b7fe8cd939297bc451a5955`.
Spec staging gate: pending.
End-to-end latency and author or coordinator usage remain unknown.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| spec forecast | One blocking repair round follows initial review. | claimed | 7 | held | unknown / unknown / spec author |
| review pass 1 | B1 leaves file identity undefined. | claimed | 7 | pending | claude-opus-5-5 / high / spec reviewer |
| review pass 2 | B1 is resolved. | claimed | 8 | pending | claude-opus-5-5 / high / spec reviewer |
| review pass 1 | Any test passed. | abstained |  | pending | claude-opus-5-5 / high / spec reviewer |
| review pass 2 | Any test passed. | abstained |  | pending | claude-opus-5-5 / high / spec reviewer |

The repair table labels the author forecast. Its Brier mean is 0.09 over one pair.
The reviewer has no labeled pair and two explicit abstentions about test results.
The coordinator retained both terminal reports. It did not infer test execution from their verdicts.
The author and coordinator model metadata and provider costs remain unknown.

| review pass | input | cache creation | cache read | output | provider cost, USD |
|---|---|---|---|---|---|
| 1 | 62 | 101965 | 1924714 | 27669 | 1.7542908 |
| 2 | 40 | 76055 | 1061289 | 15995 | 1.1407578 |

The reported provider total is USD 2.8950486.
The terminal result files supply these figures, without reconstructed message sums.
Sources: `/tmp/ft337-spec-review-1-result.json` and `/tmp/ft337-spec-review-2-result.json`.

## Coordinator catches

The first review found undefined rename identity and paths without text headers.
Git fixtures confirmed the inventory mismatch and the missing header shapes.
The author defined exact descriptors and retained ambient rename behavior.
The second review supplied four nonblocking clarifications, which the author folded without another paid review.

The coordinator kept actual review evidence outside the ticket author's acceptance criteria.
It also distinguished a green debug landing from the pending spec staging gate.
A merge invoked from the primary checkout used caller paths for the target's prose lane.
The same merge passed unchanged when the coordinator invoked it from the target worktree.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-order-review-preparation.md | 0 implementation rounds; not started | not applicable at spec staging |
| 2-page-review-diffs-by-file.md | 0 implementation rounds; not started | not applicable at spec staging |
| spec authoring | 1 blocking repair round; 1 nonblocking close fold | spec/ticket: the first draft did not settle headerless patch identity |

Authoring used three passes in total. Those passes are not ticket build rounds.
Daemon recovery restored the retained work without a new research cycle.
The separate debug retrospective owns the guidance repair's implementation count.

## Agent-experience improvements

### Bench CLI

- Resolve merge lane paths from the target checkout, independent of the caller's checkout.
  The unchanged retry exposed a caller-root defect. The capture learning retains the exact symptom.
  Feeds: new

### Skills

- Observe headerless and repeated-path output before declaring file identity in a spec.
  The existing learning records this proposed authoring rule.
  Feeds: none

### Process

- Keep the real control comparison as an orchestrator checkpoint after the ticket commit.
  A ticket author cannot deliver independent review of its own work.
  Feeds: none
