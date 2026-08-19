# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-drain` reconcile removes it. Raw capture never lands here — it
goes to `capture/IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>/spec.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

FT-to-FT ordering is single-sourced in `## Dependencies`. A literal dependency
blocks the dependent row until its prerequisite lands. A recommended dependency
does not block work; it says the dependent row will be cheaper or less likely to
churn when it is specified after the prerequisite.

## Release-readiness status

Deployment remains **NO-GO** for other projects, public npm release, and bank
use.
The 2026-07-11 release-readiness and repository-controlled compliance
assessments are evidence snapshots; this roadmap is the execution source for
their active findings. Each finding still open appears exactly once on a
`Sources` line below. `RR:` means the release-readiness assessment and `RC:`
means the repository-controlled compliance assessment.

## Next — the 2026-08 capability audit's active portfolio

The audit's remaining action items, now board rows rather than an external
pointer. A1, A2, and A3 landed; these four are what is left of the portfolio
that is actionable today. The `## Recommended sequence` draws from here first.

**FT226 (HIGH) — kit tests stop writing into the operator's real `BENCH_HOME`.**

**FT227 (HIGH) — adoption smoke: a newly adopted repository's scaffolded gate goes green.**

**FT228 (HIGH) — `/bench-debug` recovers its upstream constraints and settles its trigger on both harnesses.**

**FT229 (MEDIUM) — hygiene batch: residue, stale prose, and small CLI sharp edges.**

## Features, in priority order

**FT234 (HIGH) — a pool key whose source repository is gone can be reclaimed.**

**FT71 (HIGH on the bank track) — versioned local shift evidence.**

**FT162 (MEDIUM, decision required) — the handoff's pins resolve, and its shape has one owner.**

**FT142 (MEDIUM) — FT91 review residuals: eight open findings, two tracks.**

**FT144 (MEDIUM) — kit specs have two audiences and the discipline names neither.**

**FT158 (MEDIUM) — make cross-harness falsification standing for kit-guidance diffs.**

**FT98 (MEDIUM) — one preserve-then-discard primitive; four faces.**

**FT169 (MEDIUM) — the landing command's refusals name what to do next, and its authority is decided.**

**FT207 (MEDIUM, decision required) — worktree-mutating paths share malformed-admin refusal.**

**FT224 (MEDIUM) — a lifecycle refusal names the component that failed and the command that fixes it.**

**FT225 (MEDIUM, decision required) — reviewed landing has no path for a spec the review amended.**

**FT199 (MEDIUM) — a recovery-aware branch-retirement coordinator closes one repository-wide ref inventory.**

**FT233 (MEDIUM) — landing refusals say what to do next.**

**FT178 (MEDIUM) — `bench worktree`'s bare verb is a human porcelain that traps automation and leaks on signals.**

**FT172 (MEDIUM, decision required) — the roadmap row grammar is a contract, and `roadmap_id` has one decided source.**

**FT173 (MEDIUM, decision required) — AXI residual: the active-assignment-with-deleted-tree disclosure class.**

**FT202 (MEDIUM, decision required) — a standing test-support fence, and the census scope for process-backed fixtures.**

**FT185 (MEDIUM) — gate results join the structured Bench output contract.**

**FT89 (MEDIUM) — guidance coherence and current-state documentation.**

**FT106 (MEDIUM) — doc claims re-verified against the tree.**

**FT191 (MEDIUM) — a fixture-and-seam inventory a charge can carry for free.**

**FT213 (MEDIUM) — `craft-delegate` visit: worktree isolation, end-of-life, and probes that prove something.**

**FT214 (MEDIUM) — `craft-spec` visit: fences the build may not move, one source per fact, and rows that prove what they claim.**

**FT208 (MEDIUM, decision required) — skills-index producer-hardening residuals: one refusal grammar, per-shape marker diagnostics, and HI14's seam.**

**FT204 (LOW, decision required) — one bounded transcript/session query.**

**FT58 (LOW) — hardened pool roots.**

**FT92 (LOW) — attributed subject drift and consumer-shipped input hygiene.**

**FT99 (LOW) — spec problem-premise verification.**

**FT215 (MEDIUM) — no changed-package-scoped gate path; every diff pays the full fixed-cost floor.**

**FT217 (LOW) — one decision every adopt-lifecycle verb executes.**

**FT218 (LOW) — named git readers instead of learned CLI flags.**

**FT219 (LOW) — `/bench-deepen` refreshes a ready map's frontier to current state before handoff.**

**FT220 (LOW) — `/bench-write-spec` censuses shared decision readers before ticket slicing.**

**FT222 (LOW, decision required) — a standing per-repair-class delegate-tier preference in `projects/benchkit.md`.**

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.**

**FT101 (LOW) — per-context scope for monorepos: domain docs and profile.**

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency and dogfood loops.**

**FT108 (LOW) — a refactor lane with a mechanical exit test.**

**FT113 (LOW) — `bench commit --spec` residuals: the flip counts as a path, and the flip has one author.**

**FT130 (MEDIUM, decision required) — a capture write mid-lifecycle voids or blocks the run.**

**FT164 (MEDIUM) — repair-lane charges, and a done-claim that resolves its named owners.**

**FT200 (MEDIUM, decision required) — make preflight mechanical at the landing chokepoint.**

**FT182 (LOW) — a Planned-phase receipt over an absent target wedges the abandon retry.**

**FT166 (LOW) — `bench capture commit`: porcelain for the ambient capture set.**

**FT168 (LOW) — focused iteration evidence: a fixture-selecting canary.**

**FT140 (LOW) — review residuals that want a verdict, not a build.**

## False greens — verdicts that credit unchecked work

Five rows, one failure class: a green whose warrant is missing — a stale
binary, a dead or skipping citation, a vacuous baseline, an unchecked absence,
a dependency edge nothing resolves. Each hardens a different oracle surface, so
they stay separate builds, but they read and prioritize as one theme.

**FT133 (MEDIUM) — `bench coverage --check` verifies that red-signal citations resolve.**

**FT174 (MEDIUM) — ticket files have one enforced dependency, ownership, and mutation grammar.**

**FT177 (MEDIUM) — a stale `dist/bench` makes contract-test mutation probes silent no-ops.**

**FT103 (LOW) — existence-checked absence evidence: the gate half.**

**FT201 (LOW) — production cancel-signal registrations conform to one source.**

## Reds the diff doesn't own — inheritance, load, and harness defects

Five rows, one failure class: a red that answers for something other than the
diff in front of the gate — an inherited baseline, machine contention and a
flaky oracle, a literal deadline, a harness defect.

**FT141 (MEDIUM, decision required) — red verdicts are recorded against a baseline, so inherited reds stop reading as caused.**

**FT223 (LOW, decision required) — `bench commit`'s inherited-verdict refusal misreads as a red gate.**

**FT104 (LOW) — load-induced commit refusals: the stop rule and the pre-gate quiet check.**

**FT115 (LOW) — load-robust test and phase deadlines derived from bounds.**

**FT120 (LOW) — gate, canary, and contract test-harness defects nothing asserts.**

## Standards debt — one batched light-path pass

Three rows plus FT142's standards track are shippable together as small
one-source-per-fact and cleanup sweeps under one gate; FT117's parser-routing
half is the largest item in the batch. FT142 itself stays on the main list
because its ship track belongs to a separate `prep-release` hardening visit.

**FT117 (MEDIUM) — FT87 parser-surface follow-ups.**

**FT179 (MEDIUM) — comment quality: strip the reviewer-facing register, document high-stakes surfaces, sharpen `craft-comments`.**

**FT94 (LOW) — single-sourced `bench resume` summary golden.**

## Session tax — evidence-supplied reader rows

This row is a measured, recurring reader cost from the week-of-2026-07-19
transcript evidence and builds on surfaces that already exist.

**FT125 (LOW) — reader surfaces that return the slice, not the file.**

## Release and bank reassessment gate

The row below blocks the next release outright; the numbered conditions are the
reassessment gate itself.

A green source-tree gate is necessary but not sufficient. Reassessment attaches
to one immutable version and its generated manifest after:

1. FT82 has an executable regression contract and is closed (FT79–FT81 shipped
   with their regression contracts).
2. The same commit passes the full gate, race tests, vet, canary, vulnerability
   scan, package inspection, reproducibility comparison, and clean-room
   installed smokes.
3. Exact indexed artifacts select the right binary on every supported target,
   work offline, and agree on version, tag, changelog, commit, toolchain, and
   digest.
4. Publication is staged, resumable, digest-verified, wrapper-last, and bound
   to the repository-owned evidence bundle.
5. Setup, doctor, relink, fresh clone, an operational command, and unlink pass
   from an isolated prefix without a source checkout; setup preserves existing
   instructions, settings, and hooks and is idempotent and reversible (shipped
   with FT76; re-verified at reassessment).
6. Consumer artifacts exclude maintainer-only capabilities and include the
   supported-platform, security, data-handling, threat, support/EOL, network,
   rollback/recovery, license/notice, SBOM, checksum, and package-inventory
   records.
7. Bank evidence includes redacted local events for success, failure,
   interruption, and recovery, plus the preservation, oracle-change,
   post-agent recovery, unknown-hook, environment-minimization, offline, and
   transactional lifecycle contract results.
8. A clean-room user can complete setup from the README in one shell command
   plus at most one harness-native conversation.

Host IAM, OS sandboxing, endpoint controls, firewalls, server-side branch
protection, central CI administration, SIEM/retention, registry administration,
and signing-key custody remain outside this repository-controlled roadmap.

**FT230 (MEDIUM, release track) — the release workflow publishes only through `bench release submit/promote`.**

## Parked and scheduled work

**FT6 (LOW, parked pending evidence — leave parked):**

**FT24 (parked pending upstream) — Codex agent-line guard parity.**

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.**

**FT38 (LOW, decision required) — dashboard visual identity pass.**

**FT197 (LOW) — the wrapper's gate hop.**

**FT231 (EXPERIMENT, decision required) — a measurement harness and the instrumentation it reads.**

**FT232 (EXPERIMENT, decision required) — repair-loop tripwire: an advisory signal from gate records.**

## Dependencies

The dependent FT is named first. Only the literal table blocks work; the
recommended table is sequencing advice.

### Literal

| FT | Depends on | Why |
|---|---|---|
| FT100 | FT231 | Editorial cuts to the guidance surface without measurement are coin flips; the harness is what tells a cut from a regression. |

### Recommended

| FT | Better specified after | Why |
|---|---|---|
| FT71 | FT169 | The event schema should record the settled landing and recovery lifecycle rather than version an interim one. |
| FT100 | FT89 | Cut prose after the correctness and coherence pass establishes which guidance is still authoritative. |
| FT108 | FT89 | FT89 single-sources the skills index the new skill must join; the expand–migrate–contract and gate-cadence rules it builds on are already settled in `craft-tickets`. |
| FT172 | FT106 | Reuse the document-claim probe for semantic roadmap claims instead of designing a second checker. |
| FT162 | FT169 | Build full-run subject resolution on the settled landing primitive. |
| FT166 | FT98, FT113 | The porcelain composes over the shipped reduced-gate path allowlist; recoverable set-aside then defines the commit command's smallest sound contract. |
| FT169 | FT98 | Reuse recoverable discard in the landing contract; label resolution is already available. |
| FT174 | FT229 | The light-path close step dispositions the orphaned tickets the grammar would otherwise be designed around. |
| FT233 | FT199 | Both edit the landing refusal surface; the coordinator settles which refusals survive. |

### Goal track: guidance prose

The guidance-prose backlog follows one path. Its process precursor is landed:
FT164's ticket-contract core shipped, so every later build slices independently
green tickets through the current `craft-tickets` grammar. The payoff facts
shaping the order, verified in-tree 2026-08-02 with an independent mid-tier
refutation pass: `.agents/` and `.bench/BENCH.md` sit outside the gate's reduced
scope, so every separately-landed prose diff pays a full gate — rows batch on the
shared full gate, not just shared files; anchor-pinned files couple prose diffs to
conformance fixture updates. The reviewed Pocock-alignment Spec C has shipped and
FT107 is retired.

FT100 builds last: after FT89 establishes which guidance is authoritative, and
after FT231 supplies the measurement that tells a cut from a regression.

FT99 rides the prose batch. FT106 and FT162 remain independently sequenced by
their existing dependencies. FT133 remains parallel evidence hardening; FT71
stays deferred behind its existing FT169 recommendation. FT172 is outside this
critical path; the FT156 anchor registry shipped, so section-scoped
`.bench/BENCH.md` anchors — the exact surface the prose batch edits — are now
fixture-proven.

## Recommended sequence

1. `/bench-write-spec` — FT226: kit tests stop writing into the operator's real `BENCH_HOME`. Highest severity actionable row — the only one whose defect damages the operator's own machine, and it has no dependencies.
2. `/bench-write-spec` — FT227: adoption smoke, so a newly adopted repository's scaffolded gate can go green. No dependencies; it is the outcome evidence adoption currently lacks.
3. `/bench-write-spec` — FT228: restore `/bench-debug`'s constraints, settle the Codex trigger, check Claude parity. Unblocked, and it folds in the reproduction-economics rule and leaves FT24 parked as the upstream re-check.
