# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
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

**FT71 (HIGH on the bank track, evidence supplied) — versioned local shift
evidence.** Emit a redacted, append-only local event schema for shift/session
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
`/bench-shape-idea`. Sources: `IDEAS.md`, drained here;
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
as a consumer of that seam rather than a separate build.
Entry: `/bench-shape-idea`. Source: `IDEAS.md`,
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
diagnostic.

The handoff storage question joins the same lifecycle decision. A single
repository-level `session-handoff.md` can be clobbered by concurrent workstreams;
per-spec handoffs would isolate those pins and retire with their spec, but would
leave non-spec work needing a repository-level owner. Decide whether the
authoritative handoff remains singular, moves into `specs/<slug>/` for spec-backed
work, or becomes a generated projection over per-workstream state rather than
adding a second handoff convention by accident. Source: `IDEAS.md`, drained here.

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
generic dirty state. The handoff's `## Next command` body is then graded as
exactly one backticked harness-native invocation, so the authoritative next
action cannot drift into explanatory prose. Kit edit under the
`craft-synthesis` discipline. Sources: the craft-tickets, light-path,
artifact-suite, and artifact-hoist retros, drained here; `IDEAS.md`, drained
here and in prior runs; `.bench/learnings.md`, verdicted here.

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
vacuity-baseline finding became FT153. Source: `IDEAS.md`, drained here.

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
`craft-synthesis`. Sources: `.bench/learnings.md`, drained in a prior run and
here; `session-handoff.md`, drained in a prior run.

**FT158 (MEDIUM, evidence supplied three times) — make cross-harness
falsification standing for kit-guidance diffs.** FT91's draft, FT152's build,
and the FT123 + FT124 build each received a Codex pass charged to refute rather
than grade after the ordinary review surface had cleared them, and each pass
found a real defect. The third run meets this row's graduation trigger: its
counterexample was an exact worktree label that looked like help, a flag, or
`--`. Make the pass standing for kit-guidance diffs, where a defect compounds
through every session that loads the prose, and give it the prepared spec and
diff bundle so it does not spend the charge re-reading unrelated contracts.
Every finding still takes an explicit accept, merge, or dismiss disposition;
the pass is advisory and does not become a second oracle. Kit edit under the
`craft-synthesis` discipline. Sources: `.bench/learnings.md`, verdicted in a
prior drain; the gate-fastpath and FT123 + FT124 retros, drained here and in a
prior run.

**FT128 (MEDIUM, evidence supplied) — the agent-line guard cannot see a fork's
real model.** `check-agent-line` decides from the delegation envelope's
`resolvedModel`/`model` field alone (`internal/lines/lines.go`), so a
fork-type delegation — which inherits the parent's model and ignores any
`model` override — passes the guard on a declared cheap alias while actually
running the top tier. That is exactly the silent escalation invariant 2
forbids and the guard exists to block, and it is the one delegation shape the
guard grades backwards. Verified against the tree 2026-07-25: nothing in the
guard reads the subagent type. The fix is to make the envelope's delegation
type part of the verdict — deny a fork whose declared alias is not the
session's own tier — which first needs the type's field name pinned from a
real envelope rather than assumed. Enforcement stays exact-string with no
provider lookup.

FT97 merges in here, 2026-07-29 — same enforcement surface, one visit. The
deny message single-sources its bound-tiers listing, which leads with the
three tier ids and trails the harness aliases; inside a Claude Code session
the aliases are the only tokens the Agent tool can pass, so the error leads
with ids that harness cannot use (observed 2026-07-19). The design is already
decided: the closed `specs/ft128-agent-line-binding/decisions/multi-harness-line-binding.md` map answers the
schema question — symmetric per-harness bindings with no canonical family,
each layer reporting its own harness's tokens. One build fixes the fork
verdict and re-leads the denial from that map in the same diff.

The static half belongs in that build too: a conformance check rejects any
`.agents/commands/*.md` tier-model token against `.bench/lines.env`, so command
prose cannot reintroduce a hardcoded binding. The check must demonstrate its
bite for the right reason. Spec: `specs/ft128-agent-line-binding/spec.md`.
Source: `IDEAS.md`, drained here.

**FT135 (MEDIUM) — a pre-push guard on a guessed branch looks armed while
protecting nothing.** When the repository has no resolvable default branch,
the installed pre-push hook falls back to a baked `fallbackProtectedBranch`
of `main` (`internal/adopt/link_hook.go`). The hook re-resolves `origin/HEAD`
live on every push and reaches the baked token only when that lookup stays
empty — but in exactly that repository the guard looks armed while defending
a branch that may not be the default, FT86's own failure class one layer
down. `bench doctor` (or `bench link`'s output) should report whether the
installed guard's protected branch was resolved or guessed, so the false
armor is visible where the reviewer looks. Source: `IDEAS.md`, drained here.

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
enforcement hook. Source: `IDEAS.md`, drained here.

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
here. Source: `IDEAS.md`, drained here.

**FT98 (MEDIUM, evidence supplied four times) — one preserve-then-discard primitive;
four faces.** Three rows were faces of one missing primitive — a sanctioned,
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
The containment primitive now exists — `LandedInDefault` proves patch-id
containment for landed-branch pruning — so either route recovery-payload
landed-proof through it or add the explicit reviewer-authorized discard;
fail-closed stays the default and the cut line is a reviewer decision. This
face now owns the whole visible residue: FT148's orphan retirement shipped
2026-07-27 and its ledger compaction correctly declines these rows because
they preserve work, so 27 recovery refs across 23 preserved worktrees (24
assignment rows open, one active, at the 2026-07-29 session start) are
re-preserved at every session start and nothing else will retire them
(`bench worktree recovery <ref>` returns `retain … unlanded` on a sampled
ref, 2026-07-27; ref count re-taken 2026-07-29) — the residue is growing,
not draining. The third occurrence came when a scoped roadmap commit was
blocked by an unrelated dirty session handoff on 2026-07-30; the session used
an isolated verification worktree, the landing workaround owned by FT169,
because the sanctioned set-aside primitive still does not exist. Face
two, `bench commit`'s set-aside (was FT127): the refusal reads "working-tree
files outside the named set block the commit — name them, or set them aside",
but no set-aside route exists in the CLI, so an agent's only real exits are
committing an unrelated file into a scoped commit or reaching for
`block-dangerous-git`-blocked plain git; build the route on this same
primitive rather than rewording the advice — the need is real and recurring.
Face three, mutation-probe revert (was FT114): deliberately weakening an
implementation to prove a check bites always needs a revert, and
`block-dangerous-git` blocks `git checkout <path>`; copy-aside works but is a
papercut on a first-class activity in this repo (cf. `tests/canary/`), and a
scoped single-path revert through the same recoverable primitive replaces
both the papercut and any guard exemption. Whichever face ships first defines
the one discard semantics; the others reuse it.

The FT131 close adds the ignored-cache face. Current-source verification defaults
`GOCACHE` under `dist/`, while nested route tests strip ambient cache variables;
an isolated full gate can therefore leave more than 40 MiB and 1,000 ignored
entries that `bench worktree clean --discard-ignored --full` refuses even with
the matching fingerprint. Prefer a lifecycle-owned scratch cache or cleanup at
the fixture owner; if residue remains legitimate, let the same size-bounded,
fingerprinted discard contract summarize and authorize the generated directory
without falling back to manual deletion. Source: the FT131 implementation retro,
drained here.

**FT169 (MEDIUM, evidence supplied) — one sanctioned worktree landing command
owns the stale-base dance.** The gate-fastpath build hand-ran the same sequence
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
is why they belong in it rather than in the phase prose.

That parked half now has a demonstrated cost, and a cheaper option than the
probe it was parked behind. The 2026-07-27 drain promoted an `IDEAS.md` capture
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
out of the tree first. Source of this clause: `.bench/learnings.md`, verdicted
here. The grammar face came from `IDEAS.md`, drained here.

**FT126 (MEDIUM, evidence supplied) — recurrence tallying belongs to the
roadmap parser and context snapshot.**

The same parser and snapshot seam owns recurrence tallying. An idea, learning,
or retro should be able to cite an existing FT as its primary owner so the drain
records a new occurrence on that row instead of manufacturing a duplicate;
the current count stays visible while Git owns the event history. The first
captured occurrence is FT98's 2026-07-30 scoped-commit refusal, with FT169 as
the downstream workaround. Decide the citation grammar and malformed-reference
posture alongside the row grammar rather than adding a second roadmap parser.
Source: `IDEAS.md`, drained here.

Spec: `specs/ft126-recurrence-tallying/spec.md`.

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

Sources: `RR:S-06`, `RR:S-07`, `RR:S-08`, `RR:S-10`, `RR:S-11`, `RR:S-12`,
`RR:S-13`, `RR:S-14`, `RR:S-15`, `RR:S-16`, `RR:S-17`, `RR:S-18`; `RC:M-05`.

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

**FT107 (MEDIUM) — the standing guidance rules, batched: one
always-loaded-prose diff.** Fourteen remaining clauses edit the same
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
repair"). Seventh (drained the same day), `/bench-implement-spec`'s "Route the
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
2026-07-28 from `IDEAS.md`), shell wait-loop hygiene in `AGENTS.md`'s shell
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
integrity implementation retro, drained here.

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
their assigned posture until semantic review caught them. Source: `IDEAS.md`,
drained here.

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.** Apply the
gate's "prove it bites" standard to prose: audit the craft-skill library and
the communication protocol so each skill and always-loaded clause cites an
observed failure it prevents (from the learnings journal or session evidence),
merge overlapping craft docs, and shrink the always-loaded `BENCH.md` rules to
demonstrated-delta clauses. Distinct from FT89, which fixes guidance
*correctness*; this row cuts guidance *weight*. Kit edit under the
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

**FT113 (LOW) — a `bench commit --spec` status flip makes its own green verdict
stale.** The flip from `Status: staged` to
`implemented` is necessarily written after the gate has passed, and
`specs/*.md` is not on the capture-only allowlist, so every clean `--spec`
landing leaves the gate reporting strong-stale and sends the next session to
a re-run that finds nothing — reproduced through the accused command
2026-07-24 (`git diff-tree -r --name-only 20f0767 0faf47f`: only `IDEAS.md`
and the flipped spec file drifted). Its fix is a reviewer decision between
three options that are not equally safe: widening the allowlist to
`specs/*.md` (admits arbitrary spec edits, which the fixed exact-path
allowlist deliberately refuses); teaching the gate cache to record the
post-flip tree when `--spec` itself performs the flip (narrower — the flip is
the tool's own write and its content is known); or accepting the face as
cosmetic and documenting it. The same command has a separate avoidable usage
failure: `bench commit --spec <slug>` edits the owned spec transition itself but
still requires another named path. Count that owned transition as satisfying
the path requirement without widening the staged set. Immediate retirement on
the default branch currently follows that landing with a second full oracle run
over another short-lived tree; settle the cache handoff and retirement sequence
so the two owned transitions do not require redundant full runs. Sources: the
FT131 and decision-map integrity implementation retros, drained here.

**FT130 (LOW) — parking an idea mid-gate silently voids the run.** During
FT122's gated commit a session answered a reviewer question and ran `bench
idea` to park the tangent, which wrote `IDEAS.md` inside the gate's window;
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
prose batch. Kit edit under the `craft-synthesis` discipline. Source: the
2026-07-25 learnings entry, verdicted in this drain.

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
`/bench-shape-idea`. Source: `IDEAS.md`, drained here.

**FT164 (LOW) — ticket and repair charges preserve blast radius, postures, and
gate cadence.** Two
rule-shaped journal entries from FT154's own build, one owner file. First,
the wide-refactor branch: the first folder-layout ticket grouped a deep-unit
change with every consumer family, and the delegate exhausted its fresh
context before the fixture migration closed, so the ticket could not land
independently green — make blast-radius classification the explicit first
branch of the breakdown procedure, ahead of vertical drafting: a wide
refactor takes expand–migrate–contract sequencing, migration tickets sized by
an ownership fence (package or consumer family), the contract ticket blocked
by every migration. Second, the cadence rule: tickets carried a `project gate
is green` checkbox, so the coordinator ran the full gate, checked the box,
and then `bench commit` ran the same gate again — keep gate state out of the
ticket template (the green landing commit is its one source) and define the
cadence as focused seam checks while iterating, one full gate at the atomic
commit, one composed gate at final check. A third clause, from the light-path
retro: the independently-green rule holds, but name its proportionality
ceiling — FT139's one-line test-helper change still paid a worktree, a fresh
delegate, a focused proof, and a full gate; a one-line test-harness ticket is
the light path's ceiling, noted beside the rule so paying that cost is a
choice rather than a default. Kit edit under the `craft-synthesis`
discipline. Sources: `.bench/learnings.md`, verdicted here; the light-path
retro, drained here.

The charge itself carries two more load-bearing facts. When a delegate adds a
member to an enumerated family, the charge names every registry found by
tracing an existing sibling across the package; FT152's canary work missed the
classification registry and paid a repair round. When a repair changes a
shared environment helper, the charge names which rows opt into shared caches
and which remain hermetic, then runs the focused failing rows before the full
gate; the artifact-suite repair over-stripped that distinction. Sources: the
learnings journal and the `ft91-artifact-suite` retro, drained here.

Model-comparison charges are another instance of the same owner. Every
candidate receives one constant charge, base commit, file fence, effort,
focused suite, and independent probe; setup failures are recorded separately
from generation time. Fixed behavioral checks and the independent probe decide
whether the cheap default clears its bar, while style differences decide which
acceptable patch is better rather than justifying an expensive tier by
themselves. The existing `craft-line` cheap implementation default and
red-driven escalation ladder remain authoritative. Source: the FT129
implementation retro, drained here.

Repair charges that update a canonical command inventory also name both
currency owners in their focused proof: the documentation inventory and the
wrapper/router registry. FT131's repair updated one and missed the other two
until the atomic gate; tracing those owners in the charge would have caught the
integration misses without replacing the gate as oracle. Source: the FT131
implementation retro, drained here.

**FT165 (LOW) — fold the domain-modeling discipline into
`/bench-shape-idea`.** Upstream candidate (mattpocock/skills,
domain-modeling): as grill tickets resolve decisions, challenge fuzzy or
overloaded terms, stress-test the emerging model with concrete edge-case
scenarios, and keep `CONTEXT.md` plus applicable ADRs current inline. One
source per fact holds: the decision map owns the build decision, `CONTEXT.md`
owns ubiquitous language, ADRs own hard-to-reverse architectural state.
Integrates into the existing phase rather than adding a parallel skill. Kit
edit under the `craft-synthesis` discipline. Source: `IDEAS.md`, drained
here.

**FT166 (LOW) — `bench capture commit`: porcelain for the ambient capture
set.** Commit the capture surfaces (`.bench/learnings.md`, `IDEAS.md`,
`session-handoff.md`, `.bench/retros/`) with a conventional message under the
doc-only standing rule, so the plain-`git` step every session hand-assembles —
with FT107's empty-index hazard attached — becomes one sanctioned command.
Weigh it beside FT107's third clause, which owns the hazard prose; the
porcelain would remove the instruction rather than duplicate it.

The capture-only commit-path idea (2026-07-29) merges in — same visit, one
route: extend `internal/status`'s `captureOnlyStalePaths` from governing only
the gate-staleness signal to also governing the commit path, so a change
confined to capture-only files does not pay the full gate. Single-sourced in
the existing map, and sound only with a conformance check asserting no gate
check reads those paths; it is a path-scoped oracle exemption sitting near
the closed diff-scoped-gating ruling, so the exemption itself is a reviewer
ruling, not a default. The same-day journal entry lands here too: with two
unrelated dirty changes on `main` there is no sanctioned two-commit route
(`bench commit` refuses on out-of-set dirty paths, the guard refuses
`git stash`), and the session bundled with disclosure at `5fd3789` — until
this row ships, that is the convention to write down (bundle, lead with the
substantive change, name the ride-along, flag for veto). Sources:
`IDEAS.md`, drained here and in a prior run; `.bench/learnings.md`,
verdicted here.

**FT168 (LOW) — a fixture-selecting canary invocation.** Proving one changed
fixture currently costs the whole canary sweep: the light-path repair pass
needed evidence for a single race fixture, and the whole-sweep-only surface
invited expensive duplicate runs (the repair delegate launched one unbidden).
Add a `bench canary` path that runs one named fixture or family as iteration
evidence only — the full sweep remains the only thing the gate credits, so
this is a focused check, not a second oracle. Source: the light-path retro,
drained here.

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

Two singles ride along. `internal/gate/manifest.go`'s `dedupe` has no observable
effect — the scheduler's edge handling is already duplicate-tolerant and no
diagnostic renders `Needs` — but it implements a spec veto item literally, so it
is defensible dead code rather than a defect: keep or cut. And `bench learnings`
moved from a 5 MiB to a 2 MiB read bound: closing the divergence required
picking one number, and the slice chose the lower because `bench status`
already applied it to the same file — fail-closed and ambient-board-neutral,
but a 2–5 MiB journal that used to render now exits 1, a real behavior change
to keep or reverse. One line each closes this row.

## False greens — verdicts that credit unchecked work

Four rows, one failure class: a green whose warrant is missing — a stale
binary, a dead or skipping citation, a vacuous baseline, an unchecked absence.
Each hardens a different oracle surface, so they stay separate builds, but
they read and prioritize as one theme.

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
Source: `IDEAS.md`, drained here.

A third face, drained 2026-07-28, is row identity rather than citation
resolution: `bench coverage` emits only story/seam/red_signal, so a spec whose
rows share all three — `implement-spec-full-run`'s three story-3 hook rows —
cannot be enumerated row by row, and FT152's story-12 per-row accounting rule
is unexecutable as specified. Either the emission gains a stable row identity
(row number or the behavior field) or the rule names rows by story plus
behavior off the spec's own map — decide it alongside the check, same owner.
Found by the Codex falsification pass on `3eb1c9a`. Source: `IDEAS.md`,
drained here.

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
mean what it claims. Entry: `/bench-shape-idea`. Sources: `IDEAS.md`, drained
here; the `ft91-canary-compiled-bites` review S1, recoverable via `git show
4429b05:reviews/ft91-canary-compiled-bites.md`.

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

**FT141 (MEDIUM, evidence supplied) — `bench gate pin` records red verdicts,
so inherited reds stop reading as caused.** The pin records only the tree it
graded, not what failed there, so an agent that sees a failing check in a file
its diff never touched assumes causation and starts fixing — and that loop
does not self-terminate. Real case 2026-07-26: `main` red since `3c50349`,
and the FT91 build paid ~12 messages of git archaeology to prove the red was
inherited. Record which checks failed at which commit, so a stage starting
from a pinned baseline subtracts inherited reds automatically. This is the
containment half of the incident; the prevention half is the doc-only
shortcut's gate-anchored-surface exception.
Source: `IDEAS.md`, drained here.

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

**FT120 (LOW) — gate and canary test-harness defects nothing asserts.** Two
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
teardown before it recurs. Source: `IDEAS.md`, drained here.

## Standards debt — one batched light-path pass

Two rows plus FT142's standards track are shippable together as small
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

**FT94 (LOW, evidence supplied) — single-sourced `bench resume` summary
golden.** The resume summary line is asserted as a hardcoded exact-string
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

**FT125 (LOW, evidence supplied) — reader surfaces that return the slice, not
the file.** Two existing readers make a session pull a whole file to use one
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
requirement and harness pointer. Source: `IDEAS.md`, drained here.

## Dependencies

The dependent FT is named first. Only the literal table blocks work; the
recommended table is sequencing advice.

### Literal

| FT | Depends on | Why |
|---|---|---|
| FT107 | FT141 | Its fix-loop shrink rule cannot distinguish convergence until inherited reds are attributed to their pinned baseline. |

### Recommended

| FT | Better specified after | Why |
|---|---|---|
| FT71 | FT169 | The event schema should record the settled landing and recovery lifecycle rather than version an interim one. |
| FT100 | FT89 | Cut prose after the correctness and coherence pass establishes which guidance is still authoritative. |
| FT102 | FT128 | Audit tier-spending guidance against the corrected per-harness binding and guard semantics. |
| FT108 | FT164 | Define the refactor lane on the settled expand–migrate–contract and gate-cadence rules. |
| FT172 | FT106 | Reuse the document-claim probe for semantic roadmap claims instead of designing a second checker. |
| FT162 | FT169 | Build full-run subject resolution on the settled landing primitive. |
| FT166 | FT98, FT113 | Let recoverable set-aside and the completed capture-only path map define the commit command's smallest sound contract. |
| FT168 | FT153 | Binary freshness has shipped; expose focused canary execution after baseline meaning is settled. |
| FT169 | FT98 | Reuse recoverable discard in the landing contract; label resolution is already available. |

## Recommended sequence

1. `/bench-implement-spec` — FT126 recurrence tallying: approved with mixed story lines, so interactive orchestration owns the build.
2. `/bench-write-spec` — FT128: the line-enforcement fix now owns both the fork verdict and the static model-token sweep.
3. `/bench-shape-idea` — FT135: make installed pre-push protection report its resolved branch and template currency, then restore the sanctioned repair route.
