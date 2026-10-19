## Outcome

FT293 landed as 4d547d73b51e9f73edd49badfb0d7441bf0ba0d8 from the reviewed source pair 36b1fafc4..be01ed036. The landing gate was green, and the broker refresh and the sibling cleanup were complete. Build preflight predicts no ownership closure and checks no charge premises, and `--propose-writes` and `--breakdown` are gone. The ticket `Writes:` union plus `capture`, the spec folder, and the review pickup is the ownership fence. The landing and review preflight read that fence from the tracked tickets of the source tip. `bench spec widen` adds a path inside the directory of a signed-off entry, and `--decision` admits a path outside it. The guidance states one `Writes:` rule.

The build ran four chunks with six fresh ticket authors and eleven fresh repair sessions on opus at high effort. Every review axis ran on opus at high effort. The orchestrator ran `bench spec widen` itself seven times as dogfood, and the spec records each run.

## Gate-stage timings

- landing: commit 4d547d73b51e9f73edd49badfb0d7441bf0ba0d8, trace 91baa063d74494e05d2bfff178a34fb6
- gofmt: 180 ms
- vet: 1432 ms
- test: 166774 ms
- race: 2380 ms
- system: 57754 ms
- shellcheck: 536 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 01 | CR07-CR10 removed forms refuse and leave the help | verified | 9 | held | opus / high / implementation |
| ticket 01 | CR80-CR85, CR91 forbid anchors red each restored step | verified | 9 | held | opus / high / implementation |
| ticket 01 | CR93 changelog adds no `--breakdown` entry | claimed | 8 | held | opus / high / implementation |
| ticket 02 | CR06 the charge verdict holds no `charge-premises` row | verified | 9 | refuted | opus / high / implementation |
| ticket 02 | CR19-CR20 `ParseWrite` refuses every unclean or absolute entry | verified | 9 | refuted | opus / high / implementation |
| ticket 02 | CR21-CR22 `tickets-parse` reds an escaping entry | verified | 9 | held | opus / high / implementation |
| ticket 03 | CR23-CR34 widen adds, refuses, and admits by scope | verified | 9 | held | opus / high / implementation |
| ticket 03 | CR35 the CR35 case catches a root-parent rule | verified | 9 | refuted | opus / high / implementation |
| ticket 04 | CR55-CR59 the landing authorizes from ticket `Writes:` | verified | 9 | held | opus / high / implementation |
| ticket 04 | CR66-CR67 the audit scope reads every ticket at base and tip | verified | 9 | refuted | opus / high / implementation |
| ticket 05 | CR70 coverage ignores a fence section | verified | 9 | held | opus / high / implementation |
| ticket 05 | CR86-CR87 craft-spec anchors hold | verified | 9 | held | opus / high / implementation |
| ticket 06 | CR78 the refusal-boundary forbid anchor reds a restored sentence | verified | 9 | held | opus / high / implementation |
| ticket 06 | CR88 the operating guide states the widen scope rule | verified | 9 | refuted | opus / high / implementation |

Brier mean: 0.298 over 14 pairs. Abstentions: 0.

The ticket slices held their chunk boundaries, but four of six tickets needed `Writes:` widenings that the slicing pass did not predict. Ticket 01 needed four extra paths, ticket 02 needed six, ticket 04 needed ten, ticket 05 needed fifteen, and ticket 06 needed three. Each widening sat inside a signed-off directory, except `internal/maps/fields.go` and `decisions/architecture-planning.md`.

## Coordinator catches

- The first ticket 01 charge refused three times on the pre-FT293 `charge-premises` row: no planning inventory, missing coverage seams, and an unsealed worktree binary. The orchestrator added the inventory and built the worktree before the charge.
- The pre-FT293 closure rows redded the ticket 02 charge after ticket 01, so ticket 01 took two closure-only widenings.
- The ticket 02 author committed the CR1 freeze record before ticket 01 recorded, so ticket 01 could not stage a run at the chunk source. The orchestrator reset the worktree to the chunk source and kept the record uncommitted.
- A spec edit after a refreeze moved the chunk tip twice, so each such edit took a new refreeze.
- `bench record confirmation` refused a test-only repair, because its assignment plan commit changed the spec inside the delta. A confirming round replaced it.
- The CR3 review found that the landing fence read tickets from the worktree filesystem, so an ignored ticket-shaped file could authorize paths. The repair moved the landing and review preflight to the tracked tickets of the source tip.
- The complete checkpoint redded on a decision-map source that cited `roadmap/FT293.md`, which the landing deletes. A third ticket 06 repair moved the citation.
- Two consultations (fable medium, sonnet high) agreed to accept the refusal of a root headroom premise, so the opus tie-breaker did not run.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-delete-proposal-and-breakdown-forms.md | 1 | delegate-error |
| 02-delete-ownership-closure.md | 3 | spec-row, one-source, other |
| 03-add-spec-widen-verb.md | 2 | one-source, one-source |
| 04-derive-fence-from-tickets.md | 2 | spec-row, spec-row |
| 05-retire-spec-fence-section.md | 1 | delegate-error |
| 06-align-writes-guidance.md | 3 | delegate-error, delegate-error, ticket-slicing |

## Agent-experience improvements

### Bench CLI

- Census account for the landing (census=24): `rg` 16 and `sed` 1 have no Bench form, `cat` 2 maps to `bench worktree show`, `cd` 1 maps to `bench worktree exec`, and `cp` 4 maps to `bench worktree reset` with its envelope; raise the `bench worktree exec` spill threshold so that a ten-line response prints inline.
  Feeds: new
- Make `bench record ticket --id-prefix` skip a prefix that the requirement ID already starts with, so IDs such as `t2-t2-preflight-packages` stop.
  Feeds: new
- Let `bench record confirmation` treat a completion-plan assignment edit as a plan change, so a test-only repair closes without a confirming round.
  Feeds: new
- Print the first failure line of a failed test in `bench test` without `--full`, and print the total wall time of the run.
  Feeds: new

### Skills

- State in `/bench-implement-spec` that the last ticket of a multi-ticket chunk writes the freeze uncommitted, and each earlier ticket records at that HEAD before one record commit.
  Feeds: new
- Make the craft-tickets slicing checks list every caller of a changed fixture helper, every test that counts or pins a removed row, and every decision-map source that the delivery deletes.
  Feeds: new

### Process

- Make the spec stage of a spec that predates a removed charge row carry the planning inventory, so the first charge prepares without plan commits.
  Feeds: none
- Have the landing fence read only committed objects from the start of a design, because an authority that reads the working tree admits ignored files.
  Feeds: none