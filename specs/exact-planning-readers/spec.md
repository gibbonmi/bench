# Exact planning readers

Status: staged

Related roadmap: FT125.

Decision source: the named reviewed ready map `decisions/architecture-planning.md`, answer 4.

Verification log: 1 iteration(s) to accept — the independent ticket review accepted the graph without repair.
Earlier SPEC review found closure and related-identity issues; their frozen repairs were independently accepted.
The confirming SPEC and ticket reviews each returned zero Standards, Spec, and Coverage findings.
The preserved first draft is response artifact 1791348759131116188-a9e012873e65538a.out in the primary Bench response store.

## Problem

Current queries locate stories and declarations, but agents still slice whole files to read them.
Escaped evidence cells can exceed response limits even when the source looks small.
Worktree navigation also lacks a proven relation to a requested spec.

## Solution

Return exact file bytes, authored section spans, typed story obligations, and Go declaration spans.
Use the existing owners for source interpretation, output budgets, and tree authority.
Give every page a source identity, byte position, and executable continuation.
Refuse changed sources rather than joining different versions.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: existing owners provide interpretation and fixtures, while exact pagination requires careful byte accounting and integration.
Harder chunk: ER-C1 uses gpt-5.6-sol / xhigh for the reversible page protocol and source-change refusal.
The requested gpt-6.1-sol / high author line does not change other approved implementation lines.

### Read exact text

1. As an agent, I want a native file reader, so that the original bytes reach me without whole-file shell slicing.
2. As an agent, I want a literal H2 section reader, so that skills and commands can be read at their owning heading.
3. As an agent, I want stories, coverage, and status selections, so that a live spec returns its actual authored content.
4. As an agent, I want one story with typed coverage rows and seams, so that its obligations remain joined by the existing parser.
5. As an agent, I want an absent story distinguished from an empty selection, so that absence cannot look like successful authored evidence.
### Read declarations

6. As an agent, I want a Go declaration with its attached doc comment, so that I can inspect the complete production seam.
7. As an agent, I want adjacent declarations excluded, so that the selected body remains precise.
8. As an agent, I want ambiguous symbol matches refused with exact candidates, so that a bare name cannot choose a different seam.
9. As an agent, I want malformed Go refused, so that regex guesses cannot supply a declaration span.
### Read pages

10. As an agent, I want every page bound to one source and byte position, so that continuation cannot mix source versions.
11. As an agent, I want the complete continuation grammar in bounded help, so that the next read remains directly executable.
12. As an agent, I want escaped output budgeted after rendering, so that control characters cannot force a hidden spill.
13. As an agent, I want Unicode and original line endings preserved, so that page reconstruction remains exact.
14. As an agent, I want invalid UTF-8 and controls represented reversibly, so that arbitrary artifact bytes are not silently replaced.
15. As an agent, I want oversized physical lines paginated, so that one line cannot be truncated or make progress impossible.
16. As an agent, I want present empty files returned explicitly, so that empty content cannot be confused with a missing file.
17. As an agent, I want invalid cursors and changed sources refused, so that a continuation cannot silently restart or skip bytes.
### Read trees and artifacts

18. As an agent, I want a named worktree file reader through --in, so that the address stays within the native tree authority.
19. As an agent, I want a worktree section or exact line range, so that a delegated read needs no pool path in a shell variable.
20. As an agent, I want related worktrees selected by spec binding, so that labels cannot invent a delivery relationship.
21. As an agent, I want complete filtered navigation across bounded pages, so that large related sets do not hide later assignments.
22. As an agent, I want spills and exported evidence through the same reader, so that artifact reads preserve exact content.
23. As an agent, I want unsafe file states refused before opening, so that a reader cannot block or follow an unintended producer.
24. As an agent, I want deterministic ordering, so that repeating a query at one source yields the same selection.
### Adopt and preserve

25. As an agent, I want eligible raw reads allowed with a matching native form, so that guidance does not become another shell ban.
26. As an agent, I want census evidence for eligible raw and native reads, so that later sessions can assess adoption without a claimed saving now.
27. As an agent, I want ordinary shell processing left alone, so that a conservative advisory does not claim to translate general sed programs.
28. As an agent, I want legacy outline, worktree, spec lifecycle, and evidence authority retained, so that the new readers cannot alter mutation or qualification behavior.
29. As an agent, I want independent CLI reconstruction expectations, so that a shared producer mistake cannot erase bytes unnoticed.
30. As an agent, I want source and fixture closure kept coherent, so that new command forms cannot bypass their existing contracts.

## Implementation decisions

### Source observations and prerequisites

The production source is `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`.
The planning base is `49a28e8124b5d919b4b39d4c346af962e1c1fa26`.
Full FT125, map answer 4, structured sources, definitions, and actual consumers informed these contracts.
These observations are static and claim no read savings or runtime proof.

`spec.Resolve` owns path-first, then live-slug addressing.
`spec.metadata` owns visible status selection.
Coverage's private `parse` owns authored stories, schema versions, table cells, row identity, and membership validation.
Its current `Rows` projection omits row IDs and rationales, so the reader needs a typed projection from that same parse.
`outline` currently locates declaration lines; its regular expressions cannot establish complete declaration bodies.
Go parser and token positions own those new spans.

`responsebound.New`, its Owner, and the existing response limits own final response budgeting and spill behavior.
`bounds.RefuseLinks` already checks complete path components.
`ClassifyNoFollow` checks the leaf and reads within the canonical control-record limit.
The exact reader composes those policies without changing other callers' empty-file or link postures.

ER-C2 consumes the accepted `specs/markdown-block-reader/spec.md` delivery, particularly MB-C1 and MB-C2.
That contract supplies classified physical lines, original byte offsets, visible H2 ownership, comments, fences, and at most one fault.
Coverage consumes its own approved adapter delivery before exposing spans from the canonical parse.
This is an actual prerequisite for authored-section chunks, not a literal 293-to-375-to-318-to-125 blocker.
No chunk depends on the unlanded native-record or preflight pair.

### CLI grammar and composition

The complete forms are:

```text
bench read <path> [--section <literal-H2-title> | --lines <first>:<last>] [--artifact] [--offset <byte-offset> --source <identity>]
bench spec show <slug-or-path> [--section stories|coverage|status|<literal-H2-title> | --story <positive>] [--offset <byte-offset> --source <identity>]
bench outline [path] --symbol <name-or-receiver.method> [--production|--test] [--offset <byte-offset> --source <identity>]
bench worktree read [--in <label|primary>] <path> [--section <literal-H2-title> | --lines <first>:<last>] [--offset <byte-offset> --source <identity>]
bench worktree list --spec <slug-or-path> [--offset <byte-offset> --source <identity>]
```

The shared `--in` dispatcher retains leading-flag placement for existing root calls.
For the new family leaf, it accepts the flag immediately after `worktree read`.
Resolve the target before leaf execution.
Worktree read is a native leaf over the same page owner as read.
It never returns a pool path that the caller must interpolate into another command.

Current tree_scope routing only projects a root verb, and treetarget.Run emits that one argv token.
Extend that existing call with an optional leaf name and preserve existing root calls.
Child argv must contain separate `worktree` and `read` tokens, never one token containing a space.
The scope declaration and usage projection derive the existing tree-target flag once.

`--artifact` permits an explicit absolute regular-file path for existing spills and exported evidence.
Ordinary reads remain contained in the selected repository tree.
No reader executes file contents.

The spec owner declares the show grammar and its file, status, and H2 selections.
The existing root command adapter composes that declaration with coverage's story and coverage-span projections.
Coverage already imports spec, so spec must not import coverage.
Keep this composition in `cmd/bench/planning_reads.go`; do not add a parallel story parser or pass-through policy module.
The grammar is declared once, and help and dispatch consume that declaration.

A first page defaults to byte offset zero.
A supplied offset requires its matching source token, even when that offset is zero.
Duplicate selectors, negative or nondecimal offsets, empty required operands, and incompatible flags are usage refusals with exit 2.
Operational input, ambiguity, syntax, and source-change refusals return exit 1 without partial content.
A present empty selection and a definitive absent symbol or story return exit 0 with different presence metadata.

`outline --symbol` is incompatible with `--full`.
An explicit Go file selects that file; a directory or omitted path searches its deterministic tracked Go scope.
Existing production/test filters retain their meaning.
Each ambiguous candidate names its file and qualified declaration, with a directly executable narrowing route.
No ambiguous candidate body is guessed.

`worktree list --spec` is incompatible with the selected-view and repeated-target forms.
Keep the existing default and selected views' bytes and authority.
Related navigation names read, path, exec, or clean actions according to existing assignment state.
A missing or inactive read target still refuses through the current tree resolver.

### Exact source selections

A raw file selection is all original bytes, with no manufactured final newline.
A physical-line range uses inclusive one-based bounds and includes each selected line's original terminator.
Ranges require positive ordered bounds within the existing physical lines.
An empty file has zero lines; an invalid range refuses instead of silently clipping it.
The default empty-file read still succeeds explicitly.

A literal section means the approved Markdown owner's visible H2 title, never arbitrary CommonMark heading levels.
Include the owning heading line and end at the next visible H2 or end of file.
Comments and fences cannot introduce an owning heading.
Repeated matching headings refuse; no matching heading returns explicit absence.
A Markdown fault refuses a parsed section before any partial selection is exposed.
A raw file read does not parse or grade Markdown.

The stories shortcut selects the User stories H2 span, including groups and author-line text.
The coverage shortcut uses the coverage owner's exact accepted-map span, including its heading and table.
The status shortcut returns the raw physical line chosen by the existing metadata owner.
Preserve its last-visible-line semantics and distinguish no status from an empty status value.
These shortcuts do not reproduce metadata or table recognition in the root adapter.

The prospective coverage API is `Describe(content []byte) (Document, error)` over its existing parse.
Its document exposes a numbered story selection, typed rows, and original source spans.
Coverage exports an ordered `Story` projection and typed `CoverageRow` values from its one private parse.
A row includes its identifier, story IDs, behavior, seam, rationale, and source span.
A story includes its number, exact authored text, and source span.

Keep original byte offsets separately from normalized cells used by validation.
The story result carries its exact story text plus every matching typed row in authored order.
Its schema names seam cells explicitly, rather than claiming that outline locations are blessed seams.

A malformed map refuses through the authoritative validator.
Retain story occurrences in its one parse so repeated selected numbers refuse as ambiguous.
This new selection refusal does not tighten existing coverage Check behavior, which currently deduplicates those numbers.
A valid absent story has `present=false` and empty row arrays.
The reader does not silently relax coverage rules or invent rows for a legacy spec that supplies none.
Legacy valid four-column schemas retain the optional fields their existing schema provides.

A Go selection uses `parser.ParseFile` with comments and token byte positions.
For a function or method, include its attached Doc group through the declaration End position.
For a general declaration, return the containing GenDecl, including its attached documentation.
A grouped declaration is one Go declaration even when it contains several specs.
Receiver-qualified names distinguish methods, including pointer and generic receiver spellings.
Exclude a preceding unattached comment and all adjacent declarations.

Slice original source bytes; do not format an AST back into new text.
A syntax error refuses a selected file before paging.
Candidate order is repository-relative file order, then source offset and qualified name.
Existing outline counts and bare symbol locations remain separate projections from the same owner.

### Page protocol and reversible representation

The prospective `PageRead(ReadRequest) (ReadPage, error)` operation lives in `internal/responsebound/read.go`.
It accepts a typed selection, raw selected bytes or the artifact stream, source identity, and command continuation arguments.
It returns one bounded rendered page and its exact successor, or a refusal.
Every reader uses this operation; no caller computes its own page size.
No new generic pagination framework or spill store is introduced.

A page carries source identity, selection kind, presence, encoding, byte offset, total selected bytes, next offset, and completion state.
Content is a JSON string cell for valid UTF-8, with controls escaped reversibly.
Invalid UTF-8 uses a base64 cell with the encoding named explicitly.
Cursors always count raw selected bytes, never rendered characters or base64 characters.
For valid UTF-8, split only at rune boundaries.
Decoding and concatenating page content must recover the exact selected bytes.

A story's typed projection has one versioned coverage-owned serialization.
Paging reconstructs that complete typed payload, including exact story text and source spans.
It does not pretend that a fragment is itself a complete row table.
Raw section and declaration pages reconstruct their original source slices instead.
Presence and payload kind remain explicit on every fragment.

Budget the actual rendered TOON after both JSON and TOON escaping.
Include page metadata, the complete help continuation, and the dispatcher tree frame in that budget.
Use the existing line and byte limits, currently 10 lines and 4096 bytes.
Nonterminal pages contain one fully executable continuation with the same selector, source token, and next offset.
Terminal pages explicitly have no successor.

Help exposes all cursor flags before a content result could spill.
Oversized metadata or a command whose required continuation cannot fit produces a bounded refusal before returning content.
Never emit a recursively spilled page or an empty nonterminal page.
A physical line larger than the response budget advances by source bytes across several pages.
A control-heavy line must still advance after rendered-byte accounting.

The source identity hashes the complete source bytes and a length-framed selection identity.
It also binds the repository or worktree identity, resolved path, and projection schema.
This token is a reader identity, not a replacement for completion-plan or immutable-evidence source digests.
For each continuation, re-read and compare the complete source identity before emitting content.
The original source and selected payload lengths are separate fields when their lengths differ.
If it changed, refuse and offer the exact first-page command; never silently reset the supplied cursor.

The file owner refuses links in every component, classifies regular files before opening, and detects changed file identity during a read.
A narrow bounds-owned regular-file opener may expose the existing classification to the streaming artifact path.
It must preserve the current classifiers' signatures and callers' policies.
Use existing filesystem fixtures rather than a global opener hook.
Refuse a replaced or modified file before returning a partial page.

Authored Markdown retains `ControlRecordLimit`; Go syntax selection retains the outline owner's `OutlineFileLimit`.
Reuse its scannable-file policy without acquiring caller-specific size guesses.
Artifact mode can stream larger regular files with bounded memory while computing a whole-file identity.
It may therefore read the whole artifact for each page; no read-cost saving is promised.
It creates no snapshot file, ledger entry, or new retention lifecycle.
Evidence mode means an exported file, not a guessed private evidence-pack path.

Reading an exported file does not verify its immutable manifest or qualify the evidence.
The existing `preflight evidence --verify`, `--check-current`, and `--to` authority remains canonical.
Response spill storage, privacy, pruning, and recorded output census keep their current owners.

### Related trees and advisory raw reads

Assignment records have no Spec field.
`commitment/repository.runScope` derives a binding's deliverable or a continuation's declared scope from the commitment state.
Expose `RunScopes(state, assignmentID)` from that owner, preserving its existing inventory rendering.
Related filtering consumes that projection and the existing `tickets.Covers` scope predicate.
It must not infer relationships from labels, branch names, or a spec file's physical presence.

Read the intent ledger once through `intent.Read`, obtaining Assignments and Commitment from the same returned ledger.
Sort matched assignments by their stable ID and include total matched and unbound counts.
An unbound assignment is excluded from proven related results and disclosed as unbound.
A missing binding is not authoritative proof that no external relationship exists.
CommitmentState alone excludes assignment population and fields, so its bytes cannot identify the complete related result.

Bind related-page identity to the complete deterministic related projection, resolved selector, repository identity, and projection schema.
That projection includes all row cells, matched and unbound counts, and the complete ordered navigation actions.
Build it through existing `listAssignmentRow`, list fact owners, `actionsForRows`, and AXI rendering before page splitting.
Exclude only reader source-token and cursor annotations from this content hash, avoiding self-reference.
The same owned projection supplies both identity bytes and content; do not rebuild another payload after hashing.

Recompute the complete projection for every continuation before returning any content.
Any changed row, count, or action refuses through the shared page owner, even when commitment data is unchanged.
This covers assignment additions or removals, labels, request values, state, paths, and branch facts when they change the projection.
It also covers tree presence, lease classification, landedness, ignored-file classification, recovery operands, and shared navigation availability.
Keep the existing owners' unknown and refusal postures; do not guess missing facts.

A changed input whose complete projection remains identical supplies identical query content and may continue.
The token identifies this bounded query decision, not execution authenticity or future authority to act on a worktree.
Every offered action still resolves and authorizes its target through the existing execution owner.
This choice grants no semantic exemption for authored spec or completion-plan digests.
No new database, snapshot file, or retention lifecycle is introduced.

The follow-on Bash shim already forwards stderr when the current core permits a command.
Add advisory output after the existing pool-path and follow-on verdicts.
Preserve every existing deny surface and exit code.
Ordinary supported reads stay allowed with exit 0; pool-path authority remains the existing guard's decision.
The advisory identifies a native form and never executes or rewrites the shell command.

Eligibility is conservative and structural through the existing shellcommand owner.
Recognize direct named-file cat, head, tail, and a simple read-only sed line-print form when its path and range are unambiguous.
Name file or line-range reading as appropriate; do not translate arbitrary sed programs.
Exclude stdin-only calls, source examples, unknown parse states, pipelines requiring transformations, and controlled or unresolved operands.
Those are negative controls, not blanket hook exceptions or new registration requirements.

Keep lexical observations and head resolution in census, composing the existing shellcommand parse and benchguard routine-prefix resolver.
The command adapter supplies a verified current assignment through `intent.AssignmentForWorktree` when needed.
Eligible owned reads use a narrow census-owned observation operation with the verified assignment identifier.
Combine it with existing pool-path attribution so one invocation appends at most one raw-call record.
Reuse the existing census codec, append owner, and identifier validation.

Keep codec functions in census.go because its existing test reads that exact source path.
Do not add another command lexer, scanner, or raw-call record format.

Native reader calls retain the existing output census and verb head identity.
The existing outputHead adapter derives complete native read-form heads from canonical grammar declarations.
It distinguishes read, spec show, outline --symbol, and worktree read while preserving other heads.
These records count attempts, not successful adoption or measured savings.

Primary or unowned reads receive no invented assignment attribution.
Later session evidence compares raw eligible calls and native calls at the same assignment scope.
This spec claims no adoption percentage, latency improvement, or reduced read count.

### Placement, growth, and transitive closure

The selected planning lane uses per-file Growth through `internal/gate/lane.go`.
The default increased-file cap is 400 lines unless an exact grant applies.

Growth refuses only when a changed source exceeds both its applicable cap and its base line count.
`acceptedReason` uses exact keys; the conformance-directory accept does not waive individual file growth.
Directory crowding is inherited debt, not a new Growth refusal.
No count-equal grant or structure-budget edit is proposed.
Each placement below is part of the first chunk that needs that owner.

main.go has 436 lines and must retain commandRegistry in place.
Move only `commandsGrammar`, `commandsCommand`, and `commandsProbeIsStale` into commands.go, removing about 60 lines.
Keep the new registry and help edits within a 400-line main.go result.
Every conformance source reader that parses commandRegistry still reads main.go.
`command_registry.go` has 387 lines; put new adapter logic in planning_reads.go instead.

`command_registry_test.go` has 794 lines and no new default growth allowance.
Move its existing envelope types and runner plus fixture helpers into axi_query_envelope_test.go.
Move only `axiEnvelopeCases` into axi_query_cases_test.go, preserving one independent expectation table.
Move the existing final help tests into command_help_test.go, leaving the original source at about 371 lines.
The new files target at most 350 lines and reuse all existing helpers.
Existing declarations and callers keep their names except where a test name inaccurately says only retire and history.

Update directArchitectureTests to cover the relocated envelope harness, with no broader exemption.
The ticket owner's existing five command-registry names extend to the moved envelope, cases, and help files as required.
Its registered-command closure and registry tests remain the sole ownership check.
New production-entry tests use the same fixtures, never a copied repository builder or child benchmark harness.

`axi_query_registry_test.go` has 445 lines.
Move its existing guidance and bidirectional membership bite tests into axi_query_bite_test.go before adding approved memberships.
Leave its parser and independent membership map in the original owner, with a target below 400 lines.
Keep the parser's main.go path and exact diagnostics.

Approved additions are read, outline, spec show, and worktree read; worktree list is already approved.
Existing outline metadata also gains derived AXI continuation help while preserving its result table.
Spec retire and history remain outside the newly approved show query.

Coverage.go has 599 lines.
Move Command into command.go, validation into validation.go, and existing State/Rows projections into projection.go.
Keep the one parser and add original-byte span observations there, targeting at most 380 lines.
New typed story serialization and witnesses live in story.go and story_test.go.
Existing validators and citation helpers keep their public behavior and sufficient tests.

Outline.go has 383 lines and 17 lines of headroom.
Move its existing Command into command.go before adding symbol forms, leaving the location owner near 240 lines.
New AST span selection and its witnesses live in declaration.go and declaration_test.go.
Spec.go has 502 lines; move its existing Command/specArg and retirement effect family into command.go and retire.go.
Leave status, facts, lifecycle derivation, and addressing in their original owner, with a target below 350 lines.

Responsebound owner.go has 234 lines, leaving 166 lines of default headroom.
New page and file operations use siblings targeted below 350 lines each.
Repository.go has exactly 400 lines; move its existing runScope and typed projection together into run_scope.go.
List.go has 332 lines; keep related filtering in related.go and only route the new form through the original list adapter.

Census.go has 398 lines.
Move its existing raw-call orchestration and head-resolution family into raw_calls.go before adding advice.
Keep `composeRecord`, parsing, and the separator literal in census.go for the existing path-sensitive codec witness.
New eligibility and attribution witnesses use read_advice_test.go.
Cmd guards.go has 239 lines; the adapter fits without changing the shim or introducing another hook.

Registry_retained_workflow.go has 396 lines.
Preserve its existing target-slot needle and use a new registry_planning_readers.go for reader-specific anchors.
Registry_data.go has 479 lines, no file accept grant, and a default 400-line cap.
Extend its existing append expression without growing that inherited overage.

Docs_workflow_checks_test.go has 786 lines and an exact 786-line grant.
Fixture_bite_test.go has 872 lines against its 826-line grant and must not grow.
Their sufficient checks remain unchanged.
No redundant conformance test or additional system harness is proposed.

The existing tree-target system witness supplies the actual child boundary with its current selected binary and three-repository owner.
Tree_target_test.go has 115 lines and 285 lines of default headroom.
Add reader rows and paginated reconstruction to its existing TestTreeTargetRunsInNamedWorktree.
Reuse treeTargetFixture, runTreeTarget, and current fixture source; do not edit or copy the 445-line system owner.

The complete closure includes the root registry, five bound-command registry names, AXI membership, profile, grammar guidance, and relevant anchors.
The fixture census follows recursive BASE includes, dot-restored overlay paths, and MUTATE.json inputs.
Every affected existing fixture directory is named below.
Relocated source readers and architecture ownership stay visible in the same first green chunk.
No existing guard, phase registry, or evidence qualification policy is reassigned.

## Implementation chunks

The independent SPEC review accepted source ab74806ee43dfcbc339050d88b8c358ad7a12db3 before ticket slicing.
These five green outcomes have one ticket each; independent graph review accepted source 5b87e01c9a109f07ef7c33d019cf18d6496aea11.
Shared page and command-registry writes order each later chunk after its preceding reader delivery.
Same-file overlaps with other accepted plans require the reviewed earlier owner version or explicit composition before starting.
They do not invent dependencies on unrelated unlanded capabilities.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| ER-C1 / 1.md | Native raw-file and artifact pages reconstruct exact bytes through the CLI and expose bounded continuation grammar. | ER1, ER16-ER28, ER35-ER40, ER44-ER46 | planned page/file and real CLI reconstruction families | gpt-5.6-sol / xhigh |
| ER-C2 / 2.md | Canonical Markdown selections and typed story obligations read through the spec CLI. | ER2-ER10 | planned Markdown/spec/story families and real CLI reconstruction | no |
| ER-C3 / 3.md | Go syntax selections return complete declarations with deterministic narrowing. | ER11-ER15 | planned declaration families and real CLI reconstruction | no |
| ER-C4 / 4.md | Owned worktree reads and binding-based related navigation use the same exact page contract. | ER29-ER34, ER48 | planned related filtering and worktree CLI families | no |
| ER-C5 / 5.md | Eligible raw reads remain allowed with native advice and attributable census evidence. | ER41-ER43, ER47 | planned advice/census families and preserved guard refusals | no |

ER-C1 includes source splits and complete query-registry closure for its first usable read form.
Each later chunk adds its own approved query and all profile, anchor, fixture, and registry closure on the same green outcome.

Every chunk editing BENCH-reference or the benchkit profile co-names internal/anchors/registry_data_test.go at that first use.
ER-C1 establishes this closure; ER-C2 through ER-C5 retain it whenever their advertised forms edit either subject.
The holder's independent assertions remain intact, and its existing role adds no fixture unit.
ER-C2 waits for the named canonical Markdown delivery; other chunks do not acquire an artificial Markdown prerequisite.
ER-C4 includes its physical-line range behavior, which ER-C1 supplies as the reusable file owner.
ER-C5 changes no existing denial policy and has no dependency on a future guard scanner.

### Completion plan

This version-1 plan binds the five ticket checkpoints and their required verification inventories.
A delegated build must author its complete version-2 execution assignments through the existing validator before dispatch.
This graph does not invent execution identity or depend on the separate native-record capability.
Each chunk reviews its frozen predecessor and current tips before the next ticket starts.

Probe identities below name required future mutations and their named failing witnesses.
Run them through the existing probe and verification owners, retain actual output, restore the source, and demonstrate green.
Query-membership and help-grammar omission witnesses also remain required for each delivered form.
No runtime result is asserted by this plan.

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "ER-C1",
      "tickets": [
        "1.md"
      ],
      "verification": [
        {
          "id": "ER-C1-page",
          "command": "bench test --package ./internal/responsebound"
        },
        {
          "id": "ER-C1-bounds",
          "command": "bench test --package ./internal/bounds"
        },
        {
          "id": "ER-C1-entry",
          "command": "bench test --package ./cmd/bench --run \"TestExactReaderCLIReconstruction|TestPlanningReaderHelpAndRefusals\""
        },
        {
          "id": "ER-C1-successor-red",
          "command": "bench test --package ./cmd/bench --run TestExactReaderCLIReconstruction",
          "probe": "skip one source byte in continuation successor"
        },
        {
          "id": "ER-C1-query",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "ER-C1-routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "ER-C1-help",
          "command": "bench test --package ./cmd/bench --run Help"
        }
      ]
    },
    {
      "id": "ER-C2",
      "tickets": [
        "2.md"
      ],
      "verification": [
        {
          "id": "ER-C2-spec",
          "command": "bench test --package ./internal/spec"
        },
        {
          "id": "ER-C2-coverage",
          "command": "bench test --package ./internal/coverage"
        },
        {
          "id": "ER-C2-entry",
          "command": "bench test --package ./cmd/bench --run \"TestExactReaderCLIReconstruction|TestPlanningReaderHelpAndRefusals\""
        },
        {
          "id": "ER-C2-query",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "ER-C2-routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "ER-C2-help",
          "command": "bench test --package ./cmd/bench --run Help"
        }
      ]
    },
    {
      "id": "ER-C3",
      "tickets": [
        "3.md"
      ],
      "verification": [
        {
          "id": "ER-C3-outline",
          "command": "bench test --package ./internal/outline"
        },
        {
          "id": "ER-C3-entry",
          "command": "bench test --package ./cmd/bench --run \"TestExactReaderCLIReconstruction|TestPlanningReaderHelpAndRefusals\""
        },
        {
          "id": "ER-C3-query",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "ER-C3-routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "ER-C3-help",
          "command": "bench test --package ./cmd/bench --run Help"
        }
      ]
    },
    {
      "id": "ER-C4",
      "tickets": [
        "4.md"
      ],
      "verification": [
        {
          "id": "ER-C4-worktree",
          "command": "bench test --package ./internal/worktree"
        },
        {
          "id": "ER-C4-scopes",
          "command": "bench test --package ./internal/commitment/repository"
        },
        {
          "id": "ER-C4-targets",
          "command": "bench test --package ./internal/treetarget"
        },
        {
          "id": "ER-C4-lines",
          "command": "bench test --package ./internal/responsebound --run TestExactLineRanges"
        },
        {
          "id": "ER-C4-entry",
          "command": "bench test --package ./cmd/bench --run \"TestRelatedContinuationRefusesChangedProjection|TestPlanningReaderHelpAndRefusals|TestTreeScope|TestWorktree\""
        },
        {
          "id": "ER-C4-projection-red",
          "command": "bench test --package ./cmd/bench --run TestRelatedContinuationRefusesChangedProjection",
          "probe": "replace complete related projection identity with commitment-only identity"
        },
        {
          "id": "ER-C4-system",
          "command": "bench test --check system"
        },
        {
          "id": "ER-C4-query",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "ER-C4-routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "ER-C4-help",
          "command": "bench test --package ./cmd/bench --run Help"
        }
      ]
    },
    {
      "id": "ER-C5",
      "tickets": [
        "5.md"
      ],
      "verification": [
        {
          "id": "ER-C5-advice",
          "command": "bench test --package ./cmd/bench --run \"TestRawReadAdvice|TestExactReaderOutputHeads|TestGuard\""
        },
        {
          "id": "ER-C5-census",
          "command": "bench test --package ./internal/census"
        },
        {
          "id": "ER-C5-head-red",
          "command": "bench test --package ./cmd/bench --run TestExactReaderOutputHeads",
          "probe": "omit the native read-form qualifier at outputHead"
        },
        {
          "id": "ER-C5-query",
          "command": "bench test --check axi-query-registry"
        },
        {
          "id": "ER-C5-routing",
          "command": "bench test --check subcommand-routing"
        },
        {
          "id": "ER-C5-help",
          "command": "bench test --package ./cmd/bench --run Help"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "ER-final-ordinary",
      "command": "bench test --package ./..."
    },
    {
      "id": "ER-final-system",
      "command": "bench test --check system"
    },
    {
      "id": "ER-final-coverage",
      "command": "bench coverage specs/exact-planning-readers/spec.md --check"
    }
  ]
}
```

## Testing decisions

Choose the smallest sufficient seam for each wrong result.
Page tests grade rendered escaped bytes, successor progress, and reversible decoding at responsebound.
Markdown and story cases exercise their actual existing parser owners.
Declaration cases use actual Go syntax and original file bytes.
File-state fixtures expose no-follow and type refusal without a new global port.

The real CLI test reconstructs complete file, section, story payload, declaration, spill, and exported-evidence selections across all pages.
The existing tree-target system witness reconstructs owned worktree pages through the selected executable and its real --in child.
The command tests separately grade changed-source refusal before receiving more bytes.
TestRelatedContinuationRefusesChangedProjection uses the real list CLI and a nonterminal related page.
It keeps the commitment portion unchanged while mutating assignment population or fields and each material observed-fact family.

Cases cover a matching or unbound addition, removal, label or request change, assignment state, and rendered path or branch effects.
Other cases change tree presence, lease classification, landedness, ignored-file classification, and resulting recovery or shared navigation actions.
Use the existing linked assignment and filesystem fixtures, and ensure each case changes the complete projected payload.
The continuation must return refusal before any content cell, with the original commitment data still equal.
Do not claim coverage from a file-only cursor mutation.

Reuse existing command envelope, linked worktree, and filesystem helpers.
The unit owners' tables cover many edge predicates; do not create one test per row.
No test exists solely to prove a declaration was moved or added.

The independent CLI expectation retains independently authored original fixture bytes.
During implementation, omit one byte when computing the successor and record the named reconstruction test turning red.
Restore the exact source and demonstrate green.
Also retain existing bidirectional query-membership and help-grammar omission witnesses when registering each new form.
Also omit the read-form qualifier at outputHead and record TestExactReaderOutputHeads turning red.

For ER48, replace complete-projection identity with commitment-only identity during implementation.
Record TestRelatedContinuationRefusesChangedProjection turning red on its unchanged-commitment cases, then restore and demonstrate green.
These future reds are required evidence, not observations made during specification.

### Seam diagram

    native selection + existing tree authority
        │
        ├── spec / canonical Markdown spans
        ├── coverage typed story obligations
        └── outline / Go syntax declaration spans
                         │
                         ▼
               responsebound exact page owner
                         │
                         ▼
          source-bound content + executable successor

    existing shell observations ──▶ allowed advice + existing census codec

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| ER1 | 1 | A file read reconstructs the exact source bytes | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | A summarizer or newline-normalizing adapter produces different bytes. |
| ER2 | 2 | A literal visible H2 selection includes its heading and ends before the next visible H2 | planned TestExactMarkdownSelections in internal/spec/read_test.go | Whole-file or regex slicing leaks the neighboring section. |
| ER3 | 2 | A heading hidden in a fence or HTML comment is not selected | planned TestExactMarkdownSelections in internal/spec/read_test.go | An independent heading detector selects authored examples as owning headings. |
| ER4 | 2 | Repeated visible matching H2 headings refuse as ambiguous | planned TestExactMarkdownSelections in internal/spec/read_test.go | Choosing the first occurrence conceals the ambiguity. |
| ER5 | 3 | The stories shortcut returns the actual User stories H2 span | planned TestExactSpecSelections in internal/spec/read_test.go | An invented summary drops the author line or grouped story text. |
| ER6 | 3 | The coverage shortcut returns the authoritative coverage-map span | planned TestExactStoryProjection in internal/coverage/story_test.go | A second table detector disagrees with the validator boundary. |
| ER7 | 3 | The status shortcut returns the metadata owner's selected raw Status line | planned TestExactSpecSelections in internal/spec/read_test.go | Choosing a fenced or earlier status changes existing last-visible-line semantics. |
| ER8 | 4 | A story projection contains its exact authored text and all typed matching rows with named seams | planned TestExactStoryProjection in internal/coverage/story_test.go | A display-only Rows projection loses row identity, membership, or rationale. |
| ER9 | 4 | Malformed or ambiguous story membership refuses rather than returning a partial projection | planned TestExactStoryProjection in internal/coverage/story_test.go | Ignoring validator faults makes an invalid map look like a complete obligation. |
| ER10 | 5 | An absent story returns present false and no rows | planned TestExactStoryProjection in internal/coverage/story_test.go | An empty string alone cannot identify a missing story. |
| ER11 | 6 | A declaration span begins at its attached doc comment and ends at the AST declaration end | planned TestDeclarationSpans in internal/outline/declaration_test.go | Regex locations cannot prove a complete body or attached documentation. |
| ER12 | 7 | A selected declaration excludes neighboring comments and declarations | planned TestDeclarationSpans in internal/outline/declaration_test.go | Taking a line window includes an unrelated producer. |
| ER13 | 8 | Ambiguous same-name declarations return deterministic candidates without a guessed body | planned TestDeclarationSelection in internal/outline/declaration_test.go | A first-match implementation chooses the wrong file or receiver. |
| ER14 | 8 | A receiver-qualified selection distinguishes methods from same-name functions | planned TestDeclarationSelection in internal/outline/declaration_test.go | Name-only matching confuses independent production seams. |
| ER15 | 9 | A syntax-invalid selected Go file refuses before any declaration page | planned TestDeclarationSelection in internal/outline/declaration_test.go | A textual locator appears successful even though the Go parser cannot establish a span. |
| ER16 | 10 | Every page includes the source identity and selected-payload byte offset | planned TestExactPageProtocol in internal/responsebound/read_test.go | A continuation cannot prove which bytes precede it without both fields. |
| ER17 | 10 | Changing the source between pages refuses without emitting new content | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | A cursor alone combines the old prefix with a new suffix. |
| ER18 | 11 | A nonterminal page includes its complete executable continuation command | planned TestExactPageProtocol in internal/responsebound/read_test.go | A descriptive placeholder leaves the next page inaccessible. |
| ER19 | 11 | The pagination grammar is available before any read can spill | planned TestPlanningReaderHelpAndRefusals in cmd/bench/planning_reader_commands_test.go | Adding grammar only to a spilled body hides the recovery route. |
| ER20 | 12 | Rendered escaped output stays within the shared byte and line budget | planned TestExactPageProtocol in internal/responsebound/read_test.go | Raw source length understates TOON and JSON escaping expansion. |
| ER21 | 13 | Unicode page boundaries reconstruct the original bytes | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | Splitting a UTF-8 rune changes the round trip. |
| ER22 | 13 | CRLF, LF, and final-newline absence reconstruct without normalization | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | Line-based reassembly manufactures or drops terminators. |
| ER23 | 14 | Control bytes round-trip through the declared reversible content encoding | planned TestExactPageProtocol in internal/responsebound/read_test.go | Literal controls break the response while replacement characters lose the source. |
| ER24 | 14 | Invalid UTF-8 round-trips through base64 with raw-byte cursors | planned TestExactPageProtocol in internal/responsebound/read_test.go | Encoding invalid bytes as a Go string replaces their values. |
| ER25 | 15 | A single oversized physical line advances through multiple exact pages | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | A line-count paginator cannot advance within the line. |
| ER26 | 16 | A present empty file returns present true with empty content and terminal position zero | planned TestExactFileStates in internal/responsebound/file_test.go | Treating StateEmpty as absence erases an existing file. |
| ER27 | 17 | An offset without its source token refuses as usage | planned TestPlanningReaderHelpAndRefusals in cmd/bench/planning_reader_commands_test.go | An unbound cursor can skip bytes from an unrelated source. |
| ER28 | 17 | An invalid or out-of-range cursor refuses without resetting to page zero | planned TestExactPageProtocol in internal/responsebound/read_test.go | A forgiving parser repeats or omits already reconstructed bytes. |
| ER29 | 18 | The real --in worktree read resolves only an active owned target | planned extension of TestTreeTargetRunsInNamedWorktree in internal/systemtest/tree_target_test.go | Direct filesystem lookup bypasses existing assignment authority. |
| ER30 | 19 | A worktree section read reconstructs the same selected source span | planned extension of TestTreeTargetRunsInNamedWorktree in internal/systemtest/tree_target_test.go | A separate worktree selector drifts from the primary reader. |
| ER31 | 19 | An inclusive physical-line range preserves its original terminators | planned TestExactLineRanges in internal/responsebound/file_test.go | Splitting and joining lines silently normalizes the selected bytes. |
| ER32 | 20 | A spec filter uses exact commitment bindings or declared continuation scope | planned TestRelatedWorktreeFilter in internal/worktree/related_test.go | Matching a label or physical spec file invents a relationship. |
| ER33 | 20 | An unbound assignment is excluded with an explicit unbound count | planned TestRelatedWorktreeFilter in internal/worktree/related_test.go | Guessing missing Assignment.Spec data yields a false complete relation. |
| ER34 | 21 | Related results paginate in deterministic assignment-ID order with complete counts | planned extension of TestTreeTargetRunsInNamedWorktree in internal/systemtest/tree_target_test.go | A bounded first page alone hides later related worktrees. |
| ER35 | 22 | An existing spill file reconstructs through the same exact page protocol | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | A separate artifact renderer can exceed the shared escaped-byte bound. |
| ER36 | 22 | An exported evidence file reconstructs without granting evidence authority | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | Reading a file must not assert manifest verification or qualification. |
| ER37 | 23 | A symlink in any path component refuses before opening the target | planned TestExactFileStates in internal/responsebound/file_test.go | Checking only the final entry allows a linked parent to redirect the reader. |
| ER38 | 23 | A special or unreadable file refuses without partial content | planned TestExactFileStates in internal/responsebound/file_test.go | A reader that opens before classification can block or expose an incomplete prefix. |
| ER39 | 23 | A missing file is an operational refusal distinct from an empty file | planned TestExactFileStates in internal/responsebound/file_test.go | Collapsing both states conceals an unavailable input. |
| ER40 | 24 | Repeated selections at an unchanged source return the same bytes and cursor sequence | planned TestExactPageProtocol in internal/responsebound/read_test.go | Nondeterministic row or candidate ordering changes continuation boundaries. |
| ER41 | 25 | An eligible ordinary raw read remains allowed while stderr names the matching native form | planned TestRawReadAdvice in cmd/bench/planning_read_advice_test.go | A new deny code would turn guidance into an unapproved ban. |
| ER42 | 26 | Owned eligible raw reads and native reader calls remain attributable in the existing census | planned TestReaderCensusAttribution in internal/census/read_advice_test.go | An unrecorded advisory cannot support a later adoption comparison. |
| ER43 | 27 | A source example, stdin read, or general sed program receives no guessed translation | planned TestRawReadAdvice in cmd/bench/planning_read_advice_test.go | Name-only matching turns ordinary shell text into a false read recommendation. |
| ER44 | 28 | Legacy outline metadata and lifecycle or evidence commands retain their existing refusal authority | review-owned scope comparison at internal/outline/outline.go, internal/spec/spec.go, internal/preflight/evidencecmd/evidence.go | Reader source identity does not replace lifecycle status or immutable evidence manifests. |
| ER45 | 29 | The real CLI reconstruction fails when continuation skips one source byte | planned TestExactReaderCLIReconstruction in cmd/bench/planning_reader_commands_test.go | A producer-derived expectation repeats the same cursor omission and stays green. |
| ER46 | 30 | New command forms satisfy the existing real-envelope membership and routing checks | planned TestPlanningReaderHelpAndRefusals in cmd/bench/planning_reader_commands_test.go | Owner-only tests do not detect missing root registration or tree routing. |

| ER47 | 26 | Native reader output records identify their complete supported read form | planned TestExactReaderOutputHeads in cmd/bench/census_output_test.go | A root-only spec or outline head conflates exact reads with other operations. |

| ER48 | 21 | Related continuation refuses before content when its complete projection changes while commitment data stays unchanged | planned TestRelatedContinuationRefusesChangedProjection in cmd/bench/planning_reader_commands_test.go | Commitment-only identity misses assignment population, fields, observed tree facts, counts, and navigation changes. |

### Edge inventory

- Present empty and absent sources differ; a terminal empty page has no invented successor.
- Unicode, invalid UTF-8, NUL, tabs, quotes, backslashes, CRLF, LF, and missing final newline preserve reversible bytes.
- A huge physical line and heavily escaped controls require byte progress within both rendered limits.
- Empty H2 bodies, repeated H2 titles, hidden headings, malformed comments or fences, and absent selectors have explicit results.
- Story maps cover optional legacy columns, escaped pipes, repeated story membership, missing rows, and invalid numbering through one validator.
- Go cases cover attached versus unattached comments, grouped declarations, pointer/generic receivers, malformed files, and ambiguous names.
- Decimal overflow, repeated cursor flags, unbound offsets, out-of-range positions, oversized metadata, and impossible continuations refuse.
- Symlink parents and leaves, directories, FIFOs, missing files, permission refusal, and changing source identity return no partial page.
- Related sets cover zero results, unbound rows, inactive ownership, continuation scopes, and multiple bounded pages.
- Unchanged-commitment continuation cases change assignment population, rendered assignment fields, tree observations, counts, and navigation actions before the next page.
- Advisory controls cover ordinary text, source examples, stdin, transformations, unsupported sed programs, and current pool-path denials.

**Won't handle** — arbitrary heading levels or a full CommonMark parser — the accepted Markdown owner defines visible H2 spans.
**Won't handle** — a declaration's runtime call closure or a blessed seam label — syntax gives exact declarations, not architecture judgment.
**Won't handle** — atomic protection against hostile external filesystem writers — changed-source refusal and the current ownership contract bound these reads.
**Won't handle** — live evidence qualification through artifact content — immutable manifest and source-pair authority remain separate.
**Won't handle** — translating arbitrary shell or sed transformations — eligible read advice remains conservative and nonexecuting.

## Ownership fences

These are prospective exact owner, test, registry, anchor, and fixture paths.
New files have no runtime existence claim.
A co-owned existing file stays untouched when its current sufficient check needs no change.
No implementation is authorized by this planning artifact.



- `.agents/skills/bench-craft-cli/SKILL.md`
- `.bench/BENCH-reference.md`
- `cmd/bench/axi_query_cases_test.go`
- `cmd/bench/axi_query_envelope_test.go`
- `cmd/bench/census_output.go`
- `cmd/bench/census_output_test.go`
- `cmd/bench/command_help_test.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/commands.go`
- `cmd/bench/guards.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/planning_read_advice.go`
- `cmd/bench/planning_read_advice_test.go`
- `cmd/bench/planning_reader_commands_test.go`
- `cmd/bench/planning_reads.go`
- `cmd/bench/tree_scope.go`
- `cmd/bench/tree_scope_test.go`
- `cmd/bench/worktree_leaves.go`
- `cmd/bench/worktree_leaves_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_planning_readers.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/bounds/classify.go`
- `internal/bounds/open.go`
- `internal/bounds/open_test.go`
- `internal/census/census.go`
- `internal/census/census_test.go`
- `internal/census/output.go`
- `internal/census/raw_calls.go`
- `internal/census/read_advice.go`
- `internal/census/read_advice_test.go`
- `internal/commitment/repository/repository.go`
- `internal/commitment/repository/run_scope.go`
- `internal/conformance/axi_query_bite_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/docs_workflow_checks_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/help_inventory_single_source_test.go`
- `internal/conformance/ordinary_build_census_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/coverage/command.go`
- `internal/coverage/coverage.go`
- `internal/coverage/coverage_test.go`
- `internal/coverage/projection.go`
- `internal/coverage/story.go`
- `internal/coverage/story_test.go`
- `internal/coverage/validation.go`
- `internal/outline/command.go`
- `internal/outline/declaration.go`
- `internal/outline/declaration_test.go`
- `internal/outline/outline.go`
- `internal/outline/outline_test.go`
- `internal/responsebound/file.go`
- `internal/responsebound/file_test.go`
- `internal/responsebound/owner.go`
- `internal/responsebound/owner_test.go`
- `internal/responsebound/read.go`
- `internal/responsebound/read_test.go`
- `internal/spec/command.go`
- `internal/spec/read.go`
- `internal/spec/read_test.go`
- `internal/spec/resolve.go`
- `internal/spec/retire.go`
- `internal/spec/spec.go`
- `internal/spec/spec_test.go`
- `internal/systemtest/tree_target_test.go`
- `internal/tickets/registry_data.go`
- `internal/tickets/registry_data_test.go`
- `internal/treetarget/family_test.go`
- `internal/treetarget/run.go`
- `internal/usage/worktree.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/list_selected.go`
- `internal/worktree/list_selected_test.go`
- `internal/worktree/related.go`
- `internal/worktree/related_test.go`
- `projects/benchkit.md`
- `reviews/exact-planning-readers.md`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/package-core-guard/bounds-classify-limit-restated`
- `tests/canary/package-core-guard/bounds-read-limit-restated`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`

## Out of scope

- General text processing or summarization: estimated six edits and four gate runs. The raw-file and syntax selectors remain exact readers.
- A semantic completion-plan digest: estimated five edits and three gate runs. Reader identity never exempts prose from existing evidence digests.
- Private evidence-store access or new qualification: estimated four edits and three gate runs. Existing evidence export and verification remain the authority.
- A new immutable read snapshot store: estimated five edits and four gate runs. Continuation rechecks the source and refuses change.
- General shell-program translation: estimated five edits and four gate runs. Conservative raw-read advice has no execution authority.
- Global hook replacement, source-policy scanners, arbitrary JSON editing, and assignment dispatch remain separate capabilities.

## Further notes

The shared reviewed planning map remains in place.
No roadmap priority, FT290 field, commitment, or other approved implementation line changes.
FT125's reader contract is covered; claimed read savings require later comparable session evidence.
The source census covers current definition callers, same-package consumers, command bindings, source-reading tests, and fixture input closure.
No external consumer absence is asserted from an in-tree search.

The required planning lane and document checks freeze this complete spec before independent review.
The original independent three-axis review and this bounded repair are retained in [the spec review record](assets/spec-review.md).
Runtime tests, omission probes, benchmarks, and source relocations remain future work.
The independent graph review accepted all five tickets; dispatch still requires the complete version-2 execution binding and fresh authors.
The separate preflight pair keeps closure and staleness policy authority.
Ticket slicing must recheck overlapping source obligations and preserve each first green registry closure.

| approval item | proposed disposition |
| --- | --- |
| implementation line | gpt-5.6-sol / high, ER-C1 xhigh, preserved from accepted SPEC |
| test purpose | Exact owner-level predicates plus one real CLI reconstruction family |
| source authority | Canonical Markdown, coverage, Go syntax, response limits, and worktree binding owners |
| graph | Five serial green tickets accepted at 5b87e01c9a109f07ef7c33d019cf18d6496aea11 |
| evidence limits | Static source and document checks only; no runtime or read-savings claim |

### Accepted ticket graph

| ticket | title | Blocked by | delivered outcome |
| --- | --- | --- | --- |
| 1.md | Deliver exact native file and artifact pages | none | Accepted ER-C1 vertical outcome with its own verification and review checkpoint |
| 2.md | Read canonical spec sections and story obligations | 1.md | Accepted ER-C2 vertical outcome with its own verification and review checkpoint |
| 3.md | Read complete Go declarations from syntax spans | 1.md, 2.md | Accepted ER-C3 vertical outcome with its own verification and review checkpoint |
| 4.md | Read owned worktrees and source-bound related navigation | 1.md, 2.md, 3.md | Accepted ER-C4 vertical outcome with its own verification and review checkpoint |
| 5.md | Advise eligible raw reads through the existing census | 1.md, 2.md, 3.md, 4.md | Accepted ER-C5 vertical outcome with its own verification and review checkpoint |

All preceding overlapping tickets are explicit blockers, including shared registry, profile, anchor, fixture, and command-adapter writes.
ER-C2 also requires the accepted canonical Markdown and coverage-adapter delivery named above.
Other capabilities and advisory roadmap ordering introduce no additional literal blockers.
Ticket Writes union equals the accepted implementation fence, excluding only the review pickup.
Co-owned sufficient checks and immutable fixtures need no edit when their accepted input remains valid.
