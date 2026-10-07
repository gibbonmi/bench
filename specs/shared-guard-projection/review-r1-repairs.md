# GP-C1 decoder preservation repair

Reviewed source: 9b61a523d9a2d6cc6d3ba3f1a72f8aee6a609731
Reviewed spec SHA256: ec7e0275f2ffbbbb5d6a4c28ec8d0d5b10e11e9fae05b86219416daa700a1472
Review result: Standards 0, Spec 0, Coverage 1 material finding
Finding: GP-C1, confidence 10
Repair author: gpt-6.1-sol/high, one bounded pass
Confirmation: independently accepted at 02c799aa8e7ba9f1dfcf38a14906f51ab4a6d8db
Confirming reviewer: adoption_spec_review, retained high review line
Confirming result: Standards 0, Spec 0, Coverage 0; confidence 10

The first common reader contract omitted observable differences between the three encoding/json target shapes.
Folded lookup or a final RawMessage map alone would change guard decisions.

Current evidence read:

- internal/gitguard/gitguard.go, CommandFromEnvelope: tagged nested command string.
- internal/benchguard/benchguard.go, CommandFromEnvelope: tagged root field and exact RawMessage input map.
- internal/writeguard/writeguard.go, PathFromEnvelope: tagged scalar cwd and exact RawMessage input map.
- cmd/bench/guards.go, guardGit and guardFileWrite: real entry exits and authority reads.
- Installed Go 1.25.14 encoding/json/decode.go: folded struct matching, ordered duplicate merging, earliest type errors, and scalar-null preservation.

The repair centralizes the original standard decoder target shapes in three structural forms.
Callers retain required fields, validation, diagnostics, and authority.
No custom duplicate parser, duplicate rejection, key normalization, or strict JSON mode is introduced.

GP47–GP61 and family EP extend the existing planned reader, core, and differential tests.
Actual guard-git witnesses cover folded keys, retained null values, and retained earlier errors.
Actual guard-file-write witnesses cover retained cwd and folded root keys for a tracked relative primary path.
The existing implementation fence contains every added witness path.
No production or test source changed during this repair.

All input cases derive from the current decoder definitions and the standard-library contract.
Runtime evidence remains planned.
The confirming review accepted this repair before the separate placement amendment.

The subsequent placement amendment preserves this accepted repair and its source evidence.
Ticket slicing and independent graph review are complete; their accepted source is recorded below.

## SPEC acceptance before graph authoring

The earlier repair confirmer was /root/adoption_spec_review, declared gpt-6.1-sol/high.
It accepted both decoder/source repairs at 02c799aa8e7ba9f1dfcf38a14906f51ab4a6d8db on Standards, Spec, and Coverage with zero findings and confidence 10.
The separate placement confirmer was /root/deepen_execution, declared gpt-6.1-sol/xhigh.
It accepted both placement amendments at 6c13b4b31eb4bd00b4a888e060f1c48a82f572c8 on all three axes with zero findings and confidence 9.
This was static confirmation of lane limits, callers, source readers, both canary families, anchors, command closure, and all 89 unchanged predicates.
Backend execution identity and usage were not observed.

Accepted spec SHA256: 7e9da0d934b5e1806d7c30721a306da1f9de53cb8f3bbdb494b3b324d02240f6.
The complete accepted Git source remains retained, including all earlier repair evidence.
No full-text source copy is created.
One retained gpt-6.1-sol/high author pass slices this accepted source within a two-iteration cap.
Independent graph review is accepted below; no runtime preservation or mutation proof is claimed.

## Graph allocation and first-use closure

The stable mapping is GP1 to GP1, GP2 to GP2, and GP3 to GP3.
The final exact-once Covers partition is 9/35/17 across all 61 unchanged row lines.
GP04's combined consumer predicate closes in GP2, while GP1 retains its first-use Git W4 witness.
GP47–GP61 close in GP2; GP3 repeats their actual-entry and complete decoder matrix cases under GP46 without duplicate Covers.
Every first-use omission, fixture, and independent expectation closes before its introducing checkpoint turns green.

Canonical FixturePins, ReferencingFiles, and BoundFiles return Writes counts 4/16/14 and a 25-path union.
GP2 owns both Git source-pinned canary units when its envelope migration first changes their mutation subject.
GP3 retains those original diagnostic predicates and their registered bite/restore observations.
No additional anchor holder or fixture unit is required by the actual ticket paths.

The source fence had 29 backticked bullets, of which 26 were actual path entries and three were exact hook-anchor phrases.
Canonical FenceTokens returned 31 tokens because two inline path citations also entered that section's authority.
The graph moves only closure commentary and exact hook-anchor bytes into their own section.
It preserves all 26 path-entry bytes and the three phrases while yielding the canonical 26-path fence, including the review pickup.
This mechanical normalization grants or removes no implementation path.

The first applicable guard ticket relocates only the existing 16-line checkPackageCoreAndGuards declaration.
The current source and destination counts remain 470 and 641, with the destination's existing 709-line grant.
Either outcome may land first; the second extends the sole relocated declaration and preserves the first outcome's subcheck.
Cross-spec shared writes require serial integrated-source charges, but GS1 has no unlanded projection API dependency.
Npm producers, their formatter call, and the independently authored path/count assertion remain unchanged.
No structure grant or accepted debt grows.

The version-1 inventory is planning-only.
Before any implementation charge, the coordinator binds actual version-2 run/session identities with author-limit 1 and fresh ticket authors.
The future implementation line stays gpt-5.6-sol/high, with GP3 and GS1 at accepted xhigh effort.
Complete degraded repair and reference capture precede GP1, while final family checks retain the original exclusions.

Static grammar, row/path preservation, canonical closure, prose, coverage, full reread, and planning-lane checks precede freeze.
Same-clean-pin plan-only and per-ticket proposals then validate the frozen graph without an implementation delta.
No implementation, runtime test, probe, benchmark, full gate, priority change, commitment change, or landing ran.

## Independent graph acceptance and stage close

Reviewed graph source: aa86c64d0e6227b4ba6d094518bbd295f329e0af
Reviewed spec SHA256: c69fd16e3b1a43b262a9f4f753280c0105eb5e643cbeafd6a4dd555ab4ad235f
Reviewed record SHA256: d0a2992b3ee0df90a904d6eb9a60bbceb2584b77136067348835048e20dc7b43
Reviewer: /root/deepen_tests_parsers
Declared reviewer line: gpt-6.1-sol/xhigh

| axis | material findings | confidence | disposition |
|---|---|---|---|
| Standards | 0 | 9 | accepted |
| Spec | 0 | 9 | accepted |
| Coverage | 0 | 9 | accepted |

The independent reviewer accepted both complete graphs on this frozen source.
All complete earlier reviewed Git sources and repair records remain retained; no redundant full-text artifact copy is created.
The close preserves every predicate, ownership entry, completion-plan byte, and all four ticket files.
Backend execution identity, provider usage, charges, and calibration remain unknown.

GP1, GP2, and GP3 retain all 61 predicates with exact-once Covers counts 9/35/17.
Their Writes counts remain 4/16/14 with a 25-path union inside the 26-path spec fence, including only the review pickup as the extra fence.
All three exact hook-anchor phrases remain unchanged.

GP1 closes its first-use Git wrapper and worktree-context witness before its green checkpoint.
GP2 closes the complete standard decoder matrix, seven actual core entries, and both Git canary families at first use.
GP3 repeats the complete decoder and actual-entry regression under GP46 without duplicate Covers.
The complete degraded repair and independent pre-migration reference capture precede GP1.

Either graph may perform the first applicable 16-line orchestration move into checks_test.go, currently 641 lines with its existing 709-line grant.
The second extends that sole declaration and preserves the first outcome's subcheck.
The npm formatter producer and its independent path/count reader remain in the original file.
Cross-plan shared writes require serial integrated-source charges and refreshed source and headroom.

The future implementation line remains gpt-5.6-sol/high, with GP3 and GS1 at the accepted xhigh effort.
Actual version-2 execution identity, author-limit 1, and fresh ticket authors remain pre-charge obligations.
These were static planning reviews; runtime preservation, decoder, cancellation, mutation, and omission witnesses remain unexecuted.
This stage close grants no implementation, integration, benchmark, publication, or landing authority.

Each changed Markdown file and the complete close delta are checked before one required planning-lane commit with the scorecard observations.
The source-owned handoff follows that commit and names $bench-mainenance as the exact next command.
No retrospective is created for this planning close.

<!-- command-currency: historical -->
