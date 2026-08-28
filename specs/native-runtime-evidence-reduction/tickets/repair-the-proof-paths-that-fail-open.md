# Repair the proof paths that fail open

Blocked by: none
Writes: scripts/native-proof.sh, scripts/verify-release-artifact.mjs, internal/conformance/native_proof_scripts_test.go

## What to build

Two proof paths accept a state they cannot verify. Both fail closed after this
ticket.

`scripts/native-proof.sh` refuses a proven target whose operating system has no
platform predicate. Today the script sets `musl_status=not_applicable` above the
Linux guard and passes a literal `green` strip status. A proven non-Linux target
therefore skips every predicate and still mints a green proof.

`scripts/verify-release-artifact.mjs` builds its target-name set from the proven
list, not the shipped list. Today it counts against the proven list and checks
membership against the shipped list. An index carrying two Darwin proofs passes
both tests while proving neither Linux target.

The comment register in the touched files drops its narration and its
provenance. A comment states the current property, not what the change retired.

## Acceptance

- [ ] `native-proof.sh` exits non-zero for a proven target with no platform predicate (row B12).
- [ ] The refusal names the operating system it cannot verify.
- [ ] `native-proof.sh` still emits a green proof for a proven Linux target.
- [ ] `verify-release-artifact.mjs` rejects an index whose proof targets are not the proven targets (row B13).
- [ ] No comment in the touched files narrates the change or cites its provenance.
