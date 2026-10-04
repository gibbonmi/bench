# FT360 planning reviews

The user requested two checkpoints. GPT-6.1 Sol/high reviews the spec before slicing. GPT-5.6 Sol/high reviews the slices afterward.
The author remains in the invoking session. Reviews are read-only, with two rounds available at each checkpoint.

## Spec review

Round 1 examined spec digest `971f9706c01aace676899d12361c1f361e3ff4d5b31f02bc44148c531829258e`.
It found two blocking coverage gaps. Helper tests did not preserve each caller's capability class. The commit rows did not refuse an implicit empty commit.
The author added GF31, GF32, and GF33. The author also specified a tracked-file replacement for the FIFO refusal fixture.

Round 2 examined spec digest `fd44597fe9b61a80d2eedcc5ef4e93c17536b2708bebad9addca16d324bce156`.
The reviewer accepted the repaired spec for slicing. Its judgment was claimed with confidence 9. It performed no source mutations or implementation tests.
The author folded its nonblocking table-format correction. Both review passes used GPT-6.1 Sol/high.

The spec acceptance does not approve implementation. The final slice review and the user's sign-off remain separate checkpoints.

## Slice review

Round 1 examined commit `135e41c19bdc6a164a021740deb01eecf8ff5a85` on GPT-5.6 Sol/high.
The reviewer found one blocking size defect. GF-C5 combined independently useful package families with specialized probe classification.
The author split it into GF-C5A through GF-C5D. Each chunk has its own package checks and preservation evidence.

The reviewer found no other caller, dependency, checkpoint, or fence blocker. The confirming pass accepted the repair.
The inventory representation was compacted without changing its parsed data.

Round 2 examined commit `5c7ba80dd45823a6fde13d7a53faa059ac5fb15b` on GPT-5.6 Sol/high.
The reviewer accepted the package clusters, serial dependencies, unchanged fence union, and matching checkpoint commands.
Its judgment was claimed with confidence 9. It found no remaining blocker and reported a clean tree.
The author verified the clean source and the unchanged inventory data. The final plan has sixteen tickets and thirty-three acceptance rows.

## Final planning validation

The tested commit is `8148f91e92bd9ab0ad6370592a665e634082c7f9`. Its prose lane passed, and its tracked tree was clean.
Build preflight reported thirteen green checks, two not applicable checks, and no red check. The coverage map validated all thirty-three rows.
Twenty-three rows name planned implementation checks without current seam citations. This staged plan does not claim those future checks have passed.

Two full gate runs failed at the checkout guard. Both runs passed formatting, vet, package tests, race tests, system tests, and shell checks.
The guard reported changed generated files: `bin/bench-broker.manifest`, `dist/bench`, and `dist/bench.seal`.
The second run had no concurrent command in this worktree. No tracked path changed. These runs did not identify the writer.

The retained gate logs are `.logs/gate-20261004T101210.114634570Z-2964979.jsonl` and `.logs/gate-20261004T101927.537060541Z-3266510.jsonl`.
This record preserves a validation limit. It does not claim a green full gate or authorize a landing.
The reviewed plan remains staged for the successor milestone. Implementation and the generated-artifact failure remain separate from this planning deliverable.

## Landing diagnosis

The user authorized landing the planning artifacts on 2026-10-04. The first prospective landing reproduced the generated-file failure and refused publication.
The author narrowed the writer to the root conformance check named `package-core-guard`.
The shell selected Node 18.19.1 and npm 9.2.0. That npm version ran the package's prepare script despite `--ignore-scripts`.

A separate minimal package reproduced the same lifecycle error. Its prepare script wrote a marker during a dry-run pack with scripts disabled.

The installed Node 25.8.1 runtime includes npm 11.11.0. The same minimal package skipped prepare under that runtime.
The root conformance check then passed and left all three generated files unchanged.

The landing uses that installed runtime by prepending its binary directory to PATH. No gate check, repository code, or global installation changed.
The landing's retained gate verdict determines publication. The earlier failed runs remain part of the evidence.
