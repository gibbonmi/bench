# Scope and report document conformance

Blocked by: derive-complete-contract-document-inputs.md, retain-check-level-evidence.md
Ownership fence: `internal/gate`, `internal/conformance`, `internal/canary`, `internal/contract/runtime`, `projects/benchkit.md`, `.bench/BENCH-reference.md`
Assumptions: every live-tree assertion is registered before the suite can skip; selected checks share one process; fresh prospective and ship paths remain full; claims re-derived from the tree at pickup

## What to build

Ordinary doc-only gates execute meta plus reached conformance checks in one process,
inherit the rest, skip the ordinary suite only when its inputs are unchanged, and expose
the exact partition without narrowing any terminal oracle path.

## Acceptance

- [ ] [DS1] Every declared document input entry executes its observing checks/components: each declared file and directory class, including `ROADMAP.md`, capture, spec, decision, `.agents/`, `.claude/`, `.codex/`, `docs/adr/`, Bench-guide, profile, and README classes.
- [ ] [DS2] Selected ordinary checks execute once in registry order in one process through gate-authored `BENCH_CONFORMANCE_CHECKS`, while unchanged checks carry exact evidence and ambient or invalid plural selectors cannot narrow the set.
- [ ] [DS3] `conformance-suite` declares its module-test closure and manifest in the component-input source and profile table, and inherits for doc-only trees only after meta proves no hidden live-tree enforcement remains.
- [ ] [DS4] Output and durable verdicts enumerate every executed check and the evidence covering every inherited check.
- [ ] [DS5] Fresh, prospective, ship, and canary mutation paths execute all checks their existing authority requires.
- [ ] [DS6] A bound check implementation change runs its owning canary family, and an unbound shared helper runs all conformance canaries.
- [ ] [DS7] In-repository symlink declarations bind target content; target mutation moves identity, exact-file absence differs from a present empty file, and broken or out-of-subject targets widen execution.
- [ ] [DS8] Equivalent full and scoped document fixtures expose phase and per-check timing records showing fewer executed ordinary checks in one scoped aggregate process, without a fixed elapsed-time threshold.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| DS1 | omit one declared file entry or directory class, or map all Markdown to one generic fast path | registry-derived public gate fixture matrix | derive the matrix from declarations; mutate and delete every exact file, add and remove one descendant of every directory class, expect each executed/inherited partition mismatch |
| DS2 | launch one process per check, preserve caller order, or trust ambient plural selection | conformance invocation contract | select non-adjacent checks in reverse and inject ambient, duplicate, unknown, and tier-invalid values; expect one launch, registry-ordered timings, and rejection or gate replacement |
| DS3 | give the suite an unnamed closure or make it skippable while one hidden live-tree test remains | conformance-meta and suite contract | compare its component/profile declaration, add the hidden assertion, drive a doc tree, expect a stale-binding or meta diagnostic and suite execution |
| DS4 | report counts without per-check evidence | verdict projection contract | inspect a mixed run, expect every registry name exactly once with evidence on inherited rows |
| DS5 | honor check slots under fresh, prospective, or ship mode | public terminal-path contracts | seed slots at each entry, expect the complete applicable check inventory |
| DS6 | inherit old canary evidence after changing a check or shared helper | canary ownership contract | mutate each source class, expect the owning or complete conformance canary set respectively |
| DS7 | hash only symlink bytes, collapse absent and empty files, or accept a broken or escaping target | declared-input hostile fixture | mutate the canonical target without changing link text, compare absent with present-empty, then break and escape the link; expect identity movement or conservative execution in every case |
| DS8 | report a scoped partition while still executing every ordinary check | timing evidence contract | compare equivalent full and scoped fixture records, expect inherited names absent from scoped timing rows and exactly one aggregate conformance invocation |
