# Bounded repair policy review

## State

The implementation has one ticket and one chunk, BP-C1.
The initial Terra reviews passed; cross-harness falsification found two blockers.
Two post-review repair cycles are consumed. No reviewer extension exists.

The retained author changed both entry triggers and moved the definition to the shipped policy.
Focused verification and both omission probes pass. Current native repair coverage remains pending.

The first cycle repaired the entry trigger and glossary source.
The second cycle corrected the review record paragraph that the required prose check rejected.

## Standards

The Terra axis reported zero findings. The independent omission expectations passed its one-source review.

## Spec

The Terra axis reported zero findings after its BP1 through BP21 scenario replay.
The cross-harness pass found two repair targets. The author accepts both findings.

- BP-F1, auto-fix: The author trigger omits resumed repairs; extend both author entry pointers and their anchors.
- BP-F2, auto-fix: The glossary link is unavailable to linked consumers; keep one definition in the shipped policy.

BP-F2 needs a non-behavioral source-location adjustment under the operating guide's spec-contradiction rule.
The project glossary will point to the shipped definition. The approved cycle unit stays unchanged.
The author records this adjustment for reviewer veto and updates the spec and ticket before the repair.

## Coverage

The Terra axis reported zero findings. Its corrected return confirms the candidate identity and clean worktree.

## Optional advice

The falsification pass suggested clearer progression prose, a blank line after the review skill heading, and a delegate pointer.
These suggestions have no finding IDs, dispositions, or repair targets.
The author retains them outside the two-target repair count.

## Author verification

The baseline and changed continuation, convergence, and checkpoint regressions passed.
The workflow root and guidance budgets passed. The ticket lane passed after a pre-review paragraph correction.
The policy-text omission probe failed with the fixed-allowance diagnostic and restored cleanly.
The registry omission probe removed all 36 additions while their independent expectations remained.
It failed with one absent-anchor error per independent expectation and restored cleanly.

The system suite passed with no skips. The disposable shift replay passed and its pre-commit hook fired.
The spec retains the transport run's scope and its initial missing-adapter refusal.
Current repair verification remains pending.

## Native results

```bench-review-record
{
  "version": 1,
  "spec": "specs/bounded-repair-policy/spec.md",
  "plan_digest": "sha256:de004c846929876d3cac9611b5a41cfe05bd607ce187c1ffda1a30ae67748127",
  "implementation_session": "codex:/root",
  "chunks": [
    {
      "id": "BP-C1",
      "base": "6ce9c9bd9b269a153d970ea5859eb2643ce71861",
      "tip": "ca727952daeda7c84cb0eca38b380e42787d7fcd",
      "plan_digest": "sha256:de004c846929876d3cac9611b5a41cfe05bd607ce187c1ffda1a30ae67748127",
      "source_digest": "e023460aaf97608ff06c6ca57293e4f994f96a07",
      "acceptance_rows": [
        "BP1",
        "BP2",
        "BP3",
        "BP4",
        "BP5",
        "BP6",
        "BP7",
        "BP8",
        "BP9",
        "BP10",
        "BP11",
        "BP12",
        "BP13",
        "BP14",
        "BP15",
        "BP16",
        "BP17",
        "BP18",
        "BP19",
        "BP20",
        "BP21"
      ],
      "verification": [],
      "reviews": [
        {
          "id": "standards-initial",
          "performer": "codex:/root/standards",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "e023460aaf97608ff06c6ca57293e4f994f96a07",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/standards",
            "excerpt": "Standards — PASS\n\nExamined frozen pair `6ce9c9bd9b269a153d970ea5859eb2643ce71861..ca727952daeda7c84cb0eca38b380e42787d7fcd` for BP-C1.\n\nNo findings; no finding IDs or dispositions. Count: 0. Worst issue: none.",
            "digest": "sha256:8241b8e913253b879c1efe1ea338f6d5ceca00aadc3024c37a23e5d2896d2030"
          },
          "axis": "Standards",
          "base": "6ce9c9bd9b269a153d970ea5859eb2643ce71861",
          "tip": "ca727952daeda7c84cb0eca38b380e42787d7fcd",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "spec-initial",
          "performer": "codex:/root/spec",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "e023460aaf97608ff06c6ca57293e4f994f96a07",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/spec",
            "excerpt": "Spec **PASS** — frozen pair `6ce9c9bd9b269a153d970ea5859eb2643ce71861..ca727952daeda7c84cb0eca38b380e42787d7fcd`; worktree clean.\n\nFindings: **0**. Worst: none. Optional advice: none.",
            "digest": "sha256:61db4f9b8c89c50a98f08463eaac9b9d736fbc5722831579ae7d87efca1fe4a4"
          },
          "axis": "Spec",
          "base": "6ce9c9bd9b269a153d970ea5859eb2643ce71861",
          "tip": "ca727952daeda7c84cb0eca38b380e42787d7fcd",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "coverage-initial",
          "performer": "codex:/root/coverage",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "e023460aaf97608ff06c6ca57293e4f994f96a07",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/coverage",
            "excerpt": "Coverage — **PASS**, 0 findings; worst: none.\n\nRecheck: `bp-coverage` is clean at `ca727952daeda7c84cb0eca38b380e42787d7fcd`; `main` is `6ce9c9bd9b269a153d970ea5859eb2643ce71861`. I reviewed the candidate bytes at `ca727952…`.",
            "digest": "sha256:404028f48befc55960d186ba32838eb5aa70edb12e69cb9c33f1fa6bba06b929"
          },
          "axis": "Coverage",
          "base": "6ce9c9bd9b269a153d970ea5859eb2643ce71861",
          "tip": "ca727952daeda7c84cb0eca38b380e42787d7fcd",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "spec-falsification-initial",
          "performer": "claude:bp-falsification-initial",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "e023460aaf97608ff06c6ca57293e4f994f96a07",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "exec:49540",
            "excerpt": "Examined pair `6ce9c9bd..ca727952`, worktree `bp-falsification` clean at tip `ca727952daeda7c84cb0eca38b380e42787d7fcd`. Cap continuity, the mode census, blocking-versus-optional classification, and the completion safeguards survived falsification: each is anchored in `registry_retained_workflow.go` and covered by an independent expectation. Two findings survived.\n\n1. The author's entry trigger cannot fire on a resumed session. `.agents/commands/bench-implement-spec.md:40` and `.agents/skills/bench-craft-line/SKILL.md:108` both say \"At initial review, read [the bounded repair policy]\". These are the only author-side routes; the two state-conditioned routes (\"Before classifying repairs\") are the reviewer's. Failure scenario: a session resumes after two consumed cycles, re-enters `/bench-implement-spec`, sees a trigger tied to a moment that already passed, never opens the policy, and starts a third cycle without reading the retained count.\n\n2. The cycle definition is duplicated, and its link dangles for every consumer. `bounded-repair-policy.md:11` points at `../../../../CONTEXT.md`, then line 12 restates that definition. `CONTEXT.md:282` holds the same fact. `.bench/consumer-payload.json` ships `.agents/skills` but not `CONTEXT.md`, and `/bench-setup-repo` Section C makes the consumer's root `CONTEXT.md` an optional project glossary.\n\nOptional advice: `bench-review-implementation.md:47` narrows \"findings\" to \"unresolved blockers\" for successor progression; `bench gate --checkpoint` still refuses open findings, so this is prose ambiguity only. `bench-craft-review/SKILL.md` met its 122-line budget by deleting the blank line after its H1; condense prose instead. `craft-delegate` gets no pointer, so a delegate author relies on the orchestrator to count cycles.",
            "digest": "sha256:19665bb62e49fb49d8c4912a2c873a87aff942f30a0e27201cd051c33cb57b7f"
          },
          "axis": "Spec",
          "base": "6ce9c9bd9b269a153d970ea5859eb2643ce71861",
          "tip": "ca727952daeda7c84cb0eca38b380e42787d7fcd",
          "finding_ids": [
            "BP-F1",
            "BP-F2"
          ],
          "supersedes": [
            "spec-initial"
          ]
        }
      ]
    }
  ],
  "completion": {
    "state": "pending"
  }
}
```
