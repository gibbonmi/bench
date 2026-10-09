# Refuse a symbolic-link operand in gate-prose

Blocked by: none
Writes: specs/gate-prose-symlink-operand/tickets/01-refuse-a-symbolic-link-operand-in-gate-prose.md (new), internal/gate/gate_prose.go, internal/gate/gate_prose_test.go, internal/prose/subject.go, internal/prose/subject_test.go, internal/prose/prose_test.go
Covers: none

## What to build

The `bench gate-prose <root> -- <path>` form grades each named path. At this time, the prose package skips a named symbolic link without a grade and without a refusal. The verdict table then gives the link a `pass` row, and the verb exits 0. A reader can think that the gate graded the link.

A named symbolic link must never get a `pass` row. The verb must refuse the link with one diagnostic line and exit 1. The whole-tree walk already refuses a symbolic link to a file. The named form must use the same classifier, so the two forms state the same refusal.

The prose package makes the decision to grade, exclude, or refuse a subject. The verb must not do a second symbolic-link check. The staged form selects only regular files from the index, so a staged symbolic link gets no row. This ticket keeps that rule.

## Acceptance

- [ ] `bench gate-prose <root> -- <link>` prints a `refused subject` line that names the link, and it exits 1. This is true for a link to a file and for a dangling link.
- [ ] The output for a named link holds no `prose[` table and no `pass` row, also when a clean path is named beside the link.
- [ ] `bench gate-prose <root> --staged` gives no row to a staged symbolic link.
- [ ] A named path that does not exist keeps its `refused unreadable subject` line and exit 1.
- [ ] A recorded probe that restores the symbolic-link skip makes the focused test fail.
