# FT83 slice 2 — reproducible offline artifacts

Status: staged

Decision map: `decisions/governed-offline-release-bundle.md` at `62e50d7`.

## Problem

Bench's governed release-evidence core can build and index the wrapper and four
platform npm packages, but those bytes do not yet form a complete offline release.
There is no per-target archive for direct or local-npm use, no independent rebuild
proof for the complete artifact set, and no clean-room smoke proving that installation,
execution, and removal work without public-network access or binary repair. A locally
repeatable package build therefore cannot yet support the stronger claim that a
release is reproducible and usable offline on every advertised target.

## Solution

Build one deterministic offline archive for each supported target alongside the five
existing npm tarballs. Each archive carries the target binary, the exact wrapper and
platform tarballs, their governed component evidence, and self-contained verification
and installation instructions. Two isolated generation passes compare the complete
release-bound byte set, while each native target runner rebuilds and compares its own
binary and runs the same direct and npm smokes. Linux additionally proves the static
binary on a musl-compatible host.

The authoritative preflight remains the only release-verdict owner. It inspects the
finished archives, extends the canonical release index with their digests and the
reproducibility comparison, and derives final checksums from that index. Clean-room
smokes start with empty state, deny public-network access, disable repair, exercise
direct, local-tarball npm, and hermetic internal-registry journeys, and prove clean
removal. FT83 and publication readiness remain red until slice 3 and the selected
profile's external-owner records are complete.

## User stories

1. As a release consumer, I want one self-contained archive for my supported target,
   so that I can authenticate and run Bench without npm or network access. Line:
   `gpt-5.6-luna` / low. The archive contract is exact, the artifact seam is established,
   and allowlist inspection makes wrong output cheaply gate-observable.

2. As an npm consumer in a disconnected environment, I want the exact wrapper and
   native tarballs plus deterministic local-install instructions, so that I can install
   without a registry cache, downloads, substitution, or binary repair. Line:
   `gpt-5.6-luna` / low. The two-tarball install path already exists and clean-room
   process sentinels can observe every forbidden fallback.

3. As an internal-registry operator, I want to seed and install the same immutable npm
   tarballs platform-first and wrapper-last, so that an air-gapped registry serves the
   same bytes that release verification approved. Line: `gpt-5.6-luna` / low. A
   hermetic npm-compatible fixture makes upload order, exact-version resolution, and
   egress behavior fully observable at the package seam.

4. As a release reviewer, I want two isolated builds and a native rebuild for every
   supported binary to agree byte-for-byte, so that unavailable verification or hidden
   environmental variance blocks the release. Line: `gpt-5.6-terra` / medium. The
   comparison contract is exact, but hosted native runners are the only proof for all
   four architectures and the local gate cannot execute that complete matrix.

5. As a release maintainer, I want the one release index and its derived checksum file
   to bind the offline archives and reproducibility result, so that consumers verify
   inspected bytes from one deterministic evidence graph. Line: `gpt-5.6-luna` / low.
   Slice 1 established the typed finalization seam and exact-byte tests cover this
   schema extension.

6. As the reviewer, I want behavior contracts, workflow checks, and biting canaries to
   reject archive drift, reproducibility bypass, network fallback, and incomplete
   target proof, so that offline and reproducible are enforced claims rather than
   documentation. Line: `gpt-5.6-luna` / medium. The mutations are explicit and
   gate-observable, while the profile's cached routing spends medium effort on changes
   to the oracle.

Stories 1–3 and 5–6 are separable gate-observable slices. Story 4 supplies the hosted
proof that local fixtures cannot; the completed integration uses the highest required
line, `gpt-5.6-terra` / medium, for workflow integration and native-matrix repair while
the remaining implementation stays on its per-story line.

## Implementation decisions

- **Supported targets.** The release promise remains exactly Darwin arm64, Darwin x64,
  Linux arm64, and Linux x64. WSL2 consumes the matching Linux artifact and is not a
  fifth target. Native Windows and every unlisted tuple remain unsupported. The
  canonical platform matrix is the only target registry used by builders, evidence,
  workflows, instructions, and tests.
- **Binary contract.** All binaries use the pinned Go toolchain, canonical version
  stamp, `-trimpath`, and stripped symbol/debug tables. Linux builds set CGO off and
  must be statically linked. Native runners execute the target binary; Linux runners
  additionally execute it in a musl-compatible environment. A cross-compile without
  the matching native rebuild and smoke cannot establish target support.
- **Complete artifact set.** One generation emits five npm tarballs and four
  `redbench-<version>-<os>-<arch>.tar.gz` archives. Each archive has one canonical
  top-level directory and an exact allowlist: the executable native binary, the exact
  wrapper and matching platform npm tarballs, the governed component evidence already
  owned by slice 1, and LF-terminated offline instructions covering verification,
  direct execution, local npm, internal-registry seeding, and removal. No source tree,
  maintainer command, cache, or alternate target enters an archive.
- **Deterministic construction.** One private-staging builder owns npm packages,
  archives, modes, ordering, metadata normalization, and atomic promotion of the
  complete set. Inputs must be explicit regular files. Archive paths are byte-sorted;
  uid/gid, owner/group, modes, and timestamps are fixed; gzip metadata is normalized.
  The builder never edits an existing artifact in place and never publishes a partial
  generation.
- **Reproducibility proof.** A canonical isolated build produces the candidate set. A
  second clean generation from the same tagged commit, pinned toolchains, declared
  flags, and empty build state independently recreates all binaries, npm tarballs,
  offline archives, component manifests, SPDX SBOMs, inventories, and deterministic
  evidence. Every release-bound byte must match; there is no normalized comparison or
  ignored variance. Each target runner independently rebuilds its binary and compares
  it with the binary extracted from both its platform package and offline archive.
  Missing runners, tools, outputs, or comparisons are red.
- **Direct offline journey.** Starting from an empty home and work directory, the
  smoke verifies the externally supplied release-index and checksum relationship,
  verifies the archive digest, extracts it, verifies every embedded component
  inventory, and executes the native binary directly. It runs `version` plus one
  read-only operational command, then removes the extracted installation and proves
  that no home, cache, prefix, or work residue remains.
- **Local npm journey.** Starting from an empty home, prefix, and npm cache, the smoke
  installs the included platform tarball and wrapper tarball into an isolated prefix
  with npm offline mode enabled, exact local paths, no pre-existing cache, and Bench
  repair disabled. It runs the installed wrapper's `version` plus the same read-only
  operation, uninstalls, and proves no prefix package, Bench cache, or unexpected
  residue remains. Missing native bytes, optional-dependency fallback, fetch, repair,
  rebuild, or byte substitution is red.
- **Internal-registry journey.** A hermetic npm-compatible registry fixture is the only
  permitted socket endpoint. The smoke uploads the exact platform tarball first and
  wrapper tarball last, verifies stored integrity, then installs both by exact version
  from an empty client home/cache with public egress denied and repair disabled. It
  executes and removes the installation under the same contract as local npm. The
  fixture is a test adapter, not a second package builder or a preview of slice 3's
  public publication state machine.
- **Network denial.** Direct and local-npm smokes permit no network endpoint. The
  internal-registry smoke permits only its explicitly provisioned loopback fixture;
  all public or undeclared egress is denied. Command sentinels and observed registry
  requests make fetch attempts red. `BENCH_NO_REPAIR=1` remains the slice-local repair
  control until FT87 owns the general offline/network policy.
- **Evidence extension.** The typed release-evidence module inspects final archive and
  package bytes. The canonical deterministic release index adds the four archive
  identities, sizes, SHA-256 digests, component relationships, and the ordered
  two-build/native-target comparison results. `SHA256SUMS` remains derived from that
  index and covers every independently delivered package and archive exactly once.
  No artifact embeds its own final digest, and no second release manifest is created.
- **Failure and retry.** Missing, empty, malformed, unsafe, mismatched, drifting, or
  unproved artifact input fails before promotion and preserves the prior complete set.
  Interruption cannot leave a promoted partial build, permissive network state, or
  retained fixture process. A rerun from the same inputs is byte-identical; a rerun
  after failure replaces only through the existing complete old-or-new promotion.
- **Release readiness.** This slice prepares and verifies offline bytes but performs no
  public npm transaction. Publish mode remains red for absent selected-profile records,
  and FT83 remains on the roadmap until the governed-publication slice is implemented.

## Testing decisions

- Tests attach to the exact built package/archive set, the built authoritative
  preflight command and promoted evidence, and the native workflow. They inspect and
  execute real bytes through public scripts and commands; they do not mock private Go
  collaborators, duplicate builder policy, or infer reproducibility from source flags.
- Artifact behavior extends the existing private-stage, exact-inventory, target
  selection, atomic-promotion, and host-smoke contracts. Preflight behavior extends
  the existing deterministic-index and observed-digest tests. The workflow matrix runs
  the same checked-in smoke entry points used locally.
- The feature must pass the project gate: `.bench/gate.sh`. Hosted native-runner jobs
  are additional release evidence; an unavailable required runner fails closed rather
  than being treated as a skipped target.
- Pre-implementation probes run on 2026-07-17 are named `R1` through `R5`:
  - `R1` built the five current npm tarballs twice in separate temporary directories;
    all five compared byte-identical. This local npm subset is already covered, but it
    does not prove isolated full-set or native-runner reproducibility.
  - `R2` listed the current artifact set and found no per-target `.tar.gz` archive.
  - `R3` searched the current release index and found no offline-archive or
    reproducibility comparison evidence.
  - `R4` searched the current smoke and found no offline mode, network denial,
    internal-registry journey, repair disablement, or uninstall contract.
  - `R5` inspected the current native and release workflows and found no independent
    rebuild comparison or musl-compatible execution.

### Seam diagrams

Seam 1 — exact release artifact set:

    trigger: authoritative artifact phase builds the tagged source
        │
        ▼
    source + matrix + evidence ──▶ [ private-staging artifact builder ] ──▶ 5 npm tgz
    pinned toolchains + flags   ──▶ [ deterministic archive assembly ] ──▶ 4 tar.gz
                                               │
                                               ▼
                                  atomically promoted complete byte set
                    ◀ tests attach here: inspect exact entries/modes/digests,
                      rebuild twice, and inject unsafe, drifting, partial inputs

Seam 2 — offline consumer journeys:

    trigger: target smoke receives approved index, sums, and artifact set
        │
        ▼
    empty home/cache/prefix ──▶ [ direct | local npm | internal registry ] ──▶ bench
    denied public network   ──▶ [ verify, install, run, uninstall          ] ──▶ residue
                                           │
                                           ▼
                                  observed requests and filesystem state
                    ◀ tests attach here: run real bytes with sentinels/fixture,
                      observe version/read-only output, requests, and clean removal

Seam 3 — authoritative evidence and native proof:

    trigger: full preflight plus every target's native workflow completes
        │
        ▼
    final artifacts + rebuilds ──▶ [ typed evidence finalizer ] ──▶ verdict/index/sums
    ordered comparison results ──▶ [ native runner + gate     ] ──▶ supported matrix
                                              │
                                              ▼
                                  one deterministic evidence generation
                    ◀ tests attach here: recompute all digests, compare exact bytes,
                      omit/corrupt one proof, and require an attributed red verdict

### Demonstrated red-mutation log

Every entry below resolves to an executable mutation; the quoted diagnostic is the
observed red from that seam. The same test passes after restoring the mutation.

| mutation | executable seam command | observed red diagnostic |
|---|---|---|
| M1 archive parser | `go test -count=1 -run '^TestBuiltCommandRejectsConcatenatedGzipArchiveAndPreservesPriorGeneration$' ./internal/preflight` | `concatenated gzip members` |
| M2 hostile archive inventory | `go test -count=1 -run '^TestBuiltCommandRejectsHostilePackageArchivesAndPreservesPriorGeneration$' ./internal/preflight` | `archive contains duplicate path`; `unsafe archive path`; `archive contains special file`; `archive contains unsafe mode`; `archive contains empty member` |
| M3 promotion | `go test -count=1 -run '^TestArtifactPromotionIsAtomicAndExclusive$' ./internal/contract/surface/artifact` | `promotion interruption changed prior-generation bytes` |
| M4 artifacts and release evidence | `go test -count=1 -run '^TestReproducibilityComparatorRejectsArtifactAndEvidenceMutations$' ./internal/contract/surface/artifact` | `reproducibility mismatch:`; `missing second-build artifact`; `unexpected artifact`; `release-evidence mismatch`; `missing second-build release evidence` |
| M5 native proofs | `go test -count=1 -run '^TestNativeProofAggregationRejectsIncompleteAndRedProofs$' ./internal/preflight` | `native proof set is missing`; `native proof is incomplete or red` |
| M6 workflow | `go test -count=1 -run '^TestNativeWorkflowEvidenceEdgeBites$' ./internal/conformance` | `native verification does not run smoke from finalized evidence` |
| M7 approved offline evidence | `go test -count=1 -run '^TestOfflineSmokeRequiresApprovedReleaseEvidence$' ./internal/contract/surface/artifact` | `offline smoke: approved release evidence is missing or unsafe` |
| M8 malformed package | `go test -count=1 -run '^TestBuiltCommandControlByteArchivePathPreservesPriorGeneration$' ./internal/preflight` | `unsafe archive path` |

Demonstrated on 2026-07-17: M1 rejected a controlled trailing gzip member; M2 rejected
duplicate, traversal, special, unsafe-mode, and empty members; M4 captured equal,
different, missing, extra, and release-evidence mutations; M5 captured missing,
corrupt, wrong-runner, tool, and no-op proof mutations; M6 removed the evidence-to-smoke
workflow edge; and M7 withheld the approved index/checksums.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Exactly four target archives are emitted from the canonical four-row matrix, each named for version/os/arch and paired with its wrapper and platform package | exact release artifact set | M2: delete one planned archive; `TestOfflineArchiveProjection` reports the missing entry | Exact cardinality and matrix projection reject a wrapper-only archive, an omitted target, or a second target registry |
| 1 | Each archive contains only its executable target binary, exact two npm tarballs, governed evidence, and complete LF-terminated offline instructions under one canonical root | exact release artifact set | M2: remove or add one archive member; the inventory assertion is red | Independent archive inspection catches source leakage, wrong-target bytes, missing journeys, unsafe modes, and filename-only packaging |
| 1 | A verified archive supports direct `version` and one read-only operation with no npm, Node, cache, repair, or network, then removes without residue | offline consumer journeys | M7: remove the direct `commands --brief` probe or permit repair; smoke is red | Removing npm/Node and denying endpoints defeats an implementation whose direct path secretly depends on package tooling or repair |
| 2 | Empty-home local npm installs the included platform and wrapper tarballs by path with offline mode and repair disabled, runs both operations, uninstalls, and leaves no Bench cache or package residue | offline consumer journeys | M7: remove offline/repair/uninstall sentinel; smoke is red | Empty state plus request and residue observation catches cache dependence, download fallback, hidden repair, and incomplete removal |
| edge of 2 | Missing platform bytes, wrong target, corrupt tarball, checksum mismatch, an attempted fetch/rebuild/repair, or unexpected residue is distinctly red | offline consumer journeys | M7: one command mutation per failure class; each has attributed `offline smoke:` red output | One mutation per class prevents a happy-path-only installer from claiming the offline contract |
| 3 | The hermetic registry receives the exact platform tarball before the wrapper, retains matching integrity, and serves an exact-version install with only loopback allowed | offline consumer journeys | M7: reverse observed PUT sequence or digest; smoke reports upload order/digest red | Recorded requests and stored-byte digests reject wrapper-first upload, mutable substitution, ambient public registry use, and a fake no-op fixture |
| 3 | Empty-client internal-registry install runs both operations with repair disabled, then uninstalls without client or server residue | offline consumer journeys | M7: remove exact version or uninstall; smoke reports missing wrapper/residue red | Exercising the seeded package graph defeats an implementation that tests upload without resolvable wrapper dependencies |
| 4 | Two isolated clean generations reproduce every binary, npm tarball, offline archive, generated evidence record, release index, and checksum byte-for-byte | exact release artifact set plus authoritative evidence | M4: flip one compared byte or fabricate `match:true`; comparator/evidence validation is red | Enumerating and comparing the full set rejects normalized comparison, ignored variance, and a degenerate check of only one artifact class |
| 4 | Each of Darwin arm64, Darwin x64, Linux arm64, and Linux x64 rebuilds its binary on the declared native runner and matches the package and archive bytes exactly | native proof | M5: omit one proof or change its runner/digest; aggregate is red | One proof per enumerated tuple prevents one host or one cross-compiler from standing in for the supported matrix |
| 4 | Both Linux binaries are stripped static executables that run on their native runner and in a musl-compatible environment; Darwin binaries are stripped and run natively | native proof | M5: remove musl, operation, or strip status; aggregate is red | Runtime execution plus binary inspection catches flags that look correct in source but yield a dynamically linked, unstripped, or unusable payload |
| edge of 4 | A missing tool, runner, rebuild output, comparison record, or unsupported target is red and cannot be recorded as skipped support | native proof plus authoritative evidence | M5: use missing/wrong runner or tool; native proof exits non-zero | Negative matrix fixtures defeat best-effort verification and an always-green comparison summary |
| 5 | The canonical release index binds all nine final artifacts by identity/size/SHA-256 and binds ordered generation/native comparisons; `SHA256SUMS` derives exactly those nine rows | authoritative evidence and native proof | M4: remove or alter one index/checksum row; evidence test is red | Recomputing final bytes and sums catches a second manifest, stale digest, self-reference, omitted archive, or caller-authored success |
| 5 | Reordered enumeration, locale, timezone, home, cache, temp root, and build directory do not change deterministic evidence or artifacts | exact release artifact set plus authoritative evidence | M4: perturb one isolated root; byte comparator is red on drift | Perturbed isolated builds catch ambient metadata and ordering leaks that same-directory reruns miss |
| edge of 5 | Artifact or evidence drift between validation and promotion fails and preserves the prior complete artifact and evidence generations | exact release artifact set plus authoritative evidence | M3: mutate at promotion seam; prior-byte assertion is red on replacement | Synchronizing a mutation at promotion catches validate-then-use and mixed-generation implementations |
| 6 | Behavior contracts drive real direct/local/registry journeys under network sentinels and assert exact output, requests, uninstall, and residue | project gate | M7: mutate egress, repair, request, or residue control; smoke is red | A real-process contract makes documentation-only or mocked-network implementations fail |
| 6 | Workflow conformance derives all four jobs from the canonical matrix and requires clean generation, native rebuild comparison, native smoke, Linux musl smoke, and complete evidence upload | project gate plus native proof | M6: remove one job edge/path; root conformance is red | Structural and runnable checks together catch a bypassed target, duplicated matrix, or workflow-only claim that never produces evidence |
| 6 | Canaries corrupt one archive digest, force a reproducibility mismatch, omit one target proof, and permit one network/repair attempt, each with a distinct red message | project gate canary | M1/M4/M5/M7: mutate each surface; gate reports its distinct red | Independent mutations prove every new authorization dependency bites instead of trusting a green fixture |
| edge of 1–6 | Spaces/globs in roots, control bytes in archive facts, absent versus empty files, missing final newline, symlink/FIFO/device/socket inputs, multi-word argv, deep cwd, and symlink invocation fail safely or preserve exact behavior | all three seams | M1/M2/M8: hostile fixture mutation is red | Hostile fixtures exercise the profile's complete shell-input checklist at public seams rather than private parsers |
| edge of 1–6 | SIGINT during build, compare, install, registry service, or removal leaves no promoted partial generation, live fixture, permissive network state, or unrecoverable residue; rerun completes idempotently | all three seams | M3: interruption mutation preserves previous generation bytes | Interrupt-at-stage tests and a second run catch leaked processes/state and non-atomic recovery |
| edge of 1–6 | Dirty, foreign, nested, or identity-mismatched release checkout state remains red before artifact construction | authoritative preflight | Existing preflight ownership is already covered; the new builder must remain behind it | Keeping the build callable only through the authoritative release path prevents reproducible bytes from authorizing an untrusted source state |

Cheapest wrong implementations checked against the map: four archives copied from one
target fail exact target identity and native execution; deterministic local packages
without independent proof fail the isolated/native comparison rows; a direct smoke
that leaves npm available fails the minimal-tool journey; `npm --offline` backed by a
warm cache fails empty-state and request assertions; a registry that rewrites tarballs
fails stored integrity; an always-green comparison fails mismatch/missing-runner
fixtures; an always-red comparison fails the conformant two-build baseline; a second
checksum manifest fails index derivation; and workflow prose without biting checks
fails the four canary mutations.

### Edge inventory

- **Error path:** bad digest, wrong target, failed extraction, unavailable tool/runner,
  fetch/repair/rebuild attempt, install or uninstall failure, registry error,
  reproducibility mismatch, evidence drift, and promotion failure are coverage rows.
- **Empty/absent input:** package, archive, binary, instructions, checksum, comparison,
  runner output, home, cache, prefix, and prior generation distinguish absent from
  present-empty in coverage rows.
- **Boundary values:** the exact four-target, five-package, four-archive, and nine-row
  checksum sets; one permitted loopback endpoint; one canonical archive root; and zero
  ignored byte differences are coverage rows.
- **Malformed input:** invalid tar/gzip, traversal or duplicate normalized entries,
  wrong modes/types, corrupt npm metadata, malformed registry responses, unknown
  target/comparison status, and inconsistent index relationships are coverage rows.
- **Interrupted/partial state:** SIGINT at build, comparison, install, registry, removal,
  evidence encoding, and promotion plus prior-generation preservation are coverage
  rows.
- **Re-run idempotency:** repeated clean generation, failure recovery, existing output,
  empty/warm fixture restart, and post-uninstall reinstall are covered by exact-byte
  and residue rows.
- **Hostile environment:** spaces/globs, control bytes, missing final newline,
  absent-versus-empty, special files, multi-word argv, missing tools, symlink/deep-cwd
  invocation, ambient locale/timezone/home/cache/temp changes, dirty release checkout,
  and denied egress are coverage rows.
- **Every shipped invocation surface:** the repository CLI/preflight, artifact scripts,
  offline instructions, installed wrapper, direct binary, and native/release workflows
  are covered. Hooks and harness adapters are **Won't handle** because they neither
  construct nor consume release artifacts; direct and installed CLI callers remain
  usable under this exclusion.
- **Public npm registry behavior:** **Won't handle** because slice 3 owns registry
  staging, approval, promotion, retries, and rollback; the in-scope hermetic registry
  proves only that the exact offline tarballs can seed an internal registry.
- **External signing and acquisition trust:** **Won't handle** because signing-key and
  distribution-channel trust are separate capabilities explicitly excluded by the
  decision map; this slice verifies supplied index/checksum relationships without
  inventing a trust root.
- **FT87's general offline policy:** **Won't handle** because FT87 owns the reusable
  network/resource control record; this slice's explicit npm offline mode, egress
  denial, and repair disablement keep all three consumer journeys executable now.

## Out of scope

- **FT83 slice 3, governed npm publication:** first-publication and staged-release
  state machines, registry integrity polling, candidate/latest transitions,
  interactive approval, publication records, provenance, and deprecation rollback are
  a separate public-registry transaction capability. Estimated at 16–24 edits and 5
  gate runs.
- **FT87 general offline/network controls:** reusable repository-wide network and
  resource policy plus its producer evidence are a separate capability; this slice
  carries only the explicit controls needed to prove its artifact journeys. Estimated
  at 10–16 edits and 3 gate runs.
- **FT88 data-handling and FT71 local-event record production:** these are separate
  producer capabilities already represented by fail-closed release-evidence slots.
  Estimated at 10–18 edits and 3 gate runs per roadmap feature.
- **Native Windows support:** Windows binaries, package selection, runner contracts,
  path/process semantics, and install/removal smokes are a separate platform
  capability; WSL2 continues to consume Linux artifacts. Estimated at 20–30 edits and
  5 gate runs.
- **External signature trust and distribution-channel authentication:** signing keys,
  custody, transparency, mirrors, and acquisition trust are a separate supply-chain
  capability outside the decision map. Estimated at 18–28 edits and 5 gate runs.
