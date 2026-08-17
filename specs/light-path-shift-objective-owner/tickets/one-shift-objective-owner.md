# One shift objective owner

Blocked by: none

Ownership fence: `internal/shift`
Assumptions: historical decision source `decisions/deepening-2026-08.md` at commit `eda60dfe` #4 (exit-test rule) and #10 (candidate 5); `sanitize.Preview` is the one escape-and-bound policy; `objectiveBanner` and `validateObjective` keep their tested signatures

## What to build

The shift objective reaches five surfaces with five treatments: the banner goes
through `sanitize.Preview`, the scratch file `.bench-objective` takes the raw text
under mode 0600, the iteration prompt and the `done.sh` predicate argument take the
raw text, and the durable commit subject `shift: iteration <i> — <objective>` takes
the raw text into history. Nothing owns the objective; `loop.go` and `shift.go` each
decide a treatment where they happen to use it.

Give the objective one owner: an `objective` type in `internal/shift` that
`validateObjective` admits and that hands out one projection per surface — banner
line, prompt body, scratch bytes, predicate argument, commit subject. `loop.go` and
`shift.go` obtain every rendering from it and never spell a treatment themselves.

The one behavior change (ASSESSMENT C-08, map #10): the durable commit subject
becomes `shift: iteration <i> — <sanitize.Preview(objective)>` — escaped, 120-rune
bounded with the `… (N bytes)` suffix, the same policy as the banner. Prompt, scratch,
and predicate argument stay verbatim; the scratch file stays 0600.

Refactor exit test (map #4): the pre-existing suite passes with test logic unmodified;
mechanical renames are the only permitted test edit, and a changed assertion reverts
the move.

## Acceptance

- [ ] [SO1] Every rendering of the objective in `loop.go` and `shift.go` comes from the `objective` type's projections; no call site formats, escapes, or writes the text itself.
- [ ] [SO2] The commit subject for an objective longer than 120 runes or carrying an escape sequence is `shift: iteration <i> — <sanitize.Preview(objective)>`, and the banner shares that projection policy.
- [ ] [SO3] Prompt, scratch bytes (`<objective>\n`), and predicate argument are byte-identical to today's; the scratch file mode stays 0600.
- [ ] [SO4] Every pre-existing test in `internal/shift` passes with test logic unmodified.
