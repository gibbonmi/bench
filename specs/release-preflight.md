# FT82 — authoritative release preflight

Status: implemented

Decision map: `decisions/release-preflight.md`

## Problem

Bench publication is assembled from independent workflow steps that do not prove the
same source commit passed the full gate, race tests, vet, vulnerability policy,
artifact inspection, and an installed-runtime smoke. The tag workflow neither checks
release identity and ancestry nor retains one machine-readable verdict. A deleted,
bypassed, or drifting phase can therefore leave publication looking green without one
authoritative preflight having authorized it.

## Solution

Add one Go-core-owned `release-preflight` plumbing command with `verify` and `publish`
modes. It runs a single named phase registry, writes atomic JSON records and a terminal
v1 manifest, and fails closed on any unreadable input, missing required tool, red
phase, interruption, or release-identity mismatch. Publish is a strict superset of
verify and adds exact tag/version/binary/changelog equality plus release-line
ancestry.

Pull-request and default-branch verification, the native runner matrix, and tag
publication become thin callers of that command and its existing host smoke seam.
The gate proves the registry and workflow wiring, while the first real tag remains the
manual proof for GitHub runner execution, artifact upload, and npm publication.

## User stories

1. As a contributor, I want one `verify` preflight to run the full gate, race tests,
   vet, vulnerability policy, artifact build and inspection, and clean-room installed
   smoke against one source commit, so that a green verification has one complete
   meaning. Line: `gpt-5.6-terra` / medium. The map fixes the primary seam, but the new
   phase orchestration and evidence lifecycle need the profile's mid-tier
   gate/conformance route.

2. As a release maintainer, I want `publish` to be a strict superset of `verify` and
   reject an invalid tag, unequal release identities, a non-main-line commit, or an
   invalid changelog state before publication, so that an immutable npm release is
   authorized only from the intended source and version. Line: `gpt-5.6-terra` /
   medium. The equality and ancestry contracts are exact, but their release-critical
   semantics justify the profile's mid-tier gate/conformance route.

3. As a security reviewer, I want vulnerability findings to block the preflight
   unless a tracked, matching, unexpired exception exists, and I want stale exceptions
   to fail, so that exceptions remain reviewable loans rather than permanent silent
   waivers. Line: `gpt-5.6-terra` / medium. The schema is precise, while scanner-output
   parsing and time-dependent expiry are correctness-sensitive gate work.

4. As an operator inspecting a failed or successful preflight, I want one complete,
   atomic directory of per-phase records and a terminal manifest, so that automation
   and humans can identify exactly what ran without mistaking partial state for a
   verdict. Line: `gpt-5.6-terra` / medium. The record schema is fixed here, but
   interruption and atomic replacement require behavior the gate must prove.

5. As a release maintainer, I want pull requests, default-branch pushes, the native
   runner matrix, and tag publication to call the same preflight-owned phases and
   retain their records at the maximum allowed duration, so that workflow choice
   cannot bypass release policy. Line: `gpt-5.6-terra` / medium. Workflow execution is
   partly gate-blind, so the weak local coverage bumps this story to mid.

6. As the reviewer, I want the project gate and biting canaries to reject a deleted,
   duplicated, reordered, or bypassed required phase and every release-policy failure
   class, so that preflight completeness is demonstrated rather than inferred from
   implementation prose. Line: `gpt-5.6-terra` / medium. This changes the oracle and
   its omission proof, matching the profile's cached gate/conformance route.

The implementation is one integration sequence because the manifest, workflows, and
oracle all consume the same registry. It may land green in three commits in dependency
order: core and records, workflow wiring, then conformance and canaries. All stories
share the same `gpt-5.6-terra` / medium line.

## Implementation decisions

- **Command contract.** `bench release-preflight --mode verify|publish [--phase
  <name>]` is plumbing, not an AXI query surface. It takes no interactive input,
  resolves the repository root from any descendant cwd, emits progress and phase
  summaries on stderr, writes evidence under `dist/preflight/`, and exits 0 only for
  a complete green run, 1 for a red or interrupted preflight, and 2 for invalid
  usage. `publish` requires the exact tag ref through `GITHUB_REF` or a test-injected
  equivalent; a missing or non-tag ref is red. The thin repo entry script forwards
  argv unchanged to this compiled command.
- **Structured failure output.** Command-level failures emit one compact JSON object
  on stderr with exactly `kind` and `message`; `kind` is one of `usage`, `input`,
  `tool`, `phase`, `identity`, `interrupted`, or `evidence`, and `message` is an
  actionable single line with no dependency stack trace. Phase subprocess output
  remains stderr diagnostic text. Once evidence staging exists, the same stable kind
  and message are also stored in that phase's record, so automation need not parse
  logs. Successful runs write no stdout data.
- **One phase registry.** The Go package owns an ordered registry with verify phases
  `gate`, `race`, `vet`, `vulnerability`, `artifacts`, and `smoke`, followed in publish
  mode by `identity`, `ancestry`, and `changelog`. Publish's selected names are
  mechanically derived as the complete verify set plus the three publish-only names.
  Tests, the manifest, single-phase matrix calls, and structural conformance consume
  this registry; no workflow or test carries a second executable phase list.
- **Phase execution.** `gate`, `race`, `vet`, and `vulnerability` may execute
  concurrently. `artifacts` is serial and calls the existing exact-artifact builder;
  `smoke` depends on its promoted tarballs and calls the existing clean-room installed
  smoke. Publish-only identity checks complete before any npm publish step. A failed
  prerequisite leaves dependants `not_run`; independent phases finish so the record is
  diagnostic. A missing `govulncheck`, Go, Node, npm, Git, or required repo input is
  red, never skipped.
- **Focused phase contract.** `--phase smoke` is the runner-matrix entry and is legal
  only when the exact artifact set is present; it records that one phase but does not
  claim a complete preflight. It exits by the phase verdict and marks the manifest
  `scope: focused`. Only a full `scope: preflight` manifest can authorize verify or
  publish. Other phase names are accepted for deterministic fixture tests, not as a
  workflow substitute for a full preflight.
- **Evidence schema.** Each `<phase>.json` is compact JSON with exactly
  `schema_version`, `phase`, `mode`, `status`, `exit_code`, and `error`.
  `schema_version` is 1; `status` is `green`, `red`, `not_run`, or `interrupted`;
  `exit_code` is an integer for an attempted phase and `null` otherwise; `error` is
  `null` or the exact `{kind,message}` object from the structured failure contract.
  `manifest.json` has exactly
  `schema_version`, `mode`, `scope`, `status`, `identity`, and `phases`. `status` uses
  the same vocabulary. `identity` has exactly `tag`, `package_version`,
  `source_commit`, `binary_version`, `changelog_heading`, and `toolchain`; fields that
  do not apply in verify mode are JSON `null`. `phases` is the registry-ordered array
  of `{name,status,exit_code}` records. JSON field names and phase order are stable
  contract; log bytes and elapsed timing are deliberately not v1 manifest fields.
- **Atomic evidence lifecycle.** A run writes only into a private sibling staging
  directory. On every handled terminal state, it writes all expected phase records
  (including `not_run`), writes the manifest last, fsyncs files and directory, and
  atomically replaces `dist/preflight/`. SIGINT/SIGTERM cancel process groups and
  produce an `interrupted` terminal record when the handler can finish. An unhandled
  kill leaves the previous complete directory intact or no promoted directory; a
  staging directory is never a verdict. Re-running replaces the prior result and
  cannot mix records from two source commits.
- **Release identity.** Publish accepts only exact `vMAJOR.MINOR.PATCH` tags. The tag
  without `v`, the wrapper package version, the built binary's stamped version, and
  the exact `## vMAJOR.MINOR.PATCH (YYYY-MM-DD)` changelog heading must agree. The
  manifest source commit is the checked-out `HEAD`, and the tag must resolve exactly
  to it. The Go `toolchain` directive must be an exact patch version, and the actual
  `go env GOVERSION` used by the run must equal it. The manifest records both the
  decided identity values and the terminal phase verdicts.
- **Changelog state.** A release renames the rolling `## Unreleased` heading to the
  exact versioned heading. Publish rejects a missing or duplicate matching heading,
  an `## Unreleased` heading with non-heading content stranded before the release
  heading, and a heading whose version differs from the tag. Historical release
  sections below the matching heading are not reinterpreted.
- **Release-line ancestry.** Publish fetches or otherwise requires a resolvable
  `origin/main` and proves the tagged `HEAD` is its ancestor or equal. Missing remote
  state, shallow history that cannot prove ancestry, or any Git error is red. No
  `release/*` grammar is inferred.
- **Vulnerability policy.** The phase consumes `govulncheck -json` through one parser
  and compares reachable finding IDs to `scripts/vuln-exceptions.json`. Absence means
  no exceptions. A present file must be a JSON array of objects with exactly nonempty
  `id` and `reason` strings plus an RFC 3339 full-date `expires` value. Duplicate IDs,
  unknown fields, malformed/empty files, expired entries, entries matching no current
  finding, and uncovered findings are red. Expiry is evaluated in UTC through an
  injected clock used by tests. The scanner database and scanner binary remain
  external inputs; the exception file is the repository-owned policy.
- **Workflow ownership.** Pull requests and pushes to `main` call the full verify
  preflight once. The existing matrix derives rows from the canonical platform matrix,
  downloads the exact artifacts produced by verify, and calls `--phase smoke` on each
  host. Tag publication calls the full publish preflight, waits for the native smoke
  matrix, and only then runs the existing platform-first/wrapper-last npm publication.
  Both workflows upload `dist/preflight/`; verify uses the platform default retention,
  and publish sets the maximum retention allowed by that repository at run time rather
  than encoding a literal day count.
- **Oracle attachment.** Package tests drive the preflight command with fake phase
  executables and fixture repos. Root conformance proves the exact registry relation,
  workflow triggers, preflight calls, matrix dependency, record upload, and publish
  ordering. Canary fixtures independently omit or bypass each required phase class and
  plant each release-policy failure. Independently authored expected phase names are
  allowed only as ADR 0006's demonstrated omission oracle; production registry,
  manifests, workflows, parsers, and derived counts remain single-sourced.
- **Deep and thin modules.** The preflight package owns phase selection, execution,
  cancellation, identity, exception parsing, and evidence promotion. Existing artifact
  build and smoke scripts remain deep owners of their artifact behaviors. The entry
  script and workflow YAML are thin callers. Tests attach to command exits, exact JSON,
  artifact behavior, and workflow structure rather than private orchestration helpers.

## Testing decisions

- A good test runs the built command against throwaway repositories and controlled
  executables, then observes exit code, exact record/manifest JSON, artifact state,
  and Git ancestry. It does not assert private functions or duplicate the production
  phase registry outside the named ADR 0006 omission oracle.
- The primary seam follows the existing `gate-phases` command and runner tests. Exact
  artifact build and installed smoke reuse `build-artifacts.sh` and
  `smoke-artifacts.sh`; workflow structure follows the native-runtime and release
  conformance tests.
- The feature must pass the project gate: `.bench/gate.sh`.
- Pre-implementation probes were run on 2026-07-16 and are named `R1` through `R7`:
  - `R1` `bin/bench.sh release-preflight --mode verify` exited 2 with unknown
    subcommand.
  - `R2` a workflow search exited 1 because neither workflow references
    `release-preflight`.
  - `R3` an exact-heading probe exited 1 because `CHANGELOG.md` has no
    `## v0.2.0 (date)` release heading.
  - `R4` an exception-schema-and-consumer probe exited 1 because no tracked schema or
    parser exists.
  - `R5` a release-workflow smoke-matrix probe exited 1 because tag publication has
    neither the shared smoke nor a runner matrix.
  - `R6` a gate search exited 1 because no preflight conformance, contract, or canary
    coverage exists.
  - `R7` an enforcement search exited 1 because no executable check pins the exact Go
    patch toolchain.

### Seam diagrams

Seam 1 — authoritative preflight command:

    trigger: PR/push verify, tag publish, or fixture invokes the plumbing command
        │
        ▼
    mode + repo + tools ──▶ [ Go preflight registry and runner ] ──▶ exit 0/1/2
    tag/ref + policy     ──▶ [ identity and exception owners  ] ──▶ phase verdicts
                                      │
                                      ▼
                           dist/preflight staged then promoted
                    ◀ tests attach here: drive the built command with fixture
                      repos/tools and assert exits plus exact terminal evidence

Seam 2 — exact evidence directory:

    trigger: a phase finishes, fails, or receives an interrupt
        │
        ▼
    ordered results + identity ──▶ [ atomic evidence writer ] ──▶ <phase>.json
    prior complete directory   ──▶ [ stage, fsync, replace   ] ──▶ manifest.json
                                           │
                                           ▼
                                 complete old or new verdict
                        ◀ tests attach here: inject failures/signals and observe
                          that no partial or mixed directory becomes authoritative

Seam 3 — existing artifact build and host smoke:

    trigger: artifacts phase, local smoke dependency, or native matrix row
        │
        ▼
    source commit + matrix ──▶ [ exact artifact builder ] ──▶ wrapper/platform tgz
    tarballs + host        ──▶ [ installed host smoke   ] ──▶ executable verdict
                                       │
                                       ▼
                             phase record and exit code
                    ◀ tests attach here: reuse exact-tarball inspection and the
                      installed smoke while poisoning incomplete/wrong artifacts

Seam 4 — workflow structure:

    trigger: root conformance reads both checked-in workflows
        │
        ▼
    phase registry + YAML ──▶ [ workflow structure checker ] ──▶ green/red verdict
    platform matrix       ──▶ [ call/dependency assertions ] ──▶ cited omission
                                      │
                                      ▼
                              project gate conformance phase
                    ◀ tests attach here: mutate calls, needs edges, upload steps,
                      or publish ordering and require the project gate to go red

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | A conformant verify fixture runs all six verify phases, exits 0, and emits a `green` full-preflight manifest with six green records | preflight command plus exact evidence directory | R1: the command exits 2 because the surface is absent | An always-red aggregator, success-blind recorder, or command that never promotes a green terminal verdict cannot satisfy the positive baseline |
| 1 | Full verify runs registry phases `gate`, `race`, `vet`, `vulnerability`, `artifacts`, and `smoke` against one resolved HEAD | preflight command | R1: the command exits 2 because the surface is absent | A partial script, always-green stub, or workflow-only implementation cannot return the required phase records from the command seam |
| 1 | Artifact build is promoted before smoke; a failed prerequisite marks smoke `not_run`, while independent phases still record terminal results | preflight command plus artifact seam | R1: no command exists to express dependencies or records | A sequential happy-path port or direct smoke call fails dependency-state and evidence assertions under injected artifact failure |
| edge of 1 | Missing Go, Node, npm, Git, govulncheck, repo input, or a non-executable required tool is red with a structured actionable error and never skipped | preflight command | R1: the absent command cannot distinguish required-tool failure | PATH-controlled fixtures catch optional-tool treatment, unstructured errors, and dependency stack leakage |
| edge of 1 | Invocation from a descendant cwd, through the thin script, and from a repository path containing spaces/globs resolves one source commit and preserves argv | preflight command | R1: the new entry surface is absent | Process fixtures observe root, argv, record path, and exit, catching shell expansion or cwd assumptions |
| 2 | Publish derives all six verify phases plus exactly `identity`, `ancestry`, and `changelog`; no publish step can run before all are green | preflight command plus workflow structure | R2: workflows contain no shared preflight call | A separate publish list, omitted verify phase, or early npm step fails the strict-superset and dependency assertions |
| 2 | A conformant publish fixture with matching identity and main-line ancestry runs all nine phases, exits 0, and emits a `green` full-preflight manifest | preflight command plus exact evidence directory | R1: no publish command exists | An always-red release guard or publish mode that validates but cannot authorize a valid release fails the positive baseline |
| 2 | Only exact `vMAJOR.MINOR.PATCH` is accepted; pre-release/build-metadata tags, missing refs, and tag/HEAD mismatches are red | preflight command | R1: no publish parser exists | Table fixtures make a permissive SemVer parser or branch-ref fallback fail at the public exit and manifest |
| 2 | Tag, package version, binary version, changelog heading, source commit, and exact Go toolchain agree | preflight command and exact manifest | R3 and R7: release heading and toolchain enforcement probes exit 1 | Mutating each identity leg independently makes the equality phase red and names the mismatched manifest field |
| 2 | Tagged HEAD is an ancestor of or equal to resolvable `origin/main`; unrelated, missing, and insufficient shallow history fail closed | preflight command | R1: no ancestry phase exists | Fixture Git graphs reject a string/branch-name check and prove the actual commit relation |
| edge of 2 | A duplicate/missing matching heading or stranded content under `## Unreleased` is red; historical sections remain accepted | preflight command | R3: the current changelog has no matching version heading | Changelog byte fixtures catch a heading-only check that ignores ambiguous or stranded release content |
| 3 | Absent exceptions means none; valid matching unexpired entries cover findings by ID | preflight command | R4: schema-and-consumer probe exits 1 | Scanner JSON fixtures require the policy parser to participate in the verdict rather than merely validating JSON syntax |
| 3 | Malformed/empty files, unknown fields, duplicates, empty reasons, invalid dates, expired entries, unused entries, and uncovered findings are red | preflight command | R4: no exception owner exists | One mutation per class defeats an always-green parser and prevents stale suppressions from passing silently |
| edge of 3 | Expiry uses injected UTC date, including exact-expiry-day and next-day boundaries; scanner errors or malformed event streams are red | preflight command | R4: no scanner parser or clock seam exists | Boundary fixtures catch local-time comparisons and parsers that treat incomplete scanner output as no findings |
| 4 | Every full run emits exact-schema records for all selected phases and one registry-ordered manifest with the decided identity | exact evidence directory | R1: no records or manifest are emitted | Exact JSON decoding and set/order comparison reject logs disguised as records, missing phases, and hardcoded counts |
| 4 | Red runs promote a complete terminal record with `red`/`not_run`; focused runs are marked `scope: focused` and cannot authorize publication | exact evidence directory | R1: no evidence lifecycle exists | Failure and focused-mode fixtures catch success-only evidence and a one-phase manifest masquerading as full preflight |
| edge of 4 | SIGINT/SIGTERM cancels child process groups and yields `interrupted`; an unhandled kill, write failure, or promotion failure preserves the prior complete verdict or none | exact evidence directory | R1: no command or atomic writer exists | Signal and filesystem fault injection expose direct-to-output writes, orphaned children, and mixed-run directories |
| edge of 4 | A second run safely replaces the first without stale records, including after a prior red or abandoned stage | exact evidence directory | R1: no rerunnable evidence surface exists | Two-run fixtures compare commit identity and phase sets, catching append or partial-cleanup implementations |
| 5 | PR and main-push workflows call full verify once, upload its records, and feed its exact artifacts to every canonical native matrix row | workflow structure plus artifact seam | R2: neither workflow calls preflight | Root conformance rejects workflow-local reimplementations, missing rows, and artifacts not bound to verify output |
| 5 | Tag workflow calls full publish, waits for every native smoke row, uploads evidence, then publishes platform packages before the wrapper | workflow structure | R5: release workflow has no smoke matrix | Dependency and step-order assertions catch a green preflight job that publication can bypass |
| 5 | Verify uses default retention and publish requests the repository's maximum allowed retention without a literal day constant | workflow structure | R2: no preflight artifact upload exists | YAML structure rejects a hardcoded visibility-dependent cap while preserving the closed “maximum allowed” decision |
| 6 | Root conformance compares verify and publish registry sets, workflow calls, needs edges, uploads, and publication ordering | project gate | R6: no preflight gate coverage exists | Deleting or bypassing a phase or edge makes the gate red even when remaining commands succeed |
| 6 | Canary mutations independently omit each verify phase class (`gate`, Go analysis, vulnerability, artifact, smoke) and each publish-only class (`identity`, `ancestry`, `changelog`) | project gate canary | R6: no preflight canary exists | Class-granular omissions demonstrate the completeness oracle without duplicating every implementation detail |
| 6 | Runtime fixtures plant tag/version, changelog, ancestry, exception, toolchain, interruption, and partial-promotion failures and require cited red output | project gate contracts and canary | R6: no preflight runtime coverage exists | The cheapest always-green registry cannot survive one planted failure per Handoff red class |

Cheapest wrong implementations checked against the map: an always-green command fails
every exact-record and failure-injection row; an always-red aggregator fails both
conformant-green rows; a sequential shell port fails the
registry/superset and cancellation rows; workflow-only duplication fails the command
and one-source rows; a JSON-only exception validator fails scanner finding and stale
entry rows; and a focused smoke manifest fails the `scope: preflight` authorization
row.

### Edge inventory

- **Error path:** phase red, process start failure, scanner error, identity mismatch,
  Git error, record-write failure, and upload/publish bypass are coverage rows.
- **Empty/absent input:** absent versus empty exception files, missing tag/ref,
  missing release heading, empty scanner stream, and missing prior evidence are
  coverage rows.
- **Boundary values:** SemVer numeric components, expiry day, manifest phase count,
  exact toolchain patch, tag equal to `origin/main`, and maximum retention semantics
  are coverage rows.
- **Malformed input:** tag, changelog, scanner JSON events, exception JSON/schema,
  `go.mod`, workflow structure, and prior output residue are coverage rows.
- **Interrupted/partial state:** SIGINT/SIGTERM, child process groups, artifact
  prerequisite failure, evidence write/promotion failure, and unhandled kill are
  coverage rows.
- **Re-run idempotency:** green-after-green, green-after-red, red-after-green, and a
  run after abandoned staging are coverage rows.
- **Hostile environment:** paths with spaces/globs, descendant cwd, symlinked entry,
  missing/non-executable tools, shallow/missing remote history, control bytes in
  Git-sourced identity, and filesystem errors are coverage rows. Manifest strings
  containing ESC or BEL are rejected before JSON promotion rather than redacted, so
  the evidence remains an exact identity record.
- **No trailing newline:** changelog, exception JSON, and `go.mod` fixtures cover valid
  final lines without a newline; parsers must not require a terminator.
- **Special files:** exception, changelog, `go.mod`, package metadata, and output
  targets must be regular files/directories as appropriate; FIFOs, devices, sockets,
  and symlink escapes are rejected before reads or writes can block or leave the repo.
- **Every shipped invocation surface:** this command is repo-only release plumbing.
  The real kit CLI and thin repo entry are in scope; linked consumer CLIs, hooks, and
  harness adapters are **Won't handle** because they do not own repository release
  publication and adding them would expose a destructive irrelevant surface.
- **Actual GitHub Actions execution and artifact service retention:** **Won't handle**
  in the local gate because hosted runners and repository visibility are external;
  workflow structure is gate-covered and the first real tag manually verifies the
  matrix, uploads, maximum retention, and publish dependency.
- **Live npm registry behavior:** **Won't handle** pre-publish because publication is
  immutable and FT83 owns staged/live-registry verification; the in-scope caller still
  exercises exact local tarballs and publication ordering.
- **govulncheck database availability and future finding drift:** **Won't handle** as
  deterministic local data; scanner failure is red, recorded fixture streams cover
  policy semantics, and real CI supplies the current database.

## Out of scope

- **FT83 governed release bundle:** deterministic full manifest, SBOM, checksums,
  reproducibility comparison, staged/dist-tag publication, rollback records, and
  live-registry verification form a separate offline-verifiability and publication
  capability. Estimated later at 12–18 edits and 3 gate runs.
- **Release-branch ancestry grammar:** accepting `release/*` lines is a separate
  branching-policy capability to specify when the first such branch exists. Estimated
  later at 4–6 edits and 2 gate runs.
- **FT88 subprocess environment passlists and data inventory:** minimizing and
  documenting every child environment is a separate execution-hardening capability.
  Estimated later at 8–12 edits and 2 gate runs.
- **Durable in-bundle preflight records:** embedding records in release artifacts is
  FT83's manifest capability; FT82 retains CI artifacts only. Estimated later within
  FT83's 12–18 edits and 3 gate runs, not as a separate FT82 cut.
