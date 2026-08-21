# Retro — light-path-wrapper-home-guard

## Outcome

Audit item A4's last unshipped acceptance criterion landed at `31d64ac7`.
`bin/bench.sh` derives the pool home through
`${BENCH_HOME:-${HOME:?...}/.bench}`, so a session with neither variable set reads
the wrapper's own message naming both inputs instead of `HOME: unbound variable`.
`BENCH_HOME` still wins when set, and the default derivation is byte-identical.
`internal/systemtest/wrapper_pool_home_test.go` pins the wording independently of
the wrapper; a mutation probe of omission kind confirmed a revert reds it.

Three commits: the fix, the learnings capture (`ba48496f`), and the handoff
refresh. A light path with zero delegates — one ticket, implemented inline, per
the right-size table's standing approval. The ticket folder retired on the
landing commit through `bench commit --spec`.

Two pieces of tree state cost time before any code was written. The branch had
diverged from `origin/main`, and the git guard correctly refused the rebase, so
the reviewer reconciled it. Separately,
`docs/audits/2026-08-bench-capability/results-fable-high/next-ticket.md` still
names A1, which shipped 2026-08-18; the charge flagged it, and it stays stale.

## Gate-stage timings

From the landing commit's own gate run
(`.logs/gate-20260821T020623...jsonl`), all six phases green:

| phase | elapsed |
| --- | --- |
| gofmt | 0.08 s |
| vet | 1.1 s |
| test | 53.0 s |
| race | 4.0 s |
| system | 3.0 s |
| shellcheck | 0.4 s |

The `test` phase is 87% of the run. Every one of the four full gate runs this
landing paid sat within 1% of that figure, so a light path's real cost here is
roughly one minute per commit regardless of diff size — three of these four runs
graded a one-line shell change or a markdown file.

## Ticket-versus-spec-slice and delegate performance

No delegates. The single ticket was written and implemented in the coordinating
session, which is what the right-size table authorizes and what the charge
directed. Its evidence is therefore self-catch evidence, weaker in kind than an
accepted-claim catch, and it should not be folded into any implementer aggregate.

The one-ticket slice was correctly sized: one production line, one test, one
stale expectation to update, all demoable together.

## Coordinator catches

Two, both self-caught before the landing gate.

The first was a seam error. The pinning test went into `internal/conformance`
because that package already execs bash scripts. It passed there, and the
mutation probe bit there, but the package's own meta-check red the gate:
`conformance meta unregistered live-tree assertion`. That package admits a
live-tree test only through its executable registry or a fixture-construction
classification, and a policy assertion on the live wrapper is neither. The honest
fix was to move it to `internal/systemtest`, beside the adoption journey that
already execs the real wrapper, not to add a classification entry that would have
been false.

The second was a single-source catch. `internal/systemtest/adoption_test.go:89`
had pinned `"HOME: unbound variable"` as an incidental consequence of dropping a
declared gate input. Copying the new wording there would have put the same string
in two packages. It now asserts `"HOME:"` and defers the exact message to the new
test, with a comment naming where the wording is pinned.

The charge also asked whether any canary expectation depends on the changed line.
None does: the four `tests/canary/*/files/bin/bench.sh` copies are five-to-seven
line stubs carrying no `HOME`, so the frozen fixtures stayed untouched.

## Repair attribution

| ticket | repair rounds | causes |
| --- | --- | --- |
| `name-the-missing-pool-home-input` | 1 | `other` |

The one round was the conformance-package meta red described above. It is not
`spec-row` (no spec), not `ticket-slicing` (the slice was right), and not
`tree-drift` (the tree did not move) — the coordinator chose a package without
first reading that package's admission rule.

## Agent-experience improvements

### Bench CLI

`bench commit --spec <slug>` exited `error: landed-but-checkout-incomplete: exit
status 1` on this tickets-only light-path spec. The commit itself was fully
correct: gate green, three real files landed, ticket folder excluded from history.
Only the working-tree deletion did not complete, leaving `specs/<slug>/` as
untracked residue. The friction is that a green landing reports a red-looking
error, which reads as a failed commit until you inspect `git log` — and the
session then removes the residue by hand and pays another full gate run.
Expected effect of a fix: the light-path close step stops costing a manual
recovery and a re-gate, roughly one minute and one false alarm per light-path
landing. The fix needs a reviewer call on whether the checkout is retried or the
residue is reported as a named next action, so it is filed rather than folded in.

### Skills

`craft-seams` should say that a package which enumerates or registers its own
tests carries a seam contract, not merely a location: read that admission rule
before placing a test there. The friction cost a full gate cycle here. Expected
effect: the placement question gets asked at write time, where it is free,
instead of at gate time.

### Process

The light path's own close sequence is worth naming. `bench commit --spec` gates,
commits, and retires the ticket folder in one call, so a light-path build reaches
`/bench-final-check` with nothing left to land and the phase is an honest no-op
plus this retro. That worked; the only rough edge is the CLI defect above.

The pre-build divergence is the other note. The git guard refused the rebase,
which is right — the reconciliation is the reviewer's. The session kept working
on the pre-rebase HEAD while surfacing the choice with a recommendation, and
nothing had to be redone. Treat that as the pattern: surface the reconciliation,
recommend one option, and continue on work that does not depend on it.
