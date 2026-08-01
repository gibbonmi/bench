# Requirement identifiers for this build

`bench coverage` emits `{story, seam, red_signal}` with no stable row identity, and
this spec has rows sharing a story and a seam string, so `[R<n>]` cannot be derived
mechanically from its output. The identifiers below are assigned by hand, in the
order `bench coverage specs/reduced-gate-phase-set/spec.md` enumerates its 26 rows.
Each one names the (story, behavior) pair from the spec's own acceptance coverage
map. This is the map per-row accounting reads.

| id | story | behavior (spec coverage map) | ticket |
|---|---|---|---|
| R01 | 1 | Every consumer resolves the allowlist from the declaration rather than a private copy | Route the ambient staleness signal through the scope declaration |
| R02 | 1 | Every co-located capture surface is a member through the directory, not through its own entry | Declare the reduced gate scope |
| R03 | 1 | A path outside the declaration is never confined, including a near-miss sibling | Declare the reduced gate scope |
| R04 | 1 | A file directly under a declared directory is a member; a nested or sibling-prefixed path is not | Declare the reduced gate scope |
| R05 | 2 | The profile's pinned allowlist prose matches the declaration | Bind the scope declaration to the profile prose |
| R06 | 3 | The stripped identity is unchanged by an edit confined to allowlisted paths | Compute the stripped subject identity |
| R07 | 3 | The stripped identity moves for any edit outside the allowlist | Compute the stripped subject identity |
| R08 | 4 | An excludable phase that hard-reads an allowlisted path fails a full gate | Run excludable phases against a stripped worktree |
| R09 | 4 | An excludable phase that soft-skips on the missing path fails a full gate | Run excludable phases against a stripped worktree |
| R10 | 4 | Included phases still see the allowlisted paths on the same run | Run excludable phases against a stripped worktree |
| R11 | 4 | The materialized stripped tree is a git repository | Run excludable phases against a stripped worktree |
| R12 | 4 | No contract test reads an allowlisted path relative to the kit checkout | Confine contract-test capture reads to the subject root |
| R13 | 4 | The excludable set is not degenerate: `contract` is in it | Declare the reduced gate scope |
| R14 | 5 | A capture-only changeset runs only included phases | Reduce the gate run and inherit from a full-green ancestor |
| R15 | 5 | The verdict records the reduced marker, the phases run, and the ancestor | Expand the verdict record with the reduced shape |
| R16 | 5 | The loader rejects a record mixing the full and reduced shapes | Expand the verdict record with the reduced shape |
| R17 | 6 | A second consecutive capture commit inherits from the same full green | Reduce the gate run and inherit from a full-green ancestor |
| R18 | 6 | A reduced verdict is never itself a valid ancestor | Reduce the gate run and inherit from a full-green ancestor |
| R19 | 6 | An ancestor older than the freshness window forces a full run | Reduce the gate run and inherit from a full-green ancestor |
| R20 | 6 | An allowlist-confined change with no ancestor runs the full gate | Reduce the gate run and inherit from a full-green ancestor |
| R21 | 7 | `bench prep-release` refuses a reduced verdict and names that reason | Refuse a reduced verdict for release evidence |
| R22 | 8 | `bench commit` takes the reduced path for an allowlist-confined staged set | Take the reduced path from bench commit |
| R23 | 8 | A staged set mixing an allowlisted and an unlisted path runs the full gate | Take the reduced path from bench commit |
| R24 | 8 | The reduced run announces the skipped phases and the ancestor | Reduce the gate run and inherit from a full-green ancestor |
| R25 | 9 | The declaration, the phase set, and the prose cannot drift | Bind the scope declaration to the profile prose |
| R26 | 9 | The check bites in both directions | Bind the scope declaration to the profile prose |

R24 sits with the gate-run ticket rather than the commit ticket because the
announcement is emitted by the reduced execution path itself, which that ticket
owns; the commit ticket owns only route selection from the staged set.

## Blocker graph

    Declare the reduced gate scope
      ├── Route the ambient staleness signal through the scope declaration
      ├── Compute the stripped subject identity ──┐
      ├── Run excludable phases against a stripped worktree ──┐
      ├── Confine contract-test capture reads to the subject root
      └── Expand the verdict record with the reduced shape ──┤
                                                             ▼
                          Reduce the gate run and inherit from a full-green ancestor
                                        │                     │
                                        │                     └── Refuse a reduced verdict
                                        │                          for release evidence
                                        └── Take the reduced path from bench commit

    Bind the scope declaration to the profile prose  ← blocked by every ticket above

Verified by hand: every `Blocked by:` title names an existing sibling file's `#`
heading byte-for-byte, and the graph is acyclic.
