# Preserve strict and partial release schemas

Blocked by: 3-review-documents.md
Writes: internal/releaseevidence/json_validation.go, internal/releaseevidence/component_manifest.go, internal/releaseevidence/package_artifact.go, internal/releaseevidence/json_compatibility_test.go (new), internal/releaseevidence/package_artifact_test.go, internal/releaseevidence/testdata/strict-json-baseline.json (new)
Covers: SJ36, SJ37, SJ38, SJ39, SJ40, SJ41, SJ42, SJ43, SJ44, SJ45

## What to build

Migrate decodeStrict to jsonfile.DecodeDocument. Preserve typed unknown-field refusal and existing caller error context.
Replace component object grammar with shared map decode. Keep the registry-derived exact key checks under decodeExactObject.
Validate package.json structure through shared raw-document decode. Keep its partial json.Unmarshal identity projection and unknown-field acceptance.

Capture the accepted spec's complete release baseline family before the production edit. Compare the real validators with that permanent baseline.
Include null, framing, embedded registries, governance, producer envelopes, SPDX, reproducibility, native proofs, and competing invalid inputs.
Keep domain values, required newline checks, and observable diagnostic precedence unchanged.

Drive a valid artifact whose package.json contains unrelated fields. Then plant a duplicate inside one unrelated object and require artifact refusal.
Use the current artifact fixture owner. If a requirements override is necessary, use its existing restore function.

This ticket consumes the shared grammar accepted at SJ-C2. Remove rejectDuplicateJSONKeys after all its callers migrate in this ticket.
Leave no unreachable duplicate walker. The final ticket consumes this complete release migration for its ownership census.

## Acceptance

- [ ] Strict release results match the committed baseline family (SJ36).
- [ ] Strict release validators refuse recursive duplicate keys (SJ37).
- [ ] Strict release validators refuse unknown typed fields (SJ38).
- [ ] Requirement evidence retains its final-newline refusal (SJ39).
- [ ] A component object accepts exactly the registry's declared key set (SJ40).
- [ ] A component object refuses an extra or missing declared key (SJ41).
- [ ] A component object refuses duplicate keys recursively (SJ42).
- [ ] Artifact validation accepts package.json with unrelated fields (SJ43).
- [ ] Artifact validation refuses a duplicate package identity key (SJ44).
- [ ] Artifact validation refuses duplicate keys inside an unrelated package field (SJ45).

## Checkpoint verification

- `bench test --package ./internal/releaseevidence`

Record one red and one green for each covered row. Use a pre-edit red, a valid bench probe, or a recorded revert.
Keep each independent expectation only after its named behavioral mutation makes the assertion red. A compile failure proves no acceptance behavior.
The checkpoint commits only on the ticket lane after all focused checks pass. The orchestrator freezes the chunk for independent semantic review.

## Mutation proof

Disable the component extra-key refusal while preserving all required keys. The extra-key assertion must fail.

`bench probe internal/releaseevidence/component_manifest.go --swap 'if len(object) != len(fields) {' --with 'if false && len(object) != len(fields) {' --package ./internal/releaseevidence --run TestComponentJSONCompatibility`

The command is a planned probe spelling. Confirm the exact production expression after the implementation exists.
If the planned expression differs, derive an equivalent compiling mutation at the same owned seam. Record its site and mutation kind.
Require bench probe to report bit and a byte-exact restore. A silent or invalid result requires a discriminating replacement probe.
