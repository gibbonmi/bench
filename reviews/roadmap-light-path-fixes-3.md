# Review pickup: roadmap-light-path-fixes-3

Frozen pair: base `3809acb4`, tip `8e0800f8`. Initial three-axis pass on 2026-09-03. The repair ticket `repair-review-round-one.md` closed the auto-fix findings. The findings below need a reviewer decision.

## Standards

Count: 1 open finding, 1 repair target. Worst: the unwritable-cache fixture harness is pasted twice.

- S1 (ask-user): `unwritableCacheHome` exists at `internal/gate/cache_env_test.go:128` and `internal/testreport/check_test.go:488` with one derivation; the testreport copy also re-derives `restoreMode` and `requireDirectoryWriteDenied`. Rule: AGENTS.md "Code standard — one source per fact" names "a fixture harness pasted N times". The fold needs a shared test-helper package outside the spec fence, so the seam is the reviewer's call.

## Spec

Count: 1 open finding, 1 repair target. Worst: stories 19 and 20 are met only when the environment names a HOME.

- F1 (ask-user, recorded in the spec for veto): story 19 reads "a `gocache.Hold` error to refuse the gate, the lane, and the focused run". All three sites guard on `gocache.Declared`, because an unconditional refusal reds fifteen gate tests whose fixtures declare no HOME. The spec's decision line and rows LQ17 and LQ28 now record the exception.

## Coverage

Count: 4 open findings, 4 repair targets. Worst: `Resolve` treats every `EvalSymlinks` failure as an absent path.

- C1 (ask-user): a self-referential symlink resolves to itself with no error. An absent leaf under a symlinked parent gets one spelling before creation and another after. The spec decided "keeps the absolute spelling when the path does not exist"; the loop case is undecided. Missing: a `canonicalpath_test.go` case with a symlinked parent and a decision on the loop.
- C2 (ask-user): both new AST checks match on the package identifier, so an aliased import of `os/signal` or `path/filepath` escapes them. The tree holds no such alias today (`rg` over cmd and internal). Missing: an aliased-import case per canary family.
- C3 (ask-user): the Sources field-name guard tests for an ASCII space only; a tab or an NBSP in the name yields the unexpected-field message. The spec's boundary row decides "one space" against "none". Missing: a table row per whitespace class.
- C4 (ask-user): `puritycensus` cuts each line at the first `//`. A string literal that holds `//` blinds the rest of the line. The byte-identical copies it replaced had the same cut. Missing: a row for `//` inside a string, and a receiver-independent `Parallel` token.
