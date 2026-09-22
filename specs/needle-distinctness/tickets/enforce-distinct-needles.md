# Enforce distinct anchor needles

Blocked by: none
Writes: internal/anchors/anchor_harness_diagnostics_test.go, specs/needle-distinctness/tickets/enforce-distinct-needles.md (new), reviews/needle-distinctness.md (new)
Covers: none

## What to build

The ordinary anchor tests reject duplicate matching needles for the same file and kind.
The check consumes the complete registry through `Entries()`.
It uses the matcher's existing normalization rules.
Different groups, sections, steps, or diagnostics do not make a repeated needle distinct.
Different files and kinds remain valid.

The approved source is `roadmap/FT326.md` and the batch approval.
The current registry combines 14 families; the roadmap's count of 11 predates the current source.
Sentence boundaries and minimum needle lengths remain outside this ticket.
The coordinator owns the roadmap, changelog, capture records, and landing.

## Acceptance

- [x] The check rejects an exact duplicate from any registered family.
- [x] The check rejects duplicate needles that the existing matcher treats as equal.
- [x] Different files or kinds permit the same needle.
- [x] Different groups, sections, steps, or diagnostics cannot conceal a duplicate.
- [x] The original duplicate-row probe reports the distinctness diagnostic and restores its subject.

## Verification

Run `go test -count=1 -parallel 2 ./internal/anchors`.
Run the original duplicate-row repro through `bench probe`.
Run `bench test --check docs-currency-workflow` against the finished tree.
The coordinator runs the whole-project gate before the landing.

## Debug evidence

The baseline probe duplicates the `edge inventory` row in `generalAnchors`.
Two runs reported `silent` with 154 tests and no skips.
Their package times were 597 ms and 620 ms.
Both runs restored the subject.
The existing registry tests therefore accept the duplicated row.

The ranked hypotheses are missing union validation, per-row validation without distinctness, and validation within only one family.
The guard uses the existing registry invariant seam in `anchor_harness_diagnostics_test.go`.

## Retained author

Author: ft326_author
Line: gpt-6-astra / ultra
Pre-review cap: 3 attempts
Post-review allowance: 2 repair cycles
Expected repair rounds: 1
Confidence: 7
Base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
