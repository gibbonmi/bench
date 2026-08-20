# Worktree exec run binary

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-20. Three forks closed by the reviewer: the contradiction resolves in code rather than in prose; the ticket set is rewritten whole to cover that resolution, the gate-entry reword, and an explicit FT223 exclusion; and, after the first review round found that pointing the child at the worktree's own `dist/bench` would make `bench test` and `bench commit` reuse a possibly-stale artifact against the profile's private-exact-source rule, the child is marked wrapper-rooted so its gate owns a private exact-source build instead.

Verification log: 2 iteration(s) to accept — one `fable`/high round returned three blocking findings. The first reopened the resolution fork: re-pointing the run-binary variable makes the child inherit, and inherit verifies the binary's seal rather than its source, so an exec child's `bench test` and `bench commit` could reuse a stale worktree artifact against `projects/benchkit.md`'s private-exact-source rule. The reviewer closed the reopened fork on the wrapper marker, which dissolved that finding and the two that followed it — the empty-file edge and the fence-crossing acceptance clause both belonged to the discarded discovery predicate. Folded partials: the one-source claim corrected to name the Go-side validator as a pre-existing second derivation rather than asserting a single source, the composed-run row relabelled as ticket evidence instead of retained gate coverage, the exclusions group's fictional line dropped, the story-group effort corrected to the cached mid routing, and the gate entry's refusal count corrected from three messages to two.

## Problem

Two documents each state a correct rule, and a session that obeys both writes a
command that always fails.

`craft-delegate` addresses an assignment worktree as
`bench worktree exec "<label>" -- <command>`, never a cached path. The project
profile says to exercise a durable worktree artifact by invoking that worktree's
own `./dist/bench`, because `bench` on PATH resolves to the main checkout's
wrapper and may belong to a different source tree. Composed, they give
`bench worktree exec <target> -- ./dist/bench gate`, and that refuses every time.

`execEnv` strips `BENCH_KIT`, `BENCH_WRAPPER`, and `BENCH_RUN_BINARY` from the
child, correctly: the caller's selection was built for the caller's run, and a
leaked kit names the caller's checkout. The gate's owner lookup then reads the
absence of both routing variables as "no one owns this run" and returns nothing,
so the gate never selects a binary, its plan's environment never gains one, and
the gate entry refuses. The child is left unable to author the selection the
caller was right not to lend it.

The refusal then names only the variable. No operator sets it, nothing in the
message says what to run instead, and neither condition that produces it is
named. A session that reads it concludes the verb it just used is broken. That
is observed: one session abandoned `bench worktree exec` after meeting this
refusal and did not return to it.

## Solution

Mark the child as wrapper-rooted. After `bench worktree exec` strips the
caller's routing, it points `BENCH_WRAPPER` at the child worktree's own wrapper
when one is present there. The gate's owner lookup then owns the run: it builds
one private binary from the child's own source and removes it after teardown,
which is what the profile requires of every gate and `bench test` run.

The alternative — pointing the run-binary variable at the worktree's own
`dist/bench` — was considered and rejected. It makes the child inherit rather
than own, and the inherited path verifies the binary's own seal instead of
checking it against its source. An exec child's `bench test` or `bench commit`
could then reuse an artifact that had fallen behind the tree it grades.

Nothing executes the value exec supplies. It is a marker the owner lookup reads
and a path the adoption doctor reports, so exec vouches for no executable and
the binary that eventually runs is one the child built from the tree under
grade.

Separately, the gate entry's refusal gains the next action it lacks. The two
changes are independent: the first removes the common way to reach the refusal,
the second repairs the refusal for every other way.

## User stories

### The exec child owns its run

Line: opus / mid. Gate and conformance logic at an existing seam with an
existing harness, which is the cached routing for this area.

1. As a session working in an assignment worktree, I want `bench worktree exec <target> -- ./dist/bench gate` to run the gate, so that I can grade a worktree without leaving the sanctioned verb.
2. As a session, I want the exec child to receive `BENCH_WRAPPER` naming that worktree's own wrapper, so that the gate's owner lookup treats the run as owned.
3. As a session, I want the emitted value to be an absolute path, so that the adoption doctor reports a path it can resolve.
4. As a session, I want the child's gate to build a private binary from the child's own source, so that no pre-existing artifact grades the tree.
5. As a session, I want the exec child never to receive the caller's `BENCH_WRAPPER` value, so that the child is never attributed to the caller's wrapper.
6. As a session, I want the exec child never to receive the caller's `BENCH_KIT` or `BENCH_RUN_BINARY`, so that the child resolves its own kit and never runs a binary built for another tree.
7. As a session, I want every unrelated variable I set to reach the exec child, so that exec stays a transparent runner.
8. As a session in a worktree carrying no wrapper, I want the child to receive no `BENCH_WRAPPER`, so that an absent wrapper does not become a false marker.
9. As a session in a worktree whose wrapper path is a directory, I want the child to receive no `BENCH_WRAPPER`, so that a non-file at that path does not become a marker.
10. As a session in a worktree whose wrapper path is a FIFO, device, or socket, I want the child to receive no `BENCH_WRAPPER`, so that a special file never becomes a marker.
11. As a session in a worktree whose wrapper path is a dangling symbolic link, I want the child to receive no `BENCH_WRAPPER`, so that a broken link is not read as an authoritative wrapper.
12. As a session in a worktree whose wrapper path is a present-but-empty regular file, I want the child to receive `BENCH_WRAPPER`, so that the marker follows one stated predicate rather than a second judgment about content.
13. As a session running `bench test` or `bench commit` through exec, I want the child unchanged, so that this alters no verb that already owned its binary.
14. As a session, I want the child's environment to differ from today's by the wrapper assignment alone, so that no incidental change rides along.

### The refusal names its next action

Line: opus / mid. One message and its contract test, in the same cached routing.

15. As a session that reached the gate entry without a wrapper, I want the refusal to name the wrapper invocation, so that I can act on it without reading the gate script.
16. As a session, I want a relative `BENCH_RUN_BINARY` to meet that same refusal, so that one message covers both unusable forms.
17. As a session outside any worktree, I want the refusal to name no worktree, so that the message serves every wrapper-less invocation rather than one of them.
18. As a reviewer, I want the refusal to keep naming `BENCH_RUN_BINARY`, so that the existing gate-entry contract test keeps binding to it.
19. As a reviewer, I want a valid absolute executable to be accepted exactly as it is today, so that the reword changes no verdict.
20. As a reviewer, I want the two later refusals in that block unchanged, so that conditions an operator can already act on keep their wording.

### Reviewed exclusions

21. As a reviewer, I want FT223's inherited-verdict refusal left to FT223, so that this spec does not pre-empt a decision I own.
22. As a reviewer, I want the portable tilde form of a worktree path left alone, so that a deliberate portability affordance is not removed as a defect.
23. As a reviewer, I want no new prose advertising the pairing, so that the fix removes the contradiction instead of documenting around it.

## Implementation decisions

**A marker, not a selection.** Exec supplies the one fact the owner lookup is
missing — this run was rooted by Bench, not typed into a shell — and supplies it
in the variable that already carries that meaning. Exec names no executable and
authenticates none.

**Discovery is one predicate.** Exec sets the marker when a regular file sits at
the worktree's own wrapper path, and leaves it unset otherwise. Content is not
inspected: an empty wrapper still sets the marker, because the marker's only
readers are an owner lookup that tests for non-emptiness and a doctor that
reports a path.

**The worktree path is already validated.** The assignment ledger refuses a
stored worktree path that is not absolute and cleaned, and worktree creation
canonicalizes it through symbolic links before storing it, so the wrapper path
inherits both properties rather than re-deriving them.

**Exec adds no new derivation of the trust predicate.** The path predicates that
decide whether a binary may be launched are spelled twice in the tree already —
once in the run-binary validator and once in the gate entry — and this spec adds
no third. Collapsing the existing pair is its own change against its own risk.

**The reword is bounded by an existing contract.** The gate-entry contract test
asserts the refusal names the run-binary variable, so the reworded message keeps
that token. The gate script never ships to consumers, so kit-relative wording is
safe.

### Bootstrap authority

Each hop names how the validator authenticates the next executable before
launching it.

- The operator selects the wrapper on PATH. That choice is the trust root; Bench
  does not authenticate it.
- `bench worktree exec` resolves the target through the assignment ledger, which
  refuses an inactive assignment and a worktree path that is not absolute and
  cleaned. Exec then launches the argv the operator named, authenticating
  nothing about it, because the operator authored it.
- The child's gate reads the marker and owns its run: it builds a private binary
  from the child's own source through the build script. The executable it will
  launch is one it just produced from the tree under grade, so no pre-existing
  artifact is trusted at this hop.
- The run-binary validator then checks that built binary — absolute and cleaned,
  regular and executable, no symbolic-link traversal — and, because the
  selection is owned rather than inherited, verifies it against its source
  rather than against its own seal.
- The gate entry re-checks the same path predicates and runs a freshness check
  against the kit before it hands over the run. It refuses rather than launching
  on any failure.

## Testing decisions

- A good test exercises the environment a real child actually receives, and the
  gate entry's accept-or-refuse verdict for a supplied value — not the internal
  shape of either.
- Seams and prior art: `runWorktreeChild` already has a harness that runs a real
  child and returns its environment one assignment per line; the gate entry
  already has a conformance harness that runs the script as a real subprocess
  with a chosen environment.
- Gate seam: the conformance package observes the gate entry, and the worktree
  package observes the exec child. Both run in the ordinary gate.

### Seam diagram

    trigger: bench worktree exec <target> -- <argv>
        │
        ▼
    <target>  ──▶  [ resolveWorktree ]  ──▶  absolute cleaned worktree path
        │                                          │
        ▼                                          ▼
    caller env  ──▶  [ execEnv ]  ──▶  child env without caller routing,
                          │            with the worktree's wrapper marker
                          │            ◀ tests attach here: run a real child,
                          │              read the assignments it received
                          ▼
                    [ child argv ]  ──▶  ./dist/bench gate
                                              │
                                              ▼
                                    [ gate owner lookup ]  ──▶ owns a private
                                              │                exact-source build
                                              ▼
                                    [ .bench/gate.sh entry ]
                                        ◀ tests attach here: run the script as a
                                          real subprocess, assert accept or refuse

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WX1 | 2 | a child in a worktree holding a regular-file wrapper receives the wrapper variable naming that worktree's wrapper path | exec child environment | a build that keeps stripping the marker leaves the assignment absent |
| WX2 | 3 | the emitted value is an absolute path | exec child environment | a build joining a relative wrapper path emits a value the doctor cannot resolve |
| WX3 | 4 | a child's gate selects an owned binary rather than an inherited one | exec child environment | a build setting the run-binary variable instead makes the child reuse a stale artifact |
| WX4 | 5 | a child never receives the caller's wrapper value | exec child environment | a build that forwards instead of re-pointing attributes the child to the caller |
| WX5 | 6 | a child receives neither the caller's kit nor the caller's run-binary variable | exec child environment | a build that stops stripping routing makes the child resolve the caller's checkout |
| WX6 | 7 | a child receives every unrelated variable the caller set | exec child environment | a build that rebuilds the environment from scratch silently drops operator state |
| WX7 | 8 | a child in a worktree with no wrapper path receives no wrapper variable | exec child environment | a build that sets the marker unconditionally names a path that does not exist |
| WX8 | 9 | a child in a worktree whose wrapper path is a directory receives no wrapper variable | exec child environment | a build testing only existence accepts a directory as the wrapper |
| WX9 | 10 | a child in a worktree whose wrapper path is a FIFO receives no wrapper variable | exec child environment | a build testing only existence names a special file |
| WX10 | 11 | a child in a worktree whose wrapper path is a dangling symbolic link receives no wrapper variable | exec child environment | a build reading the link without resolving it treats a broken link as a wrapper |
| WX11 | 12 | a child in a worktree whose wrapper path is an empty regular file receives the wrapper variable | exec child environment | a build adding a content or size predicate departs from the one stated rule |
| WX12 | 13 | a child running a verb that already owned its binary is unaffected | exec child environment | a build that changes the owned path alters verbs this spec does not target |
| WX13 | 14 | a child's environment differs from today's by the wrapper assignment alone | exec child environment | a build with an incidental environment change alters unrelated behavior |
| WX14 | 15 | the gate entry with no run-binary variable refuses and names the wrapper invocation | gate entry subprocess | a build restating only the variable leaves the message unactionable |
| WX15 | 16 | the gate entry with a relative run-binary value refuses with that same message | gate entry subprocess | a build wording only the absent case leaves the relative case unactionable |
| WX16 | 17 | the refusal names no worktree | gate entry subprocess | a build wording the message for exec alone misleads every other wrapper-less caller |
| WX17 | 18 | the refusal text contains the run-binary variable name | gate entry subprocess | a build dropping the token breaks the standing gate-entry contract |
| WX18 | 19 | the gate entry with a valid absolute executable proceeds as it does today | gate entry subprocess | a build whose reword changes the accept path alters a verdict |
| WX19 | 20 | the regular-executable and physical-path refusals keep their present wording | gate entry subprocess | a build rewording the whole block changes messages an operator can already act on |
| WX20 | 1 | `bench worktree exec <target> -- ./dist/bench gate` over a clean worktree of this kit reports green | ticket evidence, one composed run | a build correct at each seam separately can still fail composed |

Not covered: story 21 — FT223 is a separate decision-required item with its own competing fixes; this spec adds no code on that path.
Not covered: story 22 — the tilde form is existing behavior this spec does not touch; `expandHomeTarget` owns it.
Not covered: story 23 — an exclusion satisfied by writing nothing; review observes it in the diff.

WX20 is ticket evidence rather than retained gate coverage: the composed run needs a
real assignment ledger, a built worktree artifact, and a full gate, and a nested full
gate inside the test suite is outside the ordinary phases. The build demonstrates it
once and records the evidence; the map claims no standing red for it.

### Edge inventory

Walked against the profile's hostile-input checklist for shell CLIs.

- Wrapper path absent versus present-but-empty: distinct behaviors, both asserted — WX7 leaves the marker unset, WX11 sets it.
- Special file at the wrapper path: FIFO covered by WX9; device and socket share its predicate.
- Dangling symbolic link at the wrapper path: WX10.
- Live symbolic link at the wrapper path: emitted as given, because nothing executes the value and the doctor resolves what it reports.
- Worktree path containing spaces or glob characters: exec passes the path as one argument to the child process rather than through a shell, and the ledger stores it cleaned.
- Control bytes in the target: refused before resolution by the existing target check.
- A path read out of a file: worktree creation canonicalizes through symbolic links and the ledger refuses a stored path that is not absolute and cleaned, so this class is disposed of before exec reads it.
- Required tool missing from PATH: unchanged — exec launches the operator's argv and adds no tool lookup.

**Won't handle:** a worktree of a linked project repository, whose wrapper sits under the vendored kit rather than at the worktree's own wrapper path — exec leaves the marker unset and the child behaves as it does today; the surviving in-scope caller is any such child invoked as `bench` or through the kit's wrapper, both of which set the marker themselves.

**Won't handle:** the inherited-verdict refusal FT223 tracks — a decision-required item with competing fixes the reviewer owns; the surviving in-scope caller is every exec child that does not meet a partial verdict.

**Won't handle:** a caller that sets the run-binary variable and expects the child to inherit it — deliberately refused since the kit-leak decision, and refused again here on freshness grounds; the surviving in-scope caller is any caller relying on the child to build from its own tree.

**Won't handle:** collapsing the two existing spellings of the launch-trust predicate into one — a pre-existing duplication between the run-binary validator and the gate entry; the surviving in-scope caller is every current caller, since both spellings agree today.

## Ownership fences

Ticket `own-the-run-binary-in-an-exec-child`:

- `internal/worktree/exec.go`
- `internal/worktree/exec_test.go`

Ticket `name-the-next-action-in-the-gate-entry-refusal`:

- `.bench/gate.sh`
- `internal/conformance/gate_entry_test.go`

The two fences are disjoint, so the tickets may be written concurrently in
separate worktrees.

## Out of scope

- **Resolve FT223's inherited-verdict refusal.** A reviewer decision between two
  named fixes, tracked with four occurrences. 4 edits, 2 gate runs.
- **Emit an absolute path from `bench worktree path`.** The tilde form is a
  deliberate portable affordance every worktree command expands; changing it is
  its own compatibility question. 3 edits, 1 gate run.
- **Collapse the two spellings of the launch-trust predicate.** The run-binary
  validator and the gate entry each state it; merging them touches the oracle's
  entry contract. 5 edits, 2 gate runs.
- **Advertise the exec pairing in the working agreement.** Dropped once the
  composition works: the prose would advertise an enforcement, which the code
  standard names as a defect. 2 edits, 1 gate run.

## Further notes

`internal/runbinary` calls the authenticated executable a *Selection*, and
`CONTEXT.md` carries no term for it. A glossary entry would settle whether
"selection", "run binary", and "private binary" name one thing. Proposed, not
taken — glossary upkeep is `craft-domain`'s and the reviewer's.
