# Roadmap

The working prioritization document: every row is open work, verified against
the tree; a row leaves when the work ships (spec-retire) or a
`/bench-what-next` reconcile removes it. Raw capture never lands here — it
goes to `IDEAS.md` and enters only through a reviewed drain. A row for spec'd
work names its spec path (`specs/<slug>.md`) — that path is what `bench status`
cross-checks against the tree, so a row that omits it is a visible choice to
stay outside the ambient check.

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

**FT91 (HIGH, evidence supplied) — gate wall-clock proportional to the
diff.** Raised from MEDIUM by the reviewer on 2026-07-25: the gate's length is
the dominant cost of working on this repo for small changes, and the waiting is
paid by a human. Four arms have shipped: the phases were parallelized;
host-only test builds (retired 2026-07-20) removed the per-stage cold-`GOCACHE`
four-platform matrix; the canary concurrency budget landed 2026-07-24
(`bounds.CanaryInnerWidth` pins each inner gate to `GOMAXPROCS=2`,
stripped-then-set, workers derived from `runtime.GOMAXPROCS(0)`); and the
dev/ship tier split landed 2026-07-26 (spec `ft91-gate-tier-split`, retired) —
the ~372 s release-evidence probe, the cross-compile stress matrix, and the
release-only package suites moved behind `bench prep-release`, per-check timing
became permanent gate output, and `internal/conformance` left the inner suite,
closing the recursion hazard. Measured after the split: whole gate
5m50s → 4m36s, dev conformance 328.8 s → 106.9 s. The fifth arm, canary
check-scoping, landed 2026-07-26 (spec `ft91-canary-check-scoping`, retired):
a registry-owned family→check table scopes each conformance fixture's inner
run to the one check it grades, vacuity baselines are shared per scope group,
and the ship sweep pays two scoped runs plus one shared baseline instead of
three full inner gates. The long pole is now `package-core-guard` (~86 s),
which runs seven toolchain steps serially inside one Go test function.

The pipeline refactor is the sixth arm, and it ships in slices. Slice B landed
2026-07-27 (spec `ft91-phase-manifest-dag`, retired): the phase list became
project-owned data instead of Go code — an optional `.bench/phases.json` in the
graded root, loaded behind the unchanged `.bench/gate.sh` entry, with an absent
manifest keeping the built-in kit table so linked repos work through the
migration and every malformed class failing closed — and a DAG scheduler
replaced the serial/concurrent split, so a phase starts when its `needs`
complete, a red phase marks its downstream dependents skipped-with-cause while
independents run to completion, and a killed gate names the phases still
running. The two interim defects this row carried shipped with it: the
conformance subprocess env now scrubs all three control variables
symmetrically, and a failed probe spills full output to a named log instead of
truncating to the last 40 lines.

Slice C is the remaining pipeline work, and it is where the duration win is.
Hoist `checkGoCore`'s compile-and-test work out of the conformance test binary
into first-class gate phases with declared dependencies, leaving conformance to
grade structure only — this removes the remaining recursion surface story and
the stale-state class the FT91 build hit, and wins duration by overlapping
gofmt/build/vet/test/race/cross-compile instead of running them serially inside
one check. Complication that stays: `checkGoCore` grades arbitrary roots (canary
fixtures, linked repos), so it must split into a host-repo phase plus a
narrower structural check, and the 19+ package-core-guard fixtures asserting
its diagnostics move to grading the phase.

Cross-language incrementality stays a separate later capability behind the
existing revive trigger. Incrementality rides on Go's content-addressed test
cache today (a cache-invalidating change swung this repo 106 s → 826 s), so
making it language-agnostic means phases declaring input globs the runner
hashes — a small build system, real cost, not built speculatively. Shape it
against regroup-app rather than generalising from this repo.

The deferred arms are unchanged. Capping the outer conformance and contract
phases stays dormant unless contention flakes persist. Removing the hardcoded
`-count=1` changes what green proves about freshness — an oracle-semantics
decision — and cache infrastructure (hermetic build cache, verdicts keyed on
the pinned gate subject) waits on re-measurement after the pipeline lands.
Diff-scoped gating stays ruled out (contract and canary are behavior contracts
with no file→test map). None of them may weaken the oracle: green must keep
meaning the same thing, and any scoped verdict must be explicit evidence,
never a silent skip.

Entry: `/bench-write-spec` for slice C, inputs this row plus
`decisions/gate-pipeline.md` (the map is closed) and
`decisions/gate-concurrency.md`'s watch-outs. Sources: `IDEAS.md`, drained
here; `decisions/cost-follows-project-size.md`.

**FT131 (MEDIUM) — a stale binary under test makes contract rows lie in both
directions.** The AXI and runtime contract suites drive the built `dist/bench`,
not the package under test, so their verdict answers for whatever binary happens
to be on disk. In a fresh or salvaged worktree that binary predates the change:
during FT86 two of three rows in a correct slice went red on nothing but
staleness, and a delegate was nearly re-charged to fix code that was already
right. The dangerous direction is the reverse one — a stale binary that happens
to satisfy an assertion makes a broken change look green, a false done-claim the
gate catches only later. Prefer the single-source fix: have the contract
helper itself fail loudly when the binary under test is older than the package
sources it exercises, which removes the instruction rather than duplicating it.
The fallback — naming the `scripts/go-build.sh` rebuild in the guidance where
the phase names these seams — stays on this row and is taken only if the
in-helper staleness check proves unreliable; it was offered to the delegation
batch that shipped 2026-07-26 and deliberately left untaken there. Source: the
2026-07-25 learnings entry, verdicted in a prior drain.

**FT141 (MEDIUM, evidence supplied) — `bench gate pin` records red verdicts,
so inherited reds stop reading as caused.** The pin records only the tree it
graded, not what failed there, so an agent that sees a failing check in a file
its diff never touched assumes causation and starts fixing — and that loop
does not self-terminate. Real case 2026-07-26: `main` red since `3c50349`,
and the FT91 build paid ~12 messages of git archaeology to prove the red was
inherited. Record which checks failed at which commit, so a stage starting
from a pinned baseline subtracts inherited reds automatically. This is the
containment half of the incident; the prevention half — the doc-only
shortcut's gate-anchored-surface exception — rides FT107's sixth clause, and
FT107's fifth clause (the fix-loop shrink measure) explicitly depends on this
row landing first, since shrink is only meaningful over reds the diff owns.
Source: `IDEAS.md`, drained here.

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
duplicating it. Related but distinct: FT124's skip-reasons face makes skips
visible to a reader, while this row makes an invisible skip fail the coverage
oracle. Source: `IDEAS.md`, drained here.

**FT129 (MEDIUM) — a panic in the inner test binary reads as a canary that
stopped biting.** FT122's first gate went red as `canary
'worktree-lifecycle-safety-bypassed' did not bite`, naming an untouched and
correct fixture. The real cause was a new runtime contract test that sliced a
subject-reported hash without a length guard: against the fixture's stub
subject the hash was empty, the test panicked, the whole
`internal/contract/runtime` binary aborted, and the contract the fixture greps
for never ran. A check that never got to bite is not a check that stopped
biting, and reporting the two identically points every diagnosis at the wrong
file — this one cost a full canary sweep at HEAD just to establish the red was
ours. Teach the canary to detect a panic or non-test abort in the inner output
and report it as its own failure class, naming the panicking test. A second
arm to weigh rather than assume: a conformance check that no
`internal/contract/runtime` test slices subject-reported output without a
length guard, since the canary subject is always a stub and any such slice is
a latent tripwire-disarming panic. Source: the 2026-07-25 learnings entry,
verdicted in this drain.

**FT116 (MEDIUM) — data races in `guards.Scan` the gate cannot see.** Running
`internal/guards` under `-race` fails three tests on `main`:
`TestScanTimeoutPreservesPartialRowsAndHonestCounts` and
`TestScanEnumerationTimeoutUsesUnknownCounts` hit genuine data races inside
`guards.Scan`, and `TestGuardRowRejectsFIFOWithoutOpening` carries a 750ms
subprocess deadline too tight for a race-instrumented binary. The gate runs
`go test` without `-race`, so it never observes them; the races are
pre-existing, found while attributing a delegate's red during the FT87 slice 3
build. Two halves: fix the `guards.Scan` races, and decide whether the gate
should run `-race` on a subset so this class stays caught — the second half is a
gate-authoring reviewer decision (`craft-gate`).

The race is now attributed, and the row's stakes rose. `guards.Scan` leaks its
`enumerateGuards` goroutine past the `ctx.Done()` return (`guards.go:169-170`,
187), and the abandoned goroutine races the next test — confirmed 2026-07-26
with `go test -race -count=1 ./internal/guards`. FT91's new ship tier found it
on its first real run, which is the tier split working as designed; the leak
blocks a green `bench prep-release` on any host, alongside the
govulncheck-not-installed gap carried on FT142. Fix shape: `Scan` must not
leave `enumerateGuards` running after the timeout path returns. Source:
`IDEAS.md`, drained here.

**FT142 (MEDIUM) — FT91 review residuals: nine open findings, two tracks.**
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
same-package test beside the constant that owns it. Ship track, riding the
prep-release revive alongside FT116's race: the ship conformance step runs
with no `-timeout`, so the ~372 s probe plus all three release-only suites
share one default-bounded run — the 600 s hang hazard the spec claimed to fix
is restaged at ship; `govulncheck` resolves four levels deep instead of in
`requiredTools`, so a host without it burns the artifact matrix and a full
ship conformance run before dying — the concrete limit on the up-front
refusal; a second `prep-release` after an interrupt never cleans the orphaned
`dist/.preflight.*` staging directory; two concurrent conformance runs on one
root interleave the timing file (needs a coverage row or an explicit
Won't-handle line); and the release-only `go test` step the decisions
promised was silently folded into ship-tier `goCoreTestPackages` — flagged
for veto, not a defect. Source: the FT91 review, promoted at spec
retirement.

**FT143 (MEDIUM) — the family→check binding is enforced late and on one
surface only.** `ft91-canary-check-scoping` story 4's amended seam put the
unbound-family red in the conformance layer's kit-scoped family check, which
is correct about audience — the sweep grades adopting repos a kit-owned table
can never bind — but it lands the red after the cost and only on one path.
Two consequences. First, an unbound family resolves to no scope, so each of
its fixtures pays a full unscoped inner gate during the canary phase before
the conformance phase reds; loud but late, tolerable at dev and worse at ship,
which is exactly the cost this arm removed. Second, standalone `bench canary
[root]` runs `canary.Run` alone and never reaches the conformance registry, so
on that surface an unbound family sweeps unscoped and silent — only `bench
gate` catches it. Candidate fix: a cheap kit-root-scoped binding assertion
before the sweep starts, reachable from both entry points, leaving the
adopting-repo path untouched. Sources: `IDEAS.md`, drained here.

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
Source: `.bench/learnings.md`, drained here.

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
provider lookup, per FT97.

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

**FT98 (MEDIUM, evidence supplied twice) — one preserve-then-discard primitive;
three faces.** Three rows were faces of one missing primitive — a sanctioned,
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
fail-closed stays the default and the cut line is a reviewer decision. Face
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

**FT132 (MEDIUM, evidence supplied) — the roadmap row grammar is undeclared,
and the parser's malformed verdict proved it.** The absent-versus-empty face
this row opened with shipped in the FT86 review resolution — an empty
`ROADMAP.md` now fails closed as its own state, distinct from absence — but
the drain that followed hit the family's next instance: `ROADMAP.md` failed
its own parser. A body paragraph inside FT86's row led with bold,
`ParseDocument` treats any line opening with `**` as a row start, and the
line surfaced as a `malformed roadmap row` parse failure while the rest of
the row's body silently vanished from the snapshot — the drain's declared
complete local evidence was partial, and nothing said so loudly. The
2026-07-26 restructure rewove the file so it parses clean, but the grammar is
the reviewer's call, and it is the same absent/empty/malformed family FT86
shipped: either the roadmap declares that a row is exactly the
`**FT<n> (…) — title.**` form and no other line may open with bold — keeping
today's malformed verdict and documenting the constraint beside the
preamble — or row bodies may lead a paragraph with bold and the parser must
absorb them. Decide once, then move the contracts with the decision. Source:
`IDEAS.md`, drained here.

**FT126 (MEDIUM, evidence supplied) — `bench roadmap --context` reports facts
where the drain needs verdicts.** `/bench-what-next` step 1 orders every row
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

Not carried here: verifying what a row *claims* about the code
(does that symbol still exist, is that path still in that allowlist). That is
per-row semantic checking with no mechanical source, and it is the same problem
FT106 solves for docs — extend that probe to row bodies when it ships rather
than building a second mechanism.

**FT137 (MEDIUM) — `/bench-what-next` gains a restructuring step.** The drain
only adds and verifies rows, so the board accretes near-duplicates that each
pay a full pipeline and gate separately. Add a restructuring step to the
phase: merge rows that edit the same owner file into one batched row,
collapse rows that are faces of one missing primitive, fold leftovers into
their parent row, and group evidence-backed rows under a shared theme
header — each restructure proposed in the same batch diff the drain already
produces, never applied silently. The 2026-07-26 drain is the dogfood: under
exactly these moves it merged four delegation rows into FT96, three discard
faces into FT98, three staleness faces into FT113, and three standing-rule
rows into FT107 — and FT96, the largest of those merges, shipped as one spec
five days later. The 2026-07-26 drain is the second dogfood: both its inputs
landed as clauses on existing rows (a coverage-citation instance onto FT133, a
journal verdict onto FT107) rather than as two new rows, so the same moves
apply to intake, not only to periodic cleanup. Kit edit under the
`craft-synthesis` discipline. Source: `IDEAS.md`, drained here.

**FT134 (MEDIUM) — enumerated-posture tables become a Spec-axis audit
obligation.** A spec that enumerates call sites with per-site postures —
FT86's eleven-row DefaultBranch table — has no check that the build took the
assigned posture: two adopt call sites silently reversed "skip the candidate"
into a new fallback-`main` literal, and only the semantic review caught it.
Make an enumerated-posture table an explicit Spec-axis audit obligation in
`craft-review`, the way the coverage map already is: the reviewer walks the
table row by row against the diff. Kit edit under the `craft-synthesis`
discipline. Source: `IDEAS.md`, drained here.

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

**FT107 (MEDIUM) — the standing guidance rules, batched: one
always-loaded-prose diff.** Three rows edited the same standing-guidance
surface — `.bench/BENCH.md` and the phase prose beside it — and collapse into
one batched kit edit under the `craft-synthesis` discipline: one spec, one
review, one gate. First, the right-sizing table (this row's original charge).
`.bench/BENCH.md` says a few-line change doesn't need the full pipeline and
that you may propose a lighter path but must get an explicit OK first — so
every small change costs a round trip that the same paragraph already licenses
away ("if I give you a standing rule for changes of a given size, follow it
and stop asking"), yet no standing rule has ever been written. Write the table
beside that paragraph: change shape → the path it takes → the trigger that
forces escalation back to the full pipeline; no new phases and no new files.
Bound the light path by blast radius rather than file count — a change
crossing a declared seam takes the full pipeline, a change strictly inside one
is a light-path candidate — and carry one escalation trigger on an observable
rather than on self-assessed confidence: a small stated read budget spent
without naming the cause with evidence reroutes to `/bench-debug`, not onward
through the patch. Second (was FT110), two generation-shaping clauses on an
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
`bench outline`, currently underused. Third (was FT119), plain-`git` commit
safety during a squash-merge landing. The doc-only plain-`git` convention is
safe only when the index is otherwise empty: after `git merge --squash` the
index already holds the whole merged slice, and a bare `git commit` of a
just-added capture file once swept 649 insertions across eleven files into a
commit labelled "capture: …", landed with no gate grading it — `bench commit`'s
attribution check, which would have caught it, was bypassed by plain `git`.
Note in the phase guidance that the convention holds only with an empty index,
and that during a squash-merge landing every plain-`git` commit uses the
explicit pathspec form (`git commit -m "…" -- <path>`). Fourth (drained
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
readings is a stop-and-ask. Fifth (parked 2026-07-26 from session
conversation), the fix-loop escalation trigger — the symmetric partner to this
row's first clause, which reroutes a *read* budget spent without traction. A
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
Ordering: it depends on the `bench gate pin` red-verdict idea parked the same
day, because the shrink measure is only meaningful over reds the diff owns —
inherited and spec-predicted reds are constant noise that would trip a false
stop. Evidence for both halves is the 2026-07-26 FT91 gate: three reds, one
inherited from `3c50349` and two predicted by the spec, none belonging to the
diff, resolved only by hand-run git archaeology. Source: session conversation,
parked here by reviewer instruction rather than through a drain. Sixth
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
decided when the batch builds. The containment half of the same incident — the
gate pin recording red verdicts so inherited reds stop reading as caused — is
FT141, which this row's fifth clause already depends on. Deliberately
not batched here: FT130 (mid-gate
`bench idea`) stays its own row because its preferred fix is mechanical, in
the CLI rather than in this prose. Background:
`docs/reporesident-distillation.md` §3 and §6.

**FT58 (LOW) — hardened pool roots.** The identity-safe lock protocol shipped
with the worktree-lifecycle build: a live owner is never aged out, competing
reclaimers serialize through a rename-and-identity-check takeover, and a
successor's lease survives release, each with a red-capable test. What remains
is pool-root hardening: permission failures on Bench-selected pool roots
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

**FT97 (LOW, evidence supplied) — harness-native agent-line denial.** The
agent-line deny message single-sources its bound-tiers listing, which leads
with the three tier ids and trails the harness aliases; inside a Claude Code
session the aliases are the only tokens the Agent tool can pass, so the error
leads with ids that harness cannot use (observed 2026-07-19). The design is
already decided: the closed `decisions/multi-harness-line-binding.md` map
answers the schema question — symmetric per-harness bindings with no canonical
family, each layer reporting its own harness's tokens. The row's work is
building that map; enforcement stays exact-string with no provider lookup.

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

One of three angles on the cost-follows-project-size complaint, with FT91 and
FT136; one `/bench-shape-idea` session should shape the three together — the
rows stay separate because the owners differ. Kit edit under the
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

**FT105 (LOW) — committed reports that contradict the tree.** A capture-style
report can outlive its factual state and misdirect a cold session. Invariant 3
says docs describe the current decided state, so a tracked report that
contradicts the tree is a defect. The work is the discipline alone — a
capture-style report carries its status at the top and is re-read at the
moment of commit, not only when written — as a kit edit under the
`craft-synthesis` discipline. The gate half is cut: a
doc-conformance check grepping for phrases like "nothing is committed" would
fire on ordinary prose in a spec's Problem section, and a check that cannot be
shown to bite does not earn gate weight.

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

**FT113 (LOW) — sanctioned capture writes pollute the ambient staleness
signals.** A handoff rewrite is doc-only capture, but it is a tracked write,
so it shows up in the very signals a resuming session reads. Two faces, one
owner (`internal/status`), so they ship together. The gate face: the rewrite
marks the gate plain-stale, costing the next commit a full 10–15min re-run;
`ROADMAP.md` and `IDEAS.md` are already on `status.captureOnlyStalePaths`, so
add `session-handoff.md`. Observed 2026-07-23 right after the FT109 close. The
board face: `bench handoff` derives its Next command from `status.Signals`,
which counts the handoff's own dirty path, so the board a session reads is not
the tree as the reviewer left it and story 18's byte-identical guarantee cannot
hold on a tracked repo (FT122 review finding, 2026-07-25). The FT122 build
already carries a bespoke subtraction for the dirty-path *count*
(`handoff.handoffIsDirty`); a path-exclusion option on the status query would
single-source that and cover the Next derivation too. A third face (was
FT121), the `bench commit --spec` flip: the flip from `Status: staged` to
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
cosmetic and documenting it.

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
32 lines over its 400-line budget after the FT91 arms, so the added test lands
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
holds LOW until an acceptance trigger that needs measurement (FT136's
cheap-tier re-test, the mid-tier default) actually blocks on it. Candidate
metrics: delegate tokens per slice, coordinator tokens, review findings per
implemented slice split by axis (Standards/Spec/Coverage), rework tokens
spent after a build already went green, gate runs per spec, iterations
against the declared cap, and tier declared versus tier actually used. The
only baseline that exists (FT86 review resolution: 26 findings, roughly 350k
delegate tokens of rework) survives by accident in one session transcript
rather than in any artifact. Open questions — which metrics earn their
capture cost; where they live (a `bench` subcommand, the gate cache, or the
journal); who writes them, given the harness-independent substrate is the
shift loop and the git hooks rather than any harness; per-slice versus
per-spec granularity; retention and pruning; agent-facing AXI or
reviewer-facing — mean the work starts as a grill. Entry:
`/bench-shape-idea`. Source: `IDEAS.md`, drained here.

**FT139 (LOW) — the runtime gate capability-skip contract test inherits
ambient strict mode.** The suite builds its subtest environments by adding to
the ambient environment rather than scrubbing it
(`internal/contract/runtime/runtime_gate_capability_skips_test.go`), so an
inherited `BENCH_REQUIRE_CAPABILITIES` reaches the unset-flag subtest and
fails it — observed under the FT86 spec-mandated strict run, where the flag
is set by design. Scrub the variable in the harness's environment setup so
each subtest's premise is exactly what it asserts. Source: `IDEAS.md`,
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

## Session tax — evidence-supplied reader rows

Three rows, one theme: each is a measured, recurring plumbing cost from the
week-of-2026-07-19 transcript evidence, and each builds a reader or resolver
over a surface that already exists — never a second oracle. (FT127, formerly
the fourth candidate here, collapsed into FT98's discard primitive instead.)

**FT123 (MEDIUM, evidence supplied) — worktree resolution by label.** A session
reaches its own worktree by pasting a 110-character hashed path, and that
plumbing is the single largest measured session tax on this repo: 1495 of 4115
Bash calls in the week of 2026-07-19 (36%) were
`cd /home/mgibs/.bench/worktrees/bench-<n>/<32hex>-<32hex> && …`. Add
`bench worktree path <label>` and `bench worktree exec <label> -- <cmd>` so the
label is the handle. Hard requirement: both emit and accept `~`-relative paths,
because an absolute `/home/<user>` path in a transcript or a committed artifact
does not port across machines. Untouched by the delegation-discipline batch
that shipped 2026-07-26: that one settled how a coordinator *cuts* parallel
worktrees (`bench worktree create --request <opaque-id>`), leaving how a
session names one it already owns exactly where it was.

**FT124 (MEDIUM, evidence supplied) — `bench test [pkg]`, structured go-test
triage.** Sessions re-derive the same `go test` filter every time: 698 of 797
`go test`/`go build` calls in the week of 2026-07-19 were piped through ad-hoc
`head`/`grep`/`tail`, and the patterns recur verbatim across sessions
(`^(ok|FAIL|---)`, `--- (SKIP|PASS|FAIL)`, `t\.Skipf?\(`, `go test failed`).
Emit TOON instead — package, failing test, first failure line, skip reasons —
under the `craft-cli` standards the rest of the surface already follows. This
is a reader over `go test`, not a second oracle: the gate stays the only thing
that calls work done.

A second face (drained 2026-07-25) makes the skip-reasons field load-bearing
rather than a nicety: a capability or environment skip is invisible without
`-v` — `TestRootConformance` without `BENCH_CONFORMANCE_ROOT` announces itself
via `capability.Environment`, but `go test ./internal/conformance` still
prints a bare package-level ok in ~0.006s, and an FT86 review delegate read
that as a passing conformance sweep. The skip line is the only thing
distinguishing an honest skip from a deleted assertion — the same property
`projects/benchkit.md`'s cold-session notes call load-bearing for the
canary — so the reader surfaces skips with reasons at default verbosity, and
the row decides whether a bare sub-second package ok deserves its own visible
marker.

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

## Recommended sequence

1. `/bench-write-spec` — FT91 slice C: split `checkGoCore` into gate phases.
   Gate wall-clock is the reviewer's stated dominant cost, slice B shipped the
   manifest and DAG mechanism with nothing yet consuming it, and the long pole
   is still `package-core-guard`'s seven serial toolchain steps inside one
   check. Inputs: the FT91 row and `decisions/gate-pipeline.md`.
2. `/bench-write-spec` — FT71, versioned local shift evidence. The remaining
   HIGH bank-track row; the repository-controlled bank evidence requirement
   keeps it active.
3. `/bench-write-spec` — FT133, red-signal citations that resolve *and*
   execute. Twice demonstrated now (a dead `-run` filter, a capability skip
   printing `ok`), and it is the coverage oracle itself reading green on
   nothing — every spec written before it lands can carry the same hole.
