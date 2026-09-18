# Name the refused lead in the before-side refusal

Blocked by: none
Writes: internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go
Covers: none

## What to build

The follow-on guard refuses a Bench call when an earlier simple command changes the
directory or the environment. The refusal line gives one repair sentence for each
side. The before-side sentence tells the reader to run the command from the current
directory. That sentence is correct only when the refused lead changes the directory.

This ticket gives the before side two repair sentences. The directory sentence stays
for a lead that changes the directory. A second sentence applies to each other refused
lead: an assignment, an `export`, a pipe, or a `||` operator. The verdict carries the
cause, so `Message` selects the sentence from the verdict and parses nothing again.
The rules that refuse a command do not change.

## Acceptance

- [ ] A `cd` lead before a Bench call gives the directory repair sentence.
- [ ] An assignment lead, an `export` lead, a pipe, and a `||` operator each give the second repair sentence.
- [ ] The refusal line keeps its prefix, its `segment=` field, and its `operator=` field.
- [ ] Each command that the guard allowed or refused before this ticket has the same verdict.
