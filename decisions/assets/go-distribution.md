# Go core distribution — research asset for `decisions/go-rewrite.md` #4

Acceptance bar (set by ticket #3): no toolchain requirement on consumer
machines, and an auditable consumer surface. Researched 2026-07-03.

## Repo facts the options were judged against

- npm ships the whole CLI as source today: `bin/`, `.bench/hooks/`,
  `.bench/lib/`, `.bench/adapters/` are all in `files[]`; `bin.bench` points
  at `bin/bench.sh`. `engines.node >= 22` (⇒ npm ≥ 10), `os: darwin, linux`.
- `bench link` copies the kit's `bin/` tree into consumer repos as
  `.bench/bin/` — consumers currently carry the full CLI source in-repo,
  fingerprinted by the link manifest.
- The shared harness hooks resolve the CLI through a chain
  (`.bench/bin/bench.sh` → kit `bin/bench.sh` → `bench` on PATH) and are
  exec'd by path from `.claude/settings.json` / `.codex/hooks.json`.
- The pre-push hook is **self-contained inline shell** generated at link
  time — it never execs `bench`. The branch-protection invariant has no
  binary dependency and survives any language change untouched.
- Precedents already in the kit: a postinstall that must never fail the
  install, and a doctor-written PATH shim (marker + `exec` target).

## Options

**A. npm platform packages (esbuild pattern) — recommended.** Prebuilt Go
binaries published as scoped packages (`@benchkit/darwin-arm64`,
`@benchkit/darwin-x64`, `@benchkit/linux-x64`, `@benchkit/linux-arm64`),
each carrying `os`/`cpu` fields and listed as `optionalDependencies` of
`benchkit`; npm installs only the matching one. `bin/bench` becomes a thin
launcher that resolves and execs the platform binary. Established ecosystem
convention (esbuild, turbo, biome, swc); Go cross-compiles all four targets
from one CI job (`GOOS`/`GOARCH`). Meets both bars: consumers need nothing
beyond the node ≥ 22 they already need, and every in-repo executable stays
readable text (see consumer surface below).

**B. postinstall download from GitHub releases — rejected.** Breaks under
`--ignore-scripts` and network-restricted CI as the *primary* path;
supply-chain smell. Acceptable only as A's repair fallback (esbuild verifies
such downloads against a `binaryHashes` map in the wrapper's package.json).

**C. `go install` / source build — rejected.** Toolchain on consumer
machines; fails the #3 bar outright.

**D. binaries committed into linked repos — rejected.** Megabytes of
unauditable blob per platform inside the consumer's own git history — worse
on both bars, plus repo bloat and update churn.

## Consumer surface under A

- `bench link` stops copying CLI source; it plants thin shell shims in
  `.bench/bin/` (marker + resolve + `exec "$target" "$@"` — same shape as
  the doctor shim) and keeps planting the content surface unchanged.
- Harness hook entries stay `.sh` files at the paths the adapters name;
  their logic moves into binary subcommands (`bench hook session-start`,
  …), the shim just execs. A consumer reads every line that runs as text
  in their repo; the binary is a versioned npm artifact built from public
  source — auditable the way their other dependencies are.
- pre-push: unchanged, still generated self-contained shell.
- Version skew: the link manifest gains a kit-version stamp;
  `session-start`/`bench doctor` compare it against `bench --version` and
  advise relink on mismatch (binary and planted assets always originate
  from the same package version, so skew only means "repo linked by an
  older/newer kit than this machine runs").

## Kit-repo consequences

Go toolchain becomes a **dev/CI dependency only**. The gate's parse layer
becomes `go build ./... && go vet && go test` for the core plus the existing
`bash -n`/shellcheck for what stays shell (shims, pre-push generator,
postinstall, doctor PATH shim). Release pipeline: one cross-compile matrix,
five npm publishes (wrapper + four platform packages).

## Residual risks (feed #5)

- npm's optional-deps lockfile edge (lockfile created on one platform can
  omit another's optional dep; largely fixed in npm ≥ 9 but still reported
  in the wild under pnpm/yarn mixes). Mitigation: esbuild-style repair
  fallback in the launcher plus a definitive error naming the missing
  package.
- `--ignore-scripts` consumers lose the repair fallback — same posture the
  kit already accepts for the PATH shim.
- Four-target release matrix is new operational surface; a broken platform
  package fails loudly at launcher exec, not silently.

## Sources

- [esbuild platform-specific binaries](https://deepwiki.com/evanw/esbuild/6.2-platform-specific-binaries)
- [esbuild getting started — install internals](https://esbuild.github.io/getting-started/)
- [npm optional-deps platform bug in the wild](https://github.com/AndyMik90/Aperant/issues/593)
- [Platform strings (Feb 2026 survey of the convention)](https://nesbitt.io/2026/02/17/platform-strings.html)
