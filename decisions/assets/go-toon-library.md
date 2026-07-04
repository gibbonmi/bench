# TOON library for the Go core — research summary

Research asset for `decisions/go-hooks-port.md` #5. Question: which official TOON
library, how it enters the repo, and what of the hand-rolled emitter survives.

## Recommended library

**`github.com/toon-format/toon-go`** — the Go implementation under the official
`toon-format` org (the org that owns the spec and the reference TypeScript SDK).

| Fact | State |
|---|---|
| License | MIT |
| Releases | none tagged — pseudo-version only (pre-v1, not API-stable) |
| Activity | young but active; official-org home; ~17 importing packages |
| API | `Marshal`/`MarshalString` + options (indent, delimiters, length markers), `Unmarshal`/`Decode` |
| Field order | via struct tags (`toon:"name"`); `map[string]any` accepted but Go maps are unordered — tabular field order requires structs or an ordered path |
| Defaults | comma delimiter, 2-space indent — matches the kit's current shape |

Rejected candidates: `sstraus/toon_go`, `alpkeskin/gotoon`, `rumpl/toon` — all
third-party mirrors of the same spec with no official-org standing; choosing any
of them re-opens the "whose implementation" question the official org answers.

## Inclusion shape

**Normal `go.mod` dependency, pinned to a pseudo-version; no vendoring.**
Consumers never build Go (they receive prebuilt platform binaries), so the
dependency exists only at kit dev/CI build time, where `go.mod` + `go.sum`
already give reproducible, hash-verified builds. Vendoring would add a tracked
copy with no consumer benefit. The pseudo-version pin is explicit and stable;
bumps are deliberate edits. This is the kit's first third-party Go dependency —
the precedent is: official-org, MIT-compatible, build-time-only.

## The compatibility finding (drives ticket #6)

`internal/toon` is **TOON-shaped but not spec-TOON**. Adopting the real library
changes AXI output bytes wherever cells hit the divergences:

| Rule | hand-rolled emitter | TOON spec |
|---|---|---|
| Inner quote escape | doubled (`""`, CSV-style) | backslash (`\"`, JSON-style) |
| Newline in cell | raw newline inside quotes | `\n` escape (control chars must be escaped) |
| Quoting triggers | comma, quote, newline, leading/trailing space | those **plus** empty string, `true`/`false`/`null`, numeric-looking strings, colon, backslash, brackets/braces, leading `-` |
| Cell typing | all cells are strings | numbers/booleans emit bare only when typed; numeric *strings* must be quoted |

Consequences of adoption: every gate contract fragment asserting exact TOON
lines updates to spec bytes; row-building call sites move from `[][]string` to
typed values (structs or typed cells) so numeric columns stay bare instead of
becoming quoted strings; downstream agent readers see spec-compliant TOON
(strictly more parseable by standard tooling). What survives of `internal/toon`:
the AXI error/usage line helpers (`Errorf`, `Usage`, `NotInRepo`) — they are the
hybrid contract, not TOON; the `Escape`/`Table` emitter and the `IsSpace`
whitespace class are superseded by the library.

Fallback flagged, not recommended: keeping the hand-rolled emitter means the
kit's output claims a name it doesn't conform to.
