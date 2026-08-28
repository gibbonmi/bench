# Collapse the artifacts job to one generation

Blocked by: switch-the-proof-matrix-to-proven-targets.md
Writes: .github/workflows/native-runtime.yml, .github/workflows/release.yml, internal/conformance, tests/canary/package-core-guard

## What to build

The `artifacts` job builds one generation and publishes one upload.

The job drops its second-source clone and its second build. It publishes one artifact whose `path` names `dist/artifacts` and `dist/reproducibility.json`, so the upload action resolves `dist` as the common root. The archive root therefore moves up one level, and every consumer downloads into `dist` instead of `dist/artifacts`. The consumers are the `evidence` and `smoke` jobs here, and the `authorize` and `publish` jobs in `release.yml`.

The `native-proof` job loses its second generation as well. It drops its second-source clone, its `-artifacts-second` download, its second `native-proof.sh` run, and its second proof upload. Its remaining download moves to `dist`.

The `evidence` job also drops its second clone, its second aggregate and preflight run, and its `compare-artifacts.sh` call. `dist/workflow-reproducibility.json` leaves the tree entirely.

`scripts/build-artifacts.sh` is unchanged. Its inner reproducibility build still writes `dist/reproducibility.json`, so the byte-for-byte artifact comparison survives and `reproducibility.builds` stays at 2.

`checkNativeRuntimeWorkflow` grades the new shape, and every check this ticket changes gets a mutation that turns it red. The canary fixtures under `package-core-guard` re-anchor on lines that survive the rewrite. Two fixtures are known to need it: `native-reproducibility-handoff-dropped` anchors the reproducibility upload path, and `preflight-native-call-bypassed` anchors a second-source preflight run that goes away.

## Acceptance

- [ ] The `artifacts` job contains no second-source clone (row A1).
- [ ] The conformance check names a diagnostic when a second-source clone returns (row A1).
- [ ] The `artifacts` job contains exactly one upload step (row A2).
- [ ] The single upload names `dist/reproducibility.json`, and a canary mutation of that path turns the guard red (row A3).
- [ ] `build-artifacts.sh` still writes `dist/reproducibility.json`, and release evidence still reads it (row A4).
- [ ] The `evidence` job contains no second-source clone (row A5).
- [ ] No tracked file under `.github`, `scripts`, `internal`, `tests`, or `docs` names `workflow-reproducibility.json` (row A6).
- [ ] The `evidence` job still runs `release-preflight.sh --mode verify` (row A7).
- [ ] The `smoke` job downloads the one artifact into `dist` (row A8).
- [ ] The `authorize` and `publish` jobs in `release.yml` download the one artifact into `dist` (row A9).
- [ ] The `native-proof` job contains no second-source clone (row A10).
- [ ] The `native-proof` job downloads the one artifact into `dist` (row A11).
- [ ] Each changed conformance check names its own diagnostic (row C1).
- [ ] Every mutation in the changed checks turns its check red (row C2).
- [ ] Every canary fixture under `package-core-guard` still bites (row C3).
