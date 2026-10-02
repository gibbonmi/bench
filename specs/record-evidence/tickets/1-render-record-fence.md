# 1. Render the review record fence through one encoder

Blocked by: none
Writes: internal/reviewrecord/files.go, internal/reviewrecord/parse.go, internal/reviewrecord/write.go (new), internal/reviewrecord/record_test.go, internal/reviewrecord/source_test.go, internal/reviewrecord/recordtest/fixture.go, internal/preflight/delegated_evidence_test.go
Covers: RE1, RE2, RE3, RE4, RE5, RE6, RE7, RE8, RE9, RE10, RE102, RE103

## What to build

Chunk: RE-C1.

Make `parse.go` the one owner of fence location. Add a locator beside `fenced` that finds the payload lines of one named fence and reports an unterminated or a duplicate fence. `fenced` returns its payload through that locator, so `parseRecord` and `ReadPlan` keep their answers.

Add `reviewrecord.Render` in the new file `internal/reviewrecord/write.go`. It takes the current document bytes and a record, and it returns the new document bytes. It encodes the record with two-space indentation and without HTML escaping. It does not validate the record. It finds the payload through the locator of `parse.go` and writes no second fence scanner.

- A document with one fence keeps every byte outside the payload lines.
- A document with no fence keeps its bytes, then gets a newline when it lacks one, a blank line, and the fence.
- A nil document becomes `# Review outcomes`, a blank line, and the fence.
- A document with an unterminated or a duplicate fence returns an error.

Move the fence name `bench-review-record` into one constant. The reader in `files.go` and `Render` both use it.

Make `recordtest.Fixture.Save` and the `recordFence` helper in `internal/preflight/delegated_evidence_test.go` render through `Render`. `Save` keeps its current file bytes for a record with no `&`, `<`, or `>`. Keep the hostile literal fences that the spec's Further notes lists.

Put the `Render` tests in `record_test.go`. Put the `Save` row in `source_test.go`, because an internal test cannot import `recordtest`. The package then holds 12 source files. Ticket 2 consumes `Render` for every write.

## Acceptance

- [ ] Prose before and after one fence keeps its bytes after `Render`.
- [ ] A document with no fence and no final newline gets the fence on its own line.
- [ ] A nil document starts with `# Review outcomes`.
- [ ] An excerpt with a tab, a return, U+001B, a newline, and `<&>` reads back byte-exact through `Read`.
- [ ] The payload has two-space indentation and an unescaped `&`.
- [ ] An unterminated fence and a duplicate fence each make `Render` return an error.
- [ ] `Save` of a record whose excerpt holds `&` writes a raw `&`.
- [ ] `TestReviewRecordSource` and `TestDelegatedEvidenceProjection` pass with no assertion changed.
- [ ] The RE7 and RE8 searches give the hits that the spec names.
