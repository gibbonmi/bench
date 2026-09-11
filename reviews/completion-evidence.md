# Completion evidence review

## Standards

Pending native return.

## Spec

Pending native return.

## Coverage

Pending native return.

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

## Machine record

```bench-review-record
{
  "version": 1,
  "spec": "specs/completion-evidence/spec.md",
  "plan_digest": "sha256:ecc3b81b5ceb5dfce8779c67628e3da53115ef86f98c94c53200541d13732436",
  "implementation_session": "/root",
  "chunks": [
    {
      "id": "1",
      "base": "de1447b31903b170679b66bd200431cc15ddc73f",
      "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
      "plan_digest": "sha256:ecc3b81b5ceb5dfce8779c67628e3da53115ef86f98c94c53200541d13732436",
      "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
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
          "state": "pending",
          "outcome": "",
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          },
          "axis": "Standards",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-spec",
          "performer": "/root/c1_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "pending",
          "outcome": "",
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          },
          "axis": "Spec",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-coverage",
          "performer": "/root/c1_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "5ebad816a00190ded84344b26c34bd07e818a01b",
          "state": "pending",
          "outcome": "",
          "native_ref": {
            "ref": "",
            "digest": "",
            "excerpt": ""
          },
          "axis": "Coverage",
          "base": "de1447b31903b170679b66bd200431cc15ddc73f",
          "tip": "205c92928a4bbf35b4540e1f867808c96310d97a",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "/root",
    "reconciliation": {},
    "verification": []
  }
}
```
