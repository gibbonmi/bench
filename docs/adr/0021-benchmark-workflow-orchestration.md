# Benchmark workflow outcomes remain independent under shared assurance

Status: accepted

The benchmark workflow delivers ticket implementation, continuation, completion evidence, and assessment as independently useful outcomes. Each outcome can land separately, while the shared assurance rules below keep their behavior compatible.

## Consequences

- Each ticket gets a fresh author session on the declared line, as ADR 0023 records. That session owns the production changes, tests, and probes of its ticket, and a fresh repair session owns each repair. Specs define stable review chunks with outcomes, acceptance rows, and tests; independent Standards, Spec, and Coverage review closes each chunk before its successor starts. Review uses mid/high by default and top/high for a mid-tier Sol implementation. A model change, a tier move, or a session beyond the planned ticket authors requires reviewer direction.
- The orchestrator of ADR 0023 can change chunk boundaries, expected writes, and gate coverage inside approved behavior. The orchestrator first updates the plan, preserves guarantees, and records the expansion for the learning drain.
- Explicit continuation keeps each ticket author of ADR 0023 active while it makes progress. After two attempts without progress, the author reassesses and can use debugging or brief read-only diagnosis. These actions do not transfer authorship or override a stop or approval boundary.
- Completion evidence binds each checkpoint and final result to the exact source, obligation, and performer. Missing, failed, or stale evidence blocks advancement or completion through the gate-owned check.
- A reviewer can opt one full run into delegated authorship. That run adds concurrent authors and the full tier range to the fresh ticket authors of ADR 0023. It keeps the implementation chunk as the review scope and integrates contributions serially. Every review axis then uses the invoking harness's configured mid tier at high effort. Without the opt-in, the ticket authors work serially on the declared line.
- Assessment records complete workflow cost and quality locally. Paid comparisons require an approved plan, and only the reviewer can adopt a different default from repeated equivalent-assurance evidence.
