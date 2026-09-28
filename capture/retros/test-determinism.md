## Outcome

Test determinism landed at 76ad133c6b852472abf7ea39cdd28a22df1c2b38 on 2026-09-28.
All 51 acceptance rows are covered, and all current review axes pass.
The reviewed source runs from 72e751b7a020a76aa35286852b4ea981c95cb865 through fbd8e5faee88c458ecc0336bbd6e6ff88a71cdd0.
The spec and review record remain in that history.
Retirement landed at 6625f5b98fa8374c0b1e7714f44c0c2088cb0ee7.

Kit test entry points now own private environments and preserve the operator cache selection.
Build fixtures use private kit copies, and the gate detects live-checkout writes.
Named wait policy distinguishes verdict windows from fixed windows and preserves test-owned limits.
Linked repositories retain their environment.

The implementation landing passed all six gate phases.
It reported eight capability skips and no environment skips.
Its census recorded 744 raw calls: python3 723, rg 19, and python 2.
Retirement recorded one python3 call.
The local native evidence is preserved in .logs/test-determinism-run-20260928.

## Gate-stage timings

These are the implementation landing measurements from gate-20260928T175653.655561437Z-2729292.
The scaffold's newer retirement measurements describe a separate landing.

| phase | implementation milliseconds | retirement milliseconds |
| --- | --- | --- |
| gofmt | 124 | 123 |
| vet | 1171 | 1199 |
| test | 136956 | 136970 |
| race | 2925 | 2928 |
| system | 42780 | 42582 |
| shellcheck | 527 | 528 |

## Ticket-versus-spec-slice and delegate performance

The user retained one inline author and delegated only reviews to gpt-6-sol at high effort.
The implementation model identifier is unknown.
Eight tickets completed across four accepted chunks.
Ten repair rounds followed initial review boundaries.
The user removed the repair cap while bench-debug governed repairs.

Thirty-nine independent axis returns cover the initial reviews and their confirmations.
Two additional Coverage dispatches failed before a reviewer started.
Every current axis passes on its frozen source.
Reviewers read source and retained evidence; the inline author executed all tests and mutations.
Usage, provider cost, and comparative latency remain unknown.

| ticket | slice result |
| --- | --- |
| 1 | Private run ownership passed after socket-path and local diagnostic corrections. |
| 2 | Runner composition passed after safe reading, Go stub forwarding, and telemetry cleanup repairs. |
| 3 | Conformance probes use isolated homes and the approved npm cache exception. |
| 4 | Private kit copies and build fixtures passed; later review required a commit-count mutation. |
| 5 | Named-check timing assertions grade the private root. |
| 6 | Live-checkout metadata guards passed their write, ordering, and administration-path probes. |
| 7 | Verdict-window switching passed with cancellation and fixed-window behavior preserved. |
| 8 | Wait enforcement passed after syntax-family repairs and the omitted Git guard fixture adaptation. |

The following labels preserve the original confidence.
Finding labels use the accepted auto-fix dispositions and recorded coordinator reproductions.
The invalid first dot-import example is excluded; the compiled replacement and its biting probes support the same resolver finding.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| Whole build | Expected repair rounds equal 1. | claimed | 7 | refuted: attribution table totals 10 | unknown / high / author and coordinator |
| TD-C1b-S1 | The discovered policy reader lacks safe classification. | claimed | 8 | held: accepted auto-fix and recorded special-file regression | gpt-6-sol / high / Standards |
| TD-C2-S1 | The one-commit expectation lacks a demonstrated red. | claimed | 8 | held: accepted auto-fix; probe 86295 reaches that assertion | gpt-6-sol / high / Standards |
| TD-C3-S1 | The discovered Go reader lacks special-file rejection. | claimed | 9 | held: accepted auto-fix; FIFO reproduction and classifier probe | gpt-6-sol / high / Standards |
| TD-C3-C1 | Several production wait forms evade the checker. | claimed | 8 | held: accepted auto-fix and registered-owner canary probes | gpt-6-sol / high / Coverage |
| TD-C3-C2 | Multiple dot imports can suppress a wait diagnostic. | claimed | 8 | held: corrected compiled fixture; probes 74637 and 34838 | gpt-6-sol / high / Coverage |

Brier mean: 0.11 over 6 labeled pairs.
The author forecast contributes 0.49; the five finding pairs contribute a mean of 0.034.
The remaining 34 axis returns lack an independent claim-level semantic label and remain abstentions.
Passing tests alone do not label a broad review judgment.

## Coordinator catches

The original behavior requirement covered every production wait.
Its coverage inventory did not enumerate reassignment, elapsed comparisons, injected sleepers, multiple dot imports, sibling aliases, or parentheses.
Later coverage-prose rules could have exposed that weak inventory during spec review.
The intended behavior required no change.

The final C3 checkpoint found an additional original census omission.
The real Git guard timeout fixture depended on a production window that the test switch removed.
A raw test setter and the unchanged fixture assertions close TD51.
The complete Git guard package now belongs to the required verification inventory.

The coordinator rejected an invalid Go example before accepting its claimed mutation evidence.
A compiled replacement demonstrated the resolver defect and restored the source.
The coordinator also required an assertion-specific mutation for the one-commit expectation.
Safe special-file reading remained a separate implementation standards defect.

The broker rehearsal repaired its manifest but pointed the installed shim at the temporary source.
After source release, bench status exited 127.
The primary launcher doctor repair restored the shim, and the same command passed.
No project source changed during that installation repair.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-open-kit-test-run | 2 | other; other |
| 2-compose-kit-test-run | 3 | other; check-gap; other |
| 3-isolate-conformance-probe-home | 0 | none |
| 4-run-build-scripts-on-kit-copy | 1 | other |
| 5-grade-named-check-on-private-root | 0 | none |
| 6-guard-live-checkout | 0 | none |
| 7-switch-verdict-windows | 0 | none |
| 8-name-every-production-wait | 4 | other; spec-row; spec-row; check-gap |

The table totals 10 rounds.
Ticket 1 covers the socket-path repair and local diagnostic formatting.
Ticket 2 covers safe policy reading, a missing Go fixture check, and telemetry child cleanup.
Ticket 4 is an evidence-only compliance repair.

Ticket 8 covers unsafe reading with missing wait forms, import-resolution coverage, parenthesized forms, and the omitted Git guard package.
The final check-gap originated in the spec's switch census, not a changed behavior decision.
The two spec-row rounds also expose incomplete input-family coverage.
In-pass compile, prose, and fence corrections are recorded separately and do not add repair rounds.

## Agent-experience improvements

### Bench CLI

- Add a bounded artifact reader with escaped page budgets and early pagination grammar, as captured in the learning "test-determinism census 744 raw calls".
  Feeds: new
- Expose the final assignment census before landing, as captured in the learning "test-determinism-retire census 1 raw call".
  Feeds: new
- Keep broker rehearsal separate from the durable shim target and verify the installed command after source release.
  Feeds: new

### Skills

- Apply the current coverage input-family census to older staged specs before their first implementation charge.
  Feeds: none
- Compile reviewer Go examples before treating their mutation failures as behavioral evidence.
  Feeds: none
- Persist each repair supplement before its first edit and read every required evidence page before action.
  Feeds: none

### Process

- Expand both ticket and spec fences with the pinned fixture closure before running an enabling plan preflight.
  Feeds: none
- Keep independent expectation probes specific to the assertion they justify.
  Feeds: none
- Include every production-window composition fixture in the switch census and required package checks.
  Feeds: none
