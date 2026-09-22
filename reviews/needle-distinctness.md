# Needle distinctness verification

Status: author verified; independent review and landing pending
Base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Assignment: batch-ft326
Author: ft326_author
Line: gpt-6-astra / ultra
Pre-review attempts consumed: 1 of 3
Post-review repair cycles consumed: 0 of 2
Expected repair rounds: 1
Confidence: 7

## Approved boundary

FT326 requires distinct matching needles for each file and kind across the complete registry.
The existing registry combines 14 families, rather than the historical count of 11.
`Entries()` supplies that union without a second family list.
The baseline census found no duplicate matching needles.
The fix adds one shared invariant check beside the existing registry invariant tests.
It changes no production code or registered anchor row.

`Satisfied` normalizes both the needle and the subject through `normalizeMatchMapped`.
The guard uses that same helper.
Its file key uses the evaluator's path semantics.
Its kind key keeps every distinct `Kind` value separate.
Group, section, step, and diagnostic do not distinguish repeated needles.
Sentence boundaries, minimum lengths, and substring overlaps remain outside this change.

## Acceptance evidence

| Row | Behavior | Executed evidence |
|---|---|---|
| ND1 | Exact duplicates fail across the registered union. | `TestRegistryNeedlesDistinct`; the original row-duplication probe changed from silent to bit. |
| ND2 | Matching-equivalent needles fail. | `TestRegistryNeedleDistinctnessBoundary` checks whitespace, scoped case, and folded emphasis through the guard. |
| ND3 | Different files, kinds, or needles remain valid. | The boundary test checks all three; the kind-omission probe failed the different-kind case. |
| ND4 | Group, section, step, and diagnostic cannot conceal a duplicate. | The boundary test requires the distinctness diagnostic for each changed field. |
| ND5 | The original repro bites and restores. | The self-probe reports `bit`, one intended failure, and `restored=yes`. |

Claim rows carry no free-text field.

```text
claims[5]{row,status,confidence}:
  ND1,verified,9
  ND2,verified,9
  ND3,verified,9
  ND4,verified,9
  ND5,verified,9
```

## Debug loop

The repro adds a second copy of the `edge inventory` row to `generalAnchors`.
Both baseline runs reported `silent`, 154 tests, no skips, and a restored subject.
The package times were 597 ms and 620 ms.
The original tests therefore accepted the duplicate.

The hypotheses were missing union validation, per-row validation without distinctness, and validation within only one family.
The new union check confirms the missing validation.
The existing test seam closes the gap without a new runtime owner.

Run every command below through the assignment owner.

```sh
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe internal/anchors/registry_data.go --swap '{File: ".agents/commands/bench-write-spec.md", Kind: Require, Section: "", Needle: "edge inventory", Diagnostic: ".agents/commands/bench-write-spec.md missing acceptance coverage anchor: edge inventory"},' --with '{File: ".agents/commands/bench-write-spec.md", Kind: Require, Section: "", Needle: "edge inventory", Diagnostic: ".agents/commands/bench-write-spec.md missing acceptance coverage anchor: edge inventory"}, {File: ".agents/commands/bench-write-spec.md", Kind: Require, Section: "", Needle: "edge inventory", Diagnostic: ".agents/commands/bench-write-spec.md missing acceptance coverage anchor: edge inventory"},' --package ./internal/anchors
```

The completed guard returned `bit`, 174 tests, one failure, and no skips.
The package time was 637 ms.
The failure named rows 6 and 7, their file, their kind, and `edge inventory`.
The subject restored successfully.

## Independent coordinator probe

The coordinator selected an omission at a different site from the author's row swap.
The retained author executed that probe.

```sh
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe internal/anchors/anchor_harness_diagnostics_test.go --omit 'kind:   anchor.Kind,' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/^different_kind$'
```

The result was `bit`, two executed tests, one failure, and no skips.
The package time was 2 ms.
The different-kind case reported a false duplicate, as the probe requires.
The omitted bytes came from the key initializer on line 222.
The before and after SHA-256 values matched.

Restore SHA-256: 1a431e79c1db0d05694015f4d3d261b35e51d3edba2a5c80e466e0a03586bcaa

## Final focused checks

| Command after `worktree exec batch-ft326 --` | Result | Package time |
|---|---|---|
| `go test -count=1 -parallel 2 ./internal/anchors -run '^TestRegistryNeedlesDistinct$' -v` | Passed; one test executed. | 5 ms |
| `env GOFLAGS='-p=4 -parallel=2' bench test --package ./internal/anchors` | Passed; no failures or skips. | 673 ms |
| `env GOFLAGS='-p=4 -parallel=2' bench test --check docs-currency-workflow` | Passed; no failures or skips. | 501 ms |
| `git diff --check` | Passed. | Not measured |

`bench preflight build needle-distinctness` requires a spec file and cannot grade this tickets-only debug route.
The coordinator confirmed the approved light/debug route after that refusal.
No spec-backed preflight pass is claimed.
The coordinator owns the independent review axes, serial commit, whole-project gate, and landing.
The ticket folder retires through the approved tickets-only landing.

## Fence and follow-up

The final implementation fence contains `internal/anchors/anchor_harness_diagnostics_test.go`, the ticket, and this record.
The temporary repro changes only `internal/anchors/registry_data.go` and restores it.
No existing row required cleanup.
No blocking defect or architecture change remains in the author's scope.

CLI improvement: distinguish the supported tickets-only route from a missing spec in build preflight output.
