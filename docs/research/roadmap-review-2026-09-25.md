# Roadmap review, 2026-09-25

Consolidate rows by the outcome they own. Remove shipped claims before adding scope.
Keep FT336 first. Promote review-record automation, test determinism, and recurring correctness failures.

Scope: all 82 owners in ROADMAP.md and roadmap/, with targeted source, spec, and history checks
Evidence status: source inspection; runtime closure remains unverified
Source repository: /home/mgibs/workspace/bench
Source branch: main
Source commit: 673ebcf8924bb2c91f78c8c40c98692eb08531b9
Report worktree: roadmap-review-20260925
Report branch: bench/assign/36b972233890af0e3ae6e436a75c8a72/5cade63c1ef354a5f15c0f43890cfd3d
Consumed by: the reviewer's next roadmap maintenance decision
Drift: refresh a finding when its row, implementation owner, or governing decision changes
Retire when: the reviewed roadmap records each accepted disposition
Applied by: `4ee3c625ef9ed4284808ff7453c85e4be78c425e`, which consolidated the board from 82 rows to 69
Link base: row links resolve at the source commit; a folded or retired row has no file on the current `main`

This report proposes dispositions. It changes no roadmap row, approved spec, or closed decision.
Both research delegates read the primary checkout. The coordinator checked the sources behind the recommendations.

## Question graph

1. Which rows share one outcome and can share an owner?
2. Which claims have shipped, lost their premise, or lack a current use?
3. Which work should move first, given recurrence, readiness, and dependencies?

The first two answers constrain the proposed order.
The row ledger below gives each current owner a disposition.

## Question 1: Which rows can share an owner?

These are proposed board consolidations. Each implementation can retain separate tickets and review chunks.

| Retained owner | Fold into it | Outcome and boundary |
| --- | --- | --- |
| [FT305](../../roadmap/FT305.md:3) | [FT296](../../roadmap/FT296.md:3) | Durable execution with explicit supervisor and worker lifetimes. Keep platform behavior as a shaping decision. |
| [FT318](../../roadmap/FT318.md:8) | [FT241's record mechanics](../../roadmap/FT241.md:14); optionally [FT317](../../roadmap/FT317.md:8) as a decision stage | Valid review records, plan identities, and amendments. Decide capability-blocked semantics before the writer encodes them. Keep FT241's end-to-end promises separate. |
| [FT338](../../roadmap/FT338.md:3) | [FT99's semantic premise checks](../../roadmap/FT99.md:3) | Verify current claims, affected readers, competing writers, and retired claims before review. Move mechanical fence checks to FT293. |
| [FT293](../../roadmap/FT293.md:17) | [FT262](../../roadmap/FT262.md:3) | Complete and explain ticket ownership before a charge. Retain uncited-row progress and structured closure output. |
| [FT172](../../roadmap/FT172.md:11) | [FT292](../../roadmap/FT292.md:3) | One roadmap identity and transition contract. The writer records decisions that the reviewer already made. |
| [FT265](../../roadmap/FT265.md:7) | [FT249's reader gap](../../roadmap/FT249.md:3) | Consistent capture reads from every checkout and one drain snapshot. Do not assume a new Git-ref store is necessary. |
| [FT326](../../roadmap/FT326.md:7) | [FT325](../../roadmap/FT325.md:3) | Anchor boundaries and canary mutations agree despite harmless line wraps. Preserve the shipped distinctness checks. |
| [FT341](../../roadmap/FT341.md:3) | [FT254's wrong-tree requirements](../../roadmap/FT254.md:39); [FT339's registry attribution](../../roadmap/FT339.md:3) | Each tree-scoped command reads and identifies its selected tree. Keep exec stdin, timeout, and script behavior with FT254. |
| [FT125](../../roadmap/FT125.md:3) | [FT254's bounded file reads](../../roadmap/FT254.md:77); [FT339's directory query](../../roadmap/FT339.md:5) | Task-shaped reader projections. Reuse FT336's response budget and FT341's target contract. |

FT307 can also absorb the structure-query requests buried in FT243 and FT6.
The evidence already names the same working-tree and scoped-headroom outcome.
[Sources: FT307](../../roadmap/FT307.md:3), [FT243](../../roadmap/FT243.md:37), [FT6](../../roadmap/FT6.md:40).

Keep these boundaries:

- FT336 bounds output; FT337 reduces repeated review evidence. Both remain separate delivery owners.
- FT115 removes timing and environment defects; FT255 controls concurrent resource use; FT332 detects future flakes.
- FT333 prepares tickets-only evidence. A review-record writer does not deliver that capability.
- FT130 governs capture writes during a gate. FT265 governs capture visibility and drain snapshots.
- FT283 controls landing transitions. FT284 controls retirement. Their lifecycle relationship does not require one implementation.
- FT306 remains the qualification and adoption milestone. Keep its security prerequisites visible.

Do not expand FT89, FT140, or FT254 into larger catch-all rows.
FT140 needs individual dispositions. FT254 needs the responsibility transfers above before new work.

## Question 2: Which claims are stale or no longer useful?

### Confirmed source facts

| Claim in the roadmap | Current evidence | Proposed disposition |
| --- | --- | --- |
| FT82 must close before reassessment. | Commit 2253e50d3af571722381150fcfce535aed7db513 retires release-preflight and removes FT82. [Current stale clause](../../ROADMAP.md:203). | Remove the stale closure condition. Preserve the immutable-release qualification requirements. |
| FT340 still needs a repair-ticket decision. | [Review guidance](../../.agents/commands/bench-review-implementation.md:50) updates each affected ticket's Covers line. [Anchors](../../internal/anchors/registry_data.go:444) enforce that form. | Drop that decision. Fold remaining drain lane wording into FT89. |
| FT106's research destination is untrackable. | [.gitignore](../../.gitignore:31) ignores only root research. [Research destination rules](../../.agents/skills/bench-craft-research/SKILL.md:83) distinguish map assets from repository research. | Remove the research-home clause. Keep semantic document verification. |
| FT89 treats commands --brief as an incomplete inventory. | [AGENTS.md](../../AGENTS.md:113) defines a liveness probe. [Implementation](../../cmd/bench/main.go:188) returns that probe. | Remove the inventory premise. Keep only current guidance contradictions. |
| FT293 still proposes closure paths, ordering edges, and new-file markers. | [Proposal output](../../internal/preflight/proposal.go:40) supplies paths and ordering. [New marker](../../internal/preflight/decision.go:361) already exists. | Narrow to residual owner closure, deleted paths, and proposal adoption policy. |
| FT98 needs a new preserve-and-restore probe primitive. | [Probe implementation](../../internal/probe/probe.go:149) preserves, mutates, executes, and restores. | Remove the duplicate probe capability. Preserve archive and cleanup requirements; FT168 owns probe extensions. |
| FT200 needs a landing authorization chokepoint. | [Landing](../../internal/worktree/land_identity.go:252) calls [preflight authorization](../../internal/preflight/gather.go:193). It checks authorized paths, not every preflight predicate. | Retain only complete-preflight and claimed-spec discovery policy. |
| FT265 lacks immutable drain ownership. | [Context](../../internal/roadmap/context_parse.go:124) reads a sealed generation. [Selected-row output](../../internal/roadmap/context_render.go:65) still omits capture bodies. | Keep the remaining retrieval and reader-consistency gap. Do not rebuild snapshot ownership. |
| FT108 says no refactor route exists. | [Deepening handoff](../../.agents/commands/bench-deepen.md:95) uses the light or spec workflow. | Retire the standalone route proposal. Fold any accepted characterization rule into existing craft guidance. |

### Recommendations, with limits

Archive FT38's visual revival unless the reviewer wants it.
Its own row describes optional aesthetics after a functional dashboard shipped.
[FT38](../../roadmap/FT38.md:3).

FT324 has no current in-tree citation for its alleged competing project rule.
Treat it as a retirement candidate unless the reviewer supplies that source.
Private harness memory was outside this review.
[FT324](../../roadmap/FT324.md:3), [current memory rule](../../.bench/BENCH.md:157).

FT261 needs a fresh reproducer before more design work.
Explicit-source facts now include untracked paths, but this review did not run the original command.
Park the old incident rather than claim a runtime fix.
[FT261](../../roadmap/FT261.md:3), [current path contract](../../internal/preflight/decision.go:55).

FT246 is a closure candidate after a focused constructor census.
The current run-binary policy supplies one exact-source executable to ordinary children.
Specialized bootstrap and publication tests still call Build because construction is their subject.
Do not replace those tests with a selected binary.
[Current policy](../../projects/benchkit.md:638), [bootstrap fixture](../../internal/runbinary/bootstrap_test.go:200), [publication test](../../internal/runbinary/runbinary_test.go:353).

Decline or defer FT267's extra dev-tier verifier unless an ordinary-build regression justifies it.
The ship tier already owns artifact verification and executes the script.
This is a proposed scope cut, not proof that the requested dev-tier behavior shipped.
[FT267](../../roadmap/FT267.md:3), [tier contract](../../projects/benchkit.md:466), [execution](../../internal/conformance/native_workflow_test.go:353).

FT222 needs a residual rewrite, not an automatic closure.
Current routing is explicit, but the row also asks for a configurable review binding and comparative implementation evidence.
Keep those decisions only if they remain desired.
[FT222](../../roadmap/FT222.md:17), [review binding proposal](../../roadmap/FT222.md:36).

FT287 proposes broader AXI policy. The current scoped policy does not prove that proposal complete.
Move it out of the defect queue and revisit it with FT336/FT341 evidence.
[FT287](../../roadmap/FT287.md:3), [current scope](../../.agents/skills/bench-craft-cli/SKILL.md:69).

Keep FT6 and FT24 parked under their existing triggers.
Keep FT327–FT330 and FT335 parked until their specified reproductions run.
FT328 deserves an early repro because the shared lock opens with write access.
A code read does not satisfy the row's runtime trigger.
[FT328 trigger](../../roadmap/FT328.md:7), [lock access](../../internal/chargeevidence/access.go:25).

FT140 cannot close as one unit.
Its old findings need individual verification, then closure, transfer, or an explicit risk decision.
This review did not re-open every historical pickup it cites.
[FT140](../../roadmap/FT140.md:3).

### Inference

The board overstates unfinished work when it carries a shipped premise beside a real residual.
A title can also hide a much larger staged spec.
FT115's low-priority deadline title hides a staged test-environment and shared-file isolation build.
[FT115](../../roadmap/FT115.md:6), [staged scope](../../specs/test-determinism/spec.md:15).

Five of the seven current staged specs omit a Roadmap line.
Only FT336 and FT290 name one.
The maintenance pass should map those specs to owners and stages before it removes any parent.
FT172 already owns that identity problem.
[FT172](../../roadmap/FT172.md:21), [staged coordination](../../specs/session-context-efficiency/spec.md:16).

## Question 3: What should move in priority?

| Order | Recommendation | Reason |
| --- | --- | --- |
| 1 | Keep FT336 first. | The reviewer explicitly set priority 1; its spec is staged. |
| 2 | Keep FT337 next; promote FT318 alongside it as a separate build. | Repeated evidence reads and manual invalid records both tax every review. |
| 3 | Promote FT115 to HIGH and FT290 near the top of ready work. | Both have staged specs. Test determinism affects trust; test projection affects diagnosis. |
| 4 | Promote FT215 and FT258's correctness residuals. Resolve FT293's remaining policy choices. | Lane omissions, merge state, and incomplete ownership create repeated repair work. |
| 5 | Shape FT341 before adding more Bench-specific exec conveniences. | One target contract can remove wrong-tree behavior across commands. |
| 6 | Shape FT305 with FT296 folded in. | Durable execution remains the next factory outcome after its required lifecycle contracts settle. |

These are ordering proposals, not new authorization or literal dependencies.
FT336 and FT337 keep their existing reviewer priority.
FT115 and FT290 can retain separate builds and reconcile their shared write surfaces.

Promote FT141's baseline-attribution residual within the reliability queue.
Retain its exact-tree distinction and remove its claim that nothing exists.
[FT141 recurrence](../../roadmap/FT141.md:23).

Use FT231's smallest useful baseline and the FT232 pilot as measured evidence work.
Do not make every improvement wait for the full experimental catalog.
The row already permits a minimal baseline before factory work.
A paid pilot still needs its stated authorization.
[FT231](../../roadmap/FT231.md:83), [FT232](../../roadmap/FT232.md:3).

Demote FT235 and FT308 to LOW unless an operational failure changes their value.
Directory naming and empty lock-file accumulation do not outrank correctness or repeated review cost.
Move FT302 below the delivery sequence until a fresh census identifies a concrete residual.
[FT235](../../roadmap/FT235.md:3), [FT308](../../roadmap/FT308.md:3), [FT302](../../roadmap/FT302.md:14).

Move FT331 and FT334 into the deferred section despite their current placement near the top.
Keep FT101, FT240, and cross-project tracing policy FT275 outside the near-term factory queue.
Each adds a new capability whose delivery does not resolve the current cost and reliability defects.

FT306 remains HIGH as a blocked milestone.
Its release qualification cannot become optional when execution work ships.
Keep FT58 and FT142 visible on that track.
[FT306](../../roadmap/FT306.md:3), [release conditions](../../ROADMAP.md:197).

## Row disposition ledger

Each row below is a proposal. A fold transfers every surviving requirement and occurrence before the original owner closes.

| Row | Disposition | Remaining outcome or reason |
| --- | --- | --- |
| [FT6](../../roadmap/FT6.md:1) | Keep parked | Move the structure-query candidate to FT307; retain the other evidence triggers. |
| [FT24](../../roadmap/FT24.md:1) | Keep parked | Upstream trigger unchanged; no fresh upstream capability assessment. |
| [FT38](../../roadmap/FT38.md:1) | Archive proposal | Optional visual revival has no current correctness outcome. |
| [FT58](../../roadmap/FT58.md:1) | Keep qualification blocker | Root permission and ownership policy remains a release-track decision. |
| [FT89](../../roadmap/FT89.md:1) | Narrow; LOW | Keep live contradictions and FT340's drain wording; remove stale inventory and shipped claims. |
| [FT98](../../roadmap/FT98.md:1) | Narrow | Keep preservation archives and cleanup edges; remove the delivered basic probe face. |
| [FT99](../../roadmap/FT99.md:1) | Fold | Semantic residuals to FT338; mechanical ownership facts to FT293. |
| [FT100](../../roadmap/FT100.md:1) | Keep blocked | Weight reduction stays behind FT231; correctness remains separate. |
| [FT101](../../roadmap/FT101.md:1) | Defer | Require a concrete multi-context adopter; scoped structure output belongs to FT307. |
| [FT106](../../roadmap/FT106.md:1) | Narrow; LOW | Keep semantic document checks; remove the fixed research-home claim. |
| [FT108](../../roadmap/FT108.md:1) | Retire proposal | Use the existing refactor route; preserve any accepted characterization rule in existing guidance. |
| [FT115](../../roadmap/FT115.md:1) | Promote; HIGH | Deliver the staged test-determinism spec; show its full scope in the row. |
| [FT125](../../roadmap/FT125.md:1) | Keep | Absorb bounded file reads and directory queries; preserve distinct reader seams. |
| [FT130](../../roadmap/FT130.md:1) | Keep; LOW | Only caller coordination during the gate window remains. |
| [FT140](../../roadmap/FT140.md:1) | Dismantle umbrella | Verify and dispose of each historical residual before retiring the owner. |
| [FT141](../../roadmap/FT141.md:1) | Promote residual | Baseline attribution remains useful; exact-tree verdict and log work already exists. |
| [FT142](../../roadmap/FT142.md:1) | Keep qualification track | Revalidate the old runtime residuals against the release contract. |
| [FT168](../../roadmap/FT168.md:1) | Keep; MEDIUM | Focused system/Markdown probes remain a distinct tooling gap. |
| [FT172](../../roadmap/FT172.md:1) | Keep; consolidate | Absorb FT292; connect staged specs, owner stages, and transitions. |
| [FT199](../../roadmap/FT199.md:1) | Keep | Recovery-aware ref classification remains distinct from removal mechanics. |
| [FT200](../../roadmap/FT200.md:1) | Narrow | Keep complete-preflight and spec-discovery policy; path authorization already runs at landing. |
| [FT204](../../roadmap/FT204.md:1) | Narrow | Keep session queries and the retro consumer; reuse delivered record readers. |
| [FT207](../../roadmap/FT207.md:1) | Keep decision | Malformed-admin protection across mutation callers remains a separate risk question. |
| [FT208](../../roadmap/FT208.md:1) | Split residuals | Refusal grammar and marker diagnostics differ from the HI14 fixture seam. |
| [FT215](../../roadmap/FT215.md:1) | Promote correctness | Prioritize lane omissions and empty merges; separate cost measurements from fixes. |
| [FT217](../../roadmap/FT217.md:1) | Narrow | Keep inventory deepening; move unrelated broker and canonical-path residuals to their owners. |
| [FT222](../../roadmap/FT222.md:1) | Revalidate; defer | Remove historical routing conflicts; retain only desired configuration and measurement decisions. |
| [FT231](../../roadmap/FT231.md:1) | Narrow; measure incrementally | Use delivered records; retain causal trials and missing measures without a giant prerequisite. |
| [FT232](../../roadmap/FT232.md:1) | Prioritize experiment decision | The collector exists; authorize a real pilot and review its terminal report. |
| [FT235](../../roadmap/FT235.md:1) | Demote; LOW | Physical directory labels have less value than command-level selection and correctness. |
| [FT240](../../roadmap/FT240.md:1) | Keep parked | External retrieval experiment remains gated by FT231 and its named subject. |
| [FT241](../../roadmap/FT241.md:1) | Split requirements | Record mechanics to FT318; versioned acceptance promises remain LOW. |
| [FT243](../../roadmap/FT243.md:1) | Narrow; LOW | Keep recurring selection policy; transfer concrete query and probe requests. |
| [FT244](../../roadmap/FT244.md:1) | Keep; LOW | Assignment-owned scratch lifecycle is distinct from output bounds. |
| [FT246](../../roadmap/FT246.md:1) | Closure candidate | Confirm the constructor census; retain builds whose subject is binary construction. |
| [FT249](../../roadmap/FT249.md:1) | Fold residual | Move reader consistency to FT265; avoid preselecting Git-ref storage. |
| [FT253](../../roadmap/FT253.md:1) | Keep | One destination lease could support durable execution; do not merge it with cleanup locks. |
| [FT254](../../roadmap/FT254.md:1) | Split responsibilities | Target attribution to FT341; file reads to FT125; retain exec and conflict ergonomics. |
| [FT255](../../roadmap/FT255.md:1) | Keep distinct | Shared capacity control differs from test correctness and stress scheduling. |
| [FT258](../../roadmap/FT258.md:1) | Promote correctness | Resolve merge-parent behavior first; retain explicit ownership for discovered paths. |
| [FT260](../../roadmap/FT260.md:1) | Reassess scope | Diff inspection can reuse reader work; patch transfer needs its own current justification. |
| [FT261](../../roadmap/FT261.md:1) | Park pending repro | Current source-range facts include untracked paths; runtime closure is unverified. |
| [FT262](../../roadmap/FT262.md:1) | Fold into FT293 | Keep progress output and structured closure diagnostics. |
| [FT265](../../roadmap/FT265.md:1) | Keep; consolidate residual | Absorb FT249; reuse sealed generations and settle retrieval consistency. |
| [FT267](../../roadmap/FT267.md:1) | Decline or defer proposal | Ship already verifies artifacts; extra dev-tier execution needs a distinct justification. |
| [FT275](../../roadmap/FT275.md:1) | Defer cross-project policy | Bench instrumentation does not justify a new universal adoption obligation. |
| [FT283](../../roadmap/FT283.md:1) | Keep decision | Phase-scoped transition behavior remains distinct from retirement. |
| [FT284](../../roadmap/FT284.md:1) | Keep decision | Atomic retirement must preserve cited decisions and occurrence-ledger integrity. |
| [FT287](../../roadmap/FT287.md:1) | Defer policy expansion | The current scoped AXI contract exists; broadening is not a missing current obligation. |
| [FT290](../../roadmap/FT290.md:1) | Promote ready work | Build the staged test projection; retain its separately noted failure-detail residual. |
| [FT292](../../roadmap/FT292.md:1) | Fold into FT172 | Atomic recording follows the roadmap identity and transition contract. |
| [FT293](../../roadmap/FT293.md:1) | Narrow; prioritize decisions | Keep generic owner closure, deleted paths, and proposal adoption policy; absorb FT262. |
| [FT296](../../roadmap/FT296.md:1) | Fold into FT305 | Worker lifetime is part of durable supervisor design. |
| [FT299](../../roadmap/FT299.md:1) | Conditional broker grouping | Group rehearsal with broker lifecycle work after the residuals and FT327 repro are confirmed. |
| [FT302](../../roadmap/FT302.md:1) | Demote from front section | Use a fresh census before scheduling the remaining deepening. |
| [FT304](../../roadmap/FT304.md:1) | Keep after identity | Observation-only view remains distinct; FT172 is its literal dependency. |
| [FT305](../../roadmap/FT305.md:1) | Keep HIGH; consolidate | Absorb FT296 and preserve explicit limits and recovery authority. |
| [FT306](../../roadmap/FT306.md:1) | Keep blocked HIGH milestone | Retain the qualified-release and adoption evidence requirements. |
| [FT307](../../roadmap/FT307.md:1) | Keep active | Absorb matching structure requests from FT243 and FT6. |
| [FT308](../../roadmap/FT308.md:1) | Demote; LOW | Stale lock accumulation needs safe ownership design but lacks a demonstrated urgent failure. |
| [FT312](../../roadmap/FT312.md:1) | Keep; clarify stages | Invocation recipes and review-venue capacity are related but independently deliverable. |
| [FT314](../../roadmap/FT314.md:1) | Defer decision | Cross-run gate reuse changes authorization and must remain separate from ordinary performance fixes. |
| [FT317](../../roadmap/FT317.md:1) | Fold option into FT318 | Set capability-blocked reconciliation semantics before the writer implements them. |
| [FT318](../../roadmap/FT318.md:1) | Promote; HIGH | Native validated records and plan amendments remove recurring review overhead. |
| [FT324](../../roadmap/FT324.md:1) | Retirement candidate | No competing project rule was identified in-tree; retain only with a current source. |
| [FT325](../../roadmap/FT325.md:1) | Fold into FT326 | Wrap-independent canary mutation belongs with anchor semantics. |
| [FT326](../../roadmap/FT326.md:1) | Keep; consolidate | Retain sentence and paragraph decisions; same-kind distinctness has shipped. |
| [FT327](../../roadmap/FT327.md:1) | Keep parked | Require its clean-worktree reproduction before a broker lifecycle build. |
| [FT328](../../roadmap/FT328.md:1) | Keep parked; prioritize repro | Write access for a shared read lock makes the incident plausible, not runtime-verified. |
| [FT329](../../roadmap/FT329.md:1) | Keep parked | The recorded current-tree replay did not reproduce the old symptom. |
| [FT330](../../roadmap/FT330.md:1) | Keep parked; focused repro | Verify the exact refusal route before changing it. |
| [FT331](../../roadmap/FT331.md:1) | Move to deferred section | The advisor pilot is a separate capability, not a prerequisite to the cost fixes. |
| [FT332](../../roadmap/FT332.md:1) | Defer | The scheduled stress job stays outside the staged test-determinism build. |
| [FT333](../../roadmap/FT333.md:1) | Keep distinct | Tickets-only prepared evidence changes the preparation subject, not just record writing. |
| [FT334](../../roadmap/FT334.md:1) | Move to deferred section | Deadline batching follows dependable execution and measurable costs. |
| [FT335](../../roadmap/FT335.md:1) | Keep parked | Verify source identity and reproduce full-gate selection before graduation. |
| [FT336](../../roadmap/FT336.md:1) | Keep priority 1 | The staged output-bound build remains the explicit first choice. |
| [FT337](../../roadmap/FT337.md:1) | Keep priority 2 | Retain file-stable review evidence and the full-retrieval control. |
| [FT338](../../roadmap/FT338.md:1) | Keep; consolidate | Absorb FT99 semantic checks; remove already-shipped sweep requirements. |
| [FT339](../../roadmap/FT339.md:1) | Fold and retire owner | Registry attribution to FT341; directory projection to FT125. |
| [FT340](../../roadmap/FT340.md:1) | Fold residual; retire owner | Repair-ticket choice is settled; drain lane wording belongs with FT89. |
| [FT341](../../roadmap/FT341.md:1) | Promote shaping | Absorb wrong-tree attribution and clarify command scope before more exec convenience work. |

## Verification record

- [x] Read every row body across the coordinator and two read-only delegates.
- [x] Re-open the sources behind the principal consolidations and stale-premise findings.
- [x] Separate source facts, inferences, and proposed dispositions.
- [x] Keep runtime closure, upstream capability, and private-memory limits explicit.
- [x] Pin the source tree and each recommendation's refresh trigger.
- [x] Retain a disposition for each current row.
- [x] Pass the report prose check and validate all 151 local links.
- [x] Match the disposition ledger to all 82 current row files.
- [x] Confirm the primary checkout remains clean at the pinned source commit.

No product test, gate, paid experiment, or external qualification ran.
This review does not certify runtime closure or release readiness.
The dependency tables express the material relationships; no additional diagram is necessary.

## Validation plan

1. Approve the desired folds, scope cuts, and priority changes.
2. Reproduce FT261, FT328, and other parked incidents before declaring runtime closure or starting a fix.
3. Recheck FT246's constructors and each surviving FT140 finding before retiring their rows.
4. Map every staged spec to its owner and stage.
5. Move surviving requirements and occurrence keys before removing an old owner.
6. Update dependency and sequence references in the same roadmap change.
7. Run the edited Markdown prose checks and the required roadmap checks.
8. Land the approved maintenance change through the normal Bench workflow.

Fresh-session continuation: Review this report against main at 673ebcf8924bb2c91f78c8c40c98692eb08531b9.
Treat every disposition as a proposal until the reviewer approves it.
Then follow .agents/commands/bench-drain.md for the authorized roadmap maintenance.
Keep FT336 first, preserve parked triggers, and retain the release NO-GO boundary.
