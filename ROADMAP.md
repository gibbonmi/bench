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

**FT71 (HIGH on the bank track) — versioned local shift
evidence.** Emit a redacted, append-only local event schema for shift/session
Occurrences: baseline-01
start and end, resolved agent and line, gate fingerprint and verdict, adapter
result, commit or recovery reference, cleanup decision, and release-evidence
relationship. Atomic append, rotation, and explicit retention are part of the
repository contract. Records distinguish completed, failed, interrupted,
recovered, and abandoned work; normal exit completes the matching intent and
crash recovery requires the active lease identity.

Local records are documented as mutable evidence inputs, not a tamper-proof
central audit system. Closure includes success, failure, interruption,
recovery, normal subshell completion, stale intent, and redaction fixtures. The
repository-controlled bank evidence requirement makes this row active.

Sources: `RR:C-05`; `RC:H-03`.

**FT162 (MEDIUM) — full-run and phase-close state has one authoritative subject
and handoff.** One lifecycle owner resolves the exact final-check tree before a
gate, exposes an active dirty assignment, and records open assignments and branch
commit state in the handoff. The implementation-retro pass also records
codification candidates with evidence, durable owner, and expected effect.

Occurrence: craft-tickets/light-path/artifact-suite/artifact-hoist retros — final-check gated the wrong checkout while the final ticket remained uncommitted (sources: those retros, drained here).

Concurrent main-tree writers need visible intent or a lease; two sessions sharing
a checkout serialize lifecycle mutations and landings, side-work uses a worktree,
and foreign dirt refuses or uses FT98's recoverable set-aside.

Recovered review artifacts retain their frozen candidate pair and terminal verdict;
resume reruns every review axis before repair or landing. Assigned worktree authoring
keeps `main` read-only except controlled commits and landings.

Occurrence: 2026-08-15 interrupted FT189 review — a committed review pickup had no terminal verdict or recovery handoff.
Occurrence: 2026-08-15 capture idea — concurrent authoring on `main` motivated assigned-worktree and serialized-landing evaluation.

The handoff decision chooses one authoritative `capture/session-handoff.md`,
per-spec handoffs, or a generated per-workstream projection; it must not create a
second convention. Source: `capture/IDEAS.md`, drained here.

Occurrence: 2026-08-02 `5d67654` landing — an unverified capture claim says a concurrent primary-checkout landing discarded uncommitted work (`capture/learnings.md`, verdicted here).

The full-run owner records staging base, terminal implementation SHA, explicit
`bench diff --full` base/head, exact subject, and subject-bound performance; a
different subject's number is omitted. Entry: `/bench-shape-idea`.

Occurrence: artifact-suite/decision-map integrity runs — the review surface had to reconstruct the base/head range from a single landing or handoff path.

The close decision resolves the conflict between `bench spec retire` removing a
row and `/bench-final-check` leaving it for `/bench-what-next`; retirement gets a
planned manifest and explicit commit path. Until retained-run fallback exists,
final-check and retro capture precede deletion at `internal/spec/spec.go:400`.

Occurrence: check-level-conformance-scoping retro — retirement removed a live spec before terminal status and retro capture no longer had its path.

The handoff `## Next command` is exactly one backticked harness-native invocation.
The terminal record survives retirement with promotion timings, and final-check
captures the retro rather than rerunning a successful promotion to manufacture it.
Sources: the check-level-conformance-scoping retro and `capture/IDEAS.md`, drained here.

Occurrence: FT194 project-green-desync retro — prepared state kept `next` at `resume promote` after review acceptance and repair assignment.

The terminal projection derives `next` from current lifecycle state, retains
promotion identity/evidence, and exposes phase heartbeats, stage timings,
recomposition checkpoint ownership, bounded compatible paths, and journal capacity.
Sources: the `ft194-project-green-desync` and `exact-prospective-landing` retros, drained here.

The same owner exposes blocking `bench gate --fresh`, retained first-failure data,
active-first status, `promote --check`, exact-tree peak-memory attribution, and
candidate-bound failure/per-phase receipt projections. Sources: covers-traceability,
FT195, `ft187-communication-surface-cut`, and canary-planted-reason-ownership retros.

The handoff offers a durable-boundary fresh-context escape that pins candidate,
assignment, evidence, dirty paths, closed decisions, constraints, and next command;
it never switches during a command or unrecorded mutation. Source: `capture/IDEAS.md` 2026-08-06.

The gate-decision rule routes pending, stale, or mismatched records through plain
`bench gate`; reviewer-approved `--fresh` is only for reusable green evidence that
cannot compose. Source: `capture/learnings.md` 2026-08-07.

Occurrence: gate-test-concurrency retro — terminal status omitted stage timings despite retaining the commit and green digest.

Branch-native state keeps a schema-migratable terminal record outside replaceable
gate code and gives final-check an exceptional route requiring prospective log,
published commit, and stated limitation. Sources: branch-native-build-test-architecture
and repair-ticket-reslicing retros, drained here.

The coherent-diff repair binds a second semantic-review preflight to the repaired
product candidate and chooses an explicit review range or repair-candidate selector;
the ownership fence is not weakened. Source: the axi-coherent-diff retro, drained here.

Final landing guidance states the source-built binary, `BENCH_HOME`, clean
destination, and declared allowance before the landing invocation. Source:
parallel-session-landings retro, drained here.

**FT142 (MEDIUM) — FT91 review residuals: eight open findings, two tracks.**
The retired FT91 review remains the canonical carrier; full citations are
recoverable via `bench spec history ft91-gate-tier-split`. Standards work keeps
diff/provenance comments, the duplicated release-only package list, and the
hard-coded `CHECK` filename. Ship work keeps bounded conformance, early
`govulncheck` refusal, interrupted `dist/.preflight.*` cleanup, and a decision
for concurrent timing-file ownership. The release-only `go test` step is closed
by `ft91-gate-phase-split` story 18 and stays an explicit `prep-release` step.

Occurrence: ft91-gate-tier-split review — twelve findings were reduced to eight open residuals after three pre-merge closures (source: FT91 review, promoted at retirement).

The second review's remaining standards and coverage choices stay in this row:
single-source the contract import prefix and comments, then decide `BENCH_SKIP_LOG`,
compile environment, timeout, and `TMPDIR` posture for bite invocations. Full text
remains recoverable at `git show 4429b05:reviews/ft91-canary-compiled-bites.md`.

Occurrence: 2026-07-28 `ft91-canary-compiled-bites` review — seventeen non-blocking findings left import-prefix, comment, environment, timeout, and `TMPDIR` choices (source: `capture/IDEAS.md`, drained here).

**FT144 (MEDIUM) — kit specs have two audiences and the discipline names
neither.** `craft-spec` needs an edge-inventory prompt that distinguishes the
Bench tree from linked repositories, because the same fail posture can be right
for one audience and a regression for the other. The workflow decision is whether
an intent-preserving seam move may ship under batch approval with a veto flag or
must return to `/bench-write-spec`; changes to what is built still stop for sign-off.

Occurrence: 2026-07-27 `ft91-gate-phase-split` — audience partitioning caught a linked-repo-name probe and manifest drift dropped story 9 under batch approval (sources: `capture/learnings.md`, drained here; `capture/session-handoff.md`, prior drain).

Occurrence: 2026-07-27 `ft91-artifact-build-tiering` — review corrected four tree-falsified approved-spec statements under veto flags (source: `capture/learnings.md`, drained here).

The edge inventory also distinguishes an absent tickets directory from a present
but empty one; only the latter declares an unowned ticket surface. Source: the
spec-authoring-and-light-path retro, drained here.

Cap-accepted verification records the cap and declined findings, then probes each
folded red signal through its named public fixture. Transaction-shaped specs classify
owned failures as pre-oracle persistence, in-oracle interruption, or terminal persistence.

Occurrence: 2026-08-14 FT189 falsification — adversarial review, not authoring prose, found deep consumer omissions across 33 rounds.
Occurrence: 2026-08-15 gate-run-transaction cap — folded mutations needed public-fixture reachability probes after capped review.
Occurrence: 2026-08-15 gate-run-transaction review — lifecycle-position coverage missed owner/pending writes and in-oracle cancellation.

**FT158 (MEDIUM) — make cross-harness
falsification standing for kit-guidance diffs.** The prepared spec/diff bundle
gets an advisory Codex pass charged to refute kit guidance before final review.
Every finding receives accept, merge, or dismiss disposition; the pass remains
advisory and stays at the public CLI seam.
Occurrences: baseline-01, baseline-02, baseline-03
Kit edit under `craft-synthesis`.

Occurrence: FT91/FT152/FT123+FT124 falsification passes — each found a real guidance defect after ordinary review had cleared it (sources: `capture/learnings.md`, gate-fastpath and FT123+FT124 retros).

Occurrence: check-level-conformance-scoping retro — cross-harness review caught status, dashboard, and handoff disagreement that package tests missed.

The pass runs before final review and turns each finding into an independently
red repair ticket. Source: the `ft156-anchor-registry` retro, drained here.

Occurrence: conformance-harness-scope retro — a reviewer lacked mutation receipts and rediscovered a valid mutation question.

The coordinator supplies mutation receipts and distinguishes permanent regression
tests from temporary red overlays. A clean three-axis review is terminal; one
deterministic promotion red may create one evidence-backed repair but never
restarts open-ended review. Sources: conformance-harness-scope and
spec-ticket-handoff-contract retros, drained here.

Premise verification checks a finding's subject and ownership before accepting its
label. Source: the Pocock-guidance-doctrine retro, drained here.

**FT98 (MEDIUM) — one preserve-then-discard primitive;
four faces.** Three rows share one recoverable discard contract.
Occurrences: 2026-07-30-scoped-roadmap-commit, baseline-01, baseline-02
Preserve a recovery ref first, require an exact fingerprint, refuse moved
content, plan one target per invocation, and distinguish proof from operator
judgment.

Occurrence: 2026-07-20/22 recovery payloads — reshaped or squash-landed content defeated commit-based landed proof and required manual ref cleanup (sources: `capture/IDEAS.md`, `capture/learnings.md`, drained here).

The landed face covers recovery payloads: cleanup must not infer landedness from
commit identity when content was reshaped.
The landed face uses `LandedInDefault` reverse-diff proof (`efb456c`) or
reviewer-supplied `bench worktree clean --discard-branch` proof (`37411a0`),
and remains fail-closed. The recovery/reclaim verbs and preservation refs are
removed; resume reconciles their namespaces away. The stale-fingerprint retire
check still needs runtime-contract coverage, and orphan deletion must carry the
planned OID through compare-and-swap. Sources: `capture/IDEAS.md` and
`capture/learnings.md` 2026-08-03, drained and verdicted here.

Occurrence: 2026-08-01 landed-proof build — squash proof and reviewer-supplied discard proof shipped while ambiguity remained not-landed (`efb456c`, `37411a0`).

Occurrence: 2026-08-04 recovery-discard build — inspected payload discard and provisional-residue reclaim shipped before lifecycle removal (`fafb049`).
Occurrence: 2026-07-30 scoped roadmap commit — foreign session-handoff dirt forced FT169's isolated verification-worktree landing workaround because no sanctioned set-aside existed.

Face two gives path-scoped `bench commit` a sanctioned set-aside and gives
`bench commit --all` an explicit full-working-set attribution while preserving
the default fence. Face three gives mutation probes a recoverable scoped revert
instead of `git checkout` or a guard exemption; both reuse the contract above.

Occurrence: 2026-07-30 scoped roadmap commit — foreign session-handoff dirt blocked the landing because no sanctioned set-aside existed (source: FT169, `capture/learnings.md`).

Occurrence: 2026-08-01 migration landing — a deliberate whole-tree change required three refusals before an explicit working-set contract existed.

The ignored-cache face accepts only lifecycle-owned or summarized generated
directories under the same size-bounded fingerprint. It also needs a stable
`--fingerprint-only` plan view and a named derived-cache posture. Sources: FT131
implementation and spec-build-lifecycle-preconditions retros, drained here.

Occurrence: 2026-08-03 FT183 landing — foreign capture files forced copy-out and restore after revert and path-scoped stash were blocked (`capture/learnings.md`, verdicted here).

Generated residue remains bounded and fingerprinted; explicit clean refuses dirty
removal unless the reviewer-visible discard contract owns the payload through
completion. Sources: repair-ticket-reslicing and remove-spec-build-lifecycle
retros, drained here.

**FT169 (MEDIUM) — one sanctioned worktree landing command
owns the stale-base dance.** The command fast-forwards the assignment onto the
current default, creates the gated commit, fast-forwards the result back, and
releases the worktree atomically; stale-base or interruption failure preserves
recovery state.
Occurrences: baseline-01

Occurrence: gate-fastpath build — eleven ticket landings used the hand-run stale-base sequence and one diverged into a manual cherry-pick.

The exact authority, interruption recovery, and relationship to `bench commit`
and `bench worktree release` remain reviewer decisions.

Occurrence: FT123+FT124/gate-fastpath close — foreign capture edits forced isolated landing worktrees after correct attribution refusal.
The command owns preparing that path-scoped landing worktree and safe transfer.
Entry: `/bench-shape-idea`. Sources: gate-fastpath and FT123+FT124 retros, drained here.

Occurrence: FT129 implementation retro — a delegate sandbox could not fast-forward shared Git metadata.

Delegation preflight fast-forwards from the coordinator, verifies `HEAD == main`,
and routes sandbox-denied metadata updates back to the coordinator. Source: FT129
implementation retro, drained here.

Occurrence: FT126 recurrence-tallying retro — delegates reached their first gated action with an untrusted derived binary and paid the rebuild refusal.

Tool-readiness preflight verifies the assignment's derived binary and reports any
rebuild; an explicit opt-in may rebuild, while gate entry remains repair-silent.

Occurrence: FT128 implementation retro — ticket chains could not create or release worktrees at their integration tip after the spec-build lifecycle disappeared.

The landing decision chooses whether `bench worktree create` gains a base-ref
route for non-spec chains.

Occurrence: 2026-08-01 per-component-gate-scoping — a spec-branch build used unowned `git worktree add` because the CLI withheld the core start-ref parameter (`capture/IDEAS.md`, drained here).

Occurrence: 2026-07-31 reviewer-approved drain — concurrent writes twice triggered `bench commit` attribution refusal (`capture/learnings.md`, verdicted here).

The landing command owns a concurrent-writer sequence or a detectable single-writer
assumption, names the destination branch, and preserves verdict across transition.
Sources: covers-traceability and spec-build-lifecycle-preconditions retros.

The hook-created assignment pool either lifts its one-active limit or documents
manual `bench worktree create --request <uuid>` as canonical. The command resolves
the primary checkout before merge, reports bounded offending paths and precise
subject-closure reasons, and separates publication from cleanup refusal. Sources:
Pocock-guidance-doctrine, parallel-session-landings, and
spec-authoring-and-light-path retros.

Landing recovery accepts a listed assignment/digest selector or reports a resumable
request token; it never requires reconstructing a historical opaque input.

Occurrence: FT189 landing recovery — a persisted request digest could not replay the original request token.

**FT207 (MEDIUM, decision required) — worktree-mutating paths share malformed-admin refusal.** Decide whether every worktree-mutating Git call pre-scans malformed private admin entries through FT189's refusal owner before Git can block. Entry: `/bench-shape-idea`.

Occurrence: 2026-08-14 FT189 probe — add, lock, unlock, and prune call sites could hang on the malformed admin entries enumeration now refuses.

**FT199 (MEDIUM) — a recovery-aware branch-retirement coordinator closes one
repository-wide ref inventory.** The existing
cleanup paths act on one known target at a time; they do not classify every
non-default branch before a reviewer decides what may leave. Add `bench branches
retire [--discard <branch>...] [--apply <fingerprint>]` as the coordinator over
the existing worktree landedness and cleanup seams (the recovery verb and its
seam are removed). Its plan classifies
each non-default ref as active, landed, recover, review, or explicit discard.
Apply binds that complete inventory to one exact fingerprint, preserves dirty
unpublished work under recovery refs, deletes only refs still at their planned
OIDs, and refuses concurrent drift. It never auto-lands unique content or falls
back to raw `git branch -D`. The destructive classifications, policy for refs
outside registered worktrees, and relationship to FT98's per-payload discard
remain reviewer decisions. Entry: `/bench-shape-idea`. Source:
`capture/IDEAS.md`, drained here.

**FT178 (MEDIUM) — `bench worktree`'s bare verb is a human porcelain that
traps automation and leaks on signals.** The agent-facing surface is the
subcommands; the bare verb must print usage, while any retained subshell gets an
explicit opt-in name plus signal-safe release or lease reclamation.

Occurrence: 2026-08-01 reviewer ruling — three bare-verb probes hung and leaked two worktrees before cleanup (`capture/IDEAS.md`, `capture/learnings.md`).

Parser-first dispatch rejects unknown flags before choosing the bare command, and
the discovery convention remains canonical in `AGENTS.md`. Source: the
parallel-session-landings retro, drained here.

**FT172 (MEDIUM) — the roadmap parser and context snapshot make the drain's
non-recurrence evidence complete.** Contract the row grammar as exactly
`**FT<n> (…) — title.**`, with bold forbidden at body starts, or explicitly make
the parser absorb bold-led bodies; document and contract-test the choice.

The context snapshot emits last-drain commit, commits since, code-touch status,
complete row verification, and a mechanical discrepancies block for empty
`roadmap_id`, missing spec paths, and retired specs; the CLI judges no policy.

Occurrence: 2026-07-25/26 drains — workload boundary and prior commits were re-derived by hand before each reconcile.

Occurrence: 2026-08-03 drain — two staged specs had empty `roadmap_id` values while `occurrence_discrepancies` was empty (`capture/learnings.md`, verdicted here).

The phase invokes `bench roadmap --context` once and must use `--full` or stop
claiming completeness when source bodies are truncated.

Occurrence: 2026-08-03 drain — the default context truncated the journal and retro bodies needed by later drain steps (`capture/learnings.md`, verdicted here).

Promote a capture's symptom as evidence only after checking its diagnosis against
the tree; otherwise label the diagnosis as the capture's claim. Source:
`capture/learnings.md`, verdicted here; grammar source: `capture/IDEAS.md`.

Occurrence: 2026-07-27 drain — a promoted capture incorrectly claimed artifact tests used the graded root, causing a false start into `/bench-write-spec`.

The discrepancy check resolves whether `roadmap_id` comes from the roadmap path,
the spec's `Roadmap:` header (`internal/spec/spec.go`), or both; authoring must
either populate or stop advertising that field.

Occurrence: 2026-08-01 discrepancy audit — FT135 omitted its staged spec path and live specs lacked the header that supplies `roadmap_id`.

The occurrence-ledger migration check must derive counts from `ROADMAP.md` or
surface its pinned map as a discrepancy; it must not become a second hand-edited
source. The checked rows remain FT71, FT158, FT98, FT169, FT141, FT94, and FT125.

Occurrence: 2026-07-31 reviewer-approved drain — a hard-coded occurrence map caused a full gate after FT128 was reconciled out (`capture/learnings.md`, verdicted here).

**FT173 (MEDIUM, decision required) — AXI residual: the
active-assignment-with-deleted-tree disclosure class.** AXI remains scoped to
high-frequency query surfaces: `bench diff` owns the coherent snapshot, nested
`worktree list` owns its disclosure, turn-preventing principles outrank format,
and there is no parallel renderer, legacy mode, or `bench git` namespace.

Occurrence: 2026-08-12 axi-query-disclosure capstone — contextual disclosure, harness-log, registry conformance, and `craft-cli` guidance were composed across approved query surfaces.

R11 is an undecided QD6 hybrid: an active assignment whose tree was deleted.
`internal/worktree/list.go` renders active/missing and advertises path/exec
actions without a clean state. A reviewer must define the disclosure class before
the light-path change. Sources: axi-query-disclosure retro and `capture/IDEAS.md`.

**FT202 (MEDIUM, decision required) — a standing test-support fence, and the
census scope for process-backed fixtures.** Decide whether a named,
census-visible test-support commons fence in `projects/benchkit.md` replaces
per-spec amendments. Decide separately whether process-backed git fixtures count
in the ordinary-test census; the oracle change needs reviewer approval. Resolve
census scope before the pending R2 harness collapse, and keep `leasedRepo` as
explicit residue if it needs a third API shape. Entry: `/bench-shape-idea`.

Occurrence: axi-query-disclosure close — two scaffold amendment cycles repeated the missing-fence decision and an A/B trial supported a named commons fence (source: `capture/IDEAS.md`, drained here).

Occurrence: cmd/bench registry envelope audit — process-backed fixtures were outside the ordinary census and the pending harness collapse could hide them.

**FT185 (MEDIUM) — gate results join the structured Bench output contract.**
`bench gate` is the last major agent-facing surface that reports phase verdicts
as ad-hoc prose while `bench test`, `bench diff`, and `bench coverage` emit TOON.
Give the gate one structured result schema without changing exit-code authority,
phase completeness, or the durable verdict it authors. The gate-pipeline map's
ticket 9 closed the scope decision: no output redesign rides that pipeline build,
so this is an independent item. `bench commit` prints the same per-phase evidence
inline (`phase conformance: green`, `gate: green`), and one schema should own both
public projections. Entry: `/bench-write-spec`. Sources:
`capture/IDEAS.md`, drained here; the `injected-interface-junctions` retro,
drained here.

Green `bench gate` and `bench commit` runs currently stream a full per-package
transcript even though the durable JSONL log retains it. The shared projection
should show phase verdicts and the log pointer on green while preserving the
complete transcript on red, removing the incentive to truncate command output
at the caller. Source: `capture/IDEAS.md`, drained here.

Occurrence: 2026-08-14 CLI inventory learning — public worktree commands reached the CLI before the canonical inventory was updated.

**FT89 (MEDIUM) — guidance coherence and current-state documentation.** Make
every documented CLI example executable; parse and validate real YAML
frontmatter; derive the skills index and inventories from one implementation;
embed design-it-twice
briefs in complete delegation charges; and admit reviewer-approved assessment
findings as a legal synthesis origin. Use the canonical iteration-cap line
definition and only recommend shifts that meet the routing contract.
Every phase exit emits one copy-paste cold-session continuation prompt with the
exact harness-native command and pinned repository, branch, commit, spec/status,
and unresolved next action, so resumption never depends on conversation history.

Clarify shape termination and the no-design-source branch, remove stale paths
and inventory omissions, retire obsolete historical reports, dogfood
first-party authoring guidance, and pin normative external references. Rewrite ADRs and README claims to the behavior proved by artifact
contracts, including the actual canary phase selection and npm prepare shape.

Occurrence: 2026-07-25 `craft-line` frontmatter — an unquoted ` #` hid two trigger clauses from YAML and the grep-only gate (`capture/learnings.md`, verdicted here).

Occurrence: 2026-07-31 cold tree re-derivation — `craft-review` frontmatter advertised forbidden inline self-review (`capture/learnings.md`, verdicted here).

The coherence pass reconciles `.bench/BENCH.md` and `BENCH-reference.md` against
`bin/bench.sh`, re-derives stale decision-document sources from the live tree,
and folds Redact-style secret-safe excerpts into existing debugging guidance.
Sources: `capture/IDEAS.md` 2026-08-05/11 and `upstream(mattpocock/skills@84fdeff)`.

Occurrence: 2026-08-05/11 coherence drain — CLI inventories and decision documents still carried stale source paths.

Sources: `RR:S-06`, `RR:S-07`, `RR:S-08`, `RR:S-10`, `RR:S-11`, `RR:S-12`,
`RR:S-13`, `RR:S-14`, `RR:S-15`, `RR:S-16`, `RR:S-17`, `RR:S-18`; `RC:M-05`;
`capture/IDEAS.md`, drained here.

**FT106 (MEDIUM) — doc claims re-verified against the tree.** Invariant 3 tells
every session to write for the teammate who just walked in, and nothing ever
re-checks that what that teammate reads is still true: `CONTEXT.md`,
`projects/<name>.md`, and `docs/adr/` are asserted once and trusted forever,
while `bench structure` budgets only code. Add a probe step to
`/bench-what-next` — the drain is already the scheduled maintenance surface
with a batch-diff verdict mechanism to hang corrections on. Sample two doc
claims, verify each against the code, land corrections as visible batch-diff
entries with a one-line why (never as silent edits), and escalate the sample
on any hit — the escalating sample is what makes the cost self-scaling on a
clean tree and self-paying on a rotten one. Sample by staleness rather than at
random: `git log` gives last-touched dates for both the doc and the code it
describes, so a doc older than the code it claims to describe is the candidate
set. Paired with it, a `(?)` marker for inferred-and-checkable claims and
`(unverified)` for asserted-and-not-currently-checkable, written by
`/bench-setup-repo`'s exploration half and by `craft-adr`'s doc discipline, so
an adopting session's inferences stop reading identically to
reviewer-confirmed facts; the probe drains marked claims before sampling
unmarked ones. The two halves ship together — a marker nobody seeks out is
just decoration, and a probe with no self-declared targets samples blind.
Open: whether the probe target list is hardcoded or reviewer-declared
alongside `.bench/structure.budgets` (recommend hardcoded — a declaration file
is a second thing to keep current). Kit edit under the `craft-synthesis`
discipline. Background: `docs/reporesident-distillation.md` §1 and §5.

The same owner covers capture-style reports at commit time. A tracked report
that contradicts the tree is a documentation defect, so it carries status at
the top and is re-read when committed, not only when written. This stays a
discipline check: phrase-grepping project prose cannot demonstrate a reliable
bite. Source: the learnings journal, verdicted in a prior drain.

The canary planted-reason close adds a generated live-current-claim path census
to documentation reconciliation. Any claim path outside its charged fence is
sliced before a delegate starts, so a required residual cannot first appear
during the documentation edit. Source: the canary-planted-reason-ownership
retro, drained here.

**FT190 (MEDIUM) — every injected interface has a real-producer test or a
written exemption.** Enumerate injected interfaces and require either
real-producer composition or a recorded reason a fake is sufficient; decide at
build time between a conformance check and a one-time sweep.
Occurrence: 2026-08-03 injected-interface-junctions review — a fake-only junction hid missing real-producer coverage.

**FT191 (MEDIUM) — a fixture-and-seam inventory a charge can carry for free.**
Extend `bench outline` or a sibling reader to emit helpers, doubles, and
prior-art fixtures with `file:line`; decide at spec time whether this is a
projection or an existing-outline mode.

Occurrence: 2026-08-03 delegate retros — hand-built fixture inventories improved first-pass charge inputs.

**FT192 (MEDIUM) — one source per fact reaches spec and ticket prose.** Specs
name enumeration sources instead of copying implementation-derived counts;
Standards grades restated counts. The retired `Assumptions:` field and parser
are gone, but the one-source-per-fact rule remains. Kit edit under
`craft-synthesis`.

Occurrence: 2026-08-03 injected-interface-junctions review — a failure-message count disagreed across implementation, spec, and ticket prose.
Occurrence: recovery-discard build — ticket `Assumptions:` copied a standing tree-verification rule.

**FT206 (MEDIUM) — exact-candidate review sees destination metadata before it
freezes.** Reconcile destination-owned staged-spec metadata before exact review;
post-review behavior or acceptance changes still return to review and sign-off.
Kit workflow edit under `craft-synthesis`, entry `/bench-shape-idea`.

Landing reports the exact staged-spec reconciliation route when destination and
source bytes differ.

Occurrence: spec-authoring-and-light-path landing — destination and reviewed source had different staged-spec bytes, forcing fresh exact review.
Occurrence: gate-run-transaction landing — destination/source staged-spec reconciliation required a fresh exact review.

**FT204 (LOW, decision required) — one bounded transcript/session query.**
Agents repeatedly shape harness transcripts and session evidence with
`head`/`tail`/`awk`/`sort` chains; consider one bounded agent-facing query for
transcript and session census plus evidence projection. This is a new
operational surface, not an AXI query-surface extension, so it needs its own
operational-surface decision before any shaping. Entry: reviewer decision.
Source: `capture/IDEAS.md`, drained here (routed out of the
axi-query-disclosure harness-log ledger).

**FT205 (LOW) — `craft-delegate` names the delegate-worktree end-of-life
pair.** Release is the creating request's default end; clean-then-`bench
resume-clean` is the recovery pair for harness-created delegate worktrees. Kit
edit under `craft-synthesis`.

Occurrence: 2026-08-12 delegate end-of-life — clean alone left a recovered assignment and blocked the next create.

**FT58 (LOW) — hardened pool roots.** Permission failures on Bench-selected pool roots
should propagate — the tree currently asserts best-effort tighten
(continue-on-chmod-failure), a fork the build must put to the reviewer — and
non-owned or symlinked roots are neither rejected nor mode-revalidated after
creation.

Closure covers a permissive pre-existing directory, chmod failure, symlink
root, and crash-safe re-entry.

Sources: `RR:C-04`; `RC:M-01`.

**FT92 (LOW) — attributed subject drift and consumer-shipped input hygiene.**
"gate subject changed during execution" names no component; the drift message
should say what moved (the tree hash versus which declared manifest path) so
the next FT90-shaped defect self-diagnoses. The gitignored-declared-input
conformance check is benchkit-only; ship it as consumer gate scaffolding so
linked repos get the same protection.

**FT99 (LOW) — spec problem-premise verification.** Verify every “today the
code does X” claim in Problem, Solution, and Implementation against the tree at
spec time: name the command/check or mark uncertainty. `/bench-write-spec`
step 9 applies the falsification pass, while review audits call-site tables and
cross-products.

Occurrence: retired minimal-subprocess-data-exposure spec — its gate premise was stale when the build reached stage 1b.
Occurrence: retired cli-grammar-and-capability-evidence spec — an asserted `go test` behavior was false when stories were built.
Occurrence: 2026-07-27 ft91-phase-manifest-dag review — a red signal relied on `os/exec` behavior the mapped test could not trigger.

Review audits each call-site posture against the diff, enumerates
publication/refusal cross-products, and exercises immutable tracked-file,
transitioned-spec, untracked-descendant, and nested-CWD mutations.

Occurrence: FT86 review — two adopt call sites had reversed postures in its DefaultBranch table.
Occurrence: FT194 review — one direct-path refusal and one recovery success did not cover every publication/refusal pairing.
Occurrence: exact-prospective-landing close — immutable-snapshot and nested-CWD mutations were added to the first review slice.



**FT100 (LOW) — prose-weight pass on the kit's guidance surface.** Apply the
gate's "prove it bites" standard to prose: audit the craft-skill library and
the communication protocol so each skill and always-loaded clause cites an
observed failure it prevents (from the learnings journal or session evidence),
merge overlapping craft docs, and shrink the always-loaded `BENCH.md` rules to
demonstrated-delta clauses. Distinct from FT89, which fixes guidance
*correctness*; this row cuts guidance *weight*. The shipped communication cut
already took the "How to talk to me" slice, so what remains here is the
craft-skill library and the demonstrated-delta audit over the rest of the
always-loaded surface. Kit edit under the
`craft-synthesis` discipline; starts as a grill (`/bench-shape-idea`) because
the cut line on always-loaded rules is a reviewer decision.

The audit evaluates guidance against both loads an agent pays — bytes loaded
into context and cognitive branches introduced after loading — and treats
environment knowledge as cacheable only while the environment still owns it.
The proposed demand test (a behavior is complete only when a real session needs
it) joins the demonstrated-delta decision rather than becoming a second
completion rule. Source: `upstream(mattpocock/skills@84fdeff)`, drained from
`capture/IDEAS.md` here.

The prose-budget mechanism must not make wrap width a compliance lever. Before
lowering `.bench/BENCH.md` from its reviewed 175-line bound, choose a word-count
or house-wrap-normalized measure and decide which doctrine moves or consolidates;
a cosmetic reflow is not a weight reduction. Sources: `capture/IDEAS.md`
2026-08-11 and the Pocock-guidance-doctrine retro, drained here.

The skill-library audit also asks whether procedural guidance can state the
required end state instead: assess every step-by-step instruction for an
outcome-shaped replacement that preserves its constraint while removing an
unnecessary branch. This is a weight question, not a second correctness owner.
Source: `capture/IDEAS.md` 2026-08-14, drained here.

**FT101 (LOW) — per-context scope for monorepos: domain docs and profile.** A
monorepo has more than one bounded context, but the kit assumes one
`CONTEXT.md` and one `projects/<name>.md`. Both halves want one resolution.
The docs half: support a root `CONTEXT-MAP.md` pointing to per-context
`CONTEXT.md` files, add a single- versus multi-context question to
`/bench-setup-repo` Section C, and teach every `CONTEXT.md` consumer (phase
commands, skills) the layout. The profile half: let a profile declare the
paths it owns, resolve the active profile from the paths a change touches,
and scope the ambient surfaces (`bench outline`, `bench status`) to it — a
single-profile repo resolves exactly as today, so the change is additive.
Whole-tree `outline` and `structure` get less useful as a tree grows, which
is the context-cost-follows-project-size failure this row exists to fix.

A path-scoped summary preserves aggregate debt while making the ticket verdict immediate.

Scoping the *gate* is the contested part and starts as a grill
(`/bench-shape-idea`), because narrowing what the oracle runs is a reviewer
decision. Two guards frame it. A declared package boundary is not a diff —
FT91 ruled diff-scoped gating unsound here because contract and canary are
behavior contracts with no file→test map, and that ruling stands; a scoped
gate is legitimate only where the reviewer has declared the boundary, and a
change touching two profiles takes the whole-tree gate. And scope is not a
speed lever: the measured cost on this repo is phase contention and a
hardcoded `-count=1` (FT91's arms), not tree size, so any wall-clock win here
is a monorepo side effect and must never be the justification. Green must
keep meaning the same thing; a scoped verdict is explicit evidence, never a
silent skip.

One of three angles on the cost-follows-project-size complaint; the delegate-
slicing angle shipped as FT136 and the wall-clock angle is FT91, so this is the
last of the three still open. Kit edit under the
`craft-synthesis` discipline. Background and the alternatives considered:
`docs/reporesident-distillation.md` §8.

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency and
dogfood loops.** `craft-synthesis` consistency names the escalation policy for
every tier-spending edit; ticket shaping dogfoods its decomposition and rejects
widened or duplicated knowledge.

Occurrence: 2026-07-22 write-spec rerouting — a widened step triggered an automatic top-tier spawn past review.


**FT108 (LOW) — a refactor lane with a mechanical exit test.** The kit has no
refactor path: a pure restructure either gets forced through
spec → implement → review, where there are no stories and no red signal so the
coverage map is a fiction, or it takes the direct fix-and-gate path with
nothing but the gate protecting behavior. Add a `craft-refactor` skill — a
skill rather than a phase, because phases are reviewer-chosen entry points and
a refactor is usually a shape the work turns out to have. It composes
`craft-seams` (reaching a better seam is usually the point) and `craft-tdd`
(characterization tests are TDD with the red signal running backwards). Four
rules: tests covering the affected behavior green before any move, with
characterization tests written first where the behavior is uncovered; an
ordered list of mechanical moves each leaving the repo green; one move at a
time with every caller enumerated by search rather than from recall; and the
exit test — the suite passes with test logic unmodified, mechanical renames
being the only permitted test edit, so a changed assertion means changed
behavior and the move reverts and reroutes to the feature path. Carry the
no-bundling rule too (a bug found mid-refactor is parked and fixed
separately) — invariant 4's smallest-diff rule aimed at a failure mode the kit
does not currently name. Propose the assertion rule as a gate check
(`bench diff --assertions`) only if the skill gets used and the rule gets
violated. Kit edit under the `craft-synthesis` discipline. Background:
`docs/reporesident-distillation.md` §2.

**FT111 (LOW) — provenance tags that outlive their specs.** Remove `FT<n> story
<n>` tags only while editing their line, reject new ones in review, and update
`craft-review`/`craft-comments`; do not sweep existing sites.

Occurrence: 2026-07-23 retired-spec review — dangling comment tags pointed at specs that no longer existed.

**FT112 (LOW) — an approximation that stays green is not a cleared bug.** A
green proxy only narrows a hypothesis; load- or environment-sensitive failures
must be reproduced through the accused command under exposing conditions before
another stand-in is trusted. Add this reproduction-economics rule to
`/bench-debug`.

Occurrence: 2026-07-23 trustworthy-gate-under-load diagnosis — synthetic load shapes stayed green while real `bench gate` host load reproduced the failure.

**FT113 (LOW) — `bench commit --spec` residuals: the flip counts as a path, and
the flip has one author.** The reduced gate now includes `specs/` and
content-addressed ancestor evidence; `bench commit --spec <slug>` owns the
reviewed `Status: staged` → `implemented` transition, counts it as the required
path, requires `-m` plus an explicit path, and makes other routes refuse.

Occurrence: FT131/decision-map-integrity landing — staged-spec drift left a strong-stale verdict until the reduced allowlist shape shipped.
Occurrence: FT128 implementation retro — multiple commands could race to author the same status transition.
Occurrence: axi-coherent-diff retro — abbreviated final-landing usage omitted required `-m` and path arguments.



**FT130 (MEDIUM) — a capture write mid-lifecycle voids or blocks the run.**
Capture writes during a gate or gated commit are queued after the run or refused
by the subject lock; choose that mechanical queue/refuse behavior so green
verdicts cannot be voided.

Occurrence: 2026-07-25 FT122 gated commit — `bench idea` changed `capture/IDEAS.md` inside the gate window and invalidated a green run.


Occurrence: 2026-08-02 gate window — a `bench idea` write self-rejected an otherwise green fresh gate by changing its subject.


**FT138 (LOW) — instrument Bench so build economics are measurable.** Keep LOW
until a cost decision needs measurement; candidate metrics cover
delegate/coordinator tokens, findings and rework by axis, gate/spec iterations,
tier use, and phase/package timings. Decide storage, ownership, granularity,
retention, and audience through `/bench-shape-idea` under `craft-synthesis`.

Occurrence: 2026-07-25 reviewer pricing — economics instrumentation was a nice-to-have until an acceptance trigger requires it.

**FT164 (MEDIUM) — repair-lane charges, and a done-claim that resolves its
named owners.** One `craft-delegate` visit governs repair/experiment charges
and done claims. Each candidate names base commit, file fence, effort, focused
suite, and independent probe; setup failures stay separate, behavior/probe
checks decide tier, and style never justifies an expensive tier. Currency
repairs name documentation and wrapper/router owners; shared helpers state
cache opt-in or hermetic rows and run focused failures before the full gate.

Done claims resolve every Red-mutation owner to a real tree artifact; repair
rounds pin any extended enumeration. Structural acceptance uses semantic
enumeration or multiple representations, and mutation probes independently
verify that the intended mutation is present.

Slice tickets by runtime projection, canary selection, and documentation
conformance. Reuse an existing ticket for an extended acceptance row or require
an independent green tracer; triage must-fix production defects versus
judgment-priced hardening before charging repairs. Accept probe evidence only
after execution through `bench test` or `go test -v -count=1`; cached or
near-zero runs are invalid, and caveats must match charged mechanisms.

Registry-entry tickets derive fences from the previous entry's landing and trace
every new-name crossing. Scope repairs from the invariant, reread each repaired
artifact whole, target false claims in red-mutation checks, and give every
prose-only phase step a `Contracts:` landing site.

Lifecycle assignments use candidate ancestry while ordinary delegations may
require ff-to-main. Headless wrappers distinguish edits from stalls; review-only
normalization is folded into substantive tickets or deleted; repair tickets name
discovery, integration mappings, public mutation ownership, and
post-authorization descendants.

List delegate-supplied `-run` alternations before execution, carry sibling
classifier registries and pre-promote focused runs, use fresh round-unique
backups with owner restoration/deletion, sweep orphaned callers, mutate the
exact production dimension, and count approved tier escalation as a repair
round. Conformance repairs distinguish synthetic from live roots and reread
aggregate candidate versus slice base. Kit edit under `craft-synthesis`.

Occurrence: 2026-08-03 FT164 ticket-contracts build — a done-claim named a mutation owner that never landed and an enumeration drifted across repair rounds (sources: retired spec and implementation retro).
Occurrence: gate-test-concurrency close — one function-name audit could miss another representation of the same carrier (source: `gate-test-concurrency` retro).
Occurrence: 2026-08-12 guidance probes — wrap-spanning mutations became no-ops until presence was checked (sources: Pocock-guidance-doctrine retro, `capture/learnings.md`).
Occurrence: check-level-conformance-scoping close — one ticket mixed runtime, canary, and documentation owners (source: retro).
Occurrence: FT194 close — a repair extended an existing acceptance row and required reuse-or-new-tracer guidance (source: `ft194-project-green-desync` retro).
Occurrence: 2026-08-03 lifecycle run — triage, probe execution, caveat, registry-fence, reread, and prose-step gaps surfaced together (source: `capture/learnings.md`).
Occurrence: ft156-anchor-registry close — lifecycle assignments were incorrectly charged with a HEAD-equals-main rule (source: retro).
Occurrence: exact-prospective-landing close — wrapper stalls, normalization records, integration mappings, and post-authorization descendants were missing from charges (source: retro).
Occurrence: covers-traceability and FT195 closes — `-run` alternation, classifier-registry, and focused pre-promote evidence were absent (sources: retros).
Occurrence: 2026-08-11 preflight/lifecycle-removal closes — mutation backups and orphaned-caller sweeps were not charged (source: retros).
Occurrence: axi-coherent-diff close — repair evidence changed a correlated dimension or erased an escalated failed charge (source: retro).
Occurrence: canary-planted-reason close — synthetic/live-root and candidate/base rereads were missing from conformance repair (source: retro).
Occurrence: gate-run-transaction retro — scorecard attribution separates delegate rework from spec-origin findings.













**FT200 (MEDIUM, decision required) — make preflight mechanical at the landing
chokepoint.** `bench preflight build|review <slug>` now provides the phase-entry
start-oracle over base currency, authorized paths, ticket row ownership,
producer-derived membership, and non-empty review diffs. The removed lifecycle
moots its repair-assignment and receipt machinery; what remains is the
large-organization posture requested after the command shipped. Decide whether
a branch claiming a staged spec must pass preflight in the gate, pre-push, or
another single landing chokepoint, with a red blocking landing and without
turning the start-oracle into a second done-oracle. The enforcement must name
how it discovers the claimed spec, handle multiple staged specs, and remain
optional only through an explicit project-profile decision. Entry:
`/bench-shape-idea`. Sources: the bench-preflight spec's reviewed out-of-scope
cut and `capture/IDEAS.md` 2026-08-11, drained here.

**FT165 (LOW) — fold the domain-modeling discipline into
`/bench-shape-idea`.** Grill tickets challenge overloaded terms with concrete
edge scenarios and update `CONTEXT.md`/ADRs under one-source ownership;
integrate the moves into the existing phase under `craft-synthesis`.

Occurrence: upstream domain-modeling refinement — active counterexamples and language-owner updates were added while unrelated dispatch ideas were dismissed (source: `capture/IDEAS.md`).


**FT180 (LOW) — a spec-optional route decided at shape-idea's exit.** Shape-
idea routes small-but-wider work to a ticket-only folder without `spec.md` or
to a full spec; the build decides the threshold and lifecycle semantics. One
owner indexes ticket-only receipts or removes the folder after durable
promotion, and `$bench-implement-spec --full` carries the chosen route to
push-ready state.

Occurrence: 2026-08-01 reviewer request — ticket-only receipts lacked history and retirement handling while the light path was shaped (sources: `capture/IDEAS.md`, `capture/learnings.md`).



**FT182 (LOW) — a Planned-phase receipt over an absent target wedges the abandon
retry.** `bench worktree` resume must clear an absent-target Planned receipt
instead of returning `errStaleFingerprint`; keep the fix beside absent-target
planning in `internal/worktree/resume.go`.

Occurrence: 2026-08-02 FT176 plan-absent-target delegate — the crash window between receipt write and `checkpoint(Removing)` wedged abandon retry (source: `capture/IDEAS.md`).

**FT197 (MEDIUM) — the Go core owns gate invocation and process lifetime.** The
shell entry currently hops from `bin/bench.sh` through `.bench/gate.sh` before
the Go-owned runner can bind the operation. Wrapper-only termination has left
orphaned gate process groups, and worktree-to-main execution makes verdict
locality harder to state. Move invocation, waiting, signal propagation, and
terminal attribution into one Go-owned path while keeping `.bench/gate.sh` as
the project-authored check body. Preserve the public gate contract and every
phase; this changes the process owner, not the oracle. Entry:
`/bench-write-spec`. Sources: `capture/IDEAS.md` and the covers-traceability
retro, drained here.

**FT166 (LOW) — `bench capture commit`: porcelain for the ambient capture set.**
Commit `capture/learnings.md`, `capture/IDEAS.md`,
`capture/session-handoff.md`, and `capture/retros/` with a conventional message
plus empty-index and explicit-path safeguards; compose over FT168's allowlist.
Distinctly, provide a sanctioned writer for hand-authored learnings/retros that
uses the parser's canonical open-entry grammar, while post-promotion retros
remain reviewable through ordinary landing.

Occurrence: 2026-08-01 capture-path review — the allowlist clause folded into FT168 and dirty unrelated changes required bundling with a named ride-along (sources: `capture/IDEAS.md`, `capture/learnings.md`).
Occurrence: 2026-08-09 cancel-signal session — a malformed `- date` journal bullet made parser-backed readers report empty (source: `capture/learnings.md`).






**FT168 (LOW) — focused iteration evidence: a fixture-selecting canary.** Add
`bench canary` for one named fixture or family as iteration evidence; the full
sweep remains the only gate oracle. Distinct face: a reduced phase set for
allowlist-confined changes shipped separately, while focused canary remains LOW
and out of that spec's scope.

Occurrence: light-path repair pass — one race fixture required whole-sweep evidence and invited an unnecessary duplicate run (source: retro).
Occurrence: 2026-08-01 reduced-gate-phase-set landing — the reduced phase set shipped, leaving focused canary open (source: retro).

**FT140 (LOW) — review residuals that want a verdict, not a build.** One
reviewer decision resolves provenance for three real-coverage tests from
retired specs. Four bounded decisions remain: orphaned-review pickup mapping,
`dedupe` dead code, the 2 MiB learnings bound, and index-lock retry; the
gate-decision seam also needs projection/helper/counter/partial-decision/
nil-error/mixed-reload hardening.

Occurrence: FT86/FT91 resolution reviews — real coverage tests lacked rows in retired maps (source: review records).
Occurrence: 2026-08-01 orphaned-review pickup — `reviews/` and `specs/<slug>/spec.md` were not both present (source: `capture/IDEAS.md`).
Occurrence: 2026-08-02 learnings/index-lock review — a lower read bound changed behavior and a held index lock exposed retry policy (source: `capture/learnings.md`).
Occurrence: 2026-08-07 gate-decision review — shipped test code retained a bounded residual bundle (source: `capture/IDEAS.md`).



## False greens — verdicts that credit unchecked work

Five rows, one failure class: a green whose warrant is missing — a stale
binary, a dead or skipping citation, a vacuous baseline, an unchecked absence,
a dependency edge nothing resolves. Each hardens a different oracle surface, so
they stay separate builds, but they read and prioritize as one theme.

**FT133 (MEDIUM) — `bench coverage --check` verifies that red-signal citations
resolve.** Resolve each cited test and require it actually executes; capability
skips fail closed. Emit stable row identity, and validate membership across
every tag with a mixed-tag fixture. Prefer mechanical checks over duplicated
authoring instructions.

Occurrence: FT86 story-16 review — a lowercase subtest citation could never go red.
Occurrence: 2026-07-26 conformance review — `TestRootConformance` matched but skipped without `BENCH_CONFORMANCE_ROOT`.
Occurrence: 2026-07-28 Codex falsification — shared story/seam/red-signal fields could not identify individual rows (source: `capture/IDEAS.md`).
Occurrence: 2026-08-11 bench-preflight review — mixed tags let membership escape the validator (source: retro).
Occurrence: FT126 recurrence — a scoped conformance command again omitted `BENCH_CONFORMANCE_ROOT` and printed `ok`.





**FT174 (MEDIUM) — ticket files have one enforced dependency, ownership, and
mutation grammar.** Parse `Blocked by:` by sibling identifier with
cycle/dangling checks; require explicit `Ownership fence:` and complete
Red-mutation rows per acceptance ID, reject duplicates, and keep historical
closed tickets exempt. Retire `Assumptions:`; refuse probes outside the fence;
compare `Closure:` with `Contracts:`, coverage, and edge inventories before
leasing. If a fence is wrong mid-run, the whole-run abandon path remains the
residual.

Every cross-fence fact carries declaration, resolving owner, exported value,
consuming path, and reciprocal edge. Review counts exact-anchor occurrences
before deletion, gives every closure member its mutation, checks lifecycle before
dependency blocking, and treats unchanged consumers as integration surfaces.
Preflight rejects malformed covers/contracts/fences, expands compound rows into
independently removable members with retained fixtures, assigns canary-universe
ownership, revalidates spec-ticket breakdowns against the tree, and includes
cross-line mutation fixtures.

Occurrence: 2026-07-31 FT164 template review — dependency and ownership grammar was reverified against the landed parser.
Occurrence: 2026-08-03/04 ticket parser work — fences, Red-mutation rows, duplicate IDs, `Assumptions:`, and mid-run re-fencing gaps surfaced (sources: captures).
Occurrence: conformance-harness close — cross-fence ownership needed a complete declaration-to-edge chain (source: retro).
Occurrence: one-shot close — exact-anchor counts were needed before promised deletion (source: retro).
Occurrence: repair-ticket-reslicing close — closure members, lifecycle state, and integration surfaces needed separate review (source: retro).
Occurrence: spec-ticket-handoff close — read-only preflight and compound-row expansion were missing (source: retro).
Occurrence: spec-authoring-and-light-path close — canary-universe changes lacked an explicit coverage owner (source: retro).
Occurrence: 2026-08-13 staged-spec handoff — breakdown approval moved earlier, with implementation-time tree revalidation (source: `capture/IDEAS.md`).
Occurrence: bench-preflight close — line-oriented grammar needed explicit cross-line hostile variants (source: retro).










**FT177 (MEDIUM) — a stale `dist/bench` makes contract-test mutation probes
silent no-ops.** Contract CLI tests must reject `$k/dist/bench` older than
tracked Go sources; stale bytes can only yield a false PASS. `bench commands
--brief` must likewise verify source/binary identity or refuse with a repair
action.

Occurrence: reduced-gate-phase-set retro — source edits were tested against stale `dist/bench` (source: `capture/IDEAS.md`).
Occurrence: spec-build-lifecycle-preconditions retro — FT176 freshness required manual rebuilds before `bench commit` (source: retro).


**FT103 (LOW) — existence-checked absence evidence: the gate half.** Check that
each consumer-payload allowlist source exists; FT85 closed empty-set vacuity,
while the per-path guard remains. The charge-side identifier-resolution rule is
already in `craft-delegate`.

Occurrence: 2026-07-26 delegation-discipline landing — a misspelled kit-only allowlist let absence evidence pass vacuously (source: retro).

**FT201 (LOW) — production cancel-signal registrations conform to one source.**
Check that production `signal.Notify`/`NotifyContext` uses
`subprocess.CancelSignals`; exclude `_test.go` fixtures. The migration is
deferred as a new gate rule, and the remaining reviewer decision is whether
`Pdeathsig` closes SIGKILL-orphaned builder groups.

Occurrence: 2026-08-09 cancel-signal parity — one source was established but detached SIGKILL process groups remained (sources: `capture/IDEAS.md`, `capture/learnings.md`).

## Reds the diff doesn't own — inheritance, load, and harness defects

Five rows, one failure class: a red that answers for something other than the
diff in front of the gate — an inherited baseline, machine contention, a
literal deadline, a harness defect, a flaky oracle.

**FT141 (MEDIUM) — `bench gate pin` records red verdicts, so inherited reds stop
reading as caused.** Pin failed checks with the baseline commit so inherited
reds subtract; refusal retains phase and full gate log, prepush/status explain
the reviewed tree, and interrupted gates preserve prior `.git/bench-last-gate`
identity.
Occurrences: baseline-01
Occurrence: 2026-07-26 FT91 inherited-main red — pinning the baseline avoided causation archaeology (source: `capture/IDEAS.md`).
Occurrence: gate-fastpath close — refusal output had discarded phase names and full logs (source: retro).
Occurrence: 2026-07-29 FT150 adoption — pin refusal surfaces quoted the command without explaining reviewed-tree state (source: retro).
Occurrence: 2026-08-07 interrupted gate — pending startup overwrote the prior diagnostic record (sources: captures).




**FT104 (LOW) — load-induced commit refusals: the stop rule and the pre-gate
quiet check.** After two refusals by a known flaky test proven green in
isolation, stop and hand the reviewer evidence instead of retrying. Before any
aggregate gate, canary, or `gate-phases` launch, ensure returned delegates have
no live tests and keep the coordinator-owned resource serialized; prefer
mechanical refusal/warning over prose.

Occurrence: FT85 review-fix landing — repeated load-coupled refusals passed unchanged on a quiet machine (source: learnings).
Occurrence: 2026-07-27 delegate completion — a reported done claim left `go test` and wait loops running, causing a false timeout (source: `capture/learnings.md`).
Occurrence: 2026-07-29 light-path close — a repair delegate launched full gate and canary concurrently (source: retro).




**FT115 (LOW) — load-robust test and phase deadlines derived from bounds.**
Derive `waitForPIDFile`, FIFO, and sibling waits from `bounds.TestDeadline`;
give long conformance phases explicit headroom or distinguish timeout verdicts
from assertions.

Occurrence: gate-fastpath and FT123/FT124 retros — deadline ordering and SIGINT liveness failed only under combined load (source: retros).
Occurrence: first full gate on a 12-online-core host — fixed FIFO bounds exhausted under load (source: AGENTS/CLAUDE contracts).




**FT120 (LOW) — gate, canary, and contract test-harness defects nothing
asserts.** Bound R12's release-file wait against a stated red signal and
guarantee process-group teardown even after a killed test; assert per-fixture
`BENCH_CANARY_PHASE` isolation under concurrent families; pin self-host child
teardown before `TempDir` cleanup. Register conformance checks that
whole-package prefix matching currently hides, and give worktree test shell-
outs deadlines and process-group teardown around FIFO fixtures.

The preflight harness must include nested-special-file refusal and diagnose exact
release-root fixture commands; focused mutation probes need a clean root entry
point; the ordinary dev gate must run the existing live-root sweep so stale rows
fail before release.

Doctor-shim process-launch failure remains a distinct harness diagnostic, with
launch-versus-child-status attribution preserved instead of collapsing into the
child exit code.

Occurrence: 2026-07-23 R12 contention fixture — an unbounded TempDir release wait left an orphan shell (source: `capture/IDEAS.md`).
Occurrence: 2026-07-26 self-host contract — TempDir cleanup raced a child still writing (source: retro).
Occurrence: 2026-08-09 worktree FIFO review — bare git shell-outs could orphan children when a test died (source: `capture/IDEAS.md`).
Occurrence: FT126 retirement gate — a doctor-shim process-start failure collapsed to exit 1 before unchanged cases passed.
Occurrence: 2026-08-11 bench-preflight close — nested FIFO and injected-port root gaps remained (sources: captures).
Occurrence: Pocock-guidance close — direct live-root probes emitted unrelated diagnostics without a clean entry point (source: retro).
Occurrence: 2026-08-14 spec-authoring review — ordinary `go test ./...` skipped `TestRootConformance` and left stale live-tree rows (source: `capture/IDEAS.md`).







## Standards debt — one batched light-path pass

Three rows plus FT142's standards track are shippable together as small
one-source-per-fact and cleanup sweeps under one gate; FT117's parser-routing
half is the largest item in the batch. FT142 itself stays on the main list
because its ship track belongs to a separate `prep-release` hardening visit.

**FT117 (MEDIUM) — FT87 parser-surface follow-ups.** Route `internal/spec`'s
`specArg`, `worktree list`, and `internal/adopt/doctor.go` through `usage.Parse`;
correct/grade `whyNested`, keep `cmd/bench/main.go` worktree dispatch exempt,
and flatten nested `bench commit` usage errors.

Occurrence: FT87 slice-3 grammar review — parser-routing leaves and nested usage output remained after centralization (source: retro).

**FT179 (MEDIUM) — comment quality: strip the reviewer-facing register,
document high-stakes surfaces, sharpen `craft-comments`.** Remove provenance
IDs, audit arguments, port narration, and in-file mutation transcripts while
keeping behavioral constraints. Document releaseevidence/preflight/gate/
contract/worktree APIs and `bench.sh` dispatch; update `craft-comments` with
identifier-provenance and one-source rules, “state the constraint,” sparse-file
qualification, and commit/spec red-record ownership. Reopen FT111's edit-in-
place rule for the larger measured scope.

High-stakes comments state contracts: `ParseTicket` documents its grammar and
contract rather than merely naming the consumer that calls it.

Occurrence: 2026-08-01 comment sweep — reviewer register and high-stakes undocumented surfaces were found across the kit (source: `capture/IDEAS.md`).
Occurrence: bench-preflight close — `gather_test.go`, `RenderError`, and ticket-ID density remained as comment debt (source: retro).


**FT94 (LOW) — single-sourced `bench resume` summary golden.** Extract one
shared expected-format helper across unit and runtime-binary seams; keep the
expectation literal single-sourced under the one-source-per-fact rule.
Occurrences: baseline-01

## Session tax — evidence-supplied reader rows

This row is a measured, recurring reader cost from the week-of-2026-07-19
transcript evidence and builds on surfaces that already exist.

**FT125 (LOW) — reader surfaces that return the slice, not the file.** Add
section-scoped `bench spec show <slug> [--section stories|coverage|status]` and
`bench outline --symbol <name>` body/context output; validate that each narrows
actual session reads, with separate seams.

Occurrence: week-of-2026-07-19 session evidence — spec and source readers required repeated whole-file slicing.
Occurrences: baseline-01

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

**FT6 (LOW, parked pending evidence — leave parked):** `bench refs`, `bench
detect`, `bench doc`, `bench specs --retired`, doctor binary-presence row,
`conformanceFamilies`-vs-dispatch reconcile meta-check, and a per-anchor
bite-proof meta-test (canaries prove one needle per family today; graduate on
observed anchor rot). `bench symbols` is not carried; restore only if agents
demonstrably burn turns on symbol search. Also parked here 2026-07-23:
concurrent `bench upgrade` runs, raised as a coverage gap by the FT85 review
and closed by decision rather than left open — `transactionalLink` already
moves tree, manifest rows, and version stamp together, so the damage is
bounded; graduate on an actual report of two upgrades interleaving badly, not
before. Two upstream candidates join the pending-evidence tier rather than
becoming new skills: a generated human-procedure wizard plus last-resort HITL
loop graduates on a Bench workflow that cannot be made agent-operable, and a
third-party questionnaire graduates when a real decision map blocks on someone
other than the reviewer and reviewer-directed grilling misroutes the question.
Source: `upstream(mattpocock/skills@84fdeff)`, drained from
`capture/IDEAS.md` here. Also parked here 2026-08-12: a one-off `bench commit`
refusal — `gate: red` / `prospective authorization refused: inherited` — on a
tree whose immediately following direct `bench gate` was green (the retried
commit landed on the fresh verdict); not reproduced across seven later
landings, and a repro through anything but `bench commit` itself proves
nothing. Workaround on recurrence: run `bench gate` directly and retry.
Graduate on a second reproduced refusal through `bench commit`; the
verdict-class plumbing around `inherited` records is the suspect. Source:
`capture/learnings.md` 2026-08-12, verdicted here.

**FT24 (parked pending upstream) — Codex agent-line guard parity.** Researched
2026-07-11: still not implementable on current Codex — delegation has no
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` cannot
stop the subagent. The current surface verdict is canonical in
`.bench/BENCH-reference.md` Hook Layers. Graduate only when the Codex changelog
adds a spawn tool name or a deny-capable SubagentStart.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

**FT38 (LOW, decision required) — dashboard visual identity pass.** `bench
dashboard` v1 shipped data-faithful and visually neutral; the original idea
wanted a rich treatment with animated characters, reference saved at
`ui_example/` (Gather-style pixel office with activity feed). The minimum
tabling period has elapsed; revival remains a reviewer decision and starts as
a grill (`/bench-shape-idea`). Decision detail is recoverable via
`bench spec history dashboard`.

**FT170 (LOW, decision required) — behavioral red/green evaluation for skill
guidance.** The scheduled revisit date has arrived. Before Bench adopts a
prose-heavy skill-testing workflow, prove the need through a benchmark substrate: choose
one narrow behavior with deterministic artifact assertions or a blinded
scoring rubric; run repeated no-guidance and candidate-guidance trials in fresh
isolated contexts with pinned model and effort; keep authoring cases separate
from held-out evaluation; and report variance. Run the benchmark advisory-only
during skill changes or assessment, never in the deterministic per-commit gate.
Only stable improvement earns the smallest necessary `craft-skills`
requirement and harness pointer. Entry: `/bench-shape-idea`. Source:
`capture/IDEAS.md`, drained here.

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

1. `/bench-shape-idea` — FT207 decides whether worktree-mutating paths share FT189's malformed-admin refusal before Git can block.
2. `/bench-write-spec` — FT185 can make gate and commit results one structured, concise public projection.
3. `/bench-shape-idea` — FT162 resolves durable interrupted-review recovery and assigned-worktree authoring after FT169's landing primitive.
