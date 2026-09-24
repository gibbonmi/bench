# 10. Apply the reviewed byte bound to each bounded response

Blocked by: 2-bound-every-public-response.md
Writes: internal/bounds/bounds.go, tests/canary/package-core-guard/bounds-duplicate-owner, internal/responsebound/ (new)
Covers: BO63, BO64, BO65

## What to build

Entry requires the measurement report and the reviewer budget decision that `specs/session-context-queries/tickets/4-apply-reviewed-budgets.md` owns. If either record is absent, stop before product writes and report the absent record. A spec-authoring pass records the approved value, or the per-surface values, and the record path in this ticket before it starts.

Add the approved value to the production policy registry in `internal/bounds`. The owner treats a response as over-bound when its bytes exceed the value, and its projection stays within the value. Choose no value in the build.

## Acceptance

- [ ] The ticket starts no product write while the decision record is absent.
- [ ] The byte value in the policy registry equals the value in the decision record.
- [ ] A response above the approved byte value spills, and its projection stays within the value.
