# Draft rows that the gate-reuse spec keeps, changes, or retires

Consumed by: ticket #12 of this map, then the slicer of the thin spec.
Drift: a change to a decision of this map.
Retire when: the FT314 spec retires.

The source is the reviewer-signed draft `specs/gate-reuse/spec.md` at commit `9b9150893`.
The draft row IDs stay, so a reader can compare each row with the draft.
"Grade key" replaces the draft's "authorization key" and "oracle identity" in each kept row.

## Kept rows

| row | behavior |
|---|---|
| GR1 | With an in-root launcher and tool, two identity roots with equal content and a space in one path give one grade key. |
| GR2 | Two identity roots with equal content still give two oracle digests, so each checkout's gate cache keeps its key. |
| GR3 | An absent baseline phase manifest and a present empty one give two grade keys. |
| GR4 | Two baseline phase manifests with different bytes give two grade keys. |
| GR5 | A manifest that is a live symbolic link gives two grade keys for two targets with different bytes. |
| GR6 | A dangling link and a FIFO at the manifest give one grade key, and the subject build returns. |
| GR7 | The capability posture changes the grade key. |
| GR9 | The reader plan admits a stored table whose identity a sibling identity root computed over equal content. |
| GR45 | With a `PATH` tool inside the checkout, an ordinary subject and a prospective subject of one tree from that checkout give one grade key. |
| GR10 | After a green prospective run of one tree in a linked worktree, a run of that tree from the primary checkout leaves the run count at 1. |
| GR11 | That primary run prints exactly `gate: green (fresh verdict reused for this tree)`. |
| GR12 | Two trees under one subject name two evidence records. |
| GR14 | After a green under the relaxed posture, a run with `BENCH_REQUIRE_CAPABILITIES=1` runs the gate, and the run count reads 2. |
| GR15 | After a green in a linked worktree, a primary run whose baseline phase manifest differs runs the gate. |
| GR17 | A linked-worktree green whose recorded time is older than the freshness window does not answer the primary run. |
| GR18 | After a linked-worktree green, a red primary run of the same key removes the green, and the next linked-worktree run runs the gate. |
| GR19 | A timed-out run removes the green under its key. |
| GR20 | An evidence record whose key field differs from its file name does not answer, and the run runs the gate. The slicer names the field. |
| GR21 | A delta checkpoint after an expired whole green of its base still runs the delta set. |
| GR22 | A primary run with no gate cache that reuses a linked-worktree green writes the cache with the tree, `green`, and the recorded time of the source run. |
| GR47 | A primary run whose gate cache holds a red verdict of the same tree replaces it with the reused green and the source run's recorded time. |
| GR23 | A reuse that the gate cache already describes leaves the cache bytes unchanged. |
| GR24 | A reuse whose gate-cache write fails returns an operational result with no reusable green. |
| GR44 | A reuse while another run holds the checkout's gate lock, with a cache of another tree, prints the reuse line and leaves the cache bytes unchanged. |
| GR25 | In a fixture with no lane, after `bench worktree merge --from main` ran the whole gate, a spec-less landing of the merge tip leaves the gate tally at one run. |
| GR26 | That landing advances the green marker to its published commit. |
| GR27 | After a green `bench gate --checkpoint specs/x/spec.md --complete` in the source worktree, the spec landing leaves the gate tally at one run. |

## Changed rows

| row | decision | new behavior |
|---|---|---|
| GR8 | #3 | Each of the tool bytes, `PATH`, an environment value, the declared path bytes, the run binary digest, and the running executable digest changes the grade key. |
| GR13 | #3 | A changed run binary digest names another evidence record. A changed running executable digest names another evidence record. |
| GR16 | #1 | After an unpurposed green prospective run of the published tree, the complete checkpoint of that tree leaves the run count at 1. A defect in its completion evidence still refuses it. |
| GR28 | #3 | When the source changes a Bench build input, the spec landing after a complete checkpoint under an equal running executable leaves the tally at one run. After a complete checkpoint under another running executable, the tally reads two. |
| GR29 | #3, #11 | An ADR states the grade key members, reuse across checkouts and verbs, the posture binding, and the broker checkpoint. |
| GR30 | #3 | ADR 16's evidence sentence, ADR 2, the profile, the reference, and the glossary state the grade key as the reuse rule. |

## Retired rows

| rows | decision | reason |
|---|---|---|
| GR31 to GR43, GR46 | #8 | The merge plan disclosure is its own intake. |

## Kept exclusions

- A prune of the evidence store by age.
- Two checkouts that grade one key at the same time. Each run holds its own gate lock, and the store keeps the outcome of the run that finishes last.
- A reuse line that names the source run.
- An edit of the baseline phase manifest while a run is in progress.

Tickets #8 and #3 retire three draft exclusions.
These are a landing plan form, plan guidance, and the run binary identity of an installed kit.
