# Session handoff emission (FT122)

Status: staged

Compiled from `decisions/session-handoff-emission.md`, which was closed in the
same session that wrote this spec under the reviewer-closed path in
`/bench-write-spec`'s entry contract. The map is flagged for reviewer veto, and
its **[veto]** marks are recorded alternatives rather than open questions.

Amended 2026-07-25 after a falsification pass. Eleven findings were confirmed
against the tree and folded in; the map absorbed the four that belonged to it,
including the gate-verdict field this spec had originally introduced with no
decision behind it.

## Problem

A phase close has to hand the next session a set of facts it cannot infer:
which repository and branch, what HEAD is, whether the tree is clean, what is
unpushed, which spec is staged and at what status, what the gate last said, and
the exact command to run next. Every one of those is derivable from the tree,
and today a session derives them by hand — reading `git status`, `git log`,
`specs/`, and the status board, then transcribing the result into prose. In the
week of 2026-07-19 the reviewer asked for a fresh-session prompt thirteen times,
and each request paid that reconstruction cost again.

The facts are also transcribed twice in two forms, because the working agreement
lets a phase close either emit a copy-paste continuation prompt or update
`session-handoff.md`. Two hand-derivations of one set of facts is the
duplicated-knowledge defect the code standard names, and it drifts the ordinary
way: a session that rewrites the body and misremembers the commit leaves a
confident, wrong pin block behind.

`bench status` already detects that the handoff has gone stale and says so. It
cannot fix it, because nothing emits the artifact.

## Solution

`bench handoff` collects the deterministic facts, renders the pin block, prints
it, and rewrites `session-handoff.md` — preserving the one section that carries
judgment and regenerating the rest. The session writes only what it actually
knows: the `## State` prose about what is true and what it means.

The command names the next phase in the form the *receiving* harness can invoke,
selected by `--harness`, so a handoff written for Codex does not send that
session to a Claude Code key it does not have. It names the repository by
identity — name and `origin` remote — with the filesystem path as a derived hint
rather than a constant, so the artifact survives being read on another machine.
Every fact it states about staleness — the gate's, and its own — is computed
from the tree rather than remembered, because a confidently wrong pin block is
the failure this capability exists to remove.

## User stories

1. As a session closing a phase, I want `bench handoff` to print the pin block
   and rewrite `session-handoff.md` in one invocation, so that both phase-close
   forms come from one derivation instead of two hand-transcriptions.
   Line: `gpt-5.6-luna` / medium. The command wires together facts other packages
   already expose, which is plumbing at a known seam, but it owns a file write
   that must not lose data and so earns the higher effort within the cheap tier.

2. As a cold reader on any machine, I want the header to name the repository by
   its name and `origin` remote, so that I can identify the repo without trusting
   a filesystem path that may not exist where I am reading.
   Line: `gpt-5.6-luna` / low. Reading two values out of `internal/git` and
   formatting them is mechanical once the fact-collection seam exists.

3. As a reviewer who moves this repo between machines, I want the path rendered
   from the actual git root and abbreviated to `~` when it sits under `$HOME`, so
   that no directory layout is baked into the kit.
   Line: `gpt-5.6-luna` / low. A single path transformation with an obvious
   boundary condition, fully covered by the gate.

4. As a reviewer whose checkout lives outside `$HOME`, I want the path rendered
   absolute rather than forced into a `~` form, so that the hint stays true
   instead of becoming a path that does not resolve.
   Line: `gpt-5.6-luna` / low. This is the single edge of story 3's transformation
   and shares its seam.

5. As a cold session, I want the block to carry branch, HEAD, whether the tree is
   clean, and how many commits are unpushed, so that I know what state I am
   resuming into before I touch anything.
   Line: `gpt-5.6-luna` / low. Four values from `internal/git`, each already
   exposed for the status board.

6. As a cold session, I want the staged spec named with its `Status:` line, or an
   explicit statement that none is staged, so that I know whether my next move is
   to write a spec or to build one.
   Line: `gpt-5.6-luna` / low. `internal/spec.Facts` already returns exactly this.

7. As a cold session, I want the gate field to carry the verdict, the tree it was
   computed on, and whether that tree is still current — or to say plainly that no
   gate has run — so that I never read a cached green as a statement about the
   tree I actually inherited.
   Line: `gpt-5.6-luna` / medium. `internal/status.GateVerdict` already returns
   `Present`, `Stale`, `CachedTree`, and `WorkTree`, so the work is rendering
   rather than derivation, but a bare verdict here would reproduce the exact
   defect this capability exists to remove, which lifts it above the low-effort
   rows.

8. As a cold session, I want the next command taken from the same signals
   `bench status` ranks — the highest-severity signal whose action is actually an
   invocable command — so that the field names something I can run rather than a
   description of a situation.
   Line: `gpt-5.6-terra` / medium. This deviates from the profile's cheap
   CLI-plumbing routing because reusing the board's precedence rather than
   inventing a parallel rule is the decision that keeps this single-sourced, and
   getting it wrong produces two competing recommendations.

   The board's `Action` is a prose hint — `fix before commit`,
   `split (craft-seams)`, `commit on green / /bench-final-check / push` — and only
   a minority of rows are invocable. This field promises an invocation, so it
   selects rather than converts: it walks the board in its own severity order and
   takes the first action that is one, leaving the board's ranking untouched. A
   board whose signals are all prose states that plainly and points at `--next`;
   it never renders a hint as though it were a command. Selection reads the
   canonical form, so it is unaffected by story 10's harness translation.

   Opening with an invocation is necessary but not sufficient. A board row may
   join the steps of a sequence into one action — the git row reads
   `/bench-final-check / push` once the tree is clean — and that string opens a
   phase invocation while naming two commands. Splitting it and taking an arm
   would be this command deciding what the board meant by a sequence, so a
   compound action does not qualify and the walk continues past it.

9. As a session whose real next step is judgment the board cannot see, I want
   `--next` to replace the derived line verbatim, so that a review pass or a
   debug loop can be named instead of a mechanical guess.
   Line: `gpt-5.6-luna` / low. One override branch at the point of rendering.

10. As a session handing off to Codex, I want `--harness codex` to render phase
    invocations as `$bench-*`, so that the receiving session can actually invoke
    what it is told to run.
    Line: `gpt-5.6-luna` / medium. The translation itself is trivial; the effort
    is in routing every invocation through one table so a second form cannot be
    hardcoded elsewhere.

11. As a session handing off within Claude Code, I want the default and
    `--harness claude` to render `/bench-*`, so that the common case needs no
    flag and no thought.
    Line: `gpt-5.6-luna` / low. The default arm of story 10's table.

12. As a session that mistypes the flag's value, I want an unknown `--harness`
    value rejected as a usage error rather than silently falling back, so that a
    wrong harness form cannot reach a handoff unnoticed.
    Line: `gpt-5.6-luna` / medium. `internal/handoff` must own this: the shared
    arg-grammar helper validates flag spelling, repetition, and arity, but accepts
    any value for a declared value-taking flag, so the validation and its
    `toon.Usage` rendering are this package's work rather than a free
    composition.

13. As the author of a handoff's judgment prose, I want `## State` preserved
    byte-for-byte, so that running the command never costs me writing no CLI
    could reproduce.
    Line: `gpt-5.6-terra` / medium. This deviates from the cheap CLI-plumbing
    routing because a defect here destroys reviewer work rather than producing a
    wrong string, and the splitter it rests on is the design's load-bearing part.

14. As a maintainer of the kit, I want the existing `## Shape` text moved into
    the binary and regenerated from there, so that the handoff's own conventions
    cannot drift per repo and a linked repo picks up revisions on `bench upgrade`.
    Line: `gpt-5.6-luna` / medium. This story is transcription — today's text
    moves verbatim, and story 24's conformance check is what makes the move safe
    — so it stays at the CLI-plumbing routing rather than the doc-authoring one.

15. As the first session in a repo that has never had a handoff, I want the
    scaffolded skeleton's prose to tell me what this document is for and how to
    keep it, so that the conventions arrive with the artifact instead of having to
    be learned from another repo.
    Line: `gpt-5.6-sol` / high. This is new guidance prose that ships to every
    linked repo on `bench upgrade` and is read cold by every session that meets
    an empty handoff, which is the leverage override in `craft-line`. The
    reviewer approved this top-tier bump explicitly; it is bounded to the
    skeleton's prose and does not extend to the surrounding write logic.

16. As a repo that has never had a handoff, I want the command to scaffold the
    file with an empty `## State` and exit zero, so that a new repo does not pay
    a manual first write.
    Line: `gpt-5.6-luna` / low. Writing a fixed skeleton when a file is absent.

17. As a reviewer with a hand-edited handoff, I want the command to refuse and
    change nothing when it cannot find `## State`, so that prose I may not have
    committed is never overwritten by a guess.
    Line: `gpt-5.6-terra` / medium. Deviates from the cheap routing for the same
    reason as story 13: the fail-closed posture is its counterpart and carries the
    same destroy-nothing consequence.

18. As a session that runs the command twice, I want byte-identical output on
    both sinks on an unchanged tree, so that an exploratory invocation is
    recoverable by inspection rather than by `git checkout`.
    Line: `gpt-5.6-luna` / medium. Idempotence is what makes the bare-command
    write safe, so it is asserted rather than assumed.

19. As a reviewer whose handoff contains an ambiguous `## State` — two such
    headings, or one inside a fenced code block — I want the command to fail
    closed rather than pick one, so that ambiguity never resolves into data loss.
    Line: `gpt-5.6-terra` / medium. Deviates from the cheap routing because
    fence-aware splitting is the subtlest logic in the command and the place a
    plausible implementation is silently wrong.

20. As a session in a degenerate git state — detached HEAD, no commits yet, or no
    `origin` remote — I want each affected field to render an explicit unknown,
    so that a missing fact reads as missing rather than as an empty value that
    looks deliberate.
    Line: `gpt-5.6-luna` / medium. Three independent degenerate states, each
    needing its own explicit rendering.

21. As a session running against an unwritable path, I want a structured error
    naming it, so that the failure is actionable instead of a bare permission
    trace.
    Line: `gpt-5.6-luna` / low. The existing structured-error helper covers the
    shape.

22. As a session that passes a repeated or unknown flag name, I want the shared
    argument grammar to reject it, so that this command behaves like every other
    Bench subcommand.
    Line: `gpt-5.6-luna` / low. Composing the FT87 helper; no new parsing.

23. As a session outside a repository, I want the existing not-in-a-repo error
    and its exit code, so that the failure matches what every other porcelain
    command does.
    Line: `gpt-5.6-luna` / low. One guard reusing `toon.NotInRepo()`.

24. As a maintainer, I want a conformance check proving the `## Shape` text
    exists in exactly one place, so that story 14's single-sourcing is enforced
    rather than merely intended.
    Line: `gpt-5.6-terra` / medium. Gate logic; the profile routes gate and
    conformance work to mid effort because a wrong oracle is the worst class of
    defect in this kit, and the model moves up with it.

25. As a maintainer, I want a conformance check proving the harness
    phase-prefix table is the only place either invocation form is produced, so
    that a hardcoded `/bench-` or `$bench-` cannot survive somewhere the runtime
    rows do not look.
    Line: `gpt-5.6-terra` / medium. Gate logic, routed as story 24, and the
    counterpart that makes the map's single-sourced-table contract enforceable.

26. As a session at the shell, I want `bench handoff` routed like every other
    porcelain command, so that it is discoverable and behaves consistently.
    Line: `gpt-5.6-luna` / low. One line in `bin/bench.sh`.

## Implementation decisions

A new `internal/handoff` package owns the capability: fact collection, section
splitting, rendering, and the write. It is deep — callers ask for a rendered
block and never see harness translation, `~` abbreviation, or fence-aware
section parsing. Its `Command` func stays a thin argument-grammar pass-through,
matching every other porcelain command, and `bin/bench.sh` gains a single
`handoff) route_porcelain "$@" ;;` case.

It composes rather than re-derives. `internal/git` supplies identity, branch,
HEAD, dirty state, and unpushed count; `internal/spec.Facts` supplies the staged
spec path and Status; `internal/status` supplies the ranked signals the next
action is selected from and, via `GateVerdict`, the gate field's verdict, cached
tree, and staleness. Selection over those signals is a syntactic test for an
invocation — the prefixes the kit's own commands take — not a second opinion
about severity, so the board keeps sole ownership of what outranks what.
`internal/status` keeps sole ownership of staleness detection — this command
does not compute or reset either the handoff's age or the gate's, and does not
stamp a date into the file, because staleness is read from git history precisely
to avoid a self-reported date that a forgetful rewrite leaves confidently wrong.
For the same reason the gate field is never rendered as a bare verdict: it
carries the tree it describes and its staleness, or states that no gate has run.

Argument handling splits by what the shared helper can actually enforce. Flag
spelling, repetition, and arity go to the FT87 arg-grammar helper. The
`--harness` *value* is validated by `internal/handoff` and rejected through
`toon.Usage`, because the helper accepts any value for a declared value-taking
flag and cannot own this.

Harness translation is one single-sourced table mapping a harness name to its
phase-invocation prefix, consulted at render time, with a conformance check
proving no call site produces either form independently.

The repository is named by two facts with different lifetimes: identity (name
and `origin` remote) is the durable anchor, and the path is a hint derived at
emit time from the git root, abbreviated to `~` only when it sits under `$HOME`.
Where the two disagree for a cold reader, identity wins — the same precedence
the working agreement already sets between the handoff and the tree.

Section handling splits on headings with fence awareness, preserves `## State`
byte-for-byte, and regenerates the header, `## Next command`, and `## Shape`.
Regeneration replaces rather than appends: the superseded text must be gone, and
exactly one of each generated heading may remain. The command creates freely and
destroys never — a missing file is scaffolded, and a file whose `## State`
cannot be located unambiguously is left untouched behind a non-zero exit.

## Testing decisions

A good test here drives the built command against a fixture repository and reads
the bytes it printed and wrote. The behaviors that matter are all observable at
that boundary — what the block says, what survived in the file, what the exit
code was — and none of them require reaching into the package's internals.

Two pure transformations also get unit tests, for their edge sets: the section
splitter and `~` rendering. Those are **additions to** the end-to-end rows, never
replacements. A correct helper the command never calls is precisely the failure a
unit-only seam would miss, so each of them keeps a runtime row that reads the
value out of the command's own output.

Derived values are asserted against at least two fixtures that differ in the
value being asserted, wherever a single fixture would let a hardcoded constant
pass. A one-fixture assertion of a derived value grades nothing.

Prior art: `internal/contract/runtime` already runs porcelain commands against
fixture repos and asserts output and exit codes, and `internal/conformance`
already hosts single-sourcing checks of the kind stories 24 and 25 need. The
runtime contract tests run against the built binary, so `dist/bench` is rebuilt
before they run.

Gate command: the project gate, `bench gate`.

### Seam diagram

    trigger: a session closing a phase, at the shell
        │
        ▼
    argv (--harness, --next)  ──▶  [ internal/handoff ]  ──▶  stdout: pin block
    the git tree              ──▶  [   facts→render→   ]  ──▶  session-handoff.md
    session-handoff.md        ──▶  [       write       ]  ──▶  exit code
                                     ◀ tests attach here: run the built binary in
                                       a fixture repo; assert stdout, the file's
                                       bytes, and the exit code

    trigger: internal/handoff, while parsing an existing handoff
        │
        ▼
    file bytes  ──▶  [ section splitter ]  ──▶  (preserved State, rest) | refusal
                        ◀ tests attach here: unit tests over handcrafted
                          bodies — absent, doubled, fenced, well-formed —
                          paired with a runtime row proving the command
                          calls this splitter and not another

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | bare invocation prints the block, writes the file, exits 0 | runtime contract | `go test ./internal/contract/runtime -run TestHandoffWritesAndPrints` — fails, no such command | Asserting both sinks and the exit code in one run means an implementation serving only one cannot pass. |
| 2 | header names repo name and `origin` remote | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNamesIdentity` | Two fixtures with different remotes; a hardcoded URL passes one and fails the other. |
| 3 | `~` abbreviation under `$HOME` | unit (`internal/handoff`) | `go test ./internal/handoff -run TestRenderPathAbbreviates` | Enumerates the transformation's edges directly, including exactly at `$HOME`. |
| 4 | path outside `$HOME` renders absolute | unit (`internal/handoff`) | `go test ./internal/handoff -run TestRenderPathOutsideHome` | Forcing `~` on an outside path yields a string that does not resolve; asserting absolute catches it. |
| 3, 4 | the *emitted* path equals the fixture's git root, `~`-abbreviated | runtime contract | `go test ./internal/contract/runtime -run TestHandoffEmitsDerivedPath` | Fixture rooted outside any `workspace` directory. Catches a command rendering `os.Getwd()` or a constant while a correct helper sits unit-tested and uncalled. |
| 5 | branch, HEAD, clean/dirty, unpushed count | runtime contract | `go test ./internal/contract/runtime -run TestHandoffCarriesTreeFacts` | Two fixtures — dirty with an unpushed commit, and clean with none — so hardcoded defaults fail one of them. |
| 6 | staged spec named with Status, or explicitly none | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNamesStagedSpec` | Three fixtures: two different `Status:` values plus an empty `specs/`. A constant `Status: staged` fails the second. |
| 7 | gate field carries verdict, cached tree, and staleness; or says none has run | runtime contract | `go test ./internal/contract/runtime -run TestHandoffGateFieldIsStaleAware` | Three fixtures — green-and-current, green-but-stale, absent. A hardcoded `gate: green` fails the stale and absent cases, which is the confident-wrong-fact defect. |
| 8 | selection returns the first invocable action, skips prose and compound ones, and reports when a board has none | unit (`internal/handoff`) | `go test ./internal/handoff -run TestFirstInvocable` | Enumerates the selection's edges directly — empty board, all-prose board, prose ahead of an invocable row, and the git row's `/bench-final-check / push` — including the all-prose case a live fixture cannot reliably produce. A bare prefix test passes every other row and fails the compound one. |
| 8 | next command is the board's highest-severity **invocable** action, or an explicit statement that none is | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNextMatchesStatus` | Three fixtures: two whose first invocable board action differs, each compared to `bench status --all` on that same fixture, plus one whose only signals carry prose actions. A constant fails the first two; taking `signals[0].Action` unconditionally fails the third by rendering a hint as a command. |
| 9 | `--next` replaces the derived line | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNextOverride` | Asserts the override present **and** the derived line absent, so appending alongside it fails. |
| 10 | `--harness codex` renders `$bench-*` | runtime contract | `go test ./internal/contract/runtime -run TestHandoffHarnessCodex` | Asserting `$bench-` present *and* `/bench-` absent catches an implementation that emits both forms. |
| 11 | default and `claude` render `/bench-*` | runtime contract | `go test ./internal/contract/runtime -run TestHandoffHarnessDefault` | The paired absence assertion catches a default that leaks the Codex form. |
| 12 | unknown `--harness` value is a usage error, exit 2, nothing written | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRejectsUnknownHarness` | A silent fallback would exit 0 and write a handoff; asserting exit 2 and an unchanged tree forbids it. |
| 13 | `## State` preserved byte-for-byte | runtime contract | `go test ./internal/contract/runtime -run TestHandoffPreservesState` | Distinctive prose compared byte-for-byte fails any regeneration or reflow of that section. |
| 14 | `## Shape` replaced, not appended | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRegeneratesShape` | Fixture carries a deliberately stale Shape; asserts the current text present, the stale text **absent**, and exactly one `## Shape` heading. Kills preserve-and-append. |
| 15 | the scaffolded skeleton's prose states the artifact's purpose and conventions | runtime contract | `go test ./internal/contract/runtime -run TestHandoffSkeletonCarriesConventions` | A skeleton that scaffolds structure without the guidance prose leaves the first session nothing to learn the conventions from. |
| 16 | missing file scaffolded with empty State, exit 0 | runtime contract | `go test ./internal/contract/runtime -run TestHandoffScaffoldsMissing` | A fixture with no handoff proves creation happens rather than an error. |
| 17 | no `## State` heading → non-zero, file unchanged | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRefusesUnparseable` | Asserting both the exit code and the unchanged bytes forbids the write-anyway implementation. |
| 18 | two runs byte-identical on stdout **and** the file | runtime contract | `go test ./internal/contract/runtime -run TestHandoffIdempotent` | Fixture's `## State` is non-empty, so the second run also proves the command re-parses its own output. Any emitted timestamp or reordering breaks equality. |
| 19 | doubled or fenced `## State` fails closed | unit (`internal/handoff`) | `go test ./internal/handoff -run TestSplitAmbiguousState` | A naive first-match splitter passes the well-formed case and fails here, which is the silent-data-loss path. |
| 19 | the *command* refuses a fenced or doubled `## State` | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRefusesAmbiguousState` | Proves the command calls the tested splitter rather than an inline fence-blind one. |
| 20 | detached HEAD, no commits, no `origin` render explicit unknowns | runtime contract | `go test ./internal/contract/runtime -run TestHandoffDegenerateGit` | Three fixtures; an implementation emitting empty fields fails because absence must read as stated, not blank. |
| 21 | unwritable target yields a structured error naming the path | runtime contract | `go test ./internal/contract/runtime -run TestHandoffUnwritableTarget` | Parent directory at mode 0555, skipped under the existing skip-ownership rule when running as root. A bare permission trace fails the structured-shape assertion. |
| 22 | repeated and unknown flag names rejected by the shared grammar | runtime contract | `go test ./internal/contract/runtime -run TestHandoffArgGrammar` | Hand-rolled parsing accepts what the shared helper rejects, so the matrix catches a bypass. |
| 23 | outside a repository, `toon.NotInRepo()` and its exit code | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNotInRepo` | Pins the error line and the code, not merely that no write happened. |
| 24 | `## Shape` text exists in exactly one place | conformance | `go test ./internal/conformance -run TestHandoffShapeSingleSourced` | With the text in the binary, a surviving copy in a tracked file turns this red. |
| 25 | the harness prefix table is the only producer of either form | conformance | `go test ./internal/conformance -run TestHarnessPrefixSingleSourced` | A trailing `strings.ReplaceAll` or an inline conditional passes rows 10 and 11 but fails here. |
| 26 | `bench handoff` routes as porcelain | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRouted` | Invoking through the wrapper rather than the binary fails if the case line is missing. |
| 20 | control bytes in a branch or spec name are refused | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRefusesControlBytes` | Fixture branch name carries a control byte; asserts the existing `toon` refusal fires rather than the byte reaching the rendered block. |

### Edge inventory

Walked per behavior; each landed as a row above or as a **Won't handle** line.

- **Absent vs empty** — rows on stories 7, 16, 17, 20: no gate, no file, no
  `## State`, and empty git facts are four distinct outcomes, none collapsing to
  another.
- **Stale vs current** — story 7: a cached verdict whose tree has moved.
- **Ambiguity / duplication** — story 19, both rows: doubled and fenced headings.
- **Replacement vs accumulation** — stories 9 and 14: appending alongside the
  thing being replaced.
- **Boundary of a transformation** — stories 3 and 4: exactly at `$HOME`, and
  outside it.
- **Degenerate environment** — stories 20, 21, 23: degenerate git states,
  unwritable target, no repository.
- **Argument grammar** — stories 12 and 22: unknown value, repeated and unknown
  flag names.
- **Hostile bytes** — story 20's control-byte row.
- **Idempotence / repeated application** — story 18.
- **Cross-source agreement** — story 8: the handoff versus the status board.

**Won't handle** — each is recorded in the map's Out of scope, not invented here:

- Concurrent invocations racing on the same file — the working agreement routes a
  second writer to a `bench worktree`, and the tripwire that catches concurrent
  writers is FT91's.
- Symlinked repository roots resolving differently than the git root reports —
  the path is a hint with identity as the anchor, so a divergence is cosmetic.
- Interpreting the contents of `## State` — preserved verbatim, so prose inside
  it that resembles a generated section is never parsed.
- Non-UTF-8 bytes in the handoff *body* — preserved verbatim through the
  splitter. Control bytes in *rendered field values* are a covered row above, not
  a cut.

## Out of scope

- **A machine-readable (TOON) form of the pin block** — a separate capability
  serving other tooling rather than a reader, and nothing consumes it yet.
  2 edits, 1 gate run.
- **`bench status` naming `bench handoff` in its staleness action** — a
  different command's output, and it should not cite this one until the command
  has been adopted in practice. 1 edit, 1 gate run.
- **Harnesses beyond `claude` and `codex`** — each is a new row in the
  translation table plus its fixtures, added when a third harness is actually in
  use. 2 edits, 1 gate run per harness.
- **Committing the rewritten handoff** — the reviewer owns the merge, and the
  existing commit path already covers it; the staleness row resetting on commit
  is correct behavior, not a gap.
