# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-07-28 — Wide refactor was drafted as one ticket  [open]

- **What happened:** FT154's first folder-layout ticket grouped a deep-unit
  change with every runtime, AXI, ambient, lifecycle, and conformance consumer.
  The delegate made the core green but exhausted its fresh context before the
  fixture migration closed, so the ticket could not land independently green.
- **Right behavior:** Classify blast-radius work before drafting tracer-bullet
  tickets. A wide refactor takes expand–migrate–contract sequencing; migration
  tickets are sized by an ownership fence such as package or consumer family,
  and the contract ticket is blocked by every migration.
- **Proposed rule change:** Make the wide-refactor test an explicit first branch
  of `craft-tickets`' breakdown procedure, before vertical ticket drafting.

## 2026-07-28 — Ticket progress duplicated the gate oracle  [open]

- **What happened:** FT154 tickets carried a `project gate is green` checkbox,
  so the coordinator ran `bench gate`, checked the box, then invoked
  `bench commit`, which ran the same full gate again for the same ticket.
- **Right behavior:** Iterate with focused seam checks, then let the atomic
  `bench commit` run the one full project gate that proves the ticket lands
  independently green. `/bench-final-check` still runs the full gate over the
  composed feature.
- **Proposed rule change:** Keep gate state out of the ticket template because
  the green landing commit is its one source; define ticket cadence as focused
  checks during iteration, one full gate at commit, and one final composed gate.

## 2026-07-28 — Full gate hides a multi-minute critical path  [open]

- **What happened:** FT154's first ticket under the new cadence spent 2.64
  seconds on focused checks, then 163.98 seconds in the atomic full gate.
  Concurrent output showed the contract phase's
  `internal/contract/surface/artifact` package at 147.91 seconds and the core
  test phase near 60 seconds, but the gate summary reported only green/red, so
  phase overlap and the true critical path had to be reconstructed by hand.
- **Right behavior:** A project this size should not accept a multi-minute dev
  gate as ambient cost. Diagnose phases one at a time with no competing gate or
  Go test processes, compare isolated normal and single-threaded runs, and
  preserve the oracle until evidence identifies duplicated work, contention,
  or an intrinsically slow fixture. Here, the isolated contract phase still
  took 120.71 seconds and the artifact package alone took 109.81 seconds in a
  normal profiled run. Forcing that package to one execution thread made it
  worse at 168.35 seconds. The core-test registry already excludes contract
  packages, so this is not a duplicate package invocation: the dominant cost is
  repeated full artifact generation inside one sequential package, with about
  28 seconds of additional contention when gate phases overlap.
- **Proposed rule change:** Make the gate report total and per-phase wall time
  in a concurrency-aware form, and treat a sustained phase or critical-path
  budget breach as a named performance problem rather than invisible green
  latency. Run only one gate orchestration at a time, but retain useful
  toolchain parallelism inside a phase. Review the artifact contract's existing
  prepared-generation seam for safe immutable fixture reuse; keep distinct full
  builds only where the build environment itself is the behavior under test.

## 2026-07-28 — `bench gate --help` started the gate  [open]

- **What happened:** While looking for a phase-inspection surface during FT154
  debugging, the coordinator ran `bench gate --help`. The command ignored the
  unrecognized flag and started a real concurrent gate, which continued after
  the shell call yielded and had to be found by PID and cancelled before
  isolated measurements could begin.
- **Right behavior:** Treat any undocumented `bench gate` argument as capable
  of executing the oracle; inspect command source or documented command
  listings instead of probing help flags. A help request must never trigger
  expensive or stateful work.
- **Proposed rule change:** Make `bench gate --help` print usage without running
  the gate, and make unknown gate arguments refuse with exit 2; pin both
  postures in the runtime CLI contract.

## 2026-07-29 — Implementation retro does not explicitly seek codification candidates  [open]

- **What happened:** Reviewing the implementation-retro prompt showed that it
  asks for coordinator catches and CLI, skill, and process improvements, but
  does not explicitly inspect the completed session for ad-hoc checks, decision
  procedures, or reconstructed logic worth codifying. Finding those candidates
  is therefore incidental.
- **Right behavior:** A spec-backed run's retro should deliberately inspect its
  session evidence for logic that was repeated, reconstructed, or required
  coordinator judgment, then name the appropriate durable owner such as the
  Bench CLI, a skill, the gate, or process guidance.
- **Proposed rule change:** Add an explicit codification-candidate pass to
  `/bench-final-check`'s implementation-retro instructions, requiring each
  candidate to name the session evidence, proposed owner, and expected effect.
