# Swap the release workflow to bench release submit and flip the conformance contract

Blocked by: wire-adapter-selection.md
Writes: .github/workflows/release.yml, internal/conformance/workflow_checks_test.go, internal/conformance/native_workflow_test.go, internal/conformance/registry_test.go, tests/canary/package-core-guard/preflight-publish-order-bypassed/, decisions/assets/gate-pipeline-fixture-inventory.md, docs/release-runbook.md

## What to build

The publish job of `.github/workflows/release.yml` becomes one Bench
invocation: check out the tag, set up Go and Node, build `dist/bench` from the
checkout, download `release-artifacts` into `dist/artifacts` and
`publish-preflight-evidence` into `dist/preflight`, run `dist/bench release
submit --version "${GITHUB_REF_NAME#v}" --profile public --path first
--adapter npm --provenance --registry https://registry.npmjs.org` with
`env: NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}` on the submit step, and
upload `dist/publication/` as the `publication-record` artifact on
`if: always()`. Both raw `npm publish` steps are gone; no promote appears
anywhere. `checkReleaseWorkflow` flips in the same slice: raw `npm publish`
in the text is a diagnostic, the submit invocation anchor (including
`--adapter npm` and `--provenance`) is required, the preflight-evidence
download and publication-record upload anchors are required and scoped to
`workflowJob(text, "publish")`, `release promote` in the text is a
diagnostic, and a mutation bite test proves each new diagnostic in the
`TestNativeWorkflowEvidenceEdgeBites` style. Two step-name contracts retire
in the same slice, because their anchor bytes vanish with the raw steps: the
platform-first/wrapper-last diagnostic that indexes the publish step names
(its concern now lives in T1's record-level ordering assertion), and the
`preflight-publish-order-bypassed` canary fixture with its registration and
its fixture-inventory row. The runbook's first-publication section names the
exact CI submit invocation and the tag-push presence rule (only the reviewer
pushes a release tag; the CI submit is the attended act's mechanical arm).
The workflow edit and the check flip land together because each is red
against the other's old state.

## Acceptance

- [ ] The bite test reds a workflow containing `npm publish` and one missing the submit anchor (spec rows R8, R9).
- [ ] The bite test reds a publish job whose `workflowJob` slice lacks the `publish-preflight-evidence` download or the `publication-record` upload, and stays red even when the authorize job carries the same artifact name (R10, R11).
- [ ] The bite test reds a workflow containing `release promote` (R12).
- [ ] `TestRootConformance` is green against the live tree with the swapped workflow; the retired step-name diagnostic, canary fixture, registration, and inventory row are gone; `bench canary` still passes.
- [ ] `docs/release-runbook.md` names the exact CI submit invocation and the tag-push presence rule.
