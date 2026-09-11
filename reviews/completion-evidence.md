# Completion evidence review

## Standards

Current repair review: zero findings; worst issue none. The native reviewer reaffirms source `5457042`.

Historical initial review:

Standards review completed with 2 findings (1 medium, 1 low); worst issue medium.

- S1 (medium, auto-fix): The fixture repeats the canonical axis inventory. Use reviewrecord.Axes().
  Citations at frozen tip 205c929: AGENTS.md:34; AGENTS.md:41; internal/reviewrecord/record.go:100; internal/reviewrecord/recordtest/fixture.go:121.

- S2 (low, auto-fix): Move the changelog entry under the typed Added heading.
  Citations at frozen tip 205c929: CHANGELOG.md:9; .agents/skills/bench-craft-synthesis/SKILL.md.

## Spec

Current repair review: zero findings; worst issue none. The native reviewer reaffirms source `5457042`.

Historical initial review:

completed with findings; Spec axis failed clean review.

- P1 (high, auto-fix): A checkpoint for chunk 2 accepts evidence with no predecessor. Require the ordered planned prefix, including amendment mappings.
  Citations at frozen tip 205c929: specs/completion-evidence/spec.md:81; specs/completion-evidence/spec.md:105; internal/reviewrecord/coverage.go:60.

- P2 (high, auto-fix): Terminal completion accepts an empty source digest or performer. Validate both fields.
  Citations at frozen tip 205c929: docs/adr/0021-benchmark-workflow-orchestration.md:12; specs/completion-evidence/spec.md:204; internal/reviewrecord/parse.go:73.

## Coverage

Current repair review: zero findings; worst issue none. The native reviewer reaffirms source `5457042`.

Historical initial review:

Coverage review completed with 3 findings; worst issue high.

- C1 (high, auto-fix): Control-byte record paths can be accepted when the file exists. Reject newline, return, and tab before traversal and exercise real files.
  Citations at frozen tip 205c929: specs/completion-evidence/spec.md:221; internal/reviewrecord/files.go:23; internal/toon/toon.go:94; internal/reviewrecord/source_test.go:159.

- C2 (medium, auto-fix): Immutable tree reads interpret literal glob characters as Git pathspec syntax. Load both a*.md and ab.md literally.
  Citations at frozen tip 205c929: projects/benchkit.md:178; internal/reviewrecord/plan.go:66; internal/git/tree.go:164.

- C3 (medium, auto-fix): The changed-plan test fails on malformed ticket syntax. Use a valid ticket amendment and prove stale identity plus explicit mapping.
  Citations at frozen tip 205c929: specs/completion-evidence/spec.md:225; internal/reviewrecord/source_test.go:54; internal/reviewrecord/plan.go:86.

## Chunk 1 reconciliation

Spec reaffirms E1, E2, E3, E13, E19, E20, E21, E22, E23, E24, and E38 at the current pair.
All three reviewers used gpt-5.6-sol/high for one repair iteration and executed no tests or probes.
The user approved waiving the extra Claude pass and continuing with the three Sol/high reviews.
Chunk 1 review and repair coverage are closed; chunk 2 may begin.

## Repair verification

The author repaired the confirmed findings and retained their earlier results below.
All three repair follow-ups completed clean against `5457042ec3e0919327dc7943aea7f002a54285b2`.
Both required test suites passed without skips. The missing-axis probe bit and restored successfully.

Coverage confirmed C2 is no-op: the literal ticket fixture passed before any Git reader edit.
The new regression preserves that behavior. The full retained-fixture test passed, including the new clean-result omission.
A separate direct TestRootConformance selection skipped because its required root environment was absent; it supplies no proof.
The registered docs-currency-workflow check passed after the guard repair.

Two Claude review attempts timed out. Automatic approval review rejected the network retry because the destination lacked user authorization.
The user approved the requested waiver. No Claude verdict is claimed.

## Current shared consumer evidence

```toon
blast[263]{changed_symbol,file,line,touched}:
  anchors.generalAnchors,internal/anchors/registry_data.go,26,true
  git.ReadTreeFile,internal/reviewrecord/plan.go,40,true
  git.ReadTreeFile,internal/reviewrecord/plan.go,82,true
  git.TreeWithoutFile,internal/reviewrecord/plan.go,134,true
  preflight.completionEvidenceTable,internal/preflight/review.go,180,true
  preflight.renderReviewPacket,internal/preflight/review.go,127,true
  recordtest.Attach,internal/reviewrecord/recordtest/fixture.go,26,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,26,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,28,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,30,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,56,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,67,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,75,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,76,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,77,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,88,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,96,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,109,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,127,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,139,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,18,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,21,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,73,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,129,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,187,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,219,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,238,true
  recordtest.Fixture.Commit,internal/reviewrecord/recordtest/fixture.go,47,true
  recordtest.Fixture.Commit,internal/reviewrecord/recordtest/fixture.go,114,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,20,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,24,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,48,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,65,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,207,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,221,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,229,true
  recordtest.Fixture.Complete,internal/reviewrecord/source_test.go,22,true
  recordtest.Fixture.Complete,internal/reviewrecord/source_test.go,74,true
  recordtest.Fixture.Evidence,internal/reviewrecord/recordtest/fixture.go,100,true
  recordtest.Fixture.Evidence,internal/reviewrecord/recordtest/fixture.go,122,true
  recordtest.Fixture.Git,internal/reviewrecord/recordtest/fixture.go,75,true
  recordtest.Fixture.Git,internal/reviewrecord/recordtest/fixture.go,76,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,19,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,23,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,47,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,130,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,220,true
  recordtest.Fixture.Tip,internal/reviewrecord/recordtest/fixture.go,112,true
  recordtest.Fixture.Tip,internal/reviewrecord/recordtest/fixture.go,119,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,29,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,34,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,39,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,66,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,222,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,234,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,243,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,247,true
  recordtest.Fixture.Tree,internal/reviewrecord/recordtest/fixture.go,48,true
  recordtest.Fixture.Tree,internal/reviewrecord/recordtest/fixture.go,115,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,29,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,34,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,39,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,42,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,49,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,66,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,208,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,230,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,234,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,243,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,247,true
  recordtest.Fixture.Verification,internal/reviewrecord/recordtest/fixture.go,120,true
  recordtest.Fixture.Verification,internal/reviewrecord/recordtest/fixture.go,135,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,39,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,45,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,46,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,113,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,145,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,64,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,158,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,160,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,163,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,165,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,193,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,203,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,206,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,228,true
  recordtest.Native,internal/reviewrecord/recordtest/fixture.go,93,true
  recordtest.Native,internal/reviewrecord/recordtest/fixture.go,102,true
  recordtest.Native,internal/reviewrecord/source_test.go,46,true
  recordtest.New,internal/reviewrecord/source_test.go,17,true
  recordtest.New,internal/reviewrecord/source_test.go,72,true
  recordtest.New,internal/reviewrecord/source_test.go,128,true
  recordtest.New,internal/reviewrecord/source_test.go,186,true
  recordtest.New,internal/reviewrecord/source_test.go,200,true
  recordtest.New,internal/reviewrecord/source_test.go,218,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,45,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,48,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,52,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,115,true
  recordtest.Spec,internal/reviewrecord/source_test.go,25,true
  recordtest.Spec,internal/reviewrecord/source_test.go,42,true
  recordtest.Spec,internal/reviewrecord/source_test.go,49,true
  recordtest.Spec,internal/reviewrecord/source_test.go,53,true
  recordtest.Spec,internal/reviewrecord/source_test.go,167,true
  recordtest.Spec,internal/reviewrecord/source_test.go,206,true
  recordtest.Spec,internal/reviewrecord/source_test.go,208,true
  recordtest.Spec,internal/reviewrecord/source_test.go,230,true
  reviewrecord.Amendment,internal/reviewrecord/coverage.go,133,true
  reviewrecord.Amendment,internal/reviewrecord/record.go,73,true
  reviewrecord.Amendment,internal/reviewrecord/source_test.go,246,true
  reviewrecord.Axes,internal/preflight/review.go,164,true
  reviewrecord.Axes,internal/reviewrecord/parse.go,57,true
  reviewrecord.Axes,internal/reviewrecord/record.go,77,true
  reviewrecord.Axes,internal/reviewrecord/recordtest/fixture.go,121,true
  reviewrecord.CheckReviews,internal/reviewrecord/coverage.go,73,true
  reviewrecord.CheckReviews,internal/reviewrecord/record_test.go,11,true
  reviewrecord.CheckReviews,internal/reviewrecord/source_test.go,120,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,29,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,34,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,39,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,66,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,234,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,243,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,247,true
  reviewrecord.Chunk,internal/reviewrecord/coverage.go,32,true
  reviewrecord.Chunk,internal/reviewrecord/record.go,71,true
  reviewrecord.Chunk,internal/reviewrecord/record.go,76,true
  reviewrecord.Chunk,internal/reviewrecord/record_test.go,11,true
  reviewrecord.Chunk,internal/reviewrecord/recordtest/fixture.go,119,true
  reviewrecord.Completion,internal/reviewrecord/record.go,72,true
  reviewrecord.Completion,internal/reviewrecord/recordtest/fixture.go,129,true
  reviewrecord.Digest,internal/reviewrecord/parse.go,129,true
  reviewrecord.Digest,internal/reviewrecord/plan.go,108,true
  reviewrecord.Digest,internal/reviewrecord/recordtest/fixture.go,85,true
  reviewrecord.ErrMissing,internal/preflight/review.go,213,true
  reviewrecord.ErrMissing,internal/reviewrecord/files.go,38,true
  reviewrecord.ErrMissing,internal/reviewrecord/parse.go,217,true
  reviewrecord.Evidence,internal/reviewrecord/parse.go,108,true
  reviewrecord.Evidence,internal/reviewrecord/record.go,30,true
  reviewrecord.Evidence,internal/reviewrecord/record.go,37,true
  reviewrecord.Evidence,internal/reviewrecord/recordtest/fixture.go,88,true
  reviewrecord.Evidence,internal/reviewrecord/recordtest/fixture.go,93,true
  reviewrecord.NativeRef,internal/reviewrecord/parse.go,128,true
  reviewrecord.NativeRef,internal/reviewrecord/record.go,20,true
  reviewrecord.NativeRef,internal/reviewrecord/record.go,27,true
  reviewrecord.NativeRef,internal/reviewrecord/recordtest/fixture.go,84,true
  reviewrecord.NativeRef,internal/reviewrecord/recordtest/fixture.go,85,true
  reviewrecord.Parse,internal/reviewrecord/files.go,67,true
  reviewrecord.Parse,internal/reviewrecord/record_test.go,19,true
  reviewrecord.Parse,internal/reviewrecord/source_test.go,76,true
  reviewrecord.Parse,internal/reviewrecord/source_test.go,101,true
  reviewrecord.Parse,internal/reviewrecord/source_test.go,113,true
  reviewrecord.Plan,internal/reviewrecord/coverage.go,116,true
  reviewrecord.Plan,internal/reviewrecord/plan.go,32,true
  reviewrecord.Plan,internal/reviewrecord/plan.go,33,true
  reviewrecord.Plan,internal/reviewrecord/recordtest/fixture.go,23,true
  reviewrecord.Plan,internal/reviewrecord/recordtest/fixture.go,31,true
  reviewrecord.PlannedChunk,internal/reviewrecord/coverage.go,116,true
  reviewrecord.PlannedChunk,internal/reviewrecord/plan.go,27,true
  reviewrecord.PlannedChunk,internal/reviewrecord/recordtest/fixture.go,34,true
  reviewrecord.Probe,internal/reviewrecord/record.go,34,true
  reviewrecord.Probe,internal/reviewrecord/recordtest/fixture.go,102,true
  reviewrecord.Read,internal/preflight/review.go,212,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,25,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,53,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,167,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,194,true
  reviewrecord.ReadPlan,internal/preflight/review.go,210,true
  reviewrecord.ReadPlan,internal/reviewrecord/coverage.go,13,true
  reviewrecord.ReadPlan,internal/reviewrecord/coverage.go,49,true
  reviewrecord.ReadPlan,internal/reviewrecord/recordtest/fixture.go,48,true
  reviewrecord.ReadPlan,internal/reviewrecord/source_test.go,208,true
  reviewrecord.ReadPlan,internal/reviewrecord/source_test.go,230,true
  reviewrecord.Record,internal/reviewrecord/coverage.go,12,true
  reviewrecord.Record,internal/reviewrecord/coverage.go,125,true
  reviewrecord.Record,internal/reviewrecord/files.go,54,true
  reviewrecord.Record,internal/reviewrecord/files.go,57,true
  reviewrecord.Record,internal/reviewrecord/files.go,61,true
  reviewrecord.Record,internal/reviewrecord/files.go,65,true
  reviewrecord.Record,internal/reviewrecord/parse.go,25,true
  reviewrecord.Record,internal/reviewrecord/parse.go,26,true
  reviewrecord.Record,internal/reviewrecord/record_test.go,18,true
  reviewrecord.Record,internal/reviewrecord/recordtest/fixture.go,22,true
  reviewrecord.Record,internal/reviewrecord/recordtest/fixture.go,52,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,81,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,83,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,84,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,85,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,86,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,87,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,88,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,89,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,90,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,91,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,92,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,93,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,97,true
  reviewrecord.RecordPath,internal/preflight/review.go,198,true
  reviewrecord.RecordPath,internal/reviewrecord/files.go,55,true
  reviewrecord.RecordPath,internal/reviewrecord/parse.go,33,true
  reviewrecord.RecordPath,internal/reviewrecord/plan.go,34,true
  reviewrecord.RecordPath,internal/reviewrecord/plan.go,127,true
  reviewrecord.RecordPath,internal/reviewrecord/source_test.go,173,true
  reviewrecord.RecordPath,internal/reviewrecord/source_test.go,182,true
  reviewrecord.Requirement,internal/reviewrecord/plan.go,22,true
  reviewrecord.Requirement,internal/reviewrecord/plan.go,28,true
  reviewrecord.Requirement,internal/reviewrecord/plan.go,112,true
  reviewrecord.Requirement,internal/reviewrecord/recordtest/fixture.go,31,true
  reviewrecord.Requirement,internal/reviewrecord/recordtest/fixture.go,34,true
  reviewrecord.Requirement,internal/reviewrecord/recordtest/fixture.go,96,true
  reviewrecord.Review,internal/reviewrecord/record.go,52,true
  reviewrecord.Review,internal/reviewrecord/record.go,78,true
  reviewrecord.Review,internal/reviewrecord/recordtest/fixture.go,122,true
  reviewrecord.SourceDigest,internal/preflight/review.go,206,true
  reviewrecord.SourceDigest,internal/reviewrecord/coverage.go,20,true
  reviewrecord.SourceDigest,internal/reviewrecord/coverage.go,42,true
  reviewrecord.SourceDigest,internal/reviewrecord/coverage.go,65,true
  reviewrecord.SourceDigest,internal/reviewrecord/recordtest/fixture.go,115,true
  reviewrecord.SourceDigest,internal/reviewrecord/source_test.go,42,true
  reviewrecord.SourceDigest,internal/reviewrecord/source_test.go,49,true
  reviewrecord.Verification,internal/reviewrecord/parse.go,85,true
  reviewrecord.Verification,internal/reviewrecord/record.go,51,true
  reviewrecord.Verification,internal/reviewrecord/record.go,59,true
  reviewrecord.Verification,internal/reviewrecord/recordtest/fixture.go,96,true
  reviewrecord.Verification,internal/reviewrecord/recordtest/fixture.go,97,true
  reviewrecord.Verification,internal/reviewrecord/recordtest/fixture.go,100,true
  reviewrecord.contains,internal/reviewrecord/coverage.go,95,true
  reviewrecord.contains,internal/reviewrecord/parse.go,57,true
  reviewrecord.contains,internal/reviewrecord/plan.go,99,true
  reviewrecord.decode,internal/reviewrecord/parse.go,27,true
  reviewrecord.decode,internal/reviewrecord/plan.go,48,true
  reviewrecord.fenced,internal/reviewrecord/files.go,63,true
  reviewrecord.fenced,internal/reviewrecord/plan.go,44,true
  reviewrecord.findChunk,internal/reviewrecord/coverage.go,56,true
  reviewrecord.findChunk,internal/reviewrecord/coverage.go,81,true
  reviewrecord.mappedIDs,internal/reviewrecord/coverage.go,76,true
  reviewrecord.objectID,internal/reviewrecord/coverage.go,24,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,46,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,46,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,46,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,60,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,60,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,76,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,119,true
  reviewrecord.objectID,internal/reviewrecord/plan.go,37,true
  reviewrecord.objectID,internal/reviewrecord/plan.go,131,true
  reviewrecord.occurrenceState,internal/reviewrecord/parse.go,73,true
  reviewrecord.occurrenceState,internal/reviewrecord/parse.go,113,true
  reviewrecord.readFile,internal/reviewrecord/files.go,59,true
  reviewrecord.requirementsValid,internal/reviewrecord/plan.go,63,true
  reviewrecord.requirementsValid,internal/reviewrecord/plan.go,77,true
  reviewrecord.safeRelative,internal/reviewrecord/files.go,18,true
  reviewrecord.safeRelative,internal/reviewrecord/files.go,29,true
  reviewrecord.safeRelative,internal/reviewrecord/plan.go,67,true
  reviewrecord.sameSet,internal/reviewrecord/coverage.go,57,true
  reviewrecord.uniqueJSON,internal/reviewrecord/parse.go,154,true
  reviewrecord.uniqueJSON,internal/reviewrecord/parse.go,190,true
  reviewrecord.validateEvidence,internal/reviewrecord/parse.go,54,true
  reviewrecord.validateEvidence,internal/reviewrecord/parse.go,87,true
  reviewrecord.validateNative,internal/reviewrecord/parse.go,100,true
  reviewrecord.validateNative,internal/reviewrecord/parse.go,122,true
  reviewrecord.validateVerification,internal/reviewrecord/parse.go,49,true
  reviewrecord.validateVerification,internal/reviewrecord/parse.go,79,true
meta[1]{packages,files,matches,rows,truncated}:
  266,10,66,263,false
citation[1]{sha,state,version,cmd,hash}:
  5457042ec3e0919327dc7943aea7f002a54285b2,clean,0.2.0,bench consumers --changed --base de1447b31903b170679b66bd200431cc15ddc73f --source-tip 5457042ec3e0919327dc7943aea7f002a54285b2 --full,82147b8690fd13f0bfed88d8a7efe1282179faf25d930434318490c8bb12c15a
help[0]{cmd,why}:
```

## Review disposition

The initial reviews reported seven raw findings. Six became repair targets; C2 was refuted by the author’s exact fixture run.
The three current reviews report zero findings. The author also repaired the obsolete clean-result guard found by the full gate.
The three native reviewers used gpt-5.6-sol/high, one iteration, and ran no tests or probes.

The separate Claude opus/high falsification attempt failed with `Request timed out` (exit 1).
This transport failure supplies no review verdict.

## Verification

The author ran the recorded commands against the frozen source. Record validation does not prove semantic judgment correctness.

## Shared consumer evidence

```toon
blast[219]{changed_symbol,file,line,touched}:
  git.ReadTreeFile,internal/reviewrecord/plan.go,40,true
  git.ReadTreeFile,internal/reviewrecord/plan.go,82,true
  git.TreeWithoutFile,internal/reviewrecord/plan.go,134,true
  preflight.completionEvidenceTable,internal/preflight/review.go,180,true
  preflight.renderReviewPacket,internal/preflight/review.go,127,true
  recordtest.Attach,internal/reviewrecord/recordtest/fixture.go,26,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,26,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,28,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,30,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,56,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,67,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,75,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,76,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,77,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,88,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,96,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,109,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,127,true
  recordtest.Fixture,internal/reviewrecord/recordtest/fixture.go,139,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,17,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,20,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,63,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,116,true
  recordtest.Fixture.Commit,internal/reviewrecord/recordtest/fixture.go,47,true
  recordtest.Fixture.Commit,internal/reviewrecord/recordtest/fixture.go,114,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,19,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,23,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,42,true
  recordtest.Fixture.Commit,internal/reviewrecord/source_test.go,55,true
  recordtest.Fixture.Complete,internal/reviewrecord/source_test.go,21,true
  recordtest.Fixture.Evidence,internal/reviewrecord/recordtest/fixture.go,100,true
  recordtest.Fixture.Evidence,internal/reviewrecord/recordtest/fixture.go,122,true
  recordtest.Fixture.Git,internal/reviewrecord/recordtest/fixture.go,75,true
  recordtest.Fixture.Git,internal/reviewrecord/recordtest/fixture.go,76,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,18,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,22,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,41,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,117,true
  recordtest.Fixture.Tip,internal/reviewrecord/recordtest/fixture.go,112,true
  recordtest.Fixture.Tip,internal/reviewrecord/recordtest/fixture.go,119,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,28,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,33,true
  recordtest.Fixture.Tip,internal/reviewrecord/source_test.go,56,true
  recordtest.Fixture.Tree,internal/reviewrecord/recordtest/fixture.go,48,true
  recordtest.Fixture.Tree,internal/reviewrecord/recordtest/fixture.go,115,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,28,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,33,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,36,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,43,true
  recordtest.Fixture.Tree,internal/reviewrecord/source_test.go,56,true
  recordtest.Fixture.Verification,internal/reviewrecord/recordtest/fixture.go,120,true
  recordtest.Fixture.Verification,internal/reviewrecord/recordtest/fixture.go,135,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,39,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,45,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,46,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,113,true
  recordtest.Fixture.Write,internal/reviewrecord/recordtest/fixture.go,145,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,54,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,145,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,147,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,150,true
  recordtest.Fixture.Write,internal/reviewrecord/source_test.go,152,true
  recordtest.Native,internal/reviewrecord/recordtest/fixture.go,93,true
  recordtest.Native,internal/reviewrecord/recordtest/fixture.go,102,true
  recordtest.Native,internal/reviewrecord/source_test.go,40,true
  recordtest.New,internal/reviewrecord/source_test.go,16,true
  recordtest.New,internal/reviewrecord/source_test.go,62,true
  recordtest.New,internal/reviewrecord/source_test.go,115,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,45,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,48,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,52,true
  recordtest.Spec,internal/reviewrecord/recordtest/fixture.go,115,true
  recordtest.Spec,internal/reviewrecord/source_test.go,24,true
  recordtest.Spec,internal/reviewrecord/source_test.go,36,true
  recordtest.Spec,internal/reviewrecord/source_test.go,43,true
  recordtest.Spec,internal/reviewrecord/source_test.go,47,true
  recordtest.Spec,internal/reviewrecord/source_test.go,154,true
  reviewrecord.Amendment,internal/reviewrecord/coverage.go,125,true
  reviewrecord.Amendment,internal/reviewrecord/record.go,73,true
  reviewrecord.Axes,internal/preflight/review.go,164,true
  reviewrecord.Axes,internal/reviewrecord/parse.go,57,true
  reviewrecord.Axes,internal/reviewrecord/record.go,77,true
  reviewrecord.CheckReviews,internal/reviewrecord/coverage.go,73,true
  reviewrecord.CheckReviews,internal/reviewrecord/record_test.go,11,true
  reviewrecord.CheckReviews,internal/reviewrecord/source_test.go,107,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,28,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,33,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,56,true
  reviewrecord.Chunk,internal/reviewrecord/coverage.go,32,true
  reviewrecord.Chunk,internal/reviewrecord/record.go,71,true
  reviewrecord.Chunk,internal/reviewrecord/record.go,76,true
  reviewrecord.Chunk,internal/reviewrecord/record_test.go,11,true
  reviewrecord.Chunk,internal/reviewrecord/recordtest/fixture.go,119,true
  reviewrecord.Completion,internal/reviewrecord/record.go,72,true
  reviewrecord.Completion,internal/reviewrecord/recordtest/fixture.go,129,true
  reviewrecord.Digest,internal/reviewrecord/parse.go,126,true
  reviewrecord.Digest,internal/reviewrecord/plan.go,108,true
  reviewrecord.Digest,internal/reviewrecord/recordtest/fixture.go,85,true
  reviewrecord.ErrMissing,internal/preflight/review.go,213,true
  reviewrecord.ErrMissing,internal/reviewrecord/files.go,37,true
  reviewrecord.ErrMissing,internal/reviewrecord/parse.go,214,true
  reviewrecord.Evidence,internal/reviewrecord/parse.go,105,true
  reviewrecord.Evidence,internal/reviewrecord/record.go,30,true
  reviewrecord.Evidence,internal/reviewrecord/record.go,37,true
  reviewrecord.Evidence,internal/reviewrecord/recordtest/fixture.go,88,true
  reviewrecord.Evidence,internal/reviewrecord/recordtest/fixture.go,93,true
  reviewrecord.NativeRef,internal/reviewrecord/parse.go,125,true
  reviewrecord.NativeRef,internal/reviewrecord/record.go,20,true
  reviewrecord.NativeRef,internal/reviewrecord/record.go,27,true
  reviewrecord.NativeRef,internal/reviewrecord/recordtest/fixture.go,84,true
  reviewrecord.NativeRef,internal/reviewrecord/recordtest/fixture.go,85,true
  reviewrecord.Parse,internal/reviewrecord/files.go,66,true
  reviewrecord.Parse,internal/reviewrecord/record_test.go,19,true
  reviewrecord.Parse,internal/reviewrecord/source_test.go,65,true
  reviewrecord.Parse,internal/reviewrecord/source_test.go,88,true
  reviewrecord.Parse,internal/reviewrecord/source_test.go,100,true
  reviewrecord.Plan,internal/reviewrecord/coverage.go,108,true
  reviewrecord.Plan,internal/reviewrecord/plan.go,32,true
  reviewrecord.Plan,internal/reviewrecord/plan.go,33,true
  reviewrecord.Plan,internal/reviewrecord/recordtest/fixture.go,23,true
  reviewrecord.Plan,internal/reviewrecord/recordtest/fixture.go,31,true
  reviewrecord.PlannedChunk,internal/reviewrecord/coverage.go,108,true
  reviewrecord.PlannedChunk,internal/reviewrecord/plan.go,27,true
  reviewrecord.PlannedChunk,internal/reviewrecord/recordtest/fixture.go,34,true
  reviewrecord.Probe,internal/reviewrecord/record.go,34,true
  reviewrecord.Probe,internal/reviewrecord/recordtest/fixture.go,102,true
  reviewrecord.Read,internal/preflight/review.go,212,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,24,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,47,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,154,true
  reviewrecord.Read,internal/reviewrecord/source_test.go,160,true
  reviewrecord.ReadPlan,internal/preflight/review.go,210,true
  reviewrecord.ReadPlan,internal/reviewrecord/coverage.go,13,true
  reviewrecord.ReadPlan,internal/reviewrecord/coverage.go,49,true
  reviewrecord.ReadPlan,internal/reviewrecord/recordtest/fixture.go,48,true
  reviewrecord.Record,internal/reviewrecord/coverage.go,12,true
  reviewrecord.Record,internal/reviewrecord/coverage.go,117,true
  reviewrecord.Record,internal/reviewrecord/files.go,53,true
  reviewrecord.Record,internal/reviewrecord/files.go,56,true
  reviewrecord.Record,internal/reviewrecord/files.go,60,true
  reviewrecord.Record,internal/reviewrecord/files.go,64,true
  reviewrecord.Record,internal/reviewrecord/parse.go,25,true
  reviewrecord.Record,internal/reviewrecord/parse.go,26,true
  reviewrecord.Record,internal/reviewrecord/record_test.go,18,true
  reviewrecord.Record,internal/reviewrecord/recordtest/fixture.go,22,true
  reviewrecord.Record,internal/reviewrecord/recordtest/fixture.go,52,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,70,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,72,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,73,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,74,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,75,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,76,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,77,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,78,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,79,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,80,true
  reviewrecord.Record,internal/reviewrecord/source_test.go,84,true
  reviewrecord.RecordPath,internal/preflight/review.go,198,true
  reviewrecord.RecordPath,internal/reviewrecord/files.go,54,true
  reviewrecord.RecordPath,internal/reviewrecord/parse.go,33,true
  reviewrecord.RecordPath,internal/reviewrecord/plan.go,34,true
  reviewrecord.RecordPath,internal/reviewrecord/plan.go,127,true
  reviewrecord.Requirement,internal/reviewrecord/plan.go,22,true
  reviewrecord.Requirement,internal/reviewrecord/plan.go,28,true
  reviewrecord.Requirement,internal/reviewrecord/plan.go,112,true
  reviewrecord.Requirement,internal/reviewrecord/recordtest/fixture.go,31,true
  reviewrecord.Requirement,internal/reviewrecord/recordtest/fixture.go,34,true
  reviewrecord.Requirement,internal/reviewrecord/recordtest/fixture.go,96,true
  reviewrecord.Review,internal/reviewrecord/record.go,52,true
  reviewrecord.Review,internal/reviewrecord/record.go,78,true
  reviewrecord.Review,internal/reviewrecord/recordtest/fixture.go,122,true
  reviewrecord.SourceDigest,internal/preflight/review.go,206,true
  reviewrecord.SourceDigest,internal/reviewrecord/coverage.go,20,true
  reviewrecord.SourceDigest,internal/reviewrecord/coverage.go,42,true
  reviewrecord.SourceDigest,internal/reviewrecord/coverage.go,65,true
  reviewrecord.SourceDigest,internal/reviewrecord/recordtest/fixture.go,115,true
  reviewrecord.SourceDigest,internal/reviewrecord/source_test.go,36,true
  reviewrecord.SourceDigest,internal/reviewrecord/source_test.go,43,true
  reviewrecord.Verification,internal/reviewrecord/parse.go,82,true
  reviewrecord.Verification,internal/reviewrecord/record.go,51,true
  reviewrecord.Verification,internal/reviewrecord/record.go,59,true
  reviewrecord.Verification,internal/reviewrecord/recordtest/fixture.go,96,true
  reviewrecord.Verification,internal/reviewrecord/recordtest/fixture.go,97,true
  reviewrecord.Verification,internal/reviewrecord/recordtest/fixture.go,100,true
  reviewrecord.contains,internal/reviewrecord/parse.go,57,true
  reviewrecord.contains,internal/reviewrecord/plan.go,99,true
  reviewrecord.decode,internal/reviewrecord/parse.go,27,true
  reviewrecord.decode,internal/reviewrecord/plan.go,48,true
  reviewrecord.fenced,internal/reviewrecord/files.go,62,true
  reviewrecord.fenced,internal/reviewrecord/plan.go,44,true
  reviewrecord.findChunk,internal/reviewrecord/coverage.go,56,true
  reviewrecord.findChunk,internal/reviewrecord/coverage.go,81,true
  reviewrecord.mappedIDs,internal/reviewrecord/coverage.go,76,true
  reviewrecord.objectID,internal/reviewrecord/coverage.go,24,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,46,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,46,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,46,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,60,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,60,true
  reviewrecord.objectID,internal/reviewrecord/parse.go,116,true
  reviewrecord.objectID,internal/reviewrecord/plan.go,37,true
  reviewrecord.objectID,internal/reviewrecord/plan.go,131,true
  reviewrecord.occurrenceState,internal/reviewrecord/parse.go,73,true
  reviewrecord.occurrenceState,internal/reviewrecord/parse.go,110,true
  reviewrecord.readFile,internal/reviewrecord/files.go,58,true
  reviewrecord.requirementsValid,internal/reviewrecord/plan.go,63,true
  reviewrecord.requirementsValid,internal/reviewrecord/plan.go,77,true
  reviewrecord.safeRelative,internal/reviewrecord/files.go,17,true
  reviewrecord.safeRelative,internal/reviewrecord/files.go,28,true
  reviewrecord.safeRelative,internal/reviewrecord/plan.go,67,true
  reviewrecord.sameSet,internal/reviewrecord/coverage.go,57,true
  reviewrecord.uniqueJSON,internal/reviewrecord/parse.go,151,true
  reviewrecord.uniqueJSON,internal/reviewrecord/parse.go,187,true
  reviewrecord.validateEvidence,internal/reviewrecord/parse.go,54,true
  reviewrecord.validateEvidence,internal/reviewrecord/parse.go,84,true
  reviewrecord.validateNative,internal/reviewrecord/parse.go,97,true
  reviewrecord.validateNative,internal/reviewrecord/parse.go,119,true
  reviewrecord.validateVerification,internal/reviewrecord/parse.go,49,true
  reviewrecord.validateVerification,internal/reviewrecord/parse.go,76,true
meta[1]{packages,files,matches,rows,truncated}:
  266,9,62,219,false
citation[1]{sha,state,version,cmd,hash}:
  205c92928a4bbf35b4540e1f867808c96310d97a,clean,0.2.0,bench consumers --changed --base de1447b31903b170679b66bd200431cc15ddc73f --source-tip 205c92928a4bbf35b4540e1f867808c96310d97a --full,c900cc6f94e0f11d91589c5fae952acdcd155c91b96ba4ec45dd079e21aea067
help[0]{cmd,why}:

```

## Chunk 2 verification and review

Three fresh native Sol/high reviews are pending on the frozen chunk 2 pair.
The author ran the checkpoint, canonical-axis, and public-route requirements against source `e22288d`.
All passed without skips. Both omission probes bit and restored successfully.
The retained plan mapping names unchanged chunk IDs and covers the reviewed verification-inventory amendment.

The shell route now delegates grammar to the gate owner; the installed wrapper still has its pre-landing grammar.
Use this worktree’s source wrapper for the new checkpoint during implementation, preserving the installed publication broker.

Shared consumer evidence is complete. The author walked the untouched command registry,
RunCommand failure/reuse tests, reviewrecord source tests, and preflight reader before dispatch.

```toon
blast[105]{changed_symbol,file,line,touched}:
  bench.commandRegistry,cmd/bench/command_registry.go,245,false
  bench.commandRegistry,cmd/bench/command_registry.go,261,false
  bench.commandRegistry,cmd/bench/command_registry.go,289,false
  bench.commandRegistry,cmd/bench/command_registry_test.go,69,false
  bench.commandRegistry,cmd/bench/help_inventory_test.go,9,true
  bench.commandRegistry,cmd/bench/help_inventory_test.go,10,true
  bench.commandRegistry,cmd/bench/help_inventory_test.go,11,true
  bench.commandRegistry,cmd/bench/otel_hook_seams_test.go,18,false
  bench.gateRoute,cmd/bench/gate_route_test.go,61,true
  bench.gateRoute,cmd/bench/gate_route_test.go,78,true
  bench.gateRoute,cmd/bench/gate_route_test.go,84,true
  gate.Checkpoint,internal/gate/checkpoint.go,24,true
  gate.Checkpoint,internal/gate/checkpoint.go,29,true
  gate.Checkpoint,internal/gate/checkpoint.go,33,true
  gate.Checkpoint,internal/gate/checkpoint.go,49,true
  gate.Checkpoint,internal/gate/checkpoint.go,50,true
  gate.Checkpoint,internal/gate/evaluation.go,26,true
  gate.Checkpoint.validate,internal/gate/checkpoint.go,82,true
  gate.Checkpoint.validate,internal/gate/checkpoint.go,86,true
  gate.Command,cmd/bench/main.go,151,true
  gate.Command,internal/gate/command_test.go,24,true
  gate.CommandUsage,cmd/bench/gate_route_test.go,62,true
  gate.CommandUsage,cmd/bench/gate_route_test.go,93,true
  gate.CommandUsage,cmd/bench/main.go,150,true
  gate.CommandUsage,internal/gate/command_test.go,17,true
  gate.CommandUsage,internal/gate/command_test.go,18,true
  gate.CommandUsage,internal/gate/command_test.go,19,true
  gate.CommandUsage,internal/gate/command_test.go,20,true
  gate.CommandUsage,internal/gate/gate.go,220,true
  gate.CommandUsage,internal/gate/gate.go,245,true
  gate.CommandUsage,internal/gate/gate.go,249,true
  gate.RunCommand,cmd/bench/main.go,152,true
  gate.RunCommand,internal/gate/gate.go,252,true
  gate.RunCommand,internal/gate/review_checkpoint_test.go,27,true
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,72,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,93,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,157,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,185,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,211,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,233,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,254,false
  gate.RunCommand,internal/gate/run_outcomes_test.go,79,false
  gate.RunCommand,internal/gate/run_outcomes_test.go,104,false
  gate.WithCheckpoint,internal/gate/gate.go,231,true
  gate.checkpointEvaluation,internal/gate/gate.go,278,true
  gate.checkpointEvaluation,internal/gate/gate.go,362,true
  gate.checkpointEvaluation,internal/gate/gate.go,366,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,53,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,68,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,91,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,107,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,126,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,156,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,168,true
  gate.checkpointKey,internal/gate/checkpoint.go,25,true
  gate.checkpointKey,internal/gate/checkpoint.go,29,true
  gate.execute,internal/gate/gate.go,267,true
  gate.execute,internal/gate/gate.go,283,true
  gate.executeAfterAcquire,internal/gate/gate.go,233,true
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,85,true
  gate.gateEvaluation,internal/gate/evaluation.go,30,true
  gate.gateEvaluation,internal/gate/evaluation.go,31,true
  gate.gateEvaluation,internal/gate/evaluation.go,49,true
  gate.gateEvaluation,internal/gate/evaluation.go,50,true
  gate.gateEvaluation,internal/gate/evaluation.go,74,true
  gate.gateEvaluation,internal/gate/evaluation.go,88,true
  gate.gateEvaluation,internal/gate/evaluation.go,98,true
  gate.gateEvaluation,internal/gate/evaluation.go,111,true
  gate.gateEvaluation.applyCheckpoint,internal/gate/evaluation.go,122,true
  gate.gateEvaluation.build,internal/gate/evaluation.go,80,true
  gate.gateEvaluation.build,internal/gate/evaluation.go,95,true
  gate.gateEvaluation.build,internal/gate/evaluation.go,104,true
  gate.parseGateArgs,internal/gate/gate.go,218,true
  gate.parseGateArgs,internal/gate/gate.go,248,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,56,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,72,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,78,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,83,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,96,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,116,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,160,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,180,true
  reviewrecord.Check,internal/gate/checkpoint.go,101,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,29,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,34,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,39,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,66,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,234,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,243,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,247,false
  reviewrecord.Read,internal/preflight/review.go,212,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,25,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,53,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,167,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,194,false
  reviewrecord.ReadTree,internal/reviewrecord/check.go,10,true
  reviewrecord.checkCompletion,internal/reviewrecord/coverage.go,123,true
  reviewrecord.checkSource,internal/reviewrecord/check.go,14,true
  reviewrecord.checkSource,internal/reviewrecord/coverage.go,13,true
  reviewrecord.checkVerification,internal/reviewrecord/check.go,56,true
  reviewrecord.checkVerification,internal/reviewrecord/coverage.go,65,true
  reviewrecord.parseRecord,internal/reviewrecord/files.go,64,true
  reviewrecord.parseRecord,internal/reviewrecord/files.go,76,true
blast_deleted[2]{changed_symbol,base_file,base_line}:
  bench.gateUsageLine,cmd/bench/gate_route_test.go,19
  gate.commandUsage,internal/gate/gate.go,244
meta[1]{packages,files,matches,rows,truncated}:
  266,18,41,105,false
citation[1]{sha,state,version,cmd,hash}:
  e22288dae5b1514ebb3f223fb72f12e03390cd24,clean,0.2.0,bench consumers --changed --base 5457042ec3e0919327dc7943aea7f002a54285b2 --source-tip e22288dae5b1514ebb3f223fb72f12e03390cd24 --full,8d749af8f8cc7d0d415c6689f96ac20315561ee888b6531947ebe3f8b4e656e3
help[4]{cmd,why}:
  bench consumers bench.commandRegistry --full,walk the consumers outside the diff
  bench consumers gate.RunCommand --full,walk the consumers outside the diff
  bench consumers reviewrecord.CheckSource --full,walk the consumers outside the diff
  bench consumers reviewrecord.Read --full,walk the consumers outside the diff
```

## Chunk 2 initial review findings

The three axes returned seven findings and six distinct repair targets. All have disposition `auto-fix`.
P5 and C2-C3 share the complete-wrapper repair. The author verified each citation against the frozen source.

## Standards

One finding; worst issue: stale/dead gate-routing ownership. Repair targets: 1.

S1 — auto-fix: Remove unused run_gate and update the stale routing comments. Sources: bin/bench.sh:6-30 and internal/gate/gate.go:6-10. Native terminal conclusion: Standards returned one actionable finding.

## Spec

Spec review completed with 3 medium findings; worst issue: medium. Disposition for each is auto-fix.

P3: Name the last covered/requested chunk in the stale-source refusal (internal/reviewrecord/coverage.go:112-113; spec.md:93). P4: Add the planned verification test/file or amend the approved seam before using the alternative (spec.md:233-240). P5: Add a positive wrapper traversal for --checkpoint … --complete (cmd/bench/gate_route_test.go:68-87). Terminal conclusion: completed with findings.

## Coverage

Coverage review completed with 3 findings; worst issue high. All three dispositions are auto-fix.

C2-C1: Add at least two plan requirements, omit one, and require refusal naming that requirement (recordtest/fixture.go:34; review_checkpoint_test.go:40; E25). C2-C2: Add a valid chunk → complete transition on unchanged source and require another oracle run (review_checkpoint_test.go:67; spec.md:89). C2-C3: Add a valid complete request through the shell using a hostile but valid slug (gate_route_test.go:68; E37). Terminal conclusion: Coverage does not close chunk 2 until C2-C1 through C2-C3 are repaired and re-reviewed.

All results examined `5457042ec3e0919327dc7943aea7f002a54285b2..e22288dae5b1514ebb3f223fb72f12e03390cd24`.
Each performer used gpt-5.6-sol, high effort, and one read-only iteration. No reviewer ran tests or probes.

## Chunk 2 repair verification

The author repaired all six targets in the retained Astra session. All three required checks pass on source `06a589b` without skips.
Both required omission probes bit and restored on that source. The three Sol repair follow-ups are pending.
The current chunk mapping keeps IDs 1, 2, and 3. Ticket r2 adds ownership closure and amends the verification test citations.

Supplemental probes on `721089e` cover the same code before the ticket-only ownership amendment.
Partial-inventory, purpose-collapse, and unquoted-wrapper mutations each failed one named assertion and restored.
These historical probes supplement the current required proofs.

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/reviewrecord/check.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/gate,TestReviewCheckpoint,passed,36
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/partial_additional_verification,"checkpoint accepted invalid evidence or lost reason: exit 0"
```

```text
bit,internal/gate/checkpoint.go,swap,failed,1,yes
package,./internal/gate,TestReviewCheckpointReuse,passed,1
TestReviewCheckpointReuse: chunk green satisfied complete purpose: 2 runs
```

```text
bit,bin/bench.sh,swap,failed,1,yes
package,./cmd/bench,TestGateCheckpointRoute,passed,1
TestGateCheckpointRoute: public chunk checkpoint: 2 usage: bench gate [--fresh] [--checkpoint <spec-path> (--chunk <id> | --complete)]
```

The shared repair charge is complete. The author checked the untouched consumers before dispatch.

```toon
blast[135]{changed_symbol,file,line,touched}:
  bench.commandRegistry,cmd/bench/command_registry.go,245,false
  bench.commandRegistry,cmd/bench/command_registry.go,261,false
  bench.commandRegistry,cmd/bench/command_registry.go,289,false
  bench.commandRegistry,cmd/bench/command_registry_test.go,69,false
  bench.commandRegistry,cmd/bench/help_inventory_test.go,9,true
  bench.commandRegistry,cmd/bench/help_inventory_test.go,10,true
  bench.commandRegistry,cmd/bench/help_inventory_test.go,11,true
  bench.commandRegistry,cmd/bench/otel_hook_seams_test.go,18,false
  bench.gateRoute,cmd/bench/gate_route_test.go,62,true
  bench.gateRoute,cmd/bench/gate_route_test.go,80,true
  bench.gateRoute,cmd/bench/gate_route_test.go,84,true
  bench.gateRoute,cmd/bench/gate_route_test.go,90,true
  gate.Checkpoint,internal/gate/checkpoint.go,24,true
  gate.Checkpoint,internal/gate/checkpoint.go,29,true
  gate.Checkpoint,internal/gate/checkpoint.go,33,true
  gate.Checkpoint,internal/gate/checkpoint.go,49,true
  gate.Checkpoint,internal/gate/checkpoint.go,50,true
  gate.Checkpoint,internal/gate/evaluation.go,26,true
  gate.Checkpoint.validate,internal/gate/checkpoint.go,82,true
  gate.Checkpoint.validate,internal/gate/checkpoint.go,86,true
  gate.Command,cmd/bench/main.go,151,true
  gate.Command,internal/gate/command_test.go,24,true
  gate.CommandUsage,cmd/bench/gate_route_test.go,63,true
  gate.CommandUsage,cmd/bench/gate_route_test.go,99,true
  gate.CommandUsage,cmd/bench/main.go,150,true
  gate.CommandUsage,internal/gate/command_test.go,17,true
  gate.CommandUsage,internal/gate/command_test.go,18,true
  gate.CommandUsage,internal/gate/command_test.go,19,true
  gate.CommandUsage,internal/gate/command_test.go,20,true
  gate.CommandUsage,internal/gate/gate.go,220,true
  gate.CommandUsage,internal/gate/gate.go,245,true
  gate.CommandUsage,internal/gate/gate.go,249,true
  gate.RunCommand,cmd/bench/main.go,152,true
  gate.RunCommand,internal/gate/gate.go,252,true
  gate.RunCommand,internal/gate/review_checkpoint_test.go,27,true
  gate.RunCommand,internal/gate/review_checkpoint_test.go,87,true
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,72,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,93,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,157,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,185,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,211,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,233,false
  gate.RunCommand,internal/gate/run_failure_outcomes_test.go,254,false
  gate.RunCommand,internal/gate/run_outcomes_test.go,79,false
  gate.RunCommand,internal/gate/run_outcomes_test.go,104,false
  gate.WithCheckpoint,internal/gate/gate.go,231,true
  gate.checkpointEvaluation,internal/gate/gate.go,278,true
  gate.checkpointEvaluation,internal/gate/gate.go,362,true
  gate.checkpointEvaluation,internal/gate/gate.go,366,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,55,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,70,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,103,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,119,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,138,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,168,true
  gate.checkpointFixture,internal/gate/review_checkpoint_test.go,180,true
  gate.checkpointKey,internal/gate/checkpoint.go,25,true
  gate.checkpointKey,internal/gate/checkpoint.go,29,true
  gate.execute,internal/gate/gate.go,267,true
  gate.execute,internal/gate/gate.go,283,true
  gate.executeAfterAcquire,internal/gate/gate.go,233,true
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,85,true
  gate.gateEvaluation,internal/gate/evaluation.go,30,true
  gate.gateEvaluation,internal/gate/evaluation.go,31,true
  gate.gateEvaluation,internal/gate/evaluation.go,49,true
  gate.gateEvaluation,internal/gate/evaluation.go,50,true
  gate.gateEvaluation,internal/gate/evaluation.go,74,true
  gate.gateEvaluation,internal/gate/evaluation.go,88,true
  gate.gateEvaluation,internal/gate/evaluation.go,98,true
  gate.gateEvaluation,internal/gate/evaluation.go,111,true
  gate.gateEvaluation.applyCheckpoint,internal/gate/evaluation.go,122,true
  gate.gateEvaluation.build,internal/gate/evaluation.go,80,true
  gate.gateEvaluation.build,internal/gate/evaluation.go,95,true
  gate.gateEvaluation.build,internal/gate/evaluation.go,104,true
  gate.parseGateArgs,internal/gate/gate.go,218,true
  gate.parseGateArgs,internal/gate/gate.go,248,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,58,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,77,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,83,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,95,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,108,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,128,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,172,true
  gate.runCheckpoint,internal/gate/review_checkpoint_test.go,192,true
  recordtest.Attach,internal/gate/review_checkpoint_test.go,17,true
  recordtest.Attach,internal/reviewrecord/recordtest/fixture.go,26,true
  recordtest.AttachAt,cmd/bench/gate_route_test.go,70,true
  recordtest.AttachAt,internal/reviewrecord/recordtest/fixture.go,29,true
  recordtest.Fixture.AddChunk,cmd/bench/gate_route_test.go,76,true
  recordtest.Fixture.AddChunk,internal/gate/review_checkpoint_test.go,18,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,18,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,21,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,73,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,129,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,187,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,219,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,238,false
  recordtest.Fixture.Save,cmd/bench/gate_route_test.go,78,true
  recordtest.Fixture.Save,cmd/bench/gate_route_test.go,89,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,19,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,57,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,72,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,94,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,127,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,146,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,149,true
  recordtest.Fixture.Save,internal/gate/review_checkpoint_test.go,191,true
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,19,false
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,23,false
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,47,false
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,130,false
  recordtest.Fixture.Save,internal/reviewrecord/source_test.go,220,false
  reviewrecord.Check,internal/gate/checkpoint.go,101,true
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,29,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,34,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,39,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,66,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,234,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,243,false
  reviewrecord.CheckSource,internal/reviewrecord/source_test.go,247,false
  reviewrecord.Read,internal/preflight/review.go,212,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,25,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,53,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,167,false
  reviewrecord.Read,internal/reviewrecord/source_test.go,194,false
  reviewrecord.ReadTree,internal/reviewrecord/check.go,10,true
  reviewrecord.checkCompletion,internal/reviewrecord/coverage.go,123,true
  reviewrecord.checkSource,internal/reviewrecord/check.go,14,true
  reviewrecord.checkSource,internal/reviewrecord/coverage.go,13,true
  reviewrecord.checkVerification,internal/reviewrecord/check.go,56,true
  reviewrecord.checkVerification,internal/reviewrecord/coverage.go,65,true
  reviewrecord.parseRecord,internal/reviewrecord/files.go,64,true
  reviewrecord.parseRecord,internal/reviewrecord/files.go,76,true
blast_deleted[2]{changed_symbol,base_file,base_line}:
  bench.gateUsageLine,cmd/bench/gate_route_test.go,19
  gate.commandUsage,internal/gate/gate.go,244
meta[1]{packages,files,matches,rows,truncated}:
  266,19,45,135,false
citation[1]{sha,state,version,cmd,hash}:
  "06a589b935ce26fd9e27d18292f20985c81db487",clean,0.2.0,bench consumers --changed --base 5457042ec3e0919327dc7943aea7f002a54285b2 --source-tip 06a589b935ce26fd9e27d18292f20985c81db487 --full,10d43936604dd34b8c7228c98c7dac47c870b8cd61ee060e36419cc789922405
help[6]{cmd,why}:
  bench consumers bench.commandRegistry --full,walk the consumers outside the diff
  bench consumers gate.RunCommand --full,walk the consumers outside the diff
  bench consumers recordtest.Fixture.AddChunk --full,walk the consumers outside the diff
  bench consumers recordtest.Fixture.Save --full,walk the consumers outside the diff
  bench consumers reviewrecord.CheckSource --full,walk the consumers outside the diff
  bench consumers reviewrecord.Read --full,walk the consumers outside the diff
```

## Chunk 2 repair review closure

All three fresh Sol/high axes returned current clean repair reaffirmations on `06a589b935ce26fd9e27d18292f20985c81db487`.
Standards closed S1. Spec closed P3–P5. Coverage closed C2-C1–C2-C3.

Each axis reports zero findings, no worst issue, and zero repair targets.
The native excerpts remain below. These reviews grant no gate or publication authority.

## Chunk 3 author evidence and review request

The source is `55520083eb588b3d9aa17667d18a06b9c6a51c5c`, based on accepted chunk 2 at `06a589b935ce26fd9e27d18292f20985c81db487`.
The required landing tests passed without skips. The landing-obligation omission bit all 11 refusal cases and restored.
The ordinary lane passed. The full worktree suite passed with two socket capability skips.
The actual dogfood shift and hook results are recorded in the spec.

The final system run failed four existing landing fixtures before their gate barriers.
Those producers need completion plans and retained fixture results. Their race and recovery assertions remain required.
Final acceptance is not complete. Three fresh Sol/high review requests are retained below.

The author walked every untouched consumer in the prepared table before axis dispatch.
Ordinary callers retain the empty obligation. The existing landing tests with injected authorization still test their own publication seams.
The shared fixture and moved census and telemetry declarations keep their consumers. The final whole-suite checks test those edges.

```toon
blast[189]{changed_symbol,file,line,touched}:
  authorization.AuthorizeWithWriters,internal/gate/authorization/authorization.go,145,true
  authorization.AuthorizeWithWriters,internal/gate/authorization/authorization_test.go,17,false
  authorization.AuthorizeWithWriters,internal/landing/completion_evidence_test.go,117,true
  authorization.AuthorizeWithWriters,internal/landing/completion_evidence_test.go,126,true
  authorization.AuthorizeWithWriters,internal/landing/landing.go,120,true
  gate.InspectTree,internal/gate/authorization/authorization.go,172,true
  gate.InspectTree,internal/gate/prospective_owner_test.go,435,false
  gate.InspectTreeContext,internal/gate/authorization/authorization.go,152,true
  gate.InspectTreeContext,internal/gate/engine.go,32,true
  gate.WithCompletion,internal/landing/landing.go,252,true
  gate.WithCompletion,internal/worktree/land_journey_test.go,82,true
  gate.beginGateSpan,internal/gate/engine.go,44,true
  gate.beginGateSpan,internal/gate/gate.go,228,true
  gate.beginGateSpan,internal/gate/gate.go,262,true
  gate.checkpointEvaluation,internal/gate/engine.go,63,true
  gate.checkpointEvaluation,internal/gate/engine.go,95,true
  gate.checkpointEvaluation,internal/gate/gate.go,275,true
  gate.checkpointEvaluation,internal/gate/gate.go,359,true
  gate.checkpointEvaluation,internal/gate/gate.go,363,true
  gate.completionSourceKey,internal/gate/checkpoint.go,30,true
  gate.completionSourceKey,internal/gate/completion.go,19,true
  gate.executeTreeWithOwner,internal/gate/engine.go,46,true
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,292,false
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,310,false
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,343,false
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,362,false
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,86,true
  gate.gateEvaluation,internal/gate/completion.go,22,true
  gate.gateEvaluation,internal/gate/completion.go,68,true
  gate.gateEvaluation,internal/gate/evaluation.go,31,true
  gate.gateEvaluation,internal/gate/evaluation.go,32,true
  gate.gateEvaluation,internal/gate/evaluation.go,50,true
  gate.gateEvaluation,internal/gate/evaluation.go,51,true
  gate.gateEvaluation,internal/gate/evaluation.go,75,true
  gate.gateEvaluation,internal/gate/evaluation.go,89,true
  gate.gateEvaluation,internal/gate/evaluation.go,99,true
  gate.gateEvaluation,internal/gate/evaluation.go,112,true
  gate.gateEvaluation.applyCheckpoint,internal/gate/evaluation.go,123,true
  gate.gateEvaluation.completionTree,internal/gate/checkpoint.go,113,true
  gate.inspectProspective,internal/gate/engine.go,37,true
  gate.operational,internal/gate/run_transaction.go,45,false
  gate.operational,internal/gate/run_transaction.go,57,false
  gate.operational,internal/gate/run_transaction.go,68,false
  gate.operational,internal/gate/run_transaction.go,73,false
  gate.operational,internal/gate/run_transaction.go,94,false
  gate.operational,internal/gate/run_transaction.go,98,false
  gate.operational,internal/gate/run_transaction.go,124,false
  gate.operational,internal/gate/run_transaction.go,142,false
  gate.operational,internal/gate/run_transaction.go,165,false
  gate.operational,internal/gate/run_transaction.go,176,false
  gate.operational,internal/gate/run_transaction.go,181,false
  gate.operational,internal/gate/run_transaction.go,207,false
  gate.otelGateEnv,internal/gate/gate.go,167,true
  gate.otelGateEnv,internal/gate/otel_env_test.go,21,false
  gate.otelGatePhaseSeam,internal/gate/runner.go,572,false
  gate.otelGateSeam,internal/gate/telemetry.go,34,true
  gate.otelGateSeam,internal/gate/telemetry.go,34,true
  gate.otelRecordRootKey,internal/gate/telemetry.go,35,true
  gate.otelRecordRootKey,internal/gate/telemetry.go,53,true
  gate.otelRootEnv,internal/gate/otel_env_test.go,14,false
  gate.otelRootEnv,internal/gate/runner.go,556,false
  gate.otelRootEnv,internal/gate/telemetry.go,24,true
  gate.otelRootEnv,internal/gate/telemetry.go,57,true
  gate.otelTraceparentEnv,internal/gate/otel_env_test.go,15,false
  gate.otelTraceparentEnv,internal/gate/runner.go,562,false
  gate.otelTraceparentEnv,internal/gate/telemetry.go,24,true
  gate.otelTraceparentEnv,internal/gate/telemetry.go,61,true
  gate.validateCompletionContext,internal/gate/checkpoint.go,87,true
  gate.withGateSpanEnv,internal/gate/engine.go,237,true
  landing.Owner.LandReviewed,internal/landing/completion_evidence_test.go,139,true
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,46,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,93,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,145,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,190,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,223,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,282,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,290,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,294,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,320,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,367,false
  landing.Owner.LandReviewed,internal/landing/lane_test.go,75,false
  landing.Owner.LandReviewed,internal/landing/state_test.go,307,false
  landing.Owner.LandReviewed,internal/worktree/joins.go,94,false
  landing.completionFixture,internal/landing/completion_evidence_test.go,22,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,43,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,46,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,68,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,70,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,71,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,75,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,80,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,84,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,89,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,93,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,97,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,101,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,106,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,110,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,120,true
  landing.completionFixture.request,internal/landing/completion_evidence_test.go,138,true
  landing.newCompletionFixture,internal/landing/completion_evidence_test.go,132,true
  recordtest.AttachAt,cmd/bench/gate_route_test.go,70,false
  recordtest.AttachAt,internal/reviewrecord/recordtest/fixture.go,29,true
  recordtest.Fixture.AddChunk,cmd/bench/gate_route_test.go,76,false
  recordtest.Fixture.AddChunk,internal/gate/review_checkpoint_test.go,18,false
  recordtest.Fixture.AddChunk,internal/landing/completion_evidence_test.go,36,true
  recordtest.Fixture.AddChunk,internal/landing/completion_evidence_test.go,39,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,18,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,21,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,73,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,129,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,187,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,219,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,238,false
  recordtest.Fixture.RecordChunk,internal/reviewrecord/recordtest/fixture.go,79,true
  recordtest.Fixture.RecordChunk,internal/reviewrecord/recordtest/fixture.go,149,true
  recordtest.Fixture.loadPlan,internal/reviewrecord/recordtest/fixture.go,57,true
  recordtest.Fixture.loadPlan,internal/reviewrecord/recordtest/fixture.go,75,true
  recordtest.Prepare,internal/reviewrecord/recordtest/fixture.go,33,true
  recordtest.Prepare,internal/worktree/land_fixtures_test.go,91,true
  recordtest.RetainSingleChunk,internal/worktree/completion_fixture_test.go,13,true
  reviewrecord.CheckTrees,internal/gate/checkpoint.go,118,true
  reviewrecord.CheckTrees,internal/reviewrecord/check.go,10,true
  worktree.brokerChangingLanding,internal/worktree/land_effects_test.go,73,true
  worktree.brokerChangingLanding,internal/worktree/land_effects_test.go,95,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_cleanup_test.go,211,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_cleanup_test.go,269,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,226,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,248,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,271,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,308,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,357,true
  worktree.censusRecordPath,internal/worktree/land_census_test.go,49,true
  worktree.censusRecordPath,internal/worktree/land_census_test.go,116,true
  worktree.censusRecordPath,internal/worktree/land_resume_test.go,316,false
  worktree.censusRecordPath,internal/worktree/worktree_test.go,573,false
  worktree.censusRecordPath,internal/worktree/worktree_test.go,650,false
  worktree.foldCompletionComposition,internal/worktree/land_surface_test.go,74,true
  worktree.foldCompletionComposition,internal/worktree/land_surface_test.go,119,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,79,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,110,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,136,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,173,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,212,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,270,true
  worktree.foldLandingSibling,internal/worktree/land_effects_test.go,272,true
  worktree.landingFixtureAtHome,internal/worktree/land_fixtures_test.go,31,true
  worktree.landingFixtureAtHome,internal/worktree/land_fixtures_test.go,40,true
  worktree.landingSpecAmendment,internal/worktree/land_spec_amendment_test.go,39,true
  worktree.landingSpecAmendment,internal/worktree/land_spec_amendment_test.go,64,true
  worktree.recordRawCalls,internal/worktree/land_census_test.go,42,true
  worktree.recordRawCalls,internal/worktree/land_census_test.go,107,true
  worktree.recordRawCalls,internal/worktree/land_resume_test.go,308,false
  worktree.recordRawCalls,internal/worktree/worktree_test.go,566,false
  worktree.recordRawCalls,internal/worktree/worktree_test.go,636,false
  worktree.recordRawCallsWithHead,internal/worktree/land_census_test.go,16,true
  worktree.recordRawCallsWithHead,internal/worktree/land_census_test.go,75,true
  worktree.recordRawCallsWithHead,internal/worktree/land_census_test.go,76,true
  worktree.redProspectiveGateLanding,internal/worktree/land_freshness_test.go,171,true
  worktree.redProspectiveGateLanding,internal/worktree/land_freshness_test.go,210,true
  worktree.refreshLandingEvidence,internal/worktree/completion_fixture_test.go,25,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_cleanup_test.go,33,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_test.go,35,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_test.go,38,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_test.go,172,true
  worktree.refreshLandingEvidence,internal/worktree/land_fixtures_test.go,95,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,36,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,69,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,106,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,142,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,308,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,45,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,129,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,167,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,191,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,220,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,230,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,332,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,382,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,412,true
  worktree.refreshLandingEvidence,internal/worktree/land_spec_amendment_test.go,30,true
  worktree.refreshLandingEvidence,internal/worktree/land_surface_test.go,65,true
  worktree.refreshLandingEvidence,internal/worktree/land_surface_test.go,99,true
  worktree.refreshLandingEvidence,internal/worktree/land_surface_test.go,132,true
  worktree.siblingReviewPath,internal/worktree/land_effects_cleanup_test.go,30,true
  worktree.siblingReviewPath,internal/worktree/land_effects_cleanup_test.go,137,true
  worktree.siblingReviewPath,internal/worktree/land_fixtures_test.go,90,true
meta[1]{packages,files,matches,rows,truncated}:
  266,34,67,189,false
citation[1]{sha,state,version,cmd,hash}:
  55520083eb588b3d9aa17667d18a06b9c6a51c5c,clean,0.2.0,bench consumers --changed --base 06a589b935ce26fd9e27d18292f20985c81db487 --source-tip 55520083eb588b3d9aa17667d18a06b9c6a51c5c --full,64bdeea3a6368cfa737395143af4101433179494f358f33171a33ba4dd694a9d
help[13]{cmd,why}:
  bench consumers authorization.AuthorizeWithWriters --full,walk the consumers outside the diff
  bench consumers gate.InspectTree --full,walk the consumers outside the diff
  bench consumers gate.executeTreeWithOwner --full,walk the consumers outside the diff
  bench consumers gate.operational --full,walk the consumers outside the diff
  bench consumers gate.otelGateEnv --full,walk the consumers outside the diff
  bench consumers gate.otelGatePhaseSeam --full,walk the consumers outside the diff
  bench consumers gate.otelRootEnv --full,walk the consumers outside the diff
  bench consumers gate.otelTraceparentEnv --full,walk the consumers outside the diff
  bench consumers landing.Owner.LandReviewed --full,walk the consumers outside the diff
  bench consumers recordtest.AttachAt --full,walk the consumers outside the diff
  bench consumers recordtest.Fixture.AddChunk --full,walk the consumers outside the diff
  bench consumers worktree.censusRecordPath --full,walk the consumers outside the diff
  bench consumers worktree.recordRawCalls --full,walk the consumers outside the diff
```

## Chunk 3 review findings

All three fresh Sol/high axes completed on `55520083`. Standards returned two findings, Spec two, and Coverage one.
The five raw findings map to three repair targets. Each finding has disposition `auto-fix`.
No axis grants completion or publication authority.

### Standards

S-C3-1 (high) cites the green-tree rule at `.bench/BENCH.md:53` and the four system failures in `codex:tool/87aeb6`.
The shared producer at `internal/systemtest/owner_land_race_test.go:174` lacks completion evidence.
S-C3-2 (medium) cites `AGENTS.md:34` and `internal/gate/tree_snapshot.go:3`: `completion.go:22-46` reads the same entry through two parsers.
Use the existing tree snapshot and one bounded blob reader for both content and mode.

### Spec

P6 (high) shares the system producer target with S-C3-1 and C3-C1.
It cites the fixture migration clause at `spec.md:443` and preserved broker tests at `spec.md:255`.
P7 (high) cites E17 and `spec.md:95,225,445`: `owner_land_race_test.go:65,113` still expects an unreviewed winner delta to publish.
Preserve its race refusal, add the composition refusal, then retain current evidence for the source that includes the winner before retrying.

### Coverage

C3-C1 (high) enumerates the four system consumers of the shared producer.
They are the public race journey and the killed, concurrent, and dead-plus-live artifact recovery journeys.
Expand the approved fence before migrating this producer through `recordtest`.
The axis refuted its other E15–E18 and E33–E36 candidates against existing tests.

The Coverage reviewer proposes adding command-boundary system consumers to the consumer census.
The current function census did not expose these binary-driven journeys.

## Chunk 3 repair closure and final acceptance

The retained Astra author repaired the three targets and retained current results.
Standards, Spec, and Coverage each returned zero findings on `943a373b`.
The coordinator checked each review venue and each cited repair before acceptance.
All 38 acceptance rows are covered. Record validation proves occurrence and source
coverage; it does not prove judgment correctness or authenticate an invented result.

The final ordinary package suite and the separate full system suite passed.
The required landing omission probe bit all 11 refusal cases and restored.
The regular-file, executable-mode, and byte-bound probes also bit and restored.
Those supplemental probes examined the same production bytes before the final plan amendment.

The changed-package selector refused a changed system-only package on `d8bd938b`.
The approved expansion selected every ordinary package and retained the separate system suite.
The refused result remains in the record. No test or pass criterion was removed.
Chunk IDs remain 1, 2, and 3. Ticket r3 adds E20 to the final chunk.

| Acceptance rows | Final evidence |
| --- | --- |
| E1–E3, E13, E19–E21, E23–E24, E38 | Record parser, source-chain, hostile-path, amendment, and supersession tests passed. |
| E4–E12, E14, E25–E32, E37 | Checkpoint, canonical-axis, ordinary-work, reuse, and wrapper tests passed. |
| E15–E18, E33–E36 | Landing completion cases and the full system race and recovery journeys passed. |
| E22 | Spec review confirmed the documented trusted-harness limit. |

The current shared charge retains the same untouched-consumer inventory as `d8bd938b`.
The coordinator walked every untouched row. Existing ordinary callers retain an empty
obligation, and existing injected-authority tests retain their publication assertions.

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/gate/completion.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/gate,TestCompletionFile,passed,6
packages[1]{package,status,elapsed_ms}:
  github.com/gibbonmi/bench/internal/gate,fail,279
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/gate,TestCompletionFile/symlink,"completion_test.go:42: symlink completion file admitted: 10 bytes"
skips[0]{package,test,reason}:

```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/gate/completion.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/gate,TestCompletionFile,passed,6
packages[1]{package,status,elapsed_ms}:
  github.com/gibbonmi/bench/internal/gate,fail,298
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/gate,TestCompletionFile/executable,"completion_test.go:51: completion file = {data:[115 111 117 114 99 101 10] mode:100644}, <nil>"
skips[0]{package,test,reason}:

```

```text
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/git/tree.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/gate,TestCompletionFile,passed,6
packages[1]{package,status,elapsed_ms}:
  github.com/gibbonmi/bench/internal/gate,fail,301
failures[1]{package,test,line}:
  github.com/gibbonmi/bench/internal/gate,TestCompletionFile/oversized,"completion_test.go:42: oversized completion file admitted: 2097153 bytes"
skips[0]{package,test,reason}:

```

```toon
blast[201]{changed_symbol,file,line,touched}:
  authorization.AuthorizeWithWriters,internal/gate/authorization/authorization.go,145,true
  authorization.AuthorizeWithWriters,internal/gate/authorization/authorization_test.go,17,false
  authorization.AuthorizeWithWriters,internal/landing/completion_evidence_test.go,117,true
  authorization.AuthorizeWithWriters,internal/landing/completion_evidence_test.go,126,true
  authorization.AuthorizeWithWriters,internal/landing/landing.go,120,true
  gate.InspectTree,internal/gate/authorization/authorization.go,172,true
  gate.InspectTree,internal/gate/prospective_owner_test.go,435,false
  gate.InspectTreeContext,internal/gate/authorization/authorization.go,152,true
  gate.InspectTreeContext,internal/gate/engine.go,32,true
  gate.WithCompletion,internal/landing/landing.go,252,true
  gate.WithCompletion,internal/worktree/land_journey_test.go,82,true
  gate.beginGateSpan,internal/gate/engine.go,44,true
  gate.beginGateSpan,internal/gate/gate.go,228,true
  gate.beginGateSpan,internal/gate/gate.go,262,true
  gate.checkpointEvaluation,internal/gate/engine.go,63,true
  gate.checkpointEvaluation,internal/gate/engine.go,95,true
  gate.checkpointEvaluation,internal/gate/gate.go,275,true
  gate.checkpointEvaluation,internal/gate/gate.go,359,true
  gate.checkpointEvaluation,internal/gate/gate.go,363,true
  gate.completionFile,internal/gate/completion.go,79,true
  gate.completionFile,internal/gate/completion.go,83,true
  gate.completionFile,internal/gate/completion.go,86,true
  gate.completionSourceKey,internal/gate/checkpoint.go,30,true
  gate.completionSourceKey,internal/gate/completion.go,19,true
  gate.executeTreeWithOwner,internal/gate/engine.go,46,true
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,292,false
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,310,false
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,343,false
  gate.executeTreeWithOwner,internal/gate/prospective_owner_test.go,362,false
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,28,true
  gate.gateEvaluation,internal/gate/checkpoint.go,86,true
  gate.gateEvaluation,internal/gate/completion.go,22,true
  gate.gateEvaluation,internal/gate/completion.go,66,true
  gate.gateEvaluation,internal/gate/completion.go,79,true
  gate.gateEvaluation,internal/gate/evaluation.go,31,true
  gate.gateEvaluation,internal/gate/evaluation.go,32,true
  gate.gateEvaluation,internal/gate/evaluation.go,50,true
  gate.gateEvaluation,internal/gate/evaluation.go,51,true
  gate.gateEvaluation,internal/gate/evaluation.go,75,true
  gate.gateEvaluation,internal/gate/evaluation.go,89,true
  gate.gateEvaluation,internal/gate/evaluation.go,99,true
  gate.gateEvaluation,internal/gate/evaluation.go,112,true
  gate.gateEvaluation.applyCheckpoint,internal/gate/evaluation.go,123,true
  gate.gateEvaluation.completionFile,internal/gate/completion.go,32,true
  gate.gateEvaluation.completionFile,internal/gate/completion.go,40,true
  gate.gateEvaluation.completionFile,internal/gate/completion_test.go,39,true
  gate.gateEvaluation.completionTree,internal/gate/checkpoint.go,113,true
  gate.inspectProspective,internal/gate/engine.go,37,true
  gate.operational,internal/gate/run_transaction.go,45,false
  gate.operational,internal/gate/run_transaction.go,57,false
  gate.operational,internal/gate/run_transaction.go,68,false
  gate.operational,internal/gate/run_transaction.go,73,false
  gate.operational,internal/gate/run_transaction.go,94,false
  gate.operational,internal/gate/run_transaction.go,98,false
  gate.operational,internal/gate/run_transaction.go,124,false
  gate.operational,internal/gate/run_transaction.go,142,false
  gate.operational,internal/gate/run_transaction.go,165,false
  gate.operational,internal/gate/run_transaction.go,176,false
  gate.operational,internal/gate/run_transaction.go,181,false
  gate.operational,internal/gate/run_transaction.go,207,false
  gate.otelGateEnv,internal/gate/gate.go,167,true
  gate.otelGateEnv,internal/gate/otel_env_test.go,21,false
  gate.otelGatePhaseSeam,internal/gate/runner.go,572,false
  gate.otelGateSeam,internal/gate/telemetry.go,34,true
  gate.otelGateSeam,internal/gate/telemetry.go,34,true
  gate.otelRecordRootKey,internal/gate/telemetry.go,35,true
  gate.otelRecordRootKey,internal/gate/telemetry.go,53,true
  gate.otelRootEnv,internal/gate/otel_env_test.go,14,false
  gate.otelRootEnv,internal/gate/runner.go,556,false
  gate.otelRootEnv,internal/gate/telemetry.go,24,true
  gate.otelRootEnv,internal/gate/telemetry.go,57,true
  gate.otelTraceparentEnv,internal/gate/otel_env_test.go,15,false
  gate.otelTraceparentEnv,internal/gate/runner.go,562,false
  gate.otelTraceparentEnv,internal/gate/telemetry.go,24,true
  gate.otelTraceparentEnv,internal/gate/telemetry.go,61,true
  gate.validateCompletionContext,internal/gate/checkpoint.go,87,true
  gate.withGateSpanEnv,internal/gate/engine.go,237,true
  git.ReadControlBlob,internal/gate/completion.go,85,true
  git.ReadControlBlob,internal/git/tree.go,174,true
  git.ReadTreeFile,internal/reviewrecord/files.go,72,false
  git.ReadTreeFile,internal/reviewrecord/plan.go,40,false
  git.ReadTreeFile,internal/reviewrecord/plan.go,82,false
  landing.Owner.LandReviewed,internal/landing/completion_evidence_test.go,139,true
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,46,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,93,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,145,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,190,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,223,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,282,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,290,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,294,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,320,false
  landing.Owner.LandReviewed,internal/landing/landing_reviewed_test.go,367,false
  landing.Owner.LandReviewed,internal/landing/lane_test.go,75,false
  landing.Owner.LandReviewed,internal/landing/state_test.go,307,false
  landing.Owner.LandReviewed,internal/worktree/joins.go,94,false
  landing.completionFixture,internal/landing/completion_evidence_test.go,22,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,43,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,46,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,68,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,70,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,71,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,75,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,80,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,84,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,89,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,93,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,97,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,101,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,106,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,110,true
  landing.completionFixture,internal/landing/completion_evidence_test.go,120,true
  landing.completionFixture.request,internal/landing/completion_evidence_test.go,138,true
  landing.newCompletionFixture,internal/landing/completion_evidence_test.go,132,true
  recordtest.AttachAt,cmd/bench/gate_route_test.go,70,false
  recordtest.AttachAt,internal/reviewrecord/recordtest/fixture.go,29,true
  recordtest.Fixture.AddChunk,cmd/bench/gate_route_test.go,76,false
  recordtest.Fixture.AddChunk,internal/gate/review_checkpoint_test.go,18,false
  recordtest.Fixture.AddChunk,internal/landing/completion_evidence_test.go,36,true
  recordtest.Fixture.AddChunk,internal/landing/completion_evidence_test.go,39,true
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,18,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,21,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,73,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,129,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,187,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,219,false
  recordtest.Fixture.AddChunk,internal/reviewrecord/source_test.go,238,false
  recordtest.Fixture.RecordChunk,internal/reviewrecord/recordtest/fixture.go,79,true
  recordtest.Fixture.RecordChunk,internal/reviewrecord/recordtest/fixture.go,149,true
  recordtest.Fixture.loadPlan,internal/reviewrecord/recordtest/fixture.go,57,true
  recordtest.Fixture.loadPlan,internal/reviewrecord/recordtest/fixture.go,75,true
  recordtest.Prepare,internal/reviewrecord/recordtest/fixture.go,33,true
  recordtest.Prepare,internal/worktree/land_fixtures_test.go,91,true
  recordtest.RetainSingleChunk,internal/worktree/completion_fixture_test.go,13,true
  reviewrecord.CheckTrees,internal/gate/checkpoint.go,118,true
  reviewrecord.CheckTrees,internal/reviewrecord/check.go,10,true
  worktree.brokerChangingLanding,internal/worktree/land_effects_test.go,73,true
  worktree.brokerChangingLanding,internal/worktree/land_effects_test.go,95,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_cleanup_test.go,211,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_cleanup_test.go,269,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,226,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,248,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,271,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,308,true
  worktree.brokerDestinationFixture,internal/worktree/land_effects_test.go,357,true
  worktree.censusRecordPath,internal/worktree/land_census_test.go,49,true
  worktree.censusRecordPath,internal/worktree/land_census_test.go,116,true
  worktree.censusRecordPath,internal/worktree/land_resume_test.go,316,false
  worktree.censusRecordPath,internal/worktree/worktree_test.go,573,false
  worktree.censusRecordPath,internal/worktree/worktree_test.go,650,false
  worktree.foldCompletionComposition,internal/worktree/land_surface_test.go,74,true
  worktree.foldCompletionComposition,internal/worktree/land_surface_test.go,119,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,79,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,110,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,136,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,173,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,212,true
  worktree.foldLandingSibling,internal/worktree/land_effects_cleanup_test.go,270,true
  worktree.foldLandingSibling,internal/worktree/land_effects_test.go,272,true
  worktree.landingFixtureAtHome,internal/worktree/land_fixtures_test.go,31,true
  worktree.landingFixtureAtHome,internal/worktree/land_fixtures_test.go,40,true
  worktree.landingSpecAmendment,internal/worktree/land_spec_amendment_test.go,39,true
  worktree.landingSpecAmendment,internal/worktree/land_spec_amendment_test.go,64,true
  worktree.recordRawCalls,internal/worktree/land_census_test.go,42,true
  worktree.recordRawCalls,internal/worktree/land_census_test.go,107,true
  worktree.recordRawCalls,internal/worktree/land_resume_test.go,308,false
  worktree.recordRawCalls,internal/worktree/worktree_test.go,566,false
  worktree.recordRawCalls,internal/worktree/worktree_test.go,636,false
  worktree.recordRawCallsWithHead,internal/worktree/land_census_test.go,16,true
  worktree.recordRawCallsWithHead,internal/worktree/land_census_test.go,75,true
  worktree.recordRawCallsWithHead,internal/worktree/land_census_test.go,76,true
  worktree.redProspectiveGateLanding,internal/worktree/land_freshness_test.go,171,true
  worktree.redProspectiveGateLanding,internal/worktree/land_freshness_test.go,210,true
  worktree.refreshLandingEvidence,internal/worktree/completion_fixture_test.go,25,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_cleanup_test.go,33,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_test.go,35,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_test.go,38,true
  worktree.refreshLandingEvidence,internal/worktree/land_effects_test.go,172,true
  worktree.refreshLandingEvidence,internal/worktree/land_fixtures_test.go,95,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,36,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,69,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,106,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,142,true
  worktree.refreshLandingEvidence,internal/worktree/land_freshness_test.go,308,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,45,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,129,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,167,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,191,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,220,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,230,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,332,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,382,true
  worktree.refreshLandingEvidence,internal/worktree/land_journey_test.go,412,true
  worktree.refreshLandingEvidence,internal/worktree/land_spec_amendment_test.go,30,true
  worktree.refreshLandingEvidence,internal/worktree/land_surface_test.go,65,true
  worktree.refreshLandingEvidence,internal/worktree/land_surface_test.go,99,true
  worktree.refreshLandingEvidence,internal/worktree/land_surface_test.go,132,true
  worktree.siblingReviewPath,internal/worktree/land_effects_cleanup_test.go,30,true
  worktree.siblingReviewPath,internal/worktree/land_effects_cleanup_test.go,137,true
  worktree.siblingReviewPath,internal/worktree/land_fixtures_test.go,90,true
blast_deleted[6]{changed_symbol,base_file,base_line}:
  systemtest.TestProspectiveArtifactRecoveryAfterKilledLanding,internal/systemtest/owner_artifact_recovery_test.go,20
  systemtest.TestWorktreeLandPublicRaceAndRerun,internal/systemtest/owner_land_race_test.go,16
  systemtest.configureArtifactLandingFixture,internal/systemtest/owner_artifact_recovery_test.go,108
  systemtest.startArtifactAuthorization,internal/systemtest/owner_artifact_recovery_test.go,440
  systemtest.systemLandingRaceFixture,internal/systemtest/owner_land_race_test.go,174
  systemtest.systemTicketBody,internal/systemtest/owner_land_race_test.go,167
meta[1]{packages,files,matches,rows,truncated}:
  266,38,72,201,false
citation[1]{sha,state,version,cmd,hash}:
  943a373b3e4b714b7074da233f7ad1c4b7a6f1e5,clean,0.2.0,bench consumers --changed --base 06a589b935ce26fd9e27d18292f20985c81db487 --source-tip 943a373b3e4b714b7074da233f7ad1c4b7a6f1e5 --full,cf570a8740b91c2fdb47f1b6f8a0605ed21e24614de81ee1bbd98dad20459c33
help[14]{cmd,why}:
  bench consumers authorization.AuthorizeWithWriters --full,walk the consumers outside the diff
  bench consumers gate.InspectTree --full,walk the consumers outside the diff
  bench consumers gate.executeTreeWithOwner --full,walk the consumers outside the diff
  bench consumers gate.operational --full,walk the consumers outside the diff
  bench consumers gate.otelGateEnv --full,walk the consumers outside the diff
  bench consumers gate.otelGatePhaseSeam --full,walk the consumers outside the diff
  bench consumers gate.otelRootEnv --full,walk the consumers outside the diff
  bench consumers gate.otelTraceparentEnv --full,walk the consumers outside the diff
  bench consumers git.ReadTreeFile --full,walk the consumers outside the diff
  bench consumers landing.Owner.LandReviewed --full,walk the consumers outside the diff
  bench consumers recordtest.AttachAt --full,walk the consumers outside the diff
  bench consumers recordtest.Fixture.AddChunk --full,walk the consumers outside the diff
  bench consumers worktree.censusRecordPath --full,walk the consumers outside the diff
  bench consumers worktree.recordRawCalls --full,walk the consumers outside the diff
```

## Machine record

```bench-review-record
{
  "version": 1,
  "spec": "specs/completion-evidence/spec.md",
  "plan_digest": "sha256:e298566dc9a9abc5e580968cae0bf69057dfe06abb8e59feb3fb286948d0c754",
  "implementation_session": "/root",
  "chunks": [
    {
      "id": "1",
      "base": "de1447b31903b170679b66bd200431cc15ddc73f",
      "tip": "5457042ec3e0919327dc7943aea7f002a54285b2",
      "plan_digest": "sha256:267c4a827157b8e6acb6a9a24a56b984305fdea06ed5601964bef8e434df0f05",
      "source_digest": "6eb528d5bbc0b629a82fcd2cb19fa9ea4e69799d",
      "acceptance_rows": [
        "E1",
        "E2",
        "E3",
        "E13",
        "E19",
        "E20",
        "E21",
        "E22",
        "E23",
        "E24",
        "E38"
      ],
      "verification": [
        {
          "id": "c1-record-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:tool/639c7c",
            "digest": "sha256:dd115b4aa2866d0c697be74efa32f14ec691b9f555c2d28b43d9794553556dc2",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,1233\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "record-tests",
          "command": "bench test --package ./internal/reviewrecord --run TestReviewRecord",
          "exit_code": 0,
          "probe": {
            "mutation": "omit missing-axis rejection",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:probe/chunk1-missing-axis",
              "digest": "sha256:e88573124ae65a3931462550fc350931a8856c4fe234570a444db27d13e123f4",
              "excerpt": "probe: bit; subject: internal/reviewrecord/record.go; mutation: swap; failed_tests: 1; restored: yes. TestReviewRecord/missing_clean_axis: missing clean review result: <nil>."
            }
          }
        },
        {
          "id": "c1-preflight-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:tool/721367",
            "digest": "sha256:bedecc25f4b65b146aaba8ffe6b645fc94367d1e7908225138b28bc443c6f796",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/preflight,pass,10399\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "preflight-tests",
          "command": "bench test --package ./internal/preflight --run TestReviewCharge",
          "exit_code": 0
        },
        {
          "id": "c1r3-record-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "6eb528d5bbc0b629a82fcd2cb19fa9ea4e69799d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:tool/0bcbdd",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2048\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:02e2ff4e74cc2fe949de205b949d0a1676fc02f9c2829e7ca512a7eae58386c3"
          },
          "requirement": "record-tests",
          "command": "bench test --package ./internal/reviewrecord --run TestReviewRecord",
          "exit_code": 0,
          "probe": {
            "mutation": "omit missing-axis rejection",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/94b1e8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestReviewRecord,passed,29\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/reviewrecord,fail,1823\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/reviewrecord,TestReviewRecord/missing_clean_axis,\"record_test.go:12: missing clean review result: <nil>\"\nskips[0]{package,test,reason}:\n",
              "digest": "sha256:bbc5f78b2f519348986e8d6d0d1f4524f263b37c76f98e4e53a7c0bbc5df3ca4"
            }
          }
        },
        {
          "id": "c1r3-preflight-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "6eb528d5bbc0b629a82fcd2cb19fa9ea4e69799d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:tool/2975c6",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/preflight,pass,10807\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:00c0cb76546d4f748d08d29f2b30392ceb85af9a21b59c8a9c2090ba2d8b12d5"
          },
          "requirement": "preflight-tests",
          "command": "bench test --package ./internal/preflight --run TestReviewCharge",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c1-standards",
          "performer": "/root/c1_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "codex:agent/c1_standards",
            "excerpt": "Standards review completed with 2 findings (1 medium, 1 low); worst issue medium.",
            "digest": "sha256:0cc67e8ec6ebcee965aaedbae50ba5caaf0281d7fc3e2eb7df136157865281aa"
          },
          "axis": "Standards",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "c1-spec",
          "performer": "/root/c1_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "codex:agent/c1_spec",
            "excerpt": "completed with findings; Spec axis failed clean review.",
            "digest": "sha256:2fcd632e4be4bd0c75889eed7260f518352f7540532fe9933f4dc15b67ef8c14"
          },
          "axis": "Spec",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "c1-coverage",
          "performer": "/root/c1_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "codex:agent/c1_coverage",
            "excerpt": "Coverage review completed with 3 findings; worst issue high.",
            "digest": "sha256:deaea2bf973fd416e6ecdc37ac80dbbd1827fe3c9e736bdd5adc405b71f4eb69"
          },
          "axis": "Coverage",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "c1r3-standards",
          "performer": "/root/c1_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "6eb528d5bbc0b629a82fcd2cb19fa9ea4e69799d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:agent/c1_standards/repair-followup",
            "excerpt": "Standards repair review completed with no findings; S1 and S2 are closed at the repaired source.",
            "digest": "sha256:c6f3c4277709be4c4e23f9f8b4031907a35922a60cd897e67a50dd16e439e54a"
          },
          "axis": "Standards",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "finding_ids": [],
          "supersedes": [
            "c1-standards"
          ]
        },
        {
          "id": "c1r3-spec",
          "performer": "/root/c1_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "6eb528d5bbc0b629a82fcd2cb19fa9ea4e69799d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:agent/c1_spec/repair-followup",
            "excerpt": "completed clean; Spec axis reaffirms the repaired current pair",
            "digest": "sha256:f10817ce3a90da4b216c7457a0cda693d06d4e3a947a01d71293ae943899f84a"
          },
          "axis": "Spec",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "finding_ids": [],
          "supersedes": [
            "c1-spec"
          ]
        },
        {
          "id": "c1r3-coverage",
          "performer": "/root/c1_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "6eb528d5bbc0b629a82fcd2cb19fa9ea4e69799d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:agent/c1_coverage/repair-followup",
            "excerpt": "Coverage completed with zero findings; C1 and C3 are resolved, and C2 is no-op.",
            "digest": "sha256:41258191ea65766980d34a08c91ed198e6bbd5df9076212c830ae327034ffe4a"
          },
          "axis": "Coverage",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "finding_ids": [],
          "supersedes": [
            "c1-coverage"
          ]
        }
      ]
    },
    {
      "id": "2",
      "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
      "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
      "plan_digest": "sha256:2dafdcc48cb42371008792230a53b6c03b9084f37d8ae5cac79d58ab919fe94f",
      "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
      "acceptance_rows": [
        "E4",
        "E5",
        "E6",
        "E7",
        "E8",
        "E9",
        "E10",
        "E11",
        "E12",
        "E14",
        "E25",
        "E26",
        "E27",
        "E28",
        "E29",
        "E30",
        "E31",
        "E32",
        "E37"
      ],
      "verification": [
        {
          "id": "c2-checkpoint-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "completed",
          "outcome": "pass",
          "requirement": "checkpoint-tests",
          "command": "bench test --package ./internal/gate --run TestReviewCheckpoint",
          "exit_code": 0,
          "native_ref": {
            "ref": "codex:tool/0696f4",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,6090\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:26b6e128a98f830845e62584212004bb975b9c0d4dbc12ccbaf10be6ba987e0f"
          },
          "probe": {
            "mutation": "omit checkpoint evidence validation",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/983e10",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/evaluation.go,swap,failed,24,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestReviewCheckpoint,passed,34\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,5298\nfailures[24]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/failed_restore,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/failed_transport,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/failed_verification,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_axis,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_probe,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_verification,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_verifier,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/pending_axis,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/pending_verification,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/review_as_verification,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/silent_probe,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/skipped_axis,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/stale_verification,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/wrong_verifier,\"review_checkpoint_test.go:58: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Coverage,\"review_checkpoint_test.go:117: canonical axis omitted without refusal: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Spec,\"review_checkpoint_test.go:117: canonical axis omitted without refusal: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Standards,\"review_checkpoint_test.go:117: canonical axis omitted without refusal: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointFindingAndReviewIdentity/findings,\"review_checkpoint_test.go:181: invalid axis accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointFindingAndReviewIdentity/pair,\"review_checkpoint_test.go:181: invalid axis accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointFindingAndReviewIdentity/source,\"review_checkpoint_test.go:181: invalid axis accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointLaterSource/dirty_source,\"review_checkpoint_test.go:97: uncovered delta accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointLaterSource/unreviewed_repair,\"review_checkpoint_test.go:97: uncovered delta accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointMissingRecord,\"review_checkpoint_test.go:161: missing legacy record passed checkpoint: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointReuse,\"review_checkpoint_test.go:76: ordinary green satisfied stronger checkpoint: 1 runs\"\nskips[0]{package,test,reason}:\n",
              "digest": "sha256:5591453d694ac62738d77e868f56ee71f004cd89c7b037bd4cc29bdfcead51ca"
            }
          }
        },
        {
          "id": "c2-axis-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "completed",
          "outcome": "pass",
          "requirement": "axis-tests",
          "command": "bench test --package ./internal/gate --run TestReviewCheckpointCanonicalAxes",
          "exit_code": 0,
          "native_ref": {
            "ref": "codex:tool/c1acb2",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,530\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:cec1fd7009e3473ceb65cf2a65f2aae7e499ba46c7a28ec20f03784c3e6552d3"
          },
          "probe": {
            "mutation": "omit canonical Coverage axis",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/fd92c7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestReviewCheckpointCanonicalAxes,passed,4\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Coverage,\"review_checkpoint_test.go:117: canonical axis omitted without refusal: 0\"\nskips[0]{package,test,reason}:",
              "digest": "sha256:2c251d29863ae93f6cd06a6298e482a346e476289d38fd037623b4f95617e7ad"
            }
          }
        },
        {
          "id": "c2-route-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "completed",
          "outcome": "pass",
          "requirement": "route-tests",
          "command": "bench test --package ./cmd/bench --run TestGateCheckpointRoute",
          "exit_code": 0,
          "native_ref": {
            "ref": "codex:tool/62bce4",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,691\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:e4a00c8dc1e699d77fe4ba807aeb3c31550bd6a1aa6ccf15bb05312156401005"
          }
        },
        {
          "id": "c2-checkpoint-tests-repair",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "completed",
          "outcome": "pass",
          "requirement": "checkpoint-tests",
          "command": "bench test --package ./internal/gate --run TestReviewCheckpoint",
          "exit_code": 0,
          "native_ref": {
            "ref": "codex:tool/217810",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,8248\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:d8f908bed782305bd9a083c4f174e4a334edb3d529087169d533ca29f2c60d3a"
          },
          "probe": {
            "mutation": "omit checkpoint evidence validation",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/18e01c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/gate/evaluation.go,swap,failed,26,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestReviewCheckpoint,passed,36\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,5838\nfailures[26]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/failed_restore,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/failed_transport,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/failed_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_axis,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_probe,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/missing_verifier,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/partial_additional_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/partial_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/pending_axis,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/pending_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/review_as_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/silent_probe,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/skipped_axis,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/stale_verification,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpoint/wrong_verifier,\"review_checkpoint_test.go:60: checkpoint accepted invalid evidence or lost reason: exit 0:\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Coverage,\"review_checkpoint_test.go:129: canonical axis omitted without refusal: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Spec,\"review_checkpoint_test.go:129: canonical axis omitted without refusal: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Standards,\"review_checkpoint_test.go:129: canonical axis omitted without refusal: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointFindingAndReviewIdentity/findings,\"review_checkpoint_test.go:193: invalid axis accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointFindingAndReviewIdentity/pair,\"review_checkpoint_test.go:193: invalid axis accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointFindingAndReviewIdentity/source,\"review_checkpoint_test.go:193: invalid axis accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointLaterSource/dirty_source,\"review_checkpoint_test.go:109: uncovered delta accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointLaterSource/unreviewed_repair,\"review_checkpoint_test.go:109: uncovered delta accepted: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointMissingRecord,\"review_checkpoint_test.go:173: missing legacy record passed checkpoint: 0\"\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointReuse,\"review_checkpoint_test.go:81: ordinary green satisfied stronger checkpoint: 1 runs\"\nskips[0]{package,test,reason}:\n",
              "digest": "sha256:55ec6ba0f19fa31321938513ea487ba3b47623ee348d4d494d8d4638fb333602"
            }
          }
        },
        {
          "id": "c2-axis-tests-repair",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "completed",
          "outcome": "pass",
          "requirement": "axis-tests",
          "command": "bench test --package ./internal/gate --run TestReviewCheckpointCanonicalAxes",
          "exit_code": 0,
          "native_ref": {
            "ref": "codex:tool/51c9b1",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,1201\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:8ca517936b22e9eea86082003efc2cdc586d1ece8c24317447e01138204ed8a9"
          },
          "probe": {
            "mutation": "omit canonical Coverage axis",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/4272f6",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/reviewrecord/record.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestReviewCheckpointCanonicalAxes,passed,4\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,fail,725\nfailures[1]{package,test,line}:\n  github.com/gibbonmi/bench/internal/gate,TestReviewCheckpointCanonicalAxes/Coverage,\"review_checkpoint_test.go:129: canonical axis omitted without refusal: 0\"\nskips[0]{package,test,reason}:\n",
              "digest": "sha256:0480348932c906afaae7ed88b71bd63e6c3c9963555ede52b7c19109e0f17081"
            }
          }
        },
        {
          "id": "c2-route-tests-repair",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "completed",
          "outcome": "pass",
          "requirement": "route-tests",
          "command": "bench test --package ./cmd/bench --run TestGateCheckpointRoute",
          "exit_code": 0,
          "native_ref": {
            "ref": "codex:tool/ac7240",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,1316\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:dca5c3f82527c6ccf4a5b7d98752853aa2f3d701236aec76f7522b548c8ef1ad"
          }
        }
      ],
      "reviews": [
        {
          "id": "c2-standards",
          "performer": "/root/c2_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "pending",
          "outcome": "",
          "axis": "Standards",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "e22288dae5b1514ebb3f223fb72f12e03390cd24",
          "finding_ids": [],
          "supersedes": [],
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          }
        },
        {
          "id": "c2-spec",
          "performer": "/root/c2_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "pending",
          "outcome": "",
          "axis": "Spec",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "e22288dae5b1514ebb3f223fb72f12e03390cd24",
          "finding_ids": [],
          "supersedes": [],
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          }
        },
        {
          "id": "c2-coverage",
          "performer": "/root/c2_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "pending",
          "outcome": "",
          "axis": "Coverage",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "e22288dae5b1514ebb3f223fb72f12e03390cd24",
          "finding_ids": [],
          "supersedes": [],
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          }
        },
        {
          "id": "c2-standards-terminal",
          "performer": "/root/c2_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "completed",
          "outcome": "findings",
          "axis": "Standards",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "e22288dae5b1514ebb3f223fb72f12e03390cd24",
          "finding_ids": [
            "S1"
          ],
          "supersedes": [
            "c2-standards"
          ],
          "native_ref": {
            "ref": "codex:agent//root/c2_standards/terminal-e22288d",
            "excerpt": "One finding; worst issue: stale/dead gate-routing ownership. Repair targets: 1.\nS1 \u2014 auto-fix: Remove unused run_gate and update the stale routing comments. Sources: bin/bench.sh:6-30 and internal/gate/gate.go:6-10. Native terminal conclusion: Standards returned one actionable finding.",
            "digest": "sha256:eefe7a4e39bf810229893fa705b781d2ffa9e327e74b8d9e92f29781a2bbe9be"
          }
        },
        {
          "id": "c2-spec-terminal",
          "performer": "/root/c2_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "completed",
          "outcome": "findings",
          "axis": "Spec",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "e22288dae5b1514ebb3f223fb72f12e03390cd24",
          "finding_ids": [
            "P3",
            "P4",
            "P5"
          ],
          "supersedes": [
            "c2-spec"
          ],
          "native_ref": {
            "ref": "codex:agent//root/c2_spec/terminal-e22288d",
            "excerpt": "Spec review completed with 3 medium findings; worst issue: medium. Disposition for each is auto-fix.\nP3: Name the last covered/requested chunk in the stale-source refusal (internal/reviewrecord/coverage.go:112-113; spec.md:93). P4: Add the planned verification test/file or amend the approved seam before using the alternative (spec.md:233-240). P5: Add a positive wrapper traversal for --checkpoint \u2026 --complete (cmd/bench/gate_route_test.go:68-87). Terminal conclusion: completed with findings.",
            "digest": "sha256:cf4d0d6888e2d424d93d03c386d5f825d7b6c617ca2b1e419ad4378c14cf62e8"
          }
        },
        {
          "id": "c2-coverage-terminal",
          "performer": "/root/c2_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "8e1b770b166fac00e706c8c60572ca9a7eff00d0",
          "state": "completed",
          "outcome": "findings",
          "axis": "Coverage",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "e22288dae5b1514ebb3f223fb72f12e03390cd24",
          "finding_ids": [
            "C2-C1",
            "C2-C2",
            "C2-C3"
          ],
          "supersedes": [
            "c2-coverage"
          ],
          "native_ref": {
            "ref": "codex:agent//root/c2_coverage/terminal-e22288d",
            "excerpt": "Coverage review completed with 3 findings; worst issue high. All three dispositions are auto-fix.\nC2-C1: Add at least two plan requirements, omit one, and require refusal naming that requirement (recordtest/fixture.go:34; review_checkpoint_test.go:40; E25). C2-C2: Add a valid chunk \u2192 complete transition on unchanged source and require another oracle run (review_checkpoint_test.go:67; spec.md:89). C2-C3: Add a valid complete request through the shell using a hostile but valid slug (gate_route_test.go:68; E37). Terminal conclusion: Coverage does not close chunk 2 until C2-C1 through C2-C3 are repaired and re-reviewed.",
            "digest": "sha256:7aad36e6c6161c8f6f1519b87da3943704a25400ddebd0bd64705076214be919"
          }
        },
        {
          "id": "c2-standards-repair-pending",
          "performer": "/root/c2_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "pending",
          "outcome": "",
          "axis": "Standards",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-terminal"
          ],
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          }
        },
        {
          "id": "c2-spec-repair-pending",
          "performer": "/root/c2_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "pending",
          "outcome": "",
          "axis": "Spec",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-terminal"
          ],
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          }
        },
        {
          "id": "c2-coverage-repair-pending",
          "performer": "/root/c2_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "pending",
          "outcome": "",
          "axis": "Coverage",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-terminal"
          ],
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          }
        },
        {
          "id": "c2-standards-repair-clean",
          "performer": "/root/c2_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "completed",
          "outcome": "pass",
          "axis": "Standards",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
          "finding_ids": [],
          "supersedes": [
            "c2-standards-repair-pending"
          ],
          "native_ref": {
            "ref": "codex:agent//root/c2_standards/repair-06a589b",
            "excerpt": "Clean reaffirmation: zero findings; worst issue: none; repair targets: 0. S1 is closed. No repair-delta or cross-chunk Standards concern remains. Native terminal conclusion: Standards repair follow-up completed with zero findings.\nFrozen pair: 5457042ec3e0919327dc7943aea7f002a54285b2..06a589b935ce26fd9e27d18292f20985c81db487. Actual line: gpt-5.6-sol / high / 1 iteration.",
            "digest": "sha256:895d6a25f857eea6490e2a71756bb8e24554c829b1076cf5bc6a1e07bae97523"
          }
        },
        {
          "id": "c2-spec-repair-clean",
          "performer": "/root/c2_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "completed",
          "outcome": "pass",
          "axis": "Spec",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
          "finding_ids": [],
          "supersedes": [
            "c2-spec-repair-pending"
          ],
          "native_ref": {
            "ref": "codex:agent//root/c2_spec/repair-06a589b",
            "excerpt": "Spec repair follow-up completed clean: 0 findings; worst issue none; repair targets 0. P3\u2013P5 are closed. Whole-spec audit: E1\u2013E14, E19\u2013E32, E37, and E38 remain satisfied at this source. E15\u2013E18 and E33\u2013E36 remain correctly deferred to chunk 3. Terminal conclusion: Spec reaffirms the repaired chunk-2 source with no findings.\nFrozen pair: 5457042ec3e0919327dc7943aea7f002a54285b2..06a589b935ce26fd9e27d18292f20985c81db487. Actual line: gpt-5.6-sol / high / 1 iteration.",
            "digest": "sha256:4b8579973127ec9911649091cb493e533eec8452d9254f00f46ba39f9a7af3bb"
          }
        },
        {
          "id": "c2-coverage-repair-clean",
          "performer": "/root/c2_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "3b22f1eb96beee54fe69fb6a9ee56c4e10b6e17e",
          "state": "completed",
          "outcome": "pass",
          "axis": "Coverage",
          "base": "5457042ec3e0919327dc7943aea7f002a54285b2",
          "tip": "06a589b935ce26fd9e27d18292f20985c81db487",
          "finding_ids": [],
          "supersedes": [
            "c2-coverage-repair-pending"
          ],
          "native_ref": {
            "ref": "codex:agent//root/c2_coverage/repair-06a589b",
            "excerpt": "Coverage repair follow-up completed with zero findings; worst issue none. No repair-delta or cross-chunk Coverage finding remains. The shared-seam amendment and r2.md retain the approved rows, and the stale-source diagnostic remains exercised with the chunk identity. Terminal conclusion: current clean Coverage reaffirmation; C2-C1 through C2-C3 are closed at 06a589b935ce26fd9e27d18292f20985c81db487.\nFrozen pair: 5457042ec3e0919327dc7943aea7f002a54285b2..06a589b935ce26fd9e27d18292f20985c81db487. Actual line: gpt-5.6-sol / high / 1 iteration.",
            "digest": "sha256:1be754e84f4597d9a374cf036ff11eb7205e44a48fafbe3476e2a13f2d091ff8"
          }
        }
      ]
    },
    {
      "id": "3",
      "base": "06a589b935ce26fd9e27d18292f20985c81db487",
      "tip": "943a373b3e4b714b7074da233f7ad1c4b7a6f1e5",
      "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
      "plan_digest": "sha256:e298566dc9a9abc5e580968cae0bf69057dfe06abb8e59feb3fb286948d0c754",
      "acceptance_rows": [
        "E15",
        "E16",
        "E17",
        "E18",
        "E33",
        "E34",
        "E35",
        "E36",
        "E20"
      ],
      "verification": [
        {
          "id": "c3-landing-tests",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:tool/1a8294",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/landing,pass,8304\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:92083675a2d8d2a3a5ad8328fbb47f8176692d9827043302a988f27030ff71ed"
          },
          "requirement": "landing-tests",
          "command": "bench test --package ./internal/landing --run TestLandingCompletionEvidence",
          "exit_code": 0,
          "probe": {
            "mutation": "omit landing completion obligation",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/4256ae",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,11,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,TestLandingCompletionEvidence,passed,13\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/landing,fail,4890\n",
              "digest": "sha256:7a3ea6fc1451a54df8ff92fd6e85e4dd964f45e4c925166b060948320065f80f"
            }
          }
        },
        {
          "id": "c3-landing-tests-repair",
          "performer": "/root",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:tool/03cddf",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/landing,pass,14843\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
            "digest": "sha256:8e9139bb84bd1aec7b9bfb10826f11eab66b55152bedffe64759339248261502"
          },
          "requirement": "landing-tests",
          "command": "bench test --package ./internal/landing --run TestLandingCompletionEvidence",
          "exit_code": 0,
          "probe": {
            "mutation": "omit landing completion obligation",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "codex:tool/71684d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,11,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,TestLandingCompletionEvidence,passed,13\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/landing,fail,5181\nfailures[11]{package,test,line}:\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/destination_addition,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"completion composition changes added.txt\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/destination_deletion,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"completion composition changes foreign\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/extra_spec_bytes,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"exact status transform\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/failed_integration,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"verification integration\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/missing_integration,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"verification integration\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/missing_reconciliation,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"reconciliation E1\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/missing_record,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"retain a valid native result record\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/pending_integration,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"verification integration\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/spec_mode_change,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"exact status transform\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/stale_integration,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"verification integration\\\": <nil>\"\n  github.com/gibbonmi/bench/internal/landing,TestLandingCompletionEvidence/wrong_final_verifier,\"completion_evidence_test.go:157: incomplete evidence published or lost reason \\\"verification integration\\\": <nil>\"\nskips[0]{package,test,reason}:\n",
              "digest": "sha256:c1ec32e8a19aa9c7bfae231950f94a1bea115f09e1b662017f4b4c86abdcf26d"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "c3-standards-pending",
          "performer": "/root/c3_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "pending",
          "outcome": "",
          "native_ref": {
            "ref": "",
            "excerpt": "",
            "digest": ""
          },
          "axis": "Standards",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "55520083eb588b3d9aa17667d18a06b9c6a51c5c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c3-spec-pending",
          "performer": "/root/c3_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "pending",
          "outcome": "",
          "native_ref": {
            "ref": "",
            "excerpt": "",
            "digest": ""
          },
          "axis": "Spec",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "55520083eb588b3d9aa17667d18a06b9c6a51c5c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c3-coverage-pending",
          "performer": "/root/c3_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "pending",
          "outcome": "",
          "native_ref": {
            "ref": "",
            "excerpt": "",
            "digest": ""
          },
          "axis": "Coverage",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "55520083eb588b3d9aa17667d18a06b9c6a51c5c",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c3-standards-terminal",
          "performer": "/root/c3_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "codex:agent//root/c3_standards/55520083",
            "excerpt": "Standards review completed with findings.\nCount: 2 hard violations, 0 judgment calls. Worst issue: high. Repair targets: 2.\nS-C3-1 \u2014 high \u2014 auto-fix: The chunk leaves the authoritative system suite red. The shared landing-race fixture creates a staged spec without its completion plan and evidence record, then invokes a spec-backed landing.\nS-C3-2 \u2014 medium \u2014 auto-fix: The new completion owner independently parses the same immutable Git-tree entry through two production paths.\nFrozen pair: 06a589b935ce26fd9e27d18292f20985c81db487..55520083eb588b3d9aa17667d18a06b9c6a51c5c. Actual line: gpt-5.6-sol / high / 1 iteration. Worktree clean; no candidate execution or authorship.",
            "digest": "sha256:04b841c687c1a2aa75aaf9a7ce2ebaa0b6d77cfd10044728441feff596253535"
          },
          "axis": "Standards",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "55520083eb588b3d9aa17667d18a06b9c6a51c5c",
          "finding_ids": [
            "S-C3-1",
            "S-C3-2"
          ],
          "supersedes": [
            "c3-standards-pending"
          ]
        },
        {
          "id": "c3-spec-terminal",
          "performer": "/root/c3_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "codex:agent//root/c3_spec/55520083",
            "excerpt": "Spec axis terminal state: completed with findings. Findings: 2. Worst issue: high. Deduplicated repair targets: 2.\nP6 \u2014 high \u2014 auto-fix. System landing fixtures were not migrated to completion evidence.\nP7 \u2014 high \u2014 auto-fix. The public race journey still expects an unreviewed destination delta to publish.\nComplete acceptance-row audit: E1\u2013E16 satisfied; E17 partial P7; E18 partial P6; E19\u2013E38 satisfied.\nFrozen pair: 06a589b935ce26fd9e27d18292f20985c81db487..55520083eb588b3d9aa17667d18a06b9c6a51c5c. Actual line: gpt-5.6-sol / high / 1 iteration. Worktree clean; no candidate execution or authorship.",
            "digest": "sha256:f8d7ab9fc6b5817b21a47e1e8af90b5d167833b12a597e7fbcdcb8a0bbbb9a90"
          },
          "axis": "Spec",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "55520083eb588b3d9aa17667d18a06b9c6a51c5c",
          "finding_ids": [
            "P6",
            "P7"
          ],
          "supersedes": [
            "c3-spec-pending"
          ]
        },
        {
          "id": "c3-coverage-terminal",
          "performer": "/root/c3_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "codex:agent//root/c3_coverage/55520083",
            "excerpt": "Coverage review completed with findings. Findings: 1. Worst issue: high. Repair targets: 1.\nC3-C1 \u2014 high \u2014 auto-fix. Input: a valid spec-backed public landing from systemLandingRaceFixture. Its staged spec has no bench-completion-plan, and its source has no terminal completion record.\nRepair: migrate the one shared system producer through the completion-record fixture, then preserve the tests\u2019 intended barrier states.\nFrozen pair: 06a589b935ce26fd9e27d18292f20985c81db487..55520083eb588b3d9aa17667d18a06b9c6a51c5c. Actual line: gpt-5.6-sol / high / 1 iteration. Worktree clean; no candidate execution or authorship.",
            "digest": "sha256:50f2c2a5969fb2f15075a7bf2902d0730602943d7cfa0bf72347d3b7214624cd"
          },
          "axis": "Coverage",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "55520083eb588b3d9aa17667d18a06b9c6a51c5c",
          "finding_ids": [
            "C3-C1"
          ],
          "supersedes": [
            "c3-coverage-pending"
          ]
        },
        {
          "id": "c3-standards-repair-terminal",
          "performer": "/root/c3_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:agent//root/c3_standards/943a373b",
            "excerpt": "Standards repair review completed clean.\n\n- Frozen pair: `06a589b935ce26fd9e27d18292f20985c81db487..943a373b3e4b714b7074da233f7ad1c4b7a6f1e5`\n- Actual line: `gpt-5.6-sol / high / 1 iteration`\n- Terminal state: positive, completed\n- Hard violations: 0\n- Judgment calls: 0\n- Worst issue: none\n- Repair targets: 0",
            "digest": "sha256:c34a790123ffb18ea5ae787a1b6479cbfc03b3d9e0b76a910147463553dd58bc"
          },
          "axis": "Standards",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "943a373b3e4b714b7074da233f7ad1c4b7a6f1e5",
          "finding_ids": [],
          "supersedes": [
            "c3-standards-terminal"
          ]
        },
        {
          "id": "c3-spec-repair-terminal",
          "performer": "/root/c3_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:agent//root/c3_spec/943a373b",
            "excerpt": "Spec repair axis terminal state: **completed, clean**.\n\n- Frozen pair: `06a589b935ce26fd9e27d18292f20985c81db487..943a373b3e4b714b7074da233f7ad1c4b7a6f1e5`\n- Repair focus: `55520083eb588b3d9aa17667d18a06b9c6a51c5c..943a373b3e4b714b7074da233f7ad1c4b7a6f1e5`\n- Actual line: `gpt-5.6-sol / high / 1 iteration`\n- Findings: 0\n- Worst issue: none\n- Deduplicated repair targets: 0",
            "digest": "sha256:3a106aed04f774e910c09f0be4375efd8cd54df88cd54ddb49b441206dc6f992"
          },
          "axis": "Spec",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "943a373b3e4b714b7074da233f7ad1c4b7a6f1e5",
          "finding_ids": [],
          "supersedes": [
            "c3-spec-terminal"
          ]
        },
        {
          "id": "c3-coverage-repair-terminal",
          "performer": "/root/c3_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "codex:agent//root/c3_coverage/943a373b",
            "excerpt": "Coverage repair review completed clean.\n\n- Frozen pair: `06a589b935ce26fd9e27d18292f20985c81db487..943a373b3e4b714b7074da233f7ad1c4b7a6f1e5`\n- Repair delta reviewed: `55520083eb588b3d9aa17667d18a06b9c6a51c5c..943a373b3e4b714b7074da233f7ad1c4b7a6f1e5`\n- Actual line: `gpt-5.6-sol` / high / 1 iteration\n- Terminal state: completed clean\n- Findings: 0\n- Worst issue: none\n- Repair targets: 0",
            "digest": "sha256:dfbbe9ba0f648fe179afcabcacee6c011a774e3a4a519756b93e589af87fd99a"
          },
          "axis": "Coverage",
          "base": "06a589b935ce26fd9e27d18292f20985c81db487",
          "tip": "943a373b3e4b714b7074da233f7ad1c4b7a6f1e5",
          "finding_ids": [],
          "supersedes": [
            "c3-coverage-terminal"
          ]
        }
      ]
    }
  ],
  "completion": {
    "state": "completed",
    "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
    "performer": "/root",
    "reconciliation": {
      "E1": "covered",
      "E2": "covered",
      "E3": "covered",
      "E4": "covered",
      "E5": "covered",
      "E6": "covered",
      "E7": "covered",
      "E8": "covered",
      "E9": "covered",
      "E10": "covered",
      "E11": "covered",
      "E12": "covered",
      "E13": "covered",
      "E14": "covered",
      "E15": "covered",
      "E16": "covered",
      "E17": "covered",
      "E18": "covered",
      "E19": "covered",
      "E20": "covered",
      "E21": "covered",
      "E22": "covered",
      "E23": "covered",
      "E24": "covered",
      "E25": "covered",
      "E26": "covered",
      "E27": "covered",
      "E28": "covered",
      "E29": "covered",
      "E30": "covered",
      "E31": "covered",
      "E32": "covered",
      "E33": "covered",
      "E34": "covered",
      "E35": "covered",
      "E36": "covered",
      "E37": "covered",
      "E38": "covered"
    },
    "verification": [
      {
        "id": "c3-integration-initial",
        "performer": "/root",
        "role": "author-verification",
        "model": "gpt-6-astra",
        "effort": "high",
        "source_digest": "c9201fa276fabba32c215fa364122fd46a775776",
        "state": "failed",
        "outcome": "fail",
        "native_ref": {
          "ref": "codex:tool/87aeb6",
          "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/systemtest,fail,41394\nfailures[4]{package,test,line}:\n",
          "digest": "sha256:8ba34572134bc79e6b1c0843500e419f5a0aadc602fb89b60a8a830e0e4013f1"
        },
        "requirement": "integration",
        "command": "bench test --check system",
        "exit_code": 1
      },
      {
        "id": "c3-acceptance-selector-refusal",
        "performer": "/root",
        "role": "author-verification",
        "model": "gpt-6-astra",
        "effort": "high",
        "source_digest": "1807c31d4fe57a17dd7b1bab70eb42ccbb68a765",
        "state": "failed",
        "outcome": "fail",
        "native_ref": {
          "ref": "codex:tool/f4a1d7",
          "excerpt": "error: changed selection failed \u2014 changed Go path is not in a current package\nworktree: /home/mgibs/.bench/worktrees/bench-2826441890/c615fc296ef1abf0cc0456b06d3b1d9c-9d9cea28c60222a3010e9e1d497d8da1\n",
          "digest": "sha256:60b0122516a1d97d96340c9333e41f71d2df4b899d19a0c1dbe59ea83829ce0d"
        },
        "requirement": "acceptance",
        "command": "bench test --changed --base de1447b31903b170679b66bd200431cc15ddc73f",
        "exit_code": 1
      },
      {
        "id": "c3-acceptance-final",
        "performer": "/root",
        "role": "author-verification",
        "model": "gpt-6-astra",
        "effort": "high",
        "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "codex:tool/e67fc9",
          "excerpt": "packages[97]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench,pass,4\n  github.com/gibbonmi/bench/cmd/bench,pass,40099\n  github.com/gibbonmi/bench/internal/adopt,pass,21270\n  github.com/gibbonmi/bench/internal/anchors,pass,526\n  github.com/gibbonmi/bench/internal/axi,pass,17\n  github.com/gibbonmi/bench/internal/axi/axitest,pass,6\n  github.com/gibbonmi/bench/internal/benchguard,pass,5\n  github.com/gibbonmi/bench/internal/benchhome,pass,5\n  github.com/gibbonmi/bench/internal/bounds,pass,524\n  github.com/gibbonmi/bench/internal/brokermanifest,no-tests,0\n  github.com/gibbonmi/bench/internal/canary,pass,81\n  github.com/gibbonmi/bench/internal/canonicalpath,pass,15\n  github.com/gibbonmi/bench/internal/capability,pass,32\n  github.com/gibbonmi/bench/internal/census,pass,849\n  github.com/gibbonmi/bench/internal/commit,pass,15122\n  github.com/gibbonmi/bench/internal/conformance,pass,42663\n  github.com/gibbonmi/bench/internal/conformance/registry,pass,3\n  github.com/gibbonmi/bench/internal/consumers,pass,7788\n  github.com/gibbonmi/bench/internal/contract,pass,3\n  github.com/gibbonmi/bench/internal/coverage,pass,7213\n  github.com/gibbonmi/bench/internal/dashboard,pass,5181\n  github.com/gibbonmi/bench/internal/diff,pass,44838\n  github.com/gibbonmi/bench/internal/env,pass,17\n  github.com/gibbonmi/bench/internal/freshness,pass,28414\n  github.com/gibbonmi/bench/internal/freshness/check,no-tests,0\n  github.com/gibbonmi/bench/internal/gate,pass,33242\n  github.com/gibbonmi/bench/internal/gate/authorization,pass,1068\n  github.com/gibbonmi/bench/internal/gate/greenmarker,pass,317\n  github.com/gibbonmi/bench/internal/gate/prospectiveartifact,pass,2911\n  github.com/gibbonmi/bench/internal/git,pass,8099\n  github.com/gibbonmi/bench/internal/gitguard,pass,20173\n  github.com/gibbonmi/bench/internal/gittest,no-tests,0\n  github.com/gibbonmi/bench/internal/gocache,pass,401\n  github.com/gibbonmi/bench/internal/gocache/cleanprobe,no-tests,0\n  github.com/gibbonmi/bench/internal/guards,pass,234\n  github.com/gibbonmi/bench/internal/handoff,pass,11770\n  github.com/gibbonmi/bench/internal/handoffdoc,pass,2031\n  github.com/gibbonmi/bench/internal/harness,no-tests,0\n  github.com/gibbonmi/bench/internal/harnesses,pass,4\n  github.com/gibbonmi/bench/internal/intent,pass,11517\n  github.com/gibbonmi/bench/internal/intent/admissionpolicy,pass,7\n  github.com/gibbonmi/bench/internal/intent/ledger,pass,5\n  github.com/gibbonmi/bench/internal/jsonfile,pass,5\n  github.com/gibbonmi/bench/internal/landing,pass,31800\n  github.com/gibbonmi/bench/internal/landing/settlepolicy,pass,5\n  github.com/gibbonmi/bench/internal/learnings,pass,323\n  github.com/gibbonmi/bench/internal/lines,pass,13\n  github.com/gibbonmi/bench/internal/maps,pass,868\n  github.com/gibbonmi/bench/internal/modelid,pass,6\n  github.com/gibbonmi/bench/internal/modelid/modelidtest,no-tests,0\n  github.com/gibbonmi/bench/internal/models,pass,558\n  github.com/gibbonmi/bench/internal/otelrecord,pass,407\n  github.com/gibbonmi/bench/internal/outline,pass,600\n  github.com/gibbonmi/bench/internal/packagesurface,pass,19\n  github.com/gibbonmi/bench/internal/poolkey,pass,64\n  github.com/gibbonmi/bench/internal/preflight,pass,64530\n  github.com/gibbonmi/bench/internal/preprelease,pass,43\n  github.com/gibbonmi/bench/internal/probe,pass,41457\n  github.com/gibbonmi/bench/internal/prose,pass,619\n  github.com/gibbonmi/bench/internal/publication,pass,9947\n  github.com/gibbonmi/bench/internal/puritycensus,pass,5\n  github.com/gibbonmi/bench/internal/racetests,no-tests,0\n  github.com/gibbonmi/bench/internal/refresh,pass,36\n  github.com/gibbonmi/bench/internal/releaseevidence,pass,2992\n  github.com/gibbonmi/bench/internal/releasepreflight,pass,7\n  github.com/gibbonmi/bench/internal/retros,pass,34\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2613\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0\n  github.com/gibbonmi/bench/internal/roadmap,pass,6029\n  github.com/gibbonmi/bench/internal/roadmap/roadmaptest,no-tests,0\n  github.com/gibbonmi/bench/internal/roadmapflow,pass,680\n  github.com/gibbonmi/bench/internal/runbinary,pass,25053\n  github.com/gibbonmi/bench/internal/sanitize,pass,5\n  github.com/gibbonmi/bench/internal/sessioninspect,pass,2026\n  github.com/gibbonmi/bench/internal/shellcommand,pass,2\n  github.com/gibbonmi/bench/internal/shift,pass,3868\n  github.com/gibbonmi/bench/internal/skillsindex,pass,194\n  github.com/gibbonmi/bench/internal/spec,pass,931\n  github.com/gibbonmi/bench/internal/status,pass,22389\n  github.com/gibbonmi/bench/internal/stophook,pass,355\n  github.com/gibbonmi/bench/internal/structure,pass,1409\n  github.com/gibbonmi/bench/internal/subprocess,pass,593\n  github.com/gibbonmi/bench/internal/terminal,no-tests,0\n  github.com/gibbonmi/bench/internal/testlines,pass,4\n  github.com/gibbonmi/bench/internal/testrepo,no-tests,0\n  github.com/gibbonmi/bench/internal/testreport,pass,26595\n  github.com/gibbonmi/bench/internal/tickets,pass,10\n  github.com/gibbonmi/bench/internal/toon,pass,3\n  github.com/gibbonmi/bench/internal/usage,pass,3\n  github.com/gibbonmi/bench/internal/worktree,pass,152226\n  github.com/gibbonmi/bench/internal/worktree/landingpolicy,pass,4\n  github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy,pass,7\n  github.com/gibbonmi/bench/internal/worktree/reclaimpolicy,pass,4\n  github.com/gibbonmi/bench/internal/writeguard,pass,3\n  github.com/gibbonmi/bench/tests/canary/guard-classifier-table/word-test-drops-xargs/files/internal/conformance,pass,2\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/default-branch-refabricated/files/internal/git,no-tests,0\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/default-branch-refabricated/files/internal/roadmap,no-tests,0\nfailures[0]{package,test,line}:\nskips[10]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/capability,TestCapabilityWritesLineBeforeSkip/file_transport/skips,\"capability_test.go:120: requires a host fifo\"\n  github.com/gibbonmi/bench/internal/capability,TestCapabilityWritesLineBeforeSkip/stdout_fallback/skips,\"capability_test.go:144: requires a host fifo\"\n  github.com/gibbonmi/bench/internal/capability,TestEnvironmentWritesLineBeforeSkip/skips,\"capability_test.go:165: subject root has no bin/bench.sh\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/TestGuidanceProseBudgetRefusesNonRegularSubjectssocket3598157310/001/.agents/skills/bench-craft-linked/SKILL.md: bind: invalid argument\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"environment: BENCH_CONFORMANCE_ROOT not set\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/worktree,TestCleanLandedSpecialPathsRetainedWithoutOpening/socket,\"capability: fifo: unix sockets unavailable: listen unix /tmp/TestCleanLandedSpecialPathsRetainedWithoutOpeningsocket3843996133/001/.bench-home/worktrees/001-2454750424/40f7b64c89793a53690db58b5f594fb8-b0524f47541b3954dc2a5cbc8e5a9c14: bind:\u2026 (257 bytes)\"\n  github.com/gibbonmi/bench/internal/worktree,TestLandedConsumersRejectSpecialGitMetadataBeforePlanning/socket,\"capability: fifo: unix sockets unavailable: listen unix /tmp/TestLandedConsumersRejectSpecialGitMetadataBeforePlanningsocket3259797845/001/.bench-home/worktrees/001-4283935671/1d094b12e696abe5e55ef5861a77ed6f-0c1ad5c5445297e4edb980c4f697bf5\u2026 (270 bytes)\"\nroot_conformance[1]{package,status,route}:\n  github.com/gibbonmi/bench/internal/conformance,skipped,bench test --check <name>\n",
          "digest": "sha256:72d121ca9fda2be59b5c606f25fc751b3fb9e828e9a61cda6a37f936a6ec571e"
        },
        "requirement": "acceptance",
        "command": "bench test --package ./...",
        "exit_code": 0
      },
      {
        "id": "c3-integration-final",
        "performer": "/root",
        "role": "author-verification",
        "model": "gpt-6-astra",
        "effort": "high",
        "source_digest": "46ce309594a597037464d2511c9d2d9da3c0e3fd",
        "state": "completed",
        "outcome": "pass",
        "native_ref": {
          "ref": "codex:tool/de0bc0",
          "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,96472\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
          "digest": "sha256:45b64836f49c508dec8549dd09f48b7111654c08c4273db70ba3621284adb627"
        },
        "requirement": "integration",
        "command": "bench test --check system",
        "exit_code": 0
      }
    ]
  },
  "amendments": [
    {
      "from": "sha256:267c4a827157b8e6acb6a9a24a56b984305fdea06ed5601964bef8e434df0f05",
      "to": "sha256:344ffe2a529e7043e5bd579bdd1bec47f28129f293c913f5d98a0e33cf8be4dc",
      "chunk_ids": {
        "1": [
          "1"
        ],
        "2": [
          "2"
        ],
        "3": [
          "3"
        ]
      }
    },
    {
      "from": "sha256:344ffe2a529e7043e5bd579bdd1bec47f28129f293c913f5d98a0e33cf8be4dc",
      "to": "sha256:2dafdcc48cb42371008792230a53b6c03b9084f37d8ae5cac79d58ab919fe94f",
      "chunk_ids": {
        "1": [
          "1"
        ],
        "2": [
          "2"
        ],
        "3": [
          "3"
        ]
      }
    },
    {
      "from": "sha256:2dafdcc48cb42371008792230a53b6c03b9084f37d8ae5cac79d58ab919fe94f",
      "to": "sha256:013da5369b1c20983a9334647a1662e26944f81926d81846258ef531f71fd2ed",
      "chunk_ids": {
        "1": [
          "1"
        ],
        "2": [
          "2"
        ],
        "3": [
          "3"
        ]
      }
    },
    {
      "from": "sha256:013da5369b1c20983a9334647a1662e26944f81926d81846258ef531f71fd2ed",
      "to": "sha256:e298566dc9a9abc5e580968cae0bf69057dfe06abb8e59feb3fb286948d0c754",
      "chunk_ids": {
        "1": [
          "1"
        ],
        "2": [
          "2"
        ],
        "3": [
          "3"
        ]
      }
    }
  ]
}
```
