# Write confirmations with bench record

Blocked by: 03-carry-confirmed-axis-results.md, 01-state-record-only-repair-cycle-rule.md
Writes: internal/reviewrecord/repair.go, internal/reviewrecord/record.go, internal/reviewrecord/coverage.go, internal/reviewrecord/confirmation.go (new), internal/reviewrecord/confirmation_test.go (new), internal/reviewrecord/recordtest/confirmation.go (new), internal/gate/review_checkpoint_confirmation_test.go (new), internal/reviewrecord/write.go, internal/reviewrecord/repair_test.go, internal/reviewrecord/recordcmd/confirmation.go (new), internal/reviewrecord/recordcmd/confirmation_test.go (new), internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/refusal.go, internal/reviewrecord/recordcmd/verification.go, internal/reviewrecord/recordcmd/ordering_test.go, internal/reviewrecord/recordcmd/verification_import_test.go, internal/reviewrecord/recordcmd/review_test.go, internal/reviewrecord/recordcmd/amendment_test.go, internal/gate/record_repair_test.go, internal/gate/review_checkpoint_commits_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RC38, RC39, RC40, RC41, RC42, RC43, RC44, RC45, RC58

## What to build

`bench record confirmation <slug> --chunk <id> --id <id> --performer <session> --model <model> --effort <effort> --exit-code <n> --ref <ref> --probe-result <file> [--mutant <n>] --mutation <text> [--finding <id>]...` writes one confirmation.
It imports the probe through the same path as the verification form.
It derives the reviewed pair from the current result of each axis.


It refuses when the three results sit at different pairs or already sit at the chunk pair.
It refuses a delta that `testgap.Prove` refuses, and a finding set that differs from the union of the three results' finding IDs.
That union is the same set that the checkpoint compares.
It supersedes the prior confirmation of the chunk.
It records the `--mutation` text as unverified caller text.

`FormConfirmation` joins the form constants, and the completion order places it after the review form.
The record command registers the form, and the help lists it.

Without a valid confirmation, an axis result can sit at an earlier pair of the chunk.
Its refusal names the confirmation form only when all three results sit at that pair and `testgap.Prove` accepts the delta.
That refusal comes before the outcome refusal.
Otherwise the refusal names the review form, as today.
A confirmation-validity refusal routes to `bench record 'confirmation' --help` when its three axis results sit at an earlier pair, as defined above, with a test-only delta.
Update the route assertion of each such ticket 03 row.

The writer and the stale-axis route call ticket 03's reviewed-pair function.
They compute no second pair, no second prover call, and no second finding union.

The moved-chunk case of the strict-evidence test has a test-only delta, so its route changes to the confirmation form.
Its reason text `chunk 1: stale Standards` stays.
The help rows derive from the record command, so `cmd/bench/command_registry.go` must not grow past its 400-line budget.
The over-budget registry test files must not grow.

## Acceptance

- [ ] RC38: the writer writes a confirmation from a complete `bench probe` answer with the derived reviewed pair and the reviewed tip's tree digest.
- [ ] RC39–RC41: the writer refuses axis results at different pairs, a non-test delta, and a finding set that differs from the open set.
- [ ] RC42: a second confirmation names the first in its supersession link.
- [ ] RC43: with Standards `fail` and Spec and Coverage `pass` at the reviewed pair and a test-only delta, the refusal routes to `bench record 'confirmation' --help`.
- [ ] RC44: with a delta that changes a non-test Go file, the refusal routes to `bench record 'review' --help`.
- [ ] RC58: the writer refuses a probe answer that did not bite or did not restore, through `probeBites`, and writes no entry.
- [ ] RC45: `bench record` help lists the confirmation form, and the completion order puts it after review.
- [ ] A confirmation that `bench record confirmation` writes for the RC9 state makes `bench gate --checkpoint` pass for that chunk.
- [ ] A refusal for an invalid confirmation, such as the RC20 silent probe, routes to `bench record 'confirmation' --help` at a test-only earlier pair.
