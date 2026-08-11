# Pocock-aligned guidance doctrine

Status: staged

Decision source: compiled map `specs/pocock-guidance-doctrine/decisions/pocock-alignment.md` (decisions #4, #5, #6, #8, #9, and #10; reviewer-resolved 2026-08-11; restored from `bench spec history remove-spec-build-lifecycle` after the dual retirement)

## Problem

Bench removed the provisional spec-build lifecycle, but the kit still teaches the
old machinery's compensating discipline. Its longest guidance surfaces encode
receipts, parsed ticket evidence, delegate ceremony, and sampled review claims
instead of putting completeness upstream in domain scenarios and re-deriving facts
at review. The result is hard to read cold and hard to keep coherent: `craft-tickets`
is 410 lines, `bench-implement-spec` is 318, `craft-delegate` is 275, and the
always-loaded operating guide is 256. Important operational lessons are mixed into
that sediment, while no gate enforces the reviewed prose budgets.

## Solution

Replace the retired machinery with a compact doctrine stack. A new domain-modeling
skill owns glossary challenges, concept-edge scenarios, code-versus-claim checks,
and inline glossary maintenance. TDD, seam design, and disposable prototyping gain
small on-demand leaves. Grilling works a whole dependency-ready frontier per round;
light-path TDD seams and every ticket breakdown return to explicit reviewer gates.
Review starts by independently deriving primary facts, then compares the candidate
to those facts on the existing Standards, Spec, and Coverage axes.

Slim `craft-tickets`, `bench-implement-spec`, `craft-delegate`, and `craft-line` to
their durable contracts, move branch-only detail behind references, and make the
profile's prose-budget table the one source a fail-closed conformance check consumes.
The same pass folds FT107's still-live operational lessons into their narrow owners.
It does not revive the removed lifecycle or weaken the gate.

## User stories

1. As a reviewer shaping a feature, I want `craft-domain` charged by grilling,
   shaping, and spec authoring so canonical terms, Avoid lists, concept-edge
   scenarios, producer-derived partitions, and code-versus-claim conflicts are
   resolved before implementation, while `CONTEXT.md` remains glossary-only and
   hard-to-reverse decisions continue through `craft-adr`.
   Line: gpt-5.6-sol / high. This is generation-shaping skill and command prose,
   so the profile's leverage override applies.
2. As an implementer working at a chosen seam, I want compact TDD examples,
   dependency-category deepening guidance, radically different interface designs,
   and a disposable prototype route so tests and interfaces are chosen from
   concrete alternatives without turning prototypes into production artifacts.
   Line: gpt-5.6-sol / high. These leaves steer future implementation choices the
   gate can only partially observe.
3. As the reviewer in the loop, I want grilling to ask the whole settled frontier
   per numbered round, light-path TDD to stop for seam confirmation, and ticket
   breakdowns to iterate as a numbered title/blocked-by/delivery list until I
   approve, while the existing batch-approval AFK carve-out remains intact.
   Line: gpt-5.6-sol / high. The story changes who owns planning decisions, so a
   confidently wrong prose edit would multiply across builds.
4. As a reviewer assessing a candidate, I want every axis to derive its facts from
   the current primary source before comparison: Coverage derives the complete
   producer partition and spec-authorized paths, Spec quotes the spec against
   behavior rather than ticket claims, and Standards keeps the Fowler baseline.
   Each finding receives a no-op, auto-fix, or ask-user disposition; review pickup
   state is written and committed before repair; handoffs report raw axis findings
   and de-duplicated repair targets separately.
   Line: gpt-5.6-sol / high. Review doctrine is semantic and weakly gate-observable,
   so the leverage override and up-bias both apply.
5. As a fresh session reading Bench guidance, I want the surviving line,
   delegation, ticket, implementation, and operating rules inside the reviewed
   prose budgets, with one profile-owned table and a canaried conformance check that
   fails closed on an over-budget or unclassified guidance file.
   Line: gpt-5.6-sol / high. The mechanical check is covered, but the pruning choices
   decide what every later session knows.
6. As an agent operating the kit, I want FT107's still-live safety lessons at their
   narrow owners: source-before-call exploration, an observable read-budget reroute,
   tree-consistent handling of non-behavioral spec contradictions, owned-red-set
   convergence before tier escalation, durable review pickup, delegation-authority
   precedence, safe wait and destructive-script conventions, whole-artifact rereads,
   acceptance-shortfall exits, verifier-bootstrap review, hidden-path sweeps, safe
   Bench verb discovery, defaulted-decision checks, and write isolation before debug
   artifacts are created.
   Line: gpt-5.6-sol / high. This is the wide always-loaded prose batch FT107 was
   created to land under one semantic review and gate.

## Implementation decisions

- `craft-domain` is a companion craft skill, not a phase and not a gate. It fires
  from `craft-grill`, `$bench-shape-idea`, and `$bench-write-spec`; every other
  phase consumes the glossary ambiently. It owns domain terms, Avoid lists,
  concrete scenarios at concept seams, producer-derived equivalence partitions,
  code-versus-claim comparison, and inline `CONTEXT.md` updates. `CONTEXT.md`
  remains glossary-only. ADR creation stays with `craft-adr`, whose present
  paragraph format already matches the reviewed upstream shape.
- TDD adds two example leaves, `tests.md` and `mocking.md`, reached only from the
  branches that need them. They teach public-interface behavior, independent
  expected values, system-seam mocking, and realistic failure shapes without
  duplicating the compact rules in `craft-tdd`. Bench keeps its green refactor step;
  upstream's red-green-only loop remains consciously unadopted.
- `craft-seams` extends its existing `references/design-it-twice.md` owner rather
  than creating a second design-it-twice source, and adds one deepening reference
  for the in-process, local-substitutable, remote-owned, and true-external dependency
  classes plus their test strategies. The uncertain-seam trigger still requires
  radically different designs; known seams do not pay the fan-out.
- The user-invoked `prototype` skill backs Prototype decision tickets. A prototype
  answers one named question, is trivial to run, keeps state in memory unless
  persistence is the question, surfaces the relevant state, records the verdict,
  and is then discarded. The refreshed upstream at `mattpocock/skills` commit
  `84fdeff` now retains prototypes on throwaway branches; that drift is deliberately
  not adopted because the reviewed Bench decision says discard.
- `craft-grill` replaces one-question-at-a-time with numbered frontier rounds: every
  question whose prerequisites are settled appears in the same round with a
  recommendation, then the skill waits and recomputes. The refreshed upstream at
  `84fdeff` has itself moved its grilling skill to the same frontier-round shape,
  so the reviewed adoption now matches rather than diverges from its source. Facts remain the agent's job,
  decisions the reviewer's, and enactment still waits for shared-understanding
  confirmation. `bench-shape-idea` uses the same frontier vocabulary.
- Spec sign-off confirms TDD seams for spec-backed work. Light-path work has no spec
  sign-off, so it stops before its first TDD test and presents its seam for reviewer
  confirmation. Ticket breakdown review is no longer a delegate verdict: before
  assignment, the coordinator presents a numbered list of title, `Blocked by:`, and
  delivered outcome, iterates it with the reviewer, and records approval. The
  existing batch-approval AFK carve-out remains the only no-round-trip route.
- Review is re-derive-then-compare. Coverage independently enumerates producer
  membership and the authorized write set from the approved spec; Spec drives the
  behavior and quotes the applicable spec line; Standards independently reads the
  current conventions and retains the Fowler smell baseline. Findings cite the
  derivation. A declaration-only confirmation is incomplete. The three axes run in
  parallel fresh contexts so one derivation cannot seed the next. Review also asks
  what authenticates the verifier before candidate-controlled execution.
- Every finding carries exactly one next-action disposition. `no-op` means the
  candidate or cited source refuted the concern and no repair target remains;
  `auto-fix` means a deterministic hard rule or exact spec predicate can be repaired
  inside already-approved scope; `ask-user` means the finding needs judgment, scope,
  authority, or an oracle change. These are repair-routing labels, not permission for
  the read-only review phase to edit.
- Actionable review findings are written to `reviews/<slug>.md` and committed as an
  ordered step before repair begins, including findings returned by another harness.
  The handoff states both raw per-axis counts and the de-duplicated repair-target
  count. A clean review still writes no pickup artifact.
- `projects/benchkit.md` carries one mechanically parseable prose-budget table:
  `.bench/BENCH.md` at 150 lines, `.agents/commands/bench-implement-spec.md` at 60,
  `.agents/skills/bench-craft-tickets/SKILL.md` at 100, and every other real regular
  file matching `.agents/skills/*/SKILL.md` at 120. That is the complete enumeration
  universe: other command files are outside the reviewed budget, and
  `.claude/skills/*` adapter symlinks are distribution surfaces, not budget subjects.
  The conformance owner parses the table rather than repeating its numbers, applies
  the exact ticket exception before the all-skills default, classifies every newly
  added skill automatically, refuses malformed or duplicate policy and any symlink or
  special file inside the enumerated canonical paths before reading, and reports
  every over-budget path in one run. A canary mutation pushes a classified file over
  its limit and requires the check's own diagnostic. Raising a budget remains a
  reviewer edit to the profile.
- `craft-tickets` shrinks to the independently-green tracer rule, `Blocked by:`,
  `What to build`, Acceptance, frontier order, reviewer-approved breakdown, and the
  advisory `Writes:` note used only to judge parallel disjointness. Contracts,
  Integration surfaces, Closure, covers annotations, red-mutation tables, handoff
  ledgers, and fence enforcement are removed rather than summarized. A meaningful
  contract is stated in What to build and Acceptance and re-derived from the tree by
  review. `bench-implement-spec` becomes the short orchestration pointer to that
  ticket shape, commit-on-green cadence, TDD, review, and final-check.
- `craft-delegate` keeps fresh context, explicit line and bounded charge, worktree
  isolation for writes, read-only review without a worktree, and independent
  verification of every done-claim against the tree and gate. It sheds receipt,
  lifecycle, duplicated mutation, and charge ceremony. When delegation is either
  unavailable or prohibited by the reviewer, the same capability-aware stop and one
  executable handoff applies; a spec-doc-only correction is not a silent exception.
  Any delegation that changes who performs requested work is surfaced before spawn.
- `craft-line` keeps the harness-local tier binding, the three starting signals,
  leverage override, declaration, cap, and ladder. Before retry or escalation it
  classifies only reds owned by the current diff against the pinned inherited
  baseline and spec-predicted reds. A non-shrinking owned-red set stops and surfaces a
  likely seam/spec contradiction; it does not buy a more expensive attempt.
- Every story's `gpt-5.6-sol / high` line is the current Codex-column resolution of
  the profile's top/high leverage override, not the retired spec-build cadence row.
  A build charged from another harness re-resolves the top tier through that
  harness's column at charge time. The obsolete “Spec-build guidance cadence” cache
  entry in `projects/benchkit.md` retires with this batch because its `promote`
  contract no longer exists.
- FT107's operational clauses land once at their narrow owners. `AGENTS.md` owns
  PID-or-sentinel waits, plan-before-apply destructive scripts, `rg --hidden` for
  repository-wide sweeps, and non-interactive Bench discovery through
  `bench commands --brief` or source rather than an unknown bare verb. `craft-seams`
  owns the declared exploration-read budget and `bench outline` reroute.
  `craft-spec` owns a complete-artifact reread after wide structured-prose edits.
  `.bench/BENCH.md` owns the tree-consistent non-behavioral spec reading and material
  acceptance-shortfall exit. `/bench-debug` establishes its isolated writing
  worktree before creating a repro artifact.
- Review treats a compiled map's defaulted decisions as authoritative unless the spec
  explicitly overrides them. A claimed repair is checked against both its coverage
  row and the applicable defaulted-decision table. FT107's former cross-fence atomic
  repair clause has no surviving enforcement target after ownership fences became
  advisory `Writes:` notes, so it is retired rather than generalized.
- Two older FT107 clauses have no surviving target. The plain-Git doc shortcut and
  its squash-merge exception disappeared when `bench commit` became the sole landing
  path. Prospective-promotion selected-binary guidance disappeared with spec-build
  promotion. They are removed from the live guidance target rather than reintroduced.
  The one-build Opus `/bench-debug` override also stays non-standing and out of this
  spec.
- Existing workflow anchors are re-derived against the slim current doctrine. Pins
  whose only subject is deleted ceremony retire with their canaries; surviving
  behavioral obligations keep one owning anchor and a red mutation. Anchor count or
  wording is never copied into prose. The generated skills index, consumer payload,
  Claude symlinks, linked-repo install behavior, and changelog are updated for both
  new skills through their existing owners.
- This is a deliberate bundle despite separable files. The reviewed map chose one
  Spec C, and all stories share the prose-budget source, generated index, anchor
  migration, synthesis audit, and one full gate. Splitting by skill would pay the
  same coupled conformance rewrite repeatedly and could not independently satisfy
  the final budgets.

## Testing decisions

- The mechanical seam is the conformance registry. Focused tests feed a temporary
  profile and guidance tree to the real budget checker, asserting exact diagnostics,
  complete enumeration, newline parity, and fail-closed input classification. The
  existing canary dispatcher proves the over-budget mutation still bites.
- The distribution seam is `bench link` into a throwaway repo. Existing adoption
  tests prove new portable skills arrive and the Claude adapter symlinks resolve to
  the canonical `.agents/skills` files without a copied mirror.
- The semantic seam is a fresh session consuming the changed guidance. A map-to-spec
  dogfood run rebuilds the ignored development binary first when it exercises CLI
  behavior, then verifies frontier grilling, domain-scenario enumeration, seam
  confirmation, and primary-source review from the public phase surfaces. This is
  semantic evidence, not a substitute for the gate.
- The gate seam is the ordinary dev gate: conformance, fixture bite, skills-index,
  adoption, and package-surface checks observe the complete batch. `craft-synthesis`
  additionally requires the legibility and consistency rereads over every changed
  guidance surface; because this batch changes a conformance behavior, the prose-only
  substitution is not enough.

### Seam diagram

    trigger: fresh agent enters shaping, spec, build, or review
        │
        ▼
    prompt + tree ──▶ [ compact doctrine skill / phase ] ──▶ decision, ticket, or finding
                            ◀ semantic tests attach here: fresh-session dogfood
                              re-derives from the named primary source
        │
        ▼
    profile budget table ──▶ [ conformance registry ] ──▶ per-path green/red diagnostics
                                  ◀ mechanical tests attach here: temporary guidance
                                    tree plus over-budget canary mutation

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| PG1 | 1 | `craft-domain` exists, is indexed, and is charged by grilling, shaping, and spec authoring while ordinary phases read `CONTEXT.md` ambiently | skills index plus fresh-session shaping/spec dogfood | observed red 2026-08-11: `.agents/skills/bench-craft-domain/SKILL.md` is absent and the generated index has no domain-modeling entry | a missing skill or trigger leaves completeness enumeration in the retired downstream machinery |
| PG2 | 1 | domain work challenges canonical terms and Avoid lists, walks concrete concept-edge scenarios including producer-derived partitions, checks code against claims, and updates glossary-only `CONTEXT.md` inline | fresh-session domain scenario dogfood | not TDD-able: no deterministic parser can judge whether a scenario exposed a concept ambiguity; the fresh dogfood transcript and falsification review must cite the scenario and resolved predicate | a vocabulary-only shell could exist yet never enumerate the domain cases moved upstream by decision #4 |
| PG3 | 1 | hard-to-reverse architectural decisions route to `craft-adr`; `craft-domain` does not become a second ADR format owner | skill cross-reference and semantic review | already covered: `craft-adr` owns the one-paragraph current-state format and the build must preserve that owner while adding the pointer | catches duplicated ADR policy that would drift between two skills |
| PG4 | 2 | `craft-tdd` reaches separate good/bad-test and mocking leaves that teach independent expectations, public-interface behavior, system-seam mocks, and realistic failures | skill references plus complete-leaf reread | observed red 2026-08-11: no `tests.md` or `mocking.md` leaf exists under `bench-craft-tdd` | a one-line pointer without the leaves, or leaves never referenced, cannot teach the adopted examples |
| PG5 | 2 | `craft-seams` classifies in-process, local-substitutable, remote-owned, and true-external dependencies to their test strategies and keeps one design-it-twice owner for radically different interfaces | seam reference leaves and high-effort design charge | observed red 2026-08-11: the design-it-twice reference exists but no deepening/dependency-category leaf exists | enumerating all four classes catches a partial port that handles only mocks and in-process code |
| PG6 | 2 | `prototype` is indexed and drives one named question through a trivial-to-run, state-visible, non-persistent artifact, records the verdict, then discards the artifact | skills index plus prototype ticket dogfood | observed red 2026-08-11: `.agents/skills/prototype/SKILL.md` is absent | catches an ordinary implementation skill mislabeled as a prototype or upstream-style branch retention |
| PG7 | 3 | grilling asks every dependency-ready question in one numbered frontier round, waits for the answers, recomputes, and still requires final shared-understanding confirmation | `craft-grill` fresh-session transcript | observed red 2026-08-11: `craft-grill` mandates one question per turn and contains no whole-frontier rule | catches both the old serial cadence and a one-shot questionnaire that never recomputes blockers |
| PG8 | 3 | light-path TDD stops before its first test for reviewer seam confirmation; spec-backed TDD consumes the signed-off spec seam without a second gate | light-path TDD phase transcript | observed red 2026-08-11: `craft-tdd` accepts only seams named by `$bench-write-spec`, leaving the sanctioned no-spec light path without a confirmation route | prevents TDD from inventing an unspecced seam while avoiding a duplicate human gate for spec-backed work |
| PG9 | 3 | before assignment, the reviewer approves an iterated numbered ticket list containing title, `Blocked by:`, and delivered outcome; batch approval is the only AFK carve-out | implement-spec planning handoff | observed red 2026-08-11: `craft-tickets` and `bench-implement-spec` route breakdown approval to a fresh read-only delegate | catches retaining delegate sign-off or omitting one of the three facts the reviewer needs to reslice |
| PG10 | 4 | Standards, Spec, and Coverage remain separate, and each finding cites facts independently re-derived from its primary source before comparison | three-axis review dogfood against a planted declaration-only candidate | observed red 2026-08-11: `craft-review` validates declared artifacts but has no primary-source re-derive-then-compare mandate | a review that trusts ticket claims passes the planted producer-membership omission |
| PG11 | 4 | Coverage enumerates the producer-derived input family and spec-authorized write set, Spec quotes the spec against behavior, Standards reads current rules and retains the Fowler baseline, bootstrap review asks what authenticates the verifier, and the three axes run in parallel fresh contexts | per-axis charged review transcripts | not TDD-able as one pre-build red: these are independent semantic derivations; the dogfood candidate plants one mismatch per axis and requires three independently cited findings | per-axis planted mismatches prevent a generic “review primary sources” sentence or one shared-context derivation from satisfying the story |
| PG12 | 4 | every finding receives exactly one `no-op`, `auto-fix`, or `ask-user` disposition; actionable findings are written and committed before repair, including cross-harness returns; the handoff reports raw axis findings separately from unique repair targets | review command handoff and Git history | observed red 2026-08-11: the current review command owns pickup state but carries no three-value disposition vocabulary, ordered cross-harness capture step, or dual-count handoff | catches lost findings, unowned judgment calls, and confusion between review volume and repair work |
| PG13 | 5 | the profile is the one source for the complete budget universe: `.bench/BENCH.md` 150, `.agents/commands/bench-implement-spec.md` 60, `bench-craft-tickets/SKILL.md` 100, and every other real `.agents/skills/*/SKILL.md` 120; every subject lands within its limit | profile parser plus complete canonical-path enumeration | observed red 2026-08-11: line counts are 256, 318, 410, and multiple other canonical skills exceed 120 | an incomplete prune, accidental adapter-symlink scan, or hard-coded second budget table stays red |
| PG14 | 5 | malformed, missing, or duplicate budget policy; a newly added canonical skill; a symlink or special file inside the canonical universe; and every over-budget subject all fail closed with path-specific diagnostics in one run | conformance checker over temporary trees | not TDD-able until the parser seam exists: the implementing ticket writes these cases first and proves each exact diagnostic before pruning prose | catches green-by-omission, wildcard drift, blocking reads, and first-error-only reporting without treating `.claude/skills` adapters as canonical inputs |
| PG15 | 5 | a canary mutation pushes a classified guidance file over budget and receives the budget check's own diagnostic; the complete dev gate is green after pruning | conformance canary dispatcher and dev gate | not TDD-able until the new checker and diagnostic exist: the implementing ticket registers the checker and canary first, observes the targeted red on today's over-budget subject, then performs the prune | proves the new oracle bites and is not an uncalled helper or uncoupled grep |
| PG16 | 5 | the slim ticket contract contains only title, `Blocked by:`, What to build, Acceptance, and advisory `Writes:`; frontier order and serial commit-on-green remain, while retired parsed evidence fields and delegate breakdown verdict disappear | complete skill/command reread plus existing CLI/gate positive controls | observed red 2026-08-11: `craft-tickets` is 410 lines and still specifies Contracts, Integration surfaces, Closure, covers annotations, red mutations, and a delegate review | catches a cosmetic trim that leaves the lifecycle's schema and authority intact |
| PG17 | 5 | slim delegation keeps fresh explicit charges, isolated writes, read-only review, and independent tree/gate verification; unavailable or reviewer-prohibited delegation stops with one executable handoff and delegation-authority changes are surfaced before spawn | delegate skill dogfood in allowed, incapable, and prohibited postures | observed red 2026-08-11: the skill handles inability but not reviewer prohibition and is 275 lines of lifecycle and mutation ceremony | exercises both fail-closed postures and prevents slimming away isolation or verification |
| PG18 | 5 | line governance keeps harness-local bindings, start signals, leverage, declaration, cap, and ladder, but classifies owned reds against the pinned baseline before retry or escalation and stops on a non-shrinking owned-red set | line-routing conformance plus planted gate-output transcript | observed red 2026-08-11: `craft-line` escalates on a second red without first classifying inherited and spec-predicted reds or testing convergence | catches an expensive escalation loop against a contradictory spec |
| PG19 | 6 | no API or function is called before its definition is read this session; a declared small exploration-read budget without traction reroutes through `bench outline` and reports the budget spent | invariant and seam-skill dogfood transcript | observed red 2026-08-11: BENCH says only “read the surrounding code,” and `craft-seams` has no observable read-budget reroute | makes source verification and the stop condition observable instead of self-graded |
| PG20 | 6 | a non-behavioral spec contradiction follows the current tree convention and is flagged for veto, while behaviorally different readings stop; only diff-owned reds count toward fix-loop convergence | spec/build stop-path transcript | observed red 2026-08-11: no current guidance states either predicate | catches silent behavioral choice and false non-convergence caused by inherited or predicted reds |
| PG21 | 6 | review capture, whole-artifact reread, material acceptance-shortfall exit, raw-versus-unique finding counts, verifier-bootstrap questioning, and defaulted-decision authority each remain explicit after the prose diet | command/skill semantic reread and review dogfood | observed red 2026-08-11: complete-artifact reconciliation, dual review counts, verifier-bootstrap review, and defaulted-decision checking are absent; pickup persistence is present but subordinate | catches pruning that meets the line number by deleting the lessons the batch exists to preserve |
| PG22 | 6 | project shell guidance uses PID or sentinel waits, plan-before-apply destructive scripts with exact sampled targets, hidden-path repository sweeps excluding `.git`, and safe Bench discovery plus non-interactive stdin | AGENTS shell-convention review with one planted counterexample per rule | observed red 2026-08-11: none of the four rules appears in `AGENTS.md` | catches the exact self-matching wait, missed dot-path, accidental mutation, and interactive bare-verb failures behind FT107 |
| PG23 | 6 | debug work that may write creates or selects its isolated worktree before the first repro artifact, and the one-build Opus delegate-debug override is not generalized | debug entry transcript and residue sweep | observed red 2026-08-11: `/bench-debug` defers delegation/isolation until the fix phase and gives repro artifacts no entry isolation rule | prevents a clean-at-start main checkout from becoming unattributably dirty mid-debug |
| PG24 | 1-6 | a fresh-session map-to-spec dogfood run and final whole-artifact reread find no contradiction among triggers, budgets, anchors, and phase handoffs; generated index, linked-repo payload, and Claude symlinks expose both new skills from their canonical files | fresh-session synthesis dogfood plus `bench link` system seam | not TDD-able as pure pre-build semantics; distribution controls are already covered, while the changed trigger behavior requires fresh context after the edit | catches a locally coherent skill that is undiscoverable, unshipped, mirrored, or contradicted by another phase |

### Edge inventory

- Error path — PG14 and PG15: malformed policy and an over-budget subject emit
  distinct path-specific diagnostics and make the gate red.
- Empty/absent input — PG14: missing policy, no skills, and a newly added unclassified
  skill fail closed; the complete enumerated skill set cannot vacuously pass.
- Boundary values — PG13 and PG15: exactly-at-budget is green; one line over is red;
  files with and without a trailing newline count logical lines consistently.
- Malformed input — PG14: duplicate rows, invalid numeric limits, unknown classes,
  and overlapping exact/default classifications refuse rather than pick a winner.
- Interrupted/partial state — PG12 and PG23: review pickup is committed before repair,
  and debug artifacts begin isolated, so a session death leaves attributable state.
- Re-run idempotency — PG7 and PG24: the next grill round recomputes from settled
  answers; skills-index generation and link convergence remain byte-identical on a
  second run through their existing controls.
- Process-boundary lifecycle — PG12 and PG24: another session can recover committed
  findings and a fresh harness discovers the same canonical skill files.
- Hostile environment — PG14, PG22, and PG24: dot-directories are enumerated,
  symlinks and special files are classified before read, missing global `bench` uses
  linked resolution, and non-interactive stdin cannot hang discovery.
- Composition degenerate — PG10, PG13, and PG24: each individual skill can read well
  while its trigger, budget classification, index entry, adapter symlink, or consuming
  phase is missing; the fresh-session and full-enumeration seams exercise the real
  composition.
- **Won't handle:** prototype persistence or UI-specific prototype polish — the
  reviewed decision requires a disposable, question-shaped artifact and this repo
  has no design source.
- **Won't handle:** tracker-backed maps and the upstream two-axis review — both are
  closed non-adoptions; Bench keeps local decision maps and all three review axes.
- **Won't handle:** byte-identical and divergent foreign adapter targets — the
  landed `bench-link` repair keeps both as hard refusals; this guidance batch does not
  touch adoption behavior.

## Ownership fences

- `AGENTS.md`, `CONTEXT.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`
- `.agents/commands/bench-shape-idea.md`, `.agents/commands/bench-write-spec.md`,
  `.agents/commands/bench-implement-spec.md`,
  `.agents/commands/bench-review-implementation.md`,
  `.agents/commands/bench-debug.md`
- `.agents/skills/bench-craft-domain`, `.agents/skills/bench-craft-grill`,
  `.agents/skills/bench-craft-tdd`, `.agents/skills/bench-craft-seams`,
  `.agents/skills/prototype`, `.agents/skills/bench-craft-review`,
  `.agents/skills/bench-craft-tickets`, `.agents/skills/bench-craft-delegate`,
  `.agents/skills/bench-craft-line`, `.agents/skills/bench-craft-spec`
- `.claude/skills/bench-craft-domain`, `.claude/skills/prototype`
- `internal/anchors`, `internal/conformance`,
  `tests/canary/workflow-guidance-anchors`, `tests/canary/guidance-prose-budgets`
- `projects/benchkit.md`, `README.md`, `CHANGELOG.md`, `ROADMAP.md`,
  `capture/session-handoff.md`

## Out of scope

- Rescoping `axi-coherent-diff` and `axi-query-disclosure`: about 6 edits and 1 gate
  run each. These are separate behavior-first capabilities sequenced after Spec C.
- Implementing `single-build-serial-gate`: about 12 edits and 3 gate runs. Its
  runtime scheduler is independent of this guidance batch.
- Consolidating FT158, FT165, FT100, or the other remaining guidance roadmap rows:
  about 20 edits and 2 gate runs for a future restructure spec. FT107 preserves
  their visible inputs but does not silently absorb them.
- Reintroducing removed spec-build commands, receipts, recovery refs, or promotion:
  a separate reversal of the reviewed lifecycle-removal decision, not unfinished
  doctrine work.
- Changing gate authority, dropping the Coverage axis, or widening unowned
  `bench-link` convergence beyond canonical same-file adapter targets: each is a
  separate reviewer decision and remains closed.
