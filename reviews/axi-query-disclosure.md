# AXI query disclosure implementation review

Candidate: `c2542a9765f91d3581ae88cf6d37d1da3ea0d2ee..bd4b184f65411e8cbc8251d9167a9a4dd2310a34`, plus the generated `capture/session-handoff.md` refresh.

Verdict: **REJECT**. The candidate has seven unique actionable findings. Accepted product decisions remain closed; every finding is a deterministic repair inside the approved specification.

## Standards

Finding count: 2. Worst issue: Major.

### Major — duplicate production-registry AST derivation

Disposition: `auto-fix`.

`internal/conformance/axi_query_registry_test.go:261` independently locates and walks the `commandRegistry` literal even though `internal/conformance/subcommand_routing_test.go:157` already owns that route. The two walkers already differ on accepted declaration shapes. This violates the repository's one-source rule and the spec's requirement that AXI membership use the same AST route as subcommand routing (`specs/axi-query-disclosure/spec.md:31`). Extract one registry-literal seam, preserve both checks' distinct expectations, and keep both-direction registry mutations red.

### Major — duplicated TOON decode and POSIX argv probe harness

Disposition: `auto-fix`.

`internal/axi/action_test.go:88`, `internal/coverage/coverage_test.go:385`, and `internal/worktree/list_actions_test.go:172` paste the same official-TOON decode plus `sh -c` argv recovery derivation, including the same positional-index assumption. The repository standard explicitly treats a fixture harness pasted across packages as duplicated knowledge. Move the shared decode-and-execute mechanism behind one test-support seam while retaining independent surface expectations and exact byte comparisons.

## Spec

Finding count: 3. Worst issue: Critical.

### Critical — QD2 empty-token fixture does not contain the pre-change bytes

Disposition: `auto-fix`.

`internal/worktree/testdata/pre-disclosure-argv-pairs.json:12` claims the old response was `unknown argument: \"\"`. At base `c2542a97`, `internal/worktree/list.go:21` passed the empty string directly to `toon.Usage`, whose interpolation produced `usage: bench worktree list (unknown argument: )\n`. The accepted QD2 rule requires every non-help argv row to preserve old stdout, stderr, and exit (`specs/axi-query-disclosure/spec.md:29`, `specs/axi-query-disclosure/spec.md:58`). Correct the captured old bytes and make the grammar-scoped parse path reproduce them; do not add a fourth accepted delta.

### Major — honest-empty fallback omits non-control unrepresentable values

Disposition: `auto-fix`.

`internal/axi/action.go:106` recognizes only unsupported control runes before rendering, while `internal/axi/action.go:243` also rejects empty known arguments and values containing `<` or `>`. Such a coverage or worktree path reaches `RenderHelp`, returns an error, replaces the computed primary response, and changes the exit instead of appending honest empty help. QD1 explicitly requires primary-response and exit preservation for any unrepresentable known argv or `why` cell, including coverage and worktree paths (`specs/axi-query-disclosure/spec.md:30`, `specs/axi-query-disclosure/spec.md:57`). Make the fallback predicate cover every renderer refusal and prove it through public surfaces.

### Medium — QD1 public oracle independently re-sorts producer rows

Disposition: `auto-fix`.

`internal/worktree/list_actions_test.go:137` compares the two generated assignment IDs and swaps the fixture values into lexical order before building the expected public response. QD1 requires intent-ledger serialized ID order followed by Git registration producer order, and explicitly forbids tests from independently re-sorting it (`specs/axi-query-disclosure/spec.md:27`, `specs/axi-query-disclosure/spec.md:57`). A production change that re-sorts assignments can therefore stay green. Consume the assignment producer's order directly and prove a production re-sort makes the public oracle red.

## Coverage

Finding count: 2. Worst issue: Medium.

### Medium — guards enumeration-timeout lacks a public old/new pair

Disposition: `auto-fix`.

QD6 separately enumerates partial-timeout and enumeration-timeout incomplete scans (`specs/axi-query-disclosure/spec.md:39`, `specs/axi-query-disclosure/spec.md:62`). `internal/guards/guards_test.go:111` drives only partial timeout through `Command`; the enumeration-timeout test at `internal/guards/guards_test.go:167` asserts internal `Scan` metadata. A change to enumeration-timeout primary bytes, exit, or help remains green. Add a checked-in public-command old/new fixture that pins bytes, streams, exit, and honest empty help.

### Medium — worktree non-actionable terminal rows lack a public old/new pair

Disposition: `auto-fix`.

QD6 separately enumerates worktree active, orphaned, empty, and non-actionable terminal states (`specs/axi-query-disclosure/spec.md:39`, `specs/axi-query-disclosure/spec.md:62`). Public fixtures cover empty and active/orphaned output, while terminal complete and present-foreign rows appear only in the direct action-derivation test at `internal/worktree/list_actions_test.go:95`. A terminal-row primary-output regression remains green. Add a checked-in public-command old/new fixture that pins the terminal primary response, exit, and honest empty help.
