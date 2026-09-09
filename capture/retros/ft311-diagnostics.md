## Outcome

The FT311 diagnostics spec extends three diagnostic verbs in place. `bench probe` runs a baseline before the mutation, prints a `selection[1]{form,target,run,baseline,ran}` row, and derives its help notes from the focused-run owner. `bench gate-prose` prints the sentence starts of a long paragraph and takes a `--staged` form that grades the index. `bench anchors` prints the line of each needle and refuses an unreadable file.

Landing `2578c9856448c93ccadadd853491066cc85c336f` completed on 2026-09-09 over base `f41917040c24637ca0c5d69441398361559ce319`. The spec status is implemented. The reviewed pair was `f4191704..ca2f008a`.

Seven tickets landed on one retained integration source. Tickets 2 and 3 ran in sibling worktrees and folded through `bench worktree merge`, and each fold ran the whole gate green.

One review round returned 15 findings across three axes and 9 repair targets. Ticket 7 repaired them, and the repair-scoped re-review was clean with four non-blocking follow-ons. Two follow-ons folded into the spec, and two parked as ideas.

The build extended the ticket fences in range eight times. Five test files were over the line budget, and the one-source repairs needed two production files. Each extension is flagged for reviewer veto in the spec's review note.

The landing refused once on an untracked `.codex/config.toml` in the primary checkout. The coordinator preserved a copy, moved the file aside, landed, and restored it byte-exact.

## Gate-stage timings

| run | stage | elapsed | result |
| --- | --- | --- | --- |
| ticket 2 fold at `965711d6` | gofmt | 104 ms | green |
| ticket 2 fold at `965711d6` | go vet | 1,417 ms | green |
| ticket 2 fold at `965711d6` | go test | 77,066 ms | green |
| ticket 2 fold at `965711d6` | race | 6,191 ms | green |
| ticket 2 fold at `965711d6` | system | 27,206 ms | green |
| ticket 2 fold at `965711d6` | shellcheck | 589 ms | green |
| ticket 3 fold at `27d19a87` | gofmt | 110 ms | green |
| ticket 3 fold at `27d19a87` | go vet | 1,201 ms | green |
| ticket 3 fold at `27d19a87` | go test | 90,997 ms | green |
| ticket 3 fold at `27d19a87` | race | 6,205 ms | green |
| ticket 3 fold at `27d19a87` | system | 32,026 ms | green |
| ticket 3 fold at `27d19a87` | shellcheck | 812 ms | green |
| repaired tip at `ad546de4` | gofmt | 126 ms | green |
| repaired tip at `ad546de4` | go vet | 1,026 ms | green |
| repaired tip at `ad546de4` | go test | 79,372 ms | green |
| repaired tip at `ad546de4` | race | 6,196 ms | green |
| repaired tip at `ad546de4` | system | 27,060 ms | green |
| repaired tip at `ad546de4` | shellcheck | 585 ms | green |
| final landing gate | gofmt | 100 ms | green |
| final landing gate | go vet | 960 ms | green |
| final landing gate | go test | 72,474 ms | green |
| final landing gate | race | 2,453 ms | green |
| final landing gate | system | 24,970 ms | green |
| final landing gate | shellcheck | 549 ms | green |

Every gate reported six capability skips and zero environment skips.

## Ticket-versus-spec-slice and delegate performance

| ticket | model / effort | result |
| --- | --- | --- |
| 1. Locate anchor needles | Sonnet / low | Landed after two corrections. The case-fold fixture was not red-capable for a byte offset, and the fixture left the fenced directory. |
| 2. Print the sentence starts | Opus / medium | Landed first-pass. The self-probe expressed an omission as a swap because the omission did not compile. |
| 3. Grade the staged Markdown | Opus / medium | Stopped at the fence when eleven tests did not fit one file, then landed after a fence extension and a mechanical move. |
| 4. Name the probe selection | Opus / medium | Landed first-pass with one minimal out-of-fence edit that the preflight named. |
| 5. Run the baseline | Opus / medium | Landed after one round for two refusal rows and a test split. |
| 6. Fold the guidance | Opus / high | Landed first-pass with every pinned sentence byte-identical. |
| 7. Repair the review findings | Opus / medium | Landed first-pass with two more fence needs the preflight would have named. |

Every charge was ticket-sized. No charge took a spec slice. The three review axes and the re-review ran as Opus / medium read-only delegates.

## Coordinator catches

- Ticket 1's case-fold fixture placed the KELVIN SIGN inside the match, so a byte-offset locate stayed silent. The repaired fixture put the fold on the line before the match.
- Ticket 1 deleted the fenced fixture directory, and the next charge failed on writes-resolve. The fixture moved back under the fenced path.
- Ticket 3 stopped at its fence instead of a write outside it. The coordinator extended the fence in range and flagged it.
- Ticket 5's two refusal rows asserted no run child, which the baseline made false by spec. The rows now assert that only the baseline started.
- Every ticket took an independent probe at a distinct site and kind, and every probe bit.
- The Coverage axis found that the help's own `bench gate-prose . --staged` example refused on every relative root. A hand run confirmed it.
- The Spec axis found a constant `baseline` cell, the defect the spec's own round-one review had blocked. The repair threads the observed kind, and the row's clause now names its structural half.
- The Standards axis found two restated prose diagnostics and an independent test expectation without a recorded red. The coordinator ran the two probes and recorded the reds in the spec.

## Repair attribution

| ticket | repair rounds | cause per round |
| --- | ---: | --- |
| 1.md | 2 | delegate-error, delegate-error |
| 2.md | 0 | none |
| 3.md | 1 | ticket-slicing |
| 4.md | 0 | none |
| 5.md | 1 | ticket-slicing |
| 6.md | 0 | none |
| 7.md | 0 | none |

The review round's nine targets attribute as follows. The relative root and the blob bounds are `spec-row`. The single-pass strip, the constant cell, and the five one-source items are `delegate-error`.

## Agent-experience improvements

### Bench CLI

- Add an `--unwrap <call>` mutation and a cause cell on an `invalid` verdict to `bench probe`, and add a package headroom view to `bench structure --growth`. Make `toon.TableTyped` return an error on a short row. The landing census is 25 raw calls, captured as `ft311-diagnostics landing census: 25 raw calls`.
  Feeds: new
- Make `bench preflight build` accept a `(new)` test file beside an over-budget sibling as a proposed write. Then a fence extension for a test split needs no spec amendment.
  Feeds: new

### Skills

- State in `craft-tickets` that a ticket with more than six command-level tests names a second test file in its `Writes:`. This applies when the sibling is near the line budget.
  Feeds: new
- State in `craft-delegate` that a charge names each fenced path that must still exist on return, so a fixture cannot leave the tree.
  Feeds: new

### Process

- Run `bench preflight build` after every ticket commit on the integration source and before the next charge, because a delegate can remove a fenced path.
  Feeds: new
- Put an untracked file in the primary checkout aside as its own preserved step before a landing, and restore it with `cmp` after. The landing refuses a dirty destination.
  Feeds: none
