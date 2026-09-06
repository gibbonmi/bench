# Software factory

Status: shaping

## Destination

Bench will develop a traceable, recoverable local factory. It starts with a
structural refactor, a module-deepening refactor, and one CLI improvement. A
shared projection then supports a useful roadmap and execution view. Durable
execution and a qualified Regroup pilot follow through existing release
requirements.

## #1: What direction and authority govern the factory?

Blocked by: none
Type: Grill

### Question

What work does the software factory organize?

### Answer

One factory roadmap organizes Bench and one later Regroup pilot. Bench retains
its gate, release, and project-owned policy. Regroup retains its domain and
design authorities. The factory does not select new product direction or permit
implementation without a target-specific approval.

## #2: What recurring maintenance cadence applies?

Blocked by: none
Type: Grill

### Question

How does recurring maintenance select work?

### Answer

FT243 checks introduced structure violations on each run and makes existing
violations ranked proposals. Module-deepening assessment runs weekly. CLI and
other deep assessments start biweekly and become monthly only after successive
runs find little actionable friction. The policy favors seams and tracer-bullet
slices, not count splitting.

## #3: Which durable-execution principles apply?

Blocked by: none
Type: Grill

### Question

Which principles guide the future execution and view work?

### Answer

Use model-free routine supervision, durable pending decisions, unknown state
with provenance, monitor continuity, typed lifecycle controls, and one shared
projection. Existing authoritative readers remain state owners. These principles
create no Firstmate dependency and make no runtime verification claim.

## #4: What qualifies the external pilot?

Blocked by: none
Type: Grill

### Question

What must precede a Regroup factory pilot?

### Answer

FT306 requires durable execution, FT71, and existing release qualification.
The pilot is one user-visible Regroup change with browser behavior and visual
evidence. It begins only after the local workflow is inspected and a change is
selected.

## #5: Which structural target starts the sequence?

Blocked by: #1
Type: Research

### Question

Which selected structural pass establishes the baseline?

### Answer

The structural-refactor-pass landing at `e86d344d` establishes this baseline.
Its retained retro records the survey, accepted scope, gate results, and cost.
The next target requires a separate module-deepening survey.

## #6: Which module-deepening target should FT302 select?

Blocked by: #5
Type: Research

### Question

Which post-structural survey target will hide behavior behind a useful seam?

### Answer

— (open)

## #7: Which CLI improvement should FT303 select?

Blocked by: #6
Type: Research

### Question

Which grouped calls, existing verbs, and attributed accessible sources justify
one CLI improvement?

### Answer

— (open)

## #8: What does the shared view project?

Blocked by: #7
Type: Research

### Question

Which existing readers provide the useful observation view without a second
work-state owner?

### Answer

— (open)

## #9: What execution target fulfills the settled principles?

Blocked by: #8
Type: Research

### Question

Which existing execution seams support durable continuity and typed lifecycle
control?

### Answer

— (open)

## Not yet specified

- The targets selected by FT302 through FT305.
- The Regroup pilot change and its browser evidence.

## Spec-writer discretion

- Choose reversible internal placement after the selected target and source
  readers are verified.
- Choose a view layout that consumes the shared projection without becoming a
  state owner.

## Out of scope

- A Firstmate runtime dependency, copied implementation, or runtime-capability
  claim.
- A release, deployment, or Regroup mutation before qualification.
- Automatic implementation from `/bench-deepen`.

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
