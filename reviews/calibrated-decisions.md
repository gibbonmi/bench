# Calibrated decisions

## CD1 ticket 1 author evidence

The ticket adds one `Claim schema` section to the delegation discipline reference. It adds one pointer sentence to the delegate skill. Seven anchors and seven omission canary fixtures grade the seven sentences.

### Scenario rows

These two rows show the claim schema shape the ticket describes. The coordinator writes the label cell.

| surface | claim | status | confidence | label |
| --- | --- | --- | --- | --- |
| delegate return | the named check runs green on the returned tree | claimed | 7 | held |
| delegate return | the guidance reads clearly for a cold reader | abstained |  |  |

The coordinator probes the first row's named check on the exact tree and labels it `held`. The coordinator labels nothing on the second row and counts one abstention.

### Done-claim table

The author does not write the label cell. Each label cell stays empty for the coordinator.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR1 | verified | 9 |  |
| CR2 | verified | 9 |  |
| CR3 | verified | 9 |  |
| CR13 | verified | 9 |  |
| CR14 | verified | 9 |  |
| CR32 | verified | 9 |  |
| CR34 | verified | 9 |  |
| CR37 | verified | 9 |  |

### Red-then-green log

Each row below uses `bench probe` with the omission kind against `docs-currency-workflow`. The verb records the green baseline, then the red under the mutation, then the proved restore. The row for CR14 uses a line-growth swap against `guidance-prose-budgets`.

| row | red line under the mutation | green |
| --- | --- | --- |
| CR1 | `gate: calibration: a done-claim row needs its status and stated confidence` | baseline passed, restored yes |
| CR2 | `gate: calibration: an unstatable confidence needs an abstained row` | baseline passed, restored yes |
| CR3 | `gate: calibration: the coordinator's tree probe labels each done-claim row` | baseline passed, restored yes |
| CR13 | `gate: calibration: the charge must name the claim schema section` | baseline passed, restored yes |
| CR14 | `gate: prose-budget exceeded: .agents/skills/bench-craft-delegate/SKILL.md is 128 lines, over its 126-line budget` | baseline passed, restored yes |
| CR32 | `gate: calibration: a claim must carry no free-text field` | baseline passed, restored yes |
| CR34 | `gate: calibration: verified and claimed need their evidence definitions` | baseline passed, restored yes |
| CR37 | `gate: calibration: an abstention must never become a refuted claim` | baseline passed, restored yes |

`TestEveryRetainedFixtureBitesThroughRegisteredOwner` proves each of the seven new canaries bites through its registered owner. The same run proves the 525 earlier fixtures keep their planted diagnostics.

### Probe verdict

The self-probe omits the CR3 sentence from the delegation discipline reference and runs `docs-currency-workflow`. The verdict line reads `bit`, and `restored` reads `yes`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | pass | 890 ms |
| `bench test --check guidance-prose-budgets` | pass | 4 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 15 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 9299 ms |
| `go vet ./...` | pass | no output |
| `bench gate-prose` on both edited files | pass | two pass rows |
