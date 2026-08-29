# Upgrade the Node runtime to 24

Blocked by: none
Writes: CHANGELOG.md, package.json, internal/publication/npm_registry.go, internal/publication/npm_registry_test.go, .github/workflows/release.yml, .github/workflows/native-runtime.yml, docs/release-runbook.md, tests/canary/compliance-hardening/mutable-workflow-action/MUTATE.json

## What to build

The package, staged publication path, and GitHub workflows use Node 24 as the minimum supported runtime. Release guidance and the workflow mutation fixture agree with that runtime.

## Acceptance

- [x] The npm package rejects Node versions earlier than 24.
- [x] Staged publication rejects Node versions earlier than 24.
- [x] Release and native-runtime workflows select Node 24.
- [x] Release guidance and workflow mutation evidence name Node 24.
