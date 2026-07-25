# Review — session-handoff-emission

Reviewed range `5ccaf6c~1..1523522` (the FT122 feature commit plus the story-8 fix),
three axes run separately. The branch diff was empty because the work landed on `main`
before review, which is the documented happy path.

Findings marked **verified** were reproduced against the source or a built `dist/bench`
by the aggregating session, not accepted on a delegate's claim.

## Standards

7 findings. Worst: a stale duplicate of the invocable-selection rule that the story-8
fix corrected in production only — armed to fail green work.

1. **verified** — `internal/contract/runtime/runtime_handoff_facts_test.go:265` re-implements
   the invocable rule as a bare prefix test and was left at the pre-fix behavior when
   `1523522` added the compound-action exclusion to `internal/handoff/facts.go:203`. A board
   leading with `/bench-final-check / push` makes this expectation demand exactly what the
   fix forbids. AGENTS.md, one source per fact; the test-expectation exception does not
   reach it, because the copy asserts the superseded rule rather than guarding a named
   mutation.
2. **verified** — `internal/handoff/facts.go:216,222`: `short` and `plural` are copies of
   `internal/status/status.go:607,366`, and `internal/handoff` already imports
   `internal/status`. The comment at `facts.go:215` names status as the owner.
3. `internal/handoff/facts.go:173`: `boardStepSeparator = " / "` duplicates the producer at
   `internal/status/status.go:264`; the comment concedes the producer/recognizer split.
   Status should export it.
4. `internal/handoff/render.go:16` claims "adding a harness is a row here and nothing else",
   but the harness key set is also written at `internal/handoff/handoff.go:19` and
   `bin/bench.sh:307`, and `checkHarnessPrefix` guards only the prefix forms. A third harness
   ships with two stale usage lines. AGENTS.md (enforcement and its advertisement) plus
   craft-comments.
5. `internal/conformance/handoff_single_source_test.go:78` — `shapeSectionBody` is a second
   markdown section parser, re-deriving the rule owned by `internal/handoff/sections.go`
   and already weaker (fence-blind). AGENTS.md names parsers explicitly as single-sourced.
   (The delegate's line citations into `sections.go` were wrong — that file is 92 lines —
   but the duplicated-parser substance holds.)
6. `internal/contract/runtime/runtime_handoff_facts_test.go:14` and
   `internal/conformance/handoff_single_source_test.go:21` restate
   `"session-handoff.md"` as literals, in the same diff that exported
   `internal/status/handoff.go:57` as the one source of that name.
7. `internal/status/handoff.go:40` — the handoff row's action is the prose
   "rewrite session-handoff.md at HEAD" now that `bench handoff` does exactly that.
   `.bench/BENCH.md` and craft-cli require an emitted action to name what the reader can
   invoke. Side effect: the row can never be selected by `firstInvocable`.

Deferred as judgment, not filed as fixes: `facts.go` composing markdown that `render.go`
is named for; `validate` enumerating the struct's fields by hand; the empty-value guard at
`handoff.go:56` belonging in the shared grammar as an attribute; `CONTEXT.md` carrying no
entry for "pin block" while binding "Handoff" to a different meaning.

## Spec

5 findings. Worst: story 18's guarantee does not hold in the ordinary configuration, and
its fixture is built to dodge the input that moves.

All 30 coverage rows name a test that exists and passes, all 26 stories are implemented,
and the story-8 re-decision landed as re-decided (`internal/handoff/facts.go:202-212`).

1. **verified** — Story 18 asks for byte-identical output on both sinks on an unchanged
   tree. On a tracked clean clone, two consecutive runs differ: `clean tree` becomes
   `1 dirty path`, because the command's own write dirties the tree and `collect`
   re-derives the count from it (`internal/handoff/facts.go:61,90-100`). The block run 1
   prints is falsified by the write run 1 performs — the confident-wrong-fact class story 7
   exists to remove. Row 18's fixture
   (`internal/contract/runtime/runtime_handoff_document_test.go:14`) creates the handoff
   **untracked** on purpose, pinning the one input that actually moves, so the row's red
   signal bites only for timestamps.
2. `internal/handoff/facts.go:104-121` — `liveSpecs` returns every non-implemented spec
   joined with `", "`, and synthesizes `Status: unknown`. Story 6 asks for *the staged*
   spec or a statement that none is staged; neither behavior is asked for, and
   `runtime_handoff_facts_test.go:171` locks the unasked one in. Row 6 also names "an empty
   `specs/`" as its third fixture; the fixture at `:175` is a `specs/` holding an
   implemented spec, so the absent-directory case is never exercised.
3. `internal/handoff/handoff.go:56-60` rejects an empty `--next` value — unasked by any
   story, and asserted under row 22, whose behavior column credits "the shared grammar"
   that the test's own comment concedes is not what rejects it.
4. **verified** — Row 9 claims it asserts the derived line absent, but
   `internal/contract/runtime/runtime_handoff_grammar_test.go:58` checks for
   "the board's leading signal" while `internal/handoff/render.go:104` renders "the board's
   leading **invocable** signal". The assertion cannot fail. (The Standards axis found this
   independently.)
5. Row 25 and story 25 state a repo-wide scope; `checkHarnessPrefix`
   (`internal/conformance/handoff_single_source_test.go:175`) reads only non-test `.go`
   files under `internal/handoff`.

## Coverage

6 findings. Worst: silent data loss in the reviewer-owned section, on exit 0.

1. **verified** — `internal/handoff/sections.go:73` ends the State body at the first
   unfenced `## ` line, with nothing distinguishing a generated heading from one the
   reviewer wrote inside their own section. A `## Notes` inside State discards everything
   below it and exits 0. This contradicts the spec's own Won't-handle line
   (`specs/session-handoff-emission.md:396`): "preserved verbatim, so prose inside it that
   resembles a generated section is never parsed." Every other ambiguity in this splitter
   fails closed; this one fails open into the data-loss path story 19 exists to forbid.
   `sections_test.go:36` exercises only a generated heading and a fenced one.
2. **verified** — `internal/handoff/render.go:110`'s `validate` delegates to
   `toon.Representable`, which by design permits `\n`, `\r`, `\t`
   (`internal/toon/toon.go:89`) because TOON escapes them. The pin block is line-structured
   markdown, which cannot — so the predicate's contract does not match this use site.
   `bench handoff --next` carrying a newline writes a second `## State` heading and
   permanently bricks the artifact against every later run. A repo directory name carrying
   a newline splits the `Repository:` and `Path:` lines.
   `TestHandoffRefusesControlBytes` exercises only the bytes the predicate rejects.
3. **verified** — `internal/handoff/handoff.go:71` writes with `os.WriteFile` (`O_TRUNC`)
   while the package doc at `handoff.go:26-30` promises "creates freely and destroys never"
   and "nothing is written on any non-zero path". ENOSPC or a signal between truncate and
   write destroys the reviewer-owned State. `TestHandoffUnwritableTarget` chmods the parent
   directory, so it fails at `open` — the one path where the claim does hold, making it
   vacuous with respect to the guarantee. Precedent for the fix is in-repo:
   `internal/publication/record.go:89` writes to a temp path and renames.
4. A present-but-empty `session-handoff.md` is unasserted; the profile's hostile-input
   checklist names absent-vs-empty as distinct behaviors both requiring assertion. Behavior
   today is coherent (refuses, exit 1) but unpinned.
5. Backtick and glob characters in the repository path break the inline-code spans — a named
   checklist class with zero fixtures. Cosmetic.
6. A symlinked `session-handoff.md` is written through, placing the write outside the repo.
   The spec's Won't-handle covers symlinked repository *roots*, not a symlinked target.
   Undecided rather than defective.

Refuted and recorded: the prior unguarded-slice panic has no sibling in production code —
`facts.go:216` guards with `min(7, len(s))`. One unguarded index survives in test code at
`runtime_handoff_facts_test.go:256`, safe only by a format-string invariant in
`internal/status`.
