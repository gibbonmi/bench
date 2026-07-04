# TOON library adoption — spec-conformant AXI output

Status: implemented

Map: `decisions/go-hooks-port.md` #4/#5/#6 (TOON spec, sequenced first, ahead of
the hooks port). Research asset: `decisions/assets/go-toon-library.md`.

## Problem

The Go core hand-rolls a TOON-*shaped* emitter (`internal/toon`) that is not
spec-TOON: it escapes inner quotes by doubling (`""`, CSV-style) instead of
backslash (`\"`), keeps raw newlines inside cells instead of `\n` escapes, and
quotes on a far narrower trigger set than the spec. Every AXI command's stdout
therefore claims the TOON name without conforming, and a standard TOON decoder
mis-reads the divergent cells. The standing prefer-official-libraries rule and
map #4/#6 direct the swap.

## Solution

Replace the hand-rolled emitter with `github.com/toon-format/toon-go` — the Go
implementation under the official spec-owning org (MIT, pinned pseudo-version,
build-time-only). AXI stdout moves to spec-TOON bytes. A thin `internal/toon`
adapter over the library preserves the kit's block *contract* — the `name[N]{fields}:`
count header, the `{fields}` schema even when the table is empty, and the
trailing newline — none of which the library gives for free. Numeric columns move
to typed cells so they stay bare instead of becoming quoted strings.

## User stories

1. As a downstream agent or tool consuming AXI TOON, I want spec-conformant cell
   escaping — backslash inner quotes, `\n` for control characters, and the spec's
   full quoting-trigger set — so a standard TOON decoder parses the output.
   Line: claude-opus-4-8 / medium. The escaping and special-value surface is the
   semantic core of the swap, and although the gate observes the exact bytes, the
   adapter shape that produces them carries the design judgment.

2. As that same consumer, I want a genuinely-numeric column (the maps `ticket`)
   emitted bare, not as a quoted string, so integer semantics survive a
   round-trip through a spec decoder.
   Line: claude-opus-4-8 / low. The change is small but it defines the typed-cell
   entry point and the mixed-column handling where `ticket` is an int for real
   tickets and the literal `handoff` for close-readiness rows.

3. As a consumer of an empty AXI table, I want the schema header
   `name[0]{fields}:` preserved, because the library drops `{fields}` for an empty
   array, so an empty result stays self-describing.
   Line: claude-opus-4-8 / low. It is a compatibility shim the library does not
   provide, and it is fully gate-observed by four existing empty-table contracts.

4. As a kit maintainer, I want every row call site migrated to the adapter and
   every affected contract fragment updated to spec bytes in the same diff, so the
   regression net matches the new output and no fragment asserts a superseded byte.
   Line: claude-sonnet-5 / low. The work is mechanical and fully gate-observed.

5. As a kit maintainer, I want the first third-party Go dependency added
   correctly — a pinned `go.mod` require, a generated `go.sum`, pure-Go so the
   cross-compile matrix stays green, build-time-only — so the precedent for kit
   Go deps is set cleanly.
   Line: claude-sonnet-5 / low. Adding and pinning the module is mechanical and
   the gate builds every platform in the matrix.

## Implementation decisions

**Dependency.** `github.com/toon-format/toon-go` at
`v0.0.0-20251202084852-7ca0e27c4e8c` (pseudo-version; the module is pre-v1 with no
tagged release, so bumps are deliberate edits). A normal `go.mod` `require` plus a
generated `go.sum`; no vendoring. The library is pure-Go (stdlib only:
`strings`/`strconv`/`unicode`/`time`), so the `scripts/platforms.json`
cross-compile matrix stays clean and cgo-free. Consumers receive prebuilt
binaries, so the dependency is kit dev/CI build-time-only. Precedent for kit Go
deps: official-org, MIT-compatible, build-time-only.

**`internal/toon` becomes a thin adapter, not an emitter.**

- API: `Table(name string, fields []string, rows [][]string)` keeps its signature
  for the four all-string callers (learnings, diff, guards, coverage) and
  delegates to a shared core; a new `TableTyped(name string, fields []string, rows [][]any)`
  serves maps' mixed `ticket` column. One emission core behind two typed entry
  points, the way `Print`/`Printf` share one formatter. Recommended over a single
  `[][]any` entry, which would ripple four `[][]string` helper signatures
  (`coverage.Rows`, `diff.changedFiles`, and their tests) for no behavior gain.
- Core: marshal `map[string]any{name: []toon.Object}` with pinned options —
  `WithIndent(2)`, comma delimiter (default), and length markers **off** (the
  default; length markers on would emit `name[#N]{...}`, off gives the kit's
  `name[N]{...}`). Each row is a `toon.NewObject(Field{Key,Value}...)` built in
  `fields` order, which is how the library produces ordered tabular output despite
  Go maps being unordered.
- **Empty-table shim.** For an empty array the library emits `name[0]:` — it
  infers `{fields}` from the elements, of which there are none, so the schema
  header is lost. The adapter hand-renders `name[0]{fields}:\n` when `rows` is
  empty (an empty table has no cells, so no escaping is at stake) to preserve the
  kit's self-describing empty contract. This is a refinement the research asset
  (#5) did not surface; it is what keeps #6's "block shape unchanged" promise true.
- **Trailing newline.** The library omits the final `\n`; the adapter appends one.
  Callers concatenate blocks and the runtime contracts expect a newline per line,
  header included. Also unprobed by #5.
- `Escape` (package-internal, the doubled-quote rule) is removed — superseded by
  the library. `IsSpace` **survives**: it has a live non-emitter consumer
  (`learnings.isSpace`, which trims fields), so it stays the shared AXI whitespace
  class. This corrects the handoff's item-1 claim that `IsSpace` is superseded —
  only the emitter's *use* of it is.
- `Errorf`, `Usage`, `NotInRepo` survive unchanged — they are the hybrid AXI error
  contract, not TOON.

**Typed cells.** `maps.fileRows` returns `[][]any`; the `ticket` cell is
`strconv.Atoi(t.num)` for a real ticket (all-digit by `ticketRe`) and the string
`"handoff"` for a close-readiness row. Per-cell typing, not a typed struct — the
column is genuinely mixed, which is why the map named "typed cells" rather than
"structs". Every other column stays a string; `coverage.story` correctly gains
quotes when it is a numeric-looking value like `3`, because it is author text that
can also be `1–3` or `edge of N`, not an integer.

**Package doc** updates from "only the flat-table form is implemented, the general
format intentionally absent" to the adapter role over the library.

## Testing decisions

- A good test here exercises the adapter at its pure-function seam (bytes in →
  bytes out) and the commands at their AXI-stdout seam (contract fragments), not
  internals. Prior art: the existing `internal/toon/toon_test.go` and the
  `.bench/gate-axi*-contracts.sh` fragments.
- Two seams get tested: (1) `internal/toon`'s `go test`, rewritten to spec-TOON
  expectations; (2) the AXI contract fragments, the integration regression net.
- Gate command: the project gate, `bench gate` — `go test ./...`, the shell
  contract fragments, and the build/cross-compile matrix.

### Seam diagram

Seam 1 — the adapter (unit):

    trigger: any AXI command building a table
        │
        ▼
    name, fields, rows(typed) ──▶ [ toon.Table / TableTyped ] ──▶ "name[N]{fields}:\n  <spec-escaped rows>\n"
                                     wraps toon-go Marshal;
                                     empty-shim; trailing \n
                                        ◀ tests attach: toon_test.go drives inputs, asserts exact bytes

Seam 2 — command stdout (integration):

    trigger: bench <learnings|diff|guards|maps|coverage>
        │
        ▼
    repo state ──▶ [ command → toon.Table/TableTyped ] ──▶ stdout TOON block
                        ◀ tests attach: gate-axi*-contracts.sh assert exact lines

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | inner quote escapes backslash, not doubled | toon_test + gate-axi-wave2:60 | rewrite wave2:60 assertion to `  A,"a\"q.txt"`; fails against the current doubled-quote emitter | pins the exact escaped byte; the doubled-quote output `"a""q.txt"` mismatches |
| 1 | comma+quote title → `"a, \"b\""` | gate-axi-contracts:181 | rewrite :181 to the backslash form; fails vs current `"a, ""b"""` | pins the escaped title bytes end to end |
| 1 | each special-value class is quoted (enumerated below) | toon_test, one row per class | add the per-class rows; red against the current narrow `needsQuote` (does not quote `true`/empty/numeric-looking/colon/leading-`-`) | one row per trigger class fails until the library's spec quoting is in force |
| 2 | maps `ticket` digit emits bare `6`; `handoff` emits bare | new gate-axi-contracts maps-row assertion `  <file>,6,Grill,open` + toon_test int-vs-string row | add the row assertion; red until `ticket` is typed `int` (string `"6"` quotes to `"6"` under spec) | a bare number only survives if the cell is a Go `int`, proving the typed-cell path |
| 3 | empty table keeps `{fields}` | gate-axi-wave2:71,133 and gate-axi-contracts:121,126,134 (five existing) | a naive library delegate emits `name[0]:`; all five existing assertions fail | the assertions already pin `{fields}`; the shim is exactly what keeps them green |
| 4 | all five call sites emit via the adapter; contracts green | full gate | pre-migration `go build` fails on the removed `Escape` and the changed row types | the gate is the oracle for the migration's completeness |

Special-value classes enumerated for the story-1 row set (one `toon_test` row
each): empty string, `true`, `false`, `null`, numeric-looking (`3`, `-1`, `1.5`),
colon, backslash, brackets/braces, leading `-`; plus the carried-over cases —
leading/trailing space, comma, inner quote (now backslash), newline (now `\n`).

### Edge inventory

Edge classes walked per behavior, each resolved as a coverage row above or a
Won't-handle line here:

- error path → Won't handle (below).
- empty / absent input → story 3 (five empty-table assertions).
- boundary values → single-row and single-column tables; add a `toon_test` row
  (learnings single row) and covered by existing single-count headers.
- malformed input → the every-trigger cells of story 1.
- interrupted / partial state → N/A; the adapter is a pure function.
- re-run idempotency → `Object` preserves field order, so output is deterministic;
  covered implicitly by exact-byte assertions.
- hostile environment → unicode cells (em-dash, existing) and CRLF-bearing source
  (existing wave2:123) pass through unchanged; the emitter is byte-transparent.

**Won't handle:**

- Field names that need quoting — every AXI field name is a bare compile-time
  identifier, never data; the empty-shim joins them verbatim while the non-empty
  path routes through the library's `encodeKey`, and for the actual field sets the
  two are identical. Safe: fields are constants.
- A cell whose *bytes* spec-TOON cannot represent — a control character other than
  the escapable tab/newline/return, which can ride in a git path through `bench
  diff` — is not emittable at all. This is reachable (a string trivially carries
  such a byte; the Go *type* being `string | int` does not bound its content), so
  the adapter does not forge output: the library refuses and `Table`/`TableTyped`
  return an error, which every command surfaces as the AXI error contract
  (`error: unrepresentable TOON cell — …`, exit 1) rather than crash or emit a
  lossy block. `TestTableUnrepresentableCellErrors` and the per-command error path
  pin it.
- General (nested) TOON documents — the kit emits only flat tables; the library
  supports more, but the adapter exposes only `Table`/`TableTyped`. Safe: outside
  the kit's output surface.

## Out of scope

- **The hooks port** (`stop.sh`, `check-agent-line`, `_line-guard`, `lines-env`,
  the shared `resolve-bench.sh`) — the sibling spec sequenced *after* this one
  (map #1–#3 and the handoff). A distinct capability with its own spec; nothing is
  shared with the emitter swap. Estimate to build later: ~15 edits, ~10 gate runs.
- **Decoding / `Unmarshal` anywhere in the kit** — the kit only emits TOON; no Go
  consumer parses it. A separate capability if a reader is ever needed. Estimate:
  ~4 edits, ~3 gate runs.
