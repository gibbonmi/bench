# Retro: parallel-landings

## Outcome

Landed as `3fb24331` on `main` (2026-08-23) through `bench worktree land --base 1a135f1b --source-tip c04bd54c --spec parallel-landings`. The gate was green, the worktree released, and the spec published `Status: implemented`. Five tickets and two repair commits sat on one retained worktree.

The review ran two loops at opus / medium. Loop 1 found 15 advisory findings, 0 blocking, and 8 repair targets. Loop 2 was repair-scoped; it found one unsatisfied predicate (the handoff commit pin) and one long comment. A final prose commit fixed both. The post-merge tail is the first dogfood: it lands spec-less through the verb.

## Gate-stage timings

Landing gate: gofmt 0.07 s, vet 1.2 s, test 65.8 s, race 4.4 s, system 8.6 s, shellcheck 0.4 s. Eight full gate runs were paid: five ticket commits, one pickup commit, two repair commits. Two red runs were also paid: one on a handoff prose bound and one on the stale-executable seal at the destination.

## Ticket-versus-spec-slice and delegate performance

Every charge was one ticket on the shared integration worktree, serial. Four opus / medium charges built tickets 1–4. Two fable / high charges built ticket 5 (guidance) and the loop-1 repair. All six returned diff-ready first-pass with a biting self-probe and a judgement list.

The ticket-2 delegate chose a resume proof by folder absence instead of a skip. The ticket-3 delegate added the mode refusal for every verb. The ticket-4 delegate kept the destination commit exact in `next=`. No charge needed a second round.

## Coordinator catches

- A coordinator probe per ticket at a distinct site and kind bit every time. The sites were the range proof, the close classification, the one-side union branch, the merge-to-rebase swap, and the reference needle.
- The ticket-5 probe was silently green under the conformance package alone. Only the gate's entry test grades the live tree. Re-probing through the gate went red, so the evidence held.
- My own handoff rewrite twice broke the prose bounds (a 31-word sentence, a 7-sentence paragraph). The gate caught both.
- The landing refused on a stale `dist/bench` seal at the destination. It then refused on an undeclared ignored residue: `.bench/dist/`, a stale copy from the day before. The named rebuild and a preserve-then-move resolved both.
- The worktree request token was not retained across the session. The mismatch refusal names `bench worktree reauthorize`. The correct token was recovered by trying the creation second.

## Repair attribution

| ticket | repair rounds | cause per round |
|---|---|---|
| make-spec-optional-on-the-landing | 0 | none |
| close-a-tickets-only-folder-on-the-landing | 0 | none |
| union-merge-the-phase-owned-journals | 1 | spec-row |
| name-the-source-repair-in-the-conflict-refusal | 1 | spec-row |
| route-every-phase-through-a-worktree | 1 | other |

The one repair commit folded eight review targets. The row above charges each target to the ticket whose seam it touched. The guidance ticket's round was the handoff and debug prose. The journals' round was the containment assertion and the prefix row. The refusal ticket's round was the unsafe-spec placeholder test.

## Agent-experience improvements

### Bench CLI

- Record the worktree's opaque request token in the assignment record, so `bench worktree path` or `list` can print it back.
  Feeds: new
- Let the landing's stale-executable refusal run the sanctioned rebuild itself and re-run, as the FT169 fix intended.
  Feeds: FT169
- Name the undeclared ignored residue's remedy in the refusal: declare it in the build-output declaration, or remove it.
  Feeds: new

### Skills

- State in craft-delegate that a probe of a live-tree anchor or marker runs through the gate, not through the conformance package alone.
  Feeds: none

### Process

- Write the phase-boundary handoff in short sentences and paragraphs of six or fewer; the prose bounds grade the handoff.
  Feeds: none
- Retain the worktree request token in the handoff's State section while a build is in flight; the landing needs it.
  Feeds: none
