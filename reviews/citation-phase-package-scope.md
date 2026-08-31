# Review: citation-phase-package-scope (FT281)

Base `55ae66f0`, tip `f463c1c5`, worktree `ft281-citation-scope`. Initial review.

## Repair-scoped re-review

Prior tip `f463c1c5`, repair tip `996e886d`. All 9 initial repair targets
verify closed. Scope is the repair delta only, per the repair-scoped rule.

### Standards

1 finding. `internal/gate/tag_census_test.go` hand-copies the race phase's
package paths (`./internal/worktree`, `./internal/guards`, `./internal/census`)
instead of deriving them from the `racetests.Tests[]` registry. AGENTS.md:
"an executable registry" stays single-sourced. The independent-expectation
exception needs a demonstrated red for the copied form; none was shown. —
`auto-fix`

### Spec

2 findings.

1. **Phase env has no acceptance row.** The approved phase-env behavior
   shipped with no coverage-map row, no `Covers:` update, and no verification
   log note — the written spec still ends at PS12. — `auto-fix` (add the row,
   the decision line, and a log note)
2. **An empty env value false-reds.** `contextFor` reads `value == "1"` for
   `CGO_ENABLED`, so `{"env":{"CGO_ENABLED":""}}` forces cgo off; real `go`
   treats an empty value as unset and still runs a `cgo`-tagged file. Same
   shape for an empty `GOOS`/`GOARCH`. — `auto-fix`

### Coverage

3 findings. Worst: a cross-`GOARCH` phase silently keeps the host's
`CgoEnabled`, so it can accept a citation the phase can never really compile.

1. **`CgoEnabled` ignores a `GOARCH` override.** A manifest phase with
   `{"env":{"GOARCH":"386"}}` and no `CGO_ENABLED` leaves `ctxt.CgoEnabled` at
   the host's `true`. On this host, `GOARCH=386 go env CGO_ENABLED` answers
   `0` — the real phase cannot build a `cgo`-tagged file, but `contextFor`
   still accepts one. — `auto-fix`
2. **Leading-dash normalization is untested.** `testFlagTakesValue` now
   accepts one or two leading dashes; no test exercises the two-dash form
   (`go test --skip X`). Behavior is correct (verified live), just unproven. —
   `auto-fix`
3. **`-test.`-prefixed flags still break the operand read** (e.g.
   `-test.run`). Pre-existing, not introduced by this repair, and outside the
   accepted repair predicates. — non-blocking follow-on, parked as an idea
   (`bench idea`); does not block this landing.

## Standards

5 findings. Worst: `packageLoadTimeout` redeclares a `bounds` policy value, and a
live conformance check already reds on it outside the fast lane.

1. **`packageLoadTimeout` redeclares bounds-policy value.** `internal/bounds/bounds.go`:
   "The const block is the production policy registry. Callers name these entries.
   They do not redeclare values or locally reimplement classification."
   `internal/coverage/citation_execution.go:25` declares
   `var packageLoadTimeout = 30 * time.Second`, the same expression as
   `GitRefreshTimeout`. Verified live: `bench test --check bounds-policy` in
   `ft281-citation-scope` reports
   `TestRootConformance: gate: internal/coverage/citation_execution.go redeclares
   GitRefreshTimeout policy value`. Not in the fast lane, so the ticket commit
   stayed green; the full gate will not. — `auto-fix`
2. **Canary fixture prose overclaims two pinned diagnostics.**
   `files/specs/unexecuted/spec.md`: "Together they prove both execution
   diagnostics reach the gate." `EXPECT` holds one line; `resolveFixtureBite`
   (`internal/conformance/fixture_bite_test.go:778-796`) asserts only that one.
   The row-1 tag diagnostic is materialized but unpinned. — `auto-fix`
3. **Duplicated Code — the `compiled.go` companion-file step, pasted twice.**
   `citations_test.go:217-221` and `citation_execution_test.go:205-209` each
   restate the same four-line "Go excludes an all-excluded package" rule and
   write the same companion file. One helper serves both call sites. —
   `auto-fix`
4. **Duplicated Code — the `-C` scan.** `goTestIndex` steps over `-C <dir>`
   pairs; `goCOf` independently re-walks `argv[1:testAt]` for the same pair.
   One scan can return both. — `auto-fix`
5. **`testFlagTakesValue` enumerates inclusions instead of naming the
   exception.** `projects/benchkit.md`'s hostile-input checklist convention:
   name the exception, don't enumerate inclusions. An unlisted value-taking
   flag in a consumer's own phase manifest (`-skip`, `-exec`, `-o`,
   `-gcflags`, `-ldflags`, `-toolexec`, `-fuzz`, `-outputdir`, ...) makes its
   value read as a package operand. — `ask-user` (converges with Spec #3 and
   Coverage #1 below — same root defect, three citations)

## Spec

4 findings. Worst: a relative spec path makes `bench coverage --check` false-red
every citation — a regression this diff introduces, invisible to the gate.

1. **Relative spec path breaks package selection.** Spec, Solution: "Accept a
   citation only when one entry selects the package and accepts the cited
   file." Traced execution (synthetic repo, one row citing
   `internal/x/present_test.go`): `coverage --check specs/demo/spec.md`
   (relative) → `which no executed test phase selects` (exit 1); the same
   binary with an absolute path → `ok`; the base commit (55ae66f0) with the
   relative path → `ok`. Cause: `specLocation` yields a relative `base`, so
   `checkCitation`'s `path` is relative, while `go list -f {{.Dir}}` always
   prints absolute directories — `containsDirectory` never matches. Only the
   documented seam's human/agent invocation breaks; preflight and the canary
   owner both pass absolute roots, so the gate stays green. — `auto-fix`
2. **Effective directory is asserted nowhere.** Spec, Testing decisions:
   "Assert each phase's tags, package operands, and effective directory."
   `internal/gate/tag_census_test.go` blanks `Dir`/`GoC` before comparing;
   `want` carries neither. No test exercises `goCOf` or the `-C`/`Dir` join. —
   `auto-fix`
3. **`-skip` missing from the value-taking flag list.** `testFlagTakesValue`
   lists `-run`/`-bench` but not `-skip`, `-fuzz`, `-exec`, `-o`. A manifest
   phase `go test -skip TestX ./...` collects `["TestX","./..."]` as package
   operands; `go list` exits nonzero; every row in that spec reports
   `could not expand packages`. The spec's Won't-handle defers filter
   *semantics*, not filtered phases. Two of three filters are already
   handled — this omission turns a legal phase into a total false red. —
   `auto-fix` (same defect as Standards #5 / Coverage #1)
4. **Canary fixture prose overclaims** (same defect as Standards #2). —
   `auto-fix`

Both reviewer-approved additions (the `go list` timeout, `gate.BaselinePolicyEnv`)
grade clean against the spec's intent.

## Coverage

4 findings. Worst: the value-taking-flag table is missing ~16 real `go test`
flags. A manifest using the separated form of any one of them false-reds
every citation in that phase.

1. **`-skip` and ~15 siblings become package operands.**
   `internal/gate/tag_census.go:149-161`. Missing: `-skip`, `-o`, `-exec`,
   `-gcflags`, `-ldflags`, `-toolexec`, `-cpuprofile`, `-memprofile`, `-trace`,
   `-outputdir`, `-fuzz`, `-gocoverdir`, `-installsuffix`, `-buildmode`,
   `-compiler`. Measured: `go list -f '{{.Dir}}' TestSlow ./...` → exit 1,
   `package TestSlow is not in std` (Go 1.25.0). No fixture anywhere uses the
   separated-flag-value form — every manifest test uses `-tags=x`. — `ask-user`
   (the reviewer call is allow-list-completion vs. degrade-to-"selects
   nothing" on an unresolvable operand; same defect as Standards #5 / Spec #3)
2. **The one production argv that depends on this table is untested.**
   `raceDriverArgv()` (`internal/gate/gate_go.go:170-178`) emits a real
   `-run <pattern>` phase. `TestExecutedTagCensus`'s "kit root" case builds a
   synthetic root holding only `go.mod`, so the race phase never
   materializes there. This table is never checked against its one real
   caller. — `auto-fix`
3. **`Phase.Env` is dropped from the census.** `manifestPhase.Env` exists
   (`internal/gate/manifest.go:59`), but `TestExecution` carries no env. A
   manifest phase with `"env":{"GOARCH":"386"}` citing a `_386_test.go` file
   on amd64 gives a false red. An env-gated `cgo` tag can symmetrically
   false-green. No row or Won't-handle entry decides phase env; this gap
   predates the diff. — `ask-user`
4. **Two `-C` occurrences.** `goCOf` takes the first, `goTestIndex` skips all;
   real `go` rejects a non-first `-C`, so the phase reds on its own before
   this matters. — `no-op`, not listed as a repair target.

Coverage's adversarial pass refuted four other angles — none are repair
targets:

- A symlinked repo root: resolved correctly.
- A `packages()` memoization race: the only caller is sequential.
- Stderr posing as a package directory: `bounds.RunOutput` separates it.
- The 30s timeout: a real test drives a real grandchild-holding stub, and
  there is one call site.
