# FT81 — distributable runtime

Status: staged

Decision map: `decisions/distributable-runtime.md`

## Problem

The npm distribution does not install the same Bench runtime that the repository
claims to ship. Packing from the source tree admits the build host's binary and can
admit generated platform-package trees into the wrapper. The wrapper then prefers
that host binary over the declared native package. Packed installs cannot run the
adoption commands, the generated stable shim sends maintenance back into a local
wrapper that refuses it, and a fresh clone loses the ignored linked binary and points
at a maintainer-only build script.

Consequently, neither the exact npm tarball nor the advertised four-platform matrix
currently proves that a user receives a usable, platform-correct, repository-pinned
runtime.

## Solution

Build the wrapper and all four platform packages through one repo-only artifact
builder that stages only explicit assets, stamps one canonical version, and promotes
only complete tarballs. Install and test those exact tarballs in an isolated prefix.
The installed wrapper selects its declared native package before repair/cache
candidates and can perform adoption. The generated stable shim keeps maintenance on
the installed kit and sends operational commands to the linked repository's tracked
launcher. That launcher reads the linked kit version, repairs the exact native package
into the user cache when necessary, and fails closed without a global-runtime
fallback.

A shared host-native smoke executes the same artifact contract on Darwin arm64,
Darwin x64, Linux arm64, and Linux x64 for pull requests and pushes to the default
branch. The project gate owns the host-executable contract and structural proof that
the remote matrix still calls the shared smoke.

## User stories

1. As a Bench release maintainer, I want one repo-only artifact builder to emit one
   wrapper tarball and one platform tarball for each canonical target from clean,
   explicit staging trees, so that source-tree residue and partial builds cannot
   become release artifacts. Line: gpt-5.6-luna / low. The map fixes the builder seam
   and the artifact contract makes a wrong implementation fully observable at the
   gate.

2. As an npm user, I want the wrapper to contain no build-host binary or nested
   platform-package tree and to select the exact declared native package before any
   repair cache, so that my command executes code built for my host. Line:
   gpt-5.6-luna / low. Tarball inspection and a poisoned-candidate launch fixture
   cheaply expose contamination and resolution-order drift.

3. As a user installing Bench from its packed artifacts, I want `version`, `link`,
   `init`, `doctor`, a project-local operational command, relink, and `unlink` to work
   through an isolated installed prefix, so that the distributed product supports its
   advertised lifecycle without a source checkout. Line: gpt-5.6-luna / low. The
   exact-tarball lifecycle is precise, uses existing adoption seams, and is fully
   gate-observable on the host.

4. As a user with a durable `bench` shim, I want `setup`, `link`, `init`, `doctor`,
   and `unlink` to reach the installed kit while every other invocation reaches the
   project-local launcher when one exists, so that maintenance can repair the kit
   without bypassing a repository's pinned runtime for operations. Line:
   gpt-5.6-luna / low. The command set and default arm are closed decisions exercised
   through a focused dispatcher contract.

5. As a teammate cloning a linked repository, I want the tracked local launcher to
   recover the manifest-pinned native runtime into my user cache, so that hooks,
   adapters, and by-path commands work without a committed binary or a global Bench
   install. Line: gpt-5.6-luna / low. The existing repair seam needs one precise
   version-source change and a clone fixture observes success, cache identity, and
   fail-closed errors.

6. As an agent using any shipped Bench entry surface, I want the installed CLI, stable
   shim, project-local by-path CLI, hooks, and adapters to converge on the same
   command-aware launcher behavior, so that harness choice cannot change runtime
   ownership. Line: gpt-5.6-luna / low. Existing surface contracts provide the seam
   and the packed lifecycle supplies the previously missing distributed-runtime
   subject.

7. As a contributor, I want the exact artifact set exercised natively on Darwin
   arm64, Darwin x64, Linux arm64, and Linux x64 on pull requests and default-branch
   pushes, so that cross-compilation format checks cannot substitute for executing
   what each platform receives. Line: gpt-5.6-terra / medium. Runner execution is
   externally observable but not reproducible inside the local gate, so the weak
   local coverage bumps this story one tier.

8. As the reviewer, I want the project gate and biting canaries to reject wrapper
   contamination, native-selection drift, and removal of the shared matrix smoke, so
   that a green gate remains evidence for the distributable contract rather than a
   source-tree approximation. Line: gpt-5.6-terra / medium. This changes the oracle
   and its omission proof, matching the profile's cached gate/conformance routing.

The implementation is one integration sequence because the exact artifact set is the
acceptance subject. If story lines collapse into one authoring pass, the pass uses the
highest line above: `gpt-5.6-terra / medium`.

## Implementation decisions

- **Artifact builder.** One repo-only deep module accepts the source tree and an
  output location, derives version and targets from their canonical machine-readable
  sources, builds in private staging, and emits one wrapper tarball plus one platform
  tarball per target. The gate and workflows are thin callers of this interface.
  Failure or interruption must not promote a partial output set. Re-running must be
  safe; byte-for-byte reproducibility evidence remains FT83's responsibility.
- **Single-sourced artifact policy.** One explicit asset manifest owns the wrapper
  allowlist and modes. The canonical platform matrix owns target package names, Go
  targets, and runner labels. The wrapper package metadata owns the version. Artifact
  count is derived as one plus the matrix length. Tests and workflows consume these
  sources instead of restating lists or counts.
- **Wrapper artifact.** The wrapper tarball declares exact-version optional
  dependencies for all matrix targets. Its explicit allowlist includes the source
  assets adoption needs, but excludes repo-only state, every root or nested platform
  binary, and every generated platform-package tree. Packing never depends on an npm
  lifecycle script mutating the source tree.
- **Platform artifacts.** Each platform tarball contains package metadata and exactly
  one nonempty executable at `package/bin/bench`. The embedded Bench version equals
  the wrapper version. Format and architecture match the matrix. Linux artifacts are
  static `CGO_ENABLED=0` builds, and all release artifacts are stripped. Release build
  flags have one executable source behind the artifact builder.
- **Runtime resolution.** The installed wrapper resolves the matching declared native
  package before repair/cache candidates and never prefers a root build artifact.
  npm's supported nested and hoisted optional-dependency layouts remain accepted.
  Empty or non-executable candidates are absent, not runnable.
- **Linked runtime.** Adoption stamps the linked kit version in the existing manifest
  and leaves the tracked local launcher as the durable entry path. A linked launcher
  with no package metadata reads `#kit`, including when the manifest lacks a trailing
  newline. It resolves or repairs exactly
  `$BENCH_HOME/cache/bin/<version>/<target>/bench`, promotes only a complete executable,
  and `exec`s it with the original argv. Missing repair prerequisites or an
  unavailable package produce an actionable exit 127 and never search for a global
  `bench`.
- **Stable-shim dispatch.** The generated stable shim is a thin classifier. Its
  explicit maintenance set is `setup`, `link`, `init`, `doctor`, and `unlink`; those
  commands always `exec` the installed target. Its default arm sends all other argv
  to the discovered local wrapper, including from a nested directory or linked
  worktree, and uses the installed target only when no local wrapper exists. `setup`
  forwarding is the whole FT81 setup contract; FT76 owns setup behavior.
- **Exact-tarball lifecycle.** Tests install only the builder's wrapper and host-native
  platform tarballs into an isolated prefix. They exercise, in order, installed
  `version`; `link`; `init`; `doctor`; stable-shim setup forwarding; one local
  operational command; relink; a committed fresh clone with an empty cache and no
  global `bench`; and installed-kit `unlink`. The fixture asserts file, executable,
  manifest, Git, cache, and exit/output state after each transition.
- **Native smoke.** One repo-only host-native smoke consumes the artifact directory,
  installs its wrapper and matching platform tarball into an isolated prefix, poisons
  lower-priority candidates, checks artifact format/mode/version, runs `bench version`,
  and proves the chosen package executable. The local gate calls it for the current
  host when available. A verification workflow builds one exact artifact set, derives
  its matrix dynamically from the canonical platform source, and invokes that same
  smoke on these currently documented standard GitHub-hosted runner labels:
  `macos-15`, `macos-15-intel`, `ubuntu-24.04-arm`, and `ubuntu-24.04`. GitHub's
  official hosted-runner reference was checked on 2026-07-15:
  <https://docs.github.com/en/actions/reference/runners/github-hosted-runners>.
- **Oracle attachment.** Runtime/package contract tests own exact tarball inspection,
  the isolated lifecycle, launcher repair, and shim routing. Root conformance proves
  the PR/default-branch workflow derives its matrix and calls the shared smoke.
  Behavior-owned canaries mutate wrapper contamination and native selection so each
  enforcement can be shown to bite. The four remote native executions remain an
  explicit local-gate blind spot until FT82 composes the same smoke into release
  preflight.
- **Deep and thin modules.** Artifact construction and launcher resolution remain
  deep. Adoption remains the deep owner of managed repository state. Stable-shim
  classification, artifact-smoke orchestration, the gate caller, and workflow jobs
  stay thin. Tests attach to tarballs, installed commands, linked executable behavior,
  and workflow structure, never staging helpers.

## Testing decisions

- A good test treats the emitted tarballs as the product. It inspects archive entries
  and metadata, installs only those tarballs, and drives public commands and generated
  launchers while observing bytes, exit codes, selected executable identity, and
  filesystem/Git state.
- The primary seam extends the existing runtime/package contracts. Prior art includes
  the package generator and npm-pack contract, safe-link contract, binary-repair
  fixtures, Go routing contract, doctor-shim contract, and unlink contract.
- Focused launcher and stable-shim seams are justified in addition to the primary
  tarball seam because repair refusal and command-owner mistakes cannot be isolated
  reliably after the full lifecycle has already failed.
- The native-smoke seam is separately justified because one host cannot execute all
  four native formats. It keeps local and remote checks on the same interface while
  root conformance detects workflow amputation.
- The feature must pass the project gate: `.bench/gate.sh`.
- Pre-implementation probes were run on 2026-07-15 and are named `R1` through `R8`
  below. They are concrete current reds, not predicted test names:
  - `R1` exact `npm pack --dry-run --json` inspection exited 1 because the wrapper
    contained `dist/bench`.
  - `R2` generated-platform inspection exited 1 because all four packages lacked a
    runnable `bin/bench`.
  - `R3` isolated install of the exact wrapper tarball exited 1 on `bench link` with
    `no source asset tree`.
  - `R4` generated stable shim invoked inside a linked repo exited 1 for `doctor`
    because it executed the planted local wrapper instead of the installed target.
  - `R5` a committed linked repo cloned without ignored files exited 127 and advised
    the unavailable maintainer script instead of repairing from `#kit`.
  - `R6` a launcher fixture with a failing root build and a successful native package
    exited 1, proving the root build was selected first.
  - `R7` the workflow-structure probe for PR triggers and the four current runner
    labels exited 1 because no native verification workflow exists.
  - `R8` the canary-anchor probe exited 1 because no wrapper-contamination or
    native-selection canary exists.

### Seam diagrams

Seam 1 — exact artifact set and installed lifecycle:

    trigger: gate or verification workflow requests FT81 artifacts
        │
        ▼
    source tree + output dir ──▶ [ repo-only artifact builder ] ──▶ wrapper.tgz
    version + target matrix   ──▶ [ clean stage, build, pack   ] ──▶ 4 platform.tgz
                                      │
                                      ▼
                              isolated npm prefix + linked repo
                                      │
                                      ▼
                         commands, files, Git state, cache, exits
                  ◀ tests attach here: inspect the exact tarballs, install only
                    those tarballs, and drive the complete lifecycle as a user

Seam 2 — linked launcher resolution and repair:

    trigger: project-local command, hook, or adapter invokes tracked bench.sh
        │
        ▼
    argv + #kit manifest ──▶ [ local launcher resolver ] ──▶ native package binary
    package/cache state  ──▶ [ repair + atomic cache   ] ──▶ or actionable exit 127
                                    │
                                    ▼
                          selected path, output, cache state
                 ◀ tests attach here: clone fixture controls manifest, registry,
                   cache, PATH, symlinks, cwd, and candidate executables

Seam 3 — generated stable-shim command dispatch:

    trigger: user invokes durable `bench` shim
        │
        ▼
    argv + cwd/repo state ──▶ [ maintenance classifier ] ──▶ installed kit target
    installed/local paths ──▶ [ default local dispatch ] ──▶ project-local launcher
                                      │
                                      ▼
                             target identity, argv, exit code
                    ◀ tests attach here: planted installed/local targets make
                      every maintenance command and default-arm case observable

Seam 4 — shared native artifact smoke:

    trigger: host gate or one job per canonical matrix row
        │
        ▼
    exact artifact dir + host ──▶ [ shared host-native smoke ] ──▶ format/mode verdict
    poisoned lower candidates  ──▶ [ install + execute       ] ──▶ version/identity
                                           │
                                           ▼
                                 one verdict per native target
                     ◀ tests attach here: local gate runs the host row; root
                       conformance proves all four remote jobs invoke this seam

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | One invocation emits exactly one wrapper tarball plus one tarball per canonical matrix row, with count derived rather than hardcoded | exact artifact set | R2: the current generator probe exits 1 with all four runnable artifacts missing | A metadata-only generator or a builder that omits any target cannot satisfy the derived inventory |
| 1 | Wrapper and platform staging accepts source and output paths containing spaces and glob characters and copies only regular allowlisted inputs | exact artifact set | R1: current source-tree packing exits 1 on forbidden residue and has no clean-stage interface | A builder that packs the source tree, expands globs, or reads a FIFO or device fails inventory, timeout, or regular-file assertions |
| edge of 1 | Build failure or interruption promotes no partial artifact set; a re-run after failure and a second successful repack are safe | exact artifact set | R2: the current generator leaves incomplete package directories and has no all-or-nothing promotion | Failure injection catches direct-to-output writes, while repeat invocation catches stale or duplicate staging state |
| 1 | Version, target names, runner labels, wrapper allowlist, release flags, and artifact count each have one executable source consumed by builder, tests, and workflows | artifact contract plus root conformance | R7: current workflow hardcodes release orchestration and has no shared builder or derived native matrix | Deleting or duplicating a source makes conformance or the artifact inventory disagree instead of silently drifting |
| 2 | Wrapper tarball entries, modes, and sizes match the explicit allowlist; `dist/bench`, nested platform packages, repo-only files, and other residue are absent | exact tarball inspection | R1: exact pack probe exits 1 and names `dist/bench` | Post-pack entry inspection rejects the cheapest fix that merely changes launcher priority while still shipping contamination |
| 2 | Wrapper metadata declares the four matrix-derived optional platform packages at exactly the wrapper version | exact tarball inspection | R1: current pack is contaminated and no clean artifact contract couples metadata to emitted platform tarballs | Missing, extra, or stale dependency metadata fails comparison against the same emitted matrix set |
| 2 | Every platform tarball contains one nonempty executable `package/bin/bench` with matching embedded version, target magic and architecture, stripped symbols, and no dynamic Linux dependency | exact platform tarball inspection | R2: current generated packages are missing all four binaries | Empty stubs, host-only copies, wrong cross-target output, unstripped builds, and dynamic Linux builds each fail a black-box artifact assertion |
| 2 | Installed launcher chooses the declared native package before a root build, nested or hoisted alternatives are accepted, and empty or non-executable candidates are skipped | launcher resolution fixture | R6: planted root failure wins over the successful native package and exits 1 | Poisoning each higher and lower candidate identifies the executable actually selected and catches an order-only fake fix |
| 3 | Installing only exact wrapper and host-platform tarballs into an isolated prefix runs `version`, `link`, `init`, `doctor`, one local operational command, relink, and `unlink` with their expected exits and output | exact installed lifecycle | R3: the exact packed wrapper installs but `link` exits 1 with `no source asset tree` | Source-tree-only success cannot pass because every command is driven through the installed tarball |
| 3 | Link and init preserve project-owned text, stamp the kit version, install executable managed surfaces without a committed platform binary, and leave expected Git state | exact installed lifecycle | R3: the packed install cannot reach link, so none of the managed state exists | File hashes, manifest rows, executable modes, ignored-state checks, and Git status expose partial or overbroad adoption |
| edge of 3 | Relink, repeated init, reinstall, and unlink are idempotent; unlink removes only owned state and reports residuals through existing postures | exact installed lifecycle | R3: no first packed link exists, so the repeated lifecycle cannot complete | Before/after hashes and second-run exits catch duplicate blocks, stale ownership, and broad deletion |
| edge of 3 | Isolated lifecycle works in prefix and repository paths containing spaces and glob characters | exact installed lifecycle | R3: current packed adoption fails before the hostile paths can be exercised | Driving the public commands from hostile paths catches unquoted staging, install, shim, and adoption arguments |
| 4 | Each maintenance command, enumerated as `setup`, `link`, `init`, `doctor`, and `unlink`, always executes the installed target even when a local wrapper exists | stable-shim dispatch | R4: `doctor` executes the planted local failure and exits 1 | Per-command target markers catch any omitted maintenance member; setup forwarding is proven without implementing setup behavior |
| 4 | Default-arm cases including no args, flags, a known operational command, and an unknown command execute the local wrapper in a repo and the installed target outside one | stable-shim dispatch | R4: the current shim has no maintenance classifier and therefore cannot satisfy both target classes | Paired planted targets make always-local and always-installed implementations fail one side of the table |
| 4 | Dispatch from a nested cwd and linked worktree finds the authoritative local wrapper, while a local self-link or absent local wrapper falls back to the installed target | stable-shim dispatch | R4: current maintenance misrouting demonstrates the classifier is absent; existing linked-worktree routing remains covered by `testGoRoutingShimLinkedWorktree` | Root, common-dir, absence, and self-resolution fixtures catch recursion and cwd-only lookup |
| edge of 4 | Stable shim preserves multi-word, glob-shaped, empty, and ordinary argv byte-for-byte and returns the selected target's exit code; symlink invocation remains supported | stable-shim dispatch | Already covered by `testDoctorShimArgPassthrough`, `testGoRoutingRepoLocalWrapperForwarding`, and wrapper symlink contracts; the maintenance matrix extends the same seam | Independent target output and exit markers fail on `$*`, lost empty args, missing `exec`, or wrong path resolution |
| 5 | A fresh clone with no package metadata or local binary reads exact version from `#kit`, including a final line without newline, fetches the matching native package, and executes it | linked launcher repair | R5: clone invocation exits 127 and advises a maintainer-only build script | Removing package metadata forces the manifest path; registry version and executed banner catch an unpinned or global substitute |
| 5 | Repair promotes the verified binary only to `$BENCH_HOME/cache/bin/<version>/<target>/bench`; absent, empty, and non-executable package or cache candidates recover, and the second run makes zero registry requests | linked launcher repair | R5: the clone has no usable version source, so it cannot create or reuse the required cache entry | Path, mode, size, request-count, and second-run assertions catch torn promotion, wrong cache keys, and refetch loops |
| edge of 5 | Missing node, unavailable or malformed package data, disabled repair, and an interrupted repair exit 127 with an actionable message, no promoted partial entry, and no global-runtime fallback | linked launcher repair | R5: current failure exits 127 with the wrong maintainer-only remedy; existing repair integrity and torn-cache tests are already covered | A planted global success executable plus an empty cache makes fallback or partial promotion observable |
| edge of 5 | Fresh-clone repair works with repository and cache paths containing spaces and glob characters, without npm, without global `bench`, and without GNU `readlink -f` | linked launcher repair | R5: current clone fails even in an ordinary path | A minimal PATH and hostile directories catch hidden npm, global-command, GNU-only, and quoting dependencies |
| 6 | Installed CLI and stable shim reach the installed owner for maintenance; by-path CLI, session hook, stop hook, and shipped adapters reach the same linked launcher for operations | shipped entry surfaces | R3 and R5: installed adoption and cloned by-path execution both fail before all surfaces can converge | Target identity markers across every shipped caller catch a harness-specific bypass or direct core invocation |
| 6 | Existing control-byte refusal and destructive-worktree safety remain owned by the Go command and operational seams rather than being reimplemented in packaging | existing Go command surfaces | Already covered by AXI control-byte, shift-objective, gate-proof, and worktree contracts | Keeping those tests green catches a packaging route that bypasses the canonical Go behavior without duplicating its policy |
| 7 | The shared smoke runs the exact host platform tarball natively, verifies mode, embedded version, target format, successful `version`, and selection over a poisoned cache | shared native smoke | R2: no generated target is runnable, and R7 finds no workflow caller | Cross-format inspection alone cannot pass the executable invocation or selected-identity checks |
| 7 | One matrix row runs on each of `macos-15`, `macos-15-intel`, `ubuntu-24.04-arm`, and `ubuntu-24.04`, derived from the canonical platform source | native verification workflow | R7: runner-label and PR-trigger probe exits 1 | A missing target, duplicated host, or hardcoded divergent matrix fails derived-row conformance or leaves a required job absent |
| 7 | Native verification triggers on pull requests and pushes to the default branch and every row invokes the shared smoke | native verification workflow plus root conformance | R7: current workflow is tag-only and has no shared smoke | Structural conformance catches trigger or call-site amputation even though the local gate cannot execute remote jobs |
| 8 | Host gate builds the exact artifact set, inspects every tarball, runs the isolated lifecycle, and invokes the shared smoke for its native host | gate contract | R1, R2, and R3 are all red while the pre-FT81 gate does not own the combined artifact subject | A gate that checks only source shape, only metadata, or only commands leaves at least one current red uncaught |
| 8 | Wrapper-contamination and native-selection mutations each make the inner gate red with a targeted diagnostic | behavior-owned canaries | R8: no corresponding canary anchor exists | Each mutation removes the cheapest way an enforcement can decay into an always-pass without detection |

Degenerate implementations checked against the map:

- A builder that runs `npm pack` in the source tree fails stories 1 and 2 on wrapper
  inventory and hostile-path staging.
- A builder that emits metadata-only platform packages fails story 2 on executable,
  format, version, and native execution.
- A lifecycle test that silently substitutes the source checkout fails story 3 because
  the installed prefix and exact tarball paths are asserted.
- An always-local or always-installed stable shim fails opposite halves of story 4.
- A linked launcher that still reads package metadata, invokes a global `bench`, or
  writes directly into cache fails story 5's package-less clone and poisoned fallback.
- A harness-specific direct-core path fails story 6's selected-wrapper identity.
- A workflow that only cross-inspects artifacts fails story 7's native `version` call.
- An enforcement without an omission oracle fails story 8's mutation canaries.

### Edge inventory

The canonical error, empty/absent, boundary, malformed, interrupted/partial,
re-run-idempotency, and hostile-environment classes were walked for every mapped
behavior. The profile's shell checklist lands as follows:

- Paths with spaces or glob characters — coverage rows for stories 1, 3, and 5.
- Control bytes in Git-sourced text — story 6's already-covered Go-command row.
- Manifest without a trailing newline — story 5's fresh-clone row.
- Absent versus present-but-empty files — stories 2 and 5 candidate/cache rows.
- Special files in discovery or allowlist paths — story 1's regular-file row.
- Multi-word and glob-shaped argv — story 4's argv row.
- Missing node, npm, Git, global `bench`, or GNU `readlink -f` — story 5's failure and
  minimal-PATH rows; the installed maintenance target remains usable without Git, and
  the by-path clone route remains usable without npm or global Bench.
- Symlink invocation — story 4's argv/path row and existing launcher contracts.
- Every shipped surface — story 6's surface row.
- Interrupt during artifact construction or repair — stories 1 and 5 partial-state
  rows.
- Repack, reinstall, relink, repeated init, and repeated repair — stories 1, 3, and 5
  idempotency rows.
- Invocation below the repository root and from a linked worktree — story 4's
  discovery row.
- **Won't handle:** control-byte behavior inside artifact names — the explicit asset
  manifest and canonical package names are repository-authored constants, not
  Git-sourced user text; no primary caller is amputated.
- **Won't handle:** destructive worktree recovery inside packaging — FT81 adds no
  cleanup or worktree mutation path; existing operational commands and their gate
  contracts retain ownership.
- **Won't handle:** byte-for-byte reproducibility across separate builds — FT83 owns
  the reproducibility evidence bundle; FT81 requires safe repeatability and equivalent
  artifact contracts, not identical gzip timestamps.
- **Won't handle:** executing a foreign target in the local gate — the shared smoke
  executes the host row locally and the other three rows on native hosted runners;
  root conformance makes remote-smoke removal locally observable.

Compatibility probe: npm is the consuming implementation. The isolated lifecycle
uses npm itself to install the exact wrapper and native tarballs and covers both
supported optional-dependency layouts through launcher fixtures; no tar or package
format compatibility is accepted from a hand-written parser alone.

## Out of scope

1. **FT76 repo-aware `bench setup` behavior** — prompt-driven exploration,
   confirmation, and transactional project bootstrap are a separate adoption
   capability; estimated 10 edits, 4 gate runs.
2. **FT82 authoritative release preflight** — toolchain/scanner policy, version/tag
   equality, ancestry, and publication authorization are a separate release-gate
   capability; estimated 10 edits, 4 gate runs.
3. **FT83 governed and offline-verifiable release bundle** — archives, offline
   installation, provenance, notices, SBOMs, checksums, digest-safe publication, and
   reproducibility evidence form a separate artifact-governance capability; estimated
   14 edits, 5 gate runs.
4. **FT87 general repair/network hardening** — independent manifest pinning, fetch
   deadlines, response/decompression limits, concurrency policy, and broader resource
   controls are a separate runtime-hardening capability; estimated 8 edits, 3 gate
   runs.
5. **Publishing packages or creating the first release tag** — registry mutation and
   immutable release creation are a separate reviewer-authorized operation after the
   release preflight and governance specs; estimated 4 edits, 2 gate runs.
