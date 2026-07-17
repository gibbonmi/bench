## Standards

5 findings. Worst: High.

1. High — the package-evidence allowlist is duplicated across `internal/preflight/requirements.json:4`, `scripts/build-release-evidence.mjs:72`, `internal/preflight/artifact_evidence.go:199`, and fixture/test lists. This violates AGENTS.md's one-source-per-fact rule; centralize the registry and derive consumers.
2. Medium — `scripts/build-artifacts.sh:26`, `:30`, and `:42` parse the platform matrix independently and `:47` separately derives its count. Parse once and reuse the rows/count.
3. High — `internal/conformance/native_workflow_test.go:81` greps a JSON tag instead of running preflight and inspecting generated evidence. Replace the tripwire with a real-path behavior check as required by `bench-craft-gate`.
4. Medium — `internal/conformance/native_workflow_test.go:31` and `:45` reuse diagnostics for absent versus malformed registries, while `:52` stops at the first incomplete record. Give distinct failure modes distinct messages and aggregate independent failures.
5. Low judgment — `internal/preflight/artifact_evidence.go:182` accepts an unused matrix and discards it at `:242`. Remove the speculative parameter or implement the intended validation.

## Spec

6 findings. Worst: P1.

1. P1 — focused publish authorizes without a profile. `specs/governed-release-evidence.md:99` and row 10 at `:219` say focused runs cannot authorize; `internal/preflight/command.go:149` permits `--mode publish --phase gate`, which traces to exit 0 and green focused evidence.
2. P1 — phase-red runs omit `release-index.json` and `SHA256SUMS`. The complete trustworthy-red contract at `specs/governed-release-evidence.md:131` and row 13 at `:223` is bypassed by `internal/preflight/release_evidence.go:75`.
3. P1 — cancellation replaces the prior complete generation despite `specs/governed-release-evidence.md:131`; the interrupted path also reaches the promotion at `internal/preflight/release_evidence.go:75`.
4. P1 — final artifact validation accepts malformed embedded evidence. `specs/governed-release-evidence.md:212` and `:246` require empty/invalid SPDX and policy inputs to be red, but `internal/preflight/artifact_evidence.go:199` checks only presence before accepting a manifest-consistent tarball.
5. P1 — requirement schemas are not closed. `specs/governed-release-evidence.md:112` requires unknown versions/fields to be red, but `internal/preflight/release_requirements.go:58` accepts arbitrary schema strings and `:230` decodes policy records into permissive maps.
6. P2 — the release index omits tracked `go.sum`, contrary to the dependency-input binding contract at `specs/governed-release-evidence.md:123`; `internal/preflight/release_requirements.go:33` includes `go.mod` only.

## Coverage

7 findings. Worst: P1.

1. P1 — add a built-CLI test proving focused publish is usage-red or non-authorizing; existing focused coverage at `internal/preflight/integration_test.go:69` exercises verify only.
2. P1 — add a full-run cancellation test proving prior-generation preservation; `internal/preflight/evidence_test.go:164` checks child termination only.
3. P1 — add a deterministic phase-red black-box test requiring complete release index/checksum/requirement evidence; `internal/preflight/preflight_test.go:187` checks phase records only.
4. P1 — add typed hostile governance cases for unknown fields/schema versions; fixtures at `internal/preflight/integration_fixture_test.go:34` seed valid records only.
5. P1 — add final-tar cases for empty evidence, unsafe modes, and malformed SPDX/policies with internally consistent manifests; `internal/contract/surface/artifact/artifact_test.go:76` covers only a source FIFO.
6. P2 — add bank-profile coverage where missing FT71 is red; `internal/preflight/preflight_test.go:97` removes only FT87/FT88 under public.
7. P2 — add black-box deterministic byte comparison, exact index/SHA256SUMS cross-check, synchronized input-drift, and rerun/prior-generation assertions required by rows 7, 9, 13, 14, and 15 (`specs/governed-release-evidence.md:216`, `:218`, `:222`, `:223`, `:224`).
