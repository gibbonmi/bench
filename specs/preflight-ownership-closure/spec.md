# Preflight ownership closure

Status: staged

Roadmap: FT293

Decision source: named reviewed artifact decisions/architecture-planning.md, resolved tickets 1 and 5, reviewed 2026-10-06

Verification log: 1 graph iteration(s) to accept — spec and PF-CROSSING-1 accepted before slicing; graph accepted at c94b7f0e.
See assets/spec-review.md for review identity, source pins, and evidence limits.

## Problem

Authors repeatedly discover ownership obligations after ticket approval.
The current proposal reports fixture, package-registry, and anchor omissions through the same facts as the verdict.
Its registry requirement follows package membership.
Its proposal source cell drops the full reason chain.
These claims come from internal/preflight/closure.go:24 and internal/preflight/proposal.go:53.

A changed unexported helper can have callers outside the ticket fence.
A moved file can stop resolving before a charge reads its ticket.
Current writes-resolve also accepts committed D-status paths without an explicit deletion marker.
That existing behavior comes from internal/preflight/decision.go:378.

## Solution

Derive a complete mechanical proposal through the existing preflight owner.
Report each required path, its source, and its reason.
Require authors to amend the approved spec fence and ticket Writes together.
A proposal supplies no approval, write, or dispatch authority.

Verify supported premises against one committed source identity.
Unknown premises remain diagnostics or charge refusals.
Mechanical proof covers declared source facts and resolved static references.
It does not establish semantic adequacy.

## User stories

Line: gpt-6-astra / high, proposed implementation line only.
Implementation-line reason: grammar classification and source-bound helper analysis are the hardest chunks.
The existing fact/verdict split and production command fixtures bound their uncertainty.
Harder chunks: OC-C2 and OC-C3.
The current author line remains gpt-6.1-sol / high.
No binding or commitment changes here.

1. As a ticket author, I want closure proposals, so that I can request exact authority before dispatch.
2. As a reviewer, I want source reasons, so that I can distinguish proof from an inferred path.
3. As an author, I want transitive closure, so that a proposed registry also brings its fixture pins.
4. As a maintainer, I want grammar-specific bindings, so that unrelated package work retains only its real core obligations.
5. As an author, I want new command obligations, so that an unbound package cannot omit public inventories.
6. As an author, I want helper callers, so that changed signatures do not escape the fence.
7. As a maintainer, I want moved helper evidence, so that before-state callers remain visible.
8. As an author, I want deletion markers, so that a deleted source remains a valid planned ownership entry.
9. As a reviewer, I want counts and glob premises, so that stale enumeration cannot pass as current.
10. As an author, I want current headroom, so that a successor cannot pay this checkpoint's size debt.
11. As a reviewer, I want symbol and import premises, so that a proposed seam resolves from its actual caller.
12. As an author, I want sentinel and bound owners, so that a familiar name does not prove the wrong rule.
13. As a reviewer, I want test evidence premises, so that reused evidence names its real mutation record.
14. As a reviewer, I want posture and retirement censuses, so that every existing lifecycle consumer has a disposition.
15. As an author, I want predicate traceability, so that each declared decision maps to a coverage row.
16. As an author, I want an uncited-row view, so that I can slice before every row has a ticket.
17. As a reviewer, I want explicit ordering, so that proposed overlaps cannot silently enter the frontier.
18. As an author, I want safe operand handling, so that special files and source movement cannot yield partial green.
19. As a maintainer, I want old refusal authority, so that proposal work does not weaken charge preparation.
20. As a linked-repository author, I want precise empty states, so that absent kit registries do not become invented obligations.

## Implementation decisions

### Owners and value contracts

Keep Gather, immutable Facts, Decide, and the current proposal renderer as the production entry chain.
Introduce internal/preflight/ownership as the cohesive fact-collection and analysis module.
It does not import its parent.
Facts carries its immutable Snapshot.
Decide grades Snapshot requirements and diagnostics without I/O.
The existing closure-family descriptor remains the one ordering and applicability source.

Expose ownership.Collect(root string, source Source, manifest Manifest, tickets []tickets.Ticket) (Snapshot, error).
Source contains Base, Tip, and the canonical committed path/status projection.
It also carries SpecPath from the canonical selected-spec read.
A typed ownership.Error carries Stage, Subject, and Cause; Unwrap preserves the cause.
Source carries exact base and tip identities plus committed path statuses.

Snapshot carries Requirements, Premises, Diagnostics, Complete, and ChargeReadiness.
Complete means full supported premise certification, not reviewer approval.
Any manual claim keeps Complete false, including a declared future native witness.

Requirement carries ticket, entry, required path, kind, and all Provenance records.
Provenance carries owner path, line, rule identifier, reason, base, and tip.

It also carries source identity and tree or kit scope.
Tree evidence binds its actual committed owner bytes.
A compiled registry location must bind the independently verified kit source before it certifies proof.
Reuse the existing freshness verification and source identity from the caller's authenticated kit context.
Do not discover or execute another kit from an authored path.

An unavailable kit source or incomparable location stays explicitly diagnostic.
A built-in path spelling alone is no source proof.
Retain several reasons for one required path.
An error returns no usable complete snapshot.

Canonical tickets, anchors, and canary owners provide policy facts.
Add typed provenance projections beside their existing accessors.
Keep BoundFiles, ReferencingFiles, and FixturePins as compatibility projections of those same facts.
Do not copy their tables, BASE parser, Go scanner, or fixture materializer.

A fixed point walks newly required files until no new path or reason appears.
Cycles terminate by exact path and rule identity.
Sort paths within each canonical family while preserving ticket and family order.

### Authored claim metadata

Add one optional bench-planning-facts JSON fence to a spec.
Its parser lives in ownership and consumes the existing maps.FieldScan fenced-line classification.
It owns this schema only.
Do not add another Markdown scanner, sentence splitter, or general artifact reader.

Version 1 has authority authored-spec, completeness declared-complete or partial, and baseline source-base.
It carries scope, changes, premises, seams, decisions, and manual_claims arrays.
Authority states who authored the claims.
Completeness states what that author declared.
Neither field authenticates review, approval, or semantic completeness.

The existing approved spec and explicit ticket fields remain the authorization source.
Independent review must compare the declared inventory with all applicable source clauses.

Each item has a unique id and a spec location.
Scope is exact paths or segment-bounded prefixes.
Changes name kind, owner symbol, caller scope, and existing or planned state.
Supported change kinds are signature, fixture-helper, grammar, posture, and liveness.

Premises name kind, subject, expected value, and current or planned state.
Seams name path, test function, state, and supporting row IDs.
Decisions name source location, exact predicate bytes, and one row ID.
Manual claims name location, reason, and explicit applicability.
Applicability distinguishes a required current premise from future native acceptance evidence.
Neither form becomes supported or green.

Parse duplicate fences, duplicate keys, duplicate item IDs, unknown fields, unsupported versions, and invalid enum spelling as diagnostics.
No case folding changes a kind or symbol.
A missing fence reports incomplete premise inventory.
Legacy known closure facts can still be proposed with that diagnostic.
A build charge requiring unproved premises refuses.
No empty proposal means complete approval.

A declared complete inventory must cover every typed seam citation and decision item.
Coverage remains parsed by internal/coverage.
Expose its typed row/citation projection through that owner instead of parsing pipe tables again.
Plain prose symbols, non-Go claims, reflection, and unsupported citation forms remain manual claims.
Preflight never guesses them from word matches.

### Manual applicability and charge readiness

ownership owns ManualApplicability with explicit CurrentPremise and NativeFuture values.
ManualClaim carries ID, Location, Reason, Applicability, and an optional NativeWitness.
NativeWitness carries Row, Chunk, Ticket, Verification, PhaseOwner, and Observation.
The tagged JSON decoder rejects native fields on current-premise items and requires every native field on native-future items.

Missing, unknown, or contradictory values remain metadata diagnostics.
No authored default chooses the less restrictive form.
Every manual item remains in the premise inventory and diagnostic projection.

A native-future item identifies one exact native acceptance row and that row's spec location.
Its row, ticket, chunk, and verification must agree with the canonical completion plan at Source.Tip.
Use reviewrecord.ReadPlan and the already parsed ticket Covers to validate those relationships.
The selected SpecPath comes from the canonical spec read, never an authored JSON path.

The named ticket owns the row exactly once and belongs to the named chunk.
The verification belongs to that chunk.
A version 2 verification also names that exact ticket.
For version 1, derive ticket ownership from Covers, not the verification command's words.

phase_owner pins the existing native implementation procedure through the existing owner-location contract.
The row must declare a planned native witness, rather than a supported current Go citation or required current-source premise.
The reviewed graph explicitly binds that witness to its future verification route.
A matching lexical phrase, authored applicability field, or inherited path does not establish this classification.

Observation is always unobserved in this schema.
A native-future item remains acceptance-required manual work until its native verification occurs after the owning implementation.
Existing completion records and independent review establish that later evidence.
The metadata does not establish execution authenticity, review approval, or a completed witness.

Invalid row ownership, an unavailable phase owner, or a missing verification remains diagnostic and refuses charge readiness.
Unknown current claims cannot be renamed native-future to gain readiness.
Missing metadata, partial inventory, and omitted supported citations retain their refusal.
No Go seam substitutes for a native witness.

ChargeReadiness carries InventoryComplete, CurrentComplete, RequiredCurrent, and PendingNative.
InventoryComplete means the declared complete inventory passed structural cross-checks, including all manual locations.
It does not mean independent review occurred or every item was mechanically proved.
RequiredCurrent contains every required current-source premise ID, whether supported or unsupported.
PendingNative contains validated native-future descriptors and their diagnostics.

CurrentComplete means all required current-source premises have complete supported proof.
Any current-premise manual item or unknown required current input makes it false.
Charge readiness requires both InventoryComplete and CurrentComplete, alongside every existing applicable readiness check.
Snapshot.Complete remains false while either manual form remains.

Facts carries the registered evidencecmd.Kind as Operation.
The dispatcher passes the actual selected kind through preparation before pure Decide runs.
Build KindPrepareEvidence and build KindCurrentEvidence enforce the current-premise readiness conjunction.
The current-evidence route uses its artifact's canonical build selection, not an authored applicability switch.

Ordinary verdict, proposal, plan-only, breakdown, and review operations retain their accepted applicability.
The global completeness projection remains independent of this charge-specific decision.
buildChargePack continues to refuse every red verdict; no publisher bypass or red exception is added.
Current checkout, assignment, source, movement, seal, commitment, and evidence refusals retain their order and force.

### Exact manifest shape

Manifest uses versioned JSON objects with no implicit defaults.
Every top-level array is required, even when empty.
Location is {path, line, needle, source}, with a positive one-based line and exact nonempty needle bytes.

Source is spec for an item location and base or tip for an owner location.
An item location points into this spec; an owner location points into its declared pinned source.
Line and glob subjects also declare base or tip; current headroom and import reachability require tip.
The needle must occur exactly once at the declared line.
Moved or ambiguous needles are diagnostics, never a nearby-match repair.

| item | required fields and allowed values |
|---|---|
| scope | id, location, paths: an array of exact repository paths or trailing-slash prefixes |
| changes | id, location, kind, owner, caller_scope, state, destination, dispositions |
| change owner | path, package: exact import path, symbol: exact consumers query, kind: exact consumers declaration kind, signature: normalized declaration syntax |
| caller_scope | production-and-tests; another context is diagnostic |
| change state | existing or planned; destination is null except an explicit relocation path |
| dispositions | an array of {symbol, action}; action is retain, migrate, or remove for each declared retirement consumer |
| premises | id, location, kind, state, subject, expected |
| premise state | current or planned; planned values never certify observed current proof |
| seams | id, location, path, package: exact import path, function, state, rows; state is existing or planned |
| decisions | id, location, source: owner location, predicate: exact bytes, row: one existing row ID |
| manual_claims | id, location, reason: nonempty unsupported operation or semantic question, applicability: current-premise or native-future |
| current-premise manual claim | no additional fields; unsupported current required proof remains a charge refusal |
| native-future manual claim | row, chunk, ticket, verification, phase_owner: owner location, observation: unobserved |

Every symbol owner and calling seam also names its exact package import path.
Scope prefixes use the existing segment-boundary rule.
Empty or contradictory scope, unknown change kind, and undeclared destination are diagnostics.
A declared grammar change names registration and inventory owners in its subject.
Those paths are checked against the canonical binding owner.

Premise subject and expected are tagged objects, not arbitrary prose.
The supported kind inventory follows.
The collector refuses extra or missing fields for each kind.
A historical count at base never proves current headroom at tip.

| kind | subject | expected |
|---|---|---|
| lines | path, source: base or tip | count: nonnegative physical newline count |
| glob-members | pattern, scope, source: base or tip | members: complete sorted exact paths |
| headroom | directory, files: creation/deletion plan | maximum, current, net: integers; canonical grant and budget owner locations |
| symbol | owner identity | kind, signature: exact declaration values |
| import | caller package, target package, owner identity | reachable: true |
| anchor-kind | owner location, kind | supported: true |
| sentinel or error | owner identity, producer identities | value: exact static construction text |
| bound | owner identity | value: exact Go constant value |
| reused-test | path, function, record_spec, source_commit, evidence_commit, chunk, verification, mutation | red_needle: recorded behavioral failure marker, restored: true |
| reused-owner | owner location | required: true |
| retirement | owner identity, consumers | dispositions: complete symbol/action pairs |
| posture | owner identity, callers, writer_probe: reused-test subject | callers: complete static producer identities |
| fence-exclusion | excluded prefix, exported owner identity | reachable: true |
| printed-field | field, producer identities | producers: complete static owner identities |

An owner identity contains path, package, symbol, kind, and signature.
Const identities have an empty signature and their exact static value.
Normalize declaration syntax through the existing consumers AST reader and Go formatter.
The typed current resolver confirms declaration origin and reachability.
Normalization proves structural identity only.
Opaque expressions, unspecified contexts, or incomplete lists remain diagnostics.

For reused evidence, call reviewrecord.CheckTrees with the exact source and evidence commits, requested spec, chunk, and complete false.
Read the same record through reviewrecord.ReadTree and select the named verification.
Its source digest, mutation, bit outcome, nonzero probe exit, and pass restoration must agree with the canonical validator.
The validated native excerpt must contain red_needle as a behavioral assertion marker.
An excerpt showing compilation failure without that marker cannot satisfy this premise.

The existing checkpoint authority and independent review own whether the recorded witness is genuine and adequate.
An excerpt digest authenticates no external execution.
No new record operation or artifact reader is added.

A partial example follows.
It declares future work without claiming a complete inventory or existing evidence.

```json
{"version":1,"authority":"authored-spec","completeness":"partial","baseline":"source-base","scope":[],"changes":[],"premises":[],"seams":[],"decisions":[],"manual_claims":[{"id":"semantic-rule","location":{"path":"specs/example/spec.md","line":12,"needle":"Preserve the consumer policy.","source":"spec"},"reason":"semantic equivalence requires review","applicability":"current-premise"}]}
```

Extend FieldLine with FenceRole and FenceLabel projections from FieldScan's existing fence state.
FenceRole is a uint8 enum with FenceNone zero, FenceOpen, and FenceClose.
FenceLabel is the existing marker's exact trailing label.
Keep Fenced and all existing field behavior unchanged.
Existing parity cases preserve every legacy field; independent role cases grade the added projection.

The metadata owner consumes only an open role with exact label bench-planning-facts.
A label inside another code fence cannot authenticate a metadata opener.

An unclosed or nested metadata fence diagnoses incomplete collection.
Its payload uses bounds.ControlRecordLimit and the existing classified byte reader.
JSON token decoding rejects repeated keys; that is schema parsing, not another Markdown splitter.

### Grammar and core bindings

Extend the existing tickets binding registry with typed core and grammar obligations.
Core obligations retain their current rule owners.
Terminal and toon seed bindings remain Core path obligations.
The existing five-file commandRegistries set becomes the Grammar obligation, rather than an unconditional package obligation.
The cmd/bench seed remains registered under that typed Grammar binding.
Kinds Core and Grammar have exact spelling; source selectors are also case sensitive.

Grammar obligations follow declarations that produce argument shape, help rows, public command inventories, or routing.
Exact identifier case is significant.
The registry owns its selectors and obligated files.
The preflight collector resolves those selectors through actual static producers.

Ordinary producers are the resolved usage.Grammar values passed to usage.Parse.
Their argument-shape and help data dependencies are walked through typed references.
Custom producer roots are evidencecmd.flagTable, evidencecmd.operations, recordcmd.forms, and cmd/bench.commandRegistry, in their actual owning files.
Keep those selectors in the existing tickets registry, beside the obligated path set.
An unresolved argument or producer dependency stays diagnostic.

A changed grammar declaration requires its grammar source and each independent public inventory expectation.
A body-only change in the same package adds no grammar obligation.

BoundFiles(path) remains the legacy all-possible-path projection of the typed registry.
The effect-aware collector applies its Core and Grammar applicability and replaces the coarse probe result.
There is one path table and no parallel policy registry.
An unchanged grammar dependency does not become a changed form merely because its caller changed.

A planned new public form names its actual registration and expectation owners before its implementation exists.
The canonical command registry supplies the shared files even when the new package has no prefix binding.
An ambiguous producer, missing registration owner, or unsupported dynamic grammar prevents complete classification.
Do not bless a form from a function named Command alone.

Retain executable, anchor, fixture, and system-root closure at every fixed-point step.
A new system test receives the existing BENCH_KIT obligation through newSystemTest.
Missing kit inventory in a linked repository is an empty fact only where the canonical owner specifies that result.
An unreadable present registry is a diagnostic.

Current discoverFixtures collapses directory read errors into absentHarnessMessage.
The canonical pin projection first classifies the present inventory, preserving absence while refusing unreadable or special directories.
This is a deliberate stricter collection boundary, not a claim that the current owner already distinguishes those errors.

### Helper and owner closure

Reuse internal/consumers for resolved Go references, including unexported symbols.
Expose consumers.Load(root string) ([]*Package, error) beside loadPackages.
The ownership Snapshot stores the immutable resolved identities, not mutable Go package objects.

It reuses the existing loadPackages, Resolve, and Rows owners.
It accepts the current clean committed tip and default build context.
It returns complete typed packages or an error.
No preflight shell call reparses bench consumers output.

Declared signature and fixture-helper changes enumerate all static callers before proposing writes.
The caller scope includes production and test packages.
A new helper has planned state and no invented existing caller set.

For a deletion or relocation, compare the pinned before declaration with the actual tip declaration or explicit destination.
Reuse the existing consumers before-state AST machinery.
Do not import a test-only helper across packages.
Unsupported build contexts or unresolved dynamic references remain diagnostics.

Declared reused owners add their exact source paths.
A liveness change names each retirement consumer and its retained or removed disposition.
A posture change names its static call-site census and writer-probe evidence.
Each exported fence exclusion names the symbol that crosses it.
A printed-field premise names every static producer.
A missing producer or retirement member is not a complete census.

### Deleted entries and authorization

Add tickets.ParseWrite as the canonical typed write-entry parser.
It returns path and Existing, New, or Deleted kind.
Existing WritesPath projects that result for its current callers.
New and deleted markers cannot coexist.

The exact deletion spelling is path (deleted).
A deletion names an exact tracked before-state path.
A deleted directory shorthand is refused.

A deletion entry resolves against the pinned source base.
It can be present before implementation or absent after the move.
Its absence at the tip is not a resolution error.
A nonexistent before path cannot be rescued by the marker.

Preserve compatibility for the existing exact committed D-status case.
Strip markers once through the canonical parser for authorization, closure, union, and overlap checks.
No marker expands the spec fence.

Build dispatch requires every proposed path to appear in both explicit ownership sources.
The spec fence must equal the ticket writes union under the existing phase and review-pickup exclusions.
Every new overlap requires an explicit blocker path.
The reviewer still approves the graph.
No proposal mutates the spec, tickets, registry, or assignment.

### Supported premises

Every fact and diagnostic prints its source tip.
Count and glob checks enumerate tracked members, including tracked dot paths.
Counts derive from that member set.
Globs use documented filepath.Match syntax, never shell expansion.
Unsupported recursive syntax, escaping scope, or incomplete member lists are diagnostics.

Current headroom uses internal/structure's budget, grant, file-count, and line-count owners.
Expose a typed sizing projection through that owner.
Do not compute a second limit or line-count rule.
An authored net file plan counts declared creations and exact deletions.
Future net-line estimates remain estimates, not observed fit.
An over-budget source cannot gain lines without a same-checkpoint cohesive move or existing approved grant.

Supported symbol checks resolve exact declaration, kind, owner file, and signature.
Import checks resolve the target from the named calling seam and reject private cross-package reachability or an import cycle.
Supported anchor-kind checks use the actual anchors kind inventory.

Sentinel and error premises identify the exact declaration and constant construction where statically resolvable.
A failure-message claim still needs its reachable producer witness.
A bound premise identifies its owning constant or constructor and expected static value.
An unknown expression stays unknown.

A reused-test premise names the test function, mutation, and a source-bound existing verification record.
Read that record through reviewrecord's validator.
Check its subject identity, recorded behavioral failure marker, and successful restoration.
A valid record proves the recorded evidence shape, not the mutation's semantic adequacy.
No newly authored text authenticates an observed red.
Unavailable or incomparable evidence refuses that premise.

### Source placement and headroom

The pinned newline counts are 475 in decision.go, 431 in gather.go, and 599 in coverage.go.
Inventory.go has 415; citations.go has 397; structure.go has 396.
These are source observations at a9c395e77fec36d1f60f6d057c654aecf5720e24.
Neither structure.budgets nor structure-accept grants these files extra lines.

Growth reds a file only when its tip exceeds its applicable limit and its before count.
Directory crowding is a separate whole-tree observation.
Existing debt alone does not impose a new file-growth refusal.

Place new analysis and multi-case tests in the ownership child.
Add no file to the existing 44-file preflight root.
Move ticket probing and system-tag collection from gather.go into existing closure.go and system_tag.go.
Move pure write authorization and resolution checks from decision.go into existing fence_writes.go.
That move replaces the gathered closure machinery; it does not duplicate its policy.

Move the existing pure coverage parser/model into rows.go and export projections through projection.go.
Coverage has nine top-level source files; those two create eleven.

Move FixturePins and its pin enumeration into pins.go, with one pins_test.go.
Canary has ten top-level source files; those two create twelve.
Reuse its existing BASE reader in inventory.go and mutation owner in mutation.go.

The typed consumers loader fits loader.go's current 60 lines without another root file.
Structure's projection belongs in its existing 34-line facts.go and calls its shared sizing helpers.
Existing source-count and row-list witnesses change only for the documented extraction or added row.
No budget or accept-list change is authorized.
Implementation must remeasure the composed checkpoint through the canonical Growth and crowding owners.
An estimate cannot prove future fit.

### CLI grammar and responses

Retain:
bench preflight build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>.

Add:
bench preflight build <slug> --breakdown [--base <commit>] [--source-tip <commit>].

The existing evidencecmd operation registry owns grammar, help, root inventory, and mutual exclusions.
Breakdown cannot combine with charge, plan-only, proposal, or evidence selectors.
Missing required values, extra operands, and unknown flags retain exit 2.
Help spellings retain the current shared grammar behavior.

Proposal output retains writes_proposal[N]{path,source,fence} and ordering[N]{ticket,other,required}.
Add closure[N]{ticket,path,requires} and closure_evidence[N]{ticket,path,kind,source,line,reason}.
Path in closure is the originating write entry.
Requires is its required ownership path.
Evidence names each canonical owner and exact source location.

Add source[1]{base,tip} and premises[N]{id,state,source,reason}.
Green premise state is complete supported proof only.
Unknown state never becomes green.
Long output uses existing response bounds and complete spills.

Breakdown prints uncited[N]{row,story,seam}, closure, and premises from the same gathered facts.
It tolerates rows-owned and repairable closure/fence reds only for this read-only projection.
It renders those diagnostics and grants no frontier authority.
Rows-membership, malformed grammar, invalid source, and unsupported operand states still refuse.
Absent-manifest legacy inputs remain explicitly incomplete diagnostics in the authoring views.

The enhanced charge requires its complete reviewed claim inventory and complete required-current-premise proof.
Unobserved native-future entries remain diagnostics without requiring execution before their own implementation charge.
Absence cannot satisfy the OC-C3 premise check or authorize dispatch.

A present empty tickets directory reports every unowned row.
An absent directory permits the authoring projection while retaining ordinary build applicability.
Neither route changes the existing red rows-owned policy.

Keep current assignment, clean-checkout, no-follow read, source-pin, movement retry, and trust checks.
Dirty source plus missing selected ticket reports the existing checkout refusal first.
Persistent movement publishes no proposal or charge.
No new metadata or path can authenticate the executable or an evidence artifact.
The current sealed binary and current-action binding remain required.

### Concrete request scenarios

Use the actual preflight dispatcher with committed fixtures and explicit source pairs.
One fixture holds internal/example/helper.go with an unexported mergeFixture and three same-package test callers.
Its approved request declares a signature change and owns only the helper and the first caller.
The proposal must name the other two exact caller files with their declaration source reasons.
Sampling the first caller fails the third-file sentinel.

Relocate that declaration to another file in the same package at the second fixture commit.
The request owns the original path with its deleted marker and the destination with its new marker.
Before-state references remain in closure after the original file disappears.
Changing the deletion to a never-tracked path refuses; no destination can repair that authority.

A separate fixture edits a body under the terminal seed without changing grammar.
Its existing Core adopt/gate obligations remain.
A body under an ordinary command package adds no Grammar inventory requirement.
A declared new preflight selector requires its operation source and all five independent command registry holders.
Dropping the root help expectation must fail even though preflight had no legacy prefix binding.

The exact request forms are these future commands.
BASE and TIP mean full IDs returned by the committed fixture, not ambient branches.

```text
bench preflight build example --propose-writes --ticket 02.md --base BASE --source-tip TIP
bench preflight build example --breakdown --base BASE --source-tip TIP
bench preflight build example --charge --ticket 02.md --base BASE --source-tip TIP
```

The first response is a proposal only.
The third refuses until the spec and ticket both own every proposed path and the required premises complete.
Its publication sentinel is absence of a prepared evidence handle.
The second can report an unowned row while the ordinary build still refuses that row.

## Implementation chunks

These are complete planned green consumer outcomes.
The ticket graph is independently accepted before implementation.
The first real proposal consumer proves the ownership seam in OC-C1.
Each successor waits for that independent chunk review.
No later consumer starts from an unreviewed provider contract.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| OC-C1 / 01-propose-canonical-closure.md | Typed canonical closure through the real proposal command | OC1–OC13, OC55–OC60, OC71–OC81, OC83–OC85, OC90–OC92 | Ownership proposal and existing command fixtures | no |
| OC-C2 / 02-classify-grammar-effects.md, 03-propose-helper-and-owner-closure.md | Grammar, helper, and deletion closure through real author requests | OC14–OC31, OC82 | Typed loader and pinned-repository proposal fixtures | yes |
| OC-C3 / 04-check-sizing-premises.md, 05-check-symbol-and-rule-premises.md, 06-check-recorded-evidence-premises.md, 07-refuse-unproved-source-charges.md | Source premises refuse unproved build charges | OC32–OC54, OC87–OC89, OC93–OC95 | Premise fixtures through actual Gather and charge | yes |
| OC-C4 / 08-project-breakdown-and-author-authority.md | Read-only breakdown and explicit author workflow | OC61–OC70, OC86 | Root grammar, breakdown, and guidance checks | no |

OC-C1 preserves legacy operations while supplying provenance through one real consumer.
OC-C2 does not require FT375 or a new record writer.
OC-C3 consumes the established source and typed closure values.
OC-C4 consumes all facts without weakening charge readiness.
The final union audit follows all migrations.

Ticket slicing supplies this spec's complete claim manifest before implementation dispatch, under independent graph review.
It records supported current-source premises, exact planned seam identities, and every unsupported native acceptance obligation.
Independent graph review validates native-future applicability and its exact completion-plan ownership.

This prevents OC-C3's own newly enforced premise check from making its continuation impossible.
Earlier chunks refresh changed current premises through the existing enabling-plan authority before OC-C3 activates its required charge check.
Historical source-census counts retain their explicit baseline identity.
OC-C1 and OC-C2 do not promise the successor premise collectors or global premise certification.

## Testing decisions

Prefer the existing production command fixture and the actual proposal/charge entries.
TestChargeManualApplicability drives the actual charge and current-evidence routes with committed metadata and completion-plan fixtures.
Its ready case has no future native execution record and still retains that item's diagnostic and globally incomplete snapshot.
Its competing refusal case adds an unsupported required current premise and asserts that no evidence handle is published.
Invalid native ownership and current claims relabeled as future also refuse.

The same fixture checks review and ordinary continuation without inventing a manual-omission decision.
Omitting the current-premise refusal must fail its publication sentinel.
Requiring a future witness record before charge must fail the ready case.
All such runtime and mutation witnesses remain future implementation work.

Pure Decide tables prove classification only.
They cannot replace collector or executable-route witnesses.
Use a few multi-case fixtures whose independent sentinels fail for different omissions.
Do not add one implementation-mirroring test per coverage row.

The new ownership package receives local-substitutable Git and source fixtures.
Its external tests can import the parent entry and the existing preflighttest fixture.
Production ownership imports no parent.
The consumers typed loader gets one actual-tree witness and pure typed-package cases.
Existing loader refusal and static-reference limits remain binding.

Planned verification:

```text
bench test --package ./internal/preflight/...
bench test --package ./internal/tickets
bench test --package ./internal/consumers
bench test --package ./internal/canary
bench test --package ./internal/anchors
bench test --package ./internal/coverage
bench test --package ./internal/structure
bench test --package ./cmd/bench
bench test --check ticket-grammar
bench test --check docs-currency-workflow
```

Each mutation runs through bench probe, restores, and reruns the same command.
Source fixtures bind actual before and after commits.
Canned typed facts pin output order without elapsed-time comparisons.
Existing baseline fixtures protect preserved refusal precedence and source identity.

### Seam diagram

    explicit source pair + authored spec + ticket
        -> Gather + canonical owners + consumers typed loader
        -> ownership.Snapshot
        -> Decide
        -> proposal / charge / breakdown response
        <- tests drive the same command and inspect independent file sentinels

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| OC1 | 1 | A missing fixture pin appears as a required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | Omitting the BASE pin cannot leave an empty proposal. |
| OC2 | 1 | A proposal leaves spec and ticket bytes unchanged. | planned internal/preflight/ownership/proposal_test.go (TestProposalReadOnly), actual preflight command | A before/after byte sentinel detects accidental authorization. |
| OC3 | 3 | A required registry's own fixture enters the fixed point. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | A second-hop pin catches one-pass closure. |
| OC4 | 3 | A closure cycle terminates without dropping its required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | The cyclical fixture detects incomplete or infinite traversal. |
| OC5 | 2 | Two independent reasons for one path survive in evidence. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Path deduplication cannot erase the second rule. |
| OC6 | 2 | An anchor reason names its actual holder and source line. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | A synthetic entry label cannot satisfy the location sentinel. |
| OC7 | 2 | A fixture include reason names the actual included BASE source. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Naming only a fixture directory misses its derivation. |
| OC8 | 2 | A registry reason names the canonical rule owner. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | A copied registry location fails the owner sentinel. |
| OC9 | 2 | A proposal prints the exact resolved source tip. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | An ambient or truncated tip fails exact identity. |
| OC10 | 20 | An absent linked-kit anchor directory adds no invented holder. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The inventoryless root catches unconditional kit closure. |
| OC11 | 20 | An absent linked-kit fixture inventory adds no invented fixture. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The canonical empty inventory stays observable. |
| OC12 | 2 | A complete empty proposal prints typed zero-row tables. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | Silence and a legacy prose-only response fail shape assertions. |
| OC13 | 18 | An incomplete collector never reports complete proof. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A refused registry source catches green-by-omission. |
| OC14 | 4 | A body-only change adds no grammar obligation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The unrelated helper edit detects whole-package binding. |
| OC15 | 4 | A core rule obligation survives a body-only change. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | Dropping all obligations cannot satisfy precision. |
| OC16 | 4 | A changed grammar producer requires its grammar source. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The actual usage declaration catches missed producer changes. |
| OC17 | 5 | A new unbound form requires the root help inventory expectation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | An unbound package fixture catches missing shared inventories. |
| OC18 | 5 | An unresolved registration producer remains diagnostic. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | A guessed Command-name binding cannot complete. |
| OC19 | 4 | Case-sensitive binding selectors remain distinct. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | A Core versus core fixture catches silent case folding. |
| OC20 | 6 | A changed unexported signature proposes each static test caller. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The third caller outside seed files catches sampled sweeps. |
| OC21 | 6 | A fixture-helper change proposes its external call-site file. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | A reused fixture outside the owner package catches local-only analysis. |
| OC22 | 7 | A relocation retains the before-state helper's caller proof. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The moved declaration catches after-only enumeration. |
| OC23 | 7 | A deleted helper retains its before-state caller proof. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | An absent tip declaration cannot erase its former users. |
| OC24 | 8 | A deleted marker resolves an exact tracked before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The source base proves ownership despite absence at the tip. |
| OC25 | 8 | A nonexistent before path refuses the deleted marker. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | A fabricated name cannot gain authority. |
| OC26 | 8 | A planned deleted marker accepts a still-present before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The preimplementation fixture prevents an impossible first checkpoint. |
| OC27 | 8 | Combined new and deleted markers refuse parsing. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | Conflicting authority kinds cannot collapse into New. |
| OC28 | 14 | A reused owner adds its exact rule source to closure. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | An omitted home-policy owner fails the rule sentinel. |
| OC29 | 14 | A liveness premise reports an omitted retirement consumer. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The second lifecycle consumer defeats a partial census. |
| OC30 | 14 | A posture premise reports an omitted static writer caller. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The independent writer site catches a reader-only census. |
| OC31 | 14 | A printed-field premise reports an omitted producer. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The second renderer branch catches one-producer evidence. |
| OC32 | 9 | A line-count premise reports its changed current value. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | One added physical line defeats a stale count. |
| OC33 | 9 | A glob-member premise reports its omitted tracked dot member. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | A hidden tracked member defeats visible-only enumeration. |
| OC34 | 9 | A glob premise retains literal metacharacter path identity. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Shell expansion cannot supply the expected exact member. |
| OC35 | 9 | An unsupported recursive glob remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | A double-star input cannot become an empty complete set. |
| OC36 | 10 | A current file-headroom premise uses the canonical limit. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Changing the owner grant defeats a copied cap. |
| OC37 | 10 | A net file plan counts its exact declared deletion. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Ignoring the deletion changes the owner-derived fit. |
| OC38 | 10 | Future net-line estimates remain unobserved estimates. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | An unbuilt target cannot yield observed-fit state. |
| OC39 | 11 | A named symbol resolves to its declared owner file. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-named symbol in another owner fails identity. |
| OC40 | 11 | A declared signature mismatch refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A changed parameter defeats name-only matching. |
| OC41 | 11 | A calling seam cannot import an unexported cross-package helper. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | The real loader catches private reachability. |
| OC42 | 11 | A cycle in the proposed import edge refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | The actual package dependency graph catches cyclic composition. |
| OC43 | 11 | A supported anchor kind resolves from the canonical inventory. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An unsupported step placement cannot use a section kind. |
| OC44 | 12 | A sentinel premise resolves its exact owner declaration. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A familiar Err name elsewhere cannot certify the rule. |
| OC45 | 12 | An unresolved failure expression stays diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An opaque dynamic Error method cannot certify exact text. |
| OC46 | 12 | A bound premise resolves its owning static value. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-valued local literal cannot substitute for the owner. |
| OC47 | 12 | An opaque bound expression remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An unknown runtime input cannot become an inferred limit. |
| OC48 | 13 | A reused-test premise requires its source-bound evidence record. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A prose claim without a valid record cannot certify an observed red. |
| OC49 | 13 | A mismatched evidence subject refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A real red for another subject cannot satisfy reuse. |
| OC50 | 13 | A failed restoration refuses reused mutation evidence. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | An un-restored behavioral probe cannot supply valid evidence. |
| OC51 | 15 | A declared predicate names exactly one existing coverage row. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A nonexistent or multirow predicate mapping is diagnostic. |
| OC52 | 15 | Changed exact predicate bytes invalidate their row premise. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A retained row ID cannot hide changed requirement bytes. |
| OC53 | 18 | An unsupported dynamic caller claim remains diagnostic. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | Reflection cannot be silently reported as an empty complete census. |
| OC54 | 19 | An unproved required premise prevents charge publication. | planned internal/preflight/ownership/premises_test.go (TestChargePremises), actual preflight command | A prepared-handle absence sentinel catches dispatch on unknown facts. |
| OC55 | 18 | A duplicate metadata fence refuses complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Selecting the first duplicate cannot pass. |
| OC56 | 18 | Duplicate JSON keys refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Last-key-wins decoding cannot pass. |
| OC57 | 18 | Unknown metadata fields refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Ignored undeclared promises cannot pass. |
| OC58 | 18 | Duplicate item IDs refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Overwriting one claim cannot produce full evidence. |
| OC59 | 18 | An unsupported metadata version remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Guessing a schema cannot certify completion. |
| OC60 | 18 | An invalid metadata enum remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Case-folded or unknown kind values cannot pass. |
| OC61 | 16 | Breakdown lists the actual uncited row's typed seam. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | A prose token cannot replace the coverage parser's row. |
| OC62 | 16 | Breakdown keeps the ordinary rows-owned refusal unchanged. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Ordinary build still reds the omitted Covers token. |
| OC63 | 16 | Breakdown refuses a foreign or phantom ownership token. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Projection tolerance cannot weaken rows-membership. |
| OC64 | 16 | A present empty ticket directory reports every unowned row. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Present empty cannot inherit absent applicability. |
| OC65 | 16 | An absent ticket directory permits read-only breakdown. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Authoring starts without pretending implementation tickets exist. |
| OC66 | 18 | Breakdown rejects an incompatible operation selector. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A charge plus breakdown operand cannot choose one silently. |
| OC67 | 18 | Help derives the new selector from the operation registry. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A missing public form fails the real root inventory. |
| OC68 | 17 | A proposed overlapping path requires an explicit blocker. | planned internal/preflight/ownership/proposal_test.go (TestProposalOrdering), actual preflight command | A second ticket's shared path catches implicit serial order. |
| OC69 | 1 | A proposed path outside the spec reports expansion required. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The unchanged spec cannot gain authority from a proposal. |
| OC70 | 1 | Dispatch refuses a proposal absent from ticket Writes. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | A wide spec cannot replace missing ticket authorization. |
| OC71 | 18 | A symlink metadata source refuses before opening. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | The outside sentinel detects any followed link. |
| OC72 | 18 | A FIFO metadata source refuses without blocking. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A bounded actual read catches special-file opening. |
| OC73 | 18 | A directory at the metadata file refuses collection. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A present wrong-type path cannot become empty metadata. |
| OC74 | 18 | An oversized metadata source refuses collection. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A partial payload cannot produce complete proof. |
| OC75 | 18 | A control-byte ownership path retains the shared refusal. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | An unrepresentable path cannot disappear from closure. |
| OC76 | 18 | A repository-escaping path refuses authority. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A lexical prefix cannot authorize an external file. |
| OC77 | 18 | Segment-prefix collisions remain outside authorization. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | An x2 path cannot inherit an x fence. |
| OC78 | 18 | Closure output spills without losing evidence reasons. | planned internal/preflight/ownership/source_test.go (TestProposalBound), actual preflight command | A full spill reconstruction detects truncated provenance. |
| OC79 | 19 | Dirty checkout wins over a missing selected ticket. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | Two competing refusals pin the existing checkout-first order. |
| OC80 | 19 | Persistent source movement publishes no proposal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | Two actual snapshot movements defeat mixed-source evidence. |
| OC81 | 19 | The current source-tip mismatch remains a refusal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | A stale full commit cannot produce current proof. |
| OC82 | 19 | Existing committed D-status ownership remains accepted. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | A differential legacy fixture protects the existing exception. |
| OC83 | 20 | An unreadable present fixture inventory remains diagnostic. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | Absence cannot mask a failed live inventory. |
| OC84 | 20 | An unreadable present anchor registry remains diagnostic. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A partial scan cannot supply complete closure. |
| OC85 | 3 | A newly required system test retains its BENCH_KIT obligation. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | A second-hop system test catches first-hop-only root checks. |
| OC86 | 1 | Dispatch refuses a ticket path absent from the spec fence. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The reciprocal union sentinel catches one-sided amendment. |
| OC87 | 13 | A failed behavioral probe is distinct from a compile failure. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | An ill-typed mutation cannot authenticate behavioral evidence. |
| OC88 | 15 | A metadata inventory cannot omit a supported seam citation. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | The second actual citation defeats an authored partial list. |
| OC89 | 15 | Manual claims prevent complete premise certification. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | Unsupported semantic work cannot be declared complete by a flag. |
| OC90 | 19 | Canonical evidence-source mismatch retains charge refusal. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | A changed required source cannot be rescued by closure green. |
| OC91 | 18 | Metadata-looking text inside another fence is not an opener. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Treating every fenced label as an opener falsely completes the inventory. |
| OC92 | 19 | Added fence roles preserve every legacy field projection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), canonical FieldScan consumers | Dropping Fenced, scope, or duplicate diagnostics fails existing parity cases. |
| OC93 | 19 | Complete current-premise proof permits charge before its future native evidence exists. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | Requiring the future witness record would block its own implementation. |
| OC94 | 19 | An unsupported required current premise refuses charge with future native work pending. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | The competing native diagnostic cannot hide current unknown proof or publish a handle. |
| OC95 | 15 | A native-future claim without exact canonical witness ownership refuses charge. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | A relabeled current claim or mismatched row, phase, ticket, or verification cannot grant an exemption. |

### Edge inventory

The audience includes Bench and repositories that link the kit.
Absent kit registry and absent fixture inventory use their canonical empty contracts.
Present unreadable inventories remain diagnostics.
Absent tickets and present empty tickets remain distinct.
Duplicate metadata, malformed Go, ambiguous symbols, and partial fixed points cannot certify completion.

Path cases include spaces, glob metacharacters, CRLF documents, segment-prefix collisions, and control bytes.
Paths escaping the repository are refused.
Symlink, FIFO, directory-at-file, and oversized metadata inputs use existing classified readers.
No prompt or external research lookup resolves an unknown claim.
The existing Go loader retains its ordinary dependency resolution and fails closed on unresolved packages.

Won't handle: reflection and plugin caller proof — the actual static caller still receives a proposal and the dynamic edge stays diagnostic.
Won't handle: non-default platform caller completeness — the default Go caller remains supported and another build context requires manual review.
Won't handle: semantic mutation adequacy — the supported test-record shape remains checkable and independent review owns adequacy.

## Ownership fences

- `.agents/skills/bench-craft-spec/references/map-discipline.md`
- `.agents/skills/bench-craft-tickets/references/slicing-checks.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/preflight_version_test.go`
- `internal/anchors/references.go`
- `internal/anchors/references_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_spec_trace_test.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/bounds/bounds.go`
- `internal/bounds/bounds_test.go`
- `internal/bounds/classify.go`
- `internal/bounds/classify_nofollow_test.go`
- `internal/canary/inventory.go`
- `internal/canary/inventory_test.go`
- `internal/canary/mutation.go`
- `internal/canary/mutation_test.go`
- `internal/canary/pins.go`
- `internal/canary/pins_test.go`
- `internal/commitment/repository/light_path.go`
- `internal/commitment/repository/light_path_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/ticket_grammar_test.go`
- `internal/consumers/blast.go`
- `internal/consumers/blast_test.go`
- `internal/consumers/loader.go`
- `internal/consumers/loader_test.go`
- `internal/consumers/resolve.go`
- `internal/consumers/rows.go`
- `internal/coverage/citation_form_test.go`
- `internal/coverage/citations.go`
- `internal/coverage/citations_test.go`
- `internal/coverage/coverage.go`
- `internal/coverage/coverage_command_test.go`
- `internal/coverage/coverage_schema_test.go`
- `internal/coverage/coverage_test.go`
- `internal/coverage/projection.go`
- `internal/coverage/rows.go`
- `internal/maps/fields.go`
- `internal/maps/maps_parse_test.go`
- `internal/maps/schema.go`
- `internal/maps/tickets.go`
- `internal/maps/tickets_test.go`
- `internal/preflight/anchor_closure_test.go`
- `internal/preflight/binary_seal_test.go`
- `internal/preflight/charge_pack.go`
- `internal/preflight/charge_test.go`
- `internal/preflight/closure.go`
- `internal/preflight/command.go`
- `internal/preflight/command_bootstrap_test.go`
- `internal/preflight/command_build_test.go`
- `internal/preflight/command_review_test.go`
- `internal/preflight/completion_plan_test.go`
- `internal/preflight/decision.go`
- `internal/preflight/decision_test.go`
- `internal/preflight/deletion_preflight_test.go`
- `internal/preflight/evidencecmd/evidence_grammar_test.go`
- `internal/preflight/evidencecmd/operations.go`
- `internal/preflight/evidencecmd/operations_test.go`
- `internal/preflight/fence_writes.go`
- `internal/preflight/fence_writes_test.go`
- `internal/preflight/gather.go`
- `internal/preflight/gather_inputs.go`
- `internal/preflight/gather_test.go`
- `internal/preflight/kit_pin_new_test.go`
- `internal/preflight/ownership`
- `internal/preflight/preflighttest/fixture.go`
- `internal/preflight/preparation.go`
- `internal/preflight/proposal.go`
- `internal/preflight/proposal_command_test.go`
- `internal/preflight/proposal_edges_test.go`
- `internal/preflight/proposal_readonly_test.go`
- `internal/preflight/proposal_test.go`
- `internal/preflight/system_tag.go`
- `internal/preflight/verdict_summary_test.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/command_test.go`
- `internal/spec/fences.go`
- `internal/spec/fences_test.go`
- `internal/spec/resolve.go`
- `internal/spec/resolve_nofollow_test.go`
- `internal/structure/facts.go`
- `internal/structure/structure.go`
- `internal/structure/structure_test.go`
- `internal/tickets/registry_data.go`
- `internal/tickets/registry_data_test.go`
- `internal/tickets/tickets.go`
- `internal/tickets/writes.go`
- `internal/tickets/writes_test.go`
- `internal/usage/parse.go`
- `internal/usage/parse_test.go`
- `reviews/preflight-ownership-closure.md`
- `specs/preflight-ownership-closure/`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/package-core-guard/bounds-classify-limit-restated`
- `tests/canary/package-core-guard/bounds-discovery-window-unwrapped`
- `tests/canary/package-core-guard/bounds-dot-import-package-alias`
- `tests/canary/package-core-guard/bounds-dot-import-wait`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/package-core-guard/bounds-intent-window-fixed`
- `tests/canary/package-core-guard/bounds-multiple-dot-import-wait`
- `tests/canary/package-core-guard/bounds-parenthesized-wait`
- `tests/canary/package-core-guard/bounds-raw-elapsed-wait`
- `tests/canary/package-core-guard/bounds-raw-injected-wait`
- `tests/canary/package-core-guard/bounds-raw-wait-deadline`
- `tests/canary/package-core-guard/bounds-raw-wait-duration`
- `tests/canary/package-core-guard/bounds-read-limit-restated`
- `tests/canary/package-core-guard/bounds-reassigned-wait-duration`
- `tests/canary/package-core-guard/bounds-redeclared-wait-duration`
- `tests/canary/package-core-guard/bounds-worktree-window-unwrapped`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows`
- `tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory`
- `tests/canary/workflow-guidance-anchors/decision-map-asset-path`
- `tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition`
- `tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules`
- `tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells`
- `tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof`
- `tests/canary/workflow-guidance-anchors/map-discipline-derived-grader`
- `tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows`
- `tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader`
- `tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller`
- `tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace`
- `tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state`
- `tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions`
- `tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep`
- `tests/canary/workflow-guidance-anchors/map-discipline-no-expectation-under-test`
- `tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace`
- `tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof`
- `tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers`
- `tests/canary/workflow-guidance-anchors/map-discipline-promise-rows`
- `tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets`
- `tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands`
- `tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim`
- `tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers`
- `tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence`
- `tests/canary/workflow-guidance-anchors/map-discipline-unexported-callers`
- `tests/canary/workflow-guidance-anchors/map-discipline-workflow-step-digests`
- `tests/canary/workflow-guidance-anchors/reader-sweep-term`

Canonical closure at a9c395e77fec36d1f60f6d057c654aecf5720e24 yields 159 fence entries and 57 fixture units.
The canonical anchor literal scan supplies 7 holder files.
The inventory was queried through BoundFiles, ReferencingFiles, and FixturePins, not a copied policy table.

Reviewer disposition: specification, applicability amendment, and ticket graph independently accepted.
The ticket union equals this fence under canonical phase exclusions.
Co-owned closure holders retain every assertion and mutation purpose.
Holder files that need no byte change stay unchanged.

## Out of scope

FT375's automatic manual-audit omission is separate: 1 outcome, estimated 12 owner edits and 4 focused verification runs.
FT318/FT317 record writes and capability values are separate: estimated 8 owner edits and 3 focused runs.
FT125 exact artifact readers are separate: estimated 9 owner edits and 4 focused runs.
These estimates describe capability cuts, not skipped acceptance.
No literal C10 order becomes a dependency.

## Further notes

Flagged additions: typed metadata and typed consumer/size projections make the reviewed mechanical obligations executable.
FieldScan adds a typed fence-state projection so the schema owner needs no second fence classifier.
The breakdown selector supplies the reviewed read-only projection.
No automatic approval or second policy registry is added.

Source-clause mapping:
ownership proposal and explicit authorization -> OC1–OC13 and OC55–OC60.
helper map, grammar obligations, reused owners, and deletion -> OC14–OC31.
counts, globs, headroom, imports, sentinels, errors, bounds, evidence, and predicates -> OC32–OC54.
uncited breakdown and row ownership -> OC61–OC70.
The full roadmap occurrences are motivating cases, not new independently implemented behaviors.

The pre-review proof checklist follows.

- Cited symbols: GatherPinned, Decide, BoundFiles, ReferencingFiles, FixturePins, Resolve, Rows, and FieldScan resolve in their cited owners.
- Import edges: existing packages were enumerated with go list.
- The new ownership child points to existing owners and never imports preflight.
- Source-row clauses and occurrences: the complete FT293 body was read and grouped above.
- Promised field labels: writes_proposal, ordering, closure, closure_evidence, source, premises, and uncited are fixed above.
- Changed-function callers: the typed consumers census found 77 static references across 17 files for six existing owner functions.
- FieldScan/FieldLine references: 14 resolved references across five source files join the preserved projection contract.
- Copy survival: OC3 and OC4 exercise the canonical fixed point through independent omission sentinels.
- Rendered-shape readers: proposal golden expectations, command registry inventories, and verdict row-list tests join the fence.
- Pin operators: exact path/kind/signature equality, exact member-set equality, and the existing structure growth comparison apply.
- Entry reads: Gather owns source reads, the typed loader owns Go loading, and current preparation owns assignment and checkout checks.
- Derived expectations: canonical owners derive facts while independent omission fixtures grade their delivery.
- Consolidated rules: current closure policy remains single-sourced and the package-only grammar rule becomes declared producer proof.
- Quantified obligations: every required source, caller, producer, retirement member, and predicate needs complete enumeration or a diagnostic.
- Workflow-step writes: author amendment still updates spec and tickets before dispatch under existing source and plan digests.

## Source evidence and limits

The named reviewed artifact is decisions/architecture-planning.md and its resolved tickets 1, 2, and 5.
The shared map remains in place.
All eight structured map sources were reopened on 2026-10-06.
The complete FT293 and FT375 roadmap bodies constrain these separate outcomes.

The research questions were ownership authority, source identity, reader closure, and the limit of mechanical proof.
Current definitions pin the claims that follow.
The proposal extends their contracts.
No runtime test, mutation experiment, native qualification, or benchmark was performed during specification.

Source definitions: internal/preflight/closure.go:24, internal/preflight/proposal.go:53, and internal/tickets/registry_data.go:21.
Go loader definitions: internal/consumers/loader.go:25, resolve.go:27, and rows.go:56.
Recorded-evidence definitions: internal/reviewrecord/check.go:15 and files.go:84.
CheckTrees grades retained source-bound occurrences; ReadTree reads their committed record through its canonical parser.
Those APIs are reachable from the new premise collector without importing recordcmd or adding a CLI record form.

Applicability source: b307c9fb41a8bc33dcee85475861b7ed89178e36.
Current charge refusal: internal/preflight/charge.go:16–23.
Preparation callers: command.go:102 and 122, charge_pack.go:108, and preparation.go:18.
Canonical operation and plan owners: evidencecmd/operations.go:17 and internal/reviewrecord/plan.go:36.
The amendment changes their planned composition, not their current implementation.

Source checkout: a9c395e77fec36d1f60f6d057c654aecf5720e24.
Production baseline: 37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68.
A source move, registry edit, parser change, or caller change invalidates this source census.
Recheck those sources before implementation.

## Completion plan

This version 1 inventory records planned obligations only, with graph review accepted.
Before implementation dispatch, record a valid version 2 plan with actual run, orchestrator, and assignment identities.
Retain author_limit 1 and one fresh author per ticket.
Existing assignment, readiness, source-pin, and independent graph approval checks remain prerequisites.

Each ticket's focused obligation names its production command and omission mutation identity.
The ticket's mutation paragraph specifies the failure and source seam; retain its exact red, restoration, and green evidence.
The broad focused checks preserve co-owned assertions and do not replace that mutation.
Every successor waits for the independently accepted predecessor chunk and its checkpoint.

```bench-completion-plan
{"version":1,"chunks":[{"id":"OC-C1","tickets":["01-propose-canonical-closure.md"],"verification":[{"id":"oc-01-focused","command":"bench test --package ./internal/preflight/... --run 'TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar'","probe":"oc-01-omission"},{"id":"oc-01-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-01-check-2","command":"bench test --package ./internal/anchors"},{"id":"oc-01-check-3","command":"bench test --package ./internal/canary"},{"id":"oc-01-check-4","command":"bench test --package ./internal/coverage"},{"id":"oc-01-check-5","command":"bench test --package ./internal/maps"},{"id":"oc-01-check-6","command":"bench test --package ./internal/bounds"},{"id":"oc-01-check-7","command":"bench test --package ./internal/spec"},{"id":"oc-01-check-8","command":"bench test --package ./internal/tickets"},{"id":"oc-01-check-9","command":"bench test --package ./cmd/bench"},{"id":"oc-01-check-10","command":"bench test --check ticket-grammar"}]},{"id":"OC-C2","tickets":["02-classify-grammar-effects.md","03-propose-helper-and-owner-closure.md"],"verification":[{"id":"oc-02-focused","command":"bench test --package ./internal/preflight/... --run TestGrammarEffects","probe":"oc-02-omission"},{"id":"oc-02-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-02-check-2","command":"bench test --package ./internal/tickets"},{"id":"oc-02-check-3","command":"bench test --package ./internal/usage"},{"id":"oc-02-check-4","command":"bench test --package ./cmd/bench"},{"id":"oc-02-check-5","command":"bench test --check ticket-grammar"},{"id":"oc-03-focused","command":"bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'","probe":"oc-03-omission"},{"id":"oc-03-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-03-check-2","command":"bench test --package ./internal/consumers"},{"id":"oc-03-check-3","command":"bench test --package ./internal/tickets"},{"id":"oc-03-check-4","command":"bench test --package ./internal/commitment/repository"},{"id":"oc-03-check-5","command":"bench test --package ./internal/reviewrecord/recordcmd"},{"id":"oc-03-check-6","command":"bench test --package ./cmd/bench"},{"id":"oc-03-check-7","command":"bench test --check ticket-grammar"}]},{"id":"OC-C3","tickets":["04-check-sizing-premises.md","05-check-symbol-and-rule-premises.md","06-check-recorded-evidence-premises.md","07-refuse-unproved-source-charges.md"],"verification":[{"id":"oc-04-focused","command":"bench test --package ./internal/preflight/... --run TestSizingPremises","probe":"oc-04-omission"},{"id":"oc-04-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-04-check-2","command":"bench test --package ./internal/structure"},{"id":"oc-05-focused","command":"bench test --package ./internal/preflight/... --run TestSymbolPremises","probe":"oc-05-omission"},{"id":"oc-05-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-05-check-2","command":"bench test --package ./internal/consumers"},{"id":"oc-05-check-3","command":"bench test --package ./internal/anchors"},{"id":"oc-06-focused","command":"bench test --package ./internal/preflight/... --run TestEvidencePremises","probe":"oc-06-omission"},{"id":"oc-06-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-06-check-2","command":"bench test --package ./internal/reviewrecord"},{"id":"oc-07-focused","command":"bench test --package ./internal/preflight/... --run 'TestChargeManualApplicability|TestChargePremises|TestPredicatePremises'","probe":"oc-07-omission"},{"id":"oc-07-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-07-check-2","command":"bench test --package ./internal/coverage"},{"id":"oc-07-check-3","command":"bench test --package ./internal/reviewrecord"}]},{"id":"OC-C4","tickets":["08-project-breakdown-and-author-authority.md"],"verification":[{"id":"oc-08-focused","command":"bench test --package ./internal/preflight/... --run 'TestBreakdownRows|TestExplicitAuthority|TestBreakdownGrammar'","probe":"oc-08-omission"},{"id":"oc-08-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"oc-08-check-2","command":"bench test --package ./cmd/bench"},{"id":"oc-08-check-3","command":"bench test --check ticket-grammar"},{"id":"oc-08-check-4","command":"bench test --check docs-currency-workflow"}]}],"final_verification":[{"id":"final-preflight","command":"bench preflight build preflight-ownership-closure"},{"id":"final-focused","command":"bench test --package ./internal/preflight/..."},{"id":"final-registry","command":"bench test --package ./cmd/bench"},{"id":"final-grammar","command":"bench test --check ticket-grammar"}]}
```

## Planning claim inventory

This authored inventory is declared complete for current source premises and the planned acceptance seams listed here.
It does not authenticate review, execution, semantic adequacy, or dispatch authority.
Acceptance source locations bind exact declared row bytes, not an independent proof that the behavior exists.
Current symbol and count premises bind the reviewed source base or tip; new Go witnesses remain planned.

Before dispatch, reconcile these current premises against the complete landed predecessor and actual composed source.
Use the existing enabling-plan procedure when a predecessor changes an owner or count.
Do not demand a planned function or future native result before the implementation that creates it.
Unsupported required current premises still refuse; inventory omissions never grant readiness.

```bench-planning-facts
{"version":1,"authority":"authored-spec","completeness":"declared-complete","baseline":"source-base","scope":[{"id":"implementation-scope","location":{"path":"specs/preflight-ownership-closure/spec.md","line":737,"needle":"## Ownership fences","source":"spec"},"paths":[".agents/skills/bench-craft-spec/references/map-discipline.md",".agents/skills/bench-craft-tickets/references/slicing-checks.md","CHANGELOG.md","CONTEXT.md","cmd/bench/command_registry.go","cmd/bench/command_registry_test.go","cmd/bench/help_inventory_test.go","cmd/bench/preflight_version_test.go","internal/anchors/references.go","internal/anchors/references_test.go","internal/anchors/registry_data.go","internal/anchors/registry_data_test.go","internal/anchors/registry_decision_maps.go","internal/anchors/registry_decision_maps_test.go","internal/anchors/registry_spec_trace_test.go","internal/anchors/registry_ticket_passes.go","internal/anchors/registry_ticket_passes_test.go","internal/bounds/bounds.go","internal/bounds/bounds_test.go","internal/bounds/classify.go","internal/bounds/classify_nofollow_test.go","internal/canary/inventory.go","internal/canary/inventory_test.go","internal/canary/mutation.go","internal/canary/mutation_test.go","internal/canary/pins.go","internal/canary/pins_test.go","internal/commitment/repository/light_path.go","internal/commitment/repository/light_path_test.go","internal/conformance/axi_query_registry_test.go","internal/conformance/subcommand_routing_table_test.go","internal/conformance/ticket_grammar_test.go","internal/consumers/blast.go","internal/consumers/blast_test.go","internal/consumers/loader.go","internal/consumers/loader_test.go","internal/consumers/resolve.go","internal/consumers/rows.go","internal/coverage/citation_form_test.go","internal/coverage/citations.go","internal/coverage/citations_test.go","internal/coverage/coverage.go","internal/coverage/coverage_command_test.go","internal/coverage/coverage_schema_test.go","internal/coverage/coverage_test.go","internal/coverage/projection.go","internal/coverage/rows.go","internal/maps/fields.go","internal/maps/maps_parse_test.go","internal/maps/schema.go","internal/maps/tickets.go","internal/maps/tickets_test.go","internal/preflight/anchor_closure_test.go","internal/preflight/binary_seal_test.go","internal/preflight/charge_pack.go","internal/preflight/charge_test.go","internal/preflight/closure.go","internal/preflight/command.go","internal/preflight/command_bootstrap_test.go","internal/preflight/command_build_test.go","internal/preflight/command_review_test.go","internal/preflight/completion_plan_test.go","internal/preflight/decision.go","internal/preflight/decision_test.go","internal/preflight/deletion_preflight_test.go","internal/preflight/evidencecmd/evidence_grammar_test.go","internal/preflight/evidencecmd/operations.go","internal/preflight/evidencecmd/operations_test.go","internal/preflight/fence_writes.go","internal/preflight/fence_writes_test.go","internal/preflight/gather.go","internal/preflight/gather_inputs.go","internal/preflight/gather_test.go","internal/preflight/kit_pin_new_test.go","internal/preflight/ownership/","internal/preflight/preflighttest/fixture.go","internal/preflight/preparation.go","internal/preflight/proposal.go","internal/preflight/proposal_command_test.go","internal/preflight/proposal_edges_test.go","internal/preflight/proposal_readonly_test.go","internal/preflight/proposal_test.go","internal/preflight/system_tag.go","internal/preflight/verdict_summary_test.go","internal/reviewrecord/recordcmd/command.go","internal/reviewrecord/recordcmd/command_test.go","internal/spec/fences.go","internal/spec/fences_test.go","internal/spec/resolve.go","internal/spec/resolve_nofollow_test.go","internal/structure/facts.go","internal/structure/structure.go","internal/structure/structure_test.go","internal/tickets/registry_data.go","internal/tickets/registry_data_test.go","internal/tickets/tickets.go","internal/tickets/writes.go","internal/tickets/writes_test.go","internal/usage/parse.go","internal/usage/parse_test.go","tests/canary/docs-currency-token-diet/signal-vocabulary-drift","tests/canary/package-core-guard/bounds-classify-limit-restated","tests/canary/package-core-guard/bounds-discovery-window-unwrapped","tests/canary/package-core-guard/bounds-dot-import-package-alias","tests/canary/package-core-guard/bounds-dot-import-wait","tests/canary/package-core-guard/bounds-duplicate-owner","tests/canary/package-core-guard/bounds-intent-window-fixed","tests/canary/package-core-guard/bounds-multiple-dot-import-wait","tests/canary/package-core-guard/bounds-parenthesized-wait","tests/canary/package-core-guard/bounds-raw-elapsed-wait","tests/canary/package-core-guard/bounds-raw-injected-wait","tests/canary/package-core-guard/bounds-raw-wait-deadline","tests/canary/package-core-guard/bounds-raw-wait-duration","tests/canary/package-core-guard/bounds-read-limit-restated","tests/canary/package-core-guard/bounds-reassigned-wait-duration","tests/canary/package-core-guard/bounds-redeclared-wait-duration","tests/canary/package-core-guard/bounds-worktree-window-unwrapped","tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns","tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary","tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary","tests/canary/workflow-guidance-anchors/context-coverage-map-term","tests/canary/workflow-guidance-anchors/context-coverage-row-parts","tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary","tests/canary/workflow-guidance-anchors/context-decision-map-term","tests/canary/workflow-guidance-anchors/context-reader-sweep-term","tests/canary/workflow-guidance-anchors/context-ticket-vocabulary","tests/canary/workflow-guidance-anchors/craft-spec-transaction-failure-rows","tests/canary/workflow-guidance-anchors/craft-spec-two-audience-inventory","tests/canary/workflow-guidance-anchors/decision-map-asset-path","tests/canary/workflow-guidance-anchors/map-discipline-addition-disposition","tests/canary/workflow-guidance-anchors/map-discipline-consolidation-both-rules","tests/canary/workflow-guidance-anchors/map-discipline-consolidation-cells","tests/canary/workflow-guidance-anchors/map-discipline-copy-survival-proof","tests/canary/workflow-guidance-anchors/map-discipline-derived-grader","tests/canary/workflow-guidance-anchors/map-discipline-either-side-rows","tests/canary/workflow-guidance-anchors/map-discipline-entry-read-grader","tests/canary/workflow-guidance-anchors/map-discipline-excluded-edge-caller","tests/canary/workflow-guidance-anchors/map-discipline-executed-root-trace","tests/canary/workflow-guidance-anchors/map-discipline-fixture-reachable-state","tests/canary/workflow-guidance-anchors/map-discipline-flagged-additions","tests/canary/workflow-guidance-anchors/map-discipline-moved-bytes-sweep","tests/canary/workflow-guidance-anchors/map-discipline-no-expectation-under-test","tests/canary/workflow-guidance-anchors/map-discipline-pin-operator-trace","tests/canary/workflow-guidance-anchors/map-discipline-pre-review-proof","tests/canary/workflow-guidance-anchors/map-discipline-pre-review-source-consumers","tests/canary/workflow-guidance-anchors/map-discipline-promise-rows","tests/canary/workflow-guidance-anchors/map-discipline-quantified-tickets","tests/canary/workflow-guidance-anchors/map-discipline-quoted-operands","tests/canary/workflow-guidance-anchors/map-discipline-shipped-surface-claim","tests/canary/workflow-guidance-anchors/map-discipline-source-sentence-table","tests/canary/workflow-guidance-anchors/map-discipline-sweep-depth-bound","tests/canary/workflow-guidance-anchors/map-discipline-sweep-direct-helpers","tests/canary/workflow-guidance-anchors/map-discipline-sweep-named-consumers","tests/canary/workflow-guidance-anchors/map-discipline-sweep-reader-fence","tests/canary/workflow-guidance-anchors/map-discipline-unexported-callers","tests/canary/workflow-guidance-anchors/map-discipline-workflow-step-digests","tests/canary/workflow-guidance-anchors/reader-sweep-term"]}],"changes":[{"id":"projection-tickets-BoundFiles","location":{"path":"specs/preflight-ownership-closure/spec.md","line":103,"needle":"Add typed provenance projections beside their existing accessors.","source":"spec"},"kind":"signature","owner":{"path":"internal/tickets/registry_data.go","package":"github.com/gibbonmi/bench/internal/tickets","symbol":"tickets.BoundFiles","kind":"func","signature":"func BoundFiles(path string) []string"},"caller_scope":"production-and-tests","state":"existing","destination":null,"dispositions":[]},{"id":"projection-anchors-ReferencingFiles","location":{"path":"specs/preflight-ownership-closure/spec.md","line":103,"needle":"Add typed provenance projections beside their existing accessors.","source":"spec"},"kind":"signature","owner":{"path":"internal/anchors/references.go","package":"github.com/gibbonmi/bench/internal/anchors","symbol":"anchors.ReferencingFiles","kind":"func","signature":"func ReferencingFiles(root string) (map[string][]string, error)"},"caller_scope":"production-and-tests","state":"existing","destination":null,"dispositions":[]},{"id":"projection-canary-FixturePins","location":{"path":"specs/preflight-ownership-closure/spec.md","line":103,"needle":"Add typed provenance projections beside their existing accessors.","source":"spec"},"kind":"signature","owner":{"path":"internal/canary/inventory.go","package":"github.com/gibbonmi/bench/internal/canary","symbol":"canary.FixturePins","kind":"func","signature":"func FixturePins(root string) (map[string][]string, error)"},"caller_scope":"production-and-tests","state":"existing","destination":null,"dispositions":[]},{"id":"relocate-canary-FixturePins","location":{"path":"specs/preflight-ownership-closure/spec.md","line":448,"needle":"Move FixturePins and its pin enumeration into pins.go","source":"spec"},"kind":"fixture-helper","owner":{"path":"internal/canary/inventory.go","package":"github.com/gibbonmi/bench/internal/canary","symbol":"canary.FixturePins","kind":"func","signature":"func FixturePins(root string) (map[string][]string, error)"},"caller_scope":"production-and-tests","state":"existing","destination":"internal/canary/pins.go","dispositions":[]},{"id":"operation-registry-change","location":{"path":"specs/preflight-ownership-closure/spec.md","line":459,"needle":"### CLI grammar and responses","source":"spec"},"kind":"grammar","owner":{"path":"internal/preflight/evidencecmd/operations.go","package":"github.com/gibbonmi/bench/internal/preflight/evidencecmd","symbol":"evidencecmd.operations","kind":"var","signature":"var operations = []Operation{\n\t{Mode: ModeReview, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,\n\t\tdescription: \"review-entry checks that a spec's artifacts agree with the tree, one count line then the red checks only\"},\n\t{Mode: ModeReview, selectors: []string{flagCharge}, required: []string{FlagBase, FlagTip}, optional: []string{flagQuota}, Kind: KindPrepareReviewEvidence, Bounded: true,\n\t\tdescription: \"prepare one immutable review evidence artifact and print its bounded orientation\"},\n\t{Mode: ModeBuild, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,\n\t\tdescription: \"build-entry checks that a spec's artifacts agree with the tree, one count line then the red checks only\"},\n\t{Mode: ModeBuild, selectors: []string{flagPlanOnly}, optional: []string{FlagBase, FlagTip}, Kind: KindPlanOnly, Bounded: true,\n\t\tdescription: \"validate the authored spec and tickets without delivery admission or a build charge\"},\n\t{Mode: ModeBuild, selectors: []string{flagCharge}, required: []string{FlagTicket, FlagBase, FlagTip}, optional: []string{flagQuota}, Kind: KindPrepareEvidence, Bounded: true,\n\t\tdescription: \"prepare one immutable build evidence artifact and print its bounded orientation\"},\n\t{Mode: ModeBuild, selectors: []string{flagPropose}, required: []string{FlagTicket, FlagBase, FlagTip}, Kind: KindProposal,\n\t\tdescription: \"propose one ticket's Writes: entries from the pinned source\"},\n\t{Mode: modeEvidence, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,\n\t\tdescription: \"print the summary of a prepared evidence artifact, or one bounded fragment at a cursor, and its exact successor\"},\n\t{Mode: modeEvidence, selectors: []string{flagSource}, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,\n\t\tdescription: \"print one bounded fragment of one declared source stream and its exact successor\"},\n\t{Mode: modeEvidence, selectors: []string{flagVerify}, Kind: KindVerifyEvidence, Bounded: true,\n\t\tdescription: \"verify every stored page and source digest of a prepared evidence artifact\"},\n\t{Mode: modeEvidence, selectors: []string{flagCurrent}, Kind: KindCurrentEvidence, Bounded: true,\n\t\tdescription: \"bind a prepared evidence artifact to the current assignment and source pair\"},\n\t{Mode: modeEvidence, selectors: []string{FlagTo}, Kind: KindExportEvidence, Bounded: true,\n\t\tdescription: \"export every verified source of a prepared evidence artifact to its own file in an absent or empty directory\"},\n\t{Mode: ModeClean, optional: []string{flagCursor}, Kind: KindCleanPlan, Bounded: true,\n\t\tdescription: \"print one bounded page of the exact evidence deletion targets and its fingerprint\"},\n\t{Mode: ModeClean, selectors: []string{flagApply}, Kind: KindCleanApply, Bounded: true,\n\t\tdescription: \"delete exactly the targets one fingerprinted cleanup plan named\"},\n}"},"caller_scope":"production-and-tests","state":"existing","destination":null,"dispositions":[]}],"premises":[{"id":"current-preflight-GatherPinned","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/preflight/gather_inputs.go","package":"github.com/gibbonmi/bench/internal/preflight","symbol":"preflight.GatherPinned","kind":"func","signature":"func GatherPinned(root, mode, slug, explicitBase, sourceTipPin string) (Facts, *BootstrapFailure)"},"expected":{"kind":"func","signature":"func GatherPinned(root, mode, slug, explicitBase, sourceTipPin string) (Facts, *BootstrapFailure)"}},{"id":"current-preflight-Decide","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/preflight/decision.go","package":"github.com/gibbonmi/bench/internal/preflight","symbol":"preflight.Decide","kind":"func","signature":"func Decide(f Facts) Verdict"},"expected":{"kind":"func","signature":"func Decide(f Facts) Verdict"}},{"id":"current-tickets-BoundFiles","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/tickets/registry_data.go","package":"github.com/gibbonmi/bench/internal/tickets","symbol":"tickets.BoundFiles","kind":"func","signature":"func BoundFiles(path string) []string"},"expected":{"kind":"func","signature":"func BoundFiles(path string) []string"}},{"id":"current-anchors-ReferencingFiles","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/anchors/references.go","package":"github.com/gibbonmi/bench/internal/anchors","symbol":"anchors.ReferencingFiles","kind":"func","signature":"func ReferencingFiles(root string) (map[string][]string, error)"},"expected":{"kind":"func","signature":"func ReferencingFiles(root string) (map[string][]string, error)"}},{"id":"current-canary-FixturePins","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/canary/inventory.go","package":"github.com/gibbonmi/bench/internal/canary","symbol":"canary.FixturePins","kind":"func","signature":"func FixturePins(root string) (map[string][]string, error)"},"expected":{"kind":"func","signature":"func FixturePins(root string) (map[string][]string, error)"}},{"id":"current-coverage-ParseSpec","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/coverage/coverage.go","package":"github.com/gibbonmi/bench/internal/coverage","symbol":"coverage.ParseSpec","kind":"func","signature":"func ParseSpec(path string) (optIn bool, ids []string, violations []string, err error)"},"expected":{"kind":"func","signature":"func ParseSpec(path string) (optIn bool, ids []string, violations []string, err error)"}},{"id":"current-tickets-ParseTicket","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/tickets/tickets.go","package":"github.com/gibbonmi/bench/internal/tickets","symbol":"tickets.ParseTicket","kind":"func","signature":"func ParseTicket(name string, content []byte, siblings []string, tag string) (Ticket, []string)"},"expected":{"kind":"func","signature":"func ParseTicket(name string, content []byte, siblings []string, tag string) (Ticket, []string)"}},{"id":"current-tickets-WritesPath","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/tickets/writes.go","package":"github.com/gibbonmi/bench/internal/tickets","symbol":"tickets.WritesPath","kind":"func","signature":"func WritesPath(entry string) (path string, isNew bool)"},"expected":{"kind":"func","signature":"func WritesPath(entry string) (path string, isNew bool)"}},{"id":"current-tickets-Covers","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/tickets/writes.go","package":"github.com/gibbonmi/bench/internal/tickets","symbol":"tickets.Covers","kind":"func","signature":"func Covers(entry, path string) bool"},"expected":{"kind":"func","signature":"func Covers(entry, path string) bool"}},{"id":"current-consumers-Resolve","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/consumers/resolve.go","package":"github.com/gibbonmi/bench/internal/consumers","symbol":"consumers.Resolve","kind":"func","signature":"func Resolve(pkgs []*Package, query string) ([]Match, error)"},"expected":{"kind":"func","signature":"func Resolve(pkgs []*Package, query string) ([]Match, error)"}},{"id":"current-consumers-Rows","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/consumers/rows.go","package":"github.com/gibbonmi/bench/internal/consumers","symbol":"consumers.Rows","kind":"func","signature":"func Rows(pkgs []*Package, target types.Object, root string) []Row"},"expected":{"kind":"func","signature":"func Rows(pkgs []*Package, target types.Object, root string) []Row"}},{"id":"current-maps-FieldScan","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/maps/fields.go","package":"github.com/gibbonmi/bench/internal/maps","symbol":"maps.FieldScan","kind":"type","signature":"type FieldScan struct {\n\tTable []FieldSpec\n\n\tScope func(line string) (identity string, changed bool)\n\n\tDuplicate func(spec FieldSpec, scope string) string\n}"},"expected":{"kind":"type","signature":"type FieldScan struct {\n\tTable []FieldSpec\n\n\tScope func(line string) (identity string, changed bool)\n\n\tDuplicate func(spec FieldSpec, scope string) string\n}"}},{"id":"current-reviewrecord-CheckTrees","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/reviewrecord/check.go","package":"github.com/gibbonmi/bench/internal/reviewrecord","symbol":"reviewrecord.CheckTrees","kind":"func","signature":"func CheckTrees(root, sourceTree, evidenceTree, tip, spec, chunk string, complete bool) error"},"expected":{"kind":"func","signature":"func CheckTrees(root, sourceTree, evidenceTree, tip, spec, chunk string, complete bool) error"}},{"id":"current-reviewrecord-ReadTree","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/reviewrecord/files.go","package":"github.com/gibbonmi/bench/internal/reviewrecord","symbol":"reviewrecord.ReadTree","kind":"func","signature":"func ReadTree(root, tree, spec string) (Record, error)"},"expected":{"kind":"func","signature":"func ReadTree(root, tree, spec string) (Record, error)"}},{"id":"current-reviewrecord-ReadPlan","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/reviewrecord/plan.go","package":"github.com/gibbonmi/bench/internal/reviewrecord","symbol":"reviewrecord.ReadPlan","kind":"func","signature":"func ReadPlan(root, tree, spec string) (Plan, error)"},"expected":{"kind":"func","signature":"func ReadPlan(root, tree, spec string) (Plan, error)"}},{"id":"current-preflight-preparedAttempts","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/preflight/preparation.go","package":"github.com/gibbonmi/bench/internal/preflight","symbol":"preflight.preparedAttempts","kind":"func","signature":"func preparedAttempts(root, mode, slug, base, sourceTip, action string, args []string, render func(Facts) (string, int)) (string, int)"},"expected":{"kind":"func","signature":"func preparedAttempts(root, mode, slug, base, sourceTip, action string, args []string, render func(Facts) (string, int)) (string, int)"}},{"id":"lines-internal-preflight-decision.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/preflight/decision.go","source":"base"},"expected":{"count":475}},{"id":"lines-internal-preflight-gather.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/preflight/gather.go","source":"base"},"expected":{"count":431}},{"id":"lines-internal-coverage-coverage.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/coverage/coverage.go","source":"base"},"expected":{"count":599}},{"id":"lines-internal-canary-inventory.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/canary/inventory.go","source":"base"},"expected":{"count":415}},{"id":"lines-internal-coverage-citations.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/coverage/citations.go","source":"base"},"expected":{"count":397}},{"id":"lines-internal-structure-structure.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/structure/structure.go","source":"base"},"expected":{"count":396}},{"id":"lines-internal-consumers-loader.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/consumers/loader.go","source":"base"},"expected":{"count":60}},{"id":"lines-internal-structure-facts.go","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"lines","state":"current","subject":{"path":"internal/structure/facts.go","source":"base"},"expected":{"count":34}},{"id":"members-internal-preflight","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"glob-members","state":"current","subject":{"pattern":"internal/preflight/*.go","scope":"internal/preflight/","source":"base"},"expected":{"members":["internal/preflight/anchor_closure_test.go","internal/preflight/binary_seal.go","internal/preflight/binary_seal_test.go","internal/preflight/charge.go","internal/preflight/charge_cells_test.go","internal/preflight/charge_edges_test.go","internal/preflight/charge_evidence_test.go","internal/preflight/charge_file_kinds_test.go","internal/preflight/charge_legacy_pack_test.go","internal/preflight/charge_nofollow_linux_test.go","internal/preflight/charge_pack.go","internal/preflight/charge_snapshot_test.go","internal/preflight/charge_test.go","internal/preflight/charge_text_test.go","internal/preflight/closure.go","internal/preflight/command.go","internal/preflight/command_bootstrap_test.go","internal/preflight/command_build_test.go","internal/preflight/command_review_test.go","internal/preflight/commitment_test.go","internal/preflight/completion_plan_test.go","internal/preflight/decision.go","internal/preflight/decision_test.go","internal/preflight/delegated_evidence_test.go","internal/preflight/deletion_preflight_test.go","internal/preflight/explicit_base_test.go","internal/preflight/fence_writes.go","internal/preflight/fence_writes_test.go","internal/preflight/gather.go","internal/preflight/gather_inputs.go","internal/preflight/gather_test.go","internal/preflight/kit_pin_new_test.go","internal/preflight/plan.go","internal/preflight/preparation.go","internal/preflight/proposal.go","internal/preflight/proposal_command_test.go","internal/preflight/proposal_edges_test.go","internal/preflight/proposal_readonly_test.go","internal/preflight/proposal_test.go","internal/preflight/review.go","internal/preflight/review_charge_test.go","internal/preflight/source_tip_test.go","internal/preflight/system_tag.go","internal/preflight/verdict_summary_test.go"]}},{"id":"members-internal-coverage","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"glob-members","state":"current","subject":{"pattern":"internal/coverage/*.go","scope":"internal/coverage/","source":"base"},"expected":{"members":["internal/coverage/citation_execution.go","internal/coverage/citation_execution_test.go","internal/coverage/citation_form_test.go","internal/coverage/citations.go","internal/coverage/citations_test.go","internal/coverage/coverage.go","internal/coverage/coverage_command_test.go","internal/coverage/coverage_schema_test.go","internal/coverage/coverage_test.go"]}},{"id":"members-internal-canary","location":{"path":"specs/preflight-ownership-closure/spec.md","line":428,"needle":"### Source placement and headroom","source":"spec"},"kind":"glob-members","state":"current","subject":{"pattern":"internal/canary/*.go","scope":"internal/canary/","source":"base"},"expected":{"members":["internal/canary/authority_read.go","internal/canary/authority_read_windows.go","internal/canary/decision.go","internal/canary/decision_test.go","internal/canary/inventory.go","internal/canary/inventory_test.go","internal/canary/mutation.go","internal/canary/mutation_test.go","internal/canary/shape_test.go","internal/canary/test_helpers_test.go"]}},{"id":"current-control-record-bound","location":{"path":"specs/preflight-ownership-closure/spec.md","line":299,"needle":"Its payload uses bounds.ControlRecordLimit","source":"spec"},"kind":"bound","state":"current","subject":{"path":"internal/bounds/bounds.go","package":"github.com/gibbonmi/bench/internal/bounds","symbol":"bounds.ControlRecordLimit","kind":"const","signature":""},"expected":{"value":"2097152"}},{"id":"current-preflight-buildChargePack","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/preflight/charge.go","package":"github.com/gibbonmi/bench/internal/preflight","symbol":"preflight.buildChargePack","kind":"func","signature":"func buildChargePack(root string, facts Facts, verdict Verdict, name string, policy []buildSourceDescriptor) (*chargeevidence.Pack, string)"},"expected":{"kind":"func","signature":"func buildChargePack(root string, facts Facts, verdict Verdict, name string, policy []buildSourceDescriptor) (*chargeevidence.Pack, string)"}},{"id":"headroom-internal-coverage","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"headroom","state":"current","subject":{"directory":"internal/coverage","files":{"create":["internal/coverage/rows.go","internal/coverage/projection.go"],"delete":[]}},"expected":{"maximum":12,"current":9,"net":2,"budget_owner":{"path":"internal/structure/structure.go","line":60,"needle":"maxFiles := envInt(\"BENCH_MAX_DIR_FILES\", 12)","source":"tip"},"grant_owner":{"path":"internal/structure/budgets.go","line":159,"needle":"func acceptedReason(","source":"tip"}}},{"id":"headroom-internal-canary","location":{"path":"specs/preflight-ownership-closure/spec.md","line":948,"needle":"## Source evidence and limits","source":"spec"},"kind":"headroom","state":"current","subject":{"directory":"internal/canary","files":{"create":["internal/canary/pins.go","internal/canary/pins_test.go"],"delete":[]}},"expected":{"maximum":12,"current":10,"net":2,"budget_owner":{"path":"internal/structure/structure.go","line":60,"needle":"maxFiles := envInt(\"BENCH_MAX_DIR_FILES\", 12)","source":"tip"},"grant_owner":{"path":"internal/structure/budgets.go","line":159,"needle":"func acceptedReason(","source":"tip"}}}],"seams":[{"id":"planned-TestProposalFixedPoint","location":{"path":"specs/preflight-ownership-closure/spec.md","line":623,"needle":"| OC1 | 1 | A missing fixture pin appears as a required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | Omitting the BASE pin cannot leave an empty proposal. |","source":"spec"},"path":"internal/preflight/ownership/proposal_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalFixedPoint","state":"planned","rows":["OC1","OC3","OC4","OC85"]},{"id":"planned-TestProposalReadOnly","location":{"path":"specs/preflight-ownership-closure/spec.md","line":624,"needle":"| OC2 | 1 | A proposal leaves spec and ticket bytes unchanged. | planned internal/preflight/ownership/proposal_test.go (TestProposalReadOnly), actual preflight command | A before/after byte sentinel detects accidental authorization. |","source":"spec"},"path":"internal/preflight/ownership/proposal_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalReadOnly","state":"planned","rows":["OC2"]},{"id":"planned-TestProposalProvenance","location":{"path":"specs/preflight-ownership-closure/spec.md","line":627,"needle":"| OC5 | 2 | Two independent reasons for one path survive in evidence. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Path deduplication cannot erase the second rule. |","source":"spec"},"path":"internal/preflight/ownership/proposal_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalProvenance","state":"planned","rows":["OC5","OC6","OC7","OC8","OC9"]},{"id":"planned-TestProposalEmpty","location":{"path":"specs/preflight-ownership-closure/spec.md","line":632,"needle":"| OC10 | 20 | An absent linked-kit anchor directory adds no invented holder. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The inventoryless root catches unconditional kit closure. |","source":"spec"},"path":"internal/preflight/ownership/proposal_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalEmpty","state":"planned","rows":["OC10","OC11","OC12"]},{"id":"planned-TestProposalUnknown","location":{"path":"specs/preflight-ownership-closure/spec.md","line":635,"needle":"| OC13 | 18 | An incomplete collector never reports complete proof. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A refused registry source catches green-by-omission. |","source":"spec"},"path":"internal/preflight/ownership/proposal_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalUnknown","state":"planned","rows":["OC13","OC83","OC84"]},{"id":"planned-TestGrammarEffects","location":{"path":"specs/preflight-ownership-closure/spec.md","line":636,"needle":"| OC14 | 4 | A body-only change adds no grammar obligation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The unrelated helper edit detects whole-package binding. |","source":"spec"},"path":"internal/preflight/ownership/effects_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestGrammarEffects","state":"planned","rows":["OC14","OC15","OC16","OC17","OC18","OC19"]},{"id":"planned-TestHelperCallers","location":{"path":"specs/preflight-ownership-closure/spec.md","line":642,"needle":"| OC20 | 6 | A changed unexported signature proposes each static test caller. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The third caller outside seed files catches sampled sweeps. |","source":"spec"},"path":"internal/preflight/ownership/effects_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestHelperCallers","state":"planned","rows":["OC20","OC21","OC22","OC23","OC53"]},{"id":"planned-TestDeletedEntry","location":{"path":"specs/preflight-ownership-closure/spec.md","line":646,"needle":"| OC24 | 8 | A deleted marker resolves an exact tracked before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The source base proves ownership despite absence at the tip. |","source":"spec"},"path":"internal/preflight/ownership/effects_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestDeletedEntry","state":"planned","rows":["OC24","OC25","OC26","OC27","OC82"]},{"id":"planned-TestOwnerPremises","location":{"path":"specs/preflight-ownership-closure/spec.md","line":650,"needle":"| OC28 | 14 | A reused owner adds its exact rule source to closure. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | An omitted home-policy owner fails the rule sentinel. |","source":"spec"},"path":"internal/preflight/ownership/effects_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestOwnerPremises","state":"planned","rows":["OC28","OC29","OC30","OC31"]},{"id":"planned-TestSizingPremises","location":{"path":"specs/preflight-ownership-closure/spec.md","line":654,"needle":"| OC32 | 9 | A line-count premise reports its changed current value. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | One added physical line defeats a stale count. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestSizingPremises","state":"planned","rows":["OC32","OC33","OC34","OC35","OC36","OC37","OC38"]},{"id":"planned-TestSymbolPremises","location":{"path":"specs/preflight-ownership-closure/spec.md","line":661,"needle":"| OC39 | 11 | A named symbol resolves to its declared owner file. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-named symbol in another owner fails identity. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestSymbolPremises","state":"planned","rows":["OC39","OC40","OC41","OC42","OC43","OC44","OC45","OC46","OC47"]},{"id":"planned-TestEvidencePremises","location":{"path":"specs/preflight-ownership-closure/spec.md","line":670,"needle":"| OC48 | 13 | A reused-test premise requires its source-bound evidence record. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A prose claim without a valid record cannot certify an observed red. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestEvidencePremises","state":"planned","rows":["OC48","OC49","OC50","OC87"]},{"id":"planned-TestPredicatePremises","location":{"path":"specs/preflight-ownership-closure/spec.md","line":673,"needle":"| OC51 | 15 | A declared predicate names exactly one existing coverage row. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A nonexistent or multirow predicate mapping is diagnostic. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestPredicatePremises","state":"planned","rows":["OC51","OC52","OC88","OC89"]},{"id":"planned-TestChargePremises","location":{"path":"specs/preflight-ownership-closure/spec.md","line":676,"needle":"| OC54 | 19 | An unproved required premise prevents charge publication. | planned internal/preflight/ownership/premises_test.go (TestChargePremises), actual preflight command | A prepared-handle absence sentinel catches dispatch on unknown facts. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestChargePremises","state":"planned","rows":["OC54"]},{"id":"planned-TestMetadataGrammar","location":{"path":"specs/preflight-ownership-closure/spec.md","line":677,"needle":"| OC55 | 18 | A duplicate metadata fence refuses complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Selecting the first duplicate cannot pass. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestMetadataGrammar","state":"planned","rows":["OC55","OC56","OC57","OC58","OC59","OC60","OC91","OC92"]},{"id":"planned-TestBreakdownRows","location":{"path":"specs/preflight-ownership-closure/spec.md","line":683,"needle":"| OC61 | 16 | Breakdown lists the actual uncited row's typed seam. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | A prose token cannot replace the coverage parser's row. |","source":"spec"},"path":"internal/preflight/ownership/authority_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestBreakdownRows","state":"planned","rows":["OC61","OC62","OC63","OC64","OC65"]},{"id":"planned-TestBreakdownGrammar","location":{"path":"specs/preflight-ownership-closure/spec.md","line":688,"needle":"| OC66 | 18 | Breakdown rejects an incompatible operation selector. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A charge plus breakdown operand cannot choose one silently. |","source":"spec"},"path":"internal/preflight/ownership/authority_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestBreakdownGrammar","state":"planned","rows":["OC66","OC67"]},{"id":"planned-TestProposalOrdering","location":{"path":"specs/preflight-ownership-closure/spec.md","line":690,"needle":"| OC68 | 17 | A proposed overlapping path requires an explicit blocker. | planned internal/preflight/ownership/proposal_test.go (TestProposalOrdering), actual preflight command | A second ticket's shared path catches implicit serial order. |","source":"spec"},"path":"internal/preflight/ownership/proposal_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalOrdering","state":"planned","rows":["OC68"]},{"id":"planned-TestExplicitAuthority","location":{"path":"specs/preflight-ownership-closure/spec.md","line":691,"needle":"| OC69 | 1 | A proposed path outside the spec reports expansion required. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The unchanged spec cannot gain authority from a proposal. |","source":"spec"},"path":"internal/preflight/ownership/authority_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestExplicitAuthority","state":"planned","rows":["OC69","OC70","OC77","OC86"]},{"id":"planned-TestSourceSafety","location":{"path":"specs/preflight-ownership-closure/spec.md","line":693,"needle":"| OC71 | 18 | A symlink metadata source refuses before opening. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | The outside sentinel detects any followed link. |","source":"spec"},"path":"internal/preflight/ownership/source_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestSourceSafety","state":"planned","rows":["OC71","OC72","OC73","OC74","OC75","OC76"]},{"id":"planned-TestProposalBound","location":{"path":"specs/preflight-ownership-closure/spec.md","line":700,"needle":"| OC78 | 18 | Closure output spills without losing evidence reasons. | planned internal/preflight/ownership/source_test.go (TestProposalBound), actual preflight command | A full spill reconstruction detects truncated provenance. |","source":"spec"},"path":"internal/preflight/ownership/source_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestProposalBound","state":"planned","rows":["OC78"]},{"id":"planned-TestSourcePrecedence","location":{"path":"specs/preflight-ownership-closure/spec.md","line":701,"needle":"| OC79 | 19 | Dirty checkout wins over a missing selected ticket. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | Two competing refusals pin the existing checkout-first order. |","source":"spec"},"path":"internal/preflight/ownership/source_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestSourcePrecedence","state":"planned","rows":["OC79","OC90"]},{"id":"planned-TestSourceMovement","location":{"path":"specs/preflight-ownership-closure/spec.md","line":702,"needle":"| OC80 | 19 | Persistent source movement publishes no proposal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | Two actual snapshot movements defeat mixed-source evidence. |","source":"spec"},"path":"internal/preflight/ownership/source_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestSourceMovement","state":"planned","rows":["OC80","OC81"]},{"id":"planned-TestChargeManualApplicability","location":{"path":"specs/preflight-ownership-closure/spec.md","line":715,"needle":"| OC93 | 19 | Complete current-premise proof permits charge before its future native evidence exists. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | Requiring the future witness record would block its own implementation. |","source":"spec"},"path":"internal/preflight/ownership/premises_test.go","package":"github.com/gibbonmi/bench/internal/preflight/ownership","function":"TestChargeManualApplicability","state":"planned","rows":["OC93","OC94","OC95"]}],"decisions":[{"id":"acceptance-OC1","location":{"path":"specs/preflight-ownership-closure/spec.md","line":623,"needle":"| OC1 | 1 | A missing fixture pin appears as a required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | Omitting the BASE pin cannot leave an empty proposal. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":623,"needle":"| OC1 | 1 | A missing fixture pin appears as a required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | Omitting the BASE pin cannot leave an empty proposal. |","source":"tip"},"predicate":"A missing fixture pin appears as a required path.","row":"OC1"},{"id":"acceptance-OC2","location":{"path":"specs/preflight-ownership-closure/spec.md","line":624,"needle":"| OC2 | 1 | A proposal leaves spec and ticket bytes unchanged. | planned internal/preflight/ownership/proposal_test.go (TestProposalReadOnly), actual preflight command | A before/after byte sentinel detects accidental authorization. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":624,"needle":"| OC2 | 1 | A proposal leaves spec and ticket bytes unchanged. | planned internal/preflight/ownership/proposal_test.go (TestProposalReadOnly), actual preflight command | A before/after byte sentinel detects accidental authorization. |","source":"tip"},"predicate":"A proposal leaves spec and ticket bytes unchanged.","row":"OC2"},{"id":"acceptance-OC3","location":{"path":"specs/preflight-ownership-closure/spec.md","line":625,"needle":"| OC3 | 3 | A required registry's own fixture enters the fixed point. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | A second-hop pin catches one-pass closure. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":625,"needle":"| OC3 | 3 | A required registry's own fixture enters the fixed point. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | A second-hop pin catches one-pass closure. |","source":"tip"},"predicate":"A required registry's own fixture enters the fixed point.","row":"OC3"},{"id":"acceptance-OC4","location":{"path":"specs/preflight-ownership-closure/spec.md","line":626,"needle":"| OC4 | 3 | A closure cycle terminates without dropping its required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | The cyclical fixture detects incomplete or infinite traversal. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":626,"needle":"| OC4 | 3 | A closure cycle terminates without dropping its required path. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | The cyclical fixture detects incomplete or infinite traversal. |","source":"tip"},"predicate":"A closure cycle terminates without dropping its required path.","row":"OC4"},{"id":"acceptance-OC5","location":{"path":"specs/preflight-ownership-closure/spec.md","line":627,"needle":"| OC5 | 2 | Two independent reasons for one path survive in evidence. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Path deduplication cannot erase the second rule. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":627,"needle":"| OC5 | 2 | Two independent reasons for one path survive in evidence. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Path deduplication cannot erase the second rule. |","source":"tip"},"predicate":"Two independent reasons for one path survive in evidence.","row":"OC5"},{"id":"acceptance-OC6","location":{"path":"specs/preflight-ownership-closure/spec.md","line":628,"needle":"| OC6 | 2 | An anchor reason names its actual holder and source line. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | A synthetic entry label cannot satisfy the location sentinel. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":628,"needle":"| OC6 | 2 | An anchor reason names its actual holder and source line. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | A synthetic entry label cannot satisfy the location sentinel. |","source":"tip"},"predicate":"An anchor reason names its actual holder and source line.","row":"OC6"},{"id":"acceptance-OC7","location":{"path":"specs/preflight-ownership-closure/spec.md","line":629,"needle":"| OC7 | 2 | A fixture include reason names the actual included BASE source. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Naming only a fixture directory misses its derivation. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":629,"needle":"| OC7 | 2 | A fixture include reason names the actual included BASE source. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | Naming only a fixture directory misses its derivation. |","source":"tip"},"predicate":"A fixture include reason names the actual included BASE source.","row":"OC7"},{"id":"acceptance-OC8","location":{"path":"specs/preflight-ownership-closure/spec.md","line":630,"needle":"| OC8 | 2 | A registry reason names the canonical rule owner. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | A copied registry location fails the owner sentinel. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":630,"needle":"| OC8 | 2 | A registry reason names the canonical rule owner. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | A copied registry location fails the owner sentinel. |","source":"tip"},"predicate":"A registry reason names the canonical rule owner.","row":"OC8"},{"id":"acceptance-OC9","location":{"path":"specs/preflight-ownership-closure/spec.md","line":631,"needle":"| OC9 | 2 | A proposal prints the exact resolved source tip. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | An ambient or truncated tip fails exact identity. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":631,"needle":"| OC9 | 2 | A proposal prints the exact resolved source tip. | planned internal/preflight/ownership/proposal_test.go (TestProposalProvenance), actual preflight command | An ambient or truncated tip fails exact identity. |","source":"tip"},"predicate":"A proposal prints the exact resolved source tip.","row":"OC9"},{"id":"acceptance-OC10","location":{"path":"specs/preflight-ownership-closure/spec.md","line":632,"needle":"| OC10 | 20 | An absent linked-kit anchor directory adds no invented holder. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The inventoryless root catches unconditional kit closure. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":632,"needle":"| OC10 | 20 | An absent linked-kit anchor directory adds no invented holder. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The inventoryless root catches unconditional kit closure. |","source":"tip"},"predicate":"An absent linked-kit anchor directory adds no invented holder.","row":"OC10"},{"id":"acceptance-OC11","location":{"path":"specs/preflight-ownership-closure/spec.md","line":633,"needle":"| OC11 | 20 | An absent linked-kit fixture inventory adds no invented fixture. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The canonical empty inventory stays observable. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":633,"needle":"| OC11 | 20 | An absent linked-kit fixture inventory adds no invented fixture. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | The canonical empty inventory stays observable. |","source":"tip"},"predicate":"An absent linked-kit fixture inventory adds no invented fixture.","row":"OC11"},{"id":"acceptance-OC12","location":{"path":"specs/preflight-ownership-closure/spec.md","line":634,"needle":"| OC12 | 2 | A complete empty proposal prints typed zero-row tables. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | Silence and a legacy prose-only response fail shape assertions. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":634,"needle":"| OC12 | 2 | A complete empty proposal prints typed zero-row tables. | planned internal/preflight/ownership/proposal_test.go (TestProposalEmpty), actual preflight command | Silence and a legacy prose-only response fail shape assertions. |","source":"tip"},"predicate":"A complete empty proposal prints typed zero-row tables.","row":"OC12"},{"id":"acceptance-OC13","location":{"path":"specs/preflight-ownership-closure/spec.md","line":635,"needle":"| OC13 | 18 | An incomplete collector never reports complete proof. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A refused registry source catches green-by-omission. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":635,"needle":"| OC13 | 18 | An incomplete collector never reports complete proof. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A refused registry source catches green-by-omission. |","source":"tip"},"predicate":"An incomplete collector never reports complete proof.","row":"OC13"},{"id":"acceptance-OC14","location":{"path":"specs/preflight-ownership-closure/spec.md","line":636,"needle":"| OC14 | 4 | A body-only change adds no grammar obligation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The unrelated helper edit detects whole-package binding. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":636,"needle":"| OC14 | 4 | A body-only change adds no grammar obligation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The unrelated helper edit detects whole-package binding. |","source":"tip"},"predicate":"A body-only change adds no grammar obligation.","row":"OC14"},{"id":"acceptance-OC15","location":{"path":"specs/preflight-ownership-closure/spec.md","line":637,"needle":"| OC15 | 4 | A core rule obligation survives a body-only change. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | Dropping all obligations cannot satisfy precision. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":637,"needle":"| OC15 | 4 | A core rule obligation survives a body-only change. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | Dropping all obligations cannot satisfy precision. |","source":"tip"},"predicate":"A core rule obligation survives a body-only change.","row":"OC15"},{"id":"acceptance-OC16","location":{"path":"specs/preflight-ownership-closure/spec.md","line":638,"needle":"| OC16 | 4 | A changed grammar producer requires its grammar source. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The actual usage declaration catches missed producer changes. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":638,"needle":"| OC16 | 4 | A changed grammar producer requires its grammar source. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | The actual usage declaration catches missed producer changes. |","source":"tip"},"predicate":"A changed grammar producer requires its grammar source.","row":"OC16"},{"id":"acceptance-OC17","location":{"path":"specs/preflight-ownership-closure/spec.md","line":639,"needle":"| OC17 | 5 | A new unbound form requires the root help inventory expectation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | An unbound package fixture catches missing shared inventories. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":639,"needle":"| OC17 | 5 | A new unbound form requires the root help inventory expectation. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | An unbound package fixture catches missing shared inventories. |","source":"tip"},"predicate":"A new unbound form requires the root help inventory expectation.","row":"OC17"},{"id":"acceptance-OC18","location":{"path":"specs/preflight-ownership-closure/spec.md","line":640,"needle":"| OC18 | 5 | An unresolved registration producer remains diagnostic. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | A guessed Command-name binding cannot complete. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":640,"needle":"| OC18 | 5 | An unresolved registration producer remains diagnostic. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | A guessed Command-name binding cannot complete. |","source":"tip"},"predicate":"An unresolved registration producer remains diagnostic.","row":"OC18"},{"id":"acceptance-OC19","location":{"path":"specs/preflight-ownership-closure/spec.md","line":641,"needle":"| OC19 | 4 | Case-sensitive binding selectors remain distinct. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | A Core versus core fixture catches silent case folding. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":641,"needle":"| OC19 | 4 | Case-sensitive binding selectors remain distinct. | planned internal/preflight/ownership/effects_test.go (TestGrammarEffects), actual preflight command | A Core versus core fixture catches silent case folding. |","source":"tip"},"predicate":"Case-sensitive binding selectors remain distinct.","row":"OC19"},{"id":"acceptance-OC20","location":{"path":"specs/preflight-ownership-closure/spec.md","line":642,"needle":"| OC20 | 6 | A changed unexported signature proposes each static test caller. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The third caller outside seed files catches sampled sweeps. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":642,"needle":"| OC20 | 6 | A changed unexported signature proposes each static test caller. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The third caller outside seed files catches sampled sweeps. |","source":"tip"},"predicate":"A changed unexported signature proposes each static test caller.","row":"OC20"},{"id":"acceptance-OC21","location":{"path":"specs/preflight-ownership-closure/spec.md","line":643,"needle":"| OC21 | 6 | A fixture-helper change proposes its external call-site file. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | A reused fixture outside the owner package catches local-only analysis. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":643,"needle":"| OC21 | 6 | A fixture-helper change proposes its external call-site file. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | A reused fixture outside the owner package catches local-only analysis. |","source":"tip"},"predicate":"A fixture-helper change proposes its external call-site file.","row":"OC21"},{"id":"acceptance-OC22","location":{"path":"specs/preflight-ownership-closure/spec.md","line":644,"needle":"| OC22 | 7 | A relocation retains the before-state helper's caller proof. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The moved declaration catches after-only enumeration. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":644,"needle":"| OC22 | 7 | A relocation retains the before-state helper's caller proof. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | The moved declaration catches after-only enumeration. |","source":"tip"},"predicate":"A relocation retains the before-state helper's caller proof.","row":"OC22"},{"id":"acceptance-OC23","location":{"path":"specs/preflight-ownership-closure/spec.md","line":645,"needle":"| OC23 | 7 | A deleted helper retains its before-state caller proof. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | An absent tip declaration cannot erase its former users. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":645,"needle":"| OC23 | 7 | A deleted helper retains its before-state caller proof. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | An absent tip declaration cannot erase its former users. |","source":"tip"},"predicate":"A deleted helper retains its before-state caller proof.","row":"OC23"},{"id":"acceptance-OC24","location":{"path":"specs/preflight-ownership-closure/spec.md","line":646,"needle":"| OC24 | 8 | A deleted marker resolves an exact tracked before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The source base proves ownership despite absence at the tip. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":646,"needle":"| OC24 | 8 | A deleted marker resolves an exact tracked before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The source base proves ownership despite absence at the tip. |","source":"tip"},"predicate":"A deleted marker resolves an exact tracked before path.","row":"OC24"},{"id":"acceptance-OC25","location":{"path":"specs/preflight-ownership-closure/spec.md","line":647,"needle":"| OC25 | 8 | A nonexistent before path refuses the deleted marker. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | A fabricated name cannot gain authority. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":647,"needle":"| OC25 | 8 | A nonexistent before path refuses the deleted marker. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | A fabricated name cannot gain authority. |","source":"tip"},"predicate":"A nonexistent before path refuses the deleted marker.","row":"OC25"},{"id":"acceptance-OC26","location":{"path":"specs/preflight-ownership-closure/spec.md","line":648,"needle":"| OC26 | 8 | A planned deleted marker accepts a still-present before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The preimplementation fixture prevents an impossible first checkpoint. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":648,"needle":"| OC26 | 8 | A planned deleted marker accepts a still-present before path. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | The preimplementation fixture prevents an impossible first checkpoint. |","source":"tip"},"predicate":"A planned deleted marker accepts a still-present before path.","row":"OC26"},{"id":"acceptance-OC27","location":{"path":"specs/preflight-ownership-closure/spec.md","line":649,"needle":"| OC27 | 8 | Combined new and deleted markers refuse parsing. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | Conflicting authority kinds cannot collapse into New. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":649,"needle":"| OC27 | 8 | Combined new and deleted markers refuse parsing. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | Conflicting authority kinds cannot collapse into New. |","source":"tip"},"predicate":"Combined new and deleted markers refuse parsing.","row":"OC27"},{"id":"acceptance-OC28","location":{"path":"specs/preflight-ownership-closure/spec.md","line":650,"needle":"| OC28 | 14 | A reused owner adds its exact rule source to closure. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | An omitted home-policy owner fails the rule sentinel. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":650,"needle":"| OC28 | 14 | A reused owner adds its exact rule source to closure. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | An omitted home-policy owner fails the rule sentinel. |","source":"tip"},"predicate":"A reused owner adds its exact rule source to closure.","row":"OC28"},{"id":"acceptance-OC29","location":{"path":"specs/preflight-ownership-closure/spec.md","line":651,"needle":"| OC29 | 14 | A liveness premise reports an omitted retirement consumer. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The second lifecycle consumer defeats a partial census. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":651,"needle":"| OC29 | 14 | A liveness premise reports an omitted retirement consumer. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The second lifecycle consumer defeats a partial census. |","source":"tip"},"predicate":"A liveness premise reports an omitted retirement consumer.","row":"OC29"},{"id":"acceptance-OC30","location":{"path":"specs/preflight-ownership-closure/spec.md","line":652,"needle":"| OC30 | 14 | A posture premise reports an omitted static writer caller. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The independent writer site catches a reader-only census. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":652,"needle":"| OC30 | 14 | A posture premise reports an omitted static writer caller. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The independent writer site catches a reader-only census. |","source":"tip"},"predicate":"A posture premise reports an omitted static writer caller.","row":"OC30"},{"id":"acceptance-OC31","location":{"path":"specs/preflight-ownership-closure/spec.md","line":653,"needle":"| OC31 | 14 | A printed-field premise reports an omitted producer. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The second renderer branch catches one-producer evidence. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":653,"needle":"| OC31 | 14 | A printed-field premise reports an omitted producer. | planned internal/preflight/ownership/effects_test.go (TestOwnerPremises), actual preflight command | The second renderer branch catches one-producer evidence. |","source":"tip"},"predicate":"A printed-field premise reports an omitted producer.","row":"OC31"},{"id":"acceptance-OC32","location":{"path":"specs/preflight-ownership-closure/spec.md","line":654,"needle":"| OC32 | 9 | A line-count premise reports its changed current value. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | One added physical line defeats a stale count. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":654,"needle":"| OC32 | 9 | A line-count premise reports its changed current value. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | One added physical line defeats a stale count. |","source":"tip"},"predicate":"A line-count premise reports its changed current value.","row":"OC32"},{"id":"acceptance-OC33","location":{"path":"specs/preflight-ownership-closure/spec.md","line":655,"needle":"| OC33 | 9 | A glob-member premise reports its omitted tracked dot member. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | A hidden tracked member defeats visible-only enumeration. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":655,"needle":"| OC33 | 9 | A glob-member premise reports its omitted tracked dot member. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | A hidden tracked member defeats visible-only enumeration. |","source":"tip"},"predicate":"A glob-member premise reports its omitted tracked dot member.","row":"OC33"},{"id":"acceptance-OC34","location":{"path":"specs/preflight-ownership-closure/spec.md","line":656,"needle":"| OC34 | 9 | A glob premise retains literal metacharacter path identity. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Shell expansion cannot supply the expected exact member. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":656,"needle":"| OC34 | 9 | A glob premise retains literal metacharacter path identity. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Shell expansion cannot supply the expected exact member. |","source":"tip"},"predicate":"A glob premise retains literal metacharacter path identity.","row":"OC34"},{"id":"acceptance-OC35","location":{"path":"specs/preflight-ownership-closure/spec.md","line":657,"needle":"| OC35 | 9 | An unsupported recursive glob remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | A double-star input cannot become an empty complete set. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":657,"needle":"| OC35 | 9 | An unsupported recursive glob remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | A double-star input cannot become an empty complete set. |","source":"tip"},"predicate":"An unsupported recursive glob remains diagnostic.","row":"OC35"},{"id":"acceptance-OC36","location":{"path":"specs/preflight-ownership-closure/spec.md","line":658,"needle":"| OC36 | 10 | A current file-headroom premise uses the canonical limit. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Changing the owner grant defeats a copied cap. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":658,"needle":"| OC36 | 10 | A current file-headroom premise uses the canonical limit. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Changing the owner grant defeats a copied cap. |","source":"tip"},"predicate":"A current file-headroom premise uses the canonical limit.","row":"OC36"},{"id":"acceptance-OC37","location":{"path":"specs/preflight-ownership-closure/spec.md","line":659,"needle":"| OC37 | 10 | A net file plan counts its exact declared deletion. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Ignoring the deletion changes the owner-derived fit. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":659,"needle":"| OC37 | 10 | A net file plan counts its exact declared deletion. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | Ignoring the deletion changes the owner-derived fit. |","source":"tip"},"predicate":"A net file plan counts its exact declared deletion.","row":"OC37"},{"id":"acceptance-OC38","location":{"path":"specs/preflight-ownership-closure/spec.md","line":660,"needle":"| OC38 | 10 | Future net-line estimates remain unobserved estimates. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | An unbuilt target cannot yield observed-fit state. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":660,"needle":"| OC38 | 10 | Future net-line estimates remain unobserved estimates. | planned internal/preflight/ownership/premises_test.go (TestSizingPremises), actual preflight command | An unbuilt target cannot yield observed-fit state. |","source":"tip"},"predicate":"Future net-line estimates remain unobserved estimates.","row":"OC38"},{"id":"acceptance-OC39","location":{"path":"specs/preflight-ownership-closure/spec.md","line":661,"needle":"| OC39 | 11 | A named symbol resolves to its declared owner file. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-named symbol in another owner fails identity. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":661,"needle":"| OC39 | 11 | A named symbol resolves to its declared owner file. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-named symbol in another owner fails identity. |","source":"tip"},"predicate":"A named symbol resolves to its declared owner file.","row":"OC39"},{"id":"acceptance-OC40","location":{"path":"specs/preflight-ownership-closure/spec.md","line":662,"needle":"| OC40 | 11 | A declared signature mismatch refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A changed parameter defeats name-only matching. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":662,"needle":"| OC40 | 11 | A declared signature mismatch refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A changed parameter defeats name-only matching. |","source":"tip"},"predicate":"A declared signature mismatch refuses that premise.","row":"OC40"},{"id":"acceptance-OC41","location":{"path":"specs/preflight-ownership-closure/spec.md","line":663,"needle":"| OC41 | 11 | A calling seam cannot import an unexported cross-package helper. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | The real loader catches private reachability. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":663,"needle":"| OC41 | 11 | A calling seam cannot import an unexported cross-package helper. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | The real loader catches private reachability. |","source":"tip"},"predicate":"A calling seam cannot import an unexported cross-package helper.","row":"OC41"},{"id":"acceptance-OC42","location":{"path":"specs/preflight-ownership-closure/spec.md","line":664,"needle":"| OC42 | 11 | A cycle in the proposed import edge refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | The actual package dependency graph catches cyclic composition. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":664,"needle":"| OC42 | 11 | A cycle in the proposed import edge refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | The actual package dependency graph catches cyclic composition. |","source":"tip"},"predicate":"A cycle in the proposed import edge refuses that premise.","row":"OC42"},{"id":"acceptance-OC43","location":{"path":"specs/preflight-ownership-closure/spec.md","line":665,"needle":"| OC43 | 11 | A supported anchor kind resolves from the canonical inventory. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An unsupported step placement cannot use a section kind. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":665,"needle":"| OC43 | 11 | A supported anchor kind resolves from the canonical inventory. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An unsupported step placement cannot use a section kind. |","source":"tip"},"predicate":"A supported anchor kind resolves from the canonical inventory.","row":"OC43"},{"id":"acceptance-OC44","location":{"path":"specs/preflight-ownership-closure/spec.md","line":666,"needle":"| OC44 | 12 | A sentinel premise resolves its exact owner declaration. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A familiar Err name elsewhere cannot certify the rule. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":666,"needle":"| OC44 | 12 | A sentinel premise resolves its exact owner declaration. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A familiar Err name elsewhere cannot certify the rule. |","source":"tip"},"predicate":"A sentinel premise resolves its exact owner declaration.","row":"OC44"},{"id":"acceptance-OC45","location":{"path":"specs/preflight-ownership-closure/spec.md","line":667,"needle":"| OC45 | 12 | An unresolved failure expression stays diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An opaque dynamic Error method cannot certify exact text. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":667,"needle":"| OC45 | 12 | An unresolved failure expression stays diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An opaque dynamic Error method cannot certify exact text. |","source":"tip"},"predicate":"An unresolved failure expression stays diagnostic.","row":"OC45"},{"id":"acceptance-OC46","location":{"path":"specs/preflight-ownership-closure/spec.md","line":668,"needle":"| OC46 | 12 | A bound premise resolves its owning static value. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-valued local literal cannot substitute for the owner. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":668,"needle":"| OC46 | 12 | A bound premise resolves its owning static value. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | A same-valued local literal cannot substitute for the owner. |","source":"tip"},"predicate":"A bound premise resolves its owning static value.","row":"OC46"},{"id":"acceptance-OC47","location":{"path":"specs/preflight-ownership-closure/spec.md","line":669,"needle":"| OC47 | 12 | An opaque bound expression remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An unknown runtime input cannot become an inferred limit. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":669,"needle":"| OC47 | 12 | An opaque bound expression remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestSymbolPremises), actual preflight command | An unknown runtime input cannot become an inferred limit. |","source":"tip"},"predicate":"An opaque bound expression remains diagnostic.","row":"OC47"},{"id":"acceptance-OC48","location":{"path":"specs/preflight-ownership-closure/spec.md","line":670,"needle":"| OC48 | 13 | A reused-test premise requires its source-bound evidence record. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A prose claim without a valid record cannot certify an observed red. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":670,"needle":"| OC48 | 13 | A reused-test premise requires its source-bound evidence record. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A prose claim without a valid record cannot certify an observed red. |","source":"tip"},"predicate":"A reused-test premise requires its source-bound evidence record.","row":"OC48"},{"id":"acceptance-OC49","location":{"path":"specs/preflight-ownership-closure/spec.md","line":671,"needle":"| OC49 | 13 | A mismatched evidence subject refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A real red for another subject cannot satisfy reuse. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":671,"needle":"| OC49 | 13 | A mismatched evidence subject refuses that premise. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | A real red for another subject cannot satisfy reuse. |","source":"tip"},"predicate":"A mismatched evidence subject refuses that premise.","row":"OC49"},{"id":"acceptance-OC50","location":{"path":"specs/preflight-ownership-closure/spec.md","line":672,"needle":"| OC50 | 13 | A failed restoration refuses reused mutation evidence. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | An un-restored behavioral probe cannot supply valid evidence. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":672,"needle":"| OC50 | 13 | A failed restoration refuses reused mutation evidence. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | An un-restored behavioral probe cannot supply valid evidence. |","source":"tip"},"predicate":"A failed restoration refuses reused mutation evidence.","row":"OC50"},{"id":"acceptance-OC51","location":{"path":"specs/preflight-ownership-closure/spec.md","line":673,"needle":"| OC51 | 15 | A declared predicate names exactly one existing coverage row. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A nonexistent or multirow predicate mapping is diagnostic. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":673,"needle":"| OC51 | 15 | A declared predicate names exactly one existing coverage row. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A nonexistent or multirow predicate mapping is diagnostic. |","source":"tip"},"predicate":"A declared predicate names exactly one existing coverage row.","row":"OC51"},{"id":"acceptance-OC52","location":{"path":"specs/preflight-ownership-closure/spec.md","line":674,"needle":"| OC52 | 15 | Changed exact predicate bytes invalidate their row premise. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A retained row ID cannot hide changed requirement bytes. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":674,"needle":"| OC52 | 15 | Changed exact predicate bytes invalidate their row premise. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | A retained row ID cannot hide changed requirement bytes. |","source":"tip"},"predicate":"Changed exact predicate bytes invalidate their row premise.","row":"OC52"},{"id":"acceptance-OC53","location":{"path":"specs/preflight-ownership-closure/spec.md","line":675,"needle":"| OC53 | 18 | An unsupported dynamic caller claim remains diagnostic. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | Reflection cannot be silently reported as an empty complete census. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":675,"needle":"| OC53 | 18 | An unsupported dynamic caller claim remains diagnostic. | planned internal/preflight/ownership/effects_test.go (TestHelperCallers), actual preflight command | Reflection cannot be silently reported as an empty complete census. |","source":"tip"},"predicate":"An unsupported dynamic caller claim remains diagnostic.","row":"OC53"},{"id":"acceptance-OC54","location":{"path":"specs/preflight-ownership-closure/spec.md","line":676,"needle":"| OC54 | 19 | An unproved required premise prevents charge publication. | planned internal/preflight/ownership/premises_test.go (TestChargePremises), actual preflight command | A prepared-handle absence sentinel catches dispatch on unknown facts. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":676,"needle":"| OC54 | 19 | An unproved required premise prevents charge publication. | planned internal/preflight/ownership/premises_test.go (TestChargePremises), actual preflight command | A prepared-handle absence sentinel catches dispatch on unknown facts. |","source":"tip"},"predicate":"An unproved required premise prevents charge publication.","row":"OC54"},{"id":"acceptance-OC55","location":{"path":"specs/preflight-ownership-closure/spec.md","line":677,"needle":"| OC55 | 18 | A duplicate metadata fence refuses complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Selecting the first duplicate cannot pass. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":677,"needle":"| OC55 | 18 | A duplicate metadata fence refuses complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Selecting the first duplicate cannot pass. |","source":"tip"},"predicate":"A duplicate metadata fence refuses complete collection.","row":"OC55"},{"id":"acceptance-OC56","location":{"path":"specs/preflight-ownership-closure/spec.md","line":678,"needle":"| OC56 | 18 | Duplicate JSON keys refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Last-key-wins decoding cannot pass. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":678,"needle":"| OC56 | 18 | Duplicate JSON keys refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Last-key-wins decoding cannot pass. |","source":"tip"},"predicate":"Duplicate JSON keys refuse complete collection.","row":"OC56"},{"id":"acceptance-OC57","location":{"path":"specs/preflight-ownership-closure/spec.md","line":679,"needle":"| OC57 | 18 | Unknown metadata fields refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Ignored undeclared promises cannot pass. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":679,"needle":"| OC57 | 18 | Unknown metadata fields refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Ignored undeclared promises cannot pass. |","source":"tip"},"predicate":"Unknown metadata fields refuse complete collection.","row":"OC57"},{"id":"acceptance-OC58","location":{"path":"specs/preflight-ownership-closure/spec.md","line":680,"needle":"| OC58 | 18 | Duplicate item IDs refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Overwriting one claim cannot produce full evidence. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":680,"needle":"| OC58 | 18 | Duplicate item IDs refuse complete collection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Overwriting one claim cannot produce full evidence. |","source":"tip"},"predicate":"Duplicate item IDs refuse complete collection.","row":"OC58"},{"id":"acceptance-OC59","location":{"path":"specs/preflight-ownership-closure/spec.md","line":681,"needle":"| OC59 | 18 | An unsupported metadata version remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Guessing a schema cannot certify completion. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":681,"needle":"| OC59 | 18 | An unsupported metadata version remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Guessing a schema cannot certify completion. |","source":"tip"},"predicate":"An unsupported metadata version remains diagnostic.","row":"OC59"},{"id":"acceptance-OC60","location":{"path":"specs/preflight-ownership-closure/spec.md","line":682,"needle":"| OC60 | 18 | An invalid metadata enum remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Case-folded or unknown kind values cannot pass. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":682,"needle":"| OC60 | 18 | An invalid metadata enum remains diagnostic. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Case-folded or unknown kind values cannot pass. |","source":"tip"},"predicate":"An invalid metadata enum remains diagnostic.","row":"OC60"},{"id":"acceptance-OC61","location":{"path":"specs/preflight-ownership-closure/spec.md","line":683,"needle":"| OC61 | 16 | Breakdown lists the actual uncited row's typed seam. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | A prose token cannot replace the coverage parser's row. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":683,"needle":"| OC61 | 16 | Breakdown lists the actual uncited row's typed seam. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | A prose token cannot replace the coverage parser's row. |","source":"tip"},"predicate":"Breakdown lists the actual uncited row's typed seam.","row":"OC61"},{"id":"acceptance-OC62","location":{"path":"specs/preflight-ownership-closure/spec.md","line":684,"needle":"| OC62 | 16 | Breakdown keeps the ordinary rows-owned refusal unchanged. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Ordinary build still reds the omitted Covers token. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":684,"needle":"| OC62 | 16 | Breakdown keeps the ordinary rows-owned refusal unchanged. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Ordinary build still reds the omitted Covers token. |","source":"tip"},"predicate":"Breakdown keeps the ordinary rows-owned refusal unchanged.","row":"OC62"},{"id":"acceptance-OC63","location":{"path":"specs/preflight-ownership-closure/spec.md","line":685,"needle":"| OC63 | 16 | Breakdown refuses a foreign or phantom ownership token. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Projection tolerance cannot weaken rows-membership. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":685,"needle":"| OC63 | 16 | Breakdown refuses a foreign or phantom ownership token. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Projection tolerance cannot weaken rows-membership. |","source":"tip"},"predicate":"Breakdown refuses a foreign or phantom ownership token.","row":"OC63"},{"id":"acceptance-OC64","location":{"path":"specs/preflight-ownership-closure/spec.md","line":686,"needle":"| OC64 | 16 | A present empty ticket directory reports every unowned row. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Present empty cannot inherit absent applicability. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":686,"needle":"| OC64 | 16 | A present empty ticket directory reports every unowned row. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Present empty cannot inherit absent applicability. |","source":"tip"},"predicate":"A present empty ticket directory reports every unowned row.","row":"OC64"},{"id":"acceptance-OC65","location":{"path":"specs/preflight-ownership-closure/spec.md","line":687,"needle":"| OC65 | 16 | An absent ticket directory permits read-only breakdown. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Authoring starts without pretending implementation tickets exist. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":687,"needle":"| OC65 | 16 | An absent ticket directory permits read-only breakdown. | planned internal/preflight/ownership/authority_test.go (TestBreakdownRows), actual preflight command | Authoring starts without pretending implementation tickets exist. |","source":"tip"},"predicate":"An absent ticket directory permits read-only breakdown.","row":"OC65"},{"id":"acceptance-OC66","location":{"path":"specs/preflight-ownership-closure/spec.md","line":688,"needle":"| OC66 | 18 | Breakdown rejects an incompatible operation selector. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A charge plus breakdown operand cannot choose one silently. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":688,"needle":"| OC66 | 18 | Breakdown rejects an incompatible operation selector. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A charge plus breakdown operand cannot choose one silently. |","source":"tip"},"predicate":"Breakdown rejects an incompatible operation selector.","row":"OC66"},{"id":"acceptance-OC67","location":{"path":"specs/preflight-ownership-closure/spec.md","line":689,"needle":"| OC67 | 18 | Help derives the new selector from the operation registry. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A missing public form fails the real root inventory. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":689,"needle":"| OC67 | 18 | Help derives the new selector from the operation registry. | planned internal/preflight/ownership/authority_test.go (TestBreakdownGrammar), actual preflight command | A missing public form fails the real root inventory. |","source":"tip"},"predicate":"Help derives the new selector from the operation registry.","row":"OC67"},{"id":"acceptance-OC68","location":{"path":"specs/preflight-ownership-closure/spec.md","line":690,"needle":"| OC68 | 17 | A proposed overlapping path requires an explicit blocker. | planned internal/preflight/ownership/proposal_test.go (TestProposalOrdering), actual preflight command | A second ticket's shared path catches implicit serial order. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":690,"needle":"| OC68 | 17 | A proposed overlapping path requires an explicit blocker. | planned internal/preflight/ownership/proposal_test.go (TestProposalOrdering), actual preflight command | A second ticket's shared path catches implicit serial order. |","source":"tip"},"predicate":"A proposed overlapping path requires an explicit blocker.","row":"OC68"},{"id":"acceptance-OC69","location":{"path":"specs/preflight-ownership-closure/spec.md","line":691,"needle":"| OC69 | 1 | A proposed path outside the spec reports expansion required. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The unchanged spec cannot gain authority from a proposal. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":691,"needle":"| OC69 | 1 | A proposed path outside the spec reports expansion required. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The unchanged spec cannot gain authority from a proposal. |","source":"tip"},"predicate":"A proposed path outside the spec reports expansion required.","row":"OC69"},{"id":"acceptance-OC70","location":{"path":"specs/preflight-ownership-closure/spec.md","line":692,"needle":"| OC70 | 1 | Dispatch refuses a proposal absent from ticket Writes. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | A wide spec cannot replace missing ticket authorization. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":692,"needle":"| OC70 | 1 | Dispatch refuses a proposal absent from ticket Writes. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | A wide spec cannot replace missing ticket authorization. |","source":"tip"},"predicate":"Dispatch refuses a proposal absent from ticket Writes.","row":"OC70"},{"id":"acceptance-OC71","location":{"path":"specs/preflight-ownership-closure/spec.md","line":693,"needle":"| OC71 | 18 | A symlink metadata source refuses before opening. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | The outside sentinel detects any followed link. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":693,"needle":"| OC71 | 18 | A symlink metadata source refuses before opening. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | The outside sentinel detects any followed link. |","source":"tip"},"predicate":"A symlink metadata source refuses before opening.","row":"OC71"},{"id":"acceptance-OC72","location":{"path":"specs/preflight-ownership-closure/spec.md","line":694,"needle":"| OC72 | 18 | A FIFO metadata source refuses without blocking. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A bounded actual read catches special-file opening. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":694,"needle":"| OC72 | 18 | A FIFO metadata source refuses without blocking. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A bounded actual read catches special-file opening. |","source":"tip"},"predicate":"A FIFO metadata source refuses without blocking.","row":"OC72"},{"id":"acceptance-OC73","location":{"path":"specs/preflight-ownership-closure/spec.md","line":695,"needle":"| OC73 | 18 | A directory at the metadata file refuses collection. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A present wrong-type path cannot become empty metadata. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":695,"needle":"| OC73 | 18 | A directory at the metadata file refuses collection. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A present wrong-type path cannot become empty metadata. |","source":"tip"},"predicate":"A directory at the metadata file refuses collection.","row":"OC73"},{"id":"acceptance-OC74","location":{"path":"specs/preflight-ownership-closure/spec.md","line":696,"needle":"| OC74 | 18 | An oversized metadata source refuses collection. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A partial payload cannot produce complete proof. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":696,"needle":"| OC74 | 18 | An oversized metadata source refuses collection. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A partial payload cannot produce complete proof. |","source":"tip"},"predicate":"An oversized metadata source refuses collection.","row":"OC74"},{"id":"acceptance-OC75","location":{"path":"specs/preflight-ownership-closure/spec.md","line":697,"needle":"| OC75 | 18 | A control-byte ownership path retains the shared refusal. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | An unrepresentable path cannot disappear from closure. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":697,"needle":"| OC75 | 18 | A control-byte ownership path retains the shared refusal. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | An unrepresentable path cannot disappear from closure. |","source":"tip"},"predicate":"A control-byte ownership path retains the shared refusal.","row":"OC75"},{"id":"acceptance-OC76","location":{"path":"specs/preflight-ownership-closure/spec.md","line":698,"needle":"| OC76 | 18 | A repository-escaping path refuses authority. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A lexical prefix cannot authorize an external file. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":698,"needle":"| OC76 | 18 | A repository-escaping path refuses authority. | planned internal/preflight/ownership/source_test.go (TestSourceSafety), actual preflight command | A lexical prefix cannot authorize an external file. |","source":"tip"},"predicate":"A repository-escaping path refuses authority.","row":"OC76"},{"id":"acceptance-OC77","location":{"path":"specs/preflight-ownership-closure/spec.md","line":699,"needle":"| OC77 | 18 | Segment-prefix collisions remain outside authorization. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | An x2 path cannot inherit an x fence. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":699,"needle":"| OC77 | 18 | Segment-prefix collisions remain outside authorization. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | An x2 path cannot inherit an x fence. |","source":"tip"},"predicate":"Segment-prefix collisions remain outside authorization.","row":"OC77"},{"id":"acceptance-OC78","location":{"path":"specs/preflight-ownership-closure/spec.md","line":700,"needle":"| OC78 | 18 | Closure output spills without losing evidence reasons. | planned internal/preflight/ownership/source_test.go (TestProposalBound), actual preflight command | A full spill reconstruction detects truncated provenance. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":700,"needle":"| OC78 | 18 | Closure output spills without losing evidence reasons. | planned internal/preflight/ownership/source_test.go (TestProposalBound), actual preflight command | A full spill reconstruction detects truncated provenance. |","source":"tip"},"predicate":"Closure output spills without losing evidence reasons.","row":"OC78"},{"id":"acceptance-OC79","location":{"path":"specs/preflight-ownership-closure/spec.md","line":701,"needle":"| OC79 | 19 | Dirty checkout wins over a missing selected ticket. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | Two competing refusals pin the existing checkout-first order. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":701,"needle":"| OC79 | 19 | Dirty checkout wins over a missing selected ticket. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | Two competing refusals pin the existing checkout-first order. |","source":"tip"},"predicate":"Dirty checkout wins over a missing selected ticket.","row":"OC79"},{"id":"acceptance-OC80","location":{"path":"specs/preflight-ownership-closure/spec.md","line":702,"needle":"| OC80 | 19 | Persistent source movement publishes no proposal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | Two actual snapshot movements defeat mixed-source evidence. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":702,"needle":"| OC80 | 19 | Persistent source movement publishes no proposal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | Two actual snapshot movements defeat mixed-source evidence. |","source":"tip"},"predicate":"Persistent source movement publishes no proposal.","row":"OC80"},{"id":"acceptance-OC81","location":{"path":"specs/preflight-ownership-closure/spec.md","line":703,"needle":"| OC81 | 19 | The current source-tip mismatch remains a refusal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | A stale full commit cannot produce current proof. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":703,"needle":"| OC81 | 19 | The current source-tip mismatch remains a refusal. | planned internal/preflight/ownership/source_test.go (TestSourceMovement), actual preflight command | A stale full commit cannot produce current proof. |","source":"tip"},"predicate":"The current source-tip mismatch remains a refusal.","row":"OC81"},{"id":"acceptance-OC82","location":{"path":"specs/preflight-ownership-closure/spec.md","line":704,"needle":"| OC82 | 19 | Existing committed D-status ownership remains accepted. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | A differential legacy fixture protects the existing exception. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":704,"needle":"| OC82 | 19 | Existing committed D-status ownership remains accepted. | planned internal/preflight/ownership/effects_test.go (TestDeletedEntry), actual preflight command | A differential legacy fixture protects the existing exception. |","source":"tip"},"predicate":"Existing committed D-status ownership remains accepted.","row":"OC82"},{"id":"acceptance-OC83","location":{"path":"specs/preflight-ownership-closure/spec.md","line":705,"needle":"| OC83 | 20 | An unreadable present fixture inventory remains diagnostic. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | Absence cannot mask a failed live inventory. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":705,"needle":"| OC83 | 20 | An unreadable present fixture inventory remains diagnostic. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | Absence cannot mask a failed live inventory. |","source":"tip"},"predicate":"An unreadable present fixture inventory remains diagnostic.","row":"OC83"},{"id":"acceptance-OC84","location":{"path":"specs/preflight-ownership-closure/spec.md","line":706,"needle":"| OC84 | 20 | An unreadable present anchor registry remains diagnostic. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A partial scan cannot supply complete closure. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":706,"needle":"| OC84 | 20 | An unreadable present anchor registry remains diagnostic. | planned internal/preflight/ownership/proposal_test.go (TestProposalUnknown), actual preflight command | A partial scan cannot supply complete closure. |","source":"tip"},"predicate":"An unreadable present anchor registry remains diagnostic.","row":"OC84"},{"id":"acceptance-OC85","location":{"path":"specs/preflight-ownership-closure/spec.md","line":707,"needle":"| OC85 | 3 | A newly required system test retains its BENCH_KIT obligation. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | A second-hop system test catches first-hop-only root checks. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":707,"needle":"| OC85 | 3 | A newly required system test retains its BENCH_KIT obligation. | planned internal/preflight/ownership/proposal_test.go (TestProposalFixedPoint), actual preflight command | A second-hop system test catches first-hop-only root checks. |","source":"tip"},"predicate":"A newly required system test retains its BENCH_KIT obligation.","row":"OC85"},{"id":"acceptance-OC86","location":{"path":"specs/preflight-ownership-closure/spec.md","line":708,"needle":"| OC86 | 1 | Dispatch refuses a ticket path absent from the spec fence. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The reciprocal union sentinel catches one-sided amendment. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":708,"needle":"| OC86 | 1 | Dispatch refuses a ticket path absent from the spec fence. | planned internal/preflight/ownership/authority_test.go (TestExplicitAuthority), actual preflight command | The reciprocal union sentinel catches one-sided amendment. |","source":"tip"},"predicate":"Dispatch refuses a ticket path absent from the spec fence.","row":"OC86"},{"id":"acceptance-OC87","location":{"path":"specs/preflight-ownership-closure/spec.md","line":709,"needle":"| OC87 | 13 | A failed behavioral probe is distinct from a compile failure. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | An ill-typed mutation cannot authenticate behavioral evidence. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":709,"needle":"| OC87 | 13 | A failed behavioral probe is distinct from a compile failure. | planned internal/preflight/ownership/premises_test.go (TestEvidencePremises), actual preflight command | An ill-typed mutation cannot authenticate behavioral evidence. |","source":"tip"},"predicate":"A failed behavioral probe is distinct from a compile failure.","row":"OC87"},{"id":"acceptance-OC88","location":{"path":"specs/preflight-ownership-closure/spec.md","line":710,"needle":"| OC88 | 15 | A metadata inventory cannot omit a supported seam citation. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | The second actual citation defeats an authored partial list. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":710,"needle":"| OC88 | 15 | A metadata inventory cannot omit a supported seam citation. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | The second actual citation defeats an authored partial list. |","source":"tip"},"predicate":"A metadata inventory cannot omit a supported seam citation.","row":"OC88"},{"id":"acceptance-OC89","location":{"path":"specs/preflight-ownership-closure/spec.md","line":711,"needle":"| OC89 | 15 | Manual claims prevent complete premise certification. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | Unsupported semantic work cannot be declared complete by a flag. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":711,"needle":"| OC89 | 15 | Manual claims prevent complete premise certification. | planned internal/preflight/ownership/premises_test.go (TestPredicatePremises), actual preflight command | Unsupported semantic work cannot be declared complete by a flag. |","source":"tip"},"predicate":"Manual claims prevent complete premise certification.","row":"OC89"},{"id":"acceptance-OC90","location":{"path":"specs/preflight-ownership-closure/spec.md","line":712,"needle":"| OC90 | 19 | Canonical evidence-source mismatch retains charge refusal. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | A changed required source cannot be rescued by closure green. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":712,"needle":"| OC90 | 19 | Canonical evidence-source mismatch retains charge refusal. | planned internal/preflight/ownership/source_test.go (TestSourcePrecedence), actual preflight command | A changed required source cannot be rescued by closure green. |","source":"tip"},"predicate":"Canonical evidence-source mismatch retains charge refusal.","row":"OC90"},{"id":"acceptance-OC91","location":{"path":"specs/preflight-ownership-closure/spec.md","line":713,"needle":"| OC91 | 18 | Metadata-looking text inside another fence is not an opener. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Treating every fenced label as an opener falsely completes the inventory. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":713,"needle":"| OC91 | 18 | Metadata-looking text inside another fence is not an opener. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), actual preflight command | Treating every fenced label as an opener falsely completes the inventory. |","source":"tip"},"predicate":"Metadata-looking text inside another fence is not an opener.","row":"OC91"},{"id":"acceptance-OC92","location":{"path":"specs/preflight-ownership-closure/spec.md","line":714,"needle":"| OC92 | 19 | Added fence roles preserve every legacy field projection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), canonical FieldScan consumers | Dropping Fenced, scope, or duplicate diagnostics fails existing parity cases. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":714,"needle":"| OC92 | 19 | Added fence roles preserve every legacy field projection. | planned internal/preflight/ownership/premises_test.go (TestMetadataGrammar), canonical FieldScan consumers | Dropping Fenced, scope, or duplicate diagnostics fails existing parity cases. |","source":"tip"},"predicate":"Added fence roles preserve every legacy field projection.","row":"OC92"},{"id":"acceptance-OC93","location":{"path":"specs/preflight-ownership-closure/spec.md","line":715,"needle":"| OC93 | 19 | Complete current-premise proof permits charge before its future native evidence exists. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | Requiring the future witness record would block its own implementation. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":715,"needle":"| OC93 | 19 | Complete current-premise proof permits charge before its future native evidence exists. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | Requiring the future witness record would block its own implementation. |","source":"tip"},"predicate":"Complete current-premise proof permits charge before its future native evidence exists.","row":"OC93"},{"id":"acceptance-OC94","location":{"path":"specs/preflight-ownership-closure/spec.md","line":716,"needle":"| OC94 | 19 | An unsupported required current premise refuses charge with future native work pending. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | The competing native diagnostic cannot hide current unknown proof or publish a handle. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":716,"needle":"| OC94 | 19 | An unsupported required current premise refuses charge with future native work pending. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | The competing native diagnostic cannot hide current unknown proof or publish a handle. |","source":"tip"},"predicate":"An unsupported required current premise refuses charge with future native work pending.","row":"OC94"},{"id":"acceptance-OC95","location":{"path":"specs/preflight-ownership-closure/spec.md","line":717,"needle":"| OC95 | 15 | A native-future claim without exact canonical witness ownership refuses charge. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | A relabeled current claim or mismatched row, phase, ticket, or verification cannot grant an exemption. |","source":"spec"},"source":{"path":"specs/preflight-ownership-closure/spec.md","line":717,"needle":"| OC95 | 15 | A native-future claim without exact canonical witness ownership refuses charge. | planned internal/preflight/ownership/premises_test.go (TestChargeManualApplicability), actual preflight command | A relabeled current claim or mismatched row, phase, ticket, or verification cannot grant an exemption. |","source":"tip"},"predicate":"A native-future claim without exact canonical witness ownership refuses charge.","row":"OC95"}],"manual_claims":[]}
```
