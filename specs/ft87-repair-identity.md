# FT87 slice 2: hardened repair, one identity, and the offline evidence record

Status: implemented

Map: `decisions/bounded-network-resource-cli.md` (tickets #1, #3, and #9, and
the Handoff), closed at `abd4a1d`. Slice 1 shipped at `d8a4ff8`.

## Problem

Binary repair is the one Bench operation that runs before any Bench binary
exists, and it is the least governed thing the product does. It starts as a
silent side effect of a failed resolution, so a user who expected a local tool
gets a network fetch they never asked for. It trusts the registry metadata it
just downloaded for the digest it verifies against, which means the transport
it is defending against is also the authority on what "correct" means. It has
no fetch deadline, no download cap, and no decompression cap, so a hung or
hostile endpoint can stall or exhaust the machine. Its failure path deletes the
installed binary itself rather than only its own temp file, so a losing racer
can destroy a winner's good install. Its cache grows without bound and nothing
ever prunes it.

Separately, three names ship for one product — `bench` is the command,
`redbench` is the npm package, and `benchkit` is what `bench version` and the
platform package descriptions print — so a user cannot tell from the output
which thing they installed. The root and generated platform `package.json`
files carry no `repository`, `homepage`, `bugs`, or `author`, so nothing in the
published artifact points back at the project.

Finally, the FT83 requirement registry declares
`public.ft87.offline_network_control` as a required producer record for the
public and bank release profiles, and nothing writes it. The offline smoke
harness proves zero unwanted network attempts today, but it produces no stable
record, so every release inspection reports the requirement as missing.

## Solution

Repair becomes an explicit, pinned, bounded, and concurrency-safe operation.
Resolution failure names `bench repair` instead of performing it; `BENCH_REPAIR=1`
restores implicit repair for automation. The expected digest moves out of
registry `dist.integrity` into a manifest generated at build time from the
actual packed tarballs and shipped inside the wrapper package, so a mismatch is
a hard red decided by a pin the transport cannot influence. Fetch, download, and
decompression each gain a bound whose breach is a distinct failure. Failure
cleanup removes only files this process created. Cache pruning becomes the
explicit `bench repair --prune`.

The user-facing identity becomes Bench everywhere output is read by a person:
`bench version` prints `bench <version> (<os>/<arch>)`, and `benchkit` survives
only as this repository's internal project-profile name. The root and every
generated platform package gain the four metadata fields, single-sourced from
the root so the platform packages cannot drift.

The offline proof gains a producer: the native/offline workflow writes the
`ft87/offline-network-control/v1` envelope that the FT83 registry already
declares, turning the existing zero-attempt sentinels into the stable release
evidence record.

## User stories

1. As an operator whose binary is missing, I want resolution failure to name
   the explicit `bench repair` action and start no network operation, while
   `BENCH_REPAIR=1` restores the implicit attempt for automation environments,
   so a routine command never fetches on my behalf without being asked.
   Line: gpt-5.6-luna / medium. The wrapper routing seam and the exact failure
   text are fixed by this spec and covered by existing surface contracts, while
   the default flip touches every porcelain subcommand.

2. As an operator repairing deliberately, I want `bench repair` to install the
   pinned binary for this platform and `bench repair --prune` to remove cache
   entries other than the current version and platform, so both the fetch and
   the deletion are acts I chose rather than side effects hidden in a resolver.
   Line: gpt-5.6-luna / medium. The subcommand is a shell-owned route beside
   the existing pre-binary commands, and the prune target set is exactly
   enumerable, but deletion semantics justify medium effort.

3. As a user installing Bench, I want repair to verify the downloaded tarball
   against a digest shipped inside the wrapper package rather than one read
   from the registry response, and to fail hard on absent, empty, malformed, or
   mismatched pins, so the transport being defended against is not the authority
   on what is correct.
   Line: gpt-5.6-terra / medium. This is the security-load-bearing change in
   the slice and its failure modes must all be closed rather than only the
   mismatch case.

4. As a release builder, I want the pin manifest generated from the actual
   packed platform tarballs and shipped in the wrapper package, so the pins
   describe the artifacts that were really published and the build fails rather
   than shipping a wrapper without them.
   Line: gpt-5.6-terra / medium. It inserts a step into an artifact build that
   already asserts reproducibility and tarball counts, so build ordering and
   determinism both have to hold.

5. As an operator on a hostile or hung network, I want repair to stop at a
   60-second total deadline, refuse a download past 100 MB, and refuse a
   decompressed payload past 200 MB, each with its own named failure, so a
   stalled or oversized endpoint cannot hang the command or exhaust the disk.
   Line: gpt-5.6-luna / medium. The three bounds attach at known points in one
   script, and fixture endpoints make each one observable.

6. As one of several Bench processes racing to repair, I want a failing process
   to remove only the temp file it created and never the installed target, so a
   loser cannot delete a winner's good binary.
   Line: gpt-5.6-terra / medium. The current outer handler removes the shared
   target, and the correct behavior has to hold across interrupt as well as
   error paths.

7. As a user reading Bench output, I want `bench version` to print
   `bench <version> (<os>/<arch>)` and no user-facing string to say `benchkit`,
   so one product name appears everywhere a person reads.
   Line: gpt-5.6-luna / low. The strings are enumerable and every consumer of
   them is an existing test or smoke assertion.

8. As a consumer of the published packages, I want the root and every generated
   platform `package.json` to carry `repository`, `homepage`, `bugs`, and
   `author`, derived from the root rather than restated per platform, so the
   artifact points back at the project and the platform packages cannot drift.
   Line: gpt-5.6-luna / low. The generator already derives version and license
   from the root package and the assertion seam is an existing conformance test.

9. As a release inspector, I want the offline proof to write the
   `ft87/offline-network-control/v1` producer record the FT83 registry declares,
   with a payload naming every suppressed operation class and the sentinel
   result, so `public.ft87.offline_network_control` stops reporting as missing.
   Line: gpt-5.6-terra / medium. It is the first producer record any owner
   writes, so the envelope helper it establishes is reused by FT88 and FT71.

10. As an agent picking up cold, I want `bench repair` present in the CLI
    inventory and the repair posture described where the pre-binary behavior is
    documented, so the new explicit action is discoverable without reading the
    wrapper.
    Line: gpt-5.6-luna / low. The inventory is a fixed list the gate already
    sweeps for stale command references.

## Implementation decisions

- Map-carried veto surface is locked here, not left open. Repair becomes
  explicit by default with `BENCH_REPAIR=1` as the opt-in (map #3). The bounds
  are 60 seconds total, 100 MB downloaded, 200 MB decompressed (map #3). The
  version line is `bench <version> (<os>/<arch>)` (map #9). The canonical URLs
  are `https://github.com/gibbonmi/bench` with `git+https://github.com/gibbonmi/bench.git`,
  `https://github.com/gibbonmi/bench#readme`, `https://github.com/gibbonmi/bench/issues`,
  and author `gibbonmi`, supplied by the reviewer at spec time against Handoff
  uncertainty flag (b). A later veto changes the owning fact rather than adding
  a second override.

- `route_porcelain` stops passing `--repair`. `route_binary`'s not-found case
  attempts repair only when `BENCH_REPAIR=1` is set exactly, matching the
  existing exact-value convention slice 1 established for `BENCH_OFFLINE`.
  Without it, the case prints a message naming `bench repair` and exits 127 as
  today. `BENCH_OFFLINE=1` continues to suppress repair and to say so, and it
  outranks `BENCH_REPAIR=1`; `BENCH_NO_REPAIR` remains the narrower lever and
  also outranks it. Precedence is stated once in the wrapper and asserted, so
  no caller infers it.

- `bench repair` is a shell-owned subcommand routed beside the existing
  pre-binary commands, because it must run when no binary exists. It accepts
  exactly `bench repair` or `bench repair --prune`; any other argument is a
  usage error exiting 2. It runs the same Node script the implicit path runs,
  with an explicit mode argument, so there is one repair implementation.

- The pin manifest replaces `dist.integrity` as the digest authority. It maps
  package name and version to the expected `sha512` integrity string per
  platform, and the repair script verifies the downloaded tarball against it
  before extraction. A manifest that is absent, empty, unparseable, or missing
  an entry for the requested name and version is a hard failure that names the
  manifest — the script never falls back to registry metadata, because a
  fallback reintroduces exactly the trust this story removes. Registry metadata
  is still read for the tarball URL; it is no longer read for the digest.

- The manifest's path inside the wrapper package is a build-time constant
  declared once in `internal/releaseevidence/requirements.json` alongside the
  other artifact facts. `build-release-evidence.mjs` includes that path in the
  wrapper `files[]` when it writes the wrapper `package.json`, before the
  content exists. A new step in `scripts/build-artifacts.sh`, running after the
  platform tarballs are packed and before the wrapper is packed, computes each
  tarball's integrity and writes the manifest. Packing the wrapper without a
  present, non-empty manifest fails the build. The repair script's own copy of
  that path is the necessary cross-runtime constant — the script runs where no
  Go binary and no build tooling exists — and conformance pins the two against
  each other, following the shell/Go `BENCH_OFFLINE` parity precedent from
  slice 1 rather than inventing a second registry.

- The manifest is derived from the packed artifacts, never hand-maintained, so
  the existing independent reproducibility build must produce a byte-identical
  manifest. Ordering is deterministic by platform name.

- The repair script gains three bounds it owns as its own facts, because it
  cannot import `internal/bounds`: a 60-second deadline covering metadata fetch
  and tarball download via `AbortSignal.timeout`; a 100 MB cap on downloaded
  bytes; and a 200 MB cap on decompressed bytes, enforced incrementally during
  decompression rather than after it, so an oversized payload is refused before
  it is fully materialized. Each bound produces a distinct message naming the
  bound it hit. The map's domain watch-out applies: these are a documented
  second runtime, not a drifted copy of the Go policy, and conformance must not
  try to unify them.

- Failure cleanup removes only paths this process created. The outer handler
  stops removing the installed target; it removes the process's own temp file
  if one exists, and nothing else. The interrupt handlers keep their existing
  temp-file removal and exit codes. Promotion stays write-to-unique-temp,
  fsync, atomic rename, which is already correct.

- `bench repair --prune` removes entries under the binary cache root other than
  the current version and platform directory, and reports what it removed. It
  removes only regular files and directories it owns under that root, refuses
  to follow symlinks, and is a no-op with an explicit report when nothing is
  stale. It never runs implicitly.

- The version line becomes `bench <version> (<os>/<arch>)`. Every user-facing
  `benchkit` string is retired: the version line, the wrapper help text, and the
  generated platform package description, which becomes `Bench prebuilt binary
  for <os>-<arch>`. `projects/benchkit.md` keeps its name as this repository's
  internal profile, and the shipped asset list keeps referring to it by that
  path. Conformance extends its existing stray-name sweep to fail on `benchkit`
  in user-facing output strings while permitting the profile path and this
  repository's internal references.

- The four metadata fields are added to the root `package.json`. The wrapper
  package already spreads the root package, so it inherits them. The platform
  package literal in `build-release-evidence.mjs` derives all four from the root
  package the way it already derives version and license, so the fields have one
  source. Conformance asserts the four fields on the root and on every generated
  platform package, and asserts the platform values equal the root values.

- The producer record is written by a small shared envelope writer that takes a
  key, owner, schema, identity, status, and payload and emits the
  `producerEnvelope` shape the release requirements already validate, including
  the payload digest. FT87 is its first caller; FT88 and FT71 adopt it when they
  ship, so the envelope mechanics have one owner from the start. The FT87
  payload names each suppressed operation class from slice 1 — wrapper binary
  repair, worktree Git refresh, the Codex discovery subprocess including its
  bundled fallback, the OpenAI models request, and the Anthropic models request
  — with the sentinel's observed attempt count for each, plus the journeys run
  and the flag exercised. The record is written by the native/offline workflow
  where the sentinels run, not by the artifact build, because the proof is
  workflow-attached; it is a producer record and therefore not packaged into any
  tarball.

## Testing decisions

- A good test here drives the shipped wrapper and repair script against a
  hermetic local registry and fixture tarballs, and observes exit code, output
  text, and filesystem state. The existing offline and repair fixture harnesses
  already provide the loopback registry and fake-endpoint machinery; new cases
  extend them rather than adding a second harness.

- Seams, all pre-agreed in the map Handoff. `bin/bench.sh` routing and the
  `bench repair` subcommand are driven as the CLI. `bin/bench-repair-binary.mjs`
  is driven through the wrapper with fixture registries and endpoints. The
  artifact build is driven as a build and its output tarballs inspected. The
  release evidence envelope is asserted through the existing requirement
  inspection. No new seam is invented.

- Prior art: `internal/contract/surface/binary_repair_test.go` owns the repair
  contract, with `binary_repair_offline_test.go` and
  `binary_repair_fixture_test.go` as siblings; `internal/conformance/package_shipped_surface_test.go`
  owns package shape and the stray-name sweep; `internal/contract/surface/artifact/artifact_test.go`
  owns tarball inspection; `internal/releaseevidence/package_artifact_test.go`
  owns evidence-record validation; `cmd/bench/main_test.go` owns the version
  line; `scripts/smoke-offline.sh` owns the zero-attempt sentinels.

- Handoff assertable disposition for this slice is complete. Oversized fixture
  tarballs, the version string, and package metadata have rows below. The
  stable evidence record has a row. The `BENCH_OFFLINE` zero-egress proof
  itself remains workflow-attached and shipped in slice 1; this slice adds only
  the record it produces. Trailing-garbage, `--`, help-exit, root-anchored
  coverage, and capability-skip assertables belong to slice 3. Those are stated
  exceptions, not silently missing tests.

- Gate command: `.bench/gate.sh`.

### Seam diagram

    trigger: operator runs any porcelain subcommand with no binary present
        │
        ▼
    argv + BENCH_REPAIR/OFFLINE/NO_REPAIR ──▶ [ bin/bench.sh route_binary ]
    resolved kit + wrapper path           ──▶ [ resolve, decide, never fetch ]
                                          ──▶ exit 127 + text naming `bench repair`
                                              ◀ tests attach here: run the wrapper in a
                                                binary-less fixture kit and assert exit,
                                                text, and that no child process started

    trigger: operator runs `bench repair` or `bench repair --prune`
        │
        ▼
    mode + version + platform ──▶ [ bin/bench-repair-binary.mjs ]
    fixture registry + tarball ──▶ [ pin verify, bounded fetch, promote, prune ]
                                ──▶ installed binary or named failure + cache state
                                    ◀ tests attach here: hung/oversized/mismatched
                                      fixture endpoints, and cache trees inspected
                                      before and after

    trigger: scripts/build-artifacts.sh packs the release
        │
        ▼
    platform tarballs + root package ──▶ [ pin manifest builder ]
                                      ──▶ wrapper package containing the manifest
                                          ◀ tests attach here: unpack the wrapper
                                            tarball and compare pins to the packed
                                            platform tarballs, twice for repro

    trigger: native/offline workflow runs the offline smoke
        │
        ▼
    BENCH_OFFLINE=1 + egress sentinels ──▶ [ envelope writer ]
                                        ──▶ release-evidence/ft87-*.json
                                            ◀ tests attach here: requirement
                                              inspection over a fixture record

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a porcelain subcommand with no binary exits 127, prints text naming `bench repair`, and starts no node or network child | bench CLI | extended surface contract with a child-process marker, run first against `d8a4ff8`, records the implicit repair child | a message-only change that still calls `repair_binary` writes the marker |
| 1 | exactly `BENCH_REPAIR=1` restores the implicit attempt; unset, empty, `0`, and `true` do not | bench CLI | new value matrix, run first, fails because the variable is unread | truthy or nonempty parsing fails at least one enumerated non-`1` case |
| 1 | `BENCH_OFFLINE=1` and `BENCH_NO_REPAIR` each suppress repair with their own message even when `BENCH_REPAIR=1` is set | bench CLI | new precedence matrix over all four variable combinations, run first, fails because `BENCH_REPAIR` does not exist | an implementation that checks `BENCH_REPAIR` first fails the two suppression rows |
| 2 | `bench repair` installs the pinned binary in a binary-less fixture kit and a following command runs from it | bench CLI | new repair-subcommand contract, run first, exits 2 on the unknown subcommand | routing repair through the Go binary cannot work when no binary exists, so the pre-binary route is forced |
| 2 | `bench repair --prune` removes exactly the cache entries outside the current version and platform, leaves the current one, and reports what it removed | bench CLI | new cache-tree fixture with three stale and one current entry, run first, exits 2 | pruning everything fails the current-entry assertion; pruning nothing fails the removal report |
| edge of 2 | `bench repair --prune` on an empty or absent cache reports a no-op and exits 0; a second identical prune is a no-op | bench CLI | idempotency and absent-cache cases in the same contract | an implementation that errors on an absent cache or re-reports removals fails one of the paired runs |
| edge of 2 | `bench repair extra`, `bench repair --unknown`, and `bench repair --prune extra` each exit 2 as usage errors | bench CLI | argument matrix in the same contract, run first, all four forms are accepted or ignored today | accepting trailing arguments fails all three rows at once |
| edge of 2 | `bench repair --prune` removes a malformed current-platform regular file while preserving a current-platform directory and refusing to follow a symlink | bench CLI | current-platform cache fixture changed from a directory to a regular file, run first, reports no stale entries and preserves the corrupt file | preserving by entry name alone mistakes corrupt owned state for the valid current platform; following a symlink would escape the cache owner |
| 3 | a tarball whose digest matches the shipped pin installs; one that does not is refused with a message naming the mismatch and nothing is installed | bench CLI | fixture registry serving a tampered tarball with a valid `dist.integrity`, run first, installs it because the registry is the authority | a check that still reads `dist.integrity` passes the tampered tarball, which is exactly the failure |
| 3 | an absent, empty, unparseable, or entry-missing pin manifest each fails closed with a message naming the manifest, and no registry digest is consulted | bench CLI | four fixture wrappers, run first, all install successfully because no manifest is read | any fallback to registry metadata turns at least one of the four rows green when it must be red |
| 4 | the packed wrapper tarball contains a manifest whose pins equal the integrity of every packed platform tarball, for every platform in the matrix | artifact contract | new artifact assertion over `dist/artifacts`, run first, the manifest path is absent from the wrapper | a manifest covering one platform, or stale pins from a prior build, fails the per-platform equality |
| 4 | a build that would pack the wrapper without a present, non-empty manifest fails instead of publishing | artifact contract | build driven with the manifest step disabled, run first, the build succeeds and ships a wrapper with no pins | a warning-only guard lets the build exit zero |
| 4 | the independent reproducibility build produces a byte-identical manifest | artifact contract | existing reproducibility comparison extended to the manifest path; starts red because the path does not exist | nondeterministic ordering or embedded timestamps differ between the two builds |
| 5 | a metadata or tarball endpoint that never responds fails at the 60-second deadline with a message naming the deadline, not a hang | bench CLI | hung fixture endpoint under an injected short deadline, run first, the test times out because `fetch` is unbounded | an unbounded fetch produces a harness timeout rather than the named failure |
| 5 | a download of exactly 100 MB is accepted and 100 MB + 1 is refused as oversized | bench CLI | paired boundary fixtures, run first, both are accepted | an off-by-one or absent cap fails one of the paired cases |
| 5 | a payload decompressing to exactly 200 MB is accepted and one exceeding it is refused before full materialization | bench CLI | paired fixtures plus a decompression-bomb fixture with peak memory observed, run first, the bomb is fully expanded | a post-hoc size check passes the incremental assertion only by materializing the payload |
| 6 | a repair that fails after another process has installed the target leaves the installed target intact and removes only its own temp file | bench CLI | two-process fixture using the existing repair-ready synchronization hook, run first, the loser deletes the winner's binary | a cleanup that removes the target fails the post-run existence assertion |
| 6 | SIGINT mid-repair removes the temp file, leaves any installed target intact, and exits 130 | bench CLI | interrupt-stage fixture extended with a pre-installed target, run first, the target is removed | reusing the outer handler for interrupt reintroduces the deletion |
| edge of 6 | a repair interrupted before any temp file exists exits 130 and removes nothing | bench CLI | earliest interrupt stage in the same fixture | an unconditional unlink of an unset path errors instead of exiting cleanly |
| 7 | `bench version` prints exactly `bench <version> (<os>/<arch>)` | cmd/bench + bench CLI | existing `cmd/bench/main_test.go` assertion inverted to the new string, run first, fails on the current `benchkit` prefix | the exact-format assertion rejects a partial rename that keeps the old prefix or drops the platform |
| 7 | no user-facing output string contains `benchkit`, while the `projects/benchkit.md` profile path and this repository's internal references remain permitted | conformance | extended stray-name sweep, run first, fails on the version line, the wrapper help text, and the platform package description | a sweep that enumerates the three current sites and the permitted exceptions cannot be passed by renaming only one |
| 8 | the root `package.json` carries `repository`, `homepage`, `bugs`, and `author` with the locked values | conformance | new assertion in `package_shipped_surface_test.go`, run first, all four fields are absent | naming three of the four fails the enumerated check |
| 8 | every generated platform package carries the same four values as the root, for every platform in the matrix | conformance + artifact contract | generator output inspected per matrix entry, run first, the platform literal has none of the four | hardcoding the values in the platform literal passes equality but fails the single-source conformance check below |
| 8 | the platform package fields are derived from the root package rather than restated | conformance | mutation check that changes the root value and asserts the generated platform value follows, run first, fails because the fields do not exist | a second hardcoded copy does not follow the root and fails the mutation |
| 9 | the offline proof writes a `ft87/offline-network-control/v1` envelope whose key, owner, schema, identity, and payload digest all validate | releaseevidence + workflow record | requirement inspection over the produced record, run first, reports `public.ft87.offline_network_control` missing | an envelope with a stale digest or wrong identity fails the existing producer validation |
| 9 | the payload names all five suppressed operation classes with the sentinel's observed attempt count for each | releaseevidence | fixture-record assertion enumerating the five classes, run first, no record exists | a record naming only the HTTP providers fails the enumerated class set |
| edge of 9 | a sentinel that observed a nonzero attempt produces a record whose status is not a pass, rather than no record or a passing one | releaseevidence | fixture with an injected attempt, run first, no record exists either way | a producer that only writes on success makes a failed proof indistinguishable from an unrun one |
| 10 | `bench repair` appears in the CLI inventory and no stale command reference remains | conformance | existing stale-command-reference sweep plus the inventory check, run first, the inventory omits the new subcommand | the sweep already fails on a documented command the wrapper does not route and the reverse |

Degenerate implementations checked: a message-only change that still repairs
fails story 1's child-process marker; a pin check that keeps a registry fallback
fails story 3's four fail-closed rows; a manifest covering one platform fails
story 4's per-platform equality; post-hoc size checks fail story 5's incremental
decompression row; a cleanup handler shared between error and interrupt fails
story 6's interrupt row; renaming only the version line fails story 7's sweep;
hardcoding the metadata into the platform literal fails story 8's mutation
check; and a producer that writes only on success fails story 9's edge row.

### Edge inventory

Walked per behavior against the profile's hostile-input checklist and the
canonical classes.

- Error path — covered: stories 1, 3, 5, 6 rows.
- Empty vs absent input — covered: the empty and absent pin manifest rows
  (story 3) and the absent cache prune row (story 2).
- Boundary values — covered: the 100 MB and 200 MB paired rows (story 5).
- Malformed input — covered: the unparseable manifest row (story 3) and the
  tampered tarball row (story 3).
- Interrupted or partial state — covered: the SIGINT rows (story 6).
- Re-run idempotency — covered: the second-prune row (story 2) and the
  post-repair command row (story 2).
- Hostile environment — covered: the two-process race row (story 6) and the
  decompression-bomb row (story 5).
- Required tool missing from PATH — covered: `node` absent keeps its existing
  named failure, now reached through the explicit `bench repair` route as well
  as the `BENCH_REPAIR=1` route; asserted in story 2's contract.
- Unquoted multi-word arguments — covered: story 2's argument matrix drives the
  wrapper with quoted multi-word and leading-dash arguments.
- Invocation through every shipped surface — covered: story 2 drives `bench
  repair` through both the real kit CLI and the installed shim in the artifact
  contract.

**Won't handle**

- Control bytes in registry-supplied names — the pin manifest is matched by
  exact package name and version against a build-generated table, so a hostile
  name matches nothing and fails closed before any string is rendered.
- Paths with spaces or glob characters in the cache root — the cache root is
  derived from `BENCH_HOME`, which slice 1 already resolves and quotes; prune
  operates on entries it enumerates rather than on globs.
- Symlinked cache entries — prune refuses to follow symlinks by construction,
  so there is no traversal behavior left to test separately.
- Registry compromise that serves both a bad tarball and a matching pin — the
  pin ships inside the wrapper the user already installed, so this requires
  compromising the wrapper itself, which is outside the repository-controlled
  scope the map fixed.
- A file lacking a trailing newline — no behavior in this slice parses
  line-oriented input.
- cwd deeper than the repository root — repair resolves from the kit directory,
  never from cwd; slice 3 owns root-anchored slug resolution.

## Out of scope

- **A shared generator for the pin manifest and the FT83 release evidence index**
  — the map's "Not yet specified" names it as a single-source candidate once both
  exist, and both exist only at the end of this slice; it is a distinct
  consolidation feature with its own correctness question about what the shared
  schema is. Estimate: 4 edits, 2 gate runs.
- **Repair resuming a partial download** — a separate capability with its own
  range-request and partial-state correctness model, unrelated to bounding or
  pinning. Estimate: 5 edits, 3 gate runs.
