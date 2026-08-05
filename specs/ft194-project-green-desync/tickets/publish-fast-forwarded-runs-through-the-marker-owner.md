# Publish fast-forwarded runs through the marker owner

Blocked by: recognize-lagging-project-green-markers.md
Ownership fence: `internal/specbuild`
Contracts: the marker-advance request crosses `internal/specbuild`→`internal/gate/authorization` through `GateOwner`, asserted by LP1-LP4 against the real authorization owner; publication supplies destination and expected lineage only after prospective green or retained recovery evidence is proven, and a refusal returns before terminal run state is authored

## What to build

Both direct promotion publication and crash-recovery publication delegate marker advancement to `GateOwner`. Real-owner lifecycle fixtures reproduce the empty-run fast-forward wedge, the checkpointed recomposition wedge, and the branch-advanced recovery wedge, while preserving refusal, reload, idempotency, and the rule that fast-forward itself spends no gate and moves no marker.

## Acceptance

- [ ] [LP1] Starting empty at a project-green tip, advancing the branch by an ordinary gated commit, fast-forwarding, then completing assign→checkpoint→integrate→review→promote ends terminal with the branch and marker at the promotion commit; a fresh service reloads that terminal state.
- [ ] [LP2] A lagging-marker run with checkpointed work whose branch tip moves again recomposes without the marker-conflict refusal and completes after the required fresh review and promote retry.
- [ ] [LP3] Re-entering a promotion interrupted after branch advance completes publication through the real owner when the marker lags the run base.
- [ ] [LP4] A divergent marker planted after review refuses publication, leaves branch, marker, candidate, and run state recoverable, and succeeds only after the marker is restored.
- [ ] [LP5] Empty-run fast-forward continues to run no gate and move no project-green marker.
- [ ] [LP6] Existing recomposition-refusal and terminal replay controls remain mutation-free and idempotent.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LP1 | restore direct publication's compare-and-swap from `run.Base` | real-owner lifecycle regression test | drive Start, ordinary green advance, Checkpoint, Integrate, Review, and Promote, run the focused lifecycle test, expect publication to wedge after branch advance instead of reloading terminal |
| LP2 | keep recomposition Bootstrap equality-only or skip the post-recomposition review retry | real-owner recomposition regression test | integrate checkpointed work, advance the branch, promote to recompose, review the recomposed candidate, promote again, expect terminal completion without the marker-conflict string |
| LP3 | restore recovery publication's compare-and-swap from `run.Base` | real-owner promotion fault test | inject `promote/branch`, recreate the service, re-enter Promote, run the focused lifecycle test, expect terminal publication instead of a permanent compare-and-swap failure |
| LP4 | swap from whichever marker is present without checking its lineage | real-owner publication refusal test | plant a sibling marker after review, invoke Promote, snapshot public refs and status, expect refusal with no protected mutation; restore the marker and expect re-entry to remain possible |
| LP5 | move the marker or execute the gate inside `fastForwardEmptyRun` | fast-forward lifecycle control | snapshot marker and gate calls, fast-forward through Checkpoint or Start, run the focused fast-forward test, expect both counts and marker position unchanged |
| LP6 | finish terminal state twice or mutate before a recomposition refusal | existing promotion replay and refusal controls | run the focused promote-twice and recomposition-refusal tests, expect identical public state across the refused or replayed operation |
