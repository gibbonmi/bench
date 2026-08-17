# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
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

## Features, in priority order

**FT71 (HIGH on the bank track) — versioned local shift evidence.**

**FT198 (MEDIUM, decision required) — a progressively loaded roadmap.**

**FT162 (MEDIUM) — full-run and phase-close state has one authoritative subject and handoff.**

**FT142 (MEDIUM) — FT91 review residuals: eight open findings, two tracks.**

**FT144 (MEDIUM) — kit specs have two audiences and the discipline names neither.**

**FT158 (MEDIUM) — make cross-harness falsification standing for kit-guidance diffs.**

**FT98 (MEDIUM) — one preserve-then-discard primitive; four faces.**

**FT169 (MEDIUM) — one sanctioned worktree landing command owns the stale-base dance.**

**FT207 (MEDIUM, decision required) — worktree-mutating paths share malformed-admin refusal.**

**FT199 (MEDIUM) — a recovery-aware branch-retirement coordinator closes one repository-wide ref inventory.**

**FT178 (MEDIUM) — `bench worktree`'s bare verb is a human porcelain that traps automation and leaks on signals.**

**FT212 (LOW) — `bench worktree clean --landed` fails "invalid invocation" though `.bench/BENCH.md`'s inventory advertises the form.**

**FT172 (MEDIUM) — the roadmap parser and context snapshot make the drain's non-recurrence evidence complete.**

**FT173 (MEDIUM, decision required) — AXI residual: the active-assignment-with-deleted-tree disclosure class.**

**FT202 (MEDIUM, decision required) — a standing test-support fence, and the census scope for process-backed fixtures.**

**FT185 (MEDIUM) — gate results join the structured Bench output contract.**

**FT89 (MEDIUM) — guidance coherence and current-state documentation.**

**FT106 (MEDIUM) — doc claims re-verified against the tree.**

**FT190 (MEDIUM) — every injected interface has a real-producer test or a written exemption.**

**FT191 (MEDIUM) — a fixture-and-seam inventory a charge can carry for free.**

**FT192 (MEDIUM) — one source per fact reaches spec and ticket prose.**

**FT206 (MEDIUM) — exact-candidate review sees destination metadata before it freezes.**

**FT208 (MEDIUM, decision required) — skills-index producer-hardening residuals: one refusal grammar, per-shape marker diagnostics, and HI14's seam.**

**FT209 (MEDIUM) — a behavior-preserving refactor proves itself by differential, and a new grouping fixes its cardinality.**

**FT204 (LOW, decision required) — one bounded transcript/session query.**

**FT205 (LOW) — `craft-delegate` names the delegate-worktree end-of-life pair.**

**FT213 (MEDIUM) — a read-only delegate reading a graded tree gets its own worktree, and a delegate's claim about a gate signal gets an oracle-verified probe.**

**FT58 (LOW) — hardened pool roots.**

**FT92 (LOW) — attributed subject drift and consumer-shipped input hygiene.**

**FT99 (LOW) — spec problem-premise verification.**

**FT214 (MEDIUM) — a build may not edit its own spec's acceptance rows, budget targets, or ownership fences.**

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.**

**FT101 (LOW) — per-context scope for monorepos: domain docs and profile.**

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency and dogfood loops.**

**FT108 (LOW) — a refactor lane with a mechanical exit test.**

**FT111 (LOW) — provenance tags that outlive their specs.**

**FT112 (LOW) — an approximation that stays green is not a cleared bug.**

**FT113 (LOW) — `bench commit --spec` residuals: the flip counts as a path, and the flip has one author.**

**FT130 (MEDIUM) — a capture write mid-lifecycle voids or blocks the run.**

**FT138 (LOW) — instrument Bench so build economics are measurable.**

**FT164 (MEDIUM) — repair-lane charges, and a done-claim that resolves its named owners.**

**FT200 (MEDIUM, decision required) — make preflight mechanical at the landing chokepoint.**

**FT165 (LOW) — fold the domain-modeling discipline into `/bench-shape-idea`.**

**FT180 (LOW) — a spec-optional route decided at shape-idea's exit.**

**FT182 (LOW) — a Planned-phase receipt over an absent target wedges the abandon retry.**

**FT197 (MEDIUM) — the Go core owns gate invocation and process lifetime.**

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
diff in front of the gate — an inherited baseline, machine contention, a
literal deadline, a harness defect, a flaky oracle.

**FT141 (MEDIUM) — `bench gate pin` records red verdicts, so inherited reds stop reading as caused.**

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

## Parked and scheduled work

**FT6 (LOW, parked pending evidence — leave parked):**

**FT24 (parked pending upstream) — Codex agent-line guard parity.**

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.**

**FT38 (LOW, decision required) — dashboard visual identity pass.**

**FT170 (LOW, decision required) — behavioral red/green evaluation for skill guidance.**

## Dependencies

The dependent FT is named first. Only the literal table blocks work; the
recommended table is sequencing advice.

### Literal

| FT | Depends on | Why |
|---|---|---|
### Recommended

| FT | Better specified after | Why |
|---|---|---|
| FT71 | FT169 | The event schema should record the settled landing and recovery lifecycle rather than version an interim one. |
| FT100 | FT89 | Cut prose after the correctness and coherence pass establishes which guidance is still authoritative. |
| FT108 | FT89 | FT89 single-sources the skills index the new skill must join; the expand–migrate–contract and gate-cadence rules it builds on are already settled in `craft-tickets`. |
| FT111 | FT179 | FT179 reopens FT111's edit-in-place-only ruling on an order-larger measured count; land them as one `craft-comments`/`craft-review` visit. |
| FT172 | FT106 | Reuse the document-claim probe for semantic roadmap claims instead of designing a second checker. |
| FT162 | FT169 | Build full-run subject resolution on the settled landing primitive. |
| FT166 | FT98, FT113 | The porcelain composes over the shipped reduced-gate path allowlist; recoverable set-aside then defines the commit command's smallest sound contract. |
| FT168 | FT153 | Expose focused canary execution after baseline meaning is settled. |
| FT169 | FT98 | Reuse recoverable discard in the landing contract; label resolution is already available. |

### Goal track: guidance prose

The guidance-prose backlog follows one path. Its process precursor is landed:
FT164's ticket-contract core shipped 2026-08-03,
so every later build slices independently green tickets through the current
`craft-tickets` grammar. The payoff
facts shaping the order, verified in-tree 2026-08-02 with an independent
mid-tier refutation pass: `.agents/` and `.bench/BENCH.md` sit outside the
gate's reduced scope, so every separately-landed prose diff pays a full gate
— rows batch on the shared full gate, not just shared files; anchor-pinned
files couple prose diffs to conformance fixture updates (`craft-delegate` 14
anchors, `bench-implement-spec.md` 35+, `.bench/BENCH.md` 17). The reviewed
Pocock-alignment Spec C has shipped and FT107 is retired.

1. Shape FT198; the
   board's 170 KB full snapshot now confirms the progressive-loading trigger.
2. FT100's remaining work grills and builds last, after FT89 establishes which
   guidance is authoritative and the reviewer resolves FT170's benchmark route.

FT99 rides prose batch 1. FT106 and FT162 remain independently sequenced by
their existing dependencies. FT133 remains parallel evidence hardening; FT71
stays deferred behind its existing FT169 recommendation. FT172 is outside this
critical path; the FT156 anchor registry shipped, so section-scoped
`.bench/BENCH.md` anchors — the exact surface the prose batch edits — are now
fixture-proven.

## Recommended sequence

1. `/bench-shape-idea` — FT198 decides the progressive roadmap shape; `ASSESSMENT.md` ranks it 0 and this drain paid the 170 KB snapshot to reach three rows.
2. `/bench-shape-idea` — FT207 decides whether worktree-mutating paths share FT189's malformed-admin refusal before Git can block.
3. `/bench-write-spec` — FT213 gives a read-only delegate its own worktree when it reads a tree also being graded, and requires an oracle-verified probe before trusting a delegate's gate-skip claim; reproduced twice in this drain's source retro.
