# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/landing-stack/spec.md",
  "plan_digest": "sha256:3a488226f798d4b679e90398c8906526574cada61446e16eb27202fe49e3d6e0",
  "implementation_session": "",
  "chunks": [
    {
      "id": "LS-C1",
      "base": "064889b693de71257f22b97e8e49bccbf83c95c8",
      "tip": "d172cd9967dcd2b4dff64de7732b5c959f388336",
      "plan_digest": "sha256:3a488226f798d4b679e90398c8906526574cada61446e16eb27202fe49e3d6e0",
      "source_digest": "7d04e4680b5236884008f6227e82a4e7453b43a5",
      "acceptance_rows": [
        "LS1",
        "LS2",
        "LS3",
        "LS4",
        "LS5",
        "LS6"
      ],
      "verification": [
        {
          "id": "ls-c1-t01-stackpolicy",
          "performer": "claude:ls-t01-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7d04e4680b5236884008f6227e82a4e7453b43a5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ls-t01-author-20261026",
            "digest": "sha256:7b983f6943e66e81f2124a65def105d8e31e7574ca68d71b9c4ffb468588b533",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing/stackpolicy,pass,2,6\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d172cd9967dcd2b4dff64de7732b5c959f388336",
          "requirement": "stackpolicy",
          "command": "bench test --package ./internal/landing/stackpolicy",
          "exit_code": 0
        },
        {
          "id": "ls-c1-t01-root-conformance",
          "performer": "claude:ls-t01-author-20261026",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7d04e4680b5236884008f6227e82a4e7453b43a5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ls-t01-author-20261026",
            "digest": "sha256:44bee3d6e2c90f772c5f847ba44042c3fa557822ccc4b2eb47b785ad4287d3f5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15104,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d172cd9967dcd2b4dff64de7732b5c959f388336",
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ls-c1-standards-r1",
          "performer": "claude:ls-c1-standards-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7d04e4680b5236884008f6227e82a4e7453b43a5",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ls-c1-standards-20261026",
            "digest": "sha256:e51342689b62669c0e3a20b3385e506ffe33baf7684b8eb0d9c7d224c4343ff4",
            "excerpt": "axis: Standards\nverdict: findings\nfindings: 1\nworst: S1, severity 3: the \"capture/\" path rule is written a second time, outside its settle-policy owner.\n\nS1 | 3 | ask-user | 5 | internal/landing/stackpolicy/stackpolicy.go:16 (`const captureDir = \"capture/\"`) and :43 (`strings.HasPrefix(path, captureDir)`) | The predicate \"a path is a capture path\" is the `strings.HasPrefix(path, \"capture/\")` branch of `settlepolicy.CaptureSide` (internal/landing/settlepolicy/settlepolicy.go:108). The new code re-derives it. AGENTS.md \"one source per fact\" calls two derivations of one fact a defect. The repair needs a decision: let stackpolicy read the settle policy, which lifts the stdlib-only line, or accept the copy.\n\nAdvice:\n- Purity census: no mandatory rule needs purity_census_test.go; it is precedent in settlepolicy and five other policy packages. Without one, the package doc's claim at stackpolicy.go:3-4 has no check behind it. A one-line census modelled on settlepolicy/purity_census_test.go:14 would back that claim; it needs `bench spec widen`.\n- stackpolicy.go:2-4 is one long semicolon sentence; two or three short sentences read better.\n- The \"(Coverage row LSn.)\" test comments match the settlepolicy_test.go convention.\n\nChecked with no defect: independent literal expectations; the spec edit changes only the LS1-LS6 cells; the Owners and Overlaps bodies are correct and deterministic.\n\nCommand contribution: no command change is necessary.\nRead: preflight bind output, the frozen-pair diff, git show d172cd996, ticket 01, both stackpolicy files, settlepolicy.go and its purity_census_test.go, internal/puritycensus/census.go header, craft-review and craft-comments skills, spec lines 138-146 and rows LS1-LS6, rg sweeps for capture-prefix literals and purity rules.\n"
          },
          "axis": "Standards",
          "base": "064889b693de71257f22b97e8e49bccbf83c95c8",
          "tip": "d172cd9967dcd2b4dff64de7732b5c959f388336",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "ls-c1-spec-r1",
          "performer": "claude:ls-c1-spec-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7d04e4680b5236884008f6227e82a4e7453b43a5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ls-c1-spec-20261026",
            "digest": "sha256:b5ab3cbd20c6f6fe9ce30049120b294f51f9f6ac976ba7a6cb63c6c23ef8ad7f",
            "excerpt": "axis: Spec\nverdict: pass\nfindings: 0\nworst: none\n\nAudit of LS1 to LS6, re-derived from stackpolicy.go and the spec text (spec.md:283-288 and \"The owner rule\", spec.md:140-144):\n- LS1: Owners (stackpolicy.go:21-36) appends every differing member in member order; the member loop and the destination-only loop are exclusive per path. TestStackOwnersNamesEachChangingMember (stackpolicy_test.go:26) checks a two-owner path, a one-owner path, and an unchanged path.\n- LS2: Overlaps skips the capture/ prefix (stackpolicy.go:43); test at stackpolicy_test.go:40.\n- LS3: test at stackpolicy_test.go:53 keeps internal/x.go and leaves out alone.go.\n- LS4: metadata compared as a whole string (stackpolicy.go:25); test at stackpolicy_test.go:65.\n- LS5: destination-only loop plus member-only branch give old and new paths; test at stackpolicy_test.go:73.\n- LS6: paths are opaque map keys; test at stackpolicy_test.go:83.\n\nScope: the code edits are the two new files only; the spec edit is coverage-cell citations. bench coverage maps LS1 to LS6 to six existing tests. No scope creep.\n\nShape for later consumers: ticket 10 calls Owners then Overlaps with one Entries type; ticket 02 fits for rules 2, 3, and the LS13 deletion.\n\nAdvice:\n- Ticket 02 must not write a second capture/ test; it can derive rules 4 and 5 from len(owners[p])>=2 and membership in Overlaps(owners).\n- Both functions return maps, so ticket 10 must sort paths before it prints refusals, and must expand a three-owner path into pairs.\n- captureDir repeats the capture/ literal at settlepolicy.go:108; a Standards decision.\n\nCommand contribution: no command change is necessary.\nRead: craft-review SKILL.md, the base-to-tip delta, stackpolicy.go, stackpolicy_test.go, spec.md lines 56-145 and rows, tickets 01, 02, 10, bench coverage output, settlepolicy.go capture lines.\n"
          },
          "axis": "Spec",
          "base": "064889b693de71257f22b97e8e49bccbf83c95c8",
          "tip": "d172cd9967dcd2b4dff64de7732b5c959f388336",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ls-c1-coverage-r1",
          "performer": "claude:ls-c1-coverage-20261026",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "7d04e4680b5236884008f6227e82a4e7453b43a5",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ls-c1-coverage-20261026",
            "digest": "sha256:e62e3959090ee4df54f4d82931a6ba289db841e4bb466ff51f8863a71e16fa44",
            "excerpt": "axis: Coverage\nverdict: findings\nfindings: 3\nworst: C1 — the deletion branch of Owners (stackpolicy.go:29-33) has no test with two members, so a last-owner-only mutant there stays green.\n\nC1 | 5 | auto-fix | 8 | internal/landing/stackpolicy/stackpolicy_test.go:73-79 against stackpolicy.go:29-33 | Only the rename test (one member) exercises the destination-only path loop. No test has a path that two members both delete. Overwriting instead of appending at line 31 keeps every test green, so a shared deletion would report one owner and hide the overlap. Add a case where two members both drop a destination path and expect {1,2}; a mixed case (member 1 deletes, member 2 modifies) also pins the index order.\nC2 | 3 | auto-fix | 8 | stackpolicy_test.go:53-61 against stackpolicy.go:43 | Every overlap fixture has exactly two owners. `len(indexes) == 2` in place of `>= 2` keeps all tests green, so a path that three members change would drop out of the overlap set. Add a three-member path to the LS3 test.\nC3 | 4 | ask-user | 7 | stackpolicy_test.go:40-48 against stackpolicy.go:16,43 | LS2 and LS3 pin only `capture/learnings.md` and `internal/x.go`. `strings.Contains(path, \"capture/\")`, `HasPrefix(path, \"capture\")`, and a case-folded match are not killed; each hides a real overlap. Decide the boundary, then add one near-prefix path that two members change and expect it to stay in the overlap set.\n\nRefutations: the 8 reported witnesses are red-capable by reading of the code.\n\nAdvice:\n- LS2 passes vacuously if Owners returned nothing; an assertion that Owners lists the capture path would make LS2 self-sufficient.\n- LS6 catches key mangling inside this package only; the real parsing risk sits in the caller tickets.\n- Add one empty-members case to pin a non-nil empty map.\n\nCommand: no command change is necessary.\nRead: stackpolicy.go, stackpolicy_test.go, ticket 01, spec rows LS1-LS6 and the chunk table, craft-review SKILL.md, the benchkit hostile-input checklist, the delta diff.\n"
          },
          "axis": "Coverage",
          "base": "064889b693de71257f22b97e8e49bccbf83c95c8",
          "tip": "d172cd9967dcd2b4dff64de7732b5c959f388336",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  }
}
```

## LS-C1 review, round 1

Frozen pair: `064889b693de71257f22b97e8e49bccbf83c95c8..d172cd9967dcd2b4dff64de7732b5c959f388336`. Repair cycles used: 0 of 2. Every accepted finding falls on ticket 01, so one combined repair charge carries them.

## Standards

Findings: 1. Worst: S1, severity 3.

- S1 (confidence 5), `auto-fix`: `internal/landing/stackpolicy/stackpolicy.go:16` and `:43` write the `capture/` prefix a second time. The settle policy owns it at `internal/landing/settlepolicy/settlepolicy.go:108`. The orchestrator decides the repair: the overlap set leaves out each path that `settlepolicy.CaptureSide` settles, so one rule names a capture path. The purity census of the policy package permits that import. The spec and ticket 01 sentences that say "imports only the standard library" take a prose correction.

## Spec

Findings: 0. Worst: none.

## Coverage

Findings: 3. Worst: C1, severity 5.

- C1 (confidence 8), `auto-fix`: no test has a path that two members delete, so an overwrite in the deletion branch of `Owners` stays green. Add a shared deletion and a mixed deletion and change.
- C2 (confidence 8), `auto-fix`: every overlap fixture has two owners, so `== 2` in place of `>= 2` stays green. Add a path that three members change.
- C3 (confidence 7), `auto-fix`: no near-prefix path pins the capture boundary. The orchestrator decides that the boundary is the settle policy's root-anchored `capture/` rule, which S1 adopts. Add `internal/capture/x.go` and `capture-notes.md` as overlaps that stay in the set.

## Advice

- Add a purity census to the new package, as `settlepolicy` has. The repair adds it, because the S1 repair changes the import set that the census grades.
- Ticket 02 derives its shared-capture rule from `Overlaps`, and it writes no second capture test.
- Ticket 10 sorts the paths before it prints the refusals, and it expands a path with three owners into pairs.
