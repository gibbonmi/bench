# npm identity and install path (redbench rename + npx-from-git)

Status: implemented

## Problem

The advertised npm identity was never ours. Unscoped `benchkit` is an unrelated
third-party package (`benchkit@1.4.3`), all four `@benchkit/<os>-<arch>` platform
packages 404, and nothing was ever published under any coded name. So the README's
primary install path (`npx benchkit link`) installs the wrong software, and the
missing-binary error (exit 127) tells the user to `npm install @benchkit/<...>` —
a 404. A complete tag-driven release pipeline already exists and is blocked only on
identity.

## Solution

One sweep, from the reviewer decisions in `decisions/npm-identity.md`:

1. **Rename the npm identity to `redbench`** (reviewer-owned org, created
   2026-07-08): wrapper package `redbench`, bins `{bench, redbench}`, platform
   packages `@redbench/<os>-<arch>`. Only npm distribution strings change — the
   product, repo, `bench` command, and profile name stay Bench/benchkit-the-project.
2. **Publishing is deferred** (reviewer, 2026-07-08). Until first publish the
   advertised install path is **npx-from-git**: a two-line `package.json` change
   (`prepare` runs `scripts/go-build.sh` → `dist/bench`; `dist/` joins `files[]`)
   makes `npx github:gibbonmi/bench#<ref> link` build the core on the consumer
   machine at install time. Proven end-to-end 2026-07-08 (decision #7).
3. **The interim 127 error prints only the build remedy** (the npm line would name
   404 packages); the npm line rejoins publish-day.

The spec must not depend on live registry state. Publish day is a small, deferred,
reviewer-initiated follow-up (out of scope below).

## User stories

1. As a maintainer publishing packages, I want the wrapper package named `redbench`
   with bins `{bench, redbench}` (the dead `benchkit` bin alias dropped), so that
   `npx redbench <cmd>` will resolve once published and no string names a package we
   do not own. Line: claude-sonnet-5 / medium. The names are exact literals the
   package-surface gate matches, so the work is fully observable at a known seam.

2. As a maintainer, I want the four platform packages named `@redbench/<os>-<arch>`,
   generated from the single `scripts/platforms.json` matrix at the wrapper version,
   so that the os/cpu-pinned optionalDependencies stay derived from one source.
   Line: claude-sonnet-5 / medium. The generator output is asserted field-by-field
   against the matrix by an existing contract test.

3. As a maintainer, I want the wrapper `optionalDependencies` to list the four
   `@redbench/*` packages at the wrapper version, so that the dependency scope
   matches the generated packages. Line: claude-sonnet-5 / medium. The generator
   test already asserts `optionalDependencies` deep-equals the matrix mapping.

4. As a consumer installing from git before publish, I want `package.json` to carry
   `"prepare": "bash scripts/go-build.sh \"$PWD\" dist/bench"` and `dist/` in
   `files[]`, so that `npx github:gibbonmi/bench#<ref> link` builds the core at
   install time and the pack carries the built binary. Line: claude-opus-4-8 /
   medium. The prepare/files posture is a behavioral install-time change whose
   effect (a git install runs green) needs an end-to-end probe, not just a string
   match.

5. As a maintainer, I want the npm-pack dry-run shape check to treat `dist/bench` as
   tree-state-independent — required nowhere, forbidden nowhere — while still
   requiring the rest of `files[]`, so that the check stays deterministic whether it
   runs on a built tree (binary present) or the CI publish checkout (binary absent).
   Line: claude-opus-4-8 / medium. This is a judgment call about the check's posture,
   not a mechanical rename, and it is the one place the rename could silently break
   determinism.

6. As a consumer on an unsupported build, I want the missing-binary error (exit 127)
   to print only the clone/git build remedy (`scripts/go-build.sh`) and no npm line
   until publish, with off-matrix staying exit 2, so that the error never points at a
   404 package. Line: claude-opus-4-8 / medium. This changes runtime behavior and
   asserts a negative (no npm line), so it carries more than a string swap.

7. As a maintainer, I want `.github/workflows/release.yml` to build and publish
   `@redbench/*` packages and the `redbench` wrapper, staying tag-driven with
   provenance and `--access public`, so that publish day needs only a token and a
   tag. Line: claude-sonnet-5 / medium. The workflow paths are exact-string renames
   and its structure is asserted by the conformance check.

8. As a worker reading the README, I want the primary install path to be
   `npx github:gibbonmi/bench#<ref> link` with prerequisites repo access + Node + Go,
   the clone fallback to state the exact build command before the symlink step, and
   the uninstall strings to name `redbench`, so that every documented path works
   today. Line: claude-opus-4-8 / medium. Install prose is semantic and only
   partially gate-graded (the stale-string sweep), so it needs a careful writer.

9. As the teammate maintaining the profile, I want `projects/benchkit.md`'s "shipped
   as the npm package `benchkit`" sentence to name `redbench`, so that the profile
   states the current npm identity. Line: claude-sonnet-5 / medium. One sentence, a
   direct string swap.

10. As a maintainer guarding against a half-done rename, I want a shipped-surface
    sweep asserting no `@benchkit`, `npx benchkit`, or `npm i -g benchkit` string
    survives outside `CHANGELOG.md`, so that a stale identity string fails the gate
    rather than shipping. Line: claude-opus-4-8 / medium. This is a new gate check —
    authoring the oracle needs care about what it matches and where it looks.

## Implementation decisions

- **Identity constants, one source each.** `package.json` (name, bin,
  optionalDependencies scope, `prepare`, `files[]`); `scripts/platforms.json` +
  `scripts/gen-platform-packages.sh` (platform matrix — names/os/cpu/version all
  derive here, never enumerated twice); `bin/bench.sh` `platform_pkg()` (the single
  shell-side scope literal `@redbench/`); `.github/workflows/release.yml` (publish
  loop path `@redbench`). README worker/maintainer section and `projects/benchkit.md`
  intro are thin prose consumers of those facts.

- **Product identity stays.** `bench version` output (`benchkit X.Y.Z (...)`), the
  `bench` command, repo, docs branding, and CHANGELOG history are untouched. Only npm
  package/org strings change.

- **npx-from-git enablement.** `package.json` gains
  `"prepare": "bash scripts/go-build.sh \"$PWD\" dist/bench"` and `dist/` in `files[]`.
  No wrapper change: `bench_binary_path`'s search order already finds `dist/bench`
  first. `scripts/go-build.sh` stays repo-only (not in `files[]`); the git install has
  the full clone at prepare time, so the build script is present regardless. `postinstall`
  is already git-install-safe (advice line, exit 0).

- **Interim 127 remedy.** `route_binary`'s 127 branch prints one line naming the build
  remedy (`scripts/go-build.sh`), no npm line, no install-mode detection. The npm
  remedy line is deferred to publish day. Off-matrix stays exit 2, unchanged.

- **dist/bench pack posture.** Move `dist/bench` out of `ForbiddenPackAssets`; do
  **not** add it to `RequiredPackAssets`. It becomes tree-state-independent so the
  dry-run check passes on both a built tree and the CI publish checkout. The
  tag-driven CI workflow (checkout never has `dist/bench`) stays the only registry
  publish path, so a registry publish never fattens the wrapper with one platform's
  binary.

- **Release workflow.** Rename `@benchkit` → `@redbench` in the cross-compile output
  path and the publish loop. Stays tag-driven, provenance, `--access public`.

## Testing decisions

- **What a good test is here:** exercise the shipped surface (generated package
  metadata, npm-pack file list, launcher routing/error text, a real git-protocol
  install), not internal helpers. Every assertion is a black box a consumer could
  observe.
- **Prior art:** `internal/contract/surface/package_test.go` (generator + npm-pack
  surface), `internal/contract/surface/binary_repair_test.go` (127 stderr against a
  fixture registry, `WithSpacePath` quote-safety), `internal/contract/surface/go_routing_test.go`
  (routing + missing-binary error), `internal/conformance/package_core_checks_test.go`
  (files[] existence, dry-run shape, release-workflow structure). Heavy contract tests
  that spin real npm/registry already exist, so the new git-install probe has precedent.
- **Gate command:** the project gate (`bench gate`).

### Seam diagram

**Seam A — platform-package generator + npm-pack surface** (stories 1–5, 7)

    trigger: gen-platform-packages.sh <out>  |  npm pack --dry-run --json
        │
        ▼
    platforms.json ──▶ [ gen-platform-packages.sh ] ──▶ @redbench/<os>-<arch>/package.json
    package.json   ──▶ [ (name/os/cpu/version)    ]      (name, os, cpu, version)
    package.json   ──▶ [ npm pack --dry-run       ] ──▶ file list (files[] present; dist/bench tree-dependent)
                          ◀ tests attach here: package_test.go asserts field equality
                            + RequiredPackAssets present, dist/bench neither required nor forbidden

**Seam B — launcher routing + interim missing-binary error** (story 6)

    trigger: bench <cmd>  (fixture kit, binary absent / off-matrix host)
        │
        ▼
    host uname ──▶ [ platform_pkg / bench_binary_path ] ──▶ exec dist/bench   (present)
    fixture    ──▶ [ route_binary 127 branch          ] ──▶ exit 127 + build remedy, NO npm line  (absent)
    off-matrix ──▶ [                                   ] ──▶ exit 2 "unsupported platform"
                      ◀ tests attach here: binary_repair_test.go / go_routing_test.go
                        assert stderr contains scripts/go-build.sh, contains no "npm install"

**Seam C — git-protocol install probe** (story 4)

    trigger: npm install git+file://<repo> into a throwaway prefix
        │
        ▼
    full clone ──▶ [ prepare: go-build.sh → dist/bench ] ──▶ dist/bench packed
               ──▶ [ pack (files[] incl. dist/)        ] ──▶ installed tree
                      ◀ tests attach here: run `bench version`, require exit 0 + version line

**Seam D — shipped-surface stale-identity sweep** (story 10)

    trigger: conformance check over shipped markdown/config
        │
        ▼
    tracked files ──▶ [ grep @benchkit / npx benchkit / npm i -g benchkit ] ──▶ diag per hit
    (excl CHANGELOG)                                                            (empty = green)
                      ◀ tests attach here: conformance asserts no diags outside CHANGELOG.md

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | wrapper name `redbench`, bins `{bench, redbench}`, no `benchkit` bin | A | conformance/pack surface check reads `package.json`; a leftover `benchkit` name or bin fails the identity assertion | pack/identity assertion reads the exact name+bin set; a missed rename mismatches |
| 2 | platform packages `@redbench/<os>-<arch>` from matrix | A | `testPackageGeneratorOutput` asserts `got.Name == "@redbench/"+os+"-"+arch`; old scope → `name is @benchkit/...` fail | test derives expected name from the matrix + new scope; stale generator scope diverges |
| 3 | optionalDependencies = four `@redbench/*` at wrapper version | A | `testPackageGeneratorOutput` deep-equals `wrapper.OptionalDependencies` to matrix `@redbench` map; `@benchkit` entry → mismatch fail | the deep-equal names both scope and version; a stale or partial rename fails equality |
| 4 | git install builds via prepare and runs green | C | new probe: `npm install git+file://…` then `bench version` — pre-change (no prepare/dist) the installed tree has no binary → non-zero exit | the probe drives the real install path; absent prepare or `dist/` in files[] yields no runnable binary |
| 5 | dry-run shape deterministic across built/unbuilt tree | A | run `npm pack --dry-run` on a built tree (dist/bench present) and simulate CI checkout (absent); `dist/bench` in Forbidden → built tree fails; in Required → checkout fails | asserting `dist/bench` in neither list is the only posture green in both states; either list reintroduces tree-dependence |
| 6 | 127 prints build remedy, no npm line; off-matrix exit 2 | B | `go_routing_test`/`binary_repair` require stderr contains `scripts/go-build.sh` and `RequireContains(... "npm install")` is inverted to absence; current `npm install @benchkit/` stderr fails the new assertion | the assertion pins the exact remedy text and the negative; leaving the npm line in fails the no-npm-line check |
| 7 | release workflow builds/publishes `@redbench`, tag-driven + provenance + access public | A | `checkReleaseWorkflow` structure check + path grep; `@benchkit` in the cross-compile/publish path → identity sweep (story 10) fails; missing `provenance`/`npm publish` → structure fail | the workflow-structure check asserts publish shape; the identity sweep catches a stale scope path |
| 8 | README primary path npx-from-git, prereqs, clone build command, uninstall `redbench` | A/D | story-10 sweep: `npx benchkit`/`npm i -g benchkit` surviving in README → diag; not TDD-able for prose accuracy beyond the sweep | the sweep catches the stale advertised strings; prose correctness of the new path is reviewer-checked |
| 9 | profile names npm package `redbench` | D | story-10 sweep flags `benchkit` npm-package claim in `projects/benchkit.md` if unrenamed | the sweep matches the npm-identity string in the profile intro |
| 10 | no `@benchkit`/`npx benchkit`/`npm i -g benchkit` outside CHANGELOG | D | new conformance check returns a diag for any surviving hit; seed a stray string → red | the check greps the shipped surface directly; any missed rename anywhere becomes one diag |

### Edge inventory

Walked per behavior; each lands as a coverage row above or a **Won't handle** line here.

- **Absent vs present-but-empty vs non-executable binary** — covered: `go_routing_test`
  already exercises all three against `dist/bench`; semantics unchanged, only the 127
  remedy text changes (row 6).
- **Off-matrix platform** — covered: exit 2 unchanged, asserted in seam B (row 6).
- **Paths with spaces** in the printed build remedy — covered: `binary_repair` fixtures
  use `WithSpacePath`; the printed `scripts/go-build.sh` remedy must stay quote-safe
  (row 6).
- **Every shipped surface renamed** (real kit CLI, linked-repo by-path CLI, hooks) —
  covered by the story-10 sweep over shipped files (row 10).
- **Built tree vs CI publish checkout** for dry-run determinism — covered (row 5).
- **git install without Go / without repo access on the consumer** — **Won't handle**:
  the prepare build fails the install honestly (non-zero); documented as a prerequisite
  (Node + Go + repo access), not a case the kit masks.
- **CHANGELOG history entries** naming `benchkit` — **Won't handle**: history is not
  rewritten (decision #4); the sweep excludes `CHANGELOG.md` by design.
- **Live registry state** (name still claimable, packages unpublished) — **Won't
  handle** in this spec: gate-blind, owned by the deferred publish-day follow-up and
  release.yml's manual pre-tag smokes.
- **git install with the unpublished `@redbench/*` optionalDependencies present**
  (npm skipping the 404s and still exiting 0) — **Won't handle** in the gate: the
  probe sets `npm_config_omit=optional` so it stays deterministic and offline. The
  optional-dep tolerance the advertised command relies on was proven end-to-end
  (decision #7) and is re-checked by release.yml's pre-tag smokes; a fixture-registry
  variant that 404s the four scopes could gate it later if the risk warrants.
- **npx git-install cache staleness** — **Won't handle** in tests: mitigated by
  advertising a pinned `#<ref>`; documented in README, not gate-enforced.

## Out of scope

- **Publish day.** Wire `NPM_TOKEN` for the `redbench` org, run release.yml's pre-tag
  manual smokes, tag (ships the package.json version), restore the npm remedy line in
  the 127 error, re-advertise `npx redbench link` / `npm i -g redbench` in the README.
  Deferred reviewer-initiated work — a *separate capability* (live publication) gated
  on a reviewer decision and a token, not the rest of this rename. Estimate: ~4 edits
  (token secret, 127 npm line, README re-advertise, tag), 1 gate run + reviewer smokes.
- **Live registry integration test.** A test that hits the real npm registry to prove
  the published shape. Separate capability (network-coupled, non-deterministic); the
  fixture registry mirrors the documented metadata shape. Estimate: 1 edit, out-of-band
  run — belongs to the pre-tag smokes, not the gate.
