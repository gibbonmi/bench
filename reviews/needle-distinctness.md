# Needle distinctness verification

Status: author acceptance satisfied; all native axes pass; landing pending
Base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Assignment: batch-ft326
Author: ft326_author
Line: gpt-6-astra / ultra
Pre-review attempts consumed: 1 of 3
Post-review repair cycles consumed: 1 of 2
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

## Initial native review

Reviewed base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Reviewed source: 09b48bf61dc71066486e4937dc873faf6e6043f6
Review route: separate native Standards, Spec, and Coverage sessions

The coordinator supplied the following native result summaries and dispositions.
These read-only review claims are `claimed`; they are not executed verification.
The original findings remain below when repairs close them.

### Standards

The axis returned two findings, with S1 as the worst issue.
S1 has confidence 9; S2 has confidence 10.

| Finding | Native concern | Binding source | Disposition | State |
|---|---|---|---|---|
| S1 | The 18 independent expectations lack demonstrated named mutation support. | `AGENTS.md` lines 41–47 | auto-fix | Open |
| S2 | The ticket and review duplicate the baseline repro, hypotheses, and timings. | `AGENTS.md` lines 34–39 | auto-fix | Open |

### Spec

The axis returned no findings.
Its conclusion is `claimed` with confidence 9.

### Coverage

The axis returned three findings, with COV-1 as the worst issue.

| Finding | Native concern | Confidence | Disposition | State |
|---|---|---|---|---|
| COV-1 | A guard narrowed from `Entries()` to `generalAnchors` escapes the demonstrated probes and excludes 13 families. | 10 | auto-fix | Open |
| COV-2 | A key that combines `RequireInSection` and `RequireInStep` can reject a permitted kind pair. | 9 | no-op | Refuted |
| COV-3 | The coverage omits NBSP and zero-width-space predicate boundaries from the profile. | 10 | auto-fix | Open |

COV-2 requests all 15 unordered pairs of the six kinds.
The actual key stores the exact `Kind` value without a classification or mapping.
The existing kind-omission probe demonstrates that the independent key field must remain present.
No second kind-specific key path exists.
The coordinator therefore treats additional pair enumeration as optional hardening, rather than a blocking defect.

COV-3 cites `projects/benchkit.md` lines 202–209.
The accepted repair proves that U+00A0 collides after normalization and U+200B remains distinct.
A normalization bypass and a widened zero-width-space predicate must each turn their corresponding expectation red.

The Coverage acceptance claims retain their original confidence values.
The native ND5 abstention carries no confidence under the claim schema.

```text
coverage_claims[6]{row,status,confidence}:
  ND1,claimed,8
  ND2,claimed,6
  ND3,claimed,5
  ND4,claimed,8
  ND5,abstained,
  gate-inclusion,claimed,9
```

### Repair plan

One coherent repair cycle addresses S1, S2, COV-1, and COV-3.
The cycle preserves every approved acceptance row.
The review remains the evidence owner; the ticket cites it.

Each retained independent expectation receives a demonstrated mutation that it must detect.
A cross-family duplicate probe exercises the complete registry union.
The two Unicode boundary cases receive their named mutation probes.
All three native axes must return current results after the repair.

Repair cycles consumed before work: 0 of 2
Raw finding count: 5
Accepted repair targets: 4
Refuted finding count: 1

## Repair cycle 1

Initial review record commit: c042112371bcebf885d97913a9763bc3bc7e9ccb
Repair cycles consumed after verification: 1 of 2

The author committed the initial findings on a green prose lane before the repair.
The repair adds NBSP and zero-width-space cases to the existing boundary table.
The ticket now cites this record instead of duplicating the debug evidence.
No production code, registered row, or existing expectation changed.
The original findings remain in the initial review section.

### Expectation mutation map

Each listed subtest failed under its named mutation.
Every probe started from a green baseline and restored its exact subject.
All valid probes reported no skipped tests.
The six boundary groups cover all 20 independent expectations.

| Probe | Mutation | Failed boundary subtests |
|---|---|---|
| duplicate_rejection | Suppress duplicate rejection with an impossible index condition. | exact; path_alias; different_group; different_section; different_step; different_diagnostic |
| normalization | Use raw needle runes instead of the existing normalizer. | required_whitespace; forbidden_whitespace; required_section_case; forbidden_section_case; step_case; emphasis_case; nonbreaking_space |
| file_key | Omit the file field from the key. | different_file |
| kind_key | Omit the kind field from the key. | different_kind |
| needle_key | Clear the normalized needle value in the key. | different_needle; substring; required_case; forbidden_case |
| zero_width_predicate | Replace U+200B with a space in the key. | zero_width_space |

The first needle-field omission was invalid because it left an unused local variable.
It failed to compile, restored successfully, and supplies no verification evidence.
The replacement clears the value while preserving compilation and produces the four intended behavioral failures.

### Complete registry proof

The cross-family probe adds an `edge inventory` row to `frontDoorAnchors`.
That row duplicates the existing row in `generalAnchors`.
`TestRegistryNeedlesDistinct` reports rows 6 and 460 with the intended diagnostic.
A guard narrowed to `generalAnchors` cannot detect this planted row.
The probe restores `registry_front_door.go` after its observed red.

### Repair probe results

| Probe | Verdict | Failed tests | Executed tests | Package ms | Wall ms |
|---|---|---|---|---|---|
| duplicate_rejection | bit; restored | 6 | 7 | 3 | 9031 |
| normalization | bit; restored | 7 | 8 | 3 | 5506 |
| file_key | bit; restored | 1 | 2 | 2 | 6167 |
| kind_key | bit; restored | 1 | 2 | 2 | 6424 |
| needle_key | bit; restored | 4 | 5 | 3 | 7118 |
| zero_width_predicate | bit; restored | 1 | 2 | 2 | 6305 |
| cross_family | bit; restored | 1 | 1 | 5 | 7014 |

The exact commands below produce the recorded mutations.

```sh
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/anchor_harness_diagnostics_test.go' '--swap' 'if previous, ok := seen[k]; ok {' '--with' 'if previous, ok := seen[k]; ok && i < 0 {' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/(exact|path_alias|different_group|different_section|different_step|different_diagnostic)$'
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/anchor_harness_diagnostics_test.go' '--swap' 'normalized, _ := normalizeMatchMapped(anchor.Kind, runes, identityOrigin(len(runes)))' '--with' 'normalized := runes' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/(required_whitespace|forbidden_whitespace|required_section_case|forbidden_section_case|step_case|emphasis_case|nonbreaking_space)$'
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/anchor_harness_diagnostics_test.go' '--omit' 'file:   filepath.Clean(filepath.FromSlash(anchor.File)),' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/^different_file$'
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/anchor_harness_diagnostics_test.go' '--omit' 'kind:   anchor.Kind,' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/^different_kind$'
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/anchor_harness_diagnostics_test.go' '--swap' 'needle: string(normalized),' '--with' 'needle: string(normalized[:0]),' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/(different_needle|substring|required_case|forbidden_case)$'
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/anchor_harness_diagnostics_test.go' '--swap' 'needle: string(normalized),' '--with' 'needle: strings.ReplaceAll(string(normalized), "\u200b", " "),' --package ./internal/anchors --run '^TestRegistryNeedleDistinctnessBoundary$/^zero_width_space$'
/home/mgibs/workspace/bench/bin/bench.sh worktree exec batch-ft326 -- env GOFLAGS='-p=4 -parallel=2' bench probe 'internal/anchors/registry_front_door.go' '--swap' 'var frontDoorAnchors = []Anchor{' '--with' 'var frontDoorAnchors = []Anchor{
{File: ".agents/commands/bench-write-spec.md", Kind: Require, Needle: "edge inventory", Diagnostic: "planted cross-family duplicate"},' --package ./internal/anchors --run '^TestRegistryNeedlesDistinct$'
```

### Current verification and dispositions

| Target | Repair evidence | Current state |
|---|---|---|
| S1 | The mutation map names every retained expectation and its observed red. | Author verified; Standards reaffirmation pending |
| S2 | The ticket points to this record as the debug evidence owner. | Author verified; Standards reaffirmation pending |
| COV-1 | The cross-family duplicate produces the union diagnostic. | Author verified; Coverage reaffirmation pending |
| COV-3 | NBSP fails under normalization bypass; U+200B fails under predicate widening. | Author verified; Coverage reaffirmation pending |

The full anchor package passes after every probe restores.
The workflow conformance check also passes.
Neither final check reports a failure or skip.

| Final command after `worktree exec batch-ft326 --` | Package ms | Wall ms |
|---|---|---|
| `env GOFLAGS='-p=4 -parallel=2' bench test --package ./internal/anchors` | 665 | 4740 |
| `env GOFLAGS='-p=4 -parallel=2' bench test --check docs-currency-workflow` | 463 | 4550 |

All three native axes must reaffirm the repaired source before the coordinator lands it.

## Terminal native reviews

Frozen base: 216ced3fb95ff10314dc72ab237609f14a6e4bab
Frozen source: 7afcd8bb3156b2a6ff2daf1dec944e2015e02ccc
Repair cycles consumed: 1 of 2
Hardening cycles consumed: 0

The coordinator reports a terminal pass from each distinct native review session.
Standards closes S1 and S2; Coverage closes COV-1 and COV-3 and retains the COV-2 refutation.
Spec confirms all five acceptance rows and the excluded sentence-boundary proposal.
These results supersede the initial occurrences without deleting their findings.

| Axis | Native session | Status | Confidence | Findings |
|---|---|---|---|---|
| Standards | `/root/ft326_standards` | claimed; completed; pass | 10 | 0 |
| Spec | `/root/ft326_spec` | claimed; completed; pass | 9 | 0 |
| Coverage | `/root/ft326_coverage` | claimed; completed; pass | 10 | 0 |

The metadata below uses the terminal Review fields from `internal/reviewrecord`.
The coordinator relayed the excerpts; their SHA-256 values bind the embedded bytes.
Initial excerpts summarize the retained findings; final excerpts preserve the coordinator's supplied native-result text.
Model and effort are the confirmed native settings.

This tickets-only route has no preflight plan or source digest.
Those unavailable fields are omitted instead of fabricated.
The JSON is light-path terminal metadata, not a claim of spec-backed record-schema verification.

```json
[
  {
    "id": "ft326-standards-initial",
    "performer": "/root/ft326_standards",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "findings",
    "native_ref": {
      "ref": "/root/ft326_standards",
      "digest": "sha256:0d697400264363d5dd7ac714c8f9967b1f3df058b8def823d11a95dab4c52f5a",
      "excerpt": "S1: The 18 independent expectations lack demonstrated named mutation support. S2: The ticket and review duplicate the baseline repro, hypotheses, and timings."
    },
    "axis": "Standards",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "09b48bf61dc71066486e4937dc873faf6e6043f6",
    "finding_ids": [
      "S1",
      "S2"
    ],
    "supersedes": []
  },
  {
    "id": "ft326-standards-final",
    "performer": "/root/ft326_standards",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_standards",
      "digest": "sha256:0529d03bc3d99cfd1db1cb3168986c379950cb759808564c77b75af34e30a4d8",
      "excerpt": "S1 all20expectations mapped6probes, S2 singleevidenceownerclosed."
    },
    "axis": "Standards",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "7afcd8bb3156b2a6ff2daf1dec944e2015e02ccc",
    "finding_ids": [],
    "supersedes": [
      "ft326-standards-initial"
    ]
  },
  {
    "id": "ft326-spec-initial",
    "performer": "/root/ft326_spec",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_spec",
      "digest": "sha256:bbd442520ffa7a865b9949ecd01f13457892d2ea136dbb654067e83d317b42ed",
      "excerpt": "The axis returned no findings."
    },
    "axis": "Spec",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "09b48bf61dc71066486e4937dc873faf6e6043f6",
    "finding_ids": [],
    "supersedes": []
  },
  {
    "id": "ft326-spec-final",
    "performer": "/root/ft326_spec",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_spec",
      "digest": "sha256:e501deeebed1a44d252fa512b5e03e8d1e84cad9548bf46f341323c188beb50a",
      "excerpt": "ND1-5preserved, only2Unicodefixture additions; sentenceboundaryexcluded."
    },
    "axis": "Spec",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "7afcd8bb3156b2a6ff2daf1dec944e2015e02ccc",
    "finding_ids": [],
    "supersedes": [
      "ft326-spec-initial"
    ]
  },
  {
    "id": "ft326-coverage-initial",
    "performer": "/root/ft326_coverage",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "findings",
    "native_ref": {
      "ref": "/root/ft326_coverage",
      "digest": "sha256:1f4ca4287f39bff97cc988cd738c19954b71eda6513bfe0e4b17dbe0f1ca1c10",
      "excerpt": "COV-1: Entries narrowed to generalAnchors excludes 13 families. COV-2: Enumerate 15 unordered kind pairs. COV-3: Add NBSP and zero-width-space boundaries."
    },
    "axis": "Coverage",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "09b48bf61dc71066486e4937dc873faf6e6043f6",
    "finding_ids": [
      "COV-1",
      "COV-2",
      "COV-3"
    ],
    "supersedes": []
  },
  {
    "id": "ft326-coverage-final",
    "performer": "/root/ft326_coverage",
    "role": "independent-review",
    "model": "gpt-5.6-sol",
    "effort": "high",
    "state": "completed",
    "outcome": "pass",
    "native_ref": {
      "ref": "/root/ft326_coverage",
      "digest": "sha256:bf31d57805d758c59e94f49051d4d4d333f9cb142a21f1aa13710877de9e94fe",
      "excerpt": "COV1 crossfrontDoor/general rows6/460, COV3 NBSPrawnorm andZWSPwidening closed; COV2no-op exactKind field."
    },
    "axis": "Coverage",
    "base": "216ced3fb95ff10314dc72ab237609f14a6e4bab",
    "tip": "7afcd8bb3156b2a6ff2daf1dec944e2015e02ccc",
    "finding_ids": [],
    "supersedes": [
      "ft326-coverage-initial"
    ]
  }
]
```

## Final author acceptance

The retained author confirms ND1, ND2, ND3, ND4, and ND5 as satisfied.
The acceptance claims retain their previously stated confidence of 9.
The repair verification and exact mutation map above supply the executed evidence.
The final review record changes no owner or test bytes and consumes no repair cycle.
The coordinator owns the remaining whole-project gate and landing.
