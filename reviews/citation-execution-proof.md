# Review pickup — citation-execution-proof

Base 006a28c3, reviewed tip 18a30271. Raw findings: 13. Repair targets after
collapse: 8. The Spec and Coverage axes found the same worst defect, so their
two findings are one target.

## Standards

Count: 6 (3 hard, 3 judgment). Worst: S2 — the citation predicate is derived
once per delegate fence, so the report and the checks can drift.

- S1 (auto-fix) — `internal/coverage/citations_test.go:399-411` and `:421-433`
  build the same three-row fixture twice. Collapse into one helper that both
  tests consume. Rule: one source per fact (`AGENTS.md`, fixture harness
  example).
- S2 (auto-fix) — `citesATest` (`citations.go:348-353`) re-derives the
  has-a-nonempty-name-list rule that `checkCitations` (`citations.go:85-91`,
  `:120-123`) already encodes. Extract one citation projection; make both
  consumers read it.
- S3 (auto-fix) — `TestHistoricalSpecSilencesTheNewChecks` and
  `TestParseSpecReturnsTheNewViolationClasses` carry diff-relative names, and
  the first one's comment says "this feature added". Rename to the state they
  pin; drop the narration. Rule: `craft-comments`, no narration.
- J1 (auto-fix) — `assertCensus` (`tag_census_test.go:74-88`) open-codes a
  deep comparison where `reflect.DeepEqual` is the package idiom.
- J2 (auto-fix) — `citations_test.go:442-450` iterates `[][]string` of
  one-string slices. Flatten to `[]string`.
- J3 (no-op) — `uncitedRows` re-checks the mapped state its caller already
  checked. Defensible for an unexported helper.

## Spec

Count: 2. Worst: F1 — the story-11 exemption keys on any `.Run(`, which
disables the story-10 diagnostic for 33 of 45 subtest-carrying files.

- F1 (ask-user) — `runNameRe` matches any selector `.Run(`
  (`citations.go:16`), so `cmd.Run()` or a raw-string `t.Run` exempts the
  whole file from segment resolution. CE10 is partial. The repair anchors the
  scan to the `testing.T` receiver, and that can newly red existing
  citations — a reviewer call. Same defect as Coverage F1.
- F2 (no-op) — a file with both a build line and a foreign suffix names the
  build line as the refused constraint. The row and file fields stay correct;
  the edge is undecided and low value.

## Coverage

Count: 5. Worst: F1 — same defect as Spec F1.

- F1 (ask-user) — see Spec F1. A probe with `cmd.Run()` beside
  `t.Run("a case", …)` returned no violation for a stale segment.
- F2 (ask-user) — the census grades tags, never package scope. A citation
  into `testdata/` or a `_`-prefixed directory passes, and no test phase
  compiles those files. The edge is undecided; a scope call.
- F3 (auto-fix) — `checkCitation` joins `base` with the cited path without a
  containment test, so `../outside/foo_test.go` resolves and greens. Refuse a
  `..` segment, an absolute path, and a backslash, the way
  `validatePayloadRows` already does.
- F4 (auto-fix) — the symlink arm of CE7 is decided ("a FIFO or a link reds")
  but untested. Add the link case to the non-regular subtest.
- F5 (no-op) — repeated `-tags` flags in one hand-written argv union instead
  of last-wins. Pathological input from a reviewed manifest.
