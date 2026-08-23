# Split the systemtest owner tests

Blocked by: none
Writes: internal/systemtest

## What to build

Split `internal/systemtest/owner_test.go` (884 lines) into four files. The harness file keeps `TestMain` and the whole `systemOwner` type. The selection-journey file also takes `systemReauthorizeEvidence`. The land-race file keeps its fixture helpers, and the teardown file takes the rest.

Every new file repeats the `//go:build system` tag. The harness stays whole because the four sibling test files call it. Pure moves only.

## Acceptance

- [ ] R07: `bench structure` no longer lists `internal/systemtest/owner_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R09: `go vet -tags system ./internal/systemtest/...` exits zero.
- [ ] R10: exactly one `TestMain` definition remains under `internal/systemtest/`.
- [ ] R08: `go test -tags system -list '.*' ./internal/systemtest/` emits the same test-name set at base and tip.
- [ ] R11: each moved fixture has exactly one definition in the package.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/systemtest/`.
- [ ] R18: `bench gate` exits zero before the commit.
