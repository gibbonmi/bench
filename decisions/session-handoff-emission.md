# Session handoff emission (FT122)

## Destination

`bench handoff` derives the cold-session pin block — the deterministic facts a
fresh session needs to resume — from the tree instead of a session
reconstructing them by hand, and lands them in both phase-close forms
`AGENTS.md` allows: the copy-paste continuation prompt and `session-handoff.md`.
Source: `IDEAS.md` 2026-07-24, from a transcript survey of 2026-07-19..25 that
counted 13 hand-written fresh-session prompts in six days.

Closed 2026-07-25 in the same session as the spec it feeds, under the
reviewer-closed path in `/bench-write-spec`'s entry contract. Every fork below
was put to the reviewer and answered; contestable calls are marked **[veto]**.

## #1: What does the command produce?

Blocked by: none
Type: Grill

### Question

`AGENTS.md` lets a phase close either emit a copy-paste continuation prompt or
update `session-handoff.md`. Does the command serve one form, or both?

### Answer

Both, on the bare invocation: `bench handoff` rewrites `session-handoff.md`
*and* prints the pin block to stdout. The reviewer's reasoning is that the block
is small enough that printing it costs nothing worth a flag, and the two sinks
carry identical facts derived once.

The write is a side effect on a bare, read-looking command, which is the one
cost of this shape. It is defused by idempotence rather than by a flag: with the
judgment section preserved verbatim (#4) and every generated section a pure
function of the tree, running the command twice on an unchanged tree produces
byte-identical output, so an exploratory invocation is recoverable by
inspection **[veto — the alternative was a print-only default with `--write`
opting into the file; the reviewer chose the simpler surface]**.

## #2: How does the command know which harness to format for?

Blocked by: #1
Type: Grill

### Question

The communication rules require a handoff to name the next phase in the form the
*receiving* harness can invoke. What decides that form?

### Answer

An explicit `--harness` flag, defaulting to `claude`. No environment sniffing:
a wrong guess sends the next session to a key its harness does not have, which
is the exact failure the translation rule exists to prevent, and detection would
have to be re-tuned per harness release.

The translation is narrow. The handoff body is harness-independent; only skill
and phase invocations change form — `claude` renders `/bench-write-spec`, and
`codex` renders `$bench-write-spec`. Harnesses beyond these two are out of scope
here and land as new cases in the same single-sourced table.

## #3: Who decides the `next` line?

Blocked by: #1
Type: Grill

### Question

The "Next command" section is the single most-copied line in a handoff. Does the
CLI derive it, or does the session supply it?

### Answer

The CLI derives it, and `--next` overrides. `bench status` already computes an
action per signal from the same tree facts, so deriving the line reuses that
source rather than adding a second one; the override exists because the genuinely
correct next step is sometimes judgment the board cannot see — a review pass, a
debug loop, a deliberate pause for the reviewer.

Derivation follows the board's existing precedence, not a new rule: whatever
`bench status` already leads with is what the handoff names.

## #4: What is preserved, and what is regenerated?

Blocked by: #1
Type: Grill

### Question

`session-handoff.md` mixes hand-written judgment with facts the tree owns. Which
sections may the command rewrite?

### Answer

`## State` is preserved verbatim — it is the session's judgment about what is
true and what it means, which no CLI can derive. Everything else is regenerated:
the header, `## Next command`, and `## Shape`.

`## Shape` is regenerated because it states the artifact's own contract (rewrite
in full, prune rather than accrete, keep these three sections). That contract is
kit-owned and identical in every linked repo, so hand-copying it is exactly the
duplicated-knowledge defect the code standard names; making the command its
source means a linked repo picks up revisions on `bench upgrade` **[veto — the
alternative was preserving Shape verbatim alongside State, trading drift for a
smaller blast radius on a hand-edited file]**.

## #5: What happens on a missing or unparseable handoff?

Blocked by: #4
Type: Grill

### Question

Preservation presumes a `## State` section exists to preserve. What does the
command do when there is no file, or no such section?

### Answer

Asymmetric, by what is recoverable. No file at all → scaffold the skeleton with
an empty `## State` for the session to fill, because creating a file that did
not exist destroys nothing and every new repo would otherwise pay a manual first
write. A file present but carrying no recognizable `## State` heading → exit
non-zero and change nothing, naming the missing section, because guessing which
prose is judgment risks overwriting hand-written work that may never have been
committed.

The general rule the two cases share: the command creates freely and destroys
never.

## #6: How is the repository named without pinning one machine's paths?

Blocked by: #1
Type: Grill

### Question

A handoff that hardcodes `/home/<user>/workspace/bench` breaks the moment the
repo is read on another machine. What does the header name instead?

### Answer

Two facts with different lifetimes. The **identity** — repository name and the
`origin` remote — is the durable anchor, true on every machine and the thing a
cold reader should trust. The **path** is a hint, derived at emit time from the
git root and rendered `~`-relative when it sits under `$HOME`, so it reflects
wherever the repo actually is on the machine that ran the command rather than
anything baked into the kit.

Nothing about a directory layout is a kit constant: a checkout at
`~/src/bench` renders that way with no configuration. The residual exposure is a
*committed* handoff carrying machine A's path being read on machine B, which the
identity anchor covers — where path and identity disagree, identity wins, the
same precedence `AGENTS.md` already sets between the handoff and the tree.

## Not yet specified

- Whether the pin block gains a machine-readable (TOON) form for other tooling
  to consume; today it is prose for a reader.
- Whether `bench status`'s handoff-staleness row should suggest `bench handoff`
  by name once the command exists (a one-line action change, but it couples the
  row to a command that may not be adopted).

## Out of scope

- The four other CLI candidates parked the same day (`bench worktree path` /
  `exec`, `bench test`, `bench spec show`, `bench outline --symbol`). Each is a
  separate capability with its own map.
- Changing what `session-handoff.md` is *for*, or its three-section shape. This
  command generates the artifact the working agreement already specifies.
- Harnesses beyond `claude` and `codex`.
- Committing the handoff. The staleness row resets when the rewrite lands in a
  commit, which stays the reviewer's call and the existing commit path's job.

## Handoff

1. **Module boundaries.** A new `internal/handoff` owns the pin block: fact
   collection, section splitting, rendering, and the write. It composes
   `internal/git` (identity, branch, HEAD, dirty, unpushed), `internal/spec`
   (`Facts` for staged spec path and Status), and `internal/status` (the derived
   next action). `internal/status` keeps sole ownership of staleness detection —
   this command does not re-derive it. `bin/bench.sh` gains one
   `handoff) route_porcelain "$@" ;;` line.
2. **Contracts.** Bare `bench handoff` writes the file and prints the block,
   exit 0. Missing file → scaffold, exit 0. File present without `## State` →
   exit non-zero, tree unchanged, message names the section. `--harness codex`
   renders phase invocations as `$bench-*`; default and `--harness claude`
   render `/bench-*`; an unknown value is a usage error, exit 2. `--next <cmd>`
   replaces the derived line verbatim. Not in a repo → the existing
   `toon.NotInRepo()` error. Re-running on an unchanged tree is byte-identical.
3. **Deep vs thin.** `internal/handoff` is deep: callers ask for a rendered
   block and never see section splitting, `~` rendering, or harness translation.
   Its `Command` func stays a thin arg-grammar pass-through like every other
   porcelain command. The harness translation is a single-sourced table, not a
   conditional at each call site.
4. **Black-box assertables.** A fixture repo with a handoff whose `## State`
   carries distinctive prose: after `bench handoff`, that prose is byte-identical
   and `## Shape` matches the kit's current text. A fixture with no
   `## State` heading: exit non-zero and the file's bytes unchanged. No file:
   skeleton exists with an empty State. `--harness codex` output contains
   `$bench-` and no `/bench-`; the default contains `/bench-` and no `$bench-`.
   A repo checked out under a non-`workspace` path renders that path, proving no
   layout constant. A repo whose root is outside `$HOME` renders absolute, not
   `~`. Two consecutive runs produce identical bytes. A staged spec is named
   with its Status; no staged spec says so. `--next` appears verbatim.
5. **Gate attachment.** Existing runtime contract phase hosts the CLI behavior
   (this is a porcelain command with a fixture repo, the shape
   `internal/contract/runtime` already runs). Conformance hosts one check: the
   `## Shape` text the command emits is the single source, so no other kit file
   restates it. Unit tests in `internal/handoff` cover section splitting and
   `~` rendering.
6. **Hostile-input owners.** Repeated or unknown flags → the shared arg-grammar
   helper (FT87 #7). A handoff with two `## State` headings, or one inside a
   fenced code block → the section splitter, which fails closed rather than
   picking one. Detached HEAD, no commits yet, no `origin` remote → fact
   collection, each rendering an explicit unknown rather than an empty field.
   Repo root outside `$HOME` → path rendering. Control bytes in branch or spec
   names → existing `toon` refusal. A read-only file system on `--write` → a
   structured error naming the path.
7. **Uncertainty flags.** None. Every fork above was answered by the reviewer in
   the session that wrote this map; the three **[veto]** marks are recorded
   alternatives, not open questions.
8. **Rejected alternatives.** Environment sniffing for the harness (#2). A
   print-only default with `--write` (#1). Preserving `## Shape` verbatim (#4).
   Overwriting a malformed file (#5), and refusing to scaffold a missing one.
   Emitting a complete fresh-session prompt including task framing — the CLI
   cannot derive the task, and a template that invites it would launder guesses
   as facts. Storing the handoff's age inside the file, which
   `internal/status` already rejected for the same remembered-not-computed
   reason.
9. **Domain watch-outs.** `internal/status` computes handoff staleness from git
   history *specifically* to avoid a self-reported date; this command must not
   reintroduce one by stamping the file. The write lands uncommitted, so
   `bench status` keeps reporting the row until the rewrite is committed — that
   is correct, not a bug to paper over. `## Shape` becoming CLI-owned means the
   text in today's `session-handoff.md` moves into the binary; the conformance
   check exists so a second copy cannot survive.

Dependency order: one slice. Fact collection, rendering, and the write are a
single capability whose parts have no independent value, and the whole is small
enough that splitting would cost more gate runs than it saves. Slicing stays the
reviewer's call.
