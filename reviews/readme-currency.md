# README currency review pickup

Status: coordinator source and diagram review accepted; landing pending
Base: 5a7fb32f2d9d19cbd093f6f6ad386e93002c7884
Assignment: readme-currency
Author: /root/readme_author
Line: gpt-5.6-sol / high
Pre-review attempts consumed: 2 of 3
Post-review repair cycles consumed: 2 of 2
Hardening cycles consumed: 0 of 1
Diagram composition attempts consumed: 2 of 2
Expected repair rounds: 1
Confidence: 6

## Result

The README now introduces Bench through the maintenance problem that first-time users face with AI-written code.
It separates a chosen seam from blast-radius evidence and shows the CLI commands that support safe changes.
The same pass corrects stale workflow, setup, installation, maintenance, example-profile, and shift-note guidance.
Four Mermaid diagrams make research and specification, chunk review and repair,
final reconciliation and landing, and the capture loop visible.

## Source-to-correction audit

| Subject | Primary source | README correction |
|---|---|---|
| Seam and blast radius | `CONTEXT.md`; `bench-craft-seams`; `bench help`; `internal/outline`; `internal/consumers` | Define a seam, add a concrete example, and state the static-Go limits beside `bench consumers`. |
| Ticket ownership | `.agents/commands/bench-write-spec.md`; `.agents/commands/bench-implement-spec.md` | Say that spec authoring slices the tickets before approval and implementation consumes the approved graph. |
| Full and delegated runs | `.agents/commands/bench-implement-spec.md` | Link the canonical `--full "<spec>"` and opt-in `--delegate` guidance. |
| Setup | `.agents/commands/bench-setup-repo.md`; `internal/adopt/setup.go`; `internal/adopt/init.go` | Use transactional `bench setup` as the default and retain link and init as low-level primitives. |
| Source requirements | `package.json`; `go.mod`; `internal/adopt/setup.go` | Point to the version owners and require a Git repository. |
| Source status | `ROADMAP.md`; npm registry response supplied by the coordinator | Require an immutable source ref and link the canonical NO-GO status. |
| Global shim | `bin/bench-postinstall.sh`; `internal/adopt/doctor.go` | Describe shim installation as best effort and route repair through `bench doctor`. |
| Upgrade route | `.bench/BENCH-reference.md`; `.bench/consumer-payload.json` | Give linked repositories `bench upgrade` and reserve update-kit for kit maintainers. |
| Example profile | `projects/benchkit.md`; `projects/gl-axi.md` | Identify gl-axi as a shipped example instead of live integration evidence. |
| Shift notes | `internal/shift/session.go` | Name `.bench-notes.md` instead of `notes.md`. |
| Workflow diagrams | `.bench/BENCH.md`; phase commands; `craft-research` | Show phase ownership, bounded repair, landing authority, and the reviewed capture loop. |

The coordinator confirmed that `README.md` ships without `ROADMAP.md`.
The release-status link therefore targets the canonical repository page.

## Author verification

| Command after `worktree exec readme-currency --` | Result | Package or wall time |
|---|---|---|
| `bench gate-prose . -- README.md` | Passed. | 60 ms wall |
| `bench gate-prose . -- specs/readme-currency/tickets/update-readme.md` | Passed. | 60 ms wall |
| `bench test --check docs-currency-workflow` | Passed; no failures or skips. | 603 ms package; 3.33 s wall |
| `bench test --check load-validity-metadata` | Passed; no failures or skips. | 104 ms package; 1.79 s wall |
| `bench test --check prose-mechanics` | Passed; no failures or skips. | 339 ms package; 1.89 s wall |
| `bench test --check retro-improvement-markers` | Passed; no failures or skips. | 4 ms package; 1.46 s wall |
| `bench help` | Confirmed every advertised CLI spelling. | 226 ms wall |
| `git diff --check` | Passed. | 138 ms wall |

The first two named-check calls could not write the sandboxed shared Go build cache.
The authorized reruns reached the checks and produced the passing results above.
Local targets for the package, Go module, implementation command, and example profile all exist.
The required README anchors remain present.

## Claims

Claim rows carry no free-text field.

```text
claims[7]{row,status,confidence}:
  maintenance-opening,claimed,9
  workflow-ownership,claimed,9
  setup-and-upgrade,claimed,9
  command-and-link-spellings,verified,10
  package-unpublished,claimed,8
  workflow-diagram-semantics,claimed,9
  workflow-diagram-rendering,verified,10
```

## Independent semantic review

The coordinator found two blocking README details after the author candidate.
The dirty-checkout wording could imply that `bench consumers --changed` runs on uncommitted edits.
The duplicated-pricing example could imply that the command detects copied rules.

Repair cycle 1 states the clean-checkout requirement and starts from one shared pricing function.
It also changes semantic source-inspection claims from `verified` to `claimed`.
The coordinator accepted the repaired source with no blocking finding.

## User-expanded diagram scope

The user expanded the accepted candidate before commit.
This expansion adds visual workflows and consumes no defect-repair cycle.
The coordinator audited the exact Mermaid candidate against the canonical phase
commands and successfully rendered all four diagrams before commit.
The accepted README SHA-256 is
`2bb83748657690716af64ae2bb9d1392647ee2afba4cf5a6794dbd9853faec69`.

## Source-lane record repair

The first source-commit attempt published no commit because
`retro-improvement-markers` rejected three improvement entries without a
destination marker. Repair cycle 2 adds `Feeds: none` to those entries while
leaving the accepted README bytes unchanged.

The coordinator owns independent review, commit, the whole-project gate, and landing.
