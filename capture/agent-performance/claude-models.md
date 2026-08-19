# Claude model scorecard

Last incorporated landing: `ft234-pool-key-reclaim` (`50495a81`, 2026-08-19) —
Opus/medium orchestrated two spec-backed landings back to back (FT226 then FT234)
with six write delegates and six review axes; the review's worst finding was a
relative-`gitdir:` path that would have deleted a live worktree, and the
coordinator confirmed it by planting the shape rather than by reading.

Nine completed landings are now recorded. Sonnet now carries a large,
higher-stakes implementer sample spanning medium and high effort on real
production Go refactor work, not only low-effort tickets. Opus has a large
implementer sample across three effort tiers and a five-landing
orchestrator sample alongside a four-landing reviewer sample. Routing still follows the project's
harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high | orchestrator on 1 landing (tickets 6–8, review, landing) + implementer, 1 prose ticket | every done-claim probed with a distinct mutation kind and site; caught three false enforcement claims in spec rows and two unattributed-writer events; as implementer it surfaced an unreachable acceptance criterion instead of forcing it | Top tier for guidance prose and for coordination when the spec's own claims are suspect; not yet observed on Go |
| Opus / medium, low | implementer, 24 medium-effort charges across 6 landings + 2 low-effort polish-repair charges (`worktree-cleanup-eligibility`) | medium: 20 of 24 first-pass accepted, 4 coordinator catches (one described a preflight-red diff as clean; one omitted a test its own production change needed), every required mutation bit and production restored exactly; low (n=2): 2 of 2 first-pass accepted, one delegate exceeded its charge with a more diagnostic mutation probe than requested | Medium is mid-tier default for gate and conformance logic; low is promising for small, low-risk polish/hardening repairs specifically — too small a sample yet to set a default there |
| Opus / high | implementer, 13 charges (2 process-lifecycle, 4 guidance-prose/anchor, 3 inline light-path, 2 Go-seam parser/conformance-oracle rewrites, 1 conformance-correctness repair) | all 13 first-pass accepted; on `ft234-pool-key-reclaim` the single high charge (the destructive apply) returned two independent removal guards and, unprompted, showed that its own acceptance row's specified seam was unreachable rather than writing a test that could not fail; the largest single charge on `spec-ticket-fence-reduction` (a full parser rewrite plus a 67-row repo migration with an empty differential proof) still self-probed cleanly and flagged a spec-vs-charge disagreement rather than silently resolving it either way | Use high for process-lifecycle, cleanup-authority, destructive-command, anchored guidance-prose, and large/foundational Go-seam rewrites |
| Opus / medium | orchestrator, 7 landings | verified every done-claim independently with a mutation at a distinct site/kind from the delegate's own each time; caught 1 silent oracle regression, 1 partial-ordering hole, 1 dependency inversion, 1 under-declared ownership fence, 1 cross-ticket merge conflict from parallel repair porting, 1 resolved review-pickup artifact left uncommitted-for-deletion that hard-blocked landing, a blocked landing whose opaque four-cause refusal it resolved from the intent ledger rather than by trial, and (these two landings) a destructive-command predicate that would have deleted a live worktree, caught by planting the shape rather than reading the code; it also probed its own two repair fixes and found both unpinned | Continue; deliberate site/kind variance per probe is holding up, and reading the authoritative state record before theorizing is the same discipline applied to operational blocks |
| Opus / medium–high | reviewer, 3 read-only axes on 6 landings + 1 delta re-review + 1 scoped low-effort follow-up | 21+19+12(3-axis)+7(follow-up) raw findings across landings; every finding cites file:line or story/row ID; `worktree-cleanup-eligibility`'s Coverage axis found a real production regression (a derived-after operator override that had migrated outside its original guard, reachable but exercised by no test or coordinator mutation probe) and Standards+Spec independently converged on the same under-closed spec row from two different angles; on `ft234` the Coverage axis found the relative-`gitdir:` deletion path that three write delegates, seven mutation probes, and a green gate all missed, while the Spec axis was wrong on one duplication claim the Standards axis got right | Reviewer default at mid tier for a first full pass; a narrowly-scoped low-effort follow-up (verify specific named fixes, not re-hunt the whole diff) held up as a real second-pass discount, not just a smaller sample |
| Sonnet / medium–high | implementer, 10 ticket-sized charges (`worktree-cleanup-eligibility`: 8 build tickets + 2 regression/seam repairs) | 10 of 10 first-pass accepted at the diff level, including that landing's two highest-risk charges (an ordered-decision-logic extraction and its final cross-file consolidation); one delegate caught and discarded its own wrong initial hypothesis via direct testing rather than reporting it as fact; one delegate flagged a judgment call already resolved in its charge rather than silently picking a different answer | Effort should scale with the seam's behavior-preservation risk, not the ticket's line count — the two highest-effort charges here were both refactors of already-correct logic, not new logic |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `skills-index-hostile-input-hardening` HI10 | Opus / high / implementer | Met a non-negotiable protocol exactly — production-reached barrier, inherited-pipe marker, no sleep, poll, or temp-name oracle — and routed every bound through the existing deadline helper unprompted, so a conformance check that grades exactly that had nothing to flag. |
| `spec-ticket-fence-reduction` realign ticket | Opus / high / implementer | Ran the charged omission probe on its own landed glossary entry, saw the gate stay green, and reported that plainly as a finding instead of a pass — which exposed three spec rows claiming enforcement that did not exist and produced two repair tickets. |
| `spec-ticket-fence-reduction` collapse ticket | Fable / high / implementer | Enumerated all 47 anchored needles, then stopped at a spec-level shortfall (60-line budget unreachable from the retained needles) with the arithmetic and three options rather than forcing a fit or silently re-pricing. |
| `bench-front-door` landing repair | Opus / medium / orchestrator | A landing refusal reading `request, assignment, or path mismatch` names four distinct causes; dumping the intent ledger falsified three outright and one digest comparison confirmed the fourth (the assignment ID had been passed where the request token belonged), so the repair was one sanctioned reauthorize call with no commit rewritten and no source implementation touched. |
| `worktree-cleanup-eligibility` Coverage axis | Opus / medium / reviewer | Found that a refactor-relocated operator override (`--discard-branch`) had silently escaped its original detached-HEAD guard — a real, reachable regression that 8 build tickets' own tests, 8 coordinator mutation probes at deliberately varied sites, and a full green gate had all missed, because nothing had ever exercised that specific evidence combination. |

## Current decisions

- Fable now has one coordinator/implementer sample on guidance prose; keep it there and
  collect a Go-seam sample before widening.
- A delegate's reading of a gate skip is a claim: plant a break and run the oracle.
- Change routing only after two comparable runs or one controlled model comparison.
- **Review axes now route to Sonnet** at the reviewer's direction (2026-08-19). The Opus
  reviewer row above is the closing sample for that role, not the current routing; the
  next landing supplies the first Sonnet reviewer evidence. Write and orchestrator
  routing is unchanged.
- Opus/medium now has a four-landing implementer sample and a five-landing orchestrator sample; keep it as the mid-tier default
  and stop treating it as provisional.
- Opus/high is the routing for process-lifecycle, signal, cleanup-authority, and
  large/foundational Go-seam work until a medium-effort sample contradicts it.
- Charges should name every state reaching a branch, not only the state the acceptance
  row tests — both observed repair rounds share that shape.
- A parallel-repair-ticket batch (multiple delegates from one shared stale base, ported
  onto the retained source serially) is real leverage but shifts merge-conflict cost onto
  the coordinator's port step; a charge for one should declare its expected touched-file
  overlap with sibling charges so ports can be sequenced deliberately.
- The zero-delegate light-path shape produces self-catches, not accepted-claim catches;
  treat its evidence as weaker in kind, not merely smaller, when comparing to delegated
  landings.
- Sonnet at medium-high effort handled an 8-ticket behavior-preserving Go refactor,
  including its two highest-risk seams, with zero coordinator-found defects; effort choice
  tracked behavior-preservation risk correctly even though every ticket's line count was
  modest. Treat Sonnet as viable for real production refactor work at the right effort,
  not only for known-shape low-effort tickets.
- Mutation probes and a full green gate both missed a real regression that an adversarial
  review pass caught by asking "what combination of evidence has nothing exercised" rather
  than "does the changed logic behave as charged." A review pass earns its cost distinctly
  from probing: probing verifies a claimed fix; review hunts for the untested combination
  nobody claimed to have covered.
- A resolved review-pickup artifact (`reviews/<slug>.md`) left in the tree past its last
  finding's fix is now a hard `bench worktree land` refusal, not merely a style miss —
  delete it in the same commit that closes the last finding, every time.
- An operational refusal is diagnosable from the state record the tool itself reads. A
  message that conflates several causes is a CLI defect to file, not a reason to guess:
  dump the authoritative record, falsify causes against it, and confirm the survivor.
- Drive Bench through its own wrapper, never `go run ./cmd/bench`. The wrapper supplies
  the run-binary variable the gate requires; without it a full gate run reports green
  phases and then refuses authorization as `infrastructure`.
