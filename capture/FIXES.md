# Prioritized roadmap fixes

Assessed 2026-08-13 against `main` at `7d22a29f`. This is the fix-only triage
view of the 24 live defect rows in `ROADMAP.md`. Features, decision-only
items, parked work, and already-shipped rows do not appear here. The roadmap
remains the canonical board, and every fix still needs live revalidation at
build entry.

## Reading the queue

- **Applicable** means the current tree still exhibits the stated defect.
- **Partial** means part of the defect shipped or disappeared; re-slice the live
  remainder before implementation.
- **Decision** means the defect is evidenced, but its repair posture still needs
  the named reviewer choice.

Priority favors a trustworthy oracle and false-green removal, then lifecycle
correctness, then bounded reliability and ergonomics. Dependency order remains
binding: FT98 precedes FT169.

The overall roadmap starts with feature FT198, which is intentionally absent
from this file. FT189 is the highest-ranked live fix.

## Ordered fixes

| Rank | ID | Verdict | Current assessment and next disposition |
| ---: | --- | --- | --- |
| 1 | FT189 | Applicable | An upstream `git worktree list` hang can stall every Bench command that performs discovery; decide the Bench-owned refusal or execution bound. |
| 2 | FT142 | Partial | Release publication still bypasses the governed resumable publisher. Revalidate the surviving FT91 residuals, then route publication through the current owner. |
| 3 | FT133 | Applicable | `bench coverage --check` can credit acceptance evidence without proving that every cited test executes or that each row has one accountable producer. |
| 4 | FT174 | Applicable | Ticket metadata can produce dependency and ownership false-greens; enforce the declared graph and exact producer fence. |
| 5 | FT177 | Partial | Branch-native build/test narrowed the stale-`dist/bench` surface, but direct contract mutation probes can still select an old binary. Reproduce the remaining public path before editing. |
| 6 | FT103 | Applicable | The packed-asset check verifies exclusions and a non-empty allowlist without proving every allowlisted source path exists. |
| 7 | FT201 | Applicable | Production cancel-signal registrations lack one conformance rule and mutation that detects an omitted cleanup path. |
| 8 | FT98 | Applicable | Preserve-then-discard behavior remains split across lifecycle commands. Establish one primitive before FT169 consumes it. |
| 9 | FT169 | Applicable | No single sanctioned worktree landing command owns preservation, integration, and cleanup. Blocked by FT98. |
| 10 | FT141 | Partial | The output-truncation guidance face is closed, but gate-pin red records still lack failing-phase inventory and stable attribution. |
| 11 | FT178 | Applicable | Bare `bench worktree` remains human porcelain with behavior too implicit for automation. |
| 12 | FT173 | Decision | An active assignment whose tree is missing receives misleading actions and no cleanup disclosure; choose the intended disclosure class before the light-path repair. |
| 13 | FT190 | Partial | Injected-port conformance covers only part of the interface inventory. Re-slice to uncovered packages and require a real-producer test or explicit exemption. |
| 14 | FT92 | Applicable | Subject-drift attribution and shipped-input hygiene remain incomplete across consumers of retained and landed state. |
| 15 | FT130 | Applicable | A capture write during an active lifecycle can invalidate the subject without a mechanical void-or-block response. |
| 16 | FT117 | Applicable | The two FT87 parser-surface defects remain focused fixes and should stay independent from broader parser redesign. |
| 17 | FT172 | Partial | The drain now captures one full trusted snapshot, but roadmap row grammar and mechanical discrepancy coverage remain incomplete. |
| 18 | FT113 | Applicable | `bench commit --spec` does not treat its owned spec-state transition as satisfying the required path set. |
| 19 | FT182 | Applicable | A Planned-phase receipt over an absent target can wedge the lifecycle instead of refusing or repairing the state. |
| 20 | FT58 | Applicable | Pool-root selection still trusts ambient ownership and permissions rather than validating them. |
| 21 | FT104 | Applicable | Load-induced commit refusals lack a stop rule and cheap pre-gate discriminator, so transient pressure can be reported as a product defect. |
| 22 | FT115 | Partial | The original wait sites are gone and marker bounds exist; re-census only surviving unexplained deadline literals. |
| 23 | FT120 | Partial | Most named harness files disappeared with branch-native and lifecycle subtraction; reconcile the row to live FIFO, timeout, and isolation defects. |
| 24 | FT125 | Applicable | Narrow readers still over-return context: `bench outline` has no symbol selector and there is no `bench spec show`. |

Every **Partial** row needs a fresh problem statement and exact live red signal
before ticketing. Do not carry vanished file names, retired lifecycle
assumptions, or already-closed clauses into a build fence.
