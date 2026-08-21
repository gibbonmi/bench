# Claude model scorecard

Last incorporated landing: `light-path-wrapper-home-guard` (`31d64ac7`,
2026-08-21) — Opus/medium ran a zero-delegate light path end to end (one shell
line, one pinning test, one stale expectation), self-caught a wrong test package
against the gate and a would-be duplicated string, and filed the CLI defect its
own close step hit rather than folding it in.

Sixteen completed landings are now recorded. Fable has a Go-seam reviewer
sample, prose, and one full-lifecycle orchestration. Opus has a large
implementer sample across three effort tiers and a twelve-landing orchestrator
sample. Sonnet has six reviewer landings and a large implementer sample.
Routing follows the project's harness-to-tier binding.

## Current routing

| model / effort | role and sample | observed quality | current use |
| --- | --- | --- | --- |
| Fable / high, medium | orchestrator on 1 landing + implementer, 3 prose tickets + reviewer, 2 rounds on one Go spec (`worktree-exec-run-binary`) | every done-claim probed with a distinct mutation kind and site; caught three false enforcement claims in spec rows; as implementer it surfaced an unreachable acceptance criterion instead of forcing it; as reviewer on Go it refuted a ticket's own stated premise by citation, then found that the reviewer-chosen design would make an exec child inherit a seal-verified binary and reuse a stale artifact against the profile's private-exact-source rule — a blocking finding that reopened a closed fork and changed the design, traced through four files the charge had not named; on `land-spec-amendments` its guidance ticket landed the right sentence at the right site first-pass and justified the site better than the charge asked, distinguishing where sessions summarize from where they act, but returned a ragged mid-paragraph wrap that cost one formatting round; as orchestrator on `ft230` it ran a full spec lifecycle (2 tickets, review, repair, land, retire) with every delegate claim verified true and no false accept, but paid 10 full gate runs — three of them on markdown-only destination diffs — and recovered two landing refusals (dirty destination, divergent spec bytes) by hand | Top tier for guidance prose, for coordination when the spec's own claims are suspect, and now for spec review where a design decision has already been made and needs adversarial testing rather than confirmation |
| Opus / medium, low | implementer, 35 medium-effort charges across 10 landings + 2 low-effort polish-repair charges (`worktree-cleanup-eligibility`) | medium: 31 of 35 first-pass accepted, 4 coordinator catches; on `ft230` the workflow-swap charge self-reported its two gate-forced out-of-fence edits and the repair charge corrected its own ticket's "resume completes" sketch to the state machine's real `in_progress`/`promote` terminal after reading it; on `land-spec-amendments` all 3 charges landed first-pass, and the repair charge rejected the fixture shape its charge prescribed — showing that shape stayed green under the very mutation it was meant to catch — then derived the shape that was red-capable and covered two production sites instead of one; on `worktree-exec-run-binary` all 3 charges landed first-pass and one declined the classifier its charge named because that helper reads file content, which its ticket's single-predicate rule forbade, flagging the divergence rather than working around it (one described a preflight-red diff as clean; one omitted a test its own production change needed), every required mutation bit and production restored exactly; on `ft227` all 3 charges landed first-pass including a 7-leg tagged system journey against an owner harness read cold, and each self-reported a judgment call rather than burying it; low (n=8): 8 of 8 first-pass accepted at the diff level, 1 coordinator catch; on `ft228` six low-effort charges covered anchor rows and canary fixtures, a 13-row conformance policy table with four graded facts, a Claude-side parity check with two completeness directions, and three Go fixes in the worktree/wrapper seam — three of them returned something the charge had not named (a second instance of the same test bug, a second guard the linked-project case demanded, a refusal to fold a variable into a shared helper because a contract test required its inheritance) | Medium is mid-tier default for gate and conformance logic; low is now the default for tickets derived from an exact spec at a known seam under a covering gate, which is most tickets a good spec produces — the eight-charge sample spans conformance logic, canary fixtures, and Go bug fixes, not only polish |
| Opus / high | implementer, 14 charges (2 process-lifecycle, 4 guidance-prose/anchor, 3 inline light-path, 3 Go-seam parser/conformance-oracle/adapter rewrites, 1 conformance-correctness repair) | all 14 first-pass accepted; on `ft230` the adapter charge ran five attributed mutation probes and revived a dead tagged test harness as the single approved-set fixture source rather than pasting a second copy; on `ft234-pool-key-reclaim` the single high charge (the destructive apply) returned two independent removal guards and, unprompted, showed that its own acceptance row's specified seam was unreachable rather than writing a test that could not fail; the largest single charge on `spec-ticket-fence-reduction` (a full parser rewrite plus a 67-row repo migration with an empty differential proof) still self-probed cleanly and flagged a spec-vs-charge disagreement rather than silently resolving it either way | Use high for process-lifecycle, cleanup-authority, destructive-command, anchored guidance-prose, and large/foundational Go-seam rewrites |
| Opus / medium | orchestrator, 13 landings | on `light-path-wrapper-home-guard` (zero delegates) it placed the pinning test in the package that already execs bash and let the gate correct it, then moved the test rather than adding the false classification entry that would have silenced the meta-check; it also refused to copy the new refusal string into a second package, rewriting the incidental assertion there to defer to the one pinned copy; on `land-executable-freshness` it read the spec's own test seam before charging it and found LF4's row unbuildable — the row named a subprocess journey plus a package-var substitution, which cannot combine — then supplied the real failing condition in the charge, which is what bought a first-pass ticket; it also ran the one probe the delegate had skipped (the mutation LF4 exists to catch) and confirmed the row bites, and it kept a review finding it could have waved off, since six passing cases all produced a nil error and the predicate's `errors.Is` half was genuinely ungraded; verified every done-claim independently with a mutation at a distinct site/kind from the delegate's own each time; caught 1 silent oracle regression, 1 partial-ordering hole, 1 dependency inversion, 1 under-declared ownership fence, 1 cross-ticket merge conflict from parallel repair porting, 1 resolved review-pickup artifact left uncommitted-for-deletion that hard-blocked landing, a blocked landing whose opaque four-cause refusal it resolved from the intent ledger rather than by trial, and (these two landings) a destructive-command predicate that would have deleted a live worktree, caught by planting the shape rather than reading the code; it also probed its own two repair fixes and found both unpinned; on `ft229` it reproduced a delegate's reported fail-open independently before accepting it and separately confirmed the escaping was the producer default rather than contrived, refuted two further delegate claims by measurement, and caught two of its own instruments lying — a mutation probe whose red was a compile error, and an exit-code harness reading `$?` after a command substitution; on `worktree-exec-run-binary` it found its own composed-run evidence worthless because the gate had answered from cache before reaching the path under test, and re-ran it as a controlled A/B; on `ft227` it refuted both surviving review findings from cited spec text rather than opening repair tickets; on `ft228` it found an acceptance row graded on only one of its two halves by probing from the row's text rather than from the implementation, and refused to accept a delegate's red-then-green "flake" attribution, which separated a real one-in-a-thousand test flake from a structural defect that had been silently dropping two gate phases inside every worktree; on `land-spec-amendments` it accepted a review finding only after refuting its stated mechanism by probe — a mode mismatch reds rather than shipping silently — while keeping the coverage gap the finding had actually exposed, and returned a delegate's correct prose to its author rather than hand-editing it to stay inside the inline allowance | Continue; deliberate site/kind variance per probe is holding up, and reading the authoritative state record before theorizing is the same discipline applied to operational blocks |
| Sonnet / medium–high | implementer, 10 ticket-sized charges (`worktree-cleanup-eligibility`: 8 build tickets + 2 regression/seam repairs) | 10 of 10 first-pass accepted at the diff level, including that landing's two highest-risk charges (an ordered-decision-logic extraction and its final cross-file consolidation); one delegate caught and discarded its own wrong initial hypothesis via direct testing rather than reporting it as fact; one delegate flagged a judgment call already resolved in its charge rather than silently picking a different answer | Effort should scale with the seam's behavior-preservation risk, not the ticket's line count — the two highest-effort charges here were both refactors of already-correct logic, not new logic |
| Sonnet / low–medium | reviewer, 3 read-only axes on 6 landings + 4 scoped re-reviews | `land-executable-freshness` (low): 3 findings, 1 repair target — Coverage found the feature's applicability gate reached only through the land surface and named the exact ungraded states; Standards and Spec each returned zero after auditing all eight rows and the one-source claim against the constant's other readers, but both drifted on line numbers (a closure cited at line 94 sat at 555) while every substantive claim held; the scoped re-review then found the repair still ungraded on its second mutation — every case produced a nil error, so `err == nil` survived them all — which is the finding that closed the gap. `ft230` (medium): 3 findings, 1 de-duplicated repair target — Coverage caught that the spec's edge inventory decided two npm-failure behaviors no coverage row tested, and its third finding was correctly refuted by the spec's own Access decision; the scoped re-review verified the repair's resume assertions against the state machine's source. `ft228` (medium): 3 axes, 1 accepted finding — Standards caught a doc paragraph re-enumerating the phase set the policy table grades, and the enumeration was already wrong by one phase; Spec cleared all 24 stories and every implementation-decision bullet against the diff; Coverage verified all 22 rows by grepping needle uniqueness and reading check order, then instrumented the check to measure the suspected fixture-noise concern at 13 diagnostics and correctly refuted it as harmless by showing no fixture EXPECT collides with the noise text. `ft227` (low): 3 raw findings, 0 repair targets: Spec audited all 15 acceptance rows and cross-checked one asserted CLI string against its production source rather than trusting the test; Standards and Coverage each filed one finding a cited source refuted, and Coverage's worst finding misread an independently authored string expectation as weaker than a behavioral one, missing that the literal is itself the mutation catch `worktree-exec-run-binary`: 5 raw findings, 4 repair targets; Coverage found a disposition the spec had decided but never rowed, and named the exact plausible mutation (`os.Stat`→`os.Lstat`) that would silently invert it; Spec caught a coverage row claimed in a commit message but recorded nowhere in the tree; a scoped re-review of the repair returned zero findings after re-verifying every moved assertion under four production mutations. `land-spec-amendments` (medium): 1 finding across 3 axes; Standards and Spec each refuted every lead handed to them rather than manufacturing findings, and Spec audited all 13 rows against code and test rather than row-ID presence; Coverage's one finding named the wrong mechanism but the right gap, and its scoped re-review then verified the two-site repair claim empirically with scratch `merge-tree` experiments instead of accepting it. | Keep at reviewer; three landings in, citation discipline holds and no axis has needed a re-run. Medium effort bought measurement over speculation — the Coverage axis resolved an open concern with instrumented evidence instead of a verdict — so prefer medium where the charge names a specific concern to settle |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft228-debug-restoration` IP10 | Opus / medium / orchestrator | A coverage row required an adapter description to carry a trigger *and* omit an anti-trigger; the check graded only the omission, and the delegate's own probe — derived from what it had built — could not see the gap. A coordinator probe written from the row's text instead set the description to prose carrying neither, and the gate stayed green: the exact silent failure the row names. |
| `land-spec-amendments` mode repair | Opus / medium / implementer | Handed a prescribed fixture shape, it showed that shape passed under the mutation it was written to catch — a one-sided mode change auto-resolves in `merge-tree` — and derived the add/add shape that actually conflicts, returning two rows covering two production sites rather than the one row charged. |
| `light-path-wrapper-home-guard` | Opus / medium / orchestrator | A test that passed, and whose mutation probe bit, still red the gate: `internal/conformance` admits a live-tree test only through its own registry, and a policy assertion on the live wrapper qualifies under neither admission route. Picking the package by "which one already execs bash" cost a full gate cycle; the package's admission rule was the fact to read first. |
| `worktree-exec-run-binary` WX20 | Opus / medium / orchestrator | Re-ran a delegate's composed-run evidence and got `gate: green (fresh verdict reused for this tree)` twice. The cache answers before the gate entry is reached, so neither run exercised the path the row grades — the evidence was worthless and had already been reported as verified. `--fresh` then produced a real controlled A/B: the unfixed driver refuses, the fixed driver runs all six phases. |

## Current decisions

- A delegate's reading of a gate skip is a claim: plant a break and run the oracle.
- Change routing only after two comparable runs or one controlled model comparison.
- **Review axes now route to Sonnet** at the reviewer's direction (2026-08-19). The Opus
  reviewer row above is the closing sample for that role, not the current routing; the
  two Sonnet reviewer samples (`ft227`, `ft228`) held the citation standard at 3 axes
  each with no re-run; keep the routing and watch whether it under-reads
  string-expectation seams. Write and orchestrator routing is unchanged.
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
- A gate verdict carrying "reused for this tree" is evidence about the tree's last real
  run, not about any code path. An evidence run aimed at a specific path uses `--fresh`;
  the provenance prints on the same line as the verdict and reads as a pass.
- The `ft230` amendment refusal was never a contract disagreement: FT225 had retired that
  refusal hours earlier, and a `dist/bench` built the day before was still enforcing it.
  A stale executable enforces the contract it was built with, and the land command refuses
  before any gate runs, so nothing downstream could catch it. `land-executable-freshness`
  closed the residual; the surviving lesson is that a refusal quoting a rule you believe is
  retired is evidence about the binary, not about the rule.
- A delegate citing a file and line it did not re-read produces a claim that looks
  checkable and is not. Three axes here held every substantive claim while two drifted on
  line numbers. Weigh the claim, verify the location.
- A charge that prescribes a fixture *shape* can buy a vacuous test. Name the property
  that must be red-capable and let the delegate derive the shape: here the prescribed
  shape stayed green under the exact mutation it was meant to catch, and only the
  delegate's own derived shape bit.
- Assert a repo convention from the check that enforces it, not from a sample of artifacts
  that happen to follow it. Reading `specs/` and finding no coverage-row citations
  produced a confident, wrong report; the rule lived in a preflight check.
