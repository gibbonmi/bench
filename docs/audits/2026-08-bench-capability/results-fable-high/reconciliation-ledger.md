# Reconciliation ledger

Subject: `58d966e2` (both prior audits and this run — EXACT MATCH, see `baseline.md`).
Every entry below was re-derived against this tree. "Fresh" = observation made in this
run; "prior" = the auditor's own evidence artifact. Classification vocabulary is the
prompt's: REPRODUCED · COMPATIBLE · CONTRADICTED · PARTIALLY SUPPORTED · UNSUPPORTED ·
STALE · UNINVESTIGATED · NOT ACTIONABLE. Basis labels: OBSERVED · INFERRED ·
PRACTITIONER EVIDENCE · HYPOTHESIS · RESEARCH-INSPIRED.

---

```yaml
- id: L-01
  topic: Root conformance suite (29 checks) does not run in the dev gate; live tree is red
  sol_position: Not found. E008 records the green gate with "seven capability skips" and E013 that 233 canary bindings pass; Sol treats conformance as proving anchors, not behavior.
  opus_position: Headline P0 (E6/I.1). Gate has no conformance phase; TestRootConformance skips without BENCH_CONFORMANCE_ROOT; only bench prep-release sets it; 10 diagnostics red at HEAD; unwired "since 72c037a1, 43 days".
  relationship: Opus-only finding; Sol silent (not contradicting).
  prior_evidence: Opus 02-conformance-gap.md (command output, diagnostics list, git archaeology).
  fresh_repository_evidence: |
    internal/gate/phases.go BenchkitPhases = gofmt, vet, test, race, system, shellcheck; no conformance phase.
    Only non-test setter of BENCH_CONFORMANCE_ROOT: internal/preprelease/preprelease.go:100 (ship tier).
    `bench gate` at HEAD: green in 64.8 s, footer "capability-skips: 7 (capability=6 environment=1)", environment class never named.
    History correction: 72c037a1 (2026-07-05) MOVED the conformance run from .bench/gate.sh into a Go phase (it added Name:"conformance"); 3701c4a0 (2026-08-09, "adopt branch-native test architecture") REMOVED that phase. So live-root grading was lost 9 days before HEAD, not 43.
    projects/benchkit.md:206 says the separate conformance dev driver was deliberately retired ("There is no separate contract or conformance dev driver…") while line 81 still claims "The gate's conformance layer enforces those contracts"; .bench/BENCH-reference.md:131 still describes "the built-in conformance phase"; bench-review-implementation.md:26 lists "conformance" among what the gate runs. The decided profile expects live conformance inside the ordinary driver; the ordinary driver skips it. Loss of coverage is accidental; the phase removal was deliberate.
    FT120 (occurrence 2026-08-14) and FT213 (2026-08-16) record teams reading the skip as ambiguous after the removal.
  experiment: BENCH_CONFORMANCE_ROOT=$PWD go test -count=1 ./internal/conformance -run '^TestRootConformance$'
  result: FAIL in 3.24 s with exactly the 10 diagnostics Opus listed (stale $bench-finalize-spec refs ×2; bench-implement-spec.md missing Entry orientation / Exit handoff; final-check and BENCH.md dropped implementation-retro owner; BENCH-reference.md:106 removed token "spec build"; canary.go does not consume bounds.CanaryInnerWidth; two dangling Sources paths in decisions/spec-build-review-gate-cadence.md). 10 of 11 phase files still carry the two headings, so the doc drifted, not the check.
  classification: REPRODUCED (with a date/attribution correction to Opus)
  confidence: high
  final_disposition: P0. Restore live-root dev-tier conformance in a way consistent with the profile's decided shape — carry BENCH_CONFORMANCE_ROOT + tier=dev into the ordinary `go test ./...` phase env (or default the root in-kit), not a resurrected separate driver; make an environment-class skip *inside the oracle* red; name every skip. Then disposition the 10 diagnostics (tree fixes, or reviewer-approved contract changes). This is the next ticket (next-ticket.md).

- id: L-02
  topic: Red gate verdicts are never drift-checked; `bench status`/`bench handoff` contradict `bench gate` on the same tree
  sol_position: Not found (E018 notes status action heuristics only).
  opus_position: P0-2 (E5/I.3): inspectSubjectAt returns on non-green before comparing trees; status.GateVerdict defines Stale only for greens; ambient dashboard shows red while oracle reuses green.
  relationship: Opus-only; Sol silent.
  prior_evidence: Opus 01-environment-and-oracle.md E5 (state0–state2 trace).
  fresh_repository_evidence: internal/gate/verdict.go:218-224 (`if rec.Status != "green" { … return }` precedes `if rec.Tree != s.Tree`); internal/status/status.go:237-238 (`nonReusableGreen := … in.Status == "green" && …; gi.Stale = nonReusableGreen`).
  experiment: fresh gate (green) → append gofmt violation to internal/toon/toon.go → `bench gate` red (bench-last-gate: status=red tree=5e44a4ee) → restore file → `bench gate` = "green (fresh verdict reused for this tree)" → `bench status`.
  result: status headline "▶ fix before commit (gate) / gate red → fix before commit" while gate reports green. Direction is fail-safe (never a false green) but the SessionStart surface and the resume artifact assert a red the oracle denies.
  classification: REPRODUCED
  confidence: high
  final_disposition: P0, small (~20 lines + regression test asserting three-surface agreement). Independent of L-01.

- id: L-03
  topic: No canonical entry point / router; `what-next` is misnamed maintenance
  sol_position: Root should be a logical `bench` router over deterministic status + a new canonical work record + context compiler; `/bench`, `$bench` thin adapters (D, Z).
  opus_position: `/bench` as a thin model-invocable phase over a new `bench status --route` flag emitting one `next[1]{state,why,command}` row from existing signals + `specs: staged` + `setup` signals; rename what-next → drain (D.4).
  relationship: Agree on the gap and on thin harness adapters; disagree on whether the router needs a new state layer and compiler.
  prior_evidence: Sol E001/E004/E018/E019; Opus E8/E10, D.2, D.5.
  fresh_repository_evidence: |
    Wrapper no-args: 44-line inventory, rc 0. Binary no-args: "no subcommand", rc 2. `bench help` accepted by wrapper, rejected by binary. `bench review|final-check` → unknown subcommand rc 2.
    `bench status` in a fresh un-adopted git repo → "bench: clean — nothing pending" rc 0 (indistinguishable from adopted-clean).
    internal/status/status.go:546 — the only `specs` signal is "%d merged spec(s) awaiting retirement"; no staged-spec form. intent row action is prose "resume interrupted work" (status.go:342).
    .agents/commands/bench-what-next.md frontmatter disable-model-invocation: true; hidden from this session's skill list (observed at session start).
    Upstream mattpocock/skills@d574778f ships skills/engineering/ask-matt (a prose router, disable-model-invocation: true — user-invoked); not referenced anywhere in Bench.
  experiment: as above (wrapper/binary/status probes; throwaway repo)
  result: gap confirmed on every surface both auditors named.
  classification: REPRODUCED (gap); Sol's need for a new work record + compiler is INFERRED/UNSUPPORTED by observed failure; Opus's projection is COMPATIBLE with existing single-source rules.
  confidence: high (gap), medium (design)
  final_disposition: ADOPT /bench — as a thin adapter over a deterministic `bench status --route` projection (Opus shape), which is also what Sol's "thin harness adapters over one logical entry" reduces to once the unproven work-record/compiler is subtracted. Rename what-next → drain. P1 after L-01/L-02.

- id: L-04
  topic: Work state — introduce a canonical Goal/Core/Verified/Open/Next store vs consolidate existing state
  sol_position: Adopt a small canonical work record (goal, subject, decisions, requirements, verified, open, next, attempts, feedback loop); handoff becomes a generated view (K).
  opus_position: Reject a new layer; every field already has an owner (spec, fences, gate record, maps, status/roadmap/handoff); only "feedback loop" (repro command) is missing; three "Next" sources need one projection (K).
  relationship: Material disagreement on architecture.
  prior_evidence: Sol E011/E012 (stale State survives handoff rewrite); Opus K table.
  fresh_repository_evidence: |
    Owners verified in tree: spec Status/coverage rows (internal/spec, bench coverage), ownership fences (bench preflight), gate evidence record (content-addressed, control), decisions maps (bench maps refuses invalid/unready), intent ledger (bench worktree list; status intent row), handoff pin block (regenerated) + free `## State` prose.
    capture/session-handoff.md at HEAD: pins HEAD 9616919a, says "about to land as one commit" (it landed as 58d966e2 — the commit that wrote the file), "specs/ is empty" (34 tickets-only dirs, 0 spec.md); `bench status` shows no handoff row because age is commit-distance since the file was written.
    No storage anywhere for a bug's repro command / next discriminator (rg "Repro" in internal/handoff and ticket template: none).
  experiment: none beyond inspection (a controlled resume trial was not run).
  result: The distributed state is real, mostly evidence-referenced, and mostly mechanically checked; the two weak spots are prose `## State` and the missing repro/next-discriminator carrier.
  classification: Opus COMPATIBLE; Sol PARTIALLY SUPPORTED (weak spots real; new store not shown necessary).
  confidence: medium-high
  final_disposition: MINIMALLY EXTEND EXISTING STATE — constrain `## State` to tree-contradictable facts; add a `Repro:`/next-discriminator line to the handoff (and blocked-delegate return); let the router be the single Next projection. No new file, no compiler.

- id: L-05
  topic: Context load and whether a context compiler is justified
  sol_position: Bench is a file loader; ~6,000 mandated cold words, 31,456-word estate; needs a compiler emitting a 200–400-token header + ≤1,500 task tokens (L).
  opus_position: 2,342 words always-loaded + 1,259-word skill descriptions ≈ 4.8k tokens; progressive disclosure works; 14,907 words for one spec-backed feature vs Pocock 3,115; only concrete waste is `bench outline` bare (24.6 KB); measure before cutting (L, J.2).
  relationship: Compatible measurements under different definitions; disagree on remedy.
  prior_evidence: Sol E014; Opus J.1/J.2/L.
  fresh_repository_evidence: |
    wc -w: AGENTS.md+.bench/BENCH.md = 2,294; + projects/benchkit.md + CONTEXT.md (which AGENTS/BENCH direct the reader to read) = 7,664; + BENCH-reference + handoff = 10,076; commands 10,991; craft skills+prototype 13,890; references 3,690; total agent-facing prose 27,758 (Opus's exact figure). Sol's 31,456 includes profile/CONTEXT/reference.
    `bench outline` bare = 24,661 bytes.
    Vague-verb probe (carefully/thoroughly/appropriately/robustly/properly) over 27,758 words: 0. Negations in always-loaded 2.3k words: 17.
  experiment: word counts and rg probes above.
  result: Both sets of numbers hold; the difference is "auto-loaded" (Opus) vs "mandated to read cold" (Sol). No observed failure attributable to context volume was produced by either audit or by this run.
  classification: Sol volume figures COMPATIBLE; Sol "context dumping/negative-instruction overload" CONTRADICTED at the measured density; Sol "compiler needed" RESEARCH-INSPIRED/UNSUPPORTED; Opus COMPATIBLE.
  confidence: high (measurements), medium (remedy)
  final_disposition: No first-class context compiler. Fix `bench outline` bare form (P2). Instrument context/tool-call cost (FT138) and run the measurement harness before FT100's editorial cuts.

- id: L-06
  topic: `bench test` / gate runs leak worktree fixtures
  sol_position: E006 — `bench test` left ten untracked `worktrees/001-*` dirs in the subject under the documented audit environment.
  opus_position: Not found.
  relationship: Sol-only.
  prior_evidence: Sol E006 (observation; environment not fully specified).
  fresh_repository_evidence: |
    Pool root = $BENCH_HOME/worktrees/<key> (internal/worktree/worktree.go:51). Sol's pollution of the *subject* implies Sol's BENCH_HOME pointed into the subject; the leak itself is real regardless of where BENCH_HOME points.
    ~/.bench/worktrees holds 759 `001-*` directories (41 MB): 731 dated 2026-08-17 (the two prior audits' gate/test runs), 20 dated 2026-08-18 (this run's two gate runs → ~10 per run). Each holds a fixture repo (tracked.txt, reviewed.txt, README.md) whose t.TempDir origin is gone.
    Owner: internal/worktree/reauthorize_test.go fixtures (reauthorizeFixture → mustCreate) run with no t.Setenv("BENCH_HOME", …), unlike sibling tests that isolate it; internal/systemtest sets its own home.
  experiment: two gate runs; ls -lt ~/.bench/worktrees; rg for "reviewed.txt".
  result: Reproduced and extended: the kit's own tests write into the operator's real pool on every gate run and never clean up.
  classification: REPRODUCED (stronger than Sol's framing)
  confidence: high
  final_disposition: FIX-NOW (small): isolate BENCH_HOME in the worktree package test helper/TestMain; add an assertion that a test run leaves the ambient pool untouched; one-off operator prune of the 759 orphans (out of audit scope).

- id: L-07
  topic: Handoff `## State` prose is stale at the exact commit that reports it fresh
  sol_position: E011/E012 — commit-distance freshness passes while State contradicts the tree; harness rewrite preserves stale prose verbatim.
  opus_position: J.3/S.9 — same; pin block regenerated correctly, prose not; constrain State.
  relationship: Agreement.
  prior_evidence: both.
  fresh_repository_evidence: see L-04 (file content vs tree; no handoff status row).
  experiment: inspection only.
  result: Confirmed. Severity is bounded by the file's own "the tree wins" rule and the regenerated pin block.
  classification: REPRODUCED
  confidence: high
  final_disposition: Fold into L-04 extension (constrain State; add repro line). Not a P0.

- id: L-08
  topic: Wrapper vs binary root/help disagreement
  sol_position: E004/E019 — wrapper prints inventory rc 0, binary "no subcommand" rc 2; help disagrees; schema failure at product boundary.
  opus_position: Not raised (treats `bench help` = inventory as consistent help).
  relationship: Sol-only.
  prior_evidence: Sol E004/E019.
  fresh_repository_evidence: reproduced verbatim (see L-03). The public entry is the wrapper; the binary is plumbing behind route_binary.
  experiment: as L-03.
  result: Real inconsistency, low user exposure.
  classification: REPRODUCED
  confidence: high
  final_disposition: FIX-WHEN-TOUCHED — the router work (L-03) redefines the wrapper's no-arg root; align the binary's no-arg/help then.

- id: L-09
  topic: Review and Final Check are prompt phases, not executable subcommands
  sol_position: E010 — `bench review`/`final-check`/`claims` unknown; "UNREPRODUCIBLE" as controls; merge Final Check into landing/status.
  opus_position: N — final-check reports the retained verdict and captures the retro; review is advisory by design; distinct roles, useful redundancy.
  relationship: Same fact, different framing.
  prior_evidence: both.
  fresh_repository_evidence: `bench review`, `bench final-check` → unknown subcommand rc 2. bench-final-check.md: runs `bench gate`/`bench commit` for non-spec work, reports retained land evidence for spec work, writes capture/retros/<slug>.md with a repair-attribution table (the one quantitative measurement channel Bench has). bench-review-implementation.md: read-only, three axes, writes reviews/<slug>.md pickup.
  experiment: command probes; file reads.
  result: They were never claimed to be commands; the roles are distinct (semantic judgment; close-out + retro). No duplicate oracle exists.
  classification: NOT ACTIONABLE as a defect (Sol's characterization of "claimed control" is UNSUPPORTED; Opus COMPATIBLE)
  confidence: high
  final_disposition: Keep both phases. Small doc fixes from L-16 (below).

- id: L-10
  topic: Codex does not surface the 11 phase adapters; `$bench-debug` cannot self-invoke
  sol_position: E015 — Codex session exposed 17 craft skills, none of the phase adapters, though files exist; parity failure.
  opus_position: H.5/Q — every Codex adapter carries allow_implicit_invocation: false (conformance-checked); Claude has 6 model-invocable phases incl. bench-debug; the strongest-evidence behavior has the weakest trigger on Codex; no Claude-side parity check.
  relationship: Compatible — Sol observed the effect of the policy Opus located.
  prior_evidence: Sol E015; Opus config reads.
  fresh_repository_evidence: all 11 .agents/skills/bench-*/agents/openai.yaml → allow_implicit_invocation: false; the same SKILL.md files also carry the Claude key disable-model-invocation: true (inert on Codex). .agents/commands: 5 phases disable-model-invocation (what-next, update-kit, deepen, assess, setup-repo); bench-debug and 5 others model-invocable — matches this session's visible skill list. checkCodexCommandAdapters (skills_index_checks_test.go:229) is itself a root-conformance check and therefore currently unrun (L-01).
  experiment: config inspection; observed skill list at session start.
  result: The asymmetry is policy, not breakage; its consequence for bench-debug is real.
  classification: COMPATIBLE / REPRODUCED
  confidence: high (config); Codex live behavior UNINVESTIGATED (no Codex session)
  final_disposition: Reviewer decision: make `$bench-debug` implicitly invocable on Codex or accept the repair-loop tripwire (L-22) as the trigger; add the Claude-side invocation-policy parity check; drop the inert key (P1/P2).

- id: L-11
  topic: `diagnosing-bugs` → `/bench-debug`: what is preserved, added, diluted
  sol_position: Mechanism preserved (loop → minimise → hypotheses → discriminate → regression → original-signal rerun); diluted by worktree/shift/line prose; "check it didn't author" wording weakens independence claim; PRESERVE, change delivery, benchmark before replacing.
  opus_position: Four load-bearing constraints (no hypothesis before red loop; exact symptom; one command already run; rerun original loop) verbatim-preserved; five real additions; three compressions (menu 10→5, two hard gates + checkbox form dropped, "tighten the loop" dropped); Codex trigger gap; restore, do not compress further; benchmark.
  relationship: Compatible; Opus more specific.
  prior_evidence: Opus H (clause-by-clause); Sol H.
  fresh_repository_evidence: |
    Upstream pinned d574778f (reference-skill-repos/skills, identical to ~/.claude/skills/diagnosing-bugs, 134 lines) vs .agents/commands/bench-debug.md (140 lines), read in full:
    Preserved: "one command … already run once (paste the invocation and its output)"; "asserts the user's exact symptom … Not 'runs without erroring'"; "If you catch yourself theorizing before this command exists, stop"; "3–5 ranked, falsifiable"; one variable at a time; regression before fix at a correct seam / no seam is the finding; "re-run the full Phase 1 loop"; cleanup + name the hypothesis.
    Added: isolate before first repro artifact; accused-command/no-lookalike rule; seam-vs-edit-owner rule; quarantine marker so the repro survives shift rollback; `bench spec history` for retired specs; reviewer-invoked / delegate blocked-report shape.
    Compressed: menu lists test, curl, CLI invocation, replay, throwaway harness (5 of 10; headless browser, property/fuzz, bisection, differential, HITL absent); "Tighten the loop" section absent (folded into three adjectives); "No red-capable command, no Phase 2" replaced by softer "Don't proceed on a theory"; "Do not proceed until you have reproduced and minimised" and the `- [ ]` completion checklist absent.
    .claude/settings.json turns the personal upstream `diagnosing-bugs` off in this repo (so /bench-debug is the only debug discipline available here).
  experiment: side-by-side read; no live trials (out of scope, self-benchmarking invalid).
  result: Opus's clause list is accurate; Sol's "worktree/shift prose competes for attention" is a fair reading of the 40-line "How it meets the rest of Bench" tail.
  classification: REPRODUCED (textual); outcome effect remains PRACTITIONER EVIDENCE
  confidence: high (text), n/a (effect)
  final_disposition: PRESERVE + STRENGTHEN: restore the menu as a pointer-fired reference file, the two hard gates + checkbox completion form, and "tighten the loop"; move the Bench-mechanics tail behind a pointer; benchmark arms D/E/F before any further compression.

- id: L-12
  topic: Prose quality — vague verbs, negative-instruction overload, duplication
  sol_position: J — negative-instruction overload ("thousands of words"), vague verbs, context dumping, eleven skills + eleven command documents duplicated.
  opus_position: J.1 — 0 vague verbs in 27,758 words; 0.7 negations/100 words; .claude/commands is one symlink; the only real duplication is the inert Codex key.
  relationship: Contradiction on quality; compatible on volume.
  prior_evidence: Sol E014/E015; Opus J.1 probes.
  fresh_repository_evidence: 27,758 words; carefully/thoroughly/appropriately/robustly/properly = 0; 17 negations in the always-loaded 2,294 words; .claude/commands → ../.agents/commands (symlink); the 11 Codex adapter SKILL.md files total 583 words (thin pointers, not duplicated bodies).
  experiment: rg counts; ls -la.
  result: Sol's quality characterizations do not survive measurement; Sol's volume figures do.
  classification: Sol CONTRADICTED (quality), COMPATIBLE (volume); Opus REPRODUCED
  confidence: high
  final_disposition: No prose-quality action; volume handled under L-05/FT100-after-measurement.

- id: L-13
  topic: Release workflow bypasses the governed publication path
  sol_position: E020, P0 — release.yml publishes with raw `npm publish`; runbook requires `bench release submit/promote`.
  opus_position: Not exercised (release-adjacent commands deliberately not run).
  relationship: Sol-only.
  prior_evidence: Sol E020; ASSESSMENT.md §6 (2026-08-13) records the same.
  fresh_repository_evidence: .github/workflows/release.yml:68,73 `npm publish … --provenance`; docs/release-runbook.md:35,47,58,70 require submit/promote; ROADMAP.md release-readiness item 4 ("Publication is staged, resumable…") and status NO-GO; no FT row owns this.
  experiment: static read only (no release command run).
  result: Confirmed statically.
  classification: REPRODUCED (static)
  confidence: high
  final_disposition: P1, before the next release (deployment is already NO-GO): make submit/promote the only publish path in the workflow; add a row (no FT owns it today).

- id: L-14
  topic: A newly adopted repo's gate cannot go green
  sol_position: Not found (D: init should be transactional; no reproduction).
  opus_position: I.7/E9 — setup scaffolds no gate-inputs.json → `HOME: unbound variable` in .bench/bin/bench.sh:14 under set -u; declaring HOME still red because `bench canary` exits 1 on an empty fixture inventory and the scaffolded gate hard-fails on it.
  relationship: Opus-only.
  prior_evidence: Opus E9 (throwaway repo).
  fresh_repository_evidence: reproduced in a throwaway repo under the scratchpad: `bench setup --yes` (plan-first, self-verifying, next: replace sentinel) → remove sentinel → `bench gate` = "line 14: HOME: unbound variable / canary inventory validation failed / gate: red" → add gate-inputs.json declaring HOME → `gate --fresh` = "canary fixture inventory is empty / gate: red". No .bench/gate-inputs.json scaffolded.
  experiment: as above (~90 s; repo removed at close).
  result: Reproduced end to end.
  classification: REPRODUCED
  confidence: high
  final_disposition: P1 — scaffold gate-inputs.json; wrapper names its own failure on unbound HOME; scaffolded canary line commented with TODO; adoption smoke fixture in the gate.

- id: L-15
  topic: Un-adopted repo reports "clean — nothing pending"
  sol_position: Not raised as such.
  opus_position: E8/D.2(1).
  relationship: Opus-only.
  prior_evidence: Opus E8.
  fresh_repository_evidence: reproduced (fresh git repo, no .bench → rc 0 "clean — nothing pending").
  experiment: as above.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: fold into router (`setup` signal when .bench/ absent).

- id: L-16
  topic: Prose contradictions in review/reference docs
  sol_position: E017 — review says commit pickup state and "runs no gate" while `bench commit` gates; "the gate and you" decide done vs gate-only authority.
  opus_position: J.3 — BENCH-reference.md:106 stale "spec build" token; implement-spec headings; final-check/BENCH.md retro-owner drift (all conformance-detected).
  relationship: Compatible (different instances).
  prior_evidence: both.
  fresh_repository_evidence: bench-review-implementation.md:26 lists "conformance" among what the gate runs (false at HEAD); :29 "The gate and you do"; :86 commit pickup state, :109 "runs no gate" (a raw `git commit -m` is guard-allowed, so the tension is wording, not a trap). BENCH-reference.md:131 describes a built-in conformance phase that no longer exists.
  experiment: reads.
  result: Real but low-impact wording drift; three of the instances are exactly what root conformance flags (L-01).
  classification: PARTIALLY SUPPORTED (Sol overstates the authority contradiction; "you" is the reviewer in Bench's voice)
  confidence: medium
  final_disposition: FIX-WHEN-TOUCHED under L-01's 10-diagnostic clean-up plus a one-line clarification of "who decides done".

- id: L-17
  topic: 34 tickets-only spec directories with no owner
  sol_position: E016 — ~40 ticket files, ticket-only dirs invisible to spec status; coverage "spec not found".
  opus_position: I.6/S.2 — 34 dirs, 43 files, oldest 2026-08-02; spec retire correctly refuses; no status signal.
  relationship: Agreement.
  prior_evidence: both.
  fresh_repository_evidence: ls specs → 34 dirs, 0 spec.md, 43 files.
  experiment: find/ls.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: DELETE the residue after a reviewer pass; add a `specs: N tickets-only` status row so the debt is visible (P2 hygiene).

- id: L-18
  topic: Advisory reds nobody must run — invalid decision map, structure 62 with stale accept rows
  sol_position: E003/E009 — 62 structure issues, one invalid map coexist with a green gate.
  opus_position: I.5/S.3 — same; decision-map-integrity is a conformance check that doesn't run; structure row is alarm fatigue.
  relationship: Agreement.
  prior_evidence: both.
  fresh_repository_evidence: `bench maps` rc 1 (spec-build-review-gate-cadence.md Sources → deleted internal/specbuild/*); `bench structure` rc 1, 62 issues, 4 "stale accept row" lines in .bench/structure-accept.
  experiment: commands run.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: map: fixed by L-01 wiring + one Sources edit; structure: reviewer decision (gate at a frozen budget or drop from status); stale accept rows FIX-NOW.

- id: L-19
  topic: `bench outline` bare form emits ~24.6 KB
  sol_position: Not raised (F: "Keep observation-only").
  opus_position: F.2/S.1 — 200 of 5,335 symbols, no ranking rule.
  relationship: Opus-only.
  prior_evidence: Opus E15.
  fresh_repository_evidence: `bench outline` → 24,661 bytes.
  experiment: run.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: P2 — bare form emits meta + per-directory counts; rows need a path or --full.

- id: L-20
  topic: Git guard — degraded rim false positives; live deny-surface gaps
  sol_position: Blocked read-only text during cold audit; keep, narrow parser, guarantee bootstrap path.
  opus_position: E2/E13 — degraded rim matches substring `git`; live gaps stash drop/clear, filter-branch, rm -rf; restore --staged correctly allowed.
  relationship: Agreement.
  prior_evidence: both.
  fresh_repository_evidence: at session start (no core) `cat .bench/hooks/block-dangerous-git.sh | head` was BLOCKED; with core: guard-git allows `git stash drop`, `git stash clear`, `git filter-branch --all`, `git rm -rf .`; blocks `git reset --hard`, `git restore .`; allows `git restore --staged .`.
  experiment: 8 envelope probes through `dist/bench guard-git`; the session-start block itself.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: P2 — tokenize the degraded rim; add the three gap spellings.

- id: L-21
  topic: Gate skip disclosure aggregates the dangerous class
  sol_position: E008 — "seven capability skips" (not attributed).
  opus_position: I.2 — the single environment skip is the whole conformance suite; reasons are held in memory (environmentReasons) but not printed.
  relationship: Compatible (Sol didn't decode it either — the same trap the 2026-08-13 assess fell into).
  prior_evidence: both.
  fresh_repository_evidence: gate footer "capability-skips: 7 (capability=6 environment=1)" then class=fifo: 3, class=privilege: 3; internal/gate/capability_skips.go keeps environmentReasons and skipRows never prints them.
  experiment: gate run; code read.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: fold into L-01 (name the skipping test + reason; environment skip inside the oracle → red).

- id: L-22
  topic: Nothing detects a repair loop / blank retry
  sol_position: Attempt ledger + blank-retry tripwire in work state (I, T.4).
  opus_position: I.4 — advisory `repair_loop` row from gate records on the second consecutive red with the same first-failing check; measure against the retro repair table before making it blocking.
  relationship: Compatible; Sol wants a schema, Opus a projection.
  prior_evidence: neither ran it.
  fresh_repository_evidence: .logs/gate-*.jsonl carries per-phase `phase.finish` with exit and elapsed per run (phase granularity, not check granularity); bench-last-gate carries only status/tree/oracle/time; retro template mandates a repair-attribution table.
  experiment: log inspection.
  result: A phase-level tripwire is computable from existing records; check-level needs more.
  classification: HYPOTHESIS (benefit) over OBSERVED (data exists)
  confidence: medium
  final_disposition: EXPERIMENT — advisory row first; measure repair rounds via the retro tally.

- id: L-23
  topic: Review independence
  sol_position: Unverifiable; needs an isolated review runner with retained inputs.
  opus_position: Real by construction (parallel fresh axes, re-derive-then-compare) except the diff arrives via `bench diff --full`, which includes the commit log; arm G settles it.
  relationship: Disagreement on degree.
  prior_evidence: neither ran trials.
  fresh_repository_evidence: `bench diff --help`: "--full also appends the commit log and verbatim tracked patch"; review phase step 3 feeds delegates the frozen diff.
  experiment: none possible without model trials.
  result: The one concrete leak is the commit log.
  classification: PARTIALLY SUPPORTED (both)
  confidence: medium
  final_disposition: EXPERIMENT (arm G); cheap change candidate: hand review delegates the patch without the log.

- id: L-24
  topic: Claims/evidence — control system or reporting system?
  sol_position: Needs a typed claim/evidence/freshness graph extending gate-record semantics (M).
  opus_position: Bench has no claim system and is right not to; the gate verdict is one content-addressed fact; three of four red-capable advisory surfaces control nothing (M).
  relationship: Disagreement.
  prior_evidence: Sol E007/E008; Opus M attack table.
  fresh_repository_evidence: control confirmed for gate green/red/stale (commit refuses, Stop blocks under BENCH_SHIFT, reuse refused); monitoring-only confirmed for structure, maps, conformance (all red at HEAD, gate green). Evidence store keyed sha256(tree‖oracle) — reuse verified in this run.
  experiment: gate reuse/red/revert sequence (L-02).
  result: The observed defect is that some deterministic checks are monitoring when they should be control (L-01), not that semantic claims lack a graph.
  classification: Opus COMPATIBLE; Sol UNSUPPORTED (RESEARCH-INSPIRED)
  confidence: high
  final_disposition: No claim graph. Convert the specific monitoring surfaces to control (L-01) and fix the reader (L-02).

- id: L-25
  topic: Delegation — fixed agent counts vs justified boundaries
  sol_position: Six assess agents, three review axes, fresh writer per ticket, three alternative designers are unmeasured ceremony; delegate only when the graph merits it.
  opus_position: Every boundary buys isolated ownership, independent rediscovery, or bounded research; nothing enforces fan-out/effort/caps.
  relationship: Compatible on "unmeasured", disagree on default.
  prior_evidence: neither measured.
  fresh_repository_evidence: craft-delegate/craft-review/tickets texts as both describe; check-agent-line enforces the model tier only; FT138 (instrument build economics) open.
  experiment: none.
  result: Unmeasured either way.
  classification: UNINVESTIGATED → BENCHMARK
  confidence: n/a
  final_disposition: No change; leave-one-out ablation in the measurement harness.

- id: L-26
  topic: Continue incrementally vs consolidate vs strangler vs rewrite
  sol_position: Strangler: new router + work record + compiler + typed evidence in front of existing modules.
  opus_position: Incremental with one bounded consolidation (verdict reader; status as single Next projection).
  relationship: Material disagreement.
  prior_evidence: Sol R table; Opus R table.
  fresh_repository_evidence: every defect reproduced in this run is local (an env var, an if-order, a missing status row, a scaffold file, a test helper); the deterministic core (tree-keyed evidence, prospective commit, fences, coverage --check, guards) held under every probe here and in both audits.
  experiment: the sum of L-01..L-25.
  result: No observed failure requires a new spine; a spine would duplicate readers Bench's one-source rule forbids.
  classification: Opus COMPATIBLE; Sol INFERRED
  confidence: high
  final_disposition: A — continue incrementally with the bounded consolidations in L-02/L-03/L-04.

- id: L-27
  topic: Final Check vs gate — complementary or duplicative
  sol_position: Merge (no independent executable differential).
  opus_position: Keep (reports, never re-grades; captures the retro).
  relationship: Disagreement.
  prior_evidence: both read the file.
  fresh_repository_evidence: see L-09.
  experiment: reads.
  result: Complementary: final-check is the invoker/reporter + retro writer; the gate is the oracle.
  classification: Sol UNSUPPORTED; Opus COMPATIBLE
  confidence: high
  final_disposition: Keep.

- id: L-28
  topic: Test/harness state visible in this run that neither audit weighted — FT212 already fixed
  sol_position: —
  opus_position: —
  relationship: —
  prior_evidence: none.
  fresh_repository_evidence: `bench worktree clean --landed` → valid empty plan rc 0; usage now lists `(<path> | --landed)`; FT216 landed 7d799ecc.
  experiment: run.
  result: Roadmap row FT212 is stale.
  classification: STALE (roadmap row)
  confidence: high
  final_disposition: DONE/ARCHIVE FT212.

- id: L-29
  topic: SessionStart is silent when the core is missing
  sol_position: E001 (127 before build).
  opus_position: I.9 — hook exits 127 with no build hint; recovery is 1.4 s and documented only in the handoff.
  relationship: Agreement.
  prior_evidence: both.
  fresh_repository_evidence: this session opened with no Bench dashboard and no hint (the visible dashboard was the user's global gl-axi hook); session-start.sh execs `session-inspect` into the missing binary.
  experiment: session start itself.
  result: Reproduced (weakly — the harness may swallow the hook's stderr).
  classification: REPRODUCED
  confidence: medium
  final_disposition: P2 — on 127, print the one-line build/repair command.

- id: L-30
  topic: Duplicate invocation-policy key on Codex adapters; Claude-side parity check missing
  sol_position: adapters duplicated (11+11).
  opus_position: only real duplication is the inert Claude key on Codex SKILL.md; add the Claude mirror check.
  relationship: Compatible.
  prior_evidence: both.
  fresh_repository_evidence: all 11 Codex SKILL.md carry disable-model-invocation: true and openai.yaml allow_implicit_invocation: false; no Claude-side check in internal/conformance.
  experiment: rg.
  result: Reproduced.
  classification: REPRODUCED
  confidence: high
  final_disposition: DELETE the inert key; add the parity check (P2), decide per phase which are model-invocable on both.

- id: L-31
  topic: no-mistakes / AXI import revisions unrecorded (Sol E022)
  sol_position: Provenance gap for two of three upstreams.
  opus_position: Pins its own local checkouts; README provenance "substantially accurate".
  relationship: Compatible.
  prior_evidence: Sol E022.
  fresh_repository_evidence: not re-derived (README read only for ask-matt).
  experiment: none.
  result: —
  classification: UNINVESTIGATED / NOT ACTIONABLE (prospective fix only)
  confidence: n/a
  final_disposition: Record source revision + delta on future imports; no retro-archaeology.
```
