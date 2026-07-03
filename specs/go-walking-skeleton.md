# go-walking-skeleton — slice 1 of the Go rewrite

Status: implemented

Map: `decisions/go-rewrite.md` (closed). This spec builds slice 1 of the
8-slice strangler: the Go module, the version subcommand, the strangler
router, the npm platform packaging, the release workflow, and the gate's Go
layer — the pipeline proven before any real logic moves.

## Problem

The go-rewrite map committed the kit's executable core to a strangler port,
but nothing exists to strangle into: no Go module, no way to ship a binary
to consumers, no router sending a subcommand to Go instead of shell, and no
gate authority over Go code. Every later slice depends on these shapes being
right; discovering a packaging or routing mistake in slice 5 means rework
across the whole program.

## Solution

A walking skeleton: one genuinely new subcommand (`bench version`)
implemented only in Go, reachable through the same `bench` entry consumers
already use, distributed through esbuild-style npm platform packages, built
and published by a tag-triggered release workflow, and graded by a gate that
now compiles, vets, tests, and cross-compiles the module — with
red-by-construction canaries proving the new checks bite.

## User stories

1. As a kit developer, I want a Go module at the repo root (`cmd/bench`,
   toolchain pinned in `go.mod`) whose binary implements a `version`
   subcommand, so that the strangler has a compilable, idiom-setting core to
   grow into. Line: claude-opus-4-8 / medium. This is the first Go code in
   the repo and every later slice copies its idioms, so it warrants the mid
   tier even though the diff is small.

2. As a bench user, I want `bench version` to print exactly one line —
   `benchkit <version> (<os>/<arch>)` — and exit 0, from any cwd including
   outside a git repo, so that humans and hooks can check what is installed.
   Line: claude-sonnet-4-6 / medium. The behavior is a small typed surface
   the gate fully asserts, which is the cheap-tier case.

3. As a bench user, I want the shell CLI to route `version` to the Go
   binary through a fixed resolution order — repo-local dev build, then the
   bundled platform package, then the hoisted global sibling — so that one
   router serves dev checkouts and npm installs alike, and every later
   ported subcommand reuses it. Line: claude-opus-4-8 / medium. This
   touches the CLI dispatch heart and a wrong resolution order breaks every
   future slice, so it is mid-tier work despite full contract coverage.

4. As a bench user whose platform package is missing or empty, I want a
   one-line stderr error naming the exact `@benchkit/<os>-<arch>` package
   and the `npm install` remedy, exiting 127 (spec-level choice: the shell convention for
   command-not-found), so that npm's optional-deps
   failure mode is self-describing instead of a stack of shell noise. Line:
   claude-sonnet-4-6 / medium. The contract is fully asserted in a
   fabricated node_modules sandbox, so the cheap tier suffices.

5. As a bench user on a platform outside the support matrix, I want a
   distinct "unsupported platform: <os>/<arch>" error (exit 2, the usage
   exit code), so that "we don't build for you" is never confused with "your
   install is broken". Line: claude-sonnet-4-6 / medium. Same sandbox, same
   assertion style as story 4.

6. As a kit maintainer, I want the four platform packages generated from
   one template plus a target matrix — never hand-maintained — with the
   wrapper's `optionalDependencies` derived from the same matrix and pinned
   to the wrapper's own version, so that package metadata has one source
   and cannot drift. Line: claude-sonnet-4-6 / medium. Mechanical
   generation validated by a gate dry-run is exactly the cheap-tier case.

7. As a kit maintainer, I want the gate to gain a Go layer — `gofmt -l`
   clean, `go vet`, `go build`, `go test`, and a build of all four
   GOOS/GOARCH targets (spec-level choices: the map named the check family,
   not the tools, and put cross-compile in CI — running it locally too
   catches a broken target before any tag) — that hard-fails when the
   toolchain is absent, so
   that the oracle has authority over the new core from day one. Line:
   claude-opus-4-8 / medium. Gate logic routes mid per the profile because
   a wrong oracle is the worst defect class in this kit.

8. As a kit maintainer, I want red-by-construction canaries for the Go
   layer (a build-breaking fixture and a test-failing fixture), so that the
   new checks provably bite instead of rotting into always-pass. Line:
   claude-sonnet-4-6 / medium. Canary fixtures follow an existing
   mechanical pattern.

9. As a kit maintainer, I want a tag-triggered GitHub Actions release
   workflow that cross-compiles the four targets, stamps the version from
   `package.json`, generates the platform packages via story 6's generator,
   and publishes all five packages to npm with provenance, so that a
   release is one `git tag && git push --tags`. Line: claude-opus-4-8 /
   high. The workflow is YAML the gate cannot execute, and per `craft-line`
   work the gate cannot grade bumps a tier and gains effort.

## Implementation decisions

**Deviation from the map — needs explicit reviewer sign-off.** The map's #4
answer and asset priced in an "esbuild-style repair fallback in the
launcher" as the mitigation for npm's optional-deps lockfile edge, and the
Handoff flagged "verify launcher repair behavior against current npm during
slice 1". This spec deviates: the launcher ships the named-package error
(stories 4–5) as the only in-slice mitigation, the auto-repair download
moves to Out of scope as a separate capability, and the npm verification
happens as the first-release manual smoke rather than inside the slice.
Reason: auto-repair imports network, registry, and hash-trust machinery
into the launcher whose failure modes are larger than the failure it
repairs, and the named error is the floor the lockfile edge actually needs.
If you want the map's original shape instead, story 4 grows into the
repair capability and this spec's estimate roughly doubles.

- **Module shape.** One Go module at the repo root; `cmd/bench` holds the
  entrypoint; packages land under `internal/` as later slices need them.
  One binary, subcommand-dispatched — mirrors the shell CLI's shape.
- **Toolchain pinning (settles Handoff uncertainty flag).** `go.mod`
  carries the `go` directive at the current stable minor and a `toolchain`
  directive pinning the exact version; the builder verifies what current
  stable is at build time. Local gate runs and CI both inherit the pin
  through `GOTOOLCHAIN=auto` semantics, and CI reads the version from
  `go.mod` rather than duplicating it.
- **Version has one source.** `package.json` `version` is canonical.
  Release and gate builds stamp it into the binary via linker flags; an
  unstamped build prints `dev`. The gate builds with the stamp and asserts
  `bench version` output equals `package.json` exactly — the one-source
  check is executable. No shell `version` subcommand exists today, so
  binary-only creates no two-implementations window (the map's "ported"
  wording was loose; this is new surface).
- **Reproducible builds.** Every binary build — gate and release — uses
  `-trimpath` with the pinned toolchain, per the map's watch-out that Go
  builds are reproducible only under both; the build flags live in one
  place both callers share.
- **The gate rebuilds `dist/` every run** (stamped, `-trimpath`) before the
  runtime contracts execute — this is what makes a torn or stale dev binary
  unable to survive into an assertion, and it backs the SIGINT exclusion in
  the edge inventory.
- **Resolution order (the strangler router).** Inside the existing CLI
  entry, a ported subcommand resolves its binary as: (1) the repo-local
  dev build in a gitignored `dist/` (kit checkout), (2) the platform
  package nested under the wrapper's own `node_modules`, (3) the hoisted
  sibling layout npm produces in global installs. First executable match
  wins; a present-but-non-executable or empty file is treated as missing
  (falls through to story 4's error, never exec'd). The router is written
  once and takes the subcommand as data — later slices add names, not
  mechanisms.
- **Platform matrix.** `darwin-arm64`, `darwin-x64`, `linux-x64`,
  `linux-arm64` — matching the wrapper's existing `os` field. The matrix
  is declared exactly once (the generator's input) and everything else —
  platform package names, wrapper `optionalDependencies`, the workflow's
  build matrix, the gate's cross-compile loop — derives from or is checked
  against it.
- **Scoped package names.** `@benchkit/<os>-<arch>`, the ecosystem
  convention. **Reviewer-owned prerequisite:** the `@benchkit` npm scope
  must be created (or confirmed owned) before the first release; the spec
  does not gate on it, the release does.
- **Generated packages are not committed.** The generator emits the four
  package directories into a build directory at release time; the gate
  runs it in dry-run mode against a temp dir to validate shape and
  idempotency. Committing them would be a second source for the matrix.
- **Gate posture for Go: hard dependency.** Unlike shellcheck's
  best-effort, a missing Go toolchain fails the kit gate once `go.mod`
  exists — the core is load-bearing, and a gate that silently skips it is
  an always-pass. (Consumers are unaffected; the kit gate never ships.)
- **Publish with provenance.** The workflow publishes via GitHub OIDC
  provenance, satisfying the auditability posture recorded in the map's #4
  asset.
- **`bench version` is binary-only.** No shell fallback that reads
  `package.json` — a fallback would mask exactly the resolution failures
  this slice exists to surface.

## Testing decisions

- A good test here drives the public `bench` entry or the generator CLI
  and asserts stdout/stderr text, exit codes, and produced file shape —
  never Go internals. Go unit tests are additional (table tests inside the
  module), but acceptance lives at the shell-visible seams.
- Seams: the **version-routing seam** (fabricated-layout sandbox, prior
  art: `gate-runtime-contracts.sh`'s tmp-sandbox pattern), the **packaging
  generator seam** (dry-run into a temp dir), and the **gate's Go layer**
  (canary fixtures, prior art: `tests/canary/`).
- Gate command: `.bench/gate.sh` (the project gate), which after this
  slice includes the Go layer itself.

### Seam diagram — version-routing seam

    trigger: user/hook runs `bench version` (dev checkout, npm local, npm global)
        │
        ▼
    argv ─────────▶ [ CLI entry: strangler router ] ──▶ exec dist|bundled|hoisted binary
    filesystem ───▶ [ resolution: dist/ → bundled  ] ──▶ stdout: `benchkit <v> (<os>/<arch>)`, exit 0
    (layouts)       [   pkg → hoisted sibling      ] ──▶ stderr: named-package remedy, exit 127
                    [                              ] ──▶ stderr: unsupported platform, exit 2
        ◀ tests attach here: sandbox fabricates each layout (marker-stub binaries
          that print which target ran), invokes `bench version`, asserts
          stdout/stderr/exit

### Seam diagram — packaging generator seam

    trigger: release workflow (real run) / gate (dry-run)
        │
        ▼
    target matrix ──▶ [ platform-package generator ] ──▶ 4 package dirs (package.json:
    wrapper pkg.json ▶[ one template, matrix-driven ]     name, version, os, cpu, bin)
                      [                             ] ──▶ wrapper optionalDependencies check
        ◀ tests attach here: gate runs dry-run into a temp dir, asserts the four
          shapes, version equality with the wrapper, exact matrix↔deps match,
          and byte-identical output on a second run

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1,2 | gate-built binary: `bench version` prints `benchkit <pkg.json version> (<GOOS>/<GOARCH>)`, one line, exit 0 | version-routing | contract run before the module exists fails on unknown subcommand | version drift or an unstamped build breaks the exact-equality assert against `package.json` |
| 2 | `bench version` succeeds with cwd outside any git repo | version-routing | run from a non-repo tmp dir before the subcommand exists → non-zero | proves the subcommand skips repo-root resolution |
| 3 | dev build preferred over bundled; bundled over hoisted (marker stubs identify which ran) | version-routing | run before the router lands → unknown subcommand | wrong precedence prints the wrong marker |
| 3 | router works when the kit path contains spaces | version-routing | sandbox under a spaced dir before quoting is right → exec failure | unquoted expansion cannot resolve the spaced path |
| 3 | router works when `bench` is invoked via symlink | version-routing | symlinked invocation before resolution handles it → wrong root, error | resolution from the symlink's target, not its location (existing runtime-contract fake-`readlink` pattern) |
| 4 | no binary anywhere → stderr contains `@benchkit/<os>-<arch>` and an `npm install` remedy, exit 127 | version-routing | empty sandbox before the error path exists → generic shell noise without the package name | asserts the exact package substring and exit code |
| 4 (edge) | present-but-empty or non-executable binary file → treated as missing, same story-4 error, never exec'd | version-routing | sandbox with `touch`ed stub before the `-x`/non-empty check → exec attempt fails differently | distinguishes the absent-vs-empty hostile-input class |
| 5 | unsupported os/arch → `unsupported platform` on stderr, exit 2, no package named | version-routing | sandbox faking an off-matrix uname/arch before the check → story-4 error names a package that doesn't exist | the two failure modes must not alias |
| 6 | dry-run emits 4 package dirs; each package.json carries name/os/cpu/bin and version == wrapper version; wrapper `optionalDependencies` == matrix exactly; second run byte-identical | generator | gate check run before the generator exists → missing files | any hand-edit or matrix drift breaks equality; non-idempotency breaks the byte compare |
| 7 | gate goes red on `gofmt` diff, `go vet` failure, failing `go test`, or any of the four cross-targets failing to build; gate goes red when `go.mod` exists but no toolchain is on PATH | gate Go layer | run the new checks before wiring → gate stays green on a broken fixture | an unwired check is an always-pass; the canary harness expects red |
| 8 | canary fixtures (build-breaking file; failing test) each turn the gate red with their targeted error substring | canary harness | canary added before the check exists → harness fails expecting red | proves the Go checks bite, per the kit's canary discipline |
| 9 | workflow file exists, triggers on version tags, derives targets from the matrix, publishes 5 packages with provenance | static asserts only | not TDD-able beyond structure — the gate cannot execute Actions; red signal is the structural assert failing before the file exists | catches deletion/rename drift; behavior is covered by the manual smoke below |

### Edge inventory

Walked per behavior: error path, empty/absent, boundary, malformed,
interrupted/partial, idempotency, hostile environment — plus the profile's
hostile-input checklist. Resolved as rows above, except:

- **Won't handle: real `npm i -g` lifecycle across the 5 packages,
  including running it twice (repeat-install idempotency, Handoff item 6)**
  — the gate fakes layouts and cannot run a registry install; one manual
  smoke on the first release covers install and re-install (already the
  map's declared gate-blind spot).
- **Won't handle: reproducing npm's cross-platform lockfile omission** —
  needs a real npm + network; story 4's named-package error is the
  designed mitigation, verified once manually at first release.
- **Won't handle: SIGINT mid-build leaving a partial `dist/` binary** — the
  gate rebuilds `dist/` before the runtime contracts every run, so a torn
  binary cannot survive into an assertion.
- **Won't handle: `go.mod` absent** — unreachable after this slice lands in
  the kit repo (the only repo the gate grades), and the canary tree always
  carries the full kit.

## Out of scope

- **Auto-repair download fallback** (esbuild-style: launcher fetches and
  hash-verifies the missing platform tarball from the registry) — a
  distinct resilience capability with its own trust decisions; ~15 edits,
  ~6 gate runs. Parked on the roadmap at close.
- **Windows targets** — the wrapper's `os` field excludes win32 today;
  adding a platform is a separate capability rippling through matrix,
  workflow, and docs; ~8 edits, ~3 gate runs.
- **Porting any existing subcommand** — slices 2–8 of the map, each its own
  spec by the map's dependency order.
