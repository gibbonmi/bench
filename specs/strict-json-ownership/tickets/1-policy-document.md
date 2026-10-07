# Refuse concealed vulnerability policy bytes

Blocked by: none
Writes: internal/releasepreflight/vulnerability.go, internal/releasepreflight/vulnerability_test.go
Covers: SJ1, SJ2, SJ3, SJ4, SJ5, SJ6, SJ7, SJ8, SJ9, SJ10, SJ11, SJ12, SJ13, SJ14, SJ15

## What to build

Use the existing jsonfile.DecodeDocument for the complete exception array. Keep the existing empty-file refusal before decode.
Keep the findings, null, required-value, uniqueness, unused-exception, full-date, and expiry decisions under ValidateVulnerabilityPolicy.
Do not change findingIDs. Its input remains a multi-value scanner stream.

Drive [] followed by nope through the real validator with a fixed date. Then drive a used exception whose expiry equals today.
The first input must fail. The second input must pass.

This ticket delivers the defect repair while every migration successor remains unbuilt. It requires no new jsonfile API.
Its successor preserves this exception-document contract.

## Acceptance

- [ ] The validator refuses [] followed by nope (SJ1).
- [ ] The validator refuses [] followed by {} (SJ2).
- [ ] The validator refuses [] followed by 7, true, null, a string, or an array (SJ3).
- [ ] The validator refuses repeated id, reason, or expires keys in one exception (SJ4).
- [ ] The validator accepts [] with leading and trailing JSON whitespace without a final newline (SJ5).
- [ ] The validator refuses an exception with an extra typed field (SJ6).
- [ ] An absent policy with no findings returns nil (SJ7).
- [ ] A present zero-byte policy returns the empty-file refusal (SJ8).
- [ ] Present [] and null policies with no findings return nil (SJ9).
- [ ] An uncovered finding remains an uncovered-vulnerabilities refusal (SJ10).
- [ ] A used exception with expiry equal to today returns nil (SJ11).
- [ ] An exception dated before today remains expired (SJ12).
- [ ] Duplicate exception entries retain the duplicate-exception refusal (SJ13).
- [ ] An unused exception retains the unused-exception refusal (SJ14).
- [ ] Blank required values and invalid full dates retain their current domain refusals (SJ15).

## Checkpoint verification

- `bench test --package ./internal/releasepreflight`
- `bench test --package ./internal/jsonfile`

Record one red and one green for each covered row. Use a pre-edit red, a valid bench probe, or a recorded revert.
Keep each independent expectation only after its named behavioral mutation makes the assertion red. A compile failure proves no acceptance behavior.
The checkpoint commits only on the ticket lane after all focused checks pass. The orchestrator freezes the chunk for independent semantic review.

## Mutation proof

Replace the real policy bytes with an empty array at the shared decode call. The suffix-refusal case must fail.

`bench probe internal/releasepreflight/vulnerability.go --swap 'jsonfile.DecodeDocument(data, &exceptions)' --with 'jsonfile.DecodeDocument([]byte("[]"), &exceptions)' --package ./internal/releasepreflight --run TestVulnerabilityPolicyDocumentContract`

The command is a planned probe spelling. Confirm the exact production expression after the implementation exists.
If the planned expression differs, derive an equivalent compiling mutation at the same owned seam. Record its site and mutation kind.
Require bench probe to report bit and a byte-exact restore. A silent or invalid result requires a discriminating replacement probe.
