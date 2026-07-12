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
their active findings. Each of the 53 release-readiness findings and 17
compliance findings appears exactly once on a `Sources` line below. `RR:` means
the release-readiness assessment and `RC:` means the repository-controlled
compliance assessment.

## Features, in priority order

**FT77 (CRITICAL, next) — ownership-safe worktree cleanup.** Bench-created
worktrees get one verifiable ownership identity. SessionStart and
`resume-clean` may act only on a matching non-live assignment with its verified
Bench lock; Bench unlocks only inside the proven cleanup transaction. Foreign
worktrees are informational. Detached work must have a durable Bench-owned
recovery ref before any removal. Explicit cleanup is path-scoped, dry-run first,
and requires affirmative destructive acknowledgement.

Closure requires black-box branch and detached-unique preservation contracts
through SessionStart, `resume-clean`, and explicit cleanup, plus a canary that
goes red when the ownership check is bypassed. Until then, automatic broad
cleanup remains unsafe.

Sources: `RR:R-01`; `RC:C-01`.
Decision map: `decisions/ownership-safe-worktree-cleanup.md`. Spec:
`specs/ownership-safe-worktree-cleanup.md`. Next: `/bench-implement-spec`.

**FT78 (CRITICAL) — oracle-bound, fail-closed gate verdicts.** Version the gate
cache and bind green reuse to the working-tree hash, fully resolved gate kind
and command, executable/script content, relevant configuration, schema, and
freshness. Verdict replacement is same-directory, synced, and atomic. A red
verdict that cannot be recorded invalidates prior reusable state and fails the
requested action.

Closure covers gate-command, script-content, and auto-detected-oracle changes,
plus malformed cache and write failure. Oracle A green must never authorize a
commit under oracle B.

Sources: `RR:R-02`; `RC:C-02`.

**FT79 (CRITICAL) — lossless shift recovery and truthful result states.** Every
failure after agent mutation preserves either a locked recovery worktree or a
durable recovery ref and prints its exact location. Staging, adapter, commit,
signal, and teardown errors propagate. Shift requires an objective, validates
a positive bounded iteration cap, has a wall deadline, and distinguishes
complete, incomplete, failed, interrupted, and no-op outcomes with honest exit
codes and structured results.

Closure is a recovery matrix for missing Git identity, failing commit hooks,
staging failure, adapter start/exit failure, signal interruption, cap
exhaustion, no-op iteration, and teardown error. Successful agent work must
remain recoverable in every case.

Sources: `RR:R-03`, `RR:R-08`; `RC:H-05`.

**FT80 (HIGH) — static, bounded guard discovery.** Guard descriptions come
from static managed metadata, not executing every shell file in the hooks
directory. Only an exact managed allowlist may execute where execution is
actually required; unknown additions are reported without running. The whole
SessionStart inspection has one aggregate deadline.

Closure requires an unwired sentinel that never executes through `bench
guards` or SessionStart, plus a canary proving the non-execution rule bites.

Sources: `RR:R-07`; `RC:H-02`.

**FT81 (HIGH) — usable, platform-correct distributable runtime.** Build wrapper
and platform packages in clean staging trees from explicit allowlists. The
wrapper contains no build-host binary or nested platform-package tree; runtime
selection reaches the declared native package first. Packed installs can run
setup, `link`, `init`, and `doctor`. The stable shim routes maintenance to the
installed kit and operational commands to the project-local runtime. A fresh
clone retains a runnable tracked local path without referring to maintainer-only
build scripts.

The acceptance seam is the exact tarball: isolated-prefix tests exercise
version, setup/link, doctor, one local operation, relink, fresh clone, and
unlink. Native smokes prove executable mode, embedded version, target format,
and selection on every advertised OS/architecture; FT83 supplies the artifact
governance and offline evidence for the same bytes.

Sources: `RR:R-04`, `RR:R-05`, `RR:R-09`, `RR:A-06`; `RC:C-04`.

**FT82 (HIGH) — one authoritative release preflight.** A repository-owned
preflight runs the full Bench gate, race tests, vet, vulnerability scan,
artifact build and inspection, and clean-room installed smoke against the same
commit. PR/push verification and tag publication both call it. The Go patch
toolchain is fixed and pinned; scanner failures block release unless a
documented, reviewable exception schema applies.

Preflight enforces exact SemVer and equality among tag, package, changelog,
binary, source commit, and release manifest, plus intended release-line
ancestry. Machine-readable phase records are retained and a canary or
structural contract fails if a required phase is deleted or bypassed. Immutable
registry publication and release versioning remain blocked until this is green.

Sources: `RR:R-06`, `RR:A-01`; `RC:C-03`.

**FT83 (HIGH on release and bank tracks) — governed, offline-verifiable release
bundle.** Every independently published package carries the license, security
and support policy, dependency/license notice, SBOM, checksums, and package
inventory. Define an exact supported-platform contract; use static builds where
supported, strip release symbols, and exercise native and musl-compatible
smokes. Produce network-independent per-platform archives and npm tarballs
with documented local and internal-registry installation.

Publication preflights every name, stages under a non-default tag, verifies
already-present immutable digests, waits for dependencies, publishes the
wrapper last, and promotes only after verification. A deterministic release
manifest binds source/version, Go/Node/npm versions and flags, dependency and
platform manifests, gate/race/vet/vulnerability results, artifact inventories
and SHA-256 digests, reproducibility comparison, rollback target, and the
machine-readable evidence produced by the other release-readiness rows.

Repository governance also defines supported versions and EOL, security
severity intake and response targets, recovery and rollback, dependency update
and license-change policy, threat model, and a non-personal support route. Every
release artifact contains the applicable records, including FT88's
data-handling inventory.

Sources: `RR:A-08`, `RR:A-09`, `RR:A-11`, `RR:A-12`; `RC:H-06`, `RC:H-07`.

**FT88 (HIGH on the bank track) — minimal subprocess data exposure.** Agents
and gates launch from separate documented environment passlists with explicit
opt-in additions. Prompt text travels through stdin or a mode-0600 file, never
argv. Durable state uses an objective identifier; commit subjects, terminal
summaries, and structured output reject or redact control characters and
sensitive text.

Ship a data-handling inventory for every repository-controlled prompt,
environment, file, log, network, cache, and retention path. Sentinel contracts
prove denied variables do not reach default subprocesses and prompt content
cannot leak into process listings, commits, or structured output.

Sources: `RR:C-08`; `RC:H-01`.

**FT87 (HIGH on the bank track; MEDIUM otherwise) — bounded, explicit
network/resource and CLI behavior.** One timeout, size, cancellation, and
iteration-cap policy governs agent and gate execution, repair, Git refresh,
model discovery, guard startup, reads, and large output. `BENCH_OFFLINE=1`
prevents every Bench-initiated network attempt; repair is explicit,
independently manifest-pinned, size-bounded, and atomically promoted; Git
refresh is explicit and noninteractive; model providers run concurrently with
bounded reads. Default outline output is a bounded summary with counts and
truncation metadata.

Centralized argument parsing rejects trailing garbage, resolves spec slugs
from the repository root, supports leading-dash and directory-scoped commit
paths through a conventional grammar, and treats help as success. Capability
skips are explicit evidence, not silent passes, and security-test deadlines
are independent from the subprocesses they bound. Complete package metadata
and one user-facing Bench/redbench identity ship with the bounded repair
policy; repair cleanup is concurrency-safe and Git refresh failures are
observable.

Sources: `RR:A-14`, `RR:C-06`, `RR:C-07`, `RR:C-09`, `RR:C-10`, `RR:C-11`,
`RR:C-12`; `RC:H-04`, `RC:M-02`.

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
action. Start with `/bench-shape-idea` to settle instruction-file selection,
inferred versus confirmed facts, conflict/rollback semantics, and the seam
between deterministic installation and prompt-driven configuration.

Sources: `RR:A-02`, `RR:A-03`, `RR:A-04`, `RR:S-09`.

**FT84 (MEDIUM) — transactional managed-asset lifecycle.** Link and relink
stage and preflight the complete write set, sync where durability matters, then
atomically promote or roll back. Upgrade and downgrade reconcile the old and
new manifests: removed clean assets leave, modified assets remain owned until
explicitly resolved, and no stale skill becomes active-but-unowned. Settings
and Git hooks compose or fail before any partial write.

Unlink and shim removal verify ownership markers, return a nonzero partial
status when residuals remain, and emit a machine-readable residual list. They
never recommend raw deletion of an executable Bench does not own. Closure
covers upgrade, downgrade, I/O failure, modified/stale assets, and repeated
link/unlink matrices.

Sources: `RR:A-05`, `RR:A-07`, `RR:A-10`, `RR:A-13`; `RC:M-04`.

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

**FT58 (MEDIUM, evidence supplied) — identity-safe intent reclamation and
private pool roots.** Replace age-only stale takeover and unconditional release
with an identity-safe lock protocol: a live owner is never aged out, competing
reclaimers serialize, and release cannot remove a successor's lock. Permission
failures on Bench-selected pool roots propagate; non-owned or symlinked roots
are rejected and modes are revalidated after creation.

Closure covers a live old PID, two-reclaimer race, successor replacement,
permissive pre-existing directory, chmod failure, symlink root, and crash-safe
re-entry. Current live-owner and successor-release defects make this row
actionable.

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

**FT89 (MEDIUM) — guidance coherence and current-state documentation.** Make
roadmap maintenance consume one full schema snapshot; make every documented
CLI example executable; parse and validate real YAML frontmatter; derive the
skills index and inventories from one implementation; embed design-it-twice
briefs in complete delegation charges; and admit reviewer-approved assessment
findings as a legal synthesis origin. Use the canonical iteration-cap line
definition and only recommend shifts that meet the routing contract.

Clarify shape termination and the no-design-source branch, remove stale paths
and inventory omissions, retire closed decision maps and obsolete historical
reports, dogfood first-party authoring guidance, and pin normative external
references. Rewrite ADRs and README claims to the behavior proved by artifact
contracts, including the actual canary phase selection and npm prepare shape.

Sources: `RR:S-06`, `RR:S-07`, `RR:S-08`, `RR:S-10`, `RR:S-11`, `RR:S-12`,
`RR:S-13`, `RR:S-14`, `RR:S-15`, `RR:S-16`, `RR:S-17`, `RR:S-18`; `RC:M-05`.

## Release and bank reassessment gate

A green source-tree gate is necessary but not sufficient. Reassessment attaches
to one immutable version and its generated manifest after:

1. FT77 through FT82 have executable regression contracts and are closed.
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

1. `/bench-implement-spec` — build FT77 from
   `specs/ownership-safe-worktree-cleanup.md` in a fresh mid-tier session; the
   three dependent stories keep their approved per-story lines.
2. `/bench-shape-idea` — specify FT78, oracle-bound gate verdicts.
3. `/bench-shape-idea` — specify FT79, lossless shift recovery.
