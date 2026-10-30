# Read spec and coverage lines from the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/spec/spec.go, internal/spec/header_fence_test.go (new), internal/coverage/coverage.go, internal/coverage/rows.go, internal/coverage/coverage_test.go
Covers: MB24, MB25, MB26, MB29, MB30, MB80, MB83

## What to build

Make `LiveSpecSlugs` (`spec.go:48`), `headerLines` (`spec.go:75`), and
`deriveImplemented` (`spec.go:233`) in `internal/spec` read one fence
classification from the block reader. `metadata` (`spec.go:92`) and
`PersistedFields` (`spec.go:129`) both read `headerLines`, so its migration
covers both. Each one matches its grammar on the raw text of the unfenced lines.
The trim inconsistency between `LiveSpecSlugs` and `headerLines` then closes.
`deriveImplemented` writes its output from the source bytes by the reader's line
offsets, so a CRLF spec flips and keeps each carriage return. `stagedRe`
(`spec.go:38`) keeps its grammar.

Make `coverage.parse` (`rows.go:128`) read the block reader on raw text. A fenced
line is no story number, no `Not covered:` line, and no map row. The historical
marker stays a raw-text match, because the marker is itself a comment. The
reader's H2 heading with the title `User stories` replaces the literal at
`rows.go:136`. The reader's H2 heading also replaces the `## ` prefix test at
`rows.go:140`, so a fenced heading does not end the story list.

The callers keep their calls, because each exported signature stays the same.
These callers need no edit: `internal/commitment/repository/sources.go:64` and
`:85`, `window.go:87`, and `coverage.MapRows` (`rows.go:279`) through
`internal/preflight/ownership/metadata.go`.

The MB25 fixture puts its `Status: implemented` line at column zero inside the
indented fence, because `metadata` reads that label only at column zero. The new
spec rows go in the new file `internal/spec/header_fence_test.go`, because
`internal/spec/spec_test.go` is over the line cap.

## Acceptance

- [ ] A spec with a real `Status: staged` line and a later tilde-fenced `Status: implemented` line reads as staged.
- [ ] An indented fence that holds a spec path and a column-zero status line hides both from `LiveSpecSlugs` and from `Facts`.
- [ ] `Implemented` flips only the real `Status: staged` line when a fenced copy exists, and it keeps each carriage return of a CRLF spec.
- [ ] `TestPersistedFieldsHeader` in `internal/spec/persisted_test.go` stays green without an edit.
- [ ] A fenced map row and a fenced story line add nothing, and a story after a fenced `## ` line counts.
- [ ] The historical marker keeps its exemption.
- [ ] `bench coverage --check` gives the same result on each staged spec at the base and at the tip.
- [ ] `internal/spec` and `internal/coverage` hold no string literal with a run of three backticks, and no literal `## `. The sites today are `spec.go:54`, `spec.go:80`, and `rows.go:140`.
- [ ] `internal/spec/spec.go` and `internal/coverage/rows.go` do not grow.
