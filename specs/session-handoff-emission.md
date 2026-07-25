# Session handoff emission (FT122)

Status: staged

Compiled from `decisions/session-handoff-emission.md`, which was closed in the
same session that wrote this spec under the reviewer-closed path in
`/bench-write-spec`'s entry contract. The map is flagged for reviewer veto, and
its three **[veto]** marks are recorded alternatives rather than open questions.

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

7. As a cold session, I want the last gate verdict in the block, so that I know
   whether the tree I inherited was green when it was handed over.
   Line: `gpt-5.6-luna` / low. Reads the same recorded verdict the status board
   reads.

8. As a cold session, I want the next command derived from the same signals
   `bench status` uses to compute its action, so that the handoff and the board
   never disagree about what to do next.
   Line: `gpt-5.6-terra` / medium. Reusing the board's precedence rather than
   inventing a parallel rule is the decision that keeps this single-sourced, and
   getting it wrong would produce two competing recommendations.

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

12. As a session that mistypes the flag, I want an unknown `--harness` value to
    be a usage error rather than a silent fallback, so that a wrong harness form
    cannot reach a handoff unnoticed.
    Line: `gpt-5.6-luna` / low. The shared arg-grammar helper owns the exit code
    and message shape.

13. As the author of a handoff's judgment prose, I want `## State` preserved
    byte-for-byte, so that running the command never costs me writing no CLI
    could reproduce.
    Line: `gpt-5.6-terra` / medium. This is the one story where a defect destroys
    reviewer work rather than producing a wrong string, and the splitter it rests
    on is the design's load-bearing part.

14. As a maintainer of the kit, I want `## Shape` regenerated from a single
    in-binary source, so that the handoff's own conventions cannot drift per repo
    and a linked repo picks up revisions on `bench upgrade`.
    Line: `gpt-5.6-terra` / medium. Moving prose into the binary is only safe if
    the old copy provably cannot survive, which is a single-sourcing judgment.

15. As a repo that has never had a handoff, I want the command to scaffold the
    skeleton with an empty `## State`, so that a new repo does not pay a manual
    first write.
    Line: `gpt-5.6-luna` / low. Writing a fixed skeleton when a file is absent.

16. As a reviewer with a hand-edited handoff, I want the command to refuse and
    change nothing when it cannot find `## State`, so that prose I may not have
    committed is never overwritten by a guess.
    Line: `gpt-5.6-terra` / medium. The fail-closed posture is the counterpart to
    story 13 and carries the same destroy-nothing consequence.

17. As a session that runs the command twice, I want byte-identical output on an
    unchanged tree, so that an exploratory invocation is recoverable by
    inspection rather than by `git checkout`.
    Line: `gpt-5.6-luna` / medium. Idempotence is what makes the bare-command
    write safe, so it is asserted rather than assumed.

18. As a reviewer whose handoff contains an ambiguous `## State` — two such
    headings, or one inside a fenced code block — I want the splitter to fail
    closed rather than pick one, so that ambiguity never resolves into data loss.
    Line: `gpt-5.6-terra` / medium. Fence-aware splitting is the subtlest logic
    in the command and the place a plausible implementation is silently wrong.

19. As a session in a degenerate git state — detached HEAD, no commits yet, or no
    `origin` remote — I want each affected field to render an explicit unknown,
    so that a missing fact reads as missing rather than as an empty value that
    looks deliberate.
    Line: `gpt-5.6-luna` / medium. Three independent degenerate states, each
    needing its own explicit rendering.

20. As a session running against a read-only file system, I want a structured
    error naming the path, so that the failure is actionable instead of a bare
    permission trace.
    Line: `gpt-5.6-luna` / low. The existing structured-error helper covers the
    shape.

21. As a session that passes a repeated or unknown flag, I want the shared
    argument grammar to reject it, so that this command behaves like every other
    Bench subcommand.
    Line: `gpt-5.6-luna` / low. Composing the FT87 helper; no new parsing.

22. As a session outside a repository, I want the existing not-in-a-repo error,
    so that the failure matches what every other porcelain command does.
    Line: `gpt-5.6-luna` / low. One guard reusing `toon.NotInRepo()`.

23. As a maintainer, I want a conformance check proving the `## Shape` text
    exists in exactly one place, so that story 14's single-sourcing is enforced
    rather than merely intended.
    Line: `gpt-5.6-terra` / medium. Gate logic, which the profile routes to mid
    effort because a wrong oracle is the worst class of defect in this kit.

24. As a session at the shell, I want `bench handoff` routed like every other
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
spec path and Status; `internal/status` supplies the derived next action.
`internal/status` keeps sole ownership of handoff staleness — this command does
not compute or reset it, and does not stamp a date into the file, because
staleness is read from git history precisely to avoid a self-reported date that
a forgetful rewrite leaves confidently wrong.

Harness translation is one single-sourced table mapping a harness name to its
phase-invocation prefix, consulted at render time. No call site hardcodes either
form.

The repository is named by two facts with different lifetimes: identity (name
and `origin` remote) is the durable anchor, and the path is a hint derived at
emit time from the git root, abbreviated to `~` only when it sits under `$HOME`.
Where the two disagree for a cold reader, identity wins — the same precedence
the working agreement already sets between the handoff and the tree.

Section handling splits on headings with fence awareness, preserves `## State`
byte-for-byte, and regenerates the header, `## Next command`, and `## Shape`.
The command creates freely and destroys never: a missing file is scaffolded, and
a file whose `## State` cannot be located unambiguously is left untouched behind
a non-zero exit.

## Testing decisions

A good test here drives the built command against a fixture repository and reads
the bytes it printed and wrote. The behaviors that matter are all observable at
that boundary — what the block says, what survived in the file, what the exit
code was — and none of them require reaching into the package's internals. Two
exceptions test below that seam because they are pure transformations with edge
sets worth enumerating directly: section splitting and `~` rendering.

Prior art: `internal/contract/runtime` already runs porcelain commands against
fixture repos and asserts output and exit codes, and `internal/conformance`
already hosts single-sourcing checks of the kind story 23 needs. The runtime
contract tests run against the built binary, so `dist/bench` is rebuilt before
they run.

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
                          bodies — absent, doubled, fenced, well-formed

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | bare invocation prints the block and writes the file | runtime contract | `go test ./internal/contract/runtime -run TestHandoffWritesAndPrints` — fails, no such command | Asserting both sinks in one run means an implementation that serves only one cannot pass. |
| 2 | header names repo name and `origin` remote | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNamesIdentity` | A fixture with a known remote proves identity is emitted, not just a path. |
| 3 | path derived from git root, `~` when under `$HOME` | unit (`internal/handoff`) | `go test ./internal/handoff -run TestRenderPathAbbreviates` | A fixture rooted at a non-`workspace` path fails any implementation carrying a layout constant. |
| 4 | path outside `$HOME` renders absolute | unit (`internal/handoff`) | `go test ./internal/handoff -run TestRenderPathOutsideHome` | Forcing `~` on an outside path yields a string that does not resolve; asserting absolute catches it. |
| 5 | branch, HEAD, clean/dirty, unpushed count present | runtime contract | `go test ./internal/contract/runtime -run TestHandoffCarriesTreeFacts` | A fixture with a dirty file and an unpushed commit distinguishes real derivation from hardcoded defaults. |
| 6 | staged spec named with Status, or explicitly none | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNamesStagedSpec` | Two fixtures — one staged, one empty — mean a stub answering one way always fails the other. |
| 7 | last gate verdict present | runtime contract | `go test ./internal/contract/runtime -run TestHandoffCarriesGateVerdict` | A fixture with a recorded verdict proves the field is read rather than omitted. |
| 8 | next command matches the board's derivation | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNextMatchesStatus` | Asserting the handoff's line equals `bench status`'s action on the same fixture makes a second, drifting rule fail. |
| 9 | `--next` replaces the derived line verbatim | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNextOverride` | A distinctive override string absent from output proves the flag was ignored. |
| 10 | `--harness codex` renders `$bench-*` | runtime contract | `go test ./internal/contract/runtime -run TestHandoffHarnessCodex` | Asserting `$bench-` present *and* `/bench-` absent catches an implementation that emits both forms. |
| 11 | default and `claude` render `/bench-*` | runtime contract | `go test ./internal/contract/runtime -run TestHandoffHarnessDefault` | The paired absence assertion catches a default that leaks the Codex form. |
| 12 | unknown `--harness` is a usage error, exit 2 | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRejectsUnknownHarness` | A silent fallback would exit 0 and write a handoff; asserting exit 2 forbids it. |
| 13 | `## State` preserved byte-for-byte | runtime contract | `go test ./internal/contract/runtime -run TestHandoffPreservesState` | Distinctive prose compared byte-for-byte fails any regeneration or reflow of that section. |
| 14 | `## Shape` regenerated from the in-binary source | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRegeneratesShape` | A fixture whose Shape text is deliberately stale must come back current, which preservation cannot do. |
| 15 | missing file is scaffolded with empty State | runtime contract | `go test ./internal/contract/runtime -run TestHandoffScaffoldsMissing` | A fixture with no handoff proves creation happens rather than an error. |
| 16 | no `## State` heading → non-zero, file unchanged | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRefusesUnparseable` | Asserting both the exit code and the unchanged bytes forbids the write-anyway implementation. |
| 17 | two runs produce byte-identical output | runtime contract | `go test ./internal/contract/runtime -run TestHandoffIdempotent` | Any emitted timestamp, counter, or reordering breaks equality on the second run. |
| 18 | doubled or fenced `## State` fails closed | unit (`internal/handoff`) | `go test ./internal/handoff -run TestSplitAmbiguousState` | A naive first-match splitter passes the well-formed case and fails here, which is exactly the silent-data-loss path. |
| 19 | detached HEAD, no commits, no `origin` render explicit unknowns | runtime contract | `go test ./internal/contract/runtime -run TestHandoffDegenerateGit` | Three fixtures; an implementation emitting empty fields fails because absence must read as stated, not blank. |
| 20 | read-only filesystem yields a structured error naming the path | runtime contract | `go test ./internal/contract/runtime -run TestHandoffReadOnlyFS` | A bare permission trace fails the assertion on the structured shape. |
| 21 | repeated and unknown flags rejected by the shared grammar | runtime contract | `go test ./internal/contract/runtime -run TestHandoffArgGrammar` | Hand-rolled parsing accepts what the shared helper rejects, so the matrix catches a bypass. |
| 22 | outside a repository, the standard error | runtime contract | `go test ./internal/contract/runtime -run TestHandoffNotInRepo` | Running in a non-repo temp dir fails any implementation that proceeds to write. |
| 23 | `## Shape` text exists in exactly one place | conformance | `go test ./internal/conformance -run TestHandoffShapeSingleSourced` | With the text in the binary, a surviving copy in a tracked file turns this red. |
| 24 | `bench handoff` routes as porcelain | runtime contract | `go test ./internal/contract/runtime -run TestHandoffRouted` | Invoking through the wrapper rather than the binary fails if the case line is missing. |

### Edge inventory

Walked per behavior; each landed as a row above or as a **Won't handle** line.

- **Absent vs empty** — rows 15, 16, 19: no file, no `## State`, and empty
  git facts are three distinct outcomes, none collapsing to another.
- **Ambiguity / duplication** — row 18: doubled and fenced headings.
- **Boundary of a transformation** — rows 3, 4: exactly at `$HOME`, and outside it.
- **Degenerate environment** — rows 19, 20, 22: degenerate git states, read-only
  filesystem, no repository.
- **Argument grammar** — rows 12, 21: unknown value, repeated and unknown flags.
- **Idempotence / repeated application** — row 17.
- **Cross-source agreement** — row 8: the handoff versus the status board.

**Won't handle:**

- Concurrent invocations racing on the same file — the working agreement already
  routes a second writer to a `bench worktree`, and the tripwire that catches
  concurrent writers is FT91's, not this command's to duplicate.
- Symlinked repository roots resolving to a different path than the git root
  reports — the path is a hint with identity as the anchor, so a divergence is
  cosmetic rather than a wrong fact.
- A `## State` section containing what looks like the generated header — the
  section is preserved verbatim, so its contents are never interpreted.
- Non-UTF-8 bytes in the handoff body — preserved verbatim through the splitter,
  and the existing `toon` refusal covers control bytes in rendered fields.

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
