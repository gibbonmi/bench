# Repair collection pilot

This repository-only pilot records bounded repair evidence for a later reviewer decision. It does not enable a detector, change a required check, emit a warning, select a model, or affect the gate.

## Activate collection

From the Bench kit repository, start one 14-day collection window:

```sh
bench repair-pilot activate
```

Activation is idempotent. Collection stops at the earlier of the 14-day deadline or ten completed source/spec/chunk sequences. Do not activate the pilot merely to test this guide.

## Record verification evidence

At each pre-review or post-review verification point, write one versioned JSON record to a file and import it:

```sh
bench repair-pilot record --input <input.json>
```

An observation identifies its source, spec, chunk, assignment, session, source revision, stage, and native evidence. A sequence begins with a blocking failure observation. Later records may describe a repair attempt, an unchanged rerun, or a verified endpoint. Preserve producer-native references to the exact command output, review finding, diff, or handoff used as evidence. When reporting assignment timing, include both interval bounds and a native interval reference; unproven, partial, adjacent, and zero-length intervals do not establish overlap.

A repair observation may propose `productive`, `stalled`, or `unchanged`, but a proposal is not an accepted label. Before the label can count, a separate reviewer session must inspect the referenced native evidence and import an audit record. The auditor records the repaired target and a verification reference. Conflicting conclusions remain unknown until a later audit explicitly resolves the cited audit IDs.

## Synthetic fixture — not pilot evidence

The following input demonstrates the record shape only. It must never be imported or counted as a real pilot example.

```json
{
  "version": 1,
  "observation": {
    "id": "synthetic-first-failure",
    "observed_at": "2026-09-03T12:00:00Z",
    "sequence": {
      "source": "synthetic-source",
      "spec": "synthetic-spec",
      "chunk": "synthetic-chunk"
    },
    "assignment_id": "synthetic-assignment",
    "session_id": "synthetic-session",
    "source_revision": "synthetic-revision",
    "stage": "pre-review",
    "kind": "failure",
    "failure_completeness": "complete",
    "failures": [
      {
        "check": "unit",
        "identity": "RP-SYNTHETIC",
        "diagnostic": "synthetic blocker",
        "ownership": "diff-owned",
        "blocking": true,
        "reference": {
          "producer": "synthetic-fixture",
          "native": "native:synthetic-failure"
        }
      }
    ],
    "references": [
      {
        "producer": "synthetic-fixture",
        "native": "native:synthetic-observation"
      }
    ]
  }
}
```

An audit import uses the same command and this top-level shape:

```json
{
  "version": 1,
  "audit": {
    "id": "synthetic-audit",
    "observation_id": "synthetic-repair-observation",
    "auditor": "independent-review-session",
    "evidence_references": [
      {
        "producer": "synthetic-fixture",
        "native": "native:synthetic-review-evidence"
      }
    ],
    "reason": "reviewed the repair and its native verification evidence",
    "conclusion": {
      "label": "productive",
      "repaired_target": "RP-SYNTHETIC",
      "verification_reference": {
        "producer": "synthetic-fixture",
        "native": "native:synthetic-verification"
      }
    }
  }
}
```

## Inspect and export

During collection, inspect the provisional summary:

```sh
bench repair-pilot report
```

After cutoff, the default report shows completed and incomplete sequence counts, unknown labels, and evidence gaps. It also shows coverage for stalled repairs, productive repairs, unchanged reruns, and overlapping assignments. Missing required classes make the terminal sample inconclusive. Complete class coverage says only that the sample contains all four examples; it makes no accuracy, threshold, or adoption claim.

Export the stable full report for the eventual review:

```sh
bench repair-pilot report --full > capture/reports/repair-collection-pilot.md
```

The full form retains observations, proposals, audits, citations, intervals, comparisons, and evidence gaps. The reviewer inspects `capture/reports/repair-collection-pilot.md` in a later, separately authorized phase and decides whether the collected evidence supports any detector proposal. The implementation phase neither activates the real pilot nor creates that empirical report.
