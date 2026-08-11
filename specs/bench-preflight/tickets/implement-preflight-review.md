# Implement bench preflight review end-to-end

Blocked by: rename-release-preflight-package.md, export-diff-review-base.md
Ownership fence: `internal/preflight/`, `internal/coverage/coverage.go`, `cmd/bench/main.go`, `bin/bench.sh`
Integration surfaces: freed package path→rename-release-preflight-package.md; review base→export-diff-review-base.md; row IDs + map violations→existing `internal/coverage/coverage.go` `ParseSpec` exercised by C4/C5; spec resolve + typed status→existing `internal/spec/spec.go` (`Resolve`, `Facts`) exercised by C2's bootstrap and the harden dependent; default branch→existing `internal/git/default_branch.go` `ResolvedDefault` exercised by C2; dispatch row + grammar→`cmd/bench/main.go`; shell route→`bin/bench.sh`; verdict core + bootstrap consumed by dependents→implement-preflight-build.md, harden-preflight-bootstrap-errors.md; routing-registry row and docs lines→advertise-preflight-kit-prose.md
Contracts: the resolved review base (base commit-ish + method + error kind/hint; method domain and absence semantics fixed by export-diff-review-base.md's contract) crosses `internal/diff`→`internal/preflight/`, asserted by C8 against the real exported resolver; the declared row set (ordered string row IDs + string violations slice; empty violations means valid; opt-in bool false means legacy map) crosses `internal/coverage`→`internal/preflight/`, asserted by C4 and C5 against the real `ParseSpec`; the resolved spec (content bytes + resolved path + typed `Status:` value from `Facts`; resolution failure is an error, never empty content) crosses `internal/spec`→`internal/preflight/`, asserted by C1's bootstrap against the real resolver; the resolved default branch (name string + ok bool; not-ok is a structured red) crosses `internal/git`→`internal/preflight/`, asserted by C2 against the real resolver
Closure: C1/base-current-row, C1/paths-authorized-row, C1/rows-owned-row, C1/rows-membership-row, C1/diff-nonempty-row, C1/rerun-byte-identical, C2/stale-base-red, C3/out-of-fence-red, C3/prefix-boundary, C4/uncited-row-red, C5/phantom-token-red, C5/foreign-tag-ignored, C6/empty-diff-red, C7/control-byte-refusal, C7/space-glob-path, C8/recorded-key-consumed, C8/diff-byte-identical, C12/tickets-absent-review-red, C12/special-file-refused, C13/not-in-repo-error, C13/missing-mode-usage, C13/unknown-mode-usage, C13/missing-slug-usage, C13/unknown-flag-usage, C14/spec-parser-newline, C14/ticket-parser-newline, C16/shell-route

## What to build

The new decision domain `internal/preflight`: a thin gatherer (git, the
exported diff-base resolution, spec resolver, coverage parser, `tickets/`
enumeration with lstat-mode refusal of special files) feeding a pure verdict
core over immutable facts, wired as `bench preflight review <slug>` —
`route_porcelain` in `bin/bench.sh`, a plain `commandRegistry` row in
`cmd/bench/main.go`, grammar through `usage.Parse`. Bootstrap fail-closed —
spec resolves, `Status: staged` (typed status via `internal/spec`'s `Facts`),
coverage map valid with row IDs opted in, `## Ownership fences` carries at
least one backticked entry outside parentheses — refusing with a structured
error on any failure; the *exact* per-artifact diagnostics are the follow-up
harden ticket's rows, so this ticket's bootstrap errors need only be
structured, fail-closed, and exit 1. Then the five checks — `base-current`,
`paths-authorized`, `rows-owned`, `rows-membership`, `diff-nonempty` — printed
as TOON (`phase:`, `spec:`, `checks[N]{check,verdict,detail}`), exit 0 all
green / 1 any red or structured error / 2 usage. Every structured error is one
`toon.Errorf` line; usage errors ride the usage parser or
`toon.Usage`/`toon.MissingArg`. Fence grammar: backticked tokens in
`## Ownership fences` outside parentheses authorize; prefix authorization only
with a `/` separator. `build` mode is not accepted yet (unknown mode, exit 2) —
the build-mode ticket expands the grammar. Repair
`internal/coverage/coverage.go:369`'s doc comment naming the retired
`assign`/`promote` consumers in the same change. Exemplars: verdict core
mirrors `internal/releasepreflight/decision.go` + `decision_test.go`; CLI
contract tests mirror `internal/coverage/coverage_test.go` `TestCommand`
(seeded throwaway repo, exact output and exit assertions). TDD at both marked
seams: CLI contract and verdict core.

## Acceptance

- [ ] [C1] (covers PF1) `bench preflight review <slug>` on a conformant seeded tree prints all five green rows — each check present by name — exits 0, and a second run is byte-identical.
- [ ] [C2] (covers PF2) with the default branch advanced past the branch point, `base-current` is the red row and exit is 1.
- [ ] [C3] (covers PF3) a tracked change outside every fence entry makes `paths-authorized` red naming the path; a path equal to a fence entry or under a fence prefix stays green, and `internal/git2` does not match fence `internal/git`.
- [ ] [C4] (covers PF4) one declared row ID cited by no ticket file makes `rows-owned` red naming the uncited ID.
- [ ] [C5] (covers PF5) a ticket token under the spec's own tag naming no declared row makes `rows-membership` red; a foreign-tag token (`FT93`) is ignored — both asserted.
- [ ] [C6] (covers PF6) an empty changed set in review mode makes `diff-nonempty` red.
- [ ] [C7] (covers PF7) a changed path carrying ESC exits 1 with the unrepresentable-TOON-cell error, never a mangled table or a silently sanitized path; a path with a space or glob character renders escaped and authorizes correctly.
- [ ] [C8] (covers PF8) with `branch.<name>.benchBase` recorded past an out-of-fence commit `paths-authorized` is green, with the key removed the same tree is red, and bare `bench diff` output stays byte-identical.
- [ ] [C12] (covers PF15) review mode with `tickets/` absent is a structured red; a FIFO inside `tickets/` is refused before reading, named in the error, without blocking.
- [ ] [C13] (covers PF16) outside a git repository the command answers the standard not-in-repo error; a missing mode, unknown mode, missing slug, and unknown flag each exit 2 with usage — all four branches asserted.
- [ ] [C14] (covers PF17) a spec whose last line lacks a trailing newline, and a ticket file whose last line lacks one, each parse identically to their terminated forms — both parsers asserted.
- [ ] [C16] (covers PF1) the tracer's invocation dispatches through the real shell: the `bin/bench.sh` route reaches the `commandRegistry` row (defense-in-depth on PF1; the conformance-graded registry-row and grammar claims stay with PF18).

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| C1/base-current-row | omit the `base-current` check from the verdict core | the CLI contract tracer test | `go test ./internal/preflight -run <tracer>`, expect the missing-row failure |
| C1/paths-authorized-row | omit the `paths-authorized` check | the tracer test | same run, expect the missing-row failure |
| C1/rows-owned-row | omit the `rows-owned` check | the tracer test | same run, expect the missing-row failure |
| C1/rows-membership-row | omit the `rows-membership` check | the tracer test | same run, expect the missing-row failure |
| C1/diff-nonempty-row | omit the `diff-nonempty` check | the tracer test | same run, expect the missing-row failure |
| C1/rerun-byte-identical | append a timestamp to the output | the tracer's second-run comparison | same run, expect the byte-identity failure |
| C2/stale-base-red | compare merge-base against HEAD instead of the default tip | the stale-base contract test | seed default advanced past branch point, run, expect the missed-red failure |
| C3/out-of-fence-red | make the authorization check return true unconditionally | the out-of-fence contract test | seed a change outside every fence, run, expect the missed-red failure |
| C3/prefix-boundary | authorize by plain string prefix without the `/` separator | the `internal/git2` boundary test | seed the boundary tree, run, expect the over-match failure |
| C4/uncited-row-red | skip the rows-owned comparison for the last declared ID | the uncited-row contract test | seed one uncited ID, run, expect the missed-red failure |
| C5/phantom-token-red | drop the membership check for tokens under the spec's tag | the phantom-token contract test | seed a `<tag>99` token, run, expect the missed-red failure |
| C5/foreign-tag-ignored | match every uppercase-tag token regardless of tag | the foreign-tag contract test | seed an `FT93` token, run, expect the false-red failure |
| C6/empty-diff-red | treat an empty changed set as green in review mode | the empty-diff contract test | seed no changes, run, expect the missed-red failure |
| C7/control-byte-refusal | strip control bytes from paths before they reach the sink, reporting a sanitized path as green | the ESC-path contract test | seed a path with ESC, run, expect the wrong-green failure — the required behavior is the unrepresentable-cell error, not a cleaned path |
| C7/space-glob-path | reject any path failing a word-character whitelist | the space/glob path contract test | seed `a b*.go` under a fence, run, expect the false-red failure |
| C8/recorded-key-consumed | re-derive the base by merge-base, discarding the exported resolver's result | the recorded-key contract test | seed benchBase past an out-of-fence commit, run, expect the wrong-verdict failure |
| C8/diff-byte-identical | change how `resolveBranchRange` consumes the export so `bench diff` output shifts | the existing `internal/diff` tests | `go test ./internal/diff`, expect the regression failure |
| C12/tickets-absent-review-red | report absent `tickets/` as not-applicable in review mode | the absent-tickets review contract test | seed no tickets dir, run review, expect the missed-red failure |
| C12/special-file-refused | open each ticket entry without the lstat-mode check | the FIFO contract test with a test-level timeout bound | seed a FIFO in `tickets/`, run under the bound, expect the hang-or-missed-refusal failure |
| C13/not-in-repo-error | skip the in-repo precondition | the not-in-repo contract test | run outside a repository, expect the missed-error failure |
| C13/missing-mode-usage | exit 1 with a structured error for a missing mode | the missing-mode usage test | run `preflight` bare, expect the exit-code failure |
| C13/unknown-mode-usage | accept any mode string as review | the unknown-mode usage test | run with mode `frob`, expect the missed-usage failure |
| C13/missing-slug-usage | default a missing slug to the newest staged spec | the missing-slug usage test | run `preflight review` bare, expect the missed-usage failure |
| C13/unknown-flag-usage | ignore unrecognized flags | the unknown-flag usage test | run with `--frob`, expect the missed-usage failure |
| C14/spec-parser-newline | key the spec scanner to `\n` so the unterminated last line drops | the spec trailing-newline test | seed a spec whose fences section ends unterminated, run, expect the dropped-line failure |
| C14/ticket-parser-newline | key the ticket scanner to `\n` | the ticket trailing-newline test | seed a ticket whose only citation of a row is the unterminated last line, run, expect the dropped-citation failure |
| C16/shell-route | remove the `bin/bench.sh` case for `preflight` | coordinator smoke through the real `bench` wrapper | `bench preflight review <slug>` in a seeded repo, expect `unknown subcommand` |
