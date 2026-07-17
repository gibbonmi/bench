## Destination

Define the governed, offline-verifiable release bundle FT83 adds around Bench's
existing `redbench` packages and authoritative release preflight, without reopening
their settled identity or ownership.

## #1: What exact target matrix does Bench support?

Blocked by: none
Type: Grill

### Question

The builder currently emits Darwin arm64/x64 and static Linux arm64/x64 packages,
while documentation promises only macOS or Linux and native CI covers those four
runners. Which targets, libc environments, archive formats, and support tiers are a
release promise, and what evidence is required before a target enters that promise?

### Answer

Bench has one supported target tier with exactly four native tuples:
`darwin/arm64`, `darwin/amd64`, `linux/arm64`, and `linux/amd64`. Linux binaries
are static and must work in glibc and musl environments. WSL2 consumes the matching
Linux artifact; it is not a separate target. Native Windows and every other tuple are
unsupported.

A target enters or remains in the matrix only when that release passes its native
binary-format check plus clean-install, operational, and network-disabled smokes on
the target runner. An unproved target is unsupported rather than best-effort.

## #2: What must every independently published artifact contain?

Blocked by: #1, #4
Type: Grill

### Question

Define the common and artifact-specific allowlists for the wrapper package, each
platform package, and each offline archive: license and security material,
dependency/license notice, SBOM, checksums, package inventory, support metadata, and
the executable payload. Set which records are embedded versus referenced, and which
missing or inconsistent record makes packaging fail closed.

### Answer

Use two evidence levels. Every wrapper package, platform package, and offline archive
carries the license, governance policies, third-party notices, an SPDX-JSON SBOM, and
a component manifest that inventories and SHA-256-checksums its internal files.
Platform packages add only their native binary; the wrapper adds its canonical
shipped assets; each offline archive adds the matching wrapper and platform tarballs
plus offline instructions.

One external release index checksums the finished artifacts and binds their component
manifests, avoiding a self-reference from embedding an artifact's final digest inside
itself. Packaging fails closed on a required record or payload that is missing, empty,
malformed, unlisted, wrong-type, or inconsistent with the source, version, or digest.

## #3: What is the offline installation and verification contract?

Blocked by: #1, #2
Type: Grill

### Question

Specify the network-independent deliverable and user journey for direct local
installation and an internal npm registry: artifact acquisition assumptions,
checksum and manifest verification, wrapper/native dependency resolution, commands,
cache or repair behavior, and the smoke that proves installation and execution with
network access denied.

### Answer

Each target gets one `redbench-<version>-<os>-<arch>.tar.gz` containing the native
binary, matching wrapper and platform npm tarballs, the component evidence from #2,
and exact local/internal-registry instructions. The separately delivered release
index and `SHA256SUMS` authenticate the finished archive; after extraction, the
component manifests authenticate every internal file.

The local npm path installs the two included tarballs into an isolated prefix with
npm offline mode enabled, no registry cache assumed, and binary repair disabled. The
direct path may execute the verified native binary without npm. The internal-registry
path uploads the exact same tarballs platform-first and wrapper-last, then installs by
exact version. Neither path rebuilds, downloads, or substitutes bytes.

Release smokes start with an empty home, prefix, and npm cache; deny network; verify
the index, archive, and component checksums; install; run `version` and one read-only
operational command; and uninstall without residue. Missing platform bytes, a fetch
attempt, repair attempt, checksum mismatch, wrong selected binary, or unexpected
residue is red. Until FT87 supplies the global offline control, this path explicitly
sets npm offline mode and disables Bench repair.

## #4: Which repository governance promises ship with a release?

Blocked by: none
Type: Grill

### Question

Set the current-state policy contract for supported versions and EOL, vulnerability
intake and response targets, dependency updates and license changes, threat model,
recovery and rollback, and a non-personal support route. Decide which values are hard
release requirements, who owns them, and which are allowed to vary by version.

### Answer

Governance records are hard release inputs. Bench supports the latest minor release
and the previous minor for 90 days, with explicit EOL state. GitHub Security
Advisories is the private vulnerability-reporting route and GitHub Issues is the
non-personal support route; the personal-email dependency leaves the policy.

Security reports are acknowledged within three business days and triaged within
seven. Critical issues require mitigation within seven days, high within 30, and
medium within 90; low severity is best-effort. Current supported-version/EOL,
security-response, dependency/license-change, threat-model, recovery/rollback, and
support policies must all exist for release. The repository owns their current
values, each release manifest binds their digests, and the release records its
version-specific rollback target.

## #5: What publication transaction can npm actually support?

Blocked by: #2
Type: Research

### Question

Using current primary npm documentation and non-publishing probes, establish the
observable semantics for preflighting names and versions, comparing an existing
immutable version by digest, publishing under a non-default tag, waiting for platform
dependencies, publishing the wrapper last, promoting tags, resuming after partial
success, and rolling back without pretending an immutable version can be removed.
Record the cited findings and probe commands in a short research asset.

### Answer

npm now supports reviewable staged publishing, but only for existing package
identities and only with npm 11.15+ and Node 22.14+. CI can use stage-only OIDC;
approval remains an interactive maintainer action with 2FA. The first publication of
Bench's still-unpublished identities must therefore use a separately governed direct
publish under a non-default tag. All publication is package-scoped and immutable, so
retries may accept only byte-identical existing versions, the wrapper stays last, and
rollback uses tags plus deprecation rather than unpublish. The cited semantics,
two-path consequence, and runnable read-only/integrity probes are recorded in
[governed-offline-release-publication-research.md](governed-offline-release-publication-research.md).

## #6: What release-state machine does Bench promise?

Blocked by: #5
Type: Grill

### Question

Given npm's proven behavior, define the repository-controlled states and transitions
for preflight, staged publication, dependency verification, wrapper publication,
promotion, retry, already-present artifacts, failure, and rollback/deprecation. Set
which transitions are automatic, which require reviewer authorization, and what
machine-readable evidence each transition leaves.

### Answer

The state machine is build and locally verify the complete immutable set; submit every
package under a version-specific candidate tag; verify or approve platform packages
first; publish or approve the wrapper last; reverify the complete live set; then
promote platform `latest` tags first and the wrapper last. First publication uses
reviewer-authorized direct publish because the identities do not yet exist. Later
releases use stage-only OIDC and interactive 2FA approval.

Retries accept an already-live package only when registry integrity exactly matches
the approved local tarball. Before approval, failed stages are rejected. After any
package is live, failure preserves the old `latest`, removes candidate tags,
deprecates the bad version, and requires a new version; automation never unpublishes.
Reviewer presence is required for first publication, staged approval or rejection,
and final promotion. Verification and bounded registry polling are automatic.

Every transition records package, version, local and registry integrity, stage ID
when applicable, authentication mode, tag state, timestamp, and result.

## #7: What counts as a reproducible release?

Blocked by: #1, #2
Type: Grill

### Question

Define the comparison subject and environment: byte-identical binaries, npm
tarballs, offline archives, generated records, or a normalized subset; same host
versus independent supported runners; number and timing of builds; permitted
variance; and the red signal when reproducibility cannot be established.

### Answer

Release-bound output must be byte-identical across two isolated builds of the tagged
commit with pinned toolchains and inputs. The canonical build is checked by a native
target-runner rebuild of each binary. A second clean generation pass independently
recreates npm tarballs, offline archives, SPDX SBOMs, inventories, checksums, and
deterministic manifests. Every compared byte must agree; there is no normalized or
ignored variance, and an unavailable verifier or mismatch blocks release.

npm provenance and publication receipts are created after the reproducible build and
are excluded from its comparison subject. Their digests are bound afterward as
external publication evidence.

## #8: What does the deterministic release manifest bind?

Blocked by: #2, #4, #6, #7
Type: Grill

### Question

Set the manifest's versioned schema, canonical serialization, and ownership for
source/version identity, Go/Node/npm versions and flags, dependency and platform
inputs, gate/race/vet/vulnerability results, artifact inventories and SHA-256
digests, reproducibility comparison, governance-record digests, publication state,
and rollback target. Decide whether it is the sole machine-readable index from which
checksums and evidence relationships are derived.

### Answer

The existing authoritative preflight remains the sole orchestrator and atomic
promoter. Its caller supplies only `RunEvidence` (mode, scope, identity, and ordered
phase results) to one `FinalizeEvidence` façade; callers cannot configure schema,
paths, ordering, hashes, status, or required records. A deep typed release-evidence
module behind that façade owns safe discovery, the canonical requirement registry,
validation, inventories, SHA-256 calculation, deterministic encoding, drift checks,
and old-or-new evidence generation. It uses concrete owner adapters only; there is no
runtime-pluggable raw-JSON evidence graph.

Versioned UTF-8 JSON with LF termination is canonical: fixed struct field order,
canonically sorted arrays, no unordered maps, ambient timestamps, locale, or
environment. Embedded component manifests bind component identity and sorted internal
path/mode/size/SHA-256 inventories. The external `release-index.json` is the sole
deterministic index and binds source commit, tag, package version, rollback target,
toolchains and flags, dependency/platform/policy inputs, ordered preflight results and
record digests, evidence-requirement status, target matrix, artifact size/SHA-256 and
component/SBOM/inventory digests, and the reproducibility comparison. `SHA256SUMS` is
derived from that index, never independently authored.

Publication creates a separate non-deterministic `publication-record.json` that
references the immutable release-index digest and carries ordered registry
transitions, stage IDs, authentication mode, tag and integrity observations,
timestamps, result, and provenance digests. It never rewrites the reproducible index
or package bytes. Trustworthy red runs produce complete red diagnostic evidence;
unreadable or unsafe inputs, cancellation before a trustworthy record, encoding
failure, drift, or promotion failure preserve the prior complete generation.

## #9: How does FT83 depend on release evidence owned by other roadmap rows?

Blocked by: #3, #4, #8
Type: Grill

### Question

FT83 must eventually carry FT88's data-handling inventory and may bind FT71 event
evidence and FT87 offline/network controls. Decide whether FT83 ships a stable slot
and fails closed until required producer records exist, omits not-yet-applicable
records with an explicit reason, or waits for those rows before it can close. Keep
record production single-sourced in the owning row.

### Answer

FT83 owns one canonical requirement registry with stable evidence keys, owning
roadmap row, applicable release profiles, requiredness, and schema version. Producer
rows own their records; FT83 only validates, digests, packages, and indexes them.
Required records can never be silently omitted or marked not applicable. A genuinely
conditional record may use `not_applicable` only when the registry permits it and a
reason is recorded.

The registry exposes `public` and `bank` profiles; `bank` is a strict superset.
`public` requires FT88's data-handling inventory and FT87's offline/network control
record. `bank` additionally requires FT71's versioned local-event evidence. Publish
mode requires an explicit profile and fails red on any missing required record.
Verify mode may record future-owner requirements as pending because it cannot
authorize publication.

FT83 may ship its machinery before FT88, FT71, or FT87, but real publication stays
red for any selected profile whose required producer record is absent. The FT83 gate
proves that fail-closed posture with missing-record fixtures and proves composition
with conforming fixtures; it does not duplicate placeholder producer content. This
lets FT83 close without falsely declaring the repository release- or bank-ready.

## #10: Should the closed map yield one spec or dependent slices?

Blocked by: #1, #2, #3, #4, #5, #6, #7, #8, #9
Type: Grill

### Question

The current tree has clean seams around artifact construction, authoritative
preflight/evidence, workflow publication, and governance records. Once their
contracts are settled, choose whether FT83 lands atomically or as reviewer-approved
slices, and set the dependency order without allowing an intermediate slice to make
the release appear governed or offline-ready prematurely.

### Answer

Use three dependent slices:

1. **Governance and deterministic evidence core.** Add the canonical requirement
   registry, policy records, SPDX SBOM and notices, component/release schemas, and
   authoritative preflight validation, including fail-closed external-owner slots.
2. **Reproducible offline artifacts.** Generate and compare the complete artifact set,
   add per-target archives and final checksums, and prove direct, local-npm, and
   internal-registry installation with network denied.
3. **Governed npm publication.** Implement the first-publication and staged-release
   state machines, integrity polling, candidate/latest tag transitions, interactive
   approval handoffs, publication records, and deprecation rollback.

Each slice keeps one source per fact and leaves the release-readiness verdict red
until its downstream dependencies and required producer records are complete. The
FT83 roadmap row retires only after all three slices ship.

## Not yet specified

- None.

## Out of scope

- Publishing a real version, claiming package names, changing registry settings, or
  handling publisher credentials during this work.
- External signature trust, signing-key custody, host IAM, central CI administration,
  registry administration, endpoint controls, firewalls, or SIEM retention.
- Reopening the `redbench` npm identity, tag-driven publication authority, or the
  single Go-core release-preflight decision.
- Producing FT88's data-handling inventory, FT71's local event log, or FT87's general
  network/resource policy; FT83 only defines how applicable records are carried and
  bound.
- FT76 setup behavior, transactional managed-asset lifecycle, and runtime binary
  repair beyond what the offline release contract directly requires.

## Handoff

1. **Module boundaries.** `internal/preflight` remains the authoritative ordered
   phase orchestrator and atomic evidence promoter. A deep typed release-evidence
   module owns the requirement registry, component/release schemas, deterministic
   assembly and verification, safe traversal, digests, reproducibility comparison,
   and bundle inspection. A deep publication module owns the resumable npm state
   machine and an npm-registry port; public npm and the hermetic registry fixture are
   its two adapters. Artifact scripts and hosted workflows are thin callers.
   Governance and future FT71/FT87/FT88 producers own record content only.
2. **Contracts.** `FinalizeEvidence(ctx, root, RunEvidence)` derives and atomically
   promotes complete green or trustworthy red evidence; typed evidence errors fail
   release and callers cannot reinterpret them. `bench release-preflight --mode
   verify [--phase <name>]` remains diagnostic; publish mode additionally requires
   `--profile public|bank`, and only its full green evidence authorizes release.
   Focused evidence never authorizes publication. New `bench release
   prepare|submit|status|promote|rollback` operations consume the release directory
   and durable publication record. They are idempotent, non-interactive, emit compact
   TOON on stdout (including structured errors and exact `next_action` for external
   2FA approval), keep progress on stderr, and use exit 0 success/no-op, 1
   unsatisfied release intent, 2 usage. They never ingest credentials into evidence.
3. **Deep vs thin.** Release evidence and publication are deep: callers do not know
   schemas, requiredness, ordering, hashing, archive layout, retry classification, or
   registry transition rules. Preflight is the deep authority composing their
   verdicts. Shell wrappers, npm commands, workflows, and AXI rendering are thin and
   contain no duplicate policy.
4. **Black-box assertables.** At preflight: exit, attributed stderr, complete
   old-or-new evidence directory, deterministic bytes, and focused-vs-authoritative
   status. At the artifact seam: exact archive/package inventories and modes, digest
   agreement, two-build byte equality, target selection, network-denied install and
   clean uninstall. At publication: TOON state/next action, registry call order,
   platform-first/wrapper-last behavior, integrity-aware retry, tag state,
   deprecation rollback, and durable record contents without secrets.
5. **Gate attachment.** Parse/validity checks validate every JSON schema and
   cross-check the one requirement registry against generated/docs surfaces.
   Behavior contracts invoke the built CLI against throwaway repositories and a
   hermetic npm-registry fixture; they exercise real package/archive bytes and deny
   network for offline smokes. Native workflows rebuild and compare each supported
   binary. Canaries plant at least missing required evidence, special-file input,
   artifact digest drift, reproducibility mismatch, publication-order bypass, and
   premature wrapper promotion, each with a distinct message. Every enforcement
   dependency fails closed; no missing tool or runner can authorize release.
6. **Hostile-input owners.** Release evidence owns spaces/globs, missing final
   newline, absent-vs-empty, special files, traversal, and input drift. Preflight
   identity owns control bytes in Git-derived tag/commit/path facts and missing
   required tools. Thin shell callers own quoting of multi-word arguments. CLI/root
   resolution owns symlink and deep-CWD invocation across source, workflow, and
   installed surfaces. Publication owns control bytes from registry data, duplicate
   or resumed transitions, partial live sets, and plan/apply drift. Preflight and
   publication jointly own SIGINT cleanup, durable recovery, and rerun idempotency.
   Existing worktree ownership continues to reject dirty, foreign, nested, or
   identity-mismatched release checkouts before bytes are built.
7. **Uncertainty flags.** No reviewer decision remains. External proofs still needed
   at implementation time: npm staged publishing requires pinned npm 11.15+ and an
   existing package; public name ownership/permissions cannot be inferred from E404;
   native-runner and live-registry evidence cannot be replaced by the hermetic gate.
8. **Rejected alternatives.** Native Windows; best-effort target tiers; embedding an
   artifact's final digest inside itself; normalized reproducibility; caller-authored
   manifests; a runtime-pluggable raw-JSON evidence graph; producer-owned release
   verdicts; direct `latest` publication; wrapper-first publication; rollback by
   unpublish; silent omission of future-owner records; and one atomic FT83 spec.
9. **Domain watch-outs.** npm permanently burns a name/version even after unpublish;
   a candidate tag is public, not private; staging cannot create the first package;
   OIDC stages but interactive 2FA approves; registry receipts are intentionally
   non-reproducible; hashes must come from bytes actually inspected; required record
   absence is red; and an unsupported or unproved target must never appear supported.

Dependency order: governance and deterministic evidence core → reproducible offline
artifacts → governed npm publication.
