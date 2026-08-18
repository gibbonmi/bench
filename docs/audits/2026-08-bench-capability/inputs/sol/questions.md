# Required adversarial questions

Evidence identifiers refer to `evidence.md`. Each answer is deliberately short;
the A–Z report carries the synthesis and recommendations.

1. **What does Bench ask modern Claude or Codex to do that they already do reliably without help?** `[INFERRED]` Read repository instructions, inspect a diff, run focused tests, avoid unrelated edits, summarize results, and choose ordinary tools. Bench repeats these as multi-page craft and phase policy without a paired no-Bench result (E014).

2. **What does Bench make more consistent?** `[OBSERVED]` Exact-tree gating, prospective atomic commits, worktree ownership, preflight fences, acceptance-row membership, status rendering, package/release evidence, and harness-prefix rendering have executable owners and passing adversarial tests (E007–E008, E013). `[REPEATED PRACTITIONER EVIDENCE]` The debug loop also makes reproduction-first behavior more consistent, but this audit did not causally measure the delta.

3. **Which consistency gains are merely claimed, and which are observed?** `[OBSERVED]` The mechanisms in answer 2 through preflight/gate/landing are observed. `[HYPOTHESIS]` Better shaping, better specs, better review independence, improved model routing, mandatory delegation value, and general skill effectiveness are claims because conformance proves prose presence rather than model behavior (E010, E013–E015).

4. **Which Pocock-derived behaviors are load-bearing?** `[OBSERVED + REPEATED PRACTITIONER EVIDENCE]` A red-capable bug loop, pre-agreed high seams, tracer-bullet slicing, facts-versus-decisions grilling, glossary/ADR separation, fixed-base review, and spec-to-ticket continuity. They survive deletion tests because removing them loses a distinct decision or executable feedback mechanism.

5. **Where has Bench weakened an upstream skill through composition?** `[OBSERVED]` `diagnosing-bugs` becomes a 140-line phase entangled with worktree/shift policy and claims a check “it didn’t author” although the same phase authors the regression. Review adds a Coverage axis but makes its persistence/commit route internally inconsistent (E017). Shaping can stop after charting even when the upstream frontier discipline would continue.

6. **Where has Bench improved an upstream skill through state or enforcement?** `[OBSERVED]` Specs gain acceptance-coverage rows, ownership fences, explicit preflight, exact landing, and gate evidence. Worktrees gain ownership and safe release. Debugging gains isolated-entry and gate integration. These are material implementation additions, not renames (E007–E008).

7. **Where does Bench ask a model to remember something software should enforce?** `[OBSERVED]` Phase progress, `--full` orchestration, review completion, diagnosis attempts, line/effort adherence in Codex, final-check duties, and which requirements remain open. No canonical active work record owns these (E010, E015).

8. **Where does Bench ask software to decide something that requires judgment?** `[OBSERVED]` `firstInvocable` treats syntactic command shape as the next semantic action, and status severity decides what the handoff presents even when higher-severity rows are prose (E018). Structure budgets also flag size as a split recommendation despite the profile conceding size is not architecture.

9. **Where does Bench confuse available context with active context?** `[OBSERVED]` Cold instructions require roughly 6,000 words before task guidance, the whole guidance estate is 31,456 words, and the profile treats availability in files as enough (E014). There is no compiler selecting only Goal/Core/Verified/Open/Next for the present action.

10. **Where does Bench confuse a plan with current state?** `[OBSERVED]` `capture/session-handoff.md` preserves a historical roadmap narrative at the current commit while status reports no lag (E011–E012). Ticket-only spec directories are also historical plans that live outside the current spec status model (E016).

11. **Where does Bench confuse assertion with observation?** `[OBSERVED]` Handoff `State`, review findings, retros, scorecards, and assumptions in tickets can carry prose without evidence references. The claim CLI named by the audit model does not exist (E010–E012).

12. **Where does Bench confuse a passing test with requirement satisfaction?** `[OBSERVED]` The gate is green while the board has an invalid map and 62 structure issues; more importantly, 233 prompt-anchor fixtures prove strings/contracts are present, not that a model follows them (E009, E013). The docs correctly call review semantic, but no executable semantic control closes the gap.

13. **Where does Bench confuse repeated review with independent review?** `[OBSERVED]` Guidance asks for fresh axes, but no runner records reviewer identity/context or a clean review. A second model pass can be called independent without verifiable isolation (E010). Current no-mistakes explicitly starts each review turn isolated; Bench only says to.

14. **Where does Bench confuse subagent count with useful parallelism?** `[OBSERVED]` Assess mandates six area delegates, review mandates three axes, deepen requires a delegate and can trigger three alternative designers, and spec builds mandate a fresh writer per ticket. Fixed counts precede a measured dependency/latency case.

15. **Where does Bench confuse detailed prose with control?** `[OBSERVED]` `--full`, Review, Final Check, What Next, line governance, and most phase transitions are Markdown contracts. Unknown-subcommand probes prove the corresponding control surface does not exist (E010, E014).

16. **Where does Bench allow blank retries?** `[OBSERVED]` `bench shift` preserves some failure evidence, but ordinary model-driven phases have no attempt schema or tripwire. The craft-line ladder discourages repetition; software does not refuse a same-class retry without new evidence.

17. **Where does Bench let an agent patch before reproducing?** `[OBSERVED]` The debug path forbids it. Generic light-path implementation and TDD outside marked seams can edit before an original red signal; `craft-tdd` says the gate may catch regressions later. Preflight checks scope, not reproduction.

18. **Where can an agent keep thinking after a discriminating test is available?** `[OBSERVED]` Everywhere outside `bench-debug`: no runtime pauses prose reasoning when a runnable discriminator exists. Review explicitly forbids tests during review, even where one could distinguish a correctness concern, deferring reality to a later phase.

19. **Where can stale evidence authorize work?** `[OBSERVED]` Exact gate evidence cannot: tree/oracle drift invalidates it (E007). Handoff state can steer the next session while semantically false at the same commit, and `firstInvocable` can select a later command from that board (E011–E012, E018).

20. **Where can local verification masquerade as global completion?** `[OBSERVED]` Dev gate green is called “shippable” in the gate script while ship-tier/native capabilities and live publication are separate; seven capability skips remained. Final Check asks the report to qualify this, but no machine-enforced global claim type exists (E008, E020).

21. **Which instructions disappear if the CLI improves?** `[INFERRED]` Command discovery, phase-entry checks, next-action selection, handoff pin rendering, line validation, review pickup persistence, attempt/failure inheritance, final-check housekeeping, and most `--full` narration can become schemas/state transitions. Judgment guidance for grilling, seams, debugging, and review remains.

22. **Which CLI commands exist because the workflow is unnecessarily complicated?** `[INFERRED]` `skills-index`, `anchors`, separate handoff/status plumbing, several gate internals, and spec lifecycle repair verbs expose maintenance needed by duplicated prose/state. Keep plumbing internal; retain status, gate, commit, worktree, preflight, coverage, test, and release state machine as public capabilities.

23. **What would Agentless remove?** `[RESEARCH + INFERRED]` Default decision mapping, mandatory specs/tickets for bounded work, per-ticket writer delegation, fixed three-axis fan-out, repeated handoffs, line ceremony, and retro/roadmap drains. It would keep localization → repair → validation and require each added layer to beat that baseline.

24. **What would SWE-agent change about the interface?** `[RESEARCH + INFERRED]` Make the root content-first, return minimal structured state/errors/next action, combine common operations, and hide plumbing. Bench's compact query commands align; the wrapper/binary split and prompt-only phases do not (E001, E004, E019).

25. **What would ReAct change about the reasoning/action cadence?** `[RESEARCH + INFERRED]` Require an observation after each material hypothesis and update the next action from it. Bench-debug already does this; shaping/review/maintenance allow long prose runs without a discriminator.

26. **What would Lost in the Middle change about context compilation?** `[RESEARCH + INFERRED]` Stop loading whole profiles and skill graphs because a fact is “somewhere.” Produce a small ordered active capsule, with source pointers for zoom-in (E014).

27. **What does J-space research suggest about active state without justifying J-space terminology?** `[RESEARCH]` Select a small working set, deliberately maintain it, and broadcast shared facts to consumers. It does not imply that a Markdown ledger is neural J-space or that latent reasoning is required.

28. **Which J-Space Suite ideas would be cargo cult if copied?** `[RESEARCH + INFERRED]` Persona phrases, “dense rails,” cognitive-module branding, latent-state metaphors, fixed Goal/Core labels without a consumer, and model-specific first-turn anchors. Adopt only observable retry, checkpoint, coverage, and selective-loading contracts after A/B tests.

29. **What state must survive a model swap?** `[INFERRED]` Goal, authorized scope/fences, immutable base/tip, accepted decisions, requirements and coverage, verified commands with subject identities, open uncertainties, next discriminator, failed attempts with diagnostics, and ownership/recovery handles.

30. **What state should deliberately not survive?** `[INFERRED]` Chain-of-thought, speculative hypotheses already falsified, conversational narration, duplicated source excerpts, model persona, transient tool chatter, and obsolete plans.

31. **Which facts should be broadcast?** `[INFERRED]` Canonical names, goal, immutable refs, scope/fences, accepted decisions, exact requirements, current verified evidence/freshness, active failure diagnosis, and the next action. Broadcast pointers, not full files.

32. **Which conclusions require independent rediscovery?** `[INFERRED]` Semantic review findings, root cause, requirement satisfaction, risky scope decisions, test adequacy, and whether a claimed fix closes the original signal. Do not seed certifiers with implementer conclusions.

33. **Is `what-next` already the right router?** `[OBSERVED]` No. It is a 210-line roadmap/capture maintenance phase, intentionally not a workflow phase, and it can write a precommit handoff that is semantically stale at HEAD (E011, E014).

34. **Is `/bench` a useful front door or only another alias?** `[INFERRED]` Useful only as a thin logical router over user intent plus deterministic state, with `/bench` and `$bench` as harness adapters. If it merely aliases help or status, delete it. The canonical product entry should be the logical `bench` router, not another workflow phase.

35. **What should `bench init` own, if anything?** `[INFERRED]` Transactional installation of the minimal platform, gate selection/scaffold, adapter/hook wiring, one project profile pointer, and initialization of the work-state store. It should not interview architecture, choose seams, or author policy.

36. **Can a new user find the right path without reading the skill graph?** `[OBSERVED]` No. Root output is an inventory, status is repository-only, and the current Codex phase adapters were not surfaced (E015, E019).

37. **Can a fresh Codex session resume Claude's work?** `[OBSERVED]` Partially. Repository pins and a derived phase command translate, but semantic `State` is unvalidated and current Codex did not expose phase adapters (E011–E012, E015).

38. **Can a fresh Claude session resume Codex's work?** `[OBSERVED]` Partially for the same durable pins/state. Claude has phase symlinks and an Agent line hook that Codex lacks, so behavior and enforcement are not symmetric.

39. **Can the `diagnosing-bugs` behavior be preserved with less prompt load?** `[HYPOTHESIS]` Yes: preserve the red-capable original loop, minimize, ranked falsifiable hypotheses, one-variable probes, regression at the real seam, original-loop rerun, and cleanup; move tool examples to on-demand reference. It requires the benchmark in W before replacement.

40. **Does reducing it lose the mechanism that breaks repair loops?** `[REPEATED PRACTITIONER EVIDENCE]` It will if reduction removes the already-run exact-symptom loop or original-signal rerun. Word-count reduction alone is not the goal; those causal constraints are the protected core.

41. **If Bench were rebuilt today, what 20% preserves 80% of demonstrated value?** `[INFERRED]` Exact gate/landing, ownership-safe worktrees, deterministic status/preflight/coverage, one compact work-state/checkpoint, a thin router/context compiler, and four practices: debug, seams/TDD, spec/tickets for genuinely wide work, independent semantic review.

42. **What existing behavior would a rewrite accidentally destroy?** `[OBSERVED + INFERRED]` Exact tree+oracle cache binding, atomic prospective commits, moved-subject refusal, worktree identity/recovery, hostile-input fixtures, release evidence/publication state, harness prefix single-sourcing, and the debugging original-signal loop.

43. **What is the smallest architectural spine that can absorb the current proven workflows?** `[INFERRED]` `init → logical bench entry → deterministic status + work-state → context compiler → selected practice → executable observation → evidence record → gate/landing → checkpoint`. Existing deep modules remain behind it.

44. **What should be deleted before anything new is added?** `[INFERRED]` Retired ticket-only spec directories, invalid maps, stale handoff narrative, prompt-only `--full` claims, fixed-count delegation rules, generic design-system guidance from non-UI profiles, duplicate command inventories, and public plumbing commands that no external consumer needs.

45. **What is the one next ticket with the highest leverage?** `[INFERRED]` Build the canonical work-state/router/context-compiler tracer described in report section Z, with old status/handoff read compatibility and an A/B benchmark before routing more phases through it.
