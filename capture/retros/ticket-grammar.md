# Retro — ticket-grammar (FT174)

## Outcome

The spec landed as commit `9f341d53` on a full green gate. One parser
(`internal/tickets`) owns the ticket-file schema. Preflight renders six grammar
rows from parsed tickets, the `ticket-grammar` conformance check sweeps every
staged spec, and a ten-fixture canary family proves the diagnostics bite. The
review round collapsed 23 raw findings into 4 repair targets; 3 were repaired,
and 1 closed as a recorded no-op. A one-line light-path landing (`40809f4a`)
migrated the citation spec's ticket to the new grammar.

## Gate-stage timings

The landing gate ran gofmt 104 ms, vet 936 ms, test 56.1 s, race 2.6 s,
system 20.7 s, and shellcheck 537 ms. The mid-build merge gate ran test 60.1 s
and system 21.9 s on the same shape.

## Ticket-versus-spec-slice and delegate performance

Seven ticket charges ran on opus at medium effort, one migration charge on
opus at low effort. Six of seven ticket charges landed first-pass on behavior.
The preflight ticket took one repair round for composed-tree reds its focused
checks could not see. Every charge ran its own biting mutation probe. Three
charges reported honest blockers instead of weakened rules: the live closure
reds, the impossible canary classes, and the compiled-binding fixture limit.
Ticket-sized charges with exact seams outperformed the wide repair charge,
which needed one scoped fix round.

## Coordinator catches

- One probe matched no source text and returned a false green; the diff check
  caught the no-op sed and the redone probe bit.
- The merge gate caught systemtest fixture reds that focused package checks
  could not see.
- The stale `dist/bench` rendered eight preflight rows and hid the closure
  rows; a `go run` rerun surfaced the real verdicts.
- The landing caught a staged spec on `main` that the new sweep reds; the
  no-migration assumption in the spec was stale by one concurrent landing.

## Repair attribution

| ticket | repair rounds | causes |
| --- | --- | --- |
| lift-shared-field-grammar | 0 | none |
| build-ticket-file-parser | 0 | none |
| read-preflight-rows-from-the-parser | 1 | ticket-slicing |
| close-the-ownership-closures | 0 | none |
| register-ticket-grammar-sweep | 0 | none |
| prove-ticket-grammar-canary-bites | 0 | none |
| advertise-the-enforced-ticket-grammar | 0 | none |
| review repair R1–R4 | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- The landing census recorded 288 raw calls, mostly worktree file reads and
  writes with no Bench form. The census learning entry proposes a worktree
  read projection.
  Feeds: new
- `bench preflight` runs through the installed `dist/bench`, so new verdict
  rows stay invisible until a landing. A freshness note in the preflight
  output would name the stale binary.
  Feeds: new

### Skills

- The craft-spec canary guidance should require each named canary class to be
  a diagnostic its check emits on an addable or mutable tree state. Two named
  classes were unbuildable.
  Feeds: new

### Process

- A grammar-enforcement spec should name every in-repo fixture its rules
  grade, because two builds paid fence deviations for shared fixture files.
  Feeds: new
- A spec that assumes an empty class of tree states should carry a landing-time
  recheck, because a concurrent landing invalidated the no-migration note.
  Feeds: none
