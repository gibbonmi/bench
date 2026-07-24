# cli-grammar-and-capability-evidence

Status: implemented

## Problem

FT87 slice 3, compiled from tickets #7 and #8 of the closed
`decisions/bounded-network-resource-cli.md` map. Two unrelated-looking defects share
one root: a fact that should have a single owner is instead re-derived at every call
site, so each site drifts on its own.

**Argument parsing is per-command.** Every Go subcommand hand-rolls its own
`switch` over `args`, and the results disagree:

- Trailing garbage is variously rejected and ignored. `bench status --all x` exits 2,
  but `bench maps --count x`, `bench guards --brief x`, `bench dashboard --stdout x`,
  and `bench structure --since base x` silently ignore `x` and do the work anyway.
- Help is not uniformly success. `bench learnings -h` exits 0; `bench commit -h` and
  `bench dashboard -h` exit 2, and `bench commands` has no help form at all.
- No subcommand supports `--`, so a path that begins with a dash is inexpressible.
  `bench commit -m msg -- -weird.txt` is rejected as an unknown flag.
- Spec-slug fallback resolution is anchored inconsistently. `internal/spec`'s own
  subcommands resolve `specs/<slug>.md` against the repository root; `bench coverage`
  passes an empty base, so the same slug resolves against the process cwd and fails
  from any subdirectory.
- `bench commit` accepts a directory path but does not authorize its children. Git's
  `:(literal)dir` pathspec already stages everything under the directory; the
  attribution block-check compares exact path strings, so every changed child reads
  as an unexplained working-tree file and the commit is refused.

**Capability-dependent security tests skip silently.** Tests that probe symlinks,
FIFOs, pid semantics, multi-CPU behavior, privilege, and POSIX signals call bare
`t.Skip` when the host cannot run them. In non-verbose `go test` a skip is
indistinguishable from a pass, so the gate reports green on a host where a whole
class of security assertion never ran. Nothing counts them, and nothing makes a
skipped security class red anywhere — including the release and native workflows,
where every capability is present and a skip means a genuine regression.

**Test deadlines are coupled to the bounds they contain.** `worktree.staleAfter` is a
private one-minute constant; test fixture ages restate it as a literal
`2 * time.Minute`, so the relationship is knowledge duplicated across files. The
two-leg marker wait's slow deadline is a literal `60 * time.Second` — numerically
equal to the staleness window it may have to outlast. An outer timeout equal to the
inner window it contains is a flake by construction.

## Solution

Give each re-derived fact one owner.

**One argument grammar.** A small hand-rolled parsing helper in `internal/usage`
owns arity, flag recognition, `--`, and help semantics for every Go subcommand.
Rules, uniformly: exact arity with trailing garbage rejected as a usage error at
exit 2; `--help` and `-h` always exit 0, and bare `help` does so as the sole
argument; `--` ends flag parsing so
leading-dash positionals are expressible; unknown flags exit 2 through the existing
`toon.Usage` rendering. A conformance check keeps the next subcommand from
hand-rolling its own parse. `bench coverage` adopts the `internal/spec` package's
root anchoring so slug fallback resolves identically from any directory, and
`bench commit`'s path grammar gains the conventional directory rule: naming a
directory authorizes its changed children.

**Capability skips become evidence.** One helper package replaces every bare
`t.Skip` in the repository. A capability skip emits a structured, recognizable line
naming its class and reason; an environment skip (absent subject binary, unset
conformance root) emits the same shape under a distinct kind. The line is appended to
a skip log named by an environment variable, not written to stdout: non-verbose
`go test` discards a passing package's stdout, so a stdout line would be invisible on
exactly the runs that matter. The gate reads that log after its phases join and
aggregates the lines into explicit `capability-skips` rows in its own output. On dev
machines those rows are visible and non-fatal. Under an explicit
strict flag — set by the release and native workflows — a nonzero capability-skip
count is red.

**Deadlines decouple by construction.** The lease-staleness window moves into
`internal/bounds` as the one named policy constant, and test fixture ages derive from
it instead of restating it. A single `bounds.TestDeadline` helper derives every outer
test deadline from the inner bound it must outlast, and is provably larger than its
argument. A conformance check forbids a literal duration at the marker-wait call
sites, so the next deadline cannot be re-hardcoded.

## User stories

1. As a Bench maintainer, I want one hand-rolled parsing helper in `internal/usage`
   that owns arity, flag recognition, `--`, and help semantics, so that the CLI's
   argument grammar has a single source instead of one hand-rolled `switch` per
   subcommand.
   Line: `gpt-5.6-luna` / medium. Every rule of the grammar is enumerated in this
   spec and every one of them is an exit code the gate observes directly, so the
   cheap tier is building to a fully specified target at a named seam.

2. As a Bench maintainer, I want every Go subcommand to route its parsing through
   that helper, and a conformance check that fails when a subcommand hand-rolls its
   own, so that the grammar cannot drift back apart the next time a subcommand is
   added.
   Line: `gpt-5.6-terra` / medium. The profile routes gate and conformance logic to
   the mid tier because a conformance check that does not bite is worse than no
   check at all.

3. As an agent driving the CLI, I want trailing garbage after an otherwise satisfied
   grammar to be a usage error at exit 2 on every subcommand, so that a mistyped
   invocation reports itself instead of quietly doing something adjacent to what I
   asked for.
   Line: `gpt-5.6-luna` / low. Once the helper exists this is mechanical rewiring of
   call sites whose correct exit codes are pinned by contract assertions.

4. As an agent driving the CLI, I want `help`, `--help`, and `-h` to exit 0 on every
   Go subcommand, so that asking a command how to use it is never reported as my
   having used it wrong.
   Line: `gpt-5.6-luna` / low. The behavior is one rule applied uniformly and each
   application is a single asserted exit code.

5. As an agent driving the CLI, I want `--` to end flag parsing on every subcommand
   that takes positionals, so that I can name a path beginning with a dash instead of
   having it rejected as an unknown flag.
   Line: `gpt-5.6-luna` / low. The rule lives entirely inside the helper from story 1
   and every call site inherits it without its own decision.

6. As an agent running `bench coverage` from a subdirectory, I want a bare spec slug
   to resolve against the repository root exactly as it does from the root, so that
   where I happen to be standing does not change whether the command finds the spec.
   Line: `gpt-5.6-luna` / medium. The fix is adopting the anchoring the `internal/spec`
   package already implements, and the equality of the two invocations is directly
   assertable.

7. As a reviewer committing a change that spans a directory, I want naming that
   directory to authorize its changed children in `bench commit`, so that the
   conventional path grammar works instead of the attribution block-check refusing
   every file under a path I explicitly named.
   Line: `gpt-5.6-terra` / medium. This deviates from the profile's cheap CLI-plumbing
   row on purpose, because the code being widened is the attribution guard that makes
   a green gate describe exactly the diff that lands.

8. As a Bench maintainer, I want hostile filenames — leading dashes, embedded spaces,
   and glob characters — exercised through the real binary for the commit and outline
   path surfaces, so that the grammar's promises hold for the filenames the profile's
   hostile-input checklist says a shell CLI actually meets.
   Line: `gpt-5.6-luna` / medium. The fixtures are mechanical to construct and every
   assertion is an observed exit code or output line.

9. As a Bench maintainer, I want one capability-skip helper that replaces every bare
   `t.Skip` in the repository and emits a structured line naming the skip's kind,
   class, and reason, so that a skipped test is visible evidence rather than a result
   indistinguishable from a pass.
   Line: `gpt-5.6-luna` / medium. The call-site sweep is broad but entirely
   mechanical, and the conformance check from story 2's family proves the sweep was
   complete.

10. As a reviewer reading a gate run, I want the gate to aggregate those lines into
    explicit `capability-skips` rows in its output, so that I can see which security
    classes did not run on this host without reading every phase's raw output.
    Line: `gpt-5.6-terra` / medium. Aggregation reads phase output as the phases
    stream concurrently, which is oracle-adjacent code the profile routes to mid.

11. As a release maintainer, I want a nonzero capability-skip count to be red under an
    explicit strict flag that the release and native workflows set, so that a security
    class silently skipped on a fully capable runner fails the release instead of
    shipping.
    Line: `gpt-5.6-terra` / medium. This changes what makes the oracle red, which is
    the class of change the profile keeps off the cheap tier.

12. As a Bench maintainer, I want the lease-staleness window to live in
    `internal/bounds` as one named constant that test fixture ages derive from, so
    that the staleness policy is stated once instead of being restated as a literal in
    every fixture that has to outlast it.
    Line: `gpt-5.6-luna` / low. A constant move plus derived call sites, with the
    existing lifecycle tests as the regression signal.

13. As a Bench maintainer, I want every outer security-test deadline derived from one
    helper that is provably larger than the inner bound it contains, and a conformance
    check that forbids a literal duration at those call sites, so that an outer timeout
    can never again be set equal to the window it exists to outlast.
    Line: `gpt-5.6-terra` / medium. The check is conformance logic, and its value is
    entirely in whether it actually bites the literal it forbids.

## Implementation decisions

**The grammar helper lives in `internal/usage`.** The map named "a shared arg-grammar
helper" as a module boundary without naming its package; `internal/usage` already owns
CLI usage strings and has no importers that would cycle through `internal/toon`, so it
is the natural home rather than a new package. This is the one seam choice the map left
open and the one thing in this spec to veto if it is wrong.

**The helper is a spec-driven parser, not a framework.** A subcommand declares its
grammar — the flags it accepts (with or without values), its positional arity, and its
help text — and receives back the parsed result or a rendered usage error with its exit
code. Rendering stays in `internal/toon` (`Usage`, `MissingArg`), which remains the one
source of the usage-line shape. No third-party CLI framework; the map records that
rejection.

**Grammar semantics, in full.** `--help` and `-h` anywhere before `--` print the
declared help text and exit 0; bare `help` does so only as the sole argument, because
a variadic grammar's free text or path list can legitimately contain the word and
recognizing it anywhere would print usage while silently discarding the rest of the
invocation. `--` ends flag parsing; every later argument is a positional, including
one that begins with a dash. An unrecognized flag is `toon.Usage` at exit 2. A flag
requiring a value with no value following is `toon.MissingArg` at exit 2. A declared
flag given twice is `toon.Usage` naming that flag at exit 2, not last-one-wins: no
flag in this CLI accumulates a list, so the only thing silently keeping the later
value buys is hiding a mistyped invocation whose two spellings disagree. An empty
positional is `toon.Usage` at exit 2 — it names nothing, it is what an unset shell
variable expands to inside quotes, and a subcommand that resolves it against the
filesystem widens to the cwd. More positionals than the declared arity is
`toon.Usage` on the first excess argument at exit 2 — this is the trailing-garbage
rule. Fewer than the required arity is `toon.MissingArg` at exit 2. A bare `-` is an
ordinary positional, not a flag and not a separator.

**Conformance enforces the routing.** A check in `internal/conformance` fails when a
Go subcommand's entry point parses `args` itself rather than through the helper. The
check is fail-closed against new subcommands: the mechanism is an explicit registry of
entry points cross-checked against the dispatch table in `cmd/bench/main.go`, so a new
name in that table with no registry entry is red rather than unexamined.

**`bench coverage` adopts root anchoring.** `internal/spec` exposes the repo-root
resolution its own subcommands already use, and `coverage` calls it instead of passing
an empty base. Resolution keeps its current fail-safe posture outside a repository.

**`bench commit`'s directory rule lives in the block-check, not in staging.** Staging
already works: `git add -A -- :(literal)dir` stages every changed path under the
directory. Two places change. The attribution allow-set in `unexplained` treats a named
directory as authorizing any working-tree path under it, matched on path segments so
that naming `sub` never authorizes `subdir/x`. `stagePlan` classifies a named directory
as stageable when it exists in the worktree, and reports the same real error as today
when a named path exists nowhere. The block-check stays exact for file paths.

**One skip helper, two kinds.** A new package owns skipping. `Capability(t, class,
reason)` covers the enumerated capability classes — symlink, fifo, pid, cpu, privilege,
signal, and tool — and is what the strict flag counts. `Environment(t, reason)` covers
the non-capability skips that already exist: absent subject binary, unset conformance
root, unmaterialized canary fixture. The existing `contract.SkipIfSubjectBenchMissing`
and `SkipIfSubjectFileMissing` route through `Environment`, so their many call sites
are unchanged. Both write their structured line to stdout before calling `t.Skip`,
because a skip message alone is invisible under non-verbose `go test`.

**Conformance forbids the bare form.** After the sweep, no `t.Skip` or `t.Skipf` may
appear in the repository outside the helper package itself. This is a fail-closed
allowlist of exactly one owner, which is what makes the sweep provably complete rather
than best-effort.

**The gate aggregates from a skip log.** The gate names a log file in the environment
it hands its phases; the skip helper appends one line per skip to it, and after the
phases join the gate reads the file and tallies the lines by kind and class. The
transport is a file rather than a tee of phase stdout because `go test` without `-v`
discards a passing package's stdout, which is where most skips happen. Phases run
concurrently and each skip is one append under the atomic write size, so the
concurrency is the filesystem's rather than a shared in-process collector's. The gate
strips the log variable from the environment it gives an inner gate run, so a canary
fixture's deliberate skips cannot contaminate the outer tally. Zero skips still emits
a row stating zero — a definitive empty state, not silence, because absent output and
zero skips must not look alike.

**Strict mode is an explicit environment flag.** `BENCH_REQUIRE_CAPABILITIES=1` makes
a nonzero capability count red with its own distinct message, and the release and
native workflows set it on the step that runs the gate. Environment skips never
contribute to the strict count: an absent subject binary is a staging fact, not a
security class. Absent or any value other than `1`, the flag is off and the rows are
informational — the same fail-safe posture the other Bench flags use.

**Deadlines.** `bounds.LeaseStale` becomes the one named staleness window and
`internal/worktree` consumes it. `bounds.TestDeadline(inner)` derives an outer test
deadline from the inner bound it must contain, and is defined so its result is
strictly greater than its argument for every bound in the policy registry. Marker-wait
call sites pass `bounds.TestDeadline(...)` instead of a literal, and a conformance
check fails on a numeric duration literal passed to `WaitForTwoLegMarkers`.

## Testing decisions

A good test here drives the real surface and reads its observable result: an exit code
and stdout from the built `bench` binary for CLI grammar; the gate's own output for
capability evidence; a returned value for the pure helpers. No test reaches into a
command's internal parse state.

Prior art in this repo, followed rather than reinvented:

- `internal/toon/toon_test.go` — table-driven unit tests over pure rendering helpers.
  The grammar helper's tests take the same shape.
- `internal/contract/axi/*_test.go` — exit-code and stdout matrices against the built
  binary in throwaway fixture repos. The grammar's per-command matrix belongs here.
- `internal/contract/runtime/runtime_commit_test.go` — commit behavior driven through
  the real binary against a real fixture repo. The directory and hostile-path stories
  extend it.
- `internal/conformance/checks_test.go` and its siblings — repository-wide
  single-sourcing checks, each paired with a canary fixture that proves it bites.

Gate command: `.bench/gate.sh` (the project gate).

The runtime contract phase executes the built `dist/bench`, so the build must be
current before those tests are hand-run outside a full gate.

### Seam diagram

**Seam 1 — the grammar helper (`internal/usage`).** Pure; no repo, no process.

    trigger: a subcommand's Command func, at entry
        │
        ▼
    grammar spec  ──▶  [ usage.Parse                      ]  ──▶  parsed flags + positionals
    raw []string  ──▶  [   help / `--` / arity / unknown  ]  ──▶  rendered usage line + exit code
                           ◀ tests attach here: call Parse with a declared grammar and an
                             argv slice; assert the returned result, line, and code

**Seam 2 — the CLI surface (`internal/contract/axi`, `internal/contract/runtime`).**
The built binary in a throwaway fixture repo.

    trigger: the contract phase execs dist/bench in a fixture repo
        │
        ▼
    argv          ──▶  [ bench <subcommand>               ]  ──▶  stdout
    fixture repo  ──▶  [   (grammar + commit + coverage)  ]  ──▶  exit code
    cwd           ──▶  [                                  ]  ──▶  repo state after the run
                           ◀ tests attach here: run the binary with a hostile argv from a
                             chosen cwd; assert exit code, stdout, and git state

**Seam 3 — capability evidence (`internal/capability`, `internal/gate`).**

    trigger: a test lacks a capability  │  the gate runs its phases
        │                               │
        ▼                               ▼
    class + reason ─▶ [ capability.Capability ] ─▶ structured stdout line ─▶ [ gate collector ]
                          │                                                        │
                          ▼                                                        ▼
                       t.Skip                                        `capability-skips` rows
                                                                     + strict-mode verdict
                           ◀ tests attach here: unit-test the line renderer directly; drive the
                             collector with scripted phase output; assert the gate's own stdout
                             on a fixture whose phase emits skip lines

**Seam 4 — single-sourcing (`internal/conformance`).** Repository-wide static checks.

    trigger: the conformance phase, with BENCH_CONFORMANCE_ROOT set
        │
        ▼
    repo tree  ──▶  [ conformance checks: routing registry,  ]  ──▶  violations
                    [ no bare t.Skip, no literal deadlines   ]  ──▶  exit code
                        ◀ tests attach here: run each check against a fixture tree carrying the
                          violation and assert it is reported; canary fixtures prove it bites

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `--` ends flag parsing: a later `-x` is returned as a positional, not an unknown flag | `internal/usage` unit | new table test asserting the positional — red today, the helper does not exist | The leading-dash path the whole grammar exists to make expressible is unreachable without this rule, so a helper that forgets `--` fails here rather than at one call site |
| 1 | `help`, `--help`, `-h` before `--` each return the declared help text at exit 0 | `internal/usage` unit | same table test over the three spellings — red against an implementation that handles only `--help` | Three spellings are enumerated because the tree today disagrees on which are honored; a partial implementation passes a one-spelling test |
| 1 | A flag requiring a value with nothing after it returns `MissingArg` at exit 2, distinct from an unknown flag's `Usage` | `internal/usage` unit | same table test asserting the rendered line and code — red against an implementation that collapses both into one message | Collapsing the two labels is the cheapest wrong implementation and it still exits 2, so only asserting the rendered line catches it |
| 1 | A bare `-` is a positional, and `--` appearing twice makes the second one a positional | `internal/usage` unit | same table test — red against an implementation that treats any dash-prefixed token as a flag | These are the two tokens a naive prefix check gets wrong, and both are legal filenames |
| 1 | Bare `help` is help only as the sole argument: in a variadic grammar `help me remember` returns three positionals, and `help` after another positional stays a positional | `internal/usage` unit | table rows over both placements — red against an implementation that recognizes `help` anywhere before `--` | Recognizing it anywhere exits 0 while discarding the invocation, so the failure is silent data loss rather than a visible error, and only the variadic placement distinguishes the two implementations |
| 1 | A declared flag given twice returns `Usage` naming that flag, for value and boolean flags alike, and before any value is consumed | `internal/usage` unit | same table test including a repeated value flag with no second value — red against last-one-wins | Last-one-wins still exits 0 on the invocation that meant two different things; asserting the flag is named excludes a generic error, and the missing-second-value row pins that the duplicate check precedes the value read |
| 1 | An empty positional returns `Usage` at exit 2 rather than being accepted as a path | `internal/usage` unit, then `internal/contract/runtime` | unit row plus a contract run of `bench commit -m msg ""` from a subdirectory — red today, the empty token resolves to the cwd and stages every changed file beneath it | The unit row pins the rule at its single owner; the contract run is what proves no path-taking subcommand re-derives the guard, since the damage is a commit that lands files nobody named |
| 2 | A subcommand whose entry point parses `args` itself, with no registry entry, is reported by the conformance check | `internal/conformance` | new check plus a canary fixture adding an unregistered subcommand — red against a check that only inspects registered names | A check that reads only its own registry passes vacuously on exactly the case it exists to catch: the newly added subcommand |
| 3 | `bench maps --count x`, `bench guards --brief x`, `bench dashboard --stdout x`, and `bench structure --since base x` each exit 2 with a usage line naming `x` | `internal/contract/axi` | extend the AXI matrix with these four invocations — red today, all four exit 0 and do the work | Naming all four is the quantifier: they are the complete set of subcommands that ignore trailing garbage today, so a fix applied to one still fails |
| 3 | The same four invocations without the trailing argument still exit 0 and produce their normal output | `internal/contract/axi` | same matrix — red against an implementation that rejects the legitimate form while tightening arity | Excludes the degenerate fix that makes the new assertion pass by rejecting the flag outright |
| 4 | `bench commit -h`, `bench dashboard -h`, `bench commands -h`, and `bench coverage -h` each exit 0 and print help | `internal/contract/axi` | extend the matrix — red today: commit and dashboard exit 2 and commands has no help form | These are the enumerated subcommands that violate the rule today, so the row cannot be satisfied by the ones that already comply |
| 5 | `bench commit -m msg -- -weird.txt` commits a file literally named `-weird.txt` | `internal/contract/runtime` | new contract test with that fixture file — red today, the argument is rejected as an unknown flag | The end-to-end path is what proves the grammar rule survived every layer between argv and the git pathspec |
| 6 | `bench coverage <slug>` run from a subdirectory produces byte-identical stdout and the same exit code as from the root | `internal/contract/axi` | new contract test running the binary from a nested cwd — red today, the subdirectory run fails to find the spec | Comparing the two runs rather than asserting success excludes an implementation that resolves both against the cwd and happens to work from the root |
| 7 | `bench commit -m msg <dir>` with two changed files under `<dir>` commits both and is not blocked | `internal/contract/runtime` | new contract test — red today, the block-check reports both children as unexplained | Two changed children rather than one excludes an implementation that special-cases a directory holding exactly one file |
| 7 | Naming `sub` does not authorize a changed file under a sibling directory named `subdir` | `internal/contract/runtime` | same test with both directories dirty — red against a prefix-string implementation | String-prefix matching is the cheapest wrong implementation of the directory rule and it silently widens what a commit stages |
| 7 | With a changed file outside the named directory, the commit is still blocked and names that file | `internal/contract/runtime` | same test — red against an implementation that widens the allow-set to everything once any directory is named | This is the safety property the whole block-check exists for; the widening story 7 grants must not become unconditional |
| 8 | `bench commit` and `bench outline` handle paths containing a space and a glob character (`*`) as literal names | `internal/contract/runtime` | new fixtures with both filenames — red against an implementation that lets the shell or git glob-expand them | The profile's hostile-input checklist names both classes; a `:(literal)` pathspec that is dropped anywhere in the chain regresses to expansion |
| 9 | A capability skip emits its structured line — kind, class, reason — on stdout and then skips | `internal/capability` unit | new test capturing the writer and asserting the line — red today, the package does not exist | Writing the line after the skip, or only into the skip message, produces nothing under non-verbose `go test`; asserting the emitted bytes is what separates the two |
| 9 | A capability class outside the enumerated set is rejected rather than emitted | `internal/capability` unit | same test with an unknown class — red against an implementation that formats whatever it is given | An open class vocabulary makes the strict count ungradeable, since a typo would silently create a class nothing counts |
| 9 | No `t.Skip` or `t.Skipf` appears in the repository outside the helper package | `internal/conformance` | new check plus a canary fixture reintroducing a bare skip — red today, the tree has many | The sweep is only provable by a check that fails on a single reintroduction; a spot-check of known files cannot make that claim |
| 10 | The gate's output carries a `capability-skips` row for each class a phase reported | `internal/gate` unit | new test driving the collector with scripted phase output — red today, no collector exists | Driving the collector directly pins the aggregation rather than the far weaker claim that some rows appeared |
| 10 | With zero skips the gate still emits a row stating zero | `internal/gate` unit | same test with clean phase output — red against an implementation that emits rows only when a skip occurred | Silence and zero must be distinguishable, or a collector that silently broke reads exactly like a fully capable host |
| 10 | Lines emitted by concurrently running phases are all counted | `internal/gate` unit | same test with interleaved scripted output across phases — red against an unguarded collector under `-race` | The phases genuinely run concurrently, so an unsynchronized tally loses counts on exactly the multi-phase runs that matter |
| 11 | With `BENCH_REQUIRE_CAPABILITIES=1` and a phase reporting one capability skip, the gate exits non-zero naming the class | `internal/contract/runtime` | new contract test against a fixture phase that emits a skip line — red today, the flag has no meaning | The whole point is the verdict flip, and asserting the class name in the message excludes a generic red that says nothing actionable |
| 11 | With the flag unset, the same fixture is green and the rows still appear | `internal/contract/runtime` | same test — red against an implementation that makes any skip red unconditionally | Dev machines legitimately lack capabilities; an unconditional red would make the gate unusable locally |
| 11 | An environment skip does not make the strict mode red | `internal/contract/runtime` | same test with a fixture emitting only environment lines — red against an implementation that counts every structured skip line | Environment skips are staging facts and are common in CI; counting them would make strict mode fire on a condition it does not describe |
| 11 | The release and native workflows set the flag on the step that runs the gate | `internal/conformance` | extend the existing release-workflow structure check — red today, the flag is absent | The behavior is worthless if no workflow turns it on, and a workflow file is exactly the surface a static check can grade |
| 12 | `internal/worktree` reads its staleness window from `bounds.LeaseStale` and the lifecycle tests derive their fixture ages from the same constant | `internal/worktree` unit | existing lifecycle table re-expressed against the constant, then run with the constant temporarily changed — red against a tree that still hardcodes `2 * time.Minute` | Changing the constant and seeing the fixtures follow is the only signal that distinguishes a real derivation from a literal that merely equals it today |
| 13 | `bounds.TestDeadline(b)` is strictly greater than `b` for every bound in the policy registry | `internal/bounds` unit | new test enumerating the registry — red against an implementation returning the bound unchanged | Enumerating the registry rather than a sample means a bound added later without a margin fails here |
| 13 | A numeric duration literal passed to `WaitForTwoLegMarkers` is reported by the conformance check | `internal/conformance` | new check plus a canary fixture with a literal call — red today, both call sites pass literals | The literal is exactly what regressed before; a check that only inspects the helper's own definition would not see the call site |

### Edge inventory

Walked against the profile's hostile-input checklist for the surfaces this spec
touches, with each class landing as a coverage row above or a **Won't handle** line
here.

- Paths with spaces and glob characters — story 8 rows.
- Leading-dash paths — story 5 row.
- Cwd deeper than the repo root — story 6 row.
- Absent vs. present-but-empty — the zero-skips row under story 10 is this class for
  the evidence surface.
- Re-run idempotency — the coverage-from-subdirectory row compares two runs, which is
  the idempotency assertion for that surface.
- Required tool missing from PATH — the `tool` capability class exists precisely for
  this, and story 9's class-vocabulary row covers its handling.
- Control bytes in git-sourced text — **Won't handle**: `toon.Table` already refuses
  them and this spec adds no new path from git text into a table.
- Non-TTY stdin — **Won't handle**: no surface in this slice prompts.
- Interrupt mid-loop — **Won't handle**: cancellation is slice 1's shipped work and
  this slice adds no new long-running subprocess.
- Special files (FIFOs, devices) in discovery paths — **Won't handle** for the commit
  directory rule: a named directory's children are enumerated by git, which lists only
  tracked and untracked regular paths, so a FIFO is never a candidate.
- Symlinked directory named as a commit path — **Won't handle**: git resolves the
  pathspec against the index, where a symlink is a blob and not a directory, so it
  takes the existing file path.
- Invocation through every shipped surface — **Won't handle**: the grammar helper sits
  inside the Go binary that every surface already routes to, so the existing
  routed-implementation checks cover reachability unchanged.

## Out of scope

- **Argument grammar for the shell-side router (`bin/bench.sh`) itself.** A separate
  capability: the shell router dispatches to the Go binary and to shell-owned paths,
  and giving it the same grammar means a POSIX-shell parser with its own portability
  rules rather than a reuse of this helper. Estimate: 6 edits, 3 gate runs.
- **Extending strict capability mode to linked consumer repos' gates.** A separate
  capability: it needs the flag, the class vocabulary, and the aggregation to be part
  of the shipped consumer surface rather than this repo's own gate, which is a
  packaging decision this slice does not open. Estimate: 5 edits, 3 gate runs.
