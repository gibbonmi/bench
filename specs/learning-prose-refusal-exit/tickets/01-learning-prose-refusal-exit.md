# `bench learning` refuses over-bound prose as a content error, not a grammar error

Blocked by: none
Writes: internal/learnings/entry.go, internal/roadmap/learning.go, internal/roadmap/learning_test.go

## What to build

A prose refusal from `bench learning` exits 1. The output leads with the prose
diagnostic and does not print the usage line. Each diagnostic names the flag
that supplied the refused sentence (`--what`, `--right`, or `--rule`). A grammar
error keeps exit 2 and the usage line. The verb writes nothing on either
refusal.

Source: the 2026-09-17 drain learning "bench learning printed its usage line for
a prose refusal". The drain reproduced the defect through `bench learning` with
a 32-word `--right` sentence: exit 2, the usage line first.

## Acceptance

- [ ] A `--what` sentence of 26 words exits 1, writes no journal, and the output
      does not contain the usage line; a test asserts it.
- [ ] The same refusal output names `--what`; an over-bound `--right` sentence
      names `--right`; a test asserts both.
- [ ] A 25-word `--what` sentence still writes the entry and exits 0.
- [ ] A missing `--right` flag still exits 2.
- [ ] `bench gate` green.
