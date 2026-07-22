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

**FT85 (HIGH) — least-privilege consumer payload and one coherent phase
contract.** Spec: `specs/consumer-payload-and-phase-contract.md`. Generate the consumer payload from one canonical allowlist and
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

**FT95 (MEDIUM, evidence supplied) — attributable, load-tolerant gate runs.**
The gate flake class recurred at the FT76 close: `bench gate` went red 5 of 7
runs on an identical tree, alternating between the conformance phase's inner
core `go test` (4×, output discarded — the diag says only `go test failed`) and
`TestBinaryRepairContracts/repair_losing-racer` (1×, a hard 2s wall-clock
`time.Now()` sync-marker deadline). Every suite passes solo and pairwise;
dmesg shows WSL2 "time jumped backwards" clock jumps under the gate's peak
parallel load, and a PATH shim cannot instrument the inner run because
`go test` prepends `$GOROOT/bin` to child PATH. Three arms, none weakening the
oracle: (1) the conformance diag carries the tail of the inner test output so
failures self-attribute; (2) wall-clock `time.Now()` test deadlines (the repair
sync marker and siblings) scale or move to monotonic-friendly generous bounds;
(3) consider capping gate phase parallelism on hosts where load provokes the
class. Deadline widths and the diag shape are gate-authoring decisions
(`craft-gate`). Retry once on the same tree and line remains the operational
response to a transient red.

**FT100 (LOW) — prose-weight pass on the kit's guidance surface.** Apply the
gate's "prove it bites" standard to prose: audit the craft-skill library and
the communication protocol so each skill and always-loaded clause cites an
observed failure it prevents (from the learnings journal or session evidence),
merge overlapping craft docs, and shrink the always-loaded `BENCH.md` rules to
demonstrated-delta clauses. Distinct from FT89, which fixes guidance
*correctness*; this row cuts guidance *weight*. Kit edit under the
`craft-synthesis` discipline; starts as a grill (`/bench-shape-idea`) because
the cut line on always-loaded rules is a reviewer decision.

**FT101 (LOW) — multi-context domain docs for monorepos.** A monorepo has more
than one bounded context, but the kit assumes one `CONTEXT.md`. Support a
root `CONTEXT-MAP.md` pointing to per-context `CONTEXT.md` files, add a
single- versus multi-context question to `/bench-setup-repo` Section C, and
teach every `CONTEXT.md` consumer (phase commands, skills) the layout. Kit
edit under the `craft-synthesis` discipline.

**FT102 (LOW) — escalation-policy cross-check in the synthesis consistency
loop.** A kit edit that instructs spending a model tier can contradict the
escalation policy without any loop catching it: the widened write-spec step-9
triggers shipped an automatic top-tier spawn past review (observed
2026-07-22; corrected in the mid-tier rerouting commit). Make
`craft-synthesis`'s consistency loop name the escalation policy as a standing
cross-check for any kit edit that spends a tier. Kit edit under the
`craft-synthesis` discipline.

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

1. `/bench-implement-spec` — FT85, least-privilege consumer payload and one
   coherent phase contract (spec staged at
   `specs/consumer-payload-and-phase-contract.md`).
2. `/bench-write-spec` — spec FT87 slice 3, command-wide parser and
   security-evidence capability.
3. `/bench-debug` — FT95, attributable and load-tolerant gate runs (named
   culprits in hand; it taxes every landing).
