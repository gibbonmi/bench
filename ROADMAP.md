# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-drain` reconcile removes it. Raw capture never lands here — it
goes to `capture/IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>/spec.md`). That path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

FT-to-FT ordering is single-sourced in `## Dependencies`. A literal dependency
blocks the dependent row until its prerequisite lands. A recommended dependency
does not block work. It says the dependent row will be cheaper or less likely to
churn when it is specified after the prerequisite.

## Release-readiness status

Deployment remains **NO-GO** for other projects, public npm release, and bank
use.
The 2026-07-11 release-readiness and repository-controlled compliance
assessments are evidence snapshots; this roadmap is the execution source for
their active findings. Each finding still open appears exactly once on a
`Sources` line below. `RR:` means the release-readiness assessment and `RC:`
means the repository-controlled compliance assessment.

## Features, in priority order

**FT246 (MEDIUM) — other binary-building test packages select the gate executable once.**

**FT247 (MEDIUM) — the next slow test packages extract pure policy from effect adapters.**

**FT249 (MEDIUM, decision required) — a primary-local idea inbox moves to a shared Git ref.**

**FT235 (MEDIUM) — a pool directory's name says what the worktree is for.**

**FT71 (HIGH on the bank track) — versioned local shift evidence.**

**FT162 (MEDIUM, decision required) — the handoff's pins resolve, and its shape has one owner.**

**FT142 (MEDIUM) — FT91 review residuals, two tracks.**

**FT144 (MEDIUM) — kit specs have two audiences and the discipline names neither.**

**FT158 (MEDIUM) — make cross-harness falsification standing for kit-guidance diffs.**

**FT98 (MEDIUM) — one preserve-then-discard primitive; four faces.**

**FT169 (MEDIUM) — the landing command names its refusals and recovery, and its authority is decided.**

**FT253 (MEDIUM) — one landing lease in the intent ledger, from composition through publish.**

**FT207 (MEDIUM, decision required) — worktree-mutating paths share malformed-admin refusal.**

**FT199 (MEDIUM) — a recovery-aware branch-retirement coordinator closes one repository-wide ref inventory.**

**FT178 (MEDIUM) — `bench worktree`'s bare verb is a human porcelain that traps automation and leaks on signals.**

**FT254 (MEDIUM) — `bench worktree exec` is the comfortable path for multi-step work.**

**FT276 (MEDIUM, decision required) — one cycle-free owner derives canonical repository paths.**

**FT277 (MEDIUM) — `bench test --changed` selects the edited packages and explains a widened set.**

**FT255 (MEDIUM, decision required) — concurrent tests share one explicit machine budget.**

**FT274 (MEDIUM) — Bench records Bench-owned FT231 measures with OpenTelemetry.**

**FT275 (MEDIUM) — code built under Bench traces its declared seams with OpenTelemetry.**

**FT172 (MEDIUM, decision required) — the roadmap row grammar is a contract, and `roadmap_id` has one decided source.**

**FT173 (MEDIUM, decision required) — AXI residual: the active-assignment-with-deleted-tree disclosure class.**

**FT202 (MEDIUM, decision required) — a shared purity-census helper, and the census scope for process-backed fixtures.**

**FT89 (MEDIUM) — guidance coherence and current-state documentation.**

**FT106 (MEDIUM) — doc claims re-verified against the tree.**

**FT208 (MEDIUM, decision required) — skills-index producer-hardening residuals: one refusal grammar, per-shape marker diagnostics, and HI14's seam.**

**FT222 (MEDIUM, decision required) — delegate-tier routing has one source in `projects/benchkit.md`.**

**FT204 (LOW, decision required) — one bounded transcript/session query.**

**FT58 (LOW) — hardened pool roots.**

**FT92 (LOW) — attributed subject drift and consumer-shipped input hygiene.**

**FT99 (LOW) — spec problem-premise verification.**

**FT241 (LOW, decision required) — versioned acceptance promises with retained evidence.**

**FT243 (LOW, decision required) — a repository-maintenance skill group.**

**FT244 (LOW) — a standard scratch directory for worktree runs.**

**FT215 (MEDIUM, decision required) — the path-aware lane's open edges: the empty-diff merge, the real-binary hop, and the unknown-path cost.**

**FT261 (MEDIUM, decision required) — preflight review classifies an in-progress untracked spec folder without blocking ticket slicing.**

**FT259 (MEDIUM) — repair coverage changes retain a ticket owner through repair-scoped re-review.**

**FT260 (LOW, decision required) — coordinator worktree diff inspection needs a scoped native path.**

**FT262 (LOW) — preflight reports uncited coverage rows while a ticket breakdown is in progress.**

**FT267 (LOW) — `scripts/verify-release-artifact.mjs` has a gate-owned execution seam.**

**FT268 (LOW) — FT251 residual test hardening: negative callers, bundle confinement, and one checkout-name source.**

**FT269 (LOW) — `craft-gate`: a check on an indirected value grades both ends.**

**FT270 (LOW) — `bench test` names its checks, prints every finding, and runs the system suite.**

**FT271 (LOW, decision required) — census-heads residuals: one record parser, and delimiter-safe head evidence.**

**FT272 (LOW) — `bench status` routes a tickets-only spec folder to a grammar `bench commit` refuses.**

**FT273 (LOW) — a charge names the live-tree check as its probe oracle, and a green probe is verified against the mutated bytes.**

**FT278 (LOW) — `craft-spec` gives every edge-inventory promise a coverage row.**

**FT217 (LOW) — one decision every adopt-lifecycle verb executes.**

**FT279 (LOW, decision required) — `bench link` refuses the kit source checkout.**

**FT280 (LOW, decision required) — a Bench-owned worktree tip projection replaces raw `git rev-parse` in the landing pin.**

**FT218 (LOW) — named git readers instead of learned CLI flags.**

**FT219 (LOW) — `/bench-deepen` refreshes a ready map's frontier to current state before handoff.**

**FT220 (LOW) — `/bench-write-spec` censuses shared decision readers before ticket slicing.**

**FT236 (LOW) — `craft-review` visit: what a string expectation proves, and where an axis under-reads.**

**FT237 (LOW) — `craft-line` states the common case behind ceiling-not-binding.**

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.**

**FT101 (LOW) — per-context scope for monorepos: domain docs and profile.**

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency and dogfood loops.**

**FT108 (LOW) — a refactor lane with a mechanical exit test.**

**FT130 (MEDIUM, decision required) — a capture write mid-lifecycle voids or blocks the run.**

**FT164 (MEDIUM) — repair-lane charges, and a done-claim that resolves its named owners.**

**FT200 (MEDIUM, decision required) — make preflight mechanical at the landing chokepoint.**

**FT258 (MEDIUM, decision required) — `bench commit` derives a complete change set and has one explicit `MERGE_HEAD` contract, without weakening its ownership fence.**

**FT264 (MEDIUM, decision required) — Go build-cache validation residuals.**

**FT265 (MEDIUM, decision required) — one coordinator-owned immutable drain evidence bundle.**

**FT266 (LOW) — `craft-tdd` re-exec helpers return silently outside their role environment.**

**FT166 (LOW) — a sanctioned writer and parser validator for primary-local learnings and retros.**

**FT168 (LOW) — focused iteration evidence: a fixture-selecting canary and mutation probe.**

**FT140 (LOW) — review residuals that want a verdict, not a build.**

## False greens — verdicts that credit unchecked work

Four rows share one failure class: a green whose warrant is missing. The
missing warrant is a stale binary, a dead or skipping citation, a vacuous
baseline, an unchecked absence, or a dependency edge nothing resolves. Each
hardens a different oracle surface. They stay separate builds, but they read
and prioritize as one theme.

**FT133 (MEDIUM) — `bench coverage --check` verifies that red-signal citations resolve.**

**FT174 (MEDIUM) — ticket files have one enforced dependency, ownership, and mutation grammar.**

**FT177 (MEDIUM, decision required) — a stale `dist/bench` invalidates tests and its promotion-broker manifest.**

**FT201 (LOW) — production cancel-signal registrations conform to one source.**

## Reds the diff doesn't own — inheritance, load, and harness defects

Five rows share one failure class: a red that answers for something other
than the diff in front of the gate. The cause is an inherited baseline,
machine contention and a flaky oracle, a literal deadline, or a harness
defect.

**FT141 (MEDIUM, decision required) — red verdicts are recorded against a baseline, so inherited reds stop reading as caused.**

**FT104 (LOW) — load-induced commit refusals: the stop rule and the pre-gate quiet check.**

**FT115 (LOW) — load-robust test and phase deadlines derived from bounds.**

**FT120 (LOW) — gate, canary, and contract test-harness defects nothing asserts.**

**FT257 (LOW) — the maps template and diagnostics state the grammar they enforce.**

## Standards debt — one batched light-path pass

Three rows plus FT142's standards track are shippable together as small
one-source-per-fact and cleanup sweeps under one gate. FT117's parser-routing
half is the largest item in the batch. FT142 itself stays on the main list,
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
   from an isolated prefix without a source checkout. Setup preserves existing
   instructions, settings, and hooks, and it is idempotent and reversible
   (shipped with FT76; re-verified at reassessment).
6. Consumer artifacts exclude maintainer-only capabilities and include the
   supported-platform, security, data-handling, threat, support/EOL, network,
   rollback/recovery, license/notice, SBOM, checksum, and package-inventory
   records.
7. Bank evidence includes redacted local events for success, failure,
   interruption, and recovery. It also includes the preservation,
   oracle-change, post-agent recovery, unknown-hook, environment-minimization,
   offline, and transactional lifecycle contract results.
8. A clean-room user can complete setup from the README in one shell command
   plus at most one harness-native conversation.

Host IAM, OS sandboxing, endpoint controls, firewalls, server-side branch
protection, central CI administration, SIEM/retention, registry administration,
and signing-key custody remain outside this repository-controlled roadmap.

## Parked and scheduled work

**FT6 (LOW, parked pending evidence — leave parked):**

**FT24 (parked pending upstream) — Codex agent-line guard parity.**

**FT8 (scheduled, decision required) — Sonnet 5 mid-tier revisit: a risk-based escalation ladder for the default implementer.**

**FT38 (LOW, decision required) — dashboard visual identity pass.**

**FT231 (EXPERIMENT, decision required) — a measurement harness and the instrumentation it reads.**

**FT232 (EXPERIMENT, decision required) — repair-loop tripwire: an advisory signal from gate records.**

**FT240 (EXPERIMENT, decision required) — iq retrieval experiment: token-budgeted search against the EKS monorepo.**

## Dependencies

The dependent FT is named first. Only the literal table blocks work; the
recommended table is sequencing advice.

### Literal

| FT | Depends on | Why |
|---|---|---|
| FT100 | FT231 | Editorial cuts to the guidance surface without measurement are coin flips; the harness is what tells a cut from a regression. |
| FT240 | FT231 | The adoption rule is a measured three-arm comparison; without the harness the experiment cannot grade its result. |

### Recommended

| FT | Better specified after | Why |
|---|---|---|
| FT71 | FT169 | The event schema should record the settled landing and recovery lifecycle rather than version an interim one. |
| FT100 | FT89 | Cut prose after the correctness and coherence pass establishes which guidance is still authoritative. |
| FT108 | FT89 | FT89 single-sources the skills index the new skill must join; the expand–migrate–contract and gate-cadence rules it builds on are already settled in `craft-tickets`. |
| FT172 | FT106 | Reuse the document-claim probe for semantic roadmap claims instead of designing a second checker. |
| FT162 | FT169 | Build full-run subject resolution on the settled landing primitive. |
| FT166 | FT98 | The porcelain composes over the shipped reduced-gate path allowlist; recoverable set-aside then defines the commit command's smallest sound contract. |
| FT169 | FT98 | Reuse recoverable discard in the landing contract; label resolution is already available. |
| FT253 | FT169 | The lease is one answer to the landing's authority questions; decide those first so the lease keys on the settled lifecycle. |
| FT241 | FT231 | Retained acceptance-run evidence should reuse the harness's record shape rather than version a second one. |
| FT71 | FT274 | The event schema should be expressed in the OpenTelemetry record rather than version a second one. |
| FT232 | FT274 | The advisory reads the red set the record retains. |
| FT204 | FT71 | The query reads the settled event schema before it reads transcript text. |
| FT254 | FT258 | The resolution slice requires the decided `MERGE_HEAD` contract. |

### Goal track: guidance prose

The guidance-prose backlog follows one path. Its process precursor is landed:
FT164's ticket-contract core shipped, so every later build slices independently
green tickets through the current `craft-tickets` grammar. The payoff facts
shaping the order, verified in-tree 2026-08-02 with an independent mid-tier
refutation pass: `.agents/` and `.bench/BENCH.md` sit outside the gate's reduced
scope, so every separately-landed prose diff pays a full gate. Rows batch on the
shared full gate, not just shared files, and anchor-pinned files couple prose
diffs to conformance fixture updates. The reviewed Pocock-alignment Spec C has
shipped and FT107 is retired.

FT100 builds last: after FT89 establishes which guidance is authoritative, and
after FT231 supplies the measurement that tells a cut from a regression.

FT99 rides the prose batch. FT106 and FT162 remain independently sequenced by
their existing dependencies. FT133 remains parallel evidence hardening; FT71
stays deferred behind its existing FT169 recommendation. FT172 is outside this
critical path; the FT156 anchor registry shipped, so section-scoped
`.bench/BENCH.md` anchors — the exact surface the prose batch edits — are now
fixture-proven.

## Recommended sequence

1. Light path — FT133: `bench coverage --check` resolves each cited test name and names the review pickup as a fence member.
2. `/bench-write-spec` — FT274: specify Bench-owned FT231 measures.
3. Reviewer decision — FT271 census-heads residuals, and the FT280 tip-projection surface.
