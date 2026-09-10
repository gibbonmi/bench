## Outcome

The FT311 landing-completion spec landed on `main` at `f50671018ec41711a89f8d73c0ae28df8af6fc71` on 2026-09-10, from source tip `d67a956f7d4374eb2462521e7f22ac368292a44f` over base `265bce07aafc93667d086628927b1b893f90297c`. The landing gate ran green on tree `32a9465c`, and the spec flipped to `Status: implemented` in that commit.
The build delivered the two landing effects, `bench handoff --state-file`, `bench retro --scaffold`, the seam-record reader, the landing trace threading, and the guidance fold. Six tickets landed: the five reviewed tickets and one repair ticket. The old installed broker ran the landing, so its effects row first printed on the resume that followed the sanctioned rebuild. That resume reported `refresh,complete` and `cleanup,complete`, and it retired the two ticket siblings and the re-review worktree.

## Gate-stage timings

The landing's own gate, from its stdout on 2026-09-10:

| stage | elapsed ms |
|---|---|
| gofmt | 116 |
| vet | 914 |
| test | 71748 |
| race | 2389 |
| system | 23654 |
| shellcheck | 499 |

The scaffold printed `commit unknown` and no stages. The newest landing span in the record was the other session's resume, and a resume carries no subject and runs no gate. The scaffold's rule selects that span by design, so the timings above come from the landing's stdout.

## Ticket-versus-spec-slice and delegate performance

Every ticket ran as one ticket-sized write charge on opus/medium, except ticket 5 on opus/high. Tickets 1 and 3 ran in parallel siblings, and tickets 2, 4, 5, and 6 ran serially on the integration source.

Ticket 3 landed first-pass with a biting self-probe and a biting coordinator probe. Ticket 1 needed two continuation rounds for out-of-fence readers of the landing stdout, and one trim for an over-budget file that grew by comment lines. Ticket 2 landed with one out-of-fence system reader. Ticket 4 landed with one seam-literal fold and, after the review, one skip-routing repair. Ticket 5 landed with one second canary fixture. Ticket 6 needed one file split for an over-budget parser.

No delegate returned a vacuous probe, and every coordinator probe bit at a distinct site and kind.

## Coordinator catches

- The spec's reader sweep missed four byte-for-byte stdout readers of the landed record, two of them in the system suite. It also missed a fifth build-input fixture whose stub build script never created its output directory.
- The spec's fence omitted the seam registry that spelled the literals the new reader exported. A second canary fixture anchored on the sentence ticket 5 replaced.
- The ticket 4 charge omitted the skip-ownership check the delegation discipline names for a test that can skip. The fold gate caught two bare `t.Skip` calls one round late.
- The Coverage axis found that a state file whose fence never closes bricked the handoff document, which the heading refusal alone did not catch.
- Another session's `bench worktree clean --landed --apply` swept the first zero-commit integration worktree thirteen seconds after creation. The landing's own narrowed cleanup later swept a zero-commit close worktree at the published commit.
- `bench doctor --fix` wrote the broker manifest for the old executable, so the sanctioned rebuild still had to run by hand before the new grammar was live.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1.md | 3 | spec-row, spec-row, delegate-error |
| 2.md | 1 | spec-row |
| 3.md | none | none |
| 4.md | 2 | spec-row, other |
| 5.md | 1 | spec-row |
| 6.md | 1 | tree-drift |

## Agent-experience improvements

### Bench CLI

- Cite the landing's census entry in `capture/learnings.md`, titled `ft311-lc census: 5 raw calls, siblings ft311-lc-t1 1 and ft311-lc-t3 1`, which folds the delegates' CLI notes.
  Feeds: new
- Add a `--path` filter to `bench structure`, so a delegate confirms its fenced files' line counts without `wc` and without a 105-row wall.
  Feeds: new
- Make `bench probe` name a non-compiling mutant as its own cause, so a delegate does not read a compiler error under the executable-selection line.
  Feeds: new
- Make `bench test --check prose` print the same packages table every other named check prints, so a caller can tell a pass from a no-op.
  Feeds: new

### Skills

- Add to the delegation discipline that the coordinator ticks its In-the-charge list against the ticket's `Writes:` before dispatch.
  Feeds: none

### Process

- Decide whether the landing's narrowed cleanup excludes a branch with no commit past the destination base, so a fresh sibling at the published commit survives.
  Feeds: FT311
- Do not create a sibling worktree between a landing's publication and its effects, or during another session's landing tail.
  Feeds: none
- Make the retrospective scaffold prefer the newest landing span that carries a subject, so a resume does not hide the landing's stages.
  Feeds: FT311
