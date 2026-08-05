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

**FT194 (HIGH, reproduced defect) — an empty-run fast-forward desynchronizes
project-green and wedges promotion repo-wide.** When an empty spec-build run
fast-forwards its recorded base to a moved branch tip, the project-green marker
stays at the prior base. Promotion later supplies the new run base as
`authorization.Bootstrap`'s expected marker and refuses with `project-green
marker conflicts with another tip`. The same transition owns no repair path:
checkpoint, assignment, review, and promotion all refuse after the tip move,
while a fresh gate cannot advance the marker and abandon discards the reviewed
composition.

The red signal is the sanctioned lifecycle itself: start an empty run from a
current project-green tip, advance the destination through an ordinary gated
commit, let the run fast-forward, then promote. The promotion must not conflict
with its recognized ancestor marker. Fix the state transition rather than add
an operator rule: either advance base and marker together or let bootstrap
accept a marker recognized as an ancestor of both the run base and destination
tip. Because this changes the promotion authorization state machine and its
only current escape is destructive, start with `/bench-write-spec`. Source:
`capture/learnings.md` 2026-08-04, verdicted here.

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

**FT156 (MEDIUM) — the anchor mechanism is weaker than its coverage rows
claim.** Two faces from FT152's build, one owner (the conformance anchor
checks), one decision. First, `requireCollapsed` does not strip HTML comments
while `forbid` does, so any anchored guidance sentence can be commented out and
still satisfy its anchor — the three new point-of-use verify hooks and every
shared-rule marker included; fixing the asymmetry touches roughly 100 existing
anchors across the kit, so it needs its own spec. The
`shape-idea-bypass-wrapped` fixture proves the forbid direction only. Second,
substring forbids die to paraphrase: "Plain invocation also starts the full
run" evades the "by default" needle, and each of the ten `--full` forbids has a
natural evading phrasing — the spec's coverage rows claim those pairs close the
cheap wrong implementation, but they close one spelling of it. Either the rows
state the weaker claim honestly or anchors get a stronger mechanism than
substring matching; that ruling also decides how much the comment fix is worth,
so both faces take one grill. A third face, drained 2026-07-28, is a coverage
hole in the same mechanism: no canary fixture proves a section-scoped anchor
over `.bench/BENCH.md` — story 9's four fixtures prove section scoping only in
`.agents/commands/`, and the shared-rule fixtures prove the marker directions
only, so FT152's story 1 (Workflow section) and story 15 (communication rules)
both rest on that untested combination. One fixture removing a section-scoped
`.bench/BENCH.md` sentence covers both; it rides whatever mechanism the grill
decides. A rider drained 2026-07-29: `bench anchors <path>` — query which
conformance anchors pin a doc file before editing it — needs exactly the
registry-or-parser seam the mechanism ruling decides, so the grill weighs it
as a consumer of that seam rather than a separate build. The mechanism now has
a measured cost signal: `internal/conformance/docs_workflow_helpers_test.go`
holds the hand-written anchor entries and stands at 864 lines against a 660
grant, a violation FT164's build widened and could not fix in place. A
declarative anchor registry is the fix; raise the grant to ~850 only as a
stopgap if structure noise blocks a pass before this row lands, never as the
answer. A verified false green adds one rider: deleting only the ticket
template's `Contracts:` line passes the real graded-root
`TestRootConformance` because no section-scoped template requirement anchors
that line. FT156 owns the oracle gap; this rider records it without changing
the check.
Entry: `/bench-write-spec`, by reviewer direction 2026-08-03: the mechanism
ruling is taken as the grill at spec entry rather than as a separate shaping
session, because the two faces are one decision and the spec cannot be written
without it. Source: `capture/IDEAS.md`,
drained here and in prior runs; found by the Codex falsification pass on
`3eb1c9a`.

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
diagnostic. The spec-backed half of the first edit has since been answered by
`spec-integration-gate-cadence`: `bench spec build promote` gates the exact
prospective implemented tree and is the sole project-green transition for a
reviewed spec build, so what stays open here is the light-path and non-spec
close, where the final-check tree is still resolved by hand.

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
retired spec currently makes `bench spec build status <slug> --full` unreadable
even though the durable run record survives, and that record omits the measured
promotion-stage timings the required retro needs. The terminal projection reads
the retained record after retirement and retains those timings. Until it does,
`/bench-final-check` captures the retro before retirement rather than rerunning a
successful promotion to manufacture evidence. Source: the
check-level-conformance-scoping retro, drained here.

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

**FT135 (MEDIUM) — a pre-push guard on a guessed branch looks armed while
protecting nothing.** When the repository has no resolvable default branch,
the installed pre-push hook falls back to a baked `fallbackProtectedBranch`
of `main` (`internal/adopt/link_hook.go`). The hook re-resolves `origin/HEAD`
live on every push and reaches the baked token only when that lookup stays
empty — but in exactly that repository the guard looks armed while defending
a branch that may not be the default, FT86's own failure class one layer
down. `bench doctor` (or `bench link`'s output) should report whether the
installed guard's protected branch was resolved or guessed, so the false
armor is visible where the reviewer looks. Source: `capture/IDEAS.md`, drained here.

A second face, drained 2026-07-29: currency, not resolution. This repo's own
installed hook is a bench-managed copy from Jul 6 predating the static
manifest header, so `bench guards` and every SessionStart banner report
`pre-push: no manifest` with an empty boundary — the guard reads as inert
while it in fact blocks pushes to `main`, and one session told the reviewer
`main` was unprotected. `bench doctor` reports ok because it checks only the
`bench:managed` marker, never content currency against the embedded template.
Same owner, same fix surface: doctor (or link's output) reports template
currency alongside resolved-versus-guessed. The local instance is repaired by
re-running `bench link` (`installPrePush` overwrites a managed hook) — offered
alongside this drain rather than performed silently, since it rewrites an
enforcement hook. Source: `capture/IDEAS.md`, drained here.

A third face, drained 2026-07-29: the sanctioned repair route refuses on this
repo. `bench link` aborts with `conflict: .claude/commands/bench-debug.md has
a symlink parent directory` — `.claude/commands` is a symlink to
`.agents/commands`, verified in-tree — and the capture records the refusal
landing before the pre-push hook refresh, so the second face's repair
(re-running `bench link`) is unavailable exactly where the stale hook lives; a
hand-copy fallback was used 2026-07-29. Whether link should traverse the
symlinked directory, skip already-converged files, or order the hook refresh
ahead of the conflict check joins the same doctor/link visit this row owns;
the abort-before-refresh sequencing is the capture's claim, not re-verified
here. Source: `capture/IDEAS.md`, drained here.

Staged spec: [`specs/pre-push-guard-visibility/spec.md`](specs/pre-push-guard-visibility/spec.md);
its decision map moved under that folder with the staging commit.

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
the default branch by diff, yet `bench worktree recovery <ref> --apply
<fingerprint>` still returns `retain` because landed-proof requires the
payload commit itself (observed 2026-07-20); recurred 2026-07-22 when
`git cherry` missed reshaped commits and the reviewer had to hand-delete refs
and intent entries, the exact manual surgery the lifecycle exists to prevent.
Both routes shipped 2026-08-01: `LandedInDefault` proves a squash-landing by
reverse-applying the branch's cumulative diff against the default tree
(`efb456c`), and `bench worktree clean --discard-branch` is the
reviewer-supplied proof for what no derivation can establish (`37411a0`) —
fail-closed stayed the default and every ambiguity still resolves to
not-landed. The unprovable half shipped 2026-08-04 (`fafb049`, from the
`recovery-discard` spec retired here): `bench worktree recovery <ref>
--discard <fingerprint>` retires one inspected payload per invocation without
asserting it landed, the plan separates an orphaned ref from an absent one and
reports how many paths the payload touches so the operator is not choosing
blind, and `bench spec build reclaim` deletes the provisional residue of
terminal spec-build runs. What remains of this face is the drain itself, and it
stays reviewer-owned because discarding preserved work is a judgment no
derivation can make: `bench resume` re-preserves the whole backlog at every
session start, and most of those rows name recovery refs that no longer exist,
so the pass is per-ref over what `bench worktree recovery` reports for
`refs/bench/recovery/` and the open assignment rows. Sources:
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

The FT128 close supplied the base-ref face, and the `bench spec build` family
has since covered only the spec-backed half of it. `bench worktree create`
always roots at the default branch with no base-ref flag, and
`bench worktree release` refuses while an assignment branch has not landed
there, so a chain of tickets on an integration branch cannot cut a worktree at
the chain tip and cannot retire any worktree until the reviewer merges; that
build formed the chain by merging each previous assignment branch into the next
worktree, which is the hand-run form of a compare-and-swap integrate. Reviewed
spec builds now get that surface from `bench spec build assign` / `checkpoint` /
`integrate`; light-path and non-spec chains still do not, and the landing
command owns whether they should. Source: the FT128 implementation retro,
drained here.

That face recurred on 2026-08-01 during per-component-gate-scoping and sharpened
into a precise gap: the core already carries the parameter, and only the CLI
withholds it. `worktree.Create` takes a variadic start ref, but the sole caller
that supplies one is `bench spec build assign`; the generic
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

**FT188 (MEDIUM) — parallel sessions land through an exact prospective tree.**
A path-scoped `bench commit` gates the whole working-tree snapshot and only then
stages its named paths, so any unrelated dirty or untracked file in the checkout
refuses the landing: checkout cleanliness is a repository-wide writer lock, and
the 2026-08-03 build recorded seven such refusals in
`capture/parallel-session-friction.md`. The staged spec introduces one exact
prospective landing substrate — start from an explicit expected base, overlay
only the caller-attributed named-path content, gate that immutable tree, publish
by compare-and-swap of the destination ref — and adopts it under ordinary
`bench commit` without changing the argument grammar or the attribution rules.
This closes the foreign-dirty face FT169 carries; FT169 keeps the landing
command's authority, recovery, and preflight questions.

`bench spec build promote` is the second consumer. Its clean-checkout
precondition composes badly with its recompose-discards-review rule: unrelated
dirt blocks promotion, and clearing it by committing moves the branch tip, which
discards the bound review and buys a full extra review round — observed once,
costing a sol round plus a second reviewer-run promote. Either promote reasons
about ownership the way the substrate does, or its refusal names stashing as the
cheap route so the reader does not reach for the expensive one.

The downstream lifecycle now has two reproduced pinning failures to close with
that same ownership model. A commit may edit the staged spec of an active run
successfully, only for checkpoint and abandon to discover the pinned-blob drift
after assignments have finished. A commit made to the destination after a run's
base may also disappear silently when promotion publishes the older candidate
composition. The multi-coordinator scope must therefore expose run-pinned paths
at write time and refuse or explicitly re-pin them, while recomposition detects
and preserves destination changes instead of overwriting them. The staged first
scope remains ordinary exact prospective landing; these lifecycle consumers stay
in the ordered downstream scope. Sources: `capture/learnings.md` 2026-08-04,
verdicted here.

The remaining decisions stay in the map at
`decisions/parallel-session-landings.md`: interrupted-mutation recovery and
abandon, which gate evidence may cross composed landing subjects, and the AXI
status surface for concurrent runs and conflicts (blocked on the first two).
Staged spec: [`specs/exact-prospective-landing/spec.md`](specs/exact-prospective-landing/spec.md).
Sources: `decisions/parallel-session-landings.md` #10; `capture/IDEAS.md` and
the `injected-interface-junctions` retro, drained here.

**FT195 (MEDIUM, staged) — compile one Bench executable and keep oracle verdicts
out of Go's test cache.** The canonical builder currently compiles the requested
binary and then `go run`s a second publication helper. The staged spec replaces
that helper with publication owned by the newly built host binary, gives
artifact callers an explicit unsealed mode, and removes reusable successful Go
test verdicts from the two oracle-owned test invocations with `-count=1` while
retaining compilation caching.

The same economics recur after a gofmt-only correction: the source seal must
change, so the safe path still rebuilds and re-gates. Neither proposed shortcut
is sound — a mutating gate changes its own subject, and a seal that ignores a
source change no longer attests to the exact build inputs. This spec instead
reduces the cost of each required rebuild without weakening subject identity;
fresh-assignment readiness remains FT169's preflight concern. Staged spec:
[`specs/go-build-cache-footprint/spec.md`](specs/go-build-cache-footprint/spec.md).
Source: reviewer-confirmed conversation and `capture/learnings.md` 2026-08-04,
verdicted here.

**FT178 (MEDIUM) — `bench worktree`'s bare verb is a human porcelain that
traps automation and leaks on signals.** Reviewer ruling 2026-08-01: humans
should not be driving the `bench` CLI in the vast majority of cases, and
worktrees are not an exception — the agent-facing surface is the subcommands,
and agent worktree creation already flows through `bench spec build assign`
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
discovery-convention prose half lives in FT107's sixteenth clause. Sources:
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
the pairing is broken from both ends. The roadmap end: the FT135 row omitted the
staged `specs/pre-push-guard-visibility/spec.md` path the preamble requires, and
was corrected here only because the session happened to know the convention. The
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

**FT173 (MEDIUM) — the AXI contract has ten principles, one derivation each.**
The kit implements its own published contract partially and unevenly, measured
2026-07-31 against the ten principles at axi.md. Five faces, one owner surface.
First, `craft-cli` enumerates only seven: content first, contextual disclosure,
and the per-subcommand help fallback are absent from the skill, so the guidance
cannot ask for what it never names. Second, contextual disclosure is emitted
nowhere — `help[` returns zero hits across non-test `cmd/` and `internal/` — and
it is the one principle with no partial implementation to consolidate. Third,
truncation has four independent derivations: `sanitize.Preview` (13 non-test
callsites across 8 packages), `internal/roadmap`'s package-private `limited()`
(9 callsites, 1 package), `internal/worktree/subshell.go`, and
`internal/outline/outline.go`. `sanitize.Preview` is already the shared one and
is the consolidation target; treating `limited()` as the seed inverts the actual
usage. Fourth, pre-computed aggregates have no shared helper and are derived per
command — `outline_meta`, `internal/publication`'s `next_action`, and
`internal/roadmap`'s byte and occurrence counts each roll their own.

Fifth, routine Git inspection still escapes the AXI surface. The 2026-08-02
decision-map gate-scope build needed repeated batches of `git status
--short --branch`, `git diff --stat`, `git diff`, `git diff --check`, `git
rev-parse`, `git log -1`, and `git show --stat` to answer two questions: what
changed, and what exact body should be reviewed. `bench diff` already owns
review-base resolution, the changed-file inventory, the landed-commit log, and
the patch body, so extend that command rather than minting a parallel `bench
git` namespace or a second derivation of those facts. Its default response
becomes one coherent snapshot with a compact revision row, pre-computed commit,
file, insertion, deletion, staged, unstaged, and untracked counts, the existing
base-to-worktree `files[N]{status,path}` inventory augmented with untracked
paths, a minimal checkout table that distinguishes index from worktree state,
and a whitespace-check result. Reuse `internal/git.Facts` — already the NUL-safe
checkout owner consumed by `bench roadmap --context` — for branch, divergence,
and XY status, and reuse `internal/diff`'s range owner for base-relative facts;
do not add another porcelain parser. `bench diff --full` returns the same
snapshot plus the existing log and an exact patch that does not silently omit
untracked regular files; `--commit` remains the post-landing view. Capture HEAD,
index, and worktree identity around the reads and retry or return a structured
drift error if they move, so one invocation cannot splice facts from concurrent
states. The acceptance target is one `bench diff` call for orientation and at
most one `bench diff --full` call for bodies, replacing the raw status/stat/
name-only/revision/check sequence. Grade it with a paired-delta fixture against
raw Git over committed, staged, unstaged, untracked, rename, deletion, binary,
hostile-filename, clean, and mid-read-drift cases.

Reviewer constraint, 2026-07-31: the consolidation changes call sites only and
makes no modification to AXI responses. The four truncation derivations and the
aggregate helpers route through one shared call without altering emitted bytes,
so the contract suite's pinned stdout stays green and retrofit scope never
becomes a compatibility question. That constraint has one consequence worth
stating rather than discovering: principle 9 has no existing implementation to
route through, so emitting `help[]` necessarily changes bytes and cannot ship
under it. Reviewer ruling 2026-08-02: the byte-preservation constraint is
relaxed for the `help[]` spec specifically, and for no other face — the
foundation still ships without changing emitted bytes. The Git-inspection
face is likewise byte-changing and therefore gets its own spec; it does not
relax byte preservation for the foundation.

A sixth face, observed 2026-08-02 during the FT164 build, is the canonical
contextual-disclosure example: `bench spec build start` refused with "no exact
green evidence: run bench gate, then retry start" on a tip whose verdict was
reduced — where plain `bench gate` can only ever re-record another reduced
verdict, so the stated remediation cannot succeed and the working command
(`bench gate --fresh`) went unnamed. A refusal or result whose sensible next
step depends on state the command already knows must emit that step in
`help[]`. Reviewer acceptance conditions, recorded 2026-08-02: the FT173 spec
is accepted only with a detailed per-command list covering every `bench` CLI
command — its follow-up commands and `help[]` contextual-disclosure surface,
command by command, no sampling — and the spec's grill must open on a deep
dive into each command's real usage drawn from the Claude session logs plus
the Codex usage logs the reviewer will share.

Two constraints shape the build rather than the diagnosis. Principle 3 should
not double-truncate a value already bounded at capture, and principle 8 is a
query-surface rule, not a binary-wide one — `craft-cli` already holds that the
contract attaches to the surface, so an operational command answering with live
data on no arguments is a regression, not conformance. The
truncation half alone is a one-source-per-fact sweep and could ride the standards
debt batch; the rest is a build. Sources: `capture/IDEAS.md`, drained here; the reviewer
constraint above, parked and drained here.

**FT185 (MEDIUM) — gate results join the structured Bench output contract.**
`bench gate` is the last major agent-facing surface that reports phase verdicts
as ad-hoc prose while `bench test`, `bench diff`, and `bench coverage` emit TOON.
Give the gate one structured result schema without changing exit-code authority,
phase completeness, or the durable verdict it authors. The gate-pipeline map's
ticket 9 closed the scope decision: no output redesign rides that pipeline build,
so this is an independent item. `bench spec build promote` is the second
surface: it emits only its TOON status line, so the phase evidence its gate run
produced never reaches a retained surface, while the `bench commit` on the same
tree prints per-phase evidence inline (`phase conformance: green`, `gate:
green`). One schema covers both. The retained-record half of the same complaint
— promotion-stage timings the required retro needs — is FT162's, not this row's.
Entry: `/bench-write-spec`. Sources:
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

**FT187 (MEDIUM) — cut the communication surface to the rules that fire.**
`.bench/BENCH.md`'s "How to talk to me" runs sixty of the file's two hundred
sixty always-loaded lines and contradicts itself in three places, so a session
resolves the conflict toward whichever clause it read last and the voice drifts
turn to turn. First, the formatting rules point opposite ways: "Format for scan:
tables and lists — use them", the cohesion clause's "bullets or tables only for
genuinely parallel facts", and "one-sentence paragraphs stacked in a row read as
a formatting error" together leave almost no sanctioned shape for a three-fact
update. Sharper, the progress clause demands bold **Status:**/**Next:** labels
"even if the entire update fits in one sentence" while the closing bullet asks
for a senior colleague's voice "not like this kit" — a one-sentence update
wearing two bold labels is the register the closing bullet forbids. Name one
governing rule (recommend the colleague-voice bullet), demote the labels to
updates that carry state worth scanning, and drop the "even if" clause. Second,
the section's own trigger is unobservable: "an in-progress Bench phase update"
asks a session to judge whether the conversation is a phase, and the 2026-08-03
review session that produced this row read the section and chose plain prose
over `## Result` on exactly that ambiguity. Bind it to something visible — a
phase is active when a `/bench-*` command was invoked this session, and not
otherwise. Third, Roles' "NEVER assume, always verify" is the same defect class
as FT107's first clause rather than a separate one: every session assumes
dozens of things per turn, so an absolute violated hourly teaches that the
file's absolutes are rhetorical, which is expensive standing beside four
invariants whose absolutes are real. Bind it to what it means — never assume the
reviewer's decisions, and never assume a claim the gate could check instead.
Fourth, invariant 2's declaration is only half enforceable: a session cannot
switch the main loop's model, only the reviewer can, so the declaration binds
delegates and headless runs and is ceremony for the main loop. Say which half
binds, or the ritual is performed without constraining anything. Fifth,
wording only: the bold "you must get an explicit OK *before* skipping canonical
steps" reads as absolute against the table directly beneath it that grants the
skip standing approval, and on a fast read the bold sentence wins and the table
goes unused — fold the exception into the sentence. Substance there defers to
FT156's named-lighter-path ruling and FT180's spec-optional route; this clause
only reconciles the two statements as they stand. Deliberately excluded: the
`CLI Inventory` section's hand-maintained duplication of `bin/bench.sh` (no
conformance check enforces the declared sync) is a second instance of FT89's
inventories half and stays there. This row takes FT100's `.bench/BENCH.md`
slice early because the conflict is self-contained in one section and does not
wait on FT89 settling what is authoritative elsewhere. Kit edit under the
`craft-synthesis` discipline. Source: session conversation 2026-08-03, parked
here by reviewer instruction rather than through a drain.
Staged spec: [`specs/ft187-communication-surface-cut/spec.md`](specs/ft187-communication-surface-cut/spec.md).

**FT107 (MEDIUM) — the standing guidance rules, batched: one
always-loaded-prose diff.** Sixteen remaining clauses edit the same
standing-guidance surface — `.bench/BENCH.md`, the phase prose beside it, and
two craft skills — and collapse into one batched kit edit under the
`craft-synthesis` discipline: one spec, one review, one gate. First (was
FT110), two generation-shaping clauses on an
observable. Invariant 4's "read the surrounding code before you write" is true
but unfalsifiable — *surrounding* is undefined and a session always believes
it read enough — so add the sharper pair: in `.bench/BENCH.md`, as a clause on
invariant 4 rather than a new invariant, never call an API or function whose
definition you have not read this session, and verify behavior in source
rather than from memory; and in `craft-seams`, beside the seam-finding
guidance, an exploration budget — a small stated number of reads without
traction means stop reading, run `bench outline`, and pick the seam from the
index, with the budget declared as spent in the transcript so the behavior is
observable rather than counted internally. The re-scope target is computed and
cannot go stale, which is the half worth having and also real pull for
`bench outline`, currently underused. Second (was FT119), plain-`git` commit
safety during a squash-merge landing. The doc-only plain-`git` convention is
safe only when the index is otherwise empty: after `git merge --squash` the
index already holds the whole merged slice, and a bare `git commit` of a
just-added capture file once swept 649 insertions across eleven files into a
commit labelled "capture: …", landed with no gate grading it — `bench commit`'s
attribution check, which would have caught it, was bypassed by plain `git`.
Note in the phase guidance that the convention holds only with an empty index,
and that during a squash-merge landing every plain-`git` commit uses the
explicit pathspec form (`git commit -m "…" -- <path>`). Third (drained
2026-07-26 from the learnings journal), the self-contradicting spec. The
batch-approval clause covers a reviewer who has gone AFK, and says nothing
about a spec whose own sections disagree — during the delegation-discipline
build the spec placed two canary fixtures in one family while its Prior art
line pointed at another, and the build corrected to the tree's convention,
flagged the choice, and continued without waiting. State the rule beside that
clause: when a spec contradicts itself and the readings are functionally
equivalent, take the one consistent with the tree's existing convention, flag
it for veto, and build on — a round trip buys nothing where no behavior
differs; the flag is mandatory, and any behavioral difference between the
readings is a stop-and-ask. Fourth (parked 2026-07-26 from session
conversation), the fix-loop escalation trigger — the symmetric partner to
`.bench/BENCH.md`'s read-budget reroute. A
fix loop has no equivalent stop, so an implementation chasing a defective seam
grinds indefinitely: every iteration edits real code against a real red, which
feels like progress the whole way down. The discriminating observable is not
iteration count — `craft-line` already caps those, and a cap governs spend
rather than diagnosis — but whether the red set *shrinks*. A converging loop
reduces it; a loop chasing a spec defect rotates it at constant size, because
the seam cannot satisfy two things the spec asked for at once. That is visible
in gate output and needs no self-assessment, which is this row's design
criterion. The reroute differs from the first clause's: not onward to
`/bench-debug` but stop and surface, since a seam judged wrong is a
spec-sign-off decision the reviewer owns. The clause must also gate
`craft-line`'s escalation ladder, which today escalates a tier on the second
red at the same tier — against a spec defect that buys a more expensive wrong
answer faster, so classification precedes escalation. Note that this is the
one clause in the batch reaching past always-loaded prose into a skill.
The shrink measure is only meaningful over reds the diff owns — inherited and
spec-predicted reds are constant noise that would trip a false stop. Evidence
for both halves is the 2026-07-26 FT91 gate: three reds, one
inherited from `3c50349` and two predicted by the spec, none belonging to the
diff, resolved only by hand-run git archaeology. Source: session conversation,
parked here by reviewer instruction rather than through a drain. Fifth
(drained 2026-07-26 from the learnings journal), the doc-only plain-`git`
shortcut needs a gate-anchored-surface exception. Commit `3c50349` rewrote a
craft-delegate paragraph whose exact sentence `checkWorkflowAnchors` pins,
plain `git` committed it with no gate, and `main` stayed red for three commits
until the FT91 build paid the diagnosis detour to prove the red was inherited.
Kit prose is gate-anchored content, not inert text: the conformance phase
asserts exact sentences from `.agents/skills/`, `.agents/commands/`,
`.bench/BENCH.md`, and `projects/*.md`. Either name the cheap targeted check
those edits must run before a plain-`git` commit (a scoped
`go test ./internal/conformance -run` over the anchor checks), or scope the
shortcut to surfaces no conformance check reads — reviewer's call which,
decided when the batch builds. The containment half of the same incident is
the gate pin recording red verdicts so inherited reds stop reading as caused.
Sixth (drained
2026-07-27 from the learnings journal), a mandatory step written as a
subordinate clause is read as advice. `/bench-review-implementation` step 5
says the `reviews/<slug>.md` pickup is "tracked state at birth", committed in
the session that writes it — the clause exists so findings survive a session
death mid-repair — but it trails a long paragraph and the FT148 review read it
as advisory, wrote the pickup, and went straight into the repair pass, so the
artifact lived and died untracked. Net tree state was identical and the failure
mode was not exercised, which is exactly why the wording rather than the rule is
the defect. Name the commit as its own ordered step ("commit the pickup, then
repair"). The 2026-08-04 recurrence is the same clause failing one step
earlier and for real: the `recovery-discard` build's second review ran in Codex
under the top binding, returned eight findings including two criticals, and
never wrote the file at all, so the findings reached the next session only as a
hand-pasted summary and four of the eight (C4, C5, C6, SP5) now exist as tags
with no text. The repair round rebuilt its evidence from the tree, so no repair
rests on the lost text, but nobody can say the repairs close the findings that
were raised. Make writing `reviews/<spec-slug>.md` a required step of
`/bench-review-implementation` rather than an artifact `/bench-implement-spec`
merely reads if present, and say so in both commands — including for a review
delegated to another harness, where the coordinator owns capturing the returned
findings to that path before acting on them. The spec-build review receipt is
not that record: it attests that a review happened and its verdict, and a
repair round consumes the findings. Seventh (drained the same day), `/bench-implement-spec`'s "Route the
venue" is unsatisfiable under a harness that forbids delegation. It requires
every spec-backed run to assign genuine write work to a write subagent before
the first implementation edit, with no inline threshold; a session carrying a
reviewer-set standing instruction against the Agent tool cannot comply, and the
FT148 session took an unsanctioned inline route for four prose edits to an
already-implemented spec. `craft-delegate`'s capability-aware policy covers a
harness that *cannot* spawn one and says nothing about one that *may not*. Give
the clause a stated precedence — name the fallback when delegation is
unavailable for either reason — and decide alongside it whether a spec-doc-only
correction (no code, no seams) is exempt outright, since routing prose edits
through an isolated worktree costs more than it catches. Like the fourth, this
clause reaches past always-loaded prose into a skill. Eighth (drained
2026-07-28 from `capture/IDEAS.md`), shell wait-loop hygiene in `AGENTS.md`'s shell
conventions: a wait-loop whose predicate matches its own command line never
terminates — `until ! pgrep -f "codex exec"; do sleep 20; done` matches the
loop's own bash process, and one such loop left an orphaned background process
through the FT152 build. One line: wait on a PID or a sentinel file, never on
a `pgrep` pattern the waiting command itself contains. Deliberately
not batched here: FT130 (mid-gate
`bench idea`) stays its own row because its preferred fix is mechanical, in
the CLI rather than in this prose. Ninth (drained 2026-07-30 from the Claude
usage-report assessment), wide edits to structured prose need an
artifact-level consistency pass. The observed failure was not the use of
item-wise edits itself; it was reporting completion without re-reading the
whole spec after those edits introduced cross-story contradictions. Put one
completion criterion in the authoring discipline for specs, roadmaps, and
decision maps: re-read the complete affected artifact and reconcile every
repeated field or cross-reference. Do not prescribe one editing mechanism —
the document's actual structure decides whether a batch or item-wise patch is
safer. Tenth (same source), an ad-hoc script that can delete, prune, release,
or otherwise discard state defaults to a plan: resolve and print the exact
targets and commands, show parsed fields for a small sample, and mutate only
on an explicit apply. This is one shell-conventions clause, not a new skill;
Bench-owned commands keep their stronger typed plan/apply and recovery
contracts. Eleventh (same source), phase exits state every material acceptance
shortfall and unverified tier with actual-versus-target evidence where a
target exists. Fold it into the existing structured exit rules: omit empty
groups rather than emitting ceremonial `none` sections. Twelfth (drained from
the FT131 implementation retro), review handoffs report both raw axis findings
and de-duplicated repair targets; eleven findings across axes and ten unique
defects are different planning facts, and reporting only one makes the repair
charge ambiguous. Thirteenth (same source), a gate-bootstrap review explicitly
asks what establishes trust in the verifier itself, not only whether the
selected binary exposes a verification command; that question found FT131's
gate self-attestation defect after ordinary surface checks passed. Background:
`docs/reporesident-distillation.md` §3 and §6; the 2026-07-30 Claude
usage-report assessment and the FT131 implementation retro, drained here.
Fourteenth (from the decision-map integrity retro), close the review gap between
a compiled decision map and its narrower acceptance rows. A fresh-session
map-to-spec dogfood run happens before semantic review, and the main checkout's
ignored development binary is rebuilt first so stale local behavior is not
reported as shipped behavior. Review charges treat defaulted decisions as
authoritative unless the spec explicitly overrides them; a claimed repair is
checked against both the applicable coverage row and the defaulted-decision
table. When diagnostic wording and its exact canary expectation cross a ticket
fence, the review may require one explicitly justified atomic repair rather
than leaving the oracle and its bite inconsistent. Source: the decision-map
integrity implementation retro, drained here. Fifteenth (drained 2026-08-01
from the learnings journal), repository-wide sweeps include hidden paths.
`rg` skips dot-directories by default, so a 278-reference migration sweep
updated the canary fixtures' literal `dot-bench/` copies and missed the
canonical `.bench/` and `.agents/` files they shadow — the exact desync the
conformance anchors exist to catch, caught at the cost of one red gate cycle.
Extend `AGENTS.md`'s shell conventions: a repository-wide sweep or audit
passes `--hidden` (excluding `.git/**`); and weigh the mechanical alternative
— a conformance check that each `dot-*` fixture still matches the file it
shadows — which would catch the class regardless of tooling and remove the
instruction rather than duplicate it. Sixteenth (same drain), bench-verb
discovery without probing. A bare `bench` verb can be a porcelain that acts,
not a parser that refuses — the bare `bench worktree` subshell hung an
automation call for two minutes and leaked worktrees when signal-killed.
Extend the same conventions: discover a subcommand's shape from
`bench commands --brief` or `bin/bench.sh`, never by running an unrecognized
bare verb; and redirect stdin from `/dev/null` for any `bench` invocation
whose interactivity is not already known. The defect half — the bare verb's
default itself — is FT178's, not prose. Seventeenth (drained 2026-08-04 from
the learnings journal), a writing session opens its worktree at entry, not
after the collision. Invariant 1 says to isolate when `git status` shows
another writer, but that reading is a single check at session start: a
`/bench-debug` session found the tree clean, edited five files in the main
checkout, and was then deadlocked against a second session's drain that dirtied
the same tree mid-run — neither changeset could commit past the other, moving
to a worktree afterwards did not help, and the reviewer had to hand-run the
`git restore` that `block-dangerous-git` correctly refuses. `/bench-debug`
routes code authorship through `craft-delegate` only at Phase 5, while Phases 1
and 2 routinely write repro harnesses with no isolation guidance at all. Either
the entry orientation names the worktree as the debug session's workspace or
the isolation rule moves ahead of the first write; the reviewer decides whether
that is `/bench-debug`'s clause or a general session-entry one.

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
the double that was written to satisfy it. `internal/specbuild` alone takes six
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

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.** Apply the
gate's "prove it bites" standard to prose: audit the craft-skill library and
the communication protocol so each skill and always-loaded clause cites an
observed failure it prevents (from the learnings journal or session evidence),
merge overlapping craft docs, and shrink the always-loaded `BENCH.md` rules to
demonstrated-delta clauses. Distinct from FT89, which fixes guidance
*correctness*; this row cuts guidance *weight*. FT187 takes the
"How to talk to me" slice ahead of this row, so what remains here is the
craft-skill library and the demonstrated-delta audit over the rest of the
always-loaded surface. Kit edit under the
`craft-synthesis` discipline; starts as a grill (`/bench-shape-idea`) because
the cut line on always-loaded rules is a reviewer decision.

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
loop.** A kit edit that instructs spending a model tier can contradict the
escalation policy without any loop catching it: the widened write-spec step-9
triggers shipped an automatic top-tier spawn past review (observed
2026-07-22; corrected in the mid-tier rerouting commit). Make
`craft-synthesis`'s consistency loop name the escalation policy as a standing
cross-check for any kit edit that spends a tier. Kit edit under the
`craft-synthesis` discipline.

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
the latter fail with `no Status: staged line`; `bench spec build promote` is a
third, and the phase contract names it the sole author for a reviewed spec
build. Name one owner per landing route and make the others refuse rather than
race. Source: the FT128 implementation retro, drained here.

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
refuse is the row's real work — which is why this row stays out of FT107's
prose batch. The second face, 2026-08-01: the spec-build lifecycle's
clean-checkout precondition refuses on a dirty `capture/IDEAS.md`, so
mid-build parking — which `.bench/BENCH.md`'s Capture section instructs —
blocked three verified tickets behind a one-line capture edit, and recording
the learning reproduced the block; the session's improvised
scratchpad-staging workaround loses the idea if the session dies first, the
exact failure `bench idea` exists to prevent. The reduced-gate allowlist
(`capture/`, `specs/`, `ROADMAP.md`) already names the paths whose dirtiness
cannot invalidate a checkpoint's evidence, because the checkpoint's subject
is the assignment worktree, not the main checkout — exempting that allowlist
from the precondition is the smaller fix and composes with the shipped
reduced-phase-set mechanism; decide it together with the gate-window face so
both surfaces share one answer. Kit edit under the `craft-synthesis`
discipline. Sources: the 2026-07-25 learnings entry, verdicted in a prior
drain; `capture/learnings.md` 2026-08-01, verdicted here. The second face
recurred across the FT176 close: `bench idea` dirtied the capture file and
blocked the next path-scoped commit until folded into the landing. Source: the
spec-build-lifecycle-preconditions retro, drained here.

The gate-window face recurred on 2026-08-02 — a `bench idea` write landed inside
a `bench gate --fresh` window, every phase came back green, and the verdict
rejected itself for a changed subject — which is the second full re-run this row
has now cost. A third face joins from FT164's build: the same write against an
*active spec-build run* is worse than a wasted gate. A capture commit before the
run's first checkpoint forced recomposition at zero checkpoints, unrecoverable
at the time (`2874d94` has since made it a rebase), and the run was abandoned and
rebuilt from snapshotted delegate diffs; even with that fix, every tip move
mid-run buys a fresh full gate for recomposition. Two things follow. The tree is
frozen from `bench spec build start` to `promote` — no capture commits, no
roadmap edits — and the `--full` route's phase-boundary handoff write, which
today lands between them, belongs before `start` or after `promote`; that
placement is a `bench-implement-spec` command-text edit and rides whichever
answer the queue-versus-refuse decision produces. Sources:
`capture/learnings.md` 2026-08-02 and 2026-08-03, verdicted here; the
`ft164-ticket-contracts` retro, drained here.

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

The check-level-conformance-scoping close adds a slicing rule to that visit: a
ticket spanning runtime projection, canary selection, and documentation
conformance splits those owners before assignment. Its combined scope/report
slice consumed most of 14 assignment attempts even though the narrower
outer-selector, conformance-meta, and evidence-retention slices converged
directly. Source: the check-level-conformance-scoping retro, drained here.

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

**FT184 (MEDIUM) — `bench spec build checkpoint` receipts are hand-assembled
against an undocumented schema.** The first full lifecycle run discovered the
row-outcome vocabulary (`passed|already-covered|not-tdd-able`) only by reading
`receiptRows`, and `error: invalid spec build receipt` names no failing
condition, so each invalid-receipt refusal cost an in-package debug harness.
The row's own re-price trigger fired on 2026-08-04: landing four repair tickets
through `bench spec build` cost more coordinator turns assembling two receipts
than verifying the two critical repairs they attested. Every field the
coordinator hand-writes is derived data — run and assignment identities, the
assignment base, a tree hash computed the way `git.TreeHash` does it, the
ticket digest, one row per acceptance ID, the changed-path set diffed against
the base, and a sha256 per assumption string — each cross-checked by
`validateReceipt`, so any one wrong is an opaque `invalid spec build receipt`.
The coordinator should supply only what it alone knows (row outcomes, the check
list, and the probe command, output, and exit) and the tool should derive the
rest, which is this row's generator half stated as a contract.
Two halves, one owner: a receipt generator (a `bench spec build receipt
<assignment>`-shaped command) that derives the row set from the ticket file —
removing the row-set-mismatch refusals paid twice in that build — and a refusal
that names the first failed check. Sources: the per-component-gate-scoping
retro, drained here; `capture/learnings.md` 2026-08-04, verdicted here — the
third build to pay the class, and the one that moved the price from LOW.

The 2026-08-03 run re-priced the generator half and added two faces. The
coordinator hand-built a `make-receipt.sh` and spent roughly 30k tokens
reverse-engineering the checkpoint and review receipt schemas, so the generator
emits the tree hash, digests, and row IDs pre-filled for a named assignment. It
also enumerates the disposition vocabulary, which is the second face: a review
receipt marking closed findings `accepted` reads to `hasAcceptedFinding` as an
unrepaired defect, `promote` refused, and clearing the recorded review cost a
recomposition round — `accepted` means a real defect the review accepts
unrepaired, and closed findings need a different word. Either the generator
enumerates the vocabulary or the schema validates the field and refuses the wrong
word at submission. Third face, same owner: `bench spec build assign --ticket`
refuses a repo-relative or absolute path with `spec build ticket must name one
regular ticket file`, which reads as "the file is missing" when the file is
present and only the *form* is wrong; it cost two failed invocations before the
parser was read. A refusal naming the expected form would have cost none.
Sources: `capture/IDEAS.md` 2026-08-03 and both 2026-08-03 retros, drained here;
`capture/learnings.md` 2026-08-03, verdicted here. The gate-evaluation-snapshot
close repeated the generator request after its coordinator manually reproduced
the tree-hash and receipt schemas. Source: the gate-evaluation-snapshot retro,
drained here.

**FT166 (LOW) — `bench capture commit`: porcelain for the ambient capture
set.** Commit the capture surfaces (`capture/learnings.md`, `capture/IDEAS.md`,
`capture/session-handoff.md`, `capture/retros/`) with a conventional message under the
doc-only standing rule, so the plain-`git` step every session hand-assembles —
with FT107's empty-index hazard attached — becomes one sanctioned command.
Weigh it beside FT107's third clause, which owns the hazard prose; the
porcelain would remove the instruction rather than duplicate it.

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
answered: the convention is live but unenforced, and FT107's sixth clause makes
writing the file a required review step, so what stays here is the signal —
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
One residual has no owner at all: when a fence does turn out to be wrong mid-run,
`abandon` is whole-run and no verb releases a single open assignment to re-fence
it, so the correction costs a sibling ticket and another delegate cycle. The
shipped authoring-time refusal makes that rarer without making it reachable.
Source: `capture/learnings.md` 2026-08-03 and 2026-08-04, verdicted here.

**FT153 (MEDIUM) — the canary's vacuity baseline is a collision screen, not a
vacuity proof.** A behavior-owned fixture's EXPECT is compared against its
contract group's empty-tree baseline, which establishes only that the string is
not infrastructure noise. A generic banner that any failure prints passes that
screen forever, so nothing grades mutation-specificity for the 33 fixtures
carrying the gate-guarding-the-gate proof — the profile now says so, but saying
so is not checking it. Real vacuity wants the unmutated twin: BASE plus `files/`
without `MUTATE.json`, run in the same shape and required to *not* show the
EXPECT. That is derivable today for the 9 MUTATE-shaped fixtures, while the 24
`files/`-only ones bake the mutation into the fixture tree and would need a
delete-or-absence op `MUTATE.json` cannot express — so the shape is genuinely
open and the row starts as a grill.

Second clause, same owner and same question: whether a non-contract vacuity
baseline needs a full inner gate. The stage-2 build silently narrowed every one
of them to a single phase, which is semantically wrong and was reverted before
landing, but the revert cost ~6 s of gate and ~1 s of canary. A deliberate
scoping with correct semantics may recover that legitimately. Decide it with the
twin question rather than separately — both rule on what a baseline must run to
mean what it claims. Entry: `/bench-shape-idea`. Sources: `capture/IDEAS.md`, drained
here; the `ft91-canary-compiled-bites` review S1, recoverable via `git show
4429b05:reviews/ft91-canary-compiled-bites.md`.

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

## Reds the diff doesn't own — inheritance, load, and harness defects

Four rows, one failure class: a red that answers for something other than the
diff in front of the gate — an inherited baseline, machine contention, a
literal deadline, a harness defect.

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
before.

**FT24 (parked pending upstream) — Codex agent-line guard parity.** Researched
2026-07-11: still not implementable on current Codex — delegation has no
matchable `tool_name` on a deny-capable hook event, and `SubagentStart` cannot
stop the subagent. The current surface verdict is canonical in
`.bench/BENCH-reference.md` Hook Layers. Graduate only when the Codex changelog
adds a spawn tool name or a deny-capable SubagentStart.

**FT8 (scheduled, not actionable) — Sonnet 5 mid-tier revisit.** Time-boxed to
2026-09-01 or the next frontier shift.

**FT38 (tabled, revisit on or after 2026-08-09) — dashboard visual identity
pass.** `bench dashboard` v1 shipped data-faithful and visually neutral; the
original idea wanted a rich treatment with animated characters, reference
saved at `ui_example/` (Gather-style pixel office with activity feed).
Reviewer tabled it 2026-07-09 for at least a month. When it revives, the work
starts as a grill (`/bench-shape-idea`); decision detail recoverable via
`bench spec history dashboard`.

**FT170 (scheduled, not actionable; revisit 2026-08-12) — behavioral
red/green evaluation for skill guidance.** Before Bench adopts a prose-heavy
skill-testing workflow, prove the need through a benchmark substrate: choose
one narrow behavior with deterministic artifact assertions or a blinded
scoring rubric; run repeated no-guidance and candidate-guidance trials in fresh
isolated contexts with pinned model and effort; keep authoring cases separate
from held-out evaluation; and report variance. Run the benchmark advisory-only
during skill changes or assessment, never in the deterministic per-commit gate.
Only stable improvement earns the smallest necessary `craft-skills`
requirement and harness pointer. Source: `capture/IDEAS.md`, drained here.

## Dependencies

The dependent FT is named first. Only the literal table blocks work; the
recommended table is sequencing advice.

### Literal

| FT | Depends on | Why |
|---|---|---|
| FT107 | FT141 | Its fix-loop shrink rule cannot distinguish convergence until inherited reds are attributed to their pinned baseline. |
| FT107 | FT187 | Both edit `.bench/BENCH.md`; FT107 is scoped as one always-loaded-prose diff, so the restructure lands first and the batch's clauses apply to the reduced surface. |
| FT186 | FT108 | The gate restructure needs the mechanical exit test before moving oracle code without behavioral stories. |

### Recommended

| FT | Better specified after | Why |
|---|---|---|
| FT71 | FT169 | The event schema should record the settled landing and recovery lifecycle rather than version an interim one. |
| FT100 | FT89 | Cut prose after the correctness and coherence pass establishes which guidance is still authoritative. |
| FT107 | FT158 | The cross-harness refute pass is the only demonstrated falsification for a wide always-loaded prose diff; make it standing before the batch lands. |
| FT107 | FT156 | No fixture proves a section-scoped `.bench/BENCH.md` anchor — the exact surface the batch edits; rule the anchor mechanism first. |
| FT187 | FT156 | Same untested anchor surface, and its fifth clause rewords the paragraph FT156's named-lighter-path ruling may replace outright. |
| FT108 | FT89 | FT89 single-sources the skills index the new skill must join; the expand–migrate–contract and gate-cadence rules it builds on are already settled in `craft-tickets`. |
| FT111 | FT179 | FT179 reopens FT111's edit-in-place-only ruling on an order-larger measured count; land them as one `craft-comments`/`craft-review` visit. |
| FT172 | FT106 | Reuse the document-claim probe for semantic roadmap claims instead of designing a second checker. |
| FT162 | FT169 | Build full-run subject resolution on the settled landing primitive. |
| FT166 | FT98, FT113 | The porcelain composes over the shipped reduced-gate path allowlist; recoverable set-aside then defines the commit command's smallest sound contract. |
| FT168 | FT153 | Expose focused canary execution after baseline meaning is settled. |
| FT169 | FT98 | Reuse recoverable discard in the landing contract; label resolution is already available. |
| FT169 | FT188 | The foreign-dirty face is the substrate's consumer; specify the landing command once exact prospective landing has settled attribution. |
| FT175 | FT173 | The ledger's read surface is AXI; settle one derivation per principle before adding a consumer that needs all ten. |

### Goal tracks: guidance prose and the claim ledger

Two reviewer goals share one path: implement the guidance-prose backlog, and
ship `bench cite` (FT175) on an AXI-compliant CLI (FT173). The process
precursor for both is landed: FT164's ticket-contract core shipped 2026-08-03,
so every later build slices its tickets through a `craft-tickets` the specbuild
parser accepts. The payoff
facts shaping the order, verified in-tree 2026-08-02 with an independent
mid-tier refutation pass: `.agents/` and `.bench/BENCH.md` sit outside the
gate's reduced scope, so every separately-landed prose diff pays a full gate
— rows batch on the shared full gate, not just shared files; anchor-pinned
files couple prose diffs to conformance fixture updates (`craft-delegate` 14
anchors, `bench-implement-spec.md` 35+, `.bench/BENCH.md` 17); and prose
compounds per-session while the AXI foundation compounds per-CLI-change, so
the prose batch outranks the foundation only if it lands whole and
falsifiable — otherwise the foundation goes first while FT141 and FT158
build.

1. Take FT156's anchor-mechanism ruling as the grill at `/bench-write-spec`
   entry (reviewer direction 2026-08-03) and FT144's one-decision-both-phases
   call as the remaining shaping item. Reviewer latency is the
   binding constraint: grills serialize on the reviewer while builds
   parallelize on agents. Four of the original six items have closed —
   FT173's principle-9 relaxation and FT175's spec-start gate were ruled
   2026-08-02, FT175's three ledger decisions moved behind the owners
   they consume (step 5), FT164's four flagged spec calls closed when
   its spec was written and implemented, and FT181 shipped 2026-08-03.
2. Drain the staged frontier — FT135, FT187, and FT188 all carry staged specs
   — before authoring any new spec; deferring one is an explicit reviewer
   override, never a silent skip. Of the three only FT187 serves a goal
   directly, and FT188 unblocks the parallel capacity the rest of the runway
   assumes. FT141
   builds in parallel where
   capacity allows: it is Go, prose-independent, and it unblocks FT107
   *whole* — splitting the fix-loop clause out would spend a second spec,
   review, and full gate on the same anchor-pinned surface, so the batch
   waits for FT141 instead.
3. Land the prose track's falsification before its wide batch: FT158 first
   (with FT156 unruled and FT170 unproven, the cross-harness refute pass is
   the only demonstrated defect-finder for prose diffs), plus FT156's
   mechanism build if its grill rules for one. Then prose batch 1: FT107
   whole + FT89 + FT99 + FT144, with FT102, FT112, and FT165 riding on the
   shared gate — disjoint files make them safer to batch, and FT165 early
   improves every later grill. Then FT179 + FT111 as one visit, with FT164's
   repair-lane residual riding it as a second `craft-delegate` visit, then
   FT108.
4. Implement FT173 in three independently reviewed specs, foundation first:
   AXI principles 8–10 into `craft-cli`, AXI kept scoped to query surfaces,
   truncation consolidated onto `sanitize.Preview` and aggregate mechanics
   extracted only where two real derivations exist, all without changing
   emitted bytes — medium, guarded by the pinned contract suite. Contextual
   disclosure (`help[]`) second: large and byte-changing (zero prior art; it
   migrates the AXI query surfaces in green batches and rewrites their
   pinned output contracts), shippable under the 2026-08-02 principle-9
   relaxation granted to this face alone. Git-inspection last or in
   parallel: it extends the existing `bench diff` owner into one coherent
   status-and-diff snapshot (untracked paths, pre-computed counts, drift-safe
   reads), but `bench diff` is not a ledger surface — landing it before
   `help[]` would rewrite the pinned diff contracts twice. FT173 remains open until all
   three land.
5. Shape FT175's decision map once the foundation and `help[]` land, and
   write the FT175 spec once the owners it consumes — truncation,
   aggregates, `help[]` — are settled, then build it as vertical green
   slices: file evidence plus strict store and `claim show/check`;
   assessments with spans and absence scopes; command evidence plus replay
   and staleness; supersede/retire reachability; the
   complete AXI list/detail/status surface using FT173's owners;
   deterministic gate and ambient status integration; then one local
   contradiction as the dogfood proof. Reviewer ruling 2026-08-02:
   "interfaces settled" means the consumed owners, not the whole FT173 row —
   the foundation and `help[]` settle them, so the FT175 spec starts one spec
   earlier and git-inspection stays off the critical path. Its three ledger
   decisions are deferred until those two implementations land, so the FT175
   shape-idea session moves to then rather than joining the front-loaded one.
6. FT100 grills and builds last, after FT89 establishes which guidance is
   authoritative and past FT170's 2026-08-12 revisit.

Fold FT106's verified-document vocabulary and FT162's exact subject binding
into the FT175 spec instead of building them as prerequisites; FT99 rides
prose batch 1, and its uncertainty obligation folds into the FT175 spec
where claims consume it. FT133 remains parallel evidence hardening; FT71
stays deferred behind its existing FT169 recommendation. FT172 is outside
this critical path; FT156 joins the decision session because no fixture
proves a section-scoped `.bench/BENCH.md` anchor — the exact surface the
prose batch edits.

## Recommended sequence

1. `/bench-write-spec` — FT194, the reproduced project-green desynchronization that blocks every spec-build promotion and leaves abandon as the only escape. Land that repair before opening another lifecycle run.
2. `/bench-write-spec` — FT156, taking the anchor-mechanism ruling as the grill at spec entry (reviewer direction 2026-08-03). Once FT194 restores promotion, reviewer latency on this ruling is again the binding constraint for both goal tracks.
3. `/bench-implement-spec` — drain the staged frontier: FT188 (`specs/exact-prospective-landing/spec.md`) first to remove the writer lock, FT195 (`specs/go-build-cache-footprint/spec.md`) next to reduce every later build's cache and publication cost, then FT187 (`specs/ft187-communication-surface-cut/spec.md`) and FT135 (`specs/pre-push-guard-visibility/spec.md`).
