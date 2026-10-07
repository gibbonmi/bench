# Process lifetime ownership

Status: staged

Roadmap: FT362

Decision source: `specs/process-lifetime/decisions/ft362-process-lifetime.md` (ready compiled map)

Verification log: specification accepted after two repairs — independent ticket graph accepted first pass

Planning stage: specification and 13-ticket graph independently accepted; implementation not started

Source tip: `a382f4848d432815968f044820fad593e856e297`

Production baseline: `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`

## Problem

Bench repeats process-group mechanics across its production callers.
Several cancellation paths wait without a final deadline.
The gate and builder also poll group absence without a deadline.
Decision ticket 1 names the current caller census and each policy.
The roadmap's raw file count does not define migration scope.

A child exit proves neither stream completion nor safe resource release.
A nested Bench child can create a separate group before its own leader exits.
Current resource owners often remove resources through deferred cleanup or owner-PID recovery.
Decision ticket 6 identifies those owners and the resulting retention obligation.

This outcome serves every repository that links the kit.
Kit-only tests use the same lifetime and resource owners.
This planning stage changes no delivery commitment.
FT290 remains the active commitment.

## Solution

One subprocess owner performs process-group lifetime mechanics.
Each caller selects its command, signal, durations, streams, and normal-exit policy.
The owner returns separate process, cancellation, stream, and cleanup facts.
Resource owners publish durable protection before a child can use their resources.
They release resources only after they prove every applicable user complete.

An expired shutdown deadline returns an incomplete result.
That result retains affected resources and refuses required cleanup or publication.
Protection survives CLI exit and reaches nested Bench launches.
Recovery validates resource ownership and observes current absence.
Recovery never signals from a stale record.

## User stories

Line: gpt-5.6-sol / high

Implementation-line reason: The hardest chunks compose durable resource protection across nested CLIs and recovery owners.
The source fixes behavior, but these production seams remain unproved.
Real-process tests expose omissions after the production interfaces exist.
The kit leverage rule selects the bound mid model at high effort.

Harder chunks: PL-C2, PL-C3, PL-C5

Author line: gpt-6.1-sol / high, one draft and at most two review repair rounds

The reviewer's author selection does not authorize implementation.

1. As an operator, I want one lifetime owner, so that every controlled group shares the same mechanics.
2. As an operator, I want prelaunch cancellation, so that a stopped run starts no child.
3. As an operator, I want startup cancellation serialized, so that an unpublished process cannot escape shutdown.
4. As a caller, I want one waiter, so that concurrent completion cannot reap a child twice.
5. As a caller, I want the raw child outcome, so that cleanup failure cannot erase a nonzero exit.
6. As a caller, I want separate stream results, so that a complete child cannot certify truncated output.
7. As an operator, I want escalation, so that a resistant child cannot defeat a catchable shutdown request.
8. As an operator, I want a final deadline, so that unresolved teardown cannot hold the lifetime caller forever.
9. As an operator, I want honest group observations, so that a permission error cannot become absence proof.
10. As a caller, I want normal policy preserved, so that migration adds no group kill to ordinary commands.
11. As a gate operator, I want interruption and timeout separated, so that their exit policies and evidence remain useful.
12. As a builder caller, I want normal group drain, so that a returned executable has no unfinished builder group.
13. As a bounds caller, I want existing status policy, so that cancellation keeps its current classification.
14. As a release operator, I want complete evidence, so that unresolved cleanup cannot authorize promotion or publication.
15. As a test operator, I want complete JSON decoding, so that focused results retain package validation.
16. As a worktree operator, I want existing command inputs, so that environment, streams, and assignment behavior remain useful.
17. As a shell operator, I want normal release and signal retention, so that migration preserves assignment policy.
18. As a shift operator, I want bounded adapter cancellation, so that startup, signal, and deadline paths share one lifetime.
19. As a resource owner, I want prelaunch durable protection, so that later cleanup cannot remove a live child's resource.
20. As a resource owner, I want partial persistence to refuse launch, so that children use only fully protected resource sets.
21. As a resource owner, I want unknown identities retained, so that a crash cannot erase an unrecorded child.
22. As a resource owner, I want every concurrent user observed, so that one completed user cannot release another user's resource.
23. As a resource owner, I want sealed admission, so that no child enters between observation and removal.
24. As an enclosing owner, I want nested groups registered, so that a nested CLI exit cannot hide its separately grouped child.
25. As an enclosing owner, I want durable unresolved disposition, so that an exit code cannot substitute for cleanup proof.
26. As a cache operator, I want protection beside the shared lock, so that CLI death cannot authorize a dangerous clean.
27. As a prospective owner, I want one retention rule, so that direct close and recovery preserve surviving users.
28. As a binary owner, I want executable retention, so that unresolved children retain the executable they can still use.
29. As a test-directory owner, I want private-directory retention, so that unresolved children retain their home and temporary files.
30. As a worktree owner, I want protected checkout retention, so that release, reuse, and recovery preserve active evidence.
31. As a recovery operator, I want validated ownership and positive absence, so that stale records grant no signal or removal authority.
32. As a caller, I want minimal worker retention, so that an incomplete return does not keep unrelated resources.
33. As a Linux builder caller, I want the parent-death signal, so that abrupt owner death still reaches its direct builder child.
34. As a Darwin caller, I want its existing group policy, so that migration invents no Linux parent-death guarantee.
35. As a maintainer, I want mechanism and caller tests, so that one group-kill witness cannot hide an unmigrated caller.
36. As a maintainer, I want direct preservation comparisons, so that each caller retains output and exit policy.
37. As a maintainer, I want explicit exclusions, so that migration cannot silently change unrelated capture or shell behavior.
38. As an operator, I want named unresolved obligations, so that retained evidence explains cleanup failure.
39. As a resource owner, I want strict protection records, so that malformed or foreign data cannot authorize removal.
40. As an operator, I want completed protections removed, so that ordinary completed runs still release resources.

## Implementation decisions

### Terms and interfaces

| Term | Definition | Avoid |
|---|---|---|
| Process group | The OS signal target named by a live launch's PGID | Complete descendant tree |
| Lifetime outcome | Immutable process, cancellation, stream, and cleanup facts | Exit code alone |
| Protection | Durable resource identity and complete user records | PID lease alone |
| Resource user | One registered launch that can use a resource | Direct child only |
| Admission | Permission to register a user before launch | Removal authority |
| Safe release | Every applicable user complete under sealed admission | Owner death |
| Unknown identity | A startup interval without durable identity or proved no-start | Dead user |

The implementation adds these terms to `CONTEXT.md`.
The map remains the behavioral source.

Subprocess owns a `Lifetime` created from a prepared command, policy, streams, and registration.
Its production methods are `Start`, `Cancel`, and `Wait`.
`Cancel` takes the caller's cause and signal.
`Wait` returns an immutable bounded outcome and an eventual completion handle.
No caller calls `Cmd.Wait`, `Run`, or `CommandContext.Cancel` on that owned command.

The owner serializes prepared, starting, running, stopping, and returned transitions.
Cancellation before starting yields no-start.
Cancellation during Start remains pending until Start succeeds or fails.
After successful Start, the pending request initiates shutdown.

A start worker performs Cmd.Start while the control path accepts cancellation.
Cancellation starts its clock even while that worker has no published identity.
The final deadline can return incomplete with pending startup protection.
An eventual successful Start receives the remembered shutdown request.
If that deadline already expired, the owner forces termination without another grace.

One successful Start creates exactly one waiter.
A failed Start creates none.

| Outcome part | Required facts |
|---|---|
| Process | Start error, known PID/PGID, pending or complete wait, raw code, terminating signal, wait error |
| Cancellation | None, interruption, or deadline cause with the caller's requested signal |
| Streams | Per-stream complete, pending, failed, or forcibly closed status and error |
| Cleanup | Not requested, complete, or incomplete status with each unresolved obligation |
| Protection | Safe, protected, or uncertain disposition with resource and user identifiers |

The owner reads ProcessState after the waiter completes.
It reads mutable captured output after the output worker completes.
A returned snapshot never changes when an unfinished worker completes.
The completion handle publishes a new snapshot instead.
Unfinished workers retain only the command, endpoints, and resources they need.
They do not retain the complete caller graph.

Subprocess imports neither bounds nor resource-owner packages.
Bounds remains the duration policy owner and imports subprocess for its runner.
Subprocess accepts selected durations and read limits from its callers.
The owner preserves caller platform attributes while it applies group setup.

### Stream ownership

Owner-created pipes separate child wait from decoder or caller-writer completion.
The child receives a file endpoint.
The owner closes its duplicate child endpoint after successful Start.
A separate consumer owns the read endpoint and reports its terminal result.
The waiter can reap the child without awaiting that consumer.

Testreport replaces StdoutPipe with this adapter.
It keeps complete JSON decoding and package validation after natural EOF.
Forced closure, malformed JSON, and missing terminal events remain failures.
A forced EOF is not complete-output proof.
No implementation calls Wait before an existing StdoutPipe reader finishes.

Direct caller files stay caller-owned.
The owner closes only endpoints it created.
An arbitrary blocked writer or reader can leave a local worker pending after lifetime return.
The owner reports incomplete state and retains its dependent resources.
It promises no bounded whole-CLI output through arbitrary blocked caller I/O.

Worktree exec keeps its separate normal-output ExecWaitDelay.
After child completion, the adapter waits that existing interval for complete output.
Expiry closes only owned endpoints and reports incomplete output.
It does not add a normal-exit group kill.

Source: Go 1.25.14 os/exec, especially Wait, awaitGoroutines, and StdoutPipe.
The map source URL and prototype asset carry the source and probe limits.
The author reread the upstream source on 2026-10-06.
The disposable probe proves no production resource retention behavior.

### Shutdown and caller policy

A catchable shutdown starts the caller's grace.
The owner sends KILL when that grace expires with unresolved group work.
A single final cleanup window starts with forced termination.
An immediate-KILL policy starts that window immediately.
The window covers all remaining teardown obligations together.
It is not a fresh deadline per obligation.

Only current ESRCH proves group absence.
An answering probe means present.
EPERM and other errors mean uncertain.
A signal error remains evidence until a current absence observation resolves it.
A pending waiter, consumer, or uncertain group at the deadline yields incomplete.
Diagnostics name each unresolved obligation once.

Bounds owns the duration entries.
Existing finite graces remain unchanged.
New catchable paths use a two-second grace where none exists.
The final cleanup window is three seconds.
FixedWindow keeps teardown finite under the unbounded-verdict test switch.
These values are policy choices, not measured minima.

| Census caller | Normal policy | Cancellation policy and preserved facts |
|---|---|---|
| Gate runner | KILL and observe remaining group | INT/130 for interruption, TERM/124 for timeout, current two-second grace |
| Runbinary builder | KILL and observe remaining group | TERM, current two-second grace, executable validation |
| Bounds runner | Preserve observed result | Immediate KILL, current status and limit classification |
| Prep-release step | Preserve step result | TERM, new two-second grace, step order and attribution |
| External preflight | Preserve phase result | TERM, new two-second grace, input checks and promotion policy |
| Vulnerability scanner | Preserve scanner validation | TERM, new two-second grace, exit 3 and finding policy |
| Focused test runner | Preserve complete decoded report | INT, current builder grace, package validation |
| Package-list runner | Preserve package data | INT, current builder grace, selection and decode policy |
| Worktree exec | Preserve child result and output delay | INT/130, current three-second grace, environment, streams, assignment |
| Worktree shell | Preserve normal assignment release | Forward actual signal, current five-second grace, retained lease |
| Shift adapter | Preserve adapter result | Current signal, new two-second grace, deadline precedence and evidence taxonomy |

Decision ticket 1 supplies the exact files and production entries for this table.
A preserve-mode caller never signals a group merely because its leader exits.
Group presence can leave protection active after normal child completion.
Required resource release then fails without adding a normal-exit kill.
The raw child result remains available beside that failure.

Worktree exec can return child zero when its own required stream and cleanup obligations completed.
It releases no assignment on that normal return.
Its active protection still prevents checkout reuse or removal.
An inherited enclosing owner must consume that protection before green or resource release.
Required cleanup and publication owners cannot use raw child zero as their operation result.

The shell's current normal nonzero mapping remains outside this change.
Its lifetime outcome retains raw status for the separate shell repair.

### Durable protection schema

A dependency-leaf package under `internal/subprocess/protection/` owns protection, admission, aggregation, and recovery observations.
It imports the standard library and the planned durablefile publication leaf.
Both subprocess and resource owners use it.
It does not replace resource derivation, assignment ownership, bundle ownership, or evidence authority.

The resource owner supplies canonical identity, existing ownership proof, and cleanup exclusion.
The CLI composition root supplies the protocol's kind-validator table.
Each validator comes from the corresponding resource-owner package.
Subprocess receives that table through its protocol dependency without importing those packages.
Nested CLIs install the same table before descriptor admission.

A private store holds resource.json, lock, and the fixed users directory.
The store and users directory have mode 0700.
Each JSON record has mode 0600.
One user.json file per random launch identifier resides in its user directory.

The strict bounded reader rejects unknown fields, trailing JSON, duplicate identities, wrong modes, special files, and unreadable data.
It does not follow record symlinks.
Invalid records retain resources and prevent admission.

| Resource field | Schema 1 value |
|---|---|
| schema | Integer 1 |
| resource_id | Random owner-generated resource generation |
| kind | cache, prospective, binary, kit-test, worktree, or release-evidence |
| resource_path | Clean canonical absolute resource path |
| binding | Canonical cache identity, or common-directory identity plus existing artifact identity |
| owner_pid | Positive initiating PID, diagnostic only |
| admission | open or sealed |
| users_dir | Fixed relative user-directory name under the validated store |

| User field | Schema 1 value |
|---|---|
| schema | Integer 1 |
| resource_id | Exact resource generation |
| user_id | Random launch identifier shared across its resource union |
| actor_pid | Positive CLI PID that owns the launch and local workers |
| identity | pending, known, or no-start |
| pid, pgid | Positive values for known, absent for pending and no-start |
| process | pending or complete |
| streams | pending, complete, or incomplete |
| cleanup | pending, complete, or incomplete |
| obligations | Exact unresolved identifiers, empty only for proved completion |

OpenResource joins an existing open generation under the protocol lock.
It never overwrites existing users or changes their generation.
A sealed generation refuses new users until recovery proves safe removal.
After safe removal, a new operation can publish a fresh generation.
Unique temporary owners create their generation before any child uses the resource.

Pending marks the durable startup interval.
Known records a successful launch.
No-start records a locally proved failed or canceled start.
The same live owner can close pending after it proves no start occurred.
Recovery cannot infer no-start from actor death.
Unknown startup identity can require indefinite retention.

| Production protocol interface | Contract |
|---|---|
| OpenResource | Validate owner identity and join or publish its generation |
| Prepare | Register one launch across inherited and local resources before Start |
| Started | Publish the actual successful PID/PGID across that union |
| Finish | Publish process, stream, and cleanup disposition |
| SealAndObserve | Close admission and return safe, protected, or uncertain |
| Recover | Apply the same observation to a validated existing resource |

Protection encodes each strict JSON record, selects mode 0600, and authorizes its destination.
It calls internal/durablefile.Replace(path, data, mode) for record publication.
It implements no second generic temporary-write and rename algorithm.
The planned leaf requires an existing parent and reports Stage, cause, and Published.
Published means rename completed; it grants no launch or release authority.

Any replacement error prevents launch, including a post-rename error with visible new bytes.
Protection retains the visible or uncertain record for observation and repair.
A partial Prepare rollback uses the same replacement owner and retains uncertain state on error.
Protection owns directory creation and terminal record deletion durability because those are resource-protocol operations.
A deletion-sync failure returns uncertain cleanup.
Caller-supplied bounds limits keep the package independent of bounds.

### Build-entry prerequisite

Complete durable-file-replacement delivery must be accepted and landed before any process-lifetime implementation begins.
The prerequisite is specs/durable-file-replacement/spec.md with its leaf, real review-record caller, and native qualification.
PL-C1 starts after that complete outcome lands.
PL-C2 consumes its landed internal/durablefile API.
The prerequisite completes Bench whole-spec acceptance, reconciliation, and landing; D-A is its internal chunk.

durable-caller-migration is not a prerequisite.
It owns the remaining caller migrations as a separate delivery.
The split prerequisite specification is accepted, but its implementation is not delivered at this checkpoint.
If the leaf contract changes, revise this dependency before process build entry.

### Ordering, admission, and failure

Resource owners acquire outer cleanup exclusions before protocol locks.
Cache users acquire the existing shared cache lock first.
Protocol operations deduplicate generations and lock stores in canonical absolute-path order.
In-process users share descriptors under the existing POSIX record-lock rule.
No operation acquires a resource-owner exclusion while it holds a protocol lock.

Prepare validates the complete descriptor vector before its first record mutation.
It publishes pending records across open resources, then releases protocol locks before Start.
A cleaner can seal after registration, but that pending user prevents removal.
Sealing refuses new registration.
It permits identity and completion updates for already registered users.

A partial Prepare failure starts no child.
It tries to publish no-start for records it already wrote.
A rollback failure retains those records and returns uncertain disposition.
It never deletes partial records to pretend the transaction did not occur.
An untouched resource has no user from that attempt.

After successful Start, Started publishes known identity across the registered union.
A partial update leaves pending records intact and the launch unresolved.
The live owner attempts shutdown with its actual in-memory handle.
It never reconstructs signal authority from records.
Only that live owner can prove completion of its known startup interval.
If it dies first, pending records retain resources.

SealAndObserve seals and grades the complete user set under cleanup exclusion and protocol locks.
Removal stays under those exclusions until deletion completes.
Registered pending users retain resources even before Start.
Known present or uncertain groups retain resources.
Recovery sends no signal to recorded identities.

For known users, release requires current ESRCH for every recorded group.
It also requires complete local workers, or current actor ESRCH after all known groups are absent.
Actor absence resolves only local workers.
It does not resolve descendant presence or unknown startup identity.
A reused identity that answers a probe retains safely.
Age grants no removal authority.

Incomplete stream evidence remains failure evidence after later absence permits removal.
Recovery cannot convert it into green.
Normal completed users remove their records under the same locks.
A last completed user permits release only after sealed admission and complete-set validation.

### Nested CLI propagation

The reserved `BENCH_PROCESS_RESOURCES` environment entry carries a strict JSON descriptor vector.
Each descriptor names store path, resource generation, kind, canonical resource path, and binding.
It contains the complete inherited resource union plus local owners.
Each nested Bench lifetime calls Prepare before its separate group can use that union.
Outer group absence cannot stand in for inner group absence.

The receiver validates descriptors against strict published records and each resource owner's identity validator.
The validator checks the resource kind's recognized store location and existing owner identity.
The protocol then checks generation equality and fixed record names.
Descriptors cannot choose arbitrary files, symlinks, generations, repositories, or cleanup targets.
Malformed, missing, foreign, and sealed descriptors prevent launch.

Environment transforms retain this entry through their current routing policy.
Worktree --env cannot replace the owner-composed vector.

Descriptors grant registration authority only.
They grant no signal, deletion, or executable authority.
Cleanup still requires the local owner's existing proof.
The protocol executes no recovered content and adds no recovery verb.

Durable user records carry nested completion across CLI exit.
The enclosing owner observes stores after its direct child's result.
Missing completion, pending users, and unreadable records remain protected or uncertain.
A nested CLI exit zero cannot authorize release while a separate registered group remains unresolved.
Concurrent users stay distinct records.

### Resource-owner cycles

| Resource owner | Store and acquisition | Required consumers |
|---|---|---|
| Cache | bench-process-users/ under canonical cache, after shared lock | Holder completion, exclusive cleaner |
| Prospective | process-users/ under validated bundle, before checkout use | Direct Close, sweep, gate, lane, evidence inspection |
| Private binary | process-users/ under selection directory, before build or execution | Build-failure cleanup, Selection.Close, inherited users |
| Private test run | process-users/ under run root, before child use | Entries propagation, Close, phase and focused-test owners |
| Worktree and shift scratch | bench-process-users/ under worktree Git admin directory | Exec, shell, shift, reuse, explicit cleanup, pool reclaim |
| Release evidence | Store under repository-bound evidence owner before phase launch | Prep-release, preflight finalization, replacement, publication |

Cache completion clears ordinary protection before releasing its shared lock.
Incomplete completion can drop the process-owned lock while durable protection remains.
The cleaner reads protection under its exclusive lock.
Unreadable or unresolved protection refuses clean before go clean starts.
An absent cache reports zero and creates nothing.
An empty cache with complete protection remains cleanable.

Prospective ownership keeps ReadPublished, canonical repository binding, and Publish as its record authority.
Its next owner schema adds a protection-generation reference.
The strict reader recognizes schema 1 as legacy ambiguity.
A legacy record without complete user proof retains the bundle.

Sweep also requires current ESRCH for the initiating bundle owner.
A live owner can materialize an unused checkout without a group registration.
Direct Close uses its live owner object instead of a stale owner-PID claim.
Close and sweep both observe protection before Git registration removal or recursive deletion.
A fresh unused bundle remains removable after admission closes.

ADR0019's retained-descendant exception becomes an explicit retention condition.
Owner-death-only removal does not safely implement that exception.
The implementation updates its recovery rule and consequence to require complete user proof.
Existing owner-death fixtures gain complete proof or expect retention.
Canonical bundle ownership and the no-recovered-execution rule remain intact.

Gate and lane propagate complete lifetime outcomes through nested entry points.
RunCaptured must stop reducing the outcome to an integer before resource owners consume it.
Unresolved protection permits no green record, retained green evidence, or lane pass.
Protection persistence failure leaves the pending verdict posture.
Phase order, subject identity, and evidence reuse authority remain unchanged.

Selection.Close and KitTestRun.Close check protection before recursive removal.
Their enclosing owners consume their errors before returning success.
Inherited binaries remain outside prospective deletion scope.
Their actual owner protects every registered user.

Worktree protection resides outside files that reset or clean can remove.
Protected checkouts cannot return to the pool or pass cleanup apply.
The existing planner includes protection evidence in its fingerprint.
Apply rechecks that evidence under cleanup exclusion.
Pool reclaim consumes the same protected-user verdict before key removal.
The implementation introduces no second cleanup predicate.

Shift stores a lifetime handle instead of a mutable command before Start.
Signal and deadline paths request cancellation through that handle.
An actual startup test pauses the production transition before canceling a resistant adapter.
Unresolved protection prevents rollback, scratch removal, release, and reuse.
Recovery retains protected scratch even when the ordinary dirty set is empty.
Deadline precedence and the existing evidence taxonomy remain unchanged.

Cleanup diagnostics name unresolved waiter, group, stream, identity, and persistence obligations.
They retain the child outcome and affected artifact paths.
The existing bounded response owner projects these facts.
Required release failure cannot report cleanup success, green, or publication.

## Implementation chunks

The independent specification review accepted the checkpoint recorded in assets/spec-review.md.
The following ticket graph is a pending approval proposal.
Each successor waits for its predecessor's frozen chunk review checkpoint.
Ticket implementation preserves the accepted interface contracts above.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| PL-C1 / 01-own-lifetime-through-bounds.md | Bounded lifetime and unchanged bounds policy | PL2–PL19, PL28, PL66–PL69, PL124–PL125 | Lifetime and bounds package tests | no |
| PL-C2 / 02-protect-nested-bounds-users.md | Durable admission and nested resource protection | PL20–PL27, PL39–PL54, PL70–PL74, PL104–PL113, PL120 | Protection tests and nested CLI witness | yes |
| PL-C3 / 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md | Protected gate, cache, bundle, binary, and test resources | PL29–PL38, PL55–PL60, PL75–PL78, PL121, PL123 | Resource-owner and second-process cleanup witnesses | yes |
| PL-C4 / 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md, 10-preserve-focused-test-completion.md, 11-preserve-worktree-exec-and-shell.md | Preserved release, testreport, and worktree policies | PL61–PL65, PL79–PL82, PL114–PL115, PL119, PL122 | Caller package tests | no |
| PL-C5 / 12-retain-shift-scratch-through-recovery.md, 13-enforce-final-lifetime-closure.md | Shift recovery and final migration closure | PL1, PL83–PL103, PL116–PL118 | Actual shift race, caller differential matrix, native checks | yes |

The first chunk supplies the lifetime and stream contract through the bounds caller.
It does not require unmigrated callers to reach the owner.
The final chunk owns the unchanged PL1 census audit beside PL91 after every caller migration.

The second supplies protection and descriptor propagation through the landed API named in Build-entry prerequisite.
The third integrates required resource release and green authority.
The fourth migrates the remaining caller vertical paths.
The last closes recovery and the complete census.

## Completion plan

The version 1 plan declares future evidence, not executed implementation results.
The coordinator adds required version 2 fresh-author assignments before dispatch.
The complete landed durability prerequisite applies before the first ticket.
The implementation line remains the approved proposal, with no binding change.

The first two chunks expose the shared lifetime and protection seams through real bounds consumers.
Subsequent tickets use those seams through existing owner entries and identity adapters.
They introduce no competing process or protection policy owner.
PL-C3 consumer tickets remain separate green checkpoints inside their stable review chunk.
All caller migrations precede the final global census and differential reconciliation.

Worktree descriptor rows PL114 and PL115 join PL-C4 with their actual consumer.
This narrower consumer assignment preserves all predicates and the five stable chunk IDs.

Every ticket co-names its canonical registry, anchor, and fixture closure.
Direct blockers order each shared Writes path, including holder files kept unchanged.
Over-budget sources gain no lines without a cohesive move into that ticket's owned files.
No later ticket can pay that checkpoint's headroom debt.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "PL-C1",
      "tickets": [
        "01-own-lifetime-through-bounds.md"
      ],
      "verification": [
        {
          "id": "ticket-1-check-1",
          "command": "bench test --package ./internal/subprocess"
        },
        {
          "id": "ticket-1-check-2",
          "command": "bench test --package ./internal/bounds"
        },
        {
          "id": "ticket-1-omission",
          "command": "bench test --package ./internal/subprocess",
          "probe": "Omit remembered cancellation for a held Start, start a second waiter, or wait unconditionally on a blocked consumer. The corresponding real-child or controlled-start assertion must fail."
        }
      ]
    },
    {
      "id": "PL-C2",
      "tickets": [
        "02-protect-nested-bounds-users.md"
      ],
      "verification": [
        {
          "id": "ticket-2-check-1",
          "command": "bench test --package ./internal/subprocess/protection"
        },
        {
          "id": "ticket-2-check-2",
          "command": "bench test --package ./internal/bounds"
        },
        {
          "id": "ticket-2-check-3",
          "command": "bench test --package ./internal/gocache"
        },
        {
          "id": "ticket-2-check-4",
          "command": "bench test --check system"
        },
        {
          "id": "ticket-2-omission",
          "command": "bench test --check system",
          "probe": "Drop inherited descriptors from the real bounds launch. The nested-provider witness must fail through bench test --check system."
        }
      ]
    },
    {
      "id": "PL-C3",
      "tickets": [
        "03-migrate-gate-process-policy.md",
        "04-retain-selected-build-binaries.md",
        "05-protect-private-test-directories.md",
        "06-refuse-clean-of-protected-cache.md",
        "07-retain-prospective-gate-artifacts.md"
      ],
      "verification": [
        {
          "id": "ticket-3-check-1",
          "command": "bench test --package ./internal/gate"
        },
        {
          "id": "ticket-3-omission",
          "command": "bench test --package ./internal/gate",
          "probe": "Replace INT with TERM, map timeout to interruption, or omit normal group drain. The actual gate signal and descendant-survival fixtures must fail."
        },
        {
          "id": "ticket-4-check-1",
          "command": "bench test --package ./internal/runbinary"
        },
        {
          "id": "ticket-4-omission",
          "command": "bench test --package ./internal/runbinary",
          "probe": "Remove normal builder drain or delete the selection before all registered groups are absent. The live descendant or executable-retention assertion must fail."
        },
        {
          "id": "ticket-5-check-1",
          "command": "bench test --package ./internal/env"
        },
        {
          "id": "ticket-5-check-2",
          "command": "bench test --package ./internal/gate"
        },
        {
          "id": "ticket-5-omission",
          "command": "bench test --package ./internal/env",
          "probe": "Drop the private-run descriptor or discard KitTestRun.Close failure. The child-owned file survives while a successful phase result must disappear."
        },
        {
          "id": "ticket-6-check-1",
          "command": "bench test --package ./internal/gocache"
        },
        {
          "id": "ticket-6-check-2",
          "command": "bench test --package ./internal/gate"
        },
        {
          "id": "ticket-6-check-3",
          "command": "bench test --check system"
        },
        {
          "id": "ticket-6-omission",
          "command": "bench test --package ./internal/gocache",
          "probe": "Use only the dropped owner lock, treat unreadable records as empty, or publish green before terminal persistence. The second CLI cleaner and real gate record witnesses must fail."
        },
        {
          "id": "ticket-7-check-1",
          "command": "bench test --package ./internal/gate/prospectiveartifact"
        },
        {
          "id": "ticket-7-check-2",
          "command": "bench test --package ./internal/gate"
        },
        {
          "id": "ticket-7-check-3",
          "command": "bench test --check system"
        },
        {
          "id": "ticket-7-check-4",
          "command": "bench test --check docs-currency-workflow"
        },
        {
          "id": "ticket-7-omission",
          "command": "bench test --package ./internal/gate/prospectiveartifact",
          "probe": "Restore owner-death-only sweep, ignore a nested registered group, or discard lane selection release failure. The real Close, second-process sweep, or green-record witness must fail."
        }
      ]
    },
    {
      "id": "PL-C4",
      "tickets": [
        "08-preserve-release-step-evidence.md",
        "09-bound-preflight-and-scanner-shutdown.md",
        "10-preserve-focused-test-completion.md",
        "11-preserve-worktree-exec-and-shell.md"
      ],
      "verification": [
        {
          "id": "ticket-8-check-1",
          "command": "bench test --package ./internal/preprelease"
        },
        {
          "id": "ticket-8-check-2",
          "command": "bench test --package ./internal/releaseevidence"
        },
        {
          "id": "ticket-8-omission",
          "command": "bench test --package ./internal/preprelease",
          "probe": "Keep the store inside a swapped stage or delete an abandoned protected stage from its owner marker alone. The resistant-child evidence sentinel must survive and the erroneous success assertion must fail."
        },
        {
          "id": "ticket-9-check-1",
          "command": "bench test --package ./internal/releasepreflight"
        },
        {
          "id": "ticket-9-omission",
          "command": "bench test --package ./internal/releasepreflight",
          "probe": "Restore an unbounded scanner wait or accept promotion from child zero with unresolved cleanup. The real resistant scanner or protected promotion assertion must fail."
        },
        {
          "id": "ticket-10-check-1",
          "command": "bench test --package ./internal/testreport"
        },
        {
          "id": "ticket-10-omission",
          "command": "bench test --package ./internal/testreport",
          "probe": "Treat forced EOF as a complete report, wait before decoder completion, or discard Selection.Close failure. The real report and retained executable witnesses must fail."
        },
        {
          "id": "ticket-11-check-1",
          "command": "bench test --package ./internal/worktree"
        },
        {
          "id": "ticket-11-omission",
          "command": "bench test --package ./internal/worktree",
          "probe": "Add a normal-exit group kill, strip the descriptor vector, or release the interrupted shell lease. The real command, environment, and lease assertions must fail."
        }
      ]
    },
    {
      "id": "PL-C5",
      "tickets": [
        "12-retain-shift-scratch-through-recovery.md",
        "13-enforce-final-lifetime-closure.md"
      ],
      "verification": [
        {
          "id": "ticket-12-check-1",
          "command": "bench test --package ./internal/shift"
        },
        {
          "id": "ticket-12-check-2",
          "command": "bench test --package ./internal/worktree"
        },
        {
          "id": "ticket-12-check-3",
          "command": "bench test --check system"
        },
        {
          "id": "ticket-12-omission",
          "command": "bench test --package ./internal/shift",
          "probe": "Read cmd.Process before the serialized transition, let an empty dirty set authorize removal, or trust a stale cleanup fingerprint. The actual resistant adapter and protected scratch witnesses must fail."
        },
        {
          "id": "ticket-13-check-1",
          "command": "bench test --package ./internal/conformance"
        },
        {
          "id": "ticket-13-check-2",
          "command": "bench test --package ./internal/subprocess"
        },
        {
          "id": "ticket-13-check-3",
          "command": "bench test --package ./internal/runbinary"
        },
        {
          "id": "ticket-13-check-4",
          "command": "bench test --check system"
        },
        {
          "id": "ticket-13-check-5",
          "command": "bench test --check docs-currency-workflow"
        },
        {
          "id": "ticket-13-omission",
          "command": "bench test --package ./internal/conformance",
          "probe": "Reintroduce a private scanner or shift loop, bypass one production caller, or certify output from child zero. The actual-tree audit or named caller differential must fail; restore and re-run the same check."
        },
        {
          "id": "native-linux",
          "command": "bench test --package ./internal/runbinary --run 'TestBuilderChildDiesWithAnOwnerThatNeverDrains'",
          "probe": "On native Linux, omit Pdeathsig from the actual builder attributes. The direct builder child survives and the native witness fails, then passes after restoration."
        },
        {
          "id": "native-darwin",
          "command": "bench test --package ./internal/runbinary --run 'DarwinPolicy'",
          "probe": "On native Darwin, omit the actual Setpgid setup. The resistant-group policy witness fails, then passes after restoration."
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "ordinary-integration",
      "command": "bench test --package ./..."
    },
    {
      "id": "system-integration",
      "command": "bench test --check system"
    },
    {
      "id": "coverage",
      "command": "bench coverage --check process-lifetime"
    },
    {
      "id": "ticket-grammar",
      "command": "bench test --check ticket-grammar"
    }
  ]
}
```

## Testing decisions

Lifetime tests use real local children for signals, descriptors, and waiter behavior.
Protection tests use real files, record locks, and cooperating processes.
Controlled faults cover persistence and uncertain probes that real processes cannot produce reliably.
Final resource witnesses invoke production owners rather than modeled outcome decisions.

The protection seam is necessary beyond the lifetime seam.
Several resource owners recover without starting a child.
Cross-process protection cannot live in one in-memory handle.
Its deletion restores substantial admission, record, and recovery mechanics across those owners.
Lifetime deletion restores start, cancellation, escalation, wait, and stream mechanics across the census.

Prior art includes testreport cancellation, worktree exec and shell, runbinary descendant drain, and prospective strict-record tests.
The Linux parent-death test provides its platform witness.
Read enforcement includes the coverage parser, conformance registry, injected-port check, and cancel-signal registration check.
The author also read the gate entry and project profile before selecting these seams.

The existing ordinary test phase grades package tests.
The existing system phase grades kit-only CLI journeys through the sealed binary and kit root.
No new phase or gate order is required.
The migration-closure package test parses actual production call sites.
It rejects group setup, negative-PID termination, and private waiter/drain ownership outside subprocess and platform adapters.
It permits read-only absence observations and test-only process fixtures.

The closure check resolves import aliases and negative operands through syntax.
It cannot treat comments as executable mechanics.
A reintroduced private scanner or shift loop must turn the actual-tree check red.
Production caller witnesses complement this static check.
One static check or one group-kill test cannot prove migration.

Existing tests retain their purpose unless replacement evidence demonstrates the same omission red.
Record that red before removing a duplicated expectation or expensive journey.
Independent expectations require the repository's recorded omission exception.
No runtime or test-cost saving is an acceptance claim.

### Seam diagram

```text
caller command, policy, streams, inherited descriptors
                         |
                         v
             [ subprocess Lifetime ] ---> immutable outcome
                |               |
            real child      owned consumer
                +------- completion handle
                         |
                         v
             [ protection protocol ] <--- nested registrations
                         |
               sealed complete user set
                         |
                         v
             existing resource owner ---> release or retained failure
```

Package tests cross the production Lifetime and protection interfaces.
Caller tests cross current entry points and observe output, exit policy, and resource disposition.
System tests start a second CLI after the initiating CLI exits.
They exercise the real cleaner, sweep, nested groups, and shift recovery.

### Acceptance coverage map

Every named new test below is planned.
The coverage command therefore reports these rows as uncited until implementation supplies their test declarations.
No row claims an executed production witness during specification.
Each seam's fixture identifier fixes a concrete scenario.
The final implementation must demonstrate its stated omission red.

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PL1 | 1 | Every census entry reaches the shared lifetime owner. | planned `internal/conformance/process_lifetime_test.go` (OwnerEntries) | A production-entry audit detects a caller that retains its private owner. |
| PL2 | 2 | A pre-canceled lifetime starts no child. | planned `internal/subprocess/lifetime_test.go` (PreCanceled) | The child's start marker appears if cancellation follows launch. |
| PL3 | 3 | Cancellation during Start reaches the eventual resistant child. | planned `internal/subprocess/lifetime_test.go` (StartingCancel) | The child's readiness PID survives if cancellation loses the unpublished handle. |
| PL4 | 4 | A failed Start creates no waiter. | planned `internal/subprocess/lifetime_test.go` (SpawnFailure) | A waiter counter exposes a failed-start wait. |
| PL5 | 4 | Concurrent cancellation and completion create exactly one waiter. | planned `internal/subprocess/lifetime_test.go` (SingleWaiter) | An audited real-command owner detects a second Wait call. |
| PL6 | 5 | Exit seven retains raw status seven and complete output. | planned `internal/subprocess/lifetime_test.go` (Nonzero) | A generic failure projection loses the child status or output. |
| PL7 | 5 | Cleanup failure retains an already observed child exit. | planned `internal/subprocess/lifetime_test.go` (ExitAndIncomplete) | A replacement cleanup error erases the process result. |
| PL8 | 6 | Captured output remains unread until its worker completes. | planned `internal/subprocess/lifetime_test.go` (DelayedCapture) | A blocked consumer exposes premature mutable-buffer access under race checks. |
| PL9 | 6 | A direct caller file remains open after completion. | planned `internal/subprocess/lifetime_test.go` (CallerFile) | A subsequent caller write fails if teardown closes its file. |
| PL10 | 6 | Forced owned-pipe closure returns incomplete output. | planned `internal/subprocess/lifetime_test.go` (ForcedClose) | An EOF-only projection accepts truncated output. |
| PL11 | 6 | A blocked caller writer returns a pending stream snapshot. | planned `internal/subprocess/lifetime_test.go` (BlockedWriter) | Waiting on the writer holds the lifetime past its final deadline. |
| PL12 | 7 | A TERM-resistant group receives KILL after its grace. | planned `internal/subprocess/lifetime_test.go` (ResistantGroup) | A live descendant-survival oracle exposes missing escalation. |
| PL13 | 8 | An unresolved waiter returns incomplete by the final deadline. | planned `internal/subprocess/lifetime_test.go` (PendingWaiter) | An unconditional receive leaves the caller blocked. |
| PL14 | 8 | A descriptor-retaining descendant cannot hold canceled decoding forever. | planned `internal/subprocess/lifetime_test.go` (InheritedDescriptor) | The live descendant and decoder remain unresolved past the final window without owned closure. |
| PL15 | 9 | Current ESRCH marks the group absent. | planned `internal/subprocess/protection/observe_test.go` (Absent) | A positive-absence fixture cannot complete cleanup if ESRCH remains uncertain. |
| PL16 | 9 | An answering probe marks the group present. | planned `internal/subprocess/protection/observe_test.go` (Present) | The protected resource disappears if presence becomes absence. |
| PL17 | 9 | EPERM retains uncertain group state. | planned `internal/subprocess/protection/observe_test.go` (Permission) | Cleanup deletes the resource if every failed probe means gone. |
| PL18 | 9 | An unknown probe error retains uncertainty. | planned `internal/subprocess/protection/observe_test.go` (UnknownError) | An ESRCH-only check omission permits deletion after an unrelated failure. |
| PL19 | 10 | Preserve mode leaves a normal surviving group unsignaled. | planned `internal/subprocess/lifetime_test.go` (NormalPreserve) | The descendant dies before fixture cleanup if migration adds a blanket normal kill. |
| PL20 | 19 | Protection persists before the child's resource-use marker. | planned `internal/subprocess/protection/registration_test.go` (Prelaunch) | The marker precedes the synced user record if launch happens first. |
| PL21 | 20 | First user-publication failure prevents Start. | planned `internal/subprocess/protection/registration_test.go` (FirstPublishFailure) | A child marker appears if the persistence error is ignored. |
| PL22 | 20 | Later resource-publication failure prevents Start. | planned `internal/subprocess/protection/registration_test.go` (PartialPublish) | A child uses an unprotected second resource if partial success permits launch. |
| PL23 | 20 | Failed rollback retains its partial pending record. | planned `internal/subprocess/protection/registration_test.go` (RollbackFailure) | Removing partial records hides an uncertain transaction. |
| PL24 | 21 | Partial known-identity publication retains the launch disposition. | planned `internal/subprocess/protection/registration_test.go` (IdentityPublishFailure) | One completed write wrongly certifies the whole resource union. |
| PL25 | 21 | Actor death retains a pending startup identity. | planned `internal/subprocess/protection/recovery_test.go` (UnknownStartup) | Owner-PID-only recovery removes a potentially live child's resource. |
| PL26 | 22 | Two concurrent users require two complete group-absence proofs. | planned `internal/subprocess/protection/recovery_test.go` (OverlappingUsers) | The first completion deletes the second live user's resource. |
| PL27 | 23 | Sealed admission refuses new registration before Start. | planned `internal/subprocess/protection/registration_test.go` (SealedAdmission) | A late start marker exposes the cleanup check-then-launch race. |
| PL28 | 13 | Bounds preserves immediate-KILL cancellation classification. | planned `internal/bounds/lifetime_test.go` (BoundsCancel) | An added grace or changed canceled status breaks its policy witness. |
| PL29 | 11 | Gate interruption retains INT and exit 130. | planned `internal/gate/lifetime_test.go` (GateInterrupt) | The child's signal log exposes an incorrect TERM or KILL opener. |
| PL30 | 11 | Gate timeout retains TERM and exit 124. | planned `internal/gate/lifetime_test.go` (GateTimeout) | The cause fixture fails if timeout becomes interruption. |
| PL31 | 11 | Gate normal completion drains the remaining group. | planned `internal/gate/lifetime_test.go` (GateNormalDrain) | A readiness descendant survives if normal drain disappears. |
| PL32 | 11, 25 | An unresolved nested phase prevents retained green evidence. | planned `internal/systemtest/process_lifetime_test.go` (NestedGate) | The real inner group survives despite outer group absence and child zero. |
| PL33 | 25 | Incomplete cleanup prevents a lane pass record. | planned `internal/gate/lifetime_resources_test.go` (LaneIncomplete) | A pass appears if lane aggregation consumes only child codes. |
| PL34 | 12 | Builder normal completion drains its remaining group. | planned `internal/runbinary/lifetime_test.go` (BuilderNormalDrain) | The selection returns while a builder descendant remains live. |
| PL35 | 12 | Builder cancellation keeps TERM and its existing grace. | planned `internal/runbinary/lifetime_test.go` (BuilderCancel) | A signal log and descendant oracle detect premature KILL or absent escalation. |
| PL36 | 28 | Selection.Close retains a protected private executable. | planned `internal/runbinary/lifetime_resources_test.go` (BinaryClose) | The selection directory disappears if Close ignores registered users. |
| PL37 | 29 | KitTestRun.Close retains a protected private directory. | planned `internal/env/kit_run_lifetime_test.go` (KitRunClose) | The private home disappears while a registered child owns its files. |
| PL38 | 29 | The gate consumes private-directory Close failure. | planned `internal/gate/lifetime_resources_test.go` (PhaseCloseError) | A successful phase hides the enclosing resource failure. |
| PL39 | 39 | A foreign binding prevents registration. | planned `internal/subprocess/protection/record_test.go` (ForeignBinding) | An unrelated store gains a user if inherited paths supply authority. |
| PL40 | 39 | Malformed protection retains the resource. | planned `internal/subprocess/protection/record_test.go` (Malformed) | A parse failure becomes an empty user set and permits removal. |
| PL41 | 39 | A protection symlink prevents admission. | planned `internal/subprocess/protection/record_test.go` (Symlink) | The outside sentinel changes if the receiver follows the descriptor. |
| PL42 | 39 | Trailing JSON prevents admission. | planned `internal/subprocess/protection/record_test.go` (TrailingJSON) | A second record hides admission state if a JSON prefix is accepted. |
| PL43 | 39 | Unknown record fields prevent admission. | planned `internal/subprocess/protection/record_test.go` (UnknownField) | An unsupported wire shape gains current-schema authority. |
| PL44 | 39 | Wrong record mode retains the resource. | planned `internal/subprocess/protection/record_test.go` (WrongMode) | A nonprivate record becomes safe-release proof. |
| PL45 | 39 | Duplicate user identity makes the set uncertain. | planned `internal/subprocess/protection/record_test.go` (DuplicateUser) | Deduplication hides the second conflicting group identity. |
| PL46 | 23 | An admitted pending user blocks removal after sealing. | planned `internal/subprocess/protection/registration_test.go` (SealPending) | The resource disappears before its admitted child starts. |
| PL47 | 23 | An admitted user can finish after sealing. | planned `internal/subprocess/protection/registration_test.go` (FinishSealed) | Completed resources leak if sealing freezes existing updates. |
| PL48 | 22 | Reversed resource vectors acquire locks without deadlock. | planned `internal/subprocess/protection/registration_test.go` (LockOrder) | Two cooperating processes expose opposing lock acquisition order. |
| PL49 | 24 | A nested CLI registers its separate group before use. | planned `internal/systemtest/process_lifetime_test.go` (NestedRegistration) | The inner use marker precedes its durable protection if propagation disappears. |
| PL50 | 24 | Overlapping nested users retain independent records. | planned `internal/systemtest/process_lifetime_test.go` (NestedOverlap) | The second child loses its resource when the first CLI completes. |
| PL51 | 25 | Missing nested completion prevents outer safe release. | planned `internal/subprocess/protection/recovery_test.go` (MissingCompletion) | Child exit zero substitutes for a missing terminal record. |
| PL52 | 31 | Recovery requires current absence for every known group. | planned `internal/subprocess/protection/recovery_test.go` (RecoverKnown) | A stale complete field authorizes deletion while its group answers. |
| PL53 | 31 | Recovery sends no signal to recorded identities. | planned `internal/subprocess/protection/recovery_test.go` (NoStaleSignal) | The signal audit changes if recovery targets a reused PGID. |
| PL54 | 31 | Record age grants no removal authority. | planned `internal/subprocess/protection/recovery_test.go` (OldRecord) | An old answering identity loses its resource under an age-only rule. |
| PL55 | 26 | A second-process cache clean refuses unresolved protection. | planned `internal/systemtest/process_lifetime_cache_test.go` (SecondProcessClean) | The initiating CLI has exited and dropped its lock before the real cleaner attempts deletion. |
| PL56 | 26 | The exclusive cleaner retains unreadable protection. | planned `internal/gocache/protection_test.go` (UnreadableProtection) | A go-clean marker appears if unreadable state becomes an empty set. |
| PL57 | 26, 40 | Completed cache protection clears before shared-lock release. | planned `internal/gocache/protection_test.go` (CompletedHold) | An observer sees an unlocked active record if release order reverses. |
| PL58 | 27 | Prospective Close retains an unresolved group user. | planned `internal/gate/prospectiveartifact/lifetime_test.go` (DirectClose) | The checkout vanishes if protection applies only during sweep. |
| PL59 | 27 | A second-process sweep retains surviving prospective users. | planned `internal/systemtest/process_lifetime_artifacts_test.go` (SecondProcessSweep) | Owner death wrongly authorizes deletion while the registered group remains alive. |
| PL60 | 27 | A legacy prospective record retains ambiguous cleanup proof. | planned `internal/gate/prospectiveartifact/lifetime_test.go` (LegacyRecord) | Owner death removes a schema-1 bundle that records no complete user set. |
| PL61 | 14 | Prep-release retains evidence after incomplete phase shutdown. | planned `internal/preprelease/lifetime_test.go` (PrepIncomplete) | The evidence directory disappears behind a resistant descendant. |
| PL62 | 14 | External preflight refuses promotion after incomplete cleanup. | planned `internal/releasepreflight/lifetime_test.go` (ExternalIncomplete) | A promoted green record appears while protection remains pending. |
| PL63 | 14 | Scanner cancellation escalates a TERM-resistant child. | planned `internal/releasepreflight/lifetime_test.go` (ScannerCancel) | The scanner survives if its old unbounded wait remains. |
| PL64 | 15 | Focused-test cancellation resolves a descriptor-retaining descendant. | planned `internal/testreport/lifetime_test.go` (FocusedPipeCancel) | The descendant-survival oracle exposes a decoder blocked before Wait. |
| PL65 | 15 | Malformed focused JSON remains a decode failure. | planned `internal/testreport/lifetime_test.go` (MalformedJSON) | Child zero falsely certifies unreadable report data. |
| PL66 | 8 | Wait and stream obligations share one final cleanup window. | planned `internal/subprocess/lifetime_test.go` (SingleFinalWindow) | Sequential fresh deadlines exceed the selected shutdown bound. |
| PL67 | 8 | The unbounded-verdict switch keeps teardown finite. | planned `internal/bounds/lifetime_test.go` (FixedShutdown) | A resistant group outlives the outer wait if teardown uses VerdictWindow. |
| PL68 | 32 | An incomplete worker holds only required resources. | planned `internal/subprocess/lifetime_test.go` (WorkerOwnership) | An unrelated closeable sentinel remains retained by the unfinished owner. |
| PL69 | 6 | Blocked caller stdin yields a pending local worker. | planned `internal/subprocess/lifetime_test.go` (BlockedInput) | Wait blocks indefinitely if independent completion covers only output. |
| PL70 | 40 | A proved no-start removes its protection durably. | planned `internal/subprocess/protection/registration_test.go` (NoStartComplete) | A canceled prepared run leaks a terminal user without the no-start close path. |
| PL71 | 19 | Resource-record sync failure prevents first launch. | planned `internal/subprocess/protection/registration_test.go` (ResourceSyncFailure) | A child marker appears if renamed but unsynced protection permits use. |
| PL72 | 38, 40 | Deletion-sync failure returns uncertain cleanup. | planned `internal/subprocess/protection/registration_test.go` (DeleteSyncFailure) | The owner reports success despite failed durable terminal publication. |
| PL73 | 31, 32 | Actor absence resolves workers only after known groups are absent. | planned `internal/subprocess/protection/recovery_test.go` (WorkerActorDeath) | A live descendant loses its resource if actor death resolves the complete user. |
| PL74 | 39 | A missing inherited resource record prevents Start. | planned `internal/subprocess/protection/record_test.go` (MissingDescriptorTarget) | The missing path silently drops from the resource union. |
| PL75 | 27, 40 | A fresh no-user prospective bundle remains removable. | planned `internal/gate/prospectiveartifact/lifetime_test.go` (EmptyCurrentBundle) | A universal retention rule leaks unused ordinary bundles. |
| PL76 | 26 | An absent cache clean creates no directory. | planned `internal/gocache/protection_test.go` (AbsentCache) | Admission before the absence branch creates new state during a zero-result clean. |
| PL77 | 26 | An empty cache with complete protection remains cleanable. | planned `internal/gocache/protection_test.go` (EmptyCache) | A blanket record-presence refusal blocks ordinary cleanup. |
| PL78 | 25 | Gate terminal persistence failure permits no retained green. | planned `internal/gate/lifetime_resources_test.go` (GateTerminalFailure) | The green store appears if certification precedes safe terminal persistence. |
| PL79 | 15 | Focused output requires terminal events for every package. | planned `internal/testreport/lifetime_test.go` (IncompletePackage) | A partial decoded report passes under an EOF-only check. |
| PL80 | 15 | Package-list cancellation keeps current interruption attribution. | planned `internal/testreport/lifetime_test.go` (ListCancel) | A generic owner error loses the go-list diagnostic. |
| PL81 | 16 | Exec normal-output delay adds no normal group kill. | planned `internal/worktree/exec_lifetime_test.go` (ExecNormalPipe) | The descendant dies before fixture cleanup if migration adds drain-on-exit. |
| PL82 | 17 | Shell interruption retains its lease after complete shutdown. | planned `internal/worktree/subshell_lifetime_test.go` (ShellSignal) | Completed shutdown incorrectly releases the signal-path assignment. |
| PL83 | 18 | Actual shift startup cancellation reaches a resistant adapter. | planned `internal/shift/lifetime_test.go` (ActualStartRace) | The live adapter survives if session cancellation still races cmd.Process. |
| PL84 | 18 | Shift deadline keeps precedence over simultaneous interruption. | planned `internal/shift/lifetime_test.go` (DeadlinePrecedence) | The result becomes interrupted/130 instead of its existing incomplete policy. |
| PL85 | 18, 30 | Recovery retains protected scratch in a clean worktree. | planned `internal/systemtest/process_lifetime_shift_test.go` (CleanProtectedRecovery) | An empty dirty set authorizes scratch removal or release. |
| PL86 | 30 | Explicit cleanup rechecks protection before apply. | planned `internal/worktree/protection_test.go` (ProtectedApply) | A user admitted after the plan loses its checkout if apply trusts the old fingerprint. |
| PL87 | 30 | Pool reuse refuses a protected checkout. | planned `internal/worktree/protection_test.go` (ProtectedAcquire) | Reset or clean changes the child's protected sentinel. |
| PL88 | 30 | Pool reclaim retains protected member users. | planned `internal/worktree/protection_test.go` (ProtectedReclaim) | Bulk key removal bypasses member protection. |
| PL89 | 33 | Linux parent death still kills the direct builder child. | planned `internal/runbinary/pdeathsig_linux_test.go` (ParentDeath) | The real child survives owner SIGKILL if Pdeathsig disappears. |
| PL90 | 34 | Native Darwin preserves its group policy. | planned `internal/runbinary/lifetime_darwin_test.go` (DarwinPolicy) | A resistant child survives if the native platform adapter loses Setpgid. |
| PL91 | 1, 35 | The final tree has no private production lifetime loop. | planned `internal/conformance/process_lifetime_test.go` (MigrationClosure) | A reintroduced private scanner or shift loop turns the graded-tree check red. |
| PL92 | 11, 36 | Gate preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (GateDifferential) | A changed normal code, output, or policy differs on the same canned fixture. |
| PL93 | 12, 36 | Builder preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (BuilderDifferential) | The same normal fixture exposes changed executable validation or output. |
| PL94 | 13, 36 | Bounds preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (BoundsDifferential) | The same fixture exposes changed start, exit, complete, canceled, or timeout projection. |
| PL95 | 14, 36 | Prep-release preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (PrepDifferential) | The same step fixture exposes changed attribution or ordering. |
| PL96 | 14, 36 | External preflight preserved fixtures match the baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (ExternalDifferential) | The same input and normal output expose a changed phase result. |
| PL97 | 14, 36 | Scanner preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (ScannerDifferential) | Valid, malformed, and exit-3 scanner fixtures expose changed finding policy. |
| PL98 | 15, 36 | Focused-test preserved fixtures match the baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (FocusedDifferential) | The same terminal, malformed, and incomplete JSON expose changed package policy. |
| PL99 | 15, 36 | Package-list preserved fixtures match the baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (ListDifferential) | The same package data and errors expose changed selection behavior. |
| PL100 | 16, 36 | Exec preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (ExecDifferential) | Canned stdin, stderr, exit, environment, and assignment fixtures expose changed command inputs. |
| PL101 | 17, 36 | Shell preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (ShellDifferential) | The same zero, nonzero, and signal fixtures expose an unapproved shell status repair. |
| PL102 | 18, 36 | Shift preserved fixtures match the production baseline. | planned `internal/systemtest/process_lifetime_preservation_test.go` (ShiftDifferential) | The same adapter result and deadline inputs expose changed evidence taxonomy. |
| PL103 | 37 | Migration preserves the declared scope exclusions. | review-owned, whole-spec fence review | Independent review rejects plain Capture migration, shell repair, and new escaped-descendant guarantees. |
| PL104 | 39 | A FIFO protection record prevents admission without blocking. | planned `internal/subprocess/protection/record_test.go` (FIFO) | A child starts or the record reader blocks if special-file checks follow open. |
| PL105 | 39 | A directory at a record path prevents admission. | planned `internal/subprocess/protection/record_test.go` (Directory) | A nonfile path silently becomes an empty record. |
| PL106 | 39 | Unsupported schema prevents admission. | planned `internal/subprocess/protection/record_test.go` (WrongSchema) | An unrecognized record gains current-schema authority. |
| PL107 | 39 | Duplicate JSON fields prevent admission. | planned `internal/subprocess/protection/record_test.go` (DuplicateField) | A later duplicate silently overrides the admission or binding field. |
| PL108 | 39 | Oversized protection data prevents admission. | planned `internal/subprocess/protection/record_test.go` (Oversized) | The reader accepts a record beyond the supplied control limit. |
| PL109 | 39 | Conflicting generations at one store prevent registration. | planned `internal/subprocess/protection/record_test.go` (GenerationConflict) | Deduplication hides a changed resource generation. |
| PL110 | 39 | Control-byte descriptor paths prevent launch. | planned `internal/subprocess/protection/record_test.go` (UnsafePath) | The child marker appears under an unsafe descriptor. |
| PL111 | 39 | Space and glob paths retain exact resource identity. | planned `internal/subprocess/protection/record_test.go` (AwkwardPath) | Splitting or expansion changes the store that receives protection. |
| PL112 | 39 | Repeated identical descriptors create one user per resource. | planned `internal/subprocess/protection/registration_test.go` (DescriptorDedup) | Two records or duplicate lock acquisition reveal inconsistent union handling. |
| PL113 | 7, 9 | A failed KILL remains unresolved until absence is proved. | planned `internal/subprocess/lifetime_test.go` (SignalFailure) | An ignored termination error falsely completes cleanup while the group remains present. |
| PL114 | 24 | Worktree environment transforms retain the resource vector. | planned `internal/worktree/exec_lifetime_test.go` (DescriptorEnvironment) | Nested CLI registration disappears if routing strips the reserved entry. |
| PL115 | 24 | Worktree --env cannot override the resource vector. | planned `internal/worktree/exec_lifetime_test.go` (DescriptorOverride) | The child sees a replacement store if user values win over owner propagation. |
| PL116 | 38 | Cleanup failure names each unresolved obligation once. | planned `internal/subprocess/lifetime_test.go` (IncompleteDiagnostic) | An empty or duplicated obligation projection hides the actual retained reason. |
| PL117 | 30 | Shift refuses rollback while resource users remain unresolved. | planned `internal/shift/lifetime_test.go` (ProtectedRollback) | Reset or clean changes evidence that the unresolved child can still modify. |
| PL118 | 40 | Ordinary complete resource users release their resources. | planned `internal/systemtest/process_lifetime_artifacts_test.go` (CompletedOwners) | A blanket retention implementation leaves cache, bundle, binary, test, or worktree protection behind. |
| PL119 | 14 | Abandoned release-stage cleanup retains protected evidence. | planned `internal/releaseevidence/protection_test.go` (ProtectedStageCleanup) | The stage disappears if cleanup treats its old owner marker as sufficient authority. |
| PL120 | 22 | OpenResource preserves users when a second owner joins. | planned `internal/subprocess/protection/registration_test.go` (ResourceJoin) | A fresh generation overwrites the first live user's protection without this join rule. |
| PL121 | 27 | Sweep retains a live prospective owner with no group users. | planned `internal/gate/prospectiveartifact/lifetime_test.go` (LiveUnusedOwner) | A second producer removes the first owner's checkout before its child registration. |
| PL122 | 25 | The focused run consumes selection Close failure. | planned `internal/testreport/lifetime_test.go` (SelectionCloseFailure) | Child zero becomes the operation result while the selected executable stays protected. |
| PL123 | 25 | The lane consumes its selection-release failure. | planned `internal/gate/lifetime_resources_test.go` (LaneSelectionCloseFailure) | A void release closure hides protected executable retention behind a lane pass. |
| PL124 | 3, 8 | Pending Start returns incomplete within the cancellation deadline. | planned `internal/subprocess/lifetime_test.go` (BlockedStart) | A held start worker blocks cancellation if control waits on the startup mutex. |
| PL125 | 3, 4 | Late successful Start receives shutdown and one waiter. | planned `internal/subprocess/lifetime_test.go` (LateStart) | Releasing the held worker leaves a live child if an incomplete return discards its pending cancellation. |

### Preservation fixtures

PL92–PL102 compare the production baseline and migration over the same inputs.
They ignore elapsed-time fields and use canned output.
Each row names one caller's differential predicate.
The compared input family contains the producer shapes below.

| Producer shape | Applicable callers | Expected preservation |
|---|---|---|
| Start failure | Every census caller | Existing attribution and no child |
| Complete zero exit | Every census caller | Complete output and current success policy |
| Complete exit seven | Every census caller | Current projection plus retained raw status |
| Valid complete JSON | Focused test, list, scanner | Current decode and domain validation |
| Malformed or incomplete JSON | Focused test, list, scanner | Current failure classification |
| Caller environment and direct files | Every caller that supplies these inputs | Exact command input and caller file lifetime |
| Surviving normal descendant | Preserve-mode callers | No added normal-exit signal and active resource protection |
| Surviving normal descendant | Gate and builder | Existing group drain before success |

The differential suite invokes each production entry.
It does not substitute a modeled lifetime result.
Defective unbounded cancellation does not become an expected result.
Existing cancellation witnesses remain until replacement omission evidence proves equivalent coverage.

### Edge inventory

The hostile-input attachment is the shell CLI checklist in `projects/benchkit.md`.
The process and resource producers add the concrete states below.
The new record reader consumes bounds' control-record limit.
It composes existing canonical-path policy instead of adding another resolver.

| Edge | Disposition and rows |
|---|---|
| Space, glob, Unicode paths | Preserve exact resource values, PL100, PL111 |
| Control bytes in paths | Refuse descriptor admission, PL110 |
| JSON-looking numeric identifiers | Preserve strict JSON types, PL106–PL107 |
| Absent versus empty cache | Keep distinct zero-result behavior, PL76–PL77 |
| FIFO, directory, symlink record | Refuse before open or mutation, PL41, PL104–PL105 |
| Unreadable, oversized, malformed data | Retain uncertainty, PL40, PL56, PL108 |
| Unknown field, schema, duplicate field | Refuse admission, PL43, PL106–PL107 |
| Duplicate users or conflicting generations | Retain uncertainty, PL45, PL109 |
| Repeated and reversed descriptors | Deduplicate and lock canonically, PL48, PL112 |
| Process exit before output EOF | Separate wait and decoder, PL14, PL64, PL81 |
| Cancellation before, during, after Start | Serialize publication and requests, PL2–PL5, PL83 |
| Failed KILL, EPERM, unknown group observation | Preserve unresolved facts, PL17–PL18, PL113 |
| Blocked stdin, output, or waiter | Return incomplete with required retention, PL11, PL13, PL69 |
| Nested groups and overlapping users | Observe the complete union, PL26, PL49–PL51 |
| Partial prepare, rollback, or identity publication | Prevent launch or retain uncertainty, PL21–PL25 |
| Seal after pending registration | Retain admitted users, PL27, PL46–PL47 |
| Owner death, reused identities, old records | Require current absence without signals, PL25, PL52–PL54, PL73 |
| Persistence before and after execution | Refuse launch or certification, PL21, PL24, PL72, PL78 |
| Clean shift worktree with protected scratch | Retain before recovery acts, PL85, PL117 |


Won't handle: Escaped unregistered descendants — the gate and builder retain their existing process-group contract.

Won't handle: Native Windows lifetime semantics — releasepreflight keeps its current non-Unix fallback outside this migration.

Won't handle: Bounded arbitrary caller I/O — exec retains caller-owned streams and explicit incomplete facts.

Won't handle: Guaranteed removal of unkillable OS tasks — the gate returns incomplete and retains affected evidence.

Won't handle: Shell nonzero normal-exit repair — shell keeps its current mapping and exposes the raw lifetime result.

Won't handle: Plain Capture migration — subprocess Capture callers without process-group ownership remain unchanged.

Won't handle: Gate order or green-authority redesign — the gate retains its current subject and evidence policy.

## Ownership fences

This union is the proposed implementation fence.
This spec-only pass writes none of its production or test paths.
Ticket design partitions it after independent spec review.
No writer receives all of internal or worktree.

The closure inventory uses the existing tickets.BoundFiles, anchors.ReferencingFiles, and canary.FixturePins owners.
Starting with this complete fence, their transitive requirements reach a fixed point.
Command-package and worktree edits co-name the five command registry files.
Profile, glossary, and data-policy edits co-name every referring anchor registry.
Individual fixture units below close their BASE, overlay, and mutation pins.
No fixture-family prefix grants unrelated writes.

The added holders retain their existing assertions, diagnostics, and mutation purpose.
Update only references or fixtures made stale by the approved ownership changes.
A holder that needs no content change remains co-owned and unchanged.
Do not weaken a check, remove an assertion, or change caller policy to close the fence.
PL103 retains whole-spec fence review, and the existing preflight closure rows grade the later ticket slices.

- `cmd/bench/process_resources.go`
- `internal/subprocess/`
- `internal/bounds/bounds.go`
- `internal/bounds/bounds_test.go`
- `internal/bounds/lifetime_test.go`
- `internal/gate/runner.go`
- `internal/gate/run_transaction.go`
- `internal/gate/engine.go`
- `internal/gate/lane.go`
- `internal/gate/phases.go`
- `internal/gate/lifetime_test.go`
- `internal/gate/lifetime_resources_test.go`
- `internal/gate/prospectiveartifact/`
- `internal/gocache/lock.go`
- `internal/gocache/clean.go`
- `internal/gocache/protection.go`
- `internal/gocache/protection_test.go`
- `internal/runbinary/runbinary.go`
- `internal/runbinary/protection.go`
- `internal/runbinary/sysprocattr_linux.go`
- `internal/runbinary/sysprocattr_darwin.go`
- `internal/runbinary/sysprocattr_other.go`
- `internal/runbinary/runbinary_test.go`
- `internal/runbinary/pdeathsig_linux_test.go`
- `internal/runbinary/lifetime_test.go`
- `internal/runbinary/lifetime_resources_test.go`
- `internal/runbinary/lifetime_darwin_test.go`
- `internal/env/kit_run.go`
- `internal/env/protection.go`
- `internal/env/kit_run_test.go`
- `internal/env/kit_run_lifetime_test.go`
- `internal/preprelease/preprelease.go`
- `internal/preprelease/lifetime_test.go`
- `internal/releasepreflight/command.go`
- `internal/releasepreflight/vulnerability.go`
- `internal/releasepreflight/lifetime_test.go`
- `internal/releaseevidence/release_evidence.go`
- `internal/releaseevidence/evidence_promotion.go`
- `internal/releaseevidence/protection.go`
- `internal/releaseevidence/protection_test.go`
- `internal/testreport/command.go`
- `internal/testreport/selection.go`
- `internal/testreport/cancel_test.go`
- `internal/testreport/selection_test.go`
- `internal/testreport/lifetime_test.go`
- `internal/worktree/exec.go`
- `internal/worktree/subshell.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/resume.go`
- `internal/worktree/pool_reclaim.go`
- `internal/worktree/protection.go`
- `internal/worktree/protection_test.go`
- `internal/worktree/exec_test.go`
- `internal/worktree/exec_lifetime_test.go`
- `internal/worktree/subshell_test.go`
- `internal/worktree/subshell_lifetime_test.go`
- `internal/shift/session.go`
- `internal/shift/loop.go`
- `internal/shift/recover.go`
- `internal/shift/lifetime_test.go`
- `internal/systemtest/process_lifetime_test.go`
- `internal/systemtest/process_lifetime_cache_test.go`
- `internal/systemtest/process_lifetime_artifacts_test.go`
- `internal/systemtest/process_lifetime_shift_test.go`
- `internal/systemtest/process_lifetime_preservation_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/conformance/process_lifetime_test.go`
- `internal/conformance/injected_ports_registry_test.go`
- `DATA_HANDLING.md`
- `CONTEXT.md`
- `projects/benchkit.md`
- `docs/adr/0019-one-owner-holds-the-prospective-artifact-bundle.md`
- `specs/process-lifetime/`
- `reviews/process-lifetime.md`

Closure holders:

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/data-handling-derivation/undocumented-passlist-var/`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift/`
- `tests/canary/guidance-prose-budgets/over-budget-skill/`
- `tests/canary/line-routing/line-binding-prose-drift/`
- `tests/canary/package-core-guard/bounds-classify-limit-restated/`
- `tests/canary/package-core-guard/bounds-discovery-window-unwrapped/`
- `tests/canary/package-core-guard/bounds-dot-import-package-alias/`
- `tests/canary/package-core-guard/bounds-dot-import-wait/`
- `tests/canary/package-core-guard/bounds-duplicate-owner/`
- `tests/canary/package-core-guard/bounds-intent-window-fixed/`
- `tests/canary/package-core-guard/bounds-multiple-dot-import-wait/`
- `tests/canary/package-core-guard/bounds-parenthesized-wait/`
- `tests/canary/package-core-guard/bounds-raw-elapsed-wait/`
- `tests/canary/package-core-guard/bounds-raw-injected-wait/`
- `tests/canary/package-core-guard/bounds-raw-wait-deadline/`
- `tests/canary/package-core-guard/bounds-raw-wait-duration/`
- `tests/canary/package-core-guard/bounds-read-limit-restated/`
- `tests/canary/package-core-guard/bounds-reassigned-wait-duration/`
- `tests/canary/package-core-guard/bounds-redeclared-wait-duration/`
- `tests/canary/package-core-guard/bounds-worktree-window-unwrapped/`
- `tests/canary/skill-description-budgets/budget-table-missing/`
- `tests/canary/skill-description-budgets/description-folded/`
- `tests/canary/skill-description-budgets/description-missing/`
- `tests/canary/skill-description-budgets/over-budget-command/`
- `tests/canary/skill-description-budgets/over-budget-description/`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading/`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner/`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing/`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership/`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route/`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary/`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term/`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts/`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary/`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term/`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term/`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary/`

Reviewer disposition: specification accepted at the retained checkpoint, ticket-graph review pending

The planning move changes only the C02 navigation link in the architecture index.
It moves the compiled map and topic folder as one unit.
The implementation fence grants no capture, roadmap, commitment, or unrelated spec edits.

## Out of scope

The source exclusions remain closed.
The shell status repair is a separate capability: estimated 3 edits, 2 gate runs.
A native Windows lifetime contract is a separate capability: estimated 12 edits, 4 gate runs.
An escaped-descendant supervisor is a separate capability: estimated 10 edits, 4 gate runs.
Plain Capture migration is a separate capability: estimated 8 edits, 3 gate runs.
These estimates claim no measured runtime cost.

This specification pass includes no benchmark, implementation, ticket file, roadmap reorder, or commitment change.

## Further notes

### Source verification

The author reread all structured map sources.
The production census still matches the pinned tip.
Platform adapters and fixture process setup do not add production caller sites.
The repository-wide sweep included dot-directories and excluded .git.
It found no moved-path reader beyond the map's source entries and the C02 navigation link.

The roadmap's fifteen-file wording counts a different surface from its eleven-copy occurrence.
Decision ticket 1 controls caller scope and includes bounds.
The roadmap's single-test proposal cannot prove migration and does not control this specification.
The profile's gate-only adapter claim must change with the approved owner.
The implementation edits that claim at its canonical source.

ADR0019's owner-death removal rule conflicts with its retained-descendant exception.
Decision ticket 6 resolves the gap through durable user proof.
This specification retains ambiguous legacy bundles.
Future production tests must qualify direct Close and the second-process sweep.
The prototype does not supply that qualification.

Read sources: Compiled map, six decision tickets, prototype evidence, and roadmap FT362

Read sources: Census callers, platform attributes, enclosing gate and lane transactions, cache, prospective bundle, binary, test directory, worktree, and shift recovery

Read sources: Testreport cancellation, worktree exec and shell, builder drain, prospective strict records, and Linux parent-death tests

Read enforcement: Coverage parser and citations, conformance registry and checks, injected-port registry, cancellation registration, gate entry, and project profile

Read closure: internal/tickets/registry_data.go, internal/preflight/closure.go, internal/preflight/gather.go, internal/anchors/references.go, and internal/canary/inventory.go

The author enumerated the fence through the existing accessors in a disposable read-only helper outside the repository.
The fixed point adds 12 registry files and 37 individually named fixture units.
No production test or gate ran for that enumeration.

Read data policy: DATA_HANDLING.md and SECURITY.md

Read dependency: specs/durable-file-replacement/spec.md from maintenance-primitives-spec
Accepted split SHA256: 05c79801d3746e5250fb23d7261f66827f971e6d26c0f9661ec7aaff48a0c2f8
Read graph-bearing SHA256: 4eacd9370ea59bf7003f897aa14e3e2702382b393b3db115727c094e61243333

Read excluded successor: specs/durable-caller-migration/spec.md from maintenance-primitives-spec
Accepted split SHA256: 32b04d2f6d1ce773170d1641eb3408a998d84656104c9fbf9702474f92c2ea66
Read graph-bearing SHA256: 2d606f650bb99d3c6daa3d040f504f432210d7e0cdf5f428b53ddbb1878be1ee

The source reads include ticket graph text added after independent split-spec acceptance.
The prerequisite retains its seventeen accepted rows and the successor retains its remaining thirty-five rows.
Build-entry prerequisite owns this process specification dependency and its implementation admission rule.

The data policy treats local records as mutable evidence and excludes same-user adversarial security guarantees.
The protection protocol shares that posture.
It stores identities, paths, statuses, and unresolved obligation identifiers locally.
It stores no environment vector, prompt, transcript, or credential.
DATA_HANDLING names the production record owner without duplicating its schema values.

### Reader and writer census

| Decision fact | Production readers or competing writers | Disposition |
|---|---|---|
| Group lifetime | Every entry in decision ticket 1 | PL1, PL91 and caller rows |
| Duration registry | Bounds, gate, builder, exec, shell, testreport | Selected caller policies, PL28, PL66–PL67 |
| Cache protection | Hold and exclusive clean, gate, lane, focused test | Shared lock before registration, PL55–PL57 |
| Prospective identity | Open, ReadPublished, Publish, Close, sweep, gate engine, lane | One strict owner, PL58–PL60, PL75 |
| Selected binary | Factory build-failure cleanup, Selection.Close, gate transaction, lane selection, focused test selection | PL36 and close-error composition |
| Private test directory | OpenKitTestRun, Entries, Close, phases, testreport, preflight | PL37–PL38 and nested registration |
| Worktree protection | Exec, shell, acquire, release, cleanup transaction, pool reclaim, shift teardown and recovery | PL82, PL85–PL88, PL117 |
| Release evidence | Preflight finalization, atomic promotion, abandoned-stage cleanup, prep-release phase | PL61–PL62 and protection owner |
| Descriptor environment | Subprocess launch, gate nested launches, worktree environment, shift adapter | PL49–PL50, PL114–PL115 |
| Moved compiled map | Its structured Sources and architecture index C02 link | Path move only |

The production reader sweep found Close results that current callers discard.
The migration consumes them wherever the operation owns cleanup or certification.
This includes focused selection cleanup and the lane's selection-release closure.
No completed child code can hide those failures.

Release evidence protection lives outside dist/preflight and temporary stage directories.
The repository common directory holds its store under bench-process-resources/release-evidence.
Its identity binds the repository and canonical dist/preflight target.
It protects phase-authored evidence before any phase starts.
FinalizeEvidence, promotion, and abandoned-stage cleanup consume its disposition before replacement or deletion.
The existing atomic exchange remains the publication owner.

### Pre-review proof checklist

Cited symbols: Existing callers and owners resolve, Lifetime and protection methods are planned interfaces

Import edges: Existing packages resolve through go list, subprocess and protection remain dependency leaves, bounds imports subprocess after implementation

Source-row clauses and occurrences: The source-to-row table below owns the relation

Promised field labels: Outcome and schema tables fix the proposed field vocabulary

Changed-function callers: The reader and writer census owns current scope, ticket design refreshes bench consumers before dispatch

Copy survival: PL91 rejects private mechanics, PL92–PL102 exercise each production caller

Rendered-shape readers: Existing cancellation attribution remains fixed, new cleanup facts have no previous rendered-shape consumers

Pin operators: Identity equality and current ESRCH remain exact, no numerical comparison pin changes

Entry reads: Resource identity resolves before protocol entry, faults enter through production interfaces

Derived expectations: Independent expectations require demonstrated omission reds before test removal

Consolidated rules: The policy table states each caller's retained rule and selected shared mechanics

Quantified obligations: The complete inherited/local union and every resource user remain explicit, PL26 and PL49–PL50 exercise the many case

Workflow-step writes: none, this outcome adds no phase or verb

The author runs no production process tests or gate during this spec-only stage.
Native Linux and macOS behavior remains a planned implementation obligation.
The injected-port registry must record every new fault port and its real-producer junction test.
The validator-table junction exercises the real resource-owner validators in a nested CLI.
No acceptance row names a test-only helper across a package seam.

### Source-to-row table

| Decision source clause | Rows or disposition |
|---|---|
| Ticket 1: current production callers include bounds | Policy table, PL28–PL38, PL61–PL65, PL79–PL102 |
| Ticket 2: shared mechanics with caller policy | PL1–PL19 and implementation interfaces |
| Ticket 2: serialize startup and own one Wait | PL2–PL5, PL83 |
| Ticket 2: bounds owns policy without a cycle | PL28, PL66–PL67, import review |
| Ticket 3: preserve normal behavior and caller contracts | PL19, PL28–PL35, PL79–PL84, PL92–PL102 |
| Ticket 3: Linux parent-death policy | PL89 |
| Ticket 3: missing grace and final cleanup window | PL12–PL14, PL66–PL67 |
| Ticket 4: bounded failure and retained resources | PL7, PL11, PL13, PL24–PL25, PL32–PL38, PL55–PL62, PL78, PL85 |
| Ticket 4: uncertainty and caller-owned streams | PL9, PL16–PL18, PL69 |
| Ticket 5: mechanism plus real caller tests | PL1, PL83, PL89–PL102 |
| Ticket 6: distinct process, streams, cleanup | PL5–PL11, PL64–PL65, PL69, PL79 |
| Ticket 6: owned pipe without early StdoutPipe Wait | PL10, PL14, PL64 |
| Ticket 6: prelaunch durable protection | PL20–PL25, PL70–PL72 |
| Ticket 6: cache shared/exclusive protocol | PL55–PL57, PL76–PL77 |
| Ticket 6: prospective direct and recovery protection | PL58–PL60, PL75 |
| Ticket 6: nested groups and enclosing owners | PL26, PL32–PL38, PL49–PL51, PL61–PL62, PL78 |
| Ticket 6: unknown users and no stale-record signals | PL25, PL52–PL54, PL73 |
| Ticket 6: shift recovery and actual startup race | PL83–PL88, PL117 |
| Ticket 6: native platform distinction | PL89–PL90 |
| Map exclusions | PL103 and Won't handle lines |

### Flagged additions

These decisions complete interfaces within the delegated architecture scope.
They add no capability beyond the map's resource and stream obligations.

- One durable launch record per resource prevents cross-user overwrite.
- Sealed admission closes the cleanup-to-launch race.
- Canonical lock order composes inherited resource unions.
- The reserved descriptor vector propagates nested registration.
- Worktree reuse and reclaim consume protected-user disposition.
- The next prospective schema retains ambiguous legacy proof.
- Release evidence keeps protection outside its exchanged output directories.
- Protection record publication consumes the complete delivery named in Build-entry prerequisite.

These decisions take PL20–PL27, PL39–PL60, PL70–PL78, PL85–PL88, PL104–PL115, and PL119–PL123.
Independent review must remove any unlisted acceptance expansion.

### Planning evidence

Document verification covers prose, coverage-map structure, source moves, links, and the authored diff.
Planned fixture citations remain uncited until implementation supplies those tests.
The frozen artifact hash and exact changed paths accompany the return to the coordinator.
The coordinator owns independent ticket review before dispatch.
No planning check supplies runtime or performance acceptance.

## Ticket approval

Author line: gpt-6.1-sol / high, one slicing pass and at most two repairs
Accepted specification SHA256: 19acc60a1ad76324de0b344052dcc9ddd16d887991baede8f361a5e22c8427a1
Independent ticket acceptance: gpt-6.1-sol / xhigh, Standards 0, Spec 0, Coverage 0, confidence 9; graph 3f8708bf2fd4c0e82045a41bc3f7c75afcadb490
No version 1 plan authorizes implementation.

| numbered ticket | Blocked by | delivered outcome |
|---|---|---|
| 1. 01-own-lifetime-through-bounds.md | none | Own bounded lifetime through the bounds runner |
| 2. 02-protect-nested-bounds-users.md | 01-own-lifetime-through-bounds.md | Protect cache users through nested bounds launches |
| 3. 03-migrate-gate-process-policy.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md | Preserve gate signals through the lifetime owner |
| 4. 04-retain-selected-build-binaries.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md | Retain selected binaries for unresolved users |
| 5. 05-protect-private-test-directories.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md | Retain private test directories through phase cleanup |
| 6. 06-refuse-clean-of-protected-cache.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md | Refuse cache clean after unresolved CLI exit |
| 7. 07-retain-prospective-gate-artifacts.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md | Retain prospective artifacts through gate and lane cleanup |
| 8. 08-preserve-release-step-evidence.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md | Retain release evidence after incomplete phase shutdown |
| 9. 09-bound-preflight-and-scanner-shutdown.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md | Bound preflight and scanner teardown without changing findings |
| 10. 10-preserve-focused-test-completion.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md | Preserve focused test decoding through bounded cancellation |
| 11. 11-preserve-worktree-exec-and-shell.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md, 10-preserve-focused-test-completion.md | Preserve worktree command policy with durable protection |
| 12. 12-retain-shift-scratch-through-recovery.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md, 10-preserve-focused-test-completion.md, 11-preserve-worktree-exec-and-shell.md | Retain shift scratch through startup cancellation and recovery |
| 13. 13-enforce-final-lifetime-closure.md | 01-own-lifetime-through-bounds.md, 02-protect-nested-bounds-users.md, 03-migrate-gate-process-policy.md, 04-retain-selected-build-binaries.md, 05-protect-private-test-directories.md, 06-refuse-clean-of-protected-cache.md, 07-retain-prospective-gate-artifacts.md, 08-preserve-release-step-evidence.md, 09-bound-preflight-and-scanner-shutdown.md, 10-preserve-focused-test-completion.md, 11-preserve-worktree-exec-and-shell.md, 12-retain-shift-scratch-through-recovery.md | Enforce the complete caller and resource lifetime census |
