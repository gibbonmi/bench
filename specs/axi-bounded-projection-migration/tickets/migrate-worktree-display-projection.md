# Migrate worktree display projection

Blocked by: none
Ownership fence: `internal/worktree`
Integration surfaces: shared projection→implemented prerequisite `axi-carriers-and-registry`; final contraction→contract-projection-routes.md
Contracts: complete sorted inventory, fingerprint input bytes, visible entries, total/at-least facts, and refusal disposition cross `internal/worktree/subshell.go`→shared projection, membership is default 20 and full 1000 views, order is enumerate-validate-fingerprint-project, and absence is empty safe inventory, asserted by WI1
Closure: WI1/sorted, WI1/safety, WI1/fingerprint, WI1/default-20, WI1/full-1000, WI1/lower-bound, WI1/stat-race, WI1/route

## What to build

worktree projects only the final display after complete safety and fingerprint enumeration.

## Acceptance

- [ ] [WI1] (covers BP3) worktree projects only the final display after complete safety and fingerprint enumeration.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| WI1/sorted | project before sorting | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/safety | drop one unsafe entry from authority input | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/fingerprint | drop one entry from the fingerprint preimage | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/default-20 | change the default visible cap | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/full-1000 | make full uncapped | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/lower-bound | replace at-least with exact total | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/stat-race | turn stat-race refusal into success | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |
| WI1/route | bypass shared display projection | the independent owner route test | apply the subject mutation, run the named package boundary test and compatibility case under its timeout, and require the specific red |

