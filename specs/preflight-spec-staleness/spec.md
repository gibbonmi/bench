# Preflight spec staleness

Status: staged

Roadmap: FT375

Decision source: named reviewed artifact decisions/architecture-planning.md, resolved tickets 2 and 5, reviewed 2026-10-06

Verification log: 2 graph iteration(s) to accept — SS-GR1 repaired and confirmed at 412046ce; original spec and applicability remain accepted.
See assets/spec-review.md for review identity, source pins, and evidence limits.

## Problem

The build-entry procedure repeats source, symbol, test, and ownership reads by hand.
It charges the staleness pass after every initial preflight.
That current procedure lives in .agents/commands/bench-implement-spec.md:19.
The existing Facts and Decide have no drift or cited-symbol/test audit result.
Their current contract is in internal/preflight/decision.go:24.

The staleness reference stops early for an empty path drift list.
That is a procedure result, not proof that every declared citation was checked.
Missing metadata or an unsupported citation cannot justify an empty mechanical audit.

## Solution

Add a complete mechanical audit to the existing build preflight owner.
Show exact baseline and current identities with source, symbol, test, and coverage evidence.
Omit the manual pass only when every required supported check completes green and drift is empty.
Otherwise route only drift, red rows, and unsupported claims to the manual procedure.
Semantic review retains its current authority.

## User stories

Line: gpt-6-astra / high, proposed implementation line only.
Implementation-line reason: phase applicability and source identity are the hardest material chunk.
The accepted ownership metadata contract and existing movement checks provide strong seams.
Harder chunks: SS-C2.
The current spec author remains gpt-6.1-sol / high.
No binding or commitment changes here.

1. As an orchestrator, I want current drift facts, so that I audit changed paths rather than replay every read.
2. As a reviewer, I want exact source identities, so that an empty result refers to the reviewed subject.
3. As an author, I want missing symbols reported, so that an old name cannot certify a current seam.
4. As an author, I want changed symbol identity reported, so that a same-named replacement cannot silently pass.
5. As an author, I want existing test citations checked, so that removed evidence reaches the actual preflight result.
6. As an author, I want planned tests distinguished, so that future tests do not become impossible entry obligations.
7. As an orchestrator, I want coverage drift checked, so that a valid source set does not hide changed row ownership.
8. As a reviewer, I want unsupported metadata explicit, so that omitted inputs cannot appear as complete green.
9. As an orchestrator, I want a precise skip decision, so that only complete clean evidence avoids the manual pass.
10. As an author, I want narrow manual inputs, so that unsupported semantic claims retain review without repeated mechanical work.
11. As a maintainer, I want continuation applicability, so that earlier ticket tests do not block later charges.
12. As an author, I want current movement refusal, so that a mixed snapshot cannot authorize dispatch.
13. As a linked-repository author, I want absent and empty evidence distinguished, so that missing source never means unchanged source.
14. As a reviewer, I want approval retained, so that an audit result cannot approve a graph or certify a mutation.
15. As an author, I want exact grammar and help, so that the orchestrator invokes an executable audit route.
16. As a maintainer, I want preserved trust contracts, so that a digest or lexical match cannot authenticate itself.

## Implementation decisions

### Complete prerequisite

This spec consumes FT293's typed planning-facts inventory, canonical closure Snapshot, and source-bound premise diagnostics.
The complete preflight-ownership-closure outcome must be independently accepted and landed before this implementation begins.
A draft, a provider chunk, or a proposed API is insufficient.
It supplies the production metadata parser, typed citation projection, and explicit authorization checks used here.
Freeze the reviewed implementation graph against that landed provider source before dispatch.
A graph drafted earlier must reconcile its pins and claims through the existing enabling-plan procedure.

The C10 order is advice.
No FT318, FT317, or FT125 outcome is a prerequisite.
This spec adds no record forms, capability values, or exact artifact reader.

### Owner and source identity

Add internal/preflight/staleness as the cohesive mechanical audit module.
It consumes FT293's immutable ownership Snapshot and planning Manifest.
It imports neither its parent nor another workflow state owner.

Expose staleness.Collect(source ownership.Source, manifest ownership.Manifest, sourceFacts ownership.Snapshot) (AuditFacts, error).
Expose staleness.Decide(facts AuditFacts) Audit.
A typed staleness.Error carries Stage, Subject, and Cause; Unwrap preserves its cause.
The parent Gather supplies the single source pair and typed source bytes.

AuditFacts carries baseline, tip, spec identity, fenced path states, symbol states, test states, coverage states, and diagnostics.
Audit carries three named mechanical rows, typed drift entries, typed diagnostics, completeness, and ManualRequired.
The three row names are spec-drift, cited-symbols, and cited-tests.
Coverage membership and ownership use the existing rows, not a second validator.
Complete means all required supported collectors ran successfully.
It says nothing about semantic equivalence.

Metadata baseline source-base refers to the existing explicit --base identity.
The response prints its fully resolved commit and the derived source tip.
The initial workflow passes its frozen reviewed graph commit as base.
The baseline must contain the committed spec and its declared metadata.
The current source must match the committed tip and existing source-pin rules.
A missing before spec cannot be treated as an empty baseline.

This selector avoids a self-referential spec commit hash.
An author-written authority field cannot authenticate its own approval.
The reviewer-approved graph and current commitment remain external authorization.
Audit consumers check those existing authorities before dispatch.
An audit success alone does not start the outcome.

Compare the baseline spec with the current spec as exact bytes.
Include changed requirements, fences, claims, and coverage in drift.
Do not grant a semantic-digest exemption to a prose correction.
A metadata amendment moves the reviewed subject and requires the existing enabling plan procedure.

### Drift scope

Enumerate every exact fence entry and every tracked member of a fence prefix at both identities.
Use the canonical tracked source range and marker grammar.
Carry addition, deletion, modification, mode change, and rename identity.
A renamed path retains its before and after names.
A prefix matching x/ does not match x2/.
Union members before comparison so a deleted path is never lost.

The audit also includes claim source files, metadata owners, coverage source, and ticket files.
A changed active spec appears as drift.
A pure record append remains visible where its exact fence includes that record.
Only existing phase-owned exclusions used by the canonical authorization owner retain their existing treatment.
No new silent exclusion hides authored spec changes.

DriftEntry carries before path, after path, change kind, before object identity, and after object identity.
An empty result is allowed only after complete enumeration and source reads.
Unreadable before or after objects refuse collection.
Worktree content cannot replace committed source bytes.

### Symbols and tests

Use the landed typed consumers loader and resolver for current static Go symbols.
A symbol identity records owner path, declaration kind, receiver, and normalized declaration syntax.
The landed consumers AST projection supplies baseline syntax from committed bytes.
The current typed resolver supplies origin and reachability.
No arbitrary historical type-loader is promised.

Compare the declared expected identity with the baseline and current declaration.
A renamed, moved, absent, ambiguous, or changed-signature symbol reports its exact discrepancy.
Aliases resolve through the canonical consumers origin rule.
A same spelling in another owner is drift.

Only supported default-context static claims can complete green.
Reflection, plugin, non-Go language, unsupported build context, and ambiguous bare names stay diagnostic.
No lexical needle proves the declaration's semantic rule.
A current constant match does not prove a bound reaches its intended caller.

Existing test metadata names exact file and function identity.
Confirm its declaration at both baseline and current tip.
Reuse the coverage citation owner for supported named subtest syntax.
A missing file and a present file without that function are separate mismatches.
Wrong build-context evidence stays diagnostic.
Existence proves declaration, not an executed or effective test.

Planned test metadata names a future file/function and its owning row IDs.
At initial audit, the named test must be absent at the reviewed baseline and current tip.
A future file can already exist for other tests.
The exact planned function is the unit of this assertion.
A supposedly planned function that already exists is already-shipped drift.
An absent planned function is valid future work, never missing current evidence.

After implementation starts, ordinary post-ticket preflight and charge preparation do not repeat this initial planned-test absence rule.
Those paths preserve their existing source/current-action validation.
A later explicit audit deliberately compares its requested pair and can report newly implemented tests as drift.
It supplies no automatic continuation exemption or record value.

### Coverage and completeness

Use the landed typed coverage rows and ticket parser.
Preserve rows-owned, rows-membership, duplicate-token, and completion-plan obligations.
A symbol/test audit cannot mask their red result.
Compare baseline and current coverage rows and ticket Covers as part of drift.
A decision predicate without a supported row remains FT293's diagnostic.

The metadata completeness declaration is an author claim.
Review validates its scope before approval.
Collection cross-checks that every supported seam citation is represented.
An omitted citation, unrecognized table, partial inventory, manual_claims entry, or absent metadata makes Complete false.
Do not silently scan prose for a replacement inventory.
Do not infer completeness from zero collected items.

ManualRequired is false only for complete supported metadata, empty drift, green symbol/test checks, and every existing required build check green.
Required checks are the canonical build checks applicable to this source and operation.
The existing build-inapplicable diff-nonempty row remains excluded by its canonical applicability rule.
Ticket checks become required because a reviewed graph contains tickets.
A not-applicable required row cannot satisfy that conjunction.
Apply the existing commitment Ready check before publishing any omission decision.

Missing or empty tickets still cannot authorize implementation.
Any diagnostic makes ManualRequired true.
This includes a validated native-future manual entry, even when all required current premises permit its implementation charge.
Charge readiness is not audit completeness or permission to omit manual review.

### Applicability and CLI

Add:
bench preflight build <slug> --audit-spec --base <commit> --source-tip <commit>.

The operation registry requires both identities for this selector.
It cannot combine with charge, proposal, plan-only, breakdown, or evidence selectors.
Missing values, extra operands, and duplicate or unknown flags retain exit 2.
The existing registry derives help and the root inventory.
No new top-level command or reader is introduced.

Only the audit-spec form computes these build-entry rows.
Ordinary build, plan-only, review, charge, and current-evidence operations retain their existing applicability.
They cannot claim a staleness skip result.
The workflow invokes audit-spec before the first implementation edit.
After-ticket ordinary preflight remains unchanged.
This explicit consumer distinction prevents successor impossibility.

Render source[1]{base,tip}, the existing checks summary, and spec_checks[N]{check,state}.
States are green, red, or diagnostic.
A supported failed predicate is red.
An unsupported input is diagnostic and keeps manual review required.

Add drift[N]{before,after,kind,before_object,after_object} and audit_diagnostics[N]{subject,reason}.
Always print staleness[1]{complete,manual_required,baseline,tip}.
Object cells hold exact object IDs; the absent side uses the explicit absent sentinel.
Mode changes retain before and after modes in their kind detail.
Empty tables are definitive only with complete true.
Long responses retain the existing complete spill and evidence bounds.

Existing readiness reds and supported audit reds exit 1.
A diagnostic-only audit exits 0 as a successful diagnostic projection with manual_required true.
It never reports readiness complete or manual omission.
The consumer must read the typed decision, not infer it from exit zero.
A collector failure uses the existing structured refusal at exit 1 and prints no complete result.

Metadata plus dirty checkout competes with source validation.
The existing checkout/source refusal wins before a staleness result.
Persistent movement replaces the attempted result and publishes no skip decision.
Current executable seal, assignment, source-pin, evidence integrity, and current-action binding stay binding.
No audit token changes trust authority.

### Workflow consumption

The initial implement procedure invokes audit-spec with the reviewed graph base and current tip.
It omits the manual pass only when the response explicitly states complete true and manual_required false.
All declared supported checks must be green and drift empty.
Exit zero, an empty table, or another consumer's audit result is insufficient.

For manual_required true, the charge receives exact changed paths, audit reds, and unsupported claim locations.
Retain the original current-code, already-shipped, contract, blocker, and ownership review procedure for those inputs.
Do not stop the delegate from inspecting a deeper source needed to resolve one listed item.
The orchestrator reviews each return under the existing contradiction and plan-expansion rules.

A metadata gap does not become approval.
The pass can amend supported metadata and rerun the same initial audit.
Unsupported semantic items remain explicit and require their existing manual judgment.
No unsupported item is relabeled green to make the skip predicate pass.
A material acceptance change returns to the reviewer under existing policy.

### Concrete source-pair scenarios

The baseline fixture commits an approved graph with one existing test and one future planned function.
Its current tip initially has identical source objects and complete supported metadata.
The actual audit-spec response prints complete true, manual_required false, and drift[0].
The workflow may omit manual dispatch only from that explicit conjunction.

A second commit removes the existing function while keeping its file.
The same audit route must print red cited-tests and the exact function discrepancy.
Another commit creates the planned function before initial dispatch.
That run reports already-shipped drift.
An ordinary after-ticket preflight remains usable after a legitimate predecessor creates its own planned test.

A third fixture keeps paths and tests unchanged but names one unsupported dynamic symbol claim.
Its diagnostic-only command exits zero and prints manual_required true.
The native workflow must charge the manual pass with that exact claim location.
Replacing the typed decision with an exit-code-only test fails this route.

```text
bench preflight build example --audit-spec --base BASE --source-tip TIP
bench preflight build example --audit-spec --charge --base BASE --source-tip TIP
```

BASE and TIP are exact committed fixture IDs.
The first command prints the audit projection.
The second is exit 2 and runs no collection or charge publication.

## Implementation chunks

These are complete planned green consumer outcomes.
The ticket graph is independently accepted before implementation.
Each successor waits for the preceding independent chunk review.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| SS-C1 / 01-audit-source-drift.md | Source drift through the real audit-spec command | SS1–SS12, SS35–SS40, SS51–SS52, SS56 | Pinned-source command fixtures | no |
| SS-C2 / 02-audit-symbol-identities.md, 03-audit-test-citations.md, 04-audit-coverage-and-completeness.md | Symbol, test, coverage, and completeness classification | SS13–SS34, SS49–SS50, SS57–SS58 | Actual loader and audit command fixtures | yes |
| SS-C3 / 05-consume-exact-manual-omission.md | Workflow consumes the exact manual omission decision | SS41–SS48, SS53–SS55 | Root dispatch and guidance witnesses | no |

SS-C1 uses the complete landed ownership outcome.
It reports manual required while successor collectors are incomplete.
It never claims that its initial drift collector alone permits omission.
SS-C2 supplies the full decision consumed by SS-C3.
The complete build-entry workflow remains manual until SS-C3 lands.

Graph slicing supplies this spec's complete current-source claim inventory and exact planned seam identities before dispatch.
Independent graph review validates that inventory; its authored completeness field cannot replace review.
All five native acceptance rows remain explicit, unobserved manual entries rather than supported Go seams.
The landed FT293 owner validates their native-future applicability and exact graph ownership.
It does not require their future execution evidence before the implementation that enables them.

| native row | chunk | proposed owning ticket | planned verification | status |
|---|---|---|---|---|
| SS47 | SS-C3 | 05-consume-exact-manual-omission.md | ss-native-diagnostic | unobserved; acceptance-required manual diagnostic |
| SS48 | SS-C3 | 05-consume-exact-manual-omission.md | ss-native-diagnostic | unobserved; acceptance-required manual diagnostic |
| SS53 | SS-C3 | 05-consume-exact-manual-omission.md | ss-native-clean | unobserved; acceptance-required manual diagnostic |
| SS54 | SS-C3 | 05-consume-exact-manual-omission.md | ss-native-diagnostic | unobserved; acceptance-required manual diagnostic |
| SS55 | SS-C3 | 05-consume-exact-manual-omission.md | ss-native-diagnostic | unobserved; acceptance-required manual diagnostic |

Each manual item pins its own original coverage-row location with the canonical exact Location fields.
Its phase_owner names the pinned native procedure at .agents/commands/bench-implement-spec.md.
Graph slicing must make these proposed ownership and verification names actual canonical plan entries before dispatch.
Independent graph approval remains required; this table and the metadata grant no execution authority.

These entries keep Audit.Complete false and ManualRequired true while they remain manual entries.
They cannot be silently removed from inventory or promoted by matching Go tests, excerpt hashes, or exit zero.
Their planned native sessions still prove actual clean and diagnostic workflow consumption after implementation.
Any unsupported current-source requirement continues to refuse a charge under the landed prerequisite.

### Source placement

Consume the landed ownership projections without modifying their policy.
Add analysis and tests only in the staleness child.
No new file joins the existing preflight root.
The landed prerequisite already made room in Gather and the pure decision owner.
Wire audit collection into that moved collection seam and retain the existing ordinary constructors.
Remeasure actual canonical Growth at each checkpoint; no new grant is authorized.

## Testing decisions

Drive the actual audit-spec operation through the existing preflight command fixture.
Use exact committed before and after objects.
A helper returning canned empty arrays cannot stand in for real source enumeration.
Pure Audit decision tables test the conjunction after collector witnesses exist.

TestAuditCompleteness supplies complete current-source proof with the exact five pending native entries through the actual audit command.
It asserts both incomplete audit and required manual work despite charge readiness.
Dropping the native diagnostic must fail these independent response sentinels.

One multi-commit fixture changes a fenced file, a cited signature, an existing test, and a planned function in separate commits.
Each run names its exact pair and expected typed row.
Omission mutations remove one collector result and must fail the corresponding independent sentinel.
Do not create a redundant test per row.

A workflow witness drives the actual returned decision into the documented dispatch route.
The complete clean case skips manual dispatch.
Diagnostic-only exit zero must still dispatch the manual pass.
An after-ticket ordinary preflight witness preserves prior continuation behavior.
Fresh-session adoption verifies the final guidance through existing anchor and canary owners.

The witness is one real native initial author session for each clean and diagnostic route, with its charge and actual command response retained.
A Go helper cannot stand in for the documented harness decision.
The executable command fixture proves returned decision values; the native witness proves their actual workflow consumption.

The clean native session targets a separate approved fixture with complete supported metadata and no manual entries.
The diagnostic session targets an approved fixture with explicit unsupported entries.
Neither session treats this spec's own pending native inventory as complete supported proof.
The existing docs-currency-workflow check protects required and forbidden guidance anchors.

Planned verification:

```text
bench test --package ./internal/preflight/...
bench test --package ./internal/consumers
bench test --package ./internal/coverage
bench test --package ./cmd/bench
bench test --check docs-currency-workflow
bench test --check ticket-grammar
```

Each decisive omission runs through bench probe with restoration and the same focused route.
No current-source claim is supported by a planned test merely existing.
Semantic equivalence remains independent review work.

### Seam diagram

    reviewed graph base + current tip + landed claim inventory
        -> Gather pinned source objects
        -> staleness.AuditFacts
        -> staleness.Decide
        -> audit-spec response
        -> initial workflow: skip or narrow manual pass
        <- tests drive actual collection and consumer routing

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SS1 | 1 | A modified fenced path appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A changed blob defeats empty enumeration. |
| SS2 | 1 | A deleted fenced path appears with its before identity. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | After-only listing cannot lose removal. |
| SS3 | 1 | A new fenced path appears with its after identity. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Before-only listing cannot lose addition. |
| SS4 | 1 | A renamed fenced path retains both names. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A move cannot appear as unchanged source. |
| SS5 | 1 | A file-mode change appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Byte equality cannot hide changed executable mode. |
| SS6 | 1 | A fence prefix includes each tracked dot-path member. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A hidden member catches incomplete scope enumeration. |
| SS7 | 2 | The audit prints exact resolved baseline and tip identities. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An ambient branch spelling cannot identify the frozen subject. |
| SS8 | 2 | An absent before spec refuses audit collection. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | Missing source cannot become a zero drift set. |
| SS9 | 2 | Changed active spec bytes appear as drift. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | A prose correction cannot claim a semantic digest exemption. |
| SS10 | 13 | An unchanged exact member set prints a complete empty drift table. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Silence or incomplete enumeration cannot pass. |
| SS11 | 1 | A segment-prefix collision stays outside the fenced scope. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | An x2 source cannot be confused with x. |
| SS12 | 2 | A claim source outside the write fence joins audit scope. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An unchanged writes set cannot hide a changed premise owner. |
| SS13 | 3 | An absent cited symbol reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Removing the actual declaration defeats name retention. |
| SS14 | 4 | A cited signature change reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Same spelling cannot hide a changed seam. |
| SS15 | 4 | A cited symbol moved to another owner reports that discrepancy. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Name-only resolution cannot establish owner identity. |
| SS16 | 4 | A wrong declaration kind reports its exact mismatch. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | A var replacing a const cannot preserve the claimed rule. |
| SS17 | 8 | An ambiguous bare symbol remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Selecting one candidate cannot yield complete proof. |
| SS18 | 3 | An alias resolves through the canonical origin rule. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Independent lexical lookup would disagree with consumers. |
| SS19 | 8 | A non-default context claim remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | A default loader cannot certify another platform. |
| SS20 | 5 | An absent existing test file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Deleting the file cannot leave green evidence. |
| SS21 | 5 | A removed function in an existing file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | File existence cannot stand in for the named test. |
| SS22 | 6 | An absent planned test function remains valid future work. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | No first checkpoint requires its new test to exist. |
| SS23 | 6 | An existing supposedly planned function reports already-shipped drift. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Future labeling cannot hide landed work. |
| SS24 | 6 | Other functions in a planned test's existing file remain allowed. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Whole-file absence would make a legal slice impossible. |
| SS25 | 5 | A named supported subtest uses the canonical citation grammar. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | A second subtest parser cannot supply inconsistent identity. |
| SS26 | 8 | An unsupported test citation remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Unknown syntax cannot produce an empty complete result. |
| SS27 | 5 | Wrong-context test evidence remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | A declaration outside the graded context is not current evidence. |
| SS28 | 7 | An unowned coverage row retains red rows-owned. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Source green cannot mask missing ticket ownership. |
| SS29 | 7 | A phantom Covers token retains red rows-membership. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Audit tolerance cannot admit an invented row. |
| SS30 | 7 | Changed coverage predicate bytes join drift. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Stable row IDs cannot hide changed acceptance. |
| SS31 | 7 | Changed ticket Covers membership joins drift. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Unchanged tests cannot hide a different row owner. |
| SS32 | 8 | Missing planning metadata reports incomplete audit evidence. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Zero parsed items cannot mean complete green. |
| SS33 | 8 | Partial declared metadata keeps complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | An author flag cannot bless omitted inventory. |
| SS34 | 8 | An omitted supported citation keeps complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | The second citation defeats claimed completeness. |
| SS35 | 12 | Dirty source refuses before a staleness decision. | planned internal/preflight/staleness/drift_test.go (TestAuditPrecedence), actual preflight command | A competing missing-symbol error cannot replace source authority. |
| SS36 | 12 | Persistent snapshot movement publishes no skip decision. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A second movement defeats mixed before/after proof. |
| SS37 | 16 | A source-tip mismatch retains its current refusal. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A stale pin cannot authenticate the subject. |
| SS38 | 16 | A required executable-seal refusal remains binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | Green static facts cannot launch an untrusted binary. |
| SS39 | 15 | Audit-spec rejects combination with another selector. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | An audit plus charge request cannot choose one silently. |
| SS40 | 15 | Audit-spec requires both explicit source identities. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | A missing base cannot imply an authoritative baseline. |
| SS41 | 9 | Complete clean supported evidence prints manual_required false. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | The actual conjunction is the only omission producer. |
| SS42 | 9 | Any supported audit red prints manual_required true. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | An empty drift set cannot override missing evidence. |
| SS43 | 9 | Diagnostic-only exit zero prints manual_required true. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | Success of the read cannot mean success of complete proof. |
| SS44 | 9 | A not-applicable required build row prevents omission. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | Missing implementation tickets cannot enter the fast route. |
| SS45 | 11 | Ordinary after-ticket preflight keeps initial test-audit inapplicability. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | An implemented predecessor test cannot block its successor. |
| SS46 | 11 | Charge preparation grants no manual-omission result. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | Continuation cannot inherit initial audit authority. |
| SS47 | 10 | The manual charge retains drift, reds, and unsupported locations. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | A narrowed charge cannot drop the remaining semantic work. |
| SS48 | 14 | Workflow guidance retains reviewer approval before dispatch. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | A complete audit cannot become graph approval. |
| SS49 | 8 | Duplicate metadata refuses complete audit collection. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | A first-fence choice cannot establish full inventory. |
| SS50 | 13 | A present empty declared claim set differs from missing metadata. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | An empty syntax container cannot certify completeness by itself. |
| SS51 | 12 | A special before/after source refuses without following it. | planned internal/preflight/staleness/drift_test.go (TestAuditSafety), actual preflight command | An outside source sentinel detects a followed link. |
| SS52 | 15 | An audit response retains all drift rows through existing spill. | planned internal/preflight/staleness/drift_test.go (TestAuditBound), actual preflight command | Full reconstruction catches truncation hidden by bounded output. |
| SS53 | 10 | A fresh initial session follows the actual clean skip decision. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | The native session must omit manual dispatch only from complete proof. |
| SS54 | 10 | A fresh initial session still audits unsupported claims. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | The real diagnostic-only command cannot suppress manual dispatch. |
| SS55 | 14 | A manual acceptance change retains the existing reviewer route. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | A mechanical gap cannot silently widen acceptance. |
| SS56 | 16 | Existing required evidence bytes retain current-action binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | A source digest cannot replace delivered current context. |
| SS57 | 8 | A pending native-future manual entry keeps audit Complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Charge-ready current proof cannot certify the entire reviewed inventory mechanically. |
| SS58 | 8 | A pending native-future manual entry keeps manual_required true. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Diagnostic-only exit zero cannot turn future evidence applicability into a manual-audit skip. |

### Edge inventory

The audience includes Bench and linked repositories.
A missing metadata fence differs from an empty declared inventory.
An absent before spec differs from unchanged bytes.
An absent test file differs from a missing function in an existing file.
Planned function absence is expected before implementation.

Hostile inputs include malformed metadata, special files, unknown citations, ambiguous names, stale pins, and persistent movement.
Safe path spelling retains spaces and metacharacters literally.
Control-byte operands retain the shared refusal.
No unsupported state can yield complete true.

Won't handle: semantic equivalence from matching signatures — the static declaration remains supported and semantic review remains manual.
Won't handle: historical arbitrary-platform type loading — current default-context symbols remain supported and other contexts stay diagnostic.
Won't handle: prose-only digest exemptions — ordinary exact-byte spec drift remains supported.

## Ownership fences

- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-implement-spec/references/staleness-pass.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/preflight_version_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/preflight/binary_seal_test.go`
- `internal/preflight/charge_pack.go`
- `internal/preflight/charge_test.go`
- `internal/preflight/command.go`
- `internal/preflight/command_build_test.go`
- `internal/preflight/completion_plan_test.go`
- `internal/preflight/decision.go`
- `internal/preflight/decision_test.go`
- `internal/preflight/evidencecmd/evidence_grammar_test.go`
- `internal/preflight/evidencecmd/operations.go`
- `internal/preflight/evidencecmd/operations_test.go`
- `internal/preflight/explicit_base_test.go`
- `internal/preflight/gather.go`
- `internal/preflight/gather_inputs.go`
- `internal/preflight/gather_test.go`
- `internal/preflight/preflighttest/fixture.go`
- `internal/preflight/source_tip_test.go`
- `internal/preflight/staleness`
- `internal/preflight/verdict_summary_test.go`
- `reviews/preflight-spec-staleness.md`
- `specs/preflight-spec-staleness/`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/dg-25`
- `tests/canary/workflow-guidance-anchors/dg-26`
- `tests/canary/workflow-guidance-anchors/dg-29`
- `tests/canary/workflow-guidance-anchors/dg-29-verification-target`
- `tests/canary/workflow-guidance-anchors/dg-30`
- `tests/canary/workflow-guidance-anchors/dg-31`
- `tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger`
- `tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`

Canonical closure at a9c395e77fec36d1f60f6d057c654aecf5720e24 yields 80 fence entries and 39 fixture units.
The canonical anchor literal scan supplies 10 holder files.
The inventory was queried through BoundFiles, ReferencingFiles, and FixturePins, not a copied policy table.

Reviewer disposition: specification, applicability amendment, and ticket graph independently accepted.
Every transitive holder preserves its assertion and mutation purpose.
The ticket union matches the canonical fence under its phase exclusions.
No live unrelated spec or ticket migration belongs to this build.

## Out of scope

FT293 implementation is an actual complete prerequisite, not a scope cut.
FT318/FT317 new record writes and capability values remain separate: estimated 8 edits and 3 focused runs.
FT125 exact readers remain separate: estimated 9 edits and 4 focused runs.
A generic semantic equivalence engine is excluded: a future spec would need at least 3 owners and 3 qualification runs.
These estimates do not replace an in-scope predicate.

## Further notes

Flagged additions: audit-spec is the concrete selector for the reviewed build-entry capability.
Its applicability preserves current charge and post-ticket behavior.
Typed diagnostic state prevents unsupported metadata from masquerading as a clean empty result.

Source-clause mapping:
changed fenced paths and exact source identity -> SS1–SS12.
cited symbols -> SS13–SS19.
existing and planned tests -> SS20–SS27.
coverage and complete evidence -> SS28–SS34.
source movement and refusal -> SS35–SS40.
manual omission and narrowed dispatch -> SS41–SS48.

The pre-review proof checklist follows.

- Cited symbols: GatherPinned, Decide, Resolve, ParseSpec, ParseTicket, and MovementCheckedRetry resolve in their current owners.
- Import edges: existing owners were enumerated with go list.
- The new staleness child consumes landed ownership values and never imports its parent.
- Source-row clauses and occurrences: the full FT375 body and its named hand-audit occurrence were read.
- Promised field labels: source, checks, spec_checks, drift, audit_diagnostics, and staleness are fixed above.
- Changed-function callers: existing parent consumers and row-list/summary tests join the fence.
- Copy survival: SS29 and SS30 retain canonical coverage ownership checks.
- Rendered-shape readers: operation help, root inventories, and workflow anchor holders join the fence.
- Pin operators: exact object identity, exact declaration identity, and an all-green conjunction apply.
- Entry reads: the existing movement-checked Gather owns source collection.
- Derived expectations: committed fixture objects and independent omission sentinels grade collector omissions.
- Consolidated rules: the existing manual procedure consumes typed drift instead of re-enumerating it.
- Quantified obligations: every required supported collector must complete before skip.
- Workflow-step writes: an enabling amendment still changes source and plan digests under the existing record procedure.

## Source evidence and limits

The named reviewed artifact is decisions/architecture-planning.md and its resolved tickets 1, 2, and 5.
The shared map remains in place.
All eight structured map sources were reopened on 2026-10-06.
The complete FT293 and FT375 roadmap bodies constrain these separate outcomes.

The research questions were ownership authority, source identity, reader closure, and the limit of mechanical proof.
Current definitions pin the claims that follow.
The proposal extends their contracts.
No runtime test, mutation experiment, native qualification, or benchmark was performed during specification.

Source definitions: internal/preflight/decision.go:24, command.go:144, and gather_inputs.go:26.
Procedure source: .agents/commands/bench-implement-spec.md:19 and its staleness-pass reference.
Typed declaration and citation sources: internal/consumers/resolve.go:27 and internal/coverage/citations.go:197.

Applicability source: b307c9fb41a8bc33dcee85475861b7ed89178e36.
The native procedure and canonical completion plan own future verification after implementation.
Their pending evidence never proves current execution or automatic manual omission.

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
{"version":1,"chunks":[{"id":"SS-C1","tickets":["01-audit-source-drift.md"],"verification":[{"id":"ss-01-focused","command":"bench test --package ./internal/preflight/... --run 'TestDriftObjects|TestDriftIdentity|TestAuditMovement'","probe":"ss-01-omission"},{"id":"ss-01-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"ss-01-check-2","command":"bench test --package ./cmd/bench"},{"id":"ss-01-check-3","command":"bench test --check ticket-grammar"}]},{"id":"SS-C2","tickets":["02-audit-symbol-identities.md","03-audit-test-citations.md","04-audit-coverage-and-completeness.md"],"verification":[{"id":"ss-02-focused","command":"bench test --package ./internal/preflight/... --run TestSymbolIdentity","probe":"ss-02-omission"},{"id":"ss-02-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"ss-02-check-2","command":"bench test --package ./internal/consumers"},{"id":"ss-03-focused","command":"bench test --package ./internal/preflight/... --run TestTestIdentity","probe":"ss-03-omission"},{"id":"ss-03-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"ss-03-check-2","command":"bench test --package ./internal/coverage"},{"id":"ss-04-focused","command":"bench test --package ./internal/preflight/... --run 'TestCoverageIdentity|TestAuditCompleteness'","probe":"ss-04-omission"},{"id":"ss-04-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"ss-04-check-2","command":"bench test --package ./internal/coverage"},{"id":"ss-04-check-3","command":"bench test --check ticket-grammar"}]},{"id":"SS-C3","tickets":["05-consume-exact-manual-omission.md"],"verification":[{"id":"ss-05-focused","command":"bench test --package ./internal/preflight/... --run 'TestAuditDecision|TestAuditApplicability'","probe":"ss-05-omission"},{"id":"ss-05-check-1","command":"bench test --package ./internal/preflight/..."},{"id":"ss-05-check-2","command":"bench test --package ./cmd/bench"},{"id":"ss-05-check-3","command":"bench test --check docs-currency-workflow"},{"id":"ss-05-check-4","command":"bench test --check ticket-grammar"},{"id":"ss-native-clean","command":"$bench-implement-spec specs/example/spec.md"},{"id":"ss-native-diagnostic","command":"$bench-implement-spec specs/example/spec.md"}]}],"final_verification":[{"id":"final-preflight","command":"bench preflight build preflight-spec-staleness"},{"id":"final-focused","command":"bench test --package ./internal/preflight/..."},{"id":"final-registry","command":"bench test --package ./cmd/bench"},{"id":"final-grammar","command":"bench test --check ticket-grammar"},{"id":"final-workflow","command":"bench test --check docs-currency-workflow"}]}
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

The five native entries remain unobserved, acceptance-required manual diagnostics with exact canonical plan ownership.
They keep global Complete false and audit ManualRequired true.
Their future completion cannot be replaced by a Go declaration, excerpt hash, or an authored field.

```bench-planning-facts
{"version":1,"authority":"authored-spec","completeness":"declared-complete","baseline":"source-base","scope":[{"id":"implementation-scope","location":{"path":"specs/preflight-spec-staleness/spec.md","line":467,"needle":"## Ownership fences","source":"spec"},"paths":[".agents/commands/bench-implement-spec.md",".agents/skills/bench-implement-spec/references/staleness-pass.md","CHANGELOG.md","CONTEXT.md","cmd/bench/command_registry.go","cmd/bench/command_registry_test.go","cmd/bench/help_inventory_test.go","cmd/bench/preflight_version_test.go","internal/anchors/registry_chunk_chain.go","internal/anchors/registry_chunk_chain_test.go","internal/anchors/registry_commitment.go","internal/anchors/registry_data.go","internal/anchors/registry_data_test.go","internal/anchors/registry_debug_loop.go","internal/anchors/registry_decision_maps.go","internal/anchors/registry_decision_maps_test.go","internal/anchors/registry_ft311_preparation.go","internal/anchors/registry_retained_workflow.go","internal/conformance/axi_query_registry_test.go","internal/conformance/subcommand_routing_table_test.go","internal/preflight/binary_seal_test.go","internal/preflight/charge_pack.go","internal/preflight/charge_test.go","internal/preflight/command.go","internal/preflight/command_build_test.go","internal/preflight/completion_plan_test.go","internal/preflight/decision.go","internal/preflight/decision_test.go","internal/preflight/evidencecmd/evidence_grammar_test.go","internal/preflight/evidencecmd/operations.go","internal/preflight/evidencecmd/operations_test.go","internal/preflight/explicit_base_test.go","internal/preflight/gather.go","internal/preflight/gather_inputs.go","internal/preflight/gather_test.go","internal/preflight/preflighttest/fixture.go","internal/preflight/source_tip_test.go","internal/preflight/staleness/","internal/preflight/verdict_summary_test.go","tests/canary/docs-currency-token-diet/signal-vocabulary-drift","tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns","tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary","tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary","tests/canary/workflow-guidance-anchors/context-coverage-map-term","tests/canary/workflow-guidance-anchors/context-coverage-row-parts","tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary","tests/canary/workflow-guidance-anchors/context-decision-map-term","tests/canary/workflow-guidance-anchors/context-reader-sweep-term","tests/canary/workflow-guidance-anchors/context-ticket-vocabulary","tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration","tests/canary/workflow-guidance-anchors/delegated-entry-refusals","tests/canary/workflow-guidance-anchors/delegated-resumption-contents","tests/canary/workflow-guidance-anchors/dg-25","tests/canary/workflow-guidance-anchors/dg-26","tests/canary/workflow-guidance-anchors/dg-29","tests/canary/workflow-guidance-anchors/dg-29-verification-target","tests/canary/workflow-guidance-anchors/dg-30","tests/canary/workflow-guidance-anchors/dg-31","tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger","tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness","tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding","tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer","tests/canary/workflow-guidance-anchors/implement-spec-entry-validation","tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness","tests/canary/workflow-guidance-anchors/implement-spec-inline-exception","tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor","tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer","tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper","tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route","tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted","tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order","tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed","tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor","tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight","tests/canary/workflow-guidance-anchors/implement-spec-write-delegation","tests/canary/workflow-guidance-anchors/line-anchor-missing","tests/canary/workflow-guidance-anchors/prepared-build-approval","tests/canary/workflow-guidance-anchors/prepared-build-freshness"]}],"changes":[{"id":"operation-registry-change","location":{"path":"specs/preflight-spec-staleness/spec.md","line":193,"needle":"### Applicability and CLI","source":"spec"},"kind":"grammar","owner":{"path":"internal/preflight/evidencecmd/operations.go","package":"github.com/gibbonmi/bench/internal/preflight/evidencecmd","symbol":"evidencecmd.operations","kind":"var","signature":"var operations = []Operation{\n\t{Mode: ModeReview, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,\n\t\tdescription: \"review-entry checks that a spec's artifacts agree with the tree, one count line then the red checks only\"},\n\t{Mode: ModeReview, selectors: []string{flagCharge}, required: []string{FlagBase, FlagTip}, optional: []string{flagQuota}, Kind: KindPrepareReviewEvidence, Bounded: true,\n\t\tdescription: \"prepare one immutable review evidence artifact and print its bounded orientation\"},\n\t{Mode: ModeBuild, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,\n\t\tdescription: \"build-entry checks that a spec's artifacts agree with the tree, one count line then the red checks only\"},\n\t{Mode: ModeBuild, selectors: []string{flagPlanOnly}, optional: []string{FlagBase, FlagTip}, Kind: KindPlanOnly, Bounded: true,\n\t\tdescription: \"validate the authored spec and tickets without delivery admission or a build charge\"},\n\t{Mode: ModeBuild, selectors: []string{flagCharge}, required: []string{FlagTicket, FlagBase, FlagTip}, optional: []string{flagQuota}, Kind: KindPrepareEvidence, Bounded: true,\n\t\tdescription: \"prepare one immutable build evidence artifact and print its bounded orientation\"},\n\t{Mode: ModeBuild, selectors: []string{flagPropose}, required: []string{FlagTicket, FlagBase, FlagTip}, Kind: KindProposal,\n\t\tdescription: \"propose one ticket's Writes: entries from the pinned source\"},\n\t{Mode: modeEvidence, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,\n\t\tdescription: \"print the summary of a prepared evidence artifact, or one bounded fragment at a cursor, and its exact successor\"},\n\t{Mode: modeEvidence, selectors: []string{flagSource}, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,\n\t\tdescription: \"print one bounded fragment of one declared source stream and its exact successor\"},\n\t{Mode: modeEvidence, selectors: []string{flagVerify}, Kind: KindVerifyEvidence, Bounded: true,\n\t\tdescription: \"verify every stored page and source digest of a prepared evidence artifact\"},\n\t{Mode: modeEvidence, selectors: []string{flagCurrent}, Kind: KindCurrentEvidence, Bounded: true,\n\t\tdescription: \"bind a prepared evidence artifact to the current assignment and source pair\"},\n\t{Mode: modeEvidence, selectors: []string{FlagTo}, Kind: KindExportEvidence, Bounded: true,\n\t\tdescription: \"export every verified source of a prepared evidence artifact to its own file in an absent or empty directory\"},\n\t{Mode: ModeClean, optional: []string{flagCursor}, Kind: KindCleanPlan, Bounded: true,\n\t\tdescription: \"print one bounded page of the exact evidence deletion targets and its fingerprint\"},\n\t{Mode: ModeClean, selectors: []string{flagApply}, Kind: KindCleanApply, Bounded: true,\n\t\tdescription: \"delete exactly the targets one fingerprinted cleanup plan named\"},\n}"},"caller_scope":"production-and-tests","state":"existing","destination":null,"dispositions":[]}],"premises":[{"id":"current-preflight-GatherPinned","location":{"path":"specs/preflight-spec-staleness/spec.md","line":598,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/preflight/gather_inputs.go","package":"github.com/gibbonmi/bench/internal/preflight","symbol":"preflight.GatherPinned","kind":"func","signature":"func GatherPinned(root, mode, slug, explicitBase, sourceTipPin string) (Facts, *BootstrapFailure)"},"expected":{"kind":"func","signature":"func GatherPinned(root, mode, slug, explicitBase, sourceTipPin string) (Facts, *BootstrapFailure)"}},{"id":"current-preflight-Decide","location":{"path":"specs/preflight-spec-staleness/spec.md","line":598,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/preflight/decision.go","package":"github.com/gibbonmi/bench/internal/preflight","symbol":"preflight.Decide","kind":"func","signature":"func Decide(f Facts) Verdict"},"expected":{"kind":"func","signature":"func Decide(f Facts) Verdict"}},{"id":"current-coverage-ParseSpec","location":{"path":"specs/preflight-spec-staleness/spec.md","line":598,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/coverage/coverage.go","package":"github.com/gibbonmi/bench/internal/coverage","symbol":"coverage.ParseSpec","kind":"func","signature":"func ParseSpec(path string) (optIn bool, ids []string, violations []string, err error)"},"expected":{"kind":"func","signature":"func ParseSpec(path string) (optIn bool, ids []string, violations []string, err error)"}},{"id":"current-tickets-ParseTicket","location":{"path":"specs/preflight-spec-staleness/spec.md","line":598,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/tickets/tickets.go","package":"github.com/gibbonmi/bench/internal/tickets","symbol":"tickets.ParseTicket","kind":"func","signature":"func ParseTicket(name string, content []byte, siblings []string, tag string) (Ticket, []string)"},"expected":{"kind":"func","signature":"func ParseTicket(name string, content []byte, siblings []string, tag string) (Ticket, []string)"}},{"id":"current-consumers-Resolve","location":{"path":"specs/preflight-spec-staleness/spec.md","line":598,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/consumers/resolve.go","package":"github.com/gibbonmi/bench/internal/consumers","symbol":"consumers.Resolve","kind":"func","signature":"func Resolve(pkgs []*Package, query string) ([]Match, error)"},"expected":{"kind":"func","signature":"func Resolve(pkgs []*Package, query string) ([]Match, error)"}},{"id":"current-reviewrecord-ReadPlan","location":{"path":"specs/preflight-spec-staleness/spec.md","line":598,"needle":"## Source evidence and limits","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/reviewrecord/plan.go","package":"github.com/gibbonmi/bench/internal/reviewrecord","symbol":"reviewrecord.ReadPlan","kind":"func","signature":"func ReadPlan(root, tree, spec string) (Plan, error)"},"expected":{"kind":"func","signature":"func ReadPlan(root, tree, spec string) (Plan, error)"}},{"id":"current-diff-MovementCheckedRetry","location":{"path":"specs/preflight-spec-staleness/spec.md","line":583,"needle":"MovementCheckedRetry resolve in their current owners.","source":"spec"},"kind":"symbol","state":"current","subject":{"path":"internal/diff/range.go","package":"github.com/gibbonmi/bench/internal/diff","symbol":"diff.MovementCheckedRetry","kind":"func","signature":"func MovementCheckedRetry(root string, read func(MovementSnapshot) (kind, hint string)) MovementResult"},"expected":{"kind":"func","signature":"func MovementCheckedRetry(root string, read func(MovementSnapshot) (kind, hint string)) MovementResult"}}],"seams":[{"id":"planned-TestDriftObjects","location":{"path":"specs/preflight-spec-staleness/spec.md","line":391,"needle":"| SS1 | 1 | A modified fenced path appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A changed blob defeats empty enumeration. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestDriftObjects","state":"planned","rows":["SS1","SS2","SS3","SS4","SS5","SS6","SS10","SS11"]},{"id":"planned-TestDriftIdentity","location":{"path":"specs/preflight-spec-staleness/spec.md","line":397,"needle":"| SS7 | 2 | The audit prints exact resolved baseline and tip identities. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An ambient branch spelling cannot identify the frozen subject. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestDriftIdentity","state":"planned","rows":["SS7","SS8","SS9","SS12"]},{"id":"planned-TestSymbolIdentity","location":{"path":"specs/preflight-spec-staleness/spec.md","line":403,"needle":"| SS13 | 3 | An absent cited symbol reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Removing the actual declaration defeats name retention. |","source":"spec"},"path":"internal/preflight/staleness/identity_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestSymbolIdentity","state":"planned","rows":["SS13","SS14","SS15","SS16","SS17","SS18","SS19"]},{"id":"planned-TestTestIdentity","location":{"path":"specs/preflight-spec-staleness/spec.md","line":410,"needle":"| SS20 | 5 | An absent existing test file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Deleting the file cannot leave green evidence. |","source":"spec"},"path":"internal/preflight/staleness/identity_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestTestIdentity","state":"planned","rows":["SS20","SS21","SS22","SS23","SS24","SS25","SS26","SS27"]},{"id":"planned-TestCoverageIdentity","location":{"path":"specs/preflight-spec-staleness/spec.md","line":418,"needle":"| SS28 | 7 | An unowned coverage row retains red rows-owned. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Source green cannot mask missing ticket ownership. |","source":"spec"},"path":"internal/preflight/staleness/identity_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestCoverageIdentity","state":"planned","rows":["SS28","SS29","SS30","SS31"]},{"id":"planned-TestAuditCompleteness","location":{"path":"specs/preflight-spec-staleness/spec.md","line":422,"needle":"| SS32 | 8 | Missing planning metadata reports incomplete audit evidence. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Zero parsed items cannot mean complete green. |","source":"spec"},"path":"internal/preflight/staleness/identity_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditCompleteness","state":"planned","rows":["SS32","SS33","SS34","SS49","SS50","SS57","SS58"]},{"id":"planned-TestAuditPrecedence","location":{"path":"specs/preflight-spec-staleness/spec.md","line":425,"needle":"| SS35 | 12 | Dirty source refuses before a staleness decision. | planned internal/preflight/staleness/drift_test.go (TestAuditPrecedence), actual preflight command | A competing missing-symbol error cannot replace source authority. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditPrecedence","state":"planned","rows":["SS35"]},{"id":"planned-TestAuditMovement","location":{"path":"specs/preflight-spec-staleness/spec.md","line":426,"needle":"| SS36 | 12 | Persistent snapshot movement publishes no skip decision. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A second movement defeats mixed before/after proof. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditMovement","state":"planned","rows":["SS36","SS37"]},{"id":"planned-TestAuditTrust","location":{"path":"specs/preflight-spec-staleness/spec.md","line":428,"needle":"| SS38 | 16 | A required executable-seal refusal remains binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | Green static facts cannot launch an untrusted binary. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditTrust","state":"planned","rows":["SS38","SS56"]},{"id":"planned-TestAuditGrammar","location":{"path":"specs/preflight-spec-staleness/spec.md","line":429,"needle":"| SS39 | 15 | Audit-spec rejects combination with another selector. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | An audit plus charge request cannot choose one silently. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditGrammar","state":"planned","rows":["SS39","SS40"]},{"id":"planned-TestAuditDecision","location":{"path":"specs/preflight-spec-staleness/spec.md","line":431,"needle":"| SS41 | 9 | Complete clean supported evidence prints manual_required false. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | The actual conjunction is the only omission producer. |","source":"spec"},"path":"internal/preflight/staleness/decision_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditDecision","state":"planned","rows":["SS41","SS42","SS43","SS44"]},{"id":"planned-TestAuditApplicability","location":{"path":"specs/preflight-spec-staleness/spec.md","line":435,"needle":"| SS45 | 11 | Ordinary after-ticket preflight keeps initial test-audit inapplicability. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | An implemented predecessor test cannot block its successor. |","source":"spec"},"path":"internal/preflight/staleness/decision_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditApplicability","state":"planned","rows":["SS45","SS46"]},{"id":"planned-TestAuditSafety","location":{"path":"specs/preflight-spec-staleness/spec.md","line":441,"needle":"| SS51 | 12 | A special before/after source refuses without following it. | planned internal/preflight/staleness/drift_test.go (TestAuditSafety), actual preflight command | An outside source sentinel detects a followed link. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditSafety","state":"planned","rows":["SS51"]},{"id":"planned-TestAuditBound","location":{"path":"specs/preflight-spec-staleness/spec.md","line":442,"needle":"| SS52 | 15 | An audit response retains all drift rows through existing spill. | planned internal/preflight/staleness/drift_test.go (TestAuditBound), actual preflight command | Full reconstruction catches truncation hidden by bounded output. |","source":"spec"},"path":"internal/preflight/staleness/drift_test.go","package":"github.com/gibbonmi/bench/internal/preflight/staleness","function":"TestAuditBound","state":"planned","rows":["SS52"]}],"decisions":[{"id":"acceptance-SS1","location":{"path":"specs/preflight-spec-staleness/spec.md","line":391,"needle":"| SS1 | 1 | A modified fenced path appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A changed blob defeats empty enumeration. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":391,"needle":"| SS1 | 1 | A modified fenced path appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A changed blob defeats empty enumeration. |","source":"tip"},"predicate":"A modified fenced path appears in drift.","row":"SS1"},{"id":"acceptance-SS2","location":{"path":"specs/preflight-spec-staleness/spec.md","line":392,"needle":"| SS2 | 1 | A deleted fenced path appears with its before identity. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | After-only listing cannot lose removal. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":392,"needle":"| SS2 | 1 | A deleted fenced path appears with its before identity. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | After-only listing cannot lose removal. |","source":"tip"},"predicate":"A deleted fenced path appears with its before identity.","row":"SS2"},{"id":"acceptance-SS3","location":{"path":"specs/preflight-spec-staleness/spec.md","line":393,"needle":"| SS3 | 1 | A new fenced path appears with its after identity. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Before-only listing cannot lose addition. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":393,"needle":"| SS3 | 1 | A new fenced path appears with its after identity. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Before-only listing cannot lose addition. |","source":"tip"},"predicate":"A new fenced path appears with its after identity.","row":"SS3"},{"id":"acceptance-SS4","location":{"path":"specs/preflight-spec-staleness/spec.md","line":394,"needle":"| SS4 | 1 | A renamed fenced path retains both names. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A move cannot appear as unchanged source. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":394,"needle":"| SS4 | 1 | A renamed fenced path retains both names. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A move cannot appear as unchanged source. |","source":"tip"},"predicate":"A renamed fenced path retains both names.","row":"SS4"},{"id":"acceptance-SS5","location":{"path":"specs/preflight-spec-staleness/spec.md","line":395,"needle":"| SS5 | 1 | A file-mode change appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Byte equality cannot hide changed executable mode. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":395,"needle":"| SS5 | 1 | A file-mode change appears in drift. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Byte equality cannot hide changed executable mode. |","source":"tip"},"predicate":"A file-mode change appears in drift.","row":"SS5"},{"id":"acceptance-SS6","location":{"path":"specs/preflight-spec-staleness/spec.md","line":396,"needle":"| SS6 | 1 | A fence prefix includes each tracked dot-path member. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A hidden member catches incomplete scope enumeration. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":396,"needle":"| SS6 | 1 | A fence prefix includes each tracked dot-path member. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | A hidden member catches incomplete scope enumeration. |","source":"tip"},"predicate":"A fence prefix includes each tracked dot-path member.","row":"SS6"},{"id":"acceptance-SS7","location":{"path":"specs/preflight-spec-staleness/spec.md","line":397,"needle":"| SS7 | 2 | The audit prints exact resolved baseline and tip identities. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An ambient branch spelling cannot identify the frozen subject. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":397,"needle":"| SS7 | 2 | The audit prints exact resolved baseline and tip identities. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An ambient branch spelling cannot identify the frozen subject. |","source":"tip"},"predicate":"The audit prints exact resolved baseline and tip identities.","row":"SS7"},{"id":"acceptance-SS8","location":{"path":"specs/preflight-spec-staleness/spec.md","line":398,"needle":"| SS8 | 2 | An absent before spec refuses audit collection. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | Missing source cannot become a zero drift set. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":398,"needle":"| SS8 | 2 | An absent before spec refuses audit collection. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | Missing source cannot become a zero drift set. |","source":"tip"},"predicate":"An absent before spec refuses audit collection.","row":"SS8"},{"id":"acceptance-SS9","location":{"path":"specs/preflight-spec-staleness/spec.md","line":399,"needle":"| SS9 | 2 | Changed active spec bytes appear as drift. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | A prose correction cannot claim a semantic digest exemption. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":399,"needle":"| SS9 | 2 | Changed active spec bytes appear as drift. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | A prose correction cannot claim a semantic digest exemption. |","source":"tip"},"predicate":"Changed active spec bytes appear as drift.","row":"SS9"},{"id":"acceptance-SS10","location":{"path":"specs/preflight-spec-staleness/spec.md","line":400,"needle":"| SS10 | 13 | An unchanged exact member set prints a complete empty drift table. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Silence or incomplete enumeration cannot pass. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":400,"needle":"| SS10 | 13 | An unchanged exact member set prints a complete empty drift table. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | Silence or incomplete enumeration cannot pass. |","source":"tip"},"predicate":"An unchanged exact member set prints a complete empty drift table.","row":"SS10"},{"id":"acceptance-SS11","location":{"path":"specs/preflight-spec-staleness/spec.md","line":401,"needle":"| SS11 | 1 | A segment-prefix collision stays outside the fenced scope. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | An x2 source cannot be confused with x. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":401,"needle":"| SS11 | 1 | A segment-prefix collision stays outside the fenced scope. | planned internal/preflight/staleness/drift_test.go (TestDriftObjects), actual preflight command | An x2 source cannot be confused with x. |","source":"tip"},"predicate":"A segment-prefix collision stays outside the fenced scope.","row":"SS11"},{"id":"acceptance-SS12","location":{"path":"specs/preflight-spec-staleness/spec.md","line":402,"needle":"| SS12 | 2 | A claim source outside the write fence joins audit scope. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An unchanged writes set cannot hide a changed premise owner. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":402,"needle":"| SS12 | 2 | A claim source outside the write fence joins audit scope. | planned internal/preflight/staleness/drift_test.go (TestDriftIdentity), actual preflight command | An unchanged writes set cannot hide a changed premise owner. |","source":"tip"},"predicate":"A claim source outside the write fence joins audit scope.","row":"SS12"},{"id":"acceptance-SS13","location":{"path":"specs/preflight-spec-staleness/spec.md","line":403,"needle":"| SS13 | 3 | An absent cited symbol reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Removing the actual declaration defeats name retention. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":403,"needle":"| SS13 | 3 | An absent cited symbol reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Removing the actual declaration defeats name retention. |","source":"tip"},"predicate":"An absent cited symbol reaches a red cited-symbols result.","row":"SS13"},{"id":"acceptance-SS14","location":{"path":"specs/preflight-spec-staleness/spec.md","line":404,"needle":"| SS14 | 4 | A cited signature change reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Same spelling cannot hide a changed seam. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":404,"needle":"| SS14 | 4 | A cited signature change reaches a red cited-symbols result. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Same spelling cannot hide a changed seam. |","source":"tip"},"predicate":"A cited signature change reaches a red cited-symbols result.","row":"SS14"},{"id":"acceptance-SS15","location":{"path":"specs/preflight-spec-staleness/spec.md","line":405,"needle":"| SS15 | 4 | A cited symbol moved to another owner reports that discrepancy. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Name-only resolution cannot establish owner identity. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":405,"needle":"| SS15 | 4 | A cited symbol moved to another owner reports that discrepancy. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Name-only resolution cannot establish owner identity. |","source":"tip"},"predicate":"A cited symbol moved to another owner reports that discrepancy.","row":"SS15"},{"id":"acceptance-SS16","location":{"path":"specs/preflight-spec-staleness/spec.md","line":406,"needle":"| SS16 | 4 | A wrong declaration kind reports its exact mismatch. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | A var replacing a const cannot preserve the claimed rule. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":406,"needle":"| SS16 | 4 | A wrong declaration kind reports its exact mismatch. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | A var replacing a const cannot preserve the claimed rule. |","source":"tip"},"predicate":"A wrong declaration kind reports its exact mismatch.","row":"SS16"},{"id":"acceptance-SS17","location":{"path":"specs/preflight-spec-staleness/spec.md","line":407,"needle":"| SS17 | 8 | An ambiguous bare symbol remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Selecting one candidate cannot yield complete proof. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":407,"needle":"| SS17 | 8 | An ambiguous bare symbol remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Selecting one candidate cannot yield complete proof. |","source":"tip"},"predicate":"An ambiguous bare symbol remains diagnostic.","row":"SS17"},{"id":"acceptance-SS18","location":{"path":"specs/preflight-spec-staleness/spec.md","line":408,"needle":"| SS18 | 3 | An alias resolves through the canonical origin rule. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Independent lexical lookup would disagree with consumers. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":408,"needle":"| SS18 | 3 | An alias resolves through the canonical origin rule. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | Independent lexical lookup would disagree with consumers. |","source":"tip"},"predicate":"An alias resolves through the canonical origin rule.","row":"SS18"},{"id":"acceptance-SS19","location":{"path":"specs/preflight-spec-staleness/spec.md","line":409,"needle":"| SS19 | 8 | A non-default context claim remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | A default loader cannot certify another platform. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":409,"needle":"| SS19 | 8 | A non-default context claim remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestSymbolIdentity), actual preflight command | A default loader cannot certify another platform. |","source":"tip"},"predicate":"A non-default context claim remains diagnostic.","row":"SS19"},{"id":"acceptance-SS20","location":{"path":"specs/preflight-spec-staleness/spec.md","line":410,"needle":"| SS20 | 5 | An absent existing test file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Deleting the file cannot leave green evidence. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":410,"needle":"| SS20 | 5 | An absent existing test file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Deleting the file cannot leave green evidence. |","source":"tip"},"predicate":"An absent existing test file reaches red cited-tests.","row":"SS20"},{"id":"acceptance-SS21","location":{"path":"specs/preflight-spec-staleness/spec.md","line":411,"needle":"| SS21 | 5 | A removed function in an existing file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | File existence cannot stand in for the named test. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":411,"needle":"| SS21 | 5 | A removed function in an existing file reaches red cited-tests. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | File existence cannot stand in for the named test. |","source":"tip"},"predicate":"A removed function in an existing file reaches red cited-tests.","row":"SS21"},{"id":"acceptance-SS22","location":{"path":"specs/preflight-spec-staleness/spec.md","line":412,"needle":"| SS22 | 6 | An absent planned test function remains valid future work. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | No first checkpoint requires its new test to exist. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":412,"needle":"| SS22 | 6 | An absent planned test function remains valid future work. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | No first checkpoint requires its new test to exist. |","source":"tip"},"predicate":"An absent planned test function remains valid future work.","row":"SS22"},{"id":"acceptance-SS23","location":{"path":"specs/preflight-spec-staleness/spec.md","line":413,"needle":"| SS23 | 6 | An existing supposedly planned function reports already-shipped drift. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Future labeling cannot hide landed work. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":413,"needle":"| SS23 | 6 | An existing supposedly planned function reports already-shipped drift. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Future labeling cannot hide landed work. |","source":"tip"},"predicate":"An existing supposedly planned function reports already-shipped drift.","row":"SS23"},{"id":"acceptance-SS24","location":{"path":"specs/preflight-spec-staleness/spec.md","line":414,"needle":"| SS24 | 6 | Other functions in a planned test's existing file remain allowed. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Whole-file absence would make a legal slice impossible. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":414,"needle":"| SS24 | 6 | Other functions in a planned test's existing file remain allowed. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Whole-file absence would make a legal slice impossible. |","source":"tip"},"predicate":"Other functions in a planned test's existing file remain allowed.","row":"SS24"},{"id":"acceptance-SS25","location":{"path":"specs/preflight-spec-staleness/spec.md","line":415,"needle":"| SS25 | 5 | A named supported subtest uses the canonical citation grammar. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | A second subtest parser cannot supply inconsistent identity. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":415,"needle":"| SS25 | 5 | A named supported subtest uses the canonical citation grammar. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | A second subtest parser cannot supply inconsistent identity. |","source":"tip"},"predicate":"A named supported subtest uses the canonical citation grammar.","row":"SS25"},{"id":"acceptance-SS26","location":{"path":"specs/preflight-spec-staleness/spec.md","line":416,"needle":"| SS26 | 8 | An unsupported test citation remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Unknown syntax cannot produce an empty complete result. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":416,"needle":"| SS26 | 8 | An unsupported test citation remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | Unknown syntax cannot produce an empty complete result. |","source":"tip"},"predicate":"An unsupported test citation remains diagnostic.","row":"SS26"},{"id":"acceptance-SS27","location":{"path":"specs/preflight-spec-staleness/spec.md","line":417,"needle":"| SS27 | 5 | Wrong-context test evidence remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | A declaration outside the graded context is not current evidence. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":417,"needle":"| SS27 | 5 | Wrong-context test evidence remains diagnostic. | planned internal/preflight/staleness/identity_test.go (TestTestIdentity), actual preflight command | A declaration outside the graded context is not current evidence. |","source":"tip"},"predicate":"Wrong-context test evidence remains diagnostic.","row":"SS27"},{"id":"acceptance-SS28","location":{"path":"specs/preflight-spec-staleness/spec.md","line":418,"needle":"| SS28 | 7 | An unowned coverage row retains red rows-owned. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Source green cannot mask missing ticket ownership. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":418,"needle":"| SS28 | 7 | An unowned coverage row retains red rows-owned. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Source green cannot mask missing ticket ownership. |","source":"tip"},"predicate":"An unowned coverage row retains red rows-owned.","row":"SS28"},{"id":"acceptance-SS29","location":{"path":"specs/preflight-spec-staleness/spec.md","line":419,"needle":"| SS29 | 7 | A phantom Covers token retains red rows-membership. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Audit tolerance cannot admit an invented row. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":419,"needle":"| SS29 | 7 | A phantom Covers token retains red rows-membership. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Audit tolerance cannot admit an invented row. |","source":"tip"},"predicate":"A phantom Covers token retains red rows-membership.","row":"SS29"},{"id":"acceptance-SS30","location":{"path":"specs/preflight-spec-staleness/spec.md","line":420,"needle":"| SS30 | 7 | Changed coverage predicate bytes join drift. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Stable row IDs cannot hide changed acceptance. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":420,"needle":"| SS30 | 7 | Changed coverage predicate bytes join drift. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Stable row IDs cannot hide changed acceptance. |","source":"tip"},"predicate":"Changed coverage predicate bytes join drift.","row":"SS30"},{"id":"acceptance-SS31","location":{"path":"specs/preflight-spec-staleness/spec.md","line":421,"needle":"| SS31 | 7 | Changed ticket Covers membership joins drift. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Unchanged tests cannot hide a different row owner. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":421,"needle":"| SS31 | 7 | Changed ticket Covers membership joins drift. | planned internal/preflight/staleness/identity_test.go (TestCoverageIdentity), actual preflight command | Unchanged tests cannot hide a different row owner. |","source":"tip"},"predicate":"Changed ticket Covers membership joins drift.","row":"SS31"},{"id":"acceptance-SS32","location":{"path":"specs/preflight-spec-staleness/spec.md","line":422,"needle":"| SS32 | 8 | Missing planning metadata reports incomplete audit evidence. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Zero parsed items cannot mean complete green. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":422,"needle":"| SS32 | 8 | Missing planning metadata reports incomplete audit evidence. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Zero parsed items cannot mean complete green. |","source":"tip"},"predicate":"Missing planning metadata reports incomplete audit evidence.","row":"SS32"},{"id":"acceptance-SS33","location":{"path":"specs/preflight-spec-staleness/spec.md","line":423,"needle":"| SS33 | 8 | Partial declared metadata keeps complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | An author flag cannot bless omitted inventory. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":423,"needle":"| SS33 | 8 | Partial declared metadata keeps complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | An author flag cannot bless omitted inventory. |","source":"tip"},"predicate":"Partial declared metadata keeps complete false.","row":"SS33"},{"id":"acceptance-SS34","location":{"path":"specs/preflight-spec-staleness/spec.md","line":424,"needle":"| SS34 | 8 | An omitted supported citation keeps complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | The second citation defeats claimed completeness. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":424,"needle":"| SS34 | 8 | An omitted supported citation keeps complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | The second citation defeats claimed completeness. |","source":"tip"},"predicate":"An omitted supported citation keeps complete false.","row":"SS34"},{"id":"acceptance-SS35","location":{"path":"specs/preflight-spec-staleness/spec.md","line":425,"needle":"| SS35 | 12 | Dirty source refuses before a staleness decision. | planned internal/preflight/staleness/drift_test.go (TestAuditPrecedence), actual preflight command | A competing missing-symbol error cannot replace source authority. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":425,"needle":"| SS35 | 12 | Dirty source refuses before a staleness decision. | planned internal/preflight/staleness/drift_test.go (TestAuditPrecedence), actual preflight command | A competing missing-symbol error cannot replace source authority. |","source":"tip"},"predicate":"Dirty source refuses before a staleness decision.","row":"SS35"},{"id":"acceptance-SS36","location":{"path":"specs/preflight-spec-staleness/spec.md","line":426,"needle":"| SS36 | 12 | Persistent snapshot movement publishes no skip decision. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A second movement defeats mixed before/after proof. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":426,"needle":"| SS36 | 12 | Persistent snapshot movement publishes no skip decision. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A second movement defeats mixed before/after proof. |","source":"tip"},"predicate":"Persistent snapshot movement publishes no skip decision.","row":"SS36"},{"id":"acceptance-SS37","location":{"path":"specs/preflight-spec-staleness/spec.md","line":427,"needle":"| SS37 | 16 | A source-tip mismatch retains its current refusal. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A stale pin cannot authenticate the subject. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":427,"needle":"| SS37 | 16 | A source-tip mismatch retains its current refusal. | planned internal/preflight/staleness/drift_test.go (TestAuditMovement), actual preflight command | A stale pin cannot authenticate the subject. |","source":"tip"},"predicate":"A source-tip mismatch retains its current refusal.","row":"SS37"},{"id":"acceptance-SS38","location":{"path":"specs/preflight-spec-staleness/spec.md","line":428,"needle":"| SS38 | 16 | A required executable-seal refusal remains binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | Green static facts cannot launch an untrusted binary. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":428,"needle":"| SS38 | 16 | A required executable-seal refusal remains binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | Green static facts cannot launch an untrusted binary. |","source":"tip"},"predicate":"A required executable-seal refusal remains binding.","row":"SS38"},{"id":"acceptance-SS39","location":{"path":"specs/preflight-spec-staleness/spec.md","line":429,"needle":"| SS39 | 15 | Audit-spec rejects combination with another selector. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | An audit plus charge request cannot choose one silently. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":429,"needle":"| SS39 | 15 | Audit-spec rejects combination with another selector. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | An audit plus charge request cannot choose one silently. |","source":"tip"},"predicate":"Audit-spec rejects combination with another selector.","row":"SS39"},{"id":"acceptance-SS40","location":{"path":"specs/preflight-spec-staleness/spec.md","line":430,"needle":"| SS40 | 15 | Audit-spec requires both explicit source identities. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | A missing base cannot imply an authoritative baseline. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":430,"needle":"| SS40 | 15 | Audit-spec requires both explicit source identities. | planned internal/preflight/staleness/drift_test.go (TestAuditGrammar), actual preflight command | A missing base cannot imply an authoritative baseline. |","source":"tip"},"predicate":"Audit-spec requires both explicit source identities.","row":"SS40"},{"id":"acceptance-SS41","location":{"path":"specs/preflight-spec-staleness/spec.md","line":431,"needle":"| SS41 | 9 | Complete clean supported evidence prints manual_required false. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | The actual conjunction is the only omission producer. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":431,"needle":"| SS41 | 9 | Complete clean supported evidence prints manual_required false. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | The actual conjunction is the only omission producer. |","source":"tip"},"predicate":"Complete clean supported evidence prints manual_required false.","row":"SS41"},{"id":"acceptance-SS42","location":{"path":"specs/preflight-spec-staleness/spec.md","line":432,"needle":"| SS42 | 9 | Any supported audit red prints manual_required true. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | An empty drift set cannot override missing evidence. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":432,"needle":"| SS42 | 9 | Any supported audit red prints manual_required true. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | An empty drift set cannot override missing evidence. |","source":"tip"},"predicate":"Any supported audit red prints manual_required true.","row":"SS42"},{"id":"acceptance-SS43","location":{"path":"specs/preflight-spec-staleness/spec.md","line":433,"needle":"| SS43 | 9 | Diagnostic-only exit zero prints manual_required true. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | Success of the read cannot mean success of complete proof. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":433,"needle":"| SS43 | 9 | Diagnostic-only exit zero prints manual_required true. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | Success of the read cannot mean success of complete proof. |","source":"tip"},"predicate":"Diagnostic-only exit zero prints manual_required true.","row":"SS43"},{"id":"acceptance-SS44","location":{"path":"specs/preflight-spec-staleness/spec.md","line":434,"needle":"| SS44 | 9 | A not-applicable required build row prevents omission. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | Missing implementation tickets cannot enter the fast route. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":434,"needle":"| SS44 | 9 | A not-applicable required build row prevents omission. | planned internal/preflight/staleness/decision_test.go (TestAuditDecision), actual preflight command | Missing implementation tickets cannot enter the fast route. |","source":"tip"},"predicate":"A not-applicable required build row prevents omission.","row":"SS44"},{"id":"acceptance-SS45","location":{"path":"specs/preflight-spec-staleness/spec.md","line":435,"needle":"| SS45 | 11 | Ordinary after-ticket preflight keeps initial test-audit inapplicability. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | An implemented predecessor test cannot block its successor. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":435,"needle":"| SS45 | 11 | Ordinary after-ticket preflight keeps initial test-audit inapplicability. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | An implemented predecessor test cannot block its successor. |","source":"tip"},"predicate":"Ordinary after-ticket preflight keeps initial test-audit inapplicability.","row":"SS45"},{"id":"acceptance-SS46","location":{"path":"specs/preflight-spec-staleness/spec.md","line":436,"needle":"| SS46 | 11 | Charge preparation grants no manual-omission result. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | Continuation cannot inherit initial audit authority. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":436,"needle":"| SS46 | 11 | Charge preparation grants no manual-omission result. | planned internal/preflight/staleness/decision_test.go (TestAuditApplicability), actual preflight command | Continuation cannot inherit initial audit authority. |","source":"tip"},"predicate":"Charge preparation grants no manual-omission result.","row":"SS46"},{"id":"acceptance-SS49","location":{"path":"specs/preflight-spec-staleness/spec.md","line":439,"needle":"| SS49 | 8 | Duplicate metadata refuses complete audit collection. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | A first-fence choice cannot establish full inventory. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":439,"needle":"| SS49 | 8 | Duplicate metadata refuses complete audit collection. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | A first-fence choice cannot establish full inventory. |","source":"tip"},"predicate":"Duplicate metadata refuses complete audit collection.","row":"SS49"},{"id":"acceptance-SS50","location":{"path":"specs/preflight-spec-staleness/spec.md","line":440,"needle":"| SS50 | 13 | A present empty declared claim set differs from missing metadata. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | An empty syntax container cannot certify completeness by itself. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":440,"needle":"| SS50 | 13 | A present empty declared claim set differs from missing metadata. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | An empty syntax container cannot certify completeness by itself. |","source":"tip"},"predicate":"A present empty declared claim set differs from missing metadata.","row":"SS50"},{"id":"acceptance-SS51","location":{"path":"specs/preflight-spec-staleness/spec.md","line":441,"needle":"| SS51 | 12 | A special before/after source refuses without following it. | planned internal/preflight/staleness/drift_test.go (TestAuditSafety), actual preflight command | An outside source sentinel detects a followed link. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":441,"needle":"| SS51 | 12 | A special before/after source refuses without following it. | planned internal/preflight/staleness/drift_test.go (TestAuditSafety), actual preflight command | An outside source sentinel detects a followed link. |","source":"tip"},"predicate":"A special before/after source refuses without following it.","row":"SS51"},{"id":"acceptance-SS52","location":{"path":"specs/preflight-spec-staleness/spec.md","line":442,"needle":"| SS52 | 15 | An audit response retains all drift rows through existing spill. | planned internal/preflight/staleness/drift_test.go (TestAuditBound), actual preflight command | Full reconstruction catches truncation hidden by bounded output. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":442,"needle":"| SS52 | 15 | An audit response retains all drift rows through existing spill. | planned internal/preflight/staleness/drift_test.go (TestAuditBound), actual preflight command | Full reconstruction catches truncation hidden by bounded output. |","source":"tip"},"predicate":"An audit response retains all drift rows through existing spill.","row":"SS52"},{"id":"acceptance-SS56","location":{"path":"specs/preflight-spec-staleness/spec.md","line":446,"needle":"| SS56 | 16 | Existing required evidence bytes retain current-action binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | A source digest cannot replace delivered current context. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":446,"needle":"| SS56 | 16 | Existing required evidence bytes retain current-action binding. | planned internal/preflight/staleness/drift_test.go (TestAuditTrust), actual preflight command | A source digest cannot replace delivered current context. |","source":"tip"},"predicate":"Existing required evidence bytes retain current-action binding.","row":"SS56"},{"id":"acceptance-SS57","location":{"path":"specs/preflight-spec-staleness/spec.md","line":447,"needle":"| SS57 | 8 | A pending native-future manual entry keeps audit Complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Charge-ready current proof cannot certify the entire reviewed inventory mechanically. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":447,"needle":"| SS57 | 8 | A pending native-future manual entry keeps audit Complete false. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Charge-ready current proof cannot certify the entire reviewed inventory mechanically. |","source":"tip"},"predicate":"A pending native-future manual entry keeps audit Complete false.","row":"SS57"},{"id":"acceptance-SS58","location":{"path":"specs/preflight-spec-staleness/spec.md","line":448,"needle":"| SS58 | 8 | A pending native-future manual entry keeps manual_required true. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Diagnostic-only exit zero cannot turn future evidence applicability into a manual-audit skip. |","source":"spec"},"source":{"path":"specs/preflight-spec-staleness/spec.md","line":448,"needle":"| SS58 | 8 | A pending native-future manual entry keeps manual_required true. | planned internal/preflight/staleness/identity_test.go (TestAuditCompleteness), actual preflight command | Diagnostic-only exit zero cannot turn future evidence applicability into a manual-audit skip. |","source":"tip"},"predicate":"A pending native-future manual entry keeps manual_required true.","row":"SS58"}],"manual_claims":[{"id":"native-SS47","location":{"path":"specs/preflight-spec-staleness/spec.md","line":437,"needle":"| SS47 | 10 | The manual charge retains drift, reds, and unsupported locations. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | A narrowed charge cannot drop the remaining semantic work. |","source":"spec"},"reason":"Unobserved native workflow consumption remains acceptance-required after its owning implementation.","applicability":"native-future","row":"SS47","chunk":"SS-C3","ticket":"05-consume-exact-manual-omission.md","verification":"ss-native-diagnostic","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":23,"needle":"The reviewer approves the spec and the whole ticket graph once,","source":"tip"},"observation":"unobserved"},{"id":"native-SS48","location":{"path":"specs/preflight-spec-staleness/spec.md","line":438,"needle":"| SS48 | 14 | Workflow guidance retains reviewer approval before dispatch. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | A complete audit cannot become graph approval. |","source":"spec"},"reason":"Unobserved native workflow consumption remains acceptance-required after its owning implementation.","applicability":"native-future","row":"SS48","chunk":"SS-C3","ticket":"05-consume-exact-manual-omission.md","verification":"ss-native-diagnostic","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":23,"needle":"The reviewer approves the spec and the whole ticket graph once,","source":"tip"},"observation":"unobserved"},{"id":"native-SS53","location":{"path":"specs/preflight-spec-staleness/spec.md","line":443,"needle":"| SS53 | 10 | A fresh initial session follows the actual clean skip decision. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | The native session must omit manual dispatch only from complete proof. |","source":"spec"},"reason":"Unobserved native workflow consumption remains acceptance-required after its owning implementation.","applicability":"native-future","row":"SS53","chunk":"SS-C3","ticket":"05-consume-exact-manual-omission.md","verification":"ss-native-clean","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":23,"needle":"The reviewer approves the spec and the whole ticket graph once,","source":"tip"},"observation":"unobserved"},{"id":"native-SS54","location":{"path":"specs/preflight-spec-staleness/spec.md","line":444,"needle":"| SS54 | 10 | A fresh initial session still audits unsupported claims. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | The real diagnostic-only command cannot suppress manual dispatch. |","source":"spec"},"reason":"Unobserved native workflow consumption remains acceptance-required after its owning implementation.","applicability":"native-future","row":"SS54","chunk":"SS-C3","ticket":"05-consume-exact-manual-omission.md","verification":"ss-native-diagnostic","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":23,"needle":"The reviewer approves the spec and the whole ticket graph once,","source":"tip"},"observation":"unobserved"},{"id":"native-SS55","location":{"path":"specs/preflight-spec-staleness/spec.md","line":445,"needle":"| SS55 | 14 | A manual acceptance change retains the existing reviewer route. | planned fresh-session adoption and docs-currency-workflow through the actual audit-spec response | A mechanical gap cannot silently widen acceptance. |","source":"spec"},"reason":"Unobserved native workflow consumption remains acceptance-required after its owning implementation.","applicability":"native-future","row":"SS55","chunk":"SS-C3","ticket":"05-consume-exact-manual-omission.md","verification":"ss-native-diagnostic","phase_owner":{"path":".agents/commands/bench-implement-spec.md","line":23,"needle":"The reviewer approves the spec and the whole ticket graph once,","source":"tip"},"observation":"unobserved"}]}
```
