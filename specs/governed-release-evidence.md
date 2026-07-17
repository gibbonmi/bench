# FT83 slice 1 — governed release evidence

Status: implemented

Decision map: `decisions/governed-offline-release-bundle.md` (reviewer-closed in
the shaping session; this same-session compilation remains open to reviewer veto).

## Problem

Bench can build and smoke wrapper and platform packages, and its authoritative
release preflight can retain a verdict, but the released bytes do not yet carry the
governance and component evidence needed to audit them offline. There is no single
requirements registry for release profiles, no deterministic release index binding
the evidence to inspected artifact bytes, and no fail-closed way to reserve evidence
owned by FT71, FT87, and FT88 without duplicating those producers.

## Solution

Add the governance and deterministic-evidence core as the first of FT83's three
dependent slices. Canonical governance records, third-party notices, an SPDX-JSON
SBOM, and typed component manifests ship in every wrapper and platform package. The
existing preflight remains the only verdict owner: one `FinalizeEvidence` facade
validates typed producer records, computes observed-byte digests, writes a canonical
`release-index.json`, derives `SHA256SUMS`, and promotes complete green or trustworthy
red evidence old-or-new.

A single requirements registry defines the `public` and `bank` profiles. Verify mode
may report future-owner records as pending, while publish mode requires an explicit
profile and stays red until all required records exist. This slice therefore adds
real machinery without claiming that FT83, the public-release track, or the bank
track is complete.

## User stories

1. As a package consumer, I want every independently published npm artifact to carry
   the license, current governance policies, third-party notices, SPDX-JSON SBOM, and
   a component manifest, so that I can inspect its obligations and contents without
   network access. Line: `gpt-5.6-luna` / high. The package seam and exact inventories
   are known and gate-observable, while the reviewer explicitly requested maximum
   effort on Luna for implementation.

2. As a release maintainer, I want one canonical registry to define stable evidence
   keys, owners, schemas, applicable profiles, and requiredness, so that release
   authorization cannot silently omit records or duplicate another roadmap row's
   content. Line: `gpt-5.6-luna` / high. The closed map fixes the registry semantics
   and fixtures can observe every rule, while the reviewer explicitly requested
   maximum effort on Luna for implementation.

3. As a release maintainer, I want verify and publish preflight to finalize one
   deterministic index from observed artifacts and evidence, so that checksums,
   identities, inventories, and policy relationships have one machine-readable
   source. Line: `gpt-5.6-luna` / high. The facade and encoding are pre-agreed and the
   gate can compare exact bytes, while the reviewer explicitly requested maximum
   effort on Luna for implementation.

4. As a release reviewer, I want public and bank authorization to fail closed on
   missing, empty, malformed, unsafe, mismatched, or drifting required evidence, so
   that an incomplete bundle cannot become publishable by omission. Line:
   `gpt-5.6-luna` / high. Failure classes are enumerated at two black-box seams, while
   the reviewer explicitly requested maximum effort on Luna for implementation.

5. As an operator diagnosing a release, I want complete green or trustworthy red
   evidence to replace the prior generation atomically, while untrustworthy failures
   preserve the prior generation, so that partial or mixed evidence is never
   authoritative. Line: `gpt-5.6-luna` / high. Existing atomic promotion provides a
   known shape and fault injection makes the extension observable, while the reviewer
   explicitly requested maximum effort on Luna for implementation.

6. As the reviewer, I want conformance checks, behavior contracts, and biting
   canaries to reject schema drift, missing package evidence, digest drift, and
   publication-profile bypass, so that the new release rules are enforced from one
   source rather than merely documented. Line: `gpt-5.6-luna` / high. This changes the
   oracle but the required mutations are explicit, while the reviewer explicitly
   requested maximum effort on Luna for implementation.

The six stories form one integration sequence because package manifests are inputs to
the release index and the index is the preflight verdict's evidence. All code-writing
passes use the reviewer-directed `gpt-5.6-luna` / high line.

## Implementation decisions

- **Governance contract.** Release inputs include current supported-version/EOL,
  security-response, dependency/license-change, threat-model, recovery/rollback, and
  support policies. Bench supports the latest minor and the previous minor for 90
  days. GitHub Security Advisories is the private vulnerability route and GitHub
  Issues is the non-personal support route. Reports are acknowledged within three
  business days and triaged within seven; critical issues require mitigation within
  seven days, high within 30, medium within 90, and low is best-effort. These values
  occur in their canonical policy records and are packaged, not re-authored in build
  scripts or validators.
- **One requirement registry.** A repository-owned, versioned registry is the sole
  source for each stable evidence key, owning roadmap row, accepted schema version,
  applicable profiles, and requiredness. `bank` is a strict superset of `public`.
  Public requires FT88's data-handling inventory and FT87's offline/network-control
  record; bank additionally requires FT71's versioned local-event evidence. Producers
  own record content; FT83 validates, hashes, packages, and indexes it. A required
  record cannot be omitted or declared not applicable. A conditional record may use
  `not_applicable` only when the registry permits it and records a reason.
- **Profile behavior.** `release-preflight --mode verify [--phase <name>]` remains
  diagnostic and may record missing future-owner requirements as `pending`. Publish
  mode requires `--profile public|bank`; missing profile is usage exit 2 and any
  unsatisfied requirement is a red release-intent exit 1. Focused runs remain
  diagnostic and cannot authorize publication. Existing hosted publication callers
  select `public` explicitly and remain red until the required owner records exist.
- **Typed deep evidence owner.** `internal/preflight` remains the sole phase
  orchestrator and atomic promoter. Its caller supplies only `RunEvidence` — mode,
  scope, identity, selected profile, and registry-ordered phase results — to
  `FinalizeEvidence(ctx, root, run)`. A deep typed release-evidence module owns safe
  discovery, schema validation, canonical requirements, inventories, SHA-256,
  encoding, drift checks, and evidence assembly. Callers cannot supply paths,
  required sets, hashes, ordering, status, or raw-JSON plugins.
- **Canonical encoding.** Versioned UTF-8 JSON with LF termination is canonical.
  Fixed typed-struct field order and canonically sorted arrays are required; unordered
  maps, locale, ambient environment, and ambient timestamps cannot affect deterministic
  bytes. Unknown fields, duplicate keys or normalized paths, unsupported schema
  versions, and control bytes in identity/path facts are red.
- **Component evidence.** The private-staging artifact builder remains the sole writer
  of package bytes. For the wrapper and each of the four supported platform packages,
  it embeds canonical governance records, license, notices, SPDX-JSON SBOM, and a
  component manifest. The manifest binds component identity plus a sorted internal
  path/mode/size/SHA-256 inventory. It excludes its own digest; the external index
  hashes the final tarball. Package tests independently inspect the exact tarballs.
- **Release index.** `release-index.json` is the sole deterministic release index. In
  this slice it binds source commit, tag/package identity when applicable, rollback
  target, toolchains and flags, dependency/platform/policy inputs, ordered preflight
  results and record digests, profile requirement status, the four-target matrix,
  current wrapper/platform artifact size and SHA-256, and component/SBOM/inventory
  digests. Later slices extend the typed schema with reproducibility and offline
  archive evidence without creating a second index. `SHA256SUMS` is mechanically
  derived from the index.
- **Trustworthy red and atomicity.** Deterministically representable defects produce a
  complete red generation naming every requirement status. Unreadable or unsafe
  inputs, cancellation before a trustworthy record, encoding failure, source drift,
  or promotion failure preserve the prior complete generation. Inputs are revalidated
  immediately before promotion; the existing private-stage, sync, and atomic-replace
  lifecycle stays the only authority path.
- **One-source enforcement.** Producers contain facts, the registry contains
  applicability policy, component manifests inventory package contents, and the index
  binds final bytes. Documentation and conformance checks consume those sources where
  possible. Independent test expectations exist only for named omission mutations
  that demonstrate the gate goes red, following ADR 0006.
- **Release-readiness status.** A green verify run proves this slice's mechanics, not
  publication readiness. Publish remains red while a selected profile lacks required
  FT71/FT87/FT88 evidence. FT83 remains on the roadmap until the reproducible-offline
  artifact and governed-publication slices also ship.

## Testing decisions

- Tests attach to the built release-preflight command and exact promoted evidence,
  plus the existing exact-tarball artifact seam. They drive throwaway repositories,
  controlled package inputs, and filesystem faults; they do not mock private Go
  collaborators or assert internal call counts.
- Artifact tests extend the existing private-staging build, exact inventory, install,
  and failed-promotion contracts. Preflight tests extend the existing command,
  registry, evidence, interruption, and workflow-conformance prior art.
- The feature must pass the project gate: `.bench/gate.sh`.
- Pre-implementation probes run on 2026-07-16 are named `R1` through `R3`:
  - `R1` built all five current npm artifacts and exited 1 because each lacked the
    required governance/evidence set, including component manifest and SPDX SBOM.
  - `R2` `bin/bench.sh release-preflight --mode publish --profile public` exited 2
    and printed usage without a profile option.
  - `R3` a nonempty-file probe exited 1 because the repository lacked current support,
    threat-model, recovery, dependency-policy, and third-party-notice records.

### Seam diagrams

Seam 1 — exact wrapper and platform package bytes:

    trigger: artifact builder stages wrapper plus four platform packages
        │
        ▼
    source + matrix + policy ──▶ [ private-staging artifact builder ] ──▶ five tgz
    dependency metadata      ──▶ [ component evidence producer     ] ──▶ manifests/SBOM
                                             │
                                             ▼
                                  promoted exact artifact set
                    ◀ tests attach here: inspect every tarball's exact entries,
                      modes, schemas, identities, inventories, and observed digests

Seam 2 — authoritative preflight evidence:

    trigger: full verify or profile-selected publish preflight finishes phases
        │
        ▼
    RunEvidence + registry ──▶ [ FinalizeEvidence / typed evidence core ] ──▶ verdict
    artifacts + records    ──▶ [ validate, hash, encode, drift-check   ] ──▶ index/sums
                                           │
                                           ▼
                              complete old-or-new evidence directory
                    ◀ tests attach here: drive the built command with fixture
                      repos/artifacts and observe exit plus exact promoted bytes

Seam 3 — project gate and workflow structure:

    trigger: root conformance and canary harness grade a candidate tree
        │
        ▼
    registry + schemas + YAML ──▶ [ conformance/contracts/canaries ] ──▶ green/red
    planted omission/drift     ──▶ [ built CLI and exact artifacts ] ──▶ cited failure
                                              │
                                              ▼
                                      project gate verdict
                    ◀ tests attach here: mutate one requirement, evidence entry,
                      profile argument, or digest and require a distinct red message

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | The wrapper and all four platform tgz files each contain the license, six governance policies, notices, one SPDX-JSON SBOM, and one component manifest | exact package bytes | R1: every current tgz lacks the required evidence set | Per-artifact inspection defeats an implementation that decorates only the wrapper or leaves one target incomplete |
| 1 | Every embedded component manifest matches package identity and a sorted path/mode/size/SHA-256 inventory of the bytes actually present | exact package bytes | R1: no component manifest exists to validate | Recomputing the inventory from extracted tar entries catches copied claims, self-hashes, unsorted paths, and unlisted payloads |
| edge of 1 | Empty policy/notice/SBOM/manifest files, duplicate or traversal paths, unsafe modes, symlinks, and special-file sources fail before artifact promotion | exact package bytes | R1: the current builder has no evidence-input validation | Hostile staging fixtures distinguish content validation from a filename-only allowlist and preserve the prior promoted set |
| 2 | One registry enumerates every core, public, and bank evidence key with owner, schema, profiles, and requiredness; bank is a strict public superset | preflight evidence plus project gate | R2: publish cannot select or report a profile | Exact registry decoding and set relations reject scattered profile lists and missing owner/schema facts |
| 2 | Verify reports absent FT88 and FT87 records pending, while a conforming fixture composes them without FT83-authored placeholder content | authoritative preflight evidence | R2: no profile requirements participate in preflight | Pending and conforming fixtures together reject both silent omission and an always-red implementation |
| 2 | Public publish is red without FT88 or FT87; bank is additionally red without FT71; permitted conditional records require a reason | authoritative preflight evidence | R2: publish ignores profile requirements | One missing-record mutation per owner catches permissive requiredness and a bank profile that is not a strict superset |
| 3 | A conformant verify fixture emits byte-identical LF-terminated release index and checksum files across two runs with perturbed enumeration order and environment | authoritative preflight evidence | R1: artifacts have no component evidence from which an index can be assembled | Exact byte comparison rejects maps, ambient time/environment, unstable traversal, and independently authored checksums |
| 3 | The release index binds observed artifact size/SHA-256, component/SBOM/inventory digests, policies, target matrix, identity/toolchains, ordered phases, and requirement status | authoritative preflight evidence | R1: no release index binds the current tarballs | Independent recomputation makes a constant manifest or partial metadata summary fail despite valid JSON |
| 3 | `SHA256SUMS` is derived from the index and lists every indexed final artifact exactly once in canonical order | preflight evidence plus exact package bytes | R1: no external checksum set exists for the built tgz files | Cross-checking both surfaces catches a second checksum source, omitted target, duplicate, or self-referential hash |
| 4 | Publish without `--profile`, with an unknown profile, or from a focused run exits 2 or non-authorizing red as specified; full public/bank fixtures can become green only with all requirements | authoritative preflight command | R2: the CLI rejects the profile flag entirely | Negative and positive baselines defeat both an always-green bypass and an always-red guard |
| 4 | Missing, empty, malformed, unknown-version, duplicate-key, mismatched-identity, or digest-mismatched producer records are distinctly red | authoritative preflight evidence | R2: no external-owner record validator exists | One fixture per failure class proves schema parsing and cross-record relationships affect the public verdict |
| edge of 4 | Spaces/globs, no final newline, absent versus empty, control bytes, symlink escape, FIFO/device/socket, descendant cwd, and missing required tools fail safely without blocking | command plus evidence and package seams | R1 and R2: current surfaces do not accept or validate this evidence | Hostile process/filesystem fixtures catch shell expansion, unsafe reads, normalization, and cwd assumptions at public exits |
| edge of 4 | Artifact inspection rejects excessive compressed bytes, member count, or aggregate expanded bytes promptly and preserves the prior generation | authoritative preflight command plus evidence | The current inspector limits only one member and otherwise reads or accumulates the entire archive | A built-command fixture for each independent budget catches compressed-input exhaustion, header floods, decompression bombs, and fail-open replacement |
| 4 | An input changed between validation and promotion produces a drift failure and never promotes an index for mixed bytes | authoritative preflight evidence | R1: no index-generation drift check exists | A synchronized mutation fixture catches validate-then-use implementations that hash different generations |
| 5 | Green and deterministically red runs promote complete generations; interrupted, unreadable, encode-failed, or promotion-failed runs preserve the prior generation | exact promoted evidence directory | Existing atomic preflight lifecycle is already covered; new release-index members are not yet part of it | Extending the established fault matrix to the new files catches direct writes and mixed old/new evidence |
| 5 | Green-after-green, green-after-red, red-after-green, and abandoned-stage reruns are idempotent and leave exactly one complete generation | exact promoted evidence directory | Existing rerun behavior is already covered; the expanded evidence set is absent | Two-run exact-set assertions catch append, stale record retention, and incomplete cleanup |
| 6 | Root conformance cross-checks registry/schema/package/workflow surfaces and rejects publish callers that omit the explicit public profile | project gate | R2: no caller or conformance rule knows a profile | Structural mutation catches a workflow-local list or a publish path that bypasses profile authorization |
| 6 | Canary mutations remove required package evidence, a future-owner requirement, or profile wiring and corrupt an artifact digest, each producing its named red message | project gate canary | R1 and R2: no FT83 evidence canaries exist | Independent mutations prove the oracle bites rather than trusting successful fixture output |

Cheapest wrong implementations checked against the map: wrapper-only evidence fails
the per-artifact row; filename-only packaging fails empty and hostile inputs; a static
index fails observed-digest and repeatability rows; an always-green finalizer fails
every missing/malformed/drift row; an always-red finalizer fails conforming verify and
profile fixtures; duplicated profile lists fail registry cross-checks; and direct
evidence writes fail the prior-generation fault matrix.

### Edge inventory

- **Error path:** invalid governance, SBOM generation failure, archive inspection
  failure, schema mismatch, digest mismatch, drift, and promotion failure are coverage
  rows.
- **Empty/absent input:** every policy, notice, producer record, artifact, manifest,
  profile argument, and prior generation distinguishes absent from present-empty in a
  coverage row.
- **Boundary values:** zero-byte payload rejection, one-entry inventory, all five npm
  artifacts, exactly four target tuples, compressed/member/aggregate archive inspection
  budgets, latest/previous-minor support, and public/bank superset boundaries are coverage
  rows.
- **Malformed input:** JSON syntax, unknown fields/schema versions, duplicate keys and
  normalized paths, invalid SPDX shape, identity disagreement, and malformed registry
  status are coverage rows.
- **Interrupted/partial state:** cancellation, unreadable input, encode failure, input
  drift, abandoned staging, fsync/rename failure, and prior-generation preservation
  are coverage rows.
- **Re-run idempotency:** repeated green, red, and cross-verdict runs with changed
  enumeration order and environment are coverage rows.
- **Hostile environment:** spaces/globs, control bytes, no final newline, descendant
  cwd, missing tools, symlink escapes, special files, unsafe modes, and source drift
  are coverage rows. Thin callers preserve multi-word argv.
- **Every shipped invocation surface:** repository CLI, artifact scripts, hosted
  verify, and hosted publish are in scope. Linked consumer hooks and harness adapters
  are **Won't handle** because they never authorize repository publication; the real
  release CLI and workflows remain usable in scope.
- **Live registry and hosted runner behavior:** **Won't handle** in this deterministic
  slice because no registry transaction occurs here; local workflow structure and
  exact artifacts are covered, while slice 3 owns hermetic and live registry proofs.
- **External FT71/FT87/FT88 record semantics:** **Won't handle** because those rows own
  their content and tests; conforming/missing schema fixtures keep this slice's
  validator callable without duplicating producer policy.

## Out of scope

- **FT83 slice 2, reproducible offline artifacts:** two isolated byte-identical builds,
  per-target offline archives, native/musl runners, and network-denied direct/local-npm/
  internal-registry installation are a separate delivery capability ordered after
  this evidence core. Estimated at 14–20 edits and 4 gate runs.
- **FT83 slice 3, governed npm publication:** first-publication and staged-release state
  machines, candidate/latest transitions, integrity polling, interactive approval,
  publication records, provenance, and deprecation rollback are a separate registry
  transaction capability. Estimated at 16–24 edits and 5 gate runs.
- **FT88 data-handling inventory and FT87 offline/network-control production:** these
  are separate producer capabilities whose records this slice only validates and
  binds. Estimated at 10–16 edits and 3 gate runs per roadmap feature.
- **FT71 versioned local-event evidence production:** this is a separate bank-profile
  audit capability whose record this slice only validates and binds. Estimated at
  12–18 edits and 3 gate runs.
- **Native Windows support:** Windows binaries, package selection, runner contracts,
  path/process semantics, and install/uninstall smokes are a separate platform
  capability; WSL2 continues to consume Linux artifacts. Estimated at 20–30 edits and
  5 gate runs.
