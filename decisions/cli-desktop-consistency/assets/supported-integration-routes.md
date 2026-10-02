# Supported CLI and desktop integration routes

Use the existing Bench integration with checks tied to each actual harness interface.
Supported hooks provide lifecycle entry points, but hook success cannot establish that a chat tool can start.
The installed probe gives a negative command result; full workflow compatibility remains unverified.

Date: 2026-10-02
Consumed by: decision ticket #12 and the subsequent specification
Drift: Refresh affected claims after changes to Bench integration, upstream documentation, runtime identity, effective configuration, or workspace identity.
Retire when: A reviewed specification and current qualification evidence supersede these findings.

## Question graph

| Node | Factual question | Prerequisite |
|---|---|---|
| Q1 | Which environments and configuration sources govern each interface? | Installed probe |
| Q2 | Which supported lifecycle routes can load and check Bench? | Q1 |
| Q3 | What recovery remains available when the chat cannot start a command? | Q1 and Q2 |
| Q4 | What does the current Bench integration prove? | Q2 and Q3 |

The report uses official upstream documentation and the Bench source tree.
Local source references describe commit `6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69`.
The shaping worktree adds documents but does not change those integration sources.

## Q1: Environment and configuration facts

The Windows app supports an agent in WSL2.
Its agent environment and integrated terminal environment are separate settings.
An agent-environment change takes effect after an app restart.
The Windows app and WSL CLI use separate configuration homes by default.

Source: [Windows app documentation](https://learn.chatgpt.com/docs/windows/windows-app), retrieved 2026-10-02.

These facts match the separate homes in the [installed probe](installed-compatibility-probe.md).
The approved outcome preserves personal settings and credentials, as [ticket #6](../tickets/6.md) requires.
The upstream option to share a whole configuration home does not implement that boundary.
This is an inference from the documented sharing behavior and the approved scope.

## Q2: Lifecycle and trust facts

Codex discovers repository hooks beside active configuration layers.
Project hooks require project trust, and non-managed hooks require trust for their exact definitions.
Hooks are enabled by default in the current documentation.
A missing explicit feature flag therefore does not establish that hooks are disabled.

`SessionStart` distinguishes startup and resume events.
`PreToolUse` can intercept shell commands and deny supported calls.
The hook command uses the session directory; the documentation recommends resolving repository hooks from the Git root.

Source: [Hooks documentation](https://learn.chatgpt.com/docs/hooks), retrieved 2026-10-02.

Bench already declares startup, stop, and command hooks through repository configuration.
Its command hook matcher is `Bash`.
These declarations support a common repository contract, but each interface still needs evidence that the hooks run.

Source: `.codex/hooks.json:3`, `.codex/hooks.json:14`, and `.codex/hooks.json:25`.

The precise desktop trust-review route for each installed version remains unverified.
The documented CLI route is `/hooks`.
Bench must not fabricate trust or overwrite trust records to complete setup.
This proposed constraint follows the approved security boundary in [ticket #4](../tickets/4.md).

## Q3: Recovery facts and limits

The official troubleshooting guide recommends checking approvals and testing a basic terminal command.
It also recommends a new focused chat for a stuck chat.
For persistent terminal trouble, it recommends waiting for active chats before restarting the app.
These are diagnostic routes, not documented repairs for this exact process-creation error.

Source: [Troubleshooting documentation](https://learn.chatgpt.com/docs/reference/troubleshooting), retrieved 2026-10-02.

The inspected documentation provides no verified repair for the absent launcher path or the Node URI rejection.
This bounded finding does not assert that no upstream repair exists.
An app restart remains untested against this incident.
The actor that removed the path remains unknown.

A Bench process cannot report its own diagnostic result when the harness never starts that process.
Therefore, a useful recovery route must remain accessible outside the failed command path.
This inference joins the documented environment separation with the [actual failed calls](installed-compatibility-probe.md).

| Route | Consequence | Evidence |
|---|---|---|
| Shared repository integration with separate user homes | Preserves the approved personal-settings boundary; requires interface-specific verification | Q1 and Q2 sources |
| Share the entire Codex home | Also shares state outside the approved Bench boundary | Windows app documentation; ticket #6 |
| Approved diagnostic shell or fresh process | Can inspect evidence; cannot certify recovery of the failed chat | Installed probe, diagnostic controls |
| Supported app recovery followed by a live retest | Can establish recovery if the actual failed interface passes | Troubleshooting documentation; retest remains unperformed |
| Edit private launcher files | Has no verified supported repair contract in this research | Inspected documentation; ticket #8 authority boundary |

## Q4: Current Bench facts

Bench's startup hook is informational and exits successfully after its diagnostic attempt.
It resolves Bench, prints command-resolution guidance, and invokes `session-inspect`.
It does not invoke the model's shell tool to prove that the tool can start.

Source: `.bench/hooks/session-start.sh:4`, `.bench/hooks/session-start.sh:40`, and `.bench/hooks/session-start.sh:52`.

Bench doctor checks the shim and repository adoption state.
Its registered rows include gate inputs, worktree administration, binary seals, and the promotion broker.
Those checks do not qualify a separate active desktop chat.

Source: `internal/adopt/doctor.go:198` and `internal/adopt/doctor_rows.go:29`.

The harness registry records one Codex row with declared hooks and a headless adapter.
It does not distinguish installed CLI and desktop execution health.
A declared capability and a passing live operation answer different questions.

Source: `internal/harnesses/harnesses.go:197`.

## Tested results

The [installed probe](installed-compatibility-probe.md) owns exact invocations, observed identities, results, and coverage gaps.
This research does not add a successful runtime compatibility claim.
The startup guidance received by this session coexists with a failed normal shell tool.
Therefore, a startup banner cannot serve as sufficient compatibility evidence.

## Proposals for reviewer decision

Ticket #12 asks the reviewer to decide these remaining boundaries:

1. Qualify capabilities per interface and invalidate relevant evidence when its runtime, configuration, or workspace changes.
2. Provide recovery guidance outside a failed command path and require a successful live retest.
3. Require desktop-specific tools only where an operation lacks a verified equivalent route.

These proposals preserve the complete workflow outcome and the existing repair authority.
They do not authorize reduced permissions or private-runtime edits.

## Unknowns and verification record

The exact desktop failure cause and its supported repair remain unknown.
The broader workflow, concurrent writers, hook enforcement, and recovery remain unqualified.
The current documentation may describe behavior absent from an installed build.
Version-specific live evidence must settle that gap.

- [x] The report opens with its recommendation, scope, and evidence status.
- [x] Facts, inferences, tested results, and proposals have explicit labels or sections.
- [x] Each factual question has a synthesized answer or an explicit unknown.
- [x] The route table states consequences and their sources.
- [x] The report preserves contradictions and residual unknowns.
- [x] The graph and tables expose the material relations; no additional diagram is needed.
- [x] Load-bearing claims cite primary URLs or exact local source lines.
- [x] The report records evidence invalidation and what remains unverified.

## Validation plan

1. Resolve ticket #12 before marking the map ready.
2. Assign acceptance coverage for each complete workflow outcome during specification.
3. Qualify the actual CLI and desktop interfaces with their required permissions and hooks.
4. Test fresh sessions, resumed sessions, configuration drift, concurrent writers, and failed-capability recovery.
5. Require a live retest before reporting that the affected operation recovered.
