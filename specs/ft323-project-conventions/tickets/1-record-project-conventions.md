# Record the project prose and host conventions

Blocked by: none
Writes: AGENTS.md, projects/benchkit.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/anchor_harness_diagnostics_test.go, internal/anchors/registry_ft311_review_dispatch.go, specs/ft323-project-conventions/tickets/1-record-project-conventions.md (new)
Covers: none

## What to build

Record the three conventions approved in roadmap/FT323.md. Keep each fact in its project-owned source.
Use the existing prose APIs and the existing named-file prose lane. Record the observed host limits from retained evidence.
This ticket adds no executable behavior and crosses no declared seam. It uses the approved light path.

## Acceptance

- [x] A Markdown check reads internal/prose and uses its shared sentence and paragraph rules.
- [x] The author runs bench gate-prose on each edited Markdown file before the commit.
- [x] The profile records the current Go parallelism pins and the desktop opener that hangs on this host.

## Verification

Read the prose owner and every surface that these conventions steer.
Run `bench gate-prose . -- AGENTS.md projects/benchkit.md specs/ft323-project-conventions/tickets/1-record-project-conventions.md`.
Run bench test --check guidance-prose-budgets and bench test --check docs-currency-workflow.

Commit this ticket on a green lane. Obtain independent Standards, Spec, and Coverage reviews of the frozen source.
Land through bench worktree land with the tickets-only slug. The landing supplies the whole-project gate.

No behavioral mutation applies to this prose-only ticket. Existing conformance checks remain unchanged.

## Charge

Retained author: the current Codex session, gpt-6-astra, ultra effort.
Initial implementation is uncapped within the approved time window and closeout budget condition.
Post-review repairs permit two cycles. The consumed count starts at zero.
The build preflight requires spec.md and refuses this tickets-only light path.
Review uses an explicit frozen base and the ticket plus roadmap row as its approved source.

## Host evidence

The current Go environment file records `GOFLAGS=-p=4 -parallel=4`.
FT255 records the cap of two test-running delegates. Retained session evidence records `-parallel 2` for each delegate test run.
The same retained evidence names xdg-open as the hanging WSL opener. This run does not reproduce the hang.
The host probe establishes no GOMAXPROCS pin.

## Author evidence

The named-file prose lane passes all three edited Markdown files.
The guidance-prose-budgets check passes. The docs-currency-workflow check passes.

The author read prose.Paragraphs, prose.Grade, and GateProseCommand before the edit.
Bench consumers found the existing prose callers. No caller or executable changes.
The inherited structure report has 104 issues. This prose-only ticket adds no Go structural debt.

The host evidence comes from the Go environment file, FT255, and a retained session memory attachment.
The ticket applies craft-synthesis through its prose-only route. The landing gate completes that route.
