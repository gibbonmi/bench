# Retro — handoff-sections

## Outcome

`handoff-sections` landed as `dae9a77e` from the reviewed source tip
`fa9ff532` over the frozen base `6fcd0882`. The build ran 23 commits and
11 folds on the integration source, and the landing released it with
`census=5`.

The session handoff is now one git-ignored file with one section per live
assignment, keyed by the request digest, plus `main`. A new leaf package
owns the grammar, the locked read-modify-write, and the removal. `bench
handoff` resolves its section from the assignment record and rewrites only
that section. It keeps a non-empty Next command and refuses a State that
pins a commit outside the tip's ancestry, with one of three reasons.

The retirement path removes the section and prints a failed removal.
`bench status` dates each section by the commits past its recorded tip.
`bench init` ignores the three capture files. The working agreement states
the section rule under an anchor and a fixture.

All 32 acceptance rows resolve to a real test. Five build decisions are
recorded in the spec for veto. No review finding stays open.

## Gate-stage timings

The landing gate, in milliseconds:

| phase | verdict | elapsed |
|---|---|---|
| gofmt | green | 122 |
| vet | green | 913 |
| test | green | 65242 |
| race | green | 2508 |
| system | green | 24765 |
| shellcheck | green | 536 |

Eleven whole-project gates ran, one per fold into the integration source and
one at the landing. Every fold gated green on the first run except the
guidance fold, which went red once on a captured snapshot.

## Ticket-versus-spec-slice and delegate performance

Nine ticket charges ran as nine Opus delegates in sibling worktrees, at most
two at a time under the test-parallel cap. Five repair charges, three
review axes, and two repair-scoped re-reviews followed. Every delegate ran
at Opus low or medium.

Every ticket charge verified its premise with citations. One found the
premise wrong: `liveSpecs` rooted at the caller's checkout, not the primary.
Two delegates reported a fence gap instead of an edit outside it: the
document path constant, and the lock deadline literal. One delegate caught
its own vacuous row when its self-probe stayed green, and added the row.

Seven of nine ticket charges landed their behavior first-pass. The
Next-command charge took one continuation, because its self-probe came back
silently green. The guidance charge took one snapshot resync, because the
anchors table fixture sat outside the fence.

## Coordinator catches

The coordinator ran an independent mutation probe on thirteen of the
fourteen accepted done-claims, at a site and a kind distinct from the
delegate's own. Fifteen probes ran, and thirteen bit. Two were vacuous: a
needle reflow the anchor matcher normalizes, and a swap that did not
compile. A word swap and a constant swap replaced them, and both bit. The
Next-command ticket took no coordinator probe; its continuation's reason
assertion was the row the silent self-probe had shown missing.

The coordinator caught the lock file residue the retirement delegate
reported, and routed it to a repair ticket before the landing. It caught
the legacy document refusal when the new verb first ran over the real
file, and added the migration ticket. It caught the captured anchors
snapshot after the guidance fold went red.

The review round returned 14 findings, which collapsed to six repair
targets. The first repair-scoped re-review found the write-time State check
unreachable, and the second found the repair clean.

## Repair attribution

| ticket | rounds | cause per round |
|---|---|---|
| add-the-handoff-document-leaf-package | 1 | spec-row |
| ignore-the-capture-files-in-linked-repos | 0 | none |
| resolve-the-callers-section-in-bench-handoff | 0 | none |
| remove-the-section-at-retirement | 1 | spec-row |
| own-the-document-path-and-reclaim-the-lock | 0 | none |
| keep-the-next-command-and-refuse-a-stale-state | 1 | delegate-error |
| date-each-section-in-bench-status | 1 | spec-row |
| state-the-section-rule-in-the-working-agreement | 1 | ticket-slicing |
| read-a-legacy-document-as-main | 0 | none |
| restate-the-shape-age-and-derive-the-ignore-path | 0 | none |
| cut-the-unresolved-section-rows | 0 | none |
| refuse-a-state-that-breaks-the-grammar | 1 | spec-row |
| print-the-removal-error-at-retirement | 0 | none |
| drop-the-unreachable-write-time-check | 0 | none |

The leaf ticket's round went to the lock residue the spec's lock decision
did not dispose of. The retirement ticket's round went to the discarded
removal error. The Next-command round went to a tree-hash row the peel could
not distinguish. The status round went to unresolved rows the spec never
asked for. The guidance round went to the captured snapshot its fence did
not name. The grammar round went to a write-time check the spec placed at
an unreachable seam.

## Agent-experience improvements

### Bench CLI

- The landing's census entry records 5 raw calls, and the build's 15
  worktrees held 36. Add `bench test --packages <list>` that owns the
  parallel and count flags, so a focused lane stays inside the CLI.
  Feeds: FT254
- Print the worktree HEAD in `bench worktree path`, so a delegate confirms
  its base in one call.
  Feeds: new
- Scope the follow-on hook's operator scan to a segment that starts with
  `bench`, and name the operator it caught in the refusal.
  Feeds: FT254
- Let `bench gate-prose` take a bare file path with the root implied.
  Feeds: new

### Skills

- Make `craft-spec` require an edge row for a pre-existing artifact the
  new grammar replaces, so a migration is decided before the build.
  Feeds: new
- Make `craft-tickets` name the captured snapshot a new registry row
  moves, in the ticket's `Writes:` line.
  Feeds: FT298

### Process

- Run the new verb over the real artifact at the first phase boundary
  after its ticket folds, because the fixture never holds the legacy shape.
  Feeds: none
- Ask the delegate to prove a write-time check reachable before the review,
  because a parser upstream can make it dead code.
  Feeds: none
