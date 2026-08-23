# Split the freshness package files

Blocked by: none
Writes: internal/freshness

## What to build

Split `internal/freshness/freshness.go` (731 lines) along responsibility. The retained `freshness.go` keeps the seal type, the seal codec, and the public digest API. `freshness_publish.go` takes the atomic-publish lifecycle. `freshness_buildinputs.go` takes build-input discovery and the hardened symlink-safe readers. `freshness_verify.go` takes verify, check, and refusal messaging.

Split `freshness_test.go` (1109 lines) into digest, publish, and verify test files. The publish file keeps the child-process helpers. Pure moves only. The split adds no `Publish` caller and no `package main` file, so the AST topology test stays green.

## Acceptance

- [ ] R02: `bench structure` no longer lists `internal/freshness/freshness.go`.
- [ ] R05: `bench structure` no longer lists `internal/freshness/freshness_test.go`.
- [ ] R03: every created file counts at most 400 newlines.
- [ ] R04: `go build ./...` exits zero with no edit outside `internal/freshness/`.
- [ ] R08: `go test -list '.*' ./internal/freshness/` emits the same test-name set at base and tip.
- [ ] R17: `go test ./internal/freshness/...` exits zero.
- [ ] R12: `bench structure` prints no `DIR CROWDED` line for `internal/freshness/`.
- [ ] R18: `bench gate` exits zero before the commit.
