# FT232 record signal evidence

## Recommendation and scope

Obtain better evidence before specifying a repair-loop detector.
Existing records do not demonstrate the distinction required by decision ticket #5.
This report examines the Bench repository only, on 2026-09-13.

Consumed by: specs/bounded-repair-policy/decisions/ft232-repair-loop/tickets/9.md
Drift: Record producers change, or new labeled repair sequences become available.
Retire when: The advisory decision rejects the detector or replaces this evidence with a validated experiment.
Code revision: be006142df15380e8a9eb9f344f07b2ddffec42f

## Question graph

1. Which failure and identity fields do the current producers retain?
2. Which repeated signatures occur in the existing repository record?
3. Do those fields and examples distinguish stalled repairs from productive iteration?

Question 3 depends on questions 1 and 2.

## Facts: available fields

The lane records the first failing check and its first diagnostic line.
It stops selecting failures after that first check.
Source: internal/gate/lane.go:161 and internal/gate/lane.go:184.

The separate lane record omits the check and diagnostic.
Source: internal/gate/lane.go:319.

The commit span can record the assignment ID, but it discards the returned tracing context.
The lane therefore has no explicit parent relationship through that context.
Source: internal/commit/commit.go:44.

The phase span records an outcome and a skip blocker.
It does not record individual failing assertions or their ownership.
Source: internal/gate/runner.go:570 and internal/gate/runner.go:578.

The declared attributes contain no repair-attempt identity, failure ownership, complete failure set, or acceptance-progress evidence.
Source: internal/otelrecord/attributes.go:12.

The processor silently drops record-write failures.
Missing records therefore cannot prove that no attempt occurred.
Source: internal/otelrecord/processor.go:50.

## Observations: existing record

Record path: /home/mgibs/.bench/otel/bench-2826441890/traces.jsonl
Snapshot prefix: 101609090 bytes; 128981 complete lines
Snapshot SHA-256: 4e4070237d9b2fcb8c057722ab7ed9a9becb70d4ae5d6c7a13f5c9ed7e1931cf

The analysis parses each JSON line and selects spans with an end time.
It groups lane failures by their exact check and diagnostic pair.
The prefix contains 490 completed lane spans, including 51 failures.
None of those lane spans has an assignment attribute or a parent span.
The parser reports no malformed lines in this prefix.

Two exact failure pairs repeat.

| Check | Record lines | Observation |
|---|---|---|
| retro-improvement-markers | 34124, 85162, 93184 | Three distinct trees share the generic diagnostic `packages[1]{package,status,elapsed_ms}:`. |
| prose | 86818, 86850 | Both failures name the same tree and the same paragraph violation. |

The prose pair names tree b240c6feb412453a4ee92e3b017e451072a6bb64 twice.
The recorded trace IDs are 0565b8f2faece5f91c16a45b6de0c4a6 and 4454de6f9f1335633665fce4fd382ea2.
A Git object check confirms that the tree remains available.
A comparison of the two recorded trees returns an empty diff.
These observations prove a repeated check of unchanged content, not two repair attempts.

The generic diagnostic contains no defect identity.
Its three occurrences cannot establish whether the same defect recurred.
Their trace IDs are 54ea0f7bf12a66e67f29d3bf17cf9ec3, 92277d1badf91696e6d19f71975b68db, and 7f6a5efbe26a349bfb398f74888da096.

## Inferences and unknowns

Time containment can suggest a relationship between a commit and a lane.
It is not an explicit relationship and becomes ambiguous when runs overlap.
A changed tree does not prove useful progress, and an unchanged tree does not prove a failed repair attempt.

This run does not establish labeled examples of both stalled repair and productive iteration with a repeated first failure.
It therefore cannot estimate detector accuracy or select a justified threshold.
The examined history also spans producer versions; current code does not explain every historical record field.

The roadmap says no record retains a failure signature.
That statement remains true for the lane record, but not for the lane span.
The span enrichment still does not satisfy the discrimination requirement.
Sources: roadmap/FT232.md:9 and internal/gate/lane.go:184.

## Proposal

Ask the reviewer whether to authorize a bounded collection experiment.
The experiment needs explicit attempt attribution, comparable failure identities, and verified progress evidence.
Do not infer authorization to expand collection from this report.

## Verification record

- The report separates source facts, record observations, inferences, and a proposal.
- The analysis reads one fixed record prefix and identifies its digest.
- The repeated signature table enumerates both repeated exact pairs in that prefix.
- The prose example checks the recorded tree through Git.
- No detector runs, and no compatibility claim or accuracy claim is made.
- The report retains the missing labels and attribution as explicit unknowns.

## Validation plan

After reviewer approval, collect a bounded set of repair sequences with explicit attempt attribution.
Include stalled repairs, productive repairs, unchanged reruns, and overlapping assignments.
Require independent evidence for the progress label of each sequence.
Evaluate the candidate signal against those labels before detector specification.
