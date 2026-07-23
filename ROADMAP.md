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

**FT87 (MEDIUM) — command-wide parser and security-evidence capability.**
Slices 1 and 2 shipped: the bounds core, explicit pinned bounded repair, one
Bench identity with complete package metadata, and the FT83
offline/network-control evidence record. What remains is slice 3: centralize
argument grammar, anchor coverage at the repository root, support leading-dash
and directory-scoped commit paths, treat help as success, make capability
skips explicit evidence, and decouple security-test deadlines from the
subprocesses they bound. The closed FT87 decision map (tickets #7 and #8) owns
the cut.

Sources: `RR:C-09`, `RR:C-11`, `RR:C-12`.

**FT86 (HIGH on the bank track; MEDIUM otherwise) — fail-closed control records
and single-sourced repository facts.** Coverage validation requires a map or
explicit historical marker and validates exact positive story membership and
forward ranges. Learnings, maps, roadmap, outline, and status distinguish
absent, empty, malformed, unreadable, wrong-type, traversal, and
unsupported-schema states; only absence is an authoritative empty state. Gate
paths fail closed.

Default-branch resolution has one owner and callers handle unknown rather than
fabricating `main`. Closure uses hostile filesystem, parse, story-gap/range,
master-only, and unknown-default fixtures in both contracts and canaries.

Sources: `RR:C-01`, `RR:C-02`, `RR:C-03`; `RC:H-08`.

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

**FT91 (MEDIUM, evidence supplied) — gate wall-clock proportional to the
diff.** Two arms have shipped: the phases were parallelized, and host-only
test builds (retired 2026-07-20) removed the per-stage cold-`GOCACHE`
four-platform matrix. The graduation trigger is now met — fresh post-host-only
measurement (2026-07-22): the gate still takes 10–15 minutes on the kit repo,
largely because the canary phase nests whole gate runs, oversubscribing a
16-core box to load average ~123, so most of the wall clock is contention,
not work; `internal/gate/phases.go` also hardcodes `-count=1`, disabling Go's
test cache unconditionally. First arm: core-count-aware gate/phase
concurrency — detect the machine's cores and scale so nested phases cannot
oversubscribe the box. That cap is also the load lever behind the retired
gate-trustworthiness work: the marker stalls and cleanup flakes that arm fixed
were all contention symptoms, so capping concurrency protects the verdict as
well as the clock. Diff-scoped gating is unsound here (contract/canary are behavior
contracts with no file→test map). The remaining arms — a shared hermetic
build cache, caching keyed on the pinned gate subject, or scoped verdicts —
must not weaken the oracle: green must keep meaning the same thing, and any
scoped verdict must be explicit evidence, never a silent skip. Starts as a
grill (`/bench-shape-idea`) because the cut line between speed and oracle
authority is a reviewer decision.

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

**FT107 (MEDIUM) — the standing right-sizing rules, written down.**
`.bench/BENCH.md` says a few-line change doesn't need the full pipeline and
that you may propose a lighter path but must get an explicit OK first — so
every small change costs a round trip and the reviewer answers the same
question with the same answer. The same paragraph already licenses the fix
("if I give you a standing rule for changes of a given size, follow it and
stop asking"), but no standing rule has ever been written. Write the table
beside that paragraph: change shape → the path it takes → the trigger that
forces escalation back to the full pipeline. No new phases and no new files.
Bound the light path by blast radius rather than file count — a change
crossing a declared seam takes the full pipeline, a change strictly inside one
is a light-path candidate — because seams are vocabulary the profile already
carries and a file count fails on a monorepo. Carry one escalation trigger on
an observable rather than on self-assessed confidence: a small stated read
budget spent without naming the cause with evidence reroutes to
`/bench-debug`, not onward through the patch. Kit edit under the
`craft-synthesis` discipline. Background:
`docs/reporesident-distillation.md` §3.

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

**FT96 (LOW) — parallel-delegate worktree assignments.** The WorktreeCreate
hook keys assignments per session, so launching several parallel write
delegates with harness worktree isolation grants one assignment and refuses
the rest ("conflicts with its existing assignment"); the manual
`bench worktree create --request <distinct-id>` route works and respects the
lifecycle but is undocumented. Either key hook assignments per delegate or
document the `--request` route as the canonical parallel-delegate path in
`craft-delegate`. The same guidance must state that a whole-tree gate run is a
serialized resource: four concurrent worktree `bench commit` gates flaked three
load-sensitive contract tests (cancellation timing, tempdir cleanup, release
reproducibility probe) that all pass serially, so delegates stop at "diff ready,
focused tests green" and the coordinator runs `bench commit` per worktree one at
a time. Kit edit under the `craft-synthesis` discipline.

**FT97 (LOW, evidence supplied) — harness-native agent-line denial.** The
agent-line deny message single-sources its bound-tiers listing, which leads
with the three tier ids and trails the harness aliases; inside a Claude Code
session the aliases are the only tokens the Agent tool can pass, so the error
leads with ids that harness cannot use (observed 2026-07-19). The design is
already decided: the closed `decisions/multi-harness-line-binding.md` map
answers the schema question — symmetric per-harness bindings with no canonical
family, each layer reporting its own harness's tokens. The row's work is
building that map; enforcement stays exact-string with no provider lookup.

**FT98 (MEDIUM, evidence supplied twice) — discard path for content-landed
recovery payloads.** Recovery cleanup fail-closes permanently when a payload's
content landed through different commits: the FT83 delegate payloads are strict
subsets of the default branch by diff, yet `bench worktree recovery <ref>
--apply <fingerprint>` still returns `retain` because landed-proof requires the
payload commit itself (observed 2026-07-20), so proven-redundant refs
accumulate with no sanctioned exit. Recurred 2026-07-22: payloads whose patch
reverse-applies cleanly on main were still retained because `git cherry` misses
reshaped commits, and the reviewer had to hand-delete refs and intent entries —
the exact manual surgery the lifecycle exists to prevent. The containment
primitive now exists — `LandedInDefault` proves patch-id containment for
landed-branch pruning — so either route recovery-payload landed-proof through
it or add an explicit reviewer-authorized discard path; fail-closed stays the
default, and the cut line is a reviewer decision.

**FT99 (LOW) — spec problem-premise verification.** A spec compiled from a
closed decision map can inherit a problem statement the tree has since
falsified: the retired `minimal-subprocess-data-exposure` spec claimed the
project gate "gets everything except `BENCH_KIT` and `BENCH_WRAPPER`", which
FT78 had already
fixed, and the build reached stage 1b on that false premise before a contract
delegate caught it. Require every "today the code does X" claim in a spec's
Problem section to be checked against the tree at spec time, with the check
named in the spec — the same standard the coverage map already applies to its
red signals. Next action is the kit edit to `/bench-write-spec` and
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

Kit edit under the `craft-synthesis` discipline. Background and the
alternatives considered: `docs/reporesident-distillation.md` §8.

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency
loop.** A kit edit that instructs spending a model tier can contradict the
escalation policy without any loop catching it: the widened write-spec step-9
triggers shipped an automatic top-tier spawn past review (observed
2026-07-22; corrected in the mid-tier rerouting commit). Make
`craft-synthesis`'s consistency loop name the escalation policy as a standing
cross-check for any kit edit that spends a tier. Kit edit under the
`craft-synthesis` discipline.

**FT103 (LOW) — existence-checked absence evidence.** A delegate's payload
slice landed with a misspelled kit-only allowlist row (`craft-synthesis` for
the actual `bench-craft-synthesis`), so its contract passed by asserting the
absence of a path that never existed and the skill kept shipping to
consumers — a vacuous green the gate cannot see. Two halves: add to
`craft-delegate`'s verification discipline that when a delegate's evidence
rests on an absence, exclusion, or withholding assertion, the named
identifiers must resolve to real things before the claim is accepted (kit
edit under the `craft-synthesis` discipline); and confirm the gate holds a
per-source existence check on the consumer-payload allowlist — the emptied-set
vacuity closed with the FT85 fix commit, the per-path existence guard is the
remaining cheap single-source check.

**FT104 (LOW) — stop rule for known-flake commit refusals.** Retrying a
recorded flake is not iteration toward green: the FT85 review-fix commit was
refused twice by `TestFT78Story5ProofLedger` under gate load (green in
isolation both times), and the third identical run passed with no code
change — ~35 minutes of wall clock bought nothing. Kit edit under the
`craft-synthesis` discipline: the gate/commit discipline states that when a
commit is refused twice by the same test already recorded as a known flake
and proven green in isolation, stop and hand the blocked commit to the
reviewer with the evidence instead of re-running. Replaces the retired FT95
"retry once" operational line.

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

**FT109 (LOW) — handoff shape, capped and computed.** The phase-close rule in
`AGENTS.md` says the closing message emits a continuation prompt or updates
`session-handoff.md`, so the file's shape is entirely up to whoever writes it:
it is currently a free-form narrative with no template, no cap, no
rewrite-in-full rule, and nothing saying what happens when it disagrees with
the tree. Put the template inside `session-handoff.md` itself, so the template
cannot drift from the artifact it describes, and add the rewrite-in-full and
conflict rules to the phase-close paragraph in `AGENTS.md` — the handoff is
pruned rather than accreted, because a fresh session pays for every line it
reads cold. Make the conflict rule computed rather than remembered: the
handoff records the commit it was written at, and `bench status` compares that
against `HEAD` and reports staleness, which the SessionStart hook already
prints. That is one machine-readable line in the handoff and one comparison in
`session-inspect`, and it turns "trust git over the handoff" from a rule the
next session must recall into an ambient fact. Kit edit under the
`craft-synthesis` discipline. Background:
`docs/reporesident-distillation.md` §4.

**FT110 (LOW) — two generation-shaping clauses on an observable.** Invariant
4's "read the surrounding code before you write" is true but unfalsifiable —
*surrounding* is undefined and a session always believes it read enough. Two
sharper clauses. First, in `.bench/BENCH.md` as a clause on invariant 4 rather
than a new invariant: never call an API or function whose definition you have
not read this session, and verify behavior in source rather than from memory.
Second, in `craft-seams` beside the seam-finding guidance: an exploration
budget — a small stated number of reads without traction means stop reading,
run `bench outline`, and pick the seam from the index. State the budget as one
the session declares it has spent, so the behavior is observable in the
transcript rather than counted internally. The re-scope target is computed and
cannot go stale, which is the half worth having and also real pull for
`bench outline`, currently underused. Kit edit under the `craft-synthesis`
discipline. Background: `docs/reporesident-distillation.md` §6.

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

1. `/bench-write-spec` — FT87 slice 3, command-wide parser and
   security-evidence capability. The gate-trustworthiness row closed and
   retired, so this is the highest row open, and its decision map (tickets #7
   and #8) is already closed.
2. `/bench-shape-idea` — FT91's first arm, core-count-aware gate/phase
   concurrency. With the verdict-trustworthiness arms landed, contention is
   the last unaddressed cost on every landing; the grill settles the cut line
   between speed and oracle authority before anything is built.
3. `/bench-write-spec` — FT86, fail-closed control records and single-sourced
   repository facts. The highest bank-track row still open.
