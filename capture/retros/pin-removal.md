# Retro: pin-removal

## Outcome

The landing `97c98c3e` published the spec on 2026-09-06 over the source pair `14cc7127`
to `abb18d7c`. The managed pre-push hook keeps one clause: it blocks a direct push to
the protected branch, and `bench.allowProtectedPush` lifts that clause for one
repository. The drift clause, the `gate unpinned` warning, the `bench gate pin` verb,
its `gate-pin` plumbing route, and every pin reader left the kit. The gate grammar is
`bench gate [--fresh]`. The reference guide, ADR 0001, the FT141 detail, and the
distillation doc state the resulting rule. The spec stays `implemented` as the veto
surface until the reviewer retires it at the drain.

## Gate-stage timings

| stage | landing gate |
| --- | --- |
| gofmt | 100 ms |
| vet | 936 ms |
| test | 70130 ms |
| race | 2488 ms |
| system | 26122 ms |
| shellcheck | 537 ms |

Three fold gates ran before the landing: the `main` fold and the two sibling folds. Each
ran green on its first pass, with the test stage between 67376 ms and 69803 ms.

## Ticket-versus-spec-slice and delegate performance

Three ticket charges ran on `opus` at medium, one per ticket, in three worktrees. All
three landed first-pass on behavior, and every self-probe bit as the ticket named it. The
hook charge reported that the topic-silence test stays green under the drift probe and
explained why, instead of hiding the limit. The verb charge found that the package-local
routing test grades no real tree and moved its probe to `TestRootConformance` with the
root override. One repair charge on `opus` at medium landed first-pass. Three review
axes ran on `opus` at low and returned eight findings, of which two became repairs.

## Coordinator catches

- The Standards axis read the threat-model sentence in ADR 0001 as stale. The guidance
  ticket keeps that sentence by decision, so the finding was a no-op.
- The Standards axis read `internal/git` as no tree-hash owner. The package exports
  `TreeHash`, which the gate calls, so the finding was a no-op.
- `bench worktree merge --from <sha>` refused a sibling tip by sha; the fold takes the
  sibling label.
- `bench worktree land --base` takes the `main` tip, not the review's frozen fold
  commit, so the first landing call refused.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | --- | --- |
| remove-the-drift-clause-from-the-hook | 1 | delegate-error |
| retire-the-gate-pin-verb | 0 | none |
| state-the-hook-rule-in-the-guidance | 0 | none |

The one round folded a narrating comment and an exact-stderr assertion into one commit.

## Agent-experience improvements

### Bench CLI

- Let `bench worktree merge --from` accept a sibling's tip sha as well as its label,
  because the handoff pins the sha and the refusal names no label route. The census
  entry `pin-removal landing census: 3 raw calls` records the three raw shell reads.
  Feeds: new
- Let the `bench worktree land` refusal for a non-ancestor base name the
  distinction between the review's frozen base and the landing base in its detail.
  Feeds: new

### Skills

- Let `bench-review-implementation` state that the landing base is the `main` tip
  merged into the source, while the review base is the fold commit.
  Feeds: none

### Process

- A review finding that contradicts a ticket's explicit keep decision is a no-op, and
  the coordinator cites the ticket line in the disposition.
  Feeds: none
