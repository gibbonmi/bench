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

**FT171 (MEDIUM, decision required) — bound outer gate-phase concurrency
against measured contention.** The artifact-split follow-up measured the
`posture` package materially slower inside the fresh full gate than in its
focused run while other outer phases overlapped. The existing gate-concurrency
decision deliberately left outer phases uncapped until that evidence existed;
the gate-critical-path map now records that the trigger fired but grants no
authority to choose a cap. Price candidate outer widths with repeated full-gate
measurements on the same tree, recording the exact commit, worktree, and run
time beside every sample. Measure enough repetitions to state the variance;
discard a sample whose subject does not match the build under test rather than
reporting it. Preserve every phase and unchanged green semantics. Scoped
gating remains ruled out and is not a speed lever. Entry:
`/bench-shape-idea`. Sources: `capture/IDEAS.md`, drained here;
`decisions/gate-concurrency.md`; `decisions/gate-critical-path.md`; the
2026-07-30 Claude usage-report assessment, drained here. A first full gate on
a 12-online-core host supplied the lower-width sample 2026-07-31: both
`TestSetupConflictContracts` FIFO cases exhausted their 15 s subprocess
deadlines while the outer phases overlapped, then the exact pair passed
focused once and 3/3 repeated at about 0.43 s per run. Treat that delta as
contention evidence for pricing the outer widths, not as authority to pick one
from a single machine.

The first demand-reduction slice (the `gate-decision-test-seam` spec) landed
2026-08-07: the exhaustive public-document mapping matrix now runs at the
read-only decision seam with representative full-engine composition proofs
retained. Return to the remaining FT171 shaping decisions from here.

The second demand-reduction slice (`conformance-harness-scope`) also landed
2026-08-07: the 83 direct conformance fixture bites now run only their resolved
ordinary check while the full-table controls remain broad. The post-reduction
census then resolved decision #20 on exact commit `eb6845f`: `internal/gate`
remained a strictly serial 140–160 s focused package, `internal/specbuild`
remained a 58–60 s serial package, and one fresh gate peaked at 97 concurrent
descendants. Decision #22 selected intra-package test concurrency before any
outer-width pricing.

The `gate-test-concurrency` slice landed 2026-08-08: 192 of 245 top-level gate
tests now run in parallel while 53 reasoned serial tests retain their process-
global constraints. The exact-candidate package median fell from 150.85 s to
56.72 s, but width one and width two still wrote roughly 1.26 million filesystem
blocks each. Its measured materialization follow-up then shared the immutable
fixture binary and narrowed setup-only work: width one fell to 111.67 s and
699,904 output blocks; width two fell to 78.79 s and 699,712 blocks. The flat
write volume across widths means concurrency compresses wall time without
removing intrinsic fixture work. Resource-constrained follow-ups retain output
blocks beside wall time and state serial coordinator/delegate cadence as an
explicit unused-capacity reason when chosen. Decision #23 and its measured
residual are complete. Lifecycle removal moots #24's specbuild parallelism;
#25's remaining sized serial cuts and #26's exact post-route census stay open
before #8 prices candidate outer widths, and the decision map first removes
#24 from #26's blockers. Sources: the `gate-test-concurrency` retro, drained here;
`decisions/gate-budget.md`; the retained FT171 implementation receipts.

The branch-native gate rebuild supplies a timing-discipline instance: compare
before-and-after runs only on the same subject and cache posture. Its 73.850 s
prospective run identifies the ordinary test driver as the current dominant
cost, while its 354.073 s predecessor remains directional rather than a speed
claim. Source: the branch-native-build-test-architecture retro, drained here.

**FT162 (MEDIUM) — full-run and phase-close state has one authoritative subject
and handoff.** Recommendations from the craft-tickets, light-path,
artifact-suite, and artifact-hoist retros converge on one lifecycle owner. The close's
coordinator gated `main` while the final ticket sat uncommitted in its
assignment worktree — the aggregate dashboard signaled dirty work but nothing
identified the authoritative final-check tree, so the first green verdict
answered for the wrong checkout and cost a full redundant gate. Three edits:
a final-check entry check that resolves active assignment worktrees and
states the exact oracle subject before any gate starts; a subject-oriented
`bench` diagnostic exposing an active dirty assignment before a
primary-checkout gate begins (the CLI half, Go work); and the implementation
handoff naming any still-open assignment and whether its branch is committed.
The implementation-retro instructions also gain
an explicit codification-candidate pass — inspect the session for repeated
ad-hoc checks, decision procedures, or reconstructed logic worth codifying,
each candidate naming its session evidence, proposed durable owner (CLI,
skill, gate, or process prose), and expected effect. Concurrent main-tree
writers need a visible intent or lease signal — that close repeatedly polled
process state and delayed gates just to learn when another session finished
landing on `main`; it is the same subject-visibility surface as the CLI
diagnostic. The spec-backed half of the first edit was answered by
`spec-integration-gate-cadence`'s promotion gate, since removed with the
provisional spec-build lifecycle (a reviewed spec-backed build now lands
through serial `bench commit`), so what stays open here is the light-path and
non-spec close, where the final-check tree is still resolved by hand.

The handoff storage question joins the same lifecycle decision. A single
repository-level `capture/session-handoff.md` can be clobbered by concurrent workstreams;
per-spec handoffs would isolate those pins and retire with their spec, but would
leave non-spec work needing a repository-level owner. Decide whether the
authoritative handoff remains singular, moves into `specs/<slug>/` for spec-backed
work, or becomes a generated projection over per-workstream state rather than
adding a second handoff convention by accident. Source: `capture/IDEAS.md`, drained here.

The concurrent-writer clause gained a data-loss face 2026-08-02, carried as the
capture's claim rather than re-verified fact: while one session held finished
uncommitted work in the primary checkout, the pcgs session's landing of
`5d67654` left the tree fully clean — the dirty paths were gone with no stash
and no recovery ref, and the work survived only because its content was still
in the other session's context. The claim's graduation trigger is naming what
that landing ran and reproducing the discard through it. Until then the
standing posture is the pcgs retro's serialization rule: two sessions sharing
one checkout serialize on every lifecycle mutation and every landing, side-work
belongs in a worktree even for main-branch landings, and a landing path that
meets another writer's dirty paths refuses or sets them aside recoverably
(FT98's primitive), never silently discards. Sources: `capture/learnings.md`
2026-08-02, verdicted here; the per-component-gate-scoping retro, drained here.

The full-run half moves `--full` orchestration from phase prose into a `bench`
subcommand only if the decision map rules that the harness-independent
substrate should own it. That command records the staging base and terminal
implementation SHA as first-class review inputs and lets `bench diff --full`
accept that explicit base/head range; the artifact-suite and decision-map
integrity runs both had to reconstruct the range after the default-branch
review surface reduced to a single landing commit or the handoff path. The same
subject record binds any performance number in the close to the exact commit
and worktree that produced it; a number from a different subject is omitted,
not carried forward as evidence. Entry: `/bench-shape-idea`.

The close itself has one unresolved ownership decision. `bench spec retire`
says to remove the roadmap row while `/bench-final-check` says to leave it for
`/bench-what-next`; one instruction yields. The same decision gives retirement
a planned promotion/removal manifest and an explicit commit path instead of
generic dirty state. It also owns a close-order precondition that does not exist
today: `bench spec retire` deletes the live spec folder outright
(`internal/spec/spec.go:400`) without checking whether terminal final-check and
retro capture have run. The check-level-conformance-scoping close retired first,
then found that terminal status still required the deleted path. Until the
retained-run fallback exists, final-check and retro capture precede retirement.
Source: the check-level-conformance-scoping retro and `capture/IDEAS.md`, drained
here. The handoff's `## Next command` body is then graded as
exactly one backticked harness-native invocation, so the authoritative next
action cannot drift into explanatory prose. Kit edit under the
`craft-synthesis` discipline. Sources: the craft-tickets, light-path,
artifact-suite, and artifact-hoist retros, drained here; `capture/IDEAS.md`, drained
here and in prior runs; `capture/learnings.md`, verdicted here.

The check-level-conformance-scoping close adds the terminal-record face. A
retired spec made the removed lifecycle's `status --full` projection unreadable
even though the durable run record survived, and that record omits the measured
promotion-stage timings the required retro needs. The terminal projection reads
the retained record after retirement and retains those timings. Until it does,
`/bench-final-check` captures the retro before retirement rather than rerunning a
successful promotion to manufacture evidence. Source: the
check-level-conformance-scoping retro, drained here.

The FT194 close adds the live prepared-state face: after review acceptance and
again during repair assignment, the removed lifecycle's `status --full` kept `next` at
`resume promote`. The same terminal projection derives the next action from the
current lifecycle state rather than retaining a stale prepared operation. Its
other CLI recommendations — retained promotion identity and evidence after
retirement, plus promotion-stage timings — are already carried above. Source:
the `ft194-project-green-desync` retro, drained here.

The exact-prospective-landing close adds the remaining lifecycle-observability
faces to this row: long promotions expose phase heartbeats and retained stage
timings, a recomposition conflict names its checkpoint owner and offers a
bounded already-covered or compatible-same-path route, and the operation
journal exposes remaining capacity before its bound blocks a repair. Its final
three-axis-review recommendation is already canonical in the spec-build
lifecycle, so it needs no second rule. Source: the
`exact-prospective-landing` retro, drained here.

The covers-traceability and FT195 closes add four more observable states to the
same owner: a blocking `bench gate --fresh` mode for callers that must sequence a
follow-up command, retained failing-phase and first-failure evidence after a red
promotion, an active-first lifecycle-status projection, and a `promote --check`
preflight for cleanliness and recomposition. The generic-worktree destination
and verdict-locality recommendations belong to FT169's landing primitive rather
than creating another subject owner. Sources: the covers-traceability and
go-build-cache-footprint retros, drained here.

Promotion evidence also retains per-phase elapsed values and peak-memory
attribution bound to the exact tree and worktree. Buffered output and a
coordinator-observed wall time are not durable enough for performance diagnosis.
Source: the `ft187-communication-surface-cut` retro, drained here.

The canary planted-reason close adds two receipt projections to the same owner:
a prospective gate failure says whether its failing package lies outside the
candidate diff, and `bench commit` renders retained per-phase timings in its
final receipt instead of requiring a JSONL-log read. Both remain bound to the
exact candidate, so neither weakens the gate or turns unrelated failures into
success. Source: the canary-planted-reason-ownership retro, drained here.

The handoff half now includes a fresh-context escape hatch, but only at a
durable lifecycle boundary. After repeated rejected repairs, a false-green
done-claim, or enough superseded evidence that the durable assignment is simpler
than the active conversation, the close recommends a fresh session and emits a
copy-paste continuation pinning the candidate, assignment, accepted and rejected
evidence, dirty paths, closed decisions, preservation constraints, and exact next
command. It never switches while a command or unrecorded mutation is in flight.
Source: `capture/IDEAS.md` 2026-08-06, drained here.

The gate-decision close adds cause-specific exact-green recovery. Pending,
stale, and mismatched records route first through plain `bench gate`, which
self-selects the necessary scope; `--fresh` is reserved for a reusable green
that exists but cannot compose, and remains reviewer-approved. The same retained
state already owns promotion timings, so the conformance-harness close's request
for exact stage timings is another instance rather than a second surface.
Sources: `capture/learnings.md` 2026-08-07, verdicted here; the
gate-decision-test-seam and conformance-harness-scope retros, drained here.

The gate-test-concurrency close repeats the terminal-evidence gap: promotion
retained its commit and green digest, but terminal status exposed no per-stage
timings, so final-check could report no trustworthy promotion timing. This is
another instance of the retained terminal projection above, not a second CLI
owner. Source: the `gate-test-concurrency` retro, drained here.

The branch-native rebuild makes the terminal-summary requirement explicit for a
lifecycle-format reset: retain a schema-migratable terminal record outside the
replaceable gate implementation, and give final-check an exceptional route that
requires the retained prospective log, published commit, and stated limitation.
The repair-reslicing close repeats that retained evidence must remain available
after retirement. Source: the branch-native-build-test-architecture and
repair-ticket-reslicing retros, drained here.

The coherent-diff repair adds one exact-subject face. After transient review
findings become a committed repair ticket and the pickup file is deleted in the
green repair, a second semantic-review preflight must bind the repaired product
candidate without treating planning provenance or the required pickup deletion
as an unauthorized product path. Decide whether the full run records an
explicit review base/head or preflight gains a repair-candidate selector; do not
weaken the ownership fence to make the composed diff pass. Source: the
axi-coherent-diff retro, drained here.

**FT142 (MEDIUM) — FT91 review residuals: eight open findings, two tracks.**
The ft91-gate-tier-split semantic review found twelve; three closed before
merge (the ship canary tier pin, the untiered-registry assertion, the
present-but-empty `CHECK` file). The spec and review files are retired, so
this row is the canonical carrier of what stays open — full citations
recoverable via `bench spec history ft91-gate-tier-split`. Standards track,
light-path candidates now: two comments that narrate the diff or cite spec
provenance (`internal/canary/fixture_tier_test.go`, the preprelease fixture
and `tier_test.go` failure message); the release-only package list derived
twice (`projects/benchkit.md` prose versus `releaseOnlyPackages` — name the
seam or anchor the string); and the `CHECK` filename hard-coded in the
same-package test beside the constant that owns it. Ship track, riding the next
`prep-release` hardening visit: the ship conformance step runs
with no `-timeout`, so the ~372 s probe plus all three release-only suites
share one default-bounded run — the 600 s hang hazard the spec claimed to fix
is restaged at ship; `govulncheck` resolves four levels deep instead of in
`requiredTools`, so a host without it burns the artifact matrix and a full
ship conformance run before dying — the concrete limit on the up-front
refusal; a second `prep-release` after an interrupt never cleans the orphaned
`dist/.preflight.*` staging directory; and two concurrent conformance runs on
one root interleave the timing file (needs a coverage row or an explicit
Won't-handle line). The ninth finding — the release-only `go test` step the
decisions promised, silently folded into ship-tier `goCoreTestPackages` —
closed 2026-07-27 with `ft91-gate-phase-split` story 18, which makes it an
explicit `prep-release` step. Source: the FT91 review, promoted at spec
retirement.

A second review's residuals join this row 2026-07-28, from
`ft91-canary-compiled-bites`: 17 non-blocking findings, full text recoverable
with `git show 4429b05:reviews/ft91-canary-compiled-bites.md` — `bench spec
retire` deletes a spec's review along with it, which is why the pointer is a
commit rather than a path. Standards track, light-path candidates: the contract
import prefix derived three ways (a one-source-per-fact violation, so the
cheapest of the batch to justify), and comments that narrate the change rather
than the code. Coverage track, each a decision rather than a defect: the bite
invocations inherit `BENCH_SKIP_LOG` from the sweep, the compile call's
environment is undecided, the bites carry no test timeout, and `TMPDIR` is left
unstripped. Two findings left this list at the drain: the profile's stale
nested-fixture description was corrected in the drain diff, and the
vacuity-baseline finding became FT153. Source: `capture/IDEAS.md`, drained here.

**FT144 (MEDIUM) — kit specs have two audiences and the discipline names
neither.** The `ft91-canary-check-scoping` build discovered mid-flight that
story 4's pinned seam was correct for the kit's own tree and a shipped
regression for every linked repo, because `bench init` scaffolds a seed canary
family a kit-owned table can never bind. The spec's edge inventory had walked
its hostile inputs without ever asking which repo was being swept. Two kit
edits, both built later under `craft-synthesis`. First, a `craft-spec`
edge-inventory prompt for kit code with two audiences (the Bench tree versus a
linked repo), since a fail posture that is right for one is a regression for
the other. Second, a reviewer decision on the workflow: the build preserved
the story's stated intent while moving the seam it pinned, and shipped that
under batch approval rather than routing back to `/bench-write-spec` with the
finding quoted. Either the existing route is right and the build should have
paid the round-trip, or the workflow wants a named lighter case for
"intent stands, seam moves" that a build may take under batch approval and
flag for veto — the reviewer's call, and the decision the row exists to get.

Both halves gained a second instance 2026-07-27, from `ft91-gate-phase-split`.
The audience half was applied by hand: that spec's edge inventory walked the
kit-versus-linked-repo split explicitly for all three of its new fail postures,
ahead of the `craft-spec` edit that would make it standing — and the semantic
review still found one probe testing for a directory *name* that any linked repo
could carry, so the prompt is worth having even where a spec remembers to ask.
The workflow half recurred in a sharper form: stories 4, 5, and 9 shipped as
probed phases rather than the manifest the spec named, and story 9 was dropped
outright as unsatisfiable — the seam moved *and* a story died, discovered
mid-build, shipped under batch approval and flagged for veto rather than routed
back. That is a larger deviation than the case this row was opened on, which
makes the named-lighter-path question the one to answer first.

A third instance arrived 2026-07-27 from `ft91-artifact-build-tiering`, and it
moves the question from the build phase to the review phase. That build's
semantic review found four statements in the *approved* spec that the code
falsified — a coverage row describing a seam that parks above the backups it
claimed to exercise, a story over-stating a deletion, and two edge-inventory
exclusions resting on invalidated reasons — and corrected all four in place,
each marked `**Post-approval correction, flagged:**` for veto. The code was
right in every case; only the prose moved. So the workflow question generalizes:
when a phase's finding lands on an approved spec rather than on the code, what
may it do under batch approval? `.bench/BENCH.md` leaves each spec in `specs/`
as post-hoc veto surface, which is what the review relied on, but a
reviewer-approved spec is a stronger artifact than a staged one, and "correct a
false citation" and "change what gets built" are not obviously the same
permission. Two candidate rules for `/bench-review-implementation`, the reviewer
picks: (a) a factual correction — a citation resolving to nothing, or a
described mechanism the tree contradicts — may be made post-approval under the
existing in-line veto-flag convention, while anything changing what gets built
stops for sign-off; or (b) all post-approval spec edits stop and the review
persists them to `reviews/<spec-slug>.md` instead. This matters standing-wise
because the Spec axis is charged to audit every coverage row, so it will keep
producing this class of finding. Answer it with the named-lighter-path question
above — one decision, both phases — then build the kit edit under
`craft-synthesis`. Sources: `capture/learnings.md`, drained in a prior run and
here; `capture/session-handoff.md`, drained in a prior run.

**FT158 (MEDIUM) — make cross-harness
falsification standing for kit-guidance diffs.** FT91's draft, FT152's build,
Occurrences: baseline-01, baseline-02, baseline-03
and the FT123 + FT124 build each received a Codex pass charged to refute rather
than grade after the ordinary review surface had cleared them, and each pass
found a real defect. The third run meets this row's graduation trigger: its
counterexample was an exact worktree label that looked like help, a flag, or
`--`. Make the pass standing for kit-guidance diffs, where a defect compounds
through every session that loads the prose, and give it the prepared spec and
diff bundle so it does not spend the charge re-reading unrelated contracts.
Every finding still takes an explicit accept, merge, or dismiss disposition;
the pass is advisory and does not become a second oracle. Kit edit under the
`craft-synthesis` discipline. Sources: `capture/learnings.md`, verdicted in a
prior drain; the gate-fastpath and FT123 + FT124 retros, drained here and in a
prior run.

The check-level-conformance-scoping close confirms that this review must stay at
the public CLI seam: the cross-harness pass caught disagreement among status,
dashboard, and handoff consumers that package-only tests did not expose. Source:
the check-level-conformance-scoping retro, drained here.

The ft156-anchor-registry build fixes the pass's placement inside a spec-build
lifecycle: run the falsification pass *before* final review and turn each
finding into its own repair ticket with an independent red mutation — that run
did exactly this and the pass found three defects the first review receipt had
missed, while the subsequent final review came back clean. The standing rule
should carry that ordering. Source: the `ft156-anchor-registry` retro, drained
here.

The conformance-harness close adds the evidence-bundle face: give a fresh
reviewer the coordinator-owned mutation receipts at review start, and distinguish
a permanent regression test from a temporary overlay that proves a ticket's red.
Withholding the receipt made the same valid mutation question get rediscovered
and resolved twice. Source: the conformance-harness-scope retro, drained here.

The spec-ticket handoff close makes a clean three-axis review the terminal
semantic boundary: a deterministic promotion red may create one evidence-backed
repair, but it does not restart open-ended review. Source: the
spec-ticket-handoff-contract retro, drained here.

The Pocock-guidance-doctrine close adds premise verification to the disposition
step. Its cross-harness refutation correctly found a risk but overstated severity
by treating a spec-time fence as the later ticket-time fence; the coordinator
checks the finding's subject and ownership premises before accepting its label.
Source: the Pocock-guidance-doctrine retro, drained here.

**FT98 (MEDIUM) — one preserve-then-discard primitive;
four faces.** Three rows were faces of one missing primitive — a sanctioned,
Occurrences: 2026-07-30-scoped-roadmap-commit, baseline-01, baseline-02
recoverable discard — and collapse to one semantics rather than three:
recovery ref written first, exact fingerprint required to apply, refusal if
the content moved, modelled on `bench worktree clean`'s existing contract, so
the operation earns its authority by being recoverable rather than by an
exemption to a guard. Face one, recovery payloads (this row's original
charge): cleanup fail-closes permanently when a payload's content landed
through different commits — the FT83 delegate payloads are strict subsets of
the default branch by diff, yet the since-removed recovery verb's
`--apply <fingerprint>` still returned `retain` because landed-proof requires the
payload commit itself (observed 2026-07-20); recurred 2026-07-22 when
`git cherry` missed reshaped commits and the reviewer had to hand-delete refs
and intent entries, the exact manual surgery the lifecycle exists to prevent.
Both routes shipped 2026-08-01: `LandedInDefault` proves a squash-landing by
reverse-applying the branch's cumulative diff against the default tree
(`efb456c`), and `bench worktree clean --discard-branch` is the
reviewer-supplied proof for what no derivation can establish (`37411a0`) —
fail-closed stayed the default and every ambiguity still resolves to
not-landed. The unprovable half shipped 2026-08-04 (`fafb049`, from the
`recovery-discard` spec retired here): the recovery verb's
`--discard <fingerprint>` retired one inspected payload per invocation without
asserting it landed, the plan separated an orphaned ref from an absent one and
reported how many paths the payload touched so the operator was not choosing
blind, and the reclaim verb deleted the provisional residue of
terminal spec-build runs. The lifecycle removal has since overtaken the drain
half of this face: the recovery and reclaim verbs are deleted, `bench resume`
no longer authors preservation refs, and its reconcile deletes
`refs/bench/recovery/` and the lifecycle assignment rows wholesale at every
session start. Sources:
`capture/IDEAS.md` and `capture/learnings.md` 2026-08-03, drained and verdicted
here — both reported release growing the preserved set with no retire route,
which is the route that has now shipped. Two residuals of that build were left
for this verdict and stay open. The retire side of the stale-fingerprint guard
has no unit coverage: mutating `applyRecoveryVerb`'s check to fire only for
discard still passes `go test ./internal/worktree`, and only a runtime contract
test kills it, so the gate bites and the hole is unit-level parity. And the
orphan deletion's compare-and-swap resolves its expected OID at delete time, so
a row-less ref is checked against a value just read rather than against the one
the plan graded; closing that means carrying the planned OID through the plan,
which is a design change and not a repair. The third occurrence came when a scoped roadmap commit was
blocked by an unrelated dirty session handoff on 2026-07-30; the session used
an isolated verification worktree, the landing workaround owned by FT169,
because the sanctioned set-aside primitive still does not exist. Face
two, `bench commit`'s set-aside (was FT127): the refusal reads "working-tree
files outside the named set block the commit — name them, or set them aside",
but no set-aside route exists in the CLI, so an agent's only real exits are
committing an unrelated file into a scoped commit or reaching for
`block-dangerous-git`-blocked plain git; build the route on this same
primitive rather than rewording the advice — the need is real and recurring.
The 2026-08-01 migration landing adds this face's complement: there is also no
sanctioned way to name "everything currently changed, deliberately" — `bench
commit -m … .` does not expand to the changed set, and the obvious
`git diff --name-only HEAD` silently hides a rename's deletion side until
`--no-renames` is passed, so a reviewer-approved whole-tree change took three
refused attempts and a hand-assembled pipeline. A `bench commit --all`
explicit opt-in naming the full working set as the attribution target keeps
the default path-scoped and the refusal intact; it is the same missing
"name this working set" vocabulary as the set-aside, so build them on one
contract. Face three, mutation-probe revert (was FT114): deliberately weakening an
implementation to prove a check bites always needs a revert, and
`block-dangerous-git` blocks `git checkout <path>`; copy-aside works but is a
papercut on a first-class activity in this repo (cf. `tests/canary/`), and a
scoped single-path revert through the same recoverable primitive replaces
both the papercut and any guard exemption. The one discard semantics is now
defined rather than pending — plan first, exact fingerprint, one target per
invocation, terminal outcomes that distinguish a proof from an operator's
judgment — and the remaining faces reuse it rather than restating it.

The FT131 close adds the ignored-cache face. Current-source verification defaults
`GOCACHE` under `dist/`, while nested route tests strip ambient cache variables;
an isolated full gate can therefore leave more than 40 MiB and 1,000 ignored
entries that `bench worktree clean --discard-ignored --full` refuses even with
the matching fingerprint. Prefer a lifecycle-owned scratch cache or cleanup at
the fixture owner; if residue remains legitimate, let the same size-bounded,
fingerprinted discard contract summarize and authorize the generated directory
without falling back to manual deletion. Source: the FT131 implementation retro,
drained here.

The FT176 close adds two `bench worktree clean` paper-cuts to the same
contract, and recurs the ignored-cache face as written. The plan output has no
stable way to extract the fingerprint — scripting it means awk over a TOON row
whose field position shifts on multi-row plans — so the discard contract wants
a `--fingerprint-only` plan mode or a keyed line. And the destructive limit
refused `clean --discard-ignored --full` outright when
`dist/.freshness-go-cache` exceeded it, leaving a manual `rm -rf` of the cache
as the only route; the derived-cache case deserves a carve-out or a named
override under the same size-bounded contract. Source: the
spec-build-lifecycle-preconditions retro, drained here.

The 2026-08-03 journal supplies another set-aside instance: landing FT183 while
capture files were foreign to the scoped commit required copying them out,
restoring their committed bytes, committing, and copying them back because both
plain revert and path-scoped stash were guard-blocked. That workaround is
recoverable only by session care, so it strengthens the existing CLI-owned
set-aside face rather than earning a guard exemption. Source:
`capture/learnings.md` 2026-08-03, verdicted here.

Repair-ticket reslicing confirms the ignored-residue face: a clean worktree can
remain retained solely because generated gate/cache inventory exceeds the
destructive limit. A bounded cleanup mode may cover known generated caches, but
must keep the existing fingerprinted, size-bounded refusal for everything else.
Source: the repair-ticket-reslicing retro, drained here.

The lifecycle-removal close exposes the destructive edge of the remaining
cleanup route: `bench worktree clean --apply` can still preserve dirty work into
`refs/bench/recovery/`, while the next resume sweep now deletes that namespace.
Explicit clean therefore refuses a dirty removal unless the same recoverable,
reviewer-visible discard contract owns the payload through completion. Source:
the remove-spec-build-lifecycle retro, drained here.

**FT169 (MEDIUM) — one sanctioned worktree landing command
owns the stale-base dance.** The gate-fastpath build hand-ran the same sequence
Occurrences: baseline-01
for eleven ticket landings: fast-forward the assignment worktree onto the
current default branch immediately before landing, create the gated commit,
fast-forward the result back, then release the worktree. One stale-base miss
diverged and required a manual no-commit cherry-pick plus a fresh gate before
the sequence could continue. A `bench worktree land` command should make the
subject transition atomic from the session's point of view and fail closed
with recovery state intact when the branches cannot fast-forward. The exact
authority, interruption recovery, and relationship to `bench commit` and
`bench worktree release` start as a reviewer decision rather than being
inferred from the hand-run sequence.

The FT123 + FT124 close supplied the foreign-dirty face: `bench commit`
correctly refused while unrelated capture edits were present, but landing the
reviewed spec and retirement set then required two isolated worktrees and exact
file transfer. The command owns preparing that path-scoped landing worktree and
the safe transfer sequence; the full-run guidance names the same pattern rather
than leaving sessions to reconstruct it. Entry: `/bench-shape-idea`. Sources:
the gate-fastpath and FT123 + FT124 retros, drained here and in a prior run.

The FT129 build supplied the delegation-entry face. A write delegate's
assignment was behind `main`, but its Codex sandbox could not update shared Git
metadata to perform the required fast-forward. The landing command should offer
a delegation-ready preflight that fast-forwards the assignment from the
coordinator side, verifies and reports `HEAD == main`, then hands the exact
subject to the delegate. The accompanying `craft-delegate` guidance names this
Codex constraint and routes a denied delegate-side fast-forward back to the
coordinator instead of spending retries on permissions. Source: the FT129
implementation retro, drained here.

The FT126 close adds tool readiness to that same coordinator preflight. Both the
projection and maintenance delegates reached their first gated action with an
untrusted derived binary and paid the refusal before rebuilding, even though
the refusal correctly named the exact local rebuild. Before an assignment
subject is handed to a write delegate, the preflight verifies that subject's
derived binary and reports the rebuild action; an explicit opt-in may perform
it, while gate entry remains repair-silent and never rebuilds automatically.
Source: the FT126 recurrence-tallying retro, drained here.

The FT128 close supplied the base-ref face, and the since-removed spec-build
family covered only the spec-backed half of it. `bench worktree create`
always roots at the default branch with no base-ref flag, and
`bench worktree release` refuses while an assignment branch has not landed
there, so a chain of tickets on an integration branch cannot cut a worktree at
the chain tip and cannot retire any worktree until the reviewer merges; that
build formed the chain by merging each previous assignment branch into the next
worktree, which is the hand-run form of a compare-and-swap integrate. Reviewed
spec-backed builds briefly got that surface from the removed lifecycle's assign /
checkpoint / integrate operations; with the lifecycle gone no route has it, and
the landing command owns whether any should. Source: the FT128 implementation retro,
drained here.

That face recurred on 2026-08-01 during per-component-gate-scoping and sharpened
into a precise gap: the core already carries the parameter, and only the CLI
withholds it. `worktree.Create` takes a variadic start ref, but the sole caller
that supplied one was the removed lifecycle's assign operation; the generic
`worktree.CreateCommand` derives a start ref only from `--refresh` and otherwise
passes none, so creation roots at the resolved default branch (falling back to
`HEAD` when none resolves) with no way to name a base. A spec-branch build
therefore cannot cut delegate worktrees on its own branch, and the delegate opens
on a base that cannot fast-forward; that session worked around it with plain
`git worktree add`, which leaves the assignment unowned. Deciding the landing
command should settle whether the fix is a base-ref flag on `bench worktree
create` or whether every non-spec chain is meant to route through the spec-build
surface. Source: `capture/IDEAS.md` 2026-08-01, drained here.

The foreign-dirty face recurred on 2026-07-31 in the shape the existing advice
does not reach. A concurrent session writing the same checkout tripped
`bench commit`'s whole-tree attribution refusal twice during a reviewer-approved
drain — once on untracked `decisions/` files from an in-flight
`/bench-shape-idea`, once on a mid-drain `bench idea` write to `capture/IDEAS.md`. The
refusal is correct and cheap, and it fires before the gate; what is missing is
the sequence after it. Invariant 1 says take side-work to a worktree, but a
drain is not side-work and its diff already lives in the main checkout, so the
landing command owns either one sanctioned sequence for landing a batch beside a
concurrent writer, or an explicit statement that the main-checkout batch assumes
a single writer plus the way to detect that it does not. FT168's oracle-scope
question is a neighbour, not this: the blocker here is attribution, not gate
scope. Source: `capture/learnings.md`, verdicted here.

The FT176 close adds a harness-side face to the same delegation preflight: the
Claude WorktreeCreate hook admits one active hook-created assignment per
session, so parallel `isolation: worktree` delegate launches fail after the
first, and the working route is a manual `bench worktree create --request
<uuid>` per delegate. Either the limit lifts or the manual route becomes the
documented canonical one. Source: the spec-build-lifecycle-preconditions
retro, drained here.

The go-build-cache-footprint work shipped and retired at `f140e94`; its cache
reduction is no longer roadmap work. What remains from that close belongs here:
generic-worktree landing must name the destination branch and preserve a usable
verdict across the worktree-to-main transition. Source: the covers-traceability
retro, drained here.

The Pocock-guidance-doctrine close adds a landing-location precondition. Twice,
the coordinator ran the final fast-forward merge from inside the ticket
worktree, merged the branch into itself, and accepted Git's "Already up to
date" until worktree release refused. The sanctioned landing command resolves
and reports the primary checkout and destination before merging instead of
trusting ambient cwd. Sources: the Pocock-guidance-doctrine retro and
`capture/learnings.md` 2026-08-12, drained here.

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
traps automation and leaks on signals.** Reviewer ruling 2026-08-01: humans
should not be driving the `bench` CLI in the vast majority of cases, and
worktrees are not an exception — the agent-facing surface is the subcommands,
and agent worktree creation already flows through `bench worktree create`
and the shift loop. The bare verb instead runs `Subshell`
(`internal/worktree/worktree.go:412`): it creates a worktree, runs `$SHELL`
with stdin inherited, and blocks on `cmd.Run()`, so an automation call
probing for usage hangs until signal-killed (observed: three invocations,
two leaked worktrees, both retired with `bench worktree clean`). The release
is not signal-safe either — it runs on the line after `cmd.Run()` returns,
so SIGTERM or SIGINT skips it and leaks a registered worktree. Make the bare
verb print usage like every other parser-first verb and move the subshell
behind an explicit opt-in name if it is kept at all; if kept, the release
needs a signal trap or a lease the resume path reclaims. The
discovery convention is now canonical in `AGENTS.md`'s shell rules. Sources:
`capture/IDEAS.md` (reviewer ruling), drained here; `capture/learnings.md`,
verdicted here.

**FT172 (MEDIUM) — the roadmap parser and context snapshot make the drain's
non-recurrence evidence complete.** The row grammar is currently implicit:
`ParseDocument` treats any line opening with bold as a new row, so one bold-led
body paragraph made the snapshot silently partial. Decide whether a row is
exactly the `**FT<n> (…) — title.**` form with bold forbidden at body starts,
or whether the parser absorbs bold-led body paragraphs; then document and
contract-test that grammar.

The context half fills the facts the phase otherwise re-derives.
`/bench-what-next` step 1 orders every row
verified against the tree, and the snapshot is declared its complete local
evidence — but two things the reconcile depends on are not in it, so each run
re-derives them by hand (observed 2026-07-25, and again in the second drain
that same day, which opened with the same hand-run `git log`; the 2026-07-26
restructure drain repeated it a third time). First, the
workload boundary: with 40 rows, "verify every row" is tractable only by
knowing what moved, which today means a hand-run `git log` to establish that
nothing but docs and specs landed since the last drain. Emit it —
last-drain commit, commits since, and which of those touched code — so the
reconcile bounds itself and the "no row shipped" verdict rests on evidence
rather than on the session's recall. Second, the snapshot reports
`specs[]{slug,status,roadmap_id}` as raw fields, so a staged spec with an empty
`roadmap_id` reads as data; the roadmap's own preamble makes that a
discrepancy, and it went uncaught until a session happened to know the
convention. Add a discrepancies block covering at least that case, a row naming
a spec path with no file, and a row naming a retired spec. The CLI does no
judgment — these are mechanical cross-checks of facts it already holds, which
is why they belong in it rather than in the phase prose. That case recurred on
2026-08-03: two staged specs carried an empty `roadmap_id`, `occurrence_discrepancies`
was empty, and the drain found both by hand.

The same drain found a third face, in the phase prose rather than the CLI.
`/bench-what-next` orders `bench roadmap --context` invoked exactly once and
declares that snapshot its complete evidence, but the default truncates long
bodies — the run's longest journal entry and both retro bodies came back cut,
which is precisely the material steps 4 and 5 must read whole. Either the phase
names `--context --full`, or the snapshot stops truncating the sources whose
completeness it asserts. Source: `capture/learnings.md` 2026-08-03, verdicted
here in the run that hit it.

That parked half now has a demonstrated cost, and a cheaper option than the
probe it was parked behind. The 2026-07-27 drain promoted an `capture/IDEAS.md` capture
into a HIGH row and carried the capture's *diagnosis* forward as fact, including
the half claiming the artifact contract tests resolve their output directory to
the graded root. The very next phase read the tree and found it false — both
call sites already used `t.TempDir()` — so half a HIGH row was fiction, and it
cost that phase a false start into `/bench-write-spec`. The gap bit inside one
phase transition rather than eventually, which is the evidence the parked half
lacked. The cheap partial worth weighing against the general prober: a
`/bench-what-next` rule that a capture's symptom is promoted as evidence while
its diagnosis is either checked against the tree or written as the capture's
claim rather than as fact. That rule is prose, needs no new mechanism, and was
applied by hand in the 2026-07-27 drain — every diagnosis drained there was read
out of the tree first. Source of this clause: `capture/learnings.md`, verdicted
here. The grammar face came from `capture/IDEAS.md`, drained here.

The discrepancies clause has its second instance, 2026-08-01, and the tree shows
the pairing is broken from both ends. The roadmap end: the FT135 row omitted its
staged spec path the preamble requires, and was corrected here only because the
session happened to know the convention. The
spec end: `roadmap_id` comes from a `Roadmap:` header inside the spec file
(`internal/spec/spec.go`), and none of the three live specs carries one — the
header was written by hand on some retired specs (`craft-tickets`,
`ft126-recurrence-tallying`) while neither `/bench-write-spec` nor `craft-spec`
mentions it, so the field the snapshot reports is one no authoring surface fills.
Whether the discrepancy check reads the roadmap path, the spec header, or both
decides whether that header is taught or dropped.

A third face joins the same owner, and it is the reconcile's own hidden
dependency. `checkOccurrenceLedgerMigration`
(`internal/conformance/docs_workflow_checks_test.go`) carries a `want` map
pinning occurrence counts for seven named rows (FT71, FT158, FT98, FT169, FT141,
FT94, FT125), so the gate goes red both when the reconcile removes one of those
rows and when a drain records a new occurrence key on one — and no phase
instruction mentions the map, which cost a reviewer-approved drain a full gate
run on 2026-07-31 when FT128 was reconciled out. The map's own bite test asserts
only the legacy heading, so the counts are a second derivation of a fact
`ROADMAP.md` already owns. Prefer making the check derive from the ledger it
grades over teaching every drain to hand-edit the map; if the pinned counts are
load-bearing migration evidence rather than a duplicate, then the phase names the
map and the snapshot surfaces the pin as a discrepancy when a graded row is about
to move. Source: `capture/learnings.md`, verdicted here; the map and its bite test
were read in the tree on 2026-08-01.

**FT198 (MEDIUM) — make `ROADMAP.md` a progressively loaded index.** The board
has reached 62 rows and 179 KB; this drain's required full context snapshot was
complete on disk but too large for one agent-tool response. Keep the main
roadmap as a concise index of ID, title, priority, state, dependencies, and the
next phase, with each row's detailed evidence and rationale stored behind one
canonical on-demand reader. `bench roadmap --context --full` remains one
consistent snapshot, but the ordinary query and phase can load the index first
and request only the detailed records needed for reconcile or drain judgments.
Decide the durable detail owner, migration and history behavior, and how the
parser proves index-to-detail completeness without creating a second source of
status. Entry: `/bench-shape-idea`. Source: `capture/IDEAS.md`, drained here;
the 2026-08-06 drain's 179397-byte snapshot transport failure.

**FT173 (MEDIUM, decision required) — AXI residual: the
active-assignment-with-deleted-tree disclosure class.** The AXI program is
otherwise landed: `bench diff` owns the coherent Git snapshot, and the
`axi-query-disclosure` capstone (implemented 2026-08-12) composed contextual
disclosure, the harness-log ledger, registry-derived conformance, and the
ten-principle `craft-cli` guidance across the approved query surfaces and
nested `worktree list`. The durable boundary holds: AXI stays scoped to the
high-frequency query surfaces — turn-preventing principles outrank format, and
a gate-checked query command carries roughly five contract behaviors, so
conformance verification exceeds the feature elsewhere; no parallel renderer,
legacy mode, or `bench git` namespace.

The residual is review finding R11, a state the spec never defined: an active
assignment whose worktree tree has been deleted is an undecided QD6 hybrid —
`internal/worktree/list.go` renders active/missing, advertises path/exec
actions, and never reports clean. Needs a reviewer-approved amendment defining
the disclosure class before any code moves; once the class is decided, the
change itself is light-path sized. Entry: reviewer decision on the class, then
the light path. Sources: the axi-query-disclosure retro and
`capture/IDEAS.md`, drained here.

**FT202 (MEDIUM, decision required) — a standing test-support fence, and the
census scope for process-backed fixtures.** Two coupled reviewer decisions from
the axi-query-disclosure close. First, cross-package test scaffolding has no
fenced home: two amendment cycles in one spec (the axitest relocation, then
`internal/gittest`) each paid a reviewer round-trip for the same question.
Propose a named commons fence in `projects/benchkit.md` — test-support only,
census-visible — instead of per-spec amendments. A/B evidence exists: round two
confirmed the vocabulary-gap hypothesis (commons prose flips stop-and-surface
into a gate-green no-round-trip build); the evidence and candidate B diff live
outside the tree at `~/fence-experiment/` (README inside, not re-read this
drain), and the census-visible sentence needs tightening before the prose
graduates. Second, the ordinary-test census does not count process-backed git
fixtures reached through `internal/git` (21 repos in the cmd/bench registry
envelope test); extending the census is an oracle change and a reviewer
decision. Order constraint: answer the census-scope question first or require
the scaffold stay census-visible — the pending R2 harness collapse (repair
ticket T1, awaiting the fence decision) would otherwise move those constructors
into a non-`_test.go` scaffold `architectureOwnedTest` never inspects, making
them permanently census-invisible. In-scope residue: `leasedRepo` stayed
outside the gittest collapse (it would need a third API shape), flagged in the
T1 ticket evidence. Entry: `/bench-shape-idea`. Sources: `capture/IDEAS.md`
(three lines) and the axi-query-disclosure retro, drained here.

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

**FT175 (MEDIUM, decision required) — a claim ledger for assertions about the
world.** The gate refuses "I believe the tests pass" and nothing refuses "I
believe this harness supports fresh-context delegation." A draft spec proposes
`bench cite` / `bench claim`: a CLI that acquires and stores evidence with its
hashes, an agent-written assessment over a four-state vocabulary, a span
verifier that confirms a quoted substring is present in the stored bytes, replay
that marks an assessment stale when its evidence moves, and a gate phase over
the whole thing. The draft was reviewed 2026-07-31 by cold re-derivation against
this tree; most of its repository claims reproduced, and the parts that did not
are recorded against the rows they touch rather than here.

It enters as a decision tree, not a spec, because three questions interact and
none is settled. The draft specifies its own gate phase two incompatible ways —
well-formedness only in one section, consuming assessment state in another — and
only the first is defensible under invariant 1, since the canary can prove
`bench claim check` rejects a malformed record but cannot prove an assessment is
honest. The CLI narrows the agent's assertion without removing it: fabrication is
blocked, but the agent still selects which evidence to acquire and which span to
quote, and spans are exempted for the two states most likely to hide a check
nobody ran. And the ledger would be the kit's first durable capture store whose
retirement rule cannot be "drain to zero" — every existing source empties, while
a claim's value is that it persists. Settle those three before any spec. The
span verifier is not prior art: `kunchenguid/no-mistakes` was read at its README,
gate-model, pipeline, auto-fix, and pipeline-steps pages plus its step
implementations, and documents no span or quotation verification — its finding
`action` vocabulary (`no-op` / `auto-fix` / `ask-user`) is agent-assigned, so the
transferable part is its fail-closed default, not the classification itself.
Entry: `/bench-shape-idea`. Source: `capture/IDEAS.md`, drained here.

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

The frontmatter half now has a demonstrated bite. `craft-line`'s description
carried an unquoted ` #`, so YAML read the rest as a comment and two of its
three trigger clauses never reached any harness — a skill silently half-armed,
fixed 2026-07-25. The gate greps frontmatter rather than parsing it, so
nothing saw it. A real parse closes this class; a narrower check rejecting an
unquoted ` #` in any skill frontmatter description is the cheap version if the
parse proves too large a step.

The stale-reference half has a second demonstrated instance, found 2026-07-31
by a cold re-derivation of the tree. `craft-review`'s frontmatter description
still lists "a self-review before commit" among the contexts that charge from
it, while `/bench-implement-spec` closes inline self-review outright and
rejects a context-inheriting delegate as the same failure under another name.
A skill description is the only part of a skill every harness loads, so a
stale one advertises a route the phase forbids. One-line fix; it belongs with
the frontmatter parse rather than alone.

The derive-inventories-from-one-implementation half gains a named enforcement
target: a gate check reconciling the `.bench/BENCH.md` CLI inventory (and the
`BENCH-reference.md` plumbing list) against `bin/bench.sh`'s case labels, so
the inventory's "kept in sync with `bin/bench.sh`" promise is enforced rather
than hand-maintained. Source: `capture/IDEAS.md` 2026-08-05, drained here.

The stale-reference half also owns the decision documents that still cite the
removed `internal/canary/canary.go` and `internal/specbuild/` owners and the
retired `$bench-finalize-spec` phase. Re-derive those Sources and phase routes
from the live tree rather than carrying lifecycle names forward. Source:
`capture/IDEAS.md` 2026-08-11, drained here.

The debugging half gains an upstream Redact candidate: shown output substitutes
`<REDACTED>`, reproduction loops address secrets through environment variables,
and excerpts retain only signal-carrying artifact lines. Fold it into the
coherence pass rather than adding a parallel security skill. Source:
`upstream(mattpocock/skills@84fdeff)`, drained from `capture/IDEAS.md` here.

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

**FT189 (MEDIUM) — an upstream `git worktree list` hang reaches every Bench
worktree read.** `git worktree list --porcelain` hangs on a FIFO gitdir placed
in any private admin entry — reproduced 2026-08-03 at exit 124 under
`timeout 5`. Every `bench` command that enumerates worktrees inherits the hang
before any Bench guard runs, so the failure is a wedged command rather than a
refusal, and no deadline the repo owns bounds it today. The work is a mitigation
Bench can own: a pre-scan that refuses a malformed admin entry by shape, a bound
on the enumeration call, or both, with the upstream behavior named so the
mitigation retires if Git fixes it. Source: `capture/IDEAS.md` 2026-08-03,
drained here.

**FT190 (MEDIUM) — every injected interface has a real-producer test or a
written exemption.** The `injected-interface-junctions` build's P1 lived at
exactly one such seam: a type takes an interface, the tests drive a fake, and
nothing composes the real producer, so the junction's contract is asserted by
the double that was written to satisfy it. Before its removal,
`internal/specbuild` alone took six
(`GateOwner`, `PromotionGateOwner`, `WorktreeOwner`, `ReleaseOwner`,
`AbandonOwner`, `Runner`). Enumerate the injected interfaces across the tree and
give each either a test that composes its real producer or a recorded exemption
saying why a fake is sufficient there. The candidate ending is a conformance
check over the enumeration, which would make a new injected interface arrive
with its junction obligation rather than without one; whether the audit ends in a
check or in a one-time sweep is the build's decision. Source: `capture/IDEAS.md`
2026-08-03, drained here.

**FT191 (MEDIUM) — a fixture-and-seam inventory a charge can carry for free.**
`craft-delegate` now requires a charge to carry its fixture inventory — helper
names, gate doubles, prior-art fixtures with `file:line` — and both 2026-08-03
retros credit that input for their first-pass delegate greens. The inventory is
assembled by hand every time. Extend `bench outline` (or a sibling reader) to
emit test helpers, doubles, and prior-art fixtures per package so the charge
input costs nothing to build. Scope decision at spec time: whether this is a new
projection or a mode of the existing outline owner. Sources: `capture/IDEAS.md`
2026-08-03 and the `ft181-precondition-residuals` retro, drained here.

**FT192 (MEDIUM) — one source per fact reaches spec and ticket prose.**
`AGENTS.md` forbids two derivations of one fact, and the standard is written for
code; specs and tickets restate implementation-derived facts freely. The
`injected-interface-junctions` review's SP4 is the demonstration: a failure-message
count lived in the spec's implementation decisions, in a ticket's `Contracts:`
line, and in the implementation's message constants, with only the implementation
as truth, and `bench coverage --check` reported a valid map on the exact candidate
where the count disagreed — it grades the map's shape, not its agreement with the
tree. A 7-line shell loop comparing the advertised count against the constants
reproduces it in 7ms. This drain found the same shape twice more in `ROADMAP.md`
itself (two rows restating how many `specs/` directories are ticket-only, both
stale), and reworded them. Extend the standard to spec and ticket prose — a spec
names the enumeration source instead of restating its count — and let
`/bench-review-implementation`'s Standards axis grade a restated count as
duplicated knowledge. Whether any of it becomes a gate check rather than a review
judgment is the reviewer's call at spec time. The `recovery-discard` build adds
the ticket-side instance: the retired `Assumptions:` field copied the standing
tree-verification rule into every ticket it governed. That field and its parser
machinery are gone; what stays here is the one-source-per-fact rule the instance
argued for. Kit edit under the `craft-synthesis`
discipline. Sources: `capture/learnings.md` 2026-08-03, verdicted here and in a
prior run.

**FT204 (LOW, decision required) — one bounded transcript/session query.**
Agents repeatedly shape harness transcripts and session evidence with
`head`/`tail`/`awk`/`sort` chains; consider one bounded agent-facing query for
transcript and session census plus evidence projection. This is a new
operational surface, not an AXI query-surface extension, so it needs its own
operational-surface decision before any shaping. Entry: reviewer decision.
Source: `capture/IDEAS.md`, drained here (routed out of the
axi-query-disclosure harness-log ledger).

**FT205 (LOW) — `craft-delegate` names the delegate-worktree end-of-life
pair.** A delegate assignment ends through `bench worktree release` by the
creating request; `bench worktree clean` alone leaves the assignment
`recovered` in the intent ledger, and the one-delegate-per-session pool then
refuses every later create until `bench resume-clean` reconciles it (observed
2026-08-12). Kit edit under the `craft-synthesis` discipline: one clause in
`craft-delegate` naming release as the default end and clean-then-
`bench resume-clean` as the recovery pair for harness-created delegate
worktrees. Source: `capture/learnings.md` 2026-08-12, verdicted here.

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

**FT99 (LOW) — spec problem-premise verification.** A spec compiled from a
closed decision map can inherit a problem statement the tree has since
falsified: the retired `minimal-subprocess-data-exposure` spec claimed the
project gate "gets everything except `BENCH_KIT` and `BENCH_WRAPPER`", which
FT78 had already
fixed, and the build reached stage 1b on that false premise before a contract
delegate caught it. Require every "today the code does X" claim in a spec's
Problem section to be checked against the tree at spec time, with the check
named in the spec — the same standard the coverage map already applies to its
red signals. The same gap has a second face in the Solution and Implementation
sections: the retired `cli-grammar-and-capability-evidence` spec asserted a
structured line "survives non-verbose `go test`", which is false, and stories 10
and 11 were built on it before an independent done-claim probe caught it. Extend the
rule to those claims too — any spec sentence asserting observable third-party
tool behavior a story's seam depends on carries either a cited command whose
output was run or an explicit uncertainty flag. `/bench-write-spec` step 9's
falsification pass now runs on every draft, but it is charged at the coverage
map and the Handoff assertables, so a Problem-section claim that reads as
obvious fact still slips through. A third instance, 2026-07-27: the
`ft91-phase-manifest-dag` spec asserted that append semantics hand a child two
values for one key, so its story 7 red signal could not occur — `os/exec` dedups
`cmd.Env` before exec and the mapped test passes identically against the old
code. Confirmed by standalone probe during the semantic review, after the build
had shipped. Next action is the kit edit to `/bench-write-spec` and
`craft-spec`, built under the `craft-synthesis` discipline.

The review half uses the same truth obligation. When a spec enumerates call
sites with per-site postures, `craft-review` makes that table an explicit
Spec-axis audit: the reviewer walks every row against the diff. FT86's
DefaultBranch table lacked that obligation, and two adopt call sites reversed
their assigned posture until semantic review caught them. Source: `capture/IDEAS.md`,
drained here.

The FT194 review adds the Coverage-axis counterpart: when behavior crosses
multiple publication sites and refusal classes, the review enumerates their
cross-product rather than treating one direct-path refusal and one recovery
success as evidence for every pairing. Source: the
`ft194-project-green-desync` retro, drained here.

The exact-prospective-landing close adds immutable-snapshot coverage: the first
review runs tracked-file, transitioned-spec, and new-untracked-descendant
mutations together, and the first public-adapter slice exercises valid
nested-CWD arguments alongside usage failures. Source: the
`exact-prospective-landing` retro, drained here.

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

The FT126 close supplies the ticket-level instance: four independently-green
tickets repeatedly carried the same 49 unrelated whole-tree structure findings.
A path-scoped structure summary for the staged ticket set keeps the aggregate
debt visible while making the ticket verdict immediate; it rides this same
per-context scoping surface rather than adding a second structure command.
Source: the FT126 recurrence-tallying retro, drained here.

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

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency
and dogfood loops.** A kit edit that instructs spending a model tier can contradict the
escalation policy without any loop catching it: the widened write-spec step-9
triggers shipped an automatic top-tier spawn past review (observed
2026-07-22; corrected in the mid-tier rerouting commit). Make
`craft-synthesis`'s consistency loop name the escalation policy as a standing
cross-check for any kit edit that spends a tier. Kit edit under the
`craft-synthesis` discipline.

Ticket-shaping guidance must also dogfood its own decomposition: apply
`craft-tickets` to the candidate, reproduce any claimed keep-together project
gate red, and compare the resulting slices, blocking edges, integration
surfaces, and fixture/anchor knowledge against the prior guidance. Reject or
reslice a change that widens unrelated tickets or creates a second source of
knowledge. Source: `capture/learnings.md` 2026-08-08, verdicted here.

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

**FT186 (LOW) — split the gate executor and type its verdict
records without changing behavior.** Two structural hotspots sit behind the
gate seam. `executeSubjectWithEngine` combines locking, subject-drift checks,
component selection, dispatch, and four persistence branches in one roughly
200-line function. Verdict field-set validity is separately hand-maintained in
`readyFieldClasses` and the record-class registry. Under the mechanical
refactor lane, split the executor by those responsibilities and make
record shape type-enforced, preserving the existing gate and mutation tests as
the exit test. These are refactors, not a gate rewrite. Source:
`capture/IDEAS.md`, drained here.

**FT111 (LOW) — provenance tags that outlive their specs.** Code comments
carrying `FT<n> story <n>` tags point at specs that retire — two retired on
2026-07-23 alone — so a tag naming a retired spec points at nothing, which is
exactly the rot `craft-comments` forbids. Reviewer-decided 2026-07-23: remove
them, but only when already editing the line, and reject new ones in review.
The work is the wording in `craft-review` and `craft-comments`, not a sweep of
the roughly a dozen existing sites. Kit edit under the `craft-synthesis`
discipline.

**FT112 (LOW) — an approximation that stays green is not a cleared bug.**
Diagnosing the gate's marker stalls (the retired
`trustworthy-gate-under-load` work), four synthetic load shapes — guest CPU
saturation, parallel contract loaders, inert memory ballast, and an `fsync`
hammer — all stayed green, and only a real `bench gate` under host-side load
reproduced the failure; each disproven shape cost a build-and-measure cycle.
`/bench-debug` phase 1 tells a session to build a tight repro loop but says
nothing about what a *green* approximation proves, so a fresh session can keep
reaching for cheaper stand-ins. Add the reproduction-economics rule to the
phase's loop-building step: a proxy that stays green narrows a hypothesis at
best and never clears the real oracle, so when the failure is load- or
environment-sensitive, reproduce through the accused command under the
conditions that exposed it before spending another cycle on a stand-in. Kit
edit under the `craft-synthesis` discipline. Source: the 2026-07-23 learnings
entry, verdicted in this drain.

**FT113 (LOW) — `bench commit --spec` residuals: the flip counts as a path, and
the flip has one author.** The row's original face — the `Status: staged` →
`implemented` flip leaving the gate verdict strong-stale because `specs/*.md`
was off the capture-only allowlist — is resolved by the shipped reduced-gate
scope: `specs/` is on the allowlist and the content-addressed ancestor lets
allowlist-confined post-flip drift inherit its evidence, so the follow-up run
is reduced rather than a full oracle re-run, and the three-way allowlist
decision the row queued is settled by that shipped shape. What remains is the
command's usage contract: `bench commit --spec <slug>` edits the owned spec
transition itself but still requires another named path — count that owned
transition as satisfying the path requirement without widening the staged
set. Sources: the FT131 and decision-map integrity implementation retros,
drained here.

The transition now has three authors, not two. `bench spec implemented <spec>`
and `bench commit --spec <slug>` both perform the `Status: staged` →
`implemented` flip and each expects to own it, so running the former first makes
the latter fail with `no Status: staged line`; the removed lifecycle's promote
verb was a third, and the phase contract now names `bench commit --spec` the
sole author for a reviewed spec-backed build. Name one owner per landing route
and make the others refuse rather than race. Source: the FT128 implementation retro, drained here.

The coherent-diff close exposes the same usage residual directly: the canonical
final-landing prose abbreviates the command to `bench commit --spec <slug>`, but
the executable requires both `-m` and an explicit path. Let the owned status
transition satisfy the path requirement, and keep one copy-pasteable command
shape in the phase contract so a conformant close does not pay two usage
refusals. Source: the axi-coherent-diff retro, drained here.

**FT130 (MEDIUM) — a capture write mid-lifecycle voids or blocks the run.** During
FT122's gated commit a session answered a reviewer question and ran `bench
idea` to park the tangent, which wrote `capture/IDEAS.md` inside the gate's window;
every phase came back green and the commit was still refused with "gate subject
changed during execution", costing a full ~15-minute re-run. Two standing rules
collide: `projects/benchkit.md` forbids mutating the repository while a gate
runs, and `.bench/BENCH.md`'s Capture section says to park a tangent the moment
it appears. Name the collision where the capture rule lives — parking defers
while a gate or gated commit is in flight — but prefer the mechanical fix,
which makes the prose unnecessary: `bench idea` can see the subject lock, so it
can queue the line and write it when the run finishes, or refuse with that
reason rather than silently voiding the verdict. Deciding between queue and
refuse is the row's real work — so this row stays mechanical rather than adding
another standing prose rule. Kit edit under the `craft-synthesis` discipline. Source: the
2026-07-25 learnings entry, verdicted in a prior drain.

The gate-window face recurred on 2026-08-02 — a `bench idea` write landed inside
a `bench gate --fresh` window, every phase came back green, and the verdict
rejected itself for a changed subject — which is the second full re-run this row
has now cost. The removed lifecycle moots the former clean-checkout,
recomposition, checkpoint, and review-to-promote faces; they add no surviving
work to this gate-lock owner.

**FT138 (LOW) — instrument Bench so build economics are measurable.**
Reviewer-priced 2026-07-25 as a nice-to-have in its current state, so it
holds LOW until an acceptance trigger that needs measurement (the cheap-tier
re-test still open as `decisions/cost-follows-project-size.md` #6, and the
mid-tier default) actually blocks on it. Candidate
metrics: delegate tokens per slice, coordinator tokens, review findings per
implemented slice split by axis (Standards/Spec/Coverage), rework tokens
spent after a build already went green, gate runs per spec, iterations
against the declared cap, tier declared versus tier actually used, and phase
plus package timings in the normal gate summary so critical-path diagnosis
does not require a separate instrumented run. The
only baseline that exists (FT86 review resolution: 26 findings, roughly 350k
delegate tokens of rework) survives by accident in one session transcript
rather than in any artifact. Open questions — which metrics earn their
capture cost; where they live (a `bench` subcommand, the gate cache, or the
journal); who writes them, given the harness-independent substrate is the
shift loop and the git hooks rather than any harness; per-slice versus
per-spec granularity; retention and pruning; agent-facing AXI or
reviewer-facing — mean the work starts as a grill. Entry:
`/bench-shape-idea`. Source: `capture/IDEAS.md`, drained here.

**FT164 (MEDIUM) — repair-lane charges, and a done-claim that resolves its named
owners.** The ticket-contract core shipped 2026-08-03, promoted at `83c630e`
and since retired; what stays open is the repair-and-experiment lane that spec
cut, plus two rules the build itself produced. The lane is one `craft-delegate` visit. Model-comparison charges give
every candidate one constant charge — base commit, file fence, effort, focused
suite, independent probe — with setup failures recorded separately from
generation time; fixed behavioral checks and the independent probe decide whether
the cheap default clears its bar, while style differences decide which acceptable
patch is better rather than justifying an expensive tier by themselves
(`craft-line`'s cheap default and red-driven ladder stay authoritative). A repair
that updates a canonical command inventory names both currency owners in its
focused proof — the documentation inventory and the wrapper/router registry —
because FT131's repair updated one and missed the other two until the atomic
gate. And a repair touching a shared environment helper names which rows opt into
shared caches and which stay hermetic, then runs the focused failing rows before
the full gate; the artifact-suite repair over-stripped that distinction.

Two rules from FT164's own build join the same visit. Verifying a ticket's
done-claim resolves every owner named in its Red-mutations table to a real
artifact in the tree, not to the delegate's demonstration of a red: the teach
ticket claimed a template-heading-depth needle and mutation row that never
landed, and three review rounds missed it because their enumerations keyed on
registered needles rather than on ticket-row owners — the same
resolve-the-identifier rule `craft-delegate` already states for absence claims.
And a repair round that extends an enumeration the spec pins trues that
enumeration up as part of the round: FT164's needle count drifted 24→27 across
repair rounds and was corrected only in round 3, and row 7.2's moved-artifact
clause was left vacuously satisfied by the same drift. Kit edit under the
`craft-synthesis` discipline. Sources: the retired `ft164-ticket-contracts`
spec's out-of-scope riders and its implementation retro, drained here;
`capture/learnings.md` 2026-08-03, verdicted here.

The gate-test-concurrency close sharpens that universal-claim rule: a structural
acceptance claim uses semantic enumeration or multiple representation mutations.
An audit bound to one reviewed function name can prove that exact restoration
while missing the same carrier expressed directly at another construction site.
Source: the `gate-test-concurrency` retro, drained here.

A mutation probe is evidence only after an independent read confirms the
intended mutation is present. Three wrap-spanning prose probes became no-ops and
looked green until the coordinator checked the file; `craft-delegate`'s
verification list therefore pairs mutation application with a presence check
before accepting the probe verdict. Sources: the Pocock-guidance-doctrine retro
and `capture/learnings.md` 2026-08-12, drained here.

The check-level-conformance-scoping close adds a slicing rule to that visit: a
ticket spanning runtime projection, canary selection, and documentation
conformance splits those owners before assignment. Its combined scope/report
slice consumed most of 14 assignment attempts even though the narrower
outer-selector, conformance-meta, and evidence-retention slices converged
directly. Source: the check-level-conformance-scoping retro, drained here.

The FT194 close adds one repair-slicing ruling to that visit: when a finding
extends an acceptance row already owned by an existing ticket, repair guidance
decides whether to reuse that ticket rather than create a second source for the
same row; a separate ticket still needs its own independently-green tracer
outcome. The run-specific instruction to preserve the review receipt before
assignment, integration, and fresh exact review is dismissed here as already
canonical in the current ticket/review cadence. Source: the
`ft194-project-green-desync` retro, drained here.

The 2026-08-03 lifecycle run adds six clauses to the same visit, all of them
paid for once already. Repair triage: when a review returns material blocking
findings, the coordinator shows the split — must-fix production defects versus
judgment-priced hardening — before charging repairs, because the finding's
disposition is the review's but build priority is the reviewer's; that run
routed eight round-1 findings into repair delegates undifferentiated.
Probe evidence: a probe's acceptance is proof it *executed*, not its exit status
— a capability-skip and a cached green both read as "ok", and a mutation probe
that skipped was read as survived twice; run probes through `bench test` or
`go test -v -count=1`, and treat `(cached)` or near-zero runtimes as invalid.
Delegate caveats: match each returned caveat against the charged rows'
*mechanisms*, not just their outcomes — "unreadable metadata refuses before the
planner is reached" contradicted a row promising real-planner composition and was
recorded as no-divergence, leaving a vacuous junction test for the review to
find. Registry fences: a ticket adding an entry to an extensible registry derives
its fence from the previous entry's landing commit, which is a complete wiring
manifest, and traces every crossing of the new entry's *name*; one such fence
missed three registration surfaces and forced a follow-up ticket. Repair tickets:
scope from the invariant the finding names rather than its cited lines, reread
each repaired artifact whole for internal agreement before resubmitting, and
express red-mutation greps against the false claim rather than the words it was
made of — an absolute "this grep must return nothing" over a phrase that
legitimately survives invites editing true prose until the tool goes quiet.
Prose-only steps: the contract-discovery skip was possible because a step with no
artifact leaves no trace of being skipped, and `Contracts:` fixed exactly that
for one step — any remaining prose-only step in the phase commands deserves the
same landing-site treatment, prefactor first. Sources: `capture/learnings.md`
2026-08-03 and both 2026-08-03 retros, verdicted here.

The ft156-anchor-registry close adds one wording repair to the same
`craft-delegate` visit: the stale-base check's charge text requires HEAD to
equal moving `main`, which is wrong for a lifecycle assignment — an assignment
builds on the run's candidate ancestry, and a delegate holding the
HEAD-equals-main rule reads a correct lifecycle base as stale. The charge
guidance should describe candidate ancestry for lifecycle work and reserve the
ff-to-main check for ordinary worktree delegations. Source: the
`ft156-anchor-registry` retro, drained here.

The exact-prospective-landing close adds four guidance faces to the same visit:
headless delegate wrappers distinguish produced edits from a stalled response;
`craft-tickets` folds review-only normalization records into substantive tickets
or deletes them before exact review; repair tickets name discovery, concrete
integration mappings, and a public mutation owner; and `craft-review` includes
post-authorization directory descendants in its concurrency inventory. Source:
the `exact-prospective-landing` retro, drained here.

The covers-traceability and FT195 closes add the charge-time proof details to
that visit: a delegate-supplied Go `-run` alternation is listed before execution
so a quoted literal `|` cannot produce a vacuous green; a ticket adding tests to
a classified family explicitly carries the sibling classifier registry as an
integration surface; and a coordinator pre-promote focused classifier run is
required when the composition adds tests to such a family. The existing
edge-inventory and production-reachability rules already cover the two FT195
repair misses, so those journal entries add evidence, not another rule. Sources:
the covers-traceability and go-build-cache-footprint retros and the 2026-08-06
learnings journal, drained here.

The bench-preflight and lifecycle-removal closes add two coordinator-proof
rules to this same visit. A mutation probe takes a fresh, round-unique backup
immediately before it changes an owned file, prefers owner-driven restoration,
and deletes the backup after restore; a wholesale deletion ticket carries an
explicit orphaned-caller sweep over every removed symbol. Both are charged
evidence, not session-local recovery lore. Sources: `capture/learnings.md`
2026-08-11 and both 2026-08-11 retros, drained and verdicted here.

The coherent-diff close adds two repair-proof clauses. Every identity-oracle
repair mutates the exact production dimension it claims to protect before the
coordinator accepts it; changing a correlated dimension is not evidence. A
delegate cap followed by a reviewer-approved tier escalation counts as a repair
round and carries `delegate-error` attribution, so the retro does not erase the
failed first charge. The existing ticket-breakdown review and three-axis
semantic implementation review remain separate cadences and need no duplicate
rule. Source: the axi-coherent-diff retro, drained here.

The canary planted-reason close adds two charge-time repair checks: a
conformance repair says whether package tests use synthetic registry data rather
than the live root and runs the public root conformance seam; after preserving a
review pickup, its reread compares the aggregate candidate with the slice base.
The producer-derived baseline probe is already required before ticket execution,
so its repeated reminder is dismissed rather than becoming a second rule.
Source: the canary-planted-reason-ownership retro, drained here.

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
`/bench-shape-idea`.** Upstream candidate (mattpocock/skills,
domain-modeling): as grill tickets resolve decisions, challenge fuzzy or
overloaded terms, stress-test the emerging model with concrete edge-case
scenarios, and keep `CONTEXT.md` plus applicable ADRs current inline. One
source per fact holds: the decision map owns the build decision, `CONTEXT.md`
owns ubiquitous language, ADRs own hard-to-reverse architectural state.
Integrates into the existing phase rather than adding a parallel skill. Kit
edit under the `craft-synthesis` discipline. Source: `capture/IDEAS.md`, drained
here.

The upstream refinement adds active modeling moves — challenge an overloaded
term with concrete counterexamples and update the language owner as the decision
settles. It folds into this row; the candidate's unrelated non-blocking fact
dispatch and git-hot-spot ideas are dismissed without session evidence.
Source: `upstream(mattpocock/skills@84fdeff)`, drained from
`capture/IDEAS.md` here.

**FT180 (LOW) — a spec-optional route decided at shape-idea's exit.** Between
the one-ticket light path and the full pipeline there is no middle route:
work too small for a spec doc but wider than one ticket either pays for a
spec whose stories are ceremony or takes the light path it doesn't qualify
for. Reviewer request 2026-08-01: when the shaped scope is small enough,
`/bench-write-spec` still creates the spec folder and its tickets but omits
the spec doc, and `/bench-shape-idea` closes by routing explicitly — direct
to tickets or a full spec — so the route is the shaping phase's exit product
rather than an ad-hoc call at build entry. The build's decision is the
threshold test the routing applies and how a spec-less folder interacts
with the `bench spec` lifecycle (which currently keys on `spec.md` status).
Kit edit under the `craft-synthesis` discipline. Source: `capture/IDEAS.md`
2026-08-01, drained here.

The lifecycle half has a live instance: shipped light-path changes hold
ticket-only folders under `specs/` with no `spec.md` — `bench spec history`
returns nothing for their slugs and `bench spec retire` cannot target them, so
the receipts are neither active specs, retained history, nor retireable state.
The route this row decides owns their terminal disposition: either one
existing spec reader indexes ticket-only receipts and states why they remain,
or the light-path close removes the folder after promoting durable content —
one owner, one policy, no second archive convention. Source:
`capture/learnings.md` 2026-08-01, verdicted here.

The route also owns the requested one-command execution surface. Once shaping
chooses spec-optional work, `$bench-implement-spec --full` accepts its ticket
folder and carries implementation, review, gate, and commit to push-ready state
without inventing a placeholder spec. This is a consumer of the route and
terminal policy above, not a second lifecycle. Source: `capture/IDEAS.md`,
drained here.

**FT182 (LOW) — a Planned-phase receipt over an absent target wedges the
abandon retry.** In `bench worktree` resume, a Planned-phase in-flight receipt
whose target is already absent returns `errStaleFingerprint` and wedges the
abandon retry — the crash window sits between the receipt write and
`checkpoint(Removing)`. Found by the FT176 plan-absent-target delegate as out
of its ownership fence; the fix lives in `internal/worktree/resume.go`, beside
the absent-target planning fix that landed at `dfcc71d`. Source:
`capture/IDEAS.md` 2026-08-02, drained here.

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

**FT166 (LOW) — `bench capture commit`: porcelain for the ambient capture
set.** Commit the capture surfaces (`capture/learnings.md`, `capture/IDEAS.md`,
`capture/session-handoff.md`, `capture/retros/`) with a conventional message under the
doc-only standing rule, so the plain-`git` step every session hand-assembles —
with the current empty-index and explicit-path safeguards attached — becomes one
sanctioned command. Weigh it against the standing guidance: the porcelain would
remove the instruction rather than duplicate it.

The capture-only commit-path clause folded into FT168 on 2026-08-01, where the
reviewer's scoping ruling now owns it; what stays here is the porcelain. The
fold also corrected the clause's premise: it proposed an exemption conditioned
on "no gate check reads those paths", and three of the four paths are graded
from the real root today, so the mechanism became a reduced phase set rather
than an exemption. This row's porcelain composes over the allowlist that spec
establishes and should be specified after it. The same-day journal entry stays
here: with two
unrelated dirty changes on `main` there is no sanctioned two-commit route
(`bench commit` refuses on out-of-set dirty paths, the guard refuses
`git stash`), and the session bundled with disclosure at `5fd3789` — until
this row ships, that is the convention to write down (bundle, lead with the
substantive change, name the ride-along, flag for veto). Sources:
`capture/IDEAS.md`, drained here and in a prior run; `capture/learnings.md`,
verdicted here.

The writer half joins 2026-08-01: capture writers are asymmetric — `bench
idea` and `bench handoff` own their files while `capture/learnings.md` and
`capture/retros/` are written by hand from phase prose that names the path,
the coupling that made the capture co-location cost ~278 references instead
of 4 constants. A sanctioned journal/retro writer is this row's second face,
distinct from the commit porcelain above. Source: `capture/IDEAS.md`, drained
here.

The FT195 close adds validation to the writer half: journal appends use the
parser's canonical open-entry grammar instead of allowing a malformed entry to
remain invisible until `bench status` notices it. Source: the
go-build-cache-footprint retro, drained here. The cancel-signal session
supplied the live instance: its 2026-08-09 journal entry used the `- date`
bullet form instead of a `## date — title  [open]` heading, so every
parser-backed reader reported an empty journal while the file carried an
entry; this drain found and verdicted it only by reading the file directly.
Source: `capture/learnings.md` 2026-08-09, verdicted here.

The phase-close capture recurrence belongs to the reduced-phase-set face: a
post-promotion retro or learning remains visible and reviewable, but does not
repay the promotion gate solely because it records that terminal outcome. Its
eventual reviewed drain still uses an ordinary authoritative landing path.
Source: `capture/learnings.md` 2026-08-09, verdicted here.

**FT168 (LOW) — focused iteration evidence: a fixture-selecting
canary.** Proving one changed
fixture costs the whole canary sweep: the light-path repair pass
needed evidence for a single race fixture, and the whole-sweep-only surface
invited expensive duplicate runs (the repair delegate launched one unbidden).
Add a `bench canary` path that runs one named fixture or family as iteration
evidence only — the full sweep remains the only thing the gate credits, so
this is a focused check, not a second oracle. The row's other face — a
reduced phase set for an allowlist-confined changeset — shipped 2026-08-01
as the reduced-gate-phase-set spec (`6f3486a`, spec since retired); this face was
explicitly out of that spec's scope and stays open, back at its pre-fold
LOW. Source: the light-path retro, drained here.

**FT140 (LOW) — review residuals that want a verdict, not a build.** Calls
from two resolution runs outlived their specs' retirement. The recurring one is
the provenance question, now at three instances: a test that is the real
coverage for a story, with no acceptance-coverage map row naming it, on a spec
that has since retired — so the map can no longer be edited and the reviewer
decides once whether such a test is accepted as standing coverage without a row
or gets its provenance noted in the family's decision record. Instances:
`TestAXILearningsWrongType` (FT86; the slice correctly declined to edit a
shipped spec's map beyond its mandate), and from `ft91-phase-manifest-dag`'s
review, `TestMergeEnvStripsThenSets` (the only non-vacuous coverage story 7 has,
per FT99's third instance) and `TestManifestDirResolvesAgainstGradedRoot` (the
real graded-root anchoring semantic; the mapped test graded a branch production
never reaches). One decision closes all three.

Four singles ride along. The orphaned-review-pickup signal
(`internal/status/status.go:534`, severity 9) pairs `reviews/*.md` against
`specs/<slug>/spec.md`, and neither side of that pairing holds today: `reviews/`
does not exist in the tree, and most `specs/` directories carry only
a `tickets/` folder with no `spec.md` at all. Half of that question is now
answered: the convention is live but unenforced, and the Pocock guidance
doctrine makes writing the file a required review step, so what stays here is the signal —
repoint it at the enforced shape once that lands, or cut it. Re-measured
2026-08-01. Source: `capture/IDEAS.md`, drained here.
`internal/gate/manifest.go`'s `dedupe` has no observable
effect — the scheduler's edge handling is already duplicate-tolerant and no
diagnostic renders `Needs` — but it implements a spec veto item literally, so it
is defensible dead code rather than a defect: keep or cut. And `bench learnings`
moved from a 5 MiB to a 2 MiB read bound: closing the divergence required
picking one number, and the slice chose the lower because `bench status`
already applied it to the same file — fail-closed and ambient-board-neutral,
but a 2–5 MiB journal that used to render now exits 1, a real behavior change
to keep or reverse. And whether `bench commit` staging retries briefly on a
held `.git/index.lock`: the stderr relay landed at `6d481eb` and showed the
field failure occurs only under a lock a concurrent session holds, and a
post-gate staging failure costs a full green run — retry, or keep fail-fast.
Source: `capture/learnings.md` 2026-08-02, verdicted here. One line each
closes this row.

The gate-decision review adds one bounded residual bundle under receipt
`80069545`: collapse duplicated projection strings, decide whether the helper
name should expose storage rather than test bytes, tighten the capture counter
to the Git tree, assert the partial-decision control, preserve the original error
in the nil-error failure, and decide whether the mixed-partition reload control
needs widening. These are reviewer-priced hardening calls on shipped test code,
not defects the retired spec still owns. Source: `capture/IDEAS.md` 2026-08-07
and the gate-decision-test-seam retro, drained here.

## False greens — verdicts that credit unchecked work

Five rows, one failure class: a green whose warrant is missing — a stale
binary, a dead or skipping citation, a vacuous baseline, an unchecked absence,
a dependency edge nothing resolves. Each hardens a different oracle surface, so
they stay separate builds, but they read and prioritize as one theme.

**FT133 (MEDIUM) — `bench coverage --check` verifies that red-signal citations
resolve.** A coverage-map row naming `go test -run TestFoo` where no such test
exists exits 0 with `[no tests to run]`, so a dead citation reads as a green
guard — a hole in the coverage oracle itself. FT86 story 16 shipped exactly
this: its named regression guard is a lowercase subtest dispatched via
`RunParallel`, and cannot go red as cited. Teach `--check` to resolve each
cited command to at least one matching test before crediting the row, and
decide the posture for citations it cannot resolve — fail closed is the
family default.

A second instance, drained 2026-07-26, widens the check past mere existence: a
row whose red signal cites a capability-gated test resolves fine and still
cannot go red as written. `go test ./internal/conformance -run
^TestRootConformance$` without `BENCH_CONFORMANCE_ROOT` prints `ok … 0.002s`
and skips (`bench-skip kind=environment`, visible only under `-v`), so a
session following the map verbatim reads a false green — reproduced through the
cited command on `main` at 3ea3abf. The check must therefore credit a row only
when the cited command actually executes its test, not merely matches one, and
a capability skip is the discriminating case. The authoring-side alternative —
a `craft-spec` rule that a conformance red signal always carries its
`BENCH_CONFORMANCE_ROOT` prefix — is the fallback if the mechanical check
proves too broad; prefer the check, which removes the instruction rather than
duplicating it. Related but distinct: `bench test` makes skips visible to a
reader, while this row makes an invisible skip fail the coverage oracle.
Source: `capture/IDEAS.md`, drained here.

The FT126 close repeated that exact false green: a scoped
`TestRootConformance` invocation omitted `BENCH_CONFORMANCE_ROOT`, skipped its
graded-root work, and still printed an `ok` summary. The mechanical check remains
the preferred single owner; if it cannot discriminate the execution posture,
the authoring-side rule explicitly requires a real graded root rather than
crediting the summary. Source: the FT126 recurrence-tallying retro, drained
here.

A third face, drained 2026-07-28, is row identity rather than citation
resolution: `bench coverage` emits only story/seam/red_signal, so a spec whose
rows share all three — `implement-spec-full-run`'s three story-3 hook rows —
cannot be enumerated row by row, and FT152's story-12 per-row accounting rule
is unexecutable as specified. Either the emission gains a stable row identity
(row number or the behavior field) or the rule names rows by story plus
behavior off the spec's own map — decide it alongside the check, same owner.
Found by the Codex falsification pass on `3eb1c9a`. Source: `capture/IDEAS.md`,
drained here.

The bench-preflight review found the same membership oracle half-disabled by a
mixed-tag map: rows-membership derives one `SpecTag` from `ids[0]`, so rows under
another tag escape the comparison. Decide here whether a coverage map must carry
exactly one tag or membership must partition and scan per tag; either answer is
enforced by the coverage validator and includes a mixed-tag red fixture. Source:
`capture/IDEAS.md` 2026-08-11 and the bench-preflight retro, drained here.

**FT174 (MEDIUM) — ticket files have one enforced dependency, ownership, and
mutation grammar.** Ticket files carry a `Blocked by:` field keyed by sibling
title, and no parser reads it. A retitle silently breaks the edge, and nothing
detects a cycle or a dangling blocker,
which under parallel assignment is not a drafting nicety: the frontier is
computed from those edges. The grammar to adopt already exists one directory
away. `internal/maps/schema.go:138` enforces `^(none|#[1-9][0-9]*(, #[1-9][0-9]*)*)$`
for decision maps and validates cycles and unresolved blockers over the ticket
graph, scoped to `decisions/` and never applied to `tickets/`. Extend that owner
rather than adding a second graph, and pair the check with `Ownership fence:`
disjointness between concurrently-eligible tickets, which is the same read and
the method step 2 of `craft-tickets` currently lacks. The doc half is done —
FT164's template teaches the identifier form and the parseable field shapes, and
the gate parses the Good example with the real parser — so this row is the parser
and the validation only. Measured 2026-07-31; re-verified against FT164's landed
template.

The parser visit closes three more real-ticket gaps. Once legacy compatibility
ends, every ticket carries an explicit `Ownership fence:`; assignment no longer
infers mutation authority from incidental `internal/<pkg>` prose, with rejection
timing left for the spec. `ParseTicket` reads Red-mutations rows from real tickets
and requires one complete row per acceptance ID, while historical closed tickets
remain exempt. Duplicate acceptance IDs are malformed input and the diagnostic
names both the ID and ticket instead of silently deduplicating evidence used by
coverage accounting, checkpoint receipts, and repair routing. Sources:
`capture/IDEAS.md` 2026-08-03, drained here; `capture/IDEAS.md` and
`capture/learnings.md` 2026-08-02, drained in prior runs.

The `Assumptions:` field is retired, and assignment now refuses a red-mutation
probe that crosses no path in its ticket's ownership fence. What stays here is
the dependency graph — `Blocked by:` parsed by identifier, cycles and dangling
blockers detected, and fence disjointness between concurrently eligible tickets.
The declared atomic closure graph is now assignment-enforced; the remaining
review check compares its `Closure:` inventory against `Contracts:`, the spec
coverage map, and the edge inventory before leasing work, because a parser
cannot detect a fact omitted from every declaration.
One residual has no owner at all: when a fence does turn out to be wrong mid-run,
`abandon` is whole-run and no verb releases a single open assignment to re-fence
it, so the correction costs a sibling ticket and another delegate cycle. The
shipped authoring-time refusal makes that rarer without making it reachable.
Sources: `capture/learnings.md` 2026-08-03 and 2026-08-04, verdicted here; the
`ft187-communication-surface-cut` retro, drained here.

The conformance-harness close adds one review-time ownership chain for every
cross-fence fact: raw declaration, resolving owner, exported resolved value,
consuming path, and reciprocal ticket edge. Assignment already enforces the
reciprocal edge; what stays here is comparing the complete chain against
`Contracts:`, coverage, and the edge inventory before leasing work. Source: the
conformance-harness-scope retro, drained here.

The one-shot close adds literal-deletion honesty to that same pre-assignment
review: when a mutation depends on deleting an exact anchor, count every
occurrence so explanatory prose cannot shadow it, and execute the promised
deletion red before accepting the ticket. Source: the `one-shot-feedback`
retro, drained here.

The repair-reslicing close adds two review obligations to the same ownership
chain: every atomic closure member gets its own mutation, and review evaluates
lifecycle state before treating a dependency as blocking. Unchanged consumers
remain integration surfaces; a path temporarily changed only for a required red
still needs fence authority. Source: the repair-ticket-reslicing retro, drained
here.

The spec-ticket handoff close adds a read-only ticket preflight, so malformed
`covers` annotations and contract or fence declarations fail before a
path-scoped planning gate. `craft-tickets` also expands each independently
removable member of a compound coverage row, including its retained-inventory
mutation fixture, before the first breakdown review. Source: the
spec-ticket-handoff-contract retro, drained here.

The bench-preflight close adds one hostile-class obligation: a line-oriented
grammar's mutation inventory names cross-line variants explicitly rather than
crediting a single-line red, and the ticket preflight carries the corresponding
multi-line fixture before review. Source: the bench-preflight retro, drained
here.

**FT177 (MEDIUM) — a stale `dist/bench` makes contract-test mutation probes
silent no-ops.** Any `internal/contract/surface/*` or
`internal/contract/runtime` test that execs the CLI grades the previous
build: `bin/bench.sh:181` execs `$k/dist/bench`, so editing Go source and
re-running the test answers for stale bytes. The failure is asymmetric — a
stale probe yields PASS, never FAIL — so a probe that reds is trustworthy
while a probe that passes may have tested nothing. Hit twice in one build:
once as a coordinator rejecting a delegate's correct evidence, once
reproducing that result only after adding the rebuild — three round trips in
the reduced-gate-phase-set build. Add a staleness check at the
contract-fixture seam that reds when `dist/bench` is older than any tracked
Go source under the subject root — one stat sweep, converting an invisible
false green into a red. Sources: `capture/IDEAS.md` and the
reduced-gate-phase-set retro, drained here. Recurred across the FT176 close:
freshness forced two manual rebuilds before `bench commit` would run. Source:
the spec-build-lifecycle-preconditions retro, drained here.

**FT103 (LOW) — existence-checked absence evidence: the gate half.** A
delegate's payload slice landed with a misspelled kit-only allowlist row
(`craft-synthesis` for the actual `bench-craft-synthesis`), so its contract
passed by asserting the absence of a path that never existed and the skill
kept shipping to consumers — a vacuous green the gate cannot see. The
charge-side half — absence evidence must name identifiers that resolve to
real things — shipped 2026-07-26 in the delegation-discipline batch, as a
clause in `craft-delegate`'s done-claim verification list. What stays here is
the gate check: confirm a per-source existence check on the consumer-payload
allowlist — the emptied-set vacuity closed with the FT85 fix commit, the
per-path existence guard is the remaining cheap single-source check.

**FT201 (LOW) — production cancel-signal registrations conform to one
source.** A conformance check that no production Go
`signal.Notify`/`NotifyContext` call site names a signal literal — the trapped
set comes from `subprocess.CancelSignals`, made the one source by the
cancel-signal-parity fix (`6fbf404`). `_test.go` files are excluded because
`internal/systemtest/owner_test.go` registers its own `os.Interrupt` as a
fixture. Deferred from the light-path ticket because the check turns one green
step into a migration plus a new gate rule. The same surface carries the open
reviewer decision that fix left behind: SIGKILL still orphans the detached
builder process group, and closing it needs `Pdeathsig`. Sources:
`capture/IDEAS.md` 2026-08-09, drained here; `capture/learnings.md`
2026-08-09, verdicted here.

## Reds the diff doesn't own — inheritance, load, and harness defects

Five rows, one failure class: a red that answers for something other than the
diff in front of the gate — an inherited baseline, machine contention, a
literal deadline, a harness defect, a flaky oracle.

**FT141 (MEDIUM) — `bench gate pin` records red verdicts,
so inherited reds stop reading as caused.** The pin records only the tree it
Occurrences: baseline-01
graded, not what failed there, so an agent that sees a failing check in a file
its diff never touched assumes causation and starts fixing — and that loop
does not self-terminate. Real case 2026-07-26: `main` red since `3c50349`,
and the FT91 build paid ~12 messages of git archaeology to prove the red was
inherited. Record which checks failed at which commit, so a stage starting
from a pinned baseline subtracts inherited reds automatically. This is the
containment half of the incident; the prevention half is the doc-only
shortcut's gate-anchored-surface exception.
Source: `capture/IDEAS.md`, drained here.

The gate-fastpath retro adds the output half of the same attribution problem:
`bench commit`'s final refusal must retain the failing phase name, and the
implementation landing guidance must capture the full gate log rather than
tail-filtering it away. Two reds with discarded phase output forced full
reruns before the failure could be attributed. Source: the gate-fastpath
retro, drained here.

FT150 folds in here, 2026-07-29: every surface that demands the pin quotes
the command and explains nothing — `internal/adopt/prepush.sh` refuses a push
three ways (unpinned, no `.bench` tree, tree mismatch) with "run `bench gate
pin`", and the `bench status` gate row names it the same way, so a first-time
user on the `bench setup` adoption path is never told the pin records the
`.bench` tree a human has reviewed, or why the refusal cannot clear
automatically. This row's build already visits the pin's data; add the one
explanatory clause at those refusal sites and the status row's action text in
the same visit, single-sourced rather than written four times.

An interrupted gate adds the destructive-diagnostic face: startup writes a
pending record over `.git/bench-last-gate`, so the prior verdict needed to
explain the interruption is lost. Preserve that record or bind its identity
inside the pending record before replacement. Sources: `capture/IDEAS.md` and
`capture/learnings.md` 2026-08-07, drained and verdicted here; the
gate-decision-test-seam retro, drained here.

**FT203 (MEDIUM) — pre-existing flake: `TestListCommandPublicRowsAndDisclosure`.**
Fails roughly 1 in 6–7 full `internal/worktree` package runs; reproduced on a
clean `git archive` HEAD baseline by the T1 repair delegate, so it predates the
axi-query-disclosure repair tickets. A flaky oracle in the gate's path is a red
that answers for timing rather than the diff, and it taxes every landing that
crosses the package. Build the tight red-capable repro loop first, then fix
against it. Entry: `/bench-debug`. Sources: `capture/IDEAS.md` and the
axi-query-disclosure retro, drained here.

**FT104 (LOW) — load-induced commit refusals: the stop rule and the pre-gate
quiet check.** Two faces of the same defect — a red answering for machine
contention rather than for the diff — and one owner, so they ship together.
First, the stop rule (this row's original charge). Retrying a recorded flake is
not iteration toward green: the FT85 review-fix commit was refused twice by
`TestFT78Story5ProofLedger` under gate load (green in isolation both times),
and the third identical run passed with no code change — ~35 minutes of wall
clock bought nothing. The gate/commit discipline states that when a commit is
refused twice by the same test already recorded as a known flake and proven
green in isolation, stop and hand the blocked commit to the reviewer with the
evidence instead of re-running. Replaces the retired FT95 "retry once"
operational line.

Second (drained 2026-07-27 from the learnings journal), the load the
coordinator itself caused. A write-delegate reported done while its own
`go test ./internal/...` sweep and two shell wait-loops were still running;
`bench commit` started immediately on its worktree and went red on
`internal/intent`'s concurrency test timing out — the same commit passed on a
quiet machine with no code change. A returned delegate is not a drained one:
the done-claim says its *report* finished, not that its subprocesses did.
`craft-delegate` already names the whole-tree gate a serialized resource for
*concurrent* delegates, so the prose half is extending that clause to the
sequential case — check for live test processes before gating, one `pgrep`
against a ~12-minute false red. Prefer the mechanical half if it holds up:
`bench commit` refusing or warning when it observes another `go test` against
the same module removes the instruction rather than duplicating it, and is the
same single-source preference FT131 and FT133 take. Kit edit under the
`craft-synthesis` discipline.

The gate-fastpath run supplied the concurrent form too: landing a ticket while
another write-delegate was still running heavy focused contract tests produced
the same load-coupled red. The serialized-resource clause covers active test
phases against the same module, not only concurrent gates or subprocesses left
behind after a done-claim; stagger the landing or retain the complete gate log
before classifying the red.

A third face, from the light-path retro (2026-07-29), widens the owner from
`bench commit` to every aggregate-evidence launch: the terminal repair
delegate started a strict full gate and a whole canary concurrently despite
its charge, and the coordinator had to stop the process groups by hand. The
serialized-resource rule covers standalone canary and `gate-phases`
invocations from delegated worktrees, not only `bench commit`;
`craft-delegate` names whole-canary and direct `gate-phases` runs as
full-gate work reserved for the coordinator; and a terminal repair pass is
one declared serialized stage — focused tests may run inside it, but
aggregate gate and canary evidence gets one coordinator-owned launch.
Source: the light-path retro, drained here.

**FT115 (LOW) — load-robust test and phase deadlines derived from bounds.** Two
literal deadlines flake under concurrent gate load.
`internal/gate/runner_test.go`'s `waitForPIDFile` hardcodes a 2s deadline for a
bash subprocess to write `.git/child-pid` — it flaked the FT87 slice 3 landing
gate (`TestFT78Story4ProofLedger/R11`), then passed 3/3 alone; this is the exact
defect class story 13 fixed for `WaitForTwoLegMarkers` (an outer test deadline
as a numeric literal rather than derived from the bound it must outlast). Extend
`bounds.TestDeadline` and the marker-wait conformance check to cover
`waitForPIDFile` and sibling literal deadlines. Separately, the gate's
conformance phase runs `go test` with no `-timeout`, inheriting the 10m default;
a 225–375s suite under a parallel worktree's load hit 600.013s and the gate went
red on a timeout that read like a failure. Give long phases explicit timeout
headroom, or have the runner distinguish a timeout verdict from an assertion
failure in its phase summary.

A third literal race reproduced about three times across twelve sanctioned
landing-gate runs in the gate-fastpath build:
`TestExecuteDeadlineRecordsDistinctTimeout` gives its stubbed gate timeout 50
ms inside a 500 ms parent context, but under the gate's concurrent artifact
load on WSL2 the parent can expire first and return GateExit 130 instead of
124. The test passed 10/10 alone. Derive both legs from the shared bounds so
the assertion preserves their ordering under the load the gate itself creates.
The FT123 + FT124 run adds the same repair at the worktree SIGINT seam: its
final descendant-liveness assertion produced one transient red during a
combined run, then passed alone and in every later gate. Poll to the existing
bounded deadline, as the test-report SIGINT contract already does, so scheduler
lag is not classified as a surviving process. Sources: the gate-fastpath
journal and retro and the FT123 + FT124 retro, drained here and in a prior run.

The first full gate run on a 12-online-core host adds two more literal-deadline
instances: the AGENTS.md and CLAUDE.md FIFO setup contracts both exhausted
their hardcoded 15 s subprocess bounds under the gate's outer-phase load, then
passed focused once and 3/3 repeated at about 0.43 s per run. Derive those
bounds from the shared test deadline rather than pricing a special-file
regression against a 16-core host; FT171 owns whether reducing outer
concurrency also lowers the contention that triggered them.

**FT120 (LOW) — gate, canary, and contract test-harness defects nothing
asserts.** Two
independent holes in the harness that grades the oracle, both found during the
FT91 canary-budget build. First, the R12 contention fixture can leak an immortal
process: its owner gate script waits on a release file with an unbounded
`while [ ! -f … ]; do sleep .01; done`, but that file lives in the test's
`TempDir`, so once `r12Contention`'s deferred release gives up after 5 s (or the
test binary is killed outright, skipping defers) and `TempDir` cleanup removes
the tree, the condition can never become true. Observed 2026-07-23: an orphaned
shell spinning 30 hours for 25 minutes of CPU, forking every 10 ms. The two
candidate fixes are not equivalent and the row must pick deliberately — bounding
the fixture's own wait is the only one that survives a killed binary, but it puts
a deadline on the owner half of a *contention* proof, so too tight a bound makes
R12 flake under exactly the load it exists to test; making the give-up path tear
down the process group (the profile's `gate-run` teardown rule, which the test
currently does not follow) is semantically free but does not cover the likely
trigger. Write the bound against a stated red signal rather than guessing a
number. Second, nothing asserts that each fixture keeps its own
`BENCH_CANARY_PHASE` pin under a concurrent multi-family sweep: `runFixture`'s
defensive `append([]string(nil), env...)` copy is the only thing keeping workers
off a shared backing array, and the gate runs `go test` without `-race`, so
deleting that copy would land green. The assertion belongs beside the existing
fake-`Runner` tests — and `internal/canary/canary_concurrency_test.go` is already
49 lines over its 400-line budget after the FT91 arms, so the added test lands
together with a `craft-seams` split (the file now carries two concerns: worker
derivation and bounds, versus inner-environment pinning) rather than an accept
entry. Third, found 2026-07-26 under full gate load:
`TestRuntimeGateContracts/bench_gate_rebuilt_self-host_contract` failed once
with a `TempDir RemoveAll` "directory not empty" cleanup error alongside the
contract failure, then passed standalone (all 14 subtests) and on an identical
whole-gate re-run. That cleanup message is the signature of a gate child
process still writing into the test's temp dir after the test returned — a
teardown race in the self-host contract, not a diff defect. An intermittently
red oracle is noise in the thing that defines done; pin the child-process
teardown before it recurs. Source: `capture/IDEAS.md`, drained here.

Fourth, the FT126 retirement gate reported one process-start-shaped
doctor-shim failure as an undifferentiated exit 1, then the unchanged case
passed 100 single-case and 50 full doctor-shim stress runs. Preserve process
launch failure as a distinct contract-harness diagnostic instead of collapsing
it into the child's exit status, so load and infrastructure reds can be
attributed without a blind stress investigation. Source: the FT126
recurrence-tallying retro, drained here.

Fifth, two conformance checks are reachable only through `conformance-suite`'s
whole-package run, matched by test-name prefix rather than registry
registration, so the registry undercounts what the gate enforces and a renamed
test silently drops a check. Register them. Source: the reduced-gate-phase-set
retro, drained here.

Sixth, carried as the capture's claim rather than verified cause:
`internal/worktree/worktree_test.go`'s `gitOutput` and its sibling `gitRun`
shell out with a bare `exec.Command().Output()` — no deadline, no
process-group teardown — while the `abandon_leftover_test.go` suite
deliberately plants no-writer FIFOs in the git administration directory, and a
git child that opens one blocks in `open()` forever without burning CPU. The
tests remove the planted record before reading, so the hazard window is a test
failing between the plant and its removal, or a killed test binary: the git
child is then orphaned with nothing to reap it. The 3+-day `git worktree list`
orphan found 2026-08-09 is not covered by the cancel-signal-parity fix, which
closed only the untrapped-signal class. The repro precedes any fix: force a
failure between the FIFO plant and its removal in
`TestReleaseRegistrationSkipsUnrelatedSpecialControlRecords`, or kill the test
binary in that window, and check whether a git child survives blocked on the
FIFO. Source: `capture/IDEAS.md` 2026-08-09, drained here.

The bench-preflight close adds two concrete harness holes. Its recursive path
walk lacks the nested-FIFO fixture needed to prove a discovered special file
refuses before read, and the release-tier root conformance run reported
`injected-port` finding no ports for `internal/canary`. Preserve the latter as
an observed `prep-release` red and diagnose it through that exact command before
changing derivation; the fixture gap can land with the focused preflight test.
Sources: `capture/IDEAS.md` 2026-08-11 and the bench-preflight retro, drained
here.

The Pocock-guidance-doctrine close adds a probe-isolation face: direct
`BENCH_CONFORMANCE_ROOT` runs expose families under a different scope and emit
about eight unrelated diagnostics, so focused mutation verification needs one
clean entry point whose executed family and result are unambiguous. This is a
harness seam, not another guidance rule. Source: the Pocock-guidance-doctrine
retro, drained here.

## Standards debt — one batched light-path pass

Three rows plus FT142's standards track are shippable together as small
one-source-per-fact and cleanup sweeps under one gate; FT117's parser-routing
half is the largest item in the batch. FT142 itself stays on the main list
because its ship track belongs to a separate `prep-release` hardening visit.

**FT117 (MEDIUM) — FT87 parser-surface follow-ups.** Two leaves left flat after
the slice 3 grammar centralization. The subcommand-routing registry's
`whyNested` exemption reason is free text nothing grades, and it currently
launders three hand-rolled parsers into "nested router": `internal/spec`'s
`specArg` (also missing the `--` rule story 5 promised), `worktree list`, and
`internal/adopt/doctor.go`. Route the two real leaves through `usage.Parse`,
correct the registry reasons, and consider grading the exemption reason itself;
`cmd/bench/main.go`'s worktree dispatch stays genuinely exempt. Separately,
`bench commit`'s usage-error line nests a second full usage line inside its
parenthetical when the fault comes from `usage.Parse` (e.g. an empty
positional) — one flat line naming the fault reads better.

**FT179 (MEDIUM) — comment quality: strip the reviewer-facing register,
document the bare high-stakes surfaces, sharpen `craft-comments`.** A
four-subsystem comment sweep (2026-08-01) graded roughly 9k comment lines:
the why/warning register dominates and restatement noise is rare, but three
fix passes fell out. First, register strip: ~145 `FT##` references plus
`Audit #N — tolerate`, `[R14]`, story/row provenance, port-narration
("as the Python shim did"), and design-advocacy essays — remove the
identifiers and arguments, keep the behavioral sentences; the one in-file
mutation transcript (`runtime_gate_env_sentinel_test.go:69-84`) moves to
history entirely. This reopens FT111's 2026-07-23 edit-in-place-only ruling
on new evidence: that ruling was priced against roughly a dozen sites, and
the measured count is an order larger. Second, doc comments where stakes are
highest and density lowest: `internal/releaseevidence` (bare syscall
numbers, the `stageOwnerBytes` sentinel), `internal/preflight` (`Setpgid`,
the undocumented `BENCH_PREFLIGHT_VULNERABILITY` override inside a security
phase), `internal/gate`'s public API (`Result`/`Execute`/`State`/
`Inspection`), `internal/contract` exports (`Env` nil-means-unset),
`worktree.Pool`'s checksum as stable on-disk identity, and `bench.sh`'s
`route_porcelain` (the porcelain-versus-plumbing dispatch axis). Third,
`craft-comments` amendments under the `craft-synthesis` discipline: name
identifier provenance as forbidden, add "state the constraint, don't defend
it", extend one-source-per-fact to comment knowledge, and qualify "a sparse
file stays sparse"; plus one ruling — where a demonstrated red's record
lives (commit or spec, never an in-file transcript). One site joins the second
pass from FT164's close: `ParseTicket`'s doc comment names its consumer rather
than its contract. Sources: `capture/IDEAS.md`, drained in a prior run; the
`ft164-ticket-contracts` retro, drained here.

The bench-preflight close adds its remaining comment debt to this pass: remove
the provenance narration in `gather_test.go`, correct the stale `RenderError`
doc line, and decide ticket-ID density in tests under the existing
identifier-provenance rule. The `fenceTokensInLine` parameter-list smell has no
behavioral defect and does not graduate independently. Source: the
bench-preflight retro, drained here.

**FT94 (LOW) — single-sourced `bench resume` summary
golden.** The resume summary line is asserted as a hardcoded exact-string
Occurrences: baseline-01
golden at four sites across three files (`internal/worktree/resume_test.go`,
`internal/worktree/lifecycle_policy_test.go`, and twice in
`internal/contract/runtime/runtime_worktree_test.go`), so a format change is a
multi-file hunt. Extract one shared expected-format helper: the unit and
runtime-binary seams stay distinct while the literal is single-sourced. This
is test-vs-test duplication, not the expectation-versus-implementation
independence the code standard protects, so collapsing it is consistent with
the one-source-per-fact rule.

## Session tax — evidence-supplied reader rows

This row is a measured, recurring reader cost from the week-of-2026-07-19
transcript evidence and builds on surfaces that already exist.

**FT125 (LOW) — reader surfaces that return the slice, not
the file.** Two existing readers make a session pull a whole file to use one
Occurrences: baseline-01
part of it. `bench spec show <slug> [--section stories|coverage|status]` — a
section-scoped spec reader; ~1 MB of spec `Read`s in the week of 2026-07-19,
with the `minimal-subprocess-data-exposure` spec alone read 15 times for 450 KB.
`bench outline --symbol <name>` — print a symbol's body with context instead of
guessing line ranges; 1247 manual `sed -n`/`cat -n`/`head` slices of source
files (2.28 MB) over the same week, and `bench outline` locates seams today but
does not return bodies. Both are refinements of surfaces that already exist, so
each earns its keep only if the row can show the slice is what sessions
actually wanted — the reads may be a search pattern that a narrower reader
would not shorten. Check that before building; the two arms share the
justification but not a seam.

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
| FT186 | FT108 | The gate restructure needs the mechanical exit test before moving oracle code without behavioral stories. |

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

### Goal tracks: guidance prose and the claim ledger

Two reviewer goals share one path: implement the guidance-prose backlog, and
ship `bench cite` (FT175) on an AXI-compliant CLI (FT173). The process
precursor for both is landed: FT164's ticket-contract core shipped 2026-08-03,
so every later build slices independently green tickets through the current
`craft-tickets` grammar. The payoff
facts shaping the order, verified in-tree 2026-08-02 with an independent
mid-tier refutation pass: `.agents/` and `.bench/BENCH.md` sit outside the
gate's reduced scope, so every separately-landed prose diff pays a full gate
— rows batch on the shared full gate, not just shared files; anchor-pinned
files couple prose diffs to conformance fixture updates (`craft-delegate` 14
anchors, `bench-implement-spec.md` 35+, `.bench/BENCH.md` 17); and AXI
action mechanics compound through every later CLI change. The reviewed
Pocock-alignment Spec C has shipped and FT107 is retired; both FT173
behavior-first builds (`axi-coherent-diff`, `axi-query-disclosure`) have now
shipped, leaving only the FT173 R11 residual. FT171's landed #23 and serial-gate
baseline make #24 moot; reconcile the map, then resume #25 and run #26 before
pricing outer concurrency.

1. Reconcile FT171's `decisions/gate-budget.md`: close landed #23, retire moot
   #24 and remove it from #26's blockers, then land #25's cuts and run #26's
   exact post-route census before pricing outer concurrency. FT141 may proceed
   independently where capacity allows.
2. Shape FT175's decision map now — its consumed owners (truncation,
   aggregates, `help[]`) landed with the axi-query-disclosure capstone, the
   condition the 2026-08-02 reviewer ruling set ("interfaces settled" means
   the consumed owners, not the whole FT173 row). Then write the FT175 spec
   and build it as vertical green slices: file evidence plus strict store and
   `claim show/check`; assessments with spans and absence scopes; command
   evidence plus replay and staleness; supersede/retire reachability; the
   complete AXI list/detail/status surface; deterministic gate and ambient
   status integration; then one local contradiction as the dogfood proof.
3. Shape FT198 once FT175 and FT171 are moving; the
   board's 170 KB full snapshot now confirms the progressive-loading trigger.
4. FT100's remaining work grills and builds last, after FT89 establishes which
   guidance is authoritative and the reviewer resolves FT170's benchmark route.

Fold FT106's verified-document vocabulary and FT162's exact subject binding
into the FT175 spec instead of building them as prerequisites; FT99 rides
prose batch 1, and its uncertainty obligation folds into the FT175 spec
where claims consume it. FT133 remains parallel evidence hardening; FT71
stays deferred behind its existing FT169 recommendation. FT172 is outside
this critical path; the FT156 anchor registry shipped, so section-scoped
`.bench/BENCH.md` anchors — the exact surface the prose batch edits — are
now fixture-proven.

## Recommended sequence

1. `/bench-shape-idea` — reconcile FT171's landed #23, moot #24, and remaining #25–#26 in `decisions/gate-budget.md`.
2. `/bench-debug` — FT203's reproduced `TestListCommandPublicRowsAndDisclosure` flake (~1 in 6–7 full `internal/worktree` runs on a clean baseline).
3. `/bench-shape-idea` — FT175's claim-ledger decision map; its consumed owners (truncation, aggregates, `help[]`) landed with the axi-query-disclosure capstone.
