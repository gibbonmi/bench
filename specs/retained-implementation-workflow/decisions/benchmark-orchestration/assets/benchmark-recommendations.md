# Orchestration recommendations from the FT311 benchmark

This asset preserves historical benchmark evidence and proposals. Resolved decision tickets in the parent map supersede its policy recommendations.

Recommend one retained implementer for each coherent group of tickets, with a coordinator where independent groups need scheduling or integration. Preserve independent Standards, Spec, and Coverage review and the required independent probes in every implementation topology. Move repeated mechanical checks into their existing Bench owners. Select models by complete accepted-change cost and verified quality.

Status: proposal, after the final comparative review. No Bench rule, implementation, model default, or candidate has changed. Origin for every recommendation below: benchmark learnings. These are recommendations for a reviewed drain and subsequent specification, not an adoption verdict.

Consumed by: the orchestration proposal in [capture/IDEAS.md](/home/mgibs/workspace/bench/capture/IDEAS.md:1), the solo-review learning in [capture/learnings.md](/home/mgibs/workspace/bench/capture/learnings.md:7), and the roadmap owners named below. Drift: refresh when the frozen evidence is corrected, comparable trials finish, or the named owners change. Retire when reviewed decisions and their implementation records replace these proposals. Repository baseline: `cf2caa6d5c195c6b6ed5852fb4478902d1d99fbd` on `main`.

## Evidence and limits

The [final Astra code comparison](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:1) owns the findings. The [four-run report](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-four-arms.md:32) owns cost assumptions and timing. The arm records preserve observed interventions, including failed attempts. The later Astra comparison supersedes earlier clean-review claims.

| Run | Strength to retain | Constraint on the conclusion |
|---|---|---|
| Solo Astra | Retained context, local repair, and correct binding of staged content in the plan fingerprint. | Its original run omitted independent review. It also duplicated production mechanisms and omitted several regression proofs. |
| Delegated Astra | Shared sweep and command owners, plus stronger adversarial regression coverage. | Independent review still missed a stale staged-state authorization defect and the shared ignored-path collision. |
| Delegated Sol | Completed the ticket sequence with relatively few intervention groups and a slightly lower recorded total than delegated Astra. | A later probe reproduced a registration repair defect. This is one task, not evidence for a general model default. |
| Delegated Terra | Its failed attempts expose specific opportunities for better diagnostics, test selection, and evidence handling. | Lower worker prices did not yield lower total cost. Coordinator effort and repairs increased. This run does not establish a capability advantage to adopt. |

Sources: [Astra code and coverage differences](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:100), [Sol interventions](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/sol/result.json:49), [Terra interventions](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/terra/result.json:716), and [comparative outcomes](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-four-arms.md:9). Treat the last row as a learning from failures, not a claim that each run must contribute a winning practice.

Recorded estimated API totals were $81.79 delegated Astra, $104.89 Terra, $79.22 Sol, and $22.47 solo Astra before the required independent review. The subsequent review assessed both Astra candidates and has not been allocated to either original run. Cached input cost alone was $60.10, $77.34, $61.07, and $15.38 respectively. These use saved pricing assumptions, not verified charges or subscription consumption. [Cost record](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-four-arms.md:32); [review accounting boundary](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:121).

There was one run per configuration. Review, iteration policy, advance fence corrections, and context retention were not held equal across all four. The results support targeted changes and further trials. They do not prove that solo work, a particular effort level, or any model always wins. [Comparison limits](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-four-arms.md:55).

## 1. Who should own implementation?

**Inference:** Context continuity is a promising source of savings. The solo result makes it worth testing, but its review omission prevents a causal cost claim. The existing chunk proposal already describes the appropriate unit: related tickets with one retained author, separate verification and commit checkpoints, and one isolated worktree. [Existing proposal](/home/mgibs/workspace/bench/capture/IDEAS.md:5).

**Proposed principles:**

1. Assign a coherent area of behavior to one implementer. A ticket marks acceptance and a green checkpoint; it need not force a new author or context.
2. Retain that implementer for related repairs while its context remains useful. Restart for a change of domain, unusable context, or repeated failure to advance the same predicate.
3. Let the author choose routine implementation details within the approved behavior, contracts, and write authority. Keep behavioral scope decisions with the reviewer.
4. Define an iteration as a coherent implementation-and-verification attempt. Do not count every tool call as an attempt. For an explicitly uncapped run, continue while useful work remains within the authorized scope. Report repeated non-progress or an external blocker. Preserve any explicit user budget or stop condition.

The Astra run spent interventions on expected-output corrections, evidence classification, and write fences. Terra's attempt accounting had to be clarified because earlier counts included incremental verification batches. These support clearer authority and accounting, not permission to weaken a test or waive a failing gate. [Astra CI1–CI9](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/astra/result.json:47); [Terra accounting](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/terra/result.json:40).

**Disposition: Fold** into the existing chunk idea and the implementation/delegation discipline. The current fresh-worker-per-ticket and cap rules remain in force until the reviewer adopts a change. The user's uncapped benchmark approval does not silently rewrite kit defaults.

## 2. When does delegation help?

**Inference:** Delegation earns its cost when it separates independent work or supplies independent judgment. A serial chain of tightly related tickets offers less parallelism and repeatedly pays for context setup. This benchmark does not isolate how much of the observed difference comes from those costs.

**Proposed principles:** Parallelize groups only when their contracts and actual write closures permit it. Keep a coordinator responsible for dependencies, integration, evidence, and decisions. Avoid routinely repeating the implementer's full analysis. For one coherent group, the primary agent can implement it and later commission independent review. For several independent groups, use a coordinator and retained authors.

Prepare references, acceptance ownership, and write closure once per frozen input. Reuse the prepared artifact while its identities remain valid. Supply complete evidence through bounded retrieval; do not replace full review with an unverified summary. Terra repeated a full review packet after truncation even though its prepared identities matched. [Recorded repeat](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/terra/result.json:873).

**Proposed fence behavior:** Derive a group's write authority from its approved tickets and required closure. Preserve the spec's outer authority and prevent simultaneous conflicting writers. Do not grant the whole spec fence to every parallel author. Whether closure-only overlap creates a dependency is an existing unresolved reviewer decision. [FT293](/home/mgibs/workspace/bench/roadmap/FT293.md:1).

**Disposition: Fold** into the chunk proposal, FT293, and existing preflight preparation. This does not justify a new general orchestration engine or a second ticket-status database.

## 3. What must review always retain?

**Tested results:** Both Astra candidates overwrote ignored bytes when the target tracked the same path. Delegated Astra accepted a stale plan after a staged-only change; solo Astra refused it. A separate fresh-process check passed for both, so that finding concerns missing permanent regression coverage, not a demonstrated runtime failure. [Spec findings and probes](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:41); [coverage distinction](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:76).

**Proposed principles:**

1. Separate implementation delegation from assurance. Solo implementation retains independent review and required independent verification. Label author probes truthfully.
2. Preserve all three fresh review axes for the next trial. This benchmark provides no evidence for reducing reviewer count or removing an axis. The existing idea of one reviewer covering all axes remains untested.
3. Review the integrated behavior against real producer states, including interactions between acceptance rows. A green gate or a completed row count does not settle an untested interaction.
4. Prefer shared policy and transaction owners across related modes. Pair that ownership with tests for state changes between plan and apply. Cover lock and registration states, ignored/tracked path intersections, conditional ref deletion, and required process boundaries.
5. Require each finding to identify the violated contract, reachable input, and evidence. Deduplicate overlapping findings across axes. Keep missing proof distinct from a reproduced bug.

The review rejected a fabricated envelope that the actual producer cannot create under the closed trust model. Preserve that decision. The ignored-path collision is reachable and exposes conflicting required outcomes. Resolve its behavior explicitly before implementation; the report recommends refusal before movement. Do not silently broaden recovery capture. [Disposition and refutations](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:134).

**Disposition: Fold** the process correction into current review and verification ownership. **Recommend** cross-state analysis and shared-owner design through existing domain, spec, and review disciplines. Do not create a second checklist containing a copied acceptance map.

## 4. Which repeated work belongs in Bench?

The following are proposed extensions or contract audits, not assertions that all listed behavior is absent today. Start with the current owner and reproduce the exact gap before adding code.

| Priority | Recommended change | Existing home and classification | Acceptance proof for the change |
|---|---|---|---|
| P0 | Make omitted independent review visible in completion evidence. Bind the author/verifier identities, completed axes, and reviewed base/tip to the existing workflow record. Distinguish author-only evidence. | Review/preflight and existing completion ownership; **Fold**. FT200's enforcement chokepoint remains a reviewer decision. | A solo-author-only case cannot be reported as independently reviewed. A changed tip makes the prior review visibly stale. No second gate verdict or duplicate receipt system. |
| P0 | Make verification outcomes unambiguous: selected tests, tests actually run, compile failure, no match, pending execution, and capability skip. Preserve baseline attribution when blaming an unchanged test. | `bench test`, existing diagnostics, [FT290](/home/mgibs/workspace/bench/roadmap/FT290.md:1), and FT311 diagnostics residuals; **Fold**. | Reproduce the mistaken selector and compile failure. Neither can be presented as the intended passing test run. A still-running check cannot supply completed evidence. |
| P1 | Make mutation evidence distinguish an intended assertion failure from a broken mutation, failed setup, or unmatched test. Use existing preservation/restoration behavior. | Existing probe owner and [FT168](/home/mgibs/workspace/bench/roadmap/FT168.md:1); **Fold**. | The benchmark's invalid mutation is classified invalid. A valid omission makes the named assertion fail. Restoration is verified on both paths. |
| P1 | Reuse prepared review/build evidence without repeating large outputs. Keep exact identities, completeness, and invalidation visible. | Existing preflight preparation and [FT254](/home/mgibs/workspace/bench/roadmap/FT254.md:1); **Fold**. | A truncated response is detectable and the unread portion can be retrieved. Reuse across unchanged inputs succeeds; changed inputs invalidate it. |
| P1 | Derive group fences and scheduling conflicts from existing ticket data and closure. | Chunk proposal and FT293; **Fold**, pending the closure-overlap decision. | The missing test-file case is identified before writing. Related sequential work avoids repeated fence amendments. Conflicting concurrent writers still refuse. |
| P1 | Export one complete assessment record from native evidence. Include phase and role costs, failures, repairs, review, verification, skips, and intervention categories. | [FT231](/home/mgibs/workspace/bench/roadmap/FT231.md:1), existing Bench spans/census plus harness usage; **Fold**. | Reconstruct the saved totals without double-counting cached input or restarted counters. Missing review cost remains unknown, never zero. |

Evidence for these priorities: [Astra fence interventions](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/astra/result.json:89), [Astra invalid mutation](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/astra/result.json:396), [Terra selection and evidence failures](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/terra/result.json:716), [solo repair record](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/direct/result.json:1), and [solo review omission](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-astra-code-review.md:121).

The workflow record can prove that a review artifact exists for a frozen pair; it cannot mechanically prove the semantic judgment is correct. The gate remains the oracle. FT200 explicitly rejects a second done-oracle and notes removed receipt machinery. Use existing owners and decide the narrow enforcement point before implementation. [FT200 boundary](/home/mgibs/workspace/bench/roadmap/FT200.md:3).

**Skip:** A new wrapper for every observed call sequence, duplicated policy registries, universal adversarial checklists, and a second landing engine. First check the already shipped FT311 preparation, diagnostics, and landing work. [FT311 scope](/home/mgibs/workspace/bench/roadmap/FT311.md:1).

## 5. How should we measure and select the workflow?

**Inference:** Cheap worker tokens can be a false economy. Terra's recorded total exceeded both other delegated runs. Sol's total was only slightly lower than delegated Astra. Cached input remained the largest dollar component. Reducing repeated context and supervision is worth testing even when cache hit rates are high. [Cost table and limitations](/home/mgibs/workspace/bench/capture/benchmarks/ft311-recoverable-reset-20260910/comparison-four-arms.md:32).

**Proposed principles:** Optimize complete accepted-change cost subject to quality, and report elapsed time separately. Count implementation, coordination, failed attempts, repairs, independent review, verification, and assessment. Separate cached input, uncached input, and output. Preserve unknown actual charges. Categorize interventions as implementation defects, evidence defects, scope/fence decisions, environment failures, or routine scheduling.

**Recommend:** Use retained Astra authorship as the next topology experiment, with the same independent review and probes as the comparison arm. Test Sol on bounded coherent work before considering a default change. Keep Terra's proposed lower-tier trials within the already identified bounded routes; this benchmark does not support moving safety judgments to it. Do not reduce coordinator effort and change implementation topology in the same experiment. [Existing trial boundaries](/home/mgibs/workspace/bench/roadmap/FT311.md:1).

The benchmark does not establish the optimum chunk size, restart threshold, coordinator effort, review model, or net savings after equal assurance. The later Sol/high reviews found meaningful defects, but there was no controlled reviewer-model comparison. All of these remain unknown.

## Validation and implementation sequence

1. **Capture and decide.** Fold this proposal into the existing chunk idea and review learning through a reviewed drain. Resolve chunk authority, closure overlaps, continuation policy, and the narrow completion-evidence contract. Do not reopen the envelope trust model or merge a benchmark candidate as part of orchestration work.
2. **Build the narrow evidence fixes first.** Reproduce each P0/P1 CLI gap against the current tree. Drop anything already satisfied. Specify only the remaining gap, its existing owner, and its named failure proof. Implement in isolated assignments with the normal review and gate.
3. **Trial retained authorship.** Compare current per-ticket authorship with retained authorship for the same coherent task. Hold model, effort, baseline, acceptance, fences, review, independent probes, and iteration policy equal. Include coupled and independently parallel tasks. Keep cases used to develop the proposal separate from held-out evaluation.
4. **Measure model and effort separately.** Repeat the promising topology with one model or effort change at a time. Preserve failures and report variation. Follow FT231's no-Bench/current-Bench/one-capability design when making a causal claim about a kit capability. Do not turn a noisy benchmark into a per-commit gate.
5. **Adopt only after verification.** Run the synthesis legibility, consistency, and dogfood loops on the proposed kit change. A trigger change needs a fresh-session trial. Require preserved quality and an observed benefit before setting a default. Then update the canonical owner and retire the consumed proposal.

Verification record: the recommendation separates observed facts, tested results, inferences, and proposals. It cites the frozen reports and arm records. It preserves conflicting requirements and unknown costs. Each implementation proposal names an existing owner.

No new performance experiment, source-code change, or dogfood run was performed for this document. Proposal delivery is complete; kit adoption is not.
