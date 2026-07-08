# npm identity and install path (FT42)

The advertised npm identity was never ours: unscoped `benchkit` is an unrelated
third-party package (v1.4.3), all four `@benchkit/<os>-<arch>` platform packages
are 404, and nothing was ever published under any of the coded names — so
`npx benchkit link` installs the wrong software and the missing-binary error
points at a 404 package. A complete tag-driven release pipeline
(`.github/workflows/release.yml`) already exists and is blocked only on
identity. Registry facts verified 2026-07-08.

## #1: Where does Bench's npm identity live?

Blocked by: —
Type: Grill

### Question
Publish under an owned scope, rename, or de-advertise npm and lead with the
clone path? Blocks all external adoption and every follow-through edit.

### Answer
Publish under the reviewer-owned npm org **`redbench`** (created 2026-07-08).
The `benchkit` org name is taken (claim attempted and refused 2026-07-08), so
the coded `@benchkit` scope is unusable, not merely unpublished. De-advertising
is rejected: the release pipeline works and npx adoption is a goal.

## #2: What is the wrapper package's name and bin set?

Blocked by: #1
Type: Grill

### Answer
Unscoped **`redbench`** (verified free 2026-07-08; claimed at first publish).
Bins become `{bench, redbench}`; the dead `benchkit` bin alias is dropped. The
`redbench` bin entry is what makes `npx redbench <cmd>` resolve (npx matches
bin name to package name when a package ships multiple bins).

## #3: What are the platform packages named?

Blocked by: #1
Type: Grill

### Answer
`@redbench/<os>-<arch>` — the coded shape kept, scope swapped. Same matrix
(darwin/linux × arm64/x64), still os/cpu-pinned optionalDependencies of the
wrapper at the identical version, still generated from the single
`scripts/platforms.json` source.

## #4: How far does the `redbench` identity reach?

Blocked by: #1
Type: Grill

### Answer
npm distribution identity only. The product, docs, repo, CLI command (`bench`),
and the project profile name stay Bench/benchkit-the-project. Only npm
package/org strings change, including the profile's "shipped as the npm
package" sentence. CHANGELOG history entries are history and are not rewritten.

## #5: What does the missing-binary error (exit 127) tell the user?

Blocked by: —
Type: Grill

### Answer
Both remedies, always, no install-mode detection in the wrapper: one line for
the npm install (`npm install @redbench/<os>-<arch>`) and one for the clone
build (`scripts/go-build.sh`). Off-matrix platform stays a distinct exit 2.

## #6: How do the sweep, the README fix, and the first publish sequence?

Blocked by: #1
Type: Grill

### Answer
One sweep spec: identity strings (package.json, gen-platform-packages/
release.yml, wrapper `platform_pkg()`), README install path (npx/global/
uninstall strings, clone fallback gains the exact `scripts/go-build.sh` build
command), the 127 error, and the contract/conformance test expectations.
After merge, reviewer op actions outside the spec: wire `NPM_TOKEN` for the
`redbench` org into the GitHub repo, run release.yml's documented pre-tag
smokes, tag. First publish ships the version in package.json at tag time
(0.2.0); the README being briefly ahead of the registry between merge and tag
is accepted. The spec must not depend on live registry state.

## Handoff

1. **Module boundaries.** Identity constants: `package.json` (name, bin,
   optionalDependencies), `scripts/gen-platform-packages.sh` +
   `scripts/platforms.json` (platform matrix), `bin/bench.sh` `platform_pkg()`
   (the one shell-side scope literal), `.github/workflows/release.yml` (publish
   loop path). Install docs: README worker/maintainer section;
   `projects/benchkit.md` intro sentence. Error surface: `route_binary`'s 127
   branch. Tests: `internal/contract/surface/*` and the conformance npm-shape /
   release-workflow checks. Outside: product branding, repo name, `bench`
   command name, CHANGELOG history.
2. **Contracts.** Wrapper: npm name `redbench`, bins `{bench, redbench}`.
   Platform packages: `@redbench/<os>-<arch>`, os/cpu-pinned, version-locked to
   the wrapper. Missing binary: exit 127, stderr names both remedies (npm
   install of the exact platform package; the clone build command). Off-matrix:
   exit 2, unchanged. README primary path `npx redbench link` /
   `npm i -g redbench`; uninstall `npm uninstall -g redbench`; clone fallback
   states the build command before the symlink step. release.yml stays
   tag-driven with provenance and `--access public`.
3. **Deep vs thin.** `gen-platform-packages.sh` + `platforms.json` is the deep
   unit — names and metadata derive from the one matrix, never enumerated
   twice. `platform_pkg()` keeps the scope literal single on the shell side.
   README and profile prose are thin consumers of those facts.
4. **Black-box assertables.** `npm pack --dry-run` name and file list;
   generated platform package.json name/os/cpu/version equality against the
   wrapper; 127 stderr contains both remedy strings against a fixture kit with
   the binary absent; exit 2 on an off-matrix host; a shipped-files grep for
   `@benchkit`/`npx benchkit`/`npm i -g benchkit` returning empty outside
   CHANGELOG.
5. **Gate attachment.** The conformance phase already checks npm dry-run
   package shape and release-workflow structure, and the surface contract tests
   exercise routing/repair against fixture registries — stale `@benchkit`
   expectations go red, which is the rename's bite. Gate-blind: the live
   registry publish itself; owned by release.yml's documented pre-tag manual
   smokes (reviewer-run).
6. **Hostile-input owners.** Required-tool-missing and symlink invocation:
   existing wrapper/routing tests, semantics unchanged, strings renamed.
   Every-shipped-surface invocation: owns sweep completeness — real kit CLI,
   linked-repo by-path CLI, and hooks must all print the new identity.
   Absent vs present-but-empty binary: already distinct in the wrapper; keep
   asserted through the rename. Paths with spaces: the printed clone-build
   remedy must be quote-safe. Remaining checklist classes n/a — no new parsing
   or git-sourced text in this diff.
7. **Uncertainty flags.** None — every fork reviewer-resolved. One op
   dependency (NPM_TOKEN + tag) sits after merge and outside the spec.
8. **Rejected alternatives.** Claim `@benchkit` org (taken, attempt
   2026-07-08); unscoped rename `bench-kit` (typo-adjacent to the squatter,
   still needs an org claim); scoped wrapper `@redbench/bench` (leaves free
   unscoped `redbench` claimable, clunkier npx string); de-advertise npm (idles
   a working pipeline, drops npx adoption); context-aware 127 remedy (branching
   in a deliberately minimal POSIX wrapper); full Redbench rebrand (doc-wide
   churn, no distribution gain).
9. **Domain watch-outs.** Unscoped `benchkit` on npm is someone else's package —
   any surviving `npm i benchkit` string installs the wrong software. Scoped
   publishes default to private — `--access public` in release.yml is
   load-bearing. npx resolves multiple bins by matching bin name to package
   name — dropping the `redbench` bin entry silently breaks the advertised
   one-liner. Registry names are secured only by publishing: until the first
   tag, `redbench` (unscoped) remains claimable by anyone.

Dependency order: n/a — single spec; the reviewer op actions (NPM_TOKEN,
pre-tag smokes, tag) follow the merge and live outside it.
