# Upgrade the Node runtime to 24

Blocked by: none
Writes: .node-version (new), CHANGELOG.md, package.json, internal/publication/npm_registry.go, internal/publication/npm_registry_test.go, internal/conformance/node_runtime_policy_test.go (new), internal/conformance/package_core_checks_test.go, .github/workflows/release.yml, .github/workflows/native-runtime.yml, docs/release-runbook.md, tests/canary/compliance-hardening/mutable-workflow-action/MUTATE.json

## What to build

The `.node-version` file sets Node 24 as the supported runtime. The package, staged publication path, workflows, and release guidance project that value.

## Acceptance

- [x] The npm package rejects Node versions earlier than 24.
- [x] Staged publication rejects Node versions earlier than 24.
- [x] Release and native-runtime workflows read `.node-version`.
- [x] Release guidance and workflow mutation evidence name Node 24.
- [x] A conformance check rejects a stale Node runtime projection.
- [x] A 22.14 staged floor makes the below-floor Node test fail.
- [x] A 24.1 staged floor makes the exact-floor Node test fail.
