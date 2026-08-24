# Retro — ft238-hygiene-batch

## Outcome

FT238 landed as `d747b9ec`: nine light-path tickets and one review repair on
one integration worktree. The landing contained:

- the absolute worktree path;
- the commit help example;
- `bench commit --dry-run`;
- the heredoc-body guard fix;
- the run-binary glossary term;
- the `bench diff --base/--source-tip` pair;
- the phase-close capture-batch rule;
- the persisted and printed request token;
- identifier and prefix operands for the worktree path verbs;
- a fingerprint-prefix `clean --apply` command.

The landing's own gate run is the retained
green evidence; the source pair is `3e2dee10` to `73e3c4e5`.

## Gate-stage timings

A full gate run held near 3 minutes across all eight runs. The test phase
dominates; `internal/worktree` alone holds near 57 seconds every run. Seven
commit gates plus one landing gate ran; two more red gate runs paid for
authored-prose bounds and a pinned envelope.

## Ticket-versus-spec-slice and delegate performance

The build ran inline in the main session under the light-path standing
approval, so no write delegates served. Three sonnet review axes served
read-only: Standards returned one accepted finding with a refuted-leads list,
Spec returned four accepted ungraded-row findings, and Coverage returned six
mutation candidates, five accepted and one folded into the Standards repair.
Every axis cited file and line, and each refuted its own weak leads.

## Coordinator catches

- The `--full` flag refuses on a missing spec; the roadmap row routes FT238 as
  light-path tickets, so the run proceeded on that standing approval.
- The capture-batch piece had gone stale: the inboxes are git-ignored since
  `a3e84401`, so the landed rule now names only the tracked capture set.
- The plaintext request token sits beside the digest in the intent ledger;
  the digest stays the authorization identity. Flagged for veto.
- The list-column change crossed two pinned surfaces the tickets never named:
  the checked-in terminal-pair fixtures and the AXI registry envelope.
- A mutation probe on the heredoc fix turned red at the strip call and green
  on restore before the ticket committed.
- The prefix relaxation initially weakened the reclaim fingerprint check.
  The reclaim test caught it. The repair split the predicates, so the
  destructive path stays strict.

## Repair attribution

| ticket | repair rounds | causes |
| --- | --- | --- |
| emit-absolute-worktree-path | 1 | other |
| add-commit-help-example | 0 | none |
| add-commit-dry-run | 0 | none |
| skip-heredoc-bodies-in-gitguard | 0 | none |
| add-run-binary-glossary-term | 0 | none |
| accept-source-tip-in-diff | 0 | none |
| batch-capture-set-at-phase-close | 1 | other |
| persist-and-print-request-token | 2 | ticket-slicing, ticket-slicing |
| accept-identifiers-in-worktree-verbs | 1 | ticket-slicing |

The two `other` rounds were one shared prose-bounds red in the authored ticket
files. The three `ticket-slicing` rounds were pinned surfaces the tickets did
not name: the list fixtures, the AXI envelope, and the strict reclaim
fingerprint test.

## Agent-experience improvements

### Bench CLI

- Land the untracking of both capture inboxes and declare
  `capture/session-handoff.md` in the landing allowance, because the landing
  refused twice on local capture state this phase.
  Feeds: new
- Make `bench test` refuse a bare `internal/<pkg>` argument or map it to the
  module path, because the std-library collision grades the wrong package.
  Feeds: new

### Skills

- None.
  Feeds: none

### Process

- A batch roadmap row that lists pinned output surfaces (fixtures, envelopes,
  registry markers) beside each piece would have cut all three ticket-slicing
  rounds.
  Feeds: none
