# A light-path change lands without a commitment, and the drain delegates each light-path fix

Status: staged

Roadmap: FT391

Decision source: the reviewer-confirmed current conversation (2026-10-05).

Verification log: 2 iteration(s) to accept — Sonnet high ran the independent review. Iteration 1 returned accept after fixes: B1 found that the predicate ran before `readyFor`. It also found S1 to S6 and N1 to N4. Iteration 2 accepted, and the fold took three wording notes.

## Problem

The commitment gates every production path today, a one-ticket change included. `AuthorizeCandidate` in `internal/commitment/repository/candidate.go` calls `readyFor` for each candidate tree that changes a production path. `readyFor` in `readiness.go` refuses an assignment that holds no delivery binding. `bench commit` and the landing admission both reach that refusal.

A one-ticket change therefore pays a roadmap row, a commitment plan, an approval, a start, and a delivery landing. On 2026-10-05 the retirement of the delivered `commitment-delivery-integrity` spec met this refusal. `bench commit` refused ADR 0027 and the deletion of `reviews/commitment-delivery-integrity.md`, because those paths are not planning paths. The ADR26 retirement before it took a new outcome and four landings.

The drain cannot ship a light-path fix either. `.agents/commands/bench-drain.md` makes a drained light-path item wait for `bench commitment approve` and `bench commitment start`. `.bench/BENCH.md` lets a learning's light-path fix ship only when the active committed outcome needs it.

## Solution

The commitment gates spec implementations only. A light-path change needs no outcome, plan, or start. A light-path change carries one tickets-only folder that holds exactly one ticket, and it lands with that folder as `--spec`. Its production paths stay inside that ticket's `Writes:` line.

The light-path change keeps every other commitment guard. A change to `.bench/commitment.json` still needs an exact approval. A pinned roadmap row still changes only through its outcome. A tickets-only folder that a milestone approves as a deliverable is committed work, so it keeps the binding rule and records its delivery fact.

`/bench-drain` dispatches each light-path fix that its verdicts keep to a fresh write delegate. The delegate runs on the mid tier at high effort, in its own bench worktree. The dispatch does not wait for the active commitment. The operating guide, the drain command, the implementation command, and an ADR state the new scope.

## User stories

Line: opus / high.
Implementation-line reason: LP-C1 is the hardest chunk. One predicate decides both the commit and the landing admission, and the landing grades a tree without the folder. The spec pins each predicate and each refusal word. The seam is one existing function with read precedents, and package tests red each row cheaply.
Harder chunks: LP-C1.

### A light-path change commits and lands without a commitment

1. As a light-path author, I want `bench commit` to accept my ticket's `Writes:` paths unbound, so that a fix needs no outcome.
2. As a light-path author, I want `bench worktree land --spec <slug>` to publish my change unbound, so that one landing ships it.
3. As a light-path author, I want the exemption while a spec outcome is active, so that a light-path change never waits on the queue.
4. As a light-path author, I want a `(new)` entry to cover its file, so that a ticket that creates a file qualifies.
5. As a light-path author, I want a directory entry to cover each path below it, so that one entry owns a folder and no sibling.
6. As a reviewer, I want a retirement-shaped change to commit as light-path work, so that a delivered spec retires with no new outcome.

### The exemption stays narrow

7. As a reviewer, I want a folder with two tickets to keep the binding refusal, so that only a one-ticket change skips the commitment.
8. As a reviewer, I want a path outside `Writes:` refused with the path named, so that a light-path change carries no unrelated work.
9. As a reviewer, I want paths that only two tickets cover together refused, so that one light-path change carries one ticket.
10. As a reviewer, I want an approved tickets-only folder to keep the binding rule, so that committed work records its delivery fact.
11. As a reviewer, I want a ticket with a grammar fault not to qualify, so that the exemption reads only a valid `Writes:` line.
12. As a reviewer, I want a ticket entry that is a symbolic link not to qualify, so that no outside bytes supply `Writes:`.
13. As a reviewer, I want a spec-less light-path landing refused with `--spec` named, so that every light-path landing closes its folder.
14. As a reviewer, I want a light-path edit of `.bench/commitment.json` refused, so that only my plan changes the commitment.
15. As a reviewer, I want a light-path edit of a pinned row refused, so that a pinned row changes only through its outcome.
16. As a reviewer, I want a legacy continuation to keep its scope refusal, so that a light-path ticket never widens that scope.
17. As a light-path author, I want to remove an unpinned roadmap row, so that the change can close the row it settles.
18. As a light-path author, I want an asset beside the ticket not to count, so that a ticket can keep its assets.
19. As a reviewer, I want a refusal to quote a hostile path, so that a control byte cannot split the refusal line.

### One owner reads the `Writes:` grammar

20. As a maintainer, I want one owner for the `Writes:` split and cover rules, so that preflight and the commitment check agree.

### The guidance states the new scope

21. As an agent, I want the guide to limit `bench commitment start` to spec work, so that I start no outcome for a light-path change.
22. As an agent, I want the guide to state the light-path observable and its boundary, so that I know which change skips the commitment.
23. As an agent, I want the guide and the drain to state the drain's delegate route, so that light-path fixes ship in the drain.
24. As an agent, I want the retired admission sentences kept out, so that no sentence tells me to admit a light-path item.
25. As an agent, I want the implementation command to exempt a light-path change, so that its entry does not send me to start.
26. As a teammate, I want an ADR to record the commitment scope, so that the decision outlives this spec.
27. As a teammate, I want ADR 0023 to name the light-path delegate route, so that the ADR states the current authorship.

### A bound assignment keeps its admission

28. As a bound author, I want a light-path folder in my tree to change nothing, so that my commit keeps its current admission.

## Implementation decisions

### The light-path predicate

- One unexported predicate in `internal/commitment/repository` decides whether a candidate is light-path work. It runs only after `readyFor` refuses, at the `readyFor` call in `authorizeCandidate`.
- When `readyFor` admits, the candidate is admitted, and the predicate reads no tree. A bound commit therefore pays no tree walk and never names another change's ticket.
- When `readyFor` refuses, the predicate either admits the candidate or returns a refusal. It returns the `readyFor` refusal unchanged when no qualifying folder exists.
- The predicate runs only on the binding fault of `readyFor`. It returns each other `readyFor` error unchanged, for example an unreadable policy.
- The transition check, `protectedCandidate`, the production-path loop, and the continuation scope check run before `readyFor`, so their refusals win first.
- A tickets-only folder qualifies when four conditions hold. `spec.TicketsOnly` accepts it in the read tree, and its `tickets/` subtree holds exactly one `.md` regular file. That ticket parses with no diagnostic through `tickets.ParseTicket` with the empty tag. No milestone of the current `main` policy approves the folder as a deliverable.
- The predicate counts the `tickets.Ext` entries at every depth below `tickets/`, so `tickets/sub/two.md` is a second ticket. Another extension is an asset, and the predicate ignores it. A `tickets.Ext` entry whose tree mode is not `100644` or `100755` refuses the folder.
- A qualifying folder admits the candidate when its ticket's `Writes:` entries cover every production path. A production path is the existing `production` slice in `authorizeCandidate`.

### The tree reader

- The new file `internal/commitment/repository/light_path.go` owns the read of every tickets-only folder in a Git tree. `tickets.Enumerate` reads the filesystem, and `spec.TicketsOnlyFolders` reads only the working tree, so neither serves a tree object.
- The reader runs one `git ls-tree -r -z` of `specs` in the read tree, after the precedent in `repository.go`. That listing gives each entry's mode and path, so the folders, the ticket count, and the modes come from one read.
- The reader confirms each folder through `spec.TicketsOnly(spec.CommitTree(root, tree), slug)`. The reader of `CommitTree` resolves `<tree>:<path>`, which Git accepts for a tree object, so `internal/spec` needs no edit.
- The reader reads the one ticket through the bounded `git.ReadTreeFile` and parses it through `tickets.ParseTicket`.

### Commit and publication modes

- `bench commit` reaches `AuthorizeCandidate`, which grades in commit mode. Commit mode reads every tickets-only folder of the candidate tree. One qualifying folder must cover all production paths.
- A landing reaches `admitPublication`, which grades in publication mode. Publication mode reads only the folder that the landing's `--spec` names, at the reviewed source commit `Delivery.Source`. `published.Tree` removes that folder from the graded tree before admission, so the graded tree cannot supply the ticket.
- A publication with no delivery is a spec-less landing. It is never light-path work. Its refusal adds the `--spec` route when the candidate tree holds a qualifying folder that covers every production path.
- In both modes, a bound assignment stays admitted through `readyFor`, whatever folders the tree holds.
- A publication whose delivery names a staged spec keeps the current rules.

### Refusal words

- A path that no qualifying ticket covers: `production path %q is outside the Writes line of light-path ticket %q; add the path to that line, or run bench commitment start`. The second operand is the ticket path. Commit mode names the first uncovered production path in change order, and the first qualifying ticket in folder order.
- Paths that each qualifying ticket covers only in part: `production paths span more than one light-path ticket; a light-path change carries one ticket`.
- A spec-less landing of a covered light-path tree: the existing binding refusal, then `; land the light-path change with --spec %q`, where the operand is the folder slug.
- The Writes refusal and the span refusal replace the binding refusal only when `readyFor` refuses and a qualifying folder exists.
- When some production path is outside every qualifying ticket, the Writes refusal wins. The span refusal applies only when each production path is inside some qualifying ticket and no one ticket covers all of them.
- Each other case keeps its current refusal unchanged, including `assignment has no current delivery binding`.
- Each operand renders through Go `%q`, so a space or a control byte stays inside the quotes.

### One owner for the `Writes:` grammar

- `internal/tickets` gains two exported functions. `WritesPath(entry)` returns the tree path that one `Writes:` entry names and whether it carries `(new)`. `Covers(entry, path)` reports whether that entry names path exactly or contains it at a `/` segment boundary.
- `splitWritesEntry` and `pathCovered` in `internal/preflight` retire, and their callers call the `internal/tickets` functions. `prefixCovers` in `internal/tickets/registry_data.go` delegates to the same segment rule, so the package holds one copy.
- `inScope` in `internal/commitment/repository/candidate.go` calls `tickets.Covers` for each scope entry. Its rule is the same segment rule, so the legacy scope check keeps its behavior.
- `Covers` takes an entry that is already split. Each caller that holds a raw `Writes:` entry calls `WritesPath` first.
- The new import edge `internal/commitment/repository` → `internal/tickets` is acyclic. `go list -deps ./internal/tickets` names no commitment, preflight, or spec package.

### Roadmap rows

- A light-path change can remove an unpinned row: its `ROADMAP.md` index line and its `roadmap/FT<n>.md` detail file. Both paths are planning paths, and `protectedCandidate` already refuses a change to a pinned row. Bench does not close a row on its own for a light-path landing.

### Guidance text

The operating guide's commitment paragraph keeps six sentences. It changes two of them to these exact sentences:

- "A spec implementation starts only through `bench commitment start` for the eligible outcome that `bench status` names."
- "Any other finding, idea, learning, or drained item that needs a spec stays uncommitted intake; minimal support for the active outcome stays in it."

A new operating-guide paragraph follows it, with these exact sentences:

- "**A light-path change needs no commitment.**"
- "It carries one tickets-only folder with exactly one ticket, and it lands with that folder as `--spec`."
- "Its production paths stay inside that ticket's `Writes:` line."
- "It may remove a roadmap row that no committed outcome pins as a source; a pinned row changes only through its outcome."

The learning paragraph keeps six sentences. Its heading becomes "**Delegate each light-path fix.**" It keeps the two anchored dispatch sentences and the done-claim sentence verbatim. It drops its last sentence, which the drain now owns. Its first sentence changes, and one sentence follows it:

- "At any point in the workflow, `/bench-implement-spec` included, a `bench learning` entry's light-path fix ships at once only when the active outcome needs it."
- "`/bench-drain` dispatches every other light-path fix that its verdicts keep, whatever the active commitment is."

The fix paragraph changes its second sentence, and the capture paragraph adds one sentence. The light-path table row stays verbatim.

- "A small defect that the active outcome does not need goes to `bench learning`, and the drain delegates its light-path fix."
- "A drain implements a light-path idea under the light-path fix rule instead."

The drain command's step 5 changes its intake sentence and three later sentences to these exact sentences, and adds a fifth:

- "Every kept item that needs a spec becomes uncommitted intake: a new row, or a merge into the row that already covers it."

- "The commitment rule in `.bench/BENCH.md` decides whether spec intake starts; a drain approval never admits it."
- "A drained light-path item needs no commitment: the drain dispatches it under the light-path fix rule in `.bench/BENCH.md`."
- "The delegate writes the one ticket file, implements it, and lands it as the light-path row states."
- "The coordinator verifies the done-claim against the ticket's acceptance rows and the gate."

Three more drain sentences change:

- The learning sentence reads: "A learning entry with a light-path fix follows the light-path fix rule in `.bench/BENCH.md`."
- The next sentence reads: "Its verdict closes the entry by implementation after the fix lands."
- `## Delegate the evidence` reads: "Implement-now delegates may run while other reads continue."
- The next sentence reads: "Route their line through `craft-line` and keep their authorship under `.bench/BENCH.md`."

The implementation command's entry paragraph already holds six sentences. A new paragraph after it holds one exact sentence: "A light-path change needs no commitment start; `.bench/BENCH.md` owns that rule."

### Anchor and conformance registry

- The registry requires exactly twelve of the new and changed sentences, one row each: LP34 to LP40, LP43 to LP46, and LP55. The other new sentences are unanchored prose, and the Spec axis reviews them.
- The required guide sentences are the start route, the intake sentence, the observable, the `Writes:` boundary, the row rule, and the drain dispatch.
- The required drain sentences are the spec-intake sentence, the delegate route, the delegate sentence, the implement-now overlap, and the plural routing sentence. The required implementation sentence is the light-path start exemption.
- The family forbids three retired fragments. In the guide, `light path and fixes included` raises `commitment guidance: operating guide restored the commitment start for light-path work`.
- In the guide, `a light-path fix that the active committed outcome needs` raises `commitment guidance: operating guide restored the learning fix for the active outcome only`. In the drain, `is implement-now work only after` raises `commitment guidance: drain restored the commitment admission of a light-path item`.
- Two forbid rows retire, because the reviewer reopened drain-time implementation. They are the guide's `or close by implementation during that same drain` and the drain's direct learning-fix sentence. The guide's `a light-path fix that needs no reviewer decision` forbid row stays.
- `TestCommitmentGuidance` re-anchors two restore cases. The `unadmitted learning fix` case anchors on `**Delegate each light-path fix.**`. The `default implementation` and `declined row` cases anchor on `The coordinator verifies the done-claim against the ticket's acceptance rows and the gate.`
- The `registry_data.go` row for the retained implement-now route takes the new delegate sentence. Its diagnostic becomes `.agents/commands/bench-drain.md dropped the delegated implement-now light-path route`, and the canary `drain-implement-now-route` expects it.
- `TestCommitmentGuidance` and `TestRecurrenceMaintenanceContractCheckBites` take the same sentences as independent copies, so a dropped row fails its own case.

### Decision records

- A new ADR 0028 records the decision. The commitment gates spec implementations, and a one-ticket light-path change needs none. A pinned row and the policy still change only through a plan. The drain dispatches each kept light-path fix to a fresh write delegate.
- ADR 0023's last sentence names the light-path delegate route of the drain and of a phase that needs the fix.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| LP-C1 / `1-own-writes-grammar.md`, `2-admit-light-path-commit.md`, `3-admit-light-path-landing.md` | A light-path change commits and lands without a binding, and every other guard holds | LP1 to LP33, LP49 to LP53, LP56 | `internal/tickets`, `internal/preflight`, `internal/commitment/repository`, `internal/commit`, `internal/worktree` package tests | yes |
| LP-C2 / `4-state-light-path-guidance.md`, `5-record-commitment-scope-adr.md` | The guide, the drain, the implementation command, and the ADRs state the new scope | LP34 to LP48, LP54, LP55 | `internal/conformance` package tests and the guidance canary | no |

LP-C2 starts after the LP-C1 checkpoint, so the guidance never states an exemption that the gate does not yet admit.

The three LP-C1 tickets run in series on one integration source, because each one writes `internal/commitment/repository/candidate.go`. The two LP-C2 tickets write no shared file, so they can author concurrently after the LP-C1 review.

| Ticket | Blocked by | Delivered coverage |
| --- | --- | --- |
| [1. Move the Writes grammar to one owner in internal/tickets](tickets/1-own-writes-grammar.md) | none | LP27, LP28, LP29, LP30, LP31, LP32, LP33 |
| [2. Admit an unbound light-path commit inside its ticket's Writes line](tickets/2-admit-light-path-commit.md) | 1-own-writes-grammar.md | LP1, LP2, LP7, LP8, LP9, LP10, LP11, LP12, LP13, LP14, LP15, LP16, LP17, LP18, LP19, LP20, LP21, LP22, LP23, LP49, LP51, LP56 |
| [3. Publish an unbound light-path landing that names its folder](tickets/3-admit-light-path-landing.md) | 2-admit-light-path-commit.md | LP3, LP4, LP5, LP6, LP24, LP25, LP26, LP50, LP52, LP53 |
| [4. State the light-path scope in the guide, the drain, and the implementation command](tickets/4-state-light-path-guidance.md) | 3-admit-light-path-landing.md | LP34, LP35, LP36, LP37, LP38, LP39, LP40, LP41, LP42, LP43, LP44, LP45, LP46, LP54, LP55 |
| [5. Record the commitment scope in ADR 0028 and the delegate route in ADR 0023](tickets/5-record-commitment-scope-adr.md) | 3-admit-light-path-landing.md | LP47, LP48 |

Ticket 1 adds `tickets.WritesPath` and `tickets.Covers`, and ticket 2 calls them from the light-path predicate. Ticket 3 adds publication mode to the reader and the predicate of ticket 2. Tickets 4 and 5 wait for the LP-C1 checkpoint.

## Testing decisions

- A good test drives the real admission owner with a real fixture repository and a real tree object. It observes the refusal text, the commit exit, or the published `main` ref.
- The predicate rows attach at `Store.AuthorizeCandidate` and `Store.AdmitPublication`, after the precedent of `TestAdmitPublicationFrozenIdentity` and `TestAdmitPublicationClosureAuthority`.
- The commit rows attach at `bench commit` through `planningCommitRepo`, after the precedent of `TestCommitmentCommitBeforeEffects`.
- The landing rows attach at `bench worktree land` in `internal/worktree`, after the precedent of `TestCommitmentTicketsOnlyClosure`. The landing rows share one new top-level test, so the worktree test-count pin moves by one.
- The grammar rows attach at `internal/tickets` as pure functions. The preflight rows are the existing tests that read the moved rules.
- The guidance rows attach at `TestCommitmentGuidance`, `TestRecurrenceMaintenanceContractCheckBites`, and the guidance canary through `TestEveryRetainedFixtureBitesThroughRegisteredOwner`.

### Seam diagram

    trigger: bench commit -- <paths>
        │
        ▼
    candidate tree  ──▶  [ AuthorizeCandidate → authorizeCandidate (commit mode):
                           transition, protected rows, production paths,
                           continuation scope, readyFor, then the light-path predicate ]  ──▶  admitted or refused
                      ◀ tests attach here: Store.AuthorizeCandidate over a fixture
                        tree; bench commit in planningCommitRepo

    trigger: bench worktree land --spec <slug>
        │
        ▼
    composed tree  ──▶  [ published.Tree removes the folder → admitPublication
                          (publication mode): the --spec folder at Delivery.Source ]  ──▶  published or refused
                      ◀ tests attach here: Store.AdmitPublication with a delivery;
                        the landing fixture in internal/worktree

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LP1 | 1 | `bench commit -- change.go` in an unbound worktree exits 0 when the tree holds one tickets-only folder whose one ticket lists `change.go` in `Writes:` | planned TestCommitmentLightPathCommit in internal/commit/commitment_test.go | Today the commit refuses with `assignment has no current delivery binding`. |
| LP2 | 8 | `bench commit -- other.go` in that worktree exits 1 with stderr that contains `"other.go"`, the ticket path, and `bench commitment start` | planned TestCommitmentLightPathCommit in internal/commit/commitment_test.go | A predicate that ignores `Writes:` admits unrelated work. |
| LP3 | 2 | `bench worktree land --spec <slug>` of an unbound light-path assignment moves the `main` ref to the landing commit | planned TestCommitmentLightPathLanding in internal/worktree/commitment_light_landing_test.go | Today the landing admission refuses the unbound assignment. |
| LP4 | 2 | The LP3 published tree holds no `specs/<slug>` entry | planned TestCommitmentLightPathLanding in internal/worktree/commitment_light_landing_test.go | An exemption that skips `published.Tree` leaves the folder open. |
| LP5 | 2 | The LP3 published `.bench/commitment.json` equals its bytes at the landing base | planned TestCommitmentLightPathLanding in internal/worktree/commitment_light_landing_test.go | A light-path landing that records a delivery fact changes the policy. |
| LP6 | 13 | A spec-less landing of the LP3 change exits nonzero with output that contains `--spec` | planned TestCommitmentLightPathLanding in internal/worktree/commitment_light_landing_test.go | A spec-less exemption lands the change and leaves the folder behind. |
| LP7 | 3 | `Store.AuthorizeCandidate` admits a light-path tree while a second assignment holds the active binding and claim | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A predicate that reads the active claim blocks the light path behind the queue. |
| LP8 | 4 | `Store.AuthorizeCandidate` admits a new production path that a `Writes:` entry with the `(new)` marker names | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A cover that keeps the marker in the path never matches the new file. |
| LP9 | 5 | `Store.AuthorizeCandidate` admits `pkg/a.go` under the `Writes:` entry `pkg` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | An exact-only cover refuses a directory entry. |
| LP10 | 5 | `Store.AuthorizeCandidate` refuses `pkgx/a.go` under the `Writes:` entry `pkg` with a refusal that names `"pkgx/a.go"` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A bare string-prefix cover admits a sibling package. |
| LP11 | 6 | `Store.AuthorizeCandidate` admits a tree that deletes `reviews/x.md` and adds `docs/adr/9999-x.md` when the ticket lists both paths | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A predicate that grades added paths only refuses the deletion that a retirement makes. |
| LP12 | 7 | `Store.AuthorizeCandidate` refuses a folder with two tickets with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A count that reads only the first ticket admits a multi-ticket folder. |
| LP13 | 9 | `Store.AuthorizeCandidate` refuses two one-ticket folders that each cover one of two production paths with a refusal that contains `more than one light-path ticket` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A union of all tickets' `Writes:` lines admits a change that carries two tickets. |
| LP14 | 10 | `Store.AuthorizeCandidate` refuses a one-ticket folder that the active milestone approves as a deliverable with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | An exemption over an approved folder lets committed work skip its binding. |
| LP15 | 11 | `Store.AuthorizeCandidate` refuses a one-ticket folder whose ticket has no `Writes:` field with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A predicate that skips the grammar reads an empty cover set as permission. |
| LP16 | 12 | `Store.AuthorizeCandidate` refuses a folder whose one ticket entry has tree mode `120000` with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A reader that follows the link lets bytes outside the tree supply `Writes:`. |
| LP17 | 14 | `Store.AuthorizeCandidate` refuses a light-path tree that also edits `.bench/commitment.json` with a refusal that contains `candidate policy has no exact approval` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | An exemption placed before the transition check lets a light path change the commitment. |
| LP18 | 15 | `Store.AuthorizeCandidate` refuses a light-path tree that deletes a pinned row's detail file with a refusal that contains `candidate changes protected commitment` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | An exemption placed before `protectedCandidate` lets a light path remove a pinned row. |
| LP19 | 16 | `Store.AuthorizeCandidate` refuses a listed continuation run whose light-path ticket covers a path outside its scope with a refusal that contains `legacy continuation scope excludes` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | An exemption placed before the scope check widens a legacy continuation. |
| LP20 | 17 | `Store.AuthorizeCandidate` admits a light-path tree that removes an unpinned row's index line and detail file beside a covered production path | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A row rule that treats every row as pinned blocks the row that the change settles. |
| LP21 | 18 | `Store.AuthorizeCandidate` admits a folder that holds one ticket and `tickets/asset.txt` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A count of every entry calls the asset a second ticket. |
| LP22 | 18 | `Store.AuthorizeCandidate` refuses a folder that holds only `tickets/asset.txt` with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A count that accepts zero tickets admits a folder with no `Writes:` line. |
| LP23 | 19 | The LP10-style refusal for a production path that holds a space and an ESC byte contains the Go-quoted spelling of that path | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | An unquoted operand lets the control byte reach the terminal raw. |
| LP24 | 2 | `Store.AdmitPublication` with a delivery that names the folder at its source commit admits a composed tree that lacks the folder | planned TestLightPathPublication in internal/commitment/repository/light_path_test.go | A predicate that reads the graded tree finds no ticket after `published.Tree`. |
| LP25 | 8 | `Store.AdmitPublication` with that delivery refuses a composed tree with an uncovered production path with a refusal that names the path | planned TestLightPathPublication in internal/commitment/repository/light_path_test.go | A publication mode that skips the cover lets a landing carry unrelated work. |
| LP26 | 13 | `Store.AdmitPublication` with no delivery refuses a covered light-path tree with a refusal that contains `--spec` | planned TestLightPathPublication in internal/commitment/repository/light_path_test.go | A publication that reuses commit mode admits a spec-less landing. |
| LP27 | 20 | `tickets.WritesPath` returns `internal/x.go` and true for the entry `internal/x.go (new)` | planned TestWritesEntryCover in internal/tickets/writes_test.go | A split that keeps the marker returns a path that matches no tree entry. |
| LP28 | 20 | `tickets.WritesPath` returns `internal/d` for the entry `internal/d/` | planned TestWritesEntryCover in internal/tickets/writes_test.go | A split that keeps the trailing slash fails the segment comparison. |
| LP29 | 20 | `tickets.Covers` reports true for the entry `internal/d` and the path `internal/d/a.go` | planned TestWritesEntryCover in internal/tickets/writes_test.go | An exact-only cover refuses each path below a directory entry. |
| LP30 | 20 | `tickets.Covers` reports false for the entry `internal/d` and the path `internal/dx/a.go` | planned TestWritesEntryCover in internal/tickets/writes_test.go | A bare string-prefix cover claims a sibling package. |
| LP31 | 20 | The preflight `writes-resolve` row stays green over a `(new)` entry that the tree lacks | `internal/preflight/decision_test.go` (`TestWritesResolveAcceptsNewMarker`) | A caller that drops the marker split after the move reds this test. |
| LP32 | 20 | The preflight review fence keeps its segment boundary | `internal/preflight/command_review_test.go` (`TestCommandFencePrefixBoundary`) | A caller that adopts a bare prefix after the move reds this test. |
| LP33 | 20 | No copy of the `(new)` split or the segment rule survives outside `internal/tickets`, `inScope` included | review-owned: the Standards axis runs `rg` over `internal/preflight`, `internal/tickets/registry_data.go`, and `internal/commitment/repository/candidate.go` | A surviving copy drifts from the owner that the commitment check reads. |
| LP34 | 21 | Deleting the spec start-route sentence from `.bench/BENCH.md` raises `commitment guidance: operating guide dropped the spec start route` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A registry that keeps the old needle passes the deletion. |
| LP35 | 21 | Deleting the rewritten intake sentence raises `commitment guidance: operating guide dropped uncommitted intake` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A registry that keeps the old intake needle passes the deletion. |
| LP36 | 22 | Deleting `It carries one tickets-only folder with exactly one ticket, and it lands with that folder as` from the guide raises `commitment guidance: operating guide dropped the light-path observable` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A guide without the observable leaves the exemption undefined for the agent. |
| LP37 | 22 | Deleting `Its production paths stay inside that ticket's` sentence raises `commitment guidance: operating guide dropped the light-path Writes boundary` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A guide without the boundary lets an agent widen a light-path change. |
| LP38 | 17 | Deleting the row-rule sentence raises `commitment guidance: operating guide dropped the light-path row rule` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A guide without the rule leaves the pinned-row limit unstated. |
| LP39 | 23 | Deleting the guide's drain dispatch sentence raises `commitment guidance: operating guide dropped the drain light-path dispatch` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A guide without the dispatch keeps light-path fixes waiting on the commitment. |
| LP40 | 23 | Deleting the drain's delegate-route sentence raises `commitment guidance: drain dropped the delegate route for a light-path item` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A drain without the route keeps the admission wait. |
| LP41 | 24 | Restoring `light path and fixes included` beside the live guide raises `commitment guidance: operating guide restored the commitment start for light-path work` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A guide that restores the old start route sends light-path work to the queue again. |
| LP42 | 24 | Restoring the drain's old admission sentence raises `commitment guidance: drain restored the commitment admission of a light-path item` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A drain that restores the admission wait blocks each light-path fix. |
| LP43 | 24 | Deleting the drain's rewritten spec-intake sentence raises `commitment guidance: drain dropped its route to the commitment rule` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A registry that keeps the old drain needle passes the deletion. |
| LP44 | 25 | Deleting the implementation command's light-path sentence raises `commitment guidance: implementation dropped the light-path start exemption` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A phase entry without the sentence sends a light-path change to `bench commitment start`. |
| LP45 | 23 | The canary `drain-implement-now-route` raises `.agents/commands/bench-drain.md dropped the delegated implement-now light-path route` | `internal/conformance/fixture_bite_test.go` (`TestEveryRetainedFixtureBitesThroughRegisteredOwner`) | A canary that still expects the retained route fails against the new registry row. |
| LP46 | 23 | Swapping `Implement-now delegates may run while other reads continue.` raises `bench-drain does not allow implement-now work to overlap remaining reads` | `internal/conformance/recurrence_maintenance_contract_test.go` (`TestRecurrenceMaintenanceContractCheckBites`) | A contract that keeps the old sentence reds the live drain. |
| LP47 | 26 | ADR 0028 states that the commitment gates spec implementations and that a one-ticket light-path change needs none | review-owned: the Spec axis reads the ADR | An ADR is prose, and no executable check grades its decision. |
| LP48 | 27 | ADR 0023 names the drain's light-path delegate route | review-owned: the Spec axis reads the ADR | An ADR is prose, and no executable check grades its decision. |
| LP49 | 28 | `Store.AuthorizeCandidate` admits a bound assignment's production path that a qualifying folder in the same tree does not cover | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A predicate placed before `readyFor` refuses bound work beside an open light-path folder. |
| LP50 | 28 | `Store.AdmitPublication` admits a bound assignment's spec delivery whose composed tree holds a qualifying folder that does not cover its production paths | planned TestLightPathPublication in internal/commitment/repository/light_path_test.go | A publication mode that reuses the commit-mode check before `readyFor` refuses a bound spec landing. |
| LP51 | 7 | `Store.AuthorizeCandidate` refuses a folder with `tickets/one.md` and `tickets/sub/two.md` with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A count of the top level only admits a second ticket at depth. |
| LP52 | 7 | `Store.AdmitPublication` with a delivery that names a two-ticket folder refuses with a refusal that contains `assignment has no current delivery binding` | planned TestLightPathPublication in internal/commitment/repository/light_path_test.go | A publication mode that skips the count lands a multi-ticket folder unbound. |
| LP53 | 19 | The LP26 refusal for the folder slug `a b` contains `--spec "a b"` | planned TestLightPathPublication in internal/commitment/repository/light_path_test.go | An unquoted slug splits the printed repair command at the space. |
| LP54 | 24 | Restoring `a light-path fix that the active committed outcome needs` beside the live guide raises `commitment guidance: operating guide restored the learning fix for the active outcome only` | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A guide that restores the old condition blocks every drain fix. |
| LP55 | 23 | Swapping `Route their line through` sentence raises `bench-drain does not route implement-now work through craft-line and retained authorship` | `internal/conformance/recurrence_maintenance_contract_test.go` (`TestRecurrenceMaintenanceContractCheckBites`) | A contract that keeps the singular sentence reds the live drain. |
| LP56 | 9 | `Store.AuthorizeCandidate` refuses with a refusal that names the uncovered path when two qualifying tickets leave one production path outside both | planned TestLightPathCandidate in internal/commitment/repository/light_path_test.go | A span refusal that wins first hides the path that no ticket covers. |

### Edge inventory

The in-scope edges are the rows above: two tickets (LP12, LP51, LP52), a path outside `Writes:` (LP2, LP10, LP25), and two partial tickets (LP13). The others are an approved folder (LP14), a grammar fault (LP15), and a symbolic link (LP16). The guard edges are the policy edit (LP17), a pinned row (LP18), a legacy continuation (LP19), and an unpinned row (LP20). The landing edges are a spec-less landing (LP6, LP26) and an active spec outcome (LP7).

The hostile-input classes of `projects/benchkit.md` reach these surfaces:

- A path with a space or a control byte reaches the refusal operand, and LP23 pins its quoting.
- A symbolic link where a ticket file belongs reaches the tree read, and LP16 pins its refusal.
- An absent ticket against a present asset reaches the count, and LP21 and LP22 pin both sides.
- A special file cannot reach the predicate, because a Git tree holds only blobs, links, trees, and gitlinks.

Each excluded edge takes a Won't handle line:

**Won't handle** — a production path whose name holds a comma — the `Writes:` grammar splits on commas, and the spec route still lands that change.

**Won't handle** — a mid-work light-path fix that the active outcome does not need — the fix goes to `bench learning`, and the drain's delegate route ships it.

**Won't handle** — a copy of an approved folder under a new slug — `protectedCandidate` ignores bindings, but the approved binding stays open in `bench status`.

**Won't handle** — a `--spec` that closes a folder another change wrote — the cover and the gate still grade this change.

**Won't handle** — a light-path change made only of planning paths — it already needs no binding today, and LP1 keeps that route.

## Ownership fences

- `internal/tickets/writes.go`
- `internal/tickets/writes_test.go`
- `internal/tickets/registry_data.go`
- `internal/preflight/decision.go`
- `internal/preflight/closure.go`
- `internal/preflight/fence_writes.go`
- `internal/preflight/proposal.go`
- `internal/commitment/repository/candidate.go`
- `internal/commitment/repository/readiness.go`
- `internal/commitment/repository/publication.go`
- `internal/commitment/repository/light_path.go`
- `internal/commitment/repository/light_path_test.go`
- `internal/commitment/commitmenttest/`
- `internal/commit/commitment_test.go`
- `internal/worktree/commitment_light_landing_test.go`
- `internal/worktree/commitment_landing_fixture_test.go`
- `internal/worktree/parallel_census_test.go`
- `.bench/BENCH.md`
- `.agents/commands/bench-drain.md`
- `.agents/commands/bench-implement-spec.md`
- `projects/benchkit.md`
- `docs/adr/0023-each-ticket-gets-a-fresh-author.md`
- `docs/adr/0028-the-commitment-gates-spec-implementations.md`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/conformance/commitment_guidance_test.go`
- `internal/conformance/recurrence_maintenance_contract_test.go`
- `tests/canary/workflow-guidance-anchors/`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/row-next-grammar/token-table-lacks-kit-edit`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `CHANGELOG.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `reviews/light-path-commitment-exemption.md`

Build preflight binds the commitment, worktree, and anchors packages to the five command-registry files above, so each ticket that writes those packages names them. No ticket expects to edit them.

Ticket 4 also names the fixture canaries that pin each guidance file it edits, and the anchor registry files that name each such file. Ticket 4 edits a canary or a registry file only when its check reds. The workflow guidance anchors entry holds the route canary of row LP45.

`.bench/BENCH.md` and `.agents/commands/bench-implement-spec.md` are at their line budgets, so ticket 4 may raise their rows in `projects/benchkit.md`.

## Out of scope

- A ticket field that names the roadmap row that a light-path landing closes on its own. Estimate: 6 edits, 2 gate runs.
- The commitment streamlining for spec outcomes: a binding-required outlook state, a drafted bind proposal, and start against the branch policy. Estimate: 15 edits, 3 gate runs.

## Further notes

### Source trace

| source sentence (reviewer, 2026-10-05) | rows |
| --- | --- |
| "Commitment delivery binding gates spec implementations only." | LP1, LP3, LP7, LP34, LP47 |
| "Light path = a change that carries a tickets-only folder holding exactly one ticket and lands with that folder as `--spec`." | LP1 to LP6, LP12, LP15, LP16, LP21, LP22, LP24, LP26, LP36 |
| "`bench commit` and `bench worktree land` need no outcome, plan, or start for it." | LP1, LP3, LP7, LP24, LP44 |
| "A light-path landing may close a roadmap row that no committed outcome pins as a source." | LP20, LP38 |
| "A pinned row changes only through its outcome." | LP18, LP38 |
| "`/bench-drain` implements light-path recommendations through a fresh write delegate (mid tier, high effort, own bench worktree) regardless of the active commitment." | LP39, LP40, LP42, LP45, LP46, LP48, LP55 |
| "Light path recommendations should be implemented as part of the drain via delegate." | LP39, LP40, LP54 |
| "It should only affect spec implementations, not light path single ticket implementations." | LP49, LP50 |
| "This replaces the 'active committed outcome needs it' condition in `.bench/BENCH.md`." | LP39, LP41, LP54 |
| "The new outcome sits in quality-1 ahead of FT376; it needs a roadmap row as its source." | none: the commitment plan after staging owns it; see the flagged calls |
| "The unapproved ADR27 plan is dropped; the tickets-only folder `specs/adr-0027-retire-cdi` becomes the first light-path delivery after this ships." | LP11 |

### Flagged calls for reviewer veto

- **FT391 instead of FT392.** `roadmap/FT391.md` already owns this capability: "a hotfix lane lands a small change at once, outside the roadmap row and commitment ceremony". A new FT392 row would duplicate it, so this spec names FT391 and opens no row. The plan after staging binds this spec to FT391 with the obligation FT391 and moves FT391 ahead of FT376. The pinned FT391 row and its criterion change only through that plan.
- **The `Writes:` boundary.** The decision source names the folder, not the paths. This spec also requires every production path to stay inside the one ticket's `Writes:` line, so a one-ticket folder cannot carry unrelated work.
- **Approved folders keep the binding.** A tickets-only folder that a milestone approves is committed work, so the exemption skips it.
- **Mid-work fixes wait for the drain.** A light-path fix that the active outcome does not need goes to `bench learning`, and the drain ships it. The reviewer-confirmed conversation of 2026-10-05 names the drain route only: "Light path recommendations should be implemented as part of the drain via delegate."
- **Two forbid rows retire.** The reviewer reopened drain-time implementation, which these rows prohibited. The new forbid rows keep the old admission sentences out.
- **Row closure stays manual.** A light-path change may remove an unpinned row, but Bench closes no row on its own. The source says "may close", and an automatic route is out of scope.

### Flagged additions

- LP19 pins the continuation order. The source does not name continuations, but the ordering rule requires a row where two refusals compete.
- LP23 pins hostile quoting under the profile checklist.
- The `(new)` and segment rules move to `internal/tickets` under the one-source standard. The source does not name the move.

### Reader sweep

- Callers of `authorizeCandidate`: `AuthorizeCandidate`, which `internal/commit/commit.go` calls, and `admitPublication`. `AdmitPublication` and `PublishAdmitted` call `admitPublication`, and the landing in `internal/landing/landing.go` reaches both through `admission.Check` and `admission.Publish`.
- Callers of `readyFor`: `authorizeCandidate`, `ready`, and `closureAuthority`. `ready` serves `Ready`, which `internal/preflight/command.go` and `internal/preflight/charge_pack.go` call. It also serves `ReadyOutcome`, which `internal/shift/loop.go` calls. Each such caller grades a spec or a shift outcome, so none changes.
- Callers of `inScope`: `authorizeCandidate` once. Callers of `splitWritesEntry`: `internal/preflight/closure.go` twice, `fence_writes.go`, and `decision.go` twice. Callers of `pathCovered`: `closure.go`, `decision.go`, and `proposal.go`. Caller of `prefixCovers`: `BoundFiles`.
- Readers of the tickets-only folder: `spec.TicketsOnly` serves `ticketsOnlyAt`, the landing close in `land_identity.go`, and the resume in `land_resume.go`. Only the new predicate reads the ticket count.
- Input constructors: `commitmenttest.WriteTickets` writes `Light path ticket.`, which has no `Writes:` field. Its fixtures therefore never qualify, and their binding rows stay valid. A new helper in `internal/commitment/commitmenttest` writes a grammatical one-ticket folder.
- Guidance readers of the changed sentences: `.bench/BENCH.md`, `.agents/commands/bench-drain.md`, `.agents/commands/bench-implement-spec.md`, the anchors registry, `TestCommitmentGuidance`, `TestRecurrenceMaintenanceContractCheckBites`, and the canary `drain-implement-now-route`. The canary `drain-implement-now-row-fallback` pins a forbid row that stays, so it does not change. `.claude/commands` is a symbolic link to `.agents/commands`, and `.agents/skills/bench-drain/SKILL.md` only routes to the command.
- Shipped-surface claim words: ADR 0014 already says that "a light-path ticket retires through the verb". ADR 0026 says that "a routine roadmap drain never displaces committed work". Both stay true.

### Pre-review proof checklist

- `Cited symbols`: each symbol below resolves at `b21d6acb`.
  - `AuthorizeCandidate`, `authorizeCandidate`, `admitPublication`, `AdmitPublication`, `PublishAdmitted`, and `closureAuthority`
  - `readyFor`, `protectedCandidate`, `approvedTransition`, `continuationScope`, `PlanningPath`, and `planningPromotions`
  - `published.Tree`, `Delivery`, `spec.TicketsOnly`, `tickets.ParseTicket`, and `tickets.Enumerate`
  - `splitWritesEntry`, `pathCovered`, `prefixCovers`, `BoundFiles`, `WriteTickets`, and `planningCommitRepo`
  - each existing test in the map
- `Import edges`: `internal/commitment/repository` → `internal/tickets`. `go list -deps ./internal/tickets` names no commitment, preflight, or spec package.
- `Source-row clauses and occurrences`: the source trace quotes each closed decision. The occurrences are the two 2026-10-05 retirements in the Problem.
- `Promised field labels`: `production path`, `is outside the Writes line of light-path ticket`, `production paths span more than one light-path ticket`, `land the light-path change with --spec`, and each guidance sentence and diagnostic above.
- `Changed-function callers`: the reader sweep lists each caller of `authorizeCandidate`, `readyFor`, `splitWritesEntry`, `pathCovered`, and `prefixCovers`.
- `Copy survival`: LP33 names the review-owned check that fails when a copy of the split or the segment rule survives.
- `Rendered-shape readers`: the old needles `light path and fixes included`, `is implement-now work only after`, and `Implement that ticket in the retained session` appear in `.bench/BENCH.md`, `.agents/commands/bench-drain.md`, `internal/anchors/registry_commitment.go`, `internal/anchors/registry_data.go`, `internal/conformance/commitment_guidance_test.go`, and the route canary. `Retained implement-now work may run` and `Route its line` appear in the drain and `recurrence_maintenance_contract_test.go`. The restore anchors `**Delegate a light-path fix for a learning.**` and `Verify the diff against the ticket's acceptance rows and the gate.` appear in `commitment_guidance_test.go`. Each file is in the fence.

### Bootstrap authority

None. The predicate decides admission before the gate, and it launches no executable.

### Source disclosure

This session read the two learnings named in the charge only through their titles in this conversation. It did not re-read `capture/learnings.md` in this phase. The completion plan fence follows the ticket slice.

### Completion plan

```bench-completion-plan
{"version":2,"chunks":[{"id":"LP-C1","tickets":["1-own-writes-grammar.md","2-admit-light-path-commit.md","3-admit-light-path-landing.md"],"verification":[{"id":"t1-tickets","command":"bench test --package ./internal/tickets","probe":"In tickets.Covers, drop the slash segment boundary so a bare string prefix covers. TestWritesEntryCover must fail and the restore must be exact.","ticket":"1-own-writes-grammar.md"},{"id":"t1-preflight","command":"bench test --package ./internal/preflight","ticket":"1-own-writes-grammar.md"},{"id":"t1-commitment-repository","command":"bench test --package ./internal/commitment/repository","ticket":"1-own-writes-grammar.md"},{"id":"t1-conformance","command":"bench test --package ./internal/conformance","ticket":"1-own-writes-grammar.md"},{"id":"t2-commitment-repository","command":"bench test --package ./internal/commitment/repository","probe":"In the light-path predicate, admit a qualifying folder without the Writes cover of the production paths. TestLightPathCandidate must fail and the restore must be exact.","ticket":"2-admit-light-path-commit.md"},{"id":"t2-commit","command":"bench test --package ./internal/commit","ticket":"2-admit-light-path-commit.md"},{"id":"t2-conformance","command":"bench test --package ./internal/conformance","ticket":"2-admit-light-path-commit.md"},{"id":"t2-bench","command":"bench test --package ./cmd/bench","ticket":"2-admit-light-path-commit.md"},{"id":"t3-commitment-repository","command":"bench test --package ./internal/commitment/repository","probe":"In publication mode, return the readyFor refusal unchanged for a delivery that names a qualifying light-path folder. TestLightPathPublication must fail and the restore must be exact.","ticket":"3-admit-light-path-landing.md"},{"id":"t3-worktree","command":"bench test --package ./internal/worktree","ticket":"3-admit-light-path-landing.md"},{"id":"t3-conformance","command":"bench test --package ./internal/conformance","ticket":"3-admit-light-path-landing.md"},{"id":"t3-bench","command":"bench test --package ./cmd/bench","ticket":"3-admit-light-path-landing.md"}]},{"id":"LP-C2","tickets":["4-state-light-path-guidance.md","5-record-commitment-scope-adr.md"],"verification":[{"id":"t4-conformance","command":"bench test --package ./internal/conformance","probe":"In internal/anchors/registry_commitment.go, delete the require row for the light-path Writes boundary sentence. TestCommitmentGuidance must fail and the restore must be exact.","ticket":"4-state-light-path-guidance.md"},{"id":"t4-anchors","command":"bench test --package ./internal/anchors","ticket":"4-state-light-path-guidance.md"},{"id":"t4-prose-budgets","command":"bench test --check guidance-prose-budgets","ticket":"4-state-light-path-guidance.md"},{"id":"t4-bench","command":"bench test --package ./cmd/bench","ticket":"4-state-light-path-guidance.md"},{"id":"t5-prose","command":"bench test --check prose","ticket":"5-record-commitment-scope-adr.md"},{"id":"t5-conformance","command":"bench test --package ./internal/conformance","ticket":"5-record-commitment-scope-adr.md"}]}],"final_verification":[{"id":"coverage-check","command":"bench coverage --check specs/light-path-commitment-exemption/spec.md"},{"id":"tickets","command":"bench test --package ./internal/tickets"},{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"commitment-repository","command":"bench test --package ./internal/commitment/repository"},{"id":"commit","command":"bench test --package ./internal/commit"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"anchors","command":"bench test --package ./internal/anchors"}],"execution":{"mode":"delegate","run_id":"lpce-full-20261005","orchestrator_session":"claude:session_0128bLurgWw2evJv2LYjfDHH","author_limit":1,"assignments":{"1-own-writes-grammar.md":[{"session":"claude:lpce_t1","assignment":"54ddba1f1f8b2ffc7dc96608ec5037f3","model":"opus","effort":"high","source":"cea9238d3f3f52d05128bdd28281a81c7ca21d88","native_ref":"claude-agent:lpce_t1"}],"2-admit-light-path-commit.md":[{"session":"claude:lpce_t2","assignment":"54ddba1f1f8b2ffc7dc96608ec5037f3","model":"opus","effort":"high","source":"838c88265dfbe30d1cf1123fbbd20e96c03d8ec6","native_ref":"claude-agent:lpce_t2"}],"3-admit-light-path-landing.md":[{"session":"claude:lpce_t3","assignment":"54ddba1f1f8b2ffc7dc96608ec5037f3","model":"opus","effort":"high","source":"d9405cfc24514cbd990ca6f757f32a016b19950d","native_ref":"claude-agent:lpce_t3"}],"4-state-light-path-guidance.md":[],"5-record-commitment-scope-adr.md":[]}}}
```
