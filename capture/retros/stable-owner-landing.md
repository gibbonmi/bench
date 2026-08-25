# Retro — stable-owner-landing

## Outcome

The landing published `0aecd5fc` on `main` and flipped the spec to
`Status: implemented`. The reviewed pair was base `7eca97bf` and tip
`f2368a97`. The source carried four commits: two build tickets, one review
pickup, and one repair. All 17 coverage rows landed. Two rows carry a
narrowed proof that the review verified against the spec's Won't-handle
text: SOL10 is data-scoped, and SOL15's executable half sits at the gate
seam. Public landing now runs under the installed promotion broker, and the
install step activates a new broker.

## Gate-stage timings

The landing gate: gofmt 0.1 s, vet 1.6 s, test 84.0 s, race 5.1 s, system
9.7 s, shellcheck 0.6 s — about 101 s total. The build paid five full gate runs:
three ticket commits, one red prose retry, and one landing. The close paid
three more: the retirement, one red prose retry, and the capture commit.

## Ticket-versus-spec-slice and delegate performance

All charges were ticket-sized; no delegate received a whole spec slice.
Ticket 01 (Sol / high) landed in one build pass plus one continued pass for
a fence-extended test rewrite. It stopped at the out-of-fence red instead of
an edit past it. Ticket 02 (Terra / medium, a reviewer tier override) landed
with one coordinator-caught coverage gap. It caught its own wrong first
SOL13 probe.

Ticket 03 (Terra / high) closed six targets first-pass and
correctly stopped on the seventh. The accepted P1 predicate was
unimplementable against the rebase-forward flow, and its counter-proposal
became the reviewer's decision. Luna / medium ran three review axes plus the
repair-scoped re-review; the citation standard held on all passes.

## Coordinator catches

- The SOL10 unit tests set the baseline variable themselves, so the
  owner-to-gate transport was unproven. A coordinator omission probe at the
  transport site stayed silently green and forced the missing landing-seam
  test.
- The ticket-01 done-claim carried a genuine out-of-fence blocker: a stale
  registry test that pinned removed behavior. The red was verified before
  the reviewer disposed of it.
- A concurrent session rewrote the shared handoff mid-phase. The clash
  became a slot-ownership decision instead of a silent clobber.

## Repair attribution

| ticket | repair rounds | causes |
|---|---|---|
| 01-pin-the-stable-promotion-owner | 1 | ticket-slicing |
| 02-gate-and-publish-under-the-stable-owner | 1 | delegate-error |
| 03-review-repairs | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Give `bench handoff` per-assignment sections, or refuse a rewrite when the
  file's content belongs to a different live request.
  Feeds: new
- Give the tagged systemtest suite one wrapped invocation that supplies the
  run-binary and kit variables, so a probe does not fail on plumbing first.
  Feeds: none

### Skills

- State in `craft-delegate` that a pure test-addition charge names the
  production mutation each new test is proven against.
  Feeds: none

### Process

- Route `bench worktree land` from the primary checkout, because
  `worktree exec` sets a routing variable the stable-owner route refuses.
  Feeds: none
