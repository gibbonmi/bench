# Resolved consumer surface

Status: staged

Roadmap: FT191

Decision source: named reviewed artifact `decisions/assets/ft191-resolved-reader-research.md` (reviewer-routed 2026-08-28; scope widened to all of FT191 plus the review blast list, reviewer, 2026-08-28).

Verification log: 2 iteration(s) to accept — round 1 returned 14 blocking findings, folded in full. Round 2 verified every fold and returned 2 small blocking findings plus 5 non-blocking ones. The author folded all 7 at the declared 2-loop cap. The citation row moved before the terminal help envelope, and the wrapper route returned to the fences.

## Problem

An agent proves a universal consumer claim with sampled greps today.
`bench outline` locates candidate declarations only, and the agent confirms
each row by hand. The FT210 miss shows the cost: a planner-local counter
stayed green while two commands reached target Git through older consumers.
Charges also lack the helper, double, and fixture inventory FT191 names, so
each charge rebuilds it by exploration.

## Solution

`bench consumers` resolves real Go reference edges and emits them as TOON
tables. A blast mode enumerates the consumers of every symbol a diff
touched, over the frozen review pair. Every response carries a replayable
citation row before its help envelope. `bench outline` gains candidate rows for helpers, doubles, and
fixtures. The review phase runs the blast enumeration and walks it before
the landing.

## User stories

### Group A — resolved consumer queries

Line: opus / low. The reviewer's 2026-08-28 routing starts every charge at
opus low.

1. As an agent, I want `bench consumers <symbol>` to list every reference to
   a resolved Go symbol, so that a universal consumer claim cites an
   enumeration.
2. As an agent, I want each row classified as call, reference, or implements,
   so that an invocation reads apart from a registration.
3. As an agent, I want each row to name its enclosing declaration, so that a
   registry entry such as `commandRegistry` is visible as one.
4. As an agent, I want an ambiguous bare name answered with a candidates
   table, so that I re-query qualified instead of guessing.
5. As an agent, I want a matched symbol with zero references answered with
   the definitive empty table, so that the absence claim is warranted.
6. As an agent, I want a symbol in a non-Go file refused with the language
   named, so that an empty table never fakes "no callers".
7. As an agent, I want an ill-typed tree refused with the first error named,
   so that no enumeration rests on a broken tree.
8. As an agent, I want a missing `go` binary refused with the remedy named,
   so that the failure is actionable.
9. As a reviewer, I want one citation row (SHA, state, version, argv, hash)
   before the help envelope, so that cited evidence replays.
10. As a reviewer, I want a dirty-checkout citation marked dirty, so that an
    unreplayable citation says so.
11. As an agent, I want two runs at one SHA byte-equal, so that a replayed
    citation hash matches.
12. As an agent, I want an alias-spelled query and an origin-spelled query
    to emit identical bytes, so that alias spelling cannot vary the output.
13. As an agent, I want `--full` to emit the complete row set, so that a
    universal claim covers the whole set.
14. As an agent, I want usage on stdout at exit 2 and structured errors at
    exit 1, so that the surface honors AXI.
15. As an agent, I want the help text to state the promise clause and the
    soundness limits, so that I never over-claim.
30. As an agent, I want candidates rows to carry the qualified re-query
    spelling and one literal re-query action each, so that recovery is one
    paste.
32. As an agent, I want an over-cap default answered with per-package
    aggregates and one `--full` action, so that a hot symbol never floods my
    context.

### Group B — blast mode

Line: opus / low. The reviewer's 2026-08-28 routing starts every charge at
opus low.

16. As a reviewer, I want `bench consumers --changed --base <b> --source-tip <t>`
    to enumerate consumers of every touched symbol, so that review walks the
    FT210 class.
17. As a reviewer, I want a declaration the diff deleted reported in its own
    deleted table without enumeration, so that the walk still sees it.
18. As a reviewer, I want an empty diff answered with the definitive empty
    blast table, so that a no-op change reads as one.
19. As a reviewer, I want blast rows derived only from the frozen pair, so
    that a recomputation at review time byte-matches.
20. As an agent, I want `--changed` to mirror the `bench test` base and
    source-tip grammar, so that one revision grammar serves both.
31. As a reviewer, I want a blast row marked `touched` when its file sits
    inside the diff's set, so that I walk outside-diff rows first.

### Group C — outline inventory

Line: opus / low. The change extends a known package under existing tests.

21. As a charge author, I want outline rows with kind `helper` for test-file
    helper functions, so that a charge carries the helper inventory for free.
22. As a charge author, I want kind `double` rows for fake, stub, mock, and
    spy names, so that a charge sees the doubles.
23. As a charge author, I want one kind `fixture` row per file under a
    `testdata/` segment, so that prior-art fixtures surface with `file:line`.
24. As a reviewer, I want the outline kind delta pinned by an old-to-new
    fixture pair, so that the byte change is exactly the reviewed one.
25. As an agent, I want outline's LOCATE promise unchanged, so that the new
    rows stay candidate-class.

### Group D — advertisement and review wiring

Line: opus / low. The reviewer's 2026-08-28 routing starts every charge at
opus low.

27. As a reviewer, I want `/bench-review-implementation` to run the blast
    enumeration and walk it, so that unlisted consumers stop slipping
    through review.
28. As a cold agent, I want `bench help` and the project profile to
    advertise the new surface, so that discovery needs no guesswork.
29. As a reviewer, I want the AXI registry disposition set for
    `bench consumers`, so that conformance grades the new surface.

Story numbers 30 through 32 were added in later passes, and story 26 was
removed as an unflagged addition. The numbering stays stable rather than
reflowing.

## Implementation decisions

- New package `internal/consumers` owns the surface. A pure analysis core
  consumes typed packages and returns rows. A loader seam wraps
  `go/packages`; the seam is package-internal and joins no audited port
  registry, whose set is closed.
- The engine is `golang.org/x/tools` (`go/packages`, `go/types`), with
  `x/mod` and `x/sync` indirect. The dependency fits the AGENTS.md standard:
  an official org, BSD-3, and a build-time-only footprint.
- The symbol grammar accepts an import-path-suffix qualification
  (`outline.Command`, `Type.Method`) and a bare identifier. A bare name with
  several matches answers with a candidates table at exit 0.
- `via` has three values: `call`, `reference`, `implements`. A queried
  interface lists its satisfying types as `implements` rows. The enclosing
  declaration column carries registry context; there is no second registry
  derivation.
- Row identity uses origin objects: alias-resolved, generic-origin, with
  positions from the file set. Rows sort by file, line, then column. No
  printed instantiation name reaches the output.
- The citation row is `citation{sha,state,version,cmd,hash}`. It emits
  immediately before the help envelope, because the AXI contract pins `help`
  as the terminal block. The hash is sha256 over all output bytes before the
  citation row. `state` is `clean` or `dirty`, and `version` is the bench
  version.
- The default response is complete up to a cap of 200 rows. Over the cap,
  the default emits per-package aggregate rows, `truncated=true` in meta,
  and one `--full` action carrying the known symbol. `--full` always emits
  every row, and the review blast step runs `--full`.
- Blast derivation is a pure function of the hunk list and the tip packages.
  Enumeration runs at the tip. A deleted declaration emits one
  `blast_deleted[N]{changed_symbol,base_file,base_line}` row. The git
  queries stay at the rim, through `internal/git`.
- `touched` is true when the row's consumer file sits in the diff's
  changed-file set. The outside-diff rows are the FT210 review signal.
- Refusals are fail-closed structured stdout errors at exit 1: an ill-typed
  tree, a missing `go` binary, and a non-Go symbol file. Refusals append no
  disclosure.
- Row schemas stay minimal per AXI: `consumers[N]{file,line,via,enclosing}`
  (the queried symbol is constant, so no symbol column),
  `consumers_candidates[N]{qualified,file,line,kind}`,
  `blast[N]{changed_symbol,file,line,touched}`, and
  `blast_deleted[N]{changed_symbol,base_file,base_line}`. The blast detail
  path is the per-symbol query, not extra columns.
- Every response ends with the `help[N]{cmd,why}:` envelope through the
  existing `internal/axi` action owner. A symbol result and an empty blast
  are terminal reads (`help[0]`). A candidates result offers one literal
  re-query action per candidate row, because every argument is known. A
  blast result offers one `bench consumers <symbol> --full` action per
  changed symbol with an untouched consumer, deduplicated in stable order.
- Registration is graded in four places in one commit. The places are
  `commandRegistry` with the AXI disposition, `approvedAXIQueries`,
  `subcommandRouting`, and the approved-query tables. The conformance pair
  lives in `internal/conformance/axi_query_registry_test.go` and
  `internal/conformance/subcommand_routing_test.go`; the tables live in
  `.agents/skills/bench-craft-cli/SKILL.md` and `projects/benchkit.md`. The
  first command-surface ticket owns all four. The wrapper adds a
  `consumers` porcelain route beside the other read verbs, so the repair
  grant matches them. Nothing grades the wrapper line, so no row exists.
- Outline helper and double rows are pattern-table entries with an optional
  path predicate. The helper forms are `_test.go` functions named with the
  `new`, `make`, or `with` prefix before an upper-case letter. The double
  forms are names with a case-insensitive `fake`, `stub`, `mock`, or `spy`
  prefix. Fixture rows come from a walk-level path classifier beside the
  extension dispatch, because `testdata/` files carry no scanned extension.
- The meta table carries packages, files, matches, rows, and truncated. The
  default build context is the only graded configuration.
- `internal/consumers` splits by responsibility under the structure budgets:
  loader, core, blast, citation, and command files each stay under the
  400-line cap. No budget grant is expected.

## Testing decisions

- A good test drives the public command behavior: argv in, TOON bytes and
  exit code out. Core tests type-check fixture source in process
  (`go/parser` plus `go/types`) and never spawn a subprocess.
- One focused loader test runs the real `go list` path against a minimal
  fixture module; it is the single subprocess site.
- The determinism probe runs the command twice at one SHA and compares
  bytes. The alias row compares an alias-spelled query against an
  origin-spelled query.
- The outline delta is proved by the old-to-new fixture pair; the existing
  outline package tests keep every unchanged byte pinned.
- The landing gate observes the feature through the ordinary `test` phase
  plus the `axi-query-registry` and `subcommand-routing` conformance
  checks.

### Seam diagram

    trigger: agent query, or the review phase's blast step
        │
        ▼
    argv ──▶ [ usage.Parse ] ──▶ [ loader seam (go/packages) ] ──▶ [ consumers core ] ──▶ result tables + citation
                                        ◀ tests attach here: in-process typed fixtures drive the core;
                                          one focused test drives the real loader; command tests inject the seam

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| CS1 | 1 | a qualified query over the typed fixture emits every planted reference row | consumers core, in-process typed fixture | a grep-shaped fake misses the renamed-import reference the fixture plants |
| CS2 | 2 | a call-position use emits `via=call` and a value use emits `via=reference` | consumers core | a classifier that labels everything `reference` reds the call fixture |
| CS3 | 2 | an interface query emits `via=implements` rows for its satisfying fixture types | consumers core | a uses-only walk finds no identifier use and emits nothing |
| CS4 | 3 | each row names its innermost enclosing named declaration | consumers core | a func-only encloser cannot name the planted var-declaration consumer |
| CS5 | 4 | a bare name with two fixture matches emits the candidates table and zero consumer rows | command surface, injected loader | a first-match guesser emits consumer rows the fixture forbids |
| CS6 | 5 | a matched symbol with zero references emits the definitive empty consumers table | command surface | a fake that errors on empty results cannot warrant absence |
| CS7 | 6 | a query resolved to a `.ts` file emits a structured stdout refusal at exit 1 naming the language | command surface | an empty-table answer states the false "no callers" the research forbids |
| CS8 | 7 | a fixture with one type error emits a structured stdout refusal at exit 1 naming that first error position | loader rim, fixture module | a tolerant loader enumerates over an ill-typed tree and under-reports |
| CS9 | 8 | with `go` absent from `PATH` the command emits a structured stdout refusal at exit 1 naming the missing binary | loader rim, env-injected | a raw exec error surfaces as an unactionable stack |
| CS10 | 9 | every success response carries one `citation{sha,state,version,cmd,hash}` row before the terminal help envelope | command surface | a response without the row leaves the claim uncitable |
| CS11 | 10 | a dirty checkout emits `state=dirty` in the citation row | command surface, dirty fixture repo | a clean-marked dirty citation promises a replay that cannot match |
| CS12 | 11 | two identical runs at one SHA produce byte-equal output | determinism probe at the command seam | unsorted map iteration leaks ordering into the bytes |
| CS13 | 12 | an alias-spelled query and an origin-spelled query emit byte-identical tables | consumers core, alias fixture | a resolver that skips `types.Unalias` emits different rows for the two spellings |
| CS14 | 32 | a symbol with 3 planted references emits 3 rows and `rows=3` in meta | command surface | a fake cap below the common case truncates the usual answer |
| CS15 | 14 | an unknown flag exits 2 with usage on stdout | usage grammar | a stderr or exit-1 usage breaks the AXI taxonomy |
| CS16 | 15 | the help text carries the soundness clause verbatim | command help test | a help without the limit invites reflection-blind over-claims |
| CS17 | 30 | each candidates row carries the exact qualified re-query spelling | command surface | a display-only name forces the agent to reconstruct the argument |
| CS18 | 30 | the candidates response ends with one literal re-query action per candidate row | command surface, axi owner | a future-input slot discards arguments the response already knows |
| CS19 | 14 | a terminal symbol result ends with the empty `help[0]` envelope | command surface | a missing envelope fails the AXI approved-result contract |
| CS20 | 32 | a symbol with more references than the cap emits per-package aggregates, `truncated=true`, and one `--full` action | command surface, over-cap fixture | an unbounded default dump floods the agent context |
| CS21 | 13 | `--full` emits every planted row past the cap | command surface, over-cap fixture | a bounded-only surface cannot warrant a universal claim |
| CS22 | 9 | a recomputed sha256 over the bytes before the citation row equals the printed hash value | command surface | a constant hash passes the presence and determinism rows together |
| CS23 | 14 | `--help`, `-h`, and bare `help` each print usage on stdout at exit 0 | usage grammar | a missing help spelling breaks the AXI consistency principle |
| CS24 | 15 | the help text carries the identifies-edges promise clause verbatim | command help test | a silent promise lets resolved rows read as blessed seams |
| BL1 | 16 | `--changed --full` over the fixture pair emits the consumers of each touched declaration | blast core, fixture hunks | a file-level blast lists whole files and misses the symbol edge |
| BL2 | 17 | a declaration the diff deleted emits one `blast_deleted` row and no consumer rows | blast core | a silent drop hides the deletion from the review walk |
| BL3 | 18 | an identical base and tip emit the definitive empty blast table | command surface | an error on the no-op case blocks reviews of doc-only diffs |
| BL4 | 19 | two blast runs over one frozen pair are byte-equal | determinism probe | live-tree input makes review recomputation unmatchable |
| BL5 | 20 | `--source-tip` without `--base` exits 2 with usage | usage grammar | a tip against a defaulted base grades the wrong pair |
| BL6 | 31 | a consumer row inside the diff's file set emits `touched=true` and one outside emits `touched=false` | blast core, fixture pair | an unmarked table hides the FT210 outside-diff signal |
| BL7 | 31 | a blast result offers one per-symbol `--full` help action for each symbol with an untouched consumer | command surface, axi owner | a bare blast table gives no executable next step |
| OI1 | 21 | a `_test.go` function named with a declared helper prefix emits kind `helper` | outline symbol table | the plain `func` kind hides helpers in the undifferentiated list |
| OI2 | 22 | a name with the fake, stub, mock, or spy prefix emits kind `double` | outline symbol table | a fake-only classifier reds the stub, mock, and spy plants |
| OI3 | 23 | a file under a `testdata/` segment emits one fixture row carrying line 1 and its base name | outline walk classifier | extension-only dispatch skips fixture files entirely |
| OI4 | 24 | the old-to-new fixture pair reds on any unreviewed outline byte delta | outline package fixture test | an unpinned migration lets kind values drift silently |
| OI5 | 25 | outline help keeps the LOCATE promise line verbatim | outline help test | a reworded promise upgrades candidate rows to verified ones |
| AD1 | 28 | `bench help` lists `consumers` with its one-line promise | registry help test | an unadvertised verb becomes a dead key for cold sessions |
| AD3 | 29 | `axi-query-registry` accepts the `consumers` disposition | conformance seam | an unregistered AXI surface escapes the ten-principle grade |

Not covered: story 27 — the blast step lands as guidance prose in the review
command doc; the spec-and-tickets review round and the landing gate's prose
checks grade it, and no behavior seam exists for a doc instruction.

### Edge inventory

The walk covers the profile's hostile-input classes that reach this surface:

- A control byte in a git-sourced path drops only its own row and sets
  `truncated=true` in meta, per the outline precedent.
- A path with spaces or glob characters survives whole through the
  NUL-framed git listing.
- A cwd deeper than the repo root resolves through `git.Root`, as outline
  does.
- A dangling or live symlink inside a loaded package stays the loader's
  concern; the command refuses on loader error rather than guessing.
- `GOFLAGS` or tag-altered build contexts sit outside the graded
  configuration; help states the default-context limit.

**Won't handle** lines:

- Reflection, `go:linkname`, plugin, and exec edges — the help states the
  limit, and textual sweeps stay the candidate-class citation for them.
- Non-Go resolution — the refusal names the language; the SCIP consumer is
  the priced future path.
- Base-side enumeration for deleted declarations — a green tip already
  edited every consumer of a deleted symbol, and the `blast_deleted` row
  records the deletion.
- Blast presence enforcement at `bench preflight review` — priced out of
  scope; the review walk stays reviewer judgment.
- Non-default build-tag graphs — the default context is the one graded
  configuration, stated in help.

## Ownership fences

- `internal/consumers/` (new)
- `internal/outline/`
- `internal/conformance/`
- `cmd/bench/`
- `bin/bench.sh`
- `go.mod`
- `go.sum` (new)
- `projects/benchkit.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-cli/SKILL.md`
- `.bench/BENCH-reference.md`
- `specs/resolved-consumer-surface/`

## Out of scope

- SCIP index consumer for non-Go linked repos — 6 edits, 3 gate runs.
- `bench-debug` symptom-map orientation layer — 4 edits, 2 gate runs.
- The A/B benchmark harness from the research asset — 8 edits, 2 gate runs.
- Blast presence check inside `bench preflight review` — 3 edits, 2 gate
  runs.
- Citation-convention edits to the `craft-research` and `craft-review`
  skills — 3 edits, 1 gate run.
- The `craft-seams` reroute wording update to name the resolver — 2 edits,
  1 gate run.
- Base-side blast enumeration across revisions — 3 edits, 1 gate run.

## Further notes

The alias row replaces the research asset's call-graph alias-tail concern.
This design prints no instantiation names, so the tail's mechanism never
reaches the output. The measured ~2 s full pass keeps the
regenerate-per-call posture; no index persists, so no freshness contract
exists. The resolver identifies edges; `projects/<name>.md` keeps seam
blessing, and outline keeps LOCATE. The `testdata/` rule never reaches
`tests/canary/`, because no canary file sits under a `testdata/` segment,
and `bench canary` keeps that inventory.

Build routing (reviewer, 2026-08-28): every ticket charge starts at opus
low. A charge that fails or stalls at that line re-charges once at fable
low. The reviewer pre-approved this top-row escalation for this spec, so it
needs no pause. No other rung exists. This routing supersedes the profile's
doc-authoring effort line for this build.
