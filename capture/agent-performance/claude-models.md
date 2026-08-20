# Claude model scorecard

Last incorporated landing: `ft229-hygiene-batch` (`c8a3fad2`, 2026-08-20) —
Opus/medium orchestrated the resume of a build whose ten tickets had already
landed on a retained source, ran the eleventh inline, and took the diff through
review, repair, and landing. Review axes ran at Opus/medium at the reviewer's
explicit direction for this diff, not at the standing Sonnet routing. Two
fail-opens at the enforcement boundary were found by review after the build
reported complete; the gate could not have found either.

Eleven completed landings are now recorded. Sonnet now carries a large,
higher-stakes implementer sample spanning medium and high effort on real
production Go refactor work, not only low-effort tickets. Opus has a large
implementer sample across three effort tiers and a nine-landing orchestrator
sample, and its low-effort implementer sample is now large enough to route from.
Its six-landing reviewer sample is closed in favor of Sonnet, which has two.
Routing still follows the project's harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high, medium | orchestrator on 1 landing (tickets 6–8, review, landing) + implementer, 2 prose tickets (1 high, 1 medium) | every done-claim probed with a distinct mutation kind and site; caught three false enforcement claims in spec rows and two unattributed-writer events; as implementer it surfaced an unreachable acceptance criterion instead of forcing it; at medium effort on `ft228` it restored a 167-line phase file against a ~165-line estimate with every one of nine future anchor needles unique on first write, and reported its own mutation probe's silent green plainly as the expected pre-pin baseline rather than dressing it as a pass | Top tier for guidance prose and for coordination when the spec's own claims are suspect; medium effort was sufficient for an exacting prose restoration against a quoted source, so reserve high for prose whose target text is not already pinned; not yet observed on Go |
| Opus / medium, low | implementer, 27 medium-effort charges across 7 landings + 2 low-effort polish-repair charges (`worktree-cleanup-eligibility`) | medium: 23 of 27 first-pass accepted, 4 coordinator catches (one described a preflight-red diff as clean; one omitted a test its own production change needed), every required mutation bit and production restored exactly; on `ft227` all 3 charges landed first-pass including a 7-leg tagged system journey against an owner harness read cold, and each self-reported a judgment call rather than burying it; low (n=8): 8 of 8 first-pass accepted at the diff level, 1 coordinator catch; on `ft228` six low-effort charges covered anchor rows and canary fixtures, a 13-row conformance policy table with four graded facts, a Claude-side parity check with two completeness directions, and three Go fixes in the worktree/wrapper seam — three of them returned something the charge had not named (a second instance of the same test bug, a second guard the linked-project case demanded, a refusal to fold a variable into a shared helper because a contract test required its inheritance) | Medium is mid-tier default for gate and conformance logic; low is now the default for tickets derived from an exact spec at a known seam under a covering gate, which is most tickets a good spec produces — the eight-charge sample spans conformance logic, canary fixtures, and Go bug fixes, not only polish |
| Opus / high | implementer, 13 charges (2 process-lifecycle, 4 guidance-prose/anchor, 3 inline light-path, 2 Go-seam parser/conformance-oracle rewrites, 1 conformance-correctness repair) | all 13 first-pass accepted; on `ft234-pool-key-reclaim` the single high charge (the destructive apply) returned two independent removal guards and, unprompted, showed that its own acceptance row's specified seam was unreachable rather than writing a test that could not fail; the largest single charge on `spec-ticket-fence-reduction` (a full parser rewrite plus a 67-row repo migration with an empty differential proof) still self-probed cleanly and flagged a spec-vs-charge disagreement rather than silently resolving it either way | Use high for process-lifecycle, cleanup-authority, destructive-command, anchored guidance-prose, and large/foundational Go-seam rewrites |
| Opus / medium | orchestrator, 9 landings | verified every done-claim independently with a mutation at a distinct site/kind from the delegate's own each time; caught 1 silent oracle regression, 1 partial-ordering hole, 1 dependency inversion, 1 under-declared ownership fence, 1 cross-ticket merge conflict from parallel repair porting, 1 resolved review-pickup artifact left uncommitted-for-deletion that hard-blocked landing, a blocked landing whose opaque four-cause refusal it resolved from the intent ledger rather than by trial, and (these two landings) a destructive-command predicate that would have deleted a live worktree, caught by planting the shape rather than reading the code; it also probed its own two repair fixes and found both unpinned; on `ft229` it reproduced a delegate's reported fail-open independently before accepting it and separately confirmed the escaping was the producer default rather than contrived, refuted two further delegate claims by measurement, and caught two of its own instruments lying — a mutation probe whose red was a compile error, and an exit-code harness reading `$?` after a command substitution; on `ft227` it refuted both surviving review findings from cited spec text rather than opening repair tickets; on `ft228` it found an acceptance row graded on only one of its two halves by probing from the row's text rather than from the implementation, and refused to accept a delegate's red-then-green "flake" attribution, which separated a real one-in-a-thousand test flake from a structural defect that had been silently dropping two gate phases inside every worktree | Continue; deliberate site/kind variance per probe is holding up, and reading the authoritative state record before theorizing is the same discipline applied to operational blocks |
| Opus / medium–high | reviewer, 3 read-only axes on 6 landings + 1 delta re-review + 1 scoped low-effort follow-up | 21+19+12(3-axis)+7(follow-up) raw findings across landings; every finding cites file:line or story/row ID; `worktree-cleanup-eligibility`'s Coverage axis found a real production regression (a derived-after operator override that had migrated outside its original guard, reachable but exercised by no test or coordinator mutation probe) and Standards+Spec independently converged on the same under-closed spec row from two different angles; on `ft234` the Coverage axis found the relative-`gitdir:` deletion path that three write delegates, seven mutation probes, and a green gate all missed, while the Spec axis was wrong on one duplication claim the Standards axis got right | Reopened at the reviewer's direction for `ft229`, whose diff carried the git guard and the gate: 3 axes returned 21 raw findings and 10 repair targets, and the Coverage axis found a reproduced fail-open in the degraded rim that a green gate, eleven tickets, and every delegate probe had passed over. A single combined-axis delegate over the 12-file repair diff then found three more, one of them a second fail-open. Reviewer default at mid tier for a first full pass; a narrowly-scoped low-effort follow-up (verify specific named fixes, not re-hunt the whole diff) held up as a real second-pass discount, not just a smaller sample |
| Sonnet / medium–high | implementer, 10 ticket-sized charges (`worktree-cleanup-eligibility`: 8 build tickets + 2 regression/seam repairs) | 10 of 10 first-pass accepted at the diff level, including that landing's two highest-risk charges (an ordered-decision-logic extraction and its final cross-file consolidation); one delegate caught and discarded its own wrong initial hypothesis via direct testing rather than reporting it as fact; one delegate flagged a judgment call already resolved in its charge rather than silently picking a different answer | Effort should scale with the seam's behavior-preservation risk, not the ticket's line count — the two highest-effort charges here were both refactors of already-correct logic, not new logic |
| Sonnet / low–medium | reviewer, 3 read-only axes on 2 landings | `ft228` (medium): 3 axes, 1 accepted finding — Standards caught a doc paragraph re-enumerating the phase set the policy table grades, and the enumeration was already wrong by one phase; Spec cleared all 24 stories and every implementation-decision bullet against the diff; Coverage verified all 22 rows by grepping needle uniqueness and reading check order, then instrumented the check to measure the suspected fixture-noise concern at 13 diagnostics and correctly refuted it as harmless by showing no fixture EXPECT collides with the noise text. `ft227` (low): 3 raw findings, 0 repair targets: Spec audited all 15 acceptance rows and cross-checked one asserted CLI string against its production source rather than trusting the test; Standards and Coverage each filed one finding a cited source refuted, and Coverage's worst finding misread an independently authored string expectation as weaker than a behavioral one, missing that the literal is itself the mutation catch | Keep at reviewer; two landings in, citation discipline holds and no axis has needed a re-run. Medium effort bought measurement over speculation — the Coverage axis resolved an open concern with instrumented evidence instead of a verdict — so prefer medium where the charge names a specific concern to settle |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required an adapter description to carry a trigger *and* omit an anti-trigger; the check graded only the omission, and the delegate's own probe — derived from what it had built — could not see the gap. A coordinator probe written from the row's text instead set the description to prose carrying neither, and the gate stayed green: the exact silent failure the row names. |
| `spec-ticket-fence-reduction` realign ticket | Opus / high / implementer | Ran the charged omission probe on its own landed glossary entry, saw the gate stay green, and reported that plainly as a finding instead of a pass — which exposed three spec rows claiming enforcement that did not exist and produced two repair tickets. |
| `ft229-hygiene-batch` Coverage axis | Opus / medium / reviewer | Found that narrowing the degraded git guard from a raw-substring test to a command parser had opened a fail-open: the new JSON decoder collapsed every `\uXXXX` escape to a placeholder, so an escaped `&&` stopped separating and the destructive half of the command left command position. The surface had zero assertions because the spec's rows were inherited from the substring rim being replaced, so no gate rigor could have reached it. |
| `bench-front-door` landing repair | Opus / medium / orchestrator | A landing refusal reading `request, assignment, or path mismatch` names four distinct causes; dumping the intent ledger falsified three outright and one digest comparison confirmed the fourth (the assignment ID had been passed where the request token belonged), so the repair was one sanctioned reauthorize call with no commit rewritten and no source implementation touched. |
| `worktree-cleanup-eligibility` Coverage axis | Opus / medium / reviewer | Found that a refactor-relocated operator override (`--discard-branch`) had silently escaped its original detached-HEAD guard — a real, reachable regression that 8 build tickets' own tests, 8 coordinator mutation probes at deliberately varied sites, and a full green gate had all missed, because nothing had ever exercised that specific evidence combination. |

## Current decisions

- Fable now has one coordinator/implementer sample on guidance prose; keep it there and
  collect a Go-seam sample before widening.
- A delegate's reading of a gate skip is a claim: plant a break and run the oracle.
- Change routing only after two comparable runs or one controlled model comparison.
- **Review axes now route to Sonnet** at the reviewer's direction (2026-08-19). The Opus
  reviewer row above is the closing sample for that role, not the current routing; the
  two Sonnet reviewer samples (`ft227`, `ft228`) held the citation standard at 3 axes
  each with no re-run; keep the routing and watch whether it under-reads
  string-expectation seams. Write and orchestrator routing is unchanged.
- Opus/medium has a seven-landing implementer sample and a nine-landing orchestrator
  sample; it is the settled mid-tier default, not provisional.
- Route a ticket from the decision table at charge time, not from its story's line. On
  `ft228` every ticket derived from mid-routed stories ran cleanly at low effort, because
  an exact spec, a known seam, and a covering gate is the cheap row — which is what a good
  spec produces. The ceiling-not-binding rule is worth a real discount, not a formality.
- A red is attributed or it is not resolved. A delegate reporting red-then-green as a
  flake is making an attribution claim; here that claim concealed both a genuine
  one-in-a-thousand test flake and a structural defect dropping two phases from every
  worktree gate run. A second green run establishes nothing about the first red.
- A self-probe derived from the implementation cannot find a missing half of a
  requirement. At least one probe per acceptance row must be written from the row's own
  text, and the coordinator is the one positioned to write it.
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
- When production text is a generated script, an independently authored string expectation
  is the mutation catch, not a weaker substitute for a behavioral test. One review axis
  discounted exactly that and filed a finding its own cited literal already refuted.
- Rows inherited from the surface a ticket replaces cannot grade the surface replacing
  it. Narrowing a blunt check into a precise one creates parse surface the blunt check
  never had; on `ft229` that shipped two fail-opens at the enforcement boundary under a
  green gate, and seven of ten repair rounds in that build were `spec-row`. A ticket that
  replaces a check adds a row for the new surface, at spec time.
- A mutation probe that reds proves nothing until the red is attributed. Two probes here
  reported a kill that was a compile error or a harness bug, and in one of them the
  delegate the coordinator was overruling had been right.
- A stale session handoff is corrected by the tree, not trusted. `ft229` opened with a
  handoff claiming nothing was implemented while ten tickets had landed on a retained
  source; the commit-age row is what made the disagreement visible.
- Serial delegates on one retained integration worktree cost wall-clock and buy one author
  per diff. Disjoint fence files across tickets are not a licence to write concurrently in
  one tree; the lever is separate worktrees.
