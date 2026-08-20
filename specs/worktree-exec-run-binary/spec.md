# Worktree exec run binary

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-20. Two forks closed by the reviewer: the contradiction resolves by having `bench worktree exec` point the child at the worktree's own `dist/bench` rather than by documenting the pairing or rewording the refusal alone; and the ticket set is rewritten whole to cover that resolution, the gate-entry reword, and an explicit FT223 exclusion.

Verification log: pending — the `fable`/high round has not run.

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
leaked kit names the caller's checkout. But the gate's owner lookup treats the
absence of both variables as "nobody owns this run" and returns nothing, so the
gate plan's environment never gains a selection and the gate entry refuses an
absent `BENCH_RUN_BINARY`. The child is left with no way to author the selection
the caller was right not to lend it.

The refusal then names only the variable. No operator sets it, nothing in the
message says what to run instead, and neither condition that produces it is
named. A session that reads it concludes the verb it just used is broken. That
is observed: one session abandoned `bench worktree exec` after meeting this
refusal and did not return to it.

## Solution

Close the gap where it opens. After `bench worktree exec` strips the caller's
routing, it points `BENCH_RUN_BINARY` at the child worktree's own `dist/bench`
when an artifact is present there. The child then owns a selection drawn from
its own tree, and the composed command works.

Discovery and trust stay separate. Exec answers only "is there an artifact
here" — a regular file at the artifact path. The gate entry keeps every trust
question it already asks: absolute, executable, not a symbolic link, physical
path equal to the given path, and fresh against the kit. Exec never vouches for
the binary it names, so the validator that launches it remains the independent
one, and the trust predicate keeps one source.

Separately, the gate entry's refusal gains the next action it lacks. The two
changes are independent: the first removes the common way to reach the refusal,
the second repairs the refusal for every other way.

## User stories

### The exec child owns its run

Line: opus / high. Kit code at an existing seam with an existing test harness;
the mid tier starts the build.

1. As a session working in an assignment worktree, I want `bench worktree exec <target> -- ./dist/bench gate` to run the gate, so that I can grade a worktree without leaving the sanctioned verb.
2. As a session, I want the exec child to receive `BENCH_RUN_BINARY` naming that worktree's own `dist/bench`, so that the gate entry finds a selection to trust.
3. As a session, I want the emitted value to be an absolute path, so that the gate entry's absolute-path check accepts it.
4. As a session, I want the exec child never to receive the caller's `BENCH_RUN_BINARY` value, so that a binary built for the caller's run never grades a different tree.
5. As a session, I want the exec child never to receive the caller's `BENCH_KIT` or `BENCH_WRAPPER`, so that the child resolves its own kit.
6. As a session, I want every unrelated variable I set to reach the exec child, so that exec stays a transparent runner.
7. As a session in a worktree carrying no `dist/bench`, I want the child to receive no `BENCH_RUN_BINARY`, so that an absent artifact does not become a false selection.
8. As a session in a worktree whose `dist/bench` is a directory, I want the child to receive no `BENCH_RUN_BINARY`, so that a non-file at the artifact path does not become a selection.
9. As a session in a worktree whose `dist/bench` is a FIFO, device, or socket, I want the child to receive no `BENCH_RUN_BINARY`, so that opening a special file cannot block the child.
10. As a session in a worktree whose `dist/bench` is a dangling symbolic link, I want the child to receive no `BENCH_RUN_BINARY`, so that a broken link is not read as an authoritative artifact.
11. As a session running a verb that never consults the selection, I want the child otherwise unchanged, so that this alters no command that does not need it.
12. As a reviewer, I want exec to assert nothing about the binary's trustworthiness, so that the gate entry stays the one authority that authenticates what it launches.

### The refusal names its next action

Line: opus / high. One message and its contract test; the mid tier starts the build.

13. As a session that reached the gate entry without a wrapper, I want the refusal to name the wrapper invocation, so that I can act on it without reading the gate script.
14. As a session, I want a relative `BENCH_RUN_BINARY` to meet that same refusal, so that one message covers both unusable forms.
15. As a session outside any worktree, I want the refusal to read correctly, so that the message serves every wrapper-less invocation rather than one of them.
16. As a reviewer, I want the refusal to keep naming `BENCH_RUN_BINARY`, so that the existing gate-entry contract test keeps binding to it.
17. As a reviewer, I want a valid absolute executable to be accepted exactly as it is today, so that the reword changes no verdict.
18. As a reviewer, I want the later refusals in that block unchanged, so that conditions an operator can already act on keep their wording.

### Reviewed exclusions

Line: opus / high. Stated so review can see what this spec declines.

19. As a reviewer, I want FT223's inherited-verdict refusal left to FT223, so that this spec does not pre-empt a decision I own.
20. As a reviewer, I want the portable tilde form of a worktree path left alone, so that a deliberate portability affordance is not removed as a defect.
21. As a reviewer, I want no new prose advertising the pairing, so that the fix removes the contradiction instead of documenting around it.

## Implementation decisions

**Exec discovers, the gate entry authenticates.** Exec's added knowledge is one
predicate: a regular file exists at the worktree's artifact path. Every trust
condition — absolute, executable, not a symbolic link, physical path equal to
the given path, fresh against the kit — stays spelled once, in the gate entry.
Restating any of them inside exec would put the trust decision in two places and
let them drift, and a component that names a path cannot authenticate that path
anyway.

**The worktree path is already validated.** The assignment ledger refuses a
stored worktree path that is not absolute and cleaned, so exec joins the
artifact path onto a value that already carries those properties rather than
re-deriving them.

**The variable's presence changes which owner the gate takes.** With the
variable present, the gate's owner lookup reuses the inherited selection instead
of returning nothing. That is the intended effect, and it is why a worktree
artifact that has fallen behind its own source must refuse rather than grade:
the gate entry's freshness check against the kit is what makes reuse safe, and
it already runs before the binary is handed the run.

**The reword is bounded by an existing contract.** The gate-entry contract test
asserts the refusal names the run-binary variable, so the reworded message keeps
that token. The gate script never ships to consumers, so wording may name
kit-relative invocations.

### Bootstrap authority

Each hop names how the validator authenticates the next executable before
launching it.

- The operator selects the wrapper on PATH. That choice is the trust root; Bench
  does not authenticate it.
- `bench worktree exec` resolves the target through the assignment ledger, which
  refuses an inactive assignment and a worktree path that is not absolute and
  cleaned. Exec then launches the argv the operator named. Exec authenticates
  nothing about that argv, because the operator authored it.
- The child reaches the gate entry. The gate entry authenticates the next
  executable — the run binary — before launching it, by five path predicates and
  a freshness check against the kit. It refuses rather than launching on any
  failure.
- The value exec supplies is an input to that check, never a warrant. Exec's
  discovery predicate cannot substitute for it, which is why the two are kept
  apart.

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
                          │            with the worktree artifact path
                          │            ◀ tests attach here: run a real child,
                          │              read the assignments it received
                          ▼
                    [ child argv ]  ──▶  ./dist/bench gate
                                              │
                                              ▼
                                    [ .bench/gate.sh entry ]
                                        ◀ tests attach here: run the script as a
                                          real subprocess, assert accept or refuse
                                              │
                                              ▼
                                    authenticated run binary

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WX1 | 2 | a child in a worktree holding a regular-file artifact receives the run-binary variable naming that worktree's artifact path | exec child environment | a build that keeps stripping the variable leaves the assignment absent |
| WX2 | 3 | the emitted value is an absolute path | exec child environment | a build joining a relative artifact path emits a value the gate entry refuses |
| WX3 | 4 | a child never receives the caller's run-binary value | exec child environment | a build that forwards instead of re-pointing lends a binary built for another tree |
| WX4 | 5 | a child receives neither the caller's kit nor the caller's wrapper variable | exec child environment | a build that stops stripping routing makes the child resolve the caller's checkout |
| WX5 | 6 | a child receives every unrelated variable the caller set | exec child environment | a build that rebuilds the environment from scratch silently drops operator state |
| WX6 | 7 | a child in a worktree with no artifact path receives no run-binary variable | exec child environment | a build that sets the variable unconditionally names a path that does not exist |
| WX7 | 8 | a child in a worktree whose artifact path is a directory receives no run-binary variable | exec child environment | a build testing only existence accepts a directory as the artifact |
| WX8 | 9 | a child in a worktree whose artifact path is a FIFO receives no run-binary variable | exec child environment | a build testing only existence names a special file that blocks on open |
| WX9 | 10 | a child in a worktree whose artifact path is a dangling symbolic link receives no run-binary variable | exec child environment | a build reading the link without resolving it treats a broken link as an artifact |
| WX10 | 11 | a child's environment differs from today's by the run-binary assignment alone | exec child environment | a build with an incidental environment change alters verbs that never consult the selection |
| WX11 | 12 | a child in a worktree whose artifact path is a live symbolic link to a regular file still receives the variable, and the gate entry refuses that value | exec child environment, gate entry subprocess | a build that filters symbolic links inside exec moves the trust decision out of its one source |
| WX12 | 13 | the gate entry with no run-binary variable refuses and names the wrapper invocation | gate entry subprocess | a build restating only the variable leaves the message unactionable |
| WX13 | 14 | the gate entry with a relative run-binary value refuses with that same message | gate entry subprocess | a build wording only the absent case leaves the relative case unactionable |
| WX14 | 15 | the refusal names no worktree | gate entry subprocess | a build wording the message for exec alone misleads every other wrapper-less caller |
| WX15 | 16 | the refusal text contains the run-binary variable name | gate entry subprocess | a build dropping the token breaks the standing gate-entry contract |
| WX16 | 17 | the gate entry with a valid absolute executable proceeds as it does today | gate entry subprocess | a build whose reword changes the accept path alters a verdict |
| WX17 | 18 | the non-executable, symbolic-link, and uncleaned-path refusals keep their present wording | gate entry subprocess | a build rewording the whole block changes messages an operator can already act on |
| WX18 | 1 | `bench worktree exec <target> -- ./dist/bench gate` over a clean worktree of this kit reports green | gate entry subprocess | a build correct at each seam separately can still fail composed |

Not covered: story 19 — FT223 is a separate decision-required item with its own competing fixes; this spec adds no code on that path.
Not covered: story 20 — the tilde form is existing behavior this spec does not touch; `expandHomeTarget` owns it.
Not covered: story 21 — an exclusion satisfied by writing nothing; review observes it in the diff.

### Edge inventory

Walked against the profile's hostile-input checklist for shell CLIs.

- Artifact path absent versus present-but-empty: both yield no assignment, rows WX6 and WX7.
- Special file at the artifact path: FIFO covered by WX8; device and socket share its predicate.
- Dangling symbolic link at the artifact path: WX9.
- Live symbolic link at the artifact path: WX11 — exec emits it, the gate entry refuses it.
- Worktree path containing spaces or glob characters: exec passes the path as one argument to the child process rather than through a shell, and the ledger stores it cleaned.
- A path read out of a file: the assignment ledger refuses a stored worktree path that is not absolute and cleaned, so this class is disposed of before exec reads it.
- Required tool missing from PATH: unchanged — exec launches the operator's argv and adds no tool lookup.

**Won't handle:** a worktree reached through a symbolic-linked ancestor — the emitted value differs from its physical path, so the gate entry refuses it with its reworded message; the surviving in-scope caller is any worktree in the ordinary pool, whose path carries no linked ancestor.

**Won't handle:** a worktree artifact that has fallen behind its own source — the gate entry's freshness check refuses it before it grades; the surviving in-scope caller is a worktree whose artifact was built from its current source.

**Won't handle:** the inherited-verdict refusal FT223 tracks — a decision-required item with competing fixes the reviewer owns; the surviving in-scope caller is every exec child that does not meet a partial verdict.

**Won't handle:** a caller that sets the run-binary variable and expects the child to inherit it — deliberately refused since the kit-leak decision; the surviving in-scope caller is any caller relying on the child to resolve its own tree.

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
- **Advertise the exec pairing in the working agreement.** Dropped once the
  composition works: the prose would advertise an enforcement, which the code
  standard names as a defect. 2 edits, 1 gate run.

## Further notes

`internal/runbinary` calls the authenticated executable a *Selection*, and
`CONTEXT.md` carries no term for it. A glossary entry would settle whether
"selection", "run binary", and "private binary" name one thing. Proposed, not
taken — glossary upkeep is `craft-domain`'s and the reviewer's.
