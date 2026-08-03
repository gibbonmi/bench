# Check-level conformance scoping

Status: implemented

Decision source: reviewer-confirmed conversation on 2026-08-03; static source research only, with no Bench command, gate, test, build, or Go command executed during authoring

## Problem

Bench now batches a bookkeeping pass into one commit and treats that commit's gate
verdict as the only landing verdict, but every remaining doc-only tree still runs both
unconditional conformance surfaces. The first executes all dev conformance checks; the
second executes almost the entire conformance package. A one-file roadmap or capture
change therefore skips most components yet retains a fixed execution floor unrelated to
the file that changed. This draft claims structural work reduction, not an absolute
wall-clock saving; implementation evidence must compare the phase and per-check timing
records for equivalent full and scoped document fixtures before advertising a speedup.

The current narrowing foundation also has two correctness gaps that block further
optimization. An ambient conformance-check selector can reach the outer gate, and the
contract component's declared document inputs appear narrower than documents its tests
read. Making the selector more aggressive before closing those gaps would turn a cost
problem into a false-green problem.

## Solution

Keep one small conformance trust kernel always active, then give each ordinary
conformance check an exact-content identity and retained evidence derived from its
declared inputs. The gate, never `bench commit` arguments or an ambient selector,
decides which checks execute. All selected checks run in one conformance process;
unchanged checks inherit named evidence, and every failure to resolve or validate the
partition widens execution.

The first implementation closes the current selection and contract-input gaps. It then
registers every live-tree conformance assertion, establishes `conformance-meta`, makes
the ordinary conformance suite skippable only after no live-tree enforcement remains
hidden there, and finally adds per-check selection. Existing full boundaries stay full:
`bench gate --fresh`, prospective spec promotion, ship preparation, and canary-owned
mutation proofs do not accept the optimized partition.

## User stories

1. As a maintainer, I can trust that only the gate's internal decision narrows an outer conformance run, and that every document consumed by a component moves that component's identity. Line: `mid / high`, resolved through the active harness binding in `projects/benchkit.md`. These are correctness prerequisites at known gate and input-declaration seams.
2. As a maintainer, I see an always-executed `conformance-meta` result proving that the registry, declarations, ownership bindings, and executed-versus-inherited partition are complete before ordinary checks receive inherited credit. Line: `mid / high`, resolved through the active harness binding in `projects/benchkit.md`. The trust-kernel responsibilities are now precise, but a wrong implementation can authorize unchecked work.
3. As a maintainer changing documentation, I run only the ordinary conformance checks whose exact declared inputs moved, while unrelated checks inherit exact-content evidence and all selected checks share one process. Line: `mid / high`, resolved through the active harness binding in `projects/benchkit.md`. Per-check identity, aggregation, and invalidation are the genuinely uncertain seam.
4. As a reviewer, I can see which checks executed and which inherited evidence, while fresh, prospective, ship, and canary paths retain their full authority. Line: `mid / high`, resolved through the active harness binding in `projects/benchkit.md`. This combines oracle behavior with high-leverage profile and operator-facing claims.

## Implementation decisions

- The selector reads the gate's exact subject and retained check identities. Raw paths
  passed to `bench commit` are attribution input to that command, not gate authority:
  they may name directories and do not express every rename, deletion, or lifecycle
  mutation.
- `conformance-meta` is a logical always-run set inside the existing conformance
  invocation, not a third package-sized subprocess. Its verdict is visible as its own
  phase/evidence row even when it shares process startup with selected checks.
- Meta owns only authorization machinery: registry/function/tier bijection; exactly one
  named input declaration or explicit catch-all per ordinary check; profile bindings;
  check-to-canary ownership; prohibition of unregistered live-tree assertions; and a
  complete, disjoint executed/inherited partition. Semantic document, routing, package,
  and decision-map policy remains ordinary scoped conformance.
- Per-check input resolution stays in `internal/gate`, keyed by the lower conformance
  registry. The registry does not import gate types, preserving the existing dependency
  direction.
- The outer gate transports its ordered ordinary-check set in a gate-authored
  `BENCH_CONFORMANCE_CHECKS` value containing canonical comma-separated registry names.
  Phase construction strips any ambient singular or plural selector before setting that
  value. The conformance entry rejects duplicate, unknown, tier-invalid, or out-of-order
  names. The existing singular `BENCH_CONFORMANCE_CHECK` remains reserved for the
  authenticated inner-canary one-check control; neither variable is user authority.
- Every registered executable check binds its subject explicitly through the existing
  `(root, kitRoot, tier)` function seam. Migrated kit-policy checks such as component
  scope binding grade `kitRoot`; checks that grade an optional fixture surface guard its
  presence so a minimal fixture fails only for its planted reason.
- Each ordinary check identity binds its name and tier, the shared conformance
  implementation closure, its registry and function binding, its named input
  derivation and resolved content, its canary ownership, and the invocation/control
  schema. Uncertain declarations start as catch-all; optimism is not an efficiency
  feature.
- One green aggregate invocation authors slots only for checks it executed. A red
  invocation retires every executed check's slot and leaves inherited slots untouched.
  Meta never has reusable evidence.
- A declared symlink input binds the canonical in-repository target content, not only the
  link bytes. A broken link, a target outside the subject, or a target that cannot be
  resolved widens execution. Exact-file absence and a present empty file are distinct
  identity states.
- Unknown paths or selectors, missing declarations, failed derivations, malformed or
  mismatched slots, incomplete partitions, and unregistered live-tree enforcement run
  the affected checks or the complete set when attribution is impossible. They never
  produce an empty green.
- A conformance implementation change invalidates every ordinary check through the
  shared implementation closure. An owning-family canary change invalidates its check;
  an unbound or shared conformance helper conservatively runs every conformance canary.
- The current batching changes `e29de36` and `843a4b7` are closed baseline. This build
  reduces the cost of the one remaining gated landing; it does not add another commit
  mode or repeat their workflow prose.

## Testing decisions

- Test outer selection through the gate-owned phase construction seam, with an ambient
  selector present, and assert that only authenticated inner-canary execution may use
  the one-check control.
- Test document input identities through the real component declaration resolvers. A
  test that independently copies their path lists would preserve the same omission.
- Test meta and per-check partitioning at the gate decision seam with retained slots,
  then exercise the public gate path in hermetic fixture repositories for representative
  document classes.
- Keep the conformance invocation aggregated: tests observe one process receiving a
  gate-authored ordered check set, not one process per check.
- Every red signal below is intentionally unobserved during this planning pass because
  the reviewer prohibited gates, tests, builds, and Go commands. Each ticket must
  demonstrate its cited first red before implementation and record the targeted
  diagnostic.
- Compare the already-emitted phase and per-check timing records for equivalent full and
  scoped document fixtures. The acceptance claim is fewer executed checks and one
  aggregate process, not a host-dependent elapsed-time threshold.

### Seam diagram: gate-owned conformance partition

    trigger: ordinary dev gate over an exact Git subject
        │
        ▼
    subject + declarations + retained slots
        │
        ▼
    [ gate-owned conformance partition ]
        │                     │
        ▼                     ▼
    always-run meta      selected ordinary checks
        │                     │
        └──────────┬──────────┘
                   ▼
          one conformance invocation
                   │
                   ▼
       executed + inherited evidence
          ◀ tests attach here: injected subjects, declarations, slots, and fixture-visible output

### Seam diagram: component document identity

    trigger: resolve a component identity for one tree
        │
        ▼
    authoritative consumer inventory
        │
        ▼
    [ component input derivation ] ──▶ sorted paths + digests ──▶ identity
          ◀ tests attach here: mutate a real managed document and observe identity movement

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | an ambient `BENCH_CONFORMANCE_CHECK` cannot narrow an outer gate | gate-owned phase construction | Unobserved by reviewer instruction; seed a valid ambient check in an outer runtime fixture and expect all dev checks, with a focused contract red today. | A scrub performed only by callers or a selector trusted from ambient state still executes one check. |
| 1 | authenticated inner-canary execution can still select exactly its owning check | gate-owned phase construction | Unobserved; run the inner fixture control after sealing the outer path and require its one-check timing record. | Deleting the selector entirely makes the outer path safe by breaking canary attribution. |
| 1 | every managed document read by the contract component moves its derived identity, including `.bench/BENCH.md` | component input derivation | Unobserved; mutate each path from the authoritative managed-asset inventory and assert the contract identity changes. | A copied `.agents/`-only list passes ordinary guidance cases while leaving the observed Bench guide undeclared. |
| 1 | a failed component input derivation runs the component | component decision | Unobserved; inject resolver failure and require the component in the executed set. | Treating failure as an empty input set grants the unsafe skip. |
| 2 | every executable live-tree conformance assertion is registered exactly once as meta or ordinary | conformance-meta | Unobserved; add a prefix-named live-tree test without a registry row and require the hidden-enforcement diagnostic. | The current suite-only reachability failure survives registry/function bijection alone. |
| 2 | every ordinary check has one named derivation or explicit catch-all and every name resolves | conformance-meta | Unobserved; remove one declaration and misspell another resolver name, expecting distinct diagnostics. | A default empty declaration or free-text provenance silently authorizes inheritance. |
| 2 | meta proves the executed and inherited sets are disjoint and cover every tier check | conformance-meta | Unobserved; omit, duplicate, and overlap one registry row in three mutations and require each partition diagnostic. | Aggregate green is vacuous if one registered check appears in neither set. |
| 2 | profile tables for reduced scope, component inputs, and check inputs exactly match their sources | conformance-meta | Existing recorded bite proofs cover the reduced-scope and component-input tables; unobserved here, add and remove one check-input token on each side and require the new row's stale diagnostic. | One-direction subset checks allow either implementation or advertisement to drift. |
| 2 | meta always executes and has no reusable evidence slot | gate evidence | Unobserved; present a forged meta slot and require meta execution plus slot rejection. | Treating meta like an ordinary component lets broken scoping skip its own authorization. |
| 3 | every declared document input entry executes its reached checks and inherits only exact matching evidence for the rest, including canonical targets reached through symlinks | gate-owned conformance partition | Unobserved; derive the fixture matrix from every declared file entry and directory class; mutate and delete each exact file, add and remove a descendant in each directory class, mutate each in-repository symlink target, and require broken or escaping links to widen. | A hand-picked path matrix, link-byte identity, or blanket Markdown fast path leaves an unrepresented declaration free to skip wrongly. |
| 3 | selected checks run in registry order in one process through the gate-authored ordered-set transport | conformance invocation | Unobserved; select two non-adjacent checks in reverse input order and assert one process plus registry-ordered timing rows; inject ambient, duplicate, unknown, and tier-invalid plural values and require rejection or gate replacement. | Process-per-check keeps correctness but loses the intended runtime win; caller order or ambient transport makes selection unstable or unauthorized. |
| 3 | no prior slot, malformed evidence, wrong check or tier, changed shared implementation, and failed derivation all widen execution | check evidence resolution | Unobserved; mutate each evidence field and the shared closure independently and require execution. | A field omitted from identity or validation produces a reusable false green. |
| 3 | green authors only executed slots, red retires executed slots, interruption authors no new slots, and inherited slot bytes remain unchanged | check evidence persistence | Unobserved; drive green, red, interruption before persistence, and repeated unchanged partitions across fresh processes and compare records. | Re-stamping inherited evidence hides age and red invalidation; partial authorship grants unearned credit; retiring inherited evidence repays unrelated work. |
| 3 | `conformance-suite` declares its module-test closure and manifest and becomes skippable only after every live-tree assertion leaves the hidden suite | conformance-suite component decision | Unobserved; bind the suite to the declared component-input source, then retain one unregistered live-tree assertion and require meta red before the suite can inherit. | Adding an unnamed module-closure slot first would hide the exact checks that currently justify the suite's unconditional run. |
| 4 | output and the durable verdict enumerate executed checks and the evidence covering each inherited check | public gate fixture and verdict record | Unobserved; inspect a mixed partition and require complete, non-overlapping names and evidence identities. | A summary count cannot attribute which policy was actually graded. |
| 4 | `bench gate --fresh`, prospective promotion, and ship execute all applicable checks and ignore check slots | public gate fixture | Unobserved; seed reusable slots at each entry and require the full check inventory. | Allowing lifecycle-final or release paths to inherit changes their existing authority. |
| 4 | a check implementation change runs its owning canary family, while an unbound shared helper runs all conformance canaries | canary ownership decision | Unobserved; mutate one bound check and one shared helper and inspect the selected canary set. | Re-running a weakened check while inheriting its old bite proof defeats the tripwire. |
| 4 | scoped document execution records fewer ordinary check timings than the equivalent full fixture without claiming a fixed wall-clock threshold | timing evidence projection | Unobserved; compare the existing phase and per-check timing records for equivalent full and scoped fixtures and require the scoped record to omit inherited checks while retaining one aggregate invocation. | A selector that still executes every check can look architecturally complete while delivering no work reduction. |

### Edge inventory

- Error path: declaration, identity, slot, and persistence failures widen or red through rows 4, 7, 9, 12, and 13.
- Empty/absent input: no prior evidence and zero reached ordinary checks are covered by rows 7, 9, and 12; exact-file absence versus present-empty identity is covered by row 10; zero executions is valid only when every ordinary check inherits valid evidence.
- Boundary values: one check, every check, one path observed by multiple checks, and a catch-all declaration are covered by rows 6, 10, and 11.
- Malformed input: unknown check, wrong tier, malformed slot, and unknown derivation are covered by rows 6, 11, and 12; broken or out-of-subject symlinks widen under row 10.
- Interrupted or partial state: interruption between aggregate execution and slot authorship is covered by row 13; no partial authorship may be credited.
- Re-run idempotency: repeated unchanged partitions preserve inherited evidence bytes under row 13.
- Process-boundary lifecycle: fresh-process evidence reload is covered by row 13, and prospective promotion by row 16.
- Hostile environment: ambient selector injection is covered by row 1; unavailable tools and failed derivations widen under rows 4 and 12.
- Rename/delete: renamed checks and deleted exact-file or directory-descendant inputs are covered by rows 5, 6, 10, and 12.
- Symlink inputs: `.claude/` adapters and skills resolve to their canonical `.agents/` targets; target mutation moves identity, while broken or escaping targets widen, under row 10.
- Paths with spaces and special files: declaration resolution must either hash the exact regular file or widen; covered by rows 4 and 12.
- Two doc changes graded separately and later composed: exact-content identities, not prior changed-path lists, decide the composed partition under rows 10, 12, and 13.

**Won't handle:** adversarial writes to repository-global Git evidence — this retains the existing trusted-local evidence posture and is a separate security capability, estimated at more than 8 edits and 2 full gate runs.

## Out of scope

- FT183 landed first: `b41b4d2` removed the reduced fallback and record class, and
  `6a3ea99` added the resolver-identity proof. This build consumes that baseline rather
  than reopening its `internal/gate`, evidence-record, and profile decisions.
- Narrowing contract or canary below their current component granularity for `.agents/` edits: separate consumer-specific evidence capability, estimated at 10-15 edits and 3 full gate runs. This build only corrects the contract input inventory and binds canary invalidation to conformance checks.
- A generic doc-only allowlist or plain-Git commit exemption: deliberately rejected, because Markdown includes executable policy and gate-anchored claims and would duplicate consumer knowledge.
- A new `bench commit` mode: the landed batching and verdict-reuse guidance already reduce landing frequency; this build changes what the existing oracle executes.
- Declarative prose-anchor mechanics owned by FT156 remain separate. FT183's named-resolver binding is a prerequisite consumed here rather than duplicated; checks without an exact bound resolver use catch-all.
- Focused canary porcelain owned by FT168: this build selects owning evidence internally and does not add a public fixture selector.
