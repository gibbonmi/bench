# Export the embed-target derivation

Blocked by: none
Writes: internal/packagesurface/assets.go, internal/packagesurface/assets_test.go

## What to build

The package `internal/packagesurface` exports one derivation over a checkout's
`//go:embed` targets. `EmbedTargets(root)` reads every non-test Go source under
the module root, under `cmd/`, and under `internal/`. It resolves each literal
embed pattern against that source's own directory. It returns the repo-relative
slash path of each target, and it returns an error the caller reports.

`RequiredBuildPackAssets` composes `EmbedTargets` and keeps its own sorted asset
list. The private helper `embeddedPackAssets` and its second walk over the Go
sources go away, so the embed rule has one source.

This ticket fixes the contract that `select-kit-lane-checks.md` consumes. The
`go-build-input` class calls `EmbedTargets` with the graded root and tests a
changed path for membership. So the returned paths carry the repo-relative slash
form that a composed change carries, and they carry no `./` prefix.

The tests attach to `internal/packagesurface/assets_test.go` beside the existing
pack-asset tests. Each test builds a fixture tree with a Go source that declares
one embed directive.

## Acceptance

- [ ] PL26: `EmbedTargets` over a source `internal/adopt/link_hook.go` that declares `//go:embed prepush.sh` returns `internal/adopt/prepush.sh`.
- [ ] PL27: `RequiredBuildPackAssets` over that fixture lists `internal/adopt/prepush.sh`.
- [ ] the gate `test` phase stays green for the whole `internal/packagesurface` package.
