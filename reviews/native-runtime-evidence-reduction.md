# Review pickup — native runtime evidence reduction

Frozen base `c59b14c8`. Reviewed tip `b2d84316`.

Raw findings: 13. De-duplicated repair targets: 9.

## Standards

Findings: 6. Worst issue: `scripts/native-proof.sh` mints a green strip status
for a proven non-Linux target with nothing behind it.

- **S1 — the proof script fails open.** `scripts/native-proof.sh:52` sets
  `musl_status=not_applicable` above the Linux guard. Line 79 passes the literal
  `green` for `strip_status`. A proven non-Linux target skips every platform
  predicate and still mints a green proof. `release-plan.mjs` accepts
  `native_proof: true` on a Darwin row, and `release_plan_test.go` certifies that
  row as supported. Disposition: `ask-user`.
- **S2 — a check grades a wording.** `internal/conformance/native_proof_scripts_test.go`
  reds the gate when `native-proof.sh` contains `Mach-O`, `nm -a`, `Darwin binary`,
  or `darwin-symbols`. The spec's own Won't-handle line says a re-proved Darwin
  target re-authors that assertion. The check therefore forbids the sanctioned
  repair. Disposition: `ask-user`.
- **S3 — the runbook re-derives the plan.** `docs/release-runbook.md` names the
  plan as the one source, then enumerates four shipped targets and two proven
  ones. `AGENTS.md` treats two derivations of one fact as a defect.
  Disposition: `ask-user`.
- **S4 — a test re-derives the plan table.** `internal/releaseevidence/release_plan_test.go`
  hardcodes all four target rows and their flags. The decode assertion alone
  serves the stated purpose. No demonstrated red justifies the enumeration.
  Disposition: `auto-fix`.
- **S5 — comment register.** Six new comments narrate the change or cite its
  provenance. The archive-root explanation is written out four times.
  Disposition: `auto-fix`.
- **S6 — a transport row decoder now exists in four places.** Go's package
  boundary makes the collapse cost an exported symbol. Disposition: `no-op`.

## Spec

Findings: 2. Rows B4, B5, and B10 are graded by something weaker than the row
claims. The other 26 rows hold. Worst issue: one token reverses the central
decision with the gate still green.

- **P1 — rows B4 and B5 grade the consumer only.** The conformance check asserts
  the job's `fromJSON(needs.preflight.outputs.proven)` reference. Nothing anchors
  the producer step to `proof-matrix-json`. Change that command in
  `.github/workflows/native-runtime.yml` and the gate stays green while both
  macOS runners start. Disposition: `auto-fix`.
- **P2 — row B10 has no red.** `internal/releaseevidence/release_evidence.go`
  reports `native target proof is incomplete`. That string occurs once in the
  tree, at its definition. The probe mutates fields inside an existing proof and
  never removes one. Disposition: `auto-fix`.

## Coverage

Findings: 5. Worst issue: the `proven` matrix producer is ungraded.

- **C1 — the producer is ungraded.** Same fix as P1. Deleting the output and its
  step also stays green, and then the job cannot start. Introduced.
  Disposition: `auto-fix`.
- **C2 — the verifier checks the count, not the membership.**
  `scripts/verify-release-artifact.mjs` builds `targetNames` from the shipped
  list. An index carrying two Darwin proofs passes the count and the membership
  test while proving neither Linux target. Introduced. Disposition: `ask-user`.
- **C3 — row B10 has no red.** Same fix as P2. Inherited, closed falsely by the
  map. Disposition: `auto-fix`.
- **C4 — the proof record shape has four underived copies.** No test executes
  the emitter. Rename an emitted key and every test stays green. Inherited, and
  the Linux branch is now the script's whole body. Disposition: `ask-user`.
- **C5 — a proven non-Linux target fails open.** Same fix as S1. Introduced.
  Disposition: flag only.

## Repair targets

| id | fix | axes | disposition |
|---|---|---|---|
| R1 | Anchor the matrix producer steps to their commands | Spec P1, Coverage C1 | auto-fix |
| R2 | Add a probe mutation that removes one proven proof | Spec P2, Coverage C3 | auto-fix |
| R3 | Make the proof script refuse a proven non-Linux target | Standards S1, Coverage C5 | ask-user |
| R4 | Replace the token ban with a property check | Standards S2 | ask-user |
| R5 | Build the verifier's target names from the proven list | Coverage C2 | ask-user |
| R6 | Stop re-deriving the target lists in the runbook | Standards S3 | ask-user |
| R7 | Drop the enumeration from the plan decode test | Standards S4 | auto-fix |
| R8 | Repair the comment register and the repeated explanation | Standards S5 | auto-fix |
| R9 | Give the proof emitter an execution seam | Coverage C4 | park |
