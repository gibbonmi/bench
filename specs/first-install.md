# First-install fixes

Status: staged

## Problem

The Bench kit never installs itself cold, so a family of defects that only appear on the advertised npm/npx adoption path are invisible in this repo:

1. `bench init` and `/bench-setup-repo` point a fresh consumer at `projects/<name>.md` profile templates, but `projects/` is not in `package.json` `files[]` — on the npm path those files don't exist and the agent improvises.
2. The tarball ships the kit's own working docs — `HANDOFF.md`, `CLAUDE.md`, `AGENTS.md` — which adoption never reads (`link` generates consumer copies from Go constants). A stale kit handoff is exported to every consumer.
3. `bench link` copies the running arch-specific binary to `.bench/dist/bench` with no ignore entry. A consumer who commits it hands a different-arch teammate a broken CLI, and the Stop-hook oracle then fails open silently.
4. The scaffolded `.bench/gate.sh` calls bare `bench canary`. A machine relying only on the repo-local CLI (the documented "global bench optional" story) gets a gate that can never go green.

Each defect is a real fresh-install failure with no test that installs cold to catch it.

## Solution

Land the four fixes at their owning seams, then close the visibility gap with two coverage layers that install cold:

1. Add `projects/` to `files[]` so the profile templates ship as-is.
2. Drop `HANDOFF.md`, `CLAUDE.md`, `AGENTS.md` from `files[]`.
3. `bench link` writes a manifest-tracked `.bench/dist/.gitignore` that ignores the copied binary, alongside the copy.
4. The scaffolded gate resolves the repo-local CLI (`.bench/bin/bench.sh`) first, falling back to `bench` on PATH.
5. Prove it stays fixed: pin `files[]` contents in the package-surface contract test, and add a fresh-install runtime contract test that links + inits a throwaway repo with PATH stripped to essentials and asserts the scaffolded gate resolves and the profile templates are reachable.

## User stories

1. As someone installing benchkit from npm, I want the `projects/` profile templates to ship, so that `bench init` points at files that actually exist and my agent doesn't improvise a profile from nothing. Adding `projects/` to `files[]` also makes the profile markdown part of the shipped surface swept by the conformance repo-only-claims check, so that sweep must stop treating `projects/` as a repo-only path (it now genuinely ships). Line: claude-opus-4-8 / medium. This touches the conformance oracle's model of what is repo-only versus shipped, and getting the sweep's path list wrong either weakens a real shipping-claim guard or falsely reds a true claim.

2. As someone installing benchkit, I want the kit's own working docs (`HANDOFF.md`, `CLAUDE.md`, `AGENTS.md`) left out of the tarball, so that I never receive a stale kit handoff or a duplicate agreement — `link` already generates my consumer-facing copies from constants. Line: claude-sonnet-5 / low. This is a `files[]` removal plus a test-list edit, fully gate-observed with no judgment call.

3. As a consumer who commits `.bench/dist/`, I want `bench link` to write a `.bench/dist/.gitignore` that ignores the arch-specific binary, so that I can't hand a different-arch teammate a broken CLI and a silently fail-open Stop hook. The ignore file is manifest-tracked and idempotent across relinks like every other managed file. Line: claude-sonnet-5 / medium. This is plumbing at the known `buildLinkPlan`/`installPlan` seam, but the manifest fingerprinting and relink idempotency need care.

4. As a machine relying only on the repo-local CLI with no global `bench` on PATH, I want the scaffolded `.bench/gate.sh` to resolve `.bench/bin/bench.sh` before falling back to `bench` on PATH, so that my gate can reach `bench canary` and go green. Line: claude-sonnet-5 / low. This is one POSIX-sh resolver line in the scaffold template string at a known seam.

5. As a Bench maintainer, I want the fresh-install path proven by contract tests — a `files[]` surface pin and a PATH-stripped fresh-install smoke — so that this whole class of cold-install defect can't regress unseen in a kit that never installs itself. Line: claude-opus-4-8 / medium. The smoke is the deep piece; getting the stripped-PATH environment and the gate-resolution assertion right is the crux that makes it a genuine oracle rather than a false green.

## Implementation decisions

**Slice 1 — `projects/` ships.** Add `projects/` to `package.json` `files[]`. Pin its presence through the shared `packagesurface.RequiredPackAssets` list (add `projects/gl-axi.md`, a stable example profile) so both the surface contract test and the conformance package check require it in the pack — one source. Because shipped markdown is now swept by `conformance.checkRepoOnlyPackageClaims`, and `projects/benchkit.md` states the profiles are "shipped templates," drop `"projects/"` from that sweep's repo-only path list (`{"projects/", "specs/", "decisions/", "tests/"}` → `{"specs/", "decisions/", "tests/"}`). This is the correct fix, not rewording prose: `projects/` is no longer a repo-only path, so a claim that it ships is now true and must not be flagged. `specs/`, `decisions/`, and `tests/` remain repo-only and stay in the list. *Author-discovered — the decision map was silent on the sweep interaction; flagged for veto.*

**Slice 2 — kit docs stop shipping.** Remove `HANDOFF.md`, `CLAUDE.md`, and `AGENTS.md` from `files[]`. Pin their absence from the tarball. To keep the "does-not-ship" fact single-sourced rather than split between the surface test's inline forbidden list and conformance's lone `.claude/settings.local.json` entry, introduce `packagesurface.ForbiddenPackAssets` (the union of the existing forbidden entries plus the three kit docs) and have both consumers iterate it. *The shared-forbidden-list consolidation is an author call the map didn't name; the minimal fallback is adding the three docs to the surface test's existing inline forbidden loop. Flagged for veto.*

**Slice 3 — link writes `.bench/dist/.gitignore`.** The kit has no `.bench/dist/.gitignore` source to copy (its own `.bench/dist/` is gitignored), so the linked file is generated content — a constant that ignores the copied binary and nothing else (a single `bench` line; it must not ignore itself, since the ignore file is meant to be committed and travel with the repo). Emit it from `buildLinkPlan` alongside the `.bench/dist/bench` entry so `installPlan` fingerprints it and records its manifest row uniformly; the cleanest fit is an inline-content plan entry (a `planEntry` that carries generated content instead of a copy source) so the manifest and idempotency logic apply without a special case. **Manifest uncertainty (Handoff item 7) resolved:** no format or version bump — the manifest is a `#kit\t<version>` header plus `<rel>\t<fingerprint>` rows, so the new managed file is one additional row, compared by fingerprint on relink like any other. **Hostile-input owner (Handoff item 6):** a fresh link writes the file; an unchanged relink rewrites it byte-identically (idempotent, no duplication); a user hand-edit trips the existing modified-managed-file conflict and link refuses rather than clobbers; a pre-existing project-owned file before the first link trips the existing project-owned conflict and link refuses. Refusing, not clobbering, is the decided posture for "must not clobber a user's wider ignore."

**Slice 4 — scaffolded gate resolves the local CLI.** In `adopt.scaffoldGate()`, prepend one POSIX-sh resolver line so the canary invocation uses a resolved command instead of bare `bench`: `bench="$(dirname "$0")/bin/bench.sh"; [ -x "$bench" ] || bench=bench`, then `"$bench" canary "$root"`. The decision named the candidate `"$(dirname "$0")/bin/bench"`, but the local wrapper that `link` installs is `.bench/bin/bench.sh`; the resolver targets the file that actually exists (`bin/bench.sh`). *Path-precision correction flagged for veto.* The two-candidate inline resolver (local wrapper → PATH `bench`) is intentional per the decision — it does not source the shared `resolve-bench.sh`, keeping the user-owned scaffold self-contained.

**Slice 5 — coverage.** Surface layer is the `files[]` pins above (stories 1 and 2). The fresh-install smoke is a new runtime contract test in `internal/contract/runtime` using the built `dist/bench` per the runtime-contract convention (guarded by `SkipIfSubjectBenchMissing`; rebuild `dist/bench` before hand-running). It creates a throwaway git repo, runs the built binary's `link` then `init` with PATH stripped to `/usr/bin:/bin` (git and bash present, no global `bench`), then runs the scaffolded `.bench/gate.sh` under the same stripped PATH and asserts it resolves `.bench/bin/bench.sh` and reaches the canary. The PATH strip is load-bearing: a smoke that leaves a global `bench` on PATH false-greens against the old bare-`bench` scaffold. The full `npm pack`+install smoke in the gate is deliberately not built (decision #5: npm and network variance for coverage the surface pin already gives).

## Testing decisions

- **What a good test is here:** exercise the shipped surface and the fresh-install behavior from outside — parse `npm pack --dry-run` output, drive `bench link`/`init` in throwaway repos, and run the scaffolded gate in a stripped environment — never assert against internal function state.
- **Seams and prior art:** four seams, all with close prior art. Seam A extends `testPackageNpmPackInstallableSurface` (`internal/contract/surface/package_test.go`) and the shared `packagesurface` lists. Seam B adds a case to `TestLinkContracts` (`internal/contract/surface/link_test.go`), following `testLinkSafeFreshRelink`. Seam C is the existing `checkRepoOnlyPackageClaims` under `TestRootConformance`. Seam D is a new fresh-install test in `internal/contract/runtime`, following `testRuntimeSymlinkedKitDir` (link in a temp repo under a restricted PATH) and `copiedCLIHookFixture` (built-binary use).
- **Gate command:** the project gate, `.bench/gate.sh`. Per-seam `go test` invocations below are the red signals; the done oracle is the full gate green.

### Seam diagram

Seam A — package pack surface (stories 1, 2):

    trigger: gate conformance/surface phase — `go test ./internal/contract/surface -run TestPackageContracts`
        │
        ▼
    package.json files[]  ──▶  [ npm pack --dry-run --json ]  ──▶  shipped file set
                                     │
                                     ▼
        assert vs packagesurface.RequiredPackAssets  (projects/gl-axi.md ∈ set)
                and packagesurface.ForbiddenPackAssets (HANDOFF/CLAUDE/AGENTS.md ∉ set)
            ◀ tests attach here: parse the pack JSON, require presence, forbid the kit docs

Seam B — link writes `.bench/dist/.gitignore` (story 3):

    trigger: gate safe-link phase — `go test ./internal/contract/surface -run TestLinkContracts`
        │
        ▼
    temp git repo  ──▶  [ bench link : buildLinkPlan → installPlan ]  ──▶  linked tree + manifest
                              │
                              ▼
        .bench/dist/bench (copied binary)  +  .bench/dist/.gitignore (ignores `bench`)
            ◀ tests attach here: assert the ignore file exists, ignores the binary,
              has a manifest row, and a second link leaves it byte-identical

Seam C — conformance repo-only-claims sweep (story 1):

    trigger: gate conformance phase — TestRootConformance → checkRepoOnlyPackageClaims
        │
        ▼
    files[] markdown (now incl. projects/*.md)  ──▶  [ repo-only-claims sweep ]  ──▶  diags (must be empty)
                                                            │
                                                            ▼
        repo-only path list drops "projects/" (now shipped, no longer repo-only)
            ◀ tests attach here: with projects/ shipped, the sweep must not flag
              projects/benchkit.md's "shipped templates" line

Seam D — fresh-install smoke (stories 4, 5):

    trigger: gate runtime phase — `go test ./internal/contract/runtime -run TestRuntimeFreshInstall`
             (PATH = /usr/bin:/bin — git + bash present, no global bench)
        │
        ▼
    temp git repo  ──▶  [ dist/bench link → init ]  ──▶  .bench/{bin/bench.sh, dist/bench, gate.sh}
                              │
                              ▼
        run scaffolded .bench/gate.sh under the stripped PATH
            ◀ tests attach here: assert the gate resolves .bench/bin/bench.sh (no
              "command not found", canary runs); assert the resolved kit dir holds projects/

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | a `projects/` profile ships in the npm tarball | A | `go test ./internal/contract/surface -run TestPackageContracts` fails `npm package missing projects/gl-axi.md` against the current tree (projects/ not in files[]) | the `RequiredPackAssets` pin reds whenever the pack omits a profile — exactly defect #1 |
| 1 | conformance stays green with `projects/` now shipped | C | adding `projects/` to `files[]` without dropping it from the sweep list makes `go test ./internal/conformance -run TestRootConformance` flag `projects/benchkit.md` as a repo-only shipping claim | the sweep's repo-only model must match the new `files[]`, or a now-true claim reds the gate |
| 2 | `HANDOFF.md`, `CLAUDE.md`, `AGENTS.md` are absent from the tarball | A | `TestPackageContracts` fails `npm package includes HANDOFF.md` against the current `files[]`, which lists all three | the forbidden pin reds on today's tree, so it can only go green once slice 2 removes them |
| 3 | `bench link` writes `.bench/dist/.gitignore` ignoring the binary, with a manifest row | B | a new `testLinkWritesDistGitignore` fails (file absent, no manifest row) before the `buildLinkPlan` edit | asserts the ignore file and its manifest row exist, so a committed binary is ignored — defect #3 |
| 3 | idempotent relink | B | a non-idempotent write reds: second `bench link` must leave exactly one byte-identical `.bench/dist/.gitignore` and no duplicate manifest row | proves the Handoff item-6 "must not duplicate or clobber" posture |
| 4 | the scaffolded gate resolves `.bench/bin/bench.sh` before PATH `bench` | D | with the old bare-`bench canary` scaffold, the PATH-stripped gate emits `bench: command not found` and the canary never runs | proves the resolver line, the only way a no-global-`bench` machine reaches a green gate — defect #4 |
| 5 | the fresh-install path works end-to-end; profile templates reachable | D | before slice 4 the PATH-stripped scaffolded gate cannot resolve the CLI; the ships-regression for templates is owned by Seam A, and the smoke additionally reds if the resolved kit dir lacks `projects/` | the smoke is the only test that installs cold with no global `bench` — the class the kit otherwise can't see |

**Degenerate-implementation check.** The always-green stub for story 1 (leave `projects/` out of `files[]`) reds Seam A; the do-nothing for story 2 (leave the docs in `files[]`) reds Seam A on today's tree; a link that skips the ignore file reds Seam B, and a non-idempotent one reds the relink row; a bare-`bench` scaffold reds Seam D under stripped PATH. The cheapest wrong smoke is one that forgets to strip PATH — it false-greens the old scaffold; the PATH strip is asserted as part of the fixture, not left implicit.

### Edge inventory

Edge classes walked per behavior; each resolved as a coverage row above or a **Won't handle** line here.

- Error path (link kit-asset missing, preflight conflict) — covered by existing `testLinkConflictWithoutManifest`; the new gitignore entry adds no new error path.
- Re-run idempotency — relink of `.bench/dist/.gitignore` is a coverage row; a second `bench init` does not rewrite an existing `.bench/gate.sh` (existing `testInitExistingGateIdempotence`), so the resolver line lands only on first scaffold.
- Hostile environment (no global `bench` on PATH) — the smoke owns it (Seam D, stripped PATH).
- **Won't handle:** pre-existing project-owned `.bench/dist/.gitignore` before the first link — the existing preflight conflict refuses safely rather than clobbering; no new surface.
- **Won't handle:** user hand-edit of the managed `.bench/dist/.gitignore` — the existing modified-managed-file protection refuses relink; a user's wider ignore is never silently overwritten.
- **Won't handle:** kit paths containing spaces or glob characters for the new gitignore write — inherits link's existing path handling proven by `testLinkMetacharKitPath`; the generated entry flows through the same `installPlannedFile` path.
- **Won't handle:** control bytes and missing-trailing-newline classes — not reachable by any of these five slices (no git-sourced text rendering, no hand-edited-line parsing).
- **Won't handle:** full `npm pack`+install smoke in the gate — rejected by decision #5 (npm and network variance in a 30s gate); the surface pin (Seam A) already owns the ships-regression.

## Out of scope

- **Deleting `HANDOFF.md` from the repo** — owned by the FT25 after-merge build, a separate capability; this spec only removes it from `files[]`, which is independently valid whether or not the file still exists (~1 edit, 1 gate run).
