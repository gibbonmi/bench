# Land the example-agreement conformance check

Blocked by: export-the-ticket-parser.md, teach-the-executable-contract-template-and-examples.md
Ownership fence: `internal/conformance/example_agreement_test.go`, `internal/conformance/registry/registry.go`, `internal/conformance/checks_test.go`, `tests/canary/example-agreement`
Assumptions: the marked region is authored by the template ticket as `<!-- ticket-example:begin -->` and `<!-- ticket-example:end -->` at column 0 immediately outside the fenced block — begin above the opening fence line and end below the closing one — so both fence lines fall inside the region for this extractor to require and strip; heading depth inside the example is `###`; `TestRegistryBindsEveryCheck` fails unless the registry row and the `checks_test.go` binding land together; `familyChecks` in `registry/registry.go` binds a canary family directory to its check with `data-handling-derivation` as the precedent; `checkSkillsIndexGenerateVerify` is the posture model for a check that writes only its own temp directory and diagnoses instead of panicking; the exported specbuild parse takes paths rather than bytes. Re-derive from the tree at pickup.

## What to build

FT164's proof that the taught shape and the accepted shape cannot drift apart: a
conformance check that reads the `craft-tickets` Good example and parses it with
the real specbuild parser, so an example that looks right but does not assign
turns the gate red.

The check extracts the region between the `ticket-example` begin and end markers,
requires the opening and closing fence lines inside it, strips them, and fails
loud when either fence line or either marker is absent, when the markers are
duplicated, or when they enclose an empty block — stricter than the passlist
precedent, which misses duplicates and would grade a second block happily. It
materializes the stripped block as a real `tickets/`
file beside a temp spec path (`os.MkdirTemp` with a deferred removal, mirroring
the skills-index check's diagnose-never-panic posture), runs the exported parse,
and asserts independently authored literals: at least two acceptance rows with
distinct ids in the taught grammar, the exact multi-entry fence, the exact
assumptions, a `Blocked by:` line present, and a `## Red mutations` section
naming every acceptance id. Their independence is what lets the named omission
— the example drifting off the parseable shape while its prose stays plausible
— turn the gate red.

The registry row and the `checks_test.go` binding land together, since the
binding test refuses a half-landed check. Bite proofs sit in the same file:
in-process mutations for the wrapped fence, the wrapped assumptions
continuation, an unlabeled row, and each marker malformation, plus an
end-of-file fixture with no trailing newline. One canary family bound in
`familyChecks` follows the data-handling-derivation precedent.

## Acceptance

- [ ] [EA1] an absent begin marker, end marker, opening fence line, or closing fence line fails the check with a named diagnostic rather than grading nothing.
- [ ] [EA2] duplicated markers and a present-but-empty marked block each fail the check with their own diagnostics.
- [ ] [EA3] the marked example parses to at least two distinct-id rows, the exact fence, the exact assumptions, a blocked-by line, and a mutations section naming every acceptance id.
- [ ] [EA4] a wrapped fence or a wrapped assumptions continuation in the example turns the check red.
- [ ] [EA5] an unlabeled acceptance row in the example turns the check red.
- [ ] [EA6] a marked block ending at end-of-file with no trailing newline parses identically to the same block followed by a newline.
- [ ] [EA7] the check is registered and bound in the same diff, and its canary family is bound in `familyChecks`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| EA1 | return no diagnostics when the begin marker or the opening fence line is missing | `TestExampleAgreementMarkersFailClosed` | strip each of the four in turn in a temp-tree copy of the skill file, run the check, expect the named diagnostic for each |
| EA2 | grade the first marked block and ignore the rest | `TestExampleAgreementMarkersFailClosed` | duplicate the marker pair, and separately empty the block, run the check, expect each diagnostic |
| EA3 | compare the parse against values read back from the same block | `TestExampleAgreementParsesAuthoredLiterals` | run the check over the tree, then mutate one fence entry in a temp-tree copy of the skill file and expect the mismatch diagnostic |
| EA4 | accept a fence split across two lines | `TestExampleAgreementRejectsWrappedFields` | wrap the fence, then the assumptions, run the check, expect the red for each |
| EA5 | relax the row assertion to any `- [ ]` bullet | `TestExampleAgreementRejectsWrappedFields` | drop the `[ID]` label from one row, run the check, expect the row-grammar diagnostic |
| EA6 | trim the block at the last newline before parsing | `TestExampleAgreementEOFWithoutNewline` | parse the block with and without a trailing newline, compare the two parses field by field |
| EA7 | register the check without binding it | `TestRegistryBindsEveryCheck` plus the canary family run | remove the `checks_test.go` binding, run the binding test, expect the unbound-check failure |
