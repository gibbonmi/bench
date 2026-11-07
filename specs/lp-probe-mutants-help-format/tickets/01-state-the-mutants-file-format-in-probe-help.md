# State the mutants file format in probe help

Blocked by: none
Writes: internal/probe/set_input.go, internal/probe/command.go, internal/probe/refusal_test.go
Covers: none

## What to build

`bench probe --help` names `--mutants <path>`, but it does not state the format of that file. An author must guess the format.

The help notes block gets one more note. The note states that the file is one JSON array of entries, that each entry is a JSON array of strings, and that each entry holds the mutation flags of the single form and an optional `--expect <test>`, with no `<file>`. The note also gives the entry cap and one example. The note lives beside the mutants file decoder, and the cap comes from the decoder's own cap constant.

## Acceptance

- [ ] `bench probe --help`, `-h`, and `help` print the mutants file format note in the notes block, with exit 0.
- [ ] The entry cap in the note comes from the decoder's cap constant, not from a second copy.
- [ ] The help test turns red when the mutants file note is removed from the help.
