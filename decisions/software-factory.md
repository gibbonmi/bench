# Software factory

Status: shaping

## Destination

Bench will develop a traceable, recoverable local factory. It starts with a
structural refactor, a module-deepening refactor, and one CLI improvement. A
shared projection then supports a useful roadmap and execution view. Durable
execution and a qualified Regroup pilot follow through existing release
requirements.

## Notes

## Decisions so far

- [What direction and authority govern the factory?](software-factory/tickets/1.md): One factory roadmap organizes Bench and one later Regroup pilot.
- [What recurring maintenance cadence applies?](software-factory/tickets/2.md): FT243 checks introduced structure violations on each run and makes existing violations ranked proposals.
- [Which durable-execution principles apply?](software-factory/tickets/3.md): Use model-free routine supervision, durable pending decisions, unknown state with provenance, monitor continuity, typed lifecycle controls, and one shared projection.
- [What qualifies the external pilot?](software-factory/tickets/4.md): FT306 requires durable execution, FT71, and existing release qualification.
- [Which structural target starts the sequence?](software-factory/tickets/5.md): The structural-refactor-pass landing at `e86d344d` establishes this baseline.
- [Which module-deepening target should FT302 select?](software-factory/tickets/6.md): The reviewer approved the four-candidate deepening batch on 2026-09-05 from the `/bench-deepen` survey.
- [Which CLI improvement should FT303 select?](software-factory/tickets/7.md): The assessment of 2026-09-06 is recorded in `decisions/software-factory/assets/ft303-cli-assessment.md`.
- [Which improvement does FT303 select?](software-factory/tickets/10.md): The reviewer selected the mutation-probe verb on 2026-09-06.
- [Where does the probe verb live, and what bounds it?](software-factory/tickets/11.md): The verb is the root verb `bench probe`, and FT98's probe face folds into it.
- [Which focused-run forms and records does the probe verb carry?](software-factory/tickets/12.md): The focused run accepts both `bench test` selection forms: a package expression with `--run <go-regex>`.

## Not yet specified

- The targets selected by FT304 and FT305.
- The Regroup pilot change and its browser evidence.

## Spec-writer discretion

- Choose reversible internal placement after the selected target and source
  readers are verified.
- Choose a view layout that consumes the shared projection without becoming a
  state owner.
- Choose the preserved copy's location for the probe verb; it is reversible and
  changes no verdict.

## Out of scope

- A Firstmate runtime dependency, copied implementation, or runtime-capability
  claim.
- A release, deployment, or Regroup mutation before qualification.
- Automatic implementation from `/bench-deepen`.
- A `--time` face on `bench worktree exec` and a heredoc route to the exec
  child: the shell `time` prefix and the exec child's stdin serve both
  (2026-09-06).
- A probe verdict record in the census or a probe ledger: the terminal verdict
  row is the complete evidence (2026-09-06).

## Sources

- Path: `ROADMAP.md`
  Supports: owner rows, dependencies, release qualification, and holds.
  Drift: re-read before a factory target is specified or retired.
- URL: https://github.com/kunchenguid/firstmate/blob/main/docs/architecture.md
  Supports: shared projection and durable-pending-decision design input.
  Drift: use only as inspiration; it establishes no Bench capability.
- URL: https://github.com/kunchenguid/firstmate/blob/main/docs/watcher-continuity.md
  Supports: monitor-continuity design input.
  Drift: use only as inspiration; it establishes no Bench capability.
- URL: https://github.com/kunchenguid/firstmate/blob/main/docs/agent-control.md
  Supports: typed lifecycle-control design input.
  Drift: use only as inspiration; it establishes no Bench capability.
- Path: `decisions/software-factory/assets/ft303-cli-assessment.md`
  Supports: #7 and the #10 selection; the machine-local inventory at `~/.bench/assessments/bench-2826441890/ft303-cli-assessment/` holds the extraction.
  Drift: re-run the extraction when a probe verb, a production filter, or a structure path filter lands, or at the next CLI assessment.
