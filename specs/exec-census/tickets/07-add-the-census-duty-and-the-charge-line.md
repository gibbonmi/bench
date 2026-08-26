# Add the census duty to the final check and the charge

Blocked by: 06-carry-and-drop-the-census-across-a-landing.md
Writes: .agents/commands/bench-final-check.md, .agents/skills/bench-craft-delegate/SKILL.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, .bench/BENCH-reference.md, projects/benchkit.md, CONTEXT.md, CHANGELOG.md

## What to build

The guidance lands last, so each account names the shipped behavior.
`.claude/commands` is a symlink to `.agents/commands`, so one edit serves
both harnesses; write nothing under `.claude/`.

The post-merge tail of `.agents/commands/bench-final-check.md` gains the
census duty. The duty reads `census=<n>` from the landed record the phase
already holds. For `n > 0` it writes exactly one `bench learning --rule`
entry per landing.

The entry's title names the assignment label and `n`. Its
`--what` lists each verb head with its count. Its `--right` names the Bench
form per head, or `none`. Its `--rule` proposes the verb or the help change
the drain should open.

For `n = 0` the close states `census: 0 raw calls`.
The duty is advisory: a nonzero count never blocks a landing and never reds
the gate. The retro section states that a spec retro cites the landing's
census entry under `### Bench CLI` with its `Feeds:` line.

Two sentences in that file are anchor needles, each on one physical line:

- `Read census=<n> from the landed record; for n > 0, write exactly one bench learning --rule entry for the landing.`
- `For n = 0, state census: 0 raw calls in the close; a nonzero count never blocks a landing.`

Write each needle with its code spans (`census=<n>`, `n > 0`, `bench
learning --rule`, `n = 0`, `census: 0 raw calls`) and register the exact
text.

The charge section of `.agents/skills/bench-craft-delegate/SKILL.md` gains
this needle:

- `Ask the delegate for zero to two Bench CLI improvements derived from its own calls, and fold them into the landing's census entry.`

The file sits at its prose budget in `projects/benchkit.md`, so the sentence
absorbs into an existing paragraph without a new line. A raised budget is a
reviewer decision this ticket does not take.

`internal/anchors/registry_data.go` gains `Require` needles in the
`AfterImplementSpec` group beside the existing final-check rows:

- the duty needle
- the zero-close needle
- the charge needle
- one census sentence each for the reference and the profile

Each needle stays on one physical
line and carries its own diagnostic. A test in
`internal/anchors/registry_data_test.go` removes each needle's text from a
copy of its file and proves the registry reports the diagnostic.

`.bench/BENCH-reference.md` and `projects/benchkit.md` each gain this needle
sentence in their signal accounts, with `census` and the path as code spans. The needle reads:

- `The census signal counts raw calls per assignment from $BENCH_HOME/census/<repo-key>/.`
 `CONTEXT.md` already lists `census` in the signal
enumeration after ticket 05 and defines `census`, `raw call`, and `verb
head`. Add the records' home to the `census` entry and change nothing else.
`CHANGELOG.md` gains one entry under `## [Unreleased]` that states the new
signal, the landed key, and the duty.

## Acceptance

- [ ] The final-check command carries the duty needle, the four learning fields, the once-per-landing rule, and the retro citation. (EC27)
- [ ] The final-check command carries the zero-close needle with the advisory rule. (EC29)
- [ ] The delegate skill's charge section carries the charge needle, and the file stays within its prose budget. (EC31)
- [ ] A copy of the final-check command without the duty needle makes the registry report its diagnostic. (EC32)
- [ ] The reference and the profile each carry the census sentence needle, and the `CONTEXT.md` enumeration names `census`. (EC33)
- [ ] `CHANGELOG.md` names the signal, the landed key, and the duty under `## [Unreleased]`.
