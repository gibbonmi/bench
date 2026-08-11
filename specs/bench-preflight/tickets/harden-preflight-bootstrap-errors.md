# Harden preflight bootstrap errors to exact diagnostics

Blocked by: implement-preflight-review.md
Ownership fence: `internal/preflight/`
Integration surfaces: bootstrap + verdict core→implement-preflight-review.md; typed status producer→existing `internal/spec/spec.go` `Facts` exercised by H4; validator messages→existing `internal/coverage/coverage.go` `ParseSpec` violations exercised by H2; hardened errors consumed by dependents→implement-preflight-build.md, advertise-preflight-kit-prose.md
Contracts: the typed `Status:` value (string as `Facts` reports it; absence: a spec with no status line is itself a structured red) crosses `internal/spec`→`internal/preflight/`, asserted by H4 against the real `Facts`; the validator's violation messages (ordered string slice, empty means valid) cross `internal/coverage`→`internal/preflight/`, asserted by H2 against the real `ParseSpec`
Closure: H1/missing-spec-names-path, H1/dangling-symlink-classified, H2/validator-message-carried, H2/optin-hint-named, H3/fences-absent-error, H3/fences-empty-error, H3/paren-token-never-authorizes, H4/found-status-named

## What to build

Refine the core ticket's generic fail-closed bootstrap errors into the exact
per-artifact diagnostics the spec's error contract requires — each one
`toon.Errorf` line, exit 1: a missing spec or dangling symlink names the spec
path (lstat before read, so a broken link is classified as broken, never as an
authoritative empty state); an invalid coverage map carries the validator's
own message, and a legacy 5-cell map's error names the row-ID opt-in; an
`## Ownership fences` section absent, empty, or holding no backticked entry
outside parentheses gets its own named error; a non-`staged` spec's error
names the found status via `internal/spec`'s `Facts`. The parenthesized-token
authorization rule itself ships with the core ticket's grammar; this ticket
owns PF14's assertion of it. Behavior change plus its tests, same CLI-contract
seam and exemplars as the core ticket.

## Acceptance

- [ ] [H1] (covers PF12) a missing spec, or a dangling symlink where the spec should be, answers a structured error naming the spec path, exit 1 — the symlink case classified as broken, not empty.
- [ ] [H2] (covers PF13) a spec whose coverage map fails validation answers a structured error carrying the validator's message; a legacy 5-cell map's error names the row-ID opt-in — both exit 1.
- [ ] [H3] (covers PF14) an `## Ownership fences` section absent, empty, or containing no backticked entry outside parentheses answers a structured error, exit 1; a backticked token inside parentheses is never an authorization.
- [ ] [H4] (covers PF22) a spec whose `Status:` is anything but `staged` answers a structured error naming the found status, exit 1.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| H1/missing-spec-names-path | drop the spec path from the missing-spec error | the missing-spec contract test | run against no spec, expect the unnamed-path failure |
| H1/dangling-symlink-classified | stat through the link instead of lstat-first | the dangling-symlink contract test | seed a dangling link at the spec path, run, expect the misclassification failure |
| H2/validator-message-carried | replace the validator's message with a generic invalid-map string | the invalid-map contract test | seed a broken map, run, expect the missing-message failure |
| H2/optin-hint-named | omit the row-ID opt-in naming from the legacy-map error | the legacy-map contract test | seed a 5-cell map, run, expect the missing-hint failure |
| H3/fences-absent-error | treat a missing fences section as unrestricted authority | the absent-fences contract test | seed a spec without the section, run, expect the missed-error failure |
| H3/fences-empty-error | treat an empty fences section identically to a populated one | the empty-fences contract test | seed an empty section, run, expect the missed-error failure |
| H3/paren-token-never-authorizes | accept parenthesized backticked tokens as fence entries | the parenthesized-token contract test | seed a fences section whose only token is parenthesized, run, expect the false-authorization failure |
| H4/found-status-named | report a generic not-staged error without the found value | the non-staged contract test | seed `Status: implemented`, run, expect the unnamed-status failure |
