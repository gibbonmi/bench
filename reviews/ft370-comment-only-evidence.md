# FT370 review outcomes

The reviewer directs one Codex session to implement all four tickets.
Three independent Sol 6.1 high sessions review each chunk.

## Standards

CG-C1: 0 findings. The axis accepts the reader move.

## Spec

CG-C1: 0 findings. The axis accepts the reader move.

## Coverage

CG-C1: 0 findings. The axis accepts the reader move.

## Verification

Ticket 1 preserves the existing change-list tests.
The source-mode and destination-mode swap failed the rename test, and Bench restored the file.
The Git and gate packages passed, as did the root conformance test.

### Classifier verification

CG-C2 passes its classifier, Git, gate, and root checks.
The count-only directive mutation fails CG15, CG17, and CG18.
A Go-rule bypass fails all 16 refusal cases in the byte-pair table.
A forced token refusal fails CG9, CG10, CG11, and CG31.
A tree-rule bypass fails all 12 refusal cases in the tree table.

The configured-reader substitution fails CG46's mixed hidden-gitlink case.
Bench restores each mutation before the next check.
The initial fixture compile failed on an integer conversion and then passed.
That failed build supplied no behavioral evidence.

## Repair allowance

No chunk has consumed a repair cycle.

```bench-review-record
{
  "version": 1,
  "spec": "specs/ft370-comment-only-evidence/spec.md",
  "plan_digest": "sha256:2d9c1d57d328b461a407e5ab8bb2886d987d970101ba8352087b48a6a8adc1bf",
  "implementation_session": "ft370-root",
  "chunks": [
    {
      "id": "CG-C1",
      "base": "b0abe2fb2a4fa3764bed97896acf13d6f99fed03",
      "tip": "67f6159ba5f6df5961bee2c5e7215e43e18e01c7",
      "plan_digest": "sha256:854ff2a41d464bc649da1d3fb4e1c9525a31ff11195d26b4cc4debcf71c9dd67",
      "source_digest": "fe9bad50f069fb1e7b06a867cc8c6107ac8b0441",
      "acceptance_rows": [
        "CG40",
        "CG41"
      ],
      "verification": [
        {
          "id": "c1-git",
          "performer": "ft370-root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe9bad50f069fb1e7b06a867cc8c6107ac8b0441",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "chat:ft370-root/c1-git",
            "digest": "sha256:3ad8a9eb6d73a20d9734f11cfa1d98125253961dc2cef4be7e3e3d91074afa91",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/git,pass,1594\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "git",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "c1-gate",
          "performer": "ft370-root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "fe9bad50f069fb1e7b06a867cc8c6107ac8b0441",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "chat:ft370-root/c1-gate",
            "digest": "sha256:a2cc81460e743fb601591c3453379e612992df732e83e28cffa3716bd771d0be",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,12321\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "c1-standards",
          "performer": "/root/c1_standards",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "fe9bad50f069fb1e7b06a867cc8c6107ac8b0441",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:/root/c1_standards",
            "digest": "sha256:eeb6e02702de89d21aaeda86f2cfdb16d18e4baf65642de61b8e0683021a4b8b",
            "excerpt": "Standards: pass. Count: 0. Worst issue: none. Confidence: 9/10.\nCG-C1: b0abe2fb2a4fa3764bed97896acf13d6f99fed03..005e12afa0480ae5d17f83047839e3c4d639ef36; code tip 67f6159ba5f6df5961bee2c5e7215e43e18e01c7.\nCG40 satisfied: internal/gate/lane_select.go:22 aliases git.TreeChange; :26 calls git.TreeChanges. The search across internal found no surviving parseComposedChange or composedChangeFields. internal/git/tree.go:217 and :236 own the moved NUL-framed reader and entry parser. The existing reverseAppliesToDefault mode check is unchanged and does not duplicate the moved path reader.\nCG41: the diff preserves the existing tests and their assertions at internal/gate/lane_select_test.go:76, :99, and :122. This read-only review did not rerun them.\nThe moved comments explain metadata and NUL framing, use present tense, and add no provenance or red transcript. The gate wrapper keeps its diagnostic. No new import or dependency appears in internal/git/tree.go.\nThe explicit user override is retained at specs/ft370-comment-only-evidence/spec.md:46-51; the version 1 plan is legitimate for this run. Later chunks remain pending and are not graded as omissions here.\nOptional advice: none.\nSources read: AGENTS.md; .bench/BENCH.md; projects/benchkit.md (seams, hostile-input checklist, gate and line rules, cold-session notes); craft-review SKILL.md and finding-discipline.md; bounded-repair-policy.md; craft-comments SKILL.md; the full FT370 spec; ticket 1; the one frozen git diff; internal/git/tree.go; internal/git/git.go:219-251; internal/gate/lane_select.go (changed wrapper and lane consumers); internal/gate/lane_select_test.go:55-135; internal/gate/authorization/authorization.go:205-255. The diff includes the review record and its verification entries.\nEvidence fetched: sha256:5a38fbf5bf08ceb7555f2d9ebff75ff94188a3d6c0c91dc0f7ec06e760c9a988; s1 index 0 and s13 index 0, both complete with no next cursor. Current-action binding checked once and returned current=true at the frozen tip. Spilled CLI projections were read to completion.\nImplementation command contribution: none; .agents/commands/bench-implement-spec.md was not read. No tests, probes, source writes, stash, or second diff.\n"
          },
          "axis": "Standards",
          "base": "b0abe2fb2a4fa3764bed97896acf13d6f99fed03",
          "tip": "67f6159ba5f6df5961bee2c5e7215e43e18e01c7",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-spec",
          "performer": "/root/c1_spec",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "fe9bad50f069fb1e7b06a867cc8c6107ac8b0441",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:/root/c1_spec",
            "digest": "sha256:6fbddb76b84d1f43d007930f9308f1a1cd2cb228e9dc94c6411e932c64b71b37",
            "excerpt": "Spec review: FT370 CG-C1\nLine: gpt-6.1-sol / high; one iteration.\nFindings: 0. Worst issue: none. Blockers: none. Optional advice: none.\n\nSubject\nBench label ft370-build. Base b0abe2fb2a4fa3764bed97896acf13d6f99fed03.\nReview tip 005e12afa0480ae5d17f83047839e3c4d639ef36; code tip 67f6159ba5f6df5961bee2c5e7215e43e18e01c7.\nEvidence sha256:5a38fbf5bf08ceb7555f2d9ebff75ff94188a3d6c0c91dc0f7ec06e760c9a988 binds current=true to this clean assignment.\nThe later commit records verification. Root-session implementation and the version 1 plan are approved in spec.md:44-51.\n\nDerivation and row conclusions\nI read the complete spec, all 45 acceptance rows, and ticket 1 before grading the implementation.\nCG40: satisfied by source inspection. The spec says \"internal/gate holds no raw-diff entry parser after the move, and gate.ComposedChanges calls git.TreeChanges\" (spec.md:395).\nThe alias and call are at internal/gate/lane_select.go:22-28. The parser and metadata field count moved to internal/git/tree.go:199-242.\nA search of internal/gate and internal/git found no parseComposedChange or composedChangeFields. CG40's seam is Standards-owned; this Spec conclusion supplies corroborating source evidence.\nCG41: source matches the requirement. The unchanged TestComposedChangesExpandsANamedDirectory pins both files, order, and modes at internal/gate/lane_select_test.go:67-82.\nThe single reviewed diff changes neither this test nor the rename and symlink tests at lines 87-125.\nThe moved parser keeps the original field mapping, NUL framing, path bytes, and append order.\nThe record's author verification reports passes for internal/git and internal/gate. I did not independently execute those packages.\n\nTicket 1 contracts\nTreeChanges has the exact exported signature and four fields (internal/git/tree.go:207-218).\nIt passes from and to unchanged to Raw. Raw returns stdout bytes verbatim (internal/git/git.go:242-250), so unnamed tree objects need no commit peel.\nComposedChanges peels only its base with ^{tree}, as required, and preserves the required gate refusal prefix.\nGit failures and rejected metadata entries return errors. The entry parser is moved, not rewritten.\nNo import was added to internal/git. tree.go is 242 lines; lane_select.go is 332 lines, down from 370.\nThe recorded mode-swap mutation is consistent with the rename test's independent expectations for 000000 and 100644 (lane_select_test.go:99-104).\nI did not rerun that mutation or go list -deps; the charge prohibits probes and test execution.\n\nOther acceptance rows\nCG-C2 retains CG9-CG32, CG42, and CG43. CG-C3 retains CG1-CG8, CG33-CG39, CG44, and CG45 (spec.md:320-325).\nThose rows remain pending later chunks. They are not missing requirements in CG-C1.\nStory 39 remains an explicit reviewed exclusion (spec.md:402,443-453).\nThe diff introduces no checkpoint or classifier behavior before those chunks.\n\nReads and command contribution\nRead craft-review, finding-discipline, and bounded repair policy; read the supplied working agreement and shared platform rules.\nRead the full spec, ticket 1, internal/git/tree.go, Raw, the lane wrapper and classification context, the three unchanged reader tests, and LaneAuthority.Authorize.\nRead exactly one git diff of the specified base and review tip, including its spill file.\nRetrieved sources s1 page 0, s13 page 0, s14 page 0, and s14 page 1; all streams end complete. The only nonzero cursor used was v1.5a38fbf5bf08ceb7555f2d9ebff75ff94188a3d6c0c91dc0f7ec06e760c9a988.s.14.1.\nSource s13 enumerates the current readers and deleted parser symbols. Source s14 independently projects all mapped behaviors and seams.\nCommands contributed evidence binding, source retrieval, the one diff, targeted reads, and symbol enumeration. No tests, mutations, repository writes, stash, or delegates ran.\nAn initial primary check-current failed on the evidence-store lock. The parent directed the authorized worktree route with escalation; that single worktree check succeeded.\nAn invalid --cursor s1 retrieval was rejected; the corrected --source s1 read succeeded.\nNative report only: /tmp/ft370-c1-spec.txt.\n"
          },
          "axis": "Spec",
          "base": "b0abe2fb2a4fa3764bed97896acf13d6f99fed03",
          "tip": "67f6159ba5f6df5961bee2c5e7215e43e18e01c7",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "c1-coverage",
          "performer": "/root/c1_coverage",
          "role": "independent-review",
          "model": "gpt-6.1-sol",
          "effort": "high",
          "source_digest": "fe9bad50f069fb1e7b06a867cc8c6107ac8b0441",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native:/root/c1_coverage",
            "digest": "sha256:a254edede823db454d5e456149ec0af256a09d45f72eb34e683d3d01fc8b10a8",
            "excerpt": "Coverage review: FT370 / CG-C1\nOutcome: pass. Findings: 0. Worst issue: none. Unresolved concrete concern: none.\nReviewer line: gpt-6.1-sol / high, one iteration, user-selected.\n\nBinding\nTrusted evidence: sha256:5a38fbf5bf08ceb7555f2d9ebff75ff94188a3d6c0c91dc0f7ec06e760c9a988.\nCurrent-action binding succeeded once. It reported ft370-build at 005e12afa0480ae5d17f83047839e3c4d639ef36, clean, with the requested base b0abe2fb2a4fa3764bed97896acf13d6f99fed03.\nThe reviewed implementation tip is 67f6159ba5f6df5961bee2c5e7215e43e18e01c7. The later commit records evidence. I collected one diff from the requested base to 005e12afa0480ae5d17f83047839e3c4d639ef36.\n\nIndependent input family and write fence\nThe producer is Git: TreeChanges invokes git diff --raw --no-renames -z through git.Raw, which returns stdout without trimming. The consumed family consists of immutable tree operands and their path-ordered raw entries: no changes; one or several changed paths; modifications, additions, deletions, and type changes; renames represented by separate addition/deletion entries; source/destination regular modes 100644 and 100755, missing-side mode 000000, symlink mode 120000, and gitlink mode 160000; path bytes including spaces, tabs, newlines, glob characters, and non-ASCII bytes. NUL cannot occur inside a Git path. Git-command errors and malformed metadata return errors. Bare tree IDs pass directly to Git; the gate wrapper peels only its existing base operand with ^{tree}.\nTicket 1 authorizes two implementation paths: internal/git/tree.go and internal/gate/lane_select.go. This delta writes those paths plus the spec's approved line metadata/citation update and reviews/ft370-comment-only-evidence.md. The later classifier, checkpoint, and guidance fences belong to successor chunks. The reader introduces no repository or publication writer.\n\nCG40 and CG41 audit\nCG40 explicitly has a review-owned Standards seam. I did not invent an automated coverage requirement for that row. Source inspection plus a search across internal/gate confirms the old parseComposedChange and composedChangeFields are absent, and gate.ComposedChanges calls benchgit.TreeChanges (internal/gate/lane_select.go:26). The moved parser has one implementation in internal/git/tree.go:236. The alias preserves the existing four-field input type.\nCG41 names TestComposedChangesExpandsANamedDirectory (internal/gate/lane_select_test.go:66). Its fixture creates and commits real Git trees and its independent expected rows assert both changed files, order, status, and modes. It calls the unchanged gate wrapper, so it now reaches TreeChanges and parseTreeChange. The diff contains no test edit or assertion weakening.\nThe two other preserved reader tests are TestComposedChangesRepresentsARenameAsDeletionAndAddition and TestComposedChangesCarriesTheSymlinkMode. The rename fixture uses git mv and commits it; its expectation asserts separate A/D entries with opposite 000000 sides. That assertion would distinguish a source/destination-mode swap independently of the author mutation claim. The symlink fixture creates a real symlink and expects mode 120000, with an explicit capability skip if the host cannot create it.\n\nAdversarial assessment\nI considered a newline/tab/non-ASCII path that would break a line-based reader, a rename that would add a third frame, a source/destination-mode swap that the ordinary-modification row alone would miss, lost entries after the first change, equal-tree empty output, and bare tree operands with no naming commit. The moved code retains the same -z and --no-renames flags, frame stride, append loop, four-field mapping, and direct operand forwarding. These source facts refute a concrete regression in this chunk. The existing directory, rename, and symlink tests observe the important moved-reader fields rather than a constructed fake entry.\nThe spec assigns unusual-path, mode-only, gitlink, empty-diff, and comment/token predicates to CG-C2. Their planned tests are not current-chunk omissions. The spec expressly leaves the lane reader's diff.ignoreSubmodules exposure out of scope; the move neither introduces nor expands it.\nThe producer uses immutable trees and git -C root, so transitioned spec bodies and untracked descendants cannot change the reader's subject unless a caller supplies a corresponding tree. Nested process cwd does not alter that explicit repository root. The production caller LaneAuthority.Authorize refuses when ComposedChanges errors, before any lane publication authority is returned (internal/gate/authorization/authorization.go:227). Legacy lane selection tests still consume the same alias and fields.\n\nAdvice\nNone.\n\nSources read\nFull FT370 spec, ticket 1, the requested complete diff, craft-review, finding-discipline, bounded repair policy, .bench/BENCH.md, the relevant benchkit profile seam/hostile-input sections, internal/git/tree.go, internal/git/git.go (including Raw), internal/gate/lane_select_test.go, internal/gate/lane_run_test.go, lane_select.go reader/classifier sections, and LaneAuthority.Authorize. The s13 consumer enumeration was read in full. The s14 coverage projection was retrieved to its stream end; the spec supplied its why-it-catches clauses.\nEvidence cursors fetched: s1 index 0 (stream end); s13 index 0 (stream end); s14 index 0, then cursor v1.5a38fbf5bf08ceb7555f2d9ebff75ff94188a3d6c0c91dc0f7ec06e760c9a988.s.14.1 / index 1 (stream end).\n\nCommand contribution and limits\nOne successful bench preflight evidence --check-current, one requested git diff, evidence source reads, source searches, and file reads. Bench worktree exec ft370-build -- owned repository reads. Existing spill files were retrieved through the coordinator-approved primary-CWD read route. No tests, mutation probes, repository edits, stashes, or delegates ran in this review. I did not independently rerun the author verification. This report establishes a read-only semantic Coverage result; the coordinator still owns unchanged-tree verification and the gate. Native report is the sole /tmp write.\n"
          },
          "axis": "Coverage",
          "base": "b0abe2fb2a4fa3764bed97896acf13d6f99fed03",
          "tip": "67f6159ba5f6df5961bee2c5e7215e43e18e01c7",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "CG-C2",
      "base": "67f6159ba5f6df5961bee2c5e7215e43e18e01c7",
      "tip": "e021e09c7d2f2d638a7ff9f1fed84bdbb5178c6d",
      "plan_digest": "sha256:9bc20dfc5468058ad2a0fac076f85a213a2484f0d033ecc3ed335563ab55ed65",
      "source_digest": "80ea4129830d508b4b4adf4cbee47dee5e9bb2c4",
      "acceptance_rows": [
        "CG9",
        "CG10",
        "CG11",
        "CG12",
        "CG13",
        "CG14",
        "CG15",
        "CG16",
        "CG17",
        "CG18",
        "CG19",
        "CG20",
        "CG21",
        "CG22",
        "CG23",
        "CG24",
        "CG25",
        "CG26",
        "CG27",
        "CG28",
        "CG29",
        "CG30",
        "CG31",
        "CG32",
        "CG42",
        "CG43",
        "CG46"
      ],
      "verification": [
        {
          "id": "c2-commentgap",
          "performer": "ft370-root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "80ea4129830d508b4b4adf4cbee47dee5e9bb2c4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "chat:ft370-root/c2-commentgap",
            "digest": "sha256:dd8111081c151a9f4326e006be8a3fb827e104840fac665fbf4f11a48686dbef",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft370-build,c0fd00993d4627aac1be378a69ea6ff976cd8c3a,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/commentgap,pass,319\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "commentgap",
          "command": "bench test --package ./internal/commentgap",
          "exit_code": 0
        },
        {
          "id": "c2-git",
          "performer": "ft370-root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "80ea4129830d508b4b4adf4cbee47dee5e9bb2c4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "chat:ft370-root/c2-git",
            "digest": "sha256:ccdaa68c3ae808588d12c6352f5ceb2a9553b1fd6bb9b42a6911b2949127a74a",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft370-build,c0fd00993d4627aac1be378a69ea6ff976cd8c3a,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/git,pass,1557\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "git",
          "command": "bench test --package ./internal/git",
          "exit_code": 0
        },
        {
          "id": "c2-gate",
          "performer": "ft370-root",
          "role": "author-verification",
          "model": "unknown",
          "effort": "unknown",
          "source_digest": "80ea4129830d508b4b4adf4cbee47dee5e9bb2c4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "chat:ft370-root/c2-gate",
            "digest": "sha256:4294308c9510bdd8ec618b42983dc95a07cbb46ffc6681af9776f7ea9b58d297",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft370-build,c0fd00993d4627aac1be378a69ea6ff976cd8c3a,true\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/gate,pass,11821\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        }
      ],
      "reviews": []
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  },
  "amendments": [
    {
      "from": "sha256:854ff2a41d464bc649da1d3fb4e1c9525a31ff11195d26b4cc4debcf71c9dd67",
      "to": "sha256:9bc20dfc5468058ad2a0fac076f85a213a2484f0d033ecc3ed335563ab55ed65",
      "chunk_ids": {
        "CG-C1": [
          "CG-C1"
        ]
      }
    },
    {
      "from": "sha256:9bc20dfc5468058ad2a0fac076f85a213a2484f0d033ecc3ed335563ab55ed65",
      "to": "sha256:2d9c1d57d328b461a407e5ab8bb2886d987d970101ba8352087b48a6a8adc1bf",
      "chunk_ids": {
        "CG-C1": [
          "CG-C1"
        ],
        "CG-C2": [
          "CG-C2"
        ]
      }
    }
  ]
}
```
