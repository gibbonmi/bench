# Prioritized roadmap fixes

Assessed 2026-08-13 against `main` at `287bfbf9`. This is a triage view of all
62 `FT` rows in `ROADMAP.md`; it does not replace the roadmap or authorize a
build. Each row must still be revalidated at build entry because the live tree,
not this snapshot, is authoritative.

## Reading the queue

- **Applicable** means the current tree still exhibits the stated gap.
- **Partial** means some named work shipped or disappeared; re-slice the live
  remainder before implementation.
- **Decision** means the opportunity still exists but no build is authorized
  until the named reviewer choice is closed.
- **Retire** means current code already satisfies or supersedes the row.
- **Parked** means the graduation condition is not currently met.

Priority favors a trustworthy oracle and false-green removal, then lifecycle
correctness, then doctrine and ergonomics. Dependency order is binding inside
the list: FT98 precedes FT169; FT89 precedes FT100 and FT108; FT108 precedes
FT186; FT113 and FT98 precede FT166; FT169 precedes FT71 and FT162.

## Current focus

FT153, FT171, FT175, and FT203 no longer appear in the live roadmap. Canary
inventory and planted-reason proof now have separate truthful owners; the serial
gate already meets its destination; the claim-ledger proposal was retired; and
the worktree-list fixture flake now has deterministic regression coverage.

The current roadmap sequence is FT198, then FT189. Shape the progressively
loaded roadmap owner first; next, decide the Bench-owned refusal or execution
bound for the reproduced upstream `git worktree list` hang.

## Ordered queue

| Rank | ID | Verdict | Current assessment and next disposition |
| ---: | --- | --- | --- |
| 1 | FT198 | Applicable | `ROADMAP.md` remains a large, eagerly loaded board. Shape the durable detail owner, migration, history, and index-completeness contract now. |
| 2 | FT189 | Applicable | An upstream `git worktree list` hang can still stall every Bench command that performs discovery; decide the Bench-owned refusal or execution bound next. |
| 3 | FT142 | Partial | Release publication still bypasses the governed resumable publisher. Revalidate the surviving FT91 residuals, then route publication through the current owner. |
| 4 | FT133 | Applicable | `bench coverage --check` still validates the acceptance map without proving ancestry, ticket fence completeness, or one accountable producer per row. |
| 5 | FT174 | Applicable | Ticket metadata can still produce dependency or ownership false-greens; enforce the declared graph and exact producer fence before expanding lifecycle automation. |
| 6 | FT177 | Partial | Branch-native build/test narrowed the stale-`dist/bench` surface, but direct contract mutation probes can still select an old binary. Reproduce the remaining public path before editing. |
| 7 | FT103 | Applicable | `checkNoKitOnlyPackedAssets` checks exclusions and a non-empty allowlist, but does not prove every allowlisted source path exists. |
| 8 | FT201 | Applicable | Production cancel-signal registrations still need one conformance rule and mutation that detects an omitted cleanup path. |
| 9 | FT98 | Applicable | Preserve-then-discard behavior remains split across lifecycle commands. Establish one primitive before FT169 and FT166 consume it. |
| 10 | FT169 | Applicable | There is still no single sanctioned worktree landing command that owns preservation, integration, and cleanup. Blocked by FT98. |
| 11 | FT162 | Applicable | Full-run and phase-close state can still name different subjects. Make the landed candidate the single authority after FT169 settles landing ownership. |
| 12 | FT141 | Partial | The output-truncation guidance face is closed, but gate-pin red records still lack failing-phase inventory and stable attribution. Re-slice to those live gaps. |
| 13 | FT178 | Applicable | Bare `bench worktree` remains human porcelain with behavior that is too implicit for automation; make the public contract explicit. |
| 14 | FT199 | Applicable | Recovery-aware branch retirement is still spread across commands and cleanup paths; one coordinator remains warranted after landing ownership is settled. |
| 15 | FT173 | Decision | The residual active-assignment/missing-tree class remains undecided. The currently present retained worktree does not close the missing-tree policy choice. |
| 16 | FT202 | Decision | The standing test-support fence and process-backed fixture census scope remain reviewer choices. Existing `axitest` and `gittest` owners outrank an invented commons where they fit. |
| 17 | FT190 | Partial | Injected-port conformance exists for three packages only. Re-slice the row to uncovered packages and require either a real-producer test or an explicit exemption. |
| 18 | FT185 | Applicable | Gate results still do not fully participate in the structured Bench output contract; make machine-readable failure details first-class. |
| 19 | FT92 | Applicable | Subject-drift attribution and shipped-input hygiene remain incomplete across consumers of retained and landed state. |
| 20 | FT89 | Applicable | The generated skills index exists, but YAML parsing, stale-reference detection, examples, and broader current-state documentation remain incomplete. |
| 21 | FT106 | Applicable | Documentation claims are not comprehensively reverified against live code. Build the claim check after FT89 identifies the canonical documentation surface. |
| 22 | FT192 | Applicable | The one-source-per-fact rule is enforced in code more strongly than in spec and ticket prose; duplicated policy can still drift there. |
| 23 | FT130 | Applicable | A capture write during an active lifecycle can still invalidate the subject without a mechanical void-or-block response. |
| 24 | FT164 | Applicable | Repair-lane charges still need a done-claim that resolves the named failure rather than merely producing a green adjacent check. |
| 25 | FT200 | Decision | Landing preflight is still partly procedural. Decide which freshness, ancestry, ownership, and exact-candidate checks become one mechanical preflight. |
| 26 | FT117 | Applicable | The two FT87 parser-surface leaves remain useful focused fixes; keep them independent from broader parser redesign. |
| 27 | FT179 | Partial | `craft-comments` and `craft-review` now prohibit reviewer-facing provenance, but the codebase sweep and high-stakes comment repair remain. |
| 28 | FT111 | Partial | Its standalone premise has narrowed to the same surviving comment debt as FT179. Fold it into FT179 rather than run a second provenance-tag project. |
| 29 | FT172 | Partial | `bench-what-next` now captures one full trusted snapshot, closing the truncation face. Parser grammar and discrepancy coverage still need work. |
| 30 | FT144 | Applicable | Kit specs still serve builders and reviewers without an explicit two-audience discipline; clarify the contract before another spec-template expansion. |
| 31 | FT158 | Applicable | Cross-harness guidance still needs one canonical derivation and conformance check so harness-specific adapters cannot drift. |
| 32 | FT71 | Applicable | Versioned local shift remains valuable on the bank track, but it should consume the sanctioned landing path rather than invent another one. Blocked by FT169. |
| 33 | FT191 | Applicable | Charges still lack a cheap fixture-and-seam inventory that makes intended test placement explicit. |
| 34 | FT113 | Applicable | `bench commit --spec` still refuses an empty path list and does not treat the spec-state flip as its own path; close this before capture porcelain depends on commit semantics. |
| 35 | FT166 | Partial | Capture porcelain remains useful, but the proposal assumes older reduced-gate behavior. Re-shape it around current branch-native exact-verdict policy after FT98 and FT113. |
| 36 | FT168 | Applicable | Baseline canary meaning is settled. Focused iteration still lacks a fixture or family selector and must remain distinct from the full gate oracle. |
| 37 | FT182 | Applicable | A Planned-phase receipt over an absent target can still wedge the lifecycle; reject or repair that state mechanically. |
| 38 | FT100 | Applicable | Guidance remains heavier than necessary in places. Run the prose-weight pass after FT89 establishes an accurate inventory. |
| 39 | FT108 | Applicable | There is still no first-class refactor lane with a mechanical exit test. Define it only after the guidance surface is coherent. |
| 40 | FT186 | Applicable | `executeSubjectWithEngine` and `readyFieldClasses` remain deep mixed-responsibility seams. Split and type them after FT108 defines the refactor contract. |
| 41 | FT99 | Applicable | Spec compilation still needs an explicit live-code check that falsifies a stale problem premise before slicing implementation. |
| 42 | FT102 | Applicable | `craft-synthesis` still lacks an explicit escalation-policy consistency check and a decomposed dogfood pass. |
| 43 | FT112 | Applicable | `bench-debug` still does not say that a green approximation is evidence about the approximation, not clearance of the original bug. |
| 44 | FT205 | Applicable | `craft-delegate` names ordinary release, but not the exact clean-to-resume-clean recovery pair for an interrupted delegate worktree. |
| 45 | FT58 | Applicable | Pool-root selection still needs hardened ownership and permission handling instead of assuming a writable ambient directory. |
| 46 | FT104 | Applicable | Load-induced commit refusals still need a stop rule and a cheap pre-gate discriminator so transient pressure is not misreported as a product defect. |
| 47 | FT115 | Partial | The original named wait sites are gone and marker-wait bounds now exist. Re-census only surviving unexplained deadline literals; do not rebuild the retired examples. |
| 48 | FT120 | Partial | Most named harness files disappeared with branch-native and lifecycle subtraction. Reconcile the row to live FIFO, timeout, and isolation defects before implementation. |
| 49 | FT94 | Partial | The duplicated resume summary has fallen from the original four sites to two test expectations. Decide whether their independence is required before extracting shared knowledge. |
| 50 | FT101 | Applicable | Monorepos still need per-context profile and domain-document scope, but this is an extension after current single-context correctness work. |
| 51 | FT138 | Applicable | Bench still lacks enough stable instrumentation to price build economics; add measurement only after the higher-priority demand reductions. |
| 52 | FT125 | Applicable | `bench outline` has no symbol selector and there is no `bench spec show`; reader surfaces still over-return context for narrow queries. |
| 53 | FT180 | Partial | A spec-optional route remains a real product question, but the proposed mechanics predate lifecycle subtraction. Re-shape from current public commands before deciding. |
| 54 | FT204 | Decision | A bounded transcript/session query could help diagnosis, but its authority, retention, and output contract require a reviewer choice before implementation. |
| 55 | FT140 | Decision | The remaining review residuals need explicit verdicts rather than an assumed build; resolve them in a grill after correctness work above. |
| 56 | FT170 | Decision | Behavioral red/green evaluation for prose guidance remains plausible, but only a narrow benchmark with held-out evaluation can justify adopting it. |
| 57 | FT38 | Decision | The dashboard identity pass is still optional product work. Keep it behind an explicit visual-direction decision. |
| 58 | FT165 | Retire | `craft-domain` exists, is charged by `bench-shape-idea`, and maintains `CONTEXT.md`; the row's requested discipline is already shipped. |
| 59 | FT197 | Retire | Public gate invocation and process lifetime are owned by the Go `gate-run` path, including process groups, signals, and terminal records. |
| 60 | FT6 | Parked | No second reproduced `bench commit` refusal establishes the suspected inherited-verdict defect. Leave parked pending the stated evidence. |
| 61 | FT24 | Parked | The roadmap's external trigger has not fired. Recheck the official hook surface only when its documented behavior changes. |
| 62 | FT8 | Parked | The scheduled Sonnet 5 revisit date is 2026-09-01; it is not actionable on 2026-08-13. |

## Retirement and re-slice boundary

The two **Retire** rows should be removed from the next canonical roadmap
drain, not implemented. Every **Partial** row needs a fresh problem statement
and exact live red signal before ticketing. Do not carry vanished file names,
old lifecycle assumptions, or already-closed clauses into a build fence.
