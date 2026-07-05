# binary-auto-repair

Status: staged

## Problem

When the `@benchkit/<os>-<arch>` optional dependency is missing — skipped by the
installer, pruned by a lockfile, or never fetched under `--ignore-scripts` — every
routed `bench` command exits 127 with "install @benchkit/…". On 2026-07-04 this
degraded the git guard hook into its fail-closed rim and blocked every git command
until a manual source build. npm consumers have no source tree to build: for them
the 127 is a dead end, and the remedy line assumes npm knowledge the harness
mid-session doesn't have.

## Solution

When the launcher cannot resolve a binary for a porcelain command, it repairs
itself: fetch the matching `@benchkit` platform package from the npm registry,
verify the tarball against the registry's SRI integrity hash, extract the one
binary, install it atomically into a version-keyed cache under `$BENCH_HOME`, and
exec it — announcing every mutation. Any failure (offline, hash mismatch, no
node) degrades to today's advice message; hook plumbing subcommands never touch
the network, so the guard rims stay fast and deterministic.

## User stories

1. As a global-install user whose optional platform package never landed, I want
   any porcelain `bench` command to fetch, verify, cache, and exec the binary
   automatically, so that the CLI works instead of dead-ending at 127.
   Line: claude-opus-4-8 / medium. The repair unit mixes network, integrity
   verification, and atomic installation, whose failure postures are subtle even
   though the fixture-registry tests observe them.
2. As a user whose download is corrupt or tampered with, I want a tarball that
   fails the registry integrity hash refused with a message naming the mismatch
   and nothing installed, so that an unverified binary never runs.
   Line: claude-opus-4-8 / medium. Verification is the security-relevant half of
   the same unit and must fail toward refusal, not toward install.
3. As an offline or air-gapped user, I want repair failure to fall back to the
   existing manual-install advice with no partial files left in the cache, so
   that the failure is clean and the remedy is still on screen.
   Line: claude-opus-4-8 / medium. Same unit; the cleanup-on-failure path is
   where torn state would otherwise accumulate.
4. As an agent relying on the guard rims, I want the hook-and-adapter plumbing
   subcommands (`tree-hash` through `worktree-lease-file` in the dispatch) to
   never attempt a network fetch, so that guards keep their current fast
   fail-closed/fail-open behavior when the core is missing.
   Line: claude-sonnet-5 / low. This is a dispatch-group membership check at the
   existing router seam, fully gate-observable.
5. As the reviewer enforcing the announce-mutations posture, I want every
   directory created and file written by repair announced to stderr with its
   path, so that a self-mutating launcher is never silent.
   Line: claude-sonnet-5 / low. Mechanical stderr lines beside each mutation,
   asserted by substring in the contract tests.
6. As a repeat user, I want the second invocation to resolve the cached binary
   with zero registry requests, so that repair is a one-time cost.
   Line: claude-sonnet-5 / low. Falls out of the resolver-candidate design and
   is asserted by the fixture registry's hit counter.
7. As an upgrading user, I want the cache keyed by the wrapper's package
   version, so that a new benchkit release fetches its matching binary instead
   of reusing a stale one.
   Line: claude-sonnet-5 / low. The version is read from the kit's package.json,
   which already exists as the single version source.
8. As a user without node on PATH, I want repair skipped with a line saying why
   and the manual advice preserved, so that the no-node failure explains itself.
   Line: claude-sonnet-5 / low. A `command -v` guard in front of the repair
   call with a distinct message.
9. As a user with a torn cache entry (present-but-empty or non-executable
   file), I want repair to treat it as missing and replace it atomically, so
   that an interrupted install self-heals on the next run.
   Line: claude-opus-4-8 / medium. Reuses the resolver's existing
   empty-is-missing rule but the atomic replace is the same careful unit as
   story 1.

## Implementation decisions

- **The repair unit is one Node script shipped under `bin/`** (already inside
  `files[]`). Node ≥22 is the npm distribution's own engine premise, and its
  built-in `https`/`crypto`/`zlib` cover fetch, SRI sha512 verification, and
  gunzip with zero new dependencies. The tar read is a minimal single-entry walk
  (esbuild `install.js` style) extracting only `package/bin/bench`.
- **Trigger point is `route_binary`'s 127 arm, porcelain commands only.** On
  repair success it re-resolves once and execs; on any failure it prints one
  advice line and falls through to today's 127 message. The plumbing dispatch
  group is excluded wholesale — the guards' existing rims stay the contract
  there. Exit codes 0/2/127 keep their current meanings on every path.
- **Cache lands at `$BENCH_HOME/cache/bin/<version>/<os>-<arch>/bench`** and is
  appended as the *last* candidate in `bench_binary_path` — the one resolver
  stays the single source of resolution order, shared with the status adapters,
  and any real package on disk still wins.
- **Registry resolution:** `$BENCH_NPM_REGISTRY` → `npm_config_registry` →
  `https://registry.npmjs.org`. The env override doubles as the test seam.
- **Version and platform derive from the existing single sources** — the kit's
  `package.json` and `platform_pkg` — never a second mapping.
- **Mutation posture inherits the shim-autoinstall decisions:** announce every
  mutation with its path, write via temp file + `mv` (SIGINT-safe, and a
  concurrent race resolves to last-writer-wins of identical verified bytes, so
  no lock), refuse to install unverified bytes, and never make the caller fail
  harder than it does today.
- **Trust model stated honestly:** the integrity hash comes from the same
  registry as the tarball, so verification defends against corruption and
  transport tampering, not registry compromise — that threat is owned by npm
  provenance/2FA upstream, not this launcher.

## Testing decisions

- A good test drives the **launcher surface end to end**: invoke `bin/bench.sh`
  against a fixture kit with no `dist/` and no `node_modules`, point
  `$BENCH_NPM_REGISTRY` at a local fixture registry, and assert exit code, cache
  contents, stderr announcements, and the registry hit count. No test reaches
  into the repair script's internals.
- Tests live in `internal/contract` (runtime-contracts family; prior art:
  `go_routing_test.go`, `runtime_git_test.go`). The fixture registry is a Go
  `httptest` server serving version metadata plus a tarball built at test time
  from a stub binary, mirroring the documented npm registry shape
  (`versions[v].dist.{tarball,integrity}`).
- **Gate:** `.bench/gate.sh` (project default).
- **Gate-blind spot, declared:** the fixture mirrors the registry's documented
  shape; a live-registry shape drift would pass the fake. One manual smoke
  against the real registry before the next release, same as the shim map's
  npm-lifecycle blind spot.

### Seam diagram

    trigger: user or agent runs `bench <porcelain>`; no binary on disk
        │
        ▼
    argv, kit package.json ──▶ [ bin/bench.sh: route_binary → 127 → repair script ]
    registry metadata + tgz ──▶ [   fetch → verify SRI sha512 → extract → mv     ]
                                        │
                                        ▼
                    $BENCH_HOME/cache/bin/<ver>/<os>-<arch>/bench → exec → output
            ◀ tests attach here: contract test runs bench.sh with a fixture kit
              and BENCH_NPM_REGISTRY=<httptest URL>; asserts exit code, cache
              file, stderr announcements, and registry hit count

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | missing binary + porcelain command → binary fetched, verified, cached, command runs | launcher | new `TestRepairDownloadsAndRuns` — starts red: wrapper exits 127 today | passes only if the whole repair chain works |
| 2 | integrity mismatch → refusal message, nothing installed, exit 127 | launcher | new `TestRepairRefusesBadHash` — starts red: no verification exists | a fallback that installs unverified bytes fails the no-file assertion |
| 3 | unreachable registry → advice + existing 127 message, cache dir left clean | launcher | new `TestRepairOffline` — starts red: advice line doesn't exist yet | catches both a hard error replacing the advice and torn partial files |
| 4 | plumbing subcommand + missing binary → zero registry hits, current exit behavior | launcher | regression guard, cannot start red — no fetch exists anywhere today; bites once repair lands | hit-counter assertion fails if repair ever leaks into the plumbing group |
| 5 | every mutation announced with its path | launcher | new `TestRepairAnnounces` — starts red | silent mutation fails the substring assertions |
| 6 | second invocation → zero registry hits, same binary | launcher | new `TestRepairIdempotent` — starts red: no cache exists to hit | refetch-every-time fails the hit counter |
| 7 | bumped wrapper version → new version path fetched, old cache untouched | launcher | new `TestRepairVersionKeyed` — starts red | stale-cache reuse fails the fetched-version assertion |
| 8 | node absent from PATH → distinct skip message + manual advice, exit 127 | launcher | new `TestRepairNoNode` — starts red: today's message lacks the skip line | catches a repair path that crashes or goes silent without node |
| 9 (edge of 1) | present-but-empty or non-executable cache file → replaced atomically, command runs | launcher | new `TestRepairReplacesTornCache` — starts red | torn state that wedges resolution fails the recovery assertion |

### Edge inventory

Walked per the benchkit hostile-input checklist; every class lands in a row
above or a line below.

- Paths with spaces/globs → `$BENCH_HOME` quoting, exercised by rows 1/6 (fixture HOME contains a space).
- Absent vs present-but-empty binary → row 9.
- Required tool missing → row 8 (node); offline network → row 3.
- Malformed input (bad tarball/hash) → row 2.
- Interrupted/partial state → atomic temp+mv decision, observed via row 9.
- Re-run idempotency → row 6.
- Hostile environment (hook context) → row 4.
- Invocation through every shipped surface → linked-repo by-path CLI resolves the same wrapper and the same cache candidate; covered by row 1's fixture-kit invocation, which is the by-path form.
- **Won't handle:** off-matrix platforms (Windows) — `platform_pkg`'s exit-2 arm is unchanged and already contract-tested.
- **Won't handle:** proxy-only networks — node's fetch ignores `HTTPS_PROXY` by default; the failure degrades to row 3's advice, and manual install remains the path.
- **Won't handle:** registry-compromise threat model — integrity and tarball come from the same origin; npm provenance is the upstream control, and pinning hashes in-repo would break every release.
- **Won't handle:** concurrent repair lock — atomic rename makes the race benign (identical verified bytes, last writer wins).

## Out of scope

- **`bench doctor` binary-presence report row** — a diagnostic capability on the
  doctor surface, distinct from the launcher's self-repair; ~4 edits, 1 gate run.
- **FT2 adversarial gate pinning** — separate threat model, already parked;
  ~6 edits, 2 gate runs.
- **Windows platform support** — a build-matrix capability, not part of this
  fallback; estimate unknown until the matrix decision is made.
