# Share structural errors with gate readers

Blocked by: 1-policy-document.md
Writes: internal/jsonfile/decode.go, internal/jsonfile/decode_test.go, internal/jsonfile/fields.go, internal/gate/verdict.go, internal/gate/json_compatibility_test.go (new), internal/gate/testdata/strict-json-baseline.json (new)
Covers: SJ16, SJ17, SJ18, SJ19, SJ20, SJ21, SJ22, SJ23, SJ24, SJ25, SJ26, SJ27, SJ28, SJ29, SJ30, SJ31

## What to build

Expose duplicate-key and trailing-data categories through errors.Is. Preserve Decode, DecodeDocument, and DecodeExactDocument semantics.
Migrate strictJSON to the owner and remove rejectDuplicateNames in this ticket. Gate adapters preserve errTrailingJSON and errDuplicateJSONName identities.
Preserve manifest unknown-field extraction and the current cache, manifest, framing, and domain policies.

Capture the accepted spec's gate baseline family before the production edit. Compare the finished callers with that permanent baseline.
Include unknown-plus-suffix and duplicate-plus-suffix inputs. Preserve their observable diagnostic precedence.

Drive a valid manifest without a final newline. Then append a second value and observe the existing trailing-content diagnostic.
Drive a duplicate key and check errors.Is against the existing duplicate sentinel.

This ticket delivers the common error contract through its first real gate caller. Later adapters consume that reviewed contract.
Remove unreachable local copies when their last caller migrates. Keep unrelated filesystem and domain code under the caller.

## Acceptance

- [ ] Document decode refuses duplicate keys inside an object in an array (SJ16).
- [ ] Document decode refuses repeated keys with equivalent escaped spellings (SJ17).
- [ ] Typed document decode refuses an unknown nested field (SJ18).
- [ ] Document decode refuses every second JSON value kind (SJ19).
- [ ] Document decode refuses malformed trailing bytes (SJ20).
- [ ] Document decode accepts JSON whitespace after its value (SJ21).
- [ ] Existing record callers retain final-newline refusal (SJ22).
- [ ] Existing document callers retain acceptance without a final newline (SJ23).
- [ ] Exact document decode refuses a case-alias key that ordinary typed decode accepts (SJ24).
- [ ] errors.Is identifies duplicate and trailing errors from the owner (SJ25).
- [ ] Gate manifest results match the committed baseline family (SJ26).
- [ ] Cache loader results match the committed baseline family (SJ27).
- [ ] Gate duplicate errors retain errDuplicateJSONName through errors.Is (SJ28).
- [ ] Gate trailing errors retain errTrailingJSON through errors.Is (SJ29).
- [ ] Gate documents valid without a final newline remain valid (SJ30).
- [ ] Unknown manifest fields retain the unknown-field key diagnostic (SJ31).

## Checkpoint verification

- `bench test --package ./internal/jsonfile`
- `bench test --package ./internal/gate`
- `bench test --package ./internal/conformance --run TestBranchNativeArchitectureCensus`

Record one red and one green for each covered row. Use a pre-edit red, a valid bench probe, or a recorded revert.
Keep each independent expectation only after its named behavioral mutation makes the assertion red. A compile failure proves no acceptance behavior.
The checkpoint commits only on the ticket lane after all focused checks pass. The orchestrator freezes the chunk for independent semantic review.

## Mutation proof

Disable the shared duplicate-memory update. The nested duplicate assertion must fail.

`bench probe internal/jsonfile/decode.go --swap 'seen[name] = true' --with 'seen[name] = false' --package ./internal/jsonfile --run TestDocumentStructuralContract`

The command is a planned probe spelling. Confirm the exact production expression after the implementation exists.
If the planned expression differs, derive an equivalent compiling mutation at the same owned seam. Record its site and mutation kind.
Require bench probe to report bit and a byte-exact restore. A silent or invalid result requires a discriminating replacement probe.
