# Greenfield build sequence

This is the dependency-ordered plan for building the complete Bench kit in a
new repository. It is a planning reference, not the live roadmap. The plan
assumes two engineers or agent workers plus a reviewer and uses twelve
two-week sprints, for roughly six months of elapsed time.

The implementation starts with the Go core. Shell remains only where the
operating environment requires it: launchers, hooks, and harness adapters.
This avoids building a large shell implementation that must later be replaced.

The dependency spine is:

> product contract → gate → workflow content → adoption → enforcement →
> isolated shifts → artifact lifecycle → query surfaces → ambient feedback →
> diagnostics and maintenance → distribution → release hardening

## Sprint plan

| Sprint | Outcome | Scope |
| --- | --- | --- |
| **1 — Contract and walking skeleton** | A runnable, testable CLI with the architecture settled. | Lock the harness-agnostic design, ubiquitous language, four invariants, seams, output and exit-code conventions, and one-source rule. Establish the Go binary, thin shell launcher, command dispatcher, version and help surfaces, shared git/subprocess/terminal utilities, TOON emitter, and unit-test foundation. |
| **2 — The gate becomes the oracle** | Every later sprint has an authoritative definition of done. | Implement the gate command, resolution and phase execution, tree hashing, verdict cache, adversarial pinning, build/vet/test checks, and the first conformance fixtures. Add a minimal canary proving a deliberately broken tree goes red. |
| **3 — Planning workflow and skills** | Bench can guide feature and bug work before autonomous execution exists. | Create the operating guide and portable working agreement; add the shape, spec, implement, review, final-check, and debug phases; add the core craft skills, decision-map/spec/review formats, acceptance-coverage-map convention, and generated skills index. |
| **4 — Safe adoption** | A fresh repository can install and remove Bench without losing project-owned content. | Build link, init, unlink, and doctor; managed-block preservation; the link manifest; conflict detection; copy and symlink modes; local CLI installation; project profiles; the setup phase; and Claude, Codex, and generic AGENTS.md adapter layouts. Gate fresh install, relink, and uninstall in throwaway repositories. |
| **5 — Enforcement and line routing** | Safety rules behave consistently across harnesses. | Add destructive-git protection, Stop-hook gate enforcement, pre-push protection, SessionStart plumbing, guard self-description, owner-defined line bindings, model-ID validation, model discovery and resolution, harness aliases, shift adapters, and line enforcement where each harness supports it. |
| **6 — Isolated autonomous execution** | Bench has its end-to-end MVP: a bounded shift that commits only on green. | Implement the worktree pool, atomic leases, stale recovery, cleanup and salvage handling, bounded shifts, iteration notes, interrupt cleanup, gate-after-each-iteration behavior, agent-touched-path staging, path-scoped commits, and green-only commits. Record the pre-shift base for later review. |
| **7 — Artifact lifecycle** | Decisions, specs, reviews, and implementation state have a coherent lifecycle. | Add closed-map Handoffs, staged/implemented/retired specs, coverage-map validation, durable review findings, implementation recording, spec retirement, and spec history. Gate command handoffs, workflow anchors, and stale artifact references. |
| **8 — Agent-facing state surface** | Agents can inspect repository state through stable, machine-oriented commands. | Complete learnings, maps, guards, diff, and coverage queries; recorded-base and after-merge diff modes; guard aggregation; definitive empty states; structured errors; TOON escaping and control-byte rejection; exit 0/1/2 contracts; and AXI conformance tests. |
| **9 — Capture, roadmap, and ambient dashboard** | Cold sessions immediately see the most important next action. | Build idea capture, the working roadmap, the learnings journal, the what-next maintenance phase, and the ambient dashboard. Add the severity ladder, five-row budget, gate-cache staleness, combined capture-drain signal, worktree/spec/review/roadmap signals, the overflow view, and SessionStart injection. |
| **10 — Diagnostics and self-maintenance** | The kit can explain and maintain itself. | Add structure budgets and accepted grants, the seam outline, the standalone dashboard page, richer model diagnostics, full canary commands, the update-kit synthesis phase, the assessment phase, assessment ownership, ADR discipline, changelog and provenance duties, and migration examples. |
| **11 — Distribution** | A new user can install the same product the maintainers dogfood. | Establish the redbench npm identity, platform packages, cross-compilation, the release workflow, binary-repair fallback, durable PATH shim, postinstall behavior, npx-from-git support, package-surface checks, and complete install and uninstall guidance. Produce a public beta candidate. |
| **12 — Release hardening** | The complete system is trustworthy across supported environments. | Exercise hostile paths and names, control bytes, missing newlines, absent versus empty files, missing tools, symlink invocation, nested working directories, every shipped surface, interrupts, and reruns. Stress leases and concurrent gate phases, serialize binary builds, expand canaries so checks demonstrably bite, verify harness parity, measure gate latency, and dogfood fresh installs. |

## Milestones

- **After Sprint 4:** installable planning-workflow alpha.
- **After Sprint 6:** autonomous gated-shift MVP.
- **After Sprint 9:** operational beta with state and feedback surfaces.
- **After Sprint 12:** complete product ready for a public release.

## Conditional work

Evidence-dependent, upstream-blocked, or deliberately tabled work does not belong
in a committed sprint:

- Give the standalone dashboard page a visual-identity sprint only after its
  product decisions are reopened.
- Add Codex line-guard parity only when Codex exposes a deny-capable event with
  the resolved delegate model.
- Strengthen lease reclaim locking and per-anchor canary coverage when their
  recorded evidence thresholds fire.
- Revisit tier bindings on their scheduled date or after a frontier-model shift.

Structure debt should not accumulate into a cleanup sprint. Enable structure
budgets early and repair or explicitly accept growth in the sprint that causes it.
