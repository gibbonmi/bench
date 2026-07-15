## Destination

Deliver FT81 as one exact-tarball contract: clean, platform-correct npm artifacts
install a usable Bench runtime, preserve installed-kit versus project-local ownership,
and leave a fresh clone with a runnable pinned local path. Sources: `RR:R-04`,
`RR:R-05`, `RR:R-09`, `RR:A-06`; `RC:C-04`.

## #1: What owns artifact construction?

Type: Grill

### Question

Choose whether wrapper and platform packages keep separate release assembly paths or
share one packaging interface.

### Answer

One repo-only artifact builder is the deep owner of clean staging, explicit
single-sourced allowlists, target/version derivation, build flags, and tarball output
for both the wrapper and all four platform packages. The gate and release automation
call that same interface. Artifact construction is not a public `bench` command.

## #2: How does a fresh clone recover its project-local runtime?

Type: Grill

### Question

The linked binary copy is ignored and does not survive clone. Choose the durable
tracked path and its missing-binary behavior.

### Answer

The tracked `.bench/bin/bench.sh` remains the local entry path. It reads the linked
kit version from the existing manifest stamp and, when its matching native binary is
absent, uses the existing repair/cache seam to fetch that exact platform package into
the user cache. A failed repair exits with an actionable error; it never falls back to
the globally installed runtime.

## #3: Which runtime owns each stable-shim command?

Type: Grill

### Question

Define command-aware routing without sending lifecycle work into the local wrapper's
known refusal path or bypassing a repository's pinned runtime.

### Answer

The stable shim always routes `setup`, `link`, `init`, `doctor`, and `unlink` to its
installed-kit target. Inside a linked repository, every other command routes to the
project-local wrapper; outside one, it routes to the installed target. An existing
local wrapper is authoritative for operational commands, including its fail-closed
repair result.

## #4: Does FT81 implement `bench setup`?

Type: Grill

### Question

The distributable must accept the setup route, while FT76 separately owns the
repo-aware bootstrap behavior.

### Answer

FT81 makes `setup` an installed-kit maintenance route and proves that forwarding
contract. FT76 alone implements the explore-present-confirm-write bootstrap. FT81
does not add a smaller setup command or a second installer seam.

## #5: Where does native platform proof run?

Type: Grill

### Question

Choose whether FT81 merely exposes a future smoke or executes it across the advertised
matrix before merge.

### Answer

FT81 supplies a reusable host-native artifact smoke and four-target automation for
Darwin arm64/x64 and Linux arm64/x64. It runs on pull requests and default-branch
pushes. FT82 later composes the same smoke into the authoritative release preflight.

## #6: Is FT81 one spec or several slices?

Type: Grill

### Question

Choose whether artifact construction, routing, fresh-clone recovery, and target
smokes can close independently.

### Answer

Use one spec. The exact packed tarball and its isolated install flow are the single
acceptance seam; splitting would require duplicated or simulated integration
contracts.

## Not yet specified

n/a — the grill left no unresolved build-shaping fog.

## Out of scope

- FT76's repo-aware `setup` implementation and prompt-driven configuration.
- FT82's authoritative release preflight, toolchain/scanner policy, and tag gate.
- FT83's publication governance, offline evidence, provenance, legal records, SBOMs,
  checksums, and reproducibility bundle.
- FT87's general network/resource policy and independently manifest-pinned repair
  hardening beyond the pinned-version fresh-clone path FT81 needs.
- Publishing packages or creating the first release tag.
- Lifecycle semantic changes not required for the packed install, relink, fresh-clone,
  and unlink flow to use the correct runtime owner.

## Handoff

1. **Module boundaries.** A repo-only artifact builder owns clean wrapper/platform
   staging and tarball production from one explicit asset manifest, the canonical
   platform matrix, and the canonical package version. The shell launcher owns native
   package selection, linked-version discovery, repair/cache lookup, and execution.
   Adoption owns the link-manifest version stamp and the installed lifecycle commands.
   Doctor's generated stable shim owns only command-aware dispatch between its
   installed target and the repository-local wrapper. A reusable artifact-smoke
   interface owns executable mode, embedded version, target format, native execution,
   and selected-package proof. The gate and workflows are thin callers of those
   interfaces.
2. **Contracts.** The artifact builder accepts the source tree and an output location,
   fails nonzero without promoting partial output, and emits one wrapper tarball plus
   one platform tarball per canonical target. The wrapper carries only its explicit
   allowlist and declares the exact-version platform packages; it carries no
   build-host binary or nested platform-package tree. Each platform tarball contains
   one nonempty executable `package/bin/bench` with the package version embedded and
   the declared target format; Linux binaries are static (`CGO_ENABLED=0`) and release
   binaries are stripped. A packed install selects the declared native dependency
   before repair/cache candidates. In a linked clone, `.bench/bin/bench.sh` reads
   `#kit` from `.bench/link-manifest.tsv`, resolves or repairs that exact native
   package into `$BENCH_HOME/cache/bin/<version>/<target>/bench`, preserves argv, and
   exits 127 on unrecoverable absence without global fallback. The stable shim sends
   `setup|link|init|doctor|unlink` to its installed target and all other commands to a
   discovered local wrapper, otherwise the installed target. FT81 proves `setup`
   forwarding, not setup behavior. Exact-tarball tests use an isolated prefix to
   exercise version, link/init, doctor, one local operational command, relink, a fresh
   clone, and unlink.
3. **Deep vs thin.** The artifact builder is deep: callers request artifacts without
   knowing staging, allowlist, build, or pack mechanics. The launcher is deep for
   runtime resolution and repair. Adoption remains deep for managed project state.
   The stable shim is a thin command classifier/dispatcher, and workflow jobs are thin
   orchestration over the artifact-builder and smoke interfaces. Tests attach to
   tarballs, installed commands, and executable behavior rather than staging helpers.
4. **Black-box assertables.** Inspect tarball entry names, modes, sizes, package
   metadata, optional-dependency versions, and forbidden-path absence. Install only
   those tarballs into an isolated prefix and assert command exit codes/output,
   selected executable identity/version, shim target choice, linked file/git state,
   relink idempotency, clone-time cache recovery, no-global-fallback failure, and
   unlink state. Per native target, assert binary magic/architecture, executable mode,
   embedded version, successful `version`, and launcher selection of that exact
   package binary.
5. **Gate attachment.** Attach the primary contract at the existing runtime/package
   surface using artifacts produced by the real repo-only builder: exact tarball
   inventory plus isolated-prefix install/link/init/doctor/operation/relink/clone/
   unlink. Add focused launcher/shim fixtures for routing and repair failures, and
   biting canaries for wrapper contamination and selection-order drift. The reusable
   host-native smoke runs from the gate on its host where available; the four-target
   workflow runs the same smoke on every advertised target. A single local gate
   cannot execute foreign binaries, so the remote matrix is an explicit platform
   blind spot until FT82 composes it into release preflight; root conformance must at
   least fail if the matrix stops invoking the shared smoke.
6. **Hostile-input owners.** Spaces and glob characters in source, staging, prefix,
   repository, and cache paths belong to the artifact-builder and launcher contracts.
   Control bytes in git-sourced operational output remain owned by the existing Go
   command/TOON seams; FT81 must not bypass them. A manifest lacking a trailing
   newline belongs to linked-version parsing. Absent versus empty binaries, manifests,
   optional dependencies, and caches belong to launcher resolution and tarball
   inspection. Special files in script-discovery paths remain with the existing
   static-discovery owner; the artifact builder rejects non-regular allowlist inputs
   before reading. Multi-word argv belongs to stable-shim and launcher `"$@"`
   forwarding. Missing `node`, `npm`, Git, native dependency, global `bench`, or
   non-portable `readlink -f` belongs to isolated-install/fresh-clone fixtures and
   actionable failure paths. Symlink invocation belongs to the existing launcher and
   stable-shim path resolvers. Installed CLI, project-local by-path CLI, hooks, and
   adapters must all resolve through the same local launcher. Destructive worktree
   state remains with existing operational commands; packaging adds no cleanup claim.
   Interrupt during artifact construction or repair must leave no promoted partial
   artifact/cache entry. Repack, reinstall, relink, and repeated repair are
   idempotent. Invocation below the repository root belongs to existing Git-root and
   local-wrapper discovery.
7. **Uncertainty flags.** None on product seams. Hosted-runner labels and availability
   are current external facts the spec-writer must verify against primary provider
   documentation without reopening the accepted four-target, pre-merge smoke policy.
8. **Rejected alternatives.** Rejected separate wrapper/platform builders; publishing
   from a developer's built tree; `npm prepare` placing a host binary in the wrapper;
   committed per-platform binaries in linked repositories; maintainer build scripts
   as the clone remedy; global-runtime fallback after local repair failure; routing
   lifecycle commands through the local wrapper; implementing a reduced `setup` in
   FT81; tag-only/manual platform smokes; deferring all native execution to FT82; and
   splitting FT81 into multiple specs.
9. **Domain watch-outs.** npm lifecycle scripts can mutate the tree being packed, so
   clean staging and post-pack inspection are load-bearing. npm may nest or hoist
   optional dependencies; resolution must support declared install layouts without
   preferring a root build artifact. The tracked local wrapper has no package.json, so
   its version source is the link manifest, not package-layout inference. Cross-built
   format inspection is not native execution; each advertised target must execute its
   own packed binary. Shell dispatch must use `exec` and preserve `"$@"` exactly.
   Version, target matrix, allowlist, and expected artifact count must each have one
   executable source, with tests deriving inputs from those sources rather than
   restating them.

Dependency order: n/a — single spec.
