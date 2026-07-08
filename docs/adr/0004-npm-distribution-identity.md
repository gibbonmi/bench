# Bench's npm distribution identity is `redbench`, distinct from the product name

The coded npm names were never ownable: unscoped `benchkit` is an unrelated
third-party package, and the `@benchkit` org is taken. So the npm distribution
identity is the reviewer-owned `redbench` — the wrapper package `redbench` (bins
`bench` and `redbench`) and the four `@redbench/<os>-<arch>` platform packages,
published under the `redbench` org. This is a distribution rename only: the product
stays Bench — the `bench` command, the version banner, the repository, the docs
branding, and the project profile name are all unchanged. Only the strings npm
resolves moved.

Publishing is deferred, so until the first tag the advertised install path is
npx-from-git: the wrapper's `prepare` script Go-builds the core on the consumer's
machine at install time and the built binary rides in the pack, so `npx
github:<repo>#<ref> link` works against the unpublished packages. This keeps the
working release pipeline and npx adoption alive without publishing before the
reviewer chooses to. The consumer therefore needs repository access plus Node and
Go; a failed build fails the install honestly rather than being masked.

## Consequences

- **The registry publish must come only from the tag-driven CI checkout.** Because
  the pack now ships the built binary, publishing from a developer's built tree
  would fatten the wrapper with one platform's binary. The CI checkout never carries
  a built binary, so it stays the single publish path.
- **`--access public` in the release workflow is load-bearing** — scoped packages
  default to private, and a private platform package would 404 the install.
- **The `redbench` wrapper bin entry is load-bearing.** npx resolves a multi-bin
  package by matching bin name to package name, so dropping the `redbench` bin
  silently breaks the advertised `npx redbench` one-liner once published.
- **The advertised npx git spec pins a ref**, because npx caches git installs and an
  unpinned spec would serve a stale build.
- **A shipped-surface gate check fails on any surviving unowned-identity string**
  (`@benchkit`, the `npx`/`npm i`/`npm uninstall` invocations of `benchkit`, or an
  npm-package claim naming `benchkit`) outside the changelog, so a half-done rename
  cannot ship.
- **The name is secured only by publishing**: until the first tag, unscoped
  `redbench` remains claimable by anyone.

## Considered options

Claiming the `@benchkit` org (taken); de-advertising npm and leading with the clone
path (idles a working pipeline and drops npx adoption); publishing immediately
(deferred by reviewer choice); npx from release-asset tarballs (the repository is
private and npx cannot authenticate to asset downloads); a context-aware missing-
binary remedy that detects install mode (branching in a deliberately minimal POSIX
wrapper); and a full product rebrand to Redbench (doc-wide churn for no distribution
gain). Until publish day, the missing-binary error prints only the clone/git build
remedy — the npm remedy would name packages that 404 — and rejoins when the packages
exist.
