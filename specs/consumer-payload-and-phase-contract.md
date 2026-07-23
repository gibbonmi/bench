# Least-privilege consumer payload and one coherent phase contract

Status: implemented

<!-- Map provenance: no `decisions/` map backs this spec. It was compiled under the
reviewer-directed batch-drain override in `/bench-write-spec`'s entry contract, from the
FT85 roadmap row (decided scope) and its Sources line — `RR:S-01`..`RR:S-05` in
ASSESSMENT.md and `RC:M-03` in COMPLIANCE_ASSESSMENT.md. Every decision this spec had to
default rather than inherit is marked **[defaulted — veto surface]** and is open to
post-hoc veto; there are six such markers. A top-tier falsification pass reviewed this
draft against the six findings and returned recommend; its findings are folded in. -->

## Problem

Two unrelated symptoms share one root cause: Bench ships and states more than one source
for the same fact.

**Payload.** A linked consumer repo receives the complete `.agents/` tree, so it also
receives `bench-assess`, `bench-update-kit`, and `craft-synthesis` — kit-maintenance
surfaces that can replace the consumer's own root assessment, pull upstream HEAD into
their repo, and read their README/CHANGELOG as Bench provenance. That is unnecessary
privilege and confusing product scope. The payload itself is derived three separate times
— `buildLinkPlan` in `internal/adopt`, `files[]` in `package.json`, and
`scripts/wrapper-assets.json` — so any exclusion added to one silently misses the others.
Consumers also have no upgrade route once the maintenance command is withheld.

**Phase contract.** Four workflow guarantees are stated twice and disagree:

- README says a clear idea may skip shaping; `/bench-write-spec` refuses every feature
  without a closed map, and `/bench-shape-idea` says even a clear idea gets one.
- `craft-delegate` allows inline authorship for exactly one source-line;
  `/bench-implement-spec` separately allows any no-spec lighter-path change to stay inline
  and calls that the sole exception.
- `/bench-implement-spec` lands a finishing `bench commit --spec`, while
  `/bench-final-check` claims it owns the landing gate, commit, and status flip. A clean
  reviewed branch can reach final-check with nothing left to commit.
- `/bench-debug` instructs committing a failing repro before a shift, which contradicts
  invariant 4 (commit only on green).

## Solution

One canonical allowlist generates the consumer payload; every derivation reads it or is
graded against it, and kit-only surfaces are named there once and excluded everywhere.
Consumers get `bench upgrade` — narrow, version-pinned, manifest-aware — as the
replacement for the withheld maintenance command.

Each of the four workflow guarantees gets exactly one owner. The other documents point at
that owner instead of restating it, and a conformance anchor plus a canary fixture holds
each agreement in place so it cannot drift back apart.

## User stories

1. As a consumer whose repo Bench is linked into, I want the kit's maintenance surfaces
   (`bench-assess`, `bench-update-kit`, `craft-synthesis`) never installed into my repo, so
   that no agent session in my project can replace my assessments or pull upstream HEAD.
   Line: `gpt-5.6-luna` / medium. The exclusion set is a data edit behind an existing plan
   builder, and the contract tests observe it directly.
2. As a Bench maintainer, I want the consumer payload derived from one canonical allowlist
   that names each asset and its audience, so that adding or withholding an asset is a
   single edit rather than three that drift. Line: `gpt-5.6-luna` / medium. This is a
   refactor to a known seam whose behavior the existing link-lifecycle contract already
   pins.
3. As a Bench maintainer, I want `scripts/wrapper-assets.json` and the `package.json`
   `files[]` payload graded against that same allowlist, so that a forbidden asset cannot
   enter the npm tarball or the wrapper artifact through a surface the link plan does not
   own. Line: `gpt-5.6-luna` / medium. The package-shipped-surface conformance suite
   already parses both inputs; this adds a derived expectation to it.
4. As a consumer running an older linked kit, I want `bench upgrade` to compare my
   manifest's pinned kit version against the installed kit, show me what would change, and
   apply the existing transactional relink, so that I have a supported upgrade route
   without the maintainer command. Line: `gpt-5.6-terra` / medium. The transaction seam
   exists, but the version comparison and refusal posture are new judgment about
   destructive state.
5. As a consumer upgrading across the release that withholds the maintenance surfaces, I
   want the previously-installed kit-only assets removed by that upgrade, so that
   least-privilege actually arrives instead of applying only to fresh installs.
   Line: `gpt-5.6-luna` / medium. Manifest reconciliation already removes clean withdrawn
   assets; this asserts the case rather than building a mechanism.
6. As an agent session in any repo, I want the shaping requirement stated identically in
   README, `/bench-shape-idea`, `/bench-write-spec`, and `/bench-implement-spec`, so that I
   cannot read one document as permission to skip a map another document requires.
   Line: `gpt-5.6-sol` / high. Guidance prose compounds through every session that loads
   it, which is the leverage override in `craft-line`.
7. As an agent session on a harness without write subagents, I want one capability-aware
   delegation policy stated only in `craft-delegate`, so that I know exactly when to author
   inline, when to delegate, and when to stop and hand off — without reconciling two
   conflicting thresholds. Line: `gpt-5.6-sol` / high. Same leverage override; the policy
   governs every build in every linked repo.
8. As an agent finishing a build, I want exactly one phase to own the landing commit and
   the spec status transition, so that the other phase reports honestly on a branch with
   nothing left to commit instead of claiming a commit it did not make.
   Line: `gpt-5.6-sol` / high. Same leverage override, and the change reassigns a phase
   contract rather than editing wording.
9. As an agent debugging, I want the failing repro preserved through a shift without
   committing a red tree, so that invariant 4 stays true and the repro still survives the
   rollback of a failed iteration. Line: `gpt-5.6-sol` / high. Same leverage override, and
   the resolution has to satisfy two rules that currently conflict.
10. As a Bench maintainer, I want each of the four repaired agreements held by a
    conformance anchor with a canary fixture proving it bites, so that a later edit that
    reopens the contradiction turns the gate red. Line: `gpt-5.6-luna` / medium.
    Mechanical work at the anchor seam, with the fixture family and helper already in
    place.

## Implementation decisions

**The canonical allowlist.** One tracked data file is the single source of the consumer
payload: each row names the source path, its file mode, whether it is a tree, and its
audience (`consumer` or `kit-only`). It is embedded into the Go core and read by
`buildLinkPlan`, replacing the hand-listed plan entries and the bare tree walks over
`.agents/commands` and `.agents/skills`. `scripts/wrapper-assets.json` is retired: its
readers point at the canonical file instead, which is why that file carries mode and tree
metadata. The reader set is exactly five, and every one moves in this diff:
`scripts/build-release-evidence.mjs`, `internal/releaseevidence/registry.json` (its inputs
list), `internal/preflight/release_index_test.go`,
`internal/preflight/integration_helpers_test.go`, and
`internal/contract/surface/artifact/artifact_test.go`. Two canary fixtures under
`tests/canary/behavior-owned/` (`wrapper-required-surface-dropped`,
`wrapper-contamination-admitted`) and the `release-evidence-probe-base.txt` baseline name
the retired path and are repointed with them. The release-evidence input list is graded by
the story 3 conformance case, so a missed reader fails the gate rather than reading a
deleted file. Reaching those readers is path-following, not scope growth: the
release-evidence pipeline's own behavior is unchanged and stays out of scope. **[defaulted — veto surface]** The file lives under `.bench/` with the other
tracked kit data files rather than inside `internal/`, so shell and Node consumers can read
it without a Go build.

**`package.json` `files[]` stays npm's literal input** — npm resolves it at pack time and
cannot read our file — but it is no longer an independent derivation of the fact:
conformance derives the expected pattern set (including the exclusions) from the canonical
allowlist and grades `files[]` against it, and the existing npm dry-run pack-shape check
proves no `kit-only` path reaches the tarball.

**Kit-only surfaces** are `bench-assess` and `bench-update-kit` (command and skill) and
`craft-synthesis` (skill). The exclusion applies to every destination the plan writes,
including the `.claude/` adapter mirrors, so a consumer cannot reach a withheld surface
through the adapter path.

**The shipped guide stays one file.** `.bench/BENCH.md` is not forked into kit and consumer
variants. Instead the generated skills index marks kit-only rows from the same allowlist,
and the guide states in one place that those surfaces ship only in the kit repository and
names `bench upgrade` as the consumer route. **[defaulted — veto surface]** The alternative
— generating a second consumer-only guide — was rejected as a new drift source.

**`bench upgrade`** is a new subcommand, deliberately narrow. It resolves the linked repo's
manifest kit version and the installed kit version, prints a TOON plan (from-version,
to-version, added/changed/removed asset counts), and applies the existing transactional
relink. `--check` reports the same plan and applies nothing. An equal version is a
definitive no-op success. A lower installed version is a downgrade and fails closed naming
`--force`, which performs it. It composes `link`'s transaction and manifest reconcile; it
introduces no second write path. The manifest's kit-version header is rewritten inside that
composed transaction, never after it, so an interrupt cannot strand a new tree under an old
version claim. **[defaulted — veto surface]** Adding a subcommand rather
than documenting `npm i -D redbench@<version> && bench link` was chosen because the version
comparison and downgrade refusal are the narrow part the roadmap row asked for, and prose
cannot enforce them.

**One owner per workflow agreement.**

- *Shaping*: every spec has a map behind it — there is no skip. `/bench-write-spec`'s entry
  contract is the single owner, including the reviewer-closed same-session path that
  records the map inline. README, `/bench-shape-idea`, and `/bench-implement-spec` point at
  it and state no threshold of their own. README's "skip shaping when there is no fog"
  sentence is replaced by the inline-map route. S-02 offered this as an either/or against
  making shaping mandatory everywhere; the no-skip branch is not a fresh default because the
  tree's live `/bench-write-spec` entry contract already enforces it, so this aligns the
  other three documents to an owner that exists rather than choosing between open options.
- *Delegation*: `craft-delegate` is the single owner and is capability-aware. It states the
  one-source-line inline allowance, the lighter-path allowance for a change admitted
  without a spec, and the rule for a harness that cannot spawn a write subagent: stop
  before editing and emit one executable resume handoff — never an inline fallback.
  `/bench-implement-spec` and `/bench-debug` delete their restatements and point at it.
  **[defaulted — veto surface]** Folding the lighter-path allowance into the policy (rather
  than deleting it, which would make every one-line no-spec fix require a delegate) keeps
  the proportionality rule in `.bench/BENCH.md` intact.
- *Landing*: `/bench-final-check` owns the landing commit and the `Status: implemented`
  transition. `/bench-implement-spec` ends at its last green build commit, never runs
  `bench commit --spec`, and hands off. Final-check states the honest no-op: a branch that
  arrives with nothing to commit is reported green with the status flip still performed via
  `bench spec implemented <slug>` when the spec has not already flipped.
  **[defaulted — veto surface]** Final-check was chosen as the owner because gate-then-land
  is its whole contract and it runs after review, so review fix commits land under the same
  owner.
- *Red observation*: `/bench-debug` no longer commits a red tree. The repro is committed in
  the project's expected-failure form — a quarantine marker naming the bug — so the gate is
  green with the repro preserved through shift rollback, and the fix's green commit removes
  the marker. A project with no expected-failure form keeps the repro out of the shift and
  runs it by hand; debug states that fallback rather than leaving it to judgment.
  **[defaulted — veto surface]** The alternative — stashing the repro outside the worktree
  and re-applying it — was rejected because it puts recovery state somewhere no gate
  observes.

**Anchors.** Each repaired agreement gets a `require(...)` anchor in `checkWorkflowAnchors`
plus a `tests/canary/workflow-guidance-anchors/<name>/` fixture whose `EXPECT` names the
reintroduced contradiction. The payload contract gets a `tests/canary/package-core-guard/`
fixture for the forbidden-asset case.

## Testing decisions

A good test here drives a real surface and observes what a consumer or a session would see:
a link plan applied to a throwaway repo, a packed tarball's member list, a rendered guide,
or a gate run against a deliberately-broken fixture. No test asserts on the shape of the
allowlist parser or on internal plan structs.

Prior art to follow, not re-invent: `internal/contract/surface/link_lifecycle_test.go` for
link/relink/unlink against throwaway repos; `internal/conformance/package_shipped_surface_test.go`
for `files[]` and packed-shape grading; `internal/conformance/docs_workflow_helpers_test.go`
(`checkWorkflowAnchors`) for prose anchors; `tests/canary/workflow-guidance-anchors/shape-idea-bypass/`
for the fixture shape.

Gate command: the project gate, `.bench/gate.sh`.

### Seam diagram

**Seam A — the payload plan.**

    trigger: bench link | bench upgrade | wrapper tarball build | npm pack
        │
        ▼
    .bench/payload allowlist  ──▶  [ buildLinkPlan / wrapper builder ]  ──▶  written repo tree
    kit tree on disk          ──▶  [ (audience filter)               ]  ──▶  packed tarball members
                                        ◀ tests attach here: apply link into a throwaway
                                          repo and enumerate the written paths; pack and
                                          enumerate tarball members; assert the kit-only
                                          set is absent from both.

**Seam B — the upgrade command.**

    trigger: bench upgrade [--check] [--force]
        │
        ▼
    link-manifest.tsv (#kit version, rows)  ──▶  [ upgrade ]  ──▶  TOON plan on stdout
    installed kit version + allowlist       ──▶  [        ]  ──▶  transactional relink, exit code
                                        ◀ tests attach here: seed a fixture repo with a
                                          manifest at a chosen version, run the command,
                                          assert stdout plan, exit code, and resulting tree.

**Seam C — the workflow anchors.**

    trigger: bench gate (conformance phase) | bench canary
        │
        ▼
    .agents/**/*.md, README.md  ──▶  [ checkWorkflowAnchors ]  ──▶  diagnostics (gate verdict)
    canary fixture overlay      ──▶  [                      ]  ──▶  EXPECT substring match
                                        ◀ tests attach here: run the conformance check over
                                          the real tree (green) and over each fixture that
                                          reintroduces a contradiction (red, named).

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A fresh `bench link` writes no `bench-assess`, `bench-update-kit`, or `craft-synthesis` asset to any destination, including the `.claude/` mirrors | A | new case in `internal/contract/surface/link_lifecycle_test.go` enumerating the linked tree; red before the exclusion lands | The current plan copies both trees wholesale, so the enumeration finds the forbidden paths and fails until the audience filter exists |
| 2 | The link plan's consumer entries equal the allowlist's `consumer` rows — a file under a kit-only subtree is not written even though the allowlisted tree containing it is linked | A | contract case adding a file under the kit-only `.agents/skills/bench-assess/` subtree and asserting link does not write it; red while the plan walks the tree with no audience filter | An unfiltered tree walk links whatever is on disk under `.agents/skills`, so the kit-only file appears and fails the assertion until the plan is allowlist-driven |
| 3 | `package.json` `files[]` and the wrapper asset list carry exactly the allowlist's consumer set, and a packed tarball contains no kit-only member | A | conformance case in `package_shipped_surface_test.go` deriving the expected set from the allowlist; red against today's wholesale `.agents/` entry | Today's `files[]` ships `.agents/` whole, so the derived comparison reports the missing exclusions until they are added |
| 4 | `bench upgrade --check` on a repo whose manifest pins an older version prints the from/to plan and writes nothing; the applying run relinks and rewrites the manifest version | B | new fixture-repo contract case; red because the subcommand does not exist (unknown-subcommand exit) | A missing subcommand cannot print the plan, so the assertion on stdout and on the untouched tree fails until it is built |
| 5 | Upgrading a repo whose manifest still owns kit-only assets removes them from the tree and from the manifest | B | fixture-repo case seeding a pre-exclusion manifest; red until the exclusion reaches the reconcile path | If the allowlist still admits those assets, reconcile keeps them and the post-upgrade enumeration still finds them |
| 6 | README, `/bench-shape-idea`, `/bench-write-spec`, and `/bench-implement-spec` state one shaping requirement with a single owner | C | anchor requiring README to name the inline-map route and forbidding the skip sentence, plus `tests/canary/workflow-guidance-anchors/readme-shaping-skip/`; red against today's README | Today's README says to skip shaping, so the anchor fires immediately, and the fixture restores that sentence to prove the anchor still bites |
| 7 | `craft-delegate` states the capability-aware policy in full and `/bench-implement-spec` states no inline threshold of its own | C | anchor pair plus `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception/`; red against today's implement-spec text | The "sole inline exception" sentence is present today, so the forbidding anchor fires until it is removed and the pointer replaces it |
| 8 | `/bench-final-check` owns the landing commit and status transition; `/bench-implement-spec` names no `--spec` commit | C | anchor pair plus `tests/canary/workflow-guidance-anchors/implement-spec-landing-commit/`; red against today's implement-spec text | `bench commit -m "<msg>" --spec <slug>` appears in implement-spec today, so the forbidding anchor fires until the phase contract moves |
| 9 | `/bench-debug` preserves the repro without committing red, and names the no-expected-failure-form fallback | C | anchor replacing the existing `before launching the shift` requirement, plus `tests/canary/workflow-guidance-anchors/debug-red-commit/`; red against today's debug text | The current sentence instructs a red commit, so the new anchor fires until the quarantine rule replaces it |
| 10 | Each new canary fixture goes red with its own EXPECT substring and the unmutated tree stays green | C | `bench canary` over the new fixtures; red for any fixture whose anchor is vacuous | A fixture that stays green proves its anchor does not bite, which is exactly the rot the canary phase exists to catch |
| 1 | The shipped guide and its generated skills index name only surfaces the consumer actually has | A | conformance anchor requiring the guide's kit-only statement and the index's kit-only marker; red before the guide is updated | The guide currently recommends `/bench-update-kit` unconditionally, so the anchor fires until the consumer route replaces it |
| 2 | A non-regular file (FIFO, device, socket) inside an allowlisted tree is rejected before it is read, rather than blocking the plan | A | contract case placing a FIFO under an allowlisted tree and asserting link completes with the FIFO refused and named | A tree-typed allowlist row still walks its directory, so an unguarded plan builder blocks on the FIFO open and the case hangs or errors until the rejection lands |

### Edge inventory

- **absent vs present-but-empty file** — row: `bench upgrade` on a repo with no manifest
  fails closed naming `bench link`; an empty manifest is a distinct malformed-state error.
  Both asserted in the story 4 fixture case.
- **hand-edited file without trailing newline** — row: a manifest whose last row lacks a
  newline still parses and upgrades, asserted in the story 5 fixture case.
- **re-run idempotency** — row: a second `bench upgrade` at the same version is a
  definitive no-op success that rewrites nothing, asserted in the story 4 fixture case.
- **paths with spaces or glob characters** — row: an allowlist tree containing a
  space-bearing path links and unlinks intact, asserted in the story 2 contract case.
- **required tool missing from PATH / invocation through a symlink / by-path linked-repo
  CLI** — row: `bench upgrade` reached through the linked repo's `.bench/bin` launcher and
  through a symlink resolves the same implementation, asserted in the story 4 fixture case.
- **cwd deeper than the repo root** — row: `bench upgrade` from a subdirectory resolves the
  repo root, asserted in the story 4 fixture case.
- **destructive state / plan-apply drift** — row: a downgrade fails closed naming `--force`
  and leaves the tree untouched; a modified consumer-owned asset survives the upgrade with
  the partial-result posture the relink contract already defines. Story 4 and 5 cases.
- **interrupt mid-operation** — **Won't handle** — `bench upgrade` composes the existing
  link transaction, whose rollback-on-interrupt behavior the link-lifecycle contract
  already owns; a second copy of that assertion would duplicate covered knowledge.
- **non-TTY stdin on a prompting command** — **Won't handle** — `bench upgrade` never
  prompts; `--check` and `--force` carry every decision, so there is no TTY contract.
- **control bytes in git-sourced text** — **Won't handle** — the allowlist is a tracked
  in-repo data file and the payload paths come from it, not from git-sourced text.
- **special files in script-discovery paths** — row: tree-typed allowlist rows still walk
  their directory, so a non-regular file inside an allowlisted tree is rejected before it is
  read. Asserted in the second story 2 case.
- **unquoted multi-word arguments** — **Won't handle** — `upgrade` takes only fixed flags,
  no free-text argument.
- **compatibility probe (older consumer trees)** — row: story 5 is exactly that probe — a
  manifest written by a pre-exclusion kit upgrades cleanly.
- **amputated caller** — row: `scripts/wrapper-assets.json` is deleted, and story 3's
  conformance case asserts every former reader now resolves the canonical allowlist, so a
  missed caller fails the gate rather than silently reading a deleted file.

## Out of scope

- **Per-consumer payload profiles** (a consumer choosing which optional skills to install)
  — a separate capability: it needs a consumer-side configuration surface, its precedence
  rules against the kit default, and doctor reporting for a drifted selection. Estimate:
  6 edits, 4 gate runs.
- **Signed or checksum-verified payload assets** — a separate capability spanning the
  release-evidence pipeline rather than the link plan, with its own key handling and
  offline posture. Estimate: 8 edits, 5 gate runs.
- **RC:M-05's ADR and decision-map drift** (ADR 0003/0004 amendments, retiring closed maps)
  — named by a different finding and already carried on its own roadmap row; folding it
  here would widen the diff without sharing a seam. Estimate: 5 edits, 3 gate runs.
