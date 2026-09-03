# Review pickup: roadmap-light-path-fixes-3

Frozen pair: base `3809acb4`, tip `8e0800f8`. Initial three-axis pass on 2026-09-03.

## Standards

Count: 4 findings, 3 repair targets. Worst: the unwritable-cache fixture harness is pasted twice.

- S1 (ask-user): `unwritableCacheHome` exists at `internal/gate/cache_env_test.go:128` and `internal/testreport/check_test.go:488` with one derivation; the testreport copy also re-derives `restoreMode` and `requireDirectoryWriteDenied`. Rule: AGENTS.md "Code standard — one source per fact" names "a fixture harness pasted N times". The fold needs a shared test-helper package outside the spec fence, so the seam is the reviewer's call.
- S2 (auto-fix): the four census wrapper comments (`internal/worktree/{landingpolicy,lifecyclepolicy,reclaimpolicy}/purity_census_test.go:9-14`, `internal/canonicalpath/purity_census_test.go:9-14`) enumerate the forbidden imports and ambient effects that `internal/puritycensus/census.go` now owns. Rule: craft-comments "One source owns a fact".
- S3 (no-op): the conformance walk skeleton now has six copies; the shape predates the diff, and a `goSourceFiles` helper is roadmap work.
- S4 (auto-fix, folds into S2): "over this package s own directory" lacks the apostrophe in the same four comments.

## Spec

Count: 5 findings, 3 repair targets. Worst: stories 19 and 20 are met only when the environment names a HOME.

- F1 (ask-user, recorded in the spec for veto): story 19 reads "a `gocache.Hold` error to refuse the gate, the lane, and the focused run". All three sites guard on `gocache.Declared`, because an unconditional refusal reds fifteen gate tests whose fixtures declare no HOME. The build records the exception in the spec's decision line and in rows LQ17 and LQ28.
- F2 (auto-fix): LQ2 and LQ5 name seams that cannot see the row's failure. The closers are `internal/canonicalpath/canonicalpath_test.go` and `TestLinkRelinkStaysGreenInAConsumerRepo`. The rows amend.
- F3 (auto-fix): LQ27 reads "no production file"; the check grades per function, because `internal/gate/subject.go` calls both functions in unrelated work. The row amends to "production function".
- F4 (no-op): `canonicalSubjectRoot` keeps an explicit existence leg, so it is not a one-line wrapper; the leg keeps the caller's refusal posture.
- F5 (auto-fix): `puritycensus.Diagnose` and the `Policy` fields are exported with no caller outside the package.

## Coverage

Count: 4 findings, 4 repair targets. Worst: `Resolve` treats every `EvalSymlinks` failure as an absent path.

- C1 (ask-user): a self-referential symlink resolves to itself with no error. An absent leaf under a symlinked parent gets one spelling before creation and another after. The spec decided "keeps the absolute spelling when the path does not exist"; the loop case is undecided. Missing: a `canonicalpath_test.go` case with a symlinked parent and a decision on the loop.
- C2 (ask-user): both new AST checks match on the package identifier, so an aliased import of `os/signal` or `path/filepath` escapes them. The tree holds no such alias today (`rg` over cmd and internal). Missing: an aliased-import case per canary family.
- C3 (ask-user): the Sources field-name guard tests for an ASCII space only; a tab or an NBSP in the name yields the unexpected-field message. The spec's boundary row decides "one space" against "none". Missing: a table row per whitespace class.
- C4 (ask-user): `puritycensus.Diagnose` cuts each line at the first `//`. A string literal that holds `//` blinds the rest of the line. The byte-identical copies it replaced had the same cut. Missing: a row for `//` inside a string, and a receiver-independent `Parallel` token.
