## Destination

Make spec-backed implementation use meaningful write subagents by default, and make
Bench phase conversation visibly structured without fragmenting related prose.

## #1: Which implementation runs must delegate?

Type: Grill

### Question

Define the mandatory delegation scope without making tiny changes pay more handoff
cost than implementation cost.

### Answer

Every spec-backed `/bench-implement-spec` run uses at least one genuine write
subagent. Independent vertical slices go to separate subagents in parallel within
the harness limit; dependent slices remain sequential. An atomic spec goes wholly
to one worktree-isolated write subagent. Read-only helpers do not satisfy the
requirement. Lighter-path changes may remain inline.

## #2: What happens when the harness cannot spawn a write subagent?

Type: Grill

### Question

Choose whether portability permits an inline fallback or makes the unavailable
implementation venue explicit.

### Answer

Stop before editing. The handoff states that mandatory delegation is unavailable
and gives explicit instructions to resume the same phase in a subagent-capable
harness, naming the exact invocation for that harness. There is no inline fallback.
The portable command contract and its gate checks are the available enforcement:
Codex does not expose `spawn_agent` on a deny-capable hook, so the kit does not claim
hard cross-harness runtime enforcement.

## #3: How is Bench phase conversation structured without becoming choppy?

Type: Grill

### Question

Set a shared rendering contract for progress and exit handoffs while preserving
cohesive prose.

### Answer

Substantial progress updates use compact bold **Status** and **Next** labels. Exit
handoffs use `## Result`, `## Details`, and `## Next` headings. Empty sections are
omitted, related sentences stay together, and bullets or tables are reserved for
genuinely parallel facts. The contract applies to Bench phase conversation, not CLI
output or repository artifacts.

## #4: Does this destination ship as one build slice or two?

Type: Grill

### Question

Choose whether the independent delegation and conversational-output contracts share
one spec.

### Answer

Use one decision map and two small specs/build slices: mandatory implementation
delegation first, shared conversational structure second. They have independent
owners and verification surfaces and remain separately reviewable.

## Not yet specified

n/a — the grill left no unresolved build-shaping fog.

## Out of scope

- CLI or TOON output schemas.
- Repository artifact formats such as specs, ADRs, decision maps, or reviews.
- Mandatory delegation for lighter-path implementation changes.
- Treating read-only research or review helpers as implementation delegation.
- A harness-specific runtime guard that the other shipped harnesses cannot observe.
- Model-tier bindings and the existing delegation charge, isolation, and verification
  rules.

## Handoff

1. **Module boundaries.** Slice one is owned by the canonical implementation phase:
   it decides when implementation delegates, while `craft-delegate` and `craft-line`
   continue to own how delegates are charged, isolated, verified, and routed. Slice
   two is owned once by the shared conversational-output rules; individual phase
   commands continue to own handoff content, not formatting. Harness adapters remain
   thin consumers when their session bootstrap loads the canonical command and shared
   rules. Root conformance and canary fixtures own regression detection for both
   guidance contracts.
2. **Contracts.** For a spec-backed implementation, delegation occurs before the
   first edit: one write subagent per independent slice in parallel, dependent slices
   sequentially, or the whole atomic build to one isolated write subagent. A
   read-only delegate does not count; a lighter-path run may stay inline. A harness
   without write-subagent support stops before editing and tells the reviewer where
   and exactly how to invoke the same phase next. Substantial phase progress renders
   non-empty **Status** and **Next** groups; exit handoffs render non-empty `Result`,
   `Details`, and `Next` sections, with cohesive prose and exact harness-native phase
   invocations. The conversation contract applies to sessions that load Bench's
   shared rules; safe-link continues to preserve a project-owned bootstrap that does
   not import them.
3. **Deep vs thin.** The implementation phase hides the mandatory venue decision
   behind one phase contract; delegation and line skills remain the deep how-to
   modules. The shared communication section is the single deep source for rendering;
   each phase's exit handoff is a thin declaration of required content. Adapters add
   no policy.
4. **Black-box assertables.** Root conformance can assert that the canonical
   implementation phase requires a genuine write subagent for spec-backed work,
   preserves the lighter-path exception, and contains the stop-with-explicit-resume
   route. It can assert that the active shared communication source declares a
   non-empty clause for each rendering and cohesion concern without copying those
   clauses into enforcement. Targeted canaries prove a declared clause bites when it
   disappears. Fresh-session dogfood can observe a
   spec-backed implementation spawning the required write subagent and can inspect
   phase progress and exit handoffs; the gate cannot honestly prove a historical
   spawn event.
5. **Gate attachment.** Attach both slices at the kit-content/root-conformance seam,
   with one behavior-owned canary family per contract. Because these edits steer
   generation, `craft-synthesis` also requires fresh-session dogfood after the gate is
   green. Runtime spawn enforcement remains unobservable to the portable gate and is
   stated as such rather than simulated.
6. **Hostile-input owners.** Paths with spaces/globs, control bytes, missing trailing
   newlines, shell argument quoting, symlink invocation, deep cwd, and CLI interrupt
   handling are n/a because neither slice introduces a shell input surface. Missing
   or empty canonical guidance is owned by root conformance and the canary fixtures.
   A missing subagent facility is owned by the implementation phase's pre-edit stop
   route. Every shipped harness session that loads Bench's shared rules is owned
   through its thin adapter reading the same canonical command and shared rules. A
   preserved project-owned bootstrap that does not load those rules is outside the
   conversational guidance claim. Re-entry after a stop is state-free;
   worktree reuse and interrupted delegate cleanup remain with existing delegation
   and worktree contracts.
7. **Uncertainty flags.** None. The reviewer accepted portable, gate-pinned generation
   guidance as the enforceable ceiling because hard cross-harness spawn policing is
   unavailable.
8. **Rejected alternatives.** Rejected mandatory delegation for lighter-path edits;
   counting read-only helpers; inline fallback on an incapable harness; duplicating
   output templates across every phase; forcing every sentence into a bullet; one
   combined spec; and claiming a portable runtime hook can observe spawning.
9. **Domain watch-outs.** Phase commands shape generation but do not emit an auditable
   spawn ledger. A structural gate can prevent the mandatory instruction from
   disappearing, while only fresh-session dogfood demonstrates that a harness follows
   it. Explicit resume instructions must name the invocation the destination harness
   actually supports, not a canonical slash command that may be a dead key there.

Dependency order: 1) mandatory implementation delegation; 2) shared conversational
structure. The slices are otherwise independent.
