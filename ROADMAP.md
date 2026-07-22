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

**FT76 (HIGH) — one-command, repo-aware Bench bootstrap.** Build the Bench
counterpart to Matt Pocock's `setup-matt-pocock-skills`: an explore → present
→ confirm → write onboarding flow behind one public entry point (`npx --yes
redbench@<immutable-version> setup`, with `bench setup` as the installed
equivalent). It composes the existing link/init/setup seams; it must not grow a
second asset installer or a second source for shared rules.

The bootstrap inspects Git state, remotes, harness files, language/build/test
signals, existing gates and project docs; proposes the smallest project
profile and gate; asks only unresolved decisions one at a time; previews the
plan; then transactionally seeds the managed Bench assets, a validated generic
project profile and template set, gate/canary, context and ADR scaffolds, and
the repo-local executable path. Every emitted pointer is validated against the
packed artifact. It adds or updates one marker-owned Bench block in the
appropriate `AGENTS.md` or `CLAUDE.md` without overwriting surrounding project
content, duplicating rules, or creating competing instruction files. A
preserved Claude instruction file that does not import Bench is an explicit red
doctor state, never a false all-harness success.

The acceptance seam is the artifact a new user receives, not the source tree:
pack → install → bootstrap in empty and established repositories, including
spaces, existing `AGENTS.md`/`CLAUDE.md`, hooks/settings, monorepos, no-global
installs, fresh clones, offline artifacts, and every advertised target. Re-runs
are idempotent and downgrade-aware; clean managed content updates, modified
content is preserved and reported, stale assets reconcile, and `unlink`
remains a true inverse. Success is at most one shell command plus one
harness-native setup conversation and ends with a runnable local `bench`, an
honest doctor report, an explicit reload instruction, and the exact next
action. The closed FT76 decision map settles instruction-file selection,
inferred-versus-confirmed facts, conflict semantics, and the seam between the
deterministic CLI and the harness conversation; FT84's transactional
lifecycle — the build-first dependency — has shipped. Next action is
`/bench-write-spec`.

Sources: `RR:A-02`, `RR:A-03`, `RR:A-04`, `RR:S-09`.

**FT85 (HIGH) — least-privilege consumer payload and one coherent phase
contract.** Generate the consumer payload from one canonical allowlist and
exclude kit-only assessment, update, and synthesis surfaces. Consumers receive
a narrow version-pinned, manifest-aware upgrade path. One capability-aware
delegation policy lives in `craft-delegate`; one phase owns the green landing
commit and spec state transition; debug preserves a red observation without
violating the commit-on-green invariant; and shaping requirements agree across
README, shape, spec, and implementation guidance.

Closure combines a forbidden-payload link/package contract with workflow
scenarios for clear ideas, unavailable write delegates, implementation review
handoff, final landing, and failing bug reproduction.

Sources: `RR:S-01`, `RR:S-02`, `RR:S-03`, `RR:S-04`, `RR:S-05`; `RC:M-03`.

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

**FT91 (LOW) — gate wall-clock proportional to the diff.** Two arms have
shipped: the phases were parallelized, and host-only test builds (retired
2026-07-20) removed the per-stage cold-`GOCACHE` four-platform matrix that
drove the gate to ~10–15 minutes. The remaining arms — a shared hermetic build
cache, caching keyed on the pinned gate subject, or scoped verdicts — must not
weaken the oracle: green must keep meaning the same thing, and any scoped
verdict must be explicit evidence, never a silent skip. Starts as a grill
(`/bench-shape-idea`) because the cut line between speed and oracle authority
is a reviewer decision. Graduate only on a fresh post-host-only measurement
showing the gate still demonstrably drags shift iteration; the pre-host-only
timings are no longer evidence.

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

**FT98 (LOW, evidence supplied) — discard path for content-landed recovery
payloads.** Recovery cleanup fail-closes permanently when a payload's content
landed through different commits: the FT83 delegate payloads are strict subsets
of the default branch by diff, yet `bench worktree recovery <ref> --apply
<fingerprint>` still returns `retain` because landed-proof requires the payload
commit itself (observed 2026-07-20), so proven-redundant refs accumulate with no
sanctioned exit. The containment primitive now exists — `LandedInDefault`
proves patch-id containment for landed-branch pruning — so either route
recovery-payload landed-proof through it or add an explicit
reviewer-authorized discard path; fail-closed stays the default, and the cut
line is a reviewer decision.

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

**FT95 (LOW, evidence supplied) — attributable compiled-core gate failures.**
The serialized root-conformance gate has again reported only `go test failed`
from its inner whole-core test run on an otherwise idle machine, while an
immediate identical-tree, uncached package run passed every package. The probe
currently discards the inner stdout and stderr, so the failing package and case
cannot be captured and the intermittent defect cannot be isolated. First make
the failure emit bounded, control-safe diagnostic evidence; then reproduce and
deflake the attributable case without weakening the oracle. Retry once on the
same tree and line remains the operational response to a transient red.

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.** Apply the
gate's "prove it bites" standard to prose: audit the craft-skill library and
the communication protocol so each skill and always-loaded clause cites an
observed failure it prevents (from the learnings journal or session evidence),
merge overlapping craft docs, and shrink the always-loaded `BENCH.md` rules to
demonstrated-delta clauses. Distinct from FT89, which fixes guidance
*correctness*; this row cuts guidance *weight*. Kit edit under the
`craft-synthesis` discipline; starts as a grill (`/bench-shape-idea`) because
the cut line on always-loaded rules is a reviewer decision.

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
   from an isolated prefix without a source checkout; FT76 preserves existing
   instructions, settings, and hooks and is idempotent and reversible.
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
demonstrably burn turns on symbol search.

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

1. `/bench-write-spec` — spec FT76, one-command repo-aware Bench bootstrap
   (shape closed, FT84 dependency shipped).
2. `/bench-write-spec` — spec FT85, least-privilege consumer payload and one
   coherent phase contract.
3. `/bench-write-spec` — spec FT87 slice 3, command-wide parser and
   security-evidence capability.
